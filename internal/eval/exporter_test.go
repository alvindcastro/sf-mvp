package eval

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestExportCasesJSONLMatchesGoldenWithoutSensitiveEvidence(t *testing.T) {
	cases := GoldenCases()
	reverseCases(cases)

	got, err := ExportCasesJSONL(cases)
	if err != nil {
		t.Fatalf("ExportCasesJSONL returned error: %v", err)
	}

	want, err := os.ReadFile("testdata/evalops_cases.golden.jsonl")
	if err != nil {
		t.Fatalf("read golden JSONL: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("exported JSONL mismatch\n got:\n%s\nwant:\n%s", got, want)
	}

	exported := string(got)
	for _, evalCase := range GoldenCases() {
		if strings.Contains(exported, evalCase.Packet.VehicleID) {
			t.Fatalf("export contains vehicle ID for %s", evalCase.Packet.IncidentID)
		}
		if strings.Contains(exported, evalCase.Packet.Route) {
			t.Fatalf("export contains route for %s", evalCase.Packet.IncidentID)
		}
		if strings.Contains(exported, evalCase.Packet.LocationLabel) {
			t.Fatalf("export contains location label for %s", evalCase.Packet.IncidentID)
		}
		for _, value := range append(append([]string{}, evalCase.Packet.TranscriptNotes...), evalCase.Packet.StillFrameNotes...) {
			if strings.Contains(exported, value) {
				t.Fatalf("export contains raw evidence note for %s: %q", evalCase.Packet.IncidentID, value)
			}
		}
		for _, value := range evalCase.Packet.MediaReferences {
			if strings.Contains(exported, value) {
				t.Fatalf("export contains raw media reference for %s: %q", evalCase.Packet.IncidentID, value)
			}
		}
		for _, value := range evalCase.Expected.SensitiveTerms {
			if strings.Contains(exported, value) {
				t.Fatalf("export contains expected sensitive term for %s: %q", evalCase.Packet.IncidentID, value)
			}
		}
	}
}

func reverseCases(cases []Case) {
	for left, right := 0, len(cases)-1; left < right; left, right = left+1, right-1 {
		cases[left], cases[right] = cases[right], cases[left]
	}
}
