// FILE: apps/api/internal/atlas/repository/postgres/nutrition_migration_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Integration tests for WAVE-05 nutrition migration smoke test: verifies migration applies without errors. Skipped when DB unavailable.
//   SCOPE: Migration smoke (up only) - TEST-W05-022. Requires test Postgres on port 17501.
//   DEPENDS: internal/testinfra for Postgres config, internal/repository/postgres for RunMigrations.
//   LINKS: V-M-API / EC-W05-006.
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

func TestWave05NutritionMigration_Smoke(t *testing.T) {
	dsn := testinfra.PostgresDSN()
	testinfra.RequireSafePostgresDSN(t, dsn)

	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("INTEGRATION_TESTS not set; skipping WAVE-05 migration smoke test")
	}

	err := postgresrepo.RunMigrations(dsn, zap.NewNop())
	require.NoError(t, err, "WAVE-05 migration should apply without error")
}