// FILE: apps/api/internal/graph/schema_resolvers_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify GraphQL user resolver behavior and admin-principal guard coverage.
//   SCOPE: User CRUD/list resolver mapping, auth guard denial, validation/not-found/error mapping, and authenticated user-domain behavior; excludes admin auth resolver operations.
//   DEPENDS: apps/api/internal/graph, apps/api/internal/middleware, apps/api/internal/service.
//   LINKS: M-API / V-M-API / M-GRAPHQL-SCHEMA / V-M-GRAPHQL-SCHEMA.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   adminCtx - Provides authenticated admin principal context for user-domain resolver tests.
//   TestProtectedUserResolvers_* - Verifies unauthenticated user GraphQL operations are denied.
//   Test*User* - Verifies authenticated user GraphQL domain mapping.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin-principal guard coverage for user GraphQL resolvers.
// END_CHANGE_SUMMARY

package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/graph/model"
	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/apps/api/internal/service"
)

type resolverUserRepo struct {
	users     map[string]*service.User
	err       error
	createErr error
	updateErr error
	deleteErr error
}

func newTestResolver(repo *resolverUserRepo) *Resolver {
	return &Resolver{UserService: service.NewUserService(repo)}
}

func newResolverRepo() *resolverUserRepo {
	return &resolverUserRepo{users: map[string]*service.User{
		"user-1": {
			ID:        "user-1",
			Email:     "one@example.com",
			Name:      "One User",
			CreatedAt: "2026-05-02T00:00:00Z",
			UpdatedAt: "2026-05-02T00:00:00Z",
		},
	}}
}

func (r *resolverUserRepo) GetByID(ctx context.Context, id string) (*service.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.users[id], nil
}

func (r *resolverUserRepo) List(ctx context.Context, first *int, after *string) ([]*service.User, int, error) {
	if r.err != nil {
		return nil, 0, r.err
	}
	users := make([]*service.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, len(users), nil
}

func (r *resolverUserRepo) Create(ctx context.Context, input service.CreateUserInput) (*service.User, error) {
	if r.createErr != nil {
		return nil, r.createErr
	}
	user := &service.User{
		ID:        "created-id",
		Email:     input.Email,
		Name:      input.Name,
		CreatedAt: "2026-05-02T00:00:00Z",
		UpdatedAt: "2026-05-02T00:00:00Z",
	}
	r.users[user.ID] = user
	return user, nil
}

func (r *resolverUserRepo) Update(ctx context.Context, id string, input service.UpdateUserInput) (*service.User, error) {
	if r.updateErr != nil {
		return nil, r.updateErr
	}
	user := r.users[id]
	if user == nil {
		return nil, nil
	}
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	return user, nil
}

func (r *resolverUserRepo) Delete(ctx context.Context, id string) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	delete(r.users, id)
	return nil
}

func TestCreateUser_ReturnsSuccess(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())

	result, err := resolver.Mutation().CreateUser(adminCtx(), model.CreateUserInput{
		Email:    "created@example.com",
		Name:     "Created User",
		Password: "plain-password",
	})

	require.NoError(t, err)
	success, ok := result.(model.CreateUserSuccess)
	require.True(t, ok)
	assert.Equal(t, "created@example.com", success.User.Email)
}

func TestCreateUser_ReturnsValidationErrorOnDuplicate(t *testing.T) {
	repo := newResolverRepo()
	repo.createErr = errors.New("duplicate email")
	resolver := newTestResolver(repo)

	result, err := resolver.Mutation().CreateUser(adminCtx(), model.CreateUserInput{
		Email:    "duplicate@example.com",
		Name:     "Duplicate",
		Password: "plain-password",
	})

	require.NoError(t, err)
	validation, ok := result.(model.ValidationError)
	require.True(t, ok)
	assert.Equal(t, "email", validation.Field)
	assert.Equal(t, "already exists", validation.Message)
}

func TestCreateUser_ReturnsUnexpectedError(t *testing.T) {
	repo := newResolverRepo()
	repo.createErr = errors.New("database down")
	resolver := newTestResolver(repo)

	result, err := resolver.Mutation().CreateUser(adminCtx(), model.CreateUserInput{
		Email:    "bad@example.com",
		Name:     "Bad",
		Password: "plain-password",
	})

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateUser_ReturnsSuccess(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())
	name := "Updated User"
	email := "updated@example.com"

	result, err := resolver.Mutation().UpdateUser(adminCtx(), "user-1", model.UpdateUserInput{
		Name:  &name,
		Email: &email,
	})

	require.NoError(t, err)
	success, ok := result.(model.UpdateUserSuccess)
	require.True(t, ok)
	assert.Equal(t, "Updated User", success.User.Name)
	assert.Equal(t, "updated@example.com", success.User.Email)
}

func TestUpdateUser_ReturnsNotFound(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())
	name := "Missing"

	result, err := resolver.Mutation().UpdateUser(adminCtx(), "missing", model.UpdateUserInput{Name: &name})

	require.NoError(t, err)
	notFound, ok := result.(model.NotFoundError)
	require.True(t, ok)
	assert.Equal(t, "User", notFound.EntityType)
	assert.Equal(t, "missing", notFound.ID)
}

func TestUpdateUser_ReturnsValidationErrorOnDuplicate(t *testing.T) {
	repo := newResolverRepo()
	repo.updateErr = errors.New("duplicate email")
	resolver := newTestResolver(repo)
	email := "duplicate@example.com"

	result, err := resolver.Mutation().UpdateUser(adminCtx(), "user-1", model.UpdateUserInput{Email: &email})

	require.NoError(t, err)
	validation, ok := result.(model.ValidationError)
	require.True(t, ok)
	assert.Equal(t, "email", validation.Field)
}

func TestUpdateUser_ReturnsUnexpectedError(t *testing.T) {
	repo := newResolverRepo()
	repo.updateErr = errors.New("database down")
	resolver := newTestResolver(repo)

	result, err := resolver.Mutation().UpdateUser(adminCtx(), "user-1", model.UpdateUserInput{})

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDeleteUser_ReturnsTrueForExistingUser(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())

	result, err := resolver.Mutation().DeleteUser(adminCtx(), "user-1")

	require.NoError(t, err)
	deleted, ok := result.(model.DeleteUserSuccess)
	require.True(t, ok)
	assert.True(t, deleted.Ok)
}

func TestDeleteUser_ReturnsFalseForMissingUser(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())

	result, err := resolver.Mutation().DeleteUser(adminCtx(), "missing")

	require.NoError(t, err)
	deleted, ok := result.(model.DeleteUserSuccess)
	require.True(t, ok)
	assert.False(t, deleted.Ok)
}

func TestDeleteUser_ReturnsLookupError(t *testing.T) {
	repo := newResolverRepo()
	repo.err = errors.New("lookup failed")
	resolver := newTestResolver(repo)

	result, err := resolver.Mutation().DeleteUser(adminCtx(), "user-1")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDeleteUser_ReturnsDeleteError(t *testing.T) {
	repo := newResolverRepo()
	repo.deleteErr = errors.New("delete failed")
	resolver := newTestResolver(repo)

	result, err := resolver.Mutation().DeleteUser(adminCtx(), "user-1")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestUsers_ReturnsPagination(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())
	first := 1

	connection, err := resolver.Query().Users(adminCtx(), &model.PaginationInput{First: &first})

	require.NoError(t, err)
	require.NotNil(t, connection)
	assert.Equal(t, 1, connection.TotalCount)
	assert.False(t, connection.PageInfo.HasNextPage)
	require.Len(t, connection.Edges, 1)
	assert.NotEmpty(t, connection.Edges[0].Cursor)
}

func TestUsers_ReturnsError(t *testing.T) {
	repo := newResolverRepo()
	repo.err = errors.New("list failed")
	resolver := newTestResolver(repo)

	connection, err := resolver.Query().Users(adminCtx(), nil)

	require.Error(t, err)
	assert.Nil(t, connection)
}

func TestUser_ReturnsNilForMissingUser(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())

	user, err := resolver.Query().User(adminCtx(), "missing")

	require.NoError(t, err)
	assert.Nil(t, user)
}

func TestUser_ReturnsUser(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())

	user, err := resolver.Query().User(adminCtx(), "user-1")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "one@example.com", user.Email)
}

func TestUser_ReturnsError(t *testing.T) {
	repo := newResolverRepo()
	repo.err = errors.New("get failed")
	resolver := newTestResolver(repo)

	user, err := resolver.Query().User(adminCtx(), "user-1")

	require.Error(t, err)
	assert.Nil(t, user)
}

func TestUsers_ReturnsAuthErrorWithoutAdminPrincipal(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())
	connection, err := resolver.Query().Users(context.Background(), nil)
	require.Error(t, err)
	assert.Nil(t, connection)
	assert.Contains(t, err.Error(), "admin authentication required")
}

func TestProtectedUserResolvers_ReturnAuthErrorWithoutAdminPrincipal(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())
	ctx := context.Background()

	user, err := resolver.Query().User(ctx, "user-1")
	assert.Nil(t, user)
	assertAuthRequired(t, err)

	connection, err := resolver.Query().Users(ctx, nil)
	assert.Nil(t, connection)
	assertAuthRequired(t, err)

	created, err := resolver.Mutation().CreateUser(ctx, model.CreateUserInput{Email: "new@example.com", Name: "New", Password: "StrongPassword123!"})
	require.NoError(t, err)
	authErr, ok := created.(model.AuthError)
	require.True(t, ok)
	assert.Contains(t, authErr.Message, "admin authentication required")

	name := "Updated"
	updated, err := resolver.Mutation().UpdateUser(ctx, "user-1", model.UpdateUserInput{Name: &name})
	require.NoError(t, err)
	updateAuthErr, ok := updated.(model.AuthError)
	require.True(t, ok)
	assert.Contains(t, updateAuthErr.Message, "admin authentication required")

	deleted, err := resolver.Mutation().DeleteUser(ctx, "user-1")
	require.NoError(t, err)
	deleteAuthErr, ok := deleted.(model.AuthError)
	require.True(t, ok)
	assert.Contains(t, deleteAuthErr.Message, "admin authentication required")
}

func TestUsers_AllowsAdminPrincipal(t *testing.T) {
	resolver := newTestResolver(newResolverRepo())
	ctx := middleware.ContextWithAdminPrincipal(context.Background(), middleware.AdminPrincipal{ID: "admin-1", Email: "admin@example.com", Role: "ADMIN"})

	connection, err := resolver.Query().Users(ctx, nil)

	require.NoError(t, err)
	require.NotNil(t, connection)
	assert.Equal(t, 1, connection.TotalCount)
}

func adminCtx() context.Context {
	return middleware.ContextWithAdminPrincipal(context.Background(), middleware.AdminPrincipal{
		ID:        "admin-1",
		Email:     "admin@example.com",
		Name:      "Admin",
		Role:      "ADMIN",
		CreatedAt: "2026-06-07T00:00:00Z",
		UpdatedAt: "2026-06-07T00:00:00Z",
	})
}

func assertAuthRequired(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admin authentication required")
}
