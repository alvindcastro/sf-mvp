# Observability And Cost Controls

Phase 9 adds the first package-level observability and budget-control surface for Fleet Incident Copilot. The surface records structured in-memory workflow events for synthetic MVP workflows and defines cost-control planning notes. It does not add an external telemetry backend, dashboards, alerts, provider billing reconciliation, live model-provider calls, persistent cache storage, CLI, HTTP API, or database behavior.

## Phase 9 Checklist

- [x] Generate a trace ID per incident workflow.
- [x] Track retrieval count, retrieved source IDs, tool-call success, approval decisions, latency, token usage, and eval score.
- [x] Redact sensitive data from logs.
- [x] Add model-call budget limits.
- [x] Define caching candidates.
- [x] Define model-routing notes for hosted, smaller, or self-hosted models.

## Runtime Surface

- Package: [internal/observability](../../../internal/observability).
- Recorder constructor: `NewRecorder(now func() time.Time, budget Budget) *Recorder`.
- Workflow entry point: `StartWorkflow(incidentID string, sensitive SensitiveData) (Workflow, error)`.
- Event history: `Events() []Event`.
- Budget error: `ErrBudgetExceeded`.
- Invalid token usage error: `ErrInvalidTokenUsage`.
- Cost-control plan: `DefaultCostPlan() CostPlan`.
- Targeted test command: `go test ./internal/observability`.
- Full test command: `go test ./...`.

The recorder is intentionally orchestration-facing. It observes existing Phase 2 through Phase 8 outputs without threading trace state through every package API.

## Structured Events

Every event includes:

- `Type`.
- `TraceID`.
- `IncidentID`.
- `OccurredAt`.
- Optional `Duration`.
- Optional string `Fields`.
- Optional numeric `Metrics`.
- Optional `SourceIDs`.
- Optional caller-supplied `TokenUsage`.

Supported event types:

- `workflow.started`.
- `retrieval.completed`.
- `tool_call.completed`.
- `approval.decision_recorded`.
- `model_call.recorded`.
- `model_call.rejected`.
- `budget.exceeded`.
- `eval.score_recorded`.

## Recorded Signals

- Trace propagation: `StartWorkflow` creates a deterministic trace ID from incident ID, timestamp, and recorder sequence.
- Retrieval: `RecordRetrieval` records retrieved `SourceID` values and retrieval count from `retrieval.Result` without logging citation snippets.
- Tool calls: `RecordToolCall` records tool name, success flag, duration, and redacted caller fields.
- Approval decisions: `RecordApprovalDecision` records approval request ID, action, decision, approver, target reference, and latency from `approval.Request`.
- Model calls: `RecordModelCall` records caller-supplied provider, model, input tokens, output tokens, total tokens, and latency when token usage is valid.
- Invalid model calls: `RecordModelCall` rejects negative caller-supplied token counts with `model_call.rejected` and `ErrInvalidTokenUsage` before budget accounting.
- Eval summaries: `RecordEvalScore` records case count, severity accuracy, citation coverage, recommendation accuracy, pass/fail state, and duration from `eval.Report`.

## Redaction Rules

`SensitiveData` accepts incident-specific terms that must not appear in structured event fields. The recorder redacts those terms and coordinate-like latitude/longitude strings before storing or returning events.

Current tests cover vehicle identifiers, route names, private location labels, coordinate-like strings, and caller-supplied tool fields. Retrieval snippets, packet transcript notes, still-frame notes, and raw evidence references should not be logged into fields unless a future strict-TDD change proves they are redacted.

## Budget Controls

`Budget` supports deterministic limits for:

- Input tokens.
- Output tokens.
- Total tokens.
- Model-call count.

`RecordModelCall` appends `model_call.recorded` when usage is non-negative and within budget. It appends `budget.exceeded` and returns `ErrBudgetExceeded` when a new call would exceed a configured limit. It appends `model_call.rejected` and returns `ErrInvalidTokenUsage` when caller-supplied input or output tokens are negative, and rejected calls do not consume model-call count or token budget.

Token usage is caller-supplied because the MVP does not call a real model provider. The behavior is a local control surface, not billing reconciliation or provider-side enforcement.

## Cache Candidates

`DefaultCostPlan()` defines the first cache candidates:

| Candidate | Key | Reason |
| --- | --- | --- |
| `retrieval_results` | `workflow+scope+normalized_query+corpus_revision` | Mock guidance retrieval is deterministic for a corpus revision and can avoid repeated lexical ranking. |
| `redacted_brief_drafts` | `incident_id+timeline_hash+severity_hash+redaction_rules_version` | Draft brief assembly is deterministic after packet, timeline, severity, and redaction inputs are fixed. |
| `eval_reports` | `golden_case_set+thresholds+code_revision` | Local eval summaries can be reused during demo packaging when fixtures and thresholds do not change. |

The package defines candidates only. It does not implement cache storage, invalidation, persistence, or eviction.

## Model Routing Notes

`DefaultCostPlan()` defines model-routing notes for:

- `hosted`: use when the highest-quality review is needed for ambiguous evidence after deterministic gates run; require budget limits, trace recording, redacted prompts, and no automatic sensitive action.
- `smaller`: use for routine drafting or classification assistance when lower capacity is acceptable; prefer lower token budgets and eval-backed fallback to deterministic rules.
- `self_hosted`: use when data residency, isolation, or predictable unit cost outweighs hosted-model quality; keep the same trace, redaction, eval, and approval gates before swapping providers.

The package records routing notes only. It does not select providers, call models, or run live routing.

## Red-To-Green Evidence

- Red: `go test ./internal/observability` initially failed because `NewRecorder`, `Budget`, `SensitiveData`, workflow, event, budget, cache, and model-routing types did not exist.
- Green: after implementing [internal/observability/observability.go](../../../internal/observability/observability.go), `go test ./internal/observability` passed.
- Regression red: `go test ./internal/observability` failed after adding the negative token usage test because `ErrInvalidTokenUsage` and `model_call.rejected` did not exist.
- Regression green: after adding the invalid token guard, `go test ./internal/observability` passed.

## Current Limits

- The recorder is an in-memory package API, not a logging framework or telemetry exporter.
- Events are not persisted and are not shipped to OpenTelemetry, a metrics backend, dashboards, or alerts.
- Budget controls use caller-supplied token usage; no real model-provider calls or provider billing reconciliation exists.
- Cache candidates and routing notes are documentation-backed structs, not live cache or model-routing behavior.
- No CLI, HTTP API, database, external observability pipeline, persistent audit store, real export tool, real escalation tool, or external-sharing integration exists.
- This phase does not provide production audit, compliance, retention, security, or cost-governance guarantees.
