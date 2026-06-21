# WAVE-06: Charts — Progress Visualization

## Status
ready-for-dev

## User Approval
user-approved (2026-06-18). Source wave from docs/prd-waves/waves/wave-06.md.

## Source Wave Summary
WAVE-06 from docs/prd-waves/waves/wave-06.md. Progress visualization for workouts, body measurements, and nutrition. Source status: user-approved (2026-06-18). Q-CHART-001 resolved: Epley formula selected by user on 2026-06-21.

## Outcome After Implementation
- OUT-W06-001: Exercise progress queries (working weight, best set, e1RM, volume, reps, working sets) — CONDITIONAL on WAVE-03
- OUT-W06-002: Body weight trend queries (date range)
- OUT-W06-003: Body measurement trend queries (single + overlay)
- OUT-W06-004: Nutrition weekly macro average queries
- OUT-W06-005: Period filtering on all chart queries

## Scope Included
- CAP-W06-001: Exercise progress query support — conditional on WAVE-03 workout_sets table existing
- CAP-W06-002: e1RM calculation — Epley formula: weight × (1 + reps / 30). Per-set calculation in-memory
- CAP-W06-003: Body weight trend queries (date range, ordered by date ascending)
- CAP-W06-004: Body measurement trend queries (single type) + measurement overlay (multiple types with alphabetical ordering)
- CAP-W06-005: Nutrition macro summary queries (weekly KJBJU averages per RULE-015)
- CAP-W06-006: Period filtering (optional from/to date parameters)
- CAP-W06-007: Chart data aggregation services (read-only, no mutations)

## Scope Excluded
- Photo charts
- Advanced analytics
- Exercise chart data infrastructure (workout_sets table) — requires WAVE-03
- Frontend chart rendering (Chart.js, Recharts, etc.)

## Dependencies And Other-Wave Fit
- WAVE-01 (Foundation): prerequisite — provides PIN auth middleware, Atlas GraphQL endpoint, gqlgen config, sqlc config
- WAVE-02 (Exercise Library): provides exercises table + ExerciseService.ListAll for exercise selector dropdown. WAVE-06 queries exercise list via existing contracts.
- WAVE-03 (Workout Diary): partial dependency — exercise chart queries (working weight, best set, e1RM, volume, reps, working sets count) require workout_sets table and daily_log_exercises bridge table that WAVE-03 creates. If WAVE-03 not deployed, exercise chart queries are stubbed as empty series with documentation. Body and nutrition charts are independent.
- WAVE-04 (Cardio and Body Tracking): provides body_weight_entries, body_check_ins, body_measurements tables with repos and services. WAVE-06 uses BodyWeightService.ListByDateRange for weight trends and adds new measurement range sqlc query for measurement trends.
- WAVE-05 (Nutrition): provides nutrition tables and NutritionMacroService. WAVE-06 wraps macro service for weekly averages.
- WAVE-07 (AI Export): WAVE-06 provides no data directly — charts share same underlying data. No collision.
- WAVE-08 (AI Review): no dependency.
- WAVE-09 (Backup): chart queries are ephemeral — no serialization concerns.

## Frontend Pages Dependencies
- PAGE-008 (Charts): primary consumer. Backend provides:
  - Exercise list for dropdown: WAVE-02 allExercises query (exists)
  - Period-filtered exercise progress: conditional on WAVE-03
  - Period-filtered body weight: WAVE-06 bodyWeightTrend query
  - Period-filtered measurement: WAVE-06 measurementTrend/measurementOverlay queries
  - Period-filtered nutrition: WAVE-06 nutritionWeeklyAverages query
- Dependency context only; no frontend pages, UI, or UX work in this wave.

## Codebase Fit And Touchpoints
- apps/api/internal/atlas/models/chart.go: new model types (ChartDataPoint, BodyWeightSeriesPoint, MeasurementTrendPoint, MeasurementOverlayGroup, NutritionWeeklyAverage, error types, result types)
- apps/api/internal/atlas/service/body_chart_service.go: new service for body weight trend, measurement trend, measurement overlay
- apps/api/internal/atlas/service/nutrition_weekly_avg_service.go: new service for weekly macro averaging (wraps NutritionMacroService)
- apps/api/internal/atlas/repository/postgres/queries/body_measurements_range.sql: new sqlc query for user+type+date range measurement queries
- apps/api/internal/atlas/repository/postgres/body_measurement_repo.go: add ListByUserTypeRange method
- apps/api/internal/atlas/graph/schema/charts.graphql: new schema file with chart types, queries, result unions
- apps/api/internal/atlas/graph/resolver/charts.go: new resolver file for chart queries with PIN auth guard
- apps/api/internal/atlas/graph/resolver/resolver.go: add ChartService field
- apps/api/internal/atlas/graph/schema/schema.graphql: add chart queries to root Query type

### Not Created In This Wave
- Exercise chart service: stubbed returning empty series. Full implementation requires WAVE-03 workout_sets table.
- No new database tables — chart queries use existing data.

## Design Contracts
- Read-only wave: WAVE-06 provides only GraphQL queries — no mutations, no new storage (DDEC-W06-001)
- e1RM formula: Epley formula — weight × (1 + reps / 30). In-memory computation, not stored in DB (DDEC-W06-002)
- Body measurement range query: JOIN body_measurements ↔ body_check_ins on check_in_id, filter by user_id and measurement_type and date range. Ordered by check_in.date (DDEC-W06-003)
- Nutrition weekly average: RULE-015 — sum of daily nutrition macro values / 7. Iterates days in each week within range (DDEC-W06-004)
- Empty data handling: All chart queries return empty series (zero-length arrays, no error) when no data matches the query criteria (DDEC-W06-005)
- Period filter: optional from/to date parameters. When not provided, default range is 4 weeks (DDEC-W06-006)
- Max date range cap: 52-week maximum enforced to prevent expensive nutrition iteration queries (DDEC-W06-011)
- Measurement overlay ordering: groups ordered alphabetically by measurementType name (DDEC-W06-007)

## Data API Integration And Operations

### GraphQL Schema (Proposed)

```graphql
type ChartDataPoint { date: Date!, value: Float! }
type BodyWeightSeriesPoint { date: Date!, weight: Float!, source: BodyWeightSource }
type MeasurementTrendPoint { date: Date!, value: Float!, side: MeasurementSide }
type MeasurementOverlayGroup { measurementType: MeasurementType!, dataPoints: [MeasurementTrendPoint!]! }
type NutritionWeeklyAverage { weekStartDate: Date!, calories: Float!, protein: Float!, fat: Float!, carbs: Float! }

type BodyWeightTrendResult { series: [BodyWeightSeriesPoint!]! }
type MeasurementTrendResult { dataPoints: [MeasurementTrendPoint!]! }
type MeasurementOverlayResult { groups: [MeasurementOverlayGroup!]! }
type NutritionWeeklyAveragesResult { averages: [NutritionWeeklyAverage!]! }

type ChartValidationError { message: String!, code: ChartErrorCode! }
type ChartAuthError { message: String!, code: ChartErrorCode! }
type ChartNotFoundError { message: String!, code: ChartErrorCode! }

enum ChartErrorCode { VALIDATION_ERROR NOT_FOUND AUTH_ERROR INTERNAL_ERROR }
```

### GraphQL Query Extensions
```
bodyWeightTrend(from: Date, to: Date): BodyWeightTrendResult!
measurementTrend(measurementType: MeasurementType!, from: Date, to: Date): MeasurementTrendResult!
measurementOverlay(measurementTypes: [MeasurementType!]!, from: Date, to: Date): MeasurementOverlayResult!
nutritionWeeklyAverages(from: Date, to: Date): NutritionWeeklyAveragesResult!
```

### REST Endpoints
- None. All chart operations via GraphQL.

### Log Markers
- [BodyChart][bodyWeightTrend] — body weight trend query
- [BodyChart][measurementTrend] — measurement trend query
- [BodyChart][measurementOverlay] — measurement overlay query
- [NutritionChart][weeklyAverages] — weekly averages query
- Sensitive data (body weight values, measurement values) NOT logged

### Operations
- PostgreSQL: no new migrations — chart queries use existing tables
- Existing Docker Compose stack, no new services
- Performance: measurement range query requires idx_body_checkin_user_date. Add if not present.
- Max date range: 52-week cap enforced to prevent expensive nutrition weekly iteration (DDEC-W06-011)

## Security Privacy And Compliance
- All chart queries protected by WAVE-01 PIN auth middleware (GraphQL)
- When PIN disabled, queries accessible without auth (consistent with TDEC-037)
- All operations scoped to default user per MVP constraint — userID extracted from context via middleware.GetAtlasUserID(ctx)
- Sensitive data NOT logged: body weight values, body measurement values, body fat %
- Nutrition macro values (calories, protein, fat, carbs) may be logged (non-sensitive per WAVE-05)
- Log markers record query type, date range, result count — NOT individual data values
- No PII in chart query responses — aggregated numeric data only
- No media, no uploads, no external API calls

## Implementation Slices

| Slice ID | Name | Description |
| --- | --- | --- |
| SLICE-W06-001 | Chart models | Define ChartDataPoint, BodyWeightSeriesPoint, MeasurementTrendPoint, MeasurementOverlayGroup, NutritionWeeklyAverage types + ChartResult unions (Success/ValidationError/AuthError) + ChartErrorCode enum in models/chart.go |
| SLICE-W06-002 | Body measurement range sqlc query | Add sqlc query body_measurements_range.sql that SELECTs from body_measurements JOIN body_check_ins by check_in_id, filtering by user_id, measurement_type, and check_in.date range. Ordered by check_in.date ASC. |
| SLICE-W06-003 | Body chart service | Implement BodyChartService with BodyWeightTrend, MeasurementTrend, MeasurementOverlay methods. Queries through existing BodyWeightEntryRepo and new body measurement range repo method. |
| SLICE-W06-004 | Nutrition weekly average service | Implement NutritionWeeklyAvgService with WeeklyAverages method. Iterates weeks in range, calls NutritionMacroService.Calculate for each day, averages per week per RULE-015. |
| SLICE-W06-005 | Chart GraphQL schema | Add charts.graphql with types, result unions, query definitions |
| SLICE-W06-006 | Chart GraphQL resolvers | Implement chart resolvers with PIN auth guard and union error returns following existing resolver patterns |
| SLICE-W06-007 | Main wiring | Add ChartService to Resolver struct, wire services in cmd/server/main.go, register route groups |
| SLICE-W06-008 | Epley e1RM helper | Stateless CalculateE1RM(weight, reps float64) float64 function in models/chart.go |

## Acceptance Criteria

| AC ID | Description | Status |
| --- | --- | --- |
| AC-W06-001 | Exercise progress query for a given exercise over a date range returns per-session time-series data | conditional (WAVE-03) |
| AC-W06-002 | e1RM calculated per-set using Epley formula: weight × (1 + reps / 30) | conditional (WAVE-03) |
| AC-W06-003 | Exercise chart query returns empty series when no data exists for the selected period/exercise | conditional (WAVE-03) |
| AC-W06-004 | Body weight trend query returns time-series data (date, weight) over a given date range, ordered by date ascending | implementable |
| AC-W06-005 | Body weight trend query returns empty series when no data exists for the selected period | implementable |
| AC-W06-006 | Body measurement trend query returns time-series data (date, measurementType, side, value) for a given measurement type over a date range | implementable |
| AC-W06-007 | Body measurement overlay query returns time-series data for multiple measurement types in a single response, grouped by measurementType, ordered alphabetically by type name | implementable |
| AC-W06-008 | Body measurement queries return empty series when no data exists for the selected period | implementable |
| AC-W06-009 | Body measurement overlay query with empty measurementTypes list returns empty groups array (no error) | implementable |
| AC-W06-010 | Nutrition weekly macro averages query returns time-series data (weekStartDate, calories, protein, fat, carbs) over a given date range | implementable |
| AC-W06-011 | Nutrition weekly averages use RULE-015: sum of daily values / 7 | implementable |
| AC-W06-012 | All chart queries accept optional from/to date parameters for period filtering | implementable |
| AC-W06-013 | Period filter with from > to date returns ValidationError | implementable |
| AC-W06-014 | All chart queries return AuthError when PIN session header is missing or invalid | implementable |
| AC-W06-015 | Nutrition weekly averages query returns empty series when no data exists for the selected period | implementable |

## Exit Criteria

| EC ID | Description |
| --- | --- |
| EC-W06-001 | AC-W06-001 through AC-W06-015 pass via TEST-W06-001 through TEST-W06-022 |
| EC-W06-002 | All chart queries protected by WAVE-01 PIN auth middleware |
| EC-W06-003 | gqlgen codegen produces valid Go code for WAVE-06 schema without drift |
| EC-W06-004 | Empty series returned for no-data periods (no errors) — verified for body, measurement, and nutrition queries |
| EC-W06-005 | Nutrition weekly average calculation matches RULE-015 |
| EC-W06-006 | Lint passes for all changed packages |
| EC-W06-007 | No sensitive data (body weight values, measurement values) in application logs |
| EC-W06-008 | Body weight trend query returns accurate date-ascending data |
| EC-W06-009 | Measurement overlay with empty types list returns empty groups array (no error) |
| EC-W06-010 | Matching exit criteria for conditional exercise chart ACs (AC-W06-001–003) — gated on WAVE-03 availability |

## Verification Obligations

| Test ID | Description | Type | Command |
| --- | --- | --- | --- |
| TEST-W06-001 | Exercise progress query returns correct time-series structure | integration | bunx nx run api:test -- --run '(?i)wave06_exercise_progress' |
| TEST-W06-002 | e1RM calculation matches Epley formula | unit | bunx nx run api:test -- --run '(?i)epley' |
| TEST-W06-003 | Exercise chart empty series for no data period | integration | bunx nx run api:test -- --run '(?i)wave06_exercise_empty' |
| TEST-W06-004 | Body weight trend query returns correct data | integration | bunx nx run api:test -- --run '(?i)body_weight_trend' |
| TEST-W06-005 | Body weight trend query returns empty series for no data | integration | bunx nx run api:test -- --run '(?i)body_weight_trend_empty' |
| TEST-W06-006 | Body measurement trend query returns correct data | integration | bunx nx run api:test -- --run '(?i)measurement_trend' |
| TEST-W06-007 | Body measurement overlay query returns multiple types | integration | bunx nx run api:test -- --run '(?i)measurement_overlay' |
| TEST-W06-008 | Nutrition weekly average query returns correct RULE-015 values | integration | bunx nx run api:test -- --run '(?i)nutrition_weekly_avg' |
| TEST-W06-009 | Nutrition weekly average query returns empty series for no data | integration | bunx nx run api:test -- --run '(?i)nutrition_weekly_avg_empty' |
| TEST-W06-010 | Period filter validation (from > to returns error) | unit | bunx nx run api:test -- --run '(?i)chart_date_validation' |
| TEST-W06-011 | Period filter default (12 weeks when not specified) | integration | bunx nx run api:test -- --run '(?i)chart_default_period' |
| TEST-W06-012 | Auth check — all chart queries return AuthError without PIN session | integration | bunx nx run api:test -- --run '(?i)wave06_auth' |
| TEST-W06-013 | Codegen drift check (gqlgen + sqlc) | codegen | bunx nx run api:codegen && bunx nx run graphql:codegen |
| TEST-W06-014 | Go lint for API package | lint | bunx nx run api:lint |
| TEST-W06-015 | GraphQL schema validate | codegen | bunx nx run graphql:validate |
| TEST-W06-016 | Log privacy: no body weight or measurement values in logs | unit | bunx nx run api:test -- --run '(?i)wave06_log_sanitize' |
| TEST-W06-017 | Max date range enforcement (52 weeks) | integration | bunx nx run api:test -- --run '(?i)chart_max_range' |
| TEST-W06-018 | Measurement trend with side filter (LEFT, RIGHT, NONE) | integration | bunx nx run api:test -- --run '(?i)measurement_side_trend' |
| TEST-W06-019 | Body weight trend with single data point | integration | bunx nx run api:test -- --run '(?i)body_weight_single_point' |
| TEST-W06-020 | Nutrition weekly average across partial week (mid-week start) | integration | bunx nx run api:test -- --run '(?i)nutrition_partial_week' |
| TEST-W06-021 | Measurement overlay with empty types list returns empty groups | integration | bunx nx run api:test -- --run '(?i)measurement_overlay_empty_types' |
| TEST-W06-022 | Measurement overlay returns groups ordered alphabetically by type | integration | bunx nx run api:test -- --run '(?i)measurement_overlay_ordering' |

## Rollout Rollback And Compatibility
- Rollout: merge PR, CI builds and runs tests, deploy via Dokploy compose update. Additive queries only — no migration needed.
- Rollback: revert PR, CI builds previous image, Dokploy compose update rolls back. No data migration needed.
- Compatibility: all new operations are additive. No existing API changes. WAVE-01, WAVE-02, WAVE-04, WAVE-05 endpoints unchanged.
- Migration: no new migrations in WAVE-06.
- Max range cap: 52 weeks enforced via server constant.

## Handoff Packets
- HANDOFF-W06-001: This wave brief document
- HANDOFF-W06-002: Planner reports (6 scopes, some with 2 attempts)
- HANDOFF-W06-003: Reviewer evidence (7 perspectives, final fit reviewer)
- HANDOFF-W06-004: Consolidated question ledger (9 entries)

## Design Decisions

| DDEC ID | Decision | Rationale |
| --- | --- | --- |
| DDEC-W06-001 | Read-only wave — no mutations, no storage changes | Chart data is aggregated from existing data. No new write endpoints needed. |
| DDEC-W06-002 | Epley formula for e1RM | User-selected on 2026-06-21. Q-CHART-001 resolved. |
| DDEC-W06-003 | Measurement range query via check_in JOIN | Body measurements are child of check_in, not date-indexed. JOIN allows date-range filtering. |
| DDEC-W06-004 | Nutrition weekly average via iteration | NutritionMacroService.Calculate is per-day. Iterating days per week and averaging is the correct RULE-015 implementation. |
| DDEC-W06-005 | Empty series for no-data periods | EDGE-008: chart with no data must not error. Empty series is the correct pattern. |
| DDEC-W06-006 | Default period configurable | Default chart period TBD per DQ-W06-004. Config const allows deployment decision. |
| DDEC-W06-007 | Measurement overlay alphabetically ordered | Predictable frontend rendering. Alphabetical by measurement type name. |
| DDEC-W06-008 | Best set = highest e1RM per session | User-selected on 2026-06-21. DQ-W06-001 resolved. |
| DDEC-W06-009 | Working weight per session from WorkoutExercise.workingWeightSnapshot | User-selected on 2026-06-21. Consistent with RULE-017. DQ-W06-003 resolved. |
| DDEC-W06-010 | Exercise chart stubs returning empty series until WAVE-03 | User-selected on 2026-06-21. DQ-W06-006 resolved. |
| DDEC-W06-011 | 52-week max date range cap | User-selected on 2026-06-21. Prevents expensive nutrition iteration. DQ-W06-005 resolved. |

## Reviewer Verdicts

| Wave | Perspective | Attempt | Verdict | Reviewer Report | Required Revisions | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| WAVE-06 | product-scope-and-ac | 1 | needs-revision | review-product-scope-and-ac-attempt-1.md | Split AC-W06-001, add empty-series AC, move default to DQ | Addressed in attempt 2 |
| WAVE-06 | product-scope-and-ac | 2 | approved | review-product-scope-and-ac-attempt-2.md | none | All concerns addressed |
| WAVE-06 | architecture-codebase-fit | 1 | approved | review-architecture-codebase-fit-attempt-1.md | none | 8 slices, pattern consistent with WAVE-04/05 |
| WAVE-06 | data-api-integration-ops | 1 | approved | review-data-api-integration-ops-attempt-1.md | none | Clean schema, additive queries |
| WAVE-06 | security-privacy-compliance | 1 | approved | review-security-privacy-compliance-attempt-1.md | none | PIN auth, log privacy covered |
| WAVE-06 | testing-exit-criteria | 1 | needs-revision | review-testing-exit-criteria-attempt-1.md | Add empty-types test, document conditional tests | Addressed in attempt 2 |
| WAVE-06 | testing-exit-criteria | 2 | approved | review-testing-exit-criteria-attempt-2.md | none | 22 tests, all AC/EC covered |
| WAVE-06 | sequencing-other-wave-fit | 1 | approved | review-sequencing-other-wave-fit-attempt-1.md | none | WAVE-03 dependency correctly identified |
| WAVE-06 | traceability-consistency | 1 | needs-revision | review-traceability-consistency-attempt-1.md | Add DQ-W06-004/005, consolidate AC refs | Addressed in attempt 2 |
| WAVE-06 | traceability-consistency | 2 | approved | review-traceability-consistency-attempt-2.md | none | All concerns addressed |
| WAVE-06 | final-wave-fit-review | 1 | approved | final-wave-fit-review-attempt-1.md | none | All 9 checks pass. Ready for user approval. |

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W06-001 | WAVE-06 | product-ac | resolved | RULE-014 | Best set definition: heaviest weight (load only) or highest e1RM (load × reps)? | Chart AC-066 defines "best set per session" but does not specify the metric. | User chose: highest e1RM. | user decision 2026-06-21 | resolved | Highest e1RM per session. DDEC-W06-008. |
| DQ-W06-002 | WAVE-06 | architecture | resolved | WAVE-03 | No workout_set table exists yet. How to handle exercise chart queries? | Exercise charts require set-level data that does not exist in the DB. | User chose: stubs returning empty series. See DQ-W06-006. | user decision 2026-06-21 | resolved | Exercise chart queries return empty series until WAVE-03. |
| DQ-W06-003 | WAVE-06 | product-ac | resolved | AC-065 | Working weight per session: static exercise.workingWeight or per-session snapshot? | Affects chart accuracy. RULE-017 defines snapshot behavior. | User chose: per-session snapshot. | user decision 2026-06-21 | resolved | WorkoutExercise.workingWeightSnapshot per RULE-017. DDEC-W06-009. |
| DQ-W06-004 | WAVE-06 | data-ops | resolved | — | Default chart date range when no from/to specified? | Initial chart display on frontend PAGE-008. | User chose: 4 weeks. | user decision 2026-06-21 | resolved | Default chart period: 4 weeks. DDEC-W06-006. |
| DQ-W06-005 | WAVE-06 | data-ops | deferred | — | Max date range to prevent expensive queries? | Nutrition weekly averages iterate per week. Large ranges could be slow. | User chose: 52-week max. | user decision 2026-06-21 | open | 52-week max enforced via server constant. |
| DQ-W06-006 | WAVE-06 | data-ops | resolved | WAVE-03 | Exercise chart queries: stubs returning empty, or omitted? | Frontend PAGE-008 needs to know if exercise chart queries exist. | User chose: stubs returning empty. | user decision 2026-06-21 | resolved | Stubs returning empty series. DDEC-W06-010. |

## Traceability
- docs/prd-waves/waves/wave-06.md: source wave boundary, outcomes, capability groups
- docs/product-verified/features/charts.md: chart feature spec
- docs/product-verified/acceptance-criteria.md: AC-020–AC-022, AC-065–AC-073
- docs/product-verified/edge-cases.md: EDGE-008, EDGE-026
- docs/product-verified/business-rules.md: RULE-012 (Epley), RULE-013 (volume), RULE-014 (best set), RULE-015 (weekly avg)
- docs/prd-waves/frontend-pages/page-008.md: PAGE-008 backend dependencies
- docs/prd-wave-details/waves/wave-04.md: WAVE-04 patterns, body tracking data contracts
- docs/prd-wave-details/waves/wave-05.md: WAVE-05 patterns, nutrition macro service contract
- docs/development-plan.xml: M-API, M-PRD-WAVE-DETAILER module contracts
- docs/knowledge-graph.xml: existing module boundaries
- apps/api/internal/atlas: existing codebase patterns for service/repository/resolver structure
