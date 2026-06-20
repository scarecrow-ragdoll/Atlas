// FILE: apps/api/internal/atlas/service/nutrition_template_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for NutritionTemplateService covering Create (upsert), GetByID (with items), GetCurrent, ListByRange, Update, Delete.
//   SCOPE: Success paths, validation errors, not-found, empty list.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mock), apps/api/internal/atlas/models.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT

package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/apps/api/internal/atlas/service"
)

type mockNutritionTemplateRepo struct {
	atlasPostgres.NutritionTemplateRepository
	upsertFn     func(ctx context.Context, userID string, weekStartDate string, title, notes *string) (*models.NutritionTemplateRecord, error)
	getByIDFn    func(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error)
	getByWeekFn  func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error)
	listByRangeFn func(ctx context.Context, userID string, startDate, endDate string) ([]models.NutritionTemplateRecord, error)
	updateFn     func(ctx context.Context, userID string, id string, title, notes *string) (*models.NutritionTemplateRecord, error)
	deleteFn     func(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error)
}

func (m *mockNutritionTemplateRepo) Upsert(ctx context.Context, userID string, weekStartDate string, title, notes *string) (*models.NutritionTemplateRecord, error) {
	return m.upsertFn(ctx, userID, weekStartDate, title, notes)
}
func (m *mockNutritionTemplateRepo) GetByID(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
	return m.getByIDFn(ctx, userID, id)
}
func (m *mockNutritionTemplateRepo) GetByWeek(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
	return m.getByWeekFn(ctx, userID, weekStartDate)
}
func (m *mockNutritionTemplateRepo) ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.NutritionTemplateRecord, error) {
	return m.listByRangeFn(ctx, userID, startDate, endDate)
}
func (m *mockNutritionTemplateRepo) Update(ctx context.Context, userID string, id string, title, notes *string) (*models.NutritionTemplateRecord, error) {
	return m.updateFn(ctx, userID, id, title, notes)
}
func (m *mockNutritionTemplateRepo) Delete(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
	return m.deleteFn(ctx, userID, id)
}

type mockNutritionTemplateItemRepo struct {
	atlasPostgres.NutritionTemplateItemRepository
	listByTemplateFn func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error)
}

func (m *mockNutritionTemplateItemRepo) ListByTemplate(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
	return m.listByTemplateFn(ctx, templateID)
}

var tmplTestRecord = &models.NutritionTemplateRecord{
	ID:            testID,
	UserID:        testUserID,
	WeekStartDate: models.MustDate("2026-06-15"),
	Title:         ptrStr("Week A"),
	Notes:         ptrStr("High protein"),
	CreatedAt:     "2026-06-20T00:00:00Z",
	UpdatedAt:     "2026-06-20T00:00:00Z",
}

var tmplItemTestRecords = []models.NutritionTemplateItemRecord{
	{
		ID: "770e8400-e29b-41d4-a716-446655440000", TemplateID: testID, ProductID: "880e8400-e29b-41d4-a716-446655440000",
		AmountGrams: 150, MealLabel: ptrStr("Lunch"), Notes: ptrStr("Grilled"),
		CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T00:00:00Z",
	},
}

func newTemplateService(tmplRepo *mockNutritionTemplateRepo, itemRepo *mockNutritionTemplateItemRepo) service.NutritionTemplateService {
	return service.NewNutritionTemplateService(tmplRepo, itemRepo, zap.NewNop())
}

// ----- Create -----

func TestNutritionTemplateService_Create_Success(t *testing.T) {
	svc := newTemplateService(&mockNutritionTemplateRepo{
		upsertFn: func(ctx context.Context, userID string, weekStartDate string, title, notes *string) (*models.NutritionTemplateRecord, error) {
			return tmplTestRecord, nil
		},
	}, &mockNutritionTemplateItemRepo{})

	tmpl, err := svc.Create(ctx, testUserID, models.CreateTemplateInput{
		WeekStartDate: models.MustDate("2026-06-15"),
		Title:         ptrStr("Week A"),
		Notes:         ptrStr("High protein"),
	})
	require.NoError(t, err)
	require.NotNil(t, tmpl)
	assert.Equal(t, "Week A", *tmpl.Title)
	assert.Equal(t, "2026-06-15", tmpl.WeekStartDate)
}

func TestNutritionTemplateService_Create_WeekEmpty(t *testing.T) {
	svc := newTemplateService(&mockNutritionTemplateRepo{}, &mockNutritionTemplateItemRepo{})
	tmpl, err := svc.Create(ctx, testUserID, models.CreateTemplateInput{
		WeekStartDate: models.Date{},
	})
	assert.ErrorIs(t, err, service.ErrTemplateWeekRequired)
	assert.Nil(t, tmpl)
}

// ----- GetByID -----

func TestNutritionTemplateService_GetByID_Success(t *testing.T) {
	svc := newTemplateService(&mockNutritionTemplateRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
			return tmplTestRecord, nil
		},
	}, &mockNutritionTemplateItemRepo{
		listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
			return tmplItemTestRecords, nil
		},
	})

	tmpl, err := svc.GetByID(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, tmpl)
	assert.Len(t, tmpl.Items, 1)
}

func TestNutritionTemplateService_GetByID_NotFound(t *testing.T) {
	svc := newTemplateService(&mockNutritionTemplateRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
			return nil, nil
		},
	}, &mockNutritionTemplateItemRepo{})

	tmpl, err := svc.GetByID(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrTemplateNotFound)
	assert.Nil(t, tmpl)
}

// ----- GetCurrent -----

func TestNutritionTemplateService_GetCurrent_Success(t *testing.T) {
	svc := newTemplateService(&mockNutritionTemplateRepo{
		getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
			return tmplTestRecord, nil
		},
	}, &mockNutritionTemplateItemRepo{
		listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
			return tmplItemTestRecords, nil
		},
	})

	tmpl, err := svc.GetCurrent(ctx, testUserID, "2026-06-15")
	require.NoError(t, err)
	require.NotNil(t, tmpl)
}

func TestNutritionTemplateService_GetCurrent_Nil(t *testing.T) {
	svc := newTemplateService(&mockNutritionTemplateRepo{
		getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
			return nil, nil
		},
	}, &mockNutritionTemplateItemRepo{})

	tmpl, err := svc.GetCurrent(ctx, testUserID, "2026-06-15")
	require.NoError(t, err)
	assert.Nil(t, tmpl)
}

// ----- ListByRange -----

func TestNutritionTemplateService_ListByRange_Success(t *testing.T) {
	svc := newTemplateService(&mockNutritionTemplateRepo{
		listByRangeFn: func(ctx context.Context, userID string, startDate, endDate string) ([]models.NutritionTemplateRecord, error) {
			return []models.NutritionTemplateRecord{*tmplTestRecord}, nil
		},
	}, &mockNutritionTemplateItemRepo{
		listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
			return tmplItemTestRecords, nil
		},
	})

	templates, err := svc.ListByRange(ctx, testUserID, "2026-06-01", "2026-06-30")
	require.NoError(t, err)
	assert.Len(t, templates, 1)
}

// ----- Update -----

func TestNutritionTemplateService_Update_Success(t *testing.T) {
	svc := newTemplateService(&mockNutritionTemplateRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
			return tmplTestRecord, nil
		},
		updateFn: func(ctx context.Context, userID string, id string, title, notes *string) (*models.NutritionTemplateRecord, error) {
			return &models.NutritionTemplateRecord{
				ID: id, UserID: userID, WeekStartDate: models.MustDate("2026-06-15"),
				Title: title, Notes: notes,
				CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T12:00:00Z",
			}, nil
		},
	}, &mockNutritionTemplateItemRepo{
		listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
			return tmplItemTestRecords, nil
		},
	})

	tmpl, err := svc.Update(ctx, testUserID, testID, models.UpdateTemplateInput{
		Title: ptrStr("Updated Week"),
	})
	require.NoError(t, err)
	require.NotNil(t, tmpl)
	assert.Equal(t, "Updated Week", *tmpl.Title)
}

func TestNutritionTemplateService_Update_NotFound(t *testing.T) {
	svc := newTemplateService(&mockNutritionTemplateRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
			return nil, nil
		},
	}, &mockNutritionTemplateItemRepo{})

	tmpl, err := svc.Update(ctx, testUserID, testID, models.UpdateTemplateInput{})
	assert.ErrorIs(t, err, service.ErrTemplateNotFound)
	assert.Nil(t, tmpl)
}

// ----- Delete -----

func TestNutritionTemplateService_Delete_Success(t *testing.T) {
	svc := newTemplateService(&mockNutritionTemplateRepo{
		deleteFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
			return tmplTestRecord, nil
		},
	}, &mockNutritionTemplateItemRepo{})

	tmpl, err := svc.Delete(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, tmpl)
	assert.Equal(t, testID, tmpl.ID)
}

func TestNutritionTemplateService_Delete_NotFound(t *testing.T) {
	svc := newTemplateService(&mockNutritionTemplateRepo{
		deleteFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
			return nil, nil
		},
	}, &mockNutritionTemplateItemRepo{})

	tmpl, err := svc.Delete(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrTemplateNotFound)
	assert.Nil(t, tmpl)
}
