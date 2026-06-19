#!/usr/bin/env python3
# FILE: .agents/skills/detail-prd-wave/scripts/validate_detail_prd_wave.py
# VERSION: 1.0.0
# START_MODULE_CONTRACT
#   PURPOSE: Validate the required detailed backend PRD wave documentation structure and ready-for-dev invariants.
#   SCOPE: Checks files, headings, selected-backend-wave status, placeholders, reviewer approvals, stable ids, source/codebase/other-backend-wave gates, frontend-boundary compliance, and open blocking question consistency; excludes semantic truth review.
#   DEPENDS: Python standard library, .agents/skills/detail-prd-wave/references/output-contract.md.
#   LINKS: M-PRD-WAVE-DETAILER / V-M-PRD-WAVE-DETAILER.
#   ROLE: SCRIPT
#   MAP_MODE: EXPORTS
# END_MODULE_CONTRACT
#
# START_MODULE_MAP
#   REQUIRED_HEADINGS - Required markdown headings by relative package path.
#   WAVE_HEADINGS - Required markdown headings for the selected wave file.
#   VALID_PACKAGE_STATUSES - Package statuses accepted in index.md.
#   VALID_WAVE_STATUSES - Wave statuses accepted in wave files.
#   FINAL_PACKAGE_STATUSES - Package statuses that require ready gates.
#   FINAL_WAVE_STATUSES - Wave statuses that require ready gates.
#   VALID_REVIEWER_VERDICTS - Reviewer verdicts accepted in final package files.
#   DRAFT_REVIEWER_VERDICTS - Reviewer verdicts accepted only when placeholders are allowed.
#   REQUIRED_REVIEWERS - Required reviewer perspectives for ready-for-dev.
#   BLOCKING_SEVERITIES - Question severities that block ready-for-dev while open.
#   BANNED_FRONTEND_PATTERNS - Frontend planning terms disallowed in backend wave details.
#   APPROVAL_MARKERS - Explicit user approval tokens accepted for wave-approved handoff.
#   normalize_wave_id - Convert user wave input to WAVE-NN.
#   wave_number - Return the two-digit numeric wave suffix.
#   resolve_source_path - Resolve source-inventory paths from cwd or package root.
#   source_inventory_paths - Extract path entries from a source inventory section.
#   extract_markdown_section - Extract a markdown heading section body.
#   status_from_heading - Extract the status that follows a Status heading.
#   has_approval_marker - Detect explicit user approval markers.
#   has_negative_approval_marker - Detect text that negates or defers approval.
#   parse_question_rows - Parse markdown question ledger rows for gate checks.
#   reviewer_verdict_errors - Validate reviewer verdict values.
#   missing_reviewer_approvals - Find missing approved reviewer perspectives.
#   required_id_errors - Validate required stable ids in ready wave sections.
#   source_wave_gate_errors - Validate source wave gate markers for final packages.
#   source_inventory_errors - Validate selected backend source wave and frontend-pages sources.
#   selected_wave_file_errors - Validate only the selected wave has a detailed wave file.
#   frontend_boundary_errors - Reject frontend planning details in backend wave files.
#   user_approval_errors - Validate explicit user approval evidence for approved waves.
#   all_question_rows - Collect question rows from package ledgers.
#   blocker_ledger_errors - Validate blocker ledger rows for blocked or questions-open packages.
#   question_errors - Validate blocking question synchronization.
#   validate_root - Validate package structure and invariants.
#   main - Parse arguments and report validation failures.
# END_MODULE_MAP
# START_CHANGE_SUMMARY
#   LAST_CHANGE: 1.0.2 - Added source-wave provenance and user-approval validation.
# END_CHANGE_SUMMARY

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path


REQUIRED_HEADINGS: dict[str, list[str]] = {
    "index.md": ["# Detailed Backend PRD Waves", "## Status", "## Selected Wave", "## Source Wave Gate", "## Current Wave Gate", "## Source Set", "## Next Action"],
    "source-inventory.md": ["# Source Inventory", "## PRD Wave Sources", "## Frontend Pages Source", "## Product Sources", "## Technical Sources", "## GRACE Sources", "## Codebase Sources", "## Source Delta", "## Source Gaps"],
    "wave-map-context.md": ["# Wave Map Context", "## Selected Backend Wave Boundary", "## Prior Backend Wave Fit", "## Future Backend Wave Fit", "## Frontend Pages Context", "## Dependency Order", "## Scope Collision Check"],
    "codebase-fit.md": ["# Codebase Fit", "## Relevant Modules", "## Relevant Files Read", "## Public Contracts", "## Generated Artifact Impact", "## Integration Points", "## Likely Graph Deltas", "## Unsupported Assumptions"],
    "open-questions.md": ["# Open Questions", "## Wave-Blocking", "## Needs Owner Decision", "## Deferred", "## Watchlist", "## Resolved This Run"],
    "waves/index.md": ["# Detailed Backend Waves", "## Wave List", "## Dependency Order", "## Approval State"],
    "appendix/reviewer-verdicts.md": ["# Reviewer Verdicts", "## Current Wave", "## Historical Waves", "## Final Fit Reviews", "## Rejected Findings"],
    "appendix/traceability.md": ["# Traceability", "## Slice Map", "## Acceptance Criteria Map", "## Exit Criteria Map", "## Verification Obligation Map", "## Code Touchpoint Map", "## Question Map", "## Source Map"],
    "appendix/question-ledger.md": ["# Question Ledger", "## Open Questions", "## Answered Questions", "## Follow-Up Questions", "## Resolved Questions", "## Deferred Questions"],
    "appendix/decision-log.md": ["# Decision Log", "## Source Wave Gate", "## User Wave Approvals", "## Scope Decisions", "## Codebase Fit Decisions", "## Deferrals", "## Rejected Assumptions"],
    "appendix/run-history.md": ["# Run History", "## Runs", "## Selected Wave History", "## Planner Cycles", "## Review Cycles", "## Source Delta History", "## Approval Gate History"],
}

WAVE_HEADINGS = [
    "## Status",
    "## User Approval",
    "## Source Wave Summary",
    "## Outcome After Implementation",
    "## Scope Included",
    "## Scope Excluded",
    "## Dependencies And Other-Wave Fit",
    "## Frontend Pages Dependencies",
    "## Codebase Fit And Touchpoints",
    "## Design Contracts",
    "## Data API Integration And Operations",
    "## Security Privacy And Compliance",
    "## Implementation Slices",
    "## Acceptance Criteria",
    "## Exit Criteria",
    "## Verification Obligations",
    "## Rollout Rollback And Compatibility",
    "## Handoff Packets",
    "## Reviewer Verdicts",
    "## Open Questions",
    "## Traceability",
]

VALID_PACKAGE_STATUSES = {"draft", "questions-open", "blocked", "wave-ready-awaiting-user-approval", "wave-approved", "superseded"}
VALID_WAVE_STATUSES = {"draft", "needs-revision", "questions-open", "blocked", "ready-for-dev", "user-approved", "superseded"}
FINAL_PACKAGE_STATUSES = {"wave-ready-awaiting-user-approval", "wave-approved"}
FINAL_WAVE_STATUSES = {"ready-for-dev", "user-approved"}
VALID_REVIEWER_VERDICTS = {"approved", "needs-revision", "blocked"}
DRAFT_REVIEWER_VERDICTS = {"pending-review"}
REQUIRED_REVIEWERS = [
    "product-scope-and-ac",
    "architecture-codebase-fit",
    "data-api-integration-ops",
    "security-privacy-compliance",
    "testing-exit-criteria",
    "sequencing-other-wave-fit",
    "traceability-consistency",
    "final-wave-fit-review",
]
BLOCKING_SEVERITIES = {"wave-blocking", "needs-owner-decision"}
BANNED_FRONTEND_PATTERNS = [
    r"\bfrontend implementation\b",
    r"\bfrontend task",
    r"\bcomponent architecture\b",
    r"\bvisual design\b",
    r"\bcopy deck\b",
    r"\broute implementation\b",
    r"\bnavigation implementation\b",
    r"\bUX state\b",
    r"\bUI state\b",
    r"\bscreen implementation\b",
    r"\bfrontend test",
]
APPROVAL_MARKERS = {"approved-by-user", "user-approved", "wave-approved-by-user"}


def normalize_wave_id(raw: str) -> str:
    match = re.search(r"(\d+)", raw)
    if not match:
        raise ValueError(f"Cannot parse wave id from {raw!r}")
    return f"WAVE-{int(match.group(1)):02d}"


def wave_number(wave_id: str) -> str:
    return normalize_wave_id(wave_id).split("-")[1]


def resolve_source_path(root: Path, raw: str) -> Path:
    path = Path(raw)
    if path.is_absolute() or path.exists():
        return path
    repo_candidate = root.parent.parent / path
    if repo_candidate.exists():
        return repo_candidate
    return path


def source_inventory_paths(section: str) -> list[str]:
    paths: list[str] = []
    for line in section.splitlines():
        stripped = line.strip()
        if not stripped.startswith("- "):
            continue
        value = stripped[2:].strip()
        if value.startswith(("SOURCE_MISSING:", "SOURCE_EMPTY:")):
            continue
        paths.append(value)
    return paths


def extract_markdown_section(text: str, heading: str) -> str:
    pattern = rf"^{re.escape(heading)}\s*\n(.*?)(?=^## |\Z)"
    match = re.search(pattern, text, flags=re.MULTILINE | re.DOTALL)
    return match.group(1).strip() if match else ""


def status_from_heading(text: str) -> str | None:
    match = re.search(r"^## Status\s*\n([^\n]+)", text, flags=re.MULTILINE)
    return match.group(1).strip() if match else None


def has_approval_marker(text: str) -> bool:
    lower = text.lower()
    return any(marker in lower for marker in APPROVAL_MARKERS)


def has_negative_approval_marker(text: str) -> bool:
    lower = text.lower()
    return bool(re.search(r"\bnot\s+(?:yet\s+)?(?:approved|user-approved)\b|\bawaiting approval\b|\bpending approval\b", lower))


# START_CONTRACT: parse_question_rows
#   PURPOSE: Extract question table rows from a detailed wave question ledger.
#   INPUTS: { text: str - markdown ledger content }
#   OUTPUTS: { list[dict[str, str]] - parsed rows keyed by expected table columns }
#   SIDE_EFFECTS: None.
#   LINKS: M-PRD-WAVE-DETAILER / V-M-PRD-WAVE-DETAILER.
# END_CONTRACT: parse_question_rows
def parse_question_rows(text: str) -> list[dict[str, str]]:
    rows: list[dict[str, str]] = []
    headers = ["id", "wave", "scope", "severity", "parent", "question", "why", "needed", "source", "status", "resolution"]
    for line in text.splitlines():
        if not line.startswith("|"):
            continue
        cells = [cell.strip() for cell in line.strip().strip("|").split("|")]
        if len(cells) != len(headers) or cells[0] in {"ID", "---"}:
            continue
        rows.append(dict(zip(headers, cells)))
    return rows


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


def missing_reviewer_approvals(text: str, wave_id: str) -> list[str]:
    approved: set[str] = set()
    for line in text.splitlines():
        if not line.startswith("|"):
            continue
        cells = [cell.strip().lower() for cell in line.strip().strip("|").split("|")]
        if len(cells) != 7 or cells[0] in {"wave", "---"}:
            continue
        wave, perspective, verdict = cells[0].upper(), cells[1], cells[3]
        if wave == wave_id and verdict == "approved":
            approved.add(perspective)
    return [perspective for perspective in REQUIRED_REVIEWERS if perspective not in approved]


def required_id_errors(wave_text: str, wave_num: str) -> list[str]:
    patterns = {
        "Implementation Slices": rf"\bSLICE-W{wave_num}-\d{{3}}\b",
        "Acceptance Criteria": rf"\bAC-W{wave_num}-\d{{3}}\b",
        "Exit Criteria": rf"\bEC-W{wave_num}-\d{{3}}\b",
        "Verification Obligations": rf"\bTEST-W{wave_num}-\d{{3}}\b",
    }
    errors: list[str] = []
    for label, pattern in patterns.items():
        section = extract_markdown_section(wave_text, f"## {label}")
        if not re.search(pattern, section):
            errors.append(f"waves/wave-{wave_num.lower()}.md: missing {label} id matching {pattern}")
    return errors


def source_wave_gate_errors(index_text: str, wave_id: str, requires_passed: bool) -> list[str]:
    section = extract_markdown_section(index_text, "## Source Wave Gate").lower()
    errors: list[str] = []
    wave_slug = f"waves/wave-{wave_number(wave_id).lower()}.md"
    if requires_passed:
        if "source-wave-gate: passed" not in section:
            errors.append("index.md: ready package requires `source-wave-gate: passed` in Source Wave Gate")
        if wave_id.lower() not in section:
            errors.append(f"index.md: ready package Source Wave Gate must name {wave_id}")
        if wave_slug not in section and f"wave-{wave_number(wave_id).lower()}.md" not in section:
            errors.append(f"index.md: ready package Source Wave Gate must name source path {wave_slug}")
    elif "source-wave-gate:" not in section and "blocked" not in section:
        errors.append("index.md: Source Wave Gate must record passed or blocked")
    return errors


def source_inventory_errors(root: Path, wave_id: str) -> list[str]:
    source_path = root / "source-inventory.md"
    if not source_path.exists():
        return ["source-inventory.md: missing source inventory"]

    text = source_path.read_text(encoding="utf-8")
    prd_section = extract_markdown_section(text, "## PRD Wave Sources")
    frontend_section = extract_markdown_section(text, "## Frontend Pages Source")
    errors: list[str] = []
    for label, section in [("PRD Wave Sources", prd_section), ("Frontend Pages Source", frontend_section)]:
        if "SOURCE_MISSING" in section or "SOURCE_EMPTY" in section:
            errors.append(f"source-inventory.md: {label} must exist and be non-empty for ready packages")

    wave_slug = f"waves/wave-{wave_number(wave_id).lower()}.md"
    selected_paths = [
        value for value in source_inventory_paths(prd_section)
        if value.endswith(wave_slug) or value.endswith(f"wave-{wave_number(wave_id).lower()}.md")
    ]
    if not selected_paths:
        errors.append(f"source-inventory.md: PRD Wave Sources must include selected source {wave_slug}")
    else:
        selected_path = resolve_source_path(root, selected_paths[0])
        if not selected_path.exists():
            errors.append(f"source-inventory.md: selected source wave path does not exist: {selected_paths[0]}")
        else:
            selected_text = selected_path.read_text(encoding="utf-8")
            selected_status = status_from_heading(selected_text)
            if selected_status not in {"top-level-ready", "user-approved"}:
                errors.append(f"source-inventory.md: selected source wave must be top-level-ready or user-approved, got {selected_status!r}")

    frontend_paths = source_inventory_paths(frontend_section)
    has_frontend_index = any(value.endswith("frontend-pages/index.md") for value in frontend_paths)
    if not has_frontend_index:
        errors.append("source-inventory.md: Frontend Pages Source must include frontend-pages/index.md")
    else:
        index_value = next(value for value in frontend_paths if value.endswith("frontend-pages/index.md"))
        if not resolve_source_path(root, index_value).exists():
            errors.append(f"source-inventory.md: frontend pages index path does not exist: {index_value}")
    return errors


def selected_wave_file_errors(root: Path, selected_relative: str) -> list[str]:
    waves_dir = root / "waves"
    if not waves_dir.exists():
        return []
    errors: list[str] = []
    for path in sorted(waves_dir.glob("wave-*.md")):
        relative = path.relative_to(root).as_posix()
        if relative != selected_relative:
            errors.append(f"{relative}: extra detailed wave file is not allowed in a selected backend-wave package")
    return errors


def frontend_boundary_errors(wave_text: str, wave_num: str) -> list[str]:
    errors: list[str] = []
    for pattern in BANNED_FRONTEND_PATTERNS:
        if re.search(pattern, wave_text, flags=re.IGNORECASE):
            errors.append(f"waves/wave-{wave_num.lower()}.md: contains frontend planning pattern {pattern!r}")
    return errors


def user_approval_errors(root: Path, wave_text: str, wave_relative: str) -> list[str]:
    errors: list[str] = []
    approval_section = extract_markdown_section(wave_text, "## User Approval")
    if has_negative_approval_marker(approval_section) or not has_approval_marker(approval_section):
        errors.append(f"{wave_relative}: user-approved wave requires approved-by-user evidence in User Approval")
    decision_path = root / "appendix" / "decision-log.md"
    decision_text = decision_path.read_text(encoding="utf-8") if decision_path.exists() else ""
    if has_negative_approval_marker(decision_text) or not has_approval_marker(decision_text):
        errors.append("appendix/decision-log.md: wave-approved requires explicit approved-by-user decision")
    return errors


def all_question_rows(root: Path, wave_id: str) -> list[dict[str, str]]:
    rows: list[dict[str, str]] = []
    for relative in ["open-questions.md", "appendix/question-ledger.md", f"waves/wave-{wave_number(wave_id).lower()}.md"]:
        path = root / relative
        if path.exists():
            rows.extend(parse_question_rows(path.read_text(encoding="utf-8")))
    return rows


def blocker_ledger_errors(root: Path, wave_id: str) -> list[str]:
    rows = all_question_rows(root, wave_id)
    for row in rows:
        if row["wave"].upper() in {wave_id, "ALL"} and row["severity"] in BLOCKING_SEVERITIES and row["status"] == "open":
            return []
    return [f"{wave_id}: blocked or questions-open package requires an open wave-blocking or needs-owner-decision question row"]


def question_errors(root: Path, wave_id: str) -> list[str]:
    errors: list[str] = []
    rows = all_question_rows(root, wave_id)
    for row in rows:
        if row["wave"].upper() not in {wave_id, "ALL"}:
            continue
        if row["severity"] in BLOCKING_SEVERITIES and row["status"] == "open":
            errors.append(f"open blocking question {row['id']} for {wave_id}")
    return errors


# START_CONTRACT: validate_root
#   PURPOSE: Validate detailed backend PRD wave package structure and ready-for-dev gates.
#   INPUTS: { root: Path - package root, wave_id: str - selected wave id, allow_placeholders: bool - draft mode flag }
#   OUTPUTS: { list[str] - validation error messages }
#   SIDE_EFFECTS: Reads markdown files.
#   LINKS: M-PRD-WAVE-DETAILER / V-M-PRD-WAVE-DETAILER.
# END_CONTRACT: validate_root
def validate_root(root: Path, wave_id: str, allow_placeholders: bool) -> list[str]:
    errors: list[str] = []
    if not root.exists():
        return [f"missing output root: {root}"]

    for relative, headings in REQUIRED_HEADINGS.items():
        path = root / relative
        if not path.exists():
            errors.append(f"missing required file: {relative}")
            continue
        text = path.read_text(encoding="utf-8")
        for heading in headings:
            if heading not in text:
                errors.append(f"{relative}: missing heading {heading!r}")
        if not allow_placeholders and "PLACEHOLDER" in text:
            errors.append(f"{relative}: contains PLACEHOLDER")
        if relative == "appendix/reviewer-verdicts.md":
            errors.extend(reviewer_verdict_errors(text, relative, allow_placeholders))

    index_text = (root / "index.md").read_text(encoding="utf-8") if (root / "index.md").exists() else ""
    package_status = status_from_heading(index_text) if index_text else None
    if package_status not in VALID_PACKAGE_STATUSES:
        errors.append(f"index.md: invalid status {package_status!r}")

    selected = extract_markdown_section((root / "index.md").read_text(encoding="utf-8"), "## Selected Wave") if (root / "index.md").exists() else ""
    if wave_id not in selected:
        errors.append(f"index.md: selected wave does not include {wave_id}")

    wave_num = wave_number(wave_id)
    wave_relative = f"waves/wave-{wave_num.lower()}.md"
    errors.extend(selected_wave_file_errors(root, wave_relative))
    wave_path = root / wave_relative
    if not wave_path.exists():
        errors.append(f"missing selected wave file: {wave_relative}")
        return errors

    wave_text = wave_path.read_text(encoding="utf-8")
    for heading in WAVE_HEADINGS:
        if heading not in wave_text:
            errors.append(f"{wave_relative}: missing heading {heading!r}")
    if not allow_placeholders and "PLACEHOLDER" in wave_text:
        errors.append(f"{wave_relative}: contains PLACEHOLDER")
    errors.extend(reviewer_verdict_errors(wave_text, wave_relative, allow_placeholders))

    wave_status = status_from_heading(wave_text)
    if wave_status not in VALID_WAVE_STATUSES:
        errors.append(f"{wave_relative}: invalid status {wave_status!r}")

    if package_status in FINAL_PACKAGE_STATUSES or wave_status in FINAL_WAVE_STATUSES:
        errors.extend(source_wave_gate_errors(index_text, wave_id, requires_passed=True))
        errors.extend(source_inventory_errors(root, wave_id))
        errors.extend(required_id_errors(wave_text, wave_num))
        errors.extend(frontend_boundary_errors(wave_text, wave_num))
        reviewer_text = (root / "appendix/reviewer-verdicts.md").read_text(encoding="utf-8")
        missing = missing_reviewer_approvals(reviewer_text + "\n" + wave_text, wave_id)
        if missing:
            errors.append(f"{wave_relative}: missing approved reviewer perspectives: {', '.join(missing)}")
        errors.extend(question_errors(root, wave_id))
        codebase_fit = (root / "codebase-fit.md").read_text(encoding="utf-8")
        if "## Relevant Modules\nPLACEHOLDER" in codebase_fit or "## Relevant Files Read\nPLACEHOLDER" in codebase_fit:
            errors.append("codebase-fit.md: ready wave needs concrete module/file evidence or explicit none-needed rationale")
        if package_status == "wave-approved" or wave_status == "user-approved":
            errors.extend(user_approval_errors(root, wave_text, wave_relative))
    elif package_status in {"blocked", "questions-open"} or wave_status in {"blocked", "questions-open"}:
        errors.extend(source_wave_gate_errors(index_text, wave_id, requires_passed=False))
        errors.extend(blocker_ledger_errors(root, wave_id))

    return errors


def main() -> None:
    parser = argparse.ArgumentParser(description="Validate detailed backend PRD wave docs.")
    parser.add_argument("root")
    parser.add_argument("--wave-id", required=True)
    parser.add_argument("--allow-placeholders", action="store_true")
    args = parser.parse_args()

    try:
        wave_id = normalize_wave_id(args.wave_id)
    except ValueError as exc:
        print(str(exc), file=sys.stderr)
        sys.exit(2)

    errors = validate_root(Path(args.root), wave_id, args.allow_placeholders)
    if errors:
        for error in errors:
            print(f"ERROR: {error}", file=sys.stderr)
        sys.exit(1)
    print(f"OK: detailed backend PRD wave package valid for {wave_id}: {args.root}")


if __name__ == "__main__":
    main()
