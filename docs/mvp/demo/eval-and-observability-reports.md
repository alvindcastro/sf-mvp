# Eval And Observability Demo Reports

Phase 16 exposes local quality and operations proof through the loopback demo API. The routes are deterministic, in memory, synthetic-only, and backed by the existing `internal/eval` and `internal/observability` packages plus thin `internal/httpapi` response projection.

This is a demo report surface, not a monitoring product.

## Phase 16 Checklist

- [x] Add a local eval report surface that runs deterministic golden cases and returns case count, metric scores, thresholds, gates, and pass/fail status.
- [x] Add an in-memory trace report surface that returns redacted events by trace ID.
- [x] Include a budget-exceeded demo path using caller-supplied token counts.
- [x] Include eval summary and notification preview tool-call events when available.
- [x] Keep reports local and ephemeral.
- [x] Do not imply dashboards, alerts, OpenTelemetry export, persistent logs, provider billing reconciliation, or model benchmarking.

## Runtime Surface

- Handler package: [internal/httpapi](../../../internal/httpapi).
- Eval package: [internal/eval](../../../internal/eval).
- Observability package: [internal/observability](../../../internal/observability).
- Local server command: [cmd/demo-api](../../../cmd/demo-api).
- Eval report route: `GET /demo/eval/latest`.
- Trace report route: `GET /demo/traces/{trace_id}`.
- Budget demo route: `POST /demo/budget/check`.
- Targeted handler test: `go test ./internal/httpapi`.
- Related package tests: `go test ./internal/eval ./internal/observability ./internal/httpapi`.
- Full test command: `go test ./...`.

Trace reports are stored only in the current `internal/httpapi.Handler` instance. Restarting the demo server, or creating a fresh handler in tests, loses prior trace reports.

## Eval Report

`GET /demo/eval/latest` runs `eval.Run(eval.GoldenCases(), eval.DefaultThresholds())` and returns:

- `trace_id` for the report run.
- `eval_report.case_count`.
- `eval_report.passed`.
- `eval_report.metrics.severity_accuracy`.
- `eval_report.metrics.citation_coverage`.
- `eval_report.metrics.recommendation_accuracy`.
- `eval_report.thresholds`.
- `eval_report.gates` with pass/fail status for metric and safety gates.

Expected response highlights:

```json
{
  "trace_id": "trace-fic-syn-eval-report-20260506t160000z-001",
  "eval_report": {
    "case_count": 5,
    "passed": true,
    "metrics": {
      "severity_accuracy": 1,
      "citation_coverage": 1,
      "recommendation_accuracy": 1
    }
  }
}
```

The route records `workflow.started` and `eval.score_recorded` events into the in-memory trace report store.

## Trace Report

`GET /demo/traces/{trace_id}` returns only events already recorded by the current handler instance for that trace ID.

Expected response highlights:

```json
{
  "trace_report": {
    "trace_id": "trace-fic-syn-eval-report-20260506t160000z-001",
    "incident_id": "FIC-SYN-EVAL-REPORT",
    "ephemeral": true,
    "events": [
      {
        "type": "workflow.started"
      },
      {
        "type": "eval.score_recorded",
        "fields": {
          "passed": "true"
        },
        "metrics": {
          "case_count": 5
        }
      }
    ]
  }
}
```

Missing trace IDs return `404` with `error.code: "trace_not_found"`.

The route can show notification preview tool-call events after `POST /demo/notifications/slack` runs in the same process. The event fields include `tool_name: "notification.slack.preview"`, `delivery_mode: "dry_run"`, `sent: "false"`, and `network_delivery_attempted: "false"`.

## Budget Demo

`POST /demo/budget/check` uses caller-supplied token counts against caller-supplied local limits. It does not call a model provider and does not reconcile provider billing.

Example request:

```json
{
  "incident_id": "FIC-SYN-001",
  "provider": "hosted",
  "model": "demo-review-model",
  "input_tokens": 90,
  "output_tokens": 20,
  "max_total_tokens": 100
}
```

Expected response highlights:

```json
{
  "budget_report": {
    "status": "budget_exceeded",
    "reason": "total token budget exceeded",
    "token_usage": {
      "input_tokens": 90,
      "output_tokens": 20,
      "total_tokens": 110
    }
  }
}
```

The route records `budget.exceeded` in the trace report when the supplied token usage exceeds the supplied budget. `sensitive_terms` may be supplied to redact matching event fields in the returned trace report.

## Current Limits

- Reports are local to the loopback demo server process.
- Trace history is in memory and ephemeral.
- Eval cases are deterministic synthetic Go fixtures.
- Budget checks use caller-supplied token counts.
- No model call, model benchmark, provider billing reconciliation, persistent log store, database, dashboard, alert, OpenTelemetry export, external observability pipeline, production audit store, identity, auth, Slack delivery, webhook, export, escalation, or external-sharing integration exists.

## Red-To-Green Evidence

- Red: added failing `internal/httpapi` tests for `GET /demo/eval/latest`, `GET /demo/traces/{trace_id}`, `POST /demo/budget/check`, missing trace handling, redacted budget event fields, notification preview tool-call trace events, and per-handler ephemeral storage; observed `go test ./internal/httpapi` fail with `404` for missing Phase 16 routes.
- Green: implemented the smallest handler wiring over `internal/eval` and `internal/observability`, added the in-memory trace store, projected eval metrics, thresholds, gates, trace events, and token usage, then verified `go test ./internal/httpapi` passed.
