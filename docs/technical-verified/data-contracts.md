# Data Contracts

## Entities And Identifiers

21 entities defined in domain model. All major entities require userId FK per DEC-007. The following entities have userId FK:

- Settings, UserProfile, DefaultUser, Exercise, ExerciseMedia, DailyLog, WorkoutExercise, CardioEntry, BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto, NutritionProduct, NutritionTemplate, DailyNutritionOverride, WeekFlag, AiExport, AiReview

Child entities (WorkoutSet, NutritionTemplateItem, DailyNutritionOverrideItem) inherit userId through parent chain.

Missing userId from entity descriptions: ExerciseMedia, WorkoutSet, NutritionTemplateItem, DailyNutritionOverrideItem, BodyMeasurement, ProgressPhoto (TQ-DATA-001).

## Persistence And Storage

- PostgreSQL for all structured data
- Redis for PIN session store
- Filesystem volume for media (images, video, export ZIPs, backup ZIPs)

No index strategy defined for p95 performance targets (TQ-DATA-002).

## Migrations

No migration strategy defined (TQ-DATA-005):
- Schema version in backup manifest exists but no migration mechanics
- No rollback strategy
- No seed data for default user bootstrap

## Retention And Privacy

No data retention policy defined (TQ-COMP-001 from product run). Media lifecycle (TQ-DATA-008) and import collision handling (TQ-DATA-010) are undefined.

## Data Questions

| ID | Question | Severity | Status |
| --- | --- | --- | --- |
| TQ-DATA-001 | userId FK missing from 6 child entity attribute descriptions | dev-blocking | **resolved** (TDEC-020) |
| TQ-DATA-002 | No index strategy for p95 query targets | dev-blocking | **resolved** (TDEC-021) |
| TQ-DATA-003 | 3 enum types undefined (source, heartRateZone, cardioType) | dev-blocking | **resolved** (TDEC-022) |
| TQ-DATA-004 | 3 enum types undefined (measurementType, flagType, mediaType) | dev-blocking | **resolved** (TDEC-023) |
| TQ-DATA-005 | No migration strategy or schema versioning mechanics | dev-blocking | **resolved** (TDEC-024) |
| TQ-DATA-006 | Stale "WorkoutDay" references remain after DEC-009 DailyLog rename | needs-owner | **resolved** (TDEC-006) |
| TQ-DATA-007 | BodyWeightEntry.source enum undefined | needs-owner | **resolved** (TDEC-007) |
| TQ-DATA-008 | No media lifecycle policy (cleanup on delete, storage limits) | deferred | deferred |
| TQ-DATA-009 | No fixture/seed data format for tests | watchlist | open |
| TQ-DATA-010 | Import collision handling (existing data + re-import) undefined | dev-blocking | **resolved** (TDEC-025) |