<!-- FILE: docs/prd-wave-details/appendix/traceability.md -->
<!-- VERSION: 1.0.1 -->

# Traceability

## Slice Map

| Slice ID | Source Wave Capability | Description |
|----------|----------------------|-------------|
| SLICE-W09-001 | CAP-W09-001 through CAP-W09-008 | Migration 00094 + BackupConfig |
| SLICE-W09-002 | CAP-W09-001, CAP-W09-003 | Backup export data aggregation service |
| SLICE-W09-003 | CAP-W09-001, CAP-W09-004 | Backup ZIP generation (BackupArchive + BackupManifest + BackupData) |
| SLICE-W09-004 | CAP-W09-001, CAP-W09-002 | Export REST handler |
| SLICE-W09-005 | CAP-W09-005, CAP-W09-006 | Import validation service (dry-run) |
| SLICE-W09-006 | CAP-W09-006, CAP-W09-007 | Import restore service (all-or-nothing transaction) |
| SLICE-W09-007 | CAP-W09-005, CAP-W09-006 | Import REST handler |
| SLICE-W09-008 | CAP-W09-001 through CAP-W09-008 | GraphQL backup metadata schema + resolvers |
| SLICE-W09-009 | CAP-W09-001 through CAP-W09-008 | Main wiring (resolver.go, main.go, atlas-gqlgen.yml) |

## Acceptance Criteria Map

| AC ID | Source | Description |
|-------|--------|-------------|
| AC-W09-001 | AC-093, AC-026 | Full backup ZIP contains manifest.json, data.json, media/ |
| AC-W09-002 | AC-094 | manifest.json includes type, schema version, app version, date, sections |
| AC-W09-003 | AC-095 | data.json includes all entities |
| AC-W09-004 | AC-096 | User can include or exclude media from backup |
| AC-W09-005 | AC-097 | Import validates manifest.json structure |
| AC-W09-006 | AC-098, AC-115 | Import validates schema version |
| AC-W09-007 | AC-099 | Import runs dry-run validation before actual restore |
| AC-W09-008 | AC-100 | Import shows summary before user confirmation |
| AC-W09-009 | AC-101, AC-116 | Import restores fully or fails (all-or-nothing) |
| AC-W09-010 | AC-102, AC-114 | Import displays clear error messages on validation failure |
| AC-W09-011 | AC-124 | Backup create → reset → import → verify data restored |
| AC-W09-012 | AC-026, RULE-028 | Backup export generated only on user request |
| AC-W09-013 | RULE-008 | Dry-run validation before import |

## Exit Criteria Map

| EC ID | Validation Type | Source |
|-------|----------------|--------|
| EC-W09-001 | Migration | SLICE-W09-001 |
| EC-W09-002 | Build | All slices |
| EC-W09-003 | Service unit tests | SLICE-W09-002, SLICE-W09-003, SLICE-W09-005, SLICE-W09-006 |
| EC-W09-004 | Handler unit tests | SLICE-W09-004, SLICE-W09-007 |
| EC-W09-005 | Repository integration tests | SLICE-W09-006 |
| EC-W09-006 | Traceability | All slices |
| EC-W09-007 | Transaction rollback proof | SLICE-W09-006 |
| EC-W09-008 | ZIP content verification | SLICE-W09-003 |

## Verification Obligation Map

| TEST ID | Layer | Description |
|---------|-------|-------------|
| TEST-W09-001 | Service | BackupService.Generate — all entity services called, ZIP built |
| TEST-W09-002 | Service | BackupService.Generate — media toggle respected |
| TEST-W09-003 | Service | BackupService.Generate — error propagation |
| TEST-W09-004 | Service | ImportService.Validate — manifest structure |
| TEST-W09-005 | Service | ImportService.Validate — schema version check |
| TEST-W09-006 | Service | ImportService.Validate — invalid ZIP rejection |
| TEST-W09-007 | Service | ImportService.Confirm — entity insertion in correct order |
| TEST-W09-008 | Service | ImportService.Confirm — transaction rollback |
| TEST-W09-009 | Service | ImportService.Confirm — media file restore |
| TEST-W09-010 | Integration | Export → import round-trip, data matches |
| TEST-W09-011 | Integration | Media toggle correct |
| TEST-W09-012 | Unit | BackupArchive manifest + data format |
| TEST-W09-013 | Handler | Handler.GenerateExport — valid request response |
| TEST-W09-014 | Handler | Handler.DownloadExport — file stream |
| TEST-W09-015 | Handler | Handler.ImportValidate — summary returned |
| TEST-W09-016 | Handler | Handler.ImportConfirm — restore executed |
| TEST-W09-017 | Handler | Handler error cases |
| TEST-W09-018 | Handler | Privacy: backup logs don't contain entity content |
| TEST-W09-019 | Benchmark | db-only export < 15s p95 |

## Code Touchpoint Map

| File | Slice | Operation |
|------|-------|-----------|
| internal/repository/postgres/migrations/00094_backup_schema_version.sql | SLICE-W09-001 | New file |
| internal/atlas/service/backup_export.go | SLICE-W09-002, SLICE-W09-003 | New file |
| internal/atlas/models/backup.go | SLICE-W09-002, SLICE-W09-003 | New file |
| internal/handler/backup_handler.go | SLICE-W09-004, SLICE-W09-007 | New file |
| internal/atlas/service/backup_import.go | SLICE-W09-005, SLICE-W09-006 | New file |
| internal/atlas/graph/schema/backup.graphql | SLICE-W09-008 | New file |
| internal/atlas/graph/resolver/backup.go | SLICE-W09-008 | New file |
| internal/atlas/graph/resolver/resolver.go | SLICE-W09-009 | Add BackupService field |
| cmd/server/main.go | SLICE-W09-009 | Wire repo→service→handler→routes |
| atlas-gqlgen.yml | SLICE-W09-009 | Add Backup type bindings |
| internal/appconfig/config.go | SLICE-W09-001 | Add BackupConfig section |

## Question Map

| Question ID | Source Report | Status |
|-------------|-------------|--------|
| DQ-W09-001 | planner/product-ac, product-verified/features/backup-and-restore.md | open (BLOCKING) |
| DQ-W09-005 | planner/architecture-codebase, product-verified/edge-cases.md | open (BLOCKING) |
| DQ-W09-002 | planner/product-ac, product-verified/functional-spec.md §20 | open (WATCHLIST) |
| DQ-W09-003 | planner/data-integration-ops | open (WATCHLIST) |
| DQ-W09-004 | planner/security-compliance | open (WATCHLIST) |

## Source Map

| Source Document | Relevant Artifacts |
|----------------|-------------------|
| docs/prd-waves/waves/wave-09.md | All slices, AC-W09-001 through AC-W09-013 |
| docs/product-verified/acceptance-criteria.md | AC-W09-001 through AC-W09-011 (AC-093-102, AC-114-116, AC-124) |
| docs/product-verified/features/backup-and-restore.md | All design contracts |
| docs/product-verified/functional-spec.md §20 | AC-W09-001 through AC-W09-010 |
| docs/product-verified/user-flows.md §26.11-§26.12 | SLICE-W09-004, SLICE-W09-007 |
| docs/product-verified/edge-cases.md | DQ-W09-005 (EDGE-028), DQ-W09-001 (EDGE-010, EDGE-021) |
| docs/product-verified/business-rules.md | AC-W09-012 (RULE-028), AC-W09-013 (RULE-008) |
| docs/product-verified/product-brief.md | Performance targets (db-only export < 15s, import < 30s) |
| docs/prd-wave-details/waves/wave-07.md | ZIP export pattern, REST handler pattern |
| docs/prd-wave-details/waves/wave-08.md | ListAllByUserID pattern |
| docs/knowledge-graph.xml | Codebase fit module references |
| docs/verification-plan.xml | Verification entry patterns |