# WAVE-05 Data-Integration-Ops Planner Attempt 1

## Sources Read
- docs/product-verified/domain-model.md (Nutrition entities §25.13–§25.17)
- docs/product-verified/business-rules.md (RULE-006, RULE-010, RULE-011, RULE-018, RULE-019, RULE-020)
- docs/technical-verified/api-contracts.md
- docs/technical-verified/operations-observability.md
- apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql
- apps/api/internal/repository/postgres/queries/atlas_settings.sql
- apps/api/internal/atlas/graph/schema/settings.graphql
- docs/prd-wave-details/waves/wave-04.md (Data/API section reference)

## Selected Backend Wave Boundary
5 new PostgreSQL tables, 0 new REST endpoints, 5 new GraphQL query/mutation families, 1 macro calculation service, 5 new log markers.

## Neighboring Backend Wave Fit
No cross-table dependencies between nutrition and existing WAVE-01/02/03/04 tables. Nutrition tables stand alone.

## Frontend Pages Context
PAGE-007 requires: GET/POST/DELETE nutrition-products, GET/POST nutrition-templates/current, GET/POST daily-overrides, GET macro-calculations.

## Codebase Evidence
Existing graph schema uses union result types (SettingsResult, PinOperationResult) with success/error variants.

## Proposed Details

### Database Schema (Single Migration 00081)

**nutrition_product**
- id (UUID, PK, default gen_random_uuid())
- user_id (UUID, FK → atlas_users, NOT NULL)
- name (VARCHAR, NOT NULL)
- calories_per_100g (REAL, NOT NULL, >= 0)
- protein_per_100g (REAL, NOT NULL, >= 0)
- fat_per_100g (REAL, NOT NULL, >= 0)
- carbs_per_100g (REAL, NOT NULL, >= 0)
- notes (TEXT, nullable)
- is_active (BOOLEAN, NOT NULL, DEFAULT true) — soft-delete flag
- created_at (TIMESTAMPTZ, NOT NULL, DEFAULT now())
- updated_at (TIMESTAMPTZ, NOT NULL, DEFAULT now())
- Index: idx_nutrition_product_user (user_id)
- Unique: none (duplicate names allowed for flexibility)

**nutrition_template**
- id (UUID, PK, default gen_random_uuid())
- user_id (UUID, FK → atlas_users, NOT NULL)
- week_start_date (DATE, NOT NULL)
- title (VARCHAR, nullable)
- notes (TEXT, nullable)
- created_at (TIMESTAMPTZ, NOT NULL, DEFAULT now())
- updated_at (TIMESTAMPTZ, NOT NULL, DEFAULT now())
- Unique: (user_id, week_start_date) — one template per user per week (upsert semantics)
- Index: idx_nutrition_template_week (user_id, week_start_date)

**nutrition_template_item**
- id (UUID, PK, default gen_random_uuid())
- template_id (UUID, FK → nutrition_template ON DELETE CASCADE, NOT NULL)
- product_id (UUID, FK → nutrition_product, NOT NULL) — no CASCADE (soft-delete)
- amount_grams (REAL, NOT NULL, > 0)
- meal_label (VARCHAR, nullable)
- notes (TEXT, nullable)
- created_at (TIMESTAMPTZ, NOT NULL, DEFAULT now())
- updated_at (TIMESTAMPTZ, NOT NULL, DEFAULT now())
- Index: idx_nutrition_template_item_template (template_id)

**daily_nutrition_override**
- id (UUID, PK, default gen_random_uuid())
- user_id (UUID, FK → atlas_users, NOT NULL)
- date (DATE, NOT NULL)
- notes (TEXT, nullable)
- created_at (TIMESTAMPTZ, NOT NULL, DEFAULT now())
- updated_at (TIMESTAMPTZ, NOT NULL, DEFAULT now())
- Unique: (user_id, date) — one override per user per date
- Index: idx_nutrition_override_date (user_id, date)

**daily_nutrition_override_item**
- id (UUID, PK, default gen_random_uuid())
- override_id (UUID, FK → daily_nutrition_override ON DELETE CASCADE, NOT NULL)
- product_id (UUID, FK → nutrition_product, NOT NULL) — no CASCADE (soft-delete)
- amount_grams (REAL, NOT NULL, > 0)
- operation (VARCHAR, NOT NULL, CHECK IN 'add', 'subtract', 'replace')
- meal_label (VARCHAR, nullable)
- notes (TEXT, nullable)
- created_at (TIMESTAMPTZ, NOT NULL, DEFAULT now())
- updated_at (TIMESTAMPTZ, NOT NULL, DEFAULT now())
- Index: idx_nutrition_override_item_override (override_id)

### GraphQL Operations

**NutritionProduct**
- Query: nutritionProducts (list all active), nutritionProduct (by ID)
- Mutation: createNutritionProduct, updateNutritionProduct, deleteNutritionProduct (soft-delete)

**NutritionTemplate**
- Query: nutritionTemplates (date range), nutritionTemplate (by ID), nutritionTemplateCurrent (by weekStartDate)
- Mutation: createNutritionTemplate (upsert), updateNutritionTemplate (title/notes), deleteNutritionTemplate (cascade)

**NutritionTemplateItem**
- No top-level query — nested under template query
- Mutation: createNutritionTemplateItem, updateNutritionTemplateItem, deleteNutritionTemplateItem

**DailyNutritionOverride**
- Query: dailyNutritionOverrides (date range), dailyNutritionOverride (by ID), dailyNutritionOverrideByDate (by date)
- Mutation: createDailyNutritionOverride, updateDailyNutritionOverride, deleteDailyNutritionOverride (cascade)

**DailyNutritionOverrideItem**
- No top-level query — nested under override query
- Mutation: createDailyNutritionOverrideItem, updateDailyNutritionOverrideItem, deleteDailyNutritionOverrideItem

**Macro Calculation**
- Query: nutritionMacros(weekStartDate: Date!, date: Date) — returns calculated macros for one week, optionally filtered to one day
- Returns: NutritionMacros { date, calories, protein, fat, carbs }

### Union Result Pattern
Each mutation returns a union result following the existing pattern:
- NutritionProductResult: NutritionProduct | ValidationError | AuthError
- NutritionTemplateResult: NutritionTemplate | ValidationError | AuthError
- DailyNutritionOverrideResult: DailyNutritionOverride | ValidationError | AuthError
- NutritionMacrosResult: NutritionMacros | ValidationError | AuthError

### Log Markers
- [NutritionProduct][create|update|delete|get|list]
- [NutritionTemplate][create|update|delete|get|list|current]
- [NutritionTemplateItem][create|update|delete]
- [DailyNutritionOverride][create|update|delete|get|list]
- [DailyNutritionOverrideItem][create|update|delete]
- [NutritionMacros][calculate]
- Nutritional values (calories, protein, fat, carbs) are NOT sensitive — may be logged
- Product names and notes: log privacy — log entity IDs only, not descriptive text
- Operation enum values may be logged (non-sensitive)

### Operations
- PostgreSQL: goose migration 00081_nutrition_tables.sql, single migration, reversible
- Existing Docker Compose stack, no new services
- No media storage needed (no binary uploads)
- No external API calls (MVP constraint per RULE-029)
- No new config sections needed (reuses WAVE-01 config)

## Risks And Rollback
- Single migration (00081): simpler but may conflict with other wave migrations in parallel. If WAVE-04 uses 00081-00087, WAVE-05 must use 00088+.
- Template upsert is destructive — old template is replaced, not versioned. Consider warning before upsert.
- Macro calculation is O(n) over template items. MVP performance is acceptable.

## Questions Raised
- DQ-W05-005: Macro calculation — should it be a separate GraphQL query or computed inline in template/override queries? Recommended: separate query for flexibility. Frontend calls once on page load and after mutations.
- DQ-W05-006: Should NutritionProduct deletion be soft (isActive flag) or hard (blocked if referenced)? Recommend soft-delete. Historical template/override items referencing a deleted product should still return data (product name, macros) but product is excluded from new item selection.

## Traceability Candidates
- docs/product-verified/domain-model.md → DB tables
- docs/technical-verified/data-contracts.md → schema design
- docs/technical-verified/operations-observability.md → log markers
- docs/product-verified/business-rules.md RULE-010 → macro calculation