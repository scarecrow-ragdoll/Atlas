-- name: UpsertNutritionTemplate :one
INSERT INTO nutrition_template (user_id, week_start_date, title, notes)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, week_start_date)
DO UPDATE SET title = COALESCE($3, nutrition_template.title),
              notes = COALESCE($4, nutrition_template.notes),
              updated_at = now()
RETURNING id, user_id, week_start_date, title, notes, created_at, updated_at;

-- name: GetNutritionTemplateByID :one
SELECT id, user_id, week_start_date, title, notes, created_at, updated_at
FROM nutrition_template
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: GetNutritionTemplateByWeek :one
SELECT id, user_id, week_start_date, title, notes, created_at, updated_at
FROM nutrition_template
WHERE user_id = $1 AND week_start_date = $2
LIMIT 1;

-- name: ListNutritionTemplatesByRange :many
SELECT id, user_id, week_start_date, title, notes, created_at, updated_at
FROM nutrition_template
WHERE user_id = $1 AND week_start_date >= $2 AND week_start_date <= $3
ORDER BY week_start_date ASC;

-- name: UpdateNutritionTemplate :one
UPDATE nutrition_template
SET title = $3,
    notes = $4,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, week_start_date, title, notes, created_at, updated_at;

-- name: DeleteNutritionTemplate :one
DELETE FROM nutrition_template
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, week_start_date, title, notes, created_at, updated_at;