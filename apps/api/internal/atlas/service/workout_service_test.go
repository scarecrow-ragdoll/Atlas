// FILE: apps/api/internal/atlas/service/workout_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for WorkoutService DailyLog validation, version checks, ordering, snapshots, and empty aggregate retention.
//   SCOPE: Service-level fake repository tests for notes, workout exercise, workout set, reorder, validation, not-found, and conflict behavior.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/models, apps/api/internal/atlas/repository/postgres.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestWorkoutService_* - Unit coverage for DailyLog notes, exercise/set mutations, validation, conflict, not-found, snapshot, reorder, and no-op behavior.
//   fakeWorkoutRepo - In-memory WorkoutRepository fake that tracks aggregate state and version changes.
//   fakeExerciseRepo - ExerciseRepository fake used to prove exercise lookup and snapshot behavior.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.4 - Documented bounded fake reindex conversions for lint-clean final gates.
// END_CHANGE_SUMMARY

package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/apps/api/internal/atlas/service"
)

const (
	testDailyLogID        = "daily-log-1"
	testWorkoutExerciseID = "workout-exercise-1"
	testWorkoutSetID      = "workout-set-1"
	testExerciseID        = "exercise-1"
)

var testWorkoutDate = models.MustDate("2026-06-19")

func ptrInt32(i int32) *int32 { return &i }

type fakeExerciseRepo struct {
	atlasPostgres.ExerciseRepository
	getByIDHits int
	getByIDFn   func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error)
}

func (f *fakeExerciseRepo) GetByID(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
	f.getByIDHits++
	if f.getByIDFn == nil {
		return nil, nil
	}
	return f.getByIDFn(ctx, userID, id)
}

type fakeWorkoutRepo struct {
	atlasPostgres.WorkoutRepository

	dailyLogs         map[string]*atlasPostgres.DailyLogAggregate
	dateIndex         map[string]string
	exerciseToLog     map[string]string
	setToExercise     map[string]string
	summaryRecords    []atlasPostgres.DailyLogSummaryRecord
	lastSummaryUserID string
	lastSummaryFrom   models.Date
	lastSummaryTo     models.Date
	nextDailyLog      int
	nextExercise      int
	nextSet           int
	listSummaryHits   int
	getOrCreateHits   int
	lockDateHit       int
	lockExerciseHit   int
	lockSetHit        int
	mutationHits      int
}

func newFakeWorkoutRepo() *fakeWorkoutRepo {
	return &fakeWorkoutRepo{
		dailyLogs:     map[string]*atlasPostgres.DailyLogAggregate{},
		dateIndex:     map[string]string{},
		exerciseToLog: map[string]string{},
		setToExercise: map[string]string{},
		nextDailyLog:  1,
		nextExercise:  1,
		nextSet:       1,
	}
}

func (f *fakeWorkoutRepo) seedDailyLog(date models.Date, version int32) *atlasPostgres.DailyLogAggregate {
	id := fmt.Sprintf("daily-log-%d", f.nextDailyLog)
	f.nextDailyLog++
	aggregate := &atlasPostgres.DailyLogAggregate{
		DailyLog: atlasPostgres.DailyLogRecord{
			ID:        id,
			UserID:    testUserID,
			Date:      date,
			Version:   version,
			CreatedAt: "2026-06-19T00:00:00Z",
			UpdatedAt: "2026-06-19T00:00:00Z",
		},
		WorkoutExercises: []atlasPostgres.WorkoutExerciseRecord{},
	}
	f.dailyLogs[id] = aggregate
	f.dateIndex[dateKey(testUserID, date)] = id
	return aggregate
}

func (f *fakeWorkoutRepo) seedWorkoutExercise(dailyLogID string, id string, exerciseID string, position int32, snapshot *float64) *atlasPostgres.WorkoutExerciseRecord {
	aggregate := f.dailyLogs[dailyLogID]
	record := atlasPostgres.WorkoutExerciseRecord{
		ID:                    id,
		UserID:                testUserID,
		DailyLogID:            dailyLogID,
		ExerciseID:            exerciseID,
		Position:              position,
		WorkingWeightSnapshot: snapshot,
		Sets:                  []atlasPostgres.WorkoutSetRecord{},
		CreatedAt:             "2026-06-19T00:00:00Z",
		UpdatedAt:             "2026-06-19T00:00:00Z",
	}
	aggregate.WorkoutExercises = append(aggregate.WorkoutExercises, record)
	f.exerciseToLog[id] = dailyLogID
	return &aggregate.WorkoutExercises[len(aggregate.WorkoutExercises)-1]
}

func (f *fakeWorkoutRepo) seedWorkoutSet(workoutExerciseID string, id string, setNumber int32) {
	dailyLogID := f.exerciseToLog[workoutExerciseID]
	exercise := findRepoWorkoutExercise(f.dailyLogs[dailyLogID], workoutExerciseID)
	exercise.Sets = append(exercise.Sets, atlasPostgres.WorkoutSetRecord{
		ID:                id,
		WorkoutExerciseID: workoutExerciseID,
		SetNumber:         setNumber,
		Weight:            100,
		Reps:              5,
		CreatedAt:         "2026-06-19T00:00:00Z",
		UpdatedAt:         "2026-06-19T00:00:00Z",
	})
	f.setToExercise[id] = workoutExerciseID
}

func (f *fakeWorkoutRepo) GetDailyLogByDate(ctx context.Context, userID string, date models.Date) (*atlasPostgres.DailyLogRecord, error) {
	id := f.dateIndex[dateKey(userID, date)]
	if id == "" {
		return nil, nil
	}
	record := f.dailyLogs[id].DailyLog
	return &record, nil
}

func (f *fakeWorkoutRepo) GetOrCreateDailyLogByDate(ctx context.Context, userID string, date models.Date) (*atlasPostgres.DailyLogRecord, error) {
	f.getOrCreateHits++
	if record, err := f.GetDailyLogByDate(ctx, userID, date); record != nil || err != nil {
		return record, err
	}
	aggregate := f.seedDailyLog(date, 0)
	record := aggregate.DailyLog
	return &record, nil
}

func (f *fakeWorkoutRepo) GetDailyLogAggregate(ctx context.Context, userID string, dailyLogID string) (*atlasPostgres.DailyLogAggregate, error) {
	return cloneAggregate(f.dailyLogs[dailyLogID]), nil
}

func (f *fakeWorkoutRepo) ListDailyLogSummaries(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]atlasPostgres.DailyLogSummaryRecord, error) {
	f.listSummaryHits++
	f.lastSummaryUserID = userID
	f.lastSummaryFrom = fromDate
	f.lastSummaryTo = toDate
	return append([]atlasPostgres.DailyLogSummaryRecord(nil), f.summaryRecords...), nil
}

func (f *fakeWorkoutRepo) WithLockedDailyLogByDate(ctx context.Context, userID string, date models.Date, fn atlasPostgres.LockedDailyLogFunc) error {
	f.lockDateHit++
	id := f.dateIndex[dateKey(userID, date)]
	if id == "" {
		return pgx.ErrNoRows
	}
	record := f.dailyLogs[id].DailyLog
	return fn(ctx, &fakeWorkoutTx{repo: f}, &record)
}

func (f *fakeWorkoutRepo) WithLockedDailyLogByWorkoutExerciseID(ctx context.Context, userID string, workoutExerciseID string, fn atlasPostgres.LockedDailyLogFunc) error {
	f.lockExerciseHit++
	dailyLogID := f.exerciseToLog[workoutExerciseID]
	if dailyLogID == "" {
		return pgx.ErrNoRows
	}
	record := f.dailyLogs[dailyLogID].DailyLog
	return fn(ctx, &fakeWorkoutTx{repo: f}, &record)
}

func (f *fakeWorkoutRepo) WithLockedDailyLogByWorkoutSetID(ctx context.Context, userID string, workoutSetID string, fn atlasPostgres.LockedDailyLogFunc) error {
	f.lockSetHit++
	workoutExerciseID := f.setToExercise[workoutSetID]
	dailyLogID := f.exerciseToLog[workoutExerciseID]
	if dailyLogID == "" {
		return pgx.ErrNoRows
	}
	record := f.dailyLogs[dailyLogID].DailyLog
	return fn(ctx, &fakeWorkoutTx{repo: f}, &record)
}

func (f *fakeWorkoutRepo) IncrementDailyLogVersion(ctx context.Context, userID string, dailyLogID string) (*atlasPostgres.DailyLogRecord, error) {
	return (&fakeWorkoutTx{repo: f}).IncrementDailyLogVersion(ctx, userID, dailyLogID)
}

func (f *fakeWorkoutRepo) UpdateDailyLogNotes(ctx context.Context, userID string, dailyLogID string, notes *string) (*atlasPostgres.DailyLogRecord, error) {
	return (&fakeWorkoutTx{repo: f}).UpdateDailyLogNotes(ctx, userID, dailyLogID, notes)
}

func (f *fakeWorkoutRepo) AddWorkoutExercise(ctx context.Context, userID string, dailyLogID string, input atlasPostgres.AddWorkoutExerciseInput) (*atlasPostgres.WorkoutExerciseRecord, error) {
	return (&fakeWorkoutTx{repo: f}).AddWorkoutExercise(ctx, userID, dailyLogID, input)
}

func (f *fakeWorkoutRepo) UpdateWorkoutExercise(ctx context.Context, userID string, workoutExerciseID string, input atlasPostgres.UpdateWorkoutExerciseInput) (*atlasPostgres.WorkoutExerciseRecord, error) {
	return (&fakeWorkoutTx{repo: f}).UpdateWorkoutExercise(ctx, userID, workoutExerciseID, input)
}

func (f *fakeWorkoutRepo) DeleteWorkoutExercise(ctx context.Context, userID string, workoutExerciseID string) (*atlasPostgres.WorkoutExerciseRecord, error) {
	return (&fakeWorkoutTx{repo: f}).DeleteWorkoutExercise(ctx, userID, workoutExerciseID)
}

func (f *fakeWorkoutRepo) ReorderWorkoutExercises(ctx context.Context, userID string, dailyLogID string, orderedIDs []string) error {
	return (&fakeWorkoutTx{repo: f}).ReorderWorkoutExercises(ctx, userID, dailyLogID, orderedIDs)
}

func (f *fakeWorkoutRepo) AddWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, input atlasPostgres.AddWorkoutSetInput) (*atlasPostgres.WorkoutSetRecord, error) {
	return (&fakeWorkoutTx{repo: f}).AddWorkoutSet(ctx, userID, workoutExerciseID, input)
}

func (f *fakeWorkoutRepo) UpdateWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, workoutSetID string, input atlasPostgres.UpdateWorkoutSetInput) (*atlasPostgres.WorkoutSetRecord, error) {
	return (&fakeWorkoutTx{repo: f}).UpdateWorkoutSet(ctx, userID, workoutExerciseID, workoutSetID, input)
}

func (f *fakeWorkoutRepo) DeleteWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, workoutSetID string) (*atlasPostgres.WorkoutSetRecord, error) {
	return (&fakeWorkoutTx{repo: f}).DeleteWorkoutSet(ctx, userID, workoutExerciseID, workoutSetID)
}

func (f *fakeWorkoutRepo) ReorderWorkoutSets(ctx context.Context, userID string, workoutExerciseID string, orderedIDs []string) error {
	return (&fakeWorkoutTx{repo: f}).ReorderWorkoutSets(ctx, userID, workoutExerciseID, orderedIDs)
}

type fakeWorkoutTx struct {
	repo *fakeWorkoutRepo
}

func (tx *fakeWorkoutTx) GetDailyLogAggregate(ctx context.Context, userID string, dailyLogID string) (*atlasPostgres.DailyLogAggregate, error) {
	return cloneAggregate(tx.repo.dailyLogs[dailyLogID]), nil
}

func (tx *fakeWorkoutTx) IncrementDailyLogVersion(ctx context.Context, userID string, dailyLogID string) (*atlasPostgres.DailyLogRecord, error) {
	tx.repo.mutationHits++
	tx.repo.dailyLogs[dailyLogID].DailyLog.Version++
	record := tx.repo.dailyLogs[dailyLogID].DailyLog
	return &record, nil
}

func (tx *fakeWorkoutTx) UpdateDailyLogNotes(ctx context.Context, userID string, dailyLogID string, notes *string) (*atlasPostgres.DailyLogRecord, error) {
	tx.repo.mutationHits++
	tx.repo.dailyLogs[dailyLogID].DailyLog.Notes = notes
	record := tx.repo.dailyLogs[dailyLogID].DailyLog
	return &record, nil
}

func (tx *fakeWorkoutTx) AddWorkoutExercise(ctx context.Context, userID string, dailyLogID string, input atlasPostgres.AddWorkoutExerciseInput) (*atlasPostgres.WorkoutExerciseRecord, error) {
	tx.repo.mutationHits++
	aggregate := tx.repo.dailyLogs[dailyLogID]
	id := fmt.Sprintf("workout-exercise-%d", tx.repo.nextExercise)
	tx.repo.nextExercise++
	record := atlasPostgres.WorkoutExerciseRecord{
		ID:                    id,
		UserID:                userID,
		DailyLogID:            dailyLogID,
		ExerciseID:            input.ExerciseID,
		Position:              input.Position,
		WorkingWeightSnapshot: input.WorkingWeightSnapshot,
		Notes:                 input.Notes,
		Sets:                  []atlasPostgres.WorkoutSetRecord{},
		CreatedAt:             "2026-06-19T00:00:00Z",
		UpdatedAt:             "2026-06-19T00:00:00Z",
	}
	insertWorkoutExercise(aggregate, record)
	tx.repo.exerciseToLog[id] = dailyLogID
	out := record
	return &out, nil
}

func (tx *fakeWorkoutTx) UpdateWorkoutExercise(ctx context.Context, userID string, workoutExerciseID string, input atlasPostgres.UpdateWorkoutExerciseInput) (*atlasPostgres.WorkoutExerciseRecord, error) {
	tx.repo.mutationHits++
	dailyLogID := tx.repo.exerciseToLog[workoutExerciseID]
	aggregate := tx.repo.dailyLogs[dailyLogID]
	if input.Position != nil {
		ordered := movedOrder(repoWorkoutExerciseIDs(aggregate.WorkoutExercises), workoutExerciseID, *input.Position)
		reorderWorkoutExercises(aggregate, ordered)
	}
	exercise := findRepoWorkoutExercise(aggregate, workoutExerciseID)
	if exercise == nil {
		return nil, nil
	}
	if input.SetNotes {
		exercise.Notes = input.Notes
	}
	out := *exercise
	return &out, nil
}

func (tx *fakeWorkoutTx) DeleteWorkoutExercise(ctx context.Context, userID string, workoutExerciseID string) (*atlasPostgres.WorkoutExerciseRecord, error) {
	tx.repo.mutationHits++
	dailyLogID := tx.repo.exerciseToLog[workoutExerciseID]
	aggregate := tx.repo.dailyLogs[dailyLogID]
	for i := range aggregate.WorkoutExercises {
		if aggregate.WorkoutExercises[i].ID == workoutExerciseID {
			out := aggregate.WorkoutExercises[i]
			for _, set := range out.Sets {
				delete(tx.repo.setToExercise, set.ID)
			}
			aggregate.WorkoutExercises = append(aggregate.WorkoutExercises[:i], aggregate.WorkoutExercises[i+1:]...)
			delete(tx.repo.exerciseToLog, workoutExerciseID)
			return &out, nil
		}
	}
	return nil, nil
}

func (tx *fakeWorkoutTx) ReorderWorkoutExercises(ctx context.Context, userID string, dailyLogID string, orderedIDs []string) error {
	tx.repo.mutationHits++
	reorderWorkoutExercises(tx.repo.dailyLogs[dailyLogID], orderedIDs)
	return nil
}

func (tx *fakeWorkoutTx) AddWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, input atlasPostgres.AddWorkoutSetInput) (*atlasPostgres.WorkoutSetRecord, error) {
	tx.repo.mutationHits++
	dailyLogID := tx.repo.exerciseToLog[workoutExerciseID]
	exercise := findRepoWorkoutExercise(tx.repo.dailyLogs[dailyLogID], workoutExerciseID)
	id := fmt.Sprintf("workout-set-%d", tx.repo.nextSet)
	tx.repo.nextSet++
	record := atlasPostgres.WorkoutSetRecord{
		ID:                id,
		WorkoutExerciseID: workoutExerciseID,
		SetNumber:         input.SetNumber,
		Weight:            input.Weight,
		Reps:              input.Reps,
		RPE:               input.RPE,
		RIR:               input.RIR,
		Notes:             input.Notes,
		CreatedAt:         "2026-06-19T00:00:00Z",
		UpdatedAt:         "2026-06-19T00:00:00Z",
	}
	insertWorkoutSet(exercise, record)
	tx.repo.setToExercise[id] = workoutExerciseID
	out := record
	return &out, nil
}

func (tx *fakeWorkoutTx) UpdateWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, workoutSetID string, input atlasPostgres.UpdateWorkoutSetInput) (*atlasPostgres.WorkoutSetRecord, error) {
	tx.repo.mutationHits++
	dailyLogID := tx.repo.exerciseToLog[workoutExerciseID]
	exercise := findRepoWorkoutExercise(tx.repo.dailyLogs[dailyLogID], workoutExerciseID)
	set := findRepoWorkoutSet(exercise, workoutSetID)
	if set == nil {
		return nil, nil
	}
	if input.Weight != nil {
		set.Weight = *input.Weight
	}
	if input.Reps != nil {
		set.Reps = *input.Reps
	}
	if input.SetRPE {
		set.RPE = input.RPE
	}
	if input.SetRIR {
		set.RIR = input.RIR
	}
	if input.SetNotes {
		set.Notes = input.Notes
	}
	if input.SetNumber != nil {
		ordered := movedOrder(repoWorkoutSetIDs(exercise.Sets), workoutSetID, *input.SetNumber)
		reorderWorkoutSets(exercise, ordered)
	}
	out := *findRepoWorkoutSet(exercise, workoutSetID)
	return &out, nil
}

func (tx *fakeWorkoutTx) DeleteWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, workoutSetID string) (*atlasPostgres.WorkoutSetRecord, error) {
	tx.repo.mutationHits++
	dailyLogID := tx.repo.exerciseToLog[workoutExerciseID]
	exercise := findRepoWorkoutExercise(tx.repo.dailyLogs[dailyLogID], workoutExerciseID)
	for i := range exercise.Sets {
		if exercise.Sets[i].ID == workoutSetID {
			out := exercise.Sets[i]
			exercise.Sets = append(exercise.Sets[:i], exercise.Sets[i+1:]...)
			delete(tx.repo.setToExercise, workoutSetID)
			return &out, nil
		}
	}
	return nil, nil
}

func (tx *fakeWorkoutTx) ReorderWorkoutSets(ctx context.Context, userID string, workoutExerciseID string, orderedIDs []string) error {
	tx.repo.mutationHits++
	dailyLogID := tx.repo.exerciseToLog[workoutExerciseID]
	reorderWorkoutSets(findRepoWorkoutExercise(tx.repo.dailyLogs[dailyLogID], workoutExerciseID), orderedIDs)
	return nil
}

func TestWorkoutService_GetDailyLog_AbsentDateDoesNotCreateDailyLog(t *testing.T) {
	repo := newFakeWorkoutRepo()
	svc := service.NewWorkoutService(repo, &fakeExerciseRepo{})

	log, err := svc.GetDailyLog(ctx, testUserID, testWorkoutDate)

	require.NoError(t, err)
	assert.Nil(t, log)
	assert.Equal(t, 0, repo.getOrCreateHits)
	assert.Empty(t, repo.dateIndex)
	assert.Empty(t, repo.dailyLogs)
}

func TestWorkoutService_ListDailyLogSummaries_MapsRepositoryRecords(t *testing.T) {
	repo := newFakeWorkoutRepo()
	fromDate := models.MustDate("2026-06-18")
	toDate := models.MustDate("2026-06-19")
	repo.summaryRecords = []atlasPostgres.DailyLogSummaryRecord{
		{
			ID:                   "daily-log-1",
			Date:                 fromDate,
			Version:              2,
			WorkoutExerciseCount: 1,
			WorkoutSetCount:      3,
			TotalVolume:          4200,
			UpdatedAt:            "2026-06-18T12:00:00Z",
		},
		{
			ID:                   "daily-log-2",
			Date:                 toDate,
			Version:              5,
			WorkoutExerciseCount: 2,
			WorkoutSetCount:      6,
			TotalVolume:          8100,
			UpdatedAt:            "2026-06-19T12:00:00Z",
		},
	}
	svc := service.NewWorkoutService(repo, &fakeExerciseRepo{})

	summaries, err := svc.ListDailyLogSummaries(ctx, testUserID, fromDate, toDate)

	require.NoError(t, err)
	require.Len(t, summaries, 2)
	assert.Equal(t, "daily-log-1", summaries[0].ID)
	assert.Equal(t, fromDate, summaries[0].Date)
	assert.Equal(t, int32(2), summaries[0].Version)
	assert.Equal(t, int32(1), summaries[0].WorkoutExerciseCount)
	assert.Equal(t, int32(3), summaries[0].WorkoutSetCount)
	assert.Equal(t, 4200.0, summaries[0].TotalVolume)
	assert.Equal(t, "2026-06-18T12:00:00Z", summaries[0].UpdatedAt)
	assert.Equal(t, "daily-log-2", summaries[1].ID)
	assert.Equal(t, int32(5), summaries[1].Version)
	assert.Equal(t, 1, repo.listSummaryHits)
	assert.Equal(t, testUserID, repo.lastSummaryUserID)
	assert.Equal(t, fromDate, repo.lastSummaryFrom)
	assert.Equal(t, toDate, repo.lastSummaryTo)
	assert.Empty(t, repo.dailyLogs)
	assert.Empty(t, repo.dateIndex)
}

func TestWorkoutService_ListDailyLogSummaries_RejectsInvalidRangeWithoutRepoInteraction(t *testing.T) {
	repo := newFakeWorkoutRepo()
	repo.seedDailyLog(testWorkoutDate, 1)
	svc := service.NewWorkoutService(repo, &fakeExerciseRepo{})

	summaries, err := svc.ListDailyLogSummaries(ctx, testUserID, models.MustDate("2026-06-20"), models.MustDate("2026-06-19"))

	require.Nil(t, summaries)
	var validationErr *models.DailyLogValidationErr
	require.ErrorAs(t, err, &validationErr)
	assert.Equal(t, 0, repo.listSummaryHits)
	assert.Equal(t, int32(1), repo.dailyLogs[testDailyLogID].DailyLog.Version)
	assert.Len(t, repo.dailyLogs, 1)
}

func TestWorkoutService_UpdateNotes_CreatesDailyLogAtExpectedVersionZero(t *testing.T) {
	repo := newFakeWorkoutRepo()
	svc := service.NewWorkoutService(repo, &fakeExerciseRepo{})

	log, err := svc.UpdateDailyLogNotes(ctx, testUserID, testWorkoutDate, 0, ptrStr("felt strong"))

	require.NoError(t, err)
	require.NotNil(t, log)
	assert.Equal(t, int32(1), log.Version)
	require.NotNil(t, log.Notes)
	assert.Equal(t, "felt strong", *log.Notes)
	assert.NotNil(t, log.WorkoutExercises)
	assert.Empty(t, log.WorkoutExercises)
	assert.Equal(t, 1, repo.getOrCreateHits)
}

func TestWorkoutService_UpdateNotes_ExistingDailyLogUpdatesAndClearsNotes(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 4)
	aggregate.DailyLog.Notes = ptrStr("old notes")
	svc := service.NewWorkoutService(repo, &fakeExerciseRepo{})

	updated, err := svc.UpdateDailyLogNotes(ctx, testUserID, testWorkoutDate, 4, ptrStr("new notes"))

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, int32(5), updated.Version)
	require.NotNil(t, updated.Notes)
	assert.Equal(t, "new notes", *updated.Notes)
	require.NotNil(t, repo.dailyLogs[testDailyLogID].DailyLog.Notes)
	assert.Equal(t, "new notes", *repo.dailyLogs[testDailyLogID].DailyLog.Notes)

	cleared, err := svc.UpdateDailyLogNotes(ctx, testUserID, testWorkoutDate, 5, nil)

	require.NoError(t, err)
	require.NotNil(t, cleared)
	assert.Equal(t, int32(6), cleared.Version)
	assert.Nil(t, cleared.Notes)
	assert.Nil(t, repo.dailyLogs[testDailyLogID].DailyLog.Notes)
	assert.Equal(t, 0, repo.getOrCreateHits)
}

func TestWorkoutService_UpdateNotes_AbsentDateWithNonZeroExpectedVersionDoesNotCreateDailyLog(t *testing.T) {
	repo := newFakeWorkoutRepo()
	svc := service.NewWorkoutService(repo, &fakeExerciseRepo{})

	log, err := svc.UpdateDailyLogNotes(ctx, testUserID, testWorkoutDate, 2, ptrStr("stale notes"))

	require.Nil(t, log)
	var conflictErr *models.DailyLogConflictErr
	require.ErrorAs(t, err, &conflictErr)
	assert.Equal(t, int32(0), conflictErr.CurrentVersion)
	assert.Nil(t, conflictErr.CurrentDailyLog)
	assert.Equal(t, 0, repo.getOrCreateHits)
	assert.Empty(t, repo.dateIndex)
	assert.Empty(t, repo.dailyLogs)
}

func TestWorkoutService_UpdateNotes_RejectsStaleVersion(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 3)
	aggregate.DailyLog.Notes = ptrStr("current notes")
	svc := service.NewWorkoutService(repo, &fakeExerciseRepo{})

	log, err := svc.UpdateDailyLogNotes(ctx, testUserID, testWorkoutDate, 2, ptrStr("overwritten"))

	require.Nil(t, log)
	var conflictErr *models.DailyLogConflictErr
	require.ErrorAs(t, err, &conflictErr)
	assert.Equal(t, int32(3), conflictErr.CurrentVersion)
	require.NotNil(t, conflictErr.CurrentDailyLog)
	assert.Equal(t, int32(3), conflictErr.CurrentDailyLog.Version)
	require.NotNil(t, conflictErr.CurrentDailyLog.Notes)
	assert.Equal(t, "current notes", *conflictErr.CurrentDailyLog.Notes)
	assert.Equal(t, int32(3), repo.dailyLogs[testDailyLogID].DailyLog.Version)
	require.NotNil(t, repo.dailyLogs[testDailyLogID].DailyLog.Notes)
	assert.Equal(t, "current notes", *repo.dailyLogs[testDailyLogID].DailyLog.Notes)
}

func TestWorkoutService_RejectsNegativeExpectedVersionBeforeRepositoryMutation(t *testing.T) {
	cases := []struct {
		name string
		call func(svc service.WorkoutService) (*models.DailyLog, error)
	}{
		{
			name: "update notes",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.UpdateDailyLogNotes(ctx, testUserID, testWorkoutDate, -1, ptrStr("blocked"))
			},
		},
		{
			name: "add exercise",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.AddWorkoutExercise(ctx, testUserID, testWorkoutDate, -1, models.AddWorkoutExerciseInput{
					ExerciseID: testExerciseID,
				})
			},
		},
		{
			name: "update exercise",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.UpdateWorkoutExercise(ctx, testUserID, testWorkoutExerciseID, -1, models.UpdateWorkoutExerciseInput{
					Notes: ptrStr("blocked"),
				})
			},
		},
		{
			name: "add set",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.AddWorkoutSet(ctx, testUserID, testWorkoutExerciseID, -1, models.AddWorkoutSetInput{
					Weight: 100,
					Reps:   5,
				})
			},
		},
		{
			name: "update set",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.UpdateWorkoutSet(ctx, testUserID, testWorkoutSetID, -1, models.UpdateWorkoutSetInput{
					Weight: ptrFloat64(105),
				})
			},
		},
		{
			name: "reorder exercises",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.ReorderWorkoutExercises(ctx, testUserID, testWorkoutDate, -1, []string{testWorkoutExerciseID})
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := newFakeWorkoutRepo()
			aggregate := repo.seedDailyLog(testWorkoutDate, 2)
			repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, ptrFloat64(82.5))
			repo.seedWorkoutSet(testWorkoutExerciseID, testWorkoutSetID, 1)
			exercises := exerciseRepoWithRecord(ptrFloat64(90))
			svc := service.NewWorkoutService(repo, exercises)

			log, err := tc.call(svc)

			require.Nil(t, log)
			var validationErr *models.DailyLogValidationErr
			require.ErrorAs(t, err, &validationErr)
			assert.Equal(t, int32(2), repo.dailyLogs[testDailyLogID].DailyLog.Version)
			assert.Equal(t, 0, repo.getOrCreateHits)
			assert.Equal(t, 0, repo.lockDateHit)
			assert.Equal(t, 0, repo.lockExerciseHit)
			assert.Equal(t, 0, repo.lockSetHit)
			assert.Equal(t, 0, repo.mutationHits)
			assert.Equal(t, 0, exercises.getByIDHits)
			assert.Equal(t, []string{testWorkoutExerciseID}, repoWorkoutExerciseIDs(repo.dailyLogs[testDailyLogID].WorkoutExercises))
			assert.Equal(t, []string{testWorkoutSetID}, repoWorkoutSetIDs(repo.dailyLogs[testDailyLogID].WorkoutExercises[0].Sets))
		})
	}
}

func TestWorkoutService_AddExercise_RequiresExistingExercise(t *testing.T) {
	repo := newFakeWorkoutRepo()
	svc := service.NewWorkoutService(repo, &fakeExerciseRepo{})

	log, err := svc.AddWorkoutExercise(ctx, testUserID, testWorkoutDate, 0, models.AddWorkoutExerciseInput{
		ExerciseID: testExerciseID,
	})

	require.Nil(t, log)
	var notFoundErr *models.DailyLogNotFoundErr
	require.ErrorAs(t, err, &notFoundErr)
	assert.Zero(t, repo.getOrCreateHits)
	assert.Empty(t, repo.dailyLogs)
}

func TestWorkoutService_AddExercise_AbsentDateWithNonZeroExpectedVersionDoesNotCreateDailyLog(t *testing.T) {
	repo := newFakeWorkoutRepo()
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(ptrFloat64(82.5)))

	log, err := svc.AddWorkoutExercise(ctx, testUserID, testWorkoutDate, 2, models.AddWorkoutExerciseInput{
		ExerciseID: testExerciseID,
	})

	require.Nil(t, log)
	var conflictErr *models.DailyLogConflictErr
	require.ErrorAs(t, err, &conflictErr)
	assert.Equal(t, int32(0), conflictErr.CurrentVersion)
	assert.Nil(t, conflictErr.CurrentDailyLog)
	assert.Equal(t, 0, repo.getOrCreateHits)
	assert.Empty(t, repo.dateIndex)
	assert.Empty(t, repo.dailyLogs)
	assert.Empty(t, repo.exerciseToLog)
}

func TestWorkoutService_AddExercise_AbsentDateInvalidAppendPositionDoesNotCreateDailyLog(t *testing.T) {
	repo := newFakeWorkoutRepo()
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(ptrFloat64(82.5)))

	log, err := svc.AddWorkoutExercise(ctx, testUserID, testWorkoutDate, 0, models.AddWorkoutExerciseInput{
		ExerciseID: testExerciseID,
		Position:   ptrInt32(2),
	})

	require.Nil(t, log)
	var validationErr *models.DailyLogValidationErr
	require.ErrorAs(t, err, &validationErr)
	assert.Equal(t, 0, repo.getOrCreateHits)
	assert.Empty(t, repo.dateIndex)
	assert.Empty(t, repo.dailyLogs)
	assert.Empty(t, repo.exerciseToLog)
}

func TestWorkoutService_AddExercise_CapturesWorkingWeightSnapshot(t *testing.T) {
	repo := newFakeWorkoutRepo()
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(ptrFloat64(82.5)))

	log, err := svc.AddWorkoutExercise(ctx, testUserID, testWorkoutDate, 0, models.AddWorkoutExerciseInput{
		ExerciseID: testExerciseID,
	})

	require.NoError(t, err)
	require.NotNil(t, log)
	assert.Equal(t, int32(1), log.Version)
	require.Len(t, log.WorkoutExercises, 1)
	assert.Equal(t, testExerciseID, log.WorkoutExercises[0].ExerciseID)
	require.NotNil(t, log.WorkoutExercises[0].WorkingWeightSnapshot)
	assert.Equal(t, 82.5, *log.WorkoutExercises[0].WorkingWeightSnapshot)
	require.NotNil(t, log.WorkoutExercises[0].Exercise)
	require.NotNil(t, log.WorkoutExercises[0].Exercise.WorkingWeight)
	assert.Equal(t, 82.5, *log.WorkoutExercises[0].Exercise.WorkingWeight)
}

func TestWorkoutService_AddExercise_AllowsDuplicateExerciseID(t *testing.T) {
	repo := newFakeWorkoutRepo()
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(ptrFloat64(60)))

	first, err := svc.AddWorkoutExercise(ctx, testUserID, testWorkoutDate, 0, models.AddWorkoutExerciseInput{
		ExerciseID: testExerciseID,
	})
	require.NoError(t, err)
	require.Len(t, first.WorkoutExercises, 1)

	second, err := svc.AddWorkoutExercise(ctx, testUserID, testWorkoutDate, 1, models.AddWorkoutExerciseInput{
		ExerciseID: testExerciseID,
	})

	require.NoError(t, err)
	require.NotNil(t, second)
	assert.Equal(t, int32(2), second.Version)
	require.Len(t, second.WorkoutExercises, 2)
	assert.Equal(t, testExerciseID, second.WorkoutExercises[0].ExerciseID)
	assert.Equal(t, testExerciseID, second.WorkoutExercises[1].ExerciseID)
	assert.Equal(t, int32(1), second.WorkoutExercises[0].Position)
	assert.Equal(t, int32(2), second.WorkoutExercises[1].Position)
}

func TestWorkoutService_AddExercise_InsertsAtPositionAndReindexes(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 2)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, "workout-exercise-1", "exercise-1", 1, nil)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, "workout-exercise-2", "exercise-2", 2, nil)
	repo.nextExercise = 3
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(ptrFloat64(77.5)))

	log, err := svc.AddWorkoutExercise(ctx, testUserID, testWorkoutDate, 2, models.AddWorkoutExerciseInput{
		ExerciseID: "exercise-new",
		Position:   ptrInt32(2),
	})

	require.NoError(t, err)
	require.NotNil(t, log)
	assert.Equal(t, int32(3), log.Version)
	require.Len(t, log.WorkoutExercises, 3)
	assert.Equal(t, []string{"workout-exercise-1", "workout-exercise-3", "workout-exercise-2"}, modelWorkoutExerciseIDs(log.WorkoutExercises))
	assert.Equal(t, []int32{1, 2, 3}, []int32{
		log.WorkoutExercises[0].Position,
		log.WorkoutExercises[1].Position,
		log.WorkoutExercises[2].Position,
	})
	assert.Equal(t, "exercise-new", log.WorkoutExercises[1].ExerciseID)
}

func TestWorkoutService_UpdateExercise_RejectsEmptyInputWithoutVersionChange(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 5)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

	log, err := svc.UpdateWorkoutExercise(ctx, testUserID, testWorkoutExerciseID, 5, models.UpdateWorkoutExerciseInput{})

	require.Nil(t, log)
	var validationErr *models.DailyLogValidationErr
	require.ErrorAs(t, err, &validationErr)
	assert.Equal(t, int32(5), repo.dailyLogs[testDailyLogID].DailyLog.Version)
	assert.Equal(t, []string{testWorkoutExerciseID}, repoWorkoutExerciseIDs(repo.dailyLogs[testDailyLogID].WorkoutExercises))
}

func TestWorkoutService_UpdateExercise_RejectsSamePositionOnlyWithoutVersionChange(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 5)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

	log, err := svc.UpdateWorkoutExercise(ctx, testUserID, testWorkoutExerciseID, 5, models.UpdateWorkoutExerciseInput{
		Position: ptrInt32(1),
	})

	require.Nil(t, log)
	var validationErr *models.DailyLogValidationErr
	require.ErrorAs(t, err, &validationErr)
	assert.Equal(t, int32(5), repo.dailyLogs[testDailyLogID].DailyLog.Version)
	require.Len(t, repo.dailyLogs[testDailyLogID].WorkoutExercises, 1)
	assert.Equal(t, int32(1), repo.dailyLogs[testDailyLogID].WorkoutExercises[0].Position)
}

func TestWorkoutService_UpdateExercise_UpdatesAndClearsNotesWithoutChangingSnapshot(t *testing.T) {
	workingWeight := 82.5
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 0)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, ptrFloat64(workingWeight))
	exercises := exerciseRepoWithMutableWorkingWeight(&workingWeight)
	svc := service.NewWorkoutService(repo, exercises)

	workingWeight = 90
	updated, err := svc.UpdateWorkoutExercise(ctx, testUserID, testWorkoutExerciseID, 0, models.UpdateWorkoutExerciseInput{
		Notes: ptrStr("tempo paused reps"),
	})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, int32(1), updated.Version)
	require.Len(t, updated.WorkoutExercises, 1)
	require.NotNil(t, updated.WorkoutExercises[0].Notes)
	assert.Equal(t, "tempo paused reps", *updated.WorkoutExercises[0].Notes)
	require.NotNil(t, updated.WorkoutExercises[0].WorkingWeightSnapshot)
	assert.Equal(t, 82.5, *updated.WorkoutExercises[0].WorkingWeightSnapshot)
	require.NotNil(t, updated.WorkoutExercises[0].Exercise)
	require.NotNil(t, updated.WorkoutExercises[0].Exercise.WorkingWeight)
	assert.Equal(t, 90.0, *updated.WorkoutExercises[0].Exercise.WorkingWeight)

	cleared, err := svc.UpdateWorkoutExercise(ctx, testUserID, testWorkoutExerciseID, 1, models.UpdateWorkoutExerciseInput{
		SetNotes: true,
	})

	require.NoError(t, err)
	require.NotNil(t, cleared)
	assert.Equal(t, int32(2), cleared.Version)
	require.Len(t, cleared.WorkoutExercises, 1)
	assert.Nil(t, cleared.WorkoutExercises[0].Notes)
	require.NotNil(t, cleared.WorkoutExercises[0].WorkingWeightSnapshot)
	assert.Equal(t, 82.5, *cleared.WorkoutExercises[0].WorkingWeightSnapshot)
}

func TestWorkoutService_RemoveExercise_KeepsEmptyDailyLog(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 1)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

	log, err := svc.RemoveWorkoutExercise(ctx, testUserID, testWorkoutExerciseID, 1)

	require.NoError(t, err)
	require.NotNil(t, log)
	assert.Equal(t, testDailyLogID, log.ID)
	assert.Equal(t, int32(2), log.Version)
	assert.NotNil(t, log.WorkoutExercises)
	assert.Empty(t, log.WorkoutExercises)
	require.Contains(t, repo.dailyLogs, testDailyLogID)
	assert.Empty(t, repo.dailyLogs[testDailyLogID].WorkoutExercises)
}

func TestWorkoutService_ChildMutationsRejectStaleVersionWithCurrentAggregateAndNoMutation(t *testing.T) {
	cases := []struct {
		name string
		call func(svc service.WorkoutService) (*models.DailyLog, error)
	}{
		{
			name: "add exercise",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.AddWorkoutExercise(ctx, testUserID, testWorkoutDate, 8, models.AddWorkoutExerciseInput{
					ExerciseID: "exercise-new",
				})
			},
		},
		{
			name: "update exercise",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.UpdateWorkoutExercise(ctx, testUserID, "workout-exercise-1", 8, models.UpdateWorkoutExerciseInput{
					Notes: ptrStr("stale notes"),
				})
			},
		},
		{
			name: "remove exercise",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.RemoveWorkoutExercise(ctx, testUserID, "workout-exercise-2", 8)
			},
		},
		{
			name: "reorder exercises",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.ReorderWorkoutExercises(ctx, testUserID, testWorkoutDate, 8, []string{"workout-exercise-2", "workout-exercise-1"})
			},
		},
		{
			name: "add set",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.AddWorkoutSet(ctx, testUserID, "workout-exercise-1", 8, models.AddWorkoutSetInput{
					Weight: 100,
					Reps:   5,
				})
			},
		},
		{
			name: "update set",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.UpdateWorkoutSet(ctx, testUserID, "set-1", 8, models.UpdateWorkoutSetInput{
					Weight: ptrFloat64(105),
				})
			},
		},
		{
			name: "remove set",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.RemoveWorkoutSet(ctx, testUserID, "set-2", 8)
			},
		},
		{
			name: "reorder sets",
			call: func(svc service.WorkoutService) (*models.DailyLog, error) {
				return svc.ReorderWorkoutSets(ctx, testUserID, "workout-exercise-1", 8, []string{"set-2", "set-1"})
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := newFakeWorkoutRepo()
			aggregate := repo.seedDailyLog(testWorkoutDate, 9)
			aggregate.DailyLog.Notes = ptrStr("current aggregate")
			repo.seedWorkoutExercise(aggregate.DailyLog.ID, "workout-exercise-1", "exercise-1", 1, ptrFloat64(82.5))
			repo.seedWorkoutExercise(aggregate.DailyLog.ID, "workout-exercise-2", "exercise-2", 2, nil)
			repo.seedWorkoutSet("workout-exercise-1", "set-1", 1)
			repo.seedWorkoutSet("workout-exercise-1", "set-2", 2)
			before := cloneAggregate(repo.dailyLogs[testDailyLogID])
			svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(ptrFloat64(90)))

			log, err := tc.call(svc)

			require.Nil(t, log)
			var conflictErr *models.DailyLogConflictErr
			require.ErrorAs(t, err, &conflictErr)
			assert.Equal(t, int32(9), conflictErr.CurrentVersion)
			require.NotNil(t, conflictErr.CurrentDailyLog)
			assert.Equal(t, int32(9), conflictErr.CurrentDailyLog.Version)
			require.NotNil(t, conflictErr.CurrentDailyLog.Notes)
			assert.Equal(t, "current aggregate", *conflictErr.CurrentDailyLog.Notes)
			assert.Equal(t, []string{"workout-exercise-1", "workout-exercise-2"}, modelWorkoutExerciseIDs(conflictErr.CurrentDailyLog.WorkoutExercises))
			assert.Equal(t, []string{"set-1", "set-2"}, modelWorkoutSetIDs(conflictErr.CurrentDailyLog.WorkoutExercises[0].Sets))
			assert.Equal(t, before, repo.dailyLogs[testDailyLogID])
			assert.Equal(t, 0, repo.mutationHits)
		})
	}
}

func TestWorkoutService_AddSet_ValidatesWeightRepsRpeRir(t *testing.T) {
	invalidCases := []struct {
		name  string
		input models.AddWorkoutSetInput
	}{
		{name: "zero weight", input: models.AddWorkoutSetInput{Weight: 0, Reps: 5}},
		{name: "negative weight", input: models.AddWorkoutSetInput{Weight: -1, Reps: 5}},
		{name: "zero reps", input: models.AddWorkoutSetInput{Weight: 100, Reps: 0}},
		{name: "negative reps", input: models.AddWorkoutSetInput{Weight: 100, Reps: -1}},
		{name: "low rpe", input: models.AddWorkoutSetInput{Weight: 100, Reps: 5, RPE: ptrFloat64(0)}},
		{name: "high rpe", input: models.AddWorkoutSetInput{Weight: 100, Reps: 5, RPE: ptrFloat64(11)}},
		{name: "low rir", input: models.AddWorkoutSetInput{Weight: 100, Reps: 5, RIR: ptrInt32(-1)}},
		{name: "high rir", input: models.AddWorkoutSetInput{Weight: 100, Reps: 5, RIR: ptrInt32(11)}},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := newFakeWorkoutRepo()
			aggregate := repo.seedDailyLog(testWorkoutDate, 0)
			repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
			svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

			log, err := svc.AddWorkoutSet(ctx, testUserID, testWorkoutExerciseID, 0, tc.input)

			require.Nil(t, log)
			var validationErr *models.DailyLogValidationErr
			require.ErrorAs(t, err, &validationErr)
			assert.Equal(t, 0, repo.lockExerciseHit)
			assert.Empty(t, repo.dailyLogs[testDailyLogID].WorkoutExercises[0].Sets)
		})
	}

	t.Run("valid boundaries", func(t *testing.T) {
		repo := newFakeWorkoutRepo()
		aggregate := repo.seedDailyLog(testWorkoutDate, 0)
		repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
		svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

		log, err := svc.AddWorkoutSet(ctx, testUserID, testWorkoutExerciseID, 0, models.AddWorkoutSetInput{
			Weight: 100,
			Reps:   5,
			RPE:    ptrFloat64(10),
			RIR:    ptrInt32(0),
		})

		require.NoError(t, err)
		require.Len(t, log.WorkoutExercises, 1)
		require.Len(t, log.WorkoutExercises[0].Sets, 1)
		assert.Equal(t, int32(1), log.WorkoutExercises[0].Sets[0].SetNumber)
		require.NotNil(t, log.WorkoutExercises[0].Sets[0].RPE)
		assert.Equal(t, 10.0, *log.WorkoutExercises[0].Sets[0].RPE)
		require.NotNil(t, log.WorkoutExercises[0].Sets[0].RIR)
		assert.Equal(t, int32(0), *log.WorkoutExercises[0].Sets[0].RIR)
		assert.Equal(t, int32(1), log.Version)
	})
}

func TestWorkoutService_AddSet_MissingParentReturnsNotFoundWithoutMutation(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 3)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

	log, err := svc.AddWorkoutSet(ctx, testUserID, "missing-workout-exercise", 3, models.AddWorkoutSetInput{
		Weight: 100,
		Reps:   5,
	})

	require.Nil(t, log)
	var notFoundErr *models.DailyLogNotFoundErr
	require.ErrorAs(t, err, &notFoundErr)
	assert.Equal(t, int32(3), repo.dailyLogs[testDailyLogID].DailyLog.Version)
	assert.Equal(t, 0, repo.mutationHits)
	assert.Empty(t, repo.dailyLogs[testDailyLogID].WorkoutExercises[0].Sets)
}

func TestWorkoutService_AddSet_InsertsAtSetNumberAndReindexes(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 3)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
	repo.seedWorkoutSet(testWorkoutExerciseID, "set-1", 1)
	repo.seedWorkoutSet(testWorkoutExerciseID, "set-2", 2)
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

	log, err := svc.AddWorkoutSet(ctx, testUserID, testWorkoutExerciseID, 3, models.AddWorkoutSetInput{
		SetNumber: ptrInt32(2),
		Weight:    125,
		Reps:      3,
		Notes:     ptrStr("top single"),
	})

	require.NoError(t, err)
	require.NotNil(t, log)
	assert.Equal(t, int32(4), log.Version)
	require.Len(t, log.WorkoutExercises, 1)
	require.Len(t, log.WorkoutExercises[0].Sets, 3)
	assert.Equal(t, []string{"set-1", "workout-set-1", "set-2"}, modelWorkoutSetIDs(log.WorkoutExercises[0].Sets))
	assert.Equal(t, []int32{1, 2, 3}, []int32{
		log.WorkoutExercises[0].Sets[0].SetNumber,
		log.WorkoutExercises[0].Sets[1].SetNumber,
		log.WorkoutExercises[0].Sets[2].SetNumber,
	})
	assert.Equal(t, 125.0, log.WorkoutExercises[0].Sets[1].Weight)
	require.NotNil(t, log.WorkoutExercises[0].Sets[1].Notes)
	assert.Equal(t, "top single", *log.WorkoutExercises[0].Sets[1].Notes)
}

func TestWorkoutService_UpdateSet_RejectsEmptyInputWithoutVersionChange(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 7)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
	repo.seedWorkoutSet(testWorkoutExerciseID, testWorkoutSetID, 1)
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

	log, err := svc.UpdateWorkoutSet(ctx, testUserID, testWorkoutSetID, 7, models.UpdateWorkoutSetInput{})

	require.Nil(t, log)
	var validationErr *models.DailyLogValidationErr
	require.ErrorAs(t, err, &validationErr)
	assert.Equal(t, int32(7), repo.dailyLogs[testDailyLogID].DailyLog.Version)
	require.Len(t, repo.dailyLogs[testDailyLogID].WorkoutExercises, 1)
	assert.Equal(t, []string{testWorkoutSetID}, repoWorkoutSetIDs(repo.dailyLogs[testDailyLogID].WorkoutExercises[0].Sets))
}

func TestWorkoutService_UpdateSet_RejectsSameSetNumberOnlyWithoutVersionChange(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 7)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
	repo.seedWorkoutSet(testWorkoutExerciseID, testWorkoutSetID, 1)
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

	log, err := svc.UpdateWorkoutSet(ctx, testUserID, testWorkoutSetID, 7, models.UpdateWorkoutSetInput{
		SetNumber: ptrInt32(1),
	})

	require.Nil(t, log)
	var validationErr *models.DailyLogValidationErr
	require.ErrorAs(t, err, &validationErr)
	assert.Equal(t, int32(7), repo.dailyLogs[testDailyLogID].DailyLog.Version)
	require.Len(t, repo.dailyLogs[testDailyLogID].WorkoutExercises, 1)
	require.Len(t, repo.dailyLogs[testDailyLogID].WorkoutExercises[0].Sets, 1)
	assert.Equal(t, int32(1), repo.dailyLogs[testDailyLogID].WorkoutExercises[0].Sets[0].SetNumber)
}

func TestWorkoutService_UpdateSet_PersistsValuesAndClearsNullableFields(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 7)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
	repo.seedWorkoutSet(testWorkoutExerciseID, testWorkoutSetID, 1)
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

	updated, err := svc.UpdateWorkoutSet(ctx, testUserID, testWorkoutSetID, 7, models.UpdateWorkoutSetInput{
		Weight:   ptrFloat64(112.5),
		Reps:     ptrInt32(8),
		RPE:      ptrFloat64(8.5),
		SetRPE:   true,
		RIR:      ptrInt32(2),
		SetRIR:   true,
		Notes:    ptrStr("clean reps"),
		SetNotes: true,
	})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, int32(8), updated.Version)
	require.Len(t, updated.WorkoutExercises, 1)
	require.Len(t, updated.WorkoutExercises[0].Sets, 1)
	set := updated.WorkoutExercises[0].Sets[0]
	assert.Equal(t, 112.5, set.Weight)
	assert.Equal(t, int32(8), set.Reps)
	require.NotNil(t, set.RPE)
	assert.Equal(t, 8.5, *set.RPE)
	require.NotNil(t, set.RIR)
	assert.Equal(t, int32(2), *set.RIR)
	require.NotNil(t, set.Notes)
	assert.Equal(t, "clean reps", *set.Notes)

	cleared, err := svc.UpdateWorkoutSet(ctx, testUserID, testWorkoutSetID, 8, models.UpdateWorkoutSetInput{
		SetRPE:   true,
		SetRIR:   true,
		SetNotes: true,
	})

	require.NoError(t, err)
	require.NotNil(t, cleared)
	assert.Equal(t, int32(9), cleared.Version)
	require.Len(t, cleared.WorkoutExercises, 1)
	require.Len(t, cleared.WorkoutExercises[0].Sets, 1)
	clearedSet := cleared.WorkoutExercises[0].Sets[0]
	assert.Equal(t, 112.5, clearedSet.Weight)
	assert.Equal(t, int32(8), clearedSet.Reps)
	assert.Nil(t, clearedSet.RPE)
	assert.Nil(t, clearedSet.RIR)
	assert.Nil(t, clearedSet.Notes)
}

func TestWorkoutService_UpdateSet_ReindexesWhenSetNumberChanges(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 7)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
	repo.seedWorkoutSet(testWorkoutExerciseID, "set-1", 1)
	repo.seedWorkoutSet(testWorkoutExerciseID, "set-2", 2)
	repo.seedWorkoutSet(testWorkoutExerciseID, "set-3", 3)
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

	log, err := svc.UpdateWorkoutSet(ctx, testUserID, "set-3", 7, models.UpdateWorkoutSetInput{
		SetNumber: ptrInt32(1),
	})

	require.NoError(t, err)
	require.NotNil(t, log)
	assert.Equal(t, int32(8), log.Version)
	require.Len(t, log.WorkoutExercises, 1)
	require.Len(t, log.WorkoutExercises[0].Sets, 3)
	assert.Equal(t, []string{"set-3", "set-1", "set-2"}, []string{
		log.WorkoutExercises[0].Sets[0].ID,
		log.WorkoutExercises[0].Sets[1].ID,
		log.WorkoutExercises[0].Sets[2].ID,
	})
	assert.Equal(t, []int32{1, 2, 3}, []int32{
		log.WorkoutExercises[0].Sets[0].SetNumber,
		log.WorkoutExercises[0].Sets[1].SetNumber,
		log.WorkoutExercises[0].Sets[2].SetNumber,
	})
}

func TestWorkoutService_RemoveSet_RemovesOneSetAndReindexes(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 5)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
	repo.seedWorkoutSet(testWorkoutExerciseID, "set-1", 1)
	repo.seedWorkoutSet(testWorkoutExerciseID, "set-2", 2)
	repo.seedWorkoutSet(testWorkoutExerciseID, "set-3", 3)
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

	log, err := svc.RemoveWorkoutSet(ctx, testUserID, "set-2", 5)

	require.NoError(t, err)
	require.NotNil(t, log)
	assert.Equal(t, int32(6), log.Version)
	require.Len(t, log.WorkoutExercises, 1)
	require.Len(t, log.WorkoutExercises[0].Sets, 2)
	assert.Equal(t, []string{"set-1", "set-3"}, modelWorkoutSetIDs(log.WorkoutExercises[0].Sets))
	assert.Equal(t, []int32{1, 2}, []int32{
		log.WorkoutExercises[0].Sets[0].SetNumber,
		log.WorkoutExercises[0].Sets[1].SetNumber,
	})
	_, stillIndexed := repo.setToExercise["set-2"]
	assert.False(t, stillIndexed)
}

func TestWorkoutService_ReorderExercises_SuccessIncrementsVersionAndReindexes(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 4)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, "workout-exercise-1", "exercise-1", 1, nil)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, "workout-exercise-2", "exercise-2", 2, nil)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, "workout-exercise-3", "exercise-3", 3, nil)
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

	log, err := svc.ReorderWorkoutExercises(ctx, testUserID, testWorkoutDate, 4, []string{"workout-exercise-3", "workout-exercise-1", "workout-exercise-2"})

	require.NoError(t, err)
	require.NotNil(t, log)
	assert.Equal(t, int32(5), log.Version)
	require.Len(t, log.WorkoutExercises, 3)
	assert.Equal(t, []string{"workout-exercise-3", "workout-exercise-1", "workout-exercise-2"}, modelWorkoutExerciseIDs(log.WorkoutExercises))
	assert.Equal(t, []int32{1, 2, 3}, []int32{
		log.WorkoutExercises[0].Position,
		log.WorkoutExercises[1].Position,
		log.WorkoutExercises[2].Position,
	})
}

func TestWorkoutService_ReorderExercises_RejectsMissingDuplicateOrForeignIDs(t *testing.T) {
	cases := []struct {
		name       string
		orderedIDs []string
	}{
		{name: "missing", orderedIDs: []string{"workout-exercise-1"}},
		{name: "duplicate", orderedIDs: []string{"workout-exercise-1", "workout-exercise-1"}},
		{name: "foreign", orderedIDs: []string{"workout-exercise-1", "foreign-exercise"}},
		{name: "extra", orderedIDs: []string{"workout-exercise-1", "workout-exercise-2", "foreign-exercise"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := newFakeWorkoutRepo()
			aggregate := repo.seedDailyLog(testWorkoutDate, 4)
			repo.seedWorkoutExercise(aggregate.DailyLog.ID, "workout-exercise-1", "exercise-1", 1, nil)
			repo.seedWorkoutExercise(aggregate.DailyLog.ID, "workout-exercise-2", "exercise-2", 2, nil)
			svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

			log, err := svc.ReorderWorkoutExercises(ctx, testUserID, testWorkoutDate, 4, tc.orderedIDs)

			require.Nil(t, log)
			var validationErr *models.DailyLogValidationErr
			require.ErrorAs(t, err, &validationErr)
			assert.Equal(t, int32(4), repo.dailyLogs[testDailyLogID].DailyLog.Version)
			assert.Equal(t, []string{"workout-exercise-1", "workout-exercise-2"}, repoWorkoutExerciseIDs(repo.dailyLogs[testDailyLogID].WorkoutExercises))
		})
	}
}

func TestWorkoutService_ReorderSets_SuccessIncrementsVersionAndReindexes(t *testing.T) {
	repo := newFakeWorkoutRepo()
	aggregate := repo.seedDailyLog(testWorkoutDate, 4)
	repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
	repo.seedWorkoutSet(testWorkoutExerciseID, "set-1", 1)
	repo.seedWorkoutSet(testWorkoutExerciseID, "set-2", 2)
	repo.seedWorkoutSet(testWorkoutExerciseID, "set-3", 3)
	svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

	log, err := svc.ReorderWorkoutSets(ctx, testUserID, testWorkoutExerciseID, 4, []string{"set-3", "set-1", "set-2"})

	require.NoError(t, err)
	require.NotNil(t, log)
	assert.Equal(t, int32(5), log.Version)
	require.Len(t, log.WorkoutExercises, 1)
	require.Len(t, log.WorkoutExercises[0].Sets, 3)
	assert.Equal(t, []string{"set-3", "set-1", "set-2"}, modelWorkoutSetIDs(log.WorkoutExercises[0].Sets))
	assert.Equal(t, []int32{1, 2, 3}, []int32{
		log.WorkoutExercises[0].Sets[0].SetNumber,
		log.WorkoutExercises[0].Sets[1].SetNumber,
		log.WorkoutExercises[0].Sets[2].SetNumber,
	})
}

func TestWorkoutService_ReorderSets_RejectsMissingDuplicateOrForeignIDs(t *testing.T) {
	cases := []struct {
		name       string
		orderedIDs []string
	}{
		{name: "missing", orderedIDs: []string{"set-1"}},
		{name: "duplicate", orderedIDs: []string{"set-1", "set-1"}},
		{name: "foreign", orderedIDs: []string{"set-1", "foreign-set"}},
		{name: "extra", orderedIDs: []string{"set-1", "set-2", "foreign-set"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := newFakeWorkoutRepo()
			aggregate := repo.seedDailyLog(testWorkoutDate, 4)
			repo.seedWorkoutExercise(aggregate.DailyLog.ID, testWorkoutExerciseID, testExerciseID, 1, nil)
			repo.seedWorkoutSet(testWorkoutExerciseID, "set-1", 1)
			repo.seedWorkoutSet(testWorkoutExerciseID, "set-2", 2)
			svc := service.NewWorkoutService(repo, exerciseRepoWithRecord(nil))

			log, err := svc.ReorderWorkoutSets(ctx, testUserID, testWorkoutExerciseID, 4, tc.orderedIDs)

			require.Nil(t, log)
			var validationErr *models.DailyLogValidationErr
			require.ErrorAs(t, err, &validationErr)
			assert.Equal(t, int32(4), repo.dailyLogs[testDailyLogID].DailyLog.Version)
			require.Len(t, repo.dailyLogs[testDailyLogID].WorkoutExercises, 1)
			assert.Equal(t, []string{"set-1", "set-2"}, repoWorkoutSetIDs(repo.dailyLogs[testDailyLogID].WorkoutExercises[0].Sets))
		})
	}
}

func exerciseRepoWithRecord(workingWeight *float64) *fakeExerciseRepo {
	return &fakeExerciseRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
			return &models.ExerciseRecord{
				ID:            id,
				UserID:        userID,
				Name:          "Bench Press",
				MuscleGroups:  []string{"chest"},
				WorkingWeight: workingWeight,
				IsActive:      true,
				CreatedAt:     "2026-06-19T00:00:00Z",
				UpdatedAt:     "2026-06-19T00:00:00Z",
			}, nil
		},
	}
}

func exerciseRepoWithMutableWorkingWeight(workingWeight *float64) *fakeExerciseRepo {
	return &fakeExerciseRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
			var snapshot *float64
			if workingWeight != nil {
				value := *workingWeight
				snapshot = &value
			}
			return &models.ExerciseRecord{
				ID:            id,
				UserID:        userID,
				Name:          "Bench Press",
				MuscleGroups:  []string{"chest"},
				WorkingWeight: snapshot,
				IsActive:      true,
				CreatedAt:     "2026-06-19T00:00:00Z",
				UpdatedAt:     "2026-06-19T00:00:00Z",
			}, nil
		},
	}
}

func dateKey(userID string, date models.Date) string {
	return userID + "|" + date.String()
}

func cloneAggregate(in *atlasPostgres.DailyLogAggregate) *atlasPostgres.DailyLogAggregate {
	if in == nil {
		return nil
	}
	out := &atlasPostgres.DailyLogAggregate{
		DailyLog:         in.DailyLog,
		WorkoutExercises: make([]atlasPostgres.WorkoutExerciseRecord, len(in.WorkoutExercises)),
	}
	for i := range in.WorkoutExercises {
		out.WorkoutExercises[i] = in.WorkoutExercises[i]
		if in.WorkoutExercises[i].Sets != nil {
			out.WorkoutExercises[i].Sets = append([]atlasPostgres.WorkoutSetRecord{}, in.WorkoutExercises[i].Sets...)
		}
	}
	return out
}

func findRepoWorkoutExercise(aggregate *atlasPostgres.DailyLogAggregate, id string) *atlasPostgres.WorkoutExerciseRecord {
	if aggregate == nil {
		return nil
	}
	for i := range aggregate.WorkoutExercises {
		if aggregate.WorkoutExercises[i].ID == id {
			return &aggregate.WorkoutExercises[i]
		}
	}
	return nil
}

func findRepoWorkoutSet(exercise *atlasPostgres.WorkoutExerciseRecord, id string) *atlasPostgres.WorkoutSetRecord {
	if exercise == nil {
		return nil
	}
	for i := range exercise.Sets {
		if exercise.Sets[i].ID == id {
			return &exercise.Sets[i]
		}
	}
	return nil
}

func insertWorkoutExercise(aggregate *atlasPostgres.DailyLogAggregate, record atlasPostgres.WorkoutExerciseRecord) {
	position := int(record.Position)
	if position < 1 || position > len(aggregate.WorkoutExercises)+1 {
		position = len(aggregate.WorkoutExercises) + 1
		record.Position = int32(position) //nolint:gosec // test fixture sizes are bounded by in-memory slices.
	}
	aggregate.WorkoutExercises = append(aggregate.WorkoutExercises, atlasPostgres.WorkoutExerciseRecord{})
	copy(aggregate.WorkoutExercises[position:], aggregate.WorkoutExercises[position-1:])
	aggregate.WorkoutExercises[position-1] = record
	reindexWorkoutExercises(aggregate)
}

func insertWorkoutSet(exercise *atlasPostgres.WorkoutExerciseRecord, record atlasPostgres.WorkoutSetRecord) {
	position := int(record.SetNumber)
	if position < 1 || position > len(exercise.Sets)+1 {
		position = len(exercise.Sets) + 1
		record.SetNumber = int32(position) //nolint:gosec // test fixture sizes are bounded by in-memory slices.
	}
	exercise.Sets = append(exercise.Sets, atlasPostgres.WorkoutSetRecord{})
	copy(exercise.Sets[position:], exercise.Sets[position-1:])
	exercise.Sets[position-1] = record
	reindexWorkoutSets(exercise)
}

func reorderWorkoutExercises(aggregate *atlasPostgres.DailyLogAggregate, orderedIDs []string) {
	byID := map[string]atlasPostgres.WorkoutExerciseRecord{}
	for _, exercise := range aggregate.WorkoutExercises {
		byID[exercise.ID] = exercise
	}
	aggregate.WorkoutExercises = make([]atlasPostgres.WorkoutExerciseRecord, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		aggregate.WorkoutExercises = append(aggregate.WorkoutExercises, byID[id])
	}
	reindexWorkoutExercises(aggregate)
}

func reorderWorkoutSets(exercise *atlasPostgres.WorkoutExerciseRecord, orderedIDs []string) {
	byID := map[string]atlasPostgres.WorkoutSetRecord{}
	for _, set := range exercise.Sets {
		byID[set.ID] = set
	}
	exercise.Sets = make([]atlasPostgres.WorkoutSetRecord, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		exercise.Sets = append(exercise.Sets, byID[id])
	}
	reindexWorkoutSets(exercise)
}

func reindexWorkoutExercises(aggregate *atlasPostgres.DailyLogAggregate) {
	for i := range aggregate.WorkoutExercises {
		aggregate.WorkoutExercises[i].Position = int32(i + 1)
	}
}

func reindexWorkoutSets(exercise *atlasPostgres.WorkoutExerciseRecord) {
	for i := range exercise.Sets {
		exercise.Sets[i].SetNumber = int32(i + 1)
	}
}

func repoWorkoutExerciseIDs(records []atlasPostgres.WorkoutExerciseRecord) []string {
	out := make([]string, len(records))
	for i := range records {
		out[i] = records[i].ID
	}
	return out
}

func repoWorkoutSetIDs(records []atlasPostgres.WorkoutSetRecord) []string {
	out := make([]string, len(records))
	for i := range records {
		out[i] = records[i].ID
	}
	return out
}

func modelWorkoutExerciseIDs(records []models.WorkoutExercise) []string {
	out := make([]string, len(records))
	for i := range records {
		out[i] = records[i].ID
	}
	return out
}

func modelWorkoutSetIDs(records []models.WorkoutSet) []string {
	out := make([]string, len(records))
	for i := range records {
		out[i] = records[i].ID
	}
	return out
}

func movedOrder(current []string, targetID string, position int32) []string {
	from := -1
	for i, id := range current {
		if id == targetID {
			from = i
			break
		}
	}
	if from == -1 {
		return current
	}
	to := int(position) - 1
	if to < 0 {
		to = 0
	}
	if to >= len(current) {
		to = len(current) - 1
	}
	ordered := append([]string(nil), current...)
	target := ordered[from]
	ordered = append(ordered[:from], ordered[from+1:]...)
	ordered = append(ordered[:to], append([]string{target}, ordered[to:]...)...)
	return ordered
}
