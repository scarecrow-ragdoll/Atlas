# Wave 09: Backup Import/Export

## Status

user-approved

## User Approval

user-approved (2026-06-18)

## Purpose

Full backup and restore capability with versioning and media support.

## Outcome After Wave

- OUT-W09-001 Full backup ZIP export
- OUT-W09-002 Import with dry-run validation
- OUT-W09-003 Data version checking
- OUT-W09-004 Media restore
- OUT-W09-005 All-or-nothing transaction

## Included Scope

- CAP-W09-001 Backup ZIP generation
- CAP-W09-002 manifest.json with data version
- CAP-W09-003 data.json with all entities
- CAP-W09-004 media/ folder in backup
- CAP-W09-005 Import dry-run validation
- CAP-W09-006 Entity restoration
- CAP-W09-007 Relationship restoration
- CAP-W09-008 Migration framework

## Excluded Scope

- Cloud backup
- Incremental backup

## Dependencies

WAVE-01 through WAVE-08

## Surface Categories

backend, data, operations, security

## Risk Class

High - Data loss prevention

## Recommended Next Planning

$detail-prd-wave for WAVE-09

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |

## Traceability

- docs/product/prd.md Section 20
- docs/product-verified/acceptance-criteria.md#Backup