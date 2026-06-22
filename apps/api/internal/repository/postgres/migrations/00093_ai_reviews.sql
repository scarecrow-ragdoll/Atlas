-- +goose Up
CREATE TABLE ai_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES atlas_users(id) ON DELETE CASCADE,
    date_range_start DATE NOT NULL,
    date_range_end DATE NOT NULL,
    ai_response_text TEXT NOT NULL,
    user_notes TEXT,
    planned_actions TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ai_reviews_user_id_date_range ON ai_reviews(user_id, date_range_start, date_range_end);

-- +goose Down
DROP TABLE IF EXISTS ai_reviews;