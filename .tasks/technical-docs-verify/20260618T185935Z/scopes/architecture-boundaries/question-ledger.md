# Architecture-Boundaries Question Ledger

## Run Metadata

- **Run ID:** 20260618T185935Z
- **Scope:** architecture-boundaries
- **Source:** docs/product-verified

## Open Architecture Questions (TQ-ARCH-*)

| ID | Question | Product Source | Severity | Why It Matters | Recommended Action |
|---|---|---|---|---|---|
| TQ-ARCH-001 | System context diagram is missing | All product docs | High | Cannot verify system boundaries without diagram | Create system context diagram in architecture-and-boundaries.md |
| TQ-ARCH-002 | Component architecture undefined (frontend/backend split, API protocol, UI framework) | All product docs | High | Blocks all implementation planning | Decide frontend framework, API protocol, and service decomposition before implementation |
| TQ-ARCH-003 | Default user bootstrap mechanism unspecified | scope.md §Assumptions | Medium | Affects first-run experience and data model initialization | Define default user creation mechanism (hardcoded UUID vs env var vs generated) |
| TQ-ARCH-004 | Deployment architecture underspecified (Docker Compose, resources, env vars, SSL) | product-brief.md §Performance Targets | High | Blocks ops readiness and SLO delivery | Define Docker Compose structure, resource requirements, env var inventory, SSL guidance |
| TQ-ARCH-005 | Service boundaries undefined (monolith vs modular, background jobs, Go role) | technology stack | High | Blocks code organization and long-running operation handling | Decide monolith vs modular monolith, background job strategy (AI export, backup) |
| TQ-ARCH-006 | Go's role in architecture undefined | technology stack (Go 1.25) | Medium | Cannot plan Go component boundaries | Define Go component responsibilities (CLI tooling, background jobs, export engine?) |

## Review-Identified Minor Gaps

| ID | Gap | Source | Severity | Recommendation |
|---|---|---|---|---|
| TQ-ARCH-007 | PIN guard architecture not analyzed (session boundary, middleware layer) | actors-and-permissions.md | Low | Address when architecture-and-boundaries.md is created |
| TQ-ARCH-008 | Media serving architecture not analyzed (photo access control via PIN) | functional-spec.md §14.3 | Low | Address when architecture-and-boundaries.md is created |

## Question Status Summary

| Status | Count |
|---|---|
| Open (High) | 4 |
| Open (Medium) | 2 |
| Open (Low) | 2 |
| Resolved | 0 |
| Total | 8 |