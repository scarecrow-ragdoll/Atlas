package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/service"
)

type UsersService interface {
	GetByID(ctx context.Context, id string) (*service.User, error)
	List(ctx context.Context, first *int, after *string) ([]*service.User, int, error)
	Create(ctx context.Context, input service.CreateUserInput) (*service.User, error)
	Update(ctx context.Context, id string, input service.UpdateUserInput) (*service.User, error)
	Delete(ctx context.Context, id string) error
}

type UsersHandler struct {
	service UsersService
	logger  *zap.Logger
}

func NewUsersHandler(service UsersService, logger *zap.Logger) *UsersHandler {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &UsersHandler{service: service, logger: logger}
}

func (h *UsersHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.listUsers)
	r.Post("/", h.createUser)
	r.Get("/{id}", h.getUser)
	r.Patch("/{id}", h.updateUser)
	r.Delete("/{id}", h.deleteUser)
	return r
}

type userResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type dataEnvelope struct {
	Data any `json:"data"`
}

type listEnvelope struct {
	Data any       `json:"data"`
	Meta metaBlock `json:"meta"`
}

type metaBlock struct {
	TotalCount int `json:"totalCount"`
}

type errorEnvelope struct {
	Error errorBlock `json:"error"`
}

type errorBlock struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Field   *string `json:"field,omitempty"`
}

type createUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Email *string `json:"email"`
	Name  *string `json:"name"`
}

func (h *UsersHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, total, err := h.service.List(r.Context(), nil, nil)
	if err != nil {
		h.logger.Error("list users failed", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error", nil)
		return
	}

	data := make([]userResponse, 0, len(users))
	for _, user := range users {
		data = append(data, mapUserResponse(user))
	}
	writeJSON(w, http.StatusOK, listEnvelope{Data: data, Meta: metaBlock{TotalCount: total}})
}

func (h *UsersHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid json body", nil)
		return
	}

	user, err := h.service.Create(r.Context(), service.CreateUserInput{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, dataEnvelope{Data: mapUserResponse(user)})
}

func (h *UsersHandler) getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	if user == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "user not found", nil)
		return
	}
	writeJSON(w, http.StatusOK, dataEnvelope{Data: mapUserResponse(user)})
}

func (h *UsersHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	var req updateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid json body", nil)
		return
	}

	id := chi.URLParam(r, "id")
	user, err := h.service.Update(r.Context(), id, service.UpdateUserInput{
		Email: req.Email,
		Name:  req.Name,
	})
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	if user == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "user not found", nil)
		return
	}
	writeJSON(w, http.StatusOK, dataEnvelope{Data: mapUserResponse(user)})
}

func (h *UsersHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	if user == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "user not found", nil)
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		h.writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UsersHandler) writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case service.IsDuplicateEmail(err):
		field := "email"
		writeError(w, http.StatusConflict, "DUPLICATE_EMAIL", "email already exists", &field)
	case service.IsNotFound(err):
		writeError(w, http.StatusNotFound, "NOT_FOUND", "user not found", nil)
	default:
		h.logger.Error("users handler failed", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error", nil)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, fmt.Sprintf("encode response: %v", err), http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, code string, message string, field *string) {
	writeJSON(w, status, errorEnvelope{Error: errorBlock{Code: code, Message: message, Field: field}})
}

func mapUserResponse(u *service.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}
