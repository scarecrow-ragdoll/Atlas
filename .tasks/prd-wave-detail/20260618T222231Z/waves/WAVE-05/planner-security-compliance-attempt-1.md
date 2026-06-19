# WAVE-05 Security-Compliance Planner Attempt 1

## Sources Read
- docs/technical-verified/auth-security-compliance.md
- docs/product-verified/actors-and-permissions.md
- docs/product-verified/business-rules.md (RULE-022, RULE-023, RULE-024)
- apps/api/internal/atlas/middleware/pin_auth.go
- apps/api/internal/atlas/middleware/user_context.go
- apps/api/cmd/server/main.go (Atlas route groups)
- docs/prd-wave-details/waves/wave-04.md (Security section reference)

## Selected Backend Wave Boundary
All WAVE-05 operations are user-scoped (single default user per MVP). All operations go through the Atlas PIN-protected `/graphql/atlas` endpoint. No REST endpoints. No binary uploads. No sensitive PII in nutrition domain (food products and macro values are not health-data sensitive in the same way as body weight/photos).

## Neighboring Backend Wave Fit
Same PIN auth pattern as WAVE-02/03/04. No additional auth boundary needed. No overlap with admin auth.

## Frontend Pages Context
PAGE-007 uses the same PIN session cookie as all other Atlas frontend pages. No special auth handling needed.

## Codebase Evidence
- AtlasGraphQL handler is registered under the `atlasGuarded` chi router group which applies `atlasMiddleware.AtlasPinGuard`
- `atlasMiddleware.AtlasUserContext` injects `atlasUserID` (single default user) into request context
- When PIN is disabled, PIN guard is still in middleware chain but session is not required

## Proposed Details

### Auth and Authorization
- All GraphQL operations go through existing `/graphql/atlas` endpoint already PIN-protected by `atlasMiddleware.AtlasPinGuard`
- No additional auth middleware needed
- Service layer extracts userID from context via `middleware.GetAtlasUserID(ctx)` (same pattern as settings resolver)
- MVP single-user: no cross-user authorization needed. All data is scoped to the default user.

### Privacy
- Nutrition product names and notes: treated as low sensitivity. They are user-created food names, not health data.
- Nutritional values (calories, protein, fat, carbs): low sensitivity — can be logged.
- Notes field: user-entered free text. Log with caution — log entity IDs only, not notes content.
- No user PII in this domain.

### Audit
- No formal audit trail needed for MVP nutrition data
- Standard created_at/updated_at timestamps on all entities
- Hard delete for templates and overrides (no isActive flag) means no recovery from deletion

### Abuse and Rate Limiting
- MVP single-user: no rate limiting needed for nutrition operations
- Deferred: rate limiting on all Atlas endpoints at WAVE-01 level if needed

### Secrets
- No secrets in nutrition domain. No external API keys.

## Risks And Rollback
- Hard delete for templates and overrides: no undo. User should be warned before delete.
- Soft-delete for products: recoverable via isActive flag re-enable.

## Questions Raised
- DQ-W05-007: Should NutritionProduct soft-deleted items be recoverable via API (e.g., listInactive query) or only via DB admin? Recommended: admin-only (direct DB) for MVP simplicity.

## Traceability Candidates
- docs/technical-verified/auth-security-compliance.md → PIN auth for all endpoints
- apps/api/internal/atlas/middleware/pin_auth.go → auth guard implementation
- apps/api/cmd/server/main.go → route group protection