package approval

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"sf-mvp/internal/severity"
)

type Decision string

const (
	DecisionPending  Decision = "pending"
	DecisionApproved Decision = "approved"
	DecisionDenied   Decision = "denied"
)

type AuditEventType string

const (
	AuditEventApprovalRequested      AuditEventType = "approval.requested"
	AuditEventApprovalDecided        AuditEventType = "approval.decided"
	AuditEventSensitiveActionBlocked AuditEventType = "sensitive_action.blocked"
	AuditEventSensitiveActionAllowed AuditEventType = "sensitive_action.allowed"
)

var (
	ErrActionBlocked         = errors.New("sensitive action blocked")
	ErrRequestAlreadyDecided = errors.New("approval request already decided")
)

type Scope struct {
	IncidentID string
	TargetRef  string
}

type RequestInput struct {
	IncidentID string
	Action     severity.SensitiveAction
	Scope      Scope
	Reason     string
}

type DecisionInput struct {
	RequestID string
	Approver  string
	Decision  Decision
	Reason    string
}

type SensitiveActionCall struct {
	IncidentID string
	Action     severity.SensitiveAction
	Scope      Scope
}

type Request struct {
	ID             string
	IncidentID     string
	Action         severity.SensitiveAction
	Scope          Scope
	Reason         string
	CreatedAt      time.Time
	Decision       Decision
	Approver       string
	DecisionReason string
	DecidedAt      time.Time
}

type ExecutionResult struct {
	Allowed   bool
	Executed  bool
	RequestID string
	Reason    string
}

type AuditEvent struct {
	Type       AuditEventType
	RequestID  string
	IncidentID string
	Action     severity.SensitiveAction
	Scope      Scope
	Actor      string
	Decision   Decision
	Reason     string
	OccurredAt time.Time
}

type Gate struct {
	now      func() time.Time
	nextID   int
	requests []Request
	audit    []AuditEvent
}

func NewGate(now func() time.Time) *Gate {
	if now == nil {
		now = time.Now
	}
	return &Gate{now: now}
}

func (g *Gate) CreateRequest(input RequestInput) (Request, error) {
	if err := validateRequestInput(input); err != nil {
		return Request{}, err
	}

	g.nextID++
	request := Request{
		ID:         fmt.Sprintf("approval-%03d", g.nextID),
		IncidentID: input.IncidentID,
		Action:     input.Action,
		Scope:      normalizeScope(input.IncidentID, input.Scope),
		Reason:     strings.TrimSpace(input.Reason),
		CreatedAt:  g.now(),
		Decision:   DecisionPending,
	}
	g.requests = append(g.requests, request)
	g.appendAudit(AuditEvent{
		Type:       AuditEventApprovalRequested,
		RequestID:  request.ID,
		IncidentID: request.IncidentID,
		Action:     request.Action,
		Scope:      request.Scope,
		Decision:   request.Decision,
		Reason:     request.Reason,
	})
	return request, nil
}

func (g *Gate) Decide(input DecisionInput) (Request, error) {
	requestIndex := g.findRequestIndex(input.RequestID)
	if requestIndex < 0 {
		return Request{}, fmt.Errorf("approval request %q not found", input.RequestID)
	}
	if strings.TrimSpace(input.Approver) == "" {
		return Request{}, errors.New("approver is required")
	}
	if input.Decision != DecisionApproved && input.Decision != DecisionDenied {
		return Request{}, errors.New("decision must be approved or denied")
	}
	if strings.TrimSpace(input.Reason) == "" {
		return Request{}, errors.New("decision reason is required")
	}

	request := g.requests[requestIndex]
	if request.Decision != DecisionPending {
		return request, ErrRequestAlreadyDecided
	}

	request.Approver = strings.TrimSpace(input.Approver)
	request.Decision = input.Decision
	request.DecisionReason = strings.TrimSpace(input.Reason)
	request.DecidedAt = g.now()
	g.requests[requestIndex] = request

	g.appendAudit(AuditEvent{
		Type:       AuditEventApprovalDecided,
		RequestID:  request.ID,
		IncidentID: request.IncidentID,
		Action:     request.Action,
		Scope:      request.Scope,
		Actor:      request.Approver,
		Decision:   request.Decision,
		Reason:     request.DecisionReason,
	})
	return request, nil
}

func (g *Gate) Execute(call SensitiveActionCall, execute func() error) (ExecutionResult, error) {
	if strings.TrimSpace(call.IncidentID) == "" || call.IncidentID != call.Scope.IncidentID {
		return g.block(call, "", DecisionPending, "requested incident is outside the supplied scope")
	}

	request, ok := g.matchingRequest(call)
	if !ok {
		return g.block(call, "", DecisionPending, "no approval exists within the requested scope")
	}

	switch request.Decision {
	case DecisionApproved:
		if execute == nil {
			return ExecutionResult{Allowed: true, RequestID: request.ID, Reason: "approved within scope"}, nil
		}
		if err := execute(); err != nil {
			return ExecutionResult{Allowed: true, RequestID: request.ID, Reason: err.Error()}, err
		}
		result := ExecutionResult{Allowed: true, Executed: true, RequestID: request.ID, Reason: "approved within scope"}
		g.appendAudit(AuditEvent{
			Type:       AuditEventSensitiveActionAllowed,
			RequestID:  request.ID,
			IncidentID: call.IncidentID,
			Action:     call.Action,
			Scope:      call.Scope,
			Decision:   DecisionApproved,
			Reason:     result.Reason,
		})
		return result, nil
	case DecisionDenied:
		return g.block(call, request.ID, DecisionDenied, "approval denied: "+request.DecisionReason)
	default:
		return g.block(call, request.ID, DecisionPending, "approval pending")
	}
}

func (g *Gate) AuditHistory() []AuditEvent {
	history := make([]AuditEvent, len(g.audit))
	copy(history, g.audit)
	return history
}

func (g *Gate) appendAudit(event AuditEvent) {
	event.OccurredAt = g.now()
	g.audit = append(g.audit, event)
}

func (g *Gate) block(call SensitiveActionCall, requestID string, decision Decision, reason string) (ExecutionResult, error) {
	result := ExecutionResult{
		Allowed:   false,
		Executed:  false,
		RequestID: requestID,
		Reason:    reason,
	}
	g.appendAudit(AuditEvent{
		Type:       AuditEventSensitiveActionBlocked,
		RequestID:  requestID,
		IncidentID: call.IncidentID,
		Action:     call.Action,
		Scope:      call.Scope,
		Decision:   decision,
		Reason:     reason,
	})
	return result, ErrActionBlocked
}

func (g *Gate) matchingRequest(call SensitiveActionCall) (Request, bool) {
	for i := len(g.requests) - 1; i >= 0; i-- {
		request := g.requests[i]
		if request.Action != call.Action {
			continue
		}
		if request.Scope != call.Scope {
			continue
		}
		return request, true
	}
	return Request{}, false
}

func (g *Gate) findRequestIndex(requestID string) int {
	for i, request := range g.requests {
		if request.ID == requestID {
			return i
		}
	}
	return -1
}

func validateRequestInput(input RequestInput) error {
	if strings.TrimSpace(input.IncidentID) == "" {
		return errors.New("incident_id is required")
	}
	if !isSensitiveAction(input.Action) {
		return fmt.Errorf("unsupported sensitive action %q", input.Action)
	}
	if strings.TrimSpace(input.Reason) == "" {
		return errors.New("approval request reason is required")
	}
	scope := normalizeScope(input.IncidentID, input.Scope)
	if strings.TrimSpace(scope.IncidentID) == "" || strings.TrimSpace(scope.TargetRef) == "" {
		return errors.New("approval scope requires incident_id and target_ref")
	}
	return nil
}

func normalizeScope(incidentID string, scope Scope) Scope {
	if strings.TrimSpace(scope.IncidentID) == "" {
		scope.IncidentID = strings.TrimSpace(incidentID)
	}
	scope.TargetRef = strings.TrimSpace(scope.TargetRef)
	return scope
}

func isSensitiveAction(action severity.SensitiveAction) bool {
	switch action {
	case severity.SensitiveActionExport, severity.SensitiveActionEscalation, severity.SensitiveActionExternalSharing:
		return true
	default:
		return false
	}
}
