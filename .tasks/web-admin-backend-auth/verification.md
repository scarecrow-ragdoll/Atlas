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
<!--   LAST_CHANGE: 1.0.4 - Recorded local dev startup defaults and stale goose history repair evidence. -->
<!-- END_CHANGE_SUMMARY -->

# Web-admin Backend Auth Verification

## Command Evidence

| Command                                                                                                                                                                                                                                                                                                                                          | Result                 | Notes                                                                                                                                       |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------- | ----------- | ------------- | ---------- | ----- | ---- | ---------- | ---------------------- | ---- | ----------------------------------------------- |
| `cd apps/api && go test ./internal/appconfig -run 'TestConfig\_(Admin                                                                                                                                                                                                                                                                            | Pagination)' -count=1` | PASS                                                                                                                                        | `ok monorepo-template/apps/api/internal/appconfig`.  |
| `bunx vitest run --config tools/vitest.config.ts tools/workspace/dev-config.test.ts`                                                                                                                                                                                                                                                             | PASS                   | Local `api:serve` now provides dev-only admin seed/session defaults for `bun run dev`.                                                      |
| `cd apps/api && go test ./internal/repository/postgres -run 'TestAdminUsersMigrationVersionAfterHistoricalLocalVersions\|TestRunMigrations_ReturnsErrorForBadDSN\|TestRunMigrations_ReturnsOpenError' -count=1`                                                                                                                                  | PASS                   | `admin_users` migration is numbered after observed local goose version `78`; migration error tests still pass.                              |
| `cd apps/api && go test ./internal/repository/postgres -run TestAdminRepo -count=1`                                                                                                                                                                                                                                                              | PASS                   | `ok monorepo-template/apps/api/internal/repository/postgres`.                                                                               |
| `cd apps/api && TEST_REDIS_PORT=19502 go test ./internal/repository/redis -run TestAdminSessionStore -count=1`                                                                                                                                                                                                                                   | PASS                   | `ok monorepo-template/apps/api/internal/repository/redis`.                                                                                  |
| `cd apps/api && go test ./internal/service -run TestAdminAuth -count=1`                                                                                                                                                                                                                                                                          | PASS                   | `ok monorepo-template/apps/api/internal/service`.                                                                                           |
| `cd apps/api && go test ./internal/middleware -run 'TestCORS                                                                                                                                                                                                                                                                                     | TestAdmin' -count=1`   | PASS                                                                                                                                        | `ok monorepo-template/apps/api/internal/middleware`. |
| `bunx nx run graphql:validate`                                                                                                                                                                                                                                                                                                                   | PASS                   | `GraphQL schema is valid`.                                                                                                                  |
| `bunx nx run api:codegen`                                                                                                                                                                                                                                                                                                                        | PASS                   | `Successfully ran target codegen for project api`.                                                                                          |
| `cd apps/api && go test ./internal/graph -run 'Test(Admin                                                                                                                                                                                                                                                                                        | LoginAdmin             | LogoutAdmin                                                                                                                                 | Me                                                   | CreateAdmin | ProtectedUser | CreateUser | Users | User | UpdateUser | DeleteUser)' -count=1` | PASS | `ok monorepo-template/apps/api/internal/graph`. |
| `bunx nx run web-admin:codegen`                                                                                                                                                                                                                                                                                                                  | PASS                   | `Successfully ran target codegen for project web-admin`.                                                                                    |
| `bunx nx run web-admin:typecheck`                                                                                                                                                                                                                                                                                                                | PASS                   | `Successfully ran target typecheck for project web-admin`.                                                                                  |
| `bunx nx test web-admin`                                                                                                                                                                                                                                                                                                                         | PASS                   | `Successfully ran target test for project web-admin`.                                                                                       |
| `cd apps/web-admin && bun test --run src/shared/api/graphql-client.test.ts`                                                                                                                                                                                                                                                                      | PASS                   | Cookie credentials and no bearer helper verified.                                                                                           |
| `TEST_RESOURCE_PREFIX=mt-admin-e2e TEST_POSTGRES_CONTAINER_NAME=mt-admin-e2e-postgres TEST_REDIS_CONTAINER_NAME=mt-admin-e2e-redis TEST_POSTGRES_VOLUME=mt-admin-e2e-pg-test-data TEST_POSTGRES_PORT=19701 TEST_REDIS_PORT=19702 E2E_API_PORT=19780 E2E_WEB_PORT=19730 bunx nx run web-admin:e2e`                                                | PASS                   | `6 passed`; fresh DB seeded e2e admin and protected UI flows ran after login cookie.                                                        |
| `bunx nx run codegen:validate`                                                                                                                                                                                                                                                                                                                   | PASS                   | Passed after API GraphQL/sqlc and web-admin generated artifacts were synchronized.                                                          |
| `bunx nx test api`                                                                                                                                                                                                                                                                                                                               | PASS                   | `Successfully ran target test for project api`.                                                                                             |
| `bunx nx build api`                                                                                                                                                                                                                                                                                                                              | PASS                   | `Successfully ran target build for project api`; also included in final `verify:coverage`.                                                  |
| `TEST_RESOURCE_PREFIX=mt-admin-coverage TEST_COMPOSE_PROJECT=mt-admin-coverage TEST_POSTGRES_CONTAINER_NAME=mt-admin-coverage-postgres TEST_REDIS_CONTAINER_NAME=mt-admin-coverage-redis TEST_POSTGRES_VOLUME=mt-admin-coverage-pg-test-data TEST_POSTGRES_PORT=19801 TEST_REDIS_PORT=19802 bun run test:coverage`                               | PASS                   | `[Coverage][gate] all thresholds passed`.                                                                                                   |
| `TEST_RESOURCE_PREFIX=mt-admin-verify TEST_COMPOSE_PROJECT=mt-admin-verify TEST_POSTGRES_CONTAINER_NAME=mt-admin-verify-postgres TEST_REDIS_CONTAINER_NAME=mt-admin-verify-redis TEST_POSTGRES_VOLUME=mt-admin-verify-pg-test-data TEST_POSTGRES_PORT=19901 TEST_REDIS_PORT=19902 E2E_API_PORT=19980 E2E_WEB_PORT=19930 bun run verify:coverage` | PASS                   | Full pre-MR gate passed: lint, codegen, typecheck, build, coverage, web/web-admin e2e, XML, and GRACE lint.                                 |
| `TEST_RESOURCE_PREFIX=mt-admin-final TEST_COMPOSE_PROJECT=mt-admin-final TEST_POSTGRES_CONTAINER_NAME=mt-admin-final-postgres TEST_REDIS_CONTAINER_NAME=mt-admin-final-redis TEST_POSTGRES_VOLUME=mt-admin-final-pg-test-data TEST_POSTGRES_PORT=20001 TEST_REDIS_PORT=20002 E2E_API_PORT=20080 E2E_WEB_PORT=20030 bun run verify:coverage`      | PASS                   | Continuation audit reran the full pre-MR gate and passed lint, codegen, typecheck, build, coverage, web/web-admin e2e, XML, and GRACE lint. |
| `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`                                                                                                                                                                           | PASS                   | Included in final `verify:coverage`.                                                                                                        |
| `grace lint --path .`                                                                                                                                                                                                                                                                                                                            | PASS                   | Included in final `verify:coverage`; 0 errors, 16 pre-existing heuristic warnings in skills/worktree scripts.                               |

## Auth Evidence

| Scenario                                                            | Result | Evidence                                                                                                                                                                         |
| ------------------------------------------------------------------- | ------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Empty table seed creates first admin                                | PASS   | Startup smoke on ports `19601/19602/19680` logged `[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] seeded initial admin`, and `/readyz` returned `{"status":"ok"}`.                     |
| Stale local dev DB applies admin migration after goose version `78` | PASS   | `bunx nx serve api` applied `00079_admin_users.sql`, migrated dev DB to version `79`, seeded one admin, `/readyz` returned `{"status":"ok"}`, and `admin_users` contained 1 row. |
| Non-empty table seed does not create env admin                      | PASS   | Restarted same DB without `ADMIN_INITIAL_*`; logged `[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] seed skipped`, and `/readyz` returned `{"status":"ok"}`.                           |
| Login sets httpOnly cookie                                          | PASS   | `TestAdminCookie_SetAndClear`, GraphQL login tests, and web-admin e2e login helper verified `web_admin_session` set-cookie behavior.                                             |
| `me` returns admin for valid active session                         | PASS   | `cd apps/api && go test ./internal/graph -run ...` covers `Me` with admin principal.                                                                                             |
| `me` returns null for missing or invalid session                    | PASS   | GraphQL and middleware focused tests cover missing principal and inactive/empty sessions.                                                                                        |
| Logout clears cookie and Redis session                              | PASS   | Service, middleware, and GraphQL focused tests cover session delete and cookie clear.                                                                                            |
| `createAdmin` requires authenticated admin                          | PASS   | GraphQL and service focused tests return auth error without principal/actor.                                                                                                     |
| Public REST `/api/users` remains public                             | PASS   | `web:e2e` passed in final `verify:coverage`; public Next flow created/listed/opened a user via `/api/users` without admin cookie.                                                |
| Admin GraphQL rejects public web origin                             | PASS   | `TestAdminOriginGuard_RejectsDisallowedOrigin` covers rejected origin before protected GraphQL execution.                                                                        |
| E2E user CRUD runs after `loginAdmin` cookie                        | PASS   | Fresh-stack web-admin e2e passed 6 tests after installing admin session cookie into browser context.                                                                             |
| Empty-table API startup requires and uses `ADMIN_INITIAL_*` env     | PASS   | Startup smoke proved empty-table bootstrap with real env values.                                                                                                                 |
| Non-empty-table API startup does not require `ADMIN_INITIAL_*`      | PASS   | Startup smoke proved restart without `ADMIN_INITIAL_*` after seed.                                                                                                               |

## Generated Coverage

| Generated path                                                       | Replacement gate                                          | Result |
| -------------------------------------------------------------------- | --------------------------------------------------------- | ------ |
| `apps/api/internal/repository/postgres/generated/admin_users.sql.go` | API codegen, API build, admin repository integration test | PASS   |

## Secret Redaction

| Secret class     | Evidence command or test                                                                                 | Result                                            |
| ---------------- | -------------------------------------------------------------------------------------------------------- | ------------------------------------------------- | ---- |
| Passwords        | `cd apps/api && go test ./internal/service -run TestAdminAuthService_LogsMarkersWithoutSecrets -count=1` | PASS                                              |
| Password hashes  | `cd apps/api && go test ./internal/service -run TestAdminAuthService_LogsMarkersWithoutSecrets -count=1` | PASS                                              |
| Raw cookies      | `cd apps/api && go test ./internal/middleware -run 'TestAdmin(SessionMiddleware                          | OriginGuard)\_LogsMarkerWithoutSecrets' -count=1` | PASS |
| Raw session ids  | `cd apps/api && go test ./internal/middleware -run 'TestAdmin(SessionMiddleware                          | OriginGuard)\_LogsMarkerWithoutSecrets' -count=1` | PASS |
| Credential input | `cd apps/api && go test ./internal/graph -run TestCreateAdmin_ReturnsAuthErrorWithoutPrincipal -count=1` | PASS                                              |

## Rollout And Rollback

### Rollout

1. Apply normal API startup migrations with `postgres.RunMigrations`, which runs goose `Up`.
2. Confirm `admin_users` exists and is empty or already contains admins.
3. When `admin_users` is empty, start the API only with real `ADMIN_INITIAL_EMAIL`, `ADMIN_INITIAL_PASSWORD`, `ADMIN_INITIAL_NAME`, and `ADMIN_SESSION_KEY_SECRET` environment values.
4. Confirm `[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN]` logged seed create or skip without password/hash output.

### Rollback

The migration down path drops `admin_users`, which deletes all admin identities. Before running rollback in any environment with real admins, export the table. The backup contains admin emails and password hashes and must be handled as a secret artifact:

```bash
pg_dump "$DATABASE_URL" --table=admin_users --data-only --column-inserts > admin_users.rollback-backup.sql
```

Rollback command from the API directory:

```bash
cd apps/api
goose -dir internal/repository/postgres/migrations postgres "$DATABASE_URL" down
```

Redis sessions use keys with the `admin_session:` prefix and the configured TTL. Before deletion, count matching keys against the configured Redis target:

```bash
redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" -a "$REDIS_PASSWORD" --scan --pattern 'admin_session:*' | wc -l
```

After rollback, revoke remaining admin sessions explicitly against the same configured Redis target:

```bash
redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" -a "$REDIS_PASSWORD" --scan --pattern 'admin_session:*' | xargs -r redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" -a "$REDIS_PASSWORD" del
```

If the feature is rolled forward again, restore admins from backup or reseed exactly one first admin through `ADMIN_INITIAL_*` only when the table is empty.

### Post-rollback Validation

1. `admin_users` table is absent or intentionally recreated by a later rollout.
2. `redis-cli --scan --pattern 'admin_session:*'` returns no keys, or all remaining keys are known pre-rollback TTL leftovers scheduled to expire.
3. API startup either succeeds without admin auth wiring in the rolled-back artifact or fails fast with the expected missing-table signal; it must not silently serve a half-wired admin GraphQL auth surface.

## Final Status

READY

Final pre-MR gate passed with isolated test services on Postgres `19901`, Redis `19902`, API `19980`, and web `19930`; the continuation audit reran the same full gate on Postgres `20001`, Redis `20002`, API `20080`, and web `20030`. Temporary Docker scopes `mt-admin-coverage`, `mt-admin-verify`, and `mt-admin-final` were removed with `docker compose down -v` after verification.
