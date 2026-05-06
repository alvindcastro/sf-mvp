# Scoped Approval Demo Retry

Phase 15 makes the approval gate visible through the loopback demo API. It adds in-memory approval request and decision routes, then reuses the same approval gate when retrying the dry-run Slack-shaped notification preview.

This is a local demo retry path, not identity-backed authorization, persistence, Slack delivery, or real external sharing.

## Phase 15 Checklist

- [x] Add a local approval request path for one synthetic incident, action, and target channel.
- [x] Show missing, pending, denied, and out-of-scope approvals fail closed.
- [x] Show approved dry-run notification preview succeeds only for the exact incident, action, and channel.
- [x] Preserve append-only in-memory audit history.
- [x] Keep approvals human-supplied and deterministic.
- [x] Do not infer approval from model output, notification text, fixture names, or test setup shortcuts.

## Runtime Surface

- Handler package: [internal/httpapi](../../../internal/httpapi).
- Approval gate package: [internal/approval](../../../internal/approval).
- Notification preview package: [internal/notification](../../../internal/notification).
- Local server command: [cmd/demo-api](../../../cmd/demo-api).
- Approval request route: `POST /demo/approvals`.
- Approval decision route: `POST /demo/approvals/decisions`.
- Retry route: `POST /demo/notifications/slack`.
- Targeted handler test: `go test ./internal/httpapi`.
- Related package tests: `go test ./internal/approval ./internal/notification ./internal/httpapi ./cmd/demo-api`.
- Full test command: `go test ./...`.

`internal/httpapi` keeps approval state in memory for the lifetime of one `cmd/demo-api` process. Restarting the local server clears approval requests and audit history.

## Request Contracts

Create a pending approval request:

```json
{
  "incident_id": "FIC-SYN-001",
  "action": "external_sharing",
  "channel": "#fleet-safety",
  "reason": "operator requested dry-run preview"
}
```

Record a human decision:

```json
{
  "request_id": "approval-001",
  "approver": "fleet-safety-lead",
  "decision": "approved",
  "reason": "redacted brief approved for #fleet-safety dry-run"
}
```

Retry the dry-run notification preview:

```json
{
  "incident_id": "FIC-SYN-001",
  "channel": "#fleet-safety",
  "delivery_mode": "dry_run"
}
```

The notification retry only matches an approved request when all of these values line up:

- `incident_id`.
- `action: "external_sharing"`.
- `target_ref: "slack:<channel>"`.

Approvals for a different incident, action, or channel remain blocked.

## Response Contract

`POST /demo/approvals` returns HTTP `201` with:

- `approval_request.id`.
- `approval_request.incident_id`.
- `approval_request.action`.
- `approval_request.target_ref`.
- `approval_request.decision`, initially `pending`.
- `approval_request.reason`.
- `audit_history`.

`POST /demo/approvals/decisions` returns HTTP `200` with the updated `approval_request` and append-only `audit_history`.

Dry-run notification retry responses continue to return:

- `notification_preview.status`, either `blocked` or `allowed`.
- `notification_preview.approval_request_id` when a matching request exists.
- `notification_preview.sent: false`.
- `notification_preview.network_delivery_attempted: false`.
- `audit_history`.

Allowed dry-run previews still do not send a Slack message or perform any network request.

## Verified Local Commands

Port `18082` was used during local verification to avoid conflicting with other local demo runs:

```bash
go run ./cmd/demo-api -addr 127.0.0.1:18082
```

In a second terminal, first show the fail-closed preview before approval:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18082/demo/notifications/slack \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","channel":"#fleet-safety","delivery_mode":"dry_run"}'
```

Expected response highlights:

- HTTP `200`.
- `notification_preview.status: "blocked"`.
- `notification_preview.reason: "no approval exists within the requested scope"`.
- `notification_preview.sent: false`.
- `notification_preview.network_delivery_attempted: false`.

Create a pending approval request:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18082/demo/approvals \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","action":"external_sharing","channel":"#fleet-safety","reason":"operator requested dry-run preview"}'
```

Expected response highlights:

- HTTP `201`.
- `approval_request.id: "approval-001"` on a fresh server process.
- `approval_request.decision: "pending"`.
- `approval_request.target_ref: "slack:#fleet-safety"`.
- `audit_history` includes `approval.requested`.

Record the approval decision:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18082/demo/approvals/decisions \
  -H "Content-Type: application/json" \
  -d '{"request_id":"approval-001","approver":"fleet-safety-lead","decision":"approved","reason":"redacted brief approved for #fleet-safety dry-run"}'
```

Expected response highlights:

- HTTP `200`.
- `approval_request.decision: "approved"`.
- `approval_request.approver: "fleet-safety-lead"`.
- `audit_history` includes `approval.requested` then `approval.decided`.

Retry the same dry-run notification preview:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18082/demo/notifications/slack \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","channel":"#fleet-safety","delivery_mode":"dry_run"}'
```

Expected response highlights:

- HTTP `200`.
- `notification_preview.status: "allowed"`.
- `notification_preview.approval_request_id: "approval-001"`.
- `notification_preview.sent: false`.
- `notification_preview.network_delivery_attempted: false`.
- `audit_history` includes `sensitive_action.allowed`.

## Fail-Closed Behavior

Handler tests cover these retry states:

- Missing approval blocks dry-run preview.
- Pending approval blocks dry-run preview.
- Denied approval blocks dry-run preview.
- Approved approval for another channel blocks dry-run preview.
- Approved approval for another incident blocks dry-run preview.
- Approved approval for another action, such as `export`, does not unlock notification preview.
- Exact scoped approval allows only the dry-run preview for that incident, action, and channel.

## No-Delivery Boundary

Phase 15 does not add Slack delivery, webhook calls, Slack tokens, environment secrets, external-sharing integrations, persistence, auth, identity, roles, dashboards, or production audit storage. It only makes the existing in-memory approval gate visible in the local loopback demo.

## Red-To-Green Evidence

- Added failing `internal/httpapi` tests for creating approval requests, pending retry, denied retry, exact approved retry, out-of-channel retry, out-of-incident retry, wrong-action approval, response audit history, and no network delivery.
- Observed `go test ./internal/httpapi` fail with `404` for missing approval routes.
- Implemented the smallest shared in-memory approval gate in `internal/httpapi`, protected by a mutex, and reused it for `POST /demo/notifications/slack`.
- Verified with targeted package tests, full Go tests, vet, coverage, and local `curl` commands.
