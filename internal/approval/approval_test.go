package approval

import (
	"errors"
	"strings"
	"testing"
	"time"

	"sf-mvp/internal/severity"
)

func TestCreatePendingApprovalRequestsForSensitiveActions(t *testing.T) {
	gate := NewGate(fixedClock())
	scope := Scope{IncidentID: "FIC-SYN-001", TargetRef: "brief:FIC-SYN-001"}

	for _, action := range []severity.SensitiveAction{
		severity.SensitiveActionExport,
		severity.SensitiveActionEscalation,
		severity.SensitiveActionExternalSharing,
	} {
		request, err := gate.CreateRequest(RequestInput{
			IncidentID: "FIC-SYN-001",
			Action:     action,
			Scope:      scope,
			Reason:     "operator requested sensitive workflow",
		})
		if err != nil {
			t.Fatalf("CreateRequest(%q) returned error: %v", action, err)
		}

		if request.ID == "" {
			t.Fatalf("CreateRequest(%q) returned empty ID", action)
		}
		if request.Action != action {
			t.Fatalf("request Action = %q, want %q", request.Action, action)
		}
		if request.Decision != DecisionPending {
			t.Fatalf("request Decision = %q, want %q", request.Decision, DecisionPending)
		}
		if !request.CreatedAt.Equal(testTime()) {
			t.Fatalf("request CreatedAt = %s, want %s", request.CreatedAt, testTime())
		}
		if request.Scope != scope {
			t.Fatalf("request Scope = %#v, want %#v", request.Scope, scope)
		}
	}

	events := gate.AuditHistory()
	if len(events) != 3 {
		t.Fatalf("AuditHistory length = %d, want 3", len(events))
	}
	for _, event := range events {
		if event.Type != AuditEventApprovalRequested {
			t.Fatalf("audit event Type = %q, want %q", event.Type, AuditEventApprovalRequested)
		}
		if event.Decision != DecisionPending {
			t.Fatalf("audit event Decision = %q, want %q", event.Decision, DecisionPending)
		}
	}
}

func TestDecisionCapturesApproverTimestampDecisionReasonAndTargetAction(t *testing.T) {
	gate := NewGate(fixedClock())
	request := createExportRequest(t, gate, "FIC-SYN-002", "brief:FIC-SYN-002")

	decided, err := gate.Decide(DecisionInput{
		RequestID: request.ID,
		Approver:  "fleet-safety-lead",
		Decision:  DecisionApproved,
		Reason:    "brief is redacted and citations are present",
	})
	if err != nil {
		t.Fatalf("Decide returned error: %v", err)
	}

	if decided.Decision != DecisionApproved {
		t.Fatalf("Decision = %q, want %q", decided.Decision, DecisionApproved)
	}
	if decided.Approver != "fleet-safety-lead" {
		t.Fatalf("Approver = %q, want fleet-safety-lead", decided.Approver)
	}
	if decided.DecisionReason != "brief is redacted and citations are present" {
		t.Fatalf("DecisionReason = %q", decided.DecisionReason)
	}
	if !decided.DecidedAt.Equal(testTime()) {
		t.Fatalf("DecidedAt = %s, want %s", decided.DecidedAt, testTime())
	}
	if decided.Action != severity.SensitiveActionExport {
		t.Fatalf("Action = %q, want export", decided.Action)
	}
	if decided.Scope.TargetRef != "brief:FIC-SYN-002" {
		t.Fatalf("Scope.TargetRef = %q, want brief:FIC-SYN-002", decided.Scope.TargetRef)
	}

	events := gate.AuditHistory()
	last := events[len(events)-1]
	if last.Type != AuditEventApprovalDecided {
		t.Fatalf("last audit Type = %q, want %q", last.Type, AuditEventApprovalDecided)
	}
	if last.Actor != "fleet-safety-lead" || last.Decision != DecisionApproved {
		t.Fatalf("last audit = %#v, want approver and approved decision", last)
	}
}

func TestExecuteBlocksSensitiveActionWhileApprovalIsPending(t *testing.T) {
	gate := NewGate(fixedClock())
	request := createExportRequest(t, gate, "FIC-SYN-003", "brief:FIC-SYN-003")

	executed := false
	result, err := gate.Execute(SensitiveActionCall{
		IncidentID: "FIC-SYN-003",
		Action:     severity.SensitiveActionExport,
		Scope:      request.Scope,
	}, func() error {
		executed = true
		return nil
	})

	if !errors.Is(err, ErrActionBlocked) {
		t.Fatalf("Execute error = %v, want ErrActionBlocked", err)
	}
	if executed {
		t.Fatal("sensitive action executed while approval was pending")
	}
	if result.Allowed || result.Executed {
		t.Fatalf("result = %#v, want blocked and not executed", result)
	}
	if !strings.Contains(result.Reason, "pending") {
		t.Fatalf("blocked reason = %q, want pending approval", result.Reason)
	}
	assertLastAudit(t, gate, AuditEventSensitiveActionBlocked, DecisionPending)
}

func TestExecuteBlocksDeniedAction(t *testing.T) {
	gate := NewGate(fixedClock())
	request := createExportRequest(t, gate, "FIC-SYN-004", "brief:FIC-SYN-004")
	_, err := gate.Decide(DecisionInput{
		RequestID: request.ID,
		Approver:  "fleet-safety-lead",
		Decision:  DecisionDenied,
		Reason:    "shareable brief still contains sensitive details",
	})
	if err != nil {
		t.Fatalf("Decide returned error: %v", err)
	}

	executed := false
	result, err := gate.Execute(SensitiveActionCall{
		IncidentID: "FIC-SYN-004",
		Action:     severity.SensitiveActionExport,
		Scope:      request.Scope,
	}, func() error {
		executed = true
		return nil
	})

	if !errors.Is(err, ErrActionBlocked) {
		t.Fatalf("Execute error = %v, want ErrActionBlocked", err)
	}
	if executed {
		t.Fatal("denied sensitive action executed")
	}
	if !strings.Contains(result.Reason, "denied") {
		t.Fatalf("blocked reason = %q, want denied approval", result.Reason)
	}
	assertLastAudit(t, gate, AuditEventSensitiveActionBlocked, DecisionDenied)
}

func TestExecuteAllowsApprovedActionWithinApprovedScope(t *testing.T) {
	gate := NewGate(fixedClock())
	request := createExportRequest(t, gate, "FIC-SYN-005", "brief:FIC-SYN-005")
	_, err := gate.Decide(DecisionInput{
		RequestID: request.ID,
		Approver:  "fleet-safety-lead",
		Decision:  DecisionApproved,
		Reason:    "redacted brief approved for the requested target only",
	})
	if err != nil {
		t.Fatalf("Decide returned error: %v", err)
	}

	executed := false
	result, err := gate.Execute(SensitiveActionCall{
		IncidentID: "FIC-SYN-005",
		Action:     severity.SensitiveActionExport,
		Scope:      request.Scope,
	}, func() error {
		executed = true
		return nil
	})

	if err != nil {
		t.Fatalf("Execute returned error for approved scoped action: %v", err)
	}
	if !executed {
		t.Fatal("approved sensitive action was not executed")
	}
	if !result.Allowed || !result.Executed {
		t.Fatalf("result = %#v, want allowed and executed", result)
	}
	if result.RequestID != request.ID {
		t.Fatalf("result RequestID = %q, want %q", result.RequestID, request.ID)
	}
	assertLastAudit(t, gate, AuditEventSensitiveActionAllowed, DecisionApproved)
}

func TestExecuteBlocksApprovedActionOutsideApprovedScope(t *testing.T) {
	gate := NewGate(fixedClock())
	request := createExportRequest(t, gate, "FIC-SYN-006", "brief:FIC-SYN-006")
	_, err := gate.Decide(DecisionInput{
		RequestID: request.ID,
		Approver:  "fleet-safety-lead",
		Decision:  DecisionApproved,
		Reason:    "only the requested incident brief may be exported",
	})
	if err != nil {
		t.Fatalf("Decide returned error: %v", err)
	}

	executed := false
	result, err := gate.Execute(SensitiveActionCall{
		IncidentID: "FIC-SYN-007",
		Action:     severity.SensitiveActionExport,
		Scope:      Scope{IncidentID: "FIC-SYN-007", TargetRef: "brief:FIC-SYN-007"},
	}, func() error {
		executed = true
		return nil
	})

	if !errors.Is(err, ErrActionBlocked) {
		t.Fatalf("Execute error = %v, want ErrActionBlocked", err)
	}
	if executed {
		t.Fatal("out-of-scope sensitive action executed")
	}
	if !strings.Contains(result.Reason, "scope") {
		t.Fatalf("blocked reason = %q, want scope explanation", result.Reason)
	}
	assertLastAudit(t, gate, AuditEventSensitiveActionBlocked, DecisionPending)
}

func TestExecuteBlocksMismatchedCallIncidentAndApprovedScope(t *testing.T) {
	gate := NewGate(fixedClock())
	request := createExportRequest(t, gate, "FIC-SYN-010", "brief:FIC-SYN-010")
	_, err := gate.Decide(DecisionInput{
		RequestID: request.ID,
		Approver:  "fleet-safety-lead",
		Decision:  DecisionApproved,
		Reason:    "approve only the requested incident brief",
	})
	if err != nil {
		t.Fatalf("Decide returned error: %v", err)
	}

	executed := false
	result, err := gate.Execute(SensitiveActionCall{
		IncidentID: "FIC-SYN-011",
		Action:     severity.SensitiveActionExport,
		Scope:      request.Scope,
	}, func() error {
		executed = true
		return nil
	})

	if !errors.Is(err, ErrActionBlocked) {
		t.Fatalf("Execute error = %v, want ErrActionBlocked", err)
	}
	if executed {
		t.Fatal("sensitive action executed with mismatched call incident and approved scope")
	}
	if !strings.Contains(result.Reason, "scope") {
		t.Fatalf("blocked reason = %q, want scope explanation", result.Reason)
	}
	assertLastAudit(t, gate, AuditEventSensitiveActionBlocked, DecisionPending)
}

func TestAuditHistoryIsAppendOnlyAndReturnedAsACopy(t *testing.T) {
	gate := NewGate(fixedClock())
	request := createExportRequest(t, gate, "FIC-SYN-008", "brief:FIC-SYN-008")
	_, err := gate.Decide(DecisionInput{
		RequestID: request.ID,
		Approver:  "fleet-safety-lead",
		Decision:  DecisionDenied,
		Reason:    "needs another reviewer",
	})
	if err != nil {
		t.Fatalf("Decide returned error: %v", err)
	}

	history := gate.AuditHistory()
	if len(history) != 2 {
		t.Fatalf("AuditHistory length = %d, want 2", len(history))
	}
	history[0].Type = "tampered"
	history[0].Reason = "mutated outside the gate"

	unchanged := gate.AuditHistory()
	if unchanged[0].Type != AuditEventApprovalRequested {
		t.Fatalf("stored audit Type = %q, want unchanged requested event", unchanged[0].Type)
	}
	if unchanged[0].Reason == "mutated outside the gate" {
		t.Fatal("AuditHistory returned mutable backing storage")
	}

	_, _ = gate.Execute(SensitiveActionCall{
		IncidentID: "FIC-SYN-008",
		Action:     severity.SensitiveActionExport,
		Scope:      request.Scope,
	}, func() error {
		t.Fatal("denied action should not execute")
		return nil
	})

	grown := gate.AuditHistory()
	if len(grown) != 3 {
		t.Fatalf("AuditHistory length after blocked call = %d, want 3", len(grown))
	}
	if grown[0].Type != AuditEventApprovalRequested || grown[1].Type != AuditEventApprovalDecided || grown[2].Type != AuditEventSensitiveActionBlocked {
		t.Fatalf("audit order = %#v, want requested, decided, blocked", grown)
	}
}

func TestFinalDecisionCannotBeRewrittenInPlace(t *testing.T) {
	gate := NewGate(fixedClock())
	request := createExportRequest(t, gate, "FIC-SYN-009", "brief:FIC-SYN-009")
	denied, err := gate.Decide(DecisionInput{
		RequestID: request.ID,
		Approver:  "fleet-safety-lead",
		Decision:  DecisionDenied,
		Reason:    "redaction review failed",
	})
	if err != nil {
		t.Fatalf("Decide returned error: %v", err)
	}

	rewritten, err := gate.Decide(DecisionInput{
		RequestID: request.ID,
		Approver:  "fleet-safety-lead",
		Decision:  DecisionApproved,
		Reason:    "overwrite denied decision",
	})

	if !errors.Is(err, ErrRequestAlreadyDecided) {
		t.Fatalf("second Decide error = %v, want ErrRequestAlreadyDecided", err)
	}
	if rewritten.Decision != DecisionDenied {
		t.Fatalf("rewritten Decision = %q, want original denied decision", rewritten.Decision)
	}

	executed := false
	_, _ = gate.Execute(SensitiveActionCall{
		IncidentID: "FIC-SYN-009",
		Action:     severity.SensitiveActionExport,
		Scope:      request.Scope,
	}, func() error {
		executed = true
		return nil
	})
	if executed {
		t.Fatal("action executed after denied request was rewritten")
	}

	history := gate.AuditHistory()
	if len(history) != 3 {
		t.Fatalf("AuditHistory length = %d, want requested, denied decision, blocked action", len(history))
	}
	if history[1].Decision != denied.Decision || history[1].Decision != DecisionDenied {
		t.Fatalf("decision audit = %#v, want denied and unchanged", history[1])
	}
}

func createExportRequest(t *testing.T, gate *Gate, incidentID, targetRef string) Request {
	t.Helper()

	request, err := gate.CreateRequest(RequestInput{
		IncidentID: incidentID,
		Action:     severity.SensitiveActionExport,
		Scope:      Scope{IncidentID: incidentID, TargetRef: targetRef},
		Reason:     "operator requested export",
	})
	if err != nil {
		t.Fatalf("CreateRequest returned error: %v", err)
	}
	return request
}

func assertLastAudit(t *testing.T, gate *Gate, eventType AuditEventType, decision Decision) {
	t.Helper()

	events := gate.AuditHistory()
	if len(events) == 0 {
		t.Fatal("AuditHistory is empty")
	}
	last := events[len(events)-1]
	if last.Type != eventType {
		t.Fatalf("last audit Type = %q, want %q", last.Type, eventType)
	}
	if last.Decision != decision {
		t.Fatalf("last audit Decision = %q, want %q", last.Decision, decision)
	}
}

func fixedClock() func() time.Time {
	return func() time.Time { return testTime() }
}

func testTime() time.Time {
	return time.Date(2026, time.May, 6, 10, 15, 0, 0, time.FixedZone("America/Vancouver", -7*60*60))
}
