# Shared Logger Package Design

**Date:** 2026-03-25
**Status:** Approved
**Scope:** Shared logging package for Go services in the monorepo

## Overview

Extract logging into a shared Go package `libs/go/logger` that provides unified Zap-based structured logging with request-id correlation and consistent `op` tracing pattern across all application layers.

## Decisions

| Decision           | Choice                               | Rationale                                                                  |
| ------------------ | ------------------------------------ | -------------------------------------------------------------------------- |
| Library            | Zap v1.27.0                          | Already integrated, zero-allocation, production-proven                     |
| Package location   | `libs/go/logger`                     | `libs/` is the shared code convention; `go/` separates Go libs from others |
| Package structure  | Monolith (single package)            | Simple, one import path; split later if gRPC needed                        |
| Logger propagation | `context.Context`                    | `FromContext(ctx)` as the single API for all layers                        |
| Request-ID         | Middleware-generated UUID v4         | Respects incoming `X-Request-ID`, falls back to generated                  |
| Operation tracing  | `const op = "StructName.MethodName"` | Logs get `op` field, errors wrapped with `fmt.Errorf("%s: %w", op, err)`   |
| Fallback logger    | `zap.NewNop()`                       | Safe default when ctx has no logger (tests, background jobs)               |

## Package Structure

```
libs/go/logger/
├── go.mod             # module monorepo-template/libs/go/logger
├── logger.go          # New(), Config
├── context.go         # FromContext(), WithContext()
├── request_id.go      # UUID generation/extraction, header constant
├── middleware.go       # RequestID(), Logging()
├── logger_test.go     # New() tests
├── context_test.go    # WithContext/FromContext tests
└── middleware_test.go  # RequestID/Logging middleware tests
```

## Public API

### logger.go

```go
// Config holds logger configuration.
type Config struct {
    Level  string // "debug", "info", "warn", "error"; default "info"
    Format string // "json" or "console"; default "json"
}

// New creates a configured *zap.Logger.
// "json" format uses zap.NewProductionConfig(), "console" uses zap.NewDevelopmentConfig().
// Returns error if level string is invalid.
func New(cfg Config) (*zap.Logger, error)
```

### context.go

```go
// WithContext returns a new context with the given logger.
func WithContext(ctx context.Context, l *zap.Logger) context.Context

// FromContext extracts the logger from context.
// Returns zap.NewNop() if no logger is present.
func FromContext(ctx context.Context) *zap.Logger
```

### middleware.go

```go
// RequestID middleware:
// 1. Extracts X-Request-ID from request header, or generates UUID v4
// 2. Sets X-Request-ID in response header
// 3. Creates base.With(zap.String("request_id", id)) and puts it in ctx via WithContext
func RequestID(base *zap.Logger) func(http.Handler) http.Handler

// Logging middleware:
// Logs each request with method, path, status, duration, remote_addr.
// Takes logger from ctx (set by RequestID middleware).
// Wraps http.ResponseWriter internally to capture the written status code.
// ORDERING: RequestID must be registered before Logging in the middleware chain.
// If no logger is found in ctx, falls back to zap.NewNop() (silent, no panic).
func Logging() func(http.Handler) http.Handler
```

### request_id.go

```go
// HeaderRequestID is the canonical header name.
const HeaderRequestID = "X-Request-ID"

// generateID and extractID are unexported helpers.
// generateID returns a new UUID v4 string.
// extractID gets the request ID from the request header, or returns "".
```

## Usage Pattern (all layers)

```go
func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) {
    const op = "UserService.GetByID"
    log := logger.FromContext(ctx).With(zap.String("op", op))

    log.Debug("getting user", zap.String("user_id", id))

    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    return user, nil
}
```

**Log output:**

```json
{
  "level": "debug",
  "op": "UserService.GetByID",
  "request_id": "abc-123",
  "msg": "getting user",
  "user_id": "42"
}
```

**Error chain:**

```
UserService.GetByID: UserRepo.FindByID: no rows in result set
```

## Middleware Chain

```go
l, err := logger.New(logger.Config{
    Level:  cfg.Log.Level,
    Format: cfg.Log.Format,
})
defer func() { _ = l.Sync() }()

r := chi.NewRouter()
r.Use(logger.RequestID(l))  // generates request-id, puts logger in ctx
r.Use(logger.Logging())      // logs request/response using logger from ctx
```

Order matters: `RequestID` must come before `Logging`.

## Migration (apps/api)

1. **Add dependency:** `apps/api/go.mod` requires `monorepo-template/libs/go/logger` via `replace` directive for local development (no `go.work` required)
2. **Remove:** `internal/middleware/logging.go` (replaced by `logger.Logging()`)
3. **Replace in main.go:** Zap initialization with `logger.New()`, middleware with `logger.RequestID(l)` + `logger.Logging()`
4. **Add to services/repositories/resolvers:** `logger.FromContext(ctx)` + `const op` pattern
5. **Config:** `LogConfig` in `internal/config/` stays, maps to `logger.Config`
6. **Out of scope:** `middleware.Auth` continues receiving `*zap.Logger` via DI — migration to `FromContext` is a separate follow-up

## Testing

| Test file            | Coverage                                                                                                                                      |
| -------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| `logger_test.go`     | `New()` with json/console formats, valid/invalid levels                                                                                       |
| `context_test.go`    | `WithContext/FromContext` round-trip, nop fallback on empty ctx                                                                               |
| `middleware_test.go` | `RequestID`: generates UUID when no header, passes through existing header, sets response header. `Logging`: logs method/path/status/duration |

Testing approach: `zap/zaptest/observer` for asserting logged fields and values.

## Dependencies

- `go.uber.org/zap v1.27.0`
- `github.com/google/uuid v1.6.0`
