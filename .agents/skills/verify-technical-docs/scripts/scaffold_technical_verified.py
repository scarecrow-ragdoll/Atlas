#!/usr/bin/env python3
# FILE: .agents/skills/verify-technical-docs/scripts/scaffold_technical_verified.py
# VERSION: 1.0.0
# START_MODULE_CONTRACT
#   PURPOSE: Scaffold the managed technical-verified documentation tree for staging or final output.
#   SCOPE: Creates required markdown files and headings; excludes semantic synthesis and question analysis.
#   DEPENDS: Python standard library, .agents/skills/verify-technical-docs/references/output-contract.md.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
#   ROLE: SCRIPT
#   MAP_MODE: EXPORTS
# END_MODULE_CONTRACT
#
# START_MODULE_MAP
#   REQUIRED_DIRS - Required subdirectories for the technical-verified package.
#   main - Parse arguments and write the scaffold.
#   build_files - Return required file content keyed by relative path.
#   source_inventory - List source files for source-inventory.md.
# END_MODULE_MAP
# START_CHANGE_SUMMARY
#   LAST_CHANGE: 1.0.0 - Added scaffold script for technical readiness output.
# END_CHANGE_SUMMARY

from __future__ import annotations

import argparse
from pathlib import Path


REQUIRED_DIRS = [
    "features",
    "appendix",
]


# START_CONTRACT: source_inventory
#   PURPOSE: Return a stable markdown list of source files included in the technical verification run.
#   INPUTS: { source: Path - product-verified input folder }
#   OUTPUTS: { str - markdown bullet list or missing-source marker }
#   SIDE_EFFECTS: None.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: source_inventory
def source_inventory(source: Path) -> str:
    if not source.exists():
        return "- SOURCE_MISSING"
    files = sorted(path for path in source.rglob("*") if path.is_file())
    if not files:
        return "- SOURCE_EMPTY"
    return "\n".join(f"- {path.as_posix()}" for path in files)


# START_CONTRACT: build_files
#   PURPOSE: Build the required technical-verified markdown file set with canonical headings.
#   INPUTS: { source: Path - product-verified input folder }
#   OUTPUTS: { dict[str, str] - relative path to markdown content }
#   SIDE_EFFECTS: None.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: build_files
def build_files(source: Path) -> dict[str, str]:
    inventory = source_inventory(source)
    return {
        "index.md": """# Technical Verified
## Status
questions-open
## Source Set
See `source-inventory.md`.
## Document Map
- Technical brief: `technical-brief.md`
- Open questions: `open-questions.md`
## Dev Handoff Readiness
Not approved until all required scopes and consistency approve with no open blocking questions.
""",
        "source-inventory.md": f"""# Source Inventory
## Included Sources
{inventory}
## Source Delta
PLACEHOLDER
## Answered Questions
PLACEHOLDER
## Excluded Or Noisy Sources
PLACEHOLDER
## Coverage Gaps
PLACEHOLDER
""",
        "technical-brief.md": """# Technical Brief
## Product Signal
PLACEHOLDER
## Technical Scope
PLACEHOLDER
## Constraints
PLACEHOLDER
## Assumptions
PLACEHOLDER
## Readiness Summary
PLACEHOLDER
""",
        "architecture-and-boundaries.md": """# Architecture And Boundaries
## System Context
PLACEHOLDER
## Components
PLACEHOLDER
## Ownership Boundaries
PLACEHOLDER
## Dependencies
PLACEHOLDER
## Architecture Questions
PLACEHOLDER
""",
        "data-contracts.md": """# Data Contracts
## Entities And Identifiers
PLACEHOLDER
## Persistence And Storage
PLACEHOLDER
## Migrations
PLACEHOLDER
## Retention And Privacy
PLACEHOLDER
## Data Questions
PLACEHOLDER
""",
        "api-contracts.md": """# API Contracts
## Surfaces
PLACEHOLDER
## Requests And Responses
PLACEHOLDER
## Error And Validation Contracts
PLACEHOLDER
## Compatibility And Idempotency
PLACEHOLDER
## API Questions
PLACEHOLDER
""",
        "auth-security-compliance.md": """# Auth Security Compliance
## Identity
PLACEHOLDER
## Authorization And Ownership
PLACEHOLDER
## Auditability
PLACEHOLDER
## Privacy And Compliance
PLACEHOLDER
## Security Questions
PLACEHOLDER
""",
        "integrations-and-events.md": """# Integrations And Events
## External Systems
PLACEHOLDER
## Events And Jobs
PLACEHOLDER
## Sync And Retry Rules
PLACEHOLDER
## Rate Limits And Failure Handling
PLACEHOLDER
## Integration Questions
PLACEHOLDER
""",
        "client-state-and-ux-contracts.md": """# Client State And UX Contracts
## User Interface States
PLACEHOLDER
## Form And Validation Behavior
PLACEHOLDER
## Cache And Realtime Behavior
PLACEHOLDER
## Accessibility And Localization
PLACEHOLDER
## Client Questions
PLACEHOLDER
""",
        "operations-observability.md": """# Operations Observability
## Environments And Config
PLACEHOLDER
## Deployment And Rollout
PLACEHOLDER
## Logs Metrics Traces
PLACEHOLDER
## Alerts Runbooks Backups
PLACEHOLDER
## Operations Questions
PLACEHOLDER
""",
        "testing-and-delivery.md": """# Testing And Delivery
## Test Strategy
PLACEHOLDER
## Fixtures And Test Data
PLACEHOLDER
## Contract And E2E Coverage
PLACEHOLDER
## Release Gates
PLACEHOLDER
## Testing Questions
PLACEHOLDER
""",
        "implementation-slices.md": """# Implementation Slices
## Slice Map
PLACEHOLDER
## Dependencies
PLACEHOLDER
## Blockers
PLACEHOLDER
## Verification
PLACEHOLDER
""",
        "open-questions.md": """# Open Questions
## Dev-Blocking
PLACEHOLDER
## Needs Owner Decision
PLACEHOLDER
## Deferred
PLACEHOLDER
## Watchlist
PLACEHOLDER
## Resolved This Run
PLACEHOLDER
""",
        "features/index.md": """# Feature Technical Inventory
## Features
PLACEHOLDER
""",
        "appendix/subagent-findings.md": """# Subagent Findings
## Scope Reports
PLACEHOLDER
## Reviewer Verdicts
| Scope | Status | Reviewer Verdict | Report |
| --- | --- | --- | --- |
## Conflicts
PLACEHOLDER
""",
        "appendix/traceability.md": """# Traceability
## Technical Requirement Map
PLACEHOLDER
## Question Map
PLACEHOLDER
## Decision Map
PLACEHOLDER
## Slice Map
PLACEHOLDER
## Source Map
PLACEHOLDER
""",
        "appendix/question-ledger.md": """# Question Ledger
## Open Questions
| ID | Scope | Severity | Parent | Question | Why It Matters | Needed Artifact Or Decision | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
## Answered Questions
PLACEHOLDER
## Follow-Up Questions
PLACEHOLDER
## Resolved Questions
PLACEHOLDER
## Deferred Questions
PLACEHOLDER
""",
        "appendix/decision-log.md": """# Decision Log
## Technical Decisions
PLACEHOLDER
## Deferrals
PLACEHOLDER
## Superseded Answers
PLACEHOLDER
## Rejected Assumptions
PLACEHOLDER
""",
        "appendix/loop-history.md": """# Loop History
## Runs
PLACEHOLDER
## Answered Question Effects
PLACEHOLDER
## Follow-Up Blockers
PLACEHOLDER
## Approval Gate History
| Gate | Status | Evidence |
| --- | --- | --- |
""",
    }


# START_CONTRACT: main
#   PURPOSE: Create a technical-verified scaffold at the requested output path.
#   INPUTS: { argv: CLI args - --source, --output, --force }
#   OUTPUTS: { int - process exit code }
#   SIDE_EFFECTS: Creates directories and markdown files under output.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: main
def main() -> int:
    parser = argparse.ArgumentParser(description="Scaffold technical-verified docs.")
    parser.add_argument("--source", default="docs/product-verified", help="Verified product docs folder.")
    parser.add_argument("--output", required=True, help="Output folder to create.")
    parser.add_argument("--force", action="store_true", help="Overwrite existing managed files.")
    args = parser.parse_args()

    source = Path(args.source)
    output = Path(args.output)
    output.mkdir(parents=True, exist_ok=True)
    for relative_dir in REQUIRED_DIRS:
        (output / relative_dir).mkdir(parents=True, exist_ok=True)

    for relative_path, content in build_files(source).items():
        target = output / relative_path
        target.parent.mkdir(parents=True, exist_ok=True)
        if target.exists() and not args.force:
            continue
        target.write_text(content, encoding="utf-8")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
