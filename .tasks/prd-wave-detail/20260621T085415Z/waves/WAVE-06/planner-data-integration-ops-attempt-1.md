# WAVE-06 Data-Integration-Ops Planner Attempt 1

## Sources Read
- docs/prd-waves/waves/wave-06.md
- docs/product-verified/features/charts.md
- docs/product-verified/edge-cases.md (EDGE-008, EDGE-026)
- docs/technical-verified/auth-security-compliance.md
- docs/technical-verified/data-contracts.md
- docs/prd-wave-details/waves/wave-04.md
- docs/prd-wave-details/waves/wave-05.md
- apps/api/internal/atlas/service/body_weight.go
- apps/api/internal/atlas/service/nutrition_macro_service.go
- apps/api/internal/atlas/repository/postgres/ — repo patterns

## Selected Backend Wave Boundary
Read-only data aggregation queries. No new tables, no mutations, no external integrations, no event/async processing.

## Neighboring Backend Wave Fit
- WAVE-03: No set data yet — exercise queries not implementable
- WAVE-04: BodyWeightEntry (range query exists), BodyMeasurement (needs range query)
- WAVE-05: Nutrition macro service exists but needs weekly-average wrapper

## Frontend Pages Context
PAGE-008 expects period-filtered data for each chart type. The "period" is a date range ([from, to]) passed as GraphQL arguments.

## Codebase Evidence
- Existing body_weight_entries table has (user_id, date, weight) — range queries via sqlc
- Existing body_measurements table has (check_in_id, measurement_type, side, value) — queries via checkInId only, no user-level range
- Existing body_check_ins table has (user_id, date) — join needed for measurement range queries
- NutritionMacroService already handles soft-deleted products (skips them) and empty templates (returns 0s)
- All current services use zap.Logger for log markers

## Proposed Details

### Data Lifecycle
- No new data storage in WAVE-06
- Chart queries transform existing data into time-series format
- e1RM is computed in-memory (Epley formula) — no derived column in DB

### GraphQL Schema Design

```graphql
# chart types
type ChartDataPoint { date: Date!, value: Float! }
type BodyWeightSeriesPoint { date: Date!, weight: Float!, source: BodyWeightSource }
type MeasurementTrendPoint { date: Date!, value: Float!, side: MeasurementSide }
type MeasurementOverlayGroup { measurementType: MeasurementType!, dataPoints: [MeasurementTrendPoint!]! }
type NutritionWeeklyAverage { weekStartDate: Date!, calories: Float!, protein: Float!, fat: Float!, carbs: Float! }

# results
type BodyWeightTrendResult { series: [BodyWeightSeriesPoint!]! }
type MeasurementTrendResult { dataPoints: [MeasurementTrendPoint!]! }
type MeasurementOverlayResult { groups: [MeasurementOverlayGroup!]! }
type NutritionWeeklyAveragesResult { averages: [NutritionWeeklyAverage!]! }

# query additions to Query
bodyWeightTrend(from: Date, to: Date): BodyWeightTrendResult!
measurementTrend(measurementType: MeasurementType!, from: Date, to: Date): MeasurementTrendResult!
measurementOverlay(measurementTypes: [MeasurementType!]!, from: Date, to: Date): MeasurementOverlayResult!
nutritionWeeklyAverages(from: Date, to: Date): NutritionWeeklyAveragesResult!
```

### Log Markers
- [BodyChart][bodyWeightTrend][list] — body weight trend query
- [BodyChart][measurementTrend][list] — measurement trend query
- [BodyChart][measurementOverlay][list] — measurement overlay query
- [NutritionChart][weeklyAverages][list] — weekly averages query
- No sensitive data logged: weight values may be logged (non-sensitive per WAVE-04), measurement values may be logged

### Operations
- Pure GraphQL queries — no REST, no batch endpoints
- No new services or infrastructure
- Performance: measurement queries JOIN body_check_ins + body_measurements — index on (user_id, date) for check-ins, (check_in_id) for measurements. Add covering index for range queries if needed.
- Nutrition weekly averages: iteration over weeks could be expensive for large ranges. Cap range to 52 weeks max (1 year) to prevent abuse.
- No database migration — no new tables

### Rollout/Rollback/Compatibility
- Rollout: merge PR, CI builds, deploy — additive queries only
- Rollback: revert PR, CI reverts — no data migration needed
- Compatibility: fully additive, no existing API changes

## Risks And Rollback
- Nutrition weekly averages over large ranges: iteration cost. Mitigation: hard cap at 52 weeks.
- Measurement range query performance without index. Mitigation: add idx_body_checkin_user_date index if not present, or include in implementation notes.
- EDGE-026 (system clock): queries use provided date range, not current time — safe.

## Questions Raised
- DQ-W06-005: Should there be a max date range to prevent expensive queries? (Proposed: 52 weeks)
- DQ-W06-006: Should exercise chart stubs return empty series or be omitted from the schema until WAVE-03 is ready?

## Traceability Candidates
- docs/technical-verified/data-contracts.md — domain entities
- docs/technical-verified/api-contracts.md — GraphQL pattern
- docs/technical-verified/operations-observability.md — log markers, error format
- apps/api/internal/atlas/repository/postgres/queries/ — existing sqlc queries