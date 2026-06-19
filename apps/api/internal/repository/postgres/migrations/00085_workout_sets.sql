-- FILE: apps/api/internal/repository/postgres/migrations/00085_workout_sets.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add workout sets for WAVE-03 strength workout logging.
--   SCOPE: Ordered sets with weight, reps, optional RPE/RIR, and notes attached to workout_exercises.
--   DEPENDS: workout_exercises.
--   LINKS: M-API / V-M-API / WAVE-03.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   workout_sets - Ordered strength set rows attached to workout exercise instances.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added workout_sets table for WAVE-03.
-- END_CHANGE_SUMMARY

-- +goose Up
CREATE TABLE workout_sets (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_exercise_id UUID NOT NULL REFERENCES workout_exercises(id) ON DELETE CASCADE,
    set_number          INTEGER NOT NULL,
    weight              REAL NOT NULL,
    reps                INTEGER NOT NULL,
    rpe                 REAL,
    rir                 INTEGER,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_workout_sets_set_number CHECK (set_number > 0),
    CONSTRAINT chk_workout_sets_weight CHECK (weight > 0),
    CONSTRAINT chk_workout_sets_reps CHECK (reps > 0),
    CONSTRAINT chk_workout_sets_rpe CHECK (rpe IS NULL OR (rpe >= 1 AND rpe <= 10)),
    CONSTRAINT chk_workout_sets_rir CHECK (rir IS NULL OR (rir >= 0 AND rir <= 10)),
    CONSTRAINT uq_workout_sets_exercise_set_number UNIQUE (workout_exercise_id, set_number)
);

CREATE INDEX idx_workout_sets_workout_exercise ON workout_sets (workout_exercise_id);

-- +goose Down
DROP TABLE IF EXISTS workout_sets;
