-- FILE: apps/api/internal/repository/postgres/queries/week_flags.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc queries for the week_flags table in WAVE-04.
--   SCOPE: Week flag CRUD, list by week start date, user-scoped access.
--   DEPENDS: week_flags table (00089_week_flags.sql).
--   LINKS: M-API / V-M-API / WAVE-04.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added week flag queries for WAVE-04.
-- END_CHANGE_SUMMARY

-- name: CreateWeekFlag :one
INSERT INTO week_flags (user_id, week_start_date, flag_type, notes)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, week_start_date, flag_type, notes, created_at, updated_at;

-- name: GetWeekFlagByID :one
SELECT id, user_id, week_start_date, flag_type, notes, created_at, updated_at
FROM week_flags
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListWeekFlagsByWeekStart :many
SELECT id, user_id, week_start_date, flag_type, notes, created_at, updated_at
FROM week_flags
WHERE week_start_date = $1 AND user_id = $2
ORDER BY flag_type ASC;

-- name: DeleteWeekFlag :one
DELETE FROM week_flags
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, week_start_date, flag_type, notes, created_at, updated_at;