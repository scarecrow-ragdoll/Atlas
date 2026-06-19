# E2E Smoke Test: Users CRUD Slice

## Goal

Implement minimal end-to-end flow (PostgreSQL -> Go repository -> GraphQL resolvers -> Next.js frontend) to verify the full stack works. Covers: list users, view single user, create user.

## Scope

### In scope
- `UserRepo`: `GetByID`, `List`, `Create` — real pgx queries against PostgreSQL
- GraphQL resolvers: `user`, `users`, `createUser` — call UserService instead of panic
- Frontend GraphQL documents (`.graphql` files) for codegen
- Client Component page `/users` — React Query list + create form
- Server Component page `/users/[id]` — server-side fetch of single user

### Out of scope
- `UpdateUser`, `DeleteUser` resolvers (remain stubbed)
- Auth middleware validation (remains placeholder)
- Redis integration
- Tests (separate task)
- Styling beyond Tailwind basics

## Backend Changes

### 1. Repository: `apps/api/internal/repository/postgres/user_repo.go`

Implement three methods with pgx:

**`Create(ctx, input)`**
- INSERT INTO users (email, name, password_hash) VALUES ($1, $2, $3) RETURNING id, email, name, created_at, updated_at
- Receives already-hashed password (hashing done in service layer)
- On duplicate email: return error (let resolver map to ValidationError)

**`GetByID(ctx, id)`**
- SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1
- pgx.ErrNoRows -> return nil, ErrNotFound

**`List(ctx, first, after)`**
- SELECT with ORDER BY created_at DESC, LIMIT (first or default 20)
- Cursor encoding: `base64(created_at.Format(RFC3339Nano))` — opaque string for the client
- Cursor decoding: `base64.Decode(after)` -> parse as RFC3339Nano -> `WHERE created_at < $decoded` predicate
- Second query: SELECT COUNT(*) FROM users for totalCount
- Return ([]*User, totalCount, error)

### 2. Resolvers: `apps/api/internal/graph/schema.resolvers.go`

**`CreateUser`**
- Map `model.CreateUserInput` -> `service.CreateUserInput` (field-by-field, same names)
- Call `r.UserService.Create(ctx, svcInput)`
- On success: return `model.CreateUserSuccess{User: mapUser(u)}`
- On validation error (duplicate email): return `model.ValidationError{Field: "email", Message: "already exists"}`

**`User`**
- Call `r.UserService.GetByID(ctx, id)`
- On not found: return nil, nil (GraphQL nullable)
- On success: return mapUser(u)

**`Users`**
- Extract `first`/`after` from `*model.PaginationInput` (nil-safe, default first=20)
- Call `r.UserService.List(ctx, first, after)`
- Build UserConnection with edges, cursors, pageInfo

### 3. Helper: service.User -> model.User mapper

Simple field-by-field mapping function in resolvers file. Timestamps formatted as RFC3339.

### 4. Service layer: password hashing

`UserService.Create` hashes the plain-text password with bcrypt before calling `repo.Create`. The repository receives only `password_hash` — no crypto dependency in the persistence layer.

### 5. Dependencies

- `golang.org/x/crypto/bcrypt` — used in `UserService.Create`

## Frontend Changes

### 6. GraphQL Documents: `apps/web/src/entities/user/api/`

**`users.graphql`**
```graphql
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
```

**`user.graphql`**
```graphql
query GetUser($id: UUID!) {
  user(id: $id) { id email name createdAt updatedAt }
}
```

**`createUser.graphql`**
```graphql
mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    ... on CreateUserSuccess { user { id email name } }
    ... on ValidationError { field message }
    ... on AuthError { message }
  }
}
```

After adding documents, run codegen: `cd tools/codegen && bun graphql-codegen` to regenerate `types.ts` with operation types.

Note: the current codegen config generates base types + operation types only (no typed-document-node plugin). Frontend code must use raw query strings with `graphql-request`, not typed document nodes.

### 7. Client Component Page: `apps/web/app/users/page.tsx`

- `'use client'` directive
- `useQuery` for GetUsers with graphql-request
- `useMutation` for CreateUser, invalidates users query on success
- Simple form: email + name + password inputs, submit button
- User list rendered as cards/rows with links to `/users/[id]`
- Loading/error states

### 8. Server Component Page: `apps/web/app/users/[id]/page.tsx`

- Server Component (no `'use client'`)
- `fetch()` POST to GraphQL endpoint with GetUser query, `{ cache: 'no-store' }` for always-fresh data
- Render user details (name, email, dates)
- Back link to `/users`

### 9. Navigation

Add link to `/users` on the home page (`app/page.tsx`).

## Data Flow

```
Browser -> /users (Client Component)
  -> React Query -> graphql-request -> POST /graphql
    -> Chi router -> gqlgen -> UsersResolver
      -> UserService.List -> UserRepo.List -> PostgreSQL
    <- UserConnection JSON
  <- React renders list

Browser -> /users/[id] (Server Component)
  -> Next.js server fetch -> POST /graphql
    -> Chi router -> gqlgen -> UserResolver
      -> UserService.GetByID -> UserRepo.GetByID -> PostgreSQL
    <- User JSON
  <- Server-rendered HTML

Browser -> Create form submit
  -> useMutation -> graphql-request -> POST /graphql
    -> Chi router -> gqlgen -> CreateUserResolver
      -> UserService.Create -> UserRepo.Create -> PostgreSQL INSERT
    <- CreateUserSuccess | ValidationError
  <- React Query invalidates -> re-fetches list
```

## File Manifest

| # | File | Action | Description |
|---|------|--------|-------------|
| 1 | `apps/api/internal/repository/postgres/user_repo.go` | Modify | Implement GetByID, List, Create with pgx |
| 2 | `apps/api/internal/service/user_service.go` | Modify | Add bcrypt hashing in Create before calling repo |
| 3 | `apps/api/internal/graph/schema.resolvers.go` | Modify | Implement user, users, createUser resolvers with model->service mapping |
| 4 | `apps/api/go.mod` / `go.sum` | Modify | Add `golang.org/x/crypto` dependency |
| 5 | `apps/web/src/entities/user/api/users.graphql` | Create | GetUsers query document |
| 6 | `apps/web/src/entities/user/api/user.graphql` | Create | GetUser query document |
| 7 | `apps/web/src/entities/user/api/createUser.graphql` | Create | CreateUser mutation document (all 3 union members) |
| 8 | `apps/web/src/shared/api/generated/types.ts` | Regenerate | Run `cd tools/codegen && bun graphql-codegen` |
| 9 | `apps/web/app/users/page.tsx` | Create | Client Component — user list + create form |
| 10 | `apps/web/app/users/[id]/page.tsx` | Create | Server Component — user detail (cache: no-store) |
| 11 | `apps/web/app/page.tsx` | Modify | Add navigation link to /users |

## Verification

After implementation, verify with Docker Compose:
1. `docker compose up` — all services start, healthz OK
2. Open GraphQL Playground at `http://localhost:8080/` — run `users` query
3. Run `createUser` mutation in playground — verify user created
4. Open `http://localhost:3000/users` — see list (client-rendered)
5. Create user via form — appears in list
6. Click user — `/users/[id]` shows detail (server-rendered)
