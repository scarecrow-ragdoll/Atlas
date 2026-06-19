// FILE: apps/api/internal/repository/redis/admin_session_store_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify Redis-backed admin session storage.
//   SCOPE: Session create/read/delete, expiry, HMAC-derived Redis keys, and unavailable Redis skip semantics; excludes HTTP cookies and GraphQL.
//   DEPENDS: apps/api/internal/repository/redis, apps/api/internal/testinfra, github.com/redis/go-redis/v9.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestAdminSessionStore_* - Real Redis coverage for admin sessions.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin session store coverage.
// END_CHANGE_SUMMARY

package redis_test

import (
	"context"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	redisrepo "monorepo-template/apps/api/internal/repository/redis"
	"monorepo-template/apps/api/internal/testinfra"
)

func adminRedisClient(t *testing.T) *redisrepo.Client {
	t.Helper()
	client, err := redisrepo.New(testinfra.RedisConfig(t), zap.NewNop())
	if err != nil && !testinfra.CoverageGateEnabled() {
		t.Skipf("redis integration service is unavailable: %v", err)
	}
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	cleanupAdminSessionKeys(t, client.RDB)
	return client
}

func cleanupAdminSessionKeys(t *testing.T, rdb *goredis.Client) {
	t.Helper()
	ctx := context.Background()
	iter := rdb.Scan(ctx, 0, "admin_session:*", 100).Iterator()
	for iter.Next(ctx) {
		require.NoError(t, rdb.Del(ctx, iter.Val()).Err())
	}
	require.NoError(t, iter.Err())
}

func TestAdminSessionStore_CreateReadDelete(t *testing.T) {
	ctx := context.Background()
	store := redisrepo.NewAdminSessionStore(adminRedisClient(t).RDB, []byte("test-key-secret"), time.Hour)

	sessionID, err := store.Create(ctx, "admin-1")
	require.NoError(t, err)
	require.NotEmpty(t, sessionID)

	adminID, err := store.Get(ctx, sessionID)
	require.NoError(t, err)
	assert.Equal(t, "admin-1", adminID)

	require.NoError(t, store.Delete(ctx, sessionID))
	adminID, err = store.Get(ctx, sessionID)
	require.NoError(t, err)
	assert.Empty(t, adminID)
}

func TestAdminSessionStore_DoesNotUseRawSessionIDAsKey(t *testing.T) {
	ctx := context.Background()
	client := adminRedisClient(t)
	store := redisrepo.NewAdminSessionStore(client.RDB, []byte("test-key-secret"), time.Hour)
	sessionID, err := store.Create(ctx, "admin-1")
	require.NoError(t, err)

	keys, err := client.RDB.Keys(ctx, "admin_session:*").Result()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	assert.NotContains(t, keys[0], sessionID)
}

func TestAdminSessionStore_Expires(t *testing.T) {
	ctx := context.Background()
	store := redisrepo.NewAdminSessionStore(adminRedisClient(t).RDB, []byte("test-key-secret"), time.Millisecond)
	sessionID, err := store.Create(ctx, "admin-1")
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		adminID, err := store.Get(ctx, sessionID)
		require.NoError(t, err)
		return adminID == ""
	}, 250*time.Millisecond, 10*time.Millisecond)
}
