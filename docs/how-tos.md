# How-Tos

Common tasks for working on the Fleet Incident Copilot MVP.

## How To Find The Current Scope

1. Start with [README.md](../README.md) for the current runtime surface.
2. Read [MVP Overview](mvp/README.md) for the demo promise and artifact map.
3. Check [Scope And Guardrails](mvp/overview/scope.md) before adding behavior.
4. Check [Phases And Tasks](mvp/execution/phases.md) when a task belongs to a planned phase.
5. Use [Strict TDD Rules](mvp/execution/tdd-rules.md) for every code change.

If a behavior is not implemented in `internal` or the thin `cmd/demo-api` loopback server, describe it as planned or deferred. Do not imply a general CLI workflow, production HTTP API, database, live model call, real export, real escalation, real external sharing, Slack delivery, identity, dashboards, or production compliance behavior exists.

## How To Add Or Change Package Behavior

1. Pick the smallest observable behavior.
2. Add or update the focused test in the matching `internal/<package>` test file.
3. Run the targeted package test and confirm it fails for the expected reason.
4. Implement the smallest production change.
5. Re-run the targeted package test.
6. Run `go test ./...`.
7. Update docs when commands, behavior, scope, or acceptance criteria changed.
8. Add a changelog entry for behavior changes or meaningful documentation updates.

Example targeted commands:

```bash
go test ./internal/ingestion
go test ./internal/retrieval
go test ./internal/timeline
go test ./internal/severity
go test ./internal/brief
go test ./internal/approval
go test ./internal/eval
go test ./internal/observability
go test ./internal/demo
go test ./internal/notification
go test ./internal/httpapi
go test ./cmd/demo-api
```

## How To Add A Synthetic Incident Fixture

1. Keep the incident ID prefixed with `FIC-SYN-`.
2. Set `SyntheticRecord` to `true`.
3. Use synthetic vehicle IDs, routes, location labels, media references, transcript notes, and still-frame notes.
4. Keep media references in the synthetic URI shape, such as `synthetic://fic-syn-001/front-camera-074218.jpg`.
5. Include at least one telemetry sample with `RelativeTime`, `SpeedMPH`, `Heading`, `Signal`, and `GPSLabel`.
6. Add normal, incomplete, or adversarial expectations when the fixture participates in eval coverage.

The ingestion package rejects non-synthetic records, unsupported event types, malformed timestamps, impossible speeds, and non-synthetic media references.

## How To Trace The Demo Workflow In Code

The package-level workflow path is:

```text
internal/ingestion -> internal/retrieval -> internal/timeline -> internal/severity -> internal/brief
```

Approval, eval, and observability sit alongside that path:

```text
internal/approval      gates sensitive actions
internal/eval          scores deterministic golden cases
internal/observability records in-memory workflow events and budget signals
```

Use tests for examples of how packages are composed. The demo package owns the Phase 12 review composition contract through `ComposeIncident` and `ComposePacket`; `internal/httpapi` owns the Phase 13 loopback review route, Phase 14 dry-run notification route, and Phase 15 approval retry routes; `internal/notification` owns the dry-run Slack-shaped preview; the eval package owns golden-case scoring through `GoldenCases` and `Run`.

## How To Run The Loopback Demo API

Start the default loopback server when port `8080` is free:

```bash
go run ./cmd/demo-api
```

Use a loopback override when port `8080` is occupied:

```bash
go run ./cmd/demo-api -addr 127.0.0.1:18080
```

Then call the review endpoint from another terminal:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18080/demo/review \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001"}'
```

Call the dry-run notification preview endpoint from another terminal:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18080/demo/notifications/slack \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","channel":"#fleet-safety","delivery_mode":"dry_run"}'
```

The preview response should be blocked before scoped approval, include a prepared Slack-shaped payload, and report `sent: false` plus `network_delivery_attempted: false`.

Create an in-memory scoped approval request:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18080/demo/approvals \
  -H "Content-Type: application/json" \
  -d '{"incident_id":"FIC-SYN-001","action":"external_sharing","channel":"#fleet-safety","reason":"operator requested dry-run preview"}'
```

On a fresh server process, approve the deterministic request ID:

```bash
curl -i --max-time 5 -X POST http://127.0.0.1:18080/demo/approvals/decisions \
  -H "Content-Type: application/json" \
  -d '{"request_id":"approval-001","approver":"fleet-safety-lead","decision":"approved","reason":"redacted brief approved for #fleet-safety dry-run"}'
```

Retry the dry-run notification preview. It should return `notification_preview.status: "allowed"` only for the same incident, `external_sharing` action, and `#fleet-safety` channel, while still reporting `sent: false` and `network_delivery_attempted: false`.

## How To Update Documentation Safely

1. Keep the implemented-versus-planned boundary explicit.
2. Link to source docs instead of repeating long behavior descriptions.
3. Update [docs/README.md](README.md) when adding top-level docs.
4. Update [README.md](../README.md) when the repo entry points change.
5. Run a Markdown whitespace check:

```bash
git diff --check -- README.md CHANGELOG.md docs
```

## How To Prepare A Local Review

Run:

```bash
go test ./...
go vet ./...
go test -cover ./...
git diff --check -- README.md CHANGELOG.md docs internal
git status --short
```

Then summarize the changed behavior, docs touched, tests run, and any residual risk.
