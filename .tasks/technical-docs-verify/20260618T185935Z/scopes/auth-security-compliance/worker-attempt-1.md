# Auth Security Compliance - Worker Attempt 1

## Sources Read
- docs/product-verified/actors-and-permissions.md
- docs/product-verified/scope.md
- docs/product-verified/domain-model.md
- docs/product-verified/product-brief.md
- docs/product-verified/edge-cases.md
- docs/product-verified/business-rules.md
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/functional-spec.md
- docs/product-verified/user-flows.md

## Source Delta Reviewed
- DEC-007: single-user with multi-user-ready data model, userId on all entities, default user at bootstrap
- Implicit single-user identity: DefaultUser entity created at bootstrap, all entities carry userId FK
- DEC-007 means no registration flow, no multi-tenant isolation logic in MVP — but userId column exists for future multi-user

## Product Signals

### Identity
- DefaultUser entity at bootstrap (id, displayName, createdAt, updatedAt)
- No registration, no login, no user creation flow
- Single user per instance, single identity

### Authentication
- Optional PIN guard (AC-029–AC-034)
- PIN disabled: no access control at all (RULE-022, EDGE-013)
- PIN enabled: cookie-based session required for all pages (RULE-023)
- PIN stored as hash, not plaintext (RULE-001, AC-031)
- PIN change requires current PIN (RULE-002, AC-032, AC-109)
- Session via cookie (AC-034)
- Media files require valid session (RULE-024, AC-111)
- PIN forgotten: no recovery mechanism (user-flows failure flows)

### Authorization
- MVP has no role system (actors-and-permissions.md)
- Single user has all permissions (actors-and-permissions.md)
- All entities belong to user via userId FK
- No permission restrictions beyond optional PIN (actors-and-permissions.md)

### Ownership
- All data belongs to single user
- userId FK on all entities: Settings, UserProfile, Exercise, ExerciseMedia, DailyLog, WorkoutExercise, CardioEntry, BodyWeightEntry, BodyCheckIn, NutritionProduct, NutritionTemplate, DailyNutritionOverride, WeekFlag, AiExport, AiReview
- No ownership checking needed in MVP (all data is the user's)

### Tenant Scoping
- Single-tenant by design
- DEC-007 userId on all entities is future-proofing, no multi-tenant scoping logic in MVP

### Audit
- No audit trail mentioned in any product-verified doc
- No audit log entity
- No change tracking on sensitive operations (PIN changes, backup imports, exercise deletions)
- No login attempt logging

### Rate Limiting
- No rate limiting mentioned
- No brute-force protection on PIN attempts
- No API rate limiting

### Abuse Prevention
- PIN brute force: no protection specified
- Backup import abuse: no size limits or validation beyond manifest schema check
- No request throttling

### Secrets
- PIN hash: algorithm unspecified
- Session secret: not mentioned
- Cookie signing: not mentioned
- Cookie security flags (Secure, HttpOnly, SameSite): not mentioned

### Privacy
- No PIN logging (RULE-117/AC-117)
- No AI export content logging (RULE-118/AC-118)
- No photo logging (RULE-119/AC-119)
- No sensitive comment logging (RULE-120/AC-120)
- Photos in AI export require opt-in (RULE-025/AC-112)
- Media not accessible without valid session (RULE-024)
- Data retention/deletion: no policy (EDGE-027)
- No data anonymization or purging mechanism

### Compliance
- No regulatory compliance requirements mentioned (GDPR, CCPA, etc.)
- Privacy expectations documented but no formal compliance obligations

### Irreversible Action Controls
- Exercise deletion: no soft-delete, no recovery (EDGE-018)
- Media deletion: no confirm/recovery (EDGE-020)
- Backup import: dry-run validation exists before actual restore (RULE-008)
- No undo/trash mechanism for any entity

## Technical Facts

1. Identity is bootstrap-only: DefaultUser created at app startup with a fixed or configured displayName
2. Session is Redis-backed (scope.md dependencies: Redis for session store)
3. PIN hashing is required, algorithm not specified
4. Cookie-based session with unspecified TTL
5. All entities carry userId FK (DEC-007)
6. No auth middleware needed when PIN is disabled — app is fully open
7. When PIN is enabled, middleware check on every route except the PIN entry page
8. Media serving must check session even when PIN is disabled (contradiction: RULE-022 says no checks when PIN disabled, but RULE-024 says media needs session — RULE-024 must mean "when PIN enabled" or media gets special treatment)
9. Backup/import has no auth-scoped controls beyond the PIN gate itself

## Technical Gaps

### TGAP-AUTH-001: PIN Implementation Contract Missing
- Hash algorithm (bcrypt? argon2? scrypt? PBKDF2?)
- Salt handling
- PIN minimum/maximum length
- PIN complexity requirements (if any)
- PIN storage schema (existing Settings.pinHash field, but no format spec)
- How PIN is verified at login vs session creation

### TGAP-AUTH-002: Session Management Contract Missing
- Session TTL (timeout duration)
- Session renewal policy (sliding vs fixed)
- Session cookie name, path, Secure/HttpOnly/SameSite flags
- Logout mechanism (clear session cookie + delete server session)
- Session invalidation on PIN change/disable
- Maximum concurrent sessions (single user, likely 1, but no spec)
- Session persistence in Redis (key format, TTL strategy)

### TGAP-AUTH-003: PIN Brute Force Protection Missing
- No rate limiting on PIN attempts
- No account lockout after N failed attempts
- No progressive delay on failed attempts
- No audit logging of failed PIN attempts

### TGAP-AUTH-004: Audit Trail Missing
- No audit log entity or schema
- No requirement to log PIN changes, setting changes, backup imports, exercise deletions
- No requirement to log failed authentication attempts
- No requirement to log sensitive data access

### TGAP-AUTH-005: Cookie Security Specification Missing
- Secure flag: required for HTTPS deployments?
- HttpOnly: required to prevent XSS access?
- SameSite: Lax/Strict/None?
- Cookie name convention
- Session token generation (crypto-random? JWT? opaque token?)

### TGAP-AUTH-006: Media Access Control Implementation Missing
- Media serving endpoint must check session when PIN is enabled
- Media serving when PIN is disabled: should it check session?
- EDGE-014: direct media URL access without valid session — implementation unclear when PIN is disabled
- Media file naming: should use unpredictable names or rely on session check?

### TGAP-AUTH-007: Redis Session Store Failure Mode Missing
- EDGE-023: Redis unavailable — what happens?
- Fallback to in-memory session? Error page? Allow PIN re-entry to re-establish session on Redis recovery?
- Session store configuration

### TGAP-AUTH-008: Backup/Restore Auth Controls Missing
- Import validation: should verify user identity before restore?
- Export download: should verify session?
- No cross-instance restore protection (importing backup from another instance is inherently allowed per product docs)

### TGAP-AUTH-009: No Data Retention/Deletion Policy
- EDGE-027: no retention or deletion policy
- User has no way to delete all their data programmatically
- Media files accumulate unbounded (EDGE-030)
- No "delete my account/data" feature

### TGAP-AUTH-010: DefaultUser Bootstrap Contract Missing
- How is DefaultUser.displayName determined? Hardcoded? Configurable via env?
- DefaultUser.id: fixed (e.g., 'default') or generated UUID?
- Relationship between DefaultUser and UserProfile: separate entities, both carry userId
- When PIN is enabled, does the user authenticate against DefaultUser? Or is DefaultUser just a data ownership marker?

## Missing Source Artifacts
- Auth specification (auth flow, session lifecycle, PIN hashing algorithm, cookie config)
- Audit log schema
- Rate limiting configuration
- Session store configuration

## Questions Raised

| ID | Scope | Severity | Parent | Question | Why It Matters | Needed Artifact Or Decision | Source Or Report |
| --- | --- | --- | --- | --- | --- | --- | --- |
| TQ-AUTH-001 | auth-security-compliance | dev-blocking | none | What PIN hashing algorithm, salt strategy, and min/max length are required? | Blocks PIN implementation and validation logic. | Auth specification with algorithm choice, salt policy, length constraints. | actors-and-permissions.md (pin hash required, no algorithm) |
| TQ-AUTH-002 | auth-security-compliance | dev-blocking | none | What is the session TTL, renewal policy, cookie security flags (Secure/HttpOnly/SameSite), and logout mechanism? | Blocks session management and PIN guard implementation. | Session specification with TTL, renewal rules, cookie config, logout endpoint. | functional-spec.md (session via cookie, no details) |
| TQ-AUTH-003 | auth-security-compliance | dev-blocking | none | What PIN brute force prevention is required? (rate limit, lockout, progressive delay) | Without protection, PIN can be trivially brute-forced. | Rate limiting / brute force specification. | edge-cases.md (no protection mentioned) |
| TQ-AUTH-004 | auth-security-compliance | needs-owner-decision | none | Is an audit trail required for PIN changes, setting changes, backup imports, and exercise deletions? | Missing audit makes irreversible actions untraceable. | Audit requirement decision. | No audit mentioned in any source. |
| TQ-AUTH-005 | auth-security-compliance | dev-blocking | none | How is the session token generated (crypto-random vs JWT) and what session store key format is used in Redis? | Blocks session middleware implementation. | Session token spec and Redis key schema. | scope.md (Redis for session, no spec) |
| TQ-AUTH-006 | auth-security-compliance | dev-blocking | none | When PIN is disabled, does media serving still require a valid session? (RULE-024 conflicts with RULE-022) | Contradictory rules: RULE-022 says no checks when PIN disabled, RULE-024 says media needs session. | Owner decision on media access when PIN is disabled. | business-rules.md RULE-022 vs RULE-024 |
| TQ-AUTH-007 | auth-security-compliance | dev-blocking | none | What happens when Redis is unavailable for session storage? (fallback, error, retry) | Blocks resilient session architecture. | Redis failure mode specification. | edge-cases.md EDGE-023 |
| TQ-AUTH-008 | auth-security-compliance | deferred | none | Should backup import include any user identity validation (e.g., verify the backup was exported by the same user/instance)? | Low risk in single-user MVP; owner can defer. | Deferral rationale. | domain-model.md, product-brief.md |
| TQ-AUTH-009 | auth-security-compliance | needs-owner-decision | none | Is any data retention or deletion policy required for MVP? (delete all data, media cleanup, user data purge) | Unbounded data growth and no user-initiated data deletion. | Data retention/deletion policy decision. | edge-cases.md EDGE-027, EDGE-030 |
| TQ-AUTH-010 | auth-security-compliance | dev-blocking | none | How is DefaultUser.displayName populated and is DefaultUser.id a fixed value or generated UUID? | Blocks bootstrap implementation. | DefaultUser bootstrap specification. | domain-model.md (DefaultUser entity, no bootstrap spec) |
| TQ-AUTH-011 | auth-security-compliance | dev-blocking | none | When PIN is enabled, does the user authenticate against DefaultUser credentials or is PIN purely a session-level gate without user identity? | Determines whether DefaultUser has a PIN field or PIN is global. | Authentication model specification. | domain-model.md (Settings.pinHash global, not per-user) |
| TQ-AUTH-012 | auth-security-compliance | deferred | none | Should PIN attempt failures be logged? If so, what log level and retention? | AC-120 says no sensitive content logging; PIN attempts may be sensitive. | Logging policy for auth events. | acceptance-criteria.md AC-117–AC-120 |

## Answer Effects
N/A — initial run, no prior answers.

## Risks
1. PIN can be trivially brute-forced without rate limiting (high risk for a privacy-focused app).
2. No session TTL defined means sessions may persist indefinitely (EDGE-015).
3. Cookie security flags unspecified risks session hijacking on non-HTTPS deployments.
4. Conflicting rules on media access when PIN disabled may leak media without auth.
5. No audit trail for irreversible actions (exercise/media deletion, backup import, PIN change).

## Suggested Decisions
1. Use bcrypt for PIN hashing (standard, well-supported, configurable cost).
2. Session TTL: 24 hours sliding, cookie HttpOnly + SameSite=Lax.
3. Crypto-random opaque session tokens, stored in Redis with key format `session:{token}`.
4. Media access requires valid session regardless of PIN state (protects RULE-024).
5. PIN brute force: 5 attempts per minute per IP, 10 total failures locks PIN entry for 1 hour.
6. DefaultUser.id: use a fixed UUID generated at first bootstrap, stored in config or DB.
7. No audit log for MVP (deferred), but seed the Settings entity with an optional `auditEnabled` field for future.
8. Data retention: add a settings option for auto-delete media older than N months (deferred in MVP).

## Traceability Candidates
- Settings.pinHash → TQ-AUTH-001 (algorithm needed)
- Settings.pinEnabled → TQ-AUTH-002, TQ-AUTH-003 (session and brute force)
- edge-cases.md EDGE-011, EDGE-012, EDGE-013, EDGE-014, EDGE-015 → TQ-AUTH-001 through TQ-AUTH-007
- business-rules.md RULE-001 → TQ-AUTH-001
- business-rules.md RULE-022, RULE-023, RULE-024 → TQ-AUTH-006
- domain-model.md DefaultUser entity → TQ-AUTH-010, TQ-AUTH-011
- acceptance-criteria.md AC-029–AC-034, AC-109–AC-111, AC-117–AC-120 → TQ-AUTH-001 through TQ-AUTH-009