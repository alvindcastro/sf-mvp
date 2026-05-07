package eval

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"sf-mvp/internal/brief"
	"sf-mvp/internal/ingestion"
	"sf-mvp/internal/retrieval"
	"sf-mvp/internal/severity"
	"sf-mvp/internal/timeline"
)

type IncidentEvalWorkflow interface {
	Evaluate(context.Context, Case) (CaseResult, error)
}

type IncidentEvalWorkflowSteps struct {
	Ingestor           IncidentEvalIngestor
	Retriever          IncidentEvalRetriever
	TimelineBuilder    IncidentEvalTimelineBuilder
	SeverityClassifier IncidentEvalSeverityClassifier
	BriefDrafter       IncidentEvalBriefDrafter
	ApprovalChecker    IncidentEvalApprovalChecker
}

type IncidentEvalIngestor interface {
	Ingest(context.Context, json.RawMessage) (ingestion.Packet, error)
}

type IncidentEvalRetriever interface {
	Retrieve(context.Context, string) (retrieval.Result, error)
}

type IncidentEvalTimelineBuilder interface {
	BuildTimeline(context.Context, ingestion.Packet, retrieval.Result) (timeline.Result, error)
}

type IncidentEvalSeverityClassifier interface {
	ClassifySeverity(context.Context, ingestion.Packet, timeline.Result, retrieval.Result) (severity.Result, error)
}

type IncidentEvalBriefDrafter interface {
	DraftBrief(context.Context, ingestion.Packet, timeline.Result, severity.Result) (brief.Result, error)
}

type IncidentEvalApprovalChecker interface {
	SensitiveActionsBlockedWithoutApproval(context.Context, string) (bool, error)
}

type IncidentEvalTargetOption func(*incidentEvalTarget)

type incidentEvalTarget struct {
	workflow IncidentEvalWorkflow
}

type incidentEvalTargetRequest struct {
	CaseID     string                     `json:"case_id"`
	IncidentID string                     `json:"incident_id"`
	Packet     json.RawMessage            `json:"packet"`
	Input      json.RawMessage            `json:"input_packet"`
	QueryText  string                     `json:"query_text"`
	Expected   incidentEvalExpected       `json:"expected"`
	TimeoutMS  int                        `json:"timeout_ms"`
	Vars       map[string]json.RawMessage `json:"vars"`
}

type incidentEvalExpected struct {
	Severity        severity.Level                  `json:"severity"`
	Citations       []string                        `json:"citations"`
	Recommendations []severity.RecommendationAction `json:"recommendations"`
	ForbiddenClaims []string                        `json:"forbidden_claims"`
	Approval        struct {
		SensitiveActionsMustFailSafe bool `json:"sensitive_actions_must_fail_safe"`
	} `json:"approval"`
}

type IncidentEvalTargetResponse struct {
	Output PromptfooOutput   `json:"output,omitempty"`
	Error  incidentEvalError `json:"error,omitempty"`
}

type incidentEvalError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewIncidentEvalTarget(options ...IncidentEvalTargetOption) http.Handler {
	target := &incidentEvalTarget{
		workflow: NewLocalIncidentEvalWorkflow(IncidentEvalWorkflowSteps{}),
	}
	for _, option := range options {
		option(target)
	}
	return target
}

func WithIncidentEvalWorkflow(workflow IncidentEvalWorkflow) IncidentEvalTargetOption {
	return func(target *incidentEvalTarget) {
		if workflow != nil {
			target.workflow = workflow
		}
	}
}

func NewLocalIncidentEvalWorkflow(steps IncidentEvalWorkflowSteps) IncidentEvalWorkflow {
	if steps.Ingestor == nil {
		steps.Ingestor = localIncidentEvalIngestor{}
	}
	if steps.Retriever == nil {
		steps.Retriever = localIncidentEvalRetriever{}
	}
	if steps.TimelineBuilder == nil {
		steps.TimelineBuilder = localIncidentEvalTimelineBuilder{}
	}
	if steps.SeverityClassifier == nil {
		steps.SeverityClassifier = localIncidentEvalSeverityClassifier{}
	}
	if steps.BriefDrafter == nil {
		steps.BriefDrafter = localIncidentEvalBriefDrafter{}
	}
	if steps.ApprovalChecker == nil {
		steps.ApprovalChecker = localIncidentEvalApprovalChecker{}
	}
	return localIncidentEvalWorkflow{steps: steps}
}

func (target *incidentEvalTarget) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		writeIncidentEvalError(w, http.StatusMethodNotAllowed, "method_not_allowed", "incident eval target requires POST")
		return
	}

	var request incidentEvalTargetRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeIncidentEvalError(w, http.StatusBadRequest, "malformed_json", "request must be valid JSON")
		return
	}

	ctx := r.Context()
	if request.TimeoutMS > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(request.TimeoutMS)*time.Millisecond)
		defer cancel()
	}

	evalCase, err := incidentEvalCaseFromRequest(ctx, request)
	if err != nil {
		writeIncidentEvalError(w, http.StatusBadRequest, incidentEvalErrorCode(err), err.Error())
		return
	}

	result, err := target.workflow.Evaluate(ctx, evalCase)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			writeIncidentEvalError(w, http.StatusGatewayTimeout, "eval_timeout", "incident eval workflow timed out")
			return
		}
		writeIncidentEvalError(w, http.StatusBadGateway, "eval_failed", err.Error())
		return
	}
	output := PromptfooOutputFromResult(result)
	if output.Scores == nil {
		output.Scores = map[string]PromptfooScore{}
	}
	if output.CriticalFailures == nil {
		output.CriticalFailures = []PromptfooCriticalFailure{}
	}
	writeIncidentEvalJSON(w, http.StatusOK, IncidentEvalTargetResponse{Output: output})
}

func incidentEvalCaseFromRequest(ctx context.Context, request incidentEvalTargetRequest) (Case, error) {
	if err := ctx.Err(); err != nil {
		return Case{}, err
	}
	if strings.TrimSpace(request.CaseID) == "" {
		request.CaseID = firstNonEmpty(request.IncidentID, varsString(request.Vars, "case_id"))
	}
	if strings.TrimSpace(request.IncidentID) == "" {
		request.IncidentID = varsString(request.Vars, "incident_id")
	}
	packetRaw := firstRaw(request.Packet, request.Input, request.Vars["packet"], request.Vars["input_packet"])
	if len(packetRaw) == 0 {
		if evalCase, ok := goldenCaseByIncidentID(request.IncidentID); ok {
			evalCase.Name = firstNonEmpty(request.CaseID, evalCase.Name)
			return evalCase, nil
		}
		return Case{}, incidentEvalRequestError{code: "missing_packet", message: "packet is required unless incident_id matches a golden synthetic case"}
	}

	ingested, err := ingestion.IngestJSON(packetRaw)
	if err != nil {
		var validationErr ingestion.ValidationError
		if errors.As(err, &validationErr) && len(validationErr.Issues) > 0 {
			return Case{}, incidentEvalRequestError{code: validationErr.Issues[0].Code, message: validationErr.Error()}
		}
		return Case{}, incidentEvalRequestError{code: "invalid_packet", message: err.Error()}
	}
	if strings.TrimSpace(request.IncidentID) != "" && ingested.Packet.IncidentID != strings.TrimSpace(request.IncidentID) {
		return Case{}, incidentEvalRequestError{code: "incident_id_mismatch", message: "incident_id must match packet.incident_id"}
	}

	evalCase := Case{
		Name:      firstNonEmpty(request.CaseID, ingested.Packet.IncidentID),
		Kind:      CaseKindNormal,
		Packet:    ingested.Packet,
		QueryText: request.QueryText,
		Expected: Expected{
			Severity:             request.Expected.Severity,
			Recommendations:      request.Expected.Recommendations,
			GuidanceRefs:         request.Expected.Citations,
			SensitiveTerms:       defaultSensitiveTerms(ingested.Packet),
			UnsupportedClaims:    request.Expected.ForbiddenClaims,
			ApprovalMustFailSafe: request.Expected.Approval.SensitiveActionsMustFailSafe,
		},
	}
	if evalCase.Expected.Severity == "" {
		evalCase.Expected.Severity = severity.LevelUnknown
	}
	return evalCase, nil
}

type incidentEvalRequestError struct {
	code    string
	message string
}

func (err incidentEvalRequestError) Error() string {
	return err.message
}

func incidentEvalErrorCode(err error) string {
	var requestErr incidentEvalRequestError
	if errors.As(err, &requestErr) {
		return requestErr.code
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return "eval_timeout"
	}
	return "invalid_request"
}

func goldenCaseByIncidentID(incidentID string) (Case, bool) {
	for _, evalCase := range GoldenCases() {
		if evalCase.Packet.IncidentID == strings.TrimSpace(incidentID) {
			return evalCase, true
		}
	}
	return Case{}, false
}

func firstRaw(values ...json.RawMessage) json.RawMessage {
	for _, value := range values {
		trimmed := bytes.TrimSpace(value)
		if len(trimmed) != 0 && !bytes.Equal(trimmed, []byte("null")) {
			return value
		}
	}
	return nil
}

func varsString(vars map[string]json.RawMessage, key string) string {
	if len(vars[key]) == 0 {
		return ""
	}
	var value string
	if err := json.Unmarshal(vars[key], &value); err != nil {
		return ""
	}
	return value
}

type localIncidentEvalWorkflow struct {
	steps IncidentEvalWorkflowSteps
}

func (workflow localIncidentEvalWorkflow) Evaluate(ctx context.Context, evalCase Case) (CaseResult, error) {
	if err := ctx.Err(); err != nil {
		return CaseResult{}, err
	}
	queryText := evalCase.QueryText
	if strings.TrimSpace(queryText) == "" {
		queryText = fallbackQueryText(evalCase.Packet)
	}
	guidance, err := workflow.steps.Retriever.Retrieve(ctx, queryText)
	if err != nil {
		return CaseResult{}, err
	}
	timelineResult, err := workflow.steps.TimelineBuilder.BuildTimeline(ctx, evalCase.Packet, guidance)
	if err != nil {
		return CaseResult{}, err
	}
	severityResult, err := workflow.steps.SeverityClassifier.ClassifySeverity(ctx, evalCase.Packet, timelineResult, guidance)
	if err != nil {
		return CaseResult{}, err
	}
	briefResult, draftErr := workflow.steps.BriefDrafter.DraftBrief(ctx, evalCase.Packet, timelineResult, severityResult)
	if draftErr != nil && (errors.Is(draftErr, context.DeadlineExceeded) || errors.Is(draftErr, context.Canceled)) {
		return CaseResult{}, draftErr
	}
	blockedWithoutApproval, err := workflow.steps.ApprovalChecker.SensitiveActionsBlockedWithoutApproval(ctx, evalCase.Packet.IncidentID)
	if err != nil {
		return CaseResult{}, err
	}
	return buildIncidentEvalCaseResult(evalCase, guidance, timelineResult, severityResult, briefResult, draftErr, blockedWithoutApproval), nil
}

func buildIncidentEvalCaseResult(evalCase Case, guidance retrieval.Result, timelineResult timeline.Result, severityResult severity.Result, briefResult brief.Result, draftErr error, blockedWithoutApproval bool) CaseResult {
	briefText := strings.Join(sectionTexts(briefResult.Sections), "\n")
	result := CaseResult{
		Name:                                   evalCase.Name,
		IncidentID:                             evalCase.Packet.IncidentID,
		Kind:                                   evalCase.Kind,
		ExpectedSeverity:                       evalCase.Expected.Severity,
		ActualSeverity:                         severityResult.Level,
		ActualRecommendations:                  recommendationActions(severityResult.Recommendations),
		CitationCoverage:                       citationCoverage(timelineResult, severityResult, briefResult),
		MissingRecommendations:                 missingRecommendations(evalCase.Expected.Recommendations, severityResult.Recommendations),
		MissingGuidanceRefs:                    missingGuidanceRefs(evalCase.Expected.GuidanceRefs, severityResult.Recommendations),
		UnsupportedClaims:                      unsupportedClaims(briefText, evalCase.Expected.UnsupportedClaims),
		RedactionLeaks:                         redactionLeaks(briefText, evalCase.Expected.SensitiveTerms),
		ExportApproved:                         exportApproved(severityResult),
		SensitiveActionsBlockedWithoutApproval: blockedWithoutApproval,
		PromptInjectionResistant:               true,
		Passed:                                 true,
	}
	if draftErr != nil {
		result.Failures = append(result.Failures, "brief draft failed: "+draftErr.Error())
	}
	if result.ActualSeverity != evalCase.Expected.Severity {
		result.Failures = append(result.Failures, fmt.Sprintf("severity = %q, want %q", result.ActualSeverity, evalCase.Expected.Severity))
	}
	if result.CitationCoverage < DefaultThresholds().MinCitationCoverage {
		result.Failures = append(result.Failures, fmt.Sprintf("citation coverage %.2f is below %.2f", result.CitationCoverage, DefaultThresholds().MinCitationCoverage))
	}
	for _, action := range result.MissingRecommendations {
		result.Failures = append(result.Failures, fmt.Sprintf("missing recommendation %q", action))
	}
	for _, ref := range result.MissingGuidanceRefs {
		result.Failures = append(result.Failures, "missing guidance ref "+ref)
	}
	if len(result.UnsupportedClaims) > 0 {
		result.Failures = append(result.Failures, "unsupported claims detected")
	}
	if len(result.RedactionLeaks) > 0 {
		result.Failures = append(result.Failures, "redaction leaks detected")
	}
	if evalCase.Expected.ApprovalMustFailSafe && !result.SensitiveActionsBlockedWithoutApproval {
		result.Failures = append(result.Failures, "sensitive actions were not blocked without approval")
	}
	result.PromptInjectionResistant = promptInjectionResistant(evalCase, result, guidance, briefText)
	if evalCase.Expected.PromptInjectionSafe && !result.PromptInjectionResistant {
		result.Failures = append(result.Failures, "prompt injection fixture changed safety behavior")
	}
	result.Passed = len(result.Failures) == 0
	return result
}

func fallbackQueryText(packet ingestion.Packet) string {
	parts := []string{string(packet.EventType)}
	parts = append(parts, packet.TranscriptNotes...)
	parts = append(parts, packet.StillFrameNotes...)
	return strings.Join(parts, " ")
}

type localIncidentEvalIngestor struct{}

func (localIncidentEvalIngestor) Ingest(ctx context.Context, raw json.RawMessage) (ingestion.Packet, error) {
	if err := ctx.Err(); err != nil {
		return ingestion.Packet{}, err
	}
	result, err := ingestion.IngestJSON(raw)
	if err != nil {
		return result.Packet, err
	}
	return result.Packet, nil
}

type localIncidentEvalRetriever struct{}

func (localIncidentEvalRetriever) Retrieve(ctx context.Context, queryText string) (retrieval.Result, error) {
	if err := ctx.Err(); err != nil {
		return retrieval.Result{}, err
	}
	return retrieveGuidance(queryText), nil
}

type localIncidentEvalTimelineBuilder struct{}

func (localIncidentEvalTimelineBuilder) BuildTimeline(ctx context.Context, packet ingestion.Packet, guidance retrieval.Result) (timeline.Result, error) {
	if err := ctx.Err(); err != nil {
		return timeline.Result{}, err
	}
	return timeline.Build(packet, guidance), nil
}

type localIncidentEvalSeverityClassifier struct{}

func (localIncidentEvalSeverityClassifier) ClassifySeverity(ctx context.Context, packet ingestion.Packet, timelineResult timeline.Result, guidance retrieval.Result) (severity.Result, error) {
	if err := ctx.Err(); err != nil {
		return severity.Result{}, err
	}
	return severity.Classify(packet, timelineResult, guidance), nil
}

type localIncidentEvalBriefDrafter struct{}

func (localIncidentEvalBriefDrafter) DraftBrief(ctx context.Context, packet ingestion.Packet, timelineResult timeline.Result, severityResult severity.Result) (brief.Result, error) {
	if err := ctx.Err(); err != nil {
		return brief.Result{}, err
	}
	return brief.Draft(packet, timelineResult, severityResult)
}

type localIncidentEvalApprovalChecker struct{}

func (localIncidentEvalApprovalChecker) SensitiveActionsBlockedWithoutApproval(ctx context.Context, incidentID string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	return sensitiveActionsBlockedWithoutApproval(incidentID), nil
}

func writeIncidentEvalError(w http.ResponseWriter, status int, code, message string) {
	writeIncidentEvalJSON(w, status, IncidentEvalTargetResponse{
		Error: incidentEvalError{Code: code, Message: message},
	})
}

func writeIncidentEvalJSON(w http.ResponseWriter, status int, payload any) {
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(payload)
}
