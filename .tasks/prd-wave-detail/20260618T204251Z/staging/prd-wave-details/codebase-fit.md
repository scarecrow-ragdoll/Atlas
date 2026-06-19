# Codebase Fit

## Relevant Modules
- apps/api: Go HTTP API — target for all WAVE-02 code
- libs/graphql/schema: GraphQL schema files — new exercises.graphql added here
- libs/go/config: Shared config package — no changes needed (WAVE-01 provides MediaConfig)

## Relevant Files Read
- apps/api/cmd/server/main.go — wiring pattern for repos, services, handlers, resolvers, route groups
- apps/api/internal/appconfig/config.go — config struct and env overlay pattern
- apps/api/internal/middleware/admin_auth.go — auth middleware pattern for PIN guard replication
- apps/api/internal/service/admin_auth.go — service layer pattern (transport-neutral, bcrypt, validation)
- apps/api/internal/repository/postgres/user_repo.go — repository adapter with sqlc queries
- apps/api/internal/repository/redis/admin_session_store.go — HMAC key derivation pattern
- apps/api/gqlgen.yml — schema glob pattern for auto-discovery
- apps/api/sqlc.yaml — query glob pattern for auto-discovery
- libs/graphql/schema/schema.graphql — extend type Query/Mutation pattern
- libs/graphql/schema/admin_auth.graphql — union result type pattern
- apps/api/internal/repository/postgres/migrations/00079_admin_users.sql — migration structure pattern

## Public Contracts
- Exercise GraphQL operations: exercises (paginated), exercise(id), allExercises (unpaginated), createExercise, updateExercise, deleteExercise
- ExerciseMedia REST endpoints: POST /api/v1/exercise-media (upload), GET /api/v1/exercise-media/{id} (download), DELETE /api/v1/exercise-media/{id} (delete)
- All operations require PIN auth session (header-based) when PIN is enabled
- Error format per TDEC-027: { "error": { "code": "ERROR_CODE", "message": "..." } }

## Generated Artifact Impact
- gqlgen: auto-discovers exercises.graphql via glob — generates Exercise, ExerciseMedia models, union result types, resolver stubs
- sqlc: auto-discovers exercises.sql via glob — generates CRUD query functions for exercise and exercise_media tables
- No existing generated artifacts affected — all additions are additive

## Integration Points
- PIN auth middleware from WAVE-01: guards all GraphQL and REST exercise endpoints
- WAVE-01 media storage: exercise-media endpoints use same file storage as WAVE-01 media scaffold
- WAVE-03: allExercises query is the stable interface for workout diary exercise selector
- GraphQL schema: extends root Query and Mutation types following existing pattern

## Likely Graph Deltas
- M-API gains: ExerciseService dependency, ExerciseMediaHandler routes, PIN-protected sub-route group
- libs/graphql/schema gains: exercise.graphql type extensions
- apps/api gains: exercise_repo.go, exercise.go (service), exercise_media.go (handler), exercise.resolvers.go
- No new Nx packages or top-level modules

## Unsupported Assumptions
- WAVE-01 will provide requirePinAuth(ctx) function and PIN middleware for chi — currently non-existent
- WAVE-01 will provide ValidationError, AuthError, NotFoundError common GraphQL types — currently non-existent
- WAVE-01 will provide MediaConfig with BasePath string — currently non-existent
- WAVE-01 media scaffold provides POST/GET /api/v1/media endpoints — currently non-existent
- WAVE-01 migration numbers end at 00079 — may shift if WAVE-01 adds more migrations