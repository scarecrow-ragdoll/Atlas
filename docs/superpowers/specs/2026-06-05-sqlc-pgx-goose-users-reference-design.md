<!-- FILE: docs/superpowers/specs/2026-06-05-sqlc-pgx-goose-users-reference-design.md -->
<!-- VERSION: 1.0.1 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Capture the approved design for making sqlc plus pgx/v5 plus goose the backend SQL standard through the users reference slice. -->
<!--   SCOPE: Design-level architecture, boundaries, data flow, error handling, testing, GRACE updates, and implementation constraints; excludes implementation code. -->
<!--   DEPENDS: docs/requirements.xml, docs/technology.xml, docs/development-plan.xml, docs/knowledge-graph.xml, docs/verification-plan.xml, apps/api persistence structure. -->
<!--   LINKS: M-API / M-WORKSPACE / M-GRACE-WORKFLOW / V-M-API / V-M-GRACE-WORKFLOW. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Goal - Defines the approved sqlc, pgx/v5, and goose backend standard. -->
<!--   Architecture - Defines the repository adapter boundary and generated-code ownership. -->
<!--   Components And Data Flow - Defines sqlc config, query files, generated package, UserRepo adapter, and Nx/root codegen flow. -->
<!--   Error Handling - Preserves existing user repository observable behavior. -->
<!--   Testing And Verification - Defines focused tests, codegen drift gates, coverage policy, and GRACE artifact updates required for implementation. -->
<!--   File-Local Markup Plan - Defines markup obligations for handwritten files and generated sqlc output. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.1 - Incorporated subagent review findings for workspace contracts, generated drift, coverage, nullable params, and safe integration tests. -->
<!-- END_CHANGE_SUMMARY -->

# sqlc pgx goose users reference design

**Status:** Approved
**Date:** 2026-06-05

## Goal

Make `sqlc + pgx/v5 + goose` the backend SQL standard for the Go API, using the existing `users` vertical slice as the first reference implementation.

The implementation must keep the current public behavior of the API stable:

- Admin GraphQL and public REST continue to call the shared `UserService`.
- `UserService` continues to depend on a transport-neutral `UserRepository` contract.
- `goose` migrations remain the source of truth for PostgreSQL schema changes.
- `sqlc` owns generated, typed query code for handwritten SQL.
- `pgx/v5` remains the runtime PostgreSQL driver and pool implementation.

This is not a new product domain. The first case is the existing `users` slice because it already proves REST, GraphQL, service, repository, PostgreSQL, test, and e2e paths.

## Current Context

The API already uses `github.com/jackc/pgx/v5` through `pgxpool`, and it already applies embedded `goose` migrations from `apps/api/internal/repository/postgres/migrations`.

The current `UserRepo` keeps SQL strings inside Go methods. That makes the reference slice less useful as a template standard because downstream teams have to infer query ownership, schema drift checks, and generated-code expectations from scattered handwritten SQL.

`sqlc` is not currently wired into the repository. The root `bun run codegen` script runs API gqlgen and web-admin GraphQL codegen, but not backend SQL generation.

## Key Decisions

| Decision                   | Choice                                               | Rationale                                                                                            |
| -------------------------- | ---------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| First sqlc case            | Existing `users` slice                               | Proves the real reference flow instead of a toy example.                                             |
| Public repository boundary | Keep `UserRepo` as the adapter used by `UserService` | Avoids leaking generated sqlc models into service, REST, or GraphQL layers.                          |
| Schema source              | `goose` migrations                                   | The API already runs embedded migrations at startup; sqlc should read the same schema source.        |
| SQL ownership              | `queries/*.sql` files consumed by sqlc               | Keeps SQL reviewable and typed without embedding it in Go code.                                      |
| Runtime driver             | `pgx/v5`                                             | Already approved in `docs/technology.xml` and supported by sqlc through `sql_package: "pgx/v5"`.     |
| Tool ownership             | Track sqlc through the API Go module and Nx target   | Keeps backend codegen reproducible locally and in CI.                                                |
| Generated code             | Commit sqlc output                                   | Matches the repository's existing generated-artifact pattern for gqlgen and web-admin GraphQL types. |
| Generated coverage         | Allowlist sqlc generated Go with replacement gates   | Generated query code is proven by codegen, build, and repository integration coverage.               |
| Drift detection            | Require a clean generated diff after codegen         | Regenerating code is not enough; stale committed output must fail a local or CI gate.                |

## Architecture

`goose` remains the owner of database schema migrations. Migration SQL stays under:

```text
apps/api/internal/repository/postgres/migrations
```

`sqlc` reads those migrations as its schema input and reads query files from:

```text
apps/api/internal/repository/postgres/queries
```

Generated Go code goes into an internal generated package, for example:

```text
apps/api/internal/repository/postgres/generated
```

`UserRepo` remains the public persistence adapter for the user domain. It should call generated `Queries` methods internally, then map generated rows and errors back into the current `service.User` and `UserRepository` behavior.

The service, REST handlers, GraphQL resolvers, web-admin client, and public web client should not import the sqlc generated package.

The intended dependency flow is:

```text
GraphQL resolver or REST handler
  -> UserService
  -> service.UserRepository interface
  -> postgres.UserRepo adapter
  -> sqlc generated Queries
  -> pgxpool.Pool
  -> PostgreSQL schema managed by goose migrations
```

## Components And Data Flow

### `apps/api/sqlc.yaml`

Add a sqlc v2 configuration owned by the API module.

Expected shape:

```yaml
version: '2'
sql:
  - engine: 'postgresql'
    schema: 'internal/repository/postgres/migrations'
    queries: 'internal/repository/postgres/queries'
    gen:
      go:
        package: 'generated'
        out: 'internal/repository/postgres/generated'
        sql_package: 'pgx/v5'
        emit_interface: true
```

The final implementation can adjust the generated package name if the repository has a clearer local convention, but it must remain internal to the postgres repository area.

### `internal/repository/postgres/queries/users.sql`

Add named sqlc queries for the existing repository behavior:

- `GetUserByID :one`
- `ListUsers :many`
- `CountUsers :one`
- `CreateUser :one`
- `UpdateUser :one`
- `DeleteUser :exec`

The SQL should preserve the current observable semantics:

- user rows return `id`, `email`, `name`, `created_at`, and `updated_at`;
- list ordering remains `created_at DESC`;
- cursor paging keeps the current `created_at < cursor` behavior;
- list still fetches `limit + 1` so the adapter can trim the extra row;
- count remains separate through `CountUsers`;
- create and update continue returning the written row;
- delete remains idempotent for missing rows.

`UpdateUser` must preserve the current optional-field contract from `service.UpdateUserInput`. The query should use explicit nullable sqlc parameters, such as `sqlc.narg('name')` and `sqlc.narg('email')`, or an equivalent documented pgx nullable mapping. The implementation must prove name-only, email-only, and no-change update behavior.

### `internal/repository/postgres/generated`

This package is sqlc-owned. Developers and agents must not edit generated files by hand.

Generated code should be committed so downstream users can build the template without an implicit generation step, while CI and local gates still prove codegen is current.

Generated sqlc files are generated artifacts, not handwritten behavior. They should not receive manual file-local GRACE markup or hand-written edits. Their contract is owned by `apps/api/sqlc.yaml`, the goose migrations, query SQL, generated-drift checks, and repository integration tests.

### `internal/repository/postgres/user_repo.go`

`UserRepo` becomes a thin adapter.

It owns:

- constructing generated `Queries` from the pgx pool;
- cursor decoding and validation;
- default page size and `limit + 1` trimming;
- mapping generated rows to `service.User`;
- formatting timestamps as RFC3339Nano strings;
- mapping `pgx.ErrNoRows` to `nil, nil` where current behavior requires it;
- mapping unique violation `23505` to the existing duplicate-email error text;
- preserving operation-name error wrapping.

It must not reintroduce inline SQL for user CRUD after the sqlc migration.

### Nx and root codegen

Update `apps/api/project.json` so `api:codegen` runs sqlc generation and gqlgen generation in a deterministic order.

Preferred command shape:

```text
cd apps/api && go run github.com/sqlc-dev/sqlc/cmd/sqlc generate && go run github.com/99designs/gqlgen generate
```

Root `bun run codegen` should continue to be the single codegen entry point and should include the API target. The API target now covers both sqlc and gqlgen.

Track `github.com/sqlc-dev/sqlc/cmd/sqlc` in the API Go tool dependencies so local and CI execution do not rely on a globally installed `sqlc` binary.

Add an explicit generated-drift assertion. Acceptable shapes include:

```text
bun run codegen
git diff --exit-code -- apps/api/internal/repository/postgres/generated apps/api/internal/graph apps/web-admin/src/shared/api/generated
```

or a repository-local `codegen:validate` target that runs codegen and fails when generated output differs from committed files. The implementation should update local and CI validation surfaces consistently if CI already has a codegen validation job.

### Coverage policy for sqlc output

Add sqlc generated Go files to the generated-code coverage policy instead of treating them as handwritten behavior.

Required implementation updates:

- `tools/coverage/coverage.config.json` includes the sqlc generated `.go` files or a glob that the coverage runner supports.
- `tools/coverage/run.mjs` is updated only if the current allowlist matching cannot express the generated sqlc paths safely.
- `docs/verification-plan.xml` documents the sqlc generated-code allowlist and replacement gate.
- The replacement gate is at least `bunx nx run api:codegen && bunx nx build api`, plus repository integration tests that exercise generated queries against the goose schema.

## Error Handling

The migration must preserve current repository behavior:

- `GetByID` returns `nil, nil` when no user exists.
- `Update` returns `nil, nil` when no user exists.
- `Delete` succeeds when the user is already missing.
- duplicate email failures return an error containing `duplicate email`.
- invalid base64 cursor and invalid timestamp cursor fail before executing SQL.
- query, scan, count, create, update, and delete failures are wrapped with stable operation context.

The implementation should keep duplicate detection based on PostgreSQL unique-violation code `23505`.

## Testing And Verification

Focused implementation checks should start with the API package:

```text
bunx nx run api:codegen
bunx nx test api
bunx nx lint api
bunx nx build api
```

Repository integration tests remain the primary proof that generated sqlc queries match the real goose schema.

Repository integration tests must keep the existing destructive-test guardrails:

- use the dedicated `monorepo_test` PostgreSQL target and test compose port `17501`;
- call `testinfra.RequireSafePostgresDSN` before cleanup;
- keep destructive cleanup, such as `TRUNCATE users`, behind the safe target guard;
- report final evidence that the sqlc/goose integration path actually ran rather than only that tests exited after skipping an unavailable database.

Adapter unit tests should stay focused on behavior that sqlc does not own:

- invalid cursor handling;
- limit defaulting and trimming;
- duplicate error mapping;
- no-row mapping;
- nullable update mapping for name-only, email-only, and empty updates;
- query or command error wrapping where practical.

The adapter rewrite should preserve or explicitly replace current repository debug observability, including operation names and stable `user_id` or `email` fields where they already exist.

Existing service, REST handler, and GraphQL resolver tests should remain stable because their public contract stays behind `UserService` and `UserRepository`.

Broader follow-up checks before final implementation closeout:

```text
bun run codegen
bunx nx run graphql:validate
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Run full root build, coverage, and e2e gates only at final closeout or when the implementation touches shared contracts that require broader evidence.

## GRACE Contract Updates

The implementation should update the shared GRACE artifacts:

- `docs/requirements.xml`: update the codegen use case so backend SQL codegen is included alongside GraphQL generation, and resolve the generated-files policy for sqlc output as committed generated code with drift gates.
- `docs/technology.xml`: add sqlc as an approved API codegen tool and dependency/tooling surface.
- `docs/development-plan.xml`: update `M-API` to describe sqlc-generated PostgreSQL query ownership and goose schema ownership; update `M-WORKSPACE` for the expanded API/root codegen command surface; update `M-COVERAGE-GATE` for sqlc generated-code allowlist ownership.
- `docs/knowledge-graph.xml`: add paths for `apps/api/sqlc.yaml`, user queries, and generated sqlc package under `M-API`; update `M-WORKSPACE` annotations for API/root codegen; update coverage graph paths if `tools/coverage/coverage.config.json` or `tools/coverage/run.mjs` changes.
- `docs/verification-plan.xml`: add `bunx nx run api:codegen` as the sqlc plus gqlgen API generation gate; add generated-drift assertions to the relevant workspace/codegen gate; include repository generated-query integration coverage in `V-M-API`; update `V-M-WORKSPACE` and `V-M-COVERAGE-GATE` for changed project targets and sqlc generated-code coverage policy.
- `docs/operational-packets.xml`: no update is required unless implementation planning introduces new worker packet fields.

New or meaningfully edited governed files outside `docs/*.xml` must include file-local GRACE markup before implementation closeout.

## File-Local Markup Plan

Handwritten governed files created or meaningfully edited by the implementation must receive or refresh file-local GRACE markup before closeout.

Expected markup obligations:

- `apps/api/sqlc.yaml`: YAML file-level contract describing sqlc config ownership, schema input, query input, generated output, and codegen gate.
- `apps/api/internal/repository/postgres/queries/users.sql`: SQL file-level contract in SQL comments describing users query ownership and links to `M-API` / `V-M-API`.
- `apps/api/internal/repository/postgres/user_repo.go`: refresh or add file-level contract plus function anchors for non-trivial adapter behavior such as list pagination, create/update error mapping, and delete semantics.
- repository tests touched for sqlc behavior: file-level test contract and focused anchors where they clarify safe destructive setup or edge-case coverage.
- `apps/api/project.json`, `tools/coverage/coverage.config.json`, and `tools/coverage/run.mjs` if edited: config/tooling contracts aligned with their existing local style and GRACE refs.

Generated sqlc files under `apps/api/internal/repository/postgres/generated` must not be manually edited only to add markup. Their generated status must instead be captured in the source config/query files, shared GRACE docs, coverage allowlist, and generated-drift gate.

## Rollback

This rollout is not intended to change the database schema. Rollback should be a code, tooling, docs, and generated-artifact revert. Do not run a goose down migration for this change unless a later implementation unexpectedly adds a schema migration and explicitly documents that migration's rollback path.

## Out Of Scope

- Adding a new backend domain only to demonstrate sqlc.
- Replacing `UserService` with generated sqlc interfaces.
- Moving REST or GraphQL handlers to generated sqlc models.
- Replacing goose with sqlc schema files.
- Reworking migrations beyond what sqlc needs to parse the existing schema.
- Running a repo-wide raw SQL migration outside the `users` reference slice.

## Risks And Stop Conditions

Stop and replan if:

- sqlc cannot parse the existing goose migration directory without schema adjustments that would change runtime behavior;
- generated types force sqlc package leakage into service, REST, or GraphQL layers;
- API codegen becomes ambiguous about whether sqlc or gqlgen failed;
- generated code drifts but the root codegen gate does not detect it;
- implementation requires weakening repository integration tests or coverage allowlists.

## Sources

- sqlc configuration docs: https://docs.sqlc.dev/en/latest/reference/config.html
- sqlc pgx guide: https://docs.sqlc.dev/en/v1.22.0/guides/using-go-and-pgx.html
- sqlc DDL and migration parsing docs: https://docs.sqlc.dev/en/latest/howto/ddl.html
- sqlc named parameters docs: https://docs.sqlc.dev/en/latest/howto/named_parameters.html
- goose embedded migrations docs: https://github.com/pressly/goose
- goose annotations docs: https://pressly.github.io/goose/documentation/annotations/

## Approval

Approved direction: make `sqlc + pgx/v5 + goose` the backend SQL standard for new API SQL work, with the existing `users` slice as the first reference case. Use `UserRepo` as a thin adapter around sqlc generated queries so the service, REST, and GraphQL contracts remain stable.
