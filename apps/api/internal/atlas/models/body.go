// FILE: apps/api/internal/atlas/models/body.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define WAVE-04 BodyWeight, BodyCheckIn, and BodyMeasurement domain models, enums, input types, and result/error types.
//   SCOPE: Internal DB records, public models, enums (BodyWeightSource, MeasurementType, MeasurementSide), input types for weight/check-in/measurement CRUD, and typed error codes.
//   DEPENDS: apps/api/internal/atlas/models/date.go.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   BodyWeightRecord - Internal DB model matching body_weight_entries table.
//   BodyWeightEntry - Public body weight entry model.
//   BodyCheckInRecord - Internal DB model matching body_check_ins table.
//   BodyCheckIn - Public body check-in model with nested measurements and photos.
//   BodyMeasurementRecord - Internal DB model matching body_measurements table.
//   BodyMeasurement - Public body measurement model.
//   BodyWeightSource, MeasurementType, MeasurementSide - Enums.
//   CreateBodyWeightInput, UpdateBodyWeightInput - Service/GraphQL inputs.
//   CreateCheckInInput, UpdateCheckInInput - Check-in inputs.
//   CreateMeasurementInput, UpdateMeasurementInput - Measurement inputs.
//   BodyWeightResult, BodyWeightEntriesResult - Weight entry results.
//   BodyCheckInResult, BodyCheckInsResult - Check-in results.
//   BodyMeasurementResult - Measurement result.
//   BodyErrorCode - Enum for body tracking error codes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added WAVE-04 body tracking domain models.
// END_MODULE_SUMMARY

package models

type BodyWeightSource string

const (
	BodyWeightSourceScale        BodyWeightSource = "SCALE"
	BodyWeightSourceManual       BodyWeightSource = "MANUAL"
	BodyWeightSourceBackupImport BodyWeightSource = "BACKUP_IMPORT"
	BodyWeightSourceUnknown      BodyWeightSource = "UNKNOWN"
)

var ValidBodyWeightSources = []BodyWeightSource{
	BodyWeightSourceScale, BodyWeightSourceManual,
	BodyWeightSourceBackupImport, BodyWeightSourceUnknown,
}

func IsValidBodyWeightSource(v string) bool {
	for _, s := range ValidBodyWeightSources {
		if string(s) == v {
			return true
		}
	}
	return false
}

type MeasurementType string

const (
	MeasurementTypeNeck      MeasurementType = "NECK"
	MeasurementTypeShoulders MeasurementType = "SHOULDERS"
	MeasurementTypeForearm   MeasurementType = "FOREARM"
	MeasurementTypeBiceps    MeasurementType = "BICEPS"
	MeasurementTypeChest     MeasurementType = "CHEST"
	MeasurementTypeWaist     MeasurementType = "WAIST"
	MeasurementTypeAbdomen   MeasurementType = "ABDOMEN"
	MeasurementTypeHips      MeasurementType = "HIPS"
	MeasurementTypeThigh     MeasurementType = "THIGH"
	MeasurementTypeCalf      MeasurementType = "CALF"
)

var ValidMeasurementTypes = []MeasurementType{
	MeasurementTypeNeck, MeasurementTypeShoulders, MeasurementTypeForearm,
	MeasurementTypeBiceps, MeasurementTypeChest, MeasurementTypeWaist,
	MeasurementTypeAbdomen, MeasurementTypeHips, MeasurementTypeThigh,
	MeasurementTypeCalf,
}

var PairedMeasurementTypes = map[MeasurementType]bool{
	MeasurementTypeForearm: true,
	MeasurementTypeBiceps:  true,
	MeasurementTypeThigh:   true,
	MeasurementTypeCalf:    true,
}

func IsValidMeasurementType(v string) bool {
	for _, mt := range ValidMeasurementTypes {
		if string(mt) == v {
			return true
		}
	}
	return false
}

func IsPairedMeasurementType(mt MeasurementType) bool {
	return PairedMeasurementTypes[mt]
}

type MeasurementSide string

const (
	MeasurementSideNone  MeasurementSide = "NONE"
	MeasurementSideLeft  MeasurementSide = "LEFT"
	MeasurementSideRight MeasurementSide = "RIGHT"
)

var ValidMeasurementSides = []MeasurementSide{
	MeasurementSideNone, MeasurementSideLeft, MeasurementSideRight,
}

func IsValidMeasurementSide(v string) bool {
	for _, s := range ValidMeasurementSides {
		if string(s) == v {
			return true
		}
	}
	return false
}

// BodyWeightEntry

type BodyWeightRecord struct {
	ID        string
	UserID    string
	Date      Date
	Weight    float64
	Source    string
	Notes     *string
	CreatedAt string
	UpdatedAt string
}

type BodyWeightEntry struct {
	ID        string            `json:"id"`
	UserID    string            `json:"userId"`
	Date      Date              `json:"date"`
	Weight    float64           `json:"weight"`
	Source    BodyWeightSource  `json:"source"`
	Notes     *string           `json:"notes"`
	CreatedAt string            `json:"createdAt"`
	UpdatedAt string            `json:"updatedAt"`
}

type CreateBodyWeightInput struct {
	Date   Date              `json:"date"`
	Weight float64           `json:"weight"`
	Source BodyWeightSource  `json:"source"`
	Notes  *string           `json:"notes"`
}

type UpdateBodyWeightInput struct {
	Weight *float64          `json:"weight"`
	Source *BodyWeightSource `json:"source"`
	Notes  *string           `json:"notes"`
}

type BodyWeightResult struct {
	Entry         *BodyWeightEntry    `json:"entry"`
	ValidationErr *BodyValidationErr  `json:"validationError"`
	NotFoundErr   *BodyNotFoundErr    `json:"notFoundError"`
	AuthErr       *BodyAuthErr        `json:"authError"`
}

type BodyWeightEntriesResult struct {
	Entries       []BodyWeightEntry   `json:"entries"`
	ValidationErr *BodyValidationErr  `json:"validationError"`
	AuthErr       *BodyAuthErr        `json:"authError"`
}

// BodyCheckIn

type BodyCheckInRecord struct {
	ID               string
	UserID           string
	Date             Date
	Weight           *float64
	BodyFatPercentage *float64
	Notes            *string
	CreatedAt        string
	UpdatedAt        string
}

type BodyCheckIn struct {
	ID                string              `json:"id"`
	UserID            string              `json:"userId"`
	Date              Date                `json:"date"`
	Weight            *float64            `json:"weight"`
	BodyFatPercentage *float64            `json:"bodyFatPercentage"`
	Notes             *string             `json:"notes"`
	Measurements      []BodyMeasurement   `json:"measurements"`
	ProgressPhotos    []ProgressPhoto     `json:"progressPhotos"`
	CreatedAt         string              `json:"createdAt"`
	UpdatedAt         string              `json:"updatedAt"`
}

type CreateCheckInInput struct {
	Date              Date     `json:"date"`
	Weight            *float64 `json:"weight"`
	BodyFatPercentage *float64 `json:"bodyFatPercentage"`
	Notes             *string  `json:"notes"`
}

type UpdateCheckInInput struct {
	Weight            *float64 `json:"weight"`
	BodyFatPercentage *float64 `json:"bodyFatPercentage"`
	Notes             *string  `json:"notes"`
}

type BodyCheckInResult struct {
	CheckIn       *BodyCheckIn         `json:"checkIn"`
	ValidationErr *BodyValidationErr   `json:"validationError"`
	NotFoundErr   *BodyNotFoundErr     `json:"notFoundError"`
	AuthErr       *BodyAuthErr         `json:"authError"`
}

type BodyCheckInsResult struct {
	CheckIns      []BodyCheckIn        `json:"checkIns"`
	ValidationErr *BodyValidationErr   `json:"validationError"`
	AuthErr       *BodyAuthErr         `json:"authError"`
}

// BodyMeasurement

type BodyMeasurementRecord struct {
	ID              string
	CheckInID       string
	MeasurementType string
	Side            *string
	Value           float64
	CreatedAt       string
	UpdatedAt       string
}

type BodyMeasurement struct {
	ID              string          `json:"id"`
	CheckInID       string          `json:"checkInId"`
	MeasurementType MeasurementType `json:"measurementType"`
	Side            *MeasurementSide `json:"side"`
	Value           float64         `json:"value"`
	CreatedAt       string          `json:"createdAt"`
	UpdatedAt       string          `json:"updatedAt"`
}

type CreateMeasurementInput struct {
	MeasurementType MeasurementType  `json:"measurementType"`
	Side            *MeasurementSide `json:"side"`
	Value           float64          `json:"value"`
}

type UpdateMeasurementInput struct {
	MeasurementType *MeasurementType `json:"measurementType"`
	Side            *MeasurementSide `json:"side"`
	Value           *float64         `json:"value"`
}

type BodyMeasurementResult struct {
	Measurement   *BodyMeasurement    `json:"measurement"`
	ValidationErr *BodyValidationErr  `json:"validationError"`
	NotFoundErr   *BodyNotFoundErr    `json:"notFoundError"`
	AuthErr       *BodyAuthErr        `json:"authError"`
}

// Common Body Error Types

type BodyValidationErr struct {
	Message string        `json:"message"`
	Code    BodyErrorCode `json:"code"`
}

func (e *BodyValidationErr) Error() string {
	if e == nil || e.Message == "" {
		return "body tracking validation error"
	}
	return e.Message
}

type BodyNotFoundErr struct {
	Message string        `json:"message"`
	Code    BodyErrorCode `json:"code"`
}

func (e *BodyNotFoundErr) Error() string {
	if e == nil || e.Message == "" {
		return "body tracking entry not found"
	}
	return e.Message
}

type BodyAuthErr struct {
	Message string        `json:"message"`
	Code    BodyErrorCode `json:"code"`
}

func (e *BodyAuthErr) Error() string {
	if e == nil || e.Message == "" {
		return "body tracking auth error"
	}
	return e.Message
}

type BodyErrorCode string

const (
	BodyErrorValidation BodyErrorCode = "VALIDATION_ERROR"
	BodyErrorNotFound   BodyErrorCode = "NOT_FOUND"
	BodyErrorAuth       BodyErrorCode = "AUTH_ERROR"
	BodyErrorInternal   BodyErrorCode = "INTERNAL_ERROR"
)

func BodyWeightEntryFromRecord(r *BodyWeightRecord) *BodyWeightEntry {
	if r == nil {
		return nil
	}
	return &BodyWeightEntry{
		ID:        r.ID,
		UserID:    r.UserID,
		Date:      r.Date,
		Weight:    r.Weight,
		Source:    BodyWeightSource(r.Source),
		Notes:     r.Notes,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func BodyCheckInFromRecord(r *BodyCheckInRecord) *BodyCheckIn {
	if r == nil {
		return nil
	}
	return &BodyCheckIn{
		ID:                r.ID,
		UserID:            r.UserID,
		Date:              r.Date,
		Weight:            r.Weight,
		BodyFatPercentage: r.BodyFatPercentage,
		Notes:             r.Notes,
		Measurements:      []BodyMeasurement{},
		ProgressPhotos:    []ProgressPhoto{},
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}

func BodyMeasurementFromRecord(r *BodyMeasurementRecord) *BodyMeasurement {
	if r == nil {
		return nil
	}
	mt := MeasurementType(r.MeasurementType)
	var side *MeasurementSide
	if r.Side != nil {
		v := MeasurementSide(*r.Side)
		side = &v
	}
	return &BodyMeasurement{
		ID:              r.ID,
		CheckInID:       r.CheckInID,
		MeasurementType: mt,
		Side:            side,
		Value:           r.Value,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}