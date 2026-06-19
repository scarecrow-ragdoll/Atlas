-- +goose Up
CREATE TABLE exercises (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES atlas_users(id),
    name            TEXT NOT NULL,
    muscle_groups   TEXT[],
    description     TEXT,
    personal_notes  TEXT,
    working_weight  REAL,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_working_weight CHECK (working_weight IS NULL OR working_weight > 0)
);

CREATE INDEX idx_exercises_user_active ON exercises (user_id, is_active);
CREATE INDEX idx_exercises_user_name ON exercises (user_id, name);
CREATE INDEX idx_exercises_user_created_at ON exercises (user_id, created_at);

-- +goose Down
DROP TABLE IF EXISTS exercises;

-- FILE: apps/api/internal/repository/postgres/migrations/00081_exercises.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add exercises table for WAVE-02 Exercise Library.
--   SCOPE: User-scoped exercise records with name, muscle groups, description, personal notes, working weight, and is_active soft-delete flag. Working weight has DB-level CHECK constraint.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql (atlas_users).
--   LINKS: M-API / V-M-API / WAVE-02.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   exercises - Stores user exercises with soft archive (is_active), working weight with CHECK constraint, user-scoped indexes.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added exercises table for WAVE-02.
-- END_CHANGE_SUMMARY
