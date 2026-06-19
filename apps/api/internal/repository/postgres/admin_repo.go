// FILE: apps/api/internal/repository/postgres/admin_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Adapt sqlc-generated PostgreSQL admin_users queries to the service.AdminRepository contract.
//   SCOPE: Admin count, create, email/id lookup, row mapping, UUID parsing, and duplicate email mapping; excludes password hashing and session storage.
//   DEPENDS: apps/api/internal/repository/postgres/generated, github.com/jackc/pgx/v5, apps/api/internal/service.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   NewAdminRepo - Constructs the production admin repository from a pgx pool.
//   AdminRepo.Count - Counts admin identities for seed bootstrap.
//   AdminRepo.Create - Inserts one normalized admin and maps duplicate email.
//   AdminRepo.GetByEmail - Reads one admin by normalized email.
//   AdminRepo.GetByID - Reads one admin by UUID.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added PostgreSQL admin repository adapter.
// END_CHANGE_SUMMARY

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/repository/postgres/generated"
	"monorepo-template/apps/api/internal/service"
)

type AdminRepo struct {
	queries generated.Querier
}

func NewAdminRepo(pool *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{queries: generated.New(pool)}
}

// START_CONTRACT: Count
//
//	PURPOSE: Count admin identities to decide first-admin bootstrap behavior.
//	INPUTS: { ctx: context.Context - request or startup context }
//	OUTPUTS: { int - admin row count, error - count failure }
//	SIDE_EFFECTS: Reads PostgreSQL through generated queries.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: Count
func (r *AdminRepo) Count(ctx context.Context) (int, error) {
	count, err := r.queries.CountAdminUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("AdminRepo.Count: %w", err)
	}
	return int(count), nil
}

// START_CONTRACT: Create
//
//	PURPOSE: Insert one normalized active admin identity and map duplicate email conflicts.
//	INPUTS: { ctx: context.Context - request or startup context, input: service.CreateAdminInput - validated admin fields with password hash }
//	OUTPUTS: { *service.Admin - persisted admin, error - duplicate email or insert failure }
//	SIDE_EFFECTS: Inserts PostgreSQL admin_users row through generated queries.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: Create
func (r *AdminRepo) Create(ctx context.Context, input service.CreateAdminInput) (*service.Admin, error) {
	row, err := r.queries.CreateAdminUser(ctx, generated.CreateAdminUserParams{
		Email:        input.Email,
		Name:         input.Name,
		PasswordHash: input.PasswordHash,
		Role:         input.Role,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, service.ErrAdminDuplicateEmail
		}
		return nil, fmt.Errorf("AdminRepo.Create: %w", err)
	}
	return adminFromFields(row.ID, row.Email, row.Name, row.PasswordHash, row.Role, row.IsActive, row.CreatedAt, row.UpdatedAt), nil
}

// START_CONTRACT: GetByEmail
//
//	PURPOSE: Read one admin by case-insensitive email and return nil when missing.
//	INPUTS: { ctx: context.Context - request context, email: string - admin email }
//	OUTPUTS: { *service.Admin - mapped admin or nil, error - query failure }
//	SIDE_EFFECTS: Reads PostgreSQL through generated queries.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: GetByEmail
func (r *AdminRepo) GetByEmail(ctx context.Context, email string) (*service.Admin, error) {
	row, err := r.queries.GetAdminUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("AdminRepo.GetByEmail: %w", err)
	}
	return adminFromFields(row.ID, row.Email, row.Name, row.PasswordHash, row.Role, row.IsActive, row.CreatedAt, row.UpdatedAt), nil
}

// START_CONTRACT: GetByID
//
//	PURPOSE: Read one admin by UUID and return nil when missing.
//	INPUTS: { ctx: context.Context - request context, id: string - admin UUID }
//	OUTPUTS: { *service.Admin - mapped admin or nil, error - UUID parsing or query failure }
//	SIDE_EFFECTS: Reads PostgreSQL through generated queries.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: GetByID
func (r *AdminRepo) GetByID(ctx context.Context, id string) (*service.Admin, error) {
	adminID, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.GetByID: invalid admin id: %w", err)
	}
	row, err := r.queries.GetAdminUserByID(ctx, adminID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("AdminRepo.GetByID: %w", err)
	}
	return adminFromFields(row.ID, row.Email, row.Name, row.PasswordHash, row.Role, row.IsActive, row.CreatedAt, row.UpdatedAt), nil
}

func adminFromFields(id pgtype.UUID, email string, name string, passwordHash string, role string, isActive bool, createdAt pgtype.Timestamptz, updatedAt pgtype.Timestamptz) *service.Admin {
	return &service.Admin{
		ID:           id.String(),
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
		Role:         role,
		IsActive:     isActive,
		CreatedAt:    formatTimestamp(createdAt),
		UpdatedAt:    formatTimestamp(updatedAt),
	}
}
