# Technical Brief

## Product Signal

Atlas is a self-hosted single-user web application for workout, nutrition, body measurement, and progress photo tracking with AI export. The verified product package (docs/product-verified) defines 11 MVP capability areas, 125 acceptance criteria, 31 edge cases, and 29 business rules.

Key technical drivers from product:
- Self-hosted via Docker Compose (single user, local filesystem)
- Optional PIN-based access (no user accounts)
- DailyLog aggregates all activity per date
- All entities owned by default user via userId
- Performance p95 targets defined for all UI and API operations
- AI export and backup/restore are user-triggered background operations
- No external API integrations in MVP

## Technical Scope

- Frontend: Next.js 15 (public web), Vite 5 + React 19 (web-admin)
- Backend: Go 1.25 with chi router, pgx, go-redis, gqlgen
- Database: PostgreSQL + Redis (session store)
- Infrastructure: Docker Compose
- Testing: vitest, Playwright, 100% coverage gate

## Constraints

- Single-user MVP with multi-user-ready schema (userId FK)
- PIN is the only access control mechanism
- Self-hosted only — no SaaS mode
- Manual AI analysis (copy-paste, not API)
- Local filesystem media storage (no cloud)
- Full backup/restore only (no incremental)
- No offline mode

## Assumptions

- Deployer comfortable with Docker CLI
- Browser with modern JavaScript
- AI analysis performed externally (ChatGPT or similar)
- Performance targets tested against expected dataset (5 years, 1500 daily logs, 30K sets)

## Readiness Summary

**questions-open.** 47 dev-blocking gaps identified. Readiness requires resolution of: API protocol, system architecture, session management, data model alignment, UI state contracts, deployment topology, and test strategy.