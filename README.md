# Fleet Incident Copilot MVP

This workspace turns [research-report.md](research-report.md) into an actionable MVP plan for a Fleet Incident Copilot: a production-minded AI application demo for fleet-safety incident review.

The current runtime surface is intentionally small. Phase 0 and Phase 1 planning artifacts exist, Phase 2 adds a strict-TDD Go ingestion package that validates synthetic incident packet JSON, Phase 3 adds a strict-TDD Go retrieval package that returns cited snippets from approved mock guidance, and Phase 4 adds a strict-TDD Go timeline package that produces cited synthetic incident timelines. No CLI, HTTP API, database, severity classifier, brief generator, approval workflow, export, escalation, or external sharing exists yet.

## Planning Artifacts

- [MVP Overview](docs/mvp/README.md)
- [Product Frame](docs/mvp/product-frame.md)
- [Synthetic Incident Packets](docs/mvp/synthetic-incident-packets.md)
- [Incident Packet Ingestion](docs/mvp/incident-packet-ingestion.md)
- [RAG Corpus And Grounding](docs/mvp/rag-corpus-and-grounding.md)
- [Incident Timeline Builder](docs/mvp/incident-timeline-builder.md)
- [Scope And Guardrails](docs/mvp/scope.md)
- [Phases And Tasks](docs/mvp/phases.md)
- [Task Prompts](docs/mvp/task-prompts.md)
- [Strict TDD Rules](docs/mvp/tdd-rules.md)
- [Eval Plan](docs/mvp/eval-plan.md)
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
- [ ] Use the phase checklist to drive future implementation.
- [ ] Keep future implementation notes synchronized with the docs when behavior changes.

## Runtime Surface

- Go module: [go.mod](go.mod).
- Ingestion package: [internal/ingestion](internal/ingestion).
- Retrieval package: [internal/retrieval](internal/retrieval).
- Timeline package: [internal/timeline](internal/timeline).
- Targeted tests: `go test ./internal/ingestion`, `go test ./internal/retrieval`, and `go test ./internal/timeline`.
- Full Go test command: `go test ./...`.
