// FILE: apps/api/internal/repository/postgres/user_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Adapt sqlc-generated PostgreSQL users queries to the service.UserRepository contract.
//   SCOPE: Users persistence CRUD, keyset pagination, generated-row mapping, and PostgreSQL error mapping; excludes transport handlers and service validation.
//   DEPENDS: apps/api/internal/repository/postgres/generated, github.com/jackc/pgx/v5, apps/api/internal/service, libs/go/logger.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   NewUserRepo - Constructs the production users repository from a pgx pool.
//   userQueries - Narrows generated sqlc methods needed by UserRepo so tests are insulated from unrelated query groups.
//   UserRepo.GetByID - Reads one user by UUID and maps missing rows to nil.
//   UserRepo.List - Reads a cursor page and total count.
//   UserRepo.Create - Inserts one user and maps duplicate email errors.
//   UserRepo.Update - Applies nullable name/email updates and maps missing rows to nil.
//   UserRepo.Delete - Deletes one user idempotently.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Narrowed UserRepo to user query methods after admin sqlc generation widened generated.Querier.
// END_CHANGE_SUMMARY

package postgres

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/repository/postgres/generated"
	"monorepo-template/apps/api/internal/service"
	"monorepo-template/libs/go/logger"
)

type UserRepo struct {
	queries userQueries
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return newUserRepoWithQueries(generated.New(pool))
}

type userQueries interface {
	GetUserByID(ctx context.Context, id pgtype.UUID) (generated.GetUserByIDRow, error)
	ListUsers(ctx context.Context, arg generated.ListUsersParams) ([]generated.ListUsersRow, error)
	CountUsers(ctx context.Context) (int64, error)
	CreateUser(ctx context.Context, arg generated.CreateUserParams) (generated.CreateUserRow, error)
	UpdateUser(ctx context.Context, arg generated.UpdateUserParams) (generated.UpdateUserRow, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) error
}

func newUserRepoWithQueries(queries userQueries) *UserRepo {
	return &UserRepo{queries: queries}
}

// START_CONTRACT: GetByID
//
//	PURPOSE: Retrieve a user by UUID and return nil when the row does not exist.
//	INPUTS: { ctx: context.Context - request context, id: string - user UUID }
//	OUTPUTS: { *service.User - mapped user or nil, error - query or UUID parsing failure }
//	SIDE_EFFECTS: Reads PostgreSQL through generated queries.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: GetByID
func (r *UserRepo) GetByID(ctx context.Context, id string) (*service.User, error) {
	const op = "UserRepo.GetByID"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("querying user by id", zap.String("user_id", id))

	userID, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid user id: %w", op, err)
	}

	row, err := r.queries.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return userFromFields(row.ID, row.Email, row.Name, row.CreatedAt, row.UpdatedAt), nil
}

// START_CONTRACT: List
//
//	PURPOSE: Return a created_at-desc user page plus total row count.
//	INPUTS: { ctx: context.Context - request context, first: *int - page size, after: *string - base64 RFC3339Nano cursor }
//	OUTPUTS: { []*service.User - page rows, int - total rows, error - cursor, list, scan, or count failure }
//	SIDE_EFFECTS: Reads PostgreSQL through generated queries.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: List
func (r *UserRepo) List(ctx context.Context, first *int, after *string) ([]*service.User, int, error) {
	const op = "UserRepo.List"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("querying users list")

	limit := int32(20)
	if first != nil && *first > 0 {
		if *first >= math.MaxInt32 {
			return nil, 0, fmt.Errorf("%s: first exceeds max supported page size", op)
		}
		limit = int32(*first)
	}

	cursor, err := cursorFromString(after)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.queries.ListUsers(ctx, generated.ListUsersParams{
		AfterCreatedAt: cursor,
		LimitRows:      limit + 1,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	users := make([]*service.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, userFromFields(row.ID, row.Email, row.Name, row.CreatedAt, row.UpdatedAt))
	}
	if len(users) > int(limit) {
		users = users[:int(limit)]
	}

	total, err := r.queries.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count: %w", op, err)
	}

	return users, int(total), nil
}

// START_CONTRACT: Create
//
//	PURPOSE: Insert one user row and map duplicate email conflicts.
//	INPUTS: { ctx: context.Context - request context, input: service.CreateUserInput - validated service input with password hash }
//	OUTPUTS: { *service.User - persisted user, error - duplicate email or insert failure }
//	SIDE_EFFECTS: Inserts PostgreSQL row through generated queries.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: Create
func (r *UserRepo) Create(ctx context.Context, input service.CreateUserInput) (*service.User, error) {
	const op = "UserRepo.Create"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("inserting user", zap.String("email", input.Email))

	row, err := r.queries.CreateUser(ctx, generated.CreateUserParams{
		Email:        input.Email,
		Name:         input.Name,
		PasswordHash: input.Password,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, fmt.Errorf("%s: duplicate email: %s", op, input.Email)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return userFromFields(row.ID, row.Email, row.Name, row.CreatedAt, row.UpdatedAt), nil
}

// START_CONTRACT: Update
//
//	PURPOSE: Apply optional name and email changes and return nil when the user does not exist.
//	INPUTS: { ctx: context.Context - request context, id: string - user UUID, input: service.UpdateUserInput - optional fields }
//	OUTPUTS: { *service.User - updated user or nil, error - duplicate email, UUID parsing, or update failure }
//	SIDE_EFFECTS: Updates PostgreSQL row through generated queries.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: Update
func (r *UserRepo) Update(ctx context.Context, id string, input service.UpdateUserInput) (*service.User, error) {
	const op = "UserRepo.Update"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("updating user", zap.String("user_id", id))

	userID, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid user id: %w", op, err)
	}

	row, err := r.queries.UpdateUser(ctx, generated.UpdateUserParams{
		ID:    userID,
		Name:  nullableText(input.Name),
		Email: nullableText(input.Email),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		if isDuplicateKeyError(err) {
			return nil, fmt.Errorf("%s: duplicate email: %s", op, derefString(input.Email))
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return userFromFields(row.ID, row.Email, row.Name, row.CreatedAt, row.UpdatedAt), nil
}

// START_CONTRACT: Delete
//
//	PURPOSE: Delete one user idempotently by UUID.
//	INPUTS: { ctx: context.Context - request context, id: string - user UUID }
//	OUTPUTS: { error - UUID parsing or delete failure }
//	SIDE_EFFECTS: Deletes PostgreSQL row through generated queries.
//	LINKS: M-API / V-M-API.
//
// END_CONTRACT: Delete
func (r *UserRepo) Delete(ctx context.Context, id string) error {
	const op = "UserRepo.Delete"
	log := logger.FromContext(ctx).With(zap.String("op", op))
	log.Debug("deleting user", zap.String("user_id", id))

	userID, err := uuidFromString(id)
	if err != nil {
		return fmt.Errorf("%s: invalid user id: %w", op, err)
	}
	if err := r.queries.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func uuidFromString(value string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(value); err != nil {
		return pgtype.UUID{}, err
	}
	return uuid, nil
}

func cursorFromString(after *string) (pgtype.Timestamptz, error) {
	if after == nil || *after == "" {
		return pgtype.Timestamptz{}, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(*after)
	if err != nil {
		return pgtype.Timestamptz{}, fmt.Errorf("invalid cursor: %w", err)
	}
	cursor, err := time.Parse(time.RFC3339Nano, string(decoded))
	if err != nil {
		return pgtype.Timestamptz{}, fmt.Errorf("invalid cursor time: %w", err)
	}
	return pgtype.Timestamptz{Time: cursor, Valid: true}, nil
}

func nullableText(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *value, Valid: true}
}

func userFromFields(id pgtype.UUID, email string, name string, createdAt pgtype.Timestamptz, updatedAt pgtype.Timestamptz) *service.User {
	return &service.User{
		ID:        id.String(),
		Email:     email,
		Name:      name,
		CreatedAt: formatTimestamp(createdAt),
		UpdatedAt: formatTimestamp(updatedAt),
	}
}

func formatTimestamp(value pgtype.Timestamptz) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format(time.RFC3339Nano)
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
