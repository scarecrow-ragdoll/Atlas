# Web Admin GraphQL and REST Web Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Split the template into a GraphQL-only admin frontend (`web-admin`) and a REST-based public frontend (`web`) while keeping one shared Go API domain implementation.

**Architecture:** Rename the current Next.js GraphQL app to `web-admin`, create a new React + Vite `web` app, and expose REST user handlers from the Go API beside the existing admin GraphQL endpoint. Both transports call `service.UserService`; GraphQL and REST only map request and response shapes.

**Tech Stack:** Go 1.25, chi, gqlgen, Nx 20, Bun, Next.js 15, Vite, React 19, React Query, Vitest, Playwright, GRACE XML.

---

## Source Spec

- Design: `docs/superpowers/specs/2026-05-24-web-admin-graphql-and-rest-web-design.md`
- Approved direction:
  - GraphQL belongs to the admin example only.
  - Public `web` is React + Vite and uses REST.
  - API serves `/graphql` and `/api/users*` simultaneously.
  - User business logic remains shared in `service.UserService`.

## File Structure

- Modify `apps/api/internal/service/user_service.go`: introduce stable domain error helpers if repository errors need transport-neutral mapping.
- Modify `apps/api/internal/graph/schema.resolvers.go`: map shared domain errors to GraphQL unions.
- Create `apps/api/internal/handler/users.go`: REST users handler.
- Create `apps/api/internal/handler/users_test.go`: REST handler tests.
- Modify `apps/api/cmd/server/main.go`: mount REST routes and keep GraphQL mounted.
- Move `apps/web` to `apps/web-admin`: preserve current Next.js GraphQL app.
- Modify `apps/web-admin/project.json`: rename Nx project to `web-admin` and point commands to `apps/web-admin`.
- Modify `apps/web-admin/package.json`: rename package to `web-admin`.
- Modify `apps/web-admin/vitest.config.ts`: write coverage to `dist/coverage/web-admin`.
- Modify `apps/web-admin/e2e/playwright.config.ts`: use `web-admin` commands and GraphQL API URL.
- Modify `tools/codegen/codegen.ts`: move generated GraphQL documents/output to `apps/web-admin`.
- Modify `package.json`: change codegen/e2e/typecheck references from `web` to `web-admin` where GraphQL/admin-specific.
- Create new `apps/web` Vite app files: `package.json`, `project.json`, `index.html`, `vite.config.ts`, `tsconfig.json`, `.eslintrc.json`, `src/main.tsx`, `src/App.tsx`, `src/app/App.test.tsx`, `src/shared/api/users.ts`, `src/shared/api/users.test.ts`, `src/shared/config.ts`, `src/shared/config.test.ts`, `src/styles.css`.
- Create new `apps/web/e2e/playwright.config.ts`, `apps/web/e2e/helpers.ts`, and `apps/web/e2e/rest-users-flow.spec.ts` for the public REST browser path.
- Modify `tools/coverage/coverage.config.json`: track both web coverage summaries and move GraphQL generated allowlist to `apps/web-admin`.
- Modify GRACE XML docs: `docs/requirements.xml`, `docs/technology.xml`, `docs/development-plan.xml`, `docs/knowledge-graph.xml`, `docs/verification-plan.xml`.

## Task 1: REST API Adapter On Shared User Service

**Files:**

- Modify: `apps/api/internal/service/user_service.go`
- Modify: `apps/api/internal/graph/schema.resolvers.go`
- Create: `apps/api/internal/handler/users.go`
- Create: `apps/api/internal/handler/users_test.go`
- Modify: `apps/api/cmd/server/main.go`

- [ ] **Step 1: Write failing REST handler tests with a fake service**

Create `apps/api/internal/handler/users_test.go` with this structure and test names:

```go
package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/handler"
	"monorepo-template/apps/api/internal/service"
)

func TestUsersHandler_ListUsers_ReturnsDataAndTotalCount(t *testing.T) {
	svc := newFakeUsersService()
	svc.users["u1"] = &service.User{ID: "u1", Email: "one@example.com", Name: "One", CreatedAt: "2026-05-24T00:00:00Z", UpdatedAt: "2026-05-24T00:00:00Z"}
	rec := serveUsersRequest(t, svc, http.MethodGet, "/api/users", nil)
	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"data":[{"id":"u1","email":"one@example.com","name":"One","createdAt":"2026-05-24T00:00:00Z","updatedAt":"2026-05-24T00:00:00Z"}],"meta":{"totalCount":1}}`, rec.Body.String())
}

func TestUsersHandler_CreateUser_ReturnsCreatedUser(t *testing.T) {
	svc := newFakeUsersService()
	rec := serveUsersRequest(t, svc, http.MethodPost, "/api/users", map[string]string{"email": "new@example.com", "name": "New", "password": "secret123"})
	require.Equal(t, http.StatusCreated, rec.Code)
	require.JSONEq(t, `{"data":{"id":"created-id","email":"new@example.com","name":"New","createdAt":"2026-05-24T00:00:00Z","updatedAt":"2026-05-24T00:00:00Z"}}`, rec.Body.String())
}

func TestUsersHandler_CreateUser_MapsDuplicateEmail(t *testing.T) {
	svc := newFakeUsersService()
	svc.createErr = service.ErrDuplicateEmail
	rec := serveUsersRequest(t, svc, http.MethodPost, "/api/users", map[string]string{"email": "taken@example.com", "name": "Taken", "password": "secret123"})
	require.Equal(t, http.StatusConflict, rec.Code)
	require.JSONEq(t, `{"error":{"code":"DUPLICATE_EMAIL","message":"email already exists","field":"email"}}`, rec.Body.String())
}

func TestUsersHandler_GetUser_ReturnsNotFound(t *testing.T) {
	svc := newFakeUsersService()
	rec := serveUsersRequest(t, svc, http.MethodGet, "/api/users/missing", nil)
	require.Equal(t, http.StatusNotFound, rec.Code)
	require.JSONEq(t, `{"error":{"code":"NOT_FOUND","message":"user not found"}}`, rec.Body.String())
}

func TestUsersHandler_UpdateUser_UpdatesNameAndEmail(t *testing.T) {
	svc := newFakeUsersService()
	svc.users["u1"] = &service.User{ID: "u1", Email: "old@example.com", Name: "Old", CreatedAt: "2026-05-24T00:00:00Z", UpdatedAt: "2026-05-24T00:00:00Z"}
	rec := serveUsersRequest(t, svc, http.MethodPatch, "/api/users/u1", map[string]string{"email": "new@example.com", "name": "New"})
	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"data":{"id":"u1","email":"new@example.com","name":"New","createdAt":"2026-05-24T00:00:00Z","updatedAt":"2026-05-24T00:00:01Z"}}`, rec.Body.String())
}

func TestUsersHandler_DeleteUser_ReturnsNoContentForExistingUser(t *testing.T) {
	svc := newFakeUsersService()
	svc.users["u1"] = &service.User{ID: "u1", Email: "one@example.com", Name: "One"}
	rec := serveUsersRequest(t, svc, http.MethodDelete, "/api/users/u1", nil)
	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Empty(t, rec.Body.String())
}
```

Add fake service helpers in the same test file:

```go
type fakeUsersService struct {
	users     map[string]*service.User
	createErr error
	updateErr error
	deleteErr error
}

func newFakeUsersService() *fakeUsersService {
	return &fakeUsersService{users: map[string]*service.User{}}
}

func serveUsersRequest(t *testing.T, svc *fakeUsersService, method string, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		payload, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(payload)
	}
	req := httptest.NewRequest(method, path, reader)
	rec := httptest.NewRecorder()
	handler.NewUsersHandler(svc, zap.NewNop()).Routes().ServeHTTP(rec, req)
	return rec
}
```

Then implement `List`, `Create`, `GetByID`, `Update`, and `Delete` on `fakeUsersService`; `Create` returns `created-id`, `Update` changes `UpdatedAt` to `2026-05-24T00:00:01Z`, and missing users return `nil, nil`.

- [ ] **Step 2: Run the failing API handler tests**

Run:

```bash
cd apps/api && go test ./internal/handler -run TestUsersHandler
```

Expected: FAIL with `undefined: handler.NewUsersHandler` or the equivalent missing REST handler symbol.

- [ ] **Step 3: Add shared domain error helpers**

Add exported helpers in `apps/api/internal/service/user_service.go`:

```go
var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrNotFound       = errors.New("not found")
)

func IsDuplicateEmail(err error) bool {
	return errors.Is(err, ErrDuplicateEmail) || strings.Contains(err.Error(), "duplicate email")
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
```

Keep repository compatibility by wrapping duplicate-email repository errors when they pass through service methods:

```go
if err != nil {
	if IsDuplicateEmail(err) {
		return nil, fmt.Errorf("%s: %w", op, ErrDuplicateEmail)
	}
	return nil, fmt.Errorf("%s: %w", op, err)
}
```

- [ ] **Step 4: Implement REST users handler with explicit response envelopes**

Create `apps/api/internal/handler/users.go` with:

```go
type UsersService interface {
	GetByID(ctx context.Context, id string) (*service.User, error)
	List(ctx context.Context, first *int, after *string) ([]*service.User, int, error)
	Create(ctx context.Context, input service.CreateUserInput) (*service.User, error)
	Update(ctx context.Context, id string, input service.UpdateUserInput) (*service.User, error)
	Delete(ctx context.Context, id string) error
}

type UsersHandler struct {
	service UsersService
	logger  *zap.Logger
}

func NewUsersHandler(service UsersService, logger *zap.Logger) *UsersHandler {
	return &UsersHandler{service: service, logger: logger}
}
```

Expose `Routes() chi.Router` and implement:

- `GET /` -> list users.
- `POST /` -> create user.
- `GET /{id}` -> get user.
- `PATCH /{id}` -> update user.
- `DELETE /{id}` -> check existence with `GetByID`, delete with `Delete`, return `204`.

Use response helpers named `writeJSON`, `writeError`, `mapUserResponse`, and `decodeJSON`. Error codes must be `DUPLICATE_EMAIL`, `NOT_FOUND`, `BAD_REQUEST`, and `INTERNAL_ERROR`.

- [ ] **Step 5: Mount REST routes in API main**

In `apps/api/cmd/server/main.go`, mount REST before GraphQL:

```go
usersHandler := healthHandler.NewUsersHandler(userService, l)
r.Mount("/api/users", usersHandler.Routes())
r.Handle("/graphql", srv)
```

Keep `/graphql` and `/playground` behavior unchanged.

- [ ] **Step 6: Update GraphQL resolver error mapping**

Replace resolver string matching with `service.IsDuplicateEmail(err)` and preserve current GraphQL union response shapes.

- [ ] **Step 7: Run API verification**

Run:

```bash
bunx nx test api
bunx nx lint api
bunx nx build api
```

Expected: PASS.

- [ ] **Step 8: Commit API adapter work**

Run:

```bash
git add apps/api/internal/service/user_service.go apps/api/internal/graph/schema.resolvers.go apps/api/internal/handler/users.go apps/api/internal/handler/users_test.go apps/api/cmd/server/main.go
git commit -m "feat(api): add REST user adapter"
```

## Task 2: Rename Existing GraphQL Web App To Web Admin

**Files:**

- Move: `apps/web` to `apps/web-admin`
- Modify: `apps/web-admin/project.json`
- Modify: `apps/web-admin/package.json`
- Modify: `apps/web-admin/vitest.config.ts`
- Modify: `apps/web-admin/e2e/playwright.config.ts`
- Modify: `tools/codegen/codegen.ts`
- Modify: `package.json`
- Modify: `tools/coverage/coverage.config.json`

- [ ] **Step 1: Move the existing app**

Run:

```bash
mv apps/web apps/web-admin
```

- [ ] **Step 2: Rename project and package ownership**

Update `apps/web-admin/project.json`:

- `"name": "web-admin"`
- `"sourceRoot": "apps/web-admin/src"`
- all `cd apps/web` commands become `cd apps/web-admin`.

Update `apps/web-admin/package.json`:

- `"name": "web-admin"`
- keep Next.js and GraphQL dependencies.

- [ ] **Step 3: Move GraphQL codegen ownership**

Update `tools/codegen/codegen.ts` so documents and generated output point at `apps/web-admin`.

Update root `package.json`:

```json
"codegen": "bunx nx run-many --target=codegen --projects=api,web-admin",
"test:e2e": "bunx nx run web-admin:e2e",
"verify:coverage": "bun run lint && bun run codegen && bunx nx run web-admin:typecheck && bunx nx run web:typecheck && bun run build && bun run test:coverage && bun run test:e2e && xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml && grace lint --path ."
```

- [ ] **Step 4: Update admin coverage paths**

Update `apps/web-admin/vitest.config.ts` reports directory to `../../dist/coverage/web-admin`.

Update `tools/coverage/coverage.config.json`:

- rename the existing TypeScript summary to `web-admin`.
- change generated allowlist path to `apps/web-admin/src/shared/api/generated/**`.
- change its gate to `bunx nx run web-admin:codegen && bunx nx run web-admin:typecheck`.

- [ ] **Step 5: Run admin verification**

Run:

```bash
bunx nx run web-admin:codegen
bunx nx test web-admin
bunx nx run web-admin:typecheck
bunx nx build web-admin
```

Expected: PASS.

- [ ] **Step 6: Commit admin rename**

Run:

```bash
git add package.json tools/codegen/codegen.ts tools/coverage/coverage.config.json apps/web-admin
git add -u apps/web
git commit -m "refactor(web): rename GraphQL app to web-admin"
```

## Task 3: Add New React + Vite REST Web App

**Files:**

- Create: `apps/web/package.json`
- Create: `apps/web/project.json`
- Create: `apps/web/index.html`
- Create: `apps/web/vite.config.ts`
- Create: `apps/web/tsconfig.json`
- Create: `apps/web/.eslintrc.json`
- Create: `apps/web/src/main.tsx`
- Create: `apps/web/src/App.tsx`
- Create: `apps/web/src/app/App.test.tsx`
- Create: `apps/web/src/shared/api/users.ts`
- Create: `apps/web/src/shared/api/users.test.ts`
- Create: `apps/web/src/shared/config.ts`
- Create: `apps/web/src/shared/config.test.ts`
- Create: `apps/web/src/styles.css`
- Modify: `tools/coverage/coverage.config.json`

- [ ] **Step 1: Write failing REST client and UI tests**

Create `apps/web/src/shared/config.test.ts`:

```ts
import { describe, expect, it, vi } from 'vitest';

describe('web config', () => {
  it('defaults to the local API base URL', async () => {
    vi.stubEnv('VITE_API_BASE_URL', '');
    vi.resetModules();
    const { appConfig } = await import('./config');
    expect(appConfig.apiBaseUrl).toBe('http://localhost:8080');
  });
});
```

Create `apps/web/src/shared/api/users.test.ts`:

```ts
import { afterEach, describe, expect, it, vi } from 'vitest';
import { createUser, listUsers } from './users';

afterEach(() => vi.restoreAllMocks());

describe('REST users client', () => {
  it('loads users from the REST endpoint', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ data: [], meta: { totalCount: 0 } }), { status: 200 }),
    );
    await expect(listUsers()).resolves.toEqual({ users: [], totalCount: 0 });
    expect(fetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/users',
      expect.objectContaining({ method: 'GET' }),
    );
  });

  it('throws API errors with code, message, and field', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(
        JSON.stringify({
          error: { code: 'DUPLICATE_EMAIL', message: 'email already exists', field: 'email' },
        }),
        { status: 409 },
      ),
    );
    await expect(
      createUser({ email: 'taken@example.com', name: 'Taken', password: 'secret123' }),
    ).rejects.toMatchObject({ code: 'DUPLICATE_EMAIL', field: 'email' });
  });
});
```

Create `apps/web/src/app/App.test.tsx` with this starter coverage:

```tsx
import '@testing-library/jest-dom/vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, describe, expect, it, vi } from 'vitest';
import App from '../App';

function renderApp() {
  const client = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
  return render(
    <QueryClientProvider client={client}>
      <App />
    </QueryClientProvider>,
  );
}

afterEach(() => vi.restoreAllMocks());

describe('REST web app', () => {
  it('renders the empty users state', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ data: [], meta: { totalCount: 0 } }), { status: 200 }),
    );
    renderApp();
    expect(await screen.findByText('No users yet.')).toBeInTheDocument();
  });

  it('creates a user and shows it in the list', async () => {
    vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(
        new Response(JSON.stringify({ data: [], meta: { totalCount: 0 } }), { status: 200 }),
      )
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
    renderApp();
    await userEvent.type(screen.getByPlaceholderText('Name'), 'New User');
    await userEvent.type(screen.getByPlaceholderText('Email'), 'new@example.com');
    await userEvent.type(screen.getByPlaceholderText('Password'), 'secret123');
    await userEvent.click(screen.getByRole('button', { name: 'Create' }));
    expect(await screen.findByText('New User')).toBeInTheDocument();
  });

  it('shows duplicate email errors from REST', async () => {
    vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(
        new Response(JSON.stringify({ data: [], meta: { totalCount: 0 } }), { status: 200 }),
      )
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify({
            error: { code: 'DUPLICATE_EMAIL', message: 'email already exists', field: 'email' },
          }),
          { status: 409 },
        ),
      );
    renderApp();
    await userEvent.type(screen.getByPlaceholderText('Name'), 'Taken');
    await userEvent.type(screen.getByPlaceholderText('Email'), 'taken@example.com');
    await userEvent.type(screen.getByPlaceholderText('Password'), 'secret123');
    await userEvent.click(screen.getByRole('button', { name: 'Create' }));
    expect(await screen.findByText('email: email already exists')).toBeInTheDocument();
  });

  it('opens a selected user detail panel', async () => {
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
    renderApp();
    await userEvent.click(await screen.findByRole('button', { name: 'One' }));
    await waitFor(() => expect(screen.getByText('one@example.com')).toBeInTheDocument());
  });
});
```

- [ ] **Step 2: Run the failing web tests**

Run:

```bash
bunx nx test web
```

Expected: FAIL because the new `web` project and source files do not exist yet.

- [ ] **Step 3: Add Vite project config**

Create a Vite app with React plugin, jsdom Vitest environment, and coverage output at `../../dist/coverage/web`.

Use `VITE_API_BASE_URL` for the REST API base URL.

- [ ] **Step 4: Implement typed REST client**

Implement `listUsers`, `createUser`, `getUser`, `updateUser`, and `deleteUser` in `apps/web/src/shared/api/users.ts`.

Use typed response envelopes:

```ts
type ApiEnvelope<T> = { data: T };
type ApiListEnvelope<T> = { data: T[]; meta: { totalCount: number } };
type ApiErrorEnvelope = { error: { code: string; message: string; field?: string } };
```

- [ ] **Step 5: Implement the Vite UI**

Implement `apps/web/src/App.tsx` with these user-visible states and labels so the unit and e2e tests have stable locators:

- heading: `REST Web`
- input placeholders: `Name`, `Email`, `Password`
- submit button: `Create`
- empty state: `No users yet.`
- each user row is a button named by the user's `name`
- detail panel renders the selected user's `email`
- duplicate email error format: `email: email already exists`
- load failure message: `Failed to load users.`

- [ ] **Step 6: Add web coverage summary**

Ensure `tools/coverage/coverage.config.json` includes both:

- `dist/coverage/web-admin/coverage-summary.json`
- `dist/coverage/web/coverage-summary.json`

- [ ] **Step 7: Run public web verification**

Run:

```bash
bunx nx test web
bunx nx run web:typecheck
bunx nx build web
```

Expected: PASS.

- [ ] **Step 8: Commit public web app**

Run:

```bash
git add apps/web tools/coverage/coverage.config.json
git commit -m "feat(web): add Vite REST user example"
```

## Task 4: Public REST Browser E2E

**Files:**

- Create: `apps/web/e2e/playwright.config.ts`
- Create: `apps/web/e2e/helpers.ts`
- Create: `apps/web/e2e/rest-users-flow.spec.ts`
- Modify: `package.json`

- [ ] **Step 1: Write failing REST browser e2e spec**

Create `apps/web/e2e/rest-users-flow.spec.ts`:

```ts
import { expect, test } from '@playwright/test';
import { uniqueEmail } from './helpers';

test('public REST web creates, lists, and opens a user detail', async ({ page }) => {
  const email = uniqueEmail('rest-web');
  await page.goto('/');
  await page.getByPlaceholder('Name').fill('REST Web User');
  await page.getByPlaceholder('Email').fill(email);
  await page.getByPlaceholder('Password').fill('secret123');
  await page.getByRole('button', { name: 'Create' }).click();
  await expect(page.getByText('REST Web User')).toBeVisible();
  await page.getByRole('button', { name: 'REST Web User' }).click();
  await expect(page.getByText(email)).toBeVisible();
});
```

- [ ] **Step 2: Run the failing REST e2e target**

Run:

```bash
bunx nx run web:e2e
```

Expected: FAIL because `web:e2e` is not configured.

- [ ] **Step 3: Add public web Playwright config and helpers**

Create `apps/web/e2e/playwright.config.ts` using API port `18080` and web port `13000`, matching the existing isolated test infrastructure. The web server command should be:

```ts
command: 'cd ../.. && bunx nx serve web';
```

Create `apps/web/e2e/helpers.ts`:

```ts
export function uniqueEmail(prefix: string) {
  return `${prefix}-${Date.now()}-${Math.random().toString(16).slice(2)}@example.com`;
}
```

- [ ] **Step 4: Add the web e2e Nx target and root e2e script**

Update `apps/web/project.json` with:

```json
"e2e": {
  "executor": "nx:run-commands",
  "options": {
    "command": "cd apps/web && bunx playwright test --config e2e/playwright.config.ts"
  }
}
```

Update root `package.json`:

```json
"test:e2e": "bunx nx run-many --target=e2e --projects=web-admin,web --parallel=1"
```

- [ ] **Step 5: Run e2e verification when Docker and Playwright are available**

Run:

```bash
bunx nx run web:e2e
bun run test:e2e
```

Expected: PASS when local Docker-backed PostgreSQL/Redis and Playwright browsers are available. If unavailable, record the exact missing dependency and run the focused unit/type/build gates from Tasks 1-3.

- [ ] **Step 6: Commit public REST e2e**

Run:

```bash
git add apps/web/e2e apps/web/project.json package.json
git commit -m "test(web): add REST browser flow"
```

## Task 5: GRACE Contract Refresh

**Files:**

- Modify: `docs/requirements.xml`
- Modify: `docs/technology.xml`
- Modify: `docs/development-plan.xml`
- Modify: `docs/knowledge-graph.xml`
- Modify: `docs/verification-plan.xml`

- [ ] **Step 1: Update requirements with explicit public REST and admin GraphQL use cases**

In `docs/requirements.xml`:

- Update `UC-001` so local startup mentions `nx serve api`, `nx serve web-admin`, and `nx serve web`.
- Split the current GraphQL user flow into an admin GraphQL use case and a public REST use case.
- Add or update constraints so GraphQL is described as admin-only after this change.

- [ ] **Step 2: Update technology with Vite and admin-only GraphQL codegen**

In `docs/technology.xml`:

- Add Vite to the frontend framework list.
- Change `preferred-web-data-client` to mention REST for `web` and GraphQL request/codegen for `web-admin`.
- Change module-level checks from `web:codegen` to `web-admin:codegen`.

- [ ] **Step 3: Update development plan and knowledge graph with split modules**

Replace the single `M-WEB` GraphQL module with:

- `M-WEB-ADMIN`: Next.js GraphQL admin app.
- `M-WEB`: Vite REST public app.

Update `M-API` to shared HTTP API with REST and GraphQL adapters.

- [ ] **Step 4: Update verification plan**

Add verification entries for:

- admin GraphQL codegen/typecheck/e2e.
- public REST web tests/typecheck/build.
- API REST handler tests and GraphQL resolver compatibility.

- [ ] **Step 5: Validate GRACE XML**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected: PASS.

```bash
git add docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml
git commit -m "docs(grace): split web admin GraphQL and REST web contracts"
```

## Task 6: Integration Gates And Final Cleanup

**Files:**

- Modify: `README.md`
- Modify: `docs/infrastructure/ci-cd.md`
- Modify: `.gitlab-ci.yml`
- Modify: `docker/web.Dockerfile`
- Modify: `deploy/dokploy/docker-compose.template.yml`

- [ ] **Step 1: Search for stale ownership**

Run:

```bash
rg -n "apps/web/src/shared/api/generated|web:codegen|GraphQL.*web|web generated|nx run web:e2e|NEXT_PUBLIC_API_URL" .
```

Expected: remaining hits are either in `web-admin`, migration notes, or intentionally REST-neutral.

- [ ] **Step 2: Run workspace checks**

Run:

```bash
bun run codegen
bun run lint
bun run test
bun run build
```

Expected: PASS.

- [ ] **Step 3: Run broader gates when environment is available**

Run:

```bash
bun run test:coverage
bun run test:e2e
```

Expected: PASS, unless Docker or Playwright environment dependencies are unavailable. If unavailable, record the exact blocker and the focused checks that passed.

- [ ] **Step 4: Final diff review**

Run:

```bash
git diff --check
git status --short
```

Expected: no whitespace errors and only intended files changed.

- [ ] **Step 5: Commit cleanup**

Run:

```bash
git add README.md docs/infrastructure/ci-cd.md .gitlab-ci.yml docker/web.Dockerfile deploy/dokploy/docker-compose.template.yml package.json tools/coverage/coverage.config.json
git commit -m "chore(workspace): align web and web-admin gates"
```
