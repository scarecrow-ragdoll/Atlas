package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasRepo "monorepo-template/apps/api/internal/atlas/repository/postgres"
)

const (
	defaultChartPeriodDays = 28
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

type constantClock struct {
	t time.Time
}

func (c constantClock) Now() time.Time { return c.t }

func NewConstantClock(t time.Time) Clock {
	return constantClock{t: t}
}

var ErrMaxRangeExceeded = fmt.Errorf("date range exceeds maximum of 52 weeks")

type BodyChartService interface {
	BodyWeightTrend(ctx context.Context, userID string, fromDate, toDate models.Date) ([]models.BodyWeightSeriesPoint, error)
	MeasurementTrend(ctx context.Context, userID string, measurementType string, fromDate, toDate models.Date) ([]models.MeasurementTrendPoint, error)
	MeasurementOverlay(ctx context.Context, userID string, measurementTypes []string, fromDate, toDate models.Date) ([]models.MeasurementOverlayGroup, error)
}

type bodyChartService struct {
	weightRepo      atlasRepo.BodyWeightEntryRepository
	measurementRepo atlasRepo.BodyMeasurementRepository
	logger          *zap.Logger
}

func NewBodyChartService(
	weightRepo atlasRepo.BodyWeightEntryRepository,
	measurementRepo atlasRepo.BodyMeasurementRepository,
	logger *zap.Logger,
) BodyChartService {
	return &bodyChartService{
		weightRepo:      weightRepo,
		measurementRepo: measurementRepo,
		logger:          logger,
	}
}

func (s *bodyChartService) BodyWeightTrend(ctx context.Context, userID string, fromDate, toDate models.Date) ([]models.BodyWeightSeriesPoint, error) {
	s.logger.Info("[BodyChart][bodyWeightTrend]",
		zap.String("user_id", userID),
		zap.String("from", fromDate.String()),
		zap.String("to", toDate.String()),
	)

	records, err := s.weightRepo.ListByDateRange(ctx, userID, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("body_chart_service.BodyWeightTrend: %w", err)
	}

	out := make([]models.BodyWeightSeriesPoint, 0, len(records))
	for _, r := range records {
		source := models.BodyWeightSource(r.Source)
		if !models.IsValidBodyWeightSource(string(source)) {
			source = models.BodyWeightSourceUnknown
		}
		out = append(out, models.BodyWeightSeriesPoint{
			Date:   r.Date,
			Weight: r.Weight,
			Source: source,
		})
	}

	s.logger.Info("[BodyChart][bodyWeightTrend][result]",
		zap.Int("count", len(out)),
	)
	return out, nil
}

func (s *bodyChartService) MeasurementTrend(ctx context.Context, userID string, measurementType string, fromDate, toDate models.Date) ([]models.MeasurementTrendPoint, error) {
	s.logger.Info("[BodyChart][measurementTrend]",
		zap.String("user_id", userID),
		zap.String("type", measurementType),
		zap.String("from", fromDate.String()),
		zap.String("to", toDate.String()),
	)

	records, err := s.measurementRepo.ListByUserTypeRange(ctx, userID, measurementType, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("body_chart_service.MeasurementTrend: %w", err)
	}

	out := make([]models.MeasurementTrendPoint, 0, len(records))
	for _, r := range records {
		var side *models.MeasurementSide
		if r.Side != nil {
			s := models.MeasurementSide(*r.Side)
			side = &s
		}
		out = append(out, models.MeasurementTrendPoint{
			Date:  r.Date,
			Value: r.Value,
			Side:  side,
		})
	}

	s.logger.Info("[BodyChart][measurementTrend][result]",
		zap.Int("count", len(out)),
	)
	return out, nil
}

func (s *bodyChartService) MeasurementOverlay(ctx context.Context, userID string, measurementTypes []string, fromDate, toDate models.Date) ([]models.MeasurementOverlayGroup, error) {
	s.logger.Info("[BodyChart][measurementOverlay]",
		zap.String("user_id", userID),
		zap.Strings("types", measurementTypes),
		zap.String("from", fromDate.String()),
		zap.String("to", toDate.String()),
	)

	var out []models.MeasurementOverlayGroup

	for _, mt := range measurementTypes {
		points, err := s.MeasurementTrend(ctx, userID, mt, fromDate, toDate)
		if err != nil {
			return nil, fmt.Errorf("body_chart_service.MeasurementOverlay: %w", err)
		}
		out = append(out, models.MeasurementOverlayGroup{
			MeasurementType: models.MeasurementType(mt),
			DataPoints:      points,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return string(out[i].MeasurementType) < string(out[j].MeasurementType)
	})

	s.logger.Info("[BodyChart][measurementOverlay][result]",
		zap.Int("count", len(out)),
	)
	return out, nil
}

func DefaultDateRange(clock Clock) (models.Date, models.Date) {
	now := clock.Now()
	to := models.MustDate(now.Format("2006-01-02"))
	from := models.MustDate(now.AddDate(0, 0, -defaultChartPeriodDays).Format("2006-01-02"))
	return from, to
}