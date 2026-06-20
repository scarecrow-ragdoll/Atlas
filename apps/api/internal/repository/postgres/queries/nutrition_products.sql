-- FILE: apps/api/internal/repository/postgres/queries/nutrition_products.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: sqlc CRUD queries for nutrition_product table.
--   SCOPE: Create, GetByID, ListActive, Update, SoftDelete, GetByIDIncludeInactive. All user-scoped.
--   DEPENDS: 00090_nutrition_tables.sql migration.
--   ROLE: SCRIPT
--   MAP_MODE: LOCALS
-- END_MODULE_CONTRACT

-- name: CreateNutritionProduct :one
INSERT INTO nutrition_product (user_id, name, calories_per_100g, protein_per_100g, fat_per_100g, carbs_per_100g, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, name, calories_per_100g, protein_per_100g, fat_per_100g, carbs_per_100g, notes, is_active, created_at, updated_at;

-- name: GetNutritionProductByID :one
SELECT id, user_id, name, calories_per_100g, protein_per_100g, fat_per_100g, carbs_per_100g, notes, is_active, created_at, updated_at
FROM nutrition_product
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListActiveNutritionProducts :many
SELECT id, user_id, name, calories_per_100g, protein_per_100g, fat_per_100g, carbs_per_100g, notes, is_active, created_at, updated_at
FROM nutrition_product
WHERE user_id = $1 AND is_active = true
ORDER BY name ASC;

-- name: UpdateNutritionProduct :one
UPDATE nutrition_product
SET name = $3,
    calories_per_100g = $4,
    protein_per_100g = $5,
    fat_per_100g = $6,
    carbs_per_100g = $7,
    notes = $8,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, name, calories_per_100g, protein_per_100g, fat_per_100g, carbs_per_100g, notes, is_active, created_at, updated_at;

-- name: SoftDeleteNutritionProduct :one
UPDATE nutrition_product
SET is_active = false, updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, name, calories_per_100g, protein_per_100g, fat_per_100g, carbs_per_100g, notes, is_active, created_at, updated_at;

-- name: GetNutritionProductByIDIncludeInactive :one
SELECT id, user_id, name, calories_per_100g, protein_per_100g, fat_per_100g, carbs_per_100g, notes, is_active, created_at, updated_at
FROM nutrition_product
WHERE id = $1 AND user_id = $2
LIMIT 1;