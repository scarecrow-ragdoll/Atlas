# Actors And Permissions

## Actors

| Actor | Description | Scope |
| --- | --- | --- |
| User | Single person who owns and uses the Atlas instance | MVP |
| Deployer | Person who sets up the Docker instance (may be same as User) | Infrastructure only |

## Roles

MVP has no role system. The single user is the only actor with all permissions.

| Role | Description |
| --- | --- |
| User (sole occupant) | Full access to all features, data, and settings |

## Permissions Matrix

The single user has all permissions. No permission restrictions exist beyond the optional PIN guard.

### Derived Permissions (all high confidence)

| Permission | Source Evidence | Derivation Rationale |
| --- | --- | --- |
| Create exercise | PRD §11.1, §26.1 | User creates exercises |
| Edit exercise | PRD §11.4 | User can edit standard exercises for themselves |
| Delete exercise | PRD §11.3 | Media can be deleted, implied exercise deletion |
| Upload exercise media | PRD §11.3, §26.1 | Images/video attached to exercise |
| Delete exercise media | PRD §11.3 | Media can be removed |
| View exercise library | PRD §11 | Exercise library is a main section |
| Create workout day | PRD §10.2, §26.2 | Workout created on first save for a date |
| Edit workout day | PRD §26.2, §26.3 | User adds exercises, sets, comments |
| Add sets to exercise | PRD §10.5, §26.2 | Sets with weight and reps |
| Add comments to exercise | PRD §10.4 | Exercise comments for AI export |
| Log cardio | PRD §12, §26.4 | Cardio entry per date |
| Log body weight | PRD §13.5, §26.6 | Weight entry per date |
| Create weekly check-in | PRD §13.2, §26.5 | Check-in with weight, measurements, photos |
| Add body measurements | PRD §13.3 | 10 measurement types |
| Upload progress photos | PRD §14, §26.5 | Photos tied to check-in |
| Create nutrition product | PRD §15.2 | Manual product creation |
| Edit nutrition product | PRD §15.2 | User manages products |
| Delete nutrition product | PRD §15.2 | Implied from CRUD |
| Create nutrition template | PRD §15.3, §26.7 | Weekly template |
| Edit nutrition template | PRD §15.3 | Template items editable |
| Override daily nutrition | PRD §15.5, §26.8 | Add/remove/modify products per day |
| View progress charts | PRD §16 | Training, body, nutrition charts |
| Generate AI prompt | PRD §18, §26.9 | Prompt builder |
| Export AI data (ZIP) | PRD §17, §26.9 | ZIP with JSON, CSVs, summary |
| Save AI review | PRD §19, §26.10 | Manual AI response entry |
| Export full backup | PRD §20.1, §26.11 | Full backup ZIP |
| Import backup | PRD §20.4, §26.12 | Dry-run + full restore |
| Enable/disable PIN | PRD §7.2 | Optional PIN guard |
| Change PIN | PRD §7.2 | PIN change allowed |
| View settings | PRD §8 | Settings section |
| Edit settings | PRD §25.1 | Settings includes units, PIN config |
| View/edit user profile | PRD §18.2, §25.2 | Goal, AI context, personal data |

## Ownership Rules

- All data belongs to the single user
- No data sharing, visibility rules, or multi-tenancy
- Full export/import enables data portability per §6.3

## Privacy And Security Expectations

- PIN is optional; when enabled, all pages require valid session
- PIN stored as hash, not plaintext (§7.2)
- No PIN logging (§24.1)
- No AI export content logging (§24.1)
- No photo logging (§24.1)
- No sensitive comment logging (§24.1)
- Media files not accessible without valid PIN session (§24.1, §14.3)
- Backup and AI export generated only on user request (§24.1)
- Photos included in AI export only when user explicitly opts in (§14.3, §17.3)