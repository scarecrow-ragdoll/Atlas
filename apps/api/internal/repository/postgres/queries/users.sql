-- name: GetUserByID :one
SELECT id, email, name, created_at, updated_at
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, email, name, created_at, updated_at
FROM users
WHERE (sqlc.narg('after_created_at')::timestamptz IS NULL OR created_at < sqlc.narg('after_created_at')::timestamptz)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit_rows');

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CreateUser :one
INSERT INTO users (email, name, password_hash)
VALUES ($1, $2, $3)
RETURNING id, email, name, created_at, updated_at;

-- name: UpdateUser :one
UPDATE users
SET name = COALESCE(sqlc.narg('name'), name),
    email = COALESCE(sqlc.narg('email'), email),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING id, email, name, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- FILE: apps/api/internal/repository/postgres/queries/users.sql
-- VERSION: 1.0.1
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc users queries used by the PostgreSQL UserRepo adapter.
--   SCOPE: CRUD and pagination queries for the users table; excludes schema ownership and transport mapping.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00001_init.sql.
--   LINKS: M-API / V-M-API.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   GetUserByID - Fetches one user row by UUID.
--   ListUsers - Fetches a created_at-desc page with optional cursor cutoff.
--   CountUsers - Counts all user rows for connection metadata.
--   CreateUser - Inserts one user and returns the persisted public row.
--   UpdateUser - Applies nullable name/email updates and returns the persisted public row.
--   DeleteUser - Deletes one user idempotently by UUID.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.1 - Moved file-local GRACE markup after queries so sqlc does not copy it into generated output.
-- END_CHANGE_SUMMARY
