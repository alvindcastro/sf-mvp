package eval

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"sf-mvp/internal/severity"
)

var allowedSharedScorers = map[string]struct{}{
	"severity":                    {},
	"citation_coverage":           {},
	"recommendation_accuracy":     {},
	"unsupported_claims":          {},
	"redaction":                   {},
	"prompt_injection_resistance": {},
	"approval_fail_safe":          {},
}

type sharedResultsEnvelope struct {
	Results json.RawMessage `json:"results"`
}

type sharedResultRecord struct {
	CaseID           string                 `json:"case_id"`
	IncidentID       string                 `json:"incident_id"`
	Kind             CaseKind               `json:"kind"`
	Vars             map[string]string      `json:"vars"`
	Scores           []sharedScoreAssertion `json:"scores"`
	AssertionResults []sharedScoreAssertion `json:"assertionResults"`
}

type sharedScoreAssertion struct {
	Scorer   string          `json:"scorer"`
	Metric   string          `json:"metric"`
	Score    json.RawMessage `json:"score"`
	Pass     bool            `json:"pass"`
	Expected string          `json:"expected"`
	Actual   string          `json:"actual"`
}

func ImportSharedResultsJSON(data []byte) ([]CaseResult, error) {
	records, err := decodeSharedResultRecords(data)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{}, len(records))
	results := make([]CaseResult, 0, len(records))
	for index, record := range records {
		caseID := firstNonEmpty(record.CaseID, record.Vars["case_id"])
		if strings.TrimSpace(caseID) == "" {
			return nil, fmt.Errorf("result %d missing case_id", index)
		}
		if _, ok := seen[caseID]; ok {
			return nil, fmt.Errorf("duplicate result for case_id %q", caseID)
		}
		seen[caseID] = struct{}{}

		result := CaseResult{
			Name:                                   caseID,
			IncidentID:                             firstNonEmpty(record.IncidentID, record.Vars["incident_id"]),
			Kind:                                   firstKind(record.Kind, CaseKind(record.Vars["kind"])),
			PromptInjectionResistant:               true,
			SensitiveActionsBlockedWithoutApproval: true,
			Passed:                                 true,
		}

		assertions := record.Scores
		if len(assertions) == 0 {
			assertions = record.AssertionResults
		}
		for _, assertion := range assertions {
			scorer := firstNonEmpty(assertion.Scorer, assertion.Metric)
			if _, ok := allowedSharedScorers[scorer]; !ok {
				return nil, fmt.Errorf("unknown scorer %q for case_id %q", scorer, caseID)
			}
			score, err := parseSharedScore(assertion.Score)
			if err != nil {
				return nil, fmt.Errorf("%s score for case_id %q: %w", scorer, caseID, err)
			}
			applySharedScore(&result, scorer, score, assertion)
		}
		result.Passed = len(result.Failures) == 0
		results = append(results, result)
	}
	return results, nil
}

func decodeSharedResultRecords(data []byte) ([]sharedResultRecord, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var envelope sharedResultsEnvelope
	if err := decoder.Decode(&envelope); err != nil {
		return nil, err
	}
	if len(envelope.Results) == 0 {
		return nil, errors.New("results are required")
	}

	var nested sharedResultsEnvelope
	if err := json.Unmarshal(envelope.Results, &nested); err == nil && len(nested.Results) > 0 {
		return decodeSharedResultArray(nested.Results)
	}
	return decodeSharedResultArray(envelope.Results)
}

func decodeSharedResultArray(data []byte) ([]sharedResultRecord, error) {
	var records []sharedResultRecord
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&records); err != nil {
		return nil, err
	}
	return records, nil
}

func parseSharedScore(raw json.RawMessage) (float64, error) {
	if len(raw) == 0 {
		return 0, errors.New("score is required")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var number json.Number
	if err := decoder.Decode(&number); err != nil {
		return 0, errors.New("score must be numeric")
	}
	score, err := number.Float64()
	if err != nil {
		return 0, errors.New("score must be numeric")
	}
	return score, nil
}

func applySharedScore(result *CaseResult, scorer string, score float64, assertion sharedScoreAssertion) {
	if !assertion.Pass {
		result.Failures = append(result.Failures, scorer+" failed")
	}

	switch scorer {
	case "severity":
		result.ExpectedSeverity = severity.Level(assertion.Expected)
		result.ActualSeverity = severity.Level(assertion.Actual)
		if assertion.Expected != "" && assertion.Actual != "" && assertion.Expected != assertion.Actual {
			result.Failures = append(result.Failures, fmt.Sprintf("severity = %q, want %q", assertion.Actual, assertion.Expected))
		}
	case "citation_coverage":
		result.CitationCoverage = score
	case "recommendation_accuracy":
		if !assertion.Pass {
			result.MissingRecommendations = []severity.RecommendationAction{"shared_result_reported_recommendation_failure"}
		}
	case "unsupported_claims":
		if !assertion.Pass {
			result.UnsupportedClaims = []string{"shared result reported unsupported claims"}
		}
	case "redaction":
		if !assertion.Pass {
			result.RedactionLeaks = []string{"shared result reported redaction failure"}
		}
	case "prompt_injection_resistance":
		result.PromptInjectionResistant = assertion.Pass
	case "approval_fail_safe":
		result.SensitiveActionsBlockedWithoutApproval = assertion.Pass
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func firstKind(values ...CaseKind) CaseKind {
	for _, value := range values {
		if strings.TrimSpace(string(value)) != "" {
			return value
		}
	}
	return ""
}
