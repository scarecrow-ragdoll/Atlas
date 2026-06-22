<!-- FILE: .tasks/prd-wave-detail/20260621T000001Z/waves/WAVE-08/planner-security-compliance-attempt-1.md -->
<!-- VERSION: 1.0.0 -->

# WAVE-08 Security-Compliance Planner Attempt 1

## Sources Read
- docs/product-verified/actors-and-permissions.md
- docs/product-verified/edge-cases.md
- docs/product-verified/business-rules.md
- docs/prd-wave-details/waves/wave-07.md (security section)
- docs/prd-waves/waves/wave-08.md

## Selected Backend Wave Boundary
Simple CRUD for user-managed text entries. No external data transmission.

## Neighboring Backend Wave Fit
WAVE-07 established PIN-guard middleware pattern and log privacy pattern. WAVE-08 follows the same.

## Frontend Pages Context
No dedicated frontend page. Backend provides GraphQL-only access.

## Codebase Evidence
- All GraphQL resolvers use middleware.GetAtlasUserID for auth
- PIN guard enforced at middleware level for all API endpoints
- AC-118 establishes log privacy pattern (no content in logs)

## Proposed Details

### Authentication
- All GraphQL mutations/queries protected by PIN-guard middleware (same as WAVE-07)
- middleware.GetAtlasUserID extracts user from session context
- Non-authenticated requests return 401

### Authorization
- AiReview operations scoped to authenticated user only
- UserID from session context used for all queries/filtering
- No cross-user access possible (single-tenant design, but defense-in-depth)
- Ownership check on GetByID (verify userId matches session)

### Privacy
- ai_response_text, user_notes, planned_actions: user-generated content stored as-is
- No sensitive data classification needed (single-user, manual entry)
- Log privacy: no ai_response_text, user_notes, or planned_actions content in application logs
- Only metadata (review ID, date range, operation type) logged
- Following AC-118 pattern established by WAVE-07

### Data Storage
- All data in PostgreSQL (encrypted at rest via volume encryption)
- No file storage needed
- No external data transmission

### Audit
- Standard audit fields (created_at, updated_at) on AiReview
- No compliance requirements beyond standard data integrity

### Rate Limiting
- No specific rate limits for MVP (single-user, local deployment)
- Standard API guardrails apply

### Abuse Prevention
- No abuse vectors identified (single-user, manual entry, no external API)
- Standard validation guards (required fields, date range validation)

## Questions Raised
Q-W08-SEC-001: Should ai_response_text be encrypted at rest? Recommendation: No for MVP (single-user local deployment, standard PostgreSQL volume encryption sufficient). Consider for multi-tenant future.

## Traceability Candidates
- PIN guard → WAVE-01 foundation, WAVE-07 pattern
- Log privacy → AC-118 pattern from WAVE-07
- User-scoped queries → standard authorization pattern