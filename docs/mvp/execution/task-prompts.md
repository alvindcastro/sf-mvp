# Task Prompts

Use these prompts for future agents. Documentation prompts may edit Markdown only. Code prompts must follow [Strict TDD Rules](tdd-rules.md).

## Reusable Code Task Prompt Template

> You are working in the Fleet Incident Copilot MVP. Follow strict TDD.
>
> Goal: `[specific user-visible or system-visible behavior]`
>
> Constraints:
>
> - Do not write production code until you have added or updated failing tests.
> - Start by reading the relevant existing files and test patterns.
> - Add the narrowest failing test that proves the missing behavior.
> - Run the targeted test and confirm it fails for the expected reason.
> - Implement only enough production code to pass.
> - Run the targeted test again and confirm green.
> - Refactor only if it improves clarity without changing behavior.
> - Run the relevant broader test suite.
> - Update documentation only if behavior, commands, architecture, or acceptance criteria changed.
>
> Required coverage:
>
> - Happy path.
> - Invalid or incomplete input.
> - Security or authorization boundary, if applicable.
> - Observability, audit, or eval signal, if applicable.
>
> Acceptance checks:
>
> - Tests fail before implementation and pass after implementation.
> - No untested production behavior is introduced.
> - No secrets, real customer data, real evidence data, or external service assumptions are added.
> - Output includes files changed, tests run, behavior added, and residual risk.

## Planning Prompt: Product Frame

> Using `docs/research/research-report.md`, write or update the Fleet Incident Copilot product frame. Include the target user, problem, MVP promise, success criteria, scope boundaries, trust boundaries, approval gates, and explicit non-goals. Do not write code. Keep claims grounded in the report and avoid saying anything is implemented unless the repo proves it.

## Planning Prompt: Synthetic Incident Specs

> Create or update `docs/mvp/workflow/synthetic-incident-packets.md` with Markdown specs for at least five synthetic incident packets. Include low, medium, high, unknown, and adversarial/missing-data cases. For each packet, list required fields, expected timeline outputs, expected severity, expected recommendations, expected brief behavior, and acceptance criteria. Use fake data only. Machine-readable fixture changes belong in the Phase 12 fixture contract and must keep strict-TDD failing tests first.

## Planning Prompt: RAG Corpus

> Design or refine `docs/mvp/workflow/rag-corpus-and-grounding.md` for Fleet Incident Copilot. Include document categories, metadata fields, retrieval questions, citation rules, no-match behavior, scope-filtering behavior, and prompt-injection test content. Do not claim runtime behavior unless code and tests already exist.

## Planning Prompt: Agent And Tool Contract

> Define the agent/tool contract conceptually. Include tool names, inputs, outputs, validation rules, approval gates, audit events, and actions the agent must never perform automatically. Treat retrieved content and incident packets as untrusted. Do not write code.

## Planning Prompt: Demo Package

> Create or refine the demo packaging materials for the Fleet Incident Copilot MVP. Include a repo narrative, demo video script, architecture diagram checklist, one-page eval summary outline, and interview talking points mapped to RAG, agents, backend APIs, evals, monitoring, security, cost controls, and production readiness. Distinguish implemented package-level behavior from planned production integrations, and do not claim a feature is implemented unless the current docs, Go packages, and tests prove it.

## Planning Prompt: Hiring-Manager Demo Surface

> Brainstorm a concrete hiring-manager demo surface for Fleet Incident Copilot without writing code. Compare a loopback-only local API, CLI fallback, dry-run Slack-shaped notification preview, scoped approval retry, local eval report, and in-memory observability proof. Create or update Markdown files with tickable phase tasks and strict-TDD future task prompts. Keep real Slack delivery, webhooks, secrets, persistence, identity, dashboards, live model calls, and production compliance claims out of scope unless future code and tests prove them.

## Code Prompt: Incident Packet Ingestion

> Implement synthetic incident packet ingestion using strict TDD. First identify the smallest observable behavior: accepting one valid synthetic packet and rejecting one invalid packet. Add failing tests for valid ingestion, missing required fields, malformed telemetry, unsupported event type, and audit event emission. Run the targeted tests and confirm they fail for the expected reason. Then implement only enough parser and validation logic to pass. Refactor only after tests are green. Acceptance requires deterministic validation errors, no real customer data, and a test summary showing red-to-green evidence.

## Code Prompt: RAG Retrieval

> Implement the first RAG retrieval slice using strict TDD. Start with failing tests for relevant SOP retrieval, citation metadata, no-match behavior, scope filtering, and hostile retrieved text that tries to override system instructions. Confirm the tests fail before production code changes. Implement only the retrieval interface needed for the MVP fixture corpus. Acceptance requires cited snippets, deterministic filtering, no unauthorized documents in context, and tests proving retrieved text is data rather than instructions.

## Code Prompt: Timeline Builder

> Implement incident timeline generation using strict TDD. Begin with failing tests for chronological ordering, source citation, uncertainty labeling, missing-data handling, and conflict handling. Run the tests and confirm they fail for the expected reason. Then implement the smallest timeline builder that passes. Acceptance requires every timeline claim to trace to packet data or retrieved source metadata, and tests proving unsupported claims are omitted or marked unknown.

## Code Prompt: Severity And Actions

> Implement or refine severity classification and recommended next actions using strict TDD. Use `docs/mvp/workflow/severity-classification-and-recommended-actions.md` as the current Phase 5 behavior contract. First add failing tests for low, medium, high, unknown, conflicting-signal, explanation, and approval-required scenarios. Confirm red before writing production code. Implement deterministic rules first and keep any model-dependent behavior out of the initial path unless it is already covered by tests. Acceptance requires explainable outputs, SOP-grounded recommendations, approval flags for sensitive actions, and tests proving no escalation or export occurs automatically.

## Code Prompt: Incident Brief

> Implement or refine incident brief drafting using strict TDD. Use `docs/mvp/workflow/incident-brief-drafting.md` as the current Phase 6 behavior contract. Add failing tests for a complete brief, citation inclusion, sensitive-field redaction, missing-evidence failure, uncertainty labeling, and approval-state display. Confirm the tests fail before production changes. Then implement the smallest drafting layer needed to pass. Acceptance requires no uncited factual claims, no sensitive leakage in shareable output, and tests showing draft creation fails closed when required evidence is absent.

## Code Prompt: Approval Workflow

> Implement or refine the human approval gate using strict TDD. Use `docs/mvp/workflow/human-approval-workflow.md` as the current Phase 7 behavior contract. First add or update failing tests for pending approval creation, denied approval, granted approval, blocked action before approval, scoped approval, and immutable audit history. Confirm red before production code. Then implement the smallest workflow that passes. Acceptance requires sensitive actions to fail closed, audit records to be append-only, final decisions to avoid in-place rewrites, and tests proving denied or out-of-scope actions cannot execute.

## Code Prompt: Eval Harness

> Implement or refine the MVP eval harness using strict TDD. Use `docs/mvp/quality/eval-plan.md` as the current Phase 8 behavior contract. Start with failing tests for loading eval cases, scoring expected severity, checking citation coverage, detecting unsupported claims, verifying redaction, checking recommendation accuracy, checking approval fail-closed behavior, and handling prompt-injection fixtures. Confirm failures before implementation. Then implement the smallest evaluator. Acceptance requires repeatable local evals, clear pass/fail thresholds, normal/adversarial/incomplete fixtures, and a test summary showing what went red and green.

## Code Prompt: Observability And Cost

> Implement or refine observability and cost controls using strict TDD. Use `docs/mvp/quality/observability-and-cost-controls.md` as the current Phase 9 behavior contract. Add failing tests for trace propagation, structured event emission, token recording, latency recording, budget-limit behavior, and sensitive-field redaction in logs. Confirm red before production changes. Then implement instrumentation with the smallest surface needed. Acceptance requires structured logs, useful debugging signals, no sensitive evidence leakage, and tests for normal and budget-exceeded paths.

## Code Prompt: Demo Review Composer

> Add or refine a demo review composition contract using strict TDD. Use `docs/mvp/demo/review-composition-contract.md` as the current Phase 12 behavior contract. Do not add an HTTP server, Slack behavior, persistence, live model call, or external service in this task.
>
> Ownership suggestion: add a small `internal/demo` package for composition and fixture-facing helpers. Reuse `internal/ingestion`, `internal/retrieval`, `internal/timeline`, `internal/severity`, `internal/brief`, `internal/approval`, and `internal/observability`. Do not move business rules into demo glue code.
>
> Red:
>
> - Add `internal/demo/review_test.go`.
> - First failing test: `TestComposeReviewReturnsSeverityBriefApprovalAndTrace` or the smallest missing composer behavior.
> - Run `go test ./internal/demo` and confirm failure is missing composer behavior or the specific regression under test.
> - Add focused tests for unknown incident ID, non-synthetic input, missing evidence, citation preservation, redaction preservation, approval-required action display, and no external action execution.
>
> Green:
>
> - Implement only enough in-memory composition logic to pass.
> - Use existing deterministic package APIs and mock guidance.
> - Return validation errors from ingestion without drafting a brief.
> - Include approval state but do not approve, export, escalate, notify, or share.
>
> Verify:
>
> - Run `go test ./internal/demo`.
> - Run `go test ./internal/ingestion ./internal/retrieval ./internal/timeline ./internal/severity ./internal/brief ./internal/approval ./internal/observability`.
> - Run `go test ./...`.
> - Update docs whenever composer behavior changes.

## Code Prompt: Machine-Readable Demo Fixtures

> Add or refine synthetic machine-readable demo fixtures using strict TDD. The goal is to make current golden cases usable by demo adapters without weakening validation or introducing real data.
>
> Ownership suggestion: add fixture loading to `internal/demo` or `internal/fixtures`. Keep eval scoring in `internal/eval`. If fixture files are added, place them under `testdata/demo` or `docs/mvp/demo/fixtures` and load them through ingestion validation.
>
> Red:
>
> - Add `internal/demo/fixtures_test.go`, `internal/demo/review_test.go`, or `internal/fixtures/fixtures_test.go`.
> - First failing test: `TestLoadDemoFixturesReturnsSyntheticNormalIncompleteAndAdversarialPackets` or the smallest missing fixture-loader behavior.
> - Run the targeted package test and confirm failure is missing fixture loader behavior or the specific regression under test.
> - Add tests for rejecting non-synthetic fixture data, malformed fixture JSON, missing media refs, and incident IDs without the `FIC-SYN-` prefix.
>
> Green:
>
> - Implement a small loader that returns typed `ingestion.Packet` values by reusing `ingestion.IngestJSON`.
> - Keep fixture names stable: low hard brake, medium stop-arm conflict, high collision signal, unknown incomplete evidence, and adversarial prompt-injection note.
> - Do not bypass ingestion validation or duplicate validation rules.
>
> Verify:
>
> - Run the targeted fixture package test.
> - Run `go test ./internal/ingestion ./internal/eval`.
> - Run `go test ./...`.
> - Confirm docs still say fixtures are synthetic-only.

## Code Prompt: Loopback Demo API

> Add or refine the local HTTP API demo surface using strict TDD. Phase 13 is implemented in `internal/httpapi` and `cmd/demo-api`; do not add persistence, auth, live model calls, Slack, webhooks, or external integrations.
>
> Ownership suggestion: add `internal/httpapi` for handlers tested with `httptest`; add `cmd/demo-api` only after handler tests are green, as thin local wiring.
>
> Red:
>
> - Add `internal/httpapi/httpapi_test.go`.
> - First failing test: `TestReviewIncidentEndpointReturnsBriefSeverityApprovalAndTrace`.
> - Run `go test ./internal/httpapi` and confirm failure is missing package or handler behavior.
> - Add focused tests for malformed JSON, non-synthetic packet rejection, unknown incident ID, unsupported method, unknown path, and no external action execution.
>
> Green:
>
> - Implement only the handler needed to pass `POST /demo/review` or the documented route chosen in the phase contract.
> - Use the demo review composer, in-memory mock guidance, and existing deterministic package APIs.
> - Return ingestion validation errors without drafting a brief.
> - Include approval state but do not approve, export, escalate, notify, or share.
>
> Verify:
>
> - Run `go test ./internal/httpapi`.
> - Run `go test ./cmd/demo-api`.
> - Run `go test ./...`.
> - Run `go vet ./...`.
> - Verify one local `curl` command only after the server wiring exists.
> - Update docs only after the API surface is real.

## Code Prompt: Dry-Run Slack Preview

> Add or refine the dry-run Slack-shaped notification preview using strict TDD. Phase 14 is implemented in `internal/notification` and exposed through the loopback-only `POST /demo/notifications/slack` route; this must not call Slack or require network access.
>
> Ownership suggestion: add `internal/notification` for dry-run message formatting and send simulation. Use `internal/approval` for gating, `internal/brief` for redacted brief content, and `internal/observability` for a tool-call event.
>
> Red:
>
> - Add or update `internal/notification/slack_test.go`.
> - First failing test: `TestDryRunSlackPreviewBlocksWithoutApprovedExternalSharingScope` or the smallest missing preview behavior.
> - Run `go test ./internal/notification` and confirm targeted failure.
> - Add tests for approved dry-run preview, denied approval, mismatched incident scope, mismatched channel target, no sensitive text leakage, mandatory `dry_run`, and no network side effects.
>
> Green:
>
> - Implement a `SlackDryRunSender` or equivalent preview adapter that accepts a channel, target ref, brief result, approval gate, and recorder.
> - Gate with the existing external-sharing sensitive action semantics.
> - Return a preview result with `dry_run: true`; never perform HTTP.
> - Record success or failure as an observability tool call without message body leakage.
>
> Verify:
>
> - Run `go test ./internal/notification`.
> - Run `go test ./internal/approval ./internal/brief ./internal/observability`.
> - Run `go test ./internal/httpapi` if the local route changes.
> - Run `go test ./...`.
> - Document the route or command only after local verification passes.

## Code Prompt: Scoped Approval Retry Demo

> Add or refine approval-gated retry behavior for sensitive demo actions using strict TDD. Phase 15 is implemented in `internal/httpapi` using `internal/approval` and `internal/notification`; keep retry state in memory and loopback-only unless future scope explicitly changes.
>
> Ownership suggestion: prefer `internal/httpapi` when the behavior is demo route orchestration, and keep `approval.Gate` as the source of truth for request, decision, execute, and audit semantics. Prefer `internal/approval` only when the core gate semantics themselves need to change.
>
> Red:
>
> - Add or update `internal/httpapi/httpapi_test.go` for route behavior, or `internal/approval/approval_test.go` for pure gate behavior.
> - First failing handler test: `TestScopedApprovalRetryAllowsExactApprovedDryRunOnly` or the smallest missing retry behavior.
> - Run `go test ./internal/httpapi -run TestScopedApprovalRetryAllowsExactApprovedDryRunOnly` and confirm failure is missing route or retry behavior.
> - Add tests for missing approval, pending retry, denied retry, approved exact retry, out-of-channel retry, out-of-incident retry, wrong-action approval, audit history, and no network delivery.
>
> Green:
>
> - Implement the smallest API needed, such as route wiring that creates an approval request, records a human decision, and passes the shared gate into the dry-run notification preview.
> - Preserve existing final-decision immutability behavior.
> - Append audit events; never edit prior decision records.
> - Do not infer approval from model output, notification payload text, fixture names, or test setup shortcuts.
>
> Verify:
>
> - Run the targeted handler or approval test.
> - Run `go test ./internal/httpapi`.
> - Run `go test ./internal/approval`.
> - Run `go test ./internal/notification`.
> - Run `go test ./...`.
> - Update demo docs only after the retry flow is locally proven.

## Code Prompt: Eval And Trace Reports

> Add a local eval report surface using strict TDD. The report should expose the existing deterministic eval harness as Markdown or JSON for demo review. It must not call a model provider, read real data, or weaken strict thresholds.
>
> Ownership suggestion: add report rendering to `internal/eval`; add a CLI or API route only as thin wiring after renderer tests pass.
>
> Red:
>
> - Add `internal/eval/report_test.go`.
> - First failing test: `TestRenderMarkdownReportIncludesSummaryThresholdsFailuresAndSafetySignals`.
> - Run `go test ./internal/eval -run TestRenderMarkdownReportIncludesSummaryThresholdsFailuresAndSafetySignals` and confirm failure.
> - Add tests for failed report output, adversarial case labels, approval fail-closed status, and no sensitive fixture terms in report text.
>
> Green:
>
> - Implement `RenderMarkdown(report Report) string` or `RenderJSON(report Report)`.
> - Include case count, severity accuracy, citation coverage, recommendation accuracy, pass/fail state, thresholds, and per-case failure summaries.
> - Keep sensitive packet fields out of the report.
> - Add CLI or API wiring only after internal rendering is green.
>
> Verify:
>
> - Run `go test ./internal/eval`.
> - Run adapter package tests if a CLI or API route is added.
> - Run `go test ./...` and `go vet ./...`.
> - Update the one-page eval summary and demo commands only after verification passes.

## Code Prompt: Demo Observability Events

> Extend observability for demo surfaces using strict TDD. The goal is to trace API, fixture, dry-run notification, approval retry, and eval-report actions without logging sensitive evidence.
>
> Ownership suggestion: update `internal/observability` for new event types and redaction rules. Demo adapters should call observability APIs but not define their own event schema.
>
> Red:
>
> - Add tests to `internal/observability/observability_test.go`.
> - First failing test: `TestRecordDemoSurfaceEventTracksSurfaceActionStatusAndLatency`.
> - Run `go test ./internal/observability -run TestRecordDemoSurfaceEventTracksSurfaceActionStatusAndLatency` and confirm targeted failure.
> - Add tests for API request event, fixture load event, dry-run notification event, approval retry event, eval report event, and redaction of channel, message, and packet-derived sensitive fields.
>
> Green:
>
> - Add a small structured event API such as `RecordDemoSurface(workflow, DemoSurfaceCall)`.
> - Include fields for `surface`, `action`, `status`, `dry_run`, `request_id` when applicable, and duration.
> - Do not log packet body, brief body, raw Slack text, transcript notes, or GPS-like coordinates.
>
> Verify:
>
> - Run `go test ./internal/observability`.
> - Run adapter package tests once those packages exist.
> - Run `go test ./...`.

## Documentation Prompt: Local Demo Script Refresh

> Refresh the Fleet Incident Copilot demo package after the local demo surface is implemented. Use only verified commands and code-backed behavior. Include a short local startup path, exact `curl` commands, expected response highlights, a fallback `go test ./...` walkthrough, recording script, and implemented-versus-planned wording. Do not describe Slack delivery, production API behavior, persistence, live model calls, dashboards, identity, or external integrations unless those features are implemented and tested.
