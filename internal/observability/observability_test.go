package observability

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"sf-mvp/internal/approval"
	"sf-mvp/internal/eval"
	"sf-mvp/internal/retrieval"
	"sf-mvp/internal/severity"
)

func TestStartWorkflowGeneratesTraceIDAndStructuredEvent(t *testing.T) {
	recorder := NewRecorder(fixedClock(), Budget{})

	workflow, err := recorder.StartWorkflow("FIC-SYN-009", SensitiveData{})
	if err != nil {
		t.Fatalf("StartWorkflow returned error: %v", err)
	}

	if workflow.TraceID != "trace-fic-syn-009-20260506t160000z-001" {
		t.Fatalf("TraceID = %q, want deterministic trace ID", workflow.TraceID)
	}
	if workflow.IncidentID != "FIC-SYN-009" {
		t.Fatalf("IncidentID = %q, want FIC-SYN-009", workflow.IncidentID)
	}

	events := recorder.Events()
	if len(events) != 1 {
		t.Fatalf("event count = %d, want 1: %#v", len(events), events)
	}
	event := events[0]
	if event.Type != EventWorkflowStarted {
		t.Fatalf("event Type = %q, want %q", event.Type, EventWorkflowStarted)
	}
	if event.TraceID != workflow.TraceID || event.IncidentID != workflow.IncidentID {
		t.Fatalf("event trace fields = %#v, want workflow identity", event)
	}
	if !event.OccurredAt.Equal(testTime()) {
		t.Fatalf("OccurredAt = %s, want %s", event.OccurredAt, testTime())
	}
}

func TestRecordRetrievalTracksCountSourcesAndLatency(t *testing.T) {
	recorder := NewRecorder(fixedClock(), Budget{})
	workflow := mustStartWorkflow(t, recorder, "FIC-SYN-010", SensitiveData{})

	event := recorder.RecordRetrieval(workflow, retrieval.Result{Matches: []retrieval.Citation{
		{SourceID: "FIC-SOP-HARD-BRAKE-001", CitationRef: "FIC-SOP-HARD-BRAKE-001#2026-02-15", Snippet: "Private evidence snippet must not be logged", ContentRole: retrieval.ContentRoleData},
		{SourceID: "FIC-TS-MISSING-MEDIA-001", CitationRef: "FIC-TS-MISSING-MEDIA-001#2026-02-17", Snippet: "Another private snippet must not be logged", ContentRole: retrieval.ContentRoleData},
	}}, 42*time.Millisecond)

	if event.Type != EventRetrievalCompleted {
		t.Fatalf("event Type = %q, want retrieval completed", event.Type)
	}
	if event.Metrics["retrieval_count"] != 2 {
		t.Fatalf("retrieval_count = %.0f, want 2", event.Metrics["retrieval_count"])
	}
	wantSources := []string{"FIC-SOP-HARD-BRAKE-001", "FIC-TS-MISSING-MEDIA-001"}
	if !reflect.DeepEqual(event.SourceIDs, wantSources) {
		t.Fatalf("SourceIDs = %#v, want %#v", event.SourceIDs, wantSources)
	}
	if event.Duration != 42*time.Millisecond {
		t.Fatalf("Duration = %s, want 42ms", event.Duration)
	}
	if len(event.Fields) != 0 {
		t.Fatalf("Fields = %#v, want no retrieval snippets or citation text logged", event.Fields)
	}
}

func TestRecordToolCallTracksSuccessAndRedactsSensitiveFields(t *testing.T) {
	recorder := NewRecorder(fixedClock(), Budget{})
	workflow := mustStartWorkflow(t, recorder, "FIC-SYN-011", SensitiveData{
		Terms: []string{
			"BUS-SECRET-214",
			"Private North Loop Route",
			"49.2827,-123.1207 private yard",
		},
	})

	event := recorder.RecordToolCall(workflow, ToolCall{
		Name:     "draft_brief",
		Success:  true,
		Duration: 18 * time.Millisecond,
		Fields: map[string]string{
			"summary": "Brief for BUS-SECRET-214 on Private North Loop Route near 49.2827,-123.1207 private yard",
		},
	})

	if event.Type != EventToolCallCompleted {
		t.Fatalf("event Type = %q, want tool_call.completed", event.Type)
	}
	if event.Fields["tool_name"] != "draft_brief" || event.Fields["success"] != "true" {
		t.Fatalf("tool fields = %#v, want name and success", event.Fields)
	}
	summary := event.Fields["summary"]
	for _, leaked := range []string{"BUS-SECRET-214", "Private North Loop Route", "49.2827,-123.1207", "private yard"} {
		if strings.Contains(summary, leaked) {
			t.Fatalf("summary leaked %q: %q", leaked, summary)
		}
	}
	if !strings.Contains(summary, RedactedValue) {
		t.Fatalf("summary = %q, want redaction marker", summary)
	}
}

func TestRecordToolCallRedactsOverlappingSensitiveTermsLongestFirst(t *testing.T) {
	recorder := NewRecorder(fixedClock(), Budget{})
	workflow := mustStartWorkflow(t, recorder, "FIC-SYN-011A", SensitiveData{
		Terms: []string{
			"Oak Street",
			"Oak Street at Pine Avenue",
		},
	})

	event := recorder.RecordToolCall(workflow, ToolCall{
		Name:    "draft_brief",
		Success: true,
		Fields: map[string]string{
			"summary": "Telemetry near Oak Street at Pine Avenue",
		},
	})

	summary := event.Fields["summary"]
	if strings.Contains(summary, "Pine Avenue") {
		t.Fatalf("summary leaked overlapping sensitive suffix: %q", summary)
	}
	if !strings.Contains(summary, RedactedValue) {
		t.Fatalf("summary = %q, want redaction marker", summary)
	}
}

func TestRecordApprovalDecisionTracksDecisionActionAndScope(t *testing.T) {
	recorder := NewRecorder(fixedClock(), Budget{})
	workflow := mustStartWorkflow(t, recorder, "FIC-SYN-012", SensitiveData{})

	request := approval.Request{
		ID:         "approval-001",
		IncidentID: "FIC-SYN-012",
		Action:     severity.SensitiveActionExport,
		Scope:      approval.Scope{IncidentID: "FIC-SYN-012", TargetRef: "brief:FIC-SYN-012"},
		Decision:   approval.DecisionDenied,
		Approver:   "fleet-safety-lead",
	}
	event := recorder.RecordApprovalDecision(workflow, request, 7*time.Millisecond)

	if event.Type != EventApprovalDecisionRecorded {
		t.Fatalf("event Type = %q, want approval decision", event.Type)
	}
	if event.Fields["approval_decision"] != string(approval.DecisionDenied) {
		t.Fatalf("approval_decision = %q, want denied", event.Fields["approval_decision"])
	}
	if event.Fields["action"] != string(severity.SensitiveActionExport) {
		t.Fatalf("action = %q, want export", event.Fields["action"])
	}
	if event.Fields["target_ref"] != "brief:FIC-SYN-012" {
		t.Fatalf("target_ref = %q, want brief target", event.Fields["target_ref"])
	}
	if event.Fields["request_id"] != "approval-001" || event.Fields["approver"] != "fleet-safety-lead" {
		t.Fatalf("approval metadata = %#v, want request ID and approver", event.Fields)
	}
	if event.Duration != 7*time.Millisecond {
		t.Fatalf("Duration = %s, want 7ms", event.Duration)
	}
}

func TestRecordModelCallRecordsTokensLatencyAndEnforcesBudget(t *testing.T) {
	recorder := NewRecorder(fixedClock(), Budget{
		MaxInputTokens:  200,
		MaxOutputTokens: 100,
		MaxTotalTokens:  300,
		MaxModelCalls:   1,
	})
	workflow := mustStartWorkflow(t, recorder, "FIC-SYN-013", SensitiveData{})

	event, err := recorder.RecordModelCall(workflow, ModelCall{
		Provider:     "hosted",
		Model:        "smaller-review-model",
		InputTokens:  120,
		OutputTokens: 50,
		Duration:     95 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("RecordModelCall returned error under budget: %v", err)
	}
	if event.Type != EventModelCallRecorded {
		t.Fatalf("event Type = %q, want model call recorded", event.Type)
	}
	if event.TokenUsage.InputTokens != 120 || event.TokenUsage.OutputTokens != 50 || event.TokenUsage.TotalTokens != 170 {
		t.Fatalf("TokenUsage = %#v, want 120 input, 50 output, 170 total", event.TokenUsage)
	}
	if event.Duration != 95*time.Millisecond {
		t.Fatalf("Duration = %s, want 95ms", event.Duration)
	}

	budgetEvent, err := recorder.RecordModelCall(workflow, ModelCall{
		Provider:     "hosted",
		Model:        "larger-review-model",
		InputTokens:  1,
		OutputTokens: 1,
		Duration:     10 * time.Millisecond,
	})
	if !errors.Is(err, ErrBudgetExceeded) {
		t.Fatalf("RecordModelCall error = %v, want ErrBudgetExceeded", err)
	}
	if budgetEvent.Type != EventBudgetExceeded {
		t.Fatalf("budget event Type = %q, want budget.exceeded", budgetEvent.Type)
	}
	if budgetEvent.Fields["budget_reason"] != "model call limit exceeded" {
		t.Fatalf("budget_reason = %q, want model call limit exceeded", budgetEvent.Fields["budget_reason"])
	}
}

func TestRecordModelCallRejectsNegativeTokenUsageWithoutConsumingBudget(t *testing.T) {
	recorder := NewRecorder(fixedClock(), Budget{MaxModelCalls: 1})
	workflow := mustStartWorkflow(t, recorder, "FIC-SYN-015", SensitiveData{})

	rejectedEvent, err := recorder.RecordModelCall(workflow, ModelCall{
		Provider:     "hosted",
		Model:        "bad-usage-model",
		InputTokens:  -1,
		OutputTokens: 10,
		Duration:     4 * time.Millisecond,
	})
	if !errors.Is(err, ErrInvalidTokenUsage) {
		t.Fatalf("RecordModelCall error = %v, want ErrInvalidTokenUsage", err)
	}
	if rejectedEvent.Type != EventModelCallRejected {
		t.Fatalf("rejected event Type = %q, want model_call.rejected", rejectedEvent.Type)
	}
	if rejectedEvent.Fields["reject_reason"] != "token usage cannot be negative" {
		t.Fatalf("reject_reason = %q, want negative token usage reason", rejectedEvent.Fields["reject_reason"])
	}

	_, err = recorder.RecordModelCall(workflow, ModelCall{
		Provider:     "hosted",
		Model:        "smaller-review-model",
		InputTokens:  25,
		OutputTokens: 15,
	})
	if err != nil {
		t.Fatalf("valid RecordModelCall after rejection returned error: %v", err)
	}

	_, err = recorder.RecordModelCall(workflow, ModelCall{
		Provider:     "hosted",
		Model:        "smaller-review-model",
		InputTokens:  1,
		OutputTokens: 1,
	})
	if !errors.Is(err, ErrBudgetExceeded) {
		t.Fatalf("second valid RecordModelCall error = %v, want ErrBudgetExceeded", err)
	}
}

func TestRecordModelCallEnforcesInputOutputAndTotalTokenBudgets(t *testing.T) {
	cases := []struct {
		name   string
		budget Budget
		call   ModelCall
		reason string
	}{
		{
			name:   "input",
			budget: Budget{MaxInputTokens: 50},
			call:   ModelCall{Provider: "hosted", Model: "review-model", InputTokens: 51},
			reason: "input token budget exceeded",
		},
		{
			name:   "output",
			budget: Budget{MaxOutputTokens: 20},
			call:   ModelCall{Provider: "hosted", Model: "review-model", OutputTokens: 21},
			reason: "output token budget exceeded",
		},
		{
			name:   "total",
			budget: Budget{MaxTotalTokens: 70},
			call:   ModelCall{Provider: "hosted", Model: "review-model", InputTokens: 50, OutputTokens: 21},
			reason: "total token budget exceeded",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := NewRecorder(fixedClock(), tc.budget)
			workflow := mustStartWorkflow(t, recorder, "FIC-SYN-016-"+strings.ToUpper(tc.name), SensitiveData{})

			event, err := recorder.RecordModelCall(workflow, tc.call)
			if !errors.Is(err, ErrBudgetExceeded) {
				t.Fatalf("RecordModelCall error = %v, want ErrBudgetExceeded", err)
			}
			if event.Type != EventBudgetExceeded {
				t.Fatalf("event Type = %q, want budget.exceeded", event.Type)
			}
			if event.Fields["budget_reason"] != tc.reason {
				t.Fatalf("budget_reason = %q, want %q", event.Fields["budget_reason"], tc.reason)
			}
		})
	}
}

func TestRecordEvalScoreTracksEvalMetrics(t *testing.T) {
	recorder := NewRecorder(fixedClock(), Budget{})
	workflow := mustStartWorkflow(t, recorder, "FIC-SYN-014", SensitiveData{})

	event := recorder.RecordEvalScore(workflow, eval.Report{
		Summary: eval.Summary{
			CaseCount:              5,
			SeverityAccuracy:       1,
			CitationCoverage:       0.98,
			RecommendationAccuracy: 1,
		},
		Passed: true,
	}, 215*time.Millisecond)

	if event.Type != EventEvalScoreRecorded {
		t.Fatalf("event Type = %q, want eval score", event.Type)
	}
	if event.Fields["passed"] != "true" {
		t.Fatalf("passed field = %q, want true", event.Fields["passed"])
	}
	if event.Metrics["case_count"] != 5 || event.Metrics["citation_coverage"] != 0.98 {
		t.Fatalf("metrics = %#v, want eval summary metrics", event.Metrics)
	}
	if event.Duration != 215*time.Millisecond {
		t.Fatalf("Duration = %s, want 215ms", event.Duration)
	}
}

func TestWorkflowAttributesMapSafeOpenTelemetryStyleFields(t *testing.T) {
	workflow := Workflow{
		TraceID:    "trace-fic-syn-017-20260506t160000z-001",
		IncidentID: "FIC-SYN-017",
	}

	attributes := WorkflowAttributes(WorkflowAttributeInput{
		Workflow:           workflow,
		RetrievedSourceIDs: []string{" FIC-TS-MISSING-MEDIA-001 ", "", "FIC-SOP-HARD-BRAKE-001", "FIC-SOP-HARD-BRAKE-001"},
		SeverityLabel:      severity.LevelMedium,
		ApprovalState:      approval.DecisionPending,
		Latency:            42*time.Millisecond + 500*time.Microsecond,
		RawTranscriptNotes: []string{"Driver said private transcript detail near Depot 9."},
		RawMediaReferences: []string{"synthetic://media/private-camera-angle-017.mp4"},
		RawStillFrameNotes: []string{"Still frame shows private side-door detail."},
	})

	want := map[string]string{
		"workflow.trace_id":             "trace-fic-syn-017-20260506t160000z-001",
		"workflow.incident_id_hash":     "sha256:cdcd7e14394344372016403370275f9d314de594321ed9cca7b7e977f058ecc1",
		"workflow.retrieved_source_ids": "FIC-SOP-HARD-BRAKE-001,FIC-TS-MISSING-MEDIA-001",
		"workflow.severity":             "medium",
		"workflow.approval_state":       "pending",
		"workflow.latency_ms":           "42.5",
	}
	if !reflect.DeepEqual(attributes, want) {
		t.Fatalf("WorkflowAttributes() = %#v, want %#v", attributes, want)
	}
	assertAttributeValuesDoNotContain(t, attributes,
		"FIC-SYN-017",
		"private transcript detail",
		"Depot 9",
		"synthetic://media/private-camera-angle-017.mp4",
		"private-camera-angle-017",
		"private side-door detail",
	)
	assertAttributeKeysDoNotContain(t, attributes, "transcript", "media", "still")
}

func TestDefaultCostPlanDefinesCacheCandidatesAndModelRoutingNotes(t *testing.T) {
	plan := DefaultCostPlan()

	if len(plan.CacheCandidates) < 3 {
		t.Fatalf("CacheCandidates length = %d, want at least 3", len(plan.CacheCandidates))
	}
	assertCacheCandidate(t, plan.CacheCandidates, "retrieval_results")
	assertCacheCandidate(t, plan.CacheCandidates, "redacted_brief_drafts")
	assertCacheCandidate(t, plan.CacheCandidates, "eval_reports")

	if len(plan.ModelRouting) < 3 {
		t.Fatalf("ModelRouting length = %d, want hosted, smaller, and self-hosted notes", len(plan.ModelRouting))
	}
	assertModelRoute(t, plan.ModelRouting, ModelRouteHosted)
	assertModelRoute(t, plan.ModelRouting, ModelRouteSmaller)
	assertModelRoute(t, plan.ModelRouting, ModelRouteSelfHosted)
}

func mustStartWorkflow(t *testing.T, recorder *Recorder, incidentID string, sensitive SensitiveData) Workflow {
	t.Helper()

	workflow, err := recorder.StartWorkflow(incidentID, sensitive)
	if err != nil {
		t.Fatalf("StartWorkflow returned error: %v", err)
	}
	return workflow
}

func assertCacheCandidate(t *testing.T, candidates []CacheCandidate, name string) {
	t.Helper()

	for _, candidate := range candidates {
		if candidate.Name == name {
			if candidate.Key == "" || candidate.Reason == "" {
				t.Fatalf("candidate %q missing key or reason: %#v", name, candidate)
			}
			return
		}
	}
	t.Fatalf("cache candidate %q not found in %#v", name, candidates)
}

func assertModelRoute(t *testing.T, routes []ModelRoutingNote, route ModelRoute) {
	t.Helper()

	for _, note := range routes {
		if note.Route == route {
			if note.UseWhen == "" || note.Control == "" {
				t.Fatalf("route %q missing use/control: %#v", route, note)
			}
			return
		}
	}
	t.Fatalf("model route %q not found in %#v", route, routes)
}

func assertAttributeValuesDoNotContain(t *testing.T, attributes map[string]string, leaked ...string) {
	t.Helper()

	for key, value := range attributes {
		for _, term := range leaked {
			if strings.Contains(value, term) {
				t.Fatalf("attribute %q leaked %q in value %q", key, term, value)
			}
		}
	}
}

func assertAttributeKeysDoNotContain(t *testing.T, attributes map[string]string, leaked ...string) {
	t.Helper()

	for key := range attributes {
		for _, term := range leaked {
			if strings.Contains(key, term) {
				t.Fatalf("attribute key %q leaked raw evidence category %q", key, term)
			}
		}
	}
}

func fixedClock() func() time.Time {
	return func() time.Time { return testTime() }
}

func testTime() time.Time {
	return time.Date(2026, time.May, 6, 16, 0, 0, 0, time.UTC)
}
