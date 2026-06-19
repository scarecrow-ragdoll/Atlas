<!-- FILE: docs/superpowers/specs/2026-06-05-swap-web-admin-vite-web-next-design.md -->
<!-- VERSION: 1.0.1 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Capture the approved design for swapping web-admin to a Vite client app and public web to a Next.js App Router app. -->
<!--   SCOPE: Design-level architecture, component boundaries, data flow, error handling, testing, migration constraints, and GRACE updates; excludes implementation code. -->
<!--   DEPENDS: docs/requirements.xml, docs/technology.xml, docs/development-plan.xml, docs/knowledge-graph.xml, docs/verification-plan.xml, apps/web-admin, apps/web, libs/graphql/schema, tools/codegen/codegen.ts, tools/coverage/preflight.mjs, docker/web.Dockerfile, .gitlab-ci.yml, deploy/dokploy/docker-compose.template.yml. -->
<!--   LINKS: M-WEB-ADMIN / M-WEB / M-GRAPHQL-SCHEMA / M-API / V-M-WEB-ADMIN / V-M-WEB. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Goal - Defines the approved target framework and transport ownership. -->
<!--   Current Context - Summarizes the existing inverted app roles. -->
<!--   Key Decisions - Captures the approved choices from brainstorming. -->
<!--   Architecture - Defines the target app boundaries while preserving project names. -->
<!--   Components And Data Flow - Defines admin GraphQL SPA flow and public Next REST flow. -->
<!--   Error Handling - Preserves GraphQL union and REST error-envelope behavior. -->
<!--   Deployment Ownership - Defines Docker, CI image, and Dokploy ownership after the framework swap. -->
<!--   Migration Plan - Lists config, dependency, route, codegen, and docs surfaces to update. -->
<!--   Testing And Verification - Defines focused and final gates for the swap. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.1 - Added admin Vite env, deployment ownership, and coverage preflight requirements from subagent review. -->
<!-- END_CHANGE_SUMMARY -->

# Swap web-admin Vite and public web Next design

**Status:** Approved
**Date:** 2026-06-05

## Goal

Swap the frontend framework ownership while preserving the public project names and transport split:

- `apps/web-admin` remains the `web-admin` Nx project and becomes a Vite + React client-side admin app.
- `apps/web` remains the `web` Nx project and becomes a Next.js App Router public web app.
- Admin data access remains GraphQL through `graphql-request`, GraphQL documents, and generated client types.
- Public web data access remains REST through `/api/users`; it must not receive GraphQL codegen or GraphQL transport dependencies.

The goal is to make the template reflect the intended product shape: admin is a client application, while the public web surface uses Next.js conventions.

## Current Context

The repository currently has the opposite framework assignment:

- `apps/web-admin` is a Next.js app using generated GraphQL client types and React Query.
- `apps/web` is a Vite app using a typed REST users client and React Query.
- `tools/codegen/codegen.ts` generates GraphQL operation types into `apps/web-admin/src/shared/api/generated/types.ts`.
- GRACE artifacts describe `M-WEB-ADMIN` as `NextWebAdminApp` and `M-WEB` as `VitePublicWebApp`.

The existing transport split is correct and should remain: GraphQL is admin-only, and public web uses REST.

## Key Decisions

| Decision             | Choice                                                 | Rationale                                                                          |
| -------------------- | ------------------------------------------------------ | ---------------------------------------------------------------------------------- |
| Project names        | Preserve `web-admin` and `web`                         | Keeps Nx commands, CI references, e2e artifact names, and product meaning stable.  |
| Admin framework      | Vite + React SPA                                       | Admin is always a client application in this template.                             |
| Admin routing        | `react-router` with `/`, `/users`, and `/users/:id`    | Preserves useful URL semantics without server routing.                             |
| Admin transport      | GraphQL + generated operation types                    | Keeps the admin GraphQL contract and codegen examples intact.                      |
| Public web framework | Next.js App Router                                     | Makes `web` the public Next surface for downstream products.                       |
| Public web transport | REST + React Query/client components where interactive | Keeps the existing public REST contract and avoids GraphQL drift into `web`.       |
| Admin GraphQL env    | Vite browser env, for example `VITE_GRAPHQL_API_URL`   | Avoids keeping Next-only public env semantics in the Vite admin app.               |
| Public REST env      | Next browser/server env for REST API base URLs         | Replaces the current Vite-only public web env contract.                            |
| Deployable web image | Retarget the existing `web` image to the new Next app  | Keeps the release image name stable while making public web the deployed frontend. |
| Admin deployment     | Define or explicitly defer static Vite admin deploy    | Prevents the Vite admin from silently inheriting the old Next standalone contract. |
| Migration style      | True framework swap inside existing app directories    | Avoids renaming projects or changing external command contracts.                   |

## Architecture

`apps/web-admin` should become a Vite SPA with this shape:

```text
apps/web-admin
  index.html
  vite.config.ts
  src/main.tsx
  src/App.tsx
  src/app/providers.tsx
  src/pages/home
  src/pages/users
  src/pages/user-detail
  src/entities/user/api/*.graphql
  src/shared/api/graphql-client.ts
  src/shared/api/generated/types.ts
```

The exact page folder names can follow the repository's local convention during implementation. The important boundary is that route ownership moves from Next App Router files to client-side React Router routes.

`apps/web` should become a Next.js App Router app with this shape:

```text
apps/web
  next.config.js
  next-env.d.ts
  app/layout.tsx
  app/page.tsx
  app/users/page.tsx or app/page.tsx as the users flow
  src/app/providers.tsx
  src/shared/api/users.ts
  src/shared/config.ts
  src/styles.css or app/globals.css
```

The public web may keep a compact single-page users flow if that remains the clearest reference example. The implementation should not add product-specific routes beyond the reference users flow.

The Go API remains unchanged at the architectural level:

- `/graphql` serves the admin GraphQL schema and generated API resolvers.
- `/api/users` serves public REST users endpoints.
- GraphQL resolvers and REST handlers continue to call the shared `UserService`.

## Components And Data Flow

### `web-admin` Vite GraphQL SPA

The admin entrypoint mounts React into the Vite `index.html` root. It wraps the app in `QueryClientProvider` and a browser router.

Routes:

- `/` renders a small admin home or redirects/links to `/users`.
- `/users` renders the admin users list and create-user form.
- `/users/:id` renders the admin user detail view.

The admin GraphQL client remains in `src/shared/api/graphql-client.ts`. GraphQL documents remain under entity or feature API folders and continue to feed generated operation types through `tools/codegen/codegen.ts`.

Where the current Next pages use inline GraphQL strings and manual interfaces, the implementation should prefer generated operation types and imported GraphQL documents when practical. The design does not require broad architecture refactors beyond the swap, but it should avoid preserving duplicated handwritten GraphQL response types when generated types are already available.

The admin config must move away from `process.env.NEXT_PUBLIC_API_URL` and other Next-only public env semantics. The Vite admin app should use an `import.meta.env`-backed browser config such as `VITE_GRAPHQL_API_URL`, defaulting to the local GraphQL endpoint when no env value is supplied. Admin config unit tests, `web-admin:e2e` Playwright env injection, and any local compose or deploy env references must use the same admin GraphQL env contract.

### `web` Next.js REST public app

The public web app should use Next.js App Router for page/layout ownership. It keeps the REST users client as the only transport client for the reference user flow.

The preferred data flow is server-first public page plus client-side interaction:

```text
Next page/layout
  -> REST list fetch for initial public users data where practical
  -> client component for create/select/refetch behavior
  -> REST users client
  -> Go API /api/users
```

If implementation constraints make initial server fetch disproportionately expensive for the template slice, the plan may use a client-side React Query flow inside the Next App Router shell. Even then, `web` must still be a real Next App Router application rather than a Vite shell copied into a Next project.

The public config should move away from `VITE_API_BASE_URL`. A suitable replacement is `NEXT_PUBLIC_API_BASE_URL` for browser-visible REST calls, with a safe default matching local API development. If server-side REST fetches are used, the implementation may introduce a server-only API base URL while keeping browser config explicit.

## Error Handling

Admin GraphQL behavior should preserve current user-visible semantics:

- validation union results render field-specific form errors;
- auth union results render stable admin errors;
- request failures surface without corrupting React Query state;
- successful mutations invalidate users queries and clear the form.

Public REST behavior should preserve current `ApiError` semantics:

- REST error envelopes expose `code`, `message`, optional `field`, and `status`;
- duplicate email and validation failures show useful form errors;
- missing or failed list loads render a load error state;
- empty lists render an empty state instead of a broken page.

The swap should not change API error contracts.

## Deployment Ownership

The current deployable `web` image is framework-specific: `docker/web.Dockerfile` builds `apps/web-admin` as a Next standalone app, CI publishes that Dockerfile as the `web` image, and the Dokploy template runs that image under the `web-admin` service.

After the swap, the existing `web` image should be retargeted to the new Next public app in `apps/web`. This preserves the external image name while aligning it with the public web product surface.

The Vite admin deployment must not silently reuse the old Next standalone image contract. The implementation must either add an explicit static-serving deployment path for `apps/web-admin`, including Docker, CI, Dokploy, and env ownership, or state that admin deployment is out of scope for the implementation wave and block release/deploy handoff until a follow-up deployment design exists.

Any CI or deploy metadata that describes the `web` image or Dokploy `web-admin` service must be updated to match the chosen deployment ownership before release or deploy gates are considered complete.

## Migration Plan

Implementation should update these surfaces in a controlled order:

1. Update `apps/web-admin` package dependencies and scripts from Next to Vite, adding `react-router` and preserving React Query, GraphQL, `graphql-request`, and codegen dependencies.
2. Update `apps/web` package dependencies and scripts from Vite to Next, preserving React Query and REST client tests while removing Vite-only runtime dependencies.
3. Move or rewrite app shell files:
   - `web-admin`: `index.html`, `vite.config.ts`, `src/main.tsx`, `src/App.tsx`, route components, and provider wiring.
   - `web`: `app/layout.tsx`, `app/page.tsx`, client users component, Next config, and Next env typing.
4. Migrate env/config ownership:
   - `web-admin`: replace `NEXT_PUBLIC_API_URL` with a Vite-compatible admin GraphQL env contract and update config tests plus Playwright env injection.
   - `web`: replace `VITE_API_BASE_URL` with the Next public REST env contract and optional server-only REST env.
5. Keep Nx project names and target names stable. Adjust target commands only where scripts change.
6. Keep GraphQL codegen output under `apps/web-admin/src/shared/api/generated/types.ts`.
7. Ensure `apps/web` has no GraphQL documents, generated GraphQL client types, or GraphQL transport dependencies.
8. Update Playwright configs and preflight behavior so `web-admin:e2e` starts Vite and `web:e2e` starts Next.
9. Update `tools/coverage/preflight.mjs` so required-file checks match the new Vite `web-admin` and Next `web` shapes.
10. Update coverage allowlists and replacement gates for Vite bootstrap in `web-admin`, Next bootstrap in `web`, and generated GraphQL output in `web-admin`.
11. Update deployment surfaces according to the Deployment Ownership section:

- `docker/web.Dockerfile`
- `.gitlab-ci.yml`
- `deploy/dokploy/docker-compose.template.yml`
- related CI helper metadata when it names the web image or service shape.

12. Update GRACE docs:

- `docs/requirements.xml`
- `docs/technology.xml`
- `docs/development-plan.xml`
- `docs/knowledge-graph.xml`
- `docs/verification-plan.xml`
- `docs/operational-packets.xml` only if packet references need this framework ownership change.

Implementation should not rename `apps/web-admin`, `apps/web`, Nx project names, or root commands as part of this change.

## Testing And Verification

Focused checks should prove the changed surfaces before broad gates:

- `bunx nx test web-admin`
- `bunx nx run web-admin:codegen`
- `bunx nx run web-admin:typecheck`
- `bunx nx build web-admin`
- `bunx nx test web`
- `bunx nx run web:typecheck`
- `bunx nx build web`
- `bunx nx run graphql:validate`
- `bunx nx run api:codegen`
- `bunx nx run codegen:validate`

E2e and coverage checks should be updated to match the new framework ownership:

- `bunx nx run web-admin:e2e` proves the admin GraphQL browser flow against the Vite admin server.
- `bunx nx run web:e2e` proves the public REST browser flow against the Next public web server.
- `bun run test:coverage` enforces updated TypeScript coverage ownership.
- `tools/coverage/preflight.mjs` required-file checks match the swapped app shapes before the coverage gate runs.
- `bun run verify:coverage` remains the final handoff gate when Docker and Playwright are available.

GRACE integrity checks before closeout:

- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`
- `grace lint --path .`

## Risks And Stop Conditions

Stop and replan if:

- GraphQL codegen output moves into `apps/web`.
- Public web begins importing `graphql-request` or GraphQL generated types.
- Admin URLs cannot preserve `/users` and `/users/:id` without a clear replacement.
- Admin Vite env cannot provide the GraphQL endpoint through a single tested browser config contract.
- The deployable `web` image cannot be retargeted to `apps/web` without a separate CI/deploy design.
- Vite admin deployment is required for release but has no explicit static-serving plan.
- `tools/coverage/preflight.mjs` cannot express the swapped required-file shape without weakening coverage preflight.
- E2e configs become ambiguous about which dev server they start.
- Coverage allowlists are broadened instead of moved to the correct generated or bootstrap files.
- CI or deployment expects framework-specific behavior that requires a separate deployment design.

## Approval

Approved direction:

- `web-admin` becomes Vite + React SPA + React Router + GraphQL/codegen.
- `web` becomes Next.js App Router + REST.
- Project names and external Nx commands remain stable.
