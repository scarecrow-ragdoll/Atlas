// FILE: apps/api/internal/handler/backup_handler_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify observable BackupHandler HTTP behavior and service error mapping.
//   SCOPE: GenerateExport, DownloadExport, ImportValidate, ImportConfirm success and error paths; excludes repository persistence.
//   DEPENDS: internal/handler, internal/atlas/service, internal/atlas/models, chi, httptest.
//   LINKS: M-API / V-M-API / WAVE-09.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   backup handler tests - Prove export, download, import validate, and import confirm at the HTTP boundary.
// END_MODULE_MAP
//
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added BackupHandler coverage for success and error paths.
// END_CHANGE_SUMMARY

package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
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

type mockBackupExportService struct {
	atlasService.BackupExportService
	generateFn       func(ctx context.Context, userID string, includeMedia bool, maxExportSize int64, exportBasePath string) (*models.BackupExportResult, error)
	getDownloadPathFn func(ctx context.Context, userID string, downloadID string) (string, error)
}

func (m *mockBackupExportService) Generate(ctx context.Context, userID string, includeMedia bool, maxExportSize int64, exportBasePath string) (*models.BackupExportResult, error) {
	return m.generateFn(ctx, userID, includeMedia, maxExportSize, exportBasePath)
}

func (m *mockBackupExportService) GetDownloadPath(ctx context.Context, userID string, downloadID string) (string, error) {
	return m.getDownloadPathFn(ctx, userID, downloadID)
}

type mockBackupImportService struct {
	atlasService.BackupImportService
	validateFn func(ctx context.Context, userID string, zipData []byte) (string, *models.BackupImportSummary, error)
	confirmFn  func(ctx context.Context, userID string, validationID string) (*models.BackupImportConfirmResult, error)
}

func (m *mockBackupImportService) Validate(ctx context.Context, userID string, zipData []byte) (string, *models.BackupImportSummary, error) {
	return m.validateFn(ctx, userID, zipData)
}

func (m *mockBackupImportService) Confirm(ctx context.Context, userID string, validationID string) (*models.BackupImportConfirmResult, error) {
	return m.confirmFn(ctx, userID, validationID)
}

func TestBackupHandler_GenerateExport_Success(t *testing.T) {
	exportSvc := &mockBackupExportService{
		generateFn: func(ctx context.Context, userID string, includeMedia bool, maxExportSize int64, exportBasePath string) (*models.BackupExportResult, error) {
			return &models.BackupExportResult{
				DownloadID: "exp-001",
				SizeBytes:  1024,
				Timestamp:  "2026-06-22T12:00:00Z",
			}, nil
		},
	}
	importSvc := &mockBackupImportService{}
	h := handler.NewBackupHandler(exportSvc, importSvc, "/tmp/exports", 104857600, 52428800)

	body := `{"includeMedia":true}`
	req := httptest.NewRequest("POST", "/api/backup/export", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.ContextWithAtlasUserID(context.Background(), handlerUserID))

	rr := httptest.NewRecorder()
	h.GenerateExport(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]any
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "exp-001", resp["downloadId"])
	assert.Equal(t, float64(1024), resp["sizeBytes"])
	assert.Equal(t, "2026-06-22T12:00:00Z", resp["timestamp"])
}

func TestBackupHandler_GenerateExport_Unauthorized(t *testing.T) {
	exportSvc := &mockBackupExportService{}
	importSvc := &mockBackupImportService{}
	h := handler.NewBackupHandler(exportSvc, importSvc, "/tmp/exports", 104857600, 52428800)

	body := `{"includeMedia":true}`
	req := httptest.NewRequest("POST", "/api/backup/export", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	h.GenerateExport(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	var resp map[string]any
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	errObj := resp["error"].(map[string]any)
	assert.Equal(t, "UNAUTHORIZED", errObj["code"])
}

func TestBackupHandler_DownloadExport_Success(t *testing.T) {
	tmpDir := t.TempDir()
	userDir := filepath.Join(tmpDir, handlerUserID)
	err := os.MkdirAll(userDir, 0755)
	require.NoError(t, err)

	filePath := filepath.Join(userDir, "atlas-backup-dl-001.zip")
	err = os.WriteFile(filePath, []byte("fake-zip-content"), 0644)
	require.NoError(t, err)

	exportSvc := &mockBackupExportService{
		getDownloadPathFn: func(ctx context.Context, userID string, downloadID string) (string, error) {
			return filePath, nil
		},
	}
	importSvc := &mockBackupImportService{}
	h := handler.NewBackupHandler(exportSvc, importSvc, tmpDir, 104857600, 52428800)

	r := chi.NewRouter()
	r.Use(authMiddlewareForUser(handlerUserID))
	r.Get("/api/backup/download", h.DownloadExport)

	req := httptest.NewRequest("GET", "/api/backup/download?downloadId=dl-001", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/zip", rr.Header().Get("Content-Type"))
	assert.Contains(t, rr.Header().Get("Content-Disposition"), "attachment")
	assert.Equal(t, "16", rr.Header().Get("Content-Length"))
}

func TestBackupHandler_DownloadExport_MissingID(t *testing.T) {
	exportSvc := &mockBackupExportService{}
	importSvc := &mockBackupImportService{}
	h := handler.NewBackupHandler(exportSvc, importSvc, "/tmp/exports", 104857600, 52428800)

	r := chi.NewRouter()
	r.Use(authMiddlewareForUser(handlerUserID))
	r.Get("/api/backup/download", h.DownloadExport)

	req := httptest.NewRequest("GET", "/api/backup/download", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var resp map[string]any
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	errObj := resp["error"].(map[string]any)
	assert.Equal(t, "MISSING_DOWNLOAD_ID", errObj["code"])
}

func TestBackupHandler_DownloadExport_NotFound(t *testing.T) {
	exportSvc := &mockBackupExportService{
		getDownloadPathFn: func(ctx context.Context, userID string, downloadID string) (string, error) {
			return "", atlasService.ErrBackupExportNotFound
		},
	}
	importSvc := &mockBackupImportService{}
	h := handler.NewBackupHandler(exportSvc, importSvc, "/tmp/exports", 104857600, 52428800)

	r := chi.NewRouter()
	r.Use(authMiddlewareForUser(handlerUserID))
	r.Get("/api/backup/download", h.DownloadExport)

	req := httptest.NewRequest("GET", "/api/backup/download?downloadId=nonexistent", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var resp map[string]any
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	errObj := resp["error"].(map[string]any)
	assert.Equal(t, "DOWNLOAD_NOT_FOUND", errObj["code"])
}

func TestBackupHandler_ImportValidate_Success(t *testing.T) {
	importSvc := &mockBackupImportService{
		validateFn: func(ctx context.Context, userID string, zipData []byte) (string, *models.BackupImportSummary, error) {
			return "val-001", &models.BackupImportSummary{
				SchemaVersion: 1,
				AppVersion:    "1.0.0",
				EntityCounts:  map[string]int{"settings": 1},
				MediaCount:    0,
				Warnings:      nil,
			}, nil
		},
	}
	exportSvc := &mockBackupExportService{}
	h := handler.NewBackupHandler(exportSvc, importSvc, "/tmp/exports", 104857600, 52428800)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "backup.zip")
	require.NoError(t, err)
	_, err = fw.Write([]byte("fake-zip-content"))
	require.NoError(t, err)
	w.Close()

	req := httptest.NewRequest("POST", "/api/backup/import/validate", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(middleware.ContextWithAtlasUserID(context.Background(), handlerUserID))

	rr := httptest.NewRecorder()
	h.ImportValidate(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]any
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "val-001", resp["validationId"])
	summary := resp["summary"].(map[string]any)
	assert.Equal(t, float64(1), summary["schemaVersion"])
	assert.Equal(t, "1.0.0", summary["appVersion"])
}

func TestBackupHandler_ImportValidate_MissingFile(t *testing.T) {
	importSvc := &mockBackupImportService{}
	exportSvc := &mockBackupExportService{}
	h := handler.NewBackupHandler(exportSvc, importSvc, "/tmp/exports", 104857600, 52428800)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.Close()

	req := httptest.NewRequest("POST", "/api/backup/import/validate", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(middleware.ContextWithAtlasUserID(context.Background(), handlerUserID))

	rr := httptest.NewRecorder()
	h.ImportValidate(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var resp map[string]any
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	errObj := resp["error"].(map[string]any)
	assert.Equal(t, "MISSING_FILE", errObj["code"])
}

func TestBackupHandler_ImportConfirm_Success(t *testing.T) {
	importSvc := &mockBackupImportService{
		confirmFn: func(ctx context.Context, userID string, validationID string) (*models.BackupImportConfirmResult, error) {
			return &models.BackupImportConfirmResult{
				Status:       "confirmed",
				EntityCounts: map[string]int{"settings": 1},
				MediaCount:   0,
			}, nil
		},
	}
	exportSvc := &mockBackupExportService{}
	h := handler.NewBackupHandler(exportSvc, importSvc, "/tmp/exports", 104857600, 52428800)

	body := `{"validationId":"val-001"}`
	req := httptest.NewRequest("POST", "/api/backup/import/confirm", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.ContextWithAtlasUserID(context.Background(), handlerUserID))

	rr := httptest.NewRecorder()
	h.ImportConfirm(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]any
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "confirmed", resp["status"])
	assert.Equal(t, float64(0), resp["mediaCount"])
}

func TestBackupHandler_ImportConfirm_MissingValidationID(t *testing.T) {
	importSvc := &mockBackupImportService{}
	exportSvc := &mockBackupExportService{}
	h := handler.NewBackupHandler(exportSvc, importSvc, "/tmp/exports", 104857600, 52428800)

	body := `{}`
	req := httptest.NewRequest("POST", "/api/backup/import/confirm", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.ContextWithAtlasUserID(context.Background(), handlerUserID))

	rr := httptest.NewRecorder()
	h.ImportConfirm(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var resp map[string]any
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	errObj := resp["error"].(map[string]any)
	assert.Equal(t, "MISSING_VALIDATION_ID", errObj["code"])
}