// FILE: apps/api/internal/service/admin_auth_internal_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify private admin auth error branches that require package-local seams.
//   SCOPE: Password hash failure propagation and nil public-admin mapping; excludes transport, PostgreSQL, and Redis behavior.
//   DEPENDS: apps/api/internal/service, bcrypt-compatible hash seam.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   testAdminRepo - Minimal package-local admin repository fake for hash-failure paths.
//   testAdminSessions - Minimal package-local admin session fake.
//   TestAdminAuthService_*HashFailure - Verifies service methods propagate hash failures.
//   TestPublicAdmin_Nil - Verifies nil admin mapping remains nil.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added package-local coverage for admin auth private error seams.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminAuthService_SeedInitialAdminReturnsHashFailure(t *testing.T) {
	restore := replaceBcryptGenerate(t, errors.New("bcrypt failed"))
	defer restore()
	svc := NewAdminAuthService(&testAdminRepo{}, &testAdminSessions{})

	created, err := svc.SeedInitialAdmin(context.Background(), InitialAdminInput{
		Email: "admin@example.com", Name: "Admin", Password: "StrongPassword123!",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "hash admin password")
	assert.False(t, created)
}

func TestAdminAuthService_CreateAdminReturnsHashFailure(t *testing.T) {
	restore := replaceBcryptGenerate(t, errors.New("bcrypt failed"))
	defer restore()
	svc := NewAdminAuthService(&testAdminRepo{}, &testAdminSessions{})

	admin, err := svc.CreateAdmin(context.Background(), &Admin{ID: "admin-1", IsActive: true}, NewAdminInput{
		Email: "new@example.com", Name: "New", Password: "StrongPassword123!",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "hash admin password")
	assert.Nil(t, admin)
}

func TestPublicAdmin_Nil(t *testing.T) {
	assert.Nil(t, publicAdmin(nil))
}

func replaceBcryptGenerate(t *testing.T, err error) func() {
	t.Helper()
	previous := bcryptGenerateFromPassword
	bcryptGenerateFromPassword = func(password []byte, cost int) ([]byte, error) {
		return nil, err
	}
	return func() {
		bcryptGenerateFromPassword = previous
	}
}

type testAdminRepo struct{}

func (r *testAdminRepo) Count(ctx context.Context) (int, error) {
	return 0, nil
}

func (r *testAdminRepo) Create(ctx context.Context, input CreateAdminInput) (*Admin, error) {
	return &Admin{ID: "admin-1", Email: input.Email, Name: input.Name, Role: input.Role, IsActive: true}, nil
}

func (r *testAdminRepo) GetByEmail(ctx context.Context, email string) (*Admin, error) {
	return nil, nil
}

func (r *testAdminRepo) GetByID(ctx context.Context, id string) (*Admin, error) {
	return nil, nil
}

type testAdminSessions struct{}

func (s *testAdminSessions) Create(ctx context.Context, adminID string) (string, error) {
	return "session-1", nil
}

func (s *testAdminSessions) Get(ctx context.Context, sessionID string) (string, error) {
	return "", nil
}

func (s *testAdminSessions) Delete(ctx context.Context, sessionID string) error {
	return nil
}
