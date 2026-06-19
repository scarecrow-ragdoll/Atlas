# Operations Observability

## Environments And Config

**Decision (TDEC-010):**
- Minimum: 1 vCPU, 1 GB RAM, 10 GB storage
- Recommended: 2 vCPU, 4 GB RAM, 50+ GB storage
- PostgreSQL and media volumes must be persistent
- Backups should be stored outside the running container
- Redis optional; if used, not source of truth

No environment/config topology defined (TQ-OPS-001).

## Deployment And Rollout

Docker Compose is the implied deployment mechanism. No CI/CD pipeline, no release process, no rollback strategy defined.

## Logs Metrics Traces

No logging framework or log format defined (TQ-OPS-002). No metrics instrumentation for DEC-008 p95 targets (TQ-OPS-003). Key observations:
- PIN/export/photo content must not be logged (§24.1)
- No log levels defined
- No metric collection or dashboard
- No tracing for API request flows

## Alerts Runbooks Backups

- No alerting strategy (deferred post-MVP — TQ-OPS-004)
- No backup scheduling or retention strategy (TQ-OPS-005)
- No operational runbook for common failures
- No resource estimate for expected dataset (TQ-OPS-006)

## Operations Questions

| ID | Question | Severity | Status |
| --- | --- | --- | --- |
| TQ-OPS-001 | No environment/config topology defined | dev-blocking | **resolved** (TDEC-052) |
| TQ-OPS-002 | No logging framework or health check contract | dev-blocking | **resolved** (TDEC-053) |
| TQ-OPS-003 | No SLO/SLI/metrics instrumentation for DEC-008 p95 targets | dev-blocking | **resolved** (TDEC-054) |
| TQ-OPS-004 | Alerting/runbooks deferred | deferred | deferred |
| TQ-OPS-005 | No backup scheduling, retention, or verification strategy | dev-blocking | **resolved** (TDEC-055) |
| TQ-OPS-006 | No resource estimates or operational run procedures | needs-owner | **resolved** (TDEC-010) |