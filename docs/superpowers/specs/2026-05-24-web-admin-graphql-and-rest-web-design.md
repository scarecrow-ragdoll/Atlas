# Web Admin GraphQL and REST Web Design

## Goal

Add two distinct frontend examples to the template:

- `web-admin`: the existing Next.js admin UI that uses GraphQL.
- `web`: a new React + Vite public UI that uses REST.

The Go API must serve both clients at the same time while keeping user-domain logic shared through the existing service and repository layers.

## Context

The current template has one frontend app, `apps/web`, implemented with Next.js, React Query, GraphQL operations, GraphQL code generation, and Playwright coverage for the users flow. The API already uses `chi` and exposes health, readiness, GraphQL, and the development GraphQL playground. The user domain is already centralized behind `service.UserService` and a PostgreSQL-backed user repository.

This change splits the transport examples without duplicating business logic.

## Recommended Architecture

Rename the current `apps/web` application to `apps/web-admin` and rename its Nx project to `web-admin`. Keep its Next.js, GraphQL, generated types, and admin-oriented user flow intact.

Create a new `apps/web` application as a Vite + React app. It should use React Query and a small typed REST client instead of GraphQL code generation. The public web example should exercise the same reference user vertical slice through REST endpoints.

Update the Go API from a GraphQL-only application contract to a shared HTTP API contract:

- Admin GraphQL surface:
  - `POST /graphql`
  - `GET /playground` outside production
- Public REST surface:
  - `GET /api/users`
  - `POST /api/users`
  - `GET /api/users/{id}`
  - `PATCH /api/users/{id}`
  - `DELETE /api/users/{id}`

Both surfaces call `service.UserService`. GraphQL resolvers and REST handlers are transport adapters only. They map transport-specific requests and responses, but they must not own user persistence or password hashing behavior.

## Backend Design

The API should introduce REST handlers in a dedicated handler package area, for example `apps/api/internal/handler/users.go`, backed by tests in `apps/api/internal/handler/users_test.go`.

The REST handler should depend on a narrow user service interface matching the methods it uses. This keeps handler tests fast and avoids booting PostgreSQL for HTTP response behavior.

REST response shapes should be stable and simple:

```json
{
  "data": {
    "id": "user-id",
    "email": "user@example.com",
    "name": "User",
    "createdAt": "2026-05-24T00:00:00Z",
    "updatedAt": "2026-05-24T00:00:00Z"
  }
}
```

List responses should include `data` and `meta.totalCount`. Error responses should include an `error` object with `code`, `message`, and optional `field`.

Common domain errors should be represented once, preferably in the service layer or a small domain error helper. GraphQL and REST should map the same duplicate-email and not-found conditions consistently instead of independently string-matching repository errors.

## Frontend Design

`web-admin` remains the GraphQL example:

- Next.js app router.
- GraphQL code generation from `libs/graphql/schema`.
- GraphQL user CRUD/admin flow.
- Existing Playwright GraphQL contract coverage moves with the app.

`web` becomes the REST example:

- Vite + React.
- React Query.
- No GraphQL documents or generated GraphQL types.
- Typed REST client for the user flow.
- A small route structure that supports users list, create, and detail. Update and delete can be exposed if the current reference flow already covers them in the implementation plan; the minimum useful public example is list/create/detail plus duplicate-email error handling.

The public web app should point to the API base URL, not the GraphQL endpoint. A likely environment variable is `VITE_API_BASE_URL`, defaulting to `http://localhost:8080`.

## Workspace And Tooling

Nx and root scripts should reflect the split:

- `web-admin` owns GraphQL codegen and generated admin client types.
- `web` has no GraphQL codegen target.
- `bun run codegen` should run `api` and `web-admin`, not `web`.
- `bun run test:e2e` should prove the relevant browser/API flows for the renamed admin app and the new REST app.
- Coverage allowlists should move GraphQL generated frontend files from `apps/web` to `apps/web-admin`.

Docker, CI, and deployment templates should continue to build the public `web` image as the primary frontend image unless a separate admin image is explicitly introduced in a later change. This change does not require a separate production admin deployment.

## GRACE Contract Updates

Update the durable GRACE artifacts to make the ownership explicit:

- `M-API`: shared Go HTTP API with admin GraphQL and public REST surfaces.
- `M-WEB-ADMIN`: Next.js GraphQL admin app.
- `M-WEB`: Vite REST public app.
- `M-GRAPHQL-SCHEMA`: admin GraphQL schema and codegen contract only.
- Coverage and e2e entries should distinguish REST web checks from admin GraphQL checks.

Requirements and technology docs should stop describing GraphQL as the web-wide public contract. GraphQL is admin-only after this change.

## Testing Strategy

Follow TDD for behavior changes:

1. Add failing API handler tests for REST users list/create/detail/update/delete and error mapping.
2. Add or update GraphQL resolver tests only where shared domain error mapping changes.
3. Add failing workspace/tooling checks for `web-admin` project naming and codegen ownership.
4. Add failing Vite web tests around REST client behavior and page states.
5. Update Playwright coverage so at least one admin GraphQL path and one public REST browser path are proven against the real API.

Focused checks should include:

- `bunx nx test api`
- `bunx nx run graphql:validate`
- `bunx nx run api:codegen`
- `bunx nx run web-admin:codegen`
- `bunx nx test web-admin`
- `bunx nx test web`
- `bunx nx run web-admin:typecheck`
- `bunx nx run web:typecheck`
- `bun run build`
- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`
- `grace lint --path .`

Broader coverage and e2e gates should run before final handoff if the local environment has Docker and Playwright available.

## Risks And Stop Conditions

Stop and replan if:

- GraphQL generated files drift outside `web-admin` ownership.
- REST handlers start duplicating user business logic instead of calling `UserService`.
- The root scripts cannot clearly distinguish `web` from `web-admin`.
- Coverage allowlists are broadened instead of being moved to the renamed generated paths.
- Docker or CI image ownership requires a separate admin deployment decision.

## Approval

Approved direction: keep GraphQL admin-only, add a new REST public web example, and share backend user logic through the existing service/repository boundary.
