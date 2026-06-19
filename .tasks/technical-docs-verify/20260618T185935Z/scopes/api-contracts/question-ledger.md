# API-Contracts Question Ledger

| ID | Scope | Severity | Parent | Question | Why It Matters | Needed Artifact Or Decision | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| TQ-API-001 | api-contracts | dev-blocking | none | Which API protocol? | Every endpoint, client, error format, and codegen depends on this decision. | Owner decision: REST vs GraphQL | worker-attempt-1.md | open | TBD |
| TQ-API-002 | api-contracts | dev-blocking | none | What is the API endpoint catalog? | Implementation cannot start without knowing which endpoints exist. | Endpoint list with methods, URL patterns, and purpose | worker-attempt-1.md | open | TBD |
| TQ-API-003 | api-contracts | dev-blocking | none | What are request/response schemas? | Client-server contract requires defined JSON shapes for every operation. | OpenAPI spec or equivalent schema definitions | worker-attempt-1.md | open | TBD |
| TQ-API-004 | api-contracts | dev-blocking | none | What is the API error format? | Clients need to handle errors uniformly. | Error response schema and HTTP status code map | worker-attempt-1.md | open | TBD |
| TQ-API-005 | api-contracts | dev-blocking | none | What is the validation mapping? | Domain validation rules must translate into API error responses. | Validation-to-error mapping and field error schema | worker-attempt-1.md | open | TBD |
| TQ-API-006 | api-contracts | needs-owner-decision | none | What pagination/filtering/sorting strategy? | List endpoints need consistent pagination, filter parameters, and sort defaults. | Pagination contract, filter parameter list, sort field enums | worker-attempt-1.md | open | TBD |
| TQ-API-007 | api-contracts | needs-owner-decision | none | How are file uploads/downloads handled? | Exercise media, progress photos, ZIPs need upload/download contracts. | File upload/download contract with size limits and content-type | worker-attempt-1.md | open | TBD |
| TQ-API-008 | api-contracts | watchlist | none | Should mutations be idempotent? | Retry-safety prevents duplicate data on network retries. | Idempotency key strategy or explicit no-idempotency decision | worker-attempt-1.md | open | TBD |
| TQ-API-009 | api-contracts | needs-owner-decision | none | What is the API versioning/compatibility policy? | Future changes need backward compatibility and breaking change handling. | Versioning strategy and compatibility rules | worker-attempt-1.md | open | TBD |
| TQ-API-010 | api-contracts | dev-blocking | none | What is the chart/aggregation query contract? | Chart implementation requires defined aggregation endpoints. | Chart query request/response schemas | worker-attempt-1.md | open | TBD |
| TQ-API-011 | api-contracts | dev-blocking | none | What is the backup import multi-step flow? | Upload→validate→dry-run→confirm→restore needs state management. | Backup import flow contract with endpoint sequence | worker-attempt-1.md | open | TBD |
| TQ-API-012 | api-contracts | dev-blocking | none | How does PIN session auth apply to API? | API endpoints must validate session. | Session validation contract for API endpoints | worker-attempt-1.md | open | TBD |
| TQ-API-013 | api-contracts | watchlist | none | Should there be a health/status endpoint? | Docker deployment needs liveness/readiness probes. | Health endpoint contract or explicit exclusion | worker-attempt-1.md | open | TBD |

## Summary

- Total questions: 13
- dev-blocking: 8
- needs-owner-decision: 3
- watchlist: 2
- deferred: 0
- resolved: 0