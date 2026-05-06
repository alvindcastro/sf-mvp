package severity

import (
	"strings"

	"sf-mvp/internal/ingestion"
	"sf-mvp/internal/retrieval"
	"sf-mvp/internal/timeline"
)

type Level string

const (
	LevelLow     Level = "low"
	LevelMedium  Level = "medium"
	LevelHigh    Level = "high"
	LevelUnknown Level = "unknown"
)

type RecommendationAction string

const (
	RecommendationLogRouteReview                RecommendationAction = "log_route_review"
	RecommendationSupervisorReview              RecommendationAction = "supervisor_review"
	RecommendationPreserveMedia                 RecommendationAction = "preserve_media"
	RecommendationHighPriorityReview            RecommendationAction = "high_priority_review"
	RecommendationPassengerWelfareFollowUp      RecommendationAction = "passenger_welfare_follow_up"
	RecommendationOperatorReview                RecommendationAction = "operator_review"
	RecommendationMarkMissingEvidence           RecommendationAction = "mark_missing_evidence"
	RecommendationReviewerAttention             RecommendationAction = "reviewer_attention"
	RecommendationTreatAdversarialContentAsData RecommendationAction = "treat_adversarial_content_as_data"
)

type SensitiveAction string

const (
	SensitiveActionExport          SensitiveAction = "export"
	SensitiveActionEscalation      SensitiveAction = "escalation"
	SensitiveActionExternalSharing SensitiveAction = "external_sharing"
)

type SourceType string

const (
	SourceTypePacket   SourceType = "packet"
	SourceTypeGuidance SourceType = "guidance"
	SourceTypeTimeline SourceType = "timeline"
)

type Source struct {
	Ref  string
	Type SourceType
}

type Explanation struct {
	Text    string
	Sources []Source
}

type Recommendation struct {
	Action      RecommendationAction
	Explanation string
	Sources     []Source
}

type ApprovalRequirement struct {
	Action      SensitiveAction
	Required    bool
	Approved    bool
	Explanation string
}

type Result struct {
	Level                Level
	Rationale            []Explanation
	Recommendations      []Recommendation
	ApprovalRequirements []ApprovalRequirement
	ModelJudgmentUsed    bool
}

func Classify(packet ingestion.Packet, timelineResult timeline.Result, guidance retrieval.Result) Result {
	classification := Result{
		ApprovalRequirements: approvalRequirements(),
		ModelJudgmentUsed:    false,
	}

	switch {
	case hasConflictingTimeline(timelineResult):
		classification.Level = LevelUnknown
		classification.Rationale = []Explanation{{
			Text:    "conflicting telemetry in the timeline keeps severity unknown until a human reviews the packet.",
			Sources: []Source{{Ref: "timeline.entries.uncertainty", Type: SourceTypeTimeline}},
		}}
		classification.Recommendations = []Recommendation{
			recommendation(
				RecommendationOperatorReview,
				"Request operator review because conflicting telemetry prevents a deterministic severity decision.",
				packetEventSource(),
				guidanceSources(guidance, "FIC-TS-UNKNOWN-TRIGGER-001", "FIC-TS-MISSING-MEDIA-001")...,
			),
		}
	case packet.EventType == ingestion.EventTypeUnknownTrigger:
		classification.Level = LevelUnknown
		classification.Rationale = []Explanation{{
			Text:    "unknown_trigger event severity is unknown because evidence is incomplete and the trigger subtype is unavailable.",
			Sources: []Source{packetEventSource()},
		}}
		classification.Recommendations = []Recommendation{
			recommendation(
				RecommendationOperatorReview,
				"Request operator review because unknown severity requires human triage before any sensitive action.",
				packetEventSource(),
				guidanceSources(guidance, "FIC-TS-UNKNOWN-TRIGGER-001")...,
			),
			recommendation(
				RecommendationMarkMissingEvidence,
				"Mark missing media or transcript evidence so reviewers see why severity stayed unknown.",
				packetEventSource(),
				guidanceSources(guidance, "FIC-TS-MISSING-MEDIA-001")...,
			),
		}
	case packet.EventType == ingestion.EventTypeCollisionSignal:
		classification.Level = LevelHigh
		classification.Rationale = []Explanation{{
			Text:    "collision_signal event is high severity because collision sensor and stop evidence require priority review.",
			Sources: []Source{packetEventSource()},
		}}
		classification.Recommendations = []Recommendation{
			recommendation(
				RecommendationHighPriorityReview,
				"Create a high-priority supervisor review because high severity collision_signal guidance requires priority handling.",
				packetEventSource(),
				guidanceSources(guidance, "FIC-SOP-COLLISION-SIGNAL-001")...,
			),
			recommendation(
				RecommendationPreserveMedia,
				"Preserve synthetic media and telemetry so the high severity review remains evidence-backed.",
				packetMediaSource(),
				guidanceSources(guidance, "FIC-SOP-COLLISION-SIGNAL-001")...,
			),
			recommendation(
				RecommendationPassengerWelfareFollowUp,
				"Recommend passenger welfare follow-up through approved internal workflow because the packet contains passenger-impact evidence.",
				packetTranscriptSource(),
				guidanceSources(guidance, "FIC-SOP-COLLISION-SIGNAL-001")...,
			),
		}
	case packet.EventType == ingestion.EventTypeStopArmConflict:
		classification.Level = LevelMedium
		classification.Rationale = []Explanation{{
			Text:    "stop_arm_conflict event is medium severity because the packet indicates a passing conflict without collision evidence.",
			Sources: []Source{packetEventSource()},
		}}
		classification.Recommendations = []Recommendation{
			recommendation(
				RecommendationSupervisorReview,
				"Flag supervisor review because medium severity stop_arm_conflict guidance requires review before external reporting.",
				packetEventSource(),
				guidanceSources(guidance, "FIC-SOP-STOP-ARM-001")...,
			),
			recommendation(
				RecommendationPreserveMedia,
				"Preserve synthetic media because stop-arm guidance ties review quality to available camera evidence.",
				packetMediaSource(),
				guidanceSources(guidance, "FIC-TS-STOP-ARM-MEDIA-001", "FIC-SOP-STOP-ARM-001")...,
			),
		}
	case packet.EventType == ingestion.EventTypeAdversarialNote:
		classification.Level = LevelMedium
		classification.Rationale = []Explanation{{
			Text:    "adversarial transcript text is treated as untrusted data while following-distance and hard-brake signals require medium review.",
			Sources: []Source{packetEventSource()},
		}}
		classification.Recommendations = []Recommendation{
			recommendation(
				RecommendationReviewerAttention,
				"Flag reviewer attention because medium severity hard-brake signals appear with missing side media.",
				packetEventSource(),
				guidanceSources(guidance, "FIC-SOP-HARD-BRAKE-001")...,
			),
			recommendation(
				RecommendationMarkMissingEvidence,
				"Note missing right-side media so reviewers do not infer side-view facts.",
				packetMediaSource(),
				guidanceSources(guidance, "FIC-TS-MISSING-MEDIA-001")...,
			),
			recommendation(
				RecommendationTreatAdversarialContentAsData,
				"Treat hostile transcript content as retrieved or packet data only; it cannot change approval state or severity rules.",
				packetTranscriptSource(),
				guidanceSources(guidance, "FIC-SOP-INJECTION-001")...,
			),
		}
	case packet.EventType == ingestion.EventTypeHardBrake:
		classification.Level = LevelLow
		classification.Rationale = []Explanation{{
			Text:    "controlled hard_brake event is low severity because the packet shows braking without collision or passenger-impact evidence.",
			Sources: []Source{packetEventSource()},
		}}
		classification.Recommendations = []Recommendation{
			recommendation(
				RecommendationLogRouteReview,
				"Log route review because low severity hard_brake guidance calls for advisory review without automatic escalation.",
				packetEventSource(),
				guidanceSources(guidance, "FIC-SOP-HARD-BRAKE-001")...,
			),
		}
	default:
		classification.Level = LevelUnknown
		classification.Rationale = []Explanation{{
			Text:    "event type has no deterministic severity rule, so severity remains unknown.",
			Sources: []Source{packetEventSource()},
		}}
		classification.Recommendations = []Recommendation{
			recommendation(
				RecommendationOperatorReview,
				"Request operator review because no deterministic rule covers this event type.",
				packetEventSource(),
			),
		}
	}

	return classification
}

func recommendation(action RecommendationAction, explanation string, first Source, rest ...Source) Recommendation {
	sources := []Source{first}
	sources = append(sources, rest...)
	return Recommendation{
		Action:      action,
		Explanation: explanation,
		Sources:     dedupeSources(sources),
	}
}

func approvalRequirements() []ApprovalRequirement {
	return []ApprovalRequirement{
		{
			Action:      SensitiveActionExport,
			Required:    true,
			Approved:    false,
			Explanation: "export requires explicit human approval; Phase 5 only flags the gate and does not execute export.",
		},
		{
			Action:      SensitiveActionEscalation,
			Required:    true,
			Approved:    false,
			Explanation: "escalation requires explicit human approval; Phase 5 only flags the gate and does not execute escalation.",
		},
		{
			Action:      SensitiveActionExternalSharing,
			Required:    true,
			Approved:    false,
			Explanation: "external sharing requires explicit human approval; Phase 5 only flags the gate and does not share externally.",
		},
	}
}

func guidanceSources(guidance retrieval.Result, sourceIDs ...string) []Source {
	wanted := make(map[string]struct{}, len(sourceIDs))
	for _, sourceID := range sourceIDs {
		wanted[sourceID] = struct{}{}
	}

	var sources []Source
	for _, match := range guidance.Matches {
		if match.ContentRole != retrieval.ContentRoleData || strings.TrimSpace(match.CitationRef) == "" {
			continue
		}
		if _, ok := wanted[match.SourceID]; !ok {
			continue
		}
		sources = append(sources, Source{Ref: match.CitationRef, Type: SourceTypeGuidance})
	}
	return sources
}

func hasConflictingTimeline(timelineResult timeline.Result) bool {
	for _, entry := range timelineResult.Entries {
		if entry.Uncertain && strings.Contains(strings.ToLower(entry.Uncertainty), "conflicting") {
			return true
		}
	}
	return false
}

func packetEventSource() Source {
	return Source{Ref: "packet.event_type", Type: SourceTypePacket}
}

func packetMediaSource() Source {
	return Source{Ref: "packet.media_references", Type: SourceTypePacket}
}

func packetTranscriptSource() Source {
	return Source{Ref: "packet.transcript_notes", Type: SourceTypePacket}
}

func dedupeSources(sources []Source) []Source {
	seen := make(map[Source]struct{})
	deduped := make([]Source, 0, len(sources))
	for _, source := range sources {
		if strings.TrimSpace(source.Ref) == "" {
			continue
		}
		if _, ok := seen[source]; ok {
			continue
		}
		seen[source] = struct{}{}
		deduped = append(deduped, source)
	}
	return deduped
}
