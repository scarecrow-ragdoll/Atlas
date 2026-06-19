# Decision Log
## Source Wave Gate
- WAVE-04 source wave gate: passed (2026-06-19)
- Source: docs/prd-waves/waves/wave-04.md (user-approved 2026-06-18)
## User Wave Approvals
- WAVE-04 source wave: user-approved (2026-06-18) by original wave map approval
- WAVE-04 detailed wave: awaiting user approval after ready-for-dev
## Scope Decisions
- DDEC-W04-001: Hard delete for all WAVE-04 entities (no soft delete like WAVE-02)
- DDEC-W04-002: Photo count 2-4 recommended (soft guidance), hard limit 10
- DDEC-W04-003: Measurement side allowed only for paired types (forearm, biceps, thigh, calf)
- DDEC-W04-004: BodyWeightEntry allows multiple entries per date (different sources)
- DDEC-W04-005: DailyLog auto-creation for CardioEntry (requires daily_log table)
## Codebase Fit Decisions
- Follow WAVE-02 patterns for repository, service, resolver, handler, and migration structure
- ProgressPhoto REST handler follows WAVE-02 exercise_media.go pattern (multipart upload, MIME/size validation, physical file deletion)
- GraphQL schema files in libs/graphql/schema/ with extend type Query/Mutation pattern
- sqlc queries in apps/api/internal/repository/postgres/queries/
- Migrations start at 00082 (following WAVE-02's 00080-00081)
## Deferrals
- Signed URLs for progress photos: deferred post-MVP
- WAVE-01 MediaConfig BasePath confirmation: deferred until WAVE-01 implementation
## Rejected Assumptions
- Photos can only be uploaded, not taken in-app (confirmed by source wave excluded scope)
- Measurements use REAL type (not INT) for body measurement values
- CardioEntry requires DailyLog attachment (domain model invariant, not optional)