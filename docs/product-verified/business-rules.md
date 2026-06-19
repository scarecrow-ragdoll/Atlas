# Business Rules

## Validation Rules

| Rule ID | Description | Source |
| --- | --- | --- |
| RULE-001 | PIN must be hashed before storage, never stored in plaintext | §7.2 |
| RULE-002 | PIN change requires current PIN verification | §7.2 (change allowed) |
| RULE-003 | Working weight must be numeric, stored in exercise reference | §10.6 |
| RULE-004 | Set weight and reps must be recorded; RPE/RIR optional | §10.5 |
| RULE-005 | Photo count 2-4 per check-in (requirement status ambiguous) | §13.2 |
| RULE-006 | Product nutritional values per 100g | §15.2 |
| RULE-007 | Backup ZIP must contain manifest.json and data.json with valid schema | §20.4 |
| RULE-008 | Dry-run validation before import | §20.4 |
| RULE-009 | No silent partial import — import must complete fully or fail | §20.4 |

## Calculation Rules

| Rule ID | Description | Source |
| --- | --- | --- |
| RULE-010 | KJBJU calculated from template items (product per 100g values × grams / 100) | §15.3 |
| RULE-011 | Daily nutrition recalculated on override changes | §15.5 |
| RULE-012 | Estimated 1RM (e1RM) — formula unspecified | §10.7 |
| RULE-013 | Volume per exercise = sum of (weight × reps) across all sets | §10.7 |
| RULE-014 | Best set — definition unspecified (heaviest weight or highest e1RM) | §10.7 |
| RULE-015 | Nutrition weekly averages = sum of daily values / 7 | §16.4 |

## State Transition Rules

| Rule ID | Description | Source |
| --- | --- | --- |
| RULE-016 | Workout day created on first save for a date | §10.2 |
| RULE-017 | Working weight snapshot captured when exercise is added to workout day; not updated retroactively | §10.6 |
| RULE-018 | Nutrition template applies to all days of its week anchor | §15.4 |
| RULE-019 | Daily override affects only its target date | §15.5 |
| RULE-020 | Single nutrition template at a time (replacement semantic unclear) | §15.3 |
| RULE-021 | AI export default period: last 4 weeks | §17.2 |

## Authorization Rules

| Rule ID | Description | Source |
| --- | --- | --- |
| RULE-022 | PIN disabled: application accessible without any check | §7.2 |
| RULE-023 | PIN enabled: all pages require valid session | §7.2 |
| RULE-024 | Media files not accessible without valid session | §14.3, §24.1 |
| RULE-025 | Photos excluded from AI export by default; opt-in required | §17.3 |

## Integration Rules

| Rule ID | Description | Source |
| --- | --- | --- |
| RULE-026 | AI export generated on user request only, not automatically | §24.1 |
| RULE-027 | AI prompt designed for manual copy-paste to external AI | §17.1 |
| RULE-028 | Backup export generated on user request only | §24.1 |
| RULE-029 | No automatic external API calls in MVP | §22, §23 |