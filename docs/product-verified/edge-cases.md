# Edge Cases

## Input And Validation

| ID | Edge Case | Domain | Source Signal |
| --- | --- | --- | --- |
| EDGE-001 | Set with 0 weight or 0 reps entered | Workout | §10.5: set has weight and reps, no validation defined |
| EDGE-002 | Exercise name duplicate | Exercise Library | §11: user creates exercises, no duplicate rule |
| EDGE-003 | Nutrition product 0 or negative nutritional values | Nutrition | §15.2: values per 100g, no validation |
| EDGE-004 | Backdated workout before instance creation date | Workout | §10.2: backdating allowed, no lower bound |
| EDGE-005 | Empty workout day (no sets, no cardio, no notes) | Workout | §10.2: created on first save — empty state undefined |
| EDGE-006 | Check-in with 0-1 photos | Body Tracking | §13.2: "2-4 photos" — requirement vs recommendation |
| EDGE-007 | Body measurement value 0 or negative | Body Tracking | §13.3: measurement values without validation rule |
| EDGE-008 | AI export date range with no data in period | AI Export | §17: period selector, empty data not addressed |
| EDGE-009 | Nutrition template with zero items | Nutrition | §15.3: template creation, empty template not addressed |
| EDGE-010 | Invalid backup ZIP uploaded (corrupted, wrong format) | Backup | §20.4: manifest + schema check, but no explicit error format |

## Permissions And Ownership

| ID | Edge Case | Domain | Source Signal |
| --- | --- | --- | --- |
| EDGE-011 | PIN enabled but pinHash missing or corrupted | Access Control | §7.2: hash optional in Settings model |
| EDGE-012 | PIN session expired during active data entry | Access Control | §7.2: session via cookie, TTL unspecified |
| EDGE-013 | PIN disabled — no access control at all | Access Control | §7.2: no auth when PIN off |
| EDGE-014 | Media URL accessed directly without valid session | Privacy | §14.3, §24.1: media requires session |
| EDGE-015 | Browser tab left open with stale session | Access Control | §7.2: session management undefined |

## State And Concurrency

| ID | Edge Case | Domain | Source Signal |
| --- | --- | --- | --- |
| EDGE-016 | Same day opened in two browser tabs — last write wins vs conflict | Workout | §10.2: single user, but multi-tab is possible |
| EDGE-017 | Nutrition template created mid-week — retroactive or forward only | Nutrition | §15.4: template applied to week, week start behavior undefined |
| EDGE-018 | Exercise with historical workout data deleted | Exercise | §11: deletion allowed, referential integrity undefined |
| EDGE-019 | Nutrition product with active template items deleted | Nutrition | §15.2: deletion behavior for referenced product |
| EDGE-020 | Media file deleted from filesystem but DB record remains | Media | §11.3: media deletion, but orphaned files risk |
| EDGE-021 | Import ZIP partially restored — transaction rollback undefined | Backup | §20.4: full restore required, but failure handling mid-import |

## External Dependencies

| ID | Edge Case | Domain | Source Signal |
| --- | --- | --- | --- |
| EDGE-022 | PostgreSQL connection lost during save | Data Persistence | §24.2: data not lost on restart, but mid-session failure |
| EDGE-023 | Redis unavailable for session store | Session | §7.2: PIN session, Redis in stack |
| EDGE-024 | Disk full during export ZIP generation | Export | §17, §20: export generates files, no quota handling |
| EDGE-025 | Docker volume full — media save fails | Media | §24.2: media in volume, no capacity management |
| EDGE-026 | System clock changes affect date-based queries | General | All date features: today, week, charts, exports |

## Data Lifecycle

| ID | Edge Case | Domain | Source Signal |
| --- | --- | --- | --- |
| EDGE-027 | No data retention or deletion policy | General | MVP: no audit, no retention rules |
| EDGE-028 | Schema migration after backup import from older version | Backup | §20.5: schema version, migration strategy undefined |
| EDGE-029 | Years of training data degrade performance | Performance | §24.3: personal volume, no specific degradation target |
| EDGE-030 | Media accumulated over years — storage growth unbounded | Media | §14: photos accumulate, no storage limit policy |
| EDGE-031 | Timezone handling for date-based features (dashboard week, training day, check-in) | General | All date entities, timezone undefined (Q-CONS-002) |