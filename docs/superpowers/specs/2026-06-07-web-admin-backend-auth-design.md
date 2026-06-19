<!-- FILE: docs/superpowers/specs/2026-06-07-web-admin-backend-auth-design.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Capture the approved design for backend-only web-admin authentication. -->
<!--   SCOPE: Design-level architecture, data model, session strategy, GraphQL auth surface, seed behavior, error handling, testing, and GRACE update expectations; excludes implementation code and login UI. -->
<!--   DEPENDS: AGENTS.md, docs/requirements.xml, docs/technology.xml, docs/development-plan.xml, docs/knowledge-graph.xml, docs/verification-plan.xml, apps/api, libs/graphql/schema, apps/web-admin. -->
<!--   LINKS: M-API / V-M-API / M-GRAPHQL-SCHEMA / V-M-GRAPHQL-SCHEMA / M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRACE-WORKFLOW / V-M-GRACE-WORKFLOW. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Goal - Defines backend-only web-admin authentication scope. -->
<!--   Current Context - Summarizes the existing placeholder auth and user-domain boundary. -->
<!--   Key Decisions - Captures approved brainstorming choices. -->
<!--   Architecture - Defines admin identity, service, repository, sessions, middleware, and GraphQL surface. -->
<!--   Data Model And Config - Defines admin_users, seed env, and cookie/session settings. -->
<!--   GraphQL And HTTP Flow - Defines login, me, logout, create-admin, and protected GraphQL behavior. -->
<!--   Web-admin Transport Readiness - Defines the non-UI generated-client and cookie transport work needed for backend auth. -->
<!--   Error Handling And Security - Defines stable error mapping and safety constraints. -->
<!--   Operational Signals - Defines non-secret log markers required for auth observability. -->
<!--   Testing And Verification - Defines focused backend checks and downstream generated-client checks. -->
<!--   Out Of Scope - Defines UI and account-management exclusions for this wave. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added the approved backend-only web-admin authentication design. -->
<!-- END_CHANGE_SUMMARY -->

# Web-admin backend auth design

**Status:** Approved
**Date:** 2026-06-07

## Goal

Add the backend foundation for real `web-admin` authentication.

The work should let an admin log in, let the API identify the current admin, protect admin GraphQL operations, seed the first admin from environment configuration, and allow an existing admin to create another admin. Registration is not part of the product. Login UI is not part of this wave.

The implementation should replace the current bearer-shaped placeholder boundary with a real admin identity and session boundary while keeping the public REST users flow available as the template reference surface.

## Current Context

The repository already has:

- a Go API under `apps/api`;
- an admin-only GraphQL schema under `libs/graphql/schema`;
- a Vite `web-admin` app that talks to `/graphql`;
- PostgreSQL and Redis in local and test infrastructure;
- a `users` table and service used as the reference user CRUD domain;
- `apps/api/internal/middleware/auth.go`, which currently accepts a `Bearer` header but only stores the raw token string as `userID`.

The existing `users` domain is a demo/reference vertical slice, not an admin identity model. It should not be reused as the authentication account table.

## Key Decisions

| Decision          | Choice                                   | Rationale                                                                     |
| ----------------- | ---------------------------------------- | ----------------------------------------------------------------------------- |
| Admin identity    | Separate `admin_users` table             | Keeps web-admin access separate from the demo `users` domain.                 |
| Browser session   | `httpOnly` cookie with opaque session id | Avoids browser token storage and supports server-side logout.                 |
| Session backend   | Redis                                    | Redis already exists in the API runtime and gives TTL-backed revocation.      |
| Registration      | Not supported                            | First admin comes from env; later admins are created by authenticated admins. |
| Auth API shape    | GraphQL-first                            | `web-admin` already uses generated GraphQL contracts.                         |
| UI scope          | Excluded                                 | No login page or route guard UI is added.                                     |
| Web-admin client  | Cookie-ready transport only              | Backend auth needs `credentials: include`, but no UI behavior yet.            |
| Public REST users | Stays public                             | The public web reference REST flow should not inherit admin auth.             |

## Architecture

Add a new backend auth module inside `apps/api`, separate from `UserService`.

Target responsibilities:

- `AdminAuthService` owns login, current-admin lookup, first-admin seed, admin creation, and logout orchestration.
- `AdminRepository` owns PostgreSQL access to `admin_users`.
- `SessionStore` owns Redis-backed opaque sessions.
- Cookie helpers own session cookie set/clear behavior and production-safe attributes.
- Admin auth middleware owns session-cookie parsing, Redis lookup, active-admin validation, and context principal injection.
- GraphQL resolvers expose admin auth operations and guard admin-only operations.

The GraphQL resolver layer must have a narrow HTTP response bridge so `loginAdmin` and `logoutAdmin` can set or clear cookies. This can be done by adding the response writer or a small cookie sink to request context before gqlgen handles the request. The resolver should not know about raw router internals beyond that narrow cookie interface.

The existing placeholder middleware should be replaced or rewritten rather than layered on top of its current behavior.

## Data Model And Config

Add `admin_users` through a goose migration and sqlc queries.

Recommended table shape:

```sql
CREATE TABLE admin_users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  role VARCHAR(32) NOT NULL DEFAULT 'ADMIN',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

The first admin seed is configured by env-backed API config:

- `ADMIN_INITIAL_EMAIL`
- `ADMIN_INITIAL_PASSWORD`
- `ADMIN_INITIAL_NAME`

Seed behavior:

- runs after migrations and before the HTTP server accepts requests;
- counts existing `admin_users` first;
- creates the env-configured admin only when the table is empty;
- does not create, update, or rewrite any admin from env once at least one admin exists;
- never logs the password or raw hash;
- fails startup when the `admin_users` table is empty and the seed config is incomplete.

Session config:

- `ADMIN_SESSION_COOKIE_NAME`, default `web_admin_session`;
- `ADMIN_SESSION_TTL`, default `168h`;
- `ADMIN_SESSION_COOKIE_SECURE`, default derived from `SERVER_ENV == production`;
- `ADMIN_SESSION_SAME_SITE`, default `Lax`.

Password policy should reject empty and trivially short passwords. A minimum length of 12 characters is a reasonable template default unless implementation evidence shows existing tests or fixtures need a stricter contract. The same password policy applies to `ADMIN_INITIAL_PASSWORD` and `createAdmin`.

Email normalization should be enforced before persistence and lookup. The database contract should prevent case-only duplicate admin emails, either by storing normalized lowercase email values only or by adding an equivalent unique expression/index such as `UNIQUE (LOWER(email))`.

## GraphQL And HTTP Flow

Add admin auth schema types and operations to the admin GraphQL schema.

Recommended operation surface:

```graphql
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

Login flow:

1. Browser calls `loginAdmin`.
2. Resolver calls `AdminAuthService.Login`.
3. Service finds active admin by normalized email.
4. Service verifies bcrypt password hash.
5. Session store writes `admin_session:<hashed-session-id>` in Redis with TTL.
6. Cookie bridge sets `web_admin_session=<session-id>` with `HttpOnly`, `SameSite=Lax`, host-only scope, no `Domain`, `Path=/graphql`, `Max-Age`, and `Secure` in production.
7. Resolver returns `LoginAdminSuccess`.

Current admin flow:

- `me` reads the admin principal from context.
- Missing, expired, revoked, inactive, or malformed session returns no admin for `me`.
- Frontend UI can later use `me` as the route-guard bootstrap contract.

Logout flow:

- `logoutAdmin` deletes the Redis session if present.
- It clears the cookie.
- It is idempotent and returns success even when the incoming session is already absent or expired.

Create admin flow:

- `createAdmin` requires a valid active admin principal.
- It validates email, name, and password.
- It hashes the password with bcrypt.
- Duplicate email maps to `ValidationError(field: "email")`.
- Password validation maps to `ValidationError(field: "password")`.

Protected GraphQL behavior:

- `loginAdmin` is public.
- `logoutAdmin` is callable without a valid session so it can always clear the browser cookie.
- `me` may be called without a session and returns `null`.
- All existing user CRUD GraphQL queries/mutations and `createAdmin` require a valid admin principal.

The implementation may enforce this through gqlgen directives, resolver guards, or a small GraphQL operation guard. The selected mechanism must be testable and must not accidentally protect public REST `/api/users`.

## Web-admin Transport Readiness

Although login UI is out of scope, the backend-only wave must make the generated web-admin GraphQL transport cookie-session-ready.

Required non-UI client work:

- configure the `apps/web-admin` GraphQL client to send and receive cookies with credentialed requests, for example `credentials: "include"` on the `graphql-request` transport;
- remove, deprecate, or stop relying on the bearer-token helper as the admin auth path;
- add or update generated operation documents for `loginAdmin`, `logoutAdmin`, `me`, and `createAdmin` only as needed for type generation and future UI integration;
- add focused `graphql-client` tests proving the client is constructed with credentialed transport options.

This does not add a login page, route guard, token storage, localStorage session, or any visible UI.

## Error Handling And Security

Stable error mapping:

- missing or invalid admin session on protected mutation returns `AuthError`;
- missing or invalid admin session on protected query returns a stable GraphQL auth error unless the schema result type can carry `AuthError`;
- invalid login credentials return `AuthError` without revealing whether the email exists;
- inactive admin returns `AuthError`;
- invalid input returns `ValidationError`;
- duplicate admin email returns `ValidationError(field: "email")`;
- Redis or PostgreSQL failures return internal resolver errors and do not expose secrets.

Security constraints:

- never log passwords, raw session ids, raw cookies, password hashes, or credential payloads;
- store only opaque random session ids in browser cookies;
- store a hash or HMAC-derived session-key suffix in Redis instead of using the raw session id directly in Redis keys;
- store only minimal session data in Redis, such as admin id and expiry metadata;
- use cryptographically secure randomness for session ids;
- normalize emails consistently before lookup and create;
- hash admin passwords with bcrypt;
- clear cookies on logout using the same name/path/site attributes used to set them;
- keep public REST users endpoints outside admin auth middleware requirements.

Cookie-auth GraphQL must also have an explicit browser-origin boundary:

- credentialed `/graphql` CORS is limited to configured web-admin origins, not the public web origins used by `/api/users`;
- wildcard origins are invalid when credentials are enabled;
- protected and session-mutating GraphQL operations reject unsafe cross-origin browser requests through strict `Origin` or `Referer` allowlisting, or through an equivalent CSRF token defense;
- local development defaults may allow `http://localhost:3100` and `http://127.0.0.1:3100` for web-admin, but must not allow the public web origin to perform credentialed admin GraphQL calls;
- public REST `/api/users` keeps its existing public CORS behavior and remains outside the admin cookie requirement.

## Operational Signals

Implementation should emit stable, non-secret log markers for auth-critical paths:

- `[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN]` for seed start, create, skip, and incomplete-config failures;
- `[AdminAuth][login][BLOCK_VERIFY_CREDENTIALS]` for login attempts and credential failures without logging raw credentials;
- `[AdminAuth][session][BLOCK_VALIDATE_SESSION]` for session lookup, expiry, revocation, inactive-admin, and Redis failure paths without logging raw session ids or cookies;
- `[AdminAuth][logout][BLOCK_REVOKE_SESSION]` for logout and cookie-clear paths;
- `[AdminAuth][guard][BLOCK_AUTHORIZE_GRAPHQL]` for protected GraphQL allow/deny decisions without logging query bodies or secrets;
- `[AdminAuth][csrf][BLOCK_VALIDATE_ORIGIN]` for admin GraphQL origin/CSRF allow and deny decisions.

Focused tests or trace assertions should prove that passwords, password hashes, raw cookies, raw session ids, and credential payloads are not written to auth logs.

## Testing And Verification

Focused implementation tests should cover:

- config validation for admin seed and session settings;
- seed creates the first admin when absent;
- seed does not create or rewrite any admin once at least one admin exists;
- seed fails startup when the table is empty and seed env is incomplete;
- repository create/get-by-email/get-by-id and duplicate email mapping;
- service login success;
- wrong email and wrong password produce indistinguishable auth failures;
- inactive admin cannot log in or use a session;
- create admin requires an authenticated principal;
- create admin hashes password and maps duplicate email;
- session store create/read/delete/expiry behavior at the Redis boundary;
- cookie set/clear attributes in development and production modes;
- GraphQL login sets cookie through the response bridge;
- GraphQL logout clears cookie and deletes Redis session;
- `me` returns the current admin for a valid active session;
- `me` returns `null` for missing, expired, revoked, malformed, or inactive-admin sessions;
- unauthenticated protected GraphQL user operations are denied;
- authenticated protected GraphQL user operations still call the existing `UserService`;
- web-admin GraphQL transport sends credentialed cookie requests;
- web-admin browser/e2e setup logs in through `loginAdmin` using the browser context or Playwright API context, preserves the returned session cookie, and then runs protected CRUD checks with that cookie;
- credentialed admin GraphQL CORS allows configured web-admin origins and rejects public web or disallowed origins;
- GraphQL origin/CSRF checks reject unsafe protected or session-mutating cross-origin requests;
- public REST `/api/users` remains reachable without admin session;
- auth log assertions cover required markers and secret redaction.

Focused commands:

- `bunx nx run api:codegen`
- `bunx nx run codegen:validate`
- `bunx nx test api`
- `bunx nx build api`
- `bunx nx run graphql:validate` when schema validation is available in the current target graph
- `bunx nx run web-admin:codegen` after schema changes, even without UI work
- `bunx nx run web-admin:typecheck` if generated web-admin types change
- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`
- `grace lint --path .`

Broader `bun run test:coverage`, `bun run test:e2e`, and `bun run verify:coverage` should be reserved for final closeout or a release handoff unless implementation evidence shows the auth change has crossed into browser UI behavior.

Generated and coverage contracts:

- new sqlc output such as `apps/api/internal/repository/postgres/generated/admins.sql.go` must be committed with API codegen;
- `tools/coverage/coverage.config.json` and `docs/verification-plan.xml` must add an allowlist entry for generated admin sqlc query output when it is created;
- every new generated allowlist entry must name a replacement gate, including API codegen, API build, and an admin repository integration test against the goose-created `monorepo_test` schema;
- `bunx nx run codegen:validate` must be part of final focused verification for this wave.

## GRACE Updates Expected During Implementation

Implementation should refresh:

- `docs/requirements.xml`: replace the current product-auth placeholder for admin GraphQL with a real web-admin backend auth requirement while preserving downstream owner/tenant authorization warnings for future product-specific resources.
- `docs/development-plan.xml`: add or update API/auth, GraphQL schema, and web-admin generated-client contracts.
- `docs/knowledge-graph.xml`: add graph facts for admin auth repository, service, session store, middleware, and GraphQL schema links.
- `docs/verification-plan.xml`: add auth scenarios, focused commands, and negative access-control tests.
- `docs/operational-packets.xml`: only if execution packets need a new auth-specific packet shape.

Meaningfully edited governed source, test, config, schema, query, migration, Docker, tooling, or docs files must carry or update file-local GRACE markup.

## Out Of Scope

This wave does not include:

- login page UI;
- frontend route guards;
- admin list/edit/deactivate UI;
- password reset;
- email verification;
- OAuth or SSO;
- multi-role permissions beyond the initial `ADMIN` role string;
- tenant or owner scoping for product-specific resources;
- making public REST users endpoints require admin auth.
