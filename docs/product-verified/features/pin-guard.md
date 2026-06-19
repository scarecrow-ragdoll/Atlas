# PIN Guard

## Source Evidence

PRD §7.2: Optional PIN code with hash storage, session via cookie, enable/disable/change.

## User Problem

Prevent unauthorized access to personal fitness data when the application URL is exposed.

## Scope

In MVP. Single-user, no registration alternative.

## Behavior

- PIN is optional; when disabled, app is open without any check
- When enabled, PIN entry required on app open
- PIN stored as hash; not logged
- User can change PIN (requires current PIN)
- User can disable PIN
- Session persists via cookie after PIN entry
- All sensitive data (including media) requires valid session

## Derived Requirements

| Requirement | Source | Rationale | Confidence |
| --- | --- | --- | --- |
| PIN hashing before storage | §7.2 | "PIN не должен храниться в открытом виде" | High |
| Cookie-based session after PIN entry | §7.2 | "сессия после ввода PIN должна сохраняться через cookie/session" | High |

## Edge Cases

- PIN enabled but pinHash missing from settings (EDGE-011)
- Session expired during active data entry (EDGE-012)
- PIN disabled — no access control (EDGE-013)

## Acceptance Criteria

AC-029 through AC-034; AC-109, AC-117.

## Dependencies

Settings entity, session storage (Redis in stack).

## Open Questions

Q-ROLE-001: PIN session lifetime and renewal policy.
Q-ROLE-002: Logout mechanism when PIN is enabled.
Q-ROLE-003: Access control when PIN is disabled.
Q-AC-01: PIN failure behavior (retries, lockout).
Q-EDGE-02: PIN brute-force protection policy.