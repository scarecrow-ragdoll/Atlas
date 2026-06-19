# WAVE-02 architecture-codebase Planner Attempt 1

## Sources Read
- apps/api/cmd/server/main.go (API wiring)
- apps/api/internal/appconfig/config.go (config extension pattern)
- apps/api/internal/graph/resolver.go (resolver dependency container)
- apps/api/internal/graph/schema.resolvers.go (gqlgen resolver pattern)
- apps/api/internal/graph/admin_auth.resolvers.go (auth resolver pattern)
- apps/api/internal/graph/admin_auth_helpers.go (requireAdmin pattern)
- apps/api/internal/service/admin_auth.go (service interface pattern)
- apps/api/internal/repository/postgres/user_repo.go (sqlc adapter pattern)
- apps/api/internal/repository/postgres/admin_repo.go (sqlc adapter pattern)
- apps/api/internal/repository/postgres/queries/users.sql (sqlc query pattern)
- apps/api/internal/repository/postgres/queries/admin_users.sql (sqlc query pattern)
- apps/api/internal/repository/postgres/migrations/00001_init.sql (migration pattern)
- apps/api/internal/repository/postgres/migrations/00079_admin_users.sql (migration pattern)
- apps/api/internal/middleware/admin_auth.go (middleware pattern)
- apps/api/internal/handler/health.go (handler pattern)
- apps/api/internal/handler/users.go (REST handler pattern)
- apps/api/gqlgen.yml (codegen config)
- apps/api/sqlc.yaml (sqlc config)
- libs/graphql/schema/schema.graphql (schema extension pattern)
- libs/graphql/schema/admin_auth.graphql (auth schema pattern)
- docs/prd-wave-details/waves/wave-01.md (prior wave detail)

## Selected Backend Wave Boundary
WAVE-02 adds Exercise CRUD (GraphQL) and ExerciseMedia association (REST) atop WAVE-01 infrastructure. This wave creates:
1. Migration files for exercises and exercise_media tables
2. GraphQL schema for Exercise type and CRUD operations (extend type Query/Mutation)
3. sqlc queries for exercise persistence
4. Exercise service layer
5. Exercise repository adapter
6. Exercise GraphQL resolvers
7. ExerciseMedia REST handler for upload/delete
8. Config extensions (possibly MediaConfig if not yet in WAVE-01)
9. PIN auth integration (reuse WAVE-01 guard)

## Neighboring Backend Wave Fit
- WAVE-01: provides the fitness-domain GraphQL schema extension pattern, PIN guard middleware pattern, sqlc/gqlgen codegen config for fitness domain, media REST scaffold, migration infrastructure
- WAVE-03: uses WAVE-02 exercises for exercise selector in workout diary. WAVE-02 must provide ListAllExercises query (simple list without pagination) for WAVE-03.

## Frontend Pages Context
PAGE-003 backend dependencies:
- GET exercises (list with filtering, single by ID)
- POST create exercise
- PUT update exercise
- DELETE soft-delete exercise
- POST upload exercise media
- DELETE exercise media

## Codebase Evidence

### Module Structure (what to create)
```
apps/api/internal/repository/postgres/queries/exercises.sql
apps/api/internal/repository/postgres/migrations/00080_exercises.sql
apps/api/internal/repository/postgres/migrations/00081_exercise_media.sql
apps/api/internal/repository/postgres/exercise_repo.go
apps/api/internal/service/exercise.go
apps/api/internal/handler/exercise_media.go
apps/api/internal/graph/exercise.resolvers.go (gqlgen-generated + custom)
libs/graphql/schema/exercises.graphql
```

### Config Changes (apps/api/internal/appconfig/config.go)
- If WAVE-01 already added MediaConfig, reuse it. If not, add MediaConfig with BasePath, MaxUploadSize.
- Add ExerciseConfig if needed (e.g., MaxExercisesPerUser or similar — likely not needed for MVP).

### main.go Wiring
- Add to Resolver: `ExerciseService *service.ExerciseService`
- Wire ExerciseService in main() after PIN auth service
- Register exercise REST media handler (new route group with PIN auth)
- ExerciseRepo instantiated from db.Pool

### gqlgen.yml
- Already includes `../../libs/graphql/schema/*.graphql` — exercises.graphql in the same directory will be auto-included. No config change needed.

### sqlc.yaml
- Already reads `internal/repository/postgres/queries` — new exercises.sql file in the same directory will be auto-included. No config change needed.

### GraphQL Schema Pattern (libs/graphql/schema/exercises.graphql)
Follow existing admin_auth.graphql pattern:
- type Exercise with all fields
- input CreateExerciseInput (name required, rest optional)
- input UpdateExerciseInput (all optional except id via URL param)
- union ExerciseResult = ExerciseSuccess | ValidationError | AuthError | NotFoundError
- union ExerciseListResult = ExerciseListSuccess | AuthError
- extend type Query { exercises(...): ExerciseListResult!, exercise(id: UUID!): ExerciseResult! }
- extend type Mutation { createExercise(input: CreateExerciseInput!): ExerciseResult!, updateExercise(id: UUID!, input: UpdateExerciseInput!): ExerciseResult!, deleteExercise(id: UUID!): DeleteExerciseResult! }

### Resolver Pattern
Follow existing schema.resolvers.go pattern:
- requireAdmin → will be replaced by requirePinAuth or similar from WAVE-01
- Mutation/Query resolver methods on *mutationResolver / *queryResolver
- Union result handling (switch on error types)

### Repository Pattern (apps/api/internal/repository/postgres/exercise_repo.go)
Follow user_repo.go pattern:
- SQL queries in exercises.sql (CRUD + list with cursor pagination + list all)
- Generated code in generated/ via sqlc
- ExerciseRepo struct with query interface narrowing
- UUID handling via pgtype.UUID
- Error mapping (pgx.ErrNoRows → nil, duplicate key → domain error)

### Service Pattern (apps/api/internal/service/exercise.go)
Follow admin_auth.go pattern:
- ExerciseRepository interface
- ExerciseService struct with repo dependency
- Methods: Create, GetByID, List (paginated), ListAll (simple for WAVE-03), Update, Delete (soft)
- Validation (name required, workingWeight > 0 if set)
- Log markers: [Exercise][action][BLOCK_NAME]

## Public Contracts
- ExerciseService: Create, GetByID, List, ListAll, Update, Delete
- ExerciseRepository interface
- ExerciseMediaRESTHandler: POST /api/v1/exercise-media (multipart upload), DELETE /api/v1/exercise-media/{id}

## Generated Artifact Impact
- gqlgen: will generate new model types, resolver methods in apps/api/internal/graph/
- sqlc: will generate query functions in apps/api/internal/repository/postgres/generated/
- Both codegen configs auto-include new files in their directories — no config changes needed

## Integration Points
- WAVE-01 media REST scaffold: WAVE-02 ExerciseMedia handler reuses the file storage path from WAVE-01's media config. Media files stored at configured base path, metadata in exercise_media table.
- WAVE-01 PIN auth guard: WAVE-02 GraphQL resolvers and REST handlers reuse the PIN auth middleware.
- WAVE-03 exercise selector: WAVE-02 provides ListAllExercises for simple exercise list.

## Likely Graph Deltas
- M-API: adding ExerciseService dependency to Resolver struct
- M-GRAPHQL-SCHEMA: adding exercises.graphql schema file
- M-EXERCISE (new module): Exercise service, repository, queries
- M-EXERCISE-MEDIA (new module): ExerciseMedia handler, REST routes
- V-M-API: adding exercise resolver, repo, service tests

## Unsupported Assumptions
- Assumes WAVE-01 media REST scaffold provides file storage and retrieval at a well-known path. If WAVE-01 uses a different approach (e.g., S3-compatible), adjust accordingly.
- No pagination contract for exercise list with cursor — using same keyset pagination as existing users list.
- ExerciseMedia file storage path derivation: exercise_{exerciseId}/{uuid}.{ext} for file organization.