# WAVE-02 data-integration-ops Planner Attempt 1

## Sources Read
- docs/technical-verified/data-contracts.md (Exercise, ExerciseMedia entities)
- docs/technical-verified/api-contracts.md (hybrid API pattern)
- docs/technical-verified/operations-observability.md (logging, metrics, config)
- docs/technical-verified/integrations-and-events.md (file size limits, rate limits)
- docs/product-verified/domain-model.md (Exercise, ExerciseMedia)
- docs/product-verified/business-rules.md (RULE-003 working weight)
- apps/api/internal/appconfig/config.go
- apps/api/cmd/server/main.go

## Selected Backend Wave Boundary
WAVE-02 data lifecycle covers exercise and exercise_media PostgreSQL tables, file storage for exercise media, and the query/upload/download operations. No external integrations, no background jobs, no events.

## Neighboring Backend Wave Fit
- WAVE-01: provides the PostgreSQL connection pool, migration runner, Redis session store. WAVE-02 adds exercise tables to existing DB.
- WAVE-03: will read exercise data (list for selector) and snap workingWeight into WorkoutExercise at log time.
- No data scope collision with later waves.

## Frontend Pages Context
PAGE-003 needs: exercise CRUD operations, media upload/download. PAGE-002 needs: exercise list read-only.

## Codebase Evidence
- Existing migration files: apps/api/internal/repository/postgres/migrations/00001_init.sql, 00079_admin_users.sql
- Existing query files: apps/api/internal/repository/postgres/queries/users.sql, admin_users.sql
- Existing media REST scaffold planned for WAVE-01

## Proposed Details

### Data Model

#### exercises table
```sql
CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    muscle_groups TEXT[] NOT NULL DEFAULT '{}',
    description TEXT,
    personal_notes TEXT,
    working_weight REAL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_exercises_active_name ON exercises (is_active, name);
CREATE INDEX idx_exercises_name_trgm ON exercises USING gin (name gin_trgm_ops);
```

#### exercise_media table
```sql
CREATE TABLE exercise_media (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    media_type VARCHAR(32) NOT NULL,
    file_path TEXT NOT NULL,
    original_file_name VARCHAR(512) NOT NULL,
    mime_type VARCHAR(128) NOT NULL,
    size_bytes BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_exercise_media_exercise ON exercise_media (exercise_id);
```

Note: ON DELETE CASCADE on exercise_media.exercise_id ensures media records are removed when exercise is hard-deleted. However, WAVE-02 uses soft delete (isActive=false), so cascade only applies to hard deletion. Media rows persist during soft delete, which is correct for WAVE-03 historical references.

### Migrations
- File: 00080_exercises.sql (creates exercises table, indexes, trigram extension)
- File: 00081_exercise_media.sql (creates exercise_media table with FK to exercises)
- Migration order: 00080, 00081 (sequential, 00080 must run first)

### Data Lifecycle
- **Create exercise**: INSERT with validated fields. Return full exercise row.
- **Read exercise**: SELECT by id, or list with cursor pagination + isActive filter.
- **Update exercise**: UPDATE optional fields (name, muscleGroups, description, personalNotes, workingWeight). isActive cannot be updated through update — use dedicated delete/reactivate mutations.
- **Soft delete (isActive=false)**: UPDATE is_active = false. Row remains for WAVE-03 historical references.
- **Hard delete**: Not implemented in WAVE-02. Deferred to WAVE-03 or future cleanup.
- **Upload media**: Validate file size (max 25MB images, 250MB video per TDEC-008), validate mime type (JPEG/PNG/WEBP for images, MP4/MOV/WEBM for video), store file at configured media path, INSERT exercise_media record.
- **Delete media**: DELETE exercise_media row, DELETE physical file from disk. If physical file deletion fails, log error but do not fail the request (file system may be inconsistent but DB is clean).

### File Storage
- Store exercise media at: <MediaConfig.BasePath>/exercise/<exercise_id>/<uuid>.<ext>
- WAVE-01 MediaConfig provides BasePath for media storage root.
- File organization: grouped by exercise ID for easy backup/cleanup.

### API Operations

#### GraphQL (exercises)
- `exercises(first: Int, after: String, includeInactive: Boolean = false): ExerciseListResult!` — cursor-paginated list
- `exercise(id: UUID!): ExerciseResult!` — single exercise by ID
- `allExercises(includeInactive: Boolean = false): [Exercise!]!` — simple list for WAVE-03 selector
- `createExercise(input: CreateExerciseInput!): ExerciseResult!`
- `updateExercise(id: UUID!, input: UpdateExerciseInput!): ExerciseResult!`
- `deleteExercise(id: UUID!): DeleteExerciseResult!` — sets isActive=false

#### REST (exercise media)
- `POST /api/v1/exercise-media` — multipart: exerciseId, file. Returns ExerciseMedia JSON.
- `DELETE /api/v1/exercise-media/{id}` — removes association + file.
- Reuses WAVE-01 media storage path and file delivery.

### Error Handling
- Exercise not found → NotFoundError (union type)
- Exercise name empty → ValidationError
- Working weight <= 0 → ValidationError
- File too large → ValidationError with field "file"
- Invalid file type → ValidationError with supported types listed
- Media not found → 404 error envelope
- PIN auth failure → AuthError or 401

### Observability
- Log markers: [Exercise][create|update|delete|list][BLOCK_VALIDATE_INPUT], [ExerciseMedia][upload|delete][BLOCK_STORE_FILE]
- Standard zap logger from logger.FromContext
- No custom metrics in WAVE-02
- Do not log: exercise personalNotes content, working weight values in error logs (minimal exposure — low sensitivity but follow pattern from WAVE-01)

### Validation
- Name: required, non-empty, trimmed, max 255 chars
- Muscle groups: optional array of strings, each max 64 chars
- Working weight: optional, if provided must be > 0
- Description: optional text, max 2000 chars
- Personal notes: optional text, max 5000 chars
- File size: enforce TDEC-008 limits (25MB images, 250MB video)
- File type: enforce allowed formats (JPEG/PNG/WEBP, MP4/MOV/WEBM)

## Acceptance Criteria Contributions
(Delegated to product-ac and architecture-codebase planners)

## Exit Criteria Contributions

| EC ID | Description |
| --- | --- |
| EC-W02-004 | WAVE-01 media REST scaffold handles exercise-media association extension |
| EC-W02-006 | Migration creates indexes for active+name search and trigram name search |
| EC-W02-009 | File size and type validation enforced for exercise media uploads |
| EC-W02-010 | ON DELETE CASCADE works correctly for hard-deleted exercises |

## Verification Contributions

| Test ID | Description | Type |
| --- | --- | --- |
| TEST-W02-009 | ExerciseMedia file size validation (reject >25MB images, >250MB video) | unit |
| TEST-W02-010 | ExerciseMedia file type validation (accept JPEG/PNG/WEBP/MP4/MOV/WEBM, reject others) | unit |
| TEST-W02-011 | Exercise name uniqueness not enforced (no constraint) | integration |
| TEST-W02-012 | Soft-deleted exercise not returned in default list | integration |
| TEST-W02-013 | Soft-deleted exercise returned with includeInactive=true | integration |

## Risks And Rollback
- Risk: large media uploads block the request thread. Mitigation: standard Go HTTP request handling with read timeout; TDEC-008 limits prevent unbounded uploads. No async processing needed for WAVE-02.
- Risk: disk space exhaustion from media uploads. Mitigation: not addressed in WAVE-02 (deferred to WAVE-01 media lifecycle or global storage monitoring).
- Rollback: revert migrations (00080, 00081), revert code. Exercise data remains in DB but is inaccessible without WAVE-02 code. WAVE-03 workout data referencing deleted exercises will still find the exercise row via ID (soft delete preserves row).

## Questions Raised

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W02-003 | WAVE-02 | data-integration-ops | needs-owner-decision | TDEC-008 | What is the exact file storage path pattern for exercise media? | Needs to align with WAVE-01 MediaConfig base path | data-integration-ops planner | open |
| DQ-W02-004 | WAVE-02 | data-integration-ops | wave-blocking | EDGE-020 | When ExerciseMedia is deleted, should the physical file be deleted immediately, logged for cleanup, or retained? | TDEC-005 says "remove metadata + physical file" but cleanup failure handling is undefined | data-integration-ops planner | open |

## Traceability Candidates
- exercises table → docs/product-verified/domain-model.md#Exercise
- exercise_media table → docs/product-verified/domain-model.md#ExerciseMedia
- File validation limits → docs/technical-verified/integrations-and-events.md (TDEC-008)
- Working weight rule → docs/product-verified/business-rules.md RULE-003