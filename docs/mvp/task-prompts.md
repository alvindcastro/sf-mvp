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

> Using `research-report.md`, write or update the Fleet Incident Copilot product frame. Include the target user, problem, MVP promise, success criteria, scope boundaries, trust boundaries, approval gates, and explicit non-goals. Do not write code. Keep claims grounded in the report and avoid saying anything is implemented unless the repo proves it.

## Planning Prompt: Synthetic Incident Specs

> Create or update `docs/mvp/synthetic-incident-packets.md` with Markdown specs for at least five synthetic incident packets. Include low, medium, high, unknown, and adversarial/missing-data cases. For each packet, list required fields, expected timeline outputs, expected severity, expected recommendations, expected brief behavior, and acceptance criteria. Use fake data only and do not create machine-readable fixtures until a future strict-TDD code phase introduces failing tests first.

## Planning Prompt: RAG Corpus

> Design or refine `docs/mvp/rag-corpus-and-grounding.md` for Fleet Incident Copilot. Include document categories, metadata fields, retrieval questions, citation rules, no-match behavior, scope-filtering behavior, and prompt-injection test content. Do not claim runtime behavior unless code and tests already exist.

## Planning Prompt: Agent And Tool Contract

> Define the agent/tool contract conceptually. Include tool names, inputs, outputs, validation rules, approval gates, audit events, and actions the agent must never perform automatically. Treat retrieved content and incident packets as untrusted. Do not write code.

## Planning Prompt: Demo Package

> Create or refine the demo packaging materials for the Fleet Incident Copilot MVP. Include a repo narrative, demo video script, architecture diagram checklist, one-page eval summary outline, and interview talking points mapped to RAG, agents, backend APIs, evals, monitoring, security, cost controls, and production readiness. Distinguish implemented package-level behavior from planned production integrations, and do not claim a feature is implemented unless the current docs, Go packages, and tests prove it.

## Code Prompt: Incident Packet Ingestion

> Implement synthetic incident packet ingestion using strict TDD. First identify the smallest observable behavior: accepting one valid synthetic packet and rejecting one invalid packet. Add failing tests for valid ingestion, missing required fields, malformed telemetry, unsupported event type, and audit event emission. Run the targeted tests and confirm they fail for the expected reason. Then implement only enough parser and validation logic to pass. Refactor only after tests are green. Acceptance requires deterministic validation errors, no real customer data, and a test summary showing red-to-green evidence.

## Code Prompt: RAG Retrieval

> Implement the first RAG retrieval slice using strict TDD. Start with failing tests for relevant SOP retrieval, citation metadata, no-match behavior, scope filtering, and hostile retrieved text that tries to override system instructions. Confirm the tests fail before production code changes. Implement only the retrieval interface needed for the MVP fixture corpus. Acceptance requires cited snippets, deterministic filtering, no unauthorized documents in context, and tests proving retrieved text is data rather than instructions.

## Code Prompt: Timeline Builder

> Implement incident timeline generation using strict TDD. Begin with failing tests for chronological ordering, source citation, uncertainty labeling, missing-data handling, and conflict handling. Run the tests and confirm they fail for the expected reason. Then implement the smallest timeline builder that passes. Acceptance requires every timeline claim to trace to packet data or retrieved source metadata, and tests proving unsupported claims are omitted or marked unknown.

## Code Prompt: Severity And Actions

> Implement or refine severity classification and recommended next actions using strict TDD. Use `docs/mvp/severity-classification-and-recommended-actions.md` as the current Phase 5 behavior contract. First add failing tests for low, medium, high, unknown, conflicting-signal, explanation, and approval-required scenarios. Confirm red before writing production code. Implement deterministic rules first and keep any model-dependent behavior out of the initial path unless it is already covered by tests. Acceptance requires explainable outputs, SOP-grounded recommendations, approval flags for sensitive actions, and tests proving no escalation or export occurs automatically.

## Code Prompt: Incident Brief

> Implement or refine incident brief drafting using strict TDD. Use `docs/mvp/incident-brief-drafting.md` as the current Phase 6 behavior contract. Add failing tests for a complete brief, citation inclusion, sensitive-field redaction, missing-evidence failure, uncertainty labeling, and approval-state display. Confirm the tests fail before production changes. Then implement the smallest drafting layer needed to pass. Acceptance requires no uncited factual claims, no sensitive leakage in shareable output, and tests showing draft creation fails closed when required evidence is absent.

## Code Prompt: Approval Workflow

> Implement or refine the human approval gate using strict TDD. Use `docs/mvp/human-approval-workflow.md` as the current Phase 7 behavior contract. First add or update failing tests for pending approval creation, denied approval, granted approval, blocked action before approval, scoped approval, and immutable audit history. Confirm red before production code. Then implement the smallest workflow that passes. Acceptance requires sensitive actions to fail closed, audit records to be append-only, final decisions to avoid in-place rewrites, and tests proving denied or out-of-scope actions cannot execute.

## Code Prompt: Eval Harness

> Implement or refine the MVP eval harness using strict TDD. Use `docs/mvp/eval-plan.md` as the current Phase 8 behavior contract. Start with failing tests for loading eval cases, scoring expected severity, checking citation coverage, detecting unsupported claims, verifying redaction, checking recommendation accuracy, checking approval fail-closed behavior, and handling prompt-injection fixtures. Confirm failures before implementation. Then implement the smallest evaluator. Acceptance requires repeatable local evals, clear pass/fail thresholds, normal/adversarial/incomplete fixtures, and a test summary showing what went red and green.

## Code Prompt: Observability And Cost

> Implement or refine observability and cost controls using strict TDD. Use `docs/mvp/observability-and-cost-controls.md` as the current Phase 9 behavior contract. Add failing tests for trace propagation, structured event emission, token recording, latency recording, budget-limit behavior, and sensitive-field redaction in logs. Confirm red before production changes. Then implement instrumentation with the smallest surface needed. Acceptance requires structured logs, useful debugging signals, no sensitive evidence leakage, and tests for normal and budget-exceeded paths.
