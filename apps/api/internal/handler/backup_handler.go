// FILE: apps/api/internal/handler/backup_handler.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide Atlas backup REST handler for WAVE-09 export/import.
//   SCOPE: POST /api/backup/export, GET /api/backup/download, POST /api/backup/import/validate, POST /api/backup/import/confirm. ZIP generation, file streaming, multipart upload.
//   DEPENDS: apps/api/internal/atlas/service.BackupExportService, apps/api/internal/atlas/service.BackupImportService, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-09.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   BackupHandler - Struct holding BackupExportService, BackupImportService, and config dependencies.
//   NewBackupHandler - Creates a new BackupHandler.
//   GenerateExport - Handles POST with optional includeMedia flag; generates backup ZIP.
//   DownloadExport - Handles GET with downloadId query param; streams ZIP file.
//   ImportValidate - Handles POST multipart form file upload; validates import ZIP.
//   ImportConfirm - Handles POST with validationId; confirms import.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added backup REST handler for WAVE-09.
// END_CHANGE_SUMMARY

package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
	"monorepo-template/libs/go/logger"
)

type BackupHandler struct {
	backupExportService atlasService.BackupExportService
	backupImportService atlasService.BackupImportService
	exportBasePath      string
	maxExportSize       int64
	maxImportSize       int64
}

func NewBackupHandler(
	backupExportService atlasService.BackupExportService,
	backupImportService atlasService.BackupImportService,
	exportBasePath string,
	maxExportSize int64,
	maxImportSize int64,
) *BackupHandler {
	return &BackupHandler{
		backupExportService: backupExportService,
		backupImportService: backupImportService,
		exportBasePath:      exportBasePath,
		maxExportSize:       maxExportSize,
		maxImportSize:       maxImportSize,
	}
}

type exportRequest struct {
	IncludeMedia *bool `json:"includeMedia"`
}

type backupExportResponse struct {
	DownloadID string `json:"downloadId"`
	SizeBytes  int64  `json:"sizeBytes"`
	Timestamp  string `json:"timestamp"`
}

type importValidateResponse struct {
	ValidationID string                   `json:"validationId"`
	Summary      *models.BackupImportSummary `json:"summary"`
}

type importConfirmRequest struct {
	ValidationID string `json:"validationId"`
}

type importConfirmResponse struct {
	Status       string         `json:"status"`
	EntityCounts map[string]int `json:"entityCounts"`
	MediaCount   int            `json:"mediaCount"`
}

func (h *BackupHandler) GenerateExport(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writeBackupError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	var req exportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBackupError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	includeMedia := false
	if req.IncludeMedia != nil {
		includeMedia = *req.IncludeMedia
	}

	result, err := h.backupExportService.Generate(r.Context(), userID, includeMedia, h.maxExportSize, h.exportBasePath)
	if err != nil {
		log.Error("[Backup][export][BLOCK_EXPORT_FAILURE] export generation failed", zap.Error(err))
		writeBackupError(w, http.StatusInternalServerError, "EXPORT_FAILED", "export generation failed")
		return
	}

	log.Info("[Backup][export][BLOCK_EXPORT_SUCCESS] export generated",
		zap.String("download_id", result.DownloadID),
	)

	resp := backupExportResponse{
		DownloadID: result.DownloadID,
		SizeBytes:  result.SizeBytes,
		Timestamp:  result.Timestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *BackupHandler) DownloadExport(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writeBackupError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	downloadID := r.URL.Query().Get("downloadId")
	if downloadID == "" {
		log.Warn("[Backup][download][BLOCK_DOWNLOAD_MISSING] missing downloadId")
		writeBackupError(w, http.StatusBadRequest, "MISSING_DOWNLOAD_ID", "downloadId is required")
		return
	}

	filePath, err := h.backupExportService.GetDownloadPath(r.Context(), userID, downloadID)
	if err != nil {
		log.Warn("[Backup][download][BLOCK_DOWNLOAD_NOT_FOUND] download not found", zap.String("download_id", downloadID))
		writeBackupError(w, http.StatusNotFound, "DOWNLOAD_NOT_FOUND", "download not found")
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Error("[Backup][download][BLOCK_FILE_MISSING] cannot open file", zap.String("path", filePath), zap.Error(err))
		writeBackupError(w, http.StatusNotFound, "FILE_NOT_FOUND", "export file not found")
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Error("[Backup][download] cannot stat file", zap.Error(err))
		writeBackupError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "storage error")
		return
	}

	filename := fmt.Sprintf("atlas-backup-%s.zip", downloadID)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)

	log.Info("[Backup][download] download success",
		zap.String("download_id", downloadID),
	)
}

func (h *BackupHandler) ImportValidate(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writeBackupError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	if err := r.ParseMultipartForm(h.maxImportSize + 1024); err != nil {
		log.Warn("[Backup][import][BLOCK_PARSE_FAILURE] failed to parse multipart form", zap.Error(err))
		writeBackupError(w, http.StatusBadRequest, "INVALID_UPLOAD", "failed to parse upload")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		log.Warn("[Backup][import][BLOCK_FILE_MISSING] missing file field")
		writeBackupError(w, http.StatusBadRequest, "MISSING_FILE", "file field is required")
		return
	}
	defer file.Close()

	limitedReader := io.LimitReader(file, h.maxImportSize+1)
	zipData, err := io.ReadAll(limitedReader)
	if err != nil {
		log.Error("[Backup][import][BLOCK_READ_FAILURE] failed to read uploaded file", zap.Error(err))
		writeBackupError(w, http.StatusInternalServerError, "READ_FAILED", "failed to read uploaded file")
		return
	}

	if int64(len(zipData)) > h.maxImportSize {
		writeBackupError(w, http.StatusBadRequest, "FILE_TOO_LARGE", fmt.Sprintf("file exceeds maximum import size of %d bytes", h.maxImportSize))
		return
	}

	validationID, summary, err := h.backupImportService.Validate(r.Context(), userID, zipData)
	if err != nil {
		log.Error("[Backup][import][BLOCK_VALIDATION_FAILURE] validation failed", zap.Error(err))
		writeBackupError(w, http.StatusInternalServerError, "VALIDATION_FAILED", "import validation failed")
		return
	}

	log.Info("[Backup][import][BLOCK_VALIDATION_SUCCESS] import validated",
		zap.String("validation_id", validationID),
	)

	resp := importValidateResponse{
		ValidationID: validationID,
		Summary:      summary,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *BackupHandler) ImportConfirm(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writeBackupError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	var req importConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBackupError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if req.ValidationID == "" {
		writeBackupError(w, http.StatusBadRequest, "MISSING_VALIDATION_ID", "validationId is required")
		return
	}

	result, err := h.backupImportService.Confirm(r.Context(), userID, req.ValidationID)
	if err != nil {
		log.Error("[Backup][import][BLOCK_CONFIRM_FAILURE] confirm failed", zap.Error(err))
		writeBackupError(w, http.StatusInternalServerError, "CONFIRM_FAILED", "import confirm failed")
		return
	}

	log.Info("[Backup][import][BLOCK_CONFIRM_SUCCESS] import confirmed",
		zap.String("validation_id", req.ValidationID),
	)

	resp := importConfirmResponse{
		Status:       result.Status,
		EntityCounts: result.EntityCounts,
		MediaCount:   result.MediaCount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func writeBackupError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": apiError{Code: code, Message: message},
	})
}