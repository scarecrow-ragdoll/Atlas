<!--
FILE: docs/superpowers/specs/2026-06-19-wave-03-workout-diary-design.md
VERSION: 1.0.0
START_MODULE_CONTRACT
  PURPOSE: Define the detailed WAVE-03 Workout Diary design for Atlas, covering DailyLog, workout exercises, workout sets, aggregate versioning, GraphQL operations, WAVE-04 compatibility, validation, and verification.
  SCOPE: DailyLog canonical daily container, strength workout logging, exercise-level notes, workout set CRUD, exercise/set ordering, working weight snapshots, and optimistic concurrency; excludes cardio CRUD, body weight tracking, charts, AI export, backup/import, and frontend implementation.
  DEPENDS: docs/prd-waves/waves/wave-03.md, docs/product-verified/features/workout-diary.md, docs/product-verified/domain-model.md, docs/technical-verified/api-contracts.md, docs/superpowers/specs/2026-06-19-wave-01-foundation-design.md, docs/superpowers/specs/2026-06-19-wave-02-exercise-library-design.md.
  LINKS: M-API / V-M-API / WAVE-03.
  ROLE: DOC
  MAP_MODE: SUMMARY
END_MODULE_CONTRACT
START_CHANGE_SUMMARY
  LAST_CHANGE: 1.0.0 - Initial WAVE-03 design after brainstorming and user approval.
END_CHANGE_SUMMARY
-->

# WAVE-03: Workout Diary - Detailed Design

**Date:** 2026-06-19
**Status:** Approved
**Approach:** Granular aggregate mutations with `expectedVersion`

## 1. Scope Boundary

WAVE-03 implements the canonical daily strength workout aggregate:

- `daily_logs`
- `workout_exercises`
- `workout_sets`
- DailyLog GraphQL queries and mutations
- adding, updating, removing, and reordering exercises inside a day
- adding, updating, removing, and reordering workout sets
- exercise-level notes inside the day
- daily notes on `daily_logs`
- working weight snapshots when an exercise is added to a day
- optimistic concurrency/versioning for the DailyLog aggregate

WAVE-03 does not implement:

- `cardio_entries`
- cardio CRUD mutations
- cardio GraphQL types or placeholder fields
- cardio validation
- `CardioType` or `HeartRateZone`
- body weight persistence
- charts visualization
- AI export
- frontend implementation

## 2. Approved Product Decisions

1. DailyLog is the canonical shared daily container. The old `WorkoutDay` name must not be used in implementation names unless only quoting legacy source text.
2. WAVE-03 must not implement Cardio CRUD. WAVE-04 owns cardio tables, schema, resolvers, validation, enums, and attachment behavior.
3. WAVE-03 must expose granular GraphQL mutations. Every mutation that changes the DailyLog aggregate requires `expectedVersion`.
4. WAVE-03 includes `daily_logs.notes`, but body weight stays out of WAVE-03 and belongs to WAVE-04 body tracking.
5. A single DailyLog may contain the same `exercise_id` multiple times. There is no unique constraint on `(daily_log_id, exercise_id)`.
6. Empty DailyLog records are retained after removing all strength workout data and notes.
7. Set validation:
   - `weight > 0`
   - `reps > 0`
   - `rpe` optional, `1..10` when present
   - `rir` optional, `0..10` when present
   - `notes` optional
8. WAVE-03 conflict handling uses optimistic version checks, not last-write-wins.

## 3. Data Model

### `daily_logs`

| Column | Type | Constraints |
| --- | --- | --- |
| id | UUID | PK, default `gen_random_uuid()` |
| user_id | UUID | NOT NULL, FK -> `atlas_users(id)` |
| date | DATE | NOT NULL |
| notes | TEXT | nullable |
| version | INTEGER | NOT NULL DEFAULT 0 |
| created_at | TIMESTAMPTZ | NOT NULL DEFAULT now() |
| updated_at | TIMESTAMPTZ | NOT NULL DEFAULT now() |

Constraints and indexes:

- `UNIQUE(user_id, date)`
- `CHECK(version >= 0)`
- `idx_daily_logs_user_date (user_id, date)`
- optional range helper: `idx_daily_logs_user_date_desc (user_id, date DESC)`

### `workout_exercises`

| Column | Type | Constraints |
| --- | --- | --- |
| id | UUID | PK, default `gen_random_uuid()` |
| user_id | UUID | NOT NULL, FK -> `atlas_users(id)` |
| daily_log_id | UUID | NOT NULL, FK -> `daily_logs(id)` ON DELETE CASCADE |
| exercise_id | UUID | NOT NULL, FK -> `exercises(id)` ON DELETE RESTRICT |
| position | INTEGER | NOT NULL |
| working_weight_snapshot | REAL | nullable |
| notes | TEXT | nullable |
| created_at | TIMESTAMPTZ | NOT NULL DEFAULT now() |
| updated_at | TIMESTAMPTZ | NOT NULL DEFAULT now() |

Constraints and indexes:

- `CHECK(position > 0)`
- `CHECK(working_weight_snapshot IS NULL OR working_weight_snapshot > 0)`
- `UNIQUE(daily_log_id, position)`
- `idx_workout_exercises_user_daily_log (user_id, daily_log_id)`
- `idx_workout_exercises_exercise (exercise_id)`

No `UNIQUE(daily_log_id, exercise_id)`: duplicate exercise instances are allowed.

### `workout_sets`

| Column | Type | Constraints |
| --- | --- | --- |
| id | UUID | PK, default `gen_random_uuid()` |
| workout_exercise_id | UUID | NOT NULL, FK -> `workout_exercises(id)` ON DELETE CASCADE |
| set_number | INTEGER | NOT NULL |
| weight | REAL | NOT NULL |
| reps | INTEGER | NOT NULL |
| rpe | REAL | nullable |
| rir | INTEGER | nullable |
| notes | TEXT | nullable |
| created_at | TIMESTAMPTZ | NOT NULL DEFAULT now() |
| updated_at | TIMESTAMPTZ | NOT NULL DEFAULT now() |

Constraints and indexes:

- `CHECK(set_number > 0)`
- `CHECK(weight > 0)`
- `CHECK(reps > 0)`
- `CHECK(rpe IS NULL OR (rpe >= 1 AND rpe <= 10))`
- `CHECK(rir IS NULL OR (rir >= 0 AND rir <= 10))`
- `UNIQUE(workout_exercise_id, set_number)`
- `idx_workout_sets_workout_exercise (workout_exercise_id)`

`workout_sets` inherits user scope through `workout_exercises -> daily_logs`.

## 4. Versioning And Concurrency

DailyLog is the versioned aggregate root. Mutations that change notes, workout exercises, workout exercise order, workout exercise notes, sets, or set order must:

1. Resolve or create the DailyLog aggregate.
2. Lock the relevant DailyLog row inside a database transaction.
3. Compare `daily_logs.version` with the input `expectedVersion`.
4. If mismatched, return `ConflictError` with `currentVersion` and the current DailyLog.
5. Apply the mutation.
6. Increment `daily_logs.version` by 1 and update `daily_logs.updated_at`.
7. Return the complete current DailyLog aggregate.

For a date with no existing DailyLog, the client uses `expectedVersion: 0`. The first successful mutation creates the DailyLog and returns `version: 1`.

WAVE-04 cardio mutations must reuse the same DailyLog versioning rule. When cardio entries change, WAVE-04 should increment the attached DailyLog version.

## 5. GraphQL Schema

WAVE-03 extends the Atlas GraphQL schema under `/graphql/atlas`.

`Date` is a domain calendar date, not a timestamp. It must parse and serialize strict `YYYY-MM-DD` values without timezone conversion. Add an Atlas gqlgen binding for a small `models.Date` type rather than reusing `Time`.

### Queries

```graphql
type Query {
  dailyLog(date: Date!): DailyLogResult!
  dailyLogs(from: Date!, to: Date!): [DailyLogSummary!]!
}
```

`dailyLog(date:)` returns the existing DailyLog or an empty date slot representation without creating a row. The first mutation creates the row.

### Mutations

```graphql
type Mutation {
  updateDailyLogNotes(date: Date!, expectedVersion: Int!, notes: String): DailyLogResult!
  addWorkoutExercise(date: Date!, expectedVersion: Int!, input: AddWorkoutExerciseInput!): DailyLogResult!
  updateWorkoutExercise(id: ID!, expectedVersion: Int!, input: UpdateWorkoutExerciseInput!): DailyLogResult!
  removeWorkoutExercise(id: ID!, expectedVersion: Int!): DailyLogResult!
  reorderWorkoutExercises(date: Date!, expectedVersion: Int!, orderedIds: [ID!]!): DailyLogResult!
  addWorkoutSet(workoutExerciseId: ID!, expectedVersion: Int!, input: AddWorkoutSetInput!): DailyLogResult!
  updateWorkoutSet(id: ID!, expectedVersion: Int!, input: UpdateWorkoutSetInput!): DailyLogResult!
  removeWorkoutSet(id: ID!, expectedVersion: Int!): DailyLogResult!
  reorderWorkoutSets(workoutExerciseId: ID!, expectedVersion: Int!, orderedIds: [ID!]!): DailyLogResult!
}
```

### Types

```graphql
scalar Date

type DailyLog {
  id: ID!
  userId: ID!
  date: Date!
  notes: String
  version: Int!
  workoutExercises: [WorkoutExercise!]!
  createdAt: Time!
  updatedAt: Time!
}

type WorkoutExercise {
  id: ID!
  userId: ID!
  dailyLogId: ID!
  exerciseId: ID!
  exercise: Exercise!
  position: Int!
  workingWeightSnapshot: Float
  notes: String
  sets: [WorkoutSet!]!
  createdAt: Time!
  updatedAt: Time!
}

type WorkoutSet {
  id: ID!
  workoutExerciseId: ID!
  setNumber: Int!
  weight: Float!
  reps: Int!
  rpe: Float
  rir: Int
  notes: String
  createdAt: Time!
  updatedAt: Time!
}

type DailyLogSummary {
  id: ID!
  date: Date!
  version: Int!
  workoutExerciseCount: Int!
  workoutSetCount: Int!
  totalVolume: Float!
  updatedAt: Time!
}
```

### Inputs

```graphql
input AddWorkoutExerciseInput {
  exerciseId: ID!
  position: Int
  notes: String
}

input UpdateWorkoutExerciseInput {
  position: Int
  notes: String
}

input AddWorkoutSetInput {
  setNumber: Int
  weight: Float!
  reps: Int!
  rpe: Float
  rir: Int
  notes: String
}

input UpdateWorkoutSetInput {
  setNumber: Int
  weight: Float
  reps: Int
  rpe: Float
  rir: Int
  notes: String
}
```

If `position` or `setNumber` is omitted on add, the service appends to the end.

### Result And Error Types

```graphql
type DailyLogResult {
  dailyLog: DailyLog
  validationError: DailyLogValidationError
  notFoundError: DailyLogNotFoundError
  conflictError: DailyLogConflictError
  authError: DailyLogAuthError
}

type DailyLogValidationError {
  message: String!
  code: DailyLogErrorCode!
}

type DailyLogNotFoundError {
  message: String!
  code: DailyLogErrorCode!
}

type DailyLogConflictError {
  message: String!
  code: DailyLogErrorCode!
  currentVersion: Int!
  currentDailyLog: DailyLog
}

type DailyLogAuthError {
  message: String!
  code: DailyLogErrorCode!
}

enum DailyLogErrorCode {
  VALIDATION_ERROR
  NOT_FOUND
  CONFLICT
  AUTH_ERROR
  INTERNAL_ERROR
}
```

## 6. Service And Repository Design

### New Models

`apps/api/internal/atlas/models/workout.go` owns:

- `Date`
- `DailyLogRecord`
- `DailyLog`
- `DailyLogSummary`
- `WorkoutExerciseRecord`
- `WorkoutExercise`
- `WorkoutSetRecord`
- `WorkoutSet`
- GraphQL input/result/error structs

### Repository

`apps/api/internal/atlas/repository/postgres/workout_repo.go` owns sqlc-backed data access:

- `GetDailyLogByDate(ctx, userID, date)`
- `GetDailyLogAggregate(ctx, userID, dailyLogID)`
- `GetOrCreateDailyLogByDate(ctx, userID, date)`
- `LockDailyLogByID(ctx, userID, id)`
- `LockDailyLogByDate(ctx, userID, date)`
- `IncrementDailyLogVersion(ctx, userID, id)`
- `UpdateDailyLogNotes(ctx, userID, date, notes)`
- `ListDailyLogSummaries(ctx, userID, from, to)`
- `AddWorkoutExercise(ctx, userID, dailyLogID, input)`
- `UpdateWorkoutExercise(ctx, userID, id, input)`
- `DeleteWorkoutExercise(ctx, userID, id)`
- `ReorderWorkoutExercises(ctx, userID, dailyLogID, orderedIDs)`
- `AddWorkoutSet(ctx, workoutExerciseID, input)`
- `UpdateWorkoutSet(ctx, id, input)`
- `DeleteWorkoutSet(ctx, id)`
- `ReorderWorkoutSets(ctx, workoutExerciseID, orderedIDs)`

The repository should use transactions for aggregate mutations. If keeping transactions outside sqlc adapters is cleaner, expose transaction-bound query helpers instead of burying cross-row operations in single functions.

### Service

`apps/api/internal/atlas/service/workout.go` owns business behavior:

- current user scoping
- date validation and normalization
- DailyLog `GetOrCreateByDate`
- optimistic concurrency checks
- exercise existence checks through WAVE-02 exercise repository/service
- working weight snapshot capture
- position and set-number append/reindex behavior
- validation-to-error mapping
- aggregate reload after mutation

Recommended service errors:

- `ErrDailyLogConflict`
- `ErrDailyLogNotFound`
- `ErrWorkoutExerciseNotFound`
- `ErrWorkoutSetNotFound`
- `ErrExerciseNotFound`
- `ErrInvalidDateRange`
- `ErrInvalidPosition`
- `ErrInvalidSetNumber`
- `ErrInvalidWeight`
- `ErrInvalidReps`
- `ErrInvalidRPE`
- `ErrInvalidRIR`

## 7. Mutation Behavior

### `dailyLog(date:)`

Returns an existing DailyLog if present. If absent, returns a date slot with no `id` only if gqlgen/model design supports it cleanly; otherwise returns `dailyLog: null` and no error. It must not create a database row.

### `updateDailyLogNotes`

Creates the DailyLog if absent and `expectedVersion` is `0`. Sets `notes` to the provided value, including null/empty for clearing.

### `addWorkoutExercise`

Creates the DailyLog if absent. Validates that the exercise exists for the current user. Captures `Exercise.workingWeight` into `working_weight_snapshot`. Inserts at requested position or appends. Reindexes positions to stay contiguous.

### `updateWorkoutExercise`

Updates notes and optionally position. Does not change `exercise_id` or `working_weight_snapshot`; replacing an exercise should be remove + add so historical snapshots stay honest.

### `removeWorkoutExercise`

Deletes the workout exercise and cascades its sets. Reindexes remaining workout exercise positions. Does not delete the DailyLog if empty.

### `reorderWorkoutExercises`

Requires `orderedIds` to exactly match the current workout exercise IDs for the DailyLog. Missing, duplicate, foreign, or extra IDs return validation errors.

### `addWorkoutSet`

Validates the parent workout exercise belongs to the current user. Appends if `setNumber` is omitted; otherwise inserts at the requested set number and reindexes to contiguous order.

### `updateWorkoutSet`

Updates weight, reps, RPE, RIR, notes, and optional set number. Reindex if the set number changes.

### `removeWorkoutSet`

Deletes the set and reindexes remaining sets for that workout exercise.

### `reorderWorkoutSets`

Requires `orderedIds` to exactly match the current set IDs for the workout exercise.

## 8. WAVE-04 Compatibility Contract

WAVE-03 must make DailyLog reusable without exposing unfinished cardio API:

- `daily_logs.id`, `user_id`, `date`, `version`, `created_at`, `updated_at` are mandatory and stable.
- `UNIQUE(user_id, date)` is mandatory.
- service/repository supports `GetOrCreateByDate`.
- version increment behavior is aggregate-level and reusable by future child-domain mutations.
- WAVE-04 cardio changes should lock the same DailyLog and increment version on cardio create/update/delete.
- WAVE-03 GraphQL leaves cardio fields out until WAVE-04 implements real behavior.

## 9. Implementation Slices

1. `SLICE-W03-001` - DB migrations for `daily_logs`, `workout_exercises`, and `workout_sets`.
2. `SLICE-W03-002` - sqlc queries for DailyLog aggregate reads, locks, version increments, workout exercise CRUD/order, and workout set CRUD/order.
3. `SLICE-W03-003` - Atlas workout models.
4. `SLICE-W03-004` - PostgreSQL repository with transaction-safe aggregate operations.
5. `SLICE-W03-005` - Workout service with validation, working weight snapshots, version checks, and reindexing.
6. `SLICE-W03-006` - Atlas GraphQL schema additions.
7. `SLICE-W03-007` - GraphQL resolvers and resolver wiring.
8. `SLICE-W03-008` - Focused tests, codegen, and WAVE-04 compatibility proof.

## 10. Verification Strategy

Focused checks for WAVE-03:

- `bunx nx run api:codegen`
- `bunx nx run api:codegen:atlas`
- `bunx nx test api`
- package-scoped Go tests while iterating, such as:
  - `cd apps/api && go test ./internal/atlas/service -run '(?i)workout' -count=1`
  - `cd apps/api && go test ./internal/atlas/repository/postgres -run '(?i)workout|dailylog' -count=1`
  - `cd apps/api && go test ./internal/atlas/graph/resolver -run '(?i)workout|dailylog' -count=1`

Required behavior coverage:

- DailyLog creates on first mutation and not on query.
- `UNIQUE(user_id, date)` prevents duplicate DailyLogs.
- stale `expectedVersion` returns conflict with current version.
- notes update increments version.
- adding exercise captures WAVE-02 working weight snapshot.
- exercise update after snapshot does not mutate historical snapshot.
- duplicate exercise instances in one DailyLog are allowed.
- exercise position insert, reorder, delete, and reindex are deterministic.
- set add/update/delete/reorder validates weight, reps, RPE, RIR, and reindexes set numbers.
- user-scoped isolation prevents cross-user access.
- deleting a workout exercise cascades sets.
- deleting all exercises leaves DailyLog intact.
- WAVE-04 compatibility: `GetOrCreateByDate` and version increment can be called without strength-specific assumptions.

End-of-wave checks after generated artifacts are synchronized:

- `bunx nx run api:codegen`
- `bunx nx run api:codegen:atlas`
- `bunx nx test api`
- `bunx nx build api`

Broader root gates are not part of the WAVE-03 active-development loop unless the handoff phase explicitly asks for them.

## 11. Deliverables

- migration files after WAVE-02 numbering, expected next:
  - `apps/api/internal/repository/postgres/migrations/00083_daily_logs.sql`
  - `apps/api/internal/repository/postgres/migrations/00084_workout_exercises.sql`
  - `apps/api/internal/repository/postgres/migrations/00085_workout_sets.sql`
- `apps/api/internal/repository/postgres/queries/workouts.sql`
- generated sqlc output
- `apps/api/internal/atlas/models/workout.go`
- `apps/api/internal/atlas/repository/postgres/workout_repo.go`
- `apps/api/internal/atlas/service/workout.go`
- `apps/api/internal/atlas/graph/schema/workouts.graphql`
- resolver additions under `apps/api/internal/atlas/graph/resolver`
- focused tests for repository, service, and resolvers

## 12. Explicit Non-Goals

- No frontend route or UI implementation.
- No public web changes.
- No web-admin changes.
- No cardio placeholder data or fake empty GraphQL fields.
- No starter workout template feature.
- No automatic working weight progression.
- No charts or e1RM chart endpoint.
- No AI export payload assembly.
