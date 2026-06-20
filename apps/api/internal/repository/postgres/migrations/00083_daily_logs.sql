-- +goose Up
CREATE TABLE daily_logs (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES atlas_users(id),
    date       DATE NOT NULL,
    notes      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, date)
);

CREATE INDEX idx_daily_logs_user_date ON daily_logs (user_id, date DESC);

-- +goose Down
DROP TABLE IF EXISTS daily_logs;

-- FILE: apps/api/internal/repository/postgres/migrations/00083_daily_logs.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add daily_logs table as prerequisite for WAVE-04 cardio FK and WAVE-03 aggregate root.
--   SCOPE: Minimal daily_logs table with user_id/date unique constraint, notes, and timestamps. No versioning or workout aggregates — WAVE-03 adds those when deployed.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql (atlas_users).
--   LINKS: M-API / V-M-API / WAVE-04.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added daily_logs for WAVE-04 cardio FK dependency.
-- END_CHANGE_SUMMARY