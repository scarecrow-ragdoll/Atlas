-- +goose Up
CREATE TABLE backup_archives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES atlas_users(id) ON DELETE CASCADE,
    include_media BOOLEAN NOT NULL DEFAULT false,
    size_bytes BIGINT NOT NULL DEFAULT 0,
    entity_counts JSONB NOT NULL DEFAULT '{}',
    archive_path TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_backup_archives_user_id ON backup_archives(user_id);

-- +goose Down
DROP TABLE IF EXISTS backup_archives;
