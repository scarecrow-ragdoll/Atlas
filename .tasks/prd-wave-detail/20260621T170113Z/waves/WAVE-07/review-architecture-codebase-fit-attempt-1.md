# Review Report: Architecture & Codebase Fit — WAVE-07

**Run**: `20260621T170113Z` | **Wave**: WAVE-07 (AI Export and Prompt Builder)
**Reviewer Role**: architecture-codebase-fit | **Attempt**: 1
**Verdict**: **needs-revision**

---

## 1. Model Pattern Fit (PLANNER: architecture-codebase — SLICES 002, 008)

**APPROVED.** The `UserProfileRecord`/`UserProfile`/`UserProfileInput` triple and `AiExportRecord`/`AiExport`/`CreateAiExportInput` triple follow the established pattern from `models/week_flag.go` and `models/settings.go`. Result/error types (`UserProfileResult`, `UserProfileNotFoundErr`, etc.) match the union-result convention. `AiExportDataProvider` as a Go interface lives in the right layer (service package, not leaked to resolvers).

**One minor naming note:** The architecture-codebase planner uses `XxxValidationErr`, `XxxNotFoundErr`, `XxxAuthErr` (WeekFlag pattern), while Settings uses a single `XxxError` type. Both exist in the codebase; the WeekFlag pattern is the more granular and more recently adopted convention. The planner correctly follows this.

---

## 2. Repository Pattern Fit (SLICES 004, 010)

**APPROVED.** Both `UserProfileRepository` and `AiExportRepository` follow the established pattern from `week_flag_repo.go`:
- Interface with `ctx, userID, ...` parameter ordering
- Private `struct` + `New*Repository(pool *pgxpool.Pool)` constructor
- Dependency on `*generated.Queries`
- Uses `uuidFromString()`, `parseTwoUUIDs()`, `nullableText()`, `modelsToPGDate()` helpers
- `*RecordFromRow()` conversion function returning `*models.XxxRecord`

---

## 3. Service Pattern Fit (SLICES 005, 011, 012)

**APPROVED.** Both `UserProfileService` and `AiExportService` follow the pattern from `service/week_flag.go`/`service/body_weight.go`:
- Interface + private struct + constructor
- Sentinel error vars (`ErrUserProfileNotFound`, etc.)
- `FromRecord()` conversions
- `fmt.Errorf("service_name.method: %w", err)` error wrapping
- Validation at the service boundary before delegating to repo

**The `AiExportDataProvider` interface** (SLICE-011) is a well-designed seam. It wraps 6+ data source dependencies behind a single injectable interface, keeping the `AiExportService` constructor signature manageable. This follows the same principle as how `NutritionMacroService` aggregates across multiple repos in WAVE-05.

**ZIP generation in `service/export_zip.go`** (SLICE-012): Keeping archive building as in-memory structs (`ExportArchive`) in the service package with no filesystem dependencies is sound. The service layer handles actual file I/O, maintaining testability.

---

## 4. Resolver Pattern Fit (SLICES 006, 013)

**APPROVED.** Both `UserProfile` and `AiExport` resolvers follow the established pattern from `resolver/week_flag.go`:
- `middleware.GetAtlasUserID(ctx)` auth extraction
- Union result types with pointer fields for auth/not-found/validation errors
- Service method delegation with error-to-result mapping
- `nil, nil` fallback for unknown errors (consistent with existing codebase)

The new services must be added to `resolver.go`:
```go
UserProfileService  service.UserProfileService
AiExportService     service.AiExportService
```

---

## 5. GraphQL Schema Pattern Fit (SLICES 006, 013)

**APPROVED.** The proposed schemas follow the `week_flag.graphql` pattern:
- Types with `ID!`, scalar fields, nullable fields marked with `String`/`Float`
- Input types for mutations
- Result types with inline error types (`XxxValidationError`, `XxxNotFoundError`, `XxxAuthError`)
- `XxxErrorCode` enum with `VALIDATION_ERROR`, `NOT_FOUND`, `AUTH_ERROR`, `INTERNAL_ERROR`
- Plural result type (`AiExportsResult`) matching `WeekFlagsResult` pattern

**Note:** The existing codebase has two error patterns:
- **Settings** (`settings.graphql`): single `SettingsError` type with field `error: SettingsError`
- **WeekFlag** (`week_flag.graphql`): separate inline types per error

The planner follows the WeekFlag pattern, which is correct for the current convention.

---

## 6. REST + GraphQL Boundary (SLICE-014)

**APPROVED.** The hybrid approach (GraphQL for mutations/queries, REST only for binary ZIP download) follows the existing `ProgressPhotoHandler.Download` pattern. The `AtlasPinGuard` middleware protection is correctly specified.

---

## 7. Migration Numbering — NEEDS REVISION

**ISSUE.** Current max migration is `00090_nutrition_tables.sql`. Two planners disagree:

| Planner | user_profiles | ai_exports |
|---------|---------------|------------|
| architecture-codebase | `00093` | `00094` |
| data-integration-ops | `00091` | `00092` |

Correct next numbers are **00091** and **00092**. The architecture-codebase planner skips two numbers (00091, 00092) unnecessarily. **Fix:** Change architecture-codebase SLICE-001 to `00091_user_profiles.sql` and SLICE-007 to `00092_ai_exports.sql`.

---

## 8. gqlgen Config — NEEDS REVISION (gap)

**ISSUE.** The architecture-codebase planner does not mention updating `atlas-gqlgen.yml`. All new GraphQL types require explicit model binding entries. The following bindings must be added:

```
UserProfile
UserProfileInput
UserProfileResult
UserProfileValidationError
UserProfileNotFoundError
UserProfileAuthError
UserProfileErrorCode
AiExport
CreateAiExportInput
AiExportResult
AiExportsResult
AiExportValidationError
AiExportNotFoundError
AiExportAuthError
AiExportErrorCode
```

**Fix:** Add a slice or explicit step to update `apps/api/atlas-gqlgen.yml` and run `bun run codegen` after schema changes.

---

## 9. Planner Inconsistency: `display_name` on UserProfile

**ISSUE.** The data-integration-ops planner includes a `display_name TEXT NOT NULL DEFAULT ''` column on the `user_profiles` table. The architecture-codebase planner does not include this field in the migration, model, or GraphQL schema. If the frontend page-009 expects a display name on the profile endpoint, this field is needed.

**Recommendation:** Either add `display_name` to the architecture-codebase UserProfile model (read from `atlas_users.display_name` via a join or sync during bootstrap), or confirm the frontend does not need it. If not adding it, document the decision and remove it from the data-integration-ops migration.

---

## 10. Planner Inconsistency: `updated_at` on `ai_exports`

The architecture-codebase planner correctly includes `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()` on the `ai_exports` table. The data-integration-ops planner omits it. The architecture-codebase version is correct — all existing tables have `updated_at`.

---

## 11. File Paths and Module Names

**APPROVED.** All proposed file paths follow existing conventions:

| Proposed path | Existing equivalent | Match |
|---|---|---|
| `models/user_profile.go` | `models/week_flag.go` | ✓ |
| `models/ai_export.go` | `models/settings.go` | ✓ |
| `service/user_profile_service.go` | `service/settings_service.go` | ✓ |
| `service/ai_export_service.go` | `service/settings_service.go` | ✓ |
| `service/export_zip.go` | `service/nutrition_macro_service.go` | ✓ |
| `repository/postgres/user_profile_repo.go` | `repository/postgres/week_flag_repo.go` | ✓ |
| `repository/postgres/ai_export_repo.go` | `repository/postgres/cardio_entry_repo.go` | ✓ |
| `graph/schema/user_profile.graphql` | `graph/schema/week_flag.graphql` | ✓ |
| `graph/schema/ai_export.graphql` | `graph/schema/settings.graphql` | ✓ |
| `graph/resolver/user_profile.go` | `graph/resolver/week_flag.go` | ✓ |
| `graph/resolver/ai_export.go` | `graph/resolver/week_flag.go` | ✓ |
| `handler/ai_export_handler.go` | `handler/progress_photo_handler.go` | ✓ |
| `repository/postgres/queries/user_profiles.sql` | `repository/postgres/queries/week_flags.sql` | ✓ |

---

## 12. sqlc Generated Code Implications

**APPROVED.** The planner correctly accounts for:
- `UNIQUE(user_id)` on `user_profiles` enables the `ON CONFLICT` upsert pattern
- `COALESCE` in upsert queries for nullable column merging
- `RETURNING *` pattern for all write queries
- Date types (`DATE`) → `pgtype.Date`, UUID → `pgtype.UUID`, REAL → `*float32` nullable handling
- The planner explicitly calls out running `bun run codegen` after sqlc query files

---

## 13. Summary of Issues Requiring Revision

| # | Severity | Issue | Affected Slice |
|---|----------|-------|----------------|
| 1 | **High** | Migration numbers wrong: must be `00091`/`00092`, not `00093`/`00094` | SLICE-001, SLICE-007 |
| 2 | **Medium** | gqlgen config update not mentioned — all new types need explicit bindings in `atlas-gqlgen.yml` | SLICE-006, SLICE-013, and a new step |
| 3 | **Medium** | Inconsistency between planners: `display_name` on `user_profiles` present in data-integration-ops but absent in architecture-codebase | SLICE-001, SLICE-002, SLICE-006 |

**All three issues are fixable without restructuring the wave.** The pattern fit, module structure, data flow, and interface design are sound. Revise migration numbers, add the gqlgen config step, and resolve the display_name inconsistency.
