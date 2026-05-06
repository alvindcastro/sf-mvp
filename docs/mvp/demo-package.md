# Demo Package

## Public Repo Narrative

- [ ] Explain the MVP as a Fleet Incident Copilot for synthetic fleet-safety incident review.
- [ ] State that the demo uses synthetic incident packets and mock SOPs.
- [ ] Highlight RAG, constrained agent tools, approval gates, evals, observability, security, and cost controls.
- [ ] Include quick-start instructions only after implementation exists.
- [ ] Avoid claiming live integrations or production evidence guarantees.

## Architecture Diagram Checklist

- [ ] Synthetic incident packet input.
- [ ] Validation and audit boundary.
- [ ] Mock SOP and troubleshooting corpus.
- [ ] Retrieval and citation layer.
- [ ] Timeline builder.
- [ ] Severity and recommendation engine.
- [ ] Brief drafting layer.
- [ ] Human approval gate.
- [ ] Mock export/escalation tools.
- [ ] Eval harness.
- [ ] Observability and cost-control signals.

## Demo Video Outline

- [ ] Open with the synthetic incident packet.
- [ ] Show validation and retrieved guidance.
- [ ] Show cited timeline.
- [ ] Show severity and rationale.
- [ ] Show recommended next actions.
- [ ] Show shareable brief with redactions.
- [ ] Attempt export or escalation and show the approval gate.
- [ ] Show eval and observability summary.
- [ ] Close with known limits and next steps.

## Interview Talking Points

- [ ] RAG: retrieval is scoped, cited, evaluated, and resilient to prompt injection.
- [ ] Agents: tools are constrained, validated, logged, and approval-gated.
- [ ] Backend APIs: incident workflow is modeled as explicit contracts and state transitions.
- [ ] Evals: quality is measured through deterministic checks and golden synthetic incidents.
- [ ] Monitoring: traces, latency, token use, retrieval quality, tool calls, and approvals are visible.
- [ ] Security: least-privilege retrieval, untrusted content boundaries, redaction, and fail-closed behavior.
- [ ] Cost: token budgets, caching candidates, routing notes, and budget-exceeded behavior are planned.
- [ ] Production readiness: the demo shows how the system would be governed, measured, and operated.

## Final Packaging Checklist

- [ ] Public repo has clear narrative.
- [ ] Short demo video script is ready.
- [ ] Architecture diagram is accurate.
- [ ] One-page eval summary is complete.
- [ ] README does not overclaim implementation.
- [ ] All fixtures are synthetic.
- [ ] Known risks and next steps are explicit.

