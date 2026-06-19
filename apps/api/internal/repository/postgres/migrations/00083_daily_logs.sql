-- FILE: apps/api/internal/repository/postgres/migrations/00083_daily_logs.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add versioned DailyLog aggregate table for WAVE-03 Workout Diary.
--   SCOPE: User-scoped daily container with unique date, nullable notes, version for optimistic concurrency, and timestamps; excludes cardio and body weight.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql.
--   LINKS: M-API / V-M-API / WAVE-03.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   daily_logs - Canonical daily aggregate container for WAVE-03 strength workouts.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added DailyLog table for WAVE-03.
-- END_CHANGE_SUMMARY

-- +goose Up
CREATE TABLE daily_logs (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES atlas_users(id),
    date       DATE NOT NULL,
    notes      TEXT,
    version    INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_daily_logs_user_date UNIQUE (user_id, date),
    CONSTRAINT chk_daily_logs_version CHECK (version >= 0)
);

CREATE INDEX idx_daily_logs_user_date ON daily_logs (user_id, date);
CREATE INDEX idx_daily_logs_user_date_desc ON daily_logs (user_id, date DESC);

-- +goose Down
DROP TABLE IF EXISTS daily_logs;
