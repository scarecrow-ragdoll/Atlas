# Reviewer Report: Product Scope & AC (WAVE-09)

**Perspective:** product-scope-and-ac
**Attempt:** 1
**Verdict:** approved-with-questions

## Review Findings
1. **AC coverage is complete** — all 13 AC-W09-001 through AC-W09-013 map directly to source ACs (AC-093-102, AC-114-116, AC-124) and business rules (RULE-007, RULE-008, RULE-009, RULE-028)
2. **Outcome mapping** — all 5 source outcomes (OUT-W09-001 through -005) covered by ACs
3. **Edge case coverage** — EDGE-010 (invalid ZIP), EDGE-021 (partial restore), EDGE-028 (schema migration) all covered
4. **Performance targets** from product-brief captured (db-only <= 15s, with media best-effort)
5. **Excluded scope** respected — no cloud backup, no incremental backup

## Required Revisions
None. Scope and AC coverage is sufficient.

## Blocking Questions
- **DQ-W09-001** (existing data behavior) needs owner decision before dev can start
- **DQ-W09-002** (CSV files) — recommended: exclude (no AC requires it)

## Verdict Rationale
Product scope is well-defined. AC mapping is complete and traceable to source. Two open questions exist but one (CSV) is recommended-exclude and not blocking.