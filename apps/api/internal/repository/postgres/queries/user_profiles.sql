-- name: GetUserProfileByUserID :one
SELECT id, user_id, goal, height, birth_date, training_experience, current_training_split, preferred_progression_style, nutrition_strategy, persistent_ai_context, created_at, updated_at
FROM user_profiles
WHERE user_id = $1
LIMIT 1;

-- name: CreateUserProfile :one
INSERT INTO user_profiles (user_id, goal, height, birth_date, training_experience, current_training_split, preferred_progression_style, nutrition_strategy, persistent_ai_context)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, user_id, goal, height, birth_date, training_experience, current_training_split, preferred_progression_style, nutrition_strategy, persistent_ai_context, created_at, updated_at;

-- name: UpsertUserProfile :one
INSERT INTO user_profiles (user_id, goal, height, birth_date, training_experience, current_training_split, preferred_progression_style, nutrition_strategy, persistent_ai_context)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (user_id)
DO UPDATE SET
    goal = COALESCE($2, user_profiles.goal),
    height = COALESCE($3, user_profiles.height),
    birth_date = COALESCE($4, user_profiles.birth_date),
    training_experience = COALESCE($5, user_profiles.training_experience),
    current_training_split = COALESCE($6, user_profiles.current_training_split),
    preferred_progression_style = COALESCE($7, user_profiles.preferred_progression_style),
    nutrition_strategy = COALESCE($8, user_profiles.nutrition_strategy),
    persistent_ai_context = COALESCE($9, user_profiles.persistent_ai_context),
    updated_at = NOW()
RETURNING id, user_id, goal, height, birth_date, training_experience, current_training_split, preferred_progression_style, nutrition_strategy, persistent_ai_context, created_at, updated_at;