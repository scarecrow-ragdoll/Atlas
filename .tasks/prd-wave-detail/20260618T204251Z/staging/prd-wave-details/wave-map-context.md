# Wave Map Context

## Selected Backend Wave Boundary
WAVE-02 (Exercise Library): Full CRUD for exercises with working weight and media management. GraphQL for CRUD operations, REST for binary media upload/download/delete. All endpoints protected by WAVE-01 PIN auth middleware.

## Prior Backend Wave Fit
WAVE-01 (Foundation): Provides required infrastructure — PIN auth middleware, media REST scaffold, goose migration infrastructure, fitness GraphQL common types, gqlgen+sqlc codegen config, config extension pattern. WAVE-02 depends on all these contracts. No scope collision: WAVE-01 creates the framework; WAVE-02 owns domain logic.

## Future Backend Wave Fit
- WAVE-03 (Workout Diary): WAVE-02 provides allExercises query for exercise selector, exercise by ID for working weight snapshot
- WAVE-04/05: No direct dependency
- WAVE-06 (Charts): WAVE-02 provides exercise metadata (name, muscleGroups) for chart labels and filters; WAVE-03 provides historical working weight data
- WAVE-07/08 (AI Export): WAVE-02 exercise + media data consumed via service layer
- WAVE-09 (Backup): WAVE-02 tables are JSON-serializable; files referenced by file_path are backup-compatible
- No scope collision: all 8 other waves checked, clean separation maintained

## Frontend Pages Context
- PAGE-003 (Exercise Library): dependent on all WAVE-02 backend endpoints (GraphQL CRUD + REST media). No frontend work in this wave.
- PAGE-002 (Workout Diary): depends on allExercises query for exercise selector in workout entry
- PAGE-011 (Settings): PIN auth configuration from WAVE-01, not dependent on WAVE-02

## Dependency Order
WAVE-01 (Foundation) → WAVE-02 (Exercise Library) → WAVE-03 (Workout Diary) → WAVE-04/WAVE-05 → WAVE-06 → WAVE-07 → WAVE-08 → WAVE-09

## Scope Collision Check
No scope collisions with any other wave. WAVE-02 owns exercise + exercise_media domain exclusively.