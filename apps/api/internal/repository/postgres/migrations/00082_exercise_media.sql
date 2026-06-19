-- +goose Up
CREATE TABLE exercise_media (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES atlas_users(id),
    exercise_id  UUID NOT NULL REFERENCES exercises(id) ON DELETE NO ACTION,
    file_name    TEXT NOT NULL,
    file_path    TEXT NOT NULL,
    mime_type    TEXT NOT NULL,
    file_size    BIGINT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_exercise_media_user_exercise ON exercise_media (user_id, exercise_id);

-- +goose Down
DROP TABLE IF EXISTS exercise_media;

-- FILE: apps/api/internal/repository/postgres/migrations/00082_exercise_media.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add exercise_media table for WAVE-02 Exercise Library media attachments.
--   SCOPE: User-scoped media records linked to exercises with ON DELETE NO ACTION (archive does NOT cascade). Stores file metadata and storage path.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00081_exercises.sql (exercises), apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql (atlas_users).
--   LINKS: M-API / V-M-API / WAVE-02.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   exercise_media - Stores exercise media metadata with user-scoped index and ON DELETE NO ACTION FK to exercises.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added exercise_media table for WAVE-02.
-- END_CHANGE_SUMMARY
