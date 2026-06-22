# Planner Report: Product & Acceptance Criteria (WAVE-09)

## Scope
Analyze WAVE-09 against product-verified docs: outcome, included/excluded scope, acceptance criteria mapping, product edge cases.

## Source Wave Summary
WAVE-09 implements full backup/restore: ZIP export with manifest.json, data.json, media/; multi-step import with dry-run validate → confirm → all-or-nothing restore; media inclusion toggle; migration framework for schema version differences.

## Outcome Mapping

| Source Outcome | WAVE-09 ID | Verification |
| --- | --- | --- |
| OUT-W09-001 Full backup ZIP export | OUT-W09-001 | POST /api/backup/export generates ZIP |
| OUT-W09-002 Import with dry-run validation | OUT-W09-002 | POST /api/backup/import/validate + /confirm |
| OUT-W09-003 Data version checking | OUT-W09-003 | manifest.json schema version checked on import |
| OUT-W09-004 Media restore | OUT-W09-004 | media/ folder in ZIP, optionally restored |
| OUT-W09-005 All-or-nothing transaction | OUT-W09-005 | Import wrapped in single DB transaction |

## Acceptance Criteria Map

| AC ID | WAVE-09 AC | Source AC | Verification |
| --- | --- | --- | --- |
| AC-W09-001 | Full backup ZIP contains manifest.json, data.json, media/ | AC-093, AC-026 | BuildZIP tests with all entity types |
| AC-W09-002 | manifest.json includes type, schema version, app version, date, sections | AC-094 | manifest structure test |
| AC-W09-003 | data.json includes all entities (settings, profile, exercises, workouts, cardio, body, nutrition, AI) | AC-095 | Entity aggregation test — verify each service called |
| AC-W09-004 | User can include or exclude media from backup | AC-096 | Toggle test, verify media/ present or absent |
| AC-W09-005 | Import validates manifest.json structure | AC-097 | Invalid manifest test |
| AC-W09-006 | Import validates schema version | AC-098, AC-115 | Schema mismatch test |
| AC-W09-007 | Import runs dry-run validation before actual restore | AC-099 | Validate step returns summary without side effects |
| AC-W09-008 | Import shows summary before user confirmation | AC-100 | Validate response includes summary fields |
| AC-W09-009 | Import restores data and media fully, or fails without partial import | AC-101, AC-116 | Transaction rollback test |
| AC-W09-010 | Import displays clear error messages on validation failure | AC-102, AC-114 | Error response codes and messages |
| AC-W09-011 | User can create backup → reset app → import backup → verify data restored | AC-124 | Round-trip e2e test |
| AC-W09-012 | Backup export generated only on user request | AC-026, RULE-028 | No auto-export |
| AC-W09-013 | Dry-run validation before import | RULE-008 | Validate step before confirm |

## Product Edge Cases Coverage

| Edge Case | WAVE-09 Coverage | Status |
| --- | --- | --- |
| EDGE-010: Invalid backup ZIP uploaded | AC-W09-005 validation check | Covered |
| EDGE-021: Import ZIP partially restored — rollback | AC-W09-009 all-or-nothing | Covered |
| EDGE-028: Schema migration after backup import from older version | AC-W09-006 schema version check | Covered |
| EDGE-024: Disk full during ZIP generation | Must be handled with error return | Needs explicit test |
| EDGE-015: Import to instance with existing data | DQ-W09-001 — behavior TBD | Open question |

## Gaps & Risks
1. **Existing data behavior not specified** — Q-ACTOR-08, Q-AC-15 need owner decision
2. **CSV files** mentioned in functional spec but not in AC-093 — recommended to exclude for MVP
3. **Performance targets** from product-brief: db-only <= 15s p95, with media best-effort