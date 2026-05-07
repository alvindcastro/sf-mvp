# Fleet Incident Copilot MVP

This workspace turns [research-report.md](docs/research/research-report.md) into an actionable MVP plan for a Fleet Incident Copilot: a production-minded AI application demo for fleet-safety incident review.

The current runtime surface is intentionally small. Phase 0 and Phase 1 planning artifacts exist, Phase 2 adds a strict-TDD Go ingestion package that validates synthetic incident packet JSON, Phase 3 adds a strict-TDD Go retrieval package that returns cited snippets from approved mock guidance, Phase 4 adds a strict-TDD Go timeline package that produces cited synthetic incident timelines, Phase 5 adds a strict-TDD Go severity package that returns deterministic severity, SOP-grounded recommendations, and approval-required flags, Phase 6 adds a strict-TDD Go brief package that returns cited, redacted, draft incident briefs for human review, Phase 7 adds a strict-TDD Go approval package that creates in-memory approval requests, records human decisions, blocks pending or denied sensitive action callbacks, and allows approved callbacks only within scope, Phase 8 adds a strict-TDD Go eval package that scores deterministic synthetic golden cases for severity, citations, recommendations, unsupported claims, redaction, prompt-injection resistance, and approval fail-closed behavior, Phase 9 adds a strict-TDD Go observability package that records in-memory structured workflow events, caller-supplied token usage, invalid token usage, budget-limit failures, cache candidates, and model-routing notes, Phase 10 adds Markdown demo packaging materials that distinguish implemented package-level behavior from planned production integrations, Phase 11 adds a Markdown roadmap for future local demo surfaces, Phase 12 adds a strict-TDD in-memory demo review composer plus machine-readable synthetic fixtures, Phase 13 adds a strict-TDD loopback-only demo API for `POST /demo/review`, Phase 14 adds a strict-TDD dry-run Slack-shaped notification preview plus `POST /demo/notifications/slack`, Phase 15 adds a strict-TDD in-memory scoped approval retry demo through `POST /demo/approvals` and `POST /demo/approvals/decisions`, Phase 16 adds strict-TDD local eval, trace, and budget report routes through `GET /demo/eval/latest`, `GET /demo/traces/{trace_id}`, and `POST /demo/budget/check`, and Phase 17 refreshes the Markdown demo package with verified local API commands, response highlights, eval numbers, recording plan, and a `go test ./...` fallback. No database, external observability pipeline, persistent log store, real model-provider call, provider billing reconciliation, real export tool, real escalation tool, Slack delivery, webhook, identity, auth, production API, persistent approval store, dashboard, alerting, OpenTelemetry export, model benchmark, or external-sharing integration exists yet.

## Documentation

- [Docs Index](docs/README.md): source research and grouped MVP planning docs.
- [MVP Overview](docs/mvp/README.md): thesis, demo promise, proof points, and artifact map.
- [How-Tos](docs/how-tos.md): common local workflows for scope checks, package changes, synthetic fixtures, documentation updates, and review prep.
- [Developer Guide](docs/developer-guide.md): repository layout, package boundaries, development workflow, and documentation rules.
- [Nice To Knows](docs/nice-to-knows.md): practical context about the current package-level runtime and demo limits.
- [Testing](docs/testing.md): targeted and full Go test commands, TDD expectations, coverage notes, and doc-only verification.
- [Troubleshooting](docs/troubleshooting.md): common local development and package-behavior issues.
- [LinkedIn Post Drafts](docs/mvp/demo/linkedin-post-drafts.md): casual repo posts for hiring managers, CTOs, and engineering leaders.

## MVP Artifacts

Product and scope:

- [Product Frame](docs/mvp/overview/product-frame.md)
- [Scope And Guardrails](docs/mvp/overview/scope.md)

Workflow behavior:

- [Synthetic Incident Packets](docs/mvp/workflow/synthetic-incident-packets.md)
- [Incident Packet Ingestion](docs/mvp/workflow/incident-packet-ingestion.md)
- [RAG Corpus And Grounding](docs/mvp/workflow/rag-corpus-and-grounding.md)
- [Incident Timeline Builder](docs/mvp/workflow/incident-timeline-builder.md)
- [Severity Classification And Recommended Actions](docs/mvp/workflow/severity-classification-and-recommended-actions.md)
- [Shareable Incident Brief Drafting](docs/mvp/workflow/incident-brief-drafting.md)
- [Human Approval Workflow](docs/mvp/workflow/human-approval-workflow.md)
- [Review Composition Contract](docs/mvp/demo/review-composition-contract.md)
- [Loopback Demo API](docs/mvp/demo/loopback-demo-api.md)
- [Dry-Run Slack-Shaped Notification Preview](docs/mvp/demo/dry-run-slack-preview.md)
- [Scoped Approval Demo Retry](docs/mvp/demo/scoped-approval-retry.md)
- [Eval And Observability Demo Reports](docs/mvp/demo/eval-and-observability-reports.md)

Quality and operations:

- [Eval Plan](docs/mvp/quality/eval-plan.md)
- [Observability And Cost Controls](docs/mvp/quality/observability-and-cost-controls.md)

Execution and packaging:

- [Phases And Tasks](docs/mvp/execution/phases.md)
- [Task Prompts](docs/mvp/execution/task-prompts.md)
- [Strict TDD Rules](docs/mvp/execution/tdd-rules.md)
- [Demo Package](docs/mvp/demo/demo-package.md)
- [Demo Surface Roadmap](docs/mvp/demo/demo-surface-roadmap.md)

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
- [x] Compose one deterministic in-memory demo review result with validation status, citations, timeline, severity, recommendations, redacted brief, approval-required actions, and trace ID.
- [x] Expose the deterministic demo review result through a loopback-only local API.
- [x] Prepare a dry-run Slack-shaped notification preview from the redacted brief and block it as external sharing before scoped approval.
- [x] Retry the dry-run notification preview after an exact in-memory human approval while keeping missing, pending, denied, out-of-scope, and wrong-action approvals blocked.
- [x] Expose local eval, trace, and caller-supplied budget report routes while keeping report state in memory and ephemeral.
- [x] Document a verified local API walkthrough with exact `curl` commands and a `go test ./...` fallback.
- [ ] Use the phase checklist to drive future implementation.
- [ ] Keep future implementation notes synchronized with the docs when behavior changes.

## Local Demo Walkthrough

Start the loopback-only demo server:

```bash
go run ./cmd/demo-api -addr 127.0.0.1:18084
```

Then follow the verified walkthrough in [Demo Package](docs/mvp/demo/demo-package.md#local-demo-surface) or [How-Tos](docs/how-tos.md#how-to-run-the-loopback-demo-api): review one synthetic incident, show the dry-run Slack-shaped preview blocked before approval, create and approve the exact scoped approval, retry the preview, fetch the eval report, and fetch a trace report from the same server process.

Fallback proof path:

```bash
go test ./...
```

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
- Demo composer package: [internal/demo](internal/demo).
- Notification preview package: [internal/notification](internal/notification).
- Demo HTTP handler package: [internal/httpapi](internal/httpapi).
- Demo API command: [cmd/demo-api](cmd/demo-api).
- Targeted tests: `go test ./internal/ingestion`, `go test ./internal/retrieval`, `go test ./internal/timeline`, `go test ./internal/severity`, `go test ./internal/brief`, `go test ./internal/approval`, `go test ./internal/eval`, `go test ./internal/observability`, `go test ./internal/demo`, `go test ./internal/notification`, `go test ./internal/httpapi`, and `go test ./cmd/demo-api`.
- Full Go test command: `go test ./...`.
