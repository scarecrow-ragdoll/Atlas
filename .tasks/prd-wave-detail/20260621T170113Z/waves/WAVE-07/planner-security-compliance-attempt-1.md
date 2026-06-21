# WAVE-07 Planner: Security / Compliance

## Auth Architecture

### PIN Auth Coverage

WAVE-07 introduces two REST endpoints: `POST /api/ai-export` (generate) and `GET /api/ai-export/download`. Both must follow the WAVE-01 PIN auth middleware pattern (`GetAtlasUserID`) established in WAVE-02:

- `POST /api/ai-export` → 401 without valid PIN session
- `GET /api/ai-export/download` → 401 without valid PIN session
- When PIN is disabled, both endpoints accessible without auth (consistent with WAVE-04 TDEC-037)

No GraphQL operations in WAVE-07. The prompt builder context (persistent AI context, goal) is managed through UserProfile (WAVE-01/02 patterns).

### Auth Test Cases Per Endpoint

- `POST /api/ai-export` without valid session → AuthError/401
- `GET /api/ai-export/download` without valid session → AuthError/401

## Privacy Considerations

### Sensitive Data in Export

The ZIP export captures ALL user data for the selected period. Sensitivity classification:

| Data | Sensitivity Level | Included? | Source |
|------|-------------------|-----------|--------|
| Workout exercises, sets, weights, RPE/RIR | Low | By default | §17.3 |
| Exercise comments | Medium | By default | §17.3, AC-042 |
| Cardio entries | Low | By default | §17.3 |
| Body weight entries | Medium | By default | §17.3 |
| Body measurements | Medium | By default | §17.3 |
| Body fat % | Medium | By default | §17.3 |
| Progress photos | **High** | **Opt-in only** | RULE-025, AC-077, AC-112 |
| Nutrition data | Low | By default | §17.3 |
| User goal/context | Medium | By default | §17.3 |
| Week flags (sleep, stress, illness, etc.) | Medium | By default | §18.4 |

### Photo Opt-In Rule

RULE-025 and AC-077/AC-112: Photos excluded by default. User must explicitly toggle photos on. The service layer must enforce `includePhotos` defaulting to `false` (domain model invariant: `AiExport.includePhotos defaults to false`).

### AC-118: No AI Export Content Logging

Consistent with WAVE-06 pattern: log markers must NOT include:
- Generated prompt content
- User comments (persistent or one-time)
- Body weight/measurement values in export metadata logs
- Photo file paths from export
- Week flag notes

Log markers record only: export generation success/failure, export ID, period start/end dates (no values), section toggles selected (boolean flags only), file size of generated ZIP.

## ZIP Export Security

### 1. Export ZIP Storage Path

ZIP files must be stored with user-scoped paths to support future multi-user isolation even in MVP:

```
<ExportBasePath>/<userId>/<export-uuid>.zip
```

- `ExportBasePath` configured via settings or env (e.g., `./data/ai-exports`)
- `userId` derived from `GetAtlasUserID(ctx)` — single default user in MVP
- `export-uuid` is a generated UUID, not sequential, to prevent enumeration
- `AiExport.exportFilePath` stores the full path for retrieval

**Rationale:** User-scoped paths prevent cross-user data leakage in future multi-user mode. UUID-based filenames prevent enumeration of export files.

### 2. Download Endpoint Ownership Validation

`GET /api/ai-export/download?id=<export-id>` must:
1. Require valid PIN session (`GetAtlasUserID`)
2. Look up `AiExport` record by ID
3. Verify `AiExport.userId` matches session user ID (no-op in single-user MVP, but enforced for future-proofing)
4. Return 404 Not Found if export does not exist or does not belong to user
5. Serve file with `Content-Disposition: attachment; filename="atlas-ai-export-YYYY-MM-DD.zip"` and correct `Content-Type: application/zip`
6. Stream file from disk — do not buffer entire ZIP in memory

### 3. Export Generation Request Validation

`POST /api/ai-export` must:
1. Require valid PIN session
2. Validate date range: startDate <= endDate, range not in far future (max +1 day from today)
3. Enforce `includePhotos` default `false` at service layer (defense-in-depth — must not be `true` unless explicitly sent)
4. Limit concurrent generation (single export at a time per user) to prevent disk/CPU abuse
5. Return `exportId` immediately asynchronously or synchronously — recommend synchronous for MVP simplicity given single-user scope

### 4. ZIP File Cleanup

**Question:** When should export ZIP files be deleted?

Options:
- **A:** Delete after successful download (one-shot download model)
- **B:** Keep until user generates a new export (implicit cleanup)
- **C:** Keep indefinitely (user can re-download)
- **D:** Delete on a TTL (e.g., 24 hours after generation)

**Recommendation for MVP:** Option **C** (keep until next generation). Single-user, self-hosted. User may want to re-download. Oldest export per user is replaced when new one is generated. Add explicit deferral for cleanup automation.

**Edge to handle:** If download never happens (e.g., user generates and leaves), the ZIP stays on disk until next generation overwrites it. No auto-cleanup in MVP.

### 5. EDGE-024: Disk Full During Export Generation

Service must:
1. Write ZIP to a temp file first (`<ExportBasePath>/<userId>/.tmp-<uuid>.zip`)
2. On successful write, atomically rename to final path
3. On failure (disk full, write error), clean up temp file, return error
4. `AiExport.exportFilePath` only set on final rename success

### 6. EDGE-008: No Data in Period

Generate ZIP with empty section files (empty JSON arrays, empty CSVs with headers, summary.md noting "no data in this period"). Must not include photos/ directory when no photos or photos not opted in. This is a functional concern with security implication: empty exports must not leak stale data from previous exports — each generation is a fresh snapshot.

### 7. RULE-027: Manual Copy-Paste

The prompt is designed for manual copy-paste to external AI (RULE-027). No automatic data transmission. This eliminates external data leakage risk. Generated prompt is returned in API response body and optionally written to `generatedPrompt` field on `AiExport` record. Prompt text is NOT logged per AC-118.

### 8. Filename Collision Prevention

Current filename pattern: `atlas-ai-export-YYYY-MM-DD.zip`. If user generates two exports on same day, the filename must be disambiguated. Options:
- Append UUID suffix: `atlas-ai-export-YYYY-MM-DD-<short-uuid>.zip`
- Use timestamp: `atlas-ai-export-YYYY-MM-DDTHHMMSS.zip`

**Recommendation:** Include short UUID suffix to guarantee uniqueness without exposing generation times.

## Compliance Considerations

### MVP Context
- Single-user deployment on own infrastructure
- No external data transmission (RULE-027 — manual copy-paste)
- No third-party data processors
- Self-hosted, user controls their data

### Recommendations
1. User-scoped export paths for future multi-user readiness
2. Ownership check on download endpoint (defense-in-depth even in MVP)
3. `includePhotos` defaults enforced at service layer, not just frontend
4. Temp-file-then-rename pattern for export generation to handle disk-full gracefully
5. UUID-based export filenames to prevent enumeration
6. No sensitive data in logs per AC-118
7. Export ZIP deleted only on next generation (MVP simplicity)

## Edge Cases

### EDGE-011: PIN Enabled But PIN Hash Corrupted
Consistent with WAVE-02/04: service returns auth error, no export accessible.

### EDGE-014: Export Download URL Accessed Directly Without Session
Protected by PIN middleware — returns 401.

### EDGE-024: Disk Full During Export Generation
Temp-file-then-rename pattern prevents partial ZIP. Error returned if generation fails.

### EDGE-015: Stale Session During Download
Session middleware handles. If session expired, user must re-authenticate via PIN. No data leakage.

## Open Security Questions

| ID | Question | Why It Matters | Source | Recommendation |
|----|----------|----------------|--------|----------------|
| Q-W07-SEC-001 | Should export ZIPs be auto-cleaned after download or on TTL? | Disk space management, user data control | EDGE-024, §17 | Defer to post-MVP. MVP: keep until next generation replaces the file. |
| Q-W07-SEC-002 | Should body measurements and weight data be included by default or opt-in? | Privacy sensitivity | §17.3, RULE-025 | Include by default per PRD §17.3 — photos are the only opt-in category per RULE-025. |
| Q-W07-SEC-003 | Should week flag notes be included in export data? | Week flags may contain sensitive health context (cycle, injury, stress) | §18.4 | Include by default — flags are user-chosen and part of AI analysis context. Do not log. |
| Q-W07-SEC-004 | What is the max export size? | Disk usage, download UX, memory during generation | Q-FEAT-010 | Recommend hard limit of 100MB uncompressed for MVP. Reject generation with error if data exceeds limit. |
| Q-W07-SEC-005 | Should re-download of same export be allowed? | User may lose the ZIP after download | Page-009, GET endpoint | Yes — keep export file until next generation. |

## New AC Contributions

These ACs are recommended additions to the WAVE-07 acceptance criteria based on security analysis:

| Proposed ID | Description | Source |
|-------------|-------------|--------|
| AC-W07-SEC-001 | `POST /api/ai-export` returns 401 without valid PIN session | WAVE-01 auth pattern |
| AC-W07-SEC-002 | `GET /api/ai-export/download` returns 401 without valid PIN session | WAVE-01 auth pattern |
| AC-W07-SEC-003 | `POST /api/ai-export` enforces `includePhotos: false` default at service layer (defense-in-depth) | RULE-025, AC-112 |
| AC-W07-SEC-004 | Export generation writes to temp file first; on success, atomically renames to final path. On failure, temp file cleaned up and error returned | EDGE-024 |
| AC-W07-SEC-005 | Export filename includes unique identifier to prevent same-day collision | §17.4 filename pattern |
| AC-W07-SEC-006 | Download endpoint validates export belongs to session user; returns 404 if mismatch or not found | Ownership principle |
| AC-W07-SEC-007 | Export ZIP stored at user-scoped path: `<ExportBasePath>/<userId>/<uuid>.zip` | Future multi-user readiness |
| AC-W07-SEC-008 | No export content (prompt text, user comments, measurement values, photo paths) logged — log only export ID, success/failure, period dates, section toggles as booleans | AC-118, WAVE-04 log privacy pattern |
| AC-W07-SEC-009 | Export generation rejected if estimated data size exceeds configured maximum (recommend 100MB uncompressed) | Q-FEAT-010, EDGE-024 |
| AC-W07-SEC-010 | Each generation creates a fresh data snapshot — no stale data from previous exports leaked into new ZIP | EDGE-008 |

## Traceability

- `docs/prd-waves/waves/wave-07.md` — source wave boundary, outcomes, CAP list
- `docs/product-verified/features/ai-export.md` — export behavior, section toggles, AC-074–AC-083
- `docs/product-verified/features/ai-prompt-builder.md` — prompt builder behavior, AC-084–AC-089
- `docs/product-verified/business-rules.md` — RULE-021 (4 weeks), RULE-025 (photos opt-in), RULE-026 (on user request), RULE-027 (manual copy-paste)
- `docs/product-verified/acceptance-criteria.md` — AC-074–AC-083 (export), AC-112 (photos opt-in), AC-118 (no content logging), AC-117–AC-120 (log privacy)
- `docs/product-verified/domain-model.md` — AiExport entity with `includePhotos` default false, lifecycle (draft/generated via exportFilePath presence)
- `docs/product-verified/edge-cases.md` — EDGE-008 (no data), EDGE-024 (disk full)
- `docs/product-verified/actors-and-permissions.md` — user permissions for AI export, privacy rules (§24.1)
- `docs/product-verified/functional-spec.md` — AI Export §17-18 behavior
- `docs/prd-waves/frontend-pages/page-009.md` — backend API contracts: POST + GET endpoints
- `docs/prd-wave-details/waves/wave-01.md` — PIN auth middleware pattern
- `docs/prd-wave-details/waves/wave-04.md` — WAVE-04 security/privacy patterns (log sanitization, photo handling, MIME validation, auth coverage)
- `.tasks/technical-docs-verify/20260618T185935Z/scopes/auth-security-compliance/worker-attempt-1.md` — full security gap analysis (TGAP-AUTH-001..012)