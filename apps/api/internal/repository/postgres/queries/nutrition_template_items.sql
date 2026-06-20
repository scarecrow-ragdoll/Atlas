-- name: CreateNutritionTemplateItem :one
INSERT INTO nutrition_template_item (template_id, product_id, amount_grams, meal_label, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, template_id, product_id, amount_grams, meal_label, notes, created_at, updated_at;

-- name: GetNutritionTemplateItemByID :one
SELECT id, template_id, product_id, amount_grams, meal_label, notes, created_at, updated_at
FROM nutrition_template_item
WHERE id = $1
LIMIT 1;

-- name: ListNutritionTemplateItemsByTemplate :many
SELECT id, template_id, product_id, amount_grams, meal_label, notes, created_at, updated_at
FROM nutrition_template_item
WHERE template_id = $1
ORDER BY created_at ASC;

-- name: UpdateNutritionTemplateItem :one
UPDATE nutrition_template_item
SET amount_grams = $2,
    meal_label = $3,
    notes = $4,
    updated_at = now()
WHERE id = $1
RETURNING id, template_id, product_id, amount_grams, meal_label, notes, created_at, updated_at;

-- name: DeleteNutritionTemplateItem :one
DELETE FROM nutrition_template_item
WHERE id = $1
RETURNING id, template_id, product_id, amount_grams, meal_label, notes, created_at, updated_at;