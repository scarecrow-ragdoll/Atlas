# Codebase Fit
## Relevant Modules
- apps/api: Go HTTP API — target for all WAVE-04 code
- libs/graphql/schema: GraphQL schema files — new cardio.graphql, body_weight.graphql, body_checkin.graphql, week_flag.graphql added here
- libs/go/config: Shared config package — no changes needed (WAVE-01 provides MediaConfig)
## Relevant Files Read
- apps/api/internal/repository/postgres/user_repo.go — repository adapter pattern with sqlc
- apps/api/internal/service/exercise.go — service layer pattern (transport-neutral, validation)
- apps/api/internal/handler/exercise_media.go — REST handler pattern for multipart upload/download/delete
- apps/api/internal/graph/exercise.resolvers.go — GraphQL resolver pattern with union result types
- apps/api/cmd/server/main.go — wiring pattern for repos, services, handlers, resolvers, route groups
- apps/api/gqlgen.yml — schema glob pattern for auto-discovery
- apps/api/sqlc.yaml — query glob pattern for auto-discovery
- libs/graphql/schema/exercises.graphql — GraphQL schema pattern
- apps/api/internal/middleware/admin_auth.go — auth middleware pattern
- apps/api/internal/appconfig/config.go — config struct extension pattern
## Public Contracts
- CardioEntry GraphQL operations: createCardioEntry, updateCardioEntry, deleteCardioEntry, cardioEntries (by date), cardioEntry (by ID)
- BodyWeightEntry GraphQL operations: createBodyWeightEntry, updateBodyWeightEntry, deleteBodyWeightEntry, bodyWeightEntries (date range), latestBodyWeight
- BodyCheckIn GraphQL operations: createBodyCheckIn, updateBodyCheckIn, deleteBodyCheckIn, bodyCheckIns (date range), bodyCheckIn (by ID with nested measurements + photos)
- BodyMeasurement GraphQL operations: createBodyMeasurement, updateBodyMeasurement, deleteBodyMeasurement
- ProgressPhoto: progressPhotos GraphQL query (by checkInId), REST POST/GET/DELETE /api/v1/progress-photos
- WeekFlag GraphQL operations: createWeekFlag, deleteWeekFlag, weekFlags (by weekStartDate)
- All operations require PIN auth session (header-based) when PIN is enabled
- Error format per TDEC-027: { "error": { "code": "ERROR_CODE", "message": "..." } }
## Generated Artifact Impact
- gqlgen: auto-discovers cardio.graphql, body_weight.graphql, body_checkin.graphql, week_flag.graphql via glob — generates new types, union results, resolver stubs
- sqlc: auto-discovers cardio_entries.sql, body_weight_entries.sql, body_check_ins.sql, body_measurements.sql, progress_photos.sql, week_flags.sql via glob — generates CRUD query functions for all 6 tables
- No existing generated artifacts affected — all additions are additive
## Integration Points
- PIN auth middleware from WAVE-01: guards all GraphQL and REST WAVE-04 endpoints
- WAVE-01 media storage: progress-photo endpoints use same file storage pattern as WAVE-01 media scaffold
- WAVE-03/04 DailyLog table: CardioEntry FK references daily_log table (either from WAVE-03 or WAVE-04's own migration)
- WAVE-06 (Charts): latestBodyWeight and bodyWeightEntries queries are the stable interface for chart data
- GraphQL schema: extends root Query and Mutation types following existing pattern
## Likely Graph Deltas
- M-API gains: CardioService, BodyWeightService, BodyCheckInService, WeekFlagService dependencies, ProgressPhotoHandler routes, PIN-protected sub-route group
- libs/graphql/schema gains: cardio.graphql, body_weight.graphql, body_checkin.graphql, week_flag.graphql type extensions
- apps/api gains: 6 migration files, 6 sqlc query files, 6 repo adapters, 4 service files, 4 resolver files, 1 handler file
- No new Nx packages or top-level modules
## Unsupported Assumptions
- WAVE-01 will provide requirePinAuth(ctx) function and PIN middleware for chi — currently non-existent
- WAVE-01 will provide ValidationError, AuthError, NotFoundError common GraphQL types — currently non-existent
- WAVE-01 will provide MediaConfig with BasePath string — currently non-existent
- WAVE-01 media scaffold provides POST/GET /api/v1/media endpoints — currently non-existent
- WAVE-01 migration numbers end at 00079 — may shift if WAVE-01 adds more migrations
- WAVE-03 daily_log table may or may not exist at WAVE-04 deployment time (DQ-W04-001)