# Web-admin Backend Auth Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build backend-only web-admin authentication with an env-bootstrapped first admin, Redis-backed httpOnly sessions, protected admin GraphQL, and cookie-ready web-admin GraphQL transport.

**Architecture:** Add a dedicated admin auth slice beside the existing reference `users` domain. PostgreSQL stores `admin_users`, Redis stores HMAC-keyed opaque sessions, the API guards admin GraphQL through cookie session context plus origin checks, and the Vite web-admin client is prepared to send credentialed GraphQL requests without adding login UI.

**Tech Stack:** Go 1.25, chi, gqlgen, sqlc, goose, pgx/v5, Redis, zap, bcrypt, Bun, Nx, Vite, React, graphql-request, Playwright, Vitest, GRACE XML.

---

<!-- FILE: docs/superpowers/plans/2026-06-07-web-admin-backend-auth.md -->
<!-- VERSION: 1.0.1 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Provide the task-by-task implementation plan for backend-only web-admin authentication. -->
<!--   SCOPE: Covers config, migrations, sqlc, admin repository/service/session store, cookies, GraphQL auth, CORS/origin checks, web-admin transport readiness, e2e setup, coverage policy, GRACE docs, verification, and commit boundaries; excludes implementation performed by this document. -->
<!--   DEPENDS: docs/superpowers/specs/2026-06-07-web-admin-backend-auth-design.md, apps/api, libs/graphql/schema, apps/web-admin, tools/coverage, docs/*.xml. -->
<!--   LINKS: M-API / V-M-API / M-GRAPHQL-SCHEMA / V-M-GRAPHQL-SCHEMA / M-WEB-ADMIN / V-M-WEB-ADMIN / M-COVERAGE-GATE / V-M-COVERAGE-GATE / M-GRACE-WORKFLOW / V-M-GRACE-WORKFLOW. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Source Spec - Anchors the approved and subagent-reviewed design. -->
<!--   Scope Check - Confirms this is one coupled backend-auth plan. -->
<!--   File Structure - Lists all planned creates and modifications with ownership boundaries. -->
<!--   Execution Discipline - Defines TDD, generated-code, GRACE, and commit rules for workers. -->
<!--   Tasks - Provides ordered implementation, test, docs, and verification steps. -->
<!--   Self-Review - Records spec coverage, placeholder scan, and type consistency checks. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.2 - Renumbered admin_users migration references to 00079 for stale local dev databases. -->
<!-- END_CHANGE_SUMMARY -->

## Source Spec

- Design: `docs/superpowers/specs/2026-06-07-web-admin-backend-auth-design.md`
- Design commits:
  - `491b63a docs: add web-admin backend auth design`
  - `6fe41eb docs: harden web-admin auth design`
- Required review-loop verdicts:
  - Intent reviewer: `APPROVE`
  - Feasibility reviewer: `APPROVE`
  - Verification reviewer: `APPROVE`
  - Security/privacy reviewer: `APPROVE`

Approved decisions:

- Use a separate `admin_users` table, not the existing reference `users` table.
- Seed the first admin from env only when `admin_users` is empty.
- Do not support registration.
- Allow later admin creation only through authenticated backend GraphQL.
- Use an httpOnly cookie containing an opaque session id.
- Store Redis session keys with a hash or HMAC-derived suffix, not the raw session id.
- Keep the public REST `/api/users` flow public.
- Protect admin GraphQL except `loginAdmin`, `logoutAdmin`, and `me`.
- Add non-UI web-admin transport readiness with credentialed GraphQL requests.
- Add CSRF/origin defense and an admin-only credentialed CORS boundary.

## Scope Check

This is one implementation plan. It touches API config, database schema, generated SQL and GraphQL code, Redis sessions, GraphQL guards, web-admin transport setup, e2e helpers, coverage policy, and GRACE docs. These surfaces are tightly coupled because the feature is not complete until cookie sessions, admin GraphQL, generated client types, browser credentials, coverage replacement gates, and GRACE verification agree.

The plan does not include:

- login page UI;
- frontend route guards;
- admin list/edit/deactivate UI;
- password reset;
- OAuth or SSO;
- tenant or owner scoping for future product-specific resources.

## File Structure

### Create

- `libs/graphql/schema/admin_auth.graphql` - Admin auth schema types and operations.
- `apps/api/internal/repository/postgres/migrations/00079_admin_users.sql` - Goose migration for `admin_users`.
- `apps/api/internal/repository/postgres/queries/admin_users.sql` - sqlc admin user queries.
- `apps/api/internal/repository/postgres/admin_repo.go` - PostgreSQL adapter implementing `service.AdminRepository`.
- `apps/api/internal/repository/postgres/admin_repo_test.go` - Real database coverage for admin repository behavior.
- `apps/api/internal/repository/redis/admin_session_store.go` - Redis session store implementing `service.AdminSessionStore`.
- `apps/api/internal/repository/redis/admin_session_store_test.go` - Redis session create/read/delete/expiry/key-hashing tests.
- `apps/api/internal/service/admin_auth.go` - Admin domain types, validation, seed, login, current admin, create admin, and logout orchestration.
- `apps/api/internal/service/admin_auth_test.go` - Service tests with fake repository/session store.
- `apps/api/internal/middleware/admin_auth.go` - Admin principal context, cookie bridge, auth middleware, and protected-operation helpers.
- `apps/api/internal/middleware/admin_auth_test.go` - Cookie/context/protected-guard tests.
- `apps/api/internal/middleware/admin_origin.go` - Admin GraphQL origin/CSRF guard.
- `apps/api/internal/middleware/admin_origin_test.go` - Allowed and rejected origin tests.
- `apps/api/internal/graph/admin_auth_resolvers_test.go` - Resolver-level admin auth and guard coverage.
- `apps/web-admin/src/entities/admin-auth/api/adminAuth.graphql` - Generated client operations for login/logout/me/createAdmin.
- `.tasks/web-admin-backend-auth/verification.md` - Execution evidence log.

### Create By Generation

- `apps/api/internal/repository/postgres/generated/admin_users.sql.go` - sqlc admin query methods.
- `apps/api/internal/graph/model/models_gen.go` - gqlgen models updated for admin auth.
- `apps/api/internal/graph/generated.go` - gqlgen execution code updated for admin auth.
- `apps/web-admin/src/shared/api/generated/types.ts` - GraphQL Codegen output updated with admin auth operations.

### Modify

- `apps/api/internal/appconfig/config.go` - Add admin seed and session config, remove unused JWT placeholder from the auth path.
- `apps/api/internal/appconfig/config_test.go` - Config validation tests for admin/session settings.
- `apps/api/config/config.yml` - Add local admin/session defaults and web-admin admin origins.
- `.env.example` - Document admin seed/session env names.
- `apps/api/cmd/server/main.go` - Wire admin repository, session store, auth service, seed bootstrap, route groups, admin CORS/origin guard, auth middleware, GraphQL cookie bridge, and resolver dependencies.
- `apps/api/internal/middleware/cors.go` - Add credentialed CORS controls and reject wildcard origins when credentials are enabled.
- `apps/api/internal/middleware/cors_test.go` - Cover credentialed and non-credentialed CORS behavior.
- `apps/api/internal/graph/resolver.go` - Add `AdminAuthService` dependency.
- `apps/api/internal/graph/schema.resolvers.go` - Add generated resolver methods after gqlgen and wire guard calls around existing user operations.
- `apps/api/internal/graph/schema_resolvers_test.go` - Add protected user CRUD denial and authenticated allowance tests.
- `apps/api/internal/repository/postgres/queries/users.sql` - No behavior change unless sqlc query file ordering requires generated output stability.
- `apps/api/sqlc.yaml` - Already reads query directory; only update `START_MODULE_MAP` if needed to mention admin query files.
- `apps/web-admin/src/shared/api/graphql-client.ts` - Enable credentialed GraphQL requests and remove bearer-token auth as the admin path.
- `apps/web-admin/src/shared/api/graphql-client.test.ts` - Prove credentialed client construction.
- `apps/web-admin/e2e/helpers.ts` - Add login helper and credentialed GraphQL context setup.
- `apps/web-admin/e2e/users-flow.spec.ts` - Log in before protected admin user CRUD flows.
- `apps/web-admin/e2e/playwright.config.ts` - Provide admin seed/session env and admin origin values to the API server.
- `docker/docker-compose.yml` - Pass admin bootstrap/session env into the local API container and remove stale JWT auth env from the auth path.
- `deploy/dokploy/docker-compose.template.yml` - Add production admin/session secret placeholders for API deployment.
- `docs/infrastructure/ci-cd.md` - Document required Dokploy/GitLab admin auth runtime secrets.
- `tools/coverage/coverage.config.json` - Add generated `admin_users.sql.go` allowlist entry with replacement gate.
- `docs/requirements.xml` - Replace the admin GraphQL auth placeholder with a real backend-auth requirement while preserving product owner/tenant warnings.
- `docs/technology.xml` - Record admin auth/session config, Redis session use, and focused commands.
- `docs/development-plan.xml` - Add/update API auth, GraphQL schema, and web-admin generated-client contracts.
- `docs/knowledge-graph.xml` - Add auth repository, service, session store, middleware, GraphQL, web-admin transport, and coverage graph facts.
- `docs/verification-plan.xml` - Add auth scenarios, test files, commands, log markers, CORS/CSRF, e2e, codegen drift, and coverage replacement gates.
- `docs/operational-packets.xml` - Update only if current packet text needs auth-specific execution/checkpoint fields.

### Delete

- `apps/api/internal/middleware/auth.go` - Delete the unused bearer-shaped placeholder after `admin_auth.go` owns the real admin session boundary.
- `apps/api/internal/middleware/auth_test.go` - Delete stale bearer-placeholder tests.

### Do Not Modify

- `apps/web/src/**` - Public web stays REST-only and public.
- `apps/web-admin/src/pages/**` - No login UI or route guard UI in this wave.
- `apps/api/internal/service/user_service.go` public behavior - Existing `UserService` behavior remains the reference user CRUD domain.
- `docker/docker-compose.dev.yml` - No runtime image change is required for backend-only auth.

## Execution Discipline

- Use a clean worktree before execution. If an isolated worktree is needed, create it at execution time with `superpowers:using-git-worktrees`.
- Run TDD within each task: write failing tests first, verify they fail for the expected reason, implement the minimal production code, then rerun the focused tests.
- Run codegen before tests that import generated GraphQL or sqlc types.
- Do not run full `bun run verify:coverage` during early implementation. Save broad coverage/e2e gates for final closeout.
- Do not commit after each tiny code step. Commit only after a verified task group keeps code, generated artifacts, docs, and coverage policy internally consistent.
- Every new or meaningfully edited governed source, test, schema, query, migration, config, tooling, or durable doc file must carry file-local GRACE markup unless the file format is generated or JSON.
- When a literal replacement snippet in this plan omits a file-local header, treat the snippet as the behavior body only: preserve the existing `FILE` / `MODULE_CONTRACT` / `MODULE_MAP` / `CHANGE_SUMMARY` block or add one before changing behavior. Do not remove semantic markup while applying snippets.

## Task 0: Preflight And Evidence Log

**Files:**

- Create: `.tasks/web-admin-backend-auth/verification.md`

- [ ] **Step 1: Inspect target-file dirtiness**

Run:

```bash
git status --short -- apps/api apps/web-admin libs/graphql tools/coverage docs .env.example .tasks/web-admin-backend-auth
```

Expected:

```text
# no output for target files, or only known changes from the current execution worktree
```

If a target file is already dirty, run:

```bash
git diff -- <path>
```

Expected:

```text
# inspect the diff manually
```

If the existing diff is unrelated to this plan, pause before editing that file.

- [ ] **Step 2: Create verification evidence log**

Create `.tasks/web-admin-backend-auth/verification.md`:

```markdown
<!-- FILE: .tasks/web-admin-backend-auth/verification.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Record verification evidence for the backend-only web-admin auth wave. -->
<!--   SCOPE: Captures focused commands, generated drift checks, e2e auth setup, secret-redaction evidence, coverage decisions, GRACE validation, and final status; excludes durable architecture contracts. -->
<!--   DEPENDS: docs/superpowers/plans/2026-06-07-web-admin-backend-auth.md, docs/superpowers/specs/2026-06-07-web-admin-backend-auth-design.md, apps/api, apps/web-admin, libs/graphql/schema. -->
<!--   LINKS: M-API / V-M-API / M-GRAPHQL-SCHEMA / V-M-GRAPHQL-SCHEMA / M-WEB-ADMIN / V-M-WEB-ADMIN / M-COVERAGE-GATE / V-M-COVERAGE-GATE. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Command Evidence - Lists each focused verification command and result. -->
<!--   Auth Evidence - Records seed, login, me, logout, createAdmin, CORS, CSRF, and e2e cookie evidence. -->
<!--   Generated Coverage - Records codegen drift and generated-file replacement gates. -->
<!--   Secret Redaction - Records evidence that credentials, raw cookies, raw sessions, and hashes are not logged. -->
<!--   Final Status - States whether the wave is ready for handoff. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added verification log for backend-only web-admin auth implementation. -->
<!-- END_CHANGE_SUMMARY -->

# Web-admin Backend Auth Verification

## Command Evidence

| Command                                                                                                                                                                | Result     | Notes                  |
| ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------- | ---------------------- | --------------------- | ----------- | ------------- | ---------- | ----- | ---- | ---------- | ------------- | ------- | --------------------- |
| `cd apps/api && go test ./internal/appconfig -run TestConfig_Admin`                                                                                                    | NOT RUN    | Filled during Task 1.  |
| `cd apps/api && go test ./internal/repository/postgres -run TestAdminRepo`                                                                                             | NOT RUN    | Filled during Task 2.  |
| `cd apps/api && go test ./internal/repository/redis -run TestAdminSessionStore`                                                                                        | NOT RUN    | Filled during Task 3.  |
| `cd apps/api && go test ./internal/service -run TestAdminAuth`                                                                                                         | NOT RUN    | Filled during Task 4.  |
| `cd apps/api && go test ./internal/middleware -run 'TestAdmin                                                                                                          | TestCORS'` | NOT RUN                | Filled during Task 5. |
| `bunx nx run graphql:validate`                                                                                                                                         | NOT RUN    | Filled during Task 6.  |
| `bunx nx run api:codegen`                                                                                                                                              | NOT RUN    | Filled during Task 6.  |
| `cd apps/api && go test ./internal/graph -run 'Test(Admin                                                                                                              | LoginAdmin | LogoutAdmin            | Me                    | CreateAdmin | ProtectedUser | CreateUser | Users | User | UpdateUser | DeleteUser)'` | NOT RUN | Filled during Task 6. |
| `bunx nx run web-admin:codegen`                                                                                                                                        | NOT RUN    | Filled during Task 7.  |
| `bunx nx run web-admin:typecheck`                                                                                                                                      | NOT RUN    | Filled during Task 7.  |
| `bunx nx test web-admin`                                                                                                                                               | NOT RUN    | Filled during Task 7.  |
| `bunx nx run web-admin:e2e`                                                                                                                                            | NOT RUN    | Filled during Task 8.  |
| `bunx nx run codegen:validate`                                                                                                                                         | NOT RUN    | Filled during Task 8.  |
| `bunx nx test api`                                                                                                                                                     | NOT RUN    | Filled during Task 10. |
| `bunx nx build api`                                                                                                                                                    | NOT RUN    | Filled during Task 10. |
| `bun run test:coverage`                                                                                                                                                | NOT RUN    | Filled during Task 10. |
| `bun run verify:coverage`                                                                                                                                              | NOT RUN    | Filled during Task 10. |
| `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml` | NOT RUN    | Filled during Task 10. |
| `grace lint --path .`                                                                                                                                                  | NOT RUN    | Filled during Task 10. |

## Auth Evidence

| Scenario                                                           | Result  | Evidence |
| ------------------------------------------------------------------ | ------- | -------- |
| Empty table seed creates first admin                               | NOT RUN |          |
| Non-empty table seed does not create env admin                     | NOT RUN |          |
| Login sets httpOnly cookie                                         | NOT RUN |          |
| `me` returns admin for valid active session                        | NOT RUN |          |
| `me` returns null for missing or invalid session                   | NOT RUN |          |
| Logout clears cookie and Redis session                             | NOT RUN |          |
| `createAdmin` requires authenticated admin                         | NOT RUN |          |
| Public REST `/api/users` remains public                            | NOT RUN |          |
| Admin GraphQL rejects public web origin                            | NOT RUN |          |
| E2E user CRUD runs after `loginAdmin` cookie                       | NOT RUN |          |
| Empty-table API startup requires and uses `ADMIN_INITIAL_*` env    | NOT RUN |          |
| Non-empty-table API startup does not require `ADMIN_INITIAL_*` env | NOT RUN |          |

## Generated Coverage

| Generated path                                                       | Replacement gate                                          | Result  |
| -------------------------------------------------------------------- | --------------------------------------------------------- | ------- |
| `apps/api/internal/repository/postgres/generated/admin_users.sql.go` | API codegen, API build, admin repository integration test | NOT RUN |

## Secret Redaction

| Secret class     | Evidence command or test                                                                        | Result                                   |
| ---------------- | ----------------------------------------------------------------------------------------------- | ---------------------------------------- | ------- |
| Passwords        | `cd apps/api && go test ./internal/service -run TestAdminAuthService_LogsMarkersWithoutSecrets` | NOT RUN                                  |
| Password hashes  | `cd apps/api && go test ./internal/service -run TestAdminAuthService_LogsMarkersWithoutSecrets` | NOT RUN                                  |
| Raw cookies      | `cd apps/api && go test ./internal/middleware -run 'TestAdmin(SessionMiddleware                 | OriginGuard)\_LogsMarkerWithoutSecrets'` | NOT RUN |
| Raw session ids  | `cd apps/api && go test ./internal/middleware -run 'TestAdmin(SessionMiddleware                 | OriginGuard)\_LogsMarkerWithoutSecrets'` | NOT RUN |
| Credential input | `cd apps/api && go test ./internal/graph -run TestCreateAdmin_ReturnsAuthErrorWithoutPrincipal` | NOT RUN                                  |

## Final Status

NOT READY
```

Run:

```bash
git diff --check -- .tasks/web-admin-backend-auth/verification.md
```

Expected:

```text
# no output
```

## Task 1: Admin And Session Config

**Files:**

- Modify: `apps/api/internal/appconfig/config.go`
- Modify: `apps/api/internal/appconfig/config_test.go`
- Modify: `apps/api/config/config.yml`
- Modify: `.env.example`

- [ ] **Step 1: Write failing app config tests**

Before adding the new tests, delete or replace any stale JWT-placeholder test such as `TestConfig_AuthJWTSecretRequired`; the backend auth path is cookie/session based and no longer has a bearer JWT contract.

Rewrite the existing `TestConfig_PaginationValid` to call `validConfig()` from this task instead of constructing removed `AuthConfig` / `JWTSecret` fields:

```go
func TestConfig_PaginationValid(t *testing.T) {
	cfg := validConfig()
	err := config.Validate(cfg)
	require.NoError(t, err)
}
```

Append these tests to `apps/api/internal/appconfig/config_test.go`:

```go
func TestConfig_AdminSeedRequiredWhenBootstrapNeeded(t *testing.T) {
	cfg := validConfig()
	cfg.Admin = appconfig.AdminConfig{}
	err := appconfig.ValidateAdminBootstrap(cfg, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "initial admin email is required")
	assert.Contains(t, err.Error(), "initial admin password is required")
	assert.Contains(t, err.Error(), "initial admin name is required")
}

func TestConfig_AdminSeedNotRequiredWhenAdminExists(t *testing.T) {
	cfg := validConfig()
	cfg.Admin = appconfig.AdminConfig{}
	err := appconfig.ValidateAdminBootstrap(cfg, false)
	require.NoError(t, err)
}

func TestConfig_AdminBootstrapRequiresEnvKeysWhenTableEmpty(t *testing.T) {
	env := map[string]string{
		"ADMIN_INITIAL_EMAIL": "admin@example.com",
		"ADMIN_INITIAL_NAME":  "Template Admin",
	}
	err := appconfig.ValidateAdminBootstrapEnv(func(key string) (string, bool) {
		value, ok := env[key]
		return value, ok
	}, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ADMIN_INITIAL_PASSWORD is required when admin_users is empty")
}

func TestConfig_AdminBootstrapEnvNotRequiredWhenAdminExists(t *testing.T) {
	err := appconfig.ValidateAdminBootstrapEnv(func(key string) (string, bool) {
		return "", false
	}, false)
	require.NoError(t, err)
}

func TestConfig_AdminBootstrapRejectsPlaceholderSeedOutsideDevelopment(t *testing.T) {
	cfg := validConfig()
	cfg.Server.Env = "production"
	cfg.Admin.InitialEmail = "admin@example.com"
	cfg.Admin.InitialPassword = "ChangeMeAdmin123!"
	cfg.Admin.InitialName = "Template Admin"
	err := appconfig.ValidateAdminBootstrap(cfg, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "initial admin bootstrap values must not use example placeholders outside development")
}

func TestConfig_AdminEnvOverlayHydratesEnvOnlyAdminFields(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
server:
  port: 8080
  env: staging
postgres:
  host: localhost
  port: 5432
  user: app
  db: test
redis:
  host: localhost
  port: 6379
admin:
  origins:
    - "http://localhost:3100"
admin_session:
  cookie_name: web_admin_session
  ttl: 168h
  cookie_secure: true
  same_site: Lax
pagination:
  default_page_size: 20
  max_page_size: 100
`), 0o600))
	t.Setenv("ADMIN_INITIAL_EMAIL", "ops-admin@example.test")
	t.Setenv("ADMIN_INITIAL_PASSWORD", "StrongPassword123!")
	t.Setenv("ADMIN_INITIAL_NAME", "Ops Admin")
	t.Setenv("ADMIN_ORIGINS", "https://admin.example.com,https://admin2.example.com")
	t.Setenv("ADMIN_SESSION_KEY_SECRET", "real-session-key-secret")

	cfg, err := config.Load[appconfig.Config](config.Options{ConfigPath: configPath})
	require.NoError(t, err)
	require.NoError(t, appconfig.ApplyAdminEnvOverlay(&cfg, os.LookupEnv))
	require.NoError(t, appconfig.ApplyAdminDefaults(&cfg))

	assert.Equal(t, "ops-admin@example.test", cfg.Admin.InitialEmail)
	assert.Equal(t, "StrongPassword123!", cfg.Admin.InitialPassword)
	assert.Equal(t, "Ops Admin", cfg.Admin.InitialName)
	assert.Equal(t, []string{"https://admin.example.com", "https://admin2.example.com"}, cfg.Admin.Origins)
	assert.Equal(t, "real-session-key-secret", cfg.AdminSession.KeySecret)
}

func TestConfig_AdminSessionDefaultsAndValidation(t *testing.T) {
	cfg := validConfig()
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.NoError(t, err)
	assert.Equal(t, "web_admin_session", cfg.AdminSession.CookieName)
	assert.Equal(t, "168h", cfg.AdminSession.TTL.String())
	assert.Equal(t, "auto", cfg.AdminSession.CookieSecure)
	assert.Equal(t, "Lax", cfg.AdminSession.SameSite)
	assert.Equal(t, []string{"http://localhost:3100", "http://127.0.0.1:3100"}, cfg.Admin.Origins)
}

func TestConfig_AdminSessionRejectsInvalidSecureMode(t *testing.T) {
	cfg := validConfig()
	cfg.AdminSession.CookieSecure = "sometimes"
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin session cookie secure must be auto, true, or false")
}

func TestConfig_AdminSessionRejectsProductionInsecureCookie(t *testing.T) {
	cfg := validConfig()
	cfg.Server.Env = "production"
	cfg.AdminSession.CookieSecure = "false"
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "production admin session cookie secure must be true or auto")
}

func TestConfig_AdminSessionRejectsSameSiteNoneWithoutSecureCookie(t *testing.T) {
	cfg := validConfig()
	cfg.AdminSession.CookieSecure = "false"
	cfg.AdminSession.SameSite = "None"
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin session SameSite=None requires a secure cookie")
}

func TestConfig_AdminSessionRejectsSameSiteNoneAutoOutsideProduction(t *testing.T) {
	cfg := validConfig()
	cfg.Server.Env = "staging"
	cfg.AdminSession.CookieSecure = "auto"
	cfg.AdminSession.SameSite = "None"
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin session SameSite=None requires a secure cookie")
}

func TestConfig_AdminSessionRejectsPlaceholderSecretOutsideDevelopment(t *testing.T) {
	cfg := validConfig()
	cfg.Server.Env = "production"
	cfg.AdminSession.KeySecret = "change-me-session-key-secret"
	err := appconfig.ApplyAdminDefaults(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin session key secret must be provided by environment outside development")
}

func validConfig() appconfig.Config {
	return appconfig.Config{
		Server:   config.ServerConfig{Port: 8080, Env: "development"},
		Postgres: config.PostgresConfig{Host: "localhost", Port: 5432, User: "app", DB: "test"},
		Redis:    config.RedisConfig{Host: "localhost", Port: 6379},
		Admin: appconfig.AdminConfig{
			InitialEmail:    "admin@example.com",
			InitialPassword: "StrongPassword123!",
			InitialName:     "Template Admin",
			Origins:         []string{"http://localhost:3100", "http://127.0.0.1:3100"},
		},
		AdminSession: appconfig.AdminSessionConfig{
			CookieName:   "web_admin_session",
			TTL:          168 * time.Hour,
			CookieSecure: "auto",
			SameSite:     "Lax",
			KeySecret:    "dev-session-key-secret",
		},
		Pagination: appconfig.PaginationConfig{DefaultPageSize: 20, MaxPageSize: 100},
	}
}
```

Update the imports in `apps/api/internal/appconfig/config_test.go`:

```go
import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/appconfig"
	"monorepo-template/libs/go/config"
)
```

- [ ] **Step 2: Run the config tests and verify they fail**

Run:

```bash
cd apps/api
go test ./internal/appconfig -run TestConfig_Admin -count=1
```

Expected:

```text
FAIL
undefined: appconfig.AdminConfig
undefined: appconfig.AdminSessionConfig
undefined: appconfig.ApplyAdminEnvOverlay
undefined: appconfig.ValidateAdminBootstrap
undefined: appconfig.ApplyAdminDefaults
```

- [ ] **Step 3: Implement admin config structs and helpers**

Replace `apps/api/internal/appconfig/config.go` with:

```go
// FILE: apps/api/internal/appconfig/config.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define API service configuration, admin bootstrap checks, and admin session defaults.
//   SCOPE: Server, logging, PostgreSQL, Redis, admin seed, admin session, and pagination config; excludes secret loading implementation and persistence behavior.
//   DEPENDS: libs/go/config.
//   LINKS: M-API / V-M-API.
//   ROLE: CONFIG
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminConfig - Env-backed initial-admin and web-admin origin settings.
//   AdminSessionConfig - Cookie/session settings for Redis-backed admin sessions.
//   Config - Full API configuration shape.
//   ApplyAdminDefaults - Applies non-secret defaults and validates session security settings.
//   ValidateAdminBootstrapEnv - Proves first-admin seed values are present in env when the table is empty.
//   ValidateAdminBootstrap - Validates first-admin config values only when bootstrap is required.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Replaced placeholder bearer auth config with admin bootstrap and session config.
// END_CHANGE_SUMMARY

package appconfig

import (
	"fmt"
	"strings"
	"time"

	"monorepo-template/libs/go/config"
)

// AdminConfig holds bootstrap and browser-origin settings for web-admin authentication.
type AdminConfig struct {
	InitialEmail    string   `mapstructure:"initial_email"`
	InitialPassword string   `mapstructure:"initial_password"`
	InitialName     string   `mapstructure:"initial_name"`
	Origins         []string `mapstructure:"origins"`
}

// AdminSessionConfig holds cookie and Redis-key settings for admin sessions.
type AdminSessionConfig struct {
	CookieName   string        `mapstructure:"cookie_name"`
	TTL          time.Duration `mapstructure:"ttl"`
	CookieSecure string        `mapstructure:"cookie_secure"`
	SameSite     string        `mapstructure:"same_site"`
	KeySecret    string        `mapstructure:"key_secret"`
}

// PaginationConfig holds pagination defaults.
type PaginationConfig struct {
	DefaultPageSize int `mapstructure:"default_page_size" validate:"gt=0"`
	MaxPageSize     int `mapstructure:"max_page_size"     validate:"gt=0"`
}

// Config is the API service configuration.
type Config struct {
	Server       config.ServerConfig   `mapstructure:"server"`
	Log          config.LogConfig      `mapstructure:"log"`
	Postgres     config.PostgresConfig `mapstructure:"postgres"`
	Redis        config.RedisConfig    `mapstructure:"redis"`
	Admin        AdminConfig           `mapstructure:"admin"`
	AdminSession AdminSessionConfig    `mapstructure:"admin_session"`
	Pagination   PaginationConfig      `mapstructure:"pagination"`
}

// ApplyAdminEnvOverlay hydrates env-only admin secrets and comma-separated origin overrides.
func ApplyAdminEnvOverlay(cfg *Config, lookup func(string) (string, bool)) error {
	if value, ok := envString(lookup, "ADMIN_INITIAL_EMAIL"); ok {
		cfg.Admin.InitialEmail = value
	}
	if value, ok := envString(lookup, "ADMIN_INITIAL_PASSWORD"); ok {
		cfg.Admin.InitialPassword = value
	}
	if value, ok := envString(lookup, "ADMIN_INITIAL_NAME"); ok {
		cfg.Admin.InitialName = value
	}
	if value, ok := envString(lookup, "ADMIN_ORIGINS"); ok {
		cfg.Admin.Origins = splitCSV(value)
	}
	if value, ok := envString(lookup, "ADMIN_SESSION_COOKIE_NAME"); ok {
		cfg.AdminSession.CookieName = value
	}
	if value, ok := envString(lookup, "ADMIN_SESSION_TTL"); ok {
		ttl, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("admin session ttl is invalid: %w", err)
		}
		cfg.AdminSession.TTL = ttl
	}
	if value, ok := envString(lookup, "ADMIN_SESSION_COOKIE_SECURE"); ok {
		cfg.AdminSession.CookieSecure = value
	}
	if value, ok := envString(lookup, "ADMIN_SESSION_SAME_SITE"); ok {
		cfg.AdminSession.SameSite = value
	}
	if value, ok := envString(lookup, "ADMIN_SESSION_KEY_SECRET"); ok {
		cfg.AdminSession.KeySecret = value
	}
	return nil
}

// ApplyAdminDefaults fills safe local defaults and validates enum-like values.
func ApplyAdminDefaults(cfg *Config) error {
	env := strings.ToLower(strings.TrimSpace(cfg.Server.Env))
	if cfg.AdminSession.CookieName == "" {
		cfg.AdminSession.CookieName = "web_admin_session"
	}
	if cfg.AdminSession.TTL == 0 {
		cfg.AdminSession.TTL = 168 * time.Hour
	}
	if cfg.AdminSession.CookieSecure == "" {
		cfg.AdminSession.CookieSecure = "auto"
	}
	if cfg.AdminSession.SameSite == "" {
		cfg.AdminSession.SameSite = "Lax"
	}
	if len(cfg.Admin.Origins) == 0 {
		cfg.Admin.Origins = []string{"http://localhost:3100", "http://127.0.0.1:3100"}
	}
	if strings.TrimSpace(cfg.AdminSession.KeySecret) == "" {
		return fmt.Errorf("admin session key secret is required")
	}
	if cfg.AdminSession.CookieSecure != "auto" && cfg.AdminSession.CookieSecure != "true" && cfg.AdminSession.CookieSecure != "false" {
		return fmt.Errorf("admin session cookie secure must be auto, true, or false")
	}
	if cfg.AdminSession.SameSite != "Lax" && cfg.AdminSession.SameSite != "Strict" && cfg.AdminSession.SameSite != "None" {
		return fmt.Errorf("admin session same site must be Lax, Strict, or None")
	}
	effectiveSecure := cfg.AdminSession.CookieSecure == "true" || (cfg.AdminSession.CookieSecure == "auto" && env == "production")
	if env == "production" && cfg.AdminSession.CookieSecure == "false" {
		return fmt.Errorf("production admin session cookie secure must be true or auto")
	}
	if cfg.AdminSession.SameSite == "None" && !effectiveSecure {
		return fmt.Errorf("admin session SameSite=None requires a secure cookie")
	}
	if env != "" && env != "development" && isPlaceholderAdminSessionSecret(cfg.AdminSession.KeySecret) {
		return fmt.Errorf("admin session key secret must be provided by environment outside development")
	}
	return nil
}

// ValidateAdminBootstrapEnv ensures the first admin identity is supplied by env when needed.
func ValidateAdminBootstrapEnv(lookup func(string) (string, bool), adminTableEmpty bool) error {
	if !adminTableEmpty {
		return nil
	}
	required := []string{"ADMIN_INITIAL_EMAIL", "ADMIN_INITIAL_PASSWORD", "ADMIN_INITIAL_NAME"}
	var missing []string
	for _, key := range required {
		value, ok := lookup(key)
		if !ok || strings.TrimSpace(value) == "" {
			missing = append(missing, key+" is required when admin_users is empty")
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf(strings.Join(missing, "; "))
	}
	return nil
}

// ValidateAdminBootstrap enforces first-admin env only when no admins exist.
func ValidateAdminBootstrap(cfg Config, adminTableEmpty bool) error {
	if !adminTableEmpty {
		return nil
	}
	env := strings.ToLower(strings.TrimSpace(cfg.Server.Env))
	var missing []string
	if strings.TrimSpace(cfg.Admin.InitialEmail) == "" {
		missing = append(missing, "initial admin email is required")
	}
	if strings.TrimSpace(cfg.Admin.InitialPassword) == "" {
		missing = append(missing, "initial admin password is required")
	}
	if strings.TrimSpace(cfg.Admin.InitialName) == "" {
		missing = append(missing, "initial admin name is required")
	}
	if len(missing) > 0 {
		return fmt.Errorf(strings.Join(missing, "; "))
	}
	if env != "" && env != "development" && isPlaceholderInitialAdmin(cfg.Admin) {
		return fmt.Errorf("initial admin bootstrap values must not use example placeholders outside development")
	}
	return nil
}

func isPlaceholderAdminSessionSecret(value string) bool {
	trimmed := strings.TrimSpace(value)
	return trimmed == "change-me-session-key-secret" || trimmed == "dev-session-key-secret"
}

func isPlaceholderInitialAdmin(admin AdminConfig) bool {
	return strings.EqualFold(strings.TrimSpace(admin.InitialEmail), "admin@example.com") ||
		strings.TrimSpace(admin.InitialPassword) == "ChangeMeAdmin123!" ||
		strings.TrimSpace(admin.InitialName) == "Template Admin"
}

func envString(lookup func(string) (string, bool), key string) (string, bool) {
	value, ok := lookup(key)
	if !ok {
		return "", false
	}
	return strings.TrimSpace(value), strings.TrimSpace(value) != ""
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
```

Update the call site in `apps/api/cmd/server/main.go` after config load during Task 8:

```go
if err := appconfig.ApplyAdminEnvOverlay(&cfg, os.LookupEnv); err != nil {
	fmt.Fprintf(os.Stderr, "failed to apply admin env overlay: %v\n", err)
	os.Exit(1)
}
if err := appconfig.ApplyAdminDefaults(&cfg); err != nil {
	fmt.Fprintf(os.Stderr, "failed to apply admin config defaults: %v\n", err)
	os.Exit(1)
}
```

- [ ] **Step 4: Update local config files**

In `apps/api/config/config.yml`, replace the old `auth:` block with:

```yaml
admin:
  origins:
    - 'http://localhost:3100'
    - 'http://127.0.0.1:3100'

admin_session:
  cookie_name: web_admin_session
  ttl: 168h
  cookie_secure: auto
  same_site: Lax
```

Do not put `initial_email`, `initial_password`, `initial_name`, or `key_secret` in committed YAML. `ADMIN_INITIAL_*` and `ADMIN_SESSION_KEY_SECRET` must come from real environment variables or a local `.env` file loaded by the API runtime. This preserves the env-only first-admin contract and prevents checked-in bootstrap credentials from seeding an empty table.

The values in `.env.example` are documentation-only placeholders. Production and staging must provide non-placeholder values through their secret/env system; `ApplyAdminDefaults` must reject placeholder `ADMIN_SESSION_KEY_SECRET` values outside development.

In `.env.example`, replace the old auth section with:

```dotenv
# Admin bootstrap and session auth
ADMIN_INITIAL_EMAIL=admin@example.com
ADMIN_INITIAL_PASSWORD=ChangeMeAdmin123!
ADMIN_INITIAL_NAME=Template Admin
ADMIN_ORIGINS=http://localhost:3100,http://127.0.0.1:3100
ADMIN_SESSION_COOKIE_NAME=web_admin_session
ADMIN_SESSION_TTL=168h
ADMIN_SESSION_COOKIE_SECURE=auto
ADMIN_SESSION_SAME_SITE=Lax
ADMIN_SESSION_KEY_SECRET=change-me-session-key-secret
```

- [ ] **Step 5: Run focused config tests**

Run:

```bash
cd apps/api
go test ./internal/appconfig -run 'TestConfig_(Admin|Pagination)' -count=1
```

Expected:

```text
ok  	monorepo-template/apps/api/internal/appconfig
```

## Task 2: Admin Schema, SQLC, And PostgreSQL Repository

**Files:**

- Create: `apps/api/internal/repository/postgres/migrations/00079_admin_users.sql`
- Create: `apps/api/internal/repository/postgres/queries/admin_users.sql`
- Create: `apps/api/internal/repository/postgres/admin_repo.go`
- Create: `apps/api/internal/repository/postgres/admin_repo_test.go`
- Modify: `apps/api/sqlc.yaml`
- Generate: `apps/api/internal/repository/postgres/generated/admin_users.sql.go`

- [ ] **Step 1: Create the admin migration**

Create `apps/api/internal/repository/postgres/migrations/00079_admin_users.sql`:

```sql
-- +goose Up
CREATE TABLE admin_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(32) NOT NULL DEFAULT 'ADMIN',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_admin_users_email_lower ON admin_users (LOWER(email));
CREATE INDEX idx_admin_users_active ON admin_users (is_active);

-- +goose Down
DROP TABLE IF EXISTS admin_users;

-- FILE: apps/api/internal/repository/postgres/migrations/00079_admin_users.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add the admin_users table for web-admin authentication.
--   SCOPE: Admin identity schema, lower-case unique email enforcement, active flag, role, and timestamps; excludes session storage and public users schema.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00001_init.sql.
--   LINKS: M-API / V-M-API.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   admin_users - Stores web-admin identities separate from public reference users.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added admin identity schema for backend auth.
-- END_CHANGE_SUMMARY
```

- [ ] **Step 2: Create admin sqlc queries**

Create `apps/api/internal/repository/postgres/queries/admin_users.sql`:

```sql
-- name: CountAdminUsers :one
SELECT COUNT(*) FROM admin_users;

-- name: CreateAdminUser :one
INSERT INTO admin_users (email, name, password_hash, role, is_active)
VALUES (LOWER($1), $2, $3, $4, TRUE)
RETURNING id, email, name, password_hash, role, is_active, created_at, updated_at;

-- name: GetAdminUserByEmail :one
SELECT id, email, name, password_hash, role, is_active, created_at, updated_at
FROM admin_users
WHERE LOWER(email) = LOWER($1);

-- name: GetAdminUserByID :one
SELECT id, email, name, password_hash, role, is_active, created_at, updated_at
FROM admin_users
WHERE id = $1;

-- FILE: apps/api/internal/repository/postgres/queries/admin_users.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc admin user queries used by the PostgreSQL AdminRepo adapter.
--   SCOPE: Count, create, and identity lookup for admin_users; excludes public users persistence and Redis sessions.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00079_admin_users.sql.
--   LINKS: M-API / V-M-API.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   CountAdminUsers - Counts admin identities for bootstrap-only seed.
--   CreateAdminUser - Inserts one normalized active admin.
--   GetAdminUserByEmail - Fetches one admin by case-insensitive email.
--   GetAdminUserByID - Fetches one admin by UUID.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added admin_users sqlc query contract.
-- END_CHANGE_SUMMARY
```

- [ ] **Step 3: Run sqlc and verify generated admin query code appears**

Run:

```bash
bunx nx run api:codegen
test -f apps/api/internal/repository/postgres/generated/admin_users.sql.go
rg -n "CreateAdminUser|GetAdminUserByEmail|CountAdminUsers" apps/api/internal/repository/postgres/generated
```

Expected:

```text
Successfully ran target codegen for project api
apps/api/internal/repository/postgres/generated/admin_users.sql.go:...
apps/api/internal/repository/postgres/generated/querier.go:...
```

- [ ] **Step 4: Write failing admin repository integration tests**

Create `apps/api/internal/repository/postgres/admin_repo_test.go`:

```go
// FILE: apps/api/internal/repository/postgres/admin_repo_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify PostgreSQL admin repository behavior against the goose-managed test database.
//   SCOPE: Admin create/count/get, case-insensitive duplicate handling, active flag mapping, and safe destructive setup; excludes service validation and Redis sessions.
//   DEPENDS: apps/api/internal/repository/postgres, apps/api/internal/service, apps/api/internal/testinfra, docker/docker-compose.test.yml.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   adminTestPool - Applies migrations, enforces safe test DSN, and truncates admin_users.
//   TestAdminRepo_* - Real database coverage for admin identity persistence.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin repository integration coverage.
// END_CHANGE_SUMMARY

package postgres_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	postgresrepo "monorepo-template/apps/api/internal/repository/postgres"
	"monorepo-template/apps/api/internal/service"
	"monorepo-template/apps/api/internal/testinfra"
)

func adminTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := testinfra.PostgresDSN()
	testinfra.RequireSafePostgresDSN(t, dsn)
	if err := postgresrepo.RunMigrations(dsn, zap.NewNop()); err != nil {
		if !testinfra.CoverageGateEnabled() {
			t.Skipf("postgres integration database is unavailable: %v", err)
		}
		require.NoError(t, err)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	_, err = pool.Exec(context.Background(), `TRUNCATE admin_users RESTART IDENTITY CASCADE`)
	require.NoError(t, err)
	return pool
}

func TestAdminRepo_CreateCountAndLookup(t *testing.T) {
	ctx := context.Background()
	repo := postgresrepo.NewAdminRepo(adminTestPool(t))

	count, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	created, err := repo.Create(ctx, service.CreateAdminInput{
		Email:        "Admin@Example.COM",
		Name:         "Template Admin",
		PasswordHash: "$2a$10$hashed",
		Role:         service.AdminRoleAdmin,
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "admin@example.com", created.Email)
	assert.True(t, created.IsActive)

	foundByEmail, err := repo.GetByEmail(ctx, "ADMIN@example.com")
	require.NoError(t, err)
	require.NotNil(t, foundByEmail)
	assert.Equal(t, created.ID, foundByEmail.ID)

	foundByID, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, foundByID)
	assert.Equal(t, created.Email, foundByID.Email)

	count, err = repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestAdminRepo_CreateDuplicateEmailCaseInsensitive(t *testing.T) {
	ctx := context.Background()
	repo := postgresrepo.NewAdminRepo(adminTestPool(t))
	input := service.CreateAdminInput{
		Email:        "duplicate@example.com",
		Name:         "First",
		PasswordHash: "$2a$10$hashed",
		Role:         service.AdminRoleAdmin,
	}
	_, err := repo.Create(ctx, input)
	require.NoError(t, err)

	input.Email = "DUPLICATE@example.com"
	_, err = repo.Create(ctx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrAdminDuplicateEmail)
}

func TestAdminRepo_GetMissingReturnsNil(t *testing.T) {
	repo := postgresrepo.NewAdminRepo(adminTestPool(t))

	byEmail, err := repo.GetByEmail(context.Background(), "missing@example.com")
	require.NoError(t, err)
	assert.Nil(t, byEmail)

	byID, err := repo.GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, byID)
}
```

- [ ] **Step 5: Run admin repository tests and verify they fail**

Run:

```bash
cd apps/api
go test ./internal/repository/postgres -run TestAdminRepo -count=1
```

Expected:

```text
FAIL
undefined: postgresrepo.NewAdminRepo
undefined: service.CreateAdminInput
undefined: service.AdminRoleAdmin
undefined: service.ErrAdminDuplicateEmail
```

- [ ] **Step 6: Add admin service types needed by the repository**

Create the first version of `apps/api/internal/service/admin_auth.go`:

```go
// FILE: apps/api/internal/service/admin_auth.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Own transport-neutral web-admin authentication service contracts and behavior.
//   SCOPE: Admin domain types, repository/session interfaces, bootstrap seed, login, current admin, create admin, and logout; excludes PostgreSQL, Redis, GraphQL, and HTTP cookie adapters.
//   DEPENDS: context, strings, time, bcrypt, libs/go/logger.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminRoleAdmin - Initial admin role value.
//   Admin - Public admin identity returned to GraphQL and context consumers.
//   CreateAdminInput - Repository/service admin creation input.
//   AdminRepository - Persistence boundary for admin identities.
//   AdminSessionStore - Session boundary for Redis-backed sessions.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin auth domain contracts.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"errors"
)

const AdminRoleAdmin = "ADMIN"

type Admin struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
	Role         string
	IsActive     bool
	CreatedAt    string
	UpdatedAt    string
}

type CreateAdminInput struct {
	Email        string
	Name         string
	PasswordHash string
	Role         string
}

type AdminRepository interface {
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateAdminInput) (*Admin, error)
	GetByEmail(ctx context.Context, email string) (*Admin, error)
	GetByID(ctx context.Context, id string) (*Admin, error)
}

type AdminSessionStore interface{}

var (
	ErrAdminDuplicateEmail = errors.New("admin duplicate email")
	ErrAdminNotFound       = errors.New("admin not found")
	ErrAdminAuth           = errors.New("admin authentication failed")
	ErrAdminValidation     = errors.New("admin validation failed")
)
```

- [ ] **Step 7: Implement PostgreSQL admin repository**

Create `apps/api/internal/repository/postgres/admin_repo.go`:

```go
// FILE: apps/api/internal/repository/postgres/admin_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Adapt sqlc-generated PostgreSQL admin_users queries to the service.AdminRepository contract.
//   SCOPE: Admin count, create, email/id lookup, row mapping, UUID parsing, and duplicate email mapping; excludes password hashing and session storage.
//   DEPENDS: apps/api/internal/repository/postgres/generated, github.com/jackc/pgx/v5, apps/api/internal/service.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   NewAdminRepo - Constructs the production admin repository from a pgx pool.
//   AdminRepo.Count - Counts admin identities for seed bootstrap.
//   AdminRepo.Create - Inserts one normalized admin and maps duplicate email.
//   AdminRepo.GetByEmail - Reads one admin by normalized email.
//   AdminRepo.GetByID - Reads one admin by UUID.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added PostgreSQL admin repository adapter.
// END_CHANGE_SUMMARY

package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/repository/postgres/generated"
	"monorepo-template/apps/api/internal/service"
)

type AdminRepo struct {
	queries generated.Querier
}

func NewAdminRepo(pool *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{queries: generated.New(pool)}
}

func (r *AdminRepo) Count(ctx context.Context) (int, error) {
	count, err := r.queries.CountAdminUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("AdminRepo.Count: %w", err)
	}
	return int(count), nil
}

func (r *AdminRepo) Create(ctx context.Context, input service.CreateAdminInput) (*service.Admin, error) {
	row, err := r.queries.CreateAdminUser(ctx, generated.CreateAdminUserParams{
		Email:        input.Email,
		Name:         input.Name,
		PasswordHash: input.PasswordHash,
		Role:         input.Role,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, service.ErrAdminDuplicateEmail
		}
		return nil, fmt.Errorf("AdminRepo.Create: %w", err)
	}
	return adminFromFields(row.ID, row.Email, row.Name, row.PasswordHash, row.Role, row.IsActive, row.CreatedAt, row.UpdatedAt), nil
}

func (r *AdminRepo) GetByEmail(ctx context.Context, email string) (*service.Admin, error) {
	row, err := r.queries.GetAdminUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("AdminRepo.GetByEmail: %w", err)
	}
	return adminFromFields(row.ID, row.Email, row.Name, row.PasswordHash, row.Role, row.IsActive, row.CreatedAt, row.UpdatedAt), nil
}

func (r *AdminRepo) GetByID(ctx context.Context, id string) (*service.Admin, error) {
	adminID, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.GetByID: invalid admin id: %w", err)
	}
	row, err := r.queries.GetAdminUserByID(ctx, adminID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("AdminRepo.GetByID: %w", err)
	}
	return adminFromFields(row.ID, row.Email, row.Name, row.PasswordHash, row.Role, row.IsActive, row.CreatedAt, row.UpdatedAt), nil
}

func adminFromFields(id pgtype.UUID, email string, name string, passwordHash string, role string, isActive bool, createdAt pgtype.Timestamptz, updatedAt pgtype.Timestamptz) *service.Admin {
	return &service.Admin{
		ID:           id.String(),
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
		Role:         role,
		IsActive:     isActive,
		CreatedAt:    adminFormatTimestamp(createdAt),
		UpdatedAt:    adminFormatTimestamp(updatedAt),
	}
}

func adminFormatTimestamp(value pgtype.Timestamptz) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format(time.RFC3339Nano)
}
```

- [ ] **Step 8: Run admin repository tests**

Run:

```bash
cd apps/api
go test ./internal/repository/postgres -run TestAdminRepo -count=1
```

Expected:

```text
ok  	monorepo-template/apps/api/internal/repository/postgres
```

## Task 3: Redis Session Store And Cookie Helpers

**Files:**

- Create: `apps/api/internal/repository/redis/admin_session_store.go`
- Create: `apps/api/internal/repository/redis/admin_session_store_test.go`
- Create: `apps/api/internal/middleware/admin_auth.go`
- Create: `apps/api/internal/middleware/admin_auth_test.go`

- [ ] **Step 1: Write failing Redis session store tests**

Create `apps/api/internal/repository/redis/admin_session_store_test.go`:

```go
// FILE: apps/api/internal/repository/redis/admin_session_store_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify Redis-backed admin session storage.
//   SCOPE: Session create/read/delete, expiry, HMAC-derived Redis keys, and unavailable Redis skip semantics; excludes HTTP cookies and GraphQL.
//   DEPENDS: apps/api/internal/repository/redis, apps/api/internal/testinfra, github.com/redis/go-redis/v9.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestAdminSessionStore_* - Real Redis coverage for admin sessions.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin session store coverage.
// END_CHANGE_SUMMARY

package redis_test

import (
	"context"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	redisrepo "monorepo-template/apps/api/internal/repository/redis"
	"monorepo-template/apps/api/internal/testinfra"
)

func adminRedisClient(t *testing.T) *redisrepo.Client {
	t.Helper()
	client, err := redisrepo.New(testinfra.RedisConfig(t), zap.NewNop())
	if err != nil && !testinfra.CoverageGateEnabled() {
		t.Skipf("redis integration service is unavailable: %v", err)
	}
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	cleanupAdminSessionKeys(t, client.RDB)
	return client
}

func cleanupAdminSessionKeys(t *testing.T, rdb *goredis.Client) {
	t.Helper()
	ctx := context.Background()
	iter := rdb.Scan(ctx, 0, "admin_session:*", 100).Iterator()
	for iter.Next(ctx) {
		require.NoError(t, rdb.Del(ctx, iter.Val()).Err())
	}
	require.NoError(t, iter.Err())
}

func TestAdminSessionStore_CreateReadDelete(t *testing.T) {
	ctx := context.Background()
	store := redisrepo.NewAdminSessionStore(adminRedisClient(t).RDB, []byte("test-key-secret"), time.Hour)

	sessionID, err := store.Create(ctx, "admin-1")
	require.NoError(t, err)
	require.NotEmpty(t, sessionID)

	adminID, err := store.Get(ctx, sessionID)
	require.NoError(t, err)
	assert.Equal(t, "admin-1", adminID)

	require.NoError(t, store.Delete(ctx, sessionID))
	adminID, err = store.Get(ctx, sessionID)
	require.NoError(t, err)
	assert.Empty(t, adminID)
}

func TestAdminSessionStore_DoesNotUseRawSessionIDAsKey(t *testing.T) {
	ctx := context.Background()
	client := adminRedisClient(t)
	store := redisrepo.NewAdminSessionStore(client.RDB, []byte("test-key-secret"), time.Hour)
	sessionID, err := store.Create(ctx, "admin-1")
	require.NoError(t, err)

	keys, err := client.RDB.Keys(ctx, "admin_session:*").Result()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	assert.NotContains(t, keys[0], sessionID)
}

func TestAdminSessionStore_Expires(t *testing.T) {
	ctx := context.Background()
	store := redisrepo.NewAdminSessionStore(adminRedisClient(t).RDB, []byte("test-key-secret"), time.Millisecond)
	sessionID, err := store.Create(ctx, "admin-1")
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		adminID, err := store.Get(ctx, sessionID)
		require.NoError(t, err)
		return adminID == ""
	}, 250*time.Millisecond, 10*time.Millisecond)
}
```

- [ ] **Step 2: Run Redis session tests and verify they fail**

Run:

```bash
cd apps/api
go test ./internal/repository/redis -run TestAdminSessionStore -count=1
```

Expected:

```text
FAIL
undefined: redisrepo.NewAdminSessionStore
```

- [ ] **Step 3: Extend session service interface**

Update `AdminSessionStore` in `apps/api/internal/service/admin_auth.go`:

```go
type AdminSessionStore interface {
	Create(ctx context.Context, adminID string) (string, error)
	Get(ctx context.Context, sessionID string) (string, error)
	Delete(ctx context.Context, sessionID string) error
}
```

- [ ] **Step 4: Implement Redis session store**

Create `apps/api/internal/repository/redis/admin_session_store.go`:

```go
// FILE: apps/api/internal/repository/redis/admin_session_store.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Store web-admin opaque sessions in Redis using HMAC-derived keys.
//   SCOPE: Session id generation, Redis key derivation, create/read/delete, and TTL; excludes cookies, GraphQL, and admin identity lookup.
//   DEPENDS: crypto/rand, crypto/hmac, crypto/sha256, encoding/base64, github.com/redis/go-redis/v9.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   NewAdminSessionStore - Constructs a Redis-backed admin session store.
//   AdminSessionStore.Create - Creates an opaque browser session id and Redis entry.
//   AdminSessionStore.Get - Resolves a session id to an admin id.
//   AdminSessionStore.Delete - Revokes a session id.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Redis admin session store.
// END_CHANGE_SUMMARY

package redis

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type AdminSessionStore struct {
	rdb       *goredis.Client
	keySecret []byte
	ttl       time.Duration
}

func NewAdminSessionStore(rdb *goredis.Client, keySecret []byte, ttl time.Duration) *AdminSessionStore {
	return &AdminSessionStore{rdb: rdb, keySecret: keySecret, ttl: ttl}
}

func (s *AdminSessionStore) Create(ctx context.Context, adminID string) (string, error) {
	sessionID, err := randomSessionID()
	if err != nil {
		return "", fmt.Errorf("AdminSessionStore.Create: generate session id: %w", err)
	}
	if err := s.rdb.Set(ctx, s.key(sessionID), adminID, s.ttl).Err(); err != nil {
		return "", fmt.Errorf("AdminSessionStore.Create: redis set: %w", err)
	}
	return sessionID, nil
}

func (s *AdminSessionStore) Get(ctx context.Context, sessionID string) (string, error) {
	if sessionID == "" {
		return "", nil
	}
	adminID, err := s.rdb.Get(ctx, s.key(sessionID)).Result()
	if errors.Is(err, goredis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("AdminSessionStore.Get: redis get: %w", err)
	}
	return adminID, nil
}

func (s *AdminSessionStore) Delete(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}
	if err := s.rdb.Del(ctx, s.key(sessionID)).Err(); err != nil {
		return fmt.Errorf("AdminSessionStore.Delete: redis del: %w", err)
	}
	return nil
}

func (s *AdminSessionStore) key(sessionID string) string {
	mac := hmac.New(sha256.New, s.keySecret)
	_, _ = mac.Write([]byte(sessionID))
	return "admin_session:" + hex.EncodeToString(mac.Sum(nil))
}

func randomSessionID() (string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}
```

- [ ] **Step 5: Run Redis session tests**

Run:

```bash
cd apps/api
go test ./internal/repository/redis -run TestAdminSessionStore -count=1
```

Expected:

```text
ok  	monorepo-template/apps/api/internal/repository/redis
```

- [ ] **Step 6: Write failing cookie/context tests**

Create `apps/api/internal/middleware/admin_auth_test.go`:

```go
// FILE: apps/api/internal/middleware/admin_auth_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify admin auth cookie, context, and protected-operation helpers.
//   SCOPE: AdminPrincipal context storage, session cookie set/clear attributes, session-cookie principal hydration, and protected guard outcomes; excludes Redis implementation and GraphQL resolver logic.
//   DEPENDS: apps/api/internal/middleware, apps/api/internal/service, httptest.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestAdminCookie_* - Verifies session cookie attributes.
//   TestAdminPrincipal_* - Verifies context principal storage.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin auth middleware helper coverage.
// END_CHANGE_SUMMARY

package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/apps/api/internal/service"
	"monorepo-template/libs/go/logger"
)

func TestAdminPrincipal_ContextRoundTrip(t *testing.T) {
	principal := middleware.AdminPrincipal{ID: "admin-1", Email: "admin@example.com", Name: "Admin", Role: "ADMIN"}
	ctx := middleware.ContextWithAdminPrincipal(context.Background(), principal)

	found, ok := middleware.GetAdminPrincipal(ctx)

	require.True(t, ok)
	assert.Equal(t, principal, found)
}

func TestAdminCookie_SetAndClear(t *testing.T) {
	rec := httptest.NewRecorder()
	cfg := middleware.AdminCookieConfig{
		Name:     "web_admin_session",
		Path:     "/graphql",
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	middleware.SetAdminSessionCookie(rec, cfg, "session-id")
	setCookie := rec.Result().Cookies()[0]
	assert.Equal(t, "web_admin_session", setCookie.Name)
	assert.Equal(t, "session-id", setCookie.Value)
	assert.True(t, setCookie.HttpOnly)
	assert.Equal(t, "/graphql", setCookie.Path)
	assert.Equal(t, http.SameSiteLaxMode, setCookie.SameSite)

	clearRec := httptest.NewRecorder()
	middleware.ClearAdminSessionCookie(clearRec, cfg)
	clearCookie := clearRec.Result().Cookies()[0]
	assert.Equal(t, "web_admin_session", clearCookie.Name)
	assert.Empty(t, clearCookie.Value)
	assert.Equal(t, -1, clearCookie.MaxAge)
	assert.True(t, clearCookie.Expires.Before(time.Now()))
	assert.Equal(t, "/graphql", clearCookie.Path)
}

func TestAdminSessionMiddleware_LoadsPrincipalFromCookie(t *testing.T) {
	resolver := &fakeAdminSessionResolver{admin: &service.Admin{
		ID: "admin-1", Email: "admin@example.com", Name: "Admin", Role: service.AdminRoleAdmin, IsActive: true,
		CreatedAt: "2026-06-07T00:00:00Z", UpdatedAt: "2026-06-07T00:00:00Z",
	}}
	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.AddCookie(&http.Cookie{Name: "web_admin_session", Value: "session-1"})
	rec := httptest.NewRecorder()
	var principal middleware.AdminPrincipal
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		principal, ok = middleware.GetAdminPrincipal(r.Context())
		require.True(t, ok)
		assert.Equal(t, "session-1", middleware.AdminSessionIDFromContext(r.Context()))
	})

	middleware.AdminSessionMiddleware(resolver, "web_admin_session")(next).ServeHTTP(rec, req)

	assert.Equal(t, "admin@example.com", principal.Email)
	assert.Equal(t, "Admin", principal.Name)
}

func TestAdminSessionMiddleware_DoesNotSetPrincipalForInactiveSession(t *testing.T) {
	resolver := &fakeAdminSessionResolver{admin: nil}
	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.AddCookie(&http.Cookie{Name: "web_admin_session", Value: "session-1"})
	rec := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := middleware.GetAdminPrincipal(r.Context())
		assert.False(t, ok)
		assert.Equal(t, "session-1", middleware.AdminSessionIDFromContext(r.Context()))
	})

	middleware.AdminSessionMiddleware(resolver, "web_admin_session")(next).ServeHTTP(rec, req)
}

func TestAdminSessionMiddleware_LogsMarkerWithoutSecrets(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	resolver := &fakeAdminSessionResolver{admin: &service.Admin{ID: "admin-1", Email: "admin@example.com", Role: service.AdminRoleAdmin, IsActive: true}}
	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req = req.WithContext(logger.WithContext(req.Context(), zap.New(core)))
	req.AddCookie(&http.Cookie{Name: "web_admin_session", Value: "raw-session-id"})
	rec := httptest.NewRecorder()

	middleware.AdminSessionMiddleware(resolver, "web_admin_session")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(rec, req)

	joined := logText(logs.All())
	assert.Contains(t, joined, "[AdminAuth][session][BLOCK_VALIDATE_SESSION]")
	assert.NotContains(t, joined, "raw-session-id")
	assert.NotContains(t, joined, "web_admin_session")
}

type fakeAdminSessionResolver struct {
	admin *service.Admin
}

func (f *fakeAdminSessionResolver) CurrentAdmin(ctx context.Context, sessionID string) (*service.Admin, error) {
	if sessionID != "session-1" {
		return nil, nil
	}
	return f.admin, nil
}

func logText(entries []observer.LoggedEntry) string {
	var out strings.Builder
	for _, entry := range entries {
		out.WriteString(entry.Message)
		out.WriteString("\n")
		for _, field := range entry.Context {
			out.WriteString(field.Key)
			out.WriteString("=")
			out.WriteString(field.String)
			out.WriteString("\n")
		}
	}
	return out.String()
}
```

- [ ] **Step 7: Run cookie/context tests and verify they fail**

Run:

```bash
cd apps/api
go test ./internal/middleware -run TestAdmin -count=1
```

Expected:

```text
FAIL
undefined: middleware.AdminPrincipal
undefined: middleware.ContextWithAdminPrincipal
undefined: middleware.AdminCookieConfig
undefined: middleware.AdminSessionMiddleware
undefined: middleware.AdminSessionIDFromContext
```

- [ ] **Step 8: Implement cookie/context helpers**

Create `apps/api/internal/middleware/admin_auth.go`:

```go
// FILE: apps/api/internal/middleware/admin_auth.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide admin principal context helpers, session lookup middleware, and session cookie helpers.
//   SCOPE: AdminPrincipal context storage, session-cookie lookup through the auth service boundary, cookie set/clear behavior, and cookie config; excludes Redis implementation and GraphQL resolver decisions.
//   DEPENDS: context, net/http, apps/api/internal/service.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminPrincipal - Request-scoped authenticated admin identity.
//   ContextWithAdminPrincipal - Stores an admin principal in context.
//   GetAdminPrincipal - Reads an admin principal from context.
//   ContextWithAdminSessionID - Stores the raw session id only inside request context for logout.
//   AdminSessionIDFromContext - Reads the request-scoped session id for logout.
//   AdminSessionMiddleware - Resolves session cookies into request-scoped admin context.
//   SetAdminSessionCookie - Sets the httpOnly session cookie.
//   ClearAdminSessionCookie - Clears the session cookie.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin principal and cookie helpers.
// END_CHANGE_SUMMARY

package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/service"
	"monorepo-template/libs/go/logger"
)

type adminPrincipalKey struct{}
type adminSessionIDKey struct{}

type AdminPrincipal struct {
	ID        string
	Email     string
	Name      string
	Role      string
	CreatedAt string
	UpdatedAt string
}

type AdminSessionResolver interface {
	CurrentAdmin(ctx context.Context, sessionID string) (*service.Admin, error)
}

type AdminCookieConfig struct {
	Name     string
	Path     string
	MaxAge   int
	Secure   bool
	SameSite http.SameSite
}

func ContextWithAdminPrincipal(ctx context.Context, principal AdminPrincipal) context.Context {
	return context.WithValue(ctx, adminPrincipalKey{}, principal)
}

func GetAdminPrincipal(ctx context.Context) (AdminPrincipal, bool) {
	principal, ok := ctx.Value(adminPrincipalKey{}).(AdminPrincipal)
	return principal, ok
}

func ContextWithAdminSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, adminSessionIDKey{}, sessionID)
}

func AdminSessionIDFromContext(ctx context.Context) string {
	sessionID, _ := ctx.Value(adminSessionIDKey{}).(string)
	return sessionID
}

func AdminSessionMiddleware(resolver AdminSessionResolver, cookieName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			cookie, err := r.Cookie(cookieName)
			if err != nil || cookie.Value == "" {
				next.ServeHTTP(w, r)
				return
			}
			log.Debug("[AdminAuth][session][BLOCK_VALIDATE_SESSION] session cookie present")
			admin, err := resolver.CurrentAdmin(r.Context(), cookie.Value)
			if err != nil {
				log.Error("[AdminAuth][session][BLOCK_VALIDATE_SESSION] session lookup failed", zap.Error(err))
				http.Error(w, "admin session lookup failed", http.StatusInternalServerError)
				return
			}
			ctx := ContextWithAdminSessionID(r.Context(), cookie.Value)
			if admin != nil {
				ctx = ContextWithAdminPrincipal(ctx, AdminPrincipal{
					ID: admin.ID, Email: admin.Email, Name: admin.Name, Role: admin.Role,
					CreatedAt: admin.CreatedAt, UpdatedAt: admin.UpdatedAt,
				})
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SetAdminSessionCookie(w http.ResponseWriter, cfg AdminCookieConfig, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.Name,
		Value:    sessionID,
		Path:     cfg.Path,
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		MaxAge:   cfg.MaxAge,
	})
}

func ClearAdminSessionCookie(w http.ResponseWriter, cfg AdminCookieConfig) {
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.Name,
		Value:    "",
		Path:     cfg.Path,
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		MaxAge:   -1,
		Expires:  time.Now().Add(-time.Hour),
	})
}
```

- [ ] **Step 9: Run middleware helper tests**

Run:

```bash
cd apps/api
go test ./internal/middleware -run TestAdmin -count=1
```

Expected:

```text
ok  	monorepo-template/apps/api/internal/middleware
```

## Task 4: Admin Auth Service

**Files:**

- Modify: `apps/api/internal/service/admin_auth.go`
- Create: `apps/api/internal/service/admin_auth_test.go`

- [ ] **Step 1: Write failing admin service tests**

Create `apps/api/internal/service/admin_auth_test.go` with fakes and tests:

```go
// FILE: apps/api/internal/service/admin_auth_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify transport-neutral admin auth service behavior.
//   SCOPE: Bootstrap seed, password validation, login, current admin, create admin, logout, duplicate mapping, inactive admin denial, session revocation, and secret-redaction marker paths; excludes PostgreSQL, Redis, HTTP cookies, and GraphQL.
//   DEPENDS: internal/service, bcrypt, testify.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   fakeAdminRepo - Admin repository test double.
//   fakeAdminSessions - Admin session store test double.
//   TestAdminAuthService_* - Service behavior coverage.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin auth service coverage.
// END_CHANGE_SUMMARY

package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"monorepo-template/apps/api/internal/service"
	"monorepo-template/libs/go/logger"
)

func TestAdminAuthService_SeedCreatesOnlyWhenEmpty(t *testing.T) {
	repo := newFakeAdminRepo()
	svc := service.NewAdminAuthService(repo, newFakeAdminSessions())

	created, err := svc.SeedInitialAdmin(context.Background(), service.InitialAdminInput{
		Email:    "Admin@Example.COM",
		Name:     "Template Admin",
		Password: "StrongPassword123!",
	})

	require.NoError(t, err)
	assert.True(t, created)
	assert.Equal(t, 1, repo.count)
	assert.Equal(t, "admin@example.com", repo.created[0].Email)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(repo.created[0].PasswordHash), []byte("StrongPassword123!")))
}

func TestAdminAuthService_SeedNoopsWhenAnyAdminExists(t *testing.T) {
	repo := newFakeAdminRepo()
	repo.count = 1
	svc := service.NewAdminAuthService(repo, newFakeAdminSessions())

	created, err := svc.SeedInitialAdmin(context.Background(), service.InitialAdminInput{
		Email:    "second@example.com",
		Name:     "Second",
		Password: "StrongPassword123!",
	})

	require.NoError(t, err)
	assert.False(t, created)
	assert.Empty(t, repo.created)
}

func TestAdminAuthService_LoginCreatesSession(t *testing.T) {
	repo := newFakeAdminRepo()
	hash, err := bcrypt.GenerateFromPassword([]byte("StrongPassword123!"), bcrypt.MinCost)
	require.NoError(t, err)
	repo.adminsByEmail["admin@example.com"] = &service.Admin{
		ID: "admin-1", Email: "admin@example.com", Name: "Admin", PasswordHash: string(hash), Role: service.AdminRoleAdmin, IsActive: true,
	}
	sessions := newFakeAdminSessions()
	svc := service.NewAdminAuthService(repo, sessions)

	result, err := svc.Login(context.Background(), service.LoginAdminInput{Email: "ADMIN@example.com", Password: "StrongPassword123!"})

	require.NoError(t, err)
	assert.Equal(t, "session-1", result.SessionID)
	assert.Equal(t, "admin-1", sessions.createdFor)
	assert.Equal(t, "admin@example.com", result.Admin.Email)
}

func TestAdminAuthService_LoginRejectsWrongPasswordAndInactiveAdmin(t *testing.T) {
	repo := newFakeAdminRepo()
	hash, err := bcrypt.GenerateFromPassword([]byte("StrongPassword123!"), bcrypt.MinCost)
	require.NoError(t, err)
	repo.adminsByEmail["admin@example.com"] = &service.Admin{
		ID: "admin-1", Email: "admin@example.com", PasswordHash: string(hash), Role: service.AdminRoleAdmin, IsActive: false,
	}
	svc := service.NewAdminAuthService(repo, newFakeAdminSessions())

	_, err = svc.Login(context.Background(), service.LoginAdminInput{Email: "admin@example.com", Password: "wrong-password"})
	require.ErrorIs(t, err, service.ErrAdminAuth)

	_, err = svc.Login(context.Background(), service.LoginAdminInput{Email: "admin@example.com", Password: "StrongPassword123!"})
	require.ErrorIs(t, err, service.ErrAdminAuth)
}

func TestAdminAuthService_CurrentAdmin(t *testing.T) {
	repo := newFakeAdminRepo()
	repo.adminsByID["admin-1"] = &service.Admin{ID: "admin-1", Email: "admin@example.com", Role: service.AdminRoleAdmin, IsActive: true}
	sessions := newFakeAdminSessions()
	sessions.sessions["session-1"] = "admin-1"
	svc := service.NewAdminAuthService(repo, sessions)

	admin, err := svc.CurrentAdmin(context.Background(), "session-1")
	require.NoError(t, err)
	require.NotNil(t, admin)
	assert.Equal(t, "admin@example.com", admin.Email)

	admin, err = svc.CurrentAdmin(context.Background(), "missing-session")
	require.NoError(t, err)
	assert.Nil(t, admin)
}

func TestAdminAuthService_CurrentAdminReturnsNilForInactiveAdmin(t *testing.T) {
	repo := newFakeAdminRepo()
	repo.adminsByID["admin-1"] = &service.Admin{ID: "admin-1", Email: "admin@example.com", Role: service.AdminRoleAdmin, IsActive: false}
	sessions := newFakeAdminSessions()
	sessions.sessions["session-1"] = "admin-1"
	svc := service.NewAdminAuthService(repo, sessions)

	admin, err := svc.CurrentAdmin(context.Background(), "session-1")

	require.NoError(t, err)
	assert.Nil(t, admin)
}

func TestAdminAuthService_CreateAdminRequiresActor(t *testing.T) {
	svc := service.NewAdminAuthService(newFakeAdminRepo(), newFakeAdminSessions())

	_, err := svc.CreateAdmin(context.Background(), nil, service.NewAdminInput{
		Email: "new@example.com", Name: "New", Password: "StrongPassword123!",
	})

	require.ErrorIs(t, err, service.ErrAdminAuth)
}

func TestAdminAuthService_CreateAdminHashesPassword(t *testing.T) {
	repo := newFakeAdminRepo()
	actor := &service.Admin{ID: "admin-1", Email: "admin@example.com", Role: service.AdminRoleAdmin, IsActive: true}
	svc := service.NewAdminAuthService(repo, newFakeAdminSessions())

	created, err := svc.CreateAdmin(context.Background(), actor, service.NewAdminInput{
		Email: "New@Example.COM", Name: "New Admin", Password: "StrongPassword123!",
	})

	require.NoError(t, err)
	assert.Equal(t, "new@example.com", created.Email)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(repo.created[0].PasswordHash), []byte("StrongPassword123!")))
}

func TestAdminAuthService_LogoutDeletesSession(t *testing.T) {
	sessions := newFakeAdminSessions()
	svc := service.NewAdminAuthService(newFakeAdminRepo(), sessions)

	require.NoError(t, svc.Logout(context.Background(), "session-1"))
	assert.Equal(t, "session-1", sessions.deleted)
}

func TestAdminAuthService_LogsMarkersWithoutSecrets(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	ctx := logger.WithContext(context.Background(), zap.New(core))
	repo := newFakeAdminRepo()
	hash, err := bcrypt.GenerateFromPassword([]byte("StrongPassword123!"), bcrypt.MinCost)
	require.NoError(t, err)
	repo.adminsByEmail["admin@example.com"] = &service.Admin{
		ID: "admin-1", Email: "admin@example.com", PasswordHash: string(hash), Role: service.AdminRoleAdmin, IsActive: true,
	}
	sessions := newFakeAdminSessions()
	svc := service.NewAdminAuthService(repo, sessions)

	_, _ = svc.Login(ctx, service.LoginAdminInput{Email: "admin@example.com", Password: "StrongPassword123!"})
	_ = svc.Logout(ctx, "session-1")

	joined := logText(logs.All())
	assert.Contains(t, joined, "[AdminAuth][login][BLOCK_VERIFY_CREDENTIALS]")
	assert.Contains(t, joined, "[AdminAuth][session][BLOCK_VALIDATE_SESSION]")
	assert.Contains(t, joined, "[AdminAuth][logout][BLOCK_REVOKE_SESSION]")
	assert.NotContains(t, joined, "admin@example.com")
	assert.NotContains(t, joined, "StrongPassword123!")
	assert.NotContains(t, joined, string(hash))
	assert.NotContains(t, joined, "session-1")
}

func logText(entries []observer.LoggedEntry) string {
	var out strings.Builder
	for _, entry := range entries {
		out.WriteString(entry.Message)
		out.WriteString("\n")
		for _, field := range entry.Context {
			out.WriteString(field.Key)
			out.WriteString("=")
			out.WriteString(field.String)
			out.WriteString("\n")
		}
	}
	return out.String()
}

type fakeAdminRepo struct {
	count        int
	adminsByID   map[string]*service.Admin
	adminsByEmail map[string]*service.Admin
	created      []service.CreateAdminInput
	err          error
}

func newFakeAdminRepo() *fakeAdminRepo {
	return &fakeAdminRepo{adminsByID: map[string]*service.Admin{}, adminsByEmail: map[string]*service.Admin{}}
}

func (r *fakeAdminRepo) Count(ctx context.Context) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.count, nil
}

func (r *fakeAdminRepo) Create(ctx context.Context, input service.CreateAdminInput) (*service.Admin, error) {
	if r.err != nil {
		return nil, r.err
	}
	if _, exists := r.adminsByEmail[input.Email]; exists {
		return nil, service.ErrAdminDuplicateEmail
	}
	r.created = append(r.created, input)
	r.count++
	admin := &service.Admin{ID: "admin-" + input.Email, Email: input.Email, Name: input.Name, PasswordHash: input.PasswordHash, Role: input.Role, IsActive: true}
	r.adminsByEmail[input.Email] = admin
	r.adminsByID[admin.ID] = admin
	return admin, nil
}

func (r *fakeAdminRepo) GetByEmail(ctx context.Context, email string) (*service.Admin, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.adminsByEmail[email], nil
}

func (r *fakeAdminRepo) GetByID(ctx context.Context, id string) (*service.Admin, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.adminsByID[id], nil
}

type fakeAdminSessions struct {
	sessions   map[string]string
	createdFor string
	deleted    string
	err        error
}

func newFakeAdminSessions() *fakeAdminSessions {
	return &fakeAdminSessions{sessions: map[string]string{}}
}

func (s *fakeAdminSessions) Create(ctx context.Context, adminID string) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	s.createdFor = adminID
	s.sessions["session-1"] = adminID
	return "session-1", nil
}

func (s *fakeAdminSessions) Get(ctx context.Context, sessionID string) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.sessions[sessionID], nil
}

func (s *fakeAdminSessions) Delete(ctx context.Context, sessionID string) error {
	if s.err != nil && !errors.Is(s.err, service.ErrAdminNotFound) {
		return s.err
	}
	s.deleted = sessionID
	delete(s.sessions, sessionID)
	return nil
}
```

- [ ] **Step 2: Run service tests and verify they fail**

Run:

```bash
cd apps/api
go test ./internal/service -run TestAdminAuth -count=1
```

Expected:

```text
FAIL
undefined: service.NewAdminAuthService
undefined: service.InitialAdminInput
undefined: service.LoginAdminInput
undefined: service.NewAdminInput
```

- [ ] **Step 3: Implement admin auth service behavior**

Extend `apps/api/internal/service/admin_auth.go` with these public types and methods:

```go
type InitialAdminInput struct {
	Email    string
	Name     string
	Password string
}

type LoginAdminInput struct {
	Email    string
	Password string
}

type LoginAdminResult struct {
	Admin     *Admin
	SessionID string
}

type NewAdminInput struct {
	Email    string
	Name     string
	Password string
}

type AdminAuthService struct {
	repo     AdminRepository
	sessions AdminSessionStore
}

func NewAdminAuthService(repo AdminRepository, sessions AdminSessionStore) *AdminAuthService {
	return &AdminAuthService{repo: repo, sessions: sessions}
}

func (s *AdminAuthService) SeedInitialAdmin(ctx context.Context, input InitialAdminInput) (bool, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, nil
	}
	if err := validateAdminInput(input.Email, input.Name, input.Password); err != nil {
		return false, err
	}
	hash, err := hashPassword(input.Password)
	if err != nil {
		return false, err
	}
	_, err = s.repo.Create(ctx, CreateAdminInput{
		Email:        normalizeEmail(input.Email),
		Name:         strings.TrimSpace(input.Name),
		PasswordHash: hash,
		Role:         AdminRoleAdmin,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *AdminAuthService) Login(ctx context.Context, input LoginAdminInput) (*LoginAdminResult, error) {
	log := logger.FromContext(ctx)
	log.Info("[AdminAuth][login][BLOCK_VERIFY_CREDENTIALS] login attempt")
	admin, err := s.repo.GetByEmail(ctx, normalizeEmail(input.Email))
	if err != nil {
		return nil, err
	}
	if admin == nil || !admin.IsActive {
		log.Warn("[AdminAuth][login][BLOCK_VERIFY_CREDENTIALS] credential check failed")
		return nil, ErrAdminAuth
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(input.Password)); err != nil {
		log.Warn("[AdminAuth][login][BLOCK_VERIFY_CREDENTIALS] credential check failed")
		return nil, ErrAdminAuth
	}
	sessionID, err := s.sessions.Create(ctx, admin.ID)
	if err != nil {
		return nil, err
	}
	log.Info("[AdminAuth][session][BLOCK_VALIDATE_SESSION] session created", zap.String("admin_id", admin.ID))
	return &LoginAdminResult{Admin: publicAdmin(admin), SessionID: sessionID}, nil
}

func (s *AdminAuthService) CurrentAdmin(ctx context.Context, sessionID string) (*Admin, error) {
	log := logger.FromContext(ctx)
	log.Debug("[AdminAuth][session][BLOCK_VALIDATE_SESSION] session lookup")
	adminID, err := s.sessions.Get(ctx, sessionID)
	if err != nil {
		log.Error("[AdminAuth][session][BLOCK_VALIDATE_SESSION] session lookup failed", zap.Error(err))
		return nil, err
	}
	if adminID == "" {
		return nil, nil
	}
	admin, err := s.repo.GetByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if admin == nil || !admin.IsActive {
		return nil, nil
	}
	return publicAdmin(admin), nil
}

func (s *AdminAuthService) CreateAdmin(ctx context.Context, actor *Admin, input NewAdminInput) (*Admin, error) {
	if actor == nil || !actor.IsActive {
		return nil, ErrAdminAuth
	}
	if err := validateAdminInput(input.Email, input.Name, input.Password); err != nil {
		return nil, err
	}
	hash, err := hashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	admin, err := s.repo.Create(ctx, CreateAdminInput{
		Email:        normalizeEmail(input.Email),
		Name:         strings.TrimSpace(input.Name),
		PasswordHash: hash,
		Role:         AdminRoleAdmin,
	})
	if err != nil {
		return nil, err
	}
	return publicAdmin(admin), nil
}

func (s *AdminAuthService) Logout(ctx context.Context, sessionID string) error {
	logger.FromContext(ctx).Info("[AdminAuth][logout][BLOCK_REVOKE_SESSION] logout requested")
	return s.sessions.Delete(ctx, sessionID)
}
```

Also add these helpers and imports:

```go
import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"go.uber.org/zap"

	"monorepo-template/libs/go/logger"
)

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateAdminInput(email string, name string, password string) error {
	if normalizeEmail(email) == "" {
		return fmt.Errorf("%w: email is required", ErrAdminValidation)
	}
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name is required", ErrAdminValidation)
	}
	if len(password) < 12 {
		return fmt.Errorf("%w: password must be at least 12 characters", ErrAdminValidation)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash admin password: %w", err)
	}
	return string(hashed), nil
}

func publicAdmin(admin *Admin) *Admin {
	if admin == nil {
		return nil
	}
	copy := *admin
	copy.PasswordHash = ""
	return &copy
}
```

- [ ] **Step 4: Run admin service tests**

Run:

```bash
cd apps/api
go test ./internal/service -run TestAdminAuth -count=1
```

Expected:

```text
ok  	monorepo-template/apps/api/internal/service
```

## Task 5: CORS, CSRF, And Admin Request Boundary

**Files:**

- Modify: `apps/api/internal/middleware/cors.go`
- Modify: `apps/api/internal/middleware/cors_test.go`
- Create: `apps/api/internal/middleware/admin_origin.go`
- Create: `apps/api/internal/middleware/admin_origin_test.go`

- [ ] **Step 1: Write failing CORS tests**

Append to `apps/api/internal/middleware/cors_test.go`:

```go
func TestCORS_CredentialedRejectsWildcardOrigin(t *testing.T) {
	handler := middleware.CORS(middleware.CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORS_CredentialedAdminOrigin(t *testing.T) {
	handler := middleware.CORS(middleware.CORSConfig{
		AllowedOrigins:   []string{"http://localhost:3100"},
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, "http://localhost:3100", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
}
```

- [ ] **Step 2: Run CORS tests and verify they fail**

Run:

```bash
cd apps/api
go test ./internal/middleware -run TestCORS_Credentialed -count=1
```

Expected:

```text
FAIL
unknown field AllowCredentials in struct literal of type middleware.CORSConfig
```

- [ ] **Step 3: Implement credentialed CORS behavior**

Modify `apps/api/internal/middleware/cors.go`:

```go
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowOrigin := ""
			for _, allowed := range cfg.AllowedOrigins {
				if allowed == "*" && !cfg.AllowCredentials {
					allowOrigin = "*"
					break
				}
				if allowed == origin {
					allowOrigin = origin
					break
				}
			}
			if allowOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
				if cfg.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
```

- [ ] **Step 4: Write failing admin origin guard tests**

Create `apps/api/internal/middleware/admin_origin_test.go`:

```go
// FILE: apps/api/internal/middleware/admin_origin_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify admin GraphQL origin and CSRF guard behavior.
//   SCOPE: Allows configured web-admin origins and rejects public or disallowed browser origins for protected/session-mutating GraphQL requests; excludes GraphQL parsing.
//   DEPENDS: apps/api/internal/middleware, httptest.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestAdminOriginGuard_* - Origin allow and deny coverage.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin origin guard coverage.
// END_CHANGE_SUMMARY

package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/libs/go/logger"
)

func TestAdminOriginGuard_AllowsConfiguredOrigin(t *testing.T) {
	handler := middleware.AdminOriginGuard([]string{"http://localhost:3100"})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"mutation { createAdmin(input:{email:\"a@example.com\",name:\"A\",password:\"StrongPassword123!\"}) { __typename } }"}`))
	req.Header.Set("Origin", "http://localhost:3100")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminOriginGuard_RejectsPublicWebOrigin(t *testing.T) {
	handler := middleware.AdminOriginGuard([]string{"http://localhost:3100"})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"mutation { createAdmin(input:{email:\"a@example.com\",name:\"A\",password:\"StrongPassword123!\"}) { __typename } }"}`))
	req.Header.Set("Origin", "http://localhost:3101")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAdminOriginGuard_RejectsMissingOriginForUnsafeRequest(t *testing.T) {
	handler := middleware.AdminOriginGuard([]string{"http://localhost:3100"})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"mutation { logoutAdmin { __typename } }"}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAdminOriginGuard_LogsMarkerWithoutSecrets(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	handler := middleware.AdminOriginGuard([]string{"http://localhost:3100"})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	body := `{"query":"mutation { loginAdmin(input:{email:\"admin@example.com\",password:\"StrongPassword123!\"}) { __typename } }"}`
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(body))
	req = req.WithContext(logger.WithContext(req.Context(), zap.New(core)))
	req.Header.Set("Origin", "http://localhost:3101")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	joined := logText(logs.All())
	assert.Contains(t, joined, "[AdminAuth][csrf][BLOCK_VALIDATE_ORIGIN]")
	assert.NotContains(t, joined, "admin@example.com")
	assert.NotContains(t, joined, "StrongPassword123!")
	assert.NotContains(t, joined, "loginAdmin")
	assert.NotContains(t, joined, body)
}
```

- [ ] **Step 5: Implement admin origin guard**

Create `apps/api/internal/middleware/admin_origin.go`:

```go
// FILE: apps/api/internal/middleware/admin_origin.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Enforce the web-admin browser-origin boundary for credentialed admin GraphQL requests.
//   SCOPE: Strict Origin/Referer allowlisting for unsafe browser requests; excludes CORS header emission and GraphQL auth decisions.
//   DEPENDS: net/http, net/url.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminOriginGuard - Rejects unsafe admin GraphQL browser requests from non-admin origins.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin origin guard.
// END_CHANGE_SUMMARY

package middleware

import (
	"net/http"
	"net/url"

	"go.uber.org/zap"

	"monorepo-template/libs/go/logger"
)

func AdminOriginGuard(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := map[string]struct{}{}
	for _, origin := range allowedOrigins {
		allowed[origin] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			if r.Method == http.MethodGet || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}
			origin := r.Header.Get("Origin")
			if origin == "" {
				origin = originFromReferer(r.Header.Get("Referer"))
			}
			if origin == "" {
				log.Warn("[AdminAuth][csrf][BLOCK_VALIDATE_ORIGIN] missing origin")
				http.Error(w, "admin origin is required", http.StatusForbidden)
				return
			}
			if _, ok := allowed[origin]; !ok {
				log.Warn("[AdminAuth][csrf][BLOCK_VALIDATE_ORIGIN] rejected origin", zap.String("origin", origin))
				http.Error(w, "admin origin is not allowed", http.StatusForbidden)
				return
			}
			log.Debug("[AdminAuth][csrf][BLOCK_VALIDATE_ORIGIN] allowed origin", zap.String("origin", origin))
			next.ServeHTTP(w, r)
		})
	}
}

func originFromReferer(value string) string {
	if value == "" {
		return ""
	}
	parsed, err := url.Parse(value)
	if err != nil {
		return ""
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	return parsed.Scheme + "://" + parsed.Host
}
```

- [ ] **Step 6: Run middleware tests**

Run:

```bash
cd apps/api
go test ./internal/middleware -run 'TestCORS|TestAdmin' -count=1
```

Expected:

```text
ok  	monorepo-template/apps/api/internal/middleware
```

## Task 6: GraphQL Admin Auth Schema, Resolvers, Cookie Bridge, And Guards

**Files:**

- Create: `libs/graphql/schema/admin_auth.graphql`
- Modify: `apps/api/internal/graph/resolver.go`
- Modify: `apps/api/internal/graph/schema.resolvers.go`
- Create: `apps/api/internal/graph/admin_auth_resolvers_test.go`
- Modify: `apps/api/internal/graph/schema_resolvers_test.go`
- Generate: `apps/api/internal/graph/generated.go`
- Generate: `apps/api/internal/graph/model/models_gen.go`

- [ ] **Step 1: Add admin auth GraphQL schema**

Create `libs/graphql/schema/admin_auth.graphql`:

```graphql
# FILE: libs/graphql/schema/admin_auth.graphql
# VERSION: 1.0.0
# START_MODULE_CONTRACT
#   PURPOSE: Define admin auth GraphQL types and operations for web-admin authentication.
#   SCOPE: Admin user object, login/logout/me/createAdmin contracts, and auth result unions; excludes resolver implementation and public REST users API.
#   DEPENDS: libs/graphql/schema/common.graphql.
#   LINKS: M-GRAPHQL-SCHEMA / V-M-GRAPHQL-SCHEMA / M-API / V-M-API / M-WEB-ADMIN / V-M-WEB-ADMIN.
#   ROLE: CONFIG
#   MAP_MODE: SUMMARY
# END_MODULE_CONTRACT
# START_MODULE_MAP
#   AdminUser - Public admin identity returned by auth operations.
#   LoginAdminResult - Login success, validation, or auth error.
#   CreateAdminResult - Admin creation success, validation, or auth error.
#   logoutAdmin - Clears backend session state and browser cookie.
# END_MODULE_MAP
# START_CHANGE_SUMMARY
#   LAST_CHANGE: 1.0.0 - Added web-admin backend auth schema.
# END_CHANGE_SUMMARY

type AdminUser {
  id: UUID!
  email: String!
  name: String!
  role: String!
  createdAt: DateTime!
  updatedAt: DateTime!
}

input LoginAdminInput {
  email: String!
  password: String!
}

input CreateAdminInput {
  email: String!
  name: String!
  password: String!
}

type LoginAdminSuccess {
  admin: AdminUser!
}

type CreateAdminSuccess {
  admin: AdminUser!
}

type LogoutAdminSuccess {
  ok: Boolean!
}

union LoginAdminResult = LoginAdminSuccess | ValidationError | AuthError
union CreateAdminResult = CreateAdminSuccess | ValidationError | AuthError
union LogoutAdminResult = LogoutAdminSuccess

extend type Query {
  me: AdminUser
}

extend type Mutation {
  loginAdmin(input: LoginAdminInput!): LoginAdminResult!
  logoutAdmin: LogoutAdminResult!
  createAdmin(input: CreateAdminInput!): CreateAdminResult!
}
```

- [ ] **Step 2: Run API codegen**

Run:

```bash
bunx nx run api:codegen
```

Expected:

```text
Successfully ran target codegen for project api
```

If gqlgen reports missing resolver methods, continue to Step 3.

- [ ] **Step 3: Extend resolver dependencies**

Modify `apps/api/internal/graph/resolver.go`:

```go
// FILE: apps/api/internal/graph/resolver.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Hold GraphQL resolver dependencies.
//   SCOPE: User and admin-auth service dependencies for generated gqlgen resolvers; excludes resolver method behavior.
//   DEPENDS: apps/api/internal/service.
//   LINKS: M-API / V-M-API / M-GRAPHQL-SCHEMA / V-M-GRAPHQL-SCHEMA.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Resolver - Dependency container for GraphQL resolver implementations.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added AdminAuthService dependency for backend admin auth.
// END_CHANGE_SUMMARY

package graph

import "monorepo-template/apps/api/internal/service"

// Resolver holds dependencies for GraphQL resolvers.
type Resolver struct {
	UserService      *service.UserService
	AdminAuthService *service.AdminAuthService
}
```

- [ ] **Step 4: Write failing admin resolver tests**

Create `apps/api/internal/graph/admin_auth_resolvers_test.go` with service fake:

```go
package graph

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"monorepo-template/apps/api/internal/graph/model"
	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/apps/api/internal/service"
	"monorepo-template/libs/go/logger"
)

func TestMe_ReturnsNilWithoutSession(t *testing.T) {
	resolver := &Resolver{}
	admin, err := resolver.Query().Me(context.Background())
	require.NoError(t, err)
	assert.Nil(t, admin)
}

func TestMe_ReturnsAdminFromPrincipal(t *testing.T) {
	resolver := &Resolver{}
	ctx := middleware.ContextWithAdminPrincipal(context.Background(), middleware.AdminPrincipal{
		ID: "admin-1", Email: "admin@example.com", Name: "Admin", Role: service.AdminRoleAdmin,
		CreatedAt: "2026-06-07T00:00:00Z", UpdatedAt: "2026-06-07T00:00:00Z",
	})

	admin, err := resolver.Query().Me(ctx)

	require.NoError(t, err)
	require.NotNil(t, admin)
	assert.Equal(t, "admin@example.com", admin.Email)
}

func TestLoginAdmin_SetsSessionCookie(t *testing.T) {
	repo := newFakeAdminRepoForGraph()
	hash, err := bcrypt.GenerateFromPassword([]byte("StrongPassword123!"), bcrypt.MinCost)
	require.NoError(t, err)
	repo.adminsByEmail["admin@example.com"] = &service.Admin{
		ID: "admin-1", Email: "admin@example.com", Name: "Admin", PasswordHash: string(hash), Role: service.AdminRoleAdmin, IsActive: true,
		CreatedAt: "2026-06-07T00:00:00Z", UpdatedAt: "2026-06-07T00:00:00Z",
	}
	resolver := &Resolver{AdminAuthService: service.NewAdminAuthService(repo, newFakeAdminSessionsForGraph())}
	rec := httptest.NewRecorder()
	ctx := middleware.ContextWithAdminCookieBridge(context.Background(), middleware.AdminCookieBridge{
		Response: rec,
		Config:   testAdminCookieConfig(),
	})

	result, err := resolver.Mutation().LoginAdmin(ctx, model.LoginAdminInput{Email: "admin@example.com", Password: "StrongPassword123!"})

	require.NoError(t, err)
	_, ok := result.(model.LoginAdminSuccess)
	require.True(t, ok)
	cookie := rec.Result().Cookies()[0]
	assert.Equal(t, "web_admin_session", cookie.Name)
	assert.True(t, cookie.HttpOnly)
	assert.Equal(t, "/graphql", cookie.Path)
}

func TestLogoutAdmin_DeletesSessionAndClearsCookie(t *testing.T) {
	sessions := newFakeAdminSessionsForGraph()
	sessions.sessions["session-1"] = "admin-1"
	resolver := &Resolver{AdminAuthService: service.NewAdminAuthService(newFakeAdminRepoForGraph(), sessions)}
	rec := httptest.NewRecorder()
	ctx := middleware.ContextWithAdminCookieBridge(context.Background(), middleware.AdminCookieBridge{
		Response: rec,
		Config:   testAdminCookieConfig(),
	})
	ctx = middleware.ContextWithAdminSessionID(ctx, "session-1")

	result, err := resolver.Mutation().LogoutAdmin(ctx)

	require.NoError(t, err)
	success, ok := result.(model.LogoutAdminSuccess)
	require.True(t, ok)
	assert.True(t, success.Ok)
	assert.Equal(t, "session-1", sessions.deleted)
	cookie := rec.Result().Cookies()[0]
	assert.Equal(t, "web_admin_session", cookie.Name)
	assert.Equal(t, -1, cookie.MaxAge)
}

func TestLogoutAdmin_SucceedsAndClearsCookieWithoutSession(t *testing.T) {
	sessions := newFakeAdminSessionsForGraph()
	resolver := &Resolver{AdminAuthService: service.NewAdminAuthService(newFakeAdminRepoForGraph(), sessions)}
	rec := httptest.NewRecorder()
	ctx := middleware.ContextWithAdminCookieBridge(context.Background(), middleware.AdminCookieBridge{
		Response: rec,
		Config:   testAdminCookieConfig(),
	})

	result, err := resolver.Mutation().LogoutAdmin(ctx)

	require.NoError(t, err)
	success, ok := result.(model.LogoutAdminSuccess)
	require.True(t, ok)
	assert.True(t, success.Ok)
	assert.Empty(t, sessions.deleted)
	cookie := rec.Result().Cookies()[0]
	assert.Equal(t, -1, cookie.MaxAge)
}

func TestCreateAdmin_ReturnsAuthErrorWithoutPrincipal(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	resolver := newAdminResolver(t)
	ctx := logger.WithContext(context.Background(), zap.New(core))
	result, err := resolver.Mutation().CreateAdmin(ctx, model.CreateAdminInput{
		Email: "new@example.com", Name: "New", Password: "StrongPassword123!",
	})
	require.NoError(t, err)
	authErr, ok := result.(model.AuthError)
	require.True(t, ok)
	assert.Contains(t, authErr.Message, "authentication required")
	joined := logText(logs.All())
	assert.Contains(t, joined, "[AdminAuth][guard][BLOCK_AUTHORIZE_GRAPHQL]")
	assert.NotContains(t, joined, "new@example.com")
	assert.NotContains(t, joined, "StrongPassword123!")
	assert.NotContains(t, joined, "createAdmin")
}

func TestCreateAdmin_ReturnsSuccessWithPrincipal(t *testing.T) {
	resolver := newAdminResolver(t)
	ctx := middleware.ContextWithAdminPrincipal(context.Background(), middleware.AdminPrincipal{ID: "admin-1", Email: "admin@example.com", Name: "Admin", Role: "ADMIN"})
	result, err := resolver.Mutation().CreateAdmin(ctx, model.CreateAdminInput{
		Email: "new@example.com", Name: "New", Password: "StrongPassword123!",
	})
	require.NoError(t, err)
	success, ok := result.(model.CreateAdminSuccess)
	require.True(t, ok)
	assert.Equal(t, "new@example.com", success.Admin.Email)
}

func newAdminResolver(t *testing.T) *Resolver {
	t.Helper()
	repo := newFakeAdminRepoForGraph()
	repo.adminsByID["admin-1"] = &service.Admin{ID: "admin-1", Email: "admin@example.com", Role: service.AdminRoleAdmin, IsActive: true}
	return &Resolver{AdminAuthService: service.NewAdminAuthService(repo, newFakeAdminSessionsForGraph())}
}

func testAdminCookieConfig() middleware.AdminCookieConfig {
	return middleware.AdminCookieConfig{Name: "web_admin_session", Path: "/graphql", MaxAge: 3600, Secure: false, SameSite: http.SameSiteLaxMode}
}

type fakeAdminRepoForGraph struct {
	adminsByID    map[string]*service.Admin
	adminsByEmail map[string]*service.Admin
}

func newFakeAdminRepoForGraph() *fakeAdminRepoForGraph {
	return &fakeAdminRepoForGraph{adminsByID: map[string]*service.Admin{}, adminsByEmail: map[string]*service.Admin{}}
}

func (r *fakeAdminRepoForGraph) Count(ctx context.Context) (int, error) {
	return len(r.adminsByID), nil
}

func (r *fakeAdminRepoForGraph) Create(ctx context.Context, input service.CreateAdminInput) (*service.Admin, error) {
	if _, exists := r.adminsByEmail[input.Email]; exists {
		return nil, service.ErrAdminDuplicateEmail
	}
	admin := &service.Admin{
		ID: "admin-" + input.Email, Email: input.Email, Name: input.Name,
		PasswordHash: input.PasswordHash, Role: input.Role, IsActive: true,
		CreatedAt: "2026-06-07T00:00:00Z", UpdatedAt: "2026-06-07T00:00:00Z",
	}
	r.adminsByID[admin.ID] = admin
	r.adminsByEmail[admin.Email] = admin
	return admin, nil
}

func (r *fakeAdminRepoForGraph) GetByEmail(ctx context.Context, email string) (*service.Admin, error) {
	return r.adminsByEmail[email], nil
}

func (r *fakeAdminRepoForGraph) GetByID(ctx context.Context, id string) (*service.Admin, error) {
	return r.adminsByID[id], nil
}

type fakeAdminSessionsForGraph struct {
	sessions map[string]string
	deleted  string
}

func newFakeAdminSessionsForGraph() *fakeAdminSessionsForGraph {
	return &fakeAdminSessionsForGraph{sessions: map[string]string{}}
}

func (s *fakeAdminSessionsForGraph) Create(ctx context.Context, adminID string) (string, error) {
	s.sessions["session-1"] = adminID
	return "session-1", nil
}

func (s *fakeAdminSessionsForGraph) Get(ctx context.Context, sessionID string) (string, error) {
	return s.sessions[sessionID], nil
}

func (s *fakeAdminSessionsForGraph) Delete(ctx context.Context, sessionID string) error {
	s.deleted = sessionID
	delete(s.sessions, sessionID)
	return nil
}

func logText(entries []observer.LoggedEntry) string {
	var out strings.Builder
	for _, entry := range entries {
		out.WriteString(entry.Message)
		out.WriteString("\n")
		for _, field := range entry.Context {
			out.WriteString(field.Key)
			out.WriteString("=")
			out.WriteString(field.String)
			out.WriteString("\n")
		}
	}
	return out.String()
}
```

- [ ] **Step 5: Write failing protected user resolver tests**

Append to `apps/api/internal/graph/schema_resolvers_test.go`:

```go
func TestUsers_ReturnsAuthErrorWithoutAdminPrincipal(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())
	connection, err := resolver.Query().Users(context.Background(), nil)
	require.Error(t, err)
	assert.Nil(t, connection)
	assert.Contains(t, err.Error(), "admin authentication required")
}

func TestProtectedUserResolvers_ReturnAuthErrorWithoutAdminPrincipal(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())
	ctx := context.Background()

	user, err := resolver.Query().User(ctx, "user-1")
	assert.Nil(t, user)
	assertAuthRequired(t, err)

	connection, err := resolver.Query().Users(ctx, nil)
	assert.Nil(t, connection)
	assertAuthRequired(t, err)

	created, err := resolver.Mutation().CreateUser(ctx, model.CreateUserInput{Email: "new@example.com", Name: "New", Password: "StrongPassword123!"})
	require.NoError(t, err)
	authErr, ok := created.(model.AuthError)
	require.True(t, ok)
	assert.Contains(t, authErr.Message, "admin authentication required")

	name := "Updated"
	updated, err := resolver.Mutation().UpdateUser(ctx, "user-1", model.UpdateUserInput{Name: &name})
	assert.Nil(t, updated)
	assertAuthRequired(t, err)

	deleted, err := resolver.Mutation().DeleteUser(ctx, "user-1")
	assert.False(t, deleted)
	assertAuthRequired(t, err)
}

func TestUsers_AllowsAdminPrincipal(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())
	ctx := middleware.ContextWithAdminPrincipal(context.Background(), middleware.AdminPrincipal{ID: "admin-1", Email: "admin@example.com", Role: "ADMIN"})

	connection, err := resolver.Query().Users(ctx, nil)

	require.NoError(t, err)
	require.NotNil(t, connection)
	assert.Equal(t, 1, connection.TotalCount)
}

func assertAuthRequired(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin authentication required")
}
```

Add import:

```go
"monorepo-template/apps/api/internal/middleware"
```

Add this helper to `apps/api/internal/graph/schema_resolvers_test.go`:

```go
func adminCtx() context.Context {
	return middleware.ContextWithAdminPrincipal(context.Background(), middleware.AdminPrincipal{
		ID: "admin-1", Email: "admin@example.com", Name: "Admin", Role: "ADMIN",
		CreatedAt: "2026-06-07T00:00:00Z", UpdatedAt: "2026-06-07T00:00:00Z",
	})
}
```

Update every existing user GraphQL test that is meant to exercise user-domain behavior to use `adminCtx()` instead of `context.Background()`. This includes success, validation, not-found, repository-error, pagination, direct user lookup, create, update, and delete cases. Keep only the new negative auth tests unauthenticated.

- [ ] **Step 6: Run graph tests and verify they fail for missing resolvers/guards**

Run:

```bash
cd apps/api
go test ./internal/graph -run 'Test(Admin|LoginAdmin|LogoutAdmin|Me|CreateAdmin|ProtectedUser|CreateUser|Users|User|UpdateUser|DeleteUser)' -count=1
```

Expected:

```text
FAIL
# failures mention missing Me/LoginAdmin/LogoutAdmin/CreateAdmin implementations or missing admin principal guards
```

- [ ] **Step 7: Implement resolver mapping and guards**

In `apps/api/internal/graph/schema.resolvers.go`, add helper functions:

```go
func requireAdmin(ctx context.Context) (*service.Admin, error) {
	log := logger.FromContext(ctx)
	principal, ok := middleware.GetAdminPrincipal(ctx)
	if !ok {
		log.Warn("[AdminAuth][guard][BLOCK_AUTHORIZE_GRAPHQL] admin principal missing")
		return nil, fmt.Errorf("admin authentication required")
	}
	log.Debug("[AdminAuth][guard][BLOCK_AUTHORIZE_GRAPHQL] admin principal accepted", zap.String("admin_id", principal.ID))
	return &service.Admin{
		ID: principal.ID, Email: principal.Email, Name: principal.Name, Role: principal.Role, IsActive: true,
		CreatedAt: principal.CreatedAt, UpdatedAt: principal.UpdatedAt,
	}, nil
}

func mapAdmin(admin *service.Admin) *model.AdminUser {
	if admin == nil {
		return nil
	}
	return &model.AdminUser{
		ID:        admin.ID,
		Email:     admin.Email,
		Name:      admin.Name,
		Role:      admin.Role,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}
}
```

Add generated resolver bodies matching gqlgen names:

```go
func (r *queryResolver) Me(ctx context.Context) (*model.AdminUser, error) {
	principal, ok := middleware.GetAdminPrincipal(ctx)
	if !ok {
		return nil, nil
	}
	return mapAdmin(&service.Admin{
		ID: principal.ID, Email: principal.Email, Name: principal.Name, Role: principal.Role, IsActive: true,
		CreatedAt: principal.CreatedAt, UpdatedAt: principal.UpdatedAt,
	}), nil
}

func (r *mutationResolver) LoginAdmin(ctx context.Context, input model.LoginAdminInput) (model.LoginAdminResult, error) {
	result, err := r.AdminAuthService.Login(ctx, service.LoginAdminInput{Email: input.Email, Password: input.Password})
	if err != nil {
		if errors.Is(err, service.ErrAdminAuth) {
			return model.AuthError{Message: "invalid email or password"}, nil
		}
		if errors.Is(err, service.ErrAdminValidation) {
			return model.ValidationError{Field: "email", Message: "invalid login input"}, nil
		}
		return nil, err
	}
	middleware.SetAdminSessionCookieFromContext(ctx, result.SessionID)
	return model.LoginAdminSuccess{Admin: mapAdmin(result.Admin)}, nil
}

func (r *mutationResolver) LogoutAdmin(ctx context.Context) (model.LogoutAdminResult, error) {
	sessionID := middleware.AdminSessionIDFromContext(ctx)
	if err := r.AdminAuthService.Logout(ctx, sessionID); err != nil {
		return nil, err
	}
	middleware.ClearAdminSessionCookieFromContext(ctx)
	return model.LogoutAdminSuccess{Ok: true}, nil
}

func (r *mutationResolver) CreateAdmin(ctx context.Context, input model.CreateAdminInput) (model.CreateAdminResult, error) {
	actor, err := requireAdmin(ctx)
	if err != nil {
		return model.AuthError{Message: "admin authentication required"}, nil
	}
	admin, err := r.AdminAuthService.CreateAdmin(ctx, actor, service.NewAdminInput{
		Email: input.Email, Name: input.Name, Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrAdminAuth) {
			return model.AuthError{Message: "admin authentication required"}, nil
		}
		if errors.Is(err, service.ErrAdminDuplicateEmail) {
			return model.ValidationError{Field: "email", Message: "already exists"}, nil
		}
		if errors.Is(err, service.ErrAdminValidation) {
			return model.ValidationError{Field: "password", Message: "invalid admin input"}, nil
		}
		return nil, err
	}
	return model.CreateAdminSuccess{Admin: mapAdmin(admin)}, nil
}
```

Also guard existing user operations at the start of each resolver. For `CreateUser`, return `model.AuthError{Message: "admin authentication required"}` because `CreateUserResult` already includes `AuthError`:

```go
if _, err := requireAdmin(ctx); err != nil {
	return model.AuthError{Message: "admin authentication required"}, nil
}
```

For `User`, `Users`, `UpdateUser`, and `DeleteUser`, the current schema cannot carry `AuthError` in the return type. Return a stable GraphQL error from `requireAdmin` for those methods and keep tests explicit about this schema limitation:

```go
if _, err := requireAdmin(ctx); err != nil {
	return nil, err
}
```

For `DeleteUser`, return `false, err` on missing admin.

Add imports:

```go
import (
	"errors"
	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/libs/go/logger"
	"go.uber.org/zap"
)
```

- [ ] **Step 8: Add context cookie bridge helper**

Extend `apps/api/internal/middleware/admin_auth.go` with a response bridge:

```go
type adminCookieBridgeKey struct{}

type AdminCookieBridge struct {
	Response http.ResponseWriter
	Config   AdminCookieConfig
}

func ContextWithAdminCookieBridge(ctx context.Context, bridge AdminCookieBridge) context.Context {
	return context.WithValue(ctx, adminCookieBridgeKey{}, bridge)
}

func SetAdminSessionCookieFromContext(ctx context.Context, sessionID string) {
	if bridge, ok := ctx.Value(adminCookieBridgeKey{}).(AdminCookieBridge); ok {
		SetAdminSessionCookie(bridge.Response, bridge.Config, sessionID)
	}
}

func ClearAdminSessionCookieFromContext(ctx context.Context) {
	if bridge, ok := ctx.Value(adminCookieBridgeKey{}).(AdminCookieBridge); ok {
		ClearAdminSessionCookie(bridge.Response, bridge.Config)
	}
}
```

- [ ] **Step 9: Run graph tests**

Run:

```bash
cd apps/api
go test ./internal/graph -run 'Test(Admin|LoginAdmin|LogoutAdmin|Me|CreateAdmin|ProtectedUser|CreateUser|Users|User|UpdateUser|DeleteUser)' -count=1
```

Expected:

```text
ok  	monorepo-template/apps/api/internal/graph
```

## Task 7: Web-admin Transport And Generated Operations

**Files:**

- Create: `apps/web-admin/src/entities/admin-auth/api/adminAuth.graphql`
- Modify: `apps/web-admin/src/shared/api/graphql-client.ts`
- Modify: `apps/web-admin/src/shared/api/graphql-client.test.ts`
- Generate: `apps/web-admin/src/shared/api/generated/types.ts`

- [ ] **Step 1: Write failing web-admin GraphQL client test**

Replace `apps/web-admin/src/shared/api/graphql-client.test.ts` with:

```ts
// FILE: apps/web-admin/src/shared/api/graphql-client.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify web-admin GraphQL client construction.
//   SCOPE: Cookie-session transport options and configured GraphQL endpoint; excludes generated operation behavior and UI flows.
//   DEPENDS: vitest, apps/web-admin/src/shared/api/graphql-client.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   graphql client suite - Proves credentialed GraphQL transport for backend admin auth.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added credentialed cookie transport coverage.
// END_CHANGE_SUMMARY

import { describe, expect, it } from 'vitest';
import { createGraphQLClient, graphqlClientOptions } from './graphql-client';

describe('graphql client', () => {
  it('creates a credentialed GraphQLClient for cookie-backed admin auth', () => {
    const client = createGraphQLClient('http://example.test/graphql');

    expect(client).toBeDefined();
    expect((client as unknown as { url: string }).url).toBe('http://example.test/graphql');
    expect(graphqlClientOptions.credentials).toBe('include');
  });
});
```

- [ ] **Step 2: Run web-admin client test and verify it fails**

Run:

```bash
cd apps/web-admin
bun run test src/shared/api/graphql-client.test.ts
```

Expected:

```text
FAIL
expected undefined to be 'include'
```

- [ ] **Step 3: Update GraphQL client transport**

Replace `apps/web-admin/src/shared/api/graphql-client.ts`:

```ts
// FILE: apps/web-admin/src/shared/api/graphql-client.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Construct the web-admin GraphQL transport.
//   SCOPE: Creates a credentialed graphql-request client for httpOnly admin session cookies; excludes UI auth state and route guards.
//   DEPENDS: graphql-request, apps/web-admin/src/shared/config/index.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   createGraphQLClient - Builds a cookie-session-ready GraphQL client.
//   graphqlClient - Default web-admin GraphQL client instance.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Switched admin transport from bearer helper to credentialed cookie requests.
// END_CHANGE_SUMMARY

import { GraphQLClient } from 'graphql-request';
import { appConfig } from '@shared/config';

export const graphqlClientOptions = {
  credentials: 'include' as const,
  headers: {},
};

export function createGraphQLClient(apiUrl = appConfig.apiUrl) {
  return new GraphQLClient(apiUrl, graphqlClientOptions);
}

export const graphqlClient = createGraphQLClient();
```

- [ ] **Step 4: Create admin auth operation documents**

Create `apps/web-admin/src/entities/admin-auth/api/adminAuth.graphql`:

```graphql
# FILE: apps/web-admin/src/entities/admin-auth/api/adminAuth.graphql
# VERSION: 1.0.0
# START_MODULE_CONTRACT
#   PURPOSE: Define web-admin auth GraphQL documents for generated client types.
#   SCOPE: Login, logout, current admin, and create-admin operation documents; excludes login UI and route guards.
#   DEPENDS: libs/graphql/schema/admin_auth.graphql, tools/codegen/codegen.ts, apps/web-admin/src/shared/api/graphql-client.ts.
#   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA / V-M-GRAPHQL-SCHEMA.
#   ROLE: RUNTIME
#   MAP_MODE: SUMMARY
# END_MODULE_CONTRACT
# START_MODULE_MAP
#   LoginAdmin - Authenticates an admin and receives a server-set session cookie.
#   LogoutAdmin - Revokes the current session and clears the cookie.
#   GetCurrentAdmin - Reads the current admin identity for future route guards.
#   CreateAdmin - Creates another admin through an authenticated session.
# END_MODULE_MAP
# START_CHANGE_SUMMARY
#   LAST_CHANGE: 1.0.0 - Added generated operation documents for backend auth readiness.
# END_CHANGE_SUMMARY

mutation LoginAdmin($input: LoginAdminInput!) {
  loginAdmin(input: $input) {
    __typename
    ... on LoginAdminSuccess {
      admin {
        id
        email
        name
        role
      }
    }
    ... on ValidationError {
      field
      message
    }
    ... on AuthError {
      message
    }
  }
}

mutation LogoutAdmin {
  logoutAdmin {
    ... on LogoutAdminSuccess {
      ok
    }
  }
}

query GetCurrentAdmin {
  me {
    id
    email
    name
    role
  }
}

mutation CreateAdmin($input: CreateAdminInput!) {
  createAdmin(input: $input) {
    __typename
    ... on CreateAdminSuccess {
      admin {
        id
        email
        name
        role
      }
    }
    ... on ValidationError {
      field
      message
    }
    ... on AuthError {
      message
    }
  }
}
```

- [ ] **Step 5: Run web-admin codegen and tests**

Run:

```bash
bunx nx run web-admin:codegen
bunx nx test web-admin
bunx nx run web-admin:typecheck
```

Expected:

```text
Successfully ran target codegen for project web-admin
Successfully ran target test for project web-admin
Successfully ran target typecheck for project web-admin
```

## Task 8: Server Wiring, E2E Auth Setup, And Generated Drift

**Files:**

- Modify: `apps/api/cmd/server/main.go`
- Modify: `apps/web-admin/e2e/playwright.config.ts`
- Modify: `apps/web-admin/e2e/helpers.ts`
- Modify: `apps/web-admin/e2e/users-flow.spec.ts`
- Modify: `tools/coverage/coverage.config.json`

- [ ] **Step 1: Wire API server dependencies and route groups**

Modify `apps/api/cmd/server/main.go` to:

1. call `appconfig.ApplyAdminEnvOverlay(&cfg, os.LookupEnv)` and then `appconfig.ApplyAdminDefaults(&cfg)` after config load;
2. construct `adminRepo := postgres.NewAdminRepo(db.Pool)`;
3. construct `sessionStore := redisRepo.NewAdminSessionStore(rdb.RDB, []byte(cfg.AdminSession.KeySecret), cfg.AdminSession.TTL)`;
4. construct `adminAuthService := service.NewAdminAuthService(adminRepo, sessionStore)`;
5. run seed bootstrap after migrations:

```go
l.Info("[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] seed bootstrap starting")
adminCount, err := adminRepo.Count(context.Background())
if err != nil {
	l.Fatal("failed to count admin users", zap.Error(err))
}
if err := appconfig.ValidateAdminBootstrapEnv(os.LookupEnv, adminCount == 0); err != nil {
	l.Fatal("[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] missing bootstrap env", zap.Error(err))
}
if err := appconfig.ValidateAdminBootstrap(cfg, adminCount == 0); err != nil {
	l.Fatal("[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] invalid bootstrap config", zap.Error(err))
}
seeded, err := adminAuthService.SeedInitialAdmin(context.Background(), service.InitialAdminInput{
	Email:    cfg.Admin.InitialEmail,
	Name:     cfg.Admin.InitialName,
	Password: cfg.Admin.InitialPassword,
})
if err != nil {
	l.Fatal("[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] seed failed", zap.Error(err))
}
if seeded {
	l.Info("[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] seeded initial admin")
} else {
	l.Info("[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] seed skipped")
}
```

6. pass `AdminAuthService` to `graph.Resolver`;
7. split route groups:

```go
publicCORS := middleware.CORS(middleware.CORSConfig{
	AllowedOrigins: cfg.Server.CORSOrigins,
	AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
	AllowedHeaders: []string{"Content-Type"},
})
adminCORS := middleware.CORS(middleware.CORSConfig{
	AllowedOrigins:   cfg.Admin.Origins,
	AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
	AllowedHeaders:   []string{"Content-Type"},
	AllowCredentials: true,
})

r.Group(func(public chi.Router) {
	public.Use(publicCORS)
	public.Get("/healthz", healthHandler.Healthz())
	public.Get("/readyz", healthHandler.Readyz(db, rdb))
	usersHandler := healthHandler.NewUsersHandler(userService, l)
	public.Mount("/api/users", usersHandler.Routes())
})

r.Group(func(admin chi.Router) {
	admin.Use(adminCORS)
	admin.Use(middleware.AdminOriginGuard(cfg.Admin.Origins))
	admin.Use(middleware.AdminSessionMiddleware(adminAuthService, cfg.AdminSession.CookieName))
	admin.Handle("/graphql", middleware.WithAdminCookieBridge(srv, middleware.AdminCookieConfigFromConfig(cfg.AdminSession, cfg.Server.Env)))
	if cfg.Server.Env != "production" {
		admin.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))
	}
})
```

Implement `AdminCookieConfigFromConfig` and `WithAdminCookieBridge` in `apps/api/internal/middleware/admin_auth.go` if not already present.

```go
func AdminCookieConfigFromConfig(cfg appconfig.AdminSessionConfig, env string) AdminCookieConfig {
	env = strings.ToLower(strings.TrimSpace(env))
	secure := cfg.CookieSecure == "true" || env == "production"
	return AdminCookieConfig{
		Name: cfg.CookieName,
		Path: "/graphql",
		MaxAge: int(cfg.TTL.Seconds()),
		Secure: secure,
		SameSite: adminSameSiteMode(cfg.SameSite),
	}
}

func WithAdminCookieBridge(next http.Handler, cfg AdminCookieConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := ContextWithAdminCookieBridge(r.Context(), AdminCookieBridge{Response: w, Config: cfg})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func adminSameSiteMode(value string) http.SameSite {
	switch value {
	case "Strict":
		return http.SameSiteStrictMode
	case "None":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
```

Add `monorepo-template/apps/api/internal/appconfig` to `admin_auth.go` imports for `AdminCookieConfigFromConfig`.

- [ ] **Step 2: Delete the stale bearer placeholder**

After the server rewrite no longer calls `middleware.Auth(cfg.Auth.JWTSecret, l)`, delete `apps/api/internal/middleware/auth.go` and `apps/api/internal/middleware/auth_test.go`. Confirm no code still references the placeholder helpers or stale JWT config:

```bash
rg -n "middleware\\.Auth\\(|GetUserID\\(|UserIDKey|Bearer my-token|JWTSecret" apps/api
```

Expected:

```text
# no output
```

- [ ] **Step 3: Update Docker, Dokploy, and env documentation surfaces**

In `docker/docker-compose.yml`, remove `AUTH_JWT_SECRET` from the API auth path and pass admin bootstrap/session env through:

```yaml
- ADMIN_INITIAL_EMAIL=${ADMIN_INITIAL_EMAIL:-}
- ADMIN_INITIAL_PASSWORD=${ADMIN_INITIAL_PASSWORD:-}
- ADMIN_INITIAL_NAME=${ADMIN_INITIAL_NAME:-}
- ADMIN_ORIGINS=${ADMIN_ORIGINS:-http://localhost:3100,http://127.0.0.1:3100}
- ADMIN_SESSION_COOKIE_NAME=${ADMIN_SESSION_COOKIE_NAME:-web_admin_session}
- ADMIN_SESSION_TTL=${ADMIN_SESSION_TTL:-168h}
- ADMIN_SESSION_COOKIE_SECURE=${ADMIN_SESSION_COOKIE_SECURE:-auto}
- ADMIN_SESSION_SAME_SITE=${ADMIN_SESSION_SAME_SITE:-Lax}
- ADMIN_SESSION_KEY_SECRET=${ADMIN_SESSION_KEY_SECRET:-dev-session-key-secret}
```

In `deploy/dokploy/docker-compose.template.yml`, replace `AUTH_JWT_SECRET` with required production placeholders:

```yaml
ADMIN_INITIAL_EMAIL: ${ADMIN_INITIAL_EMAIL}
ADMIN_INITIAL_PASSWORD: ${ADMIN_INITIAL_PASSWORD}
ADMIN_INITIAL_NAME: ${ADMIN_INITIAL_NAME}
ADMIN_ORIGINS: ${ADMIN_ORIGINS}
ADMIN_SESSION_COOKIE_NAME: ${ADMIN_SESSION_COOKIE_NAME:-web_admin_session}
ADMIN_SESSION_TTL: ${ADMIN_SESSION_TTL:-168h}
ADMIN_SESSION_COOKIE_SECURE: ${ADMIN_SESSION_COOKIE_SECURE:-auto}
ADMIN_SESSION_SAME_SITE: ${ADMIN_SESSION_SAME_SITE:-Lax}
ADMIN_SESSION_KEY_SECRET: ${ADMIN_SESSION_KEY_SECRET}
```

In `docs/infrastructure/ci-cd.md`, replace `AUTH_JWT_SECRET` in the required runtime values list with the same `ADMIN_*` variables. State that `ADMIN_INITIAL_*` is required only while `admin_users` is empty, and `ADMIN_SESSION_KEY_SECRET` is always required for API startup.

Also add a local first-run note: fresh Docker/API startup with an empty `admin_users` table requires `ADMIN_INITIAL_EMAIL`, `ADMIN_INITIAL_PASSWORD`, `ADMIN_INITIAL_NAME`, and `ADMIN_SESSION_KEY_SECRET`; after the first admin exists, `ADMIN_INITIAL_*` may be unset but `ADMIN_SESSION_KEY_SECRET` remains required.

After these edits, verify no stale JWT placeholder remains in auth runtime docs:

```bash
rg -n "AUTH_JWT_SECRET|jwt_secret|JWTSecret" docker deploy docs/infrastructure apps/api/config .env.example apps/api/internal/appconfig
```

Expected:

```text
# no output
```

- [ ] **Step 4: Add e2e server env**

In `apps/web-admin/e2e/playwright.config.ts`, add API server env values:

```ts
ADMIN_INITIAL_EMAIL: 'e2e-admin@example.test',
ADMIN_INITIAL_PASSWORD: 'StrongPassword123!',
ADMIN_INITIAL_NAME: 'E2E Admin',
ADMIN_ORIGINS: webBaseURL,
ADMIN_SESSION_COOKIE_NAME: 'web_admin_session',
ADMIN_SESSION_TTL: '168h',
ADMIN_SESSION_COOKIE_SECURE: 'false',
ADMIN_SESSION_SAME_SITE: 'Lax',
ADMIN_SESSION_KEY_SECRET: 'e2e-session-key-secret',
SERVER_CORS_ORIGINS: webBaseURL,
```

- [ ] **Step 5: Add e2e login helper**

In `apps/web-admin/e2e/helpers.ts`, add:

```ts
export async function loginAdmin(apiContext: APIRequestContext) {
  const response = await apiContext.post('/graphql', {
    data: {
      query: `mutation LoginAdmin($input: LoginAdminInput!) {
        loginAdmin(input: $input) {
          __typename
          ... on LoginAdminSuccess {
            admin { id email role }
          }
          ... on AuthError { message }
        }
      }`,
      variables: {
        input: {
          email: 'e2e-admin@example.test',
          password: 'StrongPassword123!',
        },
      },
    },
    headers: {
      'Content-Type': 'application/json',
      Origin: process.env.E2E_WEB_URL ?? 'http://localhost:13000',
    },
  });
  expect(response.ok()).toBeTruthy();
  const cookies = response.headers()['set-cookie'];
  expect(cookies).toContain('web_admin_session=');
}

export async function installAdminSessionCookie(browserContext: BrowserContext) {
  const apiContext = await playwrightRequest.newContext({
    baseURL: apiBaseURL,
    extraHTTPHeaders: { Origin: process.env.E2E_WEB_URL ?? 'http://localhost:13000' },
  });
  try {
    await loginAdmin(apiContext);
    const storage = await apiContext.storageState();
    await browserContext.addCookies(storage.cookies);
  } finally {
    await apiContext.dispose();
  }
}

export async function withAuthenticatedGraphQLContext<T>(
  fn: (context: APIRequestContext) => Promise<T>,
) {
  const context = await playwrightRequest.newContext({
    baseURL: apiBaseURL,
    extraHTTPHeaders: { Origin: process.env.E2E_WEB_URL ?? 'http://localhost:13000' },
  });
  try {
    await loginAdmin(context);
    return await fn(context);
  } finally {
    await context.dispose();
  }
}
```

Add `BrowserContext` to the type imports from `@playwright/test`. Update `createUser` calls in e2e tests to use `withAuthenticatedGraphQLContext`, which is API-bound through `apiBaseURL`.

- [ ] **Step 6: Update users e2e flow to log in before browser CRUD**

In `apps/web-admin/e2e/users-flow.spec.ts`, import `installAdminSessionCookie` and add before each protected browser flow:

```ts
test.beforeEach(async ({ page, context }) => {
  await installAdminSessionCookie(context);
  await page.goto('/');
});
```

Do not use Playwright's built-in `request` fixture for admin login here: its `baseURL` is the web app. The helper must create an API-bound request context, transfer the API cookie into the browser context, and navigate only after the cookie is installed.

- [ ] **Step 7: Add generated coverage allowlist**

In `tools/coverage/coverage.config.json`, add:

```json
{
  "path": "apps/api/internal/repository/postgres/generated/admin_users.sql.go",
  "reason": "sqlc generated admin users query methods",
  "gate": "bunx nx run api:codegen && bunx nx build api && TEST_COMPOSE_PROJECT=mt-admin-sqlc-ref TEST_POSTGRES_CONTAINER_NAME=mt-admin-sqlc-ref-postgres TEST_REDIS_CONTAINER_NAME=mt-admin-sqlc-ref-redis TEST_POSTGRES_VOLUME=mt-admin-sqlc-ref-pg-test-data docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis && cd apps/api && COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable go test ./internal/repository/postgres -run TestAdminRepo -count=1 -v"
}
```

Place it beside the existing sqlc generated entries.

- [ ] **Step 8: Verify empty and non-empty admin startup behavior**

Before final e2e closeout, use the test Postgres/Redis stack to prove both startup modes and record the evidence in `.tasks/web-admin-backend-auth/verification.md`:

1. Start the API with an empty `admin_users` table and real `ADMIN_INITIAL_EMAIL`, `ADMIN_INITIAL_PASSWORD`, `ADMIN_INITIAL_NAME`, and `ADMIN_SESSION_KEY_SECRET`; `/readyz` must become healthy and the seed marker must indicate create.
2. Stop the API without dropping the database.
3. Restart the API with `ADMIN_INITIAL_EMAIL`, `ADMIN_INITIAL_PASSWORD`, and `ADMIN_INITIAL_NAME` unset, but with `ADMIN_SESSION_KEY_SECRET` still set; `/readyz` must become healthy and the seed marker must indicate skip.
4. Drop or truncate only the test database objects created for this smoke check.

The second start is the guard against accidentally requiring seed env forever after the first admin exists.

- [ ] **Step 9: Run e2e and generated drift checks**

Run:

```bash
bunx nx run web-admin:e2e
bunx nx run codegen:validate
```

Expected:

```text
Successfully ran target e2e for project web-admin
Successfully ran target validate for project codegen
```

If e2e fails because cookies are not shared between `request` and `page`, apply the second `test.beforeEach` variant from Step 4 and rerun only:

```bash
bunx nx run web-admin:e2e
```

Expected:

```text
Successfully ran target e2e for project web-admin
```

## Task 9: GRACE Docs And Verification Log Sync

**Files:**

- Modify: `docs/requirements.xml`
- Modify: `docs/technology.xml`
- Modify: `docs/development-plan.xml`
- Modify: `docs/knowledge-graph.xml`
- Modify: `docs/verification-plan.xml`
- Modify only if needed: `docs/operational-packets.xml`
- Modify: `.tasks/web-admin-backend-auth/verification.md`

- [ ] **Step 1: Update GRACE requirements**

In `docs/requirements.xml`:

- add a use case for backend web-admin auth;
- replace the admin GraphQL auth placeholder with a concrete backend auth requirement;
- keep downstream owner/tenant authorization as a future product-specific warning;
- add risks for credentialed CORS/CSRF, session revocation, and generated admin sqlc drift.

Use this wording for the new use case:

```xml
<UC-012>
  <Actor>AdminUser</Actor>
  <Action>Logs in to the web-admin backend, keeps an httpOnly session, reads the current admin, logs out, and creates another admin through protected GraphQL.</Action>
  <Goal>Provide a real backend auth foundation before login UI is added.</Goal>
  <Preconditions>PostgreSQL and Redis are reachable, API migrations have run, and first-admin env config is present only while the admin table is empty.</Preconditions>
  <AcceptanceCriteria>`loginAdmin` sets an httpOnly cookie, `me` identifies the current admin or returns null without a session, `logoutAdmin` revokes the session, `createAdmin` requires an authenticated admin, protected user GraphQL operations reject unauthenticated callers, public REST `/api/users` remains public, and web-admin GraphQL transport sends credentialed cookie requests.</AcceptanceCriteria>
  <Priority>high</Priority>
  <RelatedFlows>DF-WEB-ADMIN-AUTH</RelatedFlows>
</UC-012>
```

- [ ] **Step 2: Update development plan and knowledge graph**

In `docs/development-plan.xml`, add or update module entries so `M-API`, `M-GRAPHQL-SCHEMA`, `M-WEB-ADMIN`, and `M-COVERAGE-GATE` mention admin auth, Redis sessions, admin GraphQL operations, credentialed client transport, and generated admin sqlc coverage gates.

In `docs/knowledge-graph.xml`, add paths:

```xml
<path>apps/api/internal/service/admin_auth.go</path>
<path>apps/api/internal/repository/postgres/admin_repo.go</path>
<path>apps/api/internal/repository/redis/admin_session_store.go</path>
<path>apps/api/internal/middleware/admin_auth.go</path>
<path>apps/api/internal/middleware/admin_origin.go</path>
<path>libs/graphql/schema/admin_auth.graphql</path>
<path>apps/web-admin/src/entities/admin-auth/api/adminAuth.graphql</path>
```

Add cross-links from `M-WEB-ADMIN` to `M-API` for credentialed admin GraphQL sessions and from `M-COVERAGE-GATE` to `M-API` for generated admin sqlc replacement gates.

- [ ] **Step 3: Update verification plan**

In `docs/verification-plan.xml`, update `V-M-API`, `V-M-GRAPHQL-SCHEMA`, `V-M-WEB-ADMIN`, and `V-M-COVERAGE-GATE` with:

- test files created in Tasks 1-8;
- `bunx nx run api:codegen`;
- `bunx nx run graphql:validate`;
- `bunx nx run codegen:validate`;
- `bunx nx test api`;
- `bunx nx build api`;
- `bunx nx run web-admin:codegen`;
- `bunx nx run web-admin:typecheck`;
- `bunx nx test web-admin`;
- `bunx nx run web-admin:e2e`;
- `bun run test:coverage`;
- `bun run verify:coverage`;
- required log markers from the spec;
- negative tests for unauthenticated admin GraphQL;
- allowed and rejected admin origin tests;
- public REST remains public;
- coverage allowlist replacement gate for generated admin sqlc output.

- [ ] **Step 4: Add rollout and rollback procedure**

In `.tasks/web-admin-backend-auth/verification.md`, add a `## Rollout And Rollback` section:

```markdown
## Rollout And Rollback

### Rollout

1. Apply normal API startup migrations with `postgres.RunMigrations`, which runs goose `Up`.
2. Confirm `admin_users` exists and is empty or already contains admins.
3. When `admin_users` is empty, start the API only with real `ADMIN_INITIAL_EMAIL`, `ADMIN_INITIAL_PASSWORD`, `ADMIN_INITIAL_NAME`, and `ADMIN_SESSION_KEY_SECRET` environment values.
4. Confirm `[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN]` logged seed create or skip without password/hash output.

### Rollback

The migration down path drops `admin_users`, which deletes all admin identities. Before running rollback in any environment with real admins, export the table. The backup contains admin emails and password hashes and must be handled as a secret artifact:

    pg_dump "$DATABASE_URL" --table=admin_users --data-only --column-inserts > admin_users.rollback-backup.sql

Rollback command from the API directory:

    cd apps/api
    goose -dir internal/repository/postgres/migrations postgres "$DATABASE_URL" down

Redis sessions use keys with the `admin_session:` prefix and the configured TTL. Before deletion, count matching keys against the configured Redis target:

    redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" -a "$REDIS_PASSWORD" --scan --pattern 'admin_session:*' | wc -l

After rollback, revoke remaining admin sessions explicitly against the same configured Redis target:

    redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" -a "$REDIS_PASSWORD" --scan --pattern 'admin_session:*' | xargs -r redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" -a "$REDIS_PASSWORD" del

If the feature is rolled forward again, restore admins from backup or reseed exactly one first admin through `ADMIN_INITIAL_*` only when the table is empty.

### Post-rollback Validation

1. `admin_users` table is absent or intentionally recreated by a later rollout.
2. `redis-cli --scan --pattern 'admin_session:*'` returns no keys, or all remaining keys are known pre-rollback TTL leftovers scheduled to expire.
3. API startup either succeeds without admin auth wiring in the rolled-back artifact or fails fast with the expected missing-table signal; it must not silently serve a half-wired admin GraphQL auth surface.
```

- [ ] **Step 5: Record verification evidence**

Update `.tasks/web-admin-backend-auth/verification.md` with commands already run and their outcomes. Use exact command text and result snippets. Do not write `PASS` for commands that were not run.

- [ ] **Step 6: Run XML and GRACE checks**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected:

```text
# xmllint prints no output and exits 0
# grace lint exits 0
```

## Task 10: Focused Final Verification And Commit

**Files:**

- All files touched in Tasks 0-9.

- [ ] **Step 1: Run focused final verification**

Run:

```bash
bunx nx run api:codegen
bunx nx run graphql:validate
bunx nx run web-admin:codegen
bunx nx run codegen:validate
bunx nx test api
bunx nx build api
bunx nx test web-admin
bunx nx run web-admin:typecheck
bunx nx run web-admin:e2e
bun run test:coverage
bun run verify:coverage
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
git diff --check
```

Expected:

```text
Successfully ran target codegen for project api
GraphQL schema is valid
Successfully ran target codegen for project web-admin
Successfully ran target validate for project codegen
Successfully ran target test for project api
Successfully ran target build for project api
Successfully ran target test for project web-admin
Successfully ran target typecheck for project web-admin
Successfully ran target e2e for project web-admin
coverage gate exits 0
verify coverage exits 0
# xmllint prints no output
# grace lint exits 0
# git diff --check prints no output
```

- [ ] **Step 2: Update verification log final status**

In `.tasks/web-admin-backend-auth/verification.md`, set every run command row to the actual result and set:

```markdown
## Final Status

READY
```

Run:

```bash
git diff --check -- .tasks/web-admin-backend-auth/verification.md
```

Expected:

```text
# no output
```

- [ ] **Step 3: Commit the verified auth wave**

Run:

```bash
git status --short
git add apps/api apps/web-admin libs/graphql tools/coverage docs .env.example .tasks/web-admin-backend-auth
git commit -m "feat(api): add web-admin backend auth"
```

Expected:

```text
[branch ...] feat(api): add web-admin backend auth
```

If commit hooks run formatting, rerun:

```bash
git status --short
git diff --check
```

Expected:

```text
# no unstaged changes from hook formatting
# git diff --check prints no output
```

## Self-Review

### Spec Coverage

- Separate `admin_users`: Task 2 migration, sqlc queries, repository tests.
- Bootstrap-only first admin from env: Task 1 config, Task 4 service tests, Task 8 server seed wiring.
- No registration and no UI: File Structure Do Not Modify, Task 7 transport-only web-admin work, Out Of Scope preserved.
- Authenticated admin creates later admins: Task 4 service, Task 6 GraphQL `createAdmin`.
- Redis opaque httpOnly sessions: Task 3 Redis store and cookies, Task 6 login/logout bridge.
- HMAC-derived Redis keys: Task 3 store tests and implementation.
- Protected admin GraphQL: Task 6 resolver guard tests and implementation.
- Public REST remains public: Task 5 route group design, Task 8 e2e/helper updates, Task 9 verification docs.
- Credentialed web-admin transport: Task 7 client test and implementation.
- CORS/CSRF admin boundary: Task 5 origin guard and CORS tests.
- E2E auth setup: Task 8 login helper and browser flow updates.
- Generated-code coverage/drift: Task 2 codegen, Task 8 coverage config, Task 10 `codegen:validate`.
- `me` behavior: Task 6 tests.
- Operational markers and secret redaction: Task 9 verification-plan updates and verification log.
- GRACE docs: Task 9.

### Placeholder Scan

No unresolved placeholder markers, deferred-work labels, unchecked placeholder prose, or shortcut task references are intentionally present in this plan.

### Type Consistency

Core names used consistently across tasks:

- `service.Admin`
- `service.CreateAdminInput`
- `service.InitialAdminInput`
- `service.LoginAdminInput`
- `service.NewAdminInput`
- `service.AdminAuthService`
- `service.AdminRepository`
- `service.AdminSessionStore`
- `postgres.NewAdminRepo`
- `redis.NewAdminSessionStore`
- `middleware.AdminPrincipal`
- `middleware.AdminCookieConfig`
- `middleware.AdminSessionMiddleware`
- `middleware.AdminSessionIDFromContext`
- `middleware.AdminOriginGuard`
- `graphqlClientOptions`
- `loginAdmin`, `logoutAdmin`, `me`, `createAdmin`
