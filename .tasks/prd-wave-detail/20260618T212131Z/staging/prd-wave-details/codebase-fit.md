# Codebase Fit

## Relevant Modules
- apps/api: Go HTTP API — target for all WAVE-03 code
- libs/graphql/schema: GraphQL schema files — new workout.graphql added here
- apps/api/internal/repository/postgres: target for sqlc migrations, queries, and repository adapters
- apps/api/internal/service: target for workout service layer
- apps/api/internal/graph: target for workout GraphQL resolvers

## Relevant Files Read
- apps/api/cmd/server/main.go — wiring pattern for repos, services, handlers, resolvers, route groups
- apps/api/internal/appconfig/config.go — config struct and env overlay pattern
- apps/api/internal/middleware/admin_auth.go — auth middleware pattern for PIN guard replication
- apps/api/internal/service/admin_auth.go — service layer pattern (transport-neutral, validation, log markers)
- apps/api/internal/repository/postgres/user_repo.go — repository adapter with sqlc queries (narrowed interface pattern)
- apps/api/internal/repository/redis/admin_session_store.go — HMAC key derivation pattern
- apps/api/gqlgen.yml — schema glob pattern for auto-discovery
- apps/api/sqlc.yaml — query glob pattern for auto-discovery
- libs/graphql/schema/schema.graphql — extend type Query/Mutation pattern
- libs/graphql/schema/admin_auth.graphql — union result type pattern
- libs/graphql/schema/common.graphql — common types (UUID, DateTime, ValidationError, AuthError, NotFoundError)
- apps/api/internal/repository/postgres/migrations/00079_admin_users.sql — migration structure pattern (goose format)

## Public Contracts
- DailyLog GraphQL operations: dailyLogByDate (query), dailyLogsByDateRange (query), upsertDailyLog (mutation), deleteDailyLog (mutation)
- WorkoutExercise GraphQL operations: addWorkoutExercise (mutation), updateWorkoutExercise (mutation), removeWorkoutExercise (mutation)
- WorkoutSet GraphQL operations: addWorkoutSet (mutation), updateWorkoutSet (mutation), removeWorkoutSet (mutation)
- CardioEntry GraphQL operations: addCardioEntry (mutation), updateCardioEntry (mutation), removeCardioEntry (mutation)
- All operations require PIN auth session (header-based) when PIN is enabled
- Error format per TDEC-027: { "error": { "code": "ERROR_CODE", "message": "..." } }

## Generated Artifact Impact
- gqlgen: auto-discovers new workout.graphql via glob — generates DailyLog, WorkoutExercise, WorkoutSet, CardioEntry model structs, input types, union result types, resolver stubs
- sqlc: auto-discovers new query .sql files via glob — generates CRUD query functions for daily_logs, workout_exercises, workout_sets, and cardio_entries tables
- No existing generated artifacts affected — all additions are additive

## Integration Points
- WAVE-01 PIN auth middleware: guards all WAVE-03 GraphQL queries and mutations
- WAVE-02 allExercises query: called from workout service to read Exercise.workingWeight for snapshot
- WAVE-01 common types: ValidationError, AuthError, NotFoundError reused in union results
- WAVE-01 exercises table: FK target for workout_exercises.exercise_id
- GraphQL schema: extends root Query and Mutation types following existing pattern

## Likely Graph Deltas
- M-API gains: DailyLogRepo, WorkoutExerciseRepo, WorkoutSetRepo, CardioEntryRepo, WorkoutService dependencies, PIN-protected sub-route group for fitness GraphQL
- libs/graphql/schema gains: workout.graphql type extensions
- apps/api gains: daily_log_repo.go, workout_exercise_repo.go, workout_set_repo.go, cardio_entry_repo.go, workout.go (service), workout.resolvers.go
- No new Nx packages or top-level modules

## Unsupported Assumptions
- WAVE-01 will provide requirePinAuth(ctx) function and PIN middleware for chi — currently non-existent
- WAVE-01 will provide ValidationError, AuthError, NotFoundError common GraphQL types — currently non-existent
- WAVE-02 will provide allExercises query returning workingWeight field — currently non-existent
- WAVE-01 migration numbers end at 00079 — may shift if WAVE-01 adds more migrations
- WAVE-02 migration numbers are 00080/00081 — may shift
