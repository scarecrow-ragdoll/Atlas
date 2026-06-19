# Open Questions

## Wave-Blocking

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |

## Needs Owner Decision

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W01-001 | WAVE-01 | security | deferred | Q-PIN-001 | PIN rate limiting implementation? | Security against brute force attacks | Decide rate limit strategy (Redis TTL, fixed delay, lockout) | docs/technical-verified/auth-security-compliance.md | open | deferred |

## Deferred

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |

## Watchlist

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W01-002 | WAVE-01 | operations | watchlist | None | First-admin bootstrap for fitness app? | Template has admin seed, fitness app may need its own setup | Confirm if fitness app needs auto-seed or manual PIN setup | docs/technical-verified/auth-security-compliance.md | open | watchlist |

## Resolved This Run

- API protocol decision (TDEC-001): hybrid GraphQL/REST — confirmed for WAVE-01
- Session auth contract (TDEC-029): PIN session with Redis — confirmed
- Error format (TDEC-027): standard error response — confirmed
- No endpoint catalog needed for WAVE-01 (settings CRUD + PIN operations)
- Media upload path for later waves: POST /api/v1/media/upload + GET /api/v1/media/{id} REST endpoints scaffolded in WAVE-01