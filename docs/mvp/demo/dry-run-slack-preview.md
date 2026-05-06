# Dry-Run Slack-Shaped Notification Preview

Phase 14 adds a local dry-run notification preview for Fleet Incident Copilot. It prepares a Slack-shaped payload from the redacted incident brief and proves external sharing fails closed unless scoped approval exists. It does not deliver to Slack.

## Phase 14 Checklist

- [x] Generate a Slack-shaped payload from the redacted brief only.
- [x] Require `delivery_mode: "dry_run"` for notification previews.
- [x] Block notification preview as external sharing unless a scoped approval exists.
- [x] Return blocked status, reason, and prepared payload when approval is missing.
- [x] Record a redacted tool-call observability event for preview generation.
- [x] Prove no Slack token, webhook URL, SDK, secret, or network request is used.

## Runtime Surface

- Package: [internal/notification](../../../internal/notification).
- Local route: `POST /demo/notifications/slack`.
- Handler package: [internal/httpapi](../../../internal/httpapi).
- Local server command: [cmd/demo-api](../../../cmd/demo-api).
- Targeted package test: `go test ./internal/notification`.
- Targeted handler test: `go test ./internal/httpapi`.
- Full test command: `go test ./...`.

## Request Contract

```json
{
  "incident_id": "FIC-SYN-001",
  "channel": "#fleet-safety",
  "delivery_mode": "dry_run"
}
```

The handler composes the known synthetic review through `internal/demo`, converts its redacted brief into the notification preview input, and uses the local in-memory approval gate managed by `internal/httpapi`. Without a matching approved request, the route returns a blocked preview before approval. Phase 15 documents the retry path in [Scoped Approval Demo Retry](scoped-approval-retry.md).

## Response Contract

Successful dry-run preview responses return HTTP `200` with:

- `trace_id`.
- `notification_preview.status`, either `blocked` before exact approval or `allowed` after exact scoped approval.
- `notification_preview.delivery_mode: "dry_run"`.
- `notification_preview.reason`.
- `notification_preview.prepared_payload.channel`.
- `notification_preview.prepared_payload.text`.
- `notification_preview.prepared_payload.blocks`.
- `notification_preview.sent: false`.
- `notification_preview.network_delivery_attempted: false`.
- redacted `notification_preview.observability_events`.

Non-dry-run delivery modes fail with HTTP `400` and error code `dry_run_required`.

## Verified Local Commands

Port `18081` was used during local verification to avoid conflicting with other local demo runs:

```bash
go run ./cmd/demo-api -addr 127.0.0.1:18081
```

In a second terminal:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18081/demo/notifications/slack \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","channel":"#fleet-safety","delivery_mode":"dry_run"}'
```

Expected response highlights:

- HTTP `200`.
- `trace_id: "trace-fic-syn-001-20260506t160000z-001"`.
- `notification_preview.status: "blocked"`.
- `notification_preview.reason: "no approval exists within the requested scope"`.
- `notification_preview.prepared_payload.channel: "#fleet-safety"`.
- `notification_preview.sent: false`.
- `notification_preview.network_delivery_attempted: false`.
- `notification_preview.observability_events` includes `workflow.started` and `tool_call.completed`.

## Approval Behavior

`internal/notification` uses the existing approval gate with `external_sharing` and a target ref of `slack:<channel>`.

Package-level tests prove:

- Missing approval blocks preview and still returns the prepared payload.
- Denied approval blocks preview.
- Out-of-scope approval blocks preview.
- Exact scoped approval allows the dry-run preview.
- Allowed dry-run previews still do not send, call a network sender, or mark network delivery attempted.

The HTTP route does not infer approval from notification text or fixture names. Approval requests and human decisions are created through `POST /demo/approvals` and `POST /demo/approvals/decisions`; retry details live in [Scoped Approval Demo Retry](scoped-approval-retry.md).

## No-Delivery Boundary

Phase 14 adds no Slack SDK, token, webhook URL, environment secret, outgoing HTTP client, webhook sender, network callback, or real external sharing behavior. The payload is Slack-shaped JSON only.

## Red-To-Green Evidence

- Added failing `internal/notification` tests for missing approval, mandatory dry-run mode, denied approval, out-of-scope approval, scoped approval, redacted observability events, and no network delivery.
- Observed `go test ./internal/notification` fail because `PreviewSlack`, request, result, payload, status, and delivery-mode types did not exist.
- Implemented only the dry-run preview package needed to pass.
- Added failing `internal/httpapi` tests for `POST /demo/notifications/slack` and non-dry-run rejection.
- Observed `go test ./internal/httpapi` fail with `404` for the missing route.
- Implemented the local route as a thin wrapper over the demo composer and notification package.
- Verified with targeted package tests, `go test ./...`, `go vet ./...`, coverage, and one local `curl` command.
