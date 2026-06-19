// FILE: apps/api/internal/appconfig/config_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify API configuration validation, admin bootstrap env checks, and admin session defaults.
//   SCOPE: Pagination validation, env-only admin bootstrap, admin session cookie/security defaults, and placeholder rejection; excludes runtime startup wiring and persistence.
//   DEPENDS: apps/api/internal/appconfig, libs/go/config.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestConfig_PaginationValid - Verifies baseline config validation still accepts valid pagination.
//   TestConfig_Admin* - Verifies admin bootstrap and session configuration contracts.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Replaced stale JWT placeholder coverage with admin bootstrap and session config tests.
// END_CHANGE_SUMMARY

package appconfig_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/appconfig"
	"monorepo-template/libs/go/config"
)

func TestConfig_PaginationValid(t *testing.T) {
	cfg := validConfig()
	err := config.Validate(cfg)
	require.NoError(t, err)
}

func TestConfig_AdminSeedRequiredWhenBootstrapNeeded(t *testing.T) {
	cfg := validConfig()
	cfg.Admin = appconfig.AdminConfig{}
	err := appconfig.ValidateAdminBootstrap(cfg, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "initial admin email is required")
	assert.Contains(t, err.Error(), "initial admin password is required")
	assert.Contains(t, err.Error(), "initial admin name is required")
}

func TestConfig_AdminSeedNotRequiredWhenAdminExists(t *testing.T) {
	cfg := validConfig()
	cfg.Admin = appconfig.AdminConfig{}
	err := appconfig.ValidateAdminBootstrap(cfg, false)
	require.NoError(t, err)
}

func TestConfig_AdminBootstrapRequiresEnvKeysWhenTableEmpty(t *testing.T) {
	env := map[string]string{
		"ADMIN_INITIAL_EMAIL": "admin@example.com",
		"ADMIN_INITIAL_NAME":  "Template Admin",
	}
	err := appconfig.ValidateAdminBootstrapEnv(func(key string) (string, bool) {
		value, ok := env[key]
		return value, ok
	}, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ADMIN_INITIAL_PASSWORD is required when admin_users is empty")
}

func TestConfig_AdminBootstrapEnvNotRequiredWhenAdminExists(t *testing.T) {
	err := appconfig.ValidateAdminBootstrapEnv(func(key string) (string, bool) {
		return "", false
	}, false)
	require.NoError(t, err)
}

func TestConfig_AdminBootstrapEnvAcceptsAllRequiredValues(t *testing.T) {
	env := map[string]string{
		"ADMIN_INITIAL_EMAIL":    "admin@example.test",
		"ADMIN_INITIAL_PASSWORD": "StrongPassword123!",
		"ADMIN_INITIAL_NAME":     "Template Admin",
	}
	err := appconfig.ValidateAdminBootstrapEnv(func(key string) (string, bool) {
		value, ok := env[key]
		return value, ok
	}, true)
	require.NoError(t, err)
}

func TestConfig_AdminBootstrapRejectsPlaceholderSeedOutsideDevelopment(t *testing.T) {
	cfg := validConfig()
	cfg.Server.Env = "production"
	cfg.Admin.InitialEmail = "admin@example.com"
	cfg.Admin.InitialPassword = "ChangeMeAdmin123!"
	cfg.Admin.InitialName = "Template Admin"
	err := appconfig.ValidateAdminBootstrap(cfg, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "initial admin bootstrap values must not use example placeholders outside development")
}

func TestConfig_AdminBootstrapAcceptsNonPlaceholderSeedOutsideDevelopment(t *testing.T) {
	cfg := validConfig()
	cfg.Server.Env = "production"
	cfg.Admin.InitialEmail = "ops@example.test"
	cfg.Admin.InitialPassword = "VeryStrongPassword123!"
	cfg.Admin.InitialName = "Ops Owner"
	err := appconfig.ValidateAdminBootstrap(cfg, true)
	require.NoError(t, err)
}

func TestConfig_AdminEnvOverlayHydratesEnvOnlyAdminFields(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
server:
  port: 8080
  env: staging
postgres:
  host: localhost
  port: 5432
  user: app
  db: test
redis:
  host: localhost
  port: 6379
admin:
  origins:
    - "http://localhost:3100"
admin_session:
  cookie_name: web_admin_session
  ttl: 168h
  cookie_secure: true
  same_site: Lax
pagination:
  default_page_size: 20
  max_page_size: 100
`), 0o600))
	t.Setenv("ADMIN_INITIAL_EMAIL", "ops-admin@example.test")
	t.Setenv("ADMIN_INITIAL_PASSWORD", "StrongPassword123!")
	t.Setenv("ADMIN_INITIAL_NAME", "Ops Admin")
	t.Setenv("ADMIN_ORIGINS", "https://admin.example.com,https://admin2.example.com")
	t.Setenv("ADMIN_SESSION_COOKIE_NAME", "ops_admin_session")
	t.Setenv("ADMIN_SESSION_TTL", "24h")
	t.Setenv("ADMIN_SESSION_COOKIE_SECURE", "1")
	t.Setenv("ADMIN_SESSION_SAME_SITE", "Strict")
	t.Setenv("ADMIN_SESSION_KEY_SECRET", "real-session-key-secret")

	cfg, err := config.Load[appconfig.Config](config.Options{ConfigPath: configPath})
	require.NoError(t, err)
	require.NoError(t, appconfig.ApplyAdminEnvOverlay(&cfg, os.LookupEnv))
	require.NoError(t, appconfig.ApplyAdminDefaults(&cfg))

	assert.Equal(t, "ops-admin@example.test", cfg.Admin.InitialEmail)
	assert.Equal(t, "StrongPassword123!", cfg.Admin.InitialPassword)
	assert.Equal(t, "Ops Admin", cfg.Admin.InitialName)
	assert.Equal(t, []string{"https://admin.example.com", "https://admin2.example.com"}, cfg.Admin.Origins)
	assert.Equal(t, "ops_admin_session", cfg.AdminSession.CookieName)
	assert.Equal(t, 24*time.Hour, cfg.AdminSession.TTL)
	assert.Equal(t, "true", cfg.AdminSession.CookieSecure)
	assert.Equal(t, "Strict", cfg.AdminSession.SameSite)
	assert.Equal(t, "real-session-key-secret", cfg.AdminSession.KeySecret)
}

func TestConfig_AdminEnvOverlayRejectsInvalidSessionTTL(t *testing.T) {
	cfg := validConfig()
	err := appconfig.ApplyAdminEnvOverlay(&cfg, func(key string) (string, bool) {
		if key == "ADMIN_SESSION_TTL" {
			return "not-a-duration", true
		}
		return "", false
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin session ttl is invalid")
}

func TestConfig_AdminSessionDefaultsAndValidation(t *testing.T) {
	cfg := validConfig()
	cfg.Admin.Origins = nil
	cfg.AdminSession.CookieName = ""
	cfg.AdminSession.TTL = 0
	cfg.AdminSession.CookieSecure = ""
	cfg.AdminSession.SameSite = ""
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.NoError(t, err)
	assert.Equal(t, "web_admin_session", cfg.AdminSession.CookieName)
	assert.Equal(t, 168*time.Hour, cfg.AdminSession.TTL)
	assert.Equal(t, "auto", cfg.AdminSession.CookieSecure)
	assert.Equal(t, "Lax", cfg.AdminSession.SameSite)
	assert.Equal(t, []string{"http://localhost:3100", "http://127.0.0.1:3100"}, cfg.Admin.Origins)
}

func TestConfig_AdminSessionAcceptsNumericSecureModes(t *testing.T) {
	cfg := validConfig()
	cfg.AdminSession.CookieSecure = "1"
	require.NoError(t, appconfig.ApplyAdminDefaults(&cfg))
	assert.Equal(t, "true", cfg.AdminSession.CookieSecure)

	cfg = validConfig()
	cfg.AdminSession.CookieSecure = "0"
	require.NoError(t, appconfig.ApplyAdminDefaults(&cfg))
	assert.Equal(t, "false", cfg.AdminSession.CookieSecure)
}

func TestConfig_AdminSessionRejectsMissingSecret(t *testing.T) {
	cfg := validConfig()
	cfg.AdminSession.KeySecret = " "
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin session key secret is required")
}

func TestConfig_AdminSessionRejectsInvalidSecureMode(t *testing.T) {
	cfg := validConfig()
	cfg.AdminSession.CookieSecure = "sometimes"
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin session cookie secure must be auto, true, or false")
}

func TestConfig_AdminSessionRejectsInvalidSameSiteMode(t *testing.T) {
	cfg := validConfig()
	cfg.AdminSession.SameSite = "Relaxed"
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin session same site must be Lax, Strict, or None")
}

func TestConfig_AdminSessionRejectsProductionInsecureCookie(t *testing.T) {
	cfg := validConfig()
	cfg.Server.Env = "production"
	cfg.AdminSession.CookieSecure = "false"
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "production admin session cookie secure must be true or auto")
}

func TestConfig_AdminSessionRejectsSameSiteNoneWithoutSecureCookie(t *testing.T) {
	cfg := validConfig()
	cfg.AdminSession.CookieSecure = "false"
	cfg.AdminSession.SameSite = "None"
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin session SameSite=None requires a secure cookie")
}

func TestConfig_AdminSessionRejectsSameSiteNoneAutoOutsideProduction(t *testing.T) {
	cfg := validConfig()
	cfg.Server.Env = "staging"
	cfg.AdminSession.CookieSecure = "auto"
	cfg.AdminSession.SameSite = "None"
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin session SameSite=None requires a secure cookie")
}

func TestConfig_AdminSessionRejectsPlaceholderSecretOutsideDevelopment(t *testing.T) {
	cfg := validConfig()
	cfg.Server.Env = "production"
	cfg.AdminSession.KeySecret = "change-me-session-key-secret"
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin session key secret must be provided by environment outside development")
}

func validConfig() appconfig.Config {
	return appconfig.Config{
		Server:   config.ServerConfig{Port: 8080, Env: "development"},
		Postgres: config.PostgresConfig{Host: "localhost", Port: 5432, User: "app", DB: "test"},
		Redis:    config.RedisConfig{Host: "localhost", Port: 6379},
		Admin: appconfig.AdminConfig{
			InitialEmail:    "admin@example.com",
			InitialPassword: "StrongPassword123!",
			InitialName:     "Template Admin",
			Origins:         []string{"http://localhost:3100", "http://127.0.0.1:3100"},
		},
		AdminSession: appconfig.AdminSessionConfig{
			CookieName:   "web_admin_session",
			TTL:          168 * time.Hour,
			CookieSecure: "auto",
			SameSite:     "Lax",
			KeySecret:    strings.Repeat("x", 24),
		},
		Pagination: appconfig.PaginationConfig{DefaultPageSize: 20, MaxPageSize: 100},
	}
}
