-- FILE: apps/api/internal/repository/postgres/migrations/00090_nutrition_tables.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Create all 5 nutrition tables for WAVE-05 with FKs, CHECK constraints, indexes, and unique constraints.
--   SCOPE: Single goose migration creating nutrition_product, nutrition_template, nutrition_template_item, daily_nutrition_override, daily_nutrition_override_item. Reversible.
--   DEPENDS: atlas_users table from WAVE-01 foundation migration.
--   ROLE: SCRIPT
--   MAP_MODE: LOCALS
-- END_MODULE_CONTRACT

-- +goose Up
CREATE TABLE nutrition_product (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES atlas_users(id),
    name VARCHAR NOT NULL,
    calories_per_100g REAL NOT NULL CHECK (calories_per_100g >= 0),
    protein_per_100g REAL NOT NULL CHECK (protein_per_100g >= 0),
    fat_per_100g REAL NOT NULL CHECK (fat_per_100g >= 0),
    carbs_per_100g REAL NOT NULL CHECK (carbs_per_100g >= 0),
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_nutrition_product_user ON nutrition_product (user_id);

CREATE TABLE nutrition_template (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES atlas_users(id),
    week_start_date DATE NOT NULL,
    title VARCHAR,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, week_start_date)
);

CREATE INDEX idx_nutrition_template_week ON nutrition_template (user_id, week_start_date);

CREATE TABLE nutrition_template_item (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES nutrition_template(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES nutrition_product(id),
    amount_grams REAL NOT NULL CHECK (amount_grams > 0),
    meal_label VARCHAR,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_nutrition_template_item_template ON nutrition_template_item (template_id);

CREATE TABLE daily_nutrition_override (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES atlas_users(id),
    date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, date)
);

CREATE INDEX idx_nutrition_override_date ON daily_nutrition_override (user_id, date);

CREATE TABLE daily_nutrition_override_item (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    override_id UUID NOT NULL REFERENCES daily_nutrition_override(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES nutrition_product(id),
    amount_grams REAL NOT NULL CHECK (amount_grams > 0),
    operation VARCHAR NOT NULL CHECK (operation IN ('add', 'subtract', 'replace')),
    meal_label VARCHAR,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_nutrition_override_item_override ON daily_nutrition_override_item (override_id);

-- +goose Down
DROP TABLE IF EXISTS daily_nutrition_override_item;
DROP TABLE IF EXISTS daily_nutrition_override;
DROP TABLE IF EXISTS nutrition_template_item;
DROP TABLE IF EXISTS nutrition_template;
DROP TABLE IF EXISTS nutrition_product;