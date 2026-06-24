// FILE: apps/api/internal/atlas/service/nutrition_template_apply_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for applying weekly nutrition templates into factual daily food logs.
//   SCOPE: seed_empty_days creation, idempotency, legacy-date skips, product conflict skips, empty-template creation, retry atomicity, concurrency, and unsupported-mode validation.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres interfaces, apps/api/internal/atlas/models.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.2 - Kept legacy resolver mock aligned with Task 6 Resolve contract.
//   LAST_CHANGE: 1.0.1 - Added quality-review coverage for product-conflict and empty-template apply semantics.
// END_CHANGE_SUMMARY

package service_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/apps/api/internal/atlas/service"
)

const (
	applyProductChicken = "880e8400-e29b-41d4-a716-446655440000"
	applyProductRice    = "990e8400-e29b-41d4-a716-446655440000"
)

var applyTemplateItems = []models.NutritionTemplateItemRecord{
	{
		ID: "770e8400-e29b-41d4-a716-446655440000", TemplateID: testID, ProductID: applyProductChicken,
		AmountGrams: 150, MealLabel: ptrStr("Lunch"), Notes: ptrStr("grilled"),
		CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T00:00:00Z",
	},
	{
		ID: "770e8400-e29b-41d4-a716-446655440001", TemplateID: testID, ProductID: applyProductRice,
		AmountGrams: 120, MealLabel: ptrStr("Lunch"), Notes: nil,
		CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T00:00:00Z",
	},
}

type mockTemplateApplyProductRepo struct {
	atlasPostgres.NutritionProductRepository
	products map[string]*models.NutritionProductRecord
}

func (m *mockTemplateApplyProductRepo) GetByID(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
	product := m.products[id]
	if product == nil {
		return nil, nil
	}
	return product, nil
}

type mockTemplateApplyLegacyResolver struct {
	legacyDates map[string]bool
}

func (m *mockTemplateApplyLegacyResolver) HasLegacyNutrition(ctx context.Context, userID string, date models.Date) (bool, error) {
	return m.legacyDates[date.String()], nil
}

func (m *mockTemplateApplyLegacyResolver) Resolve(ctx context.Context, userID string, date models.Date) (*models.DailyNutritionLegacyResolution, error) {
	return nil, nil
}

type memorySeedDailyRepo struct {
	atlasPostgres.DailyNutritionLogRepository

	mu      sync.Mutex
	entries map[string][]models.DailyNutritionSeedEntryInput
	calls   int
}

func newMemorySeedDailyRepo() *memorySeedDailyRepo {
	return &memorySeedDailyRepo{entries: map[string][]models.DailyNutritionSeedEntryInput{}}
}

func (r *memorySeedDailyRepo) SeedEntriesIfEmpty(ctx context.Context, userID string, date models.Date, items []models.DailyNutritionSeedEntryInput) (*models.DailyNutritionSeedResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.calls++
	key := date.String()
	if existing := r.entries[key]; len(existing) > 0 {
		return &models.DailyNutritionSeedResult{Created: false, EntryCount: int32(len(existing))}, nil
	}
	copied := append([]models.DailyNutritionSeedEntryInput(nil), items...)
	r.entries[key] = copied
	return &models.DailyNutritionSeedResult{Created: true, EntryCount: int32(len(copied))}, nil
}

func (r *memorySeedDailyRepo) entryCount(date string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries[date])
}

func (r *memorySeedDailyRepo) totalEntryCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	total := 0
	for _, entries := range r.entries {
		total += len(entries)
	}
	return total
}

type partialFailureDailyRepo struct {
	atlasPostgres.DailyNutritionLogRepository

	mu          sync.Mutex
	entries     map[string][]models.DailyNutritionSeedEntryInput
	failOnceFor string
	failed      bool
}

func newPartialFailureDailyRepo(date string) *partialFailureDailyRepo {
	return &partialFailureDailyRepo{entries: map[string][]models.DailyNutritionSeedEntryInput{}, failOnceFor: date}
}

func (r *partialFailureDailyRepo) SeedEntriesIfEmpty(ctx context.Context, userID string, date models.Date, items []models.DailyNutritionSeedEntryInput) (*models.DailyNutritionSeedResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := date.String()
	if existing := r.entries[key]; len(existing) > 0 {
		return &models.DailyNutritionSeedResult{Created: false, EntryCount: int32(len(existing))}, nil
	}
	if key == r.failOnceFor && !r.failed {
		r.failed = true
		r.entries[key] = append([]models.DailyNutritionSeedEntryInput(nil), items[:1]...)
		delete(r.entries, key)
		return nil, errors.New("injected seed failure after first entry")
	}
	copied := append([]models.DailyNutritionSeedEntryInput(nil), items...)
	r.entries[key] = copied
	return &models.DailyNutritionSeedResult{Created: true, EntryCount: int32(len(copied))}, nil
}

func (r *partialFailureDailyRepo) entryCount(date string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries[date])
}

func newTemplateApplyService(dailyRepo atlasPostgres.DailyNutritionLogRepository, legacy service.DailyNutritionLegacyResolver) service.NutritionTemplateApplyService {
	return newTemplateApplyServiceWithItemsAndProducts(dailyRepo, legacy, applyTemplateItems, map[string]*models.NutritionProductRecord{
		applyProductChicken: {ID: applyProductChicken, UserID: testUserID, Name: "Chicken", IsActive: true},
		applyProductRice:    {ID: applyProductRice, UserID: testUserID, Name: "Rice", IsActive: true},
	})
}

func newTemplateApplyServiceWithItemsAndProducts(
	dailyRepo atlasPostgres.DailyNutritionLogRepository,
	legacy service.DailyNutritionLegacyResolver,
	items []models.NutritionTemplateItemRecord,
	products map[string]*models.NutritionProductRecord,
) service.NutritionTemplateApplyService {
	return service.NewNutritionTemplateApplyService(
		&mockNutritionTemplateRepo{
			getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
				return tmplTestRecord, nil
			},
		},
		&mockNutritionTemplateItemRepo{
			listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
				return items, nil
			},
		},
		&mockTemplateApplyProductRepo{products: products},
		dailyRepo,
		legacy,
		zap.NewNop(),
	)
}

func assertTemplateApplyProductConflict(t *testing.T, result *models.NutritionTemplateApplyResult, dailyRepo *memorySeedDailyRepo) {
	t.Helper()

	require.NotNil(t, result)
	assert.Len(t, result.Dates, 7)
	assert.Equal(t, 0, result.CreatedCount())
	assert.Equal(t, 0, result.SkippedCount())
	assert.Equal(t, 7, result.ConflictCount())
	assert.Equal(t, 0, dailyRepo.calls)
	assert.Equal(t, 0, dailyRepo.totalEntryCount())
	for _, date := range result.Dates {
		assert.Equal(t, models.ApplyDateConflict, date.Status)
		assert.Equal(t, int32(0), date.EntryCount)
		require.NotNil(t, date.Reason)
		assert.Equal(t, "template product missing or inactive", *date.Reason)
	}
}

func TestNutritionTemplateApplyService_SeedEmptyDaysCreatesSevenEmptyDates(t *testing.T) {
	dailyRepo := newMemorySeedDailyRepo()
	svc := newTemplateApplyService(dailyRepo, &mockTemplateApplyLegacyResolver{legacyDates: map[string]bool{}})

	result, err := svc.ApplyToWeek(ctx, testUserID, testID, models.ApplyModeSeedEmptyDays)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "2026-06-15", result.WeekStartDate)
	assert.Equal(t, "2026-06-21", result.WeekEndDate)
	assert.Len(t, result.Dates, 7)
	assert.Equal(t, 7, result.CreatedCount())
	assert.Equal(t, 0, result.SkippedCount())
	assert.Equal(t, 0, result.ConflictCount())
	for _, date := range result.Dates {
		assert.Equal(t, models.ApplyDateCreated, date.Status)
		assert.Equal(t, int32(2), date.EntryCount)
		assert.Nil(t, date.Reason)
		assert.Equal(t, 2, dailyRepo.entryCount(date.Date))
	}
}

func TestNutritionTemplateApplyService_SeedEmptyDaysSkipsNonEmptyDateWithoutDuplicates(t *testing.T) {
	dailyRepo := newMemorySeedDailyRepo()
	dailyRepo.entries["2026-06-16"] = []models.DailyNutritionSeedEntryInput{{ProductID: applyProductChicken, AmountGrams: 80}}
	svc := newTemplateApplyService(dailyRepo, &mockTemplateApplyLegacyResolver{legacyDates: map[string]bool{}})

	first, err := svc.ApplyToWeek(ctx, testUserID, testID, models.ApplyModeSeedEmptyDays)
	require.NoError(t, err)
	second, err := svc.ApplyToWeek(ctx, testUserID, testID, models.ApplyModeSeedEmptyDays)
	require.NoError(t, err)

	assert.Equal(t, 6, first.CreatedCount())
	assert.Equal(t, 1, first.SkippedCount())
	assert.Equal(t, 0, second.CreatedCount())
	assert.Equal(t, 7, second.SkippedCount())
	assert.Equal(t, 1, dailyRepo.entryCount("2026-06-16"))
	assert.Equal(t, 13, dailyRepo.totalEntryCount())

	firstSkipped := first.Dates[1]
	require.NotNil(t, firstSkipped.Reason)
	assert.Equal(t, models.ApplyDateSkipped, firstSkipped.Status)
	assert.Equal(t, "day has entries", *firstSkipped.Reason)
	secondSkipped := second.Dates[1]
	require.NotNil(t, secondSkipped.Reason)
	assert.Equal(t, models.ApplyDateSkipped, secondSkipped.Status)
	assert.Equal(t, "day has entries", *secondSkipped.Reason)
}

func TestNutritionTemplateApplyService_SkipsLegacyOnlyDatesUntilCutover(t *testing.T) {
	dailyRepo := newMemorySeedDailyRepo()
	svc := newTemplateApplyService(dailyRepo, &mockTemplateApplyLegacyResolver{legacyDates: map[string]bool{"2026-06-17": true}})

	result, err := svc.ApplyToWeek(ctx, testUserID, testID, models.ApplyModeSeedEmptyDays)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 6, result.CreatedCount())
	assert.Equal(t, 1, result.SkippedCount())
	assert.Equal(t, 0, dailyRepo.entryCount("2026-06-17"))

	legacyDate := result.Dates[2]
	require.NotNil(t, legacyDate.Reason)
	assert.Equal(t, models.ApplyDateSkipped, legacyDate.Status)
	assert.Equal(t, "legacy nutrition exists; migrate or review before seeding", *legacyDate.Reason)
}

func TestNutritionTemplateApplyService_MissingProductReturnsConflictsWithoutSeeding(t *testing.T) {
	dailyRepo := newMemorySeedDailyRepo()
	svc := newTemplateApplyServiceWithItemsAndProducts(
		dailyRepo,
		&mockTemplateApplyLegacyResolver{legacyDates: map[string]bool{}},
		applyTemplateItems,
		map[string]*models.NutritionProductRecord{
			applyProductChicken: {ID: applyProductChicken, UserID: testUserID, Name: "Chicken", IsActive: true},
		},
	)

	result, err := svc.ApplyToWeek(ctx, testUserID, testID, models.ApplyModeSeedEmptyDays)

	require.NoError(t, err)
	assertTemplateApplyProductConflict(t, result, dailyRepo)
}

func TestNutritionTemplateApplyService_InactiveProductReturnsConflictsWithoutSeeding(t *testing.T) {
	dailyRepo := newMemorySeedDailyRepo()
	svc := newTemplateApplyServiceWithItemsAndProducts(
		dailyRepo,
		&mockTemplateApplyLegacyResolver{legacyDates: map[string]bool{}},
		applyTemplateItems,
		map[string]*models.NutritionProductRecord{
			applyProductChicken: {ID: applyProductChicken, UserID: testUserID, Name: "Chicken", IsActive: true},
			applyProductRice:    {ID: applyProductRice, UserID: testUserID, Name: "Rice", IsActive: false},
		},
	)

	result, err := svc.ApplyToWeek(ctx, testUserID, testID, models.ApplyModeSeedEmptyDays)

	require.NoError(t, err)
	assertTemplateApplyProductConflict(t, result, dailyRepo)
}

func TestNutritionTemplateApplyService_EmptyTemplateCreatesSevenEmptyDates(t *testing.T) {
	dailyRepo := newMemorySeedDailyRepo()
	svc := newTemplateApplyServiceWithItemsAndProducts(
		dailyRepo,
		&mockTemplateApplyLegacyResolver{legacyDates: map[string]bool{}},
		[]models.NutritionTemplateItemRecord{},
		map[string]*models.NutritionProductRecord{},
	)

	result, err := svc.ApplyToWeek(ctx, testUserID, testID, models.ApplyModeSeedEmptyDays)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Dates, 7)
	assert.Equal(t, 7, result.CreatedCount())
	assert.Equal(t, 0, result.SkippedCount())
	assert.Equal(t, 0, result.ConflictCount())
	assert.Equal(t, 7, dailyRepo.calls)
	assert.Equal(t, 0, dailyRepo.totalEntryCount())
	for _, date := range result.Dates {
		assert.Equal(t, models.ApplyDateCreated, date.Status)
		assert.Equal(t, int32(0), date.EntryCount)
		assert.Nil(t, date.Reason)
	}
}

func TestNutritionTemplateApplyService_RetryAfterPartialFailureDoesNotPreserveIncompleteSeed(t *testing.T) {
	dailyRepo := newPartialFailureDailyRepo("2026-06-15")
	svc := newTemplateApplyService(dailyRepo, &mockTemplateApplyLegacyResolver{legacyDates: map[string]bool{}})

	result, err := svc.ApplyToWeek(ctx, testUserID, testID, models.ApplyModeSeedEmptyDays)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, dailyRepo.entryCount("2026-06-15"))

	result, err = svc.ApplyToWeek(ctx, testUserID, testID, models.ApplyModeSeedEmptyDays)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 7, result.CreatedCount())
	assert.Equal(t, 2, dailyRepo.entryCount("2026-06-15"))
}

func TestNutritionTemplateApplyService_ConcurrentSeedDoesNotDuplicateEntries(t *testing.T) {
	dailyRepo := newMemorySeedDailyRepo()
	svc := newTemplateApplyService(dailyRepo, &mockTemplateApplyLegacyResolver{legacyDates: map[string]bool{}})

	var wg sync.WaitGroup
	results := make([]*models.NutritionTemplateApplyResult, 2)
	errs := make([]error, 2)
	for i := range results {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx], errs[idx] = svc.ApplyToWeek(ctx, testUserID, testID, models.ApplyModeSeedEmptyDays)
		}(i)
	}
	wg.Wait()

	require.NoError(t, errs[0])
	require.NoError(t, errs[1])
	require.NotNil(t, results[0])
	require.NotNil(t, results[1])
	assert.Equal(t, 7, results[0].CreatedCount()+results[1].CreatedCount())
	assert.Equal(t, 7, results[0].SkippedCount()+results[1].SkippedCount())
	for _, date := range []string{"2026-06-15", "2026-06-16", "2026-06-17", "2026-06-18", "2026-06-19", "2026-06-20", "2026-06-21"} {
		assert.Equal(t, 2, dailyRepo.entryCount(date), date)
	}
}

func TestNutritionTemplateApplyService_UnsupportedModeDoesNotWrite(t *testing.T) {
	dailyRepo := newMemorySeedDailyRepo()
	svc := newTemplateApplyService(dailyRepo, &mockTemplateApplyLegacyResolver{legacyDates: map[string]bool{}})

	result, err := svc.ApplyToWeek(ctx, testUserID, testID, models.NutritionTemplateApplyMode("replace_week"))

	assert.ErrorIs(t, err, service.ErrNutritionTemplateApplyModeUnsupported)
	assert.Nil(t, result)
	assert.Equal(t, 0, dailyRepo.calls)
}
