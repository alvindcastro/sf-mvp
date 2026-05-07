package eval

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sf-mvp/internal/severity"
)

func TestScoreEventsFromPromptfooOutputMapsStableScoreEvents(t *testing.T) {
	output := PromptfooOutputFromResult(CaseResult{
		Name:                                   "adversarial failure",
		IncidentID:                             "FIC-SYN-005",
		Kind:                                   CaseKindAdversarial,
		ExpectedSeverity:                       severity.LevelMedium,
		ActualSeverity:                         severity.LevelHigh,
		CitationCoverage:                       0.5,
		UnsupportedClaims:                      []string{"approved for export"},
		PromptInjectionResistant:               false,
		ExportApproved:                         true,
		SensitiveActionsBlockedWithoutApproval: false,
		Passed:                                 false,
	})

	events := ScoreEventsFromPromptfooOutput(output, ScoreEventMetadata{
		RunID:   "eval-run-13",
		TraceID: "trace-abc",
	})

	if len(events) != len(output.Scores) {
		t.Fatalf("events len = %d, want %d", len(events), len(output.Scores))
	}
	unsupported := findScoreEvent(t, events, "unsupported_claims")
	if unsupported.Name != "eval.score.unsupported_claims" {
		t.Fatalf("event name = %q, want stable scorer event name", unsupported.Name)
	}
	if unsupported.Score != 0 || unsupported.Pass || !unsupported.Critical {
		t.Fatalf("unsupported score event = %#v, want failed critical score 0", unsupported)
	}
	if unsupported.Severity != ScoreEventSeverityCritical {
		t.Fatalf("unsupported severity = %q, want critical", unsupported.Severity)
	}
	if unsupported.RunID != "eval-run-13" || unsupported.TraceID != "trace-abc" || unsupported.CaseID != "adversarial failure" || unsupported.IncidentID != "FIC-SYN-005" {
		t.Fatalf("trace/case correlation = %#v", unsupported)
	}

	severityEvent := findScoreEvent(t, events, "severity")
	if severityEvent.Name != "eval.score.severity" || severityEvent.Score != 0 || severityEvent.Pass {
		t.Fatalf("severity event = %#v, want failed severity score event", severityEvent)
	}
	if severityEvent.Severity != ScoreEventSeverityError {
		t.Fatalf("severity event severity = %q, want error", severityEvent.Severity)
	}
}

func TestNoopScoreEventExporterIsDefaultAndDisabledModeSkipsExporter(t *testing.T) {
	exporter := &recordingScoreEventExporter{}
	output := PromptfooOutputFromResult(passingScoreEventCaseResult())

	if err := ExportPromptfooScoreEvents(context.Background(), output, ScoreEventExportOptions{}); err != nil {
		t.Fatalf("default export error = %v, want nil no-op", err)
	}
	if err := ExportPromptfooScoreEvents(context.Background(), output, ScoreEventExportOptions{Exporter: failingScoreEventExporter{}}); err != nil {
		t.Fatalf("best-effort export error = %v, want nil", err)
	}
	if err := ExportPromptfooScoreEvents(context.Background(), output, ScoreEventExportOptions{
		Mode:     ScoreEventExportModeDisabled,
		Exporter: exporter,
	}); err != nil {
		t.Fatalf("disabled export error = %v, want nil", err)
	}
	if len(exporter.events) != 0 {
		t.Fatalf("disabled exporter received %d events, want none", len(exporter.events))
	}
}

func TestBestEffortScoreEventExporterFailureDoesNotBreakIncidentEvalTarget(t *testing.T) {
	target := NewIncidentEvalTarget(WithScoreEventExporter(failingScoreEventExporter{}))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/evalops/incident", strings.NewReader(`{"incident_id":"FIC-SYN-001","vars":{"trace_id":"trace-best-effort","eval_run_id":"run-best-effort"}}`))

	target.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s, want OK despite exporter failure", recorder.Code, recorder.Body.String())
	}
	response := decodeIncidentEvalTargetResponse(t, recorder)
	if !response.Output.Passed {
		t.Fatalf("output passed = false, want core incident workflow result preserved: %#v", response.Output)
	}
}

func TestReleaseGateScoreEventExporterFailureBreaksIncidentEvalTarget(t *testing.T) {
	target := NewIncidentEvalTarget(
		WithScoreEventExporter(failingScoreEventExporter{}),
		WithScoreEventExportMode(ScoreEventExportModeReleaseGate),
	)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/evalops/incident", strings.NewReader(`{"incident_id":"FIC-SYN-001","vars":{"trace_id":"trace-release-gate","eval_run_id":"run-release-gate"}}`))

	target.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("status = %d body = %s, want exporter failure to fail release gate", recorder.Code, recorder.Body.String())
	}
	response := decodeIncidentEvalTargetResponse(t, recorder)
	if response.Error.Code != "score_event_export_failed" {
		t.Fatalf("error code = %q, want score_event_export_failed", response.Error.Code)
	}
}

func passingScoreEventCaseResult() CaseResult {
	return CaseResult{
		Name:                                   "low severity hard brake",
		IncidentID:                             "FIC-SYN-001",
		Kind:                                   CaseKindNormal,
		ExpectedSeverity:                       severity.LevelLow,
		ActualSeverity:                         severity.LevelLow,
		CitationCoverage:                       1,
		PromptInjectionResistant:               true,
		SensitiveActionsBlockedWithoutApproval: true,
		Passed:                                 true,
	}
}

func findScoreEvent(t *testing.T, events []ScoreEvent, scorer string) ScoreEvent {
	t.Helper()

	for _, event := range events {
		if event.Scorer == scorer {
			return event
		}
	}
	t.Fatalf("score event for %q missing from %#v", scorer, events)
	return ScoreEvent{}
}

type recordingScoreEventExporter struct {
	events []ScoreEvent
}

func (exporter *recordingScoreEventExporter) ExportScoreEvent(ctx context.Context, event ScoreEvent) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	exporter.events = append(exporter.events, event)
	return nil
}

type failingScoreEventExporter struct{}

func (failingScoreEventExporter) ExportScoreEvent(context.Context, ScoreEvent) error {
	return errors.New("export unavailable")
}
