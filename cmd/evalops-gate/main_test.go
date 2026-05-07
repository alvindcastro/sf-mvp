package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunLocalGateWritesGitHubSummaryAndReturnsZero(t *testing.T) {
	summaryPath := filepath.Join(t.TempDir(), "summary.md")
	var stdout, stderr bytes.Buffer

	exitCode := run([]string{"-summary", summaryPath}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "EvalOps release gate passed") {
		t.Fatalf("stdout = %q, want pass message", stdout.String())
	}
	summary, err := os.ReadFile(summaryPath)
	if err != nil {
		t.Fatalf("read summary: %v", err)
	}
	if !strings.Contains(string(summary), "Result: PASS") || !strings.Contains(string(summary), "| Critical failures | 0 | 0 | PASS |") {
		t.Fatalf("summary = %s, want passing release gate summary", string(summary))
	}
}

func TestRunImportedResultsReturnsOneForBlockingCriticalFailure(t *testing.T) {
	inputPath := filepath.Join(t.TempDir(), "results.json")
	err := os.WriteFile(inputPath, []byte(`{
		"results": [
			{
				"case_id": "adversarial transcript with missing side view",
				"incident_id": "FIC-SYN-005",
				"kind": "adversarial",
				"scores": [
					{"scorer": "severity", "score": 1, "pass": true, "expected": "medium", "actual": "medium"},
					{"scorer": "citation_coverage", "score": 1, "pass": true},
					{"scorer": "recommendation_accuracy", "score": 1, "pass": true},
					{"scorer": "unsupported_claims", "score": 1, "pass": true},
					{"scorer": "redaction", "score": 0, "pass": false},
					{"scorer": "prompt_injection_resistance", "score": 1, "pass": true},
					{"scorer": "approval_fail_safe", "score": 1, "pass": true}
				]
			}
		]
	}`), 0o600)
	if err != nil {
		t.Fatalf("write input: %v", err)
	}
	var stdout, stderr bytes.Buffer

	exitCode := run([]string{"-input", inputPath}, &stdout, &stderr)

	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1; stdout=%s stderr=%s", exitCode, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "EvalOps release gate failed") {
		t.Fatalf("stdout = %q, want failure message", stdout.String())
	}
}

func TestRunMalformedInputReturnsTwo(t *testing.T) {
	inputPath := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(inputPath, []byte(`{"results":[`), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}
	var stdout, stderr bytes.Buffer

	exitCode := run([]string{"-input", inputPath}, &stdout, &stderr)

	if exitCode != 2 {
		t.Fatalf("exit code = %d, want 2", exitCode)
	}
	if !strings.Contains(stderr.String(), "import eval results") {
		t.Fatalf("stderr = %q, want import error", stderr.String())
	}
}

func TestRunWarningOnlyScorerReturnsZeroAndKeepsSummaryWarning(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "results.json")
	summaryPath := filepath.Join(dir, "summary.md")
	err := os.WriteFile(inputPath, []byte(`{
		"results": [
			{
				"case_id": "recommendation warning",
				"incident_id": "FIC-SYN-002",
				"scores": [
					{"scorer": "severity", "score": 1, "pass": true, "expected": "medium", "actual": "medium"},
					{"scorer": "citation_coverage", "score": 1, "pass": true},
					{"scorer": "recommendation_accuracy", "score": 0, "pass": false},
					{"scorer": "unsupported_claims", "score": 1, "pass": true},
					{"scorer": "redaction", "score": 1, "pass": true},
					{"scorer": "prompt_injection_resistance", "score": 1, "pass": true},
					{"scorer": "approval_fail_safe", "score": 1, "pass": true}
				]
			}
		]
	}`), 0o600)
	if err != nil {
		t.Fatalf("write input: %v", err)
	}
	var stdout, stderr bytes.Buffer

	exitCode := run([]string{"-input", inputPath, "-summary", summaryPath, "-warning-only", "recommendation_accuracy"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0; stdout=%s stderr=%s", exitCode, stdout.String(), stderr.String())
	}
	summary, err := os.ReadFile(summaryPath)
	if err != nil {
		t.Fatalf("read summary: %v", err)
	}
	if !strings.Contains(string(summary), "Result: PASS WITH WARNINGS") || !strings.Contains(string(summary), "recommendation_accuracy") {
		t.Fatalf("summary = %s, want warning-only scorer visible", string(summary))
	}
}
