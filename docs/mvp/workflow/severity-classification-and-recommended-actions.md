# Severity Classification And Recommended Actions

Phase 5 implements deterministic severity classification and recommended next actions for Fleet Incident Copilot. It consumes already-validated synthetic packet data, a Phase 4 timeline result, and Phase 3 retrieval citations.

## Phase 5 Checklist

- [x] Define low, medium, high, and unknown severity rules.
- [x] Prefer deterministic rule output for the MVP.
- [x] Isolate model-dependent judgment behind an interface if added later.
- [x] Tie recommendations to severity, event type, and retrieved SOPs.
- [x] Explain why each recommendation was made.
- [x] Flag export, escalation, and external sharing as approval-required.

## Runtime Surface

- Package: `internal/severity`.
- Entry point: `Classify(packet ingestion.Packet, timelineResult timeline.Result, guidance retrieval.Result) Result`.
- Targeted tests: `go test ./internal/severity`.
- Broad Go tests: `go test ./...`.

`Classify` is deterministic and in-memory. It does not call a model, external service, database, approval workflow, export tool, escalation tool, sharing tool, CLI, or HTTP API. The result includes `ModelJudgmentUsed: false` so downstream code can detect that Phase 5 used deterministic rules only.

## Severity Rules

The initial MVP rules are intentionally narrow and traceable to the Phase 1 synthetic packet contract:

| Packet signal | Severity | Rule |
| --- | --- | --- |
| `hard_brake` without timeline conflict | `low` | Controlled hard-brake event without collision or passenger-impact evidence. |
| `stop_arm_conflict` without collision evidence | `medium` | Stop-arm passing conflict needs supervisor review but does not prove impact. |
| `collision_signal` | `high` | Collision sensor and stop evidence require priority review. |
| `unknown_trigger` | `unknown` | Trigger subtype or evidence is incomplete. |
| timeline entry with conflicting telemetry uncertainty | `unknown` | Conflicting packet signals block deterministic severity. |
| `adversarial_note` | `medium` | Following-distance or hard-brake signals still require review, while hostile text stays untrusted data. |

The classifier does not infer injury, plate numbers, approval, export, escalation, or external sharing from packet text or retrieved snippets.

## Recommendation Output

Each recommendation includes:

- `Action`: deterministic action label such as `log_route_review`, `supervisor_review`, `preserve_media`, or `operator_review`.
- `Explanation`: why the recommendation was made, including severity and event-type context.
- `Sources`: structured packet references and retrieved guidance citation references.

Recommendations are SOP-grounded when a matching retrieved citation is present. The classifier only uses citations whose content role is `retrieved_data`; retrieved text remains data and cannot override rules or approval state.

## Approval Requirements

Every classification result flags these sensitive actions:

- `export`: approval required, not approved.
- `escalation`: approval required, not approved.
- `external_sharing`: approval required, not approved.

Phase 5 does not create approval records, approve actions, deny actions, execute tools, export evidence, escalate incidents, or share externally. The Phase 7 `internal/approval` package uses the same sensitive action labels when a human approval gate is needed.

## Red-To-Green Evidence

Smallest observable behavior: classify one low-severity hard-brake packet and return one SOP-grounded recommendation without using model judgment.

Observed red state:

- `go test ./internal/severity` failed because `Classify`, `Result`, `Level`, `Recommendation`, approval, and source types were undefined.

Green state:

- `go test ./internal/severity` passes after adding deterministic rules for low, medium, high, unknown, conflicting-signal, and adversarial cases; recommendation explanations with packet and guidance sources; and approval-required flags for export, escalation, and external sharing.

## Current Limits

- No model judgment interface is implemented because Phase 5 does not need model-dependent behavior.
- No CLI, HTTP API, database, persistent fixture loader, brief generator, export, escalation, external sharing, observability, or eval harness exists in Phase 5.
- Recommendations are advisory data for downstream review surfaces; they are not tool calls.
- Approval requirements are flags only; separate human approval decision records are created by the Phase 7 approval gate.
