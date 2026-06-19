# Monorepo Template — Production-Ready Design Spec

## Overview

Production-ready monorepo template for fullstack web applications with scalability in mind. Go backend + Next.js frontend, GraphQL as the primary API contract, Nx for orchestration.

## Tech Stack

| Layer | Technology |
|---|---|
| Monorepo orchestrator | Nx |
| Package manager | pnpm |
| Frontend | Next.js (App Router) |
| Frontend architecture | Feature-Sliced Design (FSD) |
| Frontend data fetching | TanStack Query + graphql-request |
| Frontend styling | Tailwind CSS |
| Backend | Go + chi (HTTP router) |
| Backend architecture | Clean Architecture (cmd → internal: graph → service → repository) |
| API contract | GraphQL (schema-first) |
| GraphQL Go | gqlgen |
| GraphQL codegen (TS) | graphql-codegen |
| Database | PostgreSQL |
| Cache | Redis |
| Logging | zap |
| Migrations | goose |
| Config | viper (stable values from config file, sensitive from env) |
| Testing (Go) | go test + testify |
| Testing (TS) | Vitest |
| Testing (e2e) | Playwright |
| Linting (TS) | ESLint + Prettier |
| Linting (Go) | golangci-lint |
| FSD lint | eslint-plugin-boundaries |
| Pre-commit hooks | lefthook |
| Commit convention | conventional commits (commitlint) |
| CI/CD | GitLab CI/CD |
| Containers | Docker + Docker Compose (local dev) |

## Repository Structure

```
monorepo-template/
├── apps/
│   ├── web/                        # Next.js application
│   │   ├── app/                    # App Router (routing, layouts only)
│   │   │   ├── layout.tsx          # root layout, imports providers from src/app/
│   │   │   ├── page.tsx            # imports page composition from src/pages/
│   │   │   ├── (auth)/
│   │   │   └── providers.tsx
│   │   ├── src/
│   │   │   ├── app/                # FSD: app layer
│   │   │   │   ├── styles/
│   │   │   │   │   └── globals.css
│   │   │   │   └── config.ts
│   │   │   ├── pages/              # FSD: page composition components (NOT routing)
│   │   │   ├── widgets/            # FSD: widgets (header, sidebar)
│   │   │   ├── features/           # FSD: features (login, create-post)
│   │   │   │   └── auth/
│   │   │   │       ├── api/
│   │   │   │       ├── model/
│   │   │   │       ├── ui/
│   │   │   │       └── index.ts
│   │   │   ├── entities/           # FSD: entities (user, project)
│   │   │   │   └── user/
│   │   │   │       ├── api/
│   │   │   │       ├── model/
│   │   │   │       ├── ui/
│   │   │   │       └── index.ts
│   │   │   └── shared/             # FSD: shared (no business logic)
│   │   │       ├── api/
│   │   │       │   ├── graphql-client.ts
│   │   │       │   └── generated/
│   │   │       ├── config/
│   │   │       ├── ui/
│   │   │       └── lib/
│   │   ├── e2e/                    # Playwright e2e tests
│   │   │   └── playwright.config.ts
│   │   ├── next.config.js
│   │   ├── tsconfig.json
│   │   └── project.json
│   └── api/                        # Go GraphQL server
│       ├── cmd/
│       │   └── server/
│       │       └── main.go
│       ├── internal/
│       │   ├── graph/
│       │   │   ├── resolver.go
│       │   │   ├── schema.resolvers.go
│       │   │   └── model/models_gen.go
│       │   ├── middleware/
│       │   │   ├── auth.go
│       │   │   ├── logging.go
│       │   │   └── cors.go
│       │   ├── handler/
│       │   │   └── health.go       # GET /healthz, GET /readyz
│       │   ├── service/
│       │   │   └── user_service.go
│       │   ├── repository/
│       │   │   ├── postgres/
│       │   │   │   ├── user_repo.go
│       │   │   │   └── migrations/
│       │   │   └── redis/
│       │   │       └── cache.go
│       │   └── config/
│       │       └── config.go
│       ├── config/
│       │   └── config.yml          # stable config values
│       ├── air.toml                # hot reload config for dev
│       ├── go.mod                  # module: monorepo-template/apps/api
│       ├── go.sum
│       ├── gqlgen.yml
│       └── project.json
├── libs/
│   └── graphql/                    # Shared GraphQL schema (source of truth)
│       ├── schema/
│       │   ├── schema.graphql      # Root: Query, Mutation
│       │   ├── user.graphql
│       │   └── common.graphql      # Scalars (DateTime, UUID), pagination, errors
│       └── project.json
├── tools/
│   ├── codegen/
│   │   └── codegen.ts              # graphql-codegen config
│   └── nx-go/                      # Custom Nx executor for Go
│       ├── src/
│       │   ├── executors/
│       │   │   ├── build/executor.ts
│       │   │   ├── serve/executor.ts
│       │   │   ├── test/executor.ts
│       │   │   └── lint/executor.ts
│       │   └── index.ts
│       ├── executors.json
│       ├── package.json
│       └── tsconfig.json
├── docker/
│   ├── docker-compose.yml
│   ├── api.Dockerfile
│   └── web.Dockerfile
├── .gitlab-ci.yml
├── .lefthook.yml
├── .eslintrc.json
├── .prettierrc
├── .gitignore
├── commitlint.config.js
├── nx.json
├── pnpm-workspace.yaml
├── package.json
├── tsconfig.base.json
├── .env.example
└── README.md
```

## Go API Architecture

### Clean Architecture Layers

1. **cmd/server/main.go** — entrypoint: config loading, DI, HTTP server startup
2. **internal/graph/** — gqlgen resolvers (transport layer)
3. **internal/middleware/** — HTTP middleware (auth, logging, CORS)
4. **internal/handler/** — REST endpoints (health checks)
5. **internal/service/** — business logic, interfaces for repositories
6. **internal/repository/** — data access implementations (postgres, redis)
7. **internal/config/** — viper-based configuration

### Health Check Endpoints

- `GET /healthz` — liveness probe (server is running)
- `GET /readyz` — readiness probe (database and redis connections are healthy)

Used by Docker Compose `healthcheck`, and ready for future K8s probes.

### Key Decisions

- **gqlgen** schema-first: write `.graphql` → generate Go types and resolver stubs
- **Dependency injection** via struct embedding in resolver, no DI frameworks (wire/dig overkill at start)
- **Repository pattern** — interfaces defined in `service/`, implementations in `repository/`. Enables mocking in tests
- **Migrations** — goose, SQL files in `repository/postgres/migrations/`
- **Logging** — zap (structured, performant)
- **HTTP router** — chi (lightweight, net/http compatible, middleware chain)
- **Configuration** — viper: stable values from `config/config.yml` (server port, log level, pagination defaults), sensitive values from environment variables (DB credentials, JWT secret, Redis password). Separate config module merges both sources.
- **Go module path** — `monorepo-template/apps/api` (template users replace `monorepo-template` with their project name)

### gqlgen Configuration (gqlgen.yml)

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

## Frontend Architecture (FSD)

### Feature-Sliced Design Layers

Strict import rule: a layer can only import from layers below it.

1. **app** — providers, global styles, global config
2. **pages** — page composition from features and widgets
3. **widgets** — self-contained UI blocks (header, sidebar)
4. **features** — business actions (login, create-post)
5. **entities** — business entities (user, project)
6. **shared** — reusable code without business logic (UI kit, utils, API client)

### FSD Pages vs Next.js App Router

`app/` directory (Next.js App Router) handles **routing only** — route segments, layouts, loading/error states. Each route file imports a page composition component from `src/pages/`.

`src/pages/` (FSD pages layer) handles **page composition** — assembling widgets, features, and entities into a complete page view. These are regular React components, not routing primitives.

Example: `app/dashboard/page.tsx` imports and renders `src/pages/dashboard/ui/DashboardPage.tsx`.

### Slice Structure

Each feature/entity follows the same structure:
- `api/` — GraphQL operations + generated hooks
- `model/` — state, types
- `ui/` — components
- `index.ts` — public API (re-exports)

### Key Decisions

- **App Router** — server components by default, client components only where interactivity needed
- **graphql-codegen** generates typed TanStack Query hooks from schema — write `.graphql` operation, get ready `useUsersQuery()` hook
- **Tailwind CSS** — utility-first, zero-runtime
- **eslint-plugin-boundaries** — enforces FSD import rules in CI
- Base types generated to `shared/api/generated/`, features write operations in `features/*/api/*.graphql`

## GraphQL Schema & Codegen

### Schema Location

`libs/graphql/schema/` — single source of truth for both Go and TS sides.

### Error Handling Strategy

Union-based errors in GraphQL schema. Each mutation returns a union type:

```graphql
type CreateUserSuccess {
  user: User!
}

type ValidationError {
  field: String!
  message: String!
}

type AuthError {
  message: String!
}

union CreateUserResult = CreateUserSuccess | ValidationError | AuthError
```

- Domain errors are part of the schema (union members), not transport errors
- Transport errors (500, network) surface through GraphQL top-level `errors` array
- Go resolvers return domain errors as typed values, not `error` interface
- Frontend codegen generates discriminated unions — handle with `__typename` switch

### Codegen Flow

```
libs/graphql/schema/*.graphql  (source of truth)
        │
        ├──→ gqlgen (Go)
        │     output: apps/api/internal/graph/model/ + resolver stubs
        │     config: apps/api/gqlgen.yml
        │
        └──→ graphql-codegen (TS)
              output: apps/web/src/shared/api/generated/
              - base types (TypedDocumentNode)
              - TanStack Query hooks from operations in features/*/api/*.graphql
              config: tools/codegen/codegen.ts
```

### graphql-codegen Configuration (tools/codegen/codegen.ts)

```typescript
// Plugins:
// - @graphql-codegen/typescript — base types from schema
// - @graphql-codegen/typescript-operations — types from .graphql operations
// - @graphql-codegen/typescript-react-query — TanStack Query hooks

// Schema source: libs/graphql/schema/**/*.graphql

// Documents (operations) source:
//   - apps/web/src/features/**/api/**/*.graphql
//   - apps/web/src/entities/**/api/**/*.graphql

// Output:
//   - Base types → apps/web/src/shared/api/generated/types.ts
//   - Hooks → apps/web/src/shared/api/generated/hooks.ts
```

### libs/graphql project.json

```json
{
  "targets": {
    "validate": {
      "executor": "nx:run-commands",
      "options": {
        "command": "npx graphql-inspector introspect libs/graphql/schema/*.graphql --write schema.json && npx graphql-inspector validate libs/graphql/schema/*.graphql"
      }
    }
  }
}
```

### Nx Targets

- `nx run graphql:validate` — schema validation via graphql-inspector
- `nx run api:codegen` — gqlgen generate
- `nx run web:codegen` — graphql-codegen
- `nx run codegen` — both in parallel
- Dependency: codegen depends on `libs/graphql`, Nx rebuilds on schema change

## Custom Nx Go Executor (tools/nx-go)

TypeScript-based Nx executor plugin that wraps Go CLI commands. Located at `tools/nx-go/`.

### Executors

| Executor | Go Command | Options |
|---|---|---|
| `nx-go:build` | `go build -o <outputPath> ./cmd/server` | `outputPath`, `main` (entrypoint path) |
| `nx-go:serve` | `air` (dev) / `go run ./cmd/server` (fallback) | `port`, `configPath` |
| `nx-go:test` | `go test ./...` | `coverage` (bool), `short` (bool), `packages` (string[]) |
| `nx-go:lint` | `golangci-lint run` | `fix` (bool), `config` (path to .golangci.yml) |

### apps/api/project.json

```json
{
  "targets": {
    "build": {
      "executor": "nx-go:build",
      "options": { "outputPath": "dist/apps/api", "main": "cmd/server" }
    },
    "serve": {
      "executor": "nx-go:serve",
      "options": { "port": 8080 }
    },
    "test": {
      "executor": "nx-go:test",
      "options": { "coverage": true }
    },
    "go-lint": {
      "executor": "nx-go:lint",
      "options": { "fix": false }
    },
    "codegen": {
      "executor": "nx:run-commands",
      "options": { "command": "go run github.com/99designs/gqlgen generate", "cwd": "apps/api" },
      "dependsOn": ["^validate"]
    }
  }
}
```

## Environment Variables

### .env.example

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

Stable values (log level defaults, pagination, timeouts) live in `apps/api/config/config.yml`. Sensitive values (credentials, secrets) come from env and override config file values via viper.

## Docker & Local Development

### docker-compose.yml Services

| Service | Image/Build | Ports | Purpose |
|---|---|---|---|
| `postgres` | postgres:16-alpine | 5432 | Database |
| `redis` | redis:7-alpine | 6379 | Cache |
| `api` | build: docker/api.Dockerfile | 8080 | Go GraphQL server |
| `web` | build: docker/web.Dockerfile | 3000 | Next.js dev server |

### Key Decisions

- Dockerfiles live in `docker/` directory (not inside apps)
- **api.Dockerfile** — multi-stage: builder (go build) + alpine for prod, dev stage with air (hot reload)
- **web.Dockerfile** — Node.js alpine, dev uses `pnpm dev` via compose
- Volumes: postgres data persistent, Go and Next.js sources mounted for hot reload
- `.env.example` — template for environment variables, copy to `.env`
- `docker-compose.yml` reads `.env` file
- `api` service uses `healthcheck` against `/healthz` for `depends_on` conditions
- Single `docker compose up` — full working environment

## CI/CD (GitLab)

### Pipeline Stages

```yaml
stages:
  - validate
  - test
  - build

validate:
  - nx affected --target=lint          # ESLint + Prettier (JS/TS)
  - nx affected --target=go-lint       # golangci-lint via nx-go executor
  - nx run graphql:validate            # schema validation
  - nx affected --target=typecheck     # tsc --noEmit
  - pnpm exec commitlint --from $CI_MERGE_REQUEST_DIFF_BASE_SHA  # commit message validation

test:
  - nx affected --target=test          # Vitest (unit) + Go tests
  - nx run web:e2e                     # Playwright (MR to main only)

build:
  - nx affected --target=build
  - docker build (api, web)
```

### CI Services

PostgreSQL + Redis run as GitLab services for Go integration tests.

### `nx affected`

Runs only targets affected by changes in MR. Saves CI time. Go linting runs through nx-go executor for consistent affected/caching behavior.

## Linting & Hooks

### Lefthook (pre-commit)

```yaml
pre-commit:
  parallel: true
  commands:
    lint-staged:
      run: pnpm exec lint-staged
    go-lint:
      glob: "*.go"
      run: golangci-lint run --fix
    go-test:
      glob: "*.go"
      run: go test -short ./...

commit-msg:
  commands:
    commitlint:
      run: pnpm exec commitlint --edit {1}
```

Note: pre-commit runs `go test -short` (unit tests only) for speed. Full test suite (including integration) runs in CI.

### Linters

- **ESLint** — strict config + `eslint-plugin-boundaries` (FSD imports enforcement)
- **Prettier** — formatting for JS/TS/JSON/YAML/GraphQL
- **golangci-lint** — Go linters (errcheck, gosec, govet, staticcheck, etc.)
- **commitlint** — conventional commits format (`feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `test:`, `ci:`)

### Conventional Commits

Format: `type(scope): description`
- Scope is optional
- Examples: `feat(api): add user resolver`, `fix(web): login redirect`

## Testing Strategy

### Go (apps/api)

- `go test` + **testify** — unit tests for services and resolvers
- Integration tests — postgres/redis via testcontainers or test fixtures
- `go test -cover` — coverage report in CI
- Default coverage gate: 70%

### Next.js (apps/web)

- **Vitest** — unit/integration tests for components and hooks
- Tests live next to code: `features/auth/__tests__/`, `entities/user/__tests__/`
- **Playwright** — e2e tests in `apps/web/e2e/`
- Runs in CI on MR to main
- Default coverage gate: 70% (configured in vitest.config.ts)

### GraphQL Schema

- Schema validation in CI validate stage via graphql-inspector — catches breaking changes

### Coverage Gates

- Default threshold: 70% for both Go and TS
- Go: configured via CI script flag (`go test -coverprofile=coverage.out && go tool cover -func=coverage.out`)
- TS: configured in `vitest.config.ts` under `coverage.thresholds`
- Blocks MR if below threshold
