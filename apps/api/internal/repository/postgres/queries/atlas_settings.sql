-- FILE: apps/api/internal/repository/postgres/queries/atlas_settings.sql
-- VERSION: 1.0.1
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc queries for the atlas_settings table in the Atlas fitness tracker module.
--   SCOPE: FindByUserID, UpsertSettings, UpdatePinState operations; excludes domain queries for later waves.
--   DEPENDS: atlas_settings table in migration 00080_atlas_foundation.sql.
--   LINKS: M-API / V-M-API.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.1 - Removed ON CONFLICT DO NOTHING from user upsert since atlas_users has no unique constraint on display_name.
-- END_CHANGE_SUMMARY

-- name: GetAtlasSettingsByUserID :one
SELECT id, user_id, pin_enabled, pin_hash, units, default_ai_export_weeks, created_at, updated_at
FROM atlas_settings
WHERE user_id = $1
LIMIT 1;

-- name: UpsertAtlasSettings :one
INSERT INTO atlas_settings (user_id, pin_enabled, units, default_ai_export_weeks, updated_at)
VALUES ($1, $2, $3, $4, now())
ON CONFLICT (user_id)
DO UPDATE SET units = COALESCE(NULLIF($3, ''), atlas_settings.units),
              default_ai_export_weeks = CASE WHEN $4 <> 0 THEN $4 ELSE atlas_settings.default_ai_export_weeks END,
              updated_at = now()
RETURNING id, user_id, pin_enabled, pin_hash, units, default_ai_export_weeks, created_at, updated_at;

-- name: UpdateAtlasPinState :one
UPDATE atlas_settings
SET pin_enabled = $2,
    pin_hash = $3,
    updated_at = now()
WHERE user_id = $1
RETURNING id, user_id, pin_enabled, pin_hash, units, default_ai_export_weeks, created_at, updated_at;

-- name: GetAtlasDefaultUser :one
SELECT id, display_name, created_at, updated_at
FROM atlas_users
ORDER BY created_at ASC
LIMIT 1;

-- name: InsertAtlasDefaultUser :one
INSERT INTO atlas_users (display_name)
VALUES ($1)
RETURNING id;

-- name: CreateAtlasSettings :exec
INSERT INTO atlas_settings (user_id, pin_enabled, units, default_ai_export_weeks)
VALUES ($1, false, 'metric', 4)
ON CONFLICT (user_id) DO NOTHING;