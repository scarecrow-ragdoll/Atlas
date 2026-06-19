# WAVE-02 sequencing-other-wave-fit Review Attempt 1

## Verdict
needs-revision

## Sources Read
- planner-sequencing-fit-attempt-1.md
- planner-architecture-codebase-attempt-1.md
- planner-data-integration-ops-attempt-1.md
- planner-product-ac-attempt-1.md
- planner-testing-exit-attempt-1.md
- docs/prd-waves/wave-map.md
- docs/prd-waves/waves/wave-02.md
- docs/prd-wave-details/waves/wave-01.md
- docs/prd-waves/frontend-pages/page-002.md
- docs/prd-waves/frontend-pages/page-003.md
- docs/technical-verified/implementation-slices.md

## Coverage Check
Dependency mapping covers: WAVE-01 → WAVE-02 (all dependencies), WAVE-02 → WAVE-03 (allExercises query), WAVE-02 → WAVE-06/07/08/09 (data source). No scope collisions identified.

## Evidence Check
Dependency order confirmed in wave-map.md and WAVE-01's "Dependencies And Other-Wave Fit" section. Source wave gate (passed) confirms readiness.

## Codebase Fit Check
All dependency claims are structural (no runtime conflicts). Codebase fit is confirmed by architecture-codebase planner.

### Issues Found

1. **WAVE-01 dependency completeness**: The planner correctly identifies WAVE-01 dependencies (PIN auth, media scaffold, migration infra, codegen config). However, it does not list what happens if WAVE-01 is not yet implemented. Since WAVE-01 is `ready-for-dev` but may not be user-approved or implemented yet, the sequencing plan should include a contingency: if WAVE-01 implementation is delayed, WAVE-02 cannot start. This is a valid sequencing dependency that should be explicitly documented.

2. **WAVE-01 media scaffold assumption**: The planner assumes WAVE-01's media REST scaffold stores files at a configurable base path and returns media IDs. If WAVE-01 uses a different media model (e.g., stores directly in DB as bytea, or uses a different storage pattern), WAVE-02 must adapt. The planner should specify the exact WAVE-01 media contract (input/output) that WAVE-02 depends on.

3. **WAVE-03 allExercises API contract**: The planner says WAVE-02 provides `allExercises` for WAVE-03, but doesn't specify whether this is a GraphQL query or REST endpoint. PAGE-002 lists `GET /api/exercises` as REST. The planner should align — either provide both GraphQL and REST, or choose one and document the interface. Recommendation: GraphQL only (consistent with hybrid pattern), and PAGE-003/002 frontends use GraphQL for all exercise queries.

4. **WAVE-06 chart data compatibility**: The planner says WAVE-06 (Charts) reads exercise data for training charts. WAVE-02 stores workingWeight as a single current value. WAVE-06 chart queries need historical working weight data. This is inconsistent — WAVE-06 will get historical data from WAVE-03's WorkoutExercise.workingWeightSnapshot, not from WAVE-02's current workingWeight. The planner should correct this: WAVE-02 provides exercise metadata (name, muscle groups, active status) for chart labels; WAVE-03 provides historical working weight data for trend charts.

5. **Missing WAVE-09 backup compatibility note**: The planner says WAVE-09 includes exercise data in backup, but doesn't specify the data export format. WAVE-09 needs exercises + exercise_media data in its export JSON. This is a forward contract that WAVE-02 data should be exportable. The planner should add a note that WAVE-02 data model should be backup-compatible (serializable to JSON, all fields exportable).

6. **Migration numbering collision**: WAVE-01's migration plan is not yet finalized in code. WAVE-02 planners propose 00080 and 00081, but WAVE-01 might use migrations in the 00080+ range. Need to document that migration numbering must be confirmed after WAVE-01 implementation.

## Other-Wave Fit Check
All 8 other waves checked for collision. Clean separation confirmed.

## Acceptance Criteria Check
AC-W02-021 (allExercises for WAVE-03) and AC-W02-022 (query by ID after soft-delete) are validated as correct forward contracts.

## Exit Criteria Check
EC-W02-022 (no WAVE-03 functionality) is hard to prove (also flagged by testing-exit reviewer). Move to code review checklist.

## Verification Check
Not directly applicable, but TEST-W02-007 (simple list for WAVE-03) validates the most critical forward dependency.

## Question Ledger Check
DQ-W02-008 (allExercises filtering) is valid — if WAVE-03 needs filtered lists, the query surface grows. But for MVP, unfiltered list of active exercises is sufficient. The question should be moved to watchlist rather than needs-owner-decision since MVP scope is clear.

## Unsupported Or Invented Claims
No unsupported claims. The WAVE-02 → WAVE-06 data flow assumption (that WAVE-06 uses WAVE-02 workingWeight for charts) is technically incorrect — WAVE-06 needs historical data from WAVE-03. This should be corrected in the planner.

## Required Revisions
1. **Document WAVE-01 implementation dependency**: Explicitly state WAVE-02 cannot start before WAVE-01 is implemented.
2. **Specify exact WAVE-01 media contract**: Document the precise interface WAVE-01's media scaffold provides.
3. **Align allExercises interface**: Specify GraphQL-only or add REST endpoint consistent with PAGE-002 requirements.
4. **Correct WAVE-06 data flow**: WAVE-02 provides metadata; WAVE-03 provides historical working weight data.
5. **Add WAVE-09 backup compatibility note**: WAVE-02 data must be serializable for backup export.
6. **Document migration numbering coordination**: Confirm WAVE-01 migration range before finalizing WAVE-02 migration files.
7. **Move DQ-W02-008 to watchlist**: Not a blocker for MVP.

## Approval Notes
Strong sequencing analysis. The revision items are mostly clarifying — no fundamental sequencing conflicts. After revisions, will approve.