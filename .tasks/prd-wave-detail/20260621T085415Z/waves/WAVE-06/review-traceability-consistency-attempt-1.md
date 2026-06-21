# WAVE-06 Traceability-Consistency Review Attempt 1

## Verdict
needs-revision

## Sources Read
- All 6 planner reports
- docs/prd-waves/waves/wave-06.md
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/business-rules.md
- docs/product-verified/edge-cases.md
- docs/prd-waves/frontend-pages/page-008.md

## Coverage Check
All planners include traceability sections. Source docs referenced consistently. Wave prefix used for all IDs.

## Evidence Check
- Source AC-065–073 mapped to wave-local AC-W06-001–015 ✓
- Source RULE-012 (Epley) → AC-W06-002 ✓
- Source RULE-013 (volume) → implicit in AC-W06-001 ✓
- Source RULE-015 (weekly avg) → AC-W06-009 ✓
- Source EDGE-008 → AC-W06-003, AC-W06-005, AC-W06-008, AC-W06-011 ✓
- Source EDGE-026 → noted in data-ops planner ✓

## Codebase Fit Check
Traceability to specific files (service/, resolver/, schema/, models/) included in architecture planner.

## Other-Wave Fit Check
Prior wave traceability (WAVE-04, WAVE-05) present.

## Acceptance Criteria Check
Issues:
1. AC-W06-001 is defined in product-ac planner but also referenced in architecture planner — needs single source.
2. AC numbering goes to 015 but some gaps in consistency checks.
3. AC-W06-012 default 12-week period has no traceability to any source document.

## Question Ledger Check
DQ IDs use W06 prefix correctly. Severity classifications appropriate. However, DQ-W06-004 (default period) was raised in architecture planner but does not appear in the question ledger.

## Unsupported Or Invented Claims
- 12-week default period has no source evidence.
- 52-week max range has no source evidence.

## Required Revisions
1. Add DQ-W06-004 (default chart period) to question ledger — currently missing.
2. Move 12-week default to open question (DQ) with "needs owner decision" severity.
3. Move 52-week max range to open question (DQ) or document as design decision.
4. Consolidate AC-W06-001 references to a single source (product-ac planner) and split as per product-reviewer feedback.
5. Add traceability from measurement overlay to source PRD §16.3 (AC-071).

## Approval Notes
Traceability structure is good but has gaps: missing DQ-W06-004 in ledger, unsourced defaults, and references need consolidation. Awaiting revision.