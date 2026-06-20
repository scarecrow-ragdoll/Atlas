-- FILE: apps/api/internal/repository/postgres/queries/progress_photos.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc queries for the progress_photos table in WAVE-04.
--   SCOPE: Progress photo CRUD, list by check-in ID, user-scoped via check-in join.
--   DEPENDS: progress_photos table (00088_progress_photos.sql), body_check_ins table (00086_body_check_ins.sql).
--   LINKS: M-API / V-M-API / WAVE-04.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added progress photo queries for WAVE-04.
-- END_CHANGE_SUMMARY

-- name: CreateProgressPhoto :one
INSERT INTO progress_photos (check_in_id, file_path, original_file_name, mime_type, size_bytes, angle, label, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, check_in_id, file_path, original_file_name, mime_type, size_bytes, angle, label, notes, created_at, updated_at;

-- name: GetProgressPhotoByID :one
SELECT p.id, p.check_in_id, p.file_path, p.original_file_name, p.mime_type, p.size_bytes, p.angle, p.label, p.notes, p.created_at, p.updated_at
FROM progress_photos p
JOIN body_check_ins c ON c.id = p.check_in_id
WHERE p.id = $1 AND c.user_id = $2
LIMIT 1;

-- name: ListProgressPhotosByCheckIn :many
SELECT p.id, p.check_in_id, p.file_path, p.original_file_name, p.mime_type, p.size_bytes, p.angle, p.label, p.notes, p.created_at, p.updated_at
FROM progress_photos p
JOIN body_check_ins c ON c.id = p.check_in_id
WHERE p.check_in_id = $1 AND c.user_id = $2
ORDER BY p.created_at ASC;

-- name: DeleteProgressPhoto :one
DELETE FROM progress_photos p
USING body_check_ins c
WHERE p.check_in_id = c.id AND p.id = $1 AND c.user_id = $2
RETURNING p.id, p.check_in_id, p.file_path, p.original_file_name, p.mime_type, p.size_bytes, p.angle, p.label, p.notes, p.created_at, p.updated_at;

-- name: CountProgressPhotosByCheckIn :one
SELECT COUNT(*)
FROM progress_photos
WHERE check_in_id = $1;