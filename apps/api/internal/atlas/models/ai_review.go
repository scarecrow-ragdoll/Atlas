// FILE: apps/api/internal/atlas/models/ai_review.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define WAVE-08 AiReview domain models, input types, result/error types.
//   SCOPE: AiReviewRecord (DB), AiReview (public), CreateAiReviewInput, UpdateAiReviewInput, typed error codes, AiReviewFromRecord converter.
//   DEPENDS: apps/api/internal/atlas/models/date.go.
//   LINKS: M-API / V-M-API / WAVE-08.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AiReviewRecord - Internal DB model matching ai_reviews table.
//   AiReview - Public AI review model.
//   CreateAiReviewInput - Input for creating an AI review.
//   UpdateAiReviewInput - Input for updating an AI review.
//   AiReviewResult - Union-like result for single AI review operations.
//   AiReviewsResult - Union-like result for listing AI reviews.
//   AiReviewErrorCode - Enum for AI review error codes.
//   AiReviewFromRecord - Converts DB record to public model.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added WAVE-08 AiReview domain models.
// END_CHANGE_SUMMARY

package models

type AiReviewRecord struct {
	ID              string
	UserID          string
	DateRangeStart  Date
	DateRangeEnd    Date
	AiResponseText  string
	UserNotes       *string
	PlannedActions  *string
	CreatedAt       string
	UpdatedAt       string
}

type AiReview struct {
	ID              string  `json:"id"`
	UserID          string  `json:"userId"`
	DateRangeStart  Date    `json:"dateRangeStart"`
	DateRangeEnd    Date    `json:"dateRangeEnd"`
	AiResponseText  string  `json:"aiResponseText"`
	UserNotes       *string `json:"userNotes"`
	PlannedActions  *string `json:"plannedActions"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

type CreateAiReviewInput struct {
	DateRangeStart Date    `json:"dateRangeStart"`
	DateRangeEnd   Date    `json:"dateRangeEnd"`
	AiResponseText string  `json:"aiResponseText"`
	UserNotes      *string `json:"userNotes"`
	PlannedActions *string `json:"plannedActions"`
}

type UpdateAiReviewInput struct {
	DateRangeStart *Date   `json:"dateRangeStart"`
	DateRangeEnd   *Date   `json:"dateRangeEnd"`
	AiResponseText *string `json:"aiResponseText"`
	UserNotes      *string `json:"userNotes"`
	PlannedActions *string `json:"plannedActions"`
}

type AiReviewResult struct {
	Review        *AiReview                `json:"review"`
	ValidationErr *AiReviewValidationErr   `json:"validationError"`
	NotFoundErr   *AiReviewNotFoundErr     `json:"notFoundError"`
	AuthErr       *AiReviewAuthErr         `json:"authError"`
}

type AiReviewsResult struct {
	Reviews       []AiReview               `json:"reviews"`
	ValidationErr *AiReviewValidationErr   `json:"validationError"`
	AuthErr       *AiReviewAuthErr         `json:"authError"`
}

type AiReviewValidationErr struct {
	Message string              `json:"message"`
	Code    AiReviewErrorCode   `json:"code"`
}

func (e *AiReviewValidationErr) Error() string {
	if e == nil || e.Message == "" {
		return "ai review validation error"
	}
	return e.Message
}

type AiReviewNotFoundErr struct {
	Message string              `json:"message"`
	Code    AiReviewErrorCode   `json:"code"`
}

func (e *AiReviewNotFoundErr) Error() string {
	if e == nil || e.Message == "" {
		return "ai review not found"
	}
	return e.Message
}

type AiReviewAuthErr struct {
	Message string              `json:"message"`
	Code    AiReviewErrorCode   `json:"code"`
}

func (e *AiReviewAuthErr) Error() string {
	if e == nil || e.Message == "" {
		return "ai review auth error"
	}
	return e.Message
}

type AiReviewErrorCode string

const (
	AiReviewErrorValidation AiReviewErrorCode = "VALIDATION_ERROR"
	AiReviewErrorNotFound   AiReviewErrorCode = "NOT_FOUND"
	AiReviewErrorAuth       AiReviewErrorCode = "AUTH_ERROR"
	AiReviewErrorInternal   AiReviewErrorCode = "INTERNAL_ERROR"
)

func AiReviewFromRecord(r *AiReviewRecord) *AiReview {
	if r == nil {
		return nil
	}
	return &AiReview{
		ID:              r.ID,
		UserID:          r.UserID,
		DateRangeStart:  r.DateRangeStart,
		DateRangeEnd:    r.DateRangeEnd,
		AiResponseText:  r.AiResponseText,
		UserNotes:       r.UserNotes,
		PlannedActions:  r.PlannedActions,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}