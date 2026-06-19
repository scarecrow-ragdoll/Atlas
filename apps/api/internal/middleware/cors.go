// FILE: apps/api/internal/middleware/cors.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Apply API CORS headers for public and credentialed admin browser requests.
//   SCOPE: Allowed origin/method/header matching, credentialed origin reflection, wildcard handling, and OPTIONS termination; excludes CSRF/origin authorization decisions.
//   DEPENDS: net/http, strings.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   CORSConfig - Allowed CORS origins, methods, headers, and credential mode.
//   DefaultCORSConfig - Local development CORS defaults.
//   CORS - Emits CORS response headers and handles preflight requests.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added credentialed admin CORS support without wildcard credentials.
// END_CHANGE_SUMMARY

package middleware

import (
	"net/http"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

// DefaultCORSConfig returns a CORSConfig permitting local web and web-admin
// origins with standard methods and headers. Suitable for local development.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://localhost:3001",
			"http://127.0.0.1:3001",
			"http://localhost:3002",
			"http://127.0.0.1:3002",
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}
}

func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowOrigin := ""

			for _, allowed := range cfg.AllowedOrigins {
				if allowed == "*" && !cfg.AllowCredentials {
					allowOrigin = "*"
					break
				}
				if allowed == origin {
					allowOrigin = origin
					break
				}
			}
			if allowOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
				if cfg.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
