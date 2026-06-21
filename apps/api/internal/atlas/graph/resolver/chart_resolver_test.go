package resolver_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/models"

	atlasSvc "monorepo-template/apps/api/internal/atlas/service"

	"monorepo-template/apps/api/internal/atlas/graph/resolver"
)

type mockBodyChartService struct {
	bodyWeightTrendFn      func(ctx context.Context, userID string, fromDate, toDate models.Date) ([]models.BodyWeightSeriesPoint, error)
	measurementTrendFn     func(ctx context.Context, userID string, measurementType string, fromDate, toDate models.Date) ([]models.MeasurementTrendPoint, error)
	measurementOverlayFn   func(ctx context.Context, userID string, measurementTypes []string, fromDate, toDate models.Date) ([]models.MeasurementOverlayGroup, error)
}

func (m *mockBodyChartService) BodyWeightTrend(ctx context.Context, userID string, fromDate, toDate models.Date) ([]models.BodyWeightSeriesPoint, error) {
	return m.bodyWeightTrendFn(ctx, userID, fromDate, toDate)
}

func (m *mockBodyChartService) MeasurementTrend(ctx context.Context, userID string, measurementType string, fromDate, toDate models.Date) ([]models.MeasurementTrendPoint, error) {
	return m.measurementTrendFn(ctx, userID, measurementType, fromDate, toDate)
}

func (m *mockBodyChartService) MeasurementOverlay(ctx context.Context, userID string, measurementTypes []string, fromDate, toDate models.Date) ([]models.MeasurementOverlayGroup, error) {
	return m.measurementOverlayFn(ctx, userID, measurementTypes, fromDate, toDate)
}

type mockNutritionWeeklyAvgService struct {
	weeklyAveragesFn func(ctx context.Context, userID string, fromDate, toDate models.Date) ([]models.NutritionWeeklyAverage, error)
}

func (m *mockNutritionWeeklyAvgService) WeeklyAverages(ctx context.Context, userID string, fromDate, toDate models.Date) ([]models.NutritionWeeklyAverage, error) {
	return m.weeklyAveragesFn(ctx, userID, fromDate, toDate)
}

// --- Body Weight Trend ---

func TestBodyWeightTrend_Success(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService: &mockBodyChartService{
			bodyWeightTrendFn: func(_ context.Context, _ string, _, _ models.Date) ([]models.BodyWeightSeriesPoint, error) {
				return []models.BodyWeightSeriesPoint{
					{Date: models.MustDate("2026-03-01"), Weight: 80, Source: models.BodyWeightSourceScale},
				}, nil
			},
		},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{},
	}

	result, err := r.GetBodyWeightTrend(userCtx("test-uid"), nil, nil)
	require.NoError(t, err)
	assert.Nil(t, result.AuthErr)
	assert.Nil(t, result.ValidationErr)
	assert.Len(t, result.Series, 1)
	assert.Equal(t, 80.0, result.Series[0].Weight)
}

func TestBodyWeightTrend_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService:             &mockBodyChartService{},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{},
	}

	result, err := r.GetBodyWeightTrend(context.Background(), nil, nil)
	require.NoError(t, err)
	assert.NotNil(t, result.AuthErr)
	assert.Equal(t, models.ChartErrorAuth, result.AuthErr.Code)
}

func TestBodyWeightTrend_EmptySeries(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService: &mockBodyChartService{
			bodyWeightTrendFn: func(_ context.Context, _ string, _, _ models.Date) ([]models.BodyWeightSeriesPoint, error) {
				return []models.BodyWeightSeriesPoint{}, nil
			},
		},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{},
	}

	result, err := r.GetBodyWeightTrend(userCtx("uid"), nil, nil)
	require.NoError(t, err)
	assert.Empty(t, result.Series)
}

// --- Measurement Trend ---

func TestMeasurementTrend_Success(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService: &mockBodyChartService{
			measurementTrendFn: func(_ context.Context, _ string, _ string, _, _ models.Date) ([]models.MeasurementTrendPoint, error) {
				return []models.MeasurementTrendPoint{
					{Date: models.MustDate("2026-03-01"), Value: 100},
				}, nil
			},
		},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{},
	}

	result, err := r.GetMeasurementTrend(userCtx("uid"), models.MeasurementTypeChest, nil, nil)
	require.NoError(t, err)
	assert.Len(t, result.DataPoints, 1)
}

func TestMeasurementTrend_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService:             &mockBodyChartService{},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{},
	}

	result, err := r.GetMeasurementTrend(context.Background(), models.MeasurementTypeChest, nil, nil)
	require.NoError(t, err)
	assert.NotNil(t, result.AuthErr)
}

func TestMeasurementTrend_InvalidDateRange(t *testing.T) {
	from := models.MustDate("2026-06-01")
	to := models.MustDate("2026-01-01")

	r := &resolver.Resolver{
		BodyChartService:             &mockBodyChartService{},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{},
	}

	result, err := r.GetMeasurementTrend(userCtx("uid"), models.MeasurementTypeChest, &from, &to)
	require.NoError(t, err)
	assert.NotNil(t, result.ValidationErr)
	assert.Equal(t, models.ChartErrorValidation, result.ValidationErr.Code)
}

// --- Measurement Overlay ---

func TestMeasurementOverlay_Success(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService: &mockBodyChartService{
			measurementOverlayFn: func(_ context.Context, _ string, _ []string, _, _ models.Date) ([]models.MeasurementOverlayGroup, error) {
				return []models.MeasurementOverlayGroup{
					{MeasurementType: models.MeasurementTypeChest, DataPoints: []models.MeasurementTrendPoint{}},
				}, nil
			},
		},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{},
	}

	result, err := r.GetMeasurementOverlay(userCtx("uid"), []models.MeasurementType{models.MeasurementTypeChest}, nil, nil)
	require.NoError(t, err)
	assert.Len(t, result.Groups, 1)
}

func TestMeasurementOverlay_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService:             &mockBodyChartService{},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{},
	}

	result, err := r.GetMeasurementOverlay(context.Background(), nil, nil, nil)
	require.NoError(t, err)
	assert.NotNil(t, result.AuthErr)
}

func TestMeasurementOverlay_EmptyTypes(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService: &mockBodyChartService{
			measurementOverlayFn: func(_ context.Context, _ string, _ []string, _, _ models.Date) ([]models.MeasurementOverlayGroup, error) {
				return []models.MeasurementOverlayGroup{}, nil
			},
		},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{},
	}

	result, err := r.GetMeasurementOverlay(userCtx("uid"), []models.MeasurementType{}, nil, nil)
	require.NoError(t, err)
	assert.Empty(t, result.Groups)
}

// --- Nutrition Weekly Averages ---

func TestNutritionWeeklyAverages_Success(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService: &mockBodyChartService{},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{
			weeklyAveragesFn: func(_ context.Context, _ string, _, _ models.Date) ([]models.NutritionWeeklyAverage, error) {
				return []models.NutritionWeeklyAverage{
					{WeekStartDate: models.MustDate("2026-01-05"), Calories: 2000, Protein: 150, Fat: 60, Carbs: 250},
				}, nil
			},
		},
	}

	result, err := r.GetNutritionWeeklyAverages(userCtx("uid"), nil, nil)
	require.NoError(t, err)
	assert.Len(t, result.Averages, 1)
	assert.Equal(t, 2000.0, result.Averages[0].Calories)
}

func TestNutritionWeeklyAverages_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService:             &mockBodyChartService{},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{},
	}

	result, err := r.GetNutritionWeeklyAverages(context.Background(), nil, nil)
	require.NoError(t, err)
	assert.NotNil(t, result.AuthErr)
}

func TestNutritionWeeklyAverages_EmptySeries(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService: &mockBodyChartService{},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{
			weeklyAveragesFn: func(_ context.Context, _ string, _, _ models.Date) ([]models.NutritionWeeklyAverage, error) {
				return []models.NutritionWeeklyAverage{}, nil
			},
		},
	}

	result, err := r.GetNutritionWeeklyAverages(userCtx("uid"), nil, nil)
	require.NoError(t, err)
	assert.Empty(t, result.Averages)
}

func TestNutritionWeeklyAverages_MaxRangeExceeded(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService: &mockBodyChartService{},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{
			weeklyAveragesFn: func(_ context.Context, _ string, _, _ models.Date) ([]models.NutritionWeeklyAverage, error) {
				return nil, atlasSvc.ErrMaxRangeExceeded
			},
		},
	}

	result, err := r.GetNutritionWeeklyAverages(userCtx("uid"), nil, nil)
	require.NoError(t, err)
	assert.NotNil(t, result.ValidationErr)
}

// --- Exercise Progress (stub) ---

func TestExerciseProgress_ReturnsEmptySeries(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService:             &mockBodyChartService{},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{},
	}

	result, err := r.GetExerciseProgress(userCtx("uid"), "ex-1", nil, nil)
	require.NoError(t, err)
	assert.Empty(t, result.DataPoints)
	assert.Nil(t, result.AuthErr)
}

func TestExerciseProgress_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		BodyChartService:             &mockBodyChartService{},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{},
	}

	result, err := r.GetExerciseProgress(context.Background(), "ex-1", nil, nil)
	require.NoError(t, err)
	assert.NotNil(t, result.AuthErr)
}

// --- Date validation ---

func TestBodyWeightTrend_InvalidDateRange(t *testing.T) {
	from := models.MustDate("2026-06-01")
	to := models.MustDate("2026-01-01")

	r := &resolver.Resolver{
		BodyChartService: &mockBodyChartService{},
		NutritionWeeklyAvgService: &mockNutritionWeeklyAvgService{},
	}

	result, err := r.GetBodyWeightTrend(userCtx("uid"), &from, &to)
	require.NoError(t, err)
	assert.NotNil(t, result.ValidationErr)
}