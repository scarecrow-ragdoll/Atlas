# Wave Map Context
## Selected Backend Wave Boundary
WAVE-04 (Cardio and Body Tracking): Full CRUD for cardio entries (attached to DailyLog), body weight entries (standalone by date), weekly body check-ins (with nested body measurements and progress photos), and week flags for AI export context. GraphQL for all CRUD operations, REST for binary progress photo upload/download/delete. All endpoints protected by WAVE-01 PIN auth middleware.

WAVE-05 (Nutrition): Full CRUD for nutrition products (with soft-delete), weekly nutrition templates (upsert by week, nested items), daily nutrition overrides (unique per date, nested items with add/subtract/replace operations), and server-side KJBJU macro calculation. All via GraphQL. All endpoints protected by WAVE-01 PIN auth middleware.
## Prior Backend Wave Fit
- WAVE-01 (Foundation): Provides required infrastructure — PIN auth middleware, Atlas GraphQL endpoint, goose migration infrastructure, atlas-gqlgen codegen config, sqlc config, config extension pattern. WAVE-04 depends on PIN auth middleware, media scaffold, and fitness common types. WAVE-05 depends on PIN auth middleware, Atlas GraphQL endpoint, and migration infrastructure.
- WAVE-02 (Exercise Library): No direct dependency to WAVE-04 or WAVE-05. Media handler pattern served as template for WAVE-04 ProgressPhoto handler.
- WAVE-03 (Workout Diary): Partial dependency for WAVE-04 — CardioEntry.dailyLogId FK references daily_log table. No dependency for WAVE-05 — nutrition tables are independent from daily_log.
## Future Backend Wave Fit
- WAVE-05 (Nutrition): No direct dependency with WAVE-04 — can fully parallelize per wave-map
- WAVE-06 (Charts): WAVE-04 provides body weight data (latestBodyWeight, date range), body measurements, and check-in data for chart queries. WAVE-05 provides nutrition macro data for weekly KJBJU average chart queries
- WAVE-07/08 (AI Export/Review): WAVE-04 provides cardio, weight, check-in, measurement, photo, and week flag data via service layer. WAVE-05 provides nutrition template and override data via service layer
- WAVE-09 (Backup): WAVE-04 tables are designed for JSON-serializable export compatibility; photo files referenced by file_path are backup-compatible. WAVE-05 tables are JSON-serializable for export compatibility
- No scope collision: all 8 other waves checked, clean separation maintained
## Frontend Pages Context
- PAGE-004 (Cardio): dependent on all CardioEntry backend endpoints (GraphQL CRUD). Date-based listing for cardio log. No frontend work in this wave.
- PAGE-005 (Body Measurements): dependent on BodyCheckIn, BodyWeightEntry, BodyMeasurement GraphQL queries/mutations. Nested measurements and photos via check-in query. No frontend work in this wave.
- PAGE-006 (Progress Photos): dependent on progressPhotos GraphQL query and REST upload/download/delete endpoints. Photos grouped by check-in. No frontend work in this wave.
- PAGE-007 (Nutrition): dependent on NutritionProduct, NutritionTemplate, DailyNutritionOverride, and nutritionMacros GraphQL queries/mutations. Products CRUD, template editor, daily override editor, and macro summary. No frontend work in this wave.
- PAGE-001 (Dashboard): dependent on latestBodyWeight GraphQL query for weight summary card. No frontend work in this wave.
## Dependency Order
WAVE-01 (Foundation) → WAVE-02 (Exercise Library) → WAVE-03 (Workout Diary) → WAVE-04/WAVE-05 (parallelizable) → WAVE-06 → WAVE-07 → WAVE-08 → WAVE-09
## Scope Collision Check
No scope collisions with any other wave. WAVE-04 owns CardioEntry, BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto, and WeekFlag domains exclusively. WAVE-05 owns NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, and DailyNutritionOverrideItem domains exclusively. The DailyLog dependency with WAVE-03 is a shared table (created by whichever wave deploys first) with no conflicting columns or constraints.