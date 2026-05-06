package notification

import (
	"errors"
	"strings"
	"testing"
	"time"

	"sf-mvp/internal/approval"
	"sf-mvp/internal/brief"
	"sf-mvp/internal/observability"
	"sf-mvp/internal/severity"
)

func TestDryRunSlackPreviewBlocksWithoutApprovedExternalSharingScope(t *testing.T) {
	gate := approval.NewGate(fixedNow)
	recorder, workflow := newPreviewRecorder(t, "FIC-SYN-001", observability.SensitiveData{Terms: []string{"BUS-007", "Depot 9"}})

	result, err := PreviewSlack(PreviewRequest{
		IncidentID:   "FIC-SYN-001",
		Channel:      "#fleet-safety",
		DeliveryMode: DeliveryModeDryRun,
		Brief:        redactedBrief(),
		Gate:         gate,
		Recorder:     recorder,
		Workflow:     workflow,
	})
	if err != nil {
		t.Fatalf("PreviewSlack returned error: %v", err)
	}

	if result.Status != PreviewStatusBlocked {
		t.Fatalf("Status = %q, want %q", result.Status, PreviewStatusBlocked)
	}
	if !strings.Contains(result.Reason, "approval") {
		t.Fatalf("Reason = %q, want approval explanation", result.Reason)
	}
	if result.PreparedPayload.Channel != "#fleet-safety" {
		t.Fatalf("payload channel = %q, want #fleet-safety", result.PreparedPayload.Channel)
	}
	if result.PreparedPayload.Text == "" || len(result.PreparedPayload.Blocks) == 0 {
		t.Fatalf("payload = %#v, want Slack-shaped text and blocks", result.PreparedPayload)
	}
	assertContains(t, slackPayloadText(result.PreparedPayload), "[redacted vehicle]")
	assertNotContains(t, slackPayloadText(result.PreparedPayload), "BUS-007")
	assertNotContains(t, slackPayloadText(result.PreparedPayload), "Depot 9")
	if result.NetworkDeliveryAttempted || result.Sent {
		t.Fatalf("result = %#v, want dry-run preview without network delivery", result)
	}
	assertToolCallEvent(t, recorder.Events(), "notification.slack.preview", "blocked", "false")
}

func TestDryRunSlackPreviewRequiresDryRunDeliveryMode(t *testing.T) {
	recorder, workflow := newPreviewRecorder(t, "FIC-SYN-001", observability.SensitiveData{})

	result, err := PreviewSlack(PreviewRequest{
		IncidentID:   "FIC-SYN-001",
		Channel:      "#fleet-safety",
		DeliveryMode: "send",
		Brief:        redactedBrief(),
		Gate:         approval.NewGate(fixedNow),
		Recorder:     recorder,
		Workflow:     workflow,
	})

	if !errors.Is(err, ErrDryRunRequired) {
		t.Fatalf("error = %v, want ErrDryRunRequired", err)
	}
	if result.PreparedPayload.Text != "" || result.NetworkDeliveryAttempted || result.Sent {
		t.Fatalf("result = %#v, want no prepared or sent payload", result)
	}
}

func TestDryRunSlackPreviewBlocksDeniedApproval(t *testing.T) {
	gate := approval.NewGate(fixedNow)
	request := createExternalSharingRequest(t, gate, "FIC-SYN-002", "#fleet-safety")
	_, err := gate.Decide(approval.DecisionInput{
		RequestID: request.ID,
		Approver:  "fleet-safety-lead",
		Decision:  approval.DecisionDenied,
		Reason:    "brief needs another review",
	})
	if err != nil {
		t.Fatalf("Decide returned error: %v", err)
	}
	recorder, workflow := newPreviewRecorder(t, "FIC-SYN-002", observability.SensitiveData{})

	result, err := PreviewSlack(PreviewRequest{
		IncidentID:   "FIC-SYN-002",
		Channel:      "#fleet-safety",
		DeliveryMode: DeliveryModeDryRun,
		Brief:        redactedBriefFor("FIC-SYN-002"),
		Gate:         gate,
		Recorder:     recorder,
		Workflow:     workflow,
	})
	if err != nil {
		t.Fatalf("PreviewSlack returned error: %v", err)
	}

	if result.Status != PreviewStatusBlocked {
		t.Fatalf("Status = %q, want %q", result.Status, PreviewStatusBlocked)
	}
	if !strings.Contains(result.Reason, "denied") {
		t.Fatalf("Reason = %q, want denied approval explanation", result.Reason)
	}
	if result.ApprovalRequestID != request.ID {
		t.Fatalf("ApprovalRequestID = %q, want %q", result.ApprovalRequestID, request.ID)
	}
	if result.NetworkDeliveryAttempted || result.Sent {
		t.Fatalf("result = %#v, want no network delivery", result)
	}
	assertToolCallEvent(t, recorder.Events(), "notification.slack.preview", "blocked", "false")
}

func TestDryRunSlackPreviewBlocksOutOfScopeApproval(t *testing.T) {
	gate := approval.NewGate(fixedNow)
	request := createExternalSharingRequest(t, gate, "FIC-SYN-003", "#other-channel")
	_, err := gate.Decide(approval.DecisionInput{
		RequestID: request.ID,
		Approver:  "fleet-safety-lead",
		Decision:  approval.DecisionApproved,
		Reason:    "approve only the requested target channel",
	})
	if err != nil {
		t.Fatalf("Decide returned error: %v", err)
	}
	recorder, workflow := newPreviewRecorder(t, "FIC-SYN-003", observability.SensitiveData{})

	result, err := PreviewSlack(PreviewRequest{
		IncidentID:   "FIC-SYN-003",
		Channel:      "#fleet-safety",
		DeliveryMode: DeliveryModeDryRun,
		Brief:        redactedBriefFor("FIC-SYN-003"),
		Gate:         gate,
		Recorder:     recorder,
		Workflow:     workflow,
	})
	if err != nil {
		t.Fatalf("PreviewSlack returned error: %v", err)
	}

	if result.Status != PreviewStatusBlocked {
		t.Fatalf("Status = %q, want %q", result.Status, PreviewStatusBlocked)
	}
	if !strings.Contains(result.Reason, "scope") {
		t.Fatalf("Reason = %q, want scope explanation", result.Reason)
	}
	if result.NetworkDeliveryAttempted || result.Sent {
		t.Fatalf("result = %#v, want no network delivery", result)
	}
}

func TestDryRunSlackPreviewAllowsScopedApprovalWithoutNetworkDelivery(t *testing.T) {
	gate := approval.NewGate(fixedNow)
	request := createExternalSharingRequest(t, gate, "FIC-SYN-004", "#fleet-safety")
	_, err := gate.Decide(approval.DecisionInput{
		RequestID: request.ID,
		Approver:  "fleet-safety-lead",
		Decision:  approval.DecisionApproved,
		Reason:    "redacted brief is approved for dry-run preview only",
	})
	if err != nil {
		t.Fatalf("Decide returned error: %v", err)
	}
	recorder, workflow := newPreviewRecorder(t, "FIC-SYN-004", observability.SensitiveData{})

	result, err := PreviewSlack(PreviewRequest{
		IncidentID:   "FIC-SYN-004",
		Channel:      "#fleet-safety",
		DeliveryMode: DeliveryModeDryRun,
		Brief:        redactedBriefFor("FIC-SYN-004"),
		Gate:         gate,
		Recorder:     recorder,
		Workflow:     workflow,
	})
	if err != nil {
		t.Fatalf("PreviewSlack returned error: %v", err)
	}

	if result.Status != PreviewStatusAllowed {
		t.Fatalf("Status = %q, want %q", result.Status, PreviewStatusAllowed)
	}
	if result.ApprovalRequestID != request.ID {
		t.Fatalf("ApprovalRequestID = %q, want %q", result.ApprovalRequestID, request.ID)
	}
	if result.DeliveryMode != DeliveryModeDryRun {
		t.Fatalf("DeliveryMode = %q, want %q", result.DeliveryMode, DeliveryModeDryRun)
	}
	if result.NetworkDeliveryAttempted || result.Sent {
		t.Fatalf("result = %#v, want allowed dry-run preview without network delivery", result)
	}
	assertToolCallEvent(t, recorder.Events(), "notification.slack.preview", "allowed", "true")
}

func TestDryRunSlackPreviewRecordsRedactedToolCallEvent(t *testing.T) {
	gate := approval.NewGate(fixedNow)
	recorder, workflow := newPreviewRecorder(t, "FIC-SYN-005", observability.SensitiveData{Terms: []string{"#fleet-safety"}})

	_, err := PreviewSlack(PreviewRequest{
		IncidentID:   "FIC-SYN-005",
		Channel:      "#fleet-safety",
		DeliveryMode: DeliveryModeDryRun,
		Brief:        redactedBriefFor("FIC-SYN-005"),
		Gate:         gate,
		Recorder:     recorder,
		Workflow:     workflow,
	})
	if err != nil {
		t.Fatalf("PreviewSlack returned error: %v", err)
	}

	events := recorder.Events()
	event := findToolCallEvent(t, events, "notification.slack.preview")
	if event.Fields["channel"] != observability.RedactedValue {
		t.Fatalf("channel field = %q, want redacted value", event.Fields["channel"])
	}
	if _, ok := event.Fields["payload_text"]; ok {
		t.Fatalf("tool-call event included payload_text field: %#v", event.Fields)
	}
}

func createExternalSharingRequest(t *testing.T, gate *approval.Gate, incidentID, channel string) approval.Request {
	t.Helper()

	request, err := gate.CreateRequest(approval.RequestInput{
		IncidentID: incidentID,
		Action:     severity.SensitiveActionExternalSharing,
		Scope: approval.Scope{
			IncidentID: incidentID,
			TargetRef:  SlackTargetRef(channel),
		},
		Reason: "operator requested dry-run Slack-shaped notification preview",
	})
	if err != nil {
		t.Fatalf("CreateRequest returned error: %v", err)
	}
	return request
}

func redactedBrief() brief.Result {
	return redactedBriefFor("FIC-SYN-001")
}

func redactedBriefFor(incidentID string) brief.Result {
	return brief.Result{
		Status:          brief.StatusDraft,
		IncidentID:      incidentID,
		SyntheticRecord: true,
		Shareable:       true,
		Sections: []brief.Section{
			{
				Title: "Incident Summary",
				Text:  "Synthetic draft for [redacted vehicle] near [redacted location].",
				Sources: []brief.Source{
					{Ref: "packet.incident_id", Type: brief.SourceTypePacket},
				},
			},
			{
				Title: "Recommended Actions",
				Text:  "Keep external sharing blocked until human approval exists.",
				Sources: []brief.Source{
					{Ref: "severity.approval_requirements[2]", Type: brief.SourceTypeSeverity},
				},
			},
		},
		RedactionsApplied: []brief.Redaction{
			{Field: "packet.vehicle_id", Reason: "vehicle identifier is not included in shareable drafts"},
		},
	}
}

func newPreviewRecorder(t *testing.T, incidentID string, sensitive observability.SensitiveData) (*observability.Recorder, observability.Workflow) {
	t.Helper()

	recorder := observability.NewRecorder(fixedNow, observability.Budget{})
	workflow, err := recorder.StartWorkflow(incidentID, sensitive)
	if err != nil {
		t.Fatalf("StartWorkflow returned error: %v", err)
	}
	return recorder, workflow
}

func fixedNow() time.Time {
	return time.Date(2026, time.May, 6, 16, 0, 0, 0, time.UTC)
}

func slackPayloadText(payload SlackPayload) string {
	var parts []string
	parts = append(parts, payload.Text)
	for _, block := range payload.Blocks {
		parts = append(parts, block.Text.Text)
	}
	return strings.Join(parts, "\n")
}

func assertContains(t *testing.T, text, want string) {
	t.Helper()
	if !strings.Contains(text, want) {
		t.Fatalf("%q not found in:\n%s", want, text)
	}
}

func assertNotContains(t *testing.T, text, unwanted string) {
	t.Helper()
	if strings.Contains(text, unwanted) {
		t.Fatalf("unexpected %q found in:\n%s", unwanted, text)
	}
}

func assertToolCallEvent(t *testing.T, events []observability.Event, toolName, status, success string) {
	t.Helper()
	event := findToolCallEvent(t, events, toolName)
	if event.Fields["status"] != status {
		t.Fatalf("tool status = %q, want %q in %#v", event.Fields["status"], status, event.Fields)
	}
	if event.Fields["success"] != success {
		t.Fatalf("tool success = %q, want %q in %#v", event.Fields["success"], success, event.Fields)
	}
}

func findToolCallEvent(t *testing.T, events []observability.Event, toolName string) observability.Event {
	t.Helper()
	for i := len(events) - 1; i >= 0; i-- {
		event := events[i]
		if event.Type == observability.EventToolCallCompleted && event.Fields["tool_name"] == toolName {
			return event
		}
	}
	t.Fatalf("tool-call event %q not found in %#v", toolName, events)
	return observability.Event{}
}
