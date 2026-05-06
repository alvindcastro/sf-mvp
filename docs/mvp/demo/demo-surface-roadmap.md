# Demo Surface Roadmap

Phase 11 began as a documentation-only brainstorm for a concrete hiring-manager demo surface. Phase 13 implements the loopback review API, Phase 14 implements a dry-run Slack-shaped notification preview, and Phase 15 implements an in-memory scoped approval retry demo. Real Slack delivery, webhooks, database, persistent store, live model calls, and external integrations remain unimplemented.

## Current State

- [x] The implemented runtime is package-level Go code under `internal`.
- [x] The current proof command is `go test ./...`.
- [x] The current demo package is a docs, code, package-level composer, loopback API, and tests walkthrough.
- [x] The repo has in-memory approval, eval, and observability packages.
- [x] The repo has an in-memory demo review composer in `internal/demo`.
- [x] Machine-readable demo fixtures exist under `internal/demo/testdata`.
- [x] The repo has a loopback-only local review API in `internal/httpapi`.
- [x] The repo has a thin local server command in `cmd/demo-api`.
- [x] The repo has a dry-run Slack-shaped notification preview package in `internal/notification`.
- [x] The loopback demo API exposes `POST /demo/notifications/slack` for blocked and exact-approved dry-run previews.
- [x] The loopback demo API exposes `POST /demo/approvals` and `POST /demo/approvals/decisions` for in-memory scoped approval retry.
- [ ] No general CLI workflow exists yet.
- [ ] No real Slack integration exists yet.
- [ ] No webhook or external notification delivery exists yet.

## Recommended Hiring-Manager Demo Arc

Build the smallest local demo surface that proves the workflow without implying production integrations:

1. [x] `POST /demo/review` accepts a synthetic incident ID or synthetic packet JSON and returns the composed incident review.
2. [x] `POST /demo/notifications/slack` prepares a Slack-shaped notification preview in `dry_run` mode and returns `blocked` before scoped approval.
3. [x] `POST /demo/approvals` records an in-memory human approval for the exact synthetic incident, action, and channel target.
4. [x] Retrying `POST /demo/notifications/slack` returns an allowed `dry_run` payload without sending a network request.
5. [ ] `GET /demo/eval/latest` returns a local deterministic eval report with scores, thresholds, and pass/fail status.
6. [ ] `GET /demo/traces/{trace_id}` returns redacted in-memory observability events for the review, approval, notification preview, eval summary, and budget path.

This arc shows a concrete API call, an integration-shaped action, human approval gating, eval evidence, and observability proof while staying synthetic and local.

## Recommended Build Order

Original build order, updated as phases land:

- [x] Add machine-readable synthetic demo fixtures through strict TDD.
- [x] Add the review composer.
- [ ] Add a local eval report renderer over existing deterministic golden cases.
- [ ] Add demo-surface observability events for fixture load, review, approval retry, notification preview, eval report, and budget paths.
- [x] Add the loopback API.
- [x] Add scoped approval retry behavior if the current approval package needs a clearer retry path.
- [x] Add the dry-run Slack-shaped notification preview.
- [ ] Refresh the demo script with verified commands only after the surfaces exist.

## Verified Loopback Review Endpoint

The Phase 13 loopback review endpoint is implemented and locally verified. Port `8080` was occupied during verification, so this exact command uses the loopback override:

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

```json
{
  "trace_id": "trace-fic-syn-001-20260506t160000z-001",
  "review": {
    "validation_status": "accepted",
    "incident_id": "FIC-SYN-001",
    "severity": {
      "level": "low"
    },
    "redacted_brief": {
      "status": "draft"
    }
  },
  "approval_required_actions": [
    {
      "action": "export",
      "status": "blocked",
      "approved": false
    }
  ],
  "eval_summary": {
    "available": true,
    "ref": "docs/mvp/quality/eval-plan.md"
  }
}
```

## Verified Dry-Run Notification Preview Endpoint

The Phase 14 dry-run notification route is implemented and locally verified. It returns a prepared Slack-shaped payload with `status: "blocked"` before exact scoped approval exists.

```bash
go run ./cmd/demo-api -addr 127.0.0.1:18081
```

In a second terminal:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18081/demo/notifications/slack \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","channel":"#fleet-safety","delivery_mode":"dry_run"}'
```

Expected blocked response highlights:

```json
{
  "notification_preview": {
    "status": "blocked",
    "delivery_mode": "dry_run",
    "reason": "no approval exists within the requested scope",
    "prepared_payload": {
      "channel": "#fleet-safety"
    },
    "sent": false,
    "network_delivery_attempted": false
  }
}
```

## Verified Scoped Approval Retry

The Phase 15 approval retry routes are implemented and locally verified. A fresh server starts with no approval state; `approval-001` is deterministic only when this is the first approval request in that server process.

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18082/demo/approvals \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","action":"external_sharing","channel":"#fleet-safety","reason":"operator requested dry-run preview"}'
```

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18082/demo/approvals/decisions \
  -H "Content-Type: application/json" \
  -d '{"request_id":"approval-001","approver":"fleet-safety-lead","decision":"approved","reason":"redacted brief approved for #fleet-safety dry-run"}'
```

Retrying `POST /demo/notifications/slack` for `FIC-SYN-001`, `external_sharing`, and `#fleet-safety` returns:

```json
{
  "notification_preview": {
    "status": "allowed",
    "approval_request_id": "approval-001",
    "sent": false,
    "network_delivery_attempted": false
  }
}
```

## Surface Options

| Option | Demo value | Risk | Recommendation |
| --- | --- | --- | --- |
| Local loopback API | Shows a concrete `curl` path and backend composition | Could be mistaken for production API | Build first; label loopback-only and demo-only. |
| CLI command | Simple to run and test | Less impressive than an API call for integration roles | Keep as fallback if HTTP surface is too much for the next code run. |
| Dry-run Slack preview | Shows integration judgment and approval gating | Easy to overstate as real Slack delivery | Build as preview only; no token, secret, or network request. |
| Real Slack notification | Visually memorable | Adds secrets, network dependency, and external-data risk | Defer until explicit future scope allows live external services. |
| Local eval report | Shows quality gates beyond happy-path demo | Could look abstract without the API path | Add after review API, or expose through the API once composer exists. |
| Trace report endpoint | Shows observability, redaction, and budget thinking | Needs careful wording around monitoring | Add as in-memory proof, not dashboard or telemetry platform. |

## Phase Checklist

### Phase 11: Demo Surface Planning

- [x] Brainstorm API, CLI, Slack-shaped notification, approval, eval, and observability demo surfaces.
- [x] Choose local loopback API plus dry-run notification preview as the recommended path.
- [x] Keep real Slack delivery, webhooks, model calls, persistence, identity, and dashboards out of scope.
- [x] Add future phase checklists and prompts before any code is written.
- [x] Confirm this phase is Markdown-only.

### Phase 12: Review Composition Contract

- [x] Add or reuse machine-readable synthetic fixtures only through a failing test first.
- [x] Define a demo review response that composes ingestion, retrieval, timeline, severity, brief, approval state, and trace ID.
- [x] Reject non-synthetic or real-looking evidence before downstream composition.
- [x] Preserve citations and redactions from existing package contracts.
- [x] Record package-level observability events without adding persistent logs.

Implemented output: [Review Composition Contract](review-composition-contract.md), `internal/demo`, and `internal/demo/testdata/demo-fixtures.json`. Phase 12 itself remains package-level and in-memory; Phase 13 adds the local server route.

### Phase 13: Loopback Demo API

- [x] Add `POST /demo/review` behind a local demo server or equivalent local-only transport.
- [x] Return deterministic JSON suitable for a hiring-manager `curl` demo.
- [x] Cover malformed JSON, unknown incident ID, non-synthetic input, and unsupported method paths.
- [x] Keep the API stateless or in-memory unless a later phase explicitly adds persistence.
- [x] Document commands only after tests prove they work.

Implemented output: [Loopback Demo API](loopback-demo-api.md), `internal/httpapi`, and `cmd/demo-api`. The route remains loopback-only and stateless; no auth, persistence, Slack behavior, webhook, live model call, export, escalation, or external-sharing integration exists.

### Phase 14: Dry-Run Slack Preview

- [x] Add a Slack-shaped preview payload generated from the redacted brief.
- [x] Require `delivery_mode: "dry_run"` for every notification preview.
- [x] Block previewed external sharing until scoped approval exists.
- [x] Prove no network call, token, secret, webhook, or Slack SDK is used.
- [x] Record a redacted tool-call observability event for preview generation.

Implemented output: [Dry-Run Slack-Shaped Notification Preview](dry-run-slack-preview.md), `internal/notification`, and `POST /demo/notifications/slack`. The route returns a blocked dry-run preview before scoped approval and an allowed dry-run preview after the Phase 15 exact approval retry.

### Phase 15: Scoped Approval Retry

- [x] Demonstrate missing, pending, denied, and out-of-scope approvals fail closed.
- [x] Demonstrate an approved dry-run notification succeeds only for the exact incident, action, and target channel.
- [x] Preserve in-memory audit history for approval request and decision events.
- [x] Keep approval decisions human-supplied and deterministic.
- [x] Do not infer approval from model output, notification payload content, or test fixture names.

Implemented output: [Scoped Approval Demo Retry](scoped-approval-retry.md), `POST /demo/approvals`, `POST /demo/approvals/decisions`, and shared in-memory approval gate state inside `internal/httpapi`. The retry flow is local and ephemeral; no identity, persistence, Slack delivery, webhook, or real external-sharing integration exists.

### Phase 16: Eval And Observability Demo Reports

- [ ] Add a local eval report surface that runs deterministic golden cases and returns scores, thresholds, and pass/fail state.
- [ ] Add an in-memory trace report surface that returns redacted workflow events by trace ID.
- [ ] Include a budget-exceeded demo path using caller-supplied token counts.
- [ ] Keep reports local and ephemeral; do not imply dashboards, alerts, OpenTelemetry export, or persisted history.
- [ ] Update the one-page eval summary with numbers from the latest verified run.

### Phase 17: Demo Script Refresh

- [ ] Update the demo script from a code/tests walkthrough to a local API walkthrough after implementation exists.
- [ ] Add exact local run commands and `curl` examples only after they pass locally.
- [ ] Keep the old package-level walkthrough as a fallback demo.
- [ ] Add a short recording plan with one happy path, one blocked action, one approved dry-run retry, one eval report, and one trace report.
- [ ] Confirm all implemented-versus-planned language is synchronized across `README.md`, `docs/mvp/README.md`, and the demo package.

## Wording Guardrails

Use these phrases for current and planned surfaces:

- **implemented dry-run Slack-shaped notification preview**, not Slack integration.
- **implemented in-memory scoped approval demo**, not identity-backed approval workflow.
- **planned local eval report over deterministic synthetic cases**, not model benchmark.
- **planned in-memory observability proof**, not monitoring platform.

Use these phrases after strict-TDD implementation proves the behavior:

- **implemented loopback-only demo API** if tests and local commands prove the endpoint.
- **implemented dry-run Slack-shaped preview** because no network delivery exists.
- **implemented in-memory approval gate** for scoped action callbacks.
- **implemented local eval report** if the report is generated from `internal/eval`.
- **implemented in-memory trace report** if it returns package-level observability events.

Always keep these deferred unless future scope explicitly changes:

- [ ] Real Slack delivery.
- [ ] Webhook calls.
- [ ] Secrets or Slack tokens.
- [ ] Persistent approval store.
- [ ] Authentication and authorization.
- [ ] Live model-provider calls.
- [ ] Provider billing reconciliation.
- [ ] External observability export, dashboards, alerts, or log retention.
- [ ] Production audit, compliance, or customer-data claims.
