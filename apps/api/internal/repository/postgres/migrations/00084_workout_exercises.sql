-- FILE: apps/api/internal/repository/postgres/migrations/00084_workout_exercises.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add workout exercise instances within DailyLog for WAVE-03.
--   SCOPE: Ordered user-scoped exercise instances with working weight snapshots and notes; allows duplicate exercise_id values per day.
--   DEPENDS: daily_logs, exercises, atlas_users.
--   LINKS: M-API / V-M-API / WAVE-03.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   workout_exercises - Ordered strength exercise instances attached to a DailyLog.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added workout_exercises table for WAVE-03.
-- END_CHANGE_SUMMARY

-- +goose Up
CREATE TABLE workout_exercises (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID NOT NULL REFERENCES atlas_users(id),
    daily_log_id            UUID NOT NULL REFERENCES daily_logs(id) ON DELETE CASCADE,
    exercise_id             UUID NOT NULL REFERENCES exercises(id) ON DELETE RESTRICT,
    position                INTEGER NOT NULL,
    working_weight_snapshot REAL,
    notes                   TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_workout_exercises_position CHECK (position > 0),
    CONSTRAINT chk_workout_exercises_working_weight_snapshot CHECK (working_weight_snapshot IS NULL OR working_weight_snapshot > 0),
    CONSTRAINT uq_workout_exercises_daily_log_position UNIQUE (daily_log_id, position)
);

CREATE INDEX idx_workout_exercises_user_daily_log ON workout_exercises (user_id, daily_log_id);
CREATE INDEX idx_workout_exercises_exercise ON workout_exercises (exercise_id);

-- +goose Down
DROP TABLE IF EXISTS workout_exercises;
