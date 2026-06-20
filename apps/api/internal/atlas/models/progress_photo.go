// FILE: apps/api/internal/atlas/models/progress_photo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define WAVE-04 ProgressPhoto domain models, enums, and public types.
//   SCOPE: ProgressPhotoRecord (DB), ProgressPhoto (public), ProgressPhotoAngle enum, create metadata input.
//   DEPENDS: apps/api/internal/atlas/models/date.go.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ProgressPhotoRecord - Internal DB model matching progress_photos table.
//   ProgressPhoto - Public progress photo model.
//   ProgressPhotoAngle - Enum of photo angle values.
//   UploadPhotoInput - Input for progress photo upload metadata.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added WAVE-04 progress photo domain models.
// END_CHANGE_SUMMARY

package models

type ProgressPhotoAngle string

const (
	ProgressPhotoAngleFront  ProgressPhotoAngle = "FRONT"
	ProgressPhotoAngleSide   ProgressPhotoAngle = "SIDE"
	ProgressPhotoAngleBack   ProgressPhotoAngle = "BACK"
	ProgressPhotoAngleCustom ProgressPhotoAngle = "CUSTOM"
)

var ValidProgressPhotoAngles = []ProgressPhotoAngle{
	ProgressPhotoAngleFront, ProgressPhotoAngleSide,
	ProgressPhotoAngleBack, ProgressPhotoAngleCustom,
}

func IsValidProgressPhotoAngle(v string) bool {
	for _, a := range ValidProgressPhotoAngles {
		if string(a) == v {
			return true
		}
	}
	return false
}

type ProgressPhotoRecord struct {
	ID               string
	CheckInID        string
	FilePath         string
	OriginalFileName string
	MimeType         string
	SizeBytes        int64
	Angle            *string
	Label            *string
	Notes            *string
	CreatedAt        string
	UpdatedAt        string
}

type ProgressPhoto struct {
	ID               string              `json:"id"`
	CheckInID        string              `json:"checkInId"`
	FilePath         string              `json:"-"`
	OriginalFileName string              `json:"originalFileName"`
	MimeType         string              `json:"mimeType"`
	SizeBytes        int64               `json:"sizeBytes"`
	Angle            *ProgressPhotoAngle `json:"angle"`
	Label            *string             `json:"label"`
	Notes            *string             `json:"notes"`
	CreatedAt        string              `json:"createdAt"`
	UpdatedAt        string              `json:"updatedAt"`
}

type UploadPhotoInput struct {
	CheckInID string             `json:"checkInId"`
	Angle     *ProgressPhotoAngle `json:"angle"`
	Label     *string            `json:"label"`
	Notes     *string            `json:"notes"`
}

func ProgressPhotoFromRecord(r *ProgressPhotoRecord) *ProgressPhoto {
	if r == nil {
		return nil
	}
	var angle *ProgressPhotoAngle
	if r.Angle != nil {
		v := ProgressPhotoAngle(*r.Angle)
		angle = &v
	}
	return &ProgressPhoto{
		ID:               r.ID,
		CheckInID:        r.CheckInID,
		FilePath:         r.FilePath,
		OriginalFileName: r.OriginalFileName,
		MimeType:         r.MimeType,
		SizeBytes:        r.SizeBytes,
		Angle:            angle,
		Label:            r.Label,
		Notes:            r.Notes,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
}

type ProgressPhotosResult struct {
	Photos        []ProgressPhoto            `json:"photos"`
	ValidationErr *ProgressPhotoValidationErr `json:"validationError"`
	AuthErr       *ProgressPhotoAuthErr       `json:"authError"`
}

type ProgressPhotoValidationErr struct {
	Message string                `json:"message"`
	Code    ProgressPhotoErrorCode `json:"code"`
}

func (e *ProgressPhotoValidationErr) Error() string {
	if e == nil || e.Message == "" {
		return "progress photo validation error"
	}
	return e.Message
}

type ProgressPhotoAuthErr struct {
	Message string                `json:"message"`
	Code    ProgressPhotoErrorCode `json:"code"`
}

func (e *ProgressPhotoAuthErr) Error() string {
	if e == nil || e.Message == "" {
		return "progress photo auth error"
	}
	return e.Message
}

type ProgressPhotoErrorCode string

const (
	ProgressPhotoErrorValidation ProgressPhotoErrorCode = "VALIDATION_ERROR"
	ProgressPhotoErrorNotFound   ProgressPhotoErrorCode = "NOT_FOUND"
	ProgressPhotoErrorAuth       ProgressPhotoErrorCode = "AUTH_ERROR"
	ProgressPhotoErrorInternal   ProgressPhotoErrorCode = "INTERNAL_ERROR"
)