// FILE: apps/api/internal/atlas/models/nutrition.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define all WAVE-05 Nutrition domain types: DB records, public models, inputs, result unions, error types, enums, and converter functions.
//   SCOPE: NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, DailyNutritionOverrideItem, NutritionMacros. Inputs for all CRUD operations. Union result types with ValidationError/NotFoundError/AuthError. Error types with Error() interface. Record-to-model converters. fmtErrorAs helper (shared with cardio.go).
//   DEPENDS: models/cardio.go (for fmtErrorAs shared helper in same package).
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Operation - String enum (add/subtract/replace)
//   NutritionErrorCode - String enum (VALIDATION_ERROR/NOT_FOUND/AUTH_ERROR/INTERNAL_ERROR)
//   NutritionProductRecord - DB record type for sqlc-to-service boundary
//   NutritionTemplateRecord - DB record type for template
//   NutritionTemplateItemRecord - DB record type for template item
//   DailyNutritionOverrideRecord - DB record type for override
//   DailyNutritionOverrideItemRecord - DB record type for override item
//   NutritionProduct - Public GraphQL model
//   NutritionTemplate - Public GraphQL model with Items slice
//   NutritionTemplateItem - Public GraphQL model
//   DailyNutritionOverride - Public GraphQL model with Items slice
//   DailyNutritionOverrideItem - Public GraphQL model
//   NutritionMacros - KJBJU calculation result (calories/protein/fat/carbs)
//   CreateProductInput/UpdateProductInput - Product mutation inputs
//   CreateTemplateInput/UpdateTemplateInput - Template mutation inputs
//   CreateTemplateItemInput/UpdateTemplateItemInput - Template item mutation inputs
//   CreateOverrideInput/UpdateOverrideInput - Override mutation inputs
//   CreateOverrideItemInput/UpdateOverrideItemInput - Override item mutation inputs
//   NutritionProductResult/NutritionProductsResult - Union result types for products
//   NutritionTemplateResult/NutritionTemplatesResult - Union result types for templates
//   NutritionTemplateItemResult - Union result type for template items
//   DailyNutritionOverrideResult/DailyNutritionOverridesResult - Union result types for overrides
//   DailyNutritionOverrideItemResult - Union result type for override items
//   NutritionMacrosResult - Union result type for macro calculation
//   NutritionValidationErr/NutritionNotFoundErr/NutritionAuthErr - Error types with Error() interface
//   NutritionProductFromRecord/NutritionTemplateFromRecord/etc - Record-to-model converters
//   NutritionTemplateItemsFromRecords/NutritionOverrideItemsFromRecords - Slice converters
//   NutritionProductResultFromError - Error-to-result mapper
// END_MODULE_MAP

package models

type Operation string

const (
	OperationAdd      Operation = "add"
	OperationSubtract Operation = "subtract"
	OperationReplace  Operation = "replace"
)

type NutritionErrorCode string

const (
	NutritionErrorValidation NutritionErrorCode = "VALIDATION_ERROR"
	NutritionErrorNotFound   NutritionErrorCode = "NOT_FOUND"
	NutritionErrorAuth       NutritionErrorCode = "AUTH_ERROR"
	NutritionErrorInternal   NutritionErrorCode = "INTERNAL_ERROR"
)

type NutritionProductRecord struct {
	ID              string
	UserID          string
	Name            string
	CaloriesPer100g float64
	ProteinPer100g  float64
	FatPer100g      float64
	CarbsPer100g    float64
	Notes           *string
	IsActive        bool
	CreatedAt       string
	UpdatedAt       string
}

type NutritionTemplateRecord struct {
	ID            string
	UserID        string
	WeekStartDate Date
	Title         *string
	Notes         *string
	CreatedAt     string
	UpdatedAt     string
}

type NutritionTemplateItemRecord struct {
	ID          string
	TemplateID  string
	ProductID   string
	AmountGrams float64
	MealLabel   *string
	Notes       *string
	CreatedAt   string
	UpdatedAt   string
}

type DailyNutritionOverrideRecord struct {
	ID        string
	UserID    string
	Date      Date
	Notes     *string
	CreatedAt string
	UpdatedAt string
}

type DailyNutritionOverrideItemRecord struct {
	ID          string
	OverrideID  string
	ProductID   string
	AmountGrams float64
	Operation   string
	MealLabel   *string
	Notes       *string
	CreatedAt   string
	UpdatedAt   string
}

type NutritionProduct struct {
	ID              string   `json:"id"`
	UserID          string   `json:"userId"`
	Name            string   `json:"name"`
	CaloriesPer100g float64  `json:"caloriesPer100g"`
	ProteinPer100g  float64  `json:"proteinPer100g"`
	FatPer100g      float64  `json:"fatPer100g"`
	CarbsPer100g    float64  `json:"carbsPer100g"`
	Notes           *string  `json:"notes"`
	IsActive        bool     `json:"isActive"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}

type NutritionTemplate struct {
	ID            string                   `json:"id"`
	UserID        string                   `json:"userId"`
	WeekStartDate string                   `json:"weekStartDate"`
	Title         *string                  `json:"title"`
	Notes         *string                  `json:"notes"`
	Items         []NutritionTemplateItem  `json:"items"`
	CreatedAt     string                   `json:"createdAt"`
	UpdatedAt     string                   `json:"updatedAt"`
}

type NutritionTemplateItem struct {
	ID          string   `json:"id"`
	TemplateID  string   `json:"templateId"`
	ProductID   string   `json:"productId"`
	AmountGrams float64  `json:"amountGrams"`
	MealLabel   *string  `json:"mealLabel"`
	Notes       *string  `json:"notes"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

type DailyNutritionOverride struct {
	ID        string                       `json:"id"`
	UserID    string                       `json:"userId"`
	Date      string                       `json:"date"`
	Notes     *string                      `json:"notes"`
	Items     []DailyNutritionOverrideItem `json:"items"`
	CreatedAt string                       `json:"createdAt"`
	UpdatedAt string                       `json:"updatedAt"`
}

type DailyNutritionOverrideItem struct {
	ID          string    `json:"id"`
	OverrideID  string    `json:"overrideId"`
	ProductID   string    `json:"productId"`
	AmountGrams float64   `json:"amountGrams"`
	Operation   Operation `json:"operation"`
	MealLabel   *string   `json:"mealLabel"`
	Notes       *string   `json:"notes"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}

type NutritionMacros struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
}

type CreateProductInput struct {
	Name            string   `json:"name"`
	CaloriesPer100g float64  `json:"caloriesPer100g"`
	ProteinPer100g  float64  `json:"proteinPer100g"`
	FatPer100g      float64  `json:"fatPer100g"`
	CarbsPer100g    float64  `json:"carbsPer100g"`
	Notes           *string  `json:"notes"`
}

type UpdateProductInput struct {
	Name            *string  `json:"name"`
	CaloriesPer100g *float64 `json:"caloriesPer100g"`
	ProteinPer100g  *float64 `json:"proteinPer100g"`
	FatPer100g      *float64 `json:"fatPer100g"`
	CarbsPer100g    *float64 `json:"carbsPer100g"`
	Notes           *string  `json:"notes"`
}

type CreateTemplateInput struct {
	WeekStartDate Date    `json:"weekStartDate"`
	Title         *string `json:"title"`
	Notes         *string `json:"notes"`
}

type UpdateTemplateInput struct {
	Title *string `json:"title"`
	Notes *string `json:"notes"`
}

type CreateTemplateItemInput struct {
	TemplateID  string   `json:"templateId"`
	ProductID   string   `json:"productId"`
	AmountGrams float64  `json:"amountGrams"`
	MealLabel   *string  `json:"mealLabel"`
	Notes       *string  `json:"notes"`
}

type UpdateTemplateItemInput struct {
	AmountGrams *float64 `json:"amountGrams"`
	MealLabel   *string  `json:"mealLabel"`
	Notes       *string  `json:"notes"`
}

type CreateOverrideInput struct {
	Date  Date     `json:"date"`
	Notes *string  `json:"notes"`
}

type UpdateOverrideInput struct {
	Notes *string `json:"notes"`
}

type CreateOverrideItemInput struct {
	OverrideID  string    `json:"overrideId"`
	ProductID   string    `json:"productId"`
	AmountGrams float64   `json:"amountGrams"`
	Operation   Operation `json:"operation"`
	MealLabel   *string   `json:"mealLabel"`
	Notes       *string   `json:"notes"`
}

type UpdateOverrideItemInput struct {
	AmountGrams *float64  `json:"amountGrams"`
	Operation   *Operation `json:"operation"`
	MealLabel   *string   `json:"mealLabel"`
	Notes       *string   `json:"notes"`
}

type NutritionProductResult struct {
	NutritionProduct *NutritionProduct      `json:"nutritionProduct"`
	ValidationErr    *NutritionValidationErr `json:"validationError"`
	NotFoundErr      *NutritionNotFoundErr   `json:"notFoundError"`
	AuthErr          *NutritionAuthErr       `json:"authError"`
}

type NutritionProductsResult struct {
	Products       []NutritionProduct      `json:"products"`
	ValidationErr  *NutritionValidationErr  `json:"validationError"`
	AuthErr        *NutritionAuthErr        `json:"authError"`
}

type NutritionTemplateResult struct {
	NutritionTemplate *NutritionTemplate     `json:"nutritionTemplate"`
	ValidationErr     *NutritionValidationErr `json:"validationError"`
	NotFoundErr       *NutritionNotFoundErr    `json:"notFoundError"`
	AuthErr           *NutritionAuthErr        `json:"authError"`
}

type NutritionTemplatesResult struct {
	Templates      []NutritionTemplate     `json:"templates"`
	ValidationErr  *NutritionValidationErr  `json:"validationError"`
	AuthErr        *NutritionAuthErr        `json:"authError"`
}

type NutritionTemplateItemResult struct {
	NutritionTemplateItem *NutritionTemplateItem  `json:"nutritionTemplateItem"`
	ValidationErr         *NutritionValidationErr  `json:"validationError"`
	NotFoundErr           *NutritionNotFoundErr    `json:"notFoundError"`
	AuthErr               *NutritionAuthErr        `json:"authError"`
}

type DailyNutritionOverrideResult struct {
	DailyNutritionOverride *DailyNutritionOverride `json:"dailyNutritionOverride"`
	ValidationErr          *NutritionValidationErr  `json:"validationError"`
	NotFoundErr            *NutritionNotFoundErr    `json:"notFoundError"`
	AuthErr                *NutritionAuthErr        `json:"authError"`
}

type DailyNutritionOverridesResult struct {
	Overrides      []DailyNutritionOverride `json:"overrides"`
	ValidationErr  *NutritionValidationErr   `json:"validationError"`
	AuthErr        *NutritionAuthErr         `json:"authError"`
}

type DailyNutritionOverrideItemResult struct {
	DailyNutritionOverrideItem *DailyNutritionOverrideItem `json:"dailyNutritionOverrideItem"`
	ValidationErr              *NutritionValidationErr      `json:"validationError"`
	NotFoundErr                *NutritionNotFoundErr        `json:"notFoundError"`
	AuthErr                    *NutritionAuthErr            `json:"authError"`
}

type NutritionMacrosResult struct {
	Macros         *NutritionMacros         `json:"macros"`
	ValidationErr  *NutritionValidationErr  `json:"validationError"`
	AuthErr        *NutritionAuthErr        `json:"authError"`
}

type NutritionValidationErr struct {
	Message string             `json:"message"`
	Code    NutritionErrorCode `json:"code"`
}

func (e *NutritionValidationErr) Error() string {
	if e == nil || e.Message == "" {
		return "nutrition validation error"
	}
	return e.Message
}

type NutritionNotFoundErr struct {
	Message string             `json:"message"`
	Code    NutritionErrorCode `json:"code"`
}

func (e *NutritionNotFoundErr) Error() string {
	if e == nil || e.Message == "" {
		return "nutrition not found"
	}
	return e.Message
}

type NutritionAuthErr struct {
	Message string             `json:"message"`
	Code    NutritionErrorCode `json:"code"`
}

func (e *NutritionAuthErr) Error() string {
	if e == nil || e.Message == "" {
		return "nutrition auth error"
	}
	return e.Message
}

func NutritionProductFromRecord(r *NutritionProductRecord) *NutritionProduct {
	if r == nil {
		return nil
	}
	return &NutritionProduct{
		ID:              r.ID,
		UserID:          r.UserID,
		Name:            r.Name,
		CaloriesPer100g: r.CaloriesPer100g,
		ProteinPer100g:  r.ProteinPer100g,
		FatPer100g:      r.FatPer100g,
		CarbsPer100g:    r.CarbsPer100g,
		Notes:           r.Notes,
		IsActive:        r.IsActive,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

func NutritionTemplateFromRecord(r *NutritionTemplateRecord, items []NutritionTemplateItem) *NutritionTemplate {
	if r == nil {
		return nil
	}
	if items == nil {
		items = []NutritionTemplateItem{}
	}
	return &NutritionTemplate{
		ID:            r.ID,
		UserID:        r.UserID,
		WeekStartDate: r.WeekStartDate.String(),
		Title:         r.Title,
		Notes:         r.Notes,
		Items:         items,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

func NutritionTemplateItemFromRecord(r *NutritionTemplateItemRecord) *NutritionTemplateItem {
	if r == nil {
		return nil
	}
	return &NutritionTemplateItem{
		ID:          r.ID,
		TemplateID:  r.TemplateID,
		ProductID:   r.ProductID,
		AmountGrams: r.AmountGrams,
		MealLabel:   r.MealLabel,
		Notes:       r.Notes,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func NutritionTemplateItemsFromRecords(records []NutritionTemplateItemRecord) []NutritionTemplateItem {
	out := make([]NutritionTemplateItem, len(records))
	for i := range records {
		out[i] = *NutritionTemplateItemFromRecord(&records[i])
	}
	return out
}

func DailyNutritionOverrideFromRecord(r *DailyNutritionOverrideRecord, items []DailyNutritionOverrideItem) *DailyNutritionOverride {
	if r == nil {
		return nil
	}
	if items == nil {
		items = []DailyNutritionOverrideItem{}
	}
	return &DailyNutritionOverride{
		ID:        r.ID,
		UserID:    r.UserID,
		Date:      r.Date.String(),
		Notes:     r.Notes,
		Items:     items,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func DailyNutritionOverrideItemFromRecord(r *DailyNutritionOverrideItemRecord) *DailyNutritionOverrideItem {
	if r == nil {
		return nil
	}
	op := Operation(r.Operation)
	return &DailyNutritionOverrideItem{
		ID:          r.ID,
		OverrideID:  r.OverrideID,
		ProductID:   r.ProductID,
		AmountGrams: r.AmountGrams,
		Operation:   op,
		MealLabel:   r.MealLabel,
		Notes:       r.Notes,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func DailyNutritionOverrideItemsFromRecords(records []DailyNutritionOverrideItemRecord) []DailyNutritionOverrideItem {
	out := make([]DailyNutritionOverrideItem, len(records))
	for i := range records {
		out[i] = *DailyNutritionOverrideItemFromRecord(&records[i])
	}
	return out
}

func NutritionProductResultFromError(err error) *NutritionProductResult {
	if err == nil {
		return nil
	}
	var validationErr *NutritionValidationErr
	if fmtErrorAs(err, &validationErr) {
		return &NutritionProductResult{ValidationErr: validationErr}
	}
	var notFoundErr *NutritionNotFoundErr
	if fmtErrorAs(err, &notFoundErr) {
		return &NutritionProductResult{NotFoundErr: notFoundErr}
	}
	var authErr *NutritionAuthErr
	if fmtErrorAs(err, &authErr) {
		return &NutritionProductResult{AuthErr: authErr}
	}
	return nil
}