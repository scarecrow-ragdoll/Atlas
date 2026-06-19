# sqlc pgx goose users reference Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `sqlc + pgx/v5 + goose` the backend SQL standard by migrating the existing API `users` repository slice to generated sqlc queries while preserving REST, GraphQL, service, and test behavior.

**Architecture:** `goose` remains the PostgreSQL schema source and runtime migration tool. `sqlc` reads the existing goose migration directory plus `queries/users.sql`, emits an internal generated package using `pgx/v5`, and `postgres.UserRepo` stays the public adapter consumed by `UserService`. Root/API codegen, generated-drift checks, coverage allowlists, GRACE docs, and file-local contracts are updated together.

**Tech Stack:** Go 1.25, sqlc v1.30.0, pgx/v5, goose, Bun, Nx, Vitest, Go test, PostgreSQL test compose, GRACE XML.

---

<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the task-by-task implementation plan for migrating the users repository slice to sqlc plus pgx/v5 plus goose. -->
<!--   SCOPE: Planning only; includes file ownership, TDD sequence, commands, docs, verification, and commit boundaries; excludes implementation execution. -->
<!--   DEPENDS: docs/superpowers/specs/2026-06-05-sqlc-pgx-goose-users-reference-design.md, apps/api users repository, tools/coverage, docs/*.xml. -->
<!--   LINKS: M-API / M-WORKSPACE / M-COVERAGE-GATE / M-GRACE-WORKFLOW / V-M-API / V-M-WORKSPACE / V-M-COVERAGE-GATE. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   File Structure - Lists all planned creates and modifications with ownership boundaries. -->
<!--   Task 1 - Adds sqlc configuration, users query SQL, tool dependency, and generated code. -->
<!--   Task 2 - Migrates UserRepo behind adapter tests while preserving public service contract. -->
<!--   Task 3 - Expands repository integration coverage and safe destructive-test evidence. -->
<!--   Task 4 - Adds codegen drift detection through the workspace validation target. -->
<!--   Task 5 - Updates generated-code coverage policy for sqlc output. -->
<!--   Task 6 - Synchronizes GRACE shared docs and file-local contracts. -->
<!--   Task 7 - Runs focused and final verification gates. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Initial ready-for-execution plan. -->
<!-- END_CHANGE_SUMMARY -->

## File Structure

Create:

- `apps/api/sqlc.yaml` - API-owned sqlc v2 configuration. Reads goose migrations plus users query SQL and emits `pgx/v5` generated code.
- `apps/api/internal/repository/postgres/queries/users.sql` - Named sqlc queries for existing users repository behavior.
- `apps/api/internal/repository/postgres/generated/db.go` - Generated sqlc DBTX and `Queries` wrapper.
- `apps/api/internal/repository/postgres/generated/models.go` - Generated sqlc table models.
- `apps/api/internal/repository/postgres/generated/querier.go` - Generated sqlc `Querier` interface.
- `apps/api/internal/repository/postgres/generated/users.sql.go` - Generated sqlc users query methods.

Modify:

- `apps/api/tools.go` - Track sqlc command package with the existing Go tool dependency pattern.
- `apps/api/go.mod` and `apps/api/go.sum` - Pin `github.com/sqlc-dev/sqlc` at `v1.30.0`.
- `apps/api/project.json` - Run sqlc generation before gqlgen in `api:codegen`.
- `apps/api/internal/repository/postgres/user_repo.go` - Replace inline SQL with a thin adapter around `generated.Querier`.
- `apps/api/internal/repository/postgres/user_repo_unit_test.go` - Replace pool-level fakes with generated-query adapter fakes and edge-case tests.
- `apps/api/internal/repository/postgres/user_repo_test.go` - Add sqlc/goose integration coverage for nullable update paths and real DB evidence.
- `tools/codegen/project.json` - Make `codegen:validate` run root codegen and fail on generated-file drift.
- `tools/coverage/coverage.config.json` - Add exact sqlc generated Go files to the allowlist with replacement gates.
- `docs/requirements.xml` - Include backend SQL codegen and committed generated-code drift policy.
- `docs/technology.xml` - Add sqlc tool/dependency and preferred API SQL codegen.
- `docs/development-plan.xml` - Update `M-WORKSPACE`, `M-COVERAGE-GATE`, and `M-API`.
- `docs/knowledge-graph.xml` - Add sqlc config/query/generated paths and codegen/coverage ownership links.
- `docs/verification-plan.xml` - Add sqlc codegen, drift, integration, and coverage policy checks.

Do not modify:

- `apps/api/internal/service/user_service.go` - `UserRepository` remains the service boundary.
- REST handlers and GraphQL resolvers - They keep calling `UserService`.
- `apps/api/internal/repository/postgres/migrations/00001_init.sql` - No schema change is intended.
- `tools/coverage/run.mjs` - Avoid changing the runner by allowlisting exact generated `.go` files.

Commit policy:

- Do not commit in Tasks 1-6. They intentionally move code, generated artifacts, tests, tooling, and GRACE docs toward one synchronized state.
- Commit only after Task 7 verification proves docs, code, generated output, coverage policy, and evidence are aligned.
- If a later executor splits commits, every commit must include the relevant GRACE docs and verification evidence needed for that commit to be internally consistent.

## Task 1: Add sqlc Generation Surface

**Files:**

- Create: `apps/api/sqlc.yaml`
- Create: `apps/api/internal/repository/postgres/queries/users.sql`
- Create by generation: `apps/api/internal/repository/postgres/generated/db.go`
- Create by generation: `apps/api/internal/repository/postgres/generated/models.go`
- Create by generation: `apps/api/internal/repository/postgres/generated/querier.go`
- Create by generation: `apps/api/internal/repository/postgres/generated/users.sql.go`
- Modify: `apps/api/tools.go`
- Modify: `apps/api/go.mod`
- Modify: `apps/api/go.sum`
- Modify: `apps/api/project.json`

- [ ] **Step 1: Confirm sqlc v1.31.1 is not compatible with the repo Go contract**

Run:

```bash
cd apps/api
go list -m -versions github.com/sqlc-dev/sqlc
tmp="$(mktemp -d)"
cd "$tmp"
go mod init sqlc-version-check
GOTOOLCHAIN=local go get github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1
```

Expected:

```text
github.com/sqlc-dev/sqlc ... v1.30.0 v1.31.0 v1.31.1
go: github.com/sqlc-dev/sqlc@v1.31.1 requires go >= 1.26.0
```

The `GOTOOLCHAIN=local go get` command is expected to exit non-zero. This proves the implementation must pin `v1.30.0` and must not upgrade the repository Go version for this rollout.

- [ ] **Step 2: Add sqlc to API Go tool dependencies**

Run:

```bash
cd apps/api
go get github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0
```

Modify `apps/api/tools.go` to add sqlc beside the existing tool imports:

```go
//go:build tools

// FILE: apps/api/tools.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Track Go command dependencies used by API development and code generation.
//   SCOPE: Tool-only module imports for local and CI reproducibility; excludes runtime API dependencies.
//   DEPENDS: github.com/99designs/gqlgen, github.com/pressly/goose/v3/cmd/goose, github.com/sqlc-dev/sqlc/cmd/sqlc.
//   LINKS: M-API / M-WORKSPACE / V-M-API / V-M-WORKSPACE.
//   ROLE: CONFIG
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   imports - Pins gqlgen, goose, and sqlc command packages through Go module resolution.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added sqlc command dependency for API SQL codegen.
// END_CHANGE_SUMMARY

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/pressly/goose/v3/cmd/goose"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)
```

Run:

```bash
cd apps/api
go mod tidy
go list -deps ./... >/dev/null
```

Expected:

```text
# go mod tidy prints no errors
# go list prints no output and exits 0
```

Do not run `go test -tags=tools`; the repository already uses command-package tool imports under the `tools` tag, and normal module commands must remain green.

- [ ] **Step 3: Create sqlc configuration**

Create `apps/api/sqlc.yaml`:

```yaml
# FILE: apps/api/sqlc.yaml
# VERSION: 1.0.0
# START_MODULE_CONTRACT
#   PURPOSE: Configure sqlc generation for API PostgreSQL query code.
#   SCOPE: Reads goose migrations plus repository query SQL and emits internal pgx/v5 generated code; excludes runtime migration execution.
#   DEPENDS: internal/repository/postgres/migrations, internal/repository/postgres/queries, github.com/sqlc-dev/sqlc, github.com/jackc/pgx/v5.
#   LINKS: M-API / M-WORKSPACE / V-M-API / V-M-WORKSPACE.
#   ROLE: CONFIG
#   MAP_MODE: SUMMARY
# END_MODULE_CONTRACT
# START_MODULE_MAP
#   sql - Defines the users query generation unit for PostgreSQL.
# END_MODULE_MAP
# START_CHANGE_SUMMARY
#   LAST_CHANGE: 1.0.0 - Added users sqlc generation config.
# END_CHANGE_SUMMARY
version: '2'
sql:
  - engine: 'postgresql'
    schema: 'internal/repository/postgres/migrations'
    queries: 'internal/repository/postgres/queries'
    gen:
      go:
        package: 'generated'
        out: 'internal/repository/postgres/generated'
        sql_package: 'pgx/v5'
        emit_interface: true
```

- [ ] **Step 4: Create users query SQL**

Create `apps/api/internal/repository/postgres/queries/users.sql`:

```sql
-- FILE: apps/api/internal/repository/postgres/queries/users.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc users queries used by the PostgreSQL UserRepo adapter.
--   SCOPE: CRUD and pagination queries for the users table; excludes schema ownership and transport mapping.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00001_init.sql.
--   LINKS: M-API / V-M-API.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   GetUserByID - Fetches one user row by UUID.
--   ListUsers - Fetches a created_at-desc page with optional cursor cutoff.
--   CountUsers - Counts all user rows for connection metadata.
--   CreateUser - Inserts one user and returns the persisted public row.
--   UpdateUser - Applies nullable name/email updates and returns the persisted public row.
--   DeleteUser - Deletes one user idempotently by UUID.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added users sqlc queries.
-- END_CHANGE_SUMMARY

-- name: GetUserByID :one
SELECT id, email, name, created_at, updated_at
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, email, name, created_at, updated_at
FROM users
WHERE (sqlc.narg('after_created_at')::timestamptz IS NULL OR created_at < sqlc.narg('after_created_at')::timestamptz)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit_rows');

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CreateUser :one
INSERT INTO users (email, name, password_hash)
VALUES ($1, $2, $3)
RETURNING id, email, name, created_at, updated_at;

-- name: UpdateUser :one
UPDATE users
SET name = COALESCE(sqlc.narg('name'), name),
    email = COALESCE(sqlc.narg('email'), email),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING id, email, name, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
```

- [ ] **Step 5: Run sqlc generation directly**

Run:

```bash
cd apps/api
GOTOOLCHAIN=local go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 generate
```

Expected:

```text
# no output
```

Confirm generated files:

```bash
ls internal/repository/postgres/generated
```

Expected:

```text
db.go
models.go
querier.go
users.sql.go
```

Do not manually edit generated files to add GRACE markup. Their generated contract is captured by `apps/api/sqlc.yaml`, `queries/users.sql`, shared GRACE docs, coverage allowlist entries, and the generated-drift gate.

- [ ] **Step 6: Wire api:codegen to run sqlc before gqlgen**

Modify the `codegen` target in `apps/api/project.json`:

```json
"codegen": {
  "executor": "nx:run-commands",
  "options": {
    "command": "cd apps/api && go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 generate && go run github.com/99designs/gqlgen generate"
  },
  "dependsOn": ["^validate"]
}
```

Run:

```bash
bunx nx run api:codegen
```

Expected:

```text
> nx run api:codegen
# command completes successfully
```

- [ ] **Step 7: Checkpoint generation surface without committing**

Run:

```bash
git status --short -- apps/api/sqlc.yaml \
  apps/api/internal/repository/postgres/queries/users.sql \
  apps/api/internal/repository/postgres/generated \
  apps/api/tools.go apps/api/go.mod apps/api/go.sum apps/api/project.json
```

Expected:

```text
# status shows only the intended generation-surface files from this task
```

Do not commit yet. GRACE docs and verification evidence are not synchronized until Tasks 6-7.

## Task 2: Migrate UserRepo Through a Generated-Query Adapter

**Files:**

- Modify: `apps/api/internal/repository/postgres/user_repo_unit_test.go`
- Modify: `apps/api/internal/repository/postgres/user_repo.go`

- [ ] **Step 1: Replace unit-test fakes with generated-query adapter fakes**

Replace the pool/row fake types in `apps/api/internal/repository/postgres/user_repo_unit_test.go` with this fake query surface. Remove the old pool-backed `TestUserRepo_*` unit tests that construct `&UserRepo{pool: ...}`; they are replaced by generated-query adapter tests in the next step. Keep the existing `TestDuplicateHelpers`, `TestRunMigrations_ReturnsOpenError`, and `TestNew_ReturnsPoolConstructionError` tests after these helpers.

```go
// FILE: apps/api/internal/repository/postgres/user_repo_unit_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify PostgreSQL users repository adapter behavior at the generated-query boundary.
//   SCOPE: Adapter mapping, nullable params, error translation, migration construction failures, and helper behavior; excludes real database integration.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/service, apps/api/internal/testinfra.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   fakeUserQueries - Test double for the sqlc generated Querier interface.
//   TestUserRepo_* - Adapter behavior and error mapping coverage.
//   TestRunMigrations_ReturnsOpenError - Migration driver failure coverage.
//   TestNew_ReturnsPoolConstructionError - pgx pool construction failure coverage.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Replaced pool-level fakes with sqlc generated-query adapter tests.
// END_CHANGE_SUMMARY

package postgres

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/repository/postgres/generated"
	"monorepo-template/apps/api/internal/service"
	"monorepo-template/apps/api/internal/testinfra"
)

type fakeUserQueries struct {
	getRow    generated.GetUserByIDRow
	getErr    error
	listRows  []generated.ListUsersRow
	listErr   error
	count     int64
	countErr  error
	createRow generated.CreateUserRow
	createErr error
	updateRow generated.UpdateUserRow
	updateErr error
	deleteErr error

	lastListParams   generated.ListUsersParams
	lastUpdateParams generated.UpdateUserParams
}

func (f *fakeUserQueries) GetUserByID(ctx context.Context, id pgtype.UUID) (generated.GetUserByIDRow, error) {
	return f.getRow, f.getErr
}

func (f *fakeUserQueries) ListUsers(ctx context.Context, arg generated.ListUsersParams) ([]generated.ListUsersRow, error) {
	f.lastListParams = arg
	return f.listRows, f.listErr
}

func (f *fakeUserQueries) CountUsers(ctx context.Context) (int64, error) {
	return f.count, f.countErr
}

func (f *fakeUserQueries) CreateUser(ctx context.Context, arg generated.CreateUserParams) (generated.CreateUserRow, error) {
	return f.createRow, f.createErr
}

func (f *fakeUserQueries) UpdateUser(ctx context.Context, arg generated.UpdateUserParams) (generated.UpdateUserRow, error) {
	f.lastUpdateParams = arg
	return f.updateRow, f.updateErr
}

func (f *fakeUserQueries) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	return f.deleteErr
}

func fakeUUID(t *testing.T, value string) pgtype.UUID {
	t.Helper()
	uuid, err := uuidFromString(value)
	require.NoError(t, err)
	return uuid
}

func fakeTime(value string) pgtype.Timestamptz {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		panic(err)
	}
	return pgtype.Timestamptz{Time: parsed, Valid: true}
}
```

- [ ] **Step 2: Add adapter unit tests that fail before implementation**

Add these tests to `apps/api/internal/repository/postgres/user_repo_unit_test.go` before the migration and DB construction tests:

```go
func TestUserRepo_ListMapsGeneratedRowsAndTrimsExtraRow(t *testing.T) {
	queries := &fakeUserQueries{
		listRows: []generated.ListUsersRow{
			{ID: fakeUUID(t, "00000000-0000-0000-0000-000000000001"), Email: "one@example.com", Name: "One", CreatedAt: fakeTime("2026-06-05T10:00:00Z"), UpdatedAt: fakeTime("2026-06-05T10:01:00Z")},
			{ID: fakeUUID(t, "00000000-0000-0000-0000-000000000002"), Email: "two@example.com", Name: "Two", CreatedAt: fakeTime("2026-06-05T09:00:00Z"), UpdatedAt: fakeTime("2026-06-05T09:01:00Z")},
		},
		count: 2,
	}
	repo := newUserRepoWithQueries(queries)
	first := 1

	users, total, err := repo.List(context.Background(), &first, nil)

	require.NoError(t, err)
	assert.Equal(t, 2, total)
	require.Len(t, users, 1)
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", users[0].ID)
	assert.Equal(t, int32(2), queries.lastListParams.LimitRows)
	assert.False(t, queries.lastListParams.AfterCreatedAt.Valid)
}

func TestUserRepo_ListParsesCursorIntoGeneratedParams(t *testing.T) {
	queries := &fakeUserQueries{count: 0}
	repo := newUserRepoWithQueries(queries)
	cursorTime := "2026-06-05T11:12:13Z"
	cursor := base64.StdEncoding.EncodeToString([]byte(cursorTime))

	_, _, err := repo.List(context.Background(), nil, &cursor)

	require.NoError(t, err)
	require.True(t, queries.lastListParams.AfterCreatedAt.Valid)
	assert.Equal(t, cursorTime, queries.lastListParams.AfterCreatedAt.Time.Format(time.RFC3339))
	assert.Equal(t, int32(21), queries.lastListParams.LimitRows)
}

func TestUserRepo_UpdateMapsOptionalFieldsToNullableGeneratedParams(t *testing.T) {
	queries := &fakeUserQueries{
		updateRow: generated.UpdateUserRow{
			ID:        fakeUUID(t, "00000000-0000-0000-0000-000000000001"),
			Email:     "updated@example.com",
			Name:      "Updated",
			CreatedAt: fakeTime("2026-06-05T10:00:00Z"),
			UpdatedAt: fakeTime("2026-06-05T10:02:00Z"),
		},
	}
	repo := newUserRepoWithQueries(queries)
	name := "Updated"

	user, err := repo.Update(context.Background(), "00000000-0000-0000-0000-000000000001", service.UpdateUserInput{Name: &name})

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "Updated", user.Name)
	assert.True(t, queries.lastUpdateParams.Name.Valid)
	assert.Equal(t, "Updated", queries.lastUpdateParams.Name.String)
	assert.False(t, queries.lastUpdateParams.Email.Valid)
}

func TestUserRepo_ErrorMappingUsesGeneratedErrors(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{getErr: pgx.ErrNoRows})

		user, err := repo.GetByID(context.Background(), "00000000-0000-0000-0000-000000000001")

		require.NoError(t, err)
		assert.Nil(t, user)
	})

	t.Run("list query error", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{listErr: errors.New("list failed")})

		users, total, err := repo.List(context.Background(), nil, nil)

		require.Error(t, err)
		assert.Nil(t, users)
		assert.Zero(t, total)
		assert.Contains(t, err.Error(), "UserRepo.List")
	})

	t.Run("count error", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{countErr: errors.New("count failed")})

		users, total, err := repo.List(context.Background(), nil, nil)

		require.Error(t, err)
		assert.Nil(t, users)
		assert.Zero(t, total)
		assert.Contains(t, err.Error(), "count")
	})

	t.Run("create query error", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{createErr: errors.New("insert failed")})

		user, err := repo.Create(context.Background(), service.CreateUserInput{Email: "test@example.com"})

		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "UserRepo.Create")
	})

	t.Run("create duplicate email", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{createErr: &pgconn.PgError{Code: "23505"}})

		user, err := repo.Create(context.Background(), service.CreateUserInput{Email: "dupe@example.com"})

		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "duplicate email")
	})

	t.Run("update query error", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{updateErr: errors.New("update failed")})

		user, err := repo.Update(context.Background(), "00000000-0000-0000-0000-000000000001", service.UpdateUserInput{})

		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "UserRepo.Update")
	})

	t.Run("delete query error", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{deleteErr: errors.New("delete failed")})

		err := repo.Delete(context.Background(), "00000000-0000-0000-0000-000000000001")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "UserRepo.Delete")
	})
}
```

Run:

```bash
cd apps/api
go test ./internal/repository/postgres -run 'TestUserRepo_(ListMapsGeneratedRowsAndTrimsExtraRow|ListParsesCursorIntoGeneratedParams|UpdateMapsOptionalFieldsToNullableGeneratedParams|ErrorMappingUsesGeneratedErrors)' -count=1
```

Expected:

```text
FAIL
undefined: newUserRepoWithQueries
undefined: uuidFromString
```

- [ ] **Step 3: Replace UserRepo implementation**

Replace `apps/api/internal/repository/postgres/user_repo.go` with this implementation:

```go
// FILE: apps/api/internal/repository/postgres/user_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Adapt sqlc-generated PostgreSQL users queries to the service.UserRepository contract.
//   SCOPE: Users persistence CRUD, keyset pagination, generated-row mapping, and PostgreSQL error mapping; excludes transport handlers and service validation.
//   DEPENDS: apps/api/internal/repository/postgres/generated, github.com/jackc/pgx/v5, apps/api/internal/service, libs/go/logger.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   NewUserRepo - Constructs the production users repository from a pgx pool.
//   UserRepo.GetByID - Reads one user by UUID and maps missing rows to nil.
//   UserRepo.List - Reads a cursor page and total count.
//   UserRepo.Create - Inserts one user and maps duplicate email errors.
//   UserRepo.Update - Applies nullable name/email updates and maps missing rows to nil.
//   UserRepo.Delete - Deletes one user idempotently.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Switched users persistence from inline SQL to sqlc-generated queries.
// END_CHANGE_SUMMARY

package postgres

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/repository/postgres/generated"
	"monorepo-template/apps/api/internal/service"
	"monorepo-template/libs/go/logger"
)

type UserRepo struct {
	queries generated.Querier
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return newUserRepoWithQueries(generated.New(pool))
}

func newUserRepoWithQueries(queries generated.Querier) *UserRepo {
	return &UserRepo{queries: queries}
}

// START_CONTRACT: GetByID
//   PURPOSE: Retrieve a user by UUID and return nil when the row does not exist.
//   INPUTS: { ctx: context.Context - request context, id: string - user UUID }
//   OUTPUTS: { *service.User - mapped user or nil, error - query or UUID parsing failure }
//   SIDE_EFFECTS: Reads PostgreSQL through generated queries.
//   LINKS: M-API / V-M-API.
// END_CONTRACT: GetByID
func (r *UserRepo) GetByID(ctx context.Context, id string) (*service.User, error) {
	const op = "UserRepo.GetByID"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("querying user by id", zap.String("user_id", id))

	userID, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid user id: %w", op, err)
	}

	row, err := r.queries.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return userFromFields(row.ID, row.Email, row.Name, row.CreatedAt, row.UpdatedAt), nil
}

// START_CONTRACT: List
//   PURPOSE: Return a created_at-desc user page plus total row count.
//   INPUTS: { ctx: context.Context - request context, first: *int - page size, after: *string - base64 RFC3339Nano cursor }
//   OUTPUTS: { []*service.User - page rows, int - total rows, error - cursor, list, scan, or count failure }
//   SIDE_EFFECTS: Reads PostgreSQL through generated queries.
//   LINKS: M-API / V-M-API.
// END_CONTRACT: List
func (r *UserRepo) List(ctx context.Context, first *int, after *string) ([]*service.User, int, error) {
	const op = "UserRepo.List"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("querying users list")

	limit := int32(20)
	if first != nil && *first > 0 {
		limit = int32(*first)
	}

	cursor, err := cursorFromString(after)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.queries.ListUsers(ctx, generated.ListUsersParams{
		AfterCreatedAt: cursor,
		LimitRows:      limit + 1,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	users := make([]*service.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, userFromFields(row.ID, row.Email, row.Name, row.CreatedAt, row.UpdatedAt))
	}
	if len(users) > int(limit) {
		users = users[:int(limit)]
	}

	total, err := r.queries.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count: %w", op, err)
	}

	return users, int(total), nil
}

// START_CONTRACT: Create
//   PURPOSE: Insert one user row and map duplicate email conflicts.
//   INPUTS: { ctx: context.Context - request context, input: service.CreateUserInput - validated service input with password hash }
//   OUTPUTS: { *service.User - persisted user, error - duplicate email or insert failure }
//   SIDE_EFFECTS: Inserts PostgreSQL row through generated queries.
//   LINKS: M-API / V-M-API.
// END_CONTRACT: Create
func (r *UserRepo) Create(ctx context.Context, input service.CreateUserInput) (*service.User, error) {
	const op = "UserRepo.Create"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("inserting user", zap.String("email", input.Email))

	row, err := r.queries.CreateUser(ctx, generated.CreateUserParams{
		Email:        input.Email,
		Name:         input.Name,
		PasswordHash: input.Password,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, fmt.Errorf("%s: duplicate email: %s", op, input.Email)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return userFromFields(row.ID, row.Email, row.Name, row.CreatedAt, row.UpdatedAt), nil
}

// START_CONTRACT: Update
//   PURPOSE: Apply optional name and email changes and return nil when the user does not exist.
//   INPUTS: { ctx: context.Context - request context, id: string - user UUID, input: service.UpdateUserInput - optional fields }
//   OUTPUTS: { *service.User - updated user or nil, error - duplicate email, UUID parsing, or update failure }
//   SIDE_EFFECTS: Updates PostgreSQL row through generated queries.
//   LINKS: M-API / V-M-API.
// END_CONTRACT: Update
func (r *UserRepo) Update(ctx context.Context, id string, input service.UpdateUserInput) (*service.User, error) {
	const op = "UserRepo.Update"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("updating user", zap.String("user_id", id))

	userID, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid user id: %w", op, err)
	}

	row, err := r.queries.UpdateUser(ctx, generated.UpdateUserParams{
		ID:    userID,
		Name:  nullableText(input.Name),
		Email: nullableText(input.Email),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		if isDuplicateKeyError(err) {
			return nil, fmt.Errorf("%s: duplicate email: %s", op, derefString(input.Email))
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return userFromFields(row.ID, row.Email, row.Name, row.CreatedAt, row.UpdatedAt), nil
}

// START_CONTRACT: Delete
//   PURPOSE: Delete one user idempotently by UUID.
//   INPUTS: { ctx: context.Context - request context, id: string - user UUID }
//   OUTPUTS: { error - UUID parsing or delete failure }
//   SIDE_EFFECTS: Deletes PostgreSQL row through generated queries.
//   LINKS: M-API / V-M-API.
// END_CONTRACT: Delete
func (r *UserRepo) Delete(ctx context.Context, id string) error {
	const op = "UserRepo.Delete"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("deleting user", zap.String("user_id", id))

	userID, err := uuidFromString(id)
	if err != nil {
		return fmt.Errorf("%s: invalid user id: %w", op, err)
	}
	if err := r.queries.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func uuidFromString(value string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(value); err != nil {
		return pgtype.UUID{}, err
	}
	return uuid, nil
}

func cursorFromString(after *string) (pgtype.Timestamptz, error) {
	if after == nil || *after == "" {
		return pgtype.Timestamptz{}, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(*after)
	if err != nil {
		return pgtype.Timestamptz{}, fmt.Errorf("invalid cursor: %w", err)
	}
	cursor, err := time.Parse(time.RFC3339Nano, string(decoded))
	if err != nil {
		return pgtype.Timestamptz{}, fmt.Errorf("invalid cursor time: %w", err)
	}
	return pgtype.Timestamptz{Time: cursor, Valid: true}, nil
}

func nullableText(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *value, Valid: true}
}

func userFromFields(id pgtype.UUID, email string, name string, createdAt pgtype.Timestamptz, updatedAt pgtype.Timestamptz) *service.User {
	return &service.User{
		ID:        id.String(),
		Email:     email,
		Name:      name,
		CreatedAt: formatTimestamp(createdAt),
		UpdatedAt: formatTimestamp(updatedAt),
	}
}

func formatTimestamp(value pgtype.Timestamptz) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format(time.RFC3339Nano)
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
```

- [ ] **Step 4: Run focused adapter unit tests**

Run:

```bash
cd apps/api
go test ./internal/repository/postgres -run 'TestUserRepo_(ListMapsGeneratedRowsAndTrimsExtraRow|ListParsesCursorIntoGeneratedParams|UpdateMapsOptionalFieldsToNullableGeneratedParams|ErrorMappingUsesGeneratedErrors)' -count=1
```

Expected:

```text
ok  	monorepo-template/apps/api/internal/repository/postgres
```

- [ ] **Step 5: Run the full postgres repository package tests**

Run:

```bash
cd apps/api
go test ./internal/repository/postgres -count=1
```

Expected:

```text
ok  	monorepo-template/apps/api/internal/repository/postgres
```

If the Docker test database is not running, integration tests may skip outside the coverage gate. Continue to Task 3 to force real integration evidence.

- [ ] **Step 6: Checkpoint adapter migration without committing**

Run:

```bash
git status --short -- apps/api/internal/repository/postgres/user_repo.go apps/api/internal/repository/postgres/user_repo_unit_test.go
```

Expected:

```text
# status shows only the intended adapter and unit-test files from this task
```

Do not commit yet. GRACE docs and verification evidence are not synchronized until Tasks 6-7.

## Task 3: Expand Repository Integration Coverage

**Files:**

- Modify: `apps/api/internal/repository/postgres/user_repo_test.go`

- [ ] **Step 1: Add integration-test file contract and nullable update behavior**

If `apps/api/internal/repository/postgres/user_repo_test.go` does not already have file-local GRACE markup, add this header before `package postgres_test`:

```go
// FILE: apps/api/internal/repository/postgres/user_repo_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify PostgreSQL users repository behavior against the goose-managed test database.
//   SCOPE: Users CRUD, pagination, duplicate handling, nullable update paths, safe destructive setup, and unavailable database skip semantics; excludes transport and service validation.
//   DEPENDS: apps/api/internal/repository/postgres, apps/api/internal/service, apps/api/internal/testinfra, docker/docker-compose.test.yml.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   testPool - Applies migrations, enforces safe test DSN, and truncates users for isolated integration tests.
//   TestUserRepo_* - Real database coverage for users repository behavior.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added sqlc nullable-update integration coverage.
// END_CHANGE_SUMMARY
```

Add these tests after `TestUserRepo_UpdateDuplicateEmail`:

```go
func TestUserRepo_UpdateNameOnlyPreservesEmail(t *testing.T) {
	ctx := context.Background()
	repo := postgresrepo.NewUserRepo(testPool(t))
	created, err := repo.Create(ctx, service.CreateUserInput{
		Email:    "name-only@example.com",
		Name:     "Original",
		Password: "$2a$10$hashed",
	})
	require.NoError(t, err)
	name := "Renamed"

	updated, err := repo.Update(ctx, created.ID, service.UpdateUserInput{Name: &name})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Renamed", updated.Name)
	assert.Equal(t, "name-only@example.com", updated.Email)
}

func TestUserRepo_UpdateEmailOnlyPreservesName(t *testing.T) {
	ctx := context.Background()
	repo := postgresrepo.NewUserRepo(testPool(t))
	created, err := repo.Create(ctx, service.CreateUserInput{
		Email:    "email-only@example.com",
		Name:     "Original",
		Password: "$2a$10$hashed",
	})
	require.NoError(t, err)
	email := "email-only-updated@example.com"

	updated, err := repo.Update(ctx, created.ID, service.UpdateUserInput{Email: &email})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Original", updated.Name)
	assert.Equal(t, "email-only-updated@example.com", updated.Email)
}

func TestUserRepo_UpdateEmptyInputPreservesNameAndEmail(t *testing.T) {
	ctx := context.Background()
	repo := postgresrepo.NewUserRepo(testPool(t))
	created, err := repo.Create(ctx, service.CreateUserInput{
		Email:    "empty-update@example.com",
		Name:     "Original",
		Password: "$2a$10$hashed",
	})
	require.NoError(t, err)

	updated, err := repo.Update(ctx, created.ID, service.UpdateUserInput{})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Original", updated.Name)
	assert.Equal(t, "empty-update@example.com", updated.Email)
}
```

- [ ] **Step 2: Run integration tests without the test database to confirm skip behavior is explicit**

Run:

```bash
cd apps/api
TEST_POSTGRES_PORT=1 \
API_TEST_DATABASE_DSN=postgres://app:secret@localhost:1/monorepo_test?sslmode=disable \
  go test ./internal/repository/postgres -run 'TestUserRepo_Update(NameOnly|EmailOnly|EmptyInput)' -count=1 -v
```

Expected:

```text
SKIP: postgres integration database is unavailable
```

The skip is acceptable only outside coverage mode. The next step proves the real path.

- [ ] **Step 3: Start the dedicated test database**

Run:

```bash
TEST_COMPOSE_PROJECT=mt-sqlc-ref \
TEST_POSTGRES_CONTAINER_NAME=mt-sqlc-ref-postgres \
TEST_REDIS_CONTAINER_NAME=mt-sqlc-ref-redis \
TEST_POSTGRES_VOLUME=mt-sqlc-ref-pg-test-data \
docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis
```

Expected:

```text
Container mt-sqlc-ref-postgres  Healthy
Container mt-sqlc-ref-redis     Healthy
```

Reuse this same compose scope for all later sqlc repository verification steps that bind PostgreSQL to host port `17501`; do not switch project names while the container is running.

- [ ] **Step 4: Run repository integration tests against the real test database**

Run:

```bash
cd apps/api
COVERAGE_GATE=1 \
API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable \
go test ./internal/repository/postgres -run 'TestUserRepo_(CreateGetListUpdateDelete|ListTrimsLimitAndAcceptsCursor|CreateDuplicateEmail|UpdateDuplicateEmail|UpdateNameOnlyPreservesEmail|UpdateEmailOnlyPreservesName|UpdateEmptyInputPreservesNameAndEmail|UpdateMissingReturnsNil|DeleteMissingReturnsNil)' -count=1 -v
```

Expected:

```text
PASS
ok  	monorepo-template/apps/api/internal/repository/postgres
```

This is the required evidence that sqlc generated queries work against the goose-created schema on the safe `monorepo_test` target.

- [ ] **Step 5: Checkpoint integration tests without committing**

Run:

```bash
git status --short -- apps/api/internal/repository/postgres/user_repo_test.go
```

Expected:

```text
# status shows only the intended integration-test file from this task
```

Do not commit yet. GRACE docs and verification evidence are not synchronized until Tasks 6-7.

## Task 4: Add Generated Drift Detection

**Files:**

- Modify: `tools/codegen/project.json`

- [ ] **Step 1: Update codegen validation target to assert clean generated output**

`tools/codegen/project.json` is strict JSON. Do not add comments or file-local markup to it; capture the ownership change in `docs/development-plan.xml`, `docs/knowledge-graph.xml`, and `docs/verification-plan.xml` in Task 6.

Modify `tools/codegen/project.json` so the `validate` target runs root codegen and checks generated paths:

```json
{
  "name": "codegen",
  "$schema": "../../node_modules/nx/schemas/project-schema.json",
  "sourceRoot": "tools/codegen",
  "projectType": "library",
  "targets": {
    "validate": {
      "executor": "nx:run-commands",
      "options": {
        "command": "bunx vitest run --config tools/vitest.config.ts --coverage && bun run codegen && git diff --exit-code -- apps/api/internal/repository/postgres/generated apps/api/internal/graph apps/web-admin/src/shared/api/generated"
      }
    }
  }
}
```

- [ ] **Step 2: Verify drift detection is green on current generated output**

Run:

```bash
bunx nx run codegen:validate
```

Expected:

```text
> nx run codegen:validate
# command completes successfully
```

- [ ] **Step 3: Prove the drift gate fails when query sources change generated output**

Run:

```bash
tmp_users_sql="$(mktemp)"
cp apps/api/internal/repository/postgres/queries/users.sql "$tmp_users_sql"
perl -0pi -e 's/-- name: GetUserByID :one/-- name: GetUserByIDWithPassword :one/' apps/api/internal/repository/postgres/queries/users.sql
bunx nx run codegen:validate
```

Expected:

```text
FAIL
```

Restore generated output:

```bash
cp "$tmp_users_sql" apps/api/internal/repository/postgres/queries/users.sql
bun run codegen
git diff --exit-code -- apps/api/internal/repository/postgres/queries/users.sql apps/api/internal/repository/postgres/generated
```

Expected:

```text
# no output
```

- [ ] **Step 4: Checkpoint drift validation without committing**

Run:

```bash
git status --short -- tools/codegen/project.json
```

Expected:

```text
# status shows only the intended codegen validation config from this task
```

Do not commit yet. GRACE docs and verification evidence are not synchronized until Tasks 6-7.

## Task 5: Add Coverage Policy For sqlc Generated Go

**Files:**

- Modify: `tools/coverage/coverage.config.json`

- [ ] **Step 1: Add exact generated sqlc files to the allowlist**

`tools/coverage/coverage.config.json` is strict JSON. Do not add comments or file-local markup to it; capture the generated-code contract in `docs/verification-plan.xml` and the plan-specific verification evidence.

Add these entries to `tools/coverage/coverage.config.json` after the existing gqlgen API entries:

```json
{
  "path": "apps/api/internal/repository/postgres/generated/db.go",
  "reason": "sqlc generated DBTX and Queries wrapper",
  "gate": "bunx nx run api:codegen && bunx nx build api && TEST_COMPOSE_PROJECT=mt-sqlc-ref TEST_POSTGRES_CONTAINER_NAME=mt-sqlc-ref-postgres TEST_REDIS_CONTAINER_NAME=mt-sqlc-ref-redis TEST_POSTGRES_VOLUME=mt-sqlc-ref-pg-test-data docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis && cd apps/api && COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable go test ./internal/repository/postgres -run TestUserRepo -count=1 -v"
},
{
  "path": "apps/api/internal/repository/postgres/generated/models.go",
  "reason": "sqlc generated table models",
  "gate": "bunx nx run api:codegen && bunx nx build api && TEST_COMPOSE_PROJECT=mt-sqlc-ref TEST_POSTGRES_CONTAINER_NAME=mt-sqlc-ref-postgres TEST_REDIS_CONTAINER_NAME=mt-sqlc-ref-redis TEST_POSTGRES_VOLUME=mt-sqlc-ref-pg-test-data docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis && cd apps/api && COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable go test ./internal/repository/postgres -run TestUserRepo -count=1 -v"
},
{
  "path": "apps/api/internal/repository/postgres/generated/querier.go",
  "reason": "sqlc generated query interface",
  "gate": "bunx nx run api:codegen && bunx nx build api && TEST_COMPOSE_PROJECT=mt-sqlc-ref TEST_POSTGRES_CONTAINER_NAME=mt-sqlc-ref-postgres TEST_REDIS_CONTAINER_NAME=mt-sqlc-ref-redis TEST_POSTGRES_VOLUME=mt-sqlc-ref-pg-test-data docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis && cd apps/api && COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable go test ./internal/repository/postgres -run TestUserRepo -count=1 -v"
},
{
  "path": "apps/api/internal/repository/postgres/generated/users.sql.go",
  "reason": "sqlc generated users query methods",
  "gate": "bunx nx run api:codegen && bunx nx build api && TEST_COMPOSE_PROJECT=mt-sqlc-ref TEST_POSTGRES_CONTAINER_NAME=mt-sqlc-ref-postgres TEST_REDIS_CONTAINER_NAME=mt-sqlc-ref-redis TEST_POSTGRES_VOLUME=mt-sqlc-ref-pg-test-data docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis && cd apps/api && COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable go test ./internal/repository/postgres -run TestUserRepo -count=1 -v"
}
```

Keep exact `.go` paths. Do not use a glob here because the current `tools/coverage/run.mjs` Go allowlist logic checks entries ending in `.go`.

- [ ] **Step 2: Run the coverage preflight unit for allowlist parsing through the full coverage runner dry surface**

Run:

```bash
node -e "const c=require('./tools/coverage/coverage.config.json'); const paths=c.allowlist.map(x=>x.path); for (const p of ['apps/api/internal/repository/postgres/generated/db.go','apps/api/internal/repository/postgres/generated/models.go','apps/api/internal/repository/postgres/generated/querier.go','apps/api/internal/repository/postgres/generated/users.sql.go']) { if (!paths.includes(p)) { throw new Error('missing '+p) } } console.log('sqlc allowlist ok')"
```

Expected:

```text
sqlc allowlist ok
```

- [ ] **Step 3: Run focused API coverage with the safe database up**

Run:

```bash
TEST_COMPOSE_PROJECT=mt-sqlc-ref \
TEST_POSTGRES_CONTAINER_NAME=mt-sqlc-ref-postgres \
TEST_REDIS_CONTAINER_NAME=mt-sqlc-ref-redis \
TEST_POSTGRES_VOLUME=mt-sqlc-ref-pg-test-data \
docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis

cd apps/api
mkdir -p ../../dist/coverage/go/api
COVERAGE_GATE=1 \
API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable \
go test -coverprofile=../../dist/coverage/go/api/coverage.out ./...
```

Expected:

```text
ok  	monorepo-template/apps/api/...
```

The command must not contain `SKIP`; `COVERAGE_GATE=1` must turn unavailable database setup into a failure.

- [ ] **Step 4: Checkpoint coverage policy without committing**

Run:

```bash
git status --short -- tools/coverage/coverage.config.json
```

Expected:

```text
# status shows only the intended coverage config from this task
```

Do not commit yet. GRACE docs and verification evidence are not synchronized until Tasks 6-7.

## Task 6: Synchronize GRACE Shared Docs

**Files:**

- Modify: `docs/requirements.xml`
- Modify: `docs/technology.xml`
- Modify: `docs/development-plan.xml`
- Modify: `docs/knowledge-graph.xml`
- Modify: `docs/verification-plan.xml`

- [ ] **Step 1: Update requirements codegen use case**

In `docs/requirements.xml`, update `UC-004` so it includes backend SQL generation:

```xml
    <UC-004>
      <Actor>Developer</Actor>
      <Action>Regenerates API SQL, API GraphQL, and web-admin GraphQL artifacts after schema or query changes.</Action>
      <Goal>Keep sqlc query output, gqlgen output, and web-admin client types aligned with the migrations, query SQL, and admin GraphQL schema.</Goal>
      <Preconditions>`apps/api/sqlc.yaml`, `apps/api/internal/repository/postgres/migrations/*.sql`, `apps/api/internal/repository/postgres/queries/*.sql`, `libs/graphql/schema/*.graphql`, `apps/api/gqlgen.yml`, and `tools/codegen/codegen.ts` are consistent.</Preconditions>
      <AcceptanceCriteria>`bun run codegen` completes, `bunx nx run codegen:validate` fails on generated drift, and generated files match the schema/query sources used by API and web-admin projects.</AcceptanceCriteria>
      <Priority>high</Priority>
      <RelatedFlows>DF-CODEGEN</RelatedFlows>
    </UC-004>
```

Also replace open question 2 with a resolved generated-files policy:

```xml
    <question-2>Resolved 2026-06-05: API sqlc output, API gqlgen output, and web-admin client output stay committed; CI and local validation regenerate and compare generated artifacts to catch drift.</question-2>
```

- [ ] **Step 2: Update technology stack for sqlc**

In `docs/technology.xml`, add a dependency and tool entry:

```xml
    <dep name="github.com/sqlc-dev/sqlc" version="v1.30.0" purpose="Typed Go query generation for PostgreSQL users repository code" />
```

```xml
    <tool name="api-sql-codegen" value="sqlc" version="v1.30.0" />
```

Update the preferred stack:

```xml
    <preferred-api-sql-codegen>github.com/sqlc-dev/sqlc using pgx/v5 output from goose migrations plus repository query SQL</preferred-api-sql-codegen>
```

- [ ] **Step 3: Update development-plan modules**

In `docs/development-plan.xml`, update the relevant module excerpts.

For `M-WORKSPACE`, change the codegen interface purpose:

```xml
        <export-bun-scripts PURPOSE="Expose dev, build, test, lint, SQL/API/web codegen, and coverage scripts." />
```

For `M-COVERAGE-GATE`, add sqlc generated code to the inputs:

```xml
          <param name="generated-allowlist" type="gqlgen output, sqlc output, web-admin GraphQL output, and bootstrap entrypoints listed in tools/coverage/coverage.config.json" />
```

For `M-API`, update purpose and interface:

```xml
        <purpose>Serve HTTP health, admin GraphQL, and public REST endpoints for the reference user domain with PostgreSQL, Redis, auth, CORS, goose migrations, sqlc-generated user queries, and shared user service logic.</purpose>
```

```xml
        <export-user-repository PURPOSE="Expose PostgreSQL users repository adapter backed by sqlc-generated pgx/v5 queries." />
```

Add API target sources:

```xml
        <source>apps/api/sqlc.yaml</source>
        <source>apps/api/internal/repository/postgres/queries</source>
        <source>apps/api/internal/repository/postgres/generated</source>
```

- [ ] **Step 4: Update knowledge graph paths and cross-links**

In `docs/knowledge-graph.xml`, update `M-WORKSPACE` codegen annotation:

```xml
        <export-codegen PURPOSE="Run API sqlc, API gqlgen, and web-admin GraphQL codegen targets with generated-drift validation." />
```

Add these paths to `M-API`:

```xml
      <path>apps/api/sqlc.yaml</path>
      <path>apps/api/internal/repository/postgres/queries</path>
      <path>apps/api/internal/repository/postgres/generated</path>
```

Add this annotation to `M-API`:

```xml
        <fn-userRepository PURPOSE="PostgreSQL users repository adapter backed by sqlc-generated pgx/v5 queries." />
```

Add this cross-link:

```xml
    <CrossLink from="M-COVERAGE-GATE" to="M-API" relation="allowlists sqlc generated users query output with API codegen, build, and repository integration replacement gates" />
```

- [ ] **Step 5: Update verification plan**

In `docs/verification-plan.xml`, update `V-M-WORKSPACE` module checks to include codegen validation:

```xml
        <check-7>bunx nx run codegen:validate</check-7>
```

In `V-M-COVERAGE-GATE`, add allowlist entries matching `tools/coverage/coverage.config.json`:

```xml
        <entry path="apps/api/internal/repository/postgres/generated/db.go" reason="sqlc generated DBTX and Queries wrapper" gate="bunx nx run api:codegen &amp;&amp; bunx nx build api &amp;&amp; TEST_COMPOSE_PROJECT=mt-sqlc-ref TEST_POSTGRES_CONTAINER_NAME=mt-sqlc-ref-postgres TEST_REDIS_CONTAINER_NAME=mt-sqlc-ref-redis TEST_POSTGRES_VOLUME=mt-sqlc-ref-pg-test-data docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis &amp;&amp; cd apps/api &amp;&amp; COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable go test ./internal/repository/postgres -run TestUserRepo -count=1 -v" />
        <entry path="apps/api/internal/repository/postgres/generated/models.go" reason="sqlc generated table models" gate="bunx nx run api:codegen &amp;&amp; bunx nx build api &amp;&amp; TEST_COMPOSE_PROJECT=mt-sqlc-ref TEST_POSTGRES_CONTAINER_NAME=mt-sqlc-ref-postgres TEST_REDIS_CONTAINER_NAME=mt-sqlc-ref-redis TEST_POSTGRES_VOLUME=mt-sqlc-ref-pg-test-data docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis &amp;&amp; cd apps/api &amp;&amp; COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable go test ./internal/repository/postgres -run TestUserRepo -count=1 -v" />
        <entry path="apps/api/internal/repository/postgres/generated/querier.go" reason="sqlc generated query interface" gate="bunx nx run api:codegen &amp;&amp; bunx nx build api &amp;&amp; TEST_COMPOSE_PROJECT=mt-sqlc-ref TEST_POSTGRES_CONTAINER_NAME=mt-sqlc-ref-postgres TEST_REDIS_CONTAINER_NAME=mt-sqlc-ref-redis TEST_POSTGRES_VOLUME=mt-sqlc-ref-pg-test-data docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis &amp;&amp; cd apps/api &amp;&amp; COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable go test ./internal/repository/postgres -run TestUserRepo -count=1 -v" />
        <entry path="apps/api/internal/repository/postgres/generated/users.sql.go" reason="sqlc generated users query methods" gate="bunx nx run api:codegen &amp;&amp; bunx nx build api &amp;&amp; TEST_COMPOSE_PROJECT=mt-sqlc-ref TEST_POSTGRES_CONTAINER_NAME=mt-sqlc-ref-postgres TEST_REDIS_CONTAINER_NAME=mt-sqlc-ref-redis TEST_POSTGRES_VOLUME=mt-sqlc-ref-pg-test-data docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis &amp;&amp; cd apps/api &amp;&amp; COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable go test ./internal/repository/postgres -run TestUserRepo -count=1 -v" />
```

In `V-M-API`, add repository files and codegen checks:

```xml
        <file>apps/api/sqlc.yaml</file>
        <file>apps/api/internal/repository/postgres/queries/users.sql</file>
        <file>apps/api/internal/repository/postgres/user_repo_test.go</file>
        <file>apps/api/internal/repository/postgres/user_repo_unit_test.go</file>
```

```xml
        <check-5>bunx nx run api:codegen</check-5>
        <check-6>cd apps/api &amp;&amp; COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable go test ./internal/repository/postgres -run TestUserRepo -count=1 -v</check-6>
```

Add this trace assertion under `V-M-API`:

```xml
        <assertion-4>Users repository tests must prove sqlc-generated queries against the goose-created `monorepo_test` schema before final closeout.</assertion-4>
```

- [ ] **Step 6: Validate XML and GRACE docs**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected:

```text
# xmllint exits 0
# grace lint exits 0
```

- [ ] **Step 7: Checkpoint GRACE docs without committing**

Run:

```bash
git status --short -- docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml
```

Expected:

```text
# status shows only the intended GRACE docs from this task
```

Do not commit yet. Final verification and one synchronized commit happen in Task 7.

## Task 7: Final Verification And Handoff

**Files:**

- Verify: all files changed in Tasks 1-6.

- [ ] **Step 1: Run focused API codegen and tests**

Run:

```bash
bunx nx run api:codegen
cd apps/api && go test ./internal/repository/postgres -run TestUserRepo -count=1
bunx nx test api
bunx nx build api
```

Expected:

```text
# api:codegen exits 0
ok  	monorepo-template/apps/api/internal/repository/postgres
# nx test api exits 0
# nx build api exits 0
```

- [ ] **Step 2: Run generated drift validation**

Run:

```bash
bunx nx run codegen:validate
```

Expected:

```text
> nx run codegen:validate
# exits 0 and leaves no generated diff
```

Confirm:

```bash
git diff --exit-code -- apps/api/internal/repository/postgres/generated apps/api/internal/graph apps/web-admin/src/shared/api/generated
```

Expected:

```text
# no output
```

- [ ] **Step 3: Run safe database integration evidence**

Run:

```bash
TEST_COMPOSE_PROJECT=mt-sqlc-ref \
TEST_POSTGRES_CONTAINER_NAME=mt-sqlc-ref-postgres \
TEST_REDIS_CONTAINER_NAME=mt-sqlc-ref-redis \
TEST_POSTGRES_VOLUME=mt-sqlc-ref-pg-test-data \
docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis

cd apps/api
COVERAGE_GATE=1 \
API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable \
go test ./internal/repository/postgres -run TestUserRepo -count=1 -v
```

Expected:

```text
PASS
ok  	monorepo-template/apps/api/internal/repository/postgres
```

The command must not contain `SKIP`. If it skips, the sqlc/goose integration path is not proven.

- [ ] **Step 4: Run GRACE validation**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected:

```text
# both commands exit 0
```

- [ ] **Step 5: Run final broad gate if environment supports Docker and Playwright**

Run:

```bash
bun run verify:coverage
```

Expected:

```text
[Coverage][gate] all thresholds passed
```

If this fails because Docker, Playwright, or the local toolchain is unavailable, capture the exact environment failure and keep the focused gates from Steps 1-4 as the implementation evidence.

- [ ] **Step 6: Inspect the final diff**

Run:

```bash
git status --short
git diff --stat
git diff --check
```

Expected:

```text
# git diff --check exits 0
# status shows only files intentionally changed by this plan
```

- [ ] **Step 7: Commit the synchronized rollout**

Run:

```bash
git add apps/api/sqlc.yaml \
  apps/api/internal/repository/postgres/queries/users.sql \
  apps/api/internal/repository/postgres/generated \
  apps/api/tools.go apps/api/go.mod apps/api/go.sum apps/api/project.json \
  apps/api/internal/repository/postgres/user_repo.go \
  apps/api/internal/repository/postgres/user_repo_unit_test.go \
  apps/api/internal/repository/postgres/user_repo_test.go \
  tools/codegen/project.json tools/coverage/coverage.config.json \
  docs/requirements.xml docs/technology.xml docs/development-plan.xml \
  docs/knowledge-graph.xml docs/verification-plan.xml docs/operational-packets.xml
git commit -m "feat(api): adopt sqlc users reference"
```

Expected:

```text
git prints a new commit line containing "feat(api): adopt sqlc users reference"
```

## Rollback Plan

- Revert the implementation commits from this plan.
- Do not run `goose down`; no schema migration is part of the rollout.
- After revert, run:

```bash
bunx nx run api:codegen
bunx nx test api
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected:

```text
# all commands exit 0
```
