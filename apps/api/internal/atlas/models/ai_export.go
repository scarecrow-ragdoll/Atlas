// FILE: apps/api/internal/atlas/models/ai_export.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define WAVE-07 AiExport domain models, input types, result/error types.
//   SCOPE: AiExportRecord (DB), AiExport (public), CreateAiExportInput, typed error codes, AiExportFromRecord converter.
//   DEPENDS: apps/api/internal/atlas/models/date.go.
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AiExportRecord - Internal DB model matching ai_exports table.
//   AiExport - Public AI export model.
//   CreateAiExportInput - Input for creating an AI export.
//   AiExportResult - Union-like result for single AI export operations.
//   AiExportsResult - Union-like result for listing AI exports.
//   AiExportErrorCode - Enum for AI export error codes.
//   AiExportFromRecord - Converts DB record to public model.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added WAVE-07 AI export domain models.
// END_CHANGE_SUMMARY

package models

type AiExportRecord struct {
	ID                 string
	UserID             string
	DateRangeStart     Date
	DateRangeEnd       Date
	IncludePhotos      bool
	IncludeNutrition   bool
	IncludeCardio      bool
	IncludeMeasurements bool
	UserComment        *string
	GeneratedPrompt    string
	ExportFilePath     *string
	CreatedAt          string
	UpdatedAt          string
}

type AiExport struct {
	ID                 string  `json:"id"`
	UserID             string  `json:"userId"`
	DateRangeStart     Date    `json:"dateRangeStart"`
	DateRangeEnd       Date    `json:"dateRangeEnd"`
	IncludePhotos      bool    `json:"includePhotos"`
	IncludeNutrition   bool    `json:"includeNutrition"`
	IncludeCardio      bool    `json:"includeCardio"`
	IncludeMeasurements bool   `json:"includeMeasurements"`
	UserComment        *string `json:"userComment"`
	GeneratedPrompt    string  `json:"generatedPrompt"`
	ExportFilePath     *string `json:"exportFilePath"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

type CreateAiExportInput struct {
	DateRangeStart     Date    `json:"dateRangeStart"`
	DateRangeEnd       Date    `json:"dateRangeEnd"`
	IncludePhotos      *bool   `json:"includePhotos"`
	IncludeNutrition   *bool   `json:"includeNutrition"`
	IncludeCardio      *bool   `json:"includeCardio"`
	IncludeMeasurements *bool  `json:"includeMeasurements"`
	UserComment        *string `json:"userComment"`
}

type AiExportResult struct {
	Export        *AiExport                  `json:"export"`
	ValidationErr *AiExportValidationErr      `json:"validationError"`
	NotFoundErr   *AiExportNotFoundErr        `json:"notFoundError"`
	AuthErr       *AiExportAuthErr            `json:"authError"`
}

type AiExportsResult struct {
	Exports       []AiExport                 `json:"exports"`
	ValidationErr *AiExportValidationErr      `json:"validationError"`
	AuthErr       *AiExportAuthErr            `json:"authError"`
}

type AiExportValidationErr struct {
	Message string              `json:"message"`
	Code    AiExportErrorCode   `json:"code"`
}

func (e *AiExportValidationErr) Error() string {
	if e == nil || e.Message == "" {
		return "ai export validation error"
	}
	return e.Message
}

type AiExportNotFoundErr struct {
	Message string              `json:"message"`
	Code    AiExportErrorCode   `json:"code"`
}

func (e *AiExportNotFoundErr) Error() string {
	if e == nil || e.Message == "" {
		return "ai export not found"
	}
	return e.Message
}

type AiExportAuthErr struct {
	Message string              `json:"message"`
	Code    AiExportErrorCode   `json:"code"`
}

func (e *AiExportAuthErr) Error() string {
	if e == nil || e.Message == "" {
		return "ai export auth error"
	}
	return e.Message
}

type AiExportErrorCode string

const (
	AiExportErrorValidation AiExportErrorCode = "VALIDATION_ERROR"
	AiExportErrorNotFound   AiExportErrorCode = "NOT_FOUND"
	AiExportErrorAuth       AiExportErrorCode = "AUTH_ERROR"
	AiExportErrorInternal   AiExportErrorCode = "INTERNAL_ERROR"
)

func AiExportFromRecord(r *AiExportRecord) *AiExport {
	if r == nil {
		return nil
	}
	return &AiExport{
		ID:                 r.ID,
		UserID:             r.UserID,
		DateRangeStart:     r.DateRangeStart,
		DateRangeEnd:       r.DateRangeEnd,
		IncludePhotos:      r.IncludePhotos,
		IncludeNutrition:   r.IncludeNutrition,
		IncludeCardio:      r.IncludeCardio,
		IncludeMeasurements: r.IncludeMeasurements,
		UserComment:        r.UserComment,
		GeneratedPrompt:    r.GeneratedPrompt,
		ExportFilePath:     r.ExportFilePath,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
	}
}