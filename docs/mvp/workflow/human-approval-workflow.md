# Human Approval Workflow

Phase 7 implements the first human approval gate for Fleet Incident Copilot. It creates in-memory approval requests, records human decisions, and blocks sensitive action callbacks unless an approved request matches the requested action and scope.

## Phase 7 Checklist

- [x] Create approval requests for export, escalation, and external sharing.
- [x] Capture approver, timestamp, decision, reason, and target action.
- [x] Block sensitive tool calls while approval is pending.
- [x] Block denied actions.
- [x] Allow approved actions only within the approved scope.
- [x] Preserve append-only audit history.

## Runtime Surface

- Package: `internal/approval`.
- Gate constructor: `NewGate(now func() time.Time) *Gate`.
- Request entry point: `CreateRequest(input RequestInput) (Request, error)`.
- Decision entry point: `Decide(input DecisionInput) (Request, error)`.
- Enforcement entry point: `Execute(call SensitiveActionCall, execute func() error) (ExecutionResult, error)`.
- Audit entry point: `AuditHistory() []AuditEvent`.
- Targeted tests: `go test ./internal/approval`.
- Broad Go tests: `go test ./...`.

The gate is deterministic and in-memory. It does not call a model, external service, database, export tool, escalation tool, sharing tool, CLI, or HTTP API.

## Approval Request Model

`Request` records:

- `ID`: deterministic in-memory request ID.
- `IncidentID`: synthetic incident identifier.
- `Action`: one of the Phase 5 sensitive actions: `export`, `escalation`, or `external_sharing`.
- `Scope`: exact incident and target reference approved by the human workflow.
- `Reason`: request reason.
- `CreatedAt`: request timestamp from the injected clock.
- `Decision`: `pending`, `approved`, or `denied`.
- `Approver`: human approver captured when a decision is recorded.
- `DecisionReason`: human decision rationale.
- `DecidedAt`: decision timestamp from the injected clock.

Final decisions cannot be rewritten in place. A later changed decision must be represented by a new approval request and a new audit trail.

## Enforcement Rules

Sensitive action callbacks fail closed unless a matching approval exists:

- Missing approval blocks execution.
- Pending approval blocks execution.
- Denied approval blocks execution.
- Approved approval only allows the same action within the exact approved `Scope`.
- Out-of-scope incident or target references block execution even when another request is approved.
- The action call incident ID must match the supplied scope incident ID.
- Blocked calls return `ErrActionBlocked` and do not call the supplied function.

`Execute` is a gate around a callback only. Phase 7 does not implement real export, escalation, or external-sharing tools.

## Audit History

The gate appends audit events for:

- `approval.requested`
- `approval.decided`
- `sensitive_action.blocked`
- `sensitive_action.allowed`

Audit events include request ID, incident ID, action, scope, actor when available, decision, reason, and timestamp. `AuditHistory` returns a copy so callers cannot mutate stored audit records.

## Red-To-Green Evidence

Smallest observable behavior: create a pending approval request for each sensitive action and record an audit event before any production approval code exists.

Observed red states:

- `go test ./internal/approval` failed because `NewGate`, `Scope`, `RequestInput`, `DecisionPending`, `Gate`, `Request`, `AuditEventType`, and `Decision` were undefined.
- After the first green pass, an added final-decision immutability test failed because `ErrRequestAlreadyDecided` did not exist.
- A scoped-enforcement edge test failed because a callback could run when the call incident ID disagreed with the supplied approved scope.

Green state:

- `go test ./internal/approval` passes after adding the minimal in-memory approval gate, scoped execution enforcement, call-and-scope incident consistency checks, append-only audit-copy behavior, and final-decision immutability.

## Current Limits

- No persistence, database migrations, identity provider, role model, approval expiry, or policy engine exists.
- No real export, escalation, external-sharing, email, file-writing, webhook, CLI, or HTTP API integration exists.
- Audit history is append-only within one in-memory gate instance only.
- Observability and eval harness coverage remain future phases.
