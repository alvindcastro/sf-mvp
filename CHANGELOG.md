# Changelog

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
