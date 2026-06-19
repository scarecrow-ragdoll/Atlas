# Wave Map Context

## Selected Backend Wave Boundary
WAVE-01 (Foundation): Infrastructure setup, database schema foundation, API skeleton extension, PIN guard, session management, settings service, CI/CD structure. Sets up all cross-cutting infrastructure that later waves (WAVE-02 through WAVE-09) depend on.

## Prior Backend Wave Fit
- WAVE-01 is the first wave. No prior backend waves exist.

## Future Backend Wave Fit
- WAVE-02 (Exercise Library): Depends on WAVE-01 for DB schema, API skeleton, media upload REST endpoint.
- WAVE-03 (Workout Diary): Depends on WAVE-01 DB + WAVE-02 API.
- WAVE-04 (Cardio/Body): Depends on WAVE-01 DB.
- WAVE-05 (Nutrition): Depends on WAVE-01 DB.
- Scope collision risk low: WAVE-01 establishes framework, not domain CRUD. Each later wave adds its own domain tables and resolvers.

## Frontend Pages Context
- PAGE-011 (Settings): Depends on settings API endpoints for PIN toggle, PIN change, AI context form, export preferences.
- PAGE-001 (Dashboard): Depends on API skeleton, data availability from WAVE-01 DB.
- Pin auth is a cross-cutting concern: PAGE-011 manages PIN, and PIN guard middleware protects all API and GraphQL calls.
- No frontend implementation is planned in this wave.

## Dependency Order
- WAVE-01 → WAVE-02 → WAVE-03 → WAVE-04/05 (parallel) → WAVE-06 → WAVE-07 → WAVE-08 → WAVE-09
- WAVE-01 has no predecessor dependencies.

## Scope Collision Check
- Q-PIN-001 (PIN rate limiting) tracked as DQ-W01-001 in wave-level open questions. Not blocking — severity is Medium, not wave-blocking or needs-owner-decision.
- The existing admin session/auth in apps/api is for admin operators (web-admin). The PIN auth for the fitness app is a separate concern — no collision.
- Settings service in WAVE-01 is for fitness-domain user settings (PIN, AI context, export preferences), not admin settings. No overlap with M-WEB-ADMIN admin_auth_service.