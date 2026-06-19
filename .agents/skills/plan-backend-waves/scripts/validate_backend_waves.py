#!/usr/bin/env python3
# FILE: .agents/skills/plan-backend-waves/scripts/validate_backend_waves.py
# VERSION: 1.0.0
# START_MODULE_CONTRACT
#   PURPOSE: Validate the required backend-waves documentation structure and approval-gate invariants.
#   SCOPE: Checks files, headings, placeholder policy, status values, current wave readiness, reviewer approvals, and open blocking question consistency; excludes semantic truth review.
#   DEPENDS: Python standard library, .agents/skills/plan-backend-waves/references/output-contract.md.
#   LINKS: M-BACKEND-WAVE-PLANNER / V-M-BACKEND-WAVE-PLANNER.
#   ROLE: SCRIPT
#   MAP_MODE: EXPORTS
# END_MODULE_CONTRACT
#
# START_MODULE_MAP
#   REQUIRED_HEADINGS - Required markdown headings by relative package path.
#   WAVE_HEADINGS - Required markdown headings for each wave file.
#   VALID_PACKAGE_STATUSES - Package statuses accepted in index.md.
#   VALID_WAVE_STATUSES - Wave statuses accepted in wave files.
#   BLOCKING_SEVERITIES - Question severities that block ready-for-dev while open.
#   FINAL_WAVE_STATUSES - Wave statuses that require reviewer, criteria, and question gates.
#   REQUIRED_REVIEWERS - Required reviewer perspectives for every ready-for-dev wave.
#   VALID_REVIEWER_VERDICTS - Reviewer verdicts accepted in final package files.
#   DRAFT_REVIEWER_VERDICTS - Reviewer verdicts accepted only when placeholders are allowed.
#   main - Parse arguments and report validation failures.
#   validate_root - Validate required files, headings, placeholders, and approval invariants.
#   parse_question_rows - Parse markdown question ledger rows for gate checks.
#   status_from_heading - Extract the status that follows a Status heading.
#   extract_section - Extract a markdown section body under a heading.
#   wave_code_from_path - Convert a wave file path to its canonical WAVE-NN id.
#   wave_number_from_path - Convert a wave file path to its numeric sequence.
#   has_required_reviewer_approvals - Find required reviewer perspectives missing approved verdicts.
#   reviewer_verdict_errors - Validate reviewer verdict values.
#   contains_id - Check whether a wave file contains a stable id prefix.
#   collect_question_rows - Collect question rows from aggregate, package, and wave-local ledgers.
#   question_consistency_errors - Validate blocking question synchronization and row consistency.
# END_MODULE_MAP
# START_CHANGE_SUMMARY
#   LAST_CHANGE: 1.0.0 - Added validation script for backend wave planning output.
# END_CHANGE_SUMMARY

from __future__ import annotations

import argparse
import re
from pathlib import Path


REQUIRED_HEADINGS: dict[str, list[str]] = {
    "index.md": ["# Backend Waves", "## Status", "## Technical Approval Gate", "## Current Wave Gate", "## Source Set", "## Next Action"],
    "source-inventory.md": ["# Source Inventory", "## Technical Sources", "## Product Sources", "## Prior Wave Sources", "## Source Delta", "## Coverage Gaps"],
    "wave-map.md": ["# Wave Map", "## Backend Scope Inventory", "## Tentative Wave Count", "## Sequential Wave Map", "## MVP Scope Check", "## Dependency Notes"],
    "open-questions.md": ["# Open Questions", "## Wave-Blocking", "## Needs Owner Decision", "## Deferred", "## Watchlist", "## Resolved This Run"],
    "waves/index.md": ["# Waves", "## Wave List", "## Dependency Order", "## Approval State"],
    "appendix/reviewer-verdicts.md": ["# Reviewer Verdicts", "## Current Wave", "## Historical Waves", "## Rejected Findings"],
    "appendix/traceability.md": ["# Traceability", "## Wave Task Map", "## Acceptance Criteria Map", "## Exit Criteria Map", "## Test Obligation Map", "## Question Map", "## Source Map"],
    "appendix/question-ledger.md": ["# Question Ledger", "## Open Questions", "## Answered Questions", "## Follow-Up Questions", "## Resolved Questions", "## Deferred Questions"],
    "appendix/decision-log.md": ["# Decision Log", "## Technical Approval Gate", "## User Wave Approvals", "## Scope Decisions", "## Deferrals", "## Rejected Assumptions"],
    "appendix/run-history.md": ["# Run History", "## Runs", "## Wave Planning Cycles", "## Source Delta History", "## Approval Gate History"],
}

WAVE_HEADINGS = [
    "## Status",
    "## User Approval",
    "## Outcome After Implementation",
    "## Source Evidence",
    "## Scope Included",
    "## Scope Excluded",
    "## Dependencies",
    "## Backend Design",
    "## Data And Migration Work",
    "## API Jobs And Events",
    "## Auth Security And Compliance",
    "## Operations Observability",
    "## Implementation Tasks",
    "## Acceptance Criteria",
    "## Exit Criteria",
    "## Verification Plan",
    "## Rollback And Compatibility",
    "## Jira Ready Tasks",
    "## Reviewer Verdicts",
    "## Open Questions",
    "## Traceability",
]

VALID_PACKAGE_STATUSES = {
    "draft",
    "blocked",
    "wave-ready-awaiting-user-approval",
    "wave-approved-planning-next",
    "waves-verified",
    "superseded",
}

VALID_WAVE_STATUSES = {
    "draft",
    "needs-revision",
    "questions-open",
    "blocked",
    "ready-for-dev",
    "user-approved",
    "superseded",
}

BLOCKING_SEVERITIES = {"wave-blocking", "needs-owner-decision"}
FINAL_WAVE_STATUSES = {"ready-for-dev", "user-approved"}
REQUIRED_REVIEWERS = [
    "backend-architecture",
    "data-api-contract",
    "security-integration",
    "testing-delivery",
    "sequencing-mvp",
    "traceability-consistency",
]
VALID_REVIEWER_VERDICTS = {"approved", "needs-revision", "blocked"}
DRAFT_REVIEWER_VERDICTS = {"pending-review"}


# START_CONTRACT: parse_question_rows
#   PURPOSE: Extract question table rows from a backend wave question ledger.
#   INPUTS: { text: str - markdown ledger content }
#   OUTPUTS: { list[dict[str, str]] - parsed rows keyed by expected table columns }
#   SIDE_EFFECTS: None.
#   LINKS: M-BACKEND-WAVE-PLANNER / V-M-BACKEND-WAVE-PLANNER.
# END_CONTRACT: parse_question_rows
def parse_question_rows(text: str) -> list[dict[str, str]]:
    rows: list[dict[str, str]] = []
    headers = [
        "id",
        "wave",
        "scope",
        "severity",
        "parent",
        "question",
        "why",
        "needed",
        "source",
        "status",
        "resolution",
    ]
    for line in text.splitlines():
        if not line.startswith("|"):
            continue
        cells = [cell.strip() for cell in line.strip().strip("|").split("|")]
        if len(cells) != len(headers):
            continue
        if cells[0] in {"ID", "---"}:
            continue
        rows.append(dict(zip(headers, cells)))
    return rows


# START_CONTRACT: status_from_heading
#   PURPOSE: Extract the status that follows the first Status heading in a markdown file.
#   INPUTS: { text: str - markdown content }
#   OUTPUTS: { str | None - status value or none when the heading is absent }
#   SIDE_EFFECTS: None.
#   LINKS: M-BACKEND-WAVE-PLANNER / V-M-BACKEND-WAVE-PLANNER.
# END_CONTRACT: status_from_heading
def status_from_heading(text: str) -> str | None:
    match = re.search(r"^## Status\s*\n([^\n]+)", text, flags=re.MULTILINE)
    if not match:
        return None
    return match.group(1).strip()


def extract_section(text: str, heading: str) -> str:
    pattern = rf"^{re.escape(heading)}\s*\n(.*?)(?=^## |\Z)"
    match = re.search(pattern, text, flags=re.MULTILINE | re.DOTALL)
    if not match:
        return ""
    return match.group(1).strip()


def wave_code_from_path(path: Path) -> str:
    match = re.search(r"wave-(\d+)", path.stem)
    if not match:
        return path.stem.upper()
    return f"WAVE-{match.group(1)}"


def wave_number_from_path(path: Path) -> int | None:
    match = re.search(r"wave-(\d+)", path.stem)
    if not match:
        return None
    return int(match.group(1))


def has_required_reviewer_approvals(text: str) -> list[str]:
    missing: list[str] = []
    approved_reviewers: set[str] = set()
    for line in text.splitlines():
        if not line.startswith("|"):
            continue
        cells = [cell.strip().lower() for cell in line.strip().strip("|").split("|")]
        if len(cells) != 7 or cells[0] in {"wave", "---"}:
            continue
        perspective = cells[1]
        verdict = cells[3]
        if verdict == "approved":
            approved_reviewers.add(perspective)
    for reviewer in REQUIRED_REVIEWERS:
        if reviewer not in approved_reviewers:
            missing.append(reviewer)
    return missing


def reviewer_verdict_errors(text: str, relative: str, allow_placeholders: bool) -> list[str]:
    errors: list[str] = []
    for line in text.splitlines():
        if not line.startswith("|"):
            continue
        cells = [cell.strip().lower() for cell in line.strip().strip("|").split("|")]
        if len(cells) != 7 or cells[0] in {"wave", "---"}:
            continue
        verdict = cells[3]
        if verdict in VALID_REVIEWER_VERDICTS:
            continue
        if allow_placeholders and verdict in DRAFT_REVIEWER_VERDICTS:
            continue
        errors.append(f"{relative}: invalid reviewer verdict {verdict!r}")
    return errors


def contains_id(text: str, prefix: str) -> bool:
    return re.search(rf"\b{re.escape(prefix)}[A-Z0-9-]*\b", text) is not None


def collect_question_rows(root: Path, wave_files: list[Path]) -> list[dict[str, str]]:
    rows: list[dict[str, str]] = []
    question_sources = [
        root / "appendix" / "question-ledger.md",
        root / "open-questions.md",
    ] + wave_files
    for source in question_sources:
        if not source.exists():
            continue
        relative = source.relative_to(root).as_posix()
        for row in parse_question_rows(source.read_text(encoding="utf-8")):
            row["source_file"] = relative
            if relative.startswith("waves/") and not row["wave"]:
                row["wave"] = wave_code_from_path(source)
            rows.append(row)
    return rows


def question_consistency_errors(root: Path, rows: list[dict[str, str]]) -> list[str]:
    errors: list[str] = []
    by_id: dict[str, list[dict[str, str]]] = {}
    aggregate_ids = {row["id"] for row in rows if row.get("source_file") == "appendix/question-ledger.md"}
    open_questions_text = (root / "open-questions.md").read_text(encoding="utf-8") if (root / "open-questions.md").exists() else ""

    for row in rows:
        by_id.setdefault(row["id"], []).append(row)
        severity = row["severity"].lower()
        status = row["status"].lower()
        if status == "open" and severity in BLOCKING_SEVERITIES:
            if row["id"] not in open_questions_text:
                errors.append(f"open-questions.md: missing open blocking question {row['id']}")
            if row["id"] not in aggregate_ids:
                errors.append(f"appendix/question-ledger.md: missing open blocking question {row['id']} from {row['source_file']}")

    comparable_keys = ["wave", "scope", "severity", "parent", "question", "needed", "status", "resolution"]
    for question_id, question_rows in by_id.items():
        first = question_rows[0]
        for row in question_rows[1:]:
            differing = [
                key
                for key in comparable_keys
                if first.get(key, "").strip().lower() != row.get(key, "").strip().lower()
            ]
            if differing:
                sources = ", ".join(sorted({item["source_file"] for item in question_rows}))
                errors.append(f"question {question_id}: inconsistent ledger fields {', '.join(differing)} across {sources}")
                break
    return errors


# START_CONTRACT: validate_root
#   PURPOSE: Validate structure and approval-gate consistency for a backend-waves package.
#   INPUTS: { root: Path - package root, allow_placeholders: bool - whether PLACEHOLDER text is accepted }
#   OUTPUTS: { list[str] - validation error messages }
#   SIDE_EFFECTS: None.
#   LINKS: M-BACKEND-WAVE-PLANNER / V-M-BACKEND-WAVE-PLANNER.
# END_CONTRACT: validate_root
def validate_root(root: Path, allow_placeholders: bool) -> list[str]:
    errors: list[str] = []
    if not root.exists():
        return [f"missing output folder: {root}"]

    for relative_path, headings in REQUIRED_HEADINGS.items():
        path = root / relative_path
        if not path.exists():
            errors.append(f"missing required file: {relative_path}")
            continue
        text = path.read_text(encoding="utf-8")
        for heading in headings:
            if heading not in text:
                errors.append(f"{relative_path}: missing heading {heading!r}")
        if not allow_placeholders and "PLACEHOLDER" in text:
            errors.append(f"{relative_path}: contains PLACEHOLDER")
        if not allow_placeholders and re.search(r"\b(TODO|TBD)\b", text):
            errors.append(f"{relative_path}: contains TODO/TBD")

    index = root / "index.md"
    index_text = ""
    if index.exists():
        index_text = index.read_text(encoding="utf-8")
        package_status = status_from_heading(index_text)
        if package_status not in VALID_PACKAGE_STATUSES:
            errors.append(f"index.md: status must be one of {sorted(VALID_PACKAGE_STATUSES)}")
    else:
        package_status = None

    source_inventory = root / "source-inventory.md"
    if source_inventory.exists() and not allow_placeholders:
        technical_sources = extract_section(source_inventory.read_text(encoding="utf-8"), "## Technical Sources")
        if "SOURCE_MISSING" in technical_sources or "SOURCE_EMPTY" in technical_sources:
            errors.append("source-inventory.md: technical sources must exist and be non-empty for final validation")

    wave_files = sorted((root / "waves").glob("wave-*.md"))
    if not wave_files:
        errors.append("waves/: missing at least one wave-<nn>.md file")

    wave_statuses: dict[int, str] = {}
    for wave_file in wave_files:
        relative = wave_file.relative_to(root).as_posix()
        text = wave_file.read_text(encoding="utf-8")
        for heading in WAVE_HEADINGS:
            if heading not in text:
                errors.append(f"{relative}: missing heading {heading!r}")
        if not allow_placeholders and "PLACEHOLDER" in text:
            errors.append(f"{relative}: contains PLACEHOLDER")
        if not allow_placeholders and re.search(r"\b(TODO|TBD)\b", text):
            errors.append(f"{relative}: contains TODO/TBD")

        wave_status = status_from_heading(text)
        if wave_status not in VALID_WAVE_STATUSES:
            errors.append(f"{relative}: status must be one of {sorted(VALID_WAVE_STATUSES)}")
            continue
        wave_number = wave_number_from_path(wave_file)
        if wave_number is not None:
            wave_statuses[wave_number] = wave_status
        errors.extend(reviewer_verdict_errors(text, relative, allow_placeholders))

        wave_code = wave_code_from_path(wave_file)
        blockers = [
            row["id"]
            for row in collect_question_rows(root, wave_files)
            if row["wave"] in {wave_code, wave_file.stem, ""} and row["status"].lower() == "open" and row["severity"].lower() in BLOCKING_SEVERITIES
        ]
        if wave_status in FINAL_WAVE_STATUSES:
            if blockers:
                errors.append(f"{relative}: {wave_status} with open blocking questions: {', '.join(blockers)}")
            missing_reviewers = has_required_reviewer_approvals(text)
            if missing_reviewers:
                errors.append(f"{relative}: missing approved reviewer verdicts: {', '.join(missing_reviewers)}")
            wave_num = wave_code.split("-")[-1]
            if not contains_id(text, f"BTASK-W{wave_num}-"):
                errors.append(f"{relative}: missing BTASK-W{wave_num}- implementation task ids")
            if not contains_id(text, f"AC-W{wave_num}-"):
                errors.append(f"{relative}: missing AC-W{wave_num}- acceptance criteria ids")
            if not contains_id(text, f"EXIT-W{wave_num}-"):
                errors.append(f"{relative}: missing EXIT-W{wave_num}- exit criteria ids")
            if not contains_id(text, f"BTEST-W{wave_num}-"):
                errors.append(f"{relative}: missing BTEST-W{wave_num}- test obligation ids")
        if wave_status == "user-approved" and "approved-by-user" not in text:
            errors.append(f"{relative}: user-approved status requires approved-by-user in User Approval")

    rows = collect_question_rows(root, wave_files)
    errors.extend(question_consistency_errors(root, rows))
    open_blockers = [
        row["id"]
        for row in rows
        if row["status"].lower() == "open" and row["severity"].lower() in BLOCKING_SEVERITIES
    ]

    for wave_number in sorted(wave_statuses):
        for prior in range(1, wave_number):
            prior_status = wave_statuses.get(prior)
            if prior_status != "user-approved":
                errors.append(f"waves/wave-{wave_number:02d}.md: prior wave-{prior:02d} must be user-approved before this wave file exists")
                break

    requires_technical_gate = (
        package_status in {"wave-ready-awaiting-user-approval", "wave-approved-planning-next", "waves-verified"}
        or any(status in FINAL_WAVE_STATUSES for status in wave_statuses.values())
    )
    if requires_technical_gate and not allow_placeholders:
        technical_gate = extract_section(index_text, "## Technical Approval Gate")
        if "approved-to-dev" not in technical_gate.lower():
            errors.append("index.md: ready or approved waves require approved-to-dev evidence in Technical Approval Gate")

    if package_status in {"wave-ready-awaiting-user-approval", "wave-approved-planning-next", "waves-verified"} and open_blockers:
        errors.append(f"index.md: package status {package_status} with open blocking questions: {', '.join(sorted(set(open_blockers)))}")

    if package_status == "wave-ready-awaiting-user-approval":
        if not any(status_from_heading(path.read_text(encoding="utf-8")) == "ready-for-dev" for path in wave_files):
            errors.append("index.md: wave-ready-awaiting-user-approval requires a ready-for-dev wave")
    if package_status == "wave-approved-planning-next":
        if not any(status_from_heading(path.read_text(encoding="utf-8")) == "user-approved" for path in wave_files):
            errors.append("index.md: wave-approved-planning-next requires at least one user-approved wave")
    if package_status == "waves-verified":
        not_approved = [path.name for path in wave_files if status_from_heading(path.read_text(encoding="utf-8")) != "user-approved"]
        if not_approved:
            errors.append(f"index.md: waves-verified with non-user-approved waves: {', '.join(not_approved)}")

    return errors


# START_CONTRACT: main
#   PURPOSE: Validate a backend-waves package from the command line.
#   INPUTS: { argv: CLI args - root path and --allow-placeholders }
#   OUTPUTS: { int - process exit code }
#   SIDE_EFFECTS: Prints validation results to stdout.
#   LINKS: M-BACKEND-WAVE-PLANNER / V-M-BACKEND-WAVE-PLANNER.
# END_CONTRACT: main
def main() -> int:
    parser = argparse.ArgumentParser(description="Validate backend wave planning docs.")
    parser.add_argument("root", help="Backend waves docs folder.")
    parser.add_argument("--allow-placeholders", action="store_true", help="Allow scaffold placeholders.")
    args = parser.parse_args()

    errors = validate_root(Path(args.root), args.allow_placeholders)
    if errors:
        print("backend-waves validation failed:")
        for error in errors:
            print(f"- {error}")
        return 1
    print("backend-waves validation passed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
