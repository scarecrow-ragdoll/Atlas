# Roles-Permissions-Reviewer Question Ledger

| ID | Scope | Severity | Question | Why It Matters | Source Or Report | Status | Resolution |
|---|---|---|---|---|---|---|---|
| Q-ROLE-001 | roles-permissions-reviewer | non-blocking | What is the PIN session lifetime and renewal policy? | Affects user experience and security boundary. | `worker-attempt-1.md` | open | TBD |
| Q-ROLE-002 | roles-permissions-reviewer | non-blocking | Is there a logout mechanism when PIN is enabled? | Affects session management implementation. | `worker-attempt-1.md` | open | TBD |
| Q-ROLE-003 | roles-permissions-reviewer | non-blocking | What happens when PIN is disabled — is there any remaining access control at all? | Affects the entire security model. | `worker-attempt-1.md` | open | TBD |
| Q-ROLE-004 | roles-permissions-reviewer | deferred | Are there any plans for resource-level visibility control (e.g., hide specific workouts)? | Not mentioned in MVP but affects architecture. | `worker-attempt-1.md` | deferred | Future scope consideration |
| Q-ROLE-005 | roles-permissions-reviewer | non-blocking | How does the deployer role work in practice — is there a setup wizard, env config, or seed data step? | Missing from PRD but necessary for self-hosted deployment. | `worker-attempt-1.md` | open | TBD |