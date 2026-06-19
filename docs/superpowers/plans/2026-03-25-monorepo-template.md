# Monorepo Template Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a production-ready monorepo template with Go backend, Next.js frontend, GraphQL API contract, Nx orchestration, and full CI/CD pipeline.

**Architecture:** Nx-centric monorepo with custom Go executor. Go API (chi + gqlgen) follows clean architecture. Next.js frontend follows Feature-Sliced Design. GraphQL schema in shared lib is the single source of truth, with codegen to both sides.

**Tech Stack:** Nx, pnpm, Go, chi, gqlgen, zap, goose, viper, Next.js (App Router), TanStack Query, graphql-request, graphql-codegen, Tailwind CSS, PostgreSQL, Redis, Docker Compose, GitLab CI/CD, lefthook, commitlint, ESLint, Prettier, golangci-lint, Vitest, Playwright, testify

**Spec:** `docs/superpowers/specs/2026-03-25-monorepo-template-design.md`

---

## File Map

### Root configs (Task 1)
- Create: `package.json`
- Create: `pnpm-workspace.yaml`
- Create: `nx.json`
- Create: `tsconfig.base.json`
- Create: `.gitignore`
- Create: `.env.example`
- Create: `.npmrc`

### Linting, hooks, conventional commits (Task 2)
- Create: `.eslintrc.json`
- Create: `.prettierrc`
- Create: `.prettierignore`
- Create: `commitlint.config.js`
- Create: `.lefthook.yml`
- Create: `.lintstagedrc.json`

### GraphQL schema (Task 3)
- Create: `libs/graphql/schema/common.graphql`
- Create: `libs/graphql/schema/user.graphql`
- Create: `libs/graphql/schema/schema.graphql`
- Create: `libs/graphql/project.json`

### Custom Nx Go executor (Task 4)
- Create: `tools/nx-go/package.json`
- Create: `tools/nx-go/tsconfig.json`
- Create: `tools/nx-go/executors.json`
- Create: `tools/nx-go/src/index.ts`
- Create: `tools/nx-go/src/executors/build/executor.ts`
- Create: `tools/nx-go/src/executors/build/schema.json`
- Create: `tools/nx-go/src/executors/serve/executor.ts`
- Create: `tools/nx-go/src/executors/serve/schema.json`
- Create: `tools/nx-go/src/executors/test/executor.ts`
- Create: `tools/nx-go/src/executors/test/schema.json`
- Create: `tools/nx-go/src/executors/lint/executor.ts`
- Create: `tools/nx-go/src/executors/lint/schema.json`

### Go API (Tasks 5-8)
- Create: `apps/api/go.mod`
- Create: `apps/api/gqlgen.yml`
- Create: `apps/api/air.toml`
- Create: `apps/api/project.json`
- Create: `apps/api/.golangci.yml`
- Create: `apps/api/config/config.yml`
- Create: `apps/api/internal/config/config.go`
- Create: `apps/api/internal/config/config_test.go`
- Create: `apps/api/internal/handler/health.go`
- Create: `apps/api/internal/handler/health_test.go`
- Create: `apps/api/internal/middleware/logging.go`
- Create: `apps/api/internal/middleware/cors.go`
- Create: `apps/api/internal/middleware/auth.go`
- Create: `apps/api/internal/middleware/cors_test.go`
- Create: `apps/api/internal/middleware/auth_test.go`
- Create: `apps/api/internal/repository/postgres/postgres.go`
- Create: `apps/api/internal/repository/postgres/user_repo.go`
- Create: `apps/api/internal/repository/postgres/migrations/00001_init.sql`
- Create: `apps/api/internal/repository/redis/cache.go`
- Create: `apps/api/internal/service/user_service.go`
- Create: `apps/api/internal/service/user_service_test.go`
- Create: `apps/api/tools.go`
- Create: `apps/api/internal/graph/resolver.go`
- Create: `apps/api/cmd/server/main.go`

### Next.js frontend (Task 9)
- Create: `apps/web/package.json`
- Create: `apps/web/tsconfig.json`
- Create: `apps/web/next.config.js`
- Create: `apps/web/tailwind.config.ts`
- Create: `apps/web/postcss.config.js`
- Create: `apps/web/project.json`
- Create: `apps/web/vitest.config.ts`
- Create: `apps/web/app/layout.tsx`
- Create: `apps/web/app/page.tsx`
- Create: `apps/web/app/providers.tsx`
- Create: `apps/web/src/app/styles/globals.css`
- Create: `apps/web/src/app/config.ts`
- Create: `apps/web/src/shared/api/graphql-client.ts`
- Create: `apps/web/src/shared/config/index.ts`
- Create: `apps/web/src/shared/lib/index.ts`

### GraphQL codegen pipeline (Task 10)
- Create: `tools/codegen/codegen.ts`
- Create: `tools/codegen/package.json`

### Docker (Task 11)
- Create: `docker/docker-compose.yml`
- Create: `docker/api.Dockerfile`
- Create: `docker/web.Dockerfile`

### GitLab CI/CD (Task 12)
- Create: `.gitlab-ci.yml`

### README (Task 13)
- Create: `README.md`

---

## Task 1: Root Monorepo Scaffolding

**Files:**
- Create: `package.json`
- Create: `pnpm-workspace.yaml`
- Create: `nx.json`
- Create: `tsconfig.base.json`
- Create: `.gitignore`
- Create: `.env.example`
- Create: `.npmrc`

- [ ] **Step 1: Initialize pnpm and create package.json**

```bash
cd /home/nolood/general/rnd/monorepo-template
pnpm init
```

Then replace the generated `package.json` with:

```json
{
  "name": "monorepo-template",
  "version": "0.0.0",
  "private": true,
  "scripts": {
    "dev": "nx run-many --target=serve --all",
    "build": "nx run-many --target=build --all",
    "test": "nx run-many --target=test --all",
    "lint": "nx run-many --target=lint --all",
    "codegen": "nx run-many --target=codegen --projects=api,web",
    "prepare": "lefthook install"
  },
  "devDependencies": {
    "nx": "^20.0.0",
    "@nx/js": "^20.0.0",
    "@nx/workspace": "^20.0.0",
    "typescript": "^5.5.0",
    "lefthook": "^1.7.0",
    "@commitlint/cli": "^19.0.0",
    "@commitlint/config-conventional": "^19.0.0",
    "eslint": "^9.0.0",
    "prettier": "^3.3.0",
    "lint-staged": "^15.0.0"
  }
}
```

- [ ] **Step 2: Create pnpm-workspace.yaml**

```yaml
packages:
  - "apps/*"
  - "libs/*"
  - "tools/*"
```

- [ ] **Step 3: Create nx.json**

```json
{
  "$schema": "./node_modules/nx/schemas/nx-schema.json",
  "defaultBase": "main",
  "namedInputs": {
    "default": ["{projectRoot}/**/*", "sharedGlobals"],
    "sharedGlobals": ["{workspaceRoot}/tsconfig.base.json"],
    "production": ["default", "!{projectRoot}/**/*.spec.ts", "!{projectRoot}/**/*.test.ts"]
  },
  "targetDefaults": {
    "build": {
      "dependsOn": ["^build"],
      "inputs": ["production", "^production"],
      "cache": true
    },
    "test": {
      "inputs": ["default", "^production"],
      "cache": true
    },
    "lint": {
      "inputs": ["default"],
      "cache": true
    },
    "codegen": {
      "dependsOn": ["^validate"],
      "cache": true
    }
  },
  "plugins": []
}
```

Note: The `./tools/nx-go` plugin will be added to `nx.json` plugins array in Task 4 after the executor is created.
```

- [ ] **Step 4: Create tsconfig.base.json**

```json
{
  "compileOnSave": false,
  "compilerOptions": {
    "rootDir": ".",
    "sourceMap": true,
    "declaration": false,
    "moduleResolution": "bundler",
    "emitDecoratorMetadata": true,
    "experimentalDecorators": true,
    "importHelpers": true,
    "target": "es2022",
    "module": "esnext",
    "lib": ["es2022", "dom"],
    "skipLibCheck": true,
    "skipDefaultLibCheck": true,
    "baseUrl": ".",
    "paths": {}
  },
  "exclude": ["node_modules", "tmp"]
}
```

- [ ] **Step 5: Create .gitignore**

```
# Dependencies
node_modules/
.pnpm-store/

# Build output
dist/
.next/
out/

# Nx
.nx/

# Go
apps/api/tmp/
apps/api/server

# Generated
apps/api/internal/graph/generated.go
apps/api/internal/graph/model/models_gen.go
apps/web/src/shared/api/generated/

# Environment
.env
.env.local
.env.*.local

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Coverage
coverage/
*.out

# Docker
docker/postgres-data/

# Misc
schema.json
```

- [ ] **Step 6: Create .env.example**

```bash
# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=app
POSTGRES_PASSWORD=secret
POSTGRES_DB=monorepo_dev
POSTGRES_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Go API
API_PORT=8080
API_LOG_LEVEL=debug

# Auth
JWT_SECRET=change-me-in-production

# Next.js (client-side)
NEXT_PUBLIC_API_URL=http://localhost:8080/graphql
NEXT_PUBLIC_APP_NAME=MonorepoApp
```

- [ ] **Step 7: Create .npmrc**

```
auto-install-peers=true
strict-peer-dependencies=false
```

- [ ] **Step 8: Install dependencies and commit**

```bash
pnpm install
git add -A
git commit -m "chore: initialize monorepo root with Nx, pnpm workspace, base configs"
```

---

## Task 2: Linting, Hooks, and Conventional Commits

**Files:**
- Create: `.eslintrc.json`
- Create: `.prettierrc`
- Create: `.prettierignore`
- Create: `commitlint.config.js`
- Create: `.lefthook.yml`
- Create: `.lintstagedrc.json`

- [ ] **Step 1: Create .eslintrc.json**

```json
{
  "root": true,
  "ignorePatterns": ["**/*"],
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
                "sourceTag": "*",
                "onlyDependOnLibsWithTags": ["*"]
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

- [ ] **Step 2: Create .prettierrc**

```json
{
  "singleQuote": true,
  "semi": true,
  "trailingComma": "all",
  "printWidth": 100,
  "tabWidth": 2,
  "arrowParens": "always",
  "endOfLine": "lf"
}
```

- [ ] **Step 3: Create .prettierignore**

```
dist
.next
node_modules
coverage
apps/api
*.go
pnpm-lock.yaml
apps/web/src/shared/api/generated
```

- [ ] **Step 4: Create commitlint.config.js**

```js
module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'scope-enum': [
      2,
      'always',
      ['api', 'web', 'graphql', 'codegen', 'nx-go', 'docker', 'ci', 'deps'],
    ],
    'scope-empty': [1, 'never'],
  },
};
```

- [ ] **Step 5: Create .lefthook.yml**

```yaml
pre-commit:
  parallel: true
  commands:
    lint-staged:
      run: pnpm exec lint-staged
    go-lint:
      glob: "*.go"
      run: cd apps/api && golangci-lint run --fix
    go-test:
      glob: "*.go"
      run: cd apps/api && go test -short ./...

commit-msg:
  commands:
    commitlint:
      run: pnpm exec commitlint --edit {1}
```

- [ ] **Step 6: Create .lintstagedrc.json**

```json
{
  "*.{ts,tsx,js,jsx}": ["eslint --fix", "prettier --write"],
  "*.{json,yaml,yml,md,graphql}": ["prettier --write"]
}
```

- [ ] **Step 7: Install additional dev dependencies and commit**

```bash
pnpm add -Dw @nx/eslint-plugin eslint-plugin-boundaries
lefthook install
git add -A
git commit -m "chore: add linting, formatting, lefthook hooks, and commitlint"
```

---

## Task 3: GraphQL Schema

**Files:**
- Create: `libs/graphql/schema/common.graphql`
- Create: `libs/graphql/schema/user.graphql`
- Create: `libs/graphql/schema/schema.graphql`
- Create: `libs/graphql/project.json`

- [ ] **Step 1: Create libs/graphql/schema/common.graphql**

```graphql
scalar DateTime
scalar UUID

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}

input PaginationInput {
  first: Int
  after: String
  last: Int
  before: String
}

type ValidationError {
  field: String!
  message: String!
}

type AuthError {
  message: String!
}

type NotFoundError {
  message: String!
  entityType: String!
  id: String!
}
```

- [ ] **Step 2: Create libs/graphql/schema/user.graphql**

```graphql
type User {
  id: UUID!
  email: String!
  name: String!
  createdAt: DateTime!
  updatedAt: DateTime!
}

type UserEdge {
  cursor: String!
  node: User!
}

type UserConnection {
  edges: [UserEdge!]!
  pageInfo: PageInfo!
  totalCount: Int!
}

type CreateUserSuccess {
  user: User!
}

union CreateUserResult = CreateUserSuccess | ValidationError | AuthError

input CreateUserInput {
  email: String!
  name: String!
  password: String!
}

type UpdateUserSuccess {
  user: User!
}

union UpdateUserResult = UpdateUserSuccess | ValidationError | NotFoundError

input UpdateUserInput {
  name: String
  email: String
}
```

- [ ] **Step 3: Create libs/graphql/schema/schema.graphql**

```graphql
type Query {
  user(id: UUID!): User
  users(pagination: PaginationInput): UserConnection!
}

type Mutation {
  createUser(input: CreateUserInput!): CreateUserResult!
  updateUser(id: UUID!, input: UpdateUserInput!): UpdateUserResult!
  deleteUser(id: UUID!): Boolean!
}
```

- [ ] **Step 4: Create libs/graphql/project.json**

```json
{
  "name": "graphql",
  "$schema": "../../node_modules/nx/schemas/project-schema.json",
  "sourceRoot": "libs/graphql/schema",
  "projectType": "library",
  "targets": {
    "validate": {
      "executor": "nx:run-commands",
      "options": {
        "command": "npx graphql-inspector validate libs/graphql/schema/*.graphql"
      }
    }
  }
}
```

- [ ] **Step 5: Install graphql schema tooling and commit**

```bash
pnpm add -Dw graphql @graphql-inspector/cli
git add -A
git commit -m "feat(graphql): add shared GraphQL schema with User type and common types"
```

---

## Task 4: Custom Nx Go Executor

**Files:**
- Create: `tools/nx-go/package.json`
- Create: `tools/nx-go/tsconfig.json`
- Create: `tools/nx-go/executors.json`
- Create: `tools/nx-go/src/executors/build/executor.ts`
- Create: `tools/nx-go/src/executors/build/schema.json`
- Create: `tools/nx-go/src/executors/serve/executor.ts`
- Create: `tools/nx-go/src/executors/serve/schema.json`
- Create: `tools/nx-go/src/executors/test/executor.ts`
- Create: `tools/nx-go/src/executors/test/schema.json`
- Create: `tools/nx-go/src/executors/lint/executor.ts`
- Create: `tools/nx-go/src/executors/lint/schema.json`

- [ ] **Step 1: Create tools/nx-go/package.json**

```json
{
  "name": "nx-go",
  "version": "0.0.0",
  "private": true,
  "executors": "./executors.json"
}
```

- [ ] **Step 2: Create tools/nx-go/tsconfig.json**

```json
{
  "extends": "../../tsconfig.base.json",
  "compilerOptions": {
    "module": "commonjs",
    "outDir": "../../dist/tools/nx-go",
    "declaration": true,
    "types": ["node"]
  },
  "include": ["src/**/*.ts"]
}
```

- [ ] **Step 3: Create tools/nx-go/executors.json**

```json
{
  "executors": {
    "build": {
      "implementation": "./src/executors/build/executor",
      "schema": "./src/executors/build/schema.json",
      "description": "Build a Go application"
    },
    "serve": {
      "implementation": "./src/executors/serve/executor",
      "schema": "./src/executors/serve/schema.json",
      "description": "Serve a Go application with hot reload"
    },
    "test": {
      "implementation": "./src/executors/test/executor",
      "schema": "./src/executors/test/schema.json",
      "description": "Run Go tests"
    },
    "lint": {
      "implementation": "./src/executors/lint/executor",
      "schema": "./src/executors/lint/schema.json",
      "description": "Run golangci-lint"
    }
  }
}
```

- [ ] **Step 4: Create build executor**

Create `tools/nx-go/src/executors/build/schema.json`:

```json
{
  "type": "object",
  "properties": {
    "outputPath": { "type": "string", "description": "Output binary path" },
    "main": { "type": "string", "description": "Main package path (relative to project root)", "default": "cmd/server" }
  },
  "required": ["outputPath"]
}
```

Create `tools/nx-go/src/executors/build/executor.ts`:

```typescript
import { ExecutorContext } from '@nx/devkit';
import { execSync } from 'child_process';
import * as path from 'path';

interface BuildExecutorOptions {
  outputPath: string;
  main: string;
}

export default async function runExecutor(
  options: BuildExecutorOptions,
  context: ExecutorContext,
): Promise<{ success: boolean }> {
  const projectRoot = context.projectsConfigurations!.projects[context.projectName!].root;
  const cwd = path.join(context.root, projectRoot);
  const outputPath = path.join(context.root, options.outputPath);

  console.log(`Building Go application in ${cwd}...`);

  try {
    execSync(`go build -o ${outputPath} ./${options.main}`, {
      cwd,
      stdio: 'inherit',
      env: { ...process.env },
    });
    return { success: true };
  } catch {
    return { success: false };
  }
}
```

- [ ] **Step 5: Create serve executor**

Create `tools/nx-go/src/executors/serve/schema.json`:

```json
{
  "type": "object",
  "properties": {
    "port": { "type": "number", "description": "Server port", "default": 8080 },
    "configPath": { "type": "string", "description": "Path to air.toml config" }
  }
}
```

Create `tools/nx-go/src/executors/serve/executor.ts`:

```typescript
import { ExecutorContext } from '@nx/devkit';
import { spawn } from 'child_process';
import * as path from 'path';

interface ServeExecutorOptions {
  port: number;
  configPath?: string;
}

export default async function runExecutor(
  options: ServeExecutorOptions,
  context: ExecutorContext,
): Promise<{ success: boolean }> {
  const projectRoot = context.projectsConfigurations!.projects[context.projectName!].root;
  const cwd = path.join(context.root, projectRoot);

  console.log(`Serving Go application on port ${options.port}...`);

  const airConfig = options.configPath || 'air.toml';

  return new Promise((resolve) => {
    const child = spawn('air', ['-c', airConfig], {
      cwd,
      stdio: 'inherit',
      env: { ...process.env, API_PORT: String(options.port) },
    });

    // Forward signals for clean shutdown
    const signalHandler = (signal: NodeJS.Signals) => {
      child.kill(signal);
    };
    process.on('SIGINT', signalHandler);
    process.on('SIGTERM', signalHandler);

    child.on('close', (code) => {
      process.removeListener('SIGINT', signalHandler);
      process.removeListener('SIGTERM', signalHandler);
      resolve({ success: code === 0 || code === null });
    });
  });
}
```

- [ ] **Step 6: Create test executor**

Create `tools/nx-go/src/executors/test/schema.json`:

```json
{
  "type": "object",
  "properties": {
    "coverage": { "type": "boolean", "description": "Enable coverage", "default": false },
    "short": { "type": "boolean", "description": "Run short tests only", "default": false },
    "packages": { "type": "array", "items": { "type": "string" }, "description": "Specific packages to test" }
  }
}
```

Create `tools/nx-go/src/executors/test/executor.ts`:

```typescript
import { ExecutorContext } from '@nx/devkit';
import { execSync } from 'child_process';
import * as path from 'path';

interface TestExecutorOptions {
  coverage: boolean;
  short: boolean;
  packages?: string[];
}

export default async function runExecutor(
  options: TestExecutorOptions,
  context: ExecutorContext,
): Promise<{ success: boolean }> {
  const projectRoot = context.projectsConfigurations!.projects[context.projectName!].root;
  const cwd = path.join(context.root, projectRoot);

  const args: string[] = ['go', 'test'];

  if (options.short) args.push('-short');
  if (options.coverage) args.push('-coverprofile=coverage.out');

  const pkgs = options.packages?.length ? options.packages.join(' ') : './...';
  args.push(pkgs);

  console.log(`Running Go tests in ${cwd}...`);

  try {
    execSync(args.join(' '), {
      cwd,
      stdio: 'inherit',
      env: { ...process.env },
    });
    return { success: true };
  } catch {
    return { success: false };
  }
}
```

- [ ] **Step 7: Create lint executor**

Create `tools/nx-go/src/executors/lint/schema.json`:

```json
{
  "type": "object",
  "properties": {
    "fix": { "type": "boolean", "description": "Auto-fix issues", "default": false },
    "config": { "type": "string", "description": "Path to golangci-lint config" }
  }
}
```

Create `tools/nx-go/src/executors/lint/executor.ts`:

```typescript
import { ExecutorContext } from '@nx/devkit';
import { execSync } from 'child_process';
import * as path from 'path';

interface LintExecutorOptions {
  fix: boolean;
  config?: string;
}

export default async function runExecutor(
  options: LintExecutorOptions,
  context: ExecutorContext,
): Promise<{ success: boolean }> {
  const projectRoot = context.projectsConfigurations!.projects[context.projectName!].root;
  const cwd = path.join(context.root, projectRoot);

  const args: string[] = ['golangci-lint', 'run'];
  if (options.fix) args.push('--fix');
  if (options.config) args.push(`--config=${options.config}`);

  console.log(`Linting Go code in ${cwd}...`);

  try {
    execSync(args.join(' '), {
      cwd,
      stdio: 'inherit',
      env: { ...process.env },
    });
    return { success: true };
  } catch {
    return { success: false };
  }
}
```

- [ ] **Step 8: Create tools/nx-go/src/index.ts**

```typescript
export { default as buildExecutor } from './executors/build/executor';
export { default as serveExecutor } from './executors/serve/executor';
export { default as testExecutor } from './executors/test/executor';
export { default as lintExecutor } from './executors/lint/executor';
```

- [ ] **Step 9: Register nx-go plugin in nx.json**

Update `nx.json` plugins array from `[]` to `["./tools/nx-go"]`.

- [ ] **Step 10: Commit**

```bash
git add -A
git commit -m "feat(nx-go): add custom Nx executor for Go build, serve, test, lint"
```

---

## Task 5: Go API — Config Module

**Files:**
- Create: `apps/api/go.mod`
- Create: `apps/api/config/config.yml`
- Create: `apps/api/internal/config/config.go`
- Create: `apps/api/internal/config/config_test.go`
- Create: `apps/api/.golangci.yml`

- [ ] **Step 1: Initialize Go module**

```bash
cd apps/api
go mod init monorepo-template/apps/api
```

- [ ] **Step 2: Create apps/api/.golangci.yml**

```yaml
run:
  timeout: 5m

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gosec
    - gofmt
    - goimports
    - misspell
    - unconvert

linters-settings:
  goimports:
    local-prefixes: monorepo-template
```

- [ ] **Step 3: Create apps/api/config/config.yml**

```yaml
server:
  port: 8080
  read_timeout: 10s
  write_timeout: 30s
  shutdown_timeout: 5s

log:
  level: info
  format: json

pagination:
  default_page_size: 20
  max_page_size: 100
```

- [ ] **Step 4: Write the failing test for config loading**

Create `apps/api/internal/config/config_test.go`:

```go
package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/config"
)

func TestLoad_DefaultsFromFile(t *testing.T) {
	cfg, err := config.Load("../../config/config.yml")
	require.NoError(t, err)

	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, 20, cfg.Pagination.DefaultPageSize)
	assert.Equal(t, 100, cfg.Pagination.MaxPageSize)
}

func TestLoad_EnvOverridesConfig(t *testing.T) {
	os.Setenv("POSTGRES_HOST", "testhost")
	os.Setenv("POSTGRES_PORT", "5433")
	os.Setenv("JWT_SECRET", "test-secret")
	defer func() {
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg, err := config.Load("../../config/config.yml")
	require.NoError(t, err)

	assert.Equal(t, "testhost", cfg.Postgres.Host)
	assert.Equal(t, 5433, cfg.Postgres.Port)
	assert.Equal(t, "test-secret", cfg.Auth.JWTSecret)
}
```

- [ ] **Step 5: Run test to verify it fails**

```bash
cd apps/api
go test ./internal/config/... -v
```

Expected: compilation error — `config` package doesn't exist yet.

- [ ] **Step 6: Write config implementation**

Create `apps/api/internal/config/config.go`:

```go
package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig
	Log        LogConfig
	Postgres   PostgresConfig
	Redis      RedisConfig
	Auth       AuthConfig
	Pagination PaginationConfig
}

type ServerConfig struct {
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type AuthConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

type PaginationConfig struct {
	DefaultPageSize int `mapstructure:"default_page_size"`
	MaxPageSize     int `mapstructure:"max_page_size"`
}

func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Load stable values from config file
	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// Bind sensitive values from environment
	v.BindEnv("postgres.host", "POSTGRES_HOST")
	v.BindEnv("postgres.port", "POSTGRES_PORT")
	v.BindEnv("postgres.user", "POSTGRES_USER")
	v.BindEnv("postgres.password", "POSTGRES_PASSWORD")
	v.BindEnv("postgres.db", "POSTGRES_DB")
	v.BindEnv("postgres.sslmode", "POSTGRES_SSLMODE")
	v.BindEnv("redis.host", "REDIS_HOST")
	v.BindEnv("redis.port", "REDIS_PORT")
	v.BindEnv("redis.password", "REDIS_PASSWORD")
	v.BindEnv("auth.jwt_secret", "JWT_SECRET")
	v.BindEnv("server.port", "API_PORT")
	v.BindEnv("log.level", "API_LOG_LEVEL")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
```

- [ ] **Step 7: Install Go deps and run tests**

```bash
cd apps/api
go mod tidy
go test ./internal/config/... -v
```

Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add -A
git commit -m "feat(api): add viper-based config module with env override for secrets"
```

---

## Task 6: Go API — Health Handler

**Files:**
- Create: `apps/api/internal/handler/health.go`
- Create: `apps/api/internal/handler/health_test.go`

- [ ] **Step 1: Write the failing test**

Create `apps/api/internal/handler/health_test.go`:

```go
package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"monorepo-template/apps/api/internal/handler"
)

func TestHealthz_ReturnsOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.Healthz()(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func TestReadyz_HealthyDeps(t *testing.T) {
	checker := &mockChecker{healthy: true}
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.Readyz(checker)(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func TestReadyz_UnhealthyDeps(t *testing.T) {
	checker := &mockChecker{healthy: false}
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.Readyz(checker)(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

type mockChecker struct {
	healthy bool
}

func (m *mockChecker) Ping() error {
	if !m.healthy {
		return assert.AnError
	}
	return nil
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd apps/api && go test ./internal/handler/... -v
```

Expected: compilation error.

- [ ] **Step 3: Write implementation**

Create `apps/api/internal/handler/health.go`:

```go
package handler

import (
	"encoding/json"
	"net/http"
)

type HealthChecker interface {
	Ping() error
}

func Healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func Readyz(checkers ...HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		for _, c := range checkers {
			if err := c.Ping(); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				json.NewEncoder(w).Encode(map[string]string{
					"status": "unavailable",
					"error":  err.Error(),
				})
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}
```

- [ ] **Step 4: Run tests**

```bash
cd apps/api && go test ./internal/handler/... -v
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat(api): add health check handlers /healthz and /readyz"
```

---

## Task 7: Go API — Middleware

**Files:**
- Create: `apps/api/internal/middleware/logging.go`
- Create: `apps/api/internal/middleware/cors.go`
- Create: `apps/api/internal/middleware/auth.go`

- [ ] **Step 1: Create logging middleware**

Create `apps/api/internal/middleware/logging.go`:

```go
package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Logging(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(wrapped, r)

			logger.Info("request",
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

- [ ] **Step 2: Create CORS middleware**

Create `apps/api/internal/middleware/cors.go`:

```go
package middleware

import (
	"net/http"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}
}

func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			for _, allowed := range cfg.AllowedOrigins {
				if allowed == "*" || allowed == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
```

- [ ] **Step 3: Create auth middleware (placeholder with JWT structure)**

Create `apps/api/internal/middleware/auth.go`:

```go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(jwtSecret string, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// No auth header — pass through, let resolvers decide if auth is required
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == authHeader {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// TODO: implement JWT validation using jwtSecret
			// For now, pass the token as user ID placeholder
			_ = jwtSecret
			ctx := context.WithValue(r.Context(), UserIDKey, tokenStr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
```

- [ ] **Step 4: Write CORS middleware test**

Create `apps/api/internal/middleware/cors_test.go`:

```go
package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"monorepo-template/apps/api/internal/middleware"
)

func TestCORS_AllowedOrigin(t *testing.T) {
	handler := middleware.CORS(middleware.DefaultCORSConfig())(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, "http://localhost:3000", rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_PreflightReturnsNoContent(t *testing.T) {
	handler := middleware.CORS(middleware.DefaultCORSConfig())(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestCORS_DisallowedOrigin(t *testing.T) {
	handler := middleware.CORS(middleware.DefaultCORSConfig())(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://evil.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
}
```

- [ ] **Step 5: Write auth middleware test**

Create `apps/api/internal/middleware/auth_test.go`:

```go
package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/middleware"
)

func TestAuth_NoHeader_PassesThrough(t *testing.T) {
	logger := zap.NewNop()
	handler := middleware.Auth("secret", logger)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := middleware.GetUserID(r.Context())
			assert.False(t, ok)
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuth_InvalidFormat_Returns401(t *testing.T) {
	logger := zap.NewNop()
	handler := middleware.Auth("secret", logger)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "InvalidFormat token123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuth_ValidBearer_SetsContext(t *testing.T) {
	logger := zap.NewNop()
	handler := middleware.Auth("secret", logger)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := middleware.GetUserID(r.Context())
			assert.True(t, ok)
			assert.Equal(t, "my-token", userID)
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer my-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
```

- [ ] **Step 6: Run middleware tests**

```bash
cd apps/api && go mod tidy && go test ./internal/middleware/... -v
```

Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "feat(api): add logging, CORS, and auth middleware with tests"
```

---

## Task 8: Go API — Repository, Service, Graph, and Main

**Files:**
- Create: `apps/api/tools.go`
- Create: `apps/api/internal/repository/postgres/postgres.go`
- Create: `apps/api/internal/repository/postgres/user_repo.go`
- Create: `apps/api/internal/repository/postgres/migrations/00001_init.sql`
- Create: `apps/api/internal/repository/redis/cache.go`
- Create: `apps/api/internal/service/user_service.go`
- Create: `apps/api/internal/service/user_service_test.go`
- Create: `apps/api/internal/graph/resolver.go`
- Create: `apps/api/gqlgen.yml`
- Create: `apps/api/air.toml`
- Create: `apps/api/project.json`
- Create: `apps/api/cmd/server/main.go`

- [ ] **Step 1: Create tools.go to pin Go tool dependencies**

Create `apps/api/tools.go`:

```go
//go:build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/pressly/goose/v3/cmd/goose"
)
```

This ensures gqlgen and goose are tracked in `go.mod` so `go run github.com/99designs/gqlgen generate` works.

- [ ] **Step 2: Create PostgreSQL connection module**

Create `apps/api/internal/repository/postgres/postgres.go`:

```go
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/config"
)

type DB struct {
	Pool   *pgxpool.Pool
	logger *zap.Logger
}

func New(cfg config.PostgresConfig, logger *zap.Logger) (*DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB, cfg.SSLMode,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
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

func (db *DB) Ping() error {
	return db.Pool.Ping(context.Background())
}

func (db *DB) Close() {
	db.Pool.Close()
}
```

- [ ] **Step 3: Create user repository stub**

Create `apps/api/internal/repository/postgres/user_repo.go`:

```go
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/service"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*service.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *UserRepo) List(ctx context.Context, first *int, after *string) ([]*service.User, int, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

func (r *UserRepo) Create(ctx context.Context, input service.CreateUserInput) (*service.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *UserRepo) Update(ctx context.Context, id string, input service.UpdateUserInput) (*service.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}
```

- [ ] **Step 4: Create initial migration**

Create `apps/api/internal/repository/postgres/migrations/00001_init.sql`:

```sql
-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- +goose Down
DROP TABLE IF EXISTS users;
```

- [ ] **Step 5: Create Redis connection and cache module**

Create `apps/api/internal/repository/redis/cache.go`:

```go
package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/config"
)

type Client struct {
	RDB    *redis.Client
	logger *zap.Logger
}

func New(cfg config.RedisConfig, logger *zap.Logger) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	logger.Info("connected to Redis",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
	)

	return &Client{RDB: rdb, logger: logger}, nil
}

func (c *Client) Ping() error {
	return c.RDB.Ping(context.Background()).Err()
}

func (c *Client) Close() error {
	return c.RDB.Close()
}
```

- [ ] **Step 6: Create user service interface**

Create `apps/api/internal/service/user_service.go`:

```go
package service

import "context"

type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt string
	UpdatedAt string
}

type CreateUserInput struct {
	Email    string
	Name     string
	Password string
}

type UpdateUserInput struct {
	Name  *string
	Email *string
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context, first *int, after *string) ([]*User, int, error)
	Create(ctx context.Context, input CreateUserInput) (*User, error)
	Update(ctx context.Context, id string, input UpdateUserInput) (*User, error)
	Delete(ctx context.Context, id string) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context, first *int, after *string) ([]*User, int, error) {
	return s.repo.List(ctx, first, after)
}

func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	return s.repo.Create(ctx, input)
}

func (s *UserService) Update(ctx context.Context, id string, input UpdateUserInput) (*User, error) {
	return s.repo.Update(ctx, id, input)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
```

- [ ] **Step 7: Write user service test**

Create `apps/api/internal/service/user_service_test.go`:

```go
package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/service"
)

type mockUserRepo struct {
	users map[string]*service.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*service.User)}
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*service.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepo) List(ctx context.Context, first *int, after *string) ([]*service.User, int, error) {
	var result []*service.User
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, len(result), nil
}

func (m *mockUserRepo) Create(ctx context.Context, input service.CreateUserInput) (*service.User, error) {
	u := &service.User{ID: "test-id", Email: input.Email, Name: input.Name}
	m.users[u.ID] = u
	return u, nil
}

func (m *mockUserRepo) Update(ctx context.Context, id string, input service.UpdateUserInput) (*service.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	if input.Name != nil {
		u.Name = *input.Name
	}
	if input.Email != nil {
		u.Email = *input.Email
	}
	return u, nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

func TestUserService_Create(t *testing.T) {
	repo := newMockUserRepo()
	svc := service.NewUserService(repo)

	user, err := svc.Create(context.Background(), service.CreateUserInput{
		Email: "test@example.com",
		Name:  "Test User",
	})

	require.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)
}

func TestUserService_GetByID(t *testing.T) {
	repo := newMockUserRepo()
	svc := service.NewUserService(repo)

	created, _ := svc.Create(context.Background(), service.CreateUserInput{
		Email: "test@example.com",
		Name:  "Test",
	})

	found, err := svc.GetByID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
}
```

- [ ] **Step 8: Run user service tests**

```bash
cd apps/api && go test ./internal/service/... -v
```

Expected: PASS

- [ ] **Step 9: Create gqlgen.yml and resolver**

Create `apps/api/gqlgen.yml`:

```yaml
schema:
  - ../../libs/graphql/schema/*.graphql
exec:
  filename: internal/graph/generated.go
  package: graph
model:
  filename: internal/graph/model/models_gen.go
  package: model
resolver:
  layout: follow-schema
  dir: internal/graph
  package: graph
autobind: []
```

Create `apps/api/internal/graph/resolver.go`:

```go
package graph

import "monorepo-template/apps/api/internal/service"

// Resolver holds dependencies for GraphQL resolvers.
type Resolver struct {
	UserService *service.UserService
}
```

- [ ] **Step 10: Run gqlgen to generate resolver stubs**

```bash
cd apps/api
go get github.com/99designs/gqlgen@latest
go mod tidy
go run github.com/99designs/gqlgen generate
```

This generates `internal/graph/generated.go`, `internal/graph/model/models_gen.go`, and `internal/graph/schema.resolvers.go` with TODO stubs.

- [ ] **Step 11: Create air.toml**

Create `apps/api/air.toml`:

```toml
root = "."
tmp_dir = "tmp"

[build]
  bin = "./tmp/server"
  cmd = "go build -o ./tmp/server ./cmd/server"
  delay = 1000
  exclude_dir = ["tmp", "vendor", "node_modules"]
  exclude_regex = ["_test\\.go"]
  include_ext = ["go", "yml", "yaml"]
  kill_delay = "0s"
  send_interrupt = false
  stop_on_error = true

[log]
  time = false

[misc]
  clean_on_exit = true
```

- [ ] **Step 12: Create apps/api/project.json**

```json
{
  "name": "api",
  "$schema": "../../node_modules/nx/schemas/project-schema.json",
  "sourceRoot": "apps/api",
  "projectType": "application",
  "targets": {
    "build": {
      "executor": "nx-go:build",
      "options": {
        "outputPath": "dist/apps/api",
        "main": "cmd/server"
      }
    },
    "serve": {
      "executor": "nx-go:serve",
      "options": {
        "port": 8080
      }
    },
    "test": {
      "executor": "nx-go:test",
      "options": {
        "coverage": true
      }
    },
    "go-lint": {
      "executor": "nx-go:lint",
      "options": {
        "fix": false
      }
    },
    "codegen": {
      "executor": "nx:run-commands",
      "options": {
        "command": "go run github.com/99designs/gqlgen generate",
        "cwd": "apps/api"
      },
      "dependsOn": ["^validate"]
    }
  }
}
```

- [ ] **Step 13: Create main.go entrypoint**

Create `apps/api/cmd/server/main.go`:

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/config"
	"monorepo-template/apps/api/internal/graph"
	healthHandler "monorepo-template/apps/api/internal/handler"
	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/apps/api/internal/repository/postgres"
	redisRepo "monorepo-template/apps/api/internal/repository/redis"
	"monorepo-template/apps/api/internal/service"
)

func main() {
	// Load config
	cfg, err := config.Load("config/config.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Init logger
	var logger *zap.Logger
	if cfg.Log.Format == "json" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	// Connect to PostgreSQL
	db, err := postgres.New(cfg.Postgres, logger)
	if err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	}
	defer db.Close()

	// Connect to Redis
	rdb, err := redisRepo.New(cfg.Redis, logger)
	if err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}
	defer rdb.Close()

	// Init services
	userRepo := postgres.NewUserRepo(db.Pool)
	userService := service.NewUserService(userRepo)

	// GraphQL server
	resolver := &graph.Resolver{
		UserService: userService,
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logging(logger))
	r.Use(middleware.CORS(middleware.DefaultCORSConfig()))
	r.Use(middleware.Auth(cfg.Auth.JWTSecret, logger))

	// Health checks
	r.Get("/healthz", healthHandler.Healthz())
	r.Get("/readyz", healthHandler.Readyz(db, rdb))

	// GraphQL
	r.Handle("/graphql", srv)
	r.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	// Start server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("starting server", zap.Int("port", cfg.Server.Port))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server stopped")
}
```

- [ ] **Step 14: Install all Go dependencies and commit**

```bash
cd apps/api && go mod tidy
git add -A
git commit -m "feat(api): add Go API with chi router, gqlgen, health checks, middleware, config"
```

---

## Task 9: Next.js Frontend with FSD

**Files:**
- Create: `apps/web/package.json`
- Create: `apps/web/tsconfig.json`
- Create: `apps/web/next.config.js`
- Create: `apps/web/tailwind.config.ts`
- Create: `apps/web/postcss.config.js`
- Create: `apps/web/project.json`
- Create: `apps/web/vitest.config.ts`
- Create: `apps/web/.eslintrc.json`
- Create: `apps/web/app/layout.tsx`
- Create: `apps/web/app/page.tsx`
- Create: `apps/web/app/providers.tsx`
- Create: `apps/web/src/app/styles/globals.css`
- Create: `apps/web/src/app/config.ts`
- Create: `apps/web/src/shared/api/graphql-client.ts`
- Create: `apps/web/src/shared/config/index.ts`
- Create: `apps/web/src/shared/lib/index.ts`
- Create: `apps/web/src/pages/.gitkeep`
- Create: `apps/web/src/widgets/.gitkeep`
- Create: `apps/web/src/features/.gitkeep`
- Create: `apps/web/src/entities/.gitkeep`

- [ ] **Step 1: Create apps/web/package.json**

```json
{
  "name": "web",
  "version": "0.0.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "test": "vitest run",
    "test:watch": "vitest",
    "codegen": "graphql-codegen --config ../../tools/codegen/codegen.ts"
  },
  "dependencies": {
    "next": "^15.0.0",
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "@tanstack/react-query": "^5.0.0",
    "graphql-request": "^7.0.0",
    "graphql": "^16.0.0"
  },
  "devDependencies": {
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "typescript": "^5.5.0",
    "tailwindcss": "^3.4.0",
    "postcss": "^8.4.0",
    "autoprefixer": "^10.4.0",
    "vitest": "^2.0.0",
    "@vitejs/plugin-react": "^4.0.0",
    "@testing-library/react": "^16.0.0",
    "@testing-library/jest-dom": "^6.0.0",
    "jsdom": "^25.0.0",
    "eslint-plugin-boundaries": "^4.0.0",
    "@graphql-codegen/cli": "^5.0.0",
    "@graphql-codegen/typescript": "^4.0.0",
    "@graphql-codegen/typescript-operations": "^4.0.0",
    "@graphql-codegen/typescript-react-query": "^6.0.0",
    "@playwright/test": "^1.48.0"
  }
}
```

- [ ] **Step 2: Create apps/web/tsconfig.json**

```json
{
  "extends": "../../tsconfig.base.json",
  "compilerOptions": {
    "jsx": "preserve",
    "allowJs": true,
    "esModuleInterop": true,
    "allowSyntheticDefaultImports": true,
    "forceConsistentCasingInFileNames": true,
    "strict": true,
    "noEmit": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "incremental": true,
    "plugins": [{ "name": "next" }],
    "paths": {
      "@/*": ["./src/*"],
      "@app/*": ["./src/app/*"],
      "@pages/*": ["./src/pages/*"],
      "@widgets/*": ["./src/widgets/*"],
      "@features/*": ["./src/features/*"],
      "@entities/*": ["./src/entities/*"],
      "@shared/*": ["./src/shared/*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules", ".next", "e2e"]
}
```

- [ ] **Step 3: Create next.config.js**

```js
/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  transpilePackages: [],
  output: 'standalone',
};

module.exports = nextConfig;
```

- [ ] **Step 4: Create Tailwind and PostCSS configs**

Create `apps/web/tailwind.config.ts`:

```typescript
import type { Config } from 'tailwindcss';

const config: Config = {
  content: [
    './app/**/*.{ts,tsx}',
    './src/**/*.{ts,tsx}',
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};

export default config;
```

Create `apps/web/postcss.config.js`:

```js
module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
```

- [ ] **Step 5: Create project.json and vitest.config.ts**

Create `apps/web/project.json`:

```json
{
  "name": "web",
  "$schema": "../../node_modules/nx/schemas/project-schema.json",
  "sourceRoot": "apps/web/src",
  "projectType": "application",
  "targets": {
    "build": {
      "executor": "nx:run-commands",
      "options": {
        "command": "pnpm --filter web build"
      }
    },
    "serve": {
      "executor": "nx:run-commands",
      "options": {
        "command": "pnpm --filter web dev"
      }
    },
    "test": {
      "executor": "nx:run-commands",
      "options": {
        "command": "pnpm --filter web test"
      }
    },
    "lint": {
      "executor": "nx:run-commands",
      "options": {
        "command": "eslint apps/web --ext .ts,.tsx"
      }
    },
    "typecheck": {
      "executor": "nx:run-commands",
      "options": {
        "command": "pnpm --filter web exec tsc --noEmit"
      }
    },
    "codegen": {
      "executor": "nx:run-commands",
      "options": {
        "command": "pnpm --filter web codegen"
      },
      "dependsOn": ["^validate"]
    },
    "e2e": {
      "executor": "nx:run-commands",
      "options": {
        "command": "pnpm --filter web exec playwright test"
      }
    }
  }
}
```

Create `apps/web/vitest.config.ts`:

```typescript
import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: [],
    coverage: {
      provider: 'v8',
      thresholds: {
        statements: 70,
        branches: 70,
        functions: 70,
        lines: 70,
      },
    },
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
      '@app': resolve(__dirname, './src/app'),
      '@pages': resolve(__dirname, './src/pages'),
      '@widgets': resolve(__dirname, './src/widgets'),
      '@features': resolve(__dirname, './src/features'),
      '@entities': resolve(__dirname, './src/entities'),
      '@shared': resolve(__dirname, './src/shared'),
    },
  },
});
```

- [ ] **Step 6: Create .eslintrc.json with FSD boundaries**

Create `apps/web/.eslintrc.json`:

```json
{
  "extends": ["../../.eslintrc.json", "next/core-web-vitals"],
  "plugins": ["boundaries"],
  "settings": {
    "boundaries/elements": [
      { "type": "app", "pattern": "src/app/*" },
      { "type": "pages", "pattern": "src/pages/*" },
      { "type": "widgets", "pattern": "src/widgets/*" },
      { "type": "features", "pattern": "src/features/*" },
      { "type": "entities", "pattern": "src/entities/*" },
      { "type": "shared", "pattern": "src/shared/*" }
    ],
    "boundaries/ignore": ["**/*.test.*", "**/__tests__/**"]
  },
  "rules": {
    "boundaries/element-types": [
      "error",
      {
        "default": "disallow",
        "rules": [
          { "from": "app", "allow": ["pages", "widgets", "features", "entities", "shared"] },
          { "from": "pages", "allow": ["widgets", "features", "entities", "shared"] },
          { "from": "widgets", "allow": ["features", "entities", "shared"] },
          { "from": "features", "allow": ["entities", "shared"] },
          { "from": "entities", "allow": ["shared"] },
          { "from": "shared", "allow": ["shared"] }
        ]
      }
    ]
  }
}
```

- [ ] **Step 7: Create App Router files**

Create `apps/web/src/app/styles/globals.css`:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

Create `apps/web/src/app/config.ts`:

```typescript
// App-level config re-exports from shared for convenience
export { appConfig } from '@shared/config';
```

Create `apps/web/app/providers.tsx`:

```tsx
'use client';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useState, type ReactNode } from 'react';

export function Providers({ children }: { children: ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 60 * 1000,
            refetchOnWindowFocus: false,
          },
        },
      }),
  );

  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
}
```

Create `apps/web/app/layout.tsx`:

```tsx
import type { Metadata } from 'next';
import { Providers } from './providers';
import '@app/styles/globals.css';

export const metadata: Metadata = {
  title: 'Monorepo Template',
  description: 'Production-ready monorepo template',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
```

Create `apps/web/app/page.tsx`:

```tsx
export default function HomePage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center">
      <h1 className="text-4xl font-bold">Monorepo Template</h1>
      <p className="mt-4 text-lg text-gray-600">Go + Next.js + GraphQL</p>
    </main>
  );
}
```

- [ ] **Step 8: Create shared layer**

Create `apps/web/src/shared/config/index.ts`:

```typescript
export const appConfig = {
  apiUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/graphql',
  appName: process.env.NEXT_PUBLIC_APP_NAME || 'MonorepoApp',
} as const;
```

Create `apps/web/src/shared/api/graphql-client.ts`:

```typescript
import { GraphQLClient } from 'graphql-request';
import { appConfig } from '@shared/config';

export const graphqlClient = new GraphQLClient(appConfig.apiUrl, {
  headers: {},
});

export function setAuthToken(token: string) {
  graphqlClient.setHeader('Authorization', `Bearer ${token}`);
}
```

Create `apps/web/src/shared/lib/index.ts`:

```typescript
// Shared utilities — add as needed
```

- [ ] **Step 9: Create FSD layer placeholder directories**

```bash
mkdir -p apps/web/src/pages apps/web/src/widgets apps/web/src/features apps/web/src/entities apps/web/e2e
touch apps/web/src/pages/.gitkeep
touch apps/web/src/widgets/.gitkeep
touch apps/web/src/features/.gitkeep
touch apps/web/src/entities/.gitkeep
```

- [ ] **Step 10: Create Playwright config**

Create `apps/web/e2e/playwright.config.ts`:

```typescript
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: '.',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
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
  webServer: {
    command: 'pnpm dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
    cwd: '..',
  },
});
```

- [ ] **Step 11: Install dependencies and Playwright browsers**

```bash
cd /home/nolood/general/rnd/monorepo-template
pnpm install
pnpm --filter web exec playwright install chromium
```

- [ ] **Step 12: Commit**

```bash
git add -A
git commit -m "feat(web): add Next.js app with FSD architecture, TanStack Query, Tailwind CSS"
```

---

## Task 10: GraphQL Codegen Pipeline

**Files:**
- Create: `tools/codegen/codegen.ts`
- Create: `tools/codegen/package.json`

- [ ] **Step 1: Create tools/codegen/package.json**

```json
{
  "name": "codegen",
  "version": "0.0.0",
  "private": true,
  "devDependencies": {
    "@graphql-codegen/cli": "^5.0.0",
    "@graphql-codegen/typescript": "^4.0.0",
    "@graphql-codegen/typescript-operations": "^4.0.0",
    "@graphql-codegen/typescript-react-query": "^6.0.0",
    "graphql": "^16.0.0"
  }
}
```

- [ ] **Step 2: Create tools/codegen/codegen.ts**

```typescript
import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  schema: '../../libs/graphql/schema/**/*.graphql',
  documents: [
    '../../apps/web/src/features/**/api/**/*.graphql',
    '../../apps/web/src/entities/**/api/**/*.graphql',
  ],
  ignoreNoDocuments: true,
  generates: {
    '../../apps/web/src/shared/api/generated/types.ts': {
      plugins: ['typescript'],
      config: {
        scalars: {
          DateTime: 'string',
          UUID: 'string',
        },
      },
    },
    '../../apps/web/src/shared/api/generated/hooks.ts': {
      preset: 'import-types',
      presetConfig: {
        typesPath: './types',
      },
      plugins: ['typescript-operations', 'typescript-react-query'],
      config: {
        fetcher: {
          func: '../graphql-client#graphqlClient.request',
          isReactHook: false,
        },
        reactQueryVersion: 5,
        exposeQueryKeys: true,
        scalars: {
          DateTime: 'string',
          UUID: 'string',
        },
      },
    },
  },
};

export default config;
```

- [ ] **Step 3: Install, run codegen, and commit**

```bash
cd /home/nolood/general/rnd/monorepo-template
pnpm install
mkdir -p apps/web/src/shared/api/generated
pnpm --filter codegen exec graphql-codegen --config tools/codegen/codegen.ts
git add -A
git commit -m "feat(codegen): add graphql-codegen pipeline with TanStack Query hooks generation"
```

---

## Task 11: Docker & Docker Compose

**Files:**
- Create: `docker/docker-compose.yml`
- Create: `docker/api.Dockerfile`
- Create: `docker/web.Dockerfile`

- [ ] **Step 1: Create docker/api.Dockerfile**

```dockerfile
# ---- Dev stage ----
FROM golang:1.23-alpine AS dev

RUN go install github.com/air-verse/air@latest
WORKDIR /app
COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download
COPY apps/api/ .
CMD ["air", "-c", "air.toml"]

# ---- Build stage ----
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download
COPY apps/api/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# ---- Production stage ----
FROM alpine:3.20 AS prod

RUN apk --no-cache add ca-certificates
COPY --from=builder /server /server
COPY apps/api/config/config.yml /config/config.yml
EXPOSE 8080
CMD ["/server"]
```

- [ ] **Step 2: Create docker/web.Dockerfile**

```dockerfile
# ---- Dev stage ----
FROM node:22-alpine AS dev

RUN corepack enable && corepack prepare pnpm@latest --activate
WORKDIR /app
COPY package.json pnpm-workspace.yaml pnpm-lock.yaml .npmrc ./
COPY apps/web/package.json apps/web/
COPY libs/ libs/
COPY tools/ tools/
RUN pnpm install --frozen-lockfile
COPY apps/web/ apps/web/
WORKDIR /app/apps/web
CMD ["pnpm", "dev"]

# ---- Build stage ----
FROM node:22-alpine AS builder

RUN corepack enable && corepack prepare pnpm@latest --activate
WORKDIR /app
COPY package.json pnpm-workspace.yaml pnpm-lock.yaml .npmrc ./
COPY apps/web/package.json apps/web/
COPY libs/ libs/
COPY tools/ tools/
RUN pnpm install --frozen-lockfile
COPY apps/web/ apps/web/
WORKDIR /app/apps/web
RUN pnpm build

# ---- Production stage ----
FROM node:22-alpine AS prod

WORKDIR /app
COPY --from=builder /app/apps/web/.next/standalone ./
COPY --from=builder /app/apps/web/.next/static ./.next/static
COPY --from=builder /app/apps/web/public ./public
EXPOSE 3000
CMD ["node", "server.js"]
```

- [ ] **Step 3: Create docker/docker-compose.yml**

```yaml
version: "3.8"

services:
  postgres:
    image: postgres:16-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${POSTGRES_USER:-app}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-secret}
      POSTGRES_DB: ${POSTGRES_DB:-monorepo_dev}
    ports:
      - "${POSTGRES_PORT:-5432}:5432"
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-app}"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    ports:
      - "${REDIS_PORT:-6379}:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  api:
    build:
      context: ..
      dockerfile: docker/api.Dockerfile
      target: dev
    restart: unless-stopped
    ports:
      - "${API_PORT:-8080}:8080"
    environment:
      - POSTGRES_HOST=postgres
      - POSTGRES_PORT=5432
      - POSTGRES_USER=${POSTGRES_USER:-app}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-secret}
      - POSTGRES_DB=${POSTGRES_DB:-monorepo_dev}
      - POSTGRES_SSLMODE=disable
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=${REDIS_PASSWORD:-}
      - JWT_SECRET=${JWT_SECRET:-dev-secret}
      - API_PORT=8080
      - API_LOG_LEVEL=debug
    volumes:
      - ../apps/api:/app
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/healthz"]
      interval: 10s
      timeout: 5s
      retries: 3

  web:
    build:
      context: ..
      dockerfile: docker/web.Dockerfile
      target: dev
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080/graphql
      - NEXT_PUBLIC_APP_NAME=MonorepoApp
    volumes:
      - ../apps/web:/app/apps/web
      - /app/apps/web/node_modules
      - /app/apps/web/.next
    depends_on:
      api:
        condition: service_healthy

volumes:
  postgres-data:
```

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat(docker): add Docker Compose with postgres, redis, api, and web services"
```

---

## Task 12: GitLab CI/CD

**Files:**
- Create: `.gitlab-ci.yml`

- [ ] **Step 1: Create .gitlab-ci.yml**

```yaml
image: node:22-alpine

stages:
  - validate
  - test
  - build

variables:
  PNPM_HOME: /root/.local/share/pnpm
  PATH: /root/.local/share/pnpm:$PATH
  POSTGRES_USER: test
  POSTGRES_PASSWORD: test
  POSTGRES_DB: test_db
  REDIS_HOST: redis

.setup: &setup
  before_script:
    - corepack enable && corepack prepare pnpm@latest --activate
    - pnpm install --frozen-lockfile

# --- VALIDATE ---

lint:ts:
  stage: validate
  <<: *setup
  script:
    - pnpm nx affected --target=lint --base=origin/main
  rules:
    - if: $CI_MERGE_REQUEST_IID

lint:go:
  stage: validate
  image: golang:1.23-alpine
  before_script:
    - go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  script:
    - cd apps/api && golangci-lint run
  rules:
    - if: $CI_MERGE_REQUEST_IID
      changes:
        - "apps/api/**/*.go"

validate:schema:
  stage: validate
  <<: *setup
  script:
    - pnpm nx run graphql:validate
  rules:
    - if: $CI_MERGE_REQUEST_IID
      changes:
        - "libs/graphql/**"

typecheck:
  stage: validate
  <<: *setup
  script:
    - pnpm nx affected --target=typecheck --base=origin/main
  rules:
    - if: $CI_MERGE_REQUEST_IID

commitlint:
  stage: validate
  <<: *setup
  script:
    - pnpm exec commitlint --from $CI_MERGE_REQUEST_DIFF_BASE_SHA --to HEAD
  rules:
    - if: $CI_MERGE_REQUEST_IID

# --- TEST ---

test:ts:
  stage: test
  <<: *setup
  script:
    - pnpm nx affected --target=test --base=origin/main
  rules:
    - if: $CI_MERGE_REQUEST_IID

test:go:
  stage: test
  image: golang:1.23-alpine
  services:
    - name: postgres:16-alpine
      alias: postgres
      variables:
        POSTGRES_USER: test
        POSTGRES_PASSWORD: test
        POSTGRES_DB: test_db
    - name: redis:7-alpine
      alias: redis
  variables:
    POSTGRES_HOST: postgres
    POSTGRES_PORT: "5432"
    POSTGRES_SSLMODE: disable
    REDIS_PORT: "6379"
  script:
    - cd apps/api && go test -coverprofile=coverage.out ./...
    - go tool cover -func=coverage.out
  coverage: '/total:\s+\(statements\)\s+(\d+\.\d+)%/'
  rules:
    - if: $CI_MERGE_REQUEST_IID
      changes:
        - "apps/api/**/*.go"

test:e2e:
  stage: test
  <<: *setup
  script:
    - pnpm nx run web:e2e
  rules:
    - if: $CI_MERGE_REQUEST_IID
      changes:
        - "apps/web/**"
      when: manual
  allow_failure: true

# --- BUILD ---

build:apps:
  stage: build
  <<: *setup
  script:
    - pnpm nx affected --target=build --base=origin/main
  rules:
    - if: $CI_MERGE_REQUEST_IID

build:docker:api:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -f docker/api.Dockerfile --target prod -t api:$CI_COMMIT_SHORT_SHA .
  rules:
    - if: $CI_MERGE_REQUEST_IID
      changes:
        - "apps/api/**"
        - "docker/api.Dockerfile"

build:docker:web:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -f docker/web.Dockerfile --target prod -t web:$CI_COMMIT_SHORT_SHA .
  rules:
    - if: $CI_MERGE_REQUEST_IID
      changes:
        - "apps/web/**"
        - "docker/web.Dockerfile"
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "ci: add GitLab CI/CD pipeline with validate, test, and build stages"
```

---

## Task 13: README

**Files:**
- Create: `README.md`

- [ ] **Step 1: Create README.md**

```markdown
# Monorepo Template

Production-ready monorepo template: Go backend + Next.js frontend + GraphQL.

## Tech Stack

- **Monorepo:** Nx + pnpm workspaces
- **Backend:** Go, chi, gqlgen, zap, goose, viper
- **Frontend:** Next.js (App Router), FSD, TanStack Query, Tailwind CSS
- **API:** GraphQL (schema-first, shared contract)
- **Database:** PostgreSQL + Redis
- **CI/CD:** GitLab CI/CD
- **Quality:** ESLint, Prettier, golangci-lint, lefthook, commitlint

## Quick Start

```bash
# Clone and install
git clone <repo-url>
cd monorepo-template
pnpm install
cp .env.example .env

# Start everything with Docker
cd docker && docker compose up

# Or run locally
pnpm dev
```

## Project Structure

```
apps/
  api/     — Go GraphQL API server
  web/     — Next.js frontend (FSD)
libs/
  graphql/ — Shared GraphQL schema (source of truth)
tools/
  codegen/ — GraphQL code generation config
  nx-go/   — Custom Nx executor for Go
docker/    — Docker Compose + Dockerfiles
```

## Commands

| Command | Description |
|---------|-------------|
| `pnpm dev` | Start all apps in dev mode |
| `pnpm build` | Build all apps |
| `pnpm test` | Run all tests |
| `pnpm lint` | Lint all code |
| `pnpm codegen` | Run GraphQL codegen (Go + TS) |
| `docker compose up` | Start full stack locally (from `docker/`) |

## GraphQL Workflow

1. Edit schema in `libs/graphql/schema/`
2. Run `pnpm codegen` to generate Go types + TS hooks
3. Implement resolvers in `apps/api/internal/graph/`
4. Use generated hooks in `apps/web/src/features/*/api/`

## Commit Convention

Uses [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(api): add user resolver
fix(web): login redirect
chore(deps): update dependencies
```

Scopes: `api`, `web`, `graphql`, `codegen`, `nx-go`, `docker`, `ci`, `deps`
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "docs: add README with quick start guide and project overview"
```

---

## Final Verification

- [ ] **Step 1: Run full lint and type check from root**

```bash
cd /home/nolood/general/rnd/monorepo-template
pnpm install
pnpm lint
pnpm nx affected --target=typecheck
```

- [ ] **Step 2: Verify Go builds**

```bash
cd apps/api && go build ./...
```

- [ ] **Step 3: Run Go tests**

```bash
cd apps/api && go test -short ./...
```

- [ ] **Step 4: Verify Docker Compose config is valid**

```bash
cd docker && docker compose config
```

- [ ] **Step 5: Final commit if any fixes needed**

```bash
git add -A
git commit -m "chore: final verification fixes"
```
