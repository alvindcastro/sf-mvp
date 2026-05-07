# EvalOps Release Gates

FQ14 adds a deterministic release-gate layer over the existing synthetic eval
scores. The gate is local and CI-ready: it does not call model providers,
Promptfoo, external telemetry, Slack, cloud services, or live fleet systems.

## Code Surface

- `eval.EvaluateReleaseGate(...)` checks Promptfoo/EvalOps-style score output
  against release thresholds.
- `eval.DefaultReleaseGateConfig()` requires zero critical failures and perfect
  deterministic severity, citation, and recommendation scores by default.
- `eval.ReleaseGateMarkdownSummary(...)` renders a deterministic Markdown
  summary suitable for `GITHUB_STEP_SUMMARY`.
- `eval.PromptfooOutputsFromReport(...)` converts the current in-memory eval
  report into gate-ready score output.
- `cmd/evalops-gate` runs the local golden-case gate by default, or imports a
  shared Promptfoo/EvalOps JSON result file with `-input`.
- `make evalops` runs the FQ14 package and command checks.
- `make evalops-gate` runs the local release gate.
- `.github/workflows/evalops.yml` wires those Make targets into GitHub Actions.

## Gate Rules

Default FQ14 gates require:

- zero critical failures;
- severity accuracy of `1.00`;
- citation coverage of `1.00`;
- recommendation accuracy of `1.00`;
- approval fail-safe passing for every case;
- redaction passing for every case;
- unsupported-claim checks passing for every case;
- prompt-injection resistance passing for every case.

Critical safety scorers are inherited from the Promptfoo score adapter:

- `unsupported_claims`
- `redaction`
- `prompt_injection_resistance`
- `approval_fail_safe`

Any critical failure exits non-zero. Warning-only scorer failures remain visible
in the Markdown summary but do not block when explicitly configured through the
command's `-warning-only` flag or package config.

## Local Commands

Run the release-gate test surface:

```bash
make evalops
```

Run the local golden-case release gate:

```bash
make evalops-gate
```

The command writes to `GITHUB_STEP_SUMMARY` automatically when that environment
variable is present:

```bash
GITHUB_STEP_SUMMARY=/tmp/evalops-summary.md make evalops-gate
```

Import a shared Promptfoo/EvalOps result file instead of running local golden
cases:

```bash
go run ./cmd/evalops-gate -input promptfoo-results.json -summary /tmp/evalops-summary.md
```

Exit codes:

- `0`: pass or warning-only findings;
- `1`: blocking release-gate failure;
- `2`: malformed input, malformed threshold config, or summary write failure.

## Summary Shape

The Markdown summary includes:

- overall `PASS`, `PASS WITH WARNINGS`, or `FAIL`;
- metric table for critical failures, severity accuracy, citation coverage,
  recommendation accuracy, approval fail-safe, redaction, unsupported claims,
  and prompt-injection resistance;
- blocking failure table with case IDs, scorer names, reasons, and remediation
  hints;
- warning table with the same fields;
- reproduce command.

## Safety Boundary

The release gate operates on scorer output, case IDs, incident IDs, reasons, and
remediation text. It does not include raw vehicle IDs, route names, location
labels, transcript notes, still-frame details, media references, or live
incident evidence in the summary.

## Verification

FQ14 is covered by:

```bash
go test ./internal/eval -count=1
go test ./cmd/evalops-gate -count=1
go test ./internal/eval ./cmd/evalops-target ./cmd/evalops-gate -count=1
make evalops
make evalops-gate
go test ./...
```
