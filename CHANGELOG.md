# Changelog

## 2026-05-06 - Phase 16 Eval And Observability Demo Reports

### Task: Audit Phase 16 Scope With Parallel Agents

- What: Ran parallel read-only Codex explorer agents to inspect the Phase 16 acceptance criteria, existing eval and observability APIs, HTTP route patterns, demo command constraints, and stale documentation wording.
- Where: [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/demo/demo-surface-roadmap.md](docs/mvp/demo/demo-surface-roadmap.md), [docs/mvp/demo/loopback-demo-api.md](docs/mvp/demo/loopback-demo-api.md), [internal/eval](internal/eval), [internal/observability](internal/observability), [internal/httpapi](internal/httpapi), [cmd/demo-api](cmd/demo-api).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 16 should expose deterministic local quality and operations proof without dashboards, alerts, OpenTelemetry export, provider billing reconciliation, model calls, model benchmarking, persisted history, or external integrations.
- How: Spawned two read-only explorer agents while the main thread inspected the canonical phase tracker and existing handler code; aligned route naming with the roadmap’s `GET /demo/eval/latest` and `GET /demo/traces/{trace_id}` pattern and scoped the caller-supplied budget route as `POST /demo/budget/check`.

### Task: Add Phase 16 Report TDD Coverage

- What: Added failing-first handler tests for deterministic eval report generation, case count, metric scores, threshold pass/fail gates, trace lookup by trace ID, notification preview tool-call event display, budget-exceeded event display, redacted event fields, missing trace behavior, and per-handler ephemeral trace state.
- Where: [internal/httpapi/httpapi_test.go](internal/httpapi/httpapi_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Define the Phase 16 local report contract before production code and prove the missing routes at the HTTP boundary.
- How: Added `httptest` coverage around `GET /demo/eval/latest`, `GET /demo/traces/{trace_id}`, and `POST /demo/budget/check`; observed `go test ./internal/httpapi` fail with `404` for the missing report routes before implementing the handler behavior.

### Task: Implement Local Eval, Trace, And Budget Report Routes

- What: Added loopback-only report routes, an in-memory trace event store on the handler, eval report response projection, trace report response projection, caller-supplied budget check response projection, and deterministic error mapping for missing traces and invalid budget inputs.
- Where: [internal/httpapi/httpapi.go](internal/httpapi/httpapi.go), [internal/httpapi/httpapi_test.go](internal/httpapi/httpapi_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Complete Phase 16 by exposing existing `internal/eval` and `internal/observability` behavior locally while keeping report state ephemeral and avoiding new service, persistence, model, billing, or monitoring claims.
- How: Wired `eval.Run(eval.GoldenCases(), eval.DefaultThresholds())` into `GET /demo/eval/latest`, recorded `eval.score_recorded` through `observability.Recorder`, stored cloned observability events by trace ID, returned trace reports only for events recorded by the same handler process, and used `observability.RecordModelCall` with caller-supplied token counts for `POST /demo/budget/check`; verified `go test ./internal/httpapi` passed.

### Task: Document Phase 16 Implemented Boundary

- What: Added the Phase 16 behavior doc and synchronized phase tracking, compatibility phase pointer, demo roadmap, loopback API docs, eval plan, observability docs, demo package, overview, docs indexes, contributor guides, troubleshooting, testing docs, and root README.
- Where: [docs/mvp/demo/eval-and-observability-reports.md](docs/mvp/demo/eval-and-observability-reports.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/phases.md](docs/mvp/phases.md), [docs/mvp/demo/demo-surface-roadmap.md](docs/mvp/demo/demo-surface-roadmap.md), [docs/mvp/demo/loopback-demo-api.md](docs/mvp/demo/loopback-demo-api.md), [docs/mvp/demo/demo-package.md](docs/mvp/demo/demo-package.md), [docs/mvp/quality/eval-plan.md](docs/mvp/quality/eval-plan.md), [docs/mvp/quality/observability-and-cost-controls.md](docs/mvp/quality/observability-and-cost-controls.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/README.md](docs/README.md), [docs/how-tos.md](docs/how-tos.md), [docs/developer-guide.md](docs/developer-guide.md), [docs/nice-to-knows.md](docs/nice-to-knows.md), [docs/testing.md](docs/testing.md), [docs/troubleshooting.md](docs/troubleshooting.md), [README.md](README.md), [CHANGELOG.md](CHANGELOG.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep implemented-versus-planned wording accurate now that local eval, trace, and budget report routes exist, while preserving explicit limits around dashboards, alerts, OpenTelemetry export, persistent logs, provider billing reconciliation, model calls, model benchmarking, databases, identity, auth, Slack delivery, webhooks, production monitoring, production audit storage, and real external sharing.
- How: Marked Phase 16 complete, linked the new behavior doc, documented route request and response contracts, updated latest eval summary numbers, added exact local `curl` examples, updated runtime summaries and test coverage descriptions, and changed stale “planned report surface” language to loopback-only implemented report language.

### Task: Verify Phase 16 Eval And Observability Reports

- What: Verified the report handler, related eval and observability packages, full Go suite, vet checks, coverage, Markdown whitespace, and local loopback `curl` sequence for eval, trace, and budget report routes.
- Where: [internal/httpapi](internal/httpapi), [internal/eval](internal/eval), [internal/observability](internal/observability), [cmd/demo-api](cmd/demo-api), [docs/mvp/demo/eval-and-observability-reports.md](docs/mvp/demo/eval-and-observability-reports.md), [CHANGELOG.md](CHANGELOG.md).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 16 is deterministic, locally repeatable, redacted, ephemeral, and still free of provider calls, persistence, billing reconciliation, dashboards, alerts, and external telemetry before closing the task.
- How: Ran `go test ./internal/httpapi`, `go test ./internal/eval ./internal/observability ./internal/httpapi`, `go test ./...`, `go vet ./...`, `go test -cover ./...`, and a scoped `git diff --check` over the Phase 16 edit set; noted the broader `git diff --check -- README.md CHANGELOG.md docs internal` is blocked by pre-existing trailing whitespace in `docs/research/research-report.md`; started `go run ./cmd/demo-api -addr 127.0.0.1:18083`; verified `GET /demo/eval/latest` returned `case_count: 5`, `passed: true`, and perfect metric scores; verified `GET /demo/traces/trace-fic-syn-eval-report-20260506t160000z-001` returned `workflow.started` and `eval.score_recorded`; verified `POST /demo/budget/check` returned `budget_report.status: "budget_exceeded"` with `total_tokens: 110`; then stopped the local server.

## 2026-05-06 - Phase 15 Scoped Approval Demo Retry

### Task: Audit Phase 15 Scope With Parallel Agents

- What: Ran parallel read-only Codex explorer agents to inspect the Phase 15 acceptance criteria, existing approval and notification semantics, likely HTTP route shape, and stale documentation wording.
- Where: [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/demo/demo-surface-roadmap.md](docs/mvp/demo/demo-surface-roadmap.md), [docs/mvp/demo/dry-run-slack-preview.md](docs/mvp/demo/dry-run-slack-preview.md), [docs/mvp/demo/loopback-demo-api.md](docs/mvp/demo/loopback-demo-api.md), [internal/approval](internal/approval), [internal/notification](internal/notification), [internal/httpapi](internal/httpapi).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 15 should make the existing in-memory approval gate visible through the loopback demo without adding identity, persistence, real Slack delivery, webhooks, secrets, auth, production API behavior, or real external sharing.
- How: Spawned two read-only explorer agents while the main thread read the phase tracker and handler code; used their findings to scope `POST /demo/approvals`, `POST /demo/approvals/decisions`, shared in-memory gate state, and the documentation update set.

### Task: Add Phase 15 Scoped Approval Retry TDD Coverage

- What: Added failing-first handler tests for creating an approval request, retrying a dry-run notification while pending, retrying after denial, retrying after exact scoped approval, blocking another channel, blocking another incident, blocking a wrong-action approval, preserving audit history, and proving allowed dry-run retry still performs no send or network delivery.
- Where: [internal/httpapi/httpapi_test.go](internal/httpapi/httpapi_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Define the Phase 15 HTTP-boundary contract before production code and preserve fail-closed approval semantics at the local demo surface.
- How: Added `httptest` coverage around `POST /demo/approvals`, `POST /demo/approvals/decisions`, and repeated `POST /demo/notifications/slack` calls; observed `go test ./internal/httpapi` fail with `404` for the missing approval routes after correcting the test decoder helper.

### Task: Implement In-Memory Approval Retry Routes

- What: Added shared in-memory approval gate state to the loopback handler, approval request and decision routes, approval and audit response projections, approval error mapping, and notification retry wiring that reuses the same gate.
- Where: [internal/httpapi/httpapi.go](internal/httpapi/httpapi.go), [internal/httpapi/httpapi_test.go](internal/httpapi/httpapi_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Complete the Phase 15 retry behavior while keeping `internal/approval` as the source of truth for request, decision, execute, scope matching, and append-only audit semantics.
- How: Added a mutex-protected `approval.Gate` to `Handler`, mapped approval target channels to `notification.SlackTargetRef`, validated known synthetic incidents through the existing demo composer, mapped final-decision rewrites to `409`, and passed the shared gate into `notification.PreviewSlack`; verified `go test ./internal/httpapi` passed.

### Task: Document Phase 15 Implemented Boundary

- What: Added the Phase 15 behavior doc and synchronized phase tracking, demo roadmap, loopback API docs, dry-run notification docs, demo package, overview, contributor guides, troubleshooting, docs indexes, compatibility phase pointer, task prompts, and root README.
- Where: [docs/mvp/demo/scoped-approval-retry.md](docs/mvp/demo/scoped-approval-retry.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/phases.md](docs/mvp/phases.md), [docs/mvp/demo/demo-surface-roadmap.md](docs/mvp/demo/demo-surface-roadmap.md), [docs/mvp/demo/loopback-demo-api.md](docs/mvp/demo/loopback-demo-api.md), [docs/mvp/demo/dry-run-slack-preview.md](docs/mvp/demo/dry-run-slack-preview.md), [docs/mvp/demo/demo-package.md](docs/mvp/demo/demo-package.md), [docs/mvp/demo/review-composition-contract.md](docs/mvp/demo/review-composition-contract.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md), [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md), [docs/README.md](docs/README.md), [docs/how-tos.md](docs/how-tos.md), [docs/developer-guide.md](docs/developer-guide.md), [docs/nice-to-knows.md](docs/nice-to-knows.md), [docs/testing.md](docs/testing.md), [docs/troubleshooting.md](docs/troubleshooting.md), [README.md](README.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep implemented-versus-planned wording accurate now that local approval retry exists, while preserving explicit limits around persistent approval stores, identity, roles, auth, real Slack delivery, webhooks, secrets, live external sharing, databases, production APIs, and production audit storage.
- How: Marked Phase 15 complete, linked the new behavior doc, added route request and response contracts, added exact local `curl` examples, updated runtime summaries and test coverage descriptions, and changed stale “approval retry planned” language to in-memory implemented language.

### Task: Verify Phase 15 Scoped Approval Retry

- What: Verified the handler, approval and notification integration, full Go suite, vet checks, coverage, and local `curl` retry sequence against the loopback server.
- Where: [internal/httpapi](internal/httpapi), [internal/approval](internal/approval), [internal/notification](internal/notification), [cmd/demo-api](cmd/demo-api), [docs/mvp/demo/scoped-approval-retry.md](docs/mvp/demo/scoped-approval-retry.md), [CHANGELOG.md](CHANGELOG.md).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 15 is deterministic, locally repeatable, exact-scope, and still dry-run/no-network before closing the task.
- How: Ran `go test ./internal/httpapi`, `go test ./internal/approval ./internal/notification ./internal/httpapi ./cmd/demo-api`, `go test ./...`, `go vet ./...`, and `go test -cover ./...`; started `go run ./cmd/demo-api -addr 127.0.0.1:18082`; verified missing approval returned `notification_preview.status: "blocked"`, created `approval-001` through `POST /demo/approvals`, verified pending retry stayed blocked with `approval_request_id: "approval-001"`, approved it through `POST /demo/approvals/decisions`, verified retry returned `notification_preview.status: "allowed"` with `sent: false` and `network_delivery_attempted: false`, then stopped the local server.

## 2026-05-06 - Phase 14 Dry-Run Slack-Shaped Notification Preview

### Task: Audit Phase 14 Scope With Parallel Agents

- What: Ran parallel read-only Codex explorer agents to inspect the Phase 14 acceptance criteria, likely package boundaries, approval and observability contracts, local API options, and documentation update set.
- Where: [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md), [internal/approval](internal/approval), [internal/brief](internal/brief), [internal/observability](internal/observability), [internal/demo](internal/demo), [internal/httpapi](internal/httpapi).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 14 should add a dry-run Slack-shaped notification preview without Slack tokens, webhook URLs, SDKs, secrets, outbound network calls, real delivery, approval retry wiring, or external-sharing integration.
- How: Spawned two read-only explorer agents while the main thread inspected the phase plan and package APIs; resolved the requested `docs/mvp/phases.md` path ambiguity by using [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md) as the canonical tracker and adding a compatibility pointer.

### Task: Add Phase 14 Notification Preview TDD Coverage

- What: Added failing-first tests for blocked dry-run preview without approval, mandatory `delivery_mode: "dry_run"`, denied approval, out-of-scope approval, exact scoped approval, redacted tool-call observability events, prepared Slack-shaped payloads, and no network delivery.
- Where: [internal/notification/slack_test.go](internal/notification/slack_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Define the Phase 14 behavior contract before production code and prove missing notification-preview behavior through strict TDD.
- How: Created a new `internal/notification` test surface around `PreviewSlack`, `SlackTargetRef`, dry-run result fields, approval gate semantics, and observability event assertions; ran `go test ./internal/notification`; observed the expected red build failure for missing preview types and functions.

### Task: Implement Dry-Run Slack-Shaped Preview Package

- What: Added `internal/notification` with dry-run-only preview generation, Slack-shaped payload structs, external-sharing approval gating, blocked and allowed preview statuses, redacted observability tool-call events, and explicit no-send/no-network result fields.
- Where: [internal/notification/slack.go](internal/notification/slack.go), [internal/notification/slack_test.go](internal/notification/slack_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Complete the package-level Phase 14 preview behavior while reusing the existing approval, brief, severity, and observability contracts instead of introducing a Slack integration.
- How: Implemented `PreviewSlack`, `DeliveryModeDryRun`, `SlackPayload`, `PreviewResult`, `SlackTargetRef`, redacted-brief payload projection, `external_sharing` gate calls with `slack:<channel>` target scope, and `notification.slack.preview` tool-call recording; verified `go test ./internal/notification` passed.

### Task: Add Loopback Notification Preview Route

- What: Added `POST /demo/notifications/slack` to the local demo handler and response projection for blocked dry-run notification previews.
- Where: [internal/httpapi/httpapi.go](internal/httpapi/httpapi.go), [internal/httpapi/httpapi_test.go](internal/httpapi/httpapi_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Make Phase 14 visible through the existing loopback demo server while keeping approval creation and retry behavior deferred to Phase 15.
- How: Added failing handler tests for the missing route and non-dry-run rejection; observed `go test ./internal/httpapi` return `404` for the new route; implemented route dispatch, request parsing, review composition, redacted brief conversion, empty in-memory approval gate usage, notification preview projection, and `dry_run_required` error mapping.

### Task: Document Phase 14 Implemented Boundary

- What: Added the Phase 14 behavior doc and synchronized phase tracking, demo roadmap, loopback API docs, demo package, overview, contributor guides, testing docs, troubleshooting docs, docs indexes, and root README.
- Where: [docs/mvp/demo/dry-run-slack-preview.md](docs/mvp/demo/dry-run-slack-preview.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/phases.md](docs/mvp/phases.md), [docs/mvp/demo/demo-surface-roadmap.md](docs/mvp/demo/demo-surface-roadmap.md), [docs/mvp/demo/loopback-demo-api.md](docs/mvp/demo/loopback-demo-api.md), [docs/mvp/demo/demo-package.md](docs/mvp/demo/demo-package.md), [docs/mvp/demo/review-composition-contract.md](docs/mvp/demo/review-composition-contract.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md), [docs/README.md](docs/README.md), [docs/how-tos.md](docs/how-tos.md), [docs/developer-guide.md](docs/developer-guide.md), [docs/nice-to-knows.md](docs/nice-to-knows.md), [docs/testing.md](docs/testing.md), [docs/troubleshooting.md](docs/troubleshooting.md), [README.md](README.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep implemented-versus-planned wording accurate after adding the dry-run notification package and route, while preserving explicit limits around real Slack delivery, webhooks, tokens, SDKs, secrets, outbound network calls, approval retry, persistence, auth, production APIs, and external-sharing integrations.
- How: Marked Phase 14 complete, linked the new package and route, documented the verified `curl` command, updated targeted test lists, described the blocked-before-approval route behavior, and kept Phase 15 approval retry plus Phase 16 report surfaces planned.

### Task: Verify Phase 14 Notification Preview

- What: Verified the new notification package, local handler, related package set, full Go suite, vet checks, coverage, and one local `curl` request against the loopback notification preview route.
- Where: [internal/notification](internal/notification), [internal/httpapi](internal/httpapi), [cmd/demo-api](cmd/demo-api), [docs/mvp/demo/dry-run-slack-preview.md](docs/mvp/demo/dry-run-slack-preview.md), [CHANGELOG.md](CHANGELOG.md).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 14 is deterministic, locally repeatable, and does not regress earlier package contracts before closing the task.
- How: Ran `go test ./internal/notification`, `go test ./internal/approval ./internal/brief ./internal/observability ./internal/notification ./internal/demo ./internal/httpapi`, `go test ./internal/httpapi`, `go test ./...`, `go vet ./...`, and `go test -cover ./...`; started `go run ./cmd/demo-api -addr 127.0.0.1:18081`; verified `curl -i --max-time 5 -X POST http://127.0.0.1:18081/demo/notifications/slack -H "Content-Type: application/json" -d '{"incident_id":"FIC-SYN-001","channel":"#fleet-safety","delivery_mode":"dry_run"}'` returned HTTP `200` with `notification_preview.status: "blocked"`, `sent: false`, and `network_delivery_attempted: false`; stopped the local server.

## 2026-05-06 - Phase 13 Loopback Demo API

### Task: Audit Phase 13 Scope With Parallel Agents

- What: Ran parallel read-only Codex explorer agents to inspect Phase 13 acceptance criteria, related docs, existing demo composer APIs, likely handler and server wiring, guardrails, and verification commands.
- Where: [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md), [docs/mvp/demo/demo-surface-roadmap.md](docs/mvp/demo/demo-surface-roadmap.md), [docs/mvp/demo/review-composition-contract.md](docs/mvp/demo/review-composition-contract.md), [internal/demo](internal/demo).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 13 should wrap the existing deterministic demo composer with a loopback-only local API without adding auth, persistence, Slack, webhooks, live model calls, export, escalation, or external sharing.
- How: Spawned two read-only explorer agents while the main thread inspected the phase plan, demo docs, and package APIs; used their findings to scope `internal/httpapi`, `cmd/demo-api`, and the documentation update set.

### Task: Add Phase 13 Handler TDD Coverage

- What: Added failing-first handler tests for `POST /demo/review` with a known synthetic incident ID, full synthetic packet JSON, malformed JSON, unknown incident ID, non-synthetic packet input, unsupported method, unknown path, and loopback default address.
- Where: [internal/httpapi/httpapi_test.go](internal/httpapi/httpapi_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Define the loopback API behavior at the HTTP boundary before adding production handler code and preserve the synthetic-only, deterministic, no-external-action contract.
- How: Created `httptest` coverage around `NewHandler` and `DefaultListenAddress`; ran `go test ./internal/httpapi`; observed the expected red build failure for missing handler and address APIs.

### Task: Implement The Loopback Review Handler

- What: Added an `internal/httpapi` package with `POST /demo/review`, deterministic JSON response projection, eval summary pointer, approval-required action display, loopback default address, and stable error mapping for malformed JSON, unknown IDs, non-synthetic input, missing evidence, unsupported methods, and unknown paths.
- Where: [internal/httpapi/httpapi.go](internal/httpapi/httpapi.go).
- When: 2026-05-06, America/Vancouver.
- Why: Complete the Phase 13 handler surface by exposing the existing demo review composer through a local HTTP boundary without duplicating ingestion validation or broadening business behavior.
- How: Parsed request JSON, routed packet-shaped requests through `ingestion.IngestJSON`, routed incident-ID requests through `demo.ComposeIncident`, mapped `demo` sentinel errors to JSON status responses, projected review fields into lower_snake_case JSON, and kept approval-required actions blocked.

### Task: Add Thin Local Server Wiring

- What: Added `cmd/demo-api` to start the loopback review server, keep the default bind address at `127.0.0.1:8080`, allow loopback-only address overrides, and reject non-loopback binds such as `0.0.0.0:8080`.
- Where: [cmd/demo-api/main.go](cmd/demo-api/main.go), [cmd/demo-api/main_test.go](cmd/demo-api/main_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Provide a real local process for the verified `curl` walkthrough while preserving the Phase 13 local-only guardrail.
- How: Added a failing server-wiring test for `newServer`, implemented the minimal `http.Server` setup with `httpapi.NewHandler`, added `-addr` flag parsing, validated loopback IP addresses, and used the override after port `8080` was already occupied locally.

### Task: Document Phase 13 Implemented Boundary

- What: Added the Phase 13 loopback API contract and synchronized phase tracking, demo roadmap, demo package, overview, contributor guides, testing docs, troubleshooting docs, and root README.
- Where: [docs/mvp/demo/loopback-demo-api.md](docs/mvp/demo/loopback-demo-api.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/demo/demo-surface-roadmap.md](docs/mvp/demo/demo-surface-roadmap.md), [docs/mvp/demo/demo-package.md](docs/mvp/demo/demo-package.md), [docs/mvp/demo/review-composition-contract.md](docs/mvp/demo/review-composition-contract.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md), [docs/README.md](docs/README.md), [docs/how-tos.md](docs/how-tos.md), [docs/developer-guide.md](docs/developer-guide.md), [docs/nice-to-knows.md](docs/nice-to-knows.md), [docs/testing.md](docs/testing.md), [docs/troubleshooting.md](docs/troubleshooting.md), [README.md](README.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep implemented-versus-planned wording accurate now that a loopback review API exists, while preserving explicit limits around production APIs, auth, database, persistence, Slack, webhooks, live providers, export, escalation, external sharing, dashboards, and external observability.
- How: Marked Phase 13 complete, added verified `go run` and `curl` commands, linked `internal/httpapi` and `cmd/demo-api`, changed stale “no HTTP API” wording where it applied to current repo state, and left Phase 14 through Phase 17 behavior documented as planned.

### Task: Verify Phase 13 Loopback API

- What: Verified the new handler, server command, related package set, full Go suite, vet checks, and one local `curl` request against the loopback review endpoint.
- Where: [internal/httpapi](internal/httpapi), [cmd/demo-api](cmd/demo-api), [internal/demo](internal/demo), [internal/ingestion](internal/ingestion), [docs/mvp/demo/loopback-demo-api.md](docs/mvp/demo/loopback-demo-api.md), [CHANGELOG.md](CHANGELOG.md).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 13 is locally repeatable, deterministic, and does not regress earlier package contracts before documenting runnable commands.
- How: Ran `go test ./internal/httpapi`, `go test ./cmd/demo-api`, `go test ./internal/demo ./internal/ingestion ./internal/httpapi ./cmd/demo-api`, `go test ./...`, `go vet ./...`, and `go test -cover ./...`; started `go run ./cmd/demo-api -addr 127.0.0.1:18080` because `127.0.0.1:8080` was occupied; verified `curl -i --max-time 5 -X POST http://127.0.0.1:18080/demo/review -H "Content-Type: application/json" -d '{"incident_id":"FIC-SYN-001"}'` returned HTTP `200` with the deterministic trace ID and blocked approval actions; stopped the local server; and ran whitespace checks on the Phase 13 tracked and new-file edit set.

## 2026-05-06 - Phase 12 Review Composition Contract

### Task: Audit Phase 12 Scope With Parallel Agents

- What: Ran parallel read-only Codex explorer agents to inspect Phase 12 acceptance criteria, related docs, existing package contracts, likely implementation surface, and verification commands.
- Where: [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md), [internal](internal), [README.md](README.md).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm the new work belonged in a deterministic package-level demo composer and avoid adding HTTP, Slack, persistence, live model calls, or external integrations outside Phase 12.
- How: Spawned two read-only explorer agents while the main thread inspected the phase plan and package APIs; used their findings to scope `internal/demo`, fixture loading, and documentation updates.

### Task: Add Phase 12 Demo Composer TDD Coverage

- What: Added failing-first tests for loading default synthetic demo fixtures, rejecting malformed and non-synthetic fixture data, composing a known synthetic incident, rejecting unknown incident IDs, rejecting non-synthetic typed input before downstream composition, failing closed on missing evidence, preserving retrieved citations, preserving redactions, and displaying approval-required actions.
- Where: [internal/demo/review_test.go](internal/demo/review_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Define the Phase 12 behavior contract before production code and prove the new composer surface was missing for the expected reason.
- How: Created `internal/demo` tests around `LoadDefaultFixtures`, `LoadFixtures`, `ComposeIncident`, and `ComposePacket`; ran `go test ./internal/demo`; observed the expected red build failure for missing fixture loader, composer, response, and status types.

### Task: Implement In-Memory Demo Fixtures And Review Composer

- What: Added an `internal/demo` package with machine-readable synthetic fixture loading, an embedded JSON fixture set for `FIC-SYN-001` through `FIC-SYN-005`, deterministic review composition, non-synthetic rejection, missing-evidence fail-closed errors, retrieved citation refs, timeline entries, severity, recommendations, redacted brief fields, approval-required action display, and observability events.
- Where: [internal/demo/review.go](internal/demo/review.go), [internal/demo/testdata/demo-fixtures.json](internal/demo/testdata/demo-fixtures.json).
- When: 2026-05-06, America/Vancouver.
- Why: Complete Phase 12 by composing the existing package-level workflow into one review result without moving business rules out of ingestion, retrieval, timeline, severity, brief, approval, or observability packages.
- How: Implemented fixture decoding through `ingestion.IngestJSON`, reused deterministic mock guidance retrieval, called `timeline.Build`, `severity.Classify`, `brief.Draft`, `approval.Gate.Execute`, and `observability.Recorder`, then projected the existing package outputs into a small `ReviewResult`.

### Task: Document Phase 12 Implemented Boundary

- What: Added the Phase 12 review composition behavior contract and synchronized the phase tracker, demo roadmap, demo package, overview, fixture docs, contributor docs, testing docs, troubleshooting docs, and root README.
- Where: [docs/mvp/demo/review-composition-contract.md](docs/mvp/demo/review-composition-contract.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/demo/demo-surface-roadmap.md](docs/mvp/demo/demo-surface-roadmap.md), [docs/mvp/demo/demo-package.md](docs/mvp/demo/demo-package.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/README.md](docs/README.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md), [docs/mvp/workflow/synthetic-incident-packets.md](docs/mvp/workflow/synthetic-incident-packets.md), [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md), [docs/how-tos.md](docs/how-tos.md), [docs/developer-guide.md](docs/developer-guide.md), [docs/nice-to-knows.md](docs/nice-to-knows.md), [docs/testing.md](docs/testing.md), [docs/troubleshooting.md](docs/troubleshooting.md), [README.md](README.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep implemented-versus-planned wording accurate now that the package-level composer and machine-readable fixtures exist, while preserving explicit limits around API, CLI, Slack, persistence, live providers, export, escalation, and external sharing.
- How: Marked Phase 12 complete, added the new behavior doc, linked `internal/demo`, updated fixture notes from future to implemented, added targeted test commands, and kept Phase 13 through Phase 17 surfaces documented as planned.

### Task: Verify Phase 12 Review Composition

- What: Verified the new demo package, related composed package set, full Go suite, vet checks, package coverage, and Markdown or code whitespace for the Phase 12 edit set.
- Where: [internal/demo](internal/demo), [internal/ingestion](internal/ingestion), [internal/retrieval](internal/retrieval), [internal/timeline](internal/timeline), [internal/severity](internal/severity), [internal/brief](internal/brief), [internal/approval](internal/approval), [internal/observability](internal/observability), [docs/mvp/demo/review-composition-contract.md](docs/mvp/demo/review-composition-contract.md), [CHANGELOG.md](CHANGELOG.md).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 12 is locally repeatable, does not regress earlier package contracts, and keeps the documentation sync reviewable.
- How: Ran `go test ./internal/demo`, `go test ./internal/ingestion ./internal/retrieval ./internal/timeline ./internal/severity ./internal/brief ./internal/approval ./internal/observability ./internal/demo`, `go test ./...`, `go vet ./...`, and `go test -cover ./...` successfully. Ran `git diff --check` on the tracked Phase 12 edit set and `git diff --check --no-index /dev/null ...` on new files; the new-file checks produced no whitespace warnings and returned nonzero only because the files are untracked additions.

## 2026-05-06 - Phase 11 Demo Surface Roadmap

### Task: Plan Hiring-Manager Demo Surfaces

- What: Added a documentation-only roadmap for future local demo surfaces: loopback review API, dry-run Slack-shaped notification preview, scoped approval retry, local eval report, and redacted trace report.
- Why: Give the project a concrete hiring-manager demo path beyond `go test ./...` while preserving the current boundary that no API, CLI, Slack delivery, webhook, database, live model call, or external integration exists yet.
- Where: [docs/mvp/demo/demo-surface-roadmap.md](docs/mvp/demo/demo-surface-roadmap.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md), [docs/mvp/execution/tdd-rules.md](docs/mvp/execution/tdd-rules.md), [docs/mvp/demo/demo-package.md](docs/mvp/demo/demo-package.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/README.md](docs/README.md), [README.md](README.md).
- Validation: ran `git diff --check -- README.md CHANGELOG.md docs/mvp/README.md docs/README.md docs/mvp/demo/demo-package.md docs/mvp/demo/demo-surface-roadmap.md docs/mvp/execution/phases.md docs/mvp/execution/task-prompts.md docs/mvp/execution/tdd-rules.md` and `go test ./...`.

## 2026-05-06 - Contributor Guides

### Task: Add Local Developer Documentation Set

- What: Added standalone contributor guides for how-tos, developer workflow, nice-to-know repo context, testing commands, and troubleshooting.
- Why: Give future local work a direct entry point for common commands, package boundaries, synthetic-data guardrails, strict-TDD expectations, and common failure modes without overstating the current package-level runtime surface.
- Where: [docs/how-tos.md](docs/how-tos.md), [docs/developer-guide.md](docs/developer-guide.md), [docs/nice-to-knows.md](docs/nice-to-knows.md), [docs/testing.md](docs/testing.md), [docs/troubleshooting.md](docs/troubleshooting.md), [docs/README.md](docs/README.md), [README.md](README.md).
- Validation: ran `git diff --check -- CHANGELOG.md README.md docs/README.md docs/how-tos.md docs/developer-guide.md docs/nice-to-knows.md docs/testing.md docs/troubleshooting.md` and `go test ./...`.

## 2026-05-06 - Documentation Organization

### Task: Group Markdown Artifacts By Purpose

- What: Moved the source research report under `docs/research`, grouped MVP docs under `overview`, `workflow`, `quality`, `execution`, and `demo`, and added a top-level documentation index.
- Why: Make the Markdown set easier to scan and maintain without changing the implemented runtime surface.
- Where: [docs/README.md](docs/README.md), [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp](docs/mvp), [docs/research/research-report.md](docs/research/research-report.md).
- Validation: checked relative Markdown links, ran `git diff --check -- README.md CHANGELOG.md docs`, and ran `go test ./...`.

## 2026-05-06 - Phase 10 Demo Package

### Task: Audit Implemented Surface For Demo Truthfulness

- What: Ran a parallel Codex explorer audit of the current MVP implementation surface and overclaim risks before completing Phase 10 packaging.
- Where: [README.md](README.md), [docs/mvp](docs/mvp), [internal](internal).
- When: 2026-05-06, America/Vancouver.
- Why: Ensure demo materials describe only implemented package-level behavior and do not imply a CLI, HTTP API, database, live model call, vector database, real export, real escalation, external sharing, identity, dashboards, production compliance, or real customer evidence processing.
- How: Spawned a read-only explorer agent to inspect the repo docs and internal packages while the main thread continued the Phase 10 implementation path; used its findings to keep the packaging language focused on deterministic Go package APIs and planned production surfaces.

### Task: Create Phase 10 Demo Packaging Artifact

- What: Replaced the placeholder Phase 10 checklist with usable demo packaging materials: implemented-versus-planned wording rules, repo narrative, fleet safety operator mapping, short demo script, architecture diagram checklist, one-page eval summary outline, interview talking points, and final packaging checklist.
- Where: [docs/mvp/demo/demo-package.md](docs/mvp/demo/demo-package.md).
- When: 2026-05-06, America/Vancouver.
- Why: Complete Phase 10 with materials that communicate production-readiness thinking for RAG, constrained agent boundaries, Go package APIs, evals, monitoring, security, cost controls, and known limits without overstating implementation.
- How: Spawned a parallel worker agent to draft the artifact, reviewed the resulting diff, kept current behavior described as package-level and in-memory, and marked live providers, vector search, real tools, persistence, identity, dashboards, and production audit/compliance guarantees as planned rather than implemented.

### Task: Synchronize Phase 10 Documentation State

- What: Marked Phase 10 complete, updated the MVP overview and scope completion language, tightened the reusable demo-package prompt, and noted the new Phase 10 packaging surface in the root README.
- Where: [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md), [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md), [README.md](README.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, artifact map, definition of done, future-agent prompt, and repository overview aligned after adding the final demo packaging artifact.
- How: Checked all Phase 10 tasks, added an output note, changed the artifact map from placeholder checklist language to finished packaging-material language, marked the under-five-minute demo and package deliverable done, and preserved explicit limits around end-to-end runtime behavior and production integrations.

### Task: Verify Phase 10 Documentation Changes

- What: Verified the documentation-only Phase 10 change did not break the existing Go package suite.
- Where: [docs/mvp/demo/demo-package.md](docs/mvp/demo/demo-package.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [CHANGELOG.md](CHANGELOG.md), [internal](internal).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 10 packaging is repo-consistent and that no runtime package regressed while the docs were updated.
- How: Ran `go test ./...` successfully after the documentation sync. Ran `git diff --check`; the full-worktree check reported trailing whitespace only in pre-existing dirty files outside the Phase 10 edit set (`.idea/.gitignore` and `docs/research/research-report.md`), so the Phase 10 files were checked directly instead of changing unrelated work.

## 2026-05-06 - Phase 9 Observability And Cost Controls

### Task: Add Phase 9 Observability TDD Coverage

- What: Added failing-first observability tests for trace ID creation, structured workflow events, retrieval counts and source IDs, tool-call success and redaction, approval decision metadata, token and latency recording, budget-exceeded behavior, eval score recording, cache candidates, and model-routing notes.
- Where: [internal/observability/observability_test.go](internal/observability/observability_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 9 with strict TDD and make quality, safety, latency, and spend signals observable before adding production recorder behavior.
- How: Created same-package tests around `NewRecorder`, `StartWorkflow`, retrieval, tool-call, approval, model-call, eval, and cost-plan APIs; ran `go test ./internal/observability`; observed the expected red build failure for missing observability APIs and types before adding implementation.

### Task: Implement In-Memory Observability Recorder And Cost Plan

- What: Added an observability package that creates deterministic trace IDs, records in-memory structured events for workflow start, retrieval, tool calls, approval decisions, model calls, budget failures, and eval summaries, redacts configured sensitive terms and coordinate-like strings, and exposes cache candidates plus hosted, smaller-model, and self-hosted routing notes.
- Where: [internal/observability/observability.go](internal/observability/observability.go).
- When: 2026-05-06, America/Vancouver.
- Why: Provide the first package-level Phase 9 observability and cost-control surface without adding an external telemetry backend, dashboards, alerts, persistent log store, provider billing reconciliation, live model-provider calls, cache storage, CLI, HTTP API, database behavior, real export, real escalation, or external-sharing integrations.
- How: Implemented `NewRecorder`, `StartWorkflow`, `Events`, `RecordRetrieval`, `RecordToolCall`, `RecordApprovalDecision`, `RecordModelCall`, `RecordEvalScore`, `Budget`, `SensitiveData`, event models, defensive event copying, redaction helpers, cumulative token accounting, budget-limit checks, and `DefaultCostPlan`; then confirmed `go test ./internal/observability` passed.

### Task: Tighten Phase 9 Regression Coverage

- What: Added regression assertions for retrieval snippet non-logging, approval request ID, approver, and latency recording, private location label redaction, input, output, and total token budget failures, and invalid negative token usage rejection.
- Where: [internal/observability/observability_test.go](internal/observability/observability_test.go), [internal/observability/observability.go](internal/observability/observability.go).
- When: 2026-05-06, America/Vancouver.
- Why: Close Phase 9 assertion gaps and prevent invalid caller-supplied token counts from reducing cumulative spend accounting.
- How: Added a failing test for negative token usage, observed `go test ./internal/observability` fail for missing `ErrInvalidTokenUsage` and `model_call.rejected`, implemented the smallest rejection path before budget accounting, and expanded existing tests for the documented metadata, redaction, and budget dimensions.

### Task: Document Phase 9 Observability Contract

- What: Added the Phase 9 behavior contract, runtime surface, structured event model, recorded signals, redaction rules, budget controls, cache candidates, model-routing notes, current limits, test commands, and red-to-green evidence.
- Where: [docs/mvp/quality/observability-and-cost-controls.md](docs/mvp/quality/observability-and-cost-controls.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, future-agent prompt, and observability behavior documentation synchronized after adding Phase 9 runtime behavior and the invalid-token guard.
- How: Created the dedicated Phase 9 artifact, checked off the Phase 9 tracker, linked the new Go package, changed the reusable observability prompt to “implement or refine,” documented `ErrInvalidTokenUsage` and `model_call.rejected`, and kept external telemetry, persistent logs, live providers, billing reconciliation, live routing, cache storage, CLI, HTTP API, and database behavior out of scope.

### Task: Update Repository State Documentation

- What: Updated repository overview and scope language to include the Phase 9 observability package, targeted test command, package-level structured events, caller-supplied token usage, invalid token usage, budget-limit failures, cache candidates, and model-routing notes.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Remove stale statements that observability, token tracking, and cost controls were unimplemented while preserving the limits around external telemetry, persistence, live providers, billing, real integrations, CLI, HTTP API, and database behavior.
- How: Added the Phase 9 artifact and package to documentation maps, marked observability and cost-control thinking as implemented, recorded the Phase 9 trust boundary, and clarified that the current surface is in-memory and package-level only.

### Task: Verify Phase 9 Observability Surface

- What: Verified the targeted observability package, full Go test suite, vet checks, and package coverage after implementation and documentation sync.
- Where: [internal/observability](internal/observability), [docs/mvp/quality/observability-and-cost-controls.md](docs/mvp/quality/observability-and-cost-controls.md), [CHANGELOG.md](CHANGELOG.md).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 9 is locally repeatable and does not break earlier MVP phase packages.
- How: Ran `go test ./internal/observability`, `go test ./...`, `go vet ./...`, and `go test -cover ./...`; all completed successfully, with `internal/observability` reporting 94.3% statement coverage.

## 2026-05-06 - Phase 8 Eval Harness

### Task: Add Phase 8 Eval TDD Coverage

- What: Added failing-first eval harness tests for loading golden normal, adversarial, and incomplete cases; scoring expected severity; scoring citation coverage; detecting unsupported claims; checking SOP-grounded recommendation accuracy; verifying redaction; checking prompt-injection resistance; and asserting strict release thresholds.
- Where: [internal/eval/eval_test.go](internal/eval/eval_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 8 with strict TDD and make deterministic eval behavior observable before adding production eval logic.
- How: Created same-package tests around `GoldenCases`, `Run`, and `DefaultThresholds`; ran `go test ./internal/eval`; observed the expected red build failure for missing eval APIs and types before adding implementation.

### Task: Implement Deterministic In-Memory Eval Harness

- What: Added an eval package that loads synthetic golden cases, retrieves mock guidance, composes timeline, severity, brief, and approval-gate behavior, and reports severity accuracy, citation coverage, recommendation accuracy, unsupported claims, redaction leaks, prompt-injection resistance, and approval fail-closed status.
- Where: [internal/eval/eval.go](internal/eval/eval.go).
- When: 2026-05-06, America/Vancouver.
- Why: Provide repeatable local release gates for the MVP demo path without adding a CLI, HTTP API, database, persistent eval storage, model calls, observability, token tracking, cost controls, export, escalation, or external-sharing integrations.
- How: Implemented `GoldenCases`, `Run`, `DefaultThresholds`, strict release-gate fields, in-memory synthetic packet fixtures for `FIC-SYN-001` through `FIC-SYN-005`, an approved mock guidance corpus, citation and recommendation scoring, prohibited-claim and sensitive-term scans, prompt-injection checks, and no-approval sensitive-action callback checks; then confirmed `go test ./internal/eval` passed.

### Task: Document Phase 8 Eval Contract

- What: Reworked the eval plan into the Phase 8 behavior contract, checked off Phase 8, documented golden cases, scoring rules, release thresholds, current limits, test commands, and red-to-green evidence.
- Where: [docs/mvp/quality/eval-plan.md](docs/mvp/quality/eval-plan.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, future-agent prompt, and eval behavior documentation synchronized after adding runtime eval harness behavior.
- How: Linked the new Go package from the Phase 8 tracker, changed the reusable eval prompt to “implement or refine,” recorded strict default thresholds, and kept observability, token and cost metrics, CLI reporting, HTTP API, persistence, and real integrations out of Phase 8 scope.

### Task: Update Repository State Documentation

- What: Updated repository overview and scope language to include the Phase 8 eval package and targeted test command while keeping observability and cost controls deferred.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Remove stale statements that said no eval harness existed and split eval completion from future observability work.
- How: Added the eval artifact and package to documentation maps, marked deterministic local evals as implemented, added the Phase 8 trust boundary, and clarified that trace IDs, structured logs, latency, token usage, budget metrics, CLI, HTTP API, database, and real integrations remain unimplemented.

### Task: Verify Phase 8 Eval Harness

- What: Verified the targeted eval package, full Go test suite, vet checks, and package coverage after the implementation and documentation sync.
- Where: [internal/eval](internal/eval), [README.md](README.md), [docs/mvp/quality/eval-plan.md](docs/mvp/quality/eval-plan.md).
- When: 2026-05-06, America/Vancouver.
- Why: Confirm Phase 8 is locally repeatable and does not break earlier MVP phase packages.
- How: Ran `go test ./internal/eval`, `go test ./...`, `go vet ./...`, and `go test -cover ./...`; all completed successfully, with `internal/eval` reporting 86.0% statement coverage.

## 2026-05-06 - Phase 7 Human Approval Workflow

### Task: Add Phase 7 Approval TDD Coverage

- What: Added failing-first approval workflow tests for pending request creation, approver and decision capture, blocked pending actions, denied action blocking, approved scoped execution, out-of-scope blocking, mismatched call-and-scope incident blocking, append-only audit-copy behavior, and final-decision immutability.
- Where: [internal/approval/approval_test.go](internal/approval/approval_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 7 with strict TDD and make the human approval gate observable before adding production approval logic.
- How: Created same-package tests around `export`, `escalation`, and `external_sharing` using the existing `severity.SensitiveAction` labels; ran `go test ./internal/approval`; observed the expected red build failure for missing `NewGate`, request, scope, decision, audit, and gate types, then added focused red tests for immutable final decisions and mismatched call-and-scope incidents before tightening enforcement.

### Task: Implement In-Memory Human Approval Gate

- What: Added an approval package that creates in-memory approval requests, records approved or denied human decisions, blocks missing, pending, denied, out-of-scope, and mismatched call-and-scope sensitive action callbacks, allows approved callbacks only within exact scope, prevents final decisions from being rewritten in place, and returns append-only audit history copies.
- Where: [internal/approval/approval.go](internal/approval/approval.go).
- When: 2026-05-06, America/Vancouver.
- Why: Provide the first Phase 7 enforcement boundary without adding real export, escalation, external-sharing, persistence, identity, CLI, HTTP API, or external-service behavior.
- How: Implemented `NewGate`, `CreateRequest`, `Decide`, `Execute`, and `AuditHistory` with deterministic request IDs, injected clock timestamps, `DecisionPending`, `DecisionApproved`, `DecisionDenied`, `ErrActionBlocked`, `ErrRequestAlreadyDecided`, scoped request matching, and audit events for requested, decided, blocked, and allowed actions; then confirmed `go test ./internal/approval` passed.

### Task: Document Phase 7 Approval Contract

- What: Added the Phase 7 approval workflow contract, runtime surface, approval request model, enforcement rules, audit history behavior, current limits, test command, and red-to-green evidence.
- Where: [docs/mvp/workflow/human-approval-workflow.md](docs/mvp/workflow/human-approval-workflow.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/quality/eval-plan.md](docs/mvp/quality/eval-plan.md), [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, eval planning, future-agent prompt, and behavior documentation synchronized after adding approval workflow runtime behavior.
- How: Created the dedicated Phase 7 artifact, checked off the Phase 7 tracker, linked the new Go package, added Phase 7 eval mappings for pending, denied, approved, scoped, and immutable-audit cases, and pointed the reusable approval prompt at the new contract.

### Task: Update Repository State Documentation

- What: Updated repository overview and scope language to include the Phase 7 approval package and targeted test command without claiming persistence, identity, real export, real escalation, external sharing, observability, eval harness, CLI, HTTP API, or database behavior exists.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md), [docs/mvp/overview/product-frame.md](docs/mvp/overview/product-frame.md), [docs/mvp/workflow/synthetic-incident-packets.md](docs/mvp/workflow/synthetic-incident-packets.md), [docs/mvp/workflow/severity-classification-and-recommended-actions.md](docs/mvp/workflow/severity-classification-and-recommended-actions.md), [docs/mvp/workflow/incident-brief-drafting.md](docs/mvp/workflow/incident-brief-drafting.md).
- When: 2026-05-06, America/Vancouver.
- Why: Remove stale statements that human approval records and enforcement were future work while preserving the remaining limits around actual sensitive-action integrations and production runtime surfaces.
- How: Added the Phase 7 artifact to documentation maps, marked human approval gating as implemented, recorded the fail-closed approval trust boundary, clarified Phase 5 approval flags and Phase 6 approval-state display as separate from Phase 7 decision records, and kept real export, escalation, and external sharing out of implemented scope.

## 2026-05-06 - Phase 6 Shareable Incident Brief

### Task: Add Phase 6 Brief TDD Coverage

- What: Added failing-first brief tests for complete draft output, citation coverage, sensitive-field redaction, missing-evidence fail-closed behavior, uncertainty labeling, approval-state display, and unsupported-claim omission.
- Where: [internal/brief/brief_test.go](internal/brief/brief_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 6 with strict TDD and make cited, redacted, human-reviewable brief behavior observable before adding production drafting logic.
- How: Created synthetic in-test packet fixtures using the existing ingestion, retrieval, timeline, and severity contracts; ran `go test ./internal/brief`; observed the expected red build failure for missing `Draft`, `StatusDraft`, result, section, redaction, approval-state, and fail-closed error types.

### Task: Implement Structured Brief Drafting

- What: Added a brief package that drafts structured incident summaries, cited timelines, severity rationale, recommended actions, and approval-state sections; carries source references; redacts vehicle, route, location, GPS-label, sensitive transcript, sensitive still-frame, and coordinate-like text; labels uncertainty; and fails closed when required evidence is missing.
- Where: [internal/brief/brief.go](internal/brief/brief.go).
- When: 2026-05-06, America/Vancouver.
- Why: Provide the first Phase 6 shareable draft surface without rendering documents, approving actions, exporting evidence, escalating incidents, or sharing externally.
- How: Implemented `Draft(packet ingestion.Packet, timelineResult timeline.Result, severityResult severity.Result) (Result, error)` with deterministic section assembly, source conversion, required-evidence validation, `MissingEvidenceError`, structured redaction tracking, blocked approval-state display, and `StatusDraft`; then confirmed `go test ./internal/brief` passed.

### Task: Document Phase 6 Brief Contract

- What: Added the Phase 6 brief drafting contract, runtime surface, section list, citation rules, redaction rules, approval-state behavior, current limits, test command, and red-to-green evidence.
- Where: [docs/mvp/workflow/incident-brief-drafting.md](docs/mvp/workflow/incident-brief-drafting.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/quality/eval-plan.md](docs/mvp/quality/eval-plan.md), [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, eval planning, future-agent prompts, and behavior documentation synchronized after adding brief drafting runtime behavior.
- How: Created the dedicated Phase 6 artifact, checked off the Phase 6 tracker, linked the new Go package, added Phase 6 eval mappings for citations, redaction, fail-closed behavior, uncertainty, and approval-state display, and pointed the reusable brief prompt at the new contract.

### Task: Update Repository State Documentation

- What: Updated repository overview and scope language to include the Phase 6 brief package and targeted test command without claiming a renderer, approval workflow, eval harness, export, escalation, external sharing, CLI, HTTP API, database, or persistence exists.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Avoid stale statements that said no brief generator existed while preserving the remaining workflow boundaries and later-phase limits.
- How: Added the Phase 6 artifact to documentation maps, marked cited and redacted draft brief behavior as implemented, recorded the shareable-output redaction trust boundary, and kept human approval workflow and sensitive-action execution out of implemented scope.

## 2026-05-06 - Phase 5 Severity Classification And Recommended Actions

### Task: Add Phase 5 Severity TDD Coverage

- What: Added failing-first severity tests for low hard-brake classification, medium stop-arm classification, high collision classification, unknown incomplete-evidence classification, conflicting telemetry handling, recommendation explanations, SOP citation grounding, adversarial transcript handling, deterministic no-model behavior, and approval-required flags.
- Where: [internal/severity/severity_test.go](internal/severity/severity_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 5 with strict TDD and make severity, recommendation, and approval-flag behavior observable before adding production classifier logic.
- How: Created synthetic in-test packet fixtures using the existing ingestion, retrieval, and timeline contracts; ran `go test ./internal/severity`; observed the expected red build failure for missing `Classify`, result, recommendation, source, and approval types.

### Task: Implement Deterministic Severity And Recommendation Rules

- What: Added a severity package that classifies low, medium, high, and unknown severity with deterministic rules, treats conflicting timeline telemetry as unknown, treats adversarial transcript content as untrusted data, returns explained recommendations with packet and retrieved-guidance sources, records `ModelJudgmentUsed: false`, and marks export, escalation, and external sharing as approval-required but not approved.
- Where: [internal/severity/severity.go](internal/severity/severity.go).
- When: 2026-05-06, America/Vancouver.
- Why: Provide the first explainable Phase 5 reasoning surface without calling a model, executing sensitive actions, or introducing approval workflow behavior before Phase 7.
- How: Implemented `Classify(packet ingestion.Packet, timelineResult timeline.Result, guidance retrieval.Result) Result` with event-type rules, timeline-conflict detection, `retrieved_data` citation filtering, recommendation action labels, source deduplication, and fail-closed sensitive-action approval flags; then confirmed `go test ./internal/severity` passed.

### Task: Document Phase 5 Severity Contract

- What: Added the Phase 5 severity and recommendation contract, runtime surface, deterministic rule table, recommendation output shape, approval-required flag behavior, current limits, test command, and red-to-green evidence.
- Where: [docs/mvp/workflow/severity-classification-and-recommended-actions.md](docs/mvp/workflow/severity-classification-and-recommended-actions.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/quality/eval-plan.md](docs/mvp/quality/eval-plan.md), [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, eval planning, future-agent prompts, and behavior documentation synchronized after adding severity runtime behavior.
- How: Created the dedicated Phase 5 artifact, checked off the Phase 5 tracker, linked the new Go package, added Phase 5 eval mappings for severity and recommendations, and pointed the reusable severity prompt at the new contract.

### Task: Update Repository State Documentation

- What: Updated repository overview and scope language to include the Phase 5 severity package and targeted test command without claiming a brief generator, approval workflow, export, escalation, external sharing, observability, or eval harness exists.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Avoid stale statements that said no severity classifier existed while preserving the remaining workflow boundaries and later-phase limits.
- How: Added the Phase 5 artifact to documentation maps, marked severity and recommendations as implemented, recorded the deterministic severity trust boundary, and kept human approval workflow and sensitive-action execution out of implemented scope.

## 2026-05-06 - Phase 4 Incident Timeline Builder

### Task: Add Phase 4 Timeline TDD Coverage

- What: Added failing-first timeline tests for chronological telemetry ordering, source citation coverage, guidance citation carry-forward, transcript and still-frame attribution, missing-evidence uncertainty, conflicting telemetry uncertainty, and unsupported-claim omission.
- Where: [internal/timeline/timeline_test.go](internal/timeline/timeline_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 4 with strict TDD and make grounded timeline behavior observable before adding production timeline logic.
- How: Created synthetic packet fixtures using the existing ingestion and retrieval contracts, ran `go test ./internal/timeline`, and observed the expected red build failure for missing `Build`, `Entry`, and timeline result types.

### Task: Implement Deterministic Incident Timeline Builder

- What: Added a timeline package that builds chronological entries from validated packet telemetry, labels transcript and still-frame claims by source type, preserves structured packet source references, carries approved retrieval citation references as guidance metadata, marks unavailable evidence as uncertain, and labels conflicting same-time telemetry.
- Where: [internal/timeline/timeline.go](internal/timeline/timeline.go).
- When: 2026-05-06, America/Vancouver.
- Why: Provide the first brief-ready grounded timeline surface without inventing facts, approving sensitive actions, or treating retrieved text as instructions.
- How: Implemented `Build(packet ingestion.Packet, guidance retrieval.Result) Result` with deterministic in-memory ordering, `packet.*[N]` source references, `retrieved_data` citation filtering, unavailable-evidence detection, and conflict labeling; then confirmed `go test ./internal/timeline` passed.

### Task: Document Phase 4 Timeline Contract

- What: Added the Phase 4 timeline builder contract, output format, input dependencies, source-reference rules, uncertainty behavior, unsupported-claim behavior, current limits, test command, and red-to-green evidence.
- Where: [docs/mvp/workflow/incident-timeline-builder.md](docs/mvp/workflow/incident-timeline-builder.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/quality/eval-plan.md](docs/mvp/quality/eval-plan.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, eval planning, and behavior documentation synchronized after adding timeline runtime behavior.
- How: Created the dedicated Phase 4 artifact, checked off the Phase 4 tracker, linked the new Go package, and added future eval checks for ordering, citation coverage, uncertainty, conflicts, and unsupported claims.

### Task: Update Repository State Documentation

- What: Updated repository overview language to include the Phase 4 timeline package and its targeted test command without claiming severity, brief, approval, export, escalation, or external sharing behavior.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Avoid stale statements that said no timeline builder existed while preserving the remaining unimplemented workflow boundaries.
- How: Added the Phase 4 artifact to documentation maps, marked cited timeline generation as implemented, recorded the new timeline trust boundary, and verified `go test ./internal/timeline`, `go test ./internal/ingestion ./internal/retrieval ./internal/timeline`, and `go test ./...` passed.

## 2026-05-06 - Phase 3 RAG Corpus And Grounding

### Task: Add Phase 3 Retrieval TDD Coverage

- What: Added failing-first retrieval tests for relevant SOP retrieval, citation metadata, no-match behavior, workflow and scope filtering, hostile retrieved text, and deterministic result limiting.
- Where: [internal/retrieval/retrieval_test.go](internal/retrieval/retrieval_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 3 with strict TDD and make the grounding boundary observable before adding production retrieval logic.
- How: Created synthetic in-test mock SOP and troubleshooting documents, ran `go test ./internal/retrieval`, and observed the expected red build failure for missing retrieval APIs and types.

### Task: Implement Deterministic Mock Guidance Retrieval

- What: Added an in-memory retrieval package with document/query/result types, exact workflow and scope filtering, lexical scoring, stable citation references, snippets, and `retrieved_data` content-role marking.
- Where: [internal/retrieval/retrieval.go](internal/retrieval/retrieval.go).
- When: 2026-05-06, America/Vancouver.
- Why: Allow future reasoning phases to retrieve only approved mock guidance with citations while keeping retrieved text untrusted.
- How: Implemented `NewRetriever` and `Retrieve` with fail-closed empty results for empty workflow or scope, filtering before ranking, deterministic ordering by score and source ID, and citation references formatted as `SOURCE_ID#YYYY-MM-DD`.

### Task: Document Phase 3 Corpus And Grounding Contract

- What: Added the Phase 3 mock RAG corpus contract, document metadata, mock SOPs, troubleshooting notes, citation format, no-match behavior, scope filtering behavior, prompt-injection fixture, eval questions, red-to-green evidence, and current limits.
- Where: [docs/mvp/workflow/rag-corpus-and-grounding.md](docs/mvp/workflow/rag-corpus-and-grounding.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md), [docs/mvp/quality/eval-plan.md](docs/mvp/quality/eval-plan.md), [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker, eval planning, future-agent prompts, and implementation behavior synchronized after adding the retrieval slice.
- How: Created the dedicated Phase 3 artifact, checked off the Phase 3 tracker, added retrieval eval mappings, and updated the reusable RAG planning prompt to point at the new source document.

### Task: Update Repository State Documentation

- What: Updated repository overview language to include the Phase 3 retrieval package and its targeted test command without claiming a complete end-to-end incident workflow.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Avoid stale statements that said retrieval did not exist while preserving limits around timeline, severity, brief, approval, export, escalation, and external sharing.
- How: Added the Phase 3 artifact to documentation maps, updated runtime-surface notes, marked the retrieved-content trust boundary as implemented, and kept later workflow phases unchecked.

## 2026-05-06 - Phase 2 Incident Packet Ingestion

### Task: Add Phase 2 Go Module And TDD Coverage

- What: Added the initial Go module and failing-first ingestion tests for valid packets, missing required fields, non-synthetic records, malformed incident IDs, malformed telemetry, unsupported event types, non-synthetic evidence references, and audit event emission.
- Where: [go.mod](go.mod), [internal/ingestion/ingestion_test.go](internal/ingestion/ingestion_test.go).
- When: 2026-05-06, America/Vancouver.
- Why: Start Phase 2 with strict TDD and make packet validation observable before adding production ingestion logic.
- How: Created tests using synthetic inline packet JSON from the Phase 1 contract, observed red failures for missing APIs and later missing validation rules, then used those tests to drive the implementation.

### Task: Implement Synthetic Packet Validation

- What: Added JSON ingestion, packet schema types, deterministic validation errors, supported event-type validation, synthetic-only evidence validation, telemetry validation, and accepted/rejected audit events.
- Where: [internal/ingestion/ingestion.go](internal/ingestion/ingestion.go).
- When: 2026-05-06, America/Vancouver.
- Why: Ensure no downstream reasoning step can use unvalidated or non-synthetic incident packet data.
- How: Implemented `IngestJSON` with RFC3339 timestamp parsing, required-field checks, `FIC-SYN-` incident ID enforcement, speed bounds, relative-time parsing, `synthetic://` media reference enforcement, and stable validation issue codes/messages.

### Task: Document Phase 2 Ingestion Contract

- What: Added the Phase 2 ingestion schema, validation rules, audit event contract, test commands, red-to-green evidence, and current runtime limits.
- Where: [docs/mvp/workflow/incident-packet-ingestion.md](docs/mvp/workflow/incident-packet-ingestion.md), [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep implementation behavior, acceptance criteria, and phase status synchronized after adding runtime code.
- How: Created a dedicated Phase 2 artifact and checked off the Phase 2 tracker with links to the document and Go package.

### Task: Update Repository State Documentation

- What: Replaced documentation-only current-state language with an accurate Phase 2 runtime summary and linked the ingestion artifact.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md), [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Avoid overclaiming a complete app while no longer claiming the repo has no application code or tests.
- How: Added runtime-surface notes, marked only the ingestion validation promise complete, and recorded the implemented untrusted-packet boundary while leaving later phases unchecked.

## 2026-05-06 - Phase 1 Synthetic Evidence And Workflow Contract

### Task: Create Synthetic Incident Packet Contract

- What: Added the Phase 1 synthetic incident packet specs and shared workflow-output contract.
- Where: [docs/mvp/workflow/synthetic-incident-packets.md](docs/mvp/workflow/synthetic-incident-packets.md).
- When: 2026-05-06, America/Vancouver.
- Why: Define realistic fake evidence and expected workflow outputs before any ingestion, reasoning, eval, or fixture implementation begins.
- How: Created five explicitly synthetic records covering low, medium, high, unknown, and adversarial or missing-data scenarios with required fields, telemetry samples, media references, transcript notes, still-frame notes, expected timelines, severity, recommendations, brief requirements, missing-data behavior, and acceptance criteria.

### Task: Mark Phase 1 Complete

- What: Checked off the Phase 1 planning checklist and linked the synthetic packet output artifact.
- Where: [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep the phase tracker aligned with completed Markdown planning work while avoiding any claim that runtime behavior exists.
- How: Marked only the Phase 1 documentation items complete and added a dated note that Phase 1 added Markdown planning artifacts only.

### Task: Update Documentation Indexes

- What: Linked the Phase 1 packet contract from the root and MVP README files, and refreshed the current-state language.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md).
- When: 2026-05-06, America/Vancouver.
- Why: Make the new Phase 1 artifact discoverable and keep the repo narrative clear that no application code, tests, scaffolding, dependencies, or runtime configuration were added.
- How: Added the synthetic packet artifact to the planning artifact maps and changed the current-state wording from Phase 0-only planning to Phase 0 plus Phase 1 planning.

### Task: Sync Future Prompt And Eval References

- What: Pointed future documentation agents and eval planning at the Phase 1 packet contract.
- Where: [docs/mvp/execution/task-prompts.md](docs/mvp/execution/task-prompts.md), [docs/mvp/quality/eval-plan.md](docs/mvp/quality/eval-plan.md).
- When: 2026-05-06, America/Vancouver.
- Why: Preserve the Phase 1 source of truth for future fixture, eval, and implementation work without creating machine-readable fixtures prematurely.
- How: Updated the synthetic incident planning prompt to name the Markdown contract file and added a Phase 1 packet mapping for future eval fixture work.

## 2026-05-06 - Phase 0 Product Frame And Guardrails

### Task: Create Product Frame

- What: Added the Phase 0 product frame for Fleet Incident Copilot.
- Where: [docs/mvp/overview/product-frame.md](docs/mvp/overview/product-frame.md).
- When: 2026-05-06, America/Vancouver.
- Why: Establish the primary user, problem, MVP promise, scope boundaries, trust boundaries, approval gates, prohibited claims, demo narrative, and success criteria before implementation.
- How: Grounded the frame in [research-report.md](docs/research/research-report.md), kept all data synthetic, and explicitly confirmed no code is needed for Phase 0.

### Task: Mark Phase 0 Complete

- What: Checked off the Phase 0 planning checklist and linked the output artifact.
- Where: [docs/mvp/execution/phases.md](docs/mvp/execution/phases.md).
- When: 2026-05-06, America/Vancouver.
- Why: Make the phase tracker reflect completed planning work without implying runtime behavior exists.
- How: Marked only the Phase 0 documentation items complete and added a dated note that Phase 0 added Markdown planning artifacts only.

### Task: Update Documentation Indexes

- What: Linked the Product Frame from the root and MVP README files, and recorded the Go backend direction.
- Where: [README.md](README.md), [docs/mvp/README.md](docs/mvp/README.md).
- When: 2026-05-06, America/Vancouver.
- Why: Keep navigation current and capture that future backend code phases should use Go.
- How: Added the Product Frame artifact to the planning artifact maps and clarified that current behavior is planning-only.

### Task: Refine Scope Guardrails

- What: Added explicit prohibited claims and approval-gate rules.
- Where: [docs/mvp/overview/scope.md](docs/mvp/overview/scope.md).
- When: 2026-05-06, America/Vancouver.
- Why: Make the review boundaries visible before any implementation agent starts future work.
- How: Added scope notes that distinguish planning checklists from implemented behavior, recorded Go as the future backend language, and made export, escalation, and external sharing approval requirements explicit.
