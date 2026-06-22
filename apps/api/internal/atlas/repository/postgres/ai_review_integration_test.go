// FILE: apps/api/internal/atlas/repository/postgres/ai_review_integration_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Integration tests for WAVE-08 AiReview repository and migration: verifies CRUD operations and migration apply/rollback. Skipped when DB unavailable.
//   SCOPE: TEST-W08-013 (AiReview repo CRUD), TEST-W08-014 (migration 00093 applies cleanly). Requires test Postgres on port 17501.
//   DEPENDS: internal/testinfra, internal/repository/postgres (RunMigrations), and the AiReview repository implementation.
//   LINKS: V-M-API / WAVE-08.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added AiReview integration tests for WAVE-08.
// END_CHANGE_SUMMARY

package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	postgresrepo "monorepo-template/apps/api/internal/repository/postgres"
	"monorepo-template/apps/api/internal/testinfra"
)

func TestAiReviewRepo_CRUD(t *testing.T) {
	pool := aiReviewTestPool(t)

	userID := seedTestUser(t, pool)

	repo := atlasPostgres.NewAiReviewRepository(pool)

	dateStart := models.MustDate("2026-06-01")
	dateEnd := models.MustDate("2026-06-30")
	userNotes := "User provided notes"
	plannedActions := "Increase protein, sleep 8h"

	// Create
	created, err := repo.Create(context.Background(), userID, "2026-06-01", "2026-06-30", "AI analysis result", &userNotes, &plannedActions)
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, userID, created.UserID)
	assert.Equal(t, dateStart, created.DateRangeStart)
	assert.Equal(t, dateEnd, created.DateRangeEnd)
	assert.Equal(t, "AI analysis result", created.AiResponseText)
	assert.NotNil(t, created.UserNotes)
	assert.Equal(t, "User provided notes", *created.UserNotes)
	assert.NotNil(t, created.PlannedActions)
	assert.Equal(t, "Increase protein, sleep 8h", *created.PlannedActions)

	reviewID := created.ID

	// GetByID
	found, err := repo.GetByID(context.Background(), userID, reviewID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, reviewID, found.ID)
	assert.Equal(t, "AI analysis result", found.AiResponseText)

	// GetByID with wrong user returns nil
	otherUserID := "00000000-0000-0000-0000-000000000001"
	notFound, err := repo.GetByID(context.Background(), otherUserID, reviewID)
	require.NoError(t, err)
	assert.Nil(t, notFound)

	// ListByUserID
	listAll, err := repo.ListByUserID(context.Background(), userID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(listAll), 1)

	// ListByUserIDAndDateRange
	filtered, err := repo.ListByUserIDAndDateRange(context.Background(), userID, "2026-06-01", "2026-06-30")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(filtered), 1)

	filteredEmpty, err := repo.ListByUserIDAndDateRange(context.Background(), userID, "2025-01-01", "2025-01-31")
	require.NoError(t, err)
	assert.Empty(t, filteredEmpty)

	// Update (partial)
	updatedNotes := "Updated notes"
	updated, err := repo.Update(context.Background(), userID, reviewID, models.UpdateAiReviewInput{
		UserNotes: &updatedNotes,
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Updated notes", *updated.UserNotes)
	assert.Equal(t, "AI analysis result", updated.AiResponseText)

	// Update with wrong user returns nil
	notFoundUpdate, err := repo.Update(context.Background(), otherUserID, reviewID, models.UpdateAiReviewInput{
		UserNotes: &updatedNotes,
	})
	require.NoError(t, err)
	assert.Nil(t, notFoundUpdate)

	// Delete
	deleted, err := repo.Delete(context.Background(), userID, reviewID)
	require.NoError(t, err)
	require.NotNil(t, deleted)
	assert.Equal(t, reviewID, deleted.ID)

	// Delete again returns nil
	gone, err := repo.Delete(context.Background(), userID, reviewID)
	require.NoError(t, err)
	assert.Nil(t, gone)
}

func TestWave08Migration_ApplyCleanly(t *testing.T) {
	dsn := testinfra.PostgresDSN()
	testinfra.RequireSafePostgresDSN(t, dsn)

	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("INTEGRATION_TESTS not set; skipping WAVE-08 migration smoke test")
	}

	err := postgresrepo.RunMigrations(dsn, zap.NewNop())
	require.NoError(t, err, "WAVE-08 migration should apply without error")
}

// --- helpers ---

func aiReviewTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := testinfra.PostgresDSN()
	testinfra.RequireSafePostgresDSN(t, dsn)

	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("INTEGRATION_TESTS not set; skipping AiReview repo integration test")
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

	_, err = pool.Exec(context.Background(), `TRUNCATE ai_reviews RESTART IDENTITY CASCADE`)
	require.NoError(t, err)

	return pool
}

func seedTestUser(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	var userID string
	err := pool.QueryRow(context.Background(), `
		INSERT INTO atlas_users (email, name)
		VALUES ('test-ai-review@example.com', 'Test AI Review User')
		ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`).Scan(&userID)
	require.NoError(t, err)
	return userID
}