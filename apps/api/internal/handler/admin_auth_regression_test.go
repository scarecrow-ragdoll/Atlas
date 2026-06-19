// FILE: apps/api/internal/handler/admin_auth_regression_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify that existing web-admin auth endpoints and middleware integration continue to work — valid admin cookie + origin → admin principal injected; invalid or missing cookie → unauthenticated.
//   SCOPE: Chi router integration tests with mock admin session resolver; no real GraphQL resolvers or persistence.
//   DEPENDS: apps/api/internal/middleware, apps/api/internal/service, github.com/go-chi/chi/v5, httptest.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestAdminHealth_* - Validates admin route group middleware and handler behavior.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin auth regression tests for middleware stack.
// END_CHANGE_SUMMARY

package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/apps/api/internal/service"
)

// ---- mock admin session resolver ----

type regressionAdminResolver struct {
	admin *service.Admin
	err   error
}

func (r *regressionAdminResolver) CurrentAdmin(_ context.Context, sessionID string) (*service.Admin, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.admin == nil || sessionID == "" {
		return nil, nil
	}
	return r.admin, nil
}

// ---- helpers ----

const testAdminID = "regression-admin-1"

func newAdminCookie(name, value string) *http.Cookie {
	return &http.Cookie{Name: name, Value: value}
}

func buildAdminRouter(resolver middleware.AdminSessionResolver) *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(admin chi.Router) {
		admin.Use(middleware.AdminOriginGuard([]string{"http://localhost:3100"}))
		admin.Use(middleware.AdminSessionMiddleware(resolver, "web_admin_session"))
		admin.Handle("/graphql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := middleware.GetAdminPrincipal(r.Context())
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"errors":[{"message":"unauthorized","extensions":{"code":"UNAUTHENTICATED"}}]}`))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"currentAdmin":{"id":"` + principal.ID + `","email":"` + principal.Email + `"}}}`))
		}))
	})
	return r
}

// ---- tests ----

func TestAdminHealth_Returns200(t *testing.T) {
	resolver := &regressionAdminResolver{
		admin: &service.Admin{
			ID: testAdminID, Email: "admin@example.com", Name: "Admin",
			Role: service.AdminRoleAdmin, IsActive: true,
		},
	}
	router := buildAdminRouter(resolver)

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("Origin", "http://localhost:3100")
	req.AddCookie(newAdminCookie("web_admin_session", "valid-session-id"))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), testAdminID)
	assert.Contains(t, rec.Body.String(), "admin@example.com")
}

func TestAdminHealth_MissingCookie_Returns401(t *testing.T) {
	resolver := &regressionAdminResolver{}
	router := buildAdminRouter(resolver)

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}

func TestAdminHealth_InvalidOrigin_Returns403(t *testing.T) {
	resolver := &regressionAdminResolver{}
	router := buildAdminRouter(resolver)

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("Origin", "http://malicious-origin.com")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAdminHealth_SessionLookupFailure_Returns500(t *testing.T) {
	resolver := &regressionAdminResolver{err: assert.AnError}
	router := buildAdminRouter(resolver)

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("Origin", "http://localhost:3100")
	req.AddCookie(newAdminCookie("web_admin_session", "some-session"))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAdminHealth_InactiveAdmin_Returns401(t *testing.T) {
	resolver := &regressionAdminResolver{
		admin: &service.Admin{
			ID: testAdminID, Email: "inactive@example.com", Name: "Inactive",
			Role: service.AdminRoleAdmin, IsActive: false,
		},
	}
	router := buildAdminRouter(resolver)

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("Origin", "http://localhost:3100")
	req.AddCookie(newAdminCookie("web_admin_session", "valid-session-id"))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code,
		"inactive admin must not be granted access")
}

func TestAdminHealth_GETWithoutOrigin_Allowed(t *testing.T) {
	resolver := &regressionAdminResolver{}
	router := buildAdminRouter(resolver)

	req := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Safe methods pass through origin guard; session middleware sees no cookie → passes through;
	// handler finds no principal → unauthenticated
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAdminHealth_PlaygroundPath_NotMounted(t *testing.T) {
	// Verify the admin playground is NOT mounted in the route set (it only appears
	// in main.go when Server.Env != "production").
	resolver := &regressionAdminResolver{}
	router := buildAdminRouter(resolver)

	req := httptest.NewRequest(http.MethodGet, "/playground", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code,
		"playground should not be mounted in a non-production-like test router")
}