# Decision Log
## Source Wave Gate
- WAVE-04 source wave gate: passed (2026-06-19)
- Source: docs/prd-waves/waves/wave-04.md (user-approved 2026-06-18)
- WAVE-05 source wave gate: passed (2026-06-19)
- Source: docs/prd-waves/waves/wave-05.md (user-approved 2026-06-18)
## User Wave Approvals
- WAVE-04 source wave: user-approved (2026-06-18) by original wave map approval
- WAVE-04 detailed wave: questions-open — awaiting owner decision on DQ-W04-001, then user approval
- WAVE-05 source wave: user-approved (2026-06-18) by original wave map approval
- WAVE-05 detailed wave: ready-for-dev — awaiting user approval
## Scope Decisions
- DDEC-W04-001: Hard delete for all WAVE-04 entities (no soft delete like WAVE-02)
- DDEC-W04-002: Photo count 2-4 recommended (soft guidance), hard limit 10
- DDEC-W04-003: Measurement side allowed only for paired types (forearm, biceps, thigh, calf)
- DDEC-W04-004: BodyWeightEntry allows multiple entries per date (different sources)
- DDEC-W04-005: DailyLog auto-creation for CardioEntry (requires daily_log table)
- DDEC-W05-001: Template upsert by (userId, weekStartDate) — per-week, one template per week
- DDEC-W05-002: Server-side macro calculation via separate nutritionMacros query
- DDEC-W05-003: Soft-delete for NutritionProduct with isActive flag
- DDEC-W05-004: Free-text mealLabel (not enum)
- DDEC-W05-005: Single migration for all 5 nutrition tables
## Codebase Fit Decisions
- Follow WAVE-02 patterns for repository, service, resolver, handler, and migration structure
- ProgressPhoto REST handler follows WAVE-02 exercise_media.go pattern (multipart upload, MIME/size validation, physical file deletion)
- GraphQL schema files in libs/graphql/schema/ with extend type Query/Mutation pattern
- sqlc queries in apps/api/internal/repository/postgres/queries/
- Migrations start at 00082 (following WAVE-02's 00080-00081)
- WAVE-05 uses Atlas module patterns (apps/api/internal/atlas/) for all nutrition entities
- Atlas GraphQL schema files in apps/api/internal/atlas/graph/schema/
- Nutrition services follow settings_service.go pattern (interface + private struct)
## Deferrals
- Signed URLs for progress photos: deferred post-MVP
- WAVE-01 MediaConfig BasePath confirmation: deferred until WAVE-01 implementation
- Soft-deleted product recovery via API: deferred — admin-only DB recovery for MVP
- Migration number coordination: deferred — check at implementation time
## Rejected Assumptions
- Photos can only be uploaded, not taken in-app (confirmed by source wave excluded scope)
- Measurements use REAL type (not INT) for body measurement values
- CardioEntry requires DailyLog attachment (domain model invariant, not optional)
- Macro calculation is server-side, not client-side (DDEC-W05-002)