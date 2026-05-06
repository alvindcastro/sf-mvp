package ingestion

import (
	"errors"
	"strings"
	"testing"
)

func TestIngestJSONAcceptsValidSyntheticPacket(t *testing.T) {
	result, err := IngestJSON(validPacketJSON())
	if err != nil {
		t.Fatalf("IngestJSON returned error for valid packet: %v", err)
	}

	if result.Packet.IncidentID != "FIC-SYN-001" {
		t.Fatalf("IncidentID = %q, want FIC-SYN-001", result.Packet.IncidentID)
	}
	if result.Packet.EventType != EventTypeHardBrake {
		t.Fatalf("EventType = %q, want %q", result.Packet.EventType, EventTypeHardBrake)
	}
	if len(result.Packet.TelemetrySamples) != 2 {
		t.Fatalf("TelemetrySamples length = %d, want 2", len(result.Packet.TelemetrySamples))
	}
	if result.AuditEvent.Type != AuditEventIncidentPacketIngested {
		t.Fatalf("AuditEvent.Type = %q, want %q", result.AuditEvent.Type, AuditEventIncidentPacketIngested)
	}
	if !result.AuditEvent.Accepted {
		t.Fatal("AuditEvent.Accepted = false, want true")
	}
}

func TestIngestJSONRejectsMissingRequiredFields(t *testing.T) {
	_, err := IngestJSON([]byte(`{
		"synthetic_record": true,
		"telemetry_samples": [],
		"media_references": []
	}`))

	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want ValidationError", err)
	}

	assertValidationCodes(t, validationErr,
		"incident_id.required",
		"vehicle_id.required",
		"route.required",
		"timestamp.required",
		"location_label.required",
		"event_type.required",
		"telemetry_samples.required",
		"media_references.required",
	)
}

func TestIngestJSONRejectsNonSyntheticEvidence(t *testing.T) {
	packet := validPacketJSONWith(`"synthetic_record": false`)

	_, err := IngestJSON(packet)

	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want ValidationError", err)
	}

	assertValidationCodes(t, validationErr, "synthetic_record.required")
}

func TestIngestJSONRejectsMalformedIncidentID(t *testing.T) {
	packet := validPacketJSONWith(`"incident_id": "INC-001"`)

	_, err := IngestJSON(packet)

	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want ValidationError", err)
	}

	assertValidationCodes(t, validationErr, "incident_id.synthetic_prefix.required")
}

func TestIngestJSONRejectsMalformedTelemetry(t *testing.T) {
	packet := validPacketJSONWith(`"telemetry_samples": [
		{"relative_time":"-06s","speed_mph":24,"heading":"northbound","signal":"steady speed","gps_label":"Oak St block 1200"},
		{"relative_time":"+00s","speed_mph":-4,"heading":"northbound","signal":"hard brake threshold crossed","gps_label":"Oak St at Pine Ave"}
	]`)

	_, err := IngestJSON(packet)

	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want ValidationError", err)
	}

	assertValidationCodes(t, validationErr, "telemetry_samples[1].speed_mph.impossible")
}

func TestIngestJSONRejectsIncompleteTelemetrySample(t *testing.T) {
	packet := validPacketJSONWith(`"telemetry_samples": [
		{"relative_time":"before impact","speed_mph":24,"heading":"","signal":"","gps_label":""}
	]`)

	_, err := IngestJSON(packet)

	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want ValidationError", err)
	}

	assertValidationCodes(t, validationErr,
		"telemetry_samples[0].relative_time.malformed",
		"telemetry_samples[0].heading.required",
		"telemetry_samples[0].signal.required",
		"telemetry_samples[0].gps_label.required",
	)
}

func TestIngestJSONRejectsTelemetrySampleWithoutSpeed(t *testing.T) {
	packet := validPacketJSONWith(`"telemetry_samples": [
		{"relative_time":"-06s","heading":"northbound","signal":"steady speed","gps_label":"Oak St block 1200"}
	]`)

	_, err := IngestJSON(packet)

	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want ValidationError", err)
	}

	assertValidationCodes(t, validationErr, "telemetry_samples[0].speed_mph.required")
}

func TestIngestJSONRejectsUnsupportedEventType(t *testing.T) {
	packet := validPacketJSONWith(`"event_type": "lane_departure"`)

	_, err := IngestJSON(packet)

	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want ValidationError", err)
	}

	assertValidationCodes(t, validationErr, "event_type.unsupported")
}

func TestIngestJSONRejectsNonSyntheticEvidenceReference(t *testing.T) {
	packet := validPacketJSONWith(`"media_references": [
			"https://example.com/front-camera-074218.jpg"
		]`)

	_, err := IngestJSON(packet)

	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want ValidationError", err)
	}

	assertValidationCodes(t, validationErr, "media_references[0].synthetic_uri.required")
}

func TestIngestJSONEmitsAuditEventForRejectedPacket(t *testing.T) {
	result, err := IngestJSON(validPacketJSONWith(`"timestamp": "03/12/2026 07:42:18"`))
	if err == nil {
		t.Fatal("IngestJSON returned nil error for malformed timestamp")
	}

	if result.AuditEvent.Type != AuditEventIncidentPacketRejected {
		t.Fatalf("AuditEvent.Type = %q, want %q", result.AuditEvent.Type, AuditEventIncidentPacketRejected)
	}
	if result.AuditEvent.Accepted {
		t.Fatal("AuditEvent.Accepted = true, want false")
	}
	assertValidationCodes(t, result.AuditEvent.ValidationErrors, "timestamp.malformed")
}

func assertValidationCodes(t *testing.T, got ValidationError, want ...string) {
	t.Helper()

	if len(got.Issues) != len(want) {
		t.Fatalf("validation issue count = %d, want %d: %#v", len(got.Issues), len(want), got.Issues)
	}
	for i, code := range want {
		if got.Issues[i].Code != code {
			t.Fatalf("validation issue %d code = %q, want %q; all issues: %#v", i, got.Issues[i].Code, code, got.Issues)
		}
	}
}

func validPacketJSON() []byte {
	return []byte(`{
		"synthetic_record": true,
		"incident_id": "FIC-SYN-001",
		"vehicle_id": "BUS-214",
		"route": "North Loop School Route 7",
		"timestamp": "2026-03-12T07:42:18-07:00",
		"location_label": "Oak Street near Pine Avenue",
		"event_type": "hard_brake",
		"telemetry_samples": [
			{"relative_time":"-06s","speed_mph":24,"heading":"northbound","signal":"steady speed","gps_label":"Oak St block 1200"},
			{"relative_time":"+00s","speed_mph":9,"heading":"northbound","signal":"hard brake threshold crossed","gps_label":"Oak St at Pine Ave"}
		],
		"media_references": [
			"synthetic://fic-syn-001/front-camera-074218.jpg"
		],
		"transcript_notes": [
			"Driver says cyclist slowed near the crosswalk; no contact."
		],
		"still_frame_notes": [
			"Front frame shows a cyclist ahead in the bike lane near a marked crosswalk."
		]
	}`)
}

func validPacketJSONWith(replacement string) []byte {
	packet := string(validPacketJSON())

	switch {
	case replacement == `"synthetic_record": false`:
		packet = replaceJSONField(packet, `"synthetic_record": true`, replacement)
	case replacement == `"incident_id": "INC-001"`:
		packet = replaceJSONField(packet, `"incident_id": "FIC-SYN-001"`, replacement)
	case replacement == `"event_type": "lane_departure"`:
		packet = replaceJSONField(packet, `"event_type": "hard_brake"`, replacement)
	case replacement == `"timestamp": "03/12/2026 07:42:18"`:
		packet = replaceJSONField(packet, `"timestamp": "2026-03-12T07:42:18-07:00"`, replacement)
	case strings.HasPrefix(replacement, `"media_references"`):
		packet = replaceJSONField(packet, `"media_references": [
			"synthetic://fic-syn-001/front-camera-074218.jpg"
		]`, replacement)
	default:
		packet = replaceJSONField(packet, `"telemetry_samples": [
			{"relative_time":"-06s","speed_mph":24,"heading":"northbound","signal":"steady speed","gps_label":"Oak St block 1200"},
			{"relative_time":"+00s","speed_mph":9,"heading":"northbound","signal":"hard brake threshold crossed","gps_label":"Oak St at Pine Ave"}
		]`, replacement)
	}

	return []byte(packet)
}

func replaceJSONField(packet, old, replacement string) string {
	return strings.Replace(packet, old, replacement, 1)
}
