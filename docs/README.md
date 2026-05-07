# Documentation

This directory keeps the research and MVP planning artifacts separate from the repo root.

## Research

- [Research Report](research/research-report.md): source role and product research that shaped the Fleet Incident Copilot MVP.

## Fleet Incident Copilot MVP

- [MVP Overview](mvp/README.md): thesis, demo promise, proof points, and artifact map.

## EvalOps Overlays

- [EvalOps Extension](overlays/evalops-extension.md): FQ11-FQ15 overlay for reusable eval contracts, Promptfoo bridge work, trace and score export, release gates, and review loops.
- [EvalOps Shared Case And Result Contract](overlays/evalops-shared-contract.md): FQ11 JSONL case schema, result importer shape, golden-case mapping, and verification commands.

## Contributor Guides

- [How-Tos](how-tos.md): common local workflows for scope checks, package changes, synthetic fixtures, documentation updates, and review prep.
- [Developer Guide](developer-guide.md): repository layout, package boundaries, development workflow, and documentation rules.
- [Nice To Knows](nice-to-knows.md): context that prevents common wrong assumptions about runtime, demo, retrieval, eval, observability, and notification-preview behavior.
- [Testing](testing.md): targeted and full Go test commands, TDD expectations, coverage notes, and doc-only verification.
- [Troubleshooting](troubleshooting.md): fixes for common local development, packet, retrieval, timeline, severity, brief, approval, eval, observability, and local demo API issues.

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
- [Review Composition Contract](mvp/demo/review-composition-contract.md)
- [Loopback Demo API](mvp/demo/loopback-demo-api.md)
- [Dry-Run Slack-Shaped Notification Preview](mvp/demo/dry-run-slack-preview.md)
- [Scoped Approval Demo Retry](mvp/demo/scoped-approval-retry.md)
- [Eval And Observability Demo Reports](mvp/demo/eval-and-observability-reports.md)

Quality and operations:

- [Eval Plan](mvp/quality/eval-plan.md)
- [Observability And Cost Controls](mvp/quality/observability-and-cost-controls.md)

Execution and packaging:

- [Phases And Tasks](mvp/execution/phases.md)
- [Task Prompts](mvp/execution/task-prompts.md)
- [Strict TDD Rules](mvp/execution/tdd-rules.md)
- [Demo Package](mvp/demo/demo-package.md): refreshed local API walkthrough, response highlights, fallback tests, and interview packaging.
- [Demo Surface Roadmap](mvp/demo/demo-surface-roadmap.md)
- [LinkedIn Post Drafts](mvp/demo/linkedin-post-drafts.md): casual repo posts for hiring managers, CTOs, and engineering leaders.
