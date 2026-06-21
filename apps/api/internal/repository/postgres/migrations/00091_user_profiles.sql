-- +goose Up
CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES atlas_users(id) ON DELETE CASCADE,
    goal TEXT,
    height REAL,
    birth_date DATE,
    training_experience TEXT,
    current_training_split TEXT,
    preferred_progression_style TEXT,
    nutrition_strategy TEXT,
    persistent_ai_context TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS user_profiles;