package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"sf-mvp/internal/approval"
	"sf-mvp/internal/brief"
	"sf-mvp/internal/demo"
	"sf-mvp/internal/ingestion"
	"sf-mvp/internal/notification"
	"sf-mvp/internal/observability"
)

const (
	reviewPath                   = "/demo/review"
	slackNotificationPreviewPath = "/demo/notifications/slack"
	evalSummaryRef               = "docs/mvp/quality/eval-plan.md"
	defaultListenAddr            = "127.0.0.1:8080"
)

type Handler struct {
	now func() time.Time
}

type Option func(*Handler)

func NewHandler(options ...Option) http.Handler {
	handler := &Handler{now: defaultNow}
	for _, option := range options {
		option(handler)
	}
	return handler
}

func WithNow(now func() time.Time) Option {
	return func(handler *Handler) {
		if now != nil {
			handler.now = now
		}
	}
}

func DefaultListenAddress() string {
	return defaultListenAddr
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	case reviewPath:
		h.handleReview(w, r)
	case slackNotificationPreviewPath:
		h.handleSlackNotificationPreview(w, r)
	default:
		writeError(w, http.StatusNotFound, "not_found", "supported paths are POST /demo/review and POST /demo/notifications/slack")
	}
}

func (h *Handler) handleReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "supported method is POST")
		return
	}

	raw, ok := decodeRequestObject(w, r)
	if !ok {
		return
	}
	if len(raw) == 0 {
		writeError(w, http.StatusBadRequest, "empty_request", "request body must include incident_id or a synthetic packet")
		return
	}

	review, err := h.composeReview(raw)
	if err != nil {
		h.writeComposeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponseFromReview(review))
}

func (h *Handler) handleSlackNotificationPreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "supported method is POST")
		return
	}

	var request slackPreviewRequest
	if !decodeRequestStruct(w, r, &request) {
		return
	}
	if request.DeliveryMode != notification.DeliveryModeDryRun {
		h.writeNotificationError(w, notification.ErrDryRunRequired)
		return
	}

	review, err := demo.ComposeIncident(request.IncidentID, demo.Options{Now: h.now})
	if err != nil {
		h.writeComposeError(w, err)
		return
	}

	recorder := observability.NewRecorder(h.now, observability.Budget{})
	workflow, err := recorder.StartWorkflow(review.IncidentID, observability.SensitiveData{})
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_notification_preview", "notification preview could not start a trace")
		return
	}
	preview, err := notification.PreviewSlack(notification.PreviewRequest{
		IncidentID:   review.IncidentID,
		Channel:      request.Channel,
		DeliveryMode: request.DeliveryMode,
		Brief:        briefResultFromReview(review),
		Gate:         approval.NewGate(h.now),
		Recorder:     recorder,
		Workflow:     workflow,
	})
	if err != nil {
		h.writeNotificationError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponseFromNotificationPreview(workflow.TraceID, preview, recorder.Events()))
}

func decodeRequestObject(w http.ResponseWriter, r *http.Request) (map[string]json.RawMessage, bool) {
	body := http.MaxBytesReader(w, r.Body, 1<<20)
	defer body.Close()

	var raw map[string]json.RawMessage
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&raw); err != nil {
		writeError(w, http.StatusBadRequest, "malformed_json", "request body must be valid JSON")
		return nil, false
	}
	return raw, true
}

func decodeRequestStruct(w http.ResponseWriter, r *http.Request, target any) bool {
	body := http.MaxBytesReader(w, r.Body, 1<<20)
	defer body.Close()

	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "malformed_json", "request body must be valid JSON")
		return false
	}
	return true
}

func (h *Handler) composeReview(raw map[string]json.RawMessage) (demo.ReviewResult, error) {
	options := demo.Options{Now: h.now}
	if _, ok := raw["synthetic_record"]; ok {
		packetJSON, err := json.Marshal(raw)
		if err != nil {
			return demo.ReviewResult{}, err
		}
		result, err := ingestion.IngestJSON(packetJSON)
		if err != nil {
			if validationHasNonSyntheticIssue(err) {
				return demo.ReviewResult{}, demo.ErrNonSyntheticInput
			}
			return demo.ReviewResult{}, err
		}
		return demo.ComposePacket(result.Packet, options)
	}

	var request incidentReviewRequest
	requestJSON, err := json.Marshal(raw)
	if err != nil {
		return demo.ReviewResult{}, err
	}
	decoder := json.NewDecoder(bytes.NewReader(requestJSON))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		return demo.ReviewResult{}, err
	}
	if strings.TrimSpace(request.IncidentID) == "" {
		return demo.ReviewResult{}, demo.ErrIncidentNotFound
	}
	return demo.ComposeIncident(request.IncidentID, options)
}

func (h *Handler) writeComposeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, demo.ErrIncidentNotFound):
		writeError(w, http.StatusNotFound, "incident_not_found", "synthetic incident ID was not found")
	case errors.Is(err, demo.ErrNonSyntheticInput):
		writeError(w, http.StatusUnprocessableEntity, "non_synthetic_input", "only synthetic FIC-SYN incident input is accepted")
	case errors.Is(err, demo.ErrMissingEvidence):
		writeError(w, http.StatusUnprocessableEntity, "missing_evidence", "required review evidence is missing")
	default:
		writeError(w, http.StatusBadRequest, "invalid_request", "request could not be composed")
	}
}

func (h *Handler) writeNotificationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, notification.ErrDryRunRequired):
		writeError(w, http.StatusBadRequest, "dry_run_required", "delivery_mode must be dry_run")
	case errors.Is(err, notification.ErrInvalidPreviewInput), errors.Is(err, notification.ErrNoRedactedBriefInput):
		writeError(w, http.StatusBadRequest, "invalid_notification_preview", "notification preview request is invalid")
	default:
		writeError(w, http.StatusBadRequest, "invalid_notification_preview", "notification preview could not be prepared")
	}
}

type incidentReviewRequest struct {
	IncidentID string `json:"incident_id"`
}

type slackPreviewRequest struct {
	IncidentID   string                    `json:"incident_id"`
	Channel      string                    `json:"channel"`
	DeliveryMode notification.DeliveryMode `json:"delivery_mode"`
}

type reviewResponse struct {
	ValidationStatus      demo.ValidationStatus      `json:"validation_status"`
	IncidentID            string                     `json:"incident_id"`
	RetrievedCitationRefs []string                   `json:"retrieved_citation_refs"`
	TimelineEntries       []timelineEntryResponse    `json:"timeline_entries"`
	Severity              severityResponse           `json:"severity"`
	Recommendations       []recommendationResponse   `json:"recommendations"`
	RedactedBrief         briefResponse              `json:"redacted_brief"`
	ObservabilityEvents   []observabilityEventStatus `json:"observability_events"`
}

type timelineEntryResponse struct {
	Time        string   `json:"time"`
	Claim       string   `json:"claim"`
	SourceRefs  []string `json:"source_refs"`
	Uncertain   bool     `json:"uncertain"`
	Uncertainty string   `json:"uncertainty,omitempty"`
}

type severityResponse struct {
	Level     string                `json:"level"`
	Rationale []explanationResponse `json:"rationale"`
}

type explanationResponse struct {
	Text       string   `json:"text"`
	SourceRefs []string `json:"source_refs"`
}

type recommendationResponse struct {
	Action      string   `json:"action"`
	Explanation string   `json:"explanation"`
	SourceRefs  []string `json:"source_refs"`
}

type briefResponse struct {
	Status            string                  `json:"status"`
	Shareable         bool                    `json:"shareable"`
	Sections          []briefSectionResponse  `json:"sections"`
	ApprovalState     []approvalStateResponse `json:"approval_state"`
	RedactionsApplied []redactionResponse     `json:"redactions_applied"`
	Uncertainties     []string                `json:"uncertainties"`
}

type briefSectionResponse struct {
	Title      string   `json:"title"`
	Text       string   `json:"text"`
	SourceRefs []string `json:"source_refs"`
}

type approvalStateResponse struct {
	Action  string `json:"action"`
	Blocked bool   `json:"blocked"`
	Reason  string `json:"reason"`
}

type redactionResponse struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type approvalActionResponse struct {
	Action    string `json:"action"`
	Required  bool   `json:"required"`
	Approved  bool   `json:"approved"`
	Status    string `json:"status"`
	TargetRef string `json:"target_ref"`
	RequestID string `json:"request_id,omitempty"`
	Reason    string `json:"reason"`
}

type observabilityEventStatus struct {
	Type    string `json:"type"`
	TraceID string `json:"trace_id"`
}

type evalSummaryResponse struct {
	Available bool   `json:"available"`
	Ref       string `json:"ref"`
	Command   string `json:"command"`
}

type notificationPreviewResponse struct {
	Status                   string                     `json:"status"`
	DeliveryMode             string                     `json:"delivery_mode"`
	Reason                   string                     `json:"reason"`
	ApprovalRequestID        string                     `json:"approval_request_id,omitempty"`
	PreparedPayload          notification.SlackPayload  `json:"prepared_payload"`
	Sent                     bool                       `json:"sent"`
	NetworkDeliveryAttempted bool                       `json:"network_delivery_attempted"`
	ObservabilityEvents      []observabilityEventStatus `json:"observability_events"`
}

type apiResponse struct {
	TraceID                 string                       `json:"trace_id,omitempty"`
	Review                  *reviewResponse              `json:"review,omitempty"`
	ApprovalRequiredActions []approvalActionResponse     `json:"approval_required_actions,omitempty"`
	NotificationPreview     *notificationPreviewResponse `json:"notification_preview,omitempty"`
	EvalSummary             *evalSummaryResponse         `json:"eval_summary,omitempty"`
	Error                   *errorResponse               `json:"error,omitempty"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func apiResponseFromReview(review demo.ReviewResult) apiResponse {
	return apiResponse{
		TraceID:                 review.TraceID,
		Review:                  reviewResponseFromReview(review),
		ApprovalRequiredActions: approvalActionResponses(review.ApprovalRequiredActions),
		EvalSummary: &evalSummaryResponse{
			Available: true,
			Ref:       evalSummaryRef,
			Command:   "go test ./internal/eval",
		},
	}
}

func apiResponseFromNotificationPreview(traceID string, preview notification.PreviewResult, events []observability.Event) apiResponse {
	return apiResponse{
		TraceID: traceID,
		NotificationPreview: &notificationPreviewResponse{
			Status:                   string(preview.Status),
			DeliveryMode:             string(preview.DeliveryMode),
			Reason:                   preview.Reason,
			ApprovalRequestID:        preview.ApprovalRequestID,
			PreparedPayload:          preview.PreparedPayload,
			Sent:                     preview.Sent,
			NetworkDeliveryAttempted: preview.NetworkDeliveryAttempted,
			ObservabilityEvents:      observabilityEventStatuses(events),
		},
	}
}

func reviewResponseFromReview(review demo.ReviewResult) *reviewResponse {
	return &reviewResponse{
		ValidationStatus:      review.ValidationStatus,
		IncidentID:            review.IncidentID,
		RetrievedCitationRefs: append([]string(nil), review.RetrievedCitationRefs...),
		TimelineEntries:       timelineEntryResponses(review.TimelineEntries),
		Severity: severityResponse{
			Level:     string(review.Severity.Level),
			Rationale: explanationResponses(review.Severity.Rationale),
		},
		Recommendations:     recommendationResponses(review.Recommendations),
		RedactedBrief:       briefResponseFromReview(review.RedactedBrief),
		ObservabilityEvents: observabilityEventStatuses(review.ObservabilityEvents),
	}
}

func timelineEntryResponses(entries []demo.ReviewTimelineEntry) []timelineEntryResponse {
	responses := make([]timelineEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = timelineEntryResponse{
			Time:        entry.Time.UTC().Format(time.RFC3339),
			Claim:       entry.Claim,
			SourceRefs:  append([]string(nil), entry.SourceRefs...),
			Uncertain:   entry.Uncertain,
			Uncertainty: entry.Uncertainty,
		}
	}
	return responses
}

func explanationResponses(explanations []demo.ReviewExplanation) []explanationResponse {
	responses := make([]explanationResponse, len(explanations))
	for i, explanation := range explanations {
		responses[i] = explanationResponse{
			Text:       explanation.Text,
			SourceRefs: append([]string(nil), explanation.SourceRefs...),
		}
	}
	return responses
}

func recommendationResponses(recommendations []demo.ReviewRecommendation) []recommendationResponse {
	responses := make([]recommendationResponse, len(recommendations))
	for i, recommendation := range recommendations {
		responses[i] = recommendationResponse{
			Action:      string(recommendation.Action),
			Explanation: recommendation.Explanation,
			SourceRefs:  append([]string(nil), recommendation.SourceRefs...),
		}
	}
	return responses
}

func briefResponseFromReview(brief demo.ReviewBrief) briefResponse {
	sections := make([]briefSectionResponse, len(brief.Sections))
	for i, section := range brief.Sections {
		sections[i] = briefSectionResponse{
			Title:      section.Title,
			Text:       section.Text,
			SourceRefs: append([]string(nil), section.SourceRefs...),
		}
	}
	approvalState := make([]approvalStateResponse, len(brief.ApprovalState))
	for i, state := range brief.ApprovalState {
		approvalState[i] = approvalStateResponse{
			Action:  string(state.Action),
			Blocked: state.Blocked,
			Reason:  state.Reason,
		}
	}
	redactions := make([]redactionResponse, len(brief.RedactionsApplied))
	for i, redaction := range brief.RedactionsApplied {
		redactions[i] = redactionResponse{
			Field:  redaction.Field,
			Reason: redaction.Reason,
		}
	}
	return briefResponse{
		Status:            string(brief.Status),
		Shareable:         brief.Shareable,
		Sections:          sections,
		ApprovalState:     approvalState,
		RedactionsApplied: redactions,
		Uncertainties:     append([]string(nil), brief.Uncertainties...),
	}
}

func briefResultFromReview(review demo.ReviewResult) brief.Result {
	sections := make([]brief.Section, len(review.RedactedBrief.Sections))
	for i, section := range review.RedactedBrief.Sections {
		sections[i] = brief.Section{
			Title:   section.Title,
			Text:    section.Text,
			Sources: briefSourcesFromRefs(section.SourceRefs),
		}
	}
	return brief.Result{
		Status:            review.RedactedBrief.Status,
		IncidentID:        review.IncidentID,
		SyntheticRecord:   true,
		Sections:          sections,
		ApprovalState:     append([]brief.ApprovalState(nil), review.RedactedBrief.ApprovalState...),
		RedactionsApplied: append([]brief.Redaction(nil), review.RedactedBrief.RedactionsApplied...),
		Uncertainties:     append([]string(nil), review.RedactedBrief.Uncertainties...),
		Shareable:         review.RedactedBrief.Shareable,
	}
}

func briefSourcesFromRefs(refs []string) []brief.Source {
	sources := make([]brief.Source, 0, len(refs))
	for _, ref := range refs {
		if strings.TrimSpace(ref) == "" {
			continue
		}
		sources = append(sources, brief.Source{Ref: ref, Type: brief.SourceTypePacket})
	}
	return sources
}

func approvalActionResponses(actions []demo.ApprovalRequiredAction) []approvalActionResponse {
	responses := make([]approvalActionResponse, len(actions))
	for i, action := range actions {
		responses[i] = approvalActionResponse{
			Action:    string(action.Action),
			Required:  action.Required,
			Approved:  action.Approved,
			Status:    string(action.Status),
			TargetRef: action.TargetRef,
			RequestID: action.RequestID,
			Reason:    action.Reason,
		}
	}
	return responses
}

func observabilityEventStatuses(events []observability.Event) []observabilityEventStatus {
	responses := make([]observabilityEventStatus, len(events))
	for i, event := range events {
		responses[i] = observabilityEventStatus{
			Type:    string(event.Type),
			TraceID: event.TraceID,
		}
	}
	return responses
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, apiResponse{
		Error: &errorResponse{
			Code:    code,
			Message: message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, response apiResponse) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func validationHasNonSyntheticIssue(err error) bool {
	var validationErr ingestion.ValidationError
	if !errors.As(err, &validationErr) {
		return false
	}
	for _, issue := range validationErr.Issues {
		switch issue.Code {
		case "synthetic_record.required", "incident_id.synthetic_prefix.required":
			return true
		}
	}
	return false
}

func defaultNow() time.Time {
	return time.Date(2026, time.May, 6, 16, 0, 0, 0, time.UTC)
}
