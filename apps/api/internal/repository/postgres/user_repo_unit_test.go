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
//   TestAdminUsersMigrationVersionAfterHistoricalLocalVersions - Guards stale local goose histories.
//   TestRunMigrations_ReturnsOpenError - Migration driver failure coverage.
//   TestNew_ReturnsPoolConstructionError - pgx pool construction failure coverage.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added admin migration ordering regression coverage for stale local dev DBs.
// END_CHANGE_SUMMARY

package postgres

import (
	"context"
	"encoding/base64"
	"errors"
	"math"
	"os"
	"strconv"
	"strings"
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

func TestAdminUsersMigrationVersionAfterHistoricalLocalVersions(t *testing.T) {
	entries, err := os.ReadDir("migrations")
	require.NoError(t, err)

	for _, entry := range entries {
		name := entry.Name()
		if !strings.Contains(name, "_admin_users.sql") {
			continue
		}
		versionText, _, found := strings.Cut(name, "_")
		require.True(t, found)
		version, err := strconv.Atoi(versionText)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, version, 79)
		return
	}

	t.Fatal("admin_users migration not found")
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

func TestUserRepo_ListRejectsLimitThatWouldOverflowGeneratedParams(t *testing.T) {
	queries := &fakeUserQueries{}
	repo := newUserRepoWithQueries(queries)
	first := math.MaxInt32

	users, total, err := repo.List(context.Background(), &first, nil)

	require.Error(t, err)
	assert.Nil(t, users)
	assert.Zero(t, total)
	assert.Contains(t, err.Error(), "first exceeds max supported page size")
	assert.Zero(t, queries.lastListParams.LimitRows)
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
	t.Run("get invalid id", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{})

		user, err := repo.GetByID(context.Background(), "not-a-uuid")

		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid user id")
	})

	t.Run("not found", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{getErr: pgx.ErrNoRows})

		user, err := repo.GetByID(context.Background(), "00000000-0000-0000-0000-000000000001")

		require.NoError(t, err)
		assert.Nil(t, user)
	})

	t.Run("get query error", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{getErr: errors.New("select failed")})

		user, err := repo.GetByID(context.Background(), "00000000-0000-0000-0000-000000000001")

		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "UserRepo.GetByID")
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

	t.Run("update invalid id", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{})

		user, err := repo.Update(context.Background(), "not-a-uuid", service.UpdateUserInput{})

		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid user id")
	})

	t.Run("update query error", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{updateErr: errors.New("update failed")})

		user, err := repo.Update(context.Background(), "00000000-0000-0000-0000-000000000001", service.UpdateUserInput{})

		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "UserRepo.Update")
	})

	t.Run("delete invalid id", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{})

		err := repo.Delete(context.Background(), "not-a-uuid")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user id")
	})

	t.Run("delete query error", func(t *testing.T) {
		repo := newUserRepoWithQueries(&fakeUserQueries{deleteErr: errors.New("delete failed")})

		err := repo.Delete(context.Background(), "00000000-0000-0000-0000-000000000001")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "UserRepo.Delete")
	})
}

func TestUserRepo_MapsInvalidGeneratedTimestampsToEmptyStrings(t *testing.T) {
	user := userFromFields(
		fakeUUID(t, "00000000-0000-0000-0000-000000000001"),
		"empty-time@example.com",
		"Empty Time",
		pgtype.Timestamptz{},
		pgtype.Timestamptz{},
	)

	assert.Equal(t, "", user.CreatedAt)
	assert.Equal(t, "", user.UpdatedAt)
}

func TestDuplicateHelpers(t *testing.T) {
	assert.False(t, isDuplicateKeyError(errors.New("plain error")))
	assert.Equal(t, "", derefString(nil))
	value := "value"
	assert.Equal(t, "value", derefString(&value))
}

func TestRunMigrations_ReturnsOpenError(t *testing.T) {
	previous := migrationDriver
	migrationDriver = "missing-driver"
	t.Cleanup(func() {
		migrationDriver = previous
	})

	dsn := testinfra.PostgresDSN()
	testinfra.RequireSafePostgresDSN(t, dsn)
	err := RunMigrations(dsn, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "open db for migrations")
}

func TestNew_ReturnsPoolConstructionError(t *testing.T) {
	previous := newPoolWithConfig
	newPoolWithConfig = func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
		assert.Equal(t, time.Second, config.MaxConnIdleTime)
		return nil, errors.New("pool failed")
	}
	t.Cleanup(func() {
		newPoolWithConfig = previous
	})

	cfg := testinfra.PostgresConfig(t)
	cfg.MaxConnIdleTime = time.Second
	db, err := New(cfg, zap.NewNop())

	require.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "failed to connect")
}
