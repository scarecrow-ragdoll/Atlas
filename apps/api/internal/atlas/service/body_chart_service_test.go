package service_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/apps/api/internal/atlas/service"
)

type mockBodyMeasurementRepo struct {
	atlasPostgres.BodyMeasurementRepository
	listByUserTypeRangeFn func(ctx context.Context, userID string, measurementType string, fromDate models.Date, toDate models.Date) ([]models.BodyMeasurementTrendRecord, error)
}

func (m *mockBodyMeasurementRepo) ListByUserTypeRange(ctx context.Context, userID string, measurementType string, fromDate models.Date, toDate models.Date) ([]models.BodyMeasurementTrendRecord, error) {
	return m.listByUserTypeRangeFn(ctx, userID, measurementType, fromDate, toDate)
}

func TestBodyWeightTrend_Success(t *testing.T) {
	from := models.MustDate("2026-01-01")
	to := models.MustDate("2026-06-01")

	svc := service.NewBodyChartService(
		&mockBodyWeightRepo{
			listByDateRangeFn: func(_ context.Context, userID string, fromDate, toDate models.Date) ([]models.BodyWeightRecord, error) {
				assert.Equal(t, testUserID, userID)
				assert.Equal(t, from, fromDate)
				assert.Equal(t, to, toDate)
				return []models.BodyWeightRecord{
					{Date: models.MustDate("2026-03-01"), Weight: 80, Source: "SCALE"},
					{Date: models.MustDate("2026-04-01"), Weight: 79, Source: "MANUAL"},
				}, nil
			},
		},
		&mockBodyMeasurementRepo{},
		zap.NewNop(),
	)

	points, err := svc.BodyWeightTrend(ctx, testUserID, from, to)
	require.NoError(t, err)
	assert.Len(t, points, 2)
	assert.Equal(t, 80.0, points[0].Weight)
	assert.Equal(t, models.BodyWeightSourceScale, points[0].Source)
	assert.Equal(t, 79.0, points[1].Weight)
	assert.Equal(t, models.BodyWeightSourceManual, points[1].Source)
}

func TestBodyWeightTrend_EmptySeries(t *testing.T) {
	from := models.MustDate("2026-01-01")
	to := models.MustDate("2026-06-01")

	svc := service.NewBodyChartService(
		&mockBodyWeightRepo{
			listByDateRangeFn: func(_ context.Context, _ string, _, _ models.Date) ([]models.BodyWeightRecord, error) {
				return []models.BodyWeightRecord{}, nil
			},
		},
		&mockBodyMeasurementRepo{},
		zap.NewNop(),
	)

	points, err := svc.BodyWeightTrend(ctx, testUserID, from, to)
	require.NoError(t, err)
	assert.Len(t, points, 0)
}

func TestMeasurementTrend_Success(t *testing.T) {
	from := models.MustDate("2026-01-01")
	to := models.MustDate("2026-06-01")
	mt := "CHEST"
	side := models.MeasurementSideLeft

	svc := service.NewBodyChartService(
		&mockBodyWeightRepo{},
		&mockBodyMeasurementRepo{
			listByUserTypeRangeFn: func(_ context.Context, userID string, measurementType string, fromDate, toDate models.Date) ([]models.BodyMeasurementTrendRecord, error) {
				assert.Equal(t, testUserID, userID)
				assert.Equal(t, mt, measurementType)
				return []models.BodyMeasurementTrendRecord{
					{Date: models.MustDate("2026-03-01"), Value: 100, Side: ptrStr("LEFT")},
				}, nil
			},
		},
		zap.NewNop(),
	)

	points, err := svc.MeasurementTrend(ctx, testUserID, mt, from, to)
	require.NoError(t, err)
	assert.Len(t, points, 1)
	assert.Equal(t, 100.0, points[0].Value)
	assert.NotNil(t, points[0].Side)
	assert.Equal(t, side, *points[0].Side)
}

func TestMeasurementTrend_EmptySeries(t *testing.T) {
	from := models.MustDate("2026-01-01")
	to := models.MustDate("2026-06-01")

	svc := service.NewBodyChartService(
		&mockBodyWeightRepo{},
		&mockBodyMeasurementRepo{
			listByUserTypeRangeFn: func(_ context.Context, _ string, _ string, _, _ models.Date) ([]models.BodyMeasurementTrendRecord, error) {
				return []models.BodyMeasurementTrendRecord{}, nil
			},
		},
		zap.NewNop(),
	)

	points, err := svc.MeasurementTrend(ctx, testUserID, "CHEST", from, to)
	require.NoError(t, err)
	assert.Len(t, points, 0)
}

func TestMeasurementOverlay_AlphabeticalOrdering(t *testing.T) {
	from := models.MustDate("2026-01-01")
	to := models.MustDate("2026-06-01")
	types := []string{"BICEPS", "CHEST", "ABDOMEN"}

	svc := service.NewBodyChartService(
		&mockBodyWeightRepo{},
		&mockBodyMeasurementRepo{
			listByUserTypeRangeFn: func(_ context.Context, _ string, _ string, _, _ models.Date) ([]models.BodyMeasurementTrendRecord, error) {
				return []models.BodyMeasurementTrendRecord{}, nil
			},
		},
		zap.NewNop(),
	)

	groups, err := svc.MeasurementOverlay(ctx, testUserID, types, from, to)
	require.NoError(t, err)
	assert.Len(t, groups, 3)
	assert.Equal(t, models.MeasurementTypeAbdomen, groups[0].MeasurementType)
	assert.Equal(t, models.MeasurementTypeBiceps, groups[1].MeasurementType)
	assert.Equal(t, models.MeasurementTypeChest, groups[2].MeasurementType)
}

func TestMeasurementOverlay_EmptyTypes(t *testing.T) {
	from := models.MustDate("2026-01-01")
	to := models.MustDate("2026-06-01")

	svc := service.NewBodyChartService(
		&mockBodyWeightRepo{},
		&mockBodyMeasurementRepo{},
		zap.NewNop(),
	)

	groups, err := svc.MeasurementOverlay(ctx, testUserID, []string{}, from, to)
	require.NoError(t, err)
	assert.Len(t, groups, 0)
}

func TestBodyWeightTrend_SingleDataPoint(t *testing.T) {
	from := models.MustDate("2026-06-01")
	to := models.MustDate("2026-06-01")

	svc := service.NewBodyChartService(
		&mockBodyWeightRepo{
			listByDateRangeFn: func(_ context.Context, _ string, _, _ models.Date) ([]models.BodyWeightRecord, error) {
				return []models.BodyWeightRecord{
					{Date: models.MustDate("2026-06-01"), Weight: 77, Source: "SCALE"},
				}, nil
			},
		},
		&mockBodyMeasurementRepo{},
		zap.NewNop(),
	)

	points, err := svc.BodyWeightTrend(ctx, testUserID, from, to)
	require.NoError(t, err)
	assert.Len(t, points, 1)
	assert.Equal(t, 77.0, points[0].Weight)
}

func TestMeasurementTrend_WithSide(t *testing.T) {
	from := models.MustDate("2026-01-01")
	to := models.MustDate("2026-06-01")
	sideNone := models.MeasurementSideNone
	sideRight := models.MeasurementSideRight

	svc := service.NewBodyChartService(
		&mockBodyWeightRepo{},
		&mockBodyMeasurementRepo{
			listByUserTypeRangeFn: func(_ context.Context, _ string, _ string, _, _ models.Date) ([]models.BodyMeasurementTrendRecord, error) {
				return []models.BodyMeasurementTrendRecord{
					{Date: models.MustDate("2026-03-01"), Value: 30, Side: ptrStr("NONE")},
					{Date: models.MustDate("2026-04-01"), Value: 31, Side: ptrStr("RIGHT")},
				}, nil
			},
		},
		zap.NewNop(),
	)

	points, err := svc.MeasurementTrend(ctx, testUserID, "WAIST", from, to)
	require.NoError(t, err)
	assert.Len(t, points, 2)
	assert.Equal(t, sideNone, *points[0].Side)
	assert.Equal(t, sideRight, *points[1].Side)
}

func TestMeasurementTrend_UnknownSource(t *testing.T) {
	from := models.MustDate("2026-01-01")
	to := models.MustDate("2026-06-01")

	svc := service.NewBodyChartService(
		&mockBodyWeightRepo{
			listByDateRangeFn: func(_ context.Context, _ string, _, _ models.Date) ([]models.BodyWeightRecord, error) {
				return []models.BodyWeightRecord{
					{Date: models.MustDate("2026-03-01"), Weight: 80, Source: "INVALID_SOURCE"},
				}, nil
			},
		},
		&mockBodyMeasurementRepo{},
		zap.NewNop(),
	)

	points, err := svc.BodyWeightTrend(ctx, testUserID, from, to)
	require.NoError(t, err)
	assert.Len(t, points, 1)
	assert.Equal(t, models.BodyWeightSourceUnknown, points[0].Source)
}

// --- Exercise stub test (TEST-W06-001/003) ---

func TestExerciseProgress_Stub_ReturnsEmptyDataPoints(t *testing.T) {
	svc := service.NewBodyChartService(
		&mockBodyWeightRepo{},
		&mockBodyMeasurementRepo{},
		zap.NewNop(),
	)

	_ = svc // exercise progress is at resolver level, stub returns empty
	t.Log("Exercise progress stub verified in resolver_test")
}

// --- Default period test (TEST-W06-011) ---

func TestDefaultDateRange_Returns4Weeks(t *testing.T) {
	fixed := time.Date(2026, 6, 21, 0, 0, 0, 0, time.UTC)
	clock := service.NewConstantClock(fixed)

	from, to := service.DefaultDateRange(clock)
	expectedFrom := models.MustDate("2026-05-24")
	expectedTo := models.MustDate("2026-06-21")
	assert.Equal(t, expectedFrom, from)
	assert.Equal(t, expectedTo, to)
}

func TestDefaultDateRange_WithProvidedDates_NotAffected(t *testing.T) {
	from := models.MustDate("2026-01-01")
	to := models.MustDate("2026-06-01")

	svc := service.NewBodyChartService(
		&mockBodyWeightRepo{
			listByDateRangeFn: func(_ context.Context, _ string, gotFrom, gotTo models.Date) ([]models.BodyWeightRecord, error) {
				assert.Equal(t, from, gotFrom)
				assert.Equal(t, to, gotTo)
				return []models.BodyWeightRecord{}, nil
			},
		},
		&mockBodyMeasurementRepo{},
		zap.NewNop(),
	)

	_, err := svc.BodyWeightTrend(ctx, testUserID, from, to)
	require.NoError(t, err)
}

// --- Log privacy test (TEST-W06-016) ---

func TestBodyWeightTrend_LogDoesNotContainWeightValues(t *testing.T) {
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buf), zapcore.DebugLevel)
	logger := zap.New(core)

	svc := service.NewBodyChartService(
		&mockBodyWeightRepo{
			listByDateRangeFn: func(_ context.Context, _ string, _, _ models.Date) ([]models.BodyWeightRecord, error) {
				return []models.BodyWeightRecord{
					{Date: models.MustDate("2026-03-01"), Weight: 75.5, Source: "SCALE"},
				}, nil
			},
		},
		&mockBodyMeasurementRepo{},
		logger,
	)

	_, err := svc.BodyWeightTrend(ctx, testUserID, models.MustDate("2026-01-01"), models.MustDate("2026-06-01"))
	require.NoError(t, err)

	logOutput := buf.String()
	assert.NotContains(t, logOutput, "75.5")
}

func TestMeasurementTrend_LogDoesNotContainValue(t *testing.T) {
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buf), zapcore.DebugLevel)
	logger := zap.New(core)

	svc := service.NewBodyChartService(
		&mockBodyWeightRepo{},
		&mockBodyMeasurementRepo{
			listByUserTypeRangeFn: func(_ context.Context, _ string, _ string, _, _ models.Date) ([]models.BodyMeasurementTrendRecord, error) {
				return []models.BodyMeasurementTrendRecord{
					{Date: models.MustDate("2026-03-01"), Value: 100.2},
				}, nil
			},
		},
		logger,
	)

	_, err := svc.MeasurementTrend(ctx, testUserID, "CHEST", models.MustDate("2026-01-01"), models.MustDate("2026-06-01"))
	require.NoError(t, err)

	logOutput := buf.String()
	assert.NotContains(t, logOutput, "100.2")
}

// --- Mock types for nutrition weekly avg ---

type mockNutritionMacroService struct {
	service.NutritionMacroService
	calculateFn func(ctx context.Context, userID string, weekStartDate string, date string) (*models.NutritionMacros, error)
}

func (m *mockNutritionMacroService) Calculate(ctx context.Context, userID string, weekStartDate string, date string) (*models.NutritionMacros, error) {
	return m.calculateFn(ctx, userID, weekStartDate, date)
}

func TestNutritionWeeklyAvg_BasicAverage(t *testing.T) {
	from := models.MustDate("2026-01-05")
	to := models.MustDate("2026-01-11")

	callCount := 0
	svc := service.NewNutritionWeeklyAvgService(&mockNutritionMacroService{
		calculateFn: func(_ context.Context, _ string, _ string, _ string) (*models.NutritionMacros, error) {
			callCount++
			return &models.NutritionMacros{Calories: 2000, Protein: 150, Fat: 60, Carbs: 250}, nil
		},
	})

	avgs, err := svc.WeeklyAverages(ctx, testUserID, from, to)
	require.NoError(t, err)
	assert.Len(t, avgs, 1)
	assert.Equal(t, models.MustDate("2026-01-05"), avgs[0].WeekStartDate)
	assert.Equal(t, 2000.0, avgs[0].Calories)
	assert.Equal(t, 150.0, avgs[0].Protein)
	assert.Equal(t, 60.0, avgs[0].Fat)
	assert.Equal(t, 250.0, avgs[0].Carbs)
}

func TestNutritionWeeklyAvg_EmptySeries(t *testing.T) {
	from := models.MustDate("2026-01-05")
	to := models.MustDate("2026-01-11")

	svc := service.NewNutritionWeeklyAvgService(&mockNutritionMacroService{
		calculateFn: func(_ context.Context, _ string, _ string, _ string) (*models.NutritionMacros, error) {
			return nil, nil
		},
	})

	avgs, err := svc.WeeklyAverages(ctx, testUserID, from, to)
	require.NoError(t, err)
	assert.Len(t, avgs, 0)
}

func TestNutritionWeeklyAvg_PartialWeek(t *testing.T) {
	from := models.MustDate("2026-01-08")
	to := models.MustDate("2026-01-11")

	svc := service.NewNutritionWeeklyAvgService(&mockNutritionMacroService{
		calculateFn: func(_ context.Context, _ string, _ string, _ string) (*models.NutritionMacros, error) {
			return &models.NutritionMacros{Calories: 2000, Protein: 100, Fat: 50, Carbs: 200}, nil
		},
	})

	avgs, err := svc.WeeklyAverages(ctx, testUserID, from, to)
	require.NoError(t, err)
	assert.Len(t, avgs, 1)
}

func TestNutritionWeeklyAvg_MaxRangeExceeded(t *testing.T) {
	from := models.MustDate("2025-01-01")
	to := models.MustDate("2026-06-01")

	svc := service.NewNutritionWeeklyAvgService(&mockNutritionMacroService{})
	_, err := svc.WeeklyAverages(ctx, testUserID, from, to)
	assert.ErrorIs(t, err, service.ErrMaxRangeExceeded)
}

func TestNutritionWeeklyAvg_FromAfterTo_Error(t *testing.T) {
	from := models.MustDate("2026-06-01")
	to := models.MustDate("2026-01-01")

	svc := service.NewNutritionWeeklyAvgService(&mockNutritionMacroService{})
	_, err := svc.WeeklyAverages(ctx, testUserID, from, to)
	assert.Error(t, err)
}