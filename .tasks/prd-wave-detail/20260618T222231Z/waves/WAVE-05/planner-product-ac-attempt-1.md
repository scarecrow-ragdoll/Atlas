# WAVE-05 Product-AC Planner Attempt 1

## Sources Read
- docs/prd-waves/waves/wave-05.md
- docs/product-verified/acceptance-criteria.md (AC-017–AC-019, AC-058–AC-064, AC-113)
- docs/product-verified/edge-cases.md (EDGE-003, EDGE-009, EDGE-017, EDGE-019)
- docs/product-verified/business-rules.md (RULE-006, RULE-010, RULE-011, RULE-018, RULE-019, RULE-020)
- docs/product-verified/domain-model.md (NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, DailyNutritionOverrideItem)
- docs/product-verified/functional-spec.md (§15 REQ-010/REQ-011)
- docs/product-verified/user-flows.md (§26.7 Create Nutrition Template, §26.8 Override Daily Nutrition)
- docs/product-verified/features/nutrition.md
- docs/prd-waves/frontend-pages/page-007.md

## Selected Backend Wave Boundary
WAVE-05 covers the complete Nutrition domain: product catalog, weekly template, template items, daily overrides, override items, and KJBJU macro calculation. All CRUD via GraphQL. No REST endpoints needed (no binary uploads). No barcode scanner, food database, or advanced macro tracking.

## Neighboring Backend Wave Fit
- WAVE-01 (Foundation): prerequisite — PIN auth middleware, migration infrastructure, GraphQL common types, config extension pattern. WAVE-05 cannot start until WAVE-01 provides PIN auth contract.
- WAVE-02 (Exercise Library): no dependency — fully parallelizable
- WAVE-03 (Workout Diary): no dependency — fully parallelizable
- WAVE-04 (Cardio/Body): no dependency — fully parallelizable (note: WAVE-04 creates daily_log table, but nutrition does not reference daily_log)
- WAVE-06 (Charts): WAVE-05 provides nutrition data for chart queries (weekly KJBJU averages)
- WAVE-07/08 (AI Export/Review): WAVE-05 provides nutrition template and override data via service layer
- WAVE-09 (Backup): WAVE-05 tables designed for JSON-serializable export compatibility

## Frontend Pages Context
- PAGE-007 (Nutrition): primary frontend consumer — depends on all NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, DailyNutritionOverrideItem GraphQL queries and mutations. Macro calculation endpoint for display. Dependency context only.

## Codebase Evidence
No nutrition-related code exists. All code paths are new. Patterns follow WAVE-02/WAVE-04 established conventions: goose migration → sqlc queries → repository adapter → service → GraphQL schema → resolver.

## Proposed Details

### NutritionProduct CRUD
- Fields: name (required, string), caloriesPer100g (required, float), proteinPer100g (required, float), fatPer100g (required, float), carbsPer100g (required, float), notes (optional, string)
- Per-100g values must accept 0 but reject negative (EDGE-003 resolution)
- userId FK for scoping (MVP: single default user)
- Soft delete isActive flag per product? PRD does not specify deletion behavior for products. EDGE-019 warns about deleting referenced products. Decision needed.
- Deletion: if product referenced in active template/override items, either (a) block deletion with error, (b) cascade with warning, or (c) soft-delete (isActive flag). Product ACs don't specify. Recommend soft-delete (isActive flag) to preserve referential integrity for historical data, with query filtering active=true by default.

### NutritionTemplate CRUD
- Fields: weekStartDate (required, date — Monday anchor), title (optional, string), notes (optional, string)
- One active template per week per user (RULE-020). But "at a time" semantics unclear: does creating a new template for the same week replace the old one, or is there only ever one template active?
  - Decision (DDEC-W05-001): Create replaces. Upsert by (userId, weekStartDate). Creating a template for week X replaces any existing template for that week.
- EDGE-017 (mid-week template): template applies only forward? PRD §15.4 says template applies to all week days. Decision (DDEC-W05-002): Template applies to all days of the week, including past days within the same week, consistent with RULE-018. Frontend may choose to display differently.
- EDGE-009 (empty template): should be allowed — user may create template structure first, add items later.

### NutritionTemplateItem
- Fields: templateId (FK → NutritionTemplate), productId (FK → NutritionProduct), amountGrams (required, positive float), mealLabel (optional, enum or free string), notes (optional, string)
- EDGE-019 (delete referenced product): if product deleted (soft-delete with isActive flag), template items referencing it remain in DB but are excluded from KJBJU calculation with a warning marker.
- amountGrams > 0 validation (EDGE-003 extension)
- mealLabel: free-text string, not enum (more flexible for user)

### DailyNutritionOverride CRUD
- Fields: date (required, date), notes (optional, string)
- One override per date per user (RULE-019: override affects only one date)
- Unique constraint: (userId, date)
- AC-113: override must not affect other dates

### DailyNutritionOverrideItem
- Fields: overrideId (FK → DailyNutritionOverride), productId (FK → NutritionProduct), amountGrams (required, positive float), operation (enum: add/subtract/replace), mealLabel (optional, string), notes (optional, string)
- operation enum: add (add to template), subtract (subtract from template), replace (replace template amount)

### Macro Calculation (KJBJU)
- RULE-010: KJBJU = Σ(product per-100g values × grams / 100)
- Template day calculation: sum over all template items for the week
- Override day calculation: start with template values, apply override operations
  - ADD: add (product.per100g × amountGrams / 100) to totals
  - SUBTRACT: subtract (product.per100g × amountGrams / 100) from totals
  - REPLACE: for the given product, replace template value with override value = product.per100g × amountGrams / 100
- RULE-011: recalculate on any override change
- Return all 4 macros: calories, protein, fat, carbs

## Acceptance Criteria Contributions

| Proposed AC ID | Description |
| --- | --- |
| AC-W05-001 | NutritionProduct can be created via GraphQL mutation with name (required), caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g (all required, float >= 0), optional notes |
| AC-W05-002 | NutritionProduct nutritional values accept 0 (e.g., water has 0 calories, 0 protein), but reject negative values |
| AC-W05-003 | NutritionProduct can be read by ID |
| AC-W05-004 | NutritionProducts can be listed (all products for user) |
| AC-W05-005 | NutritionProduct can be updated (name, any nutritional value, notes) |
| AC-W05-006 | NutritionProduct can be soft-deleted (isActive flag set to false). Hard deletion blocked when referenced by active template/override items |
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
| AC-W05-032 | Soft-deleted products excluded from KJBJU calculation — template/override items referencing them return 0 for that item's macros with warning marker |
| AC-W05-033 | Empty template (zero items) returns 0 for all macros for all days (EDGE-009) |
| AC-W05-034 | All WAVE-05 GraphQL mutations return AuthError when PIN session header is missing or invalid |

## Exit Criteria Contributions

| Proposed EC ID | Description |
| --- | --- |
| EC-W05-001 | AC-W05-001 through AC-W05-034 pass via TEST-W05-001 through TEST-W05-030 |
| EC-W05-002 | gqlgen codegen produces valid Go code for WAVE-05 schema without drift |
| EC-W05-003 | sqlc codegen produces valid Go code for WAVE-05 queries without drift |
| EC-W05-004 | WAVE-01 PIN auth guard protects all WAVE-05 GraphQL endpoints |
| EC-W05-005 | WAVE-01 admin auth and health test suite still passes after WAVE-05 changes |
| EC-W05-006 | All migrations apply and roll back in sequence without errors |
| EC-W05-007 | NutritionProduct values >= 0 validation enforced for all 4 nutritional fields |
| EC-W05-008 | Template item amountGrams > 0 validation enforced |
| EC-W05-009 | Override item amountGrams > 0 validation enforced and operation enum validated |
| EC-W05-010 | Nutrition round-trip integration test passes (create product → create template → add items → create override → verify macros) |
| EC-W05-011 | Lint passes for all changed packages |
| EC-W05-012 | No sensitive content (nutrition product names, notes) in application logs |

## Verification Contributions

| Proposed Test ID | Description | Type |
| --- | --- | --- |
| TEST-W05-001 | NutritionProduct repository CRUD unit tests | unit |
| TEST-W05-002 | NutritionProduct service validation (name required, values >= 0) | unit |
| TEST-W05-003 | NutritionProduct GraphQL resolver integration tests | integration |
| TEST-W05-004 | NutritionTemplate repository CRUD unit tests | unit |
| TEST-W05-005 | NutritionTemplate service validation (weekStartDate required, upsert contract) | unit |
| TEST-W05-006 | NutritionTemplate upsert (replaces existing template for same week) | integration |
| TEST-W05-007 | NutritionTemplate GraphQL resolver integration tests | integration |
| TEST-W05-008 | NutritionTemplateItem repository CRUD unit tests | unit |
| TEST-W05-009 | NutritionTemplateItem validation (productId FK, amountGrams > 0) | unit |
| TEST-W05-010 | DailyNutritionOverride repository CRUD unit tests | unit |
| TEST-W05-011 | DailyNutritionOverride service validation (unique per date) | unit |
| TEST-W05-012 | DailyNutritionOverride GraphQL resolver integration tests | integration |
| TEST-W05-013 | DailyNutritionOverrideItem repository CRUD unit tests | unit |
| TEST-W05-014 | DailyNutritionOverrideItem validation (operation enum, amountGrams > 0) | unit |
| TEST-W05-015 | Macro calculation for template week (all 4 macros per day) | unit |
| TEST-W05-016 | Macro calculation with override operations (add/subtract/replace) | unit |
| TEST-W05-017 | Macro calculation with soft-deleted products (returns 0 for that item) | unit |
| TEST-W05-018 | Macro calculation with empty template (returns 0) | integration |
| TEST-W05-019 | Override isolation (override only affects target date) | integration |
| TEST-W05-020 | Nutrition round-trip integration test | integration |
| TEST-W05-021 | All WAVE-05 GraphQL operations return AuthError without PIN session | integration |
| TEST-W05-022 | Migration smoke test | integration |
| TEST-W05-023 | Codegen drift check (gqlgen + sqlc) | codegen |
| TEST-W05-024 | Log privacy: no product details in application logs | unit |
| TEST-W05-025 | Go lint for API package | lint |
| TEST-W05-026 | GraphQL schema validate | codegen |
| TEST-W05-027 | Soft-delete product blocked when referenced in active template | integration |
| TEST-W05-028 | Template cascade delete (deleting template deletes its items) | integration |
| TEST-W05-029 | Override cascade delete (deleting override deletes its items) | integration |

## Risks And Rollback
- KJBJU calculation is in-app (Go service layer). No persistent macro storage — always calculated from raw data. Risk: performance with many items. Mitigation: MVP has limited items, calculation is O(n).
- Soft-delete for NutritionProduct introduces isActive flag complexity. Rollback: revert to hard delete if referential integrity concerns are manageable.
- Template upsert is destructive (replaces old template). User should be warned before replacement. No undo after commit.

## Questions Raised
- DQ-W05-001: Should NutritionProduct use hard or soft delete? Soft-delete recommended for referential integrity with historical template/override data. Hard delete blocked when referenced.
- DQ-W05-002: What is the exact "single template at a time" semantic (RULE-020)? Per-week or per-user? Recommended: per-week. Creating a new template for week X replaces the old one.
- DQ-W05-003: mealLabel — free text or enum? Recommended: free-text string for maximum flexibility.

## Traceability Candidates
- docs/product-verified/features/nutrition.md → AC-W05-001–AC-W05-034
- docs/product-verified/acceptance-criteria.md (AC-017–AC-019, AC-058–AC-064, AC-113) → AC-W05-001, AC-W05-007, AC-W05-018, AC-W05-028, AC-W05-029, AC-W05-030
- docs/product-verified/edge-cases.md (EDGE-003) → AC-W05-002, AC-W05-015, AC-W05-024
- docs/product-verified/edge-cases.md (EDGE-009) → AC-W05-033
- docs/product-verified/edge-cases.md (EDGE-017) → DDEC-W05-002
- docs/product-verified/edge-cases.md (EDGE-019) → AC-W05-006, AC-W05-032
- docs/product-verified/business-rules.md (RULE-006) → AC-W05-001, AC-W05-002
- docs/product-verified/business-rules.md (RULE-010) → AC-W05-029, AC-W05-030
- docs/product-verified/business-rules.md (RULE-011) → AC-W05-030
- docs/product-verified/business-rules.md (RULE-018) → AC-W05-029
- docs/product-verified/business-rules.md (RULE-019) → AC-W05-028
- docs/product-verified/business-rules.md (RULE-020) → AC-W05-008, DDEC-W05-001