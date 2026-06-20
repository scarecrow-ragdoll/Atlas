-- name: CreateDailyNutritionOverrideItem :one
INSERT INTO daily_nutrition_override_item (override_id, product_id, amount_grams, operation, meal_label, notes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, override_id, product_id, amount_grams, operation, meal_label, notes, created_at, updated_at;

-- name: GetDailyNutritionOverrideItemByID :one
SELECT id, override_id, product_id, amount_grams, operation, meal_label, notes, created_at, updated_at
FROM daily_nutrition_override_item
WHERE id = $1
LIMIT 1;

-- name: ListDailyNutritionOverrideItemsByOverride :many
SELECT id, override_id, product_id, amount_grams, operation, meal_label, notes, created_at, updated_at
FROM daily_nutrition_override_item
WHERE override_id = $1
ORDER BY created_at ASC;

-- name: UpdateDailyNutritionOverrideItem :one
UPDATE daily_nutrition_override_item
SET amount_grams = $2,
    operation = $3,
    meal_label = $4,
    notes = $5,
    updated_at = now()
WHERE id = $1
RETURNING id, override_id, product_id, amount_grams, operation, meal_label, notes, created_at, updated_at;

-- name: DeleteDailyNutritionOverrideItem :one
DELETE FROM daily_nutrition_override_item
WHERE id = $1
RETURNING id, override_id, product_id, amount_grams, operation, meal_label, notes, created_at, updated_at;