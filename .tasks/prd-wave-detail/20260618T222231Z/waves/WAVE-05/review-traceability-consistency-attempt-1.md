# WAVE-05 Traceability-Consistency Review Attempt 1

## Verdict
needs-revision

## Sources Read
- All 6 planner reports
- docs/prd-waves/waves/wave-05.md
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/edge-cases.md
- docs/product-verified/business-rules.md
- docs/product-verified/domain-model.md

## Coverage Check
Traceability candidates listed in each planner. Source-to-AC/EC/test mapping is present across all planners.

## Evidence Check
Good source back-links to product-verified docs, PRD waves, and codebase files.

## Other-Wave Fit Check
N/A for this perspective.

## Acceptance Criteria Check
Some AC descriptions have inconsistencies between planners:
1. AC-W05-006 in product-ac says "hard deletion blocked when referenced by active template/override items" but data-ops planner uses soft-delete (isActive flag). These must be consistent.
2. AC-W05-032 mentions "warning marker" — not defined in data contract or error union types. Need to either add the marker type or remove the mention.

## Exit Criteria Check
EC list is consistent.

## Verification Check
Test IDs are consistent across reports.

## Question Ledger Check
5 questions from product-ac (DQ-W05-001–003), 1 from architecture (DQ-W05-004), 2 from data-ops (DQ-W05-005–006), 1 from security (DQ-W05-007), 1 from testing (DQ-W05-008), 1 from sequencing (DQ-W05-009). Need to consolidate into a single ledger.

## Unsupported Or Invented Claims
- Some traceability between product-level AC-058–AC-064 and WAVE-05 ACs is stated but not detailed per-product-AC. Should be explicit: which WAVE-05 AC maps to which product AC.

## Required Revisions
1. Resolve inconsistency: soft-delete vs hard-blocked for product deletion. Both planners must agree.
2. Remove or implement the "warning marker" concept in AC-W05-032.
3. Consolidate all DQ questions into a single question-ledger.md with unique IDs.

## Approval Notes
Good traceability foundation but needs consistency fixes. Approve after revisions.