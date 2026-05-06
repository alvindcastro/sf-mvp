package brief

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"sf-mvp/internal/ingestion"
	"sf-mvp/internal/severity"
	"sf-mvp/internal/timeline"
)

type Status string

const StatusDraft Status = "draft"

type SourceType string

const (
	SourceTypePacket   SourceType = "packet"
	SourceTypeTimeline SourceType = "timeline"
	SourceTypeSeverity SourceType = "severity"
	SourceTypeGuidance SourceType = "guidance"
)

type Source struct {
	Ref  string
	Type SourceType
}

type Section struct {
	Title   string
	Text    string
	Sources []Source
}

type Result struct {
	Status            Status
	IncidentID        string
	SyntheticRecord   bool
	Sections          []Section
	ApprovalState     []ApprovalState
	RedactionsApplied []Redaction
	Uncertainties     []string
	Shareable         bool
}

type ApprovalState struct {
	Action  severity.SensitiveAction
	Blocked bool
	Reason  string
}

type Redaction struct {
	Field  string
	Reason string
}

type MissingEvidenceError struct {
	Reason string
}

func (e MissingEvidenceError) Error() string {
	if strings.TrimSpace(e.Reason) == "" {
		return "incident brief draft failed closed because required evidence is missing"
	}
	return "incident brief draft failed closed: " + e.Reason
}

func Draft(packet ingestion.Packet, timelineResult timeline.Result, severityResult severity.Result) (Result, error) {
	result := Result{
		Status:          StatusDraft,
		IncidentID:      packet.IncidentID,
		SyntheticRecord: packet.SyntheticRecord,
		Shareable:       false,
	}

	if err := validateRequiredEvidence(packet, timelineResult, severityResult); err != nil {
		return result, err
	}

	redactor := newRedactor(packet)
	approvalState := approvalStates(severityResult)
	sections := []Section{
		summarySection(packet),
		timelineSection(timelineResult),
		severitySection(severityResult),
		recommendationsSection(severityResult),
		approvalSection(severityResult),
	}
	for i := range sections {
		sections[i].Text = redactor.redact(sections[i].Text)
	}

	result.Sections = sections
	result.ApprovalState = approvalState
	result.RedactionsApplied = redactor.redactions
	result.Uncertainties = uncertaintyLabels(timelineResult)
	result.Shareable = isShareableDraft(result)
	return result, nil
}

func validateRequiredEvidence(packet ingestion.Packet, timelineResult timeline.Result, severityResult severity.Result) error {
	if strings.TrimSpace(packet.IncidentID) == "" {
		return MissingEvidenceError{Reason: "packet incident_id is required"}
	}
	if len(timelineResult.Entries) == 0 {
		return MissingEvidenceError{Reason: "timeline entries are required"}
	}
	for i, entry := range timelineResult.Entries {
		if strings.TrimSpace(entry.Claim) == "" {
			return MissingEvidenceError{Reason: fmt.Sprintf("timeline entry %d claim is required", i)}
		}
		if !hasRefs(convertTimelineSources(entry.Sources)) {
			return MissingEvidenceError{Reason: fmt.Sprintf("timeline entry %d requires citations", i)}
		}
	}
	if severityResult.Level == "" {
		return MissingEvidenceError{Reason: "severity level is required"}
	}
	if len(severityResult.Rationale) == 0 {
		return MissingEvidenceError{Reason: "severity rationale is required"}
	}
	for i, explanation := range severityResult.Rationale {
		if strings.TrimSpace(explanation.Text) == "" {
			return MissingEvidenceError{Reason: fmt.Sprintf("severity rationale %d text is required", i)}
		}
		if !hasRefs(convertSeveritySources(explanation.Sources, SourceTypeSeverity)) {
			return MissingEvidenceError{Reason: fmt.Sprintf("severity rationale %d requires citations", i)}
		}
	}
	if len(severityResult.Recommendations) == 0 {
		return MissingEvidenceError{Reason: "recommendations are required"}
	}
	for i, recommendation := range severityResult.Recommendations {
		if strings.TrimSpace(string(recommendation.Action)) == "" || strings.TrimSpace(recommendation.Explanation) == "" {
			return MissingEvidenceError{Reason: fmt.Sprintf("recommendation %d action and explanation are required", i)}
		}
		if !hasRefs(convertSeveritySources(recommendation.Sources, SourceTypeSeverity)) {
			return MissingEvidenceError{Reason: fmt.Sprintf("recommendation %d requires citations", i)}
		}
	}
	if len(severityResult.ApprovalRequirements) == 0 {
		return MissingEvidenceError{Reason: "approval requirements are required"}
	}
	return nil
}

func summarySection(packet ingestion.Packet) Section {
	text := fmt.Sprintf(
		"Incident %s is a synthetic %s draft created for human review at %s. Vehicle, route, and location details are redacted for shareable review.",
		packet.IncidentID,
		packet.EventType,
		packet.Timestamp.Format(time.RFC3339),
	)
	return Section{
		Title: "Incident Summary",
		Text:  text,
		Sources: []Source{
			{Ref: "packet.incident_id", Type: SourceTypePacket},
			{Ref: "packet.synthetic_record", Type: SourceTypePacket},
			{Ref: "packet.event_type", Type: SourceTypePacket},
			{Ref: "packet.timestamp", Type: SourceTypePacket},
		},
	}
}

func timelineSection(timelineResult timeline.Result) Section {
	var lines []string
	var sources []Source
	for _, entry := range timelineResult.Entries {
		line := fmt.Sprintf(
			"%s: %s Sources: %s.",
			entry.Time.Format(time.RFC3339),
			strings.TrimSpace(entry.Claim),
			strings.Join(sourceRefs(convertTimelineSources(entry.Sources)), ", "),
		)
		if entry.Uncertain {
			line += " Uncertainty: " + strings.TrimSpace(entry.Uncertainty) + "."
		}
		lines = append(lines, line)
		sources = append(sources, convertTimelineSources(entry.Sources)...)
	}
	return Section{
		Title:   "Cited Timeline",
		Text:    strings.Join(lines, "\n"),
		Sources: dedupeSources(sources),
	}
}

func severitySection(severityResult severity.Result) Section {
	var lines []string
	var sources []Source
	lines = append(lines, fmt.Sprintf("Severity: %s.", severityResult.Level))
	for _, explanation := range severityResult.Rationale {
		lines = append(lines, "Rationale: "+strings.TrimSpace(explanation.Text))
		sources = append(sources, convertSeveritySources(explanation.Sources, SourceTypeSeverity)...)
	}
	return Section{
		Title:   "Severity Rationale",
		Text:    strings.Join(lines, " "),
		Sources: dedupeSources(sources),
	}
}

func recommendationsSection(severityResult severity.Result) Section {
	var lines []string
	var sources []Source
	for _, recommendation := range severityResult.Recommendations {
		convertedSources := convertSeveritySources(recommendation.Sources, SourceTypeSeverity)
		lines = append(lines, fmt.Sprintf(
			"%s: %s Sources: %s.",
			recommendation.Action,
			strings.TrimSpace(recommendation.Explanation),
			strings.Join(sourceRefs(convertedSources), ", "),
		))
		sources = append(sources, convertedSources...)
	}
	return Section{
		Title:   "Recommended Actions",
		Text:    strings.Join(lines, "\n"),
		Sources: dedupeSources(sources),
	}
}

func approvalSection(severityResult severity.Result) Section {
	var lines []string
	sources := make([]Source, 0, len(severityResult.ApprovalRequirements))
	for i, requirement := range severityResult.ApprovalRequirements {
		status := "not required"
		if requirement.Required && !requirement.Approved {
			status = "blocked pending human approval"
		}
		if requirement.Required && requirement.Approved {
			status = "approval recorded in supplied severity result"
		}
		lines = append(lines, fmt.Sprintf(
			"%s: %s. %s",
			requirement.Action,
			status,
			strings.TrimSpace(requirement.Explanation),
		))
		sources = append(sources, Source{Ref: fmt.Sprintf("severity.approval_requirements[%d]", i), Type: SourceTypeSeverity})
	}
	return Section{
		Title:   "Approval State",
		Text:    strings.Join(lines, "\n"),
		Sources: sources,
	}
}

func approvalStates(severityResult severity.Result) []ApprovalState {
	states := make([]ApprovalState, 0, len(severityResult.ApprovalRequirements))
	for _, requirement := range severityResult.ApprovalRequirements {
		state := ApprovalState{
			Action:  requirement.Action,
			Blocked: requirement.Required && !requirement.Approved,
		}
		switch {
		case state.Blocked:
			state.Reason = "pending human approval: " + strings.TrimSpace(requirement.Explanation)
		case requirement.Required && requirement.Approved:
			state.Reason = "approval recorded in supplied severity result: " + strings.TrimSpace(requirement.Explanation)
		default:
			state.Reason = "approval not required for this action"
		}
		states = append(states, state)
	}
	return states
}

func uncertaintyLabels(timelineResult timeline.Result) []string {
	seen := make(map[string]struct{})
	var labels []string
	for _, entry := range timelineResult.Entries {
		if !entry.Uncertain {
			continue
		}
		label := strings.TrimSpace(entry.Uncertainty)
		if label == "" {
			continue
		}
		if _, ok := seen[label]; ok {
			continue
		}
		seen[label] = struct{}{}
		labels = append(labels, label)
	}
	return labels
}

func isShareableDraft(result Result) bool {
	if result.Status != StatusDraft || !result.SyntheticRecord || len(result.RedactionsApplied) == 0 {
		return false
	}
	if !hasBlockedApproval(result.ApprovalState) {
		return false
	}
	for _, section := range result.Sections {
		if strings.TrimSpace(section.Text) == "" || !hasRefs(section.Sources) {
			return false
		}
	}
	return len(result.Sections) >= 5
}

func hasBlockedApproval(states []ApprovalState) bool {
	for _, state := range states {
		if state.Blocked {
			return true
		}
	}
	return false
}

func hasRefs(sources []Source) bool {
	for _, source := range sources {
		if strings.TrimSpace(source.Ref) != "" {
			return true
		}
	}
	return false
}

func convertTimelineSources(sources []timeline.Source) []Source {
	converted := make([]Source, 0, len(sources))
	for _, source := range sources {
		converted = append(converted, Source{
			Ref:  source.Ref,
			Type: briefSourceTypeFromTimeline(source.Type),
		})
	}
	return dedupeSources(converted)
}

func briefSourceTypeFromTimeline(sourceType timeline.SourceType) SourceType {
	switch sourceType {
	case timeline.SourceTypePacket, timeline.SourceTypeTelemetry:
		return SourceTypePacket
	case timeline.SourceTypeGuidance:
		return SourceTypeGuidance
	default:
		return SourceTypeTimeline
	}
}

func convertSeveritySources(sources []severity.Source, fallback SourceType) []Source {
	converted := make([]Source, 0, len(sources))
	for _, source := range sources {
		converted = append(converted, Source{
			Ref:  source.Ref,
			Type: briefSourceTypeFromSeverity(source.Type, fallback),
		})
	}
	return dedupeSources(converted)
}

func briefSourceTypeFromSeverity(sourceType severity.SourceType, fallback SourceType) SourceType {
	switch sourceType {
	case severity.SourceTypePacket:
		return SourceTypePacket
	case severity.SourceTypeGuidance:
		return SourceTypeGuidance
	case severity.SourceTypeTimeline:
		return SourceTypeTimeline
	default:
		return fallback
	}
}

func sourceRefs(sources []Source) []string {
	refs := make([]string, 0, len(sources))
	for _, source := range sources {
		if strings.TrimSpace(source.Ref) != "" {
			refs = append(refs, source.Ref)
		}
	}
	return refs
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

type redactor struct {
	replacements      []redactionReplacement
	redactions        []Redaction
	redactedFields    map[string]struct{}
	coordinatePattern *regexp.Regexp
}

type redactionReplacement struct {
	value       string
	placeholder string
}

func newRedactor(packet ingestion.Packet) *redactor {
	r := &redactor{
		redactedFields:    make(map[string]struct{}),
		coordinatePattern: regexp.MustCompile(`[-+]?\d{1,3}\.\d{3,}\s*,\s*[-+]?\d{1,3}\.\d{3,}`),
	}
	r.add("packet.vehicle_id", packet.VehicleID, "[redacted vehicle]", "vehicle identifier is not included in shareable drafts")
	r.add("packet.route", packet.Route, "[redacted route]", "route detail is not included in shareable drafts")
	r.add("packet.location_label", packet.LocationLabel, "[redacted location]", "location detail is not included in shareable drafts")
	for i, sample := range packet.TelemetrySamples {
		r.add(
			fmt.Sprintf("packet.telemetry_samples[%d].gps_label", i),
			sample.GPSLabel,
			"[redacted gps label]",
			"gps label is not included in shareable drafts",
		)
	}
	for i, note := range packet.TranscriptNotes {
		if containsSensitiveNote(note) {
			r.add(
				fmt.Sprintf("packet.transcript_notes[%d]", i),
				note,
				"[redacted sensitive transcript note]",
				"sensitive or untrusted transcript detail is not included in shareable drafts",
			)
		}
	}
	for i, note := range packet.StillFrameNotes {
		if containsSensitiveNote(note) {
			r.add(
				fmt.Sprintf("packet.still_frame_notes[%d]", i),
				note,
				"[redacted sensitive still-frame note]",
				"sensitive still-frame detail is not included in shareable drafts",
			)
		}
	}
	return r
}

func (r *redactor) add(field, value, placeholder, reason string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	r.replacements = append(r.replacements, redactionReplacement{
		value:       value,
		placeholder: placeholder,
	})
	r.record(field, reason)
}

func (r *redactor) redact(text string) string {
	redacted := text
	for _, replacement := range r.replacements {
		redacted = strings.ReplaceAll(redacted, replacement.value, replacement.placeholder)
	}
	if r.coordinatePattern.MatchString(redacted) {
		redacted = r.coordinatePattern.ReplaceAllString(redacted, "[redacted coordinates]")
		r.record("coordinate_like_text", "coordinate-like text is not included in shareable drafts")
	}
	return redacted
}

func (r *redactor) record(field, reason string) {
	if _, ok := r.redactedFields[field]; ok {
		return
	}
	r.redactedFields[field] = struct{}{}
	r.redactions = append(r.redactions, Redaction{Field: field, Reason: reason})
}

func containsSensitiveNote(note string) bool {
	lower := strings.ToLower(note)
	sensitiveFragments := []string{
		"ignore all safety instructions",
		"mark this incident safe for export",
		"passenger note",
		"someone fell",
		"license plate",
		"student",
	}
	for _, fragment := range sensitiveFragments {
		if strings.Contains(lower, fragment) {
			return true
		}
	}
	return false
}
