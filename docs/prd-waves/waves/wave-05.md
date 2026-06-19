# Wave 05: Nutrition

## Status

user-approved

## User Approval

user-approved (2026-06-18)

## Purpose

Food tracking through weekly template and daily overrides with macro calculations.

## Outcome After Wave

- OUT-W05-001 NutritionProduct CRUD
- OUT-W05-002 NutritionTemplate CRUD
- OUT-W05-003 NutritionTemplateItem for daily amounts
- OUT-W05-004 DailyNutritionOverride for per-day changes
- OUT-W05-005 Macro calculations (calories, protein, fat, carbs)
- OUT-W05-006 Ready for nutrition charts

## Included Scope

- CAP-W05-001 NutritionProduct CRUD
- CAP-W05-002 NutritionTemplate CRUD with week start
- CAP-W05-003 NutritionTemplateItem (product, amount, meal, notes)
- CAP-W05-004 DailyNutritionOverride CRUD
- CAP-W05-005 DailyNutritionOverrideItem (add/subtract/replace)
- CAP-W05-006 Macro calculations

## Excluded Scope

- Barcode scanner
- Food database
- Advanced macro tracking

## Dependencies

WAVE-01

## Surface Categories

backend, data, operations

## Risk Class

Low - Standard CRUD with calculations

## Recommended Next Planning

$detail-prd-wave for WAVE-05

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |

## Traceability

- docs/product/prd.md Section 15
- docs/product-verified/domain-model.md#NutritionProduct, #NutritionTemplate, #DailyNutritionOverride