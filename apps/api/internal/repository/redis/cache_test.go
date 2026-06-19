package redis_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	redisrepo "monorepo-template/apps/api/internal/repository/redis"
	"monorepo-template/apps/api/internal/testinfra"
	"monorepo-template/libs/go/config"
)

func TestNew_ConnectsPingsAndCloses(t *testing.T) {
	client, err := redisrepo.New(testinfra.RedisConfig(t), zap.NewNop())
	if err != nil && !testinfra.CoverageGateEnabled() {
		t.Skipf("redis integration service is unavailable: %v", err)
	}
	require.NoError(t, err)
	require.NoError(t, client.Ping())
	require.NoError(t, client.Close())
}

func TestNew_ReturnsErrorForBadPort(t *testing.T) {
	_, err := redisrepo.New(config.RedisConfig{Host: "localhost", Port: 1}, zap.NewNop())

	require.Error(t, err)
}
