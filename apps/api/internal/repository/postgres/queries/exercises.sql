-- FILE: apps/api/internal/repository/postgres/queries/exercises.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc queries for the exercises and exercise_media tables in WAVE-02 Exercise Library.
--   SCOPE: Exercise CRUD, list with pagination, allExercises, soft archive/restore, media CRUD, and user-scoped access.
--   DEPENDS: exercises table (00081_exercises.sql), exercise_media table (00082_exercise_media.sql).
--   LINKS: M-API / V-M-API / WAVE-02.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added exercise queries for WAVE-02.
-- END_CHANGE_SUMMARY

-- name: CreateExercise :one
INSERT INTO exercises (user_id, name, muscle_groups, description, personal_notes, working_weight)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at;

-- name: GetExerciseByID :one
SELECT id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at
FROM exercises
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListExercises :many
SELECT id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at
FROM exercises
WHERE user_id = $1 AND is_active = $2
ORDER BY name ASC
LIMIT $3;

-- name: ListExercisesCursor :many
SELECT id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at
FROM exercises
WHERE user_id = $1 AND is_active = $2 AND name > $3
ORDER BY name ASC
LIMIT $4;

-- name: ListAllExercises :many
SELECT id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at
FROM exercises
WHERE user_id = $1 AND (NOT $2::bool OR is_active = true)
ORDER BY name ASC;

-- name: CountExercises :one
SELECT COUNT(*)
FROM exercises
WHERE user_id = $1 AND is_active = $2;

-- name: UpdateExercise :one
UPDATE exercises
SET name = COALESCE($3, name),
    muscle_groups = COALESCE($4, muscle_groups),
    description = COALESCE($5, description),
    personal_notes = COALESCE($6, personal_notes),
    working_weight = COALESCE($7, working_weight),
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at;

-- name: ArchiveExercise :one
UPDATE exercises
SET is_active = false, updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at;

-- name: RestoreExercise :one
UPDATE exercises
SET is_active = true, updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at;

-- name: CreateExerciseMedia :one
INSERT INTO exercise_media (user_id, exercise_id, file_name, file_path, mime_type, file_size)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, exercise_id, file_name, file_path, mime_type, file_size, created_at;

-- name: GetExerciseMediaByID :one
SELECT id, user_id, exercise_id, file_name, file_path, mime_type, file_size, created_at
FROM exercise_media
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListExerciseMediaByExercise :many
SELECT id, user_id, exercise_id, file_name, file_path, mime_type, file_size, created_at
FROM exercise_media
WHERE exercise_id = $1 AND user_id = $2
ORDER BY created_at ASC;

-- name: DeleteExerciseMedia :one
DELETE FROM exercise_media
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, exercise_id, file_name, file_path, mime_type, file_size, created_at;