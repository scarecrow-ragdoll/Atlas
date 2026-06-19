<!-- FILE: .agents/skills/verify-technical-docs/references/output-contract.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the final technical readiness package structure produced by verify-technical-docs. -->
<!--   SCOPE: Covers required files, headings, statuses, stable ids, question ledger format, and approval gate fields; excludes orchestration prompts. -->
<!--   DEPENDS: .agents/skills/verify-technical-docs/SKILL.md. -->
<!--   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Required Structure - Lists files that must exist in docs/technical-verified. -->
<!--   Required Headings - Lists minimum headings for each package file. -->
<!--   Status Model - Defines approved-to-dev, questions-open, and blocked. -->
<!--   Question Ledger Format - Defines stable question rows and loop metadata. -->
<!--   Approval Evidence Format - Defines machine-checkable evidence required for approved-to-dev. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added technical readiness output contract. -->
<!-- END_CHANGE_SUMMARY -->

# Technical Verified Output Contract

Create this structure exactly:

```text
docs/technical-verified/
  index.md
  source-inventory.md
  technical-brief.md
  architecture-and-boundaries.md
  data-contracts.md
  api-contracts.md
  auth-security-compliance.md
  integrations-and-events.md
  client-state-and-ux-contracts.md
  operations-observability.md
  testing-and-delivery.md
  implementation-slices.md
  open-questions.md
  features/
    index.md
    <feature-id>.md
  appendix/
    subagent-findings.md
    traceability.md
    question-ledger.md
    decision-log.md
    loop-history.md
```

## Required File Purposes

- `index.md`: status, source set, document map, and dev handoff readiness.
- `source-inventory.md`: product-verified inputs, source deltas, answered questions, excluded/noisy files, and coverage gaps. `## Included Sources` must list at least one existing concrete source path as a markdown bullet such as `- docs/product-verified/index.md`.
- `technical-brief.md`: implementation-relevant product summary, technical assumptions, constraints, non-goals, and readiness summary.
- `architecture-and-boundaries.md`: system boundaries, contexts, components, ownership, dependencies, and architecture unknowns.
- `data-contracts.md`: entities, identifiers, persistence, migrations, retention, imports/exports, seed data, fixtures, and data unknowns.
- `api-contracts.md`: API surfaces, request/response contracts, error shapes, pagination/filtering, idempotency, compatibility, and API unknowns.
- `auth-security-compliance.md`: identity, authorization, ownership, audit, privacy, compliance, secrets, and abuse-case unknowns.
- `integrations-and-events.md`: external systems, sync rules, async jobs, events, webhooks, retries, rate limits, and integration unknowns.
- `client-state-and-ux-contracts.md`: UI state machine, loading/empty/error/offline states, forms, accessibility, realtime, cache, and UX technical gaps.
- `operations-observability.md`: environments, config, deployment, SLOs, logs, metrics, traces, alerts, runbooks, backups, and operational gaps.
- `testing-and-delivery.md`: unit, contract, integration, e2e, fixtures, test data, coverage gates, rollout, and QA blockers.
- `implementation-slices.md`: proposed dev slices only when supported by sources or explicit decisions, with blockers and verification refs.
- `open-questions.md`: unresolved technical questions grouped by severity and scope.
- `features/index.md`: feature-level technical inventory.
- `features/<feature-id>.md`: one feature-level technical file when product-verified has material features.
- `appendix/subagent-findings.md`: summary of scope worker/reviewer outputs and conflicts.
- `appendix/traceability.md`: map technical requirements, questions, decisions, and slices to product sources, scope reports, or decisions.
- `appendix/question-ledger.md`: canonical technical question ledger with statuses, parent links, and resolutions.
- `appendix/decision-log.md`: explicit technical decisions, deferrals, superseded answers, and rationale.
- `appendix/loop-history.md`: run history showing which questions were answered and whether answers created follow-up questions.

## Required Headings

### index.md

- `# Technical Verified`
- `## Status`
- `## Source Set`
- `## Document Map`
- `## Dev Handoff Readiness`

### source-inventory.md

- `# Source Inventory`
- `## Included Sources`
- `## Source Delta`
- `## Answered Questions`
- `## Excluded Or Noisy Sources`
- `## Coverage Gaps`

### technical-brief.md

- `# Technical Brief`
- `## Product Signal`
- `## Technical Scope`
- `## Constraints`
- `## Assumptions`
- `## Readiness Summary`

### architecture-and-boundaries.md

- `# Architecture And Boundaries`
- `## System Context`
- `## Components`
- `## Ownership Boundaries`
- `## Dependencies`
- `## Architecture Questions`

### data-contracts.md

- `# Data Contracts`
- `## Entities And Identifiers`
- `## Persistence And Storage`
- `## Migrations`
- `## Retention And Privacy`
- `## Data Questions`

### api-contracts.md

- `# API Contracts`
- `## Surfaces`
- `## Requests And Responses`
- `## Error And Validation Contracts`
- `## Compatibility And Idempotency`
- `## API Questions`

### auth-security-compliance.md

- `# Auth Security Compliance`
- `## Identity`
- `## Authorization And Ownership`
- `## Auditability`
- `## Privacy And Compliance`
- `## Security Questions`

### integrations-and-events.md

- `# Integrations And Events`
- `## External Systems`
- `## Events And Jobs`
- `## Sync And Retry Rules`
- `## Rate Limits And Failure Handling`
- `## Integration Questions`

### client-state-and-ux-contracts.md

- `# Client State And UX Contracts`
- `## User Interface States`
- `## Form And Validation Behavior`
- `## Cache And Realtime Behavior`
- `## Accessibility And Localization`
- `## Client Questions`

### operations-observability.md

- `# Operations Observability`
- `## Environments And Config`
- `## Deployment And Rollout`
- `## Logs Metrics Traces`
- `## Alerts Runbooks Backups`
- `## Operations Questions`

### testing-and-delivery.md

- `# Testing And Delivery`
- `## Test Strategy`
- `## Fixtures And Test Data`
- `## Contract And E2E Coverage`
- `## Release Gates`
- `## Testing Questions`

### implementation-slices.md

- `# Implementation Slices`
- `## Slice Map`
- `## Dependencies`
- `## Blockers`
- `## Verification`

### open-questions.md

- `# Open Questions`
- `## Dev-Blocking`
- `## Needs Owner Decision`
- `## Deferred`
- `## Watchlist`
- `## Resolved This Run`

### appendix/traceability.md

- `# Traceability`
- `## Technical Requirement Map`
- `## Question Map`
- `## Decision Map`
- `## Slice Map`
- `## Source Map`

### appendix/question-ledger.md

- `# Question Ledger`
- `## Open Questions`
- `## Answered Questions`
- `## Follow-Up Questions`
- `## Resolved Questions`
- `## Deferred Questions`

### appendix/decision-log.md

- `# Decision Log`
- `## Technical Decisions`
- `## Deferrals`
- `## Superseded Answers`
- `## Rejected Assumptions`

### appendix/loop-history.md

- `# Loop History`
- `## Runs`
- `## Answered Question Effects`
- `## Follow-Up Blockers`
- `## Approval Gate History`

## Feature File Template

Each `features/<feature-id>.md` must include:

- `# <Feature Name>`
- `## Product Evidence`
- `## Technical Scope`
- `## Architecture Notes`
- `## Data Contracts`
- `## API Or Integration Contracts`
- `## Client States`
- `## Test Coverage`
- `## Open Questions`

## Status Model

Use exactly one status in `index.md`.

- `questions-open`: technical package exists but at least one open question remains or answer effects have not closed.
- `blocked`: required source package, scope report, reviewer approval, or validation evidence is missing.
- `approved-to-dev`: all required scopes and consistency are approved, no open dev-blocking or owner-decision questions remain, and answer deltas created no new blocking questions.

## Stable IDs

- Technical requirements: `TREQ-001`, `TREQ-002`, ...
- Technical decisions: `TDEC-001`, `TDEC-002`, ...
- Technical questions: `TQ-ARCH-001`, `TQ-DATA-001`, `TQ-API-001`, `TQ-AUTH-001`, `TQ-INT-001`, `TQ-CLIENT-001`, `TQ-OPS-001`, `TQ-TEST-001`, ...
- Technical gaps by missing artifact class: `TGAP-API-001`, `TGAP-AUTH-001`, ...
- Implementation slices: `TSLICE-001`, `TSLICE-002`, ...
- Test obligations: `TTEST-001`, `TTEST-002`, ...

## Question Ledger Format

Use this table shape in every ledger:

```text
| ID | Scope | Severity | Parent | Question | Why It Matters | Needed Artifact Or Decision | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| TQ-API-001 | api-contracts | dev-blocking | none | Which API contract is authoritative? | Blocks endpoint and client implementation. | OpenAPI, GraphQL SDL, RPC contract, or owner decision. | api-contracts worker attempt 1 | open | TBD |
```

Allowed severities:

- `dev-blocking`
- `needs-owner-decision`
- `deferred`
- `watchlist`

Allowed statuses:

- `open`
- `answered-by-source`
- `answered-by-user`
- `resolved-by-decision`
- `superseded`
- `deferred`

Any `open` question with severity `dev-blocking` or `needs-owner-decision` blocks `approved-to-dev`.

## Approval Evidence Format

`appendix/subagent-findings.md` must include this table under `## Reviewer Verdicts`:

```text
| Scope | Status | Reviewer Verdict | Report |
| --- | --- | --- | --- |
| architecture-boundaries | approved | approved | .tasks/technical-docs-verify/<run-id>/scopes/architecture-boundaries/review-attempt-<n>.md |
```

Rows are required for every required scope:

- `architecture-boundaries`
- `data-contracts`
- `api-contracts`
- `auth-security-compliance`
- `integrations-events`
- `client-state-ux`
- `operations-observability`
- `testing-delivery`
- `consistency-loop-reviewer`

`appendix/loop-history.md` must include this table under `## Approval Gate History`:

```text
| Gate | Status | Evidence |
| --- | --- | --- |
| required-scopes-approved | passed | Link to Reviewer Verdicts table. |
```

Rows are required for these gates:

- `required-scopes-approved`
- `consistency-approved`
- `source-deltas-reviewed`
- `answer-deltas-reviewed`
- `no-answer-spawned-blockers`
- `no-open-blocking-questions`

When `index.md` status is `approved-to-dev`, every required scope row must be `approved`, every required reviewer report value must point to an existing `.tasks/technical-docs-verify/<run-id>/scopes/<scope>/review-attempt-<n>.md` artifact, every required gate row must be `passed`, the source inventory must not contain `SOURCE_MISSING` or `SOURCE_EMPTY`, `## Included Sources` must contain at least one existing source file path, the `source-deltas-reviewed` gate must cite `source-delta.md` or source inventory review evidence, and the `answer-deltas-reviewed` gate must cite `Answered Question Effects`, concrete question ids, or `source-delta.md` evidence.

## Approval Gate

`approved-to-dev` is allowed only when:

- every required Phase 1 scope has reviewer-approved status;
- `consistency-loop-reviewer` is approved;
- source inventory is present and non-empty;
- aggregate question ledger has zero open `dev-blocking` questions;
- aggregate question ledger has zero open `needs-owner-decision` questions;
- every source delta has recorded affected-scope review evidence;
- every answered prior question has a recorded answer source and effect analysis;
- no answer in the current loop produced an open `dev-blocking` or `needs-owner-decision` follow-up question;
- final output passes `validate_technical_verified.py` without `--allow-placeholders`.
