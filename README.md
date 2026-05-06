# Fleet Incident Copilot MVP

This workspace turns [research-report.md](research-report.md) into an actionable MVP plan for a Fleet Incident Copilot: a production-minded AI application demo for fleet-safety incident review.

The current runtime surface is intentionally small. Phase 0 and Phase 1 planning artifacts exist, Phase 2 adds a strict-TDD Go ingestion package that validates synthetic incident packet JSON, Phase 3 adds a strict-TDD Go retrieval package that returns cited snippets from approved mock guidance, Phase 4 adds a strict-TDD Go timeline package that produces cited synthetic incident timelines, Phase 5 adds a strict-TDD Go severity package that returns deterministic severity, SOP-grounded recommendations, and approval-required flags, Phase 6 adds a strict-TDD Go brief package that returns cited, redacted, draft incident briefs for human review, Phase 7 adds a strict-TDD Go approval package that creates in-memory approval requests, records human decisions, blocks pending or denied sensitive action callbacks, and allows approved callbacks only within scope, Phase 8 adds a strict-TDD Go eval package that scores deterministic synthetic golden cases for severity, citations, recommendations, unsupported claims, redaction, prompt-injection resistance, and approval fail-closed behavior, Phase 9 adds a strict-TDD Go observability package that records in-memory structured workflow events, caller-supplied token usage, invalid token usage, budget-limit failures, cache candidates, and model-routing notes, and Phase 10 adds Markdown demo packaging materials that distinguish implemented package-level behavior from planned production integrations. No CLI, HTTP API, database, external observability pipeline, persistent log store, real model-provider call, provider billing reconciliation, real export tool, real escalation tool, or external-sharing integration exists yet.

## Planning Artifacts

- [MVP Overview](docs/mvp/README.md)
- [Product Frame](docs/mvp/product-frame.md)
- [Synthetic Incident Packets](docs/mvp/synthetic-incident-packets.md)
- [Incident Packet Ingestion](docs/mvp/incident-packet-ingestion.md)
- [RAG Corpus And Grounding](docs/mvp/rag-corpus-and-grounding.md)
- [Incident Timeline Builder](docs/mvp/incident-timeline-builder.md)
- [Severity Classification And Recommended Actions](docs/mvp/severity-classification-and-recommended-actions.md)
- [Shareable Incident Brief Drafting](docs/mvp/incident-brief-drafting.md)
- [Human Approval Workflow](docs/mvp/human-approval-workflow.md)
- [Eval Plan](docs/mvp/eval-plan.md)
- [Observability And Cost Controls](docs/mvp/observability-and-cost-controls.md)
- [Scope And Guardrails](docs/mvp/scope.md)
- [Phases And Tasks](docs/mvp/phases.md)
- [Task Prompts](docs/mvp/task-prompts.md)
- [Strict TDD Rules](docs/mvp/tdd-rules.md)
- [Demo Package](docs/mvp/demo-package.md)

## Working Rules

- [x] Base the MVP on the research report.
- [x] Keep all planning data synthetic.
- [x] Require human approval before any future export or escalation behavior.
- [x] Require strict TDD for every coding task.
- [x] Use Go for backend implementation when code phases begin.
- [x] Validate synthetic incident packets before downstream reasoning can use them.
- [x] Retrieve approved mock guidance with citation metadata and scope filtering.
- [x] Build cited synthetic incident timelines from validated packet data.
- [x] Classify severity and recommend next actions with deterministic rules.
- [x] Draft cited, redacted incident briefs for human review.
- [x] Create in-memory human approval records and gate sensitive action callbacks.
- [x] Run deterministic local evals for severity, citations, recommendations, unsupported claims, redaction, prompt-injection resistance, and approval fail-closed behavior.
- [x] Record package-level observability events, caller-supplied token usage, invalid token usage, budget-limit failures, cache candidates, and model-routing notes.
- [ ] Use the phase checklist to drive future implementation.
- [ ] Keep future implementation notes synchronized with the docs when behavior changes.

## Runtime Surface

- Go module: [go.mod](go.mod).
- Ingestion package: [internal/ingestion](internal/ingestion).
- Retrieval package: [internal/retrieval](internal/retrieval).
- Timeline package: [internal/timeline](internal/timeline).
- Severity package: [internal/severity](internal/severity).
- Brief package: [internal/brief](internal/brief).
- Approval package: [internal/approval](internal/approval).
- Eval package: [internal/eval](internal/eval).
- Observability package: [internal/observability](internal/observability).
- Targeted tests: `go test ./internal/ingestion`, `go test ./internal/retrieval`, `go test ./internal/timeline`, `go test ./internal/severity`, `go test ./internal/brief`, `go test ./internal/approval`, `go test ./internal/eval`, and `go test ./internal/observability`.
- Full Go test command: `go test ./...`.
