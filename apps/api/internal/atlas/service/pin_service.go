// FILE: apps/api/internal/atlas/service/pin_service.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement PIN management for the Atlas fitness tracker using Argon2id hashing.
//   SCOPE: Enable, Disable, Change, Verify, IsEnabled operations; PIN validation (4-20 digits); session revocation on disable/change.
//   DEPENDS: golang.org/x/crypto/argon2, apps/api/internal/atlas/repository/postgres.SettingsRepository, apps/api/internal/atlas/repository/redis.PinSessionStore.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Atlas PIN service with Argon2id for WAVE-01.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/crypto/argon2"

	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	atlasRedis "monorepo-template/apps/api/internal/atlas/repository/redis"
)

var (
	ErrPinWrongPin          = errors.New("wrong PIN")
	ErrPinAlreadyEnabled    = errors.New("PIN is already enabled")
	ErrPinAlreadyDisabled   = errors.New("PIN is already disabled")
	ErrPinTooShort          = errors.New("PIN too short")
	ErrPinTooLong           = errors.New("PIN too long")
	ErrPinInternal          = errors.New("internal PIN error")

	digitsOnly = regexp.MustCompile(`^\d+$`)

	// cryptoRandReader allows replacing for deterministic tests.
	cryptoRandReader io.Reader = rand.Reader
)

// SetCryptoRandReader replaces the random source for deterministic testing.
func SetCryptoRandReader(r io.Reader) {
	cryptoRandReader = r
}

// HashPinForTest hashes a PIN using a fixed salt for deterministic testing.
func HashPinForTest(pin string) string {
	oldReader := cryptoRandReader
	defer func() { cryptoRandReader = oldReader }()
	cryptoRandReader = strings.NewReader(strings.Repeat("a", 16))
	s := &pinService{params: Argon2Params{Memory: 64, Iterations: 1, Parallelism: 1, KeyLength: 32}}
	return s.hashPin(pin)
}

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	KeyLength   uint32
}

type PinService interface {
	Enable(ctx context.Context, userID string, pin string) error
	Disable(ctx context.Context, userID string, currentPin string) error
	Change(ctx context.Context, userID string, currentPin, newPin string) error
	Verify(ctx context.Context, userID string, pin string) (bool, error)
	IsEnabled(ctx context.Context, userID string) (bool, error)
}

type pinService struct {
	settingsRepo atlasPostgres.SettingsRepository
	sessionStore atlasRedis.PinSessionStore
	params       Argon2Params
	minLength    int
	maxLength    int
}

func NewPinService(
	settingsRepo atlasPostgres.SettingsRepository,
	sessionStore atlasRedis.PinSessionStore,
	params Argon2Params,
	minLength, maxLength int,
) PinService {
	return &pinService{
		settingsRepo: settingsRepo,
		sessionStore: sessionStore,
		params:       params,
		minLength:    minLength,
		maxLength:    maxLength,
	}
}

func (s *pinService) Enable(ctx context.Context, userID string, pin string) error {
	if err := s.validatePin(pin); err != nil {
		return err
	}

	record, err := s.settingsRepo.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("pin_service.Enable: %w", err)
	}

	if record.PinEnabled {
		return ErrPinAlreadyEnabled
	}

	hash := s.hashPin(pin)
	if err := s.settingsRepo.UpdatePinState(ctx, userID, true, &hash); err != nil {
		return fmt.Errorf("pin_service.Enable: %w", err)
	}

	return nil
}

func (s *pinService) Disable(ctx context.Context, userID string, currentPin string) error {
	record, err := s.settingsRepo.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("pin_service.Disable: %w", err)
	}

	if !record.PinEnabled {
		return ErrPinAlreadyDisabled
	}

	if record.PinHash == nil {
		return ErrPinInternal
	}

	match, err := s.verifyHash(*record.PinHash, currentPin)
	if err != nil {
		return fmt.Errorf("pin_service.Disable: %w", err)
	}
	if !match {
		return ErrPinWrongPin
	}

	if err := s.settingsRepo.UpdatePinState(ctx, userID, false, nil); err != nil {
		return fmt.Errorf("pin_service.Disable: %w", err)
	}

	if err := s.sessionStore.RevokeAllByUser(ctx, userID); err != nil {
		return fmt.Errorf("pin_service.Disable: revoke sessions: %w", err)
	}

	return nil
}

func (s *pinService) Change(ctx context.Context, userID string, currentPin, newPin string) error {
	if err := s.validatePin(newPin); err != nil {
		return err
	}

	record, err := s.settingsRepo.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("pin_service.Change: %w", err)
	}

	if !record.PinEnabled {
		return ErrPinAlreadyDisabled
	}

	if record.PinHash == nil {
		return ErrPinInternal
	}

	match, err := s.verifyHash(*record.PinHash, currentPin)
	if err != nil {
		return fmt.Errorf("pin_service.Change: %w", err)
	}
	if !match {
		return ErrPinWrongPin
	}

	newHash := s.hashPin(newPin)
	if err := s.settingsRepo.UpdatePinState(ctx, userID, true, &newHash); err != nil {
		return fmt.Errorf("pin_service.Change: %w", err)
	}

	if err := s.sessionStore.RevokeAllByUser(ctx, userID); err != nil {
		return fmt.Errorf("pin_service.Change: revoke sessions: %w", err)
	}

	return nil
}

func (s *pinService) Verify(ctx context.Context, userID string, pin string) (bool, error) {
	record, err := s.settingsRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("pin_service.Verify: %w", err)
	}

	if record.PinHash == nil {
		return false, nil
	}

	return s.verifyHash(*record.PinHash, pin)
}

func (s *pinService) IsEnabled(ctx context.Context, userID string) (bool, error) {
	record, err := s.settingsRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("pin_service.IsEnabled: %w", err)
	}
	return record.PinEnabled, nil
}

func (s *pinService) validatePin(pin string) error {
	if len(pin) < s.minLength {
		return ErrPinTooShort
	}
	if len(pin) > s.maxLength {
		return ErrPinTooLong
	}
	if !digitsOnly.MatchString(pin) {
		return errors.New("PIN must contain only digits")
	}
	return nil
}

func (s *pinService) hashPin(pin string) string {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(cryptoRandReader, salt); err != nil {
		panic("pin_service: failed to generate salt: " + err.Error())
	}
	hash := argon2.IDKey([]byte(pin), salt, s.params.Iterations, s.params.Memory, s.params.Parallelism, s.params.KeyLength)
	return hex.EncodeToString(hash) + "$" + hex.EncodeToString(salt)
}

func (s *pinService) verifyHash(encodedHash, pin string) (bool, error) {
	parts := splitHash(encodedHash)
	if len(parts) != 2 {
		return false, ErrPinInternal
	}

	salt, err := hex.DecodeString(parts[1])
	if err != nil {
		return false, ErrPinInternal
	}

	expectedHash, err := hex.DecodeString(parts[0])
	if err != nil {
		return false, ErrPinInternal
	}

	computedHash := argon2.IDKey([]byte(pin), salt, s.params.Iterations, s.params.Memory, s.params.Parallelism, s.params.KeyLength)

	if len(expectedHash) != len(computedHash) {
		return false, nil
	}

	return subtle.ConstantTimeCompare(expectedHash, computedHash) == 1, nil
}

func splitHash(encoded string) []string {
	result := make([]string, 0, 2)
	current := ""
	for i := 0; i < len(encoded); i++ {
		if encoded[i] == '$' {
			result = append(result, current)
			current = ""
		} else {
			current += string(encoded[i])
		}
	}
	result = append(result, current)
	return result
}