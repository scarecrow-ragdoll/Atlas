# WAVE-06 Sequencing-Fit Planner Attempt 1

## Sources Read
- docs/prd-waves/waves/wave-06.md
- docs/prd-waves/wave-map.md
- docs/prd-waves/frontend-pages/page-008.md
- docs/prd-wave-details/waves/index.md
- docs/prd-wave-details/waves/wave-04.md
- docs/prd-wave-details/waves/wave-05.md
- docs/development-plan.xml
- docs/knowledge-graph.xml
- apps/api/internal/atlas — codebase inspection

## Selected Backend Wave Boundary
Pure query wave: chart data aggregation for exercise, body, nutrition domains. No mutations, no tables, no storage.

## Prior Backend Wave Fit
- WAVE-01 (Foundation): prerequisite — provides PIN auth middleware, GraphQL endpoint, gqlgen config, sqlc config
- WAVE-02 (Exercise Library): fully implemented — provides exercises table + service. WAVE-06 needs exercise list query (exists: allExercises/listAllExercises).
- WAVE-03 (Workout Diary): NOT implemented — only migration + DailyLogRecord model exist. WAVE-06 EXERCISE CHART QUERIES ARE BLOCKED on WAVE-03 providing:
  - daily_log_exercises table (bridge between daily_log and exercises)
  - workout_sets table (weight, reps, rpe, rir per set)
  - Service methods for set-level queries
  - Working weight per-session snapshot (RULE-017)
- WAVE-04 (Cardio/Body Tracking): fully implemented — provides BodyWeightEntry, BodyCheckIn, BodyMeasurement repos/services. WAVE-06 can query body weight range and measurements.
- WAVE-05 (Nutrition): fully implemented — provides NutritionMacroService. WAVE-06 can reuse for weekly averages.

## Future Backend Wave Fit
- WAVE-07 (AI Export): WAVE-06 provides no data directly but charts share the same underlying data — no collision.
- WAVE-08 (AI Review): no dependency.
- WAVE-09 (Backup): chart queries have no serialization concerns — queries are ephemeral.

## Frontend Pages Context
- PAGE-008 (Charts): primary consumer. Backend must provide:
  - Exercise list for dropdown → WAVE-02 allExercises query (exists)
  - Period-filtered exercise progress data → blocked on WAVE-03
  - Period-filtered body weight data → WAVE-06 bodyWeightTrend query
  - Period-filtered measurement data → WAVE-06 measurementTrend/measurementOverlay queries
  - Period-filtered nutrition weekly averages → WAVE-06 nutritionWeeklyAverages query
- Dependency context only — no frontend planning.

## Dependency Order
WAVE-01 → WAVE-02 → WAVE-03 → WAVE-04 → WAVE-05 → WAVE-06
(Exercise chart queries specifically need WAVE-03. Body/nutrition charts can be done in parallel with WAVE-03.)

## Scope Collision Check
- WAVE-04 explicitly excluded chart visualization ("Chart visualization (WAVE-06)" in WAVE-04 scope excluded). ✓
- WAVE-05 explicitly excluded nutrition charts ("Nutrition charts (WAVE-06)" in WAVE-05 scope excluded). ✓
- WAVE-03 scope includes "working weight snapshot stored per exercise in workout day" — this is the data WAVE-06 needs. No collision.
- NutritionMacroService already computes per-day macros — WAVE-06 weekly average is additive wrapper, not duplicate.

## Independent Value
- Body weight chart queries: independently valuable even without WAVE-03
- Measurement chart queries: independently valuable
- Nutrition weekly averages: independently valuable
- Exercise chart queries: dependent on WAVE-03

## Deferrals
- Exercise chart queries: implement exercise chart service stub returning empty series with clear documentation. Full implementation deferred to WAVE-03 completion.

## Questions Raised
- DQ-W06-002: Confirm WAVE-06 does not create its own workout_sets table — exercise chart queries are deferred to WAVE-03+.
- DQ-W06-009: If WAVE-03 is completed before WAVE-06 is implemented, should WAVE-06 be re-scoped to include exercise chart queries? (Proposed: yes — add flag in wave doc for "exercise chart queries conditionally implementable")

## Traceability Candidates
- docs/prd-waves/wave-map.md — wave dependency map
- docs/prd-waves/waves/wave-06.md — source wave
- docs/prd-wave-details/waves/wave-04.md — WAVE-04 scope excluded charts
- docs/prd-wave-details/waves/wave-05.md — WAVE-05 scope excluded nutrition charts
- docs/development-plan.xml — module sequencing
- docs/knowledge-graph.xml — module boundaries