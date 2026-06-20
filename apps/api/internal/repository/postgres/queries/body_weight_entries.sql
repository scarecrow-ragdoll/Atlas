-- FILE: apps/api/internal/repository/postgres/queries/body_weight_entries.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc queries for the body_weight_entries table in WAVE-04.
--   SCOPE: Body weight entry CRUD, list by date range, latest entry, user-scoped access.
--   DEPENDS: body_weight_entries table (00085_body_weight_entries.sql).
--   LINKS: M-API / V-M-API / WAVE-04.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added body weight entry queries for WAVE-04.
-- END_CHANGE_SUMMARY

-- name: CreateBodyWeightEntry :one
INSERT INTO body_weight_entries (user_id, date, weight, source, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, date, weight, source, notes, created_at, updated_at;

-- name: GetBodyWeightEntryByID :one
SELECT id, user_id, date, weight, source, notes, created_at, updated_at
FROM body_weight_entries
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListBodyWeightEntriesByDateRange :many
SELECT id, user_id, date, weight, source, notes, created_at, updated_at
FROM body_weight_entries
WHERE user_id = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC, created_at DESC;

-- name: LatestBodyWeightEntry :one
SELECT id, user_id, date, weight, source, notes, created_at, updated_at
FROM body_weight_entries
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateBodyWeightEntry :one
UPDATE body_weight_entries
SET weight = COALESCE($3, weight),
    source = COALESCE($4, source),
    notes = $5,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, date, weight, source, notes, created_at, updated_at;

-- name: DeleteBodyWeightEntry :one
DELETE FROM body_weight_entries
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, date, weight, source, notes, created_at, updated_at;