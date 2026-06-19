// FILE: apps/api/internal/handler/atlas_health.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide Atlas-specific health and readiness endpoints that do NOT bootstrap the default user.
//   SCOPE: /api/v1/healthz returns OK, /api/v1/readyz checks DB and Redis connectivity; no bootstrap side effects.
//   DEPENDS: apps/api/internal/repository/postgres.DB, apps/api/internal/repository/redis.Client.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Atlas health handler for WAVE-01.
// END_CHANGE_SUMMARY

package handler

import (
	"encoding/json"
	"net/http"
)

type AtlasHealthChecker interface {
	Ping() error
}

func AtlasHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func AtlasReadyz(checkers ...AtlasHealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		for _, c := range checkers {
			if err := c.Ping(); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"status": "unavailable",
					"error":  err.Error(),
				})
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}