# Roles-Permissions-Reviewer Worker Attempt 1

## Sources Read

- `docs/product/prd.md` (lines 1–1562+ read)

## Source Delta Reviewed

No source delta present.

## Confirmed Facts

### Actor
- **Single user** (lines 44, 159, 1087): MVP is single-user, self-hosted, no registration, no multi-user. One person is both the operator and the data owner.

### Access Control
- **Optional PIN-code** (lines 171–184): Can be enabled/disabled. When enabled, PIN entry is required before accessing the app. PIN stored as hash. Session via cookie after PIN entry. Sensitive data not accessible without valid session.
- **No registration** (lines 44–46, 163–164): No sign-up, no user management, no roles, no invitations.

### Data Ownership
- **User owns all data** (lines 145–153): Full export, backup, restore, import. No shared ownership in MVP.
- **Data belongs to the single user** of the instance (lines 145–153, 1063–1074).

### Privacy / Visibility Rules
- **Media files not publicly accessible without PIN session** (line 535, 1073).
- **Photos not included in AI export by default** (line 712), only if user opts in.
- **Sensitive data not logged** (lines 1069–1074): PIN, AI export contents, photos, sensitive comments.
- **Full backup and AI export only on user request** (line 1074).

### Responsibility Boundaries
- **User creates exercises** (line 347): Self-created exercise library, no starter catalog in MVP.
- **User creates nutrition products** (line 553): Manual product creation.
- **User sets working weight** (lines 300–301): Stored in exercise reference, auto-populated in workouts.
- **User manages nutrition template** (lines 566–568): Weekly template, daily overrides.
- **User performs weekly check-in** (lines 449–458): Body measurements, photos, weight.
- **User generates AI export** (lines 1436–1446): Chooses period, selects data, generates prompt + ZIP.
- **User saves AI review** (lines 1449–1456): Pastes AI response, links to period, notes planned actions.
- **User performs backup/restore** (lines 1458–1474): Full export, import with dry-run validation.
- **User manages settings** (lines 1125–1131): PIN enable/disable, units, default export weeks.

### No Multi-Tenant / Admin Roles
- **No admin role** (lines 44–48, 165): Single user is the only actor. No need for admin/user distinction.
- **No public profiles** (line 167): No external visibility.
- **Future architecture may support multi-user** (line 48), but not in MVP.

## Derived Roles

| Role | Source Signal | Derivation Rationale | Confidence |
|---|---|---|---|
| **Instance Owner** (deployer) | Lines 43–48: self-hosted, one user, deployed by another person | The person who deploys Atlas creates an instance. This is a deployment-time role, not an application-level role. They configure the environment but have no separate application identity. | **Low** — The PRD mentions deployment by "another person" but gives no behavior for the deployer beyond standing up the instance. This role has no in-app actions. |
| **Single User** (end user) | Lines 44, 159, 1087: one user, no registration | The sole application user. Owns all data, performs all actions, controls all settings. | **High** — Explicitly stated and reinforced throughout. |
| **Instance Operator** (deployer-as-user) | Lines 43–48: deployed by another person, but used by one person | If the deployer and user are the same person, they are the same role. If the deployer runs it for someone else, the deployer role exists at infrastructure level only. | **Low** — No application-level behavior is defined for the deployer. |
| **PIN-holder** (session-authenticated user) | Lines 171–184: PIN opt-in, cookie session | When PIN is enabled, the user becomes a "PIN-holder" with an authenticated session. This is a session state, not a separate role. | **Medium** — The PIN creates an access boundary, but there is only one user, so the PIN-holder is the same person as the single user. |

## Derived Permissions

All permissions belong to the **Single User**. There is no permission separation in MVP.

| Permission | Source Signal | Derivation Rationale | Confidence |
|---|---|---|---|
| **Read all data** | Lines 1073: media requires session; 173: PIN protects data access | User with valid session (or PIN disabled) can read all application data. | **High** |
| **Create exercise** | Lines 347, 1356–1364 | Explicit user flow: user creates exercises. | **High** |
| **Edit exercise** | Lines 383–387: user can edit standard exercises for themselves | Derived from "edit standard exercise" language and CRUD context. | **High** |
| **Delete exercise** | Implied by CRUD context (Epic 2, line 1579) | "CRUD упражнений" implies delete permission. | **Medium** — "CRUD" stated but delete not explicitly demonstrated in user flows. |
| **Upload exercise media** | Lines 366–374, 1363 | Explicit flow: user uploads images/video. | **High** |
| **Delete exercise media** | Line 372: media can be deleted | Explicit statement. | **High** |
| **Log workout** | Lines 1368–1376, 1378–1385 | Explicit user flows for logging workouts. | **High** |
| **Edit workout** | Line 242: can edit backdated workouts | User can edit existing workout entries. | **High** |
| **Delete workout** | Implied by edit context | Not explicitly stated but strongly implied by data ownership. | **Low** — No explicit delete flow for workouts described. |
| **Add sets to exercise** | Lines 1372–1374, 279–291 | Explicit flow: user adds sets with weight, reps, optional RPE/RIR. | **High** |
| **Log cardio** | Lines 1388–1395 | Explicit user flow. | **High** |
| **Create body check-in** | Lines 1398–1407 | Explicit user flow. | **High** |
| **Log body weight** | Lines 1409–1414 | Explicit user flow. | **High** |
| **Upload progress photos** | Lines 1406, 508–517 | Check-in includes 2–4 photos. | **High** |
| **Create nutrition product** | Lines 1417–1419, 553 | Explicit flow: user creates products. | **High** |
| **Edit nutrition product** | Implied by CRUD context | User manages product catalog. | **Medium** — No explicit edit flow, but CRUD is implied. |
| **Create weekly nutrition template** | Lines 566–568, 1417–1425 | Explicit user flow. | **High** |
| **Edit nutrition template** | Lines 596–603: daily overrides | Template can be overridden per day. | **High** |
| **Apply daily override** | Lines 596–603 | Explicit flow: add, remove, change product per day. | **High** |
| **View charts** | Lines 624–671 | Explicit feature: progress charts for exercises, body, nutrition. | **High** |
| **Generate AI export** | Lines 1436–1446, 674–731 | Explicit user flow. | **High** |
| **Generate AI prompt** | Lines 1443, 795–873 | Explicit system behavior. | **High** |
| **Save AI review** | Lines 1449–1456, 888–897 | Explicit user flow. | **High** |
| **Export full backup** | Lines 1458–1464, 902–932 | Explicit user flow. | **High** |
| **Import full backup** | Lines 1466–1474, 972–986 | Explicit user flow. | **High** |
| **Enable/disable PIN** | Lines 177–182 | Explicit setting. | **High** |
| **Change PIN** | Line 181 | Explicit requirement. | **High** |
| **Access app without PIN** | Lines 177–178: PIN is optional | When PIN is disabled, no authentication needed. | **High** |
| **Read all data without PIN** | Line 178: app opens without check when PIN disabled | Full read access to all data when PIN is off. | **High** |
| **Change settings** | Lines 1125–1131 (Settings model), 177–182 (PIN) | User controls app settings. | **High** |
| **Update AI context/goal** | Lines 802–818 | User can change goal and persistent AI context anytime. | **High** |

## Permissions Matrix

Since MVP has **exactly one user** and **no roles**, the permissions matrix is trivial: all permissions belong to the single user. No role-based access control (RBAC) exists.

| Actor | App Access | Data CRUD | PIN Mgmt | Export/Import | Settings | Charts |
|---|---|---|---|---|---|---|
| Single User (no PIN) | Full | Full | N/A (PIN disabled) | Full | Full | Full |
| Single User (PIN enabled, session valid) | Full | Full | Full | Full | Full | Full |
| Single User (PIN enabled, no session) | Login page only | None | None | None | None | None |
| Unknown visitor (PIN disabled) | Full | Full | N/A | Full | Full | Full |
| Deployer (separate person) | No in-app access | None | None | None | None | None |

## Ownership Rules

| Rule | Source Signal | Derivation Rationale | Confidence |
|---|---|---|---|
| User owns all data | Lines 145–153: data belongs to user, full export/backup/restore | Explicit. | **High** |
| Data is tied to the instance, not a user account | Lines 43–48: no registration, single user | No user identity beyond the instance itself. | **High** |
| No shared ownership | Lines 44–46: single user, no multi-user | By definition, no sharing needed. | **High** |
| Backup export includes all owned data | Lines 902–932, 1458–1464 | Full backup is all-encompassing. | **High** |
| Backup import replaces all data in target instance | Lines 972–986, 1466–1474 | Import restores full state, not a merge. | **High** |

## Approval Authority

There is **no approval authority** in MVP. No workflows require secondary approval:
- No moderation
- No content review
- No permission granting
- No multi-step approval

The single user is the sole decision-maker for all actions.

## Visibility Rules

| Rule | Source Signal | Derivation Rationale | Confidence |
|---|---|---|---|
| All data visible to user with valid session (or PIN disabled) | Lines 173–178, 184 | PIN protects access, but once granted, all data is visible. | **High** |
| Media files require PIN session | Line 535, 1073 | Explicitly stated. | **High** |
| No public visibility of any data | Lines 167: no public profiles; 535: media not public | No data is exposed publicly. | **High** |
| Photos excluded from AI export by default | Line 712 | Opt-in only. | **High** |
| Full backup / AI export generated only on user request | Line 1074 | Explicit privacy requirement. | **High** |

## Responsibility Boundaries

| Boundary | Source | Derivation Rationale | Confidence |
|---|---|---|---|
| User creates and manages exercise library | Lines 347, 1356–1364 | Explicit. | **High** |
| User creates and manages nutrition products | Lines 553, 1417–1419 | Explicit. | **High** |
| User manages weekly nutrition template | Lines 566–568, 1417–1425 | Explicit. | **High** |
| User tracks body measurements and photos | Lines 443–490, 1398–1407 | Explicit. | **High** |
| User generates AI exports and reviews | Lines 674–897, 1436–1456 | Explicit. | **High** |
| User performs backup and restore | Lines 902–986, 1458–1474 | Explicit. | **High** |
| User controls PIN and settings | Lines 171–184, 1125–1131 | Explicit. | **High** |
| Deployer sets up infrastructure | Lines 43–48: self-hosted | Implied by self-hosted nature. | **Low** — No explicit deployer responsibilities defined. |
| Application protects sensitive data | Lines 1063–1074 | Explicit privacy requirements (no logging, session guard). | **High** |

## Contradictions

No contradictions found within the source. The single-user model is consistent throughout.

## Missing Source Artifacts

1. **No formal authorization policy**: The PRD describes PIN-based access but does not define a formal authorization model (ACL, RBAC, capabilities). The confidence level for any derived permission is subject to this gap.
2. **No session / token specification**: Session lifetime, renewal, CSRF protection, token storage are not described.
3. **No data retention or compliance policy**: No mention of data deletion, export retention, GDPR/privacy compliance.
4. **No deployer documentation**: The deployer role is mentioned but has no documented responsibilities, setup flow, or configuration surface.

## Derived Requirements

No new product behavior requirements are derived. Permissions and rules are derived as documented above.

## Missing Information

1. How is the PIN session managed? (lifetime, renewal, logout, concurrent sessions)
2. Is there a logout action, or does the session persist indefinitely?
3. When PIN is disabled, is there any authentication at all?
4. Are there any sub-resource visibility rules (e.g., hiding specific data within the app)?
5. What happens to data when the user wants to delete specific records (explicit delete flows)?

## Open Questions Raised

| ID | Question | Why It Matters |
|---|---|---|
| Q-ROLE-001 | What is the PIN session lifetime and renewal policy? | Affects user experience and security boundary. |
| Q-ROLE-002 | Is there a logout mechanism when PIN is enabled? | Affects session management implementation. |
| Q-ROLE-003 | What happens when PIN is disabled — is there any remaining access control at all? | Affects the entire security model. |
| Q-ROLE-004 | Are there any plans for resource-level visibility control (e.g., hide specific workouts)? | Not mentioned in MVP but affects architecture. |
| Q-ROLE-005 | How does the deployer role work in practice — is there a setup wizard, env config, or seed data step? | Missing from PRD but necessary for self-hosted deployment. |

## Edge Cases Or Risks

| Edge Case | Risk |
|---|---|
| PIN lockout: user forgets PIN with no recovery mechanism | Data loss / permanent lockout. PRD does not mention PIN recovery. |
| Concurrent browser sessions with PIN | Not described. |
| Deployer and user are different people | No documented access model for this scenario. |
| Browser cookie cleared while PIN is enabled | Session lost; user must re-enter PIN. Acceptable but not documented. |
| Data accessibility during backup import | Import replaces all data; no mention of downtime or read-only mode. |

## Recommended Decisions

1. Document PIN session policy (lifetime, renewal, logout) before implementation.
2. Clarify deployer responsibilities or mark as out-of-scope for MVP.
3. Document whether any access control remains when PIN is disabled.

## Traceability Candidates

Requirements trace to sections 7, 24, 25, and user flows in sections 26 of `docs/product/prd.md`.

- REQ-PIN-ENABLE: Line 177
- REQ-PIN-DISABLE: Line 182
- REQ-PIN-CHANGE: Line 181
- REQ-PIN-SESSION: Lines 183–184
- REQ-DATA-OWNERSHIP: Lines 145–153
- REQ-PRIVACY-NO-LOG: Lines 1069–1074
- REQ-MEDIA-SESSION: Lines 535, 1073
- REQ-PHOTO-EXPORT-OPTIN: Line 712