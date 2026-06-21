package service

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ManifestSection represents a section entry in manifest.json
type ManifestSection struct {
	Workouts     bool `json:"workouts"`
	Cardio       bool `json:"cardio"`
	BodyWeight   bool `json:"bodyWeight"`
	Measurements bool `json:"measurements"`
	Nutrition    bool `json:"nutrition"`
	Photos       bool `json:"photos"`
}

// Manifest represents the export manifest.json structure
type Manifest struct {
	SchemaVersion    int              `json:"schemaVersion"`
	ExportTimestamp  string           `json:"exportTimestamp"`
	DateRangeStart   string           `json:"dateRangeStart"`
	DateRangeEnd     string           `json:"dateRangeEnd"`
	IncludedSections ManifestSection `json:"includedSections"`
}

// ExportData represents the full data.json structure
type ExportData struct {
	Workouts          []any            `json:"workouts"`
	Cardio            []any            `json:"cardio"`
	BodyWeightEntries []any            `json:"bodyWeightEntries"`
	Measurements      []any            `json:"measurements"`
	Nutrition         ExportNutrition  `json:"nutrition"`
	WeekFlags         []any            `json:"weekFlags"`
	UserProfile       ExportProfile    `json:"userProfile"`
}

type ExportNutrition struct {
	Products  []any `json:"products"`
	Templates []any `json:"templates"`
	Overrides []any `json:"overrides"`
}

type ExportProfile struct {
	Goal                     *string  `json:"goal"`
	Height                   *float64 `json:"height"`
	BirthDate                *string  `json:"birthDate"`
	TrainingExperience       *string  `json:"trainingExperience"`
	CurrentTrainingSplit     *string  `json:"currentTrainingSplit"`
	PreferredProgressionStyle *string `json:"preferredProgressionStyle"`
	NutritionStrategy        *string  `json:"nutritionStrategy"`
	PersistentAiContext      *string  `json:"persistentAiContext"`
}

// ExportPhoto represents a photo entry for the ZIP
type ExportPhoto struct {
	CheckInID string
	Angle     string
	Extension string
	Data      []byte
}

// CSVData represents CSV file content
type CSVData struct {
	Headers []string
	Rows    [][]string
}

// ExportArchive holds all data for archive building
type ExportArchive struct {
	Manifest        Manifest
	Data            ExportData
	SummaryMD       string
	WorkoutsCSV     CSVData
	MeasurementsCSV CSVData
	NutritionCSV    CSVData
	CardioCSV       CSVData
	Photos          []ExportPhoto
}

// BuildZIP builds an in-memory ZIP archive from ExportArchive
func (a *ExportArchive) BuildZIP() ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	// manifest.json
	manifestBytes, err := json.MarshalIndent(a.Manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("export_zip.BuildZIP: marshal manifest: %w", err)
	}
	mf, err := zw.Create("manifest.json")
	if err != nil {
		return nil, fmt.Errorf("export_zip.BuildZIP: create manifest.json: %w", err)
	}
	if _, err := mf.Write(manifestBytes); err != nil {
		return nil, fmt.Errorf("export_zip.BuildZIP: write manifest.json: %w", err)
	}

	// data.json
	dataBytes, err := json.MarshalIndent(a.Data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("export_zip.BuildZIP: marshal data.json: %w", err)
	}
	df, err := zw.Create("data.json")
	if err != nil {
		return nil, fmt.Errorf("export_zip.BuildZIP: create data.json: %w", err)
	}
	if _, err := df.Write(dataBytes); err != nil {
		return nil, fmt.Errorf("export_zip.BuildZIP: write data.json: %w", err)
	}

	// summary.md
	sf, err := zw.Create("summary.md")
	if err != nil {
		return nil, fmt.Errorf("export_zip.BuildZIP: create summary.md: %w", err)
	}
	if _, err := sf.Write([]byte(a.SummaryMD)); err != nil {
		return nil, fmt.Errorf("export_zip.BuildZIP: write summary.md: %w", err)
	}

	// CSV files
	if err := writeCSVToZip(zw, "workouts.csv", a.WorkoutsCSV); err != nil {
		return nil, fmt.Errorf("export_zip.BuildZIP: workouts.csv: %w", err)
	}
	if err := writeCSVToZip(zw, "measurements.csv", a.MeasurementsCSV); err != nil {
		return nil, fmt.Errorf("export_zip.BuildZIP: measurements.csv: %w", err)
	}
	if err := writeCSVToZip(zw, "nutrition.csv", a.NutritionCSV); err != nil {
		return nil, fmt.Errorf("export_zip.BuildZIP: nutrition.csv: %w", err)
	}
	if err := writeCSVToZip(zw, "cardio.csv", a.CardioCSV); err != nil {
		return nil, fmt.Errorf("export_zip.BuildZIP: cardio.csv: %w", err)
	}

	// photos/
	for _, photo := range a.Photos {
		photoName := fmt.Sprintf("photos/%s_%s.%s", photo.CheckInID, photo.Angle, photo.Extension)
		pf, err := zw.Create(photoName)
		if err != nil {
			return nil, fmt.Errorf("export_zip.BuildZIP: create %s: %w", photoName, err)
		}
		if _, err := pf.Write(photo.Data); err != nil {
			return nil, fmt.Errorf("export_zip.BuildZIP: write %s: %w", photoName, err)
		}
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("export_zip.BuildZIP: close zip: %w", err)
	}

	return buf.Bytes(), nil
}

func writeCSVToZip(zw *zip.Writer, name string, data CSVData) error {
	if len(data.Headers) == 0 {
		return nil
	}
	f, err := zw.Create(name)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	if err := w.Write(data.Headers); err != nil {
		return err
	}
	for _, row := range data.Rows {
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

// NewDefaultExportArchive creates an ExportArchive with no data (for empty range)
func NewDefaultExportArchive(dateRangeStart, dateRangeEnd string, profile ExportProfile) *ExportArchive {
	return &ExportArchive{
		Manifest: Manifest{
			SchemaVersion:   1,
			ExportTimestamp: time.Now().UTC().Format(time.RFC3339),
			DateRangeStart:  dateRangeStart,
			DateRangeEnd:    dateRangeEnd,
			IncludedSections: ManifestSection{
				Workouts:     false,
				Cardio:       false,
				BodyWeight:   false,
				Measurements: false,
				Nutrition:    false,
				Photos:       false,
			},
		},
		Data: ExportData{
			Workouts:          []any{},
			Cardio:            []any{},
			BodyWeightEntries: []any{},
			Measurements:      []any{},
			Nutrition: ExportNutrition{
				Products:  []any{},
				Templates: []any{},
				Overrides: []any{},
			},
			WeekFlags:   []any{},
			UserProfile: profile,
		},
		SummaryMD:       "No data recorded for this period.\n",
		WorkoutsCSV:     CSVData{Headers: []string{"date", "exercise_name", "set_number", "weight", "reps", "rpe", "rir", "set_notes", "exercise_notes", "day_notes"}},
		MeasurementsCSV: CSVData{Headers: []string{"check_in_date", "measurement_type", "side", "value", "notes"}},
		NutritionCSV:    CSVData{Headers: []string{"date", "product_name", "amount_grams", "calories", "protein", "fat", "carbs", "meal_label", "operation"}},
		CardioCSV:       CSVData{Headers: []string{"date", "type", "duration_minutes", "avg_pulse", "heart_rate_zone", "notes"}},
		Photos:          []ExportPhoto{},
	}
}

// BuildPrompt builds a plain-text prompt string for AI analysis
func BuildPrompt(profile *UserProfileExport, dateRangeStart, dateRangeEnd string, sections SectionToggles, userComment *string, weekFlags []string, dataSummary string) string {
	var b strings.Builder

	b.WriteString("You are a fitness analysis AI. Analyze the following training and health data and provide actionable insights.\n\n")

	b.WriteString("## User Context\n")
	if profile != nil {
		if profile.Goal != nil && *profile.Goal != "" {
			b.WriteString(fmt.Sprintf("- Goal: %s\n", *profile.Goal))
		}
		if profile.Height != nil && *profile.Height > 0 {
			b.WriteString(fmt.Sprintf("- Height: %.1f cm\n", *profile.Height))
		}
		if profile.BirthDate != nil && *profile.BirthDate != "" {
			b.WriteString(fmt.Sprintf("- Birth Date: %s\n", *profile.BirthDate))
		}
		if profile.TrainingExperience != nil && *profile.TrainingExperience != "" {
			b.WriteString(fmt.Sprintf("- Training Experience: %s\n", *profile.TrainingExperience))
		}
		if profile.CurrentTrainingSplit != nil && *profile.CurrentTrainingSplit != "" {
			b.WriteString(fmt.Sprintf("- Current Training Split: %s\n", *profile.CurrentTrainingSplit))
		}
		if profile.PreferredProgressionStyle != nil && *profile.PreferredProgressionStyle != "" {
			b.WriteString(fmt.Sprintf("- Preferred Progression: %s\n", *profile.PreferredProgressionStyle))
		}
		if profile.NutritionStrategy != nil && *profile.NutritionStrategy != "" {
			b.WriteString(fmt.Sprintf("- Nutrition Strategy: %s\n", *profile.NutritionStrategy))
		}
		if profile.PersistentAiContext != nil && *profile.PersistentAiContext != "" {
			b.WriteString(fmt.Sprintf("- Persistent Context: %s\n", *profile.PersistentAiContext))
		}
	}
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf("## Period\nFrom: %s To: %s\n\n", dateRangeStart, dateRangeEnd))

	b.WriteString("## Data\n")
	b.WriteString(dataSummary)
	b.WriteString("\n")

	if len(weekFlags) > 0 {
		b.WriteString("## Week Flags\n")
		for _, f := range weekFlags {
			b.WriteString(fmt.Sprintf("- %s\n", f))
		}
		b.WriteString("\n")
	}

	if userComment != nil && *userComment != "" {
		b.WriteString(fmt.Sprintf("## User Comment\n%s\n\n", *userComment))
	}

	b.WriteString("## Analysis Requests\n")
	b.WriteString("1. Analyze overall progress during this period\n")
	b.WriteString("2. Compare body weight and measurement changes\n")
	b.WriteString("3. Evaluate volume and intensity trends\n")
	b.WriteString("4. Consider RPE/RIR data and cardio performance\n")
	b.WriteString("5. Compare training changes vs body composition changes\n")
	b.WriteString("6. Provide specific actionable recommendations\n")

	return b.String()
}

// UserProfileExport is a subset of UserProfile for prompt building
type UserProfileExport struct {
	Goal                     *string
	Height                   *float64
	BirthDate                *string
	TrainingExperience       *string
	CurrentTrainingSplit     *string
	PreferredProgressionStyle *string
	NutritionStrategy        *string
	PersistentAiContext      *string
}

// SectionToggles for prompt building
type SectionToggles struct {
	Photos        bool
	Nutrition     bool
	Cardio        bool
	Measurements  bool
}