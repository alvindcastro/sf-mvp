# Nice To Knows

Useful context that prevents common wrong assumptions.

## Runtime Surface

- There is no Makefile in this repo.
- There is a thin loopback-only demo API command at `cmd/demo-api`.
- There is no general CLI workflow, database, frontend, worker, queue, or container runtime.
- Most implementation remains Go package-level behavior under `internal`.
- Most state is in memory and deterministic by design.
- The module name is `sf-mvp`, even though the product narrative is Fleet Incident Copilot.
- `internal/demo` composes package results in memory; `internal/notification` prepares dry-run Slack-shaped previews; `internal/httpapi` wraps those packages, the in-memory approval retry flow, local eval reports, trace reports, and caller-supplied budget reports for local demo routes.

## Demo Boundaries

- The demo is synthetic-only.
- A checked item in a planning doc means the repo has some supporting artifact or package behavior, not necessarily production readiness.
- The demo package is for explaining architecture and judgment; it is not a deployable app.
- Machine-readable demo fixtures live under `internal/demo/testdata` and are loaded through ingestion validation.
- `cmd/demo-api` is a local walkthrough server, not a production service.
- Approval requirements in severity output are not the same as executing a real export, escalation, or external share.
- The approval package gates callbacks, but the callbacks are local function calls, not real integrations.
- The dry-run Slack-shaped preview is not Slack delivery. It has no Slack SDK, token, webhook URL, environment secret, outbound network request, or real external-sharing behavior.
- The notification preview route uses in-memory approval state in Phase 15, so it returns `blocked` before approval and `allowed` only after an exact approved `external_sharing` request for the same incident and Slack-shaped target channel.
- The report routes added in Phase 16 are local and ephemeral: `GET /demo/eval/latest`, `GET /demo/traces/{trace_id}`, and `POST /demo/budget/check`.

## Retrieval And Citations

- Retrieval filters by exact `Workflow` and `Scope` before ranking.
- Retrieved snippets are marked as `retrieved_data`.
- Citation refs are deterministic and include the source ID plus revision date.
- No-match retrieval returns an empty result instead of invented guidance.
- Hostile retrieved text is preserved as fixture data, not followed as instruction.

## Evaluation

- Golden cases live in code in `internal/eval`, not in external fixture files.
- Default eval thresholds require perfect severity accuracy, citation coverage, and recommendation accuracy.
- Eval checks redaction, unsupported claims, prompt-injection resistance, and approval fail-closed behavior.
- Eval is deterministic and local. It does not call a model provider.
- The eval report route runs the same deterministic golden cases and returns scores, thresholds, gates, and pass/fail state.

## Observability

- Observability records in-memory events.
- Token usage is caller-supplied; no provider billing reconciliation exists.
- Budget failures and invalid token usage produce local events and errors.
- Sensitive terms and coordinate-like values are redacted from event fields.
- Cache candidates and model-routing notes are planning signals, not active cache storage or live routing.
- Trace reports only show events already recorded by the current loopback handler process.

## Documentation

- The docs intentionally repeat the current limitations in several places to prevent overclaiming.
- Prefer adding a precise link over rewriting the same behavior in multiple docs.
- If a doc says future or planned, keep it that way until code and tests exist.
- If package behavior changes, update both the workflow doc and any top-level summary that mentions it.
