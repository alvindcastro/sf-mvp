# LinkedIn Post Drafts

These drafts package the Fleet Incident Copilot repo for hiring managers, CTOs, engineering managers, and senior technical reviewers. They take writing cues from [The Elements of Style](../../eos.md): active voice, concrete nouns, short paragraphs, and direct claims.

Use each draft as a separate post. The subjects vary so the feed does not repeat the same repo summary every time. The tone should stay casual, practical, and builder-minded.

## Writing Guardrails

- Audience: hiring managers, CTOs, engineering managers, and senior technical reviewers.
- Voice: casual, specific, direct, and curious.
- Style: use short paragraphs, but avoid making every line a standalone sentence.
- Punctuation: do not use em dashes.
- AI positioning: AI is an assistive tool for drafting, reviewing, scoping, and pressure-testing. It is not the source of engineering judgment.
- Proof points: Go packages, strict TDD phases, deterministic local API behavior, approval gates, evals, traces, budget checks, and implementation docs.
- Avoid: production claims, customer-data claims, autonomous-agent hype, real Slack delivery claims, compliance claims, and vague "AI transformation" language.

## Post 1: AI As A Practical Assistant

I have been using AI while building a small Fleet Incident Copilot repo, but not in the way the louder demos usually frame it. The useful part has been assistance around scope, docs, stale wording, test ideas, and review questions.

The engineering judgment still has to come from the builder. The repo still needs narrow boundaries, deterministic behavior, tests, and honest documentation about what exists and what does not.

That is the relationship with AI that feels durable to me. It helps move the work along, then the code and docs have to stand on their own.

## Post 2: Trust Boundaries First

I wanted this repo to start with trust boundaries instead of starting with a chat interface. The Fleet Incident Copilot takes a synthetic incident packet, validates the input, retrieves approved mock guidance, builds a cited timeline, classifies severity, and drafts a redacted brief.

The system also blocks sensitive actions unless there is a scoped human approval. That part matters because the risky moment in an AI workflow is often not the summary, but the handoff to another system.

The current demo sends nothing and processes no real customer data. It shows the control points first, which is the order I would want before connecting anything real.

## Post 3: Strict TDD For AI Workflows

This project has been a useful way to practice strict TDD around an applied-AI-shaped workflow. Each phase adds a small Go package or demo route, and the behavior gets pinned down before the implementation gets comfortable.

That means the repo is not only a narrative about grounding, approval, evals, and observability. It has tests around packet validation, retrieval, timelines, severity, brief drafting, approval gates, eval scoring, trace events, notification previews, and local HTTP routes.

For AI-heavy products, that habit matters. The system should not depend on a confident demo voice to explain whether it behaves safely.

## Post 4: Human Approval Is A Product Feature

One of the main things I wanted to model in the Fleet Incident Copilot is human approval as part of the product, not as a slide at the end. The demo can prepare a Slack-shaped notification preview from a redacted brief, but it blocks the preview until the exact incident, action, and target have been approved.

Even after approval, the local demo still sends nothing. The point is to show the shape of an integration while keeping the action safe and inspectable.

That is the kind of boring product detail I like in applied AI systems. It makes the system less flashy, but much easier to reason about.

## Post 5: Synthetic Data With Realistic Constraints

The repo uses synthetic fleet incident data on purpose. That lets the workflow explore realistic constraints without pretending to have real customers, real incidents, or real operational authority.

The synthetic packet still has to be validated before the rest of the system can use it. The review still needs citations, redaction, severity logic, recommendations, approval gates, eval checks, and trace records.

Using fake data does not mean the engineering is fake. It means the demo can focus on system shape before any real-world integration is justified.

## Post 6: Local Demo, Honest Limits

The current demo flow is intentionally local. I start a loopback-only API, review one synthetic incident, show the blocked Slack-shaped preview, create a scoped approval, retry the preview, then pull eval and trace reports.

There is no production API, database, live model call, Slack delivery, customer data, auth layer, or external observability backend. Those limits are not hidden because they are part of the point.

I would rather show a small system clearly than imply a bigger one. For hiring managers and CTOs, that feels like a more useful signal than a polished demo with fuzzy edges.

## Post 7: What The Eval Proves

The eval package in this repo is small, but it is doing important work. It scores deterministic synthetic cases for severity, citation coverage, recommendation accuracy, unsupported claims, redaction leaks, prompt-injection resistance, and approval fail-closed behavior.

That does not make the repo production-ready. It does show the type of questions I want an AI workflow to answer before anyone asks it to touch a real process.

I like evals most when they are tied to product risk. In this repo, the risks are grounding, overclaiming, unsafe sharing, and unclear evidence.

## Post 8: Observability Before Production

I added local trace and budget-report surfaces because applied AI demos should show how they would be inspected. The repo records in-memory workflow events, token counts supplied by the caller, budget failures, tool-call records, approval decisions, and eval summaries.

This is not a monitoring platform. There is no external telemetry pipeline, dashboard, billing reconciliation, or persistent log store.

The value is in the shape of the evidence. If an AI workflow cannot explain what it used, what it decided, what it blocked, and what it cost, it is not ready to move closer to production.

## Post 9: Backend Product Thinking

The Fleet Incident Copilot is a backend-heavy repo, but the product thinking is still the main thing. The workflow starts with a fleet safety operator who needs to understand what happened, what evidence supports it, how severe it looks, and which next actions are safe.

That product frame changes the engineering choices. Citations, redactions, approval gates, evals, and trace reports become core behavior rather than extras.

This is the kind of applied AI work I enjoy. The interesting part is not only generating an answer, but designing the system around how that answer will be trusted.

## Post 10: Small Scope As A Strength

I kept this repo narrow because narrow systems are easier to test, explain, and improve. The Fleet Incident Copilot does one synthetic incident-review path rather than trying to become a full fleet platform.

That constraint made the work sharper. I could focus on validation, retrieval, timelines, severity, brief drafting, approval gates, evals, traces, and a local demo API.

Small scope is not a lack of ambition here. It is a way to make the engineering choices visible.

## Post 11: What I Would Show A CTO

If I were walking a CTO through this repo, I would not start with the output text. I would start with the boundaries around the output.

The system validates synthetic inputs, uses approved mock guidance, cites evidence, redacts the brief, blocks sensitive actions, records trace events, and exposes local eval results. It also states its limits clearly, including no live model call, no customer data, no Slack delivery, and no production API.

That is the conversation I want this repo to support. Not "look what the model wrote," but "look how the workflow is constrained."

## Post 12: What I Hope Hiring Managers Notice

I hope hiring managers see this repo as a signal of how I approach AI engineering. It is not a giant application, but it shows how I think about scope, safety, tests, docs, and product risk.

The repo uses Go, strict TDD phases, synthetic fixtures, deterministic local behavior, and a documentation trail that separates implemented behavior from planned work. AI helped me draft and review parts of the work, but the repo still has to communicate the engineering choices clearly.

That is the signal I want to send. I can use AI tools, but I do not want the tools to replace judgment.
