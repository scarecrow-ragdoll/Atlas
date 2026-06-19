# Wave Map Context

## Selected Backend Wave Boundary
WAVE-03 (Workout Diary): DailyLog CRUD by date, WorkoutExercise with order and working weight snapshot, WorkoutSet with weight/reps/RPE/RIR, CardioEntry inline to DailyLog, per-exercise comments. All operations via GraphQL, PIN-protected by WAVE-01 middleware.

## Prior Backend Wave Fit
- WAVE-01 (Foundation): prerequisite — provides PIN auth middleware, common GraphQL types, gqlgen+sqlc codegen infrastructure, config extension pattern. WAVE-03 depends on all these.
- WAVE-02 (Exercise Library): prerequisite — provides allExercises query for exercise selector, Exercise.workingWeight for snapshot, exercises table FK target. WAVE-03 depends on all these.

## Future Backend Wave Fit
- WAVE-04 (Cardio/Body): CardioEntry shared entity boundary. WAVE-03 creates DailyLog-linked cardio entries. WAVE-04 adds standalone cardio and body tracking. No collision — ownership documented.
- WAVE-05 (Nutrition): No direct dependency. Independent domain.
- WAVE-06 (Charts): WAVE-03 provides data model (sets, weights, volumes) consumed by chart aggregation.
- WAVE-07/08 (AI Export/Review): WAVE-03 provides exercise comments and set data via service layer.
- WAVE-09 (Backup): WAVE-03 tables are JSON-serializable for export compatibility.

## Frontend Pages Context
- PAGE-002 (Workout Diary): primary frontend consumer — depends on all WAVE-03 GraphQL endpoints listed below. Dependency context only; no frontend pages, UI, or UX work in this wave.
- Frontend backend dependencies from PAGE-002: dailyLogByDate query, upsertDailyLog mutation, addWorkoutExercise mutation, addWorkoutSet mutation, addCardioEntry mutation, dailyLogsByDateRange query.
- WAVE-02 allExercises query also consumed by PAGE-002 (not provided by WAVE-03).

## Dependency Order
WAVE-01 (Foundation) -> WAVE-02 (Exercise Library) -> WAVE-03 (Workout Diary) -> WAVE-04/WAVE-05 -> WAVE-06 -> WAVE-07 -> WAVE-08 -> WAVE-09

## Scope Collision Check
| Wave | Collision Risk | Assessment |
| --- | --- | --- |
| WAVE-01 (Foundation) | None | Infrastructure vs domain logic |
| WAVE-02 (Exercise Library) | None | Exercise CRUD vs Workout CRUD |
| WAVE-04 (Cardio/Body) | CardioEntry | WAVE-03 creates DailyLog-linked entries; WAVE-04 handles standalone |
| WAVE-05 (Nutrition) | None | Independent domain |
| WAVE-06 (Charts) | None | Read-only consumer |
| WAVE-07 (AI Export) | None | Read-only consumer |
| WAVE-08 (AI Review) | None | Independent domain |
| WAVE-09 (Backup) | None | Read-only consumer |
