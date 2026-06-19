# Test Infrastructure Isolation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `test:coverage`, `test:e2e`, and `verify:coverage` use only isolated test PostgreSQL/Redis infrastructure and fail before destructive tests can touch the dev database.

**Architecture:** Add a dedicated `docker/docker-compose.test.yml` stack, route Playwright and coverage runners through it, and add shared Go testinfra guards that reject unsafe PostgreSQL targets before migrations or cleanup run. Keep normal local development on `docker/docker-compose.dev.yml` and `monorepo_dev`.

**Tech Stack:** Docker Compose, Bun/Nx, Playwright, Go 1.25, pgx, Redis, GRACE XML docs.

---

## File Structure

- Create `docker/docker-compose.test.yml`: isolated PostgreSQL/Redis test services.
- Create `apps/api/internal/testinfra/safe_targets.go`: shared test-only infrastructure defaults and unsafe-target guards used by integration tests.
- Create `apps/api/internal/testinfra/safe_targets_test.go`: focused guard tests for unsafe DSN rejection and test defaults.
- Modify `apps/web/e2e/preflight.mjs`: start only the test compose stack.
- Modify `apps/web/e2e/playwright.config.ts`: pass test PostgreSQL/Redis env overrides to the API web server and remove duplicate global setup.
- Delete `apps/web/e2e/global-setup.ts`: no second preflight path remains.
- Modify `tools/coverage/preflight.mjs`: validate the test compose file exists.
- Modify `tools/coverage/run.mjs`: start test infrastructure and pass safe test env to Go coverage.
- Modify `apps/api/internal/repository/postgres/user_repo_test.go`: require a safe test DSN before migrations and `TRUNCATE`.
- Modify `apps/api/internal/repository/postgres/postgres_test.go`: use test PostgreSQL config instead of dev defaults.
- Modify `apps/api/internal/repository/postgres/user_repo_unit_test.go`: remove lingering dev DSN/config literals from unit-only error paths.
- Modify `apps/api/internal/repository/redis/cache_test.go`: use test Redis config and gate-strict failure behavior.
- Modify `docs/requirements.xml`, `docs/technology.xml`, `docs/development-plan.xml`, `docs/knowledge-graph.xml`, and `docs/verification-plan.xml`: sync GRACE contracts.
- Modify `README.md`: document dev compose vs test compose responsibilities.
- Create `.tasks/test-infrastructure-isolation-verification.md`: record final verification evidence for the risky test-infrastructure change.

Current worktree note: `AGENTS.md` is an unrelated existing modification. Do not stage or revert it while implementing this plan.

---

### Task 1: Add Dedicated Test Compose And E2E Infra Routing

**Files:**

- Create: `docker/docker-compose.test.yml`
- Modify: `apps/web/e2e/preflight.mjs`
- Modify: `apps/web/e2e/playwright.config.ts`
- Delete: `apps/web/e2e/global-setup.ts`

- [ ] **Step 1: Add the isolated test compose file**

Create `docker/docker-compose.test.yml` with this complete content:

```yaml
name: mt-test

services:
  postgres:
    container_name: mt-test-postgres
    image: postgres:16-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${TEST_POSTGRES_USER:-app}
      POSTGRES_PASSWORD: ${TEST_POSTGRES_PASSWORD:-secret}
      POSTGRES_DB: ${TEST_POSTGRES_DB:-monorepo_test}
    ports:
      - '${TEST_POSTGRES_PORT:-17501}:5432'
    volumes:
      - pg-test-data:/var/lib/postgresql/data
    healthcheck:
      test:
        [
          'CMD-SHELL',
          'pg_isready -U ${TEST_POSTGRES_USER:-app} -d ${TEST_POSTGRES_DB:-monorepo_test}',
        ]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    container_name: mt-test-redis
    image: redis:7-alpine
    restart: unless-stopped
    ports:
      - '${TEST_REDIS_PORT:-17502}:6379'
    healthcheck:
      test: ['CMD', 'redis-cli', 'ping']
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  pg-test-data:
```

- [ ] **Step 2: Verify the compose file resolves to test-only names**

Run:

```bash
docker compose -f docker/docker-compose.test.yml config
```

Expected: PASS, and output contains `mt-test-postgres`, `mt-test-redis`, `monorepo_test`, `17501`, and `17502`.

- [ ] **Step 3: Replace e2e preflight with test compose startup**

Replace `apps/web/e2e/preflight.mjs` with this complete content:

```js
import { execFileSync } from 'node:child_process';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const e2eDir = dirname(fileURLToPath(import.meta.url));
const repoRoot = resolve(e2eDir, '../../..');
const composeFile = resolve(repoRoot, 'docker/docker-compose.test.yml');
const testPostgresPort = process.env.TEST_POSTGRES_PORT ?? '17501';
const testPostgresDB = process.env.TEST_POSTGRES_DB ?? 'monorepo_test';
const testRedisPort = process.env.TEST_REDIS_PORT ?? '17502';

function run(command, args) {
  console.log(`[e2e:preflight] ${command} ${args.join(' ')}`);
  execFileSync(command, args, { cwd: repoRoot, stdio: 'inherit' });
}

function assertSafeTestTarget() {
  if (testPostgresDB !== 'monorepo_test') {
    throw new Error(`[e2e:preflight] unsafe test database target: ${testPostgresDB}`);
  }
  if (testPostgresPort === '7501') {
    throw new Error('[e2e:preflight] unsafe test postgres port: 7501 is the dev port');
  }
}

assertSafeTestTarget();
run('docker', ['compose', '-f', composeFile, 'up', '-d', '--wait', 'postgres', 'redis']);
run('docker', ['compose', '-f', composeFile, 'ps', 'postgres', 'redis']);
console.log(
  `[e2e:preflight] test docker services ready: postgres:${testPostgresPort} redis:${testRedisPort}`,
);
```

- [ ] **Step 4: Replace Playwright config with test env overrides**

Replace `apps/web/e2e/playwright.config.ts` with this complete content:

```ts
import { defineConfig, devices } from '@playwright/test';
import { resolve } from 'node:path';

const repoRoot = resolve(__dirname, '../../..');
const webRoot = resolve(__dirname, '..');
const apiPort = process.env.E2E_API_PORT ?? '18080';
const apiBaseURL = process.env.E2E_API_URL ?? `http://localhost:${apiPort}`;
const webPort = process.env.E2E_WEB_PORT ?? '13000';
const webBaseURL = process.env.E2E_WEB_URL ?? `http://localhost:${webPort}`;
const testPostgresHost = process.env.TEST_POSTGRES_HOST ?? 'localhost';
const testPostgresPort = process.env.TEST_POSTGRES_PORT ?? '17501';
const testPostgresUser = process.env.TEST_POSTGRES_USER ?? 'app';
const testPostgresPassword = process.env.TEST_POSTGRES_PASSWORD ?? 'secret';
const testPostgresDB = process.env.TEST_POSTGRES_DB ?? 'monorepo_test';
const testRedisHost = process.env.TEST_REDIS_HOST ?? 'localhost';
const testRedisPort = process.env.TEST_REDIS_PORT ?? '17502';
const testRedisPassword = process.env.TEST_REDIS_PASSWORD ?? '';

export default defineConfig({
  testDir: __dirname,
  outputDir: resolve(repoRoot, 'dist/test-results/web-e2e'),
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['list'],
    ['html', { outputFolder: resolve(repoRoot, 'dist/playwright-report/web'), open: 'never' }],
  ],
  use: {
    baseURL: webBaseURL,
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: [
    {
      command: 'bun run e2e:preflight && cd ../api && go run ./cmd/server',
      url: `${apiBaseURL}/readyz`,
      reuseExistingServer: false,
      timeout: 120_000,
      cwd: webRoot,
      env: {
        SERVER_PORT: apiPort,
        SERVER_CORS_ORIGINS: webBaseURL,
        POSTGRES_HOST: testPostgresHost,
        POSTGRES_PORT: testPostgresPort,
        POSTGRES_USER: testPostgresUser,
        POSTGRES_PASSWORD: testPostgresPassword,
        POSTGRES_DB: testPostgresDB,
        POSTGRES_SSLMODE: 'disable',
        REDIS_HOST: testRedisHost,
        REDIS_PORT: testRedisPort,
        REDIS_PASSWORD: testRedisPassword,
      },
    },
    {
      command: `bun run dev -- --hostname 127.0.0.1 --port ${webPort}`,
      url: webBaseURL,
      reuseExistingServer: false,
      timeout: 120_000,
      cwd: webRoot,
      env: {
        NEXT_PUBLIC_API_URL: `${apiBaseURL}/graphql`,
      },
    },
  ],
});
```

- [ ] **Step 5: Delete the duplicate global setup file**

Remove `apps/web/e2e/global-setup.ts`.

- [ ] **Step 6: Verify e2e preflight targets the test stack**

Run:

```bash
cd apps/web && bun run e2e:preflight
```

Expected: PASS, with logs ending in `test docker services ready: postgres:17501 redis:17502`.

- [ ] **Step 7: Commit Task 1**

```bash
git add docker/docker-compose.test.yml apps/web/e2e/preflight.mjs apps/web/e2e/playwright.config.ts apps/web/e2e/global-setup.ts
git commit -m "test: add isolated e2e infrastructure"
```

---

### Task 2: Add Go Testinfra Safety Guards

**Files:**

- Create: `apps/api/internal/testinfra/safe_targets.go`
- Create: `apps/api/internal/testinfra/safe_targets_test.go`

- [ ] **Step 1: Write focused guard tests first**

Create `apps/api/internal/testinfra/safe_targets_test.go` with this complete content:

```go
package testinfra

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateSafePostgresDSNAcceptsDefaultTestTarget(t *testing.T) {
	t.Setenv("TEST_POSTGRES_DB", "monorepo_test")
	t.Setenv("TEST_POSTGRES_PORT", "17501")

	err := ValidateSafePostgresDSN("postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable")

	require.NoError(t, err)
}

func TestValidateSafePostgresDSNRejectsUnsafeTargets(t *testing.T) {
	t.Setenv("TEST_POSTGRES_DB", "monorepo_test")
	t.Setenv("TEST_POSTGRES_PORT", "17501")
	cases := map[string]string{
		"empty":       "",
		"malformed":   "://bad",
		"wrongScheme": "mysql://app:secret@localhost:17501/monorepo_test",
		"devDB":       "postgres://app:secret@localhost:17501/monorepo_dev?sslmode=disable",
		"devPort":     "postgres://app:secret@localhost:7501/monorepo_test?sslmode=disable",
		"wrongPort":   "postgres://app:secret@localhost:15432/monorepo_test?sslmode=disable",
		"missingPort": "postgres://app:secret@localhost/monorepo_test?sslmode=disable",
	}

	for name, dsn := range cases {
		t.Run(name, func(t *testing.T) {
			err := ValidateSafePostgresDSN(dsn)

			require.Error(t, err)
		})
	}
}

func TestPostgresDSNDefaultsToTestDatabase(t *testing.T) {
	t.Setenv("API_TEST_DATABASE_DSN", "")
	t.Setenv("TEST_POSTGRES_HOST", "127.0.0.1")
	t.Setenv("TEST_POSTGRES_PORT", "17501")
	t.Setenv("TEST_POSTGRES_USER", "app")
	t.Setenv("TEST_POSTGRES_PASSWORD", "secret")
	t.Setenv("TEST_POSTGRES_DB", "monorepo_test")

	dsn := PostgresDSN()

	assert.Equal(t, "postgres://app:secret@127.0.0.1:17501/monorepo_test?sslmode=disable", dsn)
}

func TestPostgresDSNUsesEnvOverride(t *testing.T) {
	t.Setenv("API_TEST_DATABASE_DSN", "postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable")

	dsn := PostgresDSN()

	assert.Equal(t, "postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable", dsn)
}

func TestPostgresConfigUsesTestDefaults(t *testing.T) {
	t.Setenv("TEST_POSTGRES_HOST", "localhost")
	t.Setenv("TEST_POSTGRES_PORT", "17501")
	t.Setenv("TEST_POSTGRES_USER", "app")
	t.Setenv("TEST_POSTGRES_PASSWORD", "secret")
	t.Setenv("TEST_POSTGRES_DB", "monorepo_test")

	cfg := PostgresConfig(t)

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 17501, cfg.Port)
	assert.Equal(t, "app", cfg.User)
	assert.Equal(t, "secret", cfg.Password)
	assert.Equal(t, "monorepo_test", cfg.DB)
	assert.Equal(t, "disable", cfg.SSLMode)
}

func TestRedisConfigUsesTestDefaults(t *testing.T) {
	t.Setenv("TEST_REDIS_HOST", "localhost")
	t.Setenv("TEST_REDIS_PORT", "17502")
	t.Setenv("TEST_REDIS_PASSWORD", "")
	t.Setenv("TEST_REDIS_DB", "0")

	cfg := RedisConfig(t)

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 17502, cfg.Port)
	assert.Equal(t, "", cfg.Password)
	assert.Equal(t, 0, cfg.DB)
}
```

- [ ] **Step 2: Run the tests and verify they fail before implementation**

Run:

```bash
cd apps/api && go test ./internal/testinfra
```

Expected: FAIL with undefined symbols such as `ValidateSafePostgresDSN`, `PostgresDSN`, `PostgresConfig`, and `RedisConfig`.

- [ ] **Step 3: Add the testinfra implementation**

Create `apps/api/internal/testinfra/safe_targets.go` with this complete content:

```go
package testinfra

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"monorepo-template/libs/go/config"
)

const (
	TestPostgresDB      = "monorepo_test"
	DefaultPostgresHost = "localhost"
	DefaultPostgresPort = "17501"
	DefaultPostgresUser = "app"
	DefaultPostgresPass = "secret"
	DefaultRedisHost    = "localhost"
	DefaultRedisPort    = "17502"
	DevPostgresPort     = "7501"
)

type TestingT interface {
	Helper()
	Fatalf(format string, args ...any)
}

func CoverageGateEnabled() bool {
	return os.Getenv("COVERAGE_GATE") == "1"
}

func PostgresDSN() string {
	if dsn := os.Getenv("API_TEST_DATABASE_DSN"); strings.TrimSpace(dsn) != "" {
		return dsn
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env("TEST_POSTGRES_USER", DefaultPostgresUser),
		env("TEST_POSTGRES_PASSWORD", DefaultPostgresPass),
		env("TEST_POSTGRES_HOST", DefaultPostgresHost),
		TestPostgresPort(),
		env("TEST_POSTGRES_DB", TestPostgresDB),
	)
}

func PostgresConfig(t TestingT) config.PostgresConfig {
	t.Helper()
	port := mustPort(t, "TEST_POSTGRES_PORT", TestPostgresPort())
	cfg := config.PostgresConfig{
		Host:     env("TEST_POSTGRES_HOST", DefaultPostgresHost),
		Port:     port,
		User:     env("TEST_POSTGRES_USER", DefaultPostgresUser),
		Password: env("TEST_POSTGRES_PASSWORD", DefaultPostgresPass),
		DB:       env("TEST_POSTGRES_DB", TestPostgresDB),
		SSLMode:  "disable",
		MaxConns: 2,
		MinConns: 1,
	}
	RequireSafePostgresDSN(t, cfg.DSN())
	return cfg
}

func RedisConfig(t TestingT) config.RedisConfig {
	t.Helper()
	port := mustPort(t, "TEST_REDIS_PORT", env("TEST_REDIS_PORT", DefaultRedisPort))
	db := mustInt(t, "TEST_REDIS_DB", env("TEST_REDIS_DB", "0"))
	return config.RedisConfig{
		Host:     env("TEST_REDIS_HOST", DefaultRedisHost),
		Port:     port,
		Password: env("TEST_REDIS_PASSWORD", ""),
		DB:       db,
	}
}

func RequireSafePostgresDSN(t TestingT, dsn string) {
	t.Helper()
	if err := ValidateSafePostgresDSN(dsn); err != nil {
		t.Fatalf("unsafe postgres test DSN: %v", err)
	}
}

func ValidateSafePostgresDSN(dsn string) error {
	if strings.TrimSpace(dsn) == "" {
		return fmt.Errorf("dsn is empty")
	}

	parsed, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("parse dsn: %w", err)
	}
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return fmt.Errorf("scheme %q is not postgres", parsed.Scheme)
	}

	db := strings.TrimPrefix(parsed.Path, "/")
	expectedDB := env("TEST_POSTGRES_DB", TestPostgresDB)
	if db != expectedDB {
		return fmt.Errorf("database %q is not %q", db, expectedDB)
	}

	port := parsed.Port()
	if port == "" {
		return fmt.Errorf("dsn must include a postgres port")
	}
	if port == DevPostgresPort {
		return fmt.Errorf("port %s is the development postgres port", DevPostgresPort)
	}
	if port != TestPostgresPort() {
		return fmt.Errorf("port %s is not the configured test postgres port %s", port, TestPostgresPort())
	}

	return nil
}

func TestPostgresPort() string {
	return env("TEST_POSTGRES_PORT", DefaultPostgresPort)
}

func mustPort(t TestingT, key string, value string) int {
	t.Helper()
	port := mustInt(t, key, value)
	if port <= 0 {
		t.Fatalf("%s must be greater than zero, got %d", key, port)
	}
	return port
}

func mustInt(t TestingT, key string, value string) int {
	t.Helper()
	parsed, err := strconv.Atoi(value)
	if err != nil {
		t.Fatalf("%s must be an integer, got %q", key, value)
	}
	return parsed
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
```

- [ ] **Step 4: Run the focused guard tests**

Run:

```bash
cd apps/api && go test ./internal/testinfra
```

Expected: PASS.

- [ ] **Step 5: Commit Task 2**

```bash
git add apps/api/internal/testinfra/safe_targets.go apps/api/internal/testinfra/safe_targets_test.go
git commit -m "test: guard unsafe database targets"
```

---

### Task 3: Route Go Integration Tests Through Testinfra

**Files:**

- Modify: `apps/api/internal/repository/postgres/user_repo_test.go`
- Modify: `apps/api/internal/repository/postgres/postgres_test.go`
- Modify: `apps/api/internal/repository/postgres/user_repo_unit_test.go`
- Modify: `apps/api/internal/repository/redis/cache_test.go`

- [ ] **Step 1: Update user repository integration setup before destructive cleanup**

In `apps/api/internal/repository/postgres/user_repo_test.go`, add this import:

```go
	"monorepo-template/apps/api/internal/testinfra"
```

Replace `testPool` with:

```go
func testPool(t *testing.T) *pgxpool.Pool {
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
	_, err = pool.Exec(context.Background(), `TRUNCATE users RESTART IDENTITY CASCADE`)
	require.NoError(t, err)
	return pool
}
```

- [ ] **Step 2: Update PostgreSQL connection integration tests**

In `apps/api/internal/repository/postgres/postgres_test.go`, remove the `time` import and add:

```go
	"monorepo-template/apps/api/internal/testinfra"
```

Replace the config in `TestNew_ConnectsAndPings` with:

```go
	cfg := testinfra.PostgresConfig(t)
```

Replace the config in `TestNew_ReturnsErrorForBadPort` with:

```go
	cfg := testinfra.PostgresConfig(t)
	cfg.Port = 1
```

Replace the config in `TestNew_ReturnsErrorForMalformedDSN` with:

```go
	cfg := testinfra.PostgresConfig(t)
	cfg.Host = "%"
```

Replace the skip condition with:

```go
	if err != nil && !testinfra.CoverageGateEnabled() {
		t.Skipf("postgres integration database is unavailable: %v", err)
	}
```

Replace the bad migration DSN with:

```go
	err := postgresrepo.RunMigrations("postgres://app:secret@localhost:1/monorepo_test?sslmode=disable", zap.NewNop())
```

- [ ] **Step 3: Remove dev database literals from postgres unit tests**

In `apps/api/internal/repository/postgres/user_repo_unit_test.go`, add:

```go
	"monorepo-template/apps/api/internal/testinfra"
```

Replace the DSN in `TestRunMigrations_ReturnsOpenError` with:

```go
	dsn := testinfra.PostgresDSN()
	testinfra.RequireSafePostgresDSN(t, dsn)
	err := RunMigrations(dsn, nil)
```

Replace the config in `TestNew_ReturnsPoolConstructionError` with:

```go
	db, err := New(testinfra.PostgresConfig(t), zap.NewNop())
```

- [ ] **Step 4: Update Redis integration test**

In `apps/api/internal/repository/redis/cache_test.go`, add:

```go
	"monorepo-template/apps/api/internal/testinfra"
```

Replace the test client setup with:

```go
	client, err := redisrepo.New(testinfra.RedisConfig(t), zap.NewNop())
	if err != nil && !testinfra.CoverageGateEnabled() {
		t.Skipf("redis integration service is unavailable: %v", err)
	}
```

- [ ] **Step 5: Run focused repository tests**

Run:

```bash
cd apps/api && COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:7501/monorepo_dev?sslmode=disable go test ./internal/repository/postgres -run TestUserRepo_CreateGetListUpdateDelete -count=1
```

Expected: FAIL before cleanup with `unsafe postgres test DSN`.

Run:

```bash
docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis
```

Expected: PASS.

Run:

```bash
cd apps/api && COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable TEST_REDIS_PORT=17502 go test ./internal/repository/postgres ./internal/repository/redis ./internal/testinfra
```

Expected: PASS.

- [ ] **Step 6: Commit Task 3**

```bash
git add apps/api/internal/repository/postgres/user_repo_test.go apps/api/internal/repository/postgres/postgres_test.go apps/api/internal/repository/postgres/user_repo_unit_test.go apps/api/internal/repository/redis/cache_test.go
git commit -m "test: route repository integration to test infra"
```

---

### Task 4: Route Coverage Gate Through Test Infrastructure

**Files:**

- Modify: `tools/coverage/preflight.mjs`
- Modify: `tools/coverage/run.mjs`

- [ ] **Step 1: Extend coverage preflight required files**

In `tools/coverage/preflight.mjs`, add the test compose file to `requiredFiles`:

```js
const requiredFiles = [
  'tools/coverage/coverage.config.json',
  'package.json',
  'apps/web/vitest.config.ts',
  'apps/web/e2e/playwright.config.ts',
  'apps/web/e2e/preflight.mjs',
  'docker/docker-compose.test.yml',
  'docs/verification-plan.xml',
];
```

- [ ] **Step 2: Add test infra env helpers to coverage runner**

In `tools/coverage/run.mjs`, add these constants after `const config = ...`:

```js
const testPostgresHost = process.env.TEST_POSTGRES_HOST ?? 'localhost';
const testPostgresPort = process.env.TEST_POSTGRES_PORT ?? '17501';
const testPostgresUser = process.env.TEST_POSTGRES_USER ?? 'app';
const testPostgresPassword = process.env.TEST_POSTGRES_PASSWORD ?? 'secret';
const testPostgresDB = process.env.TEST_POSTGRES_DB ?? 'monorepo_test';
const testRedisHost = process.env.TEST_REDIS_HOST ?? 'localhost';
const testRedisPort = process.env.TEST_REDIS_PORT ?? '17502';
const testRedisPassword = process.env.TEST_REDIS_PASSWORD ?? '';
const testDatabaseDSN =
  process.env.API_TEST_DATABASE_DSN ??
  `postgres://${testPostgresUser}:${testPostgresPassword}@${testPostgresHost}:${testPostgresPort}/${testPostgresDB}?sslmode=disable`;
```

Add this helper near `run`:

```js
function coverageGateEnv() {
  return {
    COVERAGE_GATE: '1',
    API_TEST_DATABASE_DSN: testDatabaseDSN,
    TEST_POSTGRES_HOST: testPostgresHost,
    TEST_POSTGRES_PORT: testPostgresPort,
    TEST_POSTGRES_USER: testPostgresUser,
    TEST_POSTGRES_PASSWORD: testPostgresPassword,
    TEST_POSTGRES_DB: testPostgresDB,
    TEST_REDIS_HOST: testRedisHost,
    TEST_REDIS_PORT: testRedisPort,
    TEST_REDIS_PASSWORD: testRedisPassword,
  };
}
```

Add this helper near `coverageGateEnv`:

```js
function assertSafeTestTarget() {
  if (testPostgresDB !== 'monorepo_test') {
    throw new Error(`[Coverage][run] unsafe test database target: ${testPostgresDB}`);
  }
  if (testPostgresPort === '7501') {
    throw new Error('[Coverage][run] unsafe test postgres port: 7501 is the dev port');
  }
}
```

After `run('node', ['tools/coverage/preflight.mjs']);`, add:

```js
assertSafeTestTarget();
run('docker', [
  'compose',
  '-f',
  'docker/docker-compose.test.yml',
  'up',
  '-d',
  '--wait',
  'postgres',
  'redis',
]);
```

Replace the Go coverage env object:

```js
    env: { COVERAGE_GATE: '1' },
```

with:

```js
    env: coverageGateEnv(),
```

- [ ] **Step 3: Verify unsafe coverage DSN fails**

Run:

```bash
API_TEST_DATABASE_DSN=postgres://app:secret@localhost:7501/monorepo_dev?sslmode=disable bun run test:coverage
```

Expected: FAIL in Go repository tests with `unsafe postgres test DSN`.

- [ ] **Step 4: Verify coverage gate uses test infra**

Run:

```bash
bun run test:coverage
```

Expected: PASS with `[Coverage][gate] all thresholds passed`.

- [ ] **Step 5: Commit Task 4**

```bash
git add tools/coverage/preflight.mjs tools/coverage/run.mjs
git commit -m "test: isolate coverage infrastructure"
```

---

### Task 5: Synchronize GRACE Docs And README

**Files:**

- Modify: `README.md`
- Modify: `docs/requirements.xml`
- Modify: `docs/technology.xml`
- Modify: `docs/development-plan.xml`
- Modify: `docs/knowledge-graph.xml`
- Modify: `docs/verification-plan.xml`

- [ ] **Step 1: Update README gate description**

In `README.md`, replace the e2e sentence under "Coverage and E2E Gate" with:

```markdown
`test:e2e` starts PostgreSQL and Redis from `docker/docker-compose.test.yml`, then starts the API on `18080` and the web app on `13000` unless `E2E_API_PORT`, `E2E_API_URL`, `E2E_WEB_PORT`, or `E2E_WEB_URL` are set. Local development still uses `docker/docker-compose.dev.yml` with `monorepo_dev`; coverage and e2e gates use the isolated `monorepo_test` stack on default ports `17501` and `17502`.
```

- [ ] **Step 2: Update requirements constraints and risks**

In `docs/requirements.xml`, replace constraint 7 with:

```xml
    <constraint-7>Playwright e2e and coverage-gate integration checks use dedicated test infrastructure from `docker/docker-compose.test.yml`, `monorepo_test`, PostgreSQL port `17501`, Redis port `17502`, and isolated API/web ports `18080` and `13000`.</constraint-7>
```

Add this risk after the existing risk list:

```xml
    <risk-6>Destructive repository tests can damage local development data if test gates ever fall back to `monorepo_dev`; gate runners and repository cleanup helpers must fail before cleanup when a DSN is unsafe.</risk-6>
```

- [ ] **Step 3: Update technology e2e policy**

In `docs/technology.xml`, replace the e2e policy with:

```xml
    <e2e-policy>Playwright starts PostgreSQL and Redis through `docker/docker-compose.test.yml`, targets `monorepo_test` on PostgreSQL port `17501` and Redis port `17502`, starts the API on `18080`, starts the web app on `13000`, and verifies health, readiness, GraphQL CRUD, and browser users flows.</e2e-policy>
```

- [ ] **Step 4: Update development-plan CoverageGate contract**

In `docs/development-plan.xml`, update `M-COVERAGE-GATE` inputs to include:

```xml
          <param name="test-infrastructure" type="docker/docker-compose.test.yml PostgreSQL and Redis services" />
```

Update its errors to include:

```xml
          <error code="UNSAFE_TEST_DATABASE_TARGET" />
```

Update its interface entries to include:

```xml
        <export-test-infrastructure PURPOSE="Start isolated PostgreSQL and Redis services for coverage and e2e gates." />
        <export-safe-target-guards PURPOSE="Fail destructive tests before cleanup when the database target is not `monorepo_test`." />
```

Update its target sources to include:

```xml
        <source>docker/docker-compose.test.yml</source>
        <source>apps/api/internal/testinfra</source>
```

- [ ] **Step 5: Update knowledge graph CoverageGate annotations**

In `docs/knowledge-graph.xml`, add these paths under `M-COVERAGE-GATE`:

```xml
      <path>docker/docker-compose.test.yml</path>
      <path>apps/api/internal/testinfra</path>
```

Add these annotations under `M-COVERAGE-GATE`:

```xml
        <export-testCompose PURPOSE="Start isolated PostgreSQL and Redis for coverage and e2e gates." />
        <export-safeDatabaseTarget PURPOSE="Reject `monorepo_dev`, dev port `7501`, empty DSNs, and malformed DSNs before destructive cleanup." />
```

- [ ] **Step 6: Update verification plan coverage gate**

In `docs/verification-plan.xml`, update `V-M-COVERAGE-GATE` test files to include:

```xml
        <file>docker/docker-compose.test.yml</file>
        <file>apps/api/internal/testinfra/safe_targets.go</file>
        <file>apps/api/internal/testinfra/safe_targets_test.go</file>
```

Remove `apps/web/e2e/global-setup.ts` from the same test-file list.

Add this failure scenario under `V-M-COVERAGE-GATE`:

```xml
        <scenario-4 kind="failure">Coverage and e2e gates fail before destructive cleanup when PostgreSQL target is empty, malformed, `monorepo_dev`, or dev port `7501`.</scenario-4>
```

Add this trace assertion:

```xml
        <assertion-4>Repository cleanup that truncates tables must run only after a safe `monorepo_test` DSN guard passes.</assertion-4>
```

- [ ] **Step 7: Validate XML and GRACE docs**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
```

Expected: PASS with no output.

Run:

```bash
grace lint --path .
```

Expected: PASS with `Issues: 0`.

- [ ] **Step 8: Commit Task 5**

```bash
git add README.md docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml
git commit -m "docs: document isolated test infrastructure"
```

---

### Task 6: Final Gate Verification

**Files:**

- Verify all files changed by Tasks 1 through 5.
- Create: `.tasks/test-infrastructure-isolation-verification.md`

- [ ] **Step 1: Confirm dev references remain only in dev docs/config or historical plans**

Run:

```bash
rg -n "monorepo_dev|7501|7502|docker-compose.dev.yml" apps/api/config apps/api/internal apps/web/e2e tools/coverage docker README.md docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml
```

Expected: output may include dev-only README/local development references and `apps/api/config/config.yml`, but must not show gate runners or repository integration tests using `monorepo_dev`, port `7501`, or `docker-compose.dev.yml`.

- [ ] **Step 2: Run focused coverage gate**

Run:

```bash
bun run test:coverage
```

Expected: PASS with `[Coverage][gate] all thresholds passed`.

- [ ] **Step 3: Run focused e2e gate**

Run:

```bash
bun run test:e2e
```

Expected: PASS, with Playwright artifacts under `dist/test-results/web-e2e` and `dist/playwright-report/web`.

- [ ] **Step 4: Run final handoff gate**

Run:

```bash
bun run verify:coverage
```

Expected: PASS.

- [ ] **Step 5: Record final verification evidence**

Create `.tasks/test-infrastructure-isolation-verification.md` with this content after all commands above pass:

```markdown
# Test Infrastructure Isolation Verification

Date: 2026-05-02

## Scope

The coverage and e2e gates now use isolated PostgreSQL and Redis test infrastructure. Destructive repository cleanup is guarded before it can run against a development database.

## Commands

- `docker compose -f docker/docker-compose.test.yml config`: PASS
- `cd apps/web && bun run e2e:preflight`: PASS
- `cd apps/api && COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:7501/monorepo_dev?sslmode=disable go test ./internal/repository/postgres -run TestUserRepo_CreateGetListUpdateDelete -count=1`: FAIL as expected with `unsafe postgres test DSN`
- `cd apps/api && COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable TEST_REDIS_PORT=17502 go test ./internal/repository/postgres ./internal/repository/redis ./internal/testinfra`: PASS
- `bun run test:coverage`: PASS
- `bun run test:e2e`: PASS
- `bun run verify:coverage`: PASS
- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`: PASS
- `grace lint --path .`: PASS

## Evidence Summary

Gate runners and repository integration tests target `monorepo_test` on the test PostgreSQL port. The explicit unsafe-DSN regression check fails before cleanup when pointed at `monorepo_dev` on the development port.
```

- [ ] **Step 6: Commit final verification evidence**

```bash
git add .tasks
git commit -m "docs: record test infrastructure verification"
```

Expected: commit includes `.tasks/test-infrastructure-isolation-verification.md`.

---

## Review Checklist For This Plan

- Spec coverage: Tasks 1 through 5 cover dedicated compose, e2e routing, coverage routing, unsafe target guards, repository cleanup safety, Redis test routing, README, and GRACE XML synchronization.
- Test-first coverage: Task 2 starts with guard tests that fail before implementation; Task 3 includes an explicit unsafe DSN failure check before safe test infra verification.
- Safety invariant: no task allows `COVERAGE_GATE=1` to skip missing test infrastructure or run destructive cleanup against `monorepo_dev`.
- Verification: Task 6 runs focused gates and the final `verify:coverage` handoff gate.
