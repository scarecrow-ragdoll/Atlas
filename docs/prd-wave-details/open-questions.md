# Open Questions
## Wave-Blocking
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W09-001 | WAVE-09 | product-ac | wave-blocking | Q-ACTOR-08, Q-AC-15 | Import behavior when data already exists (merge/replace/error) | Restore transaction design depends on this | Decision: merge, replace-silently, or reject-with-error | planner/product-ac | open | No resolution |
| DQ-W09-005 | WAVE-09 | architecture-codebase | wave-blocking | Q-EDGE-11 | Migration strategy for schema version differences (same-version-only vs migration runner) | Schema comparison logic depends on this | Decision: same-version-only for MVP or implement migration runner | planner/architecture-codebase | open | Recommended: same-version-only for MVP |
## Needs Owner Decision
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W09-002 | WAVE-09 | product-ac | needs-owner-decision | Q-AC-16 | CSV files in backup — mandatory or optional? | Affects ZIP content generation | Confirm: exclude CSV | planner/product-ac | open | Recommended: exclude |
| DQ-W09-003 | WAVE-09 | data-integration-ops | needs-owner-decision | — | Import ZIP upload size limit | MaxBytesReader config parameter | Max size in MB | planner/data-integration-ops | open | Recommended: 500MB |
| DQ-W09-004 | WAVE-09 | security-privacy-compliance | needs-owner-decision | — | Should backup/import operations be logged differently than AI export? | Privacy compliance (AC-117-120) | Confirm: log event metadata only (not content) | planner/security-compliance | open | Recommended: log event metadata only |
## Deferred
None.
## Watchlist
| ID | Question | Why It Matters | Timeline |
| --- | --- | --- | --- |
| — | All 14+ entity services need ListAllByUserID methods | Scope of entity service additions | Implementation slice planning |
## Resolved This Run
None — all questions remain open.