# Fleet Incident Copilot MVP

This workspace turns [research-report.md](research-report.md) into an actionable MVP plan for a Fleet Incident Copilot: a production-minded AI application demo for fleet-safety incident review.

The current runtime surface is intentionally small. Phase 0 and Phase 1 planning artifacts exist, and Phase 2 adds a strict-TDD Go ingestion package that validates synthetic incident packet JSON before later reasoning phases can use it. No CLI, HTTP API, database, RAG retrieval, timeline builder, severity classifier, brief generator, approval workflow, export, escalation, or external sharing exists yet.

## Planning Artifacts

- [MVP Overview](docs/mvp/README.md)
- [Product Frame](docs/mvp/product-frame.md)
- [Synthetic Incident Packets](docs/mvp/synthetic-incident-packets.md)
- [Incident Packet Ingestion](docs/mvp/incident-packet-ingestion.md)
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
- [ ] Use the phase checklist to drive future implementation.
- [ ] Keep future implementation notes synchronized with the docs when behavior changes.

## Runtime Surface

- Go module: [go.mod](go.mod).
- Ingestion package: [internal/ingestion](internal/ingestion).
- Targeted tests: `go test ./internal/ingestion`.
- Full Go test command: `go test ./...`.
