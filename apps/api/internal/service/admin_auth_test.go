// FILE: apps/api/internal/service/admin_auth_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify transport-neutral admin auth service behavior.
//   SCOPE: Bootstrap seed, password validation, login, current admin, create admin, logout, duplicate mapping, inactive admin denial, session revocation, and secret-redaction marker paths; excludes PostgreSQL, Redis, HTTP cookies, and GraphQL.
//   DEPENDS: internal/service, bcrypt, testify.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   fakeAdminRepo - Admin repository test double.
//   fakeAdminSessions - Admin session store test double.
//   TestAdminAuthService_* - Service behavior coverage.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin auth service coverage.
// END_CHANGE_SUMMARY

package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"golang.org/x/crypto/bcrypt"

	"monorepo-template/apps/api/internal/service"
	"monorepo-template/libs/go/logger"
)

func TestAdminAuthService_SeedCreatesOnlyWhenEmpty(t *testing.T) {
	repo := newFakeAdminRepo()
	svc := service.NewAdminAuthService(repo, newFakeAdminSessions())

	created, err := svc.SeedInitialAdmin(context.Background(), service.InitialAdminInput{
		Email:    "Admin@Example.COM",
		Name:     "Template Admin",
		Password: "StrongPassword123!",
	})

	require.NoError(t, err)
	assert.True(t, created)
	assert.Equal(t, 1, repo.count)
	assert.Equal(t, "admin@example.com", repo.created[0].Email)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(repo.created[0].PasswordHash), []byte("StrongPassword123!")))
}

func TestAdminAuthService_SeedNoopsWhenAnyAdminExists(t *testing.T) {
	repo := newFakeAdminRepo()
	repo.count = 1
	svc := service.NewAdminAuthService(repo, newFakeAdminSessions())

	created, err := svc.SeedInitialAdmin(context.Background(), service.InitialAdminInput{
		Email:    "second@example.com",
		Name:     "Second",
		Password: "StrongPassword123!",
	})

	require.NoError(t, err)
	assert.False(t, created)
	assert.Empty(t, repo.created)
}

func TestAdminAuthService_SeedReturnsCountAndCreateErrors(t *testing.T) {
	countErr := errors.New("count failed")
	repo := newFakeAdminRepo()
	repo.countErr = countErr
	svc := service.NewAdminAuthService(repo, newFakeAdminSessions())

	created, err := svc.SeedInitialAdmin(context.Background(), service.InitialAdminInput{
		Email: "admin@example.com", Name: "Admin", Password: "StrongPassword123!",
	})
	require.ErrorIs(t, err, countErr)
	assert.False(t, created)

	createErr := errors.New("create failed")
	repo = newFakeAdminRepo()
	repo.createErr = createErr
	svc = service.NewAdminAuthService(repo, newFakeAdminSessions())

	created, err = svc.SeedInitialAdmin(context.Background(), service.InitialAdminInput{
		Email: "admin@example.com", Name: "Admin", Password: "StrongPassword123!",
	})
	require.ErrorIs(t, err, createErr)
	assert.False(t, created)
}

func TestAdminAuthService_SeedRejectsInvalidInput(t *testing.T) {
	tests := []service.InitialAdminInput{
		{Name: "Admin", Password: "StrongPassword123!"},
		{Email: "admin@example.com", Password: "StrongPassword123!"},
		{Email: "admin@example.com", Name: "Admin", Password: "short"},
	}
	for _, input := range tests {
		t.Run(input.Email+input.Name+input.Password, func(t *testing.T) {
			svc := service.NewAdminAuthService(newFakeAdminRepo(), newFakeAdminSessions())

			created, err := svc.SeedInitialAdmin(context.Background(), input)

			require.ErrorIs(t, err, service.ErrAdminValidation)
			assert.False(t, created)
		})
	}
}

func TestAdminAuthService_LoginCreatesSession(t *testing.T) {
	repo := newFakeAdminRepo()
	hash, err := bcrypt.GenerateFromPassword([]byte("StrongPassword123!"), bcrypt.MinCost)
	require.NoError(t, err)
	repo.adminsByEmail["admin@example.com"] = &service.Admin{
		ID: "admin-1", Email: "admin@example.com", Name: "Admin", PasswordHash: string(hash), Role: service.AdminRoleAdmin, IsActive: true,
	}
	sessions := newFakeAdminSessions()
	svc := service.NewAdminAuthService(repo, sessions)

	result, err := svc.Login(context.Background(), service.LoginAdminInput{Email: "ADMIN@example.com", Password: "StrongPassword123!"})

	require.NoError(t, err)
	assert.Equal(t, "session-1", result.SessionID)
	assert.Equal(t, "admin-1", sessions.createdFor)
	assert.Equal(t, "admin@example.com", result.Admin.Email)
}

func TestAdminAuthService_LoginReturnsRepositoryAndSessionErrors(t *testing.T) {
	repoErr := errors.New("repo failed")
	repo := newFakeAdminRepo()
	repo.getByEmailErr = repoErr
	svc := service.NewAdminAuthService(repo, newFakeAdminSessions())

	result, err := svc.Login(context.Background(), service.LoginAdminInput{Email: "admin@example.com", Password: "StrongPassword123!"})
	require.ErrorIs(t, err, repoErr)
	assert.Nil(t, result)

	hash, err := bcrypt.GenerateFromPassword([]byte("StrongPassword123!"), bcrypt.MinCost)
	require.NoError(t, err)
	repo = newFakeAdminRepo()
	repo.adminsByEmail["admin@example.com"] = &service.Admin{
		ID: "admin-1", Email: "admin@example.com", PasswordHash: string(hash), Role: service.AdminRoleAdmin, IsActive: true,
	}
	sessionErr := errors.New("session failed")
	sessions := newFakeAdminSessions()
	sessions.createErr = sessionErr
	svc = service.NewAdminAuthService(repo, sessions)

	result, err = svc.Login(context.Background(), service.LoginAdminInput{Email: "admin@example.com", Password: "StrongPassword123!"})
	require.ErrorIs(t, err, sessionErr)
	assert.Nil(t, result)
}

func TestAdminAuthService_LoginRejectsWrongPasswordAndInactiveAdmin(t *testing.T) {
	repo := newFakeAdminRepo()
	hash, err := bcrypt.GenerateFromPassword([]byte("StrongPassword123!"), bcrypt.MinCost)
	require.NoError(t, err)
	repo.adminsByEmail["admin@example.com"] = &service.Admin{
		ID: "admin-1", Email: "admin@example.com", PasswordHash: string(hash), Role: service.AdminRoleAdmin, IsActive: true,
	}
	svc := service.NewAdminAuthService(repo, newFakeAdminSessions())

	_, err = svc.Login(context.Background(), service.LoginAdminInput{Email: "admin@example.com", Password: "wrong-password"})
	require.ErrorIs(t, err, service.ErrAdminAuth)

	repo.adminsByEmail["admin@example.com"].IsActive = false
	_, err = svc.Login(context.Background(), service.LoginAdminInput{Email: "admin@example.com", Password: "StrongPassword123!"})
	require.ErrorIs(t, err, service.ErrAdminAuth)
}

func TestAdminAuthService_LoginRejectsMissingAdmin(t *testing.T) {
	svc := service.NewAdminAuthService(newFakeAdminRepo(), newFakeAdminSessions())

	result, err := svc.Login(context.Background(), service.LoginAdminInput{Email: "missing@example.com", Password: "StrongPassword123!"})

	require.ErrorIs(t, err, service.ErrAdminAuth)
	assert.Nil(t, result)
}

func TestAdminAuthService_CurrentAdmin(t *testing.T) {
	repo := newFakeAdminRepo()
	repo.adminsByID["admin-1"] = &service.Admin{ID: "admin-1", Email: "admin@example.com", Role: service.AdminRoleAdmin, IsActive: true}
	sessions := newFakeAdminSessions()
	sessions.sessions["session-1"] = "admin-1"
	svc := service.NewAdminAuthService(repo, sessions)

	admin, err := svc.CurrentAdmin(context.Background(), "session-1")
	require.NoError(t, err)
	require.NotNil(t, admin)
	assert.Equal(t, "admin@example.com", admin.Email)

	admin, err = svc.CurrentAdmin(context.Background(), "missing-session")
	require.NoError(t, err)
	assert.Nil(t, admin)
}

func TestAdminAuthService_CurrentAdminReturnsErrors(t *testing.T) {
	sessionErr := errors.New("session failed")
	sessions := newFakeAdminSessions()
	sessions.getErr = sessionErr
	svc := service.NewAdminAuthService(newFakeAdminRepo(), sessions)

	admin, err := svc.CurrentAdmin(context.Background(), "session-1")
	require.ErrorIs(t, err, sessionErr)
	assert.Nil(t, admin)

	repoErr := errors.New("repo failed")
	repo := newFakeAdminRepo()
	repo.getByIDErr = repoErr
	sessions = newFakeAdminSessions()
	sessions.sessions["session-1"] = "admin-1"
	svc = service.NewAdminAuthService(repo, sessions)

	admin, err = svc.CurrentAdmin(context.Background(), "session-1")
	require.ErrorIs(t, err, repoErr)
	assert.Nil(t, admin)
}

func TestAdminAuthService_CurrentAdminReturnsNilForMissingAdmin(t *testing.T) {
	sessions := newFakeAdminSessions()
	sessions.sessions["session-1"] = "admin-1"
	svc := service.NewAdminAuthService(newFakeAdminRepo(), sessions)

	admin, err := svc.CurrentAdmin(context.Background(), "session-1")

	require.NoError(t, err)
	assert.Nil(t, admin)
}

func TestAdminAuthService_CurrentAdminReturnsNilForInactiveAdmin(t *testing.T) {
	repo := newFakeAdminRepo()
	repo.adminsByID["admin-1"] = &service.Admin{ID: "admin-1", Email: "admin@example.com", Role: service.AdminRoleAdmin, IsActive: false}
	sessions := newFakeAdminSessions()
	sessions.sessions["session-1"] = "admin-1"
	svc := service.NewAdminAuthService(repo, sessions)

	admin, err := svc.CurrentAdmin(context.Background(), "session-1")

	require.NoError(t, err)
	assert.Nil(t, admin)
}

func TestAdminAuthService_CreateAdminRequiresActor(t *testing.T) {
	svc := service.NewAdminAuthService(newFakeAdminRepo(), newFakeAdminSessions())

	_, err := svc.CreateAdmin(context.Background(), nil, service.NewAdminInput{
		Email: "new@example.com", Name: "New", Password: "StrongPassword123!",
	})

	require.ErrorIs(t, err, service.ErrAdminAuth)
}

func TestAdminAuthService_CreateAdminRejectsInactiveActor(t *testing.T) {
	actor := &service.Admin{ID: "admin-1", IsActive: false}
	svc := service.NewAdminAuthService(newFakeAdminRepo(), newFakeAdminSessions())

	_, err := svc.CreateAdmin(context.Background(), actor, service.NewAdminInput{
		Email: "new@example.com", Name: "New", Password: "StrongPassword123!",
	})

	require.ErrorIs(t, err, service.ErrAdminAuth)
}

func TestAdminAuthService_CreateAdminHashesPassword(t *testing.T) {
	repo := newFakeAdminRepo()
	actor := &service.Admin{ID: "admin-1", Email: "admin@example.com", Role: service.AdminRoleAdmin, IsActive: true}
	svc := service.NewAdminAuthService(repo, newFakeAdminSessions())

	created, err := svc.CreateAdmin(context.Background(), actor, service.NewAdminInput{
		Email: "New@Example.COM", Name: "New Admin", Password: "StrongPassword123!",
	})

	require.NoError(t, err)
	assert.Equal(t, "new@example.com", created.Email)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(repo.created[0].PasswordHash), []byte("StrongPassword123!")))
}

func TestAdminAuthService_CreateAdminReturnsValidationDuplicateAndRepositoryErrors(t *testing.T) {
	actor := &service.Admin{ID: "admin-1", IsActive: true}
	svc := service.NewAdminAuthService(newFakeAdminRepo(), newFakeAdminSessions())

	_, err := svc.CreateAdmin(context.Background(), actor, service.NewAdminInput{
		Email: "", Name: "New", Password: "StrongPassword123!",
	})
	require.ErrorIs(t, err, service.ErrAdminValidation)

	repo := newFakeAdminRepo()
	repo.adminsByEmail["new@example.com"] = &service.Admin{ID: "admin-2", Email: "new@example.com", IsActive: true}
	svc = service.NewAdminAuthService(repo, newFakeAdminSessions())
	_, err = svc.CreateAdmin(context.Background(), actor, service.NewAdminInput{
		Email: "new@example.com", Name: "New", Password: "StrongPassword123!",
	})
	require.ErrorIs(t, err, service.ErrAdminDuplicateEmail)

	repoErr := errors.New("repo failed")
	repo = newFakeAdminRepo()
	repo.createErr = repoErr
	svc = service.NewAdminAuthService(repo, newFakeAdminSessions())
	_, err = svc.CreateAdmin(context.Background(), actor, service.NewAdminInput{
		Email: "new@example.com", Name: "New", Password: "StrongPassword123!",
	})
	require.ErrorIs(t, err, repoErr)
}

func TestAdminAuthService_LogoutDeletesSession(t *testing.T) {
	sessions := newFakeAdminSessions()
	svc := service.NewAdminAuthService(newFakeAdminRepo(), sessions)

	require.NoError(t, svc.Logout(context.Background(), "session-1"))
	assert.Equal(t, "session-1", sessions.deleted)
}

func TestAdminAuthService_LogoutReturnsSessionError(t *testing.T) {
	sessionErr := errors.New("delete failed")
	sessions := newFakeAdminSessions()
	sessions.deleteErr = sessionErr
	svc := service.NewAdminAuthService(newFakeAdminRepo(), sessions)

	err := svc.Logout(context.Background(), "session-1")

	require.ErrorIs(t, err, sessionErr)
}

func TestAdminAuthService_LogsMarkersWithoutSecrets(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	ctx := logger.WithContext(context.Background(), zap.New(core))
	repo := newFakeAdminRepo()
	hash, err := bcrypt.GenerateFromPassword([]byte("StrongPassword123!"), bcrypt.MinCost)
	require.NoError(t, err)
	repo.adminsByEmail["admin@example.com"] = &service.Admin{
		ID: "admin-1", Email: "admin@example.com", PasswordHash: string(hash), Role: service.AdminRoleAdmin, IsActive: true,
	}
	sessions := newFakeAdminSessions()
	svc := service.NewAdminAuthService(repo, sessions)

	_, _ = svc.Login(ctx, service.LoginAdminInput{Email: "admin@example.com", Password: "StrongPassword123!"})
	_ = svc.Logout(ctx, "session-1")

	joined := adminLogText(logs.All())
	assert.Contains(t, joined, "[AdminAuth][login][BLOCK_VERIFY_CREDENTIALS]")
	assert.Contains(t, joined, "[AdminAuth][session][BLOCK_VALIDATE_SESSION]")
	assert.Contains(t, joined, "[AdminAuth][logout][BLOCK_REVOKE_SESSION]")
	assert.NotContains(t, joined, "admin@example.com")
	assert.NotContains(t, joined, "StrongPassword123!")
	assert.NotContains(t, joined, string(hash))
	assert.NotContains(t, joined, "session-1")
}

func adminLogText(entries []observer.LoggedEntry) string {
	var out strings.Builder
	for _, entry := range entries {
		out.WriteString(entry.Message)
		out.WriteString("\n")
		for _, field := range entry.Context {
			out.WriteString(field.Key)
			out.WriteString("=")
			out.WriteString(field.String)
			out.WriteString("\n")
		}
	}
	return out.String()
}

type fakeAdminRepo struct {
	count         int
	adminsByID    map[string]*service.Admin
	adminsByEmail map[string]*service.Admin
	created       []service.CreateAdminInput
	countErr      error
	createErr     error
	getByEmailErr error
	getByIDErr    error
}

func newFakeAdminRepo() *fakeAdminRepo {
	return &fakeAdminRepo{adminsByID: map[string]*service.Admin{}, adminsByEmail: map[string]*service.Admin{}}
}

func (r *fakeAdminRepo) Count(ctx context.Context) (int, error) {
	if r.countErr != nil {
		return 0, r.countErr
	}
	return r.count, nil
}

func (r *fakeAdminRepo) Create(ctx context.Context, input service.CreateAdminInput) (*service.Admin, error) {
	if r.createErr != nil {
		return nil, r.createErr
	}
	if _, exists := r.adminsByEmail[input.Email]; exists {
		return nil, service.ErrAdminDuplicateEmail
	}
	r.created = append(r.created, input)
	r.count++
	admin := &service.Admin{ID: "admin-" + input.Email, Email: input.Email, Name: input.Name, PasswordHash: input.PasswordHash, Role: input.Role, IsActive: true}
	r.adminsByEmail[input.Email] = admin
	r.adminsByID[admin.ID] = admin
	return admin, nil
}

func (r *fakeAdminRepo) GetByEmail(ctx context.Context, email string) (*service.Admin, error) {
	if r.getByEmailErr != nil {
		return nil, r.getByEmailErr
	}
	return r.adminsByEmail[email], nil
}

func (r *fakeAdminRepo) GetByID(ctx context.Context, id string) (*service.Admin, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}
	return r.adminsByID[id], nil
}

type fakeAdminSessions struct {
	sessions   map[string]string
	createdFor string
	deleted    string
	createErr  error
	getErr     error
	deleteErr  error
}

func newFakeAdminSessions() *fakeAdminSessions {
	return &fakeAdminSessions{sessions: map[string]string{}}
}

func (s *fakeAdminSessions) Create(ctx context.Context, adminID string) (string, error) {
	if s.createErr != nil {
		return "", s.createErr
	}
	s.createdFor = adminID
	s.sessions["session-1"] = adminID
	return "session-1", nil
}

func (s *fakeAdminSessions) Get(ctx context.Context, sessionID string) (string, error) {
	if s.getErr != nil {
		return "", s.getErr
	}
	return s.sessions[sessionID], nil
}

func (s *fakeAdminSessions) Delete(ctx context.Context, sessionID string) error {
	if s.deleteErr != nil && !errors.Is(s.deleteErr, service.ErrAdminNotFound) {
		return s.deleteErr
	}
	s.deleted = sessionID
	delete(s.sessions, sessionID)
	return nil
}
