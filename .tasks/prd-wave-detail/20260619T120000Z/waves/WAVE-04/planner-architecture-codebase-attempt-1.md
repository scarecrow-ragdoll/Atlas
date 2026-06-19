# WAVE-04 Planner: Architecture / Codebase

## Codebase Patterns (from WAVE-01 and WAVE-02)

Following established patterns:
- Repository → Service → Resolver/Handler architecture (Go)
- sqlc for DB queries, gqlgen for GraphQL schema
- PIN auth middleware from WAVE-01 protects all endpoints
- Media REST handler pattern from WAVE-01/WAVE-02 for progress photos
- Migration numbering: WAVE-02 uses 00080-00081. WAVE-04 migrations start at 00082.
- Extend type Query/Mutation pattern in GraphQL schema

## Proposed Modules and Files

### DB Migrations
- `apps/api/internal/repository/postgres/migrations/00082_cardio_entries.sql` — cardio_entry table
- `apps/api/internal/repository/postgres/migrations/00083_body_weight_entries.sql` — body_weight_entry table
- `apps/api/internal/repository/postgres/migrations/00084_body_check_ins.sql` — body_check_in table
- `apps/api/internal/repository/postgres/migrations/00085_body_measurements.sql` — body_measurement table
- `apps/api/internal/repository/postgres/migrations/00086_progress_photos.sql` — progress_photo table
- `apps/api/internal/repository/postgres/migrations/00087_week_flags.sql` — week_flag table
- `apps/api/internal/repository/postgres/migrations/00088_cardio_daily_log_fk.sql` — add FK for cardio→daily_log (if daily_log migration exists from WAVE-01 or later)

### sqlc Queries
- `apps/api/internal/repository/postgres/queries/cardio_entries.sql` — CRUD, list by dailyLogId
- `apps/api/internal/repository/postgres/queries/body_weight_entries.sql` — CRUD, list with date range, latest
- `apps/api/internal/repository/postgres/queries/body_check_ins.sql` — CRUD, list with date range, nested measurements/photos
- `apps/api/internal/repository/postgres/queries/body_measurements.sql` — CRUD by checkInId
- `apps/api/internal/repository/postgres/queries/progress_photos.sql` — CRUD, list by checkInId
- `apps/api/internal/repository/postgres/queries/week_flags.sql` — CRUD, list by weekStartDate

### Repository Adapters
- `apps/api/internal/repository/postgres/cardio_entry_repo.go`
- `apps/api/internal/repository/postgres/body_weight_entry_repo.go`
- `apps/api/internal/repository/postgres/body_check_in_repo.go`
- `apps/api/internal/repository/postgres/body_measurement_repo.go`
- `apps/api/internal/repository/postgres/progress_photo_repo.go`
- `apps/api/internal/repository/postgres/week_flag_repo.go`

### Services (transport-neutral)
- `apps/api/internal/service/cardio.go` — cardio validation, dailyLog auto-creation
- `apps/api/internal/service/body_weight.go` — body weight validation
- `apps/api/internal/service/body_checkin.go` — check-in business logic, measurement/photo cascade
- `apps/api/internal/service/week_flag.go` — week flag validation

### GraphQL Schema
- `libs/graphql/schema/cardio.graphql` — CardioEntry type, queries, mutations
- `libs/graphql/schema/body_weight.graphql` — BodyWeightEntry type, queries, mutations
- `libs/graphql/schema/body_checkin.graphql` — BodyCheckIn, BodyMeasurement, ProgressPhoto types, queries, mutations
- `libs/graphql/schema/week_flag.graphql` — WeekFlag type, queries, mutations

### GraphQL Resolvers
- `apps/api/internal/graph/cardio.resolvers.go` — cardio CRUD resolvers
- `apps/api/internal/graph/body_weight.resolvers.go` — body weight CRUD resolvers
- `apps/api/internal/graph/body_checkin.resolvers.go` — check-in, measurement, photo resolvers
- `apps/api/internal/graph/week_flag.resolvers.go` — week flag resolvers

### REST Handlers
- `apps/api/internal/handler/progress_photo_handler.go` — multipart upload, download, delete (following exercise_media.go pattern)

### Wiring
- `apps/api/cmd/server/main.go` — wire all new repos, services, resolvers, handlers, PIN-protected route groups

## Implementation Slices

| Slice ID | Name | Description |
|---|---|---|
| SLICE-W04-001 | DB migrations (6 tables) | Create goose migrations 00082-00087 for cardio_entry, body_weight_entry, body_check_in, body_measurement, progress_photo, week_flag tables with indexes, FKs, and cascades. |
| SLICE-W04-002 | sqlc queries | Define all CRUD queries for cardio, weight, check-ins, measurements, photos, week flags. |
| SLICE-W04-003 | Repository adapters | Implement 6 repository adapters with sqlc-generated code and error mapping. |
| SLICE-W04-004 | Services layer | Implement 4 transport-neutral services with validation. |
| SLICE-W04-005 | GraphQL schema | Add cardio.graphql, body_weight.graphql, body_checkin.graphql, week_flag.graphql with types, queries, mutations, union results. |
| SLICE-W04-006 | GraphQL resolvers | Implement cardio, body weight, body check-in, week flag resolvers. |
| SLICE-W04-007 | ProgressPhoto REST handler | Upload (multipart with validation), download, and delete endpoints following WAVE-02 exercise media pattern. |
| SLICE-W04-008 | Main wiring | Wire all repos, services, resolvers, handlers; register PIN-protected route groups. |

## Dependencies on WAVE-01
- PIN auth middleware (protects all WAVE-04 endpoints)
- Media storage config (BasePath, MaxUploadSize)
- Common GraphQL foundation and error types
- DailyLog table from WAVE-03 or auto-creation logic

## Dependencies on Other Waves
- WAVE-01: PIN auth, media scaffold, daily_log table reference for cardio
- WAVE-03: DailyLog table for cardio entry FK. If WAVE-03 not yet deployed, cardio auto-creation must handle DailyLog absence.
- WAVE-05: Can partially parallelize (no shared tables)
- WAVE-06: WAVE-04 provides body weight data for charts
- WAVE-07/08: WAVE-04 provides cardio, weight, check-in, photos for AI export

## Codebase Touchpoints
- apps/api/internal/repository/postgres/migrations/ — 6 new migration files
- apps/api/internal/repository/postgres/queries/ — 6 new query files
- apps/api/internal/repository/postgres/ — 6 new repo adapters
- apps/api/internal/service/ — 4 new service files
- apps/api/internal/graph/ — 4 new resolver files
- apps/api/internal/handler/ — 1 new REST handler
- libs/graphql/schema/ — 4 new schema files
- apps/api/cmd/server/main.go — wiring additions
- apps/api/internal/appconfig/config.go — no new config needed (reuses WAVE-01 MediaConfig)
- apps/api/gqlgen.yml — auto-discovers new schema via glob
- apps/api/sqlc.yaml — auto-discovers new queries via glob