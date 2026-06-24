-- FILE: apps/api/internal/repository/postgres/queries/nutrition_template_items.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: sqlc queries for nutrition_template_item table.
--   SCOPE: Template item CRUD plus user-scoped item lookup variants for repository ownership checks.
--   DEPENDS: 00090_nutrition_tables.sql migration.
--   LINKS: M-API-NUTRITION / V-M-API-NUTRITION
--   ROLE: SCRIPT
--   MAP_MODE: LOCALS
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   CreateNutritionTemplateItem - Creates one template-owned nutrition item.
--   GetNutritionTemplateItemByIDForUser - Loads one item only through its parent template owner.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added file-local GRACE contract and user-scoped template item lookup.
-- END_CHANGE_SUMMARY

-- name: CreateNutritionTemplateItem :one
INSERT INTO nutrition_template_item (template_id, product_id, amount_grams, meal_label, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, template_id, product_id, amount_grams, meal_label, notes, created_at, updated_at;

-- name: GetNutritionTemplateItemByID :one
SELECT id, template_id, product_id, amount_grams, meal_label, notes, created_at, updated_at
FROM nutrition_template_item
WHERE id = $1
LIMIT 1;

-- name: GetNutritionTemplateItemByIDForUser :one
SELECT i.id, i.template_id, i.product_id, i.amount_grams, i.meal_label, i.notes, i.created_at, i.updated_at
FROM nutrition_template_item i
JOIN nutrition_template t ON t.id = i.template_id
WHERE i.id = $1 AND t.user_id = $2
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
