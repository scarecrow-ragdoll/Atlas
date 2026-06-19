<!-- FILE: docs/superpowers/plans/2026-06-05-swap-web-admin-vite-web-next.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Provide the task-by-task implementation plan for swapping web-admin to Vite and public web to Next.js App Router. -->
<!--   SCOPE: Covers frontend framework migration, env/config contracts, tests, e2e, coverage preflight, deployment ownership, GRACE docs, and final verification; excludes implementation performed by this document. -->
<!--   DEPENDS: docs/superpowers/specs/2026-06-05-swap-web-admin-vite-web-next-design.md, apps/web-admin, apps/web, tools/codegen, tools/coverage, docker/web.Dockerfile, .gitlab-ci.yml, deploy/dokploy/docker-compose.template.yml, docs/*.xml. -->
<!--   LINKS: M-WEB-ADMIN / M-WEB / M-GRAPHQL-SCHEMA / M-COVERAGE-GATE / M-CI-CD / V-M-WEB-ADMIN / V-M-WEB / V-M-COVERAGE-GATE / V-M-CI-CD. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Header - Defines the execution goal, architecture, and required execution sub-skill. -->
<!--   Source Spec - Anchors the approved design and subagent review decisions. -->
<!--   File Structure - Defines the files to create, modify, and remove. -->
<!--   Tasks - Provides TDD-oriented bite-sized implementation, verification, and commit steps. -->
<!--   Self-Review - Records spec coverage, placeholder scan, and type consistency checks. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.1 - Addressed subagent review findings for GRACE order, GraphQL documents, runtime Next REST env, lockfile, test config, coverage allowlists, and deploy checks. -->
<!-- END_CHANGE_SUMMARY -->

# Swap Web Admin Vite And Public Web Next Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Convert `apps/web-admin` into a Vite + React + React Router GraphQL admin SPA and convert `apps/web` into a Next.js App Router + REST public web app while preserving project names and Nx command names.

**Architecture:** Keep the transport split intact: `web-admin` owns admin GraphQL documents, `graphql-request`, generated operation types, and client-side routes; `web` owns public REST users flow through Next App Router. Retarget the existing deployable `web` Docker image to the new Next public app, explicitly defer or define Vite admin deployment, and synchronize coverage, e2e, codegen, and GRACE contracts.

**Tech Stack:** Bun workspaces, Nx 20, TypeScript 5.5, React 19, Vite 5, React Router, Next.js 15 App Router, React Query 5, graphql-request 7, GraphQL Codegen 5, Vitest/jsdom, Playwright, Docker Buildx, Dokploy, GRACE XML.

---

## Source Spec

- Design: `docs/superpowers/specs/2026-06-05-swap-web-admin-vite-web-next-design.md`
- Approved reviewed direction:
  - `apps/web-admin` remains Nx project `web-admin` and becomes Vite + React SPA + React Router.
  - `apps/web` remains Nx project `web` and becomes Next.js App Router.
  - Admin keeps GraphQL/codegen and generated types under `apps/web-admin/src/shared/api/generated/types.ts`.
  - Public web remains REST-only and must not import GraphQL documents, generated GraphQL types, or `graphql-request`.
  - Admin GraphQL env moves to Vite browser env, for example `VITE_GRAPHQL_API_URL`.
  - Public REST env uses a Next deploy-safe model: browser code calls same-origin `/api/users`, and Next server/runtime proxy code reads `WEB_API_BASE_URL`. Do not depend on production Dokploy runtime values for `NEXT_PUBLIC_*` browser bundle configuration.
  - Existing `web` image is retargeted to `apps/web`; Vite admin deployment is explicit or blocked from release handoff.
  - `tools/coverage/preflight.mjs`, e2e configs, coverage allowlists, deployment files, and GRACE XML must be updated.

## Scope Check

This spec touches multiple surfaces, but they are not independent projects: the framework swap is only complete when app runtime, tests, e2e, coverage, deploy, and GRACE docs agree on the same ownership. Keep it as one implementation plan with task-level commits because partial completion can otherwise leave Nx, e2e, or deployment contracts lying about which app owns which framework.

## File Structure

### `apps/web-admin` Vite GraphQL SPA

- Create: `apps/web-admin/index.html` - Vite HTML root.
- Create: `apps/web-admin/vite.config.ts` - Vite/Vitest config, aliases, jsdom coverage.
- Create: `apps/web-admin/src/main.tsx` - Vite bootstrap, excluded from coverage.
- Create: `apps/web-admin/src/App.tsx` - BrowserRouter route table.
- Create: `apps/web-admin/src/app/providers.tsx` - React Query provider.
- Create: `apps/web-admin/src/pages/home.tsx` - Admin home route.
- Create: `apps/web-admin/src/pages/users-page.tsx` - GraphQL users list/create route.
- Create: `apps/web-admin/src/pages/user-detail-page.tsx` - GraphQL user detail route.
- Create: `apps/web-admin/src/pages/users-page.test.tsx` - users route tests.
- Create: `apps/web-admin/src/pages/user-detail-page.test.tsx` - detail route tests.
- Create: `apps/web-admin/src/App.test.tsx` - route smoke tests.
- Create: `apps/web-admin/src/styles.css` - shared admin styles.
- Modify: `apps/web-admin/package.json` - Vite scripts/deps and `react-router`.
- Modify: `apps/web-admin/project.json` - same Nx target names, Vite-backed commands.
- Modify: `apps/web-admin/tsconfig.json` - Vite JSX/env settings.
- Modify or remove: `apps/web-admin/vitest.config.ts` - remove stale Next-admin test config or turn it into a thin redirect to `vite.config.ts`.
- Modify: `apps/web-admin/src/shared/config/index.ts` - `import.meta.env.VITE_GRAPHQL_API_URL`.
- Modify: `apps/web-admin/src/shared/config/index.test.ts` - Vite env tests.
- Modify: `apps/web-admin/src/app/config.ts` - keep app-level config re-export.
- Review and keep/update: `apps/web-admin/src/entities/user/api/users.graphql` - `GetUsers` operation source for codegen and runtime import.
- Review and keep/update: `apps/web-admin/src/entities/user/api/createUser.graphql` - `CreateUser` operation source for codegen and runtime import.
- Review and keep/update: `apps/web-admin/src/entities/user/api/user.graphql` - `GetUser` operation source for codegen and runtime import.
- Modify: `apps/web-admin/e2e/playwright.config.ts` - Vite dev command and `VITE_GRAPHQL_API_URL`.
- Modify: `apps/web-admin/e2e/users-flow.spec.ts` - keep `/users` and `/users/:id` browser assertions.
- Remove: `apps/web-admin/app/**`, `apps/web-admin/next.config.js`, `apps/web-admin/next-env.d.ts`, `apps/web-admin/tailwind.config.ts`, `apps/web-admin/postcss.config.js` if no longer used by Vite.

### `apps/web` Next REST Public App

- Create: `apps/web/next.config.js` - Next standalone config for the deployable `web` image.
- Create: `apps/web/next-env.d.ts` - Next generated type shim.
- Create: `apps/web/app/layout.tsx` - public root layout.
- Create: `apps/web/app/page.tsx` - server-first public users page.
- Create: `apps/web/app/users-client.tsx` - client component with REST create/select/refetch.
- Create: `apps/web/app/api/users/route.ts` - runtime Next route proxy for browser same-origin list/create REST calls.
- Create: `apps/web/app/api/users/[id]/route.ts` - runtime Next route proxy for browser same-origin detail/update/delete REST calls.
- Create: `apps/web/app/__tests__/page.test.tsx` - page and layout tests.
- Create: `apps/web/app/__tests__/users-client.test.tsx` - client interaction tests.
- Create: `apps/web/app/api/users/route.test.ts` - route proxy tests.
- Create or modify: `apps/web/src/app/providers.tsx` - React Query provider.
- Create: `apps/web/vitest.config.ts` - Next-aware Vitest/jsdom config.
- Create: `apps/web/vitest.setup.ts` - jest-dom Vitest matcher setup.
- Modify: `apps/web/src/shared/config.ts` - same-origin browser REST base and runtime server `WEB_API_BASE_URL`.
- Modify: `apps/web/src/shared/config.test.ts` - Next env tests.
- Modify: `apps/web/src/shared/api/users.ts` - keep REST client, make it usable from server and browser.
- Modify: `apps/web/src/shared/api/users.test.ts` - keep REST error and URL coverage with new env.
- Create or modify: `apps/web/app/globals.css` - public CSS.
- Modify: `apps/web/package.json` - Next scripts/deps.
- Modify: `apps/web/project.json` - same Nx target names, Next-backed commands.
- Modify: `apps/web/tsconfig.json` - Next JSX/plugin/include settings.
- Modify: `apps/web/e2e/playwright.config.ts` - Next dev command and `WEB_API_BASE_URL`.
- Remove: `apps/web/index.html`, `apps/web/vite.config.ts`, `apps/web/src/main.tsx`, `apps/web/src/App.tsx`, `apps/web/src/app/App.test.tsx`, `apps/web/src/styles.css` after equivalent Next files exist.

### Tooling, Deployment, And Docs

- Modify: `tools/coverage/preflight.mjs` - required-file checks for Vite admin and Next web.
- Modify: `tools/coverage/coverage.config.json` - coverage allowlists and summaries.
- Modify: `docker/web.Dockerfile` - build `apps/web` Next standalone image.
- Modify: `docker/docker-compose.yml` - local Docker frontend service must be public `web`, not stale `web-admin` using `docker/web.Dockerfile`.
- Modify: `deploy/dokploy/docker-compose.template.yml` - run public `web` image with runtime `WEB_API_BASE_URL`; do not pretend it is admin.
- Modify: `.gitlab-ci.yml` if comments or build assumptions mention admin web image; Docker build does not need `NEXT_PUBLIC_*` build args under the runtime proxy model.
- Modify: `tools/ci/src/core.ts`, `tools/ci/src/core.test.ts`, `tools/ci/src/cli.test.ts`, `tools/ci/src/dokploy.test.ts` - deploy helper must write `WEB_API_BASE_URL` beside image vars so Dokploy env migration is automated.
- Modify if stale: `docs/infrastructure/ci-cd.md` - deployment docs must say `WEB_IMAGE` is public Next web and `WEB_API_BASE_URL` is required runtime env.
- Modify: `docs/requirements.xml`, `docs/technology.xml`, `docs/development-plan.xml`, `docs/knowledge-graph.xml`, `docs/verification-plan.xml`; update `docs/operational-packets.xml` only if current packets mention framework ownership.
- Create or update: `.tasks/swap-web-admin-vite-web-next/verification.md` - verification evidence when execution finishes.
- Modify: `bun.lock` - lockfile must be refreshed and committed with package dependency changes.

## Cross-Cutting Execution Rules

- Before the first source change, complete Task 0 so GRACE XML and `.tasks/swap-web-admin-vite-web-next/verification.md` describe the target ownership and planned evidence.
- Before every commit in Tasks 1-8, update `.tasks/swap-web-admin-vite-web-next/verification.md` with the exact commands just run and their PASS/FAIL/SKIPPED status. Do not make source-only commits after framework ownership changes.
- Every new or meaningfully edited governed file that supports comments must include or update file-local GRACE `MODULE_CONTRACT`, `MODULE_MAP`, and, where behavior/risk changed, `CHANGE_SUMMARY` before behavior code is changed.
- Code snippets below focus on behavior. When creating or replacing TS/TSX/JS/CSS/GraphQL/Docker/YAML/Markdown files, prepend the file-local GRACE header required by `AGENTS.md`; do not copy snippets without the header. JSON files cannot carry comments, so record their contract coverage in the adjacent GRACE XML and verification evidence.
- Admin GraphQL runtime code must import the existing `.graphql` documents from `apps/web-admin/src/entities/user/api/*.graphql?raw`; do not inline operation strings in route components.
- Public web must not use `NEXT_PUBLIC_API_BASE_URL` for production deployment. Browser REST calls go through same-origin `/api/users`; `WEB_API_BASE_URL` is read only by Next server/runtime code and Dokploy/CI helpers.

## Task 0: Prime GRACE Contracts And Verification Evidence

**Files:**

- Modify: `docs/requirements.xml`
- Modify: `docs/technology.xml`
- Modify: `docs/development-plan.xml`
- Modify: `docs/knowledge-graph.xml`
- Modify: `docs/verification-plan.xml`
- Modify if needed: `docs/operational-packets.xml`
- Create: `.tasks/swap-web-admin-vite-web-next/verification.md`

- [ ] **Step 1: Apply the planned GRACE ownership delta before source work**

Apply the same target ownership changes listed in Task 7 Steps 1-5 before touching app source:

- `web-admin` is a Vite React Router GraphQL SPA with generated GraphQL operation types.
- `web` is a Next App Router REST public app.
- Browser public web REST goes through same-origin `/api/users`; Next runtime proxy/server code uses `WEB_API_BASE_URL`.
- `docker/web.Dockerfile` and Dokploy `web` service deploy public Next web. Vite admin deployment remains out of scope unless a separate design adds an admin image/service.
- `V-M-WEB-ADMIN`, `V-M-WEB`, `V-M-COVERAGE-GATE`, and `V-M-CI-CD` include the planned file refs and checks from this plan, even though the files are created by later tasks.

- [ ] **Step 2: Create the evidence file with planned checks**

Create `.tasks/swap-web-admin-vite-web-next/verification.md` using the Task 7 template. Mark each command as `NOT RUN - planned before implementation` until the corresponding task executes it.

- [ ] **Step 3: Validate GRACE docs before source commits**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected: both PASS. If `grace lint --path .` has pre-existing warnings, record exact warnings in `.tasks/swap-web-admin-vite-web-next/verification.md` and distinguish them from new failures.

- [ ] **Step 4: Commit GRACE priming**

Run:

```bash
git add docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml docs/operational-packets.xml .tasks/swap-web-admin-vite-web-next/verification.md
git commit -m "docs: prime grace contracts for web swap"
```

## Task 1: Convert Admin Config And Tests To Vite Env

**Files:**

- Modify: `apps/web-admin/package.json`
- Modify: `apps/web-admin/tsconfig.json`
- Create: `apps/web-admin/vite.config.ts`
- Remove: `apps/web-admin/vitest.config.ts`
- Modify: `apps/web-admin/src/shared/config/index.ts`
- Modify: `apps/web-admin/src/shared/config/index.test.ts`
- Modify: `apps/web-admin/src/app/config.ts`
- Modify: `bun.lock`

- [ ] **Step 1: Write the failing admin Vite env tests**

Replace `apps/web-admin/src/shared/config/index.test.ts` with:

```ts
import { afterEach, describe, expect, it, vi } from 'vitest';

describe('web-admin appConfig', () => {
  afterEach(() => {
    vi.unstubAllEnvs();
    vi.resetModules();
  });

  it('defaults to the local GraphQL endpoint and app name', async () => {
    vi.stubEnv('VITE_GRAPHQL_API_URL', '');
    vi.stubEnv('VITE_APP_NAME', '');

    const { appConfig } = await import('./index');

    expect(appConfig).toEqual({
      apiUrl: 'http://localhost:8090/graphql',
      appName: 'MonorepoApp',
    });
  });

  it('uses Vite environment overrides and trims trailing slashes', async () => {
    vi.stubEnv('VITE_GRAPHQL_API_URL', ' https://api.test/graphql/// ');
    vi.stubEnv('VITE_APP_NAME', 'TemplateAdmin');

    const shared = await import('./index');
    const app = await import('../../app/config');

    expect(shared.appConfig.apiUrl).toBe('https://api.test/graphql');
    expect(shared.appConfig.appName).toBe('TemplateAdmin');
    expect(app.appConfig).toBe(shared.appConfig);
  });
});
```

- [ ] **Step 2: Run the failing admin config test**

Run:

```bash
bunx nx test web-admin -- src/shared/config/index.test.ts
```

Expected: FAIL because current config reads `NEXT_PUBLIC_API_URL` and `NEXT_PUBLIC_APP_NAME`.

- [ ] **Step 3: Implement Vite admin config**

Replace `apps/web-admin/src/shared/config/index.ts` with:

```ts
function normalizeUrl(value: string | undefined, fallback: string): string {
  const url = value?.trim() || fallback;
  return url.replace(/\/+$/, '');
}

export const appConfig = {
  apiUrl: normalizeUrl(import.meta.env.VITE_GRAPHQL_API_URL, 'http://localhost:8090/graphql'),
  appName: import.meta.env.VITE_APP_NAME?.trim() || 'MonorepoApp',
} as const;
```

Keep `apps/web-admin/src/app/config.ts` as:

```ts
export { appConfig } from '@shared/config';
```

- [ ] **Step 4: Update admin package and TypeScript config for Vite**

In `apps/web-admin/package.json`, replace scripts and dependency sets with:

```json
{
  "name": "web-admin",
  "version": "0.0.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite --host 0.0.0.0 --port 3001",
    "build": "vite build",
    "preview": "vite preview --host 0.0.0.0 --port 3001",
    "test": "vitest run --config vite.config.ts",
    "test:watch": "vitest --config vite.config.ts",
    "test:coverage": "vitest run --coverage --config vite.config.ts",
    "e2e:preflight": "node e2e/preflight.mjs",
    "codegen": "graphql-codegen --config ../../tools/codegen/codegen.ts",
    "typecheck": "tsc --noEmit --incremental false"
  },
  "dependencies": {
    "@tanstack/react-query": "^5.0.0",
    "graphql": "^16.0.0",
    "graphql-request": "^7.0.0",
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "react-router": "^7.0.0",
    "vite": "^5.4.0"
  },
  "devDependencies": {
    "@graphql-codegen/cli": "^5.0.0",
    "@graphql-codegen/typescript": "^4.0.0",
    "@graphql-codegen/typescript-operations": "^4.0.0",
    "@playwright/test": "^1.48.0",
    "@testing-library/jest-dom": "^6.0.0",
    "@testing-library/react": "^16.0.0",
    "@types/node": "25.5.0",
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "@vitejs/plugin-react": "^4.0.0",
    "@vitest/coverage-v8": "^2.1.9",
    "eslint-plugin-boundaries": "^4.0.0",
    "jsdom": "^25.0.0",
    "typescript": "^5.5.0",
    "vitest": "^2.0.0"
  }
}
```

Create `apps/web-admin/vite.config.ts` now so the package scripts use the actual Vite/Vitest config from the first Vite commit:

```ts
import react from '@vitejs/plugin-react';
import { resolve } from 'path';
import { defineConfig } from 'vitest/config';

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    globals: true,
    include: ['src/**/*.test.{ts,tsx}'],
    passWithNoTests: false,
    setupFiles: ['./vitest.setup.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json-summary'],
      reportsDirectory: '../../dist/coverage/web-admin',
      include: ['src/**/*.{ts,tsx}'],
      exclude: ['src/**/*.test.{ts,tsx}', 'src/main.tsx', 'src/shared/api/generated/**'],
      thresholds: {
        statements: 100,
        branches: 100,
        functions: 100,
        lines: 100,
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

Remove the stale Next-admin Vitest config so `vitest run` cannot pick up old `app/**` coverage paths:

```bash
git rm apps/web-admin/vitest.config.ts
```

Replace `apps/web-admin/tsconfig.json` with:

```json
{
  "extends": "../../tsconfig.base.json",
  "compilerOptions": {
    "jsx": "react-jsx",
    "allowSyntheticDefaultImports": true,
    "esModuleInterop": true,
    "forceConsistentCasingInFileNames": true,
    "strict": true,
    "noEmit": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "baseUrl": ".",
    "types": ["vite/client"],
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
  "include": ["src/**/*.ts", "src/**/*.tsx", "vite.config.ts"],
  "exclude": ["node_modules", "dist", "e2e"]
}
```

- [ ] **Step 5: Refresh and verify the workspace lockfile**

Run:

```bash
bun install
bun install --frozen-lockfile
```

Expected: `bun.lock` is updated for the admin package dependency swap, and the frozen install passes.

- [ ] **Step 6: Run the admin config test again**

Run:

```bash
bunx nx test web-admin -- src/shared/config/index.test.ts
```

Expected: PASS for both admin config tests.

- [ ] **Step 7: Commit admin env migration**

Run:

```bash
git add apps/web-admin/package.json apps/web-admin/tsconfig.json apps/web-admin/vite.config.ts apps/web-admin/vitest.config.ts apps/web-admin/src/shared/config/index.ts apps/web-admin/src/shared/config/index.test.ts apps/web-admin/src/app/config.ts bun.lock .tasks/swap-web-admin-vite-web-next/verification.md
git commit -m "feat(web-admin): move config to vite env"
```

## Task 2: Build The Admin Vite SPA Shell And Routes

**Files:**

- Create: `apps/web-admin/index.html`
- Verify/update: `apps/web-admin/vite.config.ts`
- Create: `apps/web-admin/src/main.tsx`
- Create: `apps/web-admin/src/App.tsx`
- Create: `apps/web-admin/src/app/providers.tsx`
- Create: `apps/web-admin/src/pages/home.tsx`
- Create: `apps/web-admin/src/pages/users-page.tsx`
- Create: `apps/web-admin/src/pages/user-detail-page.tsx`
- Create: `apps/web-admin/src/styles.css`
- Create: `apps/web-admin/src/App.test.tsx`
- Create: `apps/web-admin/src/pages/users-page.test.tsx`
- Create: `apps/web-admin/src/pages/user-detail-page.test.tsx`
- Keep/update: `apps/web-admin/src/entities/user/api/users.graphql`
- Keep/update: `apps/web-admin/src/entities/user/api/createUser.graphql`
- Keep/update: `apps/web-admin/src/entities/user/api/user.graphql`
- Modify: `apps/web-admin/project.json`
- Remove: `apps/web-admin/app/**`, `apps/web-admin/next.config.js`, `apps/web-admin/next-env.d.ts`, `apps/web-admin/tailwind.config.ts`, `apps/web-admin/postcss.config.js`

- [ ] **Step 1: Write failing admin route smoke tests**

Create `apps/web-admin/src/App.test.tsx`:

```tsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import App from './App';

const requestMock = vi.hoisted(() => vi.fn());

vi.mock('@shared/api/graphql-client', () => ({
  graphqlClient: {
    request: requestMock,
  },
}));

function renderApp(path = '/') {
  window.history.pushState({}, '', path);
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>,
  );
}

describe('web-admin routes', () => {
  it('renders the home route with a users link', () => {
    renderApp('/');

    expect(screen.getByRole('heading', { name: 'Monorepo Template Admin' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Users' })).toHaveAttribute('href', '/users');
  });

  it('renders the users route through the browser router', async () => {
    requestMock.mockResolvedValue({
      users: { edges: [], pageInfo: { hasNextPage: false, endCursor: null }, totalCount: 0 },
    });

    renderApp('/users');

    expect(await screen.findByText('No users yet. Create one above.')).toBeInTheDocument();
  });
});
```

- [ ] **Step 2: Run the failing route smoke tests**

Run:

```bash
bunx nx test web-admin -- src/App.test.tsx
```

Expected: FAIL because `apps/web-admin/src/App.tsx` does not exist.

- [ ] **Step 3: Write failing users route tests**

Create `apps/web-admin/src/pages/users-page.test.tsx`:

```tsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import UsersPage from './users-page';

const requestMock = vi.hoisted(() => vi.fn());

vi.mock('@shared/api/graphql-client', () => ({
  graphqlClient: {
    request: requestMock,
  },
}));

function usersResponse(edges: Array<{ cursor: string; node: unknown }>, totalCount = edges.length) {
  return {
    users: {
      edges,
      pageInfo: { hasNextPage: false, endCursor: null },
      totalCount,
    },
  };
}

function renderUsersPage() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <UsersPage />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

async function fillAndSubmit() {
  fireEvent.change(screen.getByPlaceholderText('Name'), { target: { value: 'Created User' } });
  fireEvent.change(screen.getByPlaceholderText('Email'), {
    target: { value: 'created@example.com' },
  });
  fireEvent.change(screen.getByPlaceholderText('Password'), {
    target: { value: 'Password123!' },
  });
  fireEvent.click(screen.getByRole('button', { name: 'Create' }));
}

describe('UsersPage', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('renders loading, empty, and total count states', async () => {
    requestMock.mockResolvedValue(usersResponse([]));

    renderUsersPage();

    expect(await screen.findByText('No users yet. Create one above.')).toBeInTheDocument();
    expect(screen.getByText('Total: 0')).toBeInTheDocument();
  });

  it('renders returned users as detail links', async () => {
    requestMock.mockResolvedValue(
      usersResponse([
        {
          cursor: 'cursor-1',
          node: {
            id: 'user-1',
            email: 'one@example.com',
            name: 'One User',
            createdAt: '2026-05-02T00:00:00Z',
          },
        },
      ]),
    );

    renderUsersPage();

    expect(await screen.findByRole('link', { name: 'One User' })).toHaveAttribute(
      'href',
      '/users/user-1',
    );
    expect(screen.getByText('one@example.com')).toBeInTheDocument();
    expect(screen.getByText('Total: 1')).toBeInTheDocument();
  });

  it('shows load errors', async () => {
    requestMock.mockRejectedValue(new Error('network failed'));

    renderUsersPage();

    expect(await screen.findByText('Failed to load users.')).toBeInTheDocument();
  });

  it('creates a user, invalidates the list, and clears the form', async () => {
    requestMock
      .mockResolvedValueOnce(usersResponse([]))
      .mockResolvedValueOnce({
        createUser: { __typename: 'CreateUserSuccess', user: { id: 'u2' } },
      })
      .mockResolvedValueOnce(usersResponse([]));

    renderUsersPage();
    await screen.findByText('No users yet. Create one above.');
    await fillAndSubmit();

    await waitFor(() => expect(requestMock).toHaveBeenCalledTimes(3));
    expect(screen.getByPlaceholderText('Name')).toHaveValue('');
    expect(screen.getByPlaceholderText('Email')).toHaveValue('');
    expect(screen.getByPlaceholderText('Password')).toHaveValue('');
  });

  it('shows validation and auth errors from createUser', async () => {
    requestMock.mockResolvedValueOnce(usersResponse([])).mockResolvedValueOnce({
      createUser: { __typename: 'ValidationError', field: 'email', message: 'already exists' },
    });

    renderUsersPage();
    await screen.findByText('No users yet. Create one above.');
    await fillAndSubmit();

    expect(await screen.findByText('email: already exists')).toBeInTheDocument();
  });
});
```

- [ ] **Step 4: Write failing user detail route tests**

Create `apps/web-admin/src/pages/user-detail-page.test.tsx`:

```tsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router';
import { afterEach, describe, expect, it, vi } from 'vitest';
import UserDetailPage from './user-detail-page';

const requestMock = vi.hoisted(() => vi.fn());

vi.mock('@shared/api/graphql-client', () => ({
  graphqlClient: {
    request: requestMock,
  },
}));

function renderDetail(path = '/users/user-1') {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[path]}>
        <Routes>
          <Route path="/users/:id" element={<UserDetailPage />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe('UserDetailPage', () => {
  afterEach(() => {
    vi.resetAllMocks();
  });

  it('renders a fetched user', async () => {
    requestMock.mockResolvedValue({
      user: {
        id: 'user-1',
        email: 'one@example.com',
        name: 'One User',
        createdAt: '2026-05-02T00:00:00Z',
        updatedAt: '2026-05-02T00:00:00Z',
      },
    });

    renderDetail();

    expect(await screen.findByRole('heading', { name: 'One User' })).toBeInTheDocument();
    expect(screen.getByText('one@example.com')).toBeInTheDocument();
    expect(screen.getByText('user-1')).toBeInTheDocument();
  });

  it('renders not found when GraphQL returns null', async () => {
    requestMock.mockResolvedValue({ user: null });

    renderDetail('/users/missing');

    expect(await screen.findByRole('heading', { name: 'User not found' })).toBeInTheDocument();
  });

  it('renders load failures', async () => {
    requestMock.mockRejectedValue(new Error('network failed'));

    renderDetail();

    expect(await screen.findByText('Failed to load user.')).toBeInTheDocument();
  });
});
```

- [ ] **Step 5: Run the failing admin route suites**

Run:

```bash
bunx nx test web-admin -- src/App.test.tsx src/pages/users-page.test.tsx src/pages/user-detail-page.test.tsx
```

Expected: FAIL because the Vite route files do not exist.

- [ ] **Step 6: Add HTML root, providers, route table, and styles**

Create `apps/web-admin/index.html`:

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Monorepo Template Admin</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

Keep `apps/web-admin/vite.config.ts` from Task 1 as the single admin Vitest/Vite config. If route imports require an additional alias, update this file directly; do not recreate `apps/web-admin/vitest.config.ts`.

Create `apps/web-admin/src/app/providers.tsx`:

```tsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useState, type ReactNode } from 'react';

export function Providers({ children }: { children: ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 60_000,
            refetchOnWindowFocus: false,
          },
        },
      }),
  );

  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
}
```

Create `apps/web-admin/src/App.tsx`:

```tsx
import { BrowserRouter, Navigate, Route, Routes } from 'react-router';
import HomePage from './pages/home';
import UserDetailPage from './pages/user-detail-page';
import UsersPage from './pages/users-page';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/users" element={<UsersPage />} />
        <Route path="/users/:id" element={<UserDetailPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
```

Create `apps/web-admin/src/main.tsx`:

```tsx
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import { Providers } from './app/providers';
import './styles.css';

const rootElement = document.getElementById('root');

if (!rootElement) {
  throw new Error('Root element not found');
}

createRoot(rootElement).render(
  <StrictMode>
    <Providers>
      <App />
    </Providers>
  </StrictMode>,
);
```

Create `apps/web-admin/src/pages/home.tsx`:

```tsx
import { Link } from 'react-router';

export default function HomePage() {
  return (
    <main className="page-shell centered-page">
      <h1>Monorepo Template Admin</h1>
      <p>GraphQL admin client</p>
      <Link className="primary-link" to="/users">
        Users
      </Link>
    </main>
  );
}
```

Create `apps/web-admin/src/styles.css`:

```css
:root {
  color: #172026;
  background: #f5f7f8;
  font-family:
    Inter,
    ui-sans-serif,
    system-ui,
    -apple-system,
    BlinkMacSystemFont,
    'Segoe UI',
    sans-serif;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
}

button,
input {
  font: inherit;
}

.page-shell {
  margin: 0 auto;
  max-width: 64rem;
  min-height: 100vh;
  padding: 3rem 1.5rem;
}

.centered-page {
  align-items: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.primary-link,
.text-link {
  color: #146c5f;
  font-weight: 600;
}

.toolbar {
  align-items: center;
  display: flex;
  justify-content: space-between;
  margin-bottom: 1.5rem;
}

.create-form {
  display: grid;
  gap: 0.75rem;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  margin-bottom: 1rem;
}

.create-form input,
.create-form button,
.user-card {
  border-radius: 6px;
}

.create-form input {
  background: #fff;
  border: 1px solid #c8d2d8;
  min-width: 0;
  padding: 0.75rem;
}

.create-form button {
  background: #146c5f;
  border: 0;
  color: #fff;
  cursor: pointer;
  padding: 0.75rem 1rem;
}

.create-form button:disabled {
  cursor: wait;
  opacity: 0.7;
}

.error-message {
  color: #b42318;
}

.muted {
  color: #5f6f7a;
}

.user-list {
  display: grid;
  gap: 0.75rem;
  list-style: none;
  margin: 0;
  padding: 0;
}

.user-card {
  background: #fff;
  border: 1px solid #d8e0e5;
  padding: 1rem;
}

.detail-list {
  display: grid;
  gap: 1rem;
}

@media (max-width: 760px) {
  .create-form {
    grid-template-columns: 1fr;
  }
}
```

- [ ] **Step 7: Keep GraphQL operations in codegen-visible documents**

Verify these existing documents remain the single source of GraphQL operations:

```bash
sed -n '1,160p' apps/web-admin/src/entities/user/api/users.graphql
sed -n '1,160p' apps/web-admin/src/entities/user/api/createUser.graphql
sed -n '1,160p' apps/web-admin/src/entities/user/api/user.graphql
```

Expected: documents define `GetUsers`, `CreateUser`, and `GetUser`. If operation fields need adjustment, edit these `.graphql` files instead of adding inline GraphQL strings to route components.

- [ ] **Step 8: Add users and detail route implementations**

Create `apps/web-admin/src/pages/users-page.tsx`:

```tsx
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import createUserMutationDocument from '@entities/user/api/createUser.graphql?raw';
import getUsersQueryDocument from '@entities/user/api/users.graphql?raw';
import { graphqlClient } from '@shared/api/graphql-client';
import type { CreateUserMutation, GetUsersQuery } from '@shared/api/generated/types';
import { type FormEvent, useState } from 'react';
import { Link } from 'react-router';

type FormState = {
  name: string;
  email: string;
  password: string;
};

const initialFormState: FormState = { name: '', email: '', password: '' };

export default function UsersPage() {
  const queryClient = useQueryClient();
  const [form, setForm] = useState<FormState>(initialFormState);
  const [error, setError] = useState<string | null>(null);

  const usersQuery = useQuery({
    queryKey: ['admin-users'],
    queryFn: () => graphqlClient.request<GetUsersQuery>(getUsersQueryDocument, { first: 20 }),
  });

  const mutation = useMutation({
    mutationFn: (input: FormState) =>
      graphqlClient.request<CreateUserMutation>(createUserMutationDocument, { input }),
    onSuccess: async (res) => {
      const result = res.createUser;
      if ('user' in result) {
        setForm(initialFormState);
        setError(null);
        await queryClient.invalidateQueries({ queryKey: ['admin-users'] });
        return;
      }
      if ('field' in result) {
        setError(`${result.field}: ${result.message}`);
        return;
      }
      setError(result.message);
    },
    onError: (err: Error) => setError(err.message),
  });

  function updateField(field: keyof FormState, value: string) {
    setForm((current) => ({ ...current, [field]: value }));
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    mutation.mutate(form);
  }

  const users = usersQuery.data?.users.edges || [];

  return (
    <main className="page-shell">
      <div className="toolbar">
        <h1>Users</h1>
        <Link className="text-link" to="/">
          Home
        </Link>
      </div>

      <form className="create-form" onSubmit={handleSubmit}>
        <input
          placeholder="Name"
          value={form.name}
          onChange={(event) => updateField('name', event.target.value)}
        />
        <input
          placeholder="Email"
          type="email"
          value={form.email}
          onChange={(event) => updateField('email', event.target.value)}
        />
        <input
          placeholder="Password"
          type="password"
          value={form.password}
          onChange={(event) => updateField('password', event.target.value)}
        />
        <button type="submit" disabled={mutation.isPending}>
          {mutation.isPending ? 'Creating...' : 'Create'}
        </button>
      </form>

      {error ? <p className="error-message">{error}</p> : null}
      {usersQuery.isLoading ? <p className="muted">Loading...</p> : null}
      {usersQuery.isError ? <p className="error-message">Failed to load users.</p> : null}
      {usersQuery.data ? <p className="muted">Total: {usersQuery.data.users.totalCount}</p> : null}
      {usersQuery.data && users.length === 0 ? (
        <p className="muted">No users yet. Create one above.</p>
      ) : null}

      <ul className="user-list">
        {users.map(({ node }) => (
          <li className="user-card" key={node.id}>
            <Link className="text-link" to={`/users/${node.id}`}>
              {node.name}
            </Link>
            <span className="muted"> {node.email}</span>
          </li>
        ))}
      </ul>
    </main>
  );
}
```

Create `apps/web-admin/src/pages/user-detail-page.tsx`:

```tsx
import { useQuery } from '@tanstack/react-query';
import getUserQueryDocument from '@entities/user/api/user.graphql?raw';
import { graphqlClient } from '@shared/api/graphql-client';
import type { GetUserQuery } from '@shared/api/generated/types';
import { Link, useParams } from 'react-router';

export default function UserDetailPage() {
  const { id } = useParams<{ id: string }>();
  const userQuery = useQuery({
    enabled: Boolean(id),
    queryKey: ['admin-user', id],
    queryFn: () => graphqlClient.request<GetUserQuery>(getUserQueryDocument, { id }),
  });

  const user = userQuery.data?.user || null;

  if (userQuery.isLoading) {
    return (
      <main className="page-shell">
        <p className="muted">Loading...</p>
      </main>
    );
  }

  if (userQuery.isError) {
    return (
      <main className="page-shell">
        <p className="error-message">Failed to load user.</p>
        <Link className="text-link" to="/users">
          Back to users
        </Link>
      </main>
    );
  }

  if (!user) {
    return (
      <main className="page-shell">
        <h1>User not found</h1>
        <Link className="text-link" to="/users">
          Back to users
        </Link>
      </main>
    );
  }

  return (
    <main className="page-shell">
      <div className="toolbar">
        <h1>{user.name}</h1>
        <Link className="text-link" to="/users">
          Back to users
        </Link>
      </div>

      <dl className="detail-list">
        <div>
          <dt>Email</dt>
          <dd>{user.email}</dd>
        </div>
        <div>
          <dt>Created</dt>
          <dd>{new Date(user.createdAt).toLocaleString()}</dd>
        </div>
        <div>
          <dt>Updated</dt>
          <dd>{new Date(user.updatedAt).toLocaleString()}</dd>
        </div>
        <div>
          <dt>ID</dt>
          <dd>{user.id}</dd>
        </div>
      </dl>
    </main>
  );
}
```

- [ ] **Step 8: Update Nx target commands**

Replace only the `typecheck` command in `apps/web-admin/project.json` so it uses the package script:

```json
"typecheck": {
  "executor": "nx:run-commands",
  "options": {
    "command": "cd apps/web-admin && bun run typecheck"
  }
}
```

Keep `build`, `serve`, `test`, `test-coverage`, `lint`, `codegen`, and `e2e` target names unchanged.

- [ ] **Step 9: Remove Next admin files**

Run:

```bash
git rm -r apps/web-admin/app
git rm apps/web-admin/next.config.js apps/web-admin/next-env.d.ts apps/web-admin/tailwind.config.ts apps/web-admin/postcss.config.js
```

Expected: files are removed from git tracking. If any path is already untracked or absent, skip only that absent path and keep the rest.

- [ ] **Step 10: Run admin route tests**

Run:

```bash
bunx nx test web-admin -- src/App.test.tsx src/pages/users-page.test.tsx src/pages/user-detail-page.test.tsx
```

Expected: PASS.

- [ ] **Step 11: Run admin typecheck and build**

Run:

```bash
bunx nx run web-admin:typecheck
bunx nx build web-admin
```

Expected: both commands PASS. The build should produce Vite output under `apps/web-admin/dist`.

- [ ] **Step 12: Commit admin Vite SPA**

Run:

```bash
git add apps/web-admin .tasks/swap-web-admin-vite-web-next/verification.md
git commit -m "feat(web-admin): replace next app with vite spa"
```

## Task 3: Convert Public Web To Next App Router With REST

**Files:**

- Modify: `apps/web/package.json`
- Modify: `apps/web/project.json`
- Modify: `apps/web/tsconfig.json`
- Create: `apps/web/next.config.js`
- Create: `apps/web/next-env.d.ts`
- Create: `apps/web/app/layout.tsx`
- Create: `apps/web/app/page.tsx`
- Create: `apps/web/app/users-client.tsx`
- Create: `apps/web/app/api/users/route.ts`
- Create: `apps/web/app/api/users/[id]/route.ts`
- Create: `apps/web/app/__tests__/page.test.tsx`
- Create: `apps/web/app/__tests__/users-client.test.tsx`
- Create: `apps/web/app/api/users/route.test.ts`
- Create: `apps/web/src/app/providers.tsx`
- Create: `apps/web/vitest.config.ts`
- Create: `apps/web/vitest.setup.ts`
- Modify: `apps/web/src/shared/config.ts`
- Modify: `apps/web/src/shared/config.test.ts`
- Modify: `apps/web/src/shared/api/users.ts`
- Modify: `apps/web/src/shared/api/users.test.ts`
- Create: `apps/web/app/globals.css`
- Modify: `bun.lock`
- Remove: `apps/web/index.html`, `apps/web/vite.config.ts`, `apps/web/src/main.tsx`, `apps/web/src/App.tsx`, `apps/web/src/app/App.test.tsx`, `apps/web/src/styles.css`

- [ ] **Step 1: Write failing public web config tests**

Replace `apps/web/src/shared/config.test.ts` with:

```ts
import { afterEach, describe, expect, it, vi } from 'vitest';

afterEach(() => {
  vi.unstubAllGlobals();
  vi.unstubAllEnvs();
  vi.resetModules();
});

describe('web config', () => {
  it('defaults browser REST requests to the same-origin Next proxy', async () => {
    vi.stubEnv('WEB_API_BASE_URL', 'https://api.example.test');

    const { appConfig } = await import('./config');

    expect(appConfig.apiBaseUrl).toBe('');
  });

  it('uses the runtime server API base URL when no browser window exists', async () => {
    vi.stubGlobal('window', undefined);
    vi.stubEnv('WEB_API_BASE_URL', 'https://api.example.test');

    const { appConfig } = await import('./config');

    expect(appConfig.apiBaseUrl).toBe('https://api.example.test');
  });

  it('trims whitespace and trailing slashes from the runtime server API base URL', async () => {
    vi.stubGlobal('window', undefined);
    vi.stubEnv('WEB_API_BASE_URL', ' https://api.example.test/// ');

    const { appConfig } = await import('./config');

    expect(appConfig.apiBaseUrl).toBe('https://api.example.test');
  });

  it('defaults server REST requests to the local API base URL', async () => {
    vi.stubGlobal('window', undefined);
    vi.stubEnv('WEB_API_BASE_URL', '');

    const { appConfig } = await import('./config');

    expect(appConfig.apiBaseUrl).toBe('http://localhost:8090');
  });
});
```

- [ ] **Step 2: Run the failing public config tests**

Run:

```bash
bunx nx test web -- src/shared/config.test.ts
```

Expected: FAIL because current config reads `import.meta.env.VITE_API_BASE_URL` and has no server/runtime `WEB_API_BASE_URL` branch.

- [ ] **Step 3: Implement Next public config**

Replace `apps/web/src/shared/config.ts` with:

```ts
function normalizeApiBaseUrl(value: string | undefined): string {
  const baseUrl = value?.trim() || 'http://localhost:8090';
  return baseUrl.replace(/\/+$/, '');
}

function resolveApiBaseUrl(): string {
  if (typeof window !== 'undefined') {
    return '';
  }

  return normalizeApiBaseUrl(process.env.WEB_API_BASE_URL);
}

export const appConfig = {
  apiBaseUrl: resolveApiBaseUrl(),
};
```

- [ ] **Step 4: Run the public config tests again**

Run:

```bash
bunx nx test web -- src/shared/config.test.ts
```

Expected: PASS.

- [ ] **Step 5: Write failing Next page and client tests**

Create `apps/web/app/__tests__/users-client.test.tsx`:

```tsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import UsersClient from '../users-client';
import type { User } from '../../src/shared/api/users';

function renderUsersClient(initialUsers: User[] = [], totalCount = initialUsers.length) {
  const client = new QueryClient({
    defaultOptions: {
      queries: { retry: false, staleTime: 60_000, refetchOnMount: false },
      mutations: { retry: false },
    },
  });
  return render(
    <QueryClientProvider client={client}>
      <UsersClient initialUsers={initialUsers} initialTotalCount={totalCount} />
    </QueryClientProvider>,
  );
}

afterEach(() => {
  vi.restoreAllMocks();
});

describe('UsersClient', () => {
  it('renders initial users from the server page', () => {
    renderUsersClient([
      {
        id: 'u1',
        email: 'one@example.com',
        name: 'One',
        createdAt: '2026-05-24T00:00:00Z',
        updatedAt: '2026-05-24T00:00:00Z',
      },
    ]);

    expect(screen.getByText('One')).toBeInTheDocument();
    expect(screen.getByText('one@example.com')).toBeInTheDocument();
    expect(screen.getByText('1 users')).toBeInTheDocument();
  });

  it('creates a user and refreshes the list through REST', async () => {
    vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify({
            data: {
              id: 'u1',
              email: 'new@example.com',
              name: 'New User',
              createdAt: '2026-05-24T00:00:00Z',
              updatedAt: '2026-05-24T00:00:00Z',
            },
          }),
          { status: 201 },
        ),
      )
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify({
            data: [
              {
                id: 'u1',
                email: 'new@example.com',
                name: 'New User',
                createdAt: '2026-05-24T00:00:00Z',
                updatedAt: '2026-05-24T00:00:00Z',
              },
            ],
            meta: { totalCount: 1 },
          }),
          { status: 200 },
        ),
      );

    renderUsersClient();
    fireEvent.change(screen.getByPlaceholderText('Name'), { target: { value: 'New User' } });
    fireEvent.change(screen.getByPlaceholderText('Email'), {
      target: { value: 'new@example.com' },
    });
    fireEvent.change(screen.getByPlaceholderText('Password'), { target: { value: 'secret123' } });
    fireEvent.click(screen.getByRole('button', { name: 'Create' }));

    expect(await screen.findByText('New User')).toBeInTheDocument();
    await waitFor(() => expect(screen.getByText('new@example.com')).toBeInTheDocument());
  });

  it('shows duplicate email errors from REST', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(
      new Response(
        JSON.stringify({
          error: { code: 'DUPLICATE_EMAIL', message: 'email already exists', field: 'email' },
        }),
        { status: 409 },
      ),
    );

    renderUsersClient();
    fireEvent.change(screen.getByPlaceholderText('Name'), { target: { value: 'Taken' } });
    fireEvent.change(screen.getByPlaceholderText('Email'), {
      target: { value: 'taken@example.com' },
    });
    fireEvent.change(screen.getByPlaceholderText('Password'), { target: { value: 'secret123' } });
    fireEvent.click(screen.getByRole('button', { name: 'Create' }));

    expect(await screen.findByText('email: email already exists')).toBeInTheDocument();
  });

  it('opens a selected user detail panel', () => {
    renderUsersClient([
      {
        id: 'u1',
        email: 'one@example.com',
        name: 'One',
        createdAt: '2026-05-24T00:00:00Z',
        updatedAt: '2026-05-24T00:00:00Z',
      },
    ]);

    fireEvent.click(screen.getByRole('button', { name: 'One' }));

    expect(screen.getByText('one@example.com')).toBeInTheDocument();
  });
});
```

Create `apps/web/app/__tests__/page.test.tsx`:

```tsx
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import RootLayout, { metadata } from '../layout';
import Page from '../page';

describe('public web page', () => {
  it('exports metadata and wraps children in the root layout', () => {
    expect(metadata.title).toBe('Monorepo Template');
    render(
      <RootLayout>
        <span>layout child</span>
      </RootLayout>,
    );
    expect(screen.getByText('layout child')).toBeInTheDocument();
  });

  it('renders server-fetched users into the public page', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(
        JSON.stringify({
          data: [
            {
              id: 'u1',
              email: 'one@example.com',
              name: 'One',
              createdAt: '2026-05-24T00:00:00Z',
              updatedAt: '2026-05-24T00:00:00Z',
            },
          ],
          meta: { totalCount: 1 },
        }),
        { status: 200 },
      ),
    );

    render(<RootLayout>{await Page()}</RootLayout>);

    expect(screen.getByText('REST Web')).toBeInTheDocument();
    expect(screen.getByText('One')).toBeInTheDocument();
  });

  it('renders a load error when the server REST request fails', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(
        JSON.stringify({ error: { code: 'INTERNAL_ERROR', message: 'database unavailable' } }),
        { status: 500 },
      ),
    );

    render(<RootLayout>{await Page()}</RootLayout>);

    expect(screen.getByText('Failed to load users.')).toBeInTheDocument();
  });
});
```

Create `apps/web/app/api/users/route.test.ts`:

```ts
import { afterEach, describe, expect, it, vi } from 'vitest';
import { GET, POST } from './route';
import { DELETE as DELETE_BY_ID, GET as GET_BY_ID, PATCH as PATCH_BY_ID } from './[id]/route';

afterEach(() => {
  vi.restoreAllMocks();
  vi.unstubAllEnvs();
});

describe('public web users route proxy', () => {
  it('forwards list requests to the runtime API base URL', async () => {
    vi.stubEnv('WEB_API_BASE_URL', 'https://api.example.test/');
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ data: [], meta: { totalCount: 0 } }), { status: 200 }),
    );

    const response = await GET(new Request('http://localhost/api/users'));

    expect(response.status).toBe(200);
    expect(fetch).toHaveBeenCalledWith(
      'https://api.example.test/api/users',
      expect.objectContaining({ method: 'GET' }),
    );
  });

  it('uses the default local API base URL when runtime env is absent', async () => {
    vi.stubEnv('WEB_API_BASE_URL', '');
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ data: [], meta: { totalCount: 0 } }), { status: 200 }),
    );

    await GET(new Request('http://localhost/api/users'));

    expect(fetch).toHaveBeenCalledWith(
      'http://localhost:8090/api/users',
      expect.objectContaining({ method: 'GET' }),
    );
  });

  it('forwards create bodies and detail ids', async () => {
    vi.stubEnv('WEB_API_BASE_URL', 'https://api.example.test');
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ data: { id: 'u1' } }), { status: 201 }),
    );

    await POST(
      new Request('http://localhost/api/users', {
        method: 'POST',
        body: JSON.stringify({ email: 'one@example.com', name: 'One', password: 'secret123' }),
      }),
    );
    await GET_BY_ID(new Request('http://localhost/api/users/user%2F1'), {
      params: Promise.resolve({ id: 'user/1' }),
    });

    expect(fetch).toHaveBeenNthCalledWith(
      1,
      'https://api.example.test/api/users',
      expect.objectContaining({ method: 'POST' }),
    );
    expect(fetch).toHaveBeenNthCalledWith(
      2,
      'https://api.example.test/api/users/user%2F1',
      expect.objectContaining({ method: 'GET' }),
    );
  });

  it('forwards patch bodies and delete requests for detail routes', async () => {
    vi.stubEnv('WEB_API_BASE_URL', 'https://api.example.test');
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ data: { id: 'u1' } }), { status: 200 }),
    );

    await PATCH_BY_ID(
      new Request('http://localhost/api/users/u1', {
        method: 'PATCH',
        body: JSON.stringify({ name: 'Updated' }),
      }),
      {
        params: Promise.resolve({ id: 'u1' }),
      },
    );
    await DELETE_BY_ID(new Request('http://localhost/api/users/u1'), {
      params: Promise.resolve({ id: 'u1' }),
    });

    expect(fetch).toHaveBeenNthCalledWith(
      1,
      'https://api.example.test/api/users/u1',
      expect.objectContaining({ body: JSON.stringify({ name: 'Updated' }), method: 'PATCH' }),
    );
    expect(fetch).toHaveBeenNthCalledWith(
      2,
      'https://api.example.test/api/users/u1',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });

  it('uses the default local API base URL for detail routes when runtime env is absent', async () => {
    vi.stubEnv('WEB_API_BASE_URL', '');
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(null, { status: 204 }));

    await DELETE_BY_ID(new Request('http://localhost/api/users/u1'), {
      params: Promise.resolve({ id: 'u1' }),
    });

    expect(fetch).toHaveBeenCalledWith(
      'http://localhost:8090/api/users/u1',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });
});
```

- [ ] **Step 6: Run the failing public Next tests**

Run:

```bash
bunx nx test web -- app/__tests__/users-client.test.tsx app/__tests__/page.test.tsx app/api/users/route.test.ts
```

Expected: FAIL because the Next app files do not exist.

- [ ] **Step 7: Update public package, tsconfig, and Next config**

Replace `apps/web/package.json` with:

```json
{
  "name": "web",
  "version": "0.0.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "test": "vitest run --config vitest.config.ts",
    "test:watch": "vitest --config vitest.config.ts",
    "test:coverage": "vitest run --coverage --config vitest.config.ts",
    "typecheck": "tsc --noEmit --incremental false"
  },
  "dependencies": {
    "@tanstack/react-query": "^5.0.0",
    "next": "^15.0.0",
    "react": "^19.0.0",
    "react-dom": "^19.0.0"
  },
  "devDependencies": {
    "@playwright/test": "^1.48.0",
    "@testing-library/jest-dom": "^6.0.0",
    "@testing-library/react": "^16.0.0",
    "@types/node": "25.5.0",
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "@vitest/coverage-v8": "^2.1.9",
    "jsdom": "^25.0.0",
    "typescript": "^5.5.0",
    "vitest": "^2.1.9"
  }
}
```

Replace `apps/web/tsconfig.json` with:

```json
{
  "extends": "../../tsconfig.base.json",
  "compilerOptions": {
    "jsx": "preserve",
    "allowJs": true,
    "allowSyntheticDefaultImports": true,
    "esModuleInterop": true,
    "forceConsistentCasingInFileNames": true,
    "strict": true,
    "noEmit": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "incremental": true,
    "plugins": [{ "name": "next" }],
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"],
      "@app/*": ["./src/app/*"],
      "@shared/*": ["./src/shared/*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules", ".next", "e2e"]
}
```

Create `apps/web/next.config.js`:

```js
/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'standalone',
};

module.exports = nextConfig;
```

Create `apps/web/next-env.d.ts`:

```ts
/// <reference types="next" />
/// <reference types="next/image-types/global" />

// This file is generated by Next.js conventions and kept committed for typecheck stability.
```

Create `apps/web/vitest.setup.ts`:

```ts
import '@testing-library/jest-dom/vitest';
```

Create `apps/web/vitest.config.ts`:

```ts
import { resolve } from 'path';
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'jsdom',
    include: ['app/**/*.test.{ts,tsx}', 'src/**/*.test.{ts,tsx}'],
    setupFiles: ['./vitest.setup.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json-summary'],
      reportsDirectory: '../../dist/coverage/web',
      include: ['app/**/*.{ts,tsx}', 'src/**/*.{ts,tsx}'],
      exclude: ['app/**/*.test.{ts,tsx}', 'src/**/*.test.{ts,tsx}', 'next-env.d.ts'],
      thresholds: {
        statements: 100,
        branches: 100,
        functions: 100,
        lines: 100,
      },
    },
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
      '@app': resolve(__dirname, './src/app'),
      '@shared': resolve(__dirname, './src/shared'),
    },
  },
});
```

- [ ] **Step 8: Refresh and verify the workspace lockfile**

Run:

```bash
bun install
bun install --frozen-lockfile
```

Expected: `bun.lock` is updated for the public web package dependency swap, and the frozen install passes.

- [ ] **Step 9: Implement Next REST page, provider, and client component**

Create `apps/web/src/app/providers.tsx`:

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
            staleTime: 60_000,
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
import { Providers } from '../src/app/providers';
import './globals.css';

export const metadata: Metadata = {
  title: 'Monorepo Template',
  description: 'Public REST web app',
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
import UsersClient from './users-client';
import { listUsers } from '../src/shared/api/users';

export default async function Page() {
  try {
    const initial = await listUsers();
    return (
      <main className="app-shell">
        <UsersClient initialUsers={initial.users} initialTotalCount={initial.totalCount} />
      </main>
    );
  } catch {
    return (
      <main className="app-shell">
        <section className="users-panel">
          <div className="panel-heading">
            <h1>REST Web</h1>
          </div>
          <p className="error-message">Failed to load users.</p>
        </section>
      </main>
    );
  }
}
```

Create `apps/web/app/users-client.tsx`:

```tsx
'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { type FormEvent, useState } from 'react';
import { ApiError, createUser, listUsers, type User } from '../src/shared/api/users';

type Props = {
  initialUsers: User[];
  initialTotalCount: number;
};

type FormState = {
  name: string;
  email: string;
  password: string;
};

const initialFormState: FormState = {
  name: '',
  email: '',
  password: '',
};

function formatCreateError(error: Error): string {
  if (error instanceof ApiError && error.field) {
    return `${error.field}: ${error.message}`;
  }

  return error.message;
}

export default function UsersClient({ initialUsers, initialTotalCount }: Props) {
  const queryClient = useQueryClient();
  const [form, setForm] = useState<FormState>(initialFormState);
  const [selectedUserId, setSelectedUserId] = useState<string | null>(null);
  const [createError, setCreateError] = useState<string | null>(null);

  const usersQuery = useQuery({
    queryKey: ['users'],
    queryFn: listUsers,
    initialData: { users: initialUsers, totalCount: initialTotalCount },
    staleTime: 60_000,
  });

  const createUserMutation = useMutation({
    mutationFn: createUser,
    onMutate: () => setCreateError(null),
    onSuccess: async () => {
      setForm(initialFormState);
      await queryClient.invalidateQueries({ queryKey: ['users'] });
    },
    onError: (error: Error) => setCreateError(formatCreateError(error)),
  });

  const users = usersQuery.data.users;
  const selectedUser = users.find((user) => user.id === selectedUserId) || null;

  function updateField(field: keyof FormState, value: string) {
    setForm((current) => ({ ...current, [field]: value }));
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    createUserMutation.mutate(form);
  }

  return (
    <>
      <section className="users-panel">
        <div className="panel-heading">
          <h1>REST Web</h1>
          <span>{usersQuery.data.totalCount} users</span>
        </div>

        <form className="create-form" onSubmit={handleSubmit}>
          <input
            aria-label="Name"
            placeholder="Name"
            value={form.name}
            onChange={(event) => updateField('name', event.target.value)}
          />
          <input
            aria-label="Email"
            placeholder="Email"
            type="email"
            value={form.email}
            onChange={(event) => updateField('email', event.target.value)}
          />
          <input
            aria-label="Password"
            placeholder="Password"
            type="password"
            value={form.password}
            onChange={(event) => updateField('password', event.target.value)}
          />
          <button type="submit" disabled={createUserMutation.isPending}>
            Create
          </button>
        </form>

        {createError ? <p className="error-message">{createError}</p> : null}
        {users.length === 0 ? <p className="empty-state">No users yet.</p> : null}

        <div className="user-list" aria-label="Users">
          {users.map((user) => (
            <button
              className="user-row"
              key={user.id}
              type="button"
              onClick={() => setSelectedUserId(user.id)}
            >
              <span>{user.name}</span>
              <span>{user.email}</span>
            </button>
          ))}
        </div>
      </section>

      {selectedUser ? (
        <aside className="detail-panel" aria-label="Selected user">
          <h2>{selectedUser.name}</h2>
          <p>{selectedUser.email}</p>
        </aside>
      ) : null}
    </>
  );
}
```

Create `apps/web/app/globals.css` by moving the current `apps/web/src/styles.css` contents into this file.

- [ ] **Step 10: Add runtime Next route proxy for browser REST calls**

Create `apps/web/app/api/users/route.ts`:

```ts
function apiBaseUrl(): string {
  return (process.env.WEB_API_BASE_URL?.trim() || 'http://localhost:8090').replace(/\/+$/, '');
}

async function proxyUsers(request: Request, method: 'GET' | 'POST'): Promise<Response> {
  const init: RequestInit = {
    method,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  };

  if (method === 'POST') {
    init.body = await request.text();
  }

  return fetch(`${apiBaseUrl()}/api/users`, init);
}

export function GET(request: Request) {
  return proxyUsers(request, 'GET');
}

export function POST(request: Request) {
  return proxyUsers(request, 'POST');
}
```

Create `apps/web/app/api/users/[id]/route.ts`:

```ts
type RouteContext = {
  params: Promise<{ id: string }>;
};

function apiBaseUrl(): string {
  return (process.env.WEB_API_BASE_URL?.trim() || 'http://localhost:8090').replace(/\/+$/, '');
}

async function proxyUser(
  request: Request,
  context: RouteContext,
  method: 'GET' | 'PATCH' | 'DELETE',
): Promise<Response> {
  const { id } = await context.params;
  const init: RequestInit = {
    method,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  };

  if (method === 'PATCH') {
    init.body = await request.text();
  }

  return fetch(`${apiBaseUrl()}/api/users/${encodeURIComponent(id)}`, init);
}

export function GET(request: Request, context: RouteContext) {
  return proxyUser(request, context, 'GET');
}

export function PATCH(request: Request, context: RouteContext) {
  return proxyUser(request, context, 'PATCH');
}

export function DELETE(request: Request, context: RouteContext) {
  return proxyUser(request, context, 'DELETE');
}
```

- [ ] **Step 11: Keep the REST client compatible with server and browser**

Keep `apps/web/src/shared/api/users.ts` behavior, but verify `request<T>()` uses `appConfig.apiBaseUrl` and does not reference Vite APIs:

```ts
async function request<T>(path: string, init: RequestInit): Promise<T> {
  const response = await fetch(`${appConfig.apiBaseUrl}${path}`, {
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      ...init.headers,
    },
    ...init,
  });
  // keep the current 204, error envelope, and success envelope handling
}
```

Do not import `graphql-request`, generated GraphQL types, or `.graphql` documents in `apps/web`.

Update `apps/web/src/shared/api/users.test.ts` so jsdom/browser tests expect same-origin URLs such as `/api/users`, and add or keep a config test proving server imports use `WEB_API_BASE_URL`.

- [ ] **Step 12: Remove Vite public web files**

Run:

```bash
git rm apps/web/index.html apps/web/vite.config.ts apps/web/src/main.tsx apps/web/src/App.tsx apps/web/src/app/App.test.tsx apps/web/src/styles.css
```

Expected: Vite shell files are removed after equivalent Next files are present.

- [ ] **Step 13: Run public web focused tests**

Run:

```bash
bunx nx test web -- src/shared/config.test.ts src/shared/api/users.test.ts app/__tests__/users-client.test.tsx app/__tests__/page.test.tsx app/api/users/route.test.ts
```

Expected: PASS.

- [ ] **Step 14: Run public web typecheck and build**

Run:

```bash
bunx nx run web:typecheck
bunx nx build web
```

Expected: both commands PASS. The build should produce Next output under `apps/web/.next`.

- [ ] **Step 15: Commit public Next web**

Run:

```bash
git add apps/web bun.lock .tasks/swap-web-admin-vite-web-next/verification.md
git commit -m "feat(web): replace vite app with next rest app"
```

## Task 4: Update Playwright E2E For Swapped Dev Servers

**Files:**

- Modify: `apps/web-admin/e2e/playwright.config.ts`
- Modify: `apps/web-admin/e2e/users-flow.spec.ts`
- Modify: `apps/web/e2e/playwright.config.ts`
- Modify: `apps/web/e2e/rest-users-flow.spec.ts`

- [ ] **Step 1: Update admin Playwright server env**

In `apps/web-admin/e2e/playwright.config.ts`, replace the admin browser server block with Vite-compatible command/env:

```ts
{
  command: `bun run dev -- --host 127.0.0.1 --port ${webPort}`,
  url: webBaseURL,
  reuseExistingServer: false,
  timeout: 120_000,
  cwd: webRoot,
  env: {
    VITE_GRAPHQL_API_URL: `${apiBaseURL}/graphql`,
  },
}
```

Keep the API server block unchanged except for CORS origin still using `webBaseURL`.

- [ ] **Step 2: Verify admin e2e selectors still match Vite routes**

Keep `apps/web-admin/e2e/users-flow.spec.ts` route assertions:

```ts
await page.goto('/users');
await page.getByPlaceholder('Name').fill('Browser User');
await page.getByPlaceholder('Email').fill(email);
await page.getByPlaceholder('Password').fill('Password123!');
await page.getByRole('button', { name: 'Create' }).click();
await expect(page.getByRole('link', { name: 'Browser User' })).toBeVisible();
```

If Vite renders the same labels and links from Task 2, no selector change is needed.

- [ ] **Step 3: Update public web Playwright server env**

In `apps/web/e2e/playwright.config.ts`, replace the public browser server block with Next-compatible command/env:

```ts
{
  command: `bun run dev -- --hostname 127.0.0.1 --port ${webPort}`,
  url: webBaseURL,
  reuseExistingServer: false,
  timeout: 120_000,
  cwd: webRoot,
  env: {
    WEB_API_BASE_URL: apiBaseURL,
  },
}
```

Keep the API server block, but update comments or command text only if they still refer to Vite public web.

- [ ] **Step 4: Run e2e configs in dry targeted form**

Run admin e2e if Docker test services are available:

```bash
bunx nx run web-admin:e2e
```

Expected: PASS with artifacts under `dist/test-results/web-admin-e2e` and `dist/playwright-report/web-admin`.

Run public web e2e if Docker test services are available:

```bash
bunx nx run web:e2e
```

Expected: PASS with artifacts under `dist/test-results/web-e2e` and `dist/playwright-report/web`.

If Docker or Playwright browsers are unavailable, record the environment failure in `.tasks/swap-web-admin-vite-web-next/verification.md` and do not claim e2e success.

- [ ] **Step 5: Commit e2e server updates**

Run:

```bash
git add apps/web-admin/e2e apps/web/e2e .tasks/swap-web-admin-vite-web-next/verification.md
git commit -m "test: align e2e servers with swapped web apps"
```

## Task 5: Update Coverage, Codegen, And Preflight Tooling

**Files:**

- Modify: `tools/coverage/preflight.mjs`
- Modify: `tools/coverage/coverage.config.json`
- Review: `tools/codegen/codegen.ts`
- Review: `tools/codegen/project.json`
- Modify if generated output changes: `apps/web-admin/src/shared/api/generated/types.ts`

- [ ] **Step 1: Update coverage preflight required files**

Replace the `requiredFiles` array in `tools/coverage/preflight.mjs` with:

```js
const requiredFiles = [
  'tools/coverage/coverage.config.json',
  'package.json',
  'apps/web-admin/package.json',
  'apps/web-admin/project.json',
  'apps/web-admin/vite.config.ts',
  'apps/web-admin/tsconfig.json',
  'apps/web-admin/index.html',
  'apps/web-admin/src/App.tsx',
  'apps/web-admin/src/main.tsx',
  'apps/web-admin/src/entities/user/api/users.graphql',
  'apps/web-admin/src/entities/user/api/createUser.graphql',
  'apps/web-admin/src/entities/user/api/user.graphql',
  'apps/web-admin/src/shared/api/graphql-client.ts',
  'apps/web-admin/src/shared/config/index.ts',
  'apps/web-admin/e2e/playwright.config.ts',
  'apps/web-admin/e2e/preflight.mjs',
  'apps/web/package.json',
  'apps/web/project.json',
  'apps/web/next.config.js',
  'apps/web/tsconfig.json',
  'apps/web/vitest.config.ts',
  'apps/web/vitest.setup.ts',
  'apps/web/app/page.tsx',
  'apps/web/app/users-client.tsx',
  'apps/web/app/api/users/route.ts',
  'apps/web/app/api/users/[id]/route.ts',
  'apps/web/src/shared/api/users.ts',
  'apps/web/src/shared/config.ts',
  'apps/web/e2e/playwright.config.ts',
  'docker/docker-compose.yml',
  'docker/docker-compose.test.yml',
  'docs/verification-plan.xml',
];
```

- [ ] **Step 2: Update coverage excludes**

In `apps/web-admin/vite.config.ts`, ensure coverage includes only `src/**/*.{ts,tsx}` and excludes:

```ts
exclude: [
  'src/**/*.test.{ts,tsx}',
  'src/main.tsx',
  'src/shared/api/generated/**',
],
```

Create or update `apps/web/vitest.config.ts` for Next-aware tests:

```ts
import { resolve } from 'path';
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'jsdom',
    include: ['app/**/*.test.{ts,tsx}', 'src/**/*.test.{ts,tsx}'],
    setupFiles: ['./vitest.setup.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json-summary'],
      reportsDirectory: '../../dist/coverage/web',
      include: ['app/**/*.{ts,tsx}', 'src/**/*.{ts,tsx}'],
      exclude: ['app/**/*.test.{ts,tsx}', 'src/**/*.test.{ts,tsx}', 'next-env.d.ts'],
      thresholds: {
        statements: 100,
        branches: 100,
        functions: 100,
        lines: 100,
      },
    },
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
      '@app': resolve(__dirname, './src/app'),
      '@shared': resolve(__dirname, './src/shared'),
    },
  },
});
```

Append coverage allowlist entries in `tools/coverage/coverage.config.json` for intentional generated/bootstrap exclusions:

```json
{
  "path": "apps/web-admin/src/main.tsx",
  "reason": "Vite admin bootstrap entrypoint",
  "gate": "bunx nx build web-admin && bunx nx run web-admin:e2e"
}
```

```json
{
  "path": "apps/web/next-env.d.ts",
  "reason": "Next generated type shim",
  "gate": "bunx nx run web:typecheck && bunx nx build web"
}
```

- [ ] **Step 3: Confirm GraphQL codegen stays admin-only**

Run:

```bash
rg -n "query GetUsers|mutation CreateUser|query GetUser" apps/web-admin/src/entities/user/api
if rg -n "query GetUsers|mutation CreateUser|query GetUser" apps/web-admin/src/pages; then exit 1; fi
bunx nx run web-admin:codegen
bunx nx run codegen:validate
```

Expected: the first `rg` finds all three document operations, the negated `rg` finds no inline route operation strings, codegen PASS, and `git diff -- apps/web-admin/src/shared/api/generated apps/web` shows no GraphQL generated output under `apps/web`.

- [ ] **Step 4: Run focused coverage tooling checks**

Run:

```bash
node tools/coverage/preflight.mjs
bunx nx run web-admin:test-coverage
bunx nx run web:test-coverage
```

Expected: preflight PASS; both frontend coverage targets PASS at 100 percent. If coverage fails on a real branch, add behavior tests for the uncovered branch instead of broadening allowlists.

- [ ] **Step 5: Commit coverage and codegen tooling**

Run:

```bash
git add tools/coverage apps/web-admin/vite.config.ts apps/web/vitest.config.ts apps/web/vitest.setup.ts apps/web-admin/src/shared/api/generated .tasks/swap-web-admin-vite-web-next/verification.md
git commit -m "test: update coverage and codegen gates for web swap"
```

## Task 6: Retarget Deployment To Public Next Web

**Files:**

- Modify: `docker/web.Dockerfile`
- Modify: `docker/docker-compose.yml`
- Modify: `deploy/dokploy/docker-compose.template.yml`
- Modify if stale: `.gitlab-ci.yml`
- Modify: `tools/ci/src/core.ts`
- Modify: `tools/ci/src/core.test.ts`
- Modify: `tools/ci/src/cli.test.ts`
- Modify: `tools/ci/src/dokploy.test.ts`
- Modify if stale: `docs/infrastructure/ci-cd.md`

- [ ] **Step 1: Retarget Docker web image to `apps/web`**

Replace `docker/web.Dockerfile` with:

```dockerfile
# ---- Dev stage: deployable web image is the Next.js public web app ----
FROM oven/bun:1.3.5-alpine AS dev

WORKDIR /app
COPY package.json bun.lock tsconfig.base.json .eslintrc.json ./
COPY apps/web/package.json apps/web/
COPY apps/web-admin/package.json apps/web-admin/
COPY libs/ libs/
COPY tools/ tools/
RUN bun install --frozen-lockfile --ignore-scripts
COPY apps/web/ apps/web/
WORKDIR /app/apps/web
CMD ["bun", "run", "dev"]

# ---- Build stage: deployable web image is the Next.js public web app ----
FROM oven/bun:1.3.5-alpine AS builder

WORKDIR /app
COPY package.json bun.lock tsconfig.base.json .eslintrc.json ./
COPY apps/web/package.json apps/web/
COPY apps/web-admin/package.json apps/web-admin/
COPY libs/ libs/
COPY tools/ tools/
RUN bun install --frozen-lockfile --ignore-scripts
COPY apps/web/ apps/web/
WORKDIR /app/apps/web
RUN bun run build

# ---- Production stage ----
FROM node:22-alpine AS prod

WORKDIR /app
COPY --from=builder /app/apps/web/.next/standalone/apps/web/ ./
COPY --from=builder /app/apps/web/.next/standalone/node_modules ./node_modules
COPY --from=builder /app/apps/web/.next/static ./.next/static
COPY --from=builder /app/apps/web/public ./public
EXPOSE 3000
CMD ["node", "server.js"]
```

If `apps/web/public` does not exist, create `apps/web/public/.gitkeep` before building.

- [ ] **Step 2: Update Dokploy compose service ownership**

In `deploy/dokploy/docker-compose.template.yml`, rename the frontend service from `web-admin` to `web` and replace the public env:

```yaml
web:
  image: ${WEB_IMAGE}
  pull_policy: always
  restart: unless-stopped
  environment:
    PORT: '3000'
    WEB_API_BASE_URL: ${WEB_API_BASE_URL}
  depends_on:
    - api
  expose:
    - '3000'
```

Do not add a Vite admin service in this task. The implementation wave explicitly defers admin deployment unless the user opens a separate deployment design.

- [ ] **Step 3: Update local Docker compose frontend ownership**

In `docker/docker-compose.yml`, rename the stale `web-admin` service to `web`, keep it on `docker/web.Dockerfile` target `dev`, mount `../apps/web:/app/apps/web`, and set:

```yaml
environment:
  - WEB_API_BASE_URL=http://api:8080
ports:
  - '3000:3000'
```

Remove `NEXT_PUBLIC_API_URL`, `NEXT_PUBLIC_APP_NAME`, `.next` mounts under `apps/web-admin`, and any comments that imply the local Docker frontend is the admin app.

- [ ] **Step 4: Update CI helper expectations and Dokploy env migration**

Update `tools/ci/src/core.ts` and `tools/ci/src/cli.ts` so `deploy-dokploy` requires `WEB_API_BASE_URL` and writes it beside image refs:

```ts
export function renderDokployDeployEnv(
  registryImage: string,
  imageTag: string,
  webApiBaseUrl: string,
): Record<string, string> {
  return {
    ...renderDokployImageEnv(registryImage, imageTag),
    WEB_API_BASE_URL: webApiBaseUrl,
  };
}
```

In `runCli(['deploy-dokploy', target])`, pass `requireEnv(env, 'WEB_API_BASE_URL')` into `renderDokployDeployEnv`. Update `tools/ci/src/core.test.ts`, `tools/ci/src/cli.test.ts`, and `tools/ci/src/dokploy.test.ts` so they prove:

- `WEB_API_BASE_URL` is appended when missing from existing Dokploy env.
- `WEB_API_BASE_URL` is updated when already present.
- Old `NEXT_PUBLIC_API_URL` is not required for deploy.
- Service image keys remain `API_IMAGE`, `WEB_IMAGE`, `BOT_IMAGE`.

If `.gitlab-ci.yml` or `docs/infrastructure/ci-cd.md` still says `WEB_IMAGE` is admin or `NEXT_PUBLIC_API_URL` is required, update those references to public Next web and `WEB_API_BASE_URL`.

- [ ] **Step 5: Run CI helper tests**

Run:

```bash
bunx vitest run --config tools/vitest.config.ts tools/ci/src/core.test.ts tools/ci/src/cli.test.ts tools/ci/src/dokploy.test.ts
```

Expected: PASS. If tests fail because they assert the old service name or env variable, update the expected frontend env from `NEXT_PUBLIC_API_URL` to `WEB_API_BASE_URL` and keep service image keys `API_IMAGE`, `WEB_IMAGE`, `BOT_IMAGE`.

- [ ] **Step 6: Build and run the Docker web image locally if Docker is available**

Run:

```bash
WEB_API_BASE_URL=http://host.docker.internal:8090 bun install --frozen-lockfile
docker build --target prod -f docker/web.Dockerfile -t monorepo-template-web-swap:local .
docker run --rm -d --name mt-web-swap -p 13080:3000 -e WEB_API_BASE_URL=http://host.docker.internal:8090 monorepo-template-web-swap:local
curl -fsS http://127.0.0.1:13080/ | rg "REST Web|Failed to load users"
docker rm -f mt-web-swap
```

Expected: frozen install PASS, Docker build PASS, container starts, and the homepage responds. If Docker is unavailable, record that in `.tasks/swap-web-admin-vite-web-next/verification.md` and keep `bunx nx build web` as the local non-Docker proof; do not claim runtime deploy proof without the container run.

- [ ] **Step 7: Commit deployment retarget**

Run:

```bash
git add docker/web.Dockerfile docker/docker-compose.yml deploy/dokploy/docker-compose.template.yml .gitlab-ci.yml tools/ci/src docs/infrastructure/ci-cd.md .tasks/swap-web-admin-vite-web-next/verification.md
git commit -m "chore(deploy): retarget web image to next app"
```

## Task 7: Refresh GRACE Contracts And Verification Docs

**Files:**

- Modify: `docs/requirements.xml`
- Modify: `docs/technology.xml`
- Modify: `docs/development-plan.xml`
- Modify: `docs/knowledge-graph.xml`
- Modify: `docs/verification-plan.xml`
- Modify if needed: `docs/operational-packets.xml`
- Update: `.tasks/swap-web-admin-vite-web-next/verification.md`

- [ ] **Step 1: Update requirements wording**

In `docs/requirements.xml`, make these concrete content changes:

- Project annotation says `Vite admin GraphQL app` and `Next.js public REST web app`.
- `actor-AdminUser` still uses web-admin and GraphQL.
- `actor-EndUser` says `Next.js public web app and REST user domain`.
- `UC-003` goal says `Next.js UI to REST handler`.
- `constraint-10` says public web uses REST under `/api/users` and must not add GraphQL codegen or GraphQL transport.
- Add or update a CI/deploy acceptance note that the `web` image is the public Next app and Vite admin deploy is out of scope until explicitly designed.

Do not alter unrelated product-auth non-goals or constraints.

- [ ] **Step 2: Update technology stack**

In `docs/technology.xml`, update:

- Framework description remains `Next.js 15, Vite 5, React 19`, but constraints specify:
  - `next`: public web app depends on Next.js 15 conventions.
  - `vite`: web-admin app depends on Vite conventions.
- Preferred web data client remains:

```xml
<preferred-web-data-client>@tanstack/react-query plus REST fetch client for `web`; @tanstack/react-query plus graphql-request/codegen for `web-admin`</preferred-web-data-client>
```

- Delivery shape becomes:

```xml
<shape>Nx monorepo with Go API app, Go Telegram bot app, Vite GraphQL web-admin app, Next.js REST public web app, shared Go libraries, shared GraphQL schema library, Docker Compose infrastructure, and generated artifacts.</shape>
```

- e2e policy says Playwright runs Vite `web-admin` and Next `web` sequentially.

- [ ] **Step 3: Update knowledge graph module names**

In `docs/knowledge-graph.xml`, update these modules:

```xml
<M-WEB-ADMIN NAME="ViteWebAdminApp" TYPE="UI_COMPONENT" STATUS="implemented">
  <purpose>Vite web-admin SPA using generated GraphQL client types, React Router, and React Query.</purpose>
  <path>apps/web-admin</path>
  <depends>M-GRAPHQL-SCHEMA</depends>
  <verification-ref>V-M-WEB-ADMIN</verification-ref>
  <annotations>
    <export-ViteRouter PURPOSE="Vite React Router pages for admin home, users list, create, and detail routes." />
    <export-Providers PURPOSE="Web provider wiring." />
    <export-GeneratedGraphQLClient PURPOSE="Generated GraphQL operations and types." />
  </annotations>
</M-WEB-ADMIN>
```

```xml
<M-WEB NAME="NextPublicWebApp" TYPE="UI_COMPONENT" STATUS="implemented">
  <purpose>Next.js public web app using REST users client and React Query client components.</purpose>
  <path>apps/web</path>
  <depends>M-API</depends>
  <verification-ref>V-M-WEB</verification-ref>
  <annotations>
    <export-AppRouter PURPOSE="Next.js App Router public page and layout." />
    <export-RestUsersClient PURPOSE="REST users client for `/api/users` JSON envelopes." />
    <export-Config PURPOSE="Browser same-origin REST base plus server/runtime WEB_API_BASE_URL parsing and defaults." />
  </annotations>
</M-WEB>
```

Update crosslinks so `M-CI-CD` says it builds/deploys the public Next `web` image and does not imply admin deployment is complete.

- [ ] **Step 4: Update development plan flows**

In `docs/development-plan.xml`, update:

- `M-WEB-ADMIN` name/purpose/source/tests to Vite GraphQL admin SPA.
- `M-WEB` name/purpose/config to Next REST public app.
- `DF-LOCAL-DEV` step 3 says `nx serve web-admin` starts the Vite admin app.
- `DF-USER-REST` step 1 says browser REST requests use same-origin `/api/users` through Next route handlers, while server/runtime fetches use `WEB_API_BASE_URL`.
- `DF-CI-CD-RELEASE` evidence says the `web` image is the public Next app; admin Vite deploy is not included unless a separate admin image/service is added.

- [ ] **Step 5: Update verification plan**

In `docs/verification-plan.xml`, update:

- `V-M-WEB-ADMIN` file refs from `apps/web-admin/app/**` to:

```xml
<file>apps/web-admin/src/App.tsx</file>
<file>apps/web-admin/src/pages/users-page.tsx</file>
<file>apps/web-admin/src/pages/user-detail-page.tsx</file>
<file>apps/web-admin/src/entities/user/api/users.graphql</file>
<file>apps/web-admin/src/entities/user/api/createUser.graphql</file>
<file>apps/web-admin/src/entities/user/api/user.graphql</file>
<file>apps/web-admin/src/shared/api/graphql-client.test.ts</file>
<file>apps/web-admin/src/shared/config/index.test.ts</file>
<file>apps/web-admin/e2e/graphql-contract.spec.ts</file>
<file>apps/web-admin/e2e/users-flow.spec.ts</file>
```

- `V-M-WEB` file refs replace old Vite files with:

```xml
<file>apps/web/app/page.tsx</file>
<file>apps/web/app/users-client.tsx</file>
<file>apps/web/app/api/users/route.ts</file>
<file>apps/web/app/api/users/[id]/route.ts</file>
<file>apps/web/app/__tests__/page.test.tsx</file>
<file>apps/web/app/__tests__/users-client.test.tsx</file>
<file>apps/web/app/api/users/route.test.ts</file>
<file>apps/web/src/shared/api/users.test.ts</file>
<file>apps/web/src/shared/config.test.ts</file>
<file>apps/web/e2e/rest-users-flow.spec.ts</file>
```

- `V-M-WEB` assertions say public web must use REST endpoints, browser same-origin `/api/users`, runtime `WEB_API_BASE_URL`, and no GraphQL ownership.
- Coverage allowlist entry for web-admin generated GraphQL remains under `apps/web-admin/src/shared/api/generated/**`.
- Coverage allowlist entries include `apps/web-admin/src/main.tsx` with build/e2e replacement gate and `apps/web/next-env.d.ts` with typecheck/build replacement gate.
- Coverage preflight scenario mentions `tools/coverage/preflight.mjs` checks Vite admin and Next web shapes.
- CI/CD verification mentions `docker/web.Dockerfile` builds `apps/web`, local Docker compose no longer exposes a stale `web-admin` frontend service, and Dokploy `web` service uses runtime `WEB_API_BASE_URL`.

- [ ] **Step 6: Write verification evidence file**

Update `.tasks/swap-web-admin-vite-web-next/verification.md` so it contains:

```markdown
# Swap Web Admin Vite And Public Web Next Verification

## Focused Checks

- `bunx nx test web-admin -- src/App.test.tsx src/pages/users-page.test.tsx src/pages/user-detail-page.test.tsx`
- `bunx nx run web-admin:typecheck`
- `bunx nx build web-admin`
- `bunx nx test web -- src/shared/config.test.ts src/shared/api/users.test.ts app/__tests__/users-client.test.tsx app/__tests__/page.test.tsx app/api/users/route.test.ts`
- `bunx nx run web:typecheck`
- `bunx nx build web`
- `bun install --frozen-lockfile`
- `node tools/coverage/preflight.mjs`
- `bunx nx run web-admin:codegen`
- `bunx nx run codegen:validate`

## Broad Checks

- `bun run test:coverage`
- `bun run verify:coverage`
- `bunx nx run web-admin:e2e`
- `bunx nx run web:e2e`
- `docker build --target prod -f docker/web.Dockerfile -t monorepo-template-web-swap:local .`
- `docker run --rm -d --name mt-web-swap -p 13080:3000 -e WEB_API_BASE_URL=http://host.docker.internal:8090 monorepo-template-web-swap:local && curl -fsS http://127.0.0.1:13080/ | rg "REST Web|Failed to load users" && docker rm -f mt-web-swap`
- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`
- `grace lint --path .`

## Environment Notes

- Docker availability: not checked yet
- Playwright browser availability: not checked yet
- Skipped commands: none recorded yet
```

During execution, fill the bullets with PASS/FAIL and exact failure reason. Do not mark a skipped command as passing.

- [ ] **Step 7: Run XML and GRACE validation**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected: both PASS. If `grace lint --path .` has pre-existing warnings, record exact warnings in `.tasks/swap-web-admin-vite-web-next/verification.md` and distinguish them from new failures.

- [ ] **Step 8: Commit GRACE docs**

Run:

```bash
git add docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml docs/operational-packets.xml .tasks/swap-web-admin-vite-web-next/verification.md
git commit -m "docs: sync grace contracts for web swap"
```

## Task 8: Final Integration Gates And Hygiene

**Files:**

- Review: all files changed by Tasks 1-7.
- Modify if needed: `.tasks/swap-web-admin-vite-web-next/verification.md`

- [ ] **Step 1: Check working tree and generated drift**

Run:

```bash
git status --short
bun install --frozen-lockfile
bun run codegen
git diff --exit-code -- apps/api/internal/repository/postgres/generated apps/api/internal/graph apps/web-admin/src/shared/api/generated
```

Expected: frozen lockfile check PASS and generated diff command exits 0 after committed generated artifacts are current.

- [ ] **Step 2: Run focused frontend gates**

Run:

```bash
bunx nx test web-admin
bunx nx run web-admin:typecheck
bunx nx build web-admin
bunx nx test web
bunx nx run web:typecheck
bunx nx build web
```

Expected: all PASS.

- [ ] **Step 3: Run codegen and coverage preflight gates**

Run:

```bash
bunx nx run graphql:validate
bunx nx run api:codegen
bunx nx run web-admin:codegen
bunx nx run codegen:validate
node tools/coverage/preflight.mjs
```

Expected: all PASS.

- [ ] **Step 4: Run deployment and CI helper focused checks**

Run:

```bash
bunx vitest run --config tools/vitest.config.ts tools/ci/src/core.test.ts tools/ci/src/cli.test.ts tools/ci/src/dokploy.test.ts
docker build --target prod -f docker/web.Dockerfile -t monorepo-template-web-swap:local .
docker run --rm -d --name mt-web-swap -p 13080:3000 -e WEB_API_BASE_URL=http://host.docker.internal:8090 monorepo-template-web-swap:local
curl -fsS http://127.0.0.1:13080/ | rg "REST Web|Failed to load users"
docker rm -f mt-web-swap
```

Expected: Vitest PASS. Docker build and runtime smoke PASS if Docker is available; otherwise record the environment blocker. If the container starts but `curl` fails, remove it with `docker rm -f mt-web-swap`, record the failure, and do not claim runtime deploy proof.

- [ ] **Step 5: Run broad gates when environment allows**

Run:

```bash
bun run lint
bun run test:coverage
bun run verify:coverage
bunx nx run web-admin:e2e
bunx nx run web:e2e
bun run build
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected: all PASS. If Docker/Playwright environment blocks e2e or coverage, record the exact command, exact error, and reason in `.tasks/swap-web-admin-vite-web-next/verification.md`.

- [ ] **Step 6: Scan file-local GRACE markup on touched governed files**

Run:

```bash
touched_governed_files=(
  apps/web-admin/vite.config.ts
  apps/web-admin/src/shared/config/index.ts
  apps/web-admin/src/shared/config/index.test.ts
  apps/web-admin/src/app/config.ts
  apps/web-admin/src/main.tsx
  apps/web-admin/src/App.tsx
  apps/web-admin/src/app/providers.tsx
  apps/web-admin/src/pages/home.tsx
  apps/web-admin/src/pages/users-page.tsx
  apps/web-admin/src/pages/user-detail-page.tsx
  apps/web-admin/src/entities/user/api/users.graphql
  apps/web-admin/src/entities/user/api/createUser.graphql
  apps/web-admin/src/entities/user/api/user.graphql
  apps/web-admin/src/pages/users-page.test.tsx
  apps/web-admin/src/pages/user-detail-page.test.tsx
  apps/web-admin/src/App.test.tsx
  apps/web-admin/src/styles.css
  apps/web-admin/e2e/playwright.config.ts
  apps/web-admin/e2e/users-flow.spec.ts
  apps/web/next.config.js
  apps/web/app/layout.tsx
  apps/web/app/page.tsx
  apps/web/app/users-client.tsx
  apps/web/app/api/users/route.ts
  apps/web/app/api/users/[id]/route.ts
  apps/web/app/__tests__/page.test.tsx
  apps/web/app/__tests__/users-client.test.tsx
  apps/web/app/api/users/route.test.ts
  apps/web/src/app/providers.tsx
  apps/web/src/shared/config.ts
  apps/web/src/shared/config.test.ts
  apps/web/src/shared/api/users.ts
  apps/web/src/shared/api/users.test.ts
  apps/web/app/globals.css
  apps/web/vitest.config.ts
  apps/web/vitest.setup.ts
  apps/web/e2e/playwright.config.ts
  apps/web/e2e/rest-users-flow.spec.ts
  tools/coverage/preflight.mjs
  docker/web.Dockerfile
  docker/docker-compose.yml
  deploy/dokploy/docker-compose.template.yml
  .gitlab-ci.yml
  docs/infrastructure/ci-cd.md
  tools/ci/src/core.ts
  tools/ci/src/core.test.ts
  tools/ci/src/cli.ts
  tools/ci/src/cli.test.ts
  tools/ci/src/dokploy.test.ts
)

for file in "${touched_governed_files[@]}"; do
  rg -q "START_MODULE_CONTRACT" "$file" || {
    echo "Missing file-local GRACE MODULE_CONTRACT: $file" >&2
    exit 1
  }
done
```

Expected: no missing file-local GRACE contract is reported. JSON/package/tsconfig files are covered by the Task 0 and Task 7 GRACE XML/evidence updates because JSON cannot carry comments.

- [ ] **Step 7: Scan for forbidden public web GraphQL ownership and stale env**

Run:

```bash
if rg -n "graphql-request|\\.graphql|generated/types|VITE_API_BASE_URL|import\\.meta\\.env|NEXT_PUBLIC_API_URL|NEXT_PUBLIC_API_BASE_URL" apps/web; then exit 1; fi
if rg -n "NEXT_PUBLIC_API_URL|NEXT_PUBLIC_API_BASE_URL|process\\.env\\.NEXT_PUBLIC" apps/web-admin; then exit 1; fi
rg -n "query GetUsers|mutation CreateUser|query GetUser" apps/web-admin/src/entities/user/api
if rg -n "query GetUsers|mutation CreateUser|query GetUser" apps/web-admin/src/pages; then exit 1; fi
```

Expected: public web has no GraphQL ownership, no Vite env, and no `NEXT_PUBLIC_*` API env; admin has no old Next public API env; GraphQL operations are present in admin `.graphql` documents and absent from route components.

- [ ] **Step 8: Update verification evidence and commit final hygiene**

Update `.tasks/swap-web-admin-vite-web-next/verification.md` with the exact command results from Steps 1-7.

Run:

```bash
git add .tasks/swap-web-admin-vite-web-next/verification.md
git commit -m "chore: record web swap verification"
```

If no evidence file changes remain because Task 7 already captured final results, skip this commit and state that the evidence was already current.

## Self-Review

### Spec Coverage

- Admin Vite SPA, React Router, GraphQL/codegen ownership: Tasks 1, 2, 4, 5.
- Public Next App Router REST ownership: Task 3 and Task 4.
- Project names and Nx command names preserved: Tasks 1, 2, 3.
- Admin Vite env and public Next runtime env/proxy: Tasks 1, 3, 4, 7, 8.
- `tools/coverage/preflight.mjs` and coverage allowlists: Task 5.
- Deployment ownership for `docker/web.Dockerfile`, CI, and Dokploy: Task 6.
- GRACE XML synchronization: Task 0 primes docs before source work; Task 7 refreshes final refs and evidence.
- Final focused and broad gates: Task 8.

### Placeholder Scan

This plan avoids open implementation placeholders. The verification evidence template starts with explicit `not checked yet` values so execution can replace them with real command output; the task text forbids treating skipped commands as passing.

### Type Consistency

- Admin generated GraphQL types use `GetUsersQuery`, `GetUserQuery`, and `CreateUserMutation` from `@shared/api/generated/types`.
- Public REST types use `User`, `CreateUserInput`, `UpdateUserInput`, and `ApiError` from `apps/web/src/shared/api/users.ts`.
- Env names are consistent: `VITE_GRAPHQL_API_URL` for admin GraphQL and `WEB_API_BASE_URL` for Next server/runtime public REST proxying; browser public web uses same-origin `/api/users`.
- Nx target names remain `web-admin` and `web`.
