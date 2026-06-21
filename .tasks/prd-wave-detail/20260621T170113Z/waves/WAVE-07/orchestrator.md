# WAVE-07 Orchestrator

## Run ID
20260621T170113Z

## Wave ID
WAVE-07: AI Export and Prompt Builder

## Phase
Revision consolidation — addressing 7 reviewer needs-revision verdicts

## All Reviewer Findings Consolidated

### Critical (must fix before ready-for-dev)
| ID | Finding | Source | Resolution |
|---|---|---|---|
| RF-001 | include_photos DEFAULT true violates RULE-025 (must be false) | architecture-codebase, traceability-consistency | Change to DEFAULT false |
| RF-002 | UserProfile vs Settings: 4 planners create full UserProfile, sequencing-fit says read from Settings | cross-planner contradiction | DQ-W07-001: Create UserProfile table (separate from Settings). UserProfile stores goal, height, birthDate, trainingExperience, trainingSplit, progressionStyle, nutritionStrategy, persistentAiContext. Settings stores defaultAiExportWeeks, PIN, units. Separate entities per domain-model.md. |
| RF-003 | Migration numbering: architecture-codebase uses 00093/00094, should be 00091/00092 | architecture-codebase, data-integration-ops | Correct to 00091 (user_profiles), 00092 (ai_exports) |

### High
| ID | Finding | Source | Resolution |
|---|---|---|---|
| RF-004 | gqlgen config gap — missing type bindings for UserProfile, AiExport types | architecture-codebase | Add to atlas-gqlgen.yml |
| RF-005 | display_name inconsistency between planners | architecture-codebase vs data-integration-ops | UserProfile does NOT include display_name — use Settings/atlas_users for display name |
| RF-006 | REST URL patterns conflict across planners | data-integration-ops, architecture-codebase | Use: POST /api/ai-export/generate, GET /api/ai-export/download?exportId=, GET /api/user-profile |
| RF-007 | AC-W07-XXX namespace collision across 3 planners | traceability-consistency | Consolidate all ACs with unique WAVE-07 IDs |
| RF-008 | Missing AC: generatedPrompt returned in generate response body | product-scope-and-ac | Add AC-W07-XXX for prompt response |
| RF-009 | Missing date range input validation tests | testing-exit-criteria | Add TEST-W07-XXX for invalid date ranges |
| RF-010 | Missing ownership-mismatch test for download | testing-exit-criteria | Add TEST-W07-XXX |
| RF-011 | Storage path missing userId scope | security-compliance, data-integration-ops | Use: {base}/{userId}/{exportId}.zip |
| RF-012 | Lifecycle policy conflict: keep-until-replaced vs 7-day TTL | security-compliance, data-integration-ops | DQ-W07-002: 7-day TTL + delete-on-regeneration |
| RF-013 | Temp-file-atomic-rename missing from data-integration-ops | security-compliance | Add to ZIP generation slice |
| RF-014 | Week flags REST vs GraphQL resolution undocumented | sequencing-other-wave-fit | Note: PAGE-009 uses WAVE-04 GraphQL directly. No REST endpoint needed. |
| RF-015 | Max export size limit not specified | security-compliance | Add 100MB hard limit |
| RF-016 | CAP-W07-003 (week flags CRUD) scope collision with WAVE-04 | sequencing-fit | Remove from WAVE-07 scope. WAVE-07 reads week flags via WAVE-04 service. |

### Questions Raised
| ID | Question | Severity | Status |
|---|---|---|---|
| DQ-W07-001 | UserProfile vs Settings: new table or extend Settings? | needs-owner-decision | Resolved: Create separate UserProfile table per domain-model.md. Settings exists for app config; UserProfile exists for user-specific data. |
| DQ-W07-002 | Export ZIP lifecycle: 7-day TTL or keep-until-replaced? | needs-owner-decision | Resolved: 7-day TTL + delete-on-regeneration |
| DQ-W07-003 | Sync vs async export generation? | needs-owner-decision | Sync for MVP (small data volumes). Async ready if needed. |
| DQ-W07-004 | Photo storage in export: base64 in data.json or files in photos/? | needs-owner-decision | Files in photos/ subfolder. base64 in data.json for thumbnail preview. |
| DQ-W07-005 | Schema version for manifest.json | needs-owner-decision | Use integer schemaVersion = 1 for first version |
| DQ-W07-006 | App version in manifest? | needs-owner-decision | Include appVersion from build metadata if available, omit otherwise |

## Next Action
Consolidate all planner and reviewer findings into the candidate wave package in staging. Then dispatch final fit review.