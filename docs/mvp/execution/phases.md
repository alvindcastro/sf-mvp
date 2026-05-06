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

Output: [Product Frame](../overview/product-frame.md). Phase 0 is complete as of 2026-05-06 and added Markdown planning artifacts only.

Prompt:

> Using `docs/research/research-report.md`, write or refine the product frame for Fleet Incident Copilot. Include the user, problem, demo promise, scope boundaries, approval gates, success criteria, and explicit non-goals. Do not write code. Keep the output in Markdown and preserve a checklist that future implementation agents can follow.

## Phase 1: Synthetic Evidence And Workflow Contract

Goal: define realistic fake incidents and expected workflow outputs without introducing real data.

- [x] Design at least five synthetic incident packets.
- [x] Include incident ID, vehicle ID, route, timestamp, location label, event type, telemetry samples, media references, transcript notes, and still-frame notes.
- [x] Include low, medium, high, unknown, and adversarial/missing-data cases.
- [x] Define expected timeline, severity, recommended actions, and brief requirements for each packet.
- [x] Mark all records as synthetic.
- [x] Confirm no implementation code is needed until fixtures and acceptance criteria are settled.

Output: [Synthetic Incident Packets](../workflow/synthetic-incident-packets.md). Phase 1 is complete as of 2026-05-06 and added Markdown planning artifacts only.

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

Output: [Incident Packet Ingestion](../workflow/incident-packet-ingestion.md) and the Go package [internal/ingestion](../../../internal/ingestion). Phase 2 is complete as of 2026-05-06 and added the first strict-TDD runtime implementation for synthetic packet validation only.

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

Output: [RAG Corpus And Grounding](../workflow/rag-corpus-and-grounding.md) and the Go package [internal/retrieval](../../../internal/retrieval). Phase 3 is complete as of 2026-05-06 and added the first strict-TDD retrieval implementation for approved mock guidance only.

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

Output: [Incident Timeline Builder](../workflow/incident-timeline-builder.md) and the Go package [internal/timeline](../../../internal/timeline). Phase 4 is complete as of 2026-05-06 and added the first strict-TDD timeline implementation for grounded synthetic packet timelines only.

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

Output: [Severity Classification And Recommended Actions](../workflow/severity-classification-and-recommended-actions.md) and the Go package [internal/severity](../../../internal/severity). Phase 5 is complete as of 2026-05-06 and added the first strict-TDD deterministic severity and recommendation implementation for synthetic packet review only.

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

Output: [Shareable Incident Brief Drafting](../workflow/incident-brief-drafting.md) and the Go package [internal/brief](../../../internal/brief). Phase 6 is complete as of 2026-05-06 and added the first strict-TDD structured draft brief implementation for cited, redacted human-review output only.

Code-task prompt:

> Implement incident brief drafting using strict TDD. Add failing tests for a complete brief, citation inclusion, sensitive-field redaction, missing-evidence failure, uncertainty labeling, and approval-state display. Confirm the tests fail before production changes. Then implement the smallest drafting layer needed to pass. Acceptance requires no uncited factual claims, no sensitive leakage in shareable output, and tests showing draft creation fails closed when required evidence is absent.

## Phase 7: Human Approval Workflow

Goal: block sensitive actions until a human decision exists.

- [x] Create approval requests for export, escalation, and external sharing.
- [x] Capture approver, timestamp, decision, reason, and target action.
- [x] Block sensitive tool calls while approval is pending.
- [x] Block denied actions.
- [x] Allow approved actions only within the approved scope.
- [x] Preserve append-only audit history.

Output: [Human Approval Workflow](../workflow/human-approval-workflow.md) and the Go package [internal/approval](../../../internal/approval). Phase 7 is complete as of 2026-05-06 and added the first strict-TDD in-memory human approval gate for scoped sensitive-action callbacks only.

Code-task prompt:

> Implement or refine the human approval gate using strict TDD. Use `docs/mvp/workflow/human-approval-workflow.md` as the current Phase 7 behavior contract. First add or update failing tests for pending approval creation, denied approval, granted approval, blocked action before approval, scoped approval, and immutable audit history. Confirm red before production code. Then implement the smallest workflow that passes. Acceptance requires sensitive actions to fail closed, audit records to be append-only, final decisions to avoid in-place rewrites, and tests proving denied or out-of-scope actions cannot execute.

## Phase 8: Eval Harness

Goal: measure groundedness, citations, severity, recommendations, safety, approval fail-closed behavior, and redaction before demo release.

- [x] Create golden synthetic incidents.
- [x] Score expected severity.
- [x] Score citation coverage.
- [x] Detect unsupported claims.
- [x] Check recommendation accuracy against expected SOP guidance.
- [x] Check prompt-injection resistance.
- [x] Check redaction behavior.
- [x] Define release thresholds.

Output: [Eval Plan](../quality/eval-plan.md) and the Go package [internal/eval](../../../internal/eval). Phase 8 is complete as of 2026-05-06 and added a strict-TDD in-memory eval harness for deterministic synthetic golden cases only.

Code-task prompt:

> Implement or refine the MVP eval harness using strict TDD. Use `docs/mvp/quality/eval-plan.md` as the current Phase 8 behavior contract. Start with failing tests for loading eval cases, scoring expected severity, checking citation coverage, detecting unsupported claims, verifying redaction, checking recommendation accuracy, checking approval fail-closed behavior, and handling prompt-injection fixtures. Confirm failures before implementation. Then implement the smallest evaluator. Acceptance requires repeatable local evals, clear pass/fail thresholds, normal/adversarial/incomplete fixtures, and a test summary showing what went red and green.

## Phase 9: Observability And Cost Controls

Goal: expose enough runtime signals to explain quality, safety, latency, and spend.

- [x] Generate a trace ID per incident workflow.
- [x] Track retrieval count, retrieved source IDs, tool-call success, approval decisions, latency, token usage, and eval score.
- [x] Redact sensitive data from logs.
- [x] Add model-call budget limits.
- [x] Define caching candidates.
- [x] Define model-routing notes for hosted, smaller, or self-hosted models.

Output: [Observability And Cost Controls](../quality/observability-and-cost-controls.md) and the Go package [internal/observability](../../../internal/observability). Phase 9 is complete as of 2026-05-06 and added a strict-TDD package-level observability and budget-control surface for synthetic MVP workflows only.

Code-task prompt:

> Implement or refine observability and cost controls using strict TDD. Use `docs/mvp/quality/observability-and-cost-controls.md` as the current Phase 9 behavior contract. Add failing tests for trace propagation, structured event emission, token recording, latency recording, budget-limit behavior, and sensitive-field redaction in logs. Confirm red before production changes. Then implement instrumentation with the smallest surface needed. Acceptance requires structured logs, useful debugging signals, no sensitive evidence leakage, and tests for normal and budget-exceeded paths.

## Phase 10: Demo Package

Goal: package the MVP so it communicates production readiness clearly.

- [x] Write a repo narrative that maps the MVP to the target role.
- [x] Prepare a short demo video outline.
- [x] Prepare one architecture diagram checklist.
- [x] Prepare a one-page eval summary.
- [x] Prepare interview talking points for RAG, agents, backend APIs, evals, monitoring, security, cost, and production readiness.
- [x] Confirm no feature is described as implemented unless it exists.

Output: [Demo Package](../demo/demo-package.md). Phase 10 is complete as of 2026-05-06 and added Markdown demo packaging materials only: a repo narrative, target-role mapping, short demo script, architecture diagram checklist, one-page eval summary outline, interview talking points, and implemented-versus-planned wording rules.

Prompt:

> Create the Fleet Incident Copilot demo packaging materials. Include a repo narrative, short demo video script, architecture diagram checklist, one-page eval summary outline, and interview talking points. Tie the artifacts to RAG, agents, backend APIs, evals, monitoring, security, cost, and production readiness. Do not claim implementation that does not exist.

## Phase 11: Hiring-Manager Demo Surface Roadmap

Goal: brainstorm a concrete local demo surface without adding runtime code or overstating current implementation.

- [x] Compare API, CLI, dry-run Slack notification, approval retry, eval report, and observability proof options.
- [x] Choose a local loopback API plus dry-run Slack-shaped notification preview as the recommended future demo arc.
- [x] Keep real Slack delivery, webhooks, model calls, persistence, identity, dashboards, and production compliance claims out of scope.
- [x] Add future phase checklists for the demo surface.
- [x] Add detailed future-agent prompts for strict-TDD code tasks.
- [x] Confirm this phase adds Markdown planning artifacts only.

Output: [Demo Surface Roadmap](../demo/demo-surface-roadmap.md). Phase 11 is complete as of 2026-05-06 and added Markdown planning artifacts only. Phase 11 itself did not add API, CLI, Slack integration, webhook, live provider call, database, or persistent store behavior.

Prompt:

> Brainstorm a hiring-manager demo surface for Fleet Incident Copilot without writing code. Compare a local API, CLI, dry-run Slack-shaped notification preview, approval retry flow, eval report, and observability proof. Create or update Markdown docs with tickable phase tasks, detailed future-agent prompts, and implemented-versus-planned wording guardrails. Every future code task must require strict TDD. Do not claim the demo surface is implemented unless code and tests already prove it.

## Phase 12: Review Composition Contract

Goal: compose the existing package-level workflow into one deterministic demo review result.

- [x] Add machine-readable synthetic demo fixtures only after failing tests define fixture-loading expectations.
- [x] Define the smallest review response that includes validation status, retrieved citation refs, timeline entries, severity, recommendations, redacted brief, approval-required actions, and trace ID.
- [x] Reject non-synthetic or real-looking incident input before downstream composition.
- [x] Preserve existing citation, redaction, approval, eval, and observability package contracts.
- [x] Keep the composer in-memory and deterministic.
- [x] Document behavior only after tests prove it.

Output: [Review Composition Contract](../demo/review-composition-contract.md), the Go package [internal/demo](../../../internal/demo), and machine-readable synthetic fixtures in [internal/demo/testdata/demo-fixtures.json](../../../internal/demo/testdata/demo-fixtures.json). Phase 12 is complete as of 2026-05-06 and added an in-memory deterministic demo review composer only. Phase 12 itself did not add HTTP API, CLI, Slack behavior, database, persistence, live model call, webhook, export, escalation, or external-sharing integration behavior.

Code-task prompt:

> Implement or refine a demo review composition contract using strict TDD. Start by naming the smallest observable behavior: composing one known synthetic incident into a review result with incident ID, trace ID, severity, citation refs, redacted brief fields, and approval-required actions. Add failing tests before production code for the happy path, unknown incident ID, non-synthetic input, missing evidence, citation preservation, redaction preservation, and approval-required action display. Confirm red for the expected reason. Implement only enough in-memory composition logic to pass using existing `internal/ingestion`, `internal/retrieval`, `internal/timeline`, `internal/severity`, `internal/brief`, `internal/approval`, and `internal/observability` contracts. Run the targeted tests and then `go test ./...`. Do not add an HTTP server, Slack behavior, database, persistence, live model call, or external service in this phase.

Fixture-task prompt:

> Add or refine synthetic machine-readable demo fixtures using strict TDD before demo adapters depend on fixtures. Add a failing test such as `TestLoadDemoFixturesReturnsSyntheticNormalIncompleteAndAdversarialPackets`; confirm the targeted package test fails because the loader does not exist or the expected behavior is missing. Add rejection tests for non-synthetic fixture data, malformed JSON, missing media refs, and incident IDs without `FIC-SYN-`. Implement only a small loader that returns typed `ingestion.Packet` values by reusing `ingestion.IngestJSON`; do not bypass validation or duplicate business rules. Verify with the targeted fixture package test, `go test ./internal/ingestion ./internal/eval`, and `go test ./...`.

## Phase 13: Loopback Demo API

Goal: expose the demo review result through a local-only API suitable for a `curl` walkthrough.

- [x] Add `POST /demo/review` for synthetic incident ID or synthetic packet JSON input.
- [x] Return deterministic JSON with review output, approval-required actions, eval summary pointer if available, and trace ID.
- [x] Reject malformed JSON, unknown incident IDs, non-synthetic input, unsupported methods, and unsupported paths.
- [x] Keep the API loopback-only and stateless or in-memory.
- [x] Do not add auth, database, identity, live model calls, or external integrations.
- [x] Add exact run and `curl` commands only after tests and local verification pass.

Output: [Loopback Demo API](../demo/loopback-demo-api.md), the Go package [internal/httpapi](../../../internal/httpapi), and the local server command [cmd/demo-api](../../../cmd/demo-api). Phase 13 is complete as of 2026-05-06 and added a loopback-only stateless demo route for deterministic review JSON. No auth, database, identity, Slack behavior, webhook, live model call, real export, real escalation, external-sharing integration, dashboard, external observability pipeline, or production API behavior exists yet.

Code-task prompt:

> Implement the loopback demo API using strict TDD. First add failing tests around the HTTP handler or local server boundary for `POST /demo/review`: valid synthetic incident ID, valid synthetic packet JSON, malformed JSON, unknown incident ID, non-synthetic input, wrong method, and unknown path. Confirm the tests fail before production code. Implement only enough handler wiring to call the demo review composer and return deterministic JSON. Keep the server local-only and do not add auth, persistence, Slack, webhooks, model providers, or external network calls. Verify with targeted tests, `go test ./...`, and one local `curl` command before documenting the command.

## Phase 14: Dry-Run Slack-Shaped Notification Preview

Goal: show an integration-shaped action while proving external sharing fails closed and no network delivery occurs.

- [x] Generate a Slack-shaped payload from the redacted brief only.
- [x] Require `delivery_mode: "dry_run"` for notification previews.
- [x] Block notification preview as external sharing unless a scoped approval exists.
- [x] Return blocked status, reason, and prepared payload when approval is missing.
- [x] Record a redacted tool-call observability event for preview generation.
- [x] Prove no Slack token, webhook URL, SDK, secret, or network request is used.

Output: [Dry-Run Slack-Shaped Notification Preview](../demo/dry-run-slack-preview.md), the Go package [internal/notification](../../../internal/notification), and the route `POST /demo/notifications/slack` in [internal/httpapi](../../../internal/httpapi). Phase 14 is complete as of 2026-05-06 and added dry-run Slack-shaped payload generation plus a loopback-only blocked preview route. No Slack SDK, token, webhook URL, secret, network request, real Slack delivery, approval retry route, persistent approval store, or external-sharing integration exists yet.

Code-task prompt:

> Implement a dry-run Slack-shaped notification preview using strict TDD. Start with failing tests for payload generation from a redacted brief, mandatory `dry_run` delivery mode, missing approval blocked, denied approval blocked, out-of-scope approval blocked, scoped approval allowed, redacted observability event emission, and proof that no network sender is invoked. Confirm red before production changes. Implement only enough preview logic to return a Slack-shaped payload and approval status. Do not use Slack tokens, webhook URLs, SDKs, environment secrets, or outbound network calls. Run targeted tests and `go test ./...`; document the route or command only after local verification passes.

## Phase 15: Scoped Approval Demo Retry

Goal: make the approval gate visible in the demo by showing blocked and allowed dry-run attempts.

- [ ] Add a local approval request path for one synthetic incident, action, and target channel.
- [ ] Show missing, pending, denied, and out-of-scope approvals fail closed.
- [ ] Show approved dry-run notification preview succeeds only for the exact incident, action, and channel.
- [ ] Preserve append-only in-memory audit history.
- [ ] Keep approvals human-supplied and deterministic.
- [ ] Do not infer approval from model output, notification text, fixture names, or test setup shortcuts.

Planned output: local approval demo route or command integration, plus documentation updates after implementation exists.

Code-task prompt:

> Implement the scoped approval retry demo using strict TDD. Add failing tests for creating an approval request, recording a human decision, retrying a dry-run notification while pending, retrying after denial, retrying with out-of-scope incident/action/channel, and retrying after exact scoped approval. Confirm failures before production changes. Implement only enough API or command wiring to use the existing approval gate semantics and preserve append-only audit history. Sensitive actions must fail closed by default. Run targeted tests and `go test ./...`; update demo docs only after the retry flow is locally proven.

## Phase 16: Eval And Observability Demo Reports

Goal: expose quality and operations proof through local reports.

- [ ] Add a local eval report surface that runs deterministic golden cases and returns case count, metric scores, thresholds, and pass/fail status.
- [ ] Add an in-memory trace report surface that returns redacted events by trace ID.
- [ ] Include a budget-exceeded demo path using caller-supplied token counts.
- [ ] Include eval summary and notification preview tool-call events when available.
- [ ] Keep reports local and ephemeral.
- [ ] Do not imply dashboards, alerts, OpenTelemetry export, persistent logs, provider billing reconciliation, or model benchmarking.

Planned output: local eval and trace report routes or commands, plus documentation updates after implementation exists.

Code-task prompt:

> Implement local eval and observability demo reports using strict TDD. Start with failing tests for eval report generation, threshold pass/fail fields, deterministic case count, trace lookup by trace ID, redacted event fields, missing trace behavior, budget-exceeded event display, and no persistent storage. Confirm red before production changes. Implement only enough report logic to expose existing `internal/eval` and `internal/observability` behavior locally. Do not add dashboards, alerts, OpenTelemetry export, provider billing reconciliation, model calls, or persisted history. Run targeted tests and `go test ./...`; update the one-page eval summary and demo commands only after verification passes.

## Phase 17: Demo Script Refresh

Goal: convert the interview demo from a code/tests walkthrough to a local API walkthrough after the demo surface exists.

- [ ] Update [Demo Package](../demo/demo-package.md) with verified local startup and `curl` commands.
- [ ] Add a fallback path that still works with `go test ./...`.
- [ ] Show one happy-path review, one blocked dry-run notification, one exact scoped approval retry, one eval report, and one trace report.
- [ ] Fill the eval summary with numbers from the latest verified local run.
- [ ] Update `README.md`, `docs/mvp/README.md`, and contributor guides only for commands that exist.
- [ ] Confirm implemented-versus-planned wording is synchronized.

Planned output: refreshed demo package and how-to documentation after Phases 12 through 16 are implemented.

Documentation prompt:

> Refresh the Fleet Incident Copilot demo package after the local demo surface is implemented. Use only verified commands and code-backed behavior. Include a short local startup path, exact `curl` commands, expected response highlights, a fallback `go test ./...` walkthrough, recording script, and implemented-versus-planned wording. Do not describe Slack delivery, production API behavior, persistence, live model calls, dashboards, identity, or external integrations unless those features are implemented and tested.
