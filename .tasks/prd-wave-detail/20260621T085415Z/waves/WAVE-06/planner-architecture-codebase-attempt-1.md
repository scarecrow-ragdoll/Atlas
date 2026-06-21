# WAVE-06 Architecture-Codebase Planner Attempt 1

## Sources Read
- apps/api/internal/atlas — directory structure
- apps/api/internal/atlas/models/ — existing types
- apps/api/internal/atlas/service/ — existing services
- apps/api/internal/atlas/repository/postgres/ — existing repos
- apps/api/internal/atlas/graph/resolver/resolver.go — resolver container
- apps/api/internal/atlas/graph/schema/schema.graphql — root Query type
- apps/api/internal/atlas/graph/schema/exercises.graphql — exercise schema
- apps/api/internal/atlas/graph/schema/body_tracking.graphql — body tracking schema
- apps/api/internal/atlas/graph/schema/nutrition.graphql — nutrition schema
- apps/api/internal/atlas/service/nutrition_macro_service.go — macro calculation
- apps/api/internal/atlas/service/exercise.go — exercise service
- apps/api/internal/atlas/models/daily_log.go — DailyLogRecord
- apps/api/internal/repository/postgres/migrations/00083_daily_logs.sql — daily_logs table
- docs/prd-waves/waves/wave-06.md

## Selected Backend Wave Boundary
WAVE-06 adds GraphQL query resolvers, services, and potentially sqlc queries for chart data aggregation. No new tables, no mutations.

## Neighboring Backend Wave Fit
- WAVE-02: Provides exercises table + ExerciseService.ListAll() for exercise selector dropdown
- WAVE-03: Should provide daily_log_exercises + workout_sets — NOT implemented. WAVE-06 exercise chart queries are blocked on WAVE-03.
- WAVE-04: Provides body_weight_entries, body_check_ins, body_measurements tables with repos and services — fully queryable
- WAVE-05: Provides nutrition tables with NutritionMacroService.Calculate() — reusable for weekly averages

## Frontend Pages Context
PAGE-008 needs: exercise list for dropdown (already exists), period-filtered queries for 7 chart types.

## Codebase Evidence

### Existing Relevant Contracts
- `ExerciseService.ListAll(ctx, userID, includeInactive)` — returns all exercises for selector dropdown
- `BodyWeightService.ListByDateRange(ctx, userID, fromDate, toDate)` — returns []BodyWeightEntry
- `BodyMeasurement` — repo has ListByCheckIn, no range-based measurement listing
- `NutritionMacroService.Calculate(ctx, userID, weekStartDate, date)` — returns per-day macros (single day, not week range)

### Gaps
1. No range-based measurement query: BodyMeasurementRepo only lists measurements by checkInId. Need new sqlc query for measurements by type + date range across check-ins.
2. No range-based macro service: NutritionMacroService.Calculate is per-week+per-day. Need weekly-average service that iterates weeks in a range.
3. No workout_sets table: Exercise progress queries cannot be implemented until WAVE-03.
4. No e1RM calculation service: Need a stateless helper.

## Proposed Details

### New Files
- `apps/api/internal/atlas/models/chart.go` — new model types for chart data responses
- `apps/api/internal/atlas/service/chart_service.go` — chart data aggregation service (body, nutrition parts)
- `apps/api/internal/atlas/graph/schema/charts.graphql` — GraphQL types and queries
- `apps/api/internal/atlas/graph/resolver/charts.go` — resolver implementation

### Modified Files
- `apps/api/internal/atlas/graph/schema/schema.graphql` — add chart queries to root Query
- `apps/api/internal/atlas/graph/resolver/resolver.go` — add ChartService field

### Implementation Slices

| Slice ID | Name | Description |
|---|---|---|
| SLICE-W06-001 | Chart models | Define ChartDataPoint (date, value), BodyWeightSeries (date, weight), MeasurementTrendPoint (date, type, side, value), MeasurementOverlayResult (multiple types), NutritionWeeklyAverage (weekStartDate, calories, protein, fat, carbs), ChartResult union types with auth/validation errors |
| SLICE-W06-002 | Body chart service | BodyChartService with methods: BodyWeightTrend(from, to Date) → []BodyWeightSeries, MeasurementTrend(measurementType, from, to) → []MeasurementTrendPoint, MeasurementOverlay(types, from, to) → MeasurementOverlayResult. Queries through existing BodyWeightEntryRepo and new measurement range repo method. |
| SLICE-W06-003 | Nutrition weekly average service | NutritionWeeklyAvgService with method: WeeklyAverages(from, to Date) → []NutritionWeeklyAverage. Iterates weeks in range, calls NutritionMacroService.Calculate for each day, averages per week per RULE-015. |
| SLICE-W06-004 | Measurement range sqlc query | Add sqlc query: ListMeasurementsByUserTypeRange(userId, measurementType, from, to) → rows from body_measurements JOIN body_check_ins on check_in_id. |
| SLICE-W06-005 | Exercise chart stub service | ExerciseChartService with method returning empty series + documentation that WAVE-03 infrastructure is required. Define contract for when WAVE-03 adds sets. |
| SLICE-W06-006 | Chart GraphQL schema | Add charts.graphql with types (ChartDataPoint, BodyWeightSeries, MeasurementTrendPoint, MeasurementOverlayResult, NutritionWeeklyAverage), result unions, query definitions |
| SLICE-W06-007 | Chart GraphQL resolvers | Implement chart resolvers with PIN auth guard and union error returns following existing patterns |
| SLICE-W06-008 | Epley e1RM helper | Stateless function: CalculateE1RM(weight, reps float64) float64 — returns weight × (1 + reps / 30). Placed in models/chart.go or a new utils package. |

## Acceptance Criteria Contributions
See planner-product-ac-attempt-1.md for full AC list.

## Exit Criteria Contributions
See planner-product-ac-attempt-1.md for full EC list.

## Verification Contributions
See planner-product-ac-attempt-1.md for full TEST list.

## Risks And Rollback
- WAVE-03 dependency: exercise chart queries require workout_sets table. If WAVE-03 is still in planning, exercise chart feature is blocked. Mitigation: non-exercise charts (body, nutrition) are fully implementable.
- No changes to existing tables or migrations — purely additive (new query schemas, new service files). Rollback by removing new files and reverting root Query additions.

## Questions Raised
- DQ-W06-002: Does WAVE-06 need to define its own set query schema against raw daily_logs or wait for WAVE-03?
- DQ-W06-004: What is the default period for chart queries when no from/to is provided? (Proposed: last 12 weeks)

## Traceability Candidates
- apps/api/internal/atlas: service/, models/, graph/schema/, graph/resolver/
- docs/prd-waves/waves/wave-06.md
- docs/prd-wave-details/waves/wave-04.md — WAVE-04 patterns
- docs/prd-wave-details/waves/wave-05.md — WAVE-05 patterns