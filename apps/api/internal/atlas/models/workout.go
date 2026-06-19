// FILE: apps/api/internal/atlas/models/workout.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define WAVE-03 DailyLog, workout exercise, workout set, input, result, and typed error models.
//   SCOPE: Strength diary aggregate shapes, nullable fields, result envelopes, conflict payloads, and stable error codes.
//   DEPENDS: apps/api/internal/atlas/models/date.go, apps/api/internal/atlas/models/exercise.go.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   DailyLogRecord - Internal DailyLog row model.
//   DailyLog - Public DailyLog aggregate model.
//   DailyLogSummary - Public DailyLog range summary.
//   WorkoutExerciseRecord - Internal workout exercise row model.
//   WorkoutExercise - Public workout exercise model.
//   WorkoutSetRecord - Internal workout set row model.
//   WorkoutSet - Public workout set model.
//   AddWorkoutExerciseInput - Service input for adding a workout exercise.
//   UpdateWorkoutExerciseInput - Service input for updating workout exercise position or notes.
//   AddWorkoutSetInput - Service input for adding a workout set.
//   UpdateWorkoutSetInput - Service input for updating workout set values, order, or notes.
//   DailyLogResult - Union-like DailyLog operation result.
//   DailyLogErrorCode - Stable DailyLog error code enum.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added WAVE-03 workout diary models and typed errors.
// END_CHANGE_SUMMARY

package models

type DailyLogRecord struct {
	ID        string
	UserID    string
	Date      Date
	Notes     *string
	Version   int32
	CreatedAt string
	UpdatedAt string
}

type DailyLog struct {
	ID               string            `json:"id"`
	UserID           string            `json:"userId"`
	Date             Date              `json:"date"`
	Notes            *string           `json:"notes"`
	Version          int32             `json:"version"`
	WorkoutExercises []WorkoutExercise `json:"workoutExercises"`
	CreatedAt        string            `json:"createdAt"`
	UpdatedAt        string            `json:"updatedAt"`
}

type DailyLogSummary struct {
	ID                   string  `json:"id"`
	Date                 Date    `json:"date"`
	Version              int32   `json:"version"`
	WorkoutExerciseCount int32   `json:"workoutExerciseCount"`
	WorkoutSetCount      int32   `json:"workoutSetCount"`
	TotalVolume          float64 `json:"totalVolume"`
	UpdatedAt            string  `json:"updatedAt"`
}

type WorkoutExerciseRecord struct {
	ID                    string
	UserID                string
	DailyLogID            string
	ExerciseID            string
	Position              int32
	WorkingWeightSnapshot *float64
	Notes                 *string
	Sets                  []WorkoutSetRecord
	CreatedAt             string
	UpdatedAt             string
}

type WorkoutExercise struct {
	ID                    string       `json:"id"`
	UserID                string       `json:"userId"`
	DailyLogID            string       `json:"dailyLogId"`
	ExerciseID            string       `json:"exerciseId"`
	Exercise              *Exercise    `json:"exercise"`
	Position              int32        `json:"position"`
	WorkingWeightSnapshot *float64     `json:"workingWeightSnapshot"`
	Notes                 *string      `json:"notes"`
	Sets                  []WorkoutSet `json:"sets"`
	CreatedAt             string       `json:"createdAt"`
	UpdatedAt             string       `json:"updatedAt"`
}

type WorkoutSetRecord struct {
	ID                string
	WorkoutExerciseID string
	SetNumber         int32
	Weight            float64
	Reps              int32
	RPE               *float64
	RIR               *int32
	Notes             *string
	CreatedAt         string
	UpdatedAt         string
}

type WorkoutSet struct {
	ID                string   `json:"id"`
	WorkoutExerciseID string   `json:"workoutExerciseId"`
	SetNumber         int32    `json:"setNumber"`
	Weight            float64  `json:"weight"`
	Reps              int32    `json:"reps"`
	RPE               *float64 `json:"rpe"`
	RIR               *int32   `json:"rir"`
	Notes             *string  `json:"notes"`
	CreatedAt         string   `json:"createdAt"`
	UpdatedAt         string   `json:"updatedAt"`
}

type AddWorkoutExerciseInput struct {
	ExerciseID string  `json:"exerciseId"`
	Position   *int32  `json:"position"`
	Notes      *string `json:"notes"`
}

type UpdateWorkoutExerciseInput struct {
	Position *int32  `json:"position"`
	SetNotes bool    `json:"-"`
	Notes    *string `json:"notes"`
}

type AddWorkoutSetInput struct {
	SetNumber *int32   `json:"setNumber"`
	Weight    float64  `json:"weight"`
	Reps      int32    `json:"reps"`
	RPE       *float64 `json:"rpe"`
	RIR       *int32   `json:"rir"`
	Notes     *string  `json:"notes"`
}

type UpdateWorkoutSetInput struct {
	SetNumber *int32   `json:"setNumber"`
	Weight    *float64 `json:"weight"`
	Reps      *int32   `json:"reps"`
	SetRPE    bool     `json:"-"`
	RPE       *float64 `json:"rpe"`
	SetRIR    bool     `json:"-"`
	RIR       *int32   `json:"rir"`
	SetNotes  bool     `json:"-"`
	Notes     *string  `json:"notes"`
}

type DailyLogResult struct {
	DailyLog      *DailyLog              `json:"dailyLog"`
	ValidationErr *DailyLogValidationErr `json:"validationError"`
	NotFoundErr   *DailyLogNotFoundErr   `json:"notFoundError"`
	ConflictErr   *DailyLogConflictErr   `json:"conflictError"`
	AuthErr       *DailyLogAuthErr       `json:"authError"`
}

type DailyLogValidationErr struct {
	Message string            `json:"message"`
	Code    DailyLogErrorCode `json:"code"`
}

func (e *DailyLogValidationErr) Error() string {
	if e == nil || e.Message == "" {
		return "daily log validation error"
	}
	return e.Message
}

type DailyLogNotFoundErr struct {
	Message string            `json:"message"`
	Code    DailyLogErrorCode `json:"code"`
}

func (e *DailyLogNotFoundErr) Error() string {
	if e == nil || e.Message == "" {
		return "daily log not found"
	}
	return e.Message
}

type DailyLogConflictErr struct {
	Message         string            `json:"message"`
	Code            DailyLogErrorCode `json:"code"`
	CurrentVersion  int32             `json:"currentVersion"`
	CurrentDailyLog *DailyLog         `json:"currentDailyLog"`
}

func (e *DailyLogConflictErr) Error() string {
	if e == nil || e.Message == "" {
		return "daily log version conflict"
	}
	return e.Message
}

type DailyLogAuthErr struct {
	Message string            `json:"message"`
	Code    DailyLogErrorCode `json:"code"`
}

func (e *DailyLogAuthErr) Error() string {
	if e == nil || e.Message == "" {
		return "daily log auth error"
	}
	return e.Message
}

type DailyLogErrorCode string

const (
	DailyLogErrorValidation DailyLogErrorCode = "VALIDATION_ERROR"
	DailyLogErrorNotFound   DailyLogErrorCode = "NOT_FOUND"
	DailyLogErrorConflict   DailyLogErrorCode = "CONFLICT"
	DailyLogErrorAuth       DailyLogErrorCode = "AUTH_ERROR"
	DailyLogErrorInternal   DailyLogErrorCode = "INTERNAL_ERROR"
)
