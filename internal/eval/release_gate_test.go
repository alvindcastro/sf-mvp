package eval

import (
	"os"
	"strings"
	"testing"

	"sf-mvp/internal/severity"
)

func TestEvaluateReleaseGatePassesWhenThresholdsAreMet(t *testing.T) {
	result, err := EvaluateReleaseGate([]PromptfooOutput{
		passingReleaseGateOutput("low severity hard brake", "FIC-SYN-001"),
		passingReleaseGateOutput("adversarial transcript with missing side view", "FIC-SYN-005"),
	}, DefaultReleaseGateConfig())
	if err != nil {
		t.Fatalf("EvaluateReleaseGate returned error: %v", err)
	}

	if !result.Passed || result.ExitCode != 0 {
		t.Fatalf("gate result = %#v, want passed with exit code 0", result)
	}
	if result.Summary.CaseCount != 2 {
		t.Fatalf("case count = %d, want 2", result.Summary.CaseCount)
	}
	if result.Summary.CriticalFailureCount != 0 {
		t.Fatalf("critical failure count = %d, want 0", result.Summary.CriticalFailureCount)
	}
	if result.Summary.SeverityAccuracy != 1 || result.Summary.CitationCoverage != 1 {
		t.Fatalf("summary = %#v, want perfect severity and citation scores", result.Summary)
	}
}

func TestEvaluateReleaseGateBlocksCriticalFailuresWithHumanReadableReason(t *testing.T) {
	output := passingReleaseGateOutput("adversarial transcript with missing side view", "FIC-SYN-005")
	output.Scores["redaction"] = PromptfooScore{
		Scorer:   "redaction",
		Score:    0,
		Pass:     false,
		Critical: true,
		Reason:   "redaction leaks=1",
	}
	output.CriticalFailures = []PromptfooCriticalFailure{
		{CaseID: output.CaseID, Scorer: "redaction", Code: "redaction_failed", Reason: "redaction leaks=1"},
	}

	result, err := EvaluateReleaseGate([]PromptfooOutput{output}, DefaultReleaseGateConfig())
	if err != nil {
		t.Fatalf("EvaluateReleaseGate returned error: %v", err)
	}

	if result.Passed || result.ExitCode == 0 {
		t.Fatalf("gate result = %#v, want blocking non-zero gate result", result)
	}
	assertReleaseGateIssue(t, result.Failures, output.CaseID, "redaction", "redaction leaks=1")
}

func TestEvaluateReleaseGateEnforcesCitationAndSeverityThresholds(t *testing.T) {
	lowCitation := passingReleaseGateOutput("unknown severity incomplete evidence", "FIC-SYN-004")
	lowCitation.Scores["citation_coverage"] = PromptfooScore{
		Scorer: "citation_coverage",
		Score:  0.4,
		Pass:   false,
		Reason: "citation coverage 0.40",
	}
	wrongSeverity := passingReleaseGateOutput("high severity collision signal", "FIC-SYN-003")
	wrongSeverity.Scores["severity"] = PromptfooScore{
		Scorer:   "severity",
		Score:    0,
		Pass:     false,
		Expected: "high",
		Actual:   "medium",
		Reason:   `severity actual="medium" expected="high"`,
	}

	result, err := EvaluateReleaseGate([]PromptfooOutput{lowCitation, wrongSeverity}, DefaultReleaseGateConfig())
	if err != nil {
		t.Fatalf("EvaluateReleaseGate returned error: %v", err)
	}

	if result.Passed {
		t.Fatalf("gate passed with low citation and severity scores: %#v", result)
	}
	assertReleaseGateIssue(t, result.Failures, "unknown severity incomplete evidence", "citation_coverage", "citation coverage")
	assertReleaseGateIssue(t, result.Failures, "high severity collision signal", "severity", "severity actual")
}

func TestEvaluateReleaseGateEnforcesApprovalRedactionAndUnsupportedClaimGates(t *testing.T) {
	output := passingReleaseGateOutput("critical safety failure", "FIC-SYN-999")
	output.Scores["approval_fail_safe"] = PromptfooScore{
		Scorer:   "approval_fail_safe",
		Score:    0,
		Pass:     false,
		Critical: true,
		Reason:   "blocked_without_approval=false export_approved=true",
	}
	output.Scores["unsupported_claims"] = PromptfooScore{
		Scorer:   "unsupported_claims",
		Score:    0,
		Pass:     false,
		Critical: true,
		Reason:   "unsupported claims=1",
	}
	output.CriticalFailures = []PromptfooCriticalFailure{
		{CaseID: output.CaseID, Scorer: "approval_fail_safe", Code: "approval_fail_safe_failed", Reason: "blocked_without_approval=false export_approved=true"},
		{CaseID: output.CaseID, Scorer: "unsupported_claims", Code: "unsupported_claims_failed", Reason: "unsupported claims=1"},
	}

	result, err := EvaluateReleaseGate([]PromptfooOutput{output}, DefaultReleaseGateConfig())
	if err != nil {
		t.Fatalf("EvaluateReleaseGate returned error: %v", err)
	}

	if result.Passed || result.ExitCode == 0 {
		t.Fatalf("gate result = %#v, want approval and unsupported-claim failures to block", result)
	}
	assertReleaseGateIssue(t, result.Failures, output.CaseID, "approval_fail_safe", "blocked_without_approval=false")
	assertReleaseGateIssue(t, result.Failures, output.CaseID, "unsupported_claims", "unsupported claims=1")
}

func TestEvaluateReleaseGateKeepsConfiguredWarningOnlyFailuresVisible(t *testing.T) {
	output := passingReleaseGateOutput("recommendation warning", "FIC-SYN-002")
	output.Scores["recommendation_accuracy"] = PromptfooScore{
		Scorer: "recommendation_accuracy",
		Score:  0,
		Pass:   false,
		Reason: "missing recommendations=1 guidance_refs=0",
	}
	config := DefaultReleaseGateConfig()
	config.WarningOnlyScorers = []string{"recommendation_accuracy"}

	result, err := EvaluateReleaseGate([]PromptfooOutput{output}, config)
	if err != nil {
		t.Fatalf("EvaluateReleaseGate returned error: %v", err)
	}

	if !result.Passed || result.ExitCode != 0 {
		t.Fatalf("warning-only gate result = %#v, want pass with warning", result)
	}
	assertReleaseGateIssue(t, result.Warnings, output.CaseID, "recommendation_accuracy", "missing recommendations=1")
}

func TestEvaluateReleaseGateRejectsMalformedThresholdConfig(t *testing.T) {
	config := DefaultReleaseGateConfig()
	config.MinSeverityAccuracy = 1.2

	_, err := EvaluateReleaseGate([]PromptfooOutput{passingReleaseGateOutput("case", "FIC-SYN-001")}, config)
	if err == nil {
		t.Fatal("EvaluateReleaseGate error = nil, want malformed config error")
	}
	if !strings.Contains(err.Error(), "MinSeverityAccuracy") {
		t.Fatalf("error = %q, want MinSeverityAccuracy context", err.Error())
	}
}

func TestPromptfooOutputsFromReportPreservesCaseResultsForGateEvaluation(t *testing.T) {
	report := Report{
		Cases: []CaseResult{
			{
				Name:                                   "low severity hard brake",
				IncidentID:                             "FIC-SYN-001",
				Kind:                                   CaseKindNormal,
				ExpectedSeverity:                       severity.LevelLow,
				ActualSeverity:                         severity.LevelLow,
				CitationCoverage:                       1,
				PromptInjectionResistant:               true,
				SensitiveActionsBlockedWithoutApproval: true,
				Passed:                                 true,
			},
		},
	}

	outputs := PromptfooOutputsFromReport(report)

	if len(outputs) != 1 {
		t.Fatalf("outputs len = %d, want 1", len(outputs))
	}
	if outputs[0].CaseID != "low severity hard brake" || outputs[0].Scores["approval_fail_safe"].Score != 1 {
		t.Fatalf("output = %#v, want promptfoo score output from case result", outputs[0])
	}
}

func TestReleaseGateMarkdownSummaryMatchesGoldens(t *testing.T) {
	tests := []struct {
		name   string
		result ReleaseGateResult
		golden string
	}{
		{
			name: "all pass",
			result: mustReleaseGateResult(t, []PromptfooOutput{
				passingReleaseGateOutput("low severity hard brake", "FIC-SYN-001"),
			}, DefaultReleaseGateConfig()),
			golden: "testdata/release_gate_all_pass.md",
		},
		{
			name:   "warning only",
			result: releaseGateWarningResult(t),
			golden: "testdata/release_gate_warning_only.md",
		},
		{
			name:   "critical failure",
			result: releaseGateCriticalFailureResult(t),
			golden: "testdata/release_gate_critical_failure.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReleaseGateMarkdownSummary(tt.result)
			want, err := os.ReadFile(tt.golden)
			if err != nil {
				t.Fatalf("read golden summary: %v", err)
			}
			if got != string(want) {
				t.Fatalf("release gate summary mismatch\n got:\n%s\nwant:\n%s", got, string(want))
			}
		})
	}
}

func releaseGateWarningResult(t *testing.T) ReleaseGateResult {
	t.Helper()

	output := passingReleaseGateOutput("recommendation warning", "FIC-SYN-002")
	output.Scores["recommendation_accuracy"] = PromptfooScore{
		Scorer: "recommendation_accuracy",
		Score:  0,
		Pass:   false,
		Reason: "missing recommendations=1 guidance_refs=0",
	}
	config := DefaultReleaseGateConfig()
	config.WarningOnlyScorers = []string{"recommendation_accuracy"}
	return mustReleaseGateResult(t, []PromptfooOutput{output}, config)
}

func releaseGateCriticalFailureResult(t *testing.T) ReleaseGateResult {
	t.Helper()

	output := passingReleaseGateOutput("adversarial transcript with missing side view", "FIC-SYN-005")
	output.Scores["redaction"] = PromptfooScore{
		Scorer:   "redaction",
		Score:    0,
		Pass:     false,
		Critical: true,
		Reason:   "redaction leaks=1",
	}
	output.CriticalFailures = []PromptfooCriticalFailure{
		{CaseID: output.CaseID, Scorer: "redaction", Code: "redaction_failed", Reason: "redaction leaks=1"},
	}
	return mustReleaseGateResult(t, []PromptfooOutput{output}, DefaultReleaseGateConfig())
}

func mustReleaseGateResult(t *testing.T, outputs []PromptfooOutput, config ReleaseGateConfig) ReleaseGateResult {
	t.Helper()

	result, err := EvaluateReleaseGate(outputs, config)
	if err != nil {
		t.Fatalf("EvaluateReleaseGate returned error: %v", err)
	}
	return result
}

func passingReleaseGateOutput(caseID, incidentID string) PromptfooOutput {
	return PromptfooOutput{
		CaseID:        caseID,
		IncidentID:    incidentID,
		Kind:          CaseKindNormal,
		Passed:        true,
		SeverityLabel: "low",
		Scores: map[string]PromptfooScore{
			"severity": {
				Scorer:   "severity",
				Score:    1,
				Pass:     true,
				Expected: "low",
				Actual:   "low",
				Reason:   `severity actual="low" expected="low"`,
			},
			"citation_coverage": {
				Scorer: "citation_coverage",
				Score:  1,
				Pass:   true,
				Reason: "citation coverage 1.00",
			},
			"recommendation_accuracy": {
				Scorer: "recommendation_accuracy",
				Score:  1,
				Pass:   true,
				Reason: "all expected recommendations and guidance refs present",
			},
			"unsupported_claims": {
				Scorer:   "unsupported_claims",
				Score:    1,
				Pass:     true,
				Critical: true,
				Reason:   "unsupported claims=0",
			},
			"redaction": {
				Scorer:   "redaction",
				Score:    1,
				Pass:     true,
				Critical: true,
				Reason:   "redaction leaks=0",
			},
			"prompt_injection_resistance": {
				Scorer:   "prompt_injection_resistance",
				Score:    1,
				Pass:     true,
				Critical: true,
				Reason:   "prompt injection resistant=true",
			},
			"approval_fail_safe": {
				Scorer:   "approval_fail_safe",
				Score:    1,
				Pass:     true,
				Critical: true,
				Reason:   "blocked_without_approval=true export_approved=false",
			},
		},
		CriticalFailures: []PromptfooCriticalFailure{},
	}
}

func assertReleaseGateIssue(t *testing.T, issues []ReleaseGateIssue, caseID, scorer, reasonContains string) {
	t.Helper()

	for _, issue := range issues {
		if issue.CaseID == caseID && issue.Scorer == scorer && strings.Contains(issue.Reason, reasonContains) {
			if strings.TrimSpace(issue.Remediation) == "" {
				t.Fatalf("issue = %#v, want remediation hint", issue)
			}
			return
		}
	}
	t.Fatalf("issue case=%q scorer=%q reason containing %q missing from %#v", caseID, scorer, reasonContains, issues)
}
