# WAVE-04 Review: Data / API / Integration / Ops

## Review Cycle
1

## Planner Reports Reviewed
- planner-data-integration-ops-attempt-1.md
- planner-architecture-codebase-attempt-1.md
- planner-product-ac-attempt-1.md
- planner-sequencing-fit-attempt-1.md

## Verdict
approved

## Findings

### Data Model Design
- 6 tables with correct column types, constraints, and indexes ✓
- body_check_in.date UNIQUE constraint ✓ (one check-in per date)
- body_measurement UNIQUE (check_in_id, measurement_type, side) ✓
- week_flag UNIQUE (week_start_date, flag_type) ✓
- Cascade rules: check-in delete cascades to measurements and photos ✓
- Cardio FK to daily_log with ON DELETE CASCADE ✓

### API Design
- GraphQL types, enums, inputs, unions follow WAVE-02 pattern ✓
- Union results for mutations (Success | ValidationError | AuthError) ✓
- REST endpoints for ProgressPhoto consistent with exercise media pattern ✓
- Error codes and format consistent with TDEC-027 ✓

### Operations
- Log markers follow established pattern ✓
- Migration strategy (goose, sequential) ✓
- Media storage path pattern established ✓
- File validation (MIME, size) ✓

### Open Design Decisions

1. **BodyWeightEntry duplicate per date**: The planner recommends allowing multiple entries per date (scale vs manual). The ACs should clarify this. **Recommend:** Allow duplicates, but add a note that the `bodyWeightEntries` query returns all. The dashboard `latestBodyWeight` returns most recent by created_at.

2. **Measurement side validation**: The planner specifies side allowed only for paired types (forearm, biceps, thigh, calf). This must be enforced in the service layer, not just the DB. Add AC for this (already in AC-W04-029 ✓).

3. **DailyLog auto-creation**: The createCardioEntry mutation accepts a date and auto-creates a DailyLog if none exists. This is correct per domain model invariant #2. However, if WAVE-03's daily_log table doesn't exist yet, this will fail. **Recommend:** Add a design decision (DDEC) documenting the dependency and fallback behavior.

### Required Revisions
- None. Data/API/ops design is consistent with WAVE-02 patterns.

## Notes
- BodyWeightEntry source enum: scale/manual/unknown — reasonable, but confirm with product owner
- 10 measurement types match PRD spec exactly
- Photo angle enum: front/side/back/custom — correct per PRD