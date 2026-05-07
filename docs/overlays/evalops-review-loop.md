# EvalOps Review Loop

FQ15 closes the local EvalOps loop by turning demo rehearsal, manual review, or
production-like trace findings into draft regression cases. The loop is still
synthetic-only and local: it does not ingest real fleet data, call model
providers, export telemetry, create GitHub issues, or add a persistent case
store.

## Code Surface

- `eval.ExportDraftCasesJSONL(samples []eval.ReviewTraceSample) ([]byte, error)`
  converts redacted review or trace samples into deterministic draft JSONL.
- `eval.ReviewTraceSample` carries the safe routing inputs needed for a draft:
  case ID, name, synthetic incident ID, trace ID, event type, and inherited
  tags.
- `eval.DraftCaseRecord` is the generated JSONL shape with TODO expected fields,
  `review_required=true`, and `gate_blocking=false`.
- `eval.DraftExpectedTODO` is the placeholder value that marks expected fields a
  human reviewer must complete before the case is promoted into the reviewed
  eval suite.

## Issue-To-Case Workflow

1. Capture a failure from manual review, demo rehearsal, or a production-like
   trace. Keep only synthetic incident IDs and redacted review context.
2. Create a `ReviewTraceSample` with the failure's case ID, trace ID, synthetic
   incident ID, event type, and useful inherited tags such as `manual_review`,
   `demo_rehearsal`, `prompt_injection`, or `missing_evidence`.
3. Generate draft JSONL with `ExportDraftCasesJSONL`.
4. Review the generated case. Fill in expected severity, citations,
   recommendations, and forbidden claims. Confirm the case uses only approved
   synthetic fixture data.
5. Promote the reviewed case into the shared eval JSONL or golden case source
   only after expected fields are complete and another reviewer agrees that the
   failure is a useful regression.
6. Run the focused eval and release-gate commands before marking the issue
   closed.

Draft cases are deliberately non-blocking. They carry `review_required=true`,
`gate_blocking=false`, `draft`, and `review_required` tags so release gates can
ignore them until a human promotes the case.

## Draft JSONL Shape

Each generated line is one JSON object:

```json
{
  "case_id": "demo-stop-arm-review",
  "name": "demo stop-arm review failure",
  "source_trace_id": "trace-fic-syn-215-20260507t161500z-001",
  "input_packet": {
    "incident_id": "FIC-SYN-215",
    "synthetic_record": true,
    "event_type": "stop_arm_conflict"
  },
  "expected": {
    "severity": "TODO_REVIEW",
    "citations": ["TODO_REVIEW"],
    "recommendations": ["TODO_REVIEW"],
    "approval": {
      "sensitive_actions_must_fail_safe": true
    },
    "forbidden_claims": ["TODO_REVIEW"]
  },
  "tags": ["demo_rehearsal", "draft", "review_required", "stop_arm"],
  "review_required": true,
  "gate_blocking": false
}
```

The generator deduplicates by `case_id`, preserves the first trace for duplicate
samples, sorts output deterministically, and normalizes inherited tags.

## Safety Boundary

Generated drafts include only redacted routing fields: case ID, name, trace ID,
synthetic incident ID, synthetic flag, event type, expected TODO placeholders,
tags, and draft gate flags.

The JSONL intentionally omits raw or redacted review notes, vehicle IDs, routes,
location labels, telemetry GPS labels, transcript notes, still-frame notes,
media references, model responses, reviewer free text, and production customer
data. Samples whose incident IDs do not start with `FIC-SYN-` are rejected.

## Monthly Calibration Checklist

Review this checklist monthly, or before a major demo:

- Thresholds: confirm release-gate thresholds still reflect the risk tolerance
  for citation coverage, severity accuracy, recommendation accuracy, redaction,
  unsupported claims, prompt-injection resistance, and approval fail-safe.
- Prompts and target contract: verify the Promptfoo target request and score
  adapter still match the cases under review.
- Fixtures: confirm golden and draft cases remain synthetic, minimal,
  deterministic, and free of raw incident evidence.
- Allowed claims: review forbidden-claim terms, supported citation references,
  and recommendation expectations against current SOP-style docs.
- Draft queue: promote reviewed cases, keep useful unresolved drafts visible,
  and remove duplicates that no longer add regression value.

## Verification

FQ15 is covered by:

```bash
go test ./internal/eval -run 'TestExportDraftCasesJSONL' -count=1
go test ./internal/eval -count=1
make evalops
make evalops-gate
go test ./...
```
