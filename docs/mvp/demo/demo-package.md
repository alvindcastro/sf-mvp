# Demo Package

Phase 10 packages the Fleet Incident Copilot MVP for a short review or interview demo. The materials must make the target operator workflow clear while staying honest about what is implemented as package-level Go code and what remains planned.

## Implemented Vs Planned Discipline

Use this wording rule across the repo narrative, video, diagram, and interview answers:

- Say **implemented** only for behavior backed by current docs, Go packages, and tests under `internal/ingestion`, `internal/retrieval`, `internal/timeline`, `internal/severity`, `internal/brief`, `internal/approval`, `internal/eval`, and `internal/observability`.
- Say **package-level** or **in-memory** when describing current runtime behavior unless a future phase adds a CLI, HTTP API, database, persistent store, external telemetry pipeline, or UI.
- Say **planned** for live model-provider calls, vector databases, hosted RAG services, real export tools, real escalation tools, external sharing, identity, roles, dashboards, alerts, billing reconciliation, production audit/compliance guarantees, and live fleet integrations.
- Say **agent-ready boundaries** or **constrained action gates** for the current approval and tool-call surfaces. Do not claim a live autonomous agent loop exists.
- Say **production-readiness case** for the demo's governance, eval, monitoring, security, and cost-control design. Do not claim the repository is production deployed.

## Repo Narrative

Fleet Incident Copilot is a synthetic fleet-safety incident review MVP for a fleet safety operator. The operator's job is to turn fragmented incident evidence into a grounded review: what happened, what evidence supports each claim, how severe the event appears, what actions are recommended, and which actions require human approval before anything leaves the review context.

The current repository demonstrates the backend core of that workflow with deterministic Go packages. A synthetic incident packet is validated, approved mock SOP guidance is retrieved with citations, a cited timeline is built, severity and recommended actions are derived, a redacted draft brief is produced, sensitive actions are blocked unless a human approval record matches scope, deterministic evals score golden cases, and in-memory observability records trace, retrieval, tool-call, approval, token, budget, and eval signals.

This is not a production fleet platform. There is no CLI, HTTP API, database, UI, real model call, vector database, telemetry backend, live export, live escalation, external-sharing integration, identity provider, or real customer evidence processing. The demo is valuable because it shows how a production-minded AI workflow would be grounded, constrained, measured, observed, and cost-controlled before connecting real integrations.

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
| Trust quality over time | `internal/eval` scores deterministic synthetic golden cases | CI reports, broader fixture sets, model regression tracking |
| Operate cost and reliability | `internal/observability` records in-memory events, caller-supplied token usage, budget failures, cache candidates, and routing notes | OpenTelemetry, dashboards, alerts, provider billing reconciliation |

## Demo Video Outline And Script

Target length: 3 to 5 minutes. The current demo should be framed as a docs, code, and tests walkthrough, not a live app demo.

| Time | Screen | Script |
| --- | --- | --- |
| 0:00-0:25 | `README.md` and `docs/mvp/README.md` | "This repo demonstrates a Fleet Incident Copilot for a fleet safety operator reviewing synthetic incident packets. The goal is not just summarization; it is grounded incident review with citations, approval gates, evals, monitoring, security, and cost controls." |
| 0:25-0:55 | `docs/mvp/overview/scope.md` | "The trust boundary is explicit: synthetic data only, no live fleet integrations, no autonomous export or escalation, and no production compliance claims." |
| 0:55-1:35 | `internal/ingestion`, `internal/retrieval`, `internal/timeline` tests | "The workflow starts by validating synthetic packets, retrieving scoped mock SOP guidance with citations, then building a cited timeline. Retrieved text is treated as data, including hostile fixture text." |
| 1:35-2:10 | `internal/severity` and `internal/brief` tests | "Severity and recommendations are deterministic and source-backed. The brief is a draft for human review, includes citations, redacts sensitive fields, and shows sensitive actions as blocked." |
| 2:10-2:45 | `internal/approval` tests | "Export, escalation, and external sharing are gated. Pending, denied, missing, and out-of-scope approvals fail closed. This is the current agent safety boundary." |
| 2:45-3:25 | `docs/mvp/quality/eval-plan.md` and `internal/eval` tests | "The repo includes deterministic golden eval cases for severity, citation coverage, recommendation accuracy, unsupported claims, redaction, prompt-injection resistance, and approval fail-closed behavior." |
| 3:25-4:10 | `docs/mvp/quality/observability-and-cost-controls.md` and `internal/observability` tests | "The observability package records in-memory trace, retrieval, tool-call, approval, latency, token, budget, and eval events. It also defines cache candidates and model-routing notes, but it does not call a live provider or reconcile billing." |
| 4:10-4:45 | This file | "The production-readiness case is the architecture discipline: grounded inputs, constrained actions, eval gates, redacted logs, budget limits, and clear known limits before external integrations are added." |

Optional close: "The next production step would be a thin API or CLI around these package contracts, followed by persistence, real observability export, model-provider integration, and a reviewer UI."

## Future Local Demo Surface

The recommended next demo improvement is a local, loopback-only walkthrough described in [Demo Surface Roadmap](demo-surface-roadmap.md). It is planned, not implemented.

The target hiring-manager arc is:

- [ ] Call a planned local review endpoint with a synthetic incident ID or packet.
- [ ] Show the composed review output: validation, citations, timeline, severity, recommendations, redacted brief, approval-required actions, and trace ID.
- [ ] Attempt a planned dry-run Slack-shaped notification preview and show it is blocked before scoped approval.
- [ ] Record a planned in-memory approval for the exact incident, action, and channel.
- [ ] Retry the dry-run notification preview and show only the approved dry-run payload is allowed.
- [ ] Show a planned local eval report and redacted trace report.

Until those phases are implemented with strict TDD, the current demo remains the docs, code, and tests walkthrough above.

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
- Human approval gate: implemented by `internal/approval`; label export, escalation, and external sharing as blocked unless approved in scope.
- Eval harness: implemented by `internal/eval` over deterministic synthetic golden cases.
- Observability and cost controls: implemented by `internal/observability` as in-memory events, token-budget checks, cache candidates, and routing notes.
- Planned API boundary: mark CLI, HTTP API, persistence, identity, real integrations, dashboards, and provider calls as planned, not implemented.

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

## One-Page Eval Summary Outline

Title: **Fleet Incident Copilot MVP Eval Summary**

Use this one-page structure:

- **Scope evaluated:** deterministic synthetic incident review path across ingestion, retrieval, timeline, severity, brief drafting, approval gating, eval scoring, and observability recording.
- **Fixture set:** five synthetic golden cases: hard brake, stop-arm conflict, collision signal, incomplete evidence, and adversarial transcript.
- **Metrics:** severity accuracy, citation coverage, recommendation accuracy, unsupported-claim absence, redaction leak absence, prompt-injection resistance, and approval fail-closed behavior.
- **Release thresholds:** default deterministic thresholds require perfect scores for severity, citation coverage, recommendation accuracy, and safety checks.
- **Current evidence source:** `internal/eval` package behavior and targeted Go tests; fill final numeric results from the latest test or eval run used for the demo recording.
- **Risk controls:** scoped retrieval, citations, unsupported-claim checks, redaction checks, approval gating, prompt-injection fixture, and strict budget behavior in observability.
- **Known limits:** no live LLM, no model regression suite, no external eval report, no persisted eval history, no production monitoring, and no real-world incident data.
- **Next improvements:** add CLI or API eval report, persist historical runs, add more adversarial fixtures, introduce model-provider eval comparisons, and wire eval summaries into CI and dashboards.

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

- Implemented: explicit Go package APIs for ingestion, retrieval, timeline building, severity, brief drafting, approval, evals, and observability.
- Planned: loopback-only hiring-manager demo API, CLI fallback, database, durable jobs, authentication, authorization, and integration adapters.
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
- Planned: transport layer, persistence, identity, integrations, operational telemetry, incident response, deployment hardening, and compliance review.
- Key point: the repo is not production-ready software; it is a production-readiness demonstration for a narrowly scoped applied-AI workflow.

## Final Packaging Checklist

- [x] Repo narrative maps the MVP to the fleet safety operator role.
- [x] Demo video outline and script are ready.
- [x] Architecture diagram checklist distinguishes implemented package surfaces from planned integrations.
- [x] One-page eval summary outline is ready.
- [x] Interview talking points cover RAG, agents, backend APIs, evals, monitoring, security, cost, and production readiness.
- [x] Implemented-vs-planned wording discipline is explicit.
