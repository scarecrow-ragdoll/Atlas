// FILE: apps/api/internal/atlas/repository/postgres/settings_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement the SettingsRepository interface for the Atlas fitness tracker using sqlc-generated queries.
//   SCOPE: FindByUserID, UpsertSettings, UpdatePinState; avoids generic Upsert and prevents accidental PIN hash overwrites.
//   DEPENDS: apps/api/internal/repository/postgres/generated (sqlc), apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Atlas settings repository for WAVE-01.
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

type SettingsRepository interface {
	FindByUserID(ctx context.Context, userID string) (*models.SettingsRecord, error)
	UpsertSettings(ctx context.Context, userID string, input models.SettingsInput) (*models.SettingsRecord, error)
	UpdatePinState(ctx context.Context, userID string, pinEnabled bool, pinHash *string) error
}

type settingsRepository struct {
	q *generated.Queries
}

func NewSettingsRepository(pool *pgxpool.Pool) SettingsRepository {
	return &settingsRepository{q: generated.New(pool)}
}

func (r *settingsRepository) FindByUserID(ctx context.Context, userID string) (*models.SettingsRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("settings_repo.FindByUserID: %w", err)
	}

	row, err := r.q.GetAtlasSettingsByUserID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("settings_repo.FindByUserID: %w", err)
	}

	return recordFromRow(row), nil
}

func (r *settingsRepository) UpsertSettings(ctx context.Context, userID string, input models.SettingsInput) (*models.SettingsRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("settings_repo.UpsertSettings: %w", err)
	}

	units := ""
	if input.Units != nil {
		units = *input.Units
	}
	exportWeeks := int32(0)
	if input.DefaultAiExportWeeks != nil {
		exportWeeks = int32(*input.DefaultAiExportWeeks)
	}

	row, err := r.q.UpsertAtlasSettings(ctx, generated.UpsertAtlasSettingsParams{
		UserID:               uid,
		PinEnabled:           false,
		Units:                units,
		DefaultAiExportWeeks: exportWeeks,
	})
	if err != nil {
		return nil, fmt.Errorf("settings_repo.UpsertSettings: %w", err)
	}

	return recordFromRow(row), nil
}

func (r *settingsRepository) UpdatePinState(ctx context.Context, userID string, pinEnabled bool, pinHash *string) error {
	uid, err := uuidFromString(userID)
	if err != nil {
		return fmt.Errorf("settings_repo.UpdatePinState: %w", err)
	}

	_, err = r.q.UpdateAtlasPinState(ctx, generated.UpdateAtlasPinStateParams{
		UserID:     uid,
		PinEnabled: pinEnabled,
		PinHash:    nullableText(pinHash),
	})
	if err != nil {
		return fmt.Errorf("settings_repo.UpdatePinState: %w", err)
	}

	return nil
}

func recordFromRow(row generated.AtlasSetting) *models.SettingsRecord {
	var pinHash *string
	if row.PinHash.Valid {
		pinHash = &row.PinHash.String
	}

	return &models.SettingsRecord{
		ID:                   row.ID.String(),
		UserID:               row.UserID.String(),
		PinEnabled:           row.PinEnabled,
		PinHash:              pinHash,
		Units:                row.Units,
		DefaultAiExportWeeks: row.DefaultAiExportWeeks,
		CreatedAt:            formatTimestamp(row.CreatedAt),
		UpdatedAt:            formatTimestamp(row.UpdatedAt),
	}
}

func uuidFromString(value string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(value); err != nil {
		return pgtype.UUID{}, err
	}
	return uuid, nil
}

func nullableText(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *value, Valid: true}
}

func formatTimestamp(value pgtype.Timestamptz) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format("2006-01-02T15:04:05.999999999Z07:00")
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