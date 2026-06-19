# Wave Map

## Top-Level Wave List

1. WAVE-01: Foundation - Infrastructure, database, API skeleton
2. WAVE-02: Exercise Library - CRUD exercises with media
3. WAVE-03: Workout Diary - Daily workouts with sets
4. WAVE-04: Cardio and Body Tracking - Cardio entries and measurements
5. WAVE-05: Nutrition - Products, templates, overrides
6. WAVE-06: Charts - Progress visualization
7. WAVE-07: AI Export and Prompt Builder - Export generation
8. WAVE-08: AI Review History - Save AI responses
9. WAVE-09: Backup Import/Export - Full backup operations

## Dependency Order

WAVE-01 → WAVE-02 → WAVE-03 → WAVE-04 ↔ WAVE-05 → WAVE-06 → WAVE-07 → WAVE-08 → WAVE-09

Note: WAVE-04 (cardio/body) and WAVE-05 (nutrition) can partially parallelize.

## Coverage Matrix

| PRD Section | Covered In Wave |
| --- | --- |
| 9. Dashboard | WAVE-01, WAVE-03, WAVE-04, WAVE-06 |
| 10. Workout Diary | WAVE-02, WAVE-03 |
| 11. Exercise Library | WAVE-02 |
| 12. Cardio | WAVE-04 |
| 13. Body Measurements | WAVE-04 |
| 15. Nutrition | WAVE-05 |
| 16. Charts | WAVE-06 |
| 17-18. AI Export | WAVE-07 |
| 19. AI Review | WAVE-08 |
| 20. Backup | WAVE-09 |

## More Than Eight Wave Check

9 waves - acceptable for MVP, scoped appropriately.

## Downstream Planning Recommendations

- WAVE-01 → detail-prd-wave for foundation planning
- Use $plan-backend-waves or detailed task planning for breakdowns

## User Approval

User approved the wave map on 2026-06-18. Status: waves-approved-by-user.