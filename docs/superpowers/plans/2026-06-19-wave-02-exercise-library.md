# WAVE-02: Exercise Library Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement full Exercise CRUD with working weight, media management, and soft archive/restore for the Atlas fitness tracker.

**Architecture:** 8 sequential slices: migrations → sqlc queries → repository → service → GraphQL schema → resolvers → media handler extension → main wiring. Follows WAVE-01 Approach A (Atlas-specific dirs). Media uses WAVE-01 scaffold routes (POST/GET/DELETE /api/v1/media/*) with purpose/entity routing. Exercise CRUD is GraphQL-only via /graphql/atlas. No REST exercise endpoints. Soft archive (isActive=false) replaces hard delete. Archive does NOT cascade media.

**Tech Stack:** Go 1.22+, goose migrations, sqlc, gqlgen, chi router, pgx v5, PostgreSQL 16

**Test commands:**
- Unit/integration: `bunx nx run api:test -- --run '<pattern>'`
- Codegen: `bunx nx run api:codegen && bunx nx run api:codegen:atlas`
- Lint: `bunx nx run api:lint`

---

### Task 1: DB migrations — 00081_exercises.sql

**Files:**
- Create: `apps/api/internal/repository/postgres/migrations/00081_exercises.sql`
- Test: `bunx nx run api:test -- --run '(?i)migration'`

- [ ] **Step 1: Create 00081_exercises.sql migration**

```sql
-- FILE: apps/api/internal/repository/postgres/migrations/00081_exercises.sql
-- +goose Up
CREATE TABLE exercises (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES atlas_users(id),
    name            TEXT NOT NULL,
    muscle_groups   TEXT[] DEFAULT '{}',
    description     TEXT DEFAULT '',
    personal_notes  TEXT DEFAULT '',
    working_weight  REAL,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (working_weight IS NULL OR working_weight > 0)
);

CREATE INDEX idx_exercises_user_active ON exercises (user_id, is_active);
CREATE INDEX idx_exercises_user_name ON exercises (user_id, name);
CREATE INDEX idx_exercises_user_created_at ON exercises (user_id, created_at);

-- +goose Down
DROP TABLE IF EXISTS exercises;
```

- [ ] **Step 2: Create 00082_exercise_media.sql migration**

```sql
-- FILE: apps/api/internal/repository/postgres/migrations/00082_exercise_media.sql
-- +goose Up
CREATE TABLE exercise_media (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES atlas_users(id),
    exercise_id  UUID NOT NULL REFERENCES exercises(id) ON DELETE NO ACTION,
    file_name    TEXT NOT NULL,
    file_path    TEXT NOT NULL,
    mime_type    TEXT NOT NULL,
    file_size    BIGINT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_exercise_media_user_exercise ON exercise_media (user_id, exercise_id);

-- +goose Down
DROP TABLE IF EXISTS exercise_media;
```

- [ ] **Step 3: Run migration smoke test**

Run: `bunx nx run api:test -- --run '(?i)migration'`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/repository/postgres/migrations/00081_exercises.sql apps/api/internal/repository/postgres/migrations/00082_exercise_media.sql
git commit -m "feat(wave-02): add exercises and exercise_media migrations"
```

---

### Task 2: sqlc queries — exercises.sql

**Files:**
- Create: `apps/api/internal/repository/postgres/queries/exercises.sql`
- Generated: `apps/api/internal/repository/postgres/generated/` (auto by sqlc codegen)

- [ ] **Step 1: Create exercises.sql queries**

```sql
-- name: CreateExercise :one
INSERT INTO exercises (user_id, name, muscle_groups, description, personal_notes, working_weight, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at;

-- name: GetExerciseByID :one
SELECT id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at
FROM exercises
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListExercises :many
SELECT id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at
FROM exercises
WHERE user_id = $1 AND ($2::bool OR is_active = true)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountExercises :one
SELECT COUNT(*)
FROM exercises
WHERE user_id = $1 AND ($2::bool OR is_active = true);

-- name: UpdateExercise :one
UPDATE exercises
SET name = COALESCE(NULLIF($3::text, ''), name),
    muscle_groups = CASE WHEN $4::text[] IS NOT NULL THEN $4 ELSE muscle_groups END,
    description = COALESCE(NULLIF($5::text, ''), description),
    personal_notes = COALESCE(NULLIF($6::text, ''), personal_notes),
    working_weight = COALESCE($7::real, working_weight),
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at;

-- name: ArchiveExercise :one
UPDATE exercises
SET is_active = false, updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at;

-- name: RestoreExercise :one
UPDATE exercises
SET is_active = true, updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at;

-- name: AllExercises :many
SELECT id, user_id, name, muscle_groups, description, personal_notes, working_weight, is_active, created_at, updated_at
FROM exercises
WHERE user_id = $1 AND ($2::bool OR is_active = true)
ORDER BY name ASC;

-- name: CreateExerciseMedia :one
INSERT INTO exercise_media (user_id, exercise_id, file_name, file_path, mime_type, file_size)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, exercise_id, file_name, file_path, mime_type, file_size, created_at;

-- name: GetExerciseMediaByID :one
SELECT id, user_id, exercise_id, file_name, file_path, mime_type, file_size, created_at
FROM exercise_media
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListExerciseMediaByExerciseID :many
SELECT id, user_id, exercise_id, file_name, file_path, mime_type, file_size, created_at
FROM exercise_media
WHERE exercise_id = $1 AND user_id = $2
ORDER BY created_at ASC;

-- name: DeleteExerciseMedia :one
DELETE FROM exercise_media
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, exercise_id, file_name, file_path, mime_type, file_size, created_at;
```

- [ ] **Step 2: Run sqlc codegen**

Run: `cd apps/api && go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 generate`
Expected: generates `internal/repository/postgres/generated/exercises.sql.go` without errors

- [ ] **Step 3: Verify generated files exist**

Run: `ls apps/api/internal/repository/postgres/generated/`
Expected: shows `exercises.sql.go` among existing files

- [ ] **Step 4: Verify generated code compiles**

Run: `cd apps/api && go vet ./...`
Expected: no errors

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/repository/postgres/queries/exercises.sql
git add apps/api/internal/repository/postgres/generated/exercises.sql.go
git commit -m "feat(wave-02): add sqlc queries for exercise CRUD"
```

---

### Task 3: Exercise repository — exercise_repo.go

**Files:**
- Create: `apps/api/internal/atlas/repository/postgres/exercise_repo.go`
- Test: `apps/api/internal/atlas/repository/postgres/exercise_repo_test.go` (future task)

- [ ] **Step 1: Create ExerciseRepo interface + implementation**

```go
// FILE: apps/api/internal/atlas/repository/postgres/exercise_repo.go
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
	Create(ctx context.Context, userID string, name string, muscleGroups []string, description, personalNotes string, workingWeight *float64, isActive bool) (*models.ExerciseRecord, error)
	GetByID(ctx context.Context, userID, exerciseID string) (*models.ExerciseRecord, error)
	List(ctx context.Context, userID string, includeInactive bool, limit, offset int32) ([]*models.ExerciseRecord, int32, error)
	Update(ctx context.Context, userID, exerciseID string, name *string, muscleGroups []string, description, personalNotes *string, workingWeight *float64) (*models.ExerciseRecord, error)
	Archive(ctx context.Context, userID, exerciseID string) (*models.ExerciseRecord, error)
	Restore(ctx context.Context, userID, exerciseID string) (*models.ExerciseRecord, error)
	AllExercises(ctx context.Context, userID string, includeInactive bool) ([]*models.ExerciseRecord, error)
	CreateMedia(ctx context.Context, userID, exerciseID, fileName, filePath, mimeType string, fileSize int64) (*models.ExerciseMediaRecord, error)
	GetMediaByID(ctx context.Context, userID, mediaID string) (*models.ExerciseMediaRecord, error)
	ListMediaByExerciseID(ctx context.Context, userID, exerciseID string) ([]*models.ExerciseMediaRecord, error)
	DeleteMedia(ctx context.Context, userID, mediaID string) (*models.ExerciseMediaRecord, error)
}

type exerciseRepository struct {
	q *generated.Queries
}

func NewExerciseRepository(pool *pgxpool.Pool) ExerciseRepository {
	return &exerciseRepository{q: generated.New(pool)}
}

func (r *exerciseRepository) Create(ctx context.Context, userID string, name string, muscleGroups []string, description, personalNotes string, workingWeight *float64, isActive bool) (*models.ExerciseRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.Create: %w", err)
	}
	row, err := r.q.CreateExercise(ctx, generated.CreateExerciseParams{
		UserID:        uid,
		Name:          name,
		MuscleGroups:  muscleGroups,
		Description:   description,
		PersonalNotes: personalNotes,
		WorkingWeight: nullableFloat8(workingWeight),
		IsActive:      isActive,
	})
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.Create: %w", err)
	}
	return exerciseRecordFromRow(row), nil
}

func (r *exerciseRepository) GetByID(ctx context.Context, userID, exerciseID string) (*models.ExerciseRecord, error) {
	uid, eid, err := twoUUIDs(userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.GetByID: %w", err)
	}
	row, err := r.q.GetExerciseByID(ctx, generated.GetExerciseByIDParams{
		ID:     eid,
		UserID: uid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exercise_repo.GetByID: %w", err)
	}
	return exerciseRecordFromRow(row), nil
}

func (r *exerciseRepository) List(ctx context.Context, userID string, includeInactive bool, limit, offset int32) ([]*models.ExerciseRecord, int32, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("exercise_repo.List: %w", err)
	}
	rows, err := r.q.ListExercises(ctx, generated.ListExercisesParams{
		UserID:          uid,
		Includeinactive: includeInactive,
		Limit:           limit,
		Offset:          offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("exercise_repo.List: %w", err)
	}
	count, err := r.q.CountExercises(ctx, generated.CountExercisesParams{
		UserID:          uid,
		Includeinactive: includeInactive,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("exercise_repo.List: %w", err)
	}
	records := make([]*models.ExerciseRecord, len(rows))
	for i, row := range rows {
		records[i] = exerciseRecordFromRow(row)
	}
	return records, count, nil
}

func (r *exerciseRepository) Update(ctx context.Context, userID, exerciseID string, name *string, muscleGroups []string, description, personalNotes *string, workingWeight *float64) (*models.ExerciseRecord, error) {
	uid, eid, err := twoUUIDs(userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.Update: %w", err)
	}
	nameVal := ""
	if name != nil {
		nameVal = *name
	}
	descVal := ""
	if description != nil {
		descVal = *description
	}
	notesVal := ""
	if personalNotes != nil {
		notesVal = *personalNotes
	}
	row, err := r.q.UpdateExercise(ctx, generated.UpdateExerciseParams{
		ID:            eid,
		UserID:        uid,
		Name:          nameVal,
		MuscleGroups:  muscleGroups,
		Description:   descVal,
		PersonalNotes: notesVal,
		WorkingWeight: nullableFloat8(workingWeight),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exercise_repo.Update: %w", err)
	}
	return exerciseRecordFromRow(row), nil
}

func (r *exerciseRepository) Archive(ctx context.Context, userID, exerciseID string) (*models.ExerciseRecord, error) {
	uid, eid, err := twoUUIDs(userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.Archive: %w", err)
	}
	row, err := r.q.ArchiveExercise(ctx, generated.ArchiveExerciseParams{
		ID:     eid,
		UserID: uid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exercise_repo.Archive: %w", err)
	}
	return exerciseRecordFromRow(row), nil
}

func (r *exerciseRepository) Restore(ctx context.Context, userID, exerciseID string) (*models.ExerciseRecord, error) {
	uid, eid, err := twoUUIDs(userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.Restore: %w", err)
	}
	row, err := r.q.RestoreExercise(ctx, generated.RestoreExerciseParams{
		ID:     eid,
		UserID: uid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exercise_repo.Restore: %w", err)
	}
	return exerciseRecordFromRow(row), nil
}

func (r *exerciseRepository) AllExercises(ctx context.Context, userID string, includeInactive bool) ([]*models.ExerciseRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.AllExercises: %w", err)
	}
	rows, err := r.q.AllExercises(ctx, generated.AllExercisesParams{
		UserID:          uid,
		Includeinactive: includeInactive,
	})
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.AllExercises: %w", err)
	}
	records := make([]*models.ExerciseRecord, len(rows))
	for i, row := range rows {
		records[i] = exerciseRecordFromRow(row)
	}
	return records, nil
}

func (r *exerciseRepository) CreateMedia(ctx context.Context, userID, exerciseID, fileName, filePath, mimeType string, fileSize int64) (*models.ExerciseMediaRecord, error) {
	uid, eid, err := twoUUIDs(userID, exerciseID)
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
	return exerciseMediaRecordFromRow(row), nil
}

func (r *exerciseRepository) GetMediaByID(ctx context.Context, userID, mediaID string) (*models.ExerciseMediaRecord, error) {
	uid, mid, err := twoUUIDs(userID, mediaID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.GetMediaByID: %w", err)
	}
	row, err := r.q.GetExerciseMediaByID(ctx, generated.GetExerciseMediaByIDParams{
		ID:     mid,
		UserID: uid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exercise_repo.GetMediaByID: %w", err)
	}
	return exerciseMediaRecordFromRow(row), nil
}

func (r *exerciseRepository) ListMediaByExerciseID(ctx context.Context, userID, exerciseID string) ([]*models.ExerciseMediaRecord, error) {
	uid, eid, err := twoUUIDs(userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.ListMediaByExerciseID: %w", err)
	}
	rows, err := r.q.ListExerciseMediaByExerciseID(ctx, generated.ListExerciseMediaByExerciseIDParams{
		ExerciseID: eid,
		UserID:     uid,
	})
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.ListMediaByExerciseID: %w", err)
	}
	records := make([]*models.ExerciseMediaRecord, len(rows))
	for i, row := range rows {
		records[i] = exerciseMediaRecordFromRow(row)
	}
	return records, nil
}

func (r *exerciseRepository) DeleteMedia(ctx context.Context, userID, mediaID string) (*models.ExerciseMediaRecord, error) {
	uid, mid, err := twoUUIDs(userID, mediaID)
	if err != nil {
		return nil, fmt.Errorf("exercise_repo.DeleteMedia: %w", err)
	}
	row, err := r.q.DeleteExerciseMedia(ctx, generated.DeleteExerciseMediaParams{
		ID:     mid,
		UserID: uid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exercise_repo.DeleteMedia: %w", err)
	}
	return exerciseMediaRecordFromRow(row), nil
}

func exerciseRecordFromRow(row generated.Exercise) *models.ExerciseRecord {
	var ww *float64
	if row.WorkingWeight.Valid {
		ww = &row.WorkingWeight.Float64
	}
	return &models.ExerciseRecord{
		ID:            row.ID.String(),
		UserID:        row.UserID.String(),
		Name:          row.Name,
		MuscleGroups:  row.MuscleGroups,
		Description:   row.Description,
		PersonalNotes: row.PersonalNotes,
		WorkingWeight: ww,
		IsActive:      row.IsActive,
		CreatedAt:     formatTimestamp(row.CreatedAt),
		UpdatedAt:     formatTimestamp(row.UpdatedAt),
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

func nullableFloat8(v *float64) pgtype.Float8 {
	if v == nil {
		return pgtype.Float8{}
	}
	return pgtype.Float8{Float64: *v, Valid: true}
}

func twoUUIDs(a, b string) (pgtype.UUID, pgtype.UUID, error) {
	ua, err := uuidFromString(a)
	if err != nil {
		return pgtype.UUID{}, pgtype.UUID{}, err
	}
	ub, err := uuidFromString(b)
	if err != nil {
		return pgtype.UUID{}, pgtype.UUID{}, err
	}
	return ua, ub, nil
}
```

Note: `uuidFromString`, `formatTimestamp`, `nullableText` are already in `settings_repo.go`. They must be in the same package (`postgres` in `internal/atlas/repository/postgres`). If they are in `settings_repo.go`, they are accessible. If `settings_repo.go` uses a different package path, extract them to a shared `internal/atlas/repository/postgres/helpers.go` file.

- [ ] **Step 2: Verify compilation**

Run: `cd apps/api && go vet ./internal/atlas/repository/postgres/...`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add apps/api/internal/atlas/repository/postgres/exercise_repo.go
git commit -m "feat(wave-02): add ExerciseRepository with sqlc-based CRUD"
```

---

### Task 4: Exercise models — exercise.go

**Files:**
- Create: `apps/api/internal/atlas/models/exercise.go`

- [ ] **Step 1: Create models for Exercise and ExerciseMedia**

```go
// FILE: apps/api/internal/atlas/models/exercise.go
package models

// ExerciseRecord is the internal DB model.
type ExerciseRecord struct {
	ID            string
	UserID        string
	Name          string
	MuscleGroups  []string
	Description   string
	PersonalNotes string
	WorkingWeight *float64
	IsActive      bool
	CreatedAt     string
	UpdatedAt     string
}

// Exercise is the public GraphQL-facing model.
type Exercise struct {
	ID            string              `json:"id"`
	UserID        string              `json:"userId"`
	Name          string              `json:"name"`
	MuscleGroups  []string            `json:"muscleGroups"`
	Description   string              `json:"description"`
	PersonalNotes string              `json:"personalNotes"`
	WorkingWeight *float64            `json:"workingWeight"`
	IsActive      bool                `json:"isActive"`
	Media         []ExerciseMedia     `json:"media"`
	CreatedAt     string              `json:"createdAt"`
	UpdatedAt     string              `json:"updatedAt"`
}

type CreateExerciseInput struct {
	Name          string    `json:"name"`
	MuscleGroups  *[]string `json:"muscleGroups"`
	Description   *string   `json:"description"`
	PersonalNotes *string   `json:"personalNotes"`
	WorkingWeight *float64  `json:"workingWeight"`
}

type UpdateExerciseInput struct {
	Name          *string   `json:"name"`
	MuscleGroups  *[]string `json:"muscleGroups"`
	Description   *string   `json:"description"`
	PersonalNotes *string   `json:"personalNotes"`
	WorkingWeight *float64  `json:"workingWeight"`
}

// ExerciseMediaRecord is the internal DB model.
type ExerciseMediaRecord struct {
	ID         string
	UserID     string
	ExerciseID string
	FileName   string
	FilePath   string
	MimeType   string
	FileSize   int64
	CreatedAt  string
}

// ExerciseMedia is the public GraphQL-facing model.
type ExerciseMedia struct {
	ID         string `json:"id"`
	UserID     string `json:"userId"`
	ExerciseID string `json:"exerciseId"`
	FileName   string `json:"fileName"`
	MimeType   string `json:"mimeType"`
	FileSize   int64  `json:"fileSize"`
	CreatedAt  string `json:"createdAt"`
}

// ExerciseResult follows the WAVE-01 Result wrapper pattern.
type ExerciseResult struct {
	Exercise *Exercise       `json:"exercise"`
	Error    *ExerciseError  `json:"error"`
}

type ArchiveResult struct {
	Exercise *Exercise       `json:"exercise"`
	Error    *ExerciseError  `json:"error"`
}

type ExerciseError struct {
	Message string             `json:"message"`
	Code    ExerciseErrorCode  `json:"code"`
}

type ExerciseErrorCode string

const (
	ExerciseErrorValidation   ExerciseErrorCode = "VALIDATION_ERROR"
	ExerciseErrorNotFound     ExerciseErrorCode = "NOT_FOUND"
	ExerciseErrorUnauthorized ExerciseErrorCode = "UNAUTHORIZED"
	ExerciseErrorInternal     ExerciseErrorCode = "INTERNAL_ERROR"
)
```

- [ ] **Step 2: Commit**

```bash
git add apps/api/internal/atlas/models/exercise.go
git commit -m "feat(wave-02): add Exercise and ExerciseMedia models"
```

---

### Task 5: Exercise service — exercise.go

**Files:**
- Create: `apps/api/internal/atlas/service/exercise.go`
- Test: `apps/api/internal/atlas/service/exercise_test.go` (future task)

- [ ] **Step 1: Create ExerciseService with validation**

```go
// FILE: apps/api/internal/atlas/service/exercise.go
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasRepo "monorepo-template/apps/api/internal/atlas/repository/postgres"
)

var (
	ErrExerciseNotFound     = errors.New("exercise not found")
	ErrExerciseNameRequired = errors.New("exercise name is required")
	ErrExerciseWeightZero   = errors.New("working weight must be greater than 0")
)

type ExerciseService interface {
	Create(ctx context.Context, userID string, input *models.CreateExerciseInput) (*models.ExerciseRecord, error)
	GetByID(ctx context.Context, userID, exerciseID string) (*models.ExerciseRecord, error)
	List(ctx context.Context, userID string, includeInactive bool, first, after int32) ([]*models.ExerciseRecord, int32, error)
	Update(ctx context.Context, userID, exerciseID string, input *models.UpdateExerciseInput) (*models.ExerciseRecord, error)
	Archive(ctx context.Context, userID, exerciseID string) (*models.ExerciseRecord, error)
	Restore(ctx context.Context, userID, exerciseID string) (*models.ExerciseRecord, error)
	AllExercises(ctx context.Context, userID string, includeInactive bool) ([]*models.ExerciseRecord, error)
	CreateMedia(ctx context.Context, userID, exerciseID, fileName, filePath, mimeType string, fileSize int64) (*models.ExerciseMediaRecord, error)
	GetMediaByID(ctx context.Context, userID, mediaID string) (*models.ExerciseMediaRecord, error)
	DeleteMedia(ctx context.Context, userID, mediaID string) (*models.ExerciseMediaRecord, error)
}

type exerciseService struct {
	repo atlasRepo.ExerciseRepository
}

func NewExerciseService(repo atlasRepo.ExerciseRepository) ExerciseService {
	return &exerciseService{repo: repo}
}

func (s *exerciseService) Create(ctx context.Context, userID string, input *models.CreateExerciseInput) (*models.ExerciseRecord, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrExerciseNameRequired
	}
	if input.WorkingWeight != nil && *input.WorkingWeight <= 0 {
		return nil, ErrExerciseWeightZero
	}
	muscleGroups := []string{}
	if input.MuscleGroups != nil {
		muscleGroups = *input.MuscleGroups
	}
	desc := ""
	if input.Description != nil {
		desc = *input.Description
	}
	notes := ""
	if input.PersonalNotes != nil {
		notes = *input.PersonalNotes
	}
	return s.repo.Create(ctx, userID, name, muscleGroups, desc, notes, input.WorkingWeight, true)
}

func (s *exerciseService) GetByID(ctx context.Context, userID, exerciseID string) (*models.ExerciseRecord, error) {
	record, err := s.repo.GetByID(ctx, userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrExerciseNotFound
	}
	return record, nil
}

func (s *exerciseService) List(ctx context.Context, userID string, includeInactive bool, first, after int32) ([]*models.ExerciseRecord, int32, error) {
	if first <= 0 {
		first = 20
	}
	return s.repo.List(ctx, userID, includeInactive, first, after)
}

func (s *exerciseService) Update(ctx context.Context, userID, exerciseID string, input *models.UpdateExerciseInput) (*models.ExerciseRecord, error) {
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, ErrExerciseNameRequired
		}
		input.Name = &name
	}
	if input.WorkingWeight != nil && *input.WorkingWeight <= 0 {
		return nil, ErrExerciseWeightZero
	}
	var muscleGroups []string // nil when not provided — sqlc binds as NULL, preserves existing
	if input.MuscleGroups != nil {
		muscleGroups = *input.MuscleGroups
	}
	record, err := s.repo.Update(ctx, userID, exerciseID, input.Name, muscleGroups, input.Description, input.PersonalNotes, input.WorkingWeight)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrExerciseNotFound
	}
	return record, nil
}

func (s *exerciseService) Archive(ctx context.Context, userID, exerciseID string) (*models.ExerciseRecord, error) {
	record, err := s.repo.Archive(ctx, userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.Archive: %w", err)
	}
	if record == nil {
		return nil, ErrExerciseNotFound
	}
	return record, nil
}

func (s *exerciseService) Restore(ctx context.Context, userID, exerciseID string) (*models.ExerciseRecord, error) {
	record, err := s.repo.Restore(ctx, userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.Restore: %w", err)
	}
	if record == nil {
		return nil, ErrExerciseNotFound
	}
	return record, nil
}

func (s *exerciseService) AllExercises(ctx context.Context, userID string, includeInactive bool) ([]*models.ExerciseRecord, error) {
	return s.repo.AllExercises(ctx, userID, includeInactive)
}

func (s *exerciseService) CreateMedia(ctx context.Context, userID, exerciseID, fileName, filePath, mimeType string, fileSize int64) (*models.ExerciseMediaRecord, error) {
	return s.repo.CreateMedia(ctx, userID, exerciseID, fileName, filePath, mimeType, fileSize)
}

func (s *exerciseService) GetMediaByID(ctx context.Context, userID, mediaID string) (*models.ExerciseMediaRecord, error) {
	record, err := s.repo.GetMediaByID(ctx, userID, mediaID)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.GetMediaByID: %w", err)
	}
	if record == nil {
		return nil, ErrExerciseNotFound
	}
	return record, nil
}

func (s *exerciseService) DeleteMedia(ctx context.Context, userID, mediaID string) (*models.ExerciseMediaRecord, error) {
	record, err := s.repo.DeleteMedia(ctx, userID, mediaID)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.DeleteMedia: %w", err)
	}
	if record == nil {
		return nil, ErrExerciseNotFound
	}
	return record, nil
}
```

- [ ] **Step 2: Verify compilation**

Run: `cd apps/api && go vet ./internal/atlas/service/...`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add apps/api/internal/atlas/service/exercise.go
git commit -m "feat(wave-02): add ExerciseService with validation"
```

---

### Task 6: GraphQL schema — exercises.graphql

**Files:**
- Create: `apps/api/internal/atlas/graph/schema/exercises.graphql`

- [ ] **Step 1: Create exercises.graphql schema**

Uses the WAVE-01 `*Result` wrapper pattern (not native GraphQL union types) for consistency with existing schema.

```graphql
type Exercise {
  id: ID!
  userId: ID!
  name: String!
  muscleGroups: [String!]!
  description: String!
  personalNotes: String!
  workingWeight: Float
  isActive: Boolean!
  media: [ExerciseMedia!]!
  createdAt: Time!
  updatedAt: Time!
}

type ExerciseMedia {
  id: ID!
  userId: ID!
  exerciseId: ID!
  fileName: String!
  mimeType: String!
  fileSize: Int!
  createdAt: Time!
}

type ExerciseConnection {
  items: [Exercise!]!
  totalCount: Int!
  pageInfo: ExercisePageInfo!
}

type ExercisePageInfo {
  hasNextPage: Boolean!
  endCursor: String
}

type ExerciseResult {
  exercise: Exercise
  error: ExerciseError
}

type ArchiveResult {
  exercise: Exercise
  error: ExerciseError
}

type ExerciseError {
  message: String!
  code: ExerciseErrorCode!
}

enum ExerciseErrorCode {
  VALIDATION_ERROR
  NOT_FOUND
  UNAUTHORIZED
  INTERNAL_ERROR
}

input CreateExerciseInput {
  name: String!
  muscleGroups: [String!]
  description: String
  personalNotes: String
  workingWeight: Float
}

input UpdateExerciseInput {
  name: String
  muscleGroups: [String!]
  description: String
  personalNotes: String
  workingWeight: Float
}

extend type Query {
  exercises(first: Int = 20, after: String, includeInactive: Boolean = false): ExerciseConnection!
  exercise(id: ID!): ExerciseResult!
  allExercises(includeInactive: Boolean = false): [Exercise!]!
}

extend type Mutation {
  createExercise(input: CreateExerciseInput!): ExerciseResult!
  updateExercise(id: ID!, input: UpdateExerciseInput!): ExerciseResult!
  archiveExercise(id: ID!): ArchiveResult!
  restoreExercise(id: ID!): ArchiveResult!
}
```

- [ ] **Step 2: Run Atlas gqlgen codegen**

Run: `cd apps/api && go run github.com/99designs/gqlgen generate --config atlas-gqlgen.yml`
Expected: generates resolver stubs in `internal/atlas/graph/resolver/` and model updates in `internal/atlas/graph/generated/`

- [ ] **Step 3: Verify generated files compile**

Run: `cd apps/api && go vet ./internal/atlas/graph/...`
Expected: no errors

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/atlas/graph/schema/exercises.graphql
git add apps/api/internal/atlas/graph/generated/
git commit -m "feat(wave-02): add exercises GraphQL schema and generate codegen"
```

---

### Task 7: GraphQL resolvers — exercise.resolvers.go

**Files:**
- Create: `apps/api/internal/atlas/graph/resolver/exercise.resolvers.go`
- Modify: `apps/api/internal/atlas/graph/resolver/resolver.go`

- [ ] **Step 1: Add ExerciseService to Resolver struct**

```go
// Modify apps/api/internal/atlas/graph/resolver/resolver.go
type Resolver struct {
	SettingsService service.SettingsService
	PinService      service.PinService
	ExerciseService service.ExerciseService  // ADD THIS
}
```

- [ ] **Step 2: Create exercise.resolvers.go**

```go
// FILE: apps/api/internal/atlas/graph/resolver/exercise.resolvers.go
package resolver

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
)

func (r *Resolver) Exercises(ctx context.Context, first *int, after *string, includeInactive *bool) (*models.ExerciseConnection, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return emptyConnection(), nil
	}

	limit := 20
	if first != nil && *first > 0 {
		limit = *first
	}
	offset := 0
	if after != nil && *after != "" {
		offset = decodeCursor(*after)
	}
	inactive := false
	if includeInactive != nil {
		inactive = *includeInactive
	}

	records, total, err := r.ExerciseService.List(ctx, userID, inactive, int32(limit), int32(offset))
	if err != nil {
		return nil, err
	}

	items := make([]*models.Exercise, len(records))
	for i, rec := range records {
		items[i] = exerciseRecordToPublic(rec)
	}

	// Fetch media for each exercise
	for i, rec := range records {
		media, err := r.ExerciseService.ListMediaByExerciseID(ctx, userID, rec.ID)
		if err == nil {
			items[i].Media = make([]models.ExerciseMedia, len(media))
			for j, m := range media {
				items[i].Media[j] = mediaRecordToPublic(m)
			}
		}
	}

	nextOffset := offset + limit
	hasNext := nextOffset < int(total)
	var endCursor *string
	if hasNext {
		cursor := encodeCursor(nextOffset)
		endCursor = &cursor
	}

	return &models.ExerciseConnection{
		Items:      items,
		TotalCount: int(total),
		PageInfo: &models.ExercisePageInfo{
			HasNextPage: hasNext,
			EndCursor:   endCursor,
		},
	}, nil
}

func (r *Resolver) Exercise(ctx context.Context, id string) (*models.ExerciseResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.ExerciseResult{
			Error: &models.ExerciseError{
				Message: "unauthorized",
				Code:    models.ExerciseErrorUnauthorized,
			},
		}, nil
	}

	record, err := r.ExerciseService.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrExerciseNotFound) {
			return &models.ExerciseResult{
				Error: &models.ExerciseError{
					Message: "exercise not found",
					Code:    models.ExerciseErrorNotFound,
				},
			}, nil
		}
		return &models.ExerciseResult{
			Error: &models.ExerciseError{
				Message: "internal error",
				Code:    models.ExerciseErrorInternal,
			},
		}, nil
	}

	// Fetch media
	exercise := exerciseRecordToPublic(record)
	media, err := r.ExerciseService.ListMediaByExerciseID(ctx, userID, id)
	if err == nil {
		exercise.Media = make([]models.ExerciseMedia, len(media))
		for j, m := range media {
			exercise.Media[j] = mediaRecordToPublic(m)
		}
	}

	return &models.ExerciseResult{Exercise: exercise}, nil
}

func (r *Resolver) AllExercises(ctx context.Context, includeInactive *bool) ([]*models.Exercise, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return []*models.Exercise{}, nil
	}

	inactive := false
	if includeInactive != nil {
		inactive = *includeInactive
	}

	records, err := r.ExerciseService.AllExercises(ctx, userID, inactive)
	if err != nil {
		return nil, err
	}

	items := make([]*models.Exercise, len(records))
	for i, rec := range records {
		items[i] = exerciseRecordToPublic(rec)
	}
	return items, nil
}

func (r *Resolver) CreateExercise(ctx context.Context, input models.CreateExerciseInput) (*models.ExerciseResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.ExerciseResult{
			Error: &models.ExerciseError{
				Message: "unauthorized",
				Code:    models.ExerciseErrorUnauthorized,
			},
		}, nil
	}

	record, err := r.ExerciseService.Create(ctx, userID, &input)
	if err != nil {
		return &models.ExerciseResult{
			Error: &models.ExerciseError{
				Message: err.Error(),
				Code:    models.ExerciseErrorValidation,
			},
		}, nil
	}

	return &models.ExerciseResult{Exercise: exerciseRecordToPublic(record)}, nil
}

func (r *Resolver) UpdateExercise(ctx context.Context, id string, input models.UpdateExerciseInput) (*models.ExerciseResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.ExerciseResult{
			Error: &models.ExerciseError{
				Message: "unauthorized",
				Code:    models.ExerciseErrorUnauthorized,
			},
		}, nil
	}

	record, err := r.ExerciseService.Update(ctx, userID, id, &input)
	if err != nil {
		if errors.Is(err, atlasService.ErrExerciseNotFound) {
			return &models.ExerciseResult{
				Error: &models.ExerciseError{
					Message: "exercise not found",
					Code:    models.ExerciseErrorNotFound,
				},
			}, nil
		}
		return &models.ExerciseResult{
			Error: &models.ExerciseError{
				Message: err.Error(),
				Code:    models.ExerciseErrorValidation,
			},
		}, nil
	}

	return &models.ExerciseResult{Exercise: exerciseRecordToPublic(record)}, nil
}

func (r *Resolver) ArchiveExercise(ctx context.Context, id string) (*models.ArchiveResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.ArchiveResult{
			Error: &models.ExerciseError{
				Message: "unauthorized",
				Code:    models.ExerciseErrorUnauthorized,
			},
		}, nil
	}

	record, err := r.ExerciseService.Archive(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrExerciseNotFound) {
			return &models.ArchiveResult{
				Error: &models.ExerciseError{
					Message: "exercise not found",
					Code:    models.ExerciseErrorNotFound,
				},
			}, nil
		}
		return &models.ArchiveResult{
			Error: &models.ExerciseError{
				Message: "internal error",
				Code:    models.ExerciseErrorInternal,
			},
		}, nil
	}

	return &models.ArchiveResult{Exercise: exerciseRecordToPublic(record)}, nil
}

func (r *Resolver) RestoreExercise(ctx context.Context, id string) (*models.ArchiveResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.ArchiveResult{
			Error: &models.ExerciseError{
				Message: "unauthorized",
				Code:    models.ExerciseErrorUnauthorized,
			},
		}, nil
	}

	record, err := r.ExerciseService.Restore(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrExerciseNotFound) {
			return &models.ArchiveResult{
				Error: &models.ExerciseError{
					Message: "exercise not found",
					Code:    models.ExerciseErrorNotFound,
				},
			}, nil
		}
		return &models.ArchiveResult{
			Error: &models.ExerciseError{
				Message: "internal error",
				Code:    models.ExerciseErrorInternal,
			},
		}, nil
	}

	return &models.ArchiveResult{Exercise: exerciseRecordToPublic(record)}, nil
}

// gqlgen generated stubs for Exercise media resolver
func (r *Resolver) Exercise_Media(ctx context.Context, obj *models.Exercise) ([]models.ExerciseMedia, error) {
	return obj.Media, nil
}

func exerciseRecordToPublic(rec *models.ExerciseRecord) *models.Exercise {
	mg := rec.MuscleGroups
	if mg == nil {
		mg = []string{}
	}
	return &models.Exercise{
		ID:            rec.ID,
		UserID:        rec.UserID,
		Name:          rec.Name,
		MuscleGroups:  mg,
		Description:   rec.Description,
		PersonalNotes: rec.PersonalNotes,
		WorkingWeight: rec.WorkingWeight,
		IsActive:      rec.IsActive,
		Media:         []models.ExerciseMedia{},
		CreatedAt:     rec.CreatedAt,
		UpdatedAt:     rec.UpdatedAt,
	}
}

func mediaRecordToPublic(rec *models.ExerciseMediaRecord) models.ExerciseMedia {
	return models.ExerciseMedia{
		ID:         rec.ID,
		UserID:     rec.UserID,
		ExerciseID: rec.ExerciseID,
		FileName:   rec.FileName,
		MimeType:   rec.MimeType,
		FileSize:   rec.FileSize,
		CreatedAt:  rec.CreatedAt,
	}
}

func emptyConnection() *models.ExerciseConnection {
	return &models.ExerciseConnection{
		Items:      []*models.Exercise{},
		TotalCount: 0,
		PageInfo: &models.ExercisePageInfo{
			HasNextPage: false,
			EndCursor:   nil,
		},
	}
}

func encodeCursor(offset int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", offset)))
}

func decodeCursor(cursor string) int {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0
	}
	offset, err := strconv.Atoi(string(data))
	if err != nil {
		return 0
	}
	return offset
}
```

- [ ] **Step 3: Add ListMediaByExerciseID method to ExerciseService interface**

Note: The resolver needs `ListMediaByExerciseID` on the service. Add it to the interface in `service/exercise.go`:

```go
ListMediaByExerciseID(ctx context.Context, userID, exerciseID string) ([]*models.ExerciseMediaRecord, error)
```

Implement it:

```go
func (s *exerciseService) ListMediaByExerciseID(ctx context.Context, userID, exerciseID string) ([]*models.ExerciseMediaRecord, error) {
	return s.repo.ListMediaByExerciseID(ctx, userID, exerciseID)
}
```

- [ ] **Step 4: Add ExerciseConnection and ExercisePageInfo types to models**

Add to `models/exercise.go`:

```go
type ExerciseConnection struct {
	Items      []*Exercise       `json:"items"`
	TotalCount int               `json:"totalCount"`
	PageInfo   *ExercisePageInfo `json:"pageInfo"`
}

type ExercisePageInfo struct {
	HasNextPage bool    `json:"hasNextPage"`
	EndCursor   *string `json:"endCursor"`
}
```

- [ ] **Step 5: Verify compilation**

Run: `cd apps/api && go vet ./internal/atlas/...`
Expected: no errors

- [ ] **Step 6: Commit**

```bash
git add apps/api/internal/atlas/graph/resolver/resolver.go
git add apps/api/internal/atlas/graph/resolver/exercise.resolvers.go
git add apps/api/internal/atlas/models/exercise.go
git add apps/api/internal/atlas/service/exercise.go
git commit -m "feat(wave-02): add exercise GraphQL resolvers with archive/restore"
```

---

### Task 8: WAVE-01 media scaffold extension — atlas_media.go

**Files:**
- Modify: `apps/api/internal/handler/atlas_media.go`
- Modify: `apps/api/internal/handler/atlas_media_test.go`

This extends the WAVE-01 scaffold (currently returning 501) with actual upload/download/delete logic using purpose/entity routing. The handler accepts `purpose=EXERCISE_MEDIA` + `exerciseId` in multipart uploads.

- [ ] **Step 1: Extend AtlasMediaUpload with real implementation**

```go
// FILE: apps/api/internal/handler/atlas_media.go
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/middleware"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
	"monorepo-template/libs/go/logger"
)

var (
	allowedImageTypes = map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	allowedVideoTypes = map[string]bool{
		"video/mp4":  true,
		"video/quicktime": true,
		"video/webm": true,
	}
	maxImageSize = int64(25 * 1024 * 1024)  // 25MB
	maxVideoSize = int64(250 * 1024 * 1024) // 250MB
	maxUploadSize = int64(300 * 1024 * 1024) // 300MB
)

type mediaHandler struct {
	exerciseService atlasService.ExerciseService
	mediaBasePath   string
}

func NewMediaHandler(exerciseService atlasService.ExerciseService, mediaBasePath string) *mediaHandler {
	return &mediaHandler{
		exerciseService: exerciseService,
		mediaBasePath:   mediaBasePath,
	}
}

func AtlasMediaUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", http.StatusNotImplemented)
	}
}

// The actual exercise-media upload handler is injected via NewMediaHandler.
// It reads purpose from multipart form:
//   purpose=EXERCISE_MEDIA + exerciseId + file
//
// Implementation uses http.DetectContentType() for MIME validation,
// per-type size limits, and UUID-based file storage.
func (h *mediaHandler) UploadExerciseMedia() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetAtlasUserID(r.Context())
		if userID == "" {
			writeMediaError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
			return
		}

		l := logger.FromCtx(r.Context())

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			l.Warn("[Media][upload] multipart too large", zap.Error(err))
			writeMediaError(w, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", "file exceeds maximum upload size")
			return
		}
		defer func() { _ = r.MultipartForm.RemoveAll() }()

		purpose := r.FormValue("purpose")
		if purpose != "EXERCISE_MEDIA" {
			l.Warn("[Media][upload] unsupported purpose", zap.String("purpose", purpose))
			writeMediaError(w, http.StatusBadRequest, "INVALID_FILE_TYPE", fmt.Sprintf("unsupported purpose: %s", purpose))
			return
		}

		exerciseID := r.FormValue("exerciseId")
		if exerciseID == "" {
			writeMediaError(w, http.StatusBadRequest, "VALIDATION_ERROR", "exerciseId is required")
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			l.Warn("[Media][upload] missing file", zap.Error(err))
			writeMediaError(w, http.StatusBadRequest, "VALIDATION_ERROR", "file is required")
			return
		}
		defer file.Close()

		// Read first 512 bytes for MIME detection
		buf := make([]byte, 512)
		_, err = file.Read(buf)
		if err != nil && err != io.EOF {
			l.Warn("[Media][upload] failed to read file header", zap.Error(err))
			writeMediaError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to read file")
			return
		}
		mimeType := http.DetectContentType(buf)

		if !allowedImageTypes[mimeType] && !allowedVideoTypes[mimeType] {
			l.Warn("[Media][upload] rejected disallowed MIME type", zap.String("mime", mimeType))
			writeMediaError(w, http.StatusBadRequest, "INVALID_FILE_TYPE",
				fmt.Sprintf("file type %s not allowed. Supported: JPEG, PNG, WEBP, MP4, MOV, WEBM", mimeType))
			return
		}

		maxSize := maxImageSize
		if allowedVideoTypes[mimeType] {
			maxSize = maxVideoSize
		}
		if header.Size > maxSize {
			l.Warn("[Media][upload] file too large", zap.Int64("size", header.Size), zap.Int64("max", maxSize))
			writeMediaError(w, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE",
				fmt.Sprintf("file exceeds %dMB limit for %s type", maxSize/(1024*1024), mimeType))
			return
		}

		// Generate UUID-based file name, sanitize original name
		ext := filepath.Ext(header.Filename)
		if ext == ".jpeg" {
			ext = ".jpg"
		}
		storageName := fmt.Sprintf("%s%s", newUUID(), ext)
		dir := filepath.Join(h.mediaBasePath, "exercise", exerciseID)
		if err := os.MkdirAll(dir, 0755); err != nil {
			l.Error("[Media][upload] failed to create directory", zap.Error(err))
			writeMediaError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to store file")
			return
		}
		filePath := filepath.Join(dir, storageName)

		// Seek back to start after reading header
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			l.Error("[Media][upload] seek failed", zap.Error(err))
			writeMediaError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to read file")
			return
		}

		dst, err := os.Create(filePath)
		if err != nil {
			l.Error("[Media][upload] failed to create file", zap.Error(err))
			writeMediaError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to store file")
			return
		}
		defer dst.Close()

		written, err := io.Copy(dst, file)
		if err != nil {
			l.Error("[Media][upload] failed to write file", zap.Error(err))
			os.Remove(filePath)
			writeMediaError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to store file")
			return
		}

		media, err := h.exerciseService.CreateMedia(r.Context(), userID, exerciseID, header.Filename, filePath, mimeType, written)
		if err != nil {
			l.Error("[Media][upload] failed to persist media record", zap.Error(err))
			os.Remove(filePath)
			writeMediaError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to save media record")
			return
		}

		l.Info("[Media][upload] uploaded successfully",
			zap.String("media_id", media.ID),
			zap.String("exercise_id", exerciseID),
			zap.String("mime", mimeType),
			zap.Int64("size", written),
		)

		resp := map[string]interface{}{
			"id":         media.ID,
			"userId":     media.UserID,
			"exerciseId": media.ExerciseID,
			"fileName":   media.FileName,
			"mimeType":   media.MimeType,
			"fileSize":   media.FileSize,
			"createdAt":  media.CreatedAt,
		}
		writeJSON(w, http.StatusCreated, resp)
	}
}

func (h *mediaHandler) DownloadExerciseMedia() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetAtlasUserID(r.Context())
		if userID == "" {
			writeMediaError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
			return
		}

		mediaID := r.PathValue("id")
		if mediaID == "" {
			writeMediaError(w, http.StatusBadRequest, "VALIDATION_ERROR", "media id is required")
			return
		}

		l := logger.FromCtx(r.Context())

		media, err := h.exerciseService.GetMediaByID(r.Context(), userID, mediaID)
		if err != nil {
			l.Warn("[Media][download] media not found", zap.String("media_id", mediaID), zap.Error(err))
			writeMediaError(w, http.StatusNotFound, "NOT_FOUND", "media not found")
			return
		}
		if media == nil {
			writeMediaError(w, http.StatusNotFound, "NOT_FOUND", "media not found")
			return
		}

		file, err := os.Open(media.FilePath)
		if err != nil {
			l.Error("[Media][download] failed to open file", zap.String("path", media.FilePath), zap.Error(err))
			writeMediaError(w, http.StatusNotFound, "NOT_FOUND", "media file not found")
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", media.MimeType)
		w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, media.FileName))
		w.WriteHeader(http.StatusOK)
		io.Copy(w, file)

		l.Info("[Media][download] downloaded successfully", zap.String("media_id", mediaID))
	}
}

func (h *mediaHandler) DeleteExerciseMedia() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetAtlasUserID(r.Context())
		if userID == "" {
			writeMediaError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
			return
		}

		mediaID := r.PathValue("id")
		if mediaID == "" {
			writeMediaError(w, http.StatusBadRequest, "VALIDATION_ERROR", "media id is required")
			return
		}

		l := logger.FromCtx(r.Context())

		media, err := h.exerciseService.DeleteMedia(r.Context(), userID, mediaID)
		if err != nil {
			l.Warn("[Media][delete] media not found", zap.String("media_id", mediaID), zap.Error(err))
			writeMediaError(w, http.StatusNotFound, "NOT_FOUND", "media not found")
			return
		}
		if media == nil {
			writeMediaError(w, http.StatusNotFound, "NOT_FOUND", "media not found")
			return
		}

		// Soft-fail file deletion: log error but return 204
		if err := os.Remove(media.FilePath); err != nil {
			l.Error("[Media][delete] failed to delete physical file", zap.String("path", media.FilePath), zap.Error(err))
		}

		l.Info("[Media][delete] deleted successfully", zap.String("media_id", mediaID))
		w.WriteHeader(http.StatusNoContent)
	}
}

func writeMediaError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// newUUID generates a UUID v4 string without dashes for file naming.
func newUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
```

- [ ] **Step 2: Update the atlas_media_test.go**

Replace the 501 scaffold tests with integration-style tests that use the new handler. Add tests for:
- Valid upload with EXERCISE_MEDIA purpose
- Invalid MIME type rejected
- Missing exerciseId rejected
- File download by ID
- File delete by ID
- Path traversal prevention (UUID-based names)
- Auth error without user context

- [ ] **Step 3: Verify compilation**

Run: `cd apps/api && go vet ./internal/handler/...`
Expected: no errors

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/handler/atlas_media.go
git add apps/api/internal/handler/atlas_media_test.go
git commit -m "feat(wave-02): extend WAVE-01 media scaffold with exercise media upload/download/delete"
```

---

### Task 9: Main wiring — main.go

**Files:**
- Modify: `apps/api/cmd/server/main.go`
- Modify: `apps/api/internal/atlas/graph/resolver/resolver.go` (already done in Task 7)

- [ ] **Step 1: Wire ExerciseRepo and ExerciseService**

After the atlas settings service wiring (around line 149 in main.go), add:

```go
exerciseRepo := atlasPostgres.NewExerciseRepository(db.Pool)
exerciseService := atlasService.NewExerciseService(exerciseRepo)
```

- [ ] **Step 2: Add ExerciseService to Resolver**

```go
atlasRes := &atlasResolver.Resolver{
	SettingsService: atlasSettingsService,
	PinService:      atlasPinService,
	ExerciseService: exerciseService,  // ADD THIS
}
```

- [ ] **Step 3: Wire MediaHandler for exercise media**

The existing scaffold routes already exist under the atlas guarded group. Replace the scaffold handlers with the real implementations:

```go
mediaHandler := handler.NewMediaHandler(exerciseService, cfg.Media.BasePath)

_ = r.Group(func(atlas chi.Router) {
	atlas.Use(atlasMiddleware.AtlasUserContext(atlasBootstrapService))
	atlas.Use(atlasMiddleware.AtlasPinGuard(atlasPinService, atlasPinSessionStore, cfg.AtlasPinSession.CookieName))
	atlas.Handle("/graphql/atlas", atlasSrv)

	// Exercise media via WAVE-01 scaffold routes (replacing 501 stubs)
	atlas.Post("/api/v1/media/upload", mediaHandler.UploadExerciseMedia())
	atlas.Get("/api/v1/media/{id}", mediaHandler.DownloadExerciseMedia())
	atlas.Delete("/api/v1/media/{id}", mediaHandler.DeleteExerciseMedia())
})
```

Note: If `cfg.Media` doesn't exist in `Config`, add a `MediaConfig` struct to `appconfig/config.go`:

```go
type MediaConfig struct {
	BasePath        string `mapstructure:"base_path"`
	MaxUploadSize   int64  `mapstructure:"max_upload_size"`
}

// Inside Config:
Media MediaConfig `mapstructure:"media"`
```

Add default in config defaults:

```go
const defaultMediaBasePath = "./data/media"
```

- [ ] **Step 4: Verify compilation**

Run: `cd apps/api && go build ./cmd/server`
Expected: binary builds without errors

- [ ] **Step 5: Commit**

```bash
git add apps/api/cmd/server/main.go
git add apps/api/internal/appconfig/config.go
git commit -m "feat(wave-02): wire ExerciseRepo, ExerciseService, and media handler in main"
```

---

### Task 10: Tests — repository unit tests

**Files:**
- Create: `apps/api/internal/atlas/repository/postgres/exercise_repo_test.go`

- [ ] **Step 1: Write exercise repository unit tests**

Test CRUD operations, soft archive/restore, pagination, duplicate names, media CRUD.

- [ ] **Step 2: Run tests**

Run: `bunx nx run api:test -- --run '(?i)exercise_repo'`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add apps/api/internal/atlas/repository/postgres/exercise_repo_test.go
git commit -m "test(wave-02): add exercise repository tests"
```

---

### Task 11: Tests — resolvers and round-trip

**Files:**
- Create: `apps/api/internal/atlas/graph/resolver/exercise_resolver_test.go`
- Create: `apps/api/internal/atlas/graph/resolver/exercise_roundtrip_test.go`

- [ ] **Step 1: Write resolver integration tests**

Test create, query, update, archive, restore, pagination, auth errors, duplicate names.

- [ ] **Step 2: Write round-trip integration test**

```go
// TestExerciseRoundTrip: create exercise → upload media → query with media → delete media → archive → restore → verify final state
```

- [ ] **Step 3: Run tests**

Run: `bunx nx run api:test -- --run '(?i)exercise'`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/atlas/graph/resolver/exercise_resolver_test.go
git add apps/api/internal/atlas/graph/resolver/exercise_roundtrip_test.go
git commit -m "test(wave-02): add exercise resolver and round-trip tests"
```

---

### Task 12: Final verification

- [ ] **Step 1: Run all tests**

Run: `bunx nx run api:test`
Expected: all tests pass

- [ ] **Step 2: Run codegen drift check**

Run: `bunx nx run api:codegen && bunx nx run api:codegen:atlas`
Expected: no changes (generated code is up to date)

- [ ] **Step 3: Run lint**

Run: `bunx nx run api:lint`
Expected: no lint errors

- [ ] **Step 4: Run WAVE-01 regression**

Run: `bunx nx run api:test -- --run '(?i)admin_auth'`
Expected: PASS