# How-Tos

Common tasks for working on the Fleet Incident Copilot MVP.

## How To Find The Current Scope

1. Start with [README.md](../README.md) for the current runtime surface.
2. Read [MVP Overview](mvp/README.md) for the demo promise and artifact map.
3. Check [Scope And Guardrails](mvp/overview/scope.md) before adding behavior.
4. Check [Phases And Tasks](mvp/execution/phases.md) when a task belongs to a planned phase.
5. Use [Strict TDD Rules](mvp/execution/tdd-rules.md) for every code change.

If a behavior is not implemented in `internal`, describe it as planned or deferred. Do not imply a CLI, HTTP API, database, live model call, real export, real escalation, external sharing, identity, dashboards, or production compliance behavior exists.

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

The package-level demo path is:

```text
internal/ingestion -> internal/retrieval -> internal/timeline -> internal/severity -> internal/brief
```

Approval, eval, and observability sit alongside that path:

```text
internal/approval      gates sensitive actions
internal/eval          scores deterministic golden cases
internal/observability records in-memory workflow events and budget signals
```

Use tests for examples of how packages are composed. The eval package has the broadest in-memory composition surface through `GoldenCases` and `Run`.

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
