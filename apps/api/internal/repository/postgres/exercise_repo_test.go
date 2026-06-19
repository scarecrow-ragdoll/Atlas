// FILE: apps/api/internal/repository/postgres/exercise_repo_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify PostgreSQL exercise repository behavior against the goose-managed test database.
//   SCOPE: Exercise CRUD, soft archive/restore, pagination, cursor, duplicate names, working weight validation, media CRUD, user-scoped isolation, and unavailable database skip semantics.
//   DEPENDS: apps/api/internal/repository/postgres, apps/api/internal/atlas/repository/postgres, apps/api/internal/atlas/models, apps/api/internal/testinfra.
//   LINKS: M-API / V-M-API / WAVE-02 / TEST-W02-001 / TEST-W02-005 / TEST-W02-007 / TEST-W02-008 / TEST-W02-009 / TEST-W02-010 / TEST-W02-011 / TEST-W02-012 / TEST-W02-014.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Fixed TestExerciseRepo_UserScoped to use FK-compatible otherUID; added ensureAtlasUserWithDisplayName helper.
// END_CHANGE_SUMMARY

package postgres_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	atlasModels "monorepo-template/apps/api/internal/atlas/models"
	atlasRepo "monorepo-template/apps/api/internal/atlas/repository/postgres"
	postgresrepo "monorepo-template/apps/api/internal/repository/postgres"
	"monorepo-template/apps/api/internal/testinfra"
)

var ctx = context.Background()

func strPtr(s string) *string { return &s }

func flPtr(f float64) *float64 { return &f }

func TestExerciseRepo_Create_Success(t *testing.T) {
	pool, uid := exerciseTestSetup(t)

	repo := atlasRepo.NewExerciseRepository(pool)
	result, err := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{
		Name:          "Bench Press",
		MuscleGroups:  []string{"chest", "triceps"},
		Description:   strPtr("Barbell bench press"),
		PersonalNotes: strPtr("Use a spotter"),
		WorkingWeight: flPtr(80.5),
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, "Bench Press", result.Name)
	assert.Equal(t, []string{"chest", "triceps"}, result.MuscleGroups)
	assert.Equal(t, "Barbell bench press", *result.Description)
	assert.Equal(t, "Use a spotter", *result.PersonalNotes)
	assert.Equal(t, 80.5, *result.WorkingWeight)
	assert.True(t, result.IsActive)
	assert.NotEmpty(t, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)
}

func TestExerciseRepo_Create_DuplicateNamesAllowed(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	r1, err := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "Squat"})
	require.NoError(t, err)
	require.NotNil(t, r1)

	r2, err := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "Squat"})
	require.NoError(t, err)
	require.NotNil(t, r2)
	assert.NotEqual(t, r1.ID, r2.ID)
}

func TestExerciseRepo_Create_NullableFields(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	result, err := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "Push-ups"})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Description)
	assert.Nil(t, result.PersonalNotes)
	assert.Nil(t, result.WorkingWeight)
}

func TestExerciseRepo_GetByID_Success(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	created, err := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "Deadlift", WorkingWeight: flPtr(100)})
	require.NoError(t, err)

	result, err := repo.GetByID(ctx, uid, created.ID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, "Deadlift", result.Name)
	assert.Equal(t, 100.0, *result.WorkingWeight)
}

func TestExerciseRepo_GetByID_NotFound_WrongUser(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	created, err := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "Test"})
	require.NoError(t, err)

	otherUID := "00000000-0000-0000-0000-000000000001"
	result, err := repo.GetByID(ctx, otherUID, created.ID)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestExerciseRepo_GetByID_NotFound_MissingID(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	result, err := repo.GetByID(ctx, uid, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestExerciseRepo_List_ActiveOnly(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "A"})
	repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "B"})
	repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "C"})

	results, err := repo.List(ctx, uid, true, 10)
	require.NoError(t, err)
	assert.Len(t, results, 3)
}

func TestExerciseRepo_List_ExcludesArchived(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	ex, _ := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "ToArchive"})
	repo.Archive(ctx, uid, ex.ID)

	results, err := repo.List(ctx, uid, true, 10)
	require.NoError(t, err)
	for _, r := range results {
		assert.NotEqual(t, ex.ID, r.ID)
	}
}

func TestExerciseRepo_List_OrderedByName(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "C"})
	repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "A"})
	repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "B"})

	results, err := repo.List(ctx, uid, true, 10)
	require.NoError(t, err)
	require.Len(t, results, 3)
	assert.Equal(t, "A", results[0].Name)
	assert.Equal(t, "B", results[1].Name)
	assert.Equal(t, "C", results[2].Name)
}

func TestExerciseRepo_ListCursor_Works(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "A"})
	repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "B"})
	repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "C"})

	results, err := repo.ListCursor(ctx, uid, true, "A", 10)
	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.Equal(t, "B", results[0].Name)
	assert.Equal(t, "C", results[1].Name)
}

func TestExerciseRepo_ListCursor_EmptyAtEnd(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "A"})

	results, err := repo.ListCursor(ctx, uid, true, "Z", 10)
	require.NoError(t, err)
	assert.Len(t, results, 0)
}

func TestExerciseRepo_ListAll_IncludeInactive(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	_, _ = repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "Active"})
	archived, _ := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "Archived"})
	repo.Archive(ctx, uid, archived.ID)

	allInactive, err := repo.ListAll(ctx, uid, true)
	require.NoError(t, err)
	assert.Len(t, allInactive, 2)

	activeOnly, err := repo.ListAll(ctx, uid, false)
	require.NoError(t, err)
	assert.Len(t, activeOnly, 1)
	assert.Equal(t, "Active", activeOnly[0].Name)
}

func TestExerciseRepo_Count_Matches(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "X"})
	repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "Y"})

	count, err := repo.Count(ctx, uid, true)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestExerciseRepo_Update_Fields(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	ex, _ := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "OldName"})

	updated, err := repo.Update(ctx, uid, ex.ID, atlasModels.UpdateExerciseInput{
		Name:          strPtr("NewName"),
		MuscleGroups:  &[]string{"back"},
		Description:   strPtr("Updated desc"),
		PersonalNotes: strPtr("Updated notes"),
		WorkingWeight: flPtr(90),
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "NewName", updated.Name)
	assert.Equal(t, []string{"back"}, updated.MuscleGroups)
	assert.Equal(t, "Updated desc", *updated.Description)
	assert.Equal(t, "Updated notes", *updated.PersonalNotes)
	assert.Equal(t, 90.0, *updated.WorkingWeight)
}

func TestExerciseRepo_Update_Partial(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	ex, _ := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "Original", WorkingWeight: flPtr(50)})

	updated, err := repo.Update(ctx, uid, ex.ID, atlasModels.UpdateExerciseInput{Name: strPtr("Renamed")})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Renamed", updated.Name)
	assert.Equal(t, 50.0, *updated.WorkingWeight)
}

func TestExerciseRepo_Update_NotFound(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	result, err := repo.Update(ctx, uid, "00000000-0000-0000-0000-000000000000", atlasModels.UpdateExerciseInput{Name: strPtr("Nope")})
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestExerciseRepo_Archive_SetsInactive(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	ex, _ := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "ToArchive"})

	archived, err := repo.Archive(ctx, uid, ex.ID)
	require.NoError(t, err)
	require.NotNil(t, archived)
	assert.False(t, archived.IsActive)
}

func TestExerciseRepo_Archive_Idempotent(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	ex, _ := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "DoubleArchive"})

	repo.Archive(ctx, uid, ex.ID)
	archivedAgain, err := repo.Archive(ctx, uid, ex.ID)
	require.NoError(t, err)
	require.NotNil(t, archivedAgain)
	assert.False(t, archivedAgain.IsActive)
}

func TestExerciseRepo_Archive_NotFound(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	result, err := repo.Archive(ctx, uid, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestExerciseRepo_Restore_SetsActive(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	ex, _ := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "ToRestore"})
	repo.Archive(ctx, uid, ex.ID)

	restored, err := repo.Restore(ctx, uid, ex.ID)
	require.NoError(t, err)
	require.NotNil(t, restored)
	assert.True(t, restored.IsActive)
}

func TestExerciseRepo_Restore_NotFound(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	result, err := repo.Restore(ctx, uid, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestExerciseRepo_Media_CRUD(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	ex, _ := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "MediaTest"})

	created, err := repo.CreateMedia(ctx, uid, ex.ID, "video.mp4", "/storage/video.mp4", "video/mp4", 1024)
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "video.mp4", created.FileName)
	assert.Equal(t, ex.ID, created.ExerciseID)

	fetched, err := repo.GetMediaByID(ctx, uid, created.ID)
	require.NoError(t, err)
	require.NotNil(t, fetched)
	assert.Equal(t, created.ID, fetched.ID)

	list, err := repo.ListMediaByExercise(ctx, uid, ex.ID)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, created.ID, list[0].ID)

	deleted, err := repo.DeleteMedia(ctx, uid, created.ID)
	require.NoError(t, err)
	require.NotNil(t, deleted)
	assert.Equal(t, created.ID, deleted.ID)
	assert.Equal(t, "/storage/video.mp4", deleted.FilePath)

	fetchedAfterDelete, err := repo.GetMediaByID(ctx, uid, created.ID)
	require.NoError(t, err)
	assert.Nil(t, fetchedAfterDelete)
}

func TestExerciseRepo_Archive_DoesNotCascadeDeleteMedia(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	ex, _ := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "ArchiveNoCascade"})
	repo.CreateMedia(ctx, uid, ex.ID, "img.jpg", "/storage/img.jpg", "image/jpeg", 500)

	repo.Archive(ctx, uid, ex.ID)

	archivedCheck, _ := repo.GetByID(ctx, uid, ex.ID)
	require.NotNil(t, archivedCheck)
	assert.False(t, archivedCheck.IsActive)

	mediaList, err := repo.ListMediaByExercise(ctx, uid, ex.ID)
	require.NoError(t, err)
	assert.Len(t, mediaList, 1)
}

func TestExerciseRepo_UserScoped(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	otherUID := ensureAtlasUserWithDisplayName(t, pool, "other")

	ex, _ := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "Mine"})
	other, _ := repo.Create(ctx, otherUID, atlasModels.CreateExerciseInput{Name: "Theirs"})

	foundMine, _ := repo.GetByID(ctx, uid, ex.ID)
	require.NotNil(t, foundMine)
	assert.Equal(t, "Mine", foundMine.Name)

	cantFindTheirs, _ := repo.GetByID(ctx, uid, other.ID)
	assert.Nil(t, cantFindTheirs)

	myList, _ := repo.List(ctx, uid, true, 10)
	assert.Len(t, myList, 1)
	assert.Equal(t, "Mine", myList[0].Name)
}

func TestExerciseRepo_Media_WrongUser(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	otherUID := "00000000-0000-0000-0000-000000000003"

	ex, _ := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "Test"})
	media, _ := repo.CreateMedia(ctx, uid, ex.ID, "f.txt", "/p", "text/plain", 1)

	fetchedOther, err := repo.GetMediaByID(ctx, otherUID, media.ID)
	require.NoError(t, err)
	assert.Nil(t, fetchedOther)

	deletedOther, err := repo.DeleteMedia(ctx, otherUID, media.ID)
	require.NoError(t, err)
	assert.Nil(t, deletedOther)
}

func TestExerciseRepo_MediaRecordByID(t *testing.T) {
	pool, uid := exerciseTestSetup(t)
	repo := atlasRepo.NewExerciseRepository(pool)

	ex, _ := repo.Create(ctx, uid, atlasModels.CreateExerciseInput{Name: "RecordTest"})
	created, _ := repo.CreateMedia(ctx, uid, ex.ID, "test.jpg", "/path/test.jpg", "image/jpeg", 100)

	record, err := repo.GetMediaRecordByID(ctx, uid, created.ID)
	require.NoError(t, err)
	require.NotNil(t, record)
	assert.Equal(t, "/path/test.jpg", record.FilePath)
	assert.Equal(t, "image/jpeg", record.MimeType)
}

func exerciseTestSetup(t *testing.T) (*pgxpool.Pool, string) {
	t.Helper()
	dsn := testinfra.PostgresDSN()
	testinfra.RequireSafePostgresDSN(t, dsn)
	if err := postgresrepo.RunMigrations(dsn, zap.NewNop()); err != nil {
		if !testinfra.CoverageGateEnabled() {
			t.Skipf("postgres integration database is unavailable: %v", err)
		}
		require.NoError(t, err)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	truncateExerciseTables(t, pool)

	uid := ensureAtlasUser(t, pool)

	return pool, uid
}

func ensureAtlasUser(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	ctx := context.Background()
	var id string
	err := pool.QueryRow(ctx, `INSERT INTO atlas_users (display_name) VALUES ('test') ON CONFLICT DO NOTHING RETURNING id`).Scan(&id)
	if err != nil {
		pool.QueryRow(ctx, `SELECT id FROM atlas_users ORDER BY created_at ASC LIMIT 1`).Scan(&id)
	}
	if id == "" {
		pool.QueryRow(ctx, `INSERT INTO atlas_users (display_name) VALUES ('test') RETURNING id`).Scan(&id)
	}
	return id
}

func truncateExerciseTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `TRUNCATE exercise_media, exercises RESTART IDENTITY CASCADE`)
	if err == nil {
		return
	}
	ctx := context.Background()
	pool.Exec(ctx, `DELETE FROM exercise_media`)
	pool.Exec(ctx, `DELETE FROM exercises`)
}

func ensureAtlasUserWithDisplayName(t *testing.T, pool *pgxpool.Pool, displayName string) string {
	t.Helper()
	ctx := context.Background()
	var id string
	err := pool.QueryRow(ctx, `INSERT INTO atlas_users (display_name) VALUES ($1) RETURNING id`, displayName).Scan(&id)
	if err != nil {
		pool.QueryRow(ctx, `SELECT id FROM atlas_users WHERE display_name = $1 ORDER BY created_at ASC LIMIT 1`, displayName).Scan(&id)
	}
	if id == "" {
		pool.QueryRow(ctx, `INSERT INTO atlas_users (display_name) VALUES ($1) RETURNING id`, displayName).Scan(&id)
	}
	return id
}