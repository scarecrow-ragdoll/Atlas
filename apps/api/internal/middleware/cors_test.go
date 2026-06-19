// FILE: apps/api/internal/middleware/cors_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify public and credentialed admin CORS behavior.
//   SCOPE: Allowed/disallowed origins, preflight handling, credential reflection, and wildcard credential rejection; excludes admin CSRF origin authorization.
//   DEPENDS: apps/api/internal/middleware, httptest.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestCORS_* - Covers CORS success, failure, preflight, and credential mode behavior.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added credentialed admin CORS coverage.
// END_CHANGE_SUMMARY

package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"monorepo-template/apps/api/internal/middleware"
)

func TestCORS_AllowedOrigins(t *testing.T) {
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3001",
		"http://localhost:3002",
		"http://127.0.0.1:3002",
	}

	for _, origin := range allowedOrigins {
		t.Run(origin, func(t *testing.T) {
			handler := middleware.CORS(middleware.DefaultCORSConfig())(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Origin", origin)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, origin, rec.Header().Get("Access-Control-Allow-Origin"))
		})
	}
}

func TestCORS_PreflightReturnsNoContent(t *testing.T) {
	handler := middleware.CORS(middleware.DefaultCORSConfig())(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestCORS_DisallowedOrigin(t *testing.T) {
	handler := middleware.CORS(middleware.DefaultCORSConfig())(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://evil.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_CredentialedRejectsWildcardOrigin(t *testing.T) {
	handler := middleware.CORS(middleware.CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORS_AllowsWildcardOriginWithoutCredentials(t *testing.T) {
	handler := middleware.CORS(middleware.CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Content-Type"},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("Origin", "http://anywhere.test")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORS_CredentialedAdminOrigin(t *testing.T) {
	handler := middleware.CORS(middleware.CORSConfig{
		AllowedOrigins:   []string{"http://localhost:3100"},
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, "http://localhost:3100", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
}
