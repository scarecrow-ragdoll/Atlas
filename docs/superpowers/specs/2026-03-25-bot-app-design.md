# Bot App Design Spec

**Status:** Approved
**Date:** 2026-03-25

## Goal

Add `apps/bot` (Telegram bot, Go) to the monorepo template. The bot follows the same Clean Architecture pattern as `apps/api`, shares infrastructure via `libs/go/` (config, logger), and integrates with Nx and Docker.

## Key Decisions

| Decision                | Choice                                                   | Rationale                                                                                  |
| ----------------------- | -------------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| Template approach       | One "fat" template (api + web + bot)                     | Unused apps don't hurt: Nx affected skips them, removal is trivial                         |
| Bot library             | `github.com/go-telegram/bot` v1                          | Stable semver, Bot API 9.5, stdlib `context.Context`, zero deps, official Telegram listing |
| Architecture            | Clean Architecture (mirrors API)                         | Uniform structure across all Go apps, easier onboarding                                    |
| Bot-API relationship    | Template supports both (autonomous bot OR shared domain) | Shared domain is opt-in: extract to `libs/go/domain/` when needed                          |
| Go module strategy      | Separate `go.mod` per app/lib + `go.work`                | Module isolation, no `replace` directives, easy local dev                                  |
| Shared libs integration | Out of scope (separate task)                             | `libs/go/config` and `libs/go/logger` are in progress by other agents                      |
| Testability             | `Sender` interface wrapping `*bot.Bot`                   | Handlers depend on interface, not concrete type — easy to mock                             |

## Scope

### In scope

1. `apps/bot/` directory structure (Clean Architecture)
2. `go.work` at monorepo root
3. Nx integration (`project.json` with build/serve/test/lint targets)
4. Docker (`bot.Dockerfile`, docker-compose.dev.yml service)
5. Bot skeleton (main.go wiring, /start, /help handlers, logging + recover middleware)
6. Tests (handler tests via Sender interface, middleware tests)

### Out of scope (separate tasks)

- Integration with `libs/go/logger` and `libs/go/config` — after their implementation is complete
- Shared domain layer — opt-in, extract when real need arises
- CI pipeline for bot — trivial addition by analogy with API

## Directory Structure

```
apps/bot/
├── cmd/bot/
│   └── main.go                  # wiring, graceful shutdown
├── internal/
│   ├── appconfig/
│   │   └── config.go            # BotConfig struct (token, poll_timeout)
│   ├── botapi/
│   │   └── sender.go            # Sender interface for testability
│   ├── ctxlog/
│   │   └── ctxlog.go            # FromContext/WithContext (interim, replaced by libs/go/logger)
│   ├── handler/
│   │   ├── start.go             # /start command
│   │   ├── start_test.go
│   │   ├── help.go              # /help command
│   │   ├── help_test.go
│   │   └── default.go           # catch-all no-op handler
│   ├── middleware/
│   │   ├── logging.go           # update logging, injects logger into ctx
│   │   ├── logging_test.go
│   │   ├── recover.go           # panic recovery
│   │   └── recover_test.go
│   ├── service/                 # business logic (empty skeleton)
│   └── repository/              # data access (empty skeleton)
├── config/
│   └── config.yml               # YAML config
├── go.mod                       # module monorepo-template/apps/bot
├── air.toml                     # hot-reload
├── .golangci.yml
└── project.json                 # Nx targets
```

## Go Workspace

```go
// go.work (monorepo root)
go 1.25.0

use (
    apps/api
    apps/bot
    libs/go/config
    libs/go/logger
)
```

All Go modules resolve locally via `go.work`. No `replace` directives in individual `go.mod` files. Build always from monorepo.

## Entry Point (`cmd/bot/main.go`)

```go
func main() {
    // 1. Config (inline Viper until shared config is ready)
    cfg := loadConfig("config/config.yml")

    // 2. Logger (inline zap until shared logger is ready)
    log := initLogger(cfg.Log)
    defer log.Sync()

    // 3. Bot
    // Middleware order matters: Recover wraps Logging, so panics inside
    // Logging are caught. Register Recover first.
    opts := []bot.Option{
        bot.WithDefaultHandler(handler.Default()),
        bot.WithMiddlewares(
            middleware.Recover(log),
            middleware.Logging(log),
        ),
    }
    b, err := bot.New(cfg.Bot.Token, opts...)
    if err != nil {
        log.Fatal("failed to create bot", zap.Error(err))
    }

    // 4. Handlers
    b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handler.Start(b))
    b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handler.Help(b))

    // 5. Graceful shutdown
    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    log.Info("bot started")
    b.Start(ctx) // blocking — returns when ctx is cancelled (unlike HTTP server goroutine pattern in apps/api)
}
```

Pattern identical to API: load config -> init logger -> init transport -> wire middleware -> wire handlers -> graceful shutdown. Note: `b.Start(ctx)` is blocking (long-polling loop), unlike the API where `httpServer.ListenAndServe` runs in a goroutine.

## Sender Interface (Testability)

```go
// internal/botapi/sender.go
package botapi

import (
    "context"
    "github.com/go-telegram/bot"
    "github.com/go-telegram/bot/models"
)

type Sender interface {
    SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error)
}
```

`*bot.Bot` satisfies this interface. Handlers accept `Sender`, not `*bot.Bot`:

```go
func Start(s botapi.Sender) bot.HandlerFunc {
    return func(ctx context.Context, b *bot.Bot, update *models.Update) {
        const op = "handler.Start"
        log := ctxlog.FromContext(ctx).With(zap.String("op", op))
        log.Debug("handling /start")

        if _, err := s.SendMessage(ctx, &bot.SendMessageParams{
            ChatID: update.Message.Chat.ID,
            Text:   "Welcome! Use /help to see available commands.",
        }); err != nil {
            log.Error("failed to send message", zap.Error(err))
        }
    }
}
```

## Middleware

### Logging middleware

```go
func Logging(log *zap.Logger) bot.MiddlewareFunc {
    return func(next bot.HandlerFunc) bot.HandlerFunc {
        return func(ctx context.Context, b *bot.Bot, update *models.Update) {
            const op = "middleware.Logging"
            start := time.Now()

            // Guard against non-message updates (callbacks, inline queries, etc.)
            // where update.Message is nil.
            var chatID int64
            if update.Message != nil {
                chatID = update.Message.Chat.ID
            }

            l := log.With(
                zap.String("op", op),
                zap.Int64("update_id", int64(update.ID)),
                zap.Int64("chat_id", chatID),
            )
            ctx = ctxlog.WithContext(ctx, l)

            next(ctx, b, update)

            l.Info("update processed",
                zap.Duration("duration", time.Since(start)),
            )
        }
    }
}
```

### Recover middleware

```go
func Recover(log *zap.Logger) bot.MiddlewareFunc {
    return func(next bot.HandlerFunc) bot.HandlerFunc {
        return func(ctx context.Context, b *bot.Bot, update *models.Update) {
            const op = "middleware.Recover"
            defer func() {
                if r := recover(); r != nil {
                    log.Error("panic recovered",
                        zap.String("op", op),
                        zap.Any("panic", r),
                        zap.Int64("update_id", int64(update.ID)),
                    )
                }
            }()
            next(ctx, b, update)
        }
    }
}
```

## Context Logger (Interim)

Until `libs/go/logger` is ready, bot uses a local `internal/ctxlog` package with the same API:

```go
// internal/ctxlog/ctxlog.go
package ctxlog

import (
    "context"
    "go.uber.org/zap"
)

type ctxKey struct{}

func WithContext(ctx context.Context, l *zap.Logger) context.Context {
    return context.WithValue(ctx, ctxKey{}, l)
}

func FromContext(ctx context.Context) *zap.Logger {
    if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
        return l
    }
    return zap.NewNop()
}
```

After migration: replace all `ctxlog` imports with `monorepo-template/libs/go/logger`, delete `internal/ctxlog/`.

## Config (Interim)

Until `libs/go/config` is ready, bot uses inline Viper loading:

```go
// internal/appconfig/config.go
type Config struct {
    Bot BotConfig `mapstructure:"bot"`
    Log LogConfig `mapstructure:"log"`
}

type BotConfig struct {
    Token       string        `mapstructure:"token"`
    PollTimeout time.Duration `mapstructure:"poll_timeout"`
}

// Validate checks required fields (mirrors apps/api pattern).
func (c *Config) Validate() error {
    if c.Bot.Token == "" {
        return fmt.Errorf("bot.token is required (set BOT_TOKEN env var)")
    }
    return nil
}

type LogConfig struct {
    Level  string `mapstructure:"level"`
    Format string `mapstructure:"format"`
}
```

```yaml
# config/config.yml
# BOT_TOKEN is provided via env var, bound explicitly in loadConfig().
bot:
  token: ''
  poll_timeout: 10s

log:
  level: debug
  format: console
```

The interim `loadConfig` uses explicit `BindEnv` calls (matching the current API pattern):

```go
func loadConfig(path string) *Config {
    v := viper.New()
    v.SetConfigFile(path)
    if err := v.ReadInConfig(); err != nil {
        panic(fmt.Errorf("read config: %w", err))
    }

    // Explicit env var bindings (same pattern as apps/api)
    for _, b := range []struct{ key, env string }{
        {"bot.token", "BOT_TOKEN"},
        {"log.level", "LOG_LEVEL"},
        {"log.format", "LOG_FORMAT"},
    } {
        if err := v.BindEnv(b.key, b.env); err != nil {
            panic(fmt.Errorf("bind env %s: %w", b.key, err))
        }
    }

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        panic(fmt.Errorf("unmarshal config: %w", err))
    }

    if err := cfg.Validate(); err != nil {
        panic(err)
    }
    return &cfg
}
```

After migration to shared config:

- `LogConfig` -> `config.LogConfig` (from `libs/go/config`)
- `Load` -> `config.Load[appconfig.Config](config.Options{...})`
- Add `Postgres config.PostgresConfig` if bot needs DB

## Nx Integration

```json
{
  "name": "bot",
  "$schema": "../../node_modules/nx/schemas/project-schema.json",
  "sourceRoot": "apps/bot",
  "projectType": "application",
  "tags": ["scope:bot"],
  "targets": {
    "build": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd apps/bot && go build -o ../../dist/apps/bot ./cmd/bot"
      }
    },
    "serve": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd apps/bot && air -c air.toml"
      }
    },
    "test": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd apps/bot && go test -coverprofile=coverage.out ./..."
      }
    },
    "lint": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd apps/bot && golangci-lint run"
      }
    }
  }
}
```

## Docker

**`docker/bot.Dockerfile`** — multi-stage build:

- Stage 1: `golang:1.25-alpine` — copy `go.work` and **all modules it references** (`apps/api/go.mod`, `apps/bot/`, `libs/go/`), build binary. All modules listed in `go.work` must be present even if only `apps/bot` is built, because Go resolves the full workspace.
- Stage 2: `alpine` — copy binary + `config.yml`

**`docker-compose.dev.yml`** — new `bot` service alongside `api` and `web`. DB/Redis dependencies are optional, added when bot starts using them.

## Testing Strategy

### Middleware tests (no mock needed)

```go
func TestRecover_CatchesPanic(t *testing.T) {
    core, logs := observer.New(zap.DebugLevel)
    log := zap.New(core)

    panicking := func(ctx context.Context, b *bot.Bot, update *models.Update) {
        panic("test panic")
    }

    wrapped := middleware.Recover(log)(panicking)

    assert.NotPanics(t, func() {
        wrapped(context.Background(), nil, &models.Update{ID: 1})
    })

    require.Equal(t, 1, logs.Len())
    assert.Equal(t, "panic recovered", logs.All()[0].Message)
}
```

### Handler tests (via Sender mock)

```go
type mockSender struct {
    lastParams *bot.SendMessageParams
}

func (m *mockSender) SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error) {
    m.lastParams = params
    return &models.Message{}, nil
}

func TestStart_SendsWelcome(t *testing.T) {
    s := &mockSender{}
    h := handler.Start(s)

    update := &models.Update{
        Message: &models.Message{
            Chat: models.Chat{ID: 123},
        },
    }

    h(context.Background(), nil, update)

    require.NotNil(t, s.lastParams)
    assert.Equal(t, int64(123), s.lastParams.ChatID)
    assert.Contains(t, s.lastParams.Text, "Welcome")
}
```

## Future Migration Path

When `libs/go/config` and `libs/go/logger` are complete:

1. Replace inline Viper loading with `config.Load[appconfig.Config](config.Options{...})`
2. Replace inline zap init with `logger.New(logger.Config{...})`
3. Use `logger.WithContext` / `logger.FromContext` consistently (already designed into middleware)
4. `LogConfig` struct replaced by `config.LogConfig` from shared package
5. Add shared blocks (`PostgresConfig`, `RedisConfig`) to `appconfig.Config` as needed
