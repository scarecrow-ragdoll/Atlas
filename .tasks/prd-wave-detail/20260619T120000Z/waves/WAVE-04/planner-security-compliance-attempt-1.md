# WAVE-04 Planner: Security / Compliance

## Auth Architecture

### PIN Auth Coverage
All WAVE-04 endpoints follow the WAVE-01 PIN auth pattern established in WAVE-02:
- All GraphQL mutations and queries protected by PIN auth middleware
- All REST endpoints (ProgressPhoto upload/download/delete) protected by PIN auth middleware
- When PIN is disabled, endpoints accessible without auth (consistent with TDEC-037)

### Auth Test Cases Per Endpoint
- `createCardioEntry`, `updateCardioEntry`, `deleteCardioEntry` → must return AuthError without valid session
- `cardioEntries`, `cardioEntry` → must return AuthError without valid session
- `createBodyWeightEntry`, `updateBodyWeightEntry`, `deleteBodyWeightEntry` → AuthError
- `bodyWeightEntries`, `bodyWeightEntry`, `latestBodyWeight` → AuthError
- `createBodyCheckIn`, `updateBodyCheckIn`, `deleteBodyCheckIn` → AuthError
- `createBodyMeasurement`, `updateBodyMeasurement`, `deleteBodyMeasurement` → AuthError
- `bodyCheckIns`, `bodyCheckIn` → AuthError
- `progressPhotos` → AuthError
- `createWeekFlag`, `deleteWeekFlag` → AuthError
- `weekFlags` → AuthError
- ProgressPhoto REST upload/download/delete → 401 without valid session

## Privacy Considerations

### Sensitive Data Categories
- Body weight, body fat percentage, body measurements — personal health data
- Progress photos — potentially identifiable imagery
- Weekly flags — health, stress, lifestyle data (potentially sensitive e.g. cycle)

### Protection Measures

**No sensitive data in logs:**
- Body weight values NOT logged (log only: created/updated weight entry ID, date)
- Body fat percentage NOT logged
- Body measurement values NOT logged
- Progress photo file content NOT logged
- Weekly flag notes NOT logged
- Log markers record only entity type, action, success/failure, and entity ID

**Media access control:**
- Progress photos stored on local filesystem under configured media path
- Served only through authenticated REST endpoint (PIN session required)
- No public access to photo files
- UUID-based filenames prevent enumeration

**File upload security:**
- Server-side MIME detection (`http.DetectContentType()`) as primary validation
- Allowed types: JPEG, PNG, WEBP (consistent with WAVE-02 but photos only — no video)
- Per-type size limits: 25MB
- UUID-based storage names prevent path traversal
- User-provided filenames sanitized: path separators removed, stored as metadata (originalFileName) only

### Data Sensitivity Classification

| Data | Sensitivity Level | Notes |
|---|---|---|
| Body weight | Medium | Personal health metric |
| Body fat % | Medium | Personal health metric |
| Body measurements | Medium | Body composition data |
| Progress photos | High | Potentially identifiable |
| Cardio entries | Low | Exercise data |
| Weekly flags | Medium | Health/lifestyle context |

## Compliance Considerations

### MVP Context
- Single-user deployment on own infrastructure
- No external data transmission (AI export is user-initiated download)
- No third-party data processors
- Self-hosted, user controls their data

### Recommendations
1. Add warning in AI export context that photos may contain identifiable information
2. Photos excluded from AI export by default (per RULE-025), explicit opt-in required
3. Body measurement and weight data excluded from default AI export? **Recommend:** include by default per PRD §17.3 (measurements included in export). Photos opt-in only.

## Edge Cases

### EDGE-011: PIN enabled but pinHash corrupted
- Consistent with WAVE-02 handling: service returns auth error, no data accessible

### EDGE-014: Photo URL accessed directly without session
- All media access goes through PIN-protected REST endpoint
- Physical files not served directly by web server

### EDGE-025: Disk full during photo upload
- Upload handler returns `INTERNAL_ERROR` on file write failure
- DB transaction not created unless file write succeeds — prevents orphaned records

## Open Security Questions

1. Q-W04-SEC-001: Should body weight and measurement data be excluded from AI export by default? **Recommend:** No — included by default per PRD. Photos are the only opt-in category.

2. Q-W04-SEC-002: Is there a max number of progress photos per check-in? **Recommend:** Hard limit of 10 per check-in (above the "2-4 recommended" guidance) to prevent storage abuse. Soft warning at 2, hard limit at 10.