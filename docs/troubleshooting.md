# Troubleshooting

Common problems and the fastest checks.

## `make` Commands Fail

There is no Makefile in this repo. Use Go commands directly:

```bash
go test ./...
go vet ./...
go test -cover ./...
```

## No App Starts

The repo does not currently include a server, CLI, frontend, worker, container setup, or database runtime. The runnable surface is the Go package test suite.

## Import Errors

Check that imports use the module path from [go.mod](../go.mod):

```text
sf-mvp/internal/<package>
```

Then run:

```bash
go test ./...
```

## Ingestion Rejects A Packet

Check the returned `ingestion.ValidationError` codes. Frequent causes:

- `synthetic_record` is missing or false.
- `incident_id` does not start with `FIC-SYN-`.
- `timestamp` is not RFC3339.
- `event_type` is unsupported.
- `telemetry_samples` is empty or has malformed fields.
- `media_references` is empty or uses a non-synthetic URI.

Use `internal/ingestion/ingestion_test.go` for valid packet examples.

## Retrieval Returns No Matches

Check:

- `Query.Text` has meaningful non-stopword terms.
- `Query.Workflow` is not empty.
- `Query.Scope` is not empty.
- Documents match the same workflow and scope exactly.
- The document title or body shares terms with the query.

No matches are expected when the approved mock corpus does not cover the question.

## Timeline Entries Are Marked Uncertain

Uncertainty is expected when:

- Telemetry relative time cannot be parsed.
- Media, transcript, or still-frame text says evidence is unavailable.
- Multiple telemetry entries have the same timestamp but different claims.

Do not suppress uncertainty unless the underlying packet or timeline rule changed intentionally.

## Severity Is `unknown`

Expected causes:

- The packet event type is `unknown_trigger`.
- Timeline evidence has conflicting telemetry.
- No deterministic rule covers the supplied event type.

Unknown severity should route toward human operator review instead of model-only judgment.

## Brief Drafting Fails Closed

`brief.Draft` returns `brief.MissingEvidenceError` when required evidence is missing. Check for:

- Empty incident ID.
- No timeline entries.
- Timeline entries without claims or citations.
- Missing severity level.
- Missing severity rationale.
- Recommendations without actions, explanations, or citations.
- Missing approval requirements.

Fail-closed behavior is intentional for shareable outputs.

## Demo Review Composition Fails

`demo.ComposeIncident` returns `demo.ErrIncidentNotFound` when the incident ID is not present in the loaded synthetic fixtures.

`demo.ComposePacket` returns `demo.ErrNonSyntheticInput` before downstream composition when the packet is not synthetic or the incident ID does not start with `FIC-SYN-`.

`demo.ErrMissingEvidence` means the composer reached the existing fail-closed brief contract and did not return a partial review result. Check the same missing-evidence causes listed for brief drafting.

## Sensitive Action Is Blocked

`approval.Gate.Execute` blocks by default. Confirm:

- The request action matches the call action.
- The request scope exactly matches the call scope.
- `IncidentID` matches the scope incident ID.
- A human decision approved the request.
- The request was not denied or still pending.

Denied or pending approval should remain blocked.

## Eval Report Fails

Inspect the failing `eval.CaseResult` fields:

- `ActualSeverity`
- `MissingRecommendations`
- `MissingGuidanceRefs`
- `UnsupportedClaims`
- `RedactionLeaks`
- `PromptInjectionResistant`
- `SensitiveActionsBlockedWithoutApproval`
- `Failures`

Default thresholds are strict. A small behavior change can fail the whole report when it affects citations, recommendations, redaction, or approval fail-closed behavior.

## Observability Returns Budget Or Token Errors

`observability.RecordModelCall` returns `ErrInvalidTokenUsage` when input or output token counts are negative.

It returns `ErrBudgetExceeded` when the call would exceed the configured input, output, total token, or model-call budget. These are local budget checks; no provider billing lookup exists.

## `git diff --check` Reports Existing Files

The working tree may contain unrelated local changes. If whitespace failures are outside your edit set, do not rewrite those files as part of an unrelated task. Check only the files you changed or ask before touching unrelated dirty files.
