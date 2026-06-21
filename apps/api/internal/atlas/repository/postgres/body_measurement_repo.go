// FILE: apps/api/internal/atlas/repository/postgres/body_measurement_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement BodyMeasurementRepository for WAVE-04 body measurement tracking using sqlc-generated queries.
//   SCOPE: CRUD operations for body measurements, list by check-in ID, all user-scoped via check-in join.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   BodyMeasurementRepository - Interface for body measurement data access.
//   NewBodyMeasurementRepository - Creates a new BodyMeasurementRepository.
//   Create - Creates a body measurement.
//   GetByID - Gets a body measurement by ID (user-scoped).
//   ListByCheckIn - Lists body measurements by check-in ID.
//   Update - Updates a body measurement.
//   Delete - Deletes a body measurement.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added body measurement repository for WAVE-04.
// END_CHANGE_SUMMARY

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/repository/postgres/generated"
)

type BodyMeasurementRepository interface {
	Create(ctx context.Context, checkInID string, measurementType string, side *string, value float64) (*models.BodyMeasurementRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.BodyMeasurementRecord, error)
	ListByCheckIn(ctx context.Context, userID string, checkInID string) ([]models.BodyMeasurementRecord, error)
	ListByUserTypeRange(ctx context.Context, userID string, measurementType string, fromDate models.Date, toDate models.Date) ([]models.BodyMeasurementTrendRecord, error)
	Update(ctx context.Context, userID string, id string, measurementType string, side *string, value float64) (*models.BodyMeasurementRecord, error)
	Delete(ctx context.Context, userID string, id string) (*models.BodyMeasurementRecord, error)
}

type bodyMeasurementRepository struct {
	q *generated.Queries
}

func NewBodyMeasurementRepository(pool *pgxpool.Pool) BodyMeasurementRepository {
	return &bodyMeasurementRepository{q: generated.New(pool)}
}

func (r *bodyMeasurementRepository) Create(ctx context.Context, checkInID string, measurementType string, side *string, value float64) (*models.BodyMeasurementRecord, error) {
	cid, err := uuidFromString(checkInID)
	if err != nil {
		return nil, fmt.Errorf("body_measurement_repo.Create: %w", err)
	}

	row, err := r.q.CreateBodyMeasurement(ctx, generated.CreateBodyMeasurementParams{
		CheckInID:       cid,
		MeasurementType: measurementType,
		Side:            nullableText(side),
		Value:           float32(value),
	})
	if err != nil {
		return nil, fmt.Errorf("body_measurement_repo.Create: %w", err)
	}

	return bodyMeasurementRecordFromRow(row), nil
}

func (r *bodyMeasurementRepository) GetByID(ctx context.Context, userID string, id string) (*models.BodyMeasurementRecord, error) {
	uid, mid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_measurement_repo.GetByID: %w", err)
	}

	row, err := r.q.GetBodyMeasurementByID(ctx, generated.GetBodyMeasurementByIDParams{ID: mid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("body_measurement_repo.GetByID: %w", err)
	}

	return bodyMeasurementRecordFromRow(row), nil
}

func (r *bodyMeasurementRepository) ListByCheckIn(ctx context.Context, userID string, checkInID string) ([]models.BodyMeasurementRecord, error) {
	uid, cid, err := parseTwoUUIDs(userID, checkInID)
	if err != nil {
		return nil, fmt.Errorf("body_measurement_repo.ListByCheckIn: %w", err)
	}

	rows, err := r.q.ListBodyMeasurementsByCheckIn(ctx, generated.ListBodyMeasurementsByCheckInParams{
		CheckInID: cid,
		UserID:    uid,
	})
	if err != nil {
		return nil, fmt.Errorf("body_measurement_repo.ListByCheckIn: %w", err)
	}

	out := make([]models.BodyMeasurementRecord, len(rows))
	for i, row := range rows {
		out[i] = *bodyMeasurementRecordFromRow(row)
	}
	return out, nil
}

func (r *bodyMeasurementRepository) ListByUserTypeRange(ctx context.Context, userID string, measurementType string, fromDate models.Date, toDate models.Date) ([]models.BodyMeasurementTrendRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("body_measurement_repo.ListByUserTypeRange: %w", err)
	}

	rows, err := r.q.ListBodyMeasurementsByUserTypeRange(ctx, generated.ListBodyMeasurementsByUserTypeRangeParams{
		UserID:          uid,
		MeasurementType: measurementType,
		Column3:         modelsToPGDate(fromDate),
		Column4:         modelsToPGDate(toDate),
	})
	if err != nil {
		return nil, fmt.Errorf("body_measurement_repo.ListByUserTypeRange: %w", err)
	}

	out := make([]models.BodyMeasurementTrendRecord, len(rows))
	for i, row := range rows {
		out[i] = models.BodyMeasurementTrendRecord{
			ID:              row.ID.String(),
			CheckInID:       row.CheckInID.String(),
			MeasurementType: row.MeasurementType,
			Side:            textPtr(row.Side),
			Value:           float64(row.Value),
			Date:            dateFromPGDate(row.CheckInDate),
			CreatedAt:       formatTimestamp(row.CreatedAt),
			UpdatedAt:       formatTimestamp(row.UpdatedAt),
		}
	}
	return out, nil
}

func (r *bodyMeasurementRepository) Update(ctx context.Context, userID string, id string, measurementType string, side *string, value float64) (*models.BodyMeasurementRecord, error) {
	uid, mid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_measurement_repo.Update: %w", err)
	}

	row, err := r.q.UpdateBodyMeasurement(ctx, generated.UpdateBodyMeasurementParams{
		ID:              mid,
		UserID:          uid,
		MeasurementType: measurementType,
		Side:            nullableText(side),
		Value:           float32(value),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("body_measurement_repo.Update: %w", err)
	}

	return bodyMeasurementRecordFromRow(row), nil
}

func (r *bodyMeasurementRepository) Delete(ctx context.Context, userID string, id string) (*models.BodyMeasurementRecord, error) {
	uid, mid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_measurement_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteBodyMeasurement(ctx, generated.DeleteBodyMeasurementParams{ID: mid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("body_measurement_repo.Delete: %w", err)
	}

	return bodyMeasurementRecordFromRow(row), nil
}

func bodyMeasurementRecordFromRow(row generated.BodyMeasurement) *models.BodyMeasurementRecord {
	return &models.BodyMeasurementRecord{
		ID:              row.ID.String(),
		CheckInID:       row.CheckInID.String(),
		MeasurementType: row.MeasurementType,
		Side:            textPtr(row.Side),
		Value:           float64(row.Value),
		CreatedAt:       formatTimestamp(row.CreatedAt),
		UpdatedAt:       formatTimestamp(row.UpdatedAt),
	}
}