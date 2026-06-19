// FILE: apps/api/internal/atlas/middleware/pin_guard.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide Atlas PIN guard middleware and Atlas user context middleware for route group protection.
//   SCOPE: AtlasUserContext attaches cached default user ID; AtlasPinGuard validates PIN session cookie when PIN is enabled; no string-based skip lists.
//   DEPENDS: apps/api/internal/atlas/service (PinService, BootstrapService), apps/api/internal/atlas/repository/redis (PinSessionStore).
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Atlas PIN guard middleware for WAVE-01.
// END_CHANGE_SUMMARY

package middleware

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	atlasRedis "monorepo-template/apps/api/internal/atlas/repository/redis"
	"monorepo-template/apps/api/internal/atlas/service"
	"monorepo-template/libs/go/logger"
)

type atlasUserKey struct{}
type atlasSessionTokenKey struct{}

func ContextWithAtlasUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, atlasUserKey{}, userID)
}

func GetAtlasUserID(ctx context.Context) string {
	uid, _ := ctx.Value(atlasUserKey{}).(string)
	return uid
}

func ContextWithAtlasSessionToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, atlasSessionTokenKey{}, token)
}

func GetAtlasSessionToken(ctx context.Context) string {
	token, _ := ctx.Value(atlasSessionTokenKey{}).(string)
	return token
}

func AtlasUserContext(bootstrapService service.BootstrapService) func(http.Handler) http.Handler {
	var cachedUserID string

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())

			if cachedUserID == "" {
				uid, err := bootstrapService.EnsureDefaultUser(r.Context())
				if err != nil {
					log.Error("[Atlas][user][BLOCK_ATLAS_USER_CONTEXT] failed to resolve default user", zap.Error(err))
					http.Error(w, "atlas service unavailable", http.StatusServiceUnavailable)
					return
				}
				cachedUserID = uid
				log.Debug("[Atlas][user][BLOCK_ATLAS_USER_CONTEXT] resolved default user", zap.String("user_id", uid))
			}

			ctx := ContextWithAtlasUserID(r.Context(), cachedUserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AtlasPinGuard(pinService service.PinService, sessionStore atlasRedis.PinSessionStore, cookieName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())

			userID := GetAtlasUserID(r.Context())
			if userID == "" {
				log.Warn("[Atlas][pin][BLOCK_PIN_GUARD] no user in context")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			pinEnabled, err := pinService.IsEnabled(r.Context(), userID)
			if err != nil {
				log.Error("[Atlas][pin][BLOCK_PIN_GUARD] failed to check PIN status", zap.Error(err))
				http.Error(w, "atlas service unavailable", http.StatusServiceUnavailable)
				return
			}

			if !pinEnabled {
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := r.Cookie(cookieName)
			if err != nil || cookie.Value == "" {
				log.Debug("[Atlas][pin][BLOCK_PIN_GUARD] no session cookie")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			sessionUserID, valid, err := sessionStore.Validate(r.Context(), cookie.Value)
			if err != nil {
				log.Error("[Atlas][pin][BLOCK_PIN_GUARD] session validation failed", zap.Error(err))
				http.Error(w, "atlas service unavailable", http.StatusServiceUnavailable)
				return
			}

			if !valid {
				log.Debug("[Atlas][pin][BLOCK_PIN_GUARD] invalid session")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if sessionUserID != userID {
				log.Warn("[Atlas][pin][BLOCK_PIN_GUARD] session user mismatch",
					zap.String("session_user", sessionUserID),
					zap.String("context_user", userID),
				)
				_ = sessionStore.Revoke(r.Context(), cookie.Value)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := ContextWithAtlasSessionToken(r.Context(), cookie.Value)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClearAtlasSessionCookie(w http.ResponseWriter, cookieName string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func SetAtlasSessionCookie(w http.ResponseWriter, cookieName, token string, maxAge int, secure bool, sameSite http.SameSite) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   maxAge,
	})
}