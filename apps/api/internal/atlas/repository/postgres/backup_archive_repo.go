// FILE: apps/api/internal/atlas/repository/postgres/backup_archive_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement BackupArchiveRepository for WAVE-09 backup archive management using sqlc-generated queries.
//   SCOPE: Create, GetByID, UpdateFilePath for backup archive records. All user-scoped.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-09.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   BackupArchiveRepository - Interface for backup archive data access.
//   NewBackupArchiveRepository - Creates a new BackupArchiveRepository.
//   Create - Creates a backup archive record.
//   GetByID - Gets a backup archive by ID (user-scoped).
//   UpdateFilePath - Updates the archive file path.
// END_MODULE_MAP

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

type BackupArchiveRepository interface {
	Create(ctx context.Context, userID string, includeMedia bool, sizeBytes int64, entityCounts []byte) (*models.BackupArchiveRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.BackupArchiveRecord, error)
	UpdateFilePath(ctx context.Context, id string, filePath *string) (*models.BackupArchiveRecord, error)
}

type backupArchiveRepository struct {
	q *generated.Queries
}

func NewBackupArchiveRepository(pool *pgxpool.Pool) BackupArchiveRepository {
	return &backupArchiveRepository{q: generated.New(pool)}
}

func (r *backupArchiveRepository) Create(ctx context.Context, userID string, includeMedia bool, sizeBytes int64, entityCounts []byte) (*models.BackupArchiveRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("backup_archive_repo.Create: %w", err)
	}

	row, err := r.q.CreateBackupArchive(ctx, generated.CreateBackupArchiveParams{
		UserID:       uid,
		IncludeMedia: includeMedia,
		SizeBytes:    sizeBytes,
		EntityCounts: entityCounts,
	})
	if err != nil {
		return nil, fmt.Errorf("backup_archive_repo.Create: %w", err)
	}

	return backupArchiveRecordFromRow(row), nil
}

func (r *backupArchiveRepository) GetByID(ctx context.Context, userID string, id string) (*models.BackupArchiveRecord, error) {
	uid, bid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("backup_archive_repo.GetByID: %w", err)
	}

	row, err := r.q.GetBackupArchiveByID(ctx, generated.GetBackupArchiveByIDParams{ID: bid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("backup_archive_repo.GetByID: %w", err)
	}

	return backupArchiveRecordFromRow(row), nil
}

func (r *backupArchiveRepository) UpdateFilePath(ctx context.Context, id string, filePath *string) (*models.BackupArchiveRecord, error) {
	bid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("backup_archive_repo.UpdateFilePath: %w", err)
	}

	row, err := r.q.UpdateBackupArchiveFilePath(ctx, generated.UpdateBackupArchiveFilePathParams{
		ID:          bid,
		ArchivePath: nullableText(filePath),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("backup_archive_repo.UpdateFilePath: %w", err)
	}

	return backupArchiveRecordFromRow(row), nil
}

func backupArchiveRecordFromRow(row generated.BackupArchive) *models.BackupArchiveRecord {
	return &models.BackupArchiveRecord{
		ID:           row.ID.String(),
		UserID:       row.UserID.String(),
		IncludeMedia: row.IncludeMedia,
		SizeBytes:    row.SizeBytes,
		EntityCounts: string(row.EntityCounts),
		ArchivePath:  textPtr(row.ArchivePath),
		CreatedAt:    formatTimestamp(row.CreatedAt),
		UpdatedAt:    formatTimestamp(row.UpdatedAt),
	}
}
