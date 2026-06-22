# Planner Report: Security & Compliance (WAVE-09)

## Scope
Auth, authorization, privacy, audit, compliance for backup/import operations.

## Auth & Authorization
- Atlas uses optional PIN guard (no role-based access, no multi-tenancy)
- All endpoints behind AtlasPinGuard middleware (same as existing Atlas routes)
- `/api/v1/auth/pin/unlock` and session management apply to all backup endpoints
- No additional authorization needed — single-user instance

## Endpoint Authentication

| Endpoint | Auth Required | Notes |
| --- | --- | --- |
| POST /api/backup/export | AtlasPinGuard | Same pattern as /api/ai-export/generate |
| GET /api/backup/download | AtlasPinGuard | Same pattern as /api/ai-export/download |
| POST /api/backup/import/validate | AtlasPinGuard | New endpoint, same auth |
| POST /api/backup/import/confirm | AtlasPinGuard | New endpoint, same auth |

## Privacy

### Data Content in Backup
Backup contains ALL user data: settings, profile, exercises, workouts, cardio, body measurements, nutrition, AI exports, AI reviews, media files.

### Logging Policy (per AC-117-120)
- **AC-117:** PIN not logged ✓ (backup doesn't touch PIN)
- **AC-118:** AI export content not logged — backup includes AI exports, must not log them
- **AC-119:** Photos not logged — backup includes media, file paths should not be logged
- **AC-120:** Sensitive comments not logged — backup includes exercise/nutrition comments

**Rule:** Log metadata only (operation type, status, entity count, size, timestamp). Never log entity values, file paths, or content.

### Media Files
- Backup ZIP media files follow same access rules as existing media: require valid PIN session
- Export: media copied into ZIP only during generation
- Import: media extracted from ZIP to media storage directory
- File storage pattern (temp file → atomic rename) from atlas_media.go applies

## Compliance
- No GDPR/regulatory requirements for self-hosted single-user MVP
- Data ownership principle: user can export all data and delete instance
- Import validation prevents data corruption

## Risks
1. **DQ-W09-004:** Logging policy for backup operations needs confirmation — recommended: log event metadata, not content
2. **Large import ZIP DoS:** Must limit upload size via MaxBytesReader (same pattern as media upload)
3. **Import validation token reuse:** Validation token should be one-time-use to prevent replay attacks