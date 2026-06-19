# Web-admin Login Route Guard Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a real `/login` page, protect every other web-admin route, wire the sidebar to the current backend admin session, and verify login/logout through the browser.

**Architecture:** Keep the API as the session source of truth through the existing httpOnly cookie GraphQL auth contract. Add a focused frontend auth slice for request helpers, principal mapping, safe return-to handling, and auth context. Mount `/login` outside the sidebar shell and render all other routes through `ProtectedAdminLayout`, which gates content on `CurrentAdmin` before `AdminLayout` renders.

**Tech Stack:** Vite, React 19, React Router 7, React Query 5, graphql-request, GraphQL Codegen, Vitest, Testing Library, Playwright, Nx, GRACE XML.

---

<!-- FILE: docs/superpowers/plans/2026-06-07-web-admin-login-route-guard.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Provide the task-by-task implementation plan for the web-admin login page and protected frontend routes. -->
<!--   SCOPE: Covers auth frontend helpers, provider, protected route layout, login page, sidebar logout, e2e login/logout, GRACE docs, focused verification, and commit boundaries; excludes implementation performed by this document. -->
<!--   DEPENDS: docs/superpowers/specs/2026-06-07-web-admin-login-route-guard-design.md, apps/web-admin, libs/graphql/schema/admin_auth.graphql, apps/api admin auth GraphQL. -->
<!--   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA / V-M-GRAPHQL-SCHEMA / M-GRACE-WORKFLOW / V-M-GRACE-WORKFLOW. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Source Spec - Anchors the approved and subagent-reviewed login route guard design. -->
<!--   Scope Check - Confirms this is one coupled frontend-auth plan. -->
<!--   File Structure - Lists planned creates and modifications with ownership boundaries. -->
<!--   Execution Discipline - Defines TDD, generated-code, GRACE, dirty-base, and commit rules. -->
<!--   Tasks - Provides ordered implementation, test, docs, e2e, and verification steps. -->
<!--   Self-Review - Records spec coverage, red-flag scan, and type consistency checks. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added the implementation plan for web-admin login and route protection. -->
<!-- END_CHANGE_SUMMARY -->

## Source Spec

- Design: `docs/superpowers/specs/2026-06-07-web-admin-login-route-guard-design.md`
- Design commits:
  - `bed10c2 docs: design web-admin login guard`
  - `b9a70bd docs: refine web-admin login guard spec`
- Required review-loop verdicts:
  - Intent reviewer: `APPROVE`
  - Feasibility reviewer: `APPROVE`
  - Verification reviewer: `APPROVE`
  - UX/product reviewer: `APPROVE`
  - Security/privacy reviewer: `APPROVE`

Approved decisions:

- `/login` is public and outside the sidebar shell.
- Every other route is protected by default.
- `CurrentAdmin` / `me` is the frontend auth source of truth.
- Login and logout use only the backend httpOnly cookie session.
- No access token, session id, raw cookie, or password is stored in browser storage.
- Successful login returns to a same-app safe `from` path, preserving `pathname + search + hash`, or falls back to `/`.
- Unsafe `from` values include absolute URLs, protocol-relative URLs, malformed paths, `/login`, and `/login` loops.
- Sidebar footer shows the real admin principal and exposes logout.

## Scope Check

This is one implementation plan. The login page, protected layout, auth provider, sidebar logout, and e2e tests are coupled because the feature is not correct unless the browser can move from unauthenticated protected URL to `/login`, authenticate through `loginAdmin`, render protected content through `me`, and revoke the session with `logoutAdmin`.

The plan does not include:

- backend auth schema or session redesign;
- public registration;
- forgot-password or password reset;
- admin list/edit/deactivate UI;
- role management UI;
- MFA, OAuth, or SSO;
- public web changes under `apps/web/src/**`;
- product-specific tenant or owner authorization.

## File Structure

### Create

- `apps/web-admin/src/entities/admin-auth/model.ts` - frontend admin principal types, safe return-to parsing, initials derivation, and shell user mapping.
- `apps/web-admin/src/entities/admin-auth/model.test.ts` - unit coverage for initials and safe return-to rules.
- `apps/web-admin/src/entities/admin-auth/api/loginAdmin.graphql` - single-operation login document for runtime request and codegen.
- `apps/web-admin/src/entities/admin-auth/api/logoutAdmin.graphql` - single-operation logout document for runtime request and codegen.
- `apps/web-admin/src/entities/admin-auth/api/currentAdmin.graphql` - single-operation current-admin document for runtime request and codegen.
- `apps/web-admin/src/entities/admin-auth/api/createAdmin.graphql` - single-operation create-admin document for codegen parity with the backend auth contract.
- `apps/web-admin/src/entities/admin-auth/client.ts` - typed GraphQL request helpers for `CurrentAdmin`, `LoginAdmin`, and `LogoutAdmin`.
- `apps/web-admin/src/entities/admin-auth/client.test.ts` - unit coverage for request helper documents and union error normalization.
- `apps/web-admin/src/entities/admin-auth/provider.tsx` - React Query-backed auth context with navigation-free `login`, `logout`, `admin`, and loading state.
- `apps/web-admin/src/entities/admin-auth/provider.test.tsx` - provider/hook coverage for current-admin, login success/error, and logout cache clearing.
- `apps/web-admin/src/app/protected-admin-layout.tsx` - protected route layout that redirects unauthenticated users before rendering `AdminLayout`.
- `apps/web-admin/src/pages/login-page.tsx` - quiet operational login page using `@shared/ui`.
- `apps/web-admin/src/pages/login-page.test.tsx` - login page coverage for success redirects, unsafe `from`, already-authenticated redirect, and visible errors.

### Modify

- `apps/web-admin/src/App.tsx` - restructure route table around `/login` and protected app routes.
- `apps/web-admin/src/App.test.tsx` - update route tests for auth-aware routing and no protected-content flash.
- `apps/web-admin/src/app/admin-layout.tsx` - accept authenticated admin and logout state, then pass shell user/logout props into `AdminAppShell`.
- `apps/web-admin/src/shared/ui/layout/admin-shell-types.ts` - add optional logout action state to the shell user menu contract.
- `apps/web-admin/src/shared/ui/layout/admin-app-shell.tsx` - pass logout action props through to `AppSidebar`.
- `apps/web-admin/src/shared/ui/layout/app-sidebar.tsx` - pass logout action props through to `NavUser`.
- `apps/web-admin/src/shared/ui/layout/nav-user.tsx` - replace static account/settings entries with a real logout menu item.
- `apps/web-admin/src/shared/ui/layout/admin-shell-types.ts` tests are not needed because it is type-only and covered by consumer tests.
- `apps/web-admin/src/shared/ui/layout/admin-shell.test.tsx` or nearest existing shell test - add logout rendering coverage if current shell tests do not cover `NavUser`.
- `apps/web-admin/src/entities/admin-auth/api/adminAuth.graphql` - remove after splitting operations into the single-operation files above so codegen does not see duplicate operation names.
- `apps/web-admin/e2e/helpers.ts` - expose test admin credentials and a browser login helper; keep API-bound admin login for direct GraphQL seed helpers.
- `apps/web-admin/e2e/users-flow.spec.ts` - replace per-test cookie install with real browser login for UI flows and add logout/redirect assertions.
- `docs/requirements.xml` - add frontend login/protected-route acceptance criteria under the web-admin auth use case.
- `docs/development-plan.xml` - add web-admin auth UI responsibilities under `M-WEB-ADMIN`.
- `docs/knowledge-graph.xml` - add auth slice, protected layout, login page, and sidebar logout paths/exports.
- `docs/verification-plan.xml` - add unit and e2e scenarios plus closeout commands for login, guard, return-to, logout, and GRACE lint.

### Do Not Modify

- `apps/web/src/**` - public web stays REST-only and public.
- `apps/api/**` - backend auth is the dependency for this plan, not the implementation target.
- `apps/web-admin/src/shared/api/graphql-client.ts` - it already sends `credentials: "include"`; only revisit if a test proves drift.

## Execution Discipline

- Begin from the current backend-auth branch state. If the admin auth operation documents are still in one aggregate `adminAuth.graphql`, split them into one operation per file before Task 1 so runtime `graphqlClient.request(document, variables)` calls match the existing web-admin API pattern.
- If generated `LoginAdminMutation` / `CurrentAdminQuery` types are missing after the split, run `bunx nx run web-admin:codegen` before Task 1.
- Do not edit production code before the task's failing test exists and has been run.
- Keep `AuthProvider` navigation-free. Navigation belongs in `ProtectedAdminLayout`, `LoginPage`, and logout UI components under `BrowserRouter`.
- Import UI primitives and shell compositions from bare `@shared/ui` in pages/app components.
- Do not put session ids, raw cookies, passwords, or admin seed credentials in browser storage or logs.
- Preserve existing backend-auth dirty changes. Stage and commit only files touched by each frontend task.
- Task 0 may commit only operation-document split and generated-type drift before TDD production work begins, because `client.ts` must not import uncommitted operation files.
- Every new or meaningfully edited governed source, test, e2e, and durable docs file must carry or update file-local GRACE markup.

## Task 0: Preflight And Generated Auth Contract

**Files:**

- Inspect: `apps/web-admin/src/entities/admin-auth/api/adminAuth.graphql`
- Create: `apps/web-admin/src/entities/admin-auth/api/loginAdmin.graphql`
- Create: `apps/web-admin/src/entities/admin-auth/api/logoutAdmin.graphql`
- Create: `apps/web-admin/src/entities/admin-auth/api/currentAdmin.graphql`
- Create: `apps/web-admin/src/entities/admin-auth/api/createAdmin.graphql`
- Delete: `apps/web-admin/src/entities/admin-auth/api/adminAuth.graphql`
- Inspect: `apps/web-admin/src/shared/api/generated/types.ts`
- Inspect: `apps/web-admin/src/shared/api/graphql-client.ts`
- Inspect: `git status --short`

- [ ] **Step 1: Inspect current dirty base**

Run:

```bash
git status --short -- apps/web-admin docs/requirements.xml docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml
```

Expected: existing backend-auth changes may be present. Record any unrelated frontend changes before editing the same files.

- [ ] **Step 2: Split aggregate auth operation documents if needed**

If `apps/web-admin/src/entities/admin-auth/api/adminAuth.graphql` exists and contains more than one operation, split its operations into these one-operation files:

```text
apps/web-admin/src/entities/admin-auth/api/loginAdmin.graphql
apps/web-admin/src/entities/admin-auth/api/logoutAdmin.graphql
apps/web-admin/src/entities/admin-auth/api/currentAdmin.graphql
apps/web-admin/src/entities/admin-auth/api/createAdmin.graphql
```

Preserve the existing GRACE operation comments in each new file, update each `FILE:` marker and `MODULE_MAP` entry for the single operation it owns, then remove the aggregate `adminAuth.graphql` file.

Run:

```bash
rg -n "^(query|mutation) " apps/web-admin/src/entities/admin-auth/api
```

Expected: `LoginAdmin`, `LogoutAdmin`, `CurrentAdmin`, and `CreateAdmin` each appear exactly once, in separate files. This avoids duplicate operation names during codegen and lets runtime imports follow the existing `graphqlClient.request(singleOperationDocument, variables)` pattern.

- [ ] **Step 3: Confirm generated admin auth types exist**

Run:

```bash
rg -n "LoginAdminMutation|LogoutAdminMutation|CurrentAdminQuery|CreateAdminMutation" apps/web-admin/src/shared/api/generated/types.ts
```

Expected:

```text
apps/web-admin/src/shared/api/generated/types.ts:<line>:export type LoginAdminMutationVariables = ...
apps/web-admin/src/shared/api/generated/types.ts:<line>:export type LogoutAdminMutationVariables = ...
apps/web-admin/src/shared/api/generated/types.ts:<line>:export type CurrentAdminQueryVariables = ...
apps/web-admin/src/shared/api/generated/types.ts:<line>:export type CreateAdminMutationVariables = ...
```

- [ ] **Step 4: Regenerate web-admin types if the auth types are absent or operation files changed**

Run if Step 2 changed operation files or Step 3 fails:

```bash
bunx nx run web-admin:codegen
```

Expected: `Successfully ran target codegen for project web-admin`.

- [ ] **Step 5: Confirm cookie transport is already credentialed**

Run:

```bash
rg -n "credentials: 'include'|credentials: \"include\"" apps/web-admin/src/shared/api/graphql-client.ts
```

Expected: one match in `createGraphQLClient`.

- [ ] **Step 6: Commit operation-document setup if it changed files**

Run:

```bash
git status --short -- apps/web-admin/src/entities/admin-auth/api apps/web-admin/src/shared/api/generated/types.ts
```

If Step 2 or Step 4 changed files, commit only those files:

```bash
git add -A apps/web-admin/src/entities/admin-auth/api apps/web-admin/src/shared/api/generated/types.ts
git commit -m "chore(web-admin): split admin auth operation documents"
```

Expected: either no commit because the operation documents and generated types were already current, or one setup commit containing only the split operation documents, the removed aggregate document, and generated `types.ts` if codegen changed it.

## Task 1: Auth Model And GraphQL Client Helpers

**Files:**

- Create: `apps/web-admin/src/entities/admin-auth/model.ts`
- Create: `apps/web-admin/src/entities/admin-auth/model.test.ts`
- Create: `apps/web-admin/src/entities/admin-auth/client.ts`
- Create: `apps/web-admin/src/entities/admin-auth/client.test.ts`

- [ ] **Step 1: Write failing model tests**

Create `apps/web-admin/src/entities/admin-auth/model.test.ts`:

```ts
// FILE: apps/web-admin/src/entities/admin-auth/model.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify web-admin auth model helpers.
//   SCOPE: Covers admin initials, shell-user mapping, and safe same-app return-to parsing; excludes GraphQL transport and route rendering.
//   DEPENDS: apps/web-admin/src/entities/admin-auth/model.ts, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   auth model tests - Prove initials, shell user mapping, and return-to safety rules.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added auth model coverage for login route guard behavior.
// END_CHANGE_SUMMARY

import { describe, expect, it } from 'vitest';
import {
  adminToShellUser,
  getAdminInitials,
  resolveSafeReturnTo,
  type AdminPrincipal,
} from './model';

const admin: AdminPrincipal = {
  id: 'admin-1',
  email: 'owner@example.test',
  name: 'Owner Admin',
  role: 'ADMIN',
  createdAt: '2026-06-07T00:00:00Z',
  updatedAt: '2026-06-07T00:00:00Z',
};

describe('admin auth model helpers', () => {
  it('derives readable initials from admin names and emails', () => {
    expect(getAdminInitials('Owner Admin', 'owner@example.test')).toBe('OA');
    expect(getAdminInitials('Owner', 'owner@example.test')).toBe('O');
    expect(getAdminInitials('', 'support@example.test')).toBe('S');
  });

  it('maps the current admin to the sidebar user contract', () => {
    expect(adminToShellUser(admin)).toEqual({
      name: 'Owner Admin',
      email: 'owner@example.test',
      initials: 'OA',
    });
  });

  it('accepts same-app return paths with search and hash', () => {
    expect(resolveSafeReturnTo('/users?status=active#row-2')).toBe('/users?status=active#row-2');
    expect(resolveSafeReturnTo('/ui-kit')).toBe('/ui-kit');
  });

  it('rejects unsafe or login-loop return paths', () => {
    expect(resolveSafeReturnTo(undefined)).toBe('/');
    expect(resolveSafeReturnTo('https://evil.example/users')).toBe('/');
    expect(resolveSafeReturnTo('//evil.example/users')).toBe('/');
    expect(resolveSafeReturnTo('/\\evil.example/users')).toBe('/');
    expect(resolveSafeReturnTo('users')).toBe('/');
    expect(resolveSafeReturnTo('/login')).toBe('/');
    expect(resolveSafeReturnTo('/login/')).toBe('/');
    expect(resolveSafeReturnTo('/login?from=/users')).toBe('/');
  });
});
```

- [ ] **Step 2: Run model tests and confirm they fail**

Run:

```bash
cd apps/web-admin && bun run test -- src/entities/admin-auth/model.test.ts
```

Expected: fail with an import error because `./model` does not exist.

- [ ] **Step 3: Implement auth model helpers**

Create `apps/web-admin/src/entities/admin-auth/model.ts`:

```ts
// FILE: apps/web-admin/src/entities/admin-auth/model.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define frontend admin auth model helpers for the web-admin app.
//   SCOPE: Owns current-admin principal typing, initials derivation, sidebar user mapping, and safe same-app return-to parsing; excludes GraphQL transport and React context.
//   DEPENDS: apps/web-admin/src/shared/api/generated/types.ts, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminPrincipal - Frontend current-admin principal shape.
//   LoginCredentials - Login form input shape.
//   AuthMutationError - Normalized login mutation error.
//   getAdminInitials - Derive sidebar initials from current admin data.
//   adminToShellUser - Map current admin to the shared shell user contract.
//   resolveSafeReturnTo - Accept only safe same-app return paths.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added auth model helpers for login route guard behavior.
// END_CHANGE_SUMMARY

import type { CurrentAdminQuery } from '@shared/api/generated/types';
import type { AdminUser as ShellAdminUser } from '@shared/ui';

export type AdminPrincipal = NonNullable<CurrentAdminQuery['me']>;

export type LoginCredentials = {
  email: string;
  password: string;
};

export type AuthMutationError = {
  message: string;
  field?: string;
};

// START_CONTRACT: getAdminInitials
//   PURPOSE: Derive short stable initials for sidebar avatar display.
//   INPUTS: { name: string - admin display name, email: string - admin email fallback }
//   OUTPUTS: { string - one or two uppercase initials }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: getAdminInitials
export function getAdminInitials(name: string, email: string): string {
  const nameParts = name.trim().split(/\s+/).filter(Boolean);

  const source = nameParts.length > 0 ? nameParts : [email.split('@')[0] || 'A'];
  return source
    .slice(0, 2)
    .map((part) => part.charAt(0).toUpperCase())
    .join('');
}

// START_CONTRACT: adminToShellUser
//   PURPOSE: Map backend current-admin data into the shared sidebar user contract.
//   INPUTS: { admin: AdminPrincipal - current backend admin }
//   OUTPUTS: { ShellAdminUser - sidebar display user }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: adminToShellUser
export function adminToShellUser(admin: AdminPrincipal): ShellAdminUser {
  return {
    name: admin.name,
    email: admin.email,
    initials: getAdminInitials(admin.name, admin.email),
  };
}

// START_CONTRACT: resolveSafeReturnTo
//   PURPOSE: Accept safe same-app return paths and reject external or login-loop redirects.
//   INPUTS: { rawValue: string | null | undefined - untrusted from query parameter }
//   OUTPUTS: { string - safe app-relative path }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: resolveSafeReturnTo
export function resolveSafeReturnTo(rawValue: string | null | undefined): string {
  const fallback = '/';
  const value = rawValue?.trim();

  if (!value || !value.startsWith('/') || value.startsWith('//') || value.includes('\\')) {
    return fallback;
  }

  try {
    const parsed = new URL(value, window.location.origin);
    const candidate = `${parsed.pathname}${parsed.search}${parsed.hash}`;

    if (
      parsed.origin !== window.location.origin ||
      parsed.pathname === '/login' ||
      parsed.pathname.startsWith('/login/')
    ) {
      return fallback;
    }

    return candidate || fallback;
  } catch {
    return fallback;
  }
}
```

- [ ] **Step 4: Run model tests and confirm they pass**

Run:

```bash
cd apps/web-admin && bun run test -- src/entities/admin-auth/model.test.ts
```

Expected: pass.

- [ ] **Step 5: Write failing client tests**

Create `apps/web-admin/src/entities/admin-auth/client.test.ts`:

```ts
// FILE: apps/web-admin/src/entities/admin-auth/client.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify web-admin admin-auth GraphQL client helpers.
//   SCOPE: Covers current-admin, login, logout, and union error normalization through the shared credentialed GraphQL client; excludes React context and route rendering.
//   DEPENDS: apps/web-admin/src/entities/admin-auth/client.ts, apps/web-admin/src/shared/api/graphql-client.ts, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   auth client tests - Prove generated auth operations are requested and normalized.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added auth GraphQL client helper coverage.
// END_CHANGE_SUMMARY

import { beforeEach, describe, expect, it, vi } from 'vitest';
import { graphqlClient } from '@shared/api/graphql-client';
import { fetchCurrentAdmin, loginAdmin, logoutAdmin } from './client';

vi.mock('@shared/api/graphql-client', () => ({
  graphqlClient: {
    request: vi.fn(),
  },
}));

const requestMock = vi.mocked(graphqlClient.request);

describe('admin auth GraphQL client helpers', () => {
  beforeEach(() => {
    requestMock.mockReset();
  });

  it('fetches the current admin through the CurrentAdmin operation', async () => {
    requestMock.mockResolvedValue({
      me: {
        id: 'admin-1',
        email: 'owner@example.test',
        name: 'Owner Admin',
        role: 'ADMIN',
        createdAt: '2026-06-07T00:00:00Z',
        updatedAt: '2026-06-07T00:00:00Z',
      },
    });

    await expect(fetchCurrentAdmin()).resolves.toMatchObject({
      id: 'admin-1',
      email: 'owner@example.test',
    });
    expect(requestMock).toHaveBeenCalledWith(expect.stringContaining('query CurrentAdmin'));
  });

  it('returns the admin on successful login', async () => {
    requestMock.mockResolvedValue({
      loginAdmin: {
        __typename: 'LoginAdminSuccess',
        admin: {
          id: 'admin-1',
          email: 'owner@example.test',
          name: 'Owner Admin',
          role: 'ADMIN',
          createdAt: '2026-06-07T00:00:00Z',
          updatedAt: '2026-06-07T00:00:00Z',
        },
      },
    });

    await expect(
      loginAdmin({ email: 'owner@example.test', password: 'StrongPassword123!' }),
    ).resolves.toMatchObject({ ok: true, admin: { email: 'owner@example.test' } });
    expect(requestMock).toHaveBeenCalledWith(expect.stringContaining('mutation LoginAdmin'), {
      input: { email: 'owner@example.test', password: 'StrongPassword123!' },
    });
  });

  it('normalizes auth and validation login failures', async () => {
    requestMock.mockResolvedValueOnce({
      loginAdmin: { __typename: 'AuthError', message: 'invalid credentials' },
    });
    await expect(
      loginAdmin({ email: 'owner@example.test', password: 'bad-password' }),
    ).resolves.toEqual({ ok: false, error: { message: 'invalid credentials' } });

    requestMock.mockResolvedValueOnce({
      loginAdmin: {
        __typename: 'ValidationError',
        field: 'email',
        message: 'invalid email',
      },
    });
    await expect(loginAdmin({ email: 'broken', password: 'StrongPassword123!' })).resolves.toEqual({
      ok: false,
      error: { field: 'email', message: 'invalid email' },
    });
  });

  it('logs out through the LogoutAdmin operation', async () => {
    requestMock.mockResolvedValue({ logoutAdmin: { __typename: 'LogoutAdminSuccess', ok: true } });

    await expect(logoutAdmin()).resolves.toBe(true);
    expect(requestMock).toHaveBeenCalledWith(expect.stringContaining('mutation LogoutAdmin'));
  });
});
```

- [ ] **Step 6: Run client tests and confirm they fail**

Run:

```bash
cd apps/web-admin && bun run test -- src/entities/admin-auth/client.test.ts
```

Expected: fail with an import error because `./client` does not exist.

- [ ] **Step 7: Implement auth GraphQL client helpers**

Create `apps/web-admin/src/entities/admin-auth/client.ts`:

```ts
// FILE: apps/web-admin/src/entities/admin-auth/client.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide typed web-admin GraphQL client helpers for admin authentication.
//   SCOPE: Owns current-admin, login, and logout requests plus auth union normalization; excludes React Query context and route navigation.
//   DEPENDS: apps/web-admin/src/entities/admin-auth/api/*.graphql, apps/web-admin/src/shared/api/graphql-client.ts, apps/web-admin/src/shared/api/generated/types.ts, apps/web-admin/src/entities/admin-auth/model.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   fetchCurrentAdmin - Request the current backend admin principal.
//   loginAdmin - Request backend login and normalize success/error result unions.
//   logoutAdmin - Request backend logout and return whether logout succeeded.
//   LoginAdminResult - Frontend normalized login result.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin-auth GraphQL client helpers for login route guard behavior.
// END_CHANGE_SUMMARY

import currentAdminQueryDocument from './api/currentAdmin.graphql?raw';
import loginAdminMutationDocument from './api/loginAdmin.graphql?raw';
import logoutAdminMutationDocument from './api/logoutAdmin.graphql?raw';
import { graphqlClient } from '@shared/api/graphql-client';
import type {
  CurrentAdminQuery,
  LoginAdminMutation,
  LoginAdminMutationVariables,
  LogoutAdminMutation,
} from '@shared/api/generated/types';
import type { AdminPrincipal, AuthMutationError, LoginCredentials } from './model';

export type LoginAdminResult =
  | { ok: true; admin: AdminPrincipal }
  | { ok: false; error: AuthMutationError };

// START_CONTRACT: fetchCurrentAdmin
//   PURPOSE: Load the current admin principal from the backend session cookie.
//   INPUTS: none.
//   OUTPUTS: { Promise<AdminPrincipal | null> - current admin or null when unauthenticated }
//   SIDE_EFFECTS: Sends a credentialed GraphQL request through graphqlClient.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
// END_CONTRACT: fetchCurrentAdmin
export async function fetchCurrentAdmin(): Promise<AdminPrincipal | null> {
  const response = await graphqlClient.request<CurrentAdminQuery>(currentAdminQueryDocument);
  return response.me ?? null;
}

// START_CONTRACT: loginAdmin
//   PURPOSE: Authenticate an admin and normalize backend auth result unions for UI consumption.
//   INPUTS: { credentials: LoginCredentials - email and password from the login form }
//   OUTPUTS: { Promise<LoginAdminResult> - normalized success or user-visible error }
//   SIDE_EFFECTS: Sends a credentialed GraphQL request; backend may set the httpOnly session cookie.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
// END_CONTRACT: loginAdmin
export async function loginAdmin(credentials: LoginCredentials): Promise<LoginAdminResult> {
  const response = await graphqlClient.request<LoginAdminMutation, LoginAdminMutationVariables>(
    loginAdminMutationDocument,
    { input: credentials },
  );

  const result = response.loginAdmin;

  if (result.__typename === 'LoginAdminSuccess') {
    return { ok: true, admin: result.admin };
  }

  if (result.__typename === 'ValidationError') {
    return { ok: false, error: { field: result.field, message: result.message } };
  }

  return { ok: false, error: { message: result.message } };
}

// START_CONTRACT: logoutAdmin
//   PURPOSE: Revoke the current backend admin session.
//   INPUTS: none.
//   OUTPUTS: { Promise<boolean> - backend logout ok flag }
//   SIDE_EFFECTS: Sends a credentialed GraphQL request; backend may clear the httpOnly session cookie.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
// END_CONTRACT: logoutAdmin
export async function logoutAdmin(): Promise<boolean> {
  const response = await graphqlClient.request<LogoutAdminMutation>(logoutAdminMutationDocument);
  return response.logoutAdmin.ok;
}
```

- [ ] **Step 8: Run auth helper tests and confirm they pass**

Run:

```bash
cd apps/web-admin && bun run test -- src/entities/admin-auth/model.test.ts src/entities/admin-auth/client.test.ts
```

Expected: pass.

- [ ] **Step 9: Commit Task 1**

Run:

```bash
git add apps/web-admin/src/entities/admin-auth/model.ts apps/web-admin/src/entities/admin-auth/model.test.ts apps/web-admin/src/entities/admin-auth/client.ts apps/web-admin/src/entities/admin-auth/client.test.ts
git commit -m "feat(web-admin): add admin auth client helpers"
```

Expected: one commit containing only Task 1 model/client files. Operation documents and generated type drift must already be committed in Task 0 or already current before this commit.

## Task 2: Auth Provider And Protected Route Layout

**Files:**

- Create: `apps/web-admin/src/entities/admin-auth/provider.tsx`
- Create: `apps/web-admin/src/entities/admin-auth/provider.test.tsx`
- Create: `apps/web-admin/src/app/protected-admin-layout.tsx`
- Modify: `apps/web-admin/src/app/admin-layout.tsx`
- Modify: `apps/web-admin/src/App.tsx`
- Modify: `apps/web-admin/src/App.test.tsx`

- [ ] **Step 1: Write failing auth provider tests**

Create `apps/web-admin/src/entities/admin-auth/provider.test.tsx`:

```tsx
// FILE: apps/web-admin/src/entities/admin-auth/provider.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the web-admin admin-auth provider boundary.
//   SCOPE: Covers current-admin query state, login refetch behavior, login error normalization, and logout cache clearing; excludes route navigation and page rendering.
//   DEPENDS: apps/web-admin/src/entities/admin-auth/provider.tsx, apps/web-admin/src/entities/admin-auth/client.ts, @tanstack/react-query, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   auth provider tests - Prove auth state and actions are available to route components.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added auth provider coverage for login route guard behavior.
// END_CHANGE_SUMMARY

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AuthProvider, useAdminAuth } from './provider';
import type { AdminPrincipal } from './model';

const fetchCurrentAdminMock = vi.hoisted(() => vi.fn());
const loginAdminMock = vi.hoisted(() => vi.fn());
const logoutAdminMock = vi.hoisted(() => vi.fn());

vi.mock('./client', () => ({
  fetchCurrentAdmin: fetchCurrentAdminMock,
  loginAdmin: loginAdminMock,
  logoutAdmin: logoutAdminMock,
}));

const admin: AdminPrincipal = {
  id: 'admin-1',
  email: 'owner@example.test',
  name: 'Owner Admin',
  role: 'ADMIN',
  createdAt: '2026-06-07T00:00:00Z',
  updatedAt: '2026-06-07T00:00:00Z',
};

let latestAuth: ReturnType<typeof useAdminAuth> | null = null;

function Probe() {
  latestAuth = useAdminAuth();
  return (
    <div>
      <span>{latestAuth.isLoading ? 'loading' : 'ready'}</span>
      <span>{latestAuth.admin?.email ?? 'anonymous'}</span>
    </div>
  );
}

function renderProvider({ staleTime = 0 }: { staleTime?: number } = {}) {
  latestAuth = null;
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, staleTime }, mutations: { retry: false } },
  });

  const result = render(
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <Probe />
      </AuthProvider>
    </QueryClientProvider>,
  );

  return { queryClient, ...result };
}

describe('AuthProvider', () => {
  beforeEach(() => {
    fetchCurrentAdminMock.mockReset();
    loginAdminMock.mockReset();
    logoutAdminMock.mockReset();
  });

  it('loads the current admin from the backend session', async () => {
    fetchCurrentAdminMock.mockResolvedValue(admin);

    renderProvider();

    expect(await screen.findByText('owner@example.test')).toBeInTheDocument();
    expect(screen.getByText('ready')).toBeInTheDocument();
  });

  it('logs in and refreshes the current admin even when cached anonymous auth is fresh', async () => {
    fetchCurrentAdminMock.mockResolvedValueOnce(null).mockResolvedValueOnce(admin);
    loginAdminMock.mockResolvedValue({ ok: true, admin });

    renderProvider({ staleTime: 60_000 });
    await screen.findByText('anonymous');

    await act(async () => {
      const result = await latestAuth?.login({
        email: 'owner@example.test',
        password: 'StrongPassword123!',
      });
      expect(result).toEqual({ ok: true, admin });
    });

    await waitFor(() => expect(fetchCurrentAdminMock).toHaveBeenCalledTimes(2));
    expect(await screen.findByText('owner@example.test')).toBeInTheDocument();
  });

  it('returns normalized login errors without refetching current admin', async () => {
    fetchCurrentAdminMock.mockResolvedValue(null);
    loginAdminMock.mockResolvedValue({ ok: false, error: { message: 'invalid credentials' } });

    renderProvider();
    await screen.findByText('anonymous');

    await act(async () => {
      const result = await latestAuth?.login({
        email: 'owner@example.test',
        password: 'bad-password',
      });
      expect(result).toEqual({ ok: false, error: { message: 'invalid credentials' } });
    });

    expect(fetchCurrentAdminMock).toHaveBeenCalledTimes(1);
  });

  it('logs out and clears current admin plus protected route caches', async () => {
    fetchCurrentAdminMock.mockResolvedValue(admin);
    logoutAdminMock.mockResolvedValue(true);

    const { queryClient } = renderProvider();
    queryClient.setQueryData(['admin-users'], { users: { edges: [{ node: { id: 'user-1' } }] } });
    queryClient.setQueryData(['admin-user', 'user-1'], { user: { id: 'user-1' } });
    await screen.findByText('owner@example.test');

    await act(async () => {
      await latestAuth?.logout();
    });

    expect(logoutAdminMock).toHaveBeenCalledTimes(1);
    expect(queryClient.getQueryData(['admin-users'])).toBeUndefined();
    expect(queryClient.getQueryData(['admin-user', 'user-1'])).toBeUndefined();
    expect(await screen.findByText('anonymous')).toBeInTheDocument();
  });
});
```

- [ ] **Step 2: Run provider tests and confirm they fail**

Run:

```bash
cd apps/web-admin && bun run test -- src/entities/admin-auth/provider.test.tsx
```

Expected: fail with an import error because `./provider` does not exist.

- [ ] **Step 3: Implement AuthProvider**

Create `apps/web-admin/src/entities/admin-auth/provider.tsx`:

```tsx
// FILE: apps/web-admin/src/entities/admin-auth/provider.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide frontend admin-auth state and actions for web-admin routes.
//   SCOPE: Owns React Query current-admin state, login refetch, logout cache clearing, and context access; excludes route navigation and page-specific UI.
//   DEPENDS: react, @tanstack/react-query, apps/web-admin/src/entities/admin-auth/client.ts, apps/web-admin/src/entities/admin-auth/model.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   adminAuthQueryKey - Stable React Query key for current admin.
//   isProtectedAdminQueryKey - Identifies protected admin route data that must be dropped on logout.
//   AuthProvider - Navigation-free auth state provider.
//   useAdminAuth - Context hook for route and shell components.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added auth provider for login route guard behavior.
// END_CHANGE_SUMMARY

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { createContext, useContext, type ReactNode } from 'react';
import {
  fetchCurrentAdmin,
  loginAdmin as requestLoginAdmin,
  logoutAdmin as requestLogoutAdmin,
  type LoginAdminResult,
} from './client';
import type { AdminPrincipal, LoginCredentials } from './model';

export const adminAuthQueryKey = ['admin-auth', 'current-admin'] as const;

export function isProtectedAdminQueryKey(queryKey: readonly unknown[]) {
  return queryKey[0] === 'admin-users' || queryKey[0] === 'admin-user';
}

type AuthContextValue = {
  admin: AdminPrincipal | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  isLogoutPending: boolean;
  login: (credentials: LoginCredentials) => Promise<LoginAdminResult>;
  logout: () => Promise<void>;
  refetchCurrentAdmin: () => Promise<AdminPrincipal | null>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

// START_CONTRACT: AuthProvider
//   PURPOSE: Expose current admin auth state and login/logout actions without owning route navigation.
//   INPUTS: { children: ReactNode - route tree content }
//   OUTPUTS: { JSX.Element - auth context provider }
//   SIDE_EFFECTS: Sends CurrentAdmin/LoginAdmin/LogoutAdmin GraphQL requests through React Query actions.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AuthProvider
export function AuthProvider({ children }: { children: ReactNode }) {
  const queryClient = useQueryClient();
  const currentAdminQuery = useQuery({
    queryKey: adminAuthQueryKey,
    queryFn: fetchCurrentAdmin,
  });
  const logoutMutation = useMutation({ mutationFn: requestLogoutAdmin });

  async function refetchCurrentAdmin() {
    const result = await currentAdminQuery.refetch();
    return result.data ?? null;
  }

  async function login(credentials: LoginCredentials) {
    const result = await requestLoginAdmin(credentials);
    if (result.ok) {
      await refetchCurrentAdmin();
    }
    return result;
  }

  async function logout() {
    await logoutMutation.mutateAsync();
    await queryClient.cancelQueries();
    queryClient.setQueryData(adminAuthQueryKey, null);
    queryClient.removeQueries({
      predicate: (query) => isProtectedAdminQueryKey(query.queryKey),
    });
  }

  return (
    <AuthContext.Provider
      value={{
        admin: currentAdminQuery.data ?? null,
        isAuthenticated: Boolean(currentAdminQuery.data),
        isLoading: currentAdminQuery.isPending,
        isLogoutPending: logoutMutation.isPending,
        login,
        logout,
        refetchCurrentAdmin,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

// START_CONTRACT: useAdminAuth
//   PURPOSE: Read web-admin auth context from route, page, and shell components.
//   INPUTS: none.
//   OUTPUTS: { AuthContextValue - current admin state and actions }
//   SIDE_EFFECTS: Throws if used outside AuthProvider.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: useAdminAuth
export function useAdminAuth(): AuthContextValue {
  const value = useContext(AuthContext);
  if (!value) {
    throw new Error('useAdminAuth must be used within AuthProvider');
  }
  return value;
}
```

- [ ] **Step 4: Run provider tests and confirm they pass**

Run:

```bash
cd apps/web-admin && bun run test -- src/entities/admin-auth/provider.test.tsx
```

Expected: pass.

- [ ] **Step 5: Write failing route guard tests**

Modify `apps/web-admin/src/App.test.tsx` so the request mock can respond by operation. Replace the top-level mock setup with:

```tsx
const requestMock = vi.hoisted(() => vi.fn());
const currentAdmin = {
  id: 'admin-1',
  email: 'owner@example.test',
  name: 'Owner Admin',
  role: 'ADMIN',
  createdAt: '2026-06-07T00:00:00Z',
  updatedAt: '2026-06-07T00:00:00Z',
};

type MockGraphQLRequest = string | { document?: string; operationName?: string };

function getGraphQLDocument(request: MockGraphQLRequest) {
  return typeof request === 'string' ? request : (request.document ?? '');
}

function getGraphQLOperationName(request: MockGraphQLRequest) {
  return typeof request === 'string' ? null : (request.operationName ?? null);
}

function mockGraphQLResponse(request: MockGraphQLRequest) {
  const document = getGraphQLDocument(request);
  const operationName = getGraphQLOperationName(request);

  if (operationName === 'CurrentAdmin' || document.includes('query CurrentAdmin')) {
    return Promise.resolve({ me: currentAdmin });
  }
  if (document.includes('query GetUsers')) {
    return Promise.resolve({
      users: { edges: [], pageInfo: { hasNextPage: false, endCursor: null }, totalCount: 0 },
    });
  }
  return Promise.reject(new Error(`Unhandled GraphQL document: ${document.slice(0, 80)}`));
}
```

Inside `beforeAll`, after `window.matchMedia = ...`, add:

```tsx
requestMock.mockImplementation(mockGraphQLResponse);
```

Inside `afterEach`, before `cleanup()`, add:

```tsx
requestMock.mockImplementation(mockGraphQLResponse);
```

Update the existing route tests so protected pages wait for authenticated rendering and do not replace the operation-aware mock with users-only responses. For example:

```tsx
it('renders the home route with users and UI-kit links', async () => {
  renderApp('/');

  expect(
    await screen.findByRole('heading', { name: 'Monorepo Template Admin' }),
  ).toBeInTheDocument();
  expect(screen.getByRole('navigation', { name: 'Admin navigation' })).toBeInTheDocument();
  expect(screen.getByRole('link', { name: 'Users' })).toHaveAttribute('href', '/users');
  expect(screen.getByRole('link', { name: 'UI Kit' })).toHaveAttribute('href', '/ui-kit');
  expect(screen.getByRole('button', { name: 'Toggle sidebar' })).toBeInTheDocument();
  expect(screen.getAllByText('Overview').length).toBeGreaterThan(0);
  expect(screen.getByRole('link', { name: 'Open users' })).toHaveAttribute('href', '/users');
});

it('renders the users route through the browser router', async () => {
  renderApp('/users');

  expect(await screen.findByText('No users yet. Create one above.')).toBeInTheDocument();
});

it('renders the UI-kit route through the browser router', async () => {
  renderApp('/ui-kit');

  expect(await screen.findByRole('heading', { name: 'UI Kit' })).toBeInTheDocument();
  expect(screen.getByRole('heading', { name: 'Actions' })).toBeInTheDocument();
});
```

For the theme tests, wait for the authenticated shell before reading the toggle:

```tsx
const toggle = await screen.findByRole('button', { name: 'Switch to dark theme' });
```

Add these route guard tests:

```tsx
it('redirects protected routes to login when no current admin exists', async () => {
  requestMock.mockImplementation((request: MockGraphQLRequest) => {
    const document = getGraphQLDocument(request);
    const operationName = getGraphQLOperationName(request);
    if (operationName === 'CurrentAdmin' || document.includes('query CurrentAdmin')) {
      return Promise.resolve({ me: null });
    }
    return mockGraphQLResponse(request);
  });

  renderApp('/users?status=active#directory');

  expect(await screen.findByRole('heading', { name: 'Admin sign in' })).toBeInTheDocument();
  expect(window.location.pathname).toBe('/login');
  expect(window.location.search).toBe('?from=%2Fusers%3Fstatus%3Dactive%23directory');
  expect(screen.queryByText('No users yet. Create one above.')).not.toBeInTheDocument();
});

it('does not render protected content when the current-admin check fails', async () => {
  requestMock.mockImplementation((request: MockGraphQLRequest) => {
    const document = getGraphQLDocument(request);
    const operationName = getGraphQLOperationName(request);
    if (operationName === 'CurrentAdmin' || document.includes('query CurrentAdmin')) {
      return Promise.reject(new Error('API unavailable'));
    }
    return mockGraphQLResponse(request);
  });

  renderApp('/users');

  expect(await screen.findByRole('heading', { name: 'Admin sign in' })).toBeInTheDocument();
  expect(screen.queryByText('No users yet. Create one above.')).not.toBeInTheDocument();
});
```

- [ ] **Step 6: Run route guard tests and confirm they fail**

Run:

```bash
cd apps/web-admin && bun run test -- src/App.test.tsx
```

Expected: fail because `/users` still renders through `AdminLayout` without an auth guard and `/login` does not exist.

- [ ] **Step 7: Implement protected layout and route table**

Create `apps/web-admin/src/app/protected-admin-layout.tsx`:

```tsx
// FILE: apps/web-admin/src/app/protected-admin-layout.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Protect web-admin routes behind the backend admin session.
//   SCOPE: Gates non-login routes on CurrentAdmin state, preserves safe return-to paths, and renders AdminLayout only after auth is known; excludes login page form behavior.
//   DEPENDS: react-router, apps/web-admin/src/app/admin-layout.tsx, apps/web-admin/src/entities/admin-auth/provider.tsx, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ProtectedAdminLayout - Auth gate for all non-login web-admin routes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added protected layout for login route guard behavior.
// END_CHANGE_SUMMARY

import { Navigate, useLocation } from 'react-router';
import { AdminPageShell, Skeleton } from '@shared/ui';
import { useAdminAuth } from '@entities/admin-auth/provider';
import { AdminLayout } from './admin-layout';

function buildReturnTo(location: ReturnType<typeof useLocation>) {
  return `${location.pathname}${location.search}${location.hash}`;
}

// START_CONTRACT: ProtectedAdminLayout
//   PURPOSE: Render protected admin routes only after CurrentAdmin confirms an authenticated admin.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - loading state, login redirect, or AdminLayout }
//   SIDE_EFFECTS: Navigates unauthenticated users to /login with a safe return path.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: ProtectedAdminLayout
export function ProtectedAdminLayout() {
  const location = useLocation();
  const { admin, isLoading } = useAdminAuth();

  if (isLoading) {
    return (
      <AdminPageShell>
        <div
          aria-label="Loading admin session"
          aria-live="polite"
          className="space-y-4"
          role="status"
        >
          <Skeleton className="h-9 w-64" />
          <Skeleton className="h-48 w-full" />
        </div>
      </AdminPageShell>
    );
  }

  if (!admin) {
    return <Navigate replace to={`/login?from=${encodeURIComponent(buildReturnTo(location))}`} />;
  }

  return <AdminLayout admin={admin} />;
}
```

Modify `apps/web-admin/src/app/admin-layout.tsx`:

```tsx
import { Outlet, useLocation } from 'react-router';
import { AdminAppShell } from '@shared/ui';
import { adminToShellUser, type AdminPrincipal } from '@entities/admin-auth/model';
import { resolveAdminShellState } from './admin-navigation';

type AdminLayoutProps = {
  admin: AdminPrincipal;
};

export function AdminLayout({ admin }: AdminLayoutProps) {
  const location = useLocation();
  const shellState = resolveAdminShellState(location.pathname);

  return (
    <AdminAppShell pathname={location.pathname} {...shellState} user={adminToShellUser(admin)}>
      <Outlet />
    </AdminAppShell>
  );
}
```

Modify `apps/web-admin/src/App.tsx`:

```tsx
import { BrowserRouter, Navigate, Route, Routes } from 'react-router';
import { AuthProvider } from '@entities/admin-auth/provider';
import { ProtectedAdminLayout } from './app/protected-admin-layout';
import HomePage from './pages/home';
import LoginPage from './pages/login-page';
import UiKitPage from './pages/ui-kit-page';
import UserDetailPage from './pages/user-detail-page';
import UsersPage from './pages/users-page';

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route element={<ProtectedAdminLayout />}>
            <Route path="/" element={<HomePage />} />
            <Route path="/ui-kit" element={<UiKitPage />} />
            <Route path="/users" element={<UsersPage />} />
            <Route path="/users/:id" element={<UserDetailPage />} />
          </Route>
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}
```

Add a temporary stub page so route tests can compile before Task 3:

```tsx
// FILE: apps/web-admin/src/pages/login-page.tsx
// VERSION: 0.1.0
// START_MODULE_CONTRACT
//   PURPOSE: Temporarily render the web-admin login route until the full login form task replaces it.
//   SCOPE: Provides the route heading required by protected-route tests; excludes submit behavior.
//   DEPENDS: react.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Temporary login page route.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 0.1.0 - Added temporary login route for protected layout wiring.
// END_CHANGE_SUMMARY

export default function LoginPage() {
  return <h1>Admin sign in</h1>;
}
```

- [ ] **Step 8: Run route guard tests**

Run:

```bash
cd apps/web-admin && bun run test -- src/App.test.tsx src/entities/admin-auth/provider.test.tsx
```

Expected: pass after route table and provider wiring.

- [ ] **Step 9: Commit Task 2**

Run:

```bash
git add apps/web-admin/src/entities/admin-auth/provider.tsx apps/web-admin/src/entities/admin-auth/provider.test.tsx apps/web-admin/src/app/protected-admin-layout.tsx apps/web-admin/src/app/admin-layout.tsx apps/web-admin/src/App.tsx apps/web-admin/src/App.test.tsx apps/web-admin/src/pages/login-page.tsx
git commit -m "feat(web-admin): protect admin routes"
```

Expected: one commit containing Task 2 files.

## Task 3: Login Page UI And Redirect Behavior

**Files:**

- Modify: `apps/web-admin/src/pages/login-page.tsx`
- Create: `apps/web-admin/src/pages/login-page.test.tsx`

- [ ] **Step 1: Write failing login page tests**

Create `apps/web-admin/src/pages/login-page.test.tsx`:

```tsx
// FILE: apps/web-admin/src/pages/login-page.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the web-admin login page behavior.
//   SCOPE: Covers form rendering, successful safe redirects, unsafe return fallback, already-authenticated redirect, and visible auth/validation/network errors; excludes backend session creation.
//   DEPENDS: apps/web-admin/src/pages/login-page.tsx, apps/web-admin/src/entities/admin-auth/provider.tsx, react-router, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   login page tests - Prove login UX and redirect behavior.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added login page coverage for route guard behavior.
// END_CHANGE_SUMMARY

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { Route, Routes, BrowserRouter } from 'react-router';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AuthProvider } from '@entities/admin-auth/provider';
import LoginPage from './login-page';

const fetchCurrentAdminMock = vi.hoisted(() => vi.fn());
const loginAdminMock = vi.hoisted(() => vi.fn());
const logoutAdminMock = vi.hoisted(() => vi.fn());

vi.mock('@entities/admin-auth/client', () => ({
  fetchCurrentAdmin: fetchCurrentAdminMock,
  loginAdmin: loginAdminMock,
  logoutAdmin: logoutAdminMock,
}));

const admin = {
  id: 'admin-1',
  email: 'owner@example.test',
  name: 'Owner Admin',
  role: 'ADMIN',
  createdAt: '2026-06-07T00:00:00Z',
  updatedAt: '2026-06-07T00:00:00Z',
};

function renderLogin(path = '/login') {
  window.history.pushState({}, '', path);
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, staleTime: 60_000 }, mutations: { retry: false } },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AuthProvider>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/" element={<h1>Overview</h1>} />
            <Route path="/users" element={<h1>Users</h1>} />
          </Routes>
        </AuthProvider>
      </BrowserRouter>
    </QueryClientProvider>,
  );
}

async function submitLogin() {
  fireEvent.change(screen.getByLabelText('Email'), {
    target: { value: 'owner@example.test' },
  });
  fireEvent.change(screen.getByLabelText('Password'), {
    target: { value: 'StrongPassword123!' },
  });
  fireEvent.click(screen.getByRole('button', { name: 'Sign in' }));
}

describe('LoginPage', () => {
  beforeEach(() => {
    fetchCurrentAdminMock.mockReset();
    loginAdminMock.mockReset();
    logoutAdminMock.mockReset();
    fetchCurrentAdminMock.mockResolvedValue(null);
  });

  it('renders the login form', async () => {
    renderLogin();

    expect(await screen.findByRole('heading', { name: 'Admin sign in' })).toBeInTheDocument();
    expect(screen.getByLabelText('Email')).toBeInTheDocument();
    expect(screen.getByLabelText('Password')).toHaveAttribute('type', 'password');
    expect(screen.getByRole('button', { name: 'Sign in' })).toBeInTheDocument();
  });

  it('redirects to a safe return path after login', async () => {
    fetchCurrentAdminMock.mockResolvedValueOnce(null).mockResolvedValueOnce(admin);
    loginAdminMock.mockResolvedValue({ ok: true, admin });

    renderLogin('/login?from=%2Fusers%3Fstatus%3Dactive%23directory');
    await submitLogin();

    await waitFor(() => expect(window.location.pathname).toBe('/users'));
    expect(window.location.search).toBe('?status=active');
    expect(window.location.hash).toBe('#directory');
  });

  it('falls back to overview for unsafe return paths', async () => {
    fetchCurrentAdminMock.mockResolvedValueOnce(null).mockResolvedValueOnce(admin);
    loginAdminMock.mockResolvedValue({ ok: true, admin });

    renderLogin('/login?from=https%3A%2F%2Fevil.example%2Fusers');
    await submitLogin();

    await waitFor(() => expect(window.location.pathname).toBe('/'));
  });

  it('redirects an already-authenticated admin away from login', async () => {
    fetchCurrentAdminMock.mockResolvedValue(admin);

    renderLogin('/login?from=%2Fusers');

    await waitFor(() => expect(window.location.pathname).toBe('/users'));
  });

  it('renders auth and validation errors', async () => {
    loginAdminMock.mockResolvedValueOnce({
      ok: false,
      error: { message: 'invalid credentials' },
    });

    renderLogin();
    await submitLogin();

    expect(await screen.findByText('invalid credentials')).toBeInTheDocument();

    loginAdminMock.mockResolvedValueOnce({
      ok: false,
      error: { field: 'email', message: 'invalid email' },
    });
    await submitLogin();

    expect(await screen.findByText('email: invalid email')).toBeInTheDocument();
  });

  it('renders a stable fallback for network failures', async () => {
    loginAdminMock.mockRejectedValue(new Error('network failed'));

    renderLogin();
    await submitLogin();

    expect(
      await screen.findByText('Unable to sign in. Try again after the API is available.'),
    ).toBeInTheDocument();
  });
});
```

- [ ] **Step 2: Run login page tests and confirm they fail**

Run:

```bash
cd apps/web-admin && bun run test -- src/pages/login-page.test.tsx
```

Expected: fail because the temporary login page has no form behavior.

- [ ] **Step 3: Implement the login page**

Replace `apps/web-admin/src/pages/login-page.tsx` with:

```tsx
// FILE: apps/web-admin/src/pages/login-page.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Render the public web-admin login route.
//   SCOPE: Owns login form state, error presentation, safe return-to navigation, and already-authenticated redirects; excludes backend session storage and protected shell rendering.
//   DEPENDS: react, react-router, apps/web-admin/src/entities/admin-auth/provider.tsx, apps/web-admin/src/entities/admin-auth/model.ts, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Public login page for backend cookie-backed admin sessions.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added full login page for web-admin route guard behavior.
// END_CHANGE_SUMMARY

import { type FormEvent, useEffect, useState } from 'react';
import { Navigate, useLocation, useNavigate } from 'react-router';
import { useAdminAuth } from '@entities/admin-auth/provider';
import { resolveSafeReturnTo } from '@entities/admin-auth/model';
import {
  Alert,
  AlertDescription,
  AlertTitle,
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Input,
  Label,
} from '@shared/ui';

function errorMessageFromUnknown(error: unknown): string {
  return error instanceof Error
    ? 'Unable to sign in. Try again after the API is available.'
    : 'Unable to sign in. Try again after the API is available.';
}

export default function LoginPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const { admin, isLoading, login } = useAdminAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const returnTo = resolveSafeReturnTo(new URLSearchParams(location.search).get('from'));

  useEffect(() => {
    if (admin) {
      navigate(returnTo, { replace: true });
    }
  }, [admin, navigate, returnTo]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(null);
    setIsSubmitting(true);
    try {
      const result = await login({ email, password });
      if (!result.ok) {
        setError(
          result.error.field
            ? `${result.error.field}: ${result.error.message}`
            : result.error.message,
        );
        return;
      }
      navigate(returnTo, { replace: true });
    } catch (loginError) {
      setError(errorMessageFromUnknown(loginError));
    } finally {
      setIsSubmitting(false);
    }
  }

  if (!isLoading && admin) {
    return <Navigate replace to={returnTo} />;
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-background px-4 py-10">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Admin sign in</CardTitle>
          <CardDescription>Use your admin account to continue.</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-4" onSubmit={handleSubmit}>
            <div className="space-y-2">
              <Label htmlFor="admin-email">Email</Label>
              <Input
                autoComplete="email"
                id="admin-email"
                name="email"
                onChange={(event) => setEmail(event.target.value)}
                required
                type="email"
                value={email}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="admin-password">Password</Label>
              <Input
                autoComplete="current-password"
                id="admin-password"
                name="password"
                onChange={(event) => setPassword(event.target.value)}
                required
                type="password"
                value={password}
              />
            </div>
            {error ? (
              <Alert variant="destructive">
                <AlertTitle>Sign in failed</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            ) : null}
            <Button className="w-full" disabled={isSubmitting || isLoading} type="submit">
              {isSubmitting ? 'Signing in...' : 'Sign in'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </main>
  );
}
```

- [ ] **Step 4: Run login tests and App tests**

Run:

```bash
cd apps/web-admin && bun run test -- src/pages/login-page.test.tsx src/App.test.tsx
```

Expected: pass.

- [ ] **Step 5: Commit Task 3**

Run:

```bash
git add apps/web-admin/src/pages/login-page.tsx apps/web-admin/src/pages/login-page.test.tsx apps/web-admin/src/App.test.tsx
git commit -m "feat(web-admin): add login page"
```

Expected: one commit containing Task 3 files.

## Task 4: Real Sidebar User And Logout Action

**Files:**

- Modify: `apps/web-admin/src/shared/ui/layout/admin-shell-types.ts`
- Modify: `apps/web-admin/src/shared/ui/layout/admin-app-shell.tsx`
- Modify: `apps/web-admin/src/shared/ui/layout/app-sidebar.tsx`
- Modify: `apps/web-admin/src/shared/ui/layout/nav-user.tsx`
- Modify: nearest shell tests, preferably `apps/web-admin/src/shared/ui/layout/admin-shell.test.tsx`

- [ ] **Step 1: Write failing shell logout test**

In the nearest existing shell test file, add:

```tsx
it('renders the authenticated admin and calls logout from the user menu', async () => {
  const onLogout = vi.fn();

  render(
    <MemoryRouter>
      <AdminAppShell
        breadcrumbs={[{ label: 'Overview' }]}
        navigation={[]}
        pathname="/"
        referenceItems={[]}
        teams={[]}
        user={{ name: 'Owner Admin', email: 'owner@example.test', initials: 'OA' }}
        onLogout={onLogout}
      >
        <p>Protected content</p>
      </AdminAppShell>
    </MemoryRouter>,
  );

  const userButton = screen.getByRole('button', { name: /Owner Admin owner@example\.test/i });
  fireEvent.pointerDown(userButton, { button: 0, ctrlKey: false, pointerType: 'mouse' });
  fireEvent.keyDown(userButton, { key: 'Enter' });
  fireEvent.click(await screen.findByRole('menuitem', { name: 'Logout' }));

  expect(onLogout).toHaveBeenCalledTimes(1);
  expect(screen.getByText('owner@example.test')).toBeInTheDocument();
});

it('disables logout while logout is pending', async () => {
  render(
    <MemoryRouter>
      <AdminAppShell
        breadcrumbs={[{ label: 'Overview' }]}
        navigation={[]}
        pathname="/"
        referenceItems={[]}
        teams={[]}
        user={{ name: 'Owner Admin', email: 'owner@example.test', initials: 'OA' }}
        isLogoutPending
        onLogout={vi.fn()}
      >
        <p>Protected content</p>
      </AdminAppShell>
    </MemoryRouter>,
  );

  const userButton = screen.getByRole('button', { name: /Owner Admin owner@example\.test/i });
  fireEvent.pointerDown(userButton, { button: 0, ctrlKey: false, pointerType: 'mouse' });
  fireEvent.keyDown(userButton, { key: 'Enter' });

  expect(await screen.findByRole('menuitem', { name: 'Logging out...' })).toHaveAttribute(
    'data-disabled',
  );
});
```

If the file does not already import `fireEvent`, `MemoryRouter`, and `vi`, update its imports:

```tsx
import { fireEvent, render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { describe, expect, it, vi } from 'vitest';
```

- [ ] **Step 2: Run shell test and confirm it fails**

Run:

```bash
cd apps/web-admin && bun run test -- src/shared/ui/layout/admin-shell.test.tsx
```

Expected: fail because `AdminAppShell` does not accept `onLogout` and `NavUser` has no logout menu item.

- [ ] **Step 3: Update shell user contracts and pass-through props**

Modify `apps/web-admin/src/shared/ui/layout/admin-shell-types.ts`:

```ts
export type AdminUserAction = {
  isLogoutPending?: boolean;
  onLogout?: () => void | Promise<void>;
};

export type AdminUser = {
  name: string;
  email: string;
  initials: string;
  avatarUrl?: string;
};
```

Modify `AdminAppShellProps` in `apps/web-admin/src/shared/ui/layout/admin-app-shell.tsx`:

```tsx
type AdminAppShellProps = {
  breadcrumbs: AdminBreadcrumbItem[];
  children: ReactNode;
  isLogoutPending?: boolean;
  navigation: AdminNavigationGroup[];
  onLogout?: () => void | Promise<void>;
  pathname: string;
  referenceItems: AdminProjectItem[];
  teams: AdminTeamItem[];
  user: AdminUser;
};
```

And pass it through:

```tsx
export function AdminAppShell({
  breadcrumbs,
  children,
  isLogoutPending,
  navigation,
  onLogout,
  referenceItems,
  teams,
  user,
}: AdminAppShellProps) {
  return (
    <SidebarProvider>
      <TooltipProvider>
        <AppSidebar
          isLogoutPending={isLogoutPending}
          navigation={navigation}
          onLogout={onLogout}
          referenceItems={referenceItems}
          teams={teams}
          user={user}
        />
        <SidebarInset>
          <AdminShellHeader breadcrumbs={breadcrumbs} />
          {children}
        </SidebarInset>
      </TooltipProvider>
    </SidebarProvider>
  );
}
```

Modify `AppSidebarProps` and call in `apps/web-admin/src/shared/ui/layout/app-sidebar.tsx`:

```tsx
type AppSidebarProps = {
  isLogoutPending?: boolean;
  navigation: AdminNavigationGroup[];
  onLogout?: () => void | Promise<void>;
  referenceItems: AdminProjectItem[];
  teams: AdminTeamItem[];
  user: AdminUser;
};

export function AppSidebar({
  isLogoutPending,
  navigation,
  onLogout,
  referenceItems,
  teams,
  user,
}: AppSidebarProps) {
  return (
    <Sidebar aria-label="Admin navigation" collapsible="icon" role="navigation">
      <SidebarHeader>
        <TeamSwitcher teams={teams} />
      </SidebarHeader>
      <SidebarContent>
        <NavMain groups={navigation} />
        <NavProjects projects={referenceItems} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser isLogoutPending={isLogoutPending} onLogout={onLogout} user={user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
```

- [ ] **Step 4: Replace static user menu actions with logout**

Modify `apps/web-admin/src/shared/ui/layout/nav-user.tsx` imports:

```tsx
import { ChevronsUpDownIcon, LogOutIcon } from 'lucide-react';
```

Change props:

```tsx
type NavUserProps = {
  isLogoutPending?: boolean;
  onLogout?: () => void | Promise<void>;
  user: AdminUser;
};
```

Replace the disabled menu group with:

```tsx
<DropdownMenuGroup>
  <DropdownMenuItem disabled={!onLogout || isLogoutPending} onClick={() => void onLogout?.()}>
    <LogOutIcon aria-hidden="true" />
    {isLogoutPending ? 'Logging out...' : 'Logout'}
  </DropdownMenuItem>
</DropdownMenuGroup>
```

Update the function signature:

```tsx
export function NavUser({ isLogoutPending, onLogout, user }: NavUserProps) {
```

- [ ] **Step 5: Wire AdminLayout logout navigation**

Modify `apps/web-admin/src/app/admin-layout.tsx` imports:

```tsx
import { Outlet, useLocation, useNavigate } from 'react-router';
import { AdminAppShell } from '@shared/ui';
import { adminToShellUser, type AdminPrincipal } from '@entities/admin-auth/model';
import { useAdminAuth } from '@entities/admin-auth/provider';
import { resolveAdminShellState } from './admin-navigation';
```

Add logout handling inside `AdminLayout` after `const location = useLocation();`:

```tsx
const navigate = useNavigate();
const { isLogoutPending, logout } = useAdminAuth();

async function handleLogout() {
  await logout();
  navigate('/login', { replace: true });
}
```

Pass the handler to `AdminAppShell`:

```tsx
<AdminAppShell
  pathname={location.pathname}
  {...shellState}
  isLogoutPending={isLogoutPending}
  onLogout={handleLogout}
  user={adminToShellUser(admin)}
>
  <Outlet />
</AdminAppShell>
```

- [ ] **Step 6: Run shell and app tests**

Run:

```bash
cd apps/web-admin && bun run test -- src/shared/ui/layout/admin-shell.test.tsx src/App.test.tsx
```

Expected: pass.

- [ ] **Step 7: Commit Task 4**

Run:

```bash
git add apps/web-admin/src/shared/ui/layout/admin-shell-types.ts apps/web-admin/src/shared/ui/layout/admin-app-shell.tsx apps/web-admin/src/shared/ui/layout/app-sidebar.tsx apps/web-admin/src/shared/ui/layout/nav-user.tsx apps/web-admin/src/shared/ui/layout/admin-shell.test.tsx apps/web-admin/src/app/admin-layout.tsx
git commit -m "feat(web-admin): add admin sidebar logout"
```

Expected: one commit containing Task 4 files.

## Task 5: Real Browser Login And Logout E2E

**Files:**

- Modify: `apps/web-admin/e2e/helpers.ts`
- Modify: `apps/web-admin/e2e/users-flow.spec.ts`

- [ ] **Step 1: Add browser login helpers**

Modify `apps/web-admin/e2e/helpers.ts` map and exports. Add to `START_MODULE_MAP`:

```text
//   adminEmail - Test admin email supplied by Playwright environment.
//   adminPassword - Test admin password supplied by Playwright environment.
//   loginThroughUi - Logs in through the real browser login page.
```

Change the constants from private to exported:

```ts
export const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'e2e-admin@example.test';
export const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'StrongPassword123!';
```

Add:

```ts
export async function loginThroughUi(page: import('@playwright/test').Page, from = '/users') {
  await page.goto(`/login?from=${encodeURIComponent(from)}`);
  await page.getByLabel('Email').fill(adminEmail);
  await page.getByLabel('Password').fill(adminPassword);
  await page.getByRole('button', { name: 'Sign in' }).click();
  await expect(page).toHaveURL(new RegExp(`${from.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}$`));
}
```

- [ ] **Step 2: Replace cookie install in users e2e**

Modify imports in `apps/web-admin/e2e/users-flow.spec.ts`:

```ts
import {
  createUser,
  loginThroughUi,
  uniqueEmail,
  withAuthenticatedGraphQLContext,
} from './helpers';
```

Remove the existing `beforeEach` cookie install. Add `await loginThroughUi(page, '/')` as the first browser step in each existing non-auth-flow test (`users page creates...`, `users page shows duplicate...`, `shell navigation...`, and `mobile shell...`). Keep `withAuthenticatedGraphQLContext` for direct API seeding only.

Do not add a global browser-login `beforeEach`; the auth-flow test below must start from an unauthenticated browser context.

Add the auth flow test before users CRUD:

```ts
test('protected routes require login and logout revokes browser access', async ({
  context,
  page,
}) => {
  await page.goto('/users');
  await expect(page).toHaveURL(/\/login\?from=%2Fusers$/);
  await expect(page.getByRole('heading', { name: 'Admin sign in' })).toBeVisible();

  await loginThroughUi(page, '/users');
  await expect(page.getByRole('heading', { name: 'Users' })).toBeVisible();

  await page.getByRole('button', { name: /E2E Admin/ }).click();
  await page.getByRole('menuitem', { name: 'Logout' }).click();
  await expect(page).toHaveURL(/\/login$/);

  const freshPage = await context.newPage();
  await freshPage.goto('/users');
  await expect(freshPage).toHaveURL(/\/login\?from=%2Fusers$/);
  await expect(freshPage.getByRole('heading', { name: 'Admin sign in' })).toBeVisible();
});
```

- [ ] **Step 3: Run web-admin e2e**

Run:

```bash
bunx nx run web-admin:e2e
```

Expected: pass. If a local API/web server is already using the e2e ports, rerun with isolated ports:

```bash
TEST_RESOURCE_PREFIX=mt-login-guard TEST_COMPOSE_PROJECT=mt-login-guard TEST_POSTGRES_CONTAINER_NAME=mt-login-guard-postgres TEST_REDIS_CONTAINER_NAME=mt-login-guard-redis TEST_POSTGRES_VOLUME=mt-login-guard-pg-test-data TEST_POSTGRES_PORT=20101 TEST_REDIS_PORT=20102 E2E_API_PORT=20180 E2E_WEB_PORT=20130 bunx nx run web-admin:e2e
```

Expected: Playwright list reporter shows all web-admin e2e tests passed.

- [ ] **Step 4: Commit Task 5**

Run:

```bash
git add apps/web-admin/e2e/helpers.ts apps/web-admin/e2e/users-flow.spec.ts
git commit -m "test(web-admin): cover login and logout e2e"
```

Expected: one commit containing Task 5 files.

## Task 6: GRACE Docs And Final Verification

**Files:**

- Modify: `docs/requirements.xml`
- Modify: `docs/development-plan.xml`
- Modify: `docs/knowledge-graph.xml`
- Modify: `docs/verification-plan.xml`
- Inspect: `docs/operational-packets.xml`

- [ ] **Step 1: Update requirements acceptance criteria**

In `docs/requirements.xml`, update the web-admin auth use case text for `UC-012` to include the frontend login guard. Use this wording inside the existing `AcceptanceCriteria` element:

```xml
`loginAdmin` sets an httpOnly cookie, `me` identifies the current admin or returns null without a session, `logoutAdmin` revokes the session, `createAdmin` requires an authenticated admin, protected user GraphQL operations reject unauthenticated callers, public REST `/api/users` remains public, web-admin GraphQL transport sends credentialed cookie requests, `/login` renders the admin login form without the sidebar shell, every other web-admin route redirects unauthenticated users to `/login?from=...`, successful login returns to the safe same-app requested route, sidebar logout disables while pending, clears protected frontend query caches, and returns the browser to `/login`.
```

- [ ] **Step 2: Update development plan module exports**

In `docs/development-plan.xml`, within `M-WEB-ADMIN`, add or update interface exports with:

```xml
<export-AdminAuthFrontend PURPOSE="Expose CurrentAdmin-backed auth provider, login page, protected route layout, safe return-to handling, and sidebar logout for the web-admin app." />
<export-LoginRoute PURPOSE="Render `/login` outside the sidebar shell and authenticate through cookie-backed `loginAdmin`." />
<export-ProtectedRoutes PURPOSE="Protect all non-login admin routes through `CurrentAdmin` before rendering the global admin shell." />
```

Also add source targets under `M-WEB-ADMIN`:

```xml
<source>apps/web-admin/src/entities/admin-auth</source>
<source>apps/web-admin/src/app/protected-admin-layout.tsx</source>
<source>apps/web-admin/src/pages/login-page.tsx</source>
```

Replace any stale `apps/web-admin/src/entities/admin-auth/api/adminAuth.graphql` source entry with the parent `apps/web-admin/src/entities/admin-auth` source or the four split operation files.

- [ ] **Step 3: Update knowledge graph**

In `docs/knowledge-graph.xml`, update `M-WEB-ADMIN` with these paths and exports:

```xml
<path>apps/web-admin/src/entities/admin-auth</path>
<path>apps/web-admin/src/app/protected-admin-layout.tsx</path>
<path>apps/web-admin/src/pages/login-page.tsx</path>
<export-AdminAuthFrontend PURPOSE="CurrentAdmin-backed auth provider, auth GraphQL helpers, safe return-to parsing, and shell user mapping." />
<export-ProtectedAdminLayout PURPOSE="Protects every non-login web-admin route before rendering the sidebar shell." />
<export-LoginPage PURPOSE="Public login form for backend cookie-backed admin sessions." />
<export-SidebarLogout PURPOSE="Sidebar user-menu logout action that revokes the backend session, disables while pending, clears protected frontend query caches, and returns to login." />
```

Replace any stale `apps/web-admin/src/entities/admin-auth/api/adminAuth.graphql` path with `apps/web-admin/src/entities/admin-auth` or the four split operation document paths.

Add a cross-link:

```xml
<CrossLink from="M-WEB-ADMIN" to="M-API" relation="authenticates through cookie-backed loginAdmin, reads me through CurrentAdmin, and revokes sessions through logoutAdmin" />
```

- [ ] **Step 4: Update verification plan**

In `docs/verification-plan.xml`, under `V-M-WEB-ADMIN`, add scenarios:

```xml
<scenario-frontend-auth-1 kind="success">Unauthenticated visits to `/users` redirect to `/login?from=%2Fusers` without rendering protected user content.</scenario-frontend-auth-1>
<scenario-frontend-auth-2 kind="success">Login through `/login` with the seeded e2e admin forces a fresh `me` check even after a cached anonymous response and returns to the safe requested route, preserving search and hash when present.</scenario-frontend-auth-2>
<scenario-frontend-auth-3 kind="failure">Unsafe `from` values including absolute URLs, protocol-relative URLs, malformed paths, `/login`, and login loops fall back to `/`.</scenario-frontend-auth-3>
<scenario-frontend-auth-4 kind="success">Sidebar logout calls `logoutAdmin`, disables while pending, clears frontend auth state and protected query caches, returns to `/login`, and protected routes redirect again after logout.</scenario-frontend-auth-4>
```

Add checks:

```xml
<check-frontend-auth-1>cd apps/web-admin &amp;&amp; bun run test -- src/entities/admin-auth/model.test.ts src/entities/admin-auth/client.test.ts src/entities/admin-auth/provider.test.tsx src/pages/login-page.test.tsx src/shared/ui/layout/admin-shell.test.tsx src/App.test.tsx</check-frontend-auth-1>
<check-frontend-auth-2>bunx nx run web-admin:codegen</check-frontend-auth-2>
<check-frontend-auth-3>bunx nx run web-admin:typecheck</check-frontend-auth-3>
<check-frontend-auth-4>bunx nx test web-admin</check-frontend-auth-4>
<check-frontend-auth-5>bunx nx lint web-admin</check-frontend-auth-5>
<check-frontend-auth-6>bunx nx build web-admin</check-frontend-auth-6>
<check-frontend-auth-7>bunx nx run web-admin:e2e</check-frontend-auth-7>
```

Replace any stale verification file entry for `apps/web-admin/src/entities/admin-auth/api/adminAuth.graphql` with the split operation files or the parent `apps/web-admin/src/entities/admin-auth/api` directory. Then run:

```bash
rg -n "adminAuth\\.graphql" docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml
```

Expected: no matches.

- [ ] **Step 5: Run focused unit and generated checks**

Run:

```bash
bunx nx run web-admin:codegen
bunx nx run web-admin:typecheck
bunx nx test web-admin
bunx nx lint web-admin
bunx nx build web-admin
```

Expected: all commands pass.

- [ ] **Step 6: Run e2e**

Run:

```bash
bunx nx run web-admin:e2e
```

Expected: all web-admin Playwright tests pass. If local port conflicts exist, use the isolated command from Task 5 Step 3.

- [ ] **Step 7: Run GRACE validation**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected: XML parses and `grace lint --path .` has no new errors.

- [ ] **Step 8: Commit docs and final verification**

Run:

```bash
git add docs/requirements.xml docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml
git commit -m "docs: update web-admin login guard contracts"
```

Expected: one docs commit.

## Self-Review

### Spec Coverage

- `/login` outside sidebar: Task 2 route table and Task 3 login page.
- Protect every other route: Task 2 `ProtectedAdminLayout` and App tests.
- Return to intended route with search/hash: Task 1 model tests and Task 3 login tests.
- Reject unsafe return paths: Task 1 model tests and Task 3 login tests.
- Backend cookie only: Task 1 client helpers use existing credentialed `graphqlClient`; no browser storage is introduced.
- Post-login auth refresh: Task 2 provider uses `currentAdminQuery.refetch()` and tests the production `staleTime: 60_000` case.
- Real admin in sidebar: Task 2 passes admin to `AdminLayout`; Task 4 maps and renders real user.
- Sidebar logout: Task 4 unit test and Task 5 e2e.
- Logout pending and cache cleanup: Task 2 clears protected query caches on logout; Task 4 disables the logout menu item while pending.
- Guard failure does not expose protected content: Task 2 App tests.
- Focused verification plus closeout build/XML/GRACE lint: Task 6.

### Placeholder Scan

The plan contains no incomplete steps or instructions that require guessing a missing file path.

### Type Consistency

- `AdminPrincipal` is derived from `CurrentAdminQuery['me']`.
- `LoginCredentials` feeds `LoginAdminInput` through generated `LoginAdminMutationVariables`.
- `AuthProvider.login` returns `LoginAdminResult` from `client.ts`.
- `AuthProvider.logout` exposes `isLogoutPending` and removes protected `admin-users` / `admin-user` query caches after backend session revocation.
- `AdminLayout` accepts `AdminPrincipal`, maps it with `adminToShellUser`, and passes the shared `AdminUser` shape into `AdminAppShell`.
