# WAVE-01: Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Establish Atlas foundation: settings service, default user bootstrap, optional PIN guard with Argon2id session auth, separate GraphQL endpoint at `/graphql/atlas`, media REST scaffold, and full admin auth separation.

**Architecture:** Atlas code lives under `apps/api/internal/atlas/` as a self-contained module. Shared gqlgen with separate config (`atlas-gqlgen.yml`). Three explicit route groups in `main.go`: public/system, atlas auth-public, atlas guarded. PIN guard uses Argon2id, Redis session store with sliding TTL, brute-force protection. Admin auth remains untouched.

**Tech Stack:** Go 1.25, gqlgen, sqlc, Argon2id (golang.org/x/crypto/argon2), Redis (go-redis/redigo), PostgreSQL (pgx/v5), chi router, bcrypt-style session token hashing (SHA256)

**Design doc:** `docs/superpowers/specs/2026-06-19-wave-01-foundation-design.md`

---

## File Inventory

### New Files

| # | Path | Purpose |
|---|------|---------|
| 1 | `apps/api/atlas-gqlgen.yml` | Separate gqlgen config for Atlas schema |
| 2 | `apps/api/internal/atlas/models/settings.go` | SettingsRecord (internal, with pinHash) + Settings (public, no pinHash) |
| 3 | `apps/api/internal/atlas/models/pin.go` | PinOperationResult, PinError, PinErrorCode |
| 4 | `apps/api/internal/atlas/graph/schema/schema.graphql` | Base schema with Query/Mutation roots |
| 5 | `apps/api/internal/atlas/graph/schema/settings.graphql` | Settings types, inputs, enums |
| 6 | `apps/api/internal/atlas/graph/schema/pin.graphql` | PIN management types |
| 7 | `apps/api/internal/atlas/graph/resolver/resolver.go` | Root resolver struct (injects services) |
| 8 | `apps/api/internal/atlas/graph/resolver/settings.go` | Settings query + settings/PIN mutations |
| 9 | `apps/api/internal/atlas/service/pin_service.go` | PIN business logic (Argon2id hash, verify, enable/disable/change) |
| 10 | `apps/api/internal/atlas/service/settings_service.go` | Settings CRUD (no pinHash leak) |
| 11 | `apps/api/internal/atlas/service/bootstrap_service.go` | EnsureDefaultUser + EnsureDefaultSettings (app startup only) |
| 12 | `apps/api/internal/atlas/repository/postgres/settings_repo.go` | SettingsRepository impl (FindByUserID, UpdateUserSettings, UpdatePinState) |
| 13 | `apps/api/internal/atlas/repository/redis/pin_session_store.go` | PinSessionStore impl (Create, Validate, Revoke, RevokeAllByUser) |
| 14 | `apps/api/internal/atlas/repository/redis/pin_attempt_store.go` | PinAttemptStore impl (RegisterFailure, RegisterSuccess, IsLocked) |
| 15 | `apps/api/internal/atlas/middleware/pin_guard.go` | PIN guard middleware (checks session cookie) |
| 16 | `apps/api/internal/atlas/middleware/user_context.go` | Atlas user context middleware (attaches cached default userID) |
| 17 | `apps/api/internal/handler/atlas_health.go` | GET /api/v1/healthz, GET /api/v1/readyz |
| 18 | `apps/api/internal/handler/atlas_pin_auth.go` | POST unlock, POST lock, GET session |
| 19 | `apps/api/internal/handler/atlas_media.go` | GET/POST/DELETE media scaffold (501) |
| 20 | `apps/api/internal/repository/postgres/migrations/003_atlas_foundation.sql` | atlas_users + atlas_settings tables |
| 21 | `apps/api/internal/repository/postgres/queries/atlas_settings.sql` | sqlc queries for settings repo |
| 22 | `apps/api/internal/atlas/service/pin_service_test.go` | PIN service unit tests |
| 23 | `apps/api/internal/atlas/service/settings_service_test.go` | Settings service unit tests |
| 24 | `apps/api/internal/atlas/repository/redis/pin_session_store_test.go` | Session store tests |
| 25 | `apps/api/internal/atlas/repository/redis/pin_attempt_store_test.go` | Attempt store tests |
| 26 | `apps/api/internal/atlas/repository/postgres/settings_repo_test.go` | Settings repo tests |
| 27 | `apps/api/internal/atlas/middleware/pin_guard_test.go` | PIN guard middleware tests |
| 28 | `apps/api/internal/handler/atlas_pin_auth_test.go` | Unlock/lock/session handler tests |
| 29 | `apps/api/internal/atlas/service/bootstrap_service_test.go` | Bootstrap service tests |
| 30 | `apps/api/internal/atlas/atlas_test.go` | Auth separation + GraphQL boundary integration tests |

### Modified Files

| # | Path | Change |
|---|------|--------|
| 1 | `apps/api/cmd/server/main.go` | Wire Atlas route groups + bootstrap call |
| 2 | `apps/api/internal/appconfig/config.go` | Add AtlasPinConfig, AtlasPinSessionConfig, AtlasPinAttemptConfig to Config struct + defaults |
| 3 | `apps/api/config/config.yml` | Add atlas_pin and atlas_pin_session and atlas_pin_attempt defaults |
| 4 | `apps/api/project.json` | Add `codegen:atlas` target |
| 5 | `apps/api/go.mod` | Add golang.org/x/crypto dependency (Argon2id) |
| 6 | `.env.example` | Add ATLAS_PIN_* env vars |

---

### Task 1: Config and Dependencies

**Files:**
- Modify: `apps/api/internal/appconfig/config.go`
- Modify: `apps/api/config/config.yml`
- Modify: `apps/api/go.mod`
- Modify: `.env.example`

- [ ] **Step 1: Add Atlas config structs to appconfig**

Add to `apps/api/internal/appconfig/config.go`:

```go
type AtlasPinConfig struct {
  Argon2Memory      uint32 `mapstructure:"argon2_memory"`
  Argon2Iterations  uint32 `mapstructure:"argon2_iterations"`
  Argon2Parallelism uint8  `mapstructure:"argon2_parallelism"`
  Argon2KeyLength   uint32 `mapstructure:"argon2_key_length"`
  MinLength         int    `mapstructure:"min_length"`
  MaxLength         int    `mapstructure:"max_length"`
}

type AtlasPinSessionConfig struct {
  CookieName   string        `mapstructure:"cookie_name"`
  IdleTTL      time.Duration `mapstructure:"idle_ttl"`
  AbsoluteTTL  time.Duration `mapstructure:"absolute_ttl"`
  CookieSecure string        `mapstructure:"cookie_secure"`
  SameSite     string        `mapstructure:"same_site"`
}

type AtlasPinAttemptConfig struct {
  MaxFailures       int           `mapstructure:"max_failures"`
  LockoutDuration   time.Duration `mapstructure:"lockout_duration"`
  EscalatedDuration time.Duration `mapstructure:"escalated_duration"`
}
```

Add to `Config` struct:

```go
AtlasPin        AtlasPinConfig        `mapstructure:"atlas_pin"`
AtlasPinSession AtlasPinSessionConfig `mapstructure:"atlas_pin_session"`
AtlasPinAttempt AtlasPinAttemptConfig `mapstructure:"atlas_pin_attempt"`
```

Add `ApplyAtlasDefaults` function:

```go
func ApplyAtlasDefaults(cfg *Config) error {
  if cfg.AtlasPin.Argon2Memory == 0 { cfg.AtlasPin.Argon2Memory = 65536 }
  if cfg.AtlasPin.Argon2Iterations == 0 { cfg.AtlasPin.Argon2Iterations = 3 }
  if cfg.AtlasPin.Argon2Parallelism == 0 { cfg.AtlasPin.Argon2Parallelism = 2 }
  if cfg.AtlasPin.Argon2KeyLength == 0 { cfg.AtlasPin.Argon2KeyLength = 32 }
  if cfg.AtlasPin.MinLength == 0 { cfg.AtlasPin.MinLength = 4 }
  if cfg.AtlasPin.MaxLength == 0 { cfg.AtlasPin.MaxLength = 20 }
  if cfg.AtlasPinSession.CookieName == "" { cfg.AtlasPinSession.CookieName = "atlas_pin_session" }
  if cfg.AtlasPinSession.IdleTTL == 0 { cfg.AtlasPinSession.IdleTTL = 8 * time.Hour }
  if cfg.AtlasPinSession.AbsoluteTTL == 0 { cfg.AtlasPinSession.AbsoluteTTL = 168 * time.Hour }
  if cfg.AtlasPinSession.CookieSecure == "" { cfg.AtlasPinSession.CookieSecure = "auto" }
  if cfg.AtlasPinSession.SameSite == "" { cfg.AtlasPinSession.SameSite = "Lax" }
  if cfg.AtlasPinAttempt.MaxFailures == 0 { cfg.AtlasPinAttempt.MaxFailures = 5 }
  if cfg.AtlasPinAttempt.LockoutDuration == 0 { cfg.AtlasPinAttempt.LockoutDuration = 5 * time.Minute }
  if cfg.AtlasPinAttempt.EscalatedDuration == 0 { cfg.AtlasPinAttempt.EscalatedDuration = 30 * time.Minute }
  return nil
}
```

- [ ] **Step 2: Add default config to config.yml**

Add to `apps/api/config/config.yml`:

```yaml
atlas_pin:
  argon2_memory: 65536
  argon2_iterations: 3
  argon2_parallelism: 2
  argon2_key_length: 32
  min_length: 4
  max_length: 20

atlas_pin_session:
  cookie_name: atlas_pin_session
  idle_ttl: 8h
  absolute_ttl: 168h
  cookie_secure: auto
  same_site: Lax

atlas_pin_attempt:
  max_failures: 5
  lockout_duration: 5m
  escalated_duration: 30m
```

- [ ] **Step 3: Add Argon2 dependency**

Run: `cd apps/api && go get golang.org/x/crypto`

- [ ] **Step 4: Add env vars to .env.example**

```env
ATLAS_PIN_ARGON2_MEMORY=65536
ATLAS_PIN_ARGON2_ITERATIONS=3
ATLAS_PIN_ARGON2_PARALLELISM=2
ATLAS_PIN_ARGON2_KEY_LENGTH=32
ATLAS_PIN_MIN_LENGTH=4
ATLAS_PIN_MAX_LENGTH=20
ATLAS_PIN_SESSION_COOKIE_NAME=atlas_pin_session
ATLAS_PIN_SESSION_IDLE_TTL=8h
ATLAS_PIN_SESSION_ABSOLUTE_TTL=168h
ATLAS_PIN_SESSION_COOKIE_SECURE=auto
ATLAS_PIN_SESSION_SAME_SITE=Lax
ATLAS_PIN_ATTEMPT_MAX_FAILURES=5
ATLAS_PIN_ATTEMPT_LOCKOUT_DURATION=5m
ATLAS_PIN_ATTEMPT_ESCALATED_DURATION=30m
```

- [ ] **Step 5: Run existing tests to confirm config changes don't break anything**

Run: `cd apps/api && go test ./internal/appconfig/...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add apps/api/internal/appconfig apps/api/config/config.yml apps/api/go.mod apps/api/go.sum .env.example
git commit -m "feat(atlas): add Atlas PIN config, env vars, and Argon2 dependency"
```

---

### Task 2: Database Migration and sqlc Queries

**Files:**
- Create: `apps/api/internal/repository/postgres/migrations/003_atlas_foundation.sql`
- Create: `apps/api/internal/repository/postgres/queries/atlas_settings.sql`
- Modify: (generated output after sqlc gen)

- [ ] **Step 1: Create migration file**

`apps/api/internal/repository/postgres/migrations/003_atlas_foundation.sql`:

```sql
CREATE TABLE atlas_users (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  display_name TEXT NOT NULL DEFAULT 'Default User',
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE atlas_settings (
  id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id                 UUID NOT NULL REFERENCES atlas_users(id),
  pin_enabled             BOOLEAN NOT NULL DEFAULT false,
  pin_hash                TEXT,
  units                   TEXT NOT NULL DEFAULT 'metric',
  default_ai_export_weeks INT NOT NULL DEFAULT 4,
  created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(user_id)
);
```

- [ ] **Step 2: Create sqlc queries file**

`apps/api/internal/repository/postgres/queries/atlas_settings.sql`:

```sql
-- name: GetSettingsByUserID :one
SELECT id, user_id, pin_enabled, pin_hash, units, default_ai_export_weeks, created_at, updated_at
FROM atlas_settings
WHERE user_id = $1;

-- name: UpsertSettings :one
INSERT INTO atlas_settings (user_id, pin_enabled, pin_hash, units, default_ai_export_weeks)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id)
DO UPDATE SET
  pin_enabled = COALESCE($2, atlas_settings.pin_enabled),
  pin_hash    = COALESCE($3, atlas_settings.pin_hash),
  units       = COALESCE($4, atlas_settings.units),
  default_ai_export_weeks = COALESCE($5, atlas_settings.default_ai_export_weeks),
  updated_at  = now()
RETURNING id, user_id, pin_enabled, pin_hash, units, default_ai_export_weeks, created_at, updated_at;

-- name: UpdatePinState :one
UPDATE atlas_settings
SET pin_enabled = $2, pin_hash = $3, updated_at = now()
WHERE user_id = $1
RETURNING id, user_id, pin_enabled, pin_hash, units, default_ai_export_weeks, created_at, updated_at;

-- name: CreateDefaultUser :one
INSERT INTO atlas_users (display_name)
VALUES ('Default User')
ON CONFLICT DO NOTHING
RETURNING id, display_name, created_at, updated_at;

-- name: GetDefaultUser :one
SELECT id, display_name, created_at, updated_at
FROM atlas_users
ORDER BY created_at ASC
LIMIT 1;
```

Note: `CreateDefaultUser` uses `ON CONFLICT DO NOTHING` but there's no unique constraint on display_name. The bootstrap service will call this and check if a row was returned. If no row, it uses `GetDefaultUser`.

- [ ] **Step 3: Generate sqlc code**

Run: `cd apps/api && sqlc compile && sqlc generate`
Expected: Generated files appear in `apps/api/internal/repository/postgres/generated/atlas_settings/`

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/repository/postgres/migrations/003_atlas_foundation.sql apps/api/internal/repository/postgres/queries/atlas_settings.sql apps/api/internal/repository/postgres/generated/atlas_settings
git commit -m "feat(atlas): add atlas_users and atlas_settings migration and sqlc queries"
```

---

### Task 3: Atlas Models

**Files:**
- Create: `apps/api/internal/atlas/models/settings.go`
- Create: `apps/api/internal/atlas/models/pin.go`

- [ ] **Step 1: Create settings models**

`apps/api/internal/atlas/models/settings.go`:

```go
package models

import (
  "time"
  "github.com/google/uuid"
)

// SettingsRecord is the internal DB-facing model — may include pinHash.
// Never returned through GraphQL/API.
type SettingsRecord struct {
  ID                  uuid.UUID
  UserID              uuid.UUID
  PinEnabled          bool
  PinHash             *string
  Units               string
  DefaultAiExportWeeks int
  CreatedAt           time.Time
  UpdatedAt           time.Time
}

// Settings is the public API-facing model — no pinHash.
type Settings struct {
  PinEnabled          bool   `json:"pin_enabled"`
  Units               string `json:"units"`
  DefaultAiExportWeeks int   `json:"default_ai_export_weeks"`
}

// SettingsInput is used for partial updates.
type SettingsInput struct {
  Units               *string `json:"units,omitempty"`
  DefaultAiExportWeeks *int   `json:"default_ai_export_weeks,omitempty"`
}

// SettingsResult is the GraphQL response wrapper.
type SettingsResult struct {
  Settings *Settings     `json:"settings,omitempty"`
  Error    *SettingsError `json:"error,omitempty"`
}

type SettingsError struct {
  Message string            `json:"message"`
  Code    SettingsErrorCode `json:"code"`
}

type SettingsErrorCode string

const (
  SettingsErrorValidation  SettingsErrorCode = "VALIDATION_ERROR"
  SettingsErrorUnauthorized SettingsErrorCode = "UNAUTHORIZED"
  SettingsErrorInternal    SettingsErrorCode = "INTERNAL_ERROR"
)
```

- [ ] **Step 2: Create PIN models**

`apps/api/internal/atlas/models/pin.go`:

```go
package models

type PinEnableInput struct {
  Pin string `json:"pin"`
}

type PinDisableInput struct {
  CurrentPin string `json:"current_pin"`
}

type PinChangeInput struct {
  CurrentPin string `json:"current_pin"`
  NewPin     string `json:"new_pin"`
}

type PinOperationResult struct {
  Success bool      `json:"success"`
  Error   *PinError `json:"error,omitempty"`
}

type PinError struct {
  Message string      `json:"message"`
  Code    PinErrorCode `json:"code"`
}

type PinErrorCode string

const (
  PinErrorWrongPIN          PinErrorCode = "WRONG_PIN"
  PinErrorAlreadyEnabled    PinErrorCode = "PIN_ALREADY_ENABLED"
  PinErrorAlreadyDisabled   PinErrorCode = "PIN_ALREADY_DISABLED"
  PinErrorTooShort          PinErrorCode = "PIN_TOO_SHORT"
  PinErrorTooLong           PinErrorCode = "PIN_TOO_LONG"
  PinErrorSessionExpired    PinErrorCode = "SESSION_EXPIRED"
  PinErrorInternal          PinErrorCode = "INTERNAL_ERROR"
)
```

- [ ] **Step 3: Verify build**

Run: `cd apps/api && go build ./internal/atlas/models`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/atlas/models
git commit -m "feat(atlas): add Settings and PIN models (public vs internal separation)"
```

---

### Task 4: Atlas GraphQL Schema and gqlgen Config

**Files:**
- Create: `apps/api/atlas-gqlgen.yml`
- Create: `apps/api/internal/atlas/graph/schema/schema.graphql`
- Create: `apps/api/internal/atlas/graph/schema/settings.graphql`
- Create: `apps/api/internal/atlas/graph/schema/pin.graphql`
- Modify: `apps/api/project.json`

- [ ] **Step 1: Create gqlgen config**

`apps/api/atlas-gqlgen.yml`:

```yaml
schema:
  - internal/atlas/graph/schema/*.graphql
exec:
  filename: internal/atlas/graph/generated/exec.go
  package: generated
model:
  filename: internal/atlas/graph/generated/models.go
resolver:
  layout: follow-schema
  dir: internal/atlas/graph/resolver
  package: resolver
autobind:
  - monorepo-template/apps/api/internal/atlas/models
models:
  Settings:
    model: monorepo-template/apps/api/internal/atlas/models.Settings
  PinOperationResult:
    model: monorepo-template/apps/api/internal/atlas/models.PinOperationResult
  SettingsResult:
    model: monorepo-template/apps/api/internal/atlas/models.SettingsResult
  SettingsError:
    model: monorepo-template/apps/api/internal/atlas/models.SettingsError
  PinError:
    model: monorepo-template/apps/api/internal/atlas/models.PinError
```

- [ ] **Step 2: Create GraphQL schema files**

`schema.graphql`:
```graphql
scalar Time

type Query {
  settings: SettingsResult!
}

type Mutation {
  updateSettings(input: SettingsInput!): SettingsResult!
  enablePin(input: PinEnableInput!): PinOperationResult!
  disablePin(input: PinDisableInput!): PinOperationResult!
  changePin(input: PinChangeInput!): PinOperationResult!
}
```

`settings.graphql`:
```graphql
type Settings {
  pinEnabled: Boolean!
  units: String!
  defaultAiExportWeeks: Int!
}

type SettingsResult {
  settings: Settings
  error: SettingsError
}

type SettingsError {
  message: String!
  code: SettingsErrorCode!
}

enum SettingsErrorCode {
  VALIDATION_ERROR
  UNAUTHORIZED
  INTERNAL_ERROR
}

input SettingsInput {
  units: String
  defaultAiExportWeeks: Int
}
```

`pin.graphql`:
```graphql
input PinEnableInput {
  pin: String!
}

input PinDisableInput {
  currentPin: String!
}

input PinChangeInput {
  currentPin: String!
  newPin: String!
}

type PinOperationResult {
  success: Boolean!
  error: PinError
}

type PinError {
  message: String!
  code: PinErrorCode!
}

enum PinErrorCode {
  WRONG_PIN
  PIN_ALREADY_ENABLED
  PIN_ALREADY_DISABLED
  PIN_TOO_SHORT
  PIN_TOO_LONG
  SESSION_EXPIRED
  INTERNAL_ERROR
}
```

- [ ] **Step 3: Add Nx codegen target**

In `apps/api/project.json`, add:

```json
"codegen:atlas": {
  "executor": "nx:run-commands",
  "options": {
    "command": "go run github.com/99designs/gqlgen generate --config atlas-gqlgen.yml",
    "cwd": "apps/api"
  }
}
```

- [ ] **Step 4: Generate Atlas gqlgen code**

Run: `cd apps/api && go run github.com/99designs/gqlgen generate --config atlas-gqlgen.yml`
Expected: Generated files appear in `apps/api/internal/atlas/graph/generated/` and resolver stubs in `apps/api/internal/atlas/graph/resolver/`

- [ ] **Step 5: Commit**

```bash
git add apps/api/atlas-gqlgen.yml apps/api/internal/atlas/graph/schema apps/api/internal/atlas/graph/generated apps/api/project.json
git commit -m "feat(atlas): add Atlas GraphQL schema and gqlgen config"
```

---

### Task 4.5: Service Error Types (prerequisite for resolvers)

**Files:**
- Create: `apps/api/internal/atlas/service/errors.go`

- [ ] **Step 1: Define service error types**

`apps/api/internal/atlas/service/errors.go`:

```go
package service

import "errors"

var (
  ErrPinWrongPIN          = errors.New("wrong PIN")
  ErrPinAlreadyEnabled    = errors.New("PIN already enabled")
  ErrPinAlreadyDisabled   = errors.New("PIN already disabled")
  ErrPinTooShort          = errors.New("PIN too short")
  ErrPinTooLong           = errors.New("PIN too long")
  ErrSettingsInvalidUnits = errors.New("invalid units: must be metric or imperial")
  ErrSettingsInvalidWeeks = errors.New("default AI export weeks must be >= 1")
)

type SettingsError struct {
  Err     error
  Code    string
  Message string
}

func (e *SettingsError) Error() string { return e.Err.Error() }
func (e *SettingsError) Unwrap() error { return e.Err }
```

- [ ] **Step 2: Commit**

```bash
git add apps/api/internal/atlas/service/errors.go
git commit -m "feat(atlas): add typed service error definitions"
```

### Task 5: Resolver Implementation

**Files:**
- Create: `apps/api/internal/atlas/graph/resolver/resolver.go`
- Modify: `apps/api/internal/atlas/graph/resolver/settings.go` (gqlgen stub → real implementation)

- [ ] **Step 1: Write resolver struct**

`apps/api/internal/atlas/graph/resolver/resolver.go`:

```go
package resolver

import (
  "monorepo-template/apps/api/internal/atlas/service"
)

type Resolver struct {
  SettingsService service.SettingsService
  PinService      service.PinService
}

func New(ss service.SettingsService, ps service.PinService) *Resolver {
  return &Resolver{
    SettingsService: ss,
    PinService:      ps,
  }
}
```

- [ ] **Step 2: Implement resolver with real error mapping**

`apps/api/internal/atlas/graph/resolver/settings.go` (add resolver methods and helpers to gqlgen-generated stub):

```go
package resolver

import (
  "context"
  "errors"

  "monorepo-template/apps/api/internal/atlas/models"
  "monorepo-template/apps/api/internal/atlas/middleware"
  "monorepo-template/apps/api/internal/atlas/service"
)

func serviceErrToSettingsErr(err error) *models.SettingsError {
  var se *service.SettingsError
  if errors.As(err, &se) {
    return &models.SettingsError{Message: se.Message, Code: models.SettingsErrorCode(se.Code)}
  }
  return &models.SettingsError{Message: "internal error", Code: models.SettingsErrorInternal}
}

func serviceErrToPinErr(err error) *models.PinError {
  switch {
  case errors.Is(err, service.ErrPinWrongPIN):
    return &models.PinError{Message: err.Error(), Code: models.PinErrorWrongPIN}
  case errors.Is(err, service.ErrPinAlreadyEnabled):
    return &models.PinError{Message: err.Error(), Code: models.PinErrorAlreadyEnabled}
  case errors.Is(err, service.ErrPinAlreadyDisabled):
    return &models.PinError{Message: err.Error(), Code: models.PinErrorAlreadyDisabled}
  case errors.Is(err, service.ErrPinTooShort):
    return &models.PinError{Message: err.Error(), Code: models.PinErrorTooShort}
  case errors.Is(err, service.ErrPinTooLong):
    return &models.PinError{Message: err.Error(), Code: models.PinErrorTooLong}
  default:
    return &models.PinError{Message: "internal error", Code: models.PinErrorInternal}
  }
}

- [ ] **Step 3: Build to verify resolver compiles**

Run: `cd apps/api && go build ./internal/atlas/graph/resolver`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/atlas/graph/resolver
git commit -m "feat(atlas): implement GraphQL resolvers for settings and PIN mutations"
```

---

### Task 6: Repository Implementations

**Files:**
- Create: `apps/api/internal/atlas/repository/postgres/settings_repo.go`
- Create: `apps/api/internal/atlas/repository/redis/pin_session_store.go`
- Create: `apps/api/internal/atlas/repository/redis/pin_attempt_store.go`

- [ ] **Step 1: Write settings repository**

`apps/api/internal/atlas/repository/postgres/settings_repo.go`:

```go
package postgres

import (
  "context"
  "github.com/google/uuid"
  "github.com/jackc/pgx/v5/pgxpool"
  "monorepo-template/apps/api/internal/atlas/models"
)

type SettingsRepository struct {
  pool *pgxpool.Pool
}

func NewSettingsRepository(pool *pgxpool.Pool) *SettingsRepository {
  return &SettingsRepository{pool: pool}
}

func (r *SettingsRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.SettingsRecord, error) {
  // Uses generated sqlc query GetSettingsByUserID
}

func (r *SettingsRepository) UpdateUserSettings(ctx context.Context, userID uuid.UUID, input models.SettingsInput) (*models.SettingsRecord, error) {
  // Uses generated sqlc query UpsertSettings with only non-nil fields
  // Does NOT modify pin_enabled or pin_hash
}

func (r *SettingsRepository) UpdatePinState(ctx context.Context, userID uuid.UUID, pinEnabled bool, pinHash *string) error {
  // Uses generated sqlc query UpdatePinState
}
```

- [ ] **Step 2: Write PIN session store**

`apps/api/internal/atlas/repository/redis/pin_session_store.go`:

```go
package redis

import (
  "context"
  "crypto/rand"
  "crypto/sha256"
  "encoding/hex"
  "encoding/json"
  "fmt"
  "time"
  "github.com/redis/go-redis/v9"
  "github.com/google/uuid"
)

type SessionPayload struct {
  UserID            string    `json:"userID"`
  CreatedAt         time.Time `json:"createdAt"`
  LastSeenAt        time.Time `json:"lastSeenAt"`
  ExpiresAt         time.Time `json:"expiresAt"`
  AbsoluteExpiresAt time.Time `json:"absoluteExpiresAt"`
}

type PinSessionStore struct {
  client      *redis.Client
  idleTTL     time.Duration
  absoluteTTL time.Duration
}

func NewPinSessionStore(client *redis.Client, idleTTL, absoluteTTL time.Duration) *PinSessionStore

func (s *PinSessionStore) Create(ctx context.Context, userID uuid.UUID) (string, error)
  // 1. Generate 32 random bytes → hex token
  // 2. SHA256(token) → key
  // 3. Build payload with createdAt, lastSeenAt, expiresAt (now+idleTTL), absoluteExpiresAt (now+absoluteTTL)
  // 4. SET atlas:pin_session:<key> json EX <idleTTL seconds>
  // 5. SADD atlas:pin_user_sessions:<userID> <key>
  // 6. Return raw token

func (s *PinSessionStore) Validate(ctx context.Context, token string) (uuid.UUID, bool, error)
  // 1. SHA256(token) → key
  // 2. GET atlas:pin_session:<key>
  // 3. If not found → return zero, false, nil
  // 4. Parse JSON payload
  // 5. If time.Now().After(payload.AbsoluteExpiresAt) → delete key, return zero, false, nil
  // 6. Slide: EXPIRE <key> <idleTTL seconds>, update LastSeenAt in payload, SET the updated JSON back to the same key
  // 7. Return userID, true, nil

func (s *PinSessionStore) Revoke(ctx context.Context, token string) error
  // 1. SHA256(token) → key
  // 2. DEL atlas:pin_session:<key>
  // 3. Optionally SREM from user set (requires parsing payload first)

func (s *PinSessionStore) RevokeAllByUser(ctx context.Context, userID uuid.UUID) error
  // 1. SMEMBERS atlas:pin_user_sessions:<userID>
  // 2. For each hash: DEL atlas:pin_session:<hash>
  // 3. DEL atlas:pin_user_sessions:<userID>
```

- [ ] **Step 3: Write PIN attempt store**

`apps/api/internal/atlas/repository/redis/pin_attempt_store.go`:

```go
package redis

type PinAttemptStore struct {
  client            *redis.Client
  maxFailures       int
  lockoutDuration   time.Duration
  escalatedDuration time.Duration
}

func NewPinAttemptStore(client *redis.Client, maxFailures int, lockoutDuration, escalatedDuration time.Duration) *PinAttemptStore

const attemptKeyPrefix = "atlas:pin_attempt:"

func (s *PinAttemptStore) RegisterFailure(ctx context.Context, key string) error
  // 1. key = atlas:pin_attempt:<key>
  // 2. INCR key
  // 3. If INCR result == 1: EXPIRE key <escalatedDuration> (first failure)
  // 4. If INCR result >= maxFailures: check if previous lockout happened → use escalatedDuration
  // 5. Log audit-safe message

func (s *PinAttemptStore) RegisterSuccess(ctx context.Context, key string) error
  // 1. DELETE atlas:pin_attempt:<key>

func (s *PinAttemptStore) IsLocked(ctx context.Context, key string) (bool, time.Duration, error)
  // 1. GET atlas:pin_attempt:<key>
  // 2. If count >= maxFailures → return true, remaining TTL
  // 3. Return false, 0, nil
```

- [ ] **Step 4: Write repository tests**

Create test files for each repository. Tests use the existing test infrastructure (`apps/api/internal/testinfra`).

`apps/api/internal/atlas/repository/postgres/settings_repo_test.go`:
- TestFindByUserID_ReturnsSettings
- TestUpdateUserSettings_DoesNotModifyPinHash
- TestUpdatePinState_OnlyChangesPinFields

`apps/api/internal/atlas/repository/redis/pin_session_store_test.go`:
- TestCreateAndValidate
- TestValidateExpiredSession (set short TTL)
- TestValidateRevokedSession
- TestRevokeAllByUser
- TestSlidingTTL

`apps/api/internal/atlas/repository/redis/pin_attempt_store_test.go`:
- TestMaxFailures_LocksOut
- TestRegisterSuccess_ResetsCounter
- TestIsLocked_ReturnsRemainingDuration

- [ ] **Step 5: Run repository tests**

Run: `cd apps/api && go test ./internal/atlas/repository/...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add apps/api/internal/atlas/repository
git commit -m "feat(atlas): implement Postgres settings repo, Redis session store, and attempt store"
```

---

### Task 7: Service Implementations

**Files:**
- Create: `apps/api/internal/atlas/service/pin_service.go`
- Create: `apps/api/internal/atlas/service/settings_service.go`
- Create: `apps/api/internal/atlas/service/bootstrap_service.go`

- [ ] **Step 1: Write bootstrap service**

`apps/api/internal/atlas/service/bootstrap_service.go`:

```go
package service

type BootstrapService struct {
  queries *generated.Queries // sqlc-generated
}

func NewBootstrapService(queries *generated.Queries) *BootstrapService

func (s *BootstrapService) EnsureDefaultUser(ctx context.Context) (uuid.UUID, error)
  // 1. CreateDefaultUser query (ON CONFLICT DO NOTHING)
  // 2. If no row returned → GetDefaultUser
  // 3. Return userID

func (s *BootstrapService) EnsureDefaultSettings(ctx context.Context, userID uuid.UUID) error
  // 1. UpsertSettings with defaults (pinEnabled=false, units=metric, defaultAiExportWeeks=4)
  // 2. If already exists, no-op
```

- [ ] **Step 2: Write settings service**

`apps/api/internal/atlas/service/settings_service.go`:

```go
package service

type SettingsService struct {
  repo SettingsRepository
}

func NewSettingsService(repo SettingsRepository) *SettingsService

func (s *SettingsService) Get(ctx context.Context, userID uuid.UUID) (*models.Settings, error)
  // 1. repo.FindByUserID → SettingsRecord
  // 2. Map to Settings (no pinHash)
  // 3. Return

func (s *SettingsService) Update(ctx context.Context, userID uuid.UUID, input models.SettingsInput) (*models.Settings, error)
  // 1. Validate input (units must be "metric" or "imperial", weeks >= 1)
  // 2. repo.UpdateUserSettings (does NOT touch pin_enabled or pin_hash)
  // 3. Map to Settings
  // 4. Return
```

- [ ] **Step 3: Write PIN service**

`apps/api/internal/atlas/service/pin_service.go`:

```go
package service

type PinService struct {
  settingsRepo  SettingsRepository
  sessionStore  PinSessionStore
  attemptStore  PinAttemptStore
  config        AtlasPinConfig
}

func NewPinService(settingsRepo SettingsRepository, sessionStore PinSessionStore, attemptStore PinAttemptStore, config AtlasPinConfig) *PinService

func (s *PinService) Enable(ctx context.Context, userID uuid.UUID, pin string) error
  // 1. Validate PIN length (s.config.MinLength/MaxLength)
  // 2. Get current settings: if pin_enabled → PinErrorAlreadyEnabled
  // 3. Hash PIN with Argon2id (s.config params)
  // 4. repo.UpdatePinState(enabled=true, hash=argon2Hash)
  // 5. Return nil

func (s *PinService) Disable(ctx context.Context, userID uuid.UUID, currentPin string) error
  // 1. Get settings
  // 2. If not pin_enabled → PinErrorAlreadyDisabled
  // 3. Verify currentPin against stored hash (Argon2id compare)
  // 4. If mismatch → PinErrorWrongPIN
  // 5. repo.UpdatePinState(enabled=false, hash=nil)
  // 6. sessionStore.RevokeAllByUser(userID)
  // 7. Return nil

func (s *PinService) Change(ctx context.Context, userID uuid.UUID, currentPin, newPin string) error
  // 1. Validate new PIN length
  // 2. Get settings, verify currentPin
  // 3. Hash new PIN
  // 4. repo.UpdatePinState(enabled=true, hash=newHash)
  // 5. sessionStore.RevokeAllByUser(userID)
  // 6. Return nil

func (s *PinService) Verify(ctx context.Context, userID uuid.UUID, pin string) (bool, error)
  // 1. Get settings
  // 2. If pinHash is nil → return false, nil
  // 3. Argon2id compare
  // 4. Return match result

func (s *PinService) IsEnabled(ctx context.Context, userID uuid.UUID) (bool, error)
  // 1. Get settings
  // 2. Return pinEnabled
```

Argon2id helper (in pin_service.go or a separate file):

```go
import "golang.org/x/crypto/argon2"

type Argon2Params struct {
  Memory      uint32
  Iterations  uint32
  Parallelism uint8
  KeyLength   uint32
}

func hashPIN(pin string, params Argon2Params) (string, error) {
  salt := make([]byte, 16)
  _, _ = rand.Read(salt)
  hash := argon2.IDKey([]byte(pin), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)
  // Encode as: $argon2id$v=19$m=memory,t=iterations,p=parallelism$salt_base64$hash_base64
  return encodeArgon2Hash(salt, hash, params), nil
}

func verifyPIN(pin, encodedHash string, params Argon2Params) (bool, error) {
  // Decode stored hash, extract salt + params
  // Hash candidate PIN with same salt/params
  // Constant-time compare
}
```

- [ ] **Step 4: Write service tests**

`apps/api/internal/atlas/service/pin_service_test.go`:
- TestEnable_ValidPIN
- TestEnable_AlreadyEnabled
- TestDisable_CorrectPIN_RevokesSessions
- TestDisable_WrongPIN
- TestChange_CorrectPIN_RevokesSessions
- TestVerify_CorrectPIN
- TestVerify_WrongPIN
- TestPINTooShort
- TestPINTooLong

`apps/api/internal/atlas/service/settings_service_test.go`:
- TestGet_ReturnsSettings
- TestUpdate_ValidInput
- TestUpdate_InvalidUnits

`apps/api/internal/atlas/service/bootstrap_service_test.go`:
- TestEnsureDefaultUser_CreatesOnFirstCall
- TestEnsureDefaultUser_ReturnsExisting
- TestEnsureDefaultSettings_CreatesDefaults

- [ ] **Step 5: Run service tests**

Run: `cd apps/api && go test ./internal/atlas/service/...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add apps/api/internal/atlas/service
git commit -m "feat(atlas): implement bootstrap, settings, and PIN services with Argon2id"
```

---

### Task 8: Atlas Middleware

**Files:**
- Create: `apps/api/internal/atlas/middleware/user_context.go`
- Create: `apps/api/internal/atlas/middleware/pin_guard.go`

- [ ] **Step 1: Write user context middleware**

`apps/api/internal/atlas/middleware/user_context.go`:

```go
package atlasctx

import (
  "context"
  "net/http"
  "github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "atlas_user_id"

// NewUserContextMiddleware returns middleware that attaches the bootstrapped default userID.
// The userIDProvider is called once during startup and cached.
func NewUserContextMiddleware(getDefaultUserID func() uuid.UUID) func(http.Handler) http.Handler {
  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      userID := getDefaultUserID()
      if userID == uuid.Nil {
        http.Error(w, "atlas default user not initialized", http.StatusServiceUnavailable)
        return
      }
      ctx := context.WithValue(r.Context(), userIDKey, userID)
      next.ServeHTTP(w, r.WithContext(ctx))
    })
  }
}

// GetUserID extracts the default userID from context.
func GetUserID(ctx context.Context) uuid.UUID {
  id, _ := ctx.Value(userIDKey).(uuid.UUID)
  return id
}
```

- [ ] **Step 2: Write PIN guard middleware**

`apps/api/internal/atlas/middleware/pin_guard.go`:

```go
package atlasctx

type PinGuardMiddleware struct {
  pinService    PinService
  sessionStore  PinSessionStore
  cookieName    string
}

func NewPinGuardMiddleware(pinService PinService, sessionStore PinSessionStore, cookieName string) *PinGuardMiddleware

func (m *PinGuardMiddleware) Middleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    userID := GetUserID(r.Context())
    if userID == uuid.Nil {
      http.Error(w, "unauthorized", http.StatusUnauthorized)
      return
    }

    // Check if PIN is enabled
    enabled, err := m.pinService.IsEnabled(r.Context(), userID)
    if err != nil {
      http.Error(w, "service unavailable", http.StatusServiceUnavailable)
      return
    }
    if !enabled {
      next.ServeHTTP(w, r)
      return
    }

    // PIN is enabled — validate session
    cookie, err := r.Cookie(m.cookieName)
    if err != nil {
      http.Error(w, "unauthorized", http.StatusUnauthorized)
      return
    }

    sessionUserID, valid, err := m.sessionStore.Validate(r.Context(), cookie.Value)
    if err != nil {
      http.Error(w, "service unavailable", http.StatusServiceUnavailable)
      return
    }
    if !valid {
      http.Error(w, "unauthorized", http.StatusUnauthorized)
      return
    }
    if sessionUserID != userID {
      // Session user != context user — security event
      _ = m.sessionStore.Revoke(r.Context(), cookie.Value)
      http.Error(w, "unauthorized", http.StatusUnauthorized)
      return
    }

    next.ServeHTTP(w, r)
  })
}
```

- [ ] **Step 3: Write middleware tests**

`apps/api/internal/atlas/middleware/pin_guard_test.go`:
- TestPinDisabled_AllowsRequest
- TestPinEnabled_ValidCookie_Allows
- TestPinEnabled_NoCookie_Rejects
- TestPinEnabled_ExpiredCookie_Rejects
- TestPinEnabled_SessionUserMismatch_RejectsAndRevokes

- [ ] **Step 4: Run middleware tests**

Run: `cd apps/api && go test ./internal/atlas/middleware/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/atlas/middleware
git commit -m "feat(atlas): implement user context and PIN guard middleware"
```

---

### Task 9: REST Handlers

**Files:**
- Create: `apps/api/internal/handler/atlas_health.go`
- Create: `apps/api/internal/handler/atlas_pin_auth.go`
- Create: `apps/api/internal/handler/atlas_media.go`

- [ ] **Step 1: Write health handlers**

`apps/api/internal/handler/atlas_health.go`:

```go
package handler

func AtlasHealthz(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(`{"status":"ok"}`))
}

func AtlasReadyz(w http.ResponseWriter, r *http.Request) {
  // Optionally check DB/Redis connectivity
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(`{"status":"ready"}`))
}
```

- [ ] **Step 2: Write PIN auth handlers**

`apps/api/internal/handler/atlas_pin_auth.go`:

```go
package handler

type PinAuthHandler struct {
  pinService    service.PinService
  sessionStore  redis.PinSessionStore
  attemptStore  redis.PinAttemptStore
  userContext   func() uuid.UUID  // returns cached default userID
  cookieConfig  CookieConfig
}

type CookieConfig struct {
  Name     string
  Secure   bool
  SameSite http.SameSite
}

func NewPinAuthHandler(pinService, sessionStore, attemptStore, userContext, cookieConfig) *PinAuthHandler

// POST /api/v1/auth/pin/unlock
func (h *PinAuthHandler) Unlock(w http.ResponseWriter, r *http.Request) {
  // 1. Parse body { "pin": "..." }
  // 2. Check attempt store IsLocked(ip) → 429
  // 3. Get default userID from cached provider
  // 4. PinService.Verify(userID, pin)
  //   - false → RegisterFailure, return 401
  //   - true → RegisterSuccess
  // 5. sessionStore.Create(userID) → token
  // 6. Set cookie: atlas_pin_session=<token> (HttpOnly, SameSite=Lax, Secure, Path=/)
  // 7. Return 200 { "session_valid": true }
}

// POST /api/v1/auth/pin/lock
func (h *PinAuthHandler) Lock(w http.ResponseWriter, r *http.Request) {
  // 1. Read atlas_pin_session cookie
  // 2. If exists → sessionStore.Revoke(token)
  // 3. Clear cookie (Set-Cookie with MaxAge=0)
  // 4. Return 200 { "success": true }
}

// GET /api/v1/auth/session
func (h *PinAuthHandler) Session(w http.ResponseWriter, r *http.Request) {
  // 1. Read atlas_pin_session cookie
  // 2. If no cookie → return 200 { "session_valid": false }
  // 3. sessionStore.Validate(token)
  // 4. Return 200 { "session_valid": true/false }
}
```

- [ ] **Step 3: Write media scaffold handlers**

`apps/api/internal/handler/atlas_media.go`:

```go
package handler

func MediaDownload(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusNotImplemented)
  w.Write([]byte(`{"error":"not_implemented","message":"Media download not implemented in WAVE-01"}`))
}

func MediaUpload(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusNotImplemented)
  w.Write([]byte(`{"error":"not_implemented","message":"Media upload not implemented in WAVE-01"}`))
}

func MediaDelete(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusNotImplemented)
  w.Write([]byte(`{"error":"not_implemented","message":"Media delete not implemented in WAVE-01"}`))
}
```

- [ ] **Step 4: Write handler tests**

`apps/api/internal/handler/atlas_pin_auth_test.go`:
- TestUnlock_ValidPIN_ReturnsSessionCookie
- TestUnlock_InvalidPIN_Returns401
- TestUnlock_LockedOut_Returns429
- TestLock_RevokesSessionAndClearsCookie
- TestLock_NoCookie_Idempotent
- TestSession_ValidCookie_ReturnsTrue
- TestSession_NoCookie_ReturnsFalse
- TestSession_ExpiredCookie_ReturnsFalse

- [ ] **Step 5: Run handler tests**

Run: `cd apps/api && go test ./internal/handler/...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add apps/api/internal/handler
git commit -m "feat(atlas): add health, PIN auth, and media scaffold REST handlers"
```

---

### Task 10: Route Wiring in main.go

**Files:**
- Modify: `apps/api/cmd/server/main.go`

- [ ] **Step 1: Wire Atlas routes**

In `apps/api/cmd/server/main.go`, after existing admin routes:

```go
// --- Atlas routes ---

// Bootstrap default user and settings at startup
bootstrapSvc := service.NewBootstrapService(queries)
defaultUserID, err := bootstrapSvc.EnsureDefaultUser(ctx)
if err != nil { /* fatal */ }
err = bootstrapSvc.EnsureDefaultSettings(ctx, defaultUserID)
if err != nil { /* fatal */ }

// Repos
settingsRepo := postgres.NewSettingsRepository(pool)
pinSessionStore := redis.NewPinSessionStore(rdb, cfg.AtlasPinSession.IdleTTL, cfg.AtlasPinSession.AbsoluteTTL)
pinAttemptStore := redis.NewPinAttemptStore(rdb, cfg.AtlasPinAttempt.MaxFailures, cfg.AtlasPinAttempt.LockoutDuration, cfg.AtlasPinAttempt.EscalatedDuration)

// Services
settingsSvc := service.NewSettingsService(settingsRepo)
pinSvc := service.NewPinService(settingsRepo, pinSessionStore, pinAttemptStore, cfg.AtlasPin)

// GraphQL handler — construct gqlgen executable from generated schema
import gqlgenHandler "github.com/99designs/gqlgen/graphql/handler"
import gqlgenGenerated "monorepo-template/apps/api/internal/atlas/graph/generated"

gqlCfg := gqlgenGenerated.Config{Resolvers: resolver.New(settingsSvc, pinSvc)}
atlasSchema := gqlgenGenerated.NewExecutableSchema(gqlCfg)
atlasGraphQL := gqlgenHandler.NewDefaultServer(atlasSchema)

// Middleware
userCtxMW := atlasctx.NewUserContextMiddleware(func() uuid.UUID { return defaultUserID })
pinGuardMW := atlasctx.NewPinGuardMiddleware(pinSvc, pinSessionStore, cfg.AtlasPinSession.CookieName)

// Cookie config for PIN auth handlers
cookieCfg := handler.CookieConfig{
  Name:     cfg.AtlasPinSession.CookieName,
  Secure:   cfg.AtlasPinSession.CookieSecure == "true" || (cfg.AtlasPinSession.CookieSecure == "auto" && cfg.Server.Env == "production"),
  SameSite: parseSameSite(cfg.AtlasPinSession.SameSite),
}

// Route groups
r.Group(func(r chi.Router) {
  // Public/system — no user context, no PIN guard
  r.Get("/api/v1/healthz", handler.AtlasHealthz)
  r.Get("/api/v1/readyz", handler.AtlasReadyz)
})

r.Group(func(r chi.Router) {
  // Auth-public — user context but no PIN guard
  r.Use(userCtxMW)
  
  pinAuthHandler := handler.NewPinAuthHandler(pinSvc, pinSessionStore, pinAttemptStore, func() uuid.UUID { return defaultUserID }, cookieCfg)
  r.Post("/api/v1/auth/pin/unlock", pinAuthHandler.Unlock)
  r.Post("/api/v1/auth/pin/lock", pinAuthHandler.Lock)
  r.Get("/api/v1/auth/session", pinAuthHandler.Session)
})

r.Group(func(r chi.Router) {
  // Atlas guarded — user context + PIN guard
  r.Use(userCtxMW)
  r.Use(pinGuardMW.Middleware)
  
  r.Post("/graphql/atlas", atlasGraphQL.ServeHTTP)
  r.Get("/api/v1/media/{id}", handler.MediaDownload)
  r.Post("/api/v1/media/upload", handler.MediaUpload)
  r.Delete("/api/v1/media/{id}", handler.MediaDelete)
})
```

Add helper in main.go or a utility file:

```go
func parseSameSite(value string) http.SameSite {
  switch strings.ToLower(strings.TrimSpace(value)) {
  case "strict":
    return http.SameSiteStrictMode
  case "none":
    return http.SameSiteNoneMode
  default:
    return http.SameSiteLaxMode
  }
}
```

- [ ] **Step 2: Build the full binary**

Run: `cd apps/api && go build ./cmd/server`
Expected: Binary compiles without errors

- [ ] **Step 3: Commit**

```bash
git add apps/api/cmd/server/main.go
git commit -m "feat(atlas): wire Atlas routes, middleware chains, and bootstrap in main.go"
```

---

### Task 11: Auth Separation Integration Tests

**Files:**
- Create: `apps/api/internal/atlas/atlas_test.go`

- [ ] **Step 1: Write auth separation tests**

`apps/api/internal/atlas/atlas_test.go`:

```go
package atlas_test

// These tests use the existing testinfra package for API server + DB + Redis setup.

// TestAdminGraphQL_WithoutAdminAuth_Returns401
// TestAdminGraphQL_WithAtlasPinCookie_Returns401
// TestAtlasGraphQL_PINDisabled_NoAuth_Returns200
// TestAtlasGraphQL_PINEnabled_ValidCookie_Returns200
// TestAtlasGraphQL_PINEnabled_NoCookie_Returns401
// TestAtlasGraphQL_PINEnabled_AdminCookieOnly_Returns401
// TestAtlasGraphQL_SettingsQuery_NoPinHashInResponse

// Compile-time schema boundary assertions:
// TestAtlasGraphQL_SchemaExcludesAdminTypes — introspect Atlas gqlgen schema,
//   assert that query fields like "me" and mutation fields like "loginAdmin"
//   are NOT present. Use generated schemaconfig.Types() or introspect the
//   executable schema.
// TestAdminGraphQL_SchemaExcludesAtlasTypes — introspect admin gqlgen schema,
//   assert that query fields like "settings" and mutation fields like
//   "enablePin" are NOT present.
```

- [ ] **Step 2: Run integration tests**

Run: `cd apps/api && go test ./internal/atlas/... -run "TestAdminGraphQL|TestAtlasGraphQL|TestUnlock|TestLock|TestSession" -v`
Expected: All auth separation tests PASS

- [ ] **Step 3: Run full test suite**

Run: `cd apps/api && go test ./...`
Expected: All existing tests + new tests PASS

- [ ] **Step 4: Run full lint**

Run: `cd apps/api && golangci-lint run ./...`
Expected: No lint errors

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/atlas/atlas_test.go
git commit -m "test(atlas): add auth separation and GraphQL boundary integration tests"
```

---

### Task 12: Final Verification

- [ ] **Step 1: Codegen verification**

```bash
cd apps/api && go run github.com/99designs/gqlgen generate --config atlas-gqlgen.yml
cd apps/api && sqlc compile
```
Expected: Both succeed

- [ ] **Step 2: Full build**

```bash
cd apps/api && go build ./cmd/server
```
Expected: Binary compiles

- [ ] **Step 3: Full test suite**

```bash
cd apps/api && go test ./...
```
Expected: ALL PASS

- [ ] **Step 4: Full lint**

```bash
cd apps/api && golangci-lint run ./...
```
Expected: No errors

- [ ] **Step 5: Nx target check**

```bash
bunx nx run api:codegen:atlas
bunx nx run api:build
```
Expected: Both succeed