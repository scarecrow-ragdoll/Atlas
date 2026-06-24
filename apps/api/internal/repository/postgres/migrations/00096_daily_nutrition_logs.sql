-- FILE: apps/api/internal/repository/postgres/migrations/00096_daily_nutrition_logs.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Create factual daily nutrition logs and entry snapshots.
--   SCOPE: Adds daily_nutrition_logs and daily_nutrition_entries without dropping legacy override tables.
--   DEPENDS: nutrition_product and atlas_users tables.
--   LINKS: M-API-NUTRITION / V-M-API / V-M-API-NUTRITION
--   ROLE: SCRIPT
--   MAP_MODE: LOCALS
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   daily_nutrition_logs - Date-scoped factual nutrition log per Atlas user.
--   daily_nutrition_entries - Product-and-grams facts with product macro snapshots captured at write time.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added factual daily nutrition tables with local/dev rollback safety guidance.
-- END_CHANGE_SUMMARY
-- ROLLBACK_SAFETY:
--   Goose Down is allowed in local/dev before factual daily nutrition writes are accepted.
--   After factual daily nutrition entries exist in any shared environment, destructive rollback is prohibited unless a backup/export has been generated and restore has been verified.
--   Production rollback preference is forward-fix or compatibility fallback, not dropping daily_nutrition_logs or daily_nutrition_entries.
--   Verification evidence must state whether the environment is pre-write, backed up, or not rollback-safe.

-- +goose Up
CREATE TABLE daily_nutrition_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES atlas_users(id),
    date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, date)
);

CREATE INDEX idx_daily_nutrition_logs_user_date ON daily_nutrition_logs (user_id, date);

CREATE TABLE daily_nutrition_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    daily_log_id UUID NOT NULL REFERENCES daily_nutrition_logs(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES nutrition_product(id),
    product_name_snapshot VARCHAR NOT NULL,
    calories_per_100g_snapshot REAL NOT NULL CHECK (calories_per_100g_snapshot >= 0),
    protein_per_100g_snapshot REAL NOT NULL CHECK (protein_per_100g_snapshot >= 0),
    fat_per_100g_snapshot REAL NOT NULL CHECK (fat_per_100g_snapshot >= 0),
    carbs_per_100g_snapshot REAL NOT NULL CHECK (carbs_per_100g_snapshot >= 0),
    amount_grams REAL NOT NULL CHECK (amount_grams > 0),
    meal_label VARCHAR,
    notes TEXT,
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_daily_nutrition_entries_log ON daily_nutrition_entries (daily_log_id, position, created_at);
CREATE INDEX idx_daily_nutrition_entries_product ON daily_nutrition_entries (product_id);

-- +goose Down
DROP TABLE IF EXISTS daily_nutrition_entries;
DROP TABLE IF EXISTS daily_nutrition_logs;
