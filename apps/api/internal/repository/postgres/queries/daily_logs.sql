-- FILE: apps/api/internal/repository/postgres/queries/daily_logs.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc queries for the daily_logs table used by WAVE-04 cardio FK.
--   SCOPE: GetOrCreateDailyLog (by user_id + date), GetDailyLogByDate. Minimal — WAVE-03 adds full aggregate queries.
--   DEPENDS: daily_logs table (00083_daily_logs.sql).
--   LINKS: M-API / V-M-API / WAVE-04.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added daily_log queries for WAVE-04 cardio support.
-- END_CHANGE_SUMMARY

-- name: GetDailyLogByDate :one
SELECT id, user_id, date, notes, created_at, updated_at
FROM daily_logs
WHERE user_id = $1 AND date = $2
LIMIT 1;

-- name: CreateDailyLog :one
INSERT INTO daily_logs (user_id, date, notes)
VALUES ($1, $2, $3)
RETURNING id, user_id, date, notes, created_at, updated_at;

-- name: DeleteDailyLog :one
DELETE FROM daily_logs
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, date, notes, created_at, updated_at;