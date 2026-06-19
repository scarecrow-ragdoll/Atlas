# WAVE-02 architecture-codebase-fit Review Attempt 1

## Verdict
needs-revision

## Sources Read
- planner-architecture-codebase-attempt-1.md
- planner-product-ac-attempt-1.md
- planner-data-integration-ops-attempt-1.md
- planner-sequencing-fit-attempt-1.md
- apps/api/cmd/server/main.go
- apps/api/internal/appconfig/config.go
- apps/api/internal/graph/resolver.go
- apps/api/internal/graph/schema.resolvers.go
- apps/api/internal/graph/admin_auth.resolvers.go
- apps/api/internal/graph/admin_auth_helpers.go
- apps/api/internal/service/admin_auth.go
- apps/api/internal/repository/postgres/user_repo.go
- apps/api/internal/repository/postgres/queries/users.sql
- apps/api/gqlgen.yml
- apps/api/sqlc.yaml
- docs/prd-wave-details/waves/wave-01.md

## Coverage Check
Architecture coverage is thorough: module structure, file locations, config changes, main.go wiring, codegen config, resolver pattern, repository pattern, service pattern all documented. The planner correctly identifies that gqlgen.yml and sqlc.yaml need no changes because they glob-match directories.

## Evidence Check
All codebase claims backed by existing file patterns. Module locations and pattern references are accurate.

## Codebase Fit Check
The planner correctly identifies:
- No existing exercise code exists (entirely new domain)
- WAVE-01 infrastructure dependencies (PIN auth, media scaffold, migration infra)
- Codegen autodiscovery via directory glob patterns
- Repository/Service/Middleware/Handler layering

### Issues Found

1. **Resolver dependency injection**: The planner says "Add ExerciseService to Resolver struct." This is correct, but the planner assumes Resolver dependency injection in main.go follows the current pattern. The current pattern creates dependencies at the top of main() and passes them to resolver. WAVE-02 will follow this. No issue, but the planner should clarify how ExerciseService depends on PIN auth service for guard logic — does the resolver call requirePinAuth directly, or does ExerciseService embed auth checks?

2. **PIN auth integration gap**: The planner says "requireAdmin → will be replaced by requirePinAuth or similar from WAVE-01." But WAVE-01's PIN auth guard pattern is not yet defined in code. The planner assumes WAVE-01 will provide a `requirePinAuth(ctx) error` function in the graph package, analogous to `requireAdmin`. This is a valid assumption but should be made explicit as a dependency contract.

3. **Missing file: GraphQL schema for common types**: Exercise CRUD needs `ValidationError`, `AuthError`, `NotFoundError` types in the fitness GraphQL schema. If WAVE-01 already defines these for fitness domain (in fitness.graphql or common.graphql), WAVE-02 reuses them. If not, WAVE-02 must define or import them. This should be clarified.

4. **Migration file numbering**: Planner suggests 00080_exercises.sql and 00081_exercise_media.sql. WAVE-01 may create migrations in the 00080+ range for its own tables. Need to verify WAVE-01's migration plan and adjust numbering accordingly. The planner correctly notes this is a dependency on WAVE-01.

## Other-Wave Fit Check
WAVE-01 provides the dependency layer. WAVE-03 needs ListAll exercises. The allExercises query is correctly identified.

## Acceptance Criteria Check
Not applicable to this reviewer perspective directly, but AC-W02-021 (allExercises for WAVE-03) and AC-W02-022 (query by ID after soft-delete) are architecturally validated.

## Exit Criteria Check
EC-W02-004 (media scaffold extension) requires clarification on how WAVE-01's handler is extended.

## Verification Check
No architecture-specific concerns with the test plan.

## Question Ledger Check
DQ-W02-003 (file storage path pattern) is an architecture-level question that needs resolution before coding starts. DQ-W02-004 (physical file deletion) is a data lifecycle decision with architecture implications.

## Unsupported Or Invented Claims
The planner invents a `requirePinAuth` function name and assumes it will exist in the graph package from WAVE-01. This is a reasonable assumption but should be flagged as a dependency contract that must be verified against WAVE-01 output.

## Required Revisions
1. **Explicitly state dependency contract**: Add a "WAVE-01 dependency contract" section listing what exactly WAVE-01 must provide for WAVE-02 (requirePinAuth function, media REST base path config, GraphQL common types).
2. **Clarify common GraphQL types**: Specify whether ValidationError/AuthError/NotFoundError types are shared or domain-specific.
3. **Add ExerciseMedia route registration**: The planner mentions a new route group but doesn't specify how it's registered in main.go (PIN-protected group beside existing admin/public groups).
4. **Verify migration numbering**: Add a note that migration file numbers must be confirmed after WAVE-01's migration plan is final.

## Approval Notes
Good architectural coverage. The 4 revision items are clarifying — no fundamental architecture issues. After revisions, will approve.