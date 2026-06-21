// FILE: apps/api/internal/handler/ai_export_handler.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide Atlas AI export REST handler for WAVE-07 generate, download, and user profile.
//   SCOPE: POST /api/ai-export/generate, GET /api/ai-export/download, GET /api/user-profile. Date validation, file serving, prompt generation.
//   DEPENDS: apps/api/internal/atlas/service.AiExportService, apps/api/internal/atlas/service.UserProfileService, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AiExportHandler - Struct holding AiExportService, UserProfileService, and config dependencies.
//   NewAiExportHandler - Creates a new AiExportHandler.
//   GenerateExport - Handles POST with date range, section toggles, optional comment; generates ZIP + prompt.
//   DownloadExport - Handles GET with exportId query param; streams ZIP file.
//   GetUserProfile - Handles GET; returns user profile for AI context.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added AI export REST handler for WAVE-07.
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

type AiExportHandler struct {
	aiExportService    atlasService.AiExportService
	userProfileService atlasService.UserProfileService
	exportBasePath     string
	maxRangeDays       int
	maxExportSize     int64
}

type generateRequest struct {
	DateRangeStart      string  `json:"dateRangeStart"`
	DateRangeEnd        string  `json:"dateRangeEnd"`
	IncludePhotos       *bool   `json:"includePhotos"`
	IncludeNutrition    *bool   `json:"includeNutrition"`
	IncludeCardio       *bool   `json:"includeCardio"`
	IncludeMeasurements *bool   `json:"includeMeasurements"`
	UserComment         *string `json:"userComment"`
}

type generateResponse struct {
	Export *exportResponse `json:"export,omitempty"`
	Error  *apiError       `json:"error,omitempty"`
}

type exportResponse struct {
	ID                  string  `json:"id"`
	DateRangeStart      string  `json:"dateRangeStart"`
	DateRangeEnd        string  `json:"dateRangeEnd"`
	IncludePhotos       bool    `json:"includePhotos"`
	IncludeNutrition    bool    `json:"includeNutrition"`
	IncludeCardio       bool    `json:"includeCardio"`
	IncludeMeasurements bool    `json:"includeMeasurements"`
	UserComment         *string `json:"userComment"`
	GeneratedPrompt     string  `json:"generatedPrompt"`
	ExportFilePath      *string `json:"exportFilePath"`
	CreatedAt           string  `json:"createdAt"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewAiExportHandler(aiExportService atlasService.AiExportService, userProfileService atlasService.UserProfileService, exportBasePath string, maxRangeDays int, maxExportSize int64) *AiExportHandler {
	return &AiExportHandler{
		aiExportService:    aiExportService,
		userProfileService: userProfileService,
		exportBasePath:     exportBasePath,
		maxRangeDays:       maxRangeDays,
		maxExportSize:      maxExportSize,
	}
}

func (h *AiExportHandler) GenerateExport(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writeAiExportError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAiExportError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	startDate, err := models.ParseDate(req.DateRangeStart)
	if err != nil {
		writeAiExportError(w, http.StatusBadRequest, "INVALID_DATE_RANGE", "invalid dateRangeStart")
		return
	}
	endDate, err := models.ParseDate(req.DateRangeEnd)
	if err != nil {
		writeAiExportError(w, http.StatusBadRequest, "INVALID_DATE_RANGE", "invalid dateRangeEnd")
		return
	}

	input := models.CreateAiExportInput{
		DateRangeStart:      startDate,
		DateRangeEnd:        endDate,
		IncludePhotos:       req.IncludePhotos,
		IncludeNutrition:    req.IncludeNutrition,
		IncludeCardio:       req.IncludeCardio,
		IncludeMeasurements: req.IncludeMeasurements,
		UserComment:         req.UserComment,
	}

	export, prompt, err := h.aiExportService.Generate(r.Context(), userID, input, h.maxRangeDays, h.maxExportSize, h.exportBasePath)
	if err != nil {
		log.Error("[AiExport][generate][BLOCK_EXPORT_FAILURE] generate failed", zap.Error(err))
		writeAiExportError(w, http.StatusInternalServerError, "EXPORT_GENERATION_FAILED", "export generation failed")
		return
	}

	log.Info("[AiExport][generate][BLOCK_EXPORT_SUCCESS] export generated",
		zap.String("export_id", export.ID),
	)

	resp := generateResponse{
		Export: &exportResponse{
			ID:                  export.ID,
			DateRangeStart:      export.DateRangeStart.String(),
			DateRangeEnd:        export.DateRangeEnd.String(),
			IncludePhotos:       export.IncludePhotos,
			IncludeNutrition:    export.IncludeNutrition,
			IncludeCardio:       export.IncludeCardio,
			IncludeMeasurements: export.IncludeMeasurements,
			UserComment:         export.UserComment,
			GeneratedPrompt:     prompt,
			ExportFilePath:      export.ExportFilePath,
			CreatedAt:           export.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *AiExportHandler) DownloadExport(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writeAiExportError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	exportID := r.URL.Query().Get("exportId")
	if exportID == "" {
		log.Warn("[AiExport][download][BLOCK_EXPORT_NOT_FOUND] missing exportId")
		writeAiExportError(w, http.StatusBadRequest, "MISSING_EXPORT_ID", "exportId is required")
		return
	}

	export, err := h.aiExportService.GetByID(r.Context(), userID, exportID)
	if err != nil {
		log.Warn("[AiExport][download][BLOCK_EXPORT_NOT_FOUND] export not found", zap.String("export_id", exportID))
		writeAiExportError(w, http.StatusNotFound, "EXPORT_NOT_FOUND", "export not found")
		return
	}

	if export.ExportFilePath == nil || *export.ExportFilePath == "" {
		log.Warn("[AiExport][download][BLOCK_EXPORT_FILE_MISSING] file path is nil", zap.String("export_id", exportID))
		writeAiExportError(w, http.StatusNotFound, "EXPORT_FILE_NOT_FOUND", "export file not found")
		return
	}

	file, err := os.Open(*export.ExportFilePath)
	if err != nil {
		log.Error("[AiExport][download][BLOCK_EXPORT_FILE_MISSING] cannot open file", zap.String("path", *export.ExportFilePath), zap.Error(err))
		writeAiExportError(w, http.StatusNotFound, "EXPORT_FILE_NOT_FOUND", "export file not found")
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Error("[AiExport][download] cannot stat file", zap.Error(err))
		writeAiExportError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "storage error")
		return
	}

	shortID := exportID
	if len(exportID) > 8 {
		shortID = exportID[:8]
	}
	filename := fmt.Sprintf("ai-export-%s-%s-%s.zip", export.DateRangeStart.String(), export.DateRangeEnd.String(), shortID)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)

	log.Info("[AiExport][download] download success",
		zap.String("export_id", exportID),
	)
}

func (h *AiExportHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writeAiExportError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	profile, err := h.userProfileService.Get(r.Context(), userID)
	if err != nil {
		log.Warn("[AiExport][user_profile] profile not found", zap.Error(err))
		writeAiExportError(w, http.StatusNotFound, "PROFILE_NOT_FOUND", "user profile not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

func writeAiExportError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(generateResponse{
		Error: &apiError{Code: code, Message: message},
	})
}
