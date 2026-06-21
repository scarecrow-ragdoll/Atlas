# Wave Map Context

## Selected Backend Wave Boundary
WAVE-06: Charts — Progress visualization. Read-only query wave providing time-series chart data for exercises (conditional on WAVE-03), body weight, body measurements, and nutrition macros. No mutations, no new storage, no new tables. Depends on WAVE-01 (PIN auth), WAVE-03 (workout sets), WAVE-04 (body data), WAVE-05 (nutrition data).

## Prior Backend Wave Fit
- WAVE-01 (Foundation): provides PIN auth middleware, Atlas GraphQL endpoint, migration infrastructure
- WAVE-02 (Exercise Library): provides exercises table for exercise selector — used via existing ExerciseService.ListAll
- WAVE-03 (Workout Diary): provides workout_sets and workout_exercises tables for exercise chart queries. WAVE-03 NOT fully implemented yet — exercise chart queries return empty series until WAVE-03 deploys
- WAVE-04 (Cardio and Body Tracking): provides body_weight_entries, body_check_ins, body_measurements tables. Fully implemented. WAVE-06 adds one new sqlc query (measurement range) and one repo method
- WAVE-05 (Nutrition): provides nutrition tables and NutritionMacroService. Fully implemented. WAVE-06 wraps existing macro service for weekly averages

No pattern or contract conflicts with prior detailed waves.

## Future Backend Wave Fit
- WAVE-07 (AI Export): shares underlying workout/body/nutrition data but chart queries are ephemeral — no registration needed
- WAVE-08 (AI Review): no dependency
- WAVE-09 (Backup): chart queries produce no persistent data — no serialization concerns

No scope collision with future waves.

## Frontend Pages Context
- PAGE-008 (Charts): primary frontend consumer
  - Exercise selector dropdown → WAVE-02 allExercises query
  - Period-filtered exercise progress → WAVE-06 exerciseChart query (stub until WAVE-03)
  - Period-filtered body weight → WAVE-06 bodyWeightTrend query
  - Period-filtered measurement trend → WAVE-06 measurementTrend query
  - Period-filtered measurement overlay → WAVE-06 measurementOverlay query
  - Period-filtered nutrition → WAVE-06 nutritionWeeklyAverages query
- PAGE-001 (Dashboard): depends on body weight and macro summary — chart queries provide data for dashboard sections

Dependency context only; no frontend pages, UI, or UX work in this wave.

## Dependency Order
WAVE-01 → WAVE-02 → WAVE-03 (partial) → WAVE-04 → WAVE-05 → WAVE-06

## Scope Collision Check
- Measurement range sqlc query: new file, no collision with WAVE-04 (WAVE-04 has no range query for measurements by type+date)
- Nutrition weekly average: wraps existing NutritionMacroService — no collision with WAVE-05
- No new tables, no new migrations — no migration number collision
- Read-only wave — no mutation/collision with any other wave