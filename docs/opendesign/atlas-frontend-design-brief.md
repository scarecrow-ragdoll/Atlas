<!--
FILE: docs/opendesign/atlas-frontend-design-brief.md
VERSION: 1.0.0
START_MODULE_CONTRACT
  PURPOSE: Frontend design brief for OpenDesign — describes all MVP screens, layout, components, UX flows, visual direction, and design constraints for the Atlas fitness tracking application.
  SCOPE: Covers all 11 MVP sections: Dashboard, Exercise Library, Training Log, Cardio, Body Tracking, Nutrition, Charts, AI Export/Review, Import/Export, Settings, PIN auth. Excludes future/post-MVP features, backend implementation, and API contract changes.
  DEPENDS: docs/prd-waves/frontend-pages/*.md, docs/product-verified/*.md, docs/technical-verified/*.md, docs/prd-wave-details/waves/wave-01.md, docs/prd-wave-details/waves/wave-02.md
  LINKS: M-GRACE-WORKFLOW, WAVE-01, WAVE-02, PAGE-001 through PAGE-011
  ROLE: DOC
  MAP_MODE: SUMMARY
END_MODULE_CONTRACT
START_CHANGE_SUMMARY
  LAST_CHANGE: 1.0.0 - Initial design brief for OpenDesign generation.
END_CHANGE_SUMMARY
-->

# Atlas Frontend Design Brief for OpenDesign

---

## 1. Product Summary

**Atlas** is a self-hosted web application for personal fitness tracking. A single user runs it on their own server (Docker). The core workflow is a weekly cycle:

- **During the week**: log workouts (exercises, sets, weights, reps, RPE/RIR), cardio sessions, body weight entries, and daily nutrition.
- **Weekly check-in**: record body measurements, body fat %, and 2-4 progress photos.
- **End of week**: generate an AI-ready export package (ZIP with structured data + prompt) to send to ChatGPT or another AI for analysis and recommendations.

The app should feel like a **personal analytics dashboard** and **training log**, not a generic fitness tracker. It is data-focused, private, and designed for quick data entry on desktop.

**User**: single fitness enthusiast who self-hosts Docker. Technically proficient enough to deploy, but wants a polished UI for daily use.

**Key differentiator**: the app is built around the weekly AI analysis cycle, not just around workout logging. Every screen feeds into the weekly review.

---

## 2. Design Goals

| Goal | Description |
|------|-------------|
| **Fast data entry** | Minimize clicks and page loads for common tasks. Forms should be compact, tab-friendly, and validate inline. |
| **Minimal routine overhead** | Prefill working weights from exercise library. Default to today's date. Keep weekly template auto-applied. |
| **Calm, data-focused interface** | No flashing animations, no gamification, no social feeds. Clean typography, clear hierarchy, muted palette. |
| **Weekly review flow** | Dashboard, charts, and AI export are designed around the weekly cycle. User should feel a natural rhythm: log → check-in → export → review. |
| **Privacy and ownership** | Self-hosted feel: no cloud badges, no "share" buttons, no SaaS language. The app belongs to the user. |
| **Readable tables and forms** | Desktop-first. Tables with clear column alignment, sortable headers, consistent spacing. Forms with visible labels, clear validation, and sufficient field width. |
| **Desktop-first, responsive** | Primary use case is desktop browser (large screen for data entry). Must work on tablet. Mobile is secondary but should not break. |

---

## 3. Visual Direction

### Overall Style

Modern, strict, calm. Think **personal analytics dashboard** + **training log notebook**. Not a bodybuilding magazine. Not a gamified fitness app.

- Clean card-based layout with generous whitespace
- Minimal use of color — color is reserved for data (charts, status indicators, progress signals)
- Muted backgrounds with clear elevation hierarchy (cards, sidebar, modals)
- Compact but comfortable — information-dense without feeling crowded

### Color Mood

- **Light mode**: clean white/light gray background, dark text, subtle borders, muted accent for interactive elements
- **Dark mode**: dark gray background (#1a1a2e or similar), light text, reduced contrast for secondary elements
- **Accent color**: a single calm accent (teal, slate-blue, or muted indigo) used sparingly for links, active states, and primary buttons
- **Success/progress**: green for positive signals (weight up, volume up, check-in completed)
- **Warning/caution**: amber/yellow for stagnation, missed workouts
- **Destructive**: red for delete, archive (used sparingly)
- **Data visualization**: chart color palette should be accessible (colorblind-friendly, distinct hues)

Design assumption: no brand hex codes specified in source docs. OpenDesign may propose a palette.

### Typography

- Sans-serif system font stack or a neutral open-source font (Inter, SF Pro, or similar)
- Clean, readable at all sizes
- Hierarchical: page titles (large/semibold), section headings (medium/semibold), body (regular), data values (tabular numbers preferred)
- Monospace or tabular figures for numeric data (weights, reps, macros) for column alignment

### Spacing

- Consistent 4px/8px grid
- Cards: 16-24px padding
- Sidebar: 240-280px width on desktop
- Section spacing: comfortable but not wasteful — the user works with data on one screen

### Cards

- Rounded corners (8-12px)
- Subtle border or very light shadow for elevation
- No heavy drop shadows
- Card types:
  - **Stat card**: single metric with label and optional trend arrow
  - **Entity card**: exercise, product, or check-in summary with key fields
  - **Graph card**: chart with title and optional period filter

### Charts

- Clean line charts for weight/volume trends
- Bar charts for weekly comparisons
- Minimal grid lines, clear axis labels
- Interactive tooltip on hover showing exact values
- Responsive to window resize
- Library recommendation: Recharts or lightweight charting compatible with React/Next.js

### Icons

- Simple outlined icons (Lucide, Heroicons, or similar)
- Consistent 20x20 or 24x24 sizing
- Used for sidebar navigation, action buttons, status indicators
- No filled/duotone mixing

### Empty States

- Every list/page must have a deliberate empty state
- Illustration: simple line art or geometric placeholder (not cute characters)
- Text: clear call-to-action ("Create your first exercise" with action button)
- Subtle, not patronizing

### Destructive Actions

- Delete/archive actions require confirmation dialog
- Archive: softer language ("Move to inactive" vs "Delete")
- Actual delete: red destructive button in confirmation dialog
- Undo option via toast where feasible

---

## 4. Information Architecture

### Navigation Structure

```
Atlas App
├── Dashboard              # Weekly summary, quick actions
├── Training Log           # Daily workout entry by date
├── Exercises              # Exercise library CRUD
├── Body                   # Check-ins, weight, measurements, photos
├── Nutrition              # Products, template, daily overrides
├── Charts                 # Progress visualization
├── AI Export              # Prompt builder and ZIP download
├── AI Review              # Save and view AI analysis responses
├── Import / Export        # Full backup operations
└── Settings               # PIN, units, AI context, preferences
```

### Section Details

#### Dashboard (PAGE-001)
- **Purpose**: Weekly overview at a glance. Entry point to today's actions.
- **Data**: current date, last body weight, workout days this week, cardio sessions this week, current goal, next check-in reminder
- **Actions**: add workout (→ Training Log today), add cardio, add weight, open check-in (→ Body), generate AI report (→ AI Export)
- **Transitions**: each quick action navigates to the corresponding section with relevant date pre-selected

#### Training Log (PAGE-002)
- **Purpose**: Daily workout entry. Add exercises with sets.
- **Data**: exercises for selected date, sets per exercise (weight/reps/RPE/RIR), exercise comments, cardio entries
- **Actions**: select date via calendar, add exercise from library, add set, edit/delete set, add cardio, save
- **Transitions**: add exercise opens exercise selector modal → Exercises section

#### Exercises (PAGE-003)
- **Purpose**: Manage exercise library. Foundation for workout diary.
- **Data**: exercise list (name, muscle groups, working weight, media count, status)
- **Actions**: create, edit, view detail, upload media, archive/restore, search/filter
- **Transitions**: exercise detail → inline media gallery; row click → detail view

#### Body (PAGE-005, PAGE-006 combined)
- **Purpose**: Weekly check-ins, body weight, measurements, progress photos.
- **Data**: check-in history, weight entries, 10 measurement types, 2-4 photos per check-in
- **Actions**: create check-in, add weight entry, add measurements, upload photos, view gallery
- **Transitions**: photo thumbnail → full-size viewer

#### Nutrition (PAGE-007)
- **Purpose**: Food product catalog, weekly meal template, daily overrides.
- **Data**: products list, weekly template with gram amounts, daily override, macro summary
- **Actions**: create/edit products, create/edit template, override specific day
- **Transitions**: template item → product selector modal → Products section

#### Charts (PAGE-008)
- **Purpose**: Visualize progress over time.
- **Data**: exercise progress (working weight, best set, e1RM, volume), body (weight, fat %, measurements), nutrition (weekly averages)
- **Actions**: select exercise, select chart type, select date range
- **Transitions**: chart point hover → tooltip with exact value

#### AI Export (PAGE-009)
- **Purpose**: Generate AI prompt and export package.
- **Data**: date range, section toggles, persistent context, week flags, generated prompt, ZIP download
- **Actions**: configure export → generate → preview prompt → download ZIP → copy prompt

#### AI Review (PAGE-009 related)
- **Purpose**: Save AI analysis responses.
- **Data**: review history list, individual review (date range, AI response text, notes, planned actions)
- **Actions**: create review (paste AI response), view history, edit notes

#### Import / Export (PAGE-010)
- **Purpose**: Full backup and restore.
- **Data**: export progress, import validation summary, import results
- **Actions**: export with/without media, upload backup file, dry-run validation, confirm import

#### Settings (PAGE-011)
- **Purpose**: Application configuration.
- **Data**: PIN state, AI context fields, units, export preferences
- **Actions**: enable/disable/change PIN, edit AI context, change units, change default export weeks

---

## 5. Global Layout

### Desktop Layout

```
┌─────────────────────────────────────────────────────┐
│  ┌──────────┐  ┌──────────────────────────────────┐ │
│  │          │  │  Top Bar                          │ │
│  │  Sidebar │  │  [App Name]         [Session]    │ │
│  │  Nav     │  ├──────────────────────────────────┤ │
│  │          │  │  Main Content Area                │ │
│  │  ─────── │  │                                  │ │
│  │  Dashboard│  │  ┌─ Page Header ────────────────┐│ │
│  │  Training │  │  │  Title              [Action] ││ │
│  │  Exercises│  │  └──────────────────────────────┘│ │
│  │  Body     │  │                                  │ │
│  │  Nutrition│  │  ┌─ Content ────────────────────┐│ │
│  │  Charts   │  │  │  Cards / Tables / Forms     ││ │
│  │  AI Export│  │  │                              ││ │
│  │  AI Review│  │  │                              ││ │
│  │  Import   │  │  └──────────────────────────────┘│ │
│  │  Settings │  │                                  │ │
│  │          │  │                                  │ │
│  └──────────┘  └──────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

#### Sidebar Navigation
- Fixed left sidebar, 240-280px wide
- Vertical nav list with icons + labels
- Active page highlighted with accent color
- Collapsible on smaller screens (hamburger toggle)
- Bottom section: Settings icon, version number
- No badges, no notifications, no user avatar (single user)

#### Top Bar
- Thin horizontal bar above content area (not overlapping sidebar)
- Left: app name "Atlas" (non-clickable, small/medium text)
- Right: session status indicator (green dot + "Unlocked" when PIN session active, or lock icon when PIN disabled)
- No user avatar, no logout button (session is managed via PIN unlock/lock)

#### Main Content Area
- Fills remaining width after sidebar
- Scrollable vertically (no horizontal scroll)
- Padding: 24-32px from edges
- Background: slightly different tint from cards/sidebar for depth

#### Page Header Pattern
- Every page has a consistent header row
- Left: page title (h1)
- Right: primary action button (e.g., "Add Exercise", "Create Check-in", "Generate Export")
- Optional: secondary actions as icon buttons or dropdown

#### Primary Action Button
- Accent-colored, clearly visible
- Consistent position in page header
- Label is verb-based: "Add Exercise", "Log Workout", "Create Check-in"

### Responsive Behavior

- **Desktop (>1024px)**: full sidebar + content layout
- **Tablet (768-1024px)**: sidebar collapses to icon-only or hamburger menu
- **Mobile (<768px)**: sidebar hidden behind hamburger, stacked layouts, full-width cards
- Forms should remain usable on tablet (wider modals or full-page forms)
- Data tables may need horizontal scroll on small screens or switch to card list view
- Charts should resize responsively

### Session State

When PIN is enabled:
- If no valid session cookie, the app shows a full-screen PIN unlock overlay before any content
- The overlay is a centered modal card on a blurred/dimmed background
- The sidebar and nav are hidden until unlocked
- Once unlocked, the app works normally
- A "Lock" button in the top bar (or Settings) manually revokes the session

When PIN is disabled:
- No PIN overlay — app opens directly to Dashboard
- Top bar shows "PIN: Off" indicator
- Settings shows PIN as "Not configured"

---

## 6. Common UI Patterns

### Page Header

| Aspect | Description |
|--------|-------------|
| Used on | Every section page |
| Content | Title (left), primary action button (right), optional secondary actions |
| Behavior | Sticky top of content area on scroll |
| States | Title only (no actions on detail pages) |

### Cards

| Aspect | Description |
|--------|-------------|
| Used on | Dashboard (stat cards), entity lists, detail sections |
| Content | Varies: stat value + label, entity summary, form section |
| Behavior | Clickable when representing an entity; static for stat cards |
| States | Default, hover (subtle border/shadow change), selected |

### Data Table / List

| Aspect | Description |
|--------|-------------|
| Used on | Exercises, products, check-in history, AI review history, backup list |
| Content | Rows with columns, optional row actions (edit/delete/archive) as icon buttons |
| Behavior | Row click opens detail; sortable by column header click; optional search/filter bar above |
| States | Default rows, hover row highlight, selected row |

### Filter Bar

| Aspect | Description |
|--------|-------------|
| Used on | Exercise list, check-in history, charts |
| Content | Search input, toggle chips (active/inactive), optional dropdown filters |
| Behavior | Client-side or server-side filtering; clear filter button |

### Date Picker / Calendar

| Aspect | Description |
|--------|-------------|
| Used on | Training Log (date switcher), charts (date range), AI export (date range), check-in creation |
| Content | Month grid with day cells, prev/next month arrows, "Today" button |
| Behavior | Click day to select; date range mode: click start, click end; inline or popover |
| States | Selected date (accent fill), today (outline or dot), dates with data (subtle indicator) |

### Form Drawer / Modal / Page

| Aspect | Description |
|--------|-------------|
| Used on | Exercise create/edit, check-in form, product form, settings sections |
| Content | Form fields with labels, validation messages, save/cancel buttons |
| Behavior | Create: open modal/drawer; Edit: prefill form; Save closes on success; Cancel discards with unsaved warning |
| States | Filled, focused, validation error, saving (disabled + spinner), success (auto-close or toast) |

### Confirmation Dialog

| Aspect | Description |
|--------|-------------|
| Used on | Archive exercise, delete media, delete check-in, destructive settings changes |
| Content | Title ("Archive exercise?"), description ("Exercise will be moved to inactive."), Cancel + Confirm buttons |
| Behavior | Confirm performs action and closes; Cancel closes without action |
| States | Default, confirm loading (after click while API processes) |
| Type | Modal overlay with backdrop |

### File Upload Dropzone

| Aspect | Description |
|--------|-------------|
| Used on | Exercise media, progress photos |
| Content | Dashed border area with upload icon, "Drag & drop or click to upload" text, file type/size constraints listed below |
| Behavior | Click opens file picker; drag highlights dropzone; shows preview after selection; multiple files allowed |
| Validation | Reject unsupported types with inline error; reject oversized files; show file name/size per file |
| States | Empty (dashed border), dragging (highlighted border), uploading (progress per file), success (thumbnail), error (red border + message) |

### Media Gallery

| Aspect | Description |
|--------|-------------|
| Used on | Exercise detail, check-in detail |
| Content | Grid of thumbnails, optional delete button per item |
| Behavior | Click thumbnail opens lightbox/full-size viewer; video shows play icon overlay |
| States | Empty (hidden or "No media" text), loading (skeleton grid), loaded (thumbnail grid) |

### Graph Card

| Aspect | Description |
|--------|-------------|
| Used on | Charts page |
| Content | Chart title, chart visualization, optional period/dataset selector, tooltip on hover |
| Behavior | Responsive to container width; interactive tooltip; optional legend |
| States | Loading (spinner in card area), error (error message + retry), empty ("No data for this period"), loaded |

### Empty State

| Aspect | Description |
|--------|-------------|
| Used on | All list/table pages when no data exists |
| Content | Simple illustration + heading + description + action button |
| Example | "No exercises yet. Create your first exercise to get started." with "Add Exercise" button |
| States | Only one: empty (shown when list is empty after loading completes) |

### Loading State

| Aspect | Description |
|--------|-------------|
| Used on | All pages during data fetch |
| Content | Skeleton placeholders (pulsing gray rectangles mimicking card/row shape), or subtle spinner |
| States | One per fetch opportunity; avoid nested spinners |

### Error State

| Aspect | Description |
|--------|-------------|
| Used on | All pages on API failure |
| Content | Error message (human-readable), optional error detail, retry button |
| Behavior | Retry button re-fetches data; persistent until resolved |
| States | Single error state per error boundary |

### Validation State

| Aspect | Description |
|--------|-------------|
| Used on | All forms |
| Content | Red border on invalid field, error message below field, prevent form submission |
| Behavior | Validate on blur (field-level), validate all on submit; clear error when field is corrected |
| States | Field error, form-level error (top of form) |

### Success Toast

| Aspect | Description |
|--------|-------------|
| Used on | After successful mutations: save, delete, archive, upload, export, import |
| Content | Brief success message, optional undo action |
| Behavior | Auto-dismiss after 3-4 seconds; undo available for non-destructive actions |
| States | Visible (slide-in from top-right or bottom-center), dismissing (fade out) |

### Unsaved Changes Warning

| Aspect | Description |
|--------|-------------|
| Used on | Any form with unsaved data when user tries to navigate away or close modal |
| Content | Browser native or custom dialog: "You have unsaved changes. Discard them?" |
| Behavior | "Discard" closes without saving; "Stay" keeps form open |

---

## 7. Screen Inventory

Full list of MVP screens with route suggestions:

| # | Screen | Route Suggestion | Wave |
|---|--------|-----------------|------|
| A1 | PIN Unlock | `/pin-unlock` | WAVE-01 |
| A2 | PIN Protected (app access) | implicit — all routes behind guard | WAVE-01 |
| 1 | Dashboard | `/` | WAVE-03+ |
| 2 | Training Log (daily) | `/log?date=YYYY-MM-DD` | WAVE-03 |
| 3 | Exercise List | `/exercises` | WAVE-02 |
| 4 | Exercise Detail | `/exercises/:id` | WAVE-02 |
| 5 | Exercise Create | `/exercises/new` (or modal) | WAVE-02 |
| 6 | Exercise Edit | `/exercises/:id/edit` (or modal) | WAVE-02 |
| 7 | Archived Exercises | `/exercises?filter=archived` | WAVE-02 |
| 8 | Body Overview | `/body` | WAVE-04 |
| 9 | Body Check-in Create | `/body/check-in/new` | WAVE-04 |
| 10 | Body Check-in Detail | `/body/check-in/:id` | WAVE-04 |
| 11 | Progress Photo Gallery | `/body/photos` | WAVE-04 |
| 12 | Nutrition Overview | `/nutrition` | WAVE-05 |
| 13 | Product List | `/nutrition/products` | WAVE-05 |
| 14 | Product Create/Edit | `/nutrition/products/new` (or modal) | WAVE-05 |
| 15 | Nutrition Template | `/nutrition/template` | WAVE-05 |
| 16 | Daily Override | `/nutrition/override?date=YYYY-MM-DD` | WAVE-05 |
| 17 | Charts | `/charts` | WAVE-06 |
| 18 | AI Export | `/ai/export` | WAVE-07 |
| 19 | AI Review List | `/ai/reviews` | WAVE-08 |
| 20 | AI Review Create | `/ai/reviews/new` | WAVE-08 |
| 21 | AI Review Detail | `/ai/reviews/:id` | WAVE-08 |
| 22 | Import/Export | `/data` | WAVE-09 |
| 23 | Settings | `/settings` | WAVE-01 |

---

## 8. Detailed Screen Specifications

### 8.1 PIN Unlock Screen

**Purpose**: Authenticate user when PIN is enabled. Full-screen overlay blocking access to all app content.

**Route**: `/pin-unlock` (or full-screen overlay, no route change)

**Primary user actions**:
- Enter PIN digit by digit (numeral input)
- Submit PIN to unlock
- View lockout message if too many failed attempts

**Main content blocks**:
- Centered card on dimmed backdrop
- App logo/name "Atlas" at top
- PIN input field (secured, asterisk/dot display)
- Submit button
- Error message area (generic: "Invalid PIN. Try again.")
- Lockout message if rate-limited

**Form fields**:
- PIN: password input, numeric keyboard on mobile, 4-20 digits

**Empty state**: N/A (always shows input)

**Loading state**: Submit button shows spinner during API call

**Error state**:
- Wrong PIN: generic "Invalid PIN" message below input — clear on next attempt
- Lockout: "Too many attempts. Try again in X minutes."
- Network error: "Unable to connect. Check your server."

**Validation behavior**:
- Numeric only (reject non-digit input silently)
- Length 4-20 — form is submittable only when length >= minimum
- Input clears on wrong PIN (field clears, not page reload)

**Important UX notes**:
- On successful unlock, redirect to Dashboard (or the page user was trying to access)
- On wrong PIN, input clears but does not show lockout details (security)
- If PIN is disabled, this screen is never shown

**Links/transitions**: Dashboard (on success)

**Mobile/responsive**: Centered card fills most of the screen on mobile, comfortable single-hand operation

---

### 8.2 PIN Disabled State / Direct App Access

**Purpose**: When PIN is disabled, no auth gate. App opens directly to Dashboard.

**Behavior**: No PIN screen is shown at any point. All routes and media are accessible directly. The app behaves like a fully open personal tool.

**Session state indicator**: Top bar shows a muted lock-open icon with label "PIN: Off".

**Links/transitions**: Dashboard is the default landing page.

---

### 8.3 Dashboard

**Purpose**: Weekly summary at a glance. Quick entry point to daily actions.

**Route**: `/`

**Primary user actions**:
- Click quick action buttons to navigate to other sections
- View weekly stats

**Main content blocks**:
- Page header: "Dashboard" with today's date subtitle
- **Stat card row**: Last body weight, Training days this week, Cardio sessions this week
- **Goal card**: Current goal from settings
- **Check-in reminder**: Conditional badge if check-in is due
- **Quick actions grid**: 5 buttons — Add Workout, Add Cardio, Add Weight, Weekly Check-in, Generate AI Report

**Empty state** (first launch):
- "Welcome to Atlas! Get started by creating your first exercise."
- Quick action buttons still visible but may link to relevant setup pages

**Loading state**: Skeleton cards for stat row, shimmer for goal and actions

**Error state**: Error message with retry button replacing stat cards. Quick actions may still be functional (they navigate, not load data).

**Important UX notes**:
- Stat values are large and prominent
- Week calculation: Monday-Sunday or configurable
- Check-in reminder: shown if no check-in recorded in the last 7 days
- Quick actions pre-select today's date in target section

**Links/transitions**:
- Add Workout → Training Log (today)
- Add Cardio → Training Log, cardio section open
- Add Weight → Body, new weight entry
- Weekly Check-in → Body, new check-in
- Generate AI Report → AI Export

**Mobile/responsive**: Stat cards stack vertically, quick actions grid becomes 2-column

---

### 8.4 Exercise List

**Purpose**: Browse, search, and manage all exercises.

**Route**: `/exercises`

**Primary user actions**:
- View exercise list/table
- Search by name
- Filter: active / archived / all
- Click row to view detail
- Click "Add Exercise" to create

**Main content blocks**:
- Page header: "Exercises" + "Add Exercise" button
- Filter bar: search input, toggle chips (Active | Archived | All)
- Data table with columns: Name, Muscle Groups, Working Weight, Media, Status, Last Updated
- Row actions: edit (icon), archive/restore (icon with confirmation)

**Columns**:
| Column | Type | Notes |
|--------|------|-------|
| Name | Text | Primary identifier (display), clickable |
| Muscle Groups | Tags/chips | Comma-separated or small tag chips |
| Working Weight | Number | Formatted: "80 kg" or "175 lb" |
| Media | Icon/count | Camera icon + count if media exists |
| Status | Badge | Active (green) / Inactive (gray) |
| Updated | Date | Relative: "2 days ago" |

**Empty state**:
- "No exercises yet"
- "Create your first exercise to start building your library."
- "Add Exercise" primary button

**Loading state**: Skeleton rows (5-8 placeholder lines with shimmer)

**Error state**: Error message: "Failed to load exercises. Retry."

**Important UX notes**:
- Search is client-side for small libraries (filters displayed rows)
- Row click opens detail view (same page inline or separate route)
- Archive action: confirmation dialog "Move [exercise name] to inactive? It can be restored later."
- Duplicate names are allowed — no warning needed on list, but create/edit form may show non-blocking notice
- Pagination: cursor-based, 20 items per page default. "Show more" or infinite scroll.

**Links/transitions**:
- Row click → Exercise Detail (`/exercises/:id`)
- "Add Exercise" → Create form (modal or page)

**Mobile/responsive**: Table becomes list of cards on narrow screens. Each card shows name, muscle groups, working weight, status badge.

---

### 8.5 Exercise Detail

**Purpose**: View full exercise information and manage media.

**Route**: `/exercises/:id`

**Primary user actions**:
- View all exercise fields
- Upload/delete media
- Edit exercise
- Archive/restore exercise

**Main content blocks**:
- Page header: exercise name + status badge
- **Info section**: muscle groups (tags), working weight, description, personal notes
- **Media gallery**: grid of thumbnails with upload dropzone at top
- **Action bar**: Edit button, Archive/Restore button
- **Future placeholder**: "Exercise history (coming in a future update)" — subtle muted note

**Empty state** (no media): Dropzone is prominent, "No media yet. Upload images or video to track form."

**Loading state**: Skeleton for info section, skeleton grid for media

**Error state**: "Exercise not found" with link back to list

**Important UX notes**:
- Status badge is prominent near the name
- Archived exercises: show yellow/amber status badge, all fields still editable
- Edit opens same form as create but prefilled
- No delete — only archive. Text: "Archive this exercise? It will be hidden from the exercise selector but data will be preserved."

**Links/transitions**:
- Edit → Edit form (modal)
- Back arrow/list → Exercise list
- Upload media → triggers file upload flow

**Mobile/responsive**: Stack info vertically, gallery grid adapts columns

---

### 8.6 Exercise Create / Edit Form

**Purpose**: Add new exercise or modify existing one.

**Route**: Modal or page: `/exercises/new` or `/exercises/:id/edit`

**Primary user actions**:
- Fill in exercise fields
- Upload media
- Save or cancel

**Form fields**:
| Field | Type | Required | Validation |
|-------|------|----------|------------|
| Name | Text input | Yes | Trimmed, non-empty after trim. Duplicates allowed with non-blocking warning |
| Muscle Groups | Multi-select or tag input | No | Free text tags (typed, comma/enter-separated) |
| Working Weight | Number input | No | Positive number (>0), empty = no working weight |
| Description | Textarea | No | Multi-line |
| Personal Notes | Textarea | No | Multi-line, distinctly styled as personal/internal |
| Media Upload | Dropzone | No | JPEG/PNG/WEBP/MP4/MOV/WEBM, 25MB image, 250MB video |

**Empty state**: Form is blank for create, prefilled for edit

**Loading state**: Form disabled with spinner during save

**Error state**:
- Field validation errors shown inline
- Server error: toast "Failed to save exercise"
- Network error: toast with retry

**Validation behavior**:
- Name: validate on blur, show "Name is required" if empty
- Working Weight: validate on blur, show "Must be greater than 0" if <= 0
- Duplicate name: on blur or save, query existing names. If duplicate found, show non-blocking warning: "An exercise with this name already exists. You can still save — exercise identity is ID, not name."
- All fields validate on submit

**Important UX notes**:
- Save button: "Save Exercise" (primary, accent)
- Cancel button: "Cancel" (secondary) — unsaved changes warning if form is dirty
- On save success: close modal and refresh list/detail
- Edit mode: prefill all fields, including existing media list

**Links/transitions**:
- Save → List (refresh) or Detail view
- Cancel → List

**Mobile/responsive**: Full-page form on mobile instead of modal

---

### 8.7 Archived / Inactive Exercises View

**Purpose**: View and restore archived exercises.

**Route**: `/exercises?filter=archived` (tab or filter on exercise list)

**Primary user actions**:
- View list of archived exercises
- Restore exercise

**Main content blocks**:
- Same table as active list, filtered to `isActive=false`
- Status column shows "Inactive" badge (gray)
- Restore button per row (or in row actions)

**Empty state**: "No archived exercises."

**Important UX notes**:
- Archived exercises are excluded from all exercise selectors (training log, etc.) by default
- Restore is a one-click action with confirmation: "Restore [name] to active exercises?"
- Archived detail page is still accessible via direct URL

**Links/transitions**:
- Restore → row returns to active list (snackbar: "Exercise restored")
- Detail click → Exercise Detail (same as active)

---

### 8.8 Exercise Media Management (embedded in Detail)

**Purpose**: Upload, view, and delete exercise media.

**Location**: Embedded in Exercise Detail page

**Primary user actions**:
- Upload image/video files
- View media thumbnails
- Open full-size viewer/lightbox
- Delete individual media files

**Main content blocks**:
- Upload dropzone (top of gallery area)
- Thumbnail grid below

**Upload flow**:
1. User drags files or clicks dropzone
2. Client-side validation: file type, file size
3. Show upload progress per file
4. On success: thumbnail appears in gallery
5. On error: inline error message per file

**Delete flow**:
1. Click delete icon on thumbnail
2. Confirmation: "Delete this media?"
3. On confirm: remove from UI, API call to delete
4. On success: toast "Media deleted"

**Important UX notes**:
- Supported types: JPEG, PNG, WEBP, MP4, MOV, WEBM
- Size limits: 25MB images, 250MB video
- Video thumbnails show play overlay icon
- Lightbox viewer: click thumbnail → full-size overlay with prev/next navigation

---

### 8.9 Daily Log Page

**Purpose**: Daily workout entry and management.

**Route**: `/log?date=YYYY-MM-DD`

**Primary user actions**:
- Select date via calendar/date picker
- View exercises for selected date
- Add exercise from library
- Add/edit/delete sets per exercise
- Add exercise comment
- Add cardio entry
- Save changes

**Main content blocks**:
- **Date switcher**: left/right arrows, date display, calendar popover button, "Today" button
- **Exercise list**: ordered list of exercises for the day, each in a card
- **Exercise card**: name (linked to library), working weight display, sets table, comment area, cardio section
- **Sets table**: columns — Set#, Weight, Reps, RPE (optional), RIR (optional), Actions (delete)
- **Add Set button**: below each exercise's set table
- **Add Exercise button**: below exercise list
- **Cardio section**: list of cardio entries for the day
- **Save button**: saves entire day

**Empty state** (no exercises for selected date):
- "No exercises logged for this date."
- "Add Exercise" prominent button
- Cardio section still visible

**Empty state** (date in the future):
- Not applicable — future dates may be disabled or show "Nothing logged yet"

**Loading state**: Skeleton for date switcher, skeleton exercise cards

**Error state**: Error toast on save failure

**Validation behavior**:
- Weight: >= 0
- Reps: positive integer
- RPE: 1-10 (optional)
- RIR: 0-5 (optional)
- If working weight is 0 or empty, user can still enter custom weight per set

**Important UX notes**:
- Working weight auto-populated from Exercise Library when adding exercise
- User can override per-set weight (does not update library working weight)
- "Duplicate exercise" allowed — no restrictions
- Sets are ordered; user can reorder (drag or up/down buttons)
- No workout templates — each day is built from scratch
- No "copy from previous day" in MVP
- Current date is default; past dates accessible via calendar
- Exercise comment: small textarea below exercise name, included in AI export
- Save is per-day (save entire day's changes at once)

**Design assumption**: The daily log is the most complex screen. Sets table should be compact — each row: [set#] [weight input] [reps input] [RPE input] [RIR input] [delete icon]. The set number is auto-incrementing.

**Links/transitions**:
- Exercise name → Exercise Detail (opens in new view, not modal)
- Add Exercise → Exercise selector modal (searchable list)
- Add Cardio → inline cardio form within same page
- Date navigation → loads different date's log

**Mobile/responsive**: Sets table becomes stacked rows on mobile (label:value per field). Sidebar hides. Full-width date selector.

---

### 8.10 Body Overview

**Purpose**: View check-in history, add weight entries, start new check-in.

**Route**: `/body`

**Primary user actions**:
- View check-in history
- Create new check-in
- Add standalone weight entry
- View progress photos gallery

**Main content blocks**:
- Page header: "Body" + "New Check-in" button
- **Stat row**: Last weight, Last check-in date, Days since last check-in
- **Quick actions**: "Add Weight", "New Check-in"
- **Check-in history**: list/table of past check-ins — date, weight, body fat %, photo count, view button
- **Weight entries**: optional mini-list of recent standalone weight entries

**Empty state**: "No check-ins yet. Complete your first weekly check-in."

**Loading state**: Skeleton stat cards, skeleton history list

**Error state**: Error message with retry

**Links/transitions**:
- "New Check-in" → Body Check-in Create
- Check-in row → Body Check-in Detail
- "View Photos" → Progress Photo Gallery

---

### 8.11 Body Check-in Create / Edit

**Purpose**: Record weekly body measurements.

**Route**: `/body/check-in/new` or `/body/check-in/:id/edit`

**Primary user actions**:
- Enter date (default today)
- Enter weight (optional)
- Enter body fat % (optional)
- Enter 10 measurements
- Upload 2-4 progress photos
- Add notes
- Save

**Form fields**:
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| Date | Date picker | Yes | Default today |
| Weight | Number | No | kg or lb |
| Body Fat % | Number | No | Percentage |
| Neck | Number | No | cm or inches |
| Shoulders | Number | No | |
| Forearm | (paired) Number ×2 | No | Left / Right inputs, or single |
| Biceps | (paired) Number ×2 | No | Left / Right or single |
| Chest | Number | No | |
| Waist | Number | No | |
| Abdomen | Number | No | |
| Hips | Number | No | |
| Thigh | (paired) Number ×2 | No | Left / Right or single |
| Calf | (paired) Number ×2 | No | Left / Right or single |
| Photos | File upload (2-4) | No | Progress photos |
| Notes | Textarea | No | |

**Empty state**: Form blank for create, prefilled for edit

**Loading state**: Form disabled with spinner during save

**Validation behavior**:
- All measurement values: positive numbers
- Paired measurements: if user enters left, right is optional. If user enters only "common" field, it's treated as both sides equal
- Photos: 2-4 recommended but not strictly enforced (design assumption: accept 0-10, show recommendation text)
- Date: must be valid, should not be in the future (subtle warning if future)

**Important UX notes**:
- Paired measurements: show two input fields side by side with "Left" / "Right" labels, plus a "Single Value" option that fills both
- Measurements in a grid layout: 3 columns on desktop, 2 on tablet, 1 on mobile
- Photo upload: same dropzone pattern as Exercise Media
- Photos can be added multiple times before save

**Links/transitions**: Save -> Check-in Detail or Overview

---

### 8.12 Check-in Detail

**Purpose**: View a single check-in with all measurements and photos.

**Route**: `/body/check-in/:id`

**Primary user actions**:
- View all measurements
- View/edit notes
- View photos in gallery
- Delete check-in

**Main content blocks**:
- Check-in header: date, weight, body fat %
- **Measurement grid**: all 10 measurements displayed as labeled stat cards
- **Photo gallery**: similar to Exercise Media gallery
- **Notes**: displayed text
- **Actions**: Edit button, Delete button

**Empty state**: N/A (check-in exists by definition)

**Loading state**: Skeleton

**Error state**: "Check-in not found" with link back to overview

**Important UX notes**:
- Delete: requires confirmation dialog. "Delete this check-in? All measurements and photos will be permanently deleted."
- Delete is actual delete (not soft)
- Edit opens the same form as create but prefilled

---

### 8.13 Progress Photo Gallery (Body)

**Purpose**: Browse all progress photos across check-ins.

**Route**: `/body/photos` (or tab within Body section)

**Primary user actions**:
- View all photos grouped by check-in date
- Click photo to view full-size
- Delete photo

**Main content blocks**:
- Grid of photo thumbnails with date label
- Grouped by check-in (date header above each group)
- Click opens lightbox viewer
- Delete icon per photo

**Empty state**: "No progress photos yet. Add photos during your weekly check-in."

**Important UX notes**:
- Photos are always associated with a check-in — cannot be standalone
- Photo angles displayed as small badges: Front, Side, Back, Custom

---

### 8.14 Nutrition Overview

**Purpose**: View weekly nutrition template and daily overrides.

**Route**: `/nutrition`

**Primary user actions**:
- View product list
- View/create weekly template
- Override specific day
- View macro summary

**Main content blocks**:
- Page header: "Nutrition" with sub-tabs or section links
- **Product section**: list of products, "Add Product" button
- **Template section**: if template exists, show week's plan. If not, "Create Template" button
- **Macro summary card**: daily/weekly totals: calories, protein, fat, carbs
- **Day selector**: click a day of the current week to view override

**Empty state**: "No products yet. Create your first food product to build your nutrition plan."

**Loading state**: Skeleton

**Error state**: Error toast with retry

**Links/transitions**:
- Product click → Product Edit
- "Add Product" → Product Form
- "Create Template" → Template Editor
- Day click → Daily Override

---

### 8.15 Product List / Create / Edit

**Purpose**: Manage food product catalog.

**Route**: `/nutrition/products`

**Primary user actions**:
- View product list
- Add new product
- Edit existing product
- Delete product

**Form fields** (Create/Edit):
| Field | Type | Required | Validation |
|-------|------|----------|------------|
| Name | Text | Yes | |
| Calories / 100g | Number | Yes | >= 0 |
| Protein / 100g | Number | Yes | >= 0 |
| Fat / 100g | Number | Yes | >= 0 |
| Carbs / 100g | Number | Yes | >= 0 |
| Notes | Textarea | No | |

**Important UX notes**:
- 4 macro fields side by side in a row (Calories, Protein, Fat, Carbs)
- Values are per 100g (display label: "per 100g")
- Delete with confirmation: "Delete [product name]? This may affect existing templates."

---

### 8.16 Weekly Nutrition Template

**Purpose**: Create and edit the weekly meal plan.

**Route**: `/nutrition/template`

**Primary user actions**:
- Add products with gram amounts
- Assign meal labels (optional)
- View calculated macro totals
- Save template

**Main content blocks**:
- Template header: week start date selector, template title
- **Items table**: columns — Product, Amount (g), Meal Label (optional), Actions (delete)
- **Add Item row/button**: opens product selector
- **Macro totals row**: sum of all items — Calories, Protein, Fat, Carbs

**Empty state**: "No template yet. Create a weekly nutrition plan by adding products."

**Important UX notes**:
- Only one active template at a time
- Template auto-applies to all 7 days
- Items can be reordered (drag handles)
- Meal labels: dropdown (Breakfast, Lunch, Dinner, Snack, Pre-workout, Post-workout) or free text

**Design assumption**: product selector is a searchable dropdown or modal showing product name + macros per 100g. After selecting, user enters grams.

---

### 8.17 Daily Nutrition Override

**Purpose**: Override nutrition for a specific day.

**Route**: `/nutrition/override?date=YYYY-MM-DD`

**Primary user actions**:
- View template-calculated macros for the day
- Add/subtract/replace products
- View recalculated macros

**Main content blocks**:
- Date display
- **Template items**: shown as read-only list (what the template provides)
- **Override items**: editable table — Product, Operation (add/subtract/replace), Amount (g), Meal Label
- **Add override item button**
- **Recalculated macro totals**: template ± overrides

**Empty state**: No overrides yet — shows template values only with "Add override" button

**Important UX notes**:
- Override operations: "add" (extra), "subtract" (reduce), "replace" (swap product)
- Template values are not editable here (edit template instead)
- Changes affect only the selected date

---

### 8.18 Charts

**Purpose**: Visualize all progress metrics.

**Route**: `/charts`

**Primary user actions**:
- Select chart category (Training, Body, Nutrition)
- Select specific entity (exercise, measurement type)
- Select date range
- View chart with interactive tooltips

**Main content blocks**:
- **Category tabs**: Training | Body | Nutrition
- **Entity selector**: dropdown (exercises for Training, measurements for Body, macro for Nutrition)
- **Date range picker**: preset buttons (4 weeks, 3 months, 6 months, 1 year, All) or custom range
- **Chart card(s)**: one or multiple charts for the selected entity/metric
- **Metric selector** (Training only): Working Weight, Best Set, e1RM, Volume, Total Reps, Working Sets

**Empty state**: "No data to display for this period."

**Loading state**: Spinner or shimmer in chart area

**Error state**: "Failed to load chart data. Retry."

**Important UX notes**:
- Charts should have clear labels, grid lines, and interactive tooltips
- Chart type: line chart for trends, bar chart for weekly comparisons
- Multiple metrics can be shown on one chart (e.g., working weight + volume overlay)
- Body charts support multi-measurement overlay (select multiple measurements to compare)
- Nutrition chart shows weekly average bars for calories, protein, fat, carbs

**Links/transitions**: None (self-contained page)

---

### 8.19 AI Export Builder

**Purpose**: Configure and generate AI export package.

**Route**: `/ai/export`

**Primary user actions**:
- Select date range
- Toggle data sections to include
- Review persistent AI context
- Add one-time comment
- Select week flags
- Generate export
- View/download prompt
- Copy prompt to clipboard
- Download ZIP

**Main content blocks**:
- **Date range**: preset (4 weeks default) or custom
- **Persistent context**: read-only display of settings (goal, height, age, experience, etc.)
- **Section toggles**: checkboxes — Workouts, Exercises, Sets, Comments, RPE/RIR, Cardio, Body Weight, Measurements, Photos (off by default), Nutrition
- **Week flags**: chips/buttons — Poor Sleep, High Stress, Illness, Injury/Pain, AAS/Cycle, Calorie Deficit, Surplus, Maintenance, Missed Workouts, Travel/Disrupted
- **One-time comment**: textarea
- **Generate button**: primary action
- **After generation**:
  - **Prompt preview**: code/read-only block with copy button
  - **Download button**: downloads ZIP
  - **Status**: file size, included sections summary

**Empty state**: "No data in selected period." warning if period has no logged data

**Loading state**: "Generating export..." spinner with progress indication

**Error state**: "Export failed." toast with retry

**Important UX notes**:
- Generation may take several seconds (especially with photos) — show progress
- Photos toggle has a note: "Photos will be included in ZIP. File size may be large."
- Prompt preview is scrollable/copyable
- Section toggles are all on by default except Photos

**Links/transitions**: None (self-contained)

---

### 8.20 AI Review List

**Purpose**: View saved AI analysis responses.

**Route**: `/ai/reviews`

**Primary user actions**:
- View list of saved reviews
- Click to view detail
- Delete review

**Main content blocks**:
- Page header: "AI Reviews" + "New Review" button
- **Review list**: cards/table — date range, preview of response text (truncated), date saved, planned actions count

**Empty state**: "No AI reviews saved yet. After receiving AI analysis, save it here."

**Loading state**: Skeleton list

**Error state**: Error message with retry

**Links/transitions**:
- Row click → AI Review Detail
- "New Review" → AI Review Create

---

### 8.21 AI Review Create

**Purpose**: Save an AI analysis response.

**Route**: `/ai/reviews/new`

**Primary user actions**:
- Enter date range
- Paste AI response text
- Add user notes
- Add planned actions
- Save

**Form fields**:
| Field | Type | Required |
|-------|------|----------|
| Date Range Start | Date picker | Yes |
| Date Range End | Date picker | Yes |
| AI Response | Textarea (large) | Yes |
| User Notes | Textarea | No |
| Planned Actions | Textarea | No |

**Important UX notes**:
- Date range should generally match the AI export date range
- AI Response textarea should be generous in size (paste large responses)

---

### 8.22 AI Review Detail

**Purpose**: View a saved review.

**Route**: `/ai/reviews/:id`

**Primary user actions**:
- View full AI response text
- View notes and planned actions
- Edit
- Delete

**Main content blocks**:
- Header: date range
- **AI response**: scrollable text block (monospace or styled quote)
- **Notes**: displayed text
- **Planned actions**: displayed text or checklist
- **Actions**: Edit, Delete

---

### 8.23 Import / Export

**Purpose**: Full backup and restore.

**Route**: `/data`

**Primary user actions**:
- Export full backup
- Toggle media inclusion
- Upload backup file for import
- Run dry-run validation
- Confirm import

**Main content blocks**:
- **Export section**:
  - "Export All Data" button
  - "Include media" toggle (default: on)
  - Loading spinner during generation
  - Download link when ready
- **Import section**:
  - File upload dropzone (accepts .zip)
  - "Validate" button
  - **Validation summary**: table of what will be restored — entity counts (exercises: 45, workouts: 230, media: 12 files)
  - **Import button**: appears only after successful validation
  - "Destructive replace" warning if data exists

**Empty state**: No backups yet — "Create your first backup to safeguard your data."

**Loading state**: Spinner during export generation, spinner during validation, progress bar during import

**Error state**:
- Invalid ZIP format: "The uploaded file is not a valid backup."
- Schema mismatch: "This backup was created by a different app version. Cannot import."
- Validation failure: detailed error message

**Important UX notes**:
- Import is destructive — it replaces all existing data
- Warning dialog on import: "This will replace ALL existing data. Are you sure?"
- Progress indication for long-running imports
- Clear success/failure result at the end

---

### 8.24 Settings

**Purpose**: Configure app behavior.

**Route**: `/settings`

**Primary user actions**:
- Enable/disable PIN
- Change PIN
- Edit AI context
- Change units
- Change default export weeks

**Main content blocks**:
- **PIN section**:
  - Toggle: Enable/Disable PIN (with confirmation if disabling)
  - Change PIN: requires current PIN + new PIN (×2 confirm)
  - Session status: "Session active" or "No active session"
  - "Lock Session" button
- **AI Context section**:
  - Free-text fields: Goal, Height, Age (optional), Training Experience, Split, Limitations, Progression Style, Nutrition Strategy
  - Textarea: Persistent AI Comment
  - Save button
- **Preferences section**:
  - Units: dropdown (Metric / Imperial)
  - Default Export Weeks: number input (default 4)

**Empty state**: PIN section shows "PIN not configured. Enable to protect your data."

**Loading state**: Form disabled during save, spinner on toggle

**Validation behavior**:
- PIN: 4-20 digits only
- Change PIN: current PIN must match, new PIN must match confirmation
- Height: positive number
- Age: positive integer

**Important UX notes**:
- Disabling PIN: confirmation "Are you sure? All data will be accessible without a PIN."
- PIN change: current PIN, new PIN, confirm new PIN
- Session lock: one-click, immediate, shows toast "Session locked"
- AI context fields are text inputs — no validation beyond basic length

---

## 9. Exercise Library UX Details

### 9.1 Exercise List - Detailed UX

The exercise list is the primary interface for managing the exercise library.

**View modes**: Table view (default desktop). Card/grid view optional for small screens.

**Search behavior**:
- Search input at top of filter bar
- Filters as user types (live filtering with debounce ~300ms)
- Searches by name (case-insensitive contain match)
- No full-text search in MVP — basic ILIKE/patter matching

**Filter chips**:
- Three mutually exclusive chips: Active (default) | Archived | All
- Active filter: shows `isActive=true` exercises
- Archived filter: shows `isActive=false` exercises
- All filter: shows both

**Pagination**:
- Cursor-based pagination, 20 items per page
- "Load more" button at bottom, or scroll-based infinite load
- Total count displayed: "Showing 20 of 45 exercises"

**Row actions**:
- Row click: navigate to detail
- Edit icon (pencil): opens edit modal directly from list
- Archive/Restore icon: one-click with confirmation toast
- Confirmation dialog for archive: "Archive [name]? It will be hidden from exercise selector."

**Working weight display**:
- Formatted with unit: "80 kg" or "175 lb"
- If no working weight: "—" (dash)

**Media count**:
- Small camera icon + number
- If 0: no icon shown (empty cell)

**Status badge**:
- Active: small green dot or "Active" label
- Inactive: gray "Inactive" label
- Badge next to name or in dedicated column

**Sorting**:
- Default sort: name (ascending, alphabetical)
- Clickable column headers for sort: Name, Updated At
- Single column sort at a time

### 9.2 Create/Edit Exercise - Detailed UX

**Mode detection**:
- Create: form opens with empty fields, title "New Exercise"
- Edit: form opens prefilled, title "Edit Exercise"

**Name field**:
- Text input, single line
- Required indicator (asterisk)
- On blur: trim whitespace, validate non-empty
- Duplicate warning: if user types a name that already exists, show non-blocking warning below field: "An exercise with this name already exists. Duplicate names are allowed."
- Warning is a yellow/info banner, not blocking. Save button remains enabled.

**Muscle groups**:
- Tag-style input: user types a muscle group, presses Enter or comma to create a tag
- Tags are displayed as small chips with × to remove
- Examples: Chest, Back, Shoulders, Biceps, Triceps, Quads, Hamstrings, Glutes, Abs, Calves, Forearms, Traps, Lats
- Free text — no predefined list (user can type anything)
- Empty by default

**Working weight**:
- Number input, optional
- Label with unit indicator: "Working Weight (kg)" or "(lbs)"
- Validation: if filled, must be > 0
- Step increment: 0.5 or 1 (depending on units)
- Helper text: "This weight will prefill when adding to workout"

**Description**:
- Textarea, multi-line
- No rich text
- Placeholder: "Technique notes, equipment needed, etc."

**Personal Notes**:
- Textarea, multi-line
- Visually distinct from description (different border color or icon: lock/eye-off icon)
- Placeholder: "Your private notes (not shown in AI export)"
- This field is excluded from AI export

**Media upload**:
- Dropzone at bottom of form
- Same constraints as detail page media upload
- Upload during create: media is uploaded and associated when exercise is saved
- Upload during edit: immediate upload and association
- Existing media shown as thumbnails below dropzone

**Save flow**:
- Click "Save Exercise"
- Client-side validation: name required
- Server-side validation: name non-empty, working weight > 0 if provided
- On success: close form, refresh list/detail, toast "Exercise created" or "Exercise updated"
- On error: show server error in toast, keep form open

### 9.3 Exercise Detail - Detailed UX

**Layout**:
- Top section: name (large), status badge (right), edit button (right)
- Info card: working weight, muscle groups (tags), created/updated dates
- Description card: full description text
- Personal Notes card: visually distinct (different background or border, lock icon)
- Media gallery section
- Archive/Restore button at bottom
- Future placeholder: subtle muted text "Exercise history and performance charts coming in a future update"

**Status badge behaviors**:
- Active: green badge "Active"
- Inactive: gray badge "Inactive"

**Archive button**:
- When active: "Archive Exercise" button (secondary/destructive styling)
- Confirmation: "Archive [name]? It will be hidden from the exercise selector. Media will be preserved."
- After archive: badge changes to "Inactive", button changes to "Restore"

**Restore button**:
- When inactive: "Restore Exercise" button (primary styling)
- Confirmation: "Restore [name] to active exercises?"
- After restore: badge changes to "Active", button changes to "Archive"

### 9.4 Archive/Restore - Detailed UX

**Archive behavior**:
- `isActive` set to `false`
- Exercise removed from default lists and selectors
- Exercise still accessible via direct URL or "Archived" filter
- Media remains intact and queryable
- Existing workout history referencing this exercise remains unchanged
- Toast: "[Name] archived. Undo?" with undo option (5-second window)

**Restore behavior**:
- `isActive` set to `true`
- Exercise returns to default lists and selectors
- All media and history preserved
- Toast: "[Name] restored."

### 9.5 Media Upload - Detailed UX

**Upload flow**:
1. User drops files or clicks dropzone area
2. Client validates file type and size per file
3. Each valid file shows a preview thumbnail with name and size
4. Invalid files show error message (red border, "Unsupported format" or "File too large")
5. Upload button (or auto-upload on drop)
6. Progress bar per file during upload
7. Success: thumbnail appears in gallery
8. Failure: error on specific file, retry button

**Delete flow**:
1. Click trash icon on thumbnail
2. Confirmation dialog: "Delete this media?"
3. Confirm: API call, remove from UI
4. Success toast: "Media deleted"

**Gallery view**:
- Grid layout: 3-4 columns on desktop, 2 on tablet, 1 on mobile
- Thumbnails: 150-200px with consistent aspect ratio
- Images: show actual thumbnail
- Videos: show first frame with play icon overlay
- Hover: subtle overlay with file name, delete button

**Lightbox**:
- Click thumbnail → full-screen overlay
- Image: centered, max dimensions respecting aspect ratio
- Video: embedded player with controls
- Close button (X), click outside to close
- Keyboard: Escape to close, arrow keys to navigate multiple media

### 9.6 Empty Libraries / Zero State

When user has no exercises:
- Exercise List shows empty state: "No exercises yet"
- Subtitle: "Create your first exercise to build your training library."
- Primary "Add Exercise" button
- Simple illustration (e.g., dumbbell outline)
- Training Log's exercise selector also shows empty state: "No exercises in library. Create one first." with link to Exercise Library

### 9.7 Duplicate Name UX Note

Because duplicate names are allowed:
- List shows all exercises with the same name — differentiate by creation date or working weight
- No uniqueness validation on name
- Non-blocking warning on create/edit if duplicate detected
- Identity is always `exerciseId`, never name
- In workout diary, exercise name is displayed but backed by `exerciseId`

---

## 10. Daily Log UX Details

**Date navigation**:
- Default: today
- Arrows: ← Previous Day | Next Day → (disabled if future)
- Calendar popover: click date display to open month calendar, click date to navigate
- "Today" button: return to current date
- Date with existing log data: subtle dot indicator on calendar

**Adding exercises**:
- "Add Exercise" button opens a searchable modal
- Modal shows: search input, filtered list of active exercises (name + muscle groups)
- Click exercise to add to day
- Exercise is appended to end of list
- Working weight auto-populated from Exercise Library's `workingWeight`
- No duplicate prevention — same exercise can be added multiple times

**Sets table** (per exercise):
- Columns: #, Weight, Reps, RPE (opt), RIR (opt), [delete]
- New set: last row with empty fields
- Auto-incrementing set number
- Tab through fields for rapid entry: Weight → Reps → RPE → RIR → next row
- Delete set: trash icon, no confirmation (undo via toast if needed)

**Progression signals** (future display hints):
- Currently metadata for AI export; when implemented visually, show small indicators:
  - Green up arrow: weight/volume increasing
  - Yellow dash: stable
  - Red down arrow: decreasing/regression

**Save behavior**:
- "Save" button in page header or bottom of page
- Saves entire day: exercises, sets, comments, cardio
- On save: toast "Saved" with timestamp
- No autosave (design assumption — user clicks Save explicitly)

---

## 11. Body Tracking UX Details

**Weekly check-in flow**:
1. Navigate to Body → click "New Check-in"
2. Date defaults to today
3. Enter weight (optional, numeric)
4. Enter body fat % (optional, numeric, percentage)
5. **Measurement grid**:
   - 4 columns on desktop: Measurement Name | Value | (Left) | (Right)
   - Paired measurements (forearm, biceps, thigh, calf): show L/R columns
   - Unpaired (neck, shoulders, chest, waist, abdomen, hips): single value column
   - Measurements displayed in logical order: neck → shoulders → chest → waist → abdomen → hips → forearms → biceps → thigh → calf
   - Measurement labels should be clear, consistent
6. **Photo upload**: dropzone, 2-4 recommended
7. **Photo angles**: after upload, user can tag angle — Front, Side, Back, Custom
8. **Notes**: optional textarea
9. **Save button**

**Standalone weight entry**:
- Quick form on Body overview or Dashboard
- Date + weight only
- Source: "manual" (implied)

**Progress photos grid**:
- Thumbnails grouped by check-in date
- Click opens lightbox
- Angle label on thumbnail corner
- Delete per photo with confirmation

---

## 12. Nutrition UX Details

**Product database**:
- Simple list with search
- Each row: Name, Calories, Protein, Fat, Carbs (per 100g)
- Click row to edit, swipe/icon to delete
- Create form: inline or modal

**Weekly template**:
- Shows week start date (default Monday of current week)
- Items table: Product | Amount (g) | Meal Label | Actions
- "Add Item" → product selector modal
- Product selector: searchable list with macros per 100g
- After selecting product, user enters grams and optional meal label
- Macro totals calculated and displayed below table
- Only one active template

**Daily override**:
- Click specific day in Nutrition overview
- Shows template values (read-only) + override items (editable)
- Add override item: select product, choose operation (add/subtract/replace), enter grams
- Recalculated macros displayed

**Macro summary**:
- Card showing: Calories, Protein (g), Fat (g), Carbs (g)
- Color-coded or bar indicators: protein in blue, fat in amber, carbs in purple (design choice)
- Shows both absolute values and percentage breakdown

---

## 13. Charts UX Details

**Chart types per category**:

**Training charts**:
- Metric tabs: Working Weight | Best Set | e1RM | Volume | Total Reps | Working Sets
- Exercise selector: dropdown of all exercises with data in period
- Date range selector
- Line chart for metric over time
- Each data point corresponds to a workout session

**Body charts**:
- Metric tabs: Weight | Body Fat % | Measurements
- Measurements: multi-select (checkboxes) to overlay multiple measurements
- Date range selector
- Line chart for body weight, individual measurement trends
- Overlay mode: each measurement in different color

**Nutrition charts**:
- Metric tabs: Calories | Protein | Fat | Carbs
- Bar chart showing weekly averages
- Date range selector

**Common chart controls**:
- Date range presets: 4w, 3m, 6m, 1y, All
- Custom date range: from/to date pickers
- Responsive: chart resizes with window
- Empty: "No data for this period" message

---

## 14. AI Export / AI Review UX Details

### AI Flow (Weekly)

1. User navigates to AI Export
2. Default date range: last 4 weeks (configurable in Settings)
3. User adjusts date range if needed
4. User reviews persistent context (read-only, editable in Settings)
5. User toggles sections on/off
6. User toggles "Include photos" (off by default) — warning about file size
7. User selects week flags (chips/tags)
8. User adds optional one-time comment
9. User clicks "Generate Export"
10. System shows progress: "Gathering data...", "Building prompt...", "Creating ZIP..."
11. On completion:
    - Prompt preview shown in a styled code block
    - "Copy Prompt" button
    - "Download ZIP" button
    - Summary: date range, sections included, file size

### Prompt Preview
- Monospace font in a scrollable container
- Copy button copies full prompt text to clipboard
- No editing — prompt is generated

### Export Package
- Downloaded as `.zip` file
- Filename: `atlas-export-YYYY-MM-DD_YYYY-MM-DD.zip`

### AI Review
- Create: form with date range, AI response textarea (large), notes, planned actions
- List: cards with date range, response preview (truncated), date saved
- Detail: full response, notes, planned actions

---

## 15. Import / Export UX Details

### Export
- Section title: "Export Data"
- "Export All Data" button
- Toggle: "Include media files" (default: on)
- On click: confirmation "Create a full backup of all your data?"
- Progress indicator during generation
- Download link appears when ready
- File name: `atlas-backup-YYYY-MM-DD.zip`

### Import
- Section title: "Import Data"
- Warning banner: "Importing will replace ALL existing data. This cannot be undone."
- File upload dropzone: accepts .zip files only
- After upload: "Validate" button
- Validation runs dry-run
- Summary displays: entity counts, media count, schema version
- "Start Import" button (enabled only after successful validation)
- Second confirmation: "Are you sure? This will replace all current data."
- Progress during import
- Result: success message or detailed error

### Error scenarios
- Invalid ZIP: error message "Not a valid backup file"
- Wrong schema version: "Backup incompatible with current app version"

---

## 16. Settings UX Details

Settings page is organized into sections with clear headings:

### PIN Security
- Toggle: PIN protection (On/Off)
- When OFF: explanation text "Enable PIN to protect your data with a numeric code."
- When ON:
  - "Change PIN" button → opens form: Current PIN, New PIN, Confirm PIN
  - Session status: green "Active" or gray "No active session"
  - "Lock Session" button (ends current session immediately)

### AI Context
- Text fields:
  - Goal (textarea): "e.g., Build muscle, improve strength"
  - Height (number)
  - Age (number, optional)
  - Training Experience (text): "e.g., 3 years"
  - Current Split (text, optional): "e.g., PPL, Upper/Lower"
  - Limitations (textarea): "e.g., Shoulder injury, knee issues"
  - Progression Style (text): "e.g., Double progression, 5/3/1"
  - Nutrition Strategy (text): "e.g., Maingaining, Lean bulk"
  - Persistent Comment (textarea): anything the user wants AI to always consider
- "Save" button

### Preferences
- Units: dropdown (Metric / Imperial)
- Default AI Export Weeks: number input (1-52, default 4)

### About
- App version (from build)
- Link to docs/GitHub (placeholder)

---

## 17. UI State Matrix

| Page | Loading | Empty | Error | Validation | Success | Unauthorized |
|------|---------|-------|-------|------------|---------|--------------|
| PIN Unlock | Spinner on submit button | N/A | Generic "Invalid PIN" or lockout message | Numeric only, length 4-20 | Redirect to Dashboard | N/A (auth page) |
| Dashboard | Skeleton stat cards | Welcome + CTA | Error card + retry | N/A | Data displayed | PIN overlay |
| Exercise List | Skeleton rows | "No exercises" + CTA | Error + retry | N/A | Table with rows | PIN overlay |
| Exercise Form | Disabled form + spinner | Blank form (create) | Server error toast | Inline field errors | Close + refresh toast | PIN overlay |
| Exercise Detail | Skeleton blocks | N/A (exists) | "Not found" | N/A | Data displayed | PIN overlay |
| Exercise Media | Upload progress per file | Dropzone visible | Per-file error | Type/size rejection | Thumbnail in gallery | PIN overlay |
| Daily Log | Skeleton exercises | "Add first exercise" | Save error toast | Set field errors | "Saved" toast | PIN overlay |
| Body Check-in | Skeleton form | Blank form | Save error toast | Field validation | Close + refresh | PIN overlay |
| Nutrition Template | Skeleton table | "Create template" | Save error toast | Gram amount >0 | "Saved" toast | PIN overlay |
| Charts | Spinner in chart area | "No data" message | Error + retry | N/A | Chart rendered | PIN overlay |
| AI Export | Progress spinner | "No data in period" | Generation error | N/A | Prompt + download | PIN overlay |
| Import/Export | Export/import progress | "No backup yet" | Validation error | ZIP validity | Success summary | PIN overlay |
| Settings | Form loading | Defaults shown | Save error toast | PIN field validation | "Saved" toast | PIN overlay (settings is guarded) |
| AI Review List | Skeleton cards | "No reviews" | Error + retry | N/A | Cards/list | PIN overlay |

---

## 18. Design Constraints for OpenDesign

- **No marketing landing page** as main app. Atlas has no public-facing website within the app.
- **No social/community features**. No share buttons, no public profiles, no comments from others.
- **No coach/admin multi-user UI**. No role switching, no user management, no admin panels.
- **No payment/subscription UI**. No upgrade prompts, no pricing pages, no license keys.
- **No workout template planner** in MVP. Do not design a "repeat workout" or "copy from last week" feature.
- **No recipe/barcode scanner UI** in MVP. Do not design barcode scanning or recipe creation flows.
- **No Apple Health integration UI** as implemented feature. Do not show health sync toggles or data import from Health.
- **No Telegram bot UI** — explicitly out of scope.
- **No OpenAI/API key configuration UI** — AI analysis is manual copy-paste, not in-app API call.
- **Keep app personal, private, data-focused**. Avoid social, competitive, or gamification elements.
- **No progress predictions or recommendations** generated by the app itself. Charts show historical data only.
- **No "quick add" or "start workout" timer** — the app is for retrospective logging, not real-time coaching.
- **No split/workout plan designer** — the user's training split is a text field in settings, not a structured planner.

---

## 19. Suggested Component List

| Component | Purpose |
|-----------|---------|
| **AppShell** | Top-level layout wrapper: sidebar + top bar + content area. Manages PIN state. |
| **SidebarNav** | Left navigation with icon + label per section. Active state highlight. Collapsible on small screens. |
| **TopBar** | Top bar with app name, session status indicator. |
| **PageHeader** | Page title + optional primary action button + optional secondary actions. |
| **StatCard** | Single metric display: label, large value, optional trend indicator. Used on Dashboard and Body overview. |
| **DateSwitcher** | Horizontal date navigation: left arrow, date label, right arrow, "Today" button. Used on Daily Log. |
| **CalendarPopover** | Dropdown month calendar for date selection. Supports single date and date range modes. |
| **DataTable** | Sortable column table with optional row actions. Used on Exercise List, Product List, Check-in History. |
| **EntityCard** | Card for an entity with key fields. Used as table row alternative on mobile. |
| **ExerciseForm** | Create/Edit exercise form with all fields, validation, duplicate name warning, media upload. |
| **ExerciseMediaGallery** | Thumbnail grid with upload dropzone, delete per item, lightbox viewer. |
| **MediaUploadDropzone** | File dropzone with validation, progress, preview. Reused for exercise media and progress photos. |
| **SetEditor** | Sets table: weight, reps, RPE, RIR, set number. Add/delete rows. Tab-friendly. |
| **CardioForm** | Cardio entry form: type dropdown, duration, pulse, zone selector. |
| **BodyCheckInForm** | Full check-in form: date, weight, body fat, measurement grid, photo upload, notes. |
| **MeasurementGrid** | 10 measurement fields in grid layout. Paired L/R for forearm, biceps, thigh, calf. |
| **ProgressPhotoGrid** | Photo thumbnails grouped by check-in, with angle badges, lightbox, delete. |
| **NutritionProductForm** | Product create/edit: name, KJBJU per 100g fields. |
| **NutritionTemplateTable** | Template items table: product, grams, meal label, actions. Macro totals row. |
| **MacroSummaryCard** | KJBJU display with colored bars or values. |
| **ChartCard** | Chart container with title, chart visualization, optional controls. |
| **AIExportBuilder** | Full AI export configuration panel: date range, section toggles, context display, week flags, generate button. |
| **PromptPreview** | Read-only styled code block with copy button. |
| **BackupImportPanel** | Import section: file upload, validate button, summary display, confirm button. |
| **PinUnlockForm** | Centered PIN input card on dimmed backdrop. |
| **SettingsPanel** | Settings sections: PIN, AI context, preferences. |
| **EmptyState** | Illustration + heading + description + action button. |
| **ErrorState** | Error message + optional retry button. |
| **LoadingState** | Skeleton/shimmer placeholder. |
| **ConfirmDialog** | Modal dialog: title, description, Cancel + Confirm buttons. Support destructive variant. |
| **Toast** | Brief notification: success message, optional undo action. Auto-dismiss. |

---

## 20. Prompt to OpenDesign

```
# Prompt to OpenDesign

Generate frontend design for Atlas — a self-hosted personal fitness tracking web application.

## Design Requirements

1. **App type**: Personal analytics dashboard + training log. Single user, private, data-focused.
2. **Visual style**: Clean, calm, modern. Muted palette with accent color for interactions. Desktop-first. Dark and light mode compatible.
3. **Layout**: Fixed left sidebar (icon + label nav), thin top bar, main content area. Responsive — sidebar collapses on tablet/mobile.
4. **Core workflow**: User logs workouts daily, does a weekly body check-in, then generates an AI analysis package. Design should support this rhythm.

## Screens to generate (key screens first)

1. **Dashboard** — Weekly summary with stat cards (weight, training days, cardio sessions), current goal, check-in reminder, quick action buttons.
2. **Exercise List** — Searchable/filterable table with name, muscle groups, working weight, media count, active/inactive status. "Add Exercise" button.
3. **Exercise Detail** — Full exercise info, status badge, media gallery with upload dropzone, archive/restore action.
4. **Exercise Form** — Create/edit with name, muscle groups (tag input), working weight, description, personal notes, media upload.
5. **Daily Log** — Date switcher, exercise cards with compact sets table (weight/reps/RPE/RIR), "Add Exercise" button, save action. This is the most complex screen.
6. **Body Check-in** — Weight, body fat %, 10 measurements in grid (paired L/R for some), photo upload, notes.
7. **Nutrition Template** — Products table with grams, macro totals, product selector modal.
8. **Charts** — Category tabs (Training/Body/Nutrition), entity selector, date range, line/bar charts with tooltips.
9. **AI Export** — Date range, section toggles, persistent context, week flags, generate button, prompt preview, download.
10. **PIN Unlock** — Centered PIN input on blurred dark backdrop.
11. **Settings** — PIN toggle/change, AI context form fields, preferences (units, export weeks).
12. **Import/Export** — Export with media toggle, import dropzone with validation summary.

## Common patterns to design

- **Cards**: stat cards, entity cards, chart cards with consistent elevation
- **Data tables**: sortable columns, row click, row actions
- **Forms**: compact, inline validation, tab-friendly (for data entry speed)
- **Media gallery**: thumbnail grid, upload dropzone, lightbox viewer
- **Empty states**: illustration + heading + CTA button per section
- **Loading states**: skeleton shimmers
- **Error states**: message + retry
- **Confirmation dialogs**: for destructive actions
- **Toasts**: success with optional undo

## Component library building blocks

Design the following core components:
- AppShell (sidebar + top bar + content)
- SidebarNav
- PageHeader
- StatCard, EntityCard, ChartCard
- DataTable (sortable, with row actions)
- DateSwitcher + CalendarPopover
- ExerciseForm, SetEditor, CardioForm
- BodyCheckInForm, MeasurementGrid
- MediaUploadDropzone, ExerciseMediaGallery
- NutritionTemplateTable, MacroSummaryCard
- AIExportBuilder, PromptPreview
- BackupImportPanel
- PinUnlockForm
- EmptyState, ErrorState, LoadingState
- ConfirmDialog, Toast

## Constraints

- No marketing, social, multi-user, payment, or gamification elements.
- No workout templates, recipe/barcode, Apple Health, or Telegram UI.
- No AI API configuration — AI analysis is manual copy-paste to ChatGPT.
- Desktop-first responsive design.
- Single user — no avatar, registration, or login beyond optional PIN.

## Output expectations

Generate a cohesive UI design system with key screens and component library. Light and dark variants preferred. All screens should share consistent spacing, typography, card styles, and interaction patterns.
```

---

## 21. Source Documents Status

| Document | Status |
|----------|--------|
| `prd.md` | NOT FOUND — root path missing |
| `docs/prd-wave-details/waves/wave-01.md` | READ |
| `docs/prd-wave-details/waves/wave-02.md` | READ |
| `docs/prd-waves/wave-map.md` | READ |
| `docs/prd-waves/frontend-pages/*.md` | ALL READ (11 pages) |
| `docs/product-verified/functional-spec.md` | READ |
| `docs/product-verified/domain-model.md` | READ |
| `docs/product-verified/acceptance-criteria.md` | READ |
| `docs/technical-verified/api-contracts.md` | READ |
| `docs/technical-verified/auth-security-compliance.md` | READ |
| `docs/superpowers/specs/2026-06-19-wave-01-foundation-design.md` | READ |
| `docs/superpowers/specs/2026-06-19-wave-02-exercise-library-design.md` | READ |

**Missing**: `prd.md` (root). Not critical — all product context was reconstructed from verified docs and wave details.

---

## 22. Final Checklist

- [x] All MVP screens described (11 sections + PIN auth)
- [x] WAVE-01 screens detailed (PIN unlock, PIN settings, session state, dashboard)
- [x] WAVE-02 Exercise Library detailed (list, detail, create/edit, archive/restore, media)
- [x] Global layout described (sidebar, top bar, content, responsive)
- [x] Common states described (loading, empty, error, validation, success, unauthorized)
- [x] UI state matrix provided (15 rows)
- [x] Design constraints listed (12 constraints)
- [x] OpenDesign generation prompt included
- [x] Suggested component list provided (30 components)
- [x] No out-of-scope features introduced
- [x] Source documents read and used as source of truth
- [x] Missing source documents documented