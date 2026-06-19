# Shared Config Package Design

**Date**: 2026-03-25
**Status**: Draft
**Scope**: Generic config loader + reusable config blocks as a shared Go package

## Problem

Current config implementation lives in `apps/api/internal/config/` and is not reusable. As the monorepo grows with more Go services, each would need to duplicate config loading logic (Viper setup, YAML parsing, env var binding, validation). There is also no .env file support.

## Decision Summary

| Decision            | Choice                                                |
| ------------------- | ----------------------------------------------------- |
| Config library      | Viper (keep existing)                                 |
| Source priority     | YAML defaults → .env file → env vars                  |
| API style           | Generic loader `Load[T any](opts Options) (T, error)` |
| Reusable blocks     | PostgresConfig, RedisConfig, LogConfig, ServerConfig  |
| Validation          | go-playground/validator v10 (struct tags)             |
| Package location    | `libs/go/config/` (single Go module)                  |
| Module coordination | `go.work` in monorepo root                            |
| Backward compat     | Not required (env var names will change)              |

## Package Structure

```
libs/go/config/
├── go.mod              // module monorepo-template/libs/go/config
├── config.go           // Load[T], Options struct, private loading logic
├── validate.go         // Validate() — wrapper over validator/v10
├── blocks.go           // PostgresConfig, RedisConfig, LogConfig, ServerConfig
├── env.go              // loadDotEnv() — .env reading via godotenv
└── config_test.go      // tests for loader + validation + blocks

go.work (monorepo root)
├── apps/api
└── libs/go/config
```

Nx `project.json` for `libs/go/config` with targets: `test`, `lint`.

## API

### Options

```go
type Options struct {
    ConfigPath string   // path to YAML file (required)
    EnvFile    string   // path to .env file (optional, "" = skip)
    EnvPrefix  string   // env var prefix, e.g. "API" → API_POSTGRES_HOST (optional)
}
```

### Load

```go
func Load[T any](opts Options) (T, error)
```

Execution order:

1. Read .env file (if `EnvFile` is set) via `godotenv.Load()` — **not `Overload`**, so pre-existing OS environment variables always take precedence over .env file values
2. Read YAML file via `viper.ReadInConfig()`
3. Enable `viper.AutomaticEnv()` + `SetEnvKeyReplacer("." → "_")` + optional prefix
4. Unmarshal into `T` via `viper.Unmarshal()`
5. Validate via `validator.Struct()`
6. Return `(T, nil)` or `(zero, error)`

**Priority guarantee**: real env vars > .env file > YAML defaults. This is achieved because `godotenv.Load()` skips variables already present in the OS environment.

### Validate

```go
func Validate(v any) error
```

Called automatically inside `Load`, but also available standalone.

### Reusable Blocks

```go
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

func (c PostgresConfig) DSN() string
```

```go
type RedisConfig struct {
    Host     string `mapstructure:"host"     validate:"required"`
    Port     int    `mapstructure:"port"     validate:"required,gt=0"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}
```

```go
type LogConfig struct {
    Level  string `mapstructure:"level"  validate:"omitempty,oneof=debug info warn error"`
    Format string `mapstructure:"format" validate:"omitempty,oneof=json text"`
}
```

```go
type ServerConfig struct {
    Port            int           `mapstructure:"port"             validate:"required,gt=0"`
    ReadTimeout     time.Duration `mapstructure:"read_timeout"`
    WriteTimeout    time.Duration `mapstructure:"write_timeout"`
    ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
    Env             string        `mapstructure:"env"              validate:"omitempty,oneof=development staging production"`
    CORSOrigins     []string      `mapstructure:"cors_origins"`
}
```

## Env Var Mapping

Uses `viper.AutomaticEnv()` + `SetEnvKeyReplacer(strings.NewReplacer(".", "_"))`.

**Caveat**: Viper's `AutomaticEnv` with nested keys and `Unmarshal` can be unreliable without explicit `BindEnv` calls in some Viper versions. The implementation must include a test verifying nested keys resolve via `AutomaticEnv` without manual `BindEnv`. If this proves unreliable, fallback to explicit `BindEnv` per key using reflection over the config struct.

Nested YAML keys map automatically:

- `postgres.host` → `POSTGRES_HOST`
- `server.read_timeout` → `SERVER_READ_TIMEOUT`
- `auth.jwt_secret` → `AUTH_JWT_SECRET`

With `EnvPrefix: "API"`:

- `postgres.host` → `API_POSTGRES_HOST`

### Changes from current env var names

| Current                        | New                                 |
| ------------------------------ | ----------------------------------- |
| `API_PORT`                     | `SERVER_PORT`                       |
| `API_LOG_LEVEL`                | `LOG_LEVEL`                         |
| `APP_ENV`                      | `SERVER_ENV`                        |
| `JWT_SECRET`                   | `AUTH_JWT_SECRET`                   |
| `CORS_ORIGINS`                 | `SERVER_CORS_ORIGINS`               |
| Others (POSTGRES*\*, REDIS*\*) | Same (already match YAML structure) |

## Error Handling

Three typed sentinel errors:

```go
var (
    ErrEnvFileLoad = errors.New("failed to load .env file")
    ErrConfigLoad  = errors.New("failed to load config file")
    ErrValidation  = errors.New("config validation failed")
)
```

`ErrValidation` wraps human-readable field errors:

```
config validation failed: Postgres.Host is required; Server.Port must be greater than 0
```

Consumers check via `errors.Is(err, config.ErrValidation)`.

## Usage Example

```go
// apps/api/internal/appconfig/config.go
package appconfig

import "monorepo-template/libs/go/config"

type AuthConfig struct {
    JWTSecret string `mapstructure:"jwt_secret" validate:"required"`
}

type PaginationConfig struct {
    DefaultPageSize int `mapstructure:"default_page_size" validate:"gt=0"`
    MaxPageSize     int `mapstructure:"max_page_size"     validate:"gt=0"`
}

type Config struct {
    Server     config.ServerConfig   `mapstructure:"server"`
    Log        config.LogConfig      `mapstructure:"log"`
    Postgres   config.PostgresConfig `mapstructure:"postgres"`
    Redis      config.RedisConfig    `mapstructure:"redis"`
    Auth       AuthConfig            `mapstructure:"auth"`
    Pagination PaginationConfig      `mapstructure:"pagination"`
}
```

```go
// apps/api/cmd/server/main.go
cfg, err := config.Load[appconfig.Config](config.Options{
    ConfigPath: "config/config.yml",
    EnvFile:    ".env",
})
if err != nil {
    log.Fatal(err)
}
```

## Migration (apps/api)

0. **Create** `go.work` at monorepo root via `go work init` + `go work use apps/api libs/go/config` — enables workspace-local module resolution
1. **Delete** `apps/api/internal/config/` entirely
2. **Create** `apps/api/internal/appconfig/config.go` with app-specific struct
3. **Update** `apps/api/cmd/server/main.go` to use `config.Load[appconfig.Config]`
4. **Add** `require monorepo-template/libs/go/config` in `apps/api/go.mod` — `go.work` handles resolution of the local module, but the `require` is still needed for the module graph
5. **Update** `docker-compose*.yml` and `.env.example` — rename env vars
6. **No changes** to `config/config.yml` (YAML keys unchanged)
7. **Delete** `apps/api/internal/config/config_test.go` — generic loader tests move to shared package; app-specific validation tests (e.g. `AuthConfig.JWTSecret` required) belong in `apps/api/internal/appconfig/config_test.go`

## Testing

Tests in `libs/go/config/config_test.go` with fixtures in `testdata/`:

- YAML parsing → struct populated correctly
- Env var overrides YAML value
- .env file values picked up
- Priority order: env var > .env > YAML
- Invalid config → `ErrValidation` with field descriptions
- Missing YAML → `ErrConfigLoad`
- Partial struct (only some blocks) → maps only relevant sections, omitted `oneof` fields pass validation
- EnvPrefix → variables resolved with prefix
- AutomaticEnv nested keys → verify resolution without manual BindEnv (fallback test)
