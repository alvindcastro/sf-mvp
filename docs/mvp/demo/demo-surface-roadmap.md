# Demo Surface Roadmap

Phase 11 is a documentation-only brainstorm for adding a concrete hiring-manager demo surface in a future code run. No API, CLI, Slack delivery, webhook, database, persistent store, live model call, or external integration is implemented by this document.

## Current State

- [x] The implemented runtime is package-level Go code under `internal`.
- [x] The current proof command is `go test ./...`.
- [x] The current demo package is a docs, code, package-level composer, and tests walkthrough.
- [x] The repo has in-memory approval, eval, and observability packages.
- [x] The repo has an in-memory demo review composer in `internal/demo`.
- [x] Machine-readable demo fixtures exist under `internal/demo/testdata`.
- [ ] No local HTTP API exists yet.
- [ ] No CLI exists yet.
- [ ] No Slack integration exists yet.
- [ ] No webhook or external notification delivery exists yet.

## Recommended Hiring-Manager Demo Arc

Build the smallest local demo surface that proves the workflow without implying production integrations:

1. [ ] `POST /demo/review` accepts a synthetic incident ID or synthetic packet JSON and returns the composed incident review.
2. [ ] `POST /demo/notifications/slack` prepares a Slack-shaped notification preview in `dry_run` mode and returns `blocked` before scoped approval.
3. [ ] `POST /demo/approvals` records an in-memory human approval for the exact synthetic incident, action, and channel target.
4. [ ] Retrying `POST /demo/notifications/slack` returns an allowed `dry_run` payload without sending a network request.
5. [ ] `GET /demo/eval/latest` returns a local deterministic eval report with scores, thresholds, and pass/fail status.
6. [ ] `GET /demo/traces/{trace_id}` returns redacted in-memory observability events for the review, approval, notification preview, eval summary, and budget path.

This arc shows a concrete API call, an integration-shaped action, human approval gating, eval evidence, and observability proof while staying synthetic and local.

## Recommended Build Order

Build the supporting proof surfaces before the external-looking adapter:

- [x] Add machine-readable synthetic demo fixtures through strict TDD.
- [x] Add the review composer.
- [ ] Add a local eval report renderer over existing deterministic golden cases.
- [ ] Add demo-surface observability events for fixture load, review, approval retry, notification preview, eval report, and budget paths.
- [ ] Add the loopback API.
- [ ] Add scoped approval retry behavior if the current approval package needs a clearer retry path.
- [ ] Add the dry-run Slack-shaped notification preview.
- [ ] Refresh the demo script with verified commands only after the surfaces exist.

## Planned Endpoint Sketch

These commands are planned examples only. They must not be documented as runnable until a future strict-TDD code phase implements and verifies them.

```bash
curl -s -X POST http://localhost:8080/demo/review \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001"}'
```

Expected response shape:

```json
{
  "incident_id": "FIC-SYN-001",
  "trace_id": "trace-demo-001",
  "severity": "medium",
  "timeline": [],
  "brief": {},
  "citations": [],
  "approval_required_actions": [
    {
      "action": "external_sharing",
      "target": "slack:#fleet-safety",
      "status": "blocked",
      "reason": "human approval required"
    }
  ]
}
```

```bash
curl -s -X POST http://localhost:8080/demo/notifications/slack \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","channel":"#fleet-safety","delivery_mode":"dry_run"}'
```

Expected blocked response shape before approval:

```json
{
  "status": "blocked",
  "delivery_mode": "dry_run",
  "reason": "external notification requires scoped human approval",
  "prepared_payload": {
    "channel": "#fleet-safety",
    "text": "Draft incident review ready for human approval: FIC-SYN-001"
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

Implemented output: [Review Composition Contract](review-composition-contract.md), `internal/demo`, and `internal/demo/testdata/demo-fixtures.json`. This remains package-level and in-memory; no local server, route, notification preview, persistence, or external integration exists yet.

### Phase 13: Loopback Demo API

- [ ] Add `POST /demo/review` behind a local demo server or equivalent local-only transport.
- [ ] Return deterministic JSON suitable for a hiring-manager `curl` demo.
- [ ] Cover malformed JSON, unknown incident ID, non-synthetic input, and unsupported method paths.
- [ ] Keep the API stateless or in-memory unless a later phase explicitly adds persistence.
- [ ] Document commands only after tests prove they work.

### Phase 14: Dry-Run Slack Preview

- [ ] Add a Slack-shaped preview payload generated from the redacted brief.
- [ ] Require `delivery_mode: "dry_run"` for every notification preview.
- [ ] Block previewed external sharing until scoped approval exists.
- [ ] Prove no network call, token, secret, webhook, or Slack SDK is used.
- [ ] Record a redacted tool-call observability event for preview generation.

### Phase 15: Scoped Approval Retry

- [ ] Demonstrate missing, pending, denied, and out-of-scope approvals fail closed.
- [ ] Demonstrate an approved dry-run notification succeeds only for the exact incident, action, and target channel.
- [ ] Preserve in-memory audit history for approval request and decision events.
- [ ] Keep approval decisions human-supplied and deterministic.
- [ ] Do not infer approval from model output, notification payload content, or test fixture names.

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

Use these phrases before implementation:

- **planned local demo API**, not production API.
- **planned dry-run Slack-shaped notification preview**, not Slack integration.
- **planned in-memory scoped approval demo**, not identity-backed approval workflow.
- **planned local eval report over deterministic synthetic cases**, not model benchmark.
- **planned in-memory observability proof**, not monitoring platform.

Use these phrases after future strict-TDD implementation proves the behavior:

- **implemented loopback-only demo API** if tests and local commands prove the endpoint.
- **implemented dry-run Slack-shaped preview** if no network delivery exists.
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
