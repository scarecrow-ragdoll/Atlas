-- +goose Up
CREATE TABLE ai_exports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES atlas_users(id) ON DELETE CASCADE,
    date_range_start DATE NOT NULL,
    date_range_end DATE NOT NULL,
    include_photos BOOLEAN NOT NULL DEFAULT false,
    include_nutrition BOOLEAN NOT NULL DEFAULT true,
    include_cardio BOOLEAN NOT NULL DEFAULT true,
    include_measurements BOOLEAN NOT NULL DEFAULT true,
    user_comment TEXT,
    generated_prompt TEXT NOT NULL,
    export_file_path TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS ai_exports;