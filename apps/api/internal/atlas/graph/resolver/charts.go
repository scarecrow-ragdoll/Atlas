package resolver

import (
	"context"
	"errors"
	"fmt"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
)

type defaultDateRangeFn func(atlasService.Clock) (models.Date, models.Date)

func (r *Resolver) GetBodyWeightTrend(ctx context.Context, from, to *models.Date) (*models.BodyWeightTrendResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyWeightTrendResult{
			AuthErr: &models.ChartAuthErr{Message: "unauthorized", Code: models.ChartErrorAuth},
		}, nil
	}

	f, t := defaultOrProvided(from, to, atlasService.DefaultDateRange)
	if f.Time().After(t.Time()) {
		return &models.BodyWeightTrendResult{
			ValidationErr: &models.ChartValidationErr{Message: "from date must be on or before to date", Code: models.ChartErrorValidation},
		}, nil
	}

	series, err := r.BodyChartService.BodyWeightTrend(ctx, userID, f, t)
	if err != nil {
		return nil, fmt.Errorf("body_weight_trend: %w", err)
	}
	if series == nil {
		series = []models.BodyWeightSeriesPoint{}
	}

	return &models.BodyWeightTrendResult{Series: series}, nil
}

func (r *Resolver) GetMeasurementTrend(ctx context.Context, measurementType models.MeasurementType, from, to *models.Date) (*models.MeasurementTrendResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.MeasurementTrendResult{
			AuthErr: &models.ChartAuthErr{Message: "unauthorized", Code: models.ChartErrorAuth},
		}, nil
	}

	f, t := defaultOrProvided(from, to, atlasService.DefaultDateRange)
	if f.Time().After(t.Time()) {
		return &models.MeasurementTrendResult{
			ValidationErr: &models.ChartValidationErr{Message: "from date must be on or before to date", Code: models.ChartErrorValidation},
		}, nil
	}

	points, err := r.BodyChartService.MeasurementTrend(ctx, userID, string(measurementType), f, t)
	if err != nil {
		return nil, fmt.Errorf("measurement_trend: %w", err)
	}
	if points == nil {
		points = []models.MeasurementTrendPoint{}
	}

	return &models.MeasurementTrendResult{DataPoints: points}, nil
}

func (r *Resolver) GetMeasurementOverlay(ctx context.Context, measurementTypes []models.MeasurementType, from, to *models.Date) (*models.MeasurementOverlayResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.MeasurementOverlayResult{
			AuthErr: &models.ChartAuthErr{Message: "unauthorized", Code: models.ChartErrorAuth},
		}, nil
	}

	f, t := defaultOrProvided(from, to, atlasService.DefaultDateRange)
	if f.Time().After(t.Time()) {
		return &models.MeasurementOverlayResult{
			ValidationErr: &models.ChartValidationErr{Message: "from date must be on or before to date", Code: models.ChartErrorValidation},
		}, nil
	}

	types := make([]string, len(measurementTypes))
	for i, mt := range measurementTypes {
		types[i] = string(mt)
	}

	groups, err := r.BodyChartService.MeasurementOverlay(ctx, userID, types, f, t)
	if err != nil {
		return nil, fmt.Errorf("measurement_overlay: %w", err)
	}
	if groups == nil {
		groups = []models.MeasurementOverlayGroup{}
	}

	return &models.MeasurementOverlayResult{Groups: groups}, nil
}

func (r *Resolver) GetNutritionWeeklyAverages(ctx context.Context, from, to *models.Date) (*models.NutritionWeeklyAveragesResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionWeeklyAveragesResult{
			AuthErr: &models.ChartAuthErr{Message: "unauthorized", Code: models.ChartErrorAuth},
		}, nil
	}

	f, t := defaultOrProvided(from, to, atlasService.DefaultDateRange)
	if f.Time().After(t.Time()) {
		return &models.NutritionWeeklyAveragesResult{
			ValidationErr: &models.ChartValidationErr{Message: "from date must be on or before to date", Code: models.ChartErrorValidation},
		}, nil
	}

	averages, err := r.NutritionWeeklyAvgService.WeeklyAverages(ctx, userID, f, t)
	if err != nil {
		if errors.Is(err, atlasService.ErrMaxRangeExceeded) {
			return &models.NutritionWeeklyAveragesResult{
				ValidationErr: &models.ChartValidationErr{Message: err.Error(), Code: models.ChartErrorValidation},
			}, nil
		}
		return nil, err
	}
	if averages == nil {
		averages = []models.NutritionWeeklyAverage{}
	}

	return &models.NutritionWeeklyAveragesResult{Averages: averages}, nil
}

func (r *Resolver) GetExerciseProgress(ctx context.Context, exerciseID string, from, to *models.Date) (*models.ExerciseProgressResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.ExerciseProgressResult{
			AuthErr: &models.ChartAuthErr{Message: "unauthorized", Code: models.ChartErrorAuth},
		}, nil
	}

	_ = exerciseID
	_ = from
	_ = to
	_ = userID

	return &models.ExerciseProgressResult{DataPoints: []models.ChartDataPoint{}}, nil
}

func defaultOrProvided(from, to *models.Date, fn defaultDateRangeFn) (models.Date, models.Date) {
	if from != nil && to != nil {
		return *from, *to
	}

	f, t := fn(atlasService.RealClock{})
	return f, t
}