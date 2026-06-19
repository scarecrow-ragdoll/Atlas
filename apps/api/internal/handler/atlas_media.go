// FILE: apps/api/internal/handler/atlas_media.go
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Provide Atlas media REST handlers for exercise media upload, download, and delete. Extended in WAVE-02 from WAVE-01 scaffold.
//   SCOPE: POST /api/v1/media/upload (multipart with purpose=EXERCISE_MEDIA + exerciseId), GET /api/v1/media/{id}, DELETE /api/v1/media/{id}. File validation via http.DetectContentType(), size limits, UUID-based storage.
//   DEPENDS: apps/api/internal/atlas/service.ExerciseService, apps/api/internal/appconfig.MediaConfig.
//   LINKS: M-API / V-M-API / WAVE-02.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AtlasMediaHandler - Struct holding ExerciseService and MediaConfig dependencies.
//   NewAtlasMediaHandler - Creates a new AtlasMediaHandler.
//   Upload - Handles POST /api/v1/media/upload with multipart form, MIME validation, size limits, UUID storage.
//   Download - Handles GET /api/v1/media/{id} with correct Content-Type.
//   Delete - Handles DELETE /api/v1/media/{id} with DB + physical file removal.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Extended from scaffold to full implementation for WAVE-02.
// END_CHANGE_SUMMARY

package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/service"
	"monorepo-template/libs/go/logger"
)

var allowedMIMETypes = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/webp":      ".webp",
	"video/mp4":       ".mp4",
	"video/quicktime": ".mov",
	"video/webm":      ".webm",
}

const (
	maxImageSize = 25 * 1024 * 1024
	maxVideoSize = 250 * 1024 * 1024
	maxBodySize  = 300 * 1024 * 1024
	detectBytes  = 512
)

type AtlasMediaHandler struct {
	exerciseService service.ExerciseService
	basePath        string
}

func NewAtlasMediaHandler(exerciseService service.ExerciseService, basePath string) *AtlasMediaHandler {
	return &AtlasMediaHandler{
		exerciseService: exerciseService,
		basePath:        basePath,
	}
}

type mediaUploadResponse struct {
	Media *models.ExerciseMedia `json:"media,omitempty"`
	Error *mediaError           `json:"error,omitempty"`
}

type mediaError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (h *AtlasMediaHandler) Upload(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writeMediaError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	if err := r.ParseMultipartForm(maxBodySize); err != nil {
		log.Warn("[Media][upload][BLOCK_MEDIA_FILE_TOO_LARGE] body too large", zap.Error(err))
		writeMediaError(w, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", "file size exceeds maximum allowed")
		return
	}
	defer func() { _ = r.MultipartForm.RemoveAll() }()

	purpose := r.FormValue("purpose")
	if purpose != "EXERCISE_MEDIA" {
		writeMediaError(w, http.StatusBadRequest, "INVALID_PURPOSE", "purpose must be EXERCISE_MEDIA")
		return
	}

	exerciseID := r.FormValue("exerciseId")
	if exerciseID == "" {
		writeMediaError(w, http.StatusBadRequest, "MISSING_EXERCISE_ID", "exerciseId is required")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Warn("[Media][upload][BLOCK_MEDIA_NO_FILE] no file in request", zap.Error(err))
		writeMediaError(w, http.StatusBadRequest, "NO_FILE", "file is required")
		return
	}
	defer file.Close()

	buf := make([]byte, detectBytes)
	n, _ := io.ReadFull(file, buf)
	if n == 0 {
		log.Warn("[Media][upload][BLOCK_MEDIA_READ_ERROR] cannot read file header")
		writeMediaError(w, http.StatusBadRequest, "INVALID_FILE", "cannot read file")
		return
	}

	mimeType := http.DetectContentType(buf[:n])
	ext, ok := allowedMIMETypes[mimeType]
	if !ok {
		log.Warn("[Media][upload][BLOCK_MEDIA_INVALID_TYPE] unsupported MIME", zap.String("detected", mimeType))
		writeMediaError(w, http.StatusBadRequest, "INVALID_FILE_TYPE", fmt.Sprintf("unsupported file type: %s. Supported: JPEG, PNG, WEBP, MP4, MOV, WEBM", mimeType))
		return
	}

	isVideo := strings.HasPrefix(mimeType, "video/")
	maxSize := maxImageSize
	if isVideo {
		maxSize = maxVideoSize
	}
	if header.Size > int64(maxSize) {
		log.Warn("[Media][upload][BLOCK_MEDIA_SIZE_EXCEEDED] file exceeds size limit",
			zap.Int64("size", header.Size),
			zap.Int("max", maxSize),
		)
		limit := "25MB"
		if isVideo {
			limit = "250MB"
		}
		writeMediaError(w, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", fmt.Sprintf("file size exceeds %s limit", limit))
		return
	}

	sanitized := sanitizeFileName(header.Filename)
	storageDir := filepath.Join(h.basePath, "exercise", exerciseID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Error("[Media][upload][BLOCK_MEDIA_MKDIR] cannot create storage dir", zap.Error(err))
		writeMediaError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "storage error")
		return
	}

	uuid := newUUID()
	storageName := uuid + ext
	storagePath := filepath.Join(storageDir, storageName)

	dst, err := os.Create(storagePath)
	if err != nil {
		log.Error("[Media][upload][BLOCK_MEDIA_WRITE] cannot create file", zap.Error(err))
		writeMediaError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "storage error")
		return
	}
	defer dst.Close()

	if _, err := dst.Write(buf[:n]); err != nil {
		log.Error("[Media][upload][BLOCK_MEDIA_WRITE] cannot write file header", zap.Error(err))
		writeMediaError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "storage error")
		return
	}

	if _, err := io.Copy(dst, file); err != nil {
		log.Error("[Media][upload][BLOCK_MEDIA_WRITE] cannot write file content", zap.Error(err))
		writeMediaError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "storage error")
		os.Remove(storagePath)
		return
	}

	media, err := h.exerciseService.CreateMedia(r.Context(), userID, exerciseID, sanitized, storagePath, mimeType, header.Size)
	if err != nil {
		log.Error("[Media][upload][BLOCK_MEDIA_DB] cannot save media record", zap.Error(err))
		os.Remove(storagePath)
		writeMediaError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to save media")
		return
	}

	log.Info("[Media][upload] upload success",
		zap.String("media_id", media.ID),
		zap.String("exercise_id", exerciseID),
		zap.String("mime", mimeType),
		zap.Int64("size", header.Size),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mediaUploadResponse{Media: media})
}

func (h *AtlasMediaHandler) Download(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writeMediaError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	mediaID := chi.URLParam(r, "id")
	if mediaID == "" {
		writeMediaError(w, http.StatusBadRequest, "MISSING_ID", "media id is required")
		return
	}

	media, err := h.exerciseService.GetMediaRecordByID(r.Context(), userID, mediaID)
	if err != nil {
		log.Warn("[Media][download][BLOCK_MEDIA_NOT_FOUND] media not found", zap.String("id", mediaID))
		writeMediaError(w, http.StatusNotFound, "NOT_FOUND", "media not found")
		return
	}

	file, err := os.Open(media.FilePath)
	if err != nil {
		log.Error("[Media][download][BLOCK_MEDIA_READ] cannot open file", zap.String("path", media.FilePath), zap.Error(err))
		writeMediaError(w, http.StatusNotFound, "NOT_FOUND", "file not found")
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Error("[Media][download][BLOCK_MEDIA_STAT] cannot stat file", zap.Error(err))
		writeMediaError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "storage error")
		return
	}

	w.Header().Set("Content-Type", media.MimeType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, media.FileName))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)

	log.Info("[Media][download] download success",
		zap.String("media_id", mediaID),
		zap.String("mime", media.MimeType),
	)
}

func (h *AtlasMediaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writeMediaError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	mediaID := chi.URLParam(r, "id")
	if mediaID == "" {
		writeMediaError(w, http.StatusBadRequest, "MISSING_ID", "media id is required")
		return
	}

	media, err := h.exerciseService.DeleteMedia(r.Context(), userID, mediaID)
	if err != nil {
		log.Warn("[Media][delete][BLOCK_MEDIA_NOT_FOUND] media not found", zap.String("id", mediaID))
		writeMediaError(w, http.StatusNotFound, "NOT_FOUND", "media not found")
		return
	}

	if err := os.Remove(media.FilePath); err != nil && !os.IsNotExist(err) {
		log.Error("[Media][delete][BLOCK_MEDIA_FILE_DELETE] file deletion failed, DB record removed",
			zap.String("path", media.FilePath),
			zap.Error(err),
		)
	} else {
		log.Info("[Media][delete] file deleted", zap.String("path", media.FilePath))
	}

	dir := filepath.Dir(media.FilePath)
	removeEmptyDir(dir)

	log.Info("[Media][delete] delete success", zap.String("media_id", mediaID))
	w.WriteHeader(http.StatusNoContent)
}

func writeMediaError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(mediaUploadResponse{
		Error: &mediaError{Code: code, Message: message},
	})
}

func sanitizeFileName(name string) string {
	name = filepath.Base(name)
	if name == "." || name == "/" {
		return "unnamed"
	}
	return name
}

func removeEmptyDir(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}
	if len(entries) == 0 {
		os.Remove(path)
	}
}

func newUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}