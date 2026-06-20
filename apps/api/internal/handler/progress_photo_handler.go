// FILE: apps/api/internal/handler/progress_photo_handler.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide Atlas progress photo REST handler for WAVE-04 upload, download, and delete.
//   SCOPE: POST /api/v1/progress-photos/upload, GET /api/v1/progress-photos/{id}, DELETE /api/v1/progress-photos/{id}. File validation (MIME, size), UUID storage, physical file delete.
//   DEPENDS: apps/api/internal/atlas/repository/postgres.ProgressPhotoRepository, apps/api/internal/atlas/repository/postgres.BodyCheckInRepository, apps/api/internal/appconfig.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ProgressPhotoHandler - Struct holding repo and MediaConfig dependencies.
//   NewProgressPhotoHandler - Creates a new ProgressPhotoHandler.
//   Upload - Handles POST multipart with MIME/size validation, UUID storage, check-in photo limit check.
//   Download - Handles GET with correct Content-Type.
//   Delete - Handles DELETE with DB record + physical file removal.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added progress photo REST handler for WAVE-04.
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

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/middleware"
	atlasRepo "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/libs/go/logger"
)

var allowedPhotoMIMETypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

const (
	maxPhotoSize    = 25 * 1024 * 1024
	maxPhotoBody    = 30 * 1024 * 1024
	maxPhotosPerCheckIn = 10
)

type ProgressPhotoHandler struct {
	photoRepo  atlasRepo.ProgressPhotoRepository
	checkInRepo atlasRepo.BodyCheckInRepository
	basePath   string
}

func NewProgressPhotoHandler(photoRepo atlasRepo.ProgressPhotoRepository, checkInRepo atlasRepo.BodyCheckInRepository, basePath string) *ProgressPhotoHandler {
	return &ProgressPhotoHandler{
		photoRepo:   photoRepo,
		checkInRepo: checkInRepo,
		basePath:    basePath,
	}
}

type photoUploadResponse struct {
	Photo *photoMetadata  `json:"photo,omitempty"`
	Error *photoError     `json:"error,omitempty"`
}

type photoMetadata struct {
	ID               string  `json:"id"`
	CheckInID        string  `json:"checkInId"`
	OriginalFileName string  `json:"originalFileName"`
	MimeType         string  `json:"mimeType"`
	SizeBytes        int64   `json:"sizeBytes"`
	Angle            *string `json:"angle"`
	Label            *string `json:"label"`
}

type photoError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (h *ProgressPhotoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writePhotoError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxPhotoBody)
	if err := r.ParseMultipartForm(maxPhotoBody); err != nil {
		log.Warn("[ProgressPhoto][upload][BLOCK_PHOTO_TOO_LARGE] body too large", zap.Error(err))
		writePhotoError(w, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", "file size exceeds maximum allowed")
		return
	}
	defer func() { _ = r.MultipartForm.RemoveAll() }()

	checkInID := r.FormValue("checkInId")
	if checkInID == "" {
		writePhotoError(w, http.StatusBadRequest, "MISSING_CHECKIN_ID", "checkInId is required")
		return
	}

	if _, err := h.checkInRepo.GetByID(r.Context(), userID, checkInID); err != nil {
		writePhotoError(w, http.StatusNotFound, "NOT_FOUND", "check-in not found")
		return
	}

	count, err := h.photoRepo.CountByCheckIn(r.Context(), checkInID)
	if err != nil {
		log.Error("[ProgressPhoto][upload][BLOCK_PHOTO_COUNT] cannot count photos", zap.Error(err))
		writePhotoError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "storage error")
		return
	}
	if count >= maxPhotosPerCheckIn {
		writePhotoError(w, http.StatusBadRequest, "PHOTO_LIMIT_EXCEEDED", fmt.Sprintf("maximum %d photos per check-in", maxPhotosPerCheckIn))
		return
	}

	angle := r.FormValue("angle")
	if angle != "" && !isValidPhotoAngle(angle) {
		writePhotoError(w, http.StatusBadRequest, "INVALID_ANGLE", "angle must be one of: FRONT, SIDE, BACK, CUSTOM")
		return
	}
	label := r.FormValue("label")
	notes := r.FormValue("notes")

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Warn("[ProgressPhoto][upload][BLOCK_PHOTO_NO_FILE] no file in request", zap.Error(err))
		writePhotoError(w, http.StatusBadRequest, "NO_FILE", "file is required")
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, _ := io.ReadFull(file, buf)
	if n == 0 {
		log.Warn("[ProgressPhoto][upload][BLOCK_PHOTO_READ_ERROR] cannot read file header")
		writePhotoError(w, http.StatusBadRequest, "INVALID_FILE", "cannot read file")
		return
	}

	mimeType := http.DetectContentType(buf[:n])
	ext, ok := allowedPhotoMIMETypes[mimeType]
	if !ok {
		log.Warn("[ProgressPhoto][upload][BLOCK_PHOTO_INVALID_TYPE] unsupported MIME", zap.String("detected", mimeType))
		writePhotoError(w, http.StatusBadRequest, "INVALID_FILE_TYPE", fmt.Sprintf("unsupported file type: %s. Supported: JPEG, PNG, WEBP", mimeType))
		return
	}

	if header.Size > maxPhotoSize {
		log.Warn("[ProgressPhoto][upload][BLOCK_PHOTO_SIZE_EXCEEDED] file exceeds size limit",
			zap.Int64("size", header.Size),
			zap.Int("max", maxPhotoSize),
		)
		writePhotoError(w, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", "file size exceeds 25MB limit")
		return
	}

	storageDir := filepath.Join(h.basePath, "progress-photos", checkInID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Error("[ProgressPhoto][upload][BLOCK_PHOTO_MKDIR] cannot create storage dir", zap.Error(err))
		writePhotoError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "storage error")
		return
	}

	uuid := newPhotoUUID()
	storageName := uuid + ext
	storagePath := filepath.Join(storageDir, storageName)

	dst, err := os.Create(storagePath)
	if err != nil {
		log.Error("[ProgressPhoto][upload][BLOCK_PHOTO_WRITE] cannot create file", zap.Error(err))
		writePhotoError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "storage error")
		return
	}
	defer dst.Close()

	if _, err := dst.Write(buf[:n]); err != nil {
		log.Error("[ProgressPhoto][upload][BLOCK_PHOTO_WRITE] cannot write file header", zap.Error(err))
		os.Remove(storagePath)
		writePhotoError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "storage error")
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		log.Error("[ProgressPhoto][upload][BLOCK_PHOTO_WRITE] cannot write file content", zap.Error(err))
		os.Remove(storagePath)
		writePhotoError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "storage error")
		return
	}

	anglePtr := stringPtr(angle)
	labelPtr := stringPtr(label)
	notesPtr := stringPtr(notes)
	originalName := sanitizeFileName(header.Filename)

	record, err := h.photoRepo.Create(r.Context(), checkInID, storagePath, originalName, mimeType, header.Size, anglePtr, labelPtr, notesPtr)
	if err != nil {
		log.Error("[ProgressPhoto][upload][BLOCK_PHOTO_DB] cannot save photo record", zap.Error(err))
		os.Remove(storagePath)
		writePhotoError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to save photo")
		return
	}

	log.Info("[ProgressPhoto][upload] upload success",
		zap.String("photo_id", record.ID),
		zap.String("checkin_id", checkInID),
		zap.String("mime", mimeType),
		zap.Int64("size", header.Size),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(photoUploadResponse{
		Photo: &photoMetadata{
			ID:               record.ID,
			CheckInID:        record.CheckInID,
			OriginalFileName: record.OriginalFileName,
			MimeType:         record.MimeType,
			SizeBytes:        record.SizeBytes,
			Angle:            record.Angle,
			Label:            record.Label,
		},
	})
}

func (h *ProgressPhotoHandler) Download(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writePhotoError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	photoID := chi.URLParam(r, "id")
	if photoID == "" {
		writePhotoError(w, http.StatusBadRequest, "MISSING_ID", "photo id is required")
		return
	}

	photo, err := h.photoRepo.GetByID(r.Context(), userID, photoID)
	if err != nil {
		log.Warn("[ProgressPhoto][download][BLOCK_PHOTO_NOT_FOUND] photo not found", zap.String("id", photoID))
		writePhotoError(w, http.StatusNotFound, "NOT_FOUND", "photo not found")
		return
	}

	file, err := os.Open(photo.FilePath)
	if err != nil {
		log.Error("[ProgressPhoto][download][BLOCK_PHOTO_READ] cannot open file", zap.String("path", photo.FilePath), zap.Error(err))
		writePhotoError(w, http.StatusNotFound, "NOT_FOUND", "file not found")
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Error("[ProgressPhoto][download][BLOCK_PHOTO_STAT] cannot stat file", zap.Error(err))
		writePhotoError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "storage error")
		return
	}

	w.Header().Set("Content-Type", photo.MimeType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, photo.OriginalFileName))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)

	log.Info("[ProgressPhoto][download] download success",
		zap.String("photo_id", photoID),
		zap.String("mime", photo.MimeType),
	)
}

func (h *ProgressPhotoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writePhotoError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	photoID := chi.URLParam(r, "id")
	if photoID == "" {
		writePhotoError(w, http.StatusBadRequest, "MISSING_ID", "photo id is required")
		return
	}

	photo, err := h.photoRepo.Delete(r.Context(), userID, photoID)
	if err != nil {
		log.Warn("[ProgressPhoto][delete][BLOCK_PHOTO_NOT_FOUND] photo not found", zap.String("id", photoID))
		writePhotoError(w, http.StatusNotFound, "NOT_FOUND", "photo not found")
		return
	}

	if err := os.Remove(photo.FilePath); err != nil && !os.IsNotExist(err) {
		log.Error("[ProgressPhoto][delete][BLOCK_PHOTO_FILE_DELETE] file deletion failed, DB record removed",
			zap.String("path", photo.FilePath),
			zap.Error(err),
		)
	} else {
		log.Info("[ProgressPhoto][delete] file deleted", zap.String("path", photo.FilePath))
	}

	dir := filepath.Dir(photo.FilePath)
	removeEmptyPhotoDir(dir)

	log.Info("[ProgressPhoto][delete] delete success", zap.String("photo_id", photoID))
	w.WriteHeader(http.StatusNoContent)
}

func writePhotoError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(photoUploadResponse{
		Error: &photoError{Code: code, Message: message},
	})
}

func isValidPhotoAngle(angle string) bool {
	switch angle {
	case "FRONT", "SIDE", "BACK", "CUSTOM":
		return true
	}
	return false
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func removeEmptyPhotoDir(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}
	if len(entries) == 0 {
		os.Remove(path)
	}
}

func newPhotoUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}