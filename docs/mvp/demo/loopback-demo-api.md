# Loopback Demo API

Phase 13 added the first runnable local transport for Fleet Incident Copilot. Phase 14 extends the same loopback server with a dry-run Slack-shaped notification preview. Phase 15 adds in-memory scoped approval retry routes for the dry-run preview. The API is loopback-only, deterministic, and backed by the existing `internal/demo` review composer, `internal/notification`, and `internal/approval`. It is a hiring-manager demo surface, not a production API.

## Phase 13 Checklist

- [x] Add `POST /demo/review` for synthetic incident ID or synthetic packet JSON input.
- [x] Return deterministic JSON with review output, approval-required actions, eval summary pointer, and trace ID.
- [x] Reject malformed JSON, unknown incident IDs, non-synthetic input, unsupported methods, and unsupported paths.
- [x] Keep the API loopback-only and stateless or in-memory.
- [x] Do not add auth, database, identity, live model calls, or external integrations.
- [x] Add exact run and `curl` commands only after tests and local verification pass.

## Runtime Surface

- Handler package: [internal/httpapi](../../../internal/httpapi).
- Local server command: [cmd/demo-api](../../../cmd/demo-api).
- Routes: `POST /demo/review`, `POST /demo/approvals`, `POST /demo/approvals/decisions`, and `POST /demo/notifications/slack`.
- Default listen address: `127.0.0.1:8080`.
- Loopback override: `-addr 127.0.0.1:<port>`.
- Targeted handler test: `go test ./internal/httpapi`.
- Server wiring test: `go test ./cmd/demo-api`.
- Full test command: `go test ./...`.

The server rejects non-loopback address overrides. If port `8080` is already in use, start the demo on another loopback port. Approval state is in memory and lasts only for the current server process.

## Verified Local Commands

Port `8080` was occupied during local verification, so this exact command was verified with the loopback override:

```bash
go run ./cmd/demo-api -addr 127.0.0.1:18080
```

In a second terminal:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18080/demo/review \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001"}'
```

Expected response highlights:

- HTTP `200`.
- Top-level `trace_id`, such as `trace-fic-syn-001-20260506t160000z-001`.
- `review.validation_status: "accepted"`.
- `review.incident_id: "FIC-SYN-001"`.
- `review.severity.level: "low"`.
- `review.redacted_brief.status: "draft"`.
- `approval_required_actions` for `export`, `escalation`, and `external_sharing`, all `blocked` and not approved.
- `eval_summary.ref: "docs/mvp/quality/eval-plan.md"`.

If `127.0.0.1:8080` is free, the default startup command is:

```bash
go run ./cmd/demo-api
```

Phase 14 verified the dry-run notification preview route with this command:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18081/demo/notifications/slack \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","channel":"#fleet-safety","delivery_mode":"dry_run"}'
```

Expected response highlights:

- HTTP `200`.
- `notification_preview.status: "blocked"` before approval.
- `notification_preview.delivery_mode: "dry_run"`.
- `notification_preview.prepared_payload.channel: "#fleet-safety"`.
- `notification_preview.sent: false`.
- `notification_preview.network_delivery_attempted: false`.

Phase 15 verified the scoped approval retry route with this fresh-server sequence:

```bash
go run ./cmd/demo-api -addr 127.0.0.1:18082
```

In a second terminal:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18082/demo/approvals \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","action":"external_sharing","channel":"#fleet-safety","reason":"operator requested dry-run preview"}'
```

Then:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18082/demo/approvals/decisions \
  -H "Content-Type: application/json" \
  -d '{"request_id":"approval-001","approver":"fleet-safety-lead","decision":"approved","reason":"redacted brief approved for #fleet-safety dry-run"}'
```

Retrying the same dry-run notification request returns:

- HTTP `200`.
- `notification_preview.status: "allowed"`.
- `notification_preview.approval_request_id: "approval-001"`.
- `notification_preview.sent: false`.
- `notification_preview.network_delivery_attempted: false`.

## Request Contract

Known synthetic fixture by ID:

```json
{
  "incident_id": "FIC-SYN-001"
}
```

Synthetic packet JSON input uses the same packet shape accepted by `internal/ingestion`. The body must include `synthetic_record: true`, an incident ID with the `FIC-SYN-` prefix, telemetry samples, synthetic media refs, and the required packet fields.

The handler uses `ingestion.IngestJSON` for packet requests and does not duplicate packet validation rules.

## Response Contract

Successful responses include:

- `trace_id`.
- `review.validation_status`.
- `review.incident_id`.
- `review.retrieved_citation_refs`.
- `review.timeline_entries`.
- `review.severity`.
- `review.recommendations`.
- `review.redacted_brief`.
- `review.observability_events`.
- `approval_required_actions`.
- `eval_summary`, currently a pointer to the deterministic local eval documentation and command.

`POST /demo/notifications/slack` successful dry-run preview responses include:

- `trace_id`.
- `notification_preview.status`.
- `notification_preview.delivery_mode`.
- `notification_preview.reason`.
- `notification_preview.prepared_payload`.
- `notification_preview.sent`.
- `notification_preview.network_delivery_attempted`.
- `notification_preview.observability_events`.
- `audit_history`.

`POST /demo/approvals` successful responses include:

- `approval_request.id`.
- `approval_request.incident_id`.
- `approval_request.action`.
- `approval_request.target_ref`.
- `approval_request.decision`.
- `approval_request.reason`.
- `audit_history`.

`POST /demo/approvals/decisions` successful responses include the updated `approval_request`, including `approver` and `decision_reason`, plus append-only `audit_history`.

Error responses use the shape:

```json
{
  "error": {
    "code": "non_synthetic_input",
    "message": "only synthetic FIC-SYN incident input is accepted"
  }
}
```

Status mappings:

| Condition | Status | Code |
| --- | --- | --- |
| Malformed JSON | `400` | `malformed_json` |
| Empty request | `400` | `empty_request` |
| Invalid request shape | `400` | `invalid_request` |
| Unknown path | `404` | `not_found` |
| Unknown incident ID | `404` | `incident_not_found` |
| Unsupported method | `405` | `method_not_allowed` |
| Non-synthetic packet or incident ID | `422` | `non_synthetic_input` |
| Required evidence missing | `422` | `missing_evidence` |
| Non-dry-run notification mode | `400` | `dry_run_required` |
| Invalid approval request or decision | `400` | `invalid_approval_request` |
| Unknown approval request ID | `404` | `approval_request_not_found` |
| Approval decision rewrite | `409` | `approval_already_decided` |

## Current Limits

- The API is loopback-only and demo-only.
- State remains in memory; no database or persistence is added.
- No authentication, authorization, identity, roles, tenants, or production access control exists.
- The Slack-shaped payload is dry-run only. No Slack webhook, token, SDK, network delivery, or real external sharing exists.
- The approval retry demo is in-memory only. It is not an identity-backed workflow or persistent approval store.
- No local eval report or trace lookup endpoint exists; Phase 16 owns those report surfaces.
- No live model provider, vector database, hosted RAG service, real export, real escalation, external sharing, dashboard, alerting, or production audit store exists.

## Red-To-Green Evidence

- Added failing `internal/httpapi` tests for valid synthetic incident ID, valid synthetic packet JSON, malformed JSON, unknown incident ID, non-synthetic packet input, unsupported method, unknown path, and loopback default address.
- Observed `go test ./internal/httpapi` fail because `NewHandler` and `DefaultListenAddress` did not exist.
- Implemented the smallest handler that calls `internal/demo` and maps known demo errors to deterministic JSON status responses.
- Added failing `cmd/demo-api` tests for default loopback binding, loopback override, and non-loopback rejection.
- Observed `go test ./cmd/demo-api` fail because `newServer` did not exist, then implemented the thin local server wiring.
- Verified with `go test ./internal/httpapi`, `go test ./cmd/demo-api`, related package tests, `go test ./...`, `go vet ./...`, and one local `curl` command.
