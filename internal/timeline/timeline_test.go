package timeline

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"sf-mvp/internal/ingestion"
	"sf-mvp/internal/retrieval"
)

func TestBuildOrdersTelemetryChronologically(t *testing.T) {
	packet := timelinePacket()
	packet.TelemetrySamples = []ingestion.TelemetrySample{
		{RelativeTime: "+05s", SpeedMPH: 0, Heading: "northbound", Signal: "full stop", GPSLabel: "Oak St at Pine Ave"},
		{RelativeTime: "-06s", SpeedMPH: 24, Heading: "northbound", Signal: "steady speed", GPSLabel: "Oak St block 1200"},
		{RelativeTime: "+00s", SpeedMPH: 9, Heading: "northbound", Signal: "hard brake threshold crossed", GPSLabel: "Oak St at Pine Ave"},
	}

	result := Build(packet, retrieval.Result{})

	got := entryTimes(entriesWithSourcePrefix(result.Entries, "packet.telemetry_samples["))
	want := []string{"07:42:12", "07:42:18", "07:42:23"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("entry times = %#v, want %#v", got, want)
	}
	if !strings.Contains(result.Entries[0].Claim, "24 mph") {
		t.Fatalf("first claim = %q, want speed from earliest telemetry sample", result.Entries[0].Claim)
	}
}

func TestBuildCitesEveryTimelineClaim(t *testing.T) {
	result := Build(timelinePacket(), retrieval.Result{
		Matches: []retrieval.Citation{hardBrakeCitation()},
	})

	if len(result.Entries) == 0 {
		t.Fatal("timeline entries are empty")
	}
	for _, entry := range result.Entries {
		if len(entry.Sources) == 0 {
			t.Fatalf("entry %q has no sources", entry.Claim)
		}
		for _, source := range entry.Sources {
			if strings.TrimSpace(source.Ref) == "" {
				t.Fatalf("entry %q has empty source ref: %#v", entry.Claim, entry.Sources)
			}
		}
	}

	if len(result.GuidanceSources) != 1 {
		t.Fatalf("guidance source count = %d, want 1", len(result.GuidanceSources))
	}
	if result.GuidanceSources[0].Ref != "FIC-SOP-HARD-BRAKE-001#2026-02-15" {
		t.Fatalf("guidance source ref = %q, want citation ref", result.GuidanceSources[0].Ref)
	}
}

func TestBuildKeepsTranscriptAndStillFrameClaimsGrounded(t *testing.T) {
	result := Build(timelinePacket(), retrieval.Result{})

	var transcriptEntry, stillFrameEntry Entry
	for _, entry := range result.Entries {
		if hasSourcePrefix(entry, "packet.transcript_notes[0]") {
			transcriptEntry = entry
		}
		if hasSourcePrefix(entry, "packet.still_frame_notes[0]") {
			stillFrameEntry = entry
		}
	}

	if transcriptEntry.Claim == "" {
		t.Fatal("did not find transcript entry")
	}
	if stillFrameEntry.Claim == "" {
		t.Fatal("did not find still-frame entry")
	}
	if strings.Contains(transcriptEntry.Claim, "frame shows") {
		t.Fatalf("transcript claim = %q, want no invented visual fact", transcriptEntry.Claim)
	}
	if !strings.Contains(stillFrameEntry.Claim, "Still-frame note") {
		t.Fatalf("still-frame claim = %q, want still-frame attribution", stillFrameEntry.Claim)
	}
}

func TestBuildMarksMissingEvidenceAsUncertain(t *testing.T) {
	packet := timelinePacket()
	packet.MediaReferences = []string{
		"synthetic://fic-syn-004/rear-camera-052751-unavailable.jpg",
	}
	packet.TranscriptNotes = []string{"No driver note captured.", "Background audio is marked unavailable."}
	packet.StillFrameNotes = []string{"Rear frame unavailable.", "Side frame unavailable."}

	result := Build(packet, retrieval.Result{})

	var uncertain []Entry
	for _, entry := range result.Entries {
		if entry.Uncertain {
			uncertain = append(uncertain, entry)
		}
	}
	if len(uncertain) == 0 {
		t.Fatal("no uncertain entries for unavailable evidence")
	}
	if !strings.Contains(uncertain[0].Uncertainty, "unavailable") {
		t.Fatalf("uncertainty = %q, want unavailable evidence label", uncertain[0].Uncertainty)
	}
}

func TestBuildMarksConflictingPacketDataAsUncertain(t *testing.T) {
	packet := timelinePacket()
	packet.TelemetrySamples = []ingestion.TelemetrySample{
		{RelativeTime: "+00s", SpeedMPH: 0, Heading: "northbound", Signal: "full stop", GPSLabel: "Oak St at Pine Ave"},
		{RelativeTime: "+00s", SpeedMPH: 18, Heading: "northbound", Signal: "continued movement", GPSLabel: "Oak St at Pine Ave"},
	}

	result := Build(packet, retrieval.Result{})

	telemetryEntries := entriesWithSourcePrefix(result.Entries, "packet.telemetry_samples[")
	if len(telemetryEntries) != 2 {
		t.Fatalf("telemetry entry count = %d, want 2", len(telemetryEntries))
	}
	for _, entry := range telemetryEntries {
		if !entry.Uncertain {
			t.Fatalf("entry %#v is not marked uncertain for conflicting same-time telemetry", entry)
		}
		if !strings.Contains(entry.Uncertainty, "conflicting telemetry") {
			t.Fatalf("uncertainty = %q, want conflicting telemetry label", entry.Uncertainty)
		}
	}
}

func TestBuildOmitsUnsupportedClaims(t *testing.T) {
	result := Build(timelinePacket(), retrieval.Result{
		Matches: []retrieval.Citation{hardBrakeCitation()},
	})

	joined := strings.ToLower(strings.Join(entryClaims(result.Entries), " "))
	unsupported := []string{"collision", "injury", "approved", "exported", "plate"}
	for _, word := range unsupported {
		if strings.Contains(joined, word) {
			t.Fatalf("timeline claims include unsupported word %q: %s", word, joined)
		}
	}
}

func entryTimes(entries []Entry) []string {
	times := make([]string, len(entries))
	for i, entry := range entries {
		times[i] = entry.Time.Format("15:04:05")
	}
	return times
}

func entryClaims(entries []Entry) []string {
	claims := make([]string, len(entries))
	for i, entry := range entries {
		claims[i] = entry.Claim
	}
	return claims
}

func hasSourcePrefix(entry Entry, ref string) bool {
	for _, source := range entry.Sources {
		if source.Ref == ref {
			return true
		}
	}
	return false
}

func entriesWithSourcePrefix(entries []Entry, prefix string) []Entry {
	var filtered []Entry
	for _, entry := range entries {
		for _, source := range entry.Sources {
			if strings.HasPrefix(source.Ref, prefix) {
				filtered = append(filtered, entry)
				break
			}
		}
	}
	return filtered
}

func timelinePacket() ingestion.Packet {
	return ingestion.Packet{
		SyntheticRecord: true,
		IncidentID:      "FIC-SYN-001",
		VehicleID:       "BUS-214",
		Route:           "North Loop School Route 7",
		Timestamp:       time.Date(2026, time.March, 12, 7, 42, 18, 0, time.FixedZone("PDT", -7*60*60)),
		LocationLabel:   "Oak Street near Pine Avenue",
		EventType:       ingestion.EventTypeHardBrake,
		TelemetrySamples: []ingestion.TelemetrySample{
			{RelativeTime: "-06s", SpeedMPH: 24, Heading: "northbound", Signal: "steady speed", GPSLabel: "Oak St block 1200"},
			{RelativeTime: "+00s", SpeedMPH: 9, Heading: "northbound", Signal: "hard brake threshold crossed", GPSLabel: "Oak St at Pine Ave"},
			{RelativeTime: "+05s", SpeedMPH: 0, Heading: "northbound", Signal: "full stop", GPSLabel: "Oak St at Pine Ave"},
		},
		MediaReferences: []string{
			"synthetic://fic-syn-001/front-camera-074218.jpg",
		},
		TranscriptNotes: []string{
			"Driver says cyclist slowed near the crosswalk; no contact.",
		},
		StillFrameNotes: []string{
			"Front frame shows a cyclist ahead in the bike lane near a marked crosswalk.",
		},
	}
}

func hardBrakeCitation() retrieval.Citation {
	return retrieval.Citation{
		SourceID:     "FIC-SOP-HARD-BRAKE-001",
		Title:        "Hard-Brake Review SOP",
		Workflow:     "incident_review",
		Scope:        "tenant:fic-demo",
		RevisionDate: time.Date(2026, time.February, 15, 0, 0, 0, 0, time.UTC),
		CitationRef:  "FIC-SOP-HARD-BRAKE-001#2026-02-15",
		Snippet:      "For hard brake events near a crosswalk where no contact is reported, log the event for route review.",
		ContentRole:  retrieval.ContentRoleData,
	}
}
