# Decision Log

## Technical Decisions

| ID | Decision | Source Question | Rationale | Status |
| --- | --- | --- | --- | --- |
| TDEC-001 | Pagination: cursor-based, default 50, max 200 | TQ-API-004 | Owner decision 2026-06-18 | active |
| TDEC-002 | Idempotency: long-running ops required, CRUD optional clientMutationId | TQ-API-006 | Owner decision 2026-06-18 | active |
| TDEC-003 | Versioning: REST /api/v1, GraphQL evolutionary, export schema 1.0.0 | TQ-API-009 | Owner decision 2026-06-18 | active |
| TDEC-004 | Audit trail: required for PIN, media, export, backup operations | TQ-AUTH-004 | Owner decision 2026-06-18 | active |
| TDEC-005 | Data retention: no auto-deletion; entity-specific deletion rules | TQ-AUTH-009 | Owner decision 2026-06-18 | active |
| TDEC-006 | Canonical name: DailyLog replaces WorkoutDay everywhere | TQ-DATA-006 | Owner decision 2026-06-18 | active |
| TDEC-007 | BodyWeightEntry.source enum: MANUAL, DAILY_LOG, CHECK_IN, BACKUP_IMPORT, EXTERNAL_IMPORT | TQ-DATA-007 | Owner decision 2026-06-18 | active |
| TDEC-008 | File size limits: image 25MB, video 250MB, upload 300MB, AI export 500MB soft | TQ-INT-004 | Owner decision 2026-06-18 | active |
| TDEC-009 | Session continuity: jobId, progress polling, session-validated downloads | TQ-INT-005 | Owner decision 2026-06-18 | active |
| TDEC-010 | Resource estimates: min 1vCPU/1GB/10GB, recommended 2vCPU/4GB/50GB | TQ-OPS-006 | Owner decision 2026-06-18 | active |
| TDEC-011 | Log redaction: tested, specific no-log list | TQ-TEST-005 | Owner decision 2026-06-18 | active |
| TDEC-012 | Test isolation: isolated DB/media per type, Docker Compose e2e | TQ-TEST-006 | Owner decision 2026-06-18 | active |
| TDEC-013 | Performance tests: critical paths, CI smoke, release verification | TQ-TEST-007 | Owner decision 2026-06-18 | active |
| TDEC-014 | System context: self-hosted single-user web app with documented actors | TQ-ARCH-001 | Owner decision 2026-06-18 | active |
| TDEC-015 | Component architecture: modular monorepo with domain modules | TQ-ARCH-002 | Owner decision 2026-06-18 | active |
| TDEC-016 | DefaultUser: idempotent bootstrap, display name "Atlas User" | TQ-ARCH-003 | Owner decision 2026-06-18 | active |
| TDEC-017 | Deployment: Docker Compose with 5 services, 4 volumes | TQ-ARCH-004 | Owner decision 2026-06-18 | active |
| TDEC-018 | API surfaces: web-admin private, public web limited to static | TQ-ARCH-005 | Owner decision 2026-06-18 | active |
| TDEC-019 | Go role: one service handles both GraphQL and REST | TQ-ARCH-006 | Owner decision 2026-06-18 | active |
| TDEC-020 | userId on all user-owned entities including child tables | TQ-DATA-001 | Owner decision 2026-06-18 | active |
| TDEC-021 | Index strategy: indexes around query patterns and p95 targets | TQ-DATA-002 | Owner decision 2026-06-18 | active |
| TDEC-022 | Enums: BodyWeightSource, HeartRateZone, CardioType explicitly defined | TQ-DATA-003 | Owner decision 2026-06-18 | active |
| TDEC-023 | Enums: MeasurementType, MeasurementSide, WeekFlagType, MediaType, ProgressPhotoAngle | TQ-DATA-004 | Owner decision 2026-06-18 | active |
| TDEC-024 | Migration: SQL files as source of truth, versioned, immutable | TQ-DATA-005 | Owner decision 2026-06-18 | active |
| TDEC-025 | Import collision: clean restore or destructive replace, no merge | TQ-DATA-010 | Owner decision 2026-06-18 | active |
| TDEC-026 | Endpoint catalog: REST for auth/media/export/backup/jobs, GraphQL for CRUD | TQ-API-002 | Owner decision 2026-06-18 | active |
| TDEC-027 | Error format: common envelope with codes, fields, requestId | TQ-API-003 | Owner decision 2026-06-18 | active |
| TDEC-028 | File upload: REST multipart with purpose, entity refs, validation | TQ-API-007 | Owner decision 2026-06-18 | active |
| TDEC-029 | Session auth: cookie-based, HttpOnly, SameSite=Lax, session endpoint | TQ-API-008 | Owner decision 2026-06-18 | active |
| TDEC-030 | Chart queries: GraphQL with DateRangeInput, specific output shapes | TQ-API-010 | Owner decision 2026-06-18 | active |
| TDEC-031 | Backup import: two-step REST job (dry-run + commit), CLEAN_RESTORE/DESTRUCTIVE_REPLACE | TQ-API-011 | Owner decision 2026-06-18 | active |
| TDEC-032 | AI export: async REST job with jobId, progress, download URL | TQ-API-012 | Owner decision 2026-06-18 | active |
| TDEC-033 | PIN hash: Argon2id, 64MB memory, 3 iterations, 2 parallelism | TQ-AUTH-001 | Owner decision 2026-06-18 | active |
| TDEC-034 | Session TTL: idle 8h, absolute 7d, sliding renewal | TQ-AUTH-002 | Owner decision 2026-06-18 | active |
| TDEC-035 | Brute-force: 5 attempts → 5m lock, repeated → 30m lock | TQ-AUTH-003 | Owner decision 2026-06-18 | active |
| TDEC-036 | Session token: 32 random bytes, base64url, hash stored server-side | TQ-AUTH-005 | Owner decision 2026-06-18 | active |
| TDEC-037 | Media access: same policy as app, always via API, no public dirs | TQ-AUTH-006 | Owner decision 2026-06-18 | active |
| TDEC-038 | Redis failure: fail closed when PIN enabled, no insecure fallback | TQ-AUTH-007 | Owner decision 2026-06-18 | active |
| TDEC-039 | Bootstrap security: doesn't bypass PIN, idempotent, first-run PIN setup | TQ-AUTH-010 | Owner decision 2026-06-18 | active |
| TDEC-040 | PIN scope: global per instance, single-user in MVP | TQ-AUTH-011 | Owner decision 2026-06-18 | active |
| TDEC-041 | Sync/async: long ops async, CRUD sync, job status endpoint | TQ-INT-001 | Owner decision 2026-06-18 | active |
| TDEC-042 | Progress: job states QUEUED/RUNNING/SUCCEEDED/FAILED/CANCELLED | TQ-INT-002 | Owner decision 2026-06-18 | active |
| TDEC-043 | Import transaction: atomic at logical level, media staged, no silent partial | TQ-INT-003 | Owner decision 2026-06-18 | active |
| TDEC-044 | UI state machine: idle/loading/ready/empty/saving/error/unauthorized/offline/longRunning | TQ-CLIENT-001 | Owner decision 2026-06-18 | active |
| TDEC-045 | Page states: loading <300ms subtle, >1s visible, >2s explicit, progress for jobs | TQ-CLIENT-002..010 | Owner decision 2026-06-18 | active |
| TDEC-046 | Form validation: client+server rules, pessimistic writes | TQ-CLIENT-006, TQ-CLIENT-007 | Owner decision 2026-06-18 | active |
| TDEC-047 | Error display: field/form/page/toast/modal patterns | TQ-CLIENT-008 | Owner decision 2026-06-18 | active |
| TDEC-048 | Cache: React Query-style, stale times per entity type | TQ-CLIENT-009 | Owner decision 2026-06-18 | active |
| TDEC-049 | Navigation: dirty state tracking, beforeunload for unsaved changes | TQ-CLIENT-010 | Owner decision 2026-06-18 | active |
| TDEC-050 | Accessibility: WCAG 2.1 AA, keyboard nav, focus states, ARIA | TQ-CLIENT-011 | Owner decision 2026-06-18 | active |
| TDEC-051 | Localization: MVP Russian UI, i18n-ready code, English API values | TQ-CLIENT-012 | Owner decision 2026-06-18 | active |
| TDEC-052 | Config: env vars with defaults, no committed secrets | TQ-OPS-001 | Owner decision 2026-06-18 | active |
| TDEC-053 | Logging: zap, structured JSON, healthz/readyz endpoints | TQ-OPS-002 | Owner decision 2026-06-18 | active |
| TDEC-054 | Metrics: p95 instrumentation, Prometheus-style endpoint | TQ-OPS-003 | Owner decision 2026-06-18 | active |
| TDEC-055 | Backup ops: manual only, 24h temp retention, import dry-run as verification | TQ-OPS-005 | Owner decision 2026-06-18 | active |
| TDEC-056 | Test factories: deterministic, 5 fixture tiers, no real sensitive data | TQ-TEST-001 | Owner decision 2026-06-18 | active |
| TDEC-057 | E2E: full weekly workflow scenario, 18-step flow | TQ-TEST-002 | Owner decision 2026-06-18 | active |
| TDEC-058 | AI export tests: schema snapshot, deterministic, normalized timestamps | TQ-TEST-003 | Owner decision 2026-06-18 | active |
| TDEC-059 | Backup tests: manifest schema, version validation, missing media detection | TQ-TEST-004 | Owner decision 2026-06-18 | active |

## Deferrals

| Question | Deferred By | Rationale |
| --- | --- | --- |
| TQ-DATA-008 | product owner | Media lifecycle policy can be decided post-MVP |
| TQ-API-005 | product owner | Filtering contract can be defined per-endpoint |
| TQ-OPS-004 | product owner | Alerting/runbooks not needed for MVP |
| TQ-INT-006 | product owner | Rate limiting irrelevant for single-user |
| TQ-AUTH-008 | product owner | Backup identity validation future concern |

## Superseded Answers

None.

## Rejected Assumptions

None.