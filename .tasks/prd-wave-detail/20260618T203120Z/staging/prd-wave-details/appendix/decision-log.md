# Decision Log

## Source Wave Gate
source-wave-gate: passed
Selected wave: WAVE-01
Source path: docs/prd-waves/waves/wave-01.md
Date: 2026-06-18

## User Wave Approvals
- User approved WAVE-01 source wave map on 2026-06-18 (from $decompose-prd-waves run)
- This detailed wave brief awaits user approval

## Scope Decisions
- PIN rate limiting not implemented in WAVE-01 (deferred, DQ-W01-001)
- PIN session TTL default 7 days, configurable via environment
- Settings as single-row key-value table (not per-setting table)
- Hybrid API model confirmed per TDEC-001
- Local filesystem for media storage (pluggable to S3 later)
- PIN auth and admin auth as separate middleware chains sharing /graphql endpoint

## Codebase Fit Decisions
- Fitness-domain packages follow existing admin package patterns (service/repository/middleware/handler)
- Fitness GraphQL schema files alongside admin schema files in libs/graphql/schema/
- PIN session store uses same Redis with separate key prefix (fitness:session:)
- gqlgen and sqlc config extended with fitness-domain paths; no existing config changed
- PIN guard middleware added alongside existing admin_auth middleware

## Deferrals
- PIN rate limiting (DQ-W01-001) — deferred to owner decision
- Barcode scanner — deferred (from wave map)
- Apple Health integration — deferred (from wave map)
- Telegram bot — deferred (from wave map)
- OpenAI API integration — deferred (from wave map)

## Rejected Assumptions
- "Admin auth can be reused for fitness" — rejected. Admin auth (cookie-based, for M-WEB-ADMIN operators) and PIN auth (token-based, for fitness app user) are separate concerns. Both coexist on /graphql.