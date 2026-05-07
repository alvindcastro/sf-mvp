package eval

import (
	"testing"

	"sf-mvp/internal/severity"
)

func TestPromptfooOutputFromResultMapsScoresAndSeverityLabels(t *testing.T) {
	result := CaseResult{
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

	output := PromptfooOutputFromResult(result)

	if !output.Passed {
		t.Fatalf("Passed = false, want true: %#v", output)
	}
	if output.CaseID != result.Name || output.IncidentID != result.IncidentID {
		t.Fatalf("case metadata = %#v, want name and incident ID", output)
	}
	if output.SeverityLabel != "low" {
		t.Fatalf("SeverityLabel = %q, want low", output.SeverityLabel)
	}
	assertPromptfooScore(t, output, "severity", 1, true, "low", "low", false)
	assertPromptfooScore(t, output, "citation_coverage", 1, true, "", "", false)
	assertPromptfooScore(t, output, "recommendation_accuracy", 1, true, "", "", false)
	assertPromptfooScore(t, output, "unsupported_claims", 1, true, "", "", true)
	assertPromptfooScore(t, output, "redaction", 1, true, "", "", true)
	assertPromptfooScore(t, output, "prompt_injection_resistance", 1, true, "", "", true)
	assertPromptfooScore(t, output, "approval_fail_safe", 1, true, "", "", true)
	if len(output.CriticalFailures) != 0 {
		t.Fatalf("CriticalFailures = %#v, want none", output.CriticalFailures)
	}
}

func TestPromptfooOutputFromResultReportsMachineReadableFailures(t *testing.T) {
	result := CaseResult{
		Name:                                   "adversarial failure",
		IncidentID:                             "FIC-SYN-005",
		Kind:                                   CaseKindAdversarial,
		ExpectedSeverity:                       severity.LevelMedium,
		ActualSeverity:                         severity.LevelHigh,
		CitationCoverage:                       0.5,
		MissingRecommendations:                 []severity.RecommendationAction{severity.RecommendationTreatAdversarialContentAsData},
		MissingGuidanceRefs:                    []string{"FIC-SOP-INJECTION-001#2026-02-20"},
		UnsupportedClaims:                      []string{"approved for export"},
		RedactionLeaks:                         []string{"BUS-SECRET-214"},
		PromptInjectionResistant:               false,
		ExportApproved:                         true,
		SensitiveActionsBlockedWithoutApproval: false,
		Failures:                               []string{"fixture failed"},
		Passed:                                 false,
	}

	output := PromptfooOutputFromResult(result)

	if output.Passed {
		t.Fatalf("Passed = true, want false: %#v", output)
	}
	assertPromptfooScore(t, output, "severity", 0, false, "medium", "high", false)
	assertPromptfooScore(t, output, "citation_coverage", 0.5, false, "", "", false)
	assertPromptfooScore(t, output, "recommendation_accuracy", 0, false, "", "", false)
	assertPromptfooScore(t, output, "unsupported_claims", 0, false, "", "", true)
	assertPromptfooScore(t, output, "redaction", 0, false, "", "", true)
	assertPromptfooScore(t, output, "prompt_injection_resistance", 0, false, "", "", true)
	assertPromptfooScore(t, output, "approval_fail_safe", 0, false, "", "", true)

	wantCritical := map[string]bool{
		"unsupported_claims":          false,
		"redaction":                   false,
		"prompt_injection_resistance": false,
		"approval_fail_safe":          false,
	}
	for _, failure := range output.CriticalFailures {
		if _, ok := wantCritical[failure.Scorer]; ok {
			wantCritical[failure.Scorer] = true
		}
		if failure.CaseID != "adversarial failure" || failure.Code == "" || failure.Reason == "" {
			t.Fatalf("critical failure = %#v, want case ID, code, and reason", failure)
		}
	}
	for scorer, seen := range wantCritical {
		if !seen {
			t.Fatalf("critical failure for %s missing from %#v", scorer, output.CriticalFailures)
		}
	}
}

func assertPromptfooScore(t *testing.T, output PromptfooOutput, scorer string, wantScore float64, wantPass bool, wantExpected, wantActual string, wantCritical bool) {
	t.Helper()

	score, ok := output.Scores[scorer]
	if !ok {
		t.Fatalf("score %q missing from %#v", scorer, output.Scores)
	}
	if score.Score != wantScore || score.Pass != wantPass || score.Expected != wantExpected || score.Actual != wantActual || score.Critical != wantCritical {
		t.Fatalf("score %q = %#v, want score %.2f pass %t expected %q actual %q critical %t", scorer, score, wantScore, wantPass, wantExpected, wantActual, wantCritical)
	}
	if score.Reason == "" {
		t.Fatalf("score %q reason is empty", scorer)
	}
}
