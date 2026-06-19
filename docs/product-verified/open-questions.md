# Open Questions

## Missing Source Artifacts

| ID | Missing Artifact | Impacted Scopes | Why It Blocks | What Is Needed |
| --- | --- | --- | --- | --- |
| Q-API-001 | API contract / GraphQL schema | feature-behavior, domain-model | Endpoints, payloads, errors, and auth transport cannot be defined from PRD alone | API contract or GraphQL SDL with query/mutation definitions |
| Q-AUTH-001 | Authorization and session policy | roles-permissions, edge-case-risk | No session TTL, no logout, no brute-force protection, no PIN recovery defined | Session management specification |
| Q-COMP-001 | Data retention / compliance policy | edge-case-risk | No retention, deletion, audit, or export rules for accumulated data | Data retention and compliance policy |
| Q-INT-001 | Integration contract for AI export | feature-behavior | No specification for manual vs automatic AI submission, no API integration | Integration specification for AI export workflow |

## Blocking

| ID | Question | Why It Matters | Resolution | Status |
| --- | --- | --- | --- | --- |
| Q-SCOPE-001 | What are the quantitative success metrics for the MVP? | Without success metrics, there is no way to validate the product is achieving its goals | Resolved in product-brief.md §Success Metrics (DEC-006) | resolved |
| Q-SCOPE-002 | Should the MVP architecture build for future multi-user or remain strictly single-user? | Affects data model, auth design, API structure, schema decisions | Resolved: single-user MVP with multi-user-ready data model; all entities owned via user_id; default user created at bootstrap (DEC-007) | resolved |
| Q-SCOPE-004 | What are the specific performance targets? | PRD §24.3 uses vague language | Resolved in product-brief.md §Performance Targets (DEC-008) | resolved |
| Q-SCOPE-005 | Should cardio be a separate entity or always part of a workout day? | Data model shows separate, text says part of workout day | Resolved: DailyLog replaces WorkoutDay; cardio is separate entity always attached to DailyLog via required dailyLogId (DEC-009) | resolved |

All blocking questions resolved. Open questions below are non-blocking design-time context for implementation.

## Non-Blocking

| ID | Question | Scope |
| --- | --- | --- |
| Q-SCOPE-003 | What is the target user's expected technical proficiency? | product-scope |
| Q-SCOPE-006 | What AI models/platforms must the export format support? | product-scope |
| Q-SCOPE-007 | Is there a maximum photo/media storage limit? | product-scope |
| Q-SCOPE-008 | What data portability standard is required? | product-scope |
| Q-ROLE-001 | PIN session lifetime and renewal policy | roles-permissions |
| Q-ROLE-002 | Logout mechanism when PIN is enabled | roles-permissions |
| Q-ROLE-003 | Access control when PIN is disabled | roles-permissions |
| Q-ROLE-005 | Deployer setup flow / configuration | roles-permissions |
| Q-ACTOR-01..28 | Various alternative paths, empty states, recovery flows, first-run experience | actor-journey |
| Q-DOMAIN-001..010 | Enum definitions for source, heartRateZone, cardioType, measurementType, side, flagType, mediaType, mealLabel, units, AiExport flags | domain-model |
| Q-FEAT-001..016 | Various feature behavior details (week counting, template lifecycle, cardio placement, etc.) | feature-behavior |
| Q-EDGE-01..12 | Edge cases (PIN brute-force, file size limits, concurrent access, storage, etc.) | edge-case-risk |
| Q-AC-01..18 | Acceptance criteria gaps (PIN failure behavior, e1RM formula, best set definition, timezone, etc.) | acceptance-criteria |
| Q-CONS-002 | Timezone handling for all date-based features | consistency-reviewer |

## Deferred

| ID | Question | Rationale |
| --- | --- | --- |
| Q-ROLE-004 | Resource-level visibility control (hide specific workouts) | Future scope consideration, not in MVP |