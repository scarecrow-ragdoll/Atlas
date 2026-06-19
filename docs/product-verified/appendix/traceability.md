# Traceability

## Requirement Map

| Requirement ID | Description | Source | Evidence | Type |
| --- | --- | --- | --- | --- |
| REQ-001 | PIN optional, hash-stored, session via cookie | §7.2 | Source: docs/product/prd.md | Source |
| REQ-002 | Dashboard with weekly summary and quick actions | §9 | Source: docs/product/prd.md | Source |
| REQ-003 | Exercise CRUD with name, muscle groups, working weight, media | §11 | Source: docs/product/prd.md | Source |
| REQ-004 | Workout diary by date with calendar navigation | §10 | Source: docs/product/prd.md | Source |
| REQ-005 | Sets with weight, reps, optional RPE/RIR | §10.5 | Source: docs/product/prd.md | Source |
| REQ-006 | Working weight auto-population and snapshot | §10.6 | Source: docs/product/prd.md | Source |
| REQ-007 | Cardio with type, duration, pulse/zone | §12 | Source: docs/product/prd.md | Source |
| REQ-008 | Weekly body check-in with measurements and photos | §13 | Source: docs/product/prd.md | Source |
| REQ-009 | Standalone body weight entry | §13.5 | Source: docs/product/prd.md | Source |
| REQ-010 | Nutrition product catalog with KJBJU per 100g | §15.2 | Source: docs/product/prd.md | Source |
| REQ-011 | Weekly nutrition template with daily overrides | §15.3-§15.5 | Source: docs/product/prd.md | Source |
| REQ-012 | Training, body, and nutrition progress charts | §16 | Source: docs/product/prd.md | Source |
| REQ-013 | AI export ZIP with manifest, data, summary, CSVs | §17 | Source: docs/product/prd.md | Source |
| REQ-014 | AI prompt builder with persistent context and week flags | §18 | Source: docs/product/prd.md | Source |
| REQ-015 | AI review history with manual entry | §19 | Source: docs/product/prd.md | Source |
| REQ-016 | Full backup export and import | §20 | Source: docs/product/prd.md | Source |
| REQ-017 | No data loss on container restart | §24.2 | Source: docs/product/prd.md | Source |
| REQ-018 | Media not publicly accessible without PIN session | §14.3, §24.1 | Source: docs/product/prd.md | Source |
| REQ-019 | Single-user mode (no registration, no multi-user) | §4, §7.1 | Source: docs/product/prd.md | Source |
| REQ-020 | Photo count 2-4 per check-in | §13.2 | Source: docs/product/prd.md | Source |
| REQ-021 | Paired measurements allow left/right or single value | §13.4 | Source: docs/product/prd.md | Source |
| REQ-022 | AI export default period: last 4 weeks | §17.2 | Source: docs/product/prd.md | Source |
| REQ-023 | Photos excluded from AI export by default | §17.3 | Source: docs/product/prd.md | Source |
| REQ-024 | Import dry-run validation before restore | §20.4 | Source: docs/product/prd.md | Source |
| REQ-025 | No silent partial import | §20.4 | Source: docs/product/prd.md | Source |
| REQ-026 | Schema version in backup manifest | §20.5 | Source: docs/product/prd.md | Source |
| AC-001 | User can enable and disable PIN | §7.2 | Source: docs/product/prd.md | Source |
| AC-002 | User can create exercise with name, muscle groups, working weight | §11.1, §11.2 | Source: docs/product/prd.md | Source |
| AC-003 | User can set and modify exercise working weight | §10.6, §11.2 | Source: docs/product/prd.md | Source |
| AC-004 | User can upload images and video to exercise | §11.3 | Source: docs/product/prd.md | Source |
| AC-005 | User can open current day in workout diary | §10.2 | Source: docs/product/prd.md | Source |
| AC-006 | User can select a past date via calendar | §10.2 | Source: docs/product/prd.md | Source |
| AC-007 | User can add exercise to workout day | §10.4 | Source: docs/product/prd.md | Source |
| AC-008 | User can add sets with weight and reps | §10.5 | Source: docs/product/prd.md | Source |
| AC-009 | User can optionally specify RPE per set | §10.5 | Source: docs/product/prd.md | Source |
| AC-010 | User can optionally specify RIR per set | §10.5 | Source: docs/product/prd.md | Source |
| AC-011 | User can add comment to exercise in workout | §10.4, §26.2 | Source: docs/product/prd.md | Source |
| AC-012 | User can add cardio with type and duration | §12.2, §26.4 | Source: docs/product/prd.md | Source |
| AC-013 | User can optionally specify pulse and HR zone for cardio | §12.2, §12.4 | Source: docs/product/prd.md | Source |
| AC-014 | User can create weekly body check-in | §13.2, §26.5 | Source: docs/product/prd.md | Source |
| AC-015 | Check-in includes date, weight, optional body fat %, measurements, photos, comment | §13.2 | Source: docs/product/prd.md | Source |
| AC-016 | User can enter body weight for any date | §13.5 | Source: docs/product/prd.md | Source |
| AC-017 | User can create nutrition product with name and KJBJU per 100g | §15.2 | Source: docs/product/prd.md | Source |
| AC-018 | User can create weekly nutrition template | §15.3 | Source: docs/product/prd.md | Source |
| AC-019 | User can override nutrition for specific day | §15.5 | Source: docs/product/prd.md | Source |
| AC-020 | User can view exercise progress chart | §16.2 | Source: docs/product/prd.md | Source |
| AC-021 | User can view body weight and measurement charts | §16.3 | Source: docs/product/prd.md | Source |
| AC-022 | User can view basic nutrition charts | §16.4 | Source: docs/product/prd.md | Source |
| AC-023 | User can generate AI prompt | §18 | Source: docs/product/prd.md | Source |
| AC-024 | User can download AI export ZIP for last 4 weeks | §17 | Source: docs/product/prd.md | Source |
| AC-025 | User can save AI review | §19 | Source: docs/product/prd.md | Source |
| AC-026 | User can download full backup ZIP | §20.1 | Source: docs/product/prd.md | Source |
| AC-027 | User can import backup into clean instance | §20.4 | Source: docs/product/prd.md | Source |
| AC-028 | User can run verification command with passing tests | §24.4 | Source: docs/product/prd.md | Source |
| AC-029 | PIN is optional | §7.2 | Source: docs/product/prd.md | Source |
| AC-030 | PIN enabled requires entry before access | §7.2 | Source: docs/product/prd.md | Source |
| AC-031 | PIN stored as hash | §7.2 | Source: docs/product/prd.md | Source |
| AC-032 | User can change PIN | §7.2 | Source: docs/product/prd.md | Source |
| AC-033 | User can disable PIN | §7.2 | Source: docs/product/prd.md | Source |
| AC-034 | PIN session persists via cookie | §7.2 | Source: docs/product/prd.md | Source |
| AC-035 | Diary opens to today by default | §10.2 | Source: docs/product/prd.md | Source |
| AC-036 | User can navigate via calendar | §10.2 | Source: docs/product/prd.md | Source |
| AC-037 | Existing date record opens existing | §10.2 | Source: docs/product/prd.md | Source |
| AC-038 | New date record created on first save | §10.2 | Source: docs/product/prd.md | Source |
| AC-039 | Working weight auto-populated on add exercise | §10.6 | Source: docs/product/prd.md | Source |
| AC-040 | User can add multiple sets per exercise | §10.5 | Source: docs/product/prd.md | Source |
| AC-041 | Working weight snapshot stored per exercise | §10.6 | Source: docs/product/prd.md | Source |
| AC-042 | Exercise comment included in AI export | §10.4 | Source: docs/product/prd.md | Source |
| AC-043 | User creates exercises manually | §11.1 | Source: docs/product/prd.md | Source |
| AC-044 | Exercise includes name, muscle groups, description, working weight | §11.2 | Source: docs/product/prd.md | Source |
| AC-045 | User can upload media after exercise creation | §11.3 | Source: docs/product/prd.md | Source |
| AC-046 | User can delete media from exercise | §11.3 | Source: docs/product/prd.md | Source |
| AC-047 | Exercise can be marked active or inactive | §11.2 | Source: docs/product/prd.md | Source |
| AC-048 | User can select cardio type from predefined list | §12.3 | Source: docs/product/prd.md | Source |
| AC-049 | Cardio duration recorded in minutes | §12.2 | Source: docs/product/prd.md | Source |
| AC-050 | User can optionally record avg pulse | §12.4 | Source: docs/product/prd.md | Source |
| AC-051 | User can optionally select HR zone | §12.4 | Source: docs/product/prd.md | Source |
| AC-052 | Check-in records date, weight, optional body fat | §13.2 | Source: docs/product/prd.md | Source |
| AC-053 | Check-in includes 10 measurement types | §13.3 | Source: docs/product/prd.md | Source |
| AC-054 | Paired measurements can record left, right, or single | §13.4 | Source: docs/product/prd.md | Source |
| AC-055 | Second value in paired measurement not required | §13.4 | Source: docs/product/prd.md | Source |
| AC-056 | 2-4 photos can be attached to check-in | §13.2 | Source: docs/product/prd.md | Source |
| AC-057 | Standalone weight entry for any date | §13.5 | Source: docs/product/prd.md | Source |
| AC-058 | User creates products with KJBJU per 100g | §15.2 | Source: docs/product/prd.md | Source |
| AC-059 | Template contains products with gram amounts | §15.3 | Source: docs/product/prd.md | Source |
| AC-060 | Template items can have optional meal label | §15.3 | Source: docs/product/prd.md | Source |
| AC-061 | Template auto-applies to all week days | §15.4 | Source: docs/product/prd.md | Source |
| AC-062 | Daily override can add/remove/replace products | §15.5 | Source: docs/product/prd.md | Source |
| AC-063 | KJBJU calculated from template items | §15.3 | Source: docs/product/prd.md | Source |
| AC-064 | KJBJU recalculated after override change | §15.5 | Source: docs/product/prd.md | Source |
| AC-065 | Exercise chart shows working weight over period | §16.2 | Source: docs/product/prd.md | Source |
| AC-066 | Exercise chart shows best set per session | §16.2 | Source: docs/product/prd.md | Source |
| AC-067 | Exercise chart shows e1RM over time | §16.2 | Source: docs/product/prd.md | Source |
| AC-068 | Exercise chart shows volume per session | §16.2 | Source: docs/product/prd.md | Source |
| AC-069 | Body weight chart shows weight over period | §16.3 | Source: docs/product/prd.md | Source |
| AC-070 | Measurement chart shows individual measurement | §16.3 | Source: docs/product/prd.md | Source |
| AC-071 | User can overlay multiple measurements | §16.3 | Source: docs/product/prd.md | Source |
| AC-072 | Nutrition chart shows weekly avg KJBJU | §16.4 | Source: docs/product/prd.md | Source |
| AC-073 | User can select chart period | §16.1 | Source: docs/product/prd.md | Source |
| AC-074 | Default AI export range is last 4 weeks | §17.2 | Source: docs/product/prd.md | Source |
| AC-075 | User can select custom date range | §17.2 | Source: docs/product/prd.md | Source |
| AC-076 | User can toggle data sections | §17.3 | Source: docs/product/prd.md | Source |
| AC-077 | Photos excluded from export by default | §17.3 | Source: docs/product/prd.md | Source |
| AC-078 | ZIP contains manifest.json, data.json, summary.md | §17.4 | Source: docs/product/prd.md | Source |
| AC-079 | ZIP contains CSV files | §17.8 | Source: docs/product/prd.md | Source |
| AC-080 | ZIP contains photos/ directory when included | §17.4 | Source: docs/product/prd.md | Source |
| AC-081 | manifest.json includes export metadata | §17.5 | Source: docs/product/prd.md | Source |
| AC-082 | data.json includes all selected sections | §17.6 | Source: docs/product/prd.md | Source |
| AC-083 | summary.md includes period, stats, trends | §17.7 | Source: docs/product/prd.md | Source |
| AC-084 | Persistent AI context stored and reused | §18.2 | Source: docs/product/prd.md | Source |
| AC-085 | Context includes goal, height, optional fields | §18.2 | Source: docs/product/prd.md | Source |
| AC-086 | User can update context any time | §18.2 | Source: docs/product/prd.md | Source |
| AC-087 | User can add one-time comment | §18.3 | Source: docs/product/prd.md | Source |
| AC-088 | User can select week flags | §18.4 | Source: docs/product/prd.md | Source |
| AC-089 | Prompt asks AI for analysis and recommendations | §18.5 | Source: docs/product/prd.md | Source |
| AC-090 | User can paste AI response text | §19.2 | Source: docs/product/prd.md | Source |
| AC-091 | User can link review to date range | §19.2 | Source: docs/product/prd.md | Source |
| AC-092 | User can add notes and planned actions | §19.2 | Source: docs/product/prd.md | Source |
| AC-093 | Backup ZIP contains manifest.json, data.json, media | §20.1 | Source: docs/product/prd.md | Source |
| AC-094 | manifest.json includes type, schema version, date | §20.2 | Source: docs/product/prd.md | Source |
| AC-095 | data.json includes all entities | §20.3 | Source: docs/product/prd.md | Source |
| AC-096 | User can include/exclude media | §20.1 | Source: docs/product/prd.md | Source |
| AC-097 | Import validates manifest.json | §20.4 | Source: docs/product/prd.md | Source |
| AC-098 | Import validates schema version | §20.5 | Source: docs/product/prd.md | Source |
| AC-099 | Import runs dry-run validation | §20.4 | Source: docs/product/prd.md | Source |
| AC-100 | Import shows summary before confirm | §20.4 | Source: docs/product/prd.md | Source |
| AC-101 | Import restores fully or fails | §20.4 | Source: docs/product/prd.md | Source |
| AC-102 | Import shows clear error messages | §20.4 | Source: docs/product/prd.md | Source |
| AC-103 | Dashboard shows current date | §9 | Source: docs/product/prd.md | Source |
| AC-104 | Dashboard shows last body weight | §9 | Source: docs/product/prd.md | Source |
| AC-105 | Dashboard shows training days count | §9 | Source: docs/product/prd.md | Source |
| AC-106 | Dashboard shows cardio sessions count | §9 | Source: docs/product/prd.md | Source |
| AC-107 | Dashboard shows current goal | §9 | Source: docs/product/prd.md | Source |
| AC-108 | Dashboard shows quick actions | §9 | Source: docs/product/prd.md | Source |
| AC-109 | PIN change fails without correct current PIN | §7.2 | Source: docs/product/prd.md | Source |
| AC-110 | PIN pages inaccessible without valid session | §7.2 | Source: docs/product/prd.md | Source |
| AC-111 | Media returns 401/403 without session | §14.3, §24.1 | Source: docs/product/prd.md | Source |
| AC-112 | Photos not in AI export unless opted in | §17.3 | Source: docs/product/prd.md | Source |
| AC-113 | Daily override does not affect other dates | §15.5 | Source: docs/product/prd.md | Source |
| AC-114 | Backup fails cleanly with invalid manifest | §20.4 | Source: docs/product/prd.md | Source |
| AC-115 | Backup fails cleanly with incompatible schema | §20.5 | Source: docs/product/prd.md | Source |
| AC-116 | No silent partial import | §20.4 | Source: docs/product/prd.md | Source |
| AC-117 | PIN not logged | §24.1 | Source: docs/product/prd.md | Source |
| AC-118 | AI export content not logged | §24.1 | Source: docs/product/prd.md | Source |
| AC-119 | Photos not logged | §24.1 | Source: docs/product/prd.md | Source |
| AC-120 | Sensitive comments not logged | §24.1 | Source: docs/product/prd.md | Source |
| AC-121 | Full round-trip: exercise -> workout -> cardio -> check-in -> export | §29 | Source: docs/product/prd.md | Source |
| AC-122 | PIN round-trip: enable -> close -> reopen -> enter PIN | §29 | Source: docs/product/prd.md | Source |
| AC-123 | Nutrition: template -> override -> verify different values | §29 | Source: docs/product/prd.md | Source |
| AC-124 | Backup: create -> reset -> import -> verify restored | §29 | Source: docs/product/prd.md | Source |
| AC-125 | Test suite passes with coverage | §24.4, §29 | Source: docs/product/prd.md | Source |

| EDGE-001 | Set with 0 weight or 0 reps | §10.5 | Source: docs/product/prd.md | Derivation |
| EDGE-002 | Exercise name duplicate | §11 | Source: docs/product/prd.md | Derivation |
| EDGE-003 | Nutrition product 0 or negative values | §15.2 | Source: docs/product/prd.md | Derivation |
| EDGE-004 | Backdated workout before instance creation | §10.2 | Source: docs/product/prd.md | Derivation |
| EDGE-005 | Empty workout day | §10.2 | Source: docs/product/prd.md | Derivation |
| EDGE-006 | Check-in with 0-1 photos | §13.2 | Source: docs/product/prd.md | Derivation |
| EDGE-007 | Body measurement 0 or negative | §13.3 | Source: docs/product/prd.md | Derivation |
| EDGE-008 | AI export with no data in period | §17 | Source: docs/product/prd.md | Derivation |
| EDGE-009 | Nutrition template with zero items | §15.3 | Source: docs/product/prd.md | Derivation |
| EDGE-010 | Invalid backup ZIP uploaded | §20.4 | Source: docs/product/prd.md | Derivation |
| EDGE-011 | PIN enabled but pinHash missing | §7.2 | Source: docs/product/prd.md | Derivation |
| EDGE-012 | PIN session expired during entry | §7.2 | Source: docs/product/prd.md | Derivation |
| EDGE-013 | PIN disabled no access control | §7.2 | Source: docs/product/prd.md | Derivation |
| EDGE-014 | Media URL accessed directly | §14.3, §24.1 | Source: docs/product/prd.md | Derivation |
| EDGE-015 | Browser tab stale session | §7.2 | Source: docs/product/prd.md | Derivation |
| EDGE-016 | Same day opened in two tabs | §10.2 | Source: docs/product/prd.md | Derivation |
| EDGE-017 | Template created mid-week | §15.4 | Source: docs/product/prd.md | Derivation |
| EDGE-018 | Exercise with historical data deleted | §11 | Source: docs/product/prd.md | Derivation |
| EDGE-019 | Nutrition product with active template deleted | §15.2 | Source: docs/product/prd.md | Derivation |
| EDGE-020 | Media file deleted but DB record remains | §11.3 | Source: docs/product/prd.md | Derivation |
| EDGE-021 | Import partially restored | §20.4 | Source: docs/product/prd.md | Derivation |
| EDGE-022 | PostgreSQL connection lost | §24.2 | Source: docs/product/prd.md | Derivation |
| EDGE-023 | Redis unavailable | §7.2 | Source: docs/product/prd.md | Derivation |
| EDGE-024 | Disk full during export | §17, §20 | Source: docs/product/prd.md | Derivation |
| EDGE-025 | Docker volume full | §24.2 | Source: docs/product/prd.md | Derivation |
| EDGE-026 | System clock changes | All date features | Source: docs/product/prd.md | Derivation |
| EDGE-027 | No data retention policy | General | Source: docs/product/prd.md | Derivation |
| EDGE-028 | Schema migration after backup | §20.5 | Source: docs/product/prd.md | Derivation |
| EDGE-029 | Years of training degrade performance | §24.3 | Source: docs/product/prd.md | Derivation |
| EDGE-030 | Media storage unbounded | §14 | Source: docs/product/prd.md | Derivation |
| EDGE-031 | Timezone handling for date features | All date entities | Source: docs/product/prd.md | Derivation |
| RULE-005 | Photo count 2-4 per check-in | §13.2 | Source: docs/product/prd.md | Source |
| RULE-011 | Daily nutrition recalculated on override | §15.5 | Source: docs/product/prd.md | Source |
| RULE-015 | Nutrition weekly averages | §16.4 | Source: docs/product/prd.md | Source |
| RULE-016 | Workout day created on first save | §10.2 | Source: docs/product/prd.md | Source |
| RULE-021 | AI export default 4 weeks | §17.2 | Source: docs/product/prd.md | Source |
| RULE-029 | No automatic external API calls in MVP | §22, §23 | Source: docs/product/prd.md | Source |

## Source Map

| Source File | Requirements | Evidence |
| --- | --- | --- |
| docs/product/prd.md §4 | REQ-019 | Source: docs/product/prd.md |
| docs/product/prd.md §7.1 | REQ-019 | Source: docs/product/prd.md |
| docs/product/prd.md §7.2 | REQ-001, AC-029, AC-030, AC-031, AC-032, AC-033, AC-034, RULE-001, RULE-002, RULE-022, RULE-023 | Source: docs/product/prd.md |
| docs/product/prd.md §9 | REQ-002, AC-103, AC-104, AC-105, AC-106, AC-107, AC-108 | Source: docs/product/prd.md |
| docs/product/prd.md §10.2 | REQ-004, AC-035, AC-036, AC-037, AC-038 | Source: docs/product/prd.md |
| docs/product/prd.md §10.4 | REQ-003, AC-011 | Source: docs/product/prd.md |
| docs/product/prd.md §10.5 | REQ-005, AC-008, AC-009, AC-010, RULE-004 | Source: docs/product/prd.md |
| docs/product/prd.md §10.6 | REQ-006, AC-039, AC-040, AC-041, RULE-003, RULE-017 | Source: docs/product/prd.md |
| docs/product/prd.md §10.7 | RULE-012, RULE-013, RULE-014 | Source: docs/product/prd.md |
| docs/product/prd.md §11 | REQ-003, AC-043, AC-044, AC-045, AC-046, AC-047 | Source: docs/product/prd.md |
| docs/product/prd.md §11.2 | AC-044 | Source: docs/product/prd.md |
| docs/product/prd.md §11.3 | AC-045, AC-046 | Source: docs/product/prd.md |
| docs/product/prd.md §12 | REQ-007, AC-048, AC-049, AC-050, AC-051 | Source: docs/product/prd.md |
| docs/product/prd.md §13.2 | REQ-008, REQ-020 | Source: docs/product/prd.md |
| docs/product/prd.md §13.3 | measurement list | Source: docs/product/prd.md |
| docs/product/prd.md §13.4 | REQ-021, RULE-009 | Source: docs/product/prd.md |
| docs/product/prd.md §13.5 | REQ-009, AC-057 | Source: docs/product/prd.md |
| docs/product/prd.md §14 | REQ-018 | Source: docs/product/prd.md |
| docs/product/prd.md §15.2 | REQ-010, RULE-006 | Source: docs/product/prd.md |
| docs/product/prd.md §15.3-§15.5 | REQ-011, RULE-010, RULE-018, RULE-019, RULE-020 | Source: docs/product/prd.md |
| docs/product/prd.md §16 | REQ-012 | Source: docs/product/prd.md |
| docs/product/prd.md §17 | REQ-013, REQ-022, REQ-023, AC-074, AC-075, AC-076, AC-077, AC-078, AC-079, AC-080, AC-081, AC-082, AC-083 | Source: docs/product/prd.md |
| docs/product/prd.md §18 | REQ-014 | Source: docs/product/prd.md |
| docs/product/prd.md §19 | REQ-015 | Source: docs/product/prd.md |
| docs/product/prd.md §20 | REQ-016, REQ-024, REQ-025, REQ-026, RULE-007, RULE-008, RULE-009 | Source: docs/product/prd.md |
| docs/product/prd.md §24.1 | REQ-018, RULE-024, RULE-025, RULE-026, RULE-027, RULE-028 | Source: docs/product/prd.md |
| docs/product/prd.md §24.2 | REQ-017 | Source: docs/product/prd.md |
| docs/product/prd.md §25 | Domain model entities and attributes | Source: docs/product/prd.md |

## Assumption Map

| Assumption | Decision ID | Source Rationale |
| --- | --- | --- |
| User has basic Docker deployment knowledge | DEC-005 | Self-hosted requirement without setup flow specification |
| Single-user data model with userId on all entities | DEC-007 | Product owner decision resolving Q-SCOPE-002 |
| AI analysis performed externally (ChatGPT or similar) | DEC-003 | PRD §17 structure without API integration |
| Media files on local filesystem volume | DEC-004 | PRD §24.2 requires volume storage, no cloud option |
| Backup/restore covers full lifecycle only | DEC-005 | PRD §20 specifies only full backup, no incremental |
| Cardio belongs to DailyLog (required dailyLogId) | DEC-009 | Product owner decision resolving Q-SCOPE-005 |
| Performance targets from DEC-008 | DEC-008 | Product owner decision resolving Q-SCOPE-004 |
| Success metrics from DEC-006 | DEC-006 | Product owner decision resolving Q-SCOPE-001 |

## Open Question Map

| Question ID | Impacted Requirements | Impacted Acceptance Criteria | Evidence |
| --- | --- | --- | --- |
| Q-SCOPE-001 | REQ-001 through REQ-026 | AC-001 through AC-125 | Source: docs/product/prd.md |
| Q-SCOPE-002 | REQ-019 | AC-001 through AC-125 | Source: docs/product/prd.md |
| Q-SCOPE-004 | REQ-002, REQ-012, REQ-013 | AC-103, AC-104, AC-105, AC-106, AC-107, AC-108, AC-065, AC-066, AC-067, AC-068, AC-069, AC-070, AC-071, AC-072, AC-073, AC-074, AC-075, AC-076, AC-077, AC-078, AC-079, AC-080, AC-081, AC-082, AC-083 | Source: docs/product/prd.md |
| Q-SCOPE-005 | REQ-007 | AC-048, AC-049, AC-050, AC-051, AC-012, AC-013 | Source: docs/product/prd.md |
| Q-ROLE-001 | REQ-001 | AC-029, AC-030, AC-031, AC-032, AC-033, AC-034, AC-109 | Source: docs/product/prd.md |
| Q-ROLE-002 | REQ-001 | AC-029, AC-030, AC-031, AC-032, AC-033, AC-034 | Source: docs/product/prd.md |
| Q-ROLE-003 | REQ-001 | AC-029, AC-030, AC-031, AC-032, AC-033, AC-034 | Source: docs/product/prd.md |
| Q-API-001 | REQ-003 through REQ-016 | AC-001 through AC-125 | Source: docs/product/prd.md |
| Q-COMP-001 | REQ-017 | EDGE-027, EDGE-028, EDGE-029, EDGE-030, EDGE-031 | Source: docs/product/prd.md |
| Q-AC-01 | REQ-001 | AC-029, AC-030, AC-031, AC-032, AC-033, AC-034 | Source: docs/product/prd.md |
| Q-AC-07 | REQ-006 | AC-041, AC-136 | Source: docs/product/prd.md |
| Q-AC-08 | REQ-006 | AC-039 | Source: docs/product/prd.md |
| Q-CONS-002 | REQ-004, REQ-007, REQ-008, REQ-009, REQ-013 | AC-035, AC-036, AC-037, AC-038, AC-048, AC-049, AC-050, AC-051, AC-052, AC-053, AC-054, AC-055, AC-056, AC-057, AC-074, AC-075, AC-076, AC-077, AC-078, AC-079, AC-080, AC-081, AC-082, AC-083 | Source: docs/product/prd.md |