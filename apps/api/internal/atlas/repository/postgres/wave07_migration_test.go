// FILE: apps/api/internal/atlas/repository/postgres/wave07_migration_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Integration tests for WAVE-07 migration smoke test: verifies migration applies without errors. Skipped when DB unavailable.
//   SCOPE: Migration smoke (up only). Requires test Postgres on port 17501.
//   DEPENDS: internal/testinfra for Postgres config, internal/repository/postgres for RunMigrations.
//   LINKS: V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT

package postgres_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	postgresrepo "monorepo-template/apps/api/internal/repository/postgres"
	"monorepo-template/apps/api/internal/testinfra"
)

func TestWave07Migration_ApplyCleanly(t *testing.T) {
	dsn := testinfra.PostgresDSN()
	testinfra.RequireSafePostgresDSN(t, dsn)

	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("INTEGRATION_TESTS not set; skipping WAVE-07 migration smoke test")
	}

	err := postgresrepo.RunMigrations(dsn, zap.NewNop())
	require.NoError(t, err, "WAVE-07 migration should apply without error")
}