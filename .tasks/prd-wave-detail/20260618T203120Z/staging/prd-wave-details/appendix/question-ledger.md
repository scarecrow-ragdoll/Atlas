# Question Ledger

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W01-001 | WAVE-01 | security | deferred | Q-PIN-001 | PIN rate limiting implementation? | Security against brute force attacks | Decide rate limit strategy (Redis TTL, fixed delay, lockout) | docs/technical-verified/auth-security-compliance.md | open | deferred |

## Answered Questions
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W01-003 | WAVE-01 | api | answered | TQ-API-001 | Primary API protocol undecided | Architecture decision | Confirm GraphQL vs REST | docs/technical-verified/api-contracts.md | answered | Hybrid model (TDEC-001) |
| DQ-W01-004 | WAVE-01 | auth | answered | TQ-AUTH-001 | Session auth contract | Security | Confirm session approach | docs/technical-verified/auth-security-compliance.md | answered | PIN session with Redis (TDEC-029) |

## Follow-Up Questions
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W01-005 | WAVE-01 | security | watchlist | DQ-W01-001 | First PIN setup flow | UX and security should be decided | Confirm first-time PIN creation UX flow | docs/technical-verified/auth-security-compliance.md | open | watchlist |

## Resolved Questions
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W01-006 | WAVE-01 | media | resolved | TQ-API-007 | File upload contract | Development can proceed | Define upload/download API | docs/technical-verified/api-contracts.md | resolved | REST multipart, local storage (TDEC-028) |

## Deferred Questions
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |