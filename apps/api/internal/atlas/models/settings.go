package models

// SettingsRecord is the internal DB model with pinHash.
type SettingsRecord struct {
	ID                    string
	UserID                string
	PinEnabled            bool
	PinHash               *string
	Units                 string
	DefaultAiExportWeeks  int32
	CreatedAt             string
	UpdatedAt             string
}

// Settings is the public model that never exposes pinHash.
type Settings struct {
	PinEnabled           bool   `json:"pinEnabled"`
	Units                string `json:"units"`
	DefaultAiExportWeeks int    `json:"defaultAiExportWeeks"`
}

type SettingsInput struct {
	Units               *string `json:"units"`
	DefaultAiExportWeeks *int   `json:"defaultAiExportWeeks"`
}

type SettingsResult struct {
	Settings *Settings      `json:"settings"`
	Error    *SettingsError `json:"error"`
}

type SettingsError struct {
	Message string             `json:"message"`
	Code    SettingsErrorCode  `json:"code"`
}

type SettingsErrorCode string

const (
	SettingsErrorValidation   SettingsErrorCode = "VALIDATION_ERROR"
	SettingsErrorUnauthorized SettingsErrorCode = "UNAUTHORIZED"
	SettingsErrorInternal     SettingsErrorCode = "INTERNAL_ERROR"
)