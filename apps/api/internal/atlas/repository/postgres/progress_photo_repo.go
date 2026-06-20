// FILE: apps/api/internal/atlas/repository/postgres/progress_photo_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement ProgressPhotoRepository for WAVE-04 progress photo tracking using sqlc-generated queries.
//   SCOPE: CRUD operations for progress photos, list by check-in ID, count by check-in, all user-scoped via check-in join.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ProgressPhotoRepository - Interface for progress photo data access.
//   NewProgressPhotoRepository - Creates a new ProgressPhotoRepository.
//   Create - Creates a progress photo record.
//   GetByID - Gets a progress photo by ID (user-scoped).
//   ListByCheckIn - Lists progress photos by check-in ID.
//   Delete - Deletes a progress photo.
//   CountByCheckIn - Counts progress photos by check-in ID.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added progress photo repository for WAVE-04.
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

type ProgressPhotoRepository interface {
	Create(ctx context.Context, checkInID string, filePath string, originalFileName string, mimeType string, sizeBytes int64, angle *string, label *string, notes *string) (*models.ProgressPhotoRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.ProgressPhotoRecord, error)
	ListByCheckIn(ctx context.Context, userID string, checkInID string) ([]models.ProgressPhotoRecord, error)
	Delete(ctx context.Context, userID string, id string) (*models.ProgressPhotoRecord, error)
	CountByCheckIn(ctx context.Context, checkInID string) (int64, error)
}

type progressPhotoRepository struct {
	q *generated.Queries
}

func NewProgressPhotoRepository(pool *pgxpool.Pool) ProgressPhotoRepository {
	return &progressPhotoRepository{q: generated.New(pool)}
}

func (r *progressPhotoRepository) Create(ctx context.Context, checkInID string, filePath string, originalFileName string, mimeType string, sizeBytes int64, angle *string, label *string, notes *string) (*models.ProgressPhotoRecord, error) {
	cid, err := uuidFromString(checkInID)
	if err != nil {
		return nil, fmt.Errorf("progress_photo_repo.Create: %w", err)
	}

	row, err := r.q.CreateProgressPhoto(ctx, generated.CreateProgressPhotoParams{
		CheckInID:        cid,
		FilePath:         filePath,
		OriginalFileName: originalFileName,
		MimeType:         mimeType,
		SizeBytes:        sizeBytes,
		Angle:            nullableText(angle),
		Label:            nullableText(label),
		Notes:            nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("progress_photo_repo.Create: %w", err)
	}

	return progressPhotoRecordFromRow(row), nil
}

func (r *progressPhotoRepository) GetByID(ctx context.Context, userID string, id string) (*models.ProgressPhotoRecord, error) {
	uid, pid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("progress_photo_repo.GetByID: %w", err)
	}

	row, err := r.q.GetProgressPhotoByID(ctx, generated.GetProgressPhotoByIDParams{ID: pid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("progress_photo_repo.GetByID: %w", err)
	}

	return progressPhotoRecordFromRow(row), nil
}

func (r *progressPhotoRepository) ListByCheckIn(ctx context.Context, userID string, checkInID string) ([]models.ProgressPhotoRecord, error) {
	uid, cid, err := parseTwoUUIDs(userID, checkInID)
	if err != nil {
		return nil, fmt.Errorf("progress_photo_repo.ListByCheckIn: %w", err)
	}

	rows, err := r.q.ListProgressPhotosByCheckIn(ctx, generated.ListProgressPhotosByCheckInParams{
		CheckInID: cid,
		UserID:    uid,
	})
	if err != nil {
		return nil, fmt.Errorf("progress_photo_repo.ListByCheckIn: %w", err)
	}

	out := make([]models.ProgressPhotoRecord, len(rows))
	for i, row := range rows {
		out[i] = *progressPhotoRecordFromRow(row)
	}
	return out, nil
}

func (r *progressPhotoRepository) Delete(ctx context.Context, userID string, id string) (*models.ProgressPhotoRecord, error) {
	uid, pid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("progress_photo_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteProgressPhoto(ctx, generated.DeleteProgressPhotoParams{ID: pid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("progress_photo_repo.Delete: %w", err)
	}

	return progressPhotoRecordFromRow(row), nil
}

func (r *progressPhotoRepository) CountByCheckIn(ctx context.Context, checkInID string) (int64, error) {
	cid, err := uuidFromString(checkInID)
	if err != nil {
		return 0, fmt.Errorf("progress_photo_repo.CountByCheckIn: %w", err)
	}

	count, err := r.q.CountProgressPhotosByCheckIn(ctx, cid)
	if err != nil {
		return 0, fmt.Errorf("progress_photo_repo.CountByCheckIn: %w", err)
	}

	return count, nil
}

func progressPhotoRecordFromRow(row generated.ProgressPhoto) *models.ProgressPhotoRecord {
	return &models.ProgressPhotoRecord{
		ID:               row.ID.String(),
		CheckInID:        row.CheckInID.String(),
		FilePath:         row.FilePath,
		OriginalFileName: row.OriginalFileName,
		MimeType:         row.MimeType,
		SizeBytes:        row.SizeBytes,
		Angle:            textPtr(row.Angle),
		Label:            textPtr(row.Label),
		Notes:            textPtr(row.Notes),
		CreatedAt:        formatTimestamp(row.CreatedAt),
		UpdatedAt:        formatTimestamp(row.UpdatedAt),
	}
}