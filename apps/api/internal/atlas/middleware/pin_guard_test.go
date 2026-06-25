// FILE: apps/api/internal/atlas/middleware/pin_guard_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify Atlas PIN guard middleware and context helper behavior at the HTTP boundary.
//   SCOPE: AtlasUserContext attachment, AtlasPinGuard allow/reject paths, context helpers; uses service and store fakes.
//   DEPENDS: internal/atlas/middleware, internal/atlas/service, internal/atlas/repository/redis, httptest.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   pin guard tests - Prove context helpers, user context middleware, and PIN guard allow/block paths.
// END_MODULE_MAP

package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/middleware"
	atlasRedis "monorepo-template/apps/api/internal/atlas/repository/redis"
	"monorepo-template/apps/api/internal/atlas/service"
	"monorepo-template/libs/go/logger"
)

func TestGetAtlasUserID_ReturnsEmptyWhenNotSet(t *testing.T) {
	assert.Equal(t, "", middleware.GetAtlasUserID(context.Background()))
}

func TestGetAtlasUserID_ReturnsAttachedValue(t *testing.T) {
	ctx := middleware.ContextWithAtlasUserID(context.Background(), "user-42")
	assert.Equal(t, "user-42", middleware.GetAtlasUserID(ctx))
}

func TestGetAtlasSessionToken_ReturnsEmptyWhenNotSet(t *testing.T) {
	assert.Equal(t, "", middleware.GetAtlasSessionToken(context.Background()))
}

func TestGetAtlasSessionToken_ReturnsAttachedValue(t *testing.T) {
	ctx := middleware.ContextWithAtlasSessionToken(context.Background(), "token-abc")
	assert.Equal(t, "token-abc", middleware.GetAtlasSessionToken(ctx))
}

func TestAtlasUserContext_AttachesUserID(t *testing.T) {
	bs := &fakeBootstrapService{userID: "bootstrapped-user"}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := middleware.GetAtlasUserID(r.Context())
		assert.Equal(t, "bootstrapped-user", uid)
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.AtlasUserContext(bs)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAtlasUserContext_Returns503WhenBootstrapFails(t *testing.T) {
	bs := &fakeBootstrapService{err: errors.New("db down")}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next should not be called")
	})

	handler := middleware.AtlasUserContext(bs)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Contains(t, rec.Body.String(), "atlas service unavailable")
}

func TestAtlasPinGuard_RejectsWhenNoUserInContext(t *testing.T) {
	handler := middleware.AtlasPinGuard(nil, nil, "session")(okHandler())

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAtlasPinGuard_AllowsWhenPinDisabled(t *testing.T) {
	ps := &fakePinService{pinEnabled: false}
	handler := middleware.AtlasPinGuard(ps, nil, "session")(okHandler())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAtlasPinGuard_RejectsWhenPinServiceErrors(t *testing.T) {
	ps := &fakePinService{err: errors.New("timeout")}
	handler := middleware.AtlasPinGuard(ps, nil, "session")(okHandler())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func TestAtlasPinGuard_RejectsWhenNoCookie(t *testing.T) {
	ps := &fakePinService{pinEnabled: true}
	handler := middleware.AtlasPinGuard(ps, nil, "session")(okHandler())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAtlasPinGuard_RejectsWithEmptyCookie(t *testing.T) {
	ps := &fakePinService{pinEnabled: true}
	handler := middleware.AtlasPinGuard(ps, nil, "session")(okHandler())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: ""})
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAtlasPinGuard_RejectsWhenSessionValidateErrors(t *testing.T) {
	ps := &fakePinService{pinEnabled: true}
	ss := &fakeSessionStore{validateErr: errors.New("redis down")}
	handler := middleware.AtlasPinGuard(ps, ss, "session")(okHandler())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "valid-token"})
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func TestAtlasPinGuard_RejectsWhenSessionInvalid(t *testing.T) {
	ps := &fakePinService{pinEnabled: true}
	ss := &fakeSessionStore{valid: false}
	handler := middleware.AtlasPinGuard(ps, ss, "session")(okHandler())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "bad-token"})
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAtlasPinGuard_RejectsWhenSessionUserMismatch(t *testing.T) {
	ps := &fakePinService{pinEnabled: true}
	ss := &fakeSessionStore{userID: "user-2", valid: true}
	handler := middleware.AtlasPinGuard(ps, ss, "session")(okHandler())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "other-user-token"})
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAtlasPinGuard_AllowsWhenPinEnabledAndValidSession(t *testing.T) {
	ps := &fakePinService{pinEnabled: true}
	ss := &fakeSessionStore{userID: "user-1", valid: true}
	handler := middleware.AtlasPinGuard(ps, ss, "session")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := middleware.GetAtlasSessionToken(r.Context())
		assert.Equal(t, "my-token", token)
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "my-token"})
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- fakes ---

type fakeBootstrapService struct {
	userID string
	err    error
}

func (f *fakeBootstrapService) EnsureDefaultUser(ctx context.Context) (string, error) {
	return f.userID, f.err
}

func (f *fakeBootstrapService) EnsureDefaultSettings(ctx context.Context, userID string) error {
	return f.err
}

func (f *fakeBootstrapService) EnsureDefaultUserProfile(ctx context.Context, userID string) error {
	return f.err
}

type fakePinService struct {
	pinEnabled bool
	err        error
}

func (f *fakePinService) IsEnabled(ctx context.Context, userID string) (bool, error) {
	return f.pinEnabled, f.err
}

func (f *fakePinService) Enable(ctx context.Context, userID, pin string) error {
	return f.err
}

func (f *fakePinService) Disable(ctx context.Context, userID, currentPin string) error {
	return f.err
}

func (f *fakePinService) Change(ctx context.Context, userID, currentPin, newPin string) error {
	return f.err
}

func (f *fakePinService) Verify(ctx context.Context, userID, pin string) (bool, error) {
	return false, f.err
}

type fakeSessionStore struct {
	userID       string
	valid        bool
	validateErr  error
	revokeCalled bool
}

func (f *fakeSessionStore) Create(ctx context.Context, userID string, idleTTL, absoluteTTL time.Duration) (string, error) {
	return "new-token", nil
}

func (f *fakeSessionStore) Validate(ctx context.Context, token string) (string, bool, error) {
	if f.validateErr != nil {
		return "", false, f.validateErr
	}
	return f.userID, f.valid, nil
}

func (f *fakeSessionStore) Revoke(ctx context.Context, token string) error {
	f.revokeCalled = true
	return nil
}

func (f *fakeSessionStore) RevokeAllByUser(ctx context.Context, userID string) error {
	return nil
}

// ensure fakeSessionStore implements PinSessionStore
var _ atlasRedis.PinSessionStore = (*fakeSessionStore)(nil)

// helpers

func okHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func withUserID(ctx context.Context, userID string) context.Context {
	ctx = logger.WithContext(ctx, zap.NewNop())
	return middleware.ContextWithAtlasUserID(ctx, userID)
}

// compile-time interface checks
var _ service.BootstrapService = (*fakeBootstrapService)(nil)
var _ service.PinService = (*fakePinService)(nil)
