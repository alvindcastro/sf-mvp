package timeline

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"sf-mvp/internal/ingestion"
	"sf-mvp/internal/retrieval"
)

type Result struct {
	Entries         []Entry
	GuidanceSources []Source
}

type Entry struct {
	Time        time.Time
	Claim       string
	Sources     []Source
	Uncertain   bool
	Uncertainty string
}

type Source struct {
	Ref  string
	Type SourceType
}

type SourceType string

const (
	SourceTypePacket    SourceType = "packet"
	SourceTypeGuidance  SourceType = "guidance"
	SourceTypeTelemetry SourceType = "telemetry"
)

func Build(packet ingestion.Packet, guidance retrieval.Result) Result {
	entries := telemetryEntries(packet)
	entries = append(entries, transcriptEntries(packet)...)
	entries = append(entries, stillFrameEntries(packet)...)
	entries = append(entries, mediaAvailabilityEntries(packet)...)

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Time.Before(entries[j].Time)
	})

	return Result{
		Entries:         entries,
		GuidanceSources: guidanceSources(guidance),
	}
}

func telemetryEntries(packet ingestion.Packet) []Entry {
	entries := make([]Entry, 0, len(packet.TelemetrySamples))
	for i, sample := range packet.TelemetrySamples {
		eventTime, uncertainty := sampleTime(packet.Timestamp, sample.RelativeTime)
		entry := Entry{
			Time:  eventTime,
			Claim: telemetryClaim(sample),
			Sources: []Source{{
				Ref:  fmt.Sprintf("packet.telemetry_samples[%d]", i),
				Type: SourceTypeTelemetry,
			}},
		}
		if uncertainty != "" {
			entry.Uncertain = true
			entry.Uncertainty = uncertainty
		}
		entries = append(entries, entry)
	}

	markTelemetryConflicts(entries)
	return entries
}

func transcriptEntries(packet ingestion.Packet) []Entry {
	entries := make([]Entry, 0, len(packet.TranscriptNotes))
	for i, note := range packet.TranscriptNotes {
		entry := Entry{
			Time:  packet.Timestamp,
			Claim: fmt.Sprintf("Transcript note: %s", strings.TrimSpace(note)),
			Sources: []Source{{
				Ref:  fmt.Sprintf("packet.transcript_notes[%d]", i),
				Type: SourceTypePacket,
			}},
		}
		if unavailable(note) {
			entry.Uncertain = true
			entry.Uncertainty = "evidence unavailable"
		}
		entries = append(entries, entry)
	}
	return entries
}

func stillFrameEntries(packet ingestion.Packet) []Entry {
	entries := make([]Entry, 0, len(packet.StillFrameNotes))
	for i, note := range packet.StillFrameNotes {
		entry := Entry{
			Time:  packet.Timestamp,
			Claim: fmt.Sprintf("Still-frame note: %s", strings.TrimSpace(note)),
			Sources: []Source{{
				Ref:  fmt.Sprintf("packet.still_frame_notes[%d]", i),
				Type: SourceTypePacket,
			}},
		}
		if unavailable(note) {
			entry.Uncertain = true
			entry.Uncertainty = "visual evidence unavailable"
		}
		entries = append(entries, entry)
	}
	return entries
}

func mediaAvailabilityEntries(packet ingestion.Packet) []Entry {
	var entries []Entry
	for i, mediaReference := range packet.MediaReferences {
		if !unavailable(mediaReference) {
			continue
		}
		entries = append(entries, Entry{
			Time:        packet.Timestamp,
			Claim:       fmt.Sprintf("Media reference unavailable: %s", mediaReference),
			Sources:     []Source{{Ref: fmt.Sprintf("packet.media_references[%d]", i), Type: SourceTypePacket}},
			Uncertain:   true,
			Uncertainty: "media evidence unavailable",
		})
	}
	return entries
}

func sampleTime(base time.Time, relativeTime string) (time.Time, string) {
	offset, err := time.ParseDuration(strings.TrimSpace(relativeTime))
	if err != nil {
		return base, "telemetry relative_time unavailable"
	}
	return base.Add(offset), ""
}

func telemetryClaim(sample ingestion.TelemetrySample) string {
	return fmt.Sprintf(
		"Telemetry shows %s at %.0f mph, heading %s near %s.",
		strings.TrimSpace(sample.Signal),
		sample.SpeedMPH,
		strings.TrimSpace(sample.Heading),
		strings.TrimSpace(sample.GPSLabel),
	)
}

func markTelemetryConflicts(entries []Entry) {
	for i := range entries {
		for j := i + 1; j < len(entries); j++ {
			if !entries[i].Time.Equal(entries[j].Time) {
				continue
			}
			if entries[i].Claim == entries[j].Claim {
				continue
			}
			entries[i].Uncertain = true
			entries[i].Uncertainty = appendUncertainty(entries[i].Uncertainty, "conflicting telemetry at same timestamp")
			entries[j].Uncertain = true
			entries[j].Uncertainty = appendUncertainty(entries[j].Uncertainty, "conflicting telemetry at same timestamp")
		}
	}
}

func appendUncertainty(existing, next string) string {
	if existing == "" {
		return next
	}
	if strings.Contains(existing, next) {
		return existing
	}
	return existing + "; " + next
}

func guidanceSources(guidance retrieval.Result) []Source {
	var sources []Source
	for _, match := range guidance.Matches {
		if match.ContentRole != retrieval.ContentRoleData || strings.TrimSpace(match.CitationRef) == "" {
			continue
		}
		sources = append(sources, Source{
			Ref:  match.CitationRef,
			Type: SourceTypeGuidance,
		})
	}
	return sources
}

func unavailable(value string) bool {
	return strings.Contains(strings.ToLower(value), "unavailable")
}
