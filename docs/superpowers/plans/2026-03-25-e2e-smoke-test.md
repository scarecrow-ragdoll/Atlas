# E2E Smoke Test Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement a minimal end-to-end slice (PostgreSQL → Go repo → GraphQL → Next.js) proving the full stack works with user list, detail, and create.

**Architecture:** Backend: implement stubbed UserRepo (pgx), add bcrypt in UserService, wire resolvers. Frontend: add .graphql documents, codegen, Client Component page with React Query, Server Component page with fetch.

**Tech Stack:** Go 1.22 / pgx v5 / gqlgen / bcrypt | Next.js 15 / React Query v5 / graphql-request / Tailwind

**Spec:** `docs/superpowers/specs/2026-03-25-e2e-smoke-test-design.md`

---

## File Structure

| File | Action | Responsibility |
|------|--------|---------------|
| `apps/api/internal/repository/postgres/user_repo.go` | Modify | pgx queries: Create, GetByID, List (scan time.Time → string) |
| `apps/api/internal/service/user_service.go` | Modify | Add bcrypt hashing in Create |
| `apps/api/go.mod` / `apps/api/go.sum` | Modify | Promote `golang.org/x/crypto` to direct dependency |
| `apps/api/internal/graph/schema.resolvers.go` | Modify | Wire resolvers to service, type mapping |
| `apps/web/src/entities/user/api/users.graphql` | Create | GetUsers query document |
| `apps/web/src/entities/user/api/user.graphql` | Create | GetUser query document |
| `apps/web/src/entities/user/api/createUser.graphql` | Create | CreateUser mutation document |
| `apps/web/src/shared/api/generated/types.ts` | Regenerate | Run codegen after adding .graphql docs |
| `apps/web/app/users/page.tsx` | Create | Client Component — list + create form |
| `apps/web/app/users/[id]/page.tsx` | Create | Server Component — user detail |
| `apps/web/app/page.tsx` | Modify | Add /users navigation link |

**Note:** Migration `apps/api/internal/repository/postgres/migrations/00001_init.sql` already creates the `users` table with the required schema. Goose runs migrations at startup via Docker Compose.

---

### Task 1: Implement UserRepo — Create, GetByID, List

**Files:**
- Modify: `apps/api/internal/repository/postgres/user_repo.go`

- [ ] **Step 1: Implement GetByID**

```go
func (r *UserRepo) GetByID(ctx context.Context, id string) (*service.User, error) {
	var u service.User
	var createdAt, updatedAt time.Time
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Name, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	u.CreatedAt = createdAt.Format(time.RFC3339Nano)
	u.UpdatedAt = updatedAt.Format(time.RFC3339Nano)
	return &u, nil
}
```

Add imports: `"errors"`, `"time"`, and `"github.com/jackc/pgx/v5"`.

Design note: repo returns `nil, nil` on not-found (not a sentinel error). The resolver checks `u == nil` to handle this case. This keeps it simple for the smoke test.

- [ ] **Step 2: Implement Create**

```go
func (r *UserRepo) Create(ctx context.Context, input service.CreateUserInput) (*service.User, error) {
	var u service.User
	var createdAt, updatedAt time.Time
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (email, name, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id, email, name, created_at, updated_at`,
		input.Email, input.Name, input.Password,
	).Scan(&u.ID, &u.Email, &u.Name, &createdAt, &updatedAt)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, fmt.Errorf("duplicate email: %s", input.Email)
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	u.CreatedAt = createdAt.Format(time.RFC3339Nano)
	u.UpdatedAt = updatedAt.Format(time.RFC3339Nano)
	return &u, nil
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}
```

Add import: `"github.com/jackc/pgx/v5/pgconn"`.

Note: `input.Password` here is the already-hashed password from the service layer. The field is named `Password` in `service.CreateUserInput` but carries the hash at this point.

- [ ] **Step 3: Implement List with cursor pagination**

```go
func (r *UserRepo) List(ctx context.Context, first *int, after *string) ([]*service.User, int, error) {
	limit := 20
	if first != nil && *first > 0 {
		limit = *first
	}

	// Fetch limit+1 to detect hasNextPage without off-by-one
	args := []any{limit + 1}
	query := `SELECT id, email, name, created_at, updated_at FROM users`

	if after != nil && *after != "" {
		decoded, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid cursor: %w", err)
		}
		cursor, err := time.Parse(time.RFC3339Nano, string(decoded))
		if err != nil {
			return nil, 0, fmt.Errorf("invalid cursor time: %w", err)
		}
		query += ` WHERE created_at < $2`
		args = append(args, cursor)
	}

	query += ` ORDER BY created_at DESC LIMIT $1`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []*service.User
	for rows.Next() {
		var u service.User
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &createdAt, &updatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		u.CreatedAt = createdAt.Format(time.RFC3339Nano)
		u.UpdatedAt = updatedAt.Format(time.RFC3339Nano)
		users = append(users, &u)
	}

	// Trim to requested limit — extra row was only for hasNext detection
	if len(users) > limit {
		users = users[:limit]
	}

	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	return users, total, nil
}
```

Add imports: `"encoding/base64"` and `"time"` (already added in Step 1).

- [ ] **Step 4: Verify compilation**

Run: `cd apps/api && go build ./...`
Expected: Clean build, no errors.

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/repository/postgres/user_repo.go
git commit -m "feat(api): implement UserRepo Create, GetByID, List with pgx"
```

---

### Task 2: Add bcrypt hashing in UserService.Create

**Files:**
- Modify: `apps/api/internal/service/user_service.go`

- [ ] **Step 1: Add bcrypt hashing to Create method**

Replace the existing `Create` method:

```go
func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	input.Password = string(hashed)
	return s.repo.Create(ctx, input)
}
```

Add imports: `"fmt"` and `"golang.org/x/crypto/bcrypt"`.

- [ ] **Step 2: Promote bcrypt to direct dependency**

Run: `cd apps/api && go get golang.org/x/crypto/bcrypt && go mod tidy`

`golang.org/x/crypto` is already in go.mod as indirect — this promotes it to direct.

- [ ] **Step 3: Verify compilation**

Run: `cd apps/api && go build ./...`
Expected: Clean build.

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/service/user_service.go apps/api/go.mod apps/api/go.sum
git commit -m "feat(api): add bcrypt password hashing in UserService.Create"
```

---

### Task 3: Wire GraphQL resolvers to UserService

**Files:**
- Modify: `apps/api/internal/graph/schema.resolvers.go`

- [ ] **Step 1: Add helper functions at the bottom of the file (before the type declarations)**

```go
func mapUser(u *service.User) *model.User {
	return &model.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func mapCreateUserInput(in model.CreateUserInput) service.CreateUserInput {
	return service.CreateUserInput{
		Email:    in.Email,
		Name:     in.Name,
		Password: in.Password,
	}
}
```

Add import: `"monorepo-template/apps/api/internal/service"` and `"strings"`.

- [ ] **Step 2: Implement CreateUser resolver**

Replace the panic with:

```go
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (model.CreateUserResult, error) {
	u, err := r.UserService.Create(ctx, mapCreateUserInput(input))
	if err != nil {
		if strings.Contains(err.Error(), "duplicate email") {
			return model.ValidationError{Field: "email", Message: "already exists"}, nil
		}
		return nil, err
	}
	return model.CreateUserSuccess{User: mapUser(u)}, nil
}
```

- [ ] **Step 3: Implement User resolver**

Replace the panic with:

```go
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	u, err := r.UserService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	return mapUser(u), nil
}
```

- [ ] **Step 4: Implement Users resolver**

Replace the panic with:

```go
func (r *queryResolver) Users(ctx context.Context, pagination *model.PaginationInput) (*model.UserConnection, error) {
	var first *int
	var after *string
	if pagination != nil {
		first = pagination.First
		after = pagination.After
	}

	users, total, err := r.UserService.List(ctx, first, after)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.UserEdge, len(users))
	for i, u := range users {
		edges[i] = &model.UserEdge{
			Cursor: base64.StdEncoding.EncodeToString([]byte(u.CreatedAt)),
			Node:   mapUser(u),
		}
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}

	// Repo fetches limit+1 rows and trims — if we got exactly limit, there are more
	limit := 20
	if first != nil && *first > 0 {
		limit = *first
	}
	hasNext := len(users) == limit && total > limit

	return &model.UserConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNext,
			HasPreviousPage: after != nil,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
		TotalCount: total,
	}, nil
}
```

Add import: `"encoding/base64"`.

- [ ] **Step 5: Remove unused `fmt` import if present, verify compilation**

Run: `cd apps/api && go build ./...`
Expected: Clean build.

- [ ] **Step 6: Commit**

```bash
git add apps/api/internal/graph/schema.resolvers.go
git commit -m "feat(api): wire user, users, createUser resolvers to UserService"
```

---

### Task 4: Create frontend GraphQL documents and run codegen

**Files:**
- Create: `apps/web/src/entities/user/api/users.graphql`
- Create: `apps/web/src/entities/user/api/user.graphql`
- Create: `apps/web/src/entities/user/api/createUser.graphql`
- Regenerate: `apps/web/src/shared/api/generated/types.ts`

- [ ] **Step 1: Create directory structure**

```bash
mkdir -p apps/web/src/entities/user/api
```

- [ ] **Step 2: Create `users.graphql`**

Write to `apps/web/src/entities/user/api/users.graphql`:

```graphql
query GetUsers($first: Int, $after: String) {
  users(pagination: { first: $first, after: $after }) {
    edges {
      cursor
      node {
        id
        email
        name
        createdAt
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
  }
}
```

- [ ] **Step 3: Create `user.graphql`**

Write to `apps/web/src/entities/user/api/user.graphql`:

```graphql
query GetUser($id: UUID!) {
  user(id: $id) {
    id
    email
    name
    createdAt
    updatedAt
  }
}
```

- [ ] **Step 4: Create `createUser.graphql`**

Write to `apps/web/src/entities/user/api/createUser.graphql`:

```graphql
mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    ... on CreateUserSuccess {
      user {
        id
        email
        name
      }
    }
    ... on ValidationError {
      field
      message
    }
    ... on AuthError {
      message
    }
  }
}
```

- [ ] **Step 5: Run codegen**

```bash
cd tools/codegen && bun graphql-codegen --config codegen.ts
```

Expected: `apps/web/src/shared/api/generated/types.ts` is regenerated with new operation types (`GetUsersQuery`, `GetUserQuery`, `CreateUserMutation`, etc.).

- [ ] **Step 6: Verify generated types contain operation types**

Check that `types.ts` now contains `GetUsersQuery`, `GetUserQuery`, `CreateUserMutation`.

- [ ] **Step 7: Commit**

```bash
git add apps/web/src/entities/user/api/ apps/web/src/shared/api/generated/types.ts
git commit -m "feat(web): add GraphQL documents and regenerate types"
```

---

### Task 5: Create Client Component page — `/users`

**Files:**
- Create: `apps/web/app/users/page.tsx`

- [ ] **Step 1: Create the page file**

Write to `apps/web/app/users/page.tsx`:

```tsx
'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { graphqlClient } from '@shared/api/graphql-client';
import Link from 'next/link';
import { useState } from 'react';

const GET_USERS = /* GraphQL */ `
  query GetUsers($first: Int, $after: String) {
    users(pagination: { first: $first, after: $after }) {
      edges {
        cursor
        node { id email name createdAt }
      }
      pageInfo { hasNextPage endCursor }
      totalCount
    }
  }
`;

const CREATE_USER = /* GraphQL */ `
  mutation CreateUser($input: CreateUserInput!) {
    createUser(input: $input) {
      ... on CreateUserSuccess { user { id email name } }
      ... on ValidationError { field message }
      ... on AuthError { message }
    }
  }
`;

interface UserNode {
  id: string;
  email: string;
  name: string;
  createdAt: string;
}

interface GetUsersResponse {
  users: {
    edges: Array<{ cursor: string; node: UserNode }>;
    pageInfo: { hasNextPage: boolean; endCursor: string | null };
    totalCount: number;
  };
}

interface CreateUserResponse {
  createUser:
    | { __typename: 'CreateUserSuccess'; user: { id: string; email: string; name: string } }
    | { __typename: 'ValidationError'; field: string; message: string }
    | { __typename: 'AuthError'; message: string };
}

export default function UsersPage() {
  const queryClient = useQueryClient();
  const [form, setForm] = useState({ email: '', name: '', password: '' });
  const [error, setError] = useState<string | null>(null);

  const { data, isLoading, isError } = useQuery({
    queryKey: ['users'],
    queryFn: () => graphqlClient.request<GetUsersResponse>(GET_USERS, { first: 20 }),
  });

  const mutation = useMutation({
    mutationFn: (input: { email: string; name: string; password: string }) =>
      graphqlClient.request<CreateUserResponse>(CREATE_USER, { input }),
    onSuccess: (res) => {
      const result = res.createUser;
      if ('user' in result) {
        queryClient.invalidateQueries({ queryKey: ['users'] });
        setForm({ email: '', name: '', password: '' });
        setError(null);
      } else if ('field' in result) {
        setError(`${result.field}: ${result.message}`);
      } else {
        setError(result.message);
      }
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    mutation.mutate(form);
  };

  return (
    <main className="mx-auto max-w-2xl p-8">
      <div className="mb-8 flex items-center justify-between">
        <h1 className="text-2xl font-bold">Users</h1>
        <Link href="/" className="text-blue-600 hover:underline">Home</Link>
      </div>

      <form onSubmit={handleSubmit} className="mb-8 space-y-3 rounded border p-4">
        <h2 className="text-lg font-semibold">Create User</h2>
        <input
          type="text"
          placeholder="Name"
          value={form.name}
          onChange={(e) => setForm({ ...form, name: e.target.value })}
          className="w-full rounded border px-3 py-2"
          required
        />
        <input
          type="email"
          placeholder="Email"
          value={form.email}
          onChange={(e) => setForm({ ...form, email: e.target.value })}
          className="w-full rounded border px-3 py-2"
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={form.password}
          onChange={(e) => setForm({ ...form, password: e.target.value })}
          className="w-full rounded border px-3 py-2"
          required
        />
        {error && <p className="text-sm text-red-600">{error}</p>}
        <button
          type="submit"
          disabled={mutation.isPending}
          className="rounded bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 disabled:opacity-50"
        >
          {mutation.isPending ? 'Creating...' : 'Create'}
        </button>
      </form>

      {isLoading && <p>Loading...</p>}
      {isError && <p className="text-red-600">Failed to load users.</p>}

      {data && (
        <>
          <p className="mb-4 text-sm text-gray-500">Total: {data.users.totalCount}</p>
          <ul className="space-y-2">
            {data.users.edges.map(({ node }) => (
              <li key={node.id} className="rounded border p-3">
                <Link href={`/users/${node.id}`} className="font-medium text-blue-600 hover:underline">
                  {node.name}
                </Link>
                <span className="ml-2 text-sm text-gray-500">{node.email}</span>
              </li>
            ))}
            {data.users.edges.length === 0 && (
              <li className="text-gray-500">No users yet. Create one above.</li>
            )}
          </ul>
        </>
      )}
    </main>
  );
}
```

- [ ] **Step 2: Verify TypeScript compilation**

Run: `cd apps/web && npx tsc --noEmit`
Expected: No type errors.

- [ ] **Step 3: Commit**

```bash
git add apps/web/app/users/page.tsx
git commit -m "feat(web): add /users client component page with list and create form"
```

---

### Task 6: Create Server Component page — `/users/[id]`

**Files:**
- Create: `apps/web/app/users/[id]/page.tsx`

- [ ] **Step 1: Create the page file**

Write to `apps/web/app/users/[id]/page.tsx`:

```tsx
import Link from 'next/link';
import { appConfig } from '@shared/config';

const GET_USER_QUERY = /* GraphQL */ `
  query GetUser($id: UUID!) {
    user(id: $id) { id email name createdAt updatedAt }
  }
`;

interface User {
  id: string;
  email: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

interface GetUserResponse {
  data: { user: User | null };
}

async function getUser(id: string): Promise<User | null> {
  const res = await fetch(appConfig.apiUrl, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query: GET_USER_QUERY, variables: { id } }),
    cache: 'no-store',
  });

  if (!res.ok) return null;

  const json: GetUserResponse = await res.json();
  return json.data.user;
}

export default async function UserDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  const user = await getUser(id);

  if (!user) {
    return (
      <main className="mx-auto max-w-2xl p-8">
        <h1 className="text-2xl font-bold">User not found</h1>
        <Link href="/users" className="mt-4 inline-block text-blue-600 hover:underline">
          Back to users
        </Link>
      </main>
    );
  }

  return (
    <main className="mx-auto max-w-2xl p-8">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-bold">{user.name}</h1>
        <Link href="/users" className="text-blue-600 hover:underline">Back to users</Link>
      </div>

      <dl className="space-y-3">
        <div>
          <dt className="text-sm text-gray-500">Email</dt>
          <dd>{user.email}</dd>
        </div>
        <div>
          <dt className="text-sm text-gray-500">Created</dt>
          <dd>{new Date(user.createdAt).toLocaleString()}</dd>
        </div>
        <div>
          <dt className="text-sm text-gray-500">Updated</dt>
          <dd>{new Date(user.updatedAt).toLocaleString()}</dd>
        </div>
        <div>
          <dt className="text-sm text-gray-500">ID</dt>
          <dd className="font-mono text-sm">{user.id}</dd>
        </div>
      </dl>
    </main>
  );
}
```

Note: `params` is `Promise<{ id: string }>` in Next.js 15 App Router (async params).

- [ ] **Step 2: Verify TypeScript compilation**

Run: `cd apps/web && npx tsc --noEmit`
Expected: No type errors.

- [ ] **Step 3: Commit**

```bash
git add apps/web/app/users/\[id\]/page.tsx
git commit -m "feat(web): add /users/[id] server component page"
```

---

### Task 7: Add navigation link on home page

**Files:**
- Modify: `apps/web/app/page.tsx`

- [ ] **Step 1: Add link to /users**

Replace the content of `apps/web/app/page.tsx`:

```tsx
import Link from 'next/link';

export default function HomePage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center">
      <h1 className="text-4xl font-bold">Monorepo Template</h1>
      <p className="mt-4 text-lg text-gray-600">Go + Next.js + GraphQL</p>
      <Link href="/users" className="mt-6 text-blue-600 hover:underline">
        Users (smoke test)
      </Link>
    </main>
  );
}
```

- [ ] **Step 2: Verify TypeScript compilation**

Run: `cd apps/web && npx tsc --noEmit`
Expected: No type errors.

- [ ] **Step 3: Commit**

```bash
git add apps/web/app/page.tsx
git commit -m "feat(web): add /users navigation link on home page"
```

---

## Verification Checklist

After all tasks are complete, verify with Docker Compose:

1. `docker compose -f docker/docker-compose.yml up --build` — all services healthy
2. `http://localhost:8080/healthz` — returns `{"status":"ok"}`
3. `http://localhost:8080/playground` — GraphQL Playground opens
4. Run in playground: `mutation { createUser(input: {email: "test@test.com", name: "Test User", password: "password123"}) { ... on CreateUserSuccess { user { id email name } } ... on ValidationError { field message } } }`
5. Run in playground: `{ users { edges { node { id email name } } totalCount } }`
6. `http://localhost:3000` — home page with "Users (smoke test)" link
7. `http://localhost:3000/users` — user list page (client-rendered), create form works
8. Click user name — `/users/[id]` detail page (server-rendered)
