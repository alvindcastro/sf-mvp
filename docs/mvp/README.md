# Fleet Incident Copilot MVP

## Thesis

Build a Fleet Incident Copilot that demonstrates senior applied-AI engineering for fleet-safety workflows: synthetic onboard evidence, telemetry, policy retrieval, incident summarization, workflow automation, approval gates, evals, observability, security, and cost controls.

The MVP should show production judgment, not just LLM usage. It should make it easy to explain how the system is grounded, measured, constrained, and operated.

## Demo Promise

Given a synthetic fleet incident packet, the system should:

- [ ] Ingest structured event metadata, GPS or speed samples, transcript notes, and still-frame notes.
- [ ] Retrieve relevant mock SOP or troubleshooting guidance.
- [ ] Build a cited incident timeline.
- [ ] Classify severity with rationale.
- [ ] Recommend next actions grounded in retrieved guidance.
- [ ] Draft a shareable incident brief with citations and redactions.
- [ ] Require human approval before export, escalation, or external sharing.
- [ ] Emit observability and eval signals for traces, retrieval quality, latency, token use, tool calls, and safety checks.

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

- [Scope And Guardrails](scope.md): in-scope, out-of-scope, trust boundaries, and demo path.
- [Phases And Tasks](phases.md): tickable phase plan with prompts.
- [Task Prompts](task-prompts.md): reusable prompts for future implementation or documentation work.
- [Strict TDD Rules](tdd-rules.md): non-negotiable rules for code tasks.
- [Eval Plan](eval-plan.md): metrics, fixtures, release gates, and risk checks.
- [Demo Package](demo-package.md): checklist for repo narrative, architecture diagram, demo video, and eval summary.

