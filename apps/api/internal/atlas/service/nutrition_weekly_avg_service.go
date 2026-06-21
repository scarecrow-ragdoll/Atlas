package service

import (
	"context"
	"fmt"
	"time"

	"monorepo-template/apps/api/internal/atlas/models"
)

const (
	maxNutritionWeeks = 52
)

type NutritionWeeklyAvgService interface {
	WeeklyAverages(ctx context.Context, userID string, fromDate, toDate models.Date) ([]models.NutritionWeeklyAverage, error)
}

type nutritionWeeklyAvgService struct {
	macroService NutritionMacroService
}

func NewNutritionWeeklyAvgService(macroService NutritionMacroService) NutritionWeeklyAvgService {
	return &nutritionWeeklyAvgService{macroService: macroService}
}

func (s *nutritionWeeklyAvgService) WeeklyAverages(ctx context.Context, userID string, fromDate, toDate models.Date) ([]models.NutritionWeeklyAverage, error) {
	if fromDate.Time().After(toDate.Time()) {
		return nil, fmt.Errorf("from date must be on or before to date")
	}

	days := toDate.Time().Sub(fromDate.Time()).Hours() / 24
	if days > float64(maxNutritionWeeks*7) {
		return nil, ErrMaxRangeExceeded
	}

	var out []models.NutritionWeeklyAverage

	current := fromDate.Time()
	end := toDate.Time()

	for !current.After(end) {
		weekStart := weekStartDate(current)

		var weekEnd time.Time
		if weekStart.AddDate(0, 0, 6).After(end) {
			weekEnd = end
		} else {
			weekEnd = weekStart.AddDate(0, 0, 6)
		}

		var totalCal, totalProt, totalFat, totalCarbs float64
		daysInWeek := 0

		day := weekStart
		for !day.After(weekEnd) {
			dateStr := day.Format("2006-01-02")
			wsStr := weekStart.Format("2006-01-02")

			macros, err := s.macroService.Calculate(ctx, userID, wsStr, dateStr)
			if err != nil {
				return nil, fmt.Errorf("nutrition_weekly_avg_service.WeeklyAverages: %w", err)
			}
			if macros != nil {
				totalCal += macros.Calories
				totalProt += macros.Protein
				totalFat += macros.Fat
				totalCarbs += macros.Carbs
				daysInWeek++
			}
			day = day.AddDate(0, 0, 1)
		}

		if daysInWeek > 0 {
			out = append(out, models.NutritionWeeklyAverage{
				WeekStartDate: models.MustDate(weekStart.Format("2006-01-02")),
				Calories:      roundTo2(totalCal / float64(daysInWeek)),
				Protein:       roundTo2(totalProt / float64(daysInWeek)),
				Fat:           roundTo2(totalFat / float64(daysInWeek)),
				Carbs:         roundTo2(totalCarbs / float64(daysInWeek)),
			})
		}

		current = weekStart.AddDate(0, 0, 7)
	}

	if out == nil {
		out = []models.NutritionWeeklyAverage{}
	}

	return out, nil
}

func weekStartDate(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	return t.AddDate(0, 0, -int(weekday-time.Monday))
}

func roundTo2(v float64) float64 {
	return float64(int(v*100)) / 100
}
