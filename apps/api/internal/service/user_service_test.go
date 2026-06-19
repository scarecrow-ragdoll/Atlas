// FILE: apps/api/internal/service/user_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify user service behavior at the repository boundary.
//   SCOPE: Password hashing, validation, repository error mapping, update/delete behavior, and service helper errors; excludes HTTP and PostgreSQL persistence.
//   DEPENDS: internal/service, bcrypt, testify.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   user service tests - Prove transport-neutral user service behavior and repository error translation.
// END_MODULE_MAP
//
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added service coverage for duplicate email and helper error mapping.
// END_CHANGE_SUMMARY

package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"monorepo-template/apps/api/internal/service"
)

type mockUserRepo struct {
	users           map[string]*service.User
	err             error
	lastCreateInput service.CreateUserInput
	lastUpdateID    string
	lastUpdateInput service.UpdateUserInput
	deletedIDs      []string
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*service.User)}
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*service.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	u, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepo) List(ctx context.Context, first *int, after *string) ([]*service.User, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	result := make([]*service.User, 0, len(m.users))
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, len(result), nil
}

func (m *mockUserRepo) Create(ctx context.Context, input service.CreateUserInput) (*service.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.lastCreateInput = input
	u := &service.User{ID: "test-id", Email: input.Email, Name: input.Name}
	m.users[u.ID] = u
	return u, nil
}

func (m *mockUserRepo) Update(ctx context.Context, id string, input service.UpdateUserInput) (*service.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.lastUpdateID = id
	m.lastUpdateInput = input
	u, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	if input.Name != nil {
		u.Name = *input.Name
	}
	if input.Email != nil {
		u.Email = *input.Email
	}
	return u, nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	m.deletedIDs = append(m.deletedIDs, id)
	delete(m.users, id)
	return nil
}

func TestUserService_Create_HashesPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := service.NewUserService(repo)

	user, err := svc.Create(context.Background(), service.CreateUserInput{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "plain-password",
	})

	require.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)
	require.NotEqual(t, "plain-password", repo.lastCreateInput.Password)
	require.NoError(t, bcrypt.CompareHashAndPassword(
		[]byte(repo.lastCreateInput.Password),
		[]byte("plain-password"),
	))
}

func TestUserService_Create_ReturnsHashError(t *testing.T) {
	repo := newMockUserRepo()
	svc := service.NewUserService(repo)

	_, err := svc.Create(context.Background(), service.CreateUserInput{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: string(make([]byte, 73)),
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "hash password")
}

func TestUserService_Create_MapsDuplicateEmail(t *testing.T) {
	repo := newMockUserRepo()
	repo.err = service.ErrDuplicateEmail
	svc := service.NewUserService(repo)

	_, err := svc.Create(context.Background(), service.CreateUserInput{
		Email:    "taken@example.com",
		Name:     "Taken",
		Password: "plain-password",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrDuplicateEmail)
}

func TestUserService_GetByID(t *testing.T) {
	repo := newMockUserRepo()
	svc := service.NewUserService(repo)
	created, err := svc.Create(context.Background(), service.CreateUserInput{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "plain-password",
	})
	require.NoError(t, err)

	found, err := svc.GetByID(context.Background(), created.ID)

	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
}

func TestUserService_GetByID_ReturnsNilWhenMissing(t *testing.T) {
	svc := service.NewUserService(newMockUserRepo())

	found, err := svc.GetByID(context.Background(), "missing")

	require.NoError(t, err)
	assert.Nil(t, found)
}

func TestUserService_List_DelegatesToRepo(t *testing.T) {
	repo := newMockUserRepo()
	repo.users["one"] = &service.User{ID: "one"}
	svc := service.NewUserService(repo)

	users, total, err := svc.List(context.Background(), ptr(10), nil)

	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, users, 1)
}

func TestUserService_Update_DelegatesToRepo(t *testing.T) {
	repo := newMockUserRepo()
	repo.users["test-id"] = &service.User{ID: "test-id", Email: "old@example.com", Name: "Old"}
	svc := service.NewUserService(repo)
	name := "Updated"

	user, err := svc.Update(context.Background(), "test-id", service.UpdateUserInput{Name: &name})

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "test-id", repo.lastUpdateID)
	assert.Equal(t, "Updated", user.Name)
}

func TestUserService_Update_ReturnsNilWhenMissing(t *testing.T) {
	svc := service.NewUserService(newMockUserRepo())
	name := "Updated"

	user, err := svc.Update(context.Background(), "missing", service.UpdateUserInput{Name: &name})

	require.NoError(t, err)
	assert.Nil(t, user)
}

func TestUserService_Update_MapsDuplicateEmail(t *testing.T) {
	repo := newMockUserRepo()
	repo.err = errors.New("duplicate email violates unique index")
	svc := service.NewUserService(repo)

	_, err := svc.Update(context.Background(), "test-id", service.UpdateUserInput{})

	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrDuplicateEmail)
}

func TestUserService_Delete_DelegatesToRepo(t *testing.T) {
	repo := newMockUserRepo()
	repo.users["test-id"] = &service.User{ID: "test-id"}
	svc := service.NewUserService(repo)

	require.NoError(t, svc.Delete(context.Background(), "test-id"))
	assert.Contains(t, repo.deletedIDs, "test-id")
}

func TestUserService_PropagatesRepoErrors(t *testing.T) {
	repoErr := errors.New("repo failed")
	tests := []struct {
		name string
		run  func(*service.UserService) error
	}{
		{
			name: "get",
			run: func(svc *service.UserService) error {
				_, err := svc.GetByID(context.Background(), "test-id")
				return err
			},
		},
		{
			name: "list",
			run: func(svc *service.UserService) error {
				_, _, err := svc.List(context.Background(), nil, nil)
				return err
			},
		},
		{
			name: "create",
			run: func(svc *service.UserService) error {
				_, err := svc.Create(context.Background(), service.CreateUserInput{
					Email:    "test@example.com",
					Name:     "Test",
					Password: "plain-password",
				})
				return err
			},
		},
		{
			name: "update",
			run: func(svc *service.UserService) error {
				_, err := svc.Update(context.Background(), "test-id", service.UpdateUserInput{})
				return err
			},
		},
		{
			name: "delete",
			run: func(svc *service.UserService) error {
				return svc.Delete(context.Background(), "test-id")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepo()
			repo.err = repoErr
			svc := service.NewUserService(repo)

			err := tt.run(svc)

			require.Error(t, err)
			assert.ErrorIs(t, err, repoErr)
		})
	}
}

func TestUserService_ErrorHelpers(t *testing.T) {
	assert.False(t, service.IsDuplicateEmail(nil))
	assert.True(t, service.IsNotFound(service.ErrNotFound))
	assert.False(t, service.IsNotFound(errors.New("plain error")))
}

func ptr[T any](v T) *T { return &v }
