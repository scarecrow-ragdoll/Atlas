# Roles-Permissions-Reviewer Orchestrator

## Run Metadata

- **Run ID**: 20260618T185935Z
- **Scope**: roles-permissions-reviewer
- **Source**: docs/product/prd.md
- **Source Delta**: none
- **Review Budget**: 3 cycles
- **Interruption Retry Budget**: 3 cycles

## Scope Focus

Extracting and deriving actors, roles, permissions, ownership rules, approval authority, visibility rules, responsibility boundaries, and the permissions matrix. Every derived role or permission must cite source behavior, derivation rationale, and confidence.

## Artifacts

| Artifact | Path | Status |
|---|---|---|
| Orchestrator | `.tasks/product-docs-verify/20260618T185935Z/scopes/roles-permissions-reviewer/orchestrator.md` | Created |
| Worker attempt 1 | `.tasks/product-docs-verify/20260618T185935Z/scopes/roles-permissions-reviewer/worker-attempt-1.md` | Pending |
| Review attempt 1 | `.tasks/product-docs-verify/20260618T185935Z/scopes/roles-permissions-reviewer/review-attempt-1.md` | Pending |
| Scope status | `.tasks/product-docs-verify/20260618T185935Z/scopes/roles-permissions-reviewer/scope-status.md` | Pending |
| Question ledger | `.tasks/product-docs-verify/20260618T185935Z/scopes/roles-permissions-reviewer/question-ledger.md` | Pending |

## Worker Prompts

### Worker Attempt 1 Prompt

```
Use the verify-product-docs scoped worker role: roles-permissions-reviewer.

Input folder: docs/product
Run id: 20260618T185935Z
Attempt: 1
Report path: .tasks/product-docs-verify/20260618T185935Z/scopes/roles-permissions-reviewer/worker-attempt-1.md

Role focus:
Actors, roles, permissions, ownership rules, approval authority, visibility rules, responsibility boundaries, permissions matrix.

Available source files:
- docs/product/prd.md

Source delta: none.

Write the worker report only. Do not edit docs/product-verified. Do not edit staging skeletons. Record missing information as open questions. Derive roles/permissions from described actors, actions, ownership, approvals, visibility, and denied/allowed flows; include source reference, derivation rationale, and confidence. Do not invent unrelated behavior, API details, integration contracts, or implementation contracts. Use the exact report path above.
```

## Review Iterations

| Attempt | Worker Status | Reviewer Status | Verdict |
|---|---|---|---|
| 1 | Pending | Pending | Pending |