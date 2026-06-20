# WAVE-04 Cardio and Body Tracking Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement Atlas WAVE-04 backend: CardioEntry, BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto, and WeekFlag CRUD via GraphQL + REST.

**Architecture:** Six new entities under the existing Atlas API module. CardioEntry attaches to DailyLog (auto-created if missing). BodyCheckIn owns nested BodyMeasurements. ProgressPhoto uses a REST handler for binary upload/download (reusing WAVE-01 media storage pattern). WeekFlag is standalone. All GraphQL via `/graphql/atlas`, all REST via PIN-guarded `/api/v1/progress-photos/`.

**Tech Stack:** Go 1.25, PostgreSQL, goose migrations, sqlc v1.30.0, pgx/v5, gqlgen v0.17.49, Nx/Bun command surface.

**Detailed brief:** `docs/prd-wave-details/waves/wave-04.md`

---

## File Inventory

### Create

| Path | Purpose |
| --- | --- |
| `apps/api/internal/repository/postgres/migrations/00086_cardio_entries.sql` | CardioEntry table with daily_log FK. |
| `apps/api/internal/repository/postgres/migrations/00087_body_weight_entries.sql` | BodyWeightEntry table. |
| `apps/api/internal/repository/postgres/migrations/00088_body_check_ins.sql` | BodyCheckIn table. |
| `apps/api/internal/repository/postgres/migrations/00089_body_measurements.sql` | BodyMeasurement table with check_in FK. |
| `apps/api/internal/repository/postgres/migrations/00090_progress_photos.sql` | ProgressPhoto table with check_in FK. |
| `apps/api/internal/repository/postgres/migrations/00091_week_flags.sql` | WeekFlag table. |
| `apps/api/internal/repository/postgres/queries/cardio_entries.sql` | sqlc query source for cardio. |
| `apps/api/internal/repository/postgres/queries/body_weight_entries.sql` | sqlc query source for body weight. |
| `apps/api/internal/repository/postgres/queries/body_check_ins.sql` | sqlc query source for check-ins. |
| `apps/api/internal/repository/postgres/queries/body_measurements.sql` | sqlc query source for measurements. |
| `apps/api/internal/repository/postgres/queries/progress_photos.sql` | sqlc query source for photos. |
| `apps/api/internal/repository/postgres/queries/week_flags.sql` | sqlc query source for week flags. |
| `apps/api/internal/atlas/models/cardio.go` | CardioEntry models, enums, inputs, result union types, error codes. |
| `apps/api/internal/atlas/models/body.go` | BodyWeightEntry, BodyCheckIn, BodyMeasurement models, enums, inputs, result unions, error codes. |
| `apps/api/internal/atlas/models/photo.go` | ProgressPhoto models, PhotoAngle enum, PhotoResult union. |
| `apps/api/internal/atlas/models/week_flag.go` | WeekFlag models, FlagType enum, inputs, result union, error codes. |
| `apps/api/internal/atlas/repository/postgres/cardio_entry_repo.go` | Repository adapter for cardio_entries. |
| `apps/api/internal/atlas/repository/postgres/body_weight_entry_repo.go` | Repository adapter for body_weight_entries. |
| `apps/api/internal/atlas/repository/postgres/body_check_in_repo.go` | Repository adapter for body_check_ins. |
| `apps/api/internal/atlas/repository/postgres/body_measurement_repo.go` | Repository adapter for body_measurements. |
| `apps/api/internal/atlas/repository/postgres/progress_photo_repo.go` | Repository adapter for progress_photos. |
| `apps/api/internal/atlas/repository/postgres/week_flag_repo.go` | Repository adapter for week_flags. |
| `apps/api/internal/atlas/service/cardio.go` | Cardio service with DailyLog auto-creation. |
| `apps/api/internal/atlas/service/body_weight.go` | BodyWeightEntry service. |
| `apps/api/internal/atlas/service/body_checkin.go` | BodyCheckIn service (includes nested measurement management). |
| `apps/api/internal/atlas/service/week_flag.go` | WeekFlag service. |
| `apps/api/internal/atlas/graph/schema/cardio.graphql` | Cardio GraphQL types, enums, queries, mutations. |
| `apps/api/internal/atlas/graph/schema/body.graphql` | Body weight + check-in + measurement GraphQL types, enums, queries, mutations. |
| `apps/api/internal/atlas/graph/schema/progress_photo.graphql` | ProgressPhoto GraphQL types (query only, mutations via REST). |
| `apps/api/internal/atlas/graph/schema/week_flag.graphql` | WeekFlag GraphQL types, enums, queries, mutations. |
| `apps/api/internal/atlas/handler/photo_handler.go` | REST handler for ProgressPhoto upload/download/delete. |

### Modify

| Path | Change |
| --- | --- |
| `apps/api/atlas-gqlgen.yml` | Add all WAVE-04 model bindings. |
| `apps/api/internal/atlas/graph/resolver/resolver.go` | Add CardioService, BodyCheckInService, BodyWeightService, WeekFlagService. |
| `apps/api/cmd/server/main.go` | Wire all WAVE-04 repos, services, resolvers, and photo handler. |
| `apps/api/internal/repository/postgres/generated/*.go` | Regenerated sqlc output. |
| `apps/api/internal/atlas/graph/generated/*.go` | Regenerated atlas gqlgen output. |
| `docs/development-plan.xml` | Add WAVE-04 module facts. |
| `docs/knowledge-graph.xml` | Add WAVE-04 paths and graph annotations. |
| `docs/verification-plan.xml` | Add WAVE-04 verification entries. |
| `docs/prd-waves/waves/wave-04.md` | Update status to `implemented`. |

### Do Not Touch

- Do not touch WAVE-02 exercise code unless importing from it.
- Do not touch WAVE-01 media handler (photo handler is separate).
- Do not touch admin GraphQL (`apps/api/internal/graph/`).
- Do not implement frontend pages, routes, or UI.

---

## Task 1: Domain Models and Enums

**Files:**
- Create: `apps/api/internal/atlas/models/cardio.go`
- Create: `apps/api/internal/atlas/models/body.go`
- Create: `apps/api/internal/atlas/models/photo.go`
- Create: `apps/api/internal/atlas/models/week_flag.go`

- [ ] **Step 1: Create cardio.go models**

Create `apps/api/internal/atlas/models/cardio.go`:

```go
package models

import "time"

type CardioType string
const (
	CardioTypeWalking      CardioType = "walking"
	CardioTypeRunning      CardioType = "running"
	CardioTypeBike       CardioType = "bike"
	CardioTypeElliptical CardioType = "elliptical"
	CardioTypeTreadmill  CardioType = "treadmill"
	CardioTypeOther      CardioType = "other"
)

var ValidCardioTypes = []CardioType{
	CardioTypeWalking, CardioTypeRunning, CardioTypeBike,
	CardioTypeElliptical, CardioTypeTreadmill, CardioTypeOther,
}

type HeartRateZone string
const (
	HRZone1 HeartRateZone = "zone_1"
	HRZone2 HeartRateZone = "zone_2"
	HRZone3 HeartRateZone = "zone_3"
	HRZone4 HeartRateZone = "zone_4"
	HRZone5 HeartRateZone = "zone_5"
	HRZoneUnknown HeartRateZone = "unknown"
)

var ValidHeartRateZones = []HeartRateZone{
	HRZone1, HRZone2, HRZone3, HRZone4, HRZone5, HRZoneUnknown,
}

type CardioEntryRecord struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	DailyLogID     string    `json:"dailyLogId"`
	CardioType     string    `json:"cardioType"`
	DurationMins   int       `json:"durationMinutes"`
	AvgPulse       *int      `json:"avgPulse,omitempty"`
	HeartRateZone  *string   `json:"heartRateZone,omitempty"`
	Notes          *string   `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type CardioEntry struct {
	ID            string  `json:"id"`
	DailyLogID    string  `json:"dailyLogId"`
	CardioType    string  `json:"cardioType"`
	DurationMins  int     `json:"durationMinutes"`
	AvgPulse      *int    `json:"avgPulse,omitempty"`
	HeartRateZone *string `json:"heartRateZone,omitempty"`
	Notes         *string `json:"notes,omitempty"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

type CreateCardioInput struct {
	DailyLogID    string  `json:"dailyLogId"`
	Date          string  `json:"date"`
	CardioType    string  `json:"cardioType"`
	DurationMins  int     `json:"durationMinutes"`
	AvgPulse      *int    `json:"avgPulse,omitempty"`
	HeartRateZone *string `json:"heartRateZone,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}

type UpdateCardioInput struct {
	CardioType    *string `json:"cardioType,omitempty"`
	DurationMins  *int    `json:"durationMinutes,omitempty"`
	AvgPulse      *int    `json:"avgPulse,omitempty"`
	HeartRateZone *string `json:"heartRateZone,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}

type DeleteResult struct {
	Success bool `json:"success"`
}

type CardioEntryResult struct {
	CardioEntry    *CardioEntry    `json:"cardioEntry,omitempty"`
	ValidationErr  *CardioErr      `json:"validationError,omitempty"`
	NotFoundErr    *CardioErr      `json:"notFoundError,omitempty"`
	AuthErr        *CardioErr      `json:"authError,omitempty"`
}

type CardioDeleteResult struct {
	DeleteResult   *DeleteResult   `json:"deleteResult,omitempty"`
	NotFoundErr    *CardioErr      `json:"notFoundError,omitempty"`
	AuthErr        *CardioErr      `json:"authError,omitempty"`
}

type CardioErr struct {
	Message string         `json:"message"`
	Code    CardioErrorCode `json:"code"`
}

type CardioErrorCode string
const (
	CardioErrValidation CardioErrorCode = "VALIDATION_ERROR"
	CardioErrNotFound   CardioErrorCode = "NOT_FOUND"
	CardioErrAuth       CardioErrorCode = "AUTH_ERROR"
	CardioErrInternal   CardioErrorCode = "INTERNAL_ERROR"
)

func CardioEntryFromRecord(r *CardioEntryRecord) *CardioEntry {
	if r == nil {
		return nil
	}
	return &CardioEntry{
		ID:            r.ID,
		DailyLogID:    r.DailyLogID,
		CardioType:    r.CardioType,
		DurationMins:  r.DurationMins,
		AvgPulse:      r.AvgPulse,
		HeartRateZone: r.HeartRateZone,
		Notes:         r.Notes,
		CreatedAt:     r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     r.UpdatedAt.Format(time.RFC3339),
	}
}

func CardioEntriesFromRecords(records []CardioEntryRecord) []CardioEntry {
	if records == nil {
		return nil
	}
	out := make([]CardioEntry, len(records))
	for i, r := range records {
		*out[i].fromRecord(&r)
	}
	return out
}
```

- [ ] **Step 2: Create body.go models**

Create `apps/api/internal/atlas/models/body.go`:

```go
package models

import "time"

type MeasurementType string
const (
	MeasNeck     MeasurementType = "neck"
	MeasShoulders MeasurementType = "shoulders"
	MeasForearms MeasurementType = "forearms"
	MeasBiceps   MeasurementType = "biceps"
	MeasChest    MeasurementType = "chest"
	MeasWaist    MeasurementType = "waist"
	MeasAbdomen  MeasurementType = "abdomen"
	MeasHips     MeasurementType = "hips"
	MeasThigh    MeasurementType = "thigh"
	MeasCalves   MeasurementType = "calf"
)

var ValidMeasurementTypes = []MeasurementType{
	MeasNeck, MeasShoulders, MeasForearms, MeasBiceps,
	MeasChest, MeasWaist, MeasAbdomen, MeasHips, MeasThigh, MeasCalves,
}

var PairedMeasurementTypes = map[MeasurementType]bool{
	MeasForearms: true,
	MeasBiceps:   true,
	MeasThigh:    true,
	MeasCalves:   true,
}

type MeasurementSide string
const (
	SideLeft  MeasurementSide = "left"
	SideRight MeasurementSide = "right"
)

type BodyWeightSource string
const (
	BWSourceScale   BodyWeightSource = "scale"
	BWSourceManual  BodyWeightSource = "manual"
	BWSourceUnknown BodyWeightSource = "unknown"
)

var ValidBodyWeightSources = []BodyWeightSource{
	BWSourceScale, BWSourceManual, BWSourceUnknown,
}

// BodyWeightEntry

type BodyWeightEntryRecord struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Date      time.Time `json:"date"`
	Weight    float64   `json:"weight"`
	Source    string    `json:"source"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type BodyWeightEntry struct {
	ID        string   `json:"id"`
	Date      string   `json:"date"`
	Weight    float64  `json:"weight"`
	Source    string   `json:"source"`
	Notes     *string  `json:"notes,omitempty"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

type CreateBodyWeightInput struct {
	Date   string  `json:"date"`
	Weight float64 `json:"weight"`
	Source string  `json:"source"`
	Notes  *string `json:"notes,omitempty"`
}

type UpdateBodyWeightInput struct {
	Weight *float64 `json:"weight,omitempty"`
	Source *string  `json:"source,omitempty"`
	Notes  *string  `json:"notes,omitempty"`
}

type BodyWeightResult struct {
	BodyWeightEntry *BodyWeightEntry `json:"bodyWeightEntry,omitempty"`
	ValidationErr   *BodyErr         `json:"validationError,omitempty"`
	NotFoundErr     *BodyErr         `json:"notFoundError,omitempty"`
	AuthErr         *BodyErr         `json:"authError,omitempty"`
}

type BodyWeightDeleteResult struct {
	DeleteResult *DeleteResult `json:"deleteResult,omitempty"`
	NotFoundErr  *BodyErr      `json:"notFoundError,omitempty"`
	AuthErr      *BodyErr      `json:"authError,omitempty"`
}

// BodyCheckIn

type BodyCheckInRecord struct {
	ID               string    `json:"id"`
	UserID           string    `json:"userId"`
	Date             time.Time `json:"date"`
	Weight           *float64  `json:"weight,omitempty"`
	BodyFatPct       *float64  `json:"bodyFatPercentage,omitempty"`
	Notes            *string   `json:"notes,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type BodyCheckIn struct {
	ID               string              `json:"id"`
	Date             string              `json:"date"`
	Weight           *float64            `json:"weight,omitempty"`
	BodyFatPct       *float64            `json:"bodyFatPercentage,omitempty"`
	Notes            *string             `json:"notes,omitempty"`
	Measurements     []BodyMeasurement   `json:"measurements"`
	ProgressPhotos   []ProgressPhoto     `json:"progressPhotos"`
	CreatedAt        string              `json:"createdAt"`
	UpdatedAt        string              `json:"updatedAt"`
}

type CreateCheckInInput struct {
	Date       string   `json:"date"`
	Weight     *float64 `json:"weight,omitempty"`
	BodyFatPct *float64 `json:"bodyFatPercentage,omitempty"`
	Notes      *string  `json:"notes,omitempty"`
}

type UpdateCheckInInput struct {
	Weight     *float64 `json:"weight,omitempty"`
	BodyFatPct *float64 `json:"bodyFatPercentage,omitempty"`
	Notes      *string  `json:"notes,omitempty"`
}

type BodyCheckInResult struct {
	BodyCheckIn    *BodyCheckIn `json:"bodyCheckIn,omitempty"`
	ValidationErr  *BodyErr     `json:"validationError,omitempty"`
	NotFoundErr    *BodyErr     `json:"notFoundError,omitempty"`
	AuthErr        *BodyErr     `json:"authError,omitempty"`
}

type BodyCheckInDeleteResult struct {
	DeleteResult *DeleteResult `json:"deleteResult,omitempty"`
	NotFoundErr  *BodyErr      `json:"notFoundError,omitempty"`
	AuthErr      *BodyErr      `json:"authError,omitempty"`
}

// BodyMeasurement

type BodyMeasurementRecord struct {
	ID              string    `json:"id"`
	CheckInID       string    `json:"checkInId"`
	MeasurementType string    `json:"measurementType"`
	Side            *string   `json:"side,omitempty"`
	Value           float64   `json:"value"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type BodyMeasurement struct {
	ID              string  `json:"id"`
	CheckInID       string  `json:"checkInId"`
	MeasurementType string  `json:"measurementType"`
	Side            *string `json:"side,omitempty"`
	Value           float64 `json:"value"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

type CreateMeasurementInput struct {
	MeasurementType string  `json:"measurementType"`
	Side            *string `json:"side,omitempty"`
	Value           float64 `json:"value"`
}

type UpdateMeasurementInput struct {
	Side  *string  `json:"side,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type MeasurementResult struct {
	Measurement    *BodyMeasurement `json:"measurement,omitempty"`
	ValidationErr  *BodyErr         `json:"validationError,omitempty"`
	NotFoundErr    *BodyErr         `json:"notFoundError,omitempty"`
	AuthErr        *BodyErr         `json:"authError,omitempty"`
}

type MeasurementDeleteResult struct {
	DeleteResult *DeleteResult `json:"deleteResult,omitempty"`
	NotFoundErr  *BodyErr      `json:"notFoundError,omitempty"`
	AuthErr      *BodyErr      `json:"authError,omitempty"`
}

// BodyErr shared error type for all body tracking

type BodyErr struct {
	Message string      `json:"message"`
	Code    BodyErrorCode `json:"code"`
}

type BodyErrorCode string
const (
	BodyErrValidation BodyErrorCode = "VALIDATION_ERROR"
	BodyErrNotFound   BodyErrorCode = "NOT_FOUND"
	BodyErrAuth       BodyErrorCode = "AUTH_ERROR"
	BodyErrInternal   BodyErrorCode = "INTERNAL_ERROR"
)

func BodyWeightFromRecord(r *BodyWeightEntryRecord) *BodyWeightEntry {
	if r == nil { return nil }
	return &BodyWeightEntry{
		ID: r.ID, Date: r.Date.Format("2006-01-02"),
		Weight: r.Weight, Source: r.Source,
		Notes: r.Notes, CreatedAt: r.CreatedAt.Format(time.RFC3339),
		UpdatedAt: r.UpdatedAt.Format(time.RFC3339),
	}
}
```

- [ ] **Step 3: Create photo.go models**

Create `apps/api/internal/atlas/models/photo.go`:

```go
package models

import "time"

type PhotoAngle string
const (
	PhotoAngleFront  PhotoAngle = "front"
	PhotoAngleSide   PhotoAngle = "side"
	PhotoAngleBack   PhotoAngle = "back"
	PhotoAngleCustom PhotoAngle = "custom"
)

var ValidPhotoAngles = []PhotoAngle{
	PhotoAngleFront, PhotoAngleSide, PhotoAngleBack, PhotoAngleCustom,
}

type ProgressPhotoRecord struct {
	ID               string    `json:"id"`
	CheckInID        string    `json:"checkInId"`
	FilePath         string    `json:"filePath"`
	OriginalFileName string    `json:"originalFileName"`
	MimeType         string    `json:"mimeType"`
	SizeBytes        int64     `json:"sizeBytes"`
	Angle            *string   `json:"angle,omitempty"`
	Label            *string   `json:"label,omitempty"`
	Notes            *string   `json:"notes,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type ProgressPhoto struct {
	ID               string  `json:"id"`
	CheckInID        string  `json:"checkInId"`
	OriginalFileName string  `json:"originalFileName"`
	MimeType         string  `json:"mimeType"`
	SizeBytes        int64   `json:"sizeBytes"`
	Angle            *string `json:"angle,omitempty"`
	Label            *string `json:"label,omitempty"`
	Notes            *string `json:"notes,omitempty"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

func PhotoFromRecord(r *ProgressPhotoRecord) *ProgressPhoto {
	if r == nil { return nil }
	return &ProgressPhoto{
		ID: r.ID, CheckInID: r.CheckInID,
		OriginalFileName: r.OriginalFileName, MimeType: r.MimeType,
		SizeBytes: r.SizeBytes, Angle: r.Angle, Label: r.Label, Notes: r.Notes,
		CreatedAt: r.CreatedAt.Format(time.RFC3339),
		UpdatedAt: r.UpdatedAt.Format(time.RFC3339),
	}
}
```

- [ ] **Step 4: Create week_flag.go models**

Create `apps/api/internal/atlas/models/week_flag.go`:

```go
package models

import "time"

type FlagType string
const (
	FlagPoorSleep      FlagType = "poor_sleep"
	FlagHighStress     FlagType = "high_stress"
	FlagIllness        FlagType = "illness"
	FlagInjuryPain     FlagType = "injury_pain"
	FlagCycle          FlagType = "cycle"
	FlagCalorieDeficit FlagType = "calorie_deficit"
	FlagSurplus        FlagType = "surplus"
	FlagMaintenance    FlagType = "maintenance"
	FlagMissedWorkouts FlagType = "missed_workouts"
	FlagTravel         FlagType = "travel"
)

var ValidFlagTypes = []FlagType{
	FlagPoorSleep, FlagHighStress, FlagIllness, FlagInjuryPain, FlagCycle,
	FlagCalorieDeficit, FlagSurplus, FlagMaintenance, FlagMissedWorkouts, FlagTravel,
}

type WeekFlagRecord struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	WeekStartDate time.Time `json:"weekStartDate"`
	FlagType      string    `json:"flagType"`
	Notes         *string   `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type WeekFlag struct {
	ID            string  `json:"id"`
	WeekStartDate string  `json:"weekStartDate"`
	FlagType      string  `json:"flagType"`
	Notes         *string `json:"notes,omitempty"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

type CreateWeekFlagInput struct {
	WeekStartDate string  `json:"weekStartDate"`
	FlagType      string  `json:"flagType"`
	Notes         *string `json:"notes,omitempty"`
}

type WeekFlagResult struct {
	WeekFlag       *WeekFlag    `json:"weekFlag,omitempty"`
	ValidationErr  *WeekFlagErr `json:"validationError,omitempty"`
	NotFoundErr    *WeekFlagErr `json:"notFoundError,omitempty"`
	AuthErr        *WeekFlagErr `json:"authError,omitempty"`
}

type WeekFlagDeleteResult struct {
	DeleteResult *DeleteResult `json:"deleteResult,omitempty"`
	NotFoundErr  *WeekFlagErr  `json:"notFoundError,omitempty"`
	AuthErr      *WeekFlagErr  `json:"authError,omitempty"`
}

type WeekFlagErr struct {
	Message string          `json:"message"`
	Code    WeekFlagErrorCode `json:"code"`
}

type WeekFlagErrorCode string
const (
	WfErrValidation WeekFlagErrorCode = "VALIDATION_ERROR"
	WfErrNotFound   WeekFlagErrorCode = "NOT_FOUND"
	WfErrAuth       WeekFlagErrorCode = "AUTH_ERROR"
	WfErrInternal   WeekFlagErrorCode = "INTERNAL_ERROR"
)

func WeekFlagFromRecord(r *WeekFlagRecord) *WeekFlag {
	if r == nil { return nil }
	return &WeekFlag{
		ID: r.ID,
		WeekStartDate: r.WeekStartDate.Format("2006-01-02"),
		FlagType: r.FlagType, Notes: r.Notes,
		CreatedAt: r.CreatedAt.Format(time.RFC3339),
		UpdatedAt: r.UpdatedAt.Format(time.RFC3339),
	}
}
```

- [ ] **Step 5: Run model compilation check**

Run: `cd apps/api && go build ./internal/atlas/models`
Expected: exits 0 (no output).

- [ ] **Step 6: Commit**

```bash
git add apps/api/internal/atlas/models/cardio.go apps/api/internal/atlas/models/body.go apps/api/internal/atlas/models/photo.go apps/api/internal/atlas/models/week_flag.go
git commit -m "feat(wave-04): add domain models and enums"
```

---

## Task 2: Database Migrations

**Files:**
- Create: `apps/api/internal/repository/postgres/migrations/00086_cardio_entries.sql`
- Create: `apps/api/internal/repository/postgres/migrations/00087_body_weight_entries.sql`
- Create: `apps/api/internal/repository/postgres/migrations/00088_body_check_ins.sql`
- Create: `apps/api/internal/repository/postgres/migrations/00089_body_measurements.sql`
- Create: `apps/api/internal/repository/postgres/migrations/00090_progress_photos.sql`
- Create: `apps/api/internal/repository/postgres/migrations/00091_week_flags.sql`

**Prerequisite:** `daily_logs` table must exist (WAVE-03 migration 00083 or equivalent). If WAVE-03 is not deployed, add a daily_logs migration as 00083 before these.

- [ ] **Step 1: Cardio entries migration**

Create `apps/api/internal/repository/postgres/migrations/00086_cardio_entries.sql`:

```sql
-- +goose Up
CREATE TABLE cardio_entries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES atlas_users(id),
    daily_log_id    UUID NOT NULL REFERENCES daily_logs(id) ON DELETE CASCADE,
    cardio_type     VARCHAR(50) NOT NULL,
    duration_minutes INTEGER NOT NULL CHECK (duration_minutes > 0),
    avg_pulse       INTEGER CHECK (avg_pulse IS NULL OR avg_pulse > 0),
    heart_rate_zone VARCHAR(20),
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_cardio_entries_daily_log ON cardio_entries (daily_log_id);
CREATE INDEX idx_cardio_entries_user_date ON cardio_entries (user_id, daily_log_id);

-- +goose Down
DROP TABLE IF EXISTS cardio_entries;
```

- [ ] **Step 2: Body weight entries migration**

Create `apps/api/internal/repository/postgres/migrations/00087_body_weight_entries.sql`:

```sql
-- +goose Up
CREATE TABLE body_weight_entries (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES atlas_users(id),
    date       DATE NOT NULL,
    weight     REAL NOT NULL CHECK (weight > 0),
    source     VARCHAR(20) NOT NULL,
    notes      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_body_weight_user_date ON body_weight_entries (user_id, date DESC);

-- +goose Down
DROP TABLE IF EXISTS body_weight_entries;
```

- [ ] **Step 3: Body check-ins migration**

Create `apps/api/internal/repository/postgres/migrations/00088_body_check_ins.sql`:

```sql
-- +goose Up
CREATE TABLE body_check_ins (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES atlas_users(id),
    date                DATE NOT NULL,
    weight              REAL CHECK (weight IS NULL OR weight > 0),
    body_fat_percentage REAL CHECK (body_fat_percentage IS NULL OR (body_fat_percentage > 0 AND body_fat_percentage <= 100)),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_body_check_ins_user_date UNIQUE (user_id, date)
);

CREATE INDEX idx_body_check_ins_date ON body_check_ins (user_id, date DESC);

-- +goose Down
DROP TABLE IF EXISTS body_check_ins;
```

- [ ] **Step 4: Body measurements migration**

Create `apps/api/internal/repository/postgres/migrations/00089_body_measurements.sql`:

```sql
-- +goose Up
CREATE TABLE body_measurements (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    check_in_id      UUID NOT NULL REFERENCES body_check_ins(id) ON DELETE CASCADE,
    measurement_type VARCHAR(30) NOT NULL,
    side             VARCHAR(10),
    value            REAL NOT NULL CHECK (value > 0),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_body_measurements_checkin_type_side UNIQUE (check_in_id, measurement_type, side)
);

CREATE INDEX idx_body_measurements_checkin ON body_measurements (check_in_id);

-- +goose Down
DROP TABLE IF EXISTS body_measurements;
```

- [ ] **Step 5: Progress photos migration**

Create `apps/api/internal/repository/postgres/migrations/00090_progress_photos.sql`:

```sql
-- +goose Up
CREATE TABLE progress_photos (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    check_in_id        UUID NOT NULL REFERENCES body_check_ins(id) ON DELETE CASCADE,
    file_path          VARCHAR(500) NOT NULL,
    original_file_name VARCHAR(255) NOT NULL,
    mime_type          VARCHAR(50) NOT NULL,
    size_bytes         BIGINT NOT NULL,
    angle              VARCHAR(20),
    label              VARCHAR(255),
    notes              TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_progress_photos_checkin ON progress_photos (check_in_id);

-- +goose Down
DROP TABLE IF EXISTS progress_photos;
```

- [ ] **Step 6: Week flags migration**

Create `apps/api/internal/repository/postgres/migrations/00091_week_flags.sql`:

```sql
-- +goose Up
CREATE TABLE week_flags (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES atlas_users(id),
    week_start_date DATE NOT NULL,
    flag_type       VARCHAR(30) NOT NULL,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_week_flags_week_type UNIQUE (week_start_date, flag_type)
);

CREATE INDEX idx_week_flags_week ON week_flags (week_start_date);

-- +goose Down
DROP TABLE IF EXISTS week_flags;
```

- [ ] **Step 7: Run migration syntax check**

Run: `cd apps/api && go build ./internal/repository/postgres`
Expected: exits 0.

- [ ] **Step 8: Commit**

```bash
git add apps/api/internal/repository/postgres/migrations/00086_cardio_entries.sql apps/api/internal/repository/postgres/migrations/00087_body_weight_entries.sql apps/api/internal/repository/postgres/migrations/00088_body_check_ins.sql apps/api/internal/repository/postgres/migrations/00089_body_measurements.sql apps/api/internal/repository/postgres/migrations/00090_progress_photos.sql apps/api/internal/repository/postgres/migrations/00091_week_flags.sql
git commit -m "feat(wave-04): add database migrations"
```

---

## Task 3: sqlc Queries

**Files:**
- Create: `apps/api/internal/repository/postgres/queries/cardio_entries.sql`
- Create: `apps/api/internal/repository/postgres/queries/body_weight_entries.sql`
- Create: `apps/api/internal/repository/postgres/queries/body_check_ins.sql`
- Create: `apps/api/internal/repository/postgres/queries/body_measurements.sql`
- Create: `apps/api/internal/repository/postgres/queries/progress_photos.sql`
- Create: `apps/api/internal/repository/postgres/queries/week_flags.sql`

- [ ] **Step 1: Cardio entries sqlc queries**

Create `apps/api/internal/repository/postgres/queries/cardio_entries.sql`:

```sql
-- name: CreateCardioEntry :one
INSERT INTO cardio_entries (user_id, daily_log_id, cardio_type, duration_minutes, avg_pulse, heart_rate_zone, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetCardioEntryByID :one
SELECT * FROM cardio_entries WHERE id = $1 AND user_id = $2;

-- name: ListCardioEntriesByDailyLog :many
SELECT * FROM cardio_entries WHERE daily_log_id = $1 AND user_id = $2 ORDER BY created_at ASC;

-- name: UpdateCardioEntry :one
UPDATE cardio_entries
SET cardio_type = COALESCE($3, cardio_type),
    duration_minutes = COALESCE($4, duration_minutes),
    avg_pulse = $5,
    heart_rate_zone = $6,
    notes = $7,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteCardioEntry :one
DELETE FROM cardio_entries WHERE id = $1 AND user_id = $2 RETURNING id;
```

- [ ] **Step 2: Body weight entries sqlc queries**

Create `apps/api/internal/repository/postgres/queries/body_weight_entries.sql`:

```sql
-- name: CreateBodyWeightEntry :one
INSERT INTO body_weight_entries (user_id, date, weight, source, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetBodyWeightEntryByID :one
SELECT * FROM body_weight_entries WHERE id = $1 AND user_id = $2;

-- name: ListBodyWeightEntriesByDateRange :many
SELECT * FROM body_weight_entries
WHERE user_id = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC, created_at DESC;

-- name: ListBodyWeightEntries :many
SELECT * FROM body_weight_entries
WHERE user_id = $1
ORDER BY date DESC, created_at DESC
LIMIT $2;

-- name: LatestBodyWeightEntry :one
SELECT * FROM body_weight_entries
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateBodyWeightEntry :one
UPDATE body_weight_entries
SET weight = COALESCE($3, weight),
    source = COALESCE($4, source),
    notes = $5,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteBodyWeightEntry :one
DELETE FROM body_weight_entries WHERE id = $1 AND user_id = $2 RETURNING id;
```

- [ ] **Step 3: Body check-ins sqlc queries**

Create `apps/api/internal/repository/postgres/queries/body_check_ins.sql`:

```sql
-- name: CreateBodyCheckIn :one
INSERT INTO body_check_ins (user_id, date, weight, body_fat_percentage, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetBodyCheckInByID :one
SELECT * FROM body_check_ins WHERE id = $1 AND user_id = $2;

-- name: ListBodyCheckInsByDateRange :many
SELECT * FROM body_check_ins
WHERE user_id = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC;

-- name: ListBodyCheckIns :many
SELECT * FROM body_check_ins
WHERE user_id = $1
ORDER BY date DESC
LIMIT $2;

-- name: UpdateBodyCheckIn :one
UPDATE body_check_ins
SET weight = $3,
    body_fat_percentage = $4,
    notes = $5,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteBodyCheckIn :one
DELETE FROM body_check_ins WHERE id = $1 AND user_id = $2 RETURNING id;
```

- [ ] **Step 4: Body measurements sqlc queries**

Create `apps/api/internal/repository/postgres/queries/body_measurements.sql`:

```sql
-- name: CreateBodyMeasurement :one
INSERT INTO body_measurements (check_in_id, measurement_type, side, value)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetBodyMeasurementByID :one
SELECT * FROM body_measurements WHERE id = $1;

-- name: ListBodyMeasurementsByCheckIn :many
SELECT * FROM body_measurements WHERE check_in_id = $1 ORDER BY measurement_type, side;

-- name: UpdateBodyMeasurement :one
UPDATE body_measurements
SET side = $2,
    value = COALESCE($3, value),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteBodyMeasurement :one
DELETE FROM body_measurements WHERE id = $1 RETURNING id;

-- name: DeleteBodyMeasurementsByCheckIn :exec
DELETE FROM body_measurements WHERE check_in_id = $1;
```

- [ ] **Step 5: Progress photos sqlc queries**

Create `apps/api/internal/repository/postgres/queries/progress_photos.sql`:

```sql
-- name: CreateProgressPhoto :one
INSERT INTO progress_photos (check_in_id, file_path, original_file_name, mime_type, size_bytes, angle, label, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetProgressPhotoByID :one
SELECT * FROM progress_photos WHERE id = $1;

-- name: ListProgressPhotosByCheckIn :many
SELECT * FROM progress_photos WHERE check_in_id = $1 ORDER BY created_at ASC;

-- name: DeleteProgressPhoto :one
DELETE FROM progress_photos WHERE id = $1 RETURNING *;
```

- [ ] **Step 6: Week flags sqlc queries**

Create `apps/api/internal/repository/postgres/queries/week_flags.sql`:

```sql
-- name: CreateWeekFlag :one
INSERT INTO week_flags (user_id, week_start_date, flag_type, notes)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetWeekFlagByID :one
SELECT * FROM week_flags WHERE id = $1 AND user_id = $2;

-- name: ListWeekFlagsByWeek :many
SELECT * FROM week_flags
WHERE user_id = $1 AND week_start_date = $2
ORDER BY flag_type;

-- name: DeleteWeekFlag :one
DELETE FROM week_flags WHERE id = $1 AND user_id = $2 RETURNING id;
```

- [ ] **Step 7: Run sqlc codegen**

Run: `cd apps/api && sqlc generate`
Expected: exits 0, generates files in `internal/repository/postgres/generated/`.

Verify generated files exist:
Run: `ls apps/api/internal/repository/postgres/generated/*.go | wc -l`
Expected: > 0 (generated files for new queries).

- [ ] **Step 8: Run build check after codegen**

Run: `cd apps/api && go build ./internal/repository/postgres/generated`
Expected: exits 0.

- [ ] **Step 9: Commit**

```bash
git add apps/api/internal/repository/postgres/queries/cardio_entries.sql apps/api/internal/repository/postgres/queries/body_weight_entries.sql apps/api/internal/repository/postgres/queries/body_check_ins.sql apps/api/internal/repository/postgres/queries/body_measurements.sql apps/api/internal/repository/postgres/queries/progress_photos.sql apps/api/internal/repository/postgres/queries/week_flags.sql apps/api/internal/repository/postgres/generated/
git commit -m "feat(wave-04): add sqlc queries and generated code"
```

---

## Task 4: Repository Adapters

**Files:**
- Create: `apps/api/internal/atlas/repository/postgres/cardio_entry_repo.go`
- Create: `apps/api/internal/atlas/repository/postgres/body_weight_entry_repo.go`
- Create: `apps/api/internal/atlas/repository/postgres/body_check_in_repo.go`
- Create: `apps/api/internal/atlas/repository/postgres/body_measurement_repo.go`
- Create: `apps/api/internal/atlas/repository/postgres/progress_photo_repo.go`
- Create: `apps/api/internal/atlas/repository/postgres/week_flag_repo.go`

Each repo follows the exact pattern from `exercise_repo.go`: interface + struct + constructor + methods that parse UUIDs, call sqlc generated queries, convert rows to model records via helpers.

- [ ] **Step 1: Cardio entry repo**

Create `apps/api/internal/atlas/repository/postgres/cardio_entry_repo.go`:

```go
package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"monorepo-template/apps/api/internal/atlas/models"
	generated "monorepo-template/apps/api/internal/repository/postgres/generated"
)

type CardioEntryRepository interface {
	Create(ctx context.Context, userID string, dailyLogID string, input models.CreateCardioInput) (*models.CardioEntryRecord, error)
	GetByID(ctx context.Context, userID, id string) (*models.CardioEntryRecord, error)
	ListByDailyLog(ctx context.Context, userID, dailyLogID string) ([]models.CardioEntryRecord, error)
	Update(ctx context.Context, userID, id string, input models.UpdateCardioInput) (*models.CardioEntryRecord, error)
	Delete(ctx context.Context, userID, id string) error
}

type cardioEntryRepository struct {
	q *generated.Queries
}

func NewCardioEntryRepository(pool *pgxpool.Pool) CardioEntryRepository {
	return &cardioEntryRepository{q: generated.New(pool)}
}

func (r *cardioEntryRepository) Create(ctx context.Context, userID, dailyLogID string, input models.CreateCardioInput) (*models.CardioEntryRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil { return nil, fmt.Errorf("CardioEntryRepo.Create: %w", err) }
	did, err := uuidFromString(dailyLogID)
	if err != nil { return nil, fmt.Errorf("CardioEntryRepo.Create: %w", err) }
	row, err := r.q.CreateCardioEntry(ctx, generated.CreateCardioEntryParams{
		UserID: uid, DailyLogID: did,
		CardioType: input.CardioType, DurationMinutes: int32(input.DurationMins),
		AvgPulse: int32PtrOrNil(input.AvgPulse), HeartRateZone: nullableText(input.HeartRateZone),
		Notes: nullableText(input.Notes),
	})
	if err != nil { return nil, fmt.Errorf("CardioEntryRepo.Create: %w", err) }
	return cardioRecordFromRow(row), nil
}

func (r *cardioEntryRepository) GetByID(ctx context.Context, userID, id string) (*models.CardioEntryRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil { return nil, fmt.Errorf("CardioEntryRepo.GetByID: %w", err) }
	eid, err := uuidFromString(id)
	if err != nil { return nil, fmt.Errorf("CardioEntryRepo.GetByID: %w", err) }
	row, err := r.q.GetCardioEntryByID(ctx, generated.GetCardioEntryByIDParams{ID: eid, UserID: uid})
	if err != nil { return nil, fmt.Errorf("CardioEntryRepo.GetByID: %w", err) }
	return cardioRecordFromRow(row), nil
}

func (r *cardioEntryRepository) ListByDailyLog(ctx context.Context, userID, dailyLogID string) ([]models.CardioEntryRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil { return nil, fmt.Errorf("CardioEntryRepo.ListByDailyLog: %w", err) }
	did, err := uuidFromString(dailyLogID)
	if err != nil { return nil, fmt.Errorf("CardioEntryRepo.ListByDailyLog: %w", err) }
	rows, err := r.q.ListCardioEntriesByDailyLog(ctx, generated.ListCardioEntriesByDailyLogParams{DailyLogID: did, UserID: uid})
	if err != nil { return nil, fmt.Errorf("CardioEntryRepo.ListByDailyLog: %w", err) }
	return cardioRecordsFromRows(rows), nil
}

func (r *cardioEntryRepository) Update(ctx context.Context, userID, id string, input models.UpdateCardioInput) (*models.CardioEntryRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil { return nil, fmt.Errorf("CardioEntryRepo.Update: %w", err) }
	eid, err := uuidFromString(id)
	if err != nil { return nil, fmt.Errorf("CardioEntryRepo.Update: %w", err) }
	row, err := r.q.UpdateCardioEntry(ctx, generated.UpdateCardioEntryParams{
		ID: eid, UserID: uid,
		CardioType:       nullableText(input.CardioType),
		DurationMinutes:  int32PtrOrNil(input.DurationMins),
		AvgPulse:         int32PtrOrNil(input.AvgPulse),
		HeartRateZone:    nullableText(input.HeartRateZone),
		Notes:            nullableText(input.Notes),
	})
	if err != nil { return nil, fmt.Errorf("CardioEntryRepo.Update: %w", err) }
	return cardioRecordFromRow(row), nil
}

func (r *cardioEntryRepository) Delete(ctx context.Context, userID, id string) error {
	uid, err := uuidFromString(userID)
	if err != nil { return fmt.Errorf("CardioEntryRepo.Delete: %w", err) }
	eid, err := uuidFromString(id)
	if err != nil { return fmt.Errorf("CardioEntryRepo.Delete: %w", err) }
	_, err = r.q.DeleteCardioEntry(ctx, generated.DeleteCardioEntryParams{ID: eid, UserID: uid})
	if err != nil { return fmt.Errorf("CardioEntryRepo.Delete: %w", err) }
	return nil
}

func cardioRecordFromRow(row generated.AtlasCardioEntry) *models.CardioEntryRecord {
	return &models.CardioEntryRecord{
		ID: row.ID.String(), UserID: row.UserID.String(),
		DailyLogID: row.DailyLogID.String(), CardioType: row.CardioType,
		DurationMins: int(row.DurationMinutes),
		AvgPulse: intPtrFrom32(row.AvgPulse),
		HeartRateZone: row.HeartRateZone,
		Notes: row.Notes,
		CreatedAt: row.CreatedAt.Time, UpdatedAt: row.UpdatedAt.Time,
	}
}

func cardioRecordsFromRows(rows []generated.AtlasCardioEntry) []models.CardioEntryRecord {
	out := make([]models.CardioEntryRecord, len(rows))
	for i, r := range rows {
		rec := cardioRecordFromRow(r)
		out[i] = *rec
	}
	return out
}
```

- [ ] **Step 2: Body weight entry repo**

Create `apps/api/internal/atlas/repository/postgres/body_weight_entry_repo.go` — follows same pattern as cardio, calling `CreateBodyWeightEntry`, `GetBodyWeightEntryByID`, `ListBodyWeightEntriesByDateRange`, `ListBodyWeightEntries`, `LatestBodyWeightEntry`, `UpdateBodyWeightEntry`, `DeleteBodyWeightEntry` sqlc queries.

Include helpers:
```go
func bwRecordFromRow(row generated.AtlasBodyWeightEntry) *models.BodyWeightEntryRecord {
	return &models.BodyWeightEntryRecord{
		ID: row.ID.String(), UserID: row.UserID.String(),
		Date: row.Date.Time, Weight: float64(row.Weight),
		Source: row.Source, Notes: row.Notes,
		CreatedAt: row.CreatedAt.Time, UpdatedAt: row.UpdatedAt.Time,
	}
}
```

- [ ] **Step 3: Body check-in repo**

Create `apps/api/internal/atlas/repository/postgres/body_check_in_repo.go` — calls `CreateBodyCheckIn`, `GetBodyCheckInByID`, `ListBodyCheckInsByDateRange`, `ListBodyCheckIns`, `UpdateBodyCheckIn`, `DeleteBodyCheckIn`.

- [ ] **Step 4: Body measurement repo**

Create `apps/api/internal/atlas/repository/postgres/body_measurement_repo.go` — calls `CreateBodyMeasurement`, `GetBodyMeasurementByID`, `ListBodyMeasurementsByCheckIn`, `UpdateBodyMeasurement`, `DeleteBodyMeasurement`, `DeleteBodyMeasurementsByCheckIn`.

Note: No user_id filtering for measurements — they're scoped by check_in_id (which is already user-scoped).

- [ ] **Step 5: Progress photo repo**

Create `apps/api/internal/atlas/repository/postgres/progress_photo_repo.go` — calls `CreateProgressPhoto`, `GetProgressPhotoByID`, `ListProgressPhotosByCheckIn`, `DeleteProgressPhoto`.

- [ ] **Step 6: Week flag repo**

Create `apps/api/internal/atlas/repository/postgres/week_flag_repo.go` — calls `CreateWeekFlag`, `GetWeekFlagByID`, `ListWeekFlagsByWeek`, `DeleteWeekFlag`.

- [ ] **Step 7: Build check**

Run: `cd apps/api && go build ./internal/atlas/repository/postgres`
Expected: exits 0.

- [ ] **Step 8: Commit**

```bash
git add apps/api/internal/atlas/repository/postgres/cardio_entry_repo.go apps/api/internal/atlas/repository/postgres/body_weight_entry_repo.go apps/api/internal/atlas/repository/postgres/body_check_in_repo.go apps/api/internal/atlas/repository/postgres/body_measurement_repo.go apps/api/internal/atlas/repository/postgres/progress_photo_repo.go apps/api/internal/atlas/repository/postgres/week_flag_repo.go
git commit -m "feat(wave-04): add repository adapters"
```

---

## Task 5: Service Layer

**Files:**
- Create: `apps/api/internal/atlas/service/cardio.go`
- Create: `apps/api/internal/atlas/service/body_weight.go`
- Create: `apps/api/internal/atlas/service/body_checkin.go`
- Create: `apps/api/internal/atlas/service/week_flag.go`

- [ ] **Step 1: Cardio service**

Create `apps/api/internal/atlas/service/cardio.go`:

```go
package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasRepo "monorepo-template/apps/api/internal/atlas/repository/postgres"
	generated "monorepo-template/apps/api/internal/repository/postgres/generated"
	postgresRepo "monorepo-template/apps/api/internal/repository/postgres"
)

var (
	ErrCardioInvalidType      = errors.New("invalid cardio type")
	ErrCardioInvalidDuration  = errors.New("duration must be positive")
	ErrCardioInvalidPulse     = errors.New("pulse must be positive")
	ErrCardioInvalidZone      = errors.New("invalid heart rate zone")
	ErrCardioEntryNotFound    = errors.New("cardio entry not found")
)

type CardioService interface {
	Create(ctx context.Context, userID string, input models.CreateCardioInput) (*models.CardioEntry, error)
	GetByID(ctx context.Context, userID, id string) (*models.CardioEntry, error)
	ListByDailyLog(ctx context.Context, userID, dailyLogID string) ([]models.CardioEntry, error)
	Update(ctx context.Context, userID, id string, input models.UpdateCardioInput) (*models.CardioEntry, error)
	Delete(ctx context.Context, userID, id string) error
}

type cardioService struct {
	repo       atlasRepo.CardioEntryRepository
	pool       *pgxpool.Pool
}

func NewCardioService(repo atlasRepo.CardioEntryRepository, pool *pgxpool.Pool) CardioService {
	return &cardioService{repo: repo, pool: pool}
}

func (s *cardioService) Create(ctx context.Context, userID string, input models.CreateCardioInput) (*models.CardioEntry, error) {
	if !isValidCardioType(input.CardioType) {
		return nil, fmt.Errorf("%w: %s", ErrCardioInvalidType, input.CardioType)
	}
	if input.DurationMins <= 0 {
		return nil, ErrCardioInvalidDuration
	}
	if input.AvgPulse != nil && *input.AvgPulse <= 0 {
		return nil, ErrCardioInvalidPulse
	}
	if input.HeartRateZone != nil && !isValidHeartRateZone(*input.HeartRateZone) {
		return nil, fmt.Errorf("%w: %s", ErrCardioInvalidZone, *input.HeartRateZone)
	}

	// Auto-create DailyLog if needed
	dailyLogID := input.DailyLogID
	if dailyLogID == "" && input.Date != "" {
		dlRepo := postgresRepo.NewDailyLogRepo(s.pool)
		dl, err := dlRepo.GetOrCreateByDate(ctx, userID, input.Date)
		if err != nil {
			return nil, fmt.Errorf("CardioService.Create: auto-create DailyLog: %w", err)
		}
		dailyLogID = dl.ID
	}
	if dailyLogID == "" {
		return nil, fmt.Errorf("CardioService.Create: dailyLogId or date required")
	}

	record, err := s.repo.Create(ctx, userID, dailyLogID, input)
	if err != nil {
		return nil, fmt.Errorf("CardioService.Create: %w", err)
	}
	entry := models.CardioEntryFromRecord(record)
	return entry, nil
}

func (s *cardioService) GetByID(ctx context.Context, userID, id string) (*models.CardioEntry, error) {
	record, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCardioEntryNotFound
		}
		return nil, fmt.Errorf("CardioService.GetByID: %w", err)
	}
	return models.CardioEntryFromRecord(record), nil
}

func (s *cardioService) ListByDailyLog(ctx context.Context, userID, dailyLogID string) ([]models.CardioEntry, error) {
	records, err := s.repo.ListByDailyLog(ctx, userID, dailyLogID)
	if err != nil {
		return nil, fmt.Errorf("CardioService.ListByDailyLog: %w", err)
	}
	return models.CardioEntriesFromRecords(records), nil
}

func (s *cardioService) Update(ctx context.Context, userID, id string, input models.UpdateCardioInput) (*models.CardioEntry, error) {
	if input.CardioType != nil && !isValidCardioType(*input.CardioType) {
		return nil, fmt.Errorf("%w: %s", ErrCardioInvalidType, *input.CardioType)
	}
	if input.DurationMins != nil && *input.DurationMins <= 0 {
		return nil, ErrCardioInvalidDuration
	}
	if input.AvgPulse != nil && *input.AvgPulse <= 0 {
		return nil, ErrCardioInvalidPulse
	}
	if input.HeartRateZone != nil && !isValidHeartRateZone(*input.HeartRateZone) {
		return nil, fmt.Errorf("%w: %s", ErrCardioInvalidZone, *input.HeartRateZone)
	}
	record, err := s.repo.Update(ctx, userID, id, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCardioEntryNotFound
		}
		return nil, fmt.Errorf("CardioService.Update: %w", err)
	}
	return models.CardioEntryFromRecord(record), nil
}

func (s *cardioService) Delete(ctx context.Context, userID, id string) error {
	err := s.repo.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCardioEntryNotFound
		}
		return fmt.Errorf("CardioService.Delete: %w", err)
	}
	return nil
}

func isValidCardioType(t string) bool {
	for _, vt := range models.ValidCardioTypes {
		if string(vt) == t { return true }
	}
	return false
}

func isValidHeartRateZone(z string) bool {
	for _, vz := range models.ValidHeartRateZones {
		if string(vz) == z { return true }
	}
	return false
}
```

- [ ] **Step 2: Body weight service**

Create `apps/api/internal/atlas/service/body_weight.go` — implements `BodyWeightService` with validation (weight > 0, source enum). Methods: Create, GetByID, List (date range), Latest, Update, Delete.

```go
var (
	ErrBWInvalidWeight = errors.New("weight must be positive")
	ErrBWInvalidSource = errors.New("invalid weight source")
	ErrBWEntryNotFound = errors.New("body weight entry not found")
)
```

- [ ] **Step 3: Body check-in service**

Create `apps/api/internal/atlas/service/body_checkin.go` — implements `BodyCheckInService` with:
- Create: validate date, weight>0 if provided, bodyFatPct 0-100 if provided. Auto-create BodyWeightEntry if weight provided.
- GetByID: returns check-in with nested measurements + photos.
- List (date range): returns check-ins.
- Update: validate weight/fat pct.
- Delete: delete check-in (cascades via DB FK to measurements + photos).

- [ ] **Step 4: Week flag service**

Create `apps/api/internal/atlas/service/week_flag.go` — implements `WeekFlagService` with validation (flagType enum, unique per week).

```go
var (
	ErrWeekFlagInvalidType = errors.New("invalid week flag type")
	ErrWeekFlagDuplicate   = errors.New("flag type already exists for this week")
	ErrWeekFlagNotFound    = errors.New("week flag not found")
)
```

- [ ] **Step 5: Build check**

Run: `cd apps/api && go build ./internal/atlas/service`
Expected: exits 0.

- [ ] **Step 6: Commit**

```bash
git add apps/api/internal/atlas/service/cardio.go apps/api/internal/atlas/service/body_weight.go apps/api/internal/atlas/service/body_checkin.go apps/api/internal/atlas/service/week_flag.go
git commit -m "feat(wave-04): add service layer"
```

---

## Task 6: GraphQL Schema

**Files:**
- Create: `apps/api/internal/atlas/graph/schema/cardio.graphql`
- Create: `apps/api/internal/atlas/graph/schema/body.graphql`
- Create: `apps/api/internal/atlas/graph/schema/progress_photo.graphql`
- Create: `apps/api/internal/atlas/graph/schema/week_flag.graphql`
- Modify: `apps/api/internal/atlas/graph/schema/schema.graphql`

- [ ] **Step 1: Cardio GraphQL schema**

Create `apps/api/internal/atlas/graph/schema/cardio.graphql`:

```graphql
enum CardioType {
  WALKING
  RUNNING
  BIKE
  ELLIPTICAL
  TREADMILL
  OTHER
}

enum HeartRateZone {
  ZONE_1
  ZONE_2
  ZONE_3
  ZONE_4
  ZONE_5
  UNKNOWN
}

enum CardioErrorCode {
  VALIDATION_ERROR
  NOT_FOUND
  AUTH_ERROR
  INTERNAL_ERROR
}

type CardioEntry {
  id: ID!
  dailyLogId: ID!
  cardioType: CardioType!
  durationMinutes: Int!
  avgPulse: Int
  heartRateZone: HeartRateZone
  notes: String
  createdAt: Time!
  updatedAt: Time!
}

input CreateCardioInput {
  dailyLogId: ID
  date: String
  cardioType: CardioType!
  durationMinutes: Int!
  avgPulse: Int
  heartRateZone: HeartRateZone
  notes: String
}

input UpdateCardioInput {
  cardioType: CardioType
  durationMinutes: Int
  avgPulse: Int
  heartRateZone: HeartRateZone
  notes: String
}

type CardioEntryResult {
  cardioEntry: CardioEntry
  validationError: CardioValidationError
  notFoundError: CardioNotFoundError
  authError: CardioAuthError
}

type CardioDeleteResult {
  deleteResult: DeleteResult
  notFoundError: CardioNotFoundError
  authError: CardioAuthError
}

type DeleteResult {
  success: Boolean!
}

type CardioValidationError {
  message: String!
  code: CardioErrorCode!
}

type CardioNotFoundError {
  message: String!
  code: CardioErrorCode!
}

type CardioAuthError {
  message: String!
  code: CardioErrorCode!
}
```

- [ ] **Step 2: Body GraphQL schema**

Create `apps/api/internal/atlas/graph/schema/body.graphql`:

```graphql
enum MeasurementType {
  NECK
  SHOULDERS
  FOREARMS
  BICEPS
  CHEST
  WAIST
  ABDOMEN
  HIPS
  THIGH
  CALF
}

enum MeasurementSide {
  LEFT
  RIGHT
}

enum BodyWeightSource {
  SCALE
  MANUAL
  UNKNOWN
}

enum BodyErrorCode {
  VALIDATION_ERROR
  NOT_FOUND
  AUTH_ERROR
  INTERNAL_ERROR
}

type BodyWeightEntry {
  id: ID!
  date: String!
  weight: Float!
  source: BodyWeightSource!
  notes: String
  createdAt: Time!
  updatedAt: Time!
}

type BodyCheckIn {
  id: ID!
  date: String!
  weight: Float
  bodyFatPercentage: Float
  notes: String
  measurements: [BodyMeasurement!]!
  progressPhotos: [ProgressPhoto!]!
  createdAt: Time!
  updatedAt: Time!
}

type BodyMeasurement {
  id: ID!
  checkInId: ID!
  measurementType: MeasurementType!
  side: MeasurementSide
  value: Float!
  createdAt: Time!
  updatedAt: Time!
}

input CreateBodyWeightInput {
  date: String!
  weight: Float!
  source: BodyWeightSource!
  notes: String
}

input UpdateBodyWeightInput {
  weight: Float
  source: BodyWeightSource
  notes: String
}

input CreateCheckInInput {
  date: String!
  weight: Float
  bodyFatPercentage: Float
  notes: String
}

input UpdateCheckInInput {
  weight: Float
  bodyFatPercentage: Float
  notes: String
}

input CreateMeasurementInput {
  measurementType: MeasurementType!
  side: MeasurementSide
  value: Float!
}

input UpdateMeasurementInput {
  side: MeasurementSide
  value: Float
}

type BodyWeightResult {
  bodyWeightEntry: BodyWeightEntry
  validationError: BodyValidationError
  notFoundError: BodyNotFoundError
  authError: BodyAuthError
}

type BodyWeightDeleteResult {
  deleteResult: DeleteResult
  notFoundError: BodyNotFoundError
  authError: BodyAuthError
}

type BodyCheckInResult {
  bodyCheckIn: BodyCheckIn
  validationError: BodyValidationError
  notFoundError: BodyNotFoundError
  authError: BodyAuthError
}

type BodyCheckInDeleteResult {
  deleteResult: DeleteResult
  notFoundError: BodyNotFoundError
  authError: BodyAuthError
}

type MeasurementResult {
  measurement: BodyMeasurement
  validationError: BodyValidationError
  notFoundError: BodyNotFoundError
  authError: BodyAuthError
}

type MeasurementDeleteResult {
  deleteResult: DeleteResult
  notFoundError: BodyNotFoundError
  authError: BodyAuthError
}

type BodyValidationError {
  message: String!
  code: BodyErrorCode!
}

type BodyNotFoundError {
  message: String!
  code: BodyErrorCode!
}

type BodyAuthError {
  message: String!
  code: BodyErrorCode!
}
```

- [ ] **Step 3: Progress photo GraphQL schema**

Create `apps/api/internal/atlas/graph/schema/progress_photo.graphql`:

```graphql
enum PhotoAngle {
  FRONT
  SIDE
  BACK
  CUSTOM
}

type ProgressPhoto {
  id: ID!
  checkInId: ID!
  originalFileName: String!
  mimeType: String!
  sizeBytes: Int!
  angle: PhotoAngle
  label: String
  notes: String
  createdAt: Time!
  updatedAt: Time!
}
```

- [ ] **Step 4: Week flag GraphQL schema**

Create `apps/api/internal/atlas/graph/schema/week_flag.graphql`:

```graphql
enum FlagType {
  POOR_SLEEP
  HIGH_STRESS
  ILLNESS
  INJURY_PAIN
  CYCLE
  CALORIE_DEFICIT
  SURPLUS
  MAINTENANCE
  MISSED_WORKOUTS
  TRAVEL
}

enum WeekFlagErrorCode {
  VALIDATION_ERROR
  NOT_FOUND
  AUTH_ERROR
  INTERNAL_ERROR
}

type WeekFlag {
  id: ID!
  weekStartDate: String!
  flagType: FlagType!
  notes: String
  createdAt: Time!
  updatedAt: Time!
}

input CreateWeekFlagInput {
  weekStartDate: String!
  flagType: FlagType!
  notes: String
}

type WeekFlagResult {
  weekFlag: WeekFlag
  validationError: WeekFlagValidationError
  notFoundError: WeekFlagNotFoundError
  authError: WeekFlagAuthError
}

type WeekFlagDeleteResult {
  deleteResult: DeleteResult
  notFoundError: WeekFlagNotFoundError
  authError: WeekFlagAuthError
}

type WeekFlagValidationError {
  message: String!
  code: WeekFlagErrorCode!
}

type WeekFlagNotFoundError {
  message: String!
  code: WeekFlagErrorCode!
}

type WeekFlagAuthError {
  message: String!
  code: WeekFlagErrorCode!
}
```

- [ ] **Step 5: Update schema.graphql**

Edit `apps/api/internal/atlas/graph/schema/schema.graphql` to add WAVE-04 queries and mutations:

```graphql
scalar Time

type Query {
  # Existing
  settings: SettingsResult!
  exercises(first: Int = 20, after: String, includeInactive: Boolean = false): ExerciseConnection!
  exercise(id: ID!): ExerciseResult!
  allExercises(includeInactive: Boolean = false): [Exercise!]!
  # WAVE-04
  cardioEntries(dailyLogId: ID!): [CardioEntry!]!
  cardioEntry(id: ID!): CardioEntryResult!
  bodyWeightEntries(startDate: String!, endDate: String!): [BodyWeightEntry!]!
  bodyWeightEntry(id: ID!): BodyWeightResult!
  latestBodyWeight: BodyWeightResult!
  bodyCheckIns(startDate: String!, endDate: String!): [BodyCheckIn!]!
  bodyCheckIn(id: ID!): BodyCheckInResult!
  progressPhotos(checkInId: ID!): [ProgressPhoto!]!
  weekFlags(weekStartDate: String!): [WeekFlag!]!
}

type Mutation {
  # Existing
  updateSettings(input: SettingsInput!): SettingsResult!
  enablePin(input: PinEnableInput!): PinOperationResult!
  disablePin(input: PinDisableInput!): PinOperationResult!
  changePin(input: PinChangeInput!): PinOperationResult!
  createExercise(input: CreateExerciseInput!): ExerciseResult!
  updateExercise(id: ID!, input: UpdateExerciseInput!): ExerciseResult!
  archiveExercise(id: ID!): ArchiveResult!
  restoreExercise(id: ID!): ArchiveResult!
  # WAVE-04
  createCardioEntry(input: CreateCardioInput!): CardioEntryResult!
  updateCardioEntry(id: ID!, input: UpdateCardioInput!): CardioEntryResult!
  deleteCardioEntry(id: ID!): CardioDeleteResult!
  createBodyWeightEntry(input: CreateBodyWeightInput!): BodyWeightResult!
  updateBodyWeightEntry(id: ID!, input: UpdateBodyWeightInput!): BodyWeightResult!
  deleteBodyWeightEntry(id: ID!): BodyWeightDeleteResult!
  createBodyCheckIn(input: CreateCheckInInput!): BodyCheckInResult!
  updateBodyCheckIn(id: ID!, input: UpdateCheckInInput!): BodyCheckInResult!
  deleteBodyCheckIn(id: ID!): BodyCheckInDeleteResult!
  createMeasurement(checkInId: ID!, input: CreateMeasurementInput!): MeasurementResult!
  updateMeasurement(id: ID!, input: UpdateMeasurementInput!): MeasurementResult!
  deleteMeasurement(id: ID!): MeasurementDeleteResult!
  createWeekFlag(input: CreateWeekFlagInput!): WeekFlagResult!
  deleteWeekFlag(id: ID!): WeekFlagDeleteResult!
}
```

- [ ] **Step 6: Build check**

Run: `cd apps/api && go build ./internal/atlas/graph/schema`
Expected: exits 0 (schema files are embedded, but build should pass).

- [ ] **Step 7: Commit**

```bash
git add apps/api/internal/atlas/graph/schema/cardio.graphql apps/api/internal/atlas/graph/schema/body.graphql apps/api/internal/atlas/graph/schema/progress_photo.graphql apps/api/internal/atlas/graph/schema/week_flag.graphql apps/api/internal/atlas/graph/schema/schema.graphql
git commit -m "feat(wave-04): add GraphQL schemas"
```

---

## Task 7: gqlgen Config Update

**Files:**
- Modify: `apps/api/atlas-gqlgen.yml`

- [ ] **Step 1: Add WAVE-04 model bindings**

Edit `apps/api/atlas-gqlgen.yml`, add to `models:` section:

```yaml
  # WAVE-04
  CardioType:
    model: monorepo-template/apps/api/internal/atlas/models.CardioType
  HeartRateZone:
    model: monorepo-template/apps/api/internal/atlas/models.HeartRateZone
  CardioEntry:
    model: monorepo-template/apps/api/internal/atlas/models.CardioEntry
  CreateCardioInput:
    model: monorepo-template/apps/api/internal/atlas/models.CreateCardioInput
  UpdateCardioInput:
    model: monorepo-template/apps/api/internal/atlas/models.UpdateCardioInput
  CardioEntryResult:
    model: monorepo-template/apps/api/internal/atlas/models.CardioEntryResult
  CardioDeleteResult:
    model: monorepo-template/apps/api/internal/atlas/models.CardioDeleteResult
  DeleteResult:
    model: monorepo-template/apps/api/internal/atlas/models.DeleteResult
  CardioValidationError:
    model: monorepo-template/apps/api/internal/atlas/models.CardioErr
  CardioNotFoundError:
    model: monorepo-template/apps/api/internal/atlas/models.CardioErr
  CardioAuthError:
    model: monorepo-template/apps/api/internal/atlas/models.CardioErr
  CardioErrorCode:
    model: monorepo-template/apps/api/internal/atlas/models.CardioErrorCode
  MeasurementType:
    model: monorepo-template/apps/api/internal/atlas/models.MeasurementType
  MeasurementSide:
    model: monorepo-template/apps/api/internal/atlas/models.MeasurementSide
  BodyWeightSource:
    model: monorepo-template/apps/api/internal/atlas/models.BodyWeightSource
  BodyWeightEntry:
    model: monorepo-template/apps/api/internal/atlas/models.BodyWeightEntry
  BodyCheckIn:
    model: monorepo-template/apps/api/internal/atlas/models.BodyCheckIn
  BodyMeasurement:
    model: monorepo-template/apps/api/internal/atlas/models.BodyMeasurement
  ProgressPhoto:
    model: monorepo-template/apps/api/internal/atlas/models.ProgressPhoto
  CreateBodyWeightInput:
    model: monorepo-template/apps/api/internal/atlas/models.CreateBodyWeightInput
  UpdateBodyWeightInput:
    model: monorepo-template/apps/api/internal/atlas/models.UpdateBodyWeightInput
  CreateCheckInInput:
    model: monorepo-template/apps/api/internal/atlas/models.CreateCheckInInput
  UpdateCheckInInput:
    model: monorepo-template/apps/api/internal/atlas/models.UpdateCheckInInput
  CreateMeasurementInput:
    model: monorepo-template/apps/api/internal/atlas/models.CreateMeasurementInput
  UpdateMeasurementInput:
    model: monorepo-template/apps/api/internal/atlas/models.UpdateMeasurementInput
  BodyWeightResult:
    model: monorepo-template/apps/api/internal/atlas/models.BodyWeightResult
  BodyWeightDeleteResult:
    model: monorepo-template/apps/api/internal/atlas/models.BodyWeightDeleteResult
  BodyCheckInResult:
    model: monorepo-template/apps/api/internal/atlas/models.BodyCheckInResult
  BodyCheckInDeleteResult:
    model: monorepo-template/apps/api/internal/atlas/models.BodyCheckInDeleteResult
  MeasurementResult:
    model: monorepo-template/apps/api/internal/atlas/models.MeasurementResult
  MeasurementDeleteResult:
    model: monorepo-template/apps/api/internal/atlas/models.MeasurementDeleteResult
  BodyValidationError:
    model: monorepo-template/apps/api/internal/atlas/models.BodyErr
  BodyNotFoundError:
    model: monorepo-template/apps/api/internal/atlas/models.BodyErr
  BodyAuthError:
    model: monorepo-template/apps/api/internal/atlas/models.BodyErr
  BodyErrorCode:
    model: monorepo-template/apps/api/internal/atlas/models.BodyErrorCode
  PhotoAngle:
    model: monorepo-template/apps/api/internal/atlas/models.PhotoAngle
  FlagType:
    model: monorepo-template/apps/api/internal/atlas/models.FlagType
  WeekFlag:
    model: monorepo-template/apps/api/internal/atlas/models.WeekFlag
  CreateWeekFlagInput:
    model: monorepo-template/apps/api/internal/atlas/models.CreateWeekFlagInput
  WeekFlagResult:
    model: monorepo-template/apps/api/internal/atlas/models.WeekFlagResult
  WeekFlagDeleteResult:
    model: monorepo-template/apps/api/internal/atlas/models.WeekFlagDeleteResult
  WeekFlagValidationError:
    model: monorepo-template/apps/api/internal/atlas/models.WeekFlagErr
  WeekFlagNotFoundError:
    model: monorepo-template/apps/api/internal/atlas/models.WeekFlagErr
  WeekFlagAuthError:
    model: monorepo-template/apps/api/internal/atlas/models.WeekFlagErr
  WeekFlagErrorCode:
    model: monorepo-template/apps/api/internal/atlas/models.WeekFlagErrorCode
```

- [ ] **Step 2: Run gqlgen codegen**

Run: `cd apps/api && go run github.com/99designs/gqlgen generate -c atlas-gqlgen.yml`
Expected: exits 0, generates files under `internal/atlas/graph/generated/` and resolver stubs under `internal/atlas/graph/resolver/`.

- [ ] **Step 3: Build check**

Run: `cd apps/api && go build ./internal/atlas/graph/generated`
Expected: exits 0.

- [ ] **Step 4: Commit**

```bash
git add apps/api/atlas-gqlgen.yml apps/api/internal/atlas/graph/generated/
git commit -m "feat(wave-04): update gqlgen config and regenerate"
```

---

## Task 8: GraphQL Resolvers

**Files:**
- Create: `apps/api/internal/atlas/graph/resolver/cardio.go`
- Create: `apps/api/internal/atlas/graph/resolver/body.go`
- Create: `apps/api/internal/atlas/graph/resolver/photo.go`
- Create: `apps/api/internal/atlas/graph/resolver/week_flag.go`
- Modify: `apps/api/internal/atlas/graph/resolver/resolver.go`

- [ ] **Step 1: Update resolver.go with WAVE-04 services**

Edit `apps/api/internal/atlas/graph/resolver/resolver.go`:

```go
package resolver

import "monorepo-template/apps/api/internal/atlas/service"

type Resolver struct {
	SettingsService  service.SettingsService
	PinService       service.PinService
	ExerciseService  service.ExerciseService
	CardioService    service.CardioService
	BodyWeightService  service.BodyWeightService
	BodyCheckInService service.BodyCheckInService
	WeekFlagService    service.WeekFlagService
}
```

- [ ] **Step 2: Cardio resolvers**

Create `apps/api/internal/atlas/graph/resolver/cardio.go`:

```go
package resolver

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/service"
)

func (r *Resolver) CardioEntries(ctx context.Context, dailyLogID string) ([]models.CardioEntry, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" { return nil, nil }
	entries, err := r.CardioService.ListByDailyLog(ctx, userID, dailyLogID)
	if err != nil { return nil, nil }
	return entries, nil
}

func (r *Resolver) CardioEntry(ctx context.Context, id string) (*models.CardioEntryResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.CardioEntryResult{AuthErr: &models.CardioErr{Message: "unauthorized", Code: models.CardioErrAuth}}, nil
	}
	entry, err := r.CardioService.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, service.ErrCardioEntryNotFound) {
			return &models.CardioEntryResult{NotFoundErr: &models.CardioErr{Message: err.Error(), Code: models.CardioErrNotFound}}, nil
		}
		return nil, nil
	}
	return &models.CardioEntryResult{CardioEntry: entry}, nil
}

func (r *Resolver) CreateCardioEntry(ctx context.Context, input models.CreateCardioInput) (*models.CardioEntryResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.CardioEntryResult{AuthErr: &models.CardioErr{Message: "unauthorized", Code: models.CardioErrAuth}}, nil
	}
	entry, err := r.CardioService.Create(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCardioInvalidType), errors.Is(err, service.ErrCardioInvalidDuration),
			errors.Is(err, service.ErrCardioInvalidPulse), errors.Is(err, service.ErrCardioInvalidZone):
			return &models.CardioEntryResult{ValidationErr: &models.CardioErr{Message: err.Error(), Code: models.CardioErrValidation}}, nil
		default:
			return nil, nil
		}
	}
	return &models.CardioEntryResult{CardioEntry: entry}, nil
}

func (r *Resolver) UpdateCardioEntry(ctx context.Context, id string, input models.UpdateCardioInput) (*models.CardioEntryResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.CardioEntryResult{AuthErr: &models.CardioErr{Message: "unauthorized", Code: models.CardioErrAuth}}, nil
	}
	entry, err := r.CardioService.Update(ctx, userID, id, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCardioEntryNotFound):
			return &models.CardioEntryResult{NotFoundErr: &models.CardioErr{Message: err.Error(), Code: models.CardioErrNotFound}}, nil
		case errors.Is(err, service.ErrCardioInvalidType), errors.Is(err, service.ErrCardioInvalidDuration),
			errors.Is(err, service.ErrCardioInvalidPulse), errors.Is(err, service.ErrCardioInvalidZone):
			return &models.CardioEntryResult{ValidationErr: &models.CardioErr{Message: err.Error(), Code: models.CardioErrValidation}}, nil
		default:
			return nil, nil
		}
	}
	return &models.CardioEntryResult{CardioEntry: entry}, nil
}

func (r *Resolver) DeleteCardioEntry(ctx context.Context, id string) (*models.CardioDeleteResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.CardioDeleteResult{AuthErr: &models.CardioErr{Message: "unauthorized", Code: models.CardioErrAuth}}, nil
	}
	err := r.CardioService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, service.ErrCardioEntryNotFound) {
			return &models.CardioDeleteResult{NotFoundErr: &models.CardioErr{Message: err.Error(), Code: models.CardioErrNotFound}}, nil
		}
		return nil, nil
	}
	return &models.CardioDeleteResult{DeleteResult: &models.DeleteResult{Success: true}}, nil
}
```

- [ ] **Step 3: Body resolvers (body weight + check-in + measurement)**

Create `apps/api/internal/atlas/graph/resolver/body.go` — implement all body tracker resolvers following same pattern as cardio: get userID from middleware, call service, map errors to result unions.

Methods to implement:
- `BodyWeightEntries`, `BodyWeightEntry`, `LatestBodyWeight`, `CreateBodyWeightEntry`, `UpdateBodyWeightEntry`, `DeleteBodyWeightEntry`
- `BodyCheckIns`, `BodyCheckIn`, `CreateBodyCheckIn`, `UpdateBodyCheckIn`, `DeleteBodyCheckIn`
- `CreateMeasurement`, `UpdateMeasurement`, `DeleteMeasurement`

- [ ] **Step 4: Photo resolvers (query only)**

Create `apps/api/internal/atlas/graph/resolver/photo.go` — implement `ProgressPhotos(ctx, checkInID)` that delegates to `BodyCheckInService.ListPhotos(ctx, checkInID)`.

- [ ] **Step 5: Week flag resolvers**

Create `apps/api/internal/atlas/graph/resolver/week_flag.go` — implement `WeekFlags`, `CreateWeekFlag`, `DeleteWeekFlag`.

- [ ] **Step 6: Build check**

Run: `cd apps/api && go build ./internal/atlas/graph/resolver`
Expected: exits 0.

- [ ] **Step 7: Commit**

```bash
git add apps/api/internal/atlas/graph/resolver/resolver.go apps/api/internal/atlas/graph/resolver/cardio.go apps/api/internal/atlas/graph/resolver/body.go apps/api/internal/atlas/graph/resolver/photo.go apps/api/internal/atlas/graph/resolver/week_flag.go
git commit -m "feat(wave-04): add GraphQL resolvers"
```

---

## Task 9: ProgressPhoto REST Handler

**Files:**
- Create: `apps/api/internal/atlas/handler/photo_handler.go`

- [ ] **Step 1: Create photo REST handler**

Create `apps/api/internal/atlas/handler/photo_handler.go`:

```go
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"mime/multipart"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
)

const (
	maxPhotoSize   = 25 << 20 // 25MB
	photosPerCheckIn = 10
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

type PhotoHandler struct {
	photoRepo   atlasPostgres.ProgressPhotoRepository
	checkInRepo atlasPostgres.BodyCheckInRepository
	basePath    string
	logger      *zap.Logger
}

func NewPhotoHandler(photoRepo atlasPostgres.ProgressPhotoRepository, checkInRepo atlasPostgres.BodyCheckInRepository, basePath string, logger *zap.Logger) *PhotoHandler {
	return &PhotoHandler{photoRepo: photoRepo, checkInRepo: checkInRepo, basePath: basePath, logger: logger}
}

func (h *PhotoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writePhotoError(w, http.StatusUnauthorized, "UNAUTHORIZED", "PIN session required")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxPhotoSize)
	if err := r.ParseMultipartForm(maxPhotoSize); err != nil {
		writePhotoError(w, http.StatusBadRequest, "FILE_TOO_LARGE", "File exceeds 25MB limit")
		return
	}

	checkInID := r.FormValue("checkInId")
	if checkInID == "" {
		writePhotoError(w, http.StatusBadRequest, "VALIDATION_ERROR", "checkInId is required")
		return
	}

	// Verify check-in exists and belongs to user
	checkIn, err := h.checkInRepo.GetByID(r.Context(), userID, checkInID)
	if err != nil || checkIn == nil {
		writePhotoError(w, http.StatusNotFound, "NOT_FOUND", "Check-in not found")
		return
	}

	// Enforce max photos per check-in
	existing, err := h.photoRepo.ListByCheckIn(r.Context(), checkInID)
	if err == nil && len(existing) >= photosPerCheckIn {
		writePhotoError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Maximum 10 photos per check-in")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writePhotoError(w, http.StatusBadRequest, "VALIDATION_ERROR", "File is required")
		return
	}
	defer file.Close()

	// MIME validation
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		writePhotoError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to read file")
		return
	}
	file.Seek(0, io.SeekStart)

	mimeType := http.DetectContentType(buf)
	if !allowedMimeTypes[mimeType] {
		writePhotoError(w, http.StatusBadRequest, "INVALID_FILE_TYPE", "Only JPEG, PNG, and WEBP images are allowed")
		return
	}

	// Storage path
	photoID := uuid.New().String()
	ext := filepath.Ext(header.Filename)
	storageDir := filepath.Join(h.basePath, "progress-photos", checkInID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		h.logger.Error("[ProgressPhoto][upload] failed to create directory", zap.Error(err))
		writePhotoError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Storage error")
		return
	}
	storagePath := filepath.Join(storageDir, photoID+ext)

	dst, err := os.Create(storagePath)
	if err != nil {
		h.logger.Error("[ProgressPhoto][upload] failed to create file", zap.Error(err))
		writePhotoError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Storage error")
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		os.Remove(storagePath)
		h.logger.Error("[ProgressPhoto][upload] failed to write file", zap.Error(err))
		writePhotoError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Storage error")
		return
	}

	angle := r.FormValue("angle")
	label := r.FormValue("label")
	notes := r.FormValue("notes")

	record, err := h.photoRepo.Create(r.Context(), models.CreateProgressPhotoParams{
		CheckInID: checkInID, FilePath: storagePath,
		OriginalFileName: header.Filename, MimeType: mimeType,
		SizeBytes: written, Angle: stringPtr(angle), Label: stringPtr(label),
		Notes: stringPtr(notes),
	})
	if err != nil {
		os.Remove(storagePath)
		h.logger.Error("[ProgressPhoto][upload] failed to save record", zap.Error(err))
		writePhotoError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Database error")
		return
	}

	h.logger.Info("[ProgressPhoto][upload] success", zap.String("photo_id", record.ID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.PhotoFromRecord(record))
}

func (h *PhotoHandler) Download(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writePhotoError(w, http.StatusUnauthorized, "UNAUTHORIZED", "PIN session required")
		return
	}

	photoID := chi.URLParam(r, "id")
	record, err := h.photoRepo.GetByID(r.Context(), photoID)
	if err != nil || record == nil {
		writePhotoError(w, http.StatusNotFound, "NOT_FOUND", "Photo not found")
		return
	}

	w.Header().Set("Content-Type", record.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, record.OriginalFileName))
	http.ServeFile(w, r, record.FilePath)
}

func (h *PhotoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetAtlasUserID(r.Context())
	if userID == "" {
		writePhotoError(w, http.StatusUnauthorized, "UNAUTHORIZED", "PIN session required")
		return
	}

	photoID := chi.URLParam(r, "id")
	record, err := h.photoRepo.GetByID(r.Context(), photoID)
	if err != nil || record == nil {
		writePhotoError(w, http.StatusNotFound, "NOT_FOUND", "Photo not found")
		return
	}

	if err := h.photoRepo.Delete(r.Context(), photoID); err != nil {
		h.logger.Error("[ProgressPhoto][delete] failed to delete record", zap.Error(err))
		writePhotoError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Database error")
		return
	}

	if err := os.Remove(record.FilePath); err != nil {
		h.logger.Error("[ProgressPhoto][delete] failed to delete file", zap.String("path", record.FilePath), zap.Error(err))
	}

	h.logger.Info("[ProgressPhoto][delete] success", zap.String("photo_id", photoID))
	w.WriteHeader(http.StatusNoContent)
}

func writePhotoError(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{"code": code, "message": msg},
	})
}

func stringPtr(s string) *string {
	if s == "" { return nil }
	return &s
}
```

- [ ] **Step 2: Build check**

Run: `cd apps/api && go build ./internal/atlas/handler`
Expected: exits 0.

- [ ] **Step 3: Commit**

```bash
git add apps/api/internal/atlas/handler/photo_handler.go
git commit -m "feat(wave-04): add ProgressPhoto REST handler"
```

---

## Task 10: Main Wiring and Route Registration

**Files:**
- Modify: `apps/api/cmd/server/main.go`

- [ ] **Step 1: Wire all WAVE-04 repos, services, and photo handler**

Edit `apps/api/cmd/server/main.go`:

After the existing atlasExerciseRepo/atlasExerciseService block (around line 152), add:

```go
	// WAVE-04 repos
	atlasCardioRepo := atlasPostgres.NewCardioEntryRepository(db.Pool)
	atlasBWRepo := atlasPostgres.NewBodyWeightEntryRepository(db.Pool)
	atlasCheckInRepo := atlasPostgres.NewBodyCheckInRepository(db.Pool)
	atlasMeasRepo := atlasPostgres.NewBodyMeasurementRepository(db.Pool)
	atlasPhotoRepo := atlasPostgres.NewProgressPhotoRepository(db.Pool)
	atlasWeekFlagRepo := atlasPostgres.NewWeekFlagRepository(db.Pool)

	// WAVE-04 services
	atlasCardioService := atlasService.NewCardioService(atlasCardioRepo, db.Pool)
	atlasBWService := atlasService.NewBodyWeightService(atlasBWRepo)
	atlasCheckInService := atlasService.NewBodyCheckInService(atlasCheckInRepo, atlasMeasRepo, atlasPhotoRepo, atlasBWRepo, db.Pool)
	atlasWeekFlagService := atlasService.NewWeekFlagService(atlasWeekFlagRepo)

	atlasPhotoHandler := handler.NewPhotoHandler(atlasPhotoRepo, atlasCheckInRepo, cfg.Media.BasePath, l)
```

Update the resolver to include all new services:

```go
	atlasRes := &atlasResolver.Resolver{
		SettingsService:   atlasSettingsService,
		PinService:        atlasPinService,
		ExerciseService:   atlasExerciseService,
		CardioService:     atlasCardioService,
		BodyWeightService: atlasBWService,
		BodyCheckInService: atlasCheckInService,
		WeekFlagService:   atlasWeekFlagService,
	}
```

Add photo routes to the guarded group (around line 243):

```go
		atlas.Post("/api/v1/progress-photos/upload", atlasPhotoHandler.Upload)
		atlas.Get("/api/v1/progress-photos/{id}", atlasPhotoHandler.Download)
		atlas.Delete("/api/v1/progress-photos/{id}", atlasPhotoHandler.Delete)
```

- [ ] **Step 2: Add imports**

Ensure the imports block includes:
```go
	atlasHandler "monorepo-template/apps/api/internal/atlas/handler"
```

- [ ] **Step 3: Full build check**

Run: `cd apps/api && go build ./cmd/server`
Expected: exits 0.

- [ ] **Step 4: Commit**

```bash
git add apps/api/cmd/server/main.go
git commit -m "feat(wave-04): wire repos, services, resolvers, photo handler in main.go"
```

---

## Task 11: Tests

**Files:**
- Create: `apps/api/internal/repository/postgres/cardio_entry_repo_test.go`
- Create: `apps/api/internal/repository/postgres/body_weight_entry_repo_test.go`
- Create: `apps/api/internal/repository/postgres/body_check_in_repo_test.go`
- Create: `apps/api/internal/repository/postgres/body_measurement_repo_test.go`
- Create: `apps/api/internal/repository/postgres/progress_photo_repo_test.go`
- Create: `apps/api/internal/repository/postgres/week_flag_repo_test.go`
- Create: `apps/api/internal/atlas/service/cardio_service_test.go`
- Create: `apps/api/internal/atlas/service/body_weight_service_test.go`
- Create: `apps/api/internal/atlas/service/body_checkin_service_test.go`
- Create: `apps/api/internal/atlas/service/week_flag_service_test.go`
- Create: `apps/api/internal/atlas/graph/resolver/cardio_test.go`
- Create: `apps/api/internal/atlas/graph/resolver/body_test.go`
- Create: `apps/api/internal/atlas/graph/resolver/week_flag_test.go`
- Create: `apps/api/internal/atlas/handler/photo_handler_test.go`

- [ ] **Step 1: CardioEntry repo test**

Create `apps/api/internal/repository/postgres/cardio_entry_repo_test.go` following WAVE-02 exercise_repo_test.go pattern:

```go
package postgres_test

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	atlasRepo "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/apps/api/internal/atlas/models"
	postgresRepo "monorepo-template/apps/api/internal/repository/postgres"
	"monorepo-template/apps/api/internal/testinfra"
)

func TestCardioEntryRepo_Create_Success(t *testing.T) {
	pool, userID := cardioTestSetup(t)
	repo := postgresRepo.NewCardioEntryRepository(pool)
	
	input := models.CreateCardioInput{
		CardioType: "running", DurationMins: 30,
		AvgPulse: intPtr(145), HeartRateZone: strPtr("zone_3"),
		Notes: strPtr("morning run"),
	}
	
	record, err := repo.Create(context.Background(), userID, createDailyLog(t, pool, userID), input)
	require.NoError(t, err)
	assert.NotEmpty(t, record.ID)
	assert.Equal(t, "running", record.CardioType)
	assert.Equal(t, 30, record.DurationMins)
}

func TestCardioEntryRepo_GetByID_NotFound(t *testing.T) { ... }
func TestCardioEntryRepo_ListByDailyLog_Multiple(t *testing.T) { ... }
func TestCardioEntryRepo_Update_Success(t *testing.T) { ... }
func TestCardioEntryRepo_Delete_Success(t *testing.T) { ... }
```

- [ ] **Step 2: BodyWeightEntry repo test**

Create `apps/api/internal/repository/postgres/body_weight_entry_repo_test.go` — tests for CRUD, date range listing, latest weight.

- [ ] **Step 3: BodyCheckIn repo test**

Create `apps/api/internal/repository/postgres/body_check_in_repo_test.go` — tests for CRUD, date range, unique user+date constraint.

- [ ] **Step 4: BodyMeasurement repo test**

Create `apps/api/internal/repository/postgres/body_measurement_repo_test.go` — tests for CRUD, list by check-in, unique constraint (checkIn+type+side).

- [ ] **Step 5: ProgressPhoto repo test**

Create `apps/api/internal/repository/postgres/progress_photo_repo_test.go` — tests for create, get by ID, list by check-in, delete.

- [ ] **Step 6: WeekFlag repo test**

Create `apps/api/internal/repository/postgres/week_flag_repo_test.go` — tests for create, list by week, delete, unique constraint.

- [ ] **Step 7: Cardio service test**

Create `apps/api/internal/atlas/service/cardio_service_test.go` — mock repo, test validation (invalid type, duration <= 0, invalid pulse, invalid zone), success path.

- [ ] **Step 8: Body weight service test**

Create `apps/api/internal/atlas/service/body_weight_service_test.go` — mock repo, test validation (weight <= 0, invalid source), success.

- [ ] **Step 9: Body check-in service test**

Create `apps/api/internal/atlas/service/body_checkin_service_test.go` — mock repos, test validation (weight <= 0, body fat % out of range), success.

- [ ] **Step 10: Week flag service test**

Create `apps/api/internal/atlas/service/week_flag_service_test.go` — mock repo, test validation (invalid flag type), success.

- [ ] **Step 11: Resolver tests**

Create resolver tests for cardio, body, and week_flag — mock services, test auth error, success paths.

- [ ] **Step 12: Photo handler test**

Create `apps/api/internal/atlas/handler/photo_handler_test.go` — test upload (success, invalid MIME, too large, missing checkInId), download (success, not found), delete (success, not found).

- [ ] **Step 13: Run cardio repo tests**

Run: `cd apps/api && go test ./internal/repository/postgres -run TestCardioEntryRepo -count=1 -v`
Expected: PASS.

- [ ] **Step 14: Run body weight repo tests**

Run: `cd apps/api && go test ./internal/repository/postgres -run TestBodyWeightEntryRepo -count=1 -v`
Expected: PASS.

- [ ] **Step 15: Run body check-in repo tests**

Run: `cd apps/api && go test ./internal/repository/postgres -run TestBodyCheckInRepo -count=1 -v`
Expected: PASS.

- [ ] **Step 16: Run measurement repo tests**

Run: `cd apps/api && go test ./internal/repository/postgres -run TestBodyMeasurementRepo -count=1 -v`
Expected: PASS.

- [ ] **Step 17: Run photo repo tests**

Run: `cd apps/api && go test ./internal/repository/postgres -run TestProgressPhotoRepo -count=1 -v`
Expected: PASS.

- [ ] **Step 18: Run week flag repo tests**

Run: `cd apps/api && go test ./internal/repository/postgres -run TestWeekFlagRepo -count=1 -v`
Expected: PASS.

- [ ] **Step 19: Run service tests**

Run: `cd apps/api && go test ./internal/atlas/service -run "TestCardioService|TestBodyWeightService|TestBodyCheckInService|TestWeekFlagService" -count=1 -v`
Expected: PASS.

- [ ] **Step 20: Run resolver tests**

Run: `cd apps/api && go test ./internal/atlas/graph/resolver -run "TestCardio|TestBodyWeight|TestBodyCheckIn|TestWeekFlag" -count=1 -v`
Expected: PASS.

- [ ] **Step 21: Run handler tests**

Run: `cd apps/api && go test ./internal/atlas/handler -run TestPhotoHandler -count=1 -v`
Expected: PASS.

- [ ] **Step 22: Commit all tests**

```bash
git add apps/api/internal/repository/postgres/cardio_entry_repo_test.go apps/api/internal/repository/postgres/body_weight_entry_repo_test.go apps/api/internal/repository/postgres/body_check_in_repo_test.go apps/api/internal/repository/postgres/body_measurement_repo_test.go apps/api/internal/repository/postgres/progress_photo_repo_test.go apps/api/internal/repository/postgres/week_flag_repo_test.go apps/api/internal/atlas/service/cardio_service_test.go apps/api/internal/atlas/service/body_weight_service_test.go apps/api/internal/atlas/service/body_checkin_service_test.go apps/api/internal/atlas/service/week_flag_service_test.go apps/api/internal/atlas/graph/resolver/cardio_test.go apps/api/internal/atlas/graph/resolver/body_test.go apps/api/internal/atlas/graph/resolver/week_flag_test.go apps/api/internal/atlas/handler/photo_handler_test.go
git commit -m "test(wave-04): add repository, service, resolver, and handler tests"
```

---

## Task 12: Update GRACE Artifacts

**Files:**
- Modify: `docs/development-plan.xml`
- Modify: `docs/knowledge-graph.xml`
- Modify: `docs/verification-plan.xml`
- Modify: `docs/prd-waves/waves/wave-04.md`

- [ ] **Step 1: Update wave-04.md status**

Change `docs/prd-waves/waves/wave-04.md` status from `user-approved` to `implemented`.

- [ ] **Step 2: Update development-plan.xml**

Add WAVE-04 module entries for all new files.

- [ ] **Step 3: Update knowledge-graph.xml**

Add WAVE-04 paths, module annotations, and dependency edges.

- [ ] **Step 4: Update verification-plan.xml**

Add WAVE-04 verification entries from the detailed brief.

- [ ] **Step 5: Commit**

```bash
git add docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml docs/prd-waves/waves/wave-04.md
git commit -m "docs(wave-04): update GRACE artifacts after implementation"
```

---

## Task 13: Final Build and Codegen Verification

- [ ] **Step 1: Full build**

Run: `cd apps/api && go build ./...`
Expected: exits 0.

- [ ] **Step 2: Codegen drift check**

Run: `cd apps/api && go run github.com/99designs/gqlgen generate -c atlas-gqlgen.yml && sqlc generate`
Expected: no changes (already generated).

- [ ] **Step 3: Full test run**

Run: `cd apps/api && go test ./internal/... ./internal/atlas/... -count=1`
Expected: all tests PASS.

- [ ] **Step 4: Lint**

Run: `bunx nx run api:lint`
Expected: PASS.

- [ ] **Step 5: Final commit**

```bash
git add -A && git commit -m "chore(wave-04): final build, codegen, and test verification"
```