#!/usr/bin/env python3
# FILE: .agents/skills/detail-prd-wave/scripts/scaffold_detail_prd_wave.py
# VERSION: 1.0.0
# START_MODULE_CONTRACT
#   PURPOSE: Scaffold the managed prd-wave-details documentation tree for one selected backend wave.
#   SCOPE: Creates required markdown files and headings for selected-backend-wave detailed planning; excludes semantic PRD, codebase, frontend planning, and reviewer synthesis.
#   DEPENDS: Python standard library, .agents/skills/detail-prd-wave/references/output-contract.md.
#   LINKS: M-PRD-WAVE-DETAILER / V-M-PRD-WAVE-DETAILER.
#   ROLE: SCRIPT
#   MAP_MODE: EXPORTS
# END_MODULE_CONTRACT
#
# START_MODULE_MAP
#   REQUIRED_DIRS - Required subdirectories for the detailed wave package.
#   REVIEWER_PERSPECTIVES - Canonical reviewer perspectives for ready-for-dev.
#   QUESTION_TABLE - Canonical scaffold table for detailed wave questions.
#   normalize_wave_id - Convert user wave input to WAVE-NN.
#   wave_number - Return the two-digit numeric wave suffix.
#   source_inventory - List source files for source-inventory.md.
#   path_inventory - List configured GRACE source paths for source-inventory.md.
#   reviewer_rows - Return the canonical reviewer verdict table for a wave.
#   build_files - Return required file content keyed by relative path.
#   main - Parse arguments and write the scaffold.
# END_MODULE_MAP
# START_CHANGE_SUMMARY
#   LAST_CHANGE: 1.0.2 - Clarified scaffold output is backend-only and includes frontend-pages context only.
# END_CHANGE_SUMMARY

from __future__ import annotations

import argparse
import re
from pathlib import Path


REQUIRED_DIRS = ["waves", "appendix"]

REVIEWER_PERSPECTIVES = [
    "product-scope-and-ac",
    "architecture-codebase-fit",
    "data-api-integration-ops",
    "security-privacy-compliance",
    "testing-exit-criteria",
    "sequencing-other-wave-fit",
    "traceability-consistency",
    "final-wave-fit-review",
]

QUESTION_TABLE = """| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |"""


# START_CONTRACT: normalize_wave_id
#   PURPOSE: Convert common wave inputs into the canonical WAVE-NN id.
#   INPUTS: { raw: str - user-provided wave id }
#   OUTPUTS: { str - canonical WAVE-NN id }
#   SIDE_EFFECTS: None.
#   LINKS: M-PRD-WAVE-DETAILER / V-M-PRD-WAVE-DETAILER.
# END_CONTRACT: normalize_wave_id
def normalize_wave_id(raw: str) -> str:
    match = re.search(r"(\d+)", raw)
    if not match:
        raise ValueError(f"Cannot parse wave id from {raw!r}")
    return f"WAVE-{int(match.group(1)):02d}"


def wave_number(wave_id: str) -> str:
    return normalize_wave_id(wave_id).split("-")[1]


# START_CONTRACT: source_inventory
#   PURPOSE: Return a stable markdown list of files in a source folder.
#   INPUTS: { label: str - source label, source: Path - source folder }
#   OUTPUTS: { str - markdown heading and bullet list with missing or empty markers }
#   SIDE_EFFECTS: None.
#   LINKS: M-PRD-WAVE-DETAILER / V-M-PRD-WAVE-DETAILER.
# END_CONTRACT: source_inventory
def source_inventory(label: str, source: Path) -> str:
    if not source.exists():
        return f"## {label}\n- SOURCE_MISSING: {source.as_posix()}"
    if source.is_file():
        return f"## {label}\n- {source.as_posix()}"
    files = sorted(path for path in source.rglob("*") if path.is_file())
    if not files:
        return f"## {label}\n- SOURCE_EMPTY: {source.as_posix()}"
    items = "\n".join(f"- {path.as_posix()}" for path in files)
    return f"## {label}\n{items}"


def path_inventory(label: str, paths: list[Path]) -> str:
    items = []
    for path in paths:
        prefix = "" if path.exists() else "SOURCE_MISSING: "
        items.append(f"- {prefix}{path.as_posix()}")
    return f"## {label}\n" + "\n".join(items)


def reviewer_rows(wave_id: str) -> str:
    header = """| Wave | Perspective | Attempt | Verdict | Reviewer Report | Required Revisions | Notes |
| --- | --- | --- | --- | --- | --- | --- |"""
    rows = [
        f"| {wave_id} | {perspective} | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |"
        for perspective in REVIEWER_PERSPECTIVES
    ]
    return "\n".join([header, *rows])


# START_CONTRACT: build_files
#   PURPOSE: Build the required detailed backend PRD wave markdown file set with canonical headings.
#   INPUTS: { wave_id: str - canonical backend wave id, prd_waves: Path - shallow wave docs, product_source: Path - verified product docs, technical_source: Path - technical docs, grace_sources: list[Path] - GRACE docs }
#   OUTPUTS: { dict[str, str] - relative path to markdown content }
#   SIDE_EFFECTS: None.
#   LINKS: M-PRD-WAVE-DETAILER / V-M-PRD-WAVE-DETAILER.
# END_CONTRACT: build_files
def build_files(wave_id: str, prd_waves: Path, product_source: Path, technical_source: Path, grace_sources: list[Path]) -> dict[str, str]:
    wave_num = wave_number(wave_id)
    wave_slug = f"wave-{wave_num.lower()}"
    rows = reviewer_rows(wave_id)
    return {
        "index.md": f"""# Detailed Backend PRD Waves
## Status
draft
## Selected Wave
{wave_id}
## Source Wave Gate
PLACEHOLDER
## Current Wave Gate
PLACEHOLDER
## Source Set
See `source-inventory.md`.
## Next Action
Prepare only backend {wave_id} for ready-for-dev review.
""",
        "source-inventory.md": f"""# Source Inventory
{source_inventory("PRD Wave Sources", prd_waves)}
{source_inventory("Frontend Pages Source", prd_waves / "frontend-pages")}
{source_inventory("Product Sources", product_source)}
{source_inventory("Technical Sources", technical_source)}
{path_inventory("GRACE Sources", grace_sources)}
## Codebase Sources
PLACEHOLDER
## Source Delta
PLACEHOLDER
## Source Gaps
PLACEHOLDER
""",
        "wave-map-context.md": f"""# Wave Map Context
## Selected Backend Wave Boundary
{wave_id}: PLACEHOLDER
## Prior Backend Wave Fit
PLACEHOLDER
## Future Backend Wave Fit
PLACEHOLDER
## Frontend Pages Context
PLACEHOLDER
## Dependency Order
PLACEHOLDER
## Scope Collision Check
PLACEHOLDER
""",
        "codebase-fit.md": """# Codebase Fit
## Relevant Modules
PLACEHOLDER
## Relevant Files Read
PLACEHOLDER
## Public Contracts
PLACEHOLDER
## Generated Artifact Impact
PLACEHOLDER
## Integration Points
PLACEHOLDER
## Likely Graph Deltas
PLACEHOLDER
## Unsupported Assumptions
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
        "waves/index.md": f"""# Detailed Backend Waves
## Wave List
- {wave_id}: draft
## Dependency Order
PLACEHOLDER
## Approval State
PLACEHOLDER
""",
        f"waves/{wave_slug}.md": f"""# Wave {wave_num}: PLACEHOLDER
## Status
draft
## User Approval
Not approved.
## Source Wave Summary
PLACEHOLDER
## Outcome After Implementation
PLACEHOLDER
## Scope Included
PLACEHOLDER
## Scope Excluded
PLACEHOLDER
## Dependencies And Other-Wave Fit
PLACEHOLDER
## Frontend Pages Dependencies
PLACEHOLDER
## Codebase Fit And Touchpoints
PLACEHOLDER
## Design Contracts
PLACEHOLDER
## Data API Integration And Operations
PLACEHOLDER
## Security Privacy And Compliance
PLACEHOLDER
## Implementation Slices
PLACEHOLDER
## Acceptance Criteria
PLACEHOLDER
## Exit Criteria
PLACEHOLDER
## Verification Obligations
PLACEHOLDER
## Rollout Rollback And Compatibility
PLACEHOLDER
## Handoff Packets
PLACEHOLDER
## Reviewer Verdicts
{rows}
## Open Questions
{QUESTION_TABLE}
## Traceability
PLACEHOLDER
""",
        "appendix/reviewer-verdicts.md": f"""# Reviewer Verdicts
## Current Wave
{rows}
## Historical Waves
PLACEHOLDER
## Final Fit Reviews
PLACEHOLDER
## Rejected Findings
PLACEHOLDER
""",
        "appendix/traceability.md": """# Traceability
## Slice Map
PLACEHOLDER
## Acceptance Criteria Map
PLACEHOLDER
## Exit Criteria Map
PLACEHOLDER
## Verification Obligation Map
PLACEHOLDER
## Code Touchpoint Map
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
## Source Wave Gate
PLACEHOLDER
## User Wave Approvals
PLACEHOLDER
## Scope Decisions
PLACEHOLDER
## Codebase Fit Decisions
PLACEHOLDER
## Deferrals
PLACEHOLDER
## Rejected Assumptions
PLACEHOLDER
""",
        "appendix/run-history.md": f"""# Run History
## Runs
PLACEHOLDER
## Selected Wave History
{wave_id}: scaffolded.
## Planner Cycles
PLACEHOLDER
## Review Cycles
PLACEHOLDER
## Source Delta History
PLACEHOLDER
## Approval Gate History
PLACEHOLDER
""",
    }


# START_CONTRACT: main
#   PURPOSE: Parse CLI arguments and write the detailed backend PRD wave scaffold.
#   INPUTS: { argv: CLI - wave id, source folders, and output folder }
#   OUTPUTS: { process exit code - zero on success }
#   SIDE_EFFECTS: Creates directories and writes markdown scaffold files under output.
#   LINKS: M-PRD-WAVE-DETAILER / V-M-PRD-WAVE-DETAILER.
# END_CONTRACT: main
def main() -> None:
    parser = argparse.ArgumentParser(description="Scaffold detailed backend PRD wave docs for one selected backend wave.")
    parser.add_argument("--wave-id", required=True)
    parser.add_argument("--prd-waves", default="docs/prd-waves")
    parser.add_argument("--product-source", default="docs/product-verified")
    parser.add_argument("--technical-source", default="docs/technical-verified")
    parser.add_argument("--development-plan", default="docs/development-plan.xml")
    parser.add_argument("--knowledge-graph", default="docs/knowledge-graph.xml")
    parser.add_argument("--verification-plan", default="docs/verification-plan.xml")
    parser.add_argument("--output", required=True)
    args = parser.parse_args()

    wave_id = normalize_wave_id(args.wave_id)
    output = Path(args.output)
    for directory in REQUIRED_DIRS:
        (output / directory).mkdir(parents=True, exist_ok=True)

    grace_sources = [Path(args.development_plan), Path(args.knowledge_graph), Path(args.verification_plan)]
    files = build_files(wave_id, Path(args.prd_waves), Path(args.product_source), Path(args.technical_source), grace_sources)
    for relative, content in files.items():
        path = output / relative
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content, encoding="utf-8")


if __name__ == "__main__":
    main()
