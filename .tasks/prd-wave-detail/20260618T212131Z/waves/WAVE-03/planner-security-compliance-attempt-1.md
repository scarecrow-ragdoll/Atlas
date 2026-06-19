# WAVE-03 security-compliance Planner Attempt 1

## Sources Read
- docs/technical-verified/auth-security-compliance.md
- docs/product-verified/business-rules.md
- docs/product-verified/edge-cases.md
- docs/prd-wave-details/waves/wave-02.md
- docs/prd-wave-details/waves/wave-01.md

## Selected Backend Wave Boundary
WAVE-03 is entirely backend GraphQL operations for workout diary. All endpoints must be protected by WAVE-01 PIN auth middleware. No REST endpoints, no binary uploads, no file storage.

## Neighboring Backend Wave Fit
- WAVE-01: Provides PIN auth middleware, session validation, and common error types (AuthError). WAVE-03 reuses all.
- WAVE-02: Provides allExercises query protected by same PIN auth. WAVE-03 consumes it from service layer.

## Frontend Pages Context
- PAGE-002: consumes WAVE-03 GraphQL endpoints via PIN-authenticated requests. No security concerns beyond standard PIN auth.

## Proposed Details

### Authentication
- All WAVE-03 GraphQL queries and mutations are protected by WAVE-01 PIN auth middleware
- PIN session token in Authorization header (Bearer <token>) per TDEC-029
- When PIN is disabled: endpoints accessible without auth (consistent with TDEC-037, RULE-022)
- When PIN is enabled: all requests require valid session; invalid/expired session returns AuthError
- No additional authentication mechanism needed

### Authorization
- Single-user (MVP constraint): all data owned by default user
- Service layer scopes all queries by user_id from PIN context (TDEC-015)
- exercise_id FK references exercises table — no cross-user access possible (single-user only)
- No role-based or permission-based authorization needed

### Privacy
- WorkoutExercise.notes (comments) not logged in application logs per TDEC-004
- WorkoutSet.notes (per-set comments) not logged
- CardioEntry.notes not logged
- DailyLog.notes not logged
- Log markers record only: operation type, entity id, user id, success/failure status
- No personal health information (PHI) concerns — data is stored locally, not transmitted

### Audit Trail
- Audit events per TDEC-004: DailyLog created/updated/deleted, WorkoutExercise added/removed, WorkoutSet added/updated/removed, CardioEntry added/removed
- Audit fields: event type, timestamp, user id, request id, success/failure, entity id
- No sensitive content in audit logs

### Rate Limiting
- Deferred from WAVE-01 (DQ-W01-001). No rate limiting in WAVE-03 MVP.
- Standard HTTP-level rate limiting can be applied later via middleware.

### Input Validation
- All numeric inputs validated server-side in service layer:
  - Set weight: >= 0 (0 for bodyweight exercises allowed)
  - Set reps: positive integer (> 0)
  - RPE: 1.0-10.0 scale (optional, step 0.5)
  - RIR: 0-5 integer scale (optional)
  - Duration minutes: positive integer
  - Heart rate zone: 1-5 (optional)
  - Avg pulse: positive integer (optional)
- Date format: validated as DATE type (YYYY-MM-DD)
- exercise_id: must reference existing active exercise (validated against WAVE-02 allExercises or exerciseById query)
- String fields (notes): max length 1000 characters to prevent abuse

### Data Integrity
- FK constraints enforce referential integrity:
  - workout_exercises.daily_log_id -> daily_logs(id) CASCADE
  - workout_exercises.exercise_id -> exercises(id) NO ACTION (preserves workout history)
  - workout_sets.workout_exercise_id -> workout_exercises(id) CASCADE
  - cardio_entries.daily_log_id -> daily_logs(id) CASCADE
- UNIQUE(user_id, date) on daily_logs prevents duplicate date entries
- Working weight snapshot captured at add time; not subject to concurrent update issues in MVP

### Secrets
- No new secrets introduced by WAVE-03
- Existing WAVE-01 PIN hash (bcrypt) and Redis session key are the only secrets

### Abuse Prevention
- No user-generated file uploads in WAVE-03 — no file validation or storage concerns
- String field length limits prevent abusive payloads
- GraphQL query depth/complexity limiting: deferred — can be added as middleware in future

### Security Questions Status
- TQ-AUTH-001 through TQ-AUTH-011: all resolved in technical-verified docs. WAVE-01 provides the implementation.
- No new security questions raised by WAVE-03 scope.

## Acceptance Criteria Contributions
- AC-W03-029: AuthError for invalid/missing PIN on all mutations
- AC-W03-030: AuthError for invalid/missing PIN on queries

## Exit Criteria Contributions
- EC-W03-011: All WAVE-03 mutations return AuthError without valid PIN session
- EC-W03-012: No sensitive content in application logs for WAVE-03 operations
- EC-W03-013: Input validation enforced for all numeric and string fields

## Verification Contributions
- TEST-W03-018: Auth protection for all WAVE-03 GraphQL operations
- TEST-W03-019: Input validation (weight, reps, RPE, RIR, duration bounds)
- TEST-W03-020: Log privacy — notes not appearing in logs
- TEST-W03-021: FK constraint enforcement (invalid exercise_id, invalid daily_log_id)

## Risks And Rollback
- Risk: If WAVE-01 PIN auth is not yet implemented, WAVE-03 cannot be deployed. This is an acknowledged blocking dependency.
- Risk: Input validation rules for RPE/RIR need alignment with exercise science standards. Conservative bounds used.
- Rollback: All security measures are code-based (no security-critical config changes). Rollback removes the security posture along with the feature.

## Questions Raised
- None new. All security patterns inherited from WAVE-01 decisions.

## Traceability Candidates
- docs/technical-verified/auth-security-compliance.md: TDEC-004, TDEC-029, TDEC-037
- docs/product-verified/business-rules.md: RULE-004, RULE-022, RULE-023, RULE-024
- docs/product-verified/edge-cases.md: EDGE-012, EDGE-013
