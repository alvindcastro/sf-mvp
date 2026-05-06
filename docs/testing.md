# Testing

The repository uses Go tests as the main verification surface. There is no Makefile or separate test runner.

## Quick Commands

Run the full package suite:

```bash
go test ./...
```

Run vet:

```bash
go vet ./...
```

Run coverage for all packages:

```bash
go test -cover ./...
```

Check Markdown and Go diff whitespace:

```bash
git diff --check -- README.md CHANGELOG.md docs internal
```

## Targeted Package Tests

Use targeted commands while developing:

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
go test ./internal/httpapi
go test ./cmd/demo-api
```

Run one test by name:

```bash
go test ./internal/eval -run TestRunHandlesPromptInjectionFixtures
```

Run repeated checks when diagnosing nondeterminism:

```bash
go test ./... -count=10
```

## TDD Expectations

Every code task should follow [Strict TDD Rules](mvp/execution/tdd-rules.md):

1. Add or update a focused failing test.
2. Confirm the failure is for the expected missing behavior.
3. Add the smallest production change.
4. Re-run the targeted package test.
5. Refactor only after green.
6. Run the relevant broader suite.
7. Update docs when behavior, commands, architecture, or acceptance criteria changed.

## What Each Package Test Covers

- `internal/ingestion`: packet validation, synthetic-only rules, malformed input, telemetry validation, and audit events.
- `internal/retrieval`: scope filtering, deterministic ranking, citation metadata, no-match behavior, and hostile text as data.
- `internal/timeline`: telemetry ordering, source refs, unavailable evidence, conflicting signals, and guidance source propagation.
- `internal/severity`: low, medium, high, unknown, conflicting evidence, recommendations, citations, and approval-required flags.
- `internal/brief`: cited draft sections, redaction, fail-closed missing evidence, uncertainty labels, and approval-state display.
- `internal/approval`: pending, approved, denied, missing approval, out-of-scope calls, callback execution, and immutable audit history.
- `internal/eval`: golden cases, release thresholds, unsupported claims, redaction, prompt injection, and approval fail-closed behavior.
- `internal/observability`: trace IDs, structured events, redaction, token accounting, budget limits, invalid token usage, cache candidates, and routing notes.
- `internal/demo`: machine-readable synthetic fixture loading, deterministic review composition, non-synthetic rejection, missing-evidence fail-closed behavior, citation and redaction preservation, approval-required action display, and trace propagation.
- `internal/httpapi`: loopback review handler behavior, deterministic JSON response shape, malformed JSON, unknown ID, non-synthetic rejection, unsupported methods, unknown paths, and loopback defaults.
- `cmd/demo-api`: thin local server wiring, default loopback address, loopback override, and non-loopback rejection.

## Coverage

Coverage is useful as a regression signal, but the repo does not currently enforce coverage through a Makefile or CI config. Use:

```bash
go test -cover ./...
```

For a detailed local profile:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Remove generated local coverage files before committing unless the user explicitly wants them.

## Documentation-Only Changes

For doc-only edits, run at least:

```bash
git diff --check -- README.md CHANGELOG.md docs
go test ./...
```

The Go test suite is still useful because documentation in this repo often describes exact package behavior.
