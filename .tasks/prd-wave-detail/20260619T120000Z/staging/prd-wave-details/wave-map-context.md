# Wave Map Context
## Selected Backend Wave Boundary
WAVE-04 (Cardio and Body Tracking): Full CRUD for cardio entries (attached to DailyLog), body weight entries (standalone by date), weekly body check-ins (with nested body measurements and progress photos), and week flags for AI export context. GraphQL for all CRUD operations, REST for binary progress photo upload/download/delete. All endpoints protected by WAVE-01 PIN auth middleware.
## Prior Backend Wave Fit
- WAVE-01 (Foundation): Provides required infrastructure — PIN auth middleware, media REST scaffold pattern, goose migration infrastructure, fitness GraphQL common types, gqlgen+sqlc codegen config, config extension pattern. WAVE-04 depends on PIN auth middleware and media scaffold contracts.
- WAVE-02 (Exercise Library): No direct dependency. Media handler pattern (exercise_media.go) serves as template for WAVE-04 ProgressPhoto handler.
- WAVE-03 (Workout Diary): Partial dependency — CardioEntry.dailyLogId FK references daily_log table created in WAVE-03. WAVE-04 includes daily_log migration as prerequisite or defers cardio creation until WAVE-03 is deployed (see DQ-W04-001).
## Future Backend Wave Fit
- WAVE-05 (Nutrition): No direct dependency — can fully parallelize per wave-map
- WAVE-06 (Charts): WAVE-04 provides body weight data (latestBodyWeight, date range), body measurements, and check-in data for chart queries
- WAVE-07/08 (AI Export/Review): WAVE-04 provides cardio, weight, check-in, measurement, photo, and week flag data via service layer. Photos excluded by default per RULE-025
- WAVE-09 (Backup): WAVE-04 tables are designed for JSON-serializable export compatibility; photo files referenced by file_path are backup-compatible
- No scope collision: all 8 other waves checked, clean separation maintained
## Frontend Pages Context
- PAGE-004 (Cardio): dependent on all CardioEntry backend endpoints (GraphQL CRUD). Date-based listing for cardio log. No frontend work in this wave.
- PAGE-005 (Body Measurements): dependent on BodyCheckIn, BodyWeightEntry, BodyMeasurement GraphQL queries/mutations. Nested measurements and photos via check-in query. No frontend work in this wave.
- PAGE-006 (Progress Photos): dependent on progressPhotos GraphQL query and REST upload/download/delete endpoints. Photos grouped by check-in. No frontend work in this wave.
- PAGE-001 (Dashboard): dependent on latestBodyWeight GraphQL query for weight summary card. No frontend work in this wave.
## Dependency Order
WAVE-01 (Foundation) → WAVE-02 (Exercise Library) → WAVE-03 (Workout Diary) → WAVE-04/WAVE-05 (parallelizable) → WAVE-06 → WAVE-07 → WAVE-08 → WAVE-09
## Scope Collision Check
No scope collisions with any other wave. WAVE-04 owns CardioEntry, BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto, and WeekFlag domains exclusively. The DailyLog dependency with WAVE-03 is a shared table (created by whichever wave deploys first) with no conflicting columns or constraints.