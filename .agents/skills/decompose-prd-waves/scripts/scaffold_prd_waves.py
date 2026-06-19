#!/usr/bin/env python3
# FILE: .agents/skills/decompose-prd-waves/scripts/scaffold_prd_waves.py
# VERSION: 1.0.0
# START_MODULE_CONTRACT
#   PURPOSE: Scaffold the managed prd-waves documentation tree for staging or final output.
#   SCOPE: Creates required markdown files and headings for shallow backend wave decomposition plus per-page frontend files; excludes semantic PRD analysis and reviewer synthesis.
#   DEPENDS: Python standard library, .agents/skills/decompose-prd-waves/references/output-contract.md.
#   LINKS: M-PRD-WAVE-DECOMPOSER / V-M-PRD-WAVE-DECOMPOSER.
#   ROLE: SCRIPT
#   MAP_MODE: EXPORTS
# END_MODULE_CONTRACT
#
# START_MODULE_MAP
#   REQUIRED_DIRS - Required subdirectories for the PRD waves package.
#   REVIEWER_ROWS - Canonical scaffold rows for required reviewer perspectives.
#   QUESTION_TABLE - Canonical scaffold table for PRD wave questions.
#   main - Parse arguments and write the scaffold.
#   build_files - Return required file content keyed by relative path.
#   source_inventory - List source files for source-inventory.md.
# END_MODULE_MAP
# START_CHANGE_SUMMARY
#   LAST_CHANGE: 1.0.2 - Added raw PRD inventory and per-page frontend scaffold files.
# END_CHANGE_SUMMARY

from __future__ import annotations

import argparse
from pathlib import Path


REQUIRED_DIRS = [
    "waves",
    "frontend-pages",
    "appendix",
]

REVIEWER_ROWS = """| Scope | Perspective | Attempt | Verdict | Reviewer Report | Required Revisions | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| product-capabilities | scope-review | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| user-journeys | scope-review | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| data-lifecycle | scope-review | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| integrations-operations | scope-review | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| client-experience | scope-review | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| security-compliance | scope-review | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| delivery-sequencing | scope-review | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| wave-map-consistency | consistency-review | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| all | product-scope-coverage | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| all | technical-boundary-fit | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| all | sequencing-dependencies | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| all | backend-wave-boundary-quality | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
| all | traceability-consistency | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |"""

QUESTION_TABLE = """| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |"""


# START_CONTRACT: source_inventory
#   PURPOSE: Return a stable markdown list of files in a source folder.
#   INPUTS: { label: str - source label, source: Path - source folder }
#   OUTPUTS: { str - markdown heading and bullet list with missing or empty markers }
#   SIDE_EFFECTS: None.
#   LINKS: M-PRD-WAVE-DECOMPOSER / V-M-PRD-WAVE-DECOMPOSER.
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
#   PURPOSE: Build the required PRD waves markdown file set with canonical headings.
#   INPUTS: { raw_product_source: Path - raw product docs, product_source: Path - verified product docs, technical_source: Path - technical docs, wave_count: int - number of shallow wave files, frontend_page_count: int - number of frontend page files }
#   OUTPUTS: { dict[str, str] - relative path to markdown content }
#   SIDE_EFFECTS: None.
#   LINKS: M-PRD-WAVE-DECOMPOSER / V-M-PRD-WAVE-DECOMPOSER.
# END_CONTRACT: build_files
def build_files(raw_product_source: Path, product_source: Path, technical_source: Path, wave_count: int, frontend_page_count: int) -> dict[str, str]:
    raw_product_inventory = source_inventory("Raw Product Sources", raw_product_source)
    product_inventory = source_inventory("Verified Product Sources", product_source)
    technical_inventory = source_inventory("Technical Sources", technical_source)
    files = {
        "index.md": """# PRD Waves
## Status
draft
## Source Gate
PLACEHOLDER
## Shallow Wave Gate
PLACEHOLDER
## Wave Count
PLACEHOLDER
## Source Set
See `source-inventory.md`.
## Next Action
Prepare the backend-only shallow wave map and frontend page files, then stop before detailed planning.
""",
        "source-inventory.md": f"""# Source Inventory
{raw_product_inventory}
{product_inventory}
{technical_inventory}
## Prior Wave Sources
PLACEHOLDER
## Source Delta
PLACEHOLDER
## Source Gaps
PLACEHOLDER
""",
        "scope-inventory.md": """# Scope Inventory
## Capability Groups
PLACEHOLDER
## User Journey Groups
PLACEHOLDER
## Data Lifecycle Groups
PLACEHOLDER
## Integration And Operations Groups
PLACEHOLDER
## Client Experience Groups
PLACEHOLDER
## Security Compliance Groups
PLACEHOLDER
## Explicit Deferrals
PLACEHOLDER
""",
        "wave-map.md": """# Wave Map
## Top-Level Wave List
PLACEHOLDER
## Dependency Order
PLACEHOLDER
## Coverage Matrix
PLACEHOLDER
## More Than Eight Wave Check
PLACEHOLDER
## Downstream Planning Recommendations
PLACEHOLDER
""",
        "frontend-pages/index.md": """# Frontend Pages
## Status
draft
## Scope Source
PLACEHOLDER
## Page Order
PLACEHOLDER
## Raw PRD Source Coverage
PLACEHOLDER
## Verified PRD Source Coverage
PLACEHOLDER
## Shared UX States
PLACEHOLDER
## Backend Dependencies By Page
PLACEHOLDER
## Explicit Frontend Deferrals
PLACEHOLDER
## Open Questions
PLACEHOLDER
## Traceability
PLACEHOLDER
""",
        "open-questions.md": f"""# Open Questions
## Decomposition Blocking
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
        "appendix/reviewer-verdicts.md": f"""# Reviewer Verdicts
## Scope Reviews
{REVIEWER_ROWS}
## Consistency Review
PLACEHOLDER
## Rejected Findings
PLACEHOLDER
""",
        "appendix/traceability.md": """# Traceability
## Source To Scope Map
PLACEHOLDER
## Scope To Wave Map
PLACEHOLDER
## Wave To Source Map
PLACEHOLDER
## Question Map
PLACEHOLDER
## Decision Map
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
## Source Gate
PLACEHOLDER
## Scope Decisions
PLACEHOLDER
## Deferrals
PLACEHOLDER
## User Wave Map Approvals
PLACEHOLDER
## Rejected Assumptions
PLACEHOLDER
""",
        "appendix/run-history.md": """# Run History
## Runs
PLACEHOLDER
## Scope Mapper Cycles
PLACEHOLDER
## Consistency Cycles
PLACEHOLDER
## Source Delta History
PLACEHOLDER
## Approval Gate History
PLACEHOLDER
""",
    }
    for index in range(1, wave_count + 1):
        wave_num = f"{index:02d}"
        files[f"waves/wave-{wave_num}.md"] = f"""# Wave {wave_num}: PLACEHOLDER
## Status
draft
## User Approval
Not approved.
## Purpose
PLACEHOLDER
## Outcome After Wave
- OUT-W{wave_num}-001 PLACEHOLDER
## Included Scope
- CAP-W{wave_num}-001 PLACEHOLDER
## Excluded Scope
PLACEHOLDER
## Dependencies
PLACEHOLDER
## Surface Categories
backend, data, integrations, operations, security
## Risk Class
PLACEHOLDER
## Recommended Next Planning
- HANDOFF-W{wave_num}-001 PLACEHOLDER
## Open Questions
{QUESTION_TABLE}
## Traceability
PLACEHOLDER
"""
    for index in range(1, frontend_page_count + 1):
        page_num = f"{index:03d}"
        files[f"frontend-pages/page-{page_num}.md"] = f"""# PAGE-{page_num}: PLACEHOLDER
## Status
draft
## Page Purpose
PLACEHOLDER
## What Is On This Page
PLACEHOLDER
## Functional Parts
PLACEHOLDER
## Empty States
PLACEHOLDER
## Loading And Error States
PLACEHOLDER
## Backend Dependencies
PLACEHOLDER
## Explicit Deferrals
PLACEHOLDER
## Open Questions
PLACEHOLDER
## Raw PRD Traceability
PLACEHOLDER
## Verified PRD Traceability
PLACEHOLDER
"""
    return files


# START_CONTRACT: main
#   PURPOSE: Create a PRD waves scaffold at the requested output path.
#   INPUTS: { argv: CLI args - --raw-product-source, --product-source, --technical-source, --output, --wave-count, --frontend-page-count, --force }
#   OUTPUTS: { int - process exit code }
#   SIDE_EFFECTS: Creates directories and markdown files under output.
#   LINKS: M-PRD-WAVE-DECOMPOSER / V-M-PRD-WAVE-DECOMPOSER.
# END_CONTRACT: main
def main() -> int:
    parser = argparse.ArgumentParser(description="Scaffold shallow PRD wave decomposition docs.")
    parser.add_argument("--raw-product-source", default="docs/product", help="Raw product docs folder.")
    parser.add_argument("--product-source", default="docs/product-verified", help="Verified product docs folder.")
    parser.add_argument("--technical-source", default="docs/technical-verified", help="Technical docs folder.")
    parser.add_argument("--output", required=True, help="Output folder to create.")
    parser.add_argument("--wave-count", type=int, default=1, help="Number of shallow wave files to scaffold.")
    parser.add_argument("--frontend-page-count", type=int, default=1, help="Number of frontend page files to scaffold.")
    parser.add_argument("--force", action="store_true", help="Overwrite existing managed files.")
    args = parser.parse_args()

    if args.wave_count < 1:
        raise SystemExit("--wave-count must be at least 1")
    if args.frontend_page_count < 0:
        raise SystemExit("--frontend-page-count must be at least 0")

    output = Path(args.output)
    output.mkdir(parents=True, exist_ok=True)
    for relative_dir in REQUIRED_DIRS:
        (output / relative_dir).mkdir(parents=True, exist_ok=True)

    files = build_files(
        Path(args.raw_product_source),
        Path(args.product_source),
        Path(args.technical_source),
        args.wave_count,
        args.frontend_page_count,
    )
    for relative_path, content in files.items():
        target = output / relative_path
        if target.exists() and not args.force:
            continue
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(content, encoding="utf-8")

    print(f"Scaffolded PRD waves docs at {output.as_posix()}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
