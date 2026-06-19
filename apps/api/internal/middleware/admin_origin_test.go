// FILE: apps/api/internal/middleware/admin_origin_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify admin GraphQL origin and CSRF guard behavior.
//   SCOPE: Allows configured web-admin origins and rejects public or disallowed browser origins for protected/session-mutating GraphQL requests; excludes GraphQL parsing.
//   DEPENDS: apps/api/internal/middleware, httptest.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestAdminOriginGuard_* - Origin allow and deny coverage.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin origin guard coverage.
// END_CHANGE_SUMMARY

package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/libs/go/logger"
)

func TestAdminOriginGuard_AllowsConfiguredOrigin(t *testing.T) {
	handler := middleware.AdminOriginGuard([]string{"http://localhost:3100"})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"mutation { createAdmin(input:{email:\"a@example.com\",name:\"A\",password:\"StrongPassword123!\"}) { __typename } }"}`))
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminOriginGuard_AllowsConfiguredReferer(t *testing.T) {
	handler := middleware.AdminOriginGuard([]string{"http://localhost:3100"})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"mutation { logoutAdmin { __typename } }"}`))
	req.Header.Set("Referer", "http://localhost:3100/admin/users")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminOriginGuard_AllowsSafeMethodsWithoutOrigin(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodOptions} {
		t.Run(method, func(t *testing.T) {
			handler := middleware.AdminOriginGuard([]string{"http://localhost:3100"})(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			)
			req := httptest.NewRequest(method, "/graphql", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestAdminOriginGuard_RejectsPublicWebOrigin(t *testing.T) {
	handler := middleware.AdminOriginGuard([]string{"http://localhost:3100"})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"mutation { createAdmin(input:{email:\"a@example.com\",name:\"A\",password:\"StrongPassword123!\"}) { __typename } }"}`))
	req.Header.Set("Origin", "http://localhost:3101")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAdminOriginGuard_RejectsMalformedReferer(t *testing.T) {
	handler := middleware.AdminOriginGuard([]string{"http://localhost:3100"})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"mutation { logoutAdmin { __typename } }"}`))
	req.Header.Set("Referer", "://bad-url")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAdminOriginGuard_RejectsRelativeReferer(t *testing.T) {
	handler := middleware.AdminOriginGuard([]string{"http://localhost:3100"})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"mutation { logoutAdmin { __typename } }"}`))
	req.Header.Set("Referer", "/admin/users")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAdminOriginGuard_RejectsMissingOriginForUnsafeRequest(t *testing.T) {
	handler := middleware.AdminOriginGuard([]string{"http://localhost:3100"})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"mutation { logoutAdmin { __typename } }"}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAdminOriginGuard_LogsMarkerWithoutSecrets(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	handler := middleware.AdminOriginGuard([]string{"http://localhost:3100"})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	body := `{"query":"mutation { loginAdmin(input:{email:\"admin@example.com\",password:\"StrongPassword123!\"}) { __typename } }"}`
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(body))
	req = req.WithContext(logger.WithContext(req.Context(), zap.New(core)))
	req.Header.Set("Origin", "http://localhost:3101")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	joined := logText(logs.All())
	assert.Contains(t, joined, "[AdminAuth][csrf][BLOCK_VALIDATE_ORIGIN]")
	assert.NotContains(t, joined, "admin@example.com")
	assert.NotContains(t, joined, "StrongPassword123!")
	assert.NotContains(t, joined, "loginAdmin")
	assert.NotContains(t, joined, body)
}
