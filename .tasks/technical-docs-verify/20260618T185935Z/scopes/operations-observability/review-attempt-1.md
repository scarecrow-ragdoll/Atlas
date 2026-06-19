# Operations-Observability Review Attempt 1

<!-- FILE: .tasks/technical-docs-verify/20260618T185935Z/scopes/operations-observability/review-attempt-1.md -->

## Verdict

**needs-revision**

## Sources Read

- docs/product-verified/product-brief.md
- docs/product-verified/scope.md
- docs/product-verified/domain-model.md
- docs/product-verified/features/backup-and-restore.md

## Coverage Check

All operations-observability sub-areas covered: environments, config, deployment, monitoring, logs, metrics, SLOs, backup, capacity, runbooks, upgrade/rollback. No surface-level gaps.

## Evidence Check

Every technical fact is traceable to a product source (column `Source` present in Technical Facts table). Source delta DEC-008 is correctly identified as raw p95 targets, not SLOs. Backup/restore open questions (Q-ACTOR-08, Q-AC-15, Q-AC-16, Q-EDGE-11) are faithfully carried forward.

## No-Invention Check

**Passed.** No endpoints, schemas, Dockerfile content, env var names, or deployment topologies are fabricated. The suggested decisions section uses careful language ("adopt structured JSON logging", "implement a health endpoint") framed as suggestions, not requirements.

## Source-Gap Consolidation Check

**Passed.** Missing artifact classes are consolidated:
- TQ-OPS-001: environments + config (single consolidated question)
- TQ-OPS-002: logging + health checks
- TQ-OPS-003: metrics + SLO measurement
- TQ-OPS-005: backup scheduling + retention + verification

## Question Ledger Check

**Mostly passed.** Two issues:

1. **TQ-OPS-004 severity mismatch**: Marked `deferred` but no deferral owner or rationale is recorded. Per output contract: `deferred` requires "an explicit owner, a deferral rationale, and no implementation-blocking impact." The rationale (self-hosted MVP can ship without runbooks) exists implicitly but is not stated in the ledger row. Add a `Resolution` column value like: "Deferred to post-MVP: self-hosted single user does not require production alerting infrastructure at launch. Owner: product owner."

2. **TQ-OPS-006 and TQ-OPS-007**: Both marked `needs-owner-decision`. Acceptable — these genuinely require owner input. However, TQ-OPS-006 (resource estimates) and TQ-OPS-007 (upgrade/rollback) could be consolidated into one `needs-owner-decision` question about operational run procedures, since the answer to one likely informs the other.

## Answer Effect Check

No prior answers to evaluate (initial run). Correct.

## Missing Or Unsupported Claims

None.

## Required Revisions

1. Add deferral owner and rationale to TQ-OPS-004 in the question ledger.
2. (Optional) Consider consolidating TQ-OPS-006 and TQ-OPS-007 into a single "operational run procedures" question, since resource planning and upgrade strategy are likely decided together.

## Approval Notes

Once the single revision above is applied, this report is ready for approved status. The analysis is thorough, well-consolidated, traces cleanly to sources, and invents nothing.