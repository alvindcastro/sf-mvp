# EvalOps Trace And Score Export

FQ13 adds safe trace attributes and score events that can be adapted to
OpenTelemetry, Jaeger, Langfuse, or another EvalOps backend. The repository
still runs locally by default: no external collector, Langfuse server, model
provider, cloud service, or network exporter is started by these changes.

## Code Surface

- `observability.WorkflowAttributes(...)` maps workflow facts into safe
  OpenTelemetry-style attributes.
- `eval.ScoreEventsFromPromptfooOutput(...)` maps Promptfoo-style scorer output
  into stable score events.
- `eval.ExportPromptfooScoreEvents(...)` exports score events through an
  injectable interface with disabled, best-effort, and release-gate modes.
- `eval.NewIncidentEvalTarget(...)` can receive `WithScoreEventExporter(...)`
  and `WithScoreEventExportMode(...)` options.

## Workflow Attributes

`WorkflowAttributes` returns string attributes designed for span or event
metadata:

```json
{
  "workflow.trace_id": "trace-fic-syn-017-20260506t160000z-001",
  "workflow.incident_id_hash": "sha256:cdcd7e14394344372016403370275f9d314de594321ed9cca7b7e977f058ecc1",
  "workflow.retrieved_source_ids": "FIC-SOP-HARD-BRAKE-001,FIC-TS-MISSING-MEDIA-001",
  "workflow.severity": "medium",
  "workflow.approval_state": "pending",
  "workflow.latency_ms": "42.5"
}
```

The mapper intentionally hashes the incident ID and omits transcript notes,
media references, still-frame notes, vehicle IDs, route labels, locations, and
other raw evidence. Retrieved source IDs are normalized, deduplicated, and
sorted for stable traces.

## Score Events

Score events use the stable name pattern `eval.score.<scorer>`.

```json
{
  "name": "eval.score.unsupported_claims",
  "run_id": "eval-run-13",
  "trace_id": "trace-abc",
  "case_id": "adversarial failure",
  "incident_id": "FIC-SYN-005",
  "kind": "adversarial",
  "scorer": "unsupported_claims",
  "score": 0,
  "pass": false,
  "critical": true,
  "severity": "critical",
  "reason": "unsupported claims=1"
}
```

Supported export modes:

- `disabled`: do not call the exporter.
- `best_effort`: call the exporter, but preserve the core eval result if export
  fails.
- `release_gate`: fail the target response if score event export fails.

The default target mode is `best_effort` with a no-op exporter, so local runs
remain deterministic and no-key.

Promptfoo or another caller can pass correlation data through request `vars`:

```json
{
  "incident_id": "FIC-SYN-001",
  "vars": {
    "trace_id": "trace-local-eval-001",
    "eval_run_id": "eval-run-local-001"
  }
}
```

`eval_run_id` is preferred. `run_id` is also accepted as a fallback.

## Local OTel Or Langfuse Setup

The current code exposes an exporter interface but does not include an OTel or
Langfuse SDK adapter. A local adapter can read these conventional variables
when one is added:

```bash
EVALOPS_SCORE_EVENT_EXPORT_MODE=disabled
OTEL_SERVICE_NAME=sf-mvp-evalops
OTEL_EXPORTER_OTLP_ENDPOINT=http://127.0.0.1:4318
LANGFUSE_HOST=http://127.0.0.1:3000
LANGFUSE_PUBLIC_KEY=local-public-key
LANGFUSE_SECRET_KEY=local-secret-key
```

For local collector testing, run an OpenTelemetry collector or Jaeger all-in-one
process outside this repository, point `OTEL_EXPORTER_OTLP_ENDPOINT` at the
collector, and keep `EVALOPS_SCORE_EVENT_EXPORT_MODE=best_effort` until release
gates explicitly need exporter failures to block.

For local Langfuse testing, run Langfuse outside this repository, point
`LANGFUSE_HOST` at that service, and keep the current repository default
disabled/no-op behavior unless an injected exporter adapter is present.

## Safety Boundary

FQ13 does not change the synthetic-only data boundary. Tests assert:

- workflow attributes carry trace ID, incident hash, source IDs, severity,
  approval state, and latency without raw evidence;
- score events carry pass/fail, score, criticality, run ID, and trace ID;
- disabled export mode skips the exporter;
- best-effort exporter failures do not break incident eval results;
- release-gate exporter failures are machine-readable failures.

## Verification

FQ13 is covered by:

```bash
go test ./internal/observability -count=1
go test ./internal/eval -count=1
go test ./internal/eval ./internal/observability ./cmd/evalops-target -count=1
go test ./...
```
