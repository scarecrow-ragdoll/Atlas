# WAVE-02 traceability-consistency Review Attempt 1

## Verdict
needs-revision

## Sources Read
- All 6 planner reports (attempt 1)
- docs/prd-waves/waves/wave-02.md (source wave)
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/domain-model.md
- docs/technical-verified/api-contracts.md
- docs/technical-verified/data-contracts.md
- docs/prd-wave-details/waves/wave-01.md
- docs/prd-waves/frontend-pages/page-003.md

## Coverage Check
All planners include traceability sections linking proposed ACs, ECs, tests, and questions to source documents. Product ACs (AC-002, AC-003, AC-004, AC-043 through AC-047) are mapped to new WAVE-02 ACs.

## Evidence Check
Source traces are consistent with source documents. Each planner cites its sources at the top.

## Codebase Fit Check
Not directly applicable to this perspective.

## Other-Wave Fit Check
WAVE-01 references in planners are consistent. WAVE-03 references are correct (exercise selector dependency).

### Issues Found

1. **AC ID numbering inconsistency**: The product-ac planner proposes AC-W02-001 through AC-W02-022. However, the actual count of unique ACs after deduplication is ~18 (per product-scope-and-ac reviewer). The numbering implies more ACs than actually exist in the final detailed wave. ACs should be renumbered to match the final count.

2. **Question ID duplication**: DQ-W02-001 (product-ac planner) and DQ-W02-004 (data-integration-ops planner) ask the same question: "Should deleting ExerciseMedia delete the physical file?" These should be merged into one question with a single ID. Duplicate questions create traceability confusion.

3. **Cross-planner EC duplication**: EC-W02-001 appears in both product-ac and testing-exit planners. Same with EC-W02-002, EC-W02-003, EC-W02-007. These should be defined once in a consolidated list — ideally in the testing-exit planner with references from other planners.

4. **Missing AC traceability source map**: While each AC traces to source docs, there is no consolidated traceability table mapping each AC-W02-* to its source document (PRD section, product-acceptance-criteria.md ID, or business rule). For a ready-for-dev wave, this consolidated map is required.

5. **Test ID numbering gap**: TEST-W02-009 through TEST-W02-013 appear in data-integration-ops planner but use different numbering than testing-exit planner (which stops at TEST-W02-008, then jumps to TEST-W02-020). There are gaps and overlaps. All TEST IDs should be consolidated and sequential.

6. **Missing traceability for decision IDs**: The planners reference TDEC-008 (file size), TDEC-004 (audit), TDEC-037 (media access contradiction) but don't have a DDEC-W02-XXX decision log for WAVE-02-specific decisions. The detailed wave needs a decision log section.

7. **Source wave outcomes not explicitly mapped**: OUT-W02-001 through OUT-W02-004 from the source wave are not explicitly mapped to ACs in the planners. The traceability section should map:
   - OUT-W02-001 → AC-W02-001 through AC-W02-005 (exercise CRUD)
   - OUT-W02-002 → AC-W02-006 through AC-W02-008, AC-W02-013 (media)
   - OUT-W02-003 → AC-W02-009 (working weight)
   - OUT-W02-004 → AC-W02-021, AC-W02-022 (API ready for workout diary)

## Acceptance Criteria Check
18 unique ACs after deduplication. Coverage is adequate for the scope.

## Exit Criteria Check
22 ECs after deduplication. Some ECs are vague (EC-W02-001). Need tightening per testing-exit reviewer.

## Verification Check
23 test obligations. Gaps and overlaps in numbering need consolidation. Good coverage breadth.

## Question Ledger Check
8 questions raised. DQ-W02-001 and DQ-W02-004 are duplicates — merge to DQ-W02-001. DQ-W02-006 is properly deferred. DQ-W02-005 (MIME detection) is valid but should have a tentative decision to move forward.

## Unsupported Or Invented Claims
No unsupported claims found across all planners. All proposals are grounded in source documents.

## Required Revisions
1. **Merge duplicate questions**: Combine DQ-W02-001 and DQ-W02-004 into DQ-W02-001. Update all references.
2. **Consolidate AC numbering**: Create a single, sequential, deduplicated AC list (AC-W02-001 through AC-W02-0N).
3. **Consolidate EC numbering**: Create a single, sequential EC list (EC-W02-001 through EC-W02-0N). Remove per-planner duplication.
4. **Consolidate TEST numbering**: Create a single, sequential TEST list (TEST-W02-001 through TEST-W02-0N). Fill gaps, remove overlaps.
5. **Add source outcome mapping**: Add explicit traceability from OUT-W02-* to AC-W02-*.
6. **Add decision log entries**: Create DDEC-W02-XXX entries for WAVE-02-specific decisions (e.g., soft delete vs hard delete, file storage path pattern, trigram index).
7. **Add consolidated traceability table**: Map every AC-W02 to its source (product AC, PRD section, business rule, or domain model entity).
8. **Generate consolidated question ledger**: Merge questions from all planners into a single deduplicated ledger.

## Approval Notes
Good foundation. The deduplication and consolidation work is mechanical but necessary for a ready-for-dev wave. After the 8 revisions, will approve.