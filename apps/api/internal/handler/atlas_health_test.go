// FILE: apps/api/internal/handler/atlas_health_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify Atlas-specific health and readiness HTTP handlers.
//   SCOPE: /api/v1/healthz returns 200, /api/v1/readyz returns 200 when all checkers pass and 503 when any checker fails.
//   DEPENDS: internal/handler, httptest.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   atlas health tests - Prove healthz always-ok and readyz dependency-aware 200/503 behavior.
// END_MODULE_MAP

package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"monorepo-template/apps/api/internal/handler"
)

func TestAtlasHealthz_ReturnsOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
	rec := httptest.NewRecorder()

	handler.AtlasHealthz()(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func TestAtlasReadyz_ReturnsOKWhenAllCheckersPass(t *testing.T) {
	checker := &atlasMockChecker{healthy: true}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/readyz", nil)
	rec := httptest.NewRecorder()

	handler.AtlasReadyz(checker)(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func TestAtlasReadyz_Returns503WhenCheckerFails(t *testing.T) {
	checker := &atlasMockChecker{healthy: false}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/readyz", nil)
	rec := httptest.NewRecorder()

	handler.AtlasReadyz(checker)(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.JSONEq(t, `{"status":"unavailable","error":"test error"}`, rec.Body.String())
}

func TestAtlasReadyz_Returns503WhenMultipleCheckersAndOneFails(t *testing.T) {
	pass := &atlasMockChecker{healthy: true}
	fail := &atlasMockChecker{healthy: false}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/readyz", nil)
	rec := httptest.NewRecorder()

	handler.AtlasReadyz(pass, fail)(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.JSONEq(t, `{"status":"unavailable","error":"test error"}`, rec.Body.String())
}

type atlasMockChecker struct {
	healthy bool
}

func (m *atlasMockChecker) Ping() error {
	if !m.healthy {
		return errors.New("test error")
	}
	return nil
}