// FILE: apps/api/internal/atlas/models/nutrition_daily.go
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Define factual daily nutrition log models, entry snapshots, inputs, result wrappers, and macro-total converters.
//   SCOPE: DailyNutritionLogRecord, DailyNutritionEntryRecord, public daily log/entry models, add/update inputs, repository create input, and snapshot-based total calculations; excludes legacy template and override types.
//   DEPENDS: apps/api/internal/atlas/models/nutrition.go NutritionMacros and nutrition error result types.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   DailyNutritionLogRecord - DB record type for factual daily nutrition logs.
//   DailyNutritionEntryRecord - DB record type for product snapshot food-log entries.
//   DailyNutritionEntry - Public factual food-log entry with snapshot macros.
//   DailyNutritionLog - Public factual daily log aggregate with entries and totals.
//   UpdateDailyNutritionLogNotesInput - Notes update input for a daily log.
//   AddDailyNutritionEntryInput - Service input for adding a product snapshot entry by date.
//   UpdateDailyNutritionEntryInput - Service/repository input for updating factual entry fields.
//   CreateDailyNutritionEntryRecordInput - Repository input for creating a snapshot entry.
//   DailyNutritionLogResult - Union-style result wrapper for future transport mapping.
//   DailyNutritionEntryFromRecord/DailyNutritionLogFromRecord - Record-to-model converters.
//   DailyNutritionEntriesFromRecords - Slice converter for entry records.
//   DailyNutritionEntryMacros/DailyNutritionTotalsFromEntries - Snapshot-based macro calculations.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Made daily nutrition entry update input an explicit full-replacement contract.
// END_CHANGE_SUMMARY

package models

type DailyNutritionLogRecord struct {
	ID        string
	UserID    string
	Date      Date
	Notes     *string
	CreatedAt string
	UpdatedAt string
}

type DailyNutritionEntryRecord struct {
	ID                      string
	DailyLogID              string
	ProductID               string
	ProductNameSnapshot     string
	CaloriesPer100gSnapshot float64
	ProteinPer100gSnapshot  float64
	FatPer100gSnapshot      float64
	CarbsPer100gSnapshot    float64
	AmountGrams             float64
	MealLabel               *string
	Notes                   *string
	Position                int32
	CreatedAt               string
	UpdatedAt               string
}

type DailyNutritionEntry struct {
	ID                      string          `json:"id"`
	DailyLogID              string          `json:"dailyLogId"`
	ProductID               string          `json:"productId"`
	ProductNameSnapshot     string          `json:"productNameSnapshot"`
	CaloriesPer100gSnapshot float64         `json:"caloriesPer100gSnapshot"`
	ProteinPer100gSnapshot  float64         `json:"proteinPer100gSnapshot"`
	FatPer100gSnapshot      float64         `json:"fatPer100gSnapshot"`
	CarbsPer100gSnapshot    float64         `json:"carbsPer100gSnapshot"`
	AmountGrams             float64         `json:"amountGrams"`
	MealLabel               *string         `json:"mealLabel"`
	Notes                   *string         `json:"notes"`
	Position                int32           `json:"position"`
	Macros                  NutritionMacros `json:"macros"`
	CreatedAt               string          `json:"createdAt"`
	UpdatedAt               string          `json:"updatedAt"`
}

type DailyNutritionLog struct {
	ID        string                `json:"id"`
	UserID    string                `json:"userId"`
	Date      string                `json:"date"`
	Notes     *string               `json:"notes"`
	Entries   []DailyNutritionEntry `json:"entries"`
	Totals    NutritionMacros       `json:"totals"`
	CreatedAt string                `json:"createdAt"`
	UpdatedAt string                `json:"updatedAt"`
}

type UpdateDailyNutritionLogNotesInput struct {
	Notes *string `json:"notes"`
}

type AddDailyNutritionEntryInput struct {
	Date        Date    `json:"date"`
	ProductID   string  `json:"productId"`
	AmountGrams float64 `json:"amountGrams"`
	MealLabel   *string `json:"mealLabel"`
	Notes       *string `json:"notes"`
	Position    int32   `json:"position"`
}

type UpdateDailyNutritionEntryInput struct {
	DailyLogID  string  `json:"dailyLogId"`
	AmountGrams float64 `json:"amountGrams"`
	MealLabel   *string `json:"mealLabel"`
	Notes       *string `json:"notes"`
	Position    int32   `json:"position"`
}

type CreateDailyNutritionEntryRecordInput struct {
	DailyLogID  string
	ProductID   string
	AmountGrams float64
	MealLabel   *string
	Notes       *string
	Position    int32
}

type DailyNutritionLogResult struct {
	DailyNutritionLog *DailyNutritionLog      `json:"dailyNutritionLog"`
	ValidationErr     *NutritionValidationErr `json:"validationError"`
	NotFoundErr       *NutritionNotFoundErr   `json:"notFoundError"`
	AuthErr           *NutritionAuthErr       `json:"authError"`
}

func DailyNutritionEntryFromRecord(r *DailyNutritionEntryRecord) *DailyNutritionEntry {
	if r == nil {
		return nil
	}
	entry := &DailyNutritionEntry{
		ID:                      r.ID,
		DailyLogID:              r.DailyLogID,
		ProductID:               r.ProductID,
		ProductNameSnapshot:     r.ProductNameSnapshot,
		CaloriesPer100gSnapshot: r.CaloriesPer100gSnapshot,
		ProteinPer100gSnapshot:  r.ProteinPer100gSnapshot,
		FatPer100gSnapshot:      r.FatPer100gSnapshot,
		CarbsPer100gSnapshot:    r.CarbsPer100gSnapshot,
		AmountGrams:             r.AmountGrams,
		MealLabel:               r.MealLabel,
		Notes:                   r.Notes,
		Position:                r.Position,
		CreatedAt:               r.CreatedAt,
		UpdatedAt:               r.UpdatedAt,
	}
	entry.Macros = DailyNutritionEntryMacros(*entry)
	return entry
}

func DailyNutritionEntriesFromRecords(records []DailyNutritionEntryRecord) []DailyNutritionEntry {
	out := make([]DailyNutritionEntry, len(records))
	for i := range records {
		out[i] = *DailyNutritionEntryFromRecord(&records[i])
	}
	return out
}

func DailyNutritionLogFromRecord(r *DailyNutritionLogRecord, entries []DailyNutritionEntry) *DailyNutritionLog {
	if r == nil {
		return nil
	}
	if entries == nil {
		entries = []DailyNutritionEntry{}
	}
	return &DailyNutritionLog{
		ID:        r.ID,
		UserID:    r.UserID,
		Date:      r.Date.String(),
		Notes:     r.Notes,
		Entries:   entries,
		Totals:    DailyNutritionTotalsFromEntries(entries),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func DailyNutritionEntryMacros(entry DailyNutritionEntry) NutritionMacros {
	factor := entry.AmountGrams / 100
	return NutritionMacros{
		Calories: entry.CaloriesPer100gSnapshot * factor,
		Protein:  entry.ProteinPer100gSnapshot * factor,
		Fat:      entry.FatPer100gSnapshot * factor,
		Carbs:    entry.CarbsPer100gSnapshot * factor,
	}
}

func DailyNutritionTotalsFromEntries(entries []DailyNutritionEntry) NutritionMacros {
	var totals NutritionMacros
	for _, entry := range entries {
		macros := entry.Macros
		if macros == (NutritionMacros{}) {
			macros = DailyNutritionEntryMacros(entry)
		}
		totals.Calories += macros.Calories
		totals.Protein += macros.Protein
		totals.Fat += macros.Fat
		totals.Carbs += macros.Carbs
	}
	return totals
}
