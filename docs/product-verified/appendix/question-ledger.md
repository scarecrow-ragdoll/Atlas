# Question Ledger

## Missing Source Artifacts

| ID | Missing Artifact | Impacted Scopes | Why It Blocks | Status |
| --- | --- | --- | --- | --- |
| Q-API-001 | API contract / GraphQL schema | feature-behavior, domain-model | Endpoints, payloads, errors, and auth transport cannot be defined from PRD alone | open |
| Q-AUTH-001 | Authorization and session policy | roles-permissions, edge-case-risk | No session TTL, no logout, no brute-force protection, no PIN recovery defined | open |
| Q-COMP-001 | Data retention / compliance policy | edge-case-risk | No retention, deletion, audit, or export rules for accumulated data | open |
| Q-INT-001 | Integration contract for AI export | feature-behavior | No specification for manual vs automatic AI submission | open |

## Blocking Questions

| ID | Question | Scope | Status |
| --- | --- | --- | --- |
| Q-SCOPE-001 | What are the quantitative success metrics for the MVP? | product-scope | open |
| Q-SCOPE-002 | Should the MVP architecture build for future multi-user or remain strictly single-user? | product-scope | open |
| Q-SCOPE-004 | What are the specific performance targets? | product-scope | open |
| Q-SCOPE-005 | Should cardio be a separate entity or always part of a workout day? | product-scope | open |

## Non-Blocking Questions

See aggregate question ledger at .tasks/product-docs-verify/20260618T185935Z/question-ledger.md for the full list of 94 non-blocking questions across all scopes. Key categories:

- 4 product scope questions (Q-SCOPE-003, Q-SCOPE-006, Q-SCOPE-007, Q-SCOPE-008)
- 4 roles and permissions questions (Q-ROLE-001, Q-ROLE-002, Q-ROLE-003, Q-ROLE-005, minus deferred)
- 28 actor journey questions (Q-ACTOR-01..28)
- 10 domain model questions (Q-DOMAIN-001..010)
- 16 feature behavior questions (Q-FEAT-001..016)
- 12 edge case and risk questions (Q-EDGE-01..12)
- 18 acceptance criteria questions (Q-AC-01..18)
- 1 consistency question (Q-CONS-002)

## Resolved Questions

| ID | Question | Resolution | Status |
| --- | --- | --- | --- |
| Q-SCOPE-001 | What are the quantitative success metrics for the MVP? | Resolved in product-brief.md §Success Metrics (DEC-006) | resolved |
| Q-SCOPE-002 | Should the MVP architecture build for future multi-user or remain strictly single-user? | Single-user MVP with multi-user-ready data model (DEC-007) | resolved |
| Q-SCOPE-004 | What are the specific performance targets? | Resolved in product-brief.md §Performance Targets (DEC-008) | resolved |
| Q-SCOPE-005 | Should cardio be a separate entity or always part of a workout day? | DailyLog replaces WorkoutDay; cardio attached via required dailyLogId (DEC-009) | resolved |

## Deferred Questions

| ID | Question | Rationale | Status |
| --- | --- | --- | --- |
| Q-ROLE-004 | Resource-level visibility control (hide specific workouts) | Future scope, not in MVP | deferred |