// FILE: apps/api/internal/middleware/admin_origin.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Enforce the web-admin browser-origin boundary for credentialed admin GraphQL requests.
//   SCOPE: Strict Origin/Referer allowlisting for unsafe browser requests; excludes CORS header emission and GraphQL auth decisions.
//   DEPENDS: net/http, net/url.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminOriginGuard - Rejects unsafe admin GraphQL browser requests from non-admin origins.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin origin guard.
// END_CHANGE_SUMMARY

package middleware

import (
	"net/http"
	"net/url"

	"go.uber.org/zap"

	"monorepo-template/libs/go/logger"
)

// START_CONTRACT: AdminOriginGuard
//
//	PURPOSE: Reject unsafe credentialed admin GraphQL browser requests from missing or untrusted origins.
//	INPUTS: { allowedOrigins: []string - exact web-admin origins }
//	OUTPUTS: { func(http.Handler) http.Handler - middleware }
//	SIDE_EFFECTS: May write HTTP 403 before GraphQL execution.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: AdminOriginGuard
func AdminOriginGuard(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := map[string]struct{}{}
	for _, origin := range allowedOrigins {
		allowed[origin] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			if r.Method == http.MethodGet || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			origin := r.Header.Get("Origin")
			if origin == "" {
				origin = originFromReferer(r.Header.Get("Referer"))
			}
			if origin == "" {
				log.Warn("[AdminAuth][csrf][BLOCK_VALIDATE_ORIGIN] missing origin")
				http.Error(w, "admin origin is required", http.StatusForbidden)
				return
			}
			if _, ok := allowed[origin]; !ok {
				log.Warn("[AdminAuth][csrf][BLOCK_VALIDATE_ORIGIN] rejected origin", zap.String("origin", origin))
				http.Error(w, "admin origin is not allowed", http.StatusForbidden)
				return
			}
			log.Debug("[AdminAuth][csrf][BLOCK_VALIDATE_ORIGIN] allowed origin", zap.String("origin", origin))
			next.ServeHTTP(w, r)
		})
	}
}

func originFromReferer(value string) string {
	if value == "" {
		return ""
	}
	parsed, err := url.Parse(value)
	if err != nil {
		return ""
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	return parsed.Scheme + "://" + parsed.Host
}
