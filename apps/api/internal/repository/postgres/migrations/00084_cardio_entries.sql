-- +goose Up
CREATE TABLE cardio_entries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES atlas_users(id),
    daily_log_id    UUID NOT NULL REFERENCES daily_logs(id) ON DELETE CASCADE,
    cardio_type     TEXT NOT NULL,
    duration_minutes INTEGER NOT NULL CHECK (duration_minutes > 0),
    avg_pulse       INTEGER CHECK (avg_pulse IS NULL OR avg_pulse > 0),
    heart_rate_zone TEXT,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_cardio_entries_daily_log ON cardio_entries (daily_log_id);
CREATE INDEX idx_cardio_entries_user ON cardio_entries (user_id);

-- +goose Down
DROP TABLE IF EXISTS cardio_entries;

-- FILE: apps/api/internal/repository/postgres/migrations/00084_cardio_entries.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add cardio_entries table for WAVE-04 cardio tracking.
--   SCOPE: User-scoped cardio entries with DailyLog FK, cardio type, duration, optional pulse/zone/notes. CHECK constraints on positive duration and pulse.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql (atlas_users), apps/api/internal/repository/postgres/migrations/00083_daily_logs.sql (daily_logs).
--   LINKS: M-API / V-M-API / WAVE-04.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added cardio_entries table for WAVE-04.
-- END_CHANGE_SUMMARY