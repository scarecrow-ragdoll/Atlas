# WAVE-04 Planner: Testing / Exit Criteria

## Testing Strategy

### Test Categories

1. **Unit tests** — repository, service layer, validation logic
2. **Integration tests** — resolver round-trips, middleware chain with PIN auth, REST handlers
3. **Codegen checks** — gqlgen, sqlc drift detection
4. **Lint checks** — Go lint for API package

### Test Fixture Strategy (following WAVE-02 pattern)
- Use WAVE-01 test helper infrastructure (test DB, PIN session setup)
- Each test creates its own data via repo/service calls
- Photos: use `bytes.Buffer` multipart uploads in REST handler tests

## Proposed Test Obligations

| Test ID | Description | Type | Command |
|---|---|---|---|
| TEST-W04-001 | CardioEntry repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)cardio_entry_repo' |
| TEST-W04-002 | CardioEntry service validation (type enum, duration > 0, pulse positive, zone valid) | unit | bunx nx run api:test -- --run '(?i)cardio_entry_service_validation' |
| TEST-W04-003 | CardioEntry auto-create DailyLog when dailyLogId not provided | integration | bunx nx run api:test -- --run '(?i)cardio_entry_dailylog_auto' |
| TEST-W04-004 | CardioEntry GraphQL resolver integration tests | integration | bunx nx run api:test -- --run '(?i)cardio_entry_resolver' |
| TEST-W04-005 | BodyWeightEntry repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)body_weight_repo' |
| TEST-W04-006 | BodyWeightEntry service validation (weight > 0, source enum) | unit | bunx nx run api:test -- --run '(?i)body_weight_service' |
| TEST-W04-007 | BodyWeightEntry GraphQL resolver integration tests | integration | bunx nx run api:test -- --run '(?i)body_weight_resolver' |
| TEST-W04-008 | BodyWeightEntry latest weight query | integration | bunx nx run api:test -- --run '(?i)body_weight_latest' |
| TEST-W04-009 | BodyCheckIn repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)body_checkin_repo' |
| TEST-W04-010 | BodyCheckIn service validation (date required, weight > 0, bodyFat% 0-100) | unit | bunx nx run api:test -- --run '(?i)body_checkin_service' |
| TEST-W04-011 | BodyCheckIn GraphQL resolver integration tests (nested measurements, photos) | integration | bunx nx run api:test -- --run '(?i)body_checkin_resolver' |
| TEST-W04-012 | BodyCheckIn cascade delete (cascades to measurements and photos) | integration | bunx nx run api:test -- --run '(?i)body_checkin_cascade' |
| TEST-W04-013 | BodyMeasurement repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)body_measurement_repo' |
| TEST-W04-014 | BodyMeasurement validation (type enum, value > 0, side rules for paired vs unpaired) | unit | bunx nx run api:test -- --run '(?i)body_measurement_service' |
| TEST-W04-015 | ProgressPhoto repository unit tests | unit | bunx nx run api:test -- --run '(?i)progress_photo_repo' |
| TEST-W04-016 | ProgressPhoto REST handler integration tests (upload, download, delete) | integration | bunx nx run api:test -- --run '(?i)progress_photo_handler' |
| TEST-W04-017 | ProgressPhoto file type validation (rejects non-image) | unit | bunx nx run api:test -- --run '(?i)progress_photo_filetype' |
| TEST-W04-018 | ProgressPhoto size validation (25MB limit) | unit | bunx nx run api:test -- --run '(?i)progress_photo_filesize' |
| TEST-W04-019 | ProgressPhoto physical file deletion on delete | integration | bunx nx run api:test -- --run '(?i)progress_photo_file_delete' |
| TEST-W04-020 | WeekFlag repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)week_flag_repo' |
| TEST-W04-021 | WeekFlag validation (flagType enum) | unit | bunx nx run api:test -- --run '(?i)week_flag_service' |
| TEST-W04-022 | WeekFlag unique constraint (one flag per type per week) | integration | bunx nx run api:test -- --run '(?i)week_flag_unique' |
| TEST-W04-023 | All WAVE-04 GraphQL operations return AuthError without PIN session | integration | bunx nx run api:test -- --run '(?i)wave04_auth' |
| TEST-W04-024 | All WAVE-04 REST endpoints return 401 without PIN session | integration | bunx nx run api:test -- --run '(?i)wave04_media_auth' |
| TEST-W04-025 | Migration smoke test (00082-00087 up + down) | integration | bunx nx run api:test -- --run '(?i)migration_wave04' |
| TEST-W04-026 | Codegen drift check (gqlgen + sqlc) | codegen | bunx nx run api:codegen && bunx nx run graphql:codegen |
| TEST-W04-027 | Log privacy: no weight, body fat, measurement values in logs | unit | bunx nx run api:test -- --run '(?i)wave04_log_sanitize' |
| TEST-W04-028 | Go lint for API package | lint | bunx nx run api:lint |
| TEST-W04-029 | GraphQL schema validate | codegen | bunx nx run graphql:validate |
| TEST-W04-030 | Cardio + Body + Check-in round-trip integration test | integration | bunx nx run api:test -- --run '(?i)wave04_roundtrip' |

## Exit Criteria

| EC ID | Description |
|---|---|
| EC-W04-001 | All AC-W04-001 through AC-W04-044 pass via TEST-W04-001 through TEST-W04-030 |
| EC-W04-002 | gqlgen codegen produces valid Go code for WAVE-04 schema without drift |
| EC-W04-003 | sqlc codegen produces valid Go code for WAVE-04 queries without drift |
| EC-W04-004 | WAVE-01 PIN auth guard protects all WAVE-04 GraphQL and REST endpoints. Existing admin auth unchanged. |
| EC-W04-005 | WAVE-01 admin auth and health test suite still passes after WAVE-04 changes |
| EC-W04-006 | All 6 migrations (00082-00087) apply and roll back in sequence without errors |
| EC-W04-007 | File size and type validation enforced for progress photo uploads |
| EC-W04-008 | Body measurement value > 0 validation enforced for all measurement types |
| EC-W04-009 | Measurement side validation: side allowed only for paired measurement types |
| EC-W04-010 | Cardio + Body + Check-in round-trip integration test passes |
| EC-W04-011 | Lint passes for all changed packages |
| EC-W04-012 | No sensitive content (weight, body fat, measurement values, photo content) appears in application logs |
| EC-W04-013 | DailyLog is auto-created for CardioEntry when no daily_log exists for date |
| EC-W04-014 | Cascade delete: deleting a BodyCheckIn deletes its measurements, photos (DB + physical files) |

## Coverage
- Focused tests during active development. No 100% coverage requirement in this wave.
- Evidence recorded in .tasks/ on wave completion.