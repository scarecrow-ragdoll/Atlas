-- +goose Up
CREATE TABLE progress_photos (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    check_in_id         UUID NOT NULL REFERENCES body_check_ins(id) ON DELETE CASCADE,
    file_path           TEXT NOT NULL,
    original_file_name  TEXT NOT NULL,
    mime_type           TEXT NOT NULL,
    size_bytes          BIGINT NOT NULL,
    angle               TEXT,
    label               TEXT,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_progress_photo_checkin ON progress_photos (check_in_id);

-- +goose Down
DROP TABLE IF EXISTS progress_photos;

-- FILE: apps/api/internal/repository/postgres/migrations/00088_progress_photos.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add progress_photos table for WAVE-04 progress photo tracking.
--   SCOPE: Photos nested under body_check_ins with ON DELETE CASCADE. Stores file metadata, storage path, angle, label, and notes.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00086_body_check_ins.sql (body_check_ins).
--   LINKS: M-API / V-M-API / WAVE-04.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added progress_photos table for WAVE-04.
-- END_CHANGE_SUMMARY