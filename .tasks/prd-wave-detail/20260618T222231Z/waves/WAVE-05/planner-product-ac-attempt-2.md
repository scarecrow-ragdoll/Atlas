# WAVE-05 Product-AC Planner Attempt 2 (Revised)

## Revisions From Attempt 1
- Soft-delete (isActive flag) is now the consistent approach across all planners
- AC-W05-006 updated: "NutritionProduct can be soft-deleted (isActive flag set to false). Soft-deleted products are excluded from active product list queries but remain in template/override items for historical reference."
- AC-W05-032 updated: "Macro calculation with soft-deleted products: template/override items referencing soft-deleted products contribute 0 to all macros for that item (product is treated as having 0 values)."
- Removed "warning marker" references — no special error return needed
- Added AC-W05-035: "Soft-deleted NutritionProducts are not returned in the default (active) products list"

## Updated Acceptance Criteria

| Proposed AC ID | Description |
| --- | --- |
| AC-W05-001 | NutritionProduct can be created via GraphQL mutation with name (required), caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g (all required, float >= 0), optional notes |
| AC-W05-002 | NutritionProduct nutritional values accept 0 (e.g., water has 0 calories, 0 protein), but reject negative values |
| AC-W05-003 | NutritionProduct can be read by ID |
| AC-W05-004 | NutritionProducts can be listed (all active products for user, where isActive=true) |
| AC-W05-005 | NutritionProduct can be updated (name, any nutritional value, notes) |
| AC-W05-006 | NutritionProduct can be soft-deleted (isActive flag set to false). Soft-deleted products excluded from active list but remain in template/override items for historical reference |
| AC-W05-007 | NutritionTemplate can be created with weekStartDate (required, date), optional title, optional notes |
| AC-W05-008 | Creating a NutritionTemplate for an existing week replaces the previous template (upsert by userId+weekStartDate) |
| AC-W05-009 | NutritionTemplate can be read by ID with nested template items |
| AC-W05-010 | NutritionTemplates can be listed for a date range |
| AC-W05-011 | Current active template can be queried (nutritionTemplateCurrent for a given weekStartDate) |
| AC-W05-012 | NutritionTemplate can be updated (title, notes). weekStartDate is immutable after creation |
| AC-W05-013 | NutritionTemplate can be deleted (hard delete cascades to template items) |
| AC-W05-014 | NutritionTemplateItem can be created within a template with productId (required), amountGrams (required, > 0), optional mealLabel, optional notes |
| AC-W05-015 | NutritionTemplateItem amountGrams must be > 0. 0 or negative returns ValidationError |
| AC-W05-016 | NutritionTemplateItem can be updated (amountGrams, mealLabel, notes). productId is immutable after creation |
| AC-W05-017 | NutritionTemplateItem can be deleted |
| AC-W05-018 | DailyNutritionOverride can be created with date (required, unique per user), optional notes |
| AC-W05-019 | DailyNutritionOverride can be read by ID with nested override items |
| AC-W05-020 | DailyNutritionOverride can be listed by date range |
| AC-W05-021 | DailyNutritionOverride can be updated (notes). date is immutable after creation |
| AC-W05-022 | DailyNutritionOverride can be deleted (hard delete cascades to override items) |
| AC-W05-023 | DailyNutritionOverrideItem can be created with productId (required), amountGrams (required, > 0), operation (required, add/subtract/replace), optional mealLabel, optional notes |
| AC-W05-024 | DailyNutritionOverrideItem amountGrams must be > 0. 0 or negative returns ValidationError |
| AC-W05-025 | DailyNutritionOverrideItem operation validated against allowed enum values (add/subtract/replace) |
| AC-W05-026 | DailyNutritionOverrideItem can be updated (amountGrams, operation, mealLabel, notes) |
| AC-W05-027 | DailyNutritionOverrideItem can be deleted |
| AC-W05-028 | DailyNutritionOverride affects only its target date — does not modify template or other dates (AC-113) |
| AC-W05-029 | KJBJU calculation for a template week returns calories, protein, fat, carbs totals for each day per RULE-010 |
| AC-W05-030 | KJBJU calculation for an overridden day returns recalculated values reflecting override operations per RULE-011 |
| AC-W05-031 | Macro calculation for a day with no template items returns 0 for all macros |
| AC-W05-032 | Template/override items referencing soft-deleted products contribute 0 to all macros for that item |
| AC-W05-033 | Empty template (zero items) returns 0 for all macros for all days (EDGE-009) |
| AC-W05-034 | All WAVE-05 GraphQL mutations return AuthError when PIN session header is missing or invalid |
| AC-W05-035 | Soft-deleted NutritionProducts are not returned in default active products list query |
| AC-W05-036 | Soft-deleted NutritionProducts can still be viewed by ID (returns full record) |

## Traceability Additions
- AC-W05-001–AC-W05-006 → AC-058 (product KJBJU per 100g)
- AC-W05-007–AC-W05-017 → AC-059, AC-060 (template with items, meal labels)
- AC-W05-018–AC-W05-028 → AC-062, AC-113 (overrides, isolation)
- AC-W05-029–AC-W05-030 → AC-063, AC-064 (KJBJU calculation, recalculation)
- AC-W05-031 → AC-061 (template auto-applies to all days, showing 0 when empty)
- AC-W05-006 → EDGE-019 (delete referenced product)
- AC-W05-015, AC-W05-024 → EDGE-003 (0/negative values)
- AC-W05-033 → EDGE-009 (empty template)
- AC-W05-034 → PIN auth coverage