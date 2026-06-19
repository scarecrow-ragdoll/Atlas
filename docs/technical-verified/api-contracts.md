# API Contracts

## Surfaces

**Decision (TDEC-001):** Hybrid API model.
- GraphQL is the primary application API for CRUD and queries.
- REST is used for binary uploads/downloads and long-running operations.

Required endpoints:
- `POST /graphql`
- `POST /api/v1/media/upload`
- `GET /api/v1/media/{id}`
- `POST /api/v1/exports/ai`
- `GET /api/v1/exports/ai/{id}/download`
- `POST /api/v1/backups/export`
- `GET /api/v1/backups/{id}/download`
- `POST /api/v1/backups/import/dry-run`
- `POST /api/v1/backups/import/commit`
- `GET /api/v1/jobs/{id}`

GraphQL scope: settings, PIN/session, exercises, exercise metadata, daily logs, workout exercises, workout sets, cardio entries, body weight entries, body check-ins, body measurements, nutrition products, nutrition templates, nutrition overrides, week flags, charts data, AI prompt settings, AI review history.

REST scope: media upload/download, AI export ZIP generate/download, backup export/download, import dry-run/commit, job progress polling.

All API calls scoped to the default user.

## Requests And Responses

No endpoint catalog, request schemas, or response formats defined (TQ-API-002). Missing critical details:
- Exercise CRUD endpoints/operations
- DailyLog create/read/update operations
- WorkoutExercise and WorkoutSet operations
- CardioEntry CRUD
- Body tracking operations
- Nutrition operations
- Chart data queries (TQ-API-010)
- AI export trigger and download (TQ-API-011)
- Backup export/import (TQ-API-011)
- File upload endpoints for media

## Error And Validation Contracts

No error format, validation mapping, or error codes defined (TQ-API-003).

## Compatibility And Idempotency

No API versioning strategy (TQ-API-009). No idempotency guarantees for mutations (TQ-API-006).

## API Questions

| ID | Question | Severity | Status |
| --- | --- | --- | --- |
| TQ-API-001 | Primary API protocol undecided (GraphQL vs REST) | dev-blocking | **resolved** — hybrid model (TDEC-001) |
| TQ-API-002 | No endpoint catalog or operation definitions | dev-blocking | **resolved** (TDEC-026) |
| TQ-API-003 | No error/validation response format | dev-blocking | **resolved** (TDEC-027) |
| TQ-API-004 | Pagination contract | needs-owner | **resolved** (TDEC-002) |
| TQ-API-005 | Filtering contract | deferred | deferred |
| TQ-API-006 | Idempotency guarantees for mutations | needs-owner | **resolved** (TDEC-003) |
| TQ-API-007 | No file upload contract | dev-blocking | **resolved** (TDEC-028) |
| TQ-API-008 | No session auth contract for API | dev-blocking | **resolved** (TDEC-029) |
| TQ-API-009 | API versioning strategy | needs-owner | **resolved** (TDEC-004) |
| TQ-API-010 | Chart data aggregation query undefined | dev-blocking | **resolved** (TDEC-030) |
| TQ-API-011 | Backup import flow API contract undefined | dev-blocking | **resolved** (TDEC-031) |
| TQ-API-012 | AI export operation contract undefined | dev-blocking | **resolved** (TDEC-032) |
| TQ-API-013 | No health check endpoint defined | watchlist | open |