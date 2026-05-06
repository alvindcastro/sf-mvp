# Review Composition Contract

Phase 12 adds the first deterministic demo composition boundary for Fleet Incident Copilot. It is implemented as package-level Go code under `internal/demo`; it is not an HTTP API, CLI, database-backed service, Slack integration, webhook, live model call, or external service.

## Phase 12 Checklist

- [x] Add machine-readable synthetic demo fixtures only after failing tests define fixture-loading expectations.
- [x] Define the smallest review response with validation status, retrieved citation refs, timeline entries, severity, recommendations, redacted brief, approval-required actions, and trace ID.
- [x] Reject non-synthetic or real-looking incident input before downstream composition.
- [x] Preserve existing citation, redaction, approval, eval, and observability package contracts.
- [x] Keep the composer in-memory and deterministic.
- [x] Document behavior only after tests prove it.

## Runtime Surface

- Package: `internal/demo`.
- Fixture file: `internal/demo/testdata/demo-fixtures.json`.
- Main entry points:
  - `LoadDefaultFixtures()`.
  - `LoadFixtures(data []byte)`.
  - `ComposeIncident(incidentID string, options Options)`.
  - `ComposePacket(packet ingestion.Packet, options Options)`.
- Targeted test command: `go test ./internal/demo`.
- Broader verification command: `go test ./...`.

## Fixture Contract

The machine-readable fixture file contains the five synthetic incidents from the human-readable packet contract:

| Incident ID | Fixture kind | Scenario |
| --- | --- | --- |
| `FIC-SYN-001` | `normal` | Low severity hard brake |
| `FIC-SYN-002` | `normal` | Medium severity stop-arm conflict |
| `FIC-SYN-003` | `normal` | High severity collision signal |
| `FIC-SYN-004` | `incomplete` | Unknown severity incomplete evidence |
| `FIC-SYN-005` | `adversarial` | Adversarial transcript with missing side view |

Fixture loading reuses `ingestion.IngestJSON`; it does not bypass packet validation. It rejects malformed fixture JSON, non-synthetic records, incident IDs without the `FIC-SYN-` prefix, missing media references, malformed packet fields, duplicate incident IDs, missing query text, and unsupported fixture kinds.

## Composition Flow

`ComposeIncident` finds a known synthetic fixture by incident ID. `ComposePacket` accepts a typed packet for package-level tests or future adapters, but it still rejects non-synthetic or real-looking input before retrieval, timeline, severity, or brief drafting runs.

For accepted synthetic input, the composer:

1. Starts an in-memory observability workflow and returns the deterministic trace ID.
2. Retrieves approved mock guidance for the fixture query.
3. Builds cited timeline entries from packet data.
4. Classifies severity and recommendations with deterministic rules.
5. Drafts the cited, redacted brief through `internal/brief`.
6. Uses the approval gate to confirm export, escalation, and external sharing remain blocked without scoped human approval.
7. Records package-level retrieval and composition tool-call observability events.

The composer fails closed when required evidence is missing. Missing evidence returns `ErrMissingEvidence` instead of a partial review result.

## Review Response

The returned `ReviewResult` includes:

- `ValidationStatus`: currently `accepted` for composed reviews.
- `IncidentID`.
- `TraceID`.
- `RetrievedCitationRefs`.
- `TimelineEntries` with claim text, source refs, timestamps, and uncertainty labels.
- `Severity` with level and sourced rationale.
- `Recommendations` with action, explanation, and source refs.
- `RedactedBrief` with draft status, shareable flag, redacted sections, approval state, redactions applied, and uncertainties.
- `ApprovalRequiredActions` for export, escalation, and external sharing, all blocked unless a later phase supplies scoped approval.
- `ObservabilityEvents` from the in-memory recorder.

## Current Limits

- No HTTP route or local server exists; Phase 13 owns the loopback API.
- No Slack payload, webhook, token, SDK, or network delivery exists; Phase 14 owns dry-run notification preview behavior.
- No approval retry demo exists; Phase 15 owns local approval retry wiring.
- No database, persistence, identity, authorization, live model provider, external telemetry backend, dashboard, or production audit store exists.
- The mock guidance corpus is in-memory demo data, not a vector index or customer knowledge base.

## Red-To-Green Evidence

- Added failing `internal/demo` tests for fixture loading, malformed and non-synthetic fixture rejection, happy-path review composition, unknown incident ID, non-synthetic typed input, missing evidence, citation preservation, redaction preservation, and approval-required action display.
- Observed `go test ./internal/demo` fail because `LoadDefaultFixtures`, `LoadFixtures`, `ComposeIncident`, `ComposePacket`, and response types did not exist.
- Implemented only the in-memory fixture loader and composer needed to pass.
- Verified with `go test ./internal/demo`, related package tests, and `go test ./...`.
