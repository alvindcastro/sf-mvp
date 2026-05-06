package ingestion

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"
)

type EventType string

const (
	EventTypeHardBrake       EventType = "hard_brake"
	EventTypeStopArmConflict EventType = "stop_arm_conflict"
	EventTypeCollisionSignal EventType = "collision_signal"
	EventTypeUnknownTrigger  EventType = "unknown_trigger"
	EventTypeAdversarialNote EventType = "adversarial_note"
)

type AuditEventType string

const (
	AuditEventIncidentPacketIngested AuditEventType = "incident_packet.ingested"
	AuditEventIncidentPacketRejected AuditEventType = "incident_packet.rejected"
)

type Packet struct {
	SyntheticRecord  bool              `json:"synthetic_record"`
	IncidentID       string            `json:"incident_id"`
	VehicleID        string            `json:"vehicle_id"`
	Route            string            `json:"route"`
	Timestamp        time.Time         `json:"-"`
	LocationLabel    string            `json:"location_label"`
	EventType        EventType         `json:"event_type"`
	TelemetrySamples []TelemetrySample `json:"telemetry_samples"`
	MediaReferences  []string          `json:"media_references"`
	TranscriptNotes  []string          `json:"transcript_notes"`
	StillFrameNotes  []string          `json:"still_frame_notes"`
}

type TelemetrySample struct {
	RelativeTime string  `json:"relative_time"`
	SpeedMPH     float64 `json:"speed_mph"`
	Heading      string  `json:"heading"`
	Signal       string  `json:"signal"`
	GPSLabel     string  `json:"gps_label"`
}

type Result struct {
	Packet     Packet
	AuditEvent AuditEvent
}

type AuditEvent struct {
	Type             AuditEventType
	IncidentID       string
	Accepted         bool
	ValidationErrors ValidationError
}

type ValidationIssue struct {
	Field   string
	Code    string
	Message string
}

type ValidationError struct {
	Issues []ValidationIssue
}

func (e ValidationError) Error() string {
	if len(e.Issues) == 0 {
		return "packet validation failed"
	}

	codes := make([]string, len(e.Issues))
	for i, issue := range e.Issues {
		codes[i] = issue.Code
	}
	return "packet validation failed: " + strings.Join(codes, ", ")
}

func IngestJSON(data []byte) (Result, error) {
	var raw packetJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		validationErr := ValidationError{Issues: []ValidationIssue{{
			Field:   "packet",
			Code:    "packet.malformed_json",
			Message: "packet must be valid JSON",
		}}}
		return rejectedResult("", validationErr), validationErr
	}

	packet, validationErr := validatePacket(raw)
	if len(validationErr.Issues) > 0 {
		return Result{Packet: packet, AuditEvent: auditEvent(AuditEventIncidentPacketRejected, packet.IncidentID, false, validationErr)}, validationErr
	}

	return Result{Packet: packet, AuditEvent: auditEvent(AuditEventIncidentPacketIngested, packet.IncidentID, true, ValidationError{})}, nil
}

type packetJSON struct {
	SyntheticRecord  bool                  `json:"synthetic_record"`
	IncidentID       string                `json:"incident_id"`
	VehicleID        string                `json:"vehicle_id"`
	Route            string                `json:"route"`
	Timestamp        string                `json:"timestamp"`
	LocationLabel    string                `json:"location_label"`
	EventType        EventType             `json:"event_type"`
	TelemetrySamples []telemetrySampleJSON `json:"telemetry_samples"`
	MediaReferences  []string              `json:"media_references"`
	TranscriptNotes  []string              `json:"transcript_notes"`
	StillFrameNotes  []string              `json:"still_frame_notes"`
}

type telemetrySampleJSON struct {
	RelativeTime string   `json:"relative_time"`
	SpeedMPH     *float64 `json:"speed_mph"`
	Heading      string   `json:"heading"`
	Signal       string   `json:"signal"`
	GPSLabel     string   `json:"gps_label"`
}

func validatePacket(raw packetJSON) (Packet, ValidationError) {
	var issues []ValidationIssue

	if strings.TrimSpace(raw.IncidentID) == "" {
		issues = append(issues, issue("incident_id", "incident_id.required", "incident_id is required"))
	}
	if strings.TrimSpace(raw.IncidentID) != "" && !strings.HasPrefix(raw.IncidentID, "FIC-SYN-") {
		issues = append(issues, issue("incident_id", "incident_id.synthetic_prefix.required", "incident_id must start with FIC-SYN-"))
	}
	if strings.TrimSpace(raw.VehicleID) == "" {
		issues = append(issues, issue("vehicle_id", "vehicle_id.required", "vehicle_id is required"))
	}
	if strings.TrimSpace(raw.Route) == "" {
		issues = append(issues, issue("route", "route.required", "route is required"))
	}
	if !raw.SyntheticRecord {
		issues = append(issues, issue("synthetic_record", "synthetic_record.required", "synthetic_record must be true"))
	}
	if strings.TrimSpace(raw.Timestamp) == "" {
		issues = append(issues, issue("timestamp", "timestamp.required", "timestamp is required"))
	}
	if strings.TrimSpace(raw.LocationLabel) == "" {
		issues = append(issues, issue("location_label", "location_label.required", "location_label is required"))
	}
	if strings.TrimSpace(string(raw.EventType)) == "" {
		issues = append(issues, issue("event_type", "event_type.required", "event_type is required"))
	}
	if len(raw.TelemetrySamples) == 0 {
		issues = append(issues, issue("telemetry_samples", "telemetry_samples.required", "telemetry_samples must include at least one sample"))
	}
	if len(raw.MediaReferences) == 0 {
		issues = append(issues, issue("media_references", "media_references.required", "media_references must include at least one synthetic reference"))
	}

	var parsedTimestamp time.Time
	if strings.TrimSpace(raw.Timestamp) != "" {
		var err error
		parsedTimestamp, err = time.Parse(time.RFC3339, raw.Timestamp)
		if err != nil {
			issues = append(issues, issue("timestamp", "timestamp.malformed", "timestamp must be RFC3339 with an offset"))
		}
	}

	if strings.TrimSpace(string(raw.EventType)) != "" && !isSupportedEventType(raw.EventType) {
		issues = append(issues, issue("event_type", "event_type.unsupported", "event_type is not supported"))
	}

	telemetrySamples := make([]TelemetrySample, 0, len(raw.TelemetrySamples))
	for i, sample := range raw.TelemetrySamples {
		relativeTime := strings.TrimSpace(sample.RelativeTime)
		switch {
		case relativeTime == "":
			issues = append(issues, issue(
				fmt.Sprintf("telemetry_samples[%d].relative_time", i),
				fmt.Sprintf("telemetry_samples[%d].relative_time.required", i),
				"relative_time is required",
			))
		default:
			if _, err := time.ParseDuration(relativeTime); err != nil {
				issues = append(issues, issue(
					fmt.Sprintf("telemetry_samples[%d].relative_time", i),
					fmt.Sprintf("telemetry_samples[%d].relative_time.malformed", i),
					"relative_time must parse as a Go duration such as -06s or +12s",
				))
			}
		}
		var speedMPH float64
		if sample.SpeedMPH == nil {
			issues = append(issues, issue(
				fmt.Sprintf("telemetry_samples[%d].speed_mph", i),
				fmt.Sprintf("telemetry_samples[%d].speed_mph.required", i),
				"speed_mph is required",
			))
		} else {
			speedMPH = *sample.SpeedMPH
		}
		if sample.SpeedMPH != nil && (math.IsNaN(speedMPH) || math.IsInf(speedMPH, 0) || speedMPH < 0 || speedMPH > 120) {
			issues = append(issues, issue(
				fmt.Sprintf("telemetry_samples[%d].speed_mph", i),
				fmt.Sprintf("telemetry_samples[%d].speed_mph.impossible", i),
				"speed_mph must be between 0 and 120",
			))
		}
		if strings.TrimSpace(sample.Heading) == "" {
			issues = append(issues, issue(
				fmt.Sprintf("telemetry_samples[%d].heading", i),
				fmt.Sprintf("telemetry_samples[%d].heading.required", i),
				"heading is required",
			))
		}
		if strings.TrimSpace(sample.Signal) == "" {
			issues = append(issues, issue(
				fmt.Sprintf("telemetry_samples[%d].signal", i),
				fmt.Sprintf("telemetry_samples[%d].signal.required", i),
				"signal is required",
			))
		}
		if strings.TrimSpace(sample.GPSLabel) == "" {
			issues = append(issues, issue(
				fmt.Sprintf("telemetry_samples[%d].gps_label", i),
				fmt.Sprintf("telemetry_samples[%d].gps_label.required", i),
				"gps_label is required",
			))
		}
		telemetrySamples = append(telemetrySamples, TelemetrySample{
			RelativeTime: sample.RelativeTime,
			SpeedMPH:     speedMPH,
			Heading:      sample.Heading,
			Signal:       sample.Signal,
			GPSLabel:     sample.GPSLabel,
		})
	}
	for i, mediaReference := range raw.MediaReferences {
		if !strings.HasPrefix(mediaReference, "synthetic://") {
			issues = append(issues, issue(
				fmt.Sprintf("media_references[%d]", i),
				fmt.Sprintf("media_references[%d].synthetic_uri.required", i),
				"media reference must use the synthetic:// scheme",
			))
		}
	}

	packet := Packet{
		SyntheticRecord:  raw.SyntheticRecord,
		IncidentID:       raw.IncidentID,
		VehicleID:        raw.VehicleID,
		Route:            raw.Route,
		Timestamp:        parsedTimestamp,
		LocationLabel:    raw.LocationLabel,
		EventType:        raw.EventType,
		TelemetrySamples: telemetrySamples,
		MediaReferences:  raw.MediaReferences,
		TranscriptNotes:  raw.TranscriptNotes,
		StillFrameNotes:  raw.StillFrameNotes,
	}

	return packet, ValidationError{Issues: issues}
}

func isSupportedEventType(eventType EventType) bool {
	switch eventType {
	case EventTypeHardBrake, EventTypeStopArmConflict, EventTypeCollisionSignal, EventTypeUnknownTrigger, EventTypeAdversarialNote:
		return true
	default:
		return false
	}
}

func rejectedResult(incidentID string, validationErr ValidationError) Result {
	return Result{AuditEvent: auditEvent(AuditEventIncidentPacketRejected, incidentID, false, validationErr)}
}

func auditEvent(eventType AuditEventType, incidentID string, accepted bool, validationErr ValidationError) AuditEvent {
	return AuditEvent{
		Type:             eventType,
		IncidentID:       incidentID,
		Accepted:         accepted,
		ValidationErrors: validationErr,
	}
}

func issue(field, code, message string) ValidationIssue {
	return ValidationIssue{Field: field, Code: code, Message: message}
}
