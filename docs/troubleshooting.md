# Troubleshooting

Common problems and the fastest checks.

## `make` Commands Fail

The Makefile only wraps local verification commands. If a Make target fails,
run the underlying Go command directly:

```bash
go test ./...
go vet ./...
go test -cover ./...
go test ./internal/eval ./cmd/evalops-target ./cmd/evalops-gate -count=1
go run ./cmd/evalops-gate
```

## Demo API Does Not Start

The repo includes only a loopback demo server, not a production app stack. Start it with:

```bash
go run ./cmd/demo-api
```

If `127.0.0.1:8080` is already occupied, use another loopback port:

```bash
go run ./cmd/demo-api -addr 127.0.0.1:18080
```

Non-loopback addresses such as `0.0.0.0:8080` are rejected by design. There is no frontend, worker, container setup, database runtime, or production service entry point.

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

## Demo API Request Fails

`POST /demo/review` accepts either `{"incident_id":"FIC-SYN-001"}` or a full synthetic packet JSON body.

`POST /demo/notifications/slack` accepts a known synthetic incident ID, channel, and dry-run mode:

```json
{"incident_id":"FIC-SYN-001","channel":"#fleet-safety","delivery_mode":"dry_run"}
```

`POST /demo/approvals` accepts a known synthetic incident ID, sensitive action, target channel, and human reason:

```json
{"incident_id":"FIC-SYN-001","action":"external_sharing","channel":"#fleet-safety","reason":"operator requested dry-run preview"}
```

`POST /demo/approvals/decisions` records the human decision for an existing request:

```json
{"request_id":"approval-001","approver":"fleet-safety-lead","decision":"approved","reason":"redacted brief approved for #fleet-safety dry-run"}
```

`GET /demo/eval/latest` runs the deterministic local golden-case eval report.

`GET /demo/traces/{trace_id}` returns redacted in-memory events for a trace recorded by the current server process.

`POST /demo/budget/check` accepts caller-supplied token counts and local budget limits:

```json
{"incident_id":"FIC-SYN-001","provider":"hosted","model":"demo-review-model","input_tokens":90,"output_tokens":20,"max_total_tokens":100}
```

Common response codes:

- `400 malformed_json`: the request body is not valid JSON.
- `400 dry_run_required`: notification preview requests must use `delivery_mode: "dry_run"`.
- `400 invalid_approval_request`: the approval request or decision payload is incomplete or unsupported.
- `400 invalid_budget_report`: the budget report request is incomplete or cannot start a local trace.
- `400 invalid_token_usage`: budget report token counts cannot be negative.
- `404 approval_request_not_found`: the decision route used an unknown approval request ID.
- `404 incident_not_found`: the synthetic incident ID is not in the demo fixture set.
- `404 trace_not_found`: the trace ID is missing or was not recorded by the current local server process.
- `405 method_not_allowed`: use `POST`.
- `409 approval_already_decided`: final approval decisions cannot be rewritten.
- `422 non_synthetic_input`: `synthetic_record` is false or the incident ID does not start with `FIC-SYN-`.
- `422 missing_evidence`: the underlying composer failed closed rather than returning a partial review.

## Dry-Run Notification Preview Is Blocked

Blocked notification preview is expected before exact scoped approval. Phase 15 keeps approval state in memory for one `cmd/demo-api` process, so restarting the server clears approval requests.

Check:

- `delivery_mode` is exactly `dry_run`.
- The incident ID exists in the synthetic fixture set.
- The approval request action is exactly `external_sharing`.
- The approval request target ref is exactly `slack:<channel>`.
- The decision is `approved`, not `pending` or `denied`.
- The request was created in the same local server process as the retry.
- The response still has `prepared_payload`.
- `sent` is `false`.
- `network_delivery_attempted` is `false`.

Do not add Slack tokens, webhook URLs, SDKs, or network calls to fix a blocked preview. Use `POST /demo/approvals` and `POST /demo/approvals/decisions` for the local approval retry demo.

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

For the loopback route, call `GET /demo/eval/latest` and then use the returned `trace_id` with `GET /demo/traces/{trace_id}` to inspect the local `eval.score_recorded` event.

## Observability Returns Budget Or Token Errors

`observability.RecordModelCall` returns `ErrInvalidTokenUsage` when input or output token counts are negative.

It returns `ErrBudgetExceeded` when the call would exceed the configured input, output, total token, or model-call budget. These are local budget checks; no provider billing lookup exists.

For the loopback route, `POST /demo/budget/check` stores the local `model_call.recorded`, `model_call.rejected`, or `budget.exceeded` event in the current server process only. Restarting the server clears the trace report.

## `git diff --check` Reports Existing Files

The working tree may contain unrelated local changes. If whitespace failures are outside your edit set, do not rewrite those files as part of an unrelated task. Check only the files you changed or ask before touching unrelated dirty files.
