-- name: CountAdminUsers :one
SELECT COUNT(*) FROM admin_users;

-- name: CreateAdminUser :one
INSERT INTO admin_users (email, name, password_hash, role, is_active)
VALUES (LOWER(sqlc.arg(email)), sqlc.arg(name), sqlc.arg(password_hash), sqlc.arg(role), TRUE)
RETURNING id, email, name, password_hash, role, is_active, created_at, updated_at;

-- name: GetAdminUserByEmail :one
SELECT id, email, name, password_hash, role, is_active, created_at, updated_at
FROM admin_users
WHERE LOWER(email) = LOWER($1);

-- name: GetAdminUserByID :one
SELECT id, email, name, password_hash, role, is_active, created_at, updated_at
FROM admin_users
WHERE id = $1;

-- FILE: apps/api/internal/repository/postgres/queries/admin_users.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc admin user queries used by the PostgreSQL AdminRepo adapter.
--   SCOPE: Count, create, and identity lookup for admin_users; excludes public users persistence and Redis sessions.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00079_admin_users.sql.
--   LINKS: M-API / V-M-API.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   CountAdminUsers - Counts admin identities for bootstrap-only seed.
--   CreateAdminUser - Inserts one normalized active admin.
--   GetAdminUserByEmail - Fetches one admin by case-insensitive email.
--   GetAdminUserByID - Fetches one admin by UUID.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.0 - Added admin_users sqlc query contract.
-- END_CHANGE_SUMMARY
