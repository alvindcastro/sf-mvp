# Eval Plan

## Eval Goals

- [ ] Prove outputs are grounded in packet evidence or retrieved mock guidance.
- [ ] Prove citations are present and useful.
- [ ] Prove severity classifications match expected outcomes.
- [ ] Prove recommendations align with mock SOPs.
- [ ] Prove sensitive fields are redacted in shareable output and logs.
- [ ] Prove approval-gated actions cannot run automatically.
- [ ] Prove prompt-injection content is treated as untrusted data.
- [ ] Prove latency, token use, and cost controls are visible.

## Fixture Set

Phase 1 defines human-readable candidate packet specs in [Synthetic Incident Packets](synthetic-incident-packets.md). Machine-readable eval fixtures are not implemented yet and must be introduced through a future strict-TDD code phase.

- [ ] Normal low-severity incident.
- [ ] Normal medium-severity incident.
- [ ] Normal high-severity incident.
- [ ] Unknown-severity incident with missing evidence.
- [ ] Conflicting telemetry incident.
- [ ] Prompt-injection incident through retrieved mock document text.
- [ ] Sensitive-data redaction incident.
- [ ] Budget-limit incident that exceeds model-call allowance.

Phase 1 packet mapping:

- Low severity: `FIC-SYN-001`.
- Medium severity: `FIC-SYN-002`.
- High severity: `FIC-SYN-003`.
- Unknown severity with missing evidence: `FIC-SYN-004`.
- Adversarial or missing-data behavior: `FIC-SYN-005`.

## Metrics

- [ ] Groundedness: percentage of factual claims traced to packet data or retrieved source IDs.
- [ ] Citation coverage: percentage of timeline and brief claims with citations.
- [ ] Severity accuracy: expected severity versus produced severity.
- [ ] Recommendation accuracy: expected SOP-grounded action versus produced action.
- [ ] Tool-call success: valid tool calls divided by attempted tool calls.
- [ ] Approval compliance: sensitive actions blocked without approval.
- [ ] Redaction quality: sensitive fields absent from shareable briefs and logs.
- [ ] Injection resistance: hostile retrieved text does not change system behavior.
- [ ] Latency: p50 and p95 workflow time.
- [ ] Cost: tokens per incident and budget-limit behavior.

## Release Gates

- [ ] All deterministic unit and integration tests pass.
- [ ] Eval fixture loading is repeatable locally.
- [ ] Severity accuracy passes the defined threshold.
- [ ] Citation coverage passes the defined threshold.
- [ ] No unsupported high-risk factual claims appear in briefs.
- [ ] All sensitive actions require approval.
- [ ] Prompt-injection fixtures fail safely.
- [ ] Logs contain trace IDs and required metrics without sensitive evidence.
- [ ] Demo packet outputs can be explained with citations.

## One-Page Eval Summary Outline

- [ ] MVP behavior evaluated.
- [ ] Fixture count and categories.
- [ ] Metrics tracked.
- [ ] Thresholds used.
- [ ] Results table.
- [ ] Known failure modes.
- [ ] Risk controls.
- [ ] Next improvements.
