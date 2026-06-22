<!-- FILE: docs/prd-wave-details/wave-map-context.md -->
<!-- VERSION: 1.0.0 -->

# Wave Map Context

## Selected Backend Wave Boundary
WAVE-08 (AI Review History): Store AI responses with period linkage and planned actions. Simple CRUD for AiReview entity. Backend-only: no AI call, no OpenAI integration, no file storage, no async processing.

## Prior Backend Wave Fit

### WAVE-01 (Foundation) — Required
Hard dependency on PIN auth middleware and atlas_users identity. Must be deployed.

### WAVE-02 (Exercise Library) — No dependency
### WAVE-03 (Workout Diary) — No dependency
### WAVE-04 (Cardio and Body Tracking) — No dependency
### WAVE-05 (Nutrition) — No dependency
### WAVE-06 (Charts) — No dependency

### WAVE-07 (AI Export) — Compatible
Clean boundary. WAVE-07 creates AiExport record and UserProfile. WAVE-08 creates independent AiReview record. No shared tables or services. WAVE-07 explicitly defers AiReview to WAVE-08: "AI review history (AiReview) — belongs to WAVE-08" (wave-07.md scope excluded).

## Future Backend Wave Fit

### WAVE-09 (Backup Import/Export)
WAVE-08 must provide AiReviewService.ListAllByUserID(ctx, userID) for WAVE-09 to include AiReview data in backup data.json. Read-only interface. No schema collision — WAVE-09 consumes data through service layer.

## Frontend Pages Context
PAGE-009 (AI Export): backend deps listed as POST /api/ai-export/generate, GET /api/ai-export/download, GET /api/user-profile, GET /api/week-flags. AiReview not listed as a dependency of PAGE-009. No dedicated AiReview frontend page in the pages list.

Backend provides GraphQL queries:
- Query aiReview(id: ID!): AiReviewResult!
- Query aiReviews(dateRangeStart: Date, dateRangeEnd: Date): AiReviewsResult!

These are dependency context only — no frontend planning.

## Dependency Order
WAVE-01 → WAVE-02 → WAVE-03 → WAVE-04 → WAVE-05 → WAVE-06 → WAVE-07 → WAVE-08 → WAVE-09

## Scope Collision Check
- No scope collision with WAVE-07 (AiExport and AiReview are distinct entities)
- No scope stolen from WAVE-09 (WAVE-09 consumes AiReview through read-only service interface)
- No scope overlap with WAVE-01 through WAVE-06
- AiReview is additive — no existing tables or services modified