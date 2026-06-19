# WAVE-03 product-ac Planner Attempt 1

## Sources Read
- docs/prd-waves/waves/wave-03.md
- docs/prd-waves/frontend-pages/page-002.md
- docs/product-verified/domain-model.md
- docs/product-verified/functional-spec.md
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/user-flows.md
- docs/product-verified/business-rules.md
- docs/product-verified/edge-cases.md
- docs/prd-wave-details/waves/wave-02.md
- docs/prd-wave-details/waves/wave-01.md

## Selected Backend Wave Boundary
WAVE-03 (Workout Diary) owns all backend behavior for date-based workout tracking. This includes DailyLog CRUD by date, WorkoutExercise management with ordering, WorkoutSet management with weight/reps/RPE/RIR, working weight snapshots, per-exercise comments, and cardio inline or linked. Excluded: Exercise CRUD (WAVE-02), charts (WAVE-06), AI export (WAVE-07).

## Neighboring Backend Wave Fit
- WAVE-01 (Foundation): PIN auth middleware guards all WAVE-03 endpoints. Provides common GraphQL types (ValidationError, AuthError, NotFoundError), sqlc/gqlgen codegen infrastructure, and config extension pattern.
- WAVE-02 (Exercise Library): allExercises query provides exercise selector data. Exercise.workingWeight provides snapshot source. Exercise references preserved in workout history (soft delete).
- WAVE-04 (Cardio/Body): CardioEntry entity lives in WAVE-03 scope (inline or linked to DailyLog). WAVE-04 handles standalone cardio entries and body tracking. No collision: WAVE-03 creates cardio linked to daily log; WAVE-04 creates cardio outside workout context plus body tracking.
- WAVE-06 (Charts): consumes workout data (sets, weights, volumes) for chart aggregation. WAVE-03 provides the data model.
- WAVE-07 (AI Export): consumes workout data via service layer.

## Frontend Pages Context
- PAGE-002 (Workout Diary): primary consumer. Backend dependency context only — no frontend work in this wave. Needs:
  - DailyLog by date (GET/query)
  - DailyLog upsert (create or update by date)
  - WorkoutExercise CRUD within a DailyLog
  - WorkoutSet CRUD within a WorkoutExercise
  - CardioEntry CRUD within a DailyLog
  - allExercises query (from WAVE-02) for exercise selector
  - Working weight auto-population on exercise add

## Codebase Evidence
- Domain model confirms: DailyLog (replaces WorkoutDay), WorkoutExercise, WorkoutSet, CardioEntry entities
- AC-005 through AC-011, AC-035 through AC-042 from product-verified acceptance-criteria.md map to WAVE-03
- EDGE-001 (0 weight/0 reps set), EDGE-004 (backdating), EDGE-005 (empty day), EDGE-016 (concurrent tabs) are relevant edge cases
- RULE-004 (set weight and reps required; RPE/RIR optional), RULE-016 (day created on first save), RULE-017 (working weight snapshot timing)

## Proposed Details

### Outcome After Implementation
- OUT-W03-001: Workouts can be logged by date via DailyLog upsert
- OUT-W03-002: Exercises added with snapshot working weight from Exercise library
- OUT-W03-003: Sets with weight, reps, optional RPE/RIR per exercise
- OUT-W03-004: Cardio can be added per day (inline to DailyLog)
- OUT-W03-005: Per-exercise comments stored and retrievable
- OUT-W03-006: Calendar date switching supported via DailyLog-by-date query
- OUT-W03-007: Data available for charts and AI via service layer

### Scope Included
- CAP-W03-001: DailyLog CRUD by date (upsert: unique per user per date)
- CAP-W03-002: WorkoutExercise with order field
- CAP-W03-003: WorkoutSet with weight, reps, optional RPE, optional RIR
- CAP-W03-004: Working weight snapshot from Exercise (read at add time)
- CAP-W03-005: CardioEntry inline to DailyLog (type, duration, pulse, zone)
- CAP-W03-006: Comments per exercise (notes field on WorkoutExercise)
- CAP-W03-007: Calendar date switching (query DailyLog by date)

### Scope Excluded
- Exercise CRUD (WAVE-02)
- Charts visualization (WAVE-06)
- AI export (WAVE-07)
- Standalone CardioEntry (WAVE-04 — outside workout context)
- Body weight entries (WAVE-04)
- Workout templates (future scope per PAGE-002 deferral)

### Design Decisions
- DailyLog uses UPSERT pattern: GET by date returns existing or empty; first save with exercises creates record
- Working weight snapshot: captured at WorkoutExercise creation time, not at DailyLog save time (RULE-017)
- Set order: WorkoutSet.setNumber is sequential, starts at 1
- Exercise order: WorkoutExercise.order is sequential within DailyLog
- Cardio: required dailyLogId FK. No standalone cardio in WAVE-03.
- Comments: WorkoutExercise.notes field. Not a separate comment entity.
- Backdating: no lower bound on date (DQ-W03-004 resolved)
- Empty day: DailyLog not created until first content (exercise or cardio) saved (DQ-W03-005 resolved)
- Cascade delete: DailyLog -> WorkoutExercise -> WorkoutSet; DailyLog -> CardioEntry (DQ-W03-006 resolved)

## Acceptance Criteria Contributions

| AC ID | Description |
| --- | --- |
| AC-W03-001 | DailyLog can be queried by date (YYYY-MM-DD). Returns existing log or empty result. |
| AC-W03-002 | DailyLog is created/updated via upsert mutation. Returns the full daily log with exercises, sets, and cardio. |
| AC-W03-003 | DailyLog upsert creates a new record when no record exists for the date and exercises/cardio are provided. |
| AC-W03-004 | DailyLog upsert updates an existing record when one already exists for the date. |
| AC-W03-005 | WorkoutExercise can be added to a DailyLog with exerciseId, order, and workingWeightSnapshot. |
| AC-W03-006 | WorkoutExercise workingWeightSnapshot is auto-populated from Exercise.workingWeight at add time. |
| AC-W03-007 | Multiple WorkoutExercises can be added to a single DailyLog, each with unique order values. |
| AC-W03-008 | WorkoutSet can be added to a WorkoutExercise with weight, reps, setNumber. |
| AC-W03-009 | WorkoutSet can optionally include RPE (1-10 scale). |
| AC-W03-010 | WorkoutSet can optionally include RIR (0-5 scale). |
| AC-W03-011 | WorkoutSet can optionally include notes. |
| AC-W03-012 | Multiple WorkoutSets can be added to a single WorkoutExercise with sequential setNumber. |
| AC-W03-013 | WorkoutExercise.notes (comment) can be stored and retrieved in queries. |
| AC-W03-014 | WorkoutExercise order can be updated (reordering). |
| AC-W03-015 | CardioEntry can be created within a DailyLog with cardioType, durationMinutes. |
| AC-W03-016 | CardioEntry optionally includes avgPulse and heartRateZone. |
| AC-W03-017 | WorkoutExercise can be removed from a DailyLog (cascading delete of its WorkoutSets). |
| AC-W03-018 | WorkoutSet can be updated (weight, reps, RPE, RIR, notes). |
| AC-W03-019 | WorkoutSet can be removed from a WorkoutExercise. |
| AC-W03-020 | CardioEntry can be removed from a DailyLog. |
| AC-W03-021 | DailyLog can be deleted entirely (cascades to WorkoutExercises, WorkoutSets, CardioEntries). |
| AC-W03-022 | Set weight must be >= 0 (0 allowed for warmup/bodyweight). Mutation returns ValidationError when negative. |
| AC-W03-023 | Set reps must be > 0 (positive integer). Mutation returns ValidationError when zero or negative. |
| AC-W03-024 | WorkoutExercise.exerciseId must reference an existing active Exercise. Mutation returns NotFoundError when invalid. |
| AC-W03-025 | DailyLog query by date returns empty result (not error) when no record exists for that date. |
| AC-W03-026 | DailyLog returns all WorkoutExercises ordered by WorkoutExercise.order ASC. |
| AC-W03-027 | WorkoutExercise returns all WorkoutSets ordered by WorkoutSet.setNumber ASC. |
| AC-W03-028 | DailyLog date field is validated as YYYY-MM-DD format. Mutation returns ValidationError for invalid dates. |
| AC-W03-029 | DailyLog, WorkoutExercise, WorkoutSet, CardioEntry mutations return AuthError when PIN session header is missing or invalid. |
| AC-W03-030 | DailyLog query by date returns AuthError when PIN session header is missing or invalid. |

## Exit Criteria Contributions
See planner-testing-exit-attempt-1.md — contributed AC coverage expectations.

## Verification Contributions
See planner-testing-exit-attempt-1.md — contributed test scenarios for AC validation.

## Risks And Rollback
- Risk: Concurrent tab editing could cause last-write-wins data loss (Q-WORKOUT-001). Mitigation: acknowledge as MVP risk, defer optimistic locking.
- Risk: Large workout days with many exercises/sets could cause large GraphQL payloads. Mitigation: pagination not needed for single-day scope; set count per exercise practically bounded.
- Rollback: All WAVE-03 database migrations have corresponding down migrations. New GraphQL operations are additive; old clients unaffected.

## Questions Raised
- DQ-W03-004: Backdating lower bound? Resolved: no lower bound.
- DQ-W03-005: Empty workout day creation? Resolved: created on first content save.
- DQ-W03-007: Working weight snapshot timing? Resolved: captured at exercise-add time.

## Traceability Candidates
- docs/product/domain-model.md: DailyLog, WorkoutExercise, WorkoutSet, CardioEntry
- docs/product-verified/acceptance-criteria.md: AC-005 through AC-011, AC-035 through AC-042
- docs/product-verified/user-flows.md: Enter Workout For Today, Enter Workout Backdated
- docs/product-verified/edge-cases.md: EDGE-001, EDGE-004, EDGE-005, EDGE-016
- docs/product-verified/business-rules.md: RULE-004, RULE-016, RULE-017
