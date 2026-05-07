package eval

import (
	"encoding/json"
	"strings"
	"testing"

	"sf-mvp/internal/ingestion"
)

func TestExportDraftCasesJSONLFromReviewSamplesMarksExpectedFieldsForReview(t *testing.T) {
	samples := []ReviewTraceSample{
		{
			CaseID:     "demo-stop-arm-review",
			Name:       "demo stop-arm review failure",
			IncidentID: "FIC-SYN-215",
			TraceID:    "trace-fic-syn-215-20260507t161500z-001",
			EventType:  ingestion.EventTypeStopArmConflict,
			Tags:       []string{"demo_rehearsal", "stop_arm"},
		},
	}

	got, err := ExportDraftCasesJSONL(samples)
	if err != nil {
		t.Fatalf("ExportDraftCasesJSONL returned error: %v", err)
	}

	records := decodeDraftCaseRecords(t, got)
	if len(records) != 1 {
		t.Fatalf("record count = %d, want 1", len(records))
	}
	record := records[0]
	if record.CaseID != "demo-stop-arm-review" {
		t.Fatalf("CaseID = %q, want sample case ID", record.CaseID)
	}
	if record.SourceTraceID != "trace-fic-syn-215-20260507t161500z-001" {
		t.Fatalf("SourceTraceID = %q, want trace preservation", record.SourceTraceID)
	}
	if !record.ReviewRequired || record.GateBlocking {
		t.Fatalf("review/gate flags = review_required:%t gate_blocking:%t, want true/false", record.ReviewRequired, record.GateBlocking)
	}
	if record.InputPacket.IncidentID != "FIC-SYN-215" || record.InputPacket.EventType != "stop_arm_conflict" {
		t.Fatalf("input packet = %#v, want redacted routing fields", record.InputPacket)
	}
	if record.Expected.Severity != DraftExpectedTODO {
		t.Fatalf("expected severity = %q, want TODO marker", record.Expected.Severity)
	}
	assertStringSetContains(t, record.Expected.Citations, DraftExpectedTODO)
	assertStringSetContains(t, record.Expected.Recommendations, DraftExpectedTODO)
	assertStringSetContains(t, record.Expected.ForbiddenClaims, DraftExpectedTODO)
	if !record.Expected.Approval.SensitiveActionsMustFailSafe {
		t.Fatal("sensitive approval expectation = false, want fail-safe default")
	}
	assertStringSetContains(t, record.Tags, "demo_rehearsal")
	assertStringSetContains(t, record.Tags, "stop_arm")
	assertStringSetContains(t, record.Tags, "review_required")
	assertStringSetContains(t, record.Tags, "draft")
}

func TestExportDraftCasesJSONLDeduplicatesByCaseIDAndPreservesFirstTrace(t *testing.T) {
	samples := []ReviewTraceSample{
		{
			CaseID:     "duplicate-demo-case",
			Name:       "first failure",
			IncidentID: "FIC-SYN-216",
			TraceID:    "trace-first",
			EventType:  ingestion.EventTypeHardBrake,
			Tags:       []string{"manual_review"},
		},
		{
			CaseID:     "duplicate-demo-case",
			Name:       "second failure",
			IncidentID: "FIC-SYN-217",
			TraceID:    "trace-second",
			EventType:  ingestion.EventTypeCollisionSignal,
			Tags:       []string{"demo_rehearsal"},
		},
	}

	got, err := ExportDraftCasesJSONL(samples)
	if err != nil {
		t.Fatalf("ExportDraftCasesJSONL returned error: %v", err)
	}

	records := decodeDraftCaseRecords(t, got)
	if len(records) != 1 {
		t.Fatalf("record count = %d, want deduplicated case", len(records))
	}
	if records[0].SourceTraceID != "trace-first" || records[0].InputPacket.IncidentID != "FIC-SYN-216" {
		t.Fatalf("deduplicated record = %#v, want first sample preserved", records[0])
	}
}

func TestExportDraftCasesJSONLRejectsInvalidSamples(t *testing.T) {
	tests := []struct {
		name    string
		samples []ReviewTraceSample
		want    string
	}{
		{
			name: "missing trace",
			samples: []ReviewTraceSample{{
				IncidentID: "FIC-SYN-218",
				EventType:  ingestion.EventTypeHardBrake,
			}},
			want: "trace_id is required",
		},
		{
			name: "non synthetic incident",
			samples: []ReviewTraceSample{{
				IncidentID: "REAL-218",
				TraceID:    "trace-real",
				EventType:  ingestion.EventTypeHardBrake,
			}},
			want: "incident_id must start with FIC-SYN-",
		},
		{
			name: "unsupported event type",
			samples: []ReviewTraceSample{{
				IncidentID: "FIC-SYN-218",
				TraceID:    "trace-bad-event",
				EventType:  ingestion.EventType("unsafe_live_event"),
			}},
			want: "event_type is not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ExportDraftCasesJSONL(tt.samples)
			if err == nil {
				t.Fatal("ExportDraftCasesJSONL returned nil error, want invalid sample error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want detail %q", err, tt.want)
			}
		})
	}
}

func TestExportDraftCasesJSONLDoesNotExposeReviewNotes(t *testing.T) {
	got, err := ExportDraftCasesJSONL([]ReviewTraceSample{{
		CaseID:              "privacy-review-case",
		IncidentID:          "FIC-SYN-219",
		TraceID:             "trace-fic-syn-219-20260507t162000z-001",
		EventType:           ingestion.EventTypeAdversarialNote,
		Tags:                []string{"demo_rehearsal", "demo_rehearsal", "prompt_injection"},
		RedactedReviewNotes: "Reviewer saw [REDACTED] in the model answer; raw vehicle BUS-SECRET-219 must stay out.",
	}})
	if err != nil {
		t.Fatalf("ExportDraftCasesJSONL returned error: %v", err)
	}

	exported := string(got)
	for _, leaked := range []string{"BUS-SECRET-219", "Reviewer saw", "[REDACTED]"} {
		if strings.Contains(exported, leaked) {
			t.Fatalf("draft JSONL leaked review note detail %q in %s", leaked, exported)
		}
	}
	records := decodeDraftCaseRecords(t, got)
	assertStringSetContains(t, records[0].Tags, "prompt_injection")
	if countString(records[0].Tags, "demo_rehearsal") != 1 {
		t.Fatalf("tags = %#v, want deduplicated inherited tags", records[0].Tags)
	}
}

func decodeDraftCaseRecords(t *testing.T, data []byte) []DraftCaseRecord {
	t.Helper()

	var records []DraftCaseRecord
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		var record DraftCaseRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("decode draft case JSONL line %q: %v", line, err)
		}
		records = append(records, record)
	}
	return records
}

func assertStringSetContains(t *testing.T, values []string, want string) {
	t.Helper()

	for _, value := range values {
		if value == want {
			return
		}
	}
	t.Fatalf("%q not found in %#v", want, values)
}

func countString(values []string, want string) int {
	count := 0
	for _, value := range values {
		if value == want {
			count++
		}
	}
	return count
}
