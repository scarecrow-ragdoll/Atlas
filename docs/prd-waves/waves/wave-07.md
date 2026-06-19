# Wave 07: AI Export and Prompt Builder

## Status

user-approved

## User Approval

user-approved (2026-06-18)

## Purpose

Generate AI-ready exports with structured data and prompts for ChatGPT analysis.

## Outcome After Wave

- OUT-W07-001 Prompt builder with period selection
- OUT-W07-002 ZIP export with manifest.json, data.json, summary.md, CSV
- OUT-W07-003 Week flags support
- OUT-W07-004 One-time comment support
- OUT-W07-005 Section toggles (photos optional)

## Included Scope

- CAP-W07-001 Persistent AI context
- CAP-W07-002 User goal storage
- CAP-W07-003 Week flags CRUD
- CAP-W07-004 Prompt generation
- CAP-W07-005 AI export ZIP creation
- CAP-W07-006 manifest.json with export metadata
- CAP-W07-007 data.json with all entities for period
- CAP-W07-008 summary.md with human-readable overview
- CAP-W07-009 CSV files for compatibility

## Excluded Scope

- Direct ChatGPT API call
- OpenAI API integration

## Dependencies

WAVE-01 through WAVE-06

## Surface Categories

backend, integrations, operations

## Risk Class

High - Export correctness for AI consumption

## Recommended Next Planning

$detail-prd-wave for WAVE-07

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |

## Traceability

- docs/product/prd.md Sections 17, 18
- docs/product-verified/domain-model.md#AiExport