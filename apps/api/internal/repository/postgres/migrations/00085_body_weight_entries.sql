-- +goose Up
CREATE TABLE body_weight_entries (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES atlas_users(id),
    date       DATE NOT NULL,
    weight     REAL NOT NULL CHECK (weight > 0),
    source     TEXT NOT NULL,
    notes      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_body_weight_user_date ON body_weight_entries (user_id, date DESC);

-- +goose Down
DROP TABLE IF EXISTS body_weight_entries;

-- FILE: apps/api/internal/repository/postgres/migrations/00085_body_weight_entries.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add body_weight_entries table for WAVE-04 body weight tracking.
--   SCOPE: User-scoped weight entries by date with source enum, notes, and positive weight CHECK constraint. Multiple entries per date allowed (DDEC-W04-004).
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql (atlas_users).
--   LINKS: M-API / V-M-API / WAVE-04.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added body_weight_entries table for WAVE-04.
-- END_CHANGE_SUMMARY