# Shared Config Package Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extract config loading into a reusable shared Go package at `libs/go/config/` with generic `Load[T]`, reusable blocks, .env support, and struct-tag validation.

**Architecture:** Single Go module at `libs/go/config/` exporting a generic `Load[T any](Options) (T, error)` function. Uses Viper for YAML + env, godotenv for .env files, go-playground/validator for struct-tag validation. Connected to `apps/api` via `go.work`. Reusable config blocks (Postgres, Redis, Log, Server) are composable structs that any service can embed.

**Tech Stack:** Go 1.25, Viper v1.18, godotenv, go-playground/validator/v10, testify

**Spec:** `docs/superpowers/specs/2026-03-25-shared-config-design.md`

---

## File Map

### New files (libs/go/config/)

| File                                  | Responsibility                                                                        |
| ------------------------------------- | ------------------------------------------------------------------------------------- |
| `libs/go/config/go.mod`               | Go module definition + dependencies                                                   |
| `libs/go/config/blocks.go`            | Reusable config structs: PostgresConfig, RedisConfig, LogConfig, ServerConfig + DSN() |
| `libs/go/config/validate.go`          | `Validate(v any) error` — wrapper over validator/v10, human-readable errors           |
| `libs/go/config/env.go`               | `loadDotEnv(path string) error` — .env file loading via godotenv                      |
| `libs/go/config/config.go`            | `Options`, `Load[T]`, sentinel errors, Viper orchestration                            |
| `libs/go/config/config_test.go`       | Full test suite: loader, validation, priority, .env, partial structs                  |
| `libs/go/config/testdata/valid.yml`   | YAML fixture for tests                                                                |
| `libs/go/config/testdata/minimal.yml` | Minimal YAML fixture (partial struct test)                                            |
| `libs/go/config/testdata/.env.test`   | .env fixture for tests                                                                |
| `libs/go/config/project.json`         | Nx project config with test + lint targets                                            |

### New files (monorepo root)

| File      | Responsibility                                 |
| --------- | ---------------------------------------------- |
| `go.work` | Go workspace linking apps/api + libs/go/config |

### New files (apps/api migration)

| File                                         | Responsibility                                                            |
| -------------------------------------------- | ------------------------------------------------------------------------- |
| `apps/api/internal/appconfig/config.go`      | App-specific Config struct (AuthConfig, PaginationConfig + shared blocks) |
| `apps/api/internal/appconfig/config_test.go` | App-specific validation tests (AuthConfig required, Pagination defaults)  |

### Modified files

| File                                                | Change                                                                   |
| --------------------------------------------------- | ------------------------------------------------------------------------ |
| `apps/api/cmd/server/main.go`                       | Switch from `config.Load` to `config.Load[appconfig.Config]`             |
| `apps/api/internal/repository/postgres/postgres.go` | Import `config` from shared package instead of `internal/config`         |
| `apps/api/internal/repository/redis/cache.go`       | Import `config` from shared package instead of `internal/config`         |
| `apps/api/go.mod`                                   | Add `require monorepo-template/libs/go/config`                           |
| `docker/docker-compose.yml`                         | Rename env vars (API_PORT→SERVER_PORT, JWT_SECRET→AUTH_JWT_SECRET, etc.) |
| `.env.example`                                      | Rename env vars to match new mapping                                     |

### Deleted files

| File                                      | Reason                        |
| ----------------------------------------- | ----------------------------- |
| `apps/api/internal/config/config.go`      | Replaced by shared package    |
| `apps/api/internal/config/config_test.go` | Tests moved to shared package |

---

## Task 1: Go workspace + shared module scaffold

**Files:**

- Create: `go.work`
- Create: `libs/go/config/go.mod`
- Create: `libs/go/config/project.json`

- [ ] **Step 1: Create go.work at monorepo root**

```go
go 1.25.0

use (
	./apps/api
	./libs/go/config
)
```

- [ ] **Step 2: Create libs/go/config/go.mod**

```bash
mkdir -p libs/go/config
```

```go
module monorepo-template/libs/go/config

go 1.25.0

require (
	github.com/go-playground/validator/v10 v10.25.0
	github.com/joho/godotenv v1.5.1
	github.com/spf13/viper v1.18.0
	github.com/stretchr/testify v1.9.0
)
```

- [ ] **Step 3: Run go mod tidy to resolve dependencies**

Run: `cd libs/go/config && go mod tidy`
Expected: go.sum generated, no errors

- [ ] **Step 4: Create Nx project.json**

Create `libs/go/config/project.json`:

```json
{
  "name": "go-config",
  "$schema": "../../../node_modules/nx/schemas/project-schema.json",
  "sourceRoot": "libs/go/config",
  "projectType": "library",
  "tags": ["scope:shared", "lang:go"],
  "targets": {
    "test": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd libs/go/config && go test -v -coverprofile=coverage.out ./..."
      }
    },
    "lint": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd libs/go/config && golangci-lint run"
      }
    }
  }
}
```

- [ ] **Step 5: Verify workspace resolves**

Run: `go work sync`
Expected: no errors

- [ ] **Step 6: Commit**

```bash
git add go.work libs/go/config/go.mod libs/go/config/go.sum libs/go/config/project.json
git commit -m "chore: scaffold shared config module with go.work workspace"
```

---

## Task 2: Reusable config blocks (blocks.go)

**Files:**

- Create: `libs/go/config/blocks.go`

- [ ] **Step 1: Create blocks.go with all shared structs**

Create `libs/go/config/blocks.go`:

```go
package config

import (
	"fmt"
	"time"
)

// PostgresConfig holds PostgreSQL connection settings.
type PostgresConfig struct {
	Host            string        `mapstructure:"host"              validate:"required"`
	Port            int           `mapstructure:"port"              validate:"required,gt=0"`
	User            string        `mapstructure:"user"              validate:"required"`
	Password        string        `mapstructure:"password"`
	DB              string        `mapstructure:"db"                validate:"required"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxConns        int32         `mapstructure:"max_conns"`
	MinConns        int32         `mapstructure:"min_conns"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

// DSN returns a PostgreSQL connection string.
func (c PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DB, c.SSLMode,
	)
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Host     string `mapstructure:"host"     validate:"required"`
	Port     int    `mapstructure:"port"     validate:"required,gt=0"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level  string `mapstructure:"level"  validate:"omitempty,oneof=debug info warn error"`
	Format string `mapstructure:"format" validate:"omitempty,oneof=json text"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port            int           `mapstructure:"port"             validate:"required,gt=0"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	Env             string        `mapstructure:"env"              validate:"omitempty,oneof=development staging production"`
	CORSOrigins     []string      `mapstructure:"cors_origins"`
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd libs/go/config && go build ./...`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add libs/go/config/blocks.go
git commit -m "feat(config): add reusable config blocks (Postgres, Redis, Log, Server)"
```

---

## Task 3: Validation wrapper (validate.go)

**Files:**

- Create: `libs/go/config/validate.go`

- [ ] **Step 1: Create validate.go**

Create `libs/go/config/validate.go`:

```go
package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// Validate runs struct-tag validation on v and returns a human-readable error.
func Validate(v any) error {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}

	msgs := make([]string, 0, len(validationErrors))
	for _, fe := range validationErrors {
		msgs = append(msgs, formatFieldError(fe))
	}

	return fmt.Errorf("%w: %s", ErrValidation, strings.Join(msgs, "; "))
}

func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s failed on '%s' validation", fe.Field(), fe.Tag())
	}
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd libs/go/config && go build ./...`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add libs/go/config/validate.go
git commit -m "feat(config): add struct-tag validation wrapper"
```

---

## Task 4: Generic loader + dotenv + errors (config.go + env.go)

**Files:**

- Create: `libs/go/config/config.go`
- Create: `libs/go/config/env.go`

- [ ] **Step 1: Create env.go**

Create `libs/go/config/env.go`:

```go
package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

// loadDotEnv reads a .env file into the OS environment.
// Uses godotenv.Load (not Overload) so pre-existing env vars take precedence.
func loadDotEnv(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("%w: %v", ErrEnvFileLoad, err)
	}
	return nil
}
```

- [ ] **Step 2: Create config.go with Options, errors, and Load[T]**

Create `libs/go/config/config.go`:

```go
package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var (
	// ErrEnvFileLoad indicates .env file could not be loaded.
	ErrEnvFileLoad = errors.New("failed to load .env file")

	// ErrConfigLoad indicates config file could not be read or parsed.
	ErrConfigLoad = errors.New("failed to load config file")

	// ErrValidation indicates config validation failed.
	ErrValidation = errors.New("config validation failed")
)

// Options configures the loading behavior.
type Options struct {
	// ConfigPath is the path to the YAML config file (required).
	ConfigPath string

	// EnvFile is the path to a .env file (optional, "" = skip).
	EnvFile string

	// EnvPrefix is an optional prefix for env var lookup (e.g. "API" → API_POSTGRES_HOST).
	EnvPrefix string
}

// Load reads configuration from YAML + .env + env vars into T, then validates.
// Priority: real env vars > .env file > YAML defaults.
func Load[T any](opts Options) (T, error) {
	var zero T

	// Step 1: Load .env file (if specified).
	// godotenv.Load does NOT overwrite pre-existing env vars.
	if opts.EnvFile != "" {
		if err := loadDotEnv(opts.EnvFile); err != nil {
			return zero, err
		}
	}

	// Step 2: Read YAML config file.
	v := viper.New()
	v.SetConfigFile(opts.ConfigPath)

	if err := v.ReadInConfig(); err != nil {
		return zero, fmt.Errorf("%w: %v", ErrConfigLoad, err)
	}

	// Step 3: Enable automatic env var resolution.
	// Replaces "." with "_" so postgres.host → POSTGRES_HOST.
	if opts.EnvPrefix != "" {
		v.SetEnvPrefix(opts.EnvPrefix)
	}
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Step 4: Unmarshal into T.
	var cfg T
	if err := v.Unmarshal(&cfg); err != nil {
		return zero, fmt.Errorf("%w: %v", ErrConfigLoad, err)
	}

	// Step 5: Validate struct tags.
	if err := Validate(cfg); err != nil {
		return zero, err
	}

	return cfg, nil
}
```

- [ ] **Step 3: Verify full package compiles**

Run: `cd libs/go/config && go build ./...`
Expected: no errors — all files reference each other correctly

- [ ] **Step 4: Commit**

```bash
git add libs/go/config/config.go libs/go/config/env.go
git commit -m "feat(config): add generic Load[T] with YAML + .env + env var support"
```

---

## Task 5: Test suite

**Files:**

- Create: `libs/go/config/config_test.go`
- Create: `libs/go/config/testdata/valid.yml`
- Create: `libs/go/config/testdata/minimal.yml`
- Create: `libs/go/config/testdata/.env.test`

- [ ] **Step 1: Create test fixtures**

Create `libs/go/config/testdata/valid.yml`:

```yaml
server:
  port: 8080
  read_timeout: 10s
  write_timeout: 30s
  shutdown_timeout: 5s
  env: development
  cors_origins:
    - 'http://localhost:3000'

log:
  level: info
  format: json

postgres:
  host: localhost
  port: 5432
  user: testuser
  password: testpass
  db: testdb
  sslmode: disable
  max_conns: 10
  min_conns: 2
  max_conn_idle_time: 30m

redis:
  host: localhost
  port: 6379

auth:
  jwt_secret: test-secret
```

Create `libs/go/config/testdata/minimal.yml`:

```yaml
postgres:
  host: localhost
  port: 5432
  user: app
  password: pass
  db: mydb
  sslmode: disable
```

Create `libs/go/config/testdata/.env.test`:

```env
POSTGRES_HOST=envhost
POSTGRES_PORT=9999
```

- [ ] **Step 2: Create config_test.go**

Create `libs/go/config/config_test.go`:

```go
package config_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/libs/go/config"
)

// fullConfig mirrors a typical app config with shared + custom blocks.
type fullConfig struct {
	Server   config.ServerConfig   `mapstructure:"server"`
	Log      config.LogConfig      `mapstructure:"log"`
	Postgres config.PostgresConfig `mapstructure:"postgres"`
	Redis    config.RedisConfig    `mapstructure:"redis"`
	Auth     authConfig            `mapstructure:"auth"`
}

type authConfig struct {
	JWTSecret string `mapstructure:"jwt_secret" validate:"required"`
}

// minimalConfig uses only one shared block.
type minimalConfig struct {
	Postgres config.PostgresConfig `mapstructure:"postgres"`
}

func TestLoad_YAMLParsing(t *testing.T) {
	cfg, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/valid.yml",
	})
	require.NoError(t, err)

	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "development", cfg.Server.Env)
	assert.Equal(t, []string{"http://localhost:3000"}, cfg.Server.CORSOrigins)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, 5432, cfg.Postgres.Port)
	assert.Equal(t, "testuser", cfg.Postgres.User)
	assert.Equal(t, "testdb", cfg.Postgres.DB)
	assert.Equal(t, "localhost", cfg.Redis.Host)
	assert.Equal(t, 6379, cfg.Redis.Port)
	assert.Equal(t, "test-secret", cfg.Auth.JWTSecret)
}

func TestLoad_EnvVarOverridesYAML(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "override-host")
	t.Setenv("SERVER_PORT", "9090")

	cfg, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/valid.yml",
	})
	require.NoError(t, err)

	assert.Equal(t, "override-host", cfg.Postgres.Host)
	assert.Equal(t, 9090, cfg.Server.Port)
}

func TestLoad_DotEnvFilePickedUp(t *testing.T) {
	// godotenv.Load writes to OS env — clean up vars it sets to avoid polluting other tests.
	t.Cleanup(func() {
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
	})

	cfg, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/valid.yml",
		EnvFile:    "testdata/.env.test",
	})
	require.NoError(t, err)

	assert.Equal(t, "envhost", cfg.Postgres.Host)
	assert.Equal(t, 9999, cfg.Postgres.Port)
}

func TestLoad_PriorityOrder_EnvWins(t *testing.T) {
	// Real env var should beat .env file value.
	// t.Setenv sets POSTGRES_HOST before godotenv.Load runs;
	// godotenv.Load (not Overload) skips already-set vars → real env wins.
	t.Setenv("POSTGRES_HOST", "real-env-host")
	t.Cleanup(func() {
		os.Unsetenv("POSTGRES_PORT") // .env.test also sets POSTGRES_PORT
	})

	cfg, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/valid.yml",
		EnvFile:    "testdata/.env.test", // .env has POSTGRES_HOST=envhost
	})
	require.NoError(t, err)

	assert.Equal(t, "real-env-host", cfg.Postgres.Host)
}

func TestLoad_ValidationError(t *testing.T) {
	// minimal.yml has no auth section, so Auth.JWTSecret will be "".
	// Explicitly set env to empty to prevent any leaked value from prior tests
	// from accidentally making this pass.
	t.Setenv("AUTH_JWT_SECRET", "")

	_, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/minimal.yml",
	})
	require.Error(t, err)
	assert.True(t, errors.Is(err, config.ErrValidation))
}

func TestLoad_MissingYAML(t *testing.T) {
	_, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/nonexistent.yml",
	})
	require.Error(t, err)
	assert.True(t, errors.Is(err, config.ErrConfigLoad))
}

func TestLoad_PartialStruct(t *testing.T) {
	// minimalConfig has only Postgres — omitempty oneof fields should pass
	cfg, err := config.Load[minimalConfig](config.Options{
		ConfigPath: "testdata/minimal.yml",
	})
	require.NoError(t, err)

	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, 5432, cfg.Postgres.Port)
}

func TestLoad_EnvPrefix(t *testing.T) {
	t.Setenv("MYAPP_POSTGRES_HOST", "prefix-host")

	cfg, err := config.Load[minimalConfig](config.Options{
		ConfigPath: "testdata/minimal.yml",
		EnvPrefix:  "MYAPP",
	})
	require.NoError(t, err)

	assert.Equal(t, "prefix-host", cfg.Postgres.Host)
}

func TestLoad_AutomaticEnvNestedKeys(t *testing.T) {
	// Verify AutomaticEnv resolves nested keys without manual BindEnv
	t.Setenv("SERVER_READ_TIMEOUT", "25s")
	t.Setenv("POSTGRES_MAX_CONNS", "50")

	cfg, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/valid.yml",
	})
	require.NoError(t, err)

	// If AutomaticEnv works for nested keys, these will be overridden.
	// If not, they'll have the YAML defaults (10s, 10).
	// This test documents the actual behavior — if it fails,
	// fallback to explicit BindEnv via reflection is needed.
	t.Logf("server.read_timeout = %v (env: 25s, yaml: 10s)", cfg.Server.ReadTimeout)
	t.Logf("postgres.max_conns = %d (env: 50, yaml: 10)", cfg.Postgres.MaxConns)
}

func TestValidate_Standalone(t *testing.T) {
	// Validate can be used outside of Load
	pg := config.PostgresConfig{
		Host: "",
		Port: 0,
	}
	err := config.Validate(pg)
	require.Error(t, err)
	assert.True(t, errors.Is(err, config.ErrValidation))
	assert.Contains(t, err.Error(), "Host is required")
	assert.Contains(t, err.Error(), "Port must be greater than 0")
}

func TestPostgresConfig_DSN(t *testing.T) {
	pg := config.PostgresConfig{
		Host:     "dbhost",
		Port:     5432,
		User:     "admin",
		Password: "pass",
		DB:       "mydb",
		SSLMode:  "require",
	}
	expected := "postgres://admin:pass@dbhost:5432/mydb?sslmode=require"
	assert.Equal(t, expected, pg.DSN())
}
```

- [ ] **Step 3: Run tests**

Run: `cd libs/go/config && go test -v ./...`
Expected: all tests PASS (AutomaticEnv nested key test may log fallback — that's informational)

- [ ] **Step 4: If AutomaticEnv nested key test shows env vars NOT being picked up, implement BindEnv fallback**

Check the `TestLoad_AutomaticEnvNestedKeys` output. If `server.read_timeout` still shows 10s (not 25s), Viper's AutomaticEnv doesn't resolve nested keys on unmarshal. In that case, add a `bindAllKeys(v *viper.Viper)` helper to `config.go` that calls `v.BindEnv` for each key found by `v.AllKeys()` after reading the YAML:

```go
// Add after v.ReadInConfig() in Load[T], BEFORE SetEnvPrefix/AutomaticEnv:
// Note: when using this fallback, do NOT also call v.SetEnvPrefix — pass
// the prefix manually here to avoid double-prefixing.
for _, key := range v.AllKeys() {
	envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	if opts.EnvPrefix != "" {
		envKey = opts.EnvPrefix + "_" + envKey
	}
	_ = v.BindEnv(key, envKey)
}
// Remove the v.SetEnvPrefix() call above if using this fallback.
```

Then re-run tests to confirm nested keys are resolved.

- [ ] **Step 5: Commit**

```bash
git add libs/go/config/config_test.go libs/go/config/testdata/
git commit -m "test(config): add full test suite for shared config loader"
```

---

## Task 6: Migrate apps/api — create appconfig + update imports

**Files:**

- Create: `apps/api/internal/appconfig/config.go`
- Create: `apps/api/internal/appconfig/config_test.go`
- Modify: `apps/api/cmd/server/main.go`
- Modify: `apps/api/internal/repository/postgres/postgres.go:10-11,18`
- Modify: `apps/api/internal/repository/redis/cache.go:10-11,18`
- Delete: `apps/api/internal/config/config.go`
- Delete: `apps/api/internal/config/config_test.go`

- [ ] **Step 1: Create apps/api/internal/appconfig/config.go**

```go
package appconfig

import "monorepo-template/libs/go/config"

// AuthConfig holds authentication settings specific to the API service.
type AuthConfig struct {
	JWTSecret string `mapstructure:"jwt_secret" validate:"required"`
}

// PaginationConfig holds pagination defaults.
type PaginationConfig struct {
	DefaultPageSize int `mapstructure:"default_page_size" validate:"gt=0"`
	MaxPageSize     int `mapstructure:"max_page_size"     validate:"gt=0"`
}

// Config is the API service configuration.
type Config struct {
	Server     config.ServerConfig   `mapstructure:"server"`
	Log        config.LogConfig      `mapstructure:"log"`
	Postgres   config.PostgresConfig `mapstructure:"postgres"`
	Redis      config.RedisConfig    `mapstructure:"redis"`
	Auth       AuthConfig            `mapstructure:"auth"`
	Pagination PaginationConfig      `mapstructure:"pagination"`
}
```

- [ ] **Step 2: Create apps/api/internal/appconfig/config_test.go**

App-specific validation tests that cover the fields not tested in the shared package:

```go
package appconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/libs/go/config"
	"monorepo-template/apps/api/internal/appconfig"
)

func TestConfig_AuthJWTSecretRequired(t *testing.T) {
	cfg := appconfig.Config{
		Server:   config.ServerConfig{Port: 8080},
		Postgres: config.PostgresConfig{Host: "localhost", Port: 5432, User: "app", DB: "test"},
		Redis:    config.RedisConfig{Host: "localhost", Port: 6379},
		Auth:     appconfig.AuthConfig{JWTSecret: ""},
		Pagination: appconfig.PaginationConfig{DefaultPageSize: 20, MaxPageSize: 100},
	}
	err := config.Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "JWTSecret is required")
}

func TestConfig_PaginationValid(t *testing.T) {
	cfg := appconfig.Config{
		Server:   config.ServerConfig{Port: 8080},
		Postgres: config.PostgresConfig{Host: "localhost", Port: 5432, User: "app", DB: "test"},
		Redis:    config.RedisConfig{Host: "localhost", Port: 6379},
		Auth:     appconfig.AuthConfig{JWTSecret: "secret"},
		Pagination: appconfig.PaginationConfig{DefaultPageSize: 20, MaxPageSize: 100},
	}
	err := config.Validate(cfg)
	require.NoError(t, err)
}
```

- [ ] **Step 3: Update apps/api/cmd/server/main.go**

Replace the import block and config loading. The complete new import block:

```go
import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"monorepo-template/libs/go/config"
	"monorepo-template/apps/api/internal/appconfig"
	"monorepo-template/apps/api/internal/graph"
	healthHandler "monorepo-template/apps/api/internal/handler"
	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/apps/api/internal/repository/postgres"
	redisRepo "monorepo-template/apps/api/internal/repository/redis"
	"monorepo-template/apps/api/internal/service"
)
```

**Key change:** Remove `"monorepo-template/apps/api/internal/config"`, add `"monorepo-template/libs/go/config"` and `"monorepo-template/apps/api/internal/appconfig"`. The identifier `config` now refers to the shared library.

Replace config loading (line 26):

```go
// Old:
cfg, err := config.Load("config/config.yml")

// New:
cfg, err := config.Load[appconfig.Config](config.Options{
	ConfigPath: "config/config.yml",
	EnvFile:    ".env",
})
```

All `cfg.` references stay the same (`cfg.Server.Port`, `cfg.Postgres`, `cfg.Auth.JWTSecret`, etc.) since field names match. The only difference: `Load` returns a value not a pointer, so `cfg` is `appconfig.Config` not `*config.Config`. No dereferences exist in the current code, so no changes needed.

- [ ] **Step 4: Update postgres repository import**

In `apps/api/internal/repository/postgres/postgres.go`, change:

```go
// Old:
import "monorepo-template/apps/api/internal/config"

// New:
import "monorepo-template/libs/go/config"
```

The function signature `New(cfg config.PostgresConfig, ...)` stays identical — the type name is the same.

- [ ] **Step 5: Update redis repository import and pass DB field**

In `apps/api/internal/repository/redis/cache.go`, change:

```go
// Old:
import "monorepo-template/apps/api/internal/config"

// New:
import "monorepo-template/libs/go/config"
```

The function signature `New(cfg config.RedisConfig, ...)` stays identical. Also update `redis.NewClient` to pass the new `DB` field:

```go
rdb := redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	Password: cfg.Password,
	DB:       cfg.DB,
})
```

- [ ] **Step 6: Delete old config package**

```bash
rm -rf apps/api/internal/config/
```

- [ ] **Step 7: Update apps/api/go.mod**

Run: `cd apps/api && go mod tidy`
Expected: `monorepo-template/libs/go/config` added to require block, unused internal/config references removed.

Note: `go.work` handles local module resolution, but the `require` directive is still needed for the module graph. If `go mod tidy` doesn't add it automatically, run `cd apps/api && go get monorepo-template/libs/go/config` first.

- [ ] **Step 8: Verify compilation**

Run: `cd apps/api && go build ./cmd/server`
Expected: compiles without errors

- [ ] **Step 9: Run existing Go tests**

Run: `cd apps/api && go test ./...`
Expected: all existing tests pass (health, cors, auth middleware, service tests + new appconfig tests)

- [ ] **Step 10: Commit**

```bash
git add apps/api/ go.work
git commit -m "refactor(api): migrate to shared config package"
```

---

## Task 7: Update docker-compose + .env.example env var names

**Files:**

- Modify: `docker/docker-compose.yml:39-52`
- Modify: `.env.example`

- [ ] **Step 1: Update docker/docker-compose.yml api environment**

Change env var names in the `api` service environment section:

```yaml
environment:
  - POSTGRES_HOST=postgres
  - POSTGRES_PORT=5432
  - POSTGRES_USER=${POSTGRES_USER:-app}
  - POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-secret}
  - POSTGRES_DB=${POSTGRES_DB:-monorepo_dev}
  - POSTGRES_SSLMODE=disable
  - REDIS_HOST=redis
  - REDIS_PORT=6379
  - REDIS_PASSWORD=${REDIS_PASSWORD:-}
  - AUTH_JWT_SECRET=${AUTH_JWT_SECRET:-dev-secret}
  - SERVER_PORT=8080
  - LOG_LEVEL=debug
  - SERVER_ENV=development
```

Changes:

- `JWT_SECRET` → `AUTH_JWT_SECRET`
- `API_PORT` → `SERVER_PORT`
- `API_LOG_LEVEL` → `LOG_LEVEL`
- Added `SERVER_ENV=development`

Also update the ports mapping line if it references `API_PORT`:

```yaml
ports:
  - '${SERVER_PORT:-8080}:8080'
```

- [ ] **Step 2: Update .env.example**

```env
# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=app
POSTGRES_PASSWORD=secret
POSTGRES_DB=monorepo_dev
POSTGRES_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Go API
SERVER_PORT=8080
SERVER_ENV=development
SERVER_CORS_ORIGINS=http://localhost:3000
LOG_LEVEL=debug

# Auth
AUTH_JWT_SECRET=change-me-in-production

# Next.js (client-side)
NEXT_PUBLIC_API_URL=http://localhost:8080/graphql
NEXT_PUBLIC_APP_NAME=MonorepoApp
```

- [ ] **Step 3: Commit**

```bash
git add docker/docker-compose.yml .env.example
git commit -m "chore: rename env vars to match shared config mapping"
```

---

## Task 8: Final verification

- [ ] **Step 1: Run shared config tests**

Run: `cd libs/go/config && go test -v -count=1 ./...`
Expected: all PASS

- [ ] **Step 2: Run API tests**

Run: `cd apps/api && go test -v ./...`
Expected: all PASS

- [ ] **Step 3: Run linter on shared config**

Run: `cd libs/go/config && golangci-lint run`
Expected: no issues

- [ ] **Step 4: Run linter on API**

Run: `cd apps/api && golangci-lint run`
Expected: no issues

- [ ] **Step 5: Build API binary**

Run: `cd apps/api && go build -o /dev/null ./cmd/server`
Expected: compiles cleanly

- [ ] **Step 6: Verify Nx recognizes the new project**

Run: `npx nx show project go-config`
Expected: shows project config with test + lint targets

- [ ] **Step 7: Final commit (if any fixups needed)**

Only if previous steps required fixes:

```bash
git add -A
git commit -m "fix: address lint/test issues from shared config migration"
```
