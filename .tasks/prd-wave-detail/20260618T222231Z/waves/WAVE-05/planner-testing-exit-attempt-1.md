# WAVE-05 Testing-Exit Planner Attempt 1

## Sources Read
- docs/verification-plan.xml (V-M-API, V-M-COVERAGE-GATE sections)
- docs/technical-verified/testing-and-delivery.md
- docs/prd-wave-details/waves/wave-04.md (Verification Obligations section)
- docs/prd-wave-details/waves/wave-01.md (Verification Obligations section)

## Selected Backend Wave Boundary
5 new entities + macro calculation. No REST endpoints. No binary/media operations. All operations are GraphQL CRUD + calculation.

## Neighboring Backend Wave Fit
Tests are independent from other waves. No shared test fixtures needed beyond WAVE-01 test infrastructure (test DB, test PIN auth helpers).

## Frontend Pages Context
Frontend tests are out of scope. Backend tests only.

## Codebase Evidence
Existing test patterns in apps/api: Go tests with testify/suite, sqlc-generated queries tested via repository adapter tests, service tests with mocked repos, resolver integration tests with real DB.

## Proposed Details

### Exit Criteria

| EC ID | Description |
| --- | --- |
| EC-W05-001 | AC-W05-001 through AC-W05-034 pass via TEST-W05-001 through TEST-W05-030 |
| EC-W05-002 | gqlgen codegen produces valid Go code for WAVE-05 schema without drift |
| EC-W05-003 | sqlc codegen produces valid Go code for WAVE-05 queries without drift |
| EC-W05-004 | WAVE-01 PIN auth guard protects all WAVE-05 GraphQL endpoints |
| EC-W05-005 | WAVE-01 admin auth and health test suite still passes after WAVE-05 changes |
| EC-W05-006 | Migration 00081 applies and rolls back in sequence without errors |
| EC-W05-007 | NutritionProduct values >= 0 validation enforced for all 4 nutritional fields |
| EC-W05-008 | Template item amountGrams > 0 validation enforced |
| EC-W05-009 | Override item amountGrams > 0 validation enforced and operation enum validated |
| EC-W05-010 | Nutrition round-trip integration test passes (create product → create template → add items → create override → verify macros) |
| EC-W05-011 | Lint passes for all changed packages |
| EC-W05-012 | No sensitive content (product notes) in application logs |

### Verification Obligations

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
| TEST-W05-022 | Migration smoke test (00081 up + down) | integration | bunx nx run api:test -- --run '(?i)migration_wave05' |
| TEST-W05-023 | Codegen drift check (gqlgen + sqlc) | codegen | bunx nx run api:codegen && bunx nx run graphql:codegen |
| TEST-W05-024 | Log privacy: no product notes in application logs | unit | bunx nx run api:test -- --run '(?i)wave05_log_sanitize' |
| TEST-W05-025 | Go lint for API package | lint | bunx nx run api:lint |
| TEST-W05-026 | GraphQL schema validate | codegen | bunx nx run graphql:validate |
| TEST-W05-027 | Soft-delete product blocked when referenced in active template | integration | bunx nx run api:test -- --run '(?i)nutrition_product_delete_blocked' |
| TEST-W05-028 | Template cascade delete (deleting template deletes its items) | integration | bunx nx run api:test -- --run '(?i)nutrition_template_cascade' |
| TEST-W05-029 | Override cascade delete (deleting override deletes its items) | integration | bunx nx run api:test -- --run '(?i)nutrition_override_cascade' |
| TEST-W05-030 | WAVE-01 admin auth regression tests | unit | bunx nx run api:test -- --run '(?i)admin_auth' |

## Risks And Rollback
- Template upsert integration tests must carefully set up and tear down to avoid cross-test interference
- Macro calculation tests are pure unit tests (no DB) — easy to write and fast

## Questions Raised
- DQ-W05-008: Should macro calculation have a dedicated integration test with real DB data or is unit test with mocked data sufficient? Recommended: unit tests for calculation logic (fast), integration test for round-trip (TEST-W05-020).

## Traceability Candidates
- docs/verification-plan.xml → V-M-API check patterns
- docs/prd-wave-details/waves/wave-04.md → test pattern reference