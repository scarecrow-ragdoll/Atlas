# WAVE-02 architecture-codebase Planner Attempt 2

## Cycle 1 Reviewer Feedback Addressed

### 1. WAVE-01 Dependency Contract (explicit)
WAVE-02 explicitly depends on these contracts from WAVE-01:
- **requirePinAuth(ctx) error** — Function in `apps/api/internal/graph/` that extracts PIN-authenticated principal from context. Analogous to existing `requireAdmin`. Returns error with `ErrAdminAuth` semantics. If WAVE-01 uses a different name, WAVE-02 adapts.
- **MediaConfig** — Config struct in `apps/api/internal/appconfig/config.go` with `BasePath string` for file storage root, `MaxUploadSize int64` for global upload limit.
- **POST /api/v1/media/upload** — Multipart file upload, returns `{ "id": "uuid" }`. Stores file at configured path.
- **GET /api/v1/media/{id}** — File download by media UUID. Returns file with correct Content-Type.
- **Fitness GraphQL common types** — `ValidationError`, `AuthError`, `NotFoundError` types defined in fitness-domain GraphQL schema (fitness.graphql or common.graphql). These are reused by WAVE-02 exercises.graphql.

### 2. Common GraphQL Types
WAVE-02 uses the same `ValidationError`, `AuthError`, `NotFoundError` union result types that WAVE-01 establishes for the fitness domain. No need to redefine. If WAVE-01 does not provide these, WAVE-02 defines them in exercises.graphql with the same structure.

### 3. ExerciseMedia Route Registration in main.go
```go
// After WAVE-01 PIN middleware setup:
// PIN-protected exercise media REST routes
r.Group(func(fitness chi.Router) {
    fitness.Use(pinAuthMiddleware)  // From WAVE-01
    fitness.Post("/api/v1/exercise-media", exerciseMediaHandler.Upload)
    fitness.Get("/api/v1/exercise-media/{id}", exerciseMediaHandler.Download)
    fitness.Delete("/api/v1/exercise-media/{id}", exerciseMediaHandler.Delete)
})

// PIN-protected GraphQL for fitness domain
r.Group(func(fitnessGQL chi.Router) {
    fitnessGQL.Use(pinAuthMiddleware)
    fitnessGQL.Handle("/graphql", srv)  // Same /graphql endpoint, PIN middleware
})
```

The PIN auth route group must be distinct from the admin cookie-based group. Both share the same /graphql path — routing is handled by middleware type check (admin cookie middleware for admin session, PIN token middleware for fitness session). Alternatively, use separate GraphQL endpoints: /graphql (admin), /fitness/graphql (PIN). This depends on WAVE-01's routing decision.

### 4. Migration Numbering
WAVE-02 proposes 00080_exercises.sql and 00081_exercise_media.sql. These numbers assume WAVE-01 uses 00079 as the last migration (current state). If WAVE-01 adds more migrations, WAVE-02 numbers are adjusted to follow sequentially.

### 5. PIN Auth Guard Integration
The WAVE-01 PIN auth middleware pattern is expected to provide:
- `requirePinAuth(ctx) (*PinPrincipal, error)` — analogous to `requireAdmin`
- If WAVE-01 wraps this at the middleware level (chi middleware), GraphQL resolvers use the same middleware pattern. WAVE-02 resolvers call `requirePinAuth(ctx)` for each mutation/query guard.

## Updated Module Structure
```
apps/api/internal/repository/postgres/migrations/00080_exercises.sql
apps/api/internal/repository/postgres/migrations/00081_exercise_media.sql
apps/api/internal/repository/postgres/queries/exercises.sql
apps/api/internal/repository/postgres/exercise_repo.go
apps/api/internal/service/exercise.go
apps/api/internal/handler/exercise_media.go
apps/api/internal/graph/exercise.resolvers.go
libs/graphql/schema/exercises.graphql
```

## Updated main.go Wiring Pattern
```go
// In main():
exerciseRepo := postgres.NewExerciseRepo(db.Pool)
exerciseService := service.NewExerciseService(exerciseRepo)
exerciseMediaHandler := handler.NewExerciseMediaHandler(exerciseService, mediaConfig, l)

// Resolver extension:
resolver := &graph.Resolver{
    UserService:      userService,
    AdminAuthService: adminAuthService,
    ExerciseService:  exerciseService,  // NEW
}

// Route registration:
// Admin group (unchanged)
// Fitness PIN-protected group (new):
pinGroup := r.Group(func(pin chi.Router) {
    pin.Use(pinAuthMiddleware)  // WAVE-01
    pin.Handle("/graphql", srv)
    pin.Post("/api/v1/exercise-media", exerciseMediaHandler.Upload)
    pin.Get("/api/v1/exercise-media/{id}", exerciseMediaHandler.Download)
    pin.Delete("/api/v1/exercise-media/{id}", exerciseMediaHandler.Delete)
})
```

## Resolver Dependency Injection
Resolver struct gains:
```go
type Resolver struct {
    UserService      *service.UserService
    AdminAuthService *service.AdminAuthService
    ExerciseService  *service.ExerciseService  // NEW
}
```

## CORS Configuration
Exercise media REST endpoints use the same CORS config that WAVE-01 establishes for fitness endpoints. If WAVE-01 provides a `fitnessCORS` or `publicCORS` config, WAVE-02 reuses it. The existing `publicCORS` from main.go (allows public API endpoints) is appropriate for fitness REST endpoints with PIN auth.

## Generated Artifact Impact
- gqlgen: generates Exercise model, ExerciseResult union, resolver stubs
- sqlc: generates exercise CRUD query functions
- Both auto-discover new files via directory glob patterns

## GraphQL Schema (exercises.graphql)
```graphql
type Exercise {
  id: UUID!
  name: String!
  muscleGroups: [String!]!
  description: String
  personalNotes: String
  workingWeight: Float
  isActive: Boolean!
  media: [ExerciseMedia!]!
  createdAt: DateTime!
  updatedAt: DateTime!
}

type ExerciseMedia {
  id: UUID!
  exerciseId: UUID!
  mediaType: String!
  filePath: String!
  originalFileName: String!
  mimeType: String!
  sizeBytes: Int!
  createdAt: DateTime!
}

input CreateExerciseInput {
  name: String!
  muscleGroups: [String!]
  description: String
  personalNotes: String
  workingWeight: Float
  isActive: Boolean
}

input UpdateExerciseInput {
  name: String
  muscleGroups: [String!]
  description: String
  personalNotes: String
  workingWeight: Float
}

type ExerciseSuccess {
  exercise: Exercise!
}

type ExerciseListSuccess {
  items: [Exercise!]!
  totalCount: Int!
}

type DeleteExerciseSuccess {
  ok: Boolean!
}

union ExerciseResult = ExerciseSuccess | ValidationError | AuthError | NotFoundError
union ExerciseListResult = ExerciseListSuccess | AuthError
union DeleteExerciseResult = DeleteExerciseSuccess | ValidationError | AuthError | NotFoundError

extend type Query {
  exercises(first: Int, after: String, includeInactive: Boolean): ExerciseListResult!
  exercise(id: UUID!): ExerciseResult!
  allExercises(includeInactive: Boolean): [Exercise!]!
}

extend type Mutation {
  createExercise(input: CreateExerciseInput!): ExerciseResult!
  updateExercise(id: UUID!, input: UpdateExerciseInput!): ExerciseResult!
  deleteExercise(id: UUID!): DeleteExerciseResult!
}
```