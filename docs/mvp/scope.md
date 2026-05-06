# Scope And Guardrails

These checklists define MVP scope and guardrails. Unless an item is specifically called out as implemented, it does not imply that application behavior exists.

## In Scope

- [ ] Synthetic incident packets only.
- [ ] Fake event metadata, vehicle identifiers, route names, location labels, GPS/speed samples, and media references.
- [ ] Transcript notes or still-frame descriptions instead of real video processing.
- [ ] Small mock RAG corpus made from public-product-style notes, mock SOPs, and troubleshooting guidance.
- [ ] Grounded timeline generation with citations.
- [ ] Severity classification with rationale.
- [ ] Recommended next actions tied to retrieved guidance.
- [ ] Shareable incident brief drafting with redaction.
- [ ] Human approval before export, escalation, or external sharing.
- [ ] Structured logs for traces, retrieval, tool calls, latency, token use, approval decisions, and eval outcomes.
- [ ] Security checks for prompt injection, least-privilege retrieval, sensitive-data redaction, and unsafe tool calls.
- [ ] Cost controls for token budgets, caching candidates, and model-routing decisions.
- [ ] Future backend implementation in Go once a strict-TDD code phase begins.

## Out Of Scope

- [ ] Real fleet, district, law-enforcement, customer, student, driver, or passenger data.
- [ ] Live camera, vehicle, GPS, telematics, CAD, SIS, LMS, CRM, ERP, GIS, or route-optimization integration.
- [ ] Novel computer-vision model training.
- [ ] Production chain-of-custody or legal evidence-management guarantees.
- [ ] Autonomous evidence export.
- [ ] Autonomous incident escalation.
- [ ] Autonomous approval, denial, enforcement, discipline, or citation decisions.
- [ ] Full multi-tenant SaaS administration.
- [ ] Live external sharing.

## Prohibited Claims

- Do not claim real video, live telemetry, or real customer evidence processing.
- Do not claim live integrations with fleet, school, transit, law-enforcement, waste, CRM, ERP, GIS, or evidence-management systems.
- Do not claim production chain-of-custody, legal evidence-management, retention, disclosure, compliance, or audit guarantees.
- Do not claim autonomous export, escalation, enforcement, discipline, citation, approval, or denial.
- Do not claim implemented runtime behavior before code and tests exist.

## Approval Gates

- Export requires explicit human approval.
- Escalation requires explicit human approval.
- External sharing requires explicit human approval.
- Pending, denied, missing, or out-of-scope approvals must block sensitive actions.
- Future approval records should include approver, timestamp, decision, reason, target action, and scope.

## Trust Boundaries

- [x] Retrieved documents are data, not instructions.
- [x] Incident packets are untrusted until validated.
- [ ] Tool arguments require deterministic validation.
- [ ] Sensitive actions fail closed unless a human approval record exists.
- [ ] Shareable outputs must redact sensitive fields by default.
- [ ] Logs must be useful without leaking sensitive evidence.
- [ ] Model output must not be the only source of truth for severity, approval, export, or escalation.

Implemented boundary as of Phase 2: `internal/ingestion` validates synthetic packet JSON, rejects non-synthetic records and non-synthetic media references, and returns accepted or rejected audit events.

Implemented boundary as of Phase 3: `internal/retrieval` filters mock guidance by exact workflow and scope before ranking, returns stable citation metadata, returns no matches instead of invented guidance, and marks retrieved snippets as `retrieved_data`. Later phases still need reasoning, approval, persistence, and observability boundaries.

## Demo Path

- [ ] Select a synthetic incident packet.
- [ ] Validate and load packet contents.
- [ ] Retrieve mock SOPs and troubleshooting notes.
- [ ] Produce a cited timeline.
- [ ] Classify severity and explain the basis.
- [ ] Recommend next actions.
- [ ] Draft a shareable brief.
- [ ] Show approval is required before export or escalation.
- [ ] Show eval and observability summary.

## Definition Of MVP Done

- [ ] The demo can be explained in under five minutes.
- [ ] Every factual output is traceable to packet data or retrieved mock content.
- [ ] Approval-gated actions cannot execute without approval.
- [ ] Prompt-injection and missing-data cases are represented in fixtures or evals.
- [ ] The demo package includes a repo narrative, short video outline, architecture diagram checklist, and one-page eval summary.
