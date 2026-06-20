-- name: UpsertDailyNutritionOverride :one
INSERT INTO daily_nutrition_override (user_id, date, notes)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, date)
DO UPDATE SET notes = COALESCE($3, daily_nutrition_override.notes),
              updated_at = now()
RETURNING id, user_id, date, notes, created_at, updated_at;

-- name: GetDailyNutritionOverrideByID :one
SELECT id, user_id, date, notes, created_at, updated_at
FROM daily_nutrition_override
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: GetDailyNutritionOverrideByDate :one
SELECT id, user_id, date, notes, created_at, updated_at
FROM daily_nutrition_override
WHERE user_id = $1 AND date = $2
LIMIT 1;

-- name: ListDailyNutritionOverridesByRange :many
SELECT id, user_id, date, notes, created_at, updated_at
FROM daily_nutrition_override
WHERE user_id = $1 AND date >= $2 AND date <= $3
ORDER BY date ASC;

-- name: UpdateDailyNutritionOverride :one
UPDATE daily_nutrition_override
SET notes = $3, updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, date, notes, created_at, updated_at;

-- name: DeleteDailyNutritionOverride :one
DELETE FROM daily_nutrition_override
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, date, notes, created_at, updated_at;