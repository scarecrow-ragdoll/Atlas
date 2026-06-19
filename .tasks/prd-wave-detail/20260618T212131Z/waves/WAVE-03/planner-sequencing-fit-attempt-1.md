# WAVE-03 sequencing-fit Planner Attempt 1

## Sources Read
- docs/prd-waves/wave-map.md
- docs/prd-waves/waves/wave-01.md
- docs/prd-waves/waves/wave-02.md
- docs/prd-waves/waves/wave-03.md
- docs/prd-waves/waves/wave-04.md
- docs/prd-waves/waves/wave-05.md
- docs/prd-wave-details/waves/wave-01.md
- docs/prd-wave-details/waves/wave-02.md
- docs/prd-wave-details/wave-map-context.md
- docs/prd-waves/frontend-pages/page-002.md

## Selected Backend Wave Boundary
WAVE-03 (Workout Diary) owns DailyLog CRUD by date, WorkoutExercise with order and working weight snapshot, WorkoutSet with weight/reps/RPE/RIR, CardioEntry inline to DailyLog, and per-exercise comments. This is a pure backend wave providing GraphQL operations consumed by PAGE-002.

## Neighboring Backend Wave Fit

### Prior: WAVE-01 (Foundation)
- Provides PIN auth middleware protecting all WAVE-03 endpoints
- Provides common GraphQL types: ValidationError, AuthError, NotFoundError
- Provides gqlgen/sqlc codegen infrastructure
- Provides config extension pattern and DB connection
- BLOCKING DEPENDENCY: WAVE-03 cannot be implemented or tested until WAVE-01 provides PIN middleware and common types

### Prior: WAVE-02 (Exercise Library)
- Provides allExercises query for exercise selector
- Provides Exercise.workingWeight field for snapshot
- Provides exercises table FK target for workout_exercises.exercise_id
- BLOCKING DEPENDENCY: WAVE-03 needs exercises table and allExercises query to function
- WAVE-02 soft delete (isActive) is compatible: exercises referenced in workout history remain via FK NO ACTION

### Future: WAVE-04 (Cardio and Body Tracking)
- CardioEntry entity is created by WAVE-03 (DailyLog-linked)
- WAVE-04 adds: standalone CardioEntry (outside workout context), BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto
- NO COLLISION: WAVE-03 only creates CardioEntry with required dailyLogId. WAVE-04 may add optional dailyLogId on its CardioEntry or create separate table. Current domain model shows single CardioEntry entity with required dailyLogId.

### Future: WAVE-06 (Charts)
- Consumes workout data (sets, weights, volumes) via service layer queries
- WAVE-03 provides the data model that WAVE-06 queries
- NO COLLISION: WAVE-06 reads only; WAVE-03 writes

### Future: WAVE-07/08 (AI Export/Review)
- Consumes workout data via service layer
- WAVE-03 provides exercise comments and set data for AI context
- NO COLLISION: read-only from WAVE-03 perspective

### Future: WAVE-09 (Backup)
- WAVE-03 tables are JSON-serializable for backup export
- FK relationships preserved in data.json export
- NO COLLISION: backup reads all tables, WAVE-03 writes them

## Frontend Pages Context

### PAGE-002 (Workout Diary)
- Primary (only) frontend consumer of WAVE-03 backend
- Backend dependencies:
  - GET/DailyLog by date: dailyLogByDate(date) query
  - POST/Create DailyLog: upsertDailyLog mutation
  - POST/Create WorkoutExercise: addWorkoutExercise mutation
  - POST/Create WorkoutSet: addWorkoutSet mutation
  - GET/Exercise list from WAVE-02: allExercises query
  - Working weight auto-population on exercise add
  - Cardio entry: addCardioEntry mutation
  - Calendar navigation: dailyLogByDate + dailyLogsByDateRange queries
- No frontend pages, UI, UX, or frontend tests in this wave
- Frontend dependency is purely backend contract — PAGE-002 consumes GraphQL operations exposed by this wave

## Dependency Order
WAVE-01 (Foundation) -> WAVE-02 (Exercise Library) -> WAVE-03 (Workout Diary) -> WAVE-04/WAVE-05 -> WAVE-06 -> WAVE-07/WAVE-08 -> WAVE-09

## Scope Collision Check

| Wave | Collision Risk | Assessment |
| --- | --- | --- |
| WAVE-01 (Foundation) | None | WAVE-01 provides infrastructure; WAVE-03 provides domain logic |
| WAVE-02 (Exercise Library) | None | Exercise CRUD vs Workout CRUD — clean separation |
| WAVE-04 (Cardio/Body) | CardioEntry | WAVE-03 creates DailyLog-linked CardioEntry. WAVE-04 handles standalone cardio and body tracking. Per domain model, CardioEntry has required dailyLogId — WAVE-04 may add standalone variant. No collision in current plan. |
| WAVE-05 (Nutrition) | None | Nutrition is independent domain |
| WAVE-06 (Charts) | None | Read-only consumer |
| WAVE-07 (AI Export) | None | Read-only consumer |
| WAVE-08 (AI Review) | None | Independent domain |
| WAVE-09 (Backup) | None | Read-only consumer |

## Independent Value
WAVE-03 is independently valuable: once deployed alongside WAVE-01/WAVE-02, users can log their daily workouts. This is the core value proposition of the fitness application. It has no dependency on WAVE-04+.

## Deferrals
- Workout templates: explicitly deferred (PAGE-002 deferral). Not in this wave.
- Concurrent edit handling: Q-WORKOUT-001, deferred to post-MVP or owner decision.
- Bulk exercise reorder (drag and drop): not explicitly deferred but not in current scope. Per-exercise order update is sufficient for MVP.
- Pagination for DailyLog: not needed. Single day per query. Date range query returns array (bounded by typical date range).

## Acceptance Criteria Contributions
All ACs validated for dependency fit. No AC depends on unimplemented WAVE-04+ features.

## Exit Criteria Contributions
- EC-W03-015: WAVE-01 test regression pass
- EC-W03-016: WAVE-02 test regression pass

## Verification Contributions
- TEST-W03-027: WAVE-01 admin auth regression
- TEST-W03-028: WAVE-02 exercise regression

## Risks And Rollback
- WAVE-01 delay: if WAVE-01 is not implemented, WAVE-03 cannot start. This is an explicit blocking dependency.
- WAVE-02 delay: same — allExercises and exercises table needed.
- Rollback: WAVE-03 is additive. No changes to existing WAVE-01/WAVE-02 code.

## Questions Raised
- DQ-W03-003: allExercises contract confirmed to include workingWeight.

## Traceability Candidates
- docs/prd-waves/wave-map.md: dependency order
- docs/prd-wave-details/wave-map-context.md: prior/future wave fit
- docs/prd-waves/frontend-pages/page-002.md: backend dependency contract
- docs/prd-waves/waves/wave-04.md: CardioEntry boundary
