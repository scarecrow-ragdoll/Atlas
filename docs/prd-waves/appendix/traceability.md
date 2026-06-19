# Traceability

## Source To Scope Map

| Source Section | Scope |
| ---- | ----- |
| 9. Dashboard | product-capabilities |
| 10. Workout Diary | user-journeys |
| 11. Exercise Library | data-lifecycle |
| 12. Cardio | integrations-operations |
| 13. Body Measurements | data-lifecycle |
| 15. Nutrition | data-lifecycle |
| 16. Charts | client-experience |
| 17-18. AI Export | integrations-operations |
| 19. AI Review | product-capabilities |
| 20. Backup | integrations-operations |

## Scope To Wave Map

| Scope | Wave |
| ----- | ---- |
| product-capabilities | 01, 09 |
| user-journeys | 02, 03 |
| data-lifecycle | 02, 03, 04, 05 |
| integrations-operations | 04, 07, 08 |
| client-experience | 06 |
| security-compliance | 01, 09 |
| delivery-sequencing | wave-map |

## Wave To Source Map

| Wave | Source Sections |
| ---- | -------------- |
| 01 | 5.1-5.3, 7.2 |
| 02 | 11 |
| 03 | 10 |
| 04 | 12, 13 |
| 05 | 15 |
| 06 | 16 |
| 07 | 17, 18 |
| 08 | 19 |
| 09 | 20 |

## Question Map

| Question ID | Wave | Scope | Status |
| ----------- | ---- | ----- | ------ |
| Q-PIN-001 | 01 | security-compliance | open |
| Q-CHART-001 | 06 | product-capabilities | open |
| Q-WORKOUT-001 | 03 | integrations-operations | open |

## Decision Map

- All PRD decisions preserved
- PIN guard included as security feature
- Single-user model confirmed
- No social/multi-user in MVP