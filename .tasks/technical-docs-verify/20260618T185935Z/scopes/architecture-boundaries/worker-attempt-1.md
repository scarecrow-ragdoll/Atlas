# Architecture-Boundaries Worker Report — Attempt 1

## Run Metadata

- **Run ID:** 20260618T185935Z
- **Role:** verify-technical-docs scoped worker: architecture-boundaries
- **Source:** docs/product-verified/
- **Available architecture file:** NONE — `docs/product-verified/architecture-and-boundaries.md` does not exist
- **Worker:** Main session (no subagent spawning available)

## Source Summary

### Key Resolved Decisions Affecting Architecture

| Decision | Impact |
|---|---|
| DEC-006 (Q-SCOPE-001) | Quality gates define test architecture expectations (coverage, e2e, integration) |
| DEC-007 (Q-SCOPE-002) | All entities owned via userId FK; default user at bootstrap; multi-user-ready data model |
| DEC-008 (Q-SCOPE-004) | p95 SLOs for UI, API, export, backup — affect deployment and resource planning |
| DEC-009 (Q-SCOPE-005) | DailyLog replaces WorkoutDay; CardioEntry requires dailyLogId; auto-create DailyLog on first activity |

### Architecture-Relevant Source Content

From `product-brief.md`:
- Self-hosted Docker deployment
- Single user per instance
- Performance SLOs: p95 targets for UI pages (1.0–2.0s), API mutations (300ms), queries (500ms–1.0s), AI export (5–20s), backup (15–30s)
- UX rule: operations >2s must show loading state

From `scope.md`:
- In scope: single-user self-hosted web app, optional PIN, full backup/restore
- Out of scope: SaaS, multi-user, mobile, cloud backup, Telegram bot
- Dependencies: Docker, Docker Compose, PostgreSQL, Redis, filesystem volume, Bun 1.1+, Node 22+, Go 1.25
- Assumptions: user has Docker knowledge, single user per instance, AI external, local filesystem media

From `domain-model.md`:
- 20 entities, all with userId FK except DefaultUser, WorkoutSet, BodyMeasurement, ProgressPhoto, NutritionTemplateItem, DailyNutritionOverrideItem, WeekFlag, AiExport, AiReview
- Entities without explicit userId inherit ownership through parent chain (e.g., WorkoutSet -> WorkoutExercise -> userId on WorkoutExercise)
- Cardinality and relationship invariants defined

From `actors-and-permissions.md`:
- Single actor (User), no role system
- All data belongs to single user
- 34 derived permissions covering all CRUD operations
- No multi-tenancy, no data sharing, no visibility rules

From `functional-spec.md`:
- 12 capability areas: Access Control, Dashboard, Exercise Library, Workout Diary, Cardio, Body Tracking, Nutrition, Charts, AI Export, AI Review, Backup/Restore, Settings
- No external API integrations in MVP
- AI analysis is manual copy-paste to ChatGPT
- Telegram bot library included in stack but explicitly out of MVP scope

## Architecture Analysis

### 1. System Context

**Current state:** The product-verified docs describe a self-hosted web application without a formal system context diagram. The implied boundaries are:

```
[Browser User] --HTTPS--> [Atlas Web App (Docker)] --SQL--> [PostgreSQL]
                                              --Cache--> [Redis]
                                              --Files--> [Filesystem Volume]
                                              --ZIP--> [Download]
```

**Missing:** No explicit system context diagram, no description of external systems (only AI/LLM is an external actor via manual copy-paste), no description of the Docker networking topology.

**TQ-ARCH-001:** System context diagram is missing. Architecture docs should define the system boundary, external actors (User, AI/LLM via manual export), and all data stores.

### 2. Component Boundaries

**Current state:** No component decomposition exists. The product docs imply a monolithic web app with 12 capability areas. The tech stack lists Bun, Node, Go — suggesting a polyglot monolith, but no component boundaries are defined.

**Implied components:**
- Web UI (frontend rendering)
- API layer (GraphQL or REST — not specified)
- Business logic / services layer
- Data access layer
- File/media storage abstraction
- Backup/export engine

**Missing:**
- No frontend/backend separation decision (SPA vs server-rendered)
- No API protocol decision (GraphQL vs REST vs tRPC)
- No service layer decomposition
- No component responsibility boundaries

**TQ-ARCH-002:** Component architecture is undefined. The product docs do not specify frontend framework, API protocol, or service decomposition boundaries. This blocks all downstream component-level decisions.

### 3. Ownership and Tenancy

**Current state:** Clearly defined:
- Single user per instance (DEC-007)
- All entities owned via userId FK with multi-user-ready data model
- Default user created at bootstrap
- No role system, no permission model beyond optional PIN
- All data belongs to the single user

**Implicit ownership chains:**
- WorkoutSet -> WorkoutExercise (userId) — inherited ownership
- BodyMeasurement -> BodyCheckIn (userId) — inherited ownership
- ProgressPhoto -> BodyCheckIn (userId) — inherited ownership
- NutritionTemplateItem -> NutritionTemplate (userId) — inherited ownership

**Assessment:**
- Ownership model is well-defined and clear
- "Multi-user-ready" means userId FK exists but no registration/login/session management
- The bootstrap default user approach needs an architecture decision: hardcoded UUID vs configurable vs generated

**TQ-ARCH-003:** Default user bootstrap mechanism is unspecified. Should the default user ID be hardcoded, generated from a deterministic seed, or configurable via environment variable?

### 4. Deployment Boundary

**Current state:**
- Docker Compose deployment
- PostgreSQL, Redis, filesystem volume
- Bun 1.1+, Node 22+, Go 1.25 runtime requirements
- Self-hosted by technically proficient user
- No SaaS, no cloud hosting

**Missing details:**
- No Docker Compose structure (services, networks, volumes)
- No environment variable inventory
- No resource requirements (CPU, RAM, disk) to meet p95 SLOs
- No backup/restore deployment considerations (volume persistence, migration)
- No SSL/TLS termination guidance
- No deployment verification instructions

**Performance implications from DEC-008:**
- AI export with photos (4 weeks) target: 20s p95 — implies significant I/O and CPU requirements
- Full backup with media: "best effort" — implies media size is unpredictable
- These targets need deployment guidelines (minimum RAM/CPU recommendations)

**TQ-ARCH-004:** Deployment architecture is underspecified. Docker Compose structure, environment variables, resource requirements, and SSL guidance are needed to meet the SLOs defined in DEC-008.

### 5. Service Boundaries

**Current state:**
- Monolithic deployment implied (single Docker service)
- No microservices, no internal service boundaries
- Go in stack alongside Bun/Node — unclear if Go is a separate service or embedded component

**Missing:**
- No decision on monolith vs modular monolith
- No guidance on Go vs Bun/Node responsibility boundaries
- No internal communication protocol between potential sub-components
- No background job processing architecture (AI export ZIP generation, backup/restore are long-running operations that could block the request thread)

**TQ-ARCH-005:** Service boundaries are undefined. Long-running operations (AI export, backup) may need background job processing, but no architecture decision exists. Go vs Bun/Node responsibility split is unclear.

### 6. Build-vs-Buy Boundaries

**Current state:**

| Component | Decision | Source |
|---|---|---|
| Database | PostgreSQL (buy/existing) | stack |
| Cache | Redis (buy/existing) | stack |
| Web framework | Bun/Node (build on) | stack |
| Background processing | Go (build on) | stack |
| AI analysis | External (manual copy-paste) | scope.md |
| Media storage | Local filesystem (build) | scope.md |
| Backup format | Custom ZIP format (build) | functional-spec.md |
| Telegram | Out of MVP (buy library, no integration) | scope.md |

**Assessment:**
- Build-vs-buy is mostly clear
- The Go dependency raises a question: what is Go used for? The product docs do not specify Go's role (CLI tooling, background jobs, export engine, media processing?)
- No evaluation of off-the-shelf fitness tracking solutions considered
- No build-vs-buy decision for charting library, UI framework, export ZIP generation

**TQ-ARCH-006:** Go's role in the architecture is undefined. The tech stack includes Go 1.25 but no product doc explains what Go component exists or what it owns.

### 7. Architecture Decisions Implied by Resolved Questions

| Resolved Question | Implied Architecture Decision | Missing Decision |
|---|---|---|
| Q-SCOPE-001 (DEC-006) | Quality gates require test infrastructure | No test architecture (e2e framework, integration test setup) |
| Q-SCOPE-002 (DEC-007) | Multi-user-ready data model with userId FK | No migration strategy from single-user to multi-user; no user management API |
| Q-SCOPE-004 (DEC-008) | Performance SLOs require adequate resource planning | No deployment resource recommendations |
| Q-SCOPE-005 (DEC-009) | DailyLog aggregates WorkoutExercise and CardioEntry | No decision on how auto-creation of DailyLog works (API design implication) |

## Key Gaps Summary

| ID | Gap | Severity | Impact |
|---|---|---|---|
| TQ-ARCH-001 | No system context diagram | High | Blocks system-level understanding |
| TQ-ARCH-002 | No component architecture (frontend/backend, API protocol) | High | Blocks all implementation planning |
| TQ-ARCH-003 | Default user bootstrap mechanism unspecified | Medium | Blocks bootstrap flow design |
| TQ-ARCH-004 | Deployment architecture underspecified | High | Blocks ops and SLO delivery |
| TQ-ARCH-005 | Service boundaries undefined (monolith vs modular, background jobs, Go role) | High | Blocks code organization |
| TQ-ARCH-006 | Go's role in architecture undefined | Medium | Blocks Go component planning |

## Recommended Next Scope Actions

1. **Create `architecture-and-boundaries.md`** in docs/product-verified covering: system context diagram, component architecture, deployment topology, service boundaries
2. **Resolve TQ-ARCH-002** before implementation: decide frontend framework, API protocol, and service decomposition
3. **Resolve TQ-ARCH-005** before implementation: decide monolith vs modular monolith, background job strategy
4. **Resolve TQ-ARCH-004** before ops planning: define Docker Compose structure, resource requirements, env vars
5. **Resolve TQ-ARCH-003** before bootstrap flow: define default user creation mechanism
6. **Resolve TQ-ARCH-006** before Go component work: define Go's role in the system

## Open Architecture Questions (TQ-ARCH-*)

| ID | Question | Product Source | Why It Matters |
|---|---|---|---|
| TQ-ARCH-001 | System context diagram is missing | All product docs | Cannot verify system boundaries without diagram |
| TQ-ARCH-002 | Component architecture undefined (frontend/backend split, API protocol, UI framework) | All product docs | Blocks all implementation planning |
| TQ-ARCH-003 | Default user bootstrap mechanism unspecified | scope.md §Assumptions | Affects first-run experience and data model initialization |
| TQ-ARCH-004 | Deployment architecture underspecified (Docker Compose, resources, env vars, SSL) | product-brief.md §Performance Targets | Blocks ops readiness and SLO delivery |
| TQ-ARCH-005 | Service boundaries undefined (monolith vs modular, background jobs, Go role) | technology stack | Blocks code organization and long-running operation handling |
| TQ-ARCH-006 | Go's role in architecture undefined | technology stack (Go 1.25) | Cannot plan Go component boundaries |