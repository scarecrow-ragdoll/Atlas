<!--
FILE: docs/superpowers/specs/2026-06-19-wave-01-foundation-design.md
VERSION: 1.0.0
START_MODULE_CONTRACT
  PURPOSE: Define the detailed WAVE-01 (Foundation) design for the Atlas fitness tracker, covering directory structure, GraphQL schema, service/repository interfaces, PIN auth flow, middleware chains, route wiring, Docker Compose/config extensions, and verification commands.
  SCOPE: Settings service, default user bootstrap, PIN guard with Argon2id, Redis session store, separate Atlas GraphQL endpoint at /graphql/atlas, media REST scaffold, and admin auth separation; excludes domain CRUD (exercises, workouts, cardio, body, nutrition, charts, AI, backup/import).
  DEPENDS: docs/requirements.xml, docs/development-plan.xml, docs/prd-waves/waves/wave-01.md, docs/product-verified/features/pin-guard.md, docs/technical-verified/implementation-slices.md.
  LINKS: M-API / V-M-API.
  ROLE: DOC
  MAP_MODE: SUMMARY
END_MODULE_CONTRACT
START_CHANGE_SUMMARY
  LAST_CHANGE: 1.0.0 - Initial WAVE-01 design after brainstorming and user approval.
END_CHANGE_SUMMARY
-->

# WAVE-01: Foundation — Detailed Design

**Date:** 2026-06-19
**Status:** Approved
**Approach:** A (separate gqlgen executable, self-contained Atlas module, single Go binary)

---

## Table of Contents

1. [Directory Structure](#1-directory-structure)
2. [gqlgen Config and Generation](#2-gqlgen-config-and-generation)
3. [Atlas GraphQL Schema Boundaries](#3-atlas-graphql-schema-boundaries)
4. [Resolver Package Structure](#4-resolver-package-structure)
5. [Service Interfaces](#5-service-interfaces)
6. [Repository Interfaces](#6-repository-interfaces)
7. [Migration: 003_atlas_foundation.sql](#7-migration-003_atlas_foundation.sql)
8. [PIN Service and Session Flow](#8-pin-service-and-session-flow)
9. [Atlas Middleware Chains](#9-atlas-middleware-chains)
10. [Route Wiring](#10-route-wiring)
11. [Config Extensions](#11-config-extensions)
12. [Media REST Scaffold](#12-media-rest-scaffold)
13. [Test Strategy](#13-test-strategy)
14. [Verification Commands](#14-verification-commands)
15. [WAVE-01 Deliverables](#15-wave-01-deliverables)

---

## 1. Directory Structure

```
apps/api/
├── atlas-gqlgen.yml                           # separate gqlgen config for Atlas
├── cmd/server/main.go                         # updated — wire Atlas routes
├── config/config.yml                          # updated — add atlas_pin_session config
├── internal/
│   ├── appconfig/
│   │   └── config.go                          # updated — add AtlasPinSessionConfig
│   ├── atlas/                                 # new — self-contained Atlas module
│   │   ├── models/
│   │   │   ├── settings.go                    # SettingsRecord (internal, with pinHash) + Settings (public, no pinHash)
│   │   │   └── pin.go                         # PinOperationResult, PinError, PinErrorCode
│   │   ├── graph/
│   │   │   ├── generated/                     # gqlgen output (exec.go, models.go)
│   │   │   ├── schema/
│   │   │   │   ├── schema.graphql             # shared types, enums, directives
│   │   │   │   ├── settings.graphql           # settings query + mutations
│   │   │   │   └── pin.graphql               # PIN management mutations
│   │   │   └── resolver/
│   │   │       ├── resolver.go                # root resolver struct
│   │   │       └── settings.go                # SettingsResolver
│   │   ├── middleware/
│   │   │   └── pin_guard.go                   # Atlas PIN guard middleware
│   │   ├── repository/
│   │   │   ├── postgres/
│   │   │   │   └── settings_repo.go           # SettingsRepository implementation
│   │   │   └── redis/
│   │   │       ├── pin_session_store.go       # PinSessionStore implementation
│   │   │       └── pin_attempt_store.go       # PinAttemptStore implementation
│   │   └── service/
│   │       ├── pin_service.go                 # PinService implementation (Argon2id)
│   │       ├── settings_service.go            # SettingsService implementation
│   │       └── bootstrap_service.go           # BootstrapService (EnsureDefaultUser, EnsureDefaultSettings)
│   ├── repository/postgres/
│   │   ├── migrations/
│   │   │   ├── 001_admin_users.sql            # existing
│   │   │   ├── 002_users.sql                  # existing
│   │   │   └── 003_atlas_foundation.sql       # new — atlas_users + atlas_settings
│   │   ├── queries/
│   │   │   ├── admin_users.sql                # existing
│   │   │   ├── users.sql                      # existing
│   │   │   └── atlas_settings.sql             # new — sqlc queries for settings
│   │   └── generated/
│   │       ├── admin_users/                   # existing
│   │       ├── users/                         # existing
│   │       └── atlas_settings/                # new
│   └── handler/
│       ├── health.go                          # existing
│       ├── atlas_health.go                    # new — Atlas health/readiness handlers
│       └── atlas_pin_auth.go                  # new — unlock/lock/session REST handlers
```

## 2. gqlgen Config and Generation

### `apps/api/atlas-gqlgen.yml`

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

### Generation Command

```bash
cd apps/api && go run github.com/99designs/gqlgen generate --config atlas-gqlgen.yml
```

### Nx Target (apps/api/project.json)

```json
"codegen:atlas": {
  "executor": "nx:run-commands",
  "options": {
    "command": "go run github.com/99designs/gqlgen generate --config atlas-gqlgen.yml",
    "cwd": "apps/api"
  }
}
```

## 3. Atlas GraphQL Schema Boundaries (WAVE-01 Only)

### `schema.graphql`

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

### `settings.graphql`

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

### `pin.graphql`

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

### Excluded from WAVE-01 Schema

Exercise CRUD, DailyLog/workout CRUD, cardio CRUD, body tracking CRUD, nutrition CRUD, charts, AI export, backup/import domain operations.

## 4. Resolver Package Structure

### `internal/atlas/graph/resolver/resolver.go`

```go
package resolver

import (
  "monorepo-template/apps/api/internal/atlas/service"
)

type Resolver struct {
  SettingsService service.SettingsService
  PinService      service.PinService
}
```

### `internal/atlas/graph/resolver/settings.go`

```go
func (r *Resolver) Settings(ctx context.Context) (*models.SettingsResult, error)
func (r *Resolver) UpdateSettings(ctx context.Context, input models.SettingsInput) (*models.SettingsResult, error)
func (r *Resolver) EnablePin(ctx context.Context, input models.PinEnableInput) (*models.PinOperationResult, error)
func (r *Resolver) DisablePin(ctx context.Context, input models.PinDisableInput) (*models.PinOperationResult, error)
func (r *Resolver) ChangePin(ctx context.Context, input models.PinChangeInput) (*models.PinOperationResult, error)
```

## 5. Service Interfaces

### `BootstrapService` (app startup only)

```go
type BootstrapService interface {
  EnsureDefaultUser(ctx context.Context) (uuid.UUID, error)
  EnsureDefaultSettings(ctx context.Context, userID uuid.UUID) error
}
```

Called once during app startup. Not per-request.

### `SettingsService`

```go
type SettingsService interface {
  Get(ctx context.Context, userID uuid.UUID) (*models.Settings, error)
  Update(ctx context.Context, userID uuid.UUID, input models.SettingsInput) (*models.Settings, error)
}
```

- `Settings` is the public model — never exposes `pinHash`.
- `SettingsRecord` is the internal DB model — may include `pinHash` for service-layer use only.

### `PinService`

```go
type PinService interface {
  Enable(ctx context.Context, userID uuid.UUID, pin string) error
  Disable(ctx context.Context, userID uuid.UUID, currentPin string) error
  Change(ctx context.Context, userID uuid.UUID, currentPin, newPin string) error
  Verify(ctx context.Context, userID uuid.UUID, pin string) (bool, error)
  IsEnabled(ctx context.Context, userID uuid.UUID) (bool, error)
}
```

#### PIN Hashing: Argon2id

| Parameter | Default | Configurable |
|---|---|---|
| Memory | 65536 KiB (~64 MiB) | `ATLAS_PIN_ARGON2_MEMORY` |
| Iterations | 3 | `ATLAS_PIN_ARGON2_ITERATIONS` |
| Parallelism | 2 | `ATLAS_PIN_ARGON2_PARALLELISM` |
| Key length | 32 bytes | `ATLAS_PIN_ARGON2_KEY_LENGTH` |

- PIN validation: length 4–20 digits only (configurable via `ATLAS_PIN_MIN_LENGTH` / `ATLAS_PIN_MAX_LENGTH`).
- No plaintext PIN logged or stored anywhere.
- `Disable` and `Change` revoke all existing PIN sessions for the user.

## 6. Repository Interfaces

### `SettingsRepository`

```go
type SettingsRepository interface {
  FindByUserID(ctx context.Context, userID uuid.UUID) (*models.SettingsRecord, error)
  UpdateUserSettings(ctx context.Context, userID uuid.UUID, input models.SettingsInput) (*models.SettingsRecord, error)
  UpdatePinState(ctx context.Context, userID uuid.UUID, pinEnabled bool, pinHash *string) error
}
```

No generic `Upsert` — explicit methods prevent accidental PIN hash overwrites.

### `PinSessionStore`

```go
type PinSessionStore interface {
  Create(ctx context.Context, userID uuid.UUID) (string, error)
  Validate(ctx context.Context, token string) (userID uuid.UUID, valid bool, err error)
  Revoke(ctx context.Context, token string) error
  RevokeAllByUser(ctx context.Context, userID uuid.UUID) error
}
```

- `valid=false, err=nil` = missing/expired/invalid session.
- `err != nil` = infrastructure/internal failure.
- `valid=true` = session is active, returns resolved userID.

#### Redis Key Layout

```
atlas:pin_session:<token_sha256> -> {
  "userID": "<uuid>",
  "createdAt": "<iso8601>",
  "lastSeenAt": "<iso8601>",
  "expiresAt": "<iso8601>",
  "absoluteExpiresAt": "<iso8601>"
}

atlas:pin_user_sessions:<userID> -> set of token_sha256
```

#### Session TTL Policy (Sliding)

| Property | Default | Configurable |
|---|---|---|
| Idle TTL | 8 hours | `ATLAS_PIN_SESSION_IDLE_TTL` |
| Absolute TTL | 7 days (168 hours) | `ATLAS_PIN_SESSION_ABSOLUTE_TTL` |

- Session `expiresAt` refreshes on valid activity (sliding window).
- `absoluteExpiresAt` is fixed at creation.
- `RevokeAllByUser` deletes every token hash from the user session set, then deletes the set key.

### `PinAttemptStore`

```go
type PinAttemptStore interface {
  RegisterFailure(ctx context.Context, key string) error
  RegisterSuccess(ctx context.Context, key string) error
  IsLocked(ctx context.Context, key string) (bool, time.Duration, error)
}
```

Key: `atlas:pin_attempt:<identifier>` (e.g. IP address).

#### Brute-Force Policy

| Threshold | Lockout Duration |
|---|---|
| 5 failed attempts | 5 minutes |
| Repeated lockouts | 30 minutes |

- Successful unlock resets the failed attempt counter.
- Failed attempts must be audit-logged (safe, no PIN tokens).
- Error messages must remain generic.

## 7. Migration: 003_atlas_foundation.sql

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

The default user is created by `BootstrapService.EnsureDefaultUser` during app startup.

## 8. PIN Service and Session Flow

### Unlock Flow

```
POST /api/v1/auth/pin/unlock
Body: { "pin": "1234" }
  │
  ├── IsLocked(key="atlas:pin_attempt:<ip>")?
  │     └── locked → 429 Too Many Requests
  │
  ├── Atlas User Context middleware provides default userID (startup-cached bootstrap result)
  ├── SettingsRepository.FindByUserID → SettingsRecord
  ├── PinService.Verify:
  │     └── Argon2id.Compare(pin, pinHash)
  │           ├── mismatch → PinAttemptStore.RegisterFailure → 401
  │           └── match → PinAttemptStore.RegisterSuccess
  │
  ├── PinSessionStore.Create:
  │     ├── Generate 32-byte crypto/rand token → hex
  │     ├── SHA256(token) → key
  │     ├── Store:
  │     │     SET atlas:pin_session:<hash> {payload} EX <idleTTL>
  │     │     SADD atlas:pin_user_sessions:<userID> <hash>
  │     └── Return raw token
  │
  └── Set cookie: atlas_pin_session=<token>
        HttpOnly=true, SameSite=Lax, Secure=<env-derived>, Path=/
```

### Infrastructure Failure Paths

| Failure | Behavior |
|---|---|
| Redis unreachable during PIN unlock | Return 503 Service Unavailable |
| Redis unreachable during session validate (PIN enabled) | Fail closed — return 503 |
| Redis unreachable during session validate (PIN disabled) | Allow (no auth dependency) |
| Postgres unreachable during settings read | Return 503 Service Unavailable |
| Postgres unreachable during settings write | Return 500 Internal Server Error |
| Default user missing (bootstrap failed) | Return 503 Service Unavailable |

### Validate Flow (PIN Guard Middleware)

```
Request to guarded route
  │
  ├── PIN disabled → allow (no auth required)
  │
  └── PIN enabled:
        ├── Read cookie "atlas_pin_session"
        ├── SHA256(token) → key
        ├── GET atlas:pin_session:<hash>
        │     ├── not found → 401 Unauthorized
        │     └── found:
        │           ├── Check absoluteExpiresAt → expired? → delete → 401
        │           ├── Slide idle TTL: EXPIRE <hash> <idleTTL>
        │           ├── Verify sessionUserID == context defaultUserID
        │           │     ├── mismatch → revoke session → 401 (log security event)
        │           │     └── match → attach userID to context → allow
        │           └── Update lastSeenAt
```

### Lock Flow

```
POST /api/v1/auth/pin/lock
  │
  ├── If Atlas PIN cookie exists:
  │     ├── SHA256(token) → key
  │     ├── DEL atlas:pin_session:<hash>
  │     └── SREM atlas:pin_user_sessions:<userID> <hash>
  ├── Clear atlas_pin_session cookie (Set-Cookie with MaxAge=0)
  └── Return 200 OK (always success, idempotent)
```

**Disable:**
1. Verify current PIN (Argon2id)
2. `UpdatePinState(enabled=false, hash=nil)`
3. `RevokeAllByUser` — delete all sessions

**Change:**
1. Verify current PIN (Argon2id)
2. Hash new PIN (Argon2id)
3. `UpdatePinState(enabled=true, hash=newHash)`
4. `RevokeAllByUser` — user must unlock again with new PIN

## 9. Atlas Middleware Chains

Three explicit route groups — no string-based path exceptions inside middleware.

### Public/System Chain

```
RequestID → Logging

Routes:
  GET /api/v1/healthz
  GET /api/v1/readyz
```

Health/readiness routes must NOT bootstrap default user or use Atlas User Context.

### Atlas Auth-Public Chain

```
RequestID → Logging → Atlas User Context

Routes:
  POST /api/v1/auth/pin/unlock
  POST /api/v1/auth/pin/lock
  GET  /api/v1/auth/session
```

### Atlas Guarded Chain

```
RequestID → Logging → Atlas User Context → Atlas PIN Guard

Routes:
  POST   /graphql/atlas
  GET    /api/v1/media/{id}
  POST   /api/v1/media/upload
  DELETE /api/v1/media/{id}
```

#### Atlas User Context Middleware

- Reads bootstrapped default user ID (cached from startup).
- Attaches userID to request context.
- If default user is missing unexpectedly, returns 503 Service Unavailable.
- Does NOT silently create default user per request.

#### Atlas PIN Guard Middleware

- Checks if PIN is enabled via `PinService.IsEnabled`.
- If disabled: pass (no auth required).
- If enabled: validate `atlas_pin_session` cookie via `PinSessionStore.Validate`.
  - Valid + session user matches context user → attach userID + pass.
  - Valid + session user mismatches context user → reject 401 + optionally revoke (log security event).
  - Invalid/missing → reject 401.

## 10. Route Wiring

### Existing (unchanged)

```
POST /graphql  → admin GraphQL with existing admin auth middleware chain
```

### New (wired alongside existing in `apps/api/cmd/server/main.go`)

```
Public/system group:
  GET /api/v1/healthz  → atlas health handler
  GET /api/v1/readyz   → atlas readiness handler

Atlas auth-public group:
  POST /api/v1/auth/pin/unlock  → unlock handler
  POST /api/v1/auth/pin/lock    → lock handler
  GET  /api/v1/auth/session     → session check handler

Atlas guarded group:
  POST /graphql/atlas           → Atlas GraphQL handler (Atlas gqlgen executable)
  GET  /api/v1/media/{id}       → media download scaffold
  POST /api/v1/media/upload     → media upload scaffold
  DELETE /api/v1/media/{id}     → media delete scaffold
```

## 11. Config Extensions

### `config/config.yml` — new `atlas_pin` block

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

### `internal/appconfig/config.go` — new config structs

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
  MaxFailures        int           `mapstructure:"max_failures"`
  LockoutDuration    time.Duration `mapstructure:"lockout_duration"`
  EscalatedDuration  time.Duration `mapstructure:"escalated_duration"`
}
```

Added to `Config`:

```go
type Config struct {
  // ... existing fields ...
  AtlasPin        AtlasPinConfig        `mapstructure:"atlas_pin"`
  AtlasPinSession AtlasPinSessionConfig `mapstructure:"atlas_pin_session"`
  AtlasPinAttempt AtlasPinAttemptConfig `mapstructure:"atlas_pin_attempt"`
}
```

### Environment Variables

```
ATLAS_PIN_ARGON2_MEMORY
ATLAS_PIN_ARGON2_ITERATIONS
ATLAS_PIN_ARGON2_PARALLELISM
ATLAS_PIN_ARGON2_KEY_LENGTH
ATLAS_PIN_MIN_LENGTH
ATLAS_PIN_MAX_LENGTH
ATLAS_PIN_SESSION_COOKIE_NAME
ATLAS_PIN_SESSION_IDLE_TTL
ATLAS_PIN_SESSION_ABSOLUTE_TTL
ATLAS_PIN_SESSION_COOKIE_SECURE
ATLAS_PIN_SESSION_SAME_SITE
ATLAS_PIN_ATTEMPT_MAX_FAILURES
ATLAS_PIN_ATTEMPT_LOCKOUT_DURATION
ATLAS_PIN_ATTEMPT_ESCALATED_DURATION
```

### Cookie Policy

- Name: `atlas_pin_session` (configurable via `ATLAS_PIN_SESSION_COOKIE_NAME`)
- `HttpOnly=true`
- `SameSite=Lax`
- `Secure=true` when:
  - Config `cookie_secure` is `true`, OR
  - Config `cookie_secure` is `auto` and environment is `production`
- `Path=/`

## 12. Media REST Scaffold (WAVE-01)

Foundation only — no actual file storage or serving logic in WAVE-01.

```
GET    /api/v1/media/{id}      → returns 501 Not Implemented
POST   /api/v1/media/upload    → returns 501 Not Implemented
DELETE /api/v1/media/{id}      → returns 501 Not Implemented
```

All three routes pass through the Atlas Guarded chain (PIN guard enforced when PIN is enabled). The auth boundary is established before domain logic is added in WAVE-02.

## 13. Test Strategy

### Auth Separation Tests

| Test | Expected Result |
|---|---|
| `POST /graphql` without admin auth | 401 |
| `POST /graphql` with Atlas PIN cookie only | 401 |
| `POST /graphql/atlas` PIN disabled, no auth | 200 |
| `POST /graphql/atlas` PIN enabled, valid Atlas PIN cookie | 200 |
| `POST /graphql/atlas` PIN enabled, no Atlas PIN cookie | 401 |
| `POST /graphql/atlas` PIN enabled, admin auth cookie only | 401 |
| `GET /api/v1/auth/session` PIN enabled, no session | 200 with `session_valid: false` |
| `POST /api/v1/auth/pin/unlock` without Atlas session | Allowed (unguarded) |
| `POST /api/v1/auth/pin/lock` without Atlas session | Idempotent — always clears cookie, returns success |

### Schema Boundary Tests

| Test | Expected Result |
|---|---|
| Settings GraphQL query returns no `pinHash` | Field absent from response |
| Atlas GraphQL schema excludes admin types/operations | gqlgen compile-time check |
| Admin GraphQL schema excludes Atlas types/operations | gqlgen compile-time check |

### PIN Service Tests

| Test | Expected Result |
|---|---|
| Enable PIN with valid PIN | PIN enabled, hash stored |
| Enable PIN when already enabled | Error: PIN_ALREADY_ENABLED |
| Disable PIN with correct PIN | PIN disabled, hash removed, sessions revoked |
| Disable PIN with wrong PIN | Error: WRONG_PIN |
| Change PIN with correct current PIN | New hash stored, sessions revoked |
| Disable PIN when already disabled | Error: PIN_ALREADY_DISABLED |
| Verify correct PIN | true |
| Verify incorrect PIN | false |
| PIN too short (<4 digits) | Error: PIN_TOO_SHORT |
| PIN too long (>20 digits) | Error: PIN_TOO_LONG |

### Session Store Tests

| Test | Expected Result |
|---|---|
| Create and validate session | valid=true, correct userID |
| Validate expired session | valid=false |
| Validate revoked session | valid=false |
| RevokeAllByUser removes all sessions | No active sessions remain |
| Sliding TTL renews on activity | expiresAt advances |
| Absolute TTL does not change | absoluteExpiresAt fixed |

### Brute-Force Tests

| Test | Expected Result |
|---|---|
| 5 failures → lockout | IsLocked=true |
| Successful unlock resets counter | IsLocked=false |
| Generic error messages | No hints about failure count or remaining attempts |

### Settings Repository Tests

| Test | Expected Result |
|---|---|
| FindByUserID returns settings | Correct record |
| UpdateUserSettings does not change pinHash | pinHash unchanged |
| UpdatePinState only changes pin fields | Only pin_enabled and pin_hash modified |

## 14. Verification Commands

### WAVE-01 Implementation Verification

```bash
# Atlas gqlgen codegen
cd apps/api && go run github.com/99designs/gqlgen generate --config atlas-gqlgen.yml

# sqlc compile
cd apps/api && sqlc compile

# Full API build
cd apps/api && go build ./cmd/server

# Atlas-specific tests (fast local check)
cd apps/api && go test ./internal/atlas/...

# Auth separation tests
cd apps/api && go test ./... -run "TestAtlasAuthSeparation|TestAtlasGraphQLSeparation|TestAtlasPinGuard"

# Full API test suite
cd apps/api && go test ./...

# Full lint
cd apps/api && golangci-lint run ./...

# Nx targets
bunx nx run api:codegen          # existing — admin gqlgen
bunx nx run api:codegen:atlas    # new — Atlas gqlgen
bunx nx run api:test             # existing — full api test suite
bunx nx run api:lint             # existing — golangci-lint full
```

## 15. WAVE-01 Deliverables

| ID | Deliverable | Type |
|---|---|---|
| CAP-W01-001 | Docker Compose config extensions (env vars only, no new services) | ops |
| CAP-W01-002 | Nx workspace codegen target for Atlas (`codegen:atlas`) | config |
| CAP-W01-003 | Go API binary with Atlas module under `apps/api/internal/atlas/` | code |
| CAP-W01-004 | Migration `003_atlas_foundation.sql` (`atlas_users` + `atlas_settings`) | data |
| CAP-W01-005 | Redis PIN session store with sliding TTL + user session index | code |
| CAP-W01-006 | Atlas GraphQL endpoint at `POST /graphql/atlas` | routing |
| CAP-W01-007 | PIN guard middleware + PIN attempt store | auth |
| CAP-W01-008 | Settings service + bootstrap service | domain |
| CAP-W01-009 | CI-ready verification targets | config |
| CAP-W01-010 | Media REST scaffold (GET, POST upload, DELETE) | scaffold |