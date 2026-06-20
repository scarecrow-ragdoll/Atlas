-- FILE: apps/api/internal/repository/postgres/queries/body_check_ins.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc queries for the body_check_ins table in WAVE-04.
--   SCOPE: Body check-in CRUD, list by date range, user-scoped access.
--   DEPENDS: body_check_ins table (00086_body_check_ins.sql).
--   LINKS: M-API / V-M-API / WAVE-04.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added body check-in queries for WAVE-04.
-- END_CHANGE_SUMMARY

-- name: CreateBodyCheckIn :one
INSERT INTO body_check_ins (user_id, date, weight, body_fat_percentage, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, date, weight, body_fat_percentage, notes, created_at, updated_at;

-- name: GetBodyCheckInByID :one
SELECT id, user_id, date, weight, body_fat_percentage, notes, created_at, updated_at
FROM body_check_ins
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListBodyCheckInsByDateRange :many
SELECT id, user_id, date, weight, body_fat_percentage, notes, created_at, updated_at
FROM body_check_ins
WHERE user_id = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC;

-- name: UpdateBodyCheckIn :one
UPDATE body_check_ins
SET weight = $3,
    body_fat_percentage = $4,
    notes = $5,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, date, weight, body_fat_percentage, notes, created_at, updated_at;

-- name: DeleteBodyCheckIn :one
DELETE FROM body_check_ins
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, date, weight, body_fat_percentage, notes, created_at, updated_at;