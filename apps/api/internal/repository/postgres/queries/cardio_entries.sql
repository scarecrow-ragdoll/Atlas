-- FILE: apps/api/internal/repository/postgres/queries/cardio_entries.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc queries for the cardio_entries table in WAVE-04.
--   SCOPE: Cardio entry CRUD, list by daily log ID, user-scoped access.
--   DEPENDS: cardio_entries table (00084_cardio_entries.sql), daily_logs table (00083_daily_logs.sql).
--   LINKS: M-API / V-M-API / WAVE-04.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added cardio entry queries for WAVE-04.
-- END_CHANGE_SUMMARY

-- name: CreateCardioEntry :one
INSERT INTO cardio_entries (user_id, daily_log_id, cardio_type, duration_minutes, avg_pulse, heart_rate_zone, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, daily_log_id, cardio_type, duration_minutes, avg_pulse, heart_rate_zone, notes, created_at, updated_at;

-- name: GetCardioEntryByID :one
SELECT id, user_id, daily_log_id, cardio_type, duration_minutes, avg_pulse, heart_rate_zone, notes, created_at, updated_at
FROM cardio_entries
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListCardioEntriesByDailyLog :many
SELECT id, user_id, daily_log_id, cardio_type, duration_minutes, avg_pulse, heart_rate_zone, notes, created_at, updated_at
FROM cardio_entries
WHERE daily_log_id = $1 AND user_id = $2
ORDER BY created_at ASC;

-- name: UpdateCardioEntry :one
UPDATE cardio_entries
SET cardio_type = COALESCE($3, cardio_type),
    duration_minutes = COALESCE($4, duration_minutes),
    avg_pulse = $5,
    heart_rate_zone = $6,
    notes = $7,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, daily_log_id, cardio_type, duration_minutes, avg_pulse, heart_rate_zone, notes, created_at, updated_at;

-- name: DeleteCardioEntry :one
DELETE FROM cardio_entries
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, daily_log_id, cardio_type, duration_minutes, avg_pulse, heart_rate_zone, notes, created_at, updated_at;