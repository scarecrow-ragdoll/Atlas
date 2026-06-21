## Review: Security / Privacy / Compliance

**Run**: 20260621T170113Z
**Wave**: WAVE-07 (AI Export + Prompt Builder)
**Reviewer Role**: security-privacy-compliance
**Attempt**: 1

---

### Verdict: needs-revision

---

### Summary

The planners correctly identify the major security concerns (PIN auth coverage, photo opt-in enforcement, log privacy, ownership validation) and propose appropriate ACs. However, three material gaps exist between the two planners' designs that must be reconciled before implementation begins.

---

### Gaps Requiring Revision

#### GAP-1: Temp-file-atomic-rename pattern missing from data-integration-ops planner

The security-compliance planner requires AC-W07-SEC-004 (write ZIP to temp file first, then atomic rename on success; clean up temp on failure). The data-integration-ops planner's ZIP write step (section 1.1, step 7) writes directly to `{exportBasePath}/{exportID}/export.zip` with no temp-file stage. This is the primary mitigation for EDGE-024 (disk full during generation).

The codebase has no existing atomic rename pattern (grep for `os.Rename` returned zero hits in Go source), so this would be new code that both planners need to align on.

**Resolution required**: Data-integration-ops planner must incorporate the temp-file-then-atomic-rename pattern into the ZIP write step and ZIP format section. The `AiExport` file path must only be set after the rename succeeds.

#### GAP-2: Export file lifecycle policy conflict

- **Security-compliance planner** (section 4, Option C): Keep ZIP until next generation replaces it. Explicitly defers cleanup automation to post-MVP.
- **Data-integration-ops planner** (section 2.5): 7-day TTL with scheduled cleanup task, hard-delete of both file and record.

These are incompatible. If both are implemented, the cleanup task would delete files the security model expects to persist for re-download, and the re-generation replacement logic would conflict with the TTL.

**Resolution required**: Choose one model. The security recommendation (keep until replaced) is simpler and safer for MVP, since it avoids cleanup-task risks (file-delete failure leading to dangling records). If the 7-day TTL is preferred, it must be reflected in the security ACs and the re-generation replacement logic must be removed.

#### GAP-3: Export ZIP storage path pattern — userId scope

- **Security-compliance planner** (section 1): `{ExportBasePath}/{userId}/{export-uuid}.zip` — user-scoped paths for future multi-user isolation.
- **Data-integration-ops planner** (section 2.5): `{export_base_path}/{export_uuid}/export.zip` — no userId in path.

The DB query (`GetAiExportByID`) correctly filters by `user_id`, so filesystem-level isolation is defense-in-depth, not a primary control. However, the two planners should agree on one pattern. The security-compliance planner's approach adds negligible complexity and prevents path-level confusion later.

**Resolution required**: Adopt the user-scoped path pattern `{basePath}/{userId}/{exportUuid}/export.zip` (or explicitly document why MVP omits this scope).

#### GAP-4: Max export size limit — missing from data planner ACs

Security-compliance planner proposes AC-W07-SEC-009 (reject generation if estimated size exceeds 100MB). The data-integration-ops planner has `max_range_days: 365` and `max_photos_in_export: 20` as size controls but no hard size limit or AC for rejection. Q-W07-DIO-03 discusses streaming threshold, not a hard rejection.

**Resolution required**: Add a configurable `max_export_size_bytes` (recommended 100MB) to `AiExportConfig` and a corresponding AC that generation is rejected with a clear error when estimated size exceeds the limit. This is primarily an operational reliability concern with a minor security angle (disk exhaustion).

---

### Items Confirmed Correct

| Area | Status | Evidence |
|------|--------|----------|
| PIN auth on POST `/api/ai-export` | ✅ Both planners agree | security pl: AC-W07-SEC-001, data pl: section 1.1 auth |
| PIN auth on GET `/api/ai-export/download` | ✅ Both planners agree | security pl: AC-W07-SEC-002, data pl: section 1.2 auth |
| Ownership check on download (404 if mismatch) | ✅ Both planners agree | security pl: AC-W07-SEC-006, data pl: AC-W07-DIO-16, sqlc query `AND user_id = $2` |
| `includePhotos` defaults to false at service layer | ✅ Both planners enforce | security pl: AC-W07-SEC-003, data pl: DDL `DEFAULT false` + domain model invariant |
| Log privacy — no sensitive values | ✅ Both planners align | security pl: AC-W07-SEC-008, data pl: section 5 log markers (no user data) |
| RULE-026 — manual generation only | ✅ No auto-export path | Both planners: POST requested by user only |
| RULE-027 — manual copy-paste, no auto-transmit | ✅ Prompt returned in API response body, not transmitted | security pl: section 7, data pl: section 1.1 step 8 |
| PIN-disabled accessibility consistent with WAVE-04 TDEC-037 | ✅ Both planners note this | security pl: section 1 |
| RULE-025 — photos excluded by default | ✅ Enforced at DB + service layer | DDL default false + security AC-W07-SEC-003 |
| UUID-based storage to prevent enumeration | ✅ Both planners use UUID paths | security pl: `uuid`, data pl: `export_uuid` |
| Fresh snapshot per generation (no stale data) | ✅ security pl: EDGE-008 | AC-W07-SEC-010 |
| Bootstrap UserProfile for default user | ✅ data pl: section 2.4 | Prevents 404 on first profile read |
| Cleanup safety: file deleted before record | ✅ data pl: section 7 | Correct ordering prevents dangling files |

---

### New ACs Confirmed (from security planner, compatible with both)

All 10 ACs proposed by the security-compliance planner (AC-W07-SEC-001 through AC-W07-SEC-010) are sound and should be adopted after resolving GAP-1 through GAP-4 above. No additional security ACs are needed.

---

### Open Questions from Security Perspective

| ID | Question | Impact | Recommendation |
|----|----------|--------|----------------|
| Q-W07-SEC-001 (carried) | Should export ZIPs be auto-cleaned after download or on TTL? | Disk space, data control | Resolve GAP-2 above first |
| Q-W07-SEC-004 (carried) | Max export size? | Anti-abuse, reliability | Resolve GAP-4 above; add to config |
| Q-W07-REV-001 | If temp-file-atomic-rename is new code with no prior pattern, should a shared utility function be extracted? | Code quality, testability | Yes — create `tools.WriteFileAtomic(path, data)` or similar utility so the pattern is testable and reusable |

---

### Conclusion

Approved with revision requirements (GAP-1, GAP-2, GAP-3, GAP-4). No blocking security vulnerabilities identified for single-user MVP scope, but the design inconsistencies between the two planners will cause integration issues and security gaps during implementation if not resolved upfront.