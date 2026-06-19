# PAGE-007: Nutrition

## Status

user-approved

## Page Purpose

Food tracking through weekly template and daily overrides. Products CRUD with macro nutrients.

## What Is On This Page

- Products list/CRUD
- Nutrition template for week
- Daily override view (selected date)
- Macro summary (calories, protein, fat, carbs)

## Functional Parts

- Product management section
- Weekly template editor
- Daily override editor
- Macro calculations

## Empty States

- No products - "Create your first food product"
- No template - "Create a weekly nutrition template"

## Loading And Error States

- Loading nutrition data - skeleton
- Save error - toast

## Backend Dependencies

- GET /api/nutrition-products
- POST /api/nutrition-products
- GET /api/nutrition-templates/current
- POST /api/nutrition-templates
- GET /api/daily-overrides?date=YYYY-MM-DD
- POST /api/daily-overrides

## Explicit Deferrals

- Barcode scanner - future scope
- Food database - future scope
- Water/calorie/salt tracking - future scope

## Open Questions

- None blocking

## Raw PRD Traceability

docs/product/prd.md Section 15

## Verified PRD Traceability

docs/product-verified/domain-model.md#NutritionProduct, #NutritionTemplate, #DailyNutritionOverride