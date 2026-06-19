// FILE: apps/api/internal/repository/redis/admin_session_store.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Store web-admin opaque sessions in Redis using HMAC-derived keys.
//   SCOPE: Session id generation, Redis key derivation, create/read/delete, and TTL; excludes cookies, GraphQL, and admin identity lookup.
//   DEPENDS: crypto/rand, crypto/hmac, crypto/sha256, encoding/base64, github.com/redis/go-redis/v9.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   NewAdminSessionStore - Constructs a Redis-backed admin session store.
//   AdminSessionStore.Create - Creates an opaque browser session id and Redis entry.
//   AdminSessionStore.Get - Resolves a session id to an admin id.
//   AdminSessionStore.Delete - Revokes a session id.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Redis admin session store.
// END_CHANGE_SUMMARY

package redis

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type AdminSessionStore struct {
	rdb       *goredis.Client
	keySecret []byte
	ttl       time.Duration
}

var randomSessionID = generateRandomSessionID
var readRandom = rand.Read

func NewAdminSessionStore(rdb *goredis.Client, keySecret []byte, ttl time.Duration) *AdminSessionStore {
	return &AdminSessionStore{rdb: rdb, keySecret: keySecret, ttl: ttl}
}

// START_CONTRACT: Create
//
//	PURPOSE: Create an opaque admin browser session and persist the admin id behind an HMAC-derived Redis key.
//	INPUTS: { ctx: context.Context - request context, adminID: string - authenticated admin id }
//	OUTPUTS: { string - opaque session id for cookie storage, error - random generation or Redis failure }
//	SIDE_EFFECTS: Writes one Redis key with the configured TTL.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: Create
func (s *AdminSessionStore) Create(ctx context.Context, adminID string) (string, error) {
	sessionID, err := randomSessionID()
	if err != nil {
		return "", fmt.Errorf("AdminSessionStore.Create: generate session id: %w", err)
	}
	if err := s.rdb.Set(ctx, s.key(sessionID), adminID, s.ttl).Err(); err != nil {
		return "", fmt.Errorf("AdminSessionStore.Create: redis set: %w", err)
	}
	return sessionID, nil
}

// START_CONTRACT: Get
//
//	PURPOSE: Resolve an opaque session id into an admin id.
//	INPUTS: { ctx: context.Context - request context, sessionID: string - opaque session id from cookie }
//	OUTPUTS: { string - admin id or empty when absent/expired, error - Redis failure }
//	SIDE_EFFECTS: Reads Redis through an HMAC-derived key.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: Get
func (s *AdminSessionStore) Get(ctx context.Context, sessionID string) (string, error) {
	if sessionID == "" {
		return "", nil
	}
	adminID, err := s.rdb.Get(ctx, s.key(sessionID)).Result()
	if errors.Is(err, goredis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("AdminSessionStore.Get: redis get: %w", err)
	}
	return adminID, nil
}

// START_CONTRACT: Delete
//
//	PURPOSE: Revoke one opaque admin session id.
//	INPUTS: { ctx: context.Context - request context, sessionID: string - opaque session id from cookie }
//	OUTPUTS: { error - Redis failure }
//	SIDE_EFFECTS: Deletes one Redis key derived from the session id.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: Delete
func (s *AdminSessionStore) Delete(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}
	if err := s.rdb.Del(ctx, s.key(sessionID)).Err(); err != nil {
		return fmt.Errorf("AdminSessionStore.Delete: redis del: %w", err)
	}
	return nil
}

func (s *AdminSessionStore) key(sessionID string) string {
	mac := hmac.New(sha256.New, s.keySecret)
	_, _ = mac.Write([]byte(sessionID))
	return "admin_session:" + hex.EncodeToString(mac.Sum(nil))
}

func generateRandomSessionID() (string, error) {
	var raw [32]byte
	if _, err := readRandom(raw[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}
