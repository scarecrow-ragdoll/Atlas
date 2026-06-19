-- +goose Up
CREATE TABLE atlas_users (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_name TEXT NOT NULL DEFAULT 'Default User',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE atlas_settings (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID NOT NULL REFERENCES atlas_users(id),
    pin_enabled             BOOLEAN NOT NULL DEFAULT false,
    pin_hash                TEXT,
    units                   TEXT NOT NULL DEFAULT 'metric',
    default_ai_export_weeks INT NOT NULL DEFAULT 4,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id)
);

-- +goose Down
DROP TABLE IF EXISTS atlas_settings;
DROP TABLE IF EXISTS atlas_users;

-- FILE: apps/api/internal/repository/postgres/migrations/00079_atlas_foundation.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add atlas_users and atlas_settings tables for the Atlas fitness tracker module.
--   SCOPE: Default user identity, settings with PIN hash and preferences; excludes domain tables (exercises, workouts, etc.) which belong to later waves.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00001_init.sql (uuid-ossp extension).
--   LINKS: M-API / V-M-API.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   atlas_users - Stores the Atlas default user identity.
--   atlas_settings - Stores PIN hash, PIN enabled state, and user preferences.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added Atlas foundation tables for WAVE-01.
-- END_CHANGE_SUMMARY