# Shareable Incident Brief Drafting

Phase 6 implements the first structured incident brief draft for Fleet Incident Copilot. It consumes already-validated synthetic packet data, a Phase 4 timeline result, and a Phase 5 severity result. It returns a cited, redacted draft for human review.

## Phase 6 Checklist

- [x] Include incident summary, cited timeline, severity, rationale, next actions, and approval state.
- [x] Include citations for factual claims.
- [x] Redact sensitive fields from shareable output.
- [x] Fail closed when required evidence is missing.
- [x] Label uncertainty clearly.
- [x] Keep the brief draft human-reviewable, not final by default.

## Runtime Surface

- Package: `internal/brief`.
- Entry point: `Draft(packet ingestion.Packet, timelineResult timeline.Result, severityResult severity.Result) (Result, error)`.
- Targeted tests: `go test ./internal/brief`.
- Broad Go tests: `go test ./...`.

`Draft` is deterministic and in-memory. It does not call a model, external service, database, renderer, file writer, approval workflow, export tool, escalation tool, sharing tool, CLI, or HTTP API.

## Brief Output

The result includes:

- `Status`: always `draft`.
- `IncidentID`: the synthetic incident ID.
- `SyntheticRecord`: copied from the packet.
- `Sections`: structured brief sections.
- `ApprovalState`: sensitive actions shown as blocked when approval is required and absent.
- `RedactionsApplied`: fields withheld from shareable draft text.
- `Uncertainties`: unique uncertainty labels carried from the timeline.
- `Shareable`: true only when the draft has required sections, citations, redactions, and blocked approval state.

The first Phase 6 sections are:

- `Incident Summary`
- `Cited Timeline`
- `Severity Rationale`
- `Recommended Actions`
- `Approval State`

## Citation Rules

Every factual section must include at least one structured source reference. The draft can cite:

- Packet fields such as `packet.incident_id`, `packet.event_type`, or `packet.telemetry_samples[N]`.
- Timeline sources carried from Phase 4.
- Severity rationale and recommendation sources from Phase 5.
- Retrieved guidance citation references that already passed Phase 3 filtering.
- Approval requirement references such as `severity.approval_requirements[N]`.

If required timeline, severity, recommendation, or approval evidence is missing, `Draft` returns `MissingEvidenceError` and a fail-closed result with `Shareable: false`.

## Redaction Rules

The shareable draft redacts:

- `packet.vehicle_id`
- `packet.route`
- `packet.location_label`
- `packet.telemetry_samples[N].gps_label`
- sensitive transcript notes, including passenger detail or hostile export instructions
- sensitive still-frame notes
- coordinate-like text

Redactions are recorded as structured `Redaction` values. Phase 6 does not implement policy configuration, user-specific redaction rules, or irreversible document export.

## Approval State

Phase 6 displays the approval flags produced by Phase 5:

- `export`: blocked pending human approval.
- `escalation`: blocked pending human approval.
- `external_sharing`: blocked pending human approval.

The draft does not approve, deny, export, escalate, share, or create approval records. The Phase 7 `internal/approval` package provides the separate human approval workflow and scoped enforcement gate.

## Red-To-Green Evidence

Smallest observable behavior: draft one complete human-reviewable incident brief from packet, timeline, and severity results while preserving citations and blocked approval state.

Observed red state:

- `go test ./internal/brief` failed because `Draft`, `StatusDraft`, `Result`, `Section`, `Redaction`, `ApprovalState`, and `MissingEvidenceError` were undefined.

Green state:

- `go test ./internal/brief` passes after adding the minimal `internal/brief` package with structured sections, citation checks, sensitive-field redaction, missing-evidence fail-closed behavior, uncertainty labels, and approval-state display.

## Current Limits

- No Markdown, PDF, email, download, CLI, HTTP API, database, persistence, export, escalation, external sharing, observability, or eval harness exists in Phase 6.
- The draft is structured data only; downstream rendering is a future phase or integration task.
- Redaction rules are deterministic MVP rules, not a configurable policy engine.
- `Shareable: true` means the draft is safe structured output for human review, not final approval or external release.
