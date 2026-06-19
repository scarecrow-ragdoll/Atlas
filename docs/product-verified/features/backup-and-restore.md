# Backup And Restore

## Source Evidence

PRD §20, §26.11, §26.12.

## User Problem

Guarantee complete data ownership through full export and import, enabling migration, recovery, and peace of mind.

## Scope

In MVP. Full backup only (no incremental, no partial, no cloud).

## Behavior

- Full backup ZIP: manifest.json, data.json, media/
- Optional CSV files for manual inspection
- Schema version in manifest for forward compatibility
- Import: upload ZIP → validate manifest → validate schema → dry-run → show summary → user confirms → full restore
- No silent partial import — complete success or clean failure

## Derived Requirements

| Requirement | Source | Rationale | Confidence |
| --- | --- | --- | --- |
| Dry-run validation before import | §20.4 | "dry-run validation перед импортом" | High |
| Import summary before confirmation | §20.4 | "отображение summary перед импортом" | High |

## Edge Cases

EDGE-007 (invalid ZIP), EDGE-015 (import to instance with existing data), EDGE-021 (mid-import failure), EDGE-028 (schema migration).

## Acceptance Criteria

AC-026, AC-027, AC-093 through AC-102, AC-114 through AC-116.

## Dependencies

All data entities and media storage.

## Open Questions

Q-ACTOR-08: Import when data already exists (merge/replace/error).
Q-AC-15: Import with existing data behavior.
Q-AC-16: CSV files — mandatory or optional.
Q-EDGE-11: Migration strategy for schema version changes.