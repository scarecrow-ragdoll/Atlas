<!-- FILE: docs/prd-wave-details/appendix/decision-log.md -->
<!-- VERSION: 1.0.0 -->

# Decision Log

## Source Wave Gate
- Source wave: docs/prd-waves/waves/wave-08.md
- Source wave status: user-approved (2026-06-18)
- Source wave gate result: passed
- Gate check date: 2026-06-21

## User Wave Approvals
- WAVE-08 source wave: user-approved (2026-06-18)
- WAVE-08 detailed wave: awaiting user approval

## Scope Decisions

| ID | Decision | Source | Rationale |
|----|---------|--------|-----------|
| DDEC-W08-001 | GraphQL for AiReview CRUD | WAVE-07 pattern | All existing CRUD operations use GraphQL. No file download needed. |
| DDEC-W08-002 | Migration number 00093 | Codebase convention | Latest migration is 00092_ai_exports.sql (WAVE-07). Next is 00093. |
| DDEC-W08-003 | planned_actions as TEXT field | MVP constraints | PRD says "planned actions storage" — simple TEXT for MVP. Structured child table deferred. |
| DDEC-W08-004 | GraphQL-only (no REST endpoints) | AiReview has no files | WAVE-07 uses REST for ZIP download. AiReview has no file operations. |
| DDEC-W08-005 | User-scoped queries (no admin role) | Single-tenant design | All queries filtered by userId from session context. No admin/read-all needed. |

## Codebase Fit Decisions
- Follow WeekFlag triple pattern (model/repo/service/resolver)
- Follow WAVE-07 wiring patterns (resolver.go, main.go, gqlgen)
- No new config entries needed

## Deferrals
- planned_actions structured storage: deferred to post-MVP (DQ-W08-001)
- ai_response_text encryption at rest: deferred to post-MVP (no multi-tenant requirement)
- Max reviews per user limit: unbounded for MVP (user-managed)

## Rejected Assumptions
- No assumption of external AI API integration (explicitly excluded in source wave)
- No assumption of file storage (AiReview uses database-only storage)