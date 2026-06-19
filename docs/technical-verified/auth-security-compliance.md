# Auth Security Compliance

## Identity

PIN-based access with session via cookie. Decisions:
- **TDEC-004**: Minimal audit trail for PIN events (enabled/changed/disabled, failed attempt, successful unlock)
- **TQ-AUTH-001**: PIN hash algorithm still unspecified — open
- **TQ-AUTH-002**: Session TTL, cookie flags, renewal still unspecified — open
- **TQ-AUTH-003**: Brute-force protection — open
- **TQ-AUTH-005**: Session token generation — open

## Authorization And Ownership

Single-user with userId FK on all entities. API resolvers must scope all queries to the default user. No access to other users' data.

Authorization gaps:
- PIN is potentially global per instance vs per-app (TQ-AUTH-011)
- Media access contradiction: RULE-022 says no auth when PIN disabled, RULE-024 says media needs valid session (TQ-AUTH-006)
- Session token generation mechanism undefined (TQ-AUTH-005)

## Auditability

**Decision (TDEC-004):** Minimal audit trail required in MVP.

Audit events: PIN enabled/changed/disabled, failed PIN attempt, successful PIN unlock, media uploaded/deleted, AI export generated, full backup generated, import dry-run/commit/failed, full data deletion.

Do not log: raw PIN, PIN hash, body photos, AI prompt/export contents, health/body comments, user private notes, nutrition details, backup archive contents.

Audit log fields: event type, timestamp, default user id, request id, operation id, success/failure status, entity id where safe.

## Privacy And Compliance

- **Decision (TDEC-005):** No automatic data deletion in MVP. Data remains until user deletes manually.
  - Deleting an exercise → soft-disable if referenced by past workouts
  - Deleting media → remove metadata + physical file
  - Deleting body check-in → delete measurements + photo metadata/files
  - Deleting daily log → delete workout exercises, sets, cardio entries
  - Deleting nutrition template → keep product records
  - Deleting product referenced by history → prevent or soft-disable
- Backup files are user-controlled local files. Atlas not responsible for deleting downloaded archives.
- AI export generated files: 24-hour temporary retention unless user explicitly saves.

## Security Questions

| ID | Question | Severity | Status |
| --- | --- | --- | --- |
| TQ-AUTH-001 | PIN hash algorithm unspecified | dev-blocking | **resolved** (TDEC-033) |
| TQ-AUTH-002 | PIN session TTL, cookie flags, renewal policy undefined | dev-blocking | **resolved** (TDEC-034) |
| TQ-AUTH-003 | No brute-force protection for PIN | dev-blocking | **resolved** (TDEC-035) |
| TQ-AUTH-004 | Audit trail for sensitive operations | needs-owner | **resolved** (TDEC-004) |
| TQ-AUTH-005 | Session token generation mechanism undefined | dev-blocking | **resolved** (TDEC-036) |
| TQ-AUTH-006 | Media access contradiction: RULE-022 vs RULE-024 | dev-blocking | **resolved** (TDEC-037) |
| TQ-AUTH-007 | Redis session store failure mode undefined | dev-blocking | **resolved** (TDEC-038) |
| TQ-AUTH-008 | Backup identity validation (who can import into an instance) | deferred | deferred |
| TQ-AUTH-009 | No data retention/deletion policy | needs-owner | **resolved** (TDEC-005) |
| TQ-AUTH-010 | DefaultUser bootstrap security (when/how is default user created) | dev-blocking | **resolved** (TDEC-039) |
| TQ-AUTH-011 | PIN: global per instance or per-app? | dev-blocking | **resolved** (TDEC-040) |