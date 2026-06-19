// FILE: apps/api/internal/middleware/admin_auth_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify admin auth cookie, context, and protected-operation helpers.
//   SCOPE: AdminPrincipal context storage, session cookie set/clear attributes, config mapping, GraphQL cookie bridge, session-cookie principal hydration, and protected guard outcomes; excludes Redis implementation and GraphQL resolver logic.
//   DEPENDS: apps/api/internal/appconfig, apps/api/internal/middleware, apps/api/internal/service, httptest.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestAdminCookie_* - Verifies session cookie attributes.
//   TestAdminPrincipal_* - Verifies context principal storage.
//   TestAdminCookieBridge_* - Verifies response-writer bridge for GraphQL cookie writes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added appconfig cookie mapping and GraphQL cookie bridge coverage.
// END_CHANGE_SUMMARY

package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"monorepo-template/apps/api/internal/appconfig"
	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/apps/api/internal/service"
	"monorepo-template/libs/go/logger"
)

func TestAdminPrincipal_ContextRoundTrip(t *testing.T) {
	principal := middleware.AdminPrincipal{ID: "admin-1", Email: "admin@example.com", Name: "Admin", Role: "ADMIN"}
	ctx := middleware.ContextWithAdminPrincipal(context.Background(), principal)

	found, ok := middleware.GetAdminPrincipal(ctx)

	require.True(t, ok)
	assert.Equal(t, principal, found)
}

func TestAdminCookie_SetAndClear(t *testing.T) {
	rec := httptest.NewRecorder()
	cfg := middleware.AdminCookieConfig{
		Name:     "web_admin_session",
		Path:     "/graphql",
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	middleware.SetAdminSessionCookie(rec, cfg, "session-id")
	setCookie := rec.Result().Cookies()[0]
	assert.Equal(t, "web_admin_session", setCookie.Name)
	assert.Equal(t, "session-id", setCookie.Value)
	assert.True(t, setCookie.HttpOnly)
	assert.Equal(t, "/graphql", setCookie.Path)
	assert.Equal(t, http.SameSiteLaxMode, setCookie.SameSite)

	clearRec := httptest.NewRecorder()
	middleware.ClearAdminSessionCookie(clearRec, cfg)
	clearCookie := clearRec.Result().Cookies()[0]
	assert.Equal(t, "web_admin_session", clearCookie.Name)
	assert.Empty(t, clearCookie.Value)
	assert.Equal(t, -1, clearCookie.MaxAge)
	assert.True(t, clearCookie.Expires.Before(time.Now()))
	assert.Equal(t, "/graphql", clearCookie.Path)
}

func TestAdminCookieConfigFromConfig_DerivesSecureSameSiteAndMaxAge(t *testing.T) {
	cfg := middleware.AdminCookieConfigFromConfig(appconfig.AdminSessionConfig{
		CookieName:   "web_admin_session",
		TTL:          7 * 24 * time.Hour,
		CookieSecure: "auto",
		SameSite:     "Strict",
	}, "production")

	assert.Equal(t, "web_admin_session", cfg.Name)
	assert.Equal(t, "/graphql", cfg.Path)
	assert.Equal(t, int((7 * 24 * time.Hour).Seconds()), cfg.MaxAge)
	assert.True(t, cfg.Secure)
	assert.Equal(t, http.SameSiteStrictMode, cfg.SameSite)

	devCfg := middleware.AdminCookieConfigFromConfig(appconfig.AdminSessionConfig{
		CookieName:   "web_admin_session",
		TTL:          time.Hour,
		CookieSecure: "false",
		SameSite:     "None",
	}, "development")

	assert.False(t, devCfg.Secure)
	assert.Equal(t, http.SameSiteNoneMode, devCfg.SameSite)

	laxCfg := middleware.AdminCookieConfigFromConfig(appconfig.AdminSessionConfig{
		CookieName:   "web_admin_session",
		TTL:          time.Hour,
		CookieSecure: "true",
		SameSite:     "Lax",
	}, "development")
	assert.True(t, laxCfg.Secure)
	assert.Equal(t, http.SameSiteLaxMode, laxCfg.SameSite)
}

func TestAdminCookieBridge_WritesSessionCookieFromGraphQLContext(t *testing.T) {
	cookieCfg := middleware.AdminCookieConfig{
		Name:     "web_admin_session",
		Path:     "/graphql",
		MaxAge:   3600,
		SameSite: http.SameSiteLaxMode,
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)

	handler := middleware.WithAdminCookieBridge(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.SetAdminSessionCookieFromContext(r.Context(), "session-1")
	}), cookieCfg)

	handler.ServeHTTP(rec, req)

	cookie := rec.Result().Cookies()[0]
	assert.Equal(t, "web_admin_session", cookie.Name)
	assert.Equal(t, "session-1", cookie.Value)
	assert.Equal(t, "/graphql", cookie.Path)
}

func TestAdminCookieBridge_ClearsSessionCookieFromGraphQLContext(t *testing.T) {
	cookieCfg := middleware.AdminCookieConfig{
		Name:     "web_admin_session",
		Path:     "/graphql",
		MaxAge:   3600,
		SameSite: http.SameSiteLaxMode,
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)

	handler := middleware.WithAdminCookieBridge(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.ClearAdminSessionCookieFromContext(r.Context())
	}), cookieCfg)

	handler.ServeHTTP(rec, req)

	cookie := rec.Result().Cookies()[0]
	assert.Equal(t, "web_admin_session", cookie.Name)
	assert.Equal(t, -1, cookie.MaxAge)
}

func TestAdminCookieBridge_NoopsWithoutBridge(t *testing.T) {
	ctx := context.Background()

	assert.NotPanics(t, func() {
		middleware.SetAdminSessionCookieFromContext(ctx, "session-1")
		middleware.ClearAdminSessionCookieFromContext(ctx)
	})
}

func TestAdminSessionMiddleware_SkipsWhenCookieMissing(t *testing.T) {
	resolver := &fakeAdminSessionResolver{}
	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	rec := httptest.NewRecorder()
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		_, ok := middleware.GetAdminPrincipal(r.Context())
		assert.False(t, ok)
		assert.Empty(t, middleware.AdminSessionIDFromContext(r.Context()))
	})

	middleware.AdminSessionMiddleware(resolver, "web_admin_session")(next).ServeHTTP(rec, req)

	assert.True(t, called)
	assert.False(t, resolver.called)
}

func TestAdminSessionMiddleware_LoadsPrincipalFromCookie(t *testing.T) {
	resolver := &fakeAdminSessionResolver{admin: &service.Admin{
		ID: "admin-1", Email: "admin@example.com", Name: "Admin", Role: service.AdminRoleAdmin, IsActive: true,
		CreatedAt: "2026-06-07T00:00:00Z", UpdatedAt: "2026-06-07T00:00:00Z",
	}}
	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.AddCookie(&http.Cookie{Name: "web_admin_session", Value: "session-1"})
	rec := httptest.NewRecorder()
	var principal middleware.AdminPrincipal
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		principal, ok = middleware.GetAdminPrincipal(r.Context())
		require.True(t, ok)
		assert.Equal(t, "session-1", middleware.AdminSessionIDFromContext(r.Context()))
	})

	middleware.AdminSessionMiddleware(resolver, "web_admin_session")(next).ServeHTTP(rec, req)

	assert.Equal(t, "admin@example.com", principal.Email)
	assert.Equal(t, "Admin", principal.Name)
}

func TestAdminSessionMiddleware_ReturnsInternalServerErrorOnLookupFailure(t *testing.T) {
	resolver := &fakeAdminSessionResolver{err: assert.AnError}
	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.AddCookie(&http.Cookie{Name: "web_admin_session", Value: "session-1"})
	rec := httptest.NewRecorder()

	middleware.AdminSessionMiddleware(resolver, "web_admin_session")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler must not run after session lookup failure")
	})).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "admin session lookup failed")
}

func TestAdminSessionMiddleware_DoesNotSetPrincipalForInactiveSession(t *testing.T) {
	resolver := &fakeAdminSessionResolver{admin: nil}
	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.AddCookie(&http.Cookie{Name: "web_admin_session", Value: "session-1"})
	rec := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := middleware.GetAdminPrincipal(r.Context())
		assert.False(t, ok)
		assert.Equal(t, "session-1", middleware.AdminSessionIDFromContext(r.Context()))
	})

	middleware.AdminSessionMiddleware(resolver, "web_admin_session")(next).ServeHTTP(rec, req)
}

func TestAdminSessionMiddleware_LogsMarkerWithoutSecrets(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	resolver := &fakeAdminSessionResolver{admin: &service.Admin{ID: "admin-1", Email: "admin@example.com", Role: service.AdminRoleAdmin, IsActive: true}}
	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req = req.WithContext(logger.WithContext(req.Context(), zap.New(core)))
	req.AddCookie(&http.Cookie{Name: "web_admin_session", Value: "raw-session-id"})
	rec := httptest.NewRecorder()

	middleware.AdminSessionMiddleware(resolver, "web_admin_session")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(rec, req)

	joined := logText(logs.All())
	assert.Contains(t, joined, "[AdminAuth][session][BLOCK_VALIDATE_SESSION]")
	assert.NotContains(t, joined, "raw-session-id")
	assert.NotContains(t, joined, "web_admin_session")
}

type fakeAdminSessionResolver struct {
	admin  *service.Admin
	err    error
	called bool
}

func (f *fakeAdminSessionResolver) CurrentAdmin(ctx context.Context, sessionID string) (*service.Admin, error) {
	f.called = true
	if f.err != nil {
		return nil, f.err
	}
	if sessionID != "session-1" {
		return nil, nil
	}
	return f.admin, nil
}

func logText(entries []observer.LoggedEntry) string {
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
