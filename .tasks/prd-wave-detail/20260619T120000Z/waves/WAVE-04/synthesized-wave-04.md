# WAVE-04: Cardio and Body Tracking

## Status
ready-for-dev

## User Approval
user-approved (2026-06-18). Ready for implementation.

## Source Wave Summary
WAVE-04 from docs/prd-waves/waves/wave-04.md. Cardio tracking and body measurements with weekly check-ins and progress photos. Source status: user-approved (2026-06-18).

## Outcome After Implementation
- OUT-W04-001: Cardio entries with type, duration, pulse, zone, attached to DailyLog
- OUT-W04-002: Body weight entries (standalone by date, source enum)
- OUT-W04-003: Weekly body check-ins with weight, body fat %, and nested measurements
- OUT-W04-004: Body measurements (10 types, paired left/right for applicable types)
- OUT-W04-005: Progress photos with angles, attached to check-ins
- OUT-W04-006: Week flags for AI export context

## Scope Included
- CAP-W04-001: CardioEntry CRUD via GraphQL (type, duration, avg pulse, heart rate zone, notes). Auto-create DailyLog when needed. Linked to daily log date.
- CAP-W04-002: BodyWeightEntry CRUD via GraphQL (date, weight, source enum, notes). Standalone per date. Latest weight query for dashboard.
- CAP-W04-003: BodyCheckIn CRUD via GraphQL (date, optional weight, optional bodyFatPercentage, notes). Nested BodyMeasurement and ProgressPhoto children.
- CAP-W04-004: BodyMeasurement CRUD via GraphQL (measurementType enum, side left/right/null for paired, value). 10 measurement types. Side validation: paired types only (forearm, biceps, thigh, calf).
- CAP-W04-005: ProgressPhoto CRUD via REST (multipart upload, download, delete). Angle enum (front/side/back/custom). File type/size validation. Physical file delete on record deletion.
- CAP-W04-006: WeekFlag CRUD via GraphQL (weekStartDate, flagType enum, notes). One flag per type per week.

## Scope Excluded
- Photo taken in app (upload only)
- Chart visualization (WAVE-06)
- AI export (WAVE-07)
- Photo count enforcement (2-4 recommended, soft guidance, hard limit at 10)
- Frontend pages, UI, UX, routes, navigation, components

## Dependencies And Other-Wave Fit
- WAVE-01 (Foundation): prerequisite — provides PIN auth middleware, media storage scaffold (BasePath, MaxUploadSize), migration infrastructure, GraphQL common types, codegen config, config extension pattern. WAVE-04 cannot start until WAVE-01 provides these contracts.
- WAVE-02 (Exercise Library): no direct dependency — can parallelize
- WAVE-03 (Workout Diary): partial dependency — CardioEntry FK references daily_log table created in WAVE-03. If WAVE-04 deploys before WAVE-03, WAVE-04 must create its own daily_log migration stub or defer cardio creation (see DDEC-W04-001). WAVE-04 can otherwise be authored in parallel with WAVE-03.
- WAVE-05 (Nutrition): no direct dependency — can fully parallelize
- WAVE-06 (Charts): WAVE-04 provides body weight, measurement, and check-in data for chart queries
- WAVE-07/08 (AI Export/Review): WAVE-04 provides cardio, weight, check-in, measurement, photo, and week flag data via service layer
- WAVE-09 (Backup): WAVE-04 tables are designed for JSON-serializable export compatibility

## Frontend Pages Dependencies
- PAGE-004 (Cardio): primary frontend consumer — depends on all CardioEntry GraphQL queries and mutations. Date-based listing for cardio log.
- PAGE-005 (Body Measurements): primary frontend consumer — depends on BodyCheckIn, BodyMeasurement, BodyWeightEntry GraphQL queries and mutations. Nested measurements and photos via check-in query.
- PAGE-006 (Progress Photos): depends on progressPhotos GraphQL query, REST upload/download/delete endpoints. Photos grouped by check-in.
- PAGE-001 (Dashboard): depends on latestBodyWeight GraphQL query for weight summary card.
- Dependency context only; no frontend pages, UI, or UX work in this wave.

## Codebase Fit And Touchpoints
- apps/api/internal/repository/postgres/migrations/00082_cardio_entries.sql: new migration
- apps/api/internal/repository/postgres/migrations/00083_body_weight_entries.sql: new migration
- apps/api/internal/repository/postgres/migrations/00084_body_check_ins.sql: new migration
- apps/api/internal/repository/postgres/migrations/00085_body_measurements.sql: new migration
- apps/api/internal/repository/postgres/migrations/00086_progress_photos.sql: new migration
- apps/api/internal/repository/postgres/migrations/00087_week_flags.sql: new migration
- apps/api/internal/repository/postgres/queries/cardio_entries.sql: sqlc query definitions
- apps/api/internal/repository/postgres/queries/body_weight_entries.sql: sqlc query definitions
- apps/api/internal/repository/postgres/queries/body_check_ins.sql: sqlc query definitions
- apps/api/internal/repository/postgres/queries/body_measurements.sql: sqlc query definitions
- apps/api/internal/repository/postgres/queries/progress_photos.sql: sqlc query definitions
- apps/api/internal/repository/postgres/queries/week_flags.sql: sqlc query definitions
- apps/api/internal/repository/postgres/cardio_entry_repo.go: repository adapter
- apps/api/internal/repository/postgres/body_weight_entry_repo.go: repository adapter
- apps/api/internal/repository/postgres/body_check_in_repo.go: repository adapter
- apps/api/internal/repository/postgres/body_measurement_repo.go: repository adapter
- apps/api/internal/repository/postgres/progress_photo_repo.go: repository adapter
- apps/api/internal/repository/postgres/week_flag_repo.go: repository adapter
- apps/api/internal/service/cardio.go: transport-neutral service with validation and DailyLog auto-creation
- apps/api/internal/service/body_weight.go: transport-neutral service with validation
- apps/api/internal/service/body_checkin.go: transport-neutral service with validation and cascade logic
- apps/api/internal/service/week_flag.go: transport-neutral service with validation
- apps/api/internal/graph/cardio.resolvers.go: GraphQL resolvers for cardio CRUD
- apps/api/internal/graph/body_weight.resolvers.go: GraphQL resolvers for body weight CRUD
- apps/api/internal/graph/body_checkin.resolvers.go: GraphQL resolvers for check-in, measurement, photo queries
- apps/api/internal/graph/week_flag.resolvers.go: GraphQL resolvers for week flag CRUD
- apps/api/internal/handler/progress_photo_handler.go: REST handler for upload/download/delete
- libs/graphql/schema/cardio.graphql: cardio GraphQL types and operations
- libs/graphql/schema/body_weight.graphql: body weight GraphQL types and operations
- libs/graphql/schema/body_checkin.graphql: check-in, measurement, photo GraphQL types and operations
- libs/graphql/schema/week_flag.graphql: week flag GraphQL types and operations
- apps/api/cmd/server/main.go: wire all repos, services, resolvers, handler; register PIN-protected route groups
- apps/api/gqlgen.yml: auto-discovers new schema files via glob (no change needed)
- apps/api/sqlc.yaml: auto-discovers new queries via glob (no change needed)
- apps/api/internal/appconfig/config.go: reuses WAVE-01 MediaConfig (BasePath, MaxUploadSize) — no new config needed

## Design Contracts
- Hard delete: all WAVE-04 entities use hard delete (no isActive flag). Cascade delete BodyCheckIn → BodyMeasurement + ProgressPhoto + physical photo files (DDEC-W04-001).
- Soft photo count guidance: 2-4 photos per check-in recommended but not enforced. Hard limit of 10 photos per check-in to prevent storage abuse (DDEC-W04-002).
- Measurement side rules: side (left/right) allowed only for paired measurement types (forearm, biceps, thigh, calf). Must be null for unpaired types (neck, shoulders, chest, waist, abdomen, hips). Single value for paired type without side = common (both sides same) (DDEC-W04-003).
- BodyWeightEntry per date: multiple entries per date allowed (different sources: scale vs manual). No unique constraint on (user_id, date). Latest weight determined by created_at DESC (DDEC-W04-004).
- DailyLog auto-creation: createCardioEntry(date) mutation auto-creates a DailyLog record if none exists for that user+date. If daily_log table does not exist (WAVE-03 not deployed), WAVE-04 must create daily_log migration as prerequisite (DDEC-W04-005).
- PIN auth: WAVE-01 middleware guards all WAVE-04 GraphQL and REST endpoints
- MIME detection: http.DetectContentType() server-side (512 bytes) as primary validation for progress photos
- File storage: <WAVE-01 BasePath>/progress-photos/<checkin_id>/<uuid>.<ext>
- BodyWeightEntry source enum: SCALE, MANUAL, UNKNOWN
- Photo angle enum: FRONT, SIDE, BACK, CUSTOM (custom allows elaboration via label)
- Error codes (REST): FILE_TOO_LARGE, INVALID_FILE_TYPE, NOT_FOUND, INTERNAL_ERROR, UNAUTHORIZED
- Error format: { "error": { "code": "ERROR_CODE", "message": "Human readable" } } per TDEC-027

## Data API Integration And Operations

### Database Schema

**cardio_entry**
- id (UUID, PK), user_id (UUID, FK → users), daily_log_id (UUID, FK → daily_log ON DELETE CASCADE), cardio_type (VARCHAR, NOT NULL), duration_minutes (INT, NOT NULL), avg_pulse (INT, nullable), heart_rate_zone (VARCHAR, nullable), notes (TEXT, nullable), created_at, updated_at
- Indexes: idx_cardio_entry_daily_log (daily_log_id)
- CardioType enum: walking, running, bike, elliptical, treadmill, other

**body_weight_entry**
- id (UUID, PK), user_id (UUID, FK → users), date (DATE, NOT NULL), weight (REAL, NOT NULL), source (VARCHAR, NOT NULL), notes (TEXT, nullable), created_at, updated_at
- Indexes: idx_body_weight_user_date (user_id, date DESC)
- WeightSource enum: scale, manual, unknown

**body_check_in**
- id (UUID, PK), user_id (UUID, FK → users), date (DATE, NOT NULL UNIQUE), weight (REAL, nullable), body_fat_percentage (REAL, nullable), notes (TEXT, nullable), created_at, updated_at
- Indexes: idx_body_check_in_date (date DESC)
- One check-in per date

**body_measurement**
- id (UUID, PK), check_in_id (UUID, FK → body_check_in ON DELETE CASCADE), measurement_type (VARCHAR, NOT NULL), side (VARCHAR, nullable), value (REAL, NOT NULL), created_at, updated_at
- Unique: (check_in_id, measurement_type, side)
- Indexes: idx_body_measurement_checkin (check_in_id)
- MeasurementType enum: neck, shoulders, forearms, biceps, chest, waist, abdomen, hips, thigh, calf
- MeasurementSide: left, right (nullable = common/unpaired)

**progress_photo**
- id (UUID, PK), check_in_id (UUID, FK → body_check_in ON DELETE CASCADE), file_path (VARCHAR, NOT NULL), original_file_name (VARCHAR, NOT NULL), mime_type (VARCHAR, NOT NULL), size_bytes (BIGINT, NOT NULL), angle (VARCHAR, nullable), label (VARCHAR, nullable), notes (TEXT, nullable), created_at, updated_at
- Indexes: idx_progress_photo_checkin (check_in_id)
- PhotoAngle enum: front, side, back, custom
- File size: 25MB max per upload
- Allowed MIME: image/jpeg, image/png, image/webp

**week_flag**
- id (UUID, PK), user_id (UUID, FK → users), week_start_date (DATE, NOT NULL), flag_type (VARCHAR, NOT NULL), notes (TEXT, nullable), created_at, updated_at
- Unique: (week_start_date, flag_type)
- Indexes: idx_week_flag_week (week_start_date)
- WeekFlagType enum: poor_sleep, high_stress, illness, injury_pain, cycle, calorie_deficit, surplus, maintenance, missed_workouts, travel

### GraphQL Operations
- CardioEntry: createCardioEntry, updateCardioEntry, deleteCardioEntry, cardioEntries (by date), cardioEntry (by ID)
- BodyWeightEntry: createBodyWeightEntry, updateBodyWeightEntry, deleteBodyWeightEntry, bodyWeightEntries (date range), bodyWeightEntry (by ID), latestBodyWeight
- BodyCheckIn: createBodyCheckIn, updateBodyCheckIn, deleteBodyCheckIn, bodyCheckIns (date range), bodyCheckIn (by ID with nested measurements + photos)
- BodyMeasurement: createBodyMeasurement, updateBodyMeasurement, deleteBodyMeasurement (nested under check-in)
- ProgressPhoto: progressPhotos (by checkInId, via GraphQL), deleteProgressPhoto (via GraphQL), REST upload/download
- WeekFlag: createWeekFlag, deleteWeekFlag, weekFlags (by weekStartDate)
- Union results: CardioEntryResult, BodyWeightEntryResult, BodyCheckInResult, BodyMeasurementResult, WeekFlagResult — each union is Success | ValidationError | AuthError

### REST Endpoints
- POST /api/v1/progress-photos/upload — multipart (checkInId, file, angle, label, notes). PIN-protected.
- GET /api/v1/progress-photos/{id} — download file. PIN-protected.
- DELETE /api/v1/progress-photos/{id} — delete. PIN-protected.
- File validation: server-side MIME detection, 25MB max, memory-safe upload per TDEC-008
- Physical file deletion on DELETE. Log error if file deletion fails, return 204.

### Log Markers
- [CardioEntry][create|update|delete|get|list], [BodyWeightEntry][create|update|delete|get|list|latest]
- [BodyCheckIn][create|update|delete|get|list], [BodyMeasurement][create|update|delete]
- [ProgressPhoto][upload|download|delete], [WeekFlag][create|delete|list]
- Sensitive data (weight, body fat %, measurement values, photo content, week flag notes) NOT logged

### Operations
- PostgreSQL: goose migrations 00082-00087, sequential, reversible
- Existing Docker Compose stack, no new services
- Media volume from WAVE-01 reused for progress-photos subdirectory

## Security Privacy And Compliance
- All endpoints protected by WAVE-01 PIN auth middleware (GraphQL + REST)
- When PIN disabled, endpoints accessible without auth (consistent with TDEC-037)
- Server-side MIME detection prevents Content-Type spoofing
- UUID-based parameters prevent path traversal (file path resolved from DB, not user input)
- Uploaded file names sanitized: UUID-based storage, no user-provided path segments
- Sensitive data NOT logged: body weight values, body fat %, measurement values, photo content, week flag notes
- Log markers record entity type, action, success/failure, entity ID only
- Progress photos: high sensitivity — stored in media volume, served only through authenticated endpoint
- Photos excluded from AI export by default per RULE-025 (opt-in only)
- Hard limit of 10 photos per check-in to prevent storage abuse
- No sensitive content (personal health data, file content) appears in application logs
- Cardio entry and body weight operations are scoped to default user per MVP constraint

## Implementation Slices

| Slice ID | Name | Description |
| --- | --- | --- |
| SLICE-W04-001 | DB migrations | Create goose migrations 00082-00087 for cardio_entry, body_weight_entry, body_check_in, body_measurement, progress_photo, week_flag tables with indexes, FKs, cascades, and constraints |
| SLICE-W04-002 | sqlc queries | Define CRUD queries for all 6 entities: cardio entries (by dailyLogId), body weight (date range, latest), check-ins (date range), measurements (by checkInId), photos (by checkInId), week flags (by weekStartDate) |
| SLICE-W04-003 | Repository adapters | Implement 6 repo adapters with sqlc-generated code and error mapping (not found, constraint violations) |
| SLICE-W04-004 | Services layer | Implement 4 services: cardio (DailyLog auto-creation, type/duration/zone validation), body weight (weight > 0, source enum), body check-in (weight/fat % validation, cascade rules), week flag (flagType enum, unique per week) |
| SLICE-W04-005 | GraphQL schema | Add 4 schema files: cardio.graphql, body_weight.graphql, body_checkin.graphql, week_flag.graphql with types, enums, inputs, queries, mutations, union results |
| SLICE-W04-006 | GraphQL resolvers | Implement cardio, body weight, body check-in, week flag resolvers with PIN auth guard and union error returns |
| SLICE-W04-007 | ProgressPhoto REST handler | Upload (multipart with MIME/size validation), download (file stream), delete (record + physical file) endpoints following WAVE-02 exercise_media.go pattern |
| SLICE-W04-008 | Main wiring | Wire all repos, services, resolvers, handler; register PIN-protected route groups in cmd/server/main.go |

## Acceptance Criteria

| AC ID | Description |
| --- | --- |
| AC-W04-001 | Cardio entry can be created via GraphQL mutation with cardioType (enum), durationMinutes, optional avgPulse, optional heartRateZone (enum 1-5/unknown), optional notes. |
| AC-W04-002 | Cardio entry type is validated against allowed enum values (walking/running/bike/elliptical/treadmill/other). Invalid type returns ValidationError. |
| AC-W04-003 | Cardio durationMinutes is required, positive integer. 0 or negative returns ValidationError. |
| AC-W04-004 | Cardio entry avgPulse, if provided, is positive integer. Non-positive returns ValidationError. |
| AC-W04-005 | Cardio entry heartRateZone, if provided, is 1-5 or "unknown". Invalid zone returns ValidationError. |
| AC-W04-006 | Cardio entry is linked to a DailyLog: providing date auto-creates or attaches to existing DailyLog for that user+date. |
| AC-W04-007 | Cardio entry can be read by ID. Returns full cardio entry. |
| AC-W04-008 | Cardio entries can be listed by date (DailyLog date). Returns all entries for that user+date. |
| AC-W04-009 | Cardio entry can be updated (type, duration, pulse, zone, notes). Updated entry returned. |
| AC-W04-010 | Cardio entry can be deleted (hard delete). Returns success indicator. |
| AC-W04-011 | BodyWeightEntry can be created with date (required), weight (required, > 0), source enum (required scale/manual/unknown), optional notes. |
| AC-W04-012 | BodyWeightEntry weight must be > 0. 0 or negative returns ValidationError. |
| AC-W04-013 | BodyWeightEntry source is validated against allowed enum values. Invalid source returns ValidationError. |
| AC-W04-014 | BodyWeightEntry can be read by ID. Returns full entry. |
| AC-W04-015 | BodyWeightEntry can be listed with date range filtering (startDate, endDate). Results ordered by date DESC. |
| AC-W04-016 | Latest body weight can be queried (latestBodyWeight). Returns most recent entry by created_at. Returns null when no entries exist. |
| AC-W04-017 | BodyWeightEntry can be updated (weight, source, notes). Updated entry returned. |
| AC-W04-018 | BodyWeightEntry can be deleted (hard delete). Returns success indicator. |
| AC-W04-019 | BodyCheckIn can be created with date (required), optional weight, optional bodyFatPercentage, optional notes. |
| AC-W04-020 | BodyCheckIn weight, if provided, must be > 0. Non-positive returns ValidationError. |
| AC-W04-021 | BodyCheckIn bodyFatPercentage, if provided, must be > 0 and <= 100. Out of range returns ValidationError. |
| AC-W04-022 | BodyCheckIn can be read by ID with nested measurements and photos. Returns check-in with all children. |
| AC-W04-023 | BodyCheckIn can be listed with date range filtering (startDate, endDate). Results ordered by date DESC. |
| AC-W04-024 | BodyCheckIn can be updated (weight, bodyFatPercentage, notes). Updated entry returned. |
| AC-W04-025 | BodyCheckIn can be deleted (hard delete cascades to measurements, photos, and physical photo files). Returns success indicator. |
| AC-W04-026 | BodyMeasurement can be created within a check-in with measurementType (enum), value (> 0), optional side (left/right for paired types). |
| AC-W04-027 | BodyMeasurement measurementType validated against 10 allowed types. Invalid type returns ValidationError. |
| AC-W04-028 | BodyMeasurement value must be > 0. 0 or negative returns ValidationError. |
| AC-W04-029 | BodyMeasurement side is validated: allowed only for paired types (forearm, biceps, thigh, calf). Must be null for unpaired types (neck, shoulders, chest, waist, abdomen, hips). Invalid side returns ValidationError. |
| AC-W04-030 | BodyMeasurement can be updated (value, side). Updated entry returned. |
| AC-W04-031 | BodyMeasurement can be deleted. Returns success indicator. |
| AC-W04-032 | ProgressPhoto can be uploaded (multipart REST) associated with a check-in, with angle (enum), optional label, optional notes. |
| AC-W04-033 | ProgressPhoto angle is validated: front/side/back/custom. Invalid angle returns validation error. |
| AC-W04-034 | ProgressPhoto file is stored in media storage with path <BasePath>/progress-photos/<checkin_id>/<uuid>.<ext>. |
| AC-W04-035 | ProgressPhoto file is validated server-side: allowed MIME types (JPEG/PNG/WEBP), 25MB max size. Rejected types/sizes return validation error. |
| AC-W04-036 | ProgressPhoto can be downloaded via GET endpoint (PIN-protected). Returns file with correct content type. |
| AC-W04-037 | ProgressPhoto can be deleted: DB record removed and physical file deleted from disk. 204 No Content returned. |
| AC-W04-038 | ProgressPhotos can be listed by check-in ID via GraphQL query. |
| AC-W04-039 | WeekFlag can be created with weekStartDate, flagType (enum), optional notes. |
| AC-W04-040 | WeekFlag flagType validated against allowed enum values. Invalid type returns ValidationError. |
| AC-W04-041 | WeekFlag can be listed by week start date. Returns all flags for that week. |
| AC-W04-042 | WeekFlag can be deleted. Returns success indicator. |
| AC-W04-043 | All WAVE-04 GraphQL mutations return AuthError when PIN session header is missing or invalid. |
| AC-W04-044 | All WAVE-04 REST endpoints return 401 when PIN session header is missing or invalid. |

## Exit Criteria

| EC ID | Description |
| --- | --- |
| EC-W04-001 | AC-W04-001 through AC-W04-044 pass via TEST-W04-001 through TEST-W04-030 |
| EC-W04-002 | gqlgen codegen produces valid Go code for WAVE-04 schema without drift |
| EC-W04-003 | sqlc codegen produces valid Go code for WAVE-04 queries without drift |
| EC-W04-004 | WAVE-01 PIN auth guard protects all WAVE-04 GraphQL and REST endpoints. Existing admin auth unchanged. |
| EC-W04-005 | WAVE-01 admin auth and health test suite still passes after WAVE-04 changes |
| EC-W04-006 | All 6 migrations (00082-00087) apply and roll back in sequence without errors |
| EC-W04-007 | File size (25MB) and type (image/jpeg, image/png, image/webp) validation enforced for progress photo uploads |
| EC-W04-008 | Body measurement value > 0 validation enforced for all 10 measurement types |
| EC-W04-009 | Measurement side validation: side allowed only for paired types (forearm, biceps, thigh, calf), rejected for unpaired |
| EC-W04-010 | Cardio + Body Weight + Check-In round-trip integration test passes |
| EC-W04-011 | Lint passes for all changed packages |
| EC-W04-012 | No sensitive content (weight, body fat %, measurement values, photo content, week flag notes) in application logs |
| EC-W04-013 | DailyLog is auto-created for CardioEntry when no daily_log exists for user+date |
| EC-W04-014 | Cascade delete: deleting BodyCheckIn deletes its measurements, photos (DB + physical files from disk) |

## Verification Obligations

| Test ID | Description | Type | Command |
| --- | --- | --- | --- |
| TEST-W04-001 | CardioEntry repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)cardio_entry_repo' |
| TEST-W04-002 | CardioEntry service validation (type enum, duration > 0, pulse positive, zone valid) | unit | bunx nx run api:test -- --run '(?i)cardio_entry_service_validation' |
| TEST-W04-003 | CardioEntry auto-create DailyLog when DailyLog not provided | integration | bunx nx run api:test -- --run '(?i)cardio_entry_dailylog_auto' |
| TEST-W04-004 | CardioEntry GraphQL resolver integration tests | integration | bunx nx run api:test -- --run '(?i)cardio_entry_resolver' |
| TEST-W04-005 | BodyWeightEntry repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)body_weight_repo' |
| TEST-W04-006 | BodyWeightEntry service validation (weight > 0, source enum) | unit | bunx nx run api:test -- --run '(?i)body_weight_service' |
| TEST-W04-007 | BodyWeightEntry GraphQL resolver integration tests | integration | bunx nx run api:test -- --run '(?i)body_weight_resolver' |
| TEST-W04-008 | BodyWeightEntry latest weight query (returns null when no entries) | integration | bunx nx run api:test -- --run '(?i)body_weight_latest' |
| TEST-W04-009 | BodyCheckIn repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)body_checkin_repo' |
| TEST-W04-010 | BodyCheckIn service validation (date required, weight > 0, bodyFat% 0-100) | unit | bunx nx run api:test -- --run '(?i)body_checkin_service' |
| TEST-W04-011 | BodyCheckIn GraphQL resolver integration tests (nested measurements, photos) | integration | bunx nx run api:test -- --run '(?i)body_checkin_resolver' |
| TEST-W04-012 | BodyCheckIn cascade delete (cascades to measurements, photos, physical files) | integration | bunx nx run api:test -- --run '(?i)body_checkin_cascade' |
| TEST-W04-013 | BodyMeasurement repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)body_measurement_repo' |
| TEST-W04-014 | BodyMeasurement validation (type enum, value > 0, side rules for paired vs unpaired) | unit | bunx nx run api:test -- --run '(?i)body_measurement_service' |
| TEST-W04-015 | ProgressPhoto repository unit tests | unit | bunx nx run api:test -- --run '(?i)progress_photo_repo' |
| TEST-W04-016 | ProgressPhoto REST handler integration tests (upload, download, delete) | integration | bunx nx run api:test -- --run '(?i)progress_photo_handler' |
| TEST-W04-017 | ProgressPhoto file type validation (rejects non-image MIME types) | unit | bunx nx run api:test -- --run '(?i)progress_photo_filetype' |
| TEST-W04-018 | ProgressPhoto size validation (rejects > 25MB) | unit | bunx nx run api:test -- --run '(?i)progress_photo_filesize' |
| TEST-W04-019 | ProgressPhoto physical file deletion on DELETE | integration | bunx nx run api:test -- --run '(?i)progress_photo_file_delete' |
| TEST-W04-020 | WeekFlag repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)week_flag_repo' |
| TEST-W04-021 | WeekFlag validation (flagType enum, duplicate per week rejected) | unit | bunx nx run api:test -- --run '(?i)week_flag_service' |
| TEST-W04-022 | WeekFlag unique constraint (one flag per type per week) | integration | bunx nx run api:test -- --run '(?i)week_flag_unique' |
| TEST-W04-023 | All WAVE-04 GraphQL operations return AuthError without PIN session | integration | bunx nx run api:test -- --run '(?i)wave04_auth' |
| TEST-W04-024 | All WAVE-04 REST endpoints return 401 without PIN session | integration | bunx nx run api:test -- --run '(?i)wave04_media_auth' |
| TEST-W04-025 | Migration smoke test (00082-00087 up + down) | integration | bunx nx run api:test -- --run '(?i)migration_wave04' |
| TEST-W04-026 | Codegen drift check (gqlgen + sqlc) | codegen | bunx nx run api:codegen && bunx nx run graphql:codegen |
| TEST-W04-027 | Log privacy: no weight, body fat %, measurement values, photo content in logs | unit | bunx nx run api:test -- --run '(?i)wave04_log_sanitize' |
| TEST-W04-028 | Go lint for API package | lint | bunx nx run api:lint |
| TEST-W04-029 | GraphQL schema validate | codegen | bunx nx run graphql:validate |
| TEST-W04-030 | Cardio + Body + Check-in round-trip integration test (full lifecycle) | integration | bunx nx run api:test -- --run '(?i)wave04_roundtrip' |

## Rollout Rollback And Compatibility
- Rollout: merge PR, CI builds and runs tests, deploy via Dokploy compose update. New tables created via goose migrations.
- Rollback: revert PR, CI builds previous image, Dokploy compose update rolls back. Run goose down migrations for 00082-00087.
- Compatibility: all new operations are additive. No existing API changes. WAVE-01 endpoints (health, admin GraphQL, media REST scaffold, users REST) and WAVE-02 endpoints unchanged.
- Migration: goose migrations 00082-00087 run at startup. Down migrations available for rollback.
- DailyLog dependency: if WAVE-03 not yet deployed, WAVE-04 migration must ensure daily_log table exists (either wait for WAVE-03 or create stub). See DDEC-W04-005.

## Handoff Packets
- HANDOFF-W04-001: This wave brief document
- HANDOFF-W04-002: Planner reports (6 scopes)
- HANDOFF-W04-003: Reviewer evidence (7 perspectives)
- HANDOFF-W04-004: Final fit review evidence

## Design Decisions

| DDEC ID | Decision | Rationale |
| --- | --- | --- |
| DDEC-W04-001 | Hard delete for all WAVE-04 entities | No referential integrity concerns unlike WAVE-02 exercises. Cascade delete BodyCheckIn cleans up measurements, photos, and physical files. |
| DDEC-W04-002 | Photo count: 2-4 recommended, hard limit 10 | EDGE-006 ambiguity resolved: soft guidance for MVP. Hard limit prevents storage abuse. |
| DDEC-W04-003 | Measurement side: allowed only for paired types | PRD §13.4: paired measurements (forearm, biceps, thigh, calf) may have left/right/null. Unpaired types must have null side. |
| DDEC-W04-004 | BodyWeightEntry allows multiple entries per date | Different sources (scale, manual) can produce different values on same date. Latest entry by created_at used for dashboard. |
| DDEC-W04-005 | DailyLog auto-creation for CardioEntry | Domain model invariant: cardio must belong to DailyLog. If WAVE-03 not deployed, WAVE-04 must create daily_log table migration as prerequisite. |

## Reviewer Verdicts

| Wave | Perspective | Attempt | Verdict | Reviewer Report | Required Revisions | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| WAVE-04 | product-scope-and-ac | 1 | approved | review-product-scope-and-ac-attempt-1.md | none | 44 ACs cover all scope, edge cases documented |
| WAVE-04 | architecture-codebase-fit | 1 | approved | review-architecture-codebase-fit-attempt-1.md | none | Codebase touchpoints well-documented, 8 slices |
| WAVE-04 | data-api-integration-ops | 1 | approved | review-data-api-integration-ops-attempt-1.md | none | Data/API/ops coverage adequate, design decisions documented |
| WAVE-04 | security-privacy-compliance | 1 | approved | review-security-privacy-compliance-attempt-1.md | none | Server-side MIME, PIN auth, log privacy covered |
| WAVE-04 | testing-exit-criteria | 1 | approved | review-testing-exit-criteria-attempt-1.md | none | 30 test obligations cover all AC and EC |
| WAVE-04 | sequencing-other-wave-fit | 1 | approved | review-sequencing-other-wave-fit-attempt-1.md | none | Dependency order correct, WAVE-03 DailyLog noted |
| WAVE-04 | traceability-consistency | 1 | approved | review-traceability-consistency-attempt-1.md | none | Source traceability documented, stable IDs used |

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W04-001 | WAVE-04 | operations | needs-owner-decision | WAVE-03 | Should WAVE-04 include daily_log table creation or require WAVE-03 to be deployed first? | Affects deployment ordering and migration strategy | WAVE-04 should create daily_log table migration as prerequisite | planner-sequencing-fit-attempt-1.md | open | DailyLog auto-creation logic required. WAVE-04 should include a daily_log migration if WAVE-03 not yet deployed. |
| DQ-W04-002 | WAVE-04 | operations | deferred | EDGE-006 | Is 2-4 photos per check-in a hard requirement or recommendation? | Affects validation logic | Soft recommendation (warn, don't block) with hard limit of 10 | planner-product-ac-attempt-1.md | resolved | Soft guidance. DDEC-W04-002. |
| DQ-W04-003 | WAVE-04 | data-ops | resolved | EDGE-007 | Should body measurement value 0 or negative be rejected? | Data integrity | Reject 0 and negative values. Validated in service layer. | planner-product-ac-attempt-1.md | resolved | AC-W04-028: value must be > 0. |
| DQ-W04-004 | WAVE-04 | security | deferred | TDEC-008 | Should progress photo URLs be time-limited (signed URLs)? | Session-gated access sufficient for MVP self-hosted | Signed URLs add complexity for single-user MVP | planner-security-compliance-attempt-1.md | deferred | Deferred post-MVP. |
| DQ-W04-005 | WAVE-04 | data-ops | resolved | WAVE-01 | What exact file storage path pattern does WAVE-01 MediaConfig provide for progress photos? | Drives migration and handler design | Use WAVE-01 BasePath/progress-photos/<checkin_id>/<uuid>.<ext> | planner-data-integration-ops-attempt-1.md | deferred | Confirmed after WAVE-01 implementation. WAVE-04 assumes composable BasePath. |

## Traceability
- docs/prd-waves/waves/wave-04.md: source wave boundary, outcomes, capability groups
- docs/product-verified/functional-spec.md: Cardio §12 REQ-007, Body Tracking §13 REQ-008/REQ-009
- docs/product-verified/domain-model.md: CardioEntry, BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto, WeekFlag entities
- docs/product-verified/acceptance-criteria.md: AC-012–AC-016, AC-048–AC-057
- docs/product-verified/edge-cases.md: EDGE-006, EDGE-007
- docs/product-verified/business-rules.md: RULE-005
- docs/product-verified/user-flows.md: §26.4 Add Cardio, §26.5 Weekly Check-In, §26.6 Add Body Weight
- docs/product-verified/actors-and-permissions.md: user permissions for cardio, body tracking
- docs/development-plan.xml: M-API, M-PRD-WAVE-DETAILER module contracts
- docs/knowledge-graph.xml: existing module boundaries
- docs/prd-wave-details/waves/wave-01.md: WAVE-01 dependency contracts (PIN auth, media scaffold)
- docs/prd-wave-details/waves/wave-02.md: WAVE-02 patterns (repository, service, resolver, handler structure)
- apps/api/internal: existing codebase patterns for service/repository/middleware/handler structure
- docs/prd-waves/frontend-pages/page-004.md: cardio backend dependencies
- docs/prd-waves/frontend-pages/page-005.md: body measurements backend dependencies
- docs/prd-waves/frontend-pages/page-006.md: progress photos backend dependencies
- docs/prd-waves/frontend-pages/page-001.md: dashboard backend dependency (latestBodyWeight)