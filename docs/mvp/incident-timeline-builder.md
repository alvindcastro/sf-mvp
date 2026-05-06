# Incident Timeline Builder

Phase 4 implements the first deterministic timeline builder for Fleet Incident Copilot. It turns an already-validated synthetic packet into a chronological, cited set of incident timeline entries that can later feed severity, recommendation, and brief-drafting phases.

## Phase 4 Checklist

- [x] Order packet events and telemetry chronologically.
- [x] Incorporate transcript or still-frame notes without inventing visual facts.
- [x] Preserve source references for every factual claim.
- [x] Mark uncertainty when data is missing or conflicting.
- [x] Avoid unsupported claims.
- [x] Produce a timeline format suitable for the incident brief.

## Runtime Surface

- Package: `internal/timeline`.
- Entry point: `Build(packet ingestion.Packet, guidance retrieval.Result) Result`.
- Targeted tests: `go test ./internal/timeline`.
- Broad Go tests: `go test ./...`.

`Build` consumes an `ingestion.Packet`, not raw JSON. Phase 2 ingestion remains responsible for packet validation. `Build` may also consume `retrieval.Result` from Phase 3 so later workflow layers can carry approved guidance citation references alongside the timeline.

The builder is deterministic and in-memory. It does not call a model, external service, database, file loader, export tool, escalation tool, approval workflow, or brief generator.

## Timeline Output

Each timeline entry includes:

- `Time`: absolute event time.
- `Claim`: brief factual text.
- `Sources`: one or more structured source references.
- `Uncertain`: true when evidence is missing, unavailable, malformed, or conflicting.
- `Uncertainty`: a short uncertainty label.

The result also includes `GuidanceSources`, which are retrieval citation references available to downstream reasoning. Guidance citations are carried as metadata; retrieved snippets are not treated as instructions.

## Source Reference Rules

Every factual timeline entry must cite packet data with stable structured references:

- Telemetry: `packet.telemetry_samples[N]`.
- Transcript notes: `packet.transcript_notes[N]`.
- Still-frame notes: `packet.still_frame_notes[N]`.
- Media availability notes: `packet.media_references[N]`.

Retrieved guidance is cited with Phase 3 citation references such as `FIC-SOP-HARD-BRAKE-001#2026-02-15`. The builder only carries guidance citations whose content role is `retrieved_data`.

## Behavior

- Telemetry entries are ordered by `packet.Timestamp + telemetry_samples[N].relative_time`.
- Transcript entries are explicitly labeled as transcript notes.
- Still-frame entries are explicitly labeled as still-frame notes.
- The builder does not convert transcript text into visual claims.
- Media, transcript, or still-frame values containing `unavailable` produce uncertainty labels.
- Telemetry entries with the same absolute timestamp but conflicting claims are marked uncertain.
- Unsupported claims such as collision, injury, plate number, approval, export, or escalation are not generated unless they are present in cited packet evidence in a future phase.

## Red-To-Green Evidence

Smallest observable behavior: build a chronological timeline from telemetry samples while preserving source references.

Observed red state:

- `go test ./internal/timeline` failed because `Build`, `Entry`, and timeline result types were undefined.

Green state:

- `go test ./internal/timeline` passes after adding the minimal `internal/timeline` package with chronological telemetry ordering, packet source references, guidance citation references, transcript/still-frame attribution, unavailable-evidence uncertainty, and same-time telemetry conflict labeling.

## Current Limits

- No severity classification, recommendations, brief drafting, redaction, approval workflow, export, escalation, external sharing, persistence, CLI, HTTP API, observability, or eval harness exists in Phase 4.
- The builder does not retrieve documents on its own; callers pass retrieval results.
- The builder does not infer facts from images, transcripts, or guidance.
- The builder does not decide whether an incident is safe, approved, reportable, or externally shareable.
