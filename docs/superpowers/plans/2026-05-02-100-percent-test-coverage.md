# 100 Percent Test Coverage Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enforce 100 percent coverage for handwritten production code and add mandatory full-stack e2e scenarios for the monorepo template.

**Architecture:** Add a root coverage enforcement layer under `tools/coverage`, route module coverage artifacts into `dist/coverage`, and make root scripts fail unless handwritten Go/TypeScript/tooling surfaces meet the approved contract. Complete the reference user CRUD slice, cover first-party workspace tooling, expand web unit/component tests, and add Playwright scenarios that prove infra, API, GraphQL, and browser flows work together.

**Tech Stack:** Bun workspaces, Nx 20, Go 1.25 coverage profiles, Vitest V8 coverage, Playwright, Next.js 15, React 19, gqlgen, pgx, Docker Compose, GRACE XML artifacts.

**Spec:** `docs/superpowers/specs/2026-05-02-100-percent-test-coverage-design.md`

---

## Scope Check

This is a broad verification-infrastructure change, but it remains one integrated plan because the deliverable is a single root gate: `bun run verify:coverage`. The tasks below are ordered so each commit leaves the repo in a more verifiable state and the full gate only becomes required after its foundations are present.

## File Structure

| File                                                      | Action        | Responsibility                                                                                                                  |
| --------------------------------------------------------- | ------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| `package.json`                                            | Modify        | Add root `test:coverage`, `test:e2e`, `verify:coverage`, tooling test scripts, and coverage dev dependencies.                   |
| `apps/api/project.json`                                   | Modify        | Write API coverage profile to `dist/coverage/go/api/coverage.out`.                                                              |
| `apps/bot/project.json`                                   | Modify        | Write bot coverage profile to `dist/coverage/go/bot/coverage.out`.                                                              |
| `libs/go/config/project.json`                             | Modify        | Write config library coverage profile to `dist/coverage/go/go-config/coverage.out`.                                             |
| `libs/go/logger/project.json`                             | Modify        | Write logger library coverage profile to `dist/coverage/go/go-logger/coverage.out`.                                             |
| `tools/coverage/coverage.config.json`                     | Create        | Define project coverage commands, profile paths, thresholds, and explicit allowlist entries.                                    |
| `tools/coverage/run.mjs`                                  | Create        | Run coverage commands, parse coverage artifacts, enforce thresholds, and print actionable failures.                             |
| `tools/coverage/preflight.mjs`                            | Create        | Verify expected generated files, coverage config, and required scripts exist before running the gate.                           |
| `tools/vitest.config.ts`                                  | Create        | Vitest config for first-party workspace tooling tests with 100 percent thresholds.                                              |
| `tools/nx-go/project.json`                                | Create        | Add `nx-go:test` target.                                                                                                        |
| `tools/codegen/project.json`                              | Create        | Add `codegen:validate` target.                                                                                                  |
| `tools/nx-go/src/executors/**/*.test.ts`                  | Create        | Cover executor command construction, cwd resolution, option handling, success, and failure paths.                               |
| `tools/codegen/codegen.test.ts`                           | Create        | Validate schema glob, document globs, generated output path, scalar mappings, and plugins.                                      |
| `apps/api/internal/repository/postgres/user_repo.go`      | Modify        | Implement `Update` and `Delete`, return stable not-found and duplicate-email behavior.                                          |
| `apps/api/internal/repository/postgres/user_repo_test.go` | Create        | Integration tests for Create, GetByID, List, Update, Delete, duplicate, not-found, and invalid cursor paths.                    |
| `apps/api/internal/service/user_service_test.go`          | Modify        | Cover create password hashing, repo errors, update, delete, list, get-by-id, and nil/not-found behavior.                        |
| `apps/api/internal/graph/schema.resolvers.go`             | Modify        | Implement `updateUser` and `deleteUser`; map duplicate/not-found errors to schema results.                                      |
| `apps/api/internal/graph/schema_resolvers_test.go`        | Create        | Cover GraphQL resolver success/error paths without generated-file line coverage.                                                |
| `apps/api/internal/repository/postgres/postgres_test.go`  | Create        | Integration coverage for DB constructor success/failure and `Ping`.                                                             |
| `apps/api/internal/repository/redis/cache_test.go`        | Create        | Integration coverage for Redis constructor success/failure, `Ping`, and `Close`.                                                |
| `apps/web/vitest.config.ts`                               | Modify        | Raise thresholds to 100, exclude generated/config-only files explicitly, emit coverage summary to `dist/coverage/web`.          |
| `apps/web/package.json`                                   | Modify        | Add `test:coverage`.                                                                                                            |
| `apps/web/project.json`                                   | Modify        | Add `web:test-coverage` and make `web:e2e` point at the nested Playwright config.                                               |
| `apps/web/app/__tests__/*.test.tsx`                       | Create/modify | Cover home page, users page loading/success/empty/error/create states, providers, user detail page, config, and GraphQL client. |
| `apps/web/e2e/playwright.config.ts`                       | Modify        | Run infra preflight, start API and web servers, reuse local Docker infra, emit traces and HTML reports to `dist/coverage/e2e`.  |
| `apps/web/e2e/preflight.mjs`                              | Create        | Start local Docker infra and verify PostgreSQL/Redis containers are healthy before Playwright starts API/web servers.           |
| `apps/web/e2e/users-flow.spec.ts`                         | Create        | Browser e2e for `/users` create/list/detail and duplicate/error behavior.                                                       |
| `apps/web/e2e/graphql-contract.spec.ts`                   | Create        | HTTP e2e for health, readiness, GraphQL CRUD, and validation paths.                                                             |
| `apps/web/e2e/helpers.ts`                                 | Create        | GraphQL request helper, deterministic user factory, and cleanup helper.                                                         |
| `docs/requirements.xml`                                   | Modify        | Add 100 percent handwritten coverage and e2e acceptance criteria.                                                               |
| `docs/technology.xml`                                     | Modify        | Add coverage tooling and Go/TS semantics.                                                                                       |
| `docs/development-plan.xml`                               | Modify        | Add/extend workspace verification module and implementation order.                                                              |
| `docs/knowledge-graph.xml`                                | Modify        | Add coverage tooling graph annotations.                                                                                         |
| `docs/verification-plan.xml`                              | Modify        | Add coverage contract, e2e matrix, allowlist, commands, and evidence expectations.                                              |
| `docs/operational-packets.xml`                            | Modify        | Add worker packet guidance for coverage closure tasks.                                                                          |
| `README.md`                                               | Modify        | Document `verify:coverage`, e2e prerequisites, and artifact locations.                                                          |

## Task 1: Baseline Coverage Inventory

**Files:**

- Create: `.tasks/coverage-baseline.md`
- Read: `docs/superpowers/specs/2026-05-02-100-percent-test-coverage-design.md`
- Read: `package.json`, `apps/*/project.json`, `libs/go/*/project.json`, `tools/*/package.json`

- [ ] **Step 1: Create the baseline report shell**

Create `.tasks/coverage-baseline.md` with this exact structure:

```markdown
# 100 Percent Coverage Baseline

## Commands

| Command | Exit | Notes |
| ------- | ---: | ----- |

## Current Coverage Artifacts

| Project | Artifact | Status |
| ------- | -------- | ------ |

## Handwritten Source Inventory

| Surface | Files | Current Test Surface |
| ------- | ----: | -------------------- |

## Known Gaps

- API GraphQL `updateUser` and `deleteUser` are public schema fields but are not implemented.
- `apps/web/e2e` has Playwright config but no scenario specs.
- `tools/nx-go` and `tools/codegen` have no test targets.
- Root scripts do not expose `test:coverage`, `test:e2e`, or `verify:coverage`.

## Exclusion Candidates

| File or Glob                                  | Reason                      | Required Gate                                          |
| --------------------------------------------- | --------------------------- | ------------------------------------------------------ |
| `apps/api/internal/graph/generated.go`        | gqlgen generated file       | `bunx nx run api:codegen`, `bunx nx build api`         |
| `apps/api/internal/graph/model/models_gen.go` | gqlgen generated model file | `bunx nx run api:codegen`, `bunx nx build api`         |
| `apps/web/src/shared/api/generated/**`        | GraphQL codegen output      | `bunx nx run web:codegen`, `bunx nx run web:typecheck` |
| `apps/api/cmd/server/main.go`                 | bootstrap entrypoint        | API startup and e2e health/readiness gates             |
| `apps/bot/cmd/bot/main.go`                    | bootstrap entrypoint        | `bunx nx build bot` and bot startup smoke gate         |
```

- [ ] **Step 2: Capture current project graph**

Run:

```bash
bunx nx show projects
bunx nx show project api --json
bunx nx show project bot --json
bunx nx show project web --json
bunx nx show project nx-go --json
bunx nx show project codegen --json
```

Expected:

- Projects include `api`, `bot`, `web`, `go-config`, `go-logger`, `graphql`, `nx-go`, and `codegen`.
- `nx-go` and `codegen` currently have no targets.

Append command results and findings to `.tasks/coverage-baseline.md`.

- [ ] **Step 3: Capture current test and coverage behavior**

Run:

```bash
bun run test
cd apps/web && bun run test
cd apps/api && go test -coverprofile=coverage.out ./...
cd apps/bot && go test -coverprofile=coverage.out ./...
cd libs/go/config && go test -coverprofile=coverage.out ./...
cd libs/go/logger && go test -coverprofile=coverage.out ./...
```

Expected current behavior:

- Existing tests may pass.
- Go projects emit local `coverage.out` files.
- Web test does not enforce 100 percent coverage.

Append exit codes and artifact locations to `.tasks/coverage-baseline.md`.

- [ ] **Step 4: Commit the baseline report**

Run:

```bash
git add .tasks/coverage-baseline.md
git commit -m "docs: record coverage baseline"
```

Expected: commit succeeds with only `.tasks/coverage-baseline.md` staged.

## Task 2: Coverage Enforcement Foundation

**Files:**

- Modify: `package.json`
- Modify: `apps/api/project.json`
- Modify: `apps/bot/project.json`
- Modify: `libs/go/config/project.json`
- Modify: `libs/go/logger/project.json`
- Create: `tools/coverage/coverage.config.json`
- Create: `tools/coverage/preflight.mjs`
- Create: `tools/coverage/run.mjs`

- [ ] **Step 1: Write the coverage config**

Create `tools/coverage/coverage.config.json`:

```json
{
  "thresholds": {
    "goStatements": 100,
    "typescriptStatements": 100,
    "typescriptBranches": 100,
    "typescriptFunctions": 100,
    "typescriptLines": 100
  },
  "goProjects": [
    {
      "name": "api",
      "cwd": "apps/api",
      "profile": "dist/coverage/go/api/coverage.out",
      "packages": "./..."
    },
    {
      "name": "bot",
      "cwd": "apps/bot",
      "profile": "dist/coverage/go/bot/coverage.out",
      "packages": "./..."
    },
    {
      "name": "go-config",
      "cwd": "libs/go/config",
      "profile": "dist/coverage/go/go-config/coverage.out",
      "packages": "./..."
    },
    {
      "name": "go-logger",
      "cwd": "libs/go/logger",
      "profile": "dist/coverage/go/go-logger/coverage.out",
      "packages": "./..."
    }
  ],
  "typescriptCoverageSummaries": [
    {
      "name": "web",
      "summary": "dist/coverage/web/coverage-summary.json"
    },
    {
      "name": "workspace-tools",
      "summary": "dist/coverage/tools/coverage-summary.json"
    }
  ],
  "allowlist": [
    {
      "path": "apps/api/internal/graph/generated.go",
      "reason": "gqlgen generated transport code",
      "gate": "bunx nx run api:codegen && bunx nx build api"
    },
    {
      "path": "apps/api/internal/graph/model/models_gen.go",
      "reason": "gqlgen generated model code",
      "gate": "bunx nx run api:codegen && bunx nx build api"
    },
    {
      "path": "apps/web/src/shared/api/generated/**",
      "reason": "GraphQL codegen output",
      "gate": "bunx nx run web:codegen && bunx nx run web:typecheck"
    },
    {
      "path": "apps/api/cmd/server/main.go",
      "reason": "API bootstrap entrypoint",
      "gate": "Playwright health/readiness startup scenarios"
    },
    {
      "path": "apps/bot/cmd/bot/main.go",
      "reason": "Bot bootstrap entrypoint",
      "gate": "bunx nx build bot"
    }
  ]
}
```

- [ ] **Step 2: Write the preflight script**

Create `tools/coverage/preflight.mjs`:

```js
import fs from 'node:fs';
import path from 'node:path';

const root = process.cwd();
const requiredFiles = [
  'tools/coverage/coverage.config.json',
  'package.json',
  'apps/web/vitest.config.ts',
  'apps/web/e2e/playwright.config.ts',
  'docs/verification-plan.xml',
];

const requiredScripts = ['test:coverage', 'test:e2e', 'verify:coverage'];

function fail(message) {
  console.error(`[Coverage][preflight] ${message}`);
  process.exitCode = 1;
}

for (const file of requiredFiles) {
  if (!fs.existsSync(path.join(root, file))) {
    fail(`Missing required file: ${file}`);
  }
}

const pkg = JSON.parse(fs.readFileSync(path.join(root, 'package.json'), 'utf8'));
for (const script of requiredScripts) {
  if (!pkg.scripts?.[script]) {
    fail(`Missing package.json script: ${script}`);
  }
}

if (process.exitCode) {
  process.exit(process.exitCode);
}

console.log('[Coverage][preflight] ok');
```

- [ ] **Step 3: Write the coverage runner**

Create `tools/coverage/run.mjs`:

```js
import fs from 'node:fs';
import path from 'node:path';
import { execFileSync } from 'node:child_process';

const root = process.cwd();
const configPath = path.join(root, 'tools/coverage/coverage.config.json');
const config = JSON.parse(fs.readFileSync(configPath, 'utf8'));

function fail(message) {
  console.error(`[Coverage][gate] ${message}`);
  process.exitCode = 1;
}

function run(command, args, options = {}) {
  console.log(`[Coverage][run] ${command} ${args.join(' ')}`);
  execFileSync(command, args, {
    cwd: options.cwd || root,
    stdio: 'inherit',
    env: { ...process.env, ...options.env },
  });
}

function ensureDirFor(file) {
  fs.mkdirSync(path.dirname(path.join(root, file)), { recursive: true });
}

function isAllowlistedGoFile(file) {
  return config.allowlist
    .filter((item) => item.path.endsWith('.go'))
    .some((item) => file.endsWith(item.path) || file.includes(item.path));
}

function parseGoFilteredTotal(profile) {
  const lines = fs.readFileSync(path.join(root, profile), 'utf8').trim().split('\n');
  let statements = 0;
  let covered = 0;

  for (const line of lines) {
    if (line === 'mode: set' || line === 'mode: count' || line === 'mode: atomic') {
      continue;
    }
    const match = line.match(/^(.+):\d+\.\d+,\d+\.\d+\s+(\d+)\s+(\d+)$/);
    if (!match) {
      throw new Error(`Cannot parse Go coverage line in ${profile}: ${line}`);
    }
    const [, file, statementCountRaw, hitCountRaw] = match;
    if (isAllowlistedGoFile(file)) {
      continue;
    }
    const statementCount = Number(statementCountRaw);
    const hitCount = Number(hitCountRaw);
    statements += statementCount;
    if (hitCount > 0) {
      covered += statementCount;
    }
  }

  if (statements === 0) {
    throw new Error(`No non-allowlisted Go statements found in ${profile}`);
  }

  return Number(((covered / statements) * 100).toFixed(1));
}

function readJson(file) {
  return JSON.parse(fs.readFileSync(path.join(root, file), 'utf8'));
}

run('node', ['tools/coverage/preflight.mjs']);

for (const project of config.goProjects) {
  ensureDirFor(project.profile);
  run('go', ['test', `-coverprofile=${path.join(root, project.profile)}`, project.packages], {
    cwd: path.join(root, project.cwd),
    env: { COVERAGE_GATE: '1' },
  });
  const total = parseGoFilteredTotal(project.profile);
  if (total !== config.thresholds.goStatements) {
    fail(`${project.name}: Go statement coverage ${total}% != ${config.thresholds.goStatements}%`);
  }
}

run('bunx', ['nx', 'run', 'web:test-coverage']);
run('bunx', ['nx', 'run', 'nx-go:test']);
run('bunx', ['nx', 'run', 'codegen:validate']);

for (const item of config.typescriptCoverageSummaries) {
  const summary = readJson(item.summary).total;
  const checks = [
    ['statements', summary.statements.pct, config.thresholds.typescriptStatements],
    ['branches', summary.branches.pct, config.thresholds.typescriptBranches],
    ['functions', summary.functions.pct, config.thresholds.typescriptFunctions],
    ['lines', summary.lines.pct, config.thresholds.typescriptLines],
  ];
  for (const [metric, actual, expected] of checks) {
    if (actual !== expected) {
      fail(`${item.name}: ${metric} coverage ${actual}% != ${expected}%`);
    }
  }
}

if (process.exitCode) {
  process.exit(process.exitCode);
}

console.log('[Coverage][gate] all thresholds passed');
```

- [ ] **Step 4: Update root package scripts and dependencies**

Modify `package.json`:

```json
{
  "scripts": {
    "dev": "bunx nx run-many --target=serve --all",
    "build": "bunx nx run-many --target=build --all",
    "test": "bunx nx run-many --target=test --all",
    "lint": "bunx nx run-many --target=lint --all --parallel=1",
    "codegen": "bunx nx run-many --target=codegen --projects=api,web",
    "test:coverage": "node tools/coverage/run.mjs",
    "test:e2e": "bunx nx run web:e2e",
    "verify:coverage": "bun run lint && bun run codegen && bunx nx run web:typecheck && bun run build && bun run test:coverage && bun run test:e2e && xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml && grace lint --path .",
    "postinstall": "lefthook install"
  },
  "devDependencies": {
    "@vitest/coverage-v8": "^2.1.9",
    "vitest": "^2.1.9"
  }
}
```

Keep existing `devDependencies` and add the two listed dependencies without removing current entries.

- [ ] **Step 5: Route Go coverage artifacts to `dist/coverage`**

Update Go project test commands exactly:

```json
{
  "command": "mkdir -p ../../dist/coverage/go/api && cd apps/api && go test -coverprofile=../../dist/coverage/go/api/coverage.out ./..."
}
```

For `apps/bot/project.json`:

```json
{
  "command": "mkdir -p ../../dist/coverage/go/bot && cd apps/bot && go test -coverprofile=../../dist/coverage/go/bot/coverage.out ./..."
}
```

For `libs/go/config/project.json`:

```json
{
  "command": "mkdir -p ../../../dist/coverage/go/go-config && cd libs/go/config && go test -coverprofile=../../../dist/coverage/go/go-config/coverage.out ./..."
}
```

For `libs/go/logger/project.json`:

```json
{
  "command": "mkdir -p ../../../dist/coverage/go/go-logger && cd libs/go/logger && go test -coverprofile=../../../dist/coverage/go/go-logger/coverage.out ./..."
}
```

- [ ] **Step 6: Run the red gate**

Run:

```bash
bun install
bun run test:coverage
```

Expected: FAIL. The failure must mention missing `web:test-coverage`, missing `nx-go:test`, missing `codegen:validate`, or coverage below 100. Record the first failure in `.tasks/coverage-baseline.md`.

- [ ] **Step 7: Commit the coverage foundation**

Run:

```bash
git add package.json bun.lock apps/api/project.json apps/bot/project.json libs/go/config/project.json libs/go/logger/project.json tools/coverage .tasks/coverage-baseline.md
git commit -m "feat: add coverage enforcement foundation"
```

Expected: commit succeeds while the full coverage gate is still red for known follow-up tasks.

## Task 3: Workspace Tooling Coverage

**Files:**

- Create: `tools/vitest.config.ts`
- Create: `tools/nx-go/project.json`
- Create: `tools/codegen/project.json`
- Create: `tools/nx-go/src/executors/test/executor.test.ts`
- Create: `tools/nx-go/src/executors/build/executor.test.ts`
- Create: `tools/nx-go/src/executors/lint/executor.test.ts`
- Create: `tools/nx-go/src/executors/serve/executor.test.ts`
- Create: `tools/codegen/codegen.test.ts`

- [ ] **Step 1: Add root tooling Vitest config**

Create `tools/vitest.config.ts`:

```ts
import { defineConfig } from 'vitest/config';
import { resolve } from 'node:path';

export default defineConfig({
  test: {
    environment: 'node',
    globals: true,
    include: ['tools/**/*.test.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json-summary'],
      reportsDirectory: resolve(__dirname, '../dist/coverage/tools'),
      include: ['tools/nx-go/src/**/*.ts', 'tools/codegen/**/*.ts'],
      exclude: ['tools/**/*.test.ts', 'tools/**/schema.json', 'tools/**/package.json'],
      thresholds: {
        statements: 100,
        branches: 100,
        functions: 100,
        lines: 100,
      },
    },
  },
});
```

- [ ] **Step 2: Add Nx targets for tooling**

Create `tools/nx-go/project.json`:

```json
{
  "name": "nx-go",
  "$schema": "../../node_modules/nx/schemas/project-schema.json",
  "sourceRoot": "tools/nx-go/src",
  "projectType": "library",
  "targets": {
    "test": {
      "executor": "nx:run-commands",
      "options": {
        "command": "bunx vitest run --config tools/vitest.config.ts --coverage tools/nx-go/src"
      }
    }
  }
}
```

Create `tools/codegen/project.json`:

```json
{
  "name": "codegen",
  "$schema": "../../node_modules/nx/schemas/project-schema.json",
  "sourceRoot": "tools/codegen",
  "projectType": "library",
  "targets": {
    "validate": {
      "executor": "nx:run-commands",
      "options": {
        "command": "bunx vitest run --config tools/vitest.config.ts --coverage tools/codegen && bunx nx run web:codegen"
      }
    }
  }
}
```

- [ ] **Step 3: Add executor tests**

Use Vitest module mocks. Each `execSync`-based executor test starts with:

```ts
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { ExecutorContext } from '@nx/devkit';
import runExecutor from './executor';

const execSyncMock = vi.fn();

vi.mock('child_process', () => ({
  execSync: execSyncMock,
}));

beforeEach(() => {
  vi.clearAllMocks();
});
```

Each `spawn`-based executor test starts with:

```ts
import { EventEmitter } from 'node:events';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { ExecutorContext } from '@nx/devkit';
import runExecutor from './executor';

const spawnMock = vi.fn();

vi.mock('child_process', () => ({
  spawn: spawnMock,
}));

beforeEach(() => {
  vi.clearAllMocks();
});
```

For each executor test file, build a minimal Nx context:

```ts
const context = {
  root: '/repo',
  projectName: 'api',
  projectsConfigurations: {
    projects: {
      api: { root: 'apps/api' },
    },
  },
};
```

`tools/nx-go/src/executors/test/executor.test.ts` must cover:

```ts
it('runs go test with coverage, short mode, and explicit packages', async () => {
  execSyncMock.mockReturnValue(Buffer.from('ok'));
  const result = await runExecutor(
    { coverage: true, short: true, packages: ['./internal/service', './internal/handler'] },
    context as ExecutorContext,
  );
  expect(result).toEqual({ success: true });
  expect(execSyncMock).toHaveBeenCalledWith(
    'go test -short -coverprofile=coverage.out ./internal/service ./internal/handler',
    expect.objectContaining({ cwd: '/repo/apps/api', stdio: 'inherit' }),
  );
});

it('returns false when go test fails', async () => {
  execSyncMock.mockImplementation(() => {
    throw new Error('failed');
  });
  await expect(
    runExecutor({ coverage: false, short: false }, context as ExecutorContext),
  ).resolves.toEqual({ success: false });
});
```

Create `tools/nx-go/src/executors/build/executor.test.ts` with success and failure tests:

```ts
it('runs go build with output path and main package', async () => {
  execSyncMock.mockReturnValue(Buffer.from('ok'));
  const result = await runExecutor(
    { outputPath: 'dist/apps/api', main: 'cmd/server' },
    context as ExecutorContext,
  );
  expect(result).toEqual({ success: true });
  expect(execSyncMock).toHaveBeenCalledWith(
    'go build -o /repo/dist/apps/api ./cmd/server',
    expect.objectContaining({ cwd: '/repo/apps/api', stdio: 'inherit' }),
  );
});

it('returns false when go build fails', async () => {
  execSyncMock.mockImplementation(() => {
    throw new Error('failed');
  });
  await expect(
    runExecutor({ outputPath: 'dist/apps/api', main: 'cmd/server' }, context as ExecutorContext),
  ).resolves.toEqual({ success: false });
});
```

Create `tools/nx-go/src/executors/lint/executor.test.ts` with base, `--fix`, config, and failure tests:

```ts
it('runs golangci-lint with fix and config flags', async () => {
  execSyncMock.mockReturnValue(Buffer.from('ok'));
  const result = await runExecutor(
    { fix: true, config: '.golangci.yml' },
    context as ExecutorContext,
  );
  expect(result).toEqual({ success: true });
  expect(execSyncMock).toHaveBeenCalledWith(
    'golangci-lint run --fix --config=.golangci.yml',
    expect.objectContaining({ cwd: '/repo/apps/api', stdio: 'inherit' }),
  );
});

it('returns false when golangci-lint fails', async () => {
  execSyncMock.mockImplementation(() => {
    throw new Error('failed');
  });
  await expect(runExecutor({ fix: false }, context as ExecutorContext)).resolves.toEqual({
    success: false,
  });
});
```

Create `tools/nx-go/src/executors/serve/executor.test.ts` with spawn tests:

```ts
it('spawns air with config path and API_PORT', async () => {
  const child = new EventEmitter() as EventEmitter & { kill: ReturnType<typeof vi.fn> };
  child.kill = vi.fn();
  spawnMock.mockReturnValue(child);
  const promise = runExecutor(
    { port: 8080, configPath: 'custom.air.toml' },
    context as ExecutorContext,
  );
  child.emit('close', 0);
  await expect(promise).resolves.toEqual({ success: true });
  expect(spawnMock).toHaveBeenCalledWith(
    'air',
    ['-c', 'custom.air.toml'],
    expect.objectContaining({
      cwd: '/repo/apps/api',
      stdio: 'inherit',
      env: expect.objectContaining({ API_PORT: '8080' }),
    }),
  );
});

it('returns false for non-zero close code', async () => {
  const child = new EventEmitter() as EventEmitter & { kill: ReturnType<typeof vi.fn> };
  child.kill = vi.fn();
  spawnMock.mockReturnValue(child);
  const promise = runExecutor({ port: 8080 }, context as ExecutorContext);
  child.emit('close', 1);
  await expect(promise).resolves.toEqual({ success: false });
});
```

- [ ] **Step 4: Add codegen config test**

Create `tools/codegen/codegen.test.ts`:

```ts
import { describe, expect, it } from 'vitest';
import config from './codegen';

describe('codegen config', () => {
  it('uses shared GraphQL schema and web operation documents', () => {
    expect(config.schema).toBe('../../libs/graphql/schema/**/*.graphql');
    expect(config.documents).toEqual([
      '../../apps/web/src/features/**/api/**/*.graphql',
      '../../apps/web/src/entities/**/api/**/*.graphql',
    ]);
  });

  it('generates web types with expected scalar mappings', () => {
    const output = config.generates?.['../../apps/web/src/shared/api/generated/types.ts'];
    expect(output).toBeDefined();
    expect(output).toMatchObject({
      plugins: ['typescript', 'typescript-operations'],
      config: {
        scalars: {
          DateTime: 'string',
          UUID: 'string',
        },
      },
    });
  });
});
```

- [ ] **Step 5: Run tooling tests**

Run:

```bash
bunx nx run nx-go:test
bunx nx run codegen:validate
```

Expected: both pass and `dist/coverage/tools/coverage-summary.json` exists.

- [ ] **Step 6: Commit tooling coverage**

Run:

```bash
git add tools/vitest.config.ts tools/nx-go/project.json tools/codegen/project.json tools/nx-go/src/**/*.test.ts tools/codegen/codegen.test.ts dist/coverage/tools/.gitkeep package.json bun.lock
git commit -m "test: cover workspace tooling"
```

Do not commit generated coverage reports. If `dist/coverage/tools/.gitkeep` is not needed because `dist` is ignored, omit it from `git add`.

## Task 4: API CRUD and Go Coverage Closure

**Files:**

- Modify: `apps/api/internal/repository/postgres/user_repo.go`
- Create: `apps/api/internal/repository/postgres/user_repo_test.go`
- Modify: `apps/api/internal/service/user_service_test.go`
- Modify: `apps/api/internal/graph/schema.resolvers.go`
- Create: `apps/api/internal/graph/schema_resolvers_test.go`
- Create: `apps/api/internal/repository/postgres/postgres_test.go`
- Create: `apps/api/internal/repository/redis/cache_test.go`

- [ ] **Step 1: Add repository integration test helper**

Create `apps/api/internal/repository/postgres/user_repo_test.go` with helper setup:

```go
package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	postgresrepo "monorepo-template/apps/api/internal/repository/postgres"
	"monorepo-template/apps/api/internal/service"
	"go.uber.org/zap"
)

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("API_TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://app:secret@localhost:7501/monorepo_dev?sslmode=disable"
	}
	if err := postgresrepo.RunMigrations(dsn, zap.NewNop()); err != nil {
		if os.Getenv("COVERAGE_GATE") != "1" {
			t.Skipf("postgres integration database is unavailable: %v", err)
		}
		require.NoError(t, err)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	_, err = pool.Exec(context.Background(), `TRUNCATE users RESTART IDENTITY`)
	require.NoError(t, err)
	return pool
}
```

- [ ] **Step 2: Add failing repository tests**

In the same file, add tests for:

```go
func TestUserRepo_CreateGetListUpdateDelete(t *testing.T) {
	ctx := context.Background()
	repo := postgresrepo.NewUserRepo(testPool(t))

	created, err := repo.Create(ctx, service.CreateUserInput{
		Email: "created@example.com",
		Name: "Created User",
		Password: "$2a$10$hashed",
	})
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, "created@example.com", found.Email)

	name := "Updated User"
	email := "updated@example.com"
	updated, err := repo.Update(ctx, created.ID, service.UpdateUserInput{Name: &name, Email: &email})
	require.NoError(t, err)
	require.Equal(t, "Updated User", updated.Name)
	require.Equal(t, "updated@example.com", updated.Email)

	users, total, err := repo.List(ctx, ptr(20), nil)
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, users, 1)

	require.NoError(t, repo.Delete(ctx, created.ID))
	deleted, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.Nil(t, deleted)
}

func ptr[T any](v T) *T { return &v }
```

Add separate tests for duplicate create, update duplicate email, update missing ID returning `nil, nil`, delete missing ID returning `nil`, and invalid list cursor returning an error containing `invalid cursor`.

Run:

```bash
cd apps/api && go test ./internal/repository/postgres -run TestUserRepo -count=1
```

Expected: FAIL because `Update` and `Delete` are not implemented.

- [ ] **Step 3: Implement repository update/delete**

Modify `apps/api/internal/repository/postgres/user_repo.go`:

```go
func (r *UserRepo) Update(ctx context.Context, id string, input service.UpdateUserInput) (*service.User, error) {
	const op = "UserRepo.Update"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("updating user", zap.String("user_id", id))

	var u service.User
	var createdAt, updatedAt time.Time
	err := r.pool.QueryRow(ctx,
		`UPDATE users
		 SET name = COALESCE($2, name),
		     email = COALESCE($3, email),
		     updated_at = NOW()
		 WHERE id = $1
		 RETURNING id, email, name, created_at, updated_at`,
		id, input.Name, input.Email,
	).Scan(&u.ID, &u.Email, &u.Name, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		if isDuplicateKeyError(err) {
			return nil, fmt.Errorf("%s: duplicate email: %s", op, derefString(input.Email))
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	u.CreatedAt = createdAt.Format(time.RFC3339Nano)
	u.UpdatedAt = updatedAt.Format(time.RFC3339Nano)
	return &u, nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	const op = "UserRepo.Delete"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("deleting user", zap.String("user_id", id))

	tag, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return nil
	}
	return nil
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
```

- [ ] **Step 4: Run repository tests**

Run:

```bash
cd apps/api && go test ./internal/repository/postgres -run TestUserRepo -count=1
```

Expected: PASS.

- [ ] **Step 5: Expand service tests**

Modify `apps/api/internal/service/user_service_test.go` so `mockUserRepo` stores `lastCreateInput`, `err`, and `deletedIDs`. Add tests:

```go
func TestUserService_Create_HashesPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := service.NewUserService(repo)

	_, err := svc.Create(context.Background(), service.CreateUserInput{
		Email: "test@example.com",
		Name: "Test User",
		Password: "plain-password",
	})

	require.NoError(t, err)
	require.NotEqual(t, "plain-password", repo.lastCreateInput.Password)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(repo.lastCreateInput.Password), []byte("plain-password")))
}

func TestUserService_Update_DelegatesToRepo(t *testing.T) {
	repo := newMockUserRepo()
	svc := service.NewUserService(repo)
	name := "Updated"
	user, err := svc.Update(context.Background(), "test-id", service.UpdateUserInput{Name: &name})
	require.NoError(t, err)
	require.Equal(t, "Updated", user.Name)
}

func TestUserService_Delete_DelegatesToRepo(t *testing.T) {
	repo := newMockUserRepo()
	svc := service.NewUserService(repo)
	require.NoError(t, svc.Delete(context.Background(), "test-id"))
	require.Contains(t, repo.deletedIDs, "test-id")
}
```

Add error propagation tests for `GetByID`, `List`, `Create`, `Update`, and `Delete` by setting `repo.err = errors.New("repo failed")`.

- [ ] **Step 6: Implement resolver tests and resolver behavior**

Add `apps/api/internal/graph/schema_resolvers_test.go` in package `graph`.

Test names:

- `TestCreateUser_ReturnsSuccess`
- `TestCreateUser_ReturnsValidationErrorOnDuplicate`
- `TestUpdateUser_ReturnsSuccess`
- `TestUpdateUser_ReturnsNotFound`
- `TestUpdateUser_ReturnsValidationErrorOnDuplicate`
- `TestDeleteUser_ReturnsTrueForExistingUser`
- `TestDeleteUser_ReturnsFalseForMissingUser`
- `TestUsers_ReturnsPagination`
- `TestUser_ReturnsNilForMissingUser`

Implement resolver behavior in `schema.resolvers.go`:

```go
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUserInput) (model.UpdateUserResult, error) {
	const op = "mutationResolver.UpdateUser"
	u, err := r.UserService.Update(ctx, id, mapUpdateUserInput(input))
	if err != nil {
		if strings.Contains(err.Error(), "duplicate email") {
			return model.ValidationError{Field: "email", Message: "already exists"}, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if u == nil {
		return model.NotFoundError{Message: "user not found", EntityType: "User", ID: id}, nil
	}
	return model.UpdateUserSuccess{User: mapUser(u)}, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	const op = "mutationResolver.DeleteUser"
	before, err := r.UserService.GetByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("%s: lookup: %w", op, err)
	}
	if before == nil {
		return false, nil
	}
	if err := r.UserService.Delete(ctx, id); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}

func mapUpdateUserInput(in model.UpdateUserInput) service.UpdateUserInput {
	return service.UpdateUserInput{Name: in.Name, Email: in.Email}
}
```

- [ ] **Step 7: Add DB and Redis integration tests**

Create `apps/api/internal/repository/postgres/postgres_test.go`:

```go
func TestNew_ConnectsAndPings(t *testing.T) {
	cfg := config.PostgresConfig{Host: "localhost", Port: 7501, User: "app", Password: "secret", DB: "monorepo_dev", SSLMode: "disable"}
	db, err := postgres.New(cfg, zap.NewNop())
	if err != nil && os.Getenv("COVERAGE_GATE") != "1" {
		t.Skipf("postgres integration database is unavailable: %v", err)
	}
	require.NoError(t, err)
	t.Cleanup(db.Close)
	require.NoError(t, db.Ping())
}

func TestNew_ReturnsErrorForBadPort(t *testing.T) {
	cfg := config.PostgresConfig{Host: "localhost", Port: 1, User: "app", Password: "secret", DB: "monorepo_dev", SSLMode: "disable"}
	_, err := postgres.New(cfg, zap.NewNop())
	require.Error(t, err)
}
```

Create `apps/api/internal/repository/redis/cache_test.go`:

```go
func TestNew_ConnectsPingsAndCloses(t *testing.T) {
	client, err := redisrepo.New(config.RedisConfig{Host: "localhost", Port: 7502}, zap.NewNop())
	if err != nil && os.Getenv("COVERAGE_GATE") != "1" {
		t.Skipf("redis integration service is unavailable: %v", err)
	}
	require.NoError(t, err)
	require.NoError(t, client.Ping())
	require.NoError(t, client.Close())
}

func TestNew_ReturnsErrorForBadPort(t *testing.T) {
	_, err := redisrepo.New(config.RedisConfig{Host: "localhost", Port: 1}, zap.NewNop())
	require.Error(t, err)
}
```

- [ ] **Step 8: Run Go checks**

Run:

```bash
docker compose -f docker/docker-compose.dev.yml up -d postgres redis
cd apps/api && go test -coverprofile=../../dist/coverage/go/api/coverage.out ./...
cd apps/bot && go test -coverprofile=../../dist/coverage/go/bot/coverage.out ./...
cd libs/go/config && go test -coverprofile=../../../dist/coverage/go/go-config/coverage.out ./...
cd libs/go/logger && go test -coverprofile=../../../dist/coverage/go/go-logger/coverage.out ./...
```

Expected: PASS with 100 percent Go statement coverage. A lower total blocks the task; record the uncovered function from `go tool cover -func=<profile>` and add the missing package-local test before repeating Step 8.

- [ ] **Step 9: Commit Go coverage closure**

Run:

```bash
git add apps/api/internal/repository/postgres apps/api/internal/repository/redis apps/api/internal/service/user_service_test.go apps/api/internal/graph/schema.resolvers.go apps/api/internal/graph/schema_resolvers_test.go
git commit -m "test: close api go coverage"
```

## Task 5: Web Unit and Component Coverage Closure

**Files:**

- Modify: `apps/web/package.json`
- Modify: `apps/web/project.json`
- Modify: `apps/web/vitest.config.ts`
- Modify: `apps/web/app/__tests__/page.test.tsx`
- Create: `apps/web/app/__tests__/providers.test.tsx`
- Create: `apps/web/app/__tests__/users-page.test.tsx`
- Create: `apps/web/app/__tests__/user-detail-page.test.tsx`
- Create: `apps/web/src/shared/api/graphql-client.test.ts`
- Create: `apps/web/src/shared/config/index.test.ts`

- [ ] **Step 1: Raise web coverage thresholds**

Modify `apps/web/package.json`:

```json
{
  "scripts": {
    "test": "vitest run",
    "test:watch": "vitest",
    "test:coverage": "vitest run --coverage",
    "codegen": "graphql-codegen --config ../../tools/codegen/codegen.ts"
  }
}
```

Modify `apps/web/project.json` to add a `test-coverage` target:

```json
{
  "test-coverage": {
    "executor": "nx:run-commands",
    "options": {
      "command": "cd apps/web && bun run test:coverage"
    }
  }
}
```

Modify `apps/web/vitest.config.ts`:

```ts
coverage: {
  provider: 'v8',
  reporter: ['text', 'json-summary'],
  reportsDirectory: '../../dist/coverage/web',
  include: ['app/**/*.{ts,tsx}', 'src/**/*.{ts,tsx}'],
  exclude: [
    'app/**/*.test.{ts,tsx}',
    'src/**/*.test.{ts,tsx}',
    'src/shared/api/generated/**',
    'next-env.d.ts',
  ],
  thresholds: {
    statements: 100,
    branches: 100,
    functions: 100,
    lines: 100,
  },
}
```

- [ ] **Step 2: Add failing users page tests**

Create `apps/web/app/__tests__/users-page.test.tsx`. Mock `@shared/api/graphql-client`, `next/link`, and use a fresh `QueryClientProvider` per test.

Required tests:

- `renders loading state`
- `renders empty state`
- `renders returned users`
- `shows load error`
- `creates a user and clears the form`
- `shows validation error from createUser`
- `shows auth error from createUser`

Use this response shape for success:

```ts
{
  users: {
    edges: [
      {
        cursor: 'cursor-1',
        node: {
          id: 'user-1',
          email: 'one@example.com',
          name: 'One User',
          createdAt: '2026-05-02T00:00:00Z',
        },
      },
    ],
    pageInfo: { hasNextPage: false, endCursor: null },
    totalCount: 1,
  },
}
```

- [ ] **Step 3: Add user detail tests**

Create `apps/web/app/__tests__/user-detail-page.test.tsx`. Mock `global.fetch`.

Required tests:

- successful fetch renders name, email, dates, and ID.
- non-OK response renders `User not found`.
- GraphQL `user: null` renders `User not found`.

Call the async component directly:

```ts
const ui = await UserDetailPage({ params: Promise.resolve({ id: 'user-1' }) });
render(ui);
```

- [ ] **Step 4: Add provider, config, and client tests**

Create `apps/web/app/__tests__/providers.test.tsx`:

```ts
it('renders children inside QueryClientProvider', () => {
  render(<Providers><span>child content</span></Providers>);
  expect(screen.getByText('child content')).toBeInTheDocument();
});
```

Create `apps/web/src/shared/config/index.test.ts` with default-env assertions. Use `vi.stubEnv` and dynamic import to cover env override behavior.

Refactor `apps/web/src/shared/api/graphql-client.ts` so the client factory is testable:

```ts
import { GraphQLClient } from 'graphql-request';
import { appConfig } from '@shared/config';

export function createGraphQLClient(apiUrl = appConfig.apiUrl) {
  return new GraphQLClient(apiUrl, { headers: {} });
}

export const graphqlClient = createGraphQLClient();

export function setAuthToken(token: string, client = graphqlClient) {
  client.setHeader('Authorization', `Bearer ${token}`);
}
```

Create `apps/web/src/shared/api/graphql-client.test.ts`:

```ts
import { describe, expect, it, vi } from 'vitest';
import { createGraphQLClient, setAuthToken } from './graphql-client';

describe('graphql client', () => {
  it('creates a GraphQLClient for the supplied URL', () => {
    const client = createGraphQLClient('http://example.test/graphql');
    expect(client).toBeDefined();
  });

  it('sets bearer token header on the supplied client', () => {
    const client = { setHeader: vi.fn() };
    setAuthToken('abc123', client as never);
    expect(client.setHeader).toHaveBeenCalledWith('Authorization', 'Bearer abc123');
  });
});
```

- [ ] **Step 5: Run web coverage**

Run:

```bash
bunx nx run web:test-coverage
```

Expected: PASS and `dist/coverage/web/coverage-summary.json` has 100 percent for statements, branches, functions, and lines.

- [ ] **Step 6: Commit web coverage closure**

Run:

```bash
git add apps/web/package.json apps/web/project.json apps/web/vitest.config.ts apps/web/app/__tests__ apps/web/src/shared/api apps/web/src/shared/config package.json bun.lock
git commit -m "test: close web coverage"
```

## Task 6: E2E Matrix Closure

**Files:**

- Modify: `apps/web/e2e/playwright.config.ts`
- Modify: `apps/web/project.json`
- Create: `apps/web/e2e/helpers.ts`
- Create: `apps/web/e2e/preflight.mjs`
- Create: `apps/web/e2e/graphql-contract.spec.ts`
- Create: `apps/web/e2e/users-flow.spec.ts`

- [ ] **Step 1: Update Playwright config**

Create `apps/web/e2e/preflight.mjs`:

```js
import { execFileSync } from 'node:child_process';

export default async function globalSetup() {
  execFileSync(
    'docker',
    ['compose', '-f', 'docker/docker-compose.dev.yml', 'up', '-d', 'postgres', 'redis'],
    {
      cwd: '../..',
      stdio: 'inherit',
    },
  );

  execFileSync('docker', ['compose', '-f', 'docker/docker-compose.dev.yml', 'ps'], {
    cwd: '../..',
    stdio: 'inherit',
  });

  console.log('[E2E][preflight] local infrastructure started');
}
```

Modify `apps/web/project.json` so `web:e2e` uses the nested config explicitly:

```json
{
  "command": "cd apps/web && bunx playwright test --config=e2e/playwright.config.ts"
}
```

Modify `apps/web/e2e/playwright.config.ts`:

```ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: '.',
  globalSetup: './preflight.mjs',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: [['html', { outputFolder: '../../dist/coverage/e2e/html-report' }], ['list']],
  use: {
    baseURL: 'http://localhost:3000',
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
      command: 'cd apps/api && go run ./cmd/server',
      url: 'http://localhost:8080/healthz',
      reuseExistingServer: !process.env.CI,
      cwd: '../..',
      timeout: 120_000,
    },
    {
      command: 'cd apps/web && bun run dev',
      url: 'http://localhost:3000',
      reuseExistingServer: !process.env.CI,
      cwd: '../..',
      timeout: 120_000,
    },
  ],
});
```

- [ ] **Step 2: Add e2e helpers**

Create `apps/web/e2e/helpers.ts`:

```ts
import { expect, request } from '@playwright/test';

export const apiURL = 'http://localhost:8080';

export async function gql<T>(query: string, variables: Record<string, unknown> = {}) {
  const ctx = await request.newContext({ baseURL: apiURL });
  const res = await ctx.post('/graphql', { data: { query, variables } });
  expect(res.ok()).toBeTruthy();
  const json = await res.json();
  await ctx.dispose();
  return json as T;
}

export function uniqueEmail(prefix: string) {
  return `${prefix}-${Date.now()}-${Math.random().toString(16).slice(2)}@example.com`;
}

export async function createUser(name: string, email = uniqueEmail('e2e')) {
  const json = await gql<{
    data: { createUser: { user: { id: string; email: string; name: string } } };
  }>(
    `mutation CreateUser($input: CreateUserInput!) {
      createUser(input: $input) {
        ... on CreateUserSuccess { user { id email name } }
        ... on ValidationError { field message }
        ... on AuthError { message }
      }
    }`,
    { input: { name, email, password: 'Password123!' } },
  );
  return json.data.createUser.user;
}
```

- [ ] **Step 3: Add GraphQL contract e2e tests**

Create `apps/web/e2e/graphql-contract.spec.ts`.

Required tests:

- `/healthz` returns `{ status: "ok" }`.
- `/readyz` returns `{ status: "ok" }`.
- `createUser` returns `CreateUserSuccess`.
- `users` includes created user.
- `user(id)` returns created user.
- `updateUser` changes name/email.
- duplicate create returns `ValidationError`.
- `deleteUser` returns `true`, then `user(id)` returns `null`.

- [ ] **Step 4: Add browser flow e2e tests**

Create `apps/web/e2e/users-flow.spec.ts`.

Required tests:

```ts
test('user can create, list, and open detail page', async ({ page }) => {
  const email = uniqueEmail('browser');
  await page.goto('/users');
  await page.getByPlaceholder('Name').fill('Browser User');
  await page.getByPlaceholder('Email').fill(email);
  await page.getByPlaceholder('Password').fill('Password123!');
  await page.getByRole('button', { name: 'Create' }).click();
  await expect(page.getByText(email)).toBeVisible();
  await page.getByRole('link', { name: 'Browser User' }).click();
  await expect(page.getByText(email)).toBeVisible();
});
```

Add a second test for duplicate email that creates a user through helper, submits the same email in the UI, and expects text containing `email: already exists`.

- [ ] **Step 5: Run e2e**

Run:

```bash
bunx nx run web:e2e
```

Expected: PASS. HTML report is written under `dist/coverage/e2e/html-report`.

- [ ] **Step 6: Commit e2e closure**

Run:

```bash
git add apps/web/e2e apps/web/project.json
git commit -m "test: add full stack e2e coverage"
```

## Task 7: GRACE and Documentation Sync

**Files:**

- Modify: `docs/requirements.xml`
- Modify: `docs/technology.xml`
- Modify: `docs/development-plan.xml`
- Modify: `docs/knowledge-graph.xml`
- Modify: `docs/verification-plan.xml`
- Modify: `docs/operational-packets.xml`
- Modify: `README.md`

- [ ] **Step 1: Update requirements**

In `docs/requirements.xml`, add a high-priority use case for the coverage gate:

```xml
<UC-005>
  <Actor>Developer</Actor>
  <Action>Runs the full coverage verification gate before using or releasing the template.</Action>
  <Goal>Guarantee handwritten production code has 100 percent coverage and full-stack e2e scenarios pass.</Goal>
  <Preconditions>Dependencies are installed and local Docker infrastructure is available.</Preconditions>
  <AcceptanceCriteria>`bun run verify:coverage` passes, generated artifacts are current, and Playwright e2e scenarios cover the required user CRUD flow.</AcceptanceCriteria>
  <Priority>high</Priority>
  <RelatedFlows>DF-COVERAGE-GATE</RelatedFlows>
</UC-005>
```

- [ ] **Step 2: Update technology**

In `docs/technology.xml`, add tooling entries:

```xml
<tool name="coverage-enforcement" value="tools/coverage" version="repo-local" />
<tool name="typescript-coverage" value="vitest v8 coverage" version="^2.1.9" />
<tool name="go-coverage" value="go test -coverprofile plus go tool cover" version="Go 1.25" />
```

Add testing policy text stating Go uses 100 percent statement coverage and TypeScript uses 100 percent statements/branches/functions/lines.

- [ ] **Step 3: Update development plan and graph**

In `docs/development-plan.xml`, add module `M-COVERAGE-GATE` or extend `M-WORKSPACE` with:

```xml
<export-verify-coverage PURPOSE="Run lint, codegen, typecheck, build, coverage, e2e, XML lint, and GRACE lint as the full template gate." />
```

In `docs/knowledge-graph.xml`, add paths:

```xml
<path>tools/coverage</path>
<path>tools/vitest.config.ts</path>
<path>apps/web/e2e</path>
```

Add cross-links from `M-COVERAGE-GATE` to `M-WORKSPACE`, `M-API`, `M-WEB`, `M-BOT`, `M-GRAPHQL-SCHEMA`, `M-GO-CONFIG`, and `M-GO-LOGGER`.

- [ ] **Step 4: Update verification plan**

In `docs/verification-plan.xml`, add:

```xml
<VF-COVERAGE-GATE NAME="OneHundredPercentCoverageGate" USE_CASES="UC-005" DATA_FLOW="DF-COVERAGE-GATE" PRIORITY="high">
  <scenario>Developer runs the root coverage verification command on the template.</scenario>
  <expected-outcome>Handwritten Go and TypeScript coverage is 100 percent, generated artifacts are validated, and full-stack e2e scenarios pass.</expected-outcome>
  <required-signals>
    <log-marker>[Coverage][gate] all thresholds passed</log-marker>
    <trace-sequence>preflight -> Go coverage -> web coverage -> tooling coverage -> e2e -> XML lint -> GRACE lint</trace-sequence>
  </required-signals>
</VF-COVERAGE-GATE>
```

Add the explicit allowlist entries from `tools/coverage/coverage.config.json` to the verification plan.

- [ ] **Step 5: Update operational packets**

In `docs/operational-packets.xml`, add a packet note:

```xml
<note-coverage>Workers closing coverage must list uncovered files, tests added, command output, coverage artifact path, and allowlist changes. A worker must stop if it needs to add a broad coverage exclusion glob.</note-coverage>
```

- [ ] **Step 6: Update README**

Add a section:

````markdown
## Coverage Gate

Run the full template gate:

```bash
docker compose -f docker/docker-compose.dev.yml up -d postgres redis
bun run verify:coverage
```

Coverage artifacts are written to `dist/coverage`.

The template enforces 100 percent coverage for handwritten production code. Generated files are excluded from line coverage and validated through codegen, typecheck, build, and e2e gates.
````

- [ ] **Step 7: Validate docs**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected: both pass.

- [ ] **Step 8: Commit docs sync**

Run:

```bash
git add docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml docs/operational-packets.xml README.md
git commit -m "docs: sync coverage verification contract"
```

## Task 8: Final Verification Gate

**Files:**

- Read: all changed files
- Modify only if a command exposes a real bug in an earlier task

- [ ] **Step 1: Run deterministic Nx reset**

Run:

```bash
bunx nx reset
```

Expected: Nx cache reset succeeds.

- [ ] **Step 2: Run deterministic root targets**

Run:

```bash
NX_SKIP_NX_CACHE=1 NX_DAEMON=false bunx nx run-many --target=lint --all --parallel=1
NX_SKIP_NX_CACHE=1 NX_DAEMON=false bunx nx run-many --target=test --all
NX_SKIP_NX_CACHE=1 NX_DAEMON=false bunx nx run-many --target=build --all
```

Expected: all pass.

- [ ] **Step 3: Run generated artifact gates**

Run:

```bash
bun run codegen
bunx nx run web:typecheck
git diff --exit-code -- apps/api/internal/graph apps/web/src/shared/api/generated
```

Expected: commands pass and generated artifacts have no uncommitted drift.

- [ ] **Step 4: Run coverage and e2e gates**

Run:

```bash
docker compose -f docker/docker-compose.dev.yml up -d postgres redis
bun run test:coverage
bun run test:e2e
bun run verify:coverage
```

Expected:

- `test:coverage` prints `[Coverage][gate] all thresholds passed`.
- `test:e2e` passes.
- `verify:coverage` passes.

- [ ] **Step 5: Run docs gates**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected: both pass.

- [ ] **Step 6: Record final evidence**

Create `.tasks/coverage-final-verification.md`:

```markdown
# 100 Percent Coverage Final Verification

## Commands

| Command | Exit | Evidence |
| ------- | ---: | -------- |

## Coverage Artifacts

| Artifact                                    | Status  |
| ------------------------------------------- | ------- |
| `dist/coverage/go/api/coverage.out`         | present |
| `dist/coverage/go/bot/coverage.out`         | present |
| `dist/coverage/go/go-config/coverage.out`   | present |
| `dist/coverage/go/go-logger/coverage.out`   | present |
| `dist/coverage/web/coverage-summary.json`   | present |
| `dist/coverage/tools/coverage-summary.json` | present |
| `dist/coverage/e2e/html-report`             | present |

## Result

The root coverage contract passed.
```

Fill the command table with the exact commands from steps 1-5 and the observed exit codes.

- [ ] **Step 7: Final commit**

Run:

```bash
git add .tasks/coverage-final-verification.md
git commit -m "docs: record coverage verification evidence"
```

Expected: commit succeeds and `git status --short` is clean.
