# Decision Log

## Resolved Contradictions

| ID | Contradiction | Sources | Resolution | Confidence |
| --- | --- | --- | --- | --- |
| DEC-001 | "No registration / no auth" vs optional PIN guard | §7.1 vs §7.2 | PIN is de facto auth for single-user mode — the app is accessible without registration but can be PIN-protected. No contradiction: PIN is optional access control, not user registration. | High |
| DEC-002 | Cardio as part of workout day vs separate entity | §10.1, §10.3 vs §25.8 | Resolved by DEC-009: DailyLog replaces WorkoutDay; cardio is separate entity always attached to DailyLog via required dailyLogId. | High |
| DEC-003 | Telegram bot in technology stack vs not in MVP scope | Stack vs §23 | Telegram bot library is in stack for future use, not MVP. No contradiction — library inclusion does not imply MVP requirement. | High |
| DEC-004 | Apple Health listed in future scope but not in stack | §22 vs stack | Apple Health is explicitly future scope. No contradiction. | High |

## Blocking Questions Resolved

| ID | Question | Decision | Rationale |
| --- | --- | --- | --- |
| DEC-006 | Q-SCOPE-001: Success metrics | Functional success metrics, data completeness metrics, backup/restore metrics, and quality gates defined in product-brief.md §Success Metrics | Product owner decision (PRD patch 20260618) |
| DEC-007 | Q-SCOPE-002: Multi-user architecture | Single-user MVP with multi-user-ready data model. All entities linked via userId. One default user created at bootstrap. No registration, login, or user management in MVP. | Product owner decision (PRD patch 20260618) |
| DEC-008 | Q-SCOPE-004: Performance targets | Specific p95 targets for UI, API, AI export, and backup operations defined in product-brief.md §Performance Targets | Product owner decision (PRD patch 20260618) |
| DEC-009 | Q-SCOPE-005: Cardio data model | DailyLog replaces WorkoutDay. Cardio is separate entity with required dailyLogId. WorkoutExercise also references DailyLog. System auto-creates DailyLog on first activity for a date. | Product owner decision (PRD patch 20260618) |

## Assumptions Adopted

| ID | Assumption | Rationale | Source Evidence |
| --- | --- | --- | --- |
| DEC-005 | User has basic Docker/CLI knowledge | Self-hosted deployment requires technical proficiency | PRD §1, §4 |
| DEC-006 | Single user per instance | MVP explicitly single-user | PRD §4, §7.1 |
| DEC-007 | AI analysis done externally (manual copy-paste) | No API integration in MVP | PRD §17.1, §22, §23 |
| DEC-008 | Media on local filesystem via volume | Docker volume stated, no cloud storage | PRD §24.2 |
| DEC-009 | Full backup only (no incremental) | PRD specifies only full backup format | PRD §20 |

## Rejected Or Outdated Inputs

None in this run — single source document, no outdated inputs.