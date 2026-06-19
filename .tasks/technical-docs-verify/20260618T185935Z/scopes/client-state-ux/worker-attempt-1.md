# client-state-ux Worker Attempt 1

## Sources Read
- docs/product-verified/functional-spec.md
- docs/product-verified/user-flows.md
- docs/product-verified/edge-cases.md
- docs/product-verified/product-brief.md (DEC-008 performance targets)
- docs/product-verified/scope.md
- docs/product-verified/domain-model.md
- docs/product-verified/features/dashboard.md
- docs/product-verified/features/exercise-library.md
- docs/product-verified/features/workout-diary.md
- docs/product-verified/features/cardio.md
- docs/product-verified/features/body-tracking.md
- docs/product-verified/features/progress-photos.md
- docs/product-verified/features/nutrition.md
- docs/product-verified/features/charts.md
- docs/product-verified/features/pin-guard.md
- docs/product-verified/features/ai-export.md
- docs/product-verified/features/ai-prompt-builder.md
- docs/product-verified/features/ai-review-history.md
- docs/product-verified/features/backup-and-restore.md
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/business-rules.md

## Source Delta Reviewed
- DEC-008: UI performance p95 targets per page and per API operation. UX Rule: any operation >2s must show loading state.

## Product Signals

| Signal | Source | UX Implication |
| --- | --- | --- |
| No empty states defined for any section | user-flows.md §Empty States, Q-ACTOR-12 | Every page lacks first-run rendering |
| UX Rule: >2s operations need loading state | product-brief.md §Performance Targets, DEC-008 | Loading states required by policy but undefined |
| "Field-level validation rules not specified in PRD" | functional-spec.md §Validations | Form validation entirely unspecified |
| EDGE-001: set with 0 weight or 0 reps | edge-cases.md | No validation rule |
| EDGE-003: 0 or negative nutritional values | edge-cases.md | No validation rule |
| EDGE-007: 0 or negative measurement value | edge-cases.md | No validation rule |
| EDGE-022: DB connection lost during save | edge-cases.md | Error feedback unspecified |
| EDGE-023: Redis unavailable for session | edge-cases.md | Error feedback unspecified |
| EDGE-024: Disk full during export | edge-cases.md | Error feedback unspecified |
| EDGE-025: Volume full during media save | edge-cases.md | Error feedback unspecified |
| EDGE-012: Session expired during data entry | edge-cases.md | No recovery UX |
| EDGE-016: Concurrent tab access | edge-cases.md | Last-write-wins undefined |
| No autosave mechanism defined | user-flows.md failure flows | Data loss risk on session loss |
| No formal state machines documented | domain-model.md §Lifecycle States | Even entity lifecycle states are implicit |
| No accessibility mentioned | entire source set | WCAG, ARIA, keyboard nav, screen reader all absent |
| No localization mentioned | entire source set | i18n, locale, number/date formatting absent |
| Offline-first is explicit non-goal | scope.md §Non-Goals | But offline fallback behavior not defined |
| PIN guard: session via cookie, TTL unspecified | features/pin-guard.md, EDGE-012 | Session lifetime undecided |
| PIN guard: no brute-force protection defined | features/pin-guard.md Q-AC-01 | Lockout/retry behavior undefined |
| "Last body weight" definition ambiguous | features/dashboard.md Q-AC-03 | Dashboard metric behavior unclear |
| "Training days this week" — calendar or trailing | features/dashboard.md Q-FEAT-001 | Week boundary definition unclear |
| Charts: charting library not chosen | features/charts.md Q-FEAT-009 | Renderer choice affects bundle size and UX |
| Media upload no retry/error feedback | user-flows.md failure flows | Upload UX not specified |
| Export failure no user feedback | user-flows.md failure flows | Export UX not specified |
| Import mid-failure no rollback UX | user-flows.md failure flows | Import UX not specified |

## Technical Facts

1. The app is a single-page web application (SPA) — implied by tech stack (Bun, Node, Go with GraphQL) and architecture.
2. Pages identified: PIN guard screen, Dashboard, Exercise Library, Exercise Detail, Workout Diary (with calendar), Body Tracking (check-in list, check-in detail, weight list), Nutrition (product catalog, template, daily override), Charts, AI Export, AI Review History, AI Prompt Builder, Settings, Backup/Restore.
3. Each page requires: loading state (initial data fetch), empty state (no data exists), error state (fetch/mutation failure), normal state (data displayed).
4. DEC-008 sets p95 page load targets: Dashboard <=1.5s, Daily log <=1.5s, Exercise list <=1.0s, Exercise detail <=1.5s, Body metrics <=1.5s, Nutrition <=1.5s, Charts <=2.0s.
5. The UX Rule from DEC-008 requires loading state for any operation >2s.
6. Tech stack: PostgreSQL, Redis, Bun/Node/Go. Redis used for session storage. No cache layer for UI data.

## Technical Gaps

### Gap 1: No UI State Machine Design (dev-blocking)
No page defines loading/empty/error/offline states. Every interactive page (13+) must ship with at least these states. Missing a unified or per-page state machine contract for:
- PIN guard: idle/pending/verified/error/locked-out states
- Dashboard: loading/ready/empty/error
- Exercise Library: loading/list empty/error
- Workout Diary: loading/today ready/backdated/empty day/error
- Body Tracking: loading/check-in list/check-in detail/weight list/all empty/error
- Nutrition: loading/catalog/template/day override/error
- Charts: loading/no data/error/ready
- AI Export: loading/ready/generating/success/error
- AI Review: loading/list empty/error
- Settings: loading/ready/saving/error
- Backup: loading/dry-run/confirming/exporting/restoring/error

### Gap 2: No Form Validation Contract (dev-blocking)
No validation rules for any field across all features. Missing:
- Field types, ranges, required/optional per field
- Client-side vs server-side validation responsibility
- Validation timing (on-blur, on-change, on-submit)
- Error display pattern (inline, banner, toast)
- Specific validation gaps per EDGE-001, EDGE-003, EDGE-007, EDGE-009

### Gap 3: No Loading State Design (dev-blocking)
Loading states required by DEC-008 UX Rule but no specification:
- Skeleton/spinner/progress bar design
- Per-page timeout (>2s rule enforcement)
- Loading overlay vs inline loading
- Page transition loading (SPA route changes)

### Gap 4: No Error State Design (dev-blocking)
Error recovery for all failure modes undefined:
- API errors (4xx, 5xx) — toast, inline, full-page, redirect
- Network errors (fetch timeout, connection refused)
- Backend unavailable (EDGE-022, EDGE-023)
- Disk/storage errors (EDGE-024, EDGE-025)
- Session errors (EDGE-012)
- Retry mechanism: auto-retry vs manual retry vs retry limit

### Gap 5: No Offline UX Design (dev-blocking)
Offline-first is an explicit non-goal, but the app needs graceful offline handling:
- What the user sees when the browser is offline
- Whether cached data is accessible
- Whether partial offline input is allowed
- Reconnection behavior

### Gap 6: No Cache or Data Freshness Contract (dev-blocking)
No cache strategy defined:
- Cache layer for UI data (HTTP cache, in-memory cache, React Query/SWR, or none)
- Stale-while-revalidate vs fetch-every-time
- Cache invalidation after mutations (add/update/delete exercise, save day, create check-in)
- Refetch behavior after navigation back
- Performance impact: without caching, every page navigation fetches fresh data, potentially exceeding DEC-008 targets

### Gap 7: No Realtime/Subscription Model (deferred)
- No WebSocket, SSE, or polling defined
- Multi-tab consistency (EDGE-016) unaddressed
- User could see stale data in a second tab
- Single-user self-hosted app can accept manual tab refresh. Deferrable to post-MVP.

### Gap 8: No Optimistic Update Contract (needs-owner-decision)
- Save operations (add set, create exercise, log body weight) have no optimistic behavior
- User waits for server confirmation before seeing result
- Decision: optimistic yes/no per operation; if yes, rollback pattern on failure

### Gap 9: No Accessibility Standard (dev-blocking)
- WCAG level not specified (A, AA, AAA)
- Keyboard navigation: all interactive elements must be keyboard-accessible
- ARIA roles, labels, live regions for dynamic content
- Focus management: after save, after modal close, after error
- Color contrast ratios
- Reduced motion support for animations/transitions
- Screen reader announcements for state changes (e.g., "workout saved", "error loading data")
- Form error announcements via aria-invalid, aria-describedby

### Gap 10: No Localization Strategy (dev-blocking)
- Target language(s) unknown (product brief implies Russian-speaking user from PRD Russian text, but not specified)
- Number formatting: weight (kg), measurements (cm), nutritional values (g, kcal) — locale-specific
- Date/time formatting: dd.MM.yyyy vs MM/dd/yyyy
- RTL support: probably not needed but unstated
- Translation system: i18next, react-intl, or none
- Locale detection and storage

### Gap 11: No Session Loss Recovery UX (dev-blocking)
- PIN session TTL undefined (features/pin-guard.md Q-ROLE-001)
- Session lost during data entry (EDGE-012): no autosave, no draft persistence
- When PIN guard redirects to PIN entry after session expires, unsaved work is lost
- Decision needed: autosave mechanism, draft storage, session renewal dialog

### Gap 12: No SPA Routing Contract (dev-blocking)
- Route map undefined
- Route guards (PIN guard) — interceptor pattern unspecified
- Lazy loading per page
- Error boundaries per route
- 404 handling for unknown routes
- Navigation transitions (none, fade, slide)
- Deep linking support (navigating to specific date in diary)

### Gap 13: Media Upload UX (dev-blocking)
- Upload progress feedback unspecified
- Retry on failure unspecified (EDGE-025)
- File type/accept restrictions
- File size limits (Q-FEAT-005)
- Multiple file upload behavior
- Upload cancellation

### Gap 14: Export/Import UX (dev-blocking)
- AI export: progress during ZIP generation, success/error notification, file download behavior
- Backup export: same as AI export
- Backup import: file upload UX, dry-run progress, summary display, confirmation dialog, restore progress, restore result
- Long-running operations (>2s): must show loading state per DEC-008. Some exports may exceed 2s.
- Large export (12mo with photos): "best effort, show progress" — progress indicator not designed

### Gap 15: Empty State Gap per Section (dev-blocking)
Dashboard: no data → show "Getting started" guidance per user-flows.md recommendation (not committed)
Exercise Library: no exercises → empty list with "Add first exercise" CTA
Workout Diary: no workouts → empty calendar or empty day
Body Tracking: no check-ins → empty list with "Add first check-in" CTA
Nutrition: no products → empty catalog with "Add first product" CTA
AI Exports: no exports → empty history
AI Reviews: no reviews → empty history

### Gap 16: PIN Brute-Force Lockout UX (dev-blocking)
- PIN guard has no lockout/retry behavior defined (Q-AC-01, EDGE-011)
- When PIN is locked, the UI must display: lockout state, remaining lockout timer, retry countdown
- Recovery path from lockout: timeout-based auto-unlock, no recovery PIN
- PIN entry screen states: idle → pending (loading) → error (wrong PIN, N retries remaining) → locked (timeout remaining)

### Gap 17: Calendar Navigation UX (dev-blocking)
- Workout Diary has calendar navigation to any past date
- Month-to-month loading behavior unspecified: pagination, year switching, date selection feedback
- Loading state while fetching month data
- Empty month rendering vs month with data indicators
- Future date handling (disable or allow)

### Gap 18: Nutrition Template Mid-Week Creation UX (needs-owner-decision)
- EDGE-017: template created mid-week
- UX for week-preview not defined: does the UI show the partial week? Does it apply retroactively to past days?
- Forward-only vs retroactive application affects date picker, confirmation dialog, and week visualization

## Missing Source Artifacts
1. UI state machine specification (per-page or unified)
2. Form validation contract (field types, ranges, required, error display)
3. Loading state design system (spinners, skeletons, placeholders)
4. Error state design (error types, display patterns, retry behavior)
5. Offline behavior decision record
6. Cache and data freshness contract
7. Realtime/subscription decision record
8. Optimistic update decision record
9. Accessibility standard declaration
10. Localization strategy decision record
11. Session loss recovery decision record
12. SPA routing contract
13. Media upload UX specification
14. Export/import UX specification
15. Empty state design per section
16. PIN brute-force lockout UX design
17. Calendar navigation UX specification
18. Nutrition template mid-week creation UX decision

## Questions Raised
Recorded in question-ledger.md as TQ-CLIENT-001 through TQ-CLIENT-017.

## Answer Effects
No prior technical answers to analyze. Source delta DEC-008 provides targets but not implementation decisions.

## Risks
1. **Scope risk**: 17 open questions (12 dev-blocking, 4 needs-owner-decision, 1 deferred) means this scope cannot reach approved-to-dev without owner decisions and design artifacts.
2. **Implementation risk**: Without state machine contract, validation rules, loading/error states, accessibility, and localization, every page will ship incomplete and may fail basic UX requirements.
3. **Performance risk**: Without cache strategy, pages may fail DEC-008 p95 load targets (e.g., exercise list must load in <=1.0s but without caching every navigation triggers full fetch).
4. **Accessibility risk**: Shipping without WCAG compliance could be a legal/reputational issue depending on deployment context.
5. **Localization risk**: If the user base expects Russian-language UI (PRD text is in Russian), shipping English-only with hardcoded strings requires rework.

## Suggested Decisions

| Decision | Type | Suggested Answer | Rationale |
| --- | --- | --- | --- |
| UI state convention | design | Unified state machine: loading | skeleton/spinner → empty | guidance CTA → error | inline toast + retry → ready | content | Covers all pages with one pattern |
| Validation architecture | design | Client-side validation on blur + server-side validation on submit. Display errors inline per field + summary banner. | Clear separation, good UX |
| Cache strategy | design | React Query/SWR stale-while-revalidate with 30s stale time, automatic refetch on window focus, manual invalidation after mutations | Covers DEC-008 targets, simple implementation |
| Realtime | decision | Polling every 60s on visible pages only. No WebSocket/SSE in MVP. | Simple, no infra cost, covers multi-tab use case |
| Optimistic updates | decision | Save operations: yes (instant feedback). Delete operations: yes (instant removal). Rollback on error with toast + undo option. | Good UX for fast user flow |
| Offline | decision | Graceful offline notice banner. Cached data readable. No offline writes. | Pragmatic given non-goal |
| Accessibility | decision | WCAG 2.1 AA. Keyboard nav for all interactive elements. Focus management. ARIA labels for dynamic content. | Industry standard |
| Localization | decision | Russian locale for MVP. i18next-ready with en fallback. Number format: comma decimal, kg/cm. Date: dd.MM.yyyy. | Matches PRD language, enables future i18n |
| Session loss | decision | Autosave draft to localStorage every 30s. On session expiry, save drafts, redirect to PIN, restore drafts on return. | Prevents data loss |
| PIN brute-force | decision | 3 retries → 30s lockout → 5 retries → 5m lockout. Lockout state in session store. | Pragmatic protection |
| PIN brute-force | decision | 3 retries → 30s lockout → 5 retries → 5m lockout. Lockout state displayed with countdown timer. PIN entry states: idle→pending→error→locked. | Pragmatic protection with clear UX |
| Calendar navigation | design | Month fetch on scroll/nav. Loading skeleton for month grid. Future dates disabled with muted styling. Date selection updates diary panel inline. | Covers DEC-008 loading requirement |
| Nutrition mid-week template | decision | Forward-only: template applies to current and future days. Past days remain as overridden. Confirmation dialog shows affected days. | Simple, no retroactive data changes |
| Photo count (2-4) | decision | 2-4 recommended, not required. Validation warns if <2 but allows submission. | Matches ambiguity from PRD |

## Traceability Candidates

| Trace Target | Evidence |
| --- | --- |
| Q-ACTOR-12 → TQ-CLIENT-003 | user-flows.md empty states gap formalized as technical question |
| DEC-008 → TQ-CLIENT-002 | Performance UX rule formalized as loading state gap |
| EDGE-022, EDGE-023, EDGE-024, EDGE-025 → TQ-CLIENT-004 | Error states missing for all failure modes |
| EDGE-012 → TQ-CLIENT-013 | Session loss UX gap |
| EDGE-016 → TQ-CLIENT-010 | Multi-tab consistency gap (deferred) |
| Q-AC-01, EDGE-011 → TQ-CLIENT-015 | PIN brute-force lockout UX gap |
| EDGE-017 → TQ-CLIENT-017 | Nutrition template mid-week UX gap |
| Q-FEAT-009 → TQ-CLIENT-014 | Route/lazy-load implications for charting library |
| Q-FEAT-005 → TQ-CLIENT-006 | File validation part of form validation contract |
| Workout Diary calendar → TQ-CLIENT-016 | Calendar navigation UX gap |