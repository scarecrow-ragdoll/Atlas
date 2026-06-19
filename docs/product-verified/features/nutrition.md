# Nutrition

## Source Evidence

PRD §15, §26.7, §26.8.

## User Problem

Track daily nutrition without the overhead of a full calorie tracker. Use a weekly template with daily overrides.

## Scope

In MVP. No recipes, barcode scanner, water/fiber/salt/alcohol tracking, food recognition, or public product database.

## Behavior

- User creates products manually (name, calories/protein/fat/carbs per 100g)
- Weekly template: one template at a time, products with gram amounts, optional meal label
- Template auto-applied to all days of its week
- Daily override: add/subtract/replace products per day
- KJBJU calculated from template items; recalculated on override changes

## Derived Requirements

| Requirement | Source | Rationale | Confidence |
| --- | --- | --- | --- |
| KJBJU calculated from product per-100g values × grams / 100 | §15.2, §15.3 | Product has per 100g, template has grams | High |
| Template auto-applies to all days | §15.4 | "приложение считает, что с понедельника по воскресенье питание соответствует шаблону" | High |

## Edge Cases

EDGE-003 (0/negative nutritional values), EDGE-009 (empty template), EDGE-017 (mid-week template creation), EDGE-019 (delete referenced product).

## Acceptance Criteria

AC-017 through AC-019, AC-058 through AC-064, AC-113.

## Dependencies

Nutrition product catalog, date range calculation.

## Open Questions

Q-FEAT-008: Template week-over-week lifecycle (auto-renew vs manual).
Q-FEAT-013: Template applied mid-week — retroactive or forward only?
Q-ACTOR-06: Reset daily overrides to template.
Q-ACTOR-21: Nutrition template after goal change.