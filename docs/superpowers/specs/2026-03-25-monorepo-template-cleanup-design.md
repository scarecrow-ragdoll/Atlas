# Monorepo Template Cleanup — Design Spec

**Date:** 2026-03-25
**Scope:** 16 fixes in a single commit
**Skipped:** #7 (logger error handling), #10 (Redis kept as example), #18 (FSD stubs kept)

---

## 1. Web Cleanup (#4, #5, #6, #19, #22)

### #4 — Playwright command
- **File:** `apps/web/e2e/playwright.config.ts:21`
- **Change:** `pnpm dev` → `bun run dev`

### #5 — Stale pnpm lockfile
- **Delete:** `apps/web/pnpm-lock.yaml`
- **Add to root `.gitignore`:** `**/pnpm-lock.yaml`

### #6 — eslint-config-next version mismatch
- **File:** `apps/web/package.json:34`
- **Change:** `"eslint-config-next": "14.2.0"` → `"eslint-config-next": "^15.0.0"`

### #19 — Remove `src/pages/` (App Router only)
- **Delete:** `apps/web/src/pages/` directory (contains only `.gitkeep`)

### #22 — Unused `typescript-react-query` dependency
- **Remove from:** `apps/web/package.json` (`@graphql-codegen/typescript-react-query`)
- **Remove from:** `tools/codegen/package.json` (`@graphql-codegen/typescript-react-query`)
- **Rationale:** `tools/codegen/codegen.ts` only uses `typescript` and `typescript-operations` plugins

---

## 2. API Hardening (#8, #9, #13, #14, #15)

### #8 — CORS origins from config
- **File:** `apps/api/internal/config/config.go`
  - Add `CORSOrigins []string` to `ServerConfig`
  - Bind to `CORS_ORIGINS` env var (comma-separated)
- **File:** `apps/api/internal/middleware/cors.go`
  - Remove `DefaultCORSConfig()`; `CORS()` still accepts full `CORSConfig` struct
- **File:** `apps/api/cmd/server/main.go`
  - Construct `CORSConfig` inline with all three fields:
    - `AllowedOrigins`: from `cfg.Server.CORSOrigins`
    - `AllowedMethods`: hardcoded `["GET", "POST", "OPTIONS"]`
    - `AllowedHeaders`: hardcoded `["Content-Type", "Authorization"]`
- **Default in `config.yml`:** `cors_origins: ["http://localhost:3000"]`

### #9 — Playground guard
- **File:** `apps/api/internal/config/config.go`
  - Add `Env string` to `ServerConfig`, bind to `APP_ENV`
- **File:** `apps/api/cmd/server/main.go`
  - Wrap **only** the `/playground` route (line 69) in `if cfg.Server.Env != "production"` — leave `/graphql` (line 68) unconditional
- **Default:** `"development"`

### #13 — Config validation
- **File:** `apps/api/internal/config/config.go`
  - Add `Validate() error` method on `Config`
  - Required fields: `Auth.JWTSecret`, `Server.Port > 0`, `Postgres.Host`, `Postgres.DB`
  - Call `Validate()` after `Unmarshal` in `Load()`

### #14 — DB pool configuration
- **File:** `apps/api/internal/config/config.go`
  - Add to `PostgresConfig`: `MaxConns int32`, `MinConns int32`, `MaxConnIdleTime time.Duration`
- **File:** `apps/api/internal/repository/postgres/postgres.go`
  - Use `pgxpool.ParseConfig(dsn)` instead of `pgxpool.New(dsn)`
  - Apply pool settings from config
  - Create pool with `pgxpool.NewWithConfig()`
- **Defaults in config.yml:** `max_conns: 10`, `min_conns: 2`, `max_conn_idle_time: 30m`

### #15 — Apply log level from config
- **File:** `apps/api/cmd/server/main.go`
  - Replace `zap.NewProduction()` / `zap.NewDevelopment()` with manual `zap.Config` construction
  - Parse `cfg.Log.Level` via `zap.ParseAtomicLevel()`
  - Apply parsed level to the logger config

---

## 3. Tooling / DX (#12, #16, #17, #21)

### #12 — Lefthook glob
- **File:** `.lefthook.yml:7,10`
- **Change:** `glob: '*.go'` → `glob: '**/*.go'` (both `go-lint` and `go-test` commands)

### #16 — Register nx-go plugin
- **Install:** `@nx-go/nx-go` as devDependency in root `package.json`
- **File:** `nx.json`
  - Add `"@nx-go/nx-go"` to `plugins` array
- **Note:** `apps/api/project.json` uses explicit `nx:run-commands` targets — the plugin's inference will not override them

### #17 — ESLint depConstraints
- **File:** `.eslintrc.json`
  - Replace wildcard constraint with scoped rules:
    - `scope:web` → depends on `scope:shared` only
    - `scope:api` → depends on `scope:shared` only
    - `scope:shared` → depends on `scope:shared` only
- **Files:** Add `"tags"` to each project's `project.json`:
  - `apps/web/project.json` → `["scope:web"]`
  - `apps/api/project.json` → `["scope:api"]`
  - `libs/graphql/project.json` → `["scope:shared"]`

### #21 — Engines field
- **File:** `package.json`
  - Add `"engines": { "node": ">=22.0.0", "bun": ">=1.0.0" }`

---

## 4. Testing + Docker (#11, #20)

### #11 — Web smoke test
- **Create:** `apps/web/app/__tests__/page.test.tsx`
- **Content:** Render `<HomePage />` (from `apps/web/app/page.tsx`), assert it renders heading text
- **Stack:** `vitest` + `@testing-library/react` + `jsdom` (all already in devDeps)

### #20 — Docker standalone paths
- **File:** `docker/web.Dockerfile`
- **Issue:** Monorepo standalone output preserves directory structure (`apps/web/server.js`), but `CMD ["node", "server.js"]` expects it at root
- **Fix:** Replace line 31 (`COPY --from=builder /app/apps/web/.next/standalone ./`) with two lines:
  ```dockerfile
  COPY --from=builder /app/apps/web/.next/standalone/apps/web/ ./
  COPY --from=builder /app/apps/web/.next/standalone/node_modules ./node_modules
  ```
  Lines 32-33 (`.next/static` and `public` COPY) remain unchanged.
- **Verified:** `output: 'standalone'` already exists in `next.config.js:5`

---

## Delivery

- Single commit: `fix: cleanup monorepo template (16 issues)`
- Run `bun install` after dependency changes
- Verify: `bun run lint`, `bun run test`, `bun run build`
