# Product Verified Output Contract

Create this structure exactly:

```text
docs/product-verified/
  index.md
  source-inventory.md
  product-brief.md
  scope.md
  actors-and-permissions.md
  domain-model.md
  functional-spec.md
  user-flows.md
  business-rules.md
  edge-cases.md
  acceptance-criteria.md
  open-questions.md
  features/
    index.md
    <feature-id>.md
  appendix/
    subagent-findings.md
    traceability.md
    derivation-log.md
    question-ledger.md
    decision-log.md
```

## Required File Purposes

- `index.md`: navigation, verification status, source set, and handoff readiness.
- `source-inventory.md`: input file list, source deltas, source coverage, stale/noisy material, and gaps.
- `product-brief.md`: product intent, target users, jobs to be done, value proposition, success metrics.
- `scope.md`: in scope, out of scope, non-goals, dependencies, assumptions.
- `actors-and-permissions.md`: actors, roles, permissions, ownership rules, privacy/security expectations.
- `domain-model.md`: entities, attributes, lifecycle states, relationships, identifiers, invariants.
- `functional-spec.md`: end-to-end behavior grouped by capability.
- `user-flows.md`: primary, alternative, failure, empty, and recovery flows.
- `business-rules.md`: policy, calculation, validation, lifecycle, notification, and integration rules.
- `edge-cases.md`: boundary and negative cases grouped by feature or domain area.
- `acceptance-criteria.md`: product-level and feature-level acceptance criteria.
- `open-questions.md`: unresolved blocking and non-blocking questions.
- `features/index.md`: feature inventory with status and links.
- `features/<feature-id>.md`: one file per material feature when there is enough source signal.
- `appendix/subagent-findings.md`: summary of each reviewer report and conflicts found.
- `appendix/traceability.md`: map verified requirements to source files, subagent reports, assumptions, or open questions.
- `appendix/derivation-log.md`: canonical list of derived roles, permissions, fields, states, acceptance criteria, and edge cases with source signal, rationale, and confidence.
- `appendix/question-ledger.md`: canonical question ledger copied from orchestration with statuses, resolutions, and links to final open questions.
- `appendix/decision-log.md`: decisions made while resolving contradictions and gaps.

## Required Headings

Use these headings at minimum.

### index.md

- `# Product Verified`
- `## Status`
- `## Source Set`
- `## Document Map`
- `## Handoff Readiness`

### source-inventory.md

- `# Source Inventory`
- `## Included Sources`
- `## Excluded Or Noisy Sources`
- `## Source Delta`
- `## Coverage Gaps`

### product-brief.md

- `# Product Brief`
- `## Product Intent`
- `## Target Users`
- `## Jobs To Be Done`
- `## Value Proposition`
- `## Success Metrics`

### scope.md

- `# Scope`
- `## In Scope`
- `## Out Of Scope`
- `## Non-Goals`
- `## Dependencies`
- `## Assumptions`

### actors-and-permissions.md

- `# Actors And Permissions`
- `## Actors`
- `## Roles`
- `## Permissions Matrix`
- `## Ownership Rules`
- `## Privacy And Security Expectations`

### domain-model.md

- `# Domain Model`
- `## Entities`
- `## Attributes`
- `## Relationships`
- `## Lifecycle States`
- `## Invariants`

### functional-spec.md

- `# Functional Specification`
- `## Capability Map`
- `## Feature Behavior`
- `## Validations`
- `## Notifications`
- `## Integrations`

### user-flows.md

- `# User Flows`
- `## Primary Flows`
- `## Alternative Flows`
- `## Failure And Recovery Flows`
- `## Empty States`

### business-rules.md

- `# Business Rules`
- `## Validation Rules`
- `## Calculation Rules`
- `## State Transition Rules`
- `## Authorization Rules`
- `## Integration Rules`

### edge-cases.md

- `# Edge Cases`
- `## Input And Validation`
- `## Permissions And Ownership`
- `## State And Concurrency`
- `## External Dependencies`
- `## Data Lifecycle`

### acceptance-criteria.md

- `# Acceptance Criteria`
- `## Product-Level Criteria`
- `## Feature-Level Criteria`
- `## Negative Criteria`
- `## Handoff Criteria`

### open-questions.md

- `# Open Questions`
- `## Missing Source Artifacts`
- `## Blocking`
- `## Non-Blocking`
- `## Deferred`

### appendix/traceability.md

- `# Traceability`
- `## Requirement Map`
- `## Source Map`
- `## Assumption Map`
- `## Open Question Map`

### appendix/derivation-log.md

- `# Derivation Log`
- `## Derived Roles And Permissions`
- `## Derived Data Fields`
- `## Derived States`
- `## Derived Acceptance Criteria`
- `## Derived Edge Cases`
- `## Low-Confidence Derivations`

### appendix/decision-log.md

- `# Decision Log`
- `## Resolved Contradictions`
- `## Assumptions Adopted`
- `## Rejected Or Outdated Inputs`

### appendix/question-ledger.md

- `# Question Ledger`
- `## Missing Source Artifacts`
- `## Blocking Questions`
- `## Non-Blocking Questions`
- `## Resolved Questions`
- `## Deferred Questions`

## Feature File Template

Each `features/<feature-id>.md` must include:

- `# <Feature Name>`
- `## Source Evidence`
- `## User Problem`
- `## Scope`
- `## Behavior`
- `## Derived Requirements`
- `## Edge Cases`
- `## Acceptance Criteria`
- `## Dependencies`
- `## Open Questions`

## Traceability Format

Use stable ids:

- Requirements: `REQ-001`, `REQ-002`, ...
- Business rules: `RULE-001`, `RULE-002`, ...
- Edge cases: `EDGE-001`, `EDGE-002`, ...
- Acceptance criteria: `AC-001`, `AC-002`, ...
- Decisions: `DEC-001`, `DEC-002`, ...
- Questions: `Q-001`, `Q-002`, ...
- Source-gap questions: `Q-API-001`, `Q-AUTH-001`, `Q-INT-001`, `Q-COMP-001`, ...

Each traceability row must identify one of:

- `Source: docs/product/path.md`
- `Subagent: .tasks/product-docs-verify/<run-id>/scopes/<scope>/...`
- `Derivation: appendix/derivation-log.md#<section>`
- `Assumption: DEC-###`
- `Open question: Q-###`

Every `REQ-###`, `RULE-###`, `EDGE-###`, and `AC-###` used in the final docs must appear in `appendix/traceability.md`.

## Derivation Format

Every derived item must include:

- stable id when it becomes a requirement, rule, edge case, or acceptance criterion;
- source signal;
- derivation rationale;
- confidence: `high`, `medium`, or `low`;
- linked open question when confidence is low or the item affects money, identity, authorization, compliance, irreversible transitions, or external contracts.
