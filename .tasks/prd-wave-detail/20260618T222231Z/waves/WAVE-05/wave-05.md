# WAVE-05: Nutrition

## Status
ready-for-dev

## User Approval
user-approved (2026-06-18). Source wave from docs/prd-waves/waves/wave-05.md.

## Source Wave Summary
WAVE-05 from docs/prd-waves/waves/wave-05.md. Food tracking through weekly template and daily overrides with macro calculations. Source status: user-approved (2026-06-18).

## Outcome After Implementation
- OUT-W05-001: NutritionProduct CRUD with name, KJBJU per 100g, optional notes, soft-delete
- OUT-W05-002: NutritionTemplate CRUD with weekStartDate, upsert semantics (one per week)
- OUT-W05-003: NutritionTemplateItem CRUD (product, amountGrams, mealLabel, notes)
- OUT-W05-004: DailyNutritionOverride CRUD (one per date, unique per user)
- OUT-W05-005: DailyNutritionOverrideItem CRUD (add/subtract/replace operations)
- OUT-W05-006: KJBJU macro calculation (calories, protein, fat, carbs) per RULE-010

## Scope Included
- CAP-W05-001: NutritionProduct CRUD via GraphQL (name, calories/protein/fat/carbs per 100g, soft-delete via isActive flag)
- CAP-W05-002: NutritionTemplate CRUD via GraphQL with weekStartDate, title, notes. Upsert by userId+weekStartDate (one template per week)
- CAP-W05-003: NutritionTemplateItem CRUD via GraphQL (productId, amountGrams, optional mealLabel, optional notes). Nested under template.
- CAP-W05-004: DailyNutritionOverride CRUD via GraphQL (date, optional notes). Unique per date per user.
- CAP-W05-005: DailyNutritionOverrideItem CRUD via GraphQL (productId, amountGrams, operation enum add/subtract/replace, optional mealLabel, optional notes). Nested under override.
- CAP-W05-006: Macro calculation query (nutritionMacros) returning calories, protein, fat, carbs per day. Server-side calculation per RULE-010 and RULE-011.

## Scope Excluded
- Barcode scanner
- Food database
- Advanced macro tracking (fiber, sugar, sodium, etc.)
- Recipes or meal planning beyond simple template
- Nutrition charts (WAVE-06)
- Frontend pages, UI, UX, routes, navigation, components

## Dependencies And Other-Wave Fit
- WAVE-01 (Foundation): prerequisite — provides PIN auth middleware, Atlas GraphQL endpoint, migration infrastructure, gqlgen config, sqlc config, atlas_users table, bootstrap service. WAVE-05 cannot start until WAVE-01 provides these contracts.
- WAVE-02 (Exercise Library): no direct dependency — can fully parallelize
- WAVE-03 (Workout Diary): no direct dependency — can fully parallelize. Nutrition tables are independent from daily_log.
- WAVE-04 (Cardio and Body Tracking): no direct dependency — can fully parallelize. Migration number collision is the only coordination risk (see DQ-W05-009).
- WAVE-06 (Charts): WAVE-05 provides nutrition macro data for weekly KJBJU average chart queries
- WAVE-07 (AI Export): WAVE-05 provides nutrition template and override data via service layer (JSON-serializable)
- WAVE-08 (AI Review): no dependency
- WAVE-09 (Backup): WAVE-05 tables are JSON-serializable for export compatibility

## Frontend Pages Dependencies
- PAGE-007 (Nutrition): primary frontend consumer — depends on all NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, DailyNutritionOverrideItem GraphQL queries and mutations. Also depends on nutritionMacros query for macro summary display.
- PAGE-001 (Dashboard): depends on nutrition macro summary data (deferred to WAVE-06 charts)
- Dependency context only; no frontend pages, UI, or UX work in this wave.

## Codebase Fit And Touchpoints
- apps/api/internal/atlas/graph/schema/nutrition.graphql: new schema file
- apps/api/internal/atlas/models/nutrition.go: new model types (NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, DailyNutritionOverrideItem, NutritionMacros + input/result/error types)
- apps/api/internal/atlas/service/nutrition_product_service.go: transport-neutral service with validation
- apps/api/internal/atlas/service/nutrition_template_service.go: transport-neutral service with upsert logic and cascade
- apps/api/internal/atlas/service/nutrition_override_service.go: transport-neutral service with override isolation
- apps/api/internal/atlas/service/nutrition_macro_service.go: stateless KJBJU calculation service
- apps/api/internal/atlas/repository/postgres/nutrition_product_repo.go: repository adapter
- apps/api/internal/atlas/repository/postgres/nutrition_template_repo.go: repository adapter
- apps/api/internal/atlas/repository/postgres/nutrition_template_item_repo.go: repository adapter
- apps/api/internal/atlas/repository/postgres/nutrition_override_repo.go: repository adapter
- apps/api/internal/atlas/repository/postgres/nutrition_override_item_repo.go: repository adapter
- apps/api/internal/atlas/graph/resolver/nutrition.go: GraphQL resolvers for all nutrition CRUD and macro calculation
- apps/api/internal/atlas/graph/resolver/resolver.go: add new service fields
- apps/api/internal/repository/postgres/migrations/00081_nutrition_tables.sql: new migration (or next available number)
- apps/api/internal/repository/postgres/queries/nutrition_products.sql: sqlc query definitions
- apps/api/internal/repository/postgres/queries/nutrition_templates.sql: sqlc query definitions
- apps/api/internal/repository/postgres/queries/nutrition_template_items.sql: sqlc query definitions
- apps/api/internal/repository/postgres/queries/nutrition_overrides.sql: sqlc query definitions
- apps/api/internal/repository/postgres/queries/nutrition_override_items.sql: sqlc query definitions
- apps/api/atlas-gqlgen.yml: add model bindings for all new nutrition types
- apps/api/sqlc.yaml: auto-discovers new queries via glob (no change needed)
- apps/api/cmd/server/main.go: wire all 5 repos, 5 services, add to atlasRes

## Design Contracts
- Soft-delete: NutritionProduct uses isActive flag. Soft-deleted products excluded from default list queries but remain in template/override items for historical reference. Viewable by direct ID lookup. (DDEC-W05-003)
- Template upsert: INSERT ... ON CONFLICT (user_id, week_start_date) DO UPDATE. Creating a template for an existing week replaces the previous template. (DDEC-W05-001)
- Template cascade: deleting a NutritionTemplate cascades to its NutritionTemplateItems.
- Override cascade: deleting a DailyNutritionOverride cascades to its DailyNutritionOverrideItems.
- Override isolation: DailyNutritionOverride affects only its target date. Other dates use template values unchanged. (AC-113)
- Macro calculation: RULE-010 formula. Override operations: ADD = add override item values, SUBTRACT = subtract override item values, REPLACE = for a given product, use override item value instead of template item value. (DDEC-W05-002)
- PIN auth: WAVE-01 middleware guards all WAVE-05 GraphQL operations via /graphql/atlas endpoint
- mealLabel: free-text string (not enum)
- Operation enum: ADD, SUBTRACT, REPLACE
- Error format: union result types following existing pattern (Success | ValidationError | AuthError)

## Data API Integration And Operations

### Database Schema

**nutrition_product**
- id (UUID, PK), user_id (UUID, FK → atlas_users), name (VARCHAR NOT NULL), calories_per_100g (REAL NOT NULL), protein_per_100g (REAL NOT NULL), fat_per_100g (REAL NOT NULL), carbs_per_100g (REAL NOT NULL), notes (TEXT nullable), is_active (BOOLEAN NOT NULL DEFAULT true), created_at, updated_at
- Index: idx_nutrition_product_user (user_id)
- Constraint: CHECK (calories_per_100g >= 0 AND protein_per_100g >= 0 AND fat_per_100g >= 0 AND carbs_per_100g >= 0)

**nutrition_template**
- id (UUID, PK), user_id (UUID, FK → atlas_users), week_start_date (DATE NOT NULL), title (VARCHAR nullable), notes (TEXT nullable), created_at, updated_at
- Unique: (user_id, week_start_date)
- Index: idx_nutrition_template_week (user_id, week_start_date)

**nutrition_template_item**
- id (UUID, PK), template_id (UUID, FK → nutrition_template ON DELETE CASCADE), product_id (UUID, FK → nutrition_product), amount_grams (REAL NOT NULL), meal_label (VARCHAR nullable), notes (TEXT nullable), created_at, updated_at
- Index: idx_nutrition_template_item_template (template_id)
- Constraint: CHECK (amount_grams > 0)

**daily_nutrition_override**
- id (UUID, PK), user_id (UUID, FK → atlas_users), date (DATE NOT NULL), notes (TEXT nullable), created_at, updated_at
- Unique: (user_id, date)
- Index: idx_nutrition_override_date (user_id, date)

**daily_nutrition_override_item**
- id (UUID, PK), override_id (UUID, FK → daily_nutrition_override ON DELETE CASCADE), product_id (UUID, FK → nutrition_product), amount_grams (REAL NOT NULL), operation (VARCHAR NOT NULL), meal_label (VARCHAR nullable), notes (TEXT nullable), created_at, updated_at
- Index: idx_nutrition_override_item_override (override_id)
- Constraint: CHECK (amount_grams > 0), CHECK (operation IN ('add', 'subtract', 'replace'))

### GraphQL Operations
- NutritionProduct: nutritionProducts (list active), nutritionProduct (by ID), createNutritionProduct, updateNutritionProduct, deleteNutritionProduct (soft-delete)
- NutritionTemplate: nutritionTemplates (date range), nutritionTemplate (by ID), nutritionTemplateCurrent (by weekStartDate), createNutritionTemplate (upsert), updateNutritionTemplate, deleteNutritionTemplate (cascade)
- NutritionTemplateItem: createNutritionTemplateItem, updateNutritionTemplateItem, deleteNutritionTemplateItem (nested under template)
- DailyNutritionOverride: dailyNutritionOverrides (date range), dailyNutritionOverride (by ID), dailyNutritionOverrideByDate (by date), createDailyNutritionOverride, updateDailyNutritionOverride, deleteDailyNutritionOverride (cascade)
- DailyNutritionOverrideItem: createDailyNutritionOverrideItem, updateDailyNutritionOverrideItem, deleteDailyNutritionOverrideItem (nested under override)
- Macro: nutritionMacros(weekStartDate: Date!, date: Date) → NutritionMacrosResult
- Union results: NutritionProductResult, NutritionTemplateResult, DailyNutritionOverrideResult, NutritionMacrosResult — each union is Success | ValidationError | AuthError

### REST Endpoints
- None. All operations via GraphQL.

### Log Markers
- [NutritionProduct][create|update|delete|get|list]
- [NutritionTemplate][create|update|delete|get|list|current]
- [NutritionTemplateItem][create|update|delete]
- [DailyNutritionOverride][create|update|delete|get|list]
- [DailyNutritionOverrideItem][create|update|delete]
- [NutritionMacros][calculate]
- Sensitive data (product notes, meal labels) NOT logged. Nutritional values may be logged (non-sensitive).

### Operations
- PostgreSQL: goose migration (single migration, all 5 tables), reversible
- Existing Docker Compose stack, no new services
- No media storage needed
- No new config sections

## Security Privacy And Compliance
- All endpoints protected by WAVE-01 PIN auth middleware (GraphQL)
- When PIN disabled, endpoints accessible without auth (consistent with TDEC-037)
- All operations scoped to default user per MVP constraint — userID extracted from context via middleware.GetAtlasUserID(ctx)
- Soft-deleted products not deleted from DB, only hidden from default queries — data remains recoverable by admin
- Nutritional values (calories, protein, fat, carbs) are not sensitive — may be logged
- Product notes and meal labels: log privacy — log entity IDs only, not notes or labels
- No user PII in nutrition domain
- No external API calls (RULE-029)

## Implementation Slices

| Slice ID | Name | Description |
| --- | --- | --- |
| SLICE-W05-001 | DB migrations | Create goose migration for all 5 nutrition tables with indexes, FKs, cascades, and CHECK constraints |
| SLICE-W05-002 | sqlc queries | Define CRUD queries for all 5 entities: nutrition products (list active, by ID), templates (by week, current, date range), template items (by templateId), overrides (by date, range, by date), override items (by overrideId) |
| SLICE-W05-003 | Repository adapters | Implement 5 repo adapters with sqlc-generated code and error mapping (not found, constraint violations) |
| SLICE-W05-004 | Models | Define all model types: DB records, public models, inputs, result unions, error types for all 5 entities + macro calculation |
| SLICE-W05-005 | Services layer | Implement 5 services: product (validation, soft-delete), template (upsert, cascade), override (isolation), override-item, macro (KJBJU calculation with override operations) |
| SLICE-W05-006 | GraphQL schema | Add nutrition.graphql with types, enums (Operation, MacroSummary), inputs, queries, mutations, union results |
| SLICE-W05-007 | GraphQL resolvers | Implement nutrition resolvers with PIN auth guard and union error returns following existing settings resolver pattern |
| SLICE-W05-008 | gqlgen config and wiring | Add model bindings to atlas-gqlgen.yml, wire all 5 repos and 5 services in main.go, add to Resolver struct |

## Acceptance Criteria

| AC ID | Description |
| --- | --- |
| AC-W05-001 | NutritionProduct can be created via GraphQL mutation with name (required), caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g (all required, float >= 0), optional notes |
| AC-W05-002 | NutritionProduct nutritional values accept 0 but reject negative values — negative returns ValidationError |
| AC-W05-003 | NutritionProduct can be read by ID — returns full record |
| AC-W05-004 | NutritionProducts can be listed — returns all active (isActive=true) products for user |
| AC-W05-005 | NutritionProduct can be updated (name, any nutritional value, notes) — updated entry returned |
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
| AC-W05-025 | DailyNutritionOverrideItem operation validated against allowed enum values (add/subtract/replace). Invalid value returns ValidationError |
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
| AC-W05-036 | Soft-deleted NutritionProducts can still be viewed by ID (returns full record including isActive=false) |

## Exit Criteria

| EC ID | Description |
| --- | --- |
| EC-W05-001 | AC-W05-001 through AC-W05-036 pass via TEST-W05-001 through TEST-W05-030 |
| EC-W05-002 | gqlgen codegen produces valid Go code for WAVE-05 schema without drift |
| EC-W05-003 | sqlc codegen produces valid Go code for WAVE-05 queries without drift |
| EC-W05-004 | WAVE-01 PIN auth guard protects all WAVE-05 GraphQL endpoints. Existing admin auth unchanged. |
| EC-W05-005 | WAVE-01 admin auth and health test suite still passes after WAVE-05 changes |
| EC-W05-006 | Migration applies and rolls back in sequence without errors |
| EC-W05-007 | NutritionProduct values >= 0 validation enforced for all 4 nutritional fields |
| EC-W05-008 | Template item amountGrams > 0 validation enforced |
| EC-W05-009 | Override item amountGrams > 0 validation enforced and operation enum validated |
| EC-W05-010 | Nutrition round-trip integration test passes (create product → create template → add items → create override → verify macros) |
| EC-W05-011 | Lint passes for all changed packages |
| EC-W05-012 | No sensitive content (product notes, meal labels) in application logs |

## Verification Obligations

| Test ID | Description | Type | Command |
| --- | --- | --- | --- |
| TEST-W05-001 | NutritionProduct repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)nutrition_product_repo' |
| TEST-W05-002 | NutritionProduct service validation (name required, values >= 0) | unit | bunx nx run api:test -- --run '(?i)nutrition_product_service' |
| TEST-W05-003 | NutritionProduct GraphQL resolver integration tests | integration | bunx nx run api:test -- --run '(?i)nutrition_product_resolver' |
| TEST-W05-004 | NutritionTemplate repository CRUD + upsert unit tests | unit | bunx nx run api:test -- --run '(?i)nutrition_template_repo' |
| TEST-W05-005 | NutritionTemplate service validation (weekStartDate, upsert contract) | unit | bunx nx run api:test -- --run '(?i)nutrition_template_service' |
| TEST-W05-006 | NutritionTemplate upsert replaces existing template for same week | integration | bunx nx run api:test -- --run '(?i)nutrition_template_upsert' |
| TEST-W05-007 | NutritionTemplate GraphQL resolver integration tests | integration | bunx nx run api:test -- --run '(?i)nutrition_template_resolver' |
| TEST-W05-008 | NutritionTemplateItem repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)nutrition_template_item_repo' |
| TEST-W05-009 | NutritionTemplateItem validation (amountGrams > 0) | unit | bunx nx run api:test -- --run '(?i)nutrition_template_item_service' |
| TEST-W05-010 | DailyNutritionOverride repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)nutrition_override_repo' |
| TEST-W05-011 | DailyNutritionOverride service validation (unique per date) | unit | bunx nx run api:test -- --run '(?i)nutrition_override_service' |
| TEST-W05-012 | DailyNutritionOverride GraphQL resolver integration tests | integration | bunx nx run api:test -- --run '(?i)nutrition_override_resolver' |
| TEST-W05-013 | DailyNutritionOverrideItem repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)nutrition_override_item_repo' |
| TEST-W05-014 | DailyNutritionOverrideItem validation (operation enum, amountGrams > 0) | unit | bunx nx run api:test -- --run '(?i)nutrition_override_item_service' |
| TEST-W05-015 | Macro calculation for template week (all 4 macros per day) | unit | bunx nx run api:test -- --run '(?i)nutrition_macro_template' |
| TEST-W05-016 | Macro calculation with override operations (add/subtract/replace) | unit | bunx nx run api:test -- --run '(?i)nutrition_macro_override' |
| TEST-W05-017 | Macro calculation with soft-deleted products (returns 0 for that item) | unit | bunx nx run api:test -- --run '(?i)nutrition_macro_deleted_product' |
| TEST-W05-018 | Macro calculation with empty template (returns 0 all) | unit | bunx nx run api:test -- --run '(?i)nutrition_macro_empty' |
| TEST-W05-019 | Override isolation — override affects only its target date | integration | bunx nx run api:test -- --run '(?i)nutrition_override_isolation' |
| TEST-W05-020 | Nutrition round-trip integration test (full lifecycle) | integration | bunx nx run api:test -- --run '(?i)nutrition_roundtrip' |
| TEST-W05-021 | All WAVE-05 GraphQL operations return AuthError without PIN session | integration | bunx nx run api:test -- --run '(?i)wave05_auth' |
| TEST-W05-022 | Migration smoke test (up + down) | integration | bunx nx run api:test -- --run '(?i)migration_wave05' |
| TEST-W05-023 | Codegen drift check (gqlgen + sqlc) | codegen | bunx nx run api:codegen && bunx nx run graphql:codegen |
| TEST-W05-024 | Log privacy: no product notes or meal labels in application logs | unit | bunx nx run api:test -- --run '(?i)wave05_log_sanitize' |
| TEST-W05-025 | Go lint for API package | lint | bunx nx run api:lint |
| TEST-W05-026 | GraphQL schema validate | codegen | bunx nx run graphql:validate |
| TEST-W05-027 | Soft-delete product excluded from active products list | integration | bunx nx run api:test -- --run '(?i)nutrition_product_soft_delete_list' |
| TEST-W05-028 | Template cascade delete (deleting template deletes its items) | integration | bunx nx run api:test -- --run '(?i)nutrition_template_cascade' |
| TEST-W05-029 | Override cascade delete (deleting override deletes its items) | integration | bunx nx run api:test -- --run '(?i)nutrition_override_cascade' |
| TEST-W05-030 | WAVE-01 admin auth regression tests still pass | unit | bunx nx run api:test -- --run '(?i)admin_auth' |

## Rollout Rollback And Compatibility
- Rollout: merge PR, CI builds and runs tests, deploy via Dokploy compose update. New tables created via goose migration.
- Rollback: revert PR, CI builds previous image, Dokploy compose update rolls back. Run goose down migration.
- Compatibility: all new operations are additive. No existing API changes. WAVE-01 endpoints (health, PIN auth, settings, media) unchanged. WAVE-02/03/04 endpoints unchanged.
- Migration: goose migration runs at startup. Down migration available for rollback.
- Migration number: use next available number after WAVE-01 (currently 00080). Coordinate with WAVE-04 (may use 00081-00087). If WAVE-04 not yet deployed, WAVE-05 uses 00081.

## Handoff Packets
- HANDOFF-W05-001: This wave brief document
- HANDOFF-W05-002: Planner reports (6 scopes, 2 attempts for product-ac)
- HANDOFF-W05-003: Reviewer evidence (7 perspectives, 2 attempts for product-scope-and-ac and traceability-consistency)
- HANDOFF-W05-004: Consolidated question ledger

## Design Decisions

| DDEC ID | Decision | Rationale |
| --- | --- | --- |
| DDEC-W05-001 | Template upsert by (userId, weekStartDate) | RULE-020: one template per week. Creating a new template for week X replaces the previous one. |
| DDEC-W05-002 | Server-side macro calculation | Consistency guarantee. Separate query `nutritionMacros` for flexibility. |
| DDEC-W05-003 | Soft-delete for NutritionProduct with isActive flag | EDGE-019: preserves referential integrity for historical template/override data while allowing product catalog cleanup. |
| DDEC-W05-004 | Free-text mealLabel | Maximum user flexibility. Not constrained to enum values. |
| DDEC-W05-005 | Single migration for all 5 nutrition tables | Simpler than 5 separate migrations. Tables are independent and will be deployed together. |

## Reviewer Verdicts

| Wave | Perspective | Attempt | Verdict | Reviewer Report | Required Revisions | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| WAVE-05 | product-scope-and-ac | 1 | needs-revision | review-product-scope-and-ac-attempt-1.md | Soft-delete consistency, warning marker removal | ACs revised in attempt 2 |
| WAVE-05 | product-scope-and-ac | 2 | approved | review-product-scope-and-ac-attempt-2.md | none | All concerns addressed |
| WAVE-05 | architecture-codebase-fit | 1 | approved | review-architecture-codebase-fit-attempt-1.md | none | 8 slices, pattern B confirmed |
| WAVE-05 | data-api-integration-ops | 1 | approved | review-data-api-integration-ops-attempt-1.md | none | Clean schema design |
| WAVE-05 | security-privacy-compliance | 1 | approved | review-security-privacy-compliance-attempt-1.md | none | PIN auth, log privacy covered |
| WAVE-05 | testing-exit-criteria | 1 | approved | review-testing-exit-criteria-attempt-1.md | none | 30 tests, all AC/EC covered |
| WAVE-05 | sequencing-other-wave-fit | 1 | approved | review-sequencing-other-wave-fit-attempt-1.md | none | Clean dependency analysis |
| WAVE-05 | traceability-consistency | 1 | needs-revision | review-traceability-consistency-attempt-1.md | Cross-planner consistency, question consolidation | Addressed in attempt 2 |
| WAVE-05 | traceability-consistency | 2 | approved | review-traceability-consistency-attempt-2.md | none | All concerns addressed |

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W05-001 | WAVE-05 | data | resolved | EDGE-019 | Should NutritionProduct use soft-delete (isActive flag) or hard-delete with FK block? | Referential integrity for historical template/override data | Soft-delete with isActive flag. Products remain in DB but excluded from default queries. | planner-product-ac-attempt-2.md | resolved | Soft-delete (isActive flag) selected. DDEC-W05-003. |
| DQ-W05-002 | WAVE-05 | product-ac | resolved | RULE-020 | What is exact "single template at a time" semantic — per-week or per-user? | Drives upsert behavior | Per-week: template for week X replaces previous template for that week. | planner-product-ac-attempt-1.md | resolved | Per-week upsert. DDEC-W05-001. |
| DQ-W05-003 | WAVE-05 | product-ac | resolved | — | mealLabel — free text or enum? | Flexibility vs validation | Free-text string. | planner-product-ac-attempt-1.md | resolved | Free-text. DDEC-W05-004. |
| DQ-W05-004 | WAVE-05 | architecture | resolved | — | Macro calculation server-side or client-side? | Consistency | Server-side in Go service. | planner-architecture-codebase-attempt-1.md | resolved | Server-side. DDEC-W05-002. |
| DQ-W05-005 | WAVE-05 | data-ops | resolved | — | Macro query: separate or inline? | API surface design | Separate query. | planner-data-integration-ops-attempt-1.md | resolved | nutritionMacros query. |
| DQ-W05-006 | WAVE-05 | data-ops | resolved | — | NutritionProduct deletion: soft or hard? | Same as DQ-W05-001 | Soft-delete. | planner-data-integration-ops-attempt-1.md | resolved | Soft-delete. Merged to DQ-W05-001. |
| DQ-W05-007 | WAVE-05 | security | deferred | — | Soft-deleted products recoverable via API? | Data recovery | Admin-only DB recovery. No restore API in MVP. | planner-security-compliance-attempt-1.md | open | Admin-only for MVP. |
| DQ-W05-008 | WAVE-05 | testing | resolved | — | Macro tests: unit or integration? | Test scope | Both: unit (calculation) + integration (round-trip). | planner-testing-exit-attempt-1.md | resolved | Both types used. |
| DQ-W05-009 | WAVE-05 | sequencing | deferred | WAVE-04 | Migration number to use? | Avoid WAVE-04 collision | Check current state at implementation time. | planner-sequencing-fit-attempt-1.md | open | Coordinate with WAVE-04. Use next available. |

## Traceability
- docs/prd-waves/waves/wave-05.md: source wave boundary, outcomes, capability groups
- docs/product-verified/functional-spec.md: Nutrition §15 REQ-010/REQ-011
- docs/product-verified/domain-model.md: NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, DailyNutritionOverrideItem entities
- docs/product-verified/acceptance-criteria.md: AC-017–AC-019, AC-058–AC-064, AC-113
- docs/product-verified/edge-cases.md: EDGE-003, EDGE-009, EDGE-017, EDGE-019
- docs/product-verified/business-rules.md: RULE-006, RULE-010, RULE-011, RULE-018, RULE-019, RULE-020
- docs/product-verified/user-flows.md: §26.7 Create Nutrition Template, §26.8 Override Daily Nutrition
- docs/product-verified/actors-and-permissions.md: user permissions for nutrition
- docs/technical-verified/api-contracts.md: hybrid GraphQL/REST protocol, TDEC-027 error format
- docs/technical-verified/auth-security-compliance.md: PIN auth (TDEC-037)
- docs/technical-verified/data-contracts.md: domain entities, userId FKs
- docs/technical-verified/operations-observability.md: log markers, error format
- docs/development-plan.xml: M-API, M-PRD-WAVE-DETAILER module contracts
- docs/knowledge-graph.xml: existing module boundaries
- docs/prd-wave-details/waves/wave-01.md: WAVE-01 dependency contracts (PIN auth, migration infrastructure)
- docs/prd-wave-details/waves/wave-04.md: WAVE-04 patterns (repository, service, resolver structure)
- apps/api/internal/atlas: existing codebase patterns for service/repository/resolver structure
- docs/prd-waves/frontend-pages/page-007.md: nutrition page backend dependencies