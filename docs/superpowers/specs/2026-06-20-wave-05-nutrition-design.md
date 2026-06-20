# WAVE-05: Nutrition — Implementation Design

## Status
Draft — awaiting spec review

## Overview

Implement Nutrition module for the Atlas fitness tracker: food tracking through weekly templates and daily overrides with macro calculations. All operations via GraphQL under the existing `/graphql/atlas` endpoint, PIN-guarded by WAVE-01 middleware.

## Approach

Single-pattern for all 5 entities, identical to WAVE-04 cardio/body tracking:

- sqlc query files → generated Go code
- Repository adapters (interface + struct) with UUID conversion and error mapping
- Transport-neutral service layer with validation
- GraphQL schema with union result types
- GraphQL resolvers with PIN auth guard
- Single goose migration for all 5 tables

## Architecture

```
nutrition_product --+
                    |
nutrition_template --+--> models/nutrition.go --> service/ --> resolver/ --> GraphQL
                    |
nutrition_template_item --+
                            |
daily_nutrition_override --+
                            |
daily_nutrition_override_item --+
                                    |
nutrition_macro_service (stateless KJBJU calc)
```

All entities are independent from WAVE-02/03/04 tables. Only FK dependency is `atlas_users` (WAVE-01).

## Database Schema (single migration `00090_nutrition_tables.sql`)

### nutrition_product
- id (UUID PK), user_id (UUID FK→atlas_users), name (VARCHAR NOT NULL)
- calories_per_100g (REAL NOT NULL), protein_per_100g (REAL NOT NULL), fat_per_100g (REAL NOT NULL), carbs_per_100g (REAL NOT NULL)
- notes (TEXT nullable), is_active (BOOLEAN NOT NULL DEFAULT true)
- created_at, updated_at
- CHECK (all macros >= 0), idx_nutrition_product_user (user_id)

### nutrition_template
- id (UUID PK), user_id (UUID FK→atlas_users), week_start_date (DATE NOT NULL)
- title (VARCHAR nullable), notes (TEXT nullable)
- created_at, updated_at
- UNIQUE (user_id, week_start_date), idx_nutrition_template_week (user_id, week_start_date)

### nutrition_template_item
- id (UUID PK), template_id (UUID FK→nutrition_template ON DELETE CASCADE)
- product_id (UUID FK→nutrition_product), amount_grams (REAL NOT NULL)
- meal_label (VARCHAR nullable), notes (TEXT nullable)
- created_at, updated_at
- CHECK (amount_grams > 0), idx_nutrition_template_item_template (template_id)

### daily_nutrition_override
- id (UUID PK), user_id (UUID FK→atlas_users), date (DATE NOT NULL)
- notes (TEXT nullable), created_at, updated_at
- UNIQUE (user_id, date), idx_nutrition_override_date (user_id, date)

### daily_nutrition_override_item
- id (UUID PK), override_id (UUID FK→daily_nutrition_override ON DELETE CASCADE)
- product_id (UUID FK→nutrition_product), amount_grams (REAL NOT NULL)
- operation (VARCHAR NOT NULL), meal_label (VARCHAR nullable), notes (TEXT nullable)
- created_at, updated_at
- CHECK (amount_grams > 0), CHECK (operation IN ('add','subtract','replace'))
- idx_nutrition_override_item_override (override_id)

## sqlc Queries (5 files)

### nutrition_products.sql
- `CreateNutritionProduct :one` — INSERT
- `GetNutritionProductByID :one` — SELECT by id + user_id
- `ListActiveNutritionProducts :many` — SELECT WHERE is_active = true AND user_id = $1
- `UpdateNutritionProduct :one` — UPDATE (name, macros, notes) returning full row
- `DeleteNutritionProduct :one` — UPDATE is_active = false, RETURNING

### nutrition_templates.sql
- `CreateNutritionTemplate :one` — INSERT ON CONFLICT (user_id, week_start_date) DO UPDATE SET title = COALESCE($3, nutrition_template.title), notes = COALESCE($4, nutrition_template.notes) RETURNING * (upsert — same conflict target on insert and update)
- `GetNutritionTemplateByID :one` — SELECT by id + user_id
- `GetNutritionTemplateByWeek :one` — SELECT by user_id + week_start_date (returns template header only; items are loaded separately by service layer)
- `ListNutritionTemplatesByRange :many` — SELECT by user_id + date range
- `UpdateNutritionTemplate :one` — UPDATE title, notes
- `DeleteNutritionTemplate :one` — DELETE (cascade due to FK)

### nutrition_template_items.sql
- `CreateNutritionTemplateItem :one` — INSERT
- `GetNutritionTemplateItemByID :one` — SELECT by id
- `ListNutritionTemplateItemsByTemplate :many` — SELECT by template_id
- `UpdateNutritionTemplateItem :one` — UPDATE amount_grams, meal_label, notes
- `DeleteNutritionTemplateItem :one` — DELETE

### nutrition_overrides.sql
- `CreateDailyNutritionOverride :one` — INSERT ON CONFLICT (user_id, date) DO UPDATE SET notes = EXCLUDED.notes RETURNING *. Returns the existing record on conflict (upsert semantics matching templates)
- `GetDailyNutritionOverrideByID :one` — SELECT by id + user_id
- `GetDailyNutritionOverrideByDate :one` — SELECT by user_id + date
- `ListDailyNutritionOverridesByRange :many` — SELECT by user_id + date range
- `UpdateDailyNutritionOverride :one` — UPDATE notes
- `DeleteDailyNutritionOverride :one` — DELETE

### nutrition_override_items.sql
- `CreateDailyNutritionOverrideItem :one` — INSERT
- `GetDailyNutritionOverrideItemByID :one` — SELECT by id
- `ListDailyNutritionOverrideItemsByOverride :many` — SELECT by override_id
- `UpdateDailyNutritionOverrideItem :one` — UPDATE amount_grams, operation, meal_label, notes
- `DeleteDailyNutritionOverrideItem :one` — DELETE

## Models (`models/nutrition.go`)

Types following the cardio.go pattern:

- **DB Records**: `NutritionProductRecord`, `NutritionTemplateRecord`, `NutritionTemplateItemRecord`, `DailyNutritionOverrideRecord`, `DailyNutritionOverrideItemRecord`
- **Public Models**: `NutritionProduct`, `NutritionTemplate`, `NutritionTemplateItem`, `DailyNutritionOverride`, `DailyNutritionOverrideItem`, `NutritionMacros`
- **Inputs**: `CreateProductInput`, `UpdateProductInput`, `CreateTemplateInput`, `UpdateTemplateInput`, `CreateTemplateItemInput`, `UpdateTemplateItemInput`, `CreateOverrideInput`, `UpdateOverrideInput`, `CreateOverrideItemInput`, `UpdateOverrideItemInput`
- **Result Unions**: `NutritionProductResult`, `NutritionProductsResult`, `NutritionTemplateResult`, `NutritionTemplatesResult`, `DailyNutritionOverrideResult`, `DailyNutritionOverridesResult`, `NutritionTemplateItemResult`, `DailyNutritionOverrideItemResult`, `NutritionMacrosResult`
- **Errors**: `NutritionValidationErr`, `NutritionNotFoundErr`, `NutritionAuthErr` (shared across entities, like WAVE-04 body errors)
- **Enums**: `Operation` (ADD/SUBTRACT/REPLACE), `NutritionErrorCode`
- **Converters**: `NutritionProductFromRecord`, `NutritionTemplateFromRecord`, etc.
- **Macro struct**: `Calories, Protein, Fat, Carbs float64`

## Repository Layer (5 files in `repository/postgres/`)

Each follows the exact cardio pattern:
- Interface + private struct + `New*Repository(pool)`
- `uuidFromString` / `parseTwoUUIDs` for UUID conversion
- `nullableText` / `nullableInt4` helpers
- `recordFromRow()` converter function
- `NOT_FOUND` → returns `(nil, nil)` — service layer handles not-found
- User-scoped: all queries include `user_id` filter

### NutritionProductRepository
- Create, GetByID, ListActive, Update, SoftDelete
- `ListActive` filters `is_active = true` and `user_id`
- `SoftDelete` sets `is_active = false`, returns full record

### NutritionTemplateRepository
- Create (upsert), GetByID, GetByWeek, ListByRange, Update, Delete
- `GetByWeek` returns template + items for a given weekStartDate
- `ListByRange` uses `week_start_date BETWEEN $2 AND $3`

### NutritionTemplateItemRepository
- Create, GetByID, ListByTemplate, Update, Delete
- `ListByTemplate` returns all items for a template

### DailyNutritionOverrideRepository
- Create, GetByID, GetByDate, ListByRange, Update, Delete
- `GetByDate` returns override + items for a specific date
- `ListByRange` uses `date BETWEEN $2 AND $3`

### DailyNutritionOverrideItemRepository
- Create, GetByID, ListByOverride, Update, Delete

## Service Layer (5 files in `service/`)

### NutritionProductService
- `Create`: validate name required, all macros >= 0
- `GetByID`: returns product or ErrNotFound
- `ListActive`: returns all active products
- `Update`: validate name if provided, macros >= 0 if provided
- `Delete`: soft-delete (set is_active = false)
- Error vars: `ErrProductNameRequired`, `ErrProductMacroNegative`, `ErrProductNotFound`

### NutritionTemplateService
- `Create`: validate weekStartDate required, upsert semantics
- `GetByID`: returns template with items or ErrNotFound
- `GetCurrent`: returns template by weekStartDate
- `ListByRange`: returns templates in date range
- `Update`: title/notes only, weekStartDate immutable
- `Delete`: cascade delete (FK handles DB cascade)
- Error vars: `ErrTemplateWeekRequired`, `ErrTemplateNotFound`

### NutritionTemplateItemService
- `Create`: validate amountGrams > 0
- `Update`: validate amountGrams > 0 if provided
- `Delete`
- Error vars: `ErrTemplateItemAmountInvalid`, `ErrTemplateItemNotFound`
- All operations require template_id + user context for authorization

### DailyNutritionOverrideService
- `Create`: validate date required, unique per user+date
- `GetByID`, `GetByDate`, `ListByRange`
- `Update`: notes only, date immutable
- `Delete`: cascade delete
- `GetByDate`: returns override with items or nil (no error — absence means "use template")

### NutritionMacroService
- `Calculate(ctx, userID, weekStartDate, date)` — stateless calculation
- Algorithm:
  1. Load template for weekStartDate
  2. If no template → return all zeros (not an error)
  3. For each template item: get product macros, scale by (amountGrams / 100)
  4. If a date is specified: check for override
  5. If override exists: apply override items per operation. Override item macros calculated identically to template items: (product macros × amount_grams / 100).
     - ADD: add override macro values to template totals
     - SUBTRACT: subtract override macro values from template totals
     - REPLACE: for the given product, use override values instead of template values. If the same product appears in multiple template items, REPLACE applies to all occurrences of that product_id in the template.
  6. If product is soft-deleted (isActive=false): contribute 0 for that item
  7. Return aggregated NutritionMacros

## GraphQL Schema (`schema/nutrition.graphql`)

Following existing WAVE-04 schema patterns exactly:

```graphql
# Types
type NutritionProduct { id, userId, name, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g, notes, isActive, createdAt, updatedAt }
type NutritionTemplate { id, userId, weekStartDate, title, notes, items: [NutritionTemplateItem!], createdAt, updatedAt }
type NutritionTemplateItem { id, templateId, productId, amountGrams, mealLabel, notes, createdAt, updatedAt }
type DailyNutritionOverride { id, userId, date, notes, items: [DailyNutritionOverrideItem!], createdAt, updatedAt }
type DailyNutritionOverrideItem { id, overrideId, productId, amountGrams, operation: Operation!, mealLabel, notes, createdAt, updatedAt }
type NutritionMacros { calories, protein, fat, carbs: Float! }

# Enums
enum Operation { ADD, SUBTRACT, REPLACE }
enum NutritionErrorCode { VALIDATION_ERROR, NOT_FOUND, AUTH_ERROR, INTERNAL_ERROR }

# Inputs
input CreateProductInput { name!, caloriesPer100g!, proteinPer100g!, fatPer100g!, carbsPer100g!, notes }
input UpdateProductInput { name, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g, notes }
input CreateTemplateInput { weekStartDate!, title, notes }
input UpdateTemplateInput { title, notes }
input CreateTemplateItemInput { templateId!, productId!, amountGrams!, mealLabel, notes }
input UpdateTemplateItemInput { amountGrams, mealLabel, notes }
input CreateOverrideInput { date!, notes }
input UpdateOverrideInput { notes }
input CreateOverrideItemInput { overrideId!, productId!, amountGrams!, operation!, mealLabel, notes }
input UpdateOverrideItemInput { amountGrams, operation, mealLabel, notes }

# Result Unions
type NutritionProductsResult { products: [NutritionProduct!]!, validationError: NutritionValidationError, authError: NutritionAuthError }
type NutritionProductResult { nutritionProduct, validationError: NutritionValidationError, notFoundError: NutritionNotFoundError, authError: NutritionAuthError }
type NutritionTemplatesResult { templates: [NutritionTemplate!], validationError: NutritionValidationError, authError: NutritionAuthError }
type NutritionTemplateResult { nutritionTemplate, validationError: NutritionValidationError, notFoundError: NutritionNotFoundError, authError: NutritionAuthError }
type DailyNutritionOverrideResult { dailyNutritionOverride, validationError: NutritionValidationError, notFoundError: NutritionNotFoundError, authError: NutritionAuthError }
type DailyNutritionOverridesResult { overrides: [DailyNutritionOverride!], validationError: NutritionValidationError, authError: NutritionAuthError }
type NutritionTemplateItemResult { nutritionTemplateItem, validationError: NutritionValidationError, notFoundError: NutritionNotFoundError, authError: NutritionAuthError }
type DailyNutritionOverrideItemResult { dailyNutritionOverrideItem, validationError: NutritionValidationError, notFoundError: NutritionNotFoundError, authError: NutritionAuthError }
type NutritionMacrosResult { macros: NutritionMacros, validationError: NutritionValidationError, authError: NutritionAuthError }

# Error types
type NutritionValidationError { message: String!, code: NutritionErrorCode! }
type NutritionNotFoundError { message: String!, code: NutritionErrorCode! }
type NutritionAuthError { message: String!, code: NutritionErrorCode! }
```

### Queries

```graphql
type Query {
  nutritionProducts: NutritionProductsResult!
  nutritionProduct(id: ID!): NutritionProductResult!
  nutritionTemplates(startDate: Date!, endDate: Date!): NutritionTemplatesResult!
  nutritionTemplate(id: ID!): NutritionTemplateResult!
  nutritionTemplateCurrent(weekStartDate: Date!): NutritionTemplateResult!
  dailyNutritionOverrides(startDate: Date!, endDate: Date!): DailyNutritionOverridesResult!
  dailyNutritionOverride(id: ID!): DailyNutritionOverrideResult!
  dailyNutritionOverrideByDate(date: Date!): DailyNutritionOverrideResult!
  nutritionMacros(weekStartDate: Date!, date: Date): NutritionMacrosResult!
}
```

### Mutations

```graphql
type Mutation {
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
}
```

## Resolvers (`resolver/nutrition.go`)

Each resolver follows the exact cardio.go pattern:
1. Extract userID from context via `middleware.GetAtlasUserID(ctx)`
2. If empty → return AuthError result
3. Call service method
4. Switch on error: validation → ValidationError, not-found → NotFoundError, default → return nil
5. Return success result

No `.resolvers.go` file needed — gqlgen with `follow-schema` layout will auto-generate resolver stubs in a `nutrition.resolvers.go` file.

## Wiring

### `resolver/resolver.go` — add 5 fields:
- NutritionProductService
- NutritionTemplateService
- NutritionTemplateItemService
- DailyNutritionOverrideService
- NutritionMacroService

### `atlas-gqlgen.yml` — add model bindings (20+ entries):
```
NutritionProduct → models.NutritionProduct
CreateProductInput → models.CreateProductInput
...etc all types, inputs, result unions, error types
Operation → models.Operation
NutritionErrorCode → models.NutritionErrorCode
```

### `cmd/server/main.go` — add wiring:
```go
atlasNutritionProductRepo := atlasPostgres.NewNutritionProductRepository(db.Pool)
atlasNutritionTemplateRepo := atlasPostgres.NewNutritionTemplateRepository(db.Pool)
atlasNutritionTemplateItemRepo := atlasPostgres.NewNutritionTemplateItemRepository(db.Pool)
atlasNutritionOverrideRepo := atlasPostgres.NewDailyNutritionOverrideRepository(db.Pool)
atlasNutritionOverrideItemRepo := atlasPostgres.NewDailyNutritionOverrideItemRepository(db.Pool)

atlasNutritionProductService := atlasService.NewNutritionProductService(atlasNutritionProductRepo)
atlasNutritionTemplateService := atlasService.NewNutritionTemplateService(atlasNutritionTemplateRepo, atlasNutritionTemplateItemRepo)
atlasNutritionTemplateItemService := atlasService.NewNutritionTemplateItemService(atlasNutritionTemplateItemRepo)
atlasNutritionOverrideService := atlasService.NewNutritionOverrideService(atlasNutritionOverrideRepo, atlasNutritionOverrideItemRepo)
atlasNutritionMacroService := atlasService.NewNutritionMacroService(
    atlasNutritionTemplateRepo, atlasNutritionTemplateItemRepo,
    atlasNutritionOverrideRepo, atlasNutritionOverrideItemRepo,
    atlasNutritionProductRepo,
)
```

## Implementation Order (8 Slices)

| # | Slice | Files |
|---|-------|-------|
| 1 | Migration | `00090_nutrition_tables.sql` |
| 2 | sqlc queries | 5 `.sql` files |
| 3 | Models | `models/nutrition.go` |
| 4 | Repositories | 5 repo files |
| 5 | Services | 5 service files |
| 6 | GraphQL schema | `schema/nutrition.graphql` |
| 7 | Resolvers | `resolver/nutrition.go` |
| 8 | Wiring | `atlas-gqlgen.yml`, `resolver/resolver.go`, `main.go` |

After codegen (`bunx nx run api:codegen`), resolver stubs are generated automatically.

## Testing (30 tests)

- **Unit**: repo tests with sqlc mock or test DB
- **Unit**: service validation tests (name required, macros >= 0, amount > 0)
- **Unit**: macro calculation tests (template only, with overrides, soft-deleted products, empty template)
- **Integration**: resolver round-trip (create product → create template → add items → create override → verify macros)
- **Integration**: soft-delete isolation (deleted product hidden from list, visible by ID)
- **Integration**: cascade delete (delete template removes items)
- **Integration**: auth guard (all operations return AuthError without PIN session)

## Log Markers

- `[NutritionProduct][create|update|delete|get|list]`
- `[NutritionTemplate][create|update|delete|get|list|current]`
- `[NutritionTemplateItem][create|update|delete]`
- `[DailyNutritionOverride][create|update|delete|get|list]`
- `[DailyNutritionOverrideItem][create|update|delete]`
- `[NutritionMacros][calculate]`

No sensitive content logged (no notes, meal labels).