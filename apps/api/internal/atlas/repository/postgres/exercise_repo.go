// FILE: apps/api/internal/atlas/repository/postgres/exercise_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement the ExerciseRepository interface for WAVE-02 Exercise Library using sqlc-generated queries.
//   SCOPE: Create, GetByID, List, ListAll, Update, Archive, Restore for exercises; CreateMedia, GetMediaByID, ListMediaByExercise, DeleteMedia for exercise media. All operations scoped by userID.
//   DEPENDS: apps/api/internal/repository/postgres/generated (sqlc), apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-02.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ExerciseRepository - Interface for exercise and exercise media data access.
//   NewExerciseRepository - Creates a new ExerciseRepository with sqlc-backed Queries.
//   Create - Creates a new exercise, returns ExerciseRecord.
//   GetByID - Gets exercise by ID and userID, returns ExerciseRecord or not-found error.
//   List - Lists exercises by userID and isActive with limit.
//   ListCursor - Lists exercises by userID, isActive, name cursor with limit.
//   ListAll - Lists all exercises by userID with optional inactive filter.
//   Count - Counts exercises by userID and isActive.
//   Update - Updates exercise fields, returns ExerciseRecord.
//   Archive - Sets is_active=false, returns ExerciseRecord.
//   Restore - Sets is_active=true, returns ExerciseRecord.
//   CreateMedia - Creates exercise media, returns ExerciseMedia record.
//   GetMediaByID - Gets exercise media by ID and userID.
//   ListMediaByExercise - Lists exercise media by exerciseID and userID.
//   DeleteMedia - Deletes exercise media by ID and userID, returns deleted record.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added exercise repository for WAVE-02.
// END_CHANGE_SUMMARY

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/repository/postgres/generated"
)

type ExerciseRepository interface {
	Create(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.ExerciseRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error)
	List(ctx context.Context, userID string, isActive bool, limit int32) ([]models.ExerciseRecord, error)
	ListCursor(ctx context.Context, userID string, isActive bool, cursor string, limit int32) ([]models.ExerciseRecord, error)
	ListAll(ctx context.Context, userID string, includeInactive bool) ([]models.ExerciseRecord, error)
	Count(ctx context.Context, userID string, isActive bool) (int, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.ExerciseRecord, error)
	Archive(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error)
	Restore(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error)
	CreateMedia(ctx context.Context, userID string, exerciseID string, fileName string, filePath string, mimeType string, fileSize int64) (*models.ExerciseMedia, error)
GetMediaByID(ctx context.Context, userID string, id string) (*models.ExerciseMedia, error)
	ListMediaByExercise(ctx context.Context, userID string, exerciseID string) ([]models.ExerciseMedia, error)
	DeleteMedia(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error)
	GetMediaRecordByID(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error)
}

type exerciseRepository struct {
	q *generated.Queries
}

func NewExerciseRepository(pool *pgxpool.Pool) ExerciseRepository {
	return &exerciseRepository{q: generated.New(pool)}
}

func (r *exerciseRepository) Create(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.ExerciseRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.Create: %w", err)
	}

	row, err := r.q.CreateExercise(ctx, generated.CreateExerciseParams{
		UserID:        uid,
		Name:          input.Name,
		MuscleGroups:  input.MuscleGroups,
		Description:   nullableText(input.Description),
		PersonalNotes: nullableText(input.PersonalNotes),
		WorkingWeight: nullableFloat4(input.WorkingWeight),
	})
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.Create: %w", err)
	}

	return exerciseRecordFromRow(row), nil
}

func (r *exerciseRepository) GetByID(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
	uid, eid, err := r.parseIDs(id, userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.GetByID: %w", err)
	}

	row, err := r.q.GetExerciseByID(ctx, generated.GetExerciseByIDParams{ID: eid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exercise_repo.GetByID: %w", err)
	}

	return exerciseRecordFromRow(row), nil
}

func (r *exerciseRepository) List(ctx context.Context, userID string, isActive bool, limit int32) ([]models.ExerciseRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.List: %w", err)
	}

	rows, err := r.q.ListExercises(ctx, generated.ListExercisesParams{
		UserID:   uid,
		IsActive: isActive,
		Limit:    limit,
	})
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.List: %w", err)
	}

	return exerciseRecordsFromRows(rows), nil
}

func (r *exerciseRepository) ListCursor(ctx context.Context, userID string, isActive bool, cursor string, limit int32) ([]models.ExerciseRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.ListCursor: %w", err)
	}

	rows, err := r.q.ListExercisesCursor(ctx, generated.ListExercisesCursorParams{
		UserID:   uid,
		IsActive: isActive,
		Name:     cursor,
		Limit:    limit,
	})
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.ListCursor: %w", err)
	}

	return exerciseRecordsFromRows(rows), nil
}

func (r *exerciseRepository) ListAll(ctx context.Context, userID string, includeInactive bool) ([]models.ExerciseRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.ListAll: %w", err)
	}

	rows, err := r.q.ListAllExercises(ctx, generated.ListAllExercisesParams{
		UserID:  uid,
		Column2: !includeInactive,
	})
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.ListAll: %w", err)
	}

	return exerciseRecordsFromRows(rows), nil
}

func (r *exerciseRepository) Count(ctx context.Context, userID string, isActive bool) (int, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return 0, fmt.Errorf("exercise_repo.Count: %w", err)
	}

	count, err := r.q.CountExercises(ctx, generated.CountExercisesParams{
		UserID:   uid,
		IsActive: isActive,
	})
	if err != nil {
		return 0, fmt.Errorf("exercise_repo.Count: %w", err)
	}

	return int(count), nil
}

func (r *exerciseRepository) Update(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.ExerciseRecord, error) {
	uid, eid, err := r.parseIDs(id, userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.Update: %w", err)
	}

	name := ""
	if input.Name != nil {
		name = *input.Name
	}

	var muscleGroups []string
	if input.MuscleGroups != nil {
		muscleGroups = *input.MuscleGroups
	}

	row, err := r.q.UpdateExercise(ctx, generated.UpdateExerciseParams{
		ID:            eid,
		UserID:        uid,
		Name:          name,
		MuscleGroups:  muscleGroups,
		Description:   nullableText(input.Description),
		PersonalNotes: nullableText(input.PersonalNotes),
		WorkingWeight: nullableFloat4(input.WorkingWeight),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exercise_repo.Update: %w", err)
	}

	return exerciseRecordFromRow(row), nil
}

func (r *exerciseRepository) Archive(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
	uid, eid, err := r.parseIDs(id, userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.Archive: %w", err)
	}

	row, err := r.q.ArchiveExercise(ctx, generated.ArchiveExerciseParams{ID: eid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exercise_repo.Archive: %w", err)
	}

	return exerciseRecordFromRow(row), nil
}

func (r *exerciseRepository) Restore(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
	uid, eid, err := r.parseIDs(id, userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.Restore: %w", err)
	}

	row, err := r.q.RestoreExercise(ctx, generated.RestoreExerciseParams{ID: eid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exercise_repo.Restore: %w", err)
	}

	return exerciseRecordFromRow(row), nil
}

func (r *exerciseRepository) CreateMedia(ctx context.Context, userID string, exerciseID string, fileName string, filePath string, mimeType string, fileSize int64) (*models.ExerciseMedia, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.CreateMedia: %w", err)
	}
	eid, err := uuidFromString(exerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.CreateMedia: %w", err)
	}

	row, err := r.q.CreateExerciseMedia(ctx, generated.CreateExerciseMediaParams{
		UserID:     uid,
		ExerciseID: eid,
		FileName:   fileName,
		FilePath:   filePath,
		MimeType:   mimeType,
		FileSize:   fileSize,
	})
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.CreateMedia: %w", err)
	}

	return exerciseMediaFromRow(row), nil
}

func (r *exerciseRepository) GetMediaByID(ctx context.Context, userID string, id string) (*models.ExerciseMedia, error) {
	uid, mid, err := r.parseIDs(id, userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.GetMediaByID: %w", err)
	}

	row, err := r.q.GetExerciseMediaByID(ctx, generated.GetExerciseMediaByIDParams{ID: mid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exercise_repo.GetMediaByID: %w", err)
	}

	return exerciseMediaFromRow(row), nil
}

func (r *exerciseRepository) GetMediaRecordByID(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error) {
	uid, mid, err := r.parseIDs(id, userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.GetMediaRecordByID: %w", err)
	}

	row, err := r.q.GetExerciseMediaByID(ctx, generated.GetExerciseMediaByIDParams{ID: mid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exercise_repo.GetMediaRecordByID: %w", err)
	}

	return exerciseMediaRecordFromRow(row), nil
}

func (r *exerciseRepository) ListMediaByExercise(ctx context.Context, userID string, exerciseID string) ([]models.ExerciseMedia, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.ListMediaByExercise: %w", err)
	}
	eid, err := uuidFromString(exerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.ListMediaByExercise: %w", err)
	}

	rows, err := r.q.ListExerciseMediaByExercise(ctx, generated.ListExerciseMediaByExerciseParams{
		ExerciseID: eid,
		UserID:     uid,
	})
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.ListMediaByExercise: %w", err)
	}

	out := make([]models.ExerciseMedia, len(rows))
	for i, row := range rows {
		out[i] = *exerciseMediaFromRow(row)
	}
	return out, nil
}

func (r *exerciseRepository) DeleteMedia(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error) {
	uid, mid, err := r.parseIDs(id, userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.DeleteMedia: %w", err)
	}

	row, err := r.q.DeleteExerciseMedia(ctx, generated.DeleteExerciseMediaParams{ID: mid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exercise_repo.DeleteMedia: %w", err)
	}

	return exerciseMediaRecordFromRow(row), nil
}

func (r *exerciseRepository) parseIDs(id string, userID string) (pgtype.UUID, pgtype.UUID, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return pgtype.UUID{}, pgtype.UUID{}, err
	}
	eid, err := uuidFromString(id)
	if err != nil {
		return pgtype.UUID{}, pgtype.UUID{}, err
	}
	return uid, eid, nil
}

func exerciseRecordFromRow(row generated.Exercise) *models.ExerciseRecord {
	return &models.ExerciseRecord{
		ID:            row.ID.String(),
		UserID:        row.UserID.String(),
		Name:          row.Name,
		MuscleGroups:  row.MuscleGroups,
		Description:   textPtr(row.Description),
		PersonalNotes: textPtr(row.PersonalNotes),
		WorkingWeight: float4Ptr(row.WorkingWeight),
		IsActive:      row.IsActive,
		CreatedAt:     formatTimestamp(row.CreatedAt),
		UpdatedAt:     formatTimestamp(row.UpdatedAt),
	}
}

func exerciseRecordsFromRows(rows []generated.Exercise) []models.ExerciseRecord {
	out := make([]models.ExerciseRecord, len(rows))
	for i, row := range rows {
		out[i] = *exerciseRecordFromRow(row)
	}
	return out
}

func exerciseMediaFromRow(row generated.ExerciseMedium) *models.ExerciseMedia {
	return &models.ExerciseMedia{
		ID:         row.ID.String(),
		UserID:     row.UserID.String(),
		ExerciseID: row.ExerciseID.String(),
		FileName:   row.FileName,
		MimeType:   row.MimeType,
		FileSize:   row.FileSize,
		CreatedAt:  formatTimestamp(row.CreatedAt),
	}
}

func exerciseMediaRecordFromRow(row generated.ExerciseMedium) *models.ExerciseMediaRecord {
	return &models.ExerciseMediaRecord{
		ID:         row.ID.String(),
		UserID:     row.UserID.String(),
		ExerciseID: row.ExerciseID.String(),
		FileName:   row.FileName,
		FilePath:   row.FilePath,
		MimeType:   row.MimeType,
		FileSize:   row.FileSize,
		CreatedAt:  formatTimestamp(row.CreatedAt),
	}
}

func textPtr(t pgtype.Text) *string {
	if t.Valid {
		return &t.String
	}
	return nil
}

func float4Ptr(f pgtype.Float4) *float64 {
	if f.Valid {
		v := float64(f.Float32)
		return &v
	}
	return nil
}

func nullableFloat4(v *float64) pgtype.Float4 {
	if v == nil {
		return pgtype.Float4{}
	}
	return pgtype.Float4{Float32: float32(*v), Valid: true}
}
