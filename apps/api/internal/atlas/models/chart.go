package models

type ChartDataPoint struct {
	Date  Date    `json:"date"`
	Value float64 `json:"value"`
}

type BodyWeightSeriesPoint struct {
	Date   Date            `json:"date"`
	Weight float64         `json:"weight"`
	Source BodyWeightSource `json:"source"`
}

type MeasurementTrendPoint struct {
	Date  Date             `json:"date"`
	Value float64          `json:"value"`
	Side  *MeasurementSide `json:"side"`
}

type MeasurementOverlayGroup struct {
	MeasurementType MeasurementType       `json:"measurementType"`
	DataPoints      []MeasurementTrendPoint `json:"dataPoints"`
}

type NutritionWeeklyAverage struct {
	WeekStartDate Date    `json:"weekStartDate"`
	Calories      float64 `json:"calories"`
	Protein       float64 `json:"protein"`
	Fat           float64 `json:"fat"`
	Carbs         float64 `json:"carbs"`
}

type BodyWeightTrendResult struct {
	Series        []BodyWeightSeriesPoint `json:"series"`
	ValidationErr *ChartValidationErr     `json:"validationError"`
	AuthErr       *ChartAuthErr           `json:"authError"`
}

type MeasurementTrendResult struct {
	DataPoints    []MeasurementTrendPoint `json:"dataPoints"`
	ValidationErr *ChartValidationErr     `json:"validationError"`
	AuthErr       *ChartAuthErr           `json:"authError"`
}

type MeasurementOverlayResult struct {
	Groups        []MeasurementOverlayGroup `json:"groups"`
	ValidationErr *ChartValidationErr       `json:"validationError"`
	AuthErr       *ChartAuthErr             `json:"authError"`
}

type NutritionWeeklyAveragesResult struct {
	Averages      []NutritionWeeklyAverage `json:"averages"`
	ValidationErr *ChartValidationErr      `json:"validationError"`
	AuthErr       *ChartAuthErr            `json:"authError"`
}

type ExerciseProgressResult struct {
	DataPoints    []ChartDataPoint    `json:"dataPoints"`
	ValidationErr *ChartValidationErr `json:"validationError"`
	AuthErr       *ChartAuthErr       `json:"authError"`
}

type ChartValidationErr struct {
	Message string         `json:"message"`
	Code    ChartErrorCode `json:"code"`
}

func (e *ChartValidationErr) Error() string { return e.Message }

type ChartAuthErr struct {
	Message string         `json:"message"`
	Code    ChartErrorCode `json:"code"`
}

func (e *ChartAuthErr) Error() string { return e.Message }

type ChartNotFoundErr struct {
	Message string         `json:"message"`
	Code    ChartErrorCode `json:"code"`
}

func (e *ChartNotFoundErr) Error() string { return e.Message }

type ChartErrorCode string

const (
	ChartErrorValidation ChartErrorCode = "VALIDATION_ERROR"
	ChartErrorNotFound   ChartErrorCode = "NOT_FOUND"
	ChartErrorAuth       ChartErrorCode = "AUTH_ERROR"
	ChartErrorInternal   ChartErrorCode = "INTERNAL_ERROR"
)

type BodyMeasurementTrendRecord struct {
	ID              string
	CheckInID       string
	MeasurementType string
	Side            *string
	Value           float64
	Date            Date
	CreatedAt       string
	UpdatedAt       string
}

func CalculateE1RM(weight, reps float64) float64 {
	if weight <= 0 || reps <= 0 {
		return 0
	}
	return weight * (1 + reps/30)
}
