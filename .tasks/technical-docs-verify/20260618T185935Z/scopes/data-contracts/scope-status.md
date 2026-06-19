# Data-Contracts Scope Status

## Summary
- Run: 20260618T185935Z
- Scope: data-contracts
- Status: approved (gap analysis complete, questions documented)
- Cycles: 2 (worker-attempt-1 → needs-revision → worker-attempt-2 → approved)

## Key Findings
1. **DEC-007 userId ambiguity**: 11 entities declare userId FK in key column but omit userId in attributes. 6 child entities have no userId at all. DailyNutritionOverride has direct contradiction between key column and attributes.
2. **DEC-009 stale references**: 3 references to "WorkoutDay" remain after DailyLog rename.
3. **Missing enums**: BodyWeightEntry.source, BodyMeasurement.measurementType, WeekFlag.flagType are "enum undefined".
4. **No migration strategy**: No tool, versioning, seed data, or fixture format defined.
5. **No index strategy**: p95 SLOs exist but no index design to meet them.

## Open Questions
- **dev-blocking**: 5 (TQ-DATA-002, 003, 004, 005, 010)
- **needs-owner-decision**: 2 (TQ-DATA-001, TQ-DATA-007)
- **deferred**: 1 (TQ-DATA-006)
- **watchlist**: 3 (TQ-DATA-008, 009, 011)

## Approval Evidence
- Worker report: .tasks/technical-docs-verify/20260618T185935Z/scopes/data-contracts/worker-attempt-2.md
- Reviewer report: .tasks/technical-docs-verify/20260618T185935Z/scopes/data-contracts/review-attempt-2.md
- Verdict: approved

## Cross-Scope Dependencies
- api-contracts: Entity attribute definitions affect API schemas.
- operations-observability: Migration strategy, backup format, and index strategy affect ops runbooks.
- testing-delivery: Seed data format and fixture strategy affect test data setup.