# Product Brief

## Product Intent

Atlas is a self-hosted web application for personal workout, nutrition, body measurement, progress photo tracking, and structured AI data export. The core idea is not merely to store training data but to automate a weekly analysis cycle: the user logs data during the week, performs a weekly body check-in, and the application generates an AI-ready prompt and data package for analysis and recommendations.

## Target Users

Single user per instance. Self-hosted deployment by a technically proficient individual. No multi-user mode, registration, or public profiles in MVP.

Target user characteristics:
- Fitness enthusiast tracking training progression
- Comfortable with self-hosted Docker deployment
- Wants AI-powered analysis of personal data
- Values privacy and data ownership

## Jobs To Be Done

1. Log workout data (exercises, sets, reps, weights) quickly
2. Track progression of working weights over time
3. Log cardio sessions with intensity context
4. Track body weight, measurements, and progress photos
5. Maintain a simple weekly nutrition template with daily overrides
6. View progress charts for training, body, and nutrition data
7. Export structured data for AI analysis
8. Save AI review history and planned actions
9. Perform full backup/restore of all data

## Value Proposition

Atlas reduces the friction of fitness data logging while enabling AI-powered weekly analysis. Users own their data completely (self-hosted, full export, full restore) and never depend on a third-party service for data access.

## Success Metrics

Atlas MVP is considered product-ready when all functional acceptance criteria are implemented and the following success metrics are met.

### Functional Success Metrics

The user must be able to complete the full weekly workflow without using external spreadsheets or notes:

1. Create and edit exercises.
2. Add at least one full training day with multiple exercises and sets.
3. Add cardio to the same day or another day.
4. Add body weight entries across the week.
5. Create a weekly body check-in with measurements and photos.
6. Create a weekly nutrition template.
7. Override nutrition for one specific day.
8. View charts for training, body metrics, and nutrition.
9. Generate an AI prompt and AI export package for the last 4 weeks.
10. Save an AI review response.
11. Export a full backup.
12. Restore the backup into a clean instance.

### Data Completeness Metrics

AI export for a selected period must include:
- 100% of workouts in the selected period
- 100% of exercises used in those workouts
- 100% of workout sets
- 100% of exercise comments
- 100% of cardio entries
- 100% of body weight entries
- 100% of body check-ins
- 100% of measurements
- 100% of nutrition templates relevant to the selected period
- 100% of daily nutrition overrides
- user goal
- persistent AI context
- selected one-time AI context
- week flags

### Backup/Restore Success Metrics

Full backup/restore is successful when:
- all database entities are restored
- all media files included in backup are restored
- entity relationships remain valid
- restored data produces the same charts as the original instance
- restored data can generate an equivalent AI export
- import validation prevents broken or partial imports

### Quality Gates

MVP is not ready unless:
- `bun run verify:coverage` passes
- critical user flows are covered by automated tests
- e2e tests cover the weekly workflow
- backup export/import is covered by integration tests
- PIN guard is covered by tests
- AI export schema is covered by snapshot/schema tests
- no sensitive data is written to logs

Reference: DEC-006 (resolved Q-SCOPE-001)

## Performance Targets

### Expected Personal Dataset

Performance targets tested against: 5 years, 1,500 daily logs, 300 exercises, 30,000 workout sets, 2,000 cardio entries, 2,000 body weight entries, 300 body check-ins, 1,200 progress photos metadata, 500 nutrition products, 300 nutrition templates, 1,000 daily nutrition overrides, 100 AI exports/reviews.

### UI Performance (p95, local or small VPS)

| Page | Target |
| --- | --- |
| Dashboard initial data load | <= 1.5s |
| Daily log page initial load | <= 1.5s |
| Exercise list load | <= 1.0s |
| Exercise detail with history | <= 1.5s |
| Body metrics page load | <= 1.5s |
| Nutrition page load | <= 1.5s |
| Charts page initial render after data | <= 2.0s |

### API Performance (p95)

| Operation | Target |
| --- | --- |
| Simple entity create/update mutation | <= 300ms |
| Daily log query by date | <= 500ms |
| Exercise history query | <= 700ms |
| Body metrics query for period | <= 700ms |
| Nutrition summary for period | <= 700ms |
| Chart data query for period | <= 1.0s |

### AI Export (p95)

| Scenario | Target |
| --- | --- |
| 4 weeks without photos | <= 5s |
| 4 weeks with photos | <= 20s |
| 12 months without photos | <= 15s |
| 12 months with photos | Best effort, show progress |

### Backup (p95)

| Scenario | Target |
| --- | --- |
| Database-only backup | <= 15s |
| Backup with media | Best effort (media-size dependent) |
| Dry-run validation (db-only) | <= 15s |
| Import (db-only) | <= 30s |

### UX Rule

Any operation expected to take >2s must show a loading state.

Reference: DEC-008 (resolved Q-SCOPE-004)