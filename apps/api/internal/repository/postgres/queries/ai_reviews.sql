-- FILE: apps/api/internal/repository/postgres/queries/ai_reviews.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc queries for the ai_reviews table in WAVE-08.
--   SCOPE: AiReview CRUD, list by user ID, list by user ID and date range, user-scoped access.
--   DEPENDS: ai_reviews table (00093_ai_reviews.sql).
--   LINKS: M-API / V-M-API / WAVE-08.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added AiReview queries for WAVE-08.
-- END_CHANGE_SUMMARY

-- name: CreateAiReview :one
INSERT INTO ai_reviews (user_id, date_range_start, date_range_end, ai_response_text, user_notes, planned_actions)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, date_range_start, date_range_end, ai_response_text, user_notes, planned_actions, created_at, updated_at;

-- name: GetAiReviewByID :one
SELECT id, user_id, date_range_start, date_range_end, ai_response_text, user_notes, planned_actions, created_at, updated_at
FROM ai_reviews
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListAiReviewsByUserID :many
SELECT id, user_id, date_range_start, date_range_end, ai_response_text, user_notes, planned_actions, created_at, updated_at
FROM ai_reviews
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListAiReviewsByUserIDAndDateRange :many
SELECT id, user_id, date_range_start, date_range_end, ai_response_text, user_notes, planned_actions, created_at, updated_at
FROM ai_reviews
WHERE user_id = $1
  AND date_range_start >= $2
  AND date_range_end <= $3
ORDER BY created_at DESC;

-- name: UpdateAiReview :one
UPDATE ai_reviews
SET date_range_start = $2,
    date_range_end = $3,
    ai_response_text = $4,
    user_notes = $5,
    planned_actions = $6,
    updated_at = NOW()
WHERE id = $1 AND user_id = $7
RETURNING id, user_id, date_range_start, date_range_end, ai_response_text, user_notes, planned_actions, created_at, updated_at;

-- name: DeleteAiReview :one
DELETE FROM ai_reviews
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, date_range_start, date_range_end, ai_response_text, user_notes, planned_actions, created_at, updated_at;