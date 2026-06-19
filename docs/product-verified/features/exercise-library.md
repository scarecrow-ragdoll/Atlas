# Exercise Library

## Source Evidence

PRD §11, §26.1.

## User Problem

Maintain a personal catalog of exercises with working weights, notes, and instructional media.

## Scope

In MVP. User-created exercises only (no starter catalog). Media upload for images and video.

## Behavior

- User creates exercises with name, muscle groups, description, personal notes, working weight
- User uploads images and video per exercise
- Media can be added or removed after exercise creation
- Exercise can be marked active or inactive
- Working weight stored on exercise; snapshot captured in workout day

## Acceptance Criteria

AC-002, AC-003, AC-004, AC-043 through AC-047.

## Derived Requirements

None beyond source evidence.

## Edge Cases

EDGE-002: Exercise name duplicate.
EDGE-005: Exercise media file deleted from filesystem but DB record remains.

## Dependencies

Media storage (filesystem volume).

## Open Questions

Q-FEAT-004: Muscle groups representation (free text, enum, multi-select).
Q-FEAT-005: Media file size limits and accepted formats.