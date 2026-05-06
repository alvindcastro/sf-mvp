# Incident Packet Ingestion

Phase 2 implements the first runtime gate for Fleet Incident Copilot. The ingestion package validates synthetic incident packet JSON before any future retrieval, timeline, severity, brief, export, or escalation step can use the packet.

## Phase 2 Checklist

- [x] Define packet schema and validation rules.
- [x] Reject missing incident IDs, timestamps, event types, telemetry arrays, or evidence references.
- [x] Reject malformed timestamps, impossible speed samples, unsupported event types, malformed telemetry sample fields, non-synthetic records, and non-synthetic evidence references.
- [x] Produce deterministic validation errors with field, code, and message values.
- [x] Emit an ingestion audit event for accepted and rejected packets.
- [x] Keep all examples synthetic.

## Runtime Surface

- Package: `internal/ingestion`.
- Entry point: `IngestJSON(data []byte) (Result, error)`.
- Targeted tests: `go test ./internal/ingestion`.
- Broad Go tests: `go test ./...`.

`IngestJSON` returns a `Result` whether validation succeeds or fails. Rejected packets return a `ValidationError` and an audit event with the validation issues attached.

## Packet Schema

Accepted JSON fields:

- `synthetic_record`: boolean marker. Must be `true`.
- `incident_id`: synthetic identifier. Required and must start with `FIC-SYN-`.
- `vehicle_id`: fake vehicle identifier. Required.
- `route`: fake route or service assignment. Required.
- `timestamp`: incident timestamp. Required and must parse as RFC 3339 with an offset.
- `location_label`: fake human-readable location. Required.
- `event_type`: controlled event label. Required.
- `telemetry_samples`: non-empty array of telemetry samples.
- `media_references`: non-empty array of evidence references.
- `transcript_notes`: optional synthetic text notes treated as data only.
- `still_frame_notes`: optional synthetic frame descriptions treated as data only.

Supported `event_type` values:

- `hard_brake`
- `stop_arm_conflict`
- `collision_signal`
- `unknown_trigger`
- `adversarial_note`

Telemetry sample fields:

- `relative_time`: required duration string such as `-06s`, `+00s`, or `+12s`.
- `speed_mph`: required number from `0` through `120`.
- `heading`: required synthetic heading label.
- `signal`: required synthetic signal label.
- `gps_label`: required synthetic GPS or location label.

Evidence reference rules:

- At least one `media_references` entry is required.
- Every media reference must use the `synthetic://` scheme.
- No real customer, driver, passenger, student, law-enforcement, route, GPS-coordinate, media, transcript, or evidence data is introduced.

## Validation Errors

Validation errors are deterministic and actionable:

- `Field`: exact packet field or indexed nested field.
- `Code`: stable machine-readable reason, such as `timestamp.malformed`.
- `Message`: short human-readable explanation.

Missing required-field errors are emitted in schema order so test expectations and future API responses remain stable.

## Audit Events

Accepted packet audit events use:

- Type: `incident_packet.ingested`.
- Accepted: `true`.
- Incident ID: parsed packet ID.

Rejected packet audit events use:

- Type: `incident_packet.rejected`.
- Accepted: `false`.
- Incident ID: parsed packet ID when available.
- Validation errors: same deterministic error collection returned by `IngestJSON`.

Audit events are in-memory return values in Phase 2. Persistent audit history is reserved for a later workflow or observability phase.

## Red-To-Green Evidence

Smallest observable behavior: accept one valid synthetic packet and reject invalid packet shapes before any downstream reasoning step can use them.

Observed red states:

- `go test ./internal/ingestion` failed because `IngestJSON`, event constants, validation types, and audit constants were undefined.
- Follow-up red tests showed non-synthetic records and non-`synthetic://` evidence references were accepted.
- Final schema-tightening red tests showed required packet context fields, malformed incident IDs, incomplete telemetry samples, and missing telemetry speeds were not rejected yet.

Green state:

- `go test ./internal/ingestion` passes after adding the minimal parser, validation rules, and audit event surface.
- `go test ./...` passes for the current Go module.

## Current Limits

- No CLI, HTTP API, database, fixture file loader, RAG retrieval, timeline builder, severity classifier, brief generator, approval workflow, export, escalation, or external sharing exists in Phase 2.
- Transcript and still-frame notes are accepted as untrusted packet data; they do not trigger instructions or tool actions.
- Future phases must add new behavior through failing tests first.
