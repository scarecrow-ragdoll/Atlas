# WAVE-07 Testing & Exit Criteria Report

**Run ID:** 20260621T170113Z  
**Wave ID:** WAVE-07 (AI Export and Prompt Builder)  
**Planner role:** testing-exit  
**Attempt:** 1  

---

## 1. Exit Criteria (EC-W07-XXX)

Measurable wave-completion conditions. Each must pass before WAVE-07 merges to develop.

| ID | Criterion | Source | Verification Method |
|---|---|---|---|
| EC-W07-001 | UserProfile model, service, repo, migration, and resolver exist, compile, and pass focused unit tests | Wave-07, AC-085, AC-086 | `cd apps/api && go test ./internal/atlas/service/ -run TestUserProfile -count=1 -v` |
| EC-W07-002 | UserProfile goal, height, birthDate validation rejects empty goal, invalid height (<=0 or >300), and invalid birthDate (future date) | RULE-003 analog | `cd apps/api && go test ./internal/atlas/service/ -run TestUserProfile.*Validation -count=1 -v` |
| EC-W07-003 | UserProfile returns default values (empty goal, zero height, nil birthDate) on first access with no profile row | AC-084, AC-085 | `cd apps/api && go test ./internal/atlas/service/ -run TestUserProfile.*Default -count=1 -v` |
| EC-W07-004 | AiExport model, service, repo, migration, resolver exist and pass focused unit tests | Wave-07 | `cd apps/api && go test ./internal/atlas/service/ -run TestAiExport -count=1 -v` |
| EC-W07-005 | Prompt generation produces correct text with active section toggles, persistent AI context, one-time comment, and week flags | AC-084, AC-085, AC-087, AC-088, AC-089 | `cd apps/api && go test ./internal/atlas/service/ -run TestAiExportPrompt -count=1 -v` |
| EC-W07-006 | ZIP generation produces valid ZIP with manifest.json, data.json, summary.md, CSVs (workouts, measurements, nutrition, cardio) | AC-078, AC-079, AC-081 | `cd apps/api && go test ./internal/atlas/service/ -run TestAiExportZIP -count=1 -v` |
| EC-W07-007 | manifest.json contains export type, schema version, app version, date, period, sections | AC-081 | Covered by EC-W07-006 |
| EC-W07-008 | data.json contains all selected data sections for the requested date range | AC-082 | Covered by EC-W07-006 |
| EC-W07-009 | summary.md includes period, goal, workout stats, exercise trends, weight/measurement changes, nutrition summary, cardio, comments | AC-083 | Covered by EC-W07-006 |
| EC-W07-010 | ZIP uses photos/ directory when includePhotos=true; excludes photos when includePhotos=false (default) | AC-077, AC-080, RULE-025 | `cd apps/api && go test ./internal/atlas/service/ -run TestAiExportZIP.*Photos -count=1 -v` |
| EC-W07-011 | Empty date range produces valid export with empty data sections (no crash) | EDGE-008 | `cd apps/api && go test ./internal/atlas/service/ -run TestAiExport.*EmptyDateRange -count=1 -v` |
| EC-W07-012 | POST /api/ai-export returns 200 with export ID and triggers ZIP generation on disk | AC-024 | `cd apps/api && go test ./internal/handler/ -run TestAiExportHandler -count=1 -v` |
| EC-W07-013 | GET /api/ai-export/download/{id} streams existing ZIP file with correct Content-Type | AC-024 | `cd apps/api && go test ./internal/handler/ -run TestAiExportHandler.*Download -count=1 -v` |
| EC-W07-014 | POST /api/ai-export and GET /api/ai-export/download return 401 when PIN guard is enabled and no valid session | RULE-023, AC-110 | `cd apps/api && go test ./internal/handler/ -run TestAiExportHandler.*Auth -count=1 -v` |
| EC-W07-015 | sqlc generation produces correct query methods for user_profiles and ai_exports tables | Codegen | `bunx nx run api:codegen && bunx nx build api` |
| EC-W07-016 | Log markers use [UserProfile] and [AiExport] prefixes; export content and photos not logged | AC-118, AC-119, AC-120 | `cd apps/api && go test ./internal/atlas/service/ -run "TestUserProfile.*Logs|TestAiExport.*Logs" -count=1 -v` |
| EC-W07-017 | GET /api/user-profile returns current profile with defaults when no profile exists | AC-085 | `cd apps/api && go test ./internal/handler/ -run TestUserProfileHandler -count=1 -v` |
| EC-W07-018 | Migration for user_profiles and ai_exports tables applies cleanly against test PostgreSQL | Migration | `cd apps/api && INTEGRATION_TESTS=1 go test ./internal/atlas/repository/postgres/ -run TestWave07 -count=1 -v` |
| EC-W07-019 | GraphQL resolvers for user profile and ai export (if used) pass schema validation | Schema | `bunx nx run graphql:validate && bunx nx run api:codegen` |
| EC-W07-020 | ZIP cleanup removes temp files after download or error; disk-full during ZIP generation returns 500 without partial file | EDGE-024 | `cd apps/api && go test ./internal/atlas/service/ -run TestAiExport.*Cleanup|TestAiExport.*DiskFull -count=1 -v` |

---

## 2. Verification Obligations (TEST-W07-XXX)

Test IDs with type, file scope, and suggested command.

### 2.1 UserProfile Service Tests

| ID | Type | Scope | Suggested Command | Blocking |
|---|---|---|---|---|
| TEST-W07-001 | unit | service/user_profile.go | `go test ./internal/atlas/service/ -run TestUserProfileService_Create_Success -count=1 -v` | yes |
| TEST-W07-002 | unit | service/user_profile.go | `go test ./internal/atlas/service/ -run TestUserProfileService_Create_EmptyGoal -count=1 -v` | yes |
| TEST-W07-003 | unit | service/user_profile.go | `go test ./internal/atlas/service/ -run TestUserProfileService_Create_InvalidHeight -count=1 -v` | yes |
| TEST-W07-004 | unit | service/user_profile.go | `go test ./internal/atlas/service/ -run TestUserProfileService_Update_Success -count=1 -v` | yes |
| TEST-W07-005 | unit | service/user_profile.go | `go test ./internal/atlas/service/ -run TestUserProfileService_Get_ReturnsDefaultsWhenNoProfile -count=1 -v` | yes |
| TEST-W07-006 | unit | service/user_profile.go | `go test ./internal/atlas/service/ -run TestUserProfileService_Get_ReturnsExistingProfile -count=1 -v` | yes |
| TEST-W07-007 | unit | service/user_profile.go | `go test ./internal/atlas/service/ -run TestUserProfileService_Logs_NoGoalInLog -count=1 -v` | yes |

### 2.2 AiExport Service Tests — Prompt

| ID | Type | Scope | Suggested Command | Blocking |
|---|---|---|---|---|
| TEST-W07-008 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GeneratePrompt_AllSections -count=1 -v` | yes |
| TEST-W07-009 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GeneratePrompt_OnlyWorkouts -count=1 -v` | yes |
| TEST-W07-010 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GeneratePrompt_WithPersistentContext -count=1 -v` | yes |
| TEST-W07-011 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GeneratePrompt_WithOneTimeComment -count=1 -v` | yes |
| TEST-W07-012 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GeneratePrompt_WithWeekFlags -count=1 -v` | yes |
| TEST-W07-013 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GeneratePrompt_EmptyDateRange -count=1 -v` | yes |
| TEST-W07-014 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GeneratePrompt_NoDataInPeriod -count=1 -v` | yes |

### 2.3 AiExport Service Tests — ZIP Generation

| ID | Type | Scope | Suggested Command | Blocking |
|---|---|---|---|---|
| TEST-W07-015 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GenerateZIP_ValidArchive -count=1 -v` | yes |
| TEST-W07-016 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GenerateZIP_ManifestStructure -count=1 -v` | yes |
| TEST-W07-017 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GenerateZIP_DataJSONStructure -count=1 -v` | yes |
| TEST-W07-018 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GenerateZIP_SummaryMDContent -count=1 -v` | yes |
| TEST-W07-019 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GenerateZIP_CSVFilesExist -count=1 -v` | yes |
| TEST-W07-020 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GenerateZIP_CSVHeadersAndRows -count=1 -v` | yes |
| TEST-W07-021 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GenerateZIP_PhotosIncluded -count=1 -v` | yes |
| TEST-W07-022 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GenerateZIP_PhotosExcludedByDefault -count=1 -v` | yes |
| TEST-W07-023 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GenerateZIP_NoPhotosDirWhenOptedOut -count=1 -v` | yes |
| TEST-W07-024 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_GenerateZIP_WorkoutData -count=1 -v` | yes |

### 2.4 AiExport Service Tests — Lifecycle

| ID | Type | Scope | Suggested Command | Blocking |
|---|---|---|---|---|
| TEST-W07-025 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_Cleanup_RemovesTempFiles -count=1 -v` | yes |
| TEST-W07-026 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_Cleanup_RemovesOrphanedExports -count=1 -v` | no |
| TEST-W07-027 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_DiskFull_ReturnsError -count=1 -v` | no |
| TEST-W07-028 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_LargeDateRange_Completes -count=1 -v` | no |
| TEST-W07-029 | unit | service/ai_export.go | `go test ./internal/atlas/service/ -run TestAiExportService_Logs_NoExportContent -count=1 -v` | yes |

### 2.5 REST Handler Tests

| ID | Type | Scope | Suggested Command | Blocking |
|---|---|---|---|---|
| TEST-W07-030 | integration | handler/ai_export.go | `go test ./internal/handler/ -run TestAiExportHandler_GenerateExport -count=1 -v` | yes |
| TEST-W07-031 | integration | handler/ai_export.go | `go test ./internal/handler/ -run TestAiExportHandler_DownloadExport -count=1 -v` | yes |
| TEST-W07-032 | integration | handler/ai_export.go | `go test ./internal/handler/ -run TestAiExportHandler_GenerateExport_MissingAuth -count=1 -v` | yes |
| TEST-W07-033 | integration | handler/ai_export.go | `go test ./internal/handler/ -run TestAiExportHandler_DownloadExport_NotFound -count=1 -v` | yes |
| TEST-W07-034 | integration | handler/user_profile.go | `go test ./internal/handler/ -run TestUserProfileHandler_Get -count=1 -v` | yes |
| TEST-W07-035 | integration | handler/user_profile.go | `go test ./internal/handler/ -run TestUserProfileHandler_Get_NoSession -count=1 -v` | yes |

### 2.6 Repository Integration Tests

| ID | Type | Scope | Suggested Command | Blocking |
|---|---|---|---|---|
| TEST-W07-036 | integration | repository/postgres/ | `INTEGRATION_TESTS=1 go test ./internal/atlas/repository/postgres/ -run TestWave07UserProfileRepo -count=1 -v` | yes |
| TEST-W07-037 | integration | repository/postgres/ | `INTEGRATION_TESTS=1 go test ./internal/atlas/repository/postgres/ -run TestWave07AiExportRepo -count=1 -v` | yes |
| TEST-W07-038 | integration | repository/postgres/ | `INTEGRATION_TESTS=1 go test ./internal/atlas/repository/postgres/ -run TestWave07Migration -count=1 -v` | yes |

### 2.7 Codegen Drift Checks

| ID | Type | Scope | Suggested Command | Blocking |
|---|---|---|---|---|
| TEST-W07-039 | build | sqlc | `bunx nx run api:codegen && bunx nx build api` | yes |
| TEST-W07-040 | build | gqlgen | `bunx nx run graphql:validate && bunx nx run api:codegen && bunx nx build api` | yes |

### 2.8 Resolver Tests (if GraphQL used)

| ID | Type | Scope | Suggested Command | Blocking |
|---|---|---|---|---|
| TEST-W07-041 | unit | graph/resolver/ | `go test ./internal/atlas/graph/resolver/ -run TestUserProfileResolver -count=1 -v` | yes |
| TEST-W07-042 | unit | graph/resolver/ | `go test ./internal/atlas/graph/resolver/ -run TestAiExportResolver -count=1 -v` | yes |

---

## 3. Test Coverage Recommendations

### 3.1 Module Boundary Layout

```
apps/api/internal/atlas/
├── models/
│   ├── user_profile.go        — NEW: UserProfileRecord, UserProfile, UserProfileInput, enums
│   └── ai_export.go            — NEW: AiExportRecord, AiExport, AiExportInput, ExportSection toggles
├── service/
│   ├── user_profile.go         — NEW: Create, Update, Get with validation + defaults
│   ├── user_profile_test.go    — NEW: mock repo, success/validation/default/log tests
│   ├── ai_export.go            — NEW: GeneratePrompt, GenerateZIP, Download, Cleanup
│   └── ai_export_test.go       — NEW: mock repo, prompt/ZIP/cleanup/disk-full tests
├── repository/postgres/
│   ├── queries/
│   │   ├── user_profiles.sql   — NEW: sqlc queries
│   │   └── ai_exports.sql      — NEW: sqlc queries
│   ├── migrations/
│   │   └── 0009x_user_profiles_ai_exports.sql  — NEW: tables
│   ├── user_profile_repo.go    — NEW
│   ├── user_profile_repo_test.go  — NEW: integration
│   ├── ai_export_repo.go       — NEW
│   └── ai_export_repo_test.go  — NEW: integration
├── graph/
│   ├── schema/
│   │   ├── user_profile.graphql  — NEW (if GraphQL)
│   │   └── ai_export.graphql     — NEW (if GraphQL)
│   └── resolver/
│       ├── user_profile.go       — NEW
│       ├── user_profile_test.go  — NEW
│       ├── ai_export.go          — NEW
│       └── ai_export_test.go     — NEW
└── handler/
    ├── ai_export.go              — NEW: POST /api/ai-export, GET /api/ai-export/download/{id}
    ├── ai_export_test.go         — NEW: handler tests
    ├── user_profile.go           — NEW: GET /api/user-profile (existing or new)
    └── user_profile_test.go      — NEW: handler tests
```

### 3.2 Key Test Patterns (from WAVE-05/WAVE-06)

- **Service tests**: Use testify + mock repo interface (see `exercise_service_test.go`, `settings_service_test.go`)
- **Mock repo pattern**: Embed real repo interface, override individual function fields
- **Integration tests**: Use `INTEGRATION_TESTS=1` guard and test PostgreSQL via docker-compose.test.yml
- **Handler tests**: Use httptest + middleware setup (see `handler/users_test.go`)
- **Log privacy tests**: Buffer zap logger, assert export content is absent from output
- **Migration tests**: Apply migration against test DB, verify columns, rollback

### 3.3 Priority Order

1. Models + sqlc queries + migrations (TEST-W07-036 to TEST-W07-040)
2. UserProfile service (TEST-W07-001 to TEST-W07-007)
3. AiExport prompt service (TEST-W07-008 to TEST-W07-014)
4. AiExport ZIP service (TEST-W07-015 to TEST-W07-024)
5. Handler integration (TEST-W07-030 to TEST-W07-035)
6. AiExport lifecycle/cleanup (TEST-W07-025 to TEST-W07-029)
7. GraphQL resolvers if used (TEST-W07-041, TEST-W07-042)

### 3.4 Files Excluded from Handwritten Coverage

Generated sqlc files under `apps/api/internal/atlas/repository/postgres/generated/` and gqlgen generated files under `apps/api/internal/atlas/graph/generated/` should be added to the coverage allowlist with replacement gates for codegen + build + focused integration tests.

### 3.5 Risks

1. **Large date range performance** (TEST-W07-028) — non-blocking; acceptable if response time is <30s for 2-year range
2. **Disk full during ZIP** (EDGE-024, TEST-W07-027) — non-blocking but recommended; implement filesystem space check or handle write errors gracefully
3. **ZIP cleanup** (TEST-W07-025, TEST-W07-026) — must not leak temp files on success paths; orphaned export cleanup is nice-to-have
4. **Empty date range with week flags** — week flags have no date range dependency; test separately that an export with flags but no data still produces valid output
5. **Photo file paths in ZIP** — photos are file paths, not binary content; ZIP must include the actual photo files from the file store, or skip gracefully if file is missing from disk
6. **No technical-verified docs** — API contracts for export download (streaming vs temp file vs redirect) are not formally specified

### 3.6 Questions for Technical Planner

- Q-TC-W07-001: Should POST /api/ai-export return immediately (async) or wait for ZIP generation (sync)?
- Q-TC-W07-002: Where is the export ZIP stored (local disk, configurable path, temp dir)?
- Q-TC-W07-003: Is AiExport exposed via GraphQL or only REST?
- Q-TC-W07-004: How are photos stored and referenced — absolute file paths or relative to a media root?
- Q-TC-W07-005: What app version value is written into manifest.json — hardcoded, build-time, or from config?
- Q-TC-W07-006: Should CSV generation use a shared CSV writer or be inline in the export service?
