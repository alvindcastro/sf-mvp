# EvalOps Extension — Fleet Incident Copilot

The MVP already has deterministic Go packages for ingestion, retrieval, timeline, severity, brief, approval, eval, and observability. This extension turns that package-level quality work into a reusable open-source-tool workflow:

```text
existing Go eval harness
    → shared evalops-go case/report contract
    → Promptfoo HTTP target for CI regression suites
    → OpenTelemetry trace/scores for production-style monitoring
```

## Non-negotiable code rule

Every code task follows [evalops-tdd-addendum.md](evalops-tdd-addendum.md) and the existing [Strict TDD Rules](../mvp/execution/tdd-rules.md).

## Phase overview

| Phase | Outcome |
|---|---|
| FQ11 — Shared eval contract | Existing synthetic cases export to a shared JSONL schema. |
| FQ12 — Promptfoo bridge | Promptfoo can call a local Go target and score incident workflows. |
| FQ13 — Trace and score export | Workflow steps emit safe OpenTelemetry/Langfuse-compatible attributes. |
| FQ14 — Release gates | CI blocks citation, approval, redaction, and unsupported-claim regressions. |
| FQ15 — Review loop | New demo failures become reviewed regression cases. |

---

## FQ11 — Shared eval contract

**Goal:** reuse the existing eval package through a cross-project case/result schema.

Output: [EvalOps Shared Case And Result Contract](evalops-shared-contract.md), `eval.ExportCasesJSONL`, `eval.ImportSharedResultsJSON`, and the golden fixture `internal/eval/testdata/evalops_cases.golden.jsonl`.

- [x] **FQ11-T01 — Map existing eval cases to shared schema**
  **Type:** Documentation
  **Done when:** every existing case maps `incident_id`, input packet, expected severity, expected citations, expected approval behavior, and forbidden claims.
- [x] **FQ11-T02 — Code: JSONL exporter for eval cases**
  **Prompt:** [FQ11-T02](evalops-task-prompts.md#fq11-t02-code-jsonl-exporter-for-eval-cases)
  **Done when:** tests prove deterministic JSONL output and no sensitive evidence leakage.
- [x] **FQ11-T03 — Code: shared result importer**
  **Prompt:** [FQ11-T03](evalops-task-prompts.md#fq11-t03-code-shared-result-importer)
  **Done when:** malformed results fail with useful errors and valid results can be scored by existing eval logic.

### FQ11 gate

- [x] Existing `internal/eval` behavior remains unchanged.
- [x] `go test ./internal/eval ./internal/observability` passes.
- [x] `go test ./...` passes.

---

## FQ12 — Promptfoo bridge

**Goal:** run incident workflow regression cases through Promptfoo without rewriting Go evals in JavaScript.

Output: [EvalOps Promptfoo Bridge](evalops-promptfoo-bridge.md), `eval.NewIncidentEvalTarget`, `eval.PromptfooOutputFromResult`, `cmd/evalops-target`, and `evals/promptfoo/fleet-incident.yaml`.

- [x] **FQ12-T01 — Code: local incident eval target**
  **Prompt:** [FQ12-T01](evalops-task-prompts.md#fq12-t01-code-local-incident-eval-target)
  **Done when:** `httptest` covers valid case, invalid case, missing packet, prompt-injection case, and timeout.
- [x] **FQ12-T02 — Add Promptfoo config**
  **Type:** Documentation/config
  **Files:** `evals/promptfoo/fleet-incident.yaml`
  **Done when:** config calls the local Go target and includes normal, incomplete, and adversarial cases.
- [x] **FQ12-T03 — Code: score adapter output**
  **Prompt:** [FQ12-T03](evalops-task-prompts.md#fq12-t03-code-score-adapter-output)
  **Done when:** Go returns scores for severity, citations, unsupported claims, redaction, and approval fail-closed behavior.

### FQ12 gate

- [x] Promptfoo run can be executed locally against deterministic Go code.
- [x] No model provider key is required.

---

## FQ13 — Trace and score export

**Goal:** make the incident workflow explainable in any OTel-compatible backend.

- [ ] **FQ13-T01 — Code: workflow span attributes**
  **Prompt:** [FQ13-T01](evalops-task-prompts.md#fq13-t01-code-workflow-span-attributes)
  **Done when:** tests cover trace ID, incident ID hashing, retrieval IDs, severity, approval status, and redaction.
- [ ] **FQ13-T02 — Code: eval score event export**
  **Prompt:** [FQ13-T02](evalops-task-prompts.md#fq13-t02-code-eval-score-event-export)
  **Done when:** score events are stable, safe, and can be correlated with trace IDs.
- [ ] **FQ13-T03 — Document Langfuse/OTel setup**
  **Type:** Documentation
  **Done when:** docs explain local OTel collector, Jaeger/Langfuse endpoint variables, and disabled-by-default behavior.

### FQ13 gate

- [ ] No raw incident evidence appears in spans.
- [ ] Telemetry disabled mode is tested.

---

## FQ14 — Release gates

**Goal:** prevent regressions before demo or deployment.

- [ ] **FQ14-T01 — Code: gate thresholds**
  **Prompt:** [FQ14-T01](evalops-task-prompts.md#fq14-t01-code-gate-thresholds)
  **Done when:** zero critical failures, citation coverage, approval fail-closed, and redaction gates are enforced.
- [ ] **FQ14-T02 — Add CI command**
  **Type:** Documentation/config
  **Done when:** `make evalops` and `make evalops-gate` are documented and CI-ready.
- [ ] **FQ14-T03 — Code: GitHub summary report**
  **Prompt:** [FQ14-T03](evalops-task-prompts.md#fq14-t03-code-github-summary-report)
  **Done when:** Markdown summary includes pass/fail table, failed case IDs, and remediation hints.

### FQ14 gate

- [ ] Seeded critical failure blocks release.
- [ ] Warning-only failures are visible but configurable.

---

## FQ15 — Review loop

**Goal:** convert demo/user-found issues into regression cases.

- [ ] **FQ15-T01 — Define issue-to-case workflow**
  **Type:** Documentation
  **Done when:** failures from manual review, demo rehearsal, or production-like traces become reviewed JSONL cases.
- [ ] **FQ15-T02 — Code: draft case generator**
  **Prompt:** [FQ15-T02](evalops-task-prompts.md#fq15-t02-code-draft-case-generator)
  **Done when:** trace/review samples convert to draft cases with TODO expected fields.
- [ ] **FQ15-T03 — Add monthly calibration checklist**
  **Type:** Documentation
  **Done when:** thresholds, prompts, fixtures, and allowed claims are periodically reviewed.
