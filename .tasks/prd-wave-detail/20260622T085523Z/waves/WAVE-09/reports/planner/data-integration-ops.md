# Planner Report: Data, Integration & Operations (WAVE-09)

## Scope
Data lifecycle, API endpoints, events/jobs, external integrations, observability, rollout/rollback.

## Data Lifecycle

### Export Flow
```
User clicks export → POST /api/backup/export
  → BackupService.Generate(ctx, userID, includeMedia)
    → collect all entities via each service's ListAllByUserID
    → build BackupArchive (ExportArchive with BackupManifest + BackupData)
    → call archive.BuildZIP() → get []byte
    → write to temp file → atomic rename to final path
    → return {downloadId, size, timestamp}

User downloads → GET /api/backup/download?downloadId=X
  → look up generated file on disk → stream ZIP to client
```

### Import Flow (multi-step)
```
Step 1: POST /api/backup/import/validate
  → upload ZIP file (multipart)
  → parse manifest.json → validate structure
  → parse data.json → validate schema version
  → run dry-run: count entities, check foreign keys
  → store validation result in memory (keyed by session token)
  → return {validationId, summary: {entityCounts, schemaVersion, mediaCount, warnings}}

Step 2: POST /api/backup/import/confirm
  → provide validationId from step 1
  → begin PostgreSQL transaction
  → DELETE all existing user data (if replace strategy — TBD in DQ-W09-001)
  → INSERT all entities in dependency order
  → restore media files
  → COMMIT on success, ROLLBACK on any failure
  → return {status: "success", entityCounts, mediaCount}
```

### All-or-Nothing Transaction Strategy
Use a single PostgreSQL transaction wrapping ALL restore INSERTs:
1. BEGIN
2. DELETE existing data (if replace strategy)
3. INSERT settings
4. INSERT user profile
5. INSERT exercises → exercise media
6. INSERT daily logs → workout exercises → workout sets
7. INSERT cardio entries
8. INSERT body check-ins → measurements → photos
9. INSERT body weight entries
10. INSERT nutrition products → templates → template items → overrides → override items
11. INSERT week flags
12. INSERT AI exports
13. INSERT AI reviews
14. COMMIT / ROLLBACK

### Import State Management
Validation state between step 1 and step 2:
- **Option A (recommended):** Store in-memory map keyed by a token returned to client
- **Option B:** Store in Redis with TTL (e.g., 15 minutes)
- Risk: server restart between validate and confirm loses state → user re-uploads

## API Endpoints

| Method | Path | Purpose | Request Body | Response |
| --- | --- | --- | --- | --- |
| POST | /api/backup/export | Generate backup ZIP | {includeMedia: boolean} | {downloadId, size, timestamp} |
| GET | /api/backup/download | Download generated ZIP | query: downloadId | application/zip stream |
| POST | /api/backup/import/validate | Upload ZIP, dry-run check | multipart: file | {validationId, summary} |
| POST | /api/backup/import/confirm | Confirm and execute restore | {validationId} | {status, entityCounts} |

## Observability

### Log Markers (following existing pattern)
- `[Backup][export][BLOCK_EXPORT_START]` — export begins
- `[Backup][export][BLOCK_EXPORT_DATA_QUERY]` — collecting entities
- `[Backup][export][BLOCK_EXPORT_ZIP_BUILD]` — building ZIP
- `[Backup][export][BLOCK_EXPORT_ZIP_WRITE]` — writing to disk
- `[Backup][export][BLOCK_EXPORT_SUCCESS]` — export complete
- `[Backup][export][BLOCK_EXPORT_FAILURE]` — export error
- `[Backup][import][BLOCK_IMPORT_VALIDATE]` — dry-run validation
- `[Backup][import][BLOCK_IMPORT_CONFIRM]` — confirmation received
- `[Backup][import][BLOCK_IMPORT_START]` — transaction begins
- `[Backup][import][BLOCK_IMPORT_COMMIT]` — transaction committed
- `[Backup][import][BLOCK_IMPORT_ROLLBACK]` — transaction rolled back
- `[Backup][import][BLOCK_IMPORT_FAILURE]` — import error

### Privacy-Logging Rules (per AC-117-120)
- Log backup event metadata only: operation, status, size, timestamp, entity counts
- Do NOT log backup content, entity values, file paths, or media data

## Performance
- db-only export <= 15s p95 (per product-brief)
- with media best-effort
- db-only import <= 30s p95
- Max export size should be configurable (BackupConfig.maxExportSizeBytes)
- Max import upload size should be configurable (recommended: 500MB default)

## Rollout/Rollback
- New migration (00094) — standard apply/rollback via existing migration framework
- New config section (BackupConfig) — hot-reloadable via config file restart
- New REST endpoints — no impact on existing routes
- New service dependencies — wired at startup, no runtime migration