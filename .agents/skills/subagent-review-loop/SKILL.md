---
name: subagent-review-loop
description: >-
  Run an iterative, subagent-based review loop for specs, design docs, implementation plans, PRDs, task decompositions, architecture notes, diffs, or other artifacts until independent reviewers approve. Use when the user asks to review a plan/spec with subagents, run a review loop, get multi-angle approval, validate and fix a generated artifact, or says in Russian phrases like "запусти ревью луп с субагентами", "отревьювить план субагентами", or "доведи спеку до аппрува ревьюверов." The main session remains the controller: it adjudicates findings, edits the artifact, and re-runs reviewers until approval or an explicit blocker.
---

<!-- FILE: .agents/skills/subagent-review-loop/SKILL.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the project-local skill for iterative artifact review by independent subagents. -->
<!--   SCOPE: Covers reviewer selection, subagent prompts, controller adjudication, artifact edits, and approval-loop closeout; excludes implementation of reviewed plans. -->
<!--   DEPENDS: multi-agent/subagent tool availability, user-provided or repository artifact context. -->
<!--   LINKS: M-SUBAGENT-REVIEW-LOOP-SKILL / V-M-SUBAGENT-REVIEW-LOOP-SKILL. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Workflow - Runs parallel read-only reviewers, controller adjudication, artifact edits, and targeted re-review until approval. -->
<!--   Reviewer Set - Defines default and optional reviewer scopes without overlapping ownership. -->
<!--   Subagent Prompt Template - Gives the reusable prompt contract and verdict format for each reviewer. -->
<!--   Controller Rules - Keeps context ownership, conflict resolution, and validity decisions in the main session. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added compact project-local subagent approval loop skill. -->
<!-- END_CHANGE_SUMMARY -->

# Subagent Review Loop

Use subagents as independent read-only reviewers. Keep all artifact edits in the main session.

## Workflow

1. Identify the review target: a file path, diff, pasted artifact, or the latest artifact produced in the conversation. If the target is ambiguous and cannot be inferred, ask for it.
2. Identify source-of-truth context: user request, requirements, existing docs, code contracts, tests, constraints, and acceptance criteria. Pass only the needed artifacts to reviewers.
3. Spawn reviewers in parallel with the available subagent tools. Do not simulate independent approval. If subagents are unavailable, say the loop cannot run fully.
4. Main session reads all reports, validates each finding against source context, merges duplicates, and rejects invalid or taste-only findings.
5. Main session edits the artifact directly.
6. Re-run only reviewers whose scope was affected or who did not approve. Include the revised artifact and a short resolution log.
7. Stop only when all required reviewers return `APPROVE`, or when the same unresolved external blocker prevents meaningful progress. If the controller rejects a finding, ask that reviewer to re-evaluate the controller rationale rather than silently overriding it.

## Reviewer Set

Default to three required reviewers:

- Intent reviewer: checks whether the artifact preserves the user's goal, scope, acceptance criteria, decisions, and open questions.
- Feasibility reviewer: checks architecture, sequencing, ownership, dependency order, data/API contracts, migration risk, and implementability.
- Verification reviewer: checks tests, evidence, edge cases, observability, rollback, and completion gates.

Add at most two optional reviewers only when the artifact clearly needs them:

- UX/product reviewer for user-facing flows, copy, interaction design, or product positioning.
- Security/privacy reviewer for auth, permissions, secrets, data exposure, compliance, or destructive operations.
- Data/ops reviewer for schema changes, migrations, background jobs, production rollout, infra, or SLO-sensitive paths.
- Domain reviewer for project-specific rules, GRACE/contracts, business process, or API compatibility.

Keep reviewer scopes disjoint. Tell reviewers to stay in lane and mention out-of-scope issues only if they are hard blockers.

## Subagent Prompt Template

Use this shape for each reviewer:

```text
You are the <role> reviewer in an iterative artifact review loop.

Review target:
<path, diff, or pasted artifact>

Source context:
<user goal, requirements, docs, code/test pointers, constraints>

Scope:
<role-specific responsibilities>

Rules:
- Read only and report; do not edit files.
- Stay within your scope.
- Block only on issues that would make the artifact misleading, unsafe, unimplementable, unverifiable, or materially incomplete.
- Do not require stylistic rewrites, preferred phrasing, or nice-to-have additions as blockers.
- Return a structured verdict.

Output:
VERDICT: APPROVE | CHANGES_REQUESTED
ROLE: <role>
BLOCKERS:
- id: <stable-id>
  severity: blocker | major
  evidence: <specific artifact/source reference>
  problem: <what breaks>
  required_change: <minimal fix>
NON_BLOCKING:
- <optional minor note>
APPROVAL_CONDITIONS:
- <conditions that must be true for APPROVE, or "none">
```

## Controller Rules

- Treat subagent approval as necessary but not blindly authoritative. Validate every finding before editing.
- Prefer minimal artifact edits that satisfy the underlying issue.
- Keep a short resolution log: accepted findings, rejected findings with rationale, and changed sections/files.
- Re-dispatch reviewers with the revised artifact and resolution log. Ask them to approve or identify remaining blockers only.
- Do not let reviewers fight over wording. If two valid findings conflict, resolve by source-of-truth priority: explicit user request, durable project docs/contracts, code reality, tests/verification, then reviewer preference.
- If reviewers keep disagreeing after a controller rationale and one re-check, make the conflict explicit to the user instead of fabricating approval.

## Final Response

Report the final artifact location or revised text, reviewer verdicts, key accepted changes, any rejected findings that mattered, and any residual external blockers. Keep it concise.
