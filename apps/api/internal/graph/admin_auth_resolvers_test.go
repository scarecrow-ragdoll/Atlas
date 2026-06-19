// FILE: apps/api/internal/graph/admin_auth_resolvers_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify GraphQL admin auth resolver behavior at the service and cookie bridge boundary.
//   SCOPE: me, loginAdmin, logoutAdmin, createAdmin, cookie set/clear, auth error mapping, and secret-redaction marker paths; excludes real HTTP routing, Redis, and PostgreSQL.
//   DEPENDS: apps/api/internal/graph, apps/api/internal/middleware, apps/api/internal/service.
//   LINKS: M-API / V-M-API / M-GRAPHQL-SCHEMA / V-M-GRAPHQL-SCHEMA.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   newAdminResolver - Builds resolver with fake admin repo/session boundaries.
//   TestMe_* - Verifies current-admin query behavior.
//   TestLoginAdmin_* - Verifies login result and session cookie behavior.
//   TestLogoutAdmin_* - Verifies session revocation and cookie clearing.
//   TestCreateAdmin_* - Verifies auth guard and admin creation mapping.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin auth GraphQL resolver coverage.
// END_CHANGE_SUMMARY

package graph

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"golang.org/x/crypto/bcrypt"

	"monorepo-template/apps/api/internal/graph/model"
	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/apps/api/internal/service"
	"monorepo-template/libs/go/logger"
)

func TestMe_ReturnsNilWithoutSession(t *testing.T) {
	resolver := &Resolver{}
	admin, err := resolver.Query().Me(context.Background())
	require.NoError(t, err)
	assert.Nil(t, admin)
}

func TestMe_ReturnsAdminFromPrincipal(t *testing.T) {
	resolver := &Resolver{}
	ctx := middleware.ContextWithAdminPrincipal(context.Background(), middleware.AdminPrincipal{
		ID: "admin-1", Email: "admin@example.com", Name: "Admin", Role: service.AdminRoleAdmin,
		CreatedAt: "2026-06-07T00:00:00Z", UpdatedAt: "2026-06-07T00:00:00Z",
	})

	admin, err := resolver.Query().Me(ctx)

	require.NoError(t, err)
	require.NotNil(t, admin)
	assert.Equal(t, "admin@example.com", admin.Email)
}

func TestLoginAdmin_SetsSessionCookie(t *testing.T) {
	repo := newFakeAdminRepoForGraph()
	hash, err := bcrypt.GenerateFromPassword([]byte("StrongPassword123!"), bcrypt.MinCost)
	require.NoError(t, err)
	repo.adminsByEmail["admin@example.com"] = &service.Admin{
		ID: "admin-1", Email: "admin@example.com", Name: "Admin", PasswordHash: string(hash), Role: service.AdminRoleAdmin, IsActive: true,
		CreatedAt: "2026-06-07T00:00:00Z", UpdatedAt: "2026-06-07T00:00:00Z",
	}
	resolver := &Resolver{AdminAuthService: service.NewAdminAuthService(repo, newFakeAdminSessionsForGraph())}
	rec := httptest.NewRecorder()
	ctx := middleware.ContextWithAdminCookieBridge(context.Background(), middleware.AdminCookieBridge{
		Response: rec,
		Config:   testAdminCookieConfig(),
	})

	result, err := resolver.Mutation().LoginAdmin(ctx, model.LoginAdminInput{Email: "admin@example.com", Password: "StrongPassword123!"})

	require.NoError(t, err)
	_, ok := result.(model.LoginAdminSuccess)
	require.True(t, ok)
	cookie := rec.Result().Cookies()[0]
	assert.Equal(t, "web_admin_session", cookie.Name)
	assert.True(t, cookie.HttpOnly)
	assert.Equal(t, "/graphql", cookie.Path)
}

func TestLoginAdmin_ReturnsAuthValidationAndUnexpectedErrors(t *testing.T) {
	repo := newFakeAdminRepoForGraph()
	hash, err := bcrypt.GenerateFromPassword([]byte("StrongPassword123!"), bcrypt.MinCost)
	require.NoError(t, err)
	repo.adminsByEmail["admin@example.com"] = &service.Admin{
		ID: "admin-1", Email: "admin@example.com", Name: "Admin", PasswordHash: string(hash), Role: service.AdminRoleAdmin, IsActive: true,
	}
	resolver := &Resolver{AdminAuthService: service.NewAdminAuthService(repo, newFakeAdminSessionsForGraph())}

	result, err := resolver.Mutation().LoginAdmin(context.Background(), model.LoginAdminInput{Email: "admin@example.com", Password: "wrong-password"})
	require.NoError(t, err)
	authErr, ok := result.(model.AuthError)
	require.True(t, ok)
	assert.Contains(t, authErr.Message, "invalid email or password")

	repo = newFakeAdminRepoForGraph()
	repo.getByEmailErr = fmtWrapped(service.ErrAdminValidation)
	resolver = &Resolver{AdminAuthService: service.NewAdminAuthService(repo, newFakeAdminSessionsForGraph())}
	result, err = resolver.Mutation().LoginAdmin(context.Background(), model.LoginAdminInput{Email: "", Password: ""})
	require.NoError(t, err)
	validation, ok := result.(model.ValidationError)
	require.True(t, ok)
	assert.Equal(t, "email", validation.Field)

	repo = newFakeAdminRepoForGraph()
	repo.getByEmailErr = errors.New("database down")
	resolver = &Resolver{AdminAuthService: service.NewAdminAuthService(repo, newFakeAdminSessionsForGraph())}
	result, err = resolver.Mutation().LoginAdmin(context.Background(), model.LoginAdminInput{Email: "admin@example.com", Password: "StrongPassword123!"})
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestLogoutAdmin_DeletesSessionAndClearsCookie(t *testing.T) {
	sessions := newFakeAdminSessionsForGraph()
	sessions.sessions["session-1"] = "admin-1"
	resolver := &Resolver{AdminAuthService: service.NewAdminAuthService(newFakeAdminRepoForGraph(), sessions)}
	rec := httptest.NewRecorder()
	ctx := middleware.ContextWithAdminCookieBridge(context.Background(), middleware.AdminCookieBridge{
		Response: rec,
		Config:   testAdminCookieConfig(),
	})
	ctx = middleware.ContextWithAdminSessionID(ctx, "session-1")

	result, err := resolver.Mutation().LogoutAdmin(ctx)

	require.NoError(t, err)
	success, ok := result.(model.LogoutAdminSuccess)
	require.True(t, ok)
	assert.True(t, success.Ok)
	assert.Equal(t, "session-1", sessions.deleted)
	cookie := rec.Result().Cookies()[0]
	assert.Equal(t, "web_admin_session", cookie.Name)
	assert.Equal(t, -1, cookie.MaxAge)
}

func TestLogoutAdmin_ReturnsServiceError(t *testing.T) {
	sessions := newFakeAdminSessionsForGraph()
	sessions.deleteErr = errors.New("redis down")
	resolver := &Resolver{AdminAuthService: service.NewAdminAuthService(newFakeAdminRepoForGraph(), sessions)}

	result, err := resolver.Mutation().LogoutAdmin(middleware.ContextWithAdminSessionID(context.Background(), "session-1"))

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestLogoutAdmin_SucceedsAndClearsCookieWithoutSession(t *testing.T) {
	sessions := newFakeAdminSessionsForGraph()
	resolver := &Resolver{AdminAuthService: service.NewAdminAuthService(newFakeAdminRepoForGraph(), sessions)}
	rec := httptest.NewRecorder()
	ctx := middleware.ContextWithAdminCookieBridge(context.Background(), middleware.AdminCookieBridge{
		Response: rec,
		Config:   testAdminCookieConfig(),
	})

	result, err := resolver.Mutation().LogoutAdmin(ctx)

	require.NoError(t, err)
	success, ok := result.(model.LogoutAdminSuccess)
	require.True(t, ok)
	assert.True(t, success.Ok)
	assert.Empty(t, sessions.deleted)
	cookie := rec.Result().Cookies()[0]
	assert.Equal(t, -1, cookie.MaxAge)
}

func TestCreateAdmin_ReturnsAuthErrorWithoutPrincipal(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	resolver := newAdminResolver(t)
	ctx := logger.WithContext(context.Background(), zap.New(core))
	result, err := resolver.Mutation().CreateAdmin(ctx, model.CreateAdminInput{
		Email: "new@example.com", Name: "New", Password: "StrongPassword123!",
	})
	require.NoError(t, err)
	authErr, ok := result.(model.AuthError)
	require.True(t, ok)
	assert.Contains(t, authErr.Message, "authentication required")
	joined := graphLogText(logs.All())
	assert.Contains(t, joined, "[AdminAuth][guard][BLOCK_AUTHORIZE_GRAPHQL]")
	assert.NotContains(t, joined, "new@example.com")
	assert.NotContains(t, joined, "StrongPassword123!")
	assert.NotContains(t, joined, "createAdmin")
}

func TestCreateAdmin_ReturnsSuccessWithPrincipal(t *testing.T) {
	resolver := newAdminResolver(t)
	ctx := middleware.ContextWithAdminPrincipal(context.Background(), middleware.AdminPrincipal{ID: "admin-1", Email: "admin@example.com", Name: "Admin", Role: "ADMIN"})
	result, err := resolver.Mutation().CreateAdmin(ctx, model.CreateAdminInput{
		Email: "new@example.com", Name: "New", Password: "StrongPassword123!",
	})
	require.NoError(t, err)
	success, ok := result.(model.CreateAdminSuccess)
	require.True(t, ok)
	assert.Equal(t, "new@example.com", success.Admin.Email)
}

func TestCreateAdmin_ReturnsDuplicateValidationAndUnexpectedErrors(t *testing.T) {
	repo := newFakeAdminRepoForGraph()
	repo.adminsByID["admin-1"] = &service.Admin{ID: "admin-1", Email: "admin@example.com", Role: service.AdminRoleAdmin, IsActive: true}
	repo.adminsByEmail["existing@example.com"] = &service.Admin{ID: "admin-existing", Email: "existing@example.com", IsActive: true}
	resolver := &Resolver{AdminAuthService: service.NewAdminAuthService(repo, newFakeAdminSessionsForGraph())}
	ctx := middleware.ContextWithAdminPrincipal(context.Background(), middleware.AdminPrincipal{ID: "admin-1", Email: "admin@example.com", Role: "ADMIN"})

	result, err := resolver.Mutation().CreateAdmin(ctx, model.CreateAdminInput{
		Email: "existing@example.com", Name: "Existing", Password: "StrongPassword123!",
	})
	require.NoError(t, err)
	validation, ok := result.(model.ValidationError)
	require.True(t, ok)
	assert.Equal(t, "email", validation.Field)

	result, err = resolver.Mutation().CreateAdmin(ctx, model.CreateAdminInput{
		Email: "new@example.com", Name: "New", Password: "short",
	})
	require.NoError(t, err)
	validation, ok = result.(model.ValidationError)
	require.True(t, ok)
	assert.Equal(t, "password", validation.Field)

	repo.createErr = errors.New("database down")
	result, err = resolver.Mutation().CreateAdmin(ctx, model.CreateAdminInput{
		Email: "other@example.com", Name: "Other", Password: "StrongPassword123!",
	})
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestMapAdmin_Nil(t *testing.T) {
	assert.Nil(t, mapAdmin(nil))
}

func newAdminResolver(t *testing.T) *Resolver {
	t.Helper()
	repo := newFakeAdminRepoForGraph()
	repo.adminsByID["admin-1"] = &service.Admin{ID: "admin-1", Email: "admin@example.com", Role: service.AdminRoleAdmin, IsActive: true}
	return &Resolver{AdminAuthService: service.NewAdminAuthService(repo, newFakeAdminSessionsForGraph())}
}

func testAdminCookieConfig() middleware.AdminCookieConfig {
	return middleware.AdminCookieConfig{Name: "web_admin_session", Path: "/graphql", MaxAge: 3600, Secure: false, SameSite: http.SameSiteLaxMode}
}

type fakeAdminRepoForGraph struct {
	adminsByID    map[string]*service.Admin
	adminsByEmail map[string]*service.Admin
	getByEmailErr error
	createErr     error
}

func newFakeAdminRepoForGraph() *fakeAdminRepoForGraph {
	return &fakeAdminRepoForGraph{adminsByID: map[string]*service.Admin{}, adminsByEmail: map[string]*service.Admin{}}
}

func (r *fakeAdminRepoForGraph) Count(ctx context.Context) (int, error) {
	return len(r.adminsByID), nil
}

func (r *fakeAdminRepoForGraph) Create(ctx context.Context, input service.CreateAdminInput) (*service.Admin, error) {
	if r.createErr != nil {
		return nil, r.createErr
	}
	if _, exists := r.adminsByEmail[input.Email]; exists {
		return nil, service.ErrAdminDuplicateEmail
	}
	admin := &service.Admin{
		ID: "admin-" + input.Email, Email: input.Email, Name: input.Name,
		PasswordHash: input.PasswordHash, Role: input.Role, IsActive: true,
		CreatedAt: "2026-06-07T00:00:00Z", UpdatedAt: "2026-06-07T00:00:00Z",
	}
	r.adminsByID[admin.ID] = admin
	r.adminsByEmail[admin.Email] = admin
	return admin, nil
}

func (r *fakeAdminRepoForGraph) GetByEmail(ctx context.Context, email string) (*service.Admin, error) {
	if r.getByEmailErr != nil {
		return nil, r.getByEmailErr
	}
	return r.adminsByEmail[email], nil
}

func (r *fakeAdminRepoForGraph) GetByID(ctx context.Context, id string) (*service.Admin, error) {
	return r.adminsByID[id], nil
}

type fakeAdminSessionsForGraph struct {
	sessions  map[string]string
	deleted   string
	deleteErr error
}

func newFakeAdminSessionsForGraph() *fakeAdminSessionsForGraph {
	return &fakeAdminSessionsForGraph{sessions: map[string]string{}}
}

func (s *fakeAdminSessionsForGraph) Create(ctx context.Context, adminID string) (string, error) {
	s.sessions["session-1"] = adminID
	return "session-1", nil
}

func (s *fakeAdminSessionsForGraph) Get(ctx context.Context, sessionID string) (string, error) {
	return s.sessions[sessionID], nil
}

func (s *fakeAdminSessionsForGraph) Delete(ctx context.Context, sessionID string) error {
	if s.deleteErr != nil {
		return s.deleteErr
	}
	s.deleted = sessionID
	delete(s.sessions, sessionID)
	return nil
}

func fmtWrapped(err error) error {
	return errors.Join(errors.New("wrapped"), err)
}

func graphLogText(entries []observer.LoggedEntry) string {
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
