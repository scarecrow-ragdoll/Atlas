# WAVE-02 security-compliance Planner Attempt 2

## Cycle 1 Reviewer Feedback Addressed

### 1. Server-Side MIME Detection (DQ-W02-005 Resolution)
Decision: Use `http.DetectContentType()` (reads first 512 bytes of file) as the primary MIME detection mechanism. Cross-check against allowed types:
```go
func validateMediaType(fileBytes []byte) (string, error) {
    detected := http.DetectContentType(fileBytes)
    // Allow only specific MIME types
    allowed := map[string]string{
        "image/jpeg": ".jpg",
        "image/png":  ".png",
        "image/webp": ".webp",
        "video/mp4":  ".mp4",
        "video/quicktime": ".mov",
        "video/webm": ".webm",
    }
    ext, ok := allowed[detected]
    if !ok {
        return "", fmt.Errorf("disallowed file type: %s", detected)
    }
    return ext, nil
}
```
This resolves DQ-W02-005 with a concrete approach.

### 2. Memory-Safe Upload Handling
Use `r.ParseMultipartForm(maxBytes)` with maxBytes from TDEC-008 limits (300MB per single upload request). This prevents unbounded memory allocation from large file uploads. Implementation:
```go
func (h *ExerciseMediaHandler) Upload(w http.ResponseWriter, r *http.Request) {
    const op = "ExerciseMediaHandler.Upload"
    log := logger.FromContext(r.Context()).With(zap.String("op", op))
    log.Debug("[ExerciseMedia][upload][BLOCK_PARSE_FORM] parsing upload")
    
    maxBytes := int64(300 << 20) // 300MB per TDEC-008
    if err := r.ParseMultipartForm(maxBytes); err != nil {
        writeError(w, http.StatusBadRequest, "FILE_TOO_LARGE", "upload exceeds maximum size", nil)
        return
    }
    // ... process file
}
```

### 3. PIN Disabled + Media Access Contradiction Resolution (TDEC-037)
Following TDEC-037 resolution: When PIN is disabled, all endpoints (including exercise media) are accessible without authentication. When PIN is enabled, media requires valid session. Implementation approach:
- WAVE-01 PIN middleware handles the enabled/disabled check
- If PIN is disabled, middleware passes the request through without auth
- This is consistent with RULE-022 ("no auth when PIN off")
- Media endpoints use the same middleware, inheriting this behavior

If WAVE-01 middleware does not handle the PIN-disabled case, WAVE-02 explicitly checks PIN status in its handler. But this is delegated to WAVE-01's middleware contract.

### 4. Download Path Traversal Prevention
GET /api/v1/exercise-media/{id} uses UUID parameter. UUIDs are inherently safe (alphanumeric + hyphens, no path separators). Still, the handler resolves the file path from the DB record (file_path column), not from user input. No path traversal risk exists. Principle documented for audit.

### 5. CORS Configuration
Exercise media REST endpoints use the existing `publicCORS` config from main.go (which allows configured origins with Content-Type header). The PIN auth middleware is applied before CORS — CORS preflight (OPTIONS) does not require auth. This matches the existing pattern where /healthz and /api/users use publicCORS without auth.

### 6. File Content-Type Cross-Check
After reading the file header for MIME detection, also validate:
- Extension from detected MIME matches stored extension
- File size against per-type limits:
  - Images (JPEG/PNG/WEBP): max 25MB
  - Video (MP4/MOV/WEBM): max 250MB
- Single upload total: max 300MB

## Updated Security ACs
- AC-W02-022 (MIME type validation) — now includes server-side magic byte detection
- AC-W02-023 (File size validation) — now per-type limits enforced
- AC-W02-024 (Path traversal prevention) — documented as UUID-safe, but principle maintained

## Log Privacy Verification
Log sanitization test (TEST-W02-019) verifies: personalNotes content does not appear in log output. Implementation: exercise service does not log exercise content fields. Audit markers log exercise_id, action, success/failure only.

## Updated Questions

| ID | Severity | Status | Change |
| --- | --- | --- | --- |
| DQ-W02-005 | needs-owner-decision | **resolved** | Decision: use http.DetectContentType() server-side |
| DQ-W02-006 | deferred | deferred | Decision deferred (signed URLs not needed for MVP)