# Implementation Slices

## Slice Map

No implementation slices can be defined with confidence until foundational technical artifacts are created. The following ordering is implied by architecture dependencies:

| Slice | Description | Depends On | Status | Verification |
| --- | --- | --- | --- | --- |
| Slice 0 | Foundation: Docker Compose, Nx workspace, Go API skeleton, DB schema, Redis, settings, PIN guard | Architecture decision, API protocol decision | Blocked by TQ-ARCH-001, TQ-ARCH-002, TQ-API-001 | Integration tests for settings and PIN |
| Slice 1 | Exercise Library: CRUD + media upload | Slice 0 | Blocked by API endpoints | Unit + integration tests |
| Slice 2 | DailyLog + Cardio + WorkoutExercise + WorkoutSet | Slice 0 | Blocked by API endpoints, DailyLog data model alignment | Unit + integration tests |
| Slice 3 | Body tracking: weight entries, check-ins, measurements, photos | Slice 0 | Blocked by API endpoints | Unit + integration tests |
| Slice 4 | Nutrition: products, templates, daily overrides | Slice 0 | Blocked by API endpoints | Unit + integration tests |
| Slice 5 | Charts: training, body, nutrition | Slices 1-4 | Blocked by API contract, chart aggregation query | Integration + snapshot tests |
| Slice 6 | AI Export: prompt builder, ZIP generation | Slices 1-4 | Blocked by async/sync contract | Schema snapshot + integration tests |
| Slice 7 | AI Review: manual entry, history | Slice 0 | None | Unit tests |
| Slice 8 | Backup: full export + import | All data slices | Blocked by async/sync contract, backup import flow API | Integration + e2e tests |
| Slice 9 | E2E: weekly workflow coverage | All slices | Blocked by e2e test strategy | E2e tests |

## Dependencies

Cross-cutting dependencies that block all slices:
- API protocol decision (TQ-API-001)
- Component architecture decision (TQ-ARCH-002)
- Session/auth contract (TQ-AUTH-001..003, TQ-AUTH-005)
- Data model userId alignment (TQ-DATA-001)
- Migration strategy (TQ-DATA-005)

## Blockers

47 dev-blocking questions listed in appendix/question-ledger.md. Top 7 blockers for any implementation to start:
1. No system context or component architecture (TQ-ARCH-001, TQ-ARCH-002)
2. No API protocol decision (TQ-API-001)
3. No endpoint catalog (TQ-API-002)
4. No PIN session/auth contract (TQ-AUTH-001, TQ-AUTH-002, TQ-AUTH-003)
5. No UI state machine contract (TQ-CLIENT-001)
6. No deployment/config topology (TQ-OPS-001)
7. No test data factory strategy (TQ-TEST-001)

## Verification

Every slice must pass: unit tests, lint, typecheck, and integration tests where applicable. E2E tests required before slice handoff.