# Operations-Observability Worker Attempt 1

<!-- FILE: .tasks/technical-docs-verify/20260618T185935Z/scopes/operations-observability/worker-attempt-1.md -->

## Sources Read

| Source | Key Content |
| --- | --- |
| docs/product-verified/product-brief.md | DEC-008 p95 SLOs (API, export, backup), quality gates (no sensitive data in logs), backup/restore success metrics |
| docs/product-verified/scope.md | Docker + Docker Compose, PostgreSQL, Redis, Bun 1.1+/Node 22+/Go 1.25 |
| docs/product-verified/domain-model.md | 20 entities, data model invariants |
| docs/product-verified/features/backup-and-restore.md | Full backup ZIP (manifest.json, data.json, media/), dry-run validation, schema version, import workflow |

## Source Delta Reviewed

DEC-008 adds p95 performance targets for:
- API operations (simple mutation 300ms, daily log 500ms, exercise history 700ms, chart data 1.0s)
- AI Export (4w without photos 5s, 4w with photos 20s, 12mo without photos 15s, 12mo with photos best-effort)
- Backup (db-only 15s, import 30s, dry-run 15s; with-media best-effort)
- UX rule: operations >2s must show loading state

No SLO formulation (error budget, measurement window, SLI definition) is provided — only raw p95 targets.

## Product Signals

1. Single-user self-hosted — operations burden entirely on the user.
2. Optional PIN-based access (no real auth infra).
3. Backup/restore is critical for data ownership story.
4. AI export includes optional photos — media I/O performance is state-dependent.
5. Weekly analysis cycle drives recurring data volume growth.

## Technical Facts

| Fact | Source |
| --- | --- |
| Container runtime: Docker + Docker Compose | scope.md §60-62 |
| Storage: PostgreSQL, Redis, file system volume for media | scope.md §61 |
| Runtime: Bun 1.1+, Node 22+, Go 1.25 | scope.md §62 |
| Backup format: ZIP with manifest.json, data.json, media/ | backup-and-restore.md §17 |
| Schema version in manifest for forward compatibility | backup-and-restore.md §18 |
| Import flow: upload → validate manifest → validate schema → dry-run → summary → confirm → restore | backup-and-restore.md §20 |
| No sensitive data in logs (quality gate) | product-brief.md §91 |
| p95 performance targets defined for UI, API, AI Export, Backup | product-brief.md §96-141, DEC-008 |
| No incremental/partial/cloud backup | scope.md §71 |
| Expected dataset: 5yr, 1500 daily logs, 300 exercises, 30K sets, 2K cardio, 2K body weight, 300 check-ins, 1200 photos metadata, 500 products, 300 templates, 1K overrides, 100 exports/reviews | product-brief.md §99 |

## Technical Gaps

### Gap 1: Environment Separation
No dev/staging/prod environment distinction or configuration topology. Single Docker Compose file implied for production use. No test environment strategy.

### Gap 2: Configuration Schema
No documented environment variables, configuration file format, or configuration validation. Secrets (PIN hash, database credentials) storage approach unspecified beyond "environment variables for Docker."

### Gap 3: Logging Infrastructure
Only constraint is "no sensitive data to logs." No logging framework, format (structured JSON?), levels (info/warn/error/debug), log rotation, or log transport defined.

### Gap 4: Metrics and Monitoring
No metrics collection. No health check endpoints. No monitoring infrastructure. Database query performance, cache hit rates, file I/O latency — all unobservable without instrumentation.

### Gap 5: No SLO Definition
DEC-008 p95 targets are performance objectives, not SLOs. Missing: SLI definitions, measurement window, error budget, burn rate, compliance period. "Best effort" targets (backup with media, 12mo export with photos) have no measurement boundary.

### Gap 6: No Alerting or Runbooks
No alert rules, notification channels, or on-call procedures. No runbooks for common scenarios (disk full, backup failure, migration failure, service down). Self-hosted means user is operator — no guidance exists.

### Gap 7: Backup/Import Operational Gaps
- Import with existing data behavior (merge/replace/error) — Q-ACTOR-08, Q-AC-15 remain open.
- CSV file inclusion mandatory or optional — Q-AC-16 open.
- Schema migration strategy for version changes — Q-EDGE-11 open.
- No backup scheduling, retention policy, or verification strategy.

### Gap 8: Capacity Planning
Expected dataset defined but no resource estimates (CPU, memory, disk for 5yr with photos). No guidance on when to expect degradation.

### Gap 9: Upgrade/Rollback
No migration run strategy, no zero-downtime deployment approach, no rollback procedure.

## Missing Source Artifacts

| Artifact | Why Required | Consolidated Question |
| --- | --- | --- |
| Environment topology / Docker Compose config | Foundation for all deployment | TQ-OPS-001 |
| Environment variable / config schema | Required for Docker startup and config validation | TQ-OPS-001 (consolidated) |
| Logging framework and format specification | Required for observability and the "no sensitive data to logs" quality gate | TQ-OPS-002 |
| Health check endpoint specification | Required for Docker health checks and operational monitoring | TQ-OPS-002 (consolidated) |
| Metrics instrumentation specification | Required for SLO measurement and capacity planning | TQ-OPS-003 |
| SLO definition (SLI, window, error budget) | Required to make DEC-008 targets operational | TQ-OPS-003 (consolidated) |
| Alert rules and runbook specification | Required for operational ownership | TQ-OPS-004 |
| Backup scheduling/retention/verification spec | Required for data ownership story | TQ-OPS-005 |
| Resource estimates / operational run procedures | Required for self-hosted user to plan infrastructure and upgrade safely | TQ-OPS-006 |

## Questions Raised

| ID | Severity | Parent | Question | Why It Matters | Needed Artifact Or Decision |
| --- | --- | --- | --- | --- | --- |
| TQ-OPS-001 | dev-blocking | none | What is the complete environment and configuration topology? | Every deployment, migration, and troubleshooting action depends on knowing environments, env vars, config files, and secrets. | docker-compose.yml spec, env var catalog, secrets approach |
| TQ-OPS-002 | dev-blocking | none | What logging framework, format, levels, and health check contract is used? | "No sensitive data to logs" gate cannot be enforced without logging spec; Docker health checks need endpoints. | logging framework decision, health endpoint spec |
| TQ-OPS-003 | dev-blocking | none | How are SLOs from DEC-008 measured and enforced? | p95 targets are not SLOs without SLI definitions, measurement instrumentation, error budgets, and compliance windows. | SLI definitions, metrics spec, instrumentation decision |
| TQ-OPS-004 | deferred | none | What alert rules and runbook content are needed for self-hosted operators? | Self-hosted user has no ops team; runbooks are the only fallback for failures. | alert rule spec, runbook content outline | Deferred to post-MVP: single-user self-hosted MVP can ship without production alerting infra. Owner: product owner. |
| TQ-OPS-005 | dev-blocking | none | What is the backup scheduling, retention, and verification strategy? | Data ownership promise requires automated, verifiable backup, not only manual export. | scheduling spec, retention policy, verification procedure | | open | TBD |
| TQ-OPS-006 | needs-owner-decision | none | What are the resource estimates and operational run procedures (upgrade, rollback, migration) for the expected 5-year dataset? | Self-hosted user needs infrastructure guidance and safe upgrade path. Resource planning and upgrade strategy are interdependent. | resource estimate decision, migration run procedure, rollback procedure | | open | TBD |

Backup/restore open questions from product-verified docs are carried forward:
- Q-ACTOR-08 (import with existing data behavior)
- Q-AC-15 (import with data behavior — duplicate coverage)
- Q-AC-16 (CSV file mandatory/optional)
- Q-EDGE-11 (schema version migration strategy)

## Answer Effects

No prior answers to analyze (initial run).

## Risks

1. **Backup import ambiguity**: Import with existing data behavior (Q-ACTOR-08) is unresolved — directly affects data ownership guarantee.
2. **No SLO enforcement**: Without metrics instrumentation, DEC-008 targets are aspirational and unverifiable.
3. **Logging gap conflicts with quality gate**: Product-brief §91 requires "no sensitive data to logs" but no logging framework or audit mechanism exists to enforce it.
4. **Self-hosted operational burden**: User is sole operator — missing runbooks and capacity guidance creates failure risk.

## Suggested Decisions

1. Consolidate Q-ACTOR-08 and Q-AC-15 into one import-behavior decision (replace vs error vs merge).
2. Adopt structured JSON logging (info/warn/error/debug) as the logging approach to support the quality gate.
3. Implement a health endpoint (`GET /health`) covering DB and Redis connectivity for Docker Compose health checks.
4. Measure DEC-008 p95 targets via simple request-level timing middleware rather than a full metrics platform for MVP.
5. Define backup as a manual feature only for MVP; defer automated scheduling to post-MVP.

## Traceability Candidates

| Product Source | Technical Artifact |
| --- | --- |
| product-brief.md §91 | Logging spec — enforce "no sensitive data" |
| product-brief.md §96-141 | SLO definitions, metrics instrumentation |
| scope.md §60-62 | docker-compose.yml, env config |
| scope.md §71 | Backup scheduling decision |
| backup-and-restore.md §20 | Import with existing data decision |
| backup-and-restore.md §17-18 | ZIP format spec, schema version handling |
| features/backup-and-restore.md open questions | Q-ACTOR-08, Q-AC-15, Q-AC-16, Q-EDGE-11 |