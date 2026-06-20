// FILE: apps/api/internal/atlas/models/week_flag.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define WAVE-04 WeekFlag domain models, enums, input types, and result/error types.
//   SCOPE: WeekFlagRecord (DB), WeekFlag (public), WeekFlagType enum, create/delete input, typed error codes.
//   DEPENDS: apps/api/internal/atlas/models/date.go.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   WeekFlagRecord - Internal DB model matching week_flags table.
//   WeekFlag - Public week flag model.
//   WeekFlagType - Enum of allowed week flag values.
//   CreateWeekFlagInput - Input for creating a week flag.
//   WeekFlagResult - Union-like result for week flag operations.
//   WeekFlagsResult - Union-like result for listing week flags.
//   WeekFlagErrorCode - Enum for week flag error codes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added WAVE-04 week flag domain models.
// END_CHANGE_SUMMARY

package models

type WeekFlagType string

const (
	WeekFlagTypePoorSleep       WeekFlagType = "POOR_SLEEP"
	WeekFlagTypeHighStress      WeekFlagType = "HIGH_STRESS"
	WeekFlagTypeIllness         WeekFlagType = "ILLNESS"
	WeekFlagTypeInjuryPain      WeekFlagType = "INJURY_PAIN"
	WeekFlagTypeCycle           WeekFlagType = "CYCLE"
	WeekFlagTypeCalorieDeficit  WeekFlagType = "CALORIE_DEFICIT"
	WeekFlagTypeSurplus         WeekFlagType = "SURPLUS"
	WeekFlagTypeMaintenance     WeekFlagType = "MAINTENANCE"
	WeekFlagTypeMissedWorkouts  WeekFlagType = "MISSED_WORKOUTS"
	WeekFlagTypeTravel          WeekFlagType = "TRAVEL"
)

var ValidWeekFlagTypes = []WeekFlagType{
	WeekFlagTypePoorSleep, WeekFlagTypeHighStress, WeekFlagTypeIllness,
	WeekFlagTypeInjuryPain, WeekFlagTypeCycle, WeekFlagTypeCalorieDeficit,
	WeekFlagTypeSurplus, WeekFlagTypeMaintenance, WeekFlagTypeMissedWorkouts,
	WeekFlagTypeTravel,
}

func IsValidWeekFlagType(v string) bool {
	for _, ft := range ValidWeekFlagTypes {
		if string(ft) == v {
			return true
		}
	}
	return false
}

type WeekFlagRecord struct {
	ID            string
	UserID        string
	WeekStartDate Date
	FlagType      string
	Notes         *string
	CreatedAt     string
	UpdatedAt     string
}

type WeekFlag struct {
	ID            string       `json:"id"`
	UserID        string       `json:"userId"`
	WeekStartDate Date         `json:"weekStartDate"`
	FlagType      WeekFlagType `json:"flagType"`
	Notes         *string      `json:"notes"`
	CreatedAt     string       `json:"createdAt"`
	UpdatedAt     string       `json:"updatedAt"`
}

type CreateWeekFlagInput struct {
	WeekStartDate Date         `json:"weekStartDate"`
	FlagType      WeekFlagType `json:"flagType"`
	Notes         *string      `json:"notes"`
}

type WeekFlagResult struct {
	WeekFlag      *WeekFlag             `json:"weekFlag"`
	ValidationErr *WeekFlagValidationErr `json:"validationError"`
	NotFoundErr   *WeekFlagNotFoundErr   `json:"notFoundError"`
	AuthErr       *WeekFlagAuthErr       `json:"authError"`
}

type WeekFlagsResult struct {
	Flags         []WeekFlag            `json:"flags"`
	ValidationErr *WeekFlagValidationErr `json:"validationError"`
	AuthErr       *WeekFlagAuthErr       `json:"authError"`
}

type WeekFlagValidationErr struct {
	Message string            `json:"message"`
	Code    WeekFlagErrorCode `json:"code"`
}

func (e *WeekFlagValidationErr) Error() string {
	if e == nil || e.Message == "" {
		return "week flag validation error"
	}
	return e.Message
}

type WeekFlagNotFoundErr struct {
	Message string            `json:"message"`
	Code    WeekFlagErrorCode `json:"code"`
}

func (e *WeekFlagNotFoundErr) Error() string {
	if e == nil || e.Message == "" {
		return "week flag not found"
	}
	return e.Message
}

type WeekFlagAuthErr struct {
	Message string            `json:"message"`
	Code    WeekFlagErrorCode `json:"code"`
}

func (e *WeekFlagAuthErr) Error() string {
	if e == nil || e.Message == "" {
		return "week flag auth error"
	}
	return e.Message
}

type WeekFlagErrorCode string

const (
	WeekFlagErrorValidation WeekFlagErrorCode = "VALIDATION_ERROR"
	WeekFlagErrorNotFound   WeekFlagErrorCode = "NOT_FOUND"
	WeekFlagErrorAuth       WeekFlagErrorCode = "AUTH_ERROR"
	WeekFlagErrorInternal   WeekFlagErrorCode = "INTERNAL_ERROR"
)

func WeekFlagFromRecord(r *WeekFlagRecord) *WeekFlag {
	if r == nil {
		return nil
	}
	return &WeekFlag{
		ID:            r.ID,
		UserID:        r.UserID,
		WeekStartDate: r.WeekStartDate,
		FlagType:      WeekFlagType(r.FlagType),
		Notes:         r.Notes,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}