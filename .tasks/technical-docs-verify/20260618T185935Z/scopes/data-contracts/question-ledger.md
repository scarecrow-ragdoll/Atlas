# Data-Contracts Question Ledger

| ID | Scope | Severity | Parent | Question | Why It Matters | Needed Artifact Or Decision | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| TQ-DATA-001 | data-contracts | needs-owner-decision | DEC-007 | Which entities require explicit userId FK vs. parent-chain traversal for multi-user readiness? | DEC-007 says "all entities" but 6 child entities (WorkoutSet, BodyMeasurement, ProgressPhoto, NutritionTemplateItem, DailyNutritionOverrideItem) omit userId. | Explicit depth decision from architect/owner. | worker-attempt-2 | open | TBD |
| TQ-DATA-002 | data-contracts | dev-blocking | DEC-007 | What is the userId attribute on entities where key column says "id (userId FK)" but attributes omit userId? | 11 entities affected: Settings, UserProfile, Exercise, ExerciseMedia, BodyWeightEntry, BodyCheckIn, NutritionProduct, NutritionTemplate, DailyNutritionOverride, WeekFlag, AiExport, AiReview. | Attribute list completion for all affected entities. | worker-attempt-2 | open | TBD |
| TQ-DATA-003 | data-contracts | dev-blocking | DEC-007 | Does DailyNutritionOverride have a userId field? Key column says yes (userId FK), attributes say no. | Direct schema contradiction blocks implementation. | Resolution of DailyNutritionOverride attribute list. | worker-attempt-2 | open | TBD |
| TQ-DATA-004 | data-contracts | dev-blocking | none | What are the allowed values for BodyWeightEntry.source, BodyMeasurement.measurementType, and WeekFlag.flagType? | Blocks schema column types, validation rules, and seed data. | Enum definitions (suggested sources: functional-spec.md §13, §17-18). | worker-attempt-2 | open | TBD |
| TQ-DATA-005 | data-contracts | dev-blocking | none | Which migration tool, versioning strategy, rollback policy, and seed data / fixture format are used? | Blocks all schema changes, default user bootstrap, and test data setup. | Migration tool choice (Goose, Atlas, Prisma) + fixture format. | worker-attempt-2 | open | TBD |
| TQ-DATA-006 | data-contracts | deferred | none | What Redis usage patterns are expected for MVP? | Redis in stack but no usage pattern documented. | Redis usage spec (if MVP needs it). | worker-attempt-2 | deferred | TBD |
| TQ-DATA-007 | data-contracts | needs-owner-decision | none | How are orphaned media files handled when ExerciseMedia or ProgressPhoto is deleted? | Risk of disk bloat over time; affects backup size and cleanup operations. | Media lifecycle policy (auto-cleanup, manual, or none). | worker-attempt-2 | open | TBD |
| TQ-DATA-008 | data-contracts | watchlist | TQ-DATA-007 | What is the data retention and lifecycle policy (auto-delete old logs, exports vs. manual only)? | Non-blocking for MVP, but affects long-term storage planning. | Data retention schedule (months/years/manual). | worker-attempt-2 | open | TBD |
| TQ-DATA-009 | data-contracts | watchlist | none | How does backup import handle ID collisions, FK violations, and partial restore failure? | Dry-run validation stated but recovery path after partial failure not specified. | Import collision / partial-failure strategy. | worker-attempt-2 | open | TBD |
| TQ-DATA-010 | data-contracts | dev-blocking | DEC-009 | What are the correct entity names for the 3 stale "WorkoutDay" references (relationships §line 109, invariant 6, lifecycle states §line 121)? | Blocks clean code generation and schema naming if DailyLog rename is applied inconsistently. | Domain model rename sweep from WorkoutDay to DailyLog. | worker-attempt-2 | open | TBD |
| TQ-DATA-011 | data-contracts | watchlist | DEC-008 | What index strategy is needed to meet p95 query SLOs (daily log <=500ms, exercise history <=700ms, body metrics <=700ms, nutrition <=700ms, charts <=1.0s)? | Missing index design risks SLO failure under 5-year dataset. | Composite index plan: (userId, date), (userId, exerciseId, date), etc. | worker-attempt-2 | open | TBD |

## Summary
- **dev-blocking**: 5 questions (must resolve before approved-to-dev)
- **needs-owner-decision**: 2 questions
- **deferred**: 1 question
- **watchlist**: 3 questions
- **resolved**: 0 questions (initial run)