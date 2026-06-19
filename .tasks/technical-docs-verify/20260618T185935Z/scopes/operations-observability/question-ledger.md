# Operations-Observability Question Ledger

<!-- FILE: .tasks/technical-docs-verify/20260618T185935Z/scopes/operations-observability/question-ledger.md -->

## Technical Questions

| ID | Scope | Severity | Parent | Question | Why It Matters | Needed Artifact Or Decision | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| TQ-OPS-001 | operations-observability | dev-blocking | none | What is the complete environment and configuration topology? | Every deployment, migration, and troubleshooting action depends on knowing environments, env vars, config files, and secrets. | docker-compose.yml spec, env var catalog, secrets approach | worker-attempt-1.md | open | TBD |
| TQ-OPS-002 | operations-observability | dev-blocking | none | What logging framework, format, levels, and health check contract is used? | "No sensitive data to logs" gate cannot be enforced without logging spec; Docker health checks need endpoints. | logging framework decision, health endpoint spec | worker-attempt-1.md | open | TBD |
| TQ-OPS-003 | operations-observability | dev-blocking | none | How are SLOs from DEC-008 measured and enforced? | p95 targets are not SLOs without SLI definitions, measurement instrumentation, error budgets, and compliance windows. | SLI definitions, metrics spec, instrumentation decision | worker-attempt-1.md | open | TBD |
| TQ-OPS-004 | operations-observability | deferred | none | What alert rules and runbook content are needed for self-hosted operators? | Self-hosted user has no ops team; runbooks are the only fallback for failures. | alert rule spec, runbook content outline | worker-attempt-1.md | deferred | Deferred to post-MVP: single-user self-hosted MVP can ship without production alerting infra. Owner: product owner. |
| TQ-OPS-005 | operations-observability | dev-blocking | none | What is the backup scheduling, retention, and verification strategy? | Data ownership promise requires automated, verifiable backup, not only manual export. | scheduling spec, retention policy, verification procedure | worker-attempt-1.md | open | TBD |
| TQ-OPS-006 | operations-observability | needs-owner-decision | none | What are the resource estimates and operational run procedures (upgrade, rollback, migration) for the expected 5-year dataset? | Self-hosted user needs infrastructure guidance and safe upgrade path. Resource planning and upgrade strategy are interdependent. | resource estimate decision, migration run procedure, rollback procedure | worker-attempt-1.md | open | TBD |

## Carried Forward Product Questions

| ID | Question | Severity | Status |
| --- | --- | --- | --- |
| Q-ACTOR-08 | Import when data already exists (merge/replace/error) | dev-blocking | open |
| Q-AC-15 | Import with existing data behavior | dev-blocking | open |
| Q-AC-16 | CSV files — mandatory or optional | needs-owner-decision | open |
| Q-EDGE-11 | Migration strategy for schema version changes | dev-blocking | open |