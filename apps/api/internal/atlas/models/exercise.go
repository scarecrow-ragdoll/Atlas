// FILE: apps/api/internal/atlas/models/exercise.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define Exercise and ExerciseMedia domain models for WAVE-02 Exercise Library.
//   SCOPE: Internal ExerciseRecord, public Exercise/ExerciseMedia, result types (ExerciseResult, ArchiveResult), error types (ValidationError, NotFoundError, AuthError), pagination types (PageInfo, ExerciseConnection), and input types (CreateExerciseInput, UpdateExerciseInput). Excludes DailyLog, workout, sets, cardio, body, nutrition, charts, AI export, backup models.
//   DEPENDS: none.
//   LINKS: M-API / V-M-API / WAVE-02.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ExerciseRecord - Internal DB model matching exercises table.
//   Exercise - Public exercise model (no sensitive fields exposed).
//   ExerciseMedia - Public exercise media model.
//   CreateExerciseInput - GraphQL input for creating exercises.
//   UpdateExerciseInput - GraphQL input for updating exercises.
//   ExerciseResult - Union-like result for exercise queries/mutations.
//   ArchiveResult - Union-like result for archive/restore mutations.
//   ExerciseConnection - Paginated exercise list with totalCount and pageInfo.
//   PageInfo - Cursor-based pagination metadata.
//   ExerciseErrorCode - Enum for exercise error codes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added exercise domain models for WAVE-02.
// END_CHANGE_SUMMARY

package models

// ExerciseRecord is the internal DB model mirroring the exercises table.
type ExerciseRecord struct {
	ID            string
	UserID        string
	Name          string
	MuscleGroups  []string
	Description   *string
	PersonalNotes *string
	WorkingWeight *float64
	IsActive      bool
	CreatedAt     string
	UpdatedAt     string
}

type Exercise struct {
	ID            string          `json:"id"`
	UserID        string          `json:"userId"`
	Name          string          `json:"name"`
	MuscleGroups  []string        `json:"muscleGroups"`
	Description   *string         `json:"description"`
	PersonalNotes *string         `json:"personalNotes"`
	WorkingWeight *float64        `json:"workingWeight"`
	IsActive      bool            `json:"isActive"`
	Media         []ExerciseMedia `json:"media"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
}

type ExerciseMedia struct {
	ID         string `json:"id"`
	UserID     string `json:"userId"`
	ExerciseID string `json:"exerciseId"`
	FileName   string `json:"fileName"`
	MimeType   string `json:"mimeType"`
	FileSize   int64  `json:"fileSize"`
	CreatedAt  string `json:"createdAt"`
}

type ExerciseMediaRecord struct {
	ID         string
	UserID     string
	ExerciseID string
	FileName   string
	FilePath   string
	MimeType   string
	FileSize   int64
	CreatedAt  string
}

type CreateExerciseInput struct {
	Name          string    `json:"name"`
	MuscleGroups  []string  `json:"muscleGroups"`
	Description   *string   `json:"description"`
	PersonalNotes *string   `json:"personalNotes"`
	WorkingWeight *float64  `json:"workingWeight"`
}

type UpdateExerciseInput struct {
	Name          *string   `json:"name"`
	MuscleGroups  *[]string `json:"muscleGroups"`
	Description   *string   `json:"description"`
	PersonalNotes *string   `json:"personalNotes"`
	WorkingWeight *float64  `json:"workingWeight"`
}

type PageInfo struct {
	HasNextPage bool    `json:"hasNextPage"`
	EndCursor   *string `json:"endCursor"`
}

type ExerciseConnection struct {
	Items       []Exercise `json:"items"`
	TotalCount  int        `json:"totalCount"`
	PageInfo    PageInfo   `json:"pageInfo"`
}

type ExerciseResult struct {
	Exercise      *Exercise       `json:"exercise"`
	ValidationErr *ValidationErr  `json:"validationError"`
	NotFoundErr   *NotFoundErr    `json:"notFoundError"`
	AuthErr       *AuthErr        `json:"authError"`
}

type ArchiveResult struct {
	Exercise    *Exercise    `json:"exercise"`
	NotFoundErr *NotFoundErr `json:"notFoundError"`
	AuthErr     *AuthErr     `json:"authError"`
}

type ValidationErr struct {
	Message string             `json:"message"`
	Code    ExerciseErrorCode  `json:"code"`
}

type NotFoundErr struct {
	Message string             `json:"message"`
	Code    ExerciseErrorCode  `json:"code"`
}

type AuthErr struct {
	Message string             `json:"message"`
	Code    ExerciseErrorCode  `json:"code"`
}

type ExerciseErrorCode string

const (
	ExerciseErrorValidation   ExerciseErrorCode = "VALIDATION_ERROR"
	ExerciseErrorNotFound     ExerciseErrorCode = "NOT_FOUND"
	ExerciseErrorAuth         ExerciseErrorCode = "AUTH_ERROR"
	ExerciseErrorInternal     ExerciseErrorCode = "INTERNAL_ERROR"
)