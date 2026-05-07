package eval

import (
	"strings"
	"testing"

	"sf-mvp/internal/severity"
)

func TestImportSharedResultsValidPromptfooStyleResult(t *testing.T) {
	data := []byte(`{
		"results": {
			"results": [
				{
					"vars": {
						"case_id": "low severity hard brake",
						"incident_id": "FIC-SYN-001",
						"kind": "normal"
					},
					"assertionResults": [
						{"metric": "severity", "score": 1, "pass": true, "expected": "low", "actual": "low"},
						{"metric": "citation_coverage", "score": 1, "pass": true},
						{"metric": "recommendation_accuracy", "score": 1, "pass": true},
						{"metric": "unsupported_claims", "score": 1, "pass": true},
						{"metric": "redaction", "score": 1, "pass": true},
						{"metric": "prompt_injection_resistance", "score": 1, "pass": true},
						{"metric": "approval_fail_safe", "score": 1, "pass": true}
					]
				}
			]
		}
	}`)

	results, err := ImportSharedResultsJSON(data)
	if err != nil {
		t.Fatalf("ImportSharedResultsJSON returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("imported result count = %d, want 1", len(results))
	}

	result := results[0]
	if result.Name != "low severity hard brake" {
		t.Fatalf("Name = %q, want case id", result.Name)
	}
	if result.IncidentID != "FIC-SYN-001" {
		t.Fatalf("IncidentID = %q, want FIC-SYN-001", result.IncidentID)
	}
	if result.Kind != CaseKindNormal {
		t.Fatalf("Kind = %q, want normal", result.Kind)
	}
	if result.ExpectedSeverity != severity.LevelLow || result.ActualSeverity != severity.LevelLow {
		t.Fatalf("severity = actual %q expected %q, want low/low", result.ActualSeverity, result.ExpectedSeverity)
	}
	if result.CitationCoverage != 1 {
		t.Fatalf("CitationCoverage = %.2f, want 1.00", result.CitationCoverage)
	}
	if !result.PromptInjectionResistant {
		t.Fatal("PromptInjectionResistant = false, want true")
	}
	if !result.SensitiveActionsBlockedWithoutApproval {
		t.Fatal("SensitiveActionsBlockedWithoutApproval = false, want true")
	}
	if !result.Passed {
		t.Fatalf("Passed = false, failures: %#v", result.Failures)
	}
}

func TestImportSharedResultsRejectsMissingCaseID(t *testing.T) {
	data := []byte(`{"results":[{"incident_id":"FIC-SYN-001","scores":[{"scorer":"severity","score":1,"pass":true}]}]}`)

	_, err := ImportSharedResultsJSON(data)
	if err == nil {
		t.Fatal("ImportSharedResultsJSON error = nil, want missing case_id error")
	}
	assertErrorContains(t, err, "case_id")
}

func TestImportSharedResultsRejectsMalformedScore(t *testing.T) {
	data := []byte(`{
		"results": [
			{
				"case_id": "bad score",
				"incident_id": "FIC-SYN-001",
				"scores": [{"scorer": "citation_coverage", "score": "perfect", "pass": true}]
			}
		]
	}`)

	_, err := ImportSharedResultsJSON(data)
	if err == nil {
		t.Fatal("ImportSharedResultsJSON error = nil, want malformed score error")
	}
	assertErrorContains(t, err, "citation_coverage")
	assertErrorContains(t, err, "score")
}

func TestImportSharedResultsRejectsUnknownScorer(t *testing.T) {
	data := []byte(`{
		"results": [
			{
				"case_id": "unknown scorer",
				"incident_id": "FIC-SYN-001",
				"scores": [{"scorer": "fluency", "score": 1, "pass": true}]
			}
		]
	}`)

	_, err := ImportSharedResultsJSON(data)
	if err == nil {
		t.Fatal("ImportSharedResultsJSON error = nil, want unknown scorer error")
	}
	assertErrorContains(t, err, "unknown scorer")
	assertErrorContains(t, err, "fluency")
}

func TestImportSharedResultsRejectsDuplicateResult(t *testing.T) {
	data := []byte(`{
		"results": [
			{"case_id": "duplicate", "incident_id": "FIC-SYN-001", "scores": [{"scorer": "citation_coverage", "score": 1, "pass": true}]},
			{"case_id": "duplicate", "incident_id": "FIC-SYN-001", "scores": [{"scorer": "citation_coverage", "score": 1, "pass": true}]}
		]
	}`)

	_, err := ImportSharedResultsJSON(data)
	if err == nil {
		t.Fatal("ImportSharedResultsJSON error = nil, want duplicate result error")
	}
	assertErrorContains(t, err, "duplicate")
}

func TestImportSharedResultsPreservesRecommendationFailureForGateScoring(t *testing.T) {
	data := []byte(`{
		"results": [
			{
				"case_id": "recommendation warning",
				"incident_id": "FIC-SYN-002",
				"scores": [{"scorer": "recommendation_accuracy", "score": 0, "pass": false}]
			}
		]
	}`)

	results, err := ImportSharedResultsJSON(data)
	if err != nil {
		t.Fatalf("ImportSharedResultsJSON returned error: %v", err)
	}
	output := PromptfooOutputFromResult(results[0])

	score := output.Scores["recommendation_accuracy"]
	if score.Pass || score.Score != 0 {
		t.Fatalf("recommendation score = %#v, want imported failure preserved", score)
	}
}

func assertErrorContains(t *testing.T, err error, want string) {
	t.Helper()

	if !strings.Contains(err.Error(), want) {
		t.Fatalf("error = %q, want substring %q", err.Error(), want)
	}
}
