# Changelog

## 2026-05-06 - Phase 9 Observability And Cost Controls

### Task: Add Phase 9 Observability TDD Coverage

- What: Added failing-first observability tests for trace ID creation, structured workflow events, retrieval counts and source IDs, tool-call success and redaction, approval decision metadata, token and latency recording, budget-exceeded behavior, eval score recording, cache candidates, and model-routing notes.
- Where: [internal/observability/observability_test.go](internal/observability/observability_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 9 with strict TDD and make quality, safety, latency, and spend signals observable before adding production recorder behavior.
- How: Created same-package tests around `NewRecorder`, `StartWorkflow`, retrieval, tool-call, approval, model-call, eval, and cost-plan APIs; ran `go test ./internal/observability`; observed the expected red build failure for missing observability APIs and types before adding implementation.

### Task: Implement In-Memory Observability Recorder And Cost Plan

- What: Added an observability package that creates deterministic trace IDs, records in-memory structured events for workflow start, retrieval, tool calls, approval decisions, model calls, budget failures, and eval summaries, redacts configured sensitive terms and coordinate-like strings, and exposes cache candidates plus hosted, smaller-model, and self-hosted routing notes.
- Where: [internal/observability/observability.go](internal/observability/observability.go).
- When: 2026-05-06, America/Vancouver.
- Why: Provide the first package-level Phase 9 observability and cost-control surface without adding an external telemetry backend, dashboards, alerts, persistent log store, provider billing reconciliation, live model-provider calls, cache storage, CLI, HTTP API, database behavior, real export, real escalation, or external-sharing integrations.
- How: Implemented `NewRecorder`, `StartWorkflow`, `Events`, `RecordRetrieval`, `RecordToolCall`, `RecordApprovalDecision`, `RecordModelCall`, `RecordEvalScore`, `Budget`, `SensitiveData`, event models, defensive event copying, redaction helpers, cumulative token accounting, budget-limit checks, and `DefaultCostPlan`; then confirmed `go test ./internal/observability` passed.

### Task: Tighten Phase 9 Regression Coverage

- What: Added regression assertions for retrieval snippet non-logging, approval request ID, approver, and latency recording, private location label redaction, input, output, and total token budget failures, and invalid negative token usage rejection.
- Where: [internal/observability/observability_test.go](internal/observability/observability_test.go), [internal/observability/observability.go](internal/observability/observability.go).
- When: 2026-05-06, America/Vancouver.
- Why: Close Phase 9 assertion gaps and prevent invalid caller-supplied token counts from reducing cumulative spend accounting.
- How: Added a failing test for negative token usage, observed `go test ./internal/observability` fail for missing `ErrInvalidTokenUsage` and `model_call.rejected`, implemented the smallest rejection path before budget accounting, and expanded existing tests for the documented metadata, redaction, and budget dimensions.

### Task: Document Phase 9 Observability Contract

- What: Added the Phase 9 behavior contract, runtime surface, structured event model, recorded signals, redaction rules, budget controls, cache candidates, model-routing notes, current limits, test commands, and red-to-green evidence.
- Where: [docs/mvp/observability-and-cost-controls.md](docs/mvp/observability-and-cost-controls.md), [docs/mvp/phases.md](docs/mvp/phases.md), [docs/mvp/task-prompts.md](docs/mvp/task-prompts.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, future-agent prompt, and observability behavior documentation synchronized after adding Phase 9 runtime behavior and the invalid-token guard.
- How: Created the dedicated Phase 9 artifact, checked off the Phase 9 tracker, linked the new Go package, changed the reusable observability prompt to “implement or refine,” documented `ErrInvalidTokenUsage` and `model_call.rejected`, and kept external telemetry, persistent logs, live providers, billing reconciliation, live routing, cache storage, CLI, HTTP API, and database behavior out of scope.

### Task: Update Repository State Documentation

- What: Updated repository overview and scope language to include the Phase 9 observability package, targeted test command, package-level structured events, caller-supplied token usage, invalid token usage, budget-limit failures, cache candidates, and model-routing notes.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/scope.md](docs/mvp/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Remove stale statements that observability, token tracking, and cost controls were unimplemented while preserving the limits around external telemetry, persistence, live providers, billing, real integrations, CLI, HTTP API, and database behavior.
- How: Added the Phase 9 artifact and package to documentation maps, marked observability and cost-control thinking as implemented, recorded the Phase 9 trust boundary, and clarified that the current surface is in-memory and package-level only.

### Task: Verify Phase 9 Observability Surface

- What: Verified the targeted observability package, full Go test suite, vet checks, and package coverage after implementation and documentation sync.
- Where: [internal/observability](internal/observability), [docs/mvp/observability-and-cost-controls.md](docs/mvp/observability-and-cost-controls.md), [CHANGELOG.md](CHANGELOG.md).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 9 is locally repeatable and does not break earlier MVP phase packages.
- How: Ran `go test ./internal/observability`, `go test ./...`, `go vet ./...`, and `go test -cover ./...`; all completed successfully, with `internal/observability` reporting 94.3% statement coverage.

## 2026-05-06 - Phase 8 Eval Harness

### Task: Add Phase 8 Eval TDD Coverage

- What: Added failing-first eval harness tests for loading golden normal, adversarial, and incomplete cases; scoring expected severity; scoring citation coverage; detecting unsupported claims; checking SOP-grounded recommendation accuracy; verifying redaction; checking prompt-injection resistance; and asserting strict release thresholds.
- Where: [internal/eval/eval_test.go](internal/eval/eval_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 8 with strict TDD and make deterministic eval behavior observable before adding production eval logic.
- How: Created same-package tests around `GoldenCases`, `Run`, and `DefaultThresholds`; ran `go test ./internal/eval`; observed the expected red build failure for missing eval APIs and types before adding implementation.

### Task: Implement Deterministic In-Memory Eval Harness

- What: Added an eval package that loads synthetic golden cases, retrieves mock guidance, composes timeline, severity, brief, and approval-gate behavior, and reports severity accuracy, citation coverage, recommendation accuracy, unsupported claims, redaction leaks, prompt-injection resistance, and approval fail-closed status.
- Where: [internal/eval/eval.go](internal/eval/eval.go).
- When: 2026-05-06, America/Vancouver.
- Why: Provide repeatable local release gates for the MVP demo path without adding a CLI, HTTP API, database, persistent eval storage, model calls, observability, token tracking, cost controls, export, escalation, or external-sharing integrations.
- How: Implemented `GoldenCases`, `Run`, `DefaultThresholds`, strict release-gate fields, in-memory synthetic packet fixtures for `FIC-SYN-001` through `FIC-SYN-005`, an approved mock guidance corpus, citation and recommendation scoring, prohibited-claim and sensitive-term scans, prompt-injection checks, and no-approval sensitive-action callback checks; then confirmed `go test ./internal/eval` passed.

### Task: Document Phase 8 Eval Contract

- What: Reworked the eval plan into the Phase 8 behavior contract, checked off Phase 8, documented golden cases, scoring rules, release thresholds, current limits, test commands, and red-to-green evidence.
- Where: [docs/mvp/eval-plan.md](docs/mvp/eval-plan.md), [docs/mvp/phases.md](docs/mvp/phases.md), [docs/mvp/task-prompts.md](docs/mvp/task-prompts.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, future-agent prompt, and eval behavior documentation synchronized after adding runtime eval harness behavior.
- How: Linked the new Go package from the Phase 8 tracker, changed the reusable eval prompt to “implement or refine,” recorded strict default thresholds, and kept observability, token and cost metrics, CLI reporting, HTTP API, persistence, and real integrations out of Phase 8 scope.

### Task: Update Repository State Documentation

- What: Updated repository overview and scope language to include the Phase 8 eval package and targeted test command while keeping observability and cost controls deferred.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/scope.md](docs/mvp/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Remove stale statements that said no eval harness existed and split eval completion from future observability work.
- How: Added the eval artifact and package to documentation maps, marked deterministic local evals as implemented, added the Phase 8 trust boundary, and clarified that trace IDs, structured logs, latency, token usage, budget metrics, CLI, HTTP API, database, and real integrations remain unimplemented.

### Task: Verify Phase 8 Eval Harness

- What: Verified the targeted eval package, full Go test suite, vet checks, and package coverage after the implementation and documentation sync.
- Where: [internal/eval](internal/eval), [README.md](README.md), [docs/mvp/eval-plan.md](docs/mvp/eval-plan.md).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 8 is locally repeatable and does not break earlier MVP phase packages.
- How: Ran `go test ./internal/eval`, `go test ./...`, `go vet ./...`, and `go test -cover ./...`; all completed successfully, with `internal/eval` reporting 86.0% statement coverage.

## 2026-05-06 - Phase 7 Human Approval Workflow

### Task: Add Phase 7 Approval TDD Coverage

- What: Added failing-first approval workflow tests for pending request creation, approver and decision capture, blocked pending actions, denied action blocking, approved scoped execution, out-of-scope blocking, mismatched call-and-scope incident blocking, append-only audit-copy behavior, and final-decision immutability.
- Where: [internal/approval/approval_test.go](internal/approval/approval_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 7 with strict TDD and make the human approval gate observable before adding production approval logic.
- How: Created same-package tests around `export`, `escalation`, and `external_sharing` using the existing `severity.SensitiveAction` labels; ran `go test ./internal/approval`; observed the expected red build failure for missing `NewGate`, request, scope, decision, audit, and gate types, then added focused red tests for immutable final decisions and mismatched call-and-scope incidents before tightening enforcement.

### Task: Implement In-Memory Human Approval Gate

- What: Added an approval package that creates in-memory approval requests, records approved or denied human decisions, blocks missing, pending, denied, out-of-scope, and mismatched call-and-scope sensitive action callbacks, allows approved callbacks only within exact scope, prevents final decisions from being rewritten in place, and returns append-only audit history copies.
- Where: [internal/approval/approval.go](internal/approval/approval.go).
- When: 2026-05-06, America/Vancouver.
- Why: Provide the first Phase 7 enforcement boundary without adding real export, escalation, external-sharing, persistence, identity, CLI, HTTP API, or external-service behavior.
- How: Implemented `NewGate`, `CreateRequest`, `Decide`, `Execute`, and `AuditHistory` with deterministic request IDs, injected clock timestamps, `DecisionPending`, `DecisionApproved`, `DecisionDenied`, `ErrActionBlocked`, `ErrRequestAlreadyDecided`, scoped request matching, and audit events for requested, decided, blocked, and allowed actions; then confirmed `go test ./internal/approval` passed.

### Task: Document Phase 7 Approval Contract

- What: Added the Phase 7 approval workflow contract, runtime surface, approval request model, enforcement rules, audit history behavior, current limits, test command, and red-to-green evidence.
- Where: [docs/mvp/human-approval-workflow.md](docs/mvp/human-approval-workflow.md), [docs/mvp/phases.md](docs/mvp/phases.md), [docs/mvp/eval-plan.md](docs/mvp/eval-plan.md), [docs/mvp/task-prompts.md](docs/mvp/task-prompts.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, eval planning, future-agent prompt, and behavior documentation synchronized after adding approval workflow runtime behavior.
- How: Created the dedicated Phase 7 artifact, checked off the Phase 7 tracker, linked the new Go package, added Phase 7 eval mappings for pending, denied, approved, scoped, and immutable-audit cases, and pointed the reusable approval prompt at the new contract.

### Task: Update Repository State Documentation

- What: Updated repository overview and scope language to include the Phase 7 approval package and targeted test command without claiming persistence, identity, real export, real escalation, external sharing, observability, eval harness, CLI, HTTP API, or database behavior exists.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/scope.md](docs/mvp/scope.md), [docs/mvp/product-frame.md](docs/mvp/product-frame.md), [docs/mvp/synthetic-incident-packets.md](docs/mvp/synthetic-incident-packets.md), [docs/mvp/severity-classification-and-recommended-actions.md](docs/mvp/severity-classification-and-recommended-actions.md), [docs/mvp/incident-brief-drafting.md](docs/mvp/incident-brief-drafting.md).
- When: 2026-05-06, America/Vancouver.
- Why: Remove stale statements that human approval records and enforcement were future work while preserving the remaining limits around actual sensitive-action integrations and production runtime surfaces.
- How: Added the Phase 7 artifact to documentation maps, marked human approval gating as implemented, recorded the fail-closed approval trust boundary, clarified Phase 5 approval flags and Phase 6 approval-state display as separate from Phase 7 decision records, and kept real export, escalation, and external sharing out of implemented scope.

## 2026-05-06 - Phase 6 Shareable Incident Brief

### Task: Add Phase 6 Brief TDD Coverage

- What: Added failing-first brief tests for complete draft output, citation coverage, sensitive-field redaction, missing-evidence fail-closed behavior, uncertainty labeling, approval-state display, and unsupported-claim omission.
- Where: [internal/brief/brief_test.go](internal/brief/brief_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 6 with strict TDD and make cited, redacted, human-reviewable brief behavior observable before adding production drafting logic.
- How: Created synthetic in-test packet fixtures using the existing ingestion, retrieval, timeline, and severity contracts; ran `go test ./internal/brief`; observed the expected red build failure for missing `Draft`, `StatusDraft`, result, section, redaction, approval-state, and fail-closed error types.

### Task: Implement Structured Brief Drafting

- What: Added a brief package that drafts structured incident summaries, cited timelines, severity rationale, recommended actions, and approval-state sections; carries source references; redacts vehicle, route, location, GPS-label, sensitive transcript, sensitive still-frame, and coordinate-like text; labels uncertainty; and fails closed when required evidence is missing.
- Where: [internal/brief/brief.go](internal/brief/brief.go).
- When: 2026-05-06, America/Vancouver.
- Why: Provide the first Phase 6 shareable draft surface without rendering documents, approving actions, exporting evidence, escalating incidents, or sharing externally.
- How: Implemented `Draft(packet ingestion.Packet, timelineResult timeline.Result, severityResult severity.Result) (Result, error)` with deterministic section assembly, source conversion, required-evidence validation, `MissingEvidenceError`, structured redaction tracking, blocked approval-state display, and `StatusDraft`; then confirmed `go test ./internal/brief` passed.

### Task: Document Phase 6 Brief Contract

- What: Added the Phase 6 brief drafting contract, runtime surface, section list, citation rules, redaction rules, approval-state behavior, current limits, test command, and red-to-green evidence.
- Where: [docs/mvp/incident-brief-drafting.md](docs/mvp/incident-brief-drafting.md), [docs/mvp/phases.md](docs/mvp/phases.md), [docs/mvp/eval-plan.md](docs/mvp/eval-plan.md), [docs/mvp/task-prompts.md](docs/mvp/task-prompts.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, eval planning, future-agent prompts, and behavior documentation synchronized after adding brief drafting runtime behavior.
- How: Created the dedicated Phase 6 artifact, checked off the Phase 6 tracker, linked the new Go package, added Phase 6 eval mappings for citations, redaction, fail-closed behavior, uncertainty, and approval-state display, and pointed the reusable brief prompt at the new contract.

### Task: Update Repository State Documentation

- What: Updated repository overview and scope language to include the Phase 6 brief package and targeted test command without claiming a renderer, approval workflow, eval harness, export, escalation, external sharing, CLI, HTTP API, database, or persistence exists.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/scope.md](docs/mvp/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Avoid stale statements that said no brief generator existed while preserving the remaining workflow boundaries and later-phase limits.
- How: Added the Phase 6 artifact to documentation maps, marked cited and redacted draft brief behavior as implemented, recorded the shareable-output redaction trust boundary, and kept human approval workflow and sensitive-action execution out of implemented scope.

## 2026-05-06 - Phase 5 Severity Classification And Recommended Actions

### Task: Add Phase 5 Severity TDD Coverage

- What: Added failing-first severity tests for low hard-brake classification, medium stop-arm classification, high collision classification, unknown incomplete-evidence classification, conflicting telemetry handling, recommendation explanations, SOP citation grounding, adversarial transcript handling, deterministic no-model behavior, and approval-required flags.
- Where: [internal/severity/severity_test.go](internal/severity/severity_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 5 with strict TDD and make severity, recommendation, and approval-flag behavior observable before adding production classifier logic.
- How: Created synthetic in-test packet fixtures using the existing ingestion, retrieval, and timeline contracts; ran `go test ./internal/severity`; observed the expected red build failure for missing `Classify`, result, recommendation, source, and approval types.

### Task: Implement Deterministic Severity And Recommendation Rules

- What: Added a severity package that classifies low, medium, high, and unknown severity with deterministic rules, treats conflicting timeline telemetry as unknown, treats adversarial transcript content as untrusted data, returns explained recommendations with packet and retrieved-guidance sources, records `ModelJudgmentUsed: false`, and marks export, escalation, and external sharing as approval-required but not approved.
- Where: [internal/severity/severity.go](internal/severity/severity.go).
- When: 2026-05-06, America/Vancouver.
- Why: Provide the first explainable Phase 5 reasoning surface without calling a model, executing sensitive actions, or introducing approval workflow behavior before Phase 7.
- How: Implemented `Classify(packet ingestion.Packet, timelineResult timeline.Result, guidance retrieval.Result) Result` with event-type rules, timeline-conflict detection, `retrieved_data` citation filtering, recommendation action labels, source deduplication, and fail-closed sensitive-action approval flags; then confirmed `go test ./internal/severity` passed.

### Task: Document Phase 5 Severity Contract

- What: Added the Phase 5 severity and recommendation contract, runtime surface, deterministic rule table, recommendation output shape, approval-required flag behavior, current limits, test command, and red-to-green evidence.
- Where: [docs/mvp/severity-classification-and-recommended-actions.md](docs/mvp/severity-classification-and-recommended-actions.md), [docs/mvp/phases.md](docs/mvp/phases.md), [docs/mvp/eval-plan.md](docs/mvp/eval-plan.md), [docs/mvp/task-prompts.md](docs/mvp/task-prompts.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, eval planning, future-agent prompts, and behavior documentation synchronized after adding severity runtime behavior.
- How: Created the dedicated Phase 5 artifact, checked off the Phase 5 tracker, linked the new Go package, added Phase 5 eval mappings for severity and recommendations, and pointed the reusable severity prompt at the new contract.

### Task: Update Repository State Documentation

- What: Updated repository overview and scope language to include the Phase 5 severity package and targeted test command without claiming a brief generator, approval workflow, export, escalation, external sharing, observability, or eval harness exists.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/scope.md](docs/mvp/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Avoid stale statements that said no severity classifier existed while preserving the remaining workflow boundaries and later-phase limits.
- How: Added the Phase 5 artifact to documentation maps, marked severity and recommendations as implemented, recorded the deterministic severity trust boundary, and kept human approval workflow and sensitive-action execution out of implemented scope.

## 2026-05-06 - Phase 4 Incident Timeline Builder

### Task: Add Phase 4 Timeline TDD Coverage

- What: Added failing-first timeline tests for chronological telemetry ordering, source citation coverage, guidance citation carry-forward, transcript and still-frame attribution, missing-evidence uncertainty, conflicting telemetry uncertainty, and unsupported-claim omission.
- Where: [internal/timeline/timeline_test.go](internal/timeline/timeline_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 4 with strict TDD and make grounded timeline behavior observable before adding production timeline logic.
- How: Created synthetic packet fixtures using the existing ingestion and retrieval contracts, ran `go test ./internal/timeline`, and observed the expected red build failure for missing `Build`, `Entry`, and timeline result types.

### Task: Implement Deterministic Incident Timeline Builder

- What: Added a timeline package that builds chronological entries from validated packet telemetry, labels transcript and still-frame claims by source type, preserves structured packet source references, carries approved retrieval citation references as guidance metadata, marks unavailable evidence as uncertain, and labels conflicting same-time telemetry.
- Where: [internal/timeline/timeline.go](internal/timeline/timeline.go).
- When: 2026-05-06, America/Vancouver.
- Why: Provide the first brief-ready grounded timeline surface without inventing facts, approving sensitive actions, or treating retrieved text as instructions.
- How: Implemented `Build(packet ingestion.Packet, guidance retrieval.Result) Result` with deterministic in-memory ordering, `packet.*[N]` source references, `retrieved_data` citation filtering, unavailable-evidence detection, and conflict labeling; then confirmed `go test ./internal/timeline` passed.

### Task: Document Phase 4 Timeline Contract

- What: Added the Phase 4 timeline builder contract, output format, input dependencies, source-reference rules, uncertainty behavior, unsupported-claim behavior, current limits, test command, and red-to-green evidence.
- Where: [docs/mvp/incident-timeline-builder.md](docs/mvp/incident-timeline-builder.md), [docs/mvp/phases.md](docs/mvp/phases.md), [docs/mvp/eval-plan.md](docs/mvp/eval-plan.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, eval planning, and behavior documentation synchronized after adding timeline runtime behavior.
- How: Created the dedicated Phase 4 artifact, checked off the Phase 4 tracker, linked the new Go package, and added future eval checks for ordering, citation coverage, uncertainty, conflicts, and unsupported claims.

### Task: Update Repository State Documentation

- What: Updated repository overview language to include the Phase 4 timeline package and its targeted test command without claiming severity, brief, approval, export, escalation, or external sharing behavior.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/scope.md](docs/mvp/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Avoid stale statements that said no timeline builder existed while preserving the remaining unimplemented workflow boundaries.
- How: Added the Phase 4 artifact to documentation maps, marked cited timeline generation as implemented, recorded the new timeline trust boundary, and verified `go test ./internal/timeline`, `go test ./internal/ingestion ./internal/retrieval ./internal/timeline`, and `go test ./...` passed.

## 2026-05-06 - Phase 3 RAG Corpus And Grounding

### Task: Add Phase 3 Retrieval TDD Coverage

- What: Added failing-first retrieval tests for relevant SOP retrieval, citation metadata, no-match behavior, workflow and scope filtering, hostile retrieved text, and deterministic result limiting.
- Where: [internal/retrieval/retrieval_test.go](internal/retrieval/retrieval_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 3 with strict TDD and make the grounding boundary observable before adding production retrieval logic.
- How: Created synthetic in-test mock SOP and troubleshooting documents, ran `go test ./internal/retrieval`, and observed the expected red build failure for missing retrieval APIs and types.

### Task: Implement Deterministic Mock Guidance Retrieval

- What: Added an in-memory retrieval package with document/query/result types, exact workflow and scope filtering, lexical scoring, stable citation references, snippets, and `retrieved_data` content-role marking.
- Where: [internal/retrieval/retrieval.go](internal/retrieval/retrieval.go).
- When: 2026-05-06, America/Vancouver.
- Why: Allow future reasoning phases to retrieve only approved mock guidance with citations while keeping retrieved text untrusted.
- How: Implemented `NewRetriever` and `Retrieve` with fail-closed empty results for empty workflow or scope, filtering before ranking, deterministic ordering by score and source ID, and citation references formatted as `SOURCE_ID#YYYY-MM-DD`.

### Task: Document Phase 3 Corpus And Grounding Contract

- What: Added the Phase 3 mock RAG corpus contract, document metadata, mock SOPs, troubleshooting notes, citation format, no-match behavior, scope filtering behavior, prompt-injection fixture, eval questions, red-to-green evidence, and current limits.
- Where: [docs/mvp/rag-corpus-and-grounding.md](docs/mvp/rag-corpus-and-grounding.md), [docs/mvp/phases.md](docs/mvp/phases.md), [docs/mvp/eval-plan.md](docs/mvp/eval-plan.md), [docs/mvp/task-prompts.md](docs/mvp/task-prompts.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, eval planning, future-agent prompts, and implementation behavior synchronized after adding the retrieval slice.
- How: Created the dedicated Phase 3 artifact, checked off the Phase 3 tracker, added retrieval eval mappings, and updated the reusable RAG planning prompt to point at the new source document.

### Task: Update Repository State Documentation

- What: Updated repository overview language to include the Phase 3 retrieval package and its targeted test command without claiming a complete end-to-end incident workflow.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/scope.md](docs/mvp/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Avoid stale statements that said retrieval did not exist while preserving limits around timeline, severity, brief, approval, export, escalation, and external sharing.
- How: Added the Phase 3 artifact to documentation maps, updated runtime-surface notes, marked the retrieved-content trust boundary as implemented, and kept later workflow phases unchecked.

## 2026-05-06 - Phase 2 Incident Packet Ingestion

### Task: Add Phase 2 Go Module And TDD Coverage

- What: Added the initial Go module and failing-first ingestion tests for valid packets, missing required fields, non-synthetic records, malformed incident IDs, malformed telemetry, unsupported event types, non-synthetic evidence references, and audit event emission.
- Where: [go.mod](go.mod), [internal/ingestion/ingestion_test.go](internal/ingestion/ingestion_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 2 with strict TDD and make packet validation observable before adding production ingestion logic.
- How: Created tests using synthetic inline packet JSON from the Phase 1 contract, observed red failures for missing APIs and later missing validation rules, then used those tests to drive the implementation.

### Task: Implement Synthetic Packet Validation

- What: Added JSON ingestion, packet schema types, deterministic validation errors, supported event-type validation, synthetic-only evidence validation, telemetry validation, and accepted/rejected audit events.
- Where: [internal/ingestion/ingestion.go](internal/ingestion/ingestion.go).
- When: 2026-05-06, America/Vancouver.
- Why: Ensure no downstream reasoning step can use unvalidated or non-synthetic incident packet data.
- How: Implemented `IngestJSON` with RFC3339 timestamp parsing, required-field checks, `FIC-SYN-` incident ID enforcement, speed bounds, relative-time parsing, `synthetic://` media reference enforcement, and stable validation issue codes/messages.

### Task: Document Phase 2 Ingestion Contract

- What: Added the Phase 2 ingestion schema, validation rules, audit event contract, test commands, red-to-green evidence, and current runtime limits.
- Where: [docs/mvp/incident-packet-ingestion.md](docs/mvp/incident-packet-ingestion.md), [docs/mvp/phases.md](docs/mvp/phases.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep implementation behavior, acceptance criteria, and phase status synchronized after adding runtime code.
- How: Created a dedicated Phase 2 artifact and checked off the Phase 2 tracker with links to the document and Go package.

### Task: Update Repository State Documentation

- What: Replaced documentation-only current-state language with an accurate Phase 2 runtime summary and linked the ingestion artifact.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/scope.md](docs/mvp/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Avoid overclaiming a complete app while no longer claiming the repo has no application code or tests.
- How: Added runtime-surface notes, marked only the ingestion validation promise complete, and recorded the implemented untrusted-packet boundary while leaving later phases unchecked.

## 2026-05-06 - Phase 1 Synthetic Evidence And Workflow Contract

### Task: Create Synthetic Incident Packet Contract

- What: Added the Phase 1 synthetic incident packet specs and shared workflow-output contract.
- Where: [docs/mvp/synthetic-incident-packets.md](docs/mvp/synthetic-incident-packets.md).
- When: 2026-05-06, America/Vancouver.
- Why: Define realistic fake evidence and expected workflow outputs before any ingestion, reasoning, eval, or fixture implementation begins.
- How: Created five explicitly synthetic records covering low, medium, high, unknown, and adversarial or missing-data scenarios with required fields, telemetry samples, media references, transcript notes, still-frame notes, expected timelines, severity, recommendations, brief requirements, missing-data behavior, and acceptance criteria.

### Task: Mark Phase 1 Complete

- What: Checked off the Phase 1 planning checklist and linked the synthetic packet output artifact.
- Where: [docs/mvp/phases.md](docs/mvp/phases.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker aligned with completed Markdown planning work while avoiding any claim that runtime behavior exists.
- How: Marked only the Phase 1 documentation items complete and added a dated note that Phase 1 added Markdown planning artifacts only.

### Task: Update Documentation Indexes

- What: Linked the Phase 1 packet contract from the root and MVP README files, and refreshed the current-state language.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md).
- When: 2026-05-06, America/Vancouver.
- Why: Make the new Phase 1 artifact discoverable and keep the repo narrative clear that no application code, tests, scaffolding, dependencies, or runtime configuration were added.
- How: Added the synthetic packet artifact to the planning artifact maps and changed the current-state wording from Phase 0-only planning to Phase 0 plus Phase 1 planning.

### Task: Sync Future Prompt And Eval References

- What: Pointed future documentation agents and eval planning at the Phase 1 packet contract.
- Where: [docs/mvp/task-prompts.md](docs/mvp/task-prompts.md), [docs/mvp/eval-plan.md](docs/mvp/eval-plan.md).
- When: 2026-05-06, America/Vancouver.
- Why: Preserve the Phase 1 source of truth for future fixture, eval, and implementation work without creating machine-readable fixtures prematurely.
- How: Updated the synthetic incident planning prompt to name the Markdown contract file and added a Phase 1 packet mapping for future eval fixture work.

## 2026-05-06 - Phase 0 Product Frame And Guardrails

### Task: Create Product Frame

- What: Added the Phase 0 product frame for Fleet Incident Copilot.
- Where: [docs/mvp/product-frame.md](docs/mvp/product-frame.md).
- When: 2026-05-06, America/Vancouver.
- Why: Establish the primary user, problem, MVP promise, scope boundaries, trust boundaries, approval gates, prohibited claims, demo narrative, and success criteria before implementation.
- How: Grounded the frame in [research-report.md](research-report.md), kept all data synthetic, and explicitly confirmed no code is needed for Phase 0.

### Task: Mark Phase 0 Complete

- What: Checked off the Phase 0 planning checklist and linked the output artifact.
- Where: [docs/mvp/phases.md](docs/mvp/phases.md).
- When: 2026-05-06, America/Vancouver.
- Why: Make the phase tracker reflect completed planning work without implying runtime behavior exists.
- How: Marked only the Phase 0 documentation items complete and added a dated note that Phase 0 added Markdown planning artifacts only.

### Task: Update Documentation Indexes

- What: Linked the Product Frame from the root and MVP README files, and recorded the Go backend direction.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep navigation current and capture that future backend code phases should use Go.
- How: Added the Product Frame artifact to the planning artifact maps and clarified that current behavior is planning-only.

### Task: Refine Scope Guardrails

- What: Added explicit prohibited claims and approval-gate rules.
- Where: [docs/mvp/scope.md](docs/mvp/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Make the review boundaries visible before any implementation agent starts future work.
- How: Added scope notes that distinguish planning checklists from implemented behavior, recorded Go as the future backend language, and made export, escalation, and external sharing approval requirements explicit.
