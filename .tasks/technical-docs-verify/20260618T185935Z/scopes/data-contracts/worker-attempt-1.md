# Data-Contracts Worker Attempt 1

## Sources Read
- docs/product-verified/domain-model.md (20 entities, attributes, relationships, invariants)
- docs/product-verified/scope.md (in/out scope, dependencies, assumptions)
- docs/product-verified/product-brief.md (success metrics, quality gates, performance targets)
- docs/product-verified/functional-spec.md (capability map, feature behavior, validations)
- docs/product-verified/appendix/decision-log.md (DEC-006 through DEC-009)
- .tasks/technical-docs-verify/20260618T185935Z/source-delta.md

## Source Delta Reviewed
DEC-007: Single-user MVP with multi-user-ready data model. All entities linked via userId. One default user at bootstrap.

DEC-009: DailyLog replaces WorkoutDay. Cardio is separate entity with required dailyLogId. WorkoutExercise references DailyLog. System auto-creates DailyLog on first activity for a date.

No prior docs/technical-verified/ exists — initial run.

## Product Signals
- Single-user self-hosted MVP with multi-user-ready data model (userId on all entities)
- PostgreSQL primary store, Redis mentioned in stack
- File system volume for media
- Full backup/restore with ZIP (manifest.json, data.json, media/)
- AI export with structured data package
- No entity-level retention or privacy policies stated
- DailyLog is the central date-anchored aggregate root for workouts + cardio
- 20 entities defined with varying completeness

## Technical Facts

### Entity Coverage
All 20 entities listed with id, attributes, and key identifiers.
Relationships documented (1:N, N:1) for major entity groups.
Lifecycle states partially documented (Exercise, AiExport, AiReview, NutritionTemplate, DailyNutritionOverride).
Invariants documented (10 items).

### DEC-007 Compliance Check (userId on all entities)
| Entity | Has userId in Attributes | Has userId in Key ID Col | Multi-User Ready? |
| --- | --- | --- | --- |
| Settings | NOT listed in attributes | id (userId FK) | Partial — FK stated but attribute missing |
| UserProfile | NOT listed in attributes | id (userId FK) | Partial — FK stated but attribute missing |
| DefaultUser | N/A (is the user) | id | N/A — the default user entity itself |
| Exercise | NOT listed in attributes | id (userId FK) | Partial |
| ExerciseMedia | NOT listed in attributes | id (userId FK) | Partial |
| DailyLog | userId listed | id (userId + date unique) | Yes |
| WorkoutExercise | userId listed | id (userId FK) | Yes |
| WorkoutSet | NOT listed | id | **No** — child of WorkoutExercise; needs review |
| CardioEntry | userId listed | id (dailyLogId required) | Yes |
| BodyWeightEntry | NOT listed | id (userId FK) | Partial |
| BodyCheckIn | NOT listed | id (userId FK) | Partial |
| BodyMeasurement | NOT listed | id | **No** — child of BodyCheckIn |
| ProgressPhoto | NOT listed | id | **No** — child of BodyCheckIn |
| NutritionProduct | NOT listed | id (userId FK) | Partial |
| NutritionTemplate | NOT listed | id (userId FK) | Partial |
| NutritionTemplateItem | NOT listed | id | **No** — child of NutritionTemplate |
| DailyNutritionOverride | NOT listed in attributes | id (userId FK) | **CONTRADICTION** — key ID column says (userId FK) but attributes have NO userId field |
| DailyNutritionOverrideItem | NOT listed | id | **No** — child of DailyNutritionOverride |
| WeekFlag | NOT listed | id (userId FK) | Partial |
| AiExport | NOT listed | id (userId FK) | Partial |
| AiReview | NOT listed | id (userId FK) | Partial |

### DEC-009 Compliance Check
- DailyLog exists with userId, date, unique constraint per user-dates — CONSISTENT
- CardioEntry has dailyLogId (required) — CONSISTENT
- WorkoutExercise has dailyLogId — CONSISTENT
- Relationship "WorkoutDay 0:1 BodyWeightEntry" in §Relationships line 109 still says WorkoutDay — **INCONSISTENT** — should reference DailyLog
- BodyWeightEntry invariant references "WorkoutDay" — stale reference
- "Exercise order: WorkoutExercise.order determines display order within WorkoutDay" (invariant 6) — stale reference

### Storage Engine Coverage
- PostgreSQL for relational data — stated in scope.md §Dependencies
- Redis for cache/sessions — stated in scope.md §Dependencies
- Filesystem volume for media — stated in scope.md §Dependencies
- No Redis usage pattern documented (cache keys, session store, job queue, rate limiter?)
- No PostgreSQL-specific types, index strategy, or partition strategy documented

### Migrations Status
- No migration strategy defined: tool (Goose, Atlas, Prisma?), versioning, rollback, seed data
- Default user bootstrap not documented as a migration step
- DEC-007 (adding userId to all entities) is a schema migration, not just a design note

### Seed Data / Fixtures
- Default user creation mentioned but no details: id, displayName, timestamps
- No fixture files or test data strategy documented
- BodyWeightEntry source enum undefined — needed for seed/validation
- MeasurementType enum undefined — needed for seed/validation
- WeekFlag flagType enum undefined — needed for seed/validation

### Retention & Privacy
- No retention policy: auto-delete old data? manual only?
- No privacy/privacy-by-design provisions
- Media file cleanup when ExerciseMedia is deleted? Orphan detection?
- Backup ZIP contains all user data — no privacy concerns stated (single-user mitigates this)

### Import/Export Data Lineage
- Full backup ZIP structure: manifest.json, data.json, media/ — stated
- AI export ZIP structure: manifest.json, data.json, summary.md, CSVs, photos/ — stated
- Schema version in manifest — stated
- Dry-run validation with summary — stated
- No entity-level export/import (full only)
- No partial backup or incremental backup
- Data lineage during restore: FK integrity, orphan handling, ID collision — NOT documented

## Technical Gaps

### T1: userId Attribute Inconsistency
DEC-007 mandates userId on ALL entities, but several entities lack userId in their attribute lists.

### T2: DailyNutritionOverride userId Contradiction
Key identifier column says "(userId FK)" but attribute list has no userId field.

### T3: Stale WorkoutDay References
Three references to old WorkoutDay name remain after DEC-009:
1. Relationships § line 109: "WorkoutDay 0:1 BodyWeightEntry"
2. Invariant 6: "WorkoutExercise.order determines display order within WorkoutDay"
3. Lifecycle states: "WorkoutDay has implicit states"

### T4: Undefined Enums
Three enum types marked "enum undefined":
1. BodyWeightEntry.source
2. BodyMeasurement.measurementType
3. WeekFlag.flagType

### T5: Migration Strategy Missing
No migration tool, versioning, rollback, or seed migration defined.

### T6: Redis Usage Pattern Undefined
Redis is in the stack but no usage pattern (session store, cache, job queue, rate limiting) is documented.

### T7: Media Orphan Management
No policy for cleaning up media files when ExerciseMedia, ProgressPhoto, or other media entities are deleted.

### T8: Data Retention Policy Missing
No retention/deletion policy for old daily logs, exercises, body entries, or AI exports.

### T9: Import Collision Handling
Backup import dry-run stated but ID collision, FK integrity, or partial-failure handling not documented.

### T10: Multi-User Readiness Depth
Scope says "multi-user-ready data model" but child entities (WorkoutSet, BodyMeasurement, ProgressPhoto, NutritionTemplateItem, DailyNutritionOverrideItem) lack userId. Multi-user readiness through parent-chain traversal is a design decision that should be explicit.

## Missing Source Artifacts
- SQL/DDL migrations
- Seed data / fixture files
- Index strategy
- Redis usage specification
- Storage engine configuration (connection pooling, timeouts)
- Media file lifecycle management
- Data retention schedule
- Import collision resolution
- Enum definitions (source, measurementType, flagType)
- Schema version format for manifest.json

## Questions Raised

| ID | Severity | Question | Why It Matters |
| --- | --- | --- | --- |
| TQ-DATA-001 | needs-owner-decision | Which entities require explicit userId FK vs. parent-chain traversal for multi-user readiness? | DEC-007 says "all entities" but 6 child entities omit userId. Design decision needed. |
| TQ-DATA-002 | dev-blocking | What is the userId attribute on entities that declare "id (userId FK)" in key column but omit userId in attributes? | Settings, UserProfile, Exercise, ExerciseMedia, BodyWeightEntry, BodyCheckIn, NutritionProduct, NutritionTemplate, DailyNutritionOverride, WeekFlag, AiExport, AiReview are affected. |
| TQ-DATA-003 | dev-blocking | Does DailyNutritionOverride have userId or not? Key column says yes, attributes say no. | Direct contradiction must be resolved before schema design. |
| TQ-DATA-004 | dev-blocking | What are the allowed values for BodyWeightEntry.source, BodyMeasurement.measurementType, and WeekFlag.flagType? | Blocks entity definition, validation, DB schema, and seed data. |
| TQ-DATA-005 | dev-blocking | Which migration tool and versioning strategy is used? | Blocks all schema changes, seed data, and bootstrap. |
| TQ-DATA-006 | deferred | What Redis usage patterns are expected for MVP? | Blocks cache/session/queue implementation if Redis is used for those. Deferred: not implementation-blocking if Redis is for post-MVP. |
| TQ-DATA-007 | needs-owner-decision | How are orphaned media files handled when ExerciseMedia or ProgressPhoto is deleted? | Risk of disk bloat; affects backup size and cleanup. |
| TQ-DATA-008 | watchlist | What is the data retention policy (auto-delete vs manual only)? | No compliance requirements for single-user MVP, but affects long-term storage. |
| TQ-DATA-009 | watchlist | How does backup import handle ID collisions, FK violations, and partial restore failure? | Dry-run validation stated but recovery path not specified. |

## Answer Effects
No prior technical questions exist — this is the initial run.
Source delta (DEC-007, DEC-009) was fully analyzed. No follow-up blockers discovered that directly contradict the product decisions, but DEC-007's "all entities" language is ambiguous for child entities.

## Risks
- If DEC-007 truly means EVERY entity (including child entities), schema will have redundant userId FKs on deeply nested tables. This adds write overhead but simplifies queries.
- If DEC-009 auto-create of DailyLog is not handled correctly, cardio entries could be orphaned.
- Missing migration strategy risks manual schema drift in early development.

## Suggested Decisions
1. Clarify DEC-007 depth: explicit userId on aggregate-root entities only vs. all entities.
2. Adopt Goose or Atlas for PostgreSQL migrations.
3. Define the three enums from product behavior:
   - BodyWeightEntry.source: "manual", "check-in", "import" (derive from PRD)
   - BodyMeasurement.measurementType: "neck", "shoulders", "forearms", "biceps", "chest", "waist", "abdomen", "hips", "thigh", "calves" (from functional-spec.md §13)
   - WeekFlag.flagType: "poor-sleep", "high-stress", "illness", "injury-pain", "aas-cycle", "calorie-deficit", "calorie-surplus", "maintenance", "missed-workouts", "travel-disrupted-routine" (from functional-spec.md §17-18)
4. Move measurements to a single-row per type (left and right as nullable columns) instead of side enum — simplifies queries.

## Traceability Candidates
- DEC-007 trace to TQ-DATA-001, TQ-DATA-002, TQ-DATA-003
- DEC-009 trace to TQ-DATA-010 (if created as stale references question)
- functional-spec.md §13 trace to measurementType enum
- functional-spec.md §17-18 trace to flagType enum