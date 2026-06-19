# Data-Contracts Worker Attempt 2

## Revisions Applied
Per review-attempt-1.md:
1. Added TQ-DATA-010 for stale WorkoutDay references
2. Extended TQ-DATA-005 to cover seed data / fixture format
3. Added TQ-DATA-011 for index strategy against p95 targets
4. Merged TQ-DATA-007 into broader TQ-DATA-008 (data lifecycle)

## Sources Read
Same as attempt 1: domain-model.md, scope.md, product-brief.md, functional-spec.md, decision-log.md, source-delta.md.

## Source Delta Reviewed
Same as attempt 1. No new deltas.

## Product Signals
Same as attempt 1.

## Technical Facts
Same as attempt 1 plus the following additions:

### DEC-009 Stale References
Three stale "WorkoutDay" references remain after DEC-009:
1. Relationships § line 109: "WorkoutDay 0:1 BodyWeightEntry" — should be DailyLog
2. Invariant 6: "WorkoutExercise.order determines display order within WorkoutDay" — should be DailyLog
3. Lifecycle states: "WorkoutDay has implicit states" — should be DailyLog

### Index Strategy Gap
Product-brief.md §Performance Targets defines p95 SLOs:
- Daily log query by date: <= 500ms
- Exercise history query: <= 700ms
- Body metrics query for period: <= 700ms
- Nutrition summary for period: <= 700ms
- Chart data query for period: <= 1.0s
These SLOs require index planning: composite indexes on (userId, date), (userId, exerciseId, date), etc.

### Seed Data Gap
Default user bootstrap described but no fixture file format, test data seeding strategy, or migration-anchored seed script specified.

## Technical Gaps
Same as attempt 1, with these additions:

### T11: WorkoutDay Legacy References
Three stale references to old entity name after DEC-009. Blocks schema naming and prevents clean code generation.

### T12: Index Strategy Missing
p95 SLOs require index design but no index strategy document exists.

### T13: Seed Data Format
Default user fixture and test data seeding format undefined.

## Missing Source Artifacts (Updated)
- SQL/DDL migrations ✓ (TQ-DATA-005)
- Seed data / fixture files ✓ (TQ-DATA-005 extended)
- Index strategy ✓ (TQ-DATA-011)
- Redis usage specification ✓ (TQ-DATA-006)
- Storage engine configuration — still unaddressed (connection pooling, timeouts, SSL)
- Media file lifecycle management ✓ (TQ-DATA-008)
- Data retention schedule ✓ (TQ-DATA-008)
- Import collision resolution ✓ (TQ-DATA-009)
- Enum definitions ✓ (TQ-DATA-004)
- Schema version format ✓ (ask to consolidate into TQ-DATA-005)
- WorkoutDay → DailyLog rename ✓ (TQ-DATA-010)

## Questions Raised (Updated)

| ID | Scope | Severity | Parent | Question | Why It Matters | Needed Artifact Or Decision | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| TQ-DATA-001 | data-contracts | needs-owner-decision | DEC-007 | Which entities require explicit userId FK vs. parent-chain traversal? | DEC-007 ambiguous for 6 child entities. | Explicit depth decision. | worker-attempt-2 | open | TBD |
| TQ-DATA-002 | data-contracts | dev-blocking | DEC-007 | What is the userId attribute on entities that declare "id (userId FK)" but omit userId in attributes? | 11 entities have FK in key col but no userId attribute. | Attribute list completion. | worker-attempt-2 | open | TBD |
| TQ-DATA-003 | data-contracts | dev-blocking | DEC-007 | Does DailyNutritionOverride have userId? Key col says yes, attributes say no. | Direct contradiction. | Attribute list resolution. | worker-attempt-2 | open | TBD |
| TQ-DATA-004 | data-contracts | dev-blocking | none | What are the allowed values for BodyWeightEntry.source, BodyMeasurement.measurementType, and WeekFlag.flagType? | Blocks schema, validation, seed data. | Enum definitions. | worker-attempt-2 | open | TBD |
| TQ-DATA-005 | data-contracts | dev-blocking | none | Which migration tool, versioning strategy, rollback policy, and seed data / fixture format are used? | Blocks all schema changes, bootstrap, test data. | Migration tool choice + fixture format. | worker-attempt-2 | open | TBD |
| TQ-DATA-006 | data-contracts | deferred | none | What Redis usage patterns are expected for MVP? | Blocks cache/session/queue design. | Redis usage spec (if MVP requires it). | worker-attempt-2 | deferred | TBD |
| TQ-DATA-007 | data-contracts | needs-owner-decision | none | How are orphaned media files handled when ExerciseMedia or ProgressPhoto is deleted? | Disk bloat risk; affects backup size. | Media lifecycle policy. | worker-attempt-2 | open | TBD |
| TQ-DATA-008 | data-contracts | watchlist | TQ-DATA-007 | What is the data retention and lifecycle policy (auto-delete vs manual)? | No compliance req for single-user, but affects long-term storage. | Retention schedule. | worker-attempt-2 | open | TBD |
| TQ-DATA-009 | data-contracts | watchlist | none | How does backup import handle ID collisions, FK violations, and partial restore failure? | Dry-run stated but recovery path unspecified. | Import collision strategy. | worker-attempt-2 | open | TBD |
| TQ-DATA-010 | data-contracts | dev-blocking | DEC-009 | What are the correct entity names for the 3 stale "WorkoutDay" references? | Blocks clean code generation and schema naming. | Domain model rename sweep. | worker-attempt-2 | open | TBD |
| TQ-DATA-011 | data-contracts | watchlist | DEC-008 | What index strategy is needed to meet p95 query SLOs (daily log <=500ms, exercise history <=700ms, etc.)? | Missing index design risks SLO failure. | Index strategy document. | worker-attempt-2 | open | TBD |

Consolidation note: Storage engine config (connection pooling, SSL, timeouts) is architecture-scope, not data-contracts. Not creating a separate question here.

## Answer Effects
No new answer effects. Same as attempt 1.

## Risks
Same as attempt 1. Added: Stale WorkoutDay references risk developer confusion and inconsistent code if renamed after implementation starts.

## Suggested Decisions
Same as attempt 1. Added: Sweep all source docs for "WorkoutDay" text and replace with "DailyLog" as a single pre-implementation task.

## Traceability Candidates
Same as attempt 1. Added:
- DEC-009 trace to TQ-DATA-010
- DEC-008 trace to TQ-DATA-011