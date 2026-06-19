# Progress Photos

## Source Evidence

PRD §14.

## User Problem

Visual progress tracking over time through body photos tied to weekly check-ins.

## Scope

In MVP. Photos attached to check-ins, not standalone. Media not publicly accessible without PIN session.

## Behavior

- Photos attached to weekly check-in
- Fields: date, check-in link, file, optional label, optional angle (front/side/back/custom), optional comment
- Not publicly accessible without valid session
- Included in full backup export
- Included in AI export only when explicitly opted in

## Acceptance Criteria

AC-056, AC-111, AC-119.

## Derived Requirements

None beyond source evidence.

## Edge Cases

EDGE-014: Media URL accessed directly without valid session.
EDGE-020: Media file deleted from filesystem but DB record remains.

## Open Questions

Q-FEAT-005: Media file size limits and accepted formats.

## Dependencies

Weekly check-in feature, media storage.