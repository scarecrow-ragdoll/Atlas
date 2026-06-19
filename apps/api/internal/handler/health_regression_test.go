// FILE: apps/api/internal/handler/health_regression_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify that the existing health and readiness endpoints remain unchanged. These are regression tests that mirror the existing coverage in health_test.go with explicit assertions on response structure, status codes, and header Content-Type.
//   SCOPE: Public /healthz and /readyz handlers; Atlas /api/v1/healthz and /api/v1/readyz handlers; does not test PIN auth, GraphQL, or admin routes.
//   DEPENDS: apps/api/internal/handler, apps/api/internal/repository/postgres, apps/api/internal/repository/redis, httptest.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestExistingHealthz_Returns200 - Regression: /healthz responds 200 with {"status":"ok"}.
//   TestExistingHealthz_ContentType - Regression: /healthz sets application/json.
//   TestExistingReadyz_WithMock - Regression: /readyz responds 200 when checkers pass.
//   TestExistingAtlasHealthz_Returns200 - Regression: /api/v1/healthz responds 200.
//   TestExistingAtlasReadyz_WithMock - Regression: /api/v1/readyz responds 200/503 with mock checkers.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added health endpoint regression tests.
// END_CHANGE_SUMMARY

package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"monorepo-template/apps/api/internal/handler"
)

// ---- mock checker ----

type mockHealthChecker struct {
	pingErr error
}

func (m *mockHealthChecker) Ping() error {
	return m.pingErr
}

// ---- /healthz regression ----

func TestExistingHealthz_Returns200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.Healthz()(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func TestExistingHealthz_ContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.Healthz()(rec, req)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestExistingHealthz_GETMethodOnly(t *testing.T) {
	// The handler does not reject based on method; verify the behavior for POST too.
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.Healthz()(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

// ---- /readyz regression ----

func TestExistingReadyz_WithMock(t *testing.T) {
	checker := &mockHealthChecker{}
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.Readyz(checker)(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func TestExistingReadyz_MultipleCheckersAllHealthy(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.Readyz(
		&mockHealthChecker{},
		&mockHealthChecker{},
	)(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestExistingReadyz_OneCheckerFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.Readyz(
		&mockHealthChecker{},
		&mockHealthChecker{pingErr: errors.New("db down")},
	)(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Contains(t, rec.Body.String(), "db down")
}

func TestExistingReadyz_ContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.Readyz(&mockHealthChecker{})(rec, req)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

// ---- /api/v1/healthz regression ----

func TestExistingAtlasHealthz_Returns200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
	rec := httptest.NewRecorder()

	handler.AtlasHealthz()(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func TestExistingAtlasHealthz_ContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
	rec := httptest.NewRecorder()

	handler.AtlasHealthz()(rec, req)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

// ---- /api/v1/readyz regression ----

type mockAtlasChecker struct {
	pingErr error
}

func (m *mockAtlasChecker) Ping() error {
	return m.pingErr
}

func TestExistingAtlasReadyz_WithMock(t *testing.T) {
	checker := &mockAtlasChecker{}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/readyz", nil)
	rec := httptest.NewRecorder()

	handler.AtlasReadyz(checker)(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func TestExistingAtlasReadyz_Unhealthy(t *testing.T) {
	checker := &mockAtlasChecker{pingErr: errors.New("redis unreachable")}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/readyz", nil)
	rec := httptest.NewRecorder()

	handler.AtlasReadyz(checker)(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Contains(t, rec.Body.String(), "redis unreachable")
}

func TestExistingAtlasReadyz_MultipleCheckers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/readyz", nil)
	rec := httptest.NewRecorder()

	handler.AtlasReadyz(
		&mockAtlasChecker{},
		&mockAtlasChecker{},
	)(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestExistingAtlasReadyz_SecondCheckerFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/readyz", nil)
	rec := httptest.NewRecorder()

	handler.AtlasReadyz(
		&mockAtlasChecker{},
		&mockAtlasChecker{pingErr: errors.New("db timeout")},
	)(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Contains(t, rec.Body.String(), "db timeout")
}

func TestExistingAtlasReadyz_ContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/readyz", nil)
	rec := httptest.NewRecorder()

	handler.AtlasReadyz(&mockAtlasChecker{})(rec, req)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}