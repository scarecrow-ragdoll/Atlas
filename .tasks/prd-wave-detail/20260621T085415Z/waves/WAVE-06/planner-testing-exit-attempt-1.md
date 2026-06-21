# WAVE-06 Testing-Exit Planner Attempt 1

## Sources Read
- planner-product-ac-attempt-1.md
- planner-architecture-codebase-attempt-1.md
- planner-data-integration-ops-attempt-1.md
- planner-security-compliance-attempt-1.md
- docs/prd-wave-details/waves/wave-05.md — TEST section as pattern
- docs/prd-wave-details/waves/wave-04.md — TEST section as pattern

## Selected Backend Wave Boundary
Read-only query aggregation. Tests must cover body chart queries, nutrition chart queries, measurement queries, auth, and empty-data edge cases.

## Neighboring Backend Wave Fit
Same pattern as WAVE-04/WAVE-05: unit tests for stateless computation + integration tests for DB-backed queries + auth tests for PIN guard.

## Frontend Pages Context
No frontend test implications. Backend queries should return well-typed, predictable structures.

## Codebase Evidence
- Existing test patterns: `bunx nx run api:test -- --run '(?i)pattern'` 
- Service tests use repository interfaces (mocks or in-memory implementations)
- Resolver tests use full resolver setup with auth context
- Migration tests exist for DB schema verification

## Proposed Details

### Test Strategy
1. Unit tests: e1RM calculation (Epley formula), date range validation, nutrition weekly average math
2. Integration tests: body weight trend, measurement trend, measurement overlay, nutrition weekly averages against test DB
3. Auth tests: all chart queries return AuthError without PIN session
4. Edge case tests: empty data periods, single data points, large date ranges

### Exit Criteria
| EC ID | Description |
|---|---|
| EC-W06-001 | AC-W06-001 through AC-W06-015 pass via TEST-W06-001 through TEST-W06-020 |
| EC-W06-002 | All chart queries protected by WAVE-01 PIN auth middleware |
| EC-W06-003 | gqlgen codegen produces valid Go code for WAVE-06 schema without drift |
| EC-W06-004 | Empty series returned for no-data periods (no errors) — verified for body, measurement, and nutrition queries |
| EC-W06-005 | Nutrition weekly average calculation matches RULE-015 |
| EC-W06-006 | Lint passes for all changed packages |
| EC-W06-007 | No sensitive data (body weight values, measurement values) in application logs |
| EC-W06-008 | Measurement overlay returns correct multi-type structure |

### Verification Obligations
| Test ID | Description | Type | Command |
|---|---|---|---|
| TEST-W06-001 | Epley e1RM formula calculation unit tests | unit | bunx nx run api:test -- --run '(?i)epley' |
| TEST-W06-002 | Body weight trend query returns correct time-series data | integration | bunx nx run api:test -- --run '(?i)body_weight_trend' |
| TEST-W06-003 | Body weight trend returns empty series for no data | integration | bunx nx run api:test -- --run '(?i)body_weight_trend_empty' |
| TEST-W06-004 | Measurement trend query returns correct data for single type | integration | bunx nx run api:test -- --run '(?i)measurement_trend' |
| TEST-W06-005 | Measurement overlay query returns multiple measurement types | integration | bunx nx run api:test -- --run '(?i)measurement_overlay' |
| TEST-W06-006 | Measurement queries return empty series for no data | integration | bunx nx run api:test -- --run '(?i)measurement_trend_empty' |
| TEST-W06-007 | Nutrition weekly averages query returns correct RULE-015 values | integration | bunx nx run api:test -- --run '(?i)nutrition_weekly_avg' |
| TEST-W06-008 | Nutrition weekly averages returns empty series for no data | integration | bunx nx run api:test -- --run '(?i)nutrition_weekly_avg_empty' |
| TEST-W06-009 | Date range validation (from > to returns error) | unit | bunx nx run api:test -- --run '(?i)chart_date_validation' |
| TEST-W06-010 | Default period when no date range specified (last 12 weeks) | integration | bunx nx run api:test -- --run '(?i)chart_default_period' |
| TEST-W06-011 | All chart queries return AuthError without PIN session | integration | bunx nx run api:test -- --run '(?i)wave06_auth' |
| TEST-W06-012 | Codegen drift check (gqlgen + sqlc) | codegen | bunx nx run api:codegen && bunx nx run graphql:codegen |
| TEST-W06-013 | Go lint for API package | lint | bunx nx run api:lint |
| TEST-W06-014 | GraphQL schema validate | codegen | bunx nx run graphql:validate |
| TEST-W06-015 | Log privacy: no body weight or measurement values in logs | unit | bunx nx run api:test -- --run '(?i)wave06_log_sanitize' |
| TEST-W06-016 | Max date range enforcement (52 weeks) | integration | bunx nx run api:test -- --run '(?i)chart_max_range' |
| TEST-W06-017 | Measurement trend with side filter (LEFT, RIGHT, NONE) | integration | bunx nx run api:test -- --run '(?i)measurement_side_trend' |
| TEST-W06-018 | Body weight trend with single data point | integration | bunx nx run api:test -- --run '(?i)body_weight_single_point' |
| TEST-W06-019 | Nutrition weekly average across partial week (mid-week start) | integration | bunx nx run api:test -- --run '(?i)nutrition_partial_week' |
| TEST-W06-020 | Codegen drift check after schema changes | codegen | bunx nx run api:codegen && bunx nx run graphql:codegen |

## Risks And Rollback
- Exercise chart tests cannot be written until WAVE-03 provides the workout_sets model
- Measurement range test depends on new sqlc query — integration test valid after query is implemented
- Nutrition weekly average test depends on existing macro service behavior — should be stable

## Questions Raised
- DQ-W06-008: Should measurement trend include null-side measurements (unpaired types) in overlay? (Proposed: yes — side = NONE)

## Traceability Candidates
- docs/prd-wave-details/waves/wave-04.md (TEST section pattern)
- docs/prd-wave-details/waves/wave-05.md (TEST section pattern)
- docs/verification-plan.xml (GRACE verification contracts)