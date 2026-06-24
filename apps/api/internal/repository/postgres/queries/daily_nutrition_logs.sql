-- FILE: apps/api/internal/repository/postgres/queries/daily_nutrition_logs.sql
-- VERSION: 1.0.2
-- START_MODULE_CONTRACT
--   PURPOSE: sqlc queries for factual daily nutrition logs and entries.
--   SCOPE: User-scoped log load/upsert, entry CRUD with parent and product ownership checks, and date-range export reads.
--   DEPENDS: 00096_daily_nutrition_logs.sql.
--   LINKS: M-API-NUTRITION / V-M-API / V-M-API-NUTRITION
--   ROLE: SCRIPT
--   MAP_MODE: LOCALS
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   CreateDailyNutritionLog - Gets or creates one user/date factual nutrition log.
--   CreateDailyNutritionEntry - Creates a product snapshot entry only when log and product belong to the same user.
--   ListDailyNutritionEntriesByLog - Lists entries only through parent log ownership.
--   UpdateDailyNutritionEntry/DeleteDailyNutritionEntry - Mutates entries through parent log ownership checks.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.2 - Made daily log get-or-create race-safe while preserving existing notes and updated_at.
-- END_CHANGE_SUMMARY

-- name: CreateDailyNutritionLog :one
INSERT INTO daily_nutrition_logs (user_id, date, notes)
VALUES (sqlc.arg(user_id), sqlc.arg(date), sqlc.arg(notes))
ON CONFLICT (user_id, date)
DO UPDATE SET
  notes = daily_nutrition_logs.notes,
  updated_at = daily_nutrition_logs.updated_at
RETURNING id, user_id, date, notes, created_at, updated_at;

-- name: GetDailyNutritionLogByDate :one
SELECT id, user_id, date, notes, created_at, updated_at
FROM daily_nutrition_logs
WHERE user_id = sqlc.arg(user_id) AND date = sqlc.arg(date)
LIMIT 1;

-- name: ListDailyNutritionLogsByRange :many
SELECT id, user_id, date, notes, created_at, updated_at
FROM daily_nutrition_logs
WHERE user_id = sqlc.arg(user_id) AND date >= sqlc.arg(start_date) AND date <= sqlc.arg(end_date)
ORDER BY date ASC;

-- name: UpdateDailyNutritionLogNotes :one
UPDATE daily_nutrition_logs
SET notes = sqlc.arg(notes), updated_at = now()
WHERE id = sqlc.arg(id) AND user_id = sqlc.arg(user_id)
RETURNING id, user_id, date, notes, created_at, updated_at;

-- name: CreateDailyNutritionEntry :one
INSERT INTO daily_nutrition_entries (
  daily_log_id, product_id, product_name_snapshot,
  calories_per_100g_snapshot, protein_per_100g_snapshot,
  fat_per_100g_snapshot, carbs_per_100g_snapshot,
  amount_grams, meal_label, notes, position
)
SELECT
  l.id, p.id, p.name,
  p.calories_per_100g, p.protein_per_100g,
  p.fat_per_100g, p.carbs_per_100g,
  sqlc.arg(amount_grams), sqlc.arg(meal_label), sqlc.arg(notes), sqlc.arg(position)
FROM daily_nutrition_logs l
JOIN nutrition_product p ON p.id = sqlc.arg(product_id) AND p.user_id = l.user_id
WHERE l.id = sqlc.arg(daily_log_id) AND l.user_id = sqlc.arg(user_id)
RETURNING id, daily_log_id, product_id, product_name_snapshot,
  calories_per_100g_snapshot, protein_per_100g_snapshot,
  fat_per_100g_snapshot, carbs_per_100g_snapshot,
  amount_grams, meal_label, notes, position, created_at, updated_at;

-- name: ListDailyNutritionEntriesByLog :many
SELECT e.id, e.daily_log_id, e.product_id, e.product_name_snapshot,
  e.calories_per_100g_snapshot, e.protein_per_100g_snapshot,
  e.fat_per_100g_snapshot, e.carbs_per_100g_snapshot,
  e.amount_grams, e.meal_label, e.notes, e.position, e.created_at, e.updated_at
FROM daily_nutrition_entries e
JOIN daily_nutrition_logs l ON l.id = e.daily_log_id
WHERE e.daily_log_id = sqlc.arg(daily_log_id) AND l.user_id = sqlc.arg(user_id)
ORDER BY e.position ASC, e.created_at ASC;

-- name: UpdateDailyNutritionEntry :one
UPDATE daily_nutrition_entries e
SET amount_grams = sqlc.arg(amount_grams),
    meal_label = sqlc.arg(meal_label),
    notes = sqlc.arg(notes),
    position = sqlc.arg(position),
    updated_at = now()
FROM daily_nutrition_logs l
WHERE e.id = sqlc.arg(id)
  AND e.daily_log_id = l.id
  AND l.user_id = sqlc.arg(user_id)
  AND e.daily_log_id = sqlc.arg(daily_log_id)
RETURNING e.id, e.daily_log_id, e.product_id, e.product_name_snapshot,
  e.calories_per_100g_snapshot, e.protein_per_100g_snapshot,
  e.fat_per_100g_snapshot, e.carbs_per_100g_snapshot,
  e.amount_grams, e.meal_label, e.notes, e.position, e.created_at, e.updated_at;

-- name: DeleteDailyNutritionEntry :one
DELETE FROM daily_nutrition_entries e
USING daily_nutrition_logs l
WHERE e.id = sqlc.arg(id)
  AND e.daily_log_id = l.id
  AND l.user_id = sqlc.arg(user_id)
RETURNING e.id, e.daily_log_id, e.product_id, e.product_name_snapshot,
  e.calories_per_100g_snapshot, e.protein_per_100g_snapshot,
  e.fat_per_100g_snapshot, e.carbs_per_100g_snapshot,
  e.amount_grams, e.meal_label, e.notes, e.position, e.created_at, e.updated_at;
