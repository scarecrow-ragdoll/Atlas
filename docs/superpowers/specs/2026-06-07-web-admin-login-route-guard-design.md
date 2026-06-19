<!-- FILE: docs/superpowers/specs/2026-06-07-web-admin-login-route-guard-design.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Capture the approved design for the web-admin login page and protected frontend routes. -->
<!--   SCOPE: Design-level frontend auth architecture, route guard behavior, login/logout UX, data flow, error handling, testing, and GRACE update expectations; excludes implementation code and backend auth redesign. -->
<!--   DEPENDS: docs/superpowers/specs/2026-06-07-web-admin-backend-auth-design.md, docs/superpowers/plans/2026-06-07-web-admin-backend-auth.md, apps/web-admin, libs/graphql/schema/admin_auth.graphql, apps/api admin auth GraphQL. -->
<!--   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA / V-M-GRAPHQL-SCHEMA / M-GRACE-WORKFLOW / V-M-GRACE-WORKFLOW. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Goal - Defines frontend login and route protection scope. -->
<!--   Current Context - Summarizes existing backend auth, GraphQL transport, router, shell, and e2e state. -->
<!--   Key Decisions - Captures approved route, redirect, session, and UX choices. -->
<!--   Architecture - Defines auth slice, provider, protected layout, login route, shell user wiring, and logout action. -->
<!--   Data Flow - Defines me/login/logout query and mutation sequencing with httpOnly cookie sessions. -->
<!--   Route Behavior - Defines login, protected-route, already-authenticated, unknown-route, and return-to handling. -->
<!--   Error Handling And Security - Defines user-visible errors, network behavior, and open-redirect prevention. -->
<!--   Testing And Verification - Defines focused unit, integration, e2e, codegen, typecheck, and lint checks. -->
<!--   Out Of Scope - Defines registration, password reset, admin management UI, and backend auth exclusions. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added the approved web-admin login page and route guard design. -->
<!-- END_CHANGE_SUMMARY -->

# Web-admin login and route guard design

**Status:** Approved
**Date:** 2026-06-07

## Goal

Add the visible web-admin authentication experience on top of the existing backend auth foundation.

The work should provide a `/login` page, protect every other web-admin route, preserve the user's intended destination through a safe `from` redirect, show the real current admin in the sidebar, and let the admin log out. The frontend must keep using the backend httpOnly cookie session; it must not store access tokens in browser storage.

## Current Context

The repository already has:

- `loginAdmin`, `logoutAdmin`, `me`, and `createAdmin` GraphQL operations under `apps/web-admin/src/entities/admin-auth/api/adminAuth.graphql`;
- generated web-admin GraphQL types for the admin auth operations;
- `graphql-request` configured with `credentials: "include"` in `apps/web-admin/src/shared/api/graphql-client.ts`;
- a Vite React router in `apps/web-admin/src/App.tsx`;
- one `AdminLayout` that wraps the current home, users, user detail, and UI-kit pages;
- a sidebar shell with a static placeholder user menu in `apps/web-admin/src/shared/ui/layout/nav-user.tsx`;
- Playwright e2e helpers that currently install an admin session cookie directly before protected UI flows.

The backend-auth wave intentionally excluded a login UI and route guards. This wave finishes the frontend connection to that backend contract.

## Key Decisions

| Decision            | Choice                                   | Rationale                                                                        |
| ------------------- | ---------------------------------------- | -------------------------------------------------------------------------------- |
| Session storage     | Backend httpOnly cookie only             | Keeps credentials out of frontend storage and matches the backend auth design.   |
| Login route         | `/login` outside the admin sidebar shell | Login should not render protected navigation before identity is known.           |
| Protected routes    | Guard all routes except `/login`         | New admin pages should be protected by default.                                  |
| Post-login redirect | Return to safe `from` path, fallback `/` | Preserves direct links and avoids forcing users to restart navigation.           |
| Auth state source   | `CurrentAdmin` / `me` query              | The backend is the source of truth for session validity.                         |
| Sidebar user        | Real `me` admin principal                | Removes the placeholder `Developer` identity from authenticated pages.           |
| Logout              | Sidebar user menu action                 | Keeps the session lifecycle in the expected account menu location.               |
| Registration        | Not supported                            | First admin and admin creation are backend-admin flows, not public registration. |

## Architecture

Add a small frontend auth slice under `apps/web-admin/src/entities/admin-auth` beside the existing operation documents.

Target responsibilities:

- Auth API helpers own `CurrentAdmin`, `LoginAdmin`, and `LogoutAdmin` request calls through the existing credentialed `graphqlClient`.
- Auth result helpers normalize GraphQL union responses into success or user-visible error states.
- `AuthProvider` owns the current admin query, exposes `admin`, `isLoading`, `isAuthenticated`, `login`, `logout`, and `refetchCurrentAdmin`, and keeps React Query as the backing cache.
- `ProtectedAdminLayout` checks auth state before rendering `AdminLayout`.
- `LoginPage` renders the login form outside the sidebar shell.
- `AdminLayout` accepts the authenticated admin display data and passes it into the existing `AdminAppShell`.
- `NavUser` replaces placeholder account menu items with a real logout action.

The shared UI boundary remains intact: pages and app-layer components should import visual primitives and shell components from the approved bare `@shared/ui` surface.

`AuthProvider` placement must respect the current provider/router shape: `Providers` wraps `App`, while `BrowserRouter` lives inside `App.tsx`. If the auth provider uses router hooks, mount it under `BrowserRouter`; otherwise keep the provider navigation-free and let route components perform navigation.

## Data Flow

`CurrentAdmin` is the bootstrap request for frontend auth state.

On protected route entry:

1. `ProtectedAdminLayout` reads auth state from `AuthProvider`.
2. While `CurrentAdmin` is pending, it renders a compact loading state.
3. If `me` returns an admin, the route renders inside `AdminLayout`.
4. If `me` returns `null`, the route redirects to `/login?from=<current-path>`.

The captured `from` value should preserve the same-app `pathname`, `search`, and `hash` so future filtered or anchored admin pages can return to the exact requested view.

On login:

1. `LoginPage` submits email and password through `loginAdmin`.
2. The API sets the httpOnly session cookie.
3. The frontend invalidates or refetches `CurrentAdmin`.
4. If `me` confirms an admin, the app navigates to the safe `from` path or `/`.

On logout:

1. The sidebar menu calls `logoutAdmin`.
2. The API clears the cookie and revokes the Redis session.
3. The frontend clears the auth-related React Query state.
4. The app navigates to `/login`.

## Route Behavior

The route table should be restructured around public and protected layouts:

- `/login` renders `LoginPage` and does not use `AdminLayout`.
- `/`, `/users`, `/users/:id`, `/ui-kit`, and future non-login routes render through `ProtectedAdminLayout`.
- Unknown routes still redirect to `/`, which is protected and therefore redirects unauthenticated users to login.
- If an authenticated admin opens `/login`, the app redirects to the safe `from` path or `/`.

The `from` value must accept only same-app relative paths beginning with `/`. It must reject absolute URLs, protocol-relative URLs, malformed paths, and `/login` loops. Rejected values fall back to `/`.

## Login UX

The login page should be quiet and operational, consistent with the current admin UI-kit:

- email input;
- password input;
- submit button with pending state;
- inline error area for auth, validation, or network failures;
- no registration prompt;
- no marketing hero;
- no token or cookie debug text.

The page can use existing `@shared/ui` primitives such as `Card`, `Input`, `Label`, `Button`, and `Alert`. The copy should be short and admin-oriented.

## Shell User And Logout UX

Authenticated pages should show the current admin in the sidebar footer:

- display `name`, `email`, and derived initials from `me`;
- remove disabled account/settings placeholder labels;
- add a `Logout` menu item with a standard icon from the existing UI-kit/icon surface;
- disable or show pending state while logout is in flight;
- after logout, redirect to `/login` and prevent stale protected content from flashing.

The admin shell should continue to resolve route navigation, breadcrumbs, teams, and reference items from app-owned navigation metadata.

## Error Handling And Security

Login errors:

- `AuthError` displays a generic invalid-credentials message from the backend result.
- `ValidationError` displays a field-qualified message, such as `email: ...` or `password: ...`.
- network or unexpected GraphQL errors display a stable fallback such as `Unable to sign in. Try again after the API is available.`

Guard errors:

- missing session or `me === null` redirects to login;
- an unhandled `CurrentAdmin` network failure should not expose protected page data, so the safe default is login redirect with no secret details;
- protected page queries that later fail with auth errors should rely on the guard refresh/logout path rather than adding page-specific auth handling.

Security constraints:

- do not store passwords, sessions, access tokens, or raw cookies in localStorage, sessionStorage, or React Query;
- do not log credential payloads;
- keep `from` redirect same-app and relative only;
- do not expose production admin seed credentials in frontend bundles, application code, logs, or test source. Isolated Playwright e2e may use test-only admin credentials supplied through the Playwright server/test environment.

## Testing And Verification

Focused unit and integration tests:

- route guard redirects unauthenticated `/users` to `/login?from=/users`;
- route guard renders protected content when `CurrentAdmin` returns an admin;
- route guard handles `CurrentAdmin` rejection and `me === null` without rendering protected page content before redirecting to `/login`;
- login page redirects an already-authenticated admin away from `/login`;
- login success navigates to the safe `from` path;
- login success preserves a safe `from` path with search/hash;
- login success falls back to `/` for absolute URLs, protocol-relative URLs, malformed paths, `/login` loops, and absent `from`;
- login page renders `AuthError`, `ValidationError`, and network failure states;
- sidebar user menu renders real admin name/email and calls logout;
- logout clears auth state and navigates to `/login`;
- `resolveAdminShellState` or adjacent shell mapping tests cover real-admin display wiring if app navigation ownership changes.

Focused e2e coverage:

- unauthenticated visit to `/users` redirects to `/login`;
- real browser login through `/login` reaches the original `/users` destination;
- existing users create/list/detail flow runs after the real login UI path;
- logout returns to `/login`;
- revisiting `/users` after logout redirects to `/login` again.

Required focused commands:

```bash
bunx nx run web-admin:codegen
bunx nx run web-admin:typecheck
bunx nx test web-admin
bunx nx lint web-admin
bunx nx build web-admin
bunx nx run web-admin:e2e
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Run broader gates only at closeout or when shared generated/schema surfaces change beyond the existing admin-auth GraphQL operations.

## GRACE Updates

Meaningful implementation must update:

- `docs/requirements.xml` with the frontend login and protected-route acceptance criteria;
- `docs/development-plan.xml` with the web-admin auth UI module responsibilities;
- `docs/knowledge-graph.xml` with auth slice, protected layout, login page, and sidebar logout links;
- `docs/verification-plan.xml` with unit and e2e scenarios for login, guard, return-to, and logout;
- file-local GRACE markup for new or meaningfully edited frontend source, test, and e2e files.

## Out Of Scope

This design does not add:

- public registration;
- forgot-password or password reset;
- admin list/edit/deactivate UI;
- role management UI;
- MFA, OAuth, or SSO;
- backend auth data model changes;
- product-specific tenant or owner authorization.
