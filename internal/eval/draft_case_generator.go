package eval

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"sf-mvp/internal/ingestion"
)

const DraftExpectedTODO = "TODO_REVIEW"

type ReviewTraceSample struct {
	CaseID              string
	Name                string
	IncidentID          string
	TraceID             string
	EventType           ingestion.EventType
	Tags                []string
	RedactedReviewNotes string
}

type DraftCaseRecord struct {
	CaseID         string           `json:"case_id"`
	Name           string           `json:"name"`
	SourceTraceID  string           `json:"source_trace_id"`
	InputPacket    DraftInputPacket `json:"input_packet"`
	Expected       DraftExpected    `json:"expected"`
	Tags           []string         `json:"tags"`
	ReviewRequired bool             `json:"review_required"`
	GateBlocking   bool             `json:"gate_blocking"`
}

type DraftInputPacket struct {
	IncidentID      string `json:"incident_id"`
	SyntheticRecord bool   `json:"synthetic_record"`
	EventType       string `json:"event_type"`
}

type DraftExpected struct {
	Severity        string        `json:"severity"`
	Citations       []string      `json:"citations"`
	Recommendations []string      `json:"recommendations"`
	Approval        DraftApproval `json:"approval"`
	ForbiddenClaims []string      `json:"forbidden_claims"`
}

type DraftApproval struct {
	SensitiveActionsMustFailSafe bool `json:"sensitive_actions_must_fail_safe"`
}

func ExportDraftCasesJSONL(samples []ReviewTraceSample) ([]byte, error) {
	records := make([]DraftCaseRecord, 0, len(samples))
	seen := make(map[string]struct{}, len(samples))
	for index, sample := range samples {
		record, err := draftCaseRecordFromSample(sample)
		if err != nil {
			return nil, fmt.Errorf("sample %d: %w", index, err)
		}
		if _, ok := seen[record.CaseID]; ok {
			continue
		}
		seen[record.CaseID] = struct{}{}
		records = append(records, record)
	}

	sort.SliceStable(records, func(i, j int) bool {
		return records[i].CaseID < records[j].CaseID
	})

	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	for _, record := range records {
		if err := encoder.Encode(record); err != nil {
			return nil, err
		}
	}
	return out.Bytes(), nil
}

func GenerateDraftCasesJSONL(samples []ReviewTraceSample) ([]byte, error) {
	return ExportDraftCasesJSONL(samples)
}

func draftCaseRecordFromSample(sample ReviewTraceSample) (DraftCaseRecord, error) {
	traceID := strings.TrimSpace(sample.TraceID)
	if traceID == "" {
		return DraftCaseRecord{}, errors.New("trace_id is required")
	}
	incidentID := strings.TrimSpace(sample.IncidentID)
	if incidentID == "" {
		return DraftCaseRecord{}, errors.New("incident_id is required")
	}
	if !strings.HasPrefix(incidentID, "FIC-SYN-") {
		return DraftCaseRecord{}, errors.New("incident_id must start with FIC-SYN-")
	}
	eventType := ingestion.EventType(strings.TrimSpace(string(sample.EventType)))
	if eventType == "" {
		return DraftCaseRecord{}, errors.New("event_type is required")
	}
	if !draftSupportedEventType(eventType) {
		return DraftCaseRecord{}, fmt.Errorf("event_type is not supported: %s", eventType)
	}

	caseID := firstNonEmpty(sample.CaseID, incidentID)
	return DraftCaseRecord{
		CaseID:        caseID,
		Name:          firstNonEmpty(sample.Name, caseID),
		SourceTraceID: traceID,
		InputPacket: DraftInputPacket{
			IncidentID:      incidentID,
			SyntheticRecord: true,
			EventType:       string(eventType),
		},
		Expected: DraftExpected{
			Severity:        DraftExpectedTODO,
			Citations:       []string{DraftExpectedTODO},
			Recommendations: []string{DraftExpectedTODO},
			Approval: DraftApproval{
				SensitiveActionsMustFailSafe: true,
			},
			ForbiddenClaims: []string{DraftExpectedTODO},
		},
		Tags:           draftTags(sample.Tags),
		ReviewRequired: true,
		GateBlocking:   false,
	}, nil
}

func draftTags(tags []string) []string {
	seen := map[string]struct{}{
		"draft":           {},
		"review_required": {},
	}
	values := []string{"draft", "review_required"}
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		values = append(values, tag)
	}
	sort.Strings(values)
	return values
}

func draftSupportedEventType(eventType ingestion.EventType) bool {
	switch eventType {
	case ingestion.EventTypeHardBrake,
		ingestion.EventTypeStopArmConflict,
		ingestion.EventTypeCollisionSignal,
		ingestion.EventTypeUnknownTrigger,
		ingestion.EventTypeAdversarialNote:
		return true
	default:
		return false
	}
}
