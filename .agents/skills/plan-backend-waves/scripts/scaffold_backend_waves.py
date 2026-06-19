#!/usr/bin/env python3
# FILE: .agents/skills/plan-backend-waves/scripts/scaffold_backend_waves.py
# VERSION: 1.0.0
# START_MODULE_CONTRACT
#   PURPOSE: Scaffold the managed backend-waves documentation tree for staging or final output.
#   SCOPE: Creates required markdown files and headings for one current wave; excludes semantic backend decomposition and reviewer synthesis.
#   DEPENDS: Python standard library, .agents/skills/plan-backend-waves/references/output-contract.md.
#   LINKS: M-BACKEND-WAVE-PLANNER / V-M-BACKEND-WAVE-PLANNER.
#   ROLE: SCRIPT
#   MAP_MODE: EXPORTS
# END_MODULE_CONTRACT
#
# START_MODULE_MAP
#   REQUIRED_DIRS - Required subdirectories for the backend-waves package.
#   REVIEWER_ROWS - Canonical scaffold rows for required reviewer perspectives.
#   QUESTION_TABLE - Canonical scaffold table for backend wave questions.
#   main - Parse arguments and write the scaffold.
#   build_files - Return required file content keyed by relative path.
#   source_inventory - List source files for source-inventory.md.
# END_MODULE_MAP
# START_CHANGE_SUMMARY
#   LAST_CHANGE: 1.0.0 - Added scaffold script for backend wave planning output.
# END_CHANGE_SUMMARY

from __future__ import annotations

import argparse
from pathlib import Path


REQUIRED_DIRS = [
    "waves",
    "appendix",
]

REVIEWER_ROWS = """| Wave | Perspective | Attempt | Verdict | Reviewer Report | Required Revisions | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| WAVE-01 | backend-architecture | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| WAVE-01 | data-api-contract | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| WAVE-01 | security-integration | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| WAVE-01 | testing-delivery | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| WAVE-01 | sequencing-mvp | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| WAVE-01 | traceability-consistency | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |"""

QUESTION_TABLE = """| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |"""


# START_CONTRACT: source_inventory
#   PURPOSE: Return a stable markdown list of files in a source folder.
#   INPUTS: { label: str - source label, source: Path - source folder }
#   OUTPUTS: { str - markdown heading and bullet list with missing or empty markers }
#   SIDE_EFFECTS: None.
#   LINKS: M-BACKEND-WAVE-PLANNER / V-M-BACKEND-WAVE-PLANNER.
# END_CONTRACT: source_inventory
def source_inventory(label: str, source: Path) -> str:
    if not source.exists():
        return f"## {label}\n- SOURCE_MISSING: {source.as_posix()}"
    files = sorted(path for path in source.rglob("*") if path.is_file())
    if not files:
        return f"## {label}\n- SOURCE_EMPTY: {source.as_posix()}"
    items = "\n".join(f"- {path.as_posix()}" for path in files)
    return f"## {label}\n{items}"


# START_CONTRACT: build_files
#   PURPOSE: Build the required backend-waves markdown file set with canonical headings.
#   INPUTS: { technical_source: Path - approved technical docs, product_source: Path - verified product docs, wave_id: str - current wave number }
#   OUTPUTS: { dict[str, str] - relative path to markdown content }
#   SIDE_EFFECTS: None.
#   LINKS: M-BACKEND-WAVE-PLANNER / V-M-BACKEND-WAVE-PLANNER.
# END_CONTRACT: build_files
def build_files(technical_source: Path, product_source: Path, wave_id: str) -> dict[str, str]:
    wave_num = wave_id.removeprefix("wave-").upper()
    wave_code = f"WAVE-{wave_num}"
    technical_inventory = source_inventory("Technical Sources", technical_source)
    product_inventory = source_inventory("Product Sources", product_source)
    reviewer_rows = REVIEWER_ROWS.replace("WAVE-01", wave_code)
    return {
        "index.md": """# Backend Waves
## Status
draft
## Technical Approval Gate
PLACEHOLDER
## Current Wave Gate
PLACEHOLDER
## Source Set
See `source-inventory.md`.
## Next Action
Prepare the current wave and stop for user approval after it reaches ready-for-dev.
""",
        "source-inventory.md": f"""# Source Inventory
{technical_inventory}
{product_inventory}
## Prior Wave Sources
PLACEHOLDER
## Source Delta
PLACEHOLDER
## Coverage Gaps
PLACEHOLDER
""",
        "wave-map.md": """# Wave Map
## Backend Scope Inventory
PLACEHOLDER
## Tentative Wave Count
PLACEHOLDER
## Sequential Wave Map
PLACEHOLDER
## MVP Scope Check
PLACEHOLDER
## Dependency Notes
PLACEHOLDER
""",
        "open-questions.md": f"""# Open Questions
## Wave-Blocking
{QUESTION_TABLE}
## Needs Owner Decision
PLACEHOLDER
## Deferred
PLACEHOLDER
## Watchlist
PLACEHOLDER
## Resolved This Run
PLACEHOLDER
""",
        "waves/index.md": """# Waves
## Wave List
PLACEHOLDER
## Dependency Order
PLACEHOLDER
## Approval State
PLACEHOLDER
""",
        f"waves/{wave_id}.md": f"""# Wave {wave_num}: PLACEHOLDER
## Status
draft
## User Approval
Not approved.
## Outcome After Implementation
PLACEHOLDER
## Source Evidence
PLACEHOLDER
## Scope Included
PLACEHOLDER
## Scope Excluded
PLACEHOLDER
## Dependencies
PLACEHOLDER
## Backend Design
PLACEHOLDER
## Data And Migration Work
PLACEHOLDER
## API Jobs And Events
PLACEHOLDER
## Auth Security And Compliance
PLACEHOLDER
## Operations Observability
PLACEHOLDER
## Implementation Tasks
- BTASK-W{wave_num}-001 PLACEHOLDER
## Acceptance Criteria
- AC-W{wave_num}-001 PLACEHOLDER
## Exit Criteria
- EXIT-W{wave_num}-001 PLACEHOLDER
## Verification Plan
- BTEST-W{wave_num}-001 PLACEHOLDER
## Rollback And Compatibility
PLACEHOLDER
## Jira Ready Tasks
PLACEHOLDER
## Reviewer Verdicts
{reviewer_rows}
## Open Questions
{QUESTION_TABLE}
## Traceability
PLACEHOLDER
""",
        "appendix/reviewer-verdicts.md": f"""# Reviewer Verdicts
## Current Wave
{reviewer_rows}
## Historical Waves
PLACEHOLDER
## Rejected Findings
PLACEHOLDER
""",
        "appendix/traceability.md": """# Traceability
## Wave Task Map
PLACEHOLDER
## Acceptance Criteria Map
PLACEHOLDER
## Exit Criteria Map
PLACEHOLDER
## Test Obligation Map
PLACEHOLDER
## Question Map
PLACEHOLDER
## Source Map
PLACEHOLDER
""",
        "appendix/question-ledger.md": f"""# Question Ledger
## Open Questions
{QUESTION_TABLE}
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
## Technical Approval Gate
PLACEHOLDER
## User Wave Approvals
PLACEHOLDER
## Scope Decisions
PLACEHOLDER
## Deferrals
PLACEHOLDER
## Rejected Assumptions
PLACEHOLDER
""",
        "appendix/run-history.md": """# Run History
## Runs
PLACEHOLDER
## Wave Planning Cycles
PLACEHOLDER
## Source Delta History
PLACEHOLDER
## Approval Gate History
PLACEHOLDER
""",
    }


# START_CONTRACT: main
#   PURPOSE: Create a backend-waves scaffold at the requested output path.
#   INPUTS: { argv: CLI args - --technical-source, --product-source, --output, --wave-id, --force }
#   OUTPUTS: { int - process exit code }
#   SIDE_EFFECTS: Creates directories and markdown files under output.
#   LINKS: M-BACKEND-WAVE-PLANNER / V-M-BACKEND-WAVE-PLANNER.
# END_CONTRACT: main
def main() -> int:
    parser = argparse.ArgumentParser(description="Scaffold backend wave planning docs.")
    parser.add_argument("--technical-source", default="docs/technical-verified", help="Approved technical docs folder.")
    parser.add_argument("--product-source", default="docs/product-verified", help="Verified product docs folder.")
    parser.add_argument("--output", required=True, help="Output folder to create.")
    parser.add_argument("--wave-id", default="wave-01", help="Current wave file id such as wave-01.")
    parser.add_argument("--force", action="store_true", help="Overwrite existing managed files.")
    args = parser.parse_args()

    output = Path(args.output)
    output.mkdir(parents=True, exist_ok=True)
    for relative_dir in REQUIRED_DIRS:
        (output / relative_dir).mkdir(parents=True, exist_ok=True)

    files = build_files(Path(args.technical_source), Path(args.product_source), args.wave_id)
    for relative_path, content in files.items():
        target = output / relative_path
        target.parent.mkdir(parents=True, exist_ok=True)
        if target.exists() and not args.force:
            continue
        target.write_text(content, encoding="utf-8")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
