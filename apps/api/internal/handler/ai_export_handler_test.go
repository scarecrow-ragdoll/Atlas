package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
	"monorepo-template/apps/api/internal/handler"
)

var (
	handlerCtx       = context.Background()
	handlerUserID    = "550e8400-e29b-41d4-a716-446655440000"
	handlerExportID  = "660e8400-e29b-41d4-a716-446655440001"
)

func authCtxForUser(userID string) context.Context {
	return middleware.ContextWithAtlasUserID(context.Background(), userID)
}

func authMiddlewareForUser(userID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := middleware.ContextWithAtlasUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type mockAiExportSvc struct {
	atlasService.AiExportService
	generateFn func(ctx context.Context, userID string, input models.CreateAiExportInput, maxRangeDays int, maxExportSize int64, exportBasePath string) (*models.AiExport, string, error)
	getByIDFn  func(ctx context.Context, userID string, id string) (*models.AiExport, error)
	listFn     func(ctx context.Context, userID string) ([]models.AiExport, error)
	deleteFn   func(ctx context.Context, userID string, id string) (*models.AiExport, error)
}

func (m *mockAiExportSvc) Generate(ctx context.Context, userID string, input models.CreateAiExportInput, maxRangeDays int, maxExportSize int64, exportBasePath string) (*models.AiExport, string, error) {
	return m.generateFn(ctx, userID, input, maxRangeDays, maxExportSize, exportBasePath)
}

func (m *mockAiExportSvc) GetByID(ctx context.Context, userID string, id string) (*models.AiExport, error) {
	return m.getByIDFn(ctx, userID, id)
}

func (m *mockAiExportSvc) List(ctx context.Context, userID string) ([]models.AiExport, error) {
	return m.listFn(ctx, userID)
}

func (m *mockAiExportSvc) Delete(ctx context.Context, userID string, id string) (*models.AiExport, error) {
	return m.deleteFn(ctx, userID, id)
}

type mockUserProfileSvc struct {
	atlasService.UserProfileService
	getFn    func(ctx context.Context, userID string) (*models.UserProfile, error)
	updateFn func(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfile, error)
}

func (m *mockUserProfileSvc) Get(ctx context.Context, userID string) (*models.UserProfile, error) {
	return m.getFn(ctx, userID)
}

func (m *mockUserProfileSvc) Update(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfile, error) {
	return m.updateFn(ctx, userID, input)
}

func TestAiExportHandler_GenerateExport(t *testing.T) {
	mockExport := &models.AiExport{
		ID:              handlerExportID,
		UserID:          handlerUserID,
		GeneratedPrompt: "test prompt",
	}
	svc := &mockAiExportSvc{
		generateFn: func(ctx context.Context, userID string, input models.CreateAiExportInput, maxRangeDays int, maxExportSize int64, exportBasePath string) (*models.AiExport, string, error) {
			return mockExport, "test prompt", nil
		},
	}
	profileSvc := &mockUserProfileSvc{}
	h := handler.NewAiExportHandler(svc, profileSvc, "./test-exports", 365, 104857600)

	body := `{"dateRangeStart":"2026-01-01","dateRangeEnd":"2026-01-28"}`
	req := httptest.NewRequest("POST", "/api/ai-export/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(authCtxForUser(handlerUserID))

	rr := httptest.NewRecorder()
	h.GenerateExport(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]any
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp, "export")
	export := resp["export"].(map[string]any)
	assert.Equal(t, handlerExportID, export["id"])
	assert.Equal(t, "test prompt", export["generatedPrompt"])
}

func TestAiExportHandler_GenerateExport_MissingAuth(t *testing.T) {
	svc := &mockAiExportSvc{}
	profileSvc := &mockUserProfileSvc{}
	h := handler.NewAiExportHandler(svc, profileSvc, "./test-exports", 365, 104857600)

	body := `{"dateRangeStart":"2026-01-01","dateRangeEnd":"2026-01-28"}`
	req := httptest.NewRequest("POST", "/api/ai-export/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	h.GenerateExport(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAiExportHandler_DownloadExport(t *testing.T) {
	tmpDir := t.TempDir()
	exportDir := filepath.Join(tmpDir, handlerUserID)
	err := os.MkdirAll(exportDir, 0755)
	require.NoError(t, err)

	filePath := filepath.Join(exportDir, handlerExportID+".zip")
	err = os.WriteFile(filePath, []byte("fake-zip-content"), 0644)
	require.NoError(t, err)

	filePathStr := filePath
	mockExport := &models.AiExport{
		ID:             handlerExportID,
		UserID:         handlerUserID,
		ExportFilePath: &filePathStr,
	}

	svc := &mockAiExportSvc{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.AiExport, error) {
			return mockExport, nil
		},
	}
	profileSvc := &mockUserProfileSvc{}
	h := handler.NewAiExportHandler(svc, profileSvc, tmpDir, 365, 104857600)

	r := chi.NewRouter()
	r.Use(authMiddlewareForUser(handlerUserID))
	r.Get("/api/ai-export/download", h.DownloadExport)

	req := httptest.NewRequest("GET", "/api/ai-export/download?exportId="+handlerExportID, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/zip", rr.Header().Get("Content-Type"))
	assert.Contains(t, rr.Header().Get("Content-Disposition"), "attachment")
}

func TestAiExportHandler_DownloadExport_NotFound(t *testing.T) {
	svc := &mockAiExportSvc{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.AiExport, error) {
			return nil, atlasService.ErrAiExportNotFound
		},
	}
	profileSvc := &mockUserProfileSvc{}
	h := handler.NewAiExportHandler(svc, profileSvc, "./test-exports", 365, 104857600)

	r := chi.NewRouter()
	r.Use(authMiddlewareForUser(handlerUserID))
	r.Get("/api/ai-export/download", h.DownloadExport)

	req := httptest.NewRequest("GET", "/api/ai-export/download?exportId=nonexistent", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestAiExportHandler_DownloadExport_OwnershipMismatch(t *testing.T) {
	otherUserID := "770e8400-e29b-41d4-a716-446655449999"

	svc := &mockAiExportSvc{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.AiExport, error) {
			return nil, atlasService.ErrAiExportNotFound
		},
	}
	profileSvc := &mockUserProfileSvc{}
	h := handler.NewAiExportHandler(svc, profileSvc, "./test-exports", 365, 104857600)

	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := middleware.ContextWithAtlasUserID(r.Context(), otherUserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Get("/api/ai-export/download", h.DownloadExport)

	req := httptest.NewRequest("GET", "/api/ai-export/download?exportId="+handlerExportID, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUserProfileHandler_Get(t *testing.T) {
	goal := "Build muscle"
	mockProfile := &models.UserProfile{
		ID:     handlerExportID,
		UserID: handlerUserID,
		Goal:   &goal,
	}

	svc := &mockAiExportSvc{}
	profileSvc := &mockUserProfileSvc{
		getFn: func(ctx context.Context, userID string) (*models.UserProfile, error) {
			return mockProfile, nil
		},
	}
	h := handler.NewAiExportHandler(svc, profileSvc, "./test-exports", 365, 104857600)

	r := chi.NewRouter()
	r.Use(authMiddlewareForUser(handlerUserID))
	r.Get("/api/user-profile", h.GetUserProfile)

	req := httptest.NewRequest("GET", "/api/user-profile", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var profile models.UserProfile
	err := json.Unmarshal(rr.Body.Bytes(), &profile)
	require.NoError(t, err)
	assert.Equal(t, "Build muscle", *profile.Goal)
}

func TestUserProfileHandler_Get_NoSession(t *testing.T) {
	svc := &mockAiExportSvc{}
	profileSvc := &mockUserProfileSvc{}
	h := handler.NewAiExportHandler(svc, profileSvc, "./test-exports", 365, 104857600)

	req := httptest.NewRequest("GET", "/api/user-profile", nil)
	rr := httptest.NewRecorder()
	h.GetUserProfile(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}