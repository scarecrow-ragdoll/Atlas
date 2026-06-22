<!-- FILE: .tasks/prd-wave-detail/20260621T000001Z/waves/WAVE-08/planner-product-ac-attempt-1.md -->
<!-- VERSION: 1.0.0 -->

# WAVE-08 Product-AC Planner Attempt 1

## Sources Read
- docs/prd-waves/waves/wave-08.md (source wave, user-approved)
- docs/product-verified/domain-model.md (AiReview entity, lifecycle)
- docs/product-verified/functional-spec.md §19 (REQ-015)
- docs/product-verified/acceptance-criteria.md (AC-025, AC-090, AC-091, AC-092)
- docs/product-verified/features/ai-review-history.md
- docs/product-verified/edge-cases.md
- docs/product-verified/actors-and-permissions.md
- docs/prd-waves/frontend-pages/page-009.md (consumer context)

## Selected Backend Wave Boundary
WAVE-08 delivers backend for storing AI review entries with manual text entry, date range linkage, user notes, and planned actions. Simple CRUD — no AI integration, no API calls.

## Neighboring Backend Wave Fit
- WAVE-07 (prior): Clean boundary. WAVE-07 owns AiExport and UserProfile. WAVE-08 owns AiReview. WAVE-07 explicitly defers AiReview to WAVE-08.
- WAVE-09 (future): WAVE-08 must provide service layer for WAVE-09 backup consumption.

## Frontend Pages Context
PAGE-009 (AI Export) may reference AiReview history. No dedicated AiReview frontend page exists in the page list. Backend-only wave.

## Codebase Evidence
AiReview entity has: id, userId, dateRangeStart, dateRangeEnd, aiResponseText, userNotes (optional), plannedActions (optional), createdAt, updatedAt. Lifecycle: created with response (simple create/read). No file storage, no external integration.

## Proposed Details

### Outcome After Implementation
- Save AI response text (manual paste)
- Link to date range (dateRangeStart, dateRangeEnd)
- Add user notes (optional)
- Track planned actions (optional TEXT field, MVP)
- Queryable review history (filterable by date range, sorted by creation date desc)

### Scope Included
- AiReview CRUD (Create, Read, List, Update, Delete)
- Period linkage (dateRangeStart, dateRangeEnd required)
- User notes (optional text)
- Planned actions storage (TEXT, MVP)
- Review history queries (list by user, filter by date range)
- Service layer interface for WAVE-09 consumption

### Scope Excluded
- Automatic AI call (explicitly excluded by source wave)
- OpenAI integration (explicitly excluded)
- Frontend pages, UI components, routes

### Acceptance Criteria Contributions

**AC-W08-001 — User can create an AI review entry**
Source: AC-025, AC-090. GraphQL createAiReview mutation accepts aiResponseText (required), dateRangeStart (required), dateRangeEnd (required), userNotes (optional), plannedActions (optional). Returns created AiReview with id.

**AC-W08-002 — User can paste AI response text**
Source: AC-090. aiResponseText field accepts arbitrary text. No server-side validation on content (plain text).

**AC-W08-003 — User can link review to a date range**
Source: AC-091. dateRangeStart and dateRangeEnd are required DATE fields. Validation: dateRangeEnd >= dateRangeStart.

**AC-W08-004 — User can add notes and planned actions**
Source: AC-092. userNotes (optional TEXT) and plannedActions (optional TEXT) accepted on create and update.

**AC-W08-005 — User can view review history**
Source: source wave outcome W08-005. aiReviews GraphQL query returns all reviews for the authenticated user, ordered by createdAt DESC.

**AC-W08-006 — User can filter reviews by date range**
Source: functional spec review history. aiReviews query accepts optional dateRangeStart and dateRangeEnd filters.

**AC-W08-007 — User can update a review entry**
Source: CRUD from source wave. updateAiReview mutation: update aiResponseText, userNotes, plannedActions, dateRangeStart, dateRangeEnd. Cannot change userId.

**AC-W08-008 — User can delete a review entry**
Source: CRUD from source wave. deleteAiReview mutation: deletes by id, user-scoped.

### Edge Cases Considered
- No feature-specific edge cases in docs/product-verified/edge-cases.md for AiReview
- Empty aiResponseText: should be rejected (required field)
- Date range end < start: validation error
- No reviews: aiReviews returns empty list
- Very long text in aiResponseText: no server-side limit (PostgreSQL TEXT handles up to 1GB; frontend may impose limit)
- Update after creation: standard CRUD pattern, no side effects

## Questions Raised
Q-W08-PAC-001: Should planned_actions be a simple TEXT field or a structured child table? PRD says "planned actions storage" — structured enables queryable action tracking. Simple TEXT matches MVP constraints. Recommendation: TEXT for MVP.

## Traceability Candidates
- AC-W08-001 → AC-025, §19
- AC-W08-002 → AC-090, §19.2
- AC-W08-003 → AC-091, §19.2
- AC-W08-004 → AC-092, §19.2
- AC-W08-005 → W08-005, functional spec §19
- AC-W08-006 → functional spec §19 "review history view"
- AC-W08-007 → CAP-W08-001 source wave
- AC-W08-008 → CAP-W08-001 source wave