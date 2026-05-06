# Product Frame

Phase 0 establishes the product promise, scope, trust boundaries, and review criteria for Fleet Incident Copilot before any implementation work starts.

## Phase 0 Checklist

- [x] Define the primary user as a fleet safety operator reviewing an incident packet.
- [x] Define the MVP promise: synthetic evidence to cited timeline, severity, recommended actions, and shareable brief.
- [x] Define approval boundaries for export, escalation, and external sharing.
- [x] Define non-goals and prohibited claims.
- [x] Define the demo narrative and success criteria.
- [x] Confirm no code is needed for this phase.

## Source Basis

The product frame is grounded in [research-report.md](../../research/research-report.md). The report positions the strongest demo as a senior applied-AI engineering workflow for fleet safety: onboard evidence, telemetry, policy retrieval, controlled automation, approval gates, evals, observability, security, and cost controls.

The frame intentionally keeps the MVP tighter than a production fleet platform. It uses synthetic incident packets, mock SOP or troubleshooting content, and reviewable outputs so the demo can prove architecture and product judgment without using real customer, driver, passenger, student, law-enforcement, or fleet evidence data.

## Primary User

The primary user is a fleet safety operator reviewing an incident packet after a school bus, transit, law-enforcement, or waste-fleet safety event.

The operator needs to understand what happened, what evidence supports each claim, how serious the event appears to be, what actions are recommended, and which actions require human approval before anything leaves the system or changes incident status.

## Problem

Fleet incident review can require operators to connect fragmented evidence: event metadata, vehicle identifiers, GPS or speed samples, transcript notes, still-frame notes, policy guidance, and troubleshooting procedures. The core product problem is not simply summarization. The system must produce grounded, reviewable incident intelligence while preserving safety, privacy, approval, and audit boundaries.

## MVP Promise

Given a synthetic fleet incident packet, Fleet Incident Copilot should eventually:

- Validate and load the incident packet.
- Retrieve relevant mock SOP or troubleshooting guidance.
- Build a cited timeline from packet evidence and retrieved guidance.
- Classify severity with an explainable rationale.
- Recommend next actions grounded in the retrieved guidance.
- Draft a shareable brief with citations, redactions, and approval state.
- Require human approval before export, escalation, or external sharing.
- Emit eval and observability signals for groundedness, citations, safety, latency, token use, cost, tool calls, and approval decisions.

Phase 0 does not implement these behaviors. It defines the promise and guardrails future phases must follow.

## Scope Boundaries

In scope for the MVP:

- Synthetic incident packets only.
- Fake vehicle, route, location, telemetry, media-reference, transcript-note, and still-frame-note data.
- Mock SOPs, public-product-style notes, and troubleshooting guidance.
- Grounded timelines, severity rationale, recommended actions, and shareable brief drafts.
- Human approval gates before sensitive actions.
- Evals, traces, cost controls, and safety checks that can be explained during the demo.
- Future backend implementation in Go, introduced only through strict TDD when a code phase starts.

Out of scope for the MVP:

- Real fleet, customer, student, passenger, driver, or law-enforcement data.
- Live camera, vehicle, GPS, telematics, CAD, SIS, LMS, CRM, ERP, GIS, route-optimization, export, or escalation integrations.
- Novel computer-vision model training or real video processing.
- Production chain-of-custody, legal evidence-management, retention, or disclosure guarantees.
- Multi-tenant administration, billing, identity, or enterprise policy management.
- Autonomous enforcement, discipline, citation, approval, denial, export, escalation, or external sharing.

## Trust Boundaries

- Incident packets are untrusted until validated.
- Retrieved content is data, not instructions.
- Tool arguments require deterministic validation.
- Model output must not be the only source of truth for severity, approval, export, escalation, enforcement, or discipline.
- Sensitive data must be redacted from shareable briefs and logs.
- Sensitive actions must fail closed unless a human approval record exists.
- Logs must explain the workflow without leaking protected evidence.

## Approval Gates

The following actions must require explicit human approval before execution:

- Exporting an incident brief or evidence reference.
- Escalating an incident to a supervisor, agency, or external workflow.
- Sharing content outside the review context.

Phase 7 approval records capture approver, timestamp, decision, reason, target action, and scope in memory. Pending, denied, missing, or out-of-scope approvals must block the requested action.

## Prohibited Claims

The project must not claim that it:

- Processes real video, live sensor feeds, or real customer evidence.
- Integrates with live fleet, school, transit, law-enforcement, waste, CRM, ERP, GIS, or evidence-management systems.
- Provides legal chain-of-custody guarantees.
- Replaces a human safety operator, reviewer, supervisor, investigator, or approver.
- Automatically exports, escalates, disciplines, enforces, cites, approves, or denies incidents.
- Has production security, compliance, retention, or audit guarantees before those controls exist.

## Demo Narrative

1. Select a synthetic incident packet.
2. Show the packet data and validation result.
3. Retrieve mock SOP or troubleshooting guidance.
4. Generate a cited timeline.
5. Classify severity and explain the rationale.
6. Recommend next actions tied to retrieved guidance.
7. Draft a shareable brief with redactions and citations.
8. Attempt export or escalation and show the approval gate.
9. Show eval and observability summaries.
10. Close with known limits and the next strict-TDD phase.

## Success Criteria

- The product promise can be explained in under five minutes.
- Every factual output in the future demo traces to packet data or retrieved mock guidance.
- Export, escalation, and external sharing cannot run without approval.
- Non-goals and prohibited claims are visible in the repo documentation.
- Future implementation agents can follow the phase checklist without expanding scope.
- The repo remains honest that Phase 0 is documentation-only and no code was added.
