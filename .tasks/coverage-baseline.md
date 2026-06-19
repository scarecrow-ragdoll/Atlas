# 100 Percent Coverage Baseline

## Commands

| Command                                                         | Exit | Notes                                                                                                                                                                       |
| --------------------------------------------------------------- | ---: | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `bunx nx show projects`                                         |    0 | Projects: `go-config`, `go-logger`, `codegen`, `graphql`, `nx-go`, `api`, `bot`, `web`.                                                                                     |
| `bunx nx show project api --json`                               |    0 | `api` has `build`, `serve`, `test`, `lint`, and `codegen`; `test` wrote local `coverage.out`.                                                                               |
| `bunx nx show project bot --json`                               |    0 | `bot` has `build`, `serve`, `test`, and `lint`; `test` wrote local `coverage.out`.                                                                                          |
| `bunx nx show project web --json`                               |    0 | `web` has unit, typecheck, codegen, and e2e targets; no `test-coverage` target yet.                                                                                         |
| `bunx nx show project nx-go --json`                             |    0 | `nx-go` is discoverable but has no targets.                                                                                                                                 |
| `bunx nx show project codegen --json`                           |    0 | `codegen` is discoverable but has no targets.                                                                                                                               |
| `bun run test`                                                  |    1 | Initial run failed on Go coverage because the local Go 1.25 toolchain was missing `covdata`; before the failure, `go-config` reported 81.1% and `go-logger` reported 97.4%. |
| `cd apps/web && bun run test`                                   |    0 | One existing page test passed; web coverage is not enforced at 100%.                                                                                                        |
| `cd apps/api && go test -coverprofile=coverage.out ./...`       |    1 | Initial run failed with `go: no such tool "covdata"` for packages without tests; service package reported 41.0%.                                                            |
| `cd libs/go/config && go test -coverprofile=coverage.out ./...` |    0 | Local `coverage.out` emitted; package reported 81.1%.                                                                                                                       |
| `cd libs/go/logger && go test -coverprofile=coverage.out ./...` |    0 | Local `coverage.out` emitted; package reported 97.4%.                                                                                                                       |

## Current Coverage Artifacts

| Project     | Artifact                      | Status                                             |
| ----------- | ----------------------------- | -------------------------------------------------- |
| `api`       | `apps/api/coverage.out`       | ignored local artifact from current project target |
| `bot`       | `apps/bot/coverage.out`       | ignored local artifact from current project target |
| `go-config` | `libs/go/config/coverage.out` | ignored local artifact from current project target |
| `go-logger` | `libs/go/logger/coverage.out` | ignored local artifact from current project target |
| `web`       | none                          | unit tests run without coverage output             |
| `nx-go`     | none                          | no target yet                                      |
| `codegen`   | none                          | no target yet                                      |

## Handwritten Source Inventory

| Surface                       | Files | Current Test Surface                                     |
| ----------------------------- | ----: | -------------------------------------------------------- |
| API Go app                    |    15 | appconfig, health, middleware, and partial service tests |
| Bot Go app                    |     9 | handler and middleware tests                             |
| Shared Go config              |     4 | config tests with env/YAML cases                         |
| Shared Go logger              |     4 | logger, context, and middleware tests                    |
| Web app and shared TypeScript |    10 | one home page test                                       |
| Workspace tooling             |     5 | no tests                                                 |

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
