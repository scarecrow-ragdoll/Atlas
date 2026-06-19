#!/usr/bin/env python3
# FILE: .agents/skills/verify-technical-docs/scripts/validate_technical_verified.py
# VERSION: 1.0.0
# START_MODULE_CONTRACT
#   PURPOSE: Validate the required technical-verified documentation structure and approval-gate invariants.
#   SCOPE: Checks files, headings, placeholder policy, status values, and open blocking question consistency; excludes semantic truth review.
#   DEPENDS: Python standard library, .agents/skills/verify-technical-docs/references/output-contract.md.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
#   ROLE: SCRIPT
#   MAP_MODE: EXPORTS
# END_MODULE_CONTRACT
#
# START_MODULE_MAP
#   REQUIRED_HEADINGS - Required markdown headings by relative package path.
#   REQUIRED_SCOPES - Technical scopes that require reviewer approval for approved-to-dev.
#   REQUIRED_APPROVAL_GATES - Approval gate rows required for approved-to-dev.
#   DELTA_GATE_KEYWORDS - Evidence keywords required for delta-review gates.
#   SOURCE_INVENTORY_HEADING - Heading that owns included source evidence.
#   ALLOWED_SEVERITIES - Question severities accepted by the technical ledger.
#   ALLOWED_STATUSES - Question statuses accepted by the technical ledger.
#   VALID_STATUSES - Status values accepted in index.md.
#   BLOCKING_SEVERITIES - Question severities that block approved-to-dev while open.
#   MISSING_SOURCE_MARKERS - Source inventory markers that cannot be approved-to-dev.
#   main - Parse arguments and report validation failures.
#   validate_root - Validate required files, headings, placeholders, and approval status.
#   parse_question_rows - Parse markdown question ledger rows for gate checks.
#   parse_table_rows - Parse simple markdown tables by expected header.
#   extract_section_lines - Extract markdown lines under a specific heading.
#   included_source_entries - Extract concrete included source bullets from source-inventory.md.
#   source_entry_is_concrete - Check whether an included source bullet looks like a source file path.
#   source_entry_exists - Check whether an included source path exists from a validation root.
#   duplicate_values - Find duplicated values in parsed table rows.
#   evidence_matches_keywords - Check whether approval evidence names required source classes.
#   report_path_valid - Check whether scope reviewer evidence points to the expected attempt artifact.
#   status_from_index - Extract the package status from index.md.
# END_MODULE_MAP
# START_CHANGE_SUMMARY
#   LAST_CHANGE: 1.0.0 - Added validation script for technical readiness output.
# END_CHANGE_SUMMARY

from __future__ import annotations

import argparse
import re
from pathlib import Path, PurePosixPath


REQUIRED_HEADINGS: dict[str, list[str]] = {
    "index.md": ["# Technical Verified", "## Status", "## Source Set", "## Document Map", "## Dev Handoff Readiness"],
    "source-inventory.md": ["# Source Inventory", "## Included Sources", "## Source Delta", "## Answered Questions", "## Excluded Or Noisy Sources", "## Coverage Gaps"],
    "technical-brief.md": ["# Technical Brief", "## Product Signal", "## Technical Scope", "## Constraints", "## Assumptions", "## Readiness Summary"],
    "architecture-and-boundaries.md": ["# Architecture And Boundaries", "## System Context", "## Components", "## Ownership Boundaries", "## Dependencies", "## Architecture Questions"],
    "data-contracts.md": ["# Data Contracts", "## Entities And Identifiers", "## Persistence And Storage", "## Migrations", "## Retention And Privacy", "## Data Questions"],
    "api-contracts.md": ["# API Contracts", "## Surfaces", "## Requests And Responses", "## Error And Validation Contracts", "## Compatibility And Idempotency", "## API Questions"],
    "auth-security-compliance.md": ["# Auth Security Compliance", "## Identity", "## Authorization And Ownership", "## Auditability", "## Privacy And Compliance", "## Security Questions"],
    "integrations-and-events.md": ["# Integrations And Events", "## External Systems", "## Events And Jobs", "## Sync And Retry Rules", "## Rate Limits And Failure Handling", "## Integration Questions"],
    "client-state-and-ux-contracts.md": ["# Client State And UX Contracts", "## User Interface States", "## Form And Validation Behavior", "## Cache And Realtime Behavior", "## Accessibility And Localization", "## Client Questions"],
    "operations-observability.md": ["# Operations Observability", "## Environments And Config", "## Deployment And Rollout", "## Logs Metrics Traces", "## Alerts Runbooks Backups", "## Operations Questions"],
    "testing-and-delivery.md": ["# Testing And Delivery", "## Test Strategy", "## Fixtures And Test Data", "## Contract And E2E Coverage", "## Release Gates", "## Testing Questions"],
    "implementation-slices.md": ["# Implementation Slices", "## Slice Map", "## Dependencies", "## Blockers", "## Verification"],
    "open-questions.md": ["# Open Questions", "## Dev-Blocking", "## Needs Owner Decision", "## Deferred", "## Watchlist", "## Resolved This Run"],
    "features/index.md": ["# Feature Technical Inventory", "## Features"],
    "appendix/subagent-findings.md": ["# Subagent Findings", "## Scope Reports", "## Reviewer Verdicts", "## Conflicts"],
    "appendix/traceability.md": ["# Traceability", "## Technical Requirement Map", "## Question Map", "## Decision Map", "## Slice Map", "## Source Map"],
    "appendix/question-ledger.md": ["# Question Ledger", "## Open Questions", "## Answered Questions", "## Follow-Up Questions", "## Resolved Questions", "## Deferred Questions"],
    "appendix/decision-log.md": ["# Decision Log", "## Technical Decisions", "## Deferrals", "## Superseded Answers", "## Rejected Assumptions"],
    "appendix/loop-history.md": ["# Loop History", "## Runs", "## Answered Question Effects", "## Follow-Up Blockers", "## Approval Gate History"],
}

VALID_STATUSES = {"questions-open", "blocked", "approved-to-dev"}
ALLOWED_SEVERITIES = {"dev-blocking", "needs-owner-decision", "deferred", "watchlist"}
ALLOWED_STATUSES = {"open", "answered-by-source", "answered-by-user", "resolved-by-decision", "superseded", "deferred"}
BLOCKING_SEVERITIES = {"dev-blocking", "needs-owner-decision"}
MISSING_SOURCE_MARKERS = {"SOURCE_MISSING", "SOURCE_EMPTY"}
SOURCE_INVENTORY_HEADING = "## Included Sources"
REQUIRED_SCOPES = {
    "architecture-boundaries",
    "data-contracts",
    "api-contracts",
    "auth-security-compliance",
    "integrations-events",
    "client-state-ux",
    "operations-observability",
    "testing-delivery",
    "consistency-loop-reviewer",
}
REQUIRED_APPROVAL_GATES = {
    "required-scopes-approved",
    "consistency-approved",
    "source-deltas-reviewed",
    "answer-deltas-reviewed",
    "no-answer-spawned-blockers",
    "no-open-blocking-questions",
}
DELTA_GATE_KEYWORDS = {
    "source-deltas-reviewed": ("source-delta", "source inventory", "product verified changes", "changed sources"),
    "answer-deltas-reviewed": ("answered question effects", "question id", "tq-", "source-delta"),
}


# START_CONTRACT: parse_table_rows
#   PURPOSE: Extract rows from markdown tables that match a known header exactly.
#   INPUTS: { text: str - markdown content, headers: list[str] - expected header names }
#   OUTPUTS: { list[dict[str, str]] - parsed table rows keyed by normalized header }
#   SIDE_EFFECTS: None.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: parse_table_rows
def parse_table_rows(text: str, headers: list[str]) -> list[dict[str, str]]:
    rows: list[dict[str, str]] = []
    normalized_headers = [header.lower().replace(" ", "_") for header in headers]
    lines = text.splitlines()
    for index, line in enumerate(lines):
        if not line.startswith("|"):
            continue
        cells = [cell.strip() for cell in line.strip().strip("|").split("|")]
        if cells != headers:
            continue
        if index + 1 >= len(lines):
            continue
        separator = [cell.strip() for cell in lines[index + 1].strip().strip("|").split("|")]
        if len(separator) != len(headers) or any(set(cell) - {"-", ":"} or "-" not in cell for cell in separator):
            continue
        for row_line in lines[index + 2 :]:
            if not row_line.startswith("|"):
                break
            row_cells = [cell.strip() for cell in row_line.strip().strip("|").split("|")]
            if len(row_cells) != len(headers):
                break
            rows.append(dict(zip(normalized_headers, row_cells)))
    return rows


# START_CONTRACT: duplicate_values
#   PURPOSE: Find duplicated values in parsed markdown table rows.
#   INPUTS: { rows: list[dict[str, str]] - parsed rows, key: str - column key to inspect }
#   OUTPUTS: { set[str] - duplicated values }
#   SIDE_EFFECTS: None.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: duplicate_values
def duplicate_values(rows: list[dict[str, str]], key: str) -> set[str]:
    seen: set[str] = set()
    duplicates: set[str] = set()
    for row in rows:
        value = row[key]
        if value in seen:
            duplicates.add(value)
        seen.add(value)
    return duplicates


# START_CONTRACT: evidence_matches_keywords
#   PURPOSE: Check whether approval gate evidence names the required source class instead of shallow pass text.
#   INPUTS: { evidence: str - approval gate evidence, keywords: tuple[str, ...] - accepted case-insensitive markers }
#   OUTPUTS: { bool - true when at least one required marker is present }
#   SIDE_EFFECTS: None.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: evidence_matches_keywords
def evidence_matches_keywords(evidence: str, keywords: tuple[str, ...]) -> bool:
    normalized = evidence.lower()
    return any(keyword in normalized for keyword in keywords)


# START_CONTRACT: extract_section_lines
#   PURPOSE: Extract lines below a markdown heading until the next same-or-higher-level heading.
#   INPUTS: { text: str - markdown content, heading: str - heading line to extract }
#   OUTPUTS: { list[str] - section lines without the heading }
#   SIDE_EFFECTS: None.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: extract_section_lines
def extract_section_lines(text: str, heading: str) -> list[str]:
    lines = text.splitlines()
    try:
        start = lines.index(heading) + 1
    except ValueError:
        return []
    heading_level = len(heading) - len(heading.lstrip("#"))
    section: list[str] = []
    for line in lines[start:]:
        stripped = line.lstrip()
        if stripped.startswith("#"):
            level = len(stripped) - len(stripped.lstrip("#"))
            if level <= heading_level:
                break
        section.append(line)
    return section


# START_CONTRACT: included_source_entries
#   PURPOSE: Extract concrete existing source path bullets from the Included Sources section.
#   INPUTS: { text: str - source-inventory.md content, package_root: Path - technical package root }
#   OUTPUTS: { list[str] - non-placeholder source entries }
#   SIDE_EFFECTS: None.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: included_source_entries
def included_source_entries(text: str, package_root: Path) -> list[str]:
    entries: list[str] = []
    for line in extract_section_lines(text, SOURCE_INVENTORY_HEADING):
        stripped = line.strip()
        if not stripped.startswith("- "):
            continue
        value = stripped[2:].strip()
        if not value:
            continue
        if not source_entry_is_concrete(value):
            continue
        if not source_entry_exists(value, package_root):
            continue
        entries.append(value)
    return entries


# START_CONTRACT: source_entry_is_concrete
#   PURPOSE: Check whether an Included Sources bullet is a concrete local source file path.
#   INPUTS: { value: str - source inventory bullet text }
#   OUTPUTS: { bool - true when the entry resembles a concrete file path, not prose or a placeholder }
#   SIDE_EFFECTS: None.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: source_entry_is_concrete
def source_entry_is_concrete(value: str) -> bool:
    normalized = value.lower()
    if value in MISSING_SOURCE_MARKERS:
        return False
    if any(marker in normalized for marker in ("<", ">")):
        return False
    if "/" not in value:
        return False
    path = PurePosixPath(value)
    if not path.name or "." not in path.name:
        return False
    parts = [part for part in path.parts if part != "/"]
    placeholder_parts = {"placeholder", "tbd", "todo", "none", "none."}
    if any(part.lower() in placeholder_parts for part in parts):
        return False
    return bool(parts) and all(part not in {"", ".", ".."} for part in parts)


# START_CONTRACT: source_entry_exists
#   PURPOSE: Check whether an Included Sources path exists as an absolute path or relative to likely validation roots.
#   INPUTS: { value: str - source inventory path, package_root: Path - technical package root }
#   OUTPUTS: { bool - true when the path resolves to an existing file }
#   SIDE_EFFECTS: None.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: source_entry_exists
def source_entry_exists(value: str, package_root: Path) -> bool:
    path = Path(value)
    if path.is_absolute():
        return path.is_file()
    for candidate_root in [Path.cwd(), *package_root.parents]:
        if (candidate_root / path).is_file():
            return True
    return False


# START_CONTRACT: report_path_valid
#   PURPOSE: Check whether a scope reviewer report value names the expected orchestration artifact path.
#   INPUTS: { scope: str - technical review scope, report: str - reviewer report path, package_root: Path - technical package root }
#   OUTPUTS: { bool - true when the path is a scope-local review-attempt markdown artifact }
#   SIDE_EFFECTS: None.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: report_path_valid
def report_path_valid(scope: str, report: str, package_root: Path) -> bool:
    normalized = report.lower()
    if any(marker in normalized for marker in ("<", ">")):
        return False
    placeholder_parts = {"placeholder", "tbd", "todo", "none"}
    if any(part.lower() in placeholder_parts for part in PurePosixPath(report).parts):
        return False
    pattern = rf"^\.tasks/technical-docs-verify/[^/]+/scopes/{re.escape(scope)}/review-attempt-\d+\.md$"
    if re.match(pattern, report) is None:
        return False
    for candidate_root in [Path.cwd(), *package_root.parents]:
        if (candidate_root / report).exists():
            return True
    return False


# START_CONTRACT: parse_question_rows
#   PURPOSE: Extract question table rows from a technical question ledger.
#   INPUTS: { text: str - markdown ledger content }
#   OUTPUTS: { list[dict[str, str]] - parsed rows keyed by expected table columns }
#   SIDE_EFFECTS: None.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: parse_question_rows
def parse_question_rows(text: str) -> list[dict[str, str]]:
    headers = ["ID", "Scope", "Severity", "Parent", "Question", "Why It Matters", "Needed Artifact Or Decision", "Source Or Report", "Status", "Resolution"]
    rows = parse_table_rows(text, headers)
    seen = {(row["id"], row["question"]) for row in rows}
    normalized_headers = [header.lower().replace(" ", "_") for header in headers]
    for line in text.splitlines():
        if not line.startswith("|"):
            continue
        cells = [cell.strip() for cell in line.strip().strip("|").split("|")]
        if len(cells) != len(headers):
            continue
        if cells[0] in {"ID", "---"} or set(cells[0]) == {"-"}:
            continue
        key = (cells[0], cells[4])
        if key in seen:
            continue
        rows.append(dict(zip(normalized_headers, cells)))
        seen.add(key)
    return rows


# START_CONTRACT: status_from_index
#   PURPOSE: Extract the package status that follows the Status heading in index.md.
#   INPUTS: { text: str - index.md content }
#   OUTPUTS: { str | None - status value or none when the heading is absent }
#   SIDE_EFFECTS: None.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: status_from_index
def status_from_index(text: str) -> str | None:
    match = re.search(r"^## Status\s*\n([^\n]+)", text, flags=re.MULTILINE)
    if not match:
        return None
    return match.group(1).strip()


# START_CONTRACT: validate_root
#   PURPOSE: Validate structure and approval-gate consistency for a technical-verified package.
#   INPUTS: { root: Path - package root, allow_placeholders: bool - whether PLACEHOLDER text is accepted }
#   OUTPUTS: { list[str] - validation error messages }
#   SIDE_EFFECTS: None.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
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
    ledger = root / "appendix" / "question-ledger.md"
    inventory = root / "source-inventory.md"
    if index.exists():
        status = status_from_index(index.read_text(encoding="utf-8"))
        if status not in VALID_STATUSES:
            errors.append(f"index.md: status must be one of {sorted(VALID_STATUSES)}")
    else:
        status = None

    rows: list[dict[str, str]] = []
    if ledger.exists():
        rows = parse_question_rows(ledger.read_text(encoding="utf-8"))
        for row in rows:
            if row["severity"] not in ALLOWED_SEVERITIES:
                errors.append(f"appendix/question-ledger.md: {row['id']} has invalid severity {row['severity']!r}")
            if row["status"] not in ALLOWED_STATUSES:
                errors.append(f"appendix/question-ledger.md: {row['id']} has invalid status {row['status']!r}")
            if row["status"] == "open" and row["severity"] in BLOCKING_SEVERITIES:
                open_questions = (root / "open-questions.md").read_text(encoding="utf-8") if (root / "open-questions.md").exists() else ""
                if row["id"] not in open_questions:
                    errors.append(f"open-questions.md: missing open blocking question {row['id']}")
    if status == "approved-to-dev":
        if inventory.exists():
            inventory_text = inventory.read_text(encoding="utf-8")
            found_markers = sorted(marker for marker in MISSING_SOURCE_MARKERS if marker in inventory_text)
            if found_markers:
                errors.append(f"source-inventory.md: approved-to-dev with missing source marker(s): {', '.join(found_markers)}")
            if not included_source_entries(inventory_text, root):
                errors.append("source-inventory.md: approved-to-dev requires at least one concrete existing included source entry")
        blockers = [row["id"] for row in rows if row["status"] == "open" and row["severity"] in BLOCKING_SEVERITIES]
        if blockers:
            errors.append(f"index.md: approved-to-dev with open blocking questions: {', '.join(blockers)}")
        findings = root / "appendix" / "subagent-findings.md"
        if findings.exists():
            reviewer_rows = parse_table_rows(findings.read_text(encoding="utf-8"), ["Scope", "Status", "Reviewer Verdict", "Report"])
            for duplicate in sorted(duplicate_values(reviewer_rows, "scope")):
                errors.append(f"appendix/subagent-findings.md: duplicate reviewer verdict row for {duplicate}")
            by_scope = {row["scope"]: row for row in reviewer_rows}
            for scope in sorted(REQUIRED_SCOPES):
                row = by_scope.get(scope)
                if row is None:
                    errors.append(f"appendix/subagent-findings.md: approved-to-dev missing reviewer verdict for {scope}")
                    continue
                if row["status"] != "approved" or row["reviewer_verdict"] != "approved":
                    errors.append(f"appendix/subagent-findings.md: {scope} must be approved for approved-to-dev")
                if not row["report"] or row["report"] in {"PLACEHOLDER", "TBD", "none"}:
                    errors.append(f"appendix/subagent-findings.md: {scope} missing reviewer report evidence")
                elif not report_path_valid(scope, row["report"], root):
                    errors.append(f"appendix/subagent-findings.md: {scope} reviewer report must point to an existing .tasks/technical-docs-verify/<run-id>/scopes/{scope}/review-attempt-<n>.md")
        loop_history = root / "appendix" / "loop-history.md"
        if loop_history.exists():
            gate_rows = parse_table_rows(loop_history.read_text(encoding="utf-8"), ["Gate", "Status", "Evidence"])
            for duplicate in sorted(duplicate_values(gate_rows, "gate")):
                errors.append(f"appendix/loop-history.md: duplicate approval gate row for {duplicate}")
            by_gate = {row["gate"]: row for row in gate_rows}
            for gate in sorted(REQUIRED_APPROVAL_GATES):
                row = by_gate.get(gate)
                if row is None:
                    errors.append(f"appendix/loop-history.md: approved-to-dev missing approval gate {gate}")
                    continue
                if row["status"] != "passed":
                    errors.append(f"appendix/loop-history.md: approval gate {gate} must be passed for approved-to-dev")
                if not row["evidence"] or row["evidence"] in {"PLACEHOLDER", "TBD", "none"}:
                    errors.append(f"appendix/loop-history.md: approval gate {gate} missing evidence")
                if gate in DELTA_GATE_KEYWORDS and not evidence_matches_keywords(row["evidence"], DELTA_GATE_KEYWORDS[gate]):
                    errors.append(f"appendix/loop-history.md: approval gate {gate} evidence must cite source-delta or question-effect review evidence")

    return errors


# START_CONTRACT: main
#   PURPOSE: Validate a technical-verified package from the command line.
#   INPUTS: { argv: CLI args - root path and --allow-placeholders }
#   OUTPUTS: { int - process exit code }
#   SIDE_EFFECTS: Prints validation results to stdout.
#   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER.
# END_CONTRACT: main
def main() -> int:
    parser = argparse.ArgumentParser(description="Validate technical-verified docs.")
    parser.add_argument("root", help="Technical verified docs folder.")
    parser.add_argument("--allow-placeholders", action="store_true", help="Allow scaffold placeholders.")
    args = parser.parse_args()

    errors = validate_root(Path(args.root), args.allow_placeholders)
    if errors:
        print("technical-verified validation failed:")
        for error in errors:
            print(f"- {error}")
        return 1
    print("technical-verified validation passed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
