package eval

import (
	"context"
	"fmt"
	"sort"
)

const (
	ScoreEventExportModeBestEffort  ScoreEventExportMode = "best_effort"
	ScoreEventExportModeDisabled    ScoreEventExportMode = "disabled"
	ScoreEventExportModeReleaseGate ScoreEventExportMode = "release_gate"

	ScoreEventSeverityInfo     ScoreEventSeverity = "info"
	ScoreEventSeverityError    ScoreEventSeverity = "error"
	ScoreEventSeverityCritical ScoreEventSeverity = "critical"
)

type ScoreEventExportMode string

type ScoreEventSeverity string

type ScoreEventExporter interface {
	ExportScoreEvent(context.Context, ScoreEvent) error
}

type ScoreEvent struct {
	Name       string             `json:"name"`
	RunID      string             `json:"run_id,omitempty"`
	TraceID    string             `json:"trace_id,omitempty"`
	CaseID     string             `json:"case_id"`
	IncidentID string             `json:"incident_id"`
	Kind       CaseKind           `json:"kind,omitempty"`
	Scorer     string             `json:"scorer"`
	Score      float64            `json:"score"`
	Pass       bool               `json:"pass"`
	Critical   bool               `json:"critical"`
	Severity   ScoreEventSeverity `json:"severity"`
	Expected   string             `json:"expected,omitempty"`
	Actual     string             `json:"actual,omitempty"`
	Reason     string             `json:"reason"`
}

type ScoreEventMetadata struct {
	RunID   string
	TraceID string
}

type ScoreEventExportOptions struct {
	Mode     ScoreEventExportMode
	Exporter ScoreEventExporter
	Metadata ScoreEventMetadata
}

type noopScoreEventExporter struct{}

func (noopScoreEventExporter) ExportScoreEvent(context.Context, ScoreEvent) error {
	return nil
}

func ScoreEventsFromPromptfooOutput(output PromptfooOutput, metadata ScoreEventMetadata) []ScoreEvent {
	scorers := make([]string, 0, len(output.Scores))
	for scorer := range output.Scores {
		scorers = append(scorers, scorer)
	}
	sort.Strings(scorers)

	events := make([]ScoreEvent, 0, len(scorers))
	for _, scorer := range scorers {
		score := output.Scores[scorer]
		events = append(events, ScoreEvent{
			Name:       "eval.score." + scorer,
			RunID:      metadata.RunID,
			TraceID:    metadata.TraceID,
			CaseID:     output.CaseID,
			IncidentID: output.IncidentID,
			Kind:       output.Kind,
			Scorer:     scorer,
			Score:      score.Score,
			Pass:       score.Pass,
			Critical:   score.Critical,
			Severity:   scoreEventSeverity(score),
			Expected:   score.Expected,
			Actual:     score.Actual,
			Reason:     score.Reason,
		})
	}
	return events
}

func ExportPromptfooScoreEvents(ctx context.Context, output PromptfooOutput, options ScoreEventExportOptions) error {
	mode := options.Mode
	if mode == "" {
		mode = ScoreEventExportModeBestEffort
	}
	if mode == ScoreEventExportModeDisabled {
		return nil
	}
	exporter := options.Exporter
	if exporter == nil {
		exporter = noopScoreEventExporter{}
	}
	for _, event := range ScoreEventsFromPromptfooOutput(output, options.Metadata) {
		if err := exporter.ExportScoreEvent(ctx, event); err != nil {
			if mode == ScoreEventExportModeReleaseGate {
				return fmt.Errorf("export score event %q for case_id %q: %w", event.Scorer, event.CaseID, err)
			}
		}
	}
	return nil
}

func scoreEventSeverity(score PromptfooScore) ScoreEventSeverity {
	if score.Pass {
		return ScoreEventSeverityInfo
	}
	if score.Critical {
		return ScoreEventSeverityCritical
	}
	return ScoreEventSeverityError
}
