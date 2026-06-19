# Shared Logger Package Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create a shared Go logging package at `libs/go/logger` (Zap-based) with context propagation, request-id correlation, and HTTP middleware, then migrate `apps/api` to use it.

**Architecture:** Single Go package providing `New()` (logger factory), `FromContext()`/`WithContext()` (context propagation), `RequestID()` (middleware that generates/extracts request-id and injects logger into ctx), `Logging()` (HTTP request logging middleware). The `apps/api` replaces its inline Zap initialization and `internal/middleware/logging.go` with the shared package.

**Tech Stack:** Go 1.25, go.uber.org/zap v1.27.0, github.com/google/uuid v1.6.0, testify v1.9

---

## File Map

### New files (libs/go/logger/)

| File                                | Responsibility                                                    |
| ----------------------------------- | ----------------------------------------------------------------- |
| `libs/go/logger/go.mod`             | Go module: `monorepo-template/libs/go/logger`                     |
| `libs/go/logger/logger.go`          | `Config` struct, `New(Config) (*zap.Logger, error)`               |
| `libs/go/logger/context.go`         | `FromContext(ctx)`, `WithContext(ctx, l)`                         |
| `libs/go/logger/request_id.go`      | `HeaderRequestID` const, unexported `generateID()`, `extractID()` |
| `libs/go/logger/middleware.go`      | `RequestID(base)`, `Logging()`, unexported `responseWriter`       |
| `libs/go/logger/logger_test.go`     | Tests for `New()`                                                 |
| `libs/go/logger/context_test.go`    | Tests for `FromContext`/`WithContext`                             |
| `libs/go/logger/middleware_test.go` | Tests for `RequestID` and `Logging` middlewares                   |

### Modified files (apps/api/)

| File                                                 | Change                                                                       |
| ---------------------------------------------------- | ---------------------------------------------------------------------------- |
| `apps/api/go.mod`                                    | Add `require` + `replace` directive for `monorepo-template/libs/go/logger`   |
| `apps/api/cmd/server/main.go`                        | Replace inline Zap init with `logger.New()`, replace middleware registration |
| `apps/api/internal/middleware/logging.go`            | **Delete** (replaced by `logger.Logging()`)                                  |
| `apps/api/internal/service/user_service.go`          | Add `logger.FromContext(ctx)` + `const op` pattern                           |
| `apps/api/internal/repository/postgres/user_repo.go` | Add `logger.FromContext(ctx)` + `const op` pattern                           |
| `apps/api/internal/graph/schema.resolvers.go`        | Add `logger.FromContext(ctx)` + `const op` pattern to implemented resolvers  |

---

## Task 1: Create logger package — go.mod and `New()`

**Files:**

- Create: `libs/go/logger/go.mod`
- Create: `libs/go/logger/logger.go`
- Create: `libs/go/logger/logger_test.go`

- [ ] **Step 1: Create `go.mod`**

```
libs/go/logger/go.mod
```

```
module monorepo-template/libs/go/logger

go 1.25.0

require (
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.9.0
	go.uber.org/zap v1.27.0
)
```

Run: `cd libs/go/logger && go mod tidy`

- [ ] **Step 2: Write failing tests for `New()`**

Create `libs/go/logger/logger_test.go`:

```go
package logger_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/libs/go/logger"
)

func TestNew_JSONFormat(t *testing.T) {
	l, err := logger.New(logger.Config{Level: "info", Format: "json"})
	require.NoError(t, err)
	assert.NotNil(t, l)
}

func TestNew_ConsoleFormat(t *testing.T) {
	l, err := logger.New(logger.Config{Level: "debug", Format: "console"})
	require.NoError(t, err)
	assert.NotNil(t, l)
}

func TestNew_DefaultsToJSONInfo(t *testing.T) {
	l, err := logger.New(logger.Config{})
	require.NoError(t, err)
	assert.NotNil(t, l)
}

func TestNew_InvalidLevel_ReturnsError(t *testing.T) {
	_, err := logger.New(logger.Config{Level: "invalid_level"})
	assert.Error(t, err)
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd libs/go/logger && go test -v ./...`
Expected: Compilation error — `logger.New` and `logger.Config` not found.

- [ ] **Step 4: Implement `logger.go`**

Create `libs/go/logger/logger.go`:

```go
package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Config holds logger configuration.
type Config struct {
	Level  string // "debug", "info", "warn", "error"; default "info"
	Format string // "json" or "console"; default "json"
}

// New creates a configured *zap.Logger.
// "json" format uses zap.NewProductionConfig(), "console" uses zap.NewDevelopmentConfig().
// Returns error if level string is invalid.
func New(cfg Config) (*zap.Logger, error) {
	if cfg.Level == "" {
		cfg.Level = "info"
	}

	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("parse log level %q: %w", cfg.Level, err)
	}

	var zapCfg zap.Config
	if cfg.Format == "console" {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}
	zapCfg.Level = level

	l, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return l, nil
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd libs/go/logger && go test -v ./...`
Expected: All 4 tests PASS.

- [ ] **Step 6: Commit**

```bash
git add libs/go/logger/go.mod libs/go/logger/go.sum libs/go/logger/logger.go libs/go/logger/logger_test.go
git commit -m "feat(logger): add shared logger package with New() and Config"
```

---

## Task 2: Context propagation — `FromContext` / `WithContext`

**Files:**

- Create: `libs/go/logger/context.go`
- Create: `libs/go/logger/context_test.go`

- [ ] **Step 1: Write failing tests**

Create `libs/go/logger/context_test.go`:

```go
package logger_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"monorepo-template/libs/go/logger"
)

func TestWithContext_FromContext_RoundTrip(t *testing.T) {
	core, _ := observer.New(zap.DebugLevel)
	l := zap.New(core)

	ctx := logger.WithContext(context.Background(), l)
	got := logger.FromContext(ctx)

	assert.Equal(t, l, got)
}

func TestFromContext_EmptyCtx_ReturnsNop(t *testing.T) {
	got := logger.FromContext(context.Background())
	assert.NotNil(t, got)
	// nop logger should not panic on usage
	got.Info("should not panic")
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd libs/go/logger && go test -v -run TestWithContext -run TestFromContext ./...`
Expected: Compilation error — `logger.WithContext` and `logger.FromContext` not found.

- [ ] **Step 3: Implement `context.go`**

Create `libs/go/logger/context.go`:

```go
package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey struct{}

// WithContext returns a new context with the given logger.
func WithContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext extracts the logger from context.
// Returns zap.NewNop() if no logger is present.
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	}
	return zap.NewNop()
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd libs/go/logger && go test -v ./...`
Expected: All 6 tests PASS (4 logger + 2 context).

- [ ] **Step 5: Commit**

```bash
git add libs/go/logger/context.go libs/go/logger/context_test.go
git commit -m "feat(logger): add context propagation (FromContext/WithContext)"
```

---

## Task 3: Request-ID and HTTP middleware

**Files:**

- Create: `libs/go/logger/request_id.go`
- Create: `libs/go/logger/middleware.go`
- Create: `libs/go/logger/middleware_test.go`

- [ ] **Step 1: Write failing tests for RequestID middleware**

Create `libs/go/logger/middleware_test.go`:

```go
package logger_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"monorepo-template/libs/go/logger"
)

func TestRequestID_GeneratesUUID_WhenNoHeader(t *testing.T) {
	core, _ := observer.New(zap.DebugLevel)
	base := zap.New(core)

	var capturedID string
	handler := logger.RequestID(base)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = w.Header().Get("X-Request-ID")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.NotEmpty(t, capturedID)
	assert.Len(t, capturedID, 36) // UUID v4 format
	assert.Equal(t, capturedID, rec.Header().Get("X-Request-ID"))
}

func TestRequestID_PassesThrough_ExistingHeader(t *testing.T) {
	core, _ := observer.New(zap.DebugLevel)
	base := zap.New(core)

	var capturedID string
	handler := logger.RequestID(base)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = w.Header().Get("X-Request-ID")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "existing-id-123")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, "existing-id-123", capturedID)
	assert.Equal(t, "existing-id-123", rec.Header().Get("X-Request-ID"))
}

func TestRequestID_PutsLoggerWithRequestID_InContext(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	base := zap.New(core)

	handler := logger.RequestID(base)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		l.Info("test message")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "test-req-id")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	assert.Equal(t, "test message", entry.Message)
	fields := entry.ContextMap()
	assert.Equal(t, "test-req-id", fields["request_id"])
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd libs/go/logger && go test -v -run TestRequestID ./...`
Expected: Compilation error — `logger.RequestID` not found.

- [ ] **Step 3: Implement `request_id.go`**

Create `libs/go/logger/request_id.go`:

```go
package logger

import "github.com/google/uuid"

// HeaderRequestID is the canonical header name.
const HeaderRequestID = "X-Request-ID"

// generateID returns a new UUID v4 string.
func generateID() string {
	return uuid.NewString()
}

// extractID gets the request ID from the request header, or returns "".
func extractID(header string) string {
	return header
}
```

- [ ] **Step 4: Implement `middleware.go`**

Create `libs/go/logger/middleware.go`:

```go
package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// RequestID middleware:
// 1. Extracts X-Request-ID from request header, or generates UUID v4
// 2. Sets X-Request-ID in response header
// 3. Creates base.With(zap.String("request_id", id)) and puts it in ctx via WithContext
func RequestID(base *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := extractID(r.Header.Get(HeaderRequestID))
			if id == "" {
				id = generateID()
			}

			w.Header().Set(HeaderRequestID, id)

			l := base.With(zap.String("request_id", id))
			ctx := WithContext(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Logging middleware:
// Logs each request with method, path, status, duration, remote_addr.
// Takes logger from ctx (set by RequestID middleware).
// Wraps http.ResponseWriter internally to capture the written status code.
// ORDERING: RequestID must be registered before Logging in the middleware chain.
// If no logger is found in ctx, falls back to zap.NewNop() (silent, no panic).
func Logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(wrapped, r)

			l := FromContext(r.Context())
			l.Info("request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", wrapped.statusCode),
				zap.Duration("duration", time.Since(start)),
				zap.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
```

- [ ] **Step 5: Run RequestID tests to verify they pass**

Run: `cd libs/go/logger && go test -v -run TestRequestID ./...`
Expected: All 3 RequestID tests PASS.

- [ ] **Step 6: Add Logging middleware tests to `middleware_test.go`**

Append to `libs/go/logger/middleware_test.go`:

```go
func TestLogging_LogsRequestFields(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	base := zap.New(core)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	// Chain: RequestID -> Logging -> inner
	handler := logger.RequestID(base)(logger.Logging()(inner))

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	assert.Equal(t, "request", entry.Message)

	fields := entry.ContextMap()
	assert.Equal(t, "POST", fields["method"])
	assert.Equal(t, "/graphql", fields["path"])
	assert.Equal(t, int64(201), fields["status"])
	assert.Contains(t, fields, "duration")
	assert.Contains(t, fields, "request_id")
	assert.Contains(t, fields, "remote_addr")
}

func TestLogging_DefaultStatus200(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	base := zap.New(core)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// no explicit WriteHeader — default 200
		_, _ = w.Write([]byte("ok"))
	})

	handler := logger.RequestID(base)(logger.Logging()(inner))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, 1, logs.Len())
	fields := logs.All()[0].ContextMap()
	assert.Equal(t, int64(200), fields["status"])
}
```

- [ ] **Step 7: Run all logger package tests**

Run: `cd libs/go/logger && go test -v ./...`
Expected: All 11 tests PASS (4 logger + 2 context + 3 RequestID + 2 Logging).

- [ ] **Step 8: Commit**

```bash
git add libs/go/logger/request_id.go libs/go/logger/middleware.go libs/go/logger/middleware_test.go
git commit -m "feat(logger): add RequestID and Logging HTTP middleware"
```

---

## Task 4: Migrate apps/api to shared logger

**Files:**

- Modify: `apps/api/go.mod` (add require + replace)
- Modify: `apps/api/cmd/server/main.go:1-120` (replace Zap init + middleware)
- Delete: `apps/api/internal/middleware/logging.go`

- [ ] **Step 1: Add dependency to `apps/api/go.mod`**

Add to `apps/api/go.mod` in the `require` block:

```
monorepo-template/libs/go/logger v0.0.0
```

Add a `replace` directive at the end:

```
replace monorepo-template/libs/go/logger => ../../libs/go/logger
```

Run: `cd apps/api && go mod tidy`

- [ ] **Step 2: Update `main.go` — replace logger initialization**

In `apps/api/cmd/server/main.go`, replace the imports and Zap initialization.

Update the import block: keep `"go.uber.org/zap"` as a direct import (`main.go` still calls `zap.Error`, `zap.Int`, `zap.String` directly), add `"monorepo-template/libs/go/logger"`.

Replace lines 32-50 (Zap initialization):

```go
	level, err := zap.ParseAtomicLevel(cfg.Log.Level)
	if err != nil {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	var zapCfg zap.Config
	if cfg.Log.Format == "json" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}
	zapCfg.Level = level

	logger, err := zapCfg.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = logger.Sync() }()
```

With:

```go
	l, err := logger.New(logger.Config{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = l.Sync() }()
```

Then rename all `logger` references in `main.go` to `l` (since `logger` is now the package name). This affects:

- `postgres.New(cfg.Postgres, l)`
- `postgres.RunMigrations(cfg.Postgres.DSN(), l)`
- `redisRepo.New(cfg.Redis, l)`
- `middleware.Auth(cfg.Auth.JWTSecret, l)`
- `l.Fatal(...)`, `l.Info(...)` throughout

- [ ] **Step 3: Update `main.go` — replace middleware registration**

Replace:

```go
	r.Use(middleware.Logging(logger))
```

With:

```go
	r.Use(logger.RequestID(l))
	r.Use(logger.Logging())
```

- [ ] **Step 4: Delete `internal/middleware/logging.go`**

```bash
rm apps/api/internal/middleware/logging.go
```

- [ ] **Step 5: Verify build compiles**

Run: `cd apps/api && go build ./cmd/server/`
Expected: Compiles without errors.

- [ ] **Step 6: Run existing tests**

Run: `cd apps/api && go test ./...`
Expected: All existing tests pass (auth, health, cors, config, user_service).

- [ ] **Step 7: Commit**

```bash
git add apps/api/go.mod apps/api/go.sum apps/api/cmd/server/main.go
git rm apps/api/internal/middleware/logging.go
git commit -m "refactor(api): migrate to shared logger package"
```

---

## Task 5: Add `op` pattern to service layer

**Files:**

- Modify: `apps/api/internal/service/user_service.go:1-68`

- [ ] **Step 1: Add logger import and `op` pattern to `UserService`**

Replace the import block in `user_service.go` with:

```go
import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"go.uber.org/zap"

	"monorepo-template/libs/go/logger"
)
```

Update each method to include `const op` and `logger.FromContext(ctx)`. Example for `GetByID`:

```go
func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) {
	const op = "UserService.GetByID"
	log := logger.FromContext(ctx).With(zap.String("op", op))

	log.Debug("getting user", zap.String("user_id", id))

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}
```

Apply same pattern to: `List`, `Create`, `Update`, `Delete`.

For `Create` — preserve existing bcrypt logic, add `op` wrapping:

```go
func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	const op = "UserService.Create"
	log := logger.FromContext(ctx).With(zap.String("op", op))

	log.Debug("creating user", zap.String("email", input.Email))

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%s: hash password: %w", op, err)
	}
	input.Password = string(hashed)

	user, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}
```

- [ ] **Step 2: Verify build compiles**

Run: `cd apps/api && go build ./cmd/server/`
Expected: Compiles without errors.

- [ ] **Step 3: Run existing tests**

Run: `cd apps/api && go test ./internal/service/ -v`
Expected: Existing user_service tests pass. `FromContext` returns nop in tests (no ctx logger set), so no log output — tests unaffected.

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/service/user_service.go
git commit -m "feat(api): add op tracing to UserService methods"
```

---

## Task 6: Add `op` pattern to repository layer

**Files:**

- Modify: `apps/api/internal/repository/postgres/user_repo.go:1-142`

- [ ] **Step 1: Add logger import and `op` pattern to `UserRepo`**

Replace the import block in `user_repo.go` with:

```go
import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/service"
	"monorepo-template/libs/go/logger"
)
```

Update each method. Example for `GetByID`:

```go
func (r *UserRepo) GetByID(ctx context.Context, id string) (*service.User, error) {
	const op = "UserRepo.GetByID"
	log := logger.FromContext(ctx).With(zap.String("op", op))

	log.Debug("querying user by id", zap.String("user_id", id))

	var u service.User
	var createdAt, updatedAt time.Time
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Name, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	u.CreatedAt = createdAt.Format(time.RFC3339Nano)
	u.UpdatedAt = updatedAt.Format(time.RFC3339Nano)
	return &u, nil
}
```

Apply same pattern to: `Create`, `Update`, `Delete`.

For `List` (has multiple error sites — show full rewrite):

```go
func (r *UserRepo) List(ctx context.Context, first *int, after *string) ([]*service.User, int, error) {
	const op = "UserRepo.List"
	log := logger.FromContext(ctx).With(zap.String("op", op))

	log.Debug("listing users")

	limit := 20
	if first != nil && *first > 0 {
		limit = *first
	}

	args := []any{limit + 1}
	query := `SELECT id, email, name, created_at, updated_at FROM users`

	if after != nil && *after != "" {
		decoded, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: invalid cursor: %w", op, err)
		}
		cursor, err := time.Parse(time.RFC3339Nano, string(decoded))
		if err != nil {
			return nil, 0, fmt.Errorf("%s: invalid cursor time: %w", op, err)
		}
		query += ` WHERE created_at < $2`
		args = append(args, cursor)
	}

	query += ` ORDER BY created_at DESC LIMIT $1`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var users []*service.User
	for rows.Next() {
		var u service.User
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &createdAt, &updatedAt); err != nil {
			return nil, 0, fmt.Errorf("%s: scan: %w", op, err)
		}
		u.CreatedAt = createdAt.Format(time.RFC3339Nano)
		u.UpdatedAt = updatedAt.Format(time.RFC3339Nano)
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: iterate: %w", op, err)
	}

	if len(users) > limit {
		users = users[:limit]
	}

	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%s: count: %w", op, err)
	}

	return users, total, nil
}
```

For `Update` and `Delete` stubs — add `op` wrapping to the error return but skip `logger.FromContext` since they return immediately:

```go
func (r *UserRepo) Update(ctx context.Context, id string, input service.UpdateUserInput) (*service.User, error) {
	const op = "UserRepo.Update"
	return nil, fmt.Errorf("%s: not implemented", op)
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	const op = "UserRepo.Delete"
	return fmt.Errorf("%s: not implemented", op)
}
```

- [ ] **Step 2: Verify build compiles**

Run: `cd apps/api && go build ./cmd/server/`
Expected: Compiles without errors.

- [ ] **Step 3: Run all tests**

Run: `cd apps/api && go test ./... -v`
Expected: All tests pass.

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/repository/postgres/user_repo.go
git commit -m "feat(api): add op tracing to UserRepo methods"
```

---

## Task 7: Add `op` pattern to GraphQL resolvers

**Files:**

- Modify: `apps/api/internal/graph/schema.resolvers.go:1-123`

- [ ] **Step 1: Add logger to implemented resolvers**

Add `"go.uber.org/zap"` and `"monorepo-template/libs/go/logger"` to the import block.

**Note:** This file is code-generated by gqlgen. The resolver implementations (and their imports) are preserved on regeneration, but if `gqlgen generate` is re-run, verify the logger imports are still present.

Update `CreateUser`, `User`, `Users` resolvers (the ones that have real implementations, not `panic`).

Example for `User` resolver:

```go
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	const op = "queryResolver.User"
	log := logger.FromContext(ctx).With(zap.String("op", op))

	log.Debug("resolving user", zap.String("user_id", id))

	u, err := r.UserService.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if u == nil {
		return nil, nil
	}
	return mapUser(u), nil
}
```

Do **not** modify `UpdateUser` and `DeleteUser` — they are stubs with `panic`.

- [ ] **Step 2: Verify build compiles**

Run: `cd apps/api && go build ./cmd/server/`
Expected: Compiles without errors.

- [ ] **Step 3: Run all tests**

Run: `cd apps/api && go test ./... -v`
Expected: All tests pass.

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/graph/schema.resolvers.go
git commit -m "feat(api): add op tracing to GraphQL resolvers"
```

---

## Task 8: Final verification

- [ ] **Step 1: Run full logger package tests**

Run: `cd libs/go/logger && go test -v -count=1 ./...`
Expected: All 11 tests PASS.

- [ ] **Step 2: Run full apps/api tests**

Run: `cd apps/api && go test -v -count=1 ./...`
Expected: All existing tests PASS.

- [ ] **Step 3: Build the binary**

Run: `cd apps/api && go build -o /dev/null ./cmd/server/`
Expected: Builds without errors.

- [ ] **Step 4: Run linter**

Run: `cd apps/api && golangci-lint run ./...`
Expected: No new lint errors.

- [ ] **Step 5: Verify no remaining references to old logging middleware**

Run: `grep -r "middleware\.Logging" apps/api/` — should return nothing.
Check that `apps/api/internal/middleware/logging.go` does not exist.
