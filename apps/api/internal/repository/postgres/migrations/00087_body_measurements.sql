-- +goose Up
CREATE TABLE body_measurements (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    check_in_id       UUID NOT NULL REFERENCES body_check_ins(id) ON DELETE CASCADE,
    measurement_type  TEXT NOT NULL,
    side              TEXT,
    value             REAL NOT NULL CHECK (value > 0),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(check_in_id, measurement_type, side)
);

CREATE INDEX idx_body_measurement_checkin ON body_measurements (check_in_id);

-- +goose Down
DROP TABLE IF EXISTS body_measurements;

-- FILE: apps/api/internal/repository/postgres/migrations/00087_body_measurements.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add body_measurements table for WAVE-04 body measurement tracking.
--   SCOPE: Measurements nested under body_check_ins with measurement type, optional side (for paired types), and positive value CHECK. Unique constraint on (check_in_id, measurement_type, side).
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00086_body_check_ins.sql (body_check_ins).
--   LINKS: M-API / V-M-API / WAVE-04.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added body_measurements table for WAVE-04.
-- END_CHANGE_SUMMARY