# Documentation

This directory keeps the research and MVP planning artifacts separate from the repo root.

## Research

- [Research Report](research/research-report.md): source role and product research that shaped the Fleet Incident Copilot MVP.

## Fleet Incident Copilot MVP

- [MVP Overview](mvp/README.md): thesis, demo promise, proof points, and artifact map.

## Contributor Guides

- [How-Tos](how-tos.md): common local workflows for scope checks, package changes, synthetic fixtures, documentation updates, and review prep.
- [Developer Guide](developer-guide.md): repository layout, package boundaries, development workflow, and documentation rules.
- [Nice To Knows](nice-to-knows.md): context that prevents common wrong assumptions about runtime, demo, retrieval, eval, and observability behavior.
- [Testing](testing.md): targeted and full Go test commands, TDD expectations, coverage notes, and doc-only verification.
- [Troubleshooting](troubleshooting.md): fixes for common local development, packet, retrieval, timeline, severity, brief, approval, eval, and observability issues.

Product and scope:

- [Product Frame](mvp/overview/product-frame.md)
- [Scope And Guardrails](mvp/overview/scope.md)

Workflow behavior:

- [Synthetic Incident Packets](mvp/workflow/synthetic-incident-packets.md)
- [Incident Packet Ingestion](mvp/workflow/incident-packet-ingestion.md)
- [RAG Corpus And Grounding](mvp/workflow/rag-corpus-and-grounding.md)
- [Incident Timeline Builder](mvp/workflow/incident-timeline-builder.md)
- [Severity Classification And Recommended Actions](mvp/workflow/severity-classification-and-recommended-actions.md)
- [Shareable Incident Brief Drafting](mvp/workflow/incident-brief-drafting.md)
- [Human Approval Workflow](mvp/workflow/human-approval-workflow.md)

Quality and operations:

- [Eval Plan](mvp/quality/eval-plan.md)
- [Observability And Cost Controls](mvp/quality/observability-and-cost-controls.md)

Execution and packaging:

- [Phases And Tasks](mvp/execution/phases.md)
- [Task Prompts](mvp/execution/task-prompts.md)
- [Strict TDD Rules](mvp/execution/tdd-rules.md)
- [Demo Package](mvp/demo/demo-package.md)
