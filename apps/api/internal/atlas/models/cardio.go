// FILE: apps/api/internal/atlas/models/cardio.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define WAVE-04 CardioEntry domain models, enums, input types, and result/error types.
//   SCOPE: Internal CardioRecord, public CardioEntry, CardioType enum, HeartRateZone enum, input types (CreateCardioInput, UpdateCardioInput), result types (CardioEntryResult, CardioEntriesResult), and typed error codes.
//   DEPENDS: apps/api/internal/atlas/models/date.go.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   CardioRecord - Internal DB model matching cardio_entries table.
//   CardioEntry - Public cardio entry model.
//   CreateCardioInput - Service/GraphQL input for creating cardio entries.
//   UpdateCardioInput - Service/GraphQL input for updating cardio entries.
//   CardioEntriesInput - Input for listing cardio entries by date range.
//   CardioEntryResult - Union-like result for single cardio entry operations.
//   CardioEntriesResult - Union-like result for listing cardio entries.
//   CardioType - Enum of allowed cardio activity types.
//   HeartRateZone - Enum of heart rate zone values.
//   CardioErrorCode - Enum for cardio error codes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added WAVE-04 cardio domain models.
// END_CHANGE_SUMMARY

package models

type CardioType string

const (
	CardioTypeWalking    CardioType = "WALKING"
	CardioTypeRunning    CardioType = "RUNNING"
	CardioTypeTreadmill  CardioType = "TREADMILL"
	CardioTypeBike       CardioType = "BIKE"
	CardioTypeElliptical CardioType = "ELLIPTICAL"
	CardioTypeStairs     CardioType = "STAIRS"
	CardioTypeRowing     CardioType = "ROWING"
	CardioTypeOther      CardioType = "OTHER"
)

var ValidCardioTypes = []CardioType{
	CardioTypeWalking, CardioTypeRunning, CardioTypeTreadmill,
	CardioTypeBike, CardioTypeElliptical, CardioTypeStairs,
	CardioTypeRowing, CardioTypeOther,
}

func IsValidCardioType(v string) bool {
	for _, t := range ValidCardioTypes {
		if string(t) == v {
			return true
		}
	}
	return false
}

type HeartRateZone string

const (
	HeartRateZoneUnknown HeartRateZone = "UNKNOWN"
	HeartRateZone1       HeartRateZone = "ZONE_1"
	HeartRateZone2       HeartRateZone = "ZONE_2"
	HeartRateZone3       HeartRateZone = "ZONE_3"
	HeartRateZone4       HeartRateZone = "ZONE_4"
	HeartRateZone5       HeartRateZone = "ZONE_5"
)

var ValidHeartRateZones = []HeartRateZone{
	HeartRateZoneUnknown, HeartRateZone1, HeartRateZone2,
	HeartRateZone3, HeartRateZone4, HeartRateZone5,
}

func IsValidHeartRateZone(v string) bool {
	for _, z := range ValidHeartRateZones {
		if string(z) == v {
			return true
		}
	}
	return false
}

type CardioRecord struct {
	ID              string
	UserID          string
	DailyLogID      string
	CardioType      string
	DurationMinutes int32
	AvgPulse        *int32
	HeartRateZone   *string
	Notes           *string
	CreatedAt       string
	UpdatedAt       string
}

type CardioEntry struct {
	ID              string         `json:"id"`
	UserID          string         `json:"userId"`
	DailyLogID      string         `json:"dailyLogId"`
	CardioType      CardioType     `json:"cardioType"`
	DurationMinutes int32          `json:"durationMinutes"`
	AvgPulse        *int32         `json:"avgPulse"`
	HeartRateZone   *HeartRateZone `json:"heartRateZone"`
	Notes           *string        `json:"notes"`
	CreatedAt       string         `json:"createdAt"`
	UpdatedAt       string         `json:"updatedAt"`
}

type CreateCardioInput struct {
	Date            Date            `json:"date"`
	CardioType      CardioType      `json:"cardioType"`
	DurationMinutes int32           `json:"durationMinutes"`
	AvgPulse        *int32          `json:"avgPulse"`
	HeartRateZone   *HeartRateZone  `json:"heartRateZone"`
	Notes           *string         `json:"notes"`
}

type UpdateCardioInput struct {
	CardioType      *CardioType     `json:"cardioType"`
	DurationMinutes *int32          `json:"durationMinutes"`
	AvgPulse        *int32          `json:"avgPulse"`
	HeartRateZone   *HeartRateZone  `json:"heartRateZone"`
	Notes           *string         `json:"notes"`
}

type CardioEntryResult struct {
	CardioEntry   *CardioEntry        `json:"cardioEntry"`
	ValidationErr *CardioValidationErr `json:"validationError"`
	NotFoundErr   *CardioNotFoundErr   `json:"notFoundError"`
	AuthErr       *CardioAuthErr       `json:"authError"`
}

type CardioEntriesResult struct {
	Entries       []CardioEntry       `json:"entries"`
	ValidationErr *CardioValidationErr `json:"validationError"`
	AuthErr       *CardioAuthErr       `json:"authError"`
}

type CardioValidationErr struct {
	Message string          `json:"message"`
	Code    CardioErrorCode `json:"code"`
}

func (e *CardioValidationErr) Error() string {
	if e == nil || e.Message == "" {
		return "cardio validation error"
	}
	return e.Message
}

type CardioNotFoundErr struct {
	Message string          `json:"message"`
	Code    CardioErrorCode `json:"code"`
}

func (e *CardioNotFoundErr) Error() string {
	if e == nil || e.Message == "" {
		return "cardio not found"
	}
	return e.Message
}

type CardioAuthErr struct {
	Message string          `json:"message"`
	Code    CardioErrorCode `json:"code"`
}

func (e *CardioAuthErr) Error() string {
	if e == nil || e.Message == "" {
		return "cardio auth error"
	}
	return e.Message
}

type CardioErrorCode string

const (
	CardioErrorValidation CardioErrorCode = "VALIDATION_ERROR"
	CardioErrorNotFound   CardioErrorCode = "NOT_FOUND"
	CardioErrorAuth       CardioErrorCode = "AUTH_ERROR"
	CardioErrorInternal   CardioErrorCode = "INTERNAL_ERROR"
)

func CardioEntryFromRecord(r *CardioRecord) *CardioEntry {
	if r == nil {
		return nil
	}
	ct := CardioType(r.CardioType)
	var hrz *HeartRateZone
	if r.HeartRateZone != nil {
		v := HeartRateZone(*r.HeartRateZone)
		hrz = &v
	}
	return &CardioEntry{
		ID:              r.ID,
		UserID:          r.UserID,
		DailyLogID:      r.DailyLogID,
		CardioType:      ct,
		DurationMinutes: r.DurationMinutes,
		AvgPulse:        r.AvgPulse,
		HeartRateZone:   hrz,
		Notes:           r.Notes,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

func CardioEntryResultFromError(err error) *CardioEntryResult {
	if err == nil {
		return nil
	}
	var validationErr *CardioValidationErr
	if fmtErrorAs(err, &validationErr) {
		return &CardioEntryResult{ValidationErr: validationErr}
	}
	var notFoundErr *CardioNotFoundErr
	if fmtErrorAs(err, &notFoundErr) {
		return &CardioEntryResult{NotFoundErr: notFoundErr}
	}
	var authErr *CardioAuthErr
	if fmtErrorAs(err, &authErr) {
		return &CardioEntryResult{AuthErr: authErr}
	}
	return nil
}

func fmtErrorAs(err error, target any) bool {
	for {
		if e, ok := err.(interface{ As(any) bool }); ok {
			if e.As(target) {
				return true
			}
		}
		if u, ok := err.(interface{ Unwrap() error }); ok {
			err = u.Unwrap()
		} else {
			return false
		}
	}
}