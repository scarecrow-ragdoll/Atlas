// FILE: apps/api/internal/atlas/middleware/auth_separation_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify strict route-level separation between web-admin GraphQL auth (session cookie from AdminSessionMiddleware) and Atlas PIN auth (session cookie from AtlasPinGuard). No cookie from one zone should unlock the other.
//   SCOPE: Chi router integration tests with mock services for admin session, pin service, bootstrap, and pin session store; excludes real GraphQL resolvers and persistence.
//   DEPENDS: apps/api/internal/middleware, apps/api/internal/atlas/middleware, apps/api/internal/atlas/service, apps/api/internal/atlas/repository/redis, github.com/go-chi/chi/v5, httptest.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestAtlasAuthSeparation_* - Validates admin vs atlas auth cookie isolation across route groups.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added credentialed CORS preflight coverage for Atlas browser route groups.
//   LAST_CHANGE: 1.0.0 - Added auth separation integration tests for admin vs Atlas route groups.
// END_CHANGE_SUMMARY

package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	atlasMiddleware "monorepo-template/apps/api/internal/atlas/middleware"
	atlasRedis "monorepo-template/apps/api/internal/atlas/repository/redis"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/apps/api/internal/service"
)

// ---- mock types ----

type mockAdminSessionResolver struct {
	admin *service.Admin
	err   error
}

func (m *mockAdminSessionResolver) CurrentAdmin(_ context.Context, sessionID string) (*service.Admin, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.admin == nil || sessionID == "" {
		return nil, nil
	}
	return m.admin, nil
}

type mockPinService struct {
	enabled   bool
	enabledFn func() bool
}

func (m *mockPinService) Enable(_ context.Context, _ string, _ string) error  { return nil }
func (m *mockPinService) Disable(_ context.Context, _ string, _ string) error { return nil }
func (m *mockPinService) Change(_ context.Context, _ string, _, _ string) error {
	return nil
}
func (m *mockPinService) Verify(_ context.Context, _ string, _ string) (bool, error) {
	return false, nil
}
func (m *mockPinService) IsEnabled(_ context.Context, _ string) (bool, error) {
	if m.enabledFn != nil {
		return m.enabledFn(), nil
	}
	return m.enabled, nil
}

type mockBootstrapService struct {
	userID string
	err    error
}

func (m *mockBootstrapService) EnsureDefaultUser(_ context.Context) (string, error) {
	return m.userID, m.err
}
func (m *mockBootstrapService) EnsureDefaultSettings(_ context.Context, _ string) error {
	return m.err
}
func (m *mockBootstrapService) EnsureDefaultUserProfile(_ context.Context, _ string) error {
	return m.err
}

type mockPinSessionStore struct {
	validUserID string
	valid       bool
	err         error
}

func (m *mockPinSessionStore) Create(_ context.Context, _ string, _, _ time.Duration) (string, error) {
	return "atlas-session-token", nil
}
func (m *mockPinSessionStore) Validate(_ context.Context, token string) (string, bool, error) {
	if m.err != nil {
		return "", false, m.err
	}
	return m.validUserID, m.valid, nil
}
func (m *mockPinSessionStore) Revoke(_ context.Context, _ string) error { return nil }
func (m *mockPinSessionStore) RevokeAllByUser(_ context.Context, _ string) error {
	return nil
}

// ---- helpers ----

const testUserID = "00000000-0000-0000-0000-000000000001"

func newAdminCookie(name, value string) *http.Cookie {
	return &http.Cookie{Name: name, Value: value}
}

func newAtlasCookie(name, value string) *http.Cookie {
	return &http.Cookie{Name: name, Value: value}
}

// buildTestRouter creates a chi router with the same logical route groups as main.go,
// backed by mock services so auth separation can be tested without persistence.
func buildTestRouter(
	t *testing.T,
	adminResolver middleware.AdminSessionResolver,
	pinService atlasService.PinService,
	bootstrapSvc atlasService.BootstrapService,
	pinSessionStore atlasRedis.PinSessionStore,
) *chi.Mux {
	t.Helper()
	r := chi.NewRouter()
	adminCORS := middleware.CORS(middleware.CORSConfig{
		AllowedOrigins:   []string{"http://localhost:3100"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	r.Group(func(admin chi.Router) {
		admin.Use(adminCORS)
		admin.Use(middleware.AdminOriginGuard([]string{"http://localhost:3100"}))
		admin.Use(middleware.AdminSessionMiddleware(adminResolver, "web_admin_session"))
		admin.Handle("/graphql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := middleware.GetAdminPrincipal(r.Context())
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"errors":[{"message":"unauthorized","extensions":{"code":"UNAUTHENTICATED"}}]}`))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"__typename":"Query"}}`))
		}))
	})

	r.Group(func(atlasAuth chi.Router) {
		atlasAuth.Use(adminCORS)
		atlasAuth.Use(atlasMiddleware.AtlasUserContext(bootstrapSvc))
		atlasAuth.Options("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		atlasAuth.Post("/api/v1/auth/pin/unlock", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
	})

	r.Group(func(atlas chi.Router) {
		atlas.Use(adminCORS)
		atlas.Use(middleware.AdminOriginGuard([]string{"http://localhost:3100"}))
		atlas.Use(atlasMiddleware.AtlasUserContext(bootstrapSvc))
		atlas.Use(atlasMiddleware.AtlasPinGuard(pinService, pinSessionStore, "atlas_pin_session"))
		atlas.Options("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		atlas.Handle("/graphql/atlas", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"__typename":"AtlasQuery"}}`))
		}))
		atlas.Post("/api/ai-export/generate", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		atlas.Delete("/api/v1/media/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
	})

	return r
}

// ---- tests ----

func TestAtlasAuthSeparation_AdminGraphQL_NoAuth_Returns401(t *testing.T) {
	router := buildTestRouter(t,
		&mockAdminSessionResolver{},
		&mockPinService{},
		&mockBootstrapService{userID: testUserID},
		&mockPinSessionStore{},
	)

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}

func TestAtlasAuthSeparation_AtlasGraphQL_PINDisabled_Returns200(t *testing.T) {
	router := buildTestRouter(t,
		&mockAdminSessionResolver{},
		&mockPinService{enabled: false},
		&mockBootstrapService{userID: testUserID},
		&mockPinSessionStore{},
	)

	req := httptest.NewRequest(http.MethodPost, "/graphql/atlas", nil)
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "AtlasQuery")
}

func TestAtlasAuthSeparation_AtlasAiExportPreflight_AllowsWebAdminOrigin(t *testing.T) {
	router := buildTestRouter(t,
		&mockAdminSessionResolver{},
		&mockPinService{enabled: true},
		&mockBootstrapService{userID: testUserID},
		&mockPinSessionStore{},
	)

	req := httptest.NewRequest(http.MethodOptions, "/api/ai-export/generate", nil)
	req.Header.Set("Origin", "http://localhost:3100")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "http://localhost:3100", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), http.MethodPost)
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
}

func TestAtlasAuthSeparation_AtlasDeletePreflight_AllowsWebAdminOrigin(t *testing.T) {
	router := buildTestRouter(t,
		&mockAdminSessionResolver{},
		&mockPinService{enabled: true},
		&mockBootstrapService{userID: testUserID},
		&mockPinSessionStore{},
	)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/media/media-1", nil)
	req.Header.Set("Origin", "http://localhost:3100")
	req.Header.Set("Access-Control-Request-Method", http.MethodDelete)
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "http://localhost:3100", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), http.MethodDelete)
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
}

func TestAtlasAuthSeparation_AtlasGraphQL_PINEnabled_NoCookie_Returns401(t *testing.T) {
	router := buildTestRouter(t,
		&mockAdminSessionResolver{},
		&mockPinService{enabled: true},
		&mockBootstrapService{userID: testUserID},
		&mockPinSessionStore{},
	)

	req := httptest.NewRequest(http.MethodPost, "/graphql/atlas", nil)
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}

func TestAtlasAuthSeparation_AdminCookie_DoesNotWorkOnAtlasGraphQL(t *testing.T) {
	router := buildTestRouter(t,
		&mockAdminSessionResolver{admin: &service.Admin{
			ID: "admin-1", Email: "admin@example.com", Name: "Admin",
			Role: service.AdminRoleAdmin, IsActive: true,
		}},
		&mockPinService{enabled: true},
		&mockBootstrapService{userID: testUserID},
		&mockPinSessionStore{},
	)

	req := httptest.NewRequest(http.MethodPost, "/graphql/atlas", nil)
	req.AddCookie(newAdminCookie("web_admin_session", "valid-admin-session"))
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code,
		"admin cookie must NOT unlock the Atlas PIN-guarded route")
}

func TestAtlasAuthSeparation_AtlasCookie_DoesNotWorkOnAdminGraphQL(t *testing.T) {
	router := buildTestRouter(t,
		&mockAdminSessionResolver{},
		&mockPinService{enabled: true},
		&mockBootstrapService{userID: testUserID},
		&mockPinSessionStore{valid: true, validUserID: testUserID},
	)

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.AddCookie(newAtlasCookie("atlas_pin_session", "valid-atlas-token"))
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code,
		"atlas PIN session cookie must NOT authenticate the admin GraphQL route")
}

func TestAtlasAuthSeparation_AdminCookie_WithAdminSession_Succeeds(t *testing.T) {
	router := buildTestRouter(t,
		&mockAdminSessionResolver{admin: &service.Admin{
			ID: "admin-1", Email: "admin@example.com", Name: "Admin",
			Role: service.AdminRoleAdmin, IsActive: true,
		}},
		&mockPinService{},
		&mockBootstrapService{userID: testUserID},
		&mockPinSessionStore{},
	)

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.AddCookie(newAdminCookie("web_admin_session", "valid-admin-session"))
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Query")
}

func TestAtlasAuthSeparation_AtlasPINEnabled_WithValidCookie_Succeeds(t *testing.T) {
	router := buildTestRouter(t,
		&mockAdminSessionResolver{},
		&mockPinService{enabled: true},
		&mockBootstrapService{userID: testUserID},
		&mockPinSessionStore{valid: true, validUserID: testUserID},
	)

	req := httptest.NewRequest(http.MethodPost, "/graphql/atlas", nil)
	req.AddCookie(newAtlasCookie("atlas_pin_session", "valid-atlas-token"))
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "AtlasQuery")
}

// Test that a request with _both_ cookies still respects route-group boundaries
// — admin cookie unused on atlas guarded routes, atlas cookie unused on admin routes.
func TestAtlasAuthSeparation_AdminCookie_IgnoredByAtlasGuard(t *testing.T) {
	pinSvc := &mockPinService{enabled: true}
	sessionStore := &mockPinSessionStore{valid: true, validUserID: testUserID}
	router := buildTestRouter(t,
		&mockAdminSessionResolver{admin: &service.Admin{
			ID: "admin-1", Email: "admin@example.com", Name: "Admin",
			Role: service.AdminRoleAdmin, IsActive: true,
		}},
		pinSvc,
		&mockBootstrapService{userID: testUserID},
		sessionStore,
	)

	// Prove Atlas route still validates PIN guard even when admin cookie present
	req := httptest.NewRequest(http.MethodPost, "/graphql/atlas", nil)
	req.AddCookie(newAdminCookie("web_admin_session", "valid-admin-session"))
	req.AddCookie(newAtlasCookie("atlas_pin_session", "valid-atlas-token"))
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code,
		"admin cookie present should not affect Atlas PIN guard when valid atlas cookie also present")
}
