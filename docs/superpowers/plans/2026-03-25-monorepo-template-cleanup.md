# Monorepo Template Cleanup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix 16 issues across the monorepo template in a single commit — covering migration leftovers, security hardening, tooling, testing, and Docker.

**Architecture:** All changes are independent patches grouped into one commit. No new modules or architectural changes — only config, middleware, and tooling fixes.

**Tech Stack:** Next.js 15, Go (chi + gqlgen + zap + pgxpool), Nx, Vitest, Playwright, Docker, Lefthook

**Spec:** `docs/superpowers/specs/2026-03-25-monorepo-template-cleanup-design.md`

---

## File Map

| Action | File | Purpose |
|--------|------|---------|
| Modify | `apps/web/e2e/playwright.config.ts` | #4: Fix pnpm → bun |
| Delete | `apps/web/pnpm-lock.yaml` | #5: Remove stale lockfile |
| Modify | `.gitignore` | #5: Ignore pnpm lockfiles |
| Modify | `apps/web/package.json` | #6, #22: Fix eslint-config version, remove unused dep |
| Modify | `tools/codegen/package.json` | #22: Remove unused dep |
| Delete | `apps/web/src/pages/` | #19: Remove Pages Router dir |
| Modify | `apps/web/vitest.config.ts` | #19: Remove @pages alias |
| Modify | `apps/api/internal/config/config.go` | #8, #9, #13, #14, #15: Add CORS/env/pool/validation config fields |
| Modify | `apps/api/internal/middleware/cors.go` | #8: Remove DefaultCORSConfig |
| Modify | `apps/api/cmd/server/main.go` | #8, #9, #15: CORS from config, playground guard, log level |
| Modify | `apps/api/internal/repository/postgres/postgres.go` | #14: Pool config |
| Modify | `apps/api/config/config.yml` | #8, #9, #14: Add cors_origins, env, pool defaults |
| Modify | `.lefthook.yml` | #12: Fix glob pattern |
| Modify | `nx.json` | #16: Register nx-go plugin |
| Modify | `package.json` | #16, #21: Add nx-go dep, engines field |
| Modify | `.eslintrc.json` | #17: Real depConstraints |
| Modify | `apps/api/project.json` | #17: Add scope tag |
| Modify | `apps/web/project.json` | #17: Add scope tag |
| Modify | `libs/graphql/project.json` | #17: Add scope tag |
| Create | `apps/web/vitest.setup.ts` | #11: Register jest-dom matchers |
| Modify | `apps/web/vitest.config.ts` | #11, #19: Setup file + remove @pages alias |
| Create | `apps/web/app/__tests__/page.test.tsx` | #11: Smoke test |
| Modify | `docker/web.Dockerfile` | #20: Fix standalone COPY paths |

---

## Task 1: Web Cleanup (#4, #5, #6, #19, #22)

**Files:**
- Modify: `apps/web/e2e/playwright.config.ts:21`
- Delete: `apps/web/pnpm-lock.yaml`
- Modify: `.gitignore` (append)
- Modify: `apps/web/package.json:25,34`
- Modify: `tools/codegen/package.json:9`
- Delete: `apps/web/src/pages/` (directory)
- Modify: `apps/web/vitest.config.ts:26`

- [ ] **Step 1: Fix Playwright command**

In `apps/web/e2e/playwright.config.ts`, change line 21:

```typescript
// Before:
command: 'pnpm dev',
// After:
command: 'bun run dev',
```

- [ ] **Step 2: Delete stale pnpm lockfile**

```bash
rm apps/web/pnpm-lock.yaml
```

- [ ] **Step 3: Add pnpm-lock.yaml to .gitignore**

Append to `.gitignore` at the end:

```gitignore
# Stale lockfiles
**/pnpm-lock.yaml
```

- [ ] **Step 4: Fix eslint-config-next version**

In `apps/web/package.json`, change line 34:

```json
// Before:
"eslint-config-next": "14.2.0",
// After:
"eslint-config-next": "^15.0.0",
```

- [ ] **Step 5: Remove unused typescript-react-query from web**

In `apps/web/package.json`, delete line 25:

```json
"@graphql-codegen/typescript-react-query": "^6.0.0",
```

- [ ] **Step 6: Remove unused typescript-react-query from codegen**

In `tools/codegen/package.json`, delete line 9:

```json
"@graphql-codegen/typescript-react-query": "^6.0.0",
```

- [ ] **Step 7: Delete src/pages/ directory**

```bash
rm -rf apps/web/src/pages/
```

- [ ] **Step 8: Remove @pages alias from vitest config**

In `apps/web/vitest.config.ts`, delete line 26:

```typescript
// Remove this line:
'@pages': resolve(__dirname, './src/pages'),
```

- [ ] **Step 9: Verify web project**

```bash
cd apps/web && bunx tsc --noEmit
```

Expected: no errors.

---

## Task 2: API Config Expansion (#8, #9, #13, #14, #15)

**Files:**
- Modify: `apps/api/internal/config/config.go`
- Modify: `apps/api/config/config.yml`

- [ ] **Step 1: Add new fields to ServerConfig**

In `apps/api/internal/config/config.go`, replace `ServerConfig`:

```go
type ServerConfig struct {
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	Env             string        `mapstructure:"env"`
	CORSOrigins     []string      `mapstructure:"cors_origins"`
}
```

- [ ] **Step 2: Add pool fields to PostgresConfig**

In the same file, replace `PostgresConfig`:

```go
type PostgresConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DB              string        `mapstructure:"db"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxConns        int32         `mapstructure:"max_conns"`
	MinConns        int32         `mapstructure:"min_conns"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}
```

- [ ] **Step 3: Add Validate method**

First, replace the import block at the top of `config.go`:

```go
import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)
```

Then append after the `Load` function:

```go
func (c *Config) Validate() error {
	if c.Server.Port <= 0 {
		return fmt.Errorf("server.port must be > 0")
	}
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("auth.jwt_secret is required")
	}
	if c.Postgres.Host == "" {
		return fmt.Errorf("postgres.host is required")
	}
	if c.Postgres.DB == "" {
		return fmt.Errorf("postgres.db is required")
	}
	return nil
}
```

- [ ] **Step 4: Call Validate in Load**

In `Load()`, replace the tail (from `var cfg Config` to `return &cfg, nil`):

```go
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
```

- [ ] **Step 5: Add env bindings**

In the `binds` slice, add:

```go
{"server.env", "APP_ENV"},
{"server.cors_origins", "CORS_ORIGINS"},
```

- [ ] **Step 6: Update config.yml with new defaults**

Replace `apps/api/config/config.yml`:

```yaml
server:
  port: 8080
  read_timeout: 10s
  write_timeout: 30s
  shutdown_timeout: 5s
  env: development
  cors_origins:
    - "http://localhost:3000"

log:
  level: info
  format: json

postgres:
  max_conns: 10
  min_conns: 2
  max_conn_idle_time: 30m

pagination:
  default_page_size: 20
  max_page_size: 100
```

- [ ] **Step 7: Verify Go compiles**

```bash
cd apps/api && go build ./...
```

Expected: no errors.

---

## Task 3: CORS Middleware (#8)

**Files:**
- Modify: `apps/api/internal/middleware/cors.go:14-20`

- [ ] **Step 1: Remove DefaultCORSConfig**

In `apps/api/internal/middleware/cors.go`, delete the entire `DefaultCORSConfig()` function (lines 14-20):

```go
// DELETE this entire function:
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}
}
```

---

## Task 4: DB Pool Configuration (#14)

**Files:**
- Modify: `apps/api/internal/repository/postgres/postgres.go:18-27`

- [ ] **Step 1: Use pgxpool.ParseConfig with pool settings**

Replace the `New` function:

```go
func New(cfg config.PostgresConfig, logger *zap.Logger) (*DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB, cfg.SSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	if cfg.MaxConns > 0 {
		poolCfg.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns > 0 {
		poolCfg.MinConns = cfg.MinConns
	}
	if cfg.MaxConnIdleTime > 0 {
		poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	logger.Info("connected to PostgreSQL",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("db", cfg.DB),
	)

	return &DB{Pool: pool, logger: logger}, nil
}
```

---

## Task 5: Main.go Updates (#8, #9, #15)

**Files:**
- Modify: `apps/api/cmd/server/main.go`

- [ ] **Step 1: Replace logger initialization with configurable level**

Replace lines 32-37:

```go
// Before:
var logger *zap.Logger
if cfg.Log.Format == "json" {
    logger, _ = zap.NewProduction()
} else {
    logger, _ = zap.NewDevelopment()
}

// After:
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
```

- [ ] **Step 2: Replace CORS middleware call**

Replace line 62:

```go
// Before:
r.Use(middleware.CORS(middleware.DefaultCORSConfig()))

// After:
r.Use(middleware.CORS(middleware.CORSConfig{
    AllowedOrigins: cfg.Server.CORSOrigins,
    AllowedMethods: []string{"GET", "POST", "OPTIONS"},
    AllowedHeaders: []string{"Content-Type", "Authorization"},
}))
```

- [ ] **Step 3: Add playground guard**

Replace line 69:

```go
// Before:
r.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

// After:
if cfg.Server.Env != "production" {
    r.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))
}
```

- [ ] **Step 4: Verify Go compiles**

```bash
cd apps/api && go build ./cmd/server
```

Expected: no errors.

---

## Task 6: Tooling & DX (#12, #16, #17, #21)

**Files:**
- Modify: `.lefthook.yml:7,10`
- Modify: `nx.json:28`
- Modify: `package.json`
- Modify: `.eslintrc.json:13-16`
- Modify: `apps/api/project.json`
- Modify: `apps/web/project.json`
- Modify: `libs/graphql/project.json`

- [ ] **Step 1: Fix Lefthook glob**

In `.lefthook.yml`, change both `glob: '*.go'` to `glob: '**/*.go'`:

```yaml
pre-commit:
  parallel: true
  commands:
    lint-staged:
      run: bunx lint-staged
    go-lint:
      glob: '**/*.go'
      run: cd apps/api && golangci-lint run --fix
    go-test:
      glob: '**/*.go'
      run: cd apps/api && go test -short ./...

commit-msg:
  commands:
    commitlint:
      run: bunx commitlint --edit {1}
```

- [ ] **Step 2: Add engines to root package.json**

In `package.json`, add after `"private": true,`:

```json
"engines": {
  "node": ">=22.0.0",
  "bun": ">=1.0.0"
},
```

- [ ] **Step 3: Add @nx-go/nx-go to devDependencies**

In `package.json`, add to `devDependencies`:

```json
"@nx-go/nx-go": "^3.0.0",
```

- [ ] **Step 4: Register nx-go plugin in nx.json**

In `nx.json`, change line 28:

```json
// Before:
"plugins": []

// After:
"plugins": ["@nx-go/nx-go"]
```

- [ ] **Step 5: Add scope tags to project.json files**

In `apps/api/project.json`, add after `"projectType": "application",`:

```json
"tags": ["scope:api"],
```

In `apps/web/project.json`, add after `"projectType": "application",`:

```json
"tags": ["scope:web"],
```

In `libs/graphql/project.json`, add after `"projectType": "library",`:

```json
"tags": ["scope:shared"],
```

- [ ] **Step 6: Configure real depConstraints**

Replace `.eslintrc.json`:

```json
{
  "root": true,
  "ignorePatterns": ["node_modules", "dist", ".next", "*.config.*"],
  "plugins": ["@nx"],
  "overrides": [
    {
      "files": ["*.ts", "*.tsx", "*.js", "*.jsx"],
      "rules": {
        "@nx/enforce-module-boundaries": [
          "error",
          {
            "allow": [],
            "depConstraints": [
              {
                "sourceTag": "scope:web",
                "onlyDependOnLibsWithTags": ["scope:shared"]
              },
              {
                "sourceTag": "scope:api",
                "onlyDependOnLibsWithTags": ["scope:shared"]
              },
              {
                "sourceTag": "scope:shared",
                "onlyDependOnLibsWithTags": ["scope:shared"]
              }
            ]
          }
        ]
      }
    },
    {
      "files": ["*.ts", "*.tsx"],
      "extends": ["plugin:@nx/typescript"],
      "rules": {}
    },
    {
      "files": ["*.js", "*.jsx"],
      "extends": ["plugin:@nx/javascript"],
      "rules": {}
    }
  ]
}
```

- [ ] **Step 7: Install new dependency**

```bash
bun add -d @nx-go/nx-go
```

---

## Task 7: Web Smoke Test (#11)

**Files:**
- Create: `apps/web/app/__tests__/page.test.tsx`

- [ ] **Step 1: Create vitest setup file for jest-dom matchers**

Create `apps/web/vitest.setup.ts`:

```typescript
import '@testing-library/jest-dom';
```

Then update `apps/web/vitest.config.ts` — change `setupFiles: []` to:

```typescript
setupFiles: ['./vitest.setup.ts'],
```

- [ ] **Step 2: Create smoke test**

Create `apps/web/app/__tests__/page.test.tsx`:

```tsx
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import HomePage from '../page';

describe('HomePage', () => {
  it('renders heading', () => {
    render(<HomePage />);
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent(
      'Monorepo Template',
    );
  });
});
```

- [ ] **Step 3: Run test**

```bash
cd apps/web && bun run test
```

Expected: 1 test passes.

---

## Task 8: Docker Standalone Fix (#20)

**Files:**
- Modify: `docker/web.Dockerfile:31`

- [ ] **Step 1: Fix standalone COPY paths**

In `docker/web.Dockerfile`, replace line 31:

```dockerfile
# Before:
COPY --from=builder /app/apps/web/.next/standalone ./

# After:
COPY --from=builder /app/apps/web/.next/standalone/apps/web/ ./
COPY --from=builder /app/apps/web/.next/standalone/node_modules ./node_modules
```

Lines 32-33 (`.next/static` and `public`) remain unchanged.

---

## Task 9: Install Dependencies & Final Verification

- [ ] **Step 1: Install all dependencies**

```bash
bun install
```

- [ ] **Step 2: Run lint**

```bash
bun run lint
```

Expected: passes (or only pre-existing warnings).

- [ ] **Step 3: Run tests**

```bash
bun run test
```

Expected: web test passes, api tests pass.

- [ ] **Step 4: Run build**

```bash
bun run build
```

Expected: both apps build successfully.

- [ ] **Step 5: Commit all changes**

```bash
git add -A
git commit -m "fix: cleanup monorepo template (16 issues)"
```
