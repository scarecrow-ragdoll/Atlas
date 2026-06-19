// FILE: apps/api/internal/appconfig/config.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define API service configuration, admin bootstrap checks, and admin session defaults.
//   SCOPE: Server, logging, PostgreSQL, Redis, admin seed, admin session, and pagination config; excludes secret loading implementation and persistence behavior.
//   DEPENDS: libs/go/config.
//   LINKS: M-API / V-M-API.
//   ROLE: CONFIG
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminConfig - Env-backed initial-admin and web-admin origin settings.
//   AdminSessionConfig - Cookie/session settings for Redis-backed admin sessions.
//   Config - Full API configuration shape.
//   ApplyAdminEnvOverlay - Hydrates env-only admin fields after YAML config load.
//   ApplyAdminDefaults - Applies non-secret defaults and validates session security settings.
//   ValidateAdminBootstrapEnv - Proves first-admin seed values are present in env when the table is empty.
//   ValidateAdminBootstrap - Validates first-admin config values only when bootstrap is required.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Replaced placeholder bearer auth config with admin bootstrap and session config.
// END_CHANGE_SUMMARY

package appconfig

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"monorepo-template/libs/go/config"
)

// AdminConfig holds bootstrap and browser-origin settings for web-admin authentication.
type AdminConfig struct {
	InitialEmail    string   `mapstructure:"initial_email"`
	InitialPassword string   `mapstructure:"initial_password"`
	InitialName     string   `mapstructure:"initial_name"`
	Origins         []string `mapstructure:"origins"`
}

// AdminSessionConfig holds cookie and Redis-key settings for admin sessions.
type AdminSessionConfig struct {
	CookieName   string        `mapstructure:"cookie_name"`
	TTL          time.Duration `mapstructure:"ttl"`
	CookieSecure string        `mapstructure:"cookie_secure"`
	SameSite     string        `mapstructure:"same_site"`
	KeySecret    string        `mapstructure:"key_secret"`
}

// PaginationConfig holds pagination defaults.
type PaginationConfig struct {
	DefaultPageSize int `mapstructure:"default_page_size" validate:"gt=0"`
	MaxPageSize     int `mapstructure:"max_page_size"     validate:"gt=0"`
}

// Config is the API service configuration.
type Config struct {
	Server       config.ServerConfig   `mapstructure:"server"`
	Log          config.LogConfig      `mapstructure:"log"`
	Postgres     config.PostgresConfig `mapstructure:"postgres"`
	Redis        config.RedisConfig    `mapstructure:"redis"`
	Admin        AdminConfig           `mapstructure:"admin"`
	AdminSession AdminSessionConfig    `mapstructure:"admin_session"`
	Pagination   PaginationConfig      `mapstructure:"pagination"`
	AtlasPin        AtlasPinConfig        `mapstructure:"atlas_pin"`
	AtlasPinSession AtlasPinSessionConfig `mapstructure:"atlas_pin_session"`
	AtlasPinAttempt AtlasPinAttemptConfig `mapstructure:"atlas_pin_attempt"`
}

// AtlasPinConfig holds Argon2id parameters and PIN validation settings.
type AtlasPinConfig struct {
	Argon2Memory      uint32 `mapstructure:"argon2_memory"`
	Argon2Iterations  uint32 `mapstructure:"argon2_iterations"`
	Argon2Parallelism uint8  `mapstructure:"argon2_parallelism"`
	Argon2KeyLength   uint32 `mapstructure:"argon2_key_length"`
	MinLength         int    `mapstructure:"min_length"`
	MaxLength         int    `mapstructure:"max_length"`
}

// AtlasPinSessionConfig holds cookie and Redis session settings for Atlas PIN sessions.
type AtlasPinSessionConfig struct {
	CookieName   string        `mapstructure:"cookie_name"`
	IdleTTL      time.Duration `mapstructure:"idle_ttl"`
	AbsoluteTTL  time.Duration `mapstructure:"absolute_ttl"`
	CookieSecure string        `mapstructure:"cookie_secure"`
	SameSite     string        `mapstructure:"same_site"`
}

// AtlasPinAttemptConfig holds rate limiting settings for PIN brute-force protection.
type AtlasPinAttemptConfig struct {
	MaxFailures       int           `mapstructure:"max_failures"`
	LockoutDuration   time.Duration `mapstructure:"lockout_duration"`
	EscalatedDuration time.Duration `mapstructure:"escalated_duration"`
}

// ApplyAdminEnvOverlay hydrates env-only admin secrets and comma-separated origin overrides.
func ApplyAdminEnvOverlay(cfg *Config, lookup func(string) (string, bool)) error {
	if value, ok := envString(lookup, "ADMIN_INITIAL_EMAIL"); ok {
		cfg.Admin.InitialEmail = value
	}
	if value, ok := envString(lookup, "ADMIN_INITIAL_PASSWORD"); ok {
		cfg.Admin.InitialPassword = value
	}
	if value, ok := envString(lookup, "ADMIN_INITIAL_NAME"); ok {
		cfg.Admin.InitialName = value
	}
	if value, ok := envString(lookup, "ADMIN_ORIGINS"); ok {
		cfg.Admin.Origins = splitCSV(value)
	}
	if value, ok := envString(lookup, "ADMIN_SESSION_COOKIE_NAME"); ok {
		cfg.AdminSession.CookieName = value
	}
	if value, ok := envString(lookup, "ADMIN_SESSION_TTL"); ok {
		ttl, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("admin session ttl is invalid: %w", err)
		}
		cfg.AdminSession.TTL = ttl
	}
	if value, ok := envString(lookup, "ADMIN_SESSION_COOKIE_SECURE"); ok {
		cfg.AdminSession.CookieSecure = value
	}
	if value, ok := envString(lookup, "ADMIN_SESSION_SAME_SITE"); ok {
		cfg.AdminSession.SameSite = value
	}
	if value, ok := envString(lookup, "ADMIN_SESSION_KEY_SECRET"); ok {
		cfg.AdminSession.KeySecret = value
	}
	return nil
}

// ApplyAdminDefaults fills safe local defaults and validates enum-like values.
func ApplyAdminDefaults(cfg *Config) error {
	env := strings.ToLower(strings.TrimSpace(cfg.Server.Env))
	if cfg.AdminSession.CookieName == "" {
		cfg.AdminSession.CookieName = "web_admin_session"
	}
	if cfg.AdminSession.TTL == 0 {
		cfg.AdminSession.TTL = 168 * time.Hour
	}
	if cfg.AdminSession.CookieSecure == "" {
		cfg.AdminSession.CookieSecure = "auto"
	}
	cfg.AdminSession.CookieSecure = normalizeCookieSecure(cfg.AdminSession.CookieSecure)
	if cfg.AdminSession.SameSite == "" {
		cfg.AdminSession.SameSite = "Lax"
	}
	if len(cfg.Admin.Origins) == 0 {
		cfg.Admin.Origins = []string{"http://localhost:3100", "http://127.0.0.1:3100"}
	}
	if strings.TrimSpace(cfg.AdminSession.KeySecret) == "" {
		return fmt.Errorf("admin session key secret is required")
	}
	if cfg.AdminSession.CookieSecure != "auto" && cfg.AdminSession.CookieSecure != "true" && cfg.AdminSession.CookieSecure != "false" {
		return fmt.Errorf("admin session cookie secure must be auto, true, or false")
	}
	if cfg.AdminSession.SameSite != "Lax" && cfg.AdminSession.SameSite != "Strict" && cfg.AdminSession.SameSite != "None" {
		return fmt.Errorf("admin session same site must be Lax, Strict, or None")
	}
	effectiveSecure := cfg.AdminSession.CookieSecure == "true" || (cfg.AdminSession.CookieSecure == "auto" && env == "production")
	if env == "production" && cfg.AdminSession.CookieSecure == "false" {
		return fmt.Errorf("production admin session cookie secure must be true or auto")
	}
	if cfg.AdminSession.SameSite == "None" && !effectiveSecure {
		return fmt.Errorf("admin session SameSite=None requires a secure cookie")
	}
	if env != "" && env != "development" && isPlaceholderAdminSessionSecret(cfg.AdminSession.KeySecret) {
		return fmt.Errorf("admin session key secret must be provided by environment outside development")
	}
	return nil
}

// ValidateAdminBootstrapEnv ensures the first admin identity is supplied by env when needed.
func ValidateAdminBootstrapEnv(lookup func(string) (string, bool), adminTableEmpty bool) error {
	if !adminTableEmpty {
		return nil
	}
	required := []string{"ADMIN_INITIAL_EMAIL", "ADMIN_INITIAL_PASSWORD", "ADMIN_INITIAL_NAME"}
	var missing []string
	for _, key := range required {
		value, ok := lookup(key)
		if !ok || strings.TrimSpace(value) == "" {
			missing = append(missing, key+" is required when admin_users is empty")
		}
	}
	if len(missing) > 0 {
		return errors.New(strings.Join(missing, "; "))
	}
	return nil
}

// ValidateAdminBootstrap enforces first-admin env only when no admins exist.
func ValidateAdminBootstrap(cfg Config, adminTableEmpty bool) error {
	if !adminTableEmpty {
		return nil
	}
	env := strings.ToLower(strings.TrimSpace(cfg.Server.Env))
	var missing []string
	if strings.TrimSpace(cfg.Admin.InitialEmail) == "" {
		missing = append(missing, "initial admin email is required")
	}
	if strings.TrimSpace(cfg.Admin.InitialPassword) == "" {
		missing = append(missing, "initial admin password is required")
	}
	if strings.TrimSpace(cfg.Admin.InitialName) == "" {
		missing = append(missing, "initial admin name is required")
	}
	if len(missing) > 0 {
		return errors.New(strings.Join(missing, "; "))
	}
	if env != "" && env != "development" && isPlaceholderInitialAdmin(cfg.Admin) {
		return fmt.Errorf("initial admin bootstrap values must not use example placeholders outside development")
	}
	return nil
}

func isPlaceholderAdminSessionSecret(value string) bool {
	trimmed := strings.TrimSpace(value)
	return trimmed == "change-me-session-key-secret" || trimmed == "dev-session-key-secret"
}

func normalizeCookieSecure(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1":
		return "true"
	case "0":
		return "false"
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}

func isPlaceholderInitialAdmin(admin AdminConfig) bool {
	return strings.EqualFold(strings.TrimSpace(admin.InitialEmail), "admin@example.com") ||
		strings.TrimSpace(admin.InitialPassword) == "ChangeMeAdmin123!" ||
		strings.TrimSpace(admin.InitialName) == "Template Admin"
}

func envString(lookup func(string) (string, bool), key string) (string, bool) {
	value, ok := lookup(key)
	if !ok {
		return "", false
	}
	trimmed := strings.TrimSpace(value)
	return trimmed, trimmed != ""
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
