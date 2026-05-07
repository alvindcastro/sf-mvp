# EvalOps Promptfoo Bridge

FQ12 adds a local Promptfoo-compatible bridge over the deterministic Go eval
path. It does not call model providers, live network services, Slack, cloud
search, or external telemetry systems.

## Code Surface

- `eval.NewIncidentEvalTarget()` returns an HTTP handler for
  `POST /evalops/incident`.
- `eval.NewLocalIncidentEvalWorkflow(...)` composes the local ingestion,
  retrieval, timeline, severity, brief, approval, and eval scoring steps behind
  injectable interfaces.
- `eval.PromptfooOutputFromResult(...)` converts `eval.CaseResult` into
  Promptfoo/EvalOps-style scorer output.
- `cmd/evalops-target` starts the loopback-only local target process.
- `evals/promptfoo/fleet-incident.yaml` calls the local Go target with normal,
  incomplete, and adversarial golden cases.

## Local Run

Start the target in one terminal:

```bash
go run ./cmd/evalops-target -addr 127.0.0.1:18085
```

Run Promptfoo in another terminal:

```bash
npx promptfoo eval -c evals/promptfoo/fleet-incident.yaml
```

The config uses Promptfoo's HTTP provider against
`http://127.0.0.1:18085/evalops/incident`. No OpenAI, Anthropic, Google, or
other model-provider key is required.

## Request Shape

Golden case lookup:

```json
{
  "case_id": "adversarial transcript with missing side view",
  "incident_id": "FIC-SYN-005",
  "timeout_ms": 5000
}
```

Packet-based eval:

```json
{
  "case_id": "packet hard brake",
  "packet": {
    "synthetic_record": true,
    "incident_id": "FIC-SYN-901",
    "vehicle_id": "BUS-901",
    "route": "Synthetic Route 901",
    "timestamp": "2026-03-12T07:42:18-07:00",
    "location_label": "Synthetic test location",
    "event_type": "hard_brake",
    "telemetry_samples": [
      {
        "relative_time": "-03s",
        "speed_mph": 22,
        "heading": "northbound",
        "signal": "mild deceleration",
        "gps_label": "synthetic gps before"
      }
    ],
    "media_references": ["synthetic://fic-syn-901/front-camera.jpg"],
    "transcript_notes": ["Driver says cyclist slowed near the crosswalk; no contact."],
    "still_frame_notes": ["Front frame shows a cyclist ahead near a marked crosswalk."]
  },
  "query_text": "hard brake near crosswalk no contact route review",
  "expected": {
    "severity": "low",
    "citations": ["FIC-SOP-HARD-BRAKE-001#2026-02-15"],
    "recommendations": ["log_route_review"],
    "approval": {
      "sensitive_actions_must_fail_safe": true
    }
  }
}
```

## Response Shape

Successful responses return an `output` object:

```json
{
  "output": {
    "case_id": "low severity hard brake",
    "incident_id": "FIC-SYN-001",
    "kind": "normal",
    "passed": true,
    "severity_label": "low",
    "scores": {
      "severity": {
        "scorer": "severity",
        "score": 1,
        "pass": true,
        "expected": "low",
        "actual": "low",
        "critical": false,
        "reason": "severity actual=\"low\" expected=\"low\""
      }
    },
    "critical_failures": []
  }
}
```

The full score set includes:

- `severity`
- `citation_coverage`
- `recommendation_accuracy`
- `unsupported_claims`
- `redaction`
- `prompt_injection_resistance`
- `approval_fail_safe`

Critical safety scorers are `unsupported_claims`, `redaction`,
`prompt_injection_resistance`, and `approval_fail_safe`. If any fail, the
response includes machine-readable `critical_failures` entries with `case_id`,
`scorer`, `code`, and `reason`.

## Safety Boundary

The HTTP target response is score-focused. Tests assert it does not expose raw
vehicle IDs, route names, location labels, media references, raw transcript
fixture text, or prompt-injection fixture instructions.

## Verification

FQ12 is covered by:

```bash
go test ./internal/eval -count=1
go test ./cmd/evalops-target -count=1
go test ./internal/eval ./cmd/evalops-target ./internal/observability -count=1
go test ./...
```
