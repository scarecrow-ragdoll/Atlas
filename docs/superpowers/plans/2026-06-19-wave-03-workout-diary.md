# WAVE-03 Workout Diary Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the Atlas WAVE-03 strength workout diary backend: DailyLog, workout exercises, workout sets, optimistic aggregate versioning, and Atlas GraphQL operations.

**Architecture:** WAVE-03 adds DailyLog as the canonical daily aggregate root under the existing Atlas API module. PostgreSQL owns the normalized tables, sqlc owns typed query access, the repository owns transaction-safe aggregate mutations, the service owns validation and optimistic concurrency, and gqlgen exposes granular GraphQL operations under `/graphql/atlas`. Cardio, body weight, charts, AI export, and frontend work remain out of scope.

**Tech Stack:** Go 1.25, PostgreSQL, goose migrations, sqlc v1.30.0, pgx/v5, gqlgen v0.17.49, Nx/Bun command surface.

**Design doc:** `docs/superpowers/specs/2026-06-19-wave-03-workout-diary-design.md`

---

## File Inventory

### Create

| Path | Purpose |
| --- | --- |
| `apps/api/internal/repository/postgres/migrations/00083_daily_logs.sql` | Create versioned DailyLog aggregate table. |
| `apps/api/internal/repository/postgres/migrations/00084_workout_exercises.sql` | Create strength exercise instances inside a DailyLog. |
| `apps/api/internal/repository/postgres/migrations/00085_workout_sets.sql` | Create sets for each workout exercise. |
| `apps/api/internal/repository/postgres/queries/workouts.sql` | sqlc query source for DailyLog, workout exercise, and workout set operations. |
| `apps/api/internal/atlas/models/date.go` | Strict `YYYY-MM-DD` GraphQL Date scalar model. |
| `apps/api/internal/atlas/models/date_test.go` | Date scalar parse/serialize tests. |
| `apps/api/internal/atlas/models/workout.go` | WAVE-03 models, inputs, result envelopes, and error codes. |
| `apps/api/internal/atlas/repository/postgres/workout_repo.go` | sqlc-backed repository and transaction helpers for DailyLog aggregate mutations. |
| `apps/api/internal/repository/postgres/workout_repo_test.go` | PostgreSQL integration tests for WAVE-03 repository behavior. |
| `apps/api/internal/atlas/service/workout.go` | Workout service with validation, snapshots, reindexing, and version checks. |
| `apps/api/internal/atlas/service/workout_service_test.go` | Service unit tests using fakes for WAVE-03 behavior. |
| `apps/api/internal/atlas/graph/schema/workouts.graphql` | Atlas GraphQL schema extensions for DailyLog/workout operations. |
| `apps/api/internal/atlas/graph/resolver/workout.go` | Handwritten resolver methods for DailyLog/workout operations. |
| `apps/api/internal/atlas/graph/resolver/workout_test.go` | Resolver tests for auth/error mapping and service delegation. |

### Modify

| Path | Change |
| --- | --- |
| `apps/api/atlas-gqlgen.yml` | Add WAVE-03 model bindings, including `Date`. |
| `apps/api/internal/atlas/graph/resolver/resolver.go` | Add `WorkoutService service.WorkoutService`. |
| `apps/api/cmd/server/main.go` | Wire workout repository and service into the Atlas resolver. |
| `apps/api/internal/repository/postgres/generated/*.go` | Regenerated sqlc output. |
| `apps/api/internal/atlas/graph/generated/*.go` | Regenerated Atlas gqlgen output. |
| `apps/api/internal/atlas/graph/resolver/*generated*.go` | Regenerated gqlgen resolver stubs as needed. |
| `docs/development-plan.xml` | Add or refresh WAVE-03 API module facts if implementation changes shared contract. |
| `docs/knowledge-graph.xml` | Add WAVE-03 paths and graph annotations if implementation adds public module surfaces. |
| `docs/verification-plan.xml` | Add WAVE-03 verification entries and commands if the implementation introduces new durable checks. |

### Do Not Touch

- Do not implement `cardio_entries`, cardio GraphQL, `CardioType`, or `HeartRateZone`.
- Do not add fake empty cardio fields to `DailyLog`.
- Do not add `body_weight` to `daily_logs`.
- Do not implement frontend routes or UI.
- Do not edit unrelated dirty files such as `opencode.json`.

---

## Task 1: Date Scalar Model

**Files:**
- Create: `apps/api/internal/atlas/models/date.go`
- Create: `apps/api/internal/atlas/models/date_test.go`

- [ ] **Step 1: Write failing Date scalar tests**

Create `apps/api/internal/atlas/models/date_test.go`:

```go
package models

import (
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDate_UnmarshalStrictYYYYMMDD(t *testing.T) {
	var d Date

	require.NoError(t, d.UnmarshalGQL("2026-06-19"))
	assert.Equal(t, "2026-06-19", d.String())
}

func TestDate_RejectsTimestamp(t *testing.T) {
	var d Date

	require.Error(t, d.UnmarshalGQL("2026-06-19T10:00:00Z"))
}

func TestDate_MarshalGQL(t *testing.T) {
	d := MustDate("2026-06-19")

	var out graphql.Marshaler = d.MarshalGQL()
	assert.NotNil(t, out)
}
```

- [ ] **Step 2: Run Date tests and confirm failure**

Run: `cd apps/api && go test ./internal/atlas/models -run TestDate -count=1`

Expected: FAIL because `Date` and `MustDate` are undefined.

- [ ] **Step 3: Implement Date scalar**

Create `apps/api/internal/atlas/models/date.go`:

```go
// FILE: apps/api/internal/atlas/models/date.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define strict calendar-date handling for Atlas GraphQL Date scalar values.
//   SCOPE: Parse, format, marshal, and unmarshal YYYY-MM-DD dates without timezone conversion; excludes timestamp handling.
//   DEPENDS: github.com/99designs/gqlgen/graphql.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Date - Calendar date wrapper for GraphQL and service input.
//   MustDate - Test/helper constructor that panics on invalid date strings.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Date scalar model for WAVE-03.
// END_CHANGE_SUMMARY

package models

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

const dateLayout = "2006-01-02"

type Date struct {
	value time.Time
}

func MustDate(raw string) Date {
	d, err := ParseDate(raw)
	if err != nil {
		panic(err)
	}
	return d
}

func ParseDate(raw string) (Date, error) {
	t, err := time.Parse(dateLayout, raw)
	if err != nil {
		return Date{}, fmt.Errorf("invalid date: %w", err)
	}
	return Date{value: t}, nil
}

func (d Date) String() string {
	if d.value.IsZero() {
		return ""
	}
	return d.value.Format(dateLayout)
}

func (d Date) Time() time.Time {
	return d.value
}

func (d Date) MarshalGQL() graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, `"`+d.String()+`"`)
	})
}

func (d *Date) UnmarshalGQL(v any) error {
	raw, ok := v.(string)
	if !ok {
		return fmt.Errorf("date must be a string")
	}
	parsed, err := ParseDate(raw)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}
```

- [ ] **Step 4: Run Date tests and confirm pass**

Run: `cd apps/api && go test ./internal/atlas/models -run TestDate -count=1`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/atlas/models/date.go apps/api/internal/atlas/models/date_test.go
git commit -m "feat(wave-03): add Atlas Date scalar model"
```

---

## Task 2: Database Migrations

**Files:**
- Create: `apps/api/internal/repository/postgres/migrations/00083_daily_logs.sql`
- Create: `apps/api/internal/repository/postgres/migrations/00084_workout_exercises.sql`
- Create: `apps/api/internal/repository/postgres/migrations/00085_workout_sets.sql`

- [ ] **Step 1: Add DailyLog migration**

Create `apps/api/internal/repository/postgres/migrations/00083_daily_logs.sql`:

```sql
-- FILE: apps/api/internal/repository/postgres/migrations/00083_daily_logs.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add versioned DailyLog aggregate table for WAVE-03 Workout Diary.
--   SCOPE: User-scoped daily container with unique date, nullable notes, version for optimistic concurrency, and timestamps; excludes cardio and body weight.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql.
--   LINKS: M-API / V-M-API / WAVE-03.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   daily_logs - Canonical daily aggregate container shared by strength workouts and later WAVE-04 cardio.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added DailyLog table for WAVE-03.
-- END_CHANGE_SUMMARY

-- +goose Up
CREATE TABLE daily_logs (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES atlas_users(id),
    date       DATE NOT NULL,
    notes      TEXT,
    version    INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_daily_logs_user_date UNIQUE (user_id, date),
    CONSTRAINT chk_daily_logs_version CHECK (version >= 0)
);

CREATE INDEX idx_daily_logs_user_date ON daily_logs (user_id, date);
CREATE INDEX idx_daily_logs_user_date_desc ON daily_logs (user_id, date DESC);

-- +goose Down
DROP TABLE IF EXISTS daily_logs;
```

- [ ] **Step 2: Add WorkoutExercise migration**

Create `apps/api/internal/repository/postgres/migrations/00084_workout_exercises.sql`:

```sql
-- FILE: apps/api/internal/repository/postgres/migrations/00084_workout_exercises.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add workout exercise instances within DailyLog for WAVE-03.
--   SCOPE: Ordered user-scoped exercise instances with working weight snapshots and notes; allows duplicate exercise_id values per day.
--   DEPENDS: daily_logs, exercises, atlas_users.
--   LINKS: M-API / V-M-API / WAVE-03.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   workout_exercises - Ordered strength exercise instances attached to a DailyLog.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added workout_exercises table for WAVE-03.
-- END_CHANGE_SUMMARY

-- +goose Up
CREATE TABLE workout_exercises (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID NOT NULL REFERENCES atlas_users(id),
    daily_log_id            UUID NOT NULL REFERENCES daily_logs(id) ON DELETE CASCADE,
    exercise_id             UUID NOT NULL REFERENCES exercises(id) ON DELETE RESTRICT,
    position                INTEGER NOT NULL,
    working_weight_snapshot REAL,
    notes                   TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_workout_exercises_position CHECK (position > 0),
    CONSTRAINT chk_workout_exercises_working_weight_snapshot CHECK (working_weight_snapshot IS NULL OR working_weight_snapshot > 0),
    CONSTRAINT uq_workout_exercises_daily_log_position UNIQUE (daily_log_id, position)
);

CREATE INDEX idx_workout_exercises_user_daily_log ON workout_exercises (user_id, daily_log_id);
CREATE INDEX idx_workout_exercises_exercise ON workout_exercises (exercise_id);

-- +goose Down
DROP TABLE IF EXISTS workout_exercises;
```

- [ ] **Step 3: Add WorkoutSet migration**

Create `apps/api/internal/repository/postgres/migrations/00085_workout_sets.sql`:

```sql
-- FILE: apps/api/internal/repository/postgres/migrations/00085_workout_sets.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add workout sets for WAVE-03 strength workout logging.
--   SCOPE: Ordered sets with weight, reps, optional RPE/RIR, and notes attached to workout_exercises.
--   DEPENDS: workout_exercises.
--   LINKS: M-API / V-M-API / WAVE-03.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   workout_sets - Ordered strength set rows attached to workout exercise instances.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added workout_sets table for WAVE-03.
-- END_CHANGE_SUMMARY

-- +goose Up
CREATE TABLE workout_sets (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_exercise_id UUID NOT NULL REFERENCES workout_exercises(id) ON DELETE CASCADE,
    set_number          INTEGER NOT NULL,
    weight              REAL NOT NULL,
    reps                INTEGER NOT NULL,
    rpe                 REAL,
    rir                 INTEGER,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_workout_sets_set_number CHECK (set_number > 0),
    CONSTRAINT chk_workout_sets_weight CHECK (weight > 0),
    CONSTRAINT chk_workout_sets_reps CHECK (reps > 0),
    CONSTRAINT chk_workout_sets_rpe CHECK (rpe IS NULL OR (rpe >= 1 AND rpe <= 10)),
    CONSTRAINT chk_workout_sets_rir CHECK (rir IS NULL OR (rir >= 0 AND rir <= 10)),
    CONSTRAINT uq_workout_sets_exercise_set_number UNIQUE (workout_exercise_id, set_number)
);

CREATE INDEX idx_workout_sets_workout_exercise ON workout_sets (workout_exercise_id);

-- +goose Down
DROP TABLE IF EXISTS workout_sets;
```

- [ ] **Step 4: Run migration-backed smoke test**

Run: `cd apps/api && go test ./internal/repository/postgres -run TestNew_ConnectsAndPings -count=1`

Expected: PASS when the integration database is available, or SKIP only when local PostgreSQL test infrastructure is unavailable outside coverage gate.

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/repository/postgres/migrations/00083_daily_logs.sql apps/api/internal/repository/postgres/migrations/00084_workout_exercises.sql apps/api/internal/repository/postgres/migrations/00085_workout_sets.sql
git commit -m "feat(wave-03): add workout diary migrations"
```

---

## Task 3: sqlc Queries And Codegen

**Files:**
- Create: `apps/api/internal/repository/postgres/queries/workouts.sql`
- Modify: `apps/api/internal/repository/postgres/generated/*.go`

- [ ] **Step 1: Create failing repository compile target**

Run: `cd apps/api && go test ./internal/atlas/repository/postgres -run TestNonExistentWorkout -count=1`

Expected: command compiles existing package and reports no tests to run. This establishes the package compiles before sqlc changes.

- [ ] **Step 2: Add workouts.sql query source**

Create `apps/api/internal/repository/postgres/queries/workouts.sql` with these query groups:

```sql
-- FILE: apps/api/internal/repository/postgres/queries/workouts.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc queries for DailyLog, workout_exercises, and workout_sets.
--   SCOPE: Aggregate reads, date range summaries, row locks, version increments, CRUD, and reorder helpers for WAVE-03.
--   DEPENDS: daily_logs, workout_exercises, workout_sets, exercises.
--   LINKS: M-API / V-M-API / WAVE-03.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added WAVE-03 workout diary sqlc queries.
-- END_CHANGE_SUMMARY

-- name: GetDailyLogByDate :one
SELECT id, user_id, date, notes, version, created_at, updated_at
FROM daily_logs
WHERE user_id = $1 AND date = $2
LIMIT 1;

-- name: CreateDailyLog :one
INSERT INTO daily_logs (user_id, date, notes)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, date) DO UPDATE SET updated_at = daily_logs.updated_at
RETURNING id, user_id, date, notes, version, created_at, updated_at;

-- name: LockDailyLogByID :one
SELECT id, user_id, date, notes, version, created_at, updated_at
FROM daily_logs
WHERE user_id = $1 AND id = $2
FOR UPDATE;

-- name: LockDailyLogByDate :one
SELECT id, user_id, date, notes, version, created_at, updated_at
FROM daily_logs
WHERE user_id = $1 AND date = $2
FOR UPDATE;

-- name: IncrementDailyLogVersion :one
UPDATE daily_logs
SET version = version + 1, updated_at = now()
WHERE user_id = $1 AND id = $2
RETURNING id, user_id, date, notes, version, created_at, updated_at;

-- name: UpdateDailyLogNotes :one
UPDATE daily_logs
SET notes = $3, updated_at = now()
WHERE user_id = $1 AND id = $2
RETURNING id, user_id, date, notes, version, created_at, updated_at;

-- name: ListDailyLogSummaries :many
SELECT
  dl.id,
  dl.date,
  dl.version,
  COUNT(DISTINCT we.id)::int AS workout_exercise_count,
  COUNT(ws.id)::int AS workout_set_count,
  COALESCE(SUM(ws.weight * ws.reps), 0)::float8 AS total_volume,
  dl.updated_at
FROM daily_logs dl
LEFT JOIN workout_exercises we ON we.daily_log_id = dl.id
LEFT JOIN workout_sets ws ON ws.workout_exercise_id = we.id
WHERE dl.user_id = $1 AND dl.date >= $2 AND dl.date <= $3
GROUP BY dl.id
ORDER BY dl.date ASC;
```

Then add CRUD/select queries for:

- list workout exercises by `daily_log_id`, ordered by `position`
- create/update/delete workout exercise
- shift/rewrite workout exercise positions
- list workout sets by workout exercise ids, ordered by `set_number`
- create/update/delete workout set
- shift/rewrite set numbers

Use `RETURNING` on mutation queries so repository tests can assert saved values immediately.

- [ ] **Step 3: Run sqlc codegen**

Run: `cd apps/api && go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 generate`

Expected: generated code includes workout query methods and compiles.

- [ ] **Step 4: Run API package compile check**

Run: `cd apps/api && go test ./internal/repository/postgres -run TestNew_ReturnsErrorForBadPort -count=1`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/repository/postgres/queries/workouts.sql apps/api/internal/repository/postgres/generated
git commit -m "feat(wave-03): add workout sqlc queries"
```

---

## Task 4: Repository Integration Tests And Implementation

**Files:**
- Create: `apps/api/internal/atlas/repository/postgres/workout_repo.go`
- Create: `apps/api/internal/repository/postgres/workout_repo_test.go`

- [ ] **Step 1: Write failing repository tests**

Create `apps/api/internal/repository/postgres/workout_repo_test.go` with tests for:

```go
func TestWorkoutRepo_GetOrCreateDailyLog_UniquePerUserDate(t *testing.T)
func TestWorkoutRepo_DailyLog_UserScopedIsolation(t *testing.T)
func TestWorkoutRepo_AddWorkoutExercise_AllowsDuplicateExercise(t *testing.T)
func TestWorkoutRepo_AddWorkoutExercise_CapturesWorkingWeightSnapshot(t *testing.T)
func TestWorkoutRepo_ReorderWorkoutExercises_ReindexesContiguously(t *testing.T)
func TestWorkoutRepo_DeleteWorkoutExercise_CascadesSetsAndKeepsDailyLog(t *testing.T)
func TestWorkoutRepo_AddWorkoutSet_ValidatesDBConstraints(t *testing.T)
func TestWorkoutRepo_ReorderWorkoutSets_ReindexesContiguously(t *testing.T)
func TestWorkoutRepo_IncrementDailyLogVersion(t *testing.T)
```

Use the existing test setup pattern from `apps/api/internal/repository/postgres/exercise_repo_test.go`: `testinfra.PostgresConfig`, `postgresrepo.New`, `postgresrepo.RunMigrations`, and explicit Atlas user setup.

- [ ] **Step 2: Run repository tests and confirm failure**

Run: `cd apps/api && go test ./internal/repository/postgres -run 'TestWorkoutRepo|TestDailyLog' -count=1`

Expected: FAIL because `NewWorkoutRepository` and repository methods are undefined.

- [ ] **Step 3: Implement repository**

Create `apps/api/internal/atlas/repository/postgres/workout_repo.go`:

```go
// FILE: apps/api/internal/atlas/repository/postgres/workout_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement WAVE-03 DailyLog aggregate repository operations with sqlc and pgx transactions.
//   SCOPE: DailyLog get/create/lock/version, workout exercise CRUD/reorder, workout set CRUD/reorder, and aggregate reads; excludes service validation and GraphQL mapping.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/atlas/models, pgx/v5.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   WorkoutRepository - Data access interface for WAVE-03 DailyLog aggregate operations.
//   NewWorkoutRepository - Creates the sqlc-backed repository.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added workout repository for WAVE-03.
// END_CHANGE_SUMMARY
```

Interface responsibilities:

- `GetDailyLogByDate`
- `GetOrCreateDailyLogByDate`
- `GetDailyLogAggregate`
- `WithLockedDailyLogByDate`
- `WithLockedDailyLogByWorkoutExerciseID`
- `WithLockedDailyLogByWorkoutSetID`
- `IncrementDailyLogVersion`
- workout exercise CRUD/reorder helpers
- workout set CRUD/reorder helpers

Implement transactions using `pgxpool.Pool.BeginTx` and `generated.New(tx)`. Keep all position rewrites and version increments inside the same transaction.

- [ ] **Step 4: Run repository tests and confirm pass**

Run: `cd apps/api && go test ./internal/repository/postgres -run 'TestWorkoutRepo|TestDailyLog' -count=1`

Expected: PASS or SKIP only if local integration database is unavailable outside coverage gate.

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/atlas/repository/postgres/workout_repo.go apps/api/internal/repository/postgres/workout_repo_test.go
git commit -m "feat(wave-03): add workout repository"
```

---

## Task 5: Service Tests And Implementation

**Files:**
- Create: `apps/api/internal/atlas/models/workout.go`
- Create: `apps/api/internal/atlas/service/workout.go`
- Create: `apps/api/internal/atlas/service/workout_service_test.go`

- [ ] **Step 1: Write service tests first**

Create `apps/api/internal/atlas/service/workout_service_test.go` using fake repositories. Cover:

```go
func TestWorkoutService_UpdateNotes_CreatesDailyLogAtExpectedVersionZero(t *testing.T)
func TestWorkoutService_UpdateNotes_RejectsStaleVersion(t *testing.T)
func TestWorkoutService_AddExercise_RequiresExistingExercise(t *testing.T)
func TestWorkoutService_AddExercise_CapturesWorkingWeightSnapshot(t *testing.T)
func TestWorkoutService_AddExercise_AllowsDuplicateExerciseID(t *testing.T)
func TestWorkoutService_RemoveExercise_KeepsEmptyDailyLog(t *testing.T)
func TestWorkoutService_AddSet_ValidatesWeightRepsRpeRir(t *testing.T)
func TestWorkoutService_UpdateSet_ReindexesWhenSetNumberChanges(t *testing.T)
func TestWorkoutService_ReorderExercises_RejectsMissingDuplicateOrForeignIDs(t *testing.T)
func TestWorkoutService_ReorderSets_RejectsMissingDuplicateOrForeignIDs(t *testing.T)
```

- [ ] **Step 2: Run service tests and confirm failure**

Run: `cd apps/api && go test ./internal/atlas/service -run 'TestWorkoutService' -count=1`

Expected: FAIL because workout models/service are undefined.

- [ ] **Step 3: Add workout models**

Create `apps/api/internal/atlas/models/workout.go` with file-local GRACE markup and these exported types:

- `DailyLogRecord`
- `DailyLog`
- `DailyLogSummary`
- `WorkoutExerciseRecord`
- `WorkoutExercise`
- `WorkoutSetRecord`
- `WorkoutSet`
- `AddWorkoutExerciseInput`
- `UpdateWorkoutExerciseInput`
- `AddWorkoutSetInput`
- `UpdateWorkoutSetInput`
- `DailyLogResult`
- `DailyLogValidationErr`
- `DailyLogNotFoundErr`
- `DailyLogConflictErr`
- `DailyLogAuthErr`
- `DailyLogErrorCode`

Use pointer fields for nullable GraphQL and database values. Keep `DailyLogConflictErr.CurrentVersion` non-null and include optional `CurrentDailyLog`.

- [ ] **Step 4: Implement workout service**

Create `apps/api/internal/atlas/service/workout.go` with file-local GRACE markup and:

```go
type WorkoutService interface {
	GetDailyLog(ctx context.Context, userID string, date models.Date) (*models.DailyLog, error)
	ListDailyLogSummaries(ctx context.Context, userID string, from models.Date, to models.Date) ([]models.DailyLogSummary, error)
	UpdateDailyLogNotes(ctx context.Context, userID string, date models.Date, expectedVersion int32, notes *string) (*models.DailyLog, error)
	AddWorkoutExercise(ctx context.Context, userID string, date models.Date, expectedVersion int32, input models.AddWorkoutExerciseInput) (*models.DailyLog, error)
	UpdateWorkoutExercise(ctx context.Context, userID string, id string, expectedVersion int32, input models.UpdateWorkoutExerciseInput) (*models.DailyLog, error)
	RemoveWorkoutExercise(ctx context.Context, userID string, id string, expectedVersion int32) (*models.DailyLog, error)
	ReorderWorkoutExercises(ctx context.Context, userID string, date models.Date, expectedVersion int32, orderedIDs []string) (*models.DailyLog, error)
	AddWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, expectedVersion int32, input models.AddWorkoutSetInput) (*models.DailyLog, error)
	UpdateWorkoutSet(ctx context.Context, userID string, id string, expectedVersion int32, input models.UpdateWorkoutSetInput) (*models.DailyLog, error)
	RemoveWorkoutSet(ctx context.Context, userID string, id string, expectedVersion int32) (*models.DailyLog, error)
	ReorderWorkoutSets(ctx context.Context, userID string, workoutExerciseID string, expectedVersion int32, orderedIDs []string) (*models.DailyLog, error)
}
```

Service rules:

- validate `expectedVersion >= 0`
- reject date ranges where `from > to`
- validate set values before repository calls
- use ExerciseService or ExerciseRepository to fetch current exercise and snapshot `WorkingWeight`
- compare locked DailyLog version before every mutation
- return conflict error with current aggregate when versions differ
- keep empty DailyLog rows

- [ ] **Step 5: Run service tests and confirm pass**

Run: `cd apps/api && go test ./internal/atlas/service -run 'TestWorkoutService' -count=1`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add apps/api/internal/atlas/models/workout.go apps/api/internal/atlas/service/workout.go apps/api/internal/atlas/service/workout_service_test.go
git commit -m "feat(wave-03): add workout service"
```

---

## Task 6: GraphQL Schema And Atlas Codegen

**Files:**
- Create: `apps/api/internal/atlas/graph/schema/workouts.graphql`
- Modify: `apps/api/atlas-gqlgen.yml`
- Modify: `apps/api/internal/atlas/graph/generated/*.go`
- Modify: generated resolver stubs under `apps/api/internal/atlas/graph/resolver/`

- [ ] **Step 1: Add workouts GraphQL schema**

Create `apps/api/internal/atlas/graph/schema/workouts.graphql`:

```graphql
scalar Date

extend type Query {
  dailyLog(date: Date!): DailyLogResult!
  dailyLogs(from: Date!, to: Date!): [DailyLogSummary!]!
}

extend type Mutation {
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

Add the types, inputs, and errors exactly from the approved design doc. Do not add cardio fields.

- [ ] **Step 2: Add gqlgen bindings**

Modify `apps/api/atlas-gqlgen.yml`:

```yaml
  Date:
    model: monorepo-template/apps/api/internal/atlas/models.Date
  DailyLog:
    model: monorepo-template/apps/api/internal/atlas/models.DailyLog
  DailyLogSummary:
    model: monorepo-template/apps/api/internal/atlas/models.DailyLogSummary
  WorkoutExercise:
    model: monorepo-template/apps/api/internal/atlas/models.WorkoutExercise
  WorkoutSet:
    model: monorepo-template/apps/api/internal/atlas/models.WorkoutSet
  AddWorkoutExerciseInput:
    model: monorepo-template/apps/api/internal/atlas/models.AddWorkoutExerciseInput
  UpdateWorkoutExerciseInput:
    model: monorepo-template/apps/api/internal/atlas/models.UpdateWorkoutExerciseInput
  AddWorkoutSetInput:
    model: monorepo-template/apps/api/internal/atlas/models.AddWorkoutSetInput
  UpdateWorkoutSetInput:
    model: monorepo-template/apps/api/internal/atlas/models.UpdateWorkoutSetInput
  DailyLogResult:
    model: monorepo-template/apps/api/internal/atlas/models.DailyLogResult
```

Bind all WAVE-03 error types and enum types too.

- [ ] **Step 3: Run Atlas gqlgen and confirm generated stubs**

Run: `bunx nx run api:codegen:atlas`

Expected: generated Atlas GraphQL code updates without schema errors and creates resolver stubs for WAVE-03 operations.

- [ ] **Step 4: Run compile check**

Run: `cd apps/api && go test ./internal/atlas/graph/generated -run TestNonExistent -count=1`

Expected: package compiles and reports no tests to run.

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/atlas/graph/schema/workouts.graphql apps/api/atlas-gqlgen.yml apps/api/internal/atlas/graph/generated apps/api/internal/atlas/graph/resolver
git commit -m "feat(wave-03): add workout GraphQL schema"
```

---

## Task 7: Resolver Tests And Implementation

**Files:**
- Create: `apps/api/internal/atlas/graph/resolver/workout.go`
- Create: `apps/api/internal/atlas/graph/resolver/workout_test.go`
- Modify: `apps/api/internal/atlas/graph/resolver/resolver.go`
- Modify: generated resolver stubs as needed by gqlgen preservation.

- [ ] **Step 1: Write resolver tests**

Create `apps/api/internal/atlas/graph/resolver/workout_test.go` with fake `WorkoutService`. Cover:

```go
func TestDailyLogResolver_UnauthorizedReturnsAuthError(t *testing.T)
func TestDailyLogResolver_DelegatesAuthenticatedDailyLog(t *testing.T)
func TestUpdateDailyLogNotesResolver_MapsConflictError(t *testing.T)
func TestAddWorkoutExerciseResolver_MapsValidationError(t *testing.T)
func TestWorkoutSetResolvers_MapNotFoundError(t *testing.T)
```

Use existing auth context helper patterns from WAVE-02 resolver tests.

- [ ] **Step 2: Run resolver tests and confirm failure**

Run: `cd apps/api && go test ./internal/atlas/graph/resolver -run 'Test.*DailyLog|Test.*Workout' -count=1`

Expected: FAIL because resolver methods and `WorkoutService` field are missing.

- [ ] **Step 3: Add WorkoutService to root resolver**

Modify `apps/api/internal/atlas/graph/resolver/resolver.go`:

```go
type Resolver struct {
	SettingsService service.SettingsService
	PinService      service.PinService
	ExerciseService service.ExerciseService
	WorkoutService  service.WorkoutService
}
```

- [ ] **Step 4: Implement handwritten resolver methods**

Create `apps/api/internal/atlas/graph/resolver/workout.go` with file-local GRACE markup. Implement thin methods:

- read user ID via `middleware.GetAtlasUserID(ctx)`
- return `DailyLogAuthErr` when missing
- call `WorkoutService`
- map known service errors to `DailyLogResult`
- return nil error for typed domain errors
- avoid logging notes or sensitive payloads

- [ ] **Step 5: Wire generated stubs to handwritten methods**

In generated-preserved resolver stub files, forward gqlgen methods to `r.Resolver.<MethodName>`. Follow the WAVE-02 pattern in `apps/api/internal/atlas/graph/resolver/exercise.go` and generated `schema.resolvers.go`.

- [ ] **Step 6: Run resolver tests and confirm pass**

Run: `cd apps/api && go test ./internal/atlas/graph/resolver -run 'Test.*DailyLog|Test.*Workout' -count=1`

Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add apps/api/internal/atlas/graph/resolver/resolver.go apps/api/internal/atlas/graph/resolver/workout.go apps/api/internal/atlas/graph/resolver/workout_test.go apps/api/internal/atlas/graph/resolver
git commit -m "feat(wave-03): add workout GraphQL resolvers"
```

---

## Task 8: Runtime Wiring

**Files:**
- Modify: `apps/api/cmd/server/main.go`

- [ ] **Step 1: Add failing startup compile check**

Run: `cd apps/api && go test ./cmd/server -run TestNonExistent -count=1`

Expected: package compile may fail until wiring is updated if resolver now requires `WorkoutService`.

- [ ] **Step 2: Wire workout service into Atlas resolver**

Modify `apps/api/cmd/server/main.go` near existing Atlas repository/service wiring:

```go
atlasWorkoutRepo := atlasPostgres.NewWorkoutRepository(db.Pool)
atlasWorkoutService := atlasService.NewWorkoutService(atlasWorkoutRepo, atlasExerciseService)
```

Then add it to the resolver:

```go
atlasRes := &atlasResolver.Resolver{
	SettingsService: atlasSettingsService,
	PinService:      atlasPinService,
	ExerciseService: atlasExerciseService,
	WorkoutService:  atlasWorkoutService,
}
```

Update file-local `START_CHANGE_SUMMARY` in `main.go`.

- [ ] **Step 3: Run startup compile check**

Run: `cd apps/api && go test ./cmd/server -run TestNonExistent -count=1`

Expected: package compiles and reports no tests to run.

- [ ] **Step 4: Commit**

```bash
git add apps/api/cmd/server/main.go
git commit -m "feat(wave-03): wire workout service"
```

---

## Task 9: Focused Integration Verification

**Files:**
- Modify only files needed to fix failures from focused checks.

- [ ] **Step 1: Run sqlc and Atlas gqlgen together**

Run: `bunx nx run api:codegen && bunx nx run api:codegen:atlas`

Expected: both commands exit 0 and generated files are current.

- [ ] **Step 2: Run WAVE-03 focused Go tests**

Run:

```bash
cd apps/api && go test ./internal/atlas/models -run TestDate -count=1
cd apps/api && go test ./internal/atlas/service -run 'TestWorkoutService' -count=1
cd apps/api && go test ./internal/atlas/graph/resolver -run 'Test.*DailyLog|Test.*Workout' -count=1
cd apps/api && go test ./internal/repository/postgres -run 'TestWorkoutRepo|TestDailyLog' -count=1
```

Expected: PASS, except repository integration tests may SKIP only when local test PostgreSQL is unavailable outside coverage gate.

- [ ] **Step 3: Run affected API test target**

Run: `bunx nx test api`

Expected: PASS. If repository tests skip due unavailable local database, record the exact skip reason and run again with test Docker infrastructure before handoff.

- [ ] **Step 4: Run API build**

Run: `bunx nx build api`

Expected: PASS.

- [ ] **Step 5: Commit verification fixes if any**

```bash
git add <changed-files>
git commit -m "fix(wave-03): close workout verification gaps"
```

Skip commit if no files changed.

---

## Task 10: GRACE Docs And Handoff Evidence

**Files:**
- Modify: `docs/development-plan.xml`
- Modify: `docs/knowledge-graph.xml`
- Modify: `docs/verification-plan.xml`
- Optional create: `.tasks/WAVE-03/HANDOFF.md`

- [ ] **Step 1: Update shared GRACE artifacts if implementation changed public contracts**

Update the API module descriptions with WAVE-03 paths and verification facts:

- `docs/development-plan.xml`: mention DailyLog/workout service, repository, schema, and aggregate versioning under `M-API` or a new Atlas submodule if project convention prefers that.
- `docs/knowledge-graph.xml`: add WAVE-03 source/test paths and annotations.
- `docs/verification-plan.xml`: add WAVE-03 focused commands and scenarios.

- [ ] **Step 2: Write handoff evidence**

Create `.tasks/WAVE-03/HANDOFF.md` only if the session is closing the implementation wave or preparing MR/review. Include:

- scope implemented
- commands run
- pass/fail/skip evidence
- generated artifact status
- WAVE-04 compatibility note
- known risks

- [ ] **Step 3: Validate XML artifacts**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
```

Expected: no output and exit 0.

- [ ] **Step 4: Run GRACE lint**

Run: `grace lint --path .`

Expected: PASS or only known unrelated baseline issues recorded with exact output.

- [ ] **Step 5: Commit docs and handoff artifacts**

```bash
git add docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml .tasks/WAVE-03/HANDOFF.md
git commit -m "docs(wave-03): update GRACE workout handoff"
```

Skip unchanged or nonexistent paths.

---

## Task 11: Final WAVE-03 Closeout

**Files:**
- No planned code edits. Only verification-driven fixes if a command fails.

- [ ] **Step 1: Confirm no forbidden scope was added**

Run:

```bash
rg -n "cardio_entries|CardioType|HeartRateZone|body_weight|bodyWeight" apps/api/internal/atlas apps/api/internal/repository/postgres/migrations apps/api/internal/repository/postgres/queries
```

Expected: no WAVE-03-introduced cardio implementation or body weight persistence. Existing docs or future-wave references are acceptable only outside product code.

- [ ] **Step 2: Run final focused gates**

Run:

```bash
bunx nx run api:codegen
bunx nx run api:codegen:atlas
bunx nx test api
bunx nx build api
```

Expected: all exit 0.

- [ ] **Step 3: Confirm git status and commit history**

Run:

```bash
git status --short --branch
git log --oneline -5
```

Expected: branch contains WAVE-03 commits; only explicitly unrelated local files may remain dirty.

- [ ] **Step 4: Push**

Run:

```bash
git pull --rebase
git push
git status --short --branch
```

Expected: push succeeds and branch is up to date with origin. Any unrelated pre-existing dirty file must remain unstaged unless the user explicitly includes it.

