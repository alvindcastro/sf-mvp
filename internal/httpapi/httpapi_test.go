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
