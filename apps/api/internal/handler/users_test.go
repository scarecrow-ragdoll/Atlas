// FILE: apps/api/internal/handler/users_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify observable users HTTP handler behavior and service error mapping.
//   SCOPE: Users route responses, JSON envelopes, mount behavior, and handler-local service fakes; excludes repository persistence.
//   DEPENDS: internal/handler, internal/service, chi, httptest.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   users handler tests - Prove list, create, get, update, delete, and route-mount behavior at the HTTP boundary.
// END_MODULE_MAP
//
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added users handler coverage for service mapping, success, and error paths.
// END_CHANGE_SUMMARY

package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/handler"
	"monorepo-template/apps/api/internal/service"
)

func TestUsersHandler_NewUsersHandler_DefaultsNilLogger(t *testing.T) {
	svc := newFakeUsersService()

	rec := serveUsersRequestWithLogger(t, svc, nil, http.MethodGet, "/api/users", nil)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUsersHandler_ListUsers_ReturnsDataAndTotalCount(t *testing.T) {
	svc := newFakeUsersService()
	svc.users["u1"] = &service.User{ID: "u1", Email: "one@example.com", Name: "One", CreatedAt: "2026-05-24T00:00:00Z", UpdatedAt: "2026-05-24T00:00:00Z"}
	rec := serveUsersRequest(t, svc, http.MethodGet, "/api/users", nil)
	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"data":[{"id":"u1","email":"one@example.com","name":"One","createdAt":"2026-05-24T00:00:00Z","updatedAt":"2026-05-24T00:00:00Z"}],"meta":{"totalCount":1}}`, rec.Body.String())
}

func TestUsersHandler_ListUsers_MapsServiceError(t *testing.T) {
	svc := newFakeUsersService()
	svc.listErr = errors.New("list failed")

	rec := serveUsersRequest(t, svc, http.MethodGet, "/api/users", nil)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.JSONEq(t, `{"error":{"code":"INTERNAL_ERROR","message":"internal error"}}`, rec.Body.String())
}

func TestUsersHandler_CreateUser_ReturnsCreatedUser(t *testing.T) {
	svc := newFakeUsersService()
	rec := serveUsersRequest(t, svc, http.MethodPost, "/api/users", map[string]string{"email": "new@example.com", "name": "New", "password": "secret123"})
	require.Equal(t, http.StatusCreated, rec.Code)
	require.JSONEq(t, `{"data":{"id":"created-id","email":"new@example.com","name":"New","createdAt":"2026-05-24T00:00:00Z","updatedAt":"2026-05-24T00:00:00Z"}}`, rec.Body.String())
}

func TestUsersHandler_CreateUser_RejectsInvalidJSON(t *testing.T) {
	svc := newFakeUsersService()

	rec := serveUsersRawRequest(t, svc, http.MethodPost, "/api/users", []byte(`{"email":`))

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.JSONEq(t, `{"error":{"code":"BAD_REQUEST","message":"invalid json body"}}`, rec.Body.String())
}

func TestUsersHandler_CreateUser_MapsDuplicateEmail(t *testing.T) {
	svc := newFakeUsersService()
	svc.createErr = service.ErrDuplicateEmail
	rec := serveUsersRequest(t, svc, http.MethodPost, "/api/users", map[string]string{"email": "taken@example.com", "name": "Taken", "password": "secret123"})
	require.Equal(t, http.StatusConflict, rec.Code)
	require.JSONEq(t, `{"error":{"code":"DUPLICATE_EMAIL","message":"email already exists","field":"email"}}`, rec.Body.String())
}

func TestUsersHandler_CreateUser_MapsServiceError(t *testing.T) {
	svc := newFakeUsersService()
	svc.createErr = errors.New("create failed")

	rec := serveUsersRequest(t, svc, http.MethodPost, "/api/users", map[string]string{"email": "new@example.com", "name": "New", "password": "secret123"})

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.JSONEq(t, `{"error":{"code":"INTERNAL_ERROR","message":"internal error"}}`, rec.Body.String())
}

func TestUsersHandler_GetUser_MapsServiceError(t *testing.T) {
	svc := newFakeUsersService()
	svc.getErr = errors.New("get failed")

	rec := serveUsersRequest(t, svc, http.MethodGet, "/api/users/u1", nil)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.JSONEq(t, `{"error":{"code":"INTERNAL_ERROR","message":"internal error"}}`, rec.Body.String())
}

func TestUsersHandler_GetUser_ReturnsUser(t *testing.T) {
	svc := newFakeUsersService()
	svc.users["u1"] = &service.User{ID: "u1", Email: "one@example.com", Name: "One", CreatedAt: "2026-05-24T00:00:00Z", UpdatedAt: "2026-05-24T00:00:00Z"}

	rec := serveUsersRequest(t, svc, http.MethodGet, "/api/users/u1", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"data":{"id":"u1","email":"one@example.com","name":"One","createdAt":"2026-05-24T00:00:00Z","updatedAt":"2026-05-24T00:00:00Z"}}`, rec.Body.String())
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

func TestUsersHandler_UpdateUser_RejectsInvalidJSON(t *testing.T) {
	svc := newFakeUsersService()

	rec := serveUsersRawRequest(t, svc, http.MethodPatch, "/api/users/u1", []byte(`{"email":`))

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.JSONEq(t, `{"error":{"code":"BAD_REQUEST","message":"invalid json body"}}`, rec.Body.String())
}

func TestUsersHandler_UpdateUser_MapsServiceNotFoundError(t *testing.T) {
	svc := newFakeUsersService()
	svc.updateErr = service.ErrNotFound

	rec := serveUsersRequest(t, svc, http.MethodPatch, "/api/users/u1", map[string]string{"name": "New"})

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.JSONEq(t, `{"error":{"code":"NOT_FOUND","message":"user not found"}}`, rec.Body.String())
}

func TestUsersHandler_UpdateUser_MapsServiceError(t *testing.T) {
	svc := newFakeUsersService()
	svc.updateErr = errors.New("update failed")

	rec := serveUsersRequest(t, svc, http.MethodPatch, "/api/users/u1", map[string]string{"name": "New"})

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.JSONEq(t, `{"error":{"code":"INTERNAL_ERROR","message":"internal error"}}`, rec.Body.String())
}

func TestUsersHandler_UpdateUser_ReturnsNotFoundForMissingUser(t *testing.T) {
	svc := newFakeUsersService()

	rec := serveUsersRequest(t, svc, http.MethodPatch, "/api/users/missing", map[string]string{"name": "New"})

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.JSONEq(t, `{"error":{"code":"NOT_FOUND","message":"user not found"}}`, rec.Body.String())
}

func TestUsersHandler_DeleteUser_ReturnsNoContentForExistingUser(t *testing.T) {
	svc := newFakeUsersService()
	svc.users["u1"] = &service.User{ID: "u1", Email: "one@example.com", Name: "One"}
	rec := serveUsersRequest(t, svc, http.MethodDelete, "/api/users/u1", nil)
	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Empty(t, rec.Body.String())
}

func TestUsersHandler_DeleteUser_MapsGetServiceError(t *testing.T) {
	svc := newFakeUsersService()
	svc.getErr = errors.New("get failed")

	rec := serveUsersRequest(t, svc, http.MethodDelete, "/api/users/u1", nil)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.JSONEq(t, `{"error":{"code":"INTERNAL_ERROR","message":"internal error"}}`, rec.Body.String())
}

func TestUsersHandler_DeleteUser_MapsDeleteServiceError(t *testing.T) {
	svc := newFakeUsersService()
	svc.users["u1"] = &service.User{ID: "u1", Email: "one@example.com", Name: "One"}
	svc.deleteErr = errors.New("delete failed")

	rec := serveUsersRequest(t, svc, http.MethodDelete, "/api/users/u1", nil)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.JSONEq(t, `{"error":{"code":"INTERNAL_ERROR","message":"internal error"}}`, rec.Body.String())
}

func TestUsersHandler_DeleteUser_ReturnsNotFound(t *testing.T) {
	svc := newFakeUsersService()
	rec := serveUsersRequest(t, svc, http.MethodDelete, "/api/users/missing", nil)
	require.Equal(t, http.StatusNotFound, rec.Code)
	require.JSONEq(t, `{"error":{"code":"NOT_FOUND","message":"user not found"}}`, rec.Body.String())
}

func TestUsersHandler_Routes_AreMountRelative(t *testing.T) {
	svc := newFakeUsersService()
	rec := serveUsersRequestWithoutMount(t, svc, http.MethodGet, "/api/users", nil)
	require.Equal(t, http.StatusNotFound, rec.Code)
}

type fakeUsersService struct {
	users     map[string]*service.User
	getErr    error
	listErr   error
	createErr error
	updateErr error
	deleteErr error
}

func newFakeUsersService() *fakeUsersService {
	return &fakeUsersService{users: map[string]*service.User{}}
}

func serveUsersRequestWithLogger(t *testing.T, svc *fakeUsersService, logger *zap.Logger, method string, path string, body any) *httptest.ResponseRecorder {
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
	r := chi.NewRouter()
	r.Mount("/api/users", handler.NewUsersHandler(svc, logger).Routes())
	r.ServeHTTP(rec, req)
	return rec
}

func serveUsersRequest(t *testing.T, svc *fakeUsersService, method string, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	return serveUsersRequestWithLogger(t, svc, zap.NewNop(), method, path, body)
}

func serveUsersRawRequest(t *testing.T, svc *fakeUsersService, method string, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Mount("/api/users", handler.NewUsersHandler(svc, zap.NewNop()).Routes())
	r.ServeHTTP(rec, req)
	return rec
}

func serveUsersRequestWithoutMount(t *testing.T, svc *fakeUsersService, method string, path string, body any) *httptest.ResponseRecorder {
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

func (s *fakeUsersService) GetByID(_ context.Context, id string) (*service.User, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.users[id], nil
}

func (s *fakeUsersService) List(_ context.Context, _ *int, _ *string) ([]*service.User, int, error) {
	if s.listErr != nil {
		return nil, 0, s.listErr
	}
	users := make([]*service.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, len(users), nil
}

func (s *fakeUsersService) Create(_ context.Context, input service.CreateUserInput) (*service.User, error) {
	if s.createErr != nil {
		return nil, s.createErr
	}
	return &service.User{
		ID:        "created-id",
		Email:     input.Email,
		Name:      input.Name,
		CreatedAt: "2026-05-24T00:00:00Z",
		UpdatedAt: "2026-05-24T00:00:00Z",
	}, nil
}

func (s *fakeUsersService) Update(_ context.Context, id string, input service.UpdateUserInput) (*service.User, error) {
	if s.updateErr != nil {
		return nil, s.updateErr
	}
	user := s.users[id]
	if user == nil {
		return nil, nil
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Name != nil {
		user.Name = *input.Name
	}
	user.UpdatedAt = "2026-05-24T00:00:01Z"
	return user, nil
}

func (s *fakeUsersService) Delete(_ context.Context, id string) error {
	if s.deleteErr != nil {
		return s.deleteErr
	}
	delete(s.users, id)
	return nil
}
