// FILE: apps/api/internal/atlas/repository/postgres/user_profile_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement UserProfileRepository for WAVE-07 AI Export user profile management using sqlc-generated queries.
//   SCOPE: FindByUserID, Upsert, Create - all user-scoped.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   UserProfileRepository - Interface for user profile data access.
//   NewUserProfileRepository - Creates a new UserProfileRepository.
//   FindByUserID - Gets a user profile by user ID.
//   Upsert - Creates or updates a user profile (partial update on conflict).
//   Create - Creates a new user profile.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added user profile repository for WAVE-07.
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

type UserProfileRepository interface {
	FindByUserID(ctx context.Context, userID string) (*models.UserProfileRecord, error)
	Upsert(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfileRecord, error)
	Create(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfileRecord, error)
}

type userProfileRepository struct {
	q *generated.Queries
}

func NewUserProfileRepository(pool *pgxpool.Pool) UserProfileRepository {
	return &userProfileRepository{q: generated.New(pool)}
}

func (r *userProfileRepository) FindByUserID(ctx context.Context, userID string) (*models.UserProfileRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("user_profile_repo.FindByUserID: %w", err)
	}

	row, err := r.q.GetUserProfileByUserID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("user_profile_repo.FindByUserID: %w", err)
	}

	return userProfileRecordFromRow(row), nil
}

func (r *userProfileRepository) Upsert(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfileRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("user_profile_repo.Upsert: %w", err)
	}

	row, err := r.q.UpsertUserProfile(ctx, generated.UpsertUserProfileParams{
		UserID:                    uid,
		Goal:                      nullableText(input.Goal),
		Height:                    nullableFloat4(input.Height),
		BirthDate:                 nullableDate(input.BirthDate),
		TrainingExperience:        nullableText(input.TrainingExperience),
		CurrentTrainingSplit:      nullableText(input.CurrentTrainingSplit),
		PreferredProgressionStyle: nullableText(input.PreferredProgressionStyle),
		NutritionStrategy:         nullableText(input.NutritionStrategy),
		PersistentAiContext:       nullableText(input.PersistentAiContext),
	})
	if err != nil {
		return nil, fmt.Errorf("user_profile_repo.Upsert: %w", err)
	}

	return userProfileRecordFromRow(row), nil
}

func (r *userProfileRepository) Create(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfileRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("user_profile_repo.Create: %w", err)
	}

	row, err := r.q.CreateUserProfile(ctx, generated.CreateUserProfileParams{
		UserID:                    uid,
		Goal:                      nullableText(input.Goal),
		Height:                    nullableFloat4(input.Height),
		BirthDate:                 nullableDate(input.BirthDate),
		TrainingExperience:        nullableText(input.TrainingExperience),
		CurrentTrainingSplit:      nullableText(input.CurrentTrainingSplit),
		PreferredProgressionStyle: nullableText(input.PreferredProgressionStyle),
		NutritionStrategy:         nullableText(input.NutritionStrategy),
		PersistentAiContext:       nullableText(input.PersistentAiContext),
	})
	if err != nil {
		return nil, fmt.Errorf("user_profile_repo.Create: %w", err)
	}

	return userProfileRecordFromRow(row), nil
}

func userProfileRecordFromRow(row generated.UserProfile) *models.UserProfileRecord {
	var birthDate *models.Date
	if row.BirthDate.Valid {
		d := dateFromPGDate(row.BirthDate)
		birthDate = &d
	}

	return &models.UserProfileRecord{
		ID:                       row.ID.String(),
		UserID:                   row.UserID.String(),
		Goal:                     textPtr(row.Goal),
		Height:                   float4Ptr(row.Height),
		BirthDate:                birthDate,
		TrainingExperience:       textPtr(row.TrainingExperience),
		CurrentTrainingSplit:     textPtr(row.CurrentTrainingSplit),
		PreferredProgressionStyle: textPtr(row.PreferredProgressionStyle),
		NutritionStrategy:        textPtr(row.NutritionStrategy),
		PersistentAiContext:      textPtr(row.PersistentAiContext),
		CreatedAt:                formatTimestamp(row.CreatedAt),
		UpdatedAt:                formatTimestamp(row.UpdatedAt),
	}
}

func nullableDate(d *models.Date) pgtype.Date {
	if d == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: d.Time(), Valid: true}
}