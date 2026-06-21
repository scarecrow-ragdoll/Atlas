# Open Questions

## Wave-Blocking
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
|---|---|---|---|---|---|---|---|---|---|---|

## Needs Owner Decision
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
|---|---|---|---|---|---|---|---|---|---|---|
| DQ-W07-001 | WAVE-07 | architecture-codebase-fit | Medium | None | Schema version format for manifest.json — integer (1) or semver (1.0.0)? | Schema evolution compatibility for downstream consumers | Adopt integer schemaVersion = 1 | planner-architecture-codebase-attempt-1 Q-W07-002 | open | MVP: integer 1 |
| DQ-W07-002 | WAVE-07 | data-api-integration-ops | Medium | None | How to inject app version into manifest.json? | No existing build version injection in codebase | Add -ldflags for main.appVersion or omit from manifest | planner-architecture-codebase-attempt-1 Q-W07-006 | open | Omit appVersion for MVP |
| DQ-W07-003 | WAVE-07 | security-privacy-compliance | Low | None | Max AiExport records per user — unbounded or capped? | Disk usage unbounded; cleanup TTL handles old records but count not limited | Configurable max_records_per_user (default: 50) | planner-product-ac-attempt-1 Q-W07-004 | open | Follow-up: add after MVP |
| DQ-W07-004 | WAVE-07 | data-api-integration-ops | Medium | None | ZIP streaming threshold — when to switch from in-memory to streaming? | Large exports with many photos could OOM if built fully in memory | Set threshold: if estimated size > 100MB, stream to temp file | planner-data-integration-ops-attempt-1 Q-W07-DIO-03 | open | Implement size check before build |
| DQ-W07-005 | WAVE-07 | data-api-integration-ops | Low | None | Photo naming convention in export ZIP — UUID-based or descriptive? | Descriptive names preserve context; UUID avoids collisions | Use {checkInId}_{angle}.{ext} | planner-product-ac-attempt-1 Q-W07-005 | open | Adopt descriptive naming |
| DQ-W07-006 | WAVE-07 | architecture-codebase-fit | Low | None | Build WeekFlagsByDateRange query or let client call per week? | Reduces N+1 queries for frontend | Build lightweight query or defer | planner-sequencing-fit-attempt-1 §1.1 | open | Defer: client calls per week for MVP |

## Deferred
None.

## Watchlist
- WAVE-03 workout data stub pattern — verified against WAVE-06 precedent
- ai_exports table generated_prompt field character limit — add 5000 char truncation
- Concurrent generation guard — frontend debouncing sufficient for MVP

## Resolved This Run
- DQ-W07-001 (UserProfile vs Settings): Resolved — Create separate UserProfile table. See DDEC-W07-001.
- DQ-W07-002 (lifecycle 7-day TTL vs keep-until-replaced): Resolved — 7-day TTL + delete-on-regeneration. See DDEC-W07-007.
- CAP-W07-003 (week flags CRUD): Removed from WAVE-07. WAVE-04 owns it. See DDEC-W07-002.
- include_photos DEFAULT true bug: Resolved — DEFAULT false. See DDEC-W07-004.
- Migration numbering: Resolved — 00091, 00092. See DDEC-W07-005.
- display_name in user_profiles: Resolved — NOT included. Use atlas_users. See DDEC-W07-015.