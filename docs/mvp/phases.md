# Phases And Tasks

Each phase is a planning unit for future work. Code phases must follow [Strict TDD Rules](tdd-rules.md): add or update a failing test first, observe red, implement the smallest green change, refactor only after green, then document what was tested.

## Phase 0: Product Frame And Guardrails

Goal: establish the product promise, scope, trust boundaries, and review criteria before implementation.

- [x] Define the primary user as a fleet safety operator reviewing an incident packet.
- [x] Define the MVP promise: synthetic evidence to cited timeline, severity, recommended actions, and shareable brief.
- [x] Define approval boundaries for export, escalation, and external sharing.
- [x] Define non-goals and prohibited claims.
- [x] Define the demo narrative and success criteria.
- [x] Confirm no code is needed for this phase.

Output: [Product Frame](product-frame.md). Phase 0 is complete as of 2026-05-06 and added Markdown planning artifacts only.

Prompt:

> Using `research-report.md`, write or refine the product frame for Fleet Incident Copilot. Include the user, problem, demo promise, scope boundaries, approval gates, success criteria, and explicit non-goals. Do not write code. Keep the output in Markdown and preserve a checklist that future implementation agents can follow.

## Phase 1: Synthetic Evidence And Workflow Contract

Goal: define realistic fake incidents and expected workflow outputs without introducing real data.

- [x] Design at least five synthetic incident packets.
- [x] Include incident ID, vehicle ID, route, timestamp, location label, event type, telemetry samples, media references, transcript notes, and still-frame notes.
- [x] Include low, medium, high, unknown, and adversarial/missing-data cases.
- [x] Define expected timeline, severity, recommended actions, and brief requirements for each packet.
- [x] Mark all records as synthetic.
- [x] Confirm no implementation code is needed until fixtures and acceptance criteria are settled.

Output: [Synthetic Incident Packets](synthetic-incident-packets.md). Phase 1 is complete as of 2026-05-06 and added Markdown planning artifacts only.

Prompt:

> Create Markdown specs for synthetic Fleet Incident Copilot incident packets. Use realistic but fake fleet-safety scenarios. For each packet, include required fields, expected outputs, missing-data behavior, and acceptance criteria. Do not write application code. If fixture files are created in a future code run, they must be introduced through failing tests first.

## Phase 2: Incident Packet Ingestion

Goal: validate synthetic incident packets before any reasoning step can use them.

- [x] Define packet schema and validation rules.
- [x] Reject missing incident IDs, timestamps, event types, telemetry arrays, or evidence references.
- [x] Reject malformed timestamps, impossible speed samples, and unsupported event types.
- [x] Produce actionable validation errors.
- [x] Emit an ingestion audit event.
- [x] Keep all examples synthetic.

Output: [Incident Packet Ingestion](incident-packet-ingestion.md) and the Go package [internal/ingestion](../../internal/ingestion). Phase 2 is complete as of 2026-05-06 and added the first strict-TDD runtime implementation for synthetic packet validation only.

Code-task prompt:

> Implement synthetic incident packet ingestion using strict TDD. First identify the smallest observable behavior: accepting one valid synthetic packet and rejecting one invalid packet. Add failing tests for valid ingestion, missing required fields, malformed telemetry, unsupported event type, and audit event emission. Run the targeted tests and confirm they fail for the expected reason. Then implement only enough parser and validation logic to pass. Refactor only after tests are green. Acceptance requires deterministic validation errors, no real customer data, and a test summary showing red-to-green evidence.

## Phase 3: RAG Corpus And Grounding

Goal: future implementation retrieves only approved mock guidance and preserves citations.

- [x] Define mock SOP documents.
- [x] Define troubleshooting notes.
- [x] Define document metadata: source ID, title, workflow, tenant/scope marker, and revision date.
- [x] Define citation format for retrieved snippets.
- [x] Define no-match behavior.
- [x] Include a prompt-injection document fixture that must be treated as untrusted content.
- [x] Define retrieval eval questions before implementation.

Output: [RAG Corpus And Grounding](rag-corpus-and-grounding.md) and the Go package [internal/retrieval](../../internal/retrieval). Phase 3 is complete as of 2026-05-06 and added the first strict-TDD retrieval implementation for approved mock guidance only.

Code-task prompt:

> Implement the first RAG retrieval slice using strict TDD. Start with failing tests for relevant SOP retrieval, citation metadata, no-match behavior, scope filtering, and hostile retrieved text that tries to override system instructions. Confirm the tests fail before production code changes. Implement only the retrieval interface needed for the MVP fixture corpus. Acceptance requires cited snippets, deterministic filtering, no unauthorized documents in context, and tests proving retrieved text is data rather than instructions.

## Phase 4: Incident Timeline Builder

Goal: future implementation builds a chronological, grounded account from packet data and retrieved guidance.

- [x] Order packet events and telemetry chronologically.
- [x] Incorporate transcript or still-frame notes without inventing visual facts.
- [x] Preserve source references for every factual claim.
- [x] Mark uncertainty when data is missing or conflicting.
- [x] Avoid unsupported claims.
- [x] Produce a timeline format suitable for the incident brief.

Output: [Incident Timeline Builder](incident-timeline-builder.md) and the Go package [internal/timeline](../../internal/timeline). Phase 4 is complete as of 2026-05-06 and added the first strict-TDD timeline implementation for grounded synthetic packet timelines only.

Code-task prompt:

> Implement incident timeline generation using strict TDD. Begin with failing tests for chronological ordering, source citation, uncertainty labeling, missing-data handling, and conflict handling. Run the tests and confirm they fail for the expected reason. Then implement the smallest timeline builder that passes. Acceptance requires every timeline claim to trace to packet data or retrieved source metadata, and tests proving unsupported claims are omitted or marked unknown.

## Phase 5: Severity Classification And Recommended Actions

Goal: classify severity and recommend next actions through explainable rules before optional model judgment.

- [x] Define low, medium, high, and unknown severity rules.
- [x] Prefer deterministic rule output for the MVP.
- [x] Isolate model-dependent judgment behind an interface if added later.
- [x] Tie recommendations to severity, event type, and retrieved SOPs.
- [x] Explain why each recommendation was made.
- [x] Flag export, escalation, and external sharing as approval-required.

Output: [Severity Classification And Recommended Actions](severity-classification-and-recommended-actions.md) and the Go package [internal/severity](../../internal/severity). Phase 5 is complete as of 2026-05-06 and added the first strict-TDD deterministic severity and recommendation implementation for synthetic packet review only.

Code-task prompt:

> Implement severity classification and recommended next actions using strict TDD. First add failing tests for low, medium, high, unknown, conflicting-signal, explanation, and approval-required scenarios. Confirm red before writing production code. Implement deterministic rules first and keep any model-dependent behavior out of the initial path unless it is already covered by tests. Acceptance requires explainable outputs, SOP-grounded recommendations, approval flags for sensitive actions, and tests proving no escalation or export occurs automatically.

## Phase 6: Shareable Incident Brief

Goal: draft a concise, cited, redacted brief for human review.

- [x] Include incident summary, cited timeline, severity, rationale, next actions, and approval state.
- [x] Include citations for factual claims.
- [x] Redact sensitive fields from shareable output.
- [x] Fail closed when required evidence is missing.
- [x] Label uncertainty clearly.
- [x] Keep the brief draft human-reviewable, not final by default.

Output: [Shareable Incident Brief Drafting](incident-brief-drafting.md) and the Go package [internal/brief](../../internal/brief). Phase 6 is complete as of 2026-05-06 and added the first strict-TDD structured draft brief implementation for cited, redacted human-review output only.

Code-task prompt:

> Implement incident brief drafting using strict TDD. Add failing tests for a complete brief, citation inclusion, sensitive-field redaction, missing-evidence failure, uncertainty labeling, and approval-state display. Confirm the tests fail before production changes. Then implement the smallest drafting layer needed to pass. Acceptance requires no uncited factual claims, no sensitive leakage in shareable output, and tests showing draft creation fails closed when required evidence is absent.

## Phase 7: Human Approval Workflow

Goal: future implementation blocks sensitive actions until a human decision exists.

- [ ] Create approval requests for export, escalation, and external sharing.
- [ ] Capture approver, timestamp, decision, reason, and target action.
- [ ] Block sensitive tool calls while approval is pending.
- [ ] Block denied actions.
- [ ] Allow approved actions only within the approved scope.
- [ ] Preserve append-only audit history.

Code-task prompt:

> Implement the human approval gate using strict TDD. First write failing tests for pending approval creation, denied approval, granted approval, blocked action before approval, scoped approval, and immutable audit history. Confirm red before production code. Then implement the smallest workflow that passes. Acceptance requires sensitive actions to fail closed, audit records to be append-only, and tests proving denied or out-of-scope actions cannot execute.

## Phase 8: Eval Harness

Goal: future implementation measures groundedness, citations, severity, recommendations, safety, and redaction before demo release.

- [ ] Create golden synthetic incidents.
- [ ] Score expected severity.
- [ ] Score citation coverage.
- [ ] Detect unsupported claims.
- [ ] Check recommendation accuracy against expected SOP guidance.
- [ ] Check prompt-injection resistance.
- [ ] Check redaction behavior.
- [ ] Define release thresholds.

Code-task prompt:

> Implement the MVP eval harness using strict TDD. Start with failing tests for loading eval cases, scoring expected severity, checking citation coverage, detecting unsupported claims, verifying redaction, and handling prompt-injection fixtures. Confirm failures before implementation. Then implement the smallest evaluator. Acceptance requires repeatable local evals, clear pass/fail thresholds, normal/adversarial/incomplete fixtures, and a test summary showing what went red and green.

## Phase 9: Observability And Cost Controls

Goal: future implementation exposes enough runtime signals to explain quality, safety, latency, and spend.

- [ ] Generate a trace ID per incident workflow.
- [ ] Track retrieval count, retrieved source IDs, tool-call success, approval decisions, latency, token usage, and eval score.
- [ ] Redact sensitive data from logs.
- [ ] Add model-call budget limits.
- [ ] Define caching candidates.
- [ ] Define model-routing notes for hosted, smaller, or self-hosted models.

Code-task prompt:

> Implement observability and cost controls using strict TDD. Add failing tests for trace propagation, structured event emission, token recording, latency recording, budget-limit behavior, and sensitive-field redaction in logs. Confirm red before production changes. Then implement instrumentation with the smallest surface needed. Acceptance requires structured logs, useful debugging signals, no sensitive evidence leakage, and tests for normal and budget-exceeded paths.

## Phase 10: Demo Package

Goal: package the MVP so it communicates production readiness clearly.

- [ ] Write a repo narrative that maps the MVP to the target role.
- [ ] Prepare a short demo video outline.
- [ ] Prepare one architecture diagram checklist.
- [ ] Prepare a one-page eval summary.
- [ ] Prepare interview talking points for RAG, agents, backend APIs, evals, monitoring, security, cost, and production readiness.
- [ ] Confirm no feature is described as implemented unless it exists.

Prompt:

> Create the Fleet Incident Copilot demo packaging materials. Include a repo narrative, short demo video script, architecture diagram checklist, one-page eval summary outline, and interview talking points. Tie the artifacts to RAG, agents, backend APIs, evals, monitoring, security, cost, and production readiness. Do not claim implementation that does not exist.
