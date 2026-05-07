# EvalOps Task Prompts — Fleet Incident Copilot

## Reusable code prompt

> You are working in `sf-mvp`. Follow the existing strict TDD rules plus `docs/overlays/evalops-tdd-addendum.md`. Read the current package tests before editing. Add a failing test first, run it, confirm red, implement the smallest code, run the narrow test green, then run `go test ./...`. Do not add real fleet data, live model calls, or external service assumptions. Preserve the implemented-versus-planned wording discipline in docs.

## FQ11-T02 Code: JSONL exporter for eval cases

> Add a JSONL exporter for existing synthetic eval cases. Start with failing tests for deterministic ordering, stable case IDs, expected severity/citation fields, adversarial tags, and no raw sensitive evidence leakage. Confirm red. Implement the exporter as a package function, not a CLI-only behavior. Acceptance requires a golden JSONL fixture and `go test ./internal/eval` passing.

## FQ11-T03 Code: shared result importer

> Add a shared result importer that can read Promptfoo/evalops-style result JSON and convert it into the existing eval scoring model. Start with failing tests for valid result, missing case ID, malformed score, unknown scorer, and duplicate result. Confirm red. Implement minimal parsing and validation. Acceptance requires actionable errors and no behavior change to existing eval scoring.

## FQ12-T01 Code: local incident eval target

> Implement a local HTTP target for incident evals. Use `httptest` failing tests for valid synthetic packet, invalid JSON, missing incident ID, prompt-injection fixture, timeout, and response shape. Confirm red. Implement the handler through interfaces for ingestion, retrieval, timeline, severity, brief, approval, and eval packages. Acceptance requires deterministic output and no live network or model calls.

## FQ12-T03 Code: score adapter output

> Adapt existing scores into the Promptfoo/evalops response shape. Add failing tests for severity score, citation coverage, unsupported-claim detection, redaction, recommendation accuracy, prompt-injection resistance, and approval fail-closed behavior. Confirm red. Implement score mapping and severity labels. Acceptance requires critical safety failures to be machine-readable.

## FQ13-T01 Code: workflow span attributes

> Implement safe OpenTelemetry-style workflow attribute mapping. Add failing tests for trace ID propagation, hashed incident ID, retrieved source IDs, severity label, approval state, latency, and redaction of raw transcript/media/still-frame details. Confirm red. Implement mapping as pure functions first. Acceptance requires no raw evidence in attributes.

## FQ13-T02 Code: eval score event export

> Implement score event export for eval runs. Add failing tests for stable event names, score values, pass/fail flags, critical severity, trace correlation, and disabled exporter mode. Confirm red. Implement no-op exporter default plus interface for future OTel/Langfuse exporters. Acceptance requires exporter failures not to break the core incident workflow unless running in release-gate mode.

## FQ14-T01 Code: gate thresholds

> Implement release thresholds for incident quality. Add failing tests for zero critical failures, minimum citation coverage, approval fail-closed, no redaction failures, minimum severity accuracy, and malformed threshold config. Confirm red. Implement threshold evaluation behind a package function. Acceptance requires non-zero gate result for critical failures and a human-readable reason.

## FQ14-T03 Code: GitHub summary report

> Implement Markdown summary suitable for CI. Add golden-file tests for all-pass, warning-only, and critical-failure runs. Include failed case IDs, scorer names, remediation hints, and commands to reproduce. Confirm red. Implement deterministic sorting. Acceptance requires stable diff-friendly output.

## FQ15-T02 Code: draft case generator

> Implement conversion from redacted review/trace samples into draft eval cases. Add failing tests for required TODO expected fields, tag inheritance, deduplication, trace ID preservation, and explicit `review_required=true`. Confirm red. Implement generator. Acceptance requires generated cases to be non-blocking until human-reviewed.
