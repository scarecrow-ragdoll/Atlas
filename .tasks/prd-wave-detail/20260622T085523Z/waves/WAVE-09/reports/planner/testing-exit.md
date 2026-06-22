# Planner Report: Testing & Exit Criteria (WAVE-09)

## Scope
Exit criteria, backend verification obligations, focused commands for WAVE-09.

## Exit Criteria

| ID | Criterion | Verification Method |
| --- | --- | --- |
| EC-W09-001 | Migration 00094 applies cleanly, rollback succeeds | `bun run migrate` + rollback test |
| EC-W09-002 | SQLc codegen succeeds | `bun run codegen` |
| EC-W09-003 | Go build succeeds | `bun run build` |
| EC-W09-004 | GraphQL codegen succeeds | `bun run codegen` |
| EC-W09-005 | Backup service unit tests pass | `bun run test -- --run TestBackupService` |
| EC-W09-006 | Backup REST handler tests pass | `bun run test -- --run TestBackupHandler` |
| EC-W09-007 | Import service unit tests pass | `bun run test -- --run TestBackupImportService` |
| EC-W09-008 | Repository integration tests pass | `bun run test -- --run TestBackupRepository` |
| EC-W09-009 | All ACs covered by tests | Trace AC-W09-001 through AC-W09-013 in test files |
| EC-W09-010 | Open questions resolved or deferred with documented rationale | question-ledger.md shows resolution |
| EC-W09-011 | No partial import — transaction rollback proven | Test verifies rollback on any step failure |
| EC-W09-012 | ZIP content verified (manifest.json, data.json, media/) | BuildZIP output verification test |

## Verification Obligations

| TEST ID | Type | Scope | What It Proves |
| --- | --- | --- | --- |
| TEST-W09-001 | Unit | BackupService.Generate | All entity services called, ZIP built |
| TEST-W09-002 | Unit | BackupService.Generate | Media toggle respected |
| TEST-W09-003 | Unit | BackupService.Generate | Error from entity service propagated |
| TEST-W09-004 | Unit | BackupImportService.Validate | Manifest structure parsed correctly |
| TEST-W09-005 | Unit | BackupImportService.Validate | Schema version check passes/fails |
| TEST-W09-006 | Unit | BackupImportService.Validate | Invalid ZIP rejected with clear error |
| TEST-W09-007 | Unit | BackupImportService.Confirm | All entities inserted in correct order |
| TEST-W09-008 | Unit | BackupImportService.Confirm | Transaction rollback on any failure |
| TEST-W09-009 | Unit | BackupImportService.Confirm | Media files restored to correct paths |
| TEST-W09-010 | Integration | Backup ZIP/import round-trip | Export → validate → confirm restores identical data |
| TEST-W09-011 | Integration | Backup with media | Media included/excluded correctly |
| TEST-W09-012 | Unit | ExportArchive + BuildZIP (backup variant) | manifest.json and data.json format |
| TEST-W09-013 | Unit | BackupHandler.GenerateExport | Valid request → correct response |
| TEST-W09-014 | Unit | BackupHandler.DownloadExport | Valid downloadId → file streamed |
| TEST-W09-015 | Unit | BackupHandler.ImportValidate | Valid ZIP → summary returned |
| TEST-W09-016 | Unit | BackupHandler.ImportConfirm | Valid token → restore executed |
| TEST-W09-017 | Unit | BackupHandler error cases | Invalid/missing ZIP, expired/missing token, auth errors |

## Focused Commands
```bash
# Service unit tests
cd apps/api && go test ./internal/atlas/service/ -run TestBackup -v

# Handler tests
cd apps/api && go test ./internal/handler/ -run TestBackup -v

# Repository/integration tests
cd apps/api && go test ./internal/atlas/repository/postgres/ -run TestBackup -v

# Migration
cd apps/api && go test ./internal/repository/postgres/ -run Migration -v

# Build verification
cd apps/api && go build ./...
```