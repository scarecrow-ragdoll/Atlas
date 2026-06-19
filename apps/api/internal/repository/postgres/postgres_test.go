package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	postgresrepo "monorepo-template/apps/api/internal/repository/postgres"
	"monorepo-template/apps/api/internal/testinfra"
)

func TestNew_ConnectsAndPings(t *testing.T) {
	cfg := testinfra.PostgresConfig(t)
	db, err := postgresrepo.New(cfg, zap.NewNop())
	if err != nil && !testinfra.CoverageGateEnabled() {
		t.Skipf("postgres integration database is unavailable: %v", err)
	}
	require.NoError(t, err)
	t.Cleanup(db.Close)
	require.NoError(t, db.Ping())
}

func TestNew_ReturnsErrorForBadPort(t *testing.T) {
	cfg := testinfra.PostgresConfig(t)
	cfg.Port = 1

	_, err := postgresrepo.New(cfg, zap.NewNop())

	require.Error(t, err)
}

func TestNew_ReturnsErrorForMalformedDSN(t *testing.T) {
	cfg := testinfra.PostgresConfig(t)
	cfg.Host = "%"

	_, err := postgresrepo.New(cfg, zap.NewNop())

	require.Error(t, err)
}

func TestRunMigrations_ReturnsErrorForBadDSN(t *testing.T) {
	err := postgresrepo.RunMigrations("postgres://app:secret@localhost:1/monorepo_test?sslmode=disable", zap.NewNop())

	require.Error(t, err)
}
