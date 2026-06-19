# WAVE-04 Planner: Sequencing / Other-Wave Fit

## Dependency Check

### WAVE-01 (Foundation) — Prerequisite
- **PIN auth middleware** — required for all WAVE-04 endpoints ✓
- **Media storage scaffold** — required for ProgressPhoto upload/download ✓
- **Config pattern** — MediaConfig extension ✓
- **GraphQL foundation** — common types, extend type pattern ✓
- **Test infrastructure** — test helpers, test DB ✓
- **Status:** WAVE-01 must be implemented before WAVE-04 can start

### WAVE-02 (Exercise Library) — No direct dependency
- WAVE-04 does not depend on exercise data, exercise media, or exercise queries
- **Parallelization possible:** Yes — WAVE-04 and WAVE-02 share no tables or logic
- Exception: both depend on WAVE-01 PIN auth middleware

### WAVE-03 (Workout Diary) — Dependency for DailyLog
- **Critical dependency:** WAVE-04 CardioEntry requires `dailyLogId` FK to `daily_log` table
- If WAVE-03 creates `daily_log`, WAVE-04 must wait for WAVE-03 migrations OR...
- **Alternative:** WAVE-04 includes migration for daily_log table if not yet present
- **Parallelization possible:** Partial — if daily_log table defined in WAVE-03, WAVE-04 must sequence after WAVE-03 migrations. If WAVE-04 creates daily_log, they can parallelize.
- **Recommendation:** WAVE-04 should assume daily_log exists (from WAVE-03) and include a migration dependency note. If WAVE-03 not deployed, WAVE-04 creates daily_log migration or defers cardio auto-creation to handle missing table gracefully.

### WAVE-05 (Nutrition) — No direct dependency
- WAVE-04 and WAVE-05 share no tables, no data flow
- **Parallelization possible:** Yes — fully independent

### WAVE-06 (Charts) — Consumer
- WAVE-06 consumes body weight data (body_weight_entry), body measurements (body_measurement), check-in data (body_check_in)
- WAVE-04 provides the data; no reverse dependency
- WAVE-04 API contracts must be stable before WAVE-06 implementation

### WAVE-07/08 (AI Export/Review) — Consumer
- WAVE-07 consumes cardio entries, body weight, check-in data, progress photos, week flags
- WAVE-04 provides the data through GraphQL queries and service layer
- WAVE-04 API contracts must be stable before WAVE-07 implementation

### WAVE-09 (Backup) — Consumer
- WAVE-09 consumes all WAVE-04 tables for full backup export/import
- WAVE-04 repository/service layer used by backup module
- WAVE-04 API contracts must be stable before WAVE-09 implementation

## Frontend Page Dependencies (read-only context)

### PAGE-004 (Cardio)
- Backend deps: GET/POST/PUT/DELETE /api/cardio (via GraphQL mutations/queries)
- WAVE-04 provides: createCardioEntry, updateCardioEntry, deleteCardioEntry, cardioEntries, cardioEntry
- Date filtering on cardioEntries query required

### PAGE-005 (Body Measurements)
- Backend deps: GET /api/body-check-ins, POST /api/body-check-ins, POST /api/measurements, POST /api/progress-photos, GET /api/body-weight
- WAVE-04 provides: bodyCheckIns, bodyCheckIn, createBodyCheckIn, updateBodyCheckIn, deleteBodyCheckIn, createBodyMeasurement, updateBodyMeasurement, deleteBodyMeasurement, bodyWeightEntries, createBodyWeightEntry, latestBodyWeight
- Nested measurements and photos in check-in query required

### PAGE-006 (Progress Photos)
- Backend deps: GET /api/progress-photos?checkin={id}, DELETE /api/progress-photos/{id}
- WAVE-04 provides: progressPhotos query, REST upload/download/delete

### PAGE-001 (Dashboard)
- Backend deps: GET /api/body-weight/latest
- WAVE-04 provides: latestBodyWeight query

## Wave Boundary Verification

### No Scope Collision
- CardioEntry: unique to WAVE-04
- BodyWeightEntry: unique to WAVE-04
- BodyCheckIn: unique to WAVE-04
- BodyMeasurement: unique to WAVE-04
- ProgressPhoto: unique to WAVE-04
- WeekFlag: unique to WAVE-04
- No overlap with WAVE-01, WAVE-02, WAVE-03, WAVE-05, WAVE-06

### Stability Requirements
- All WAVE-04 GraphQL type definitions and query signatures must be stable before WAVE-06/07/08 implementation
- Repository/service interfaces used by backup module must be stable before WAVE-09

## Risk Assessment

| Risk | Impact | Mitigation |
|---|---|---|
| WAVE-03 daily_log table not available | Blocking for cardio FK | WAVE-04 includes conditional daily_log migration or auto-creation logic |
| WAVE-01 media scaffold not available | Blocking for progress photo upload | Document as WAVE-01 dependency — cannot start until WAVE-01 complete |
| WAVE-04/05 parallelization conflicts | Low | No shared tables. Independent schema files. Sequential migration numbering. |
| WAVE-06 relies on WAVE-04 data model changes | Medium | Document WAVE-04 model decisions. Keep measurement types enum extensible. |