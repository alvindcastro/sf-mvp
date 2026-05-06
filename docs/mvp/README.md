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
- [ ] Require human approval before export, escalation, or external sharing.
- [ ] Emit observability and eval signals for traces, retrieval quality, latency, token use, tool calls, and safety checks.

This is the target promise, not complete end-to-end runtime behavior. The current repository state has Phase 0 and Phase 1 planning artifacts, the Phase 2 Go ingestion validator, the Phase 3 Go retrieval package for approved mock guidance, the Phase 4 Go timeline package for cited synthetic incident timelines, the Phase 5 Go severity package for deterministic classification, SOP-grounded recommendations, and approval-required flags, and the Phase 6 Go brief package for cited, redacted draft incident briefs. No CLI, HTTP API, database, approval workflow, eval harness, export, escalation, or external-sharing runtime exists yet.

## Primary User

A fleet safety operator reviewing an incident packet after a school bus, transit, law-enforcement, or waste-fleet safety event.

## What This Proves

- [ ] RAG design over proprietary-style operational content.
- [ ] Agent/tool orchestration with constrained actions.
- [ ] Backend-oriented product thinking.
- [ ] Approval boundaries for sensitive workflows.
- [ ] Prompt-injection and least-privilege security awareness.
- [ ] Eval discipline for groundedness, citations, severity, redaction, and tool calls.
- [ ] Observability and cost-control thinking.
- [ ] Clear demo packaging for interviews and review.

## Artifact Map

- [Product Frame](product-frame.md): Phase 0 product promise, primary user, approval gates, non-goals, demo narrative, and success criteria.
- [Synthetic Incident Packets](synthetic-incident-packets.md): Phase 1 synthetic evidence records and workflow-output contract.
- [Incident Packet Ingestion](incident-packet-ingestion.md): Phase 2 packet schema, validation rules, audit events, test commands, and red-to-green evidence.
- [RAG Corpus And Grounding](rag-corpus-and-grounding.md): Phase 3 mock guidance corpus, citation rules, scope filtering, no-match behavior, prompt-injection fixture, eval questions, and retrieval test evidence.
- [Incident Timeline Builder](incident-timeline-builder.md): Phase 4 timeline output format, source-reference rules, uncertainty behavior, unsupported-claim behavior, test commands, and red-to-green evidence.
- [Severity Classification And Recommended Actions](severity-classification-and-recommended-actions.md): Phase 5 deterministic severity rules, recommendation output, approval-required flags, test commands, and red-to-green evidence.
- [Shareable Incident Brief Drafting](incident-brief-drafting.md): Phase 6 structured draft brief sections, citation rules, redaction behavior, fail-closed behavior, approval-state display, test commands, and red-to-green evidence.
- [Scope And Guardrails](scope.md): in-scope, out-of-scope, trust boundaries, and demo path.
- [Phases And Tasks](phases.md): tickable phase plan with prompts.
- [Task Prompts](task-prompts.md): reusable prompts for future implementation or documentation work.
- [Strict TDD Rules](tdd-rules.md): non-negotiable rules for code tasks.
- [Eval Plan](eval-plan.md): metrics, fixtures, release gates, and risk checks.
- [Demo Package](demo-package.md): checklist for repo narrative, architecture diagram, demo video, and eval summary.
