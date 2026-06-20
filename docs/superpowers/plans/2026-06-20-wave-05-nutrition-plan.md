# WAVE-05: Nutrition Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement Nutrition module with 5 entities (NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, DailyNutritionOverrideItem) + KJBJU macro calculation, all via GraphQL.

**Architecture:** Single goose migration for 5 tables. sqlc-generated CRUD queries. Repository adapters with sqlc-backed data access. Transport-neutral service layer with validation. GraphQL schema with union result types. Resolvers with PIN auth guard. Follows exact WAVE-04 cardio/body patterns.

**Tech Stack:** Go, sqlc, gqlgen, goose, pgx/v5, GraphQL, PostgreSQL

---

## File Structure

### New Files (18)
```
apps/api/internal/repository/postgres/migrations/00090_nutrition_tables.sql
apps/api/internal/repository/postgres/queries/nutrition_products.sql
apps/api/internal/repository/postgres/queries/nutrition_templates.sql
apps/api/internal/repository/postgres/queries/nutrition_template_items.sql
apps/api/internal/repository/postgres/queries/nutrition_overrides.sql
apps/api/internal/repository/postgres/queries/nutrition_override_items.sql
apps/api/internal/atlas/models/nutrition.go
apps/api/internal/atlas/repository/postgres/nutrition_product_repo.go
apps/api/internal/atlas/repository/postgres/nutrition_template_repo.go
apps/api/internal/atlas/repository/postgres/nutrition_template_item_repo.go
apps/api/internal/atlas/repository/postgres/nutrition_override_repo.go
apps/api/internal/atlas/repository/postgres/nutrition_override_item_repo.go
apps/api/internal/atlas/service/nutrition_product_service.go
apps/api/internal/atlas/service/nutrition_template_service.go
apps/api/internal/atlas/service/nutrition_template_item_service.go
apps/api/internal/atlas/service/nutrition_override_service.go
apps/api/internal/atlas/service/nutrition_macro_service.go
apps/api/internal/atlas/graph/schema/nutrition.graphql
apps/api/internal/atlas/graph/resolver/nutrition.go
```

### Modified Files (4)
```
apps/api/internal/atlas/graph/schema/schema.graphql   (+Query/Mutation fields)
apps/api/internal/atlas/graph/resolver/resolver.go     (+5 service fields)
apps/api/atlas-gqlgen.yml                              (+20 model bindings)
apps/api/cmd/server/main.go                            (+5 repos + 5 services + wiring)
```

---

### Task 1: Migration — `00090_nutrition_tables.sql`

**Files:**
- Create: `apps/api/internal/repository/postgres/migrations/00090_nutrition_tables.sql`

- [ ] **Step 1: Create migration file**

Write the migration file. Single migration creates all 5 tables with FKs, CHECK constraints, indexes, and unique constraints. Reversible.

```sql
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
```

- [ ] **Step 2: Verify migration syntax**

Run: `bunx nx run api:codegen` (sqlc will parse the migration)
Expected: sqlc generates query code without errors.

- [ ] **Step 3: Commit**

```bash
git add apps/api/internal/repository/postgres/migrations/00090_nutrition_tables.sql
git commit -m "feat(wave-05): add nutrition tables migration"
```

---

### Task 2: sqlc Query Files (5 files)

**Files:**
- Create: `apps/api/internal/repository/postgres/queries/nutrition_products.sql`
- Create: `apps/api/internal/repository/postgres/queries/nutrition_templates.sql`
- Create: `apps/api/internal/repository/postgres/queries/nutrition_template_items.sql`
- Create: `apps/api/internal/repository/postgres/queries/nutrition_overrides.sql`
- Create: `apps/api/internal/repository/postgres/queries/nutrition_override_items.sql`

- [ ] **Step 1: Create nutrition_products.sql**

```sql
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
```

- [ ] **Step 2: Create nutrition_templates.sql**

```sql
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
```

- [ ] **Step 3: Create nutrition_template_items.sql**

```sql
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
```

- [ ] **Step 4: Create nutrition_overrides.sql**

```sql
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
```

- [ ] **Step 5: Create nutrition_override_items.sql**

```sql
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
```

- [ ] **Step 6: Generate sqlc code**

Run: `bunx nx run api:codegen`
Expected: sqlc generates Go code in `apps/api/internal/repository/postgres/generated/` — new types for all 5 entities.

- [ ] **Step 7: Commit**

```bash
git add apps/api/internal/repository/postgres/queries/nutrition_*.sql
git commit -m "feat(wave-05): add nutrition sqlc query files"
```

---

### Task 3: Models — `models/nutrition.go`

**Files:**
- Create: `apps/api/internal/atlas/models/nutrition.go`

- [ ] **Step 1: Create models file**

Write all nutrition types following the cardio.go pattern: DB records, public models, inputs, result unions, error types, enums, converter functions, error-to-result helpers.

```go
package models

type Operation string

const (
	OperationAdd      Operation = "add"
	OperationSubtract Operation = "subtract"
	OperationReplace  Operation = "replace"
)

type NutritionErrorCode string

const (
	NutritionErrorValidation NutritionErrorCode = "VALIDATION_ERROR"
	NutritionErrorNotFound   NutritionErrorCode = "NOT_FOUND"
	NutritionErrorAuth       NutritionErrorCode = "AUTH_ERROR"
	NutritionErrorInternal   NutritionErrorCode = "INTERNAL_ERROR"
)

// DB Records
type NutritionProductRecord struct {
	ID              string
	UserID          string
	Name            string
	CaloriesPer100g float64
	ProteinPer100g  float64
	FatPer100g      float64
	CarbsPer100g    float64
	Notes           *string
	IsActive        bool
	CreatedAt       string
	UpdatedAt       string
}

type NutritionTemplateRecord struct {
	ID            string
	UserID        string
	WeekStartDate string
	Title         *string
	Notes         *string
	CreatedAt     string
	UpdatedAt     string
}

type NutritionTemplateItemRecord struct {
	ID          string
	TemplateID  string
	ProductID   string
	AmountGrams float64
	MealLabel   *string
	Notes       *string
	CreatedAt   string
	UpdatedAt   string
}

type DailyNutritionOverrideRecord struct {
	ID        string
	UserID    string
	Date      string
	Notes     *string
	CreatedAt string
	UpdatedAt string
}

type DailyNutritionOverrideItemRecord struct {
	ID          string
	OverrideID  string
	ProductID   string
	AmountGrams float64
	Operation   string
	MealLabel   *string
	Notes       *string
	CreatedAt   string
	UpdatedAt   string
}

// Public Models
type NutritionProduct struct {
	ID              string   `json:"id"`
	UserID          string   `json:"userId"`
	Name            string   `json:"name"`
	CaloriesPer100g float64  `json:"caloriesPer100g"`
	ProteinPer100g  float64  `json:"proteinPer100g"`
	FatPer100g      float64  `json:"fatPer100g"`
	CarbsPer100g    float64  `json:"carbsPer100g"`
	Notes           *string  `json:"notes"`
	IsActive        bool     `json:"isActive"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}

type NutritionTemplate struct {
	ID            string                   `json:"id"`
	UserID        string                   `json:"userId"`
	WeekStartDate string                   `json:"weekStartDate"`
	Title         *string                  `json:"title"`
	Notes         *string                  `json:"notes"`
	Items         []NutritionTemplateItem  `json:"items"`
	CreatedAt     string                   `json:"createdAt"`
	UpdatedAt     string                   `json:"updatedAt"`
}

type NutritionTemplateItem struct {
	ID          string   `json:"id"`
	TemplateID  string   `json:"templateId"`
	ProductID   string   `json:"productId"`
	AmountGrams float64  `json:"amountGrams"`
	MealLabel   *string  `json:"mealLabel"`
	Notes       *string  `json:"notes"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

type DailyNutritionOverride struct {
	ID        string                       `json:"id"`
	UserID    string                       `json:"userId"`
	Date      string                       `json:"date"`
	Notes     *string                      `json:"notes"`
	Items     []DailyNutritionOverrideItem `json:"items"`
	CreatedAt string                       `json:"createdAt"`
	UpdatedAt string                       `json:"updatedAt"`
}

type DailyNutritionOverrideItem struct {
	ID          string    `json:"id"`
	OverrideID  string    `json:"overrideId"`
	ProductID   string    `json:"productId"`
	AmountGrams float64   `json:"amountGrams"`
	Operation   Operation `json:"operation"`
	MealLabel   *string   `json:"mealLabel"`
	Notes       *string   `json:"notes"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}

type NutritionMacros struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
}

// Inputs
type CreateProductInput struct {
	Name            string   `json:"name"`
	CaloriesPer100g float64  `json:"caloriesPer100g"`
	ProteinPer100g  float64  `json:"proteinPer100g"`
	FatPer100g      float64  `json:"fatPer100g"`
	CarbsPer100g    float64  `json:"carbsPer100g"`
	Notes           *string  `json:"notes"`
}

type UpdateProductInput struct {
	Name            *string  `json:"name"`
	CaloriesPer100g *float64 `json:"caloriesPer100g"`
	ProteinPer100g  *float64 `json:"proteinPer100g"`
	FatPer100g      *float64 `json:"fatPer100g"`
	CarbsPer100g    *float64 `json:"carbsPer100g"`
	Notes           *string  `json:"notes"`
}

type CreateTemplateInput struct {
	WeekStartDate Date    `json:"weekStartDate"`
	Title         *string `json:"title"`
	Notes         *string `json:"notes"`
}

type UpdateTemplateInput struct {
	Title *string `json:"title"`
	Notes *string `json:"notes"`
}

type CreateTemplateItemInput struct {
	TemplateID  string   `json:"templateId"`
	ProductID   string   `json:"productId"`
	AmountGrams float64  `json:"amountGrams"`
	MealLabel   *string  `json:"mealLabel"`
	Notes       *string  `json:"notes"`
}

type UpdateTemplateItemInput struct {
	AmountGrams *float64 `json:"amountGrams"`
	MealLabel   *string  `json:"mealLabel"`
	Notes       *string  `json:"notes"`
}

type CreateOverrideInput struct {
	Date  Date     `json:"date"`
	Notes *string  `json:"notes"`
}

type UpdateOverrideInput struct {
	Notes *string `json:"notes"`
}

type CreateOverrideItemInput struct {
	OverrideID  string    `json:"overrideId"`
	ProductID   string    `json:"productId"`
	AmountGrams float64   `json:"amountGrams"`
	Operation   Operation `json:"operation"`
	MealLabel   *string   `json:"mealLabel"`
	Notes       *string   `json:"notes"`
}

type UpdateOverrideItemInput struct {
	AmountGrams *float64  `json:"amountGrams"`
	Operation   *Operation `json:"operation"`
	MealLabel   *string   `json:"mealLabel"`
	Notes       *string   `json:"notes"`
}

// Result Types
type NutritionProductResult struct {
	NutritionProduct *NutritionProduct      `json:"nutritionProduct"`
	ValidationErr    *NutritionValidationErr `json:"validationError"`
	NotFoundErr      *NutritionNotFoundErr   `json:"notFoundError"`
	AuthErr          *NutritionAuthErr       `json:"authError"`
}

type NutritionProductsResult struct {
	Products       []NutritionProduct      `json:"products"`
	ValidationErr  *NutritionValidationErr  `json:"validationError"`
	AuthErr        *NutritionAuthErr        `json:"authError"`
}

type NutritionTemplateResult struct {
	NutritionTemplate *NutritionTemplate     `json:"nutritionTemplate"`
	ValidationErr     *NutritionValidationErr `json:"validationError"`
	NotFoundErr       *NutritionNotFoundErr    `json:"notFoundError"`
	AuthErr           *NutritionAuthErr        `json:"authError"`
}

type NutritionTemplatesResult struct {
	Templates      []NutritionTemplate     `json:"templates"`
	ValidationErr  *NutritionValidationErr  `json:"validationError"`
	AuthErr        *NutritionAuthErr        `json:"authError"`
}

type NutritionTemplateItemResult struct {
	NutritionTemplateItem *NutritionTemplateItem  `json:"nutritionTemplateItem"`
	ValidationErr         *NutritionValidationErr  `json:"validationError"`
	NotFoundErr           *NutritionNotFoundErr    `json:"notFoundError"`
	AuthErr               *NutritionAuthErr        `json:"authError"`
}

type DailyNutritionOverrideResult struct {
	DailyNutritionOverride *DailyNutritionOverride `json:"dailyNutritionOverride"`
	ValidationErr          *NutritionValidationErr  `json:"validationError"`
	NotFoundErr            *NutritionNotFoundErr    `json:"notFoundError"`
	AuthErr                *NutritionAuthErr        `json:"authError"`
}

type DailyNutritionOverridesResult struct {
	Overrides      []DailyNutritionOverride `json:"overrides"`
	ValidationErr  *NutritionValidationErr   `json:"validationError"`
	AuthErr        *NutritionAuthErr         `json:"authError"`
}

type DailyNutritionOverrideItemResult struct {
	DailyNutritionOverrideItem *DailyNutritionOverrideItem `json:"dailyNutritionOverrideItem"`
	ValidationErr              *NutritionValidationErr      `json:"validationError"`
	NotFoundErr                *NutritionNotFoundErr        `json:"notFoundError"`
	AuthErr                    *NutritionAuthErr            `json:"authError"`
}

type NutritionMacrosResult struct {
	Macros         *NutritionMacros         `json:"macros"`
	ValidationErr  *NutritionValidationErr  `json:"validationError"`
	AuthErr        *NutritionAuthErr        `json:"authError"`
}

// Error Types
type NutritionValidationErr struct {
	Message string             `json:"message"`
	Code    NutritionErrorCode `json:"code"`
}

func (e *NutritionValidationErr) Error() string {
	if e == nil || e.Message == "" {
		return "nutrition validation error"
	}
	return e.Message
}

type NutritionNotFoundErr struct {
	Message string             `json:"message"`
	Code    NutritionErrorCode `json:"code"`
}

func (e *NutritionNotFoundErr) Error() string {
	if e == nil || e.Message == "" {
		return "nutrition not found"
	}
	return e.Message
}

type NutritionAuthErr struct {
	Message string             `json:"message"`
	Code    NutritionErrorCode `json:"code"`
}

func (e *NutritionAuthErr) Error() string {
	if e == nil || e.Message == "" {
		return "nutrition auth error"
	}
	return e.Message
}

// Converters
func NutritionProductFromRecord(r *NutritionProductRecord) *NutritionProduct {
	if r == nil {
		return nil
	}
	return &NutritionProduct{
		ID:              r.ID,
		UserID:          r.UserID,
		Name:            r.Name,
		CaloriesPer100g: r.CaloriesPer100g,
		ProteinPer100g:  r.ProteinPer100g,
		FatPer100g:      r.FatPer100g,
		CarbsPer100g:    r.CarbsPer100g,
		Notes:           r.Notes,
		IsActive:        r.IsActive,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

func NutritionTemplateFromRecord(r *NutritionTemplateRecord, items []NutritionTemplateItem) *NutritionTemplate {
	if r == nil {
		return nil
	}
	if items == nil {
		items = []NutritionTemplateItem{}
	}
	return &NutritionTemplate{
		ID:            r.ID,
		UserID:        r.UserID,
		WeekStartDate: r.WeekStartDate,
		Title:         r.Title,
		Notes:         r.Notes,
		Items:         items,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

func NutritionTemplateItemFromRecord(r *NutritionTemplateItemRecord) *NutritionTemplateItem {
	if r == nil {
		return nil
	}
	return &NutritionTemplateItem{
		ID:          r.ID,
		TemplateID:  r.TemplateID,
		ProductID:   r.ProductID,
		AmountGrams: r.AmountGrams,
		MealLabel:   r.MealLabel,
		Notes:       r.Notes,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func NutritionTemplateItemsFromRecords(records []NutritionTemplateItemRecord) []NutritionTemplateItem {
	out := make([]NutritionTemplateItem, len(records))
	for i := range records {
		out[i] = *NutritionTemplateItemFromRecord(&records[i])
	}
	return out
}

func DailyNutritionOverrideFromRecord(r *DailyNutritionOverrideRecord, items []DailyNutritionOverrideItem) *DailyNutritionOverride {
	if r == nil {
		return nil
	}
	if items == nil {
		items = []DailyNutritionOverrideItem{}
	}
	return &DailyNutritionOverride{
		ID:        r.ID,
		UserID:    r.UserID,
		Date:      r.Date,
		Notes:     r.Notes,
		Items:     items,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func DailyNutritionOverrideItemFromRecord(r *DailyNutritionOverrideItemRecord) *DailyNutritionOverrideItem {
	if r == nil {
		return nil
	}
	op := Operation(r.Operation)
	return &DailyNutritionOverrideItem{
		ID:          r.ID,
		OverrideID:  r.OverrideID,
		ProductID:   r.ProductID,
		AmountGrams: r.AmountGrams,
		Operation:   op,
		MealLabel:   r.MealLabel,
		Notes:       r.Notes,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func DailyNutritionOverrideItemsFromRecords(records []DailyNutritionOverrideItemRecord) []DailyNutritionOverrideItem {
	out := make([]DailyNutritionOverrideItem, len(records))
	for i := range records {
		out[i] = *DailyNutritionOverrideItemFromRecord(&records[i])
	}
	return out
}

// Error-to-result helpers
func NutritionProductResultFromError(err error) *NutritionProductResult {
	if err == nil {
		return nil
	}
	var validationErr *NutritionValidationErr
	if fmtErrorAs(err, &validationErr) {
		return &NutritionProductResult{ValidationErr: validationErr}
	}
	var notFoundErr *NutritionNotFoundErr
	if fmtErrorAs(err, &notFoundErr) {
		return &NutritionProductResult{NotFoundErr: notFoundErr}
	}
	var authErr *NutritionAuthErr
	if fmtErrorAs(err, &authErr) {
		return &NutritionProductResult{AuthErr: authErr}
	}
	return nil
}
```

Note: `fmtErrorAs` is already defined in `cardio.go` in the same package — it is a shared utility and does not need to be redefined.

- [ ] **Step 2: Commit**

```bash
git add apps/api/internal/atlas/models/nutrition.go
git commit -m "feat(wave-05): add nutrition domain models"
```

---

### Task 4: Repository Layer (5 files)

**Files:**
- Create: `apps/api/internal/atlas/repository/postgres/nutrition_product_repo.go`
- Create: `apps/api/internal/atlas/repository/postgres/nutrition_template_repo.go`
- Create: `apps/api/internal/atlas/repository/postgres/nutrition_template_item_repo.go`
- Create: `apps/api/internal/atlas/repository/postgres/nutrition_override_repo.go`
- Create: `apps/api/internal/atlas/repository/postgres/nutrition_override_item_repo.go`

Each repo follows the exact pattern from `cardio_entry_repo.go`:
- Interface + private struct + `New*Repository(pool *pgxpool.Pool)`
- Uses `uuidFromString` and `parseTwoUUIDs` from settings_repo.go
- Uses `nullableText` and `formatTimestamp` helpers from settings_repo.go
- Record converter function from generated row
- NOT FOUND → return `(nil, nil)`

- [ ] **Step 1: Create nutrition_product_repo.go**

```go
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/repository/postgres/generated"
)

type NutritionProductRepository interface {
	Create(ctx context.Context, userID string, name string, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g float64, notes *string) (*models.NutritionProductRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error)
	ListActive(ctx context.Context, userID string) ([]models.NutritionProductRecord, error)
	Update(ctx context.Context, userID string, id string, name string, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g float64, notes *string) (*models.NutritionProductRecord, error)
	SoftDelete(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error)
	GetByIDIncludeInactive(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error)
}

type nutritionProductRepository struct {
	q *generated.Queries
}

func NewNutritionProductRepository(pool *pgxpool.Pool) NutritionProductRepository {
	return &nutritionProductRepository{q: generated.New(pool)}
}

func (r *nutritionProductRepository) Create(ctx context.Context, userID string, name string, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g float64, notes *string) (*models.NutritionProductRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.Create: %w", err)
	}

	row, err := r.q.CreateNutritionProduct(ctx, generated.CreateNutritionProductParams{
		UserID:          uid,
		Name:            name,
		CaloriesPer100g: float32(caloriesPer100g),
		ProteinPer100g:  float32(proteinPer100g),
		FatPer100g:      float32(fatPer100g),
		CarbsPer100g:    float32(carbsPer100g),
		Notes:           nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.Create: %w", err)
	}

	return nutritionProductRecordFromRow(row), nil
}

func (r *nutritionProductRepository) GetByID(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
	uid, pid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.GetByID: %w", err)
	}

	row, err := r.q.GetNutritionProductByID(ctx, generated.GetNutritionProductByIDParams{ID: pid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_product_repo.GetByID: %w", err)
	}

	return nutritionProductRecordFromRow(row), nil
}

func (r *nutritionProductRepository) ListActive(ctx context.Context, userID string) ([]models.NutritionProductRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.ListActive: %w", err)
	}

	rows, err := r.q.ListActiveNutritionProducts(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.ListActive: %w", err)
	}

	out := make([]models.NutritionProductRecord, len(rows))
	for i, row := range rows {
		out[i] = *nutritionProductRecordFromRow(row)
	}
	return out, nil
}

func (r *nutritionProductRepository) Update(ctx context.Context, userID string, id string, name string, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g float64, notes *string) (*models.NutritionProductRecord, error) {
	uid, pid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.Update: %w", err)
	}

	row, err := r.q.UpdateNutritionProduct(ctx, generated.UpdateNutritionProductParams{
		ID:              pid,
		UserID:          uid,
		Name:            name,
		CaloriesPer100g: float32(caloriesPer100g),
		ProteinPer100g:  float32(proteinPer100g),
		FatPer100g:      float32(fatPer100g),
		CarbsPer100g:    float32(carbsPer100g),
		Notes:           nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_product_repo.Update: %w", err)
	}

	return nutritionProductRecordFromRow(row), nil
}

func (r *nutritionProductRepository) SoftDelete(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
	uid, pid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.SoftDelete: %w", err)
	}

	row, err := r.q.SoftDeleteNutritionProduct(ctx, generated.SoftDeleteNutritionProductParams{ID: pid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_product_repo.SoftDelete: %w", err)
	}

	return nutritionProductRecordFromRow(row), nil
}

func (r *nutritionProductRepository) GetByIDIncludeInactive(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
	uid, pid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.GetByIDIncludeInactive: %w", err)
	}

	row, err := r.q.GetNutritionProductByIDIncludeInactive(ctx, generated.GetNutritionProductByIDIncludeInactiveParams{ID: pid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_product_repo.GetByIDIncludeInactive: %w", err)
	}

	return nutritionProductRecordFromRow(row), nil
}

func nutritionProductRecordFromRow(row generated.NutritionProduct) *models.NutritionProductRecord {
	return &models.NutritionProductRecord{
		ID:              row.ID.String(),
		UserID:          row.UserID.String(),
		Name:            row.Name,
		CaloriesPer100g: float64(row.CaloriesPer100g),
		ProteinPer100g:  float64(row.ProteinPer100g),
		FatPer100g:      float64(row.FatPer100g),
		CarbsPer100g:    float64(row.CarbsPer100g),
		Notes:           textPtr(row.Notes),
		IsActive:        row.IsActive,
		CreatedAt:       formatTimestamp(row.CreatedAt),
		UpdatedAt:       formatTimestamp(row.UpdatedAt),
	}
}
```

- [ ] **Step 2: Create nutrition_template_repo.go**

```go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/repository/postgres/generated"
)

type NutritionTemplateRepository interface {
	Upsert(ctx context.Context, userID string, weekStartDate string, title, notes *string) (*models.NutritionTemplateRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error)
	GetByWeek(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error)
	ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.NutritionTemplateRecord, error)
	Update(ctx context.Context, userID string, id string, title, notes *string) (*models.NutritionTemplateRecord, error)
	Delete(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error)
}

type nutritionTemplateRepository struct {
	q *generated.Queries
}

func NewNutritionTemplateRepository(pool *pgxpool.Pool) NutritionTemplateRepository {
	return &nutritionTemplateRepository{q: generated.New(pool)}
}

func (r *nutritionTemplateRepository) Upsert(ctx context.Context, userID string, weekStartDate string, title, notes *string) (*models.NutritionTemplateRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.Upsert: %w", err)
	}

	wd, err := parseDate(weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.Upsert: %w", err)
	}

	row, err := r.q.UpsertNutritionTemplate(ctx, generated.UpsertNutritionTemplateParams{
		UserID:        uid,
		WeekStartDate: wd,
		Title:         nullableText(title),
		Notes:         nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.Upsert: %w", err)
	}

	return nutritionTemplateRecordFromRow(row), nil
}

func (r *nutritionTemplateRepository) GetByID(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
	uid, tid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.GetByID: %w", err)
	}

	row, err := r.q.GetNutritionTemplateByID(ctx, generated.GetNutritionTemplateByIDParams{ID: tid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_repo.GetByID: %w", err)
	}

	return nutritionTemplateRecordFromRow(row), nil
}

func (r *nutritionTemplateRepository) GetByWeek(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.GetByWeek: %w", err)
	}

	wd, err := parseDate(weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.GetByWeek: %w", err)
	}

	row, err := r.q.GetNutritionTemplateByWeek(ctx, generated.GetNutritionTemplateByWeekParams{
		UserID:        uid,
		WeekStartDate: wd,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_repo.GetByWeek: %w", err)
	}

	return nutritionTemplateRecordFromRow(row), nil
}

func (r *nutritionTemplateRepository) ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.NutritionTemplateRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.ListByRange: %w", err)
	}

	sd, err := parseDate(startDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.ListByRange: %w", err)
	}

	ed, err := parseDate(endDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.ListByRange: %w", err)
	}

	rows, err := r.q.ListNutritionTemplatesByRange(ctx, generated.ListNutritionTemplatesByRangeParams{
		UserID:       uid,
		WeekStartDate: sd,
		WeekStartDate_2: ed,
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.ListByRange: %w", err)
	}

	out := make([]models.NutritionTemplateRecord, len(rows))
	for i, row := range rows {
		out[i] = *nutritionTemplateRecordFromRow(row)
	}
	return out, nil
}

func (r *nutritionTemplateRepository) Update(ctx context.Context, userID string, id string, title, notes *string) (*models.NutritionTemplateRecord, error) {
	uid, tid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.Update: %w", err)
	}

	row, err := r.q.UpdateNutritionTemplate(ctx, generated.UpdateNutritionTemplateParams{
		ID:    tid,
		UserID: uid,
		Title: nullableText(title),
		Notes: nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_repo.Update: %w", err)
	}

	return nutritionTemplateRecordFromRow(row), nil
}

func (r *nutritionTemplateRepository) Delete(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
	uid, tid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteNutritionTemplate(ctx, generated.DeleteNutritionTemplateParams{ID: tid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_repo.Delete: %w", err)
	}

	return nutritionTemplateRecordFromRow(row), nil
}

func nutritionTemplateRecordFromRow(row generated.NutritionTemplate) *models.NutritionTemplateRecord {
	return &models.NutritionTemplateRecord{
		ID:            row.ID.String(),
		UserID:        row.UserID.String(),
		WeekStartDate: dateFromPGDatePG(row.WeekStartDate),
		Title:         textPtr(row.Title),
		Notes:         textPtr(row.Notes),
		CreatedAt:     formatTimestamp(row.CreatedAt),
		UpdatedAt:     formatTimestamp(row.UpdatedAt),
	}
}

func parseDate(value string) (pgtype.Date, error) {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("parse date: %w", err)
	}
	return pgtype.Date{Time: t, Valid: true}, nil
}

func dateFromPGDatePG(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("2006-01-02")
}
```

- [ ] **Step 3: Create nutrition_template_item_repo.go**

```go
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/repository/postgres/generated"
)

type NutritionTemplateItemRepository interface {
	Create(ctx context.Context, templateID string, productID string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error)
	GetByID(ctx context.Context, id string) (*models.NutritionTemplateItemRecord, error)
	ListByTemplate(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error)
	Update(ctx context.Context, id string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error)
	Delete(ctx context.Context, id string) (*models.NutritionTemplateItemRecord, error)
}

type nutritionTemplateItemRepository struct {
	q *generated.Queries
}

func NewNutritionTemplateItemRepository(pool *pgxpool.Pool) NutritionTemplateItemRepository {
	return &nutritionTemplateItemRepository{q: generated.New(pool)}
}

func (r *nutritionTemplateItemRepository) Create(ctx context.Context, templateID string, productID string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error) {
	tid, err := uuidFromString(templateID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.Create: %w", err)
	}
	pid, err := uuidFromString(productID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.Create: %w", err)
	}

	row, err := r.q.CreateNutritionTemplateItem(ctx, generated.CreateNutritionTemplateItemParams{
		TemplateID:  tid,
		ProductID:   pid,
		AmountGrams: float32(amountGrams),
		MealLabel:   nullableText(mealLabel),
		Notes:       nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.Create: %w", err)
	}

	return nutritionTemplateItemRecordFromRow(row), nil
}

func (r *nutritionTemplateItemRepository) GetByID(ctx context.Context, id string) (*models.NutritionTemplateItemRecord, error) {
	iid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.GetByID: %w", err)
	}

	row, err := r.q.GetNutritionTemplateItemByID(ctx, iid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_item_repo.GetByID: %w", err)
	}

	return nutritionTemplateItemRecordFromRow(row), nil
}

func (r *nutritionTemplateItemRepository) ListByTemplate(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
	tid, err := uuidFromString(templateID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.ListByTemplate: %w", err)
	}

	rows, err := r.q.ListNutritionTemplateItemsByTemplate(ctx, tid)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.ListByTemplate: %w", err)
	}

	out := make([]models.NutritionTemplateItemRecord, len(rows))
	for i, row := range rows {
		out[i] = *nutritionTemplateItemRecordFromRow(row)
	}
	return out, nil
}

func (r *nutritionTemplateItemRepository) Update(ctx context.Context, id string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error) {
	iid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.Update: %w", err)
	}

	row, err := r.q.UpdateNutritionTemplateItem(ctx, generated.UpdateNutritionTemplateItemParams{
		ID:          iid,
		AmountGrams: float32(amountGrams),
		MealLabel:   nullableText(mealLabel),
		Notes:       nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_item_repo.Update: %w", err)
	}

	return nutritionTemplateItemRecordFromRow(row), nil
}

func (r *nutritionTemplateItemRepository) Delete(ctx context.Context, id string) (*models.NutritionTemplateItemRecord, error) {
	iid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteNutritionTemplateItem(ctx, iid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_item_repo.Delete: %w", err)
	}

	return nutritionTemplateItemRecordFromRow(row), nil
}

func nutritionTemplateItemRecordFromRow(row generated.NutritionTemplateItem) *models.NutritionTemplateItemRecord {
	return &models.NutritionTemplateItemRecord{
		ID:          row.ID.String(),
		TemplateID:  row.TemplateID.String(),
		ProductID:   row.ProductID.String(),
		AmountGrams: float64(row.AmountGrams),
		MealLabel:   textPtr(row.MealLabel),
		Notes:       textPtr(row.Notes),
		CreatedAt:   formatTimestamp(row.CreatedAt),
		UpdatedAt:   formatTimestamp(row.UpdatedAt),
	}
}
```

- [ ] **Step 4: Create nutrition_override_repo.go**

```go
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/repository/postgres/generated"
)

type DailyNutritionOverrideRepository interface {
	Upsert(ctx context.Context, userID string, date string, notes *string) (*models.DailyNutritionOverrideRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error)
	GetByDate(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error)
	ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.DailyNutritionOverrideRecord, error)
	Update(ctx context.Context, userID string, id string, notes *string) (*models.DailyNutritionOverrideRecord, error)
	Delete(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error)
}

type dailyNutritionOverrideRepository struct {
	q *generated.Queries
}

func NewDailyNutritionOverrideRepository(pool *pgxpool.Pool) DailyNutritionOverrideRepository {
	return &dailyNutritionOverrideRepository{q: generated.New(pool)}
}

func (r *dailyNutritionOverrideRepository) Upsert(ctx context.Context, userID string, date string, notes *string) (*models.DailyNutritionOverrideRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.Upsert: %w", err)
	}

	d, err := parseDate(date)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.Upsert: %w", err)
	}

	row, err := r.q.UpsertDailyNutritionOverride(ctx, generated.UpsertDailyNutritionOverrideParams{
		UserID: uid,
		Date:   d,
		Notes:  nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.Upsert: %w", err)
	}

	return dailyNutritionOverrideRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideRepository) GetByID(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error) {
	uid, oid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.GetByID: %w", err)
	}

	row, err := r.q.GetDailyNutritionOverrideByID(ctx, generated.GetDailyNutritionOverrideByIDParams{ID: oid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_repo.GetByID: %w", err)
	}

	return dailyNutritionOverrideRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideRepository) GetByDate(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.GetByDate: %w", err)
	}

	d, err := parseDate(date)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.GetByDate: %w", err)
	}

	row, err := r.q.GetDailyNutritionOverrideByDate(ctx, generated.GetDailyNutritionOverrideByDateParams{
		UserID: uid,
		Date:   d,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_repo.GetByDate: %w", err)
	}

	return dailyNutritionOverrideRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideRepository) ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.DailyNutritionOverrideRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.ListByRange: %w", err)
	}

	sd, err := parseDate(startDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.ListByRange: %w", err)
	}

	ed, err := parseDate(endDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.ListByRange: %w", err)
	}

	rows, err := r.q.ListDailyNutritionOverridesByRange(ctx, generated.ListDailyNutritionOverridesByRangeParams{
		UserID: uid,
		Date:   sd,
		Date_2: ed,
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.ListByRange: %w", err)
	}

	out := make([]models.DailyNutritionOverrideRecord, len(rows))
	for i, row := range rows {
		out[i] = *dailyNutritionOverrideRecordFromRow(row)
	}
	return out, nil
}

func (r *dailyNutritionOverrideRepository) Update(ctx context.Context, userID string, id string, notes *string) (*models.DailyNutritionOverrideRecord, error) {
	uid, oid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.Update: %w", err)
	}

	row, err := r.q.UpdateDailyNutritionOverride(ctx, generated.UpdateDailyNutritionOverrideParams{
		ID:    oid,
		UserID: uid,
		Notes: nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_repo.Update: %w", err)
	}

	return dailyNutritionOverrideRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideRepository) Delete(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error) {
	uid, oid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteDailyNutritionOverride(ctx, generated.DeleteDailyNutritionOverrideParams{ID: oid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_repo.Delete: %w", err)
	}

	return dailyNutritionOverrideRecordFromRow(row), nil
}

func dailyNutritionOverrideRecordFromRow(row generated.DailyNutritionOverride) *models.DailyNutritionOverrideRecord {
	return &models.DailyNutritionOverrideRecord{
		ID:        row.ID.String(),
		UserID:    row.UserID.String(),
		Date:      dateFromPGDatePG(row.Date),
		Notes:     textPtr(row.Notes),
		CreatedAt: formatTimestamp(row.CreatedAt),
		UpdatedAt: formatTimestamp(row.UpdatedAt),
	}
}
```

- [ ] **Step 5: Create nutrition_override_item_repo.go**

```go
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/repository/postgres/generated"
)

type DailyNutritionOverrideItemRepository interface {
	Create(ctx context.Context, overrideID string, productID string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error)
	GetByID(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error)
	ListByOverride(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error)
	Update(ctx context.Context, id string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error)
	Delete(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error)
}

type dailyNutritionOverrideItemRepository struct {
	q *generated.Queries
}

func NewDailyNutritionOverrideItemRepository(pool *pgxpool.Pool) DailyNutritionOverrideItemRepository {
	return &dailyNutritionOverrideItemRepository{q: generated.New(pool)}
}

func (r *dailyNutritionOverrideItemRepository) Create(ctx context.Context, overrideID string, productID string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error) {
	oid, err := uuidFromString(overrideID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.Create: %w", err)
	}
	pid, err := uuidFromString(productID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.Create: %w", err)
	}

	row, err := r.q.CreateDailyNutritionOverrideItem(ctx, generated.CreateDailyNutritionOverrideItemParams{
		OverrideID:  oid,
		ProductID:   pid,
		AmountGrams: float32(amountGrams),
		Operation:   operation,
		MealLabel:   nullableText(mealLabel),
		Notes:       nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.Create: %w", err)
	}

	return dailyNutritionOverrideItemRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideItemRepository) GetByID(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error) {
	iid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.GetByID: %w", err)
	}

	row, err := r.q.GetDailyNutritionOverrideItemByID(ctx, iid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_item_repo.GetByID: %w", err)
	}

	return dailyNutritionOverrideItemRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideItemRepository) ListByOverride(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
	oid, err := uuidFromString(overrideID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.ListByOverride: %w", err)
	}

	rows, err := r.q.ListDailyNutritionOverrideItemsByOverride(ctx, oid)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.ListByOverride: %w", err)
	}

	out := make([]models.DailyNutritionOverrideItemRecord, len(rows))
	for i, row := range rows {
		out[i] = *dailyNutritionOverrideItemRecordFromRow(row)
	}
	return out, nil
}

func (r *dailyNutritionOverrideItemRepository) Update(ctx context.Context, id string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error) {
	iid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.Update: %w", err)
	}

	row, err := r.q.UpdateDailyNutritionOverrideItem(ctx, generated.UpdateDailyNutritionOverrideItemParams{
		ID:          iid,
		AmountGrams: float32(amountGrams),
		Operation:   operation,
		MealLabel:   nullableText(mealLabel),
		Notes:       nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_item_repo.Update: %w", err)
	}

	return dailyNutritionOverrideItemRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideItemRepository) Delete(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error) {
	iid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteDailyNutritionOverrideItem(ctx, iid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_item_repo.Delete: %w", err)
	}

	return dailyNutritionOverrideItemRecordFromRow(row), nil
}

func dailyNutritionOverrideItemRecordFromRow(row generated.DailyNutritionOverrideItem) *models.DailyNutritionOverrideItemRecord {
	return &models.DailyNutritionOverrideItemRecord{
		ID:          row.ID.String(),
		OverrideID:  row.OverrideID.String(),
		ProductID:   row.ProductID.String(),
		AmountGrams: float64(row.AmountGrams),
		Operation:   row.Operation,
		MealLabel:   textPtr(row.MealLabel),
		Notes:       textPtr(row.Notes),
		CreatedAt:   formatTimestamp(row.CreatedAt),
		UpdatedAt:   formatTimestamp(row.UpdatedAt),
	}
}
```

- [ ] **Step 6: Run sqlc codegen to verify all queries compile**

Run: `bunx nx run api:codegen`
Expected: sqlc generates code without errors, all 5 new query files produce Go types.

- [ ] **Step 7: Commit**

```bash
git add apps/api/internal/atlas/repository/postgres/nutrition_*.go
git commit -m "feat(wave-05): add nutrition repository layer"
```

---

### Task 5: Service Layer (5 files)

**Files:**
- Create: `apps/api/internal/atlas/service/nutrition_product_service.go`
- Create: `apps/api/internal/atlas/service/nutrition_template_service.go`
- Create: `apps/api/internal/atlas/service/nutrition_template_item_service.go`
- Create: `apps/api/internal/atlas/service/nutrition_override_service.go`
- Create: `apps/api/internal/atlas/service/nutrition_macro_service.go`

Each service follows the cardio.go pattern: interface + struct + constructor.

**Log markers:** Every public method in each service must emit a log marker via `s.logger.Info("[NutritionX][action]")` following the spec's log marker definitions. This enables traceability and verification-plan compliance. All services receive a `*zap.Logger` via constructor. Example:
- Create methods: `s.logger.Info("[NutritionProduct][create]")`
- Delete methods: `s.logger.Info("[NutritionProduct][delete]")`
- GetByID/ListActive: `s.logger.Info("[NutritionProduct][get]")` / `s.logger.Info("[NutritionProduct][list]")`

- [ ] **Step 1: Create nutrition_product_service.go**

```go
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
)

var (
	ErrProductNameRequired  = errors.New("product name is required")
	ErrProductMacroNegative = errors.New("nutritional values must be >= 0")
	ErrProductNotFound      = errors.New("nutrition product not found")
	ErrProductNameTooLong   = errors.New("product name must not exceed 255 characters")
)

type NutritionProductService interface {
	Create(ctx context.Context, userID string, input models.CreateProductInput) (*models.NutritionProduct, error)
	GetByID(ctx context.Context, userID string, id string) (*models.NutritionProduct, error)
	ListActive(ctx context.Context, userID string) ([]models.NutritionProduct, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateProductInput) (*models.NutritionProduct, error)
	Delete(ctx context.Context, userID string, id string) (*models.NutritionProduct, error)
}

type nutritionProductService struct {
	repo   postgres.NutritionProductRepository
	logger *zap.Logger
}

func NewNutritionProductService(repo postgres.NutritionProductRepository, logger *zap.Logger) NutritionProductService {
	return &nutritionProductService{repo: repo, logger: logger}
}

func (s *nutritionProductService) Create(ctx context.Context, userID string, input models.CreateProductInput) (*models.NutritionProduct, error) {
	s.logger.Info("[NutritionProduct][create]")
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrProductNameRequired
	}
	if len(name) > 255 {
		return nil, ErrProductNameTooLong
	}
	if input.CaloriesPer100g < 0 || input.ProteinPer100g < 0 || input.FatPer100g < 0 || input.CarbsPer100g < 0 {
		return nil, ErrProductMacroNegative
	}

	record, err := s.repo.Create(ctx, userID, name, input.CaloriesPer100g, input.ProteinPer100g, input.FatPer100g, input.CarbsPer100g, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_service.Create: %w", err)
	}

	return models.NutritionProductFromRecord(record), nil
}

func (s *nutritionProductService) GetByID(ctx context.Context, userID string, id string) (*models.NutritionProduct, error) {
	record, err := s.repo.GetByIDIncludeInactive(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrProductNotFound
	}
	return models.NutritionProductFromRecord(record), nil
}

func (s *nutritionProductService) ListActive(ctx context.Context, userID string) ([]models.NutritionProduct, error) {
	records, err := s.repo.ListActive(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_service.ListActive: %w", err)
	}

	out := make([]models.NutritionProduct, len(records))
	for i := range records {
		out[i] = *models.NutritionProductFromRecord(&records[i])
	}
	return out, nil
}

func (s *nutritionProductService) Update(ctx context.Context, userID string, id string, input models.UpdateProductInput) (*models.NutritionProduct, error) {
	existing, err := s.repo.GetByIDIncludeInactive(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_service.Update: %w", err)
	}
	if existing == nil {
		return nil, ErrProductNotFound
	}

	name := existing.Name
	if input.Name != nil {
		name = strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, ErrProductNameRequired
		}
		if len(name) > 255 {
			return nil, ErrProductNameTooLong
		}
	}

	calories := existing.CaloriesPer100g
	if input.CaloriesPer100g != nil {
		calories = *input.CaloriesPer100g
	}
	protein := existing.ProteinPer100g
	if input.ProteinPer100g != nil {
		protein = *input.ProteinPer100g
	}
	fat := existing.FatPer100g
	if input.FatPer100g != nil {
		fat = *input.FatPer100g
	}
	carbs := existing.CarbsPer100g
	if input.CarbsPer100g != nil {
		carbs = *input.CarbsPer100g
	}

	if calories < 0 || protein < 0 || fat < 0 || carbs < 0 {
		return nil, ErrProductMacroNegative
	}

	notes := input.Notes
	if notes == nil {
		notes = existing.Notes
	}

	record, err := s.repo.Update(ctx, userID, id, name, calories, protein, fat, carbs, notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrProductNotFound
	}

	return models.NutritionProductFromRecord(record), nil
}

func (s *nutritionProductService) Delete(ctx context.Context, userID string, id string) (*models.NutritionProduct, error) {
	record, err := s.repo.SoftDelete(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrProductNotFound
	}
	return models.NutritionProductFromRecord(record), nil
}
```

- [ ] **Step 2: Create nutrition_template_service.go**

```go
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
)

var (
	ErrTemplateWeekRequired = errors.New("weekStartDate is required")
	ErrTemplateNotFound     = errors.New("nutrition template not found")
	ErrTemplateItemNotFound = errors.New("nutrition template item not found")
)

type NutritionTemplateService interface {
	Create(ctx context.Context, userID string, input models.CreateTemplateInput) (*models.NutritionTemplate, error)
	GetByID(ctx context.Context, userID string, id string) (*models.NutritionTemplate, error)
	GetCurrent(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplate, error)
	ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.NutritionTemplate, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateTemplateInput) (*models.NutritionTemplate, error)
	Delete(ctx context.Context, userID string, id string) (*models.NutritionTemplate, error)
}

type nutritionTemplateService struct {
	repo     postgres.NutritionTemplateRepository
	itemRepo postgres.NutritionTemplateItemRepository
	logger   *zap.Logger
}

func NewNutritionTemplateService(repo postgres.NutritionTemplateRepository, itemRepo postgres.NutritionTemplateItemRepository, logger *zap.Logger) NutritionTemplateService {
	return &nutritionTemplateService{repo: repo, itemRepo: itemRepo, logger: logger}
}

func (s *nutritionTemplateService) loadItems(ctx context.Context, templateID string) []models.NutritionTemplateItem {
	records, err := s.itemRepo.ListByTemplate(ctx, templateID)
	if err != nil {
		return []models.NutritionTemplateItem{}
	}
	return models.NutritionTemplateItemsFromRecords(records)
}

func (s *nutritionTemplateService) Create(ctx context.Context, userID string, input models.CreateTemplateInput) (*models.NutritionTemplate, error) {
	wd := strings.TrimSpace(input.WeekStartDate.String())
	if wd == "" {
		return nil, ErrTemplateWeekRequired
	}

	if _, err := time.Parse("2006-01-02", wd); err != nil {
		return nil, fmt.Errorf("%w: invalid date format", ErrTemplateWeekRequired)
	}

	record, err := s.repo.Upsert(ctx, userID, wd, input.Title, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.Create: %w", err)
	}

	return models.NutritionTemplateFromRecord(record, nil), nil
}

func (s *nutritionTemplateService) GetByID(ctx context.Context, userID string, id string) (*models.NutritionTemplate, error) {
	record, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrTemplateNotFound
	}

	items := s.loadItems(ctx, record.ID)
	return models.NutritionTemplateFromRecord(record, items), nil
}

func (s *nutritionTemplateService) GetCurrent(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplate, error) {
	record, err := s.repo.GetByWeek(ctx, userID, weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.GetCurrent: %w", err)
	}
	if record == nil {
		return nil, nil
	}

	items := s.loadItems(ctx, record.ID)
	return models.NutritionTemplateFromRecord(record, items), nil
}

func (s *nutritionTemplateService) ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.NutritionTemplate, error) {
	records, err := s.repo.ListByRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.ListByRange: %w", err)
	}

	out := make([]models.NutritionTemplate, len(records))
	for i := range records {
		items := s.loadItems(ctx, records[i].ID)
		out[i] = *models.NutritionTemplateFromRecord(&records[i], items)
	}
	return out, nil
}

func (s *nutritionTemplateService) Update(ctx context.Context, userID string, id string, input models.UpdateTemplateInput) (*models.NutritionTemplate, error) {
	existing, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.Update: %w", err)
	}
	if existing == nil {
		return nil, ErrTemplateNotFound
	}

	title := input.Title
	if title == nil {
		title = existing.Title
	}
	notes := input.Notes
	if notes == nil {
		notes = existing.Notes
	}

	record, err := s.repo.Update(ctx, userID, id, title, notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrTemplateNotFound
	}

	items := s.loadItems(ctx, record.ID)
	return models.NutritionTemplateFromRecord(record, items), nil
}

func (s *nutritionTemplateService) Delete(ctx context.Context, userID string, id string) (*models.NutritionTemplate, error) {
	record, err := s.repo.Delete(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrTemplateNotFound
	}
	return models.NutritionTemplateFromRecord(record, nil), nil
}
```

- [ ] **Step 3: Create nutrition_template_item_service.go**

```go
package service

import (
	"context"
	"errors"
	"fmt"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
)

var (
	ErrTemplateItemAmountInvalid = errors.New("amountGrams must be greater than 0")
)

type NutritionTemplateItemService interface {
	Create(ctx context.Context, userID string, input models.CreateTemplateItemInput) (*models.NutritionTemplateItem, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateTemplateItemInput) (*models.NutritionTemplateItem, error)
	Delete(ctx context.Context, userID string, id string) (*models.NutritionTemplateItem, error)
}

type nutritionTemplateItemService struct {
	itemRepo postgres.NutritionTemplateItemRepository
	tmplRepo postgres.NutritionTemplateRepository
}

func NewNutritionTemplateItemService(itemRepo postgres.NutritionTemplateItemRepository, tmplRepo postgres.NutritionTemplateRepository) NutritionTemplateItemService {
	return &nutritionTemplateItemService{itemRepo: itemRepo, tmplRepo: tmplRepo}
}

func (s *nutritionTemplateItemService) Create(ctx context.Context, userID string, input models.CreateTemplateItemInput) (*models.NutritionTemplateItem, error) {
	if input.AmountGrams <= 0 {
		return nil, ErrTemplateItemAmountInvalid
	}

	// Verify template exists and belongs to user
	tmpl, err := s.tmplRepo.GetByID(ctx, userID, input.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_service.Create: %w", err)
	}
	if tmpl == nil {
		return nil, ErrTemplateNotFound
	}

	record, err := s.itemRepo.Create(ctx, input.TemplateID, input.ProductID, input.AmountGrams, input.MealLabel, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_service.Create: %w", err)
	}

	return models.NutritionTemplateItemFromRecord(record), nil
}

func (s *nutritionTemplateItemService) Update(ctx context.Context, userID string, id string, input models.UpdateTemplateItemInput) (*models.NutritionTemplateItem, error) {
	existing, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_service.Update: %w", err)
	}
	if existing == nil {
		return nil, ErrTemplateItemNotFound
	}

	amount := existing.AmountGrams
	if input.AmountGrams != nil {
		if *input.AmountGrams <= 0 {
			return nil, ErrTemplateItemAmountInvalid
		}
		amount = *input.AmountGrams
	}

	mealLabel := input.MealLabel
	if mealLabel == nil {
		mealLabel = existing.MealLabel
	}

	notes := input.Notes
	if notes == nil {
		notes = existing.Notes
	}

	record, err := s.itemRepo.Update(ctx, id, amount, mealLabel, notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrTemplateItemNotFound
	}

	return models.NutritionTemplateItemFromRecord(record), nil
}

func (s *nutritionTemplateItemService) Delete(ctx context.Context, userID string, id string) (*models.NutritionTemplateItem, error) {
	record, err := s.itemRepo.Delete(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrTemplateItemNotFound
	}
	return models.NutritionTemplateItemFromRecord(record), nil
}
```

- [ ] **Step 4: Create nutrition_override_service.go**

```go
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
)

var (
	ErrOverrideDateRequired = errors.New("date is required")
	ErrOverrideNotFound     = errors.New("daily nutrition override not found")
	ErrOverrideItemNotFound = errors.New("daily nutrition override item not found")
	ErrOverrideItemAmountInvalid = errors.New("amountGrams must be greater than 0")
	ErrOverrideItemOperationInvalid = errors.New("operation must be add, subtract, or replace")
)

type DailyNutritionOverrideService interface {
	Create(ctx context.Context, userID string, input models.CreateOverrideInput) (*models.DailyNutritionOverride, error)
	GetByID(ctx context.Context, userID string, id string) (*models.DailyNutritionOverride, error)
	GetByDate(ctx context.Context, userID string, date string) (*models.DailyNutritionOverride, error)
	ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.DailyNutritionOverride, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateOverrideInput) (*models.DailyNutritionOverride, error)
	Delete(ctx context.Context, userID string, id string) (*models.DailyNutritionOverride, error)
	CreateItem(ctx context.Context, userID string, input models.CreateOverrideItemInput) (*models.DailyNutritionOverrideItem, error)
	UpdateItem(ctx context.Context, userID string, itemID string, input models.UpdateOverrideItemInput) (*models.DailyNutritionOverrideItem, error)
	DeleteItem(ctx context.Context, userID string, itemID string) (*models.DailyNutritionOverrideItem, error)
}

type dailyNutritionOverrideService struct {
	repo     postgres.DailyNutritionOverrideRepository
	itemRepo postgres.DailyNutritionOverrideItemRepository
}

func NewNutritionOverrideService(repo postgres.DailyNutritionOverrideRepository, itemRepo postgres.DailyNutritionOverrideItemRepository) DailyNutritionOverrideService {
	return &dailyNutritionOverrideService{repo: repo, itemRepo: itemRepo}
}

func (s *dailyNutritionOverrideService) loadItems(ctx context.Context, overrideID string) []models.DailyNutritionOverrideItem {
	records, err := s.itemRepo.ListByOverride(ctx, overrideID)
	if err != nil {
		return []models.DailyNutritionOverrideItem{}
	}
	return models.DailyNutritionOverrideItemsFromRecords(records)
}

func (s *dailyNutritionOverrideService) Create(ctx context.Context, userID string, input models.CreateOverrideInput) (*models.DailyNutritionOverride, error) {
	d := strings.TrimSpace(input.Date.String())
	if d == "" {
		return nil, ErrOverrideDateRequired
	}

	if _, err := time.Parse("2006-01-02", d); err != nil {
		return nil, fmt.Errorf("%w: invalid date", ErrOverrideDateRequired)
	}

	record, err := s.repo.Upsert(ctx, userID, d, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.Create: %w", err)
	}

	return models.DailyNutritionOverrideFromRecord(record, nil), nil
}

func (s *dailyNutritionOverrideService) GetByID(ctx context.Context, userID string, id string) (*models.DailyNutritionOverride, error) {
	record, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrOverrideNotFound
	}

	items := s.loadItems(ctx, record.ID)
	return models.DailyNutritionOverrideFromRecord(record, items), nil
}

func (s *dailyNutritionOverrideService) GetByDate(ctx context.Context, userID string, date string) (*models.DailyNutritionOverride, error) {
	record, err := s.repo.GetByDate(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.GetByDate: %w", err)
	}
	if record == nil {
		return nil, nil
	}

	items := s.loadItems(ctx, record.ID)
	return models.DailyNutritionOverrideFromRecord(record, items), nil
}

func (s *dailyNutritionOverrideService) ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.DailyNutritionOverride, error) {
	records, err := s.repo.ListByRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.ListByRange: %w", err)
	}

	out := make([]models.DailyNutritionOverride, len(records))
	for i := range records {
		items := s.loadItems(ctx, records[i].ID)
		out[i] = *models.DailyNutritionOverrideFromRecord(&records[i], items)
	}
	return out, nil
}

func (s *dailyNutritionOverrideService) Update(ctx context.Context, userID string, id string, input models.UpdateOverrideInput) (*models.DailyNutritionOverride, error) {
	existing, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.Update: %w", err)
	}
	if existing == nil {
		return nil, ErrOverrideNotFound
	}

	notes := input.Notes
	if notes == nil {
		notes = existing.Notes
	}

	record, err := s.repo.Update(ctx, userID, id, notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrOverrideNotFound
	}

	items := s.loadItems(ctx, record.ID)
	return models.DailyNutritionOverrideFromRecord(record, items), nil
}

func (s *dailyNutritionOverrideService) Delete(ctx context.Context, userID string, id string) (*models.DailyNutritionOverride, error) {
	record, err := s.repo.Delete(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrOverrideNotFound
	}
	return models.DailyNutritionOverrideFromRecord(record, nil), nil
}

func (s *dailyNutritionOverrideService) CreateItem(ctx context.Context, userID string, input models.CreateOverrideItemInput) (*models.DailyNutritionOverrideItem, error) {
	if input.AmountGrams <= 0 {
		return nil, ErrOverrideItemAmountInvalid
	}
	if input.Operation != models.OperationAdd && input.Operation != models.OperationSubtract && input.Operation != models.OperationReplace {
		return nil, ErrOverrideItemOperationInvalid
	}

	// Verify override exists and belongs to user
	override, err := s.repo.GetByID(ctx, userID, input.OverrideID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.CreateItem: %w", err)
	}
	if override == nil {
		return nil, ErrOverrideNotFound
	}

	record, err := s.itemRepo.Create(ctx, input.OverrideID, input.ProductID, input.AmountGrams, string(input.Operation), input.MealLabel, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.CreateItem: %w", err)
	}

	return models.DailyNutritionOverrideItemFromRecord(record), nil
}

func (s *dailyNutritionOverrideService) UpdateItem(ctx context.Context, userID string, itemID string, input models.UpdateOverrideItemInput) (*models.DailyNutritionOverrideItem, error) {
	existing, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.UpdateItem: %w", err)
	}
	if existing == nil {
		return nil, ErrOverrideItemNotFound
	}

	amount := existing.AmountGrams
	if input.AmountGrams != nil {
		if *input.AmountGrams <= 0 {
			return nil, ErrOverrideItemAmountInvalid
		}
		amount = *input.AmountGrams
	}

	op := existing.Operation
	if input.Operation != nil {
		if *input.Operation != models.OperationAdd && *input.Operation != models.OperationSubtract && *input.Operation != models.OperationReplace {
			return nil, ErrOverrideItemOperationInvalid
		}
		op = string(*input.Operation)
	}

	mealLabel := input.MealLabel
	if mealLabel == nil {
		mealLabel = existing.MealLabel
	}

	notes := input.Notes
	if notes == nil {
		notes = existing.Notes
	}

	record, err := s.itemRepo.Update(ctx, itemID, amount, op, mealLabel, notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.UpdateItem: %w", err)
	}
	if record == nil {
		return nil, ErrOverrideItemNotFound
	}

	return models.DailyNutritionOverrideItemFromRecord(record), nil
}

func (s *dailyNutritionOverrideService) DeleteItem(ctx context.Context, userID string, itemID string) (*models.DailyNutritionOverrideItem, error) {
	record, err := s.itemRepo.Delete(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.DeleteItem: %w", err)
	}
	if record == nil {
		return nil, ErrOverrideItemNotFound
	}
	return models.DailyNutritionOverrideItemFromRecord(record), nil
}
```

- [ ] **Step 5: Create nutrition_macro_service.go**

```go
package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
)

type NutritionMacroService interface {
	Calculate(ctx context.Context, userID string, weekStartDate string, date string) (*models.NutritionMacros, error)
}

type nutritionMacroService struct {
	tmplRepo         postgres.NutritionTemplateRepository
	itemRepo         postgres.NutritionTemplateItemRepository
	overrideRepo     postgres.DailyNutritionOverrideRepository
	overrideItemRepo postgres.DailyNutritionOverrideItemRepository
	productRepo      postgres.NutritionProductRepository
	logger           *zap.Logger
}

func NewNutritionMacroService(
	tmplRepo postgres.NutritionTemplateRepository,
	itemRepo postgres.NutritionTemplateItemRepository,
	overrideRepo postgres.DailyNutritionOverrideRepository,
	overrideItemRepo postgres.DailyNutritionOverrideItemRepository,
	productRepo postgres.NutritionProductRepository,
	logger *zap.Logger,
) NutritionMacroService {
	return &nutritionMacroService{
		tmplRepo:         tmplRepo,
		itemRepo:         itemRepo,
		overrideRepo:     overrideRepo,
		overrideItemRepo: overrideItemRepo,
		productRepo:      productRepo,
		logger:           logger,
	}
}

func (s *nutritionMacroService) Calculate(ctx context.Context, userID string, weekStartDate string, date string) (*models.NutritionMacros, error) {
	s.logger.Info("[NutritionMacros][calculate]",
	tmpl, err := s.tmplRepo.GetByWeek(ctx, userID, weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_macro_service.Calculate: %w", err)
	}
	if tmpl == nil {
		return &models.NutritionMacros{}, nil
	}

	items, err := s.itemRepo.ListByTemplate(ctx, tmpl.ID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_macro_service.Calculate: %w", err)
	}

	result := &models.NutritionMacros{}

	// Calculate per-item macros and build a map of product_id -> template macros
	type templateMacro struct {
		calories, protein, fat, carbs float64
	}
	tmplMacrosByProduct := make(map[string]templateMacro)

	for _, item := range items {
		product, err := s.productRepo.GetByIDIncludeInactive(ctx, userID, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("nutrition_macro_service.Calculate: %w", err)
		}
		if product == nil || !product.IsActive {
			continue
		}

		factor := item.AmountGrams / 100.0
		m := templateMacro{
			calories: product.CaloriesPer100g * factor,
			protein:  product.ProteinPer100g * factor,
			fat:      product.FatPer100g * factor,
			carbs:    product.CarbsPer100g * factor,
		}
		result.Calories += m.calories
		result.Protein += m.protein
		result.Fat += m.fat
		result.Carbs += m.carbs
		// Accumulate per-product totals (product may appear in multiple items)
		prev := tmplMacrosByProduct[item.ProductID]
		prev.calories += m.calories
		prev.protein += m.protein
		prev.fat += m.fat
		prev.carbs += m.carbs
		tmplMacrosByProduct[item.ProductID] = prev
	}

	// Apply override for the specific date if requested
	if date == "" {
		return result, nil
	}

	override, err := s.overrideRepo.GetByDate(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("nutrition_macro_service.Calculate: %w", err)
	}
	if override == nil {
		return result, nil
	}

	overrideItems, err := s.overrideItemRepo.ListByOverride(ctx, override.ID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_macro_service.Calculate: %w", err)
	}

	for _, oi := range overrideItems {
		product, err := s.productRepo.GetByIDIncludeInactive(ctx, userID, oi.ProductID)
		if err != nil {
			return nil, fmt.Errorf("nutrition_macro_service.Calculate: %w", err)
		}
		if product == nil || !product.IsActive {
			continue
		}

		factor := oi.AmountGrams / 100.0
		overrideCal := product.CaloriesPer100g * factor
		overrideProtein := product.ProteinPer100g * factor
		overrideFat := product.FatPer100g * factor
		overrideCarbs := product.CarbsPer100g * factor

		switch oi.Operation {
		case string(models.OperationAdd):
			result.Calories += overrideCal
			result.Protein += overrideProtein
			result.Fat += overrideFat
			result.Carbs += overrideCarbs
		case string(models.OperationSubtract):
			result.Calories -= overrideCal
			result.Protein -= overrideProtein
			result.Fat -= overrideFat
			result.Carbs -= overrideCarbs
		case string(models.OperationReplace):
			tmpl, exists := tmplMacrosByProduct[oi.ProductID]
			if exists {
				result.Calories = result.Calories - tmpl.calories + overrideCal
				result.Protein = result.Protein - tmpl.protein + overrideProtein
				result.Fat = result.Fat - tmpl.fat + overrideFat
				result.Carbs = result.Carbs - tmpl.carbs + overrideCarbs
			} else {
				result.Calories += overrideCal
				result.Protein += overrideProtein
				result.Fat += overrideFat
				result.Carbs += overrideCarbs
			}
		}
	}

	return result, nil
}
```

- [ ] **Step 6: Verify services compile**

Run: `bunx nx run api:typecheck` (or `go build ./apps/api/...`)
Expected: All services compile without errors.

- [ ] **Step 7: Commit**

```bash
git add apps/api/internal/atlas/service/nutrition_*.go
git commit -m "feat(wave-05): add nutrition service layer"
```

---

### Task 6: GraphQL Schema — `nutrition.graphql`

**Files:**
- Create: `apps/api/internal/atlas/graph/schema/nutrition.graphql`
- Modify: `apps/api/internal/atlas/graph/schema/schema.graphql`

- [ ] **Step 1: Create nutrition.graphql**

```graphql
type NutritionProduct {
  id: ID!
  userId: ID!
  name: String!
  caloriesPer100g: Float!
  proteinPer100g: Float!
  fatPer100g: Float!
  carbsPer100g: Float!
  notes: String
  isActive: Boolean!
  createdAt: Time!
  updatedAt: Time!
}

type NutritionTemplate {
  id: ID!
  userId: ID!
  weekStartDate: Date!
  title: String
  notes: String
  items: [NutritionTemplateItem!]!
  createdAt: Time!
  updatedAt: Time!
}

type NutritionTemplateItem {
  id: ID!
  templateId: ID!
  productId: ID!
  amountGrams: Float!
  mealLabel: String
  notes: String
  createdAt: Time!
  updatedAt: Time!
}

type DailyNutritionOverride {
  id: ID!
  userId: ID!
  date: Date!
  notes: String
  items: [DailyNutritionOverrideItem!]!
  createdAt: Time!
  updatedAt: Time!
}

type DailyNutritionOverrideItem {
  id: ID!
  overrideId: ID!
  productId: ID!
  amountGrams: Float!
  operation: Operation!
  mealLabel: String
  notes: String
  createdAt: Time!
  updatedAt: Time!
}

type NutritionMacros {
  calories: Float!
  protein: Float!
  fat: Float!
  carbs: Float!
}

enum Operation {
  ADD
  SUBTRACT
  REPLACE
}

enum NutritionErrorCode {
  VALIDATION_ERROR
  NOT_FOUND
  AUTH_ERROR
  INTERNAL_ERROR
}

input CreateProductInput {
  name: String!
  caloriesPer100g: Float!
  proteinPer100g: Float!
  fatPer100g: Float!
  carbsPer100g: Float!
  notes: String
}

input UpdateProductInput {
  name: String
  caloriesPer100g: Float
  proteinPer100g: Float
  fatPer100g: Float
  carbsPer100g: Float
  notes: String
}

input CreateTemplateInput {
  weekStartDate: Date!
  title: String
  notes: String
}

input UpdateTemplateInput {
  title: String
  notes: String
}

input CreateTemplateItemInput {
  templateId: ID!
  productId: ID!
  amountGrams: Float!
  mealLabel: String
  notes: String
}

input UpdateTemplateItemInput {
  amountGrams: Float
  mealLabel: String
  notes: String
}

input CreateOverrideInput {
  date: Date!
  notes: String
}

input UpdateOverrideInput {
  notes: String
}

input CreateOverrideItemInput {
  overrideId: ID!
  productId: ID!
  amountGrams: Float!
  operation: Operation!
  mealLabel: String
  notes: String
}

input UpdateOverrideItemInput {
  amountGrams: Float
  operation: Operation
  mealLabel: String
  notes: String
}

type NutritionProductResult {
  nutritionProduct: NutritionProduct
  validationError: NutritionValidationError
  notFoundError: NutritionNotFoundError
  authError: NutritionAuthError
}

type NutritionProductsResult {
  products: [NutritionProduct!]!
  validationError: NutritionValidationError
  authError: NutritionAuthError
}

type NutritionTemplateResult {
  nutritionTemplate: NutritionTemplate
  validationError: NutritionValidationError
  notFoundError: NutritionNotFoundError
  authError: NutritionAuthError
}

type NutritionTemplatesResult {
  templates: [NutritionTemplate!]!
  validationError: NutritionValidationError
  authError: NutritionAuthError
}

type NutritionTemplateItemResult {
  nutritionTemplateItem: NutritionTemplateItem
  validationError: NutritionValidationError
  notFoundError: NutritionNotFoundError
  authError: NutritionAuthError
}

type DailyNutritionOverrideResult {
  dailyNutritionOverride: DailyNutritionOverride
  validationError: NutritionValidationError
  notFoundError: NutritionNotFoundError
  authError: NutritionAuthError
}

type DailyNutritionOverridesResult {
  overrides: [DailyNutritionOverride!]!
  validationError: NutritionValidationError
  authError: NutritionAuthError
}

type DailyNutritionOverrideItemResult {
  dailyNutritionOverrideItem: DailyNutritionOverrideItem
  validationError: NutritionValidationError
  notFoundError: NutritionNotFoundError
  authError: NutritionAuthError
}

type NutritionMacrosResult {
  macros: NutritionMacros
  validationError: NutritionValidationError
  authError: NutritionAuthError
}

type NutritionValidationError {
  message: String!
  code: NutritionErrorCode!
}

type NutritionNotFoundError {
  message: String!
  code: NutritionErrorCode!
}

type NutritionAuthError {
  message: String!
  code: NutritionErrorCode!
}
```

- [ ] **Step 2: Add Query and Mutation fields to schema.graphql**

```graphql
  # WAVE-05 Nutrition
  nutritionProducts: NutritionProductsResult!
  nutritionProduct(id: ID!): NutritionProductResult!
  nutritionTemplates(startDate: Date!, endDate: Date!): NutritionTemplatesResult!
  nutritionTemplate(id: ID!): NutritionTemplateResult!
  nutritionTemplateCurrent(weekStartDate: Date!): NutritionTemplateResult!
  dailyNutritionOverrides(startDate: Date!, endDate: Date!): DailyNutritionOverridesResult!
  dailyNutritionOverride(id: ID!): DailyNutritionOverrideResult!
  dailyNutritionOverrideByDate(date: Date!): DailyNutritionOverrideResult!
  nutritionMacros(weekStartDate: Date!, date: Date): NutritionMacrosResult!
```

And to the Mutation type:

```graphql
  # WAVE-05 Nutrition
  createNutritionProduct(input: CreateProductInput!): NutritionProductResult!
  updateNutritionProduct(id: ID!, input: UpdateProductInput!): NutritionProductResult!
  deleteNutritionProduct(id: ID!): NutritionProductResult!
  createNutritionTemplate(input: CreateTemplateInput!): NutritionTemplateResult!
  updateNutritionTemplate(id: ID!, input: UpdateTemplateInput!): NutritionTemplateResult!
  deleteNutritionTemplate(id: ID!): NutritionTemplateResult!
  createNutritionTemplateItem(input: CreateTemplateItemInput!): NutritionTemplateItemResult!
  updateNutritionTemplateItem(id: ID!, input: UpdateTemplateItemInput!): NutritionTemplateItemResult!
  deleteNutritionTemplateItem(id: ID!): NutritionTemplateItemResult!
  createDailyNutritionOverride(input: CreateOverrideInput!): DailyNutritionOverrideResult!
  updateDailyNutritionOverride(id: ID!, input: UpdateOverrideInput!): DailyNutritionOverrideResult!
  deleteDailyNutritionOverride(id: ID!): DailyNutritionOverrideResult!
  createDailyNutritionOverrideItem(input: CreateOverrideItemInput!): DailyNutritionOverrideItemResult!
  updateDailyNutritionOverrideItem(id: ID!, input: UpdateOverrideItemInput!): DailyNutritionOverrideItemResult!
  deleteDailyNutritionOverrideItem(id: ID!): DailyNutritionOverrideItemResult!
```

- [ ] **Step 3: Validate GraphQL schema**

Run: `bunx nx run graphql:validate`
Expected: Schema validation passes without errors.

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/atlas/graph/schema/nutrition.graphql apps/api/internal/atlas/graph/schema/schema.graphql
git commit -m "feat(wave-05): add nutrition GraphQL schema"
```

---

### Task 7: Resolvers — `resolver/nutrition.go`

**Files:**
- Create: `apps/api/internal/atlas/graph/resolver/nutrition.go`
- Modify: `apps/api/internal/atlas/graph/resolver/resolver.go`

- [ ] **Step 1: Add service fields to resolver.go**

```go
type Resolver struct {
	// ... existing fields ...
	NutritionProductService      service.NutritionProductService
	NutritionTemplateService     service.NutritionTemplateService
	NutritionTemplateItemService  service.NutritionTemplateItemService
	DailyNutritionOverrideService service.DailyNutritionOverrideService
	NutritionMacroService        service.NutritionMacroService
}
```

- [ ] **Step 2: Create nutrition.go resolver file**

Each resolver follows the cardio.go pattern: extract userID, check auth, call service, handle errors.

```go
package resolver

import (
	"context"
	"errors"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
)

// NutritionProduct resolvers
func (r *Resolver) GetNutritionProducts(ctx context.Context) (*models.NutritionProductsResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionProductsResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	products, err := r.NutritionProductService.ListActive(ctx, userID)
	if err != nil {
		return nil, nil
	}

	return &models.NutritionProductsResult{Products: products}, nil
}

func (r *Resolver) GetNutritionProduct(ctx context.Context, id string) (*models.NutritionProductResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionProductResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	product, err := r.NutritionProductService.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrProductNotFound) {
			return &models.NutritionProductResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "product not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.NutritionProductResult{NutritionProduct: product}, nil
}

func (r *Resolver) CreateNutritionProduct(ctx context.Context, input models.CreateProductInput) (*models.NutritionProductResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionProductResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	product, err := r.NutritionProductService.Create(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrProductNameRequired),
			errors.Is(err, atlasService.ErrProductMacroNegative),
			errors.Is(err, atlasService.ErrProductNameTooLong):
			return &models.NutritionProductResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.NutritionProductResult{NutritionProduct: product}, nil
}

func (r *Resolver) UpdateNutritionProduct(ctx context.Context, id string, input models.UpdateProductInput) (*models.NutritionProductResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionProductResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	product, err := r.NutritionProductService.Update(ctx, userID, id, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrProductNameRequired),
			errors.Is(err, atlasService.ErrProductMacroNegative),
			errors.Is(err, atlasService.ErrProductNameTooLong):
			return &models.NutritionProductResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		case errors.Is(err, atlasService.ErrProductNotFound):
			return &models.NutritionProductResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "product not found", Code: models.NutritionErrorNotFound},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.NutritionProductResult{NutritionProduct: product}, nil
}

func (r *Resolver) DeleteNutritionProduct(ctx context.Context, id string) (*models.NutritionProductResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionProductResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	product, err := r.NutritionProductService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrProductNotFound) {
			return &models.NutritionProductResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "product not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.NutritionProductResult{NutritionProduct: product}, nil
}

// NutritionTemplate resolvers
func (r *Resolver) GetNutritionTemplates(ctx context.Context, startDate, endDate string) (*models.NutritionTemplatesResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplatesResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	templates, err := r.NutritionTemplateService.ListByRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, nil
	}

	return &models.NutritionTemplatesResult{Templates: templates}, nil
}

func (r *Resolver) GetNutritionTemplate(ctx context.Context, id string) (*models.NutritionTemplateResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	tmpl, err := r.NutritionTemplateService.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrTemplateNotFound) {
			return &models.NutritionTemplateResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "template not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.NutritionTemplateResult{NutritionTemplate: tmpl}, nil
}

func (r *Resolver) GetNutritionTemplateCurrent(ctx context.Context, weekStartDate string) (*models.NutritionTemplateResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	tmpl, err := r.NutritionTemplateService.GetCurrent(ctx, userID, weekStartDate)
	if err != nil {
		return nil, nil
	}
	if tmpl == nil {
		return &models.NutritionTemplateResult{}, nil
	}

	return &models.NutritionTemplateResult{NutritionTemplate: tmpl}, nil
}

func (r *Resolver) CreateNutritionTemplate(ctx context.Context, input models.CreateTemplateInput) (*models.NutritionTemplateResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	tmpl, err := r.NutritionTemplateService.Create(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrTemplateWeekRequired):
			return &models.NutritionTemplateResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.NutritionTemplateResult{NutritionTemplate: tmpl}, nil
}

func (r *Resolver) UpdateNutritionTemplate(ctx context.Context, id string, input models.UpdateTemplateInput) (*models.NutritionTemplateResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	tmpl, err := r.NutritionTemplateService.Update(ctx, userID, id, input)
	if err != nil {
		if errors.Is(err, atlasService.ErrTemplateNotFound) {
			return &models.NutritionTemplateResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "template not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.NutritionTemplateResult{NutritionTemplate: tmpl}, nil
}

func (r *Resolver) DeleteNutritionTemplate(ctx context.Context, id string) (*models.NutritionTemplateResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	tmpl, err := r.NutritionTemplateService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrTemplateNotFound) {
			return &models.NutritionTemplateResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "template not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.NutritionTemplateResult{NutritionTemplate: tmpl}, nil
}

// NutritionTemplateItem resolvers
func (r *Resolver) CreateNutritionTemplateItem(ctx context.Context, input models.CreateTemplateItemInput) (*models.NutritionTemplateItemResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateItemResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	item, err := r.NutritionTemplateItemService.Create(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrTemplateItemAmountInvalid):
			return &models.NutritionTemplateItemResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		case errors.Is(err, atlasService.ErrTemplateNotFound):
			return &models.NutritionTemplateItemResult{
				ValidationErr: &models.NutritionValidationErr{Message: "template not found", Code: models.NutritionErrorValidation},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.NutritionTemplateItemResult{NutritionTemplateItem: item}, nil
}

func (r *Resolver) UpdateNutritionTemplateItem(ctx context.Context, id string, input models.UpdateTemplateItemInput) (*models.NutritionTemplateItemResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateItemResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	item, err := r.NutritionTemplateItemService.Update(ctx, userID, id, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrTemplateItemAmountInvalid):
			return &models.NutritionTemplateItemResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		case errors.Is(err, atlasService.ErrTemplateItemNotFound):
			return &models.NutritionTemplateItemResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "item not found", Code: models.NutritionErrorNotFound},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.NutritionTemplateItemResult{NutritionTemplateItem: item}, nil
}

func (r *Resolver) DeleteNutritionTemplateItem(ctx context.Context, id string) (*models.NutritionTemplateItemResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateItemResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	item, err := r.NutritionTemplateItemService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrTemplateItemNotFound) {
			return &models.NutritionTemplateItemResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "item not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.NutritionTemplateItemResult{NutritionTemplateItem: item}, nil
}

// DailyNutritionOverride resolvers
func (r *Resolver) GetDailyNutritionOverrides(ctx context.Context, startDate, endDate string) (*models.DailyNutritionOverridesResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverridesResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	overrides, err := r.DailyNutritionOverrideService.ListByRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, nil
	}

	return &models.DailyNutritionOverridesResult{Overrides: overrides}, nil
}

func (r *Resolver) GetDailyNutritionOverride(ctx context.Context, id string) (*models.DailyNutritionOverrideResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	override, err := r.DailyNutritionOverrideService.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrOverrideNotFound) {
			return &models.DailyNutritionOverrideResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "override not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.DailyNutritionOverrideResult{DailyNutritionOverride: override}, nil
}

func (r *Resolver) GetDailyNutritionOverrideByDate(ctx context.Context, date string) (*models.DailyNutritionOverrideResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	override, err := r.DailyNutritionOverrideService.GetByDate(ctx, userID, date)
	if err != nil {
		return nil, nil
	}
	if override == nil {
		return &models.DailyNutritionOverrideResult{}, nil
	}

	return &models.DailyNutritionOverrideResult{DailyNutritionOverride: override}, nil
}

func (r *Resolver) CreateDailyNutritionOverride(ctx context.Context, input models.CreateOverrideInput) (*models.DailyNutritionOverrideResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	override, err := r.DailyNutritionOverrideService.Create(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrOverrideDateRequired):
			return &models.DailyNutritionOverrideResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.DailyNutritionOverrideResult{DailyNutritionOverride: override}, nil
}

func (r *Resolver) UpdateDailyNutritionOverride(ctx context.Context, id string, input models.UpdateOverrideInput) (*models.DailyNutritionOverrideResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	override, err := r.DailyNutritionOverrideService.Update(ctx, userID, id, input)
	if err != nil {
		if errors.Is(err, atlasService.ErrOverrideNotFound) {
			return &models.DailyNutritionOverrideResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "override not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.DailyNutritionOverrideResult{DailyNutritionOverride: override}, nil
}

func (r *Resolver) DeleteDailyNutritionOverride(ctx context.Context, id string) (*models.DailyNutritionOverrideResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	override, err := r.DailyNutritionOverrideService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrOverrideNotFound) {
			return &models.DailyNutritionOverrideResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "override not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.DailyNutritionOverrideResult{DailyNutritionOverride: override}, nil
}

// DailyNutritionOverrideItem resolvers
func (r *Resolver) CreateDailyNutritionOverrideItem(ctx context.Context, input models.CreateOverrideItemInput) (*models.DailyNutritionOverrideItemResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideItemResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	// Delegate to override service's item CRUD
	item, err := r.DailyNutritionOverrideService.CreateItem(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrOverrideNotFound):
			return &models.DailyNutritionOverrideItemResult{
				ValidationErr: &models.NutritionValidationErr{Message: "override not found", Code: models.NutritionErrorValidation},
			}, nil
		case errors.Is(err, atlasService.ErrOverrideItemAmountInvalid),
			errors.Is(err, atlasService.ErrOverrideItemOperationInvalid):
			return &models.DailyNutritionOverrideItemResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.DailyNutritionOverrideItemResult{DailyNutritionOverrideItem: item}, nil
}

func (r *Resolver) UpdateDailyNutritionOverrideItem(ctx context.Context, id string, input models.UpdateOverrideItemInput) (*models.DailyNutritionOverrideItemResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideItemResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	item, err := r.DailyNutritionOverrideService.UpdateItem(ctx, userID, id, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrOverrideItemAmountInvalid),
			errors.Is(err, atlasService.ErrOverrideItemOperationInvalid):
			return &models.DailyNutritionOverrideItemResult{
				ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
			}, nil
		case errors.Is(err, atlasService.ErrOverrideItemNotFound):
			return &models.DailyNutritionOverrideItemResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "override item not found", Code: models.NutritionErrorNotFound},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.DailyNutritionOverrideItemResult{DailyNutritionOverrideItem: item}, nil
}

func (r *Resolver) DeleteDailyNutritionOverrideItem(ctx context.Context, id string) (*models.DailyNutritionOverrideItemResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.DailyNutritionOverrideItemResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	item, err := r.DailyNutritionOverrideService.DeleteItem(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrOverrideItemNotFound) {
			return &models.DailyNutritionOverrideItemResult{
				NotFoundErr: &models.NutritionNotFoundErr{Message: "override item not found", Code: models.NutritionErrorNotFound},
			}, nil
		}
		return nil, nil
	}

	return &models.DailyNutritionOverrideItemResult{DailyNutritionOverrideItem: item}, nil
}

// NutritionMacros resolver
func (r *Resolver) GetNutritionMacros(ctx context.Context, weekStartDate string, date *string) (*models.NutritionMacrosResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionMacrosResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	d := ""
	if date != nil {
		d = *date
	}

	macros, err := r.NutritionMacroService.Calculate(ctx, userID, weekStartDate, d)
	if err != nil {
		return nil, nil
	}

	return &models.NutritionMacrosResult{Macros: macros}, nil
}
```

Note: Override item resolvers need proper implementation through a separate service. After gqlgen generates stubs, fill in the override item CRUD through the override service or a separate override item service.

- [ ] **Step 3: Run gqlgen codegen**

Run: `bunx nx run api:codegen`
Expected: gqlgen generates resolver stubs in `nutrition.resolvers.go`. Then fill in the override item resolvers with proper service delegations.

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/atlas/graph/resolver/nutrition.go apps/api/internal/atlas/graph/resolver/resolver.go
git commit -m "feat(wave-05): add nutrition GraphQL resolvers"
```

---

### Task 8: Wiring — `atlas-gqlgen.yml` + `main.go`

**Files:**
- Modify: `apps/api/atlas-gqlgen.yml`
- Modify: `apps/api/cmd/server/main.go`

- [ ] **Step 1: Add model bindings to atlas-gqlgen.yml**

Add after the WeekFlag section:

```yaml
  # WAVE-05 Nutrition
  NutritionProduct:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionProduct
  CreateProductInput:
    model: monorepo-template/apps/api/internal/atlas/models.CreateProductInput
  UpdateProductInput:
    model: monorepo-template/apps/api/internal/atlas/models.UpdateProductInput
  NutritionProductResult:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionProductResult
  NutritionProductsResult:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionProductsResult
  NutritionTemplate:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionTemplate
  CreateTemplateInput:
    model: monorepo-template/apps/api/internal/atlas/models.CreateTemplateInput
  UpdateTemplateInput:
    model: monorepo-template/apps/api/internal/atlas/models.UpdateTemplateInput
  NutritionTemplateResult:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionTemplateResult
  NutritionTemplatesResult:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionTemplatesResult
  NutritionTemplateItem:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionTemplateItem
  CreateTemplateItemInput:
    model: monorepo-template/apps/api/internal/atlas/models.CreateTemplateItemInput
  UpdateTemplateItemInput:
    model: monorepo-template/apps/api/internal/atlas/models.UpdateTemplateItemInput
  NutritionTemplateItemResult:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionTemplateItemResult
  DailyNutritionOverride:
    model: monorepo-template/apps/api/internal/atlas/models.DailyNutritionOverride
  CreateOverrideInput:
    model: monorepo-template/apps/api/internal/atlas/models.CreateOverrideInput
  UpdateOverrideInput:
    model: monorepo-template/apps/api/internal/atlas/models.UpdateOverrideInput
  DailyNutritionOverrideResult:
    model: monorepo-template/apps/api/internal/atlas/models.DailyNutritionOverrideResult
  DailyNutritionOverridesResult:
    model: monorepo-template/apps/api/internal/atlas/models.DailyNutritionOverridesResult
  DailyNutritionOverrideItem:
    model: monorepo-template/apps/api/internal/atlas/models.DailyNutritionOverrideItem
  CreateOverrideItemInput:
    model: monorepo-template/apps/api/internal/atlas/models.CreateOverrideItemInput
  UpdateOverrideItemInput:
    model: monorepo-template/apps/api/internal/atlas/models.UpdateOverrideItemInput
  DailyNutritionOverrideItemResult:
    model: monorepo-template/apps/api/internal/atlas/models.DailyNutritionOverrideItemResult
  NutritionMacros:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionMacros
  NutritionMacrosResult:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionMacrosResult
  Operation:
    model: monorepo-template/apps/api/internal/atlas/models.Operation
  NutritionErrorCode:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionErrorCode
  NutritionValidationError:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionValidationErr
  NutritionNotFoundError:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionNotFoundErr
  NutritionAuthError:
    model: monorepo-template/apps/api/internal/atlas/models.NutritionAuthErr
```

- [ ] **Step 2: Wire repos, services, and resolver in main.go**

Add repo creation (near other atlas repo lines):
```go
atlasNutritionProductRepo := atlasPostgres.NewNutritionProductRepository(db.Pool)
atlasNutritionTemplateRepo := atlasPostgres.NewNutritionTemplateRepository(db.Pool)
atlasNutritionTemplateItemRepo := atlasPostgres.NewNutritionTemplateItemRepository(db.Pool)
atlasNutritionOverrideRepo := atlasPostgres.NewDailyNutritionOverrideRepository(db.Pool)
atlasNutritionOverrideItemRepo := atlasPostgres.NewDailyNutritionOverrideItemRepository(db.Pool)
```

Add service creation (near other atlas service lines). All nutrition services need the zap logger `l` for log markers:
```go
atlasNutritionProductService := atlasService.NewNutritionProductService(atlasNutritionProductRepo, l)
atlasNutritionTemplateService := atlasService.NewNutritionTemplateService(atlasNutritionTemplateRepo, atlasNutritionTemplateItemRepo, l)
atlasNutritionTemplateItemService := atlasService.NewNutritionTemplateItemService(atlasNutritionTemplateItemRepo, atlasNutritionTemplateRepo, l)
atlasNutritionOverrideService := atlasService.NewNutritionOverrideService(atlasNutritionOverrideRepo, atlasNutritionOverrideItemRepo, l)
atlasNutritionMacroService := atlasService.NewNutritionMacroService(
    atlasNutritionTemplateRepo,
    atlasNutritionTemplateItemRepo,
    atlasNutritionOverrideRepo,
    atlasNutritionOverrideItemRepo,
    atlasNutritionProductRepo,
    l,
)
```

Add to atlasRes struct:
```go
NutritionProductService:      atlasNutritionProductService,
NutritionTemplateService:     atlasNutritionTemplateService,
NutritionTemplateItemService:  atlasNutritionTemplateItemService,
DailyNutritionOverrideService: atlasNutritionOverrideService,
NutritionMacroService:        atlasNutritionMacroService,
```

- [ ] **Step 3: Run full codegen and build**

Run: `bunx nx run api:codegen`
Run: `bunx nx run api:build` (or `go build ./apps/api/...`)
Expected: Build succeeds with no errors.

- [ ] **Step 4: Commit**

```bash
git add apps/api/atlas-gqlgen.yml apps/api/cmd/server/main.go
git commit -m "feat(wave-05): wire nutrition repos, services, and resolvers"
```

---

### Task 9: Tests

**Files:**
- Create: `apps/api/internal/atlas/service/nutrition_product_service_test.go`
- Create: `apps/api/internal/atlas/service/nutrition_template_service_test.go`
- Create: `apps/api/internal/atlas/service/nutrition_macro_service_test.go`
- Create: `apps/api/internal/atlas/service/nutrition_override_service_test.go`
- Create: `apps/api/internal/atlas/service/nutrition_override_item_service_test.go`

Each test file follows the existing `cardio_service_test.go` pattern: in-memory mocks or test DB, table-driven tests for success + error paths.

- [ ] **Step 1: Write nutrition_product_service_test.go**

Scenarios:
- Create: success with all fields, name empty → ErrProductNameRequired, negative macros → ErrProductMacroNegative, name too long → ErrProductNameTooLong
- GetByID: existing product returns product, non-existent returns ErrProductNotFound, soft-deleted product still returns by ID
- ListActive: returns only active products
- Update: success, name empty → error, negative macros → error, non-existent → ErrProductNotFound
- Delete (soft-delete): success, non-existent → ErrProductNotFound

- [ ] **Step 2: Write nutrition_template_service_test.go**

Scenarios:
- Create: success with upsert (second create for same week replaces), weekStartDate empty → error
- GetByID: returns template with items loaded, non-existent → ErrTemplateNotFound
- GetCurrent: returns template by week (nil if none)
- ListByRange: returns templates in date range
- Update: title/notes, weekStartDate immutable
- Delete: cascade (verify via mock that items are loaded but template is gone)

- [ ] **Step 3: Write nutrition_macro_service_test.go**

Scenarios:
- Empty week (no template) → all zeros
- Template with items → correct per-day macros (calories = calories_per_100g * amount_grams / 100)
- Template + override ADD: override macros added to template total
- Template + override SUBTRACT: override macros subtracted from template total
- Template + override REPLACE: template product macros replaced by override
- Same product in multiple template items with REPLACE → all occurrences replaced correctly
- Soft-deleted product in template → contributes 0
- No override for requested date → template values only

- [ ] **Step 4: Write nutrition_override_service_test.go**

Scenarios:
- Create: success, duplicate date returns existing (upsert), date empty → error
- GetByID: returns override with items, non-existent → ErrOverrideNotFound
- GetByDate: returns override (nil if none)
- ListByRange: returns overrides in range
- Update: notes only, date immutable
- Delete: cascade (verify items deleted)

- [ ] **Step 5: Write override item CRUD tests alongside override_service_test.go**

Scenarios:
- CreateItem: success, amountGrams <= 0 → error, invalid operation → error, non-existent override → ErrOverrideNotFound
- UpdateItem: success, invalid amount → error, invalid operation → error
- DeleteItem: non-existent → ErrOverrideItemNotFound

- [ ] **Step 6: Run all WAVE-05 tests**

Run: `bunx nx run api:test -- --run '(?i)nutrition'`
Expected: All tests pass.

- [ ] **Step 7: Commit**

```bash
git add apps/api/internal/atlas/service/nutrition_*_test.go
git commit -m "feat(wave-05): add nutrition service tests"
```

---

### Task 10: Verification

- [ ] **Step 1: Run full lint**

Run: `bunx nx run api:lint`
Expected: No lint errors.

- [ ] **Step 2: Run codegen drift check**

Run: `bunx nx run api:codegen && git diff --stat`
Expected: No uncommitted generated file changes.

- [ ] **Step 3: Run all tests**

Run: `bunx nx run api:test`
Expected: All tests (existing + new) pass.