// FILE: apps/api/internal/atlas/service/daily_nutrition_legacy_resolver.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Resolve legacy weekly template plus daily override nutrition rows into factual-style daily food-log metadata.
//   SCOPE: Legacy override existence checks, deterministic ADD/REPLACE/SUBTRACT resolution, raw operation diagnostics, and unresolved classification; excludes DB schema changes and frontend editing flows.
//   DEPENDS: atlas postgres nutrition template, template item, override, override item, and product repositories; atlas/models.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   DailyNutritionLegacyResolver - Interface for legacy override detection and deterministic resolution.
//   NewDailyNutritionLegacyResolver - Creates the concrete resolver used by server wiring and template apply skips.
//   HasLegacyNutrition - Detects whether a date has pre-cutover daily override rows.
//   Resolve - Converts one legacy daily override into export-compatible factual entries or unresolved diagnostics.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Task 6 legacy override resolver and diagnostics model integration.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"fmt"
	"math"
	"strings"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
)

const legacyResolutionTolerance = 0.001

type DailyNutritionLegacyResolver interface {
	HasLegacyNutrition(ctx context.Context, userID string, date models.Date) (bool, error)
	Resolve(ctx context.Context, userID string, date models.Date) (*models.DailyNutritionLegacyResolution, error)
}

type dailyNutritionLegacyResolver struct {
	templateRepo     postgres.NutritionTemplateRepository
	templateItemRepo postgres.NutritionTemplateItemRepository
	overrideRepo     postgres.DailyNutritionOverrideRepository
	overrideItemRepo postgres.DailyNutritionOverrideItemRepository
	productRepo      postgres.NutritionProductRepository
}

func NewDailyNutritionLegacyResolver(
	templateRepo postgres.NutritionTemplateRepository,
	templateItemRepo postgres.NutritionTemplateItemRepository,
	overrideRepo postgres.DailyNutritionOverrideRepository,
	overrideItemRepo postgres.DailyNutritionOverrideItemRepository,
	productRepo postgres.NutritionProductRepository,
) DailyNutritionLegacyResolver {
	return &dailyNutritionLegacyResolver{
		templateRepo:     templateRepo,
		templateItemRepo: templateItemRepo,
		overrideRepo:     overrideRepo,
		overrideItemRepo: overrideItemRepo,
		productRepo:      productRepo,
	}
}

func (r *dailyNutritionLegacyResolver) HasLegacyNutrition(ctx context.Context, userID string, date models.Date) (bool, error) {
	if r == nil || r.overrideRepo == nil {
		return false, nil
	}
	record, err := r.overrideRepo.GetByDate(ctx, userID, date.String())
	if err != nil {
		return false, fmt.Errorf("daily_nutrition_legacy_resolver.HasLegacyNutrition: %w", err)
	}
	return record != nil, nil
}

// START_CONTRACT: Resolve
//
//	PURPOSE: Resolve one legacy override day into factual-style entries when the old macro operations are deterministic.
//	INPUTS: { userID: string - owner scope, date: models.Date - calendar day to resolve }
//	OUTPUTS: { *models.DailyNutritionLegacyResolution - nil when no legacy override exists }
//	SIDE_EFFECTS: Reads legacy repositories only.
//	LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//
// END_CONTRACT: Resolve
func (r *dailyNutritionLegacyResolver) Resolve(ctx context.Context, userID string, date models.Date) (*models.DailyNutritionLegacyResolution, error) {
	if r == nil || r.overrideRepo == nil {
		return nil, nil
	}

	override, err := r.overrideRepo.GetByDate(ctx, userID, date.String())
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_legacy_resolver.Resolve: %w", err)
	}
	if override == nil {
		return nil, nil
	}

	weekStart := legacyWeekStart(date)
	resolution := &models.DailyNutritionLegacyResolution{
		Status:           models.LegacyResolutionResolved,
		Date:             date.String(),
		WeekStartDate:    weekStart.String(),
		SourceOverrideID: override.ID,
	}
	reasons := newLegacyResolutionReasons()

	overrideItems, err := r.loadOverrideItems(ctx, override.ID, reasons)
	if err != nil {
		return nil, err
	}
	resolution.RawOperations = legacyRawOperationsFromRecords(overrideItems)

	products := make(map[string]*models.NutritionProductRecord)
	loadProduct := func(productID string) (*models.NutritionProductRecord, error) {
		if product, ok := products[productID]; ok {
			return product, nil
		}
		if r.productRepo == nil {
			reasons.add("missing product context")
			return nil, nil
		}
		product, err := r.productRepo.GetByIDIncludeInactive(ctx, userID, productID)
		if err != nil {
			return nil, fmt.Errorf("daily_nutrition_legacy_resolver.Resolve: %w", err)
		}
		if product == nil {
			reasons.add("missing product context")
			return nil, nil
		}
		products[productID] = product
		return product, nil
	}

	templateItems, templateMissing, err := r.loadTemplateItems(ctx, userID, weekStart)
	if err != nil {
		return nil, err
	}
	if templateMissing && hasBaseDependentLegacyOperation(overrideItems) {
		reasons.add("missing template context")
	}

	entries, templateMacrosByProduct, legacyTotals, err := r.buildBaseLegacyEntries(templateItems, loadProduct)
	if err != nil {
		return nil, err
	}

	replacementSeen := make(map[string]bool)
	for _, item := range overrideItems {
		product, err := loadProduct(item.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil || !product.IsActive {
			continue
		}

		overrideMacros := legacyMacrosForProduct(product, item.AmountGrams)
		switch item.Operation {
		case string(models.OperationAdd):
			addMacros(&legacyTotals, overrideMacros)
			entries = append(entries, legacyEntryFromProduct("legacy:override:"+item.ID, product, item.AmountGrams, item.MealLabel, item.Notes))
		case string(models.OperationSubtract):
			subtractMacros(&legacyTotals, overrideMacros)
			if legacyTotalAmountForProduct(entries, item.ProductID)+legacyResolutionTolerance < item.AmountGrams {
				reasons.add("subtract exceeds base amount")
				continue
			}
			entries = subtractLegacyEntryAmount(entries, item.ProductID, item.AmountGrams)
		case string(models.OperationReplace):
			if macros, exists := templateMacrosByProduct[item.ProductID]; exists {
				subtractMacros(&legacyTotals, macros)
				addMacros(&legacyTotals, overrideMacros)
			} else {
				addMacros(&legacyTotals, overrideMacros)
			}
			if replacementSeen[item.ProductID] {
				reasons.add("multiple conflicting replacements")
				continue
			}
			replacementSeen[item.ProductID] = true
			if templateMissing {
				continue
			}
			entries = removeLegacyTemplateProductEntries(entries, item.ProductID)
			entries = append(entries, legacyEntryFromProduct("legacy:override:"+item.ID, product, item.AmountGrams, item.MealLabel, item.Notes))
		default:
			reasons.add("ambiguous legacy operation")
		}
	}

	resolution.LegacyTotals = legacyTotals
	if reasons.any() {
		resolution.Status = models.LegacyResolutionUnresolved
		resolution.UnresolvedReasons = reasons.list()
		return resolution, nil
	}

	resolution.ResolvedEntries = normalizeLegacyEntryPositions(entries)
	resolution.Totals = models.DailyNutritionTotalsFromEntries(resolution.ResolvedEntries)
	if !legacyMacrosClose(resolution.Totals, resolution.LegacyTotals) {
		reasons.add("resolved totals do not match legacy macros")
		resolution.Status = models.LegacyResolutionUnresolved
		resolution.ResolvedEntries = nil
		resolution.Totals = models.NutritionMacros{}
		resolution.UnresolvedReasons = reasons.list()
	}
	return resolution, nil
}

func (r *dailyNutritionLegacyResolver) loadOverrideItems(ctx context.Context, overrideID string, reasons *legacyResolutionReasons) ([]models.DailyNutritionOverrideItemRecord, error) {
	if r.overrideItemRepo == nil {
		reasons.add("missing override item context")
		return nil, nil
	}
	items, err := r.overrideItemRepo.ListByOverride(ctx, overrideID)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_legacy_resolver.Resolve: %w", err)
	}
	return items, nil
}

func (r *dailyNutritionLegacyResolver) loadTemplateItems(ctx context.Context, userID string, weekStart models.Date) ([]models.NutritionTemplateItemRecord, bool, error) {
	if r.templateRepo == nil {
		return nil, true, nil
	}
	template, err := r.templateRepo.GetByWeek(ctx, userID, weekStart.String())
	if err != nil {
		return nil, false, fmt.Errorf("daily_nutrition_legacy_resolver.Resolve: %w", err)
	}
	if template == nil {
		return nil, true, nil
	}
	if r.templateItemRepo == nil {
		return nil, false, nil
	}
	items, err := r.templateItemRepo.ListByTemplate(ctx, template.ID)
	if err != nil {
		return nil, false, fmt.Errorf("daily_nutrition_legacy_resolver.Resolve: %w", err)
	}
	return items, false, nil
}

func (r *dailyNutritionLegacyResolver) buildBaseLegacyEntries(
	items []models.NutritionTemplateItemRecord,
	loadProduct func(productID string) (*models.NutritionProductRecord, error),
) ([]models.DailyNutritionEntry, map[string]models.NutritionMacros, models.NutritionMacros, error) {
	entries := make([]models.DailyNutritionEntry, 0, len(items))
	indexByProduct := make(map[string]int)
	templateMacrosByProduct := make(map[string]models.NutritionMacros)
	var totals models.NutritionMacros

	for _, item := range items {
		product, err := loadProduct(item.ProductID)
		if err != nil {
			return nil, nil, models.NutritionMacros{}, err
		}
		if product == nil || !product.IsActive {
			continue
		}
		macros := legacyMacrosForProduct(product, item.AmountGrams)
		addMacros(&totals, macros)
		existingMacros := templateMacrosByProduct[item.ProductID]
		addMacros(&existingMacros, macros)
		templateMacrosByProduct[item.ProductID] = existingMacros

		if index, ok := indexByProduct[item.ProductID]; ok {
			entries[index].AmountGrams += item.AmountGrams
			entries[index].MealLabel = mergeLegacyText(entries[index].MealLabel, item.MealLabel)
			entries[index].Notes = mergeLegacyText(entries[index].Notes, item.Notes)
			entries[index].Macros = models.DailyNutritionEntryMacros(entries[index])
			continue
		}
		entry := legacyEntryFromProduct("legacy:template:"+item.ProductID, product, item.AmountGrams, item.MealLabel, item.Notes)
		indexByProduct[item.ProductID] = len(entries)
		entries = append(entries, entry)
	}

	return entries, templateMacrosByProduct, totals, nil
}

func legacyWeekStart(date models.Date) models.Date {
	t := date.Time()
	offset := (int(t.Weekday()) + 6) % 7
	return models.MustDate(t.AddDate(0, 0, -offset).Format("2006-01-02"))
}

func hasBaseDependentLegacyOperation(items []models.DailyNutritionOverrideItemRecord) bool {
	for _, item := range items {
		if item.Operation == string(models.OperationReplace) || item.Operation == string(models.OperationSubtract) {
			return true
		}
	}
	return false
}

func legacyEntryFromProduct(id string, product *models.NutritionProductRecord, amountGrams float64, mealLabel, notes *string) models.DailyNutritionEntry {
	entry := models.DailyNutritionEntry{
		ID:                      id,
		ProductID:               product.ID,
		ProductNameSnapshot:     product.Name,
		CaloriesPer100gSnapshot: product.CaloriesPer100g,
		ProteinPer100gSnapshot:  product.ProteinPer100g,
		FatPer100gSnapshot:      product.FatPer100g,
		CarbsPer100gSnapshot:    product.CarbsPer100g,
		AmountGrams:             amountGrams,
		MealLabel:               mealLabel,
		Notes:                   notes,
		CreatedAt:               product.CreatedAt,
		UpdatedAt:               product.UpdatedAt,
	}
	entry.Macros = models.DailyNutritionEntryMacros(entry)
	return entry
}

func legacyMacrosForProduct(product *models.NutritionProductRecord, amountGrams float64) models.NutritionMacros {
	factor := amountGrams / 100
	return models.NutritionMacros{
		Calories: product.CaloriesPer100g * factor,
		Protein:  product.ProteinPer100g * factor,
		Fat:      product.FatPer100g * factor,
		Carbs:    product.CarbsPer100g * factor,
	}
}

func addMacros(target *models.NutritionMacros, delta models.NutritionMacros) {
	target.Calories += delta.Calories
	target.Protein += delta.Protein
	target.Fat += delta.Fat
	target.Carbs += delta.Carbs
}

func subtractMacros(target *models.NutritionMacros, delta models.NutritionMacros) {
	target.Calories -= delta.Calories
	target.Protein -= delta.Protein
	target.Fat -= delta.Fat
	target.Carbs -= delta.Carbs
}

func legacyTotalAmountForProduct(entries []models.DailyNutritionEntry, productID string) float64 {
	var amount float64
	for _, entry := range entries {
		if entry.ProductID == productID {
			amount += entry.AmountGrams
		}
	}
	return amount
}

func subtractLegacyEntryAmount(entries []models.DailyNutritionEntry, productID string, amountGrams float64) []models.DailyNutritionEntry {
	remaining := amountGrams
	out := make([]models.DailyNutritionEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.ProductID != productID || remaining <= legacyResolutionTolerance {
			out = append(out, entry)
			continue
		}
		deduct := math.Min(entry.AmountGrams, remaining)
		entry.AmountGrams -= deduct
		remaining -= deduct
		if entry.AmountGrams > legacyResolutionTolerance {
			entry.Macros = models.DailyNutritionEntryMacros(entry)
			out = append(out, entry)
		}
	}
	return out
}

func removeLegacyTemplateProductEntries(entries []models.DailyNutritionEntry, productID string) []models.DailyNutritionEntry {
	out := make([]models.DailyNutritionEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.ProductID != productID || !strings.HasPrefix(entry.ID, "legacy:template:") {
			out = append(out, entry)
		}
	}
	return out
}

func normalizeLegacyEntryPositions(entries []models.DailyNutritionEntry) []models.DailyNutritionEntry {
	out := make([]models.DailyNutritionEntry, 0, len(entries))
	for i, entry := range entries {
		entry.Position = int32(i)
		entry.Macros = models.DailyNutritionEntryMacros(entry)
		out = append(out, entry)
	}
	return out
}

func legacyRawOperationsFromRecords(records []models.DailyNutritionOverrideItemRecord) []models.DailyNutritionLegacyOperation {
	out := make([]models.DailyNutritionLegacyOperation, len(records))
	for i, record := range records {
		out[i] = models.DailyNutritionLegacyOperation{
			ID:          record.ID,
			OverrideID:  record.OverrideID,
			ProductID:   record.ProductID,
			AmountGrams: record.AmountGrams,
			Operation:   record.Operation,
			MealLabel:   record.MealLabel,
			Notes:       record.Notes,
			CreatedAt:   record.CreatedAt,
			UpdatedAt:   record.UpdatedAt,
		}
	}
	return out
}

func mergeLegacyText(current, next *string) *string {
	if current == nil || next == nil {
		return nil
	}
	if *current != *next {
		return nil
	}
	return current
}

func legacyMacrosClose(left, right models.NutritionMacros) bool {
	return math.Abs(left.Calories-right.Calories) <= legacyResolutionTolerance &&
		math.Abs(left.Protein-right.Protein) <= legacyResolutionTolerance &&
		math.Abs(left.Fat-right.Fat) <= legacyResolutionTolerance &&
		math.Abs(left.Carbs-right.Carbs) <= legacyResolutionTolerance
}

type legacyResolutionReasons struct {
	seen  map[string]bool
	items []string
}

func newLegacyResolutionReasons() *legacyResolutionReasons {
	return &legacyResolutionReasons{seen: make(map[string]bool)}
}

func (r *legacyResolutionReasons) add(reason string) {
	if r.seen[reason] {
		return
	}
	r.seen[reason] = true
	r.items = append(r.items, reason)
}

func (r *legacyResolutionReasons) any() bool {
	return len(r.items) > 0
}

func (r *legacyResolutionReasons) list() []string {
	return append([]string(nil), r.items...)
}
