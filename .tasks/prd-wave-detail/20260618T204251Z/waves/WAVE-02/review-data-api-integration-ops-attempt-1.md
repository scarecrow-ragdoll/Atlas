# WAVE-02 data-api-integration-ops Review Attempt 1

## Verdict
needs-revision

## Sources Read
- planner-data-integration-ops-attempt-1.md
- planner-architecture-codebase-attempt-1.md
- planner-product-ac-attempt-1.md
- planner-security-compliance-attempt-1.md
- planner-testing-exit-attempt-1.md
- docs/technical-verified/api-contracts.md
- docs/technical-verified/data-contracts.md
- docs/technical-verified/integrations-and-events.md (TDEC-008 file limits)
- docs/product-verified/domain-model.md
- docs/product-verified/business-rules.md

## Coverage Check
Data model (exercises + exercise_media tables), migrations, file storage, validation, error handling, observability all covered. The planner addresses the hybrid GraphQL/REST pattern correctly.

## Evidence Check
Data model aligns with domain model (Exercise, ExerciseMedia). Table columns match entity attributes. File validation limits match TDEC-008. All claims source-backed.

## Codebase Fit Check
SQL query patterns match existing .sql files (users.sql, admin_users.sql). The proposed exercise_repo.go follows user_repo.go pattern. Error handling follows existing pattern (pgx.ErrNoRows → nil, duplicate key → domain error).

### Issues Found

1. **Trigram extension dependency**: The planner proposes `pg_trgm` extension for name search indexing (`idx_exercises_name_trgm`). This requires `CREATE EXTENSION IF NOT EXISTS pg_trgm` in the migration. This extension may not be available in all PostgreSQL deployments. Should either (a) add extension creation to migration or (b) use a simpler index (e.g., `CREATE INDEX idx_exercises_name ON exercises (name)`). Recommendation: add extension to migration.

2. **REAL vs NUMERIC for working weight**: The planner uses `REAL` for working_weight. REAL is a 4-byte floating point which may cause precision issues. `NUMERIC(8,2)` or `REAL` with appropriate rounding is preferred. Since working weight is a human-readable value (kg/lbs), `REAL` is acceptable but should be explicitly decided.

3. **ON DELETE CASCADE risk**: The planner uses `ON DELETE CASCADE` on exercise_media.exercise_id FK. Since WAVE-02 uses soft delete, this cascade only applies to hard deletes (which don't happen in WAVE-02). However, if a future wave implements hard cleanup, media files would be lost. Consider `ON DELETE NO ACTION` or `ON DELETE SET NULL` instead to prevent accidental data loss. Recommend `ON DELETE NO ACTION` (safer default, explicit cleanup needed).

4. **File storage path pattern**: The planner says `exercise/<exercise_id>/<uuid>.<ext>` but doesn't confirm this aligns with WAVE-01's media storage path scheme. WAVE-01 might use `media/<uuid>.<ext>` without grouping. Need WAVE-01 path pattern verified.

5. **Missing REST endpoint for media download**: The planner lists POST and DELETE for exercise-media but doesn't specify the GET endpoint for downloading exercise media. PAGE-003 needs to display media. Where does GET /api/v1/exercise-media/{id} map to? WAVE-01's GET /api/v1/media/{id} may not know about exercise association. Either (a) reuse WAVE-01's GET /api/v1/media/{id} for all media (including exercise) or (b) add dedicated GET /api/v1/exercise-media/{id}. Clarify.

6. **Error envelope format**: The existing API uses `{ "error": { "code", "message", "field"? } }` per TDEC-027. The planner says "ValidationError with field 'file'" in GraphQL context but "404 error envelope" for REST. This is inconsistent — REST should also use the standard error envelope. Clarify that GraphQL uses union types (already correct), REST uses the JSON error envelope.

## Other-Wave Fit Check
WAVE-01 provides base path config. WAVE-03 uses the exercise data. The planner correctly identifies the data dependency boundaries.

## Acceptance Criteria Check
Not directly applicable, but AC-W02-006 (media upload) and AC-W02-007 (media delete) rely on the REST endpoints being correctly designed.

## Exit Criteria Check
EC-W02-010 (ON DELETE CASCADE) should be revised if the FK constraint changes to NO ACTION.

## Verification Check
TEST-W02-009 through TEST-W02-013 cover data operations well. File validation tests are correctly placed.

## Question Ledger Check
DQ-W02-003 (file storage path) and DQ-W02-004 (physical file deletion) are critical data-level questions. DQ-W02-003 blocks the migration design.

## Unsupported Or Invented Claims
The trigram extension and REAL type are valid choices but should be documented as decisions with rationale.

## Required Revisions
1. **Add pg_trgm extension or simplify index**: Document pg_trgm as a decision with fallback plan.
2. **Change ON DELETE CASCADE to NO ACTION**: Safer for soft-delete pattern, prevents accidental cascade deletion.
3. **Add GET endpoint for exercise media download**: Clarify whether WAVE-01's GET /media/{id} suffices or a dedicated endpoint is needed.
4. **Clarify working_weight data type**: REAL vs NUMERIC decision.
5. **Align REST error format with TDEC-027**: Consistent error envelope across all REST endpoints.

## Approval Notes
Data model is well-designed. The 5 revision items are all resolvable. After revisions, will approve.