# PAGE-010: Import/Export

## Status

user-approved

## Page Purpose

Full backup operations: export all data to ZIP, import from ZIP with validation.

## What Is On This Page

- Export button (with media option)
- Import file picker
- Dry-run summary display
- Import confirmation
- Backup history (optional)

## Functional Parts

- Export all data
- Import with validation
- Media inclusion toggle
- Summary preview

## Empty States

- No backup yet - "Create your first backup"

## Loading And Error States

- Export in progress - spinner
- Import validation - spinner
- Invalid archive error - clear message

## Backend Dependencies

- POST /api/backup/export
- POST /api/backup/import (dry-run then confirm)

## Explicit Deferrals

- Cloud backup - future scope

## Open Questions

- None blocking

## Raw PRD Traceability

docs/product/prd.md Sections 20

## Verified PRD Traceability

docs/product-verified/domain-model.md#AiExport, #AiReview