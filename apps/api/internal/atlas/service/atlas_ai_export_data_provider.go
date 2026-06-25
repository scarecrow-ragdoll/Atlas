// FILE: apps/api/internal/atlas/service/atlas_ai_export_data_provider.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Assemble local Atlas repository and service data for AI export archives.
//   SCOPE: Detailed factual daily nutrition logs, weekly planned nutrition templates with product snapshots, unresolved legacy nutrition diagnostics, and empty non-nutrition placeholders until those sections get concrete providers.
//   DEPENDS: DailyNutritionLogService, DailyNutritionLegacyResolver, atlas postgres nutrition repositories, apps/api/internal/atlas/models.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   NewAtlasAiExportDataProvider - Creates the local/internal AI export data provider used by runtime wiring.
//   GetDailyNutritionExport/GetNutritionTemplateExport/GetLegacyNutritionExport - Export detailed nutrition payload slices for archive data.json and CSV flattening.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Task 11 repo/service-backed nutrition AI export provider without external AI/API calls.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"fmt"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
)

type atlasAiExportDataProvider struct {
	dailyNutrition   DailyNutritionLogService
	templateRepo     postgres.NutritionTemplateRepository
	templateItemRepo postgres.NutritionTemplateItemRepository
	productRepo      postgres.NutritionProductRepository
	legacyResolver   DailyNutritionLegacyResolver
}

func NewAtlasAiExportDataProvider(
	dailyNutrition DailyNutritionLogService,
	templateRepo postgres.NutritionTemplateRepository,
	templateItemRepo postgres.NutritionTemplateItemRepository,
	productRepo postgres.NutritionProductRepository,
	legacyResolver DailyNutritionLegacyResolver,
) AiExportDataProvider {
	return &atlasAiExportDataProvider{
		dailyNutrition:   dailyNutrition,
		templateRepo:     templateRepo,
		templateItemRepo: templateItemRepo,
		productRepo:      productRepo,
		legacyResolver:   legacyResolver,
	}
}

func (p *atlasAiExportDataProvider) GetWorkoutSummary(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *atlasAiExportDataProvider) GetCardioEntries(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *atlasAiExportDataProvider) GetBodyWeightEntries(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *atlasAiExportDataProvider) GetBodyCheckIns(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *atlasAiExportDataProvider) GetBodyMeasurements(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *atlasAiExportDataProvider) GetWeekFlags(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

// START_CONTRACT: GetDailyNutritionExport
//
//	PURPOSE: Export factual daily nutrition logs with per-entry product snapshots and calculated macros.
//	INPUTS: { userID: string - owner scope, from: models.Date - start date, to: models.Date - end date }
//	OUTPUTS: { []any - JSON-ready daily log payloads }
//	SIDE_EFFECTS: Reads DailyNutritionLogService only.
//	LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//
// END_CONTRACT: GetDailyNutritionExport
func (p *atlasAiExportDataProvider) GetDailyNutritionExport(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	if p == nil || p.dailyNutrition == nil {
		return []any{}, nil
	}
	logs, err := p.dailyNutrition.ListByRange(ctx, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("atlas_ai_export_data_provider.GetDailyNutritionExport: %w", err)
	}
	out := make([]any, 0, len(logs))
	for _, log := range logs {
		entries := make([]any, 0, len(log.Entries))
		for _, entry := range log.Entries {
			entries = append(entries, nutritionEntryExportMap(entry))
		}
		out = append(out, map[string]any{
			"id":      log.ID,
			"userId":  log.UserID,
			"date":    log.Date,
			"notes":   log.Notes,
			"totals":  nutritionMacrosExportMap(log.Totals),
			"entries": entries,
		})
	}
	return out, nil
}

// START_CONTRACT: GetNutritionTemplateExport
//
//	PURPOSE: Export weekly nutrition templates with planned product snapshot entries.
//	INPUTS: { userID: string - owner scope, from: models.Date - start date, to: models.Date - end date }
//	OUTPUTS: { []any - JSON-ready weekly template payloads }
//	SIDE_EFFECTS: Reads template, template-item, and product repositories only.
//	LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//
// END_CONTRACT: GetNutritionTemplateExport
func (p *atlasAiExportDataProvider) GetNutritionTemplateExport(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	if p == nil || p.templateRepo == nil {
		return []any{}, nil
	}
	templates, err := p.templateRepo.ListByRange(ctx, userID, from.String(), to.String())
	if err != nil {
		return nil, fmt.Errorf("atlas_ai_export_data_provider.GetNutritionTemplateExport: %w", err)
	}
	out := make([]any, 0, len(templates))
	for _, tmpl := range templates {
		plannedEntries, err := p.nutritionTemplatePlannedEntries(ctx, userID, tmpl.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, map[string]any{
			"id":             tmpl.ID,
			"userId":         tmpl.UserID,
			"weekStartDate":  tmpl.WeekStartDate.String(),
			"weekEndDate":    tmpl.WeekStartDate.Time().AddDate(0, 0, 6).Format("2006-01-02"),
			"title":          tmpl.Title,
			"notes":          tmpl.Notes,
			"plannedEntries": plannedEntries,
		})
	}
	return out, nil
}

// START_CONTRACT: GetLegacyNutritionExport
//
//	PURPOSE: Export unresolved legacy daily nutrition resolution diagnostics under a separate legacy block.
//	INPUTS: { userID: string - owner scope, from: models.Date - start date, to: models.Date - end date }
//	OUTPUTS: { []any - unresolved legacy day payloads only }
//	SIDE_EFFECTS: Reads DailyNutritionLegacyResolver only.
//	LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//
// END_CONTRACT: GetLegacyNutritionExport
func (p *atlasAiExportDataProvider) GetLegacyNutritionExport(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	if p == nil || p.legacyResolver == nil {
		return []any{}, nil
	}
	out := []any{}
	for date := from; !date.Time().After(to.Time()); date = models.MustDate(date.Time().AddDate(0, 0, 1).Format("2006-01-02")) {
		resolution, err := p.legacyResolver.Resolve(ctx, userID, date)
		if err != nil {
			return nil, fmt.Errorf("atlas_ai_export_data_provider.GetLegacyNutritionExport: %w", err)
		}
		if resolution == nil || resolution.Status != models.LegacyResolutionUnresolved {
			continue
		}
		out = append(out, legacyNutritionExportMap(resolution))
	}
	return out, nil
}

func (p *atlasAiExportDataProvider) GetProgressPhotos(ctx context.Context, userID string, from, to models.Date) ([]ExportPhoto, error) {
	return []ExportPhoto{}, nil
}

func (p *atlasAiExportDataProvider) nutritionTemplatePlannedEntries(ctx context.Context, userID string, templateID string) ([]any, error) {
	if p.templateItemRepo == nil {
		return []any{}, nil
	}
	items, err := p.templateItemRepo.ListByTemplate(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("atlas_ai_export_data_provider.GetNutritionTemplateExport: %w", err)
	}
	out := make([]any, 0, len(items))
	for _, item := range items {
		product, err := p.loadNutritionProductSnapshot(ctx, userID, item.ProductID)
		if err != nil {
			return nil, err
		}
		out = append(out, plannedNutritionEntryExportMap(item, product))
	}
	return out, nil
}

func (p *atlasAiExportDataProvider) loadNutritionProductSnapshot(ctx context.Context, userID string, productID string) (*models.NutritionProductRecord, error) {
	if p.productRepo == nil {
		return nil, nil
	}
	product, err := p.productRepo.GetByIDIncludeInactive(ctx, userID, productID)
	if err != nil {
		return nil, fmt.Errorf("atlas_ai_export_data_provider.loadNutritionProductSnapshot: %w", err)
	}
	return product, nil
}

func nutritionEntryExportMap(entry models.DailyNutritionEntry) map[string]any {
	macros := entry.Macros
	if macros == (models.NutritionMacros{}) {
		macros = models.DailyNutritionEntryMacros(entry)
	}
	return map[string]any{
		"id":                      entry.ID,
		"dailyLogId":              entry.DailyLogID,
		"productId":               entry.ProductID,
		"productNameSnapshot":     entry.ProductNameSnapshot,
		"amountGrams":             entry.AmountGrams,
		"caloriesPer100gSnapshot": entry.CaloriesPer100gSnapshot,
		"proteinPer100gSnapshot":  entry.ProteinPer100gSnapshot,
		"fatPer100gSnapshot":      entry.FatPer100gSnapshot,
		"carbsPer100gSnapshot":    entry.CarbsPer100gSnapshot,
		"entryCalories":           macros.Calories,
		"entryProtein":            macros.Protein,
		"entryFat":                macros.Fat,
		"entryCarbs":              macros.Carbs,
		"mealLabel":               entry.MealLabel,
		"notes":                   entry.Notes,
		"position":                entry.Position,
		"createdAt":               entry.CreatedAt,
		"updatedAt":               entry.UpdatedAt,
	}
}

func plannedNutritionEntryExportMap(item models.NutritionTemplateItemRecord, product *models.NutritionProductRecord) map[string]any {
	productName := ""
	var caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g float64
	if product != nil {
		productName = product.Name
		caloriesPer100g = product.CaloriesPer100g
		proteinPer100g = product.ProteinPer100g
		fatPer100g = product.FatPer100g
		carbsPer100g = product.CarbsPer100g
	}
	factor := item.AmountGrams / 100
	return map[string]any{
		"id":                      item.ID,
		"templateId":              item.TemplateID,
		"productId":               item.ProductID,
		"productNameSnapshot":     productName,
		"amountGrams":             item.AmountGrams,
		"caloriesPer100gSnapshot": caloriesPer100g,
		"proteinPer100gSnapshot":  proteinPer100g,
		"fatPer100gSnapshot":      fatPer100g,
		"carbsPer100gSnapshot":    carbsPer100g,
		"entryCalories":           caloriesPer100g * factor,
		"entryProtein":            proteinPer100g * factor,
		"entryFat":                fatPer100g * factor,
		"entryCarbs":              carbsPer100g * factor,
		"mealLabel":               item.MealLabel,
		"notes":                   item.Notes,
		"createdAt":               item.CreatedAt,
		"updatedAt":               item.UpdatedAt,
	}
}

func legacyNutritionExportMap(resolution *models.DailyNutritionLegacyResolution) map[string]any {
	rawOperations := make([]any, 0, len(resolution.RawOperations))
	for _, op := range resolution.RawOperations {
		rawOperations = append(rawOperations, map[string]any{
			"id":          op.ID,
			"overrideId":  op.OverrideID,
			"productId":   op.ProductID,
			"amountGrams": op.AmountGrams,
			"operation":   op.Operation,
			"mealLabel":   op.MealLabel,
			"notes":       op.Notes,
			"createdAt":   op.CreatedAt,
			"updatedAt":   op.UpdatedAt,
		})
	}
	return map[string]any{
		"legacyResolutionStatus": string(resolution.Status),
		"date":                   resolution.Date,
		"weekStartDate":          resolution.WeekStartDate,
		"sourceOverrideId":       resolution.SourceOverrideID,
		"legacyTotals":           nutritionMacrosExportMap(resolution.LegacyTotals),
		"rawOperations":          rawOperations,
		"unresolvedReasons":      resolution.UnresolvedReasons,
	}
}

func nutritionMacrosExportMap(macros models.NutritionMacros) map[string]any {
	return map[string]any{
		"calories": macros.Calories,
		"protein":  macros.Protein,
		"fat":      macros.Fat,
		"carbs":    macros.Carbs,
	}
}
