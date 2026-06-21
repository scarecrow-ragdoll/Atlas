# WAVE-06 Security-Compliance Planner Attempt 1

## Sources Read
- docs/prd-waves/waves/wave-06.md
- docs/technical-verified/auth-security-compliance.md
- docs/product-verified/actors-and-permissions.md
- docs/prd-wave-details/waves/wave-04.md (Security section)
- docs/prd-wave-details/waves/wave-05.md (Security section)
- apps/api/internal/atlas/graph/resolver/resolver.go
- apps/api/internal/atlas/middleware/

## Selected Backend Wave Boundary
Read-only queries returning aggregated user data. No mutations, no media, no external calls.

## Neighboring Backend Wave Fit
Follows same PIN auth pattern as WAVE-04 and WAVE-05.

## Frontend Pages Context
PAGE-008 chart queries will be called from frontend with PIN session header.

## Codebase Evidence
- PIN auth middleware: `middleware.GetAtlasUserID(ctx)` — returns "" if unauthorized. All existing resolvers follow the pattern: check userID == "" → return AuthError.
- MVP single-user: all operations scoped to default user extracted from context.
- No role-based authorization.
- When PIN disabled, endpoints accessible without auth.

## Proposed Details

### Auth
- All chart queries protected by WAVE-01 PIN auth middleware (GraphQL)
- Pattern: same as all existing resolvers — check middleware.GetAtlasUserID(ctx), return AuthError if empty
- AuthError type: use existing ChartErrorCode enum, or body/nutrition error codes

### Data Privacy
- Body weight values: WAVE-04 marks these as log-privacy-sensitive (NOT logged). Chart query may log query count but not weight values.
- Measurement values: same privacy treatment as WAVE-04 — NOT logged individually.
- Nutrition macro values (calories/protein/fat/carbs): WAVE-05 marks as non-sensitive — may be logged.
- No PII exposed in chart queries. Queries return aggregated numeric data only.
- Chart queries do not expose user identities or other user data.

### Audit
- Log markers record query type, date range, success/failure
- No detailed data values in logs

### Abuse/ Rate Limits
- Max date range of 52 weeks recommended to prevent expensive iteration (nutrition weekly averages)
- GraphQL query complexity — chart queries return bounded data (max 365 points for daily queries, 52 points for weekly queries)

### Compliance
- MVP single-user: no multi-tenant data isolation concerns beyond userID scoping
- No compliance framework applicable (no PII, no GDPR-sensitive data in chart responses)

## Risks And Rollback
- WAVE-04 body data (weight, measurements) is classified as sensitive — chart service must NOT log individual values
- Zero risk for existing security posture — all new code is additive queries with existing auth wrapper

## Questions Raised
- DQ-W06-007: Should chart queries log the returned data count or only the query request? (Proposed: log query + data point count, not values)

## Traceability Candidates
- docs/technical-verified/auth-security-compliance.md — TDEC-037 PIN auth
- apps/api/internal/atlas/middleware/ — auth middleware
- docs/prd-wave-details/waves/wave-04.md Security section (sensitive data logging rules)
- docs/prd-wave-details/waves/wave-05.md Security section (nutrition log privacy)