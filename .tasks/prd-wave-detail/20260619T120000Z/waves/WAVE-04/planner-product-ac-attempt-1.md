# WAVE-04 Planner: Product / AC

## Scope Analysis

### Capability Groups (from source wave CAPs)
- CAP-W04-001: CardioEntry CRUD — type (walk/run/bike/elliptical/treadmill/other), duration (minutes), avg pulse (optional), heart rate zone (1-5/unknown, optional), notes (optional). Attached to DailyLog via dailyLogId.
- CAP-W04-002: BodyWeightEntry CRUD — date, weight, source enum, notes (optional). Standalone per date, user-level.
- CAP-W04-003: BodyCheckIn CRUD — date, weight (optional), bodyFatPercentage (optional), notes. 1:N to BodyMeasurement and ProgressPhoto.
- CAP-W04-004: BodyMeasurement CRUD — checkInId, measurementType (neck/shoulders/forearms/biceps/chest/waist/abdomen/hips/thigh/calf), side (left/right/null for unpaired), value.
- CAP-W04-005: ProgressPhoto CRUD — checkInId, filePath, originalFileName, mimeType, sizeBytes, angle (front/side/back/custom), label (optional), notes (optional).
- CAP-W04-006: WeekFlag CRUD — weekStartDate, flagType (poor-sleep/high-stress/illness/injury-pain/cycle/calorie-deficit/surplus/maintenance/missed-workouts/travel), notes (optional).

### Product ACs Mapped to WAVE-04

| Product AC | WAVE-04 Mapping | Notes |
|---|---|---|
| AC-012 | SLICE Cardio — user can add cardio with type and duration | Covered |
| AC-013 | SLICE Cardio — user can optionally specify pulse and heart rate zone | Covered (optional fields) |
| AC-014 | SLICE CheckIn — user can create weekly body check-in | Covered |
| AC-015 | SLICE CheckIn + Measurements + Photos — check-in includes date, weight, optional body fat %, measurements, 2-4 photos, comment | Photo count 2-4 is requirement vs recommendation (EDGE-006) |
| AC-016 | SLICE BodyWeight — user can enter body weight for any date | Covered |
| AC-048 | SLICE Cardio — user can select cardio type from predefined list | Covered (enum validation) |
| AC-049 | SLICE Cardio — cardio duration in minutes | Covered (durationMinutes) |
| AC-050 | SLICE Cardio — optionally record average pulse | Covered (avgPulse optional) |
| AC-051 | SLICE Cardio — optionally select heart rate zone 1-5/unknown | Covered (heartRateZone optional) |
| AC-052 | SLICE CheckIn — weekly check-in records date, weight, optional body fat % | Covered |
| AC-053 | SLICE Measurements — check-in includes 10 measurement types | Covered |
| AC-054 | SLICE Measurements — paired measurements can record left, right, or single | Covered (side field) |
| AC-055 | SLICE Measurements — second value in paired measurement not required | Covered (side nullable) |
| AC-056 | SLICE Photos — 2-4 photos per check-in | EDGE-006: requirement vs recommendation |
| AC-057 | SLICE BodyWeight — standalone weight entry for any date | Covered |

### Edge Cases

| Edge Case | WAVE-04 Impact | Resolution |
|---|---|---|
| EDGE-006 | Check-in with 0-1 photos — "2-4 photos" requirement ambiguous | Recommend: soft guidance (2-4 recommended), no hard block. Photo count ≥ 1 for check-in to have photos; allow 0 photos to not prevent check-in creation. |
| EDGE-007 | Body measurement value 0 or negative | Recommend: value > 0 validation on all measurement types. Reject 0 and negative. |

### Domain Model Invariants

1. CardioEntry.dailyLogId is required per domain model invariant #2 (Cardio must belong to a DailyLog). System auto-creates DailyLog when needed.
2. BodyCheckIn 1:N BodyMeasurement — measurements exist only within check-in context. Deleting check-in cascades to measurements and photos.
3. BodyCheckIn 1:N ProgressPhoto — photos exist only within check-in context.
4. BodyMeasurement.side: paired types (forearm, biceps, thigh, calf) may have left/right/null. Unpaired (neck, shoulders, chest, waist, abdomen, hips) — side must be null.
5. Check-in photo range: 2-4 per business rule RULE-005 (requirement ambiguous — see EDGE-006).

### Open Questions (pre-planning)

1. EDGE-006: Is 2-4 photos a hard requirement or recommendation? **Recommendation** — soft guidance, warn but allow.
2. What is the BodyWeightEntry.source enum? Not defined in domain model. **Recommend:** scale/manual/unknown.
3. BodyCheckIn.weight is optional per domain model — should we allow a check-in with no weight AND no bodyFatPercentage? Yes, minimal valid check-in has date + notes.
4. ProgressPhoto angle: front/side/back/custom — is custom a free-text label or fixed enum? **Recommend:** fixed enum plus optional label for elaboration.

### Proposed WAVE-04 ACs

| AC ID | Description | Source |
|---|---|---|
| AC-W04-001 | Cardio entry can be created via GraphQL mutation with cardioType (enum), durationMinutes, optional avgPulse, optional heartRateZone (enum 1-5/unknown), optional notes. |
| AC-W04-002 | Cardio entry type is validated against allowed enum values. Invalid type returns ValidationError. |
| AC-W04-003 | Cardio durationMinutes is required, positive integer. 0 or negative returns ValidationError. |
| AC-W04-004 | Cardio entry avgPulse, if provided, is positive integer. |
| AC-W04-005 | Cardio entry heartRateZone, if provided, is 1-5 or "unknown". Invalid zone returns ValidationError. |
| AC-W04-006 | Cardio entry is linked to a DailyLog (auto-created if needed for the date). |
| AC-W04-007 | Cardio entry can be read by ID. |
| AC-W04-008 | Cardio entries can be listed by date (DailyLog). |
| AC-W04-009 | Cardio entry can be updated (type, duration, pulse, zone, notes). |
| AC-W04-010 | Cardio entry can be deleted. |
| AC-W04-011 | BodyWeightEntry can be created with date (required), weight (required, > 0), source enum (required), optional notes. |
| AC-W04-012 | BodyWeightEntry weight must be > 0. 0 or negative returns ValidationError. |
| AC-W04-013 | BodyWeightEntry source is validated against allowed enum (scale/manual/unknown). |
| AC-W04-014 | BodyWeightEntry can be read by ID. |
| AC-W04-015 | BodyWeightEntry can be listed with date range filtering. |
| AC-W04-016 | BodyWeightEntry latest weight can be queried (for dashboard). |
| AC-W04-017 | BodyWeightEntry can be updated (weight, source, notes). |
| AC-W04-018 | BodyWeightEntry can be deleted. |
| AC-W04-019 | BodyCheckIn can be created with date (required), optional weight, optional bodyFatPercentage, optional notes. |
| AC-W04-020 | BodyCheckIn weight, if provided, must be > 0. |
| AC-W04-021 | BodyCheckIn bodyFatPercentage, if provided, must be > 0 and <= 100. |
| AC-W04-022 | BodyCheckIn can be read by ID with nested measurements and photos. |
| AC-W04-023 | BodyCheckIn can be listed with date range ordering (descending). |
| AC-W04-024 | BodyCheckIn can be updated (weight, bodyFatPercentage, notes). |
| AC-W04-025 | BodyCheckIn can be deleted (cascade to measurements and photos including physical file deletion). |
| AC-W04-026 | BodyMeasurement can be created within a check-in with measurementType (enum), value (> 0), optional side (left/right for paired types). |
| AC-W04-027 | BodyMeasurement measurementType validated against 10 allowed types. Invalid type returns ValidationError. |
| AC-W04-028 | BodyMeasurement value must be > 0. 0 or negative returns ValidationError. |
| AC-W04-029 | BodyMeasurement side is validated: allowed only for paired types (forearm, biceps, thigh, calf). Must be null for unpaired types. |
| AC-W04-030 | BodyMeasurement can be updated (value, side). |
| AC-W04-031 | BodyMeasurement can be deleted. |
| AC-W04-032 | ProgressPhoto can be uploaded (multipart REST) associated with a check-in, with angle (enum), optional label, optional notes. |
| AC-W04-033 | ProgressPhoto angle is validated: front/side/back/custom. Invalid angle returns ValidationError. |
| AC-W04-034 | ProgressPhoto file is stored in media storage with path <BasePath>/progress-photos/<checkin_id>/<uuid>.<ext>. |
| AC-W04-035 | ProgressPhoto file is validated server-side: allowed MIME types (JPEG/PNG/WEBP), per-type size limits (25MB). |
| AC-W04-036 | ProgressPhoto can be downloaded via GET endpoint (PIN-protected). |
| AC-W04-037 | ProgressPhoto can be deleted (cascade removes physical file). |
| AC-W04-038 | ProgressPhotos can be listed by check-in ID. |
| AC-W04-039 | WeekFlag can be created with weekStartDate, flagType (enum), optional notes. |
| AC-W04-040 | WeekFlag flagType validated against allowed enum values (poor-sleep/high-stress/illness/injury-pain/cycle/calorie-deficit/surplus/maintenance/missed-workouts/travel). |
| AC-W04-041 | WeekFlag can be listed by week start date. |
| AC-W04-042 | WeekFlag can be deleted. |
| AC-W04-043 | All WAVE-04 GraphQL mutations return AuthError when PIN session header is missing or invalid. |
| AC-W04-044 | All WAVE-04 REST endpoints return 401 when PIN session header is missing or invalid. |