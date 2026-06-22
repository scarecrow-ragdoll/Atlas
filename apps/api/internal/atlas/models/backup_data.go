// FILE: apps/api/internal/atlas/models/backup_data.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define the backup data transfer types for WAVE-09 full user data export and import.
//   SCOPE: BackupManifest, BackupDataArchive, BackupData structs and helpers for serializing/deserializing backup archives.
//   DEPENDS: All model types used in BackupData fields.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   BackupSchemaVersion - Current schema version constant.
//   BackupAppVersion - Current app version constant.
//   BackupManifest - Manifest metadata for backup archives.
//   BackupData - Structured container for all user data sections.
//   BackupDataArchive - ZIP archive structure with Manifest + Data.
//   NewBackupManifest - Creates a BackupManifest with provided options.
// END_MODULE_MAP

package models

import "time"

const (
	BackupSchemaVersion = 1
	BackupAppVersion    = "1.0.0"
)

type BackupManifest struct {
	Type             string         `json:"type"`
	SchemaVersion    int            `json:"schemaVersion"`
	AppVersion       string         `json:"appVersion"`
	ExportedAt       string         `json:"exportedAt"`
	IncludedSections []string       `json:"includedSections"`
	MediaIncluded    bool           `json:"mediaIncluded"`
	EntityCounts     map[string]int `json:"entityCounts,omitempty"`
}

type BackupData struct {
	Settings               *SettingsRecord                  `json:"settings,omitempty"`
	UserProfile            *UserProfileRecord               `json:"userProfile,omitempty"`
	Exercises              []ExerciseRecord                 `json:"exercises,omitempty"`
	ExerciseMedia          []ExerciseMediaRecord            `json:"exerciseMedia,omitempty"`
	DailyLogs              []DailyLogRecord                 `json:"dailyLogs,omitempty"`
	CardioEntries          []CardioRecord                   `json:"cardioEntries,omitempty"`
	BodyWeightEntries      []BodyWeightRecord               `json:"bodyWeightEntries,omitempty"`
	BodyCheckIns           []BodyCheckInRecord              `json:"bodyCheckIns,omitempty"`
	BodyMeasurements       []BodyMeasurementRecord          `json:"bodyMeasurements,omitempty"`
	ProgressPhotos         []ProgressPhotoRecord            `json:"progressPhotos,omitempty"`
	NutritionProducts      []NutritionProductRecord         `json:"nutritionProducts,omitempty"`
	NutritionTemplates     []NutritionTemplateRecord        `json:"nutritionTemplates,omitempty"`
	NutritionTemplateItems []NutritionTemplateItemRecord    `json:"nutritionTemplateItems,omitempty"`
	NutritionOverrides     []DailyNutritionOverrideRecord   `json:"nutritionOverrides,omitempty"`
	NutritionOverrideItems []DailyNutritionOverrideItemRecord `json:"nutritionOverrideItems,omitempty"`
	WeekFlags              []WeekFlagRecord                 `json:"weekFlags,omitempty"`
	AiExports              []AiExportRecord                 `json:"aiExports,omitempty"`
	AiReviews              []AiReviewRecord                 `json:"aiReviews,omitempty"`
}

type BackupDataArchive struct {
	Manifest BackupManifest `json:"manifest"`
	Data     BackupData     `json:"data"`
}

func NewBackupManifest(mediaIncluded bool, entityCounts map[string]int) BackupManifest {
	sections := []string{}
	for s := range entityCounts {
		sections = append(sections, s)
	}
	return BackupManifest{
		Type:             "full_backup",
		SchemaVersion:    BackupSchemaVersion,
		AppVersion:       BackupAppVersion,
		ExportedAt:       time.Now().UTC().Format(time.RFC3339),
		IncludedSections: sections,
		MediaIncluded:    mediaIncluded,
		EntityCounts:     entityCounts,
	}
}
