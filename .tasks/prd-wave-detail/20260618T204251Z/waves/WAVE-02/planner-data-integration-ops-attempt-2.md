# WAVE-02 data-integration-ops Planner Attempt 2

## Cycle 1 Reviewer Feedback Addressed

### 1. pg_trgm Extension → Simplified Index
Decision: Use simple B-tree indexes instead of pg_trgm to avoid extension dependency. Full-text exercise name search can be added later when needed (or use SQL `ILIKE` with `pg_trgm` in a future optimization wave).

Updated indexes:
```sql
CREATE INDEX idx_exercises_is_active ON exercises (is_active);
CREATE INDEX idx_exercises_name ON exercises (name);
```

### 2. ON DELETE CASCADE → ON DELETE NO ACTION
Changed FK constraint to `ON DELETE NO ACTION` to prevent accidental data loss from cascade deletion:
```sql
exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE NO ACTION,
```

This ensures exercise_media records are protected from accidental hard deletion of exercises. WAVE-02 soft delete (isActive=false) preserves both exercises and their media.

### 3. GET Endpoint for Exercise Media Download
Adding dedicated `GET /api/v1/exercise-media/{id}` endpoint. This is necessary because WAVE-01's `GET /api/v1/media/{id}` stores media by a flat UUID without exercise association. WAVE-02 needs to:
1. Look up exercise_media by ID
2. Verify the exercise association
3. Serve the file from the configured path

This is a thin wrapper around WAVE-01's storage — uses the same file serving mechanism but adds exercise association enforcement.

### 4. Working Weight Data Type
Decision: Use `REAL` for working_weight. Rationale: working weight is a human-readable value (kg/lbs), typically with 0.5 or 0.25 increments. REAL provides sufficient precision. No currency/precision-critical calculations depend on this field.

### 5. REST Error Format → Consistent with TDEC-027
All REST error responses use the standard envelope:
```json
{ "error": { "code": "FILE_TOO_LARGE", "message": "File exceeds maximum size of 25MB for images", "field": "file" } }
```

Error codes for exercise-media REST:
- `FILE_TOO_LARGE` — file exceeds size limit
- `INVALID_FILE_TYPE` — disallowed MIME type
- `NOT_FOUND` — exercise_media ID not found
- `INTERNAL_ERROR` — file storage/disk failure
- `UNAUTHORIZED` — missing/invalid PIN session

### 6. Updated Migration: 00080_exercises.sql
```sql
-- +goose Up
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
CREATE INDEX idx_exercises_is_active ON exercises (is_active);
CREATE INDEX idx_exercises_name ON exercises (name);
-- +goose Down
DROP TABLE IF EXISTS exercises;
```

### 7. Updated Migration: 00081_exercise_media.sql
```sql
-- +goose Up
CREATE TABLE exercise_media (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE NO ACTION,
    media_type VARCHAR(32) NOT NULL,
    file_path TEXT NOT NULL,
    original_file_name VARCHAR(512) NOT NULL,
    mime_type VARCHAR(128) NOT NULL,
    size_bytes BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_exercise_media_exercise ON exercise_media (exercise_id);
-- +goose Down
DROP TABLE IF EXISTS exercise_media;
```

### 8. File Storage Path Pattern
Exercise media stored at: `<WAVE-01 BasePath>/exercise/<exercise_id>/<uuid>.<ext>`
- Grouped by exercise for easy backup and cleanup
- UUID generated server-side, not derived from file name
- Extension derived from detected MIME type server-side

### 9. Physical File Deletion on Media Delete
When `DELETE /api/v1/exercise-media/{id}` is called:
1. DELETE exercise_media row
2. Attempt physical file deletion from disk
3. If file deletion fails: log error with [ExerciseMedia][delete][BLOCK_CLEANUP_FILE] marker, but return 204 to client
4. Orphaned files are acceptable — they occupy disk space but don't affect application behavior

This resolves DQ-W02-001 with the "soft fail" approach consistent with single-user MVP philosophy.

### 10. Memory-Safe Upload Handling
```go
const maxMemory = 32 << 20 // 32MB memory buffer, rest to temp files
err := r.ParseMultipartForm(maxBytesLimit) // maxBytesLimit from TDEC-008
```

## Updated sqlc Queries

```sql
-- name: CreateExercise :one
INSERT INTO exercises (name, muscle_groups, description, personal_notes, working_weight, is_active)
VALUES ($1, $2, $3, $4, $5, COALESCE(sqlc.narg('is_active'), TRUE))
RETURNING id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at;

-- name: GetExerciseByID :one
SELECT id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at
FROM exercises WHERE id = $1;

-- name: ListExercises :many
SELECT id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at
FROM exercises
WHERE (sqlc.narg('include_inactive')::boolean IS TRUE OR is_active = TRUE)
AND (sqlc.narg('after_created_at')::timestamptz IS NULL OR created_at < sqlc.narg('after_created_at')::timestamptz)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit_rows');

-- name: CountExercises :one
SELECT COUNT(*) FROM exercises
WHERE (sqlc.narg('include_inactive')::boolean IS TRUE OR is_active = TRUE);

-- name: ListAllExercises :many
SELECT id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at
FROM exercises
WHERE (sqlc.narg('include_inactive')::boolean IS TRUE OR is_active = TRUE)
ORDER BY name ASC;

-- name: UpdateExercise :one
UPDATE exercises
SET name = COALESCE(sqlc.narg('name'), name),
    muscle_groups = COALESCE(sqlc.narg('muscle_groups'), muscle_groups),
    description = COALESCE(sqlc.narg('description'), description),
    personal_notes = COALESCE(sqlc.narg('personal_notes'), personal_notes),
    working_weight = COALESCE(sqlc.narg('working_weight'), working_weight),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at;

-- name: SoftDeleteExercise :one
UPDATE exercises
SET is_active = FALSE, updated_at = NOW()
WHERE id = $1
RETURNING id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at;

-- name: CreateExerciseMedia :one
INSERT INTO exercise_media (exercise_id, media_type, file_path, original_file_name, mime_type, size_bytes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, exercise_id, media_type, file_path, original_file_name, mime_type, size_bytes, created_at;

-- name: GetExerciseMediaByID :one
SELECT id, exercise_id, media_type, file_path, original_file_name, mime_type, size_bytes, created_at
FROM exercise_media WHERE id = $1;

-- name: ListExerciseMediaByExerciseID :many
SELECT id, exercise_id, media_type, file_path, original_file_name, mime_type, size_bytes, created_at
FROM exercise_media WHERE exercise_id = $1 ORDER BY created_at DESC;

-- name: DeleteExerciseMedia :exec
DELETE FROM exercise_media WHERE id = $1;
```