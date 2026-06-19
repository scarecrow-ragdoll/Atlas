# WAVE-02 sequencing-fit Planner Attempt 2

## Cycle 1 Reviewer Feedback Addressed

### 1. WAVE-01 Implementation Dependency
WAVE-02 **cannot start implementation** until WAVE-01 is implemented and provides:
- PIN auth middleware (requirePinAuth function, middleware for chi)
- Media REST scaffold (POST/GET /api/v1/media)
- Migration infrastructure (goose runner, migration file pattern)
- Fitness GraphQL common types (ValidationError, AuthError, NotFoundError)
- Codegen config for fitness domain (gqlgen + sqlc)
- Config extension pattern (MediaConfig, SessionConfig)

If WAVE-01 is delayed, WAVE-02 is blocked. This is an explicit sequencing dependency.

### 2. Exact WAVE-01 Media Contract
WAVE-01's media scaffold must provide:
- **Storage path**: Configurable base path via MediaConfig.BasePath
- **Upload**: POST /api/v1/media/upload accepts multipart file, returns `{ "id": "uuid" }`. Stores file at `<BasePath>/<uuid>.<ext>`.
- **Download**: GET /api/v1/media/{id} returns file with correct Content-Type from stored metadata.

WAVE-02 builds on top: stores exercise association metadata in exercise_media table, uses same file storage for physical files, and provides dedicated exercise-media endpoints.

### 3. allExercises Interface Alignment
Decision: `allExercises` is a **GraphQL query only** (consistent with hybrid pattern where CRUD goes through GraphQL). The PAGE-002/003 frontends use GraphQL to fetch the exercise list for the selector. No REST endpoint for exercise listing.

Interface:
```graphql
allExercises(includeInactive: Boolean = false): [Exercise!]!
```

Returns list ordered by name ASC (alphabetical for selector UX).

### 4. WAVE-06 Data Flow Corrected
Correction from attempt 1: WAVE-06 (Charts) derives historical working weight data from **WAVE-03's WorkoutExercise.workingWeightSnapshot**, not from WAVE-02's current workingWeight.

- WAVE-02 provides: exercise metadata (name, muscle groups, isActive status) for chart labels and filters
- WAVE-03 provides: historical working weight per exercise per session (workingWeightSnapshot)
- WAVE-06 combines: exercise identity from WAVE-02 + historical data from WAVE-03

Updated dependency flow:
```
WAVE-02 (Exercise name, isActive, muscleGroups) ──→ WAVE-06 (Chart labels, filters)
WAVE-03 (workingWeightSnapshot per session) ──→ WAVE-06 (Trend data)
```

### 5. WAVE-09 Backup Compatibility
WAVE-02 data model is designed for backup compatibility:
- Exercises table: all columns are JSON-serializable (UUID, string, float, boolean, timestamptz, text array)
- ExerciseMedia table: all columns are JSON-serializable
- Files referenced by file_path are included in backup media/ directory

No structural changes needed for WAVE-09 compatibility. WAVE-02 design is export-ready.

### 6. Migration Numbering Coordination
WAVE-02 proposes 00080 and 00081, but numbers must be confirmed after WAVE-01 implementation. If WAVE-01 uses 00080 or higher, WAVE-02 numbers are adjusted to the next available numbers (00082, 00083, etc.).

### 7. DQ-W02-008 Moved to Watchlist
`allExercises` filtering beyond isActive is not needed for MVP. WAVE-03 only needs the list of active exercises. Moving DQ-W02-008 to watchlist (not a blocker).

## Updated Dependency Order
```
WAVE-01 (Foundation) → WAVE-02 (Exercise Library) → WAVE-03 (Workout Diary)
                                                      ↓
WAVE-02 (metadata) ──────────────────────────────→ WAVE-06 (Charts)
WAVE-03 (historical weight data) ─────────────────→ WAVE-06 (Charts)
WAVE-02 (exercise + media data) ─────────────────→ WAVE-07/08 (AI Export)
WAVE-02 (table data) ─────────────────────────────→ WAVE-09 (Backup)
```

## Contract Summary (WAVE-02 Interface Exposed to Other Waves)

| Consumer Wave | Interface | Description |
| --- | --- | --- |
| WAVE-03 | `allExercises(includeInactive: Boolean): [Exercise!]!` | Exercise selector list |
| WAVE-03 | `exercise(id: UUID!): ExerciseResult!` | Exercise detail + working weight for snapshot |
| WAVE-03 | Exercise accessible by ID even if isActive=false | Historical workout data preservation |
| WAVE-06 | Exercise metadata (name, muscleGroups) | Chart labels and filters |
| WAVE-07/08 | Exercise + media data via service layer | AI export data assembly |
| WAVE-09 | Exercise + exercise_media tables | Backup export |

## No Scope Collision Confirmed
All 8 other waves checked. No collision with WAVE-02's exercise CRUD scope. Clean separation maintained.