<!-- FILE: .tasks/prd-wave-detail/20260621T000001Z/waves/WAVE-08/planner-sequencing-fit-attempt-1.md -->
<!-- VERSION: 1.0.0 -->

# WAVE-08 Sequencing-Fit Planner Attempt 1

## Sources Read
- docs/prd-waves/waves/wave-08.md (source wave)
- docs/prd-waves/wave-map.md
- docs/prd-waves/frontend-pages/page-009.md
- docs/prd-wave-details/waves/wave-01.md (Foundation)
- docs/prd-wave-details/waves/wave-07.md (AI Export — direct neighbor)
- docs/prd-wave-details/index.md
- docs/knowledge-graph.xml

## Selected Backend Wave Boundary
WAVE-08 delivers AiReview CRUD with period linkage. Simple storage, no AI integration.

## Neighboring Backend Wave Fit

### Prior Wave Compatibility

**WAVE-01 (Foundation)** — Required. PIN auth middleware for all GraphQL operations. atlas_users for user identity. Settings for any future defaults. WAVE-01 is fully deployed.

**WAVE-02 (Exercise Library)** — Compatible. No dependency.

**WAVE-03 (Workout Diary)** — Compatible. No dependency.

**WAVE-04 (Cardio and Body Tracking)** — Compatible. No dependency.

**WAVE-05 (Nutrition)** — Compatible. No dependency.

**WAVE-06 (Charts)** — Compatible. No dependency.

**WAVE-07 (AI Export and Prompt Builder)** — Compatible. Clean boundary per wave-07.md "AI review history (AiReview) — belongs to WAVE-08". WAVE-07 creates AiExport; WAVE-08 creates independent AiReview. No shared tables, no shared services.

### Future Wave Compatibility

**WAVE-09 (Backup Import/Export)** — WAVE-08 must provide service layer for WAVE-09 consumption. AiReviewService must expose ListAllByUserID(ctx, userID) so WAVE-09 can include AiReview data in backup data.json. This is a read-only interface — no other coupling needed.

## Dependency Order
WAVE-01 → WAVE-02 → WAVE-03 → WAVE-04 → WAVE-05 → WAVE-06 → WAVE-07 → **WAVE-08** → WAVE-09

WAVE-08 depends on WAVE-01 (mandatory) and WAVE-07 (compatibility boundary only, not code dependency). WAVE-09 depends on WAVE-08 service interface.

## Scope Collision Check
- No scope collision with any prior wave
- No scope stolen from WAVE-09 (WAVE-09 consumes AiReview data, WAVE-08 provides service layer)
- No scope stolen from WAVE-07 (WAVE-07 explicitly defers AiReview)

## Frontend Pages Context
PAGE-009 (AI Export): backend deps listed as POST /api/ai-export/generate, GET /api/ai-export/download, GET /api/user-profile, GET /api/week-flags. AiReview not listed as a dependency of PAGE-009. No dedicated AiReview frontend page in pages list.

Backend provides GraphQL queries for any future AiReview frontend component:
- Query aiReviews(dateRangeStart, Date, dateRangeEnd: Date): AiReviewsResult!
- Query aiReview(id: ID!): AiReviewResult!

These are dependency context only — no frontend planning.

## Independent Deliverability
- Can be implemented without WAVE-02 through WAVE-07 (no code dependency)
- **Cannot be implemented without WAVE-01** — hard dependency on PIN auth and user identity

## Questions Raised
Q-W08-SEQ-001: WAVE-09 needs ListAllByUserID on AiReviewService for backup inclusion. Confirm this contract. Recommendation: expose ListAllByUserID(ctx, userID) returning []AiReview.

## Traceability Candidates
- WAVE-07 boundary → wave-07.md "AI review history (AiReview) — belongs to WAVE-08"
- WAVE-09 contract → wave-07.md "Future Wave Compatibility"
- Dependency order → docs/knowledge-graph.xml module order