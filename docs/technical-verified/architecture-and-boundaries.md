# Architecture And Boundaries

## System Context

System context is not documented. No diagram or component map exists in product-verified sources.

Implied context from product and technology docs:
- **User**: single person accessing via browser
- **Atlas app**: Go backend + Next.js frontend + PostgreSQL + Redis
- **External AI**: ChatGPT or similar (manual copy-paste, no API)
- **Docker host**: self-hosted environment

## Components

No component architecture is defined. Technology stack implies:
- Go API server (chi router + gqlgen)
- PostgreSQL database
- Redis session store
- Web-admin (Vite + React 19 + shadcn)
- Public web (Next.js 15 + React 19)
- Docker Compose orchestration

Missing architecture decisions:
- GraphQL vs REST for primary API (TQ-API-001)
- Go backend role: GraphQL resolver or REST API provider (TQ-ARCH-006)
- File upload architecture: direct to filesystem or through API
- Background job execution model for export/import

## Ownership Boundaries

Clear: single-tenant, single-user. All entities owned via userId FK. No multi-tenancy, no sharing.

## Dependencies

| Dependency | Purpose | Boundary |
| --- | --- | --- |
| PostgreSQL | Primary data store | Internal, Docker-managed |
| Redis | Session store (PIN sessions) | Internal, Docker-managed |
| Docker Compose | Service orchestration | Deployment |
| Filesystem volume | Media storage | Internal, Docker volume |

## Architecture Questions

| ID | Question | Severity | Status |
| --- | --- | --- | --- |
| TQ-ARCH-001 | No system context diagram | dev-blocking | **resolved** (TDEC-014) |
| TQ-ARCH-002 | No component architecture or service boundaries | dev-blocking | **resolved** (TDEC-015) |
| TQ-ARCH-003 | Default user bootstrap mechanism undefined | dev-blocking | **resolved** (TDEC-016) |
| TQ-ARCH-004 | No deployment topology defined | dev-blocking | **resolved** (TDEC-017) |
| TQ-ARCH-005 | No API surface boundaries between public web and admin | dev-blocking | **resolved** (TDEC-018) |
| TQ-ARCH-006 | Go's role in the architecture undefined (GraphQL provider vs REST API) | dev-blocking | **resolved** (TDEC-019) |