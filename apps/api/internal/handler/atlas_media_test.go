// FILE: apps/api/internal/handler/atlas_media_test.go
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Verify Atlas media handler for exercise media upload, download, and delete operations.
//   SCOPE: POST upload (MIME validation, size limits, auth guard), GET download (content-type, not found), DELETE delete (not found). Uses mock ExerciseService.
//   DEPENDS: internal/handler, internal/atlas/middleware, internal/atlas/models, internal/atlas/service, httptest.
//   LINKS: M-API / V-M-API / WAVE-02 / TEST-W02-004 / TEST-W02-015 / TEST-W02-016 / TEST-W02-017.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Updated from scaffold test to full handler unit tests for WAVE-02.
// END_CHANGE_SUMMARY

package handler_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasSvc "monorepo-template/apps/api/internal/atlas/service"

	"monorepo-template/apps/api/internal/handler"
)

type mockExerciseServiceForMedia struct {
	createMediaFn        func(ctx context.Context, userID, exerciseID, fileName, filePath, mimeType string, fileSize int64) (*models.ExerciseMedia, error)
	getMediaByIDFn       func(ctx context.Context, userID string, id string) (*models.ExerciseMedia, error)
	deleteMediaFn        func(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error)
	getMediaRecordByIDFn func(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error)
}

func (m *mockExerciseServiceForMedia) Create(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.Exercise, error) { return nil, nil }

func (m *mockExerciseServiceForMedia) GetByID(ctx context.Context, userID string, id string) (*models.Exercise, error) { return nil, nil }

func (m *mockExerciseServiceForMedia) List(ctx context.Context, userID string, first int32, after *string, includeInactive bool) (*models.ExerciseConnection, error) { return nil, nil }

func (m *mockExerciseServiceForMedia) ListAll(ctx context.Context, userID string, includeInactive bool) ([]models.Exercise, error) { return nil, nil }

func (m *mockExerciseServiceForMedia) Update(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.Exercise, error) { return nil, nil }

func (m *mockExerciseServiceForMedia) Archive(ctx context.Context, userID string, id string) (*models.Exercise, error) { return nil, nil }

func (m *mockExerciseServiceForMedia) Restore(ctx context.Context, userID string, id string) (*models.Exercise, error) { return nil, nil }

func (m *mockExerciseServiceForMedia) CreateMedia(ctx context.Context, userID, exerciseID, fileName, filePath, mimeType string, fileSize int64) (*models.ExerciseMedia, error) {
	return m.createMediaFn(ctx, userID, exerciseID, fileName, filePath, mimeType, fileSize)
}

func (m *mockExerciseServiceForMedia) GetMediaByID(ctx context.Context, userID string, id string) (*models.ExerciseMedia, error) {
	return m.getMediaByIDFn(ctx, userID, id)
}

func (m *mockExerciseServiceForMedia) GetMediaRecordByID(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error) {
	return m.getMediaRecordByIDFn(ctx, userID, id)
}

func (m *mockExerciseServiceForMedia) DeleteMedia(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error) {
	return m.deleteMediaFn(ctx, userID, id)
}

func TestAtlasMediaUpload_RequiresAuth(t *testing.T) {
	h := handler.NewAtlasMediaHandler(&mockExerciseServiceForMedia{}, t.TempDir())
	req := httptest.NewRequest(http.MethodPost, "/api/v1/media/upload", nil)
	rec := httptest.NewRecorder()
	h.Upload(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAtlasMediaUpload_RejectsInvalidFileType(t *testing.T) {
	h := handler.NewAtlasMediaHandler(&mockExerciseServiceForMedia{}, t.TempDir())

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("purpose", "EXERCISE_MEDIA")
	w.WriteField("exerciseId", "00000000-0000-0000-0000-000000000001")
	fw, _ := w.CreateFormFile("file", "test.txt")
	fw.Write([]byte("this is plain text"))
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/media/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(authCtx())
	rec := httptest.NewRecorder()
	h.Upload(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "INVALID_FILE_TYPE")
}

func TestAtlasMediaUpload_RejectsMissingExerciseID(t *testing.T) {
	h := handler.NewAtlasMediaHandler(&mockExerciseServiceForMedia{}, t.TempDir())

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("purpose", "EXERCISE_MEDIA")
	fw, _ := w.CreateFormFile("file", "test.jpg")
	fw.Write([]byte("not a real jpeg"))
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/media/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(authCtx())
	rec := httptest.NewRecorder()
	h.Upload(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAtlasMediaUpload_RejectsInvalidPurpose(t *testing.T) {
	h := handler.NewAtlasMediaHandler(&mockExerciseServiceForMedia{}, t.TempDir())

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("purpose", "WRONG_PURPOSE")
	w.WriteField("exerciseId", "00000000-0000-0000-0000-000000000001")
	fw, _ := w.CreateFormFile("file", "test.jpg")
	fw.Write([]byte("not real"))
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/media/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(authCtx())
	rec := httptest.NewRecorder()
	h.Upload(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAtlasMediaUpload_PathTraversal_SanitizesFilename(t *testing.T) {
	baseDir := t.TempDir()
	h := handler.NewAtlasMediaHandler(&mockExerciseServiceForMedia{
		createMediaFn: func(ctx context.Context, userID, exerciseID, fileName, filePath, mimeType string, fileSize int64) (*models.ExerciseMedia, error) {
			assert.Equal(t, "malicious.jpg", fileName)
			return &models.ExerciseMedia{ID: "media-1", ExerciseID: exerciseID, FileName: fileName, MimeType: mimeType, FileSize: fileSize}, nil
		},
	}, baseDir)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("purpose", "EXERCISE_MEDIA")
	w.WriteField("exerciseId", "00000000-0000-0000-0000-000000000001")

	jpegHeader := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01}
	fw, _ := w.CreateFormFile("file", "../../../etc/malicious.jpg")
	fw.Write(jpegHeader)
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/media/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(authCtx())
	rec := httptest.NewRecorder()
	h.Upload(rec, req)

	if rec.Code == http.StatusCreated {
		assert.NotContains(t, rec.Body.String(), "etc")
	}
}

func createMinJPEG() []byte {
	return []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01}
}

func TestAtlasMediaUpload_AcceptsValidJPEG(t *testing.T) {
	baseDir := t.TempDir()
	h := handler.NewAtlasMediaHandler(&mockExerciseServiceForMedia{
		createMediaFn: func(ctx context.Context, userID, exerciseID, fileName, filePath, mimeType string, fileSize int64) (*models.ExerciseMedia, error) {
			return &models.ExerciseMedia{ID: "m1", ExerciseID: exerciseID, FileName: fileName, MimeType: mimeType, FileSize: fileSize}, nil
		},
	}, baseDir)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("purpose", "EXERCISE_MEDIA")
	w.WriteField("exerciseId", "00000000-0000-0000-0000-000000000001")
	fw, _ := w.CreateFormFile("file", "photo.jpg")
	fw.Write(createMinJPEG())
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/media/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(authCtx())
	rec := httptest.NewRecorder()
	h.Upload(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestAtlasMediaDownload_RequiresAuth(t *testing.T) {
	h := handler.NewAtlasMediaHandler(&mockExerciseServiceForMedia{}, t.TempDir())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/media/123", nil)
	rec := httptest.NewRecorder()
	h.Download(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAtlasMediaDownload_NotFound(t *testing.T) {
	h := handler.NewAtlasMediaHandler(&mockExerciseServiceForMedia{
		getMediaRecordByIDFn: func(ctx context.Context, userID, id string) (*models.ExerciseMediaRecord, error) {
			return nil, atlasSvc.ErrExerciseNotFound
		},
	}, t.TempDir())

	r := chi.NewRouter()
	r.Use(authMiddleware)
	r.Get("/api/v1/media/{id}", h.Download)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/media/123", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAtlasMediaDownload_ReturnsFileWithContentType(t *testing.T) {
	baseDir := t.TempDir()
	storageDir := filepath.Join(baseDir, "exercise", "ex-1")
	os.MkdirAll(storageDir, 0755)
	filePath := filepath.Join(storageDir, "test-photo.jpg")
	os.WriteFile(filePath, createMinJPEG(), 0644)

	h := handler.NewAtlasMediaHandler(&mockExerciseServiceForMedia{
		getMediaRecordByIDFn: func(ctx context.Context, userID, id string) (*models.ExerciseMediaRecord, error) {
			return &models.ExerciseMediaRecord{
				ID:       id,
				UserID:   userID,
				FilePath: filePath,
				MimeType: "image/jpeg",
				FileName: "photo.jpg",
			}, nil
		},
	}, baseDir)

	r := chi.NewRouter()
	r.Use(authMiddleware)
	r.Get("/api/v1/media/{id}", h.Download)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/media/media-1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "image/jpeg", rec.Header().Get("Content-Type"))
	body, _ := io.ReadAll(rec.Body)
	assert.NotEmpty(t, body)
}

func TestAtlasMediaDelete_RequiresAuth(t *testing.T) {
	h := handler.NewAtlasMediaHandler(&mockExerciseServiceForMedia{}, t.TempDir())
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/media/123", nil)
	rec := httptest.NewRecorder()
	h.Delete(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAtlasMediaDelete_NotFound(t *testing.T) {
	h := handler.NewAtlasMediaHandler(&mockExerciseServiceForMedia{
		deleteMediaFn: func(ctx context.Context, userID, id string) (*models.ExerciseMediaRecord, error) {
			return nil, atlasSvc.ErrExerciseNotFound
		},
	}, t.TempDir())

	r := chi.NewRouter()
	r.Use(authMiddleware)
	r.Delete("/api/v1/media/{id}", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/media/123", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAtlasMediaDelete_DeletesFileAndReturns204(t *testing.T) {
	baseDir := t.TempDir()
	storageDir := filepath.Join(baseDir, "exercise", "ex-2")
	os.MkdirAll(storageDir, 0755)
	filePath := filepath.Join(storageDir, "test-photo.jpg")
	os.WriteFile(filePath, createMinJPEG(), 0644)

	h := handler.NewAtlasMediaHandler(&mockExerciseServiceForMedia{
		deleteMediaFn: func(ctx context.Context, userID, id string) (*models.ExerciseMediaRecord, error) {
			return &models.ExerciseMediaRecord{
				ID:       id,
				UserID:   userID,
				FilePath: filePath,
				MimeType: "image/jpeg",
				FileName: "photo.jpg",
			}, nil
		},
	}, baseDir)

	r := chi.NewRouter()
	r.Use(authMiddleware)
	r.Delete("/api/v1/media/{id}", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/media/media-1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	_, err := os.Stat(filePath)
	assert.True(t, os.IsNotExist(err))
}

func authCtx() context.Context {
	return middleware.ContextWithAtlasUserID(context.Background(), "test-uid-0000-0000-0000-000000000001")
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := middleware.ContextWithAtlasUserID(r.Context(), "test-uid-0000-0000-0000-000000000001")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}