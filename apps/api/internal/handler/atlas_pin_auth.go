// FILE: apps/api/internal/handler/atlas_pin_auth.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide Atlas PIN authentication REST handlers for unlock, lock, and session check.
//   SCOPE: POST /api/v1/auth/pin/unlock validates PIN via PinService and creates session; POST /api/v1/auth/pin/lock revokes session idempotently; GET /api/v1/auth/session returns session status.
//   DEPENDS: apps/api/internal/atlas/service (PinService, SettingsService), apps/api/internal/atlas/middleware, apps/api/internal/atlas/repository/redis (PinSessionStore, PinAttemptStore).
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Atlas PIN auth handlers for WAVE-01.
// END_CHANGE_SUMMARY

package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	atlasMiddleware "monorepo-template/apps/api/internal/atlas/middleware"
	atlasRedis "monorepo-template/apps/api/internal/atlas/repository/redis"
	"monorepo-template/apps/api/internal/atlas/service"
	"monorepo-template/libs/go/logger"
)

type PinUnlockRequest struct {
	Pin string `json:"pin"`
}

type PinUnlockResponse struct {
	Success bool `json:"success"`
}

type SessionCheckResponse struct {
	SessionValid bool `json:"session_valid"`
}

type PinAuthHandler struct {
	pinService     service.PinService
	sessionStore   atlasRedis.PinSessionStore
	attemptStore   atlasRedis.PinAttemptStore
	cookieName     string
	idleTTL        time.Duration
	absoluteTTL    time.Duration
	cookieSecure   bool
	sameSite       http.SameSite
}

func NewPinAuthHandler(
	pinService service.PinService,
	sessionStore atlasRedis.PinSessionStore,
	attemptStore atlasRedis.PinAttemptStore,
	cookieName string,
	idleTTL, absoluteTTL time.Duration,
	cookieSecure bool,
	sameSite http.SameSite,
) *PinAuthHandler {
	return &PinAuthHandler{
		pinService:   pinService,
		sessionStore: sessionStore,
		attemptStore: attemptStore,
		cookieName:   cookieName,
		idleTTL:      idleTTL,
		absoluteTTL:  absoluteTTL,
		cookieSecure: cookieSecure,
		sameSite:     sameSite,
	}
}

func (h *PinAuthHandler) Unlock(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	var req PinUnlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.Pin == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "PIN is required"})
		return
	}

	userID := atlasMiddleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	ip := r.RemoteAddr
	locked, _, err := h.attemptStore.IsLocked(r.Context(), ip)
	if err != nil {
		log.Error("[Atlas][pin][BLOCK_PIN_UNLOCK] failed to check lockout", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if locked {
		writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "too many requests"})
		return
	}

	validPIN, err := h.pinService.Verify(r.Context(), userID, req.Pin)
	if err != nil {
		log.Error("[Atlas][pin][BLOCK_PIN_UNLOCK] verification failed", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	if !validPIN {
		if err := h.attemptStore.RegisterFailure(r.Context(), ip); err != nil {
			log.Error("[Atlas][pin][BLOCK_PIN_UNLOCK] failed to register attempt", zap.Error(err))
		}
		log.Warn("[Atlas][pin][BLOCK_PIN_UNLOCK] wrong PIN")
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid PIN"})
		return
	}

	if err := h.attemptStore.RegisterSuccess(r.Context(), ip); err != nil {
		log.Error("[Atlas][pin][BLOCK_PIN_UNLOCK] failed to clear attempts", zap.Error(err))
	}

	token, err := h.sessionStore.Create(r.Context(), userID, h.idleTTL, h.absoluteTTL)
	if err != nil {
		log.Error("[Atlas][pin][BLOCK_PIN_UNLOCK] failed to create session", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	atlasMiddleware.SetAtlasSessionCookie(w, h.cookieName, token, int(h.idleTTL.Seconds()), h.cookieSecure, h.sameSite)
	writeJSON(w, http.StatusOK, PinUnlockResponse{Success: true})
}

func (h *PinAuthHandler) Lock(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	cookie, err := r.Cookie(h.cookieName)
	if err == nil && cookie.Value != "" {
		if err := h.sessionStore.Revoke(r.Context(), cookie.Value); err != nil {
			log.Error("[Atlas][pin][BLOCK_PIN_LOCK] failed to revoke session", zap.Error(err))
		}
	}

	atlasMiddleware.ClearAtlasSessionCookie(w, h.cookieName)
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *PinAuthHandler) SessionCheck(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.cookieName)
	if err != nil || cookie.Value == "" {
		writeJSON(w, http.StatusOK, SessionCheckResponse{SessionValid: false})
		return
	}

	_, valid, err := h.sessionStore.Validate(r.Context(), cookie.Value)
	if err != nil {
		writeJSON(w, http.StatusOK, SessionCheckResponse{SessionValid: false})
		return
	}

	writeJSON(w, http.StatusOK, SessionCheckResponse{SessionValid: valid})
}