# API-Contracts Worker Attempt 1

## Sources Read

- docs/product-verified/functional-spec.md (122 lines, 19 capability areas)
- docs/product-verified/domain-model.md (135 lines, 20 entities, relationships, invariants)
- docs/product-verified/actors-and-permissions.md (75 lines, single-user model)
- docs/product-verified/scope.md (72 lines, in/out scope, assumptions)
- Source delta (DEC-006 through DEC-009)

## Source Delta Reviewed

- **DEC-007** (userId FK): Every entity carries a userId. API must support userId scoping at the resource level even though only one user exists. No explicit API-auth mechanism documented.
- **DEC-009** (DailyLog replaces WorkoutDay): Cardio requires dailyLogId. DailyLog is the daily aggregate. The API surface for workout vs cardio must represent this through DailyLog as the primary date-anchored container.

## Product Signals

1. **Single-user self-hosted app** — no multi-tenancy, no registration, no public API
2. **19 entities** with CRUD operations implicitly required for all
3. **Optional PIN guard** — session-based, no API tokens or OAuth
4. **File uploads** — exercise media (images/video), progress photos, backup ZIP
5. **Backup ZIP import** — requires file upload + dry-run + confirm restore
6. **AI export ZIP** — server-generated ZIP download
7. **No external API integrations** — AI is manual copy-paste, Telegram/Apple Health out of scope
8. **All data belongs to single user** — but userId FK exists for future multi-user readiness
9. **Validation rules** — partially implied (numeric fields, positive values, required optional)
10. **Charts** — implied query surfaces for aggregated data (training progress, body measurements, nutrition averages)

## Technical Facts

### Missing API Contract Surface (Absence)
- No REST, GraphQL, gRPC, or RPC protocol decision documented.
- No endpoint list, no method signatures, no URL patterns, no GraphQL schema.
- No request/response shapes for any entity operation.
- No error format (status codes, error body structure, field-level errors).
- No validation mapping (how domain validation errors translate to API errors).
- No pagination contract for list operations (cursor vs offset, default page size, max page size).
- No filtering contract (date ranges, exercise names, measurement types, flag types).
- No sorting contract (order by field, direction defaults).
- No idempotency contract (retry-safety for set creation, backup import).
- No versioning strategy (URL prefix, header, GraphQL schema version, or no versioning).
- No compatibility policy (backward-compatible changes, breaking change process).
- No API documentation contract (OpenAPI, SDL, or other format).
- No rate limiting or request size limits documented.
- No health/status endpoint documented.

### Protocol Implications From Product Behavior
1. **CRUD-heavy**: 19 entities, most with create/read/update/delete/list — REST is the natural fit. GraphQL would reduce endpoint count but adds complexity for a single-developer MVP.
2. **File operations**: Two upload paths (exercise media, progress photos), two ZIP generation paths (AI export, backup), one ZIP upload (backup import). REST with multipart/form-data is the standard approach.
3. **Charts**: Aggregate query endpoints needed (training progress per exercise, body measurements over time, nutrition averages). These are not simple CRUD — they require server-side aggregation.
4. **Backup import flow**: Upload → validate → dry-run → show summary → user confirms → restore. This is a multi-step process that needs a stateful operation or a multi-endpoint flow.
5. **PIN guard**: Session-based auth. API endpoints must validate session cookie. No Bearer tokens, no JWT documented.

### Entity-Level API Signals

| Entity | Implied Operations | Notes |
|---|---|---|
| Settings | R, U (single record) | Upsert pattern likely |
| UserProfile | R, U | Single profile |
| Exercise | CRUD + list + activate/deactivate | Filterable by isActive |
| ExerciseMedia | C, D (per exercise) | File upload |
| DailyLog | R by date, U (partial) | Unique per (userId, date) |
| WorkoutExercise | CRUD within DailyLog | Order field |
| WorkoutSet | CRUD within WorkoutExercise | Ordered by setNumber |
| CardioEntry | CRUD within DailyLog | Attached to DailyLog |
| BodyWeightEntry | C, list | Standalone + embedded in DailyLog |
| BodyCheckIn | CRUD | With measurements and photos |
| BodyMeasurement | CRUD within check-in | 10 types, left/right for paired |
| ProgressPhoto | C, D within check-in | File upload |
| NutritionProduct | CRUD + list | Catalog filterable by name |
| NutritionTemplate | CRUD, single active rule | Week-start-date-based |
| NutritionTemplateItem | CRUD within template | |
| DailyNutritionOverride | C, R by date, D | One override per date |
| DailyNutritionOverrideItem | CRUD within override | operation: add/subtract/replace |
| WeekFlag | CRUD by week | Per-week flags |
| AiExport | C (generate ZIP), R | Multi-step: build prompt → generate → download |
| AiReview | CRUD | Manual response storage |
| DefaultUser | Created at bootstrap (no user-facing API) | System entity |

### DEC-007 (userId FK) API Implications
- Every create/read/update/list must scope to the authenticated user's userId.
- Even though it's single-user, the API should accept and validate userId from session, not from request body.
- List operations must filter by userId automatically.
- This is a server-enforced rule, not a client responsibility.

### DEC-009 (DailyLog Replaces WorkoutDay) API Implications
- DailyLog is the primary date anchor. Workout exercises, cardio, body weight, and notes all live under it.
- CardioEntry requires dailyLogId — cardio cannot exist without a DailyLog. The system must auto-create DailyLog when cardio is saved for a new date.
- API surface must expose DailyLog as the top-level date resource, not separate "workout day" and "cardio day" endpoints.
- Backward compatibility: if any earlier plan documented WorkoutDay endpoints, the API contract must unify them.

### Validation Gaps (API-Level)
- Working weight: numeric positive
- Sets: weight >= 0 numeric, reps positive int
- Duration: positive int (minutes)
- Photos per check-in: 2-4 (requirement vs recommendation unclear)
- Nutrition values: numeric per 100g
- Amount grams: positive numeric
- Backup ZIP: must contain manifest.json and data.json
- No validation error format defined (field-level vs general errors)
- No constraint violation mapping (unique constraint on DailyLog(userId, date), FK violations for exerciseId, etc.)

### Error Contract Gaps
- HTTP status codes not specified for any scenario (200, 201, 204, 400, 404, 409, 422, 500)
- Error response body format not defined (JSON structure, error codes, messages, field-level errors)
- Validation error shape not defined (what fields failed, what rules were violated)
- Not-found behavior not defined (return 404 vs empty list vs null)
- Conflict behavior not defined (duplicate DailyLog for date, duplicate product name)
- Server error format not defined (what gets exposed to client vs logged)

### Pagination And Filtering Gaps
- No pagination for any list endpoint (exercises, weight entries, check-ins, nutrition products, ai reviews)
- No filtering contract (by date range, by name, by isActive, by flag type, by measurement type)
- No sorting contract (by date, by name, by createdAt)
- No search contract (by exercise name, by product name)
- Chart query surfaces undefined (period filters, exercise selection, measurement selection)

### Idempotency Gaps
- No idempotency key strategy for any mutation endpoint
- Backup import: repeated submission risk without idempotency
- Set creation: duplicate on retry without idempotency
- Exercise creation: duplicate entries on network retry

### Versioning And Compatibility Gaps
- No API versioning strategy
- No backward-compatibility policy
- Backup ZIP has schema version in manifest — this is data-level versioning, not API-level
- No deprecation policy for future changes

## Technical Gaps

1. TGAP-API-001: Missing protocol decision (REST vs GraphQL vs other)
2. TGAP-API-002: Missing endpoint catalog for all 19 entities
3. TGAP-API-003: Missing request/response schemas
4. TGAP-API-004: Missing error format and status code contract
5. TGAP-API-005: Missing validation mapping
6. TGAP-API-006: Missing pagination/filtering/sorting contract
7. TGAP-API-007: Missing file upload/download contract
8. TGAP-API-008: Missing idempotency strategy
9. TGAP-API-009: Missing versioning and compatibility policy
10. TGAP-API-010: Missing chart/aggregation query surfaces
11. TGAP-API-011: Missing backup import multi-step flow contract
12. TGAP-API-012: Missing API auth/session validation contract
13. TGAP-API-013: Missing health/status endpoint

## Questions Raised

| ID | Severity | Parent | Question | Why It Matters | Needed Artifact Or Decision | Status |
|---|---|---|---|---|---|---|
| TQ-API-001 | dev-blocking | none | Which API protocol? | Every endpoint, client, error format, and codegen depends on this decision. | Owner decision: REST (typical for CRUD + file ops + single dev) vs GraphQL (schema-driven, single endpoint, but more complex for file ops). REST recommended for this scope given CRUD-heavy, file uploads, single developer. | open |
| TQ-API-002 | dev-blocking | none | What is the API endpoint catalog? | Implementation cannot start without knowing which endpoints exist for 19 entities, chart queries, backup flow, and settings. | Endpoint list with methods, URL patterns, and purpose per entity/flow. | open |
| TQ-API-003 | dev-blocking | none | What are the request/response schemas? | Client-server contract requires defined JSON shapes for every operation. | OpenAPI spec or equivalent schema definitions. | open |
| TQ-API-004 | dev-blocking | none | What is the API error format? | Clients need to handle errors uniformly — status codes, error body, field-level validation errors, not-found, conflict, server error. | Error response schema and HTTP status code map. | open |
| TQ-API-005 | dev-blocking | none | What is the validation mapping? | Domain validation rules must translate into API error responses with consistent structure. | Validation-to-error mapping and field error schema. | open |
| TQ-API-006 | needs-owner-decision | none | What pagination/filtering/sorting strategy? | List endpoints need consistent pagination (cursor vs offset), filter parameters, and sort defaults. | Pagination contract, filter parameter list, sort field enum per entity. | open |
| TQ-API-007 | needs-owner-decision | none | How are file uploads/downloads handled? | Exercise media, progress photos, backup ZIP, AI export ZIP all need upload/download contracts (multipart, content-type, size limits, streaming). | File upload/download endpoint contract with size limits and content-type constraints. | open |
| TQ-API-008 | watchlist | none | Should mutations be idempotent? | Retry-safety for set creation, exercise creation, and backup import prevents duplicate data on network retries. | Idempotency key strategy or explicit no-idempotency decision. | open |
| TQ-API-009 | needs-owner-decision | none | What is the API versioning/compatibility policy? | Future changes need a contract for backward compatibility and breaking change handling. | Versioning strategy (URL prefix / header / no versioning) and compatibility rules. | open |
| TQ-API-010 | dev-blocking | none | What is the chart/aggregation query contract? | Training progress, body measurements, nutrition averages need server-side aggregation endpoints. Without defined shapes, client chart implementation is blocked. | Chart query request/response schemas and aggregation logic contract. | open |
| TQ-API-011 | dev-blocking | none | What is the backup import multi-step flow? | Upload → validate → dry-run → summary → confirm → restore needs state management. Endpoints, flow states, and error handling undefined. | Backup import flow contract with endpoint sequence and state management. | open |
| TQ-API-012 | dev-blocking | none | How does PIN session auth apply to API? | API endpoints must validate session. Is session checked via cookie, header, or both? What endpoints are public vs protected? | Session validation contract for API endpoints. | open |
| TQ-API-013 | watchlist | none | Should there be a health/status endpoint? | Docker deployment needs liveness/readiness probes. Single-user app may also want status endpoint for frontend bootstrap. | Health endpoint contract or explicit exclusion. | open |

## Answer Effects

No prior answered questions exist for this scope — first run.

## Risks

1. **Protocol choice latency**: If REST is chosen, endpoint catalog is large but well-understood. If GraphQL is chosen, schema-first approach reduces endpoint count but adds query complexity and file-upload challenges. Either decision is viable but must be made before endpoint design.
2. **Scope creep on chart APIs**: Chart aggregation may expand into a significant backend surface. Without defined query contracts, implementation time is unbounded.
3. **Backup flow state management**: Multi-step import (upload → validate → dry-run → confirm → restore) requires either server-side flow state (session/temp table) or client-driven sequential endpoints. The former is more robust for large backups; the latter is simpler.

## Suggested Decisions

1. **Protocol: REST** — CRUD-heavy, file uploads, single developer, existing ecosystem tools (OpenAPI, swagger, Postman). GraphQL adds schema-publishing overhead with limited benefit for single-user app.
2. **Endpoint URL pattern**: `/api/v1/` prefix for versioning readiness, even before API versioning is formally needed.
3. **Error format**: Standard JSON error body with HTTP status codes. Use `{ "error": { "code": "...", "message": "...", "details": [...] } }` shape. Field-level validation in `details` array.
4. **Pagination**: Cursor-based for list endpoints that grow indefinitely (exercise list, weight entries, check-ins). Offset-based for small/lookup lists (nutrition products, settings).
5. **Idempotency**: Skip idempotency keys for MVP. Accept risk of duplicate on retry. Document known gaps.
6. **Session**: Validate PIN session via cookie on all API endpoints except health. No Bearer token for MVP.
7. **Health endpoint**: Add `/api/v1/health` returning `{ "status": "ok" }` for Docker probes.

## Traceability Candidates

| Source Element | Artifact Candidate |
|---|---|
| PRD §10-20 (19 entities) | REST endpoints for each entity |
| PRD §10.5 (sets) | WorkoutSet CRUD under DailyLog |
| PRD §11.3 (media upload) | Multipart file upload endpoint |
| PRD §13.2-13.4 (check-ins) | BodyCheckIn CRUD with nested measurements/photos |
| PRD §15.3-15.5 (nutrition) | Nutrition template + override endpoints |
| PRD §16 (charts) | Aggregation query endpoints |
| PRD §17-18 (AI export) | Generate → download ZIP flow |
| PRD §20 (backup) | Upload → validate → dry-run → import flow |
| PRD §7.2 (PIN) | Session validation middleware for API |
| DEC-007 (userId FK) | Server-enforced userId scoping on all endpoints |
| DEC-009 (DailyLog) | DailyLog as primary date resource unifying workout, cardio, body weight, notes |