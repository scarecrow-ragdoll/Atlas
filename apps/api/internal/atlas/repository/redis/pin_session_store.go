// FILE: apps/api/internal/atlas/repository/redis/pin_session_store.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement PIN session storage in Redis with sliding TTL and user session index.
//   SCOPE: Create, Validate, Revoke, RevokeAllByUser operations; uses HMAC-derived keys and SHA256 token hashing for opaque session tokens.
//   DEPENDS: crypto/rand, crypto/hmac, crypto/sha256, encoding/hex, github.com/redis/go-redis/v9.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Atlas PIN session store for WAVE-01.
// END_CHANGE_SUMMARY

package redis

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type PinSessionStore interface {
	Create(ctx context.Context, userID string, idleTTL, absoluteTTL time.Duration) (string, error)
	Validate(ctx context.Context, token string) (userID string, valid bool, err error)
	Revoke(ctx context.Context, token string) error
	RevokeAllByUser(ctx context.Context, userID string) error
}

type pinSessionPayload struct {
	UserID           string `json:"userID"`
	CreatedAt        string `json:"createdAt"`
	LastSeenAt       string `json:"lastSeenAt"`
	ExpiresAt        string `json:"expiresAt"`
	AbsoluteExpiresAt string `json:"absoluteExpiresAt"`
}

type pinSessionStore struct {
	rdb       *goredis.Client
	keySecret []byte
}

func NewPinSessionStore(rdb *goredis.Client, keySecret []byte) PinSessionStore {
	return &pinSessionStore{rdb: rdb, keySecret: keySecret}
}

func (s *pinSessionStore) Create(ctx context.Context, userID string, idleTTL, absoluteTTL time.Duration) (string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("pin_session_store.Create: generate token: %w", err)
	}
	token := hex.EncodeToString(raw[:])
	now := time.Now().UTC()

	payload := pinSessionPayload{
		UserID:            userID,
		CreatedAt:         now.Format(time.RFC3339),
		LastSeenAt:        now.Format(time.RFC3339),
		ExpiresAt:         now.Add(idleTTL).Format(time.RFC3339),
		AbsoluteExpiresAt: now.Add(absoluteTTL).Format(time.RFC3339),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("pin_session_store.Create: marshal payload: %w", err)
	}

	key := s.sessionKey(token)
	userKey := s.userSessionsKey(userID)

	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, key, data, idleTTL)
	pipe.SAdd(ctx, userKey, tokenHash(token))
	pipe.Expire(ctx, userKey, absoluteTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		return "", fmt.Errorf("pin_session_store.Create: pipeline exec: %w", err)
	}

	return token, nil
}

func (s *pinSessionStore) Validate(ctx context.Context, token string) (string, bool, error) {
	if token == "" {
		return "", false, nil
	}

	key := s.sessionKey(token)
	data, err := s.rdb.Get(ctx, key).Bytes()
	if errors.Is(err, goredis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("pin_session_store.Validate: redis get: %w", err)
	}

	var payload pinSessionPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", false, fmt.Errorf("pin_session_store.Validate: unmarshal: %w", err)
	}

	absExpiresAt, err := time.Parse(time.RFC3339, payload.AbsoluteExpiresAt)
	if err != nil {
		return "", false, fmt.Errorf("pin_session_store.Validate: parse absoluteExpiresAt: %w", err)
	}

	if time.Now().UTC().After(absExpiresAt) {
		s.rdb.Del(ctx, key)
		return "", false, nil
	}

	expiresAt, err := time.Parse(time.RFC3339, payload.ExpiresAt)
	if err != nil {
		return "", false, fmt.Errorf("pin_session_store.Validate: parse expiresAt: %w", err)
	}

	now := time.Now().UTC()
	remainingIdleTTL := expiresAt.Sub(now)
	if remainingIdleTTL <= 0 {
		s.rdb.Del(ctx, key)
		return "", false, nil
	}

	payload.LastSeenAt = now.Format(time.RFC3339)
	payload.ExpiresAt = now.Add(remainingIdleTTL).Format(time.RFC3339)

	updated, err := json.Marshal(payload)
	if err != nil {
		return "", false, fmt.Errorf("pin_session_store.Validate: marshal updated payload: %w", err)
	}

	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, key, updated, remainingIdleTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		return "", false, fmt.Errorf("pin_session_store.Validate: pipeline exec: %w", err)
	}

	return payload.UserID, true, nil
}

func (s *pinSessionStore) Revoke(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}

	key := s.sessionKey(token)
	data, err := s.rdb.Get(ctx, key).Bytes()
	if errors.Is(err, goredis.Nil) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("pin_session_store.Revoke: redis get: %w", err)
	}

	var payload pinSessionPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("pin_session_store.Revoke: unmarshal: %w", err)
	}

	userKey := s.userSessionsKey(payload.UserID)
	pipe := s.rdb.Pipeline()
	pipe.Del(ctx, key)
	pipe.SRem(ctx, userKey, tokenHash(token))
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("pin_session_store.Revoke: pipeline exec: %w", err)
	}

	return nil
}

func (s *pinSessionStore) RevokeAllByUser(ctx context.Context, userID string) error {
	userKey := s.userSessionsKey(userID)
	tokenHashes, err := s.rdb.SMembers(ctx, userKey).Result()
	if err != nil {
		return fmt.Errorf("pin_session_store.RevokeAllByUser: smembers: %w", err)
	}

	if len(tokenHashes) == 0 {
		return nil
	}

	pipe := s.rdb.Pipeline()
	for _, th := range tokenHashes {
		pipe.Del(ctx, "atlas:pin_session:"+th)
	}
	pipe.Del(ctx, userKey)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("pin_session_store.RevokeAllByUser: pipeline exec: %w", err)
	}

	return nil
}

func (s *pinSessionStore) sessionKey(token string) string {
	return "atlas:pin_session:" + tokenHash(token)
}

func (s *pinSessionStore) userSessionsKey(userID string) string {
	return "atlas:pin_user_sessions:" + userID
}

func tokenHash(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func keyedHash(token string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}