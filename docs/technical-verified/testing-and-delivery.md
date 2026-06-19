# Testing And Delivery

## Test Strategy

No complete test strategy exists for Atlas features. Current repo has vitest, Playwright, and a 100% coverage gate requirement. Quality gates from DEC-006 require:
- `bun run verify:coverage` passes
- Critical user flows covered by automated tests
- E2E tests cover the weekly workflow
- Backup export/import covered by integration tests
- PIN guard covered by tests
- AI export schema covered by snapshot/schema tests
- No sensitive data written to logs

## Fixtures And Test Data

No test data factory or fixture infrastructure defined (TQ-TEST-001). Missing:
- Exercise fixtures with working weights
- DailyLog+WorkoutExercise+WorkoutSet fixtures
- CardioEntry fixtures
- Body check-in fixtures with measurements and photos
- Nutrition product, template, and override fixtures
- AI export and backup manifest fixtures

## Contract And E2E Coverage

- Weekly workflow e2e strategy undefined (TQ-TEST-002)
- AI export schema snapshot test strategy undefined (TQ-TEST-003)
- Backup manifest schema test strategy undefined (TQ-TEST-004)
- No performance test policy for DEC-008 targets (TQ-TEST-007)
- No log redaction test policy (TQ-TEST-005)
- No test isolation strategy (TQ-TEST-006)

## Release Gates

- Coverage gate: `bun run verify:coverage`
- Pre-MR checks: lint, typecheck, unit tests
- Acceptance criteria coverage not formalized

## Testing Questions

| ID | Question | Severity | Status |
| --- | --- | --- | --- |
| TQ-TEST-001 | No test data factory or fixture strategy | dev-blocking | **resolved** (TDEC-056) |
| TQ-TEST-002 | Weekly workflow e2e strategy undefined | dev-blocking | **resolved** (TDEC-057) |
| TQ-TEST-003 | AI export schema snapshot test strategy undefined | dev-blocking | **resolved** (TDEC-058) |
| TQ-TEST-004 | Backup manifest schema test strategy undefined | dev-blocking | **resolved** (TDEC-059) |
| TQ-TEST-005 | Log redaction test policy undefined | needs-owner | **resolved** (TDEC-011) |
| TQ-TEST-006 | Test isolation strategy undefined | needs-owner | **resolved** (TDEC-012) |
| TQ-TEST-007 | Performance test policy for DEC-008 targets undefined | needs-owner | **resolved** (TDEC-013) |