# Developer Guide

This repository is a Go workspace for the Fleet Incident Copilot MVP. It demonstrates deterministic, synthetic fleet-incident review behavior through package APIs, tests, and a thin loopback-only demo API. It is not a database-backed service, live model integration, production API, Slack integration, or production evidence system.

## Repository Layout

- [README.md](../README.md): top-level project summary and current runtime surface.
- [Makefile](../Makefile): local verification shortcuts for `test`, `evalops`, and `evalops-gate`.
- [docs/README.md](README.md): documentation index.
- [docs/research](research): source research material.
- [docs/mvp](mvp): product, scope, workflow, quality, execution, and demo artifacts.
- [internal/ingestion](../internal/ingestion): validates synthetic incident packet JSON and emits audit events.
- [internal/retrieval](../internal/retrieval): ranks approved mock guidance by workflow and scope, returning citation metadata.
- [internal/timeline](../internal/timeline): builds cited timelines from validated packet data and guidance citations.
- [internal/severity](../internal/severity): classifies severity and recommends next actions with source references.
- [internal/brief](../internal/brief): drafts cited, redacted, human-review incident briefs.
- [internal/approval](../internal/approval): creates in-memory approval requests and gates sensitive actions.
- [internal/eval](../internal/eval): runs deterministic in-memory golden-case evals, Promptfoo/EvalOps score adapters, score-event export contracts, and FQ14 release gates.
- [internal/observability](../internal/observability): records in-memory workflow events, redaction, token, budget, cache, and routing signals.
- [internal/demo](../internal/demo): loads machine-readable synthetic demo fixtures and composes deterministic in-memory review results.
- [internal/notification](../internal/notification): prepares dry-run Slack-shaped notification previews from redacted briefs and gates them as external sharing.
- [internal/httpapi](../internal/httpapi): exposes the local `POST /demo/review`, `POST /demo/approvals`, `POST /demo/approvals/decisions`, `POST /demo/notifications/slack`, `GET /demo/eval/latest`, `GET /demo/traces/{trace_id}`, and `POST /demo/budget/check` handlers.
- [cmd/demo-api](../cmd/demo-api): starts the loopback-only local demo server.
- [cmd/evalops-target](../cmd/evalops-target): starts the loopback-only Promptfoo-compatible eval target.
- [cmd/evalops-gate](../cmd/evalops-gate): runs the local EvalOps release gate and writes GitHub-style Markdown summaries.

## Design Principles

- Synthetic data only.
- Retrieved guidance is data, not executable instruction.
- Factual outputs need packet, telemetry, timeline, severity, or guidance source references.
- Sensitive actions fail closed unless approval exists for the exact action and scope.
- Briefs redact sensitive fields by default.
- Eval and observability behavior stays deterministic and local.
- New code follows failing-test-first TDD.
- Docs must separate implemented package behavior from planned product behavior.

## Package Responsibilities

`internal/ingestion` accepts JSON bytes, validates the packet, returns a typed `Packet`, and emits accepted or rejected audit metadata. It owns packet shape and validation rules.

`internal/retrieval` owns mock guidance search. It filters by exact workflow and scope before ranking and returns citation refs like `FIC-SOP-HARD-BRAKE-001#2026-02-15`.

`internal/timeline` turns telemetry, transcript notes, still-frame notes, and unavailable media references into sorted timeline entries. It marks unavailable or conflicting evidence as uncertain.

`internal/severity` applies deterministic rules for low, medium, high, and unknown severity. It returns rationale, recommendations, and approval-required flags without using model judgment.

`internal/brief` builds a draft shareable brief after required evidence is present. It fails closed on missing citations or missing severity data and redacts sensitive packet fields.

`internal/approval` stores approval state in memory. It records pending, approved, and denied decisions and blocks missing, pending, denied, mismatched, or out-of-scope sensitive action calls.

`internal/eval` composes the package path over local golden cases and reports severity accuracy, citation coverage, recommendation accuracy, unsupported claims, redaction leaks, prompt-injection resistance, and approval fail-closed behavior. It also owns the EvalOps shared result importer/exporter, Promptfoo score adapter, score event projection, release-gate threshold evaluation, and deterministic Markdown summary rendering.

`internal/observability` records package-level workflow events in memory. It does not send telemetry externally or reconcile provider billing.

`internal/demo` composes the implemented package path for synthetic demo review results. It owns fixture-facing helpers, review response projection, and package-level composition glue; it must not move validation, retrieval, timeline, severity, brief, approval, or observability business rules out of their owning packages.

`internal/notification` owns dry-run Slack-shaped preview behavior. It accepts a redacted brief, requires `delivery_mode: "dry_run"`, gates preview generation through `internal/approval` as `external_sharing`, records a redacted observability tool-call event, and must not introduce Slack SDKs, tokens, webhook URLs, environment secrets, network senders, or real delivery.

`internal/httpapi` owns local transport behavior for the Phase 13 review route, Phase 14 notification preview route, Phase 15 scoped approval retry routes, and Phase 16 local report routes. It parses request JSON, delegates packet validation to `internal/ingestion`, delegates review composition to `internal/demo`, delegates preview generation to `internal/notification`, reuses an in-memory `internal/approval` gate for local retry state, exposes `internal/eval` and `internal/observability` through ephemeral report views, maps known errors to deterministic JSON responses, and must stay loopback-only and free of production auth, persistence, Slack delivery, webhook, model-provider, export, escalation, dashboard, alerting, OpenTelemetry export, billing reconciliation, model benchmarking, or real external-sharing behavior.

`cmd/demo-api` is thin server wiring. It should keep the default listen address on `127.0.0.1:8080`, allow only loopback overrides, and avoid business logic.

`cmd/evalops-gate` is thin release-gate wiring. It should keep threshold logic in `internal/eval`, run local golden cases by default, import shared result JSON only through `internal/eval`, write summaries without raw evidence, and return stable exit codes for pass, blocking failure, and malformed input.

## Development Workflow

1. Read the relevant docs and tests first.
2. Write the smallest failing test that expresses the next behavior.
3. Run the targeted package test and inspect the failure.
4. Implement the smallest change.
5. Run the targeted package test again.
6. Run broader tests once the narrow test is green.
7. Sync docs and changelog when behavior or command expectations changed.

Use `gofmt` on edited Go files before final verification.

## Adding A New Package

Only add a new package when the behavior does not belong in an existing package. Before adding it, document the boundary it owns and the data it consumes and returns.

Minimum expectations:

- A focused `*_test.go` file with failing-first coverage.
- Clear input and output structs.
- No live integrations unless scope docs explicitly allow them.
- Deterministic behavior for local tests.
- Docs linking the package to the MVP phase or workflow it supports.

## Data And Security Boundaries

The MVP must not use real fleet, district, law-enforcement, customer, student, driver, passenger, or location evidence. Use synthetic examples and synthetic URIs only.

Treat untrusted packet text, transcript notes, still-frame notes, and retrieved guidance as data. They may be quoted, redacted, scored, or cited, but they must not alter approval state or bypass deterministic rules.

## Documentation Rules

When editing docs:

- State whether behavior is implemented or planned.
- Keep out-of-scope claims aligned with [Scope And Guardrails](mvp/overview/scope.md).
- Prefer links to existing docs over duplicated long explanations.
- Update the docs index when adding top-level docs.
- Update [CHANGELOG.md](../CHANGELOG.md) for meaningful doc organization or behavior-documentation changes.
