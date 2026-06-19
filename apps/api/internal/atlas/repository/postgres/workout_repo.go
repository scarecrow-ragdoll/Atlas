// FILE: apps/api/internal/atlas/repository/postgres/workout_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement WAVE-03 DailyLog aggregate repository operations with sqlc and pgx transactions.
//   SCOPE: DailyLog get/create/lock/version, workout exercise CRUD/reorder, workout set CRUD/reorder, and aggregate reads; excludes service validation and GraphQL mapping.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/atlas/models, pgx/v5.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   WorkoutRepository - Data access interface for WAVE-03 DailyLog aggregate operations.
//   NewWorkoutRepository - Creates the sqlc-backed repository.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.3 - Routed update moves through transactional reorder helpers and guarded no-op version increments.
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

type WorkoutRepository interface {
	GetDailyLogByDate(ctx context.Context, userID string, date models.Date) (*DailyLogRecord, error)
	GetOrCreateDailyLogByDate(ctx context.Context, userID string, date models.Date) (*DailyLogRecord, error)
	GetDailyLogAggregate(ctx context.Context, userID string, dailyLogID string) (*DailyLogAggregate, error)
	ListDailyLogSummaries(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]DailyLogSummaryRecord, error)
	WithLockedDailyLogByDate(ctx context.Context, userID string, date models.Date, fn LockedDailyLogFunc) error
	WithLockedDailyLogByWorkoutExerciseID(ctx context.Context, userID string, workoutExerciseID string, fn LockedDailyLogFunc) error
	WithLockedDailyLogByWorkoutSetID(ctx context.Context, userID string, workoutSetID string, fn LockedDailyLogFunc) error
	IncrementDailyLogVersion(ctx context.Context, userID string, dailyLogID string) (*DailyLogRecord, error)
	UpdateDailyLogNotes(ctx context.Context, userID string, dailyLogID string, notes *string) (*DailyLogRecord, error)
	AddWorkoutExercise(ctx context.Context, userID string, dailyLogID string, input AddWorkoutExerciseInput) (*WorkoutExerciseRecord, error)
	UpdateWorkoutExercise(ctx context.Context, userID string, workoutExerciseID string, input UpdateWorkoutExerciseInput) (*WorkoutExerciseRecord, error)
	DeleteWorkoutExercise(ctx context.Context, userID string, workoutExerciseID string) (*WorkoutExerciseRecord, error)
	ReorderWorkoutExercises(ctx context.Context, userID string, dailyLogID string, orderedIDs []string) error
	AddWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, input AddWorkoutSetInput) (*WorkoutSetRecord, error)
	UpdateWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, workoutSetID string, input UpdateWorkoutSetInput) (*WorkoutSetRecord, error)
	DeleteWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, workoutSetID string) (*WorkoutSetRecord, error)
	ReorderWorkoutSets(ctx context.Context, userID string, workoutExerciseID string, orderedIDs []string) error
}

type WorkoutTx interface {
	GetDailyLogAggregate(ctx context.Context, userID string, dailyLogID string) (*DailyLogAggregate, error)
	IncrementDailyLogVersion(ctx context.Context, userID string, dailyLogID string) (*DailyLogRecord, error)
	UpdateDailyLogNotes(ctx context.Context, userID string, dailyLogID string, notes *string) (*DailyLogRecord, error)
	AddWorkoutExercise(ctx context.Context, userID string, dailyLogID string, input AddWorkoutExerciseInput) (*WorkoutExerciseRecord, error)
	UpdateWorkoutExercise(ctx context.Context, userID string, workoutExerciseID string, input UpdateWorkoutExerciseInput) (*WorkoutExerciseRecord, error)
	DeleteWorkoutExercise(ctx context.Context, userID string, workoutExerciseID string) (*WorkoutExerciseRecord, error)
	ReorderWorkoutExercises(ctx context.Context, userID string, dailyLogID string, orderedIDs []string) error
	AddWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, input AddWorkoutSetInput) (*WorkoutSetRecord, error)
	UpdateWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, workoutSetID string, input UpdateWorkoutSetInput) (*WorkoutSetRecord, error)
	DeleteWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, workoutSetID string) (*WorkoutSetRecord, error)
	ReorderWorkoutSets(ctx context.Context, userID string, workoutExerciseID string, orderedIDs []string) error
}

type LockedDailyLogFunc func(ctx context.Context, tx WorkoutTx, dailyLog *DailyLogRecord) error

type DailyLogRecord struct {
	ID        string
	UserID    string
	Date      models.Date
	Notes     *string
	Version   int32
	CreatedAt string
	UpdatedAt string
}

type DailyLogSummaryRecord struct {
	ID                   string
	Date                 models.Date
	Version              int32
	WorkoutExerciseCount int32
	WorkoutSetCount      int32
	TotalVolume          float64
	UpdatedAt            string
}

type DailyLogAggregate struct {
	DailyLog         DailyLogRecord
	WorkoutExercises []WorkoutExerciseRecord
}

type WorkoutExerciseRecord struct {
	ID                    string
	UserID                string
	DailyLogID            string
	ExerciseID            string
	Position              int32
	WorkingWeightSnapshot *float64
	Notes                 *string
	Sets                  []WorkoutSetRecord
	CreatedAt             string
	UpdatedAt             string
}

type WorkoutSetRecord struct {
	ID                string
	WorkoutExerciseID string
	SetNumber         int32
	Weight            float64
	Reps              int32
	RPE               *float64
	RIR               *int32
	Notes             *string
	CreatedAt         string
	UpdatedAt         string
}

type AddWorkoutExerciseInput struct {
	ExerciseID            string
	Position              int32
	WorkingWeightSnapshot *float64
	Notes                 *string
}

type UpdateWorkoutExerciseInput struct {
	Position *int32
	SetNotes bool
	Notes    *string
}

type AddWorkoutSetInput struct {
	SetNumber int32
	Weight    float64
	Reps      int32
	RPE       *float64
	RIR       *int32
	Notes     *string
}

type UpdateWorkoutSetInput struct {
	SetNumber *int32
	Weight    *float64
	Reps      *int32
	SetRPE    bool
	RPE       *float64
	SetRIR    bool
	RIR       *int32
	SetNotes  bool
	Notes     *string
}

type workoutRepository struct {
	pool *pgxpool.Pool
	q    *generated.Queries
}

type workoutTx struct {
	q *generated.Queries
}

func NewWorkoutRepository(pool *pgxpool.Pool) WorkoutRepository {
	return &workoutRepository{
		pool: pool,
		q:    generated.New(pool),
	}
}

func (r *workoutRepository) GetDailyLogByDate(ctx context.Context, userID string, date models.Date) (*DailyLogRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("workout_repo.GetDailyLogByDate: %w", err)
	}

	row, err := r.q.GetDailyLogByDate(ctx, generated.GetDailyLogByDateParams{
		UserID: uid,
		Date:   dateParam(date),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("workout_repo.GetDailyLogByDate: %w", err)
	}
	return dailyLogRecordFromRow(row), nil
}

func (r *workoutRepository) GetOrCreateDailyLogByDate(ctx context.Context, userID string, date models.Date) (*DailyLogRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("workout_repo.GetOrCreateDailyLogByDate: %w", err)
	}

	var out *DailyLogRecord
	err = r.withTx(ctx, func(q *generated.Queries) error {
		row, err := q.CreateDailyLog(ctx, generated.CreateDailyLogParams{
			UserID: uid,
			Date:   dateParam(date),
			Notes:  pgtype.Text{},
		})
		if err != nil {
			return err
		}
		out = dailyLogRecordFromRow(row)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("workout_repo.GetOrCreateDailyLogByDate: %w", err)
	}
	return out, nil
}

func (r *workoutRepository) GetDailyLogAggregate(ctx context.Context, userID string, dailyLogID string) (*DailyLogAggregate, error) {
	return getDailyLogAggregate(ctx, r.q, userID, dailyLogID)
}

func (r *workoutRepository) ListDailyLogSummaries(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]DailyLogSummaryRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("workout_repo.ListDailyLogSummaries: %w", err)
	}

	rows, err := r.q.ListDailyLogSummaries(ctx, generated.ListDailyLogSummariesParams{
		UserID:   uid,
		FromDate: dateParam(fromDate),
		ToDate:   dateParam(toDate),
	})
	if err != nil {
		return nil, fmt.Errorf("workout_repo.ListDailyLogSummaries: %w", err)
	}

	out := make([]DailyLogSummaryRecord, len(rows))
	for i, row := range rows {
		out[i] = DailyLogSummaryRecord{
			ID:                   row.ID.String(),
			Date:                 dateFromRow(row.Date),
			Version:              row.Version,
			WorkoutExerciseCount: row.WorkoutExerciseCount,
			WorkoutSetCount:      row.WorkoutSetCount,
			TotalVolume:          row.TotalVolume,
			UpdatedAt:            formatTimestamp(row.UpdatedAt),
		}
	}
	return out, nil
}

func (r *workoutRepository) WithLockedDailyLogByDate(ctx context.Context, userID string, date models.Date, fn LockedDailyLogFunc) error {
	uid, err := uuidFromString(userID)
	if err != nil {
		return fmt.Errorf("workout_repo.WithLockedDailyLogByDate: %w", err)
	}

	err = r.withTx(ctx, func(q *generated.Queries) error {
		row, err := q.LockDailyLogByDate(ctx, generated.LockDailyLogByDateParams{UserID: uid, Date: dateParam(date)})
		if err != nil {
			return err
		}
		return fn(ctx, &workoutTx{q: q}, dailyLogRecordFromRow(row))
	})
	if err != nil {
		return fmt.Errorf("workout_repo.WithLockedDailyLogByDate: %w", err)
	}
	return nil
}

func (r *workoutRepository) WithLockedDailyLogByWorkoutExerciseID(ctx context.Context, userID string, workoutExerciseID string, fn LockedDailyLogFunc) error {
	uid, weid, err := parseTwoUUIDs(userID, workoutExerciseID)
	if err != nil {
		return fmt.Errorf("workout_repo.WithLockedDailyLogByWorkoutExerciseID: %w", err)
	}

	err = r.withTx(ctx, func(q *generated.Queries) error {
		row, err := q.LockDailyLogByWorkoutExerciseID(ctx, generated.LockDailyLogByWorkoutExerciseIDParams{
			UserID:            uid,
			WorkoutExerciseID: weid,
		})
		if err != nil {
			return err
		}
		return fn(ctx, &workoutTx{q: q}, dailyLogRecordFromRow(row))
	})
	if err != nil {
		return fmt.Errorf("workout_repo.WithLockedDailyLogByWorkoutExerciseID: %w", err)
	}
	return nil
}

func (r *workoutRepository) WithLockedDailyLogByWorkoutSetID(ctx context.Context, userID string, workoutSetID string, fn LockedDailyLogFunc) error {
	uid, wsid, err := parseTwoUUIDs(userID, workoutSetID)
	if err != nil {
		return fmt.Errorf("workout_repo.WithLockedDailyLogByWorkoutSetID: %w", err)
	}

	err = r.withTx(ctx, func(q *generated.Queries) error {
		row, err := q.LockDailyLogByWorkoutSetID(ctx, generated.LockDailyLogByWorkoutSetIDParams{
			UserID:       uid,
			WorkoutSetID: wsid,
		})
		if err != nil {
			return err
		}
		return fn(ctx, &workoutTx{q: q}, dailyLogRecordFromRow(row))
	})
	if err != nil {
		return fmt.Errorf("workout_repo.WithLockedDailyLogByWorkoutSetID: %w", err)
	}
	return nil
}

func (r *workoutRepository) IncrementDailyLogVersion(ctx context.Context, userID string, dailyLogID string) (*DailyLogRecord, error) {
	return incrementDailyLogVersion(ctx, r.q, userID, dailyLogID, "workout_repo.IncrementDailyLogVersion")
}

func (r *workoutRepository) UpdateDailyLogNotes(ctx context.Context, userID string, dailyLogID string, notes *string) (*DailyLogRecord, error) {
	var out *DailyLogRecord
	err := r.withTx(ctx, func(q *generated.Queries) error {
		locked, err := lockDailyLogByID(ctx, q, userID, dailyLogID)
		if err != nil || locked == nil {
			return err
		}
		if _, err := q.UpdateDailyLogNotes(ctx, generated.UpdateDailyLogNotesParams{
			Notes:  nullableText(notes),
			UserID: mustUUID(userID),
			ID:     mustUUID(dailyLogID),
		}); err != nil {
			return err
		}
		incremented, err := incrementDailyLogVersionInTx(ctx, q, userID, dailyLogID)
		if err != nil {
			return err
		}
		out = dailyLogRecordFromRow(incremented)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("workout_repo.UpdateDailyLogNotes: %w", err)
	}
	return out, nil
}

func (r *workoutRepository) AddWorkoutExercise(ctx context.Context, userID string, dailyLogID string, input AddWorkoutExerciseInput) (*WorkoutExerciseRecord, error) {
	var out *WorkoutExerciseRecord
	err := r.withTx(ctx, func(q *generated.Queries) error {
		locked, err := lockDailyLogByID(ctx, q, userID, dailyLogID)
		if err != nil || locked == nil {
			return err
		}
		out, err = addWorkoutExerciseInTx(ctx, q, userID, dailyLogID, input)
		if err != nil {
			return err
		}
		_, err = incrementDailyLogVersionInTx(ctx, q, userID, dailyLogID)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("workout_repo.AddWorkoutExercise: %w", err)
	}
	return out, nil
}

func (r *workoutRepository) UpdateWorkoutExercise(ctx context.Context, userID string, workoutExerciseID string, input UpdateWorkoutExerciseInput) (*WorkoutExerciseRecord, error) {
	var out *WorkoutExerciseRecord
	err := r.withTx(ctx, func(q *generated.Queries) error {
		locked, err := lockDailyLogByWorkoutExerciseID(ctx, q, userID, workoutExerciseID)
		if err != nil || locked == nil {
			return err
		}
		var changed bool
		out, changed, err = updateWorkoutExerciseInTx(ctx, q, userID, locked.ID, workoutExerciseID, input)
		if err != nil {
			return err
		}
		if !changed || out == nil {
			return nil
		}
		_, err = incrementDailyLogVersionInTx(ctx, q, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("workout_repo.UpdateWorkoutExercise: %w", err)
	}
	return out, nil
}

func (r *workoutRepository) DeleteWorkoutExercise(ctx context.Context, userID string, workoutExerciseID string) (*WorkoutExerciseRecord, error) {
	var out *WorkoutExerciseRecord
	err := r.withTx(ctx, func(q *generated.Queries) error {
		locked, err := lockDailyLogByWorkoutExerciseID(ctx, q, userID, workoutExerciseID)
		if err != nil || locked == nil {
			return err
		}
		out, err = deleteWorkoutExerciseInTx(ctx, q, userID, workoutExerciseID)
		if err != nil || out == nil {
			return err
		}
		_, err = q.TempShiftWorkoutExercisePositionsAfterDelete(ctx, generated.TempShiftWorkoutExercisePositionsAfterDeleteParams{
			UserID:          mustUUID(userID),
			DailyLogID:      mustUUID(out.DailyLogID),
			DeletedPosition: out.Position,
		})
		if err != nil {
			return err
		}
		_, err = q.NormalizeWorkoutExercisePositionsAfterDelete(ctx, generated.NormalizeWorkoutExercisePositionsAfterDeleteParams{
			UserID:          mustUUID(userID),
			DailyLogID:      mustUUID(out.DailyLogID),
			DeletedPosition: out.Position,
		})
		if err != nil {
			return err
		}
		_, err = incrementDailyLogVersionInTx(ctx, q, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("workout_repo.DeleteWorkoutExercise: %w", err)
	}
	return out, nil
}

func (r *workoutRepository) ReorderWorkoutExercises(ctx context.Context, userID string, dailyLogID string, orderedIDs []string) error {
	err := r.withTx(ctx, func(q *generated.Queries) error {
		locked, err := lockDailyLogByID(ctx, q, userID, dailyLogID)
		if err != nil || locked == nil {
			return err
		}
		if err := reorderWorkoutExercisesInTx(ctx, q, userID, dailyLogID, orderedIDs); err != nil {
			return err
		}
		_, err = incrementDailyLogVersionInTx(ctx, q, userID, dailyLogID)
		return err
	})
	if err != nil {
		return fmt.Errorf("workout_repo.ReorderWorkoutExercises: %w", err)
	}
	return nil
}

func (r *workoutRepository) AddWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, input AddWorkoutSetInput) (*WorkoutSetRecord, error) {
	var out *WorkoutSetRecord
	err := r.withTx(ctx, func(q *generated.Queries) error {
		locked, err := lockDailyLogByWorkoutExerciseID(ctx, q, userID, workoutExerciseID)
		if err != nil || locked == nil {
			return err
		}
		out, err = addWorkoutSetInTx(ctx, q, workoutExerciseID, input)
		if err != nil {
			return err
		}
		_, err = incrementDailyLogVersionInTx(ctx, q, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("workout_repo.AddWorkoutSet: %w", err)
	}
	return out, nil
}

func (r *workoutRepository) UpdateWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, workoutSetID string, input UpdateWorkoutSetInput) (*WorkoutSetRecord, error) {
	var out *WorkoutSetRecord
	err := r.withTx(ctx, func(q *generated.Queries) error {
		locked, err := lockDailyLogByWorkoutSetID(ctx, q, userID, workoutSetID)
		if err != nil || locked == nil {
			return err
		}
		var changed bool
		out, changed, err = updateWorkoutSetInTx(ctx, q, workoutExerciseID, workoutSetID, input)
		if err != nil {
			return err
		}
		if !changed || out == nil {
			return nil
		}
		_, err = incrementDailyLogVersionInTx(ctx, q, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("workout_repo.UpdateWorkoutSet: %w", err)
	}
	return out, nil
}

func (r *workoutRepository) DeleteWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, workoutSetID string) (*WorkoutSetRecord, error) {
	var out *WorkoutSetRecord
	err := r.withTx(ctx, func(q *generated.Queries) error {
		locked, err := lockDailyLogByWorkoutSetID(ctx, q, userID, workoutSetID)
		if err != nil || locked == nil {
			return err
		}
		out, err = deleteWorkoutSetInTx(ctx, q, workoutExerciseID, workoutSetID)
		if err != nil || out == nil {
			return err
		}
		_, err = q.TempShiftWorkoutSetNumbersAfterDelete(ctx, generated.TempShiftWorkoutSetNumbersAfterDeleteParams{
			WorkoutExerciseID: mustUUID(workoutExerciseID),
			DeletedSetNumber:  out.SetNumber,
		})
		if err != nil {
			return err
		}
		_, err = q.NormalizeWorkoutSetNumbersAfterDelete(ctx, generated.NormalizeWorkoutSetNumbersAfterDeleteParams{
			WorkoutExerciseID: mustUUID(workoutExerciseID),
			DeletedSetNumber:  out.SetNumber,
		})
		if err != nil {
			return err
		}
		_, err = incrementDailyLogVersionInTx(ctx, q, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("workout_repo.DeleteWorkoutSet: %w", err)
	}
	return out, nil
}

func (r *workoutRepository) ReorderWorkoutSets(ctx context.Context, userID string, workoutExerciseID string, orderedIDs []string) error {
	err := r.withTx(ctx, func(q *generated.Queries) error {
		locked, err := lockDailyLogByWorkoutExerciseID(ctx, q, userID, workoutExerciseID)
		if err != nil || locked == nil {
			return err
		}
		if err := reorderWorkoutSetsInTx(ctx, q, workoutExerciseID, orderedIDs); err != nil {
			return err
		}
		_, err = incrementDailyLogVersionInTx(ctx, q, userID, locked.ID)
		return err
	})
	if err != nil {
		return fmt.Errorf("workout_repo.ReorderWorkoutSets: %w", err)
	}
	return nil
}

func (r *workoutRepository) withTx(ctx context.Context, fn func(q *generated.Queries) error) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(generated.New(tx)); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (tx *workoutTx) GetDailyLogAggregate(ctx context.Context, userID string, dailyLogID string) (*DailyLogAggregate, error) {
	return getDailyLogAggregate(ctx, tx.q, userID, dailyLogID)
}

func (tx *workoutTx) IncrementDailyLogVersion(ctx context.Context, userID string, dailyLogID string) (*DailyLogRecord, error) {
	return incrementDailyLogVersion(ctx, tx.q, userID, dailyLogID, "workout_tx.IncrementDailyLogVersion")
}

func (tx *workoutTx) UpdateDailyLogNotes(ctx context.Context, userID string, dailyLogID string, notes *string) (*DailyLogRecord, error) {
	row, err := tx.q.UpdateDailyLogNotes(ctx, generated.UpdateDailyLogNotesParams{
		Notes:  nullableText(notes),
		UserID: mustUUID(userID),
		ID:     mustUUID(dailyLogID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("workout_tx.UpdateDailyLogNotes: %w", err)
	}
	return dailyLogRecordFromRow(row), nil
}

func (tx *workoutTx) AddWorkoutExercise(ctx context.Context, userID string, dailyLogID string, input AddWorkoutExerciseInput) (*WorkoutExerciseRecord, error) {
	return addWorkoutExerciseInTx(ctx, tx.q, userID, dailyLogID, input)
}

func (tx *workoutTx) UpdateWorkoutExercise(ctx context.Context, userID string, workoutExerciseID string, input UpdateWorkoutExerciseInput) (*WorkoutExerciseRecord, error) {
	locked, err := lockDailyLogByWorkoutExerciseID(ctx, tx.q, userID, workoutExerciseID)
	if err != nil || locked == nil {
		return nil, err
	}
	out, _, err := updateWorkoutExerciseInTx(ctx, tx.q, userID, locked.ID, workoutExerciseID, input)
	return out, err
}

func (tx *workoutTx) DeleteWorkoutExercise(ctx context.Context, userID string, workoutExerciseID string) (*WorkoutExerciseRecord, error) {
	return deleteWorkoutExerciseInTx(ctx, tx.q, userID, workoutExerciseID)
}

func (tx *workoutTx) ReorderWorkoutExercises(ctx context.Context, userID string, dailyLogID string, orderedIDs []string) error {
	return reorderWorkoutExercisesInTx(ctx, tx.q, userID, dailyLogID, orderedIDs)
}

func (tx *workoutTx) AddWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, input AddWorkoutSetInput) (*WorkoutSetRecord, error) {
	return addWorkoutSetInTx(ctx, tx.q, workoutExerciseID, input)
}

func (tx *workoutTx) UpdateWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, workoutSetID string, input UpdateWorkoutSetInput) (*WorkoutSetRecord, error) {
	out, _, err := updateWorkoutSetInTx(ctx, tx.q, workoutExerciseID, workoutSetID, input)
	return out, err
}

func (tx *workoutTx) DeleteWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, workoutSetID string) (*WorkoutSetRecord, error) {
	return deleteWorkoutSetInTx(ctx, tx.q, workoutExerciseID, workoutSetID)
}

func (tx *workoutTx) ReorderWorkoutSets(ctx context.Context, userID string, workoutExerciseID string, orderedIDs []string) error {
	return reorderWorkoutSetsInTx(ctx, tx.q, workoutExerciseID, orderedIDs)
}

func getDailyLogAggregate(ctx context.Context, q *generated.Queries, userID string, dailyLogID string) (*DailyLogAggregate, error) {
	uid, dlid, err := parseTwoUUIDs(userID, dailyLogID)
	if err != nil {
		return nil, fmt.Errorf("workout_repo.GetDailyLogAggregate: %w", err)
	}

	dailyLog, err := q.GetDailyLogByID(ctx, generated.GetDailyLogByIDParams{UserID: uid, ID: dlid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("workout_repo.GetDailyLogAggregate: %w", err)
	}

	exerciseRows, err := q.ListWorkoutExercisesByDailyLog(ctx, generated.ListWorkoutExercisesByDailyLogParams{
		UserID:     uid,
		DailyLogID: dlid,
	})
	if err != nil {
		return nil, fmt.Errorf("workout_repo.GetDailyLogAggregate: %w", err)
	}

	exercises := workoutExerciseRecordsFromRows(exerciseRows)
	if len(exercises) > 0 {
		ids := make([]pgtype.UUID, len(exercises))
		for i, exercise := range exercises {
			ids[i] = mustUUID(exercise.ID)
		}
		setRows, err := q.ListWorkoutSetsByExerciseIDs(ctx, ids)
		if err != nil {
			return nil, fmt.Errorf("workout_repo.GetDailyLogAggregate: %w", err)
		}
		setsByExercise := make(map[string][]WorkoutSetRecord)
		for _, row := range setRows {
			set := workoutSetRecordFromRow(row)
			setsByExercise[set.WorkoutExerciseID] = append(setsByExercise[set.WorkoutExerciseID], *set)
		}
		for i := range exercises {
			exercises[i].Sets = setsByExercise[exercises[i].ID]
		}
	}

	return &DailyLogAggregate{
		DailyLog:         *dailyLogRecordFromRow(dailyLog),
		WorkoutExercises: exercises,
	}, nil
}

func lockDailyLogByID(ctx context.Context, q *generated.Queries, userID string, dailyLogID string) (*DailyLogRecord, error) {
	uid, dlid, err := parseTwoUUIDs(userID, dailyLogID)
	if err != nil {
		return nil, err
	}
	row, err := q.LockDailyLogByID(ctx, generated.LockDailyLogByIDParams{UserID: uid, ID: dlid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return dailyLogRecordFromRow(row), nil
}

func lockDailyLogByWorkoutExerciseID(ctx context.Context, q *generated.Queries, userID string, workoutExerciseID string) (*DailyLogRecord, error) {
	uid, weid, err := parseTwoUUIDs(userID, workoutExerciseID)
	if err != nil {
		return nil, err
	}
	row, err := q.LockDailyLogByWorkoutExerciseID(ctx, generated.LockDailyLogByWorkoutExerciseIDParams{
		UserID:            uid,
		WorkoutExerciseID: weid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return dailyLogRecordFromRow(row), nil
}

func lockDailyLogByWorkoutSetID(ctx context.Context, q *generated.Queries, userID string, workoutSetID string) (*DailyLogRecord, error) {
	uid, wsid, err := parseTwoUUIDs(userID, workoutSetID)
	if err != nil {
		return nil, err
	}
	row, err := q.LockDailyLogByWorkoutSetID(ctx, generated.LockDailyLogByWorkoutSetIDParams{UserID: uid, WorkoutSetID: wsid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return dailyLogRecordFromRow(row), nil
}

func addWorkoutExerciseInTx(ctx context.Context, q *generated.Queries, userID string, dailyLogID string, input AddWorkoutExerciseInput) (*WorkoutExerciseRecord, error) {
	uid, dlid, err := parseTwoUUIDs(userID, dailyLogID)
	if err != nil {
		return nil, err
	}
	exid, err := uuidFromString(input.ExerciseID)
	if err != nil {
		return nil, err
	}

	_, err = q.TempShiftWorkoutExercisePositionsForInsert(ctx, generated.TempShiftWorkoutExercisePositionsForInsertParams{
		UserID:     uid,
		DailyLogID: dlid,
		Position:   input.Position,
	})
	if err != nil {
		return nil, err
	}
	row, err := q.CreateWorkoutExercise(ctx, generated.CreateWorkoutExerciseParams{
		UserID:                uid,
		DailyLogID:            dlid,
		ExerciseID:            exid,
		Position:              input.Position,
		WorkingWeightSnapshot: nullableFloat4(input.WorkingWeightSnapshot),
		Notes:                 nullableText(input.Notes),
	})
	if err != nil {
		return nil, err
	}
	_, err = q.NormalizeWorkoutExercisePositionsForInsert(ctx, generated.NormalizeWorkoutExercisePositionsForInsertParams{
		UserID:     uid,
		DailyLogID: dlid,
		Position:   input.Position,
	})
	if err != nil {
		return nil, err
	}
	return workoutExerciseRecordFromRow(row), nil
}

func updateWorkoutExerciseInTx(ctx context.Context, q *generated.Queries, userID string, dailyLogID string, workoutExerciseID string, input UpdateWorkoutExerciseInput) (*WorkoutExerciseRecord, bool, error) {
	uid, dlid, err := parseTwoUUIDs(userID, dailyLogID)
	if err != nil {
		return nil, false, err
	}
	weid, err := uuidFromString(workoutExerciseID)
	if err != nil {
		return nil, false, err
	}

	current, err := q.ListWorkoutExercisesByDailyLog(ctx, generated.ListWorkoutExercisesByDailyLogParams{
		UserID:     uid,
		DailyLogID: dlid,
	})
	if err != nil {
		return nil, false, err
	}
	currentRow := workoutExerciseRowByID(current, workoutExerciseID)
	if currentRow == nil {
		return nil, false, nil
	}

	var out *WorkoutExerciseRecord
	changed := false
	if input.Position != nil {
		orderedIDs, moved, ok := movedIDOrder(currentWorkoutExerciseIDs(current), workoutExerciseID, *input.Position)
		if !ok {
			return nil, false, nil
		}
		if moved {
			if err := reorderWorkoutExercisesInTx(ctx, q, userID, dailyLogID, orderedIDs); err != nil {
				return nil, false, err
			}
			changed = true
		}
	}

	if input.SetNotes {
		row, err := q.UpdateWorkoutExercise(ctx, generated.UpdateWorkoutExerciseParams{
			SetNotes: input.SetNotes,
			Notes:    nullableText(input.Notes),
			UserID:   uid,
			ID:       weid,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, false, nil
			}
			return nil, false, err
		}
		out = workoutExerciseRecordFromRow(row)
		changed = true
	}

	if out == nil {
		if changed {
			out, err = getWorkoutExerciseRecordByID(ctx, q, uid, dlid, workoutExerciseID)
			if err != nil {
				return nil, false, err
			}
		} else {
			out = workoutExerciseRecordFromRow(*currentRow)
		}
	}
	return out, changed, nil
}

func deleteWorkoutExerciseInTx(ctx context.Context, q *generated.Queries, userID string, workoutExerciseID string) (*WorkoutExerciseRecord, error) {
	uid, weid, err := parseTwoUUIDs(userID, workoutExerciseID)
	if err != nil {
		return nil, err
	}
	row, err := q.DeleteWorkoutExercise(ctx, generated.DeleteWorkoutExerciseParams{UserID: uid, ID: weid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return workoutExerciseRecordFromRow(row), nil
}

func reorderWorkoutExercisesInTx(ctx context.Context, q *generated.Queries, userID string, dailyLogID string, orderedIDs []string) error {
	uid, dlid, err := parseTwoUUIDs(userID, dailyLogID)
	if err != nil {
		return err
	}
	current, err := q.ListWorkoutExercisesByDailyLog(ctx, generated.ListWorkoutExercisesByDailyLogParams{
		UserID:     uid,
		DailyLogID: dlid,
	})
	if err != nil {
		return err
	}
	if err := requireSameIDs(currentWorkoutExerciseIDs(current), orderedIDs); err != nil {
		return err
	}
	for i, id := range orderedIDs {
		if _, err := q.SetWorkoutExercisePosition(ctx, generated.SetWorkoutExercisePositionParams{
			Position:   int32(1000000 + i + 1),
			UserID:     uid,
			DailyLogID: dlid,
			ID:         mustUUID(id),
		}); err != nil {
			return err
		}
	}
	for i, id := range orderedIDs {
		if _, err := q.SetWorkoutExercisePosition(ctx, generated.SetWorkoutExercisePositionParams{
			Position:   int32(i + 1),
			UserID:     uid,
			DailyLogID: dlid,
			ID:         mustUUID(id),
		}); err != nil {
			return err
		}
	}
	return nil
}

func addWorkoutSetInTx(ctx context.Context, q *generated.Queries, workoutExerciseID string, input AddWorkoutSetInput) (*WorkoutSetRecord, error) {
	weid, err := uuidFromString(workoutExerciseID)
	if err != nil {
		return nil, err
	}
	_, err = q.TempShiftWorkoutSetNumbersForInsert(ctx, generated.TempShiftWorkoutSetNumbersForInsertParams{
		WorkoutExerciseID: weid,
		SetNumber:         input.SetNumber,
	})
	if err != nil {
		return nil, err
	}
	row, err := q.CreateWorkoutSet(ctx, generated.CreateWorkoutSetParams{
		WorkoutExerciseID: weid,
		SetNumber:         input.SetNumber,
		Weight:            float32(input.Weight),
		Reps:              input.Reps,
		Rpe:               nullableFloat4(input.RPE),
		Rir:               nullableInt4(input.RIR),
		Notes:             nullableText(input.Notes),
	})
	if err != nil {
		return nil, err
	}
	_, err = q.NormalizeWorkoutSetNumbersForInsert(ctx, generated.NormalizeWorkoutSetNumbersForInsertParams{
		WorkoutExerciseID: weid,
		SetNumber:         input.SetNumber,
	})
	if err != nil {
		return nil, err
	}
	return workoutSetRecordFromRow(row), nil
}

func updateWorkoutSetInTx(ctx context.Context, q *generated.Queries, workoutExerciseID string, workoutSetID string, input UpdateWorkoutSetInput) (*WorkoutSetRecord, bool, error) {
	weid, wsid, err := parseTwoUUIDs(workoutExerciseID, workoutSetID)
	if err != nil {
		return nil, false, err
	}

	var current []generated.WorkoutSet
	var currentRow *generated.WorkoutSet
	if input.SetNumber != nil {
		current, err = q.ListWorkoutSetsByExerciseIDs(ctx, []pgtype.UUID{weid})
		if err != nil {
			return nil, false, err
		}
		currentRow = workoutSetRowByID(current, workoutSetID)
		if currentRow == nil {
			return nil, false, nil
		}
		orderedIDs, moved, ok := movedIDOrder(currentWorkoutSetIDs(current), workoutSetID, *input.SetNumber)
		if !ok {
			return nil, false, nil
		}
		if moved {
			if err := reorderWorkoutSetsInTx(ctx, q, workoutExerciseID, orderedIDs); err != nil {
				return nil, false, err
			}
		}
		if !moved && !hasWorkoutSetFieldUpdate(input) {
			return workoutSetRecordFromRow(*currentRow), false, nil
		}
	}

	changed := input.SetNumber != nil
	var out *WorkoutSetRecord
	if hasWorkoutSetFieldUpdate(input) {
		row, err := q.UpdateWorkoutSet(ctx, generated.UpdateWorkoutSetParams{
			Weight:            nullableFloat4(input.Weight),
			Reps:              nullableInt4(input.Reps),
			SetRpe:            input.SetRPE,
			Rpe:               nullableFloat4(input.RPE),
			SetRir:            input.SetRIR,
			Rir:               nullableInt4(input.RIR),
			SetNotes:          input.SetNotes,
			Notes:             nullableText(input.Notes),
			WorkoutExerciseID: weid,
			ID:                wsid,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, false, nil
			}
			return nil, false, err
		}
		out = workoutSetRecordFromRow(row)
		changed = true
	}

	if !changed {
		return nil, false, nil
	}
	if out != nil {
		return out, true, nil
	}
	out, err = getWorkoutSetRecordByID(ctx, q, weid, workoutSetID)
	if err != nil {
		return nil, false, err
	}
	return out, out != nil, nil
}

func deleteWorkoutSetInTx(ctx context.Context, q *generated.Queries, workoutExerciseID string, workoutSetID string) (*WorkoutSetRecord, error) {
	weid, wsid, err := parseTwoUUIDs(workoutExerciseID, workoutSetID)
	if err != nil {
		return nil, err
	}
	row, err := q.DeleteWorkoutSet(ctx, generated.DeleteWorkoutSetParams{WorkoutExerciseID: weid, ID: wsid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return workoutSetRecordFromRow(row), nil
}

func reorderWorkoutSetsInTx(ctx context.Context, q *generated.Queries, workoutExerciseID string, orderedIDs []string) error {
	weid, err := uuidFromString(workoutExerciseID)
	if err != nil {
		return err
	}
	current, err := q.ListWorkoutSetsByExerciseIDs(ctx, []pgtype.UUID{weid})
	if err != nil {
		return err
	}
	if err := requireSameIDs(currentWorkoutSetIDs(current), orderedIDs); err != nil {
		return err
	}
	for i, id := range orderedIDs {
		if _, err := q.SetWorkoutSetNumber(ctx, generated.SetWorkoutSetNumberParams{
			SetNumber:         int32(1000000 + i + 1),
			WorkoutExerciseID: weid,
			ID:                mustUUID(id),
		}); err != nil {
			return err
		}
	}
	for i, id := range orderedIDs {
		if _, err := q.SetWorkoutSetNumber(ctx, generated.SetWorkoutSetNumberParams{
			SetNumber:         int32(i + 1),
			WorkoutExerciseID: weid,
			ID:                mustUUID(id),
		}); err != nil {
			return err
		}
	}
	return nil
}

func incrementDailyLogVersion(ctx context.Context, q *generated.Queries, userID string, dailyLogID string, op string) (*DailyLogRecord, error) {
	row, err := incrementDailyLogVersionInTx(ctx, q, userID, dailyLogID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return dailyLogRecordFromRow(row), nil
}

func incrementDailyLogVersionInTx(ctx context.Context, q *generated.Queries, userID string, dailyLogID string) (generated.DailyLog, error) {
	uid, dlid, err := parseTwoUUIDs(userID, dailyLogID)
	if err != nil {
		return generated.DailyLog{}, err
	}
	return q.IncrementDailyLogVersion(ctx, generated.IncrementDailyLogVersionParams{UserID: uid, ID: dlid})
}

func currentWorkoutExerciseIDs(rows []generated.WorkoutExercise) []string {
	out := make([]string, len(rows))
	for i, row := range rows {
		out[i] = row.ID.String()
	}
	return out
}

func currentWorkoutSetIDs(rows []generated.WorkoutSet) []string {
	out := make([]string, len(rows))
	for i, row := range rows {
		out[i] = row.ID.String()
	}
	return out
}

func movedIDOrder(current []string, targetID string, position int32) ([]string, bool, bool) {
	if position < 1 || int(position) > len(current) {
		return nil, false, false
	}

	from := -1
	for i, id := range current {
		if id == targetID {
			from = i
			break
		}
	}
	if from == -1 {
		return nil, false, false
	}

	to := int(position) - 1
	ordered := append([]string(nil), current...)
	if from == to {
		return ordered, false, true
	}

	movedID := ordered[from]
	ordered = append(ordered[:from], ordered[from+1:]...)
	if to >= len(ordered) {
		ordered = append(ordered, movedID)
	} else {
		ordered = append(ordered[:to], append([]string{movedID}, ordered[to:]...)...)
	}
	return ordered, true, true
}

func workoutExerciseRowByID(rows []generated.WorkoutExercise, id string) *generated.WorkoutExercise {
	for i := range rows {
		if rows[i].ID.String() == id {
			return &rows[i]
		}
	}
	return nil
}

func workoutSetRowByID(rows []generated.WorkoutSet, id string) *generated.WorkoutSet {
	for i := range rows {
		if rows[i].ID.String() == id {
			return &rows[i]
		}
	}
	return nil
}

func getWorkoutExerciseRecordByID(ctx context.Context, q *generated.Queries, userID pgtype.UUID, dailyLogID pgtype.UUID, workoutExerciseID string) (*WorkoutExerciseRecord, error) {
	rows, err := q.ListWorkoutExercisesByDailyLog(ctx, generated.ListWorkoutExercisesByDailyLogParams{
		UserID:     userID,
		DailyLogID: dailyLogID,
	})
	if err != nil {
		return nil, err
	}
	row := workoutExerciseRowByID(rows, workoutExerciseID)
	if row == nil {
		return nil, nil
	}
	return workoutExerciseRecordFromRow(*row), nil
}

func getWorkoutSetRecordByID(ctx context.Context, q *generated.Queries, workoutExerciseID pgtype.UUID, workoutSetID string) (*WorkoutSetRecord, error) {
	rows, err := q.ListWorkoutSetsByExerciseIDs(ctx, []pgtype.UUID{workoutExerciseID})
	if err != nil {
		return nil, err
	}
	row := workoutSetRowByID(rows, workoutSetID)
	if row == nil {
		return nil, nil
	}
	return workoutSetRecordFromRow(*row), nil
}

func hasWorkoutSetFieldUpdate(input UpdateWorkoutSetInput) bool {
	return input.Weight != nil || input.Reps != nil || input.SetRPE || input.SetRIR || input.SetNotes
}

func requireSameIDs(current []string, ordered []string) error {
	if len(current) != len(ordered) {
		return fmt.Errorf("ordered ids length %d does not match current length %d", len(ordered), len(current))
	}
	seen := make(map[string]int, len(current))
	for _, id := range current {
		seen[id]++
	}
	for _, id := range ordered {
		if seen[id] == 0 {
			return fmt.Errorf("ordered id %s is not part of current aggregate", id)
		}
		seen[id]--
	}
	return nil
}

func dailyLogRecordFromRow(row generated.DailyLog) *DailyLogRecord {
	return &DailyLogRecord{
		ID:        row.ID.String(),
		UserID:    row.UserID.String(),
		Date:      dateFromRow(row.Date),
		Notes:     textPtr(row.Notes),
		Version:   row.Version,
		CreatedAt: formatTimestamp(row.CreatedAt),
		UpdatedAt: formatTimestamp(row.UpdatedAt),
	}
}

func workoutExerciseRecordFromRow(row generated.WorkoutExercise) *WorkoutExerciseRecord {
	return &WorkoutExerciseRecord{
		ID:                    row.ID.String(),
		UserID:                row.UserID.String(),
		DailyLogID:            row.DailyLogID.String(),
		ExerciseID:            row.ExerciseID.String(),
		Position:              row.Position,
		WorkingWeightSnapshot: float4Ptr(row.WorkingWeightSnapshot),
		Notes:                 textPtr(row.Notes),
		CreatedAt:             formatTimestamp(row.CreatedAt),
		UpdatedAt:             formatTimestamp(row.UpdatedAt),
	}
}

func workoutExerciseRecordsFromRows(rows []generated.WorkoutExercise) []WorkoutExerciseRecord {
	out := make([]WorkoutExerciseRecord, len(rows))
	for i, row := range rows {
		out[i] = *workoutExerciseRecordFromRow(row)
	}
	return out
}

func workoutSetRecordFromRow(row generated.WorkoutSet) *WorkoutSetRecord {
	return &WorkoutSetRecord{
		ID:                row.ID.String(),
		WorkoutExerciseID: row.WorkoutExerciseID.String(),
		SetNumber:         row.SetNumber,
		Weight:            float64(row.Weight),
		Reps:              row.Reps,
		RPE:               float4Ptr(row.Rpe),
		RIR:               int4Ptr(row.Rir),
		Notes:             textPtr(row.Notes),
		CreatedAt:         formatTimestamp(row.CreatedAt),
		UpdatedAt:         formatTimestamp(row.UpdatedAt),
	}
}

func dateParam(date models.Date) pgtype.Date {
	return pgtype.Date{Time: date.Time(), Valid: !date.Time().IsZero()}
}

func dateFromRow(date pgtype.Date) models.Date {
	if !date.Valid {
		return models.Date{}
	}
	return models.MustDate(date.Time.Format("2006-01-02"))
}

func parseTwoUUIDs(first string, second string) (pgtype.UUID, pgtype.UUID, error) {
	one, err := uuidFromString(first)
	if err != nil {
		return pgtype.UUID{}, pgtype.UUID{}, err
	}
	two, err := uuidFromString(second)
	if err != nil {
		return pgtype.UUID{}, pgtype.UUID{}, err
	}
	return one, two, nil
}

func mustUUID(value string) pgtype.UUID {
	uuid, err := uuidFromString(value)
	if err != nil {
		panic(err)
	}
	return uuid
}

func nullableInt4(value *int32) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *value, Valid: true}
}

func int4Ptr(value pgtype.Int4) *int32 {
	if !value.Valid {
		return nil
	}
	return &value.Int32
}
