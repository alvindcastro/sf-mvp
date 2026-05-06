# Strict TDD Rules

These rules apply to every future coding task in the Fleet Incident Copilot MVP.

## Non-Negotiable Checklist

- [ ] No production code before a failing test exists.
- [ ] Each task begins by naming the smallest observable behavior to test.
- [ ] The first test run must fail for the expected reason.
- [ ] Production changes must be limited to the smallest code needed to pass the failing test.
- [ ] Refactoring happens only after the targeted tests are green.
- [ ] Relevant broader tests run after the targeted tests pass.
- [ ] Documentation is updated when behavior, commands, architecture, or acceptance criteria change.
- [ ] The task summary includes tests run, red-to-green evidence, changed behavior, and residual risk.

## Required Coverage By Task Type

- [ ] Ingestion tasks cover valid input, malformed input, missing required fields, and audit output.
- [ ] Retrieval tasks cover relevant matches, no matches, citation metadata, scope filtering, and prompt-injection fixtures.
- [ ] Reasoning tasks cover grounded claims, missing data, conflicting data, and uncertainty labeling.
- [ ] Severity tasks cover low, medium, high, unknown, conflicting signals, and approval-required actions.
- [ ] Briefing tasks cover citation inclusion, redaction, missing evidence, and approval state.
- [ ] Approval tasks cover pending, denied, granted, out-of-scope, and immutable audit history.
- [ ] Eval tasks cover normal, adversarial, and incomplete fixtures.
- [ ] Observability tasks cover trace propagation, structured event fields, budget limits, and log redaction.

## Red / Green / Refactor Flow

- [ ] Red: add or update a focused test and run it.
- [ ] Red: confirm the failure is caused by missing behavior, not a broken test.
- [ ] Green: implement only enough production code to pass.
- [ ] Green: run the targeted test again.
- [ ] Refactor: simplify only after green, without broadening behavior.
- [ ] Verify: run the relevant broader suite.
- [ ] Record: summarize what changed and what was tested.

## Prohibited Shortcuts

- [ ] Do not write production code based only on a plan.
- [ ] Do not add broad scaffolding that is not required by the current failing test.
- [ ] Do not claim model quality without eval evidence.
- [ ] Do not use real fleet, student, passenger, driver, law-enforcement, or customer data.
- [ ] Do not connect live external services unless a future scope document explicitly allows it.
- [ ] Do not allow model output alone to approve, export, escalate, or externally share evidence.

