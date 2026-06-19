# WAVE-02 security-privacy-compliance Review Attempt 1

## Verdict
needs-revision

## Sources Read
- planner-security-compliance-attempt-1.md
- planner-data-integration-ops-attempt-1.md
- planner-architecture-codebase-attempt-1.md
- planner-product-ac-attempt-1.md
- planner-testing-exit-attempt-1.md
- docs/technical-verified/auth-security-compliance.md
- docs/product-verified/business-rules.md (RULE-022, RULE-023, RULE-024)
- docs/product-verified/edge-cases.md (EDGE-011, EDGE-012, EDGE-014, EDGE-020)
- docs/prd-wave-details/waves/wave-01.md

## Coverage Check
Security coverage addresses: PIN auth protection, file upload validation (MIME types, size limits, path traversal), audit logging, privacy (no personalNotes in logs), and abuse vectors. Authorization is single-user — no multi-tenancy concerns.

## Evidence Check
All security claims trace to source docs. PIN protection → RULE-022/023, file limits → TDEC-008, audit → TDEC-004, privacy → TDEC-004.

## Codebase Fit Check
WAVE-02 security relies entirely on WAVE-01 PIN auth middleware. No new auth mechanism. File validation is specific to WAVE-02. Audit log markers follow WAVE-01 pattern ([Domain][action][BLOCK_NAME]).

### Issues Found

1. **Server-side MIME detection**: The planner correctly raises DQ-W02-005 about Content-Type spoofing. The planner should make a recommendation, not just raise a question. Recommendation: use `http.DetectContentType()` (reads first 512 bytes) as primary MIME detection, reject if it doesn't match the declared Content-Type. This is a security best practice.

2. **File upload memory safety**: The planner mentions standard Go HTTP read timeout but doesn't address memory-bounded upload handling. For files up to 250MB, `multipart.Form` parses into memory by default. Should use `r.ParseMultipartForm(maxBytes)` with appropriate limit to prevent memory exhaustion. Add this to the design.

3. **Media access when PIN is disabled**: RULE-022 says "no auth when PIN off." RULE-024 says "media files not accessible without valid session." This contradiction (TQ-AUTH-006, resolved by TDEC-037) needs explicit resolution in WAVE-02. The planner doesn't mention how this contradiction is handled.

4. **Audit log markers inconsistent across planners**: Product-ac planner mentions audit; security-compliance planner mentions [Exercise][*] markers; data-integration-ops planner doesn't mention audit. All domains should consistently reference the same audit pattern.

5. **ExerciseMedia file access path traversal**: The planner addresses upload path traversal but not download path traversal. The GET handler must also sanitize the {id} parameter to prevent path traversal when resolving the file path. This is especially important if IDs are user-facing (they are UUIDs, so low risk, but worth noting).

6. **No mention of CORS**: If exercise media REST endpoints are served from a different path than existing admin endpoints, CORS configuration must be verified. WAVE-02 should document which CORS config applies to the exercise media endpoints (publicCORS from main.go or a new config).

## Other-Wave Fit Check
Security model is inherited from WAVE-01. No new attack surface beyond file upload handling. WAVE-03 inherits WAVE-02 security posture.

## Acceptance Criteria Check
AC-W02-014 through AC-W02-020 cover security adequately. AC-W02-019 (log sanitize) is better as EC (moved by product-ac reviewer).

## Exit Criteria Check
EC-W02-011 through EC-W02-014 cover security. EC-W02-012 (file rejection) needs explicit test coverage for both MIME type and size.

## Verification Check
TEST-W02-014 through TEST-W02-019 cover auth, file validation, path traversal, and log privacy. Good coverage.

## Question Ledger Check
DQ-W02-005 (MIME detection) critical — recommend resolving with server-side detection. DQ-W02-006 (signed URLs) is properly deferred.

## Unsupported Or Invented Claims
No unsupported claims. Security analysis is grounded in source documents.

## Required Revisions
1. **Recommend server-side MIME detection**: Resolve DQ-W02-005 with explicit recommendation (use http.DetectContentType + cross-check).
2. **Add memory-safe upload handling**: Document ParseMultipartForm with maxBytes limit.
3. **Resolve PIN disabled + media access contradiction**: Explicitly document how TDEC-037 resolution maps to WAVE-02 implementation (media still requires session even if PIN is disabled? or open when PIN is disabled?).
4. **Add download path traversal check**: Document that GET endpoint uses UUID parameter which is inherently safe, but note the principle.
5. **Add CORS documentation**: Specify which CORS config applies to exercise media REST endpoints.

## Approval Notes
Solid security baseline. 5 revision items are clarifications and hardening — no fundamental security flaws. After revisions, will approve.