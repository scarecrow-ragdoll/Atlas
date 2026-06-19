package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"monorepo-template/apps/api/internal/handler"
)

func TestHealthz_ReturnsOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.Healthz()(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func TestReadyz_HealthyDeps(t *testing.T) {
	checker := &mockChecker{healthy: true}
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.Readyz(checker)(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func TestReadyz_UnhealthyDeps(t *testing.T) {
	checker := &mockChecker{healthy: false}
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.Readyz(checker)(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

type mockChecker struct {
	healthy bool
}

func (m *mockChecker) Ping() error {
	if !m.healthy {
		return assert.AnError
	}
	return nil
}
