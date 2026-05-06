package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"sf-mvp/internal/approval"
	"sf-mvp/internal/brief"
	"sf-mvp/internal/demo"
	"sf-mvp/internal/eval"
	"sf-mvp/internal/ingestion"
	"sf-mvp/internal/notification"
	"sf-mvp/internal/observability"
	"sf-mvp/internal/severity"
)

const (
	reviewPath                   = "/demo/review"
	approvalRequestPath          = "/demo/approvals"
	approvalDecisionPath         = "/demo/approvals/decisions"
	slackNotificationPreviewPath = "/demo/notifications/slack"
	evalReportPath               = "/demo/eval/latest"
	traceReportPathPrefix        = "/demo/traces/"
	budgetReportPath             = "/demo/budget/check"
	evalSummaryRef               = "docs/mvp/quality/eval-plan.md"
	defaultListenAddr            = "127.0.0.1:8080"
)

type Handler struct {
	now          func() time.Time
	approvalMu   sync.Mutex
	approvalGate *approval.Gate
	traceMu      sync.Mutex
	traces       map[string][]observability.Event
}

type Option func(*Handler)

func NewHandler(options ...Option) http.Handler {
	handler := &Handler{
		now:    defaultNow,
		traces: make(map[string][]observability.Event),
	}
	for _, option := range options {
		option(handler)
	}
	if handler.approvalGate == nil {
		handler.approvalGate = approval.NewGate(handler.now)
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

	switch {
	case r.URL.Path == reviewPath:
		h.handleReview(w, r)
	case r.URL.Path == approvalRequestPath:
		h.handleApprovalRequest(w, r)
	case r.URL.Path == approvalDecisionPath:
		h.handleApprovalDecision(w, r)
	case r.URL.Path == slackNotificationPreviewPath:
		h.handleSlackNotificationPreview(w, r)
	case r.URL.Path == evalReportPath:
		h.handleEvalReport(w, r)
	case r.URL.Path == budgetReportPath:
		h.handleBudgetReport(w, r)
	case strings.HasPrefix(r.URL.Path, traceReportPathPrefix):
		h.handleTraceReport(w, r)
	default:
		writeError(w, http.StatusNotFound, "not_found", "supported paths are POST /demo/review, POST /demo/approvals, POST /demo/approvals/decisions, POST /demo/notifications/slack, GET /demo/eval/latest, GET /demo/traces/{trace_id}, and POST /demo/budget/check")
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
	h.storeTraceEvents(review.ObservabilityEvents)
	writeJSON(w, http.StatusOK, apiResponseFromReview(review))
}

func (h *Handler) handleApprovalRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "supported method is POST")
		return
	}

	var request approvalCreateRequest
	if !decodeRequestStruct(w, r, &request) {
		return
	}
	if strings.TrimSpace(request.Channel) == "" {
		h.writeApprovalError(w, errors.New("approval channel is required"))
		return
	}
	if _, err := demo.ComposeIncident(request.IncidentID, demo.Options{Now: h.now}); err != nil {
		h.writeComposeError(w, err)
		return
	}

	h.approvalMu.Lock()
	approvalRequest, err := h.approvalGate.CreateRequest(approval.RequestInput{
		IncidentID: strings.TrimSpace(request.IncidentID),
		Action:     request.Action,
		Scope: approval.Scope{
			IncidentID: strings.TrimSpace(request.IncidentID),
			TargetRef:  notification.SlackTargetRef(request.Channel),
		},
		Reason: request.Reason,
	})
	audit := h.approvalGate.AuditHistory()
	h.approvalMu.Unlock()
	if err != nil {
		h.writeApprovalError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, apiResponseFromApprovalRequest(approvalRequest, audit))
}

func (h *Handler) handleApprovalDecision(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "supported method is POST")
		return
	}

	var request approvalDecisionRequest
	if !decodeRequestStruct(w, r, &request) {
		return
	}

	h.approvalMu.Lock()
	approvalRequest, err := h.approvalGate.Decide(approval.DecisionInput{
		RequestID: request.RequestID,
		Approver:  request.Approver,
		Decision:  request.Decision,
		Reason:    request.Reason,
	})
	audit := h.approvalGate.AuditHistory()
	h.approvalMu.Unlock()
	if err != nil {
		h.writeApprovalError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponseFromApprovalRequest(approvalRequest, audit))
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
	h.approvalMu.Lock()
	preview, err := notification.PreviewSlack(notification.PreviewRequest{
		IncidentID:   review.IncidentID,
		Channel:      request.Channel,
		DeliveryMode: request.DeliveryMode,
		Brief:        briefResultFromReview(review),
		Gate:         h.approvalGate,
		Recorder:     recorder,
		Workflow:     workflow,
	})
	audit := h.approvalGate.AuditHistory()
	h.approvalMu.Unlock()
	if err != nil {
		h.writeNotificationError(w, err)
		return
	}

	events := recorder.Events()
	h.storeTraceEvents(events)
	writeJSON(w, http.StatusOK, apiResponseFromNotificationPreview(workflow.TraceID, preview, events, audit))
}

func (h *Handler) handleEvalReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "supported method is GET")
		return
	}

	report := eval.Run(eval.GoldenCases(), eval.DefaultThresholds())
	recorder := observability.NewRecorder(h.now, observability.Budget{})
	workflow, err := recorder.StartWorkflow("FIC-SYN-EVAL-REPORT", observability.SensitiveData{})
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_eval_report", "eval report could not start a trace")
		return
	}
	recorder.RecordEvalScore(workflow, report, 0)
	events := recorder.Events()
	h.storeTraceEvents(events)

	writeJSON(w, http.StatusOK, apiResponse{
		TraceID:    workflow.TraceID,
		EvalReport: evalReportResponseFromReport(report),
	})
}

func (h *Handler) handleTraceReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "supported method is GET")
		return
	}

	traceID := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, traceReportPathPrefix))
	events, ok := h.traceEvents(traceID)
	if !ok {
		writeError(w, http.StatusNotFound, "trace_not_found", "trace report was not found in this local handler")
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		TraceID:     traceID,
		TraceReport: traceReportResponseFromEvents(traceID, events),
	})
}

func (h *Handler) handleBudgetReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "supported method is POST")
		return
	}

	var request budgetReportRequest
	if !decodeRequestStruct(w, r, &request) {
		return
	}
	if strings.TrimSpace(request.IncidentID) == "" {
		writeError(w, http.StatusBadRequest, "invalid_budget_report", "incident_id is required")
		return
	}

	recorder := observability.NewRecorder(h.now, observability.Budget{
		MaxInputTokens:  request.MaxInputTokens,
		MaxOutputTokens: request.MaxOutputTokens,
		MaxTotalTokens:  request.MaxTotalTokens,
		MaxModelCalls:   request.MaxModelCalls,
	})
	workflow, err := recorder.StartWorkflow(request.IncidentID, observability.SensitiveData{Terms: request.SensitiveTerms})
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_budget_report", "budget report could not start a trace")
		return
	}
	event, err := recorder.RecordModelCall(workflow, observability.ModelCall{
		Provider:     request.Provider,
		Model:        request.Model,
		InputTokens:  request.InputTokens,
		OutputTokens: request.OutputTokens,
		Duration:     0,
	})
	if errors.Is(err, observability.ErrInvalidTokenUsage) {
		h.storeTraceEvents(recorder.Events())
		writeError(w, http.StatusBadRequest, "invalid_token_usage", "token usage cannot be negative")
		return
	}
	if err != nil && !errors.Is(err, observability.ErrBudgetExceeded) {
		h.storeTraceEvents(recorder.Events())
		writeError(w, http.StatusBadRequest, "invalid_budget_report", "budget report could not be recorded")
		return
	}
	events := recorder.Events()
	h.storeTraceEvents(events)

	writeJSON(w, http.StatusOK, apiResponse{
		TraceID:      workflow.TraceID,
		BudgetReport: budgetReportResponseFromEvent(event, err),
	})
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

func (h *Handler) writeApprovalError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, approval.ErrRequestAlreadyDecided):
		writeError(w, http.StatusConflict, "approval_already_decided", "approval request already has a final decision")
	case strings.Contains(err.Error(), "not found"):
		writeError(w, http.StatusNotFound, "approval_request_not_found", "approval request was not found")
	default:
		writeError(w, http.StatusBadRequest, "invalid_approval_request", "approval request could not be processed")
	}
}

type incidentReviewRequest struct {
	IncidentID string `json:"incident_id"`
}

type approvalCreateRequest struct {
	IncidentID string                   `json:"incident_id"`
	Action     severity.SensitiveAction `json:"action"`
	Channel    string                   `json:"channel"`
	Reason     string                   `json:"reason"`
}

type approvalDecisionRequest struct {
	RequestID string            `json:"request_id"`
	Approver  string            `json:"approver"`
	Decision  approval.Decision `json:"decision"`
	Reason    string            `json:"reason"`
}

type slackPreviewRequest struct {
	IncidentID   string                    `json:"incident_id"`
	Channel      string                    `json:"channel"`
	DeliveryMode notification.DeliveryMode `json:"delivery_mode"`
}

type budgetReportRequest struct {
	IncidentID      string   `json:"incident_id"`
	Provider        string   `json:"provider"`
	Model           string   `json:"model"`
	InputTokens     int      `json:"input_tokens"`
	OutputTokens    int      `json:"output_tokens"`
	MaxInputTokens  int      `json:"max_input_tokens"`
	MaxOutputTokens int      `json:"max_output_tokens"`
	MaxTotalTokens  int      `json:"max_total_tokens"`
	MaxModelCalls   int      `json:"max_model_calls"`
	SensitiveTerms  []string `json:"sensitive_terms"`
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

type approvalRequestResponse struct {
	ID             string `json:"id"`
	IncidentID     string `json:"incident_id"`
	Action         string `json:"action"`
	TargetRef      string `json:"target_ref"`
	Decision       string `json:"decision"`
	Reason         string `json:"reason"`
	Approver       string `json:"approver,omitempty"`
	DecisionReason string `json:"decision_reason,omitempty"`
}

type auditEventResponse struct {
	Type       string `json:"type"`
	RequestID  string `json:"request_id,omitempty"`
	IncidentID string `json:"incident_id,omitempty"`
	Action     string `json:"action,omitempty"`
	TargetRef  string `json:"target_ref,omitempty"`
	Actor      string `json:"actor,omitempty"`
	Decision   string `json:"decision,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

type evalSummaryResponse struct {
	Available bool   `json:"available"`
	Ref       string `json:"ref"`
	Command   string `json:"command"`
}

type evalReportResponse struct {
	CaseCount  int                    `json:"case_count"`
	Passed     bool                   `json:"passed"`
	Metrics    evalMetricsResponse    `json:"metrics"`
	Thresholds evalThresholdsResponse `json:"thresholds"`
	Gates      []evalGateResponse     `json:"gates"`
}

type evalMetricsResponse struct {
	SeverityAccuracy       float64 `json:"severity_accuracy"`
	CitationCoverage       float64 `json:"citation_coverage"`
	RecommendationAccuracy float64 `json:"recommendation_accuracy"`
}

type evalThresholdsResponse struct {
	MinSeverityAccuracy              float64 `json:"min_severity_accuracy"`
	MinCitationCoverage              float64 `json:"min_citation_coverage"`
	MinRecommendationAccuracy        float64 `json:"min_recommendation_accuracy"`
	RequireNoUnsupportedClaims       bool    `json:"require_no_unsupported_claims"`
	RequireRedaction                 bool    `json:"require_redaction"`
	RequirePromptInjectionResistance bool    `json:"require_prompt_injection_resistance"`
	RequireApprovalFailSafe          bool    `json:"require_approval_fail_safe"`
}

type evalGateResponse struct {
	Name      string  `json:"name"`
	Actual    float64 `json:"actual,omitempty"`
	Threshold float64 `json:"threshold,omitempty"`
	Required  bool    `json:"required,omitempty"`
	Passed    bool    `json:"passed"`
}

type traceReportResponse struct {
	TraceID    string               `json:"trace_id"`
	IncidentID string               `json:"incident_id"`
	Events     []traceEventResponse `json:"events"`
	Ephemeral  bool                 `json:"ephemeral"`
}

type traceEventResponse struct {
	Type       string             `json:"type"`
	TraceID    string             `json:"trace_id"`
	IncidentID string             `json:"incident_id"`
	OccurredAt string             `json:"occurred_at"`
	DurationMS int64              `json:"duration_ms,omitempty"`
	Fields     map[string]string  `json:"fields,omitempty"`
	Metrics    map[string]float64 `json:"metrics,omitempty"`
	SourceIDs  []string           `json:"source_ids,omitempty"`
	TokenUsage tokenUsageResponse `json:"token_usage,omitempty"`
}

type budgetReportResponse struct {
	Status     string             `json:"status"`
	Reason     string             `json:"reason"`
	TokenUsage tokenUsageResponse `json:"token_usage"`
}

type tokenUsageResponse struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
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
	ApprovalRequest         *approvalRequestResponse     `json:"approval_request,omitempty"`
	AuditHistory            []auditEventResponse         `json:"audit_history,omitempty"`
	EvalSummary             *evalSummaryResponse         `json:"eval_summary,omitempty"`
	EvalReport              *evalReportResponse          `json:"eval_report,omitempty"`
	TraceReport             *traceReportResponse         `json:"trace_report,omitempty"`
	BudgetReport            *budgetReportResponse        `json:"budget_report,omitempty"`
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

func apiResponseFromNotificationPreview(traceID string, preview notification.PreviewResult, events []observability.Event, audit []approval.AuditEvent) apiResponse {
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
		AuditHistory: auditEventResponses(audit),
	}
}

func apiResponseFromApprovalRequest(request approval.Request, audit []approval.AuditEvent) apiResponse {
	return apiResponse{
		ApprovalRequest: approvalRequestResponseFromApproval(request),
		AuditHistory:    auditEventResponses(audit),
	}
}

func evalReportResponseFromReport(report eval.Report) *evalReportResponse {
	return &evalReportResponse{
		CaseCount: report.Summary.CaseCount,
		Passed:    report.Passed,
		Metrics: evalMetricsResponse{
			SeverityAccuracy:       report.Summary.SeverityAccuracy,
			CitationCoverage:       report.Summary.CitationCoverage,
			RecommendationAccuracy: report.Summary.RecommendationAccuracy,
		},
		Thresholds: evalThresholdsResponse{
			MinSeverityAccuracy:              report.Thresholds.MinSeverityAccuracy,
			MinCitationCoverage:              report.Thresholds.MinCitationCoverage,
			MinRecommendationAccuracy:        report.Thresholds.MinRecommendationAccuracy,
			RequireNoUnsupportedClaims:       report.Thresholds.RequireNoUnsupportedClaims,
			RequireRedaction:                 report.Thresholds.RequireRedaction,
			RequirePromptInjectionResistance: report.Thresholds.RequirePromptInjectionResistance,
			RequireApprovalFailSafe:          report.Thresholds.RequireApprovalFailSafe,
		},
		Gates: evalGateResponses(report),
	}
}

func evalGateResponses(report eval.Report) []evalGateResponse {
	return []evalGateResponse{
		{
			Name:      "severity_accuracy",
			Actual:    report.Summary.SeverityAccuracy,
			Threshold: report.Thresholds.MinSeverityAccuracy,
			Passed:    report.Summary.SeverityAccuracy >= report.Thresholds.MinSeverityAccuracy,
		},
		{
			Name:      "citation_coverage",
			Actual:    report.Summary.CitationCoverage,
			Threshold: report.Thresholds.MinCitationCoverage,
			Passed:    report.Summary.CitationCoverage >= report.Thresholds.MinCitationCoverage,
		},
		{
			Name:      "recommendation_accuracy",
			Actual:    report.Summary.RecommendationAccuracy,
			Threshold: report.Thresholds.MinRecommendationAccuracy,
			Passed:    report.Summary.RecommendationAccuracy >= report.Thresholds.MinRecommendationAccuracy,
		},
		{
			Name:     "no_unsupported_claims",
			Required: report.Thresholds.RequireNoUnsupportedClaims,
			Passed:   !report.Thresholds.RequireNoUnsupportedClaims || allCasesPass(report, func(result eval.CaseResult) bool { return len(result.UnsupportedClaims) == 0 }),
		},
		{
			Name:     "redaction",
			Required: report.Thresholds.RequireRedaction,
			Passed:   !report.Thresholds.RequireRedaction || allCasesPass(report, func(result eval.CaseResult) bool { return len(result.RedactionLeaks) == 0 }),
		},
		{
			Name:     "prompt_injection_resistance",
			Required: report.Thresholds.RequirePromptInjectionResistance,
			Passed:   !report.Thresholds.RequirePromptInjectionResistance || allCasesPass(report, func(result eval.CaseResult) bool { return result.PromptInjectionResistant }),
		},
		{
			Name:     "approval_fail_safe",
			Required: report.Thresholds.RequireApprovalFailSafe,
			Passed:   !report.Thresholds.RequireApprovalFailSafe || allCasesPass(report, func(result eval.CaseResult) bool { return result.SensitiveActionsBlockedWithoutApproval }),
		},
	}
}

func allCasesPass(report eval.Report, check func(eval.CaseResult) bool) bool {
	if len(report.Cases) == 0 {
		return false
	}
	for _, result := range report.Cases {
		if !check(result) {
			return false
		}
	}
	return true
}

func traceReportResponseFromEvents(traceID string, events []observability.Event) *traceReportResponse {
	report := &traceReportResponse{
		TraceID:   traceID,
		Events:    traceEventResponses(events),
		Ephemeral: true,
	}
	for _, event := range events {
		if event.TraceID == traceID && strings.TrimSpace(event.IncidentID) != "" {
			report.IncidentID = event.IncidentID
			break
		}
	}
	return report
}

func traceEventResponses(events []observability.Event) []traceEventResponse {
	responses := make([]traceEventResponse, len(events))
	for i, event := range events {
		responses[i] = traceEventResponse{
			Type:       string(event.Type),
			TraceID:    event.TraceID,
			IncidentID: event.IncidentID,
			OccurredAt: event.OccurredAt.UTC().Format(time.RFC3339),
			DurationMS: event.Duration.Milliseconds(),
			Fields:     cloneFields(event.Fields),
			Metrics:    cloneMetrics(event.Metrics),
			SourceIDs:  append([]string(nil), event.SourceIDs...),
			TokenUsage: tokenUsageResponseFromUsage(event.TokenUsage),
		}
	}
	return responses
}

func budgetReportResponseFromEvent(event observability.Event, err error) *budgetReportResponse {
	status := "within_budget"
	reason := "caller-supplied token usage is within configured local budget"
	if errors.Is(err, observability.ErrBudgetExceeded) || event.Type == observability.EventBudgetExceeded {
		status = "budget_exceeded"
		reason = event.Fields["budget_reason"]
	}
	return &budgetReportResponse{
		Status:     status,
		Reason:     reason,
		TokenUsage: tokenUsageResponseFromUsage(event.TokenUsage),
	}
}

func tokenUsageResponseFromUsage(usage observability.TokenUsage) tokenUsageResponse {
	return tokenUsageResponse{
		InputTokens:  usage.InputTokens,
		OutputTokens: usage.OutputTokens,
		TotalTokens:  usage.TotalTokens,
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

func approvalRequestResponseFromApproval(request approval.Request) *approvalRequestResponse {
	return &approvalRequestResponse{
		ID:             request.ID,
		IncidentID:     request.IncidentID,
		Action:         string(request.Action),
		TargetRef:      request.Scope.TargetRef,
		Decision:       string(request.Decision),
		Reason:         request.Reason,
		Approver:       request.Approver,
		DecisionReason: request.DecisionReason,
	}
}

func auditEventResponses(events []approval.AuditEvent) []auditEventResponse {
	responses := make([]auditEventResponse, len(events))
	for i, event := range events {
		responses[i] = auditEventResponse{
			Type:       string(event.Type),
			RequestID:  event.RequestID,
			IncidentID: event.IncidentID,
			Action:     string(event.Action),
			TargetRef:  event.Scope.TargetRef,
			Actor:      event.Actor,
			Decision:   string(event.Decision),
			Reason:     event.Reason,
		}
	}
	return responses
}

func (h *Handler) storeTraceEvents(events []observability.Event) {
	if len(events) == 0 {
		return
	}
	byTraceID := make(map[string][]observability.Event)
	for _, event := range events {
		if strings.TrimSpace(event.TraceID) == "" {
			continue
		}
		byTraceID[event.TraceID] = append(byTraceID[event.TraceID], cloneObservabilityEvent(event))
	}
	if len(byTraceID) == 0 {
		return
	}

	h.traceMu.Lock()
	defer h.traceMu.Unlock()
	for traceID, traceEvents := range byTraceID {
		h.traces[traceID] = append(h.traces[traceID], traceEvents...)
	}
}

func (h *Handler) traceEvents(traceID string) ([]observability.Event, bool) {
	traceID = strings.TrimSpace(traceID)
	if traceID == "" {
		return nil, false
	}

	h.traceMu.Lock()
	defer h.traceMu.Unlock()
	events, ok := h.traces[traceID]
	if !ok || len(events) == 0 {
		return nil, false
	}
	return cloneObservabilityEvents(events), true
}

func cloneObservabilityEvents(events []observability.Event) []observability.Event {
	cloned := make([]observability.Event, len(events))
	for i, event := range events {
		cloned[i] = cloneObservabilityEvent(event)
	}
	return cloned
}

func cloneObservabilityEvent(event observability.Event) observability.Event {
	event.Fields = cloneFields(event.Fields)
	event.Metrics = cloneMetrics(event.Metrics)
	event.SourceIDs = append([]string(nil), event.SourceIDs...)
	return event
}

func cloneFields(fields map[string]string) map[string]string {
	if len(fields) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(fields))
	for key, value := range fields {
		cloned[key] = value
	}
	return cloned
}

func cloneMetrics(metrics map[string]float64) map[string]float64 {
	if len(metrics) == 0 {
		return nil
	}
	cloned := make(map[string]float64, len(metrics))
	for key, value := range metrics {
		cloned[key] = value
	}
	return cloned
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
