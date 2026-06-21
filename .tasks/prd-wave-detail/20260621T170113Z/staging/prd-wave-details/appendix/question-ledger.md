# Question Ledger

## Open Questions
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
|---|---|---|---|---|---|---|---|---|---|---|
| DQ-W07-001 | WAVE-07 | architecture-codebase-fit | Medium | None | Schema version format for manifest.json — integer (1) or semver (1.0.0)? | Schema evolution compatibility | Adopt integer schemaVersion = 1 | planner-architecture-codebase-attempt-1 Q-W07-002 | open | MVP: integer 1 |
| DQ-W07-002 | WAVE-07 | data-api-integration-ops | Medium | None | How to inject app version into manifest.json? | No existing build version injection | Add -ldflags or omit for MVP | planner-architecture-codebase-attempt-1 Q-W07-006 | open | Omit for MVP |
| DQ-W07-003 | WAVE-07 | security-privacy-compliance | Low | None | Max AiExport records per user — unbounded or capped? | Disk usage unbounded | Configurable max_records_per_user | planner-product-ac-attempt-1 Q-W07-004 | open | Follow-up after MVP |
| DQ-W07-004 | WAVE-07 | data-api-integration-ops | Medium | None | ZIP streaming threshold — when to switch from in-memory to streaming? | Large exports could OOM | Set threshold at 100MB | planner-data-integration-ops-attempt-1 Q-W07-DIO-03 | open | Implement size check |
| DQ-W07-005 | WAVE-07 | data-api-integration-ops | Low | None | Photo naming convention in export ZIP? | Preserve context vs avoid collisions | Use {checkInId}_{angle}.{ext} | planner-product-ac-attempt-1 Q-W07-005 | open | Descriptive naming |
| DQ-W07-006 | WAVE-07 | architecture-codebase-fit | Low | None | Build WeekFlagsByDateRange query or let client call per week? | N+1 queries for frontend | Build or defer | planner-sequencing-fit-attempt-1 §1.1 | open | Defer for MVP |

## Answered Questions
None — first detail run.

## Follow-Up Questions
- Q-1RM-FORMULA: Which 1RM formula to use? (WAVE-06 scope, not WAVE-07)
- Q-PIN-001: PIN rate limiting implementation? (WAVE-01 scope, not WAVE-07)
- Q-WORKOUT-001: Concurrent edit handling? (WAVE-03 scope, not WAVE-07)

## Resolved Questions
| ID | Question | Resolution | Evidence |
|---|---|---|---|
| DQ-W07-S1 | UserProfile vs Settings: new table or extend Settings? | Create separate user_profiles table. Settings=app config; UserProfile=user data. | DDEC-W07-001, domain-model.md |
| DQ-W07-S2 | Export lifecycle: 7-day TTL or keep-until-replaced? | 7-day TTL + delete-on-regeneration | DDEC-W07-007 |
| DQ-W07-S3 | include_photos default value? | DEFAULT false per RULE-025 | DDEC-W07-004, domain-model.md invariant #10 |
| DQ-W07-S4 | Migration numbers for WAVE-07 tables? | 00091_user_profiles.sql, 00092_ai_exports.sql | DDEC-W07-005 |
| DQ-W07-S5 | display_name in user_profiles? | NOT included. Use atlas_users.display_name. | DDEC-W07-015 |
| DQ-W07-S6 | CAP-W07-003 scope collision with WAVE-04? | REMOVED from WAVE-07. WAVE-04 owns week flags CRUD. | DDEC-W07-002 |
| DQ-W07-S7 | Sync vs async ZIP generation? | Sync for MVP. Architecture supports future async. | DDEC-W07-013 |
| DQ-W07-S8 | Max export size limit? | 100MB hard limit. Configurable. | DDEC-W07-009 |
| DQ-W07-S9 | Temp-file-atomic-rename pattern? | Yes — use for ZIP generation (EDGE-024) | DDEC-W07-008 |
| DQ-W07-S10 | Storage path pattern? | {ExportBasePath}/{userId}/{exportId}.zip | DDEC-W07-006 |

## Deferred Questions
- Multi-user future scope (out of MVP scope)
- OpenAI API integration (excluded per RULE-029)
- Apple Health integration (deferred per PDEC)
- Barcode scanner (deferred per PDEC)