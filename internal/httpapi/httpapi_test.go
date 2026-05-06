package httpapi

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReviewIncidentEndpointReturnsBriefSeverityApprovalAndTrace(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/demo/review", strings.NewReader(`{"incident_id":"FIC-SYN-001"}`))

	NewHandler().ServeHTTP(recorder, request)

	response := decodeResponse(t, recorder)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %#v", recorder.Code, http.StatusOK, response)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	if response.TraceID != "trace-fic-syn-001-20260506t160000z-001" {
		t.Fatalf("trace_id = %q, want deterministic trace", response.TraceID)
	}
	if response.Review.IncidentID != "FIC-SYN-001" {
		t.Fatalf("review.incident_id = %q, want FIC-SYN-001", response.Review.IncidentID)
	}
	if response.Review.Severity.Level != "low" {
		t.Fatalf("review.severity.level = %q, want low", response.Review.Severity.Level)
	}
	if response.Review.RedactedBrief.Status != "draft" {
		t.Fatalf("review.redacted_brief.status = %q, want draft", response.Review.RedactedBrief.Status)
	}
	if len(response.ApprovalRequiredActions) != 3 {
		t.Fatalf("approval_required_actions count = %d, want 3", len(response.ApprovalRequiredActions))
	}
	for _, action := range response.ApprovalRequiredActions {
		if action.Status != "blocked" || action.Approved {
			t.Fatalf("approval action = %#v, want blocked and unapproved", action)
		}
	}
	if !response.EvalSummary.Available {
		t.Fatal("eval_summary.available = false, want true")
	}
	if response.EvalSummary.Ref != "docs/mvp/quality/eval-plan.md" {
		t.Fatalf("eval_summary.ref = %q, want docs/mvp/quality/eval-plan.md", response.EvalSummary.Ref)
	}
}

func TestReviewEndpointAcceptsSyntheticPacketJSONInput(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/demo/review", strings.NewReader(syntheticPacketJSON()))

	NewHandler().ServeHTTP(recorder, request)

	response := decodeResponse(t, recorder)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %#v", recorder.Code, http.StatusOK, response)
	}
	if response.Review.IncidentID != "FIC-SYN-901" {
		t.Fatalf("review.incident_id = %q, want FIC-SYN-901", response.Review.IncidentID)
	}
	if response.TraceID != "trace-fic-syn-901-20260506t160000z-001" {
		t.Fatalf("trace_id = %q, want deterministic packet trace", response.TraceID)
	}
	if response.Review.ValidationStatus != "accepted" {
		t.Fatalf("review.validation_status = %q, want accepted", response.Review.ValidationStatus)
	}
}

func TestReviewEndpointRejectsMalformedJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/demo/review", strings.NewReader(`{"incident_id":`))

	NewHandler().ServeHTTP(recorder, request)

	response := decodeResponse(t, recorder)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %#v", recorder.Code, http.StatusBadRequest, response)
	}
	if response.Error.Code != "malformed_json" {
		t.Fatalf("error.code = %q, want malformed_json", response.Error.Code)
	}
	if response.TraceID != "" || response.Review.IncidentID != "" {
		t.Fatalf("response = %#v, want no trace or review", response)
	}
}

func TestReviewEndpointRejectsUnknownIncidentID(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/demo/review", strings.NewReader(`{"incident_id":"FIC-SYN-999"}`))

	NewHandler().ServeHTTP(recorder, request)

	response := decodeResponse(t, recorder)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body = %#v", recorder.Code, http.StatusNotFound, response)
	}
	if response.Error.Code != "incident_not_found" {
		t.Fatalf("error.code = %q, want incident_not_found", response.Error.Code)
	}
}

func TestReviewEndpointRejectsNonSyntheticPacketInput(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/demo/review", strings.NewReader(strings.Replace(syntheticPacketJSON(), `"synthetic_record": true`, `"synthetic_record": false`, 1)))

	NewHandler().ServeHTTP(recorder, request)

	response := decodeResponse(t, recorder)
	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d; body = %#v", recorder.Code, http.StatusUnprocessableEntity, response)
	}
	if response.Error.Code != "non_synthetic_input" {
		t.Fatalf("error.code = %q, want non_synthetic_input", response.Error.Code)
	}
	if response.Review.RedactedBrief.Status != "" {
		t.Fatalf("review.redacted_brief.status = %q, want no brief drafting", response.Review.RedactedBrief.Status)
	}
}

func TestReviewEndpointRejectsUnsupportedMethod(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/demo/review", nil)

	NewHandler().ServeHTTP(recorder, request)

	response := decodeResponse(t, recorder)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d; body = %#v", recorder.Code, http.StatusMethodNotAllowed, response)
	}
	if recorder.Header().Get("Allow") != http.MethodPost {
		t.Fatalf("Allow = %q, want POST", recorder.Header().Get("Allow"))
	}
	if response.Error.Code != "method_not_allowed" {
		t.Fatalf("error.code = %q, want method_not_allowed", response.Error.Code)
	}
}

func TestReviewEndpointRejectsUnknownPath(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/demo/unknown", strings.NewReader(`{"incident_id":"FIC-SYN-001"}`))

	NewHandler().ServeHTTP(recorder, request)

	response := decodeResponse(t, recorder)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body = %#v", recorder.Code, http.StatusNotFound, response)
	}
	if response.Error.Code != "not_found" {
		t.Fatalf("error.code = %q, want not_found", response.Error.Code)
	}
}

func TestSlackNotificationPreviewEndpointReturnsBlockedDryRunPayload(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/demo/notifications/slack", strings.NewReader(`{
		"incident_id": "FIC-SYN-001",
		"channel": "#fleet-safety",
		"delivery_mode": "dry_run"
	}`))

	NewHandler().ServeHTTP(recorder, request)

	response := decodeResponse(t, recorder)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %#v", recorder.Code, http.StatusOK, response)
	}
	if response.NotificationPreview.Status != "blocked" {
		t.Fatalf("notification status = %q, want blocked", response.NotificationPreview.Status)
	}
	if response.NotificationPreview.DeliveryMode != "dry_run" {
		t.Fatalf("delivery_mode = %q, want dry_run", response.NotificationPreview.DeliveryMode)
	}
	if !strings.Contains(response.NotificationPreview.Reason, "approval") {
		t.Fatalf("reason = %q, want approval explanation", response.NotificationPreview.Reason)
	}
	if response.NotificationPreview.PreparedPayload.Channel != "#fleet-safety" {
		t.Fatalf("payload channel = %q, want #fleet-safety", response.NotificationPreview.PreparedPayload.Channel)
	}
	if response.NotificationPreview.PreparedPayload.Text == "" || len(response.NotificationPreview.PreparedPayload.Blocks) == 0 {
		t.Fatalf("prepared payload = %#v, want text and blocks", response.NotificationPreview.PreparedPayload)
	}
	if response.NotificationPreview.Sent || response.NotificationPreview.NetworkDeliveryAttempted {
		t.Fatalf("notification preview = %#v, want no send and no network delivery", response.NotificationPreview)
	}
	if len(response.NotificationPreview.ObservabilityEvents) == 0 {
		t.Fatal("notification preview observability events are empty")
	}
}

func TestSlackNotificationPreviewEndpointRequiresDryRunMode(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/demo/notifications/slack", strings.NewReader(`{
		"incident_id": "FIC-SYN-001",
		"channel": "#fleet-safety",
		"delivery_mode": "send"
	}`))

	NewHandler().ServeHTTP(recorder, request)

	response := decodeResponse(t, recorder)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %#v", recorder.Code, http.StatusBadRequest, response)
	}
	if response.Error.Code != "dry_run_required" {
		t.Fatalf("error.code = %q, want dry_run_required", response.Error.Code)
	}
	if response.NotificationPreview.Status != "" || response.NotificationPreview.PreparedPayload.Text != "" {
		t.Fatalf("notification preview = %#v, want empty preview", response.NotificationPreview)
	}
}

func TestScopedApprovalRetryCreatesRequestAndBlocksDryRunWhilePending(t *testing.T) {
	handler := NewHandler()

	createRecorder, createResponse := postJSON(t, handler, "/demo/approvals", `{
		"incident_id": "FIC-SYN-001",
		"action": "external_sharing",
		"channel": "#fleet-safety",
		"reason": "operator wants to preview the redacted brief"
	}`)
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("create approval status = %d, want %d; body = %#v", createRecorder.Code, http.StatusCreated, createResponse)
	}
	if createResponse.ApprovalRequest.ID == "" {
		t.Fatal("approval_request.id is empty")
	}
	if createResponse.ApprovalRequest.Decision != "pending" {
		t.Fatalf("approval decision = %q, want pending", createResponse.ApprovalRequest.Decision)
	}
	if createResponse.ApprovalRequest.TargetRef != "slack:#fleet-safety" {
		t.Fatalf("target_ref = %q, want slack:#fleet-safety", createResponse.ApprovalRequest.TargetRef)
	}
	assertAuditTypes(t, createResponse.AuditHistory, "approval.requested")

	previewRecorder, previewResponse := postJSON(t, handler, "/demo/notifications/slack", `{
		"incident_id": "FIC-SYN-001",
		"channel": "#fleet-safety",
		"delivery_mode": "dry_run"
	}`)
	if previewRecorder.Code != http.StatusOK {
		t.Fatalf("preview status = %d, want %d; body = %#v", previewRecorder.Code, http.StatusOK, previewResponse)
	}
	if previewResponse.NotificationPreview.Status != "blocked" {
		t.Fatalf("notification status = %q, want blocked", previewResponse.NotificationPreview.Status)
	}
	if previewResponse.NotificationPreview.ApprovalRequestID != createResponse.ApprovalRequest.ID {
		t.Fatalf("approval request id = %q, want %q", previewResponse.NotificationPreview.ApprovalRequestID, createResponse.ApprovalRequest.ID)
	}
	if !strings.Contains(previewResponse.NotificationPreview.Reason, "pending") {
		t.Fatalf("reason = %q, want pending approval explanation", previewResponse.NotificationPreview.Reason)
	}
	if previewResponse.NotificationPreview.Sent || previewResponse.NotificationPreview.NetworkDeliveryAttempted {
		t.Fatalf("notification preview = %#v, want no send and no network delivery", previewResponse.NotificationPreview)
	}
	assertAuditTypes(t, previewResponse.AuditHistory, "approval.requested", "sensitive_action.blocked")
}

func TestScopedApprovalRetryBlocksAfterDeniedDecision(t *testing.T) {
	handler := NewHandler()
	_, createResponse := postJSON(t, handler, "/demo/approvals", `{
		"incident_id": "FIC-SYN-001",
		"action": "external_sharing",
		"channel": "#fleet-safety",
		"reason": "operator wants to preview the redacted brief"
	}`)

	decisionRecorder, decisionResponse := postJSON(t, handler, "/demo/approvals/decisions", `{
		"request_id": "`+createResponse.ApprovalRequest.ID+`",
		"approver": "fleet-safety-lead",
		"decision": "denied",
		"reason": "redacted brief needs another review"
	}`)
	if decisionRecorder.Code != http.StatusOK {
		t.Fatalf("decision status = %d, want %d; body = %#v", decisionRecorder.Code, http.StatusOK, decisionResponse)
	}
	if decisionResponse.ApprovalRequest.Decision != "denied" {
		t.Fatalf("approval decision = %q, want denied", decisionResponse.ApprovalRequest.Decision)
	}
	assertAuditTypes(t, decisionResponse.AuditHistory, "approval.requested", "approval.decided")

	previewRecorder, previewResponse := postJSON(t, handler, "/demo/notifications/slack", `{
		"incident_id": "FIC-SYN-001",
		"channel": "#fleet-safety",
		"delivery_mode": "dry_run"
	}`)
	if previewRecorder.Code != http.StatusOK {
		t.Fatalf("preview status = %d, want %d; body = %#v", previewRecorder.Code, http.StatusOK, previewResponse)
	}
	if previewResponse.NotificationPreview.Status != "blocked" {
		t.Fatalf("notification status = %q, want blocked", previewResponse.NotificationPreview.Status)
	}
	if !strings.Contains(previewResponse.NotificationPreview.Reason, "denied") {
		t.Fatalf("reason = %q, want denied approval explanation", previewResponse.NotificationPreview.Reason)
	}
	if previewResponse.NotificationPreview.ApprovalRequestID != createResponse.ApprovalRequest.ID {
		t.Fatalf("approval request id = %q, want %q", previewResponse.NotificationPreview.ApprovalRequestID, createResponse.ApprovalRequest.ID)
	}
	assertAuditTypes(t, previewResponse.AuditHistory, "approval.requested", "approval.decided", "sensitive_action.blocked")
}

func TestScopedApprovalRetryAllowsExactApprovedDryRunOnly(t *testing.T) {
	handler := NewHandler()
	_, createResponse := postJSON(t, handler, "/demo/approvals", `{
		"incident_id": "FIC-SYN-001",
		"action": "external_sharing",
		"channel": "#fleet-safety",
		"reason": "operator wants to preview the redacted brief"
	}`)
	_, decisionResponse := postJSON(t, handler, "/demo/approvals/decisions", `{
		"request_id": "`+createResponse.ApprovalRequest.ID+`",
		"approver": "fleet-safety-lead",
		"decision": "approved",
		"reason": "redacted brief is approved for this dry-run channel"
	}`)
	if decisionResponse.ApprovalRequest.Decision != "approved" {
		t.Fatalf("approval decision = %q, want approved", decisionResponse.ApprovalRequest.Decision)
	}

	allowedRecorder, allowedResponse := postJSON(t, handler, "/demo/notifications/slack", `{
		"incident_id": "FIC-SYN-001",
		"channel": "#fleet-safety",
		"delivery_mode": "dry_run"
	}`)
	if allowedRecorder.Code != http.StatusOK {
		t.Fatalf("allowed preview status = %d, want %d; body = %#v", allowedRecorder.Code, http.StatusOK, allowedResponse)
	}
	if allowedResponse.NotificationPreview.Status != "allowed" {
		t.Fatalf("notification status = %q, want allowed", allowedResponse.NotificationPreview.Status)
	}
	if allowedResponse.NotificationPreview.ApprovalRequestID != createResponse.ApprovalRequest.ID {
		t.Fatalf("approval request id = %q, want %q", allowedResponse.NotificationPreview.ApprovalRequestID, createResponse.ApprovalRequest.ID)
	}
	if allowedResponse.NotificationPreview.Sent || allowedResponse.NotificationPreview.NetworkDeliveryAttempted {
		t.Fatalf("notification preview = %#v, want allowed dry-run without network delivery", allowedResponse.NotificationPreview)
	}
	assertAuditTypes(t, allowedResponse.AuditHistory, "approval.requested", "approval.decided", "sensitive_action.allowed")

	outOfChannelRecorder, outOfChannelResponse := postJSON(t, handler, "/demo/notifications/slack", `{
		"incident_id": "FIC-SYN-001",
		"channel": "#dispatch-preview",
		"delivery_mode": "dry_run"
	}`)
	if outOfChannelRecorder.Code != http.StatusOK {
		t.Fatalf("out-of-channel status = %d, want %d; body = %#v", outOfChannelRecorder.Code, http.StatusOK, outOfChannelResponse)
	}
	if outOfChannelResponse.NotificationPreview.Status != "blocked" {
		t.Fatalf("out-of-channel status = %q, want blocked", outOfChannelResponse.NotificationPreview.Status)
	}
	if outOfChannelResponse.NotificationPreview.ApprovalRequestID == createResponse.ApprovalRequest.ID {
		t.Fatalf("out-of-channel preview reused request id %q", outOfChannelResponse.NotificationPreview.ApprovalRequestID)
	}

	outOfIncidentRecorder, outOfIncidentResponse := postJSON(t, handler, "/demo/notifications/slack", `{
		"incident_id": "FIC-SYN-002",
		"channel": "#fleet-safety",
		"delivery_mode": "dry_run"
	}`)
	if outOfIncidentRecorder.Code != http.StatusOK {
		t.Fatalf("out-of-incident status = %d, want %d; body = %#v", outOfIncidentRecorder.Code, http.StatusOK, outOfIncidentResponse)
	}
	if outOfIncidentResponse.NotificationPreview.Status != "blocked" {
		t.Fatalf("out-of-incident status = %q, want blocked", outOfIncidentResponse.NotificationPreview.Status)
	}
}

func TestScopedApprovalRetryDoesNotTreatExportApprovalAsNotificationApproval(t *testing.T) {
	handler := NewHandler()
	_, createResponse := postJSON(t, handler, "/demo/approvals", `{
		"incident_id": "FIC-SYN-001",
		"action": "export",
		"channel": "#fleet-safety",
		"reason": "operator approved only export-shaped action"
	}`)
	_, decisionResponse := postJSON(t, handler, "/demo/approvals/decisions", `{
		"request_id": "`+createResponse.ApprovalRequest.ID+`",
		"approver": "fleet-safety-lead",
		"decision": "approved",
		"reason": "export approval is not notification approval"
	}`)
	if decisionResponse.ApprovalRequest.Action != "export" || decisionResponse.ApprovalRequest.Decision != "approved" {
		t.Fatalf("approval request = %#v, want approved export", decisionResponse.ApprovalRequest)
	}

	previewRecorder, previewResponse := postJSON(t, handler, "/demo/notifications/slack", `{
		"incident_id": "FIC-SYN-001",
		"channel": "#fleet-safety",
		"delivery_mode": "dry_run"
	}`)
	if previewRecorder.Code != http.StatusOK {
		t.Fatalf("preview status = %d, want %d; body = %#v", previewRecorder.Code, http.StatusOK, previewResponse)
	}
	if previewResponse.NotificationPreview.Status != "blocked" {
		t.Fatalf("notification status = %q, want blocked", previewResponse.NotificationPreview.Status)
	}
	if previewResponse.NotificationPreview.ApprovalRequestID == createResponse.ApprovalRequest.ID {
		t.Fatalf("notification preview reused out-of-action approval request %q", previewResponse.NotificationPreview.ApprovalRequestID)
	}
}

func TestDefaultListenAddressUsesLoopbackHost(t *testing.T) {
	host, _, err := net.SplitHostPort(DefaultListenAddress())
	if err != nil {
		t.Fatalf("DefaultListenAddress() is not host:port: %v", err)
	}
	if ip := net.ParseIP(host); ip == nil || !ip.IsLoopback() {
		t.Fatalf("DefaultListenAddress() host = %q, want loopback IP", host)
	}
}

type testAPIResponse struct {
	TraceID                 string                       `json:"trace_id"`
	Review                  testReviewResponse           `json:"review"`
	ApprovalRequiredActions []testApprovalActionResponse `json:"approval_required_actions"`
	NotificationPreview     testNotificationPreview      `json:"notification_preview"`
	ApprovalRequest         testApprovalRequestResponse  `json:"approval_request"`
	AuditHistory            []testAuditEventResponse     `json:"audit_history"`
	EvalSummary             testEvalSummaryResponse      `json:"eval_summary"`
	Error                   testErrorResponse            `json:"error"`
}

type testReviewResponse struct {
	ValidationStatus string               `json:"validation_status"`
	IncidentID       string               `json:"incident_id"`
	Severity         testSeverityResponse `json:"severity"`
	RedactedBrief    testRedactedBrief    `json:"redacted_brief"`
}

type testSeverityResponse struct {
	Level string `json:"level"`
}

type testRedactedBrief struct {
	Status string `json:"status"`
}

type testApprovalActionResponse struct {
	Status   string `json:"status"`
	Approved bool   `json:"approved"`
}

type testApprovalRequestResponse struct {
	ID             string `json:"id"`
	IncidentID     string `json:"incident_id"`
	Action         string `json:"action"`
	TargetRef      string `json:"target_ref"`
	Decision       string `json:"decision"`
	Approver       string `json:"approver"`
	DecisionReason string `json:"decision_reason"`
}

type testAuditEventResponse struct {
	Type       string `json:"type"`
	RequestID  string `json:"request_id"`
	IncidentID string `json:"incident_id"`
	Action     string `json:"action"`
	TargetRef  string `json:"target_ref"`
	Actor      string `json:"actor"`
	Decision   string `json:"decision"`
	Reason     string `json:"reason"`
}

type testNotificationPreview struct {
	Status                   string                     `json:"status"`
	DeliveryMode             string                     `json:"delivery_mode"`
	Reason                   string                     `json:"reason"`
	ApprovalRequestID        string                     `json:"approval_request_id"`
	PreparedPayload          testSlackPayload           `json:"prepared_payload"`
	Sent                     bool                       `json:"sent"`
	NetworkDeliveryAttempted bool                       `json:"network_delivery_attempted"`
	ObservabilityEvents      []observabilityEventStatus `json:"observability_events"`
}

type testSlackPayload struct {
	Channel string           `json:"channel"`
	Text    string           `json:"text"`
	Blocks  []testSlackBlock `json:"blocks"`
}

type testSlackBlock struct {
	Type string        `json:"type"`
	Text testSlackText `json:"text"`
}

type testSlackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type testEvalSummaryResponse struct {
	Available bool   `json:"available"`
	Ref       string `json:"ref"`
}

type testErrorResponse struct {
	Code string `json:"code"`
}

func decodeResponse(t *testing.T, recorder *httptest.ResponseRecorder) testAPIResponse {
	t.Helper()

	var response testAPIResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("response JSON did not decode: %v\nbody: %s", err, recorder.Body.String())
	}
	return response
}

func postJSON(t *testing.T, handler http.Handler, path, body string) (*httptest.ResponseRecorder, testAPIResponse) {
	t.Helper()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	handler.ServeHTTP(recorder, request)
	return recorder, decodeResponse(t, recorder)
}

func assertAuditTypes(t *testing.T, events []testAuditEventResponse, want ...string) {
	t.Helper()
	if len(events) != len(want) {
		t.Fatalf("audit history length = %d, want %d: %#v", len(events), len(want), events)
	}
	for i, event := range events {
		if event.Type != want[i] {
			t.Fatalf("audit[%d].type = %q, want %q; history = %#v", i, event.Type, want[i], events)
		}
	}
}

func syntheticPacketJSON() string {
	return `{
		"synthetic_record": true,
		"incident_id": "FIC-SYN-901",
		"vehicle_id": "BUS-901",
		"route": "Synthetic Route 901",
		"timestamp": "2026-03-12T07:42:18-07:00",
		"location_label": "Synthetic Location 901",
		"event_type": "hard_brake",
		"telemetry_samples": [
			{"relative_time":"-02s","speed_mph":16,"heading":"northbound","signal":"steady speed","gps_label":"Synthetic GPS 901"},
			{"relative_time":"+00s","speed_mph":4,"heading":"northbound","signal":"hard brake threshold crossed","gps_label":"Synthetic GPS 901"}
		],
		"media_references": [
			"synthetic://fic-syn-901/front-camera.jpg"
		],
		"transcript_notes": [
			"Driver says this is a synthetic hard-brake API fixture."
		],
		"still_frame_notes": [
			"Front frame shows synthetic road context."
		]
	}`
}
