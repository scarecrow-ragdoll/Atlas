# Decision Log

## Source Wave Gate
- WAVE-02 source wave gate: passed (2026-06-18T204251Z)
- Selected wave: WAVE-02 (docs/prd-waves/waves/wave-02.md)
- Source status: user-approved

## User Wave Approvals
- WAVE-02: not yet approved by user (wave is ready-for-dev awaiting review)

## Scope Decisions
| ID | Decision | Rationale | Source |
| --- | --- | --- | --- |
| DDEC-W02-001 | Soft delete (isActive=false) instead of hard delete | Preserves referential integrity for WAVE-03 workout history | CAP-W02-005, EDGE-018 |
| DDEC-W02-002 | ExerciseMedia file deletion removes physical file from disk | Consistent with TDEC-005 media deletion principle | TDEC-005, DQ-W02-001 resolved |
| DDEC-W02-003 | Duplicate exercise names allowed | PRD §11 does not specify uniqueness; EDGE-002 | EDGE-002 |
| DDEC-W02-004 | allExercises is a GraphQL query (not REST) | Consistent with hybrid pattern; frontend uses GraphQL | PAGE-002, TDEC-001 |

## Codebase Fit Decisions
- WAVE-02 uses WAVE-01's PIN auth middleware contract — assumes requirePinAuth(ctx) and middleware for chi
- WAVE-02 uses WAVE-01's media storage — assumes MediaConfig.BasePath and POST/GET /api/v1/media
- WAVE-02 reuses existing gqlgen and sqlc glob patterns — no config changes needed
- WAVE-02 routes exercise-media REST under PIN-protected sub-route group

## Deferrals
- Full-text search for exercise names: deferred post-MVP
- Signed/media URLs: deferred post-MVP (DQ-W02-006)
- WAVE-01 MediaConfig path confirmation: deferred to WAVE-01 implementation (DQ-W02-003)
- allExercises filtering beyond isActive: deferred to watchlist (DQ-W02-008)

## Rejected Assumptions
- pg_trgm extension: rejected — use simple B-tree indexes instead (planner-data-integration-ops attempt 2)
- ON DELETE CASCADE for exercise_media FK: rejected — use ON DELETE NO ACTION to prevent accidental cascade from hard deletes (planner-data-integration-ops attempt 2)
- Trusting Content-Type header for MIME validation: rejected — use http.DetectContentType() server-side (planner-security-compliance attempt 2)