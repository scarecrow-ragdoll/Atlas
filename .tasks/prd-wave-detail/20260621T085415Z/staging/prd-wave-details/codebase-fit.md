# Codebase Fit

## Relevant Modules
- apps/api: Go HTTP API — target for all WAVE-06 code
- apps/api/internal/atlas: Atlas fitness module — all new files added here
- apps/api/internal/atlas/graph/schema: Chart GraphQL schema added here
- apps/api/internal/atlas/service: Chart services (body chart, nutrition weekly avg)
- apps/api/internal/atlas/repository/postgres: Measurement range sqlc query + repo method

No new Nx packages or top-level modules.

## Relevant Files Read
- apps/api/internal/atlas/service/exercise.go — service pattern
- apps/api/internal/atlas/service/nutrition_macro_service.go — calculation service pattern
- apps/api/internal/atlas/service/body_weight_service.go — body weight service (existing WAVE-04)
- apps/api/internal/atlas/service/body_checkin_service.go — body measurement access patterns
- apps/api/internal/atlas/repository/postgres/body_weight_entry_repo.go — repo pattern
- apps/api/internal/atlas/repository/postgres/body_measurement_repo.go — measurement repo (needs ListByUserTypeRange addition)
- apps/api/internal/atlas/repository/postgres/queries/body_weight_entries.sql — sqlc query pattern
- apps/api/internal/atlas/repository/postgres/queries/body_measurements.sql — existing measurement queries
- apps/api/internal/atlas/graph/resolver/resolver.go — resolver container
- apps/api/internal/atlas/graph/schema/body_tracking.graphql — existing body schema for type reference
- apps/api/cmd/server/main.go — wiring pattern
- apps/api/atlas-gqlgen.yml — gqlgen config with model bindings
- apps/api/sqlc.yaml — sqlc config (auto-discovers new queries via glob)

## Public Contracts
- bodyWeightEntries(dateFrom, dateTo): existing WAVE-04 query — used by WAVE-06 for chart data
- measurements by checkInId: existing WAVE-04 — WAVE-06 adds new user+type+date range query
- nutritionMacros(weekStartDate, date): existing WAVE-05 — used by WAVE-06 for weekly averages
- All operations require PIN auth session when PIN is enabled
- Error format per existing pattern: union result types (Success | ValidationError | AuthError)

## Generated Artifact Impact
- gqlgen (atlas-gqlgen.yml): auto-discovers charts.graphql via glob — generates new chart types, query stubs, result unions
- sqlc: auto-discovers body_measurements_range.sql via glob — generates new measurement range query function
- No existing generated artifacts affected — all additions are additive
- atlas-gqlgen.yml: needs new model binding entries for chart types (ChartDataPoint, etc.)

## Integration Points
- PIN auth middleware from WAVE-01: guards all WAVE-06 GraphQL queries via /graphql/atlas
- WAVE-02 ExerciseService: used for exercise list dropdown (allExercises query)
- WAVE-04 BodyWeightService: used for weight trend queries
- WAVE-04 BodyMeasurementRepo: extended with ListByUserTypeRange method
- WAVE-05 NutritionMacroService: wrapped for weekly average calculation
- WAVE-06 feeds data to PAGE-008 frontend via GraphQL

## Likely Graph Deltas
- M-API (atlas module) gains: BodyChartService, NutritionWeeklyAvgService
- apps/api/internal/atlas/graph/schema gains: charts.graphql
- apps/api/internal/atlas/graph/resolver/resolver.go gains: new service fields
- apps/api/cmd/server/main.go gains: 2 service wiring additions
- apps/api/internal/atlas/repository/postgres/queries gains: body_measurements_range.sql
- apps/api/internal/atlas/repository/postgres gains: body_measurement_repo.go method addition
- apps/api/internal/atlas/models gains: chart.go with ChartDataPoint, body chart types, nutrition weekly types, error types, result types

## Unsupported Assumptions
- WAVE-03 workout_sets table does not exist — exercise chart queries are stubs returning empty series
- WAVE-03 WorkoutExercise table may use different schema than assumed — stubs are independent of schema
- WAVE-04 body_measurement schema (measurementType enum) assumed to exist — verify enum values before implementation (checked in codebase audit)
- Measurement range query assumes body_check_ins table has user_id index and date column — verified present in WAVE-04 migration
- Nutrition weekly average assumes NutritionMacroService.Calculate works per-day with date parameter — verified present in WAVE-05
- sqlc auto-discovery works via glob pattern — verified in WAVE-04 and WAVE-05
- No migration number collision — WAVE-06 adds no migrations