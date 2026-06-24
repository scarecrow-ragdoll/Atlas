// FILE: apps/api/internal/atlas/repository/postgres/nutrition_migration_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Integration tests for WAVE-05 nutrition migration smoke test and factual daily nutrition schema constraints. Skipped when DB unavailable.
//   SCOPE: Migration smoke plus daily_nutrition_logs/daily_nutrition_entries shape and constraints. Requires test Postgres on port 17501.
//   DEPENDS: internal/testinfra for Postgres config, internal/repository/postgres for RunMigrations.
//   LINKS: M-API-NUTRITION / V-M-API / V-M-API-NUTRITION / EC-W05-006.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.3 - Added concurrent first-write get-or-create race coverage.
// END_CHANGE_SUMMARY

package postgres_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	postgresrepo "monorepo-template/apps/api/internal/repository/postgres"
	"monorepo-template/apps/api/internal/repository/postgres/generated"
	"monorepo-template/apps/api/internal/testinfra"
)

func TestWave05NutritionMigration_Smoke(t *testing.T) {
	dsn := testinfra.PostgresDSN()
	testinfra.RequireSafePostgresDSN(t, dsn)

	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("INTEGRATION_TESTS not set; skipping WAVE-05 migration smoke test")
	}

	err := postgresrepo.RunMigrations(dsn, zap.NewNop())
	require.NoError(t, err, "WAVE-05 migration should apply without error")
}

func TestDailyNutritionMigrationHasSnapshotAndAmountConstraints(t *testing.T) {
	pool := nutritionMigrationTestPool(t)
	ctx := context.Background()
	truncateDailyNutritionTables(t, pool)

	userID := seedNutritionMigrationUser(t, pool, "daily nutrition migration user")
	productID := seedNutritionMigrationProduct(t, pool, userID)

	var logID string
	err := pool.QueryRow(ctx, `
		INSERT INTO daily_nutrition_logs (user_id, date, notes)
		VALUES ($1, DATE '2026-06-24', 'day note')
		RETURNING id
	`, userID).Scan(&logID)
	require.NoError(t, err)

	var productName string
	var calories, protein, fat, carbs, amount float32
	err = pool.QueryRow(ctx, `
		INSERT INTO daily_nutrition_entries (
			daily_log_id,
			product_id,
			product_name_snapshot,
			calories_per_100g_snapshot,
			protein_per_100g_snapshot,
			fat_per_100g_snapshot,
			carbs_per_100g_snapshot,
			amount_grams,
			meal_label,
			notes,
			position
		)
		VALUES ($1, $2, 'Greek yogurt', 59, 10.2, 0.4, 3.6, 150, 'breakfast', 'with berries', 1)
		RETURNING product_name_snapshot, calories_per_100g_snapshot, protein_per_100g_snapshot,
			fat_per_100g_snapshot, carbs_per_100g_snapshot, amount_grams
	`, logID, productID).Scan(&productName, &calories, &protein, &fat, &carbs, &amount)
	require.NoError(t, err)
	require.Equal(t, "Greek yogurt", productName)
	require.Equal(t, float32(59), calories)
	require.Equal(t, float32(10.2), protein)
	require.Equal(t, float32(0.4), fat)
	require.Equal(t, float32(3.6), carbs)
	require.Equal(t, float32(150), amount)

	_, err = pool.Exec(ctx, `
		INSERT INTO daily_nutrition_entries (
			daily_log_id,
			product_id,
			product_name_snapshot,
			calories_per_100g_snapshot,
			protein_per_100g_snapshot,
			fat_per_100g_snapshot,
			carbs_per_100g_snapshot,
			amount_grams
		)
		VALUES ($1, $2, 'Invalid amount', 10, 1, 1, 1, 0)
	`, logID, productID)
	requirePgConstraint(t, err, "daily_nutrition_entries_amount_grams_check")

	for _, tc := range []struct {
		name       string
		column     string
		constraint string
	}{
		{name: "calories", column: "calories_per_100g_snapshot", constraint: "daily_nutrition_entries_calories_per_100g_snapshot_check"},
		{name: "protein", column: "protein_per_100g_snapshot", constraint: "daily_nutrition_entries_protein_per_100g_snapshot_check"},
		{name: "fat", column: "fat_per_100g_snapshot", constraint: "daily_nutrition_entries_fat_per_100g_snapshot_check"},
		{name: "carbs", column: "carbs_per_100g_snapshot", constraint: "daily_nutrition_entries_carbs_per_100g_snapshot_check"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err = pool.Exec(ctx, `
				INSERT INTO daily_nutrition_entries (
					daily_log_id,
					product_id,
					product_name_snapshot,
					calories_per_100g_snapshot,
					protein_per_100g_snapshot,
					fat_per_100g_snapshot,
					carbs_per_100g_snapshot,
					amount_grams
				)
				VALUES (
					$1,
					$2,
					'Invalid macro',
					CASE WHEN $3 = 'calories_per_100g_snapshot' THEN -1 ELSE 10 END,
					CASE WHEN $3 = 'protein_per_100g_snapshot' THEN -1 ELSE 1 END,
					CASE WHEN $3 = 'fat_per_100g_snapshot' THEN -1 ELSE 1 END,
					CASE WHEN $3 = 'carbs_per_100g_snapshot' THEN -1 ELSE 1 END,
					100
				)
			`, logID, productID, tc.column)
			requirePgConstraint(t, err, tc.constraint)
		})
	}
}

func TestDailyNutritionRepository_GetOrCreateIsUniqueAndListsByRange(t *testing.T) {
	pool := nutritionMigrationTestPool(t)
	truncateDailyNutritionTables(t, pool)
	ctx := context.Background()
	q := generated.New(pool)

	userID := seedNutritionMigrationUser(t, pool, "daily nutrition range user")
	userUUID := testUUID(t, userID)
	firstDate := testDate(2026, time.June, 24)
	secondDate := testDate(2026, time.June, 25)

	firstLog, err := q.CreateDailyNutritionLog(ctx, generated.CreateDailyNutritionLogParams{
		UserID: userUUID,
		Date:   firstDate,
		Notes:  pgtype.Text{},
	})
	require.NoError(t, err)

	duplicateLog, err := q.CreateDailyNutritionLog(ctx, generated.CreateDailyNutritionLogParams{
		UserID: userUUID,
		Date:   firstDate,
		Notes:  pgtype.Text{String: "duplicate write", Valid: true},
	})
	require.NoError(t, err)
	require.Equal(t, firstLog.ID, duplicateLog.ID)
	require.False(t, duplicateLog.Notes.Valid, "get-or-create conflict path must not update existing notes")
	require.Equal(t, firstLog.UpdatedAt, duplicateLog.UpdatedAt)

	secondLog, err := q.CreateDailyNutritionLog(ctx, generated.CreateDailyNutritionLogParams{
		UserID: userUUID,
		Date:   secondDate,
		Notes:  pgtype.Text{String: "second day", Valid: true},
	})
	require.NoError(t, err)

	logs, err := q.ListDailyNutritionLogsByRange(ctx, generated.ListDailyNutritionLogsByRangeParams{
		UserID:    userUUID,
		StartDate: firstDate,
		EndDate:   secondDate,
	})
	require.NoError(t, err)
	require.Len(t, logs, 2)
	require.Equal(t, firstLog.ID, logs[0].ID)
	require.Equal(t, secondLog.ID, logs[1].ID)
}

func TestDailyNutritionRepository_GetOrCreateConcurrentFirstWriteReturnsExistingLog(t *testing.T) {
	pool := nutritionMigrationTestPool(t)
	truncateDailyNutritionTables(t, pool)
	ctx := context.Background()
	q := generated.New(pool)

	userID := seedNutritionMigrationUser(t, pool, "daily nutrition concurrent user")
	userUUID := testUUID(t, userID)
	targetDate := testDate(2026, time.June, 26)

	tx, err := pool.Begin(ctx)
	require.NoError(t, err)
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(context.Background())
		}
	}()

	var firstID pgtype.UUID
	var firstNotes pgtype.Text
	var firstUpdatedAt pgtype.Timestamptz
	err = tx.QueryRow(ctx, `
		INSERT INTO daily_nutrition_logs (user_id, date, notes)
		VALUES ($1, $2, 'first writer')
		RETURNING id, notes, updated_at
	`, userUUID, targetDate).Scan(&firstID, &firstNotes, &firstUpdatedAt)
	require.NoError(t, err)

	type createResult struct {
		id        pgtype.UUID
		notes     pgtype.Text
		updatedAt pgtype.Timestamptz
		err       error
	}

	resultCh := make(chan createResult, 1)
	createCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		log, err := q.CreateDailyNutritionLog(createCtx, generated.CreateDailyNutritionLogParams{
			UserID: userUUID,
			Date:   targetDate,
			Notes:  pgtype.Text{String: "second writer", Valid: true},
		})
		resultCh <- createResult{
			id:        log.ID,
			notes:     log.Notes,
			updatedAt: log.UpdatedAt,
			err:       err,
		}
	}()

	select {
	case result := <-resultCh:
		require.Failf(t, "concurrent create returned before first writer committed", "err=%v id=%v", result.err, result.id)
	case <-time.After(200 * time.Millisecond):
	}

	err = tx.Commit(ctx)
	require.NoError(t, err)
	committed = true

	select {
	case result := <-resultCh:
		require.NoError(t, result.err)
		require.Equal(t, firstID, result.id)
		require.True(t, result.notes.Valid)
		require.Equal(t, firstNotes.String, result.notes.String)
		require.True(t, result.updatedAt.Valid)
		require.True(t, firstUpdatedAt.Time.Equal(result.updatedAt.Time))
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for concurrent daily nutrition log create")
	}
}

func TestDailyNutritionRepository_EntryCRUDPreservesSnapshotsAndOwnership(t *testing.T) {
	pool := nutritionMigrationTestPool(t)
	truncateDailyNutritionTables(t, pool)
	ctx := context.Background()
	q := generated.New(pool)

	userAID := seedNutritionMigrationUser(t, pool, "daily nutrition entry user A")
	userBID := seedNutritionMigrationUser(t, pool, "daily nutrition entry user B")
	productID := seedNutritionMigrationProduct(t, pool, userAID)
	userAUUID := testUUID(t, userAID)
	userBUUID := testUUID(t, userBID)

	log, err := q.CreateDailyNutritionLog(ctx, generated.CreateDailyNutritionLogParams{
		UserID: userAUUID,
		Date:   testDate(2026, time.June, 24),
		Notes:  pgtype.Text{},
	})
	require.NoError(t, err)

	entry, err := q.CreateDailyNutritionEntry(ctx, generated.CreateDailyNutritionEntryParams{
		DailyLogID:  log.ID,
		ProductID:   testUUID(t, productID),
		AmountGrams: 150,
		MealLabel:   pgtype.Text{String: "breakfast", Valid: true},
		Notes:       pgtype.Text{String: "with berries", Valid: true},
		Position:    1,
		UserID:      userAUUID,
	})
	require.NoError(t, err)
	require.Equal(t, "Greek yogurt source", entry.ProductNameSnapshot)
	require.Equal(t, float32(59), entry.CaloriesPer100gSnapshot)
	require.Equal(t, float32(10.2), entry.ProteinPer100gSnapshot)
	require.Equal(t, float32(0.4), entry.FatPer100gSnapshot)
	require.Equal(t, float32(3.6), entry.CarbsPer100gSnapshot)

	_, err = q.UpdateDailyNutritionEntry(ctx, generated.UpdateDailyNutritionEntryParams{
		ID:          entry.ID,
		UserID:      userBUUID,
		DailyLogID:  log.ID,
		AmountGrams: 125,
		MealLabel:   pgtype.Text{String: "snack", Valid: true},
		Notes:       pgtype.Text{String: "wrong user", Valid: true},
		Position:    2,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = q.DeleteDailyNutritionEntry(ctx, generated.DeleteDailyNutritionEntryParams{
		ID:     entry.ID,
		UserID: userBUUID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	updated, err := q.UpdateDailyNutritionEntry(ctx, generated.UpdateDailyNutritionEntryParams{
		ID:          entry.ID,
		UserID:      userAUUID,
		DailyLogID:  log.ID,
		AmountGrams: 125,
		MealLabel:   pgtype.Text{String: "snack", Valid: true},
		Notes:       pgtype.Text{String: "after training", Valid: true},
		Position:    2,
	})
	require.NoError(t, err)
	require.Equal(t, "Greek yogurt source", updated.ProductNameSnapshot)
	require.Equal(t, float32(59), updated.CaloriesPer100gSnapshot)
	require.Equal(t, float32(10.2), updated.ProteinPer100gSnapshot)
	require.Equal(t, float32(0.4), updated.FatPer100gSnapshot)
	require.Equal(t, float32(3.6), updated.CarbsPer100gSnapshot)
	require.Equal(t, float32(125), updated.AmountGrams)
	require.Equal(t, int32(2), updated.Position)

	deleted, err := q.DeleteDailyNutritionEntry(ctx, generated.DeleteDailyNutritionEntryParams{
		ID:     entry.ID,
		UserID: userAUUID,
	})
	require.NoError(t, err)
	require.Equal(t, entry.ID, deleted.ID)
}

func TestDailyNutritionRepository_ListEntriesByLogIsUserScoped(t *testing.T) {
	pool := nutritionMigrationTestPool(t)
	truncateDailyNutritionTables(t, pool)
	ctx := context.Background()
	q := generated.New(pool)

	userAID := seedNutritionMigrationUser(t, pool, "daily nutrition list user A")
	userBID := seedNutritionMigrationUser(t, pool, "daily nutrition list user B")
	productID := seedNutritionMigrationProduct(t, pool, userAID)
	userAUUID := testUUID(t, userAID)
	userBUUID := testUUID(t, userBID)

	log, err := q.CreateDailyNutritionLog(ctx, generated.CreateDailyNutritionLogParams{
		UserID: userAUUID,
		Date:   testDate(2026, time.June, 24),
		Notes:  pgtype.Text{},
	})
	require.NoError(t, err)

	entry, err := q.CreateDailyNutritionEntry(ctx, generated.CreateDailyNutritionEntryParams{
		DailyLogID:  log.ID,
		ProductID:   testUUID(t, productID),
		AmountGrams: 150,
		MealLabel:   pgtype.Text{String: "breakfast", Valid: true},
		Notes:       pgtype.Text{},
		Position:    1,
		UserID:      userAUUID,
	})
	require.NoError(t, err)

	entries, err := q.ListDailyNutritionEntriesByLog(ctx, generated.ListDailyNutritionEntriesByLogParams{
		DailyLogID: log.ID,
		UserID:     userAUUID,
	})
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, entry.ID, entries[0].ID)

	entries, err = q.ListDailyNutritionEntriesByLog(ctx, generated.ListDailyNutritionEntriesByLogParams{
		DailyLogID: log.ID,
		UserID:     userBUUID,
	})
	require.NoError(t, err)
	require.Empty(t, entries)
}

func TestDailyNutritionRepository_RejectsCrossUserProductAttachment(t *testing.T) {
	pool := nutritionMigrationTestPool(t)
	truncateDailyNutritionTables(t, pool)
	ctx := context.Background()
	q := generated.New(pool)

	userAID := seedNutritionMigrationUser(t, pool, "daily nutrition attach user A")
	userBID := seedNutritionMigrationUser(t, pool, "daily nutrition attach user B")
	userAUUID := testUUID(t, userAID)
	userBProductID := seedNutritionMigrationProduct(t, pool, userBID)

	log, err := q.CreateDailyNutritionLog(ctx, generated.CreateDailyNutritionLogParams{
		UserID: userAUUID,
		Date:   testDate(2026, time.June, 24),
		Notes:  pgtype.Text{},
	})
	require.NoError(t, err)

	_, err = q.CreateDailyNutritionEntry(ctx, generated.CreateDailyNutritionEntryParams{
		DailyLogID:  log.ID,
		ProductID:   testUUID(t, userBProductID),
		AmountGrams: 100,
		MealLabel:   pgtype.Text{},
		Notes:       pgtype.Text{},
		Position:    0,
		UserID:      userAUUID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	var entryCount int
	err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM daily_nutrition_entries WHERE daily_log_id = $1`, log.ID).Scan(&entryCount)
	require.NoError(t, err)
	require.Zero(t, entryCount)
}

func nutritionMigrationTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := testinfra.PostgresDSN()
	testinfra.RequireSafePostgresDSN(t, dsn)

	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("INTEGRATION_TESTS not set; skipping nutrition migration integration test")
	}

	if err := postgresrepo.RunMigrations(dsn, zap.NewNop()); err != nil {
		if !testinfra.CoverageGateEnabled() {
			t.Skipf("postgres integration database is unavailable: %v", err)
		}
		require.NoError(t, err)
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	return pool
}

func truncateDailyNutritionTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(context.Background(), `
		TRUNCATE daily_nutrition_entries, daily_nutrition_logs, nutrition_product, atlas_users RESTART IDENTITY CASCADE
	`)
	require.NoError(t, err)
}

func seedNutritionMigrationUser(t *testing.T, pool *pgxpool.Pool, displayName string) string {
	t.Helper()

	var userID string
	err := pool.QueryRow(context.Background(), `
		INSERT INTO atlas_users (display_name)
		VALUES ($1)
		RETURNING id
	`, displayName).Scan(&userID)
	require.NoError(t, err)
	return userID
}

func seedNutritionMigrationProduct(t *testing.T, pool *pgxpool.Pool, userID string) string {
	t.Helper()

	var productID string
	err := pool.QueryRow(context.Background(), `
		INSERT INTO nutrition_product (
			user_id,
			name,
			calories_per_100g,
			protein_per_100g,
			fat_per_100g,
			carbs_per_100g,
			notes
		)
		VALUES ($1, 'Greek yogurt source', 59, 10.2, 0.4, 3.6, 'seed product')
		RETURNING id
	`, userID).Scan(&productID)
	require.NoError(t, err)
	return productID
}

func requirePgConstraint(t *testing.T, err error, constraintName string) {
	t.Helper()
	require.Error(t, err)

	var pgErr *pgconn.PgError
	require.True(t, errors.As(err, &pgErr), "expected pg error, got %T: %v", err, err)
	require.Equal(t, constraintName, pgErr.ConstraintName)
}

func testUUID(t *testing.T, value string) pgtype.UUID {
	t.Helper()

	var uuid pgtype.UUID
	require.NoError(t, uuid.Scan(value))
	return uuid
}

func testDate(year int, month time.Month, day int) pgtype.Date {
	return pgtype.Date{Time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC), Valid: true}
}
