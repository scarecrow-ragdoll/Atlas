// FILE: apps/api/internal/repository/postgres/admin_repo_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify PostgreSQL admin repository behavior against the goose-managed test database.
//   SCOPE: Admin create/count/get, case-insensitive duplicate handling, active flag mapping, and safe destructive setup; excludes service validation and Redis sessions.
//   DEPENDS: apps/api/internal/repository/postgres, apps/api/internal/service, apps/api/internal/testinfra, docker/docker-compose.test.yml.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   adminTestPool - Applies migrations, enforces safe test DSN, and truncates admin_users.
//   TestAdminRepo_* - Real database coverage for admin identity persistence.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin repository integration coverage.
// END_CHANGE_SUMMARY

package postgres_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	postgresrepo "monorepo-template/apps/api/internal/repository/postgres"
	"monorepo-template/apps/api/internal/service"
	"monorepo-template/apps/api/internal/testinfra"
)

func adminTestPool(t *testing.T) *pgxpool.Pool {
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
	_, err = pool.Exec(context.Background(), `TRUNCATE admin_users RESTART IDENTITY CASCADE`)
	require.NoError(t, err)
	return pool
}

func TestAdminRepo_CreateCountAndLookup(t *testing.T) {
	ctx := context.Background()
	repo := postgresrepo.NewAdminRepo(adminTestPool(t))

	count, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	created, err := repo.Create(ctx, service.CreateAdminInput{
		Email:        "Admin@Example.COM",
		Name:         "Template Admin",
		PasswordHash: "$2a$10$hashed",
		Role:         service.AdminRoleAdmin,
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "admin@example.com", created.Email)
	assert.True(t, created.IsActive)

	foundByEmail, err := repo.GetByEmail(ctx, "ADMIN@example.com")
	require.NoError(t, err)
	require.NotNil(t, foundByEmail)
	assert.Equal(t, created.ID, foundByEmail.ID)

	foundByID, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, foundByID)
	assert.Equal(t, created.Email, foundByID.Email)

	count, err = repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestAdminRepo_CreateDuplicateEmailCaseInsensitive(t *testing.T) {
	ctx := context.Background()
	repo := postgresrepo.NewAdminRepo(adminTestPool(t))
	input := service.CreateAdminInput{
		Email:        "duplicate@example.com",
		Name:         "First",
		PasswordHash: "$2a$10$hashed",
		Role:         service.AdminRoleAdmin,
	}
	_, err := repo.Create(ctx, input)
	require.NoError(t, err)

	input.Email = "DUPLICATE@example.com"
	_, err = repo.Create(ctx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrAdminDuplicateEmail)
}

func TestAdminRepo_GetMissingReturnsNil(t *testing.T) {
	repo := postgresrepo.NewAdminRepo(adminTestPool(t))

	byEmail, err := repo.GetByEmail(context.Background(), "missing@example.com")
	require.NoError(t, err)
	assert.Nil(t, byEmail)

	byID, err := repo.GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, byID)
}

func TestAdminRepo_GetByIDRejectsInvalidUUID(t *testing.T) {
	repo := postgresrepo.NewAdminRepo(adminTestPool(t))

	admin, err := repo.GetByID(context.Background(), "not-a-uuid")

	require.Error(t, err)
	assert.Nil(t, admin)
	assert.Contains(t, err.Error(), "invalid admin id")
}

func TestAdminRepo_ReturnsQueryErrors(t *testing.T) {
	ctx := context.Background()
	pool := adminTestPool(t)
	repo := postgresrepo.NewAdminRepo(pool)
	pool.Close()

	count, err := repo.Count(ctx)
	require.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Contains(t, err.Error(), "AdminRepo.Count")

	created, err := repo.Create(ctx, service.CreateAdminInput{
		Email: "admin@example.com", Name: "Admin", PasswordHash: "$2a$10$hashed", Role: service.AdminRoleAdmin,
	})
	require.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "AdminRepo.Create")

	byEmail, err := repo.GetByEmail(ctx, "admin@example.com")
	require.Error(t, err)
	assert.Nil(t, byEmail)
	assert.Contains(t, err.Error(), "AdminRepo.GetByEmail")

	byID, err := repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")
	require.Error(t, err)
	assert.Nil(t, byID)
	assert.Contains(t, err.Error(), "AdminRepo.GetByID")
}
