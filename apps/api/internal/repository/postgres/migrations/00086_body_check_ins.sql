-- +goose Up
CREATE TABLE body_check_ins (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES atlas_users(id),
    date                DATE NOT NULL UNIQUE,
    weight              REAL CHECK (weight IS NULL OR weight > 0),
    body_fat_percentage REAL CHECK (body_fat_percentage IS NULL OR (body_fat_percentage > 0 AND body_fat_percentage <= 100)),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_body_check_in_date ON body_check_ins (date DESC);
CREATE INDEX idx_body_check_in_user ON body_check_ins (user_id);

-- +goose Down
DROP TABLE IF EXISTS body_check_ins;

-- FILE: apps/api/internal/repository/postgres/migrations/00086_body_check_ins.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add body_check_ins table for WAVE-04 weekly body check-in tracking.
--   SCOPE: User-scoped check-ins by date with optional weight and body fat percentage. One check-in per date (UNIQUE constraint). CHECK constraints on weight and body fat percentage ranges.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql (atlas_users).
--   LINKS: M-API / V-M-API / WAVE-04.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added body_check_ins table for WAVE-04.
-- END_CHANGE_SUMMARY