# Developer Guide

This repository is a Go package workspace for the Fleet Incident Copilot MVP. It demonstrates deterministic, synthetic fleet-incident review behavior through package APIs and tests. It is not currently an application server, CLI, database-backed service, live model integration, or production evidence system.

## Repository Layout

- [README.md](../README.md): top-level project summary and current runtime surface.
- [docs/README.md](README.md): documentation index.
- [docs/research](research): source research material.
- [docs/mvp](mvp): product, scope, workflow, quality, execution, and demo artifacts.
- [internal/ingestion](../internal/ingestion): validates synthetic incident packet JSON and emits audit events.
- [internal/retrieval](../internal/retrieval): ranks approved mock guidance by workflow and scope, returning citation metadata.
- [internal/timeline](../internal/timeline): builds cited timelines from validated packet data and guidance citations.
- [internal/severity](../internal/severity): classifies severity and recommends next actions with source references.
- [internal/brief](../internal/brief): drafts cited, redacted, human-review incident briefs.
- [internal/approval](../internal/approval): creates in-memory approval requests and gates sensitive actions.
- [internal/eval](../internal/eval): runs deterministic in-memory golden-case evals.
- [internal/observability](../internal/observability): records in-memory workflow events, redaction, token, budget, cache, and routing signals.
- [internal/demo](../internal/demo): loads machine-readable synthetic demo fixtures and composes deterministic in-memory review results.

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

`internal/eval` composes the package path over local golden cases and reports severity accuracy, citation coverage, recommendation accuracy, unsupported claims, redaction leaks, prompt-injection resistance, and approval fail-closed behavior.

`internal/observability` records package-level workflow events in memory. It does not send telemetry externally or reconcile provider billing.

`internal/demo` composes the implemented package path for synthetic demo review results. It owns fixture-facing helpers, review response projection, and package-level composition glue; it must not move validation, retrieval, timeline, severity, brief, approval, or observability business rules out of their owning packages.

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
