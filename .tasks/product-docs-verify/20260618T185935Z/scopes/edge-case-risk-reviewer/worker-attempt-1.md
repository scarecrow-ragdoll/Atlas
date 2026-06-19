# Edge-Case & Risk Review — Worker Report (Attempt 1)

**Run ID:** 20260618T185935Z
**Source:** docs/product/prd.md
**Scope:** edge-case-risk-reviewer
**Worker:** autonomous agent

---

## Classification Legend

- **GAP** — Missing behavior that needs specification.
- **RISK** — Known failure mode without mitigation.
- **BOUNDARY** — Unspecified boundary or limit.
- **Q-EDGE-N** — Open question requiring product clarification.

---

## 1. PIN/Authentication

### GAP-01: PIN brute-force / lockout
No rate limit, max-attempt threshold, or lockout is specified for PIN entry. A trivial brute-force attack against a numeric PIN would succeed without detection.

### GAP-02: Session lifetime and renewal
No explicit session TTL, sliding expiration, or maximum session duration is defined. A user could have stale sessions indefinitely, or sessions could expire mid-data-entry.

### GAP-03: PIN change revokes existing sessions
Changing or disabling the PIN does not specify whether existing active sessions are invalidated. A previously authenticated browser tab could retain access.

### Q-EDGE-01: What happens when PIN is enabled but no PIN hash exists yet?
The model draft allows `pinHash` to be optional even when `pinEnabled = true`. Is this state possible, and how does the system behave?

### Q-EDGE-02: What is the minimum/maximum PIN length? Allowed characters?
Only numeric? Alphanumeric? Length range? This affects brute-force surface.

---

## 2. Workout Diary

### BOUNDARY-01: Unspecified max exercises per workout day
No limit on how many distinct exercises can be added to one day. With thousands of possible exercises, a UI or payload could degrade.

### BOUNDARY-02: Unspecified max sets per exercise
No limit on set count per exercise. Infinite sets could cause rendering or export bloat.

### BOUNDARY-03: Future dates
Can the user select a future date for a workout? The spec says "задним числом" (backdating) but does not forbid forward-dating.

### BOUNDARY-04: Zero/negative values
No validation rules for weight (0, negative), reps (0, negative), or set numbers. Zero-weight sets with positive reps or zero-rep sets are plausible.

### GAP-04: Duplicate exercises within a workout day
Can the same exercise be added twice to the same day with different orders? Not addressed. This could confuse progression tracking and AI export.

### GAP-05: Delete workout day with referenced data
If a workout day is deleted, are the associated exercises, sets, cardio entries, and comments cascade-deleted or orphaned? Not specified.

### GAP-06: Save failure — partial state
If saving a workout day with 20 sets fails after 15 sets are written, what happens? Optimistic save? Transactional rollback? Not specified.

### Q-EDGE-03: Can a workout day exist with zero exercises and zero cardio?
Is an "empty day" a valid entity, or must every workout day contain at least one exercise or cardio entry?

---

## 3. Exercise Library

### GAP-07: Duplicate exercise name
Can the user create two exercises with the same name? This creates ambiguity in workouts and graphs.

### GAP-08: Exercise deletion with active references
What happens when an exercise used in existing workout days is deleted? Data integrity concern — orphaned WorkoutExercise records referencing a deleted exercise.

### BOUNDARY-05: Exercise name length
No maximum length for exercise name, description, or notes. Could overflow DB column or UI.

### BOUNDARY-06: Working weight = 0 or negative
No validation on `workingWeight`. A zero working weight could break auto-suggest logic and progression signals.

---

## 4. Media (Exercise & Progress Photos)

### BOUNDARY-07: File size limits
No max file size for images/videos. A user could upload multi-GB videos, exhausting disk or causing request timeouts.

### BOUNDARY-08: Unsupported file types
No allowed MIME type whitelist. Malicious or corrupt files could be uploaded.

### GAP-09: Orphaned media on exercise/check-in deletion
Media files stored on disk are not explicitly cleaned up when an exercise or body check-in is deleted. Over years, orphaned files accumulate.

### GAP-10: Media count limit
No limit on how many images/videos can be attached to one exercise. Could cause export bloat and storage exhaustion.

### Q-EDGE-04: What is the maximum total storage for media? Is there a per-instance quota?
Self-hosted could run out of disk. No guidance on storage management.

---

## 5. Cardio

### BOUNDARY-09: Duration = 0 or negative
No validation on durationMinutes.

### BOUNDARY-10: Heart rate bounds
No plausible min/max for avgPulse. Value of 0 or 300 would be stored without validation.

### BOUNDARY-11: Unknown zone with known pulse
If heartRateZone is "unknown" but avgPulse is provided, is the zone auto-derived? Not specified.

### GAP-11: Cardio without workout day linkage
CardioEntry has an optional `workoutDayId`. Orphaned cardio entries (no date match, no workout day) — are they visible in reports? Not clear.

---

## 6. Body Measurements & Check-Ins

### BOUNDARY-12: Weight = 0 or negative
No boundary validation on body weight entries.

### BOUNDARY-13: Body fat percentage bounds
No range (0-100) validation on bodyFatPercentage. Values like 150% or -5% could be stored.

### BOUNDARY-14: Measurement value = 0 or negative
No validation for measurement values (e.g., waist 0 cm).

### GAP-12: Duplicate check-in on same date
Can the user create two body check-ins on the same date? If so, which one is the "current" for graphs and exports? Not specified.

### GAP-13: Photo count enforcement
Check-in spec says 2-4 photos, but no enforcement is specified. User could upload 0 or 50 photos.

### Q-EDGE-05: Are bilateral measurements required for paired body parts?
The spec says "if the user provides only one, it's treated as common." But for graphs and AI analysis, is a missing side treated as zero or omitted?

---

## 7. Nutrition

### BOUNDARY-15: Zero/negative grams in template or override
No validation on `amountGrams` in templates or override items. Zero-gram entries would be meaningless but stored.

### BOUNDARY-16: Macro values = 0 or negative
No validation on caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g. All-zero macros produce no meaningful KJBZU.

### GAP-14: Duplicate product name
Two products with the same name. Unclear which one is used in templates.

### GAP-15: Overlapping weekly templates
Can two templates exist for overlapping week ranges? Which one applies? Not specified.

### GAP-16: Override operation "replace" — full or partial replacement?
`replace` semantics are ambiguous: does it replace the product entirely for that day, or replace only the specified meal label?

### GAP-17: Unbounded nutrition template size
No limit on the number of items per template. Hundreds of products could cause calculation/export issues.

### BOUNDARY-17: Product name length
No max length. Could overflow UI or DB constraints.

### Q-EDGE-06: What happens when a product used in a template is deleted?
Orphan references in template items cascade or block deletion?

---

## 8. Graphs

### GAP-18: No data state
When no data exists for the selected period and filter, is the graph empty, does it show a message, or does it error? Not specified.

### GAP-19: Single data point
A graph with exactly one data point — does it render a line, a dot, or nothing? Not specified.

### BOUNDARY-18: Period boundaries
What happens when start date > end date? What happens for a period of length zero?

---

## 9. AI Export

### RISK-01: Large export generation blocks request
A multi-year export with hundreds of photos could take minutes to generate a ZIP. Synchronous generation could timeout the HTTP request.

### RISK-02: Export contains corrupted data
If data.json is malformed (e.g., null reference, broken encoding), the AI prompt/export still generates a ZIP without validation.

### GAP-20: Disk space failure during ZIP generation
Insufficient disk space during export — partial ZIP file created, user downloads incomplete archive.

### GAP-21: Export photos referenced but deleted
User selects "include photos" but referenced photos were deleted. Missing files in ZIP? No behavior specified.

### GAP-22: Concurrent export requests
User triggers two AI exports simultaneously. Race condition on temp file creation? Two ZIPs overwriting each other?

### BOUNDARY-19: Export file name collision
If two exports happen on the same day, the filename `atlas-ai-export-YYYY-MM-DD.zip` collides.

### Q-EDGE-07: What is the maximum period length for an AI export?
Arbitrary — years of data? Is there a hard limit to prevent resource exhaustion?

### Q-EDGE-08: Are exported photos stripped of EXIF metadata?
Photos may contain GPS location and device info. Privacy-sensitive unless stripped.

---

## 10. Import / Full Restore

### RISK-03: Partial import failure
If import fails after inserting 50% of records, what is the rollback strategy? The spec says "запрет silent partial import" but doesn't specify atomicity (transaction per entity, per batch, or full transaction).

### RISK-04: Import from newer schema version
If the backup was made by a newer app version, the current version may fail to parse or silently misinterpret fields.

### GAP-23: Duplicate detection on import
If the same backup is imported twice into a clean instance (or if data already exists), are rows duplicated? No unique constraint strategy specified (e.g., UPSERT by ID, skip existing, fail on conflict).

### GAP-24: Dry-run success but import fails
Dry-run passes validation but actual insert fails (e.g., constraint violation hidden by dry-run, or disk full during media restore). User-facing error handling not specified.

### GAP-25: Corrupted ZIP upload
Uploaded ZIP is password-protected, truncated, or otherwise damaged. Error message specificity not specified.

### BOUNDARY-20: Max upload file size
No limit on import ZIP size. Multi-GB backup could exhaust memory or timeout.

### BOUNDARY-21: Schema version compatibility range
Is only exact version match allowed? Semver range (^1.0.0)? Minor version tolerance? Not specified.

### Q-EDGE-09: What is the rollback strategy for partial import failures?
Full transaction abort? Best-effort with warning? User sees inconsistent state?

---

## 11. Data Integrity & Concurrency

### RISK-05: Two browser tabs — lost update
Single-user, but browser-level concurrency is real. If user edits a workout in two tabs, the last save wins, potentially losing changes from the first tab. No optimistic locking or conflict detection mentioned.

### RISK-06: Accidental data deletion
No soft-delete, undo, or trash/recycle bin mentioned. An accidental delete of an exercise, workout, or check-in is irreversible.

### GAP-26: Data consistency across linked entities
No cascading rules specified for: deleting a product removes template items referencing it; deleting a check-in removes related measurements and photos; deleting a workout day removes associated sets and cardio linkage.

---

## 12. Privacy & Auditability

### GAP-27: No audit log
No logging of sensitive operations: PIN change, PIN failed attempts, backup export, backup import, AI export generation. Without an audit trail, unauthorized access (e.g., someone who guesses the PIN) is untraceable.

### GAP-28: Photo access without PIN
Spec says "не отдавать media-файлы без авторизации/PIN-сессии". But no mechanism described for how this enforcement works at the API/storage level (e.g., signed URLs, middleware check).

### GAP-29: Export log — URL or file path exposure
If the export ZIP path is logged (e.g., in application logs), ZIP files could be accessible if the log leaks or the path is guessable.

### Q-EDGE-10: Are AI export contents logged anywhere?
Spec says "не логировать содержимое AI export", but is the file path logged? Is the file kept on disk after download?

---

## 13. Data Retention & Cleanup

### GAP-30: Old export file cleanup
AI export ZIPs and backup ZIPs accumulate on the server disk. No retention policy or auto-cleanup is specified. Over years, disk fills without user awareness.

### GAP-31: No data purge capability
Can the user delete all data and start fresh? Factory-reset/wipe feature not mentioned.

### Q-EDGE-11: Are temporary files (export ZIPs) deleted after download?
If ZIPs remain on disk, they accumulate and potentially leak data if the server is compromised.

---

## 14. Migration & Versioning

### GAP-32: Schema migration for live data
The schema version applies only to backup format. What about live database schema migration when the app is updated? Backup/restore alone is not a migration strategy.

### GAP-33: Forward compatibility of old exports
Will exports from version 1.0.0 be importable into version 2.0.0? What if the entity model changes?

### Q-EDGE-12: Is there a migration path for users who update the app without a full backup/restore cycle?
PG schema migrations via migrations framework not mentioned.

---

## 15. External Dependencies & Retries

### GAP-34: Database connection loss mid-operation
No retry or circuit-breaker policy for transient DB failures during data save, export, or import.

### GAP-35: File system errors
Disk full, permission denied, or inode exhaustion during media upload, export generation, or import media restore. Error handling not specified.

### GAP-36: Redis dependency failure
Redis is listed but no fallback behavior if Redis is unavailable (session storage fails => user cannot authenticate even with valid PIN?).

### RISK-07: No self-healing / health checks
No mention of startup probes, liveness checks, or automated recovery for Docker-based services.

---

## Summary

| Category | GAPs | RISKs | BOUNDARYs | Q-EDGEs |
|---|---|---|---|---|
| PIN/Authentication | 3 | 0 | 0 | 2 |
| Workout Diary | 3 | 0 | 4 | 1 |
| Exercise Library | 2 | 0 | 2 | 0 |
| Media | 2 | 0 | 2 | 1 |
| Cardio | 1 | 0 | 3 | 0 |
| Body Measurements | 2 | 0 | 3 | 1 |
| Nutrition | 4 | 0 | 2 | 1 |
| Graphs | 2 | 0 | 1 | 0 |
| AI Export | 3 | 2 | 1 | 2 |
| Import/Restore | 3 | 2 | 2 | 1 |
| Data Integrity | 2 | 2 | 0 | 0 |
| Privacy & Audit | 3 | 0 | 0 | 1 |
| Data Retention | 2 | 0 | 0 | 1 |
| Migration | 2 | 0 | 0 | 1 |
| External Dependencies | 3 | 1 | 0 | 0 |
| **Total** | **37** | **7** | **20** | **12** |

**Verdict for worker:** The PRD defines core business flows clearly but omits most boundary validation, error recovery paths, concurrency handling, and operational failure modes. The 37 GAPs and 12 open questions represent significant risk for a production self-hosted application that stores sensitive data.
