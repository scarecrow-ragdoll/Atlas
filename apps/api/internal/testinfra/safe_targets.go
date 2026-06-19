package testinfra

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"monorepo-template/libs/go/config"
)

const (
	TestPostgresDB      = "monorepo_test"
	DefaultPostgresHost = "localhost"
	DefaultPostgresPort = "17501"
	DefaultPostgresUser = "app"
	DefaultPostgresPass = "secret"
	DefaultRedisHost    = "localhost"
	DefaultRedisPort    = "17502"
	DevPostgresPort     = "7501"
	DevRedisPort        = "7502"
)

type TestingT interface {
	Helper()
	Fatalf(format string, args ...any)
}

func CoverageGateEnabled() bool {
	return os.Getenv("COVERAGE_GATE") == "1"
}

func PostgresDSN() string {
	if dsn := os.Getenv("API_TEST_DATABASE_DSN"); strings.TrimSpace(dsn) != "" {
		return dsn
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env("TEST_POSTGRES_USER", DefaultPostgresUser),
		env("TEST_POSTGRES_PASSWORD", DefaultPostgresPass),
		env("TEST_POSTGRES_HOST", DefaultPostgresHost),
		TestPostgresPort(),
		env("TEST_POSTGRES_DB", TestPostgresDB),
	)
}

func PostgresConfig(t TestingT) config.PostgresConfig {
	t.Helper()
	port := mustPort(t, "TEST_POSTGRES_PORT", TestPostgresPort())
	cfg := config.PostgresConfig{
		Host:     env("TEST_POSTGRES_HOST", DefaultPostgresHost),
		Port:     port,
		User:     env("TEST_POSTGRES_USER", DefaultPostgresUser),
		Password: env("TEST_POSTGRES_PASSWORD", DefaultPostgresPass),
		DB:       env("TEST_POSTGRES_DB", TestPostgresDB),
		SSLMode:  "disable",
		MaxConns: 2,
		MinConns: 1,
	}
	RequireSafePostgresDSN(t, cfg.DSN())
	return cfg
}

func RedisConfig(t TestingT) config.RedisConfig {
	t.Helper()
	port := mustPort(t, "TEST_REDIS_PORT", env("TEST_REDIS_PORT", DefaultRedisPort))
	if strconv.Itoa(port) == DevRedisPort {
		t.Fatalf("TEST_REDIS_PORT %s is the development redis port", DevRedisPort)
	}
	db := mustInt(t, "TEST_REDIS_DB", env("TEST_REDIS_DB", "0"))
	return config.RedisConfig{
		Host:     env("TEST_REDIS_HOST", DefaultRedisHost),
		Port:     port,
		Password: env("TEST_REDIS_PASSWORD", ""),
		DB:       db,
	}
}

func RequireSafePostgresDSN(t TestingT, dsn string) {
	t.Helper()
	if err := ValidateSafePostgresDSN(dsn); err != nil {
		t.Fatalf("unsafe postgres test DSN: %v", err)
	}
}

func ValidateSafePostgresDSN(dsn string) error {
	if strings.TrimSpace(dsn) == "" {
		return fmt.Errorf("dsn is empty")
	}

	parsed, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("parse dsn: %w", err)
	}
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return fmt.Errorf("scheme %q is not postgres", parsed.Scheme)
	}

	db := strings.TrimPrefix(parsed.Path, "/")
	if db == "monorepo_dev" {
		return fmt.Errorf("database %q is the development postgres database", db)
	}
	expectedDB := env("TEST_POSTGRES_DB", TestPostgresDB)
	if db != expectedDB {
		return fmt.Errorf("database %q is not %q", db, expectedDB)
	}

	port := parsed.Port()
	if port == "" {
		return fmt.Errorf("dsn must include a postgres port")
	}
	if port == DevPostgresPort {
		return fmt.Errorf("port %s is the development postgres port", DevPostgresPort)
	}
	if port != TestPostgresPort() {
		return fmt.Errorf("port %s is not the configured test postgres port %s", port, TestPostgresPort())
	}

	return nil
}

func TestPostgresPort() string {
	return env("TEST_POSTGRES_PORT", DefaultPostgresPort)
}

func mustPort(t TestingT, key string, value string) int {
	t.Helper()
	port := mustInt(t, key, value)
	if port <= 0 {
		t.Fatalf("%s must be greater than zero, got %d", key, port)
	}
	return port
}

func mustInt(t TestingT, key string, value string) int {
	t.Helper()
	parsed, err := strconv.Atoi(value)
	if err != nil {
		t.Fatalf("%s must be an integer, got %q", key, value)
	}
	return parsed
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
