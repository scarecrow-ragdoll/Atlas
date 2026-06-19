# WAVE-03 architecture-codebase Planner Attempt 1

## Sources Read
- apps/api/cmd/server/main.go
- apps/api/internal/appconfig/config.go
- apps/api/internal/middleware/admin_auth.go
- apps/api/internal/service/admin_auth.go
- apps/api/internal/repository/postgres/user_repo.go
- apps/api/internal/repository/redis/admin_session_store.go
- apps/api/gqlgen.yml
- apps/api/sqlc.yaml
- libs/graphql/schema/schema.graphql
- libs/graphql/schema/admin_auth.graphql
- libs/graphql/schema/common.graphql
- libs/graphql/schema/user.graphql
- apps/api/internal/repository/postgres/migrations/00079_admin_users.sql
- apps/api/internal/repository/postgres/queries/users.sql
- docs/prd-wave-details/codebase-fit.md
- docs/technical-verified/architecture-and-boundaries.md
- docs/technical-verified/api-contracts.md
- docs/technical-verified/data-contracts.md

## Selected Backend Wave Boundary
WAVE-03 adds four new database tables (daily_logs, workout_exercises, workout_sets, cardio_entries), their sqlc query definitions, GraphQL schema types and operations, repository adapter, service layer, and GraphQL resolvers. All endpoints protected by WAVE-01 PIN auth middleware.

## Neighboring Backend Wave Fit
- WAVE-01: Provides PIN middleware, common error types, config extension, codegen infra. WAVE-03 reuses all.
- WAVE-02: Provides allExercises query. WAVE-03 uses exerciseId FK to reference exercises. Exercise.workingWeight is read at snapshot time.
- WAVE-04: CardioEntry entity shared boundary. WAVE-03 creates CardioEntry linked to DailyLog (dailyLogId required). WAVE-04 adds standalone CardioEntry and body tracking tables.

## Frontend Pages Context
- PAGE-002 (Workout Diary): consumes all WAVE-03 GraphQL operations. Backend dependency only.

## Codebase Evidence

### Existing Architecture Patterns
The Go API follows a consistent layered architecture:
- `main.go` (wiring) -> `appconfig/` (config) -> middleware -> `handler/` (REST) -> `service/` (business logic) -> `repository/postgres/` (sqlc queries + adapter)
- GraphQL resolvers in `internal/graph/` are wired through the resolver struct in main.go
- gqlgen.yml auto-discovers schema files via `../../libs/graphql/schema/*.graphql` glob
- sqlc.yaml reads `internal/repository/postgres/migrations` for schema and `queries/` for query SQL
- Migration pattern: sequential `NNNNN_name.sql` with +goose Up/Down markers

### File Touchpoints for WAVE-03
1. **DB Migrations**: `apps/api/internal/repository/postgres/migrations/00082_daily_logs.sql` (starts after WAVE-02's 00080/00081)
2. **DB Migrations**: `apps/api/internal/repository/postgres/migrations/00083_workout_exercises.sql`
3. **DB Migrations**: `apps/api/internal/repository/postgres/migrations/00084_workout_sets.sql`
4. **DB Migrations**: `apps/api/internal/repository/postgres/migrations/00085_cardio_entries.sql`
5. **SQLC Queries**: `apps/api/internal/repository/postgres/queries/daily_logs.sql`
6. **SQLC Queries**: `apps/api/internal/repository/postgres/queries/workout_exercises.sql`
7. **SQLC Queries**: `apps/api/internal/repository/postgres/queries/workout_sets.sql`
8. **SQLC Queries**: `apps/api/internal/repository/postgres/queries/cardio_entries.sql`
9. **Repository**: `apps/api/internal/repository/postgres/daily_log_repo.go`
10. **Repository**: `apps/api/internal/repository/postgres/workout_exercise_repo.go`
11. **Repository**: `apps/api/internal/repository/postgres/workout_set_repo.go`
12. **Repository**: `apps/api/internal/repository/postgres/cardio_entry_repo.go`
13. **Service**: `apps/api/internal/service/workout.go` (single service for all workout domain operations)
14. **GraphQL Schema**: `libs/graphql/schema/workout.graphql` (DailyLog, WorkoutExercise, WorkoutSet, CardioEntry types + operations)
15. **GraphQL Resolvers**: `apps/api/internal/graph/workout.resolvers.go`
16. **Main wiring**: `apps/api/cmd/server/main.go` (wire repos, service, resolvers)
17. **Config**: `apps/api/internal/appconfig/config.go` (no new config needed — reuses WAVE-01 sessions)

### Generated Artifact Impact
- gqlgen: auto-discovers new workout.graphql via glob. Generates DailyLog, WorkoutExercise, WorkoutSet, CardioEntry model structs and resolver stubs.
- sqlc: auto-discovers new query .sql files via glob. Generates CRUD functions for all four tables.

### Integration Points
- WAVE-01 PIN auth middleware: wraps fitness GraphQL endpoint and guards all mutations/queries
- WAVE-02 allExercises query: used within workout service to read Exercise.workingWeight for snapshot
- WAVE-01 common types: ValidationError, AuthError, NotFoundError reused in union results

## Proposed Details

### Implementation Slices

| Slice ID | Name | Description |
| --- | --- | --- |
| SLICE-W03-001 | DB migration: daily_logs | Create goose migration 00082_daily_logs.sql with id, user_id, date (UNIQUE per user), notes, body_weight, created_at, updated_at. Index on (user_id, date). |
| SLICE-W03-002 | DB migration: workout_exercises | Create goose migration 00083_workout_exercises.sql with id, user_id, daily_log_id FK, exercise_id FK, order, working_weight_snapshot, notes, created_at, updated_at. Indexes and FK cascade. |
| SLICE-W03-003 | DB migration: workout_sets | Create goose migration 00084_workout_sets.sql with id, workout_exercise_id FK, set_number, weight REAL, reps INT, rpe REAL NULL, rir INT NULL, notes TEXT NULL, created_at, updated_at. FK cascade. |
| SLICE-W03-004 | DB migration: cardio_entries | Create goose migration 00085_cardio_entries.sql with id, user_id, daily_log_id FK, cardio_type, duration_minutes, avg_pulse INT NULL, heart_rate_zone INT NULL, notes TEXT NULL, created_at, updated_at. FK cascade. |
| SLICE-W03-005 | sqlc queries: daily_logs | CRUD queries for daily_logs table: get by user_id + date, upsert, delete. |
| SLICE-W03-006 | sqlc queries: workout_exercises | CRUD queries: list by daily_log_id ordered by order, create, update order/notes, delete (cascading to sets). |
| SLICE-W03-007 | sqlc queries: workout_sets | CRUD queries: list by workout_exercise_id ordered by set_number, create (auto set_number), update, delete. |
| SLICE-W03-008 | sqlc queries: cardio_entries | CRUD queries: list by daily_log_id, create, update, delete. |
| SLICE-W03-009 | DailyLog repository | Repository adapter for daily_logs with GetByDate(userID, date), Upsert, Delete. Uses sqlc-generated queries. |
| SLICE-W03-010 | WorkoutExercise repository | Repository adapter for workout_exercises with ListByDailyLog, Create, Update, Delete. |
| SLICE-W03-011 | WorkoutSet repository | Repository adapter for workout_sets with ListByWorkoutExercise, Create (auto-number), Update, Delete. |
| SLICE-W03-012 | CardioEntry repository | Repository adapter for cardio_entries with ListByDailyLog, Create, Update, Delete. |
| SLICE-W03-013 | Workout service | Transport-neutral service: DailyLogByDate (with full nested data), UpsertDailyLog, AddExercise (snapshot weight), RemoveExercise, UpdateExercise, AddSet, UpdateSet, RemoveSet, AddCardio, RemoveCardio, DeleteDailyLog. Validates inputs, calls WAVE-02 allExercises for snapshot. |
| SLICE-W03-014 | GraphQL schema | Add workout.graphql with DailyLog, WorkoutExercise, WorkoutSet, CardioEntry types, input types, union results, extend Query/Mutation. |
| SLICE-W03-015 | GraphQL resolvers + main wiring | Implement resolvers for all workout operations. Wire DailyLogRepo, WorkoutExerciseRepo, WorkoutSetRepo, CardioEntryRepo, WorkoutService into main.go. Register PIN-protected fitness GraphQL endpoint. |

## Acceptance Criteria Contributions
Contributed via planner-product-ac. All 30 ACs (AC-W03-001 through AC-W03-030) are supported by the proposed implementation slices.

## Exit Criteria Contributions
- EC-W03-001 through EC-W03-018 (see testing-exit planner)

## Verification Contributions
- TEST-W03-001 through TEST-W03-022 (see testing-exit planner)

## Risks And Rollback
- Migration numbering: must start at 00082 after WAVE-02's 00080/00081. If WAVE-01 adds intervening migrations, renumber.
- FK to exercises table: exercise_id FK prevents deletion of exercises with workout history. WAVE-02 soft delete (isActive) is compatible.
- FK to daily_logs: CASCADE delete ensures no orphaned records.
- Rollback: goose down migrations for 00082-00085 remove all WAVE-03 tables.

## Questions Raised
- DQ-W03-002: Migration numbering resolved at 00082.
- DQ-W03-006: Cascade delete confirmed CASCADE.

## Traceability Candidates
- apps/api/cmd/server/main.go: wiring pattern
- apps/api/internal/repository/postgres/user_repo.go: repository adapter pattern
- apps/api/internal/service/admin_auth.go: service layer pattern
- libs/graphql/schema/admin_auth.graphql: union result type and extend pattern
- apps/api/gqlgen.yml: schema glob
- apps/api/sqlc.yaml: query glob
