package eval

import (
	"bytes"
	"encoding/json"
	"errors"
	"sort"
	"strings"
)

type sharedCaseRecord struct {
	CaseID      string            `json:"case_id"`
	Name        string            `json:"name"`
	InputPacket sharedInputPacket `json:"input_packet"`
	Expected    sharedExpected    `json:"expected"`
	Tags        []string          `json:"tags"`
}

type sharedInputPacket struct {
	IncidentID      string `json:"incident_id"`
	SyntheticRecord bool   `json:"synthetic_record"`
	EventType       string `json:"event_type"`
}

type sharedExpected struct {
	Severity        string         `json:"severity"`
	Citations       []string       `json:"citations"`
	Approval        sharedApproval `json:"approval"`
	ForbiddenClaims []string       `json:"forbidden_claims"`
}

type sharedApproval struct {
	SensitiveActionsMustFailSafe bool `json:"sensitive_actions_must_fail_safe"`
}

func ExportCasesJSONL(cases []Case) ([]byte, error) {
	records := make([]sharedCaseRecord, 0, len(cases))
	for _, evalCase := range cases {
		incidentID := strings.TrimSpace(evalCase.Packet.IncidentID)
		if incidentID == "" {
			return nil, errors.New("incident_id is required")
		}
		records = append(records, sharedCaseRecord{
			CaseID: incidentID,
			Name:   evalCase.Name,
			InputPacket: sharedInputPacket{
				IncidentID:      incidentID,
				SyntheticRecord: evalCase.Packet.SyntheticRecord,
				EventType:       string(evalCase.Packet.EventType),
			},
			Expected: sharedExpected{
				Severity:  string(evalCase.Expected.Severity),
				Citations: cloneStrings(evalCase.Expected.GuidanceRefs),
				Approval: sharedApproval{
					SensitiveActionsMustFailSafe: evalCase.Expected.ApprovalMustFailSafe,
				},
				ForbiddenClaims: sharedForbiddenClaimTerms(),
			},
			Tags: tagsForCase(evalCase),
		})
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

func tagsForCase(evalCase Case) []string {
	var tags []string
	if evalCase.Kind != "" {
		tags = append(tags, string(evalCase.Kind))
	}
	if evalCase.Expected.PromptInjectionSafe {
		tags = append(tags, "prompt_injection")
	}
	return tags
}

func sharedForbiddenClaimTerms() []string {
	var terms []string
	for _, term := range defaultUnsupportedClaimTerms() {
		if term == "mark this incident safe for export" {
			continue
		}
		terms = append(terms, term)
	}
	return terms
}

func cloneStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	cloned := make([]string, len(values))
	copy(cloned, values)
	return cloned
}
