# Demo Package

Phase 10 packaged the Fleet Incident Copilot MVP for a short review or interview demo. Phase 17 refreshes that package around the verified loopback API walkthrough now that the local demo surface exists. The materials must make the target operator workflow clear while staying honest about what is implemented as package-level Go code and what remains planned.

## Implemented Vs Planned Discipline

Use this wording rule across the repo narrative, video, diagram, and interview answers:

- Say **implemented** only for behavior backed by current docs, Go packages, and tests under `internal/ingestion`, `internal/retrieval`, `internal/timeline`, `internal/severity`, `internal/brief`, `internal/approval`, `internal/eval`, `internal/observability`, `internal/demo`, `internal/notification`, `internal/httpapi`, and `cmd/demo-api`.
- Say **package-level**, **in-memory**, or **loopback-only local API** when describing current runtime behavior. Do not imply persistence, production deployment, identity, live integrations, or external delivery.
- Say **planned** for live model-provider calls, vector databases, hosted RAG services, real export tools, real escalation tools, external sharing, identity, roles, dashboards, alerts, billing reconciliation, production audit/compliance guarantees, production APIs, and live fleet integrations.
- Say **agent-ready boundaries** or **constrained action gates** for the current approval and tool-call surfaces. Do not claim a live autonomous agent loop exists.
- Say **production-readiness case** for the demo's governance, eval, monitoring, security, and cost-control design. Do not claim the repository is production deployed.

## Repo Narrative

Fleet Incident Copilot is a synthetic fleet-safety incident review MVP for a fleet safety operator. The operator's job is to turn fragmented incident evidence into a grounded review: what happened, what evidence supports each claim, how severe the event appears, what actions are recommended, and which actions require human approval before anything leaves the review context.

The current repository demonstrates the backend core of that workflow with deterministic Go packages. A synthetic incident packet is validated, approved mock SOP guidance is retrieved with citations, a cited timeline is built, severity and recommended actions are derived, a redacted draft brief is produced, sensitive actions are blocked unless a human approval record matches scope, deterministic evals score golden cases, in-memory observability records trace, retrieval, tool-call, approval, token, budget, and eval signals, and a dry-run Slack-shaped notification preview can be retried after an exact in-memory scoped approval without sending anything.

This is not a production fleet platform. It has a loopback-only local demo API, but no general CLI workflow, database, UI, real model call, vector database, telemetry backend, live export, live escalation, Slack delivery, external-sharing integration, identity provider, production API, or real customer evidence processing. The demo is valuable because it shows how a production-minded AI workflow would be grounded, constrained, measured, observed, and cost-controlled before connecting real integrations.

## Target Role Mapping

For a fleet safety operator, the MVP maps to these review responsibilities:

| Operator need | Current MVP support | Planned production surface |
| --- | --- | --- |
| Know whether an incident packet can be trusted | `internal/ingestion` validates synthetic packet structure and rejects non-synthetic evidence | Intake API, storage, identity, retention, and source-system checks |
| Find relevant policy context | `internal/retrieval` returns scoped mock SOP citations using deterministic lexical retrieval | Managed corpus, embeddings or hybrid search, corpus versioning, access control |
| Reconstruct what happened | `internal/timeline` builds cited timeline entries and marks uncertainty | Workflow API, UI timeline, persisted evidence graph |
| Triage severity and next steps | `internal/severity` classifies severity and recommends SOP-grounded actions | Reviewer-configurable policies, escalation queues, audit trails |
| Prepare a safe review artifact | `internal/brief` drafts cited, redacted, human-reviewable brief sections | Rendering, export, sharing, templates, review UI |
| Prevent unsafe automation | `internal/approval` blocks export, escalation, and external sharing without scoped approval | Roles, expiry, policy engine, real action tools |
| Preview integration-shaped actions safely | `internal/notification` prepares dry-run Slack-shaped payloads from redacted briefs and `internal/httpapi` retries them only after exact scoped approval | Notification adapters, channel policy, delivery audit, identity-backed approval flow |
| Trust quality over time | `internal/eval` scores deterministic synthetic golden cases | CI reports, broader fixture sets, model regression tracking |
| Operate cost and reliability | `internal/observability` records in-memory events, caller-supplied token usage, budget failures, cache candidates, and routing notes | OpenTelemetry, dashboards, alerts, provider billing reconciliation |

## Demo Video Outline And Script

Target length: 4 to 6 minutes. The current primary demo is a local loopback API walkthrough backed by package tests; it is not a deployed app.

| Time | Screen | Script |
| --- | --- | --- |
| 0:00-0:25 | `README.md` and `docs/mvp/README.md` | "This repo demonstrates a Fleet Incident Copilot for a fleet safety operator reviewing synthetic incident packets. The goal is not just summarization; it is grounded incident review with citations, approval gates, evals, observability, security, and cost controls." |
| 0:25-0:50 | `docs/mvp/overview/scope.md` | "The trust boundary is explicit: synthetic data only, no live fleet integrations, no autonomous export or escalation, no Slack delivery, and no production compliance claims." |
| 0:50-1:15 | Terminal startup | "I start the local demo server on loopback only. The server rejects non-loopback addresses and keeps approval and trace state in memory for this process." |
| 1:15-1:55 | `POST /demo/review` response | "The first call composes a deterministic review for a synthetic incident: accepted validation, cited timeline, low severity, redacted draft brief, blocked approval-required actions, and a trace ID." |
| 1:55-2:35 | `POST /demo/notifications/slack` before approval | "The integration-shaped action is a dry-run Slack-shaped preview. Before scoped approval, it returns a prepared payload but remains blocked, with `sent: false` and no network delivery attempted." |
| 2:35-3:20 | `POST /demo/approvals`, `POST /demo/approvals/decisions`, retry notification | "A human approval is recorded for the exact incident, `external_sharing` action, and Slack-shaped channel target. Retrying the same dry-run preview becomes allowed, but still sends nothing." |
| 3:20-4:05 | `GET /demo/eval/latest` | "The eval report runs deterministic golden cases and returns five cases, pass status, perfect severity accuracy, citation coverage, and recommendation accuracy under strict gates." |
| 4:05-4:45 | `GET /demo/traces/{trace_id}` and optional budget check | "Trace reports show redacted in-memory events from the same server process. The budget demo uses caller-supplied token counts and shows a budget-exceeded event without calling a model provider." |
| 4:45-5:20 | `go test ./...` fallback | "If the API server cannot be shown live, the fallback proof is the Go test suite across the same package boundaries. The tests are the source of truth for behavior." |

Optional close: "Future production work would add persistence, identity, live model-provider integration, real observability export, integration adapters, and a reviewer UI only after explicit scope and tests exist."

## Local Demo Surface

The current local walkthrough is described in [Loopback Demo API](loopback-demo-api.md), [Dry-Run Slack-Shaped Notification Preview](dry-run-slack-preview.md), [Scoped Approval Demo Retry](scoped-approval-retry.md), [Eval And Observability Demo Reports](eval-and-observability-reports.md), and [Demo Surface Roadmap](demo-surface-roadmap.md). The package-level review composer is implemented in `internal/demo`, dry-run notification preview is implemented in `internal/notification`, deterministic evals are implemented in `internal/eval`, observability and budget controls are implemented in `internal/observability`, and the loopback-only API plus in-memory approval retry and report routes are implemented in `internal/httpapi` with a thin `cmd/demo-api` server.

Use one running server process for the whole walkthrough. Approval state, audit history, and trace reports are in memory and are cleared when the process stops. The default server address is `127.0.0.1:8080`; the verified Phase 17 run used a loopback override because local ports can be occupied.

Start the server:

```bash
go run ./cmd/demo-api -addr 127.0.0.1:18084
```

Review one synthetic incident:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18084/demo/review \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001"}'
```

Expected highlights:

- HTTP `200`.
- `trace_id: "trace-fic-syn-001-20260506t160000z-001"`.
- `review.validation_status: "accepted"`.
- `review.severity.level: "low"`.
- `review.redacted_brief.status: "draft"`.
- three `approval_required_actions`, all `blocked` and unapproved.
- `eval_summary.ref: "docs/mvp/quality/eval-plan.md"`.

Show blocked dry-run notification before approval:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18084/demo/notifications/slack \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","channel":"#fleet-safety","delivery_mode":"dry_run"}'
```

Expected highlights:

- HTTP `200`.
- `notification_preview.status: "blocked"`.
- `notification_preview.reason: "no approval exists within the requested scope"`.
- `notification_preview.prepared_payload.channel: "#fleet-safety"`.
- `notification_preview.sent: false`.
- `notification_preview.network_delivery_attempted: false`.

Create the exact scoped approval request:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18084/demo/approvals \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","action":"external_sharing","channel":"#fleet-safety","reason":"operator requested dry-run preview"}'
```

Expected highlights on a fresh server:

- HTTP `201`.
- `approval_request.id: "approval-001"`.
- `approval_request.decision: "pending"`.
- `approval_request.target_ref: "slack:#fleet-safety"`.
- `audit_history` includes `approval.requested`.

Record the human decision:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18084/demo/approvals/decisions \
  -H "Content-Type: application/json" \
  -d '{"request_id":"approval-001","approver":"fleet-safety-lead","decision":"approved","reason":"redacted brief approved for #fleet-safety dry-run"}'
```

Expected highlights:

- HTTP `200`.
- `approval_request.decision: "approved"`.
- `approval_request.approver: "fleet-safety-lead"`.
- `audit_history` includes `approval.requested` then `approval.decided`.

Retry the same dry-run notification:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18084/demo/notifications/slack \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","channel":"#fleet-safety","delivery_mode":"dry_run"}'
```

Expected highlights:

- HTTP `200`.
- `notification_preview.status: "allowed"`.
- `notification_preview.approval_request_id: "approval-001"`.
- `notification_preview.sent: false`.
- `notification_preview.network_delivery_attempted: false`.
- `audit_history` includes `sensitive_action.allowed`.

Run the deterministic eval report:

```bash
curl -i --max-time 5 http://127.0.0.1:18084/demo/eval/latest
```

Expected highlights:

- HTTP `200`.
- `trace_id: "trace-fic-syn-eval-report-20260506t160000z-001"`.
- `eval_report.case_count: 5`.
- `eval_report.passed: true`.
- `eval_report.metrics.severity_accuracy: 1`.
- `eval_report.metrics.citation_coverage: 1`.
- `eval_report.metrics.recommendation_accuracy: 1`.

Fetch the eval trace from the same running process:

```bash
curl -i --max-time 5 http://127.0.0.1:18084/demo/traces/trace-fic-syn-eval-report-20260506t160000z-001
```

Expected highlights:

- HTTP `200`.
- `trace_report.ephemeral: true`.
- `trace_report.events` includes `workflow.started`.
- `trace_report.events` includes `eval.score_recorded`.

Optional budget check:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18084/demo/budget/check \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","provider":"hosted","model":"demo-review-model","input_tokens":90,"output_tokens":20,"max_total_tokens":100}'
```

Expected highlights:

- HTTP `200`.
- `budget_report.status: "budget_exceeded"`.
- `budget_report.reason: "total token budget exceeded"`.
- `budget_report.token_usage.total_tokens: 110`.

Fallback proof when a live API walkthrough is impractical:

```bash
go test ./...
```

The fallback walks the package tests that prove ingestion, retrieval, timeline, severity, brief drafting, approval gating, evals, observability, demo composition, dry-run notification preview, local HTTP routes, and loopback server wiring.

Recording plan:

1. Open `README.md`, then start `go run ./cmd/demo-api -addr 127.0.0.1:18084`.
2. Show `POST /demo/review` and call out validation, citations, severity, redaction, blocked actions, and trace ID.
3. Show `POST /demo/notifications/slack` blocked before approval.
4. Create and approve `approval-001`, then retry the exact same notification preview and show `allowed`, `sent: false`, and `network_delivery_attempted: false`.
5. Show `GET /demo/eval/latest`, then `GET /demo/traces/{trace_id}` from the same process.
6. End with `go test ./...` as the fallback verification path and the implemented-versus-planned boundary.

The target hiring-manager arc is:

- [x] Call the implemented local review endpoint with a synthetic incident ID or packet.
- [x] Show the composed review output at package level: validation, citations, timeline, severity, recommendations, redacted brief, approval-required actions, and trace ID.
- [x] Attempt the dry-run Slack-shaped notification preview and show it is blocked before scoped approval.
- [x] Record an in-memory approval for the exact incident, action, and channel.
- [x] Retry the dry-run notification preview and show only the approved dry-run payload is allowed.
- [x] Show a local eval report, redacted trace report, and caller-supplied budget-exceeded demo.
- [x] Keep `go test ./...` as the fallback proof path.

## Architecture Diagram Checklist

Diagram title: **Fleet Incident Copilot MVP: Synthetic Incident Review With Safety Gates**

Include these boxes and labels:

- Synthetic incident packet input: implemented as in-memory or JSON packet validation through `internal/ingestion`.
- Validation and audit boundary: implemented package-level accepted or rejected audit event.
- Mock SOP and troubleshooting corpus: implemented as supplied in-memory mock documents for the MVP corpus.
- Retrieval and citation layer: implemented by `internal/retrieval`; label as deterministic lexical retrieval, not vector search.
- Cited timeline builder: implemented by `internal/timeline`.
- Severity and recommended-action engine: implemented by `internal/severity`.
- Shareable brief drafting: implemented by `internal/brief`; label as structured draft data, not rendered export.
- Demo review composer: implemented by `internal/demo`; label as deterministic package-level composition.
- Loopback review API: implemented by `internal/httpapi` and `cmd/demo-api`; label as demo-only and local, not a production API.
- Dry-run notification preview: implemented by `internal/notification` and exposed through `internal/httpapi`; label as Slack-shaped payload only, not Slack delivery.
- Scoped approval retry: implemented by `internal/httpapi` with `internal/approval`; label as in-memory and exact-scope, not identity-backed authorization.
- Human approval gate: implemented by `internal/approval`; label export, escalation, and external sharing as blocked unless approved in scope.
- Eval harness: implemented by `internal/eval` over deterministic synthetic golden cases.
- Observability and cost controls: implemented by `internal/observability` as in-memory events, token-budget checks, cache candidates, and routing notes.
- Planned production boundary: mark general CLI workflows, persistence, identity, real integrations, dashboards, provider calls, and production API behavior as planned, not implemented.

Use edge labels:

- "validated packet"
- "scoped guidance query"
- "citation refs"
- "timeline entries with sources"
- "severity, rationale, recommended actions"
- "redacted draft sections"
- "approval request or blocked sensitive action"
- "eval report"
- "trace and budget events"

Show trust boundaries:

- Untrusted packet data before validation.
- Retrieved content treated as data, not instructions.
- Sensitive actions fail closed at the approval gate.
- Logs and shareable outputs redact sensitive fields.

## One-Page Eval Summary

Title: **Fleet Incident Copilot MVP Eval Summary**

- **Scope evaluated:** deterministic synthetic incident review path across ingestion, retrieval, timeline, severity, brief drafting, approval gating, eval scoring, and observability recording.
- **Fixture set:** five synthetic golden cases: hard brake, stop-arm conflict, collision signal, incomplete evidence, and adversarial transcript.
- **Metrics:** severity accuracy, citation coverage, recommendation accuracy, unsupported-claim absence, redaction leak absence, prompt-injection resistance, and approval fail-closed behavior.
- **Release thresholds:** default deterministic thresholds require perfect scores for severity, citation coverage, recommendation accuracy, and safety checks.
- **Current evidence source:** `internal/eval` package behavior, targeted Go tests, and `GET /demo/eval/latest`.
- **Latest local result:** verified on 2026-05-06 through `GET /demo/eval/latest` on `127.0.0.1:18084`: `case_count: 5`, `passed: true`, `severity_accuracy: 1`, `citation_coverage: 1`, and `recommendation_accuracy: 1`.
- **Risk controls:** scoped retrieval, citations, unsupported-claim checks, redaction checks, approval gating, prompt-injection fixture, and strict budget behavior in observability.
- **Known limits:** no live LLM, no model regression suite, no external eval platform, no persisted eval history, no production monitoring, and no real-world incident data.
- **Next improvements:** persist historical runs only if future scope allows it, add more adversarial fixtures, introduce model-provider eval comparisons, and wire eval summaries into CI and dashboards only as later production work.

## Interview Talking Points

### RAG

- Implemented: scoped retrieval over approved mock SOP and troubleshooting documents, stable citation references, no-match behavior, and hostile retrieved text preserved as untrusted data.
- Planned: embeddings, vector or hybrid search, corpus ingestion, access-controlled document stores, freshness checks, and production corpus governance.
- Key point: the repo treats RAG as a trust boundary, not just context stuffing.

### Agents

- Implemented: deterministic approval gates and observable tool-call records that model how sensitive actions should be constrained.
- Planned: real agent runner, model tool-calling loop, tool registry, retries, and real export or escalation tools.
- Key point: any future agent should inherit the current fail-closed approval semantics and scoped action validation.

### Backend APIs

- Implemented: explicit Go package APIs for ingestion, retrieval, timeline building, severity, brief drafting, approval, evals, observability, demo composition, dry-run notification preview, scoped approval retry, and the loopback-only demo API.
- Planned: production API, general CLI fallback, database, durable jobs, authentication, authorization, and integration adapters.
- Key point: the current code favors testable domain boundaries before service transport.

### Evals

- Implemented: deterministic golden-case evals for severity, citations, recommendations, unsupported claims, redaction, prompt injection, and approval fail-closed behavior.
- Planned: persisted reports, CI artifacts, model-based regression comparisons, larger fixture library, and reviewer-labeled datasets.
- Key point: quality is defined before model-provider integration.

### Monitoring

- Implemented: in-memory structured events for trace IDs, retrieval counts, source IDs, tool-call success, approval decisions, latency, caller-supplied token usage, budget failures, and eval summaries.
- Planned: OpenTelemetry, metrics backend, dashboards, alerting, log retention, and incident debugging workflows.
- Key point: observability fields are selected to explain safety, quality, latency, and cost.

### Security

- Implemented: synthetic-only scope, validation before downstream use, retrieved content as data, redaction in brief and observability surfaces, prompt-injection fixtures, and fail-closed sensitive actions.
- Planned: identity, roles, tenant isolation, secret management, audit retention, policy engine, security review, and compliance controls.
- Key point: the demo shows security posture through boundaries and tests, not unsupported compliance claims.

### Cost

- Implemented: caller-supplied token recording, token-budget limits, budget-exceeded behavior, invalid token rejection, cache-candidate definitions, and model-routing notes.
- Planned: provider billing reconciliation, live routing, prompt compression, cache storage, cost dashboards, and usage quotas.
- Key point: the system has a cost-control design before it has live provider spend.

### Production Readiness

- Implemented: clear scope, deterministic package contracts, strict-TDD evidence, approval gates, evals, redaction, observability, and cost-control planning.
- Planned: production transport layer, persistence, identity, integrations, operational telemetry, incident response, deployment hardening, and compliance review.
- Key point: the repo is not production-ready software; it is a production-readiness demonstration for a narrowly scoped applied-AI workflow.

## Final Packaging Checklist

- [x] Repo narrative maps the MVP to the fleet safety operator role.
- [x] Demo video outline and local API script are ready.
- [x] Verified local startup and `curl` commands are included.
- [x] Fallback `go test ./...` walkthrough is included.
- [x] Recording plan covers happy-path review, blocked dry-run notification, exact scoped approval retry, eval report, trace report, and optional budget check.
- [x] Architecture diagram checklist distinguishes implemented package surfaces from planned integrations.
- [x] One-page eval summary includes numbers from the latest verified local run.
- [x] Interview talking points cover RAG, agents, backend APIs, evals, monitoring, security, cost, and production readiness.
- [x] Implemented-vs-planned wording discipline is explicit.
