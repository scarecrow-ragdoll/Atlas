// FILE: apps/api/internal/atlas/models/user_profile.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define WAVE-07 UserProfile domain models, input types, result/error types.
//   SCOPE: UserProfileRecord (DB), UserProfile (public), UserProfileInput, typed error codes, UserProfileFromRecord converter.
//   DEPENDS: apps/api/internal/atlas/models/date.go.
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   UserProfileRecord - Internal DB model matching user_profiles table.
//   UserProfile - Public user profile model.
//   UserProfileInput - Input for creating/updating a user profile.
//   UserProfileResult - Union-like result for user profile operations.
//   UserProfileErrorCode - Enum for user profile error codes.
//   UserProfileFromRecord - Converts DB record to public model.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added WAVE-07 user profile domain models.
// END_CHANGE_SUMMARY

package models

type UserProfileRecord struct {
	ID                       string
	UserID                   string
	Goal                     *string
	Height                   *float64
	BirthDate                *Date
	TrainingExperience       *string
	CurrentTrainingSplit     *string
	PreferredProgressionStyle *string
	NutritionStrategy        *string
	PersistentAiContext      *string
	CreatedAt                string
	UpdatedAt                string
}

type UserProfile struct {
	ID                       string   `json:"id"`
	UserID                   string   `json:"userId"`
	Goal                     *string  `json:"goal"`
	Height                   *float64 `json:"height"`
	BirthDate                *Date    `json:"birthDate"`
	TrainingExperience       *string  `json:"trainingExperience"`
	CurrentTrainingSplit     *string  `json:"currentTrainingSplit"`
	PreferredProgressionStyle *string  `json:"preferredProgressionStyle"`
	NutritionStrategy        *string  `json:"nutritionStrategy"`
	PersistentAiContext      *string  `json:"persistentAiContext"`
	CreatedAt                string   `json:"createdAt"`
	UpdatedAt                string   `json:"updatedAt"`
}

type UserProfileInput struct {
	Goal                     *string `json:"goal"`
	Height                   *float64 `json:"height"`
	BirthDate                *Date   `json:"birthDate"`
	TrainingExperience       *string `json:"trainingExperience"`
	CurrentTrainingSplit     *string `json:"currentTrainingSplit"`
	PreferredProgressionStyle *string `json:"preferredProgressionStyle"`
	NutritionStrategy        *string `json:"nutritionStrategy"`
	PersistentAiContext      *string `json:"persistentAiContext"`
}

type UserProfileResult struct {
	Profile        *UserProfile                  `json:"profile"`
	ValidationErr  *UserProfileValidationErr      `json:"validationError"`
	NotFoundErr    *UserProfileNotFoundErr        `json:"notFoundError"`
	AuthErr        *UserProfileAuthErr            `json:"authError"`
}

type UserProfileValidationErr struct {
	Message string                  `json:"message"`
	Code    UserProfileErrorCode    `json:"code"`
}

func (e *UserProfileValidationErr) Error() string {
	if e == nil || e.Message == "" {
		return "user profile validation error"
	}
	return e.Message
}

type UserProfileNotFoundErr struct {
	Message string                  `json:"message"`
	Code    UserProfileErrorCode    `json:"code"`
}

func (e *UserProfileNotFoundErr) Error() string {
	if e == nil || e.Message == "" {
		return "user profile not found"
	}
	return e.Message
}

type UserProfileAuthErr struct {
	Message string                  `json:"message"`
	Code    UserProfileErrorCode    `json:"code"`
}

func (e *UserProfileAuthErr) Error() string {
	if e == nil || e.Message == "" {
		return "user profile auth error"
	}
	return e.Message
}

type UserProfileErrorCode string

const (
	UserProfileErrorValidation UserProfileErrorCode = "VALIDATION_ERROR"
	UserProfileErrorNotFound   UserProfileErrorCode = "NOT_FOUND"
	UserProfileErrorAuth       UserProfileErrorCode = "AUTH_ERROR"
	UserProfileErrorInternal   UserProfileErrorCode = "INTERNAL_ERROR"
)

func UserProfileFromRecord(r *UserProfileRecord) *UserProfile {
	if r == nil {
		return nil
	}
	return &UserProfile{
		ID:                        r.ID,
		UserID:                    r.UserID,
		Goal:                      r.Goal,
		Height:                    r.Height,
		BirthDate:                 r.BirthDate,
		TrainingExperience:        r.TrainingExperience,
		CurrentTrainingSplit:      r.CurrentTrainingSplit,
		PreferredProgressionStyle: r.PreferredProgressionStyle,
		NutritionStrategy:         r.NutritionStrategy,
		PersistentAiContext:       r.PersistentAiContext,
		CreatedAt:                 r.CreatedAt,
		UpdatedAt:                 r.UpdatedAt,
	}
}