-- +goose Up
CREATE TABLE week_flags (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL REFERENCES atlas_users(id),
    week_start_date  DATE NOT NULL,
    flag_type        TEXT NOT NULL,
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(week_start_date, flag_type)
);

CREATE INDEX idx_week_flag_week ON week_flags (week_start_date);
CREATE INDEX idx_week_flag_user ON week_flags (user_id);

-- +goose Down
DROP TABLE IF EXISTS week_flags;

-- FILE: apps/api/internal/repository/postgres/migrations/00089_week_flags.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add week_flags table for WAVE-04 week-level flag tracking.
--   SCOPE: User-scoped week flags with UNIQUE constraint on (week_start_date, flag_type). Notes optional.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql (atlas_users).
--   LINKS: M-API / V-M-API / WAVE-04.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added week_flags table for WAVE-04.
-- END_CHANGE_SUMMARY