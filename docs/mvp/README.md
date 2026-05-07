# Fleet Incident Copilot MVP

## Thesis

Build a Fleet Incident Copilot that demonstrates senior applied-AI engineering for fleet-safety workflows: synthetic onboard evidence, telemetry, policy retrieval, incident summarization, workflow automation, approval gates, evals, observability, security, and cost controls.

The MVP should show production judgment, not just LLM usage. It should make it easy to explain how the system is grounded, measured, constrained, and operated.

## Demo Promise

Given a synthetic fleet incident packet, the system should:

- [x] Validate and ingest structured synthetic event metadata, GPS or speed samples, transcript notes, and still-frame notes.
- [x] Retrieve relevant mock SOP or troubleshooting guidance.
- [x] Build a cited incident timeline.
- [x] Classify severity with rationale.
- [x] Recommend next actions grounded in retrieved guidance.
- [x] Draft a shareable incident brief with citations and redactions.
- [x] Require human approval before export, escalation, or external sharing.
- [x] Run local eval checks for severity, citation coverage, recommendation accuracy, unsupported claims, redaction, prompt-injection resistance, and approval fail-closed behavior.
- [x] Emit package-level observability signals for traces, retrieval quality, latency, token use, tool calls, approval decisions, eval summaries, and budget checks.
- [x] Compose one deterministic in-memory demo review result with validation status, retrieved citations, timeline entries, severity, recommendations, redacted brief, approval-required actions, and trace ID.
- [x] Expose the deterministic demo review result through a loopback-only local API.
- [x] Prepare a dry-run Slack-shaped notification preview from the redacted brief while blocking external sharing before scoped approval.
- [x] Retry the dry-run notification preview after an exact in-memory human approval while missing, pending, denied, out-of-scope, and wrong-action approvals stay blocked.
- [x] Expose local eval, trace, and caller-supplied budget report routes while keeping report state in memory and ephemeral.

This is the target promise, not complete end-to-end runtime behavior. The current repository state has Phase 0 and Phase 1 planning artifacts, the Phase 2 Go ingestion validator, the Phase 3 Go retrieval package for approved mock guidance, the Phase 4 Go timeline package for cited synthetic incident timelines, the Phase 5 Go severity package for deterministic classification, SOP-grounded recommendations, and approval-required flags, the Phase 6 Go brief package for cited, redacted draft incident briefs, the Phase 7 Go approval package for in-memory human approval records and scoped sensitive-action gating, the Phase 8 Go eval package for deterministic local golden-case scoring, the Phase 9 Go observability package for in-memory structured events, caller-supplied token usage, invalid token usage, budget-limit failures, cache candidates, and model-routing notes, the Phase 10 Markdown demo package, the Phase 11 Markdown roadmap for future local demo surfaces, the Phase 12 Go demo composer plus machine-readable synthetic fixtures, the Phase 13 loopback-only demo API, the Phase 14 dry-run Slack-shaped notification preview, the Phase 15 in-memory scoped approval retry demo, the Phase 16 local eval, trace, and budget report routes, and the Phase 17 refreshed local API demo script with verified commands and `go test ./...` fallback. No database, external observability pipeline, persistent log store, real model-provider call, provider billing reconciliation, real export tool, real escalation tool, Slack delivery, webhook, identity, auth, production API, persistent approval store, dashboard, alerting, OpenTelemetry export, model benchmark, or external-sharing integration exists yet.

## Primary User

A fleet safety operator reviewing an incident packet after a school bus, transit, law-enforcement, or waste-fleet safety event.

## What This Proves

- [ ] RAG design over proprietary-style operational content.
- [ ] Agent/tool orchestration with constrained actions.
- [ ] Backend-oriented product thinking.
- [x] Approval boundaries for sensitive workflows.
- [ ] Prompt-injection and least-privilege security awareness.
- [x] Eval discipline for groundedness, citations, severity, recommendations, redaction, prompt injection, and approval fail-closed behavior.
- [x] Observability and cost-control thinking.
- [x] Clear demo packaging for interviews and review.

## Artifact Map

Product and scope:

- [Product Frame](overview/product-frame.md): Phase 0 product promise, primary user, approval gates, non-goals, demo narrative, and success criteria.
- [Scope And Guardrails](overview/scope.md): in-scope, out-of-scope, trust boundaries, and demo path.

Workflow behavior:

- [Synthetic Incident Packets](workflow/synthetic-incident-packets.md): Phase 1 synthetic evidence records and workflow-output contract.
- [Incident Packet Ingestion](workflow/incident-packet-ingestion.md): Phase 2 packet schema, validation rules, audit events, test commands, and red-to-green evidence.
- [RAG Corpus And Grounding](workflow/rag-corpus-and-grounding.md): Phase 3 mock guidance corpus, citation rules, scope filtering, no-match behavior, prompt-injection fixture, eval questions, and retrieval test evidence.
- [Incident Timeline Builder](workflow/incident-timeline-builder.md): Phase 4 timeline output format, source-reference rules, uncertainty behavior, unsupported-claim behavior, test commands, and red-to-green evidence.
- [Severity Classification And Recommended Actions](workflow/severity-classification-and-recommended-actions.md): Phase 5 deterministic severity rules, recommendation output, approval-required flags, test commands, and red-to-green evidence.
- [Shareable Incident Brief Drafting](workflow/incident-brief-drafting.md): Phase 6 structured draft brief sections, citation rules, redaction behavior, fail-closed behavior, approval-state display, test commands, and red-to-green evidence.
- [Human Approval Workflow](workflow/human-approval-workflow.md): Phase 7 approval request model, decision capture, scoped enforcement rules, append-only audit behavior, test commands, and red-to-green evidence.
- [Review Composition Contract](demo/review-composition-contract.md): Phase 12 package-level demo composer, fixture loading, review response shape, safety boundaries, test commands, and red-to-green evidence.
- [Loopback Demo API](demo/loopback-demo-api.md): Phase 13 local handler and server command, request and response shape, error mapping, run commands, current limits, and red-to-green evidence.
- [Dry-Run Slack-Shaped Notification Preview](demo/dry-run-slack-preview.md): Phase 14 dry-run preview package and route, approval blocking behavior, no-delivery boundary, run commands, and red-to-green evidence.
- [Scoped Approval Demo Retry](demo/scoped-approval-retry.md): Phase 15 local approval request and decision routes, exact-scope retry behavior, in-memory audit history, run commands, and red-to-green evidence.
- [Eval And Observability Demo Reports](demo/eval-and-observability-reports.md): Phase 16 local eval report, redacted trace report, caller-supplied budget demo, run commands, current limits, and red-to-green evidence.

Quality and operations:

- [Eval Plan](quality/eval-plan.md): Phase 8 local eval harness behavior, golden cases, scoring rules, release thresholds, test commands, and red-to-green evidence.
- [Observability And Cost Controls](quality/observability-and-cost-controls.md): Phase 9 trace IDs, structured events, latency, token usage, invalid token usage, budget limits, redaction behavior, cache candidates, model-routing notes, test commands, and red-to-green evidence.

Execution and packaging:

- [Phases And Tasks](execution/phases.md): tickable phase plan with prompts.
- [Task Prompts](execution/task-prompts.md): reusable prompts for future implementation or documentation work.
- [Strict TDD Rules](execution/tdd-rules.md): non-negotiable rules for code tasks.
- [Demo Package](demo/demo-package.md): Phase 10 repo narrative and Phase 17 refreshed local API walkthrough, including verified startup, `curl` commands, expected response highlights, fallback `go test ./...`, recording plan, one-page eval summary, interview talking points, and implemented-versus-planned wording rules.
- [Demo Surface Roadmap](demo/demo-surface-roadmap.md): Phase 11 roadmap, now updated with the implemented Phase 13 loopback review API, implemented Phase 14 dry-run notification preview, implemented Phase 15 scoped approval retry, implemented Phase 16 eval and trace reports, and implemented Phase 17 demo-script refresh.
- [LinkedIn Post Drafts](demo/linkedin-post-drafts.md): casual repo posts aimed at hiring managers, CTOs, and engineering leaders, with AI framed as an assistive tool and the repo framed around concrete engineering judgment.
