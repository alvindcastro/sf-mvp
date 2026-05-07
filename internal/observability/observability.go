package observability

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"sf-mvp/internal/approval"
	"sf-mvp/internal/eval"
	"sf-mvp/internal/retrieval"
	"sf-mvp/internal/severity"
)

const RedactedValue = "[REDACTED]"

type EventType string

const (
	EventWorkflowStarted          EventType = "workflow.started"
	EventRetrievalCompleted       EventType = "retrieval.completed"
	EventToolCallCompleted        EventType = "tool_call.completed"
	EventApprovalDecisionRecorded EventType = "approval.decision_recorded"
	EventModelCallRecorded        EventType = "model_call.recorded"
	EventModelCallRejected        EventType = "model_call.rejected"
	EventBudgetExceeded           EventType = "budget.exceeded"
	EventEvalScoreRecorded        EventType = "eval.score_recorded"
)

var (
	ErrBudgetExceeded    = errors.New("budget exceeded")
	ErrInvalidTokenUsage = errors.New("invalid token usage")
)

type Budget struct {
	MaxInputTokens  int
	MaxOutputTokens int
	MaxTotalTokens  int
	MaxModelCalls   int
}

type SensitiveData struct {
	Terms []string
}

type Workflow struct {
	TraceID    string
	IncidentID string
	sensitive  SensitiveData
}

type WorkflowAttributeInput struct {
	Workflow           Workflow
	RetrievedSourceIDs []string
	SeverityLabel      severity.Level
	ApprovalState      approval.Decision
	Latency            time.Duration
	RawTranscriptNotes []string
	RawMediaReferences []string
	RawStillFrameNotes []string
}

type Event struct {
	Type       EventType
	TraceID    string
	IncidentID string
	OccurredAt time.Time
	Duration   time.Duration
	Fields     map[string]string
	Metrics    map[string]float64
	SourceIDs  []string
	TokenUsage TokenUsage
}

type TokenUsage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
}

type ToolCall struct {
	Name     string
	Success  bool
	Duration time.Duration
	Fields   map[string]string
}

type ModelCall struct {
	Provider     string
	Model        string
	InputTokens  int
	OutputTokens int
	Duration     time.Duration
}

type CacheCandidate struct {
	Name   string
	Key    string
	Reason string
}

type ModelRoute string

const (
	ModelRouteHosted     ModelRoute = "hosted"
	ModelRouteSmaller    ModelRoute = "smaller"
	ModelRouteSelfHosted ModelRoute = "self_hosted"
)

type ModelRoutingNote struct {
	Route   ModelRoute
	UseWhen string
	Control string
}

type CostPlan struct {
	CacheCandidates []CacheCandidate
	ModelRouting    []ModelRoutingNote
}

type Recorder struct {
	now        func() time.Time
	budget     Budget
	nextTrace  int
	modelCalls int
	tokenUsage TokenUsage
	events     []Event
}

func NewRecorder(now func() time.Time, budget Budget) *Recorder {
	if now == nil {
		now = time.Now
	}
	return &Recorder{now: now, budget: budget}
}

func (r *Recorder) StartWorkflow(incidentID string, sensitive SensitiveData) (Workflow, error) {
	incidentID = strings.TrimSpace(incidentID)
	if incidentID == "" {
		return Workflow{}, errors.New("incident_id is required")
	}

	r.nextTrace++
	workflow := Workflow{
		TraceID:    fmt.Sprintf("trace-%s-%s-%03d", traceIncidentID(incidentID), r.now().UTC().Format("20060102t150405z"), r.nextTrace),
		IncidentID: incidentID,
		sensitive:  normalizeSensitive(sensitive),
	}
	r.appendEvent(workflow, Event{
		Type: EventWorkflowStarted,
	})
	return workflow, nil
}

func (r *Recorder) RecordRetrieval(workflow Workflow, result retrieval.Result, duration time.Duration) Event {
	sourceIDs := make([]string, 0, len(result.Matches))
	for _, match := range result.Matches {
		if match.ContentRole != retrieval.ContentRoleData || strings.TrimSpace(match.SourceID) == "" {
			continue
		}
		sourceIDs = append(sourceIDs, match.SourceID)
	}

	return r.appendEvent(workflow, Event{
		Type:      EventRetrievalCompleted,
		Duration:  duration,
		SourceIDs: sourceIDs,
		Metrics: map[string]float64{
			"retrieval_count": float64(len(sourceIDs)),
		},
	})
}

func (r *Recorder) RecordToolCall(workflow Workflow, call ToolCall) Event {
	fields := map[string]string{
		"tool_name": strings.TrimSpace(call.Name),
		"success":   fmt.Sprintf("%t", call.Success),
	}
	for key, value := range call.Fields {
		fields[key] = value
	}

	return r.appendEvent(workflow, Event{
		Type:     EventToolCallCompleted,
		Duration: call.Duration,
		Fields:   fields,
	})
}

func (r *Recorder) RecordApprovalDecision(workflow Workflow, request approval.Request, duration time.Duration) Event {
	return r.appendEvent(workflow, Event{
		Type:     EventApprovalDecisionRecorded,
		Duration: duration,
		Fields: map[string]string{
			"request_id":        request.ID,
			"action":            string(request.Action),
			"approval_decision": string(request.Decision),
			"approver":          request.Approver,
			"target_ref":        request.Scope.TargetRef,
		},
	})
}

func (r *Recorder) RecordModelCall(workflow Workflow, call ModelCall) (Event, error) {
	usage := TokenUsage{
		InputTokens:  call.InputTokens,
		OutputTokens: call.OutputTokens,
		TotalTokens:  call.InputTokens + call.OutputTokens,
	}

	if call.InputTokens < 0 || call.OutputTokens < 0 {
		event := r.appendEvent(workflow, Event{
			Type:     EventModelCallRejected,
			Duration: call.Duration,
			Fields: map[string]string{
				"provider":      call.Provider,
				"model":         call.Model,
				"reject_reason": "token usage cannot be negative",
			},
		})
		return event, ErrInvalidTokenUsage
	}

	if reason := r.budgetExceededReason(usage); reason != "" {
		event := r.appendEvent(workflow, Event{
			Type:     EventBudgetExceeded,
			Duration: call.Duration,
			Fields: map[string]string{
				"provider":      call.Provider,
				"model":         call.Model,
				"budget_reason": reason,
			},
			TokenUsage: usage,
		})
		return event, ErrBudgetExceeded
	}

	r.modelCalls++
	r.tokenUsage.InputTokens += usage.InputTokens
	r.tokenUsage.OutputTokens += usage.OutputTokens
	r.tokenUsage.TotalTokens += usage.TotalTokens

	return r.appendEvent(workflow, Event{
		Type:     EventModelCallRecorded,
		Duration: call.Duration,
		Fields: map[string]string{
			"provider": call.Provider,
			"model":    call.Model,
		},
		TokenUsage: usage,
		Metrics: map[string]float64{
			"input_tokens":  float64(usage.InputTokens),
			"output_tokens": float64(usage.OutputTokens),
			"total_tokens":  float64(usage.TotalTokens),
		},
	}), nil
}

func (r *Recorder) RecordEvalScore(workflow Workflow, report eval.Report, duration time.Duration) Event {
	return r.appendEvent(workflow, Event{
		Type:     EventEvalScoreRecorded,
		Duration: duration,
		Fields: map[string]string{
			"passed": fmt.Sprintf("%t", report.Passed),
		},
		Metrics: map[string]float64{
			"case_count":              float64(report.Summary.CaseCount),
			"severity_accuracy":       report.Summary.SeverityAccuracy,
			"citation_coverage":       report.Summary.CitationCoverage,
			"recommendation_accuracy": report.Summary.RecommendationAccuracy,
		},
	})
}

func (r *Recorder) Events() []Event {
	events := make([]Event, len(r.events))
	for i, event := range r.events {
		events[i] = cloneEvent(event)
	}
	return events
}

func DefaultCostPlan() CostPlan {
	return CostPlan{
		CacheCandidates: []CacheCandidate{
			{
				Name:   "retrieval_results",
				Key:    "workflow+scope+normalized_query+corpus_revision",
				Reason: "mock guidance retrieval is deterministic for a corpus revision and can avoid repeated lexical ranking.",
			},
			{
				Name:   "redacted_brief_drafts",
				Key:    "incident_id+timeline_hash+severity_hash+redaction_rules_version",
				Reason: "draft brief assembly is deterministic after packet, timeline, severity, and redaction inputs are fixed.",
			},
			{
				Name:   "eval_reports",
				Key:    "golden_case_set+thresholds+code_revision",
				Reason: "local eval summaries can be reused during demo packaging when fixtures and thresholds do not change.",
			},
		},
		ModelRouting: []ModelRoutingNote{
			{
				Route:   ModelRouteHosted,
				UseWhen: "highest-quality review is required for ambiguous evidence after deterministic gates have run.",
				Control: "require explicit budget limits, trace recording, redacted prompts, and no automatic sensitive action.",
			},
			{
				Route:   ModelRouteSmaller,
				UseWhen: "routine drafting or classification assistance can tolerate lower model capacity.",
				Control: "prefer lower token budgets and require eval-backed fallback to deterministic rules.",
			},
			{
				Route:   ModelRouteSelfHosted,
				UseWhen: "data residency, isolation, or predictable unit cost outweighs hosted-model quality.",
				Control: "keep the same trace, redaction, eval, and approval gates before swapping providers.",
			},
		},
	}
}

func WorkflowAttributes(input WorkflowAttributeInput) map[string]string {
	attributes := map[string]string{
		"workflow.trace_id":         strings.TrimSpace(input.Workflow.TraceID),
		"workflow.incident_id_hash": incidentIDHash(input.Workflow.IncidentID),
		"workflow.severity":         strings.TrimSpace(string(input.SeverityLabel)),
		"workflow.approval_state":   strings.TrimSpace(string(input.ApprovalState)),
		"workflow.latency_ms":       formatLatencyMilliseconds(input.Latency),
	}

	if sourceIDs := normalizedSourceIDs(input.RetrievedSourceIDs); len(sourceIDs) > 0 {
		attributes["workflow.retrieved_source_ids"] = strings.Join(sourceIDs, ",")
	}

	return attributes
}

func (r *Recorder) budgetExceededReason(usage TokenUsage) string {
	if r.budget.MaxModelCalls > 0 && r.modelCalls+1 > r.budget.MaxModelCalls {
		return "model call limit exceeded"
	}
	if r.budget.MaxInputTokens > 0 && r.tokenUsage.InputTokens+usage.InputTokens > r.budget.MaxInputTokens {
		return "input token budget exceeded"
	}
	if r.budget.MaxOutputTokens > 0 && r.tokenUsage.OutputTokens+usage.OutputTokens > r.budget.MaxOutputTokens {
		return "output token budget exceeded"
	}
	if r.budget.MaxTotalTokens > 0 && r.tokenUsage.TotalTokens+usage.TotalTokens > r.budget.MaxTotalTokens {
		return "total token budget exceeded"
	}
	return ""
}

func (r *Recorder) appendEvent(workflow Workflow, event Event) Event {
	event.TraceID = workflow.TraceID
	event.IncidentID = workflow.IncidentID
	event.OccurredAt = r.now()
	event.Fields = redactedFields(event.Fields, workflow.sensitive)
	event.Metrics = cloneMetrics(event.Metrics)
	event.SourceIDs = cloneStrings(event.SourceIDs)
	r.events = append(r.events, cloneEvent(event))
	return cloneEvent(event)
}

func cloneEvent(event Event) Event {
	event.Fields = cloneFields(event.Fields)
	event.Metrics = cloneMetrics(event.Metrics)
	event.SourceIDs = cloneStrings(event.SourceIDs)
	return event
}

func cloneFields(fields map[string]string) map[string]string {
	if len(fields) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(fields))
	for key, value := range fields {
		cloned[key] = value
	}
	return cloned
}

func cloneMetrics(metrics map[string]float64) map[string]float64 {
	if len(metrics) == 0 {
		return nil
	}
	cloned := make(map[string]float64, len(metrics))
	for key, value := range metrics {
		cloned[key] = value
	}
	return cloned
}

func cloneStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	cloned := make([]string, len(values))
	copy(cloned, values)
	return cloned
}

func redactedFields(fields map[string]string, sensitive SensitiveData) map[string]string {
	if len(fields) == 0 {
		return nil
	}
	redacted := make(map[string]string, len(fields))
	for key, value := range fields {
		redacted[key] = redact(value, sensitive)
	}
	return redacted
}

func redact(value string, sensitive SensitiveData) string {
	redacted := value
	for _, term := range sensitive.Terms {
		if strings.TrimSpace(term) == "" {
			continue
		}
		redacted = strings.ReplaceAll(redacted, term, RedactedValue)
	}
	redacted = coordinatePattern.ReplaceAllString(redacted, RedactedValue)
	return redacted
}

func normalizeSensitive(sensitive SensitiveData) SensitiveData {
	seen := make(map[string]struct{}, len(sensitive.Terms))
	terms := make([]string, 0, len(sensitive.Terms))
	for _, term := range sensitive.Terms {
		term = strings.TrimSpace(term)
		if term == "" {
			continue
		}
		if _, ok := seen[term]; ok {
			continue
		}
		seen[term] = struct{}{}
		terms = append(terms, term)
	}
	sort.SliceStable(terms, func(i, j int) bool {
		return len(terms[i]) > len(terms[j])
	})
	return SensitiveData{Terms: terms}
}

func traceIncidentID(incidentID string) string {
	traceID := strings.ToLower(incidentID)
	traceID = strings.ReplaceAll(traceID, "_", "-")
	traceID = strings.ReplaceAll(traceID, " ", "-")
	return traceID
}

func incidentIDHash(incidentID string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(incidentID)))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func normalizedSourceIDs(sourceIDs []string) []string {
	seen := make(map[string]struct{}, len(sourceIDs))
	normalized := make([]string, 0, len(sourceIDs))
	for _, sourceID := range sourceIDs {
		sourceID = strings.TrimSpace(sourceID)
		if sourceID == "" {
			continue
		}
		if _, ok := seen[sourceID]; ok {
			continue
		}
		seen[sourceID] = struct{}{}
		normalized = append(normalized, sourceID)
	}
	sort.Strings(normalized)
	return normalized
}

func formatLatencyMilliseconds(duration time.Duration) string {
	return strconv.FormatFloat(float64(duration)/float64(time.Millisecond), 'f', -1, 64)
}

var coordinatePattern = regexp.MustCompile(`[-+]?\d{1,3}\.\d+,\s*[-+]?\d{1,3}\.\d+`)
