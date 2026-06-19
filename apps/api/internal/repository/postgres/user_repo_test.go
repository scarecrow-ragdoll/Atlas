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

package postgres_test

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	postgresrepo "monorepo-template/apps/api/internal/repository/postgres"
	"monorepo-template/apps/api/internal/service"
	"monorepo-template/apps/api/internal/testinfra"
)

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := testinfra.PostgresDSN()
	testinfra.RequireSafePostgresDSN(t, dsn)
	if err := postgresrepo.RunMigrations(dsn, zap.NewNop()); err != nil {
		if !testinfra.CoverageGateEnabled() {
			t.Skipf("postgres integration database is unavailable: %v", err)
		}
		require.NoError(t, err)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	_, err = pool.Exec(context.Background(), `TRUNCATE users RESTART IDENTITY CASCADE`)
	require.NoError(t, err)
	return pool
}

func TestUserRepo_CreateGetListUpdateDelete(t *testing.T) {
	ctx := context.Background()
	repo := postgresrepo.NewUserRepo(testPool(t))

	created, err := repo.Create(ctx, service.CreateUserInput{
		Email:    "created@example.com",
		Name:     "Created User",
		Password: "$2a$10$hashed",
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	found, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, "created@example.com", found.Email)

	name := "Updated User"
	email := "updated@example.com"
	updated, err := repo.Update(ctx, created.ID, service.UpdateUserInput{Name: &name, Email: &email})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Updated User", updated.Name)
	assert.Equal(t, "updated@example.com", updated.Email)

	users, total, err := repo.List(ctx, ptr(20), nil)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.Len(t, users, 1)
	assert.Equal(t, updated.ID, users[0].ID)

	require.NoError(t, repo.Delete(ctx, created.ID))
	deleted, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Nil(t, deleted)
}

func TestUserRepo_ListTrimsLimitAndAcceptsCursor(t *testing.T) {
	ctx := context.Background()
	repo := postgresrepo.NewUserRepo(testPool(t))
	_, err := repo.Create(ctx, service.CreateUserInput{
		Email:    "one@example.com",
		Name:     "One",
		Password: "$2a$10$hashed",
	})
	require.NoError(t, err)
	_, err = repo.Create(ctx, service.CreateUserInput{
		Email:    "two@example.com",
		Name:     "Two",
		Password: "$2a$10$hashed",
	})
	require.NoError(t, err)
	cursor := base64.StdEncoding.EncodeToString([]byte(time.Now().Add(time.Hour).Format(time.RFC3339Nano)))

	users, total, err := repo.List(ctx, ptr(1), &cursor)

	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, users, 1)
}

func TestUserRepo_CreateDuplicateEmail(t *testing.T) {
	ctx := context.Background()
	repo := postgresrepo.NewUserRepo(testPool(t))
	input := service.CreateUserInput{
		Email:    "duplicate@example.com",
		Name:     "First",
		Password: "$2a$10$hashed",
	}
	_, err := repo.Create(ctx, input)
	require.NoError(t, err)

	_, err = repo.Create(ctx, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate email")
}

func TestUserRepo_UpdateDuplicateEmail(t *testing.T) {
	ctx := context.Background()
	repo := postgresrepo.NewUserRepo(testPool(t))
	first, err := repo.Create(ctx, service.CreateUserInput{
		Email:    "first@example.com",
		Name:     "First",
		Password: "$2a$10$hashed",
	})
	require.NoError(t, err)
	_, err = repo.Create(ctx, service.CreateUserInput{
		Email:    "second@example.com",
		Name:     "Second",
		Password: "$2a$10$hashed",
	})
	require.NoError(t, err)
	duplicate := "second@example.com"

	_, err = repo.Update(ctx, first.ID, service.UpdateUserInput{Email: &duplicate})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate email")
}

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

func TestUserRepo_UpdateMissingReturnsNil(t *testing.T) {
	repo := postgresrepo.NewUserRepo(testPool(t))
	name := "Nobody"

	user, err := repo.Update(context.Background(), "00000000-0000-0000-0000-000000000000", service.UpdateUserInput{Name: &name})

	require.NoError(t, err)
	assert.Nil(t, user)
}

func TestUserRepo_DeleteMissingReturnsNil(t *testing.T) {
	repo := postgresrepo.NewUserRepo(testPool(t))

	err := repo.Delete(context.Background(), "00000000-0000-0000-0000-000000000000")

	require.NoError(t, err)
}

func TestUserRepo_ListInvalidCursor(t *testing.T) {
	repo := postgresrepo.NewUserRepo(testPool(t))
	invalid := "not-base64"

	_, _, err := repo.List(context.Background(), nil, &invalid)

	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "invalid cursor"))
}

func TestUserRepo_ListInvalidCursorTime(t *testing.T) {
	repo := postgresrepo.NewUserRepo(testPool(t))
	invalid := "bm90LWEtdGltZQ=="

	_, _, err := repo.List(context.Background(), nil, &invalid)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid cursor time")
}

func ptr[T any](v T) *T { return &v }
