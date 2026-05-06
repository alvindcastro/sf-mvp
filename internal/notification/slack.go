package notification

import (
	"errors"
	"fmt"
	"strings"

	"sf-mvp/internal/approval"
	"sf-mvp/internal/brief"
	"sf-mvp/internal/observability"
	"sf-mvp/internal/severity"
)

type DeliveryMode string

const DeliveryModeDryRun DeliveryMode = "dry_run"

type PreviewStatus string

const (
	PreviewStatusBlocked PreviewStatus = "blocked"
	PreviewStatusAllowed PreviewStatus = "allowed"
)

var (
	ErrDryRunRequired       = errors.New("delivery_mode must be dry_run")
	ErrInvalidPreviewInput  = errors.New("invalid notification preview input")
	ErrNoRedactedBriefInput = errors.New("redacted brief is required")
)

type PreviewRequest struct {
	IncidentID   string
	Channel      string
	DeliveryMode DeliveryMode
	Brief        brief.Result
	Gate         *approval.Gate
	Recorder     *observability.Recorder
	Workflow     observability.Workflow
}

type PreviewResult struct {
	Status                   PreviewStatus
	DeliveryMode             DeliveryMode
	Reason                   string
	ApprovalRequestID        string
	PreparedPayload          SlackPayload
	Sent                     bool
	NetworkDeliveryAttempted bool
}

type SlackPayload struct {
	Channel string       `json:"channel"`
	Text    string       `json:"text"`
	Blocks  []SlackBlock `json:"blocks"`
}

type SlackBlock struct {
	Type string    `json:"type"`
	Text SlackText `json:"text"`
}

type SlackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func PreviewSlack(request PreviewRequest) (PreviewResult, error) {
	if request.DeliveryMode != DeliveryModeDryRun {
		return PreviewResult{}, ErrDryRunRequired
	}
	if err := validatePreviewRequest(request); err != nil {
		return PreviewResult{}, err
	}

	payload := slackPayloadFromBrief(request.Channel, request.Brief)
	gate := request.Gate
	if gate == nil {
		gate = approval.NewGate(nil)
	}

	execution, err := gate.Execute(approval.SensitiveActionCall{
		IncidentID: strings.TrimSpace(request.IncidentID),
		Action:     severity.SensitiveActionExternalSharing,
		Scope: approval.Scope{
			IncidentID: strings.TrimSpace(request.IncidentID),
			TargetRef:  SlackTargetRef(request.Channel),
		},
	}, func() error {
		return nil
	})

	result := PreviewResult{
		Status:                   PreviewStatusAllowed,
		DeliveryMode:             DeliveryModeDryRun,
		Reason:                   execution.Reason,
		ApprovalRequestID:        execution.RequestID,
		PreparedPayload:          payload,
		Sent:                     false,
		NetworkDeliveryAttempted: false,
	}
	success := true
	if err != nil || !execution.Allowed {
		result.Status = PreviewStatusBlocked
		success = false
		if strings.TrimSpace(result.Reason) == "" {
			result.Reason = "external sharing requires scoped human approval"
		}
	}

	recordPreviewToolCall(request, result, success)
	return result, nil
}

func SlackTargetRef(channel string) string {
	return "slack:" + strings.TrimSpace(channel)
}

func validatePreviewRequest(request PreviewRequest) error {
	incidentID := strings.TrimSpace(request.IncidentID)
	if incidentID == "" {
		return fmt.Errorf("%w: incident_id is required", ErrInvalidPreviewInput)
	}
	if strings.TrimSpace(request.Channel) == "" {
		return fmt.Errorf("%w: channel is required", ErrInvalidPreviewInput)
	}
	if request.Brief.Status == "" || len(request.Brief.Sections) == 0 {
		return ErrNoRedactedBriefInput
	}
	if strings.TrimSpace(request.Brief.IncidentID) != "" && strings.TrimSpace(request.Brief.IncidentID) != incidentID {
		return fmt.Errorf("%w: brief incident_id does not match requested incident", ErrInvalidPreviewInput)
	}
	return nil
}

func slackPayloadFromBrief(channel string, redacted brief.Result) SlackPayload {
	blocks := make([]SlackBlock, 0, len(redacted.Sections))
	for _, section := range redacted.Sections {
		title := strings.TrimSpace(section.Title)
		text := strings.TrimSpace(section.Text)
		if title == "" || text == "" {
			continue
		}
		blocks = append(blocks, SlackBlock{
			Type: "section",
			Text: SlackText{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*%s*\n%s", title, text),
			},
		})
	}

	return SlackPayload{
		Channel: strings.TrimSpace(channel),
		Text:    fallbackTextFromBrief(redacted),
		Blocks:  blocks,
	}
}

func fallbackTextFromBrief(redacted brief.Result) string {
	incidentID := strings.TrimSpace(redacted.IncidentID)
	for _, section := range redacted.Sections {
		text := strings.TrimSpace(section.Text)
		if text == "" {
			continue
		}
		if incidentID == "" {
			return truncate(text, 180)
		}
		return truncate("Dry-run incident brief preview for "+incidentID+": "+text, 220)
	}
	if incidentID == "" {
		return "Dry-run incident brief preview"
	}
	return "Dry-run incident brief preview for " + incidentID
}

func truncate(text string, max int) string {
	if len(text) <= max {
		return text
	}
	if max <= 3 {
		return text[:max]
	}
	return strings.TrimSpace(text[:max-3]) + "..."
}

func recordPreviewToolCall(request PreviewRequest, result PreviewResult, success bool) {
	if request.Recorder == nil {
		return
	}
	request.Recorder.RecordToolCall(request.Workflow, observability.ToolCall{
		Name:    "notification.slack.preview",
		Success: success,
		Fields: map[string]string{
			"status":                     string(result.Status),
			"delivery_mode":              string(result.DeliveryMode),
			"channel":                    strings.TrimSpace(request.Channel),
			"target_ref":                 SlackTargetRef(request.Channel),
			"approval_request_id":        result.ApprovalRequestID,
			"network_delivery_attempted": fmt.Sprintf("%t", result.NetworkDeliveryAttempted),
			"sent":                       fmt.Sprintf("%t", result.Sent),
			"payload_block_count":        fmt.Sprintf("%d", len(result.PreparedPayload.Blocks)),
		},
	})
}
