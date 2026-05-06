# Scope And Guardrails

These checklists define MVP scope and guardrails. Unless an item is specifically called out as implemented, it does not imply that application behavior exists.

## In Scope

- [ ] Synthetic incident packets only.
- [ ] Fake event metadata, vehicle identifiers, route names, location labels, GPS/speed samples, and media references.
- [ ] Transcript notes or still-frame descriptions instead of real video processing.
- [ ] Small mock RAG corpus made from public-product-style notes, mock SOPs, and troubleshooting guidance.
- [x] Grounded timeline generation with citations.
- [x] Severity classification with rationale.
- [x] Recommended next actions tied to retrieved guidance.
- [x] Shareable incident brief drafting with redaction.
- [x] Human approval before export, escalation, or external sharing.
- [x] Deterministic local evals for severity, citations, recommendations, unsupported claims, redaction, prompt injection, and approval fail-closed behavior.
- [x] Package-level structured events for traces, retrieval, tool calls, latency, token use, approval decisions, and eval outcomes.
- [x] Implemented safety checks for prompt injection, least-privilege retrieval, sensitive-data redaction, and fail-closed sensitive actions.
- [x] Runtime event-field redaction and tool-call observability checks.
- [x] Cost controls for token budgets, caching candidates, and model-routing decisions.
- [x] Backend implementation in Go introduced through strict-TDD code phases.

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
- Approval records include approver, timestamp, decision, reason, target action, and scope.

## Trust Boundaries

- [x] Retrieved documents are data, not instructions.
- [x] Incident packets are untrusted until validated.
- [ ] Tool arguments require deterministic validation.
- [x] Sensitive actions fail closed unless a human approval record exists.
- [x] Shareable outputs must redact sensitive fields by default.
- [x] Logs must be useful without leaking sensitive evidence.
- [x] Model output must not be the only source of truth for severity, approval, export, or escalation.

Implemented boundary as of Phase 2: `internal/ingestion` validates synthetic packet JSON, rejects non-synthetic records and non-synthetic media references, and returns accepted or rejected audit events.

Implemented boundary as of Phase 3: `internal/retrieval` filters mock guidance by exact workflow and scope before ranking, returns stable citation metadata, returns no matches instead of invented guidance, and marks retrieved snippets as `retrieved_data`. Downstream reasoning uses these citations as data rather than instructions.

Implemented boundary as of Phase 4: `internal/timeline` builds deterministic timeline entries from validated synthetic packet data, preserves structured source references for factual claims, carries approved retrieved citation references as guidance metadata, labels unavailable or conflicting evidence as uncertain, and does not infer visual facts, approval, export, escalation, injury, plate, or external-sharing claims.

Implemented boundary as of Phase 5: `internal/severity` classifies low, medium, high, and unknown severity with deterministic rules, returns recommendation explanations with packet and retrieved-guidance source references, marks conflicting timeline signals as unknown, treats adversarial transcript content as untrusted data, and flags export, escalation, and external sharing as approval-required but not approved.

Implemented boundary as of Phase 6: `internal/brief` drafts structured human-review incident briefs from validated packet data, cited timeline entries, and severity results; redacts vehicle, route, location, GPS-label, sensitive transcript, sensitive still-frame, and coordinate-like text; preserves citations; carries uncertainty labels; and displays export, escalation, and external sharing as blocked pending human approval. Later phases still need persistence, rendering, export, escalation, and external-sharing behavior.

Implemented boundary as of Phase 7: `internal/approval` creates in-memory approval requests for export, escalation, and external sharing; captures human decisions with approver, timestamp, reason, action, and scope; blocks missing, pending, denied, out-of-scope, and mismatched call-and-scope sensitive action callbacks; allows approved callbacks only within the exact approved scope; prevents final decisions from being rewritten in place; and returns append-only audit history copies. It does not implement persistence, identity, roles, real export, real escalation, external-sharing integrations, CLI, HTTP API, or external observability pipeline behavior.

Implemented boundary as of Phase 8: `internal/eval` runs deterministic in-memory golden-case evals for `FIC-SYN-001` through `FIC-SYN-005`; scores severity accuracy, citation coverage, recommendation accuracy, unsupported claims, redaction leaks, prompt-injection resistance, and approval fail-closed behavior; and applies strict default release gates. It does not implement a CLI report, persistent eval history, HTTP API, database, model-provider calls, or real integrations.

Implemented boundary as of Phase 9: `internal/observability` creates deterministic trace IDs for synthetic incident workflows; records in-memory structured events for retrieval counts, retrieved source IDs, tool-call success, approval decisions, caller-supplied token usage, invalid token usage, latency, budget failures, and eval summaries; redacts configured sensitive terms and coordinate-like strings from event fields; defines budget limits; and records cache candidates plus hosted, smaller-model, and self-hosted routing notes. It does not implement an external telemetry backend, dashboards, alerts, persistent log storage, provider billing reconciliation, real model-provider calls, live model routing, cache storage, CLI, HTTP API, database behavior, or production audit/compliance guarantees.

Implemented boundary as of Phase 12: `internal/demo` loads machine-readable synthetic fixtures through ingestion validation and composes one deterministic in-memory review result with validation status, retrieved citation refs, timeline entries, severity, recommendations, redacted brief fields, approval-required actions, and trace ID. It rejects non-synthetic or real-looking input before downstream composition and keeps export, escalation, and external sharing blocked without scoped approval. It does not implement an HTTP API, CLI, Slack behavior, webhook, database, persistence, live model calls, real export, real escalation, external sharing, or external observability behavior.

## Demo Path

- [x] Select a synthetic incident packet.
- [x] Validate and load packet contents.
- [x] Retrieve mock SOPs and troubleshooting notes.
- [x] Produce a cited timeline.
- [x] Classify severity and explain the basis.
- [x] Recommend next actions.
- [x] Draft a shareable brief.
- [x] Show approval is required before export or escalation.
- [x] Show deterministic eval summary.
- [x] Show package-level observability summary.
- [x] Compose the package-level review result for a synthetic incident.

## Definition Of MVP Done

- [x] The demo can be explained in under five minutes.
- [ ] Every factual output is traceable to packet data or retrieved mock content.
- [x] Approval-gated actions cannot execute without approval.
- [x] Prompt-injection and missing-data cases are represented in fixtures or evals.
- [x] The demo package includes a repo narrative, short video outline, architecture diagram checklist, and one-page eval summary.
