// FILE: apps/api/internal/middleware/admin_auth.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide admin principal context helpers, session lookup middleware, and session cookie helpers.
//   SCOPE: AdminPrincipal context storage, session-cookie lookup through the auth service boundary, cookie set/clear behavior, and cookie config; excludes Redis implementation and GraphQL resolver decisions.
//   DEPENDS: context, net/http, apps/api/internal/service.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminPrincipal - Request-scoped authenticated admin identity.
//   ContextWithAdminPrincipal - Stores an admin principal in context.
//   GetAdminPrincipal - Reads an admin principal from context.
//   ContextWithAdminSessionID - Stores the raw session id only inside request context for logout.
//   AdminSessionIDFromContext - Reads the request-scoped session id for logout.
//   AdminSessionMiddleware - Resolves session cookies into request-scoped admin context.
//   SetAdminSessionCookie - Sets the httpOnly session cookie.
//   ClearAdminSessionCookie - Clears the session cookie.
//   ContextWithAdminCookieBridge - Stores the response/cookie sink for GraphQL resolvers.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added GraphQL cookie bridge helpers.
// END_CHANGE_SUMMARY

package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/appconfig"
	"monorepo-template/apps/api/internal/service"
	"monorepo-template/libs/go/logger"
)

type adminPrincipalKey struct{}
type adminSessionIDKey struct{}
type adminCookieBridgeKey struct{}

type AdminPrincipal struct {
	ID        string
	Email     string
	Name      string
	Role      string
	CreatedAt string
	UpdatedAt string
}

type AdminSessionResolver interface {
	CurrentAdmin(ctx context.Context, sessionID string) (*service.Admin, error)
}

type AdminCookieConfig struct {
	Name     string
	Path     string
	MaxAge   int
	Secure   bool
	SameSite http.SameSite
}

type AdminCookieBridge struct {
	Response http.ResponseWriter
	Config   AdminCookieConfig
}

func ContextWithAdminPrincipal(ctx context.Context, principal AdminPrincipal) context.Context {
	return context.WithValue(ctx, adminPrincipalKey{}, principal)
}

func GetAdminPrincipal(ctx context.Context) (AdminPrincipal, bool) {
	principal, ok := ctx.Value(adminPrincipalKey{}).(AdminPrincipal)
	return principal, ok
}

func ContextWithAdminSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, adminSessionIDKey{}, sessionID)
}

func AdminSessionIDFromContext(ctx context.Context) string {
	sessionID, _ := ctx.Value(adminSessionIDKey{}).(string)
	return sessionID
}

func ContextWithAdminCookieBridge(ctx context.Context, bridge AdminCookieBridge) context.Context {
	return context.WithValue(ctx, adminCookieBridgeKey{}, bridge)
}

func SetAdminSessionCookieFromContext(ctx context.Context, sessionID string) {
	if bridge, ok := ctx.Value(adminCookieBridgeKey{}).(AdminCookieBridge); ok {
		SetAdminSessionCookie(bridge.Response, bridge.Config, sessionID)
	}
}

func ClearAdminSessionCookieFromContext(ctx context.Context) {
	if bridge, ok := ctx.Value(adminCookieBridgeKey{}).(AdminCookieBridge); ok {
		ClearAdminSessionCookie(bridge.Response, bridge.Config)
	}
}

// START_CONTRACT: AdminSessionMiddleware
//
//	PURPOSE: Resolve an admin session cookie into request-scoped admin principal context.
//	INPUTS: { resolver: AdminSessionResolver - service boundary, cookieName: string - session cookie name }
//	OUTPUTS: { func(http.Handler) http.Handler - middleware }
//	SIDE_EFFECTS: Reads session/admin state through resolver and may write 500 on lookup failure.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: AdminSessionMiddleware
func AdminSessionMiddleware(resolver AdminSessionResolver, cookieName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			cookie, err := r.Cookie(cookieName)
			if err != nil || cookie.Value == "" {
				next.ServeHTTP(w, r)
				return
			}

			log.Debug("[AdminAuth][session][BLOCK_VALIDATE_SESSION] session cookie present")
			admin, err := resolver.CurrentAdmin(r.Context(), cookie.Value)
			if err != nil {
				log.Error("[AdminAuth][session][BLOCK_VALIDATE_SESSION] session lookup failed", zap.Error(err))
				http.Error(w, "admin session lookup failed", http.StatusInternalServerError)
				return
			}

			ctx := ContextWithAdminSessionID(r.Context(), cookie.Value)
			if admin != nil && admin.IsActive {
				ctx = ContextWithAdminPrincipal(ctx, AdminPrincipal{
					ID:        admin.ID,
					Email:     admin.Email,
					Name:      admin.Name,
					Role:      admin.Role,
					CreatedAt: admin.CreatedAt,
					UpdatedAt: admin.UpdatedAt,
				})
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminCookieConfigFromConfig(cfg appconfig.AdminSessionConfig, env string) AdminCookieConfig {
	secureSetting := strings.ToLower(strings.TrimSpace(cfg.CookieSecure))
	environment := strings.ToLower(strings.TrimSpace(env))
	secure := secureSetting == "true" || (secureSetting == "auto" && environment == "production")
	return AdminCookieConfig{
		Name:     cfg.CookieName,
		Path:     "/graphql",
		MaxAge:   int(cfg.TTL.Seconds()),
		Secure:   secure,
		SameSite: adminSameSiteMode(cfg.SameSite),
	}
}

func WithAdminCookieBridge(next http.Handler, cfg AdminCookieConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := ContextWithAdminCookieBridge(r.Context(), AdminCookieBridge{Response: w, Config: cfg})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetAdminSessionCookie(w http.ResponseWriter, cfg AdminCookieConfig, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.Name,
		Value:    sessionID,
		Path:     cfg.Path,
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		MaxAge:   cfg.MaxAge,
	})
}

func adminSameSiteMode(value string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

func ClearAdminSessionCookie(w http.ResponseWriter, cfg AdminCookieConfig) {
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.Name,
		Value:    "",
		Path:     cfg.Path,
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		MaxAge:   -1,
		Expires:  time.Now().Add(-time.Hour),
	})
}
