// FILE: apps/api/internal/atlas/repository/redis/pin_attempt_store.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement PIN brute-force attempt tracking in Redis.
//   SCOPE: RegisterFailure, RegisterSuccess, IsLocked operations with configurable thresholds and lockout durations.
//   DEPENDS: github.com/redis/go-redis/v9.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Atlas PIN attempt store for WAVE-01.
// END_CHANGE_SUMMARY

package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type PinAttemptStore interface {
	RegisterFailure(ctx context.Context, key string) error
	RegisterSuccess(ctx context.Context, key string) error
	IsLocked(ctx context.Context, key string) (bool, time.Duration, error)
}

type pinAttemptStore struct {
	rdb               *goredis.Client
	maxFailures       int
	lockoutDuration   time.Duration
	escalatedDuration time.Duration
}

func NewPinAttemptStore(rdb *goredis.Client, maxFailures int, lockoutDuration, escalatedDuration time.Duration) PinAttemptStore {
	return &pinAttemptStore{
		rdb:               rdb,
		maxFailures:       maxFailures,
		lockoutDuration:   lockoutDuration,
		escalatedDuration: escalatedDuration,
	}
}

func (s *pinAttemptStore) attemptKey(identifier string) string {
	return "atlas:pin_attempt:" + identifier
}

func (s *pinAttemptStore) lockoutKey(identifier string) string {
	return "atlas:pin_lockout:" + identifier
}

func (s *pinAttemptStore) RegisterFailure(ctx context.Context, identifier string) error {
	attemptKey := s.attemptKey(identifier)
	count, err := s.rdb.Incr(ctx, attemptKey).Result()
	if err != nil {
		return fmt.Errorf("pin_attempt_store.RegisterFailure: incr: %w", err)
	}

	if count == 1 {
		s.rdb.Expire(ctx, attemptKey, s.lockoutDuration)
	}

	if int(count) >= s.maxFailures {
		lockoutKey := s.lockoutKey(identifier)
		existingTTL, err := s.rdb.TTL(ctx, lockoutKey).Result()
		if err != nil {
			return fmt.Errorf("pin_attempt_store.RegisterFailure: ttl check: %w", err)
		}

		duration := s.lockoutDuration
		if existingTTL > 0 {
			duration = s.escalatedDuration
		}

		if err := s.rdb.Set(ctx, lockoutKey, "1", duration).Err(); err != nil {
			return fmt.Errorf("pin_attempt_store.RegisterFailure: set lockout: %w", err)
		}
	}

	return nil
}

func (s *pinAttemptStore) RegisterSuccess(ctx context.Context, identifier string) error {
	attemptKey := s.attemptKey(identifier)
	lockoutKey := s.lockoutKey(identifier)

	pipe := s.rdb.Pipeline()
	pipe.Del(ctx, attemptKey)
	pipe.Del(ctx, lockoutKey)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("pin_attempt_store.RegisterSuccess: pipeline: %w", err)
	}

	return nil
}

func (s *pinAttemptStore) IsLocked(ctx context.Context, identifier string) (bool, time.Duration, error) {
	lockoutKey := s.lockoutKey(identifier)
	ttl, err := s.rdb.TTL(ctx, lockoutKey).Result()
	if err != nil {
		return false, 0, fmt.Errorf("pin_attempt_store.IsLocked: ttl: %w", err)
	}

	if ttl > 0 {
		return true, ttl, nil
	}

	return false, 0, nil
}