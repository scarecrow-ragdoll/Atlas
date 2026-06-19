# API-Contracts Scope Orchestrator

## Run Metadata

- **Run ID:** 20260618T185935Z
- **Source:** docs/product-verified (functional-spec.md, domain-model.md, actors-and-permissions.md, scope.md)
- **Source Delta:** All 4 blocking product questions resolved (DEC-006 through DEC-009)
- **Key delta for this scope:**
  - DEC-007 (userId FK on all entities) — affects API request/response schemas and resource identifiers
  - DEC-009 (DailyLog replaces WorkoutDay) — affects workout/cardio API surfaces and invariants

## Worker Assignment

- **Role:** verify-technical-docs scoped worker: api-contracts
- **Attempt:** 1
- **Report path:** worker-attempt-1.md
- **Focus:** API surfaces, GraphQL/REST/RPC, request/response, error formats, validation, pagination, versioning, idempotency, compatibility

## Review Gate

- **Reviewer role:** verify-technical-docs scoped reviewer: api-contracts
- **Verdicts:** approved | needs-revision | blocked

## Budget

- REVIEW_BUDGET: 3 cycles
- INTERRUPTION_RETRY_BUDGET: 3 relaunches

## Output Artifacts

| Artifact | Path |
|---|---|
| orchestrator.md | This file |
| worker-attempt-1.md | Scope worker report |
| review-attempt-1.md | Scope review verdict |
| scope-status.md | Final scope status |
| question-ledger.md | API-contracts open questions |