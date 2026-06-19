# Client State And UX Contracts

## User Interface States

No UI state machine defined for any page. Each page must handle:
- **Loading state**: initial data fetch (must show within p95 targets)
- **Empty state**: no data for current context (first-run, no exercises, no workouts)
- **Error state**: API failure, network loss
- **Success state**: data loaded and displayed
- **Offline state**: browser loses connectivity during session

Missing contracts affect all 11 MVP sections (TQ-CLIENT-001).

## Form And Validation Behavior

No form validation contract for any field (TQ-CLIENT-006). Required decisions:
- Field-level validations (required, format, range)
- Submission error display
- Optimistic vs pessimistic saves
- "Copy previous set" behavior
- Date picker behavior for calendar navigation

## Cache And Realtime Behavior

No cache strategy or data freshness policy defined (TQ-CLIENT-009). Risk: UI performance targets (DEC-008) may not be met without caching. No realtime requirements in MVP.

## Accessibility And Localization

No accessibility standard defined (TQ-CLIENT-011). No localization strategy (TQ-CLIENT-012). Both deferred.

## Client Questions

| ID | Question | Severity | Status |
| --- | --- | --- | --- |
| TQ-CLIENT-001 | No UI state machine for any page | dev-blocking | **resolved** (TDEC-044) |
| TQ-CLIENT-002 | No page-specific loading state contract | dev-blocking | **resolved** (TDEC-045) |
| TQ-CLIENT-003 | No empty state contract for any section | dev-blocking | **resolved** (TDEC-045) |
| TQ-CLIENT-004 | No error state contract for API failures | dev-blocking | **resolved** (TDEC-045) |
| TQ-CLIENT-005 | No offline/unavailable state contract | dev-blocking | **resolved** (TDEC-045) |
| TQ-CLIENT-006 | No form validation contract for any field | dev-blocking | **resolved** (TDEC-046) |
| TQ-CLIENT-007 | No CRUD operation feedback pattern (optimistic vs pessimistic) | dev-blocking | **resolved** (TDEC-046) |
| TQ-CLIENT-008 | No client-side error display pattern | dev-blocking | **resolved** (TDEC-047) |
| TQ-CLIENT-009 | No cache or data freshness contract | dev-blocking | **resolved** (TDEC-048) |
| TQ-CLIENT-010 | No navigation state management (back/forward, unsaved data) | dev-blocking | **resolved** (TDEC-049) |
| TQ-CLIENT-011 | No accessibility standard (WCAG, ARIA) | dev-blocking | **resolved** (TDEC-050) |
| TQ-CLIENT-012 | No localization strategy | dev-blocking | **resolved** (TDEC-051) |