// FILE: apps/api/internal/repository/redis/admin_session_store_internal_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify private Redis admin session error branches.
//   SCOPE: Session id generation failures, closed Redis client errors, empty-id idempotency, and random session id formatting; excludes Docker-backed Redis success behavior.
//   DEPENDS: apps/api/internal/repository/redis, github.com/redis/go-redis/v9.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestAdminSessionStore_*Error - Verifies failure-path wrapping and empty input behavior.
//   TestGenerateRandomSessionID - Verifies generated opaque ids are URL-safe and non-empty.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added package-local coverage for admin session private error seams.
// END_CHANGE_SUMMARY

package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminSessionStore_CreateReturnsRandomGenerationError(t *testing.T) {
	previous := randomSessionID
	randomSessionID = func() (string, error) {
		return "", errors.New("random failed")
	}
	defer func() { randomSessionID = previous }()
	store := NewAdminSessionStore(goredis.NewClient(&goredis.Options{Addr: "127.0.0.1:1"}), []byte("secret"), time.Hour)

	sessionID, err := store.Create(context.Background(), "admin-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "generate session id")
	assert.Empty(t, sessionID)
}

func TestAdminSessionStore_ReturnsRedisErrors(t *testing.T) {
	ctx := context.Background()
	rdb := goredis.NewClient(&goredis.Options{Addr: "127.0.0.1:1"})
	require.NoError(t, rdb.Close())
	store := NewAdminSessionStore(rdb, []byte("secret"), time.Hour)

	sessionID, err := store.Create(ctx, "admin-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "redis set")
	assert.Empty(t, sessionID)

	adminID, err := store.Get(ctx, "session-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "redis get")
	assert.Empty(t, adminID)

	err = store.Delete(ctx, "session-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "redis del")
}

func TestAdminSessionStore_EmptySessionIDIsIdempotent(t *testing.T) {
	store := NewAdminSessionStore(goredis.NewClient(&goredis.Options{Addr: "127.0.0.1:1"}), []byte("secret"), time.Hour)

	adminID, err := store.Get(context.Background(), "")
	require.NoError(t, err)
	assert.Empty(t, adminID)

	require.NoError(t, store.Delete(context.Background(), ""))
}

func TestGenerateRandomSessionID(t *testing.T) {
	sessionID, err := generateRandomSessionID()

	require.NoError(t, err)
	assert.NotEmpty(t, sessionID)
	assert.NotContains(t, sessionID, "+")
	assert.NotContains(t, sessionID, "/")
	assert.NotContains(t, sessionID, "=")
}

func TestGenerateRandomSessionIDReturnsReadError(t *testing.T) {
	previous := readRandom
	readRandom = func(b []byte) (int, error) {
		return 0, errors.New("entropy unavailable")
	}
	defer func() { readRandom = previous }()

	sessionID, err := generateRandomSessionID()

	require.Error(t, err)
	assert.Empty(t, sessionID)
}
