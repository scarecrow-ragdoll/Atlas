# WAVE-02 security-compliance Planner Attempt 1

## Sources Read
- docs/technical-verified/auth-security-compliance.md (PIN auth, audit, privacy, media access)
- docs/product-verified/business-rules.md (RULE-022, RULE-023, RULE-024)
- docs/product-verified/edge-cases.md (EDGE-011, EDGE-012, EDGE-014, EDGE-020)
- docs/prd-wave-details/waves/wave-01.md (WAVE-01 security: PIN bcrypt, session TTL, rate limiting deferred)
- apps/api/internal/middleware/admin_auth.go (auth middleware pattern)
- apps/api/internal/graph/admin_auth_helpers.go (requireAdmin guard)

## Selected Backend Wave Boundary
WAVE-02 exercises and exercise media are protected by WAVE-01 PIN auth. No new authentication mechanism. Authorization is single-user (default user). Access control is: PIN enabled → all endpoints require valid session; PIN disabled → no access control (RULE-022).

## Neighboring Backend Wave Fit
- WAVE-01: provides PIN auth guard middleware for fitness-domain endpoints. WAVE-02 reuses it.
- Security gaps from WAVE-01 that affect WAVE-02: rate limiting deferred (DQ-W01-001), PIN session TTL configurable but default 7 days.

## Frontend Pages Context
PAGE-003: exercises accessible only through PIN-protected API. Media download requires valid session (RULE-024).

## Codebase Evidence
- No security-specific code for exercises yet. All WAVE-02 security comes from WAVE-01 PIN auth middleware.
- Existing admin auth uses separate middleware (admin_session) with cookie-backed sessions. WAVE-02 uses PIN auth from WAVE-01 (token-based, Authorization header).

## Proposed Details

### Authentication
- All WAVE-02 GraphQL operations (exercise CRUD) protected by WAVE-01 PIN auth middleware.
- All WAVE-02 REST operations (exercise media upload/delete) protected by WAVE-01 PIN auth middleware.
- Exercise media download (GET /api/v1/exercise-media/{id}) requires valid PIN session (RULE-024).
- When PIN is disabled, no authentication required (RULE-022) — endpoints are open but only accessible on local network per self-hosted deployment.

### Authorization
- Single-user context: no per-exercise ownership needed. All exercises belong to the default user.
- No role-based access for exercises.

### Input Validation Security
- SQL injection: prevented by sqlc parameterized queries.
- File upload validation: enforce allowed MIME types (JPEG/PNG/WEBP/MP4/MOV/WEBM) server-side. Do not trust client-provided content-type.
- File size validation: enforce TDEC-008 limits (25MB images, 250MB video) before reading into memory.
- File name sanitization: strip path separators, limit length, generate UUID-based storage filename to prevent path traversal.
- Exercise name: strip leading/trailing whitespace, limit length.
- All GraphQL inputs: gqlgen type system enforces scalar types.

### Privacy
- Exercise personalNotes: user-provided content stored in DB, transmitted through GraphQL. Not logged (follow pattern from WAVE-01 TDEC-004 — sensitive content not logged).
- Exercise media: image/video files stored on local filesystem, served only through authenticated REST endpoint. No public URL.
- Exercise media metadata (original file name, size, etc.) is stored — no PII concern for self-hosted single-user app.

### Audit
- Exercise CRUD operations: log event type, timestamp, success/failure, exercise ID. Do not log personalNotes content.
- Exercise media upload/delete: log event type, timestamp, success/failure, exercise ID, file size.
- Follow TDEC-004 audit pattern: [Exercise][create|update|delete][BLOCK_*] markers.
- Audit log format: [Exercise][action][BLOCK_NAME] + context fields (exercise_id, ok/error).

### Rate Limiting
- Rate limiting deferred from WAVE-01 (DQ-W01-001). WAVE-02 endpoints inherit the same lack of rate limiting.
- Low risk for single-user deployment.

### Abuse Vectors
- **Unlimited media uploads**: single user, but disk space could be exhausted. No quota enforcement in WAVE-02. Acceptable for MVP.
- **Large file uploads blocking server**: single request size limited to 300MB (TDEC-008). Standard Go HTTP read timeout prevents indefinite blocking.
- **Brute force exercise creation**: low-value target (no sensitive data). No rate limiting.

## Acceptance Criteria Contributions

| AC ID | Description |
| --- | --- |
| AC-W02-014 | Exercise GraphQL mutations return AuthError when PIN session is invalid |
| AC-W02-015 | Exercise media REST endpoints return 401 when PIN session is invalid |
| AC-W02-016 | File upload rejects files with disallowed MIME types (only JPEG/PNG/WEBP/MP4/MOV/WEBM) |
| AC-W02-017 | File upload rejects files larger than 25MB (images) or 250MB (video) |
| AC-W02-018 | Uploaded file name is sanitized (path traversal prevented, UUID-based storage name) |
| AC-W02-019 | No personalNotes content appears in application logs |
| AC-W02-020 | Exercise and ExerciseMedia operations appear in audit log markers |

## Exit Criteria Contributions

| EC ID | Description |
| --- | --- |
| EC-W02-011 | All exercise endpoints protected by WAVE-01 PIN auth (GraphQL + REST) |
| EC-W02-012 | File upload rejects invalid MIME types and sizes > limits |
| EC-W02-013 | No sensitive content (personalNotes, media content) logged |
| EC-W02-014 | Path traversal prevention verified for uploaded file names |

## Verification Contributions

| Test ID | Description | Type | Command |
| --- | --- | --- | --- |
| TEST-W02-014 | Exercise GraphQL returns AuthError without valid PIN session | integration | bunx nx run api:test -- --run '(?i)exercise_auth' |
| TEST-W02-015 | ExerciseMedia upload returns 401 without valid PIN session | integration | bunx nx run api:test -- --run '(?i)exercise_media_auth' |
| TEST-W02-016 | File type rejection for unauthorized MIME types | unit | bunx nx run api:test -- --run '(?i)exercise_media_filetype' |
| TEST-W02-017 | File size rejection for oversized uploads | unit | bunx nx run api:test -- --run '(?i)exercise_media_filesize' |
| TEST-W02-018 | Path traversal prevention in upload handler | unit | bunx nx run api:test -- --run '(?i)exercise_media_path_traversal' |
| TEST-W02-019 | Sensitive content not appearing in log output | unit | bunx nx run api:test -- --run '(?i)exercise_log_sanitize' |

## Risks And Rollback
- Risk: PIN disabled means no auth at all (RULE-022). Acceptable per product decision for self-hosted single-user deployment.
- Risk: file upload validation bypass with content-type spoofing. Mitigation: server-side MIME detection using file magic bytes (e.g., http.DetectContentType), not just Content-Type header.
- Risk: no rate limiting on media upload. Acceptable for single-user MVP.
- Rollback: same as WAVE-01 — revert code and migrations. No security state changes that persist after rollback.

## Questions Raised

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W02-005 | WAVE-02 | security-compliance | needs-owner-decision | TDEC-008 | Should WAVE-02 use server-side MIME detection (file magic bytes) or trust the Content-Type header for upload validation? | Content-Type can be spoofed; magic bytes are more secure but require reading file bytes | security-compliance planner | open |
| DQ-W02-006 | WAVE-02 | security-compliance | deferred | EDGE-014 | Should exercise media URLs be time-limited (signed URLs) or always accessible with valid session? | Signed URLs add complexity. Session-gated access may be sufficient for self-hosted single-user. | security-compliance planner | open |

## Traceability Candidates
- PIN protection → WAVE-01 PIN auth middleware, RULE-022, RULE-023, RULE-024
- File validation → TDEC-008 file size limits
- Audit logging → TDEC-004 audit trail
- Privacy (no logging sensitive content) → TDEC-004, AC-117, AC-118, AC-119, AC-120