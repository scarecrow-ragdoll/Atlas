# Wave 01: Foundation

## Status

user-approved

## User Approval

user-approved (2026-06-18)

## Purpose

Establish core infrastructure: Docker Compose, Nx workspace, Go API skeleton, PostgreSQL, Redis, GraphQL, PIN guard.

## Outcome After Wave

- OUT-W01-001 Running Docker Compose stack
- OUT-W01-002 Basic Go API with GraphQL service
- OUT-W01-003 PostgreSQL database with table foundation
- OUT-W01-004 Session-based PIN auth guard
- OUT-W01-005 Test infrastructure

## Included Scope

- CAP-W01-001 Docker Compose setup
- CAP-W01-002 Nx workspace configuration
- CAP-W01-003 Go API skeleton
- CAP-W01-004 PostgreSQL core tables
- CAP-W01-005 Redis for sessions
- CAP-W01-006 GraphQL foundation
- CAP-W01-007 PIN guard middleware
- CAP-W01-008 Basic settings service
- CAP-W01-009 CI/CD ready structure

## Excluded Scope

- Exercise CRUD
- Workout diary
- Nutrition
- Charts
- AI export

## Dependencies

None (first wave)

## Surface Categories

backend, data, integrations, operations, security

## Risk Class

Medium - Infrastructure stability, session management

## Recommended Next Planning

$detail-prd-wave for WAVE-01 specifics

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Q-PIN-001 | 01 | security | Medium | None | PIN rate limiting implementation? | Security against brute force | docs/product/prd.md Section 24.1 | open | needs-owner-decision |

## Traceability

- docs/product/prd.md Sections 5, 7.2
- docs/product-verified/actors-and-permissions.md