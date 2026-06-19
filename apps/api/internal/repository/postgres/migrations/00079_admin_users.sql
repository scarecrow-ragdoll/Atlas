-- +goose Up
CREATE TABLE admin_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(32) NOT NULL DEFAULT 'ADMIN',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_admin_users_email_lower ON admin_users (LOWER(email));
CREATE INDEX idx_admin_users_active ON admin_users (is_active);

-- +goose Down
DROP TABLE IF EXISTS admin_users;

-- FILE: apps/api/internal/repository/postgres/migrations/00079_admin_users.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Add the admin_users table for web-admin authentication.
--   SCOPE: Admin identity schema, lower-case unique email enforcement, active flag, role, and timestamps; excludes session storage and public users schema.
--   DEPENDS: apps/api/internal/repository/postgres/migrations/00001_init.sql.
--   LINKS: M-API / V-M-API.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   admin_users - Stores web-admin identities separate from public reference users.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.2 - Numbered after historical local goose version 78 so existing databases apply the admin table.
-- END_CHANGE_SUMMARY
