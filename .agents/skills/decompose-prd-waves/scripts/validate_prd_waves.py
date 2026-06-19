#!/usr/bin/env python3
# FILE: .agents/skills/decompose-prd-waves/scripts/validate_prd_waves.py
# VERSION: 1.0.0
# START_MODULE_CONTRACT
#   PURPOSE: Validate the required PRD waves documentation structure and shallow decomposition invariants.
#   SCOPE: Checks files, headings, placeholder policy, status values, reviewer approvals, shallow-only banned terms, stable ids, frontend page files, source traceability, and open blocking question consistency; excludes semantic truth review.
#   DEPENDS: Python standard library, .agents/skills/decompose-prd-waves/references/output-contract.md.
#   LINKS: M-PRD-WAVE-DECOMPOSER / V-M-PRD-WAVE-DECOMPOSER.
#   ROLE: SCRIPT
#   MAP_MODE: EXPORTS
# END_MODULE_CONTRACT
#
# START_MODULE_MAP
#   REQUIRED_HEADINGS - Required markdown headings by relative package path.
#   WAVE_HEADINGS - Required markdown headings for each wave file.
#   FRONTEND_PAGE_HEADINGS - Required markdown headings for each frontend page file.
#   VALID_PACKAGE_STATUSES - Package statuses accepted in index.md.
#   VALID_WAVE_STATUSES - Wave statuses accepted in wave files.
#   BLOCKING_SEVERITIES - Question severities that block top-level approval while open.
#   FINAL_PACKAGE_STATUSES - Package statuses that require reviewer and shallow gates.
#   FINAL_WAVE_STATUSES - Wave statuses that require stable ids and shallow gates.
#   REQUIRED_SCOPES - Required primary scopes and consistency scope for final approval.
#   REQUIRED_REVIEWERS - Required reviewer perspectives for final approval.
#   VALID_REVIEWER_VERDICTS - Reviewer verdicts accepted in final package files.
#   DRAFT_REVIEWER_VERDICTS - Reviewer verdicts accepted only when placeholders are allowed.
#   BANNED_DETAIL_PATTERNS - Regex patterns that identify implementation-level detail in shallow output.
#   FRONTEND_IN_BACKEND_WAVE_PATTERNS - Regex patterns that identify frontend scope leaking into backend waves.
#   APPROVAL_MARKERS - Explicit user approval tokens accepted for waves-approved handoff.
#   artifact_path - Resolve package-relative evidence paths against the run workspace root.
#   extract_markdown_section - Extract a markdown heading section body.
#   has_approval_marker - Detect explicit user approval markers in text.
#   has_negative_approval_marker - Detect text that negates or defers approval.
#   has_meaningful_final_evidence - Detect non-placeholder final source or section evidence.
#   main - Parse arguments and report validation failures.
#   validate_root - Validate required files, headings, placeholders, and approval invariants.
#   parse_question_rows - Parse markdown question ledger rows for gate checks.
#   status_from_heading - Extract the status that follows a Status heading.
#   wave_code_from_path - Convert a wave file path to its canonical WAVE-NN id.
#   wave_number_from_path - Convert a wave file path to its numeric id string.
#   page_code_from_path - Convert a page file path to its canonical PAGE-NNN id.
#   reviewer_verdict_errors - Validate reviewer verdict values.
#   has_required_reviewer_approvals - Find missing approved reviewer perspectives.
#   reviewer_report_path_errors - Validate required reviewer report paths and artifact existence.
#   candidate_package_errors - Validate consistency review evidence against a concrete candidate package.
#   scope_review_errors - Validate required scope-review approvals and report-path evidence.
#   collect_question_rows - Collect question rows from aggregate, package, and wave-local ledgers.
#   question_consistency_errors - Validate blocking question synchronization and row consistency.
#   shallow_detail_errors - Reject implementation-level detail in shallow wave outputs.
#   package_shallow_detail_errors - Reject implementation-level detail in all package markdown files.
#   backend_wave_scope_errors - Reject frontend scope in backend-only wave outputs.
#   frontend_pages_errors - Validate the per-page frontend handoff files.
#   user_approval_errors - Validate explicit user approval evidence for waves-approved packages.
# END_MODULE_MAP
# START_CHANGE_SUMMARY
#   LAST_CHANGE: 1.0.4 - Replaced the single frontend sequence output with per-page frontend files sourced from raw and verified PRDs.
# END_CHANGE_SUMMARY

from __future__ import annotations

import argparse
import re
from pathlib import Path


REQUIRED_HEADINGS: dict[str, list[str]] = {
    "index.md": ["# PRD Waves", "## Status", "## Source Gate", "## Shallow Wave Gate", "## Wave Count", "## Source Set", "## Next Action"],
    "source-inventory.md": ["# Source Inventory", "## Raw Product Sources", "## Verified Product Sources", "## Technical Sources", "## Prior Wave Sources", "## Source Delta", "## Source Gaps"],
    "scope-inventory.md": ["# Scope Inventory", "## Capability Groups", "## User Journey Groups", "## Data Lifecycle Groups", "## Integration And Operations Groups", "## Client Experience Groups", "## Security Compliance Groups", "## Explicit Deferrals"],
    "wave-map.md": ["# Wave Map", "## Top-Level Wave List", "## Dependency Order", "## Coverage Matrix", "## More Than Eight Wave Check", "## Downstream Planning Recommendations"],
    "frontend-pages/index.md": ["# Frontend Pages", "## Status", "## Scope Source", "## Page Order", "## Raw PRD Source Coverage", "## Verified PRD Source Coverage", "## Shared UX States", "## Backend Dependencies By Page", "## Explicit Frontend Deferrals", "## Open Questions", "## Traceability"],
    "open-questions.md": ["# Open Questions", "## Decomposition Blocking", "## Needs Owner Decision", "## Deferred", "## Watchlist", "## Resolved This Run"],
    "waves/index.md": ["# Waves", "## Wave List", "## Dependency Order", "## Approval State"],
    "appendix/reviewer-verdicts.md": ["# Reviewer Verdicts", "## Scope Reviews", "## Consistency Review", "## Rejected Findings"],
    "appendix/traceability.md": ["# Traceability", "## Source To Scope Map", "## Scope To Wave Map", "## Wave To Source Map", "## Question Map", "## Decision Map"],
    "appendix/question-ledger.md": ["# Question Ledger", "## Open Questions", "## Answered Questions", "## Follow-Up Questions", "## Resolved Questions", "## Deferred Questions"],
    "appendix/decision-log.md": ["# Decision Log", "## Source Gate", "## Scope Decisions", "## Deferrals", "## User Wave Map Approvals", "## Rejected Assumptions"],
    "appendix/run-history.md": ["# Run History", "## Runs", "## Scope Mapper Cycles", "## Consistency Cycles", "## Source Delta History", "## Approval Gate History"],
}

WAVE_HEADINGS = [
    "## Status",
    "## User Approval",
    "## Purpose",
    "## Outcome After Wave",
    "## Included Scope",
    "## Excluded Scope",
    "## Dependencies",
    "## Surface Categories",
    "## Risk Class",
    "## Recommended Next Planning",
    "## Open Questions",
    "## Traceability",
]

FRONTEND_PAGE_HEADINGS = [
    "## Status",
    "## Page Purpose",
    "## What Is On This Page",
    "## Functional Parts",
    "## Empty States",
    "## Loading And Error States",
    "## Backend Dependencies",
    "## Explicit Deferrals",
    "## Open Questions",
    "## Raw PRD Traceability",
    "## Verified PRD Traceability",
]

VALID_PACKAGE_STATUSES = {
    "draft",
    "questions-open",
    "blocked",
    "waves-ready-awaiting-user-approval",
    "waves-approved",
    "superseded",
}

VALID_WAVE_STATUSES = {
    "draft",
    "needs-revision",
    "questions-open",
    "blocked",
    "top-level-ready",
    "user-approved",
    "superseded",
}

BLOCKING_SEVERITIES = {"decomposition-blocking", "needs-owner-decision"}
FINAL_PACKAGE_STATUSES = {"waves-ready-awaiting-user-approval", "waves-approved"}
FINAL_WAVE_STATUSES = {"top-level-ready", "user-approved"}
REQUIRED_SCOPES = [
    "product-capabilities",
    "user-journeys",
    "data-lifecycle",
    "integrations-operations",
    "client-experience",
    "security-compliance",
    "delivery-sequencing",
    "wave-map-consistency",
]
REQUIRED_REVIEWERS = [
    "product-scope-coverage",
    "technical-boundary-fit",
    "sequencing-dependencies",
    "backend-wave-boundary-quality",
    "traceability-consistency",
]
VALID_REVIEWER_VERDICTS = {"approved", "needs-revision", "blocked"}
DRAFT_REVIEWER_VERDICTS = {"pending-review"}

BANNED_DETAIL_PATTERNS = [
    r"\bimplementation tasks?\b",
    r"\bmodule designs?\b",
    r"\bmigration plans?\b",
    r"\bacceptance criteria\b",
    r"\bexit criteria\b",
    r"\btest cases?\b",
    r"\btest plans?\b",
    r"\bapi payloads?\b",
    r"\bschemas?\b",
    r"\bcomponent architecture\b",
    r"\bcomponent designs?\b",
    r"\bBTASK-[A-Z0-9-]+\b",
    r"\bAC-W\d+-\d+\b",
    r"\bEXIT-W\d+-\d+\b",
    r"\bBTEST-W\d+-\d+\b",
    r"\bJira\b",
    r"\bBeads?\b",
    r"\bCREATE TABLE\b",
    r"\bALTER TABLE\b",
    r"\bmutation\b",
    r"\bresolver\b",
    r"\bendpoint\b",
    r"\brequest payload\b",
    r"\bresponse payload\b",
]

FRONTEND_IN_BACKEND_WAVE_PATTERNS = [
    r"\bfront[- ]?end\b",
    r"\bui\b",
    r"\bux\b",
    r"\buser interfaces?\b",
    r"\bclient[- ]facing\b",
    r"\bpages?\b",
    r"\bscreens?\b",
    r"\broutes?\b",
    r"\bnavigation\b",
    r"\bmobile\b",
    r"\bclient[- ]experience\b",
    r"\bclient experience\b",
]

APPROVAL_MARKERS = {"approved-by-user", "user-approved", "waves-approved-by-user"}


def artifact_path(root: Path, raw: str) -> Path:
    path = Path(raw)
    return path if path.is_absolute() else root.parent.parent / path


def extract_markdown_section(text: str, heading: str) -> str:
    pattern = rf"^{re.escape(heading)}\s*\n(.*?)(?=^## |\Z)"
    match = re.search(pattern, text, flags=re.MULTILINE | re.DOTALL)
    return match.group(1).strip() if match else ""


def has_approval_marker(text: str) -> bool:
    lower = text.lower()
    return any(marker in lower for marker in APPROVAL_MARKERS)


def has_negative_approval_marker(text: str) -> bool:
    lower = text.lower()
    return bool(re.search(r"\bnot\s+(?:yet\s+)?(?:approved|user-approved)\b|\bawaiting approval\b|\bpending approval\b", lower))


# START_CONTRACT: parse_question_rows
#   PURPOSE: Extract question table rows from a PRD wave question ledger.
#   INPUTS: { text: str - markdown ledger content }
#   OUTPUTS: { list[dict[str, str]] - parsed rows keyed by expected table columns }
#   SIDE_EFFECTS: None.
#   LINKS: M-PRD-WAVE-DECOMPOSER / V-M-PRD-WAVE-DECOMPOSER.
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
#   LINKS: M-PRD-WAVE-DECOMPOSER / V-M-PRD-WAVE-DECOMPOSER.
# END_CONTRACT: status_from_heading
def status_from_heading(text: str) -> str | None:
    match = re.search(r"^## Status\s*\n([^\n]+)", text, flags=re.MULTILINE)
    if not match:
        return None
    return match.group(1).strip()


def wave_code_from_path(path: Path) -> str:
    match = re.search(r"wave-(\d+)", path.stem)
    if not match:
        return path.stem.upper()
    return f"WAVE-{match.group(1)}"


def wave_number_from_path(path: Path) -> str:
    match = re.search(r"wave-(\d+)", path.stem)
    if not match:
        return "XX"
    return match.group(1)


def page_code_from_path(path: Path) -> str:
    match = re.search(r"page-(\d+)", path.stem)
    if not match:
        return path.stem.upper()
    return f"PAGE-{match.group(1).zfill(3)}"


def has_meaningful_final_evidence(text: str) -> bool:
    if re.search(r"\b(PLACEHOLDER|SOURCE_MISSING|SOURCE_EMPTY|TODO|TBD)\b", text):
        return False
    empty_markers = {"none", "n/a", "na", "not applicable", "not specified", "no evidence", "no traceability"}
    for raw_line in text.splitlines():
        line = raw_line.strip()
        for _ in range(3):
            normalized_marker = re.sub(r"^>\s*", "", line)
            normalized_marker = re.sub(r"^(?:[-*]|\d+[.)])\s*", "", normalized_marker)
            normalized_marker = normalized_marker.strip()
            if normalized_marker == line:
                break
            line = normalized_marker
        normalized = line.lower().strip(" .,:;")
        if normalized and normalized not in empty_markers:
            return True
    return False


def reviewer_verdict_errors(text: str, relative: str, allow_placeholders: bool) -> list[str]:
    errors: list[str] = []
    for line in text.splitlines():
        if not line.startswith("|"):
            continue
        cells = [cell.strip().lower() for cell in line.strip().strip("|").split("|")]
        if len(cells) != 7 or cells[0] in {"scope", "---"}:
            continue
        verdict = cells[3]
        if verdict in VALID_REVIEWER_VERDICTS:
            continue
        if allow_placeholders and verdict in DRAFT_REVIEWER_VERDICTS:
            continue
        errors.append(f"{relative}: invalid reviewer verdict {verdict!r}")
    return errors


def has_required_reviewer_approvals(text: str) -> list[str]:
    approved_reviewers: set[str] = set()
    for line in text.splitlines():
        if not line.startswith("|"):
            continue
        cells = [cell.strip().lower() for cell in line.strip().strip("|").split("|")]
        if len(cells) != 7 or cells[0] in {"scope", "---"}:
            continue
        perspective = cells[1]
        verdict = cells[3]
        if verdict == "approved":
            approved_reviewers.add(perspective)
    return [reviewer for reviewer in REQUIRED_REVIEWERS if reviewer not in approved_reviewers]


def reviewer_report_path_errors(report: str, expected_pattern: str, root: Path, label: str) -> list[str]:
    errors: list[str] = []
    if not re.search(expected_pattern, report):
        errors.append(f"appendix/reviewer-verdicts.md: invalid reviewer report path for {label}")
        return errors
    candidate = artifact_path(root, report)
    if not candidate.exists():
        errors.append(f"appendix/reviewer-verdicts.md: reviewer report path for {label} does not exist: {report}")
    return errors


def candidate_package_errors(report: str, root: Path, label: str) -> list[str]:
    errors: list[str] = []
    report_path = artifact_path(root, report)
    if not report_path.exists():
        return errors
    report_text = report_path.read_text(encoding="utf-8")
    match = re.search(r"\.tasks/prd-wave-decomposition/[^/\s)]+/staging/prd-waves", report_text)
    if not match:
        errors.append(f"appendix/reviewer-verdicts.md: consistency report for {label} must name reviewed candidate package")
        return errors
    candidate = artifact_path(root, match.group(0))
    if not candidate.exists():
        errors.append(f"appendix/reviewer-verdicts.md: reviewed candidate package for {label} does not exist: {match.group(0)}")
    return errors


def scope_review_errors(text: str, root: Path) -> list[str]:
    errors: list[str] = []
    approved_scope_reports: dict[str, str] = {}
    approved_perspective_reports: dict[str, str] = {}
    for line in text.splitlines():
        if not line.startswith("|"):
            continue
        cells = [cell.strip() for cell in line.strip().strip("|").split("|")]
        if len(cells) != 7 or cells[0].lower() in {"scope", "---"}:
            continue
        scope = cells[0].lower()
        perspective = cells[1].lower()
        verdict = cells[3].lower()
        report = cells[4]
        if verdict != "approved":
            continue
        if perspective in {"scope-review", "consistency-review"}:
            approved_scope_reports[scope] = report
        if perspective in REQUIRED_REVIEWERS:
            approved_perspective_reports[perspective] = report

    for scope in REQUIRED_SCOPES:
        report = approved_scope_reports.get(scope)
        if not report:
            errors.append(f"appendix/reviewer-verdicts.md: missing approved scope review for {scope}")
            continue
        escaped_scope = re.escape(scope)
        if scope == "wave-map-consistency":
            expected = rf"\.tasks/prd-wave-decomposition/[^/]+/scopes/{escaped_scope}/consistency-attempt-\d+\.md"
        else:
            expected = rf"\.tasks/prd-wave-decomposition/[^/]+/scopes/{escaped_scope}/review-attempt-\d+\.md"
        errors.extend(reviewer_report_path_errors(report, expected, root, scope))

    consistency_expected = r"\.tasks/prd-wave-decomposition/[^/]+/scopes/wave-map-consistency/consistency-attempt-\d+\.md"
    for perspective in REQUIRED_REVIEWERS:
        report = approved_perspective_reports.get(perspective)
        if not report:
            continue
        errors.extend(reviewer_report_path_errors(report, consistency_expected, root, perspective))

    candidate_reports: dict[str, str] = {}
    consistency_report = approved_scope_reports.get("wave-map-consistency")
    if consistency_report:
        candidate_reports["wave-map-consistency"] = consistency_report
    for perspective, report in approved_perspective_reports.items():
        candidate_reports[perspective] = report
    checked_reports: set[str] = set()
    for label, report in candidate_reports.items():
        if report in checked_reports:
            continue
        checked_reports.add(report)
        errors.extend(candidate_package_errors(report, root, label))

    return errors


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


def shallow_detail_errors(text: str, relative: str, allow_placeholders: bool) -> list[str]:
    if allow_placeholders:
        return []
    errors: list[str] = []
    scan_text = re.sub(r"\b[\w./-]+\.(?:md|xml|yaml|yml|json|ts|tsx|js|jsx|py|go|sql)\b", " ", text)
    for pattern in BANNED_DETAIL_PATTERNS:
        if re.search(pattern, scan_text, flags=re.IGNORECASE):
            errors.append(f"{relative}: contains implementation-level detail matching {pattern}")
    return errors


def package_shallow_detail_errors(root: Path, allow_placeholders: bool) -> list[str]:
    errors: list[str] = []
    for path in sorted(root.rglob("*.md")):
        if not path.is_file():
            continue
        relative = path.relative_to(root).as_posix()
        if relative == "source-inventory.md":
            continue
        errors.extend(shallow_detail_errors(path.read_text(encoding="utf-8"), relative, allow_placeholders))
    return errors


def backend_wave_scope_errors(root: Path, wave_files: list[Path], allow_placeholders: bool) -> list[str]:
    if allow_placeholders:
        return []
    errors: list[str] = []
    backend_wave_files = [root / "wave-map.md", root / "waves" / "index.md"] + wave_files
    for path in backend_wave_files:
        if not path.exists():
            continue
        relative = path.relative_to(root).as_posix()
        text = path.read_text(encoding="utf-8")
        scan_text = re.sub(r"\b[\w./-]+\.(?:md|xml|yaml|yml|json|ts|tsx|js|jsx|py|go|sql)\b", " ", text)
        for pattern in FRONTEND_IN_BACKEND_WAVE_PATTERNS:
            if re.search(pattern, scan_text, flags=re.IGNORECASE):
                errors.append(f"{relative}: backend waves must not include frontend/page scope matching {pattern}")
    return errors


def frontend_pages_errors(root: Path, package_status: str | None, allow_placeholders: bool) -> list[str]:
    legacy_path = root / "frontend-page-sequence.md"
    index_path = root / "frontend-pages" / "index.md"
    errors: list[str] = []
    if legacy_path.exists():
        errors.append("frontend-page-sequence.md: deprecated output; use frontend-pages/index.md and frontend-pages/page-<nnn>.md")
    if allow_placeholders:
        return errors
    if package_status not in FINAL_PACKAGE_STATUSES:
        return errors

    if not index_path.exists():
        return errors

    text = index_path.read_text(encoding="utf-8")
    status = status_from_heading(text)
    if status not in FINAL_WAVE_STATUSES:
        errors.append("frontend-pages/index.md: final package requires frontend pages status top-level-ready or user-approved")

    page_order = extract_markdown_section(text, "## Page Order")
    raw_source_coverage = extract_markdown_section(text, "## Raw PRD Source Coverage")
    verified_source_coverage = extract_markdown_section(text, "## Verified PRD Source Coverage")
    frontend_deferrals = extract_markdown_section(text, "## Explicit Frontend Deferrals")
    deferral_text = f"{status}\n{frontend_deferrals}"
    has_explicit_deferral = bool(re.search(r"\b(FRONTEND_NONE_CONFIRMED|FRONTEND_DEFERRED)\b", deferral_text))
    page_ids = sorted(set(re.findall(r"\bPAGE-\d{3}\b", page_order)))

    if not has_meaningful_final_evidence(raw_source_coverage):
        errors.append("frontend-pages/index.md: final package requires raw PRD source coverage")
    if not has_meaningful_final_evidence(verified_source_coverage):
        errors.append("frontend-pages/index.md: final package requires verified PRD source coverage")

    page_files = sorted((root / "frontend-pages").glob("page-*.md"))
    if has_explicit_deferral and not page_files and not page_ids:
        return errors
    if has_explicit_deferral and page_files and not page_ids:
        errors.append("frontend-pages/index.md: page files exist but Page Order has no PAGE-... entries")
    if not page_ids and not page_files:
        errors.append("frontend-pages/index.md: final package must list PAGE-... entries in Page Order or explicit FRONTEND_NONE_CONFIRMED/FRONTEND_DEFERRED")
        return errors

    for page_id in page_ids:
        number = page_id.split("-")[1].lower()
        expected = root / "frontend-pages" / f"page-{number}.md"
        if not expected.exists():
            errors.append(f"frontend-pages/index.md: {page_id} from Page Order missing frontend-pages/page-{number}.md")

    for page_file in page_files:
        relative = page_file.relative_to(root).as_posix()
        page_text = page_file.read_text(encoding="utf-8")
        page_id = page_code_from_path(page_file)
        if page_id not in page_order:
            errors.append(f"{relative}: {page_id} must appear in frontend-pages/index.md Page Order")
        if page_id not in page_text:
            errors.append(f"{relative}: missing canonical page id {page_id}")
        for heading in FRONTEND_PAGE_HEADINGS:
            if heading not in page_text:
                errors.append(f"{relative}: missing heading {heading!r}")
        if "PLACEHOLDER" in page_text:
            errors.append(f"{relative}: contains PLACEHOLDER")
        if re.search(r"\b(TODO|TBD)\b", page_text):
            errors.append(f"{relative}: contains TODO/TBD")
        page_status = status_from_heading(page_text)
        if page_status not in FINAL_WAVE_STATUSES:
            errors.append(f"{relative}: final package requires page status top-level-ready or user-approved")
        detail_sections = [
            "## Page Purpose",
            "## What Is On This Page",
            "## Functional Parts",
            "## Empty States",
            "## Loading And Error States",
            "## Backend Dependencies",
        ]
        for heading in detail_sections:
            if not has_meaningful_final_evidence(extract_markdown_section(page_text, heading)):
                errors.append(f"{relative}: final page file requires non-placeholder content under {heading}")
        raw_traceability = extract_markdown_section(page_text, "## Raw PRD Traceability")
        verified_traceability = extract_markdown_section(page_text, "## Verified PRD Traceability")
        if not has_meaningful_final_evidence(raw_traceability):
            errors.append(f"{relative}: final page file requires raw PRD traceability")
        if not has_meaningful_final_evidence(verified_traceability):
            errors.append(f"{relative}: final page file requires verified PRD traceability")
    return errors


def user_approval_errors(root: Path, wave_files: list[Path]) -> list[str]:
    errors: list[str] = []
    decision_log = root / "appendix" / "decision-log.md"
    decision_text = decision_log.read_text(encoding="utf-8") if decision_log.exists() else ""
    if has_negative_approval_marker(decision_text) or not has_approval_marker(decision_text):
        errors.append("appendix/decision-log.md: waves-approved requires explicit approved-by-user decision")

    for wave_file in wave_files:
        relative = wave_file.relative_to(root).as_posix()
        text = wave_file.read_text(encoding="utf-8")
        wave_status = status_from_heading(text)
        approval_section = extract_markdown_section(text, "## User Approval")
        approval_lower = approval_section.lower()
        if wave_status != "user-approved":
            errors.append(f"{relative}: waves-approved requires wave status user-approved")
        if has_negative_approval_marker(approval_lower) or not has_approval_marker(approval_section):
            errors.append(f"{relative}: waves-approved requires approved-by-user in User Approval")
    return errors


# START_CONTRACT: validate_root
#   PURPOSE: Validate structure and shallow approval-gate consistency for a PRD waves package.
#   INPUTS: { root: Path - package root, allow_placeholders: bool - whether PLACEHOLDER text is accepted }
#   OUTPUTS: { list[str] - validation error messages }
#   SIDE_EFFECTS: None.
#   LINKS: M-PRD-WAVE-DECOMPOSER / V-M-PRD-WAVE-DECOMPOSER.
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
    package_status = None
    if index.exists():
        index_text = index.read_text(encoding="utf-8")
        package_status = status_from_heading(index_text)
        if package_status not in VALID_PACKAGE_STATUSES:
            errors.append(f"index.md: status must be one of {sorted(VALID_PACKAGE_STATUSES)}")

    source_inventory = root / "source-inventory.md"
    if source_inventory.exists() and not allow_placeholders:
        source_text = source_inventory.read_text(encoding="utf-8")
        for heading, label in [
            ("## Raw Product Sources", "raw product sources"),
            ("## Verified Product Sources", "verified product sources"),
        ]:
            source_section = extract_markdown_section(source_text, heading)
            if not has_meaningful_final_evidence(source_section):
                errors.append(f"source-inventory.md: {label} must exist and be non-empty for final validation")

    wave_files = sorted((root / "waves").glob("wave-*.md"))
    if not wave_files:
        errors.append("waves/: at least one wave-<nn>.md file is required")

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
        if package_status in FINAL_PACKAGE_STATUSES or wave_status in FINAL_WAVE_STATUSES:
            wave_number = wave_number_from_path(wave_file)
            if not re.search(rf"\bCAP-W{wave_number}-\d+\b", text):
                errors.append(f"{relative}: final shallow wave must contain a CAP-W{wave_number}-... id")
            if not re.search(rf"\bOUT-W{wave_number}-\d+\b", text):
                errors.append(f"{relative}: final shallow wave must contain an OUT-W{wave_number}-... id")
            if wave_status not in FINAL_WAVE_STATUSES:
                errors.append(f"{relative}: package final status requires wave status top-level-ready or user-approved")

    errors.extend(package_shallow_detail_errors(root, allow_placeholders))
    errors.extend(backend_wave_scope_errors(root, wave_files, allow_placeholders))
    errors.extend(frontend_pages_errors(root, package_status, allow_placeholders))

    reviewer_path = root / "appendix" / "reviewer-verdicts.md"
    if reviewer_path.exists():
        reviewer_text = reviewer_path.read_text(encoding="utf-8")
        errors.extend(reviewer_verdict_errors(reviewer_text, "appendix/reviewer-verdicts.md", allow_placeholders))
        if package_status in FINAL_PACKAGE_STATUSES:
            errors.extend(scope_review_errors(reviewer_text, root))
            missing = has_required_reviewer_approvals(reviewer_text)
            for reviewer in missing:
                errors.append(f"appendix/reviewer-verdicts.md: missing approved reviewer {reviewer}")

    if package_status in FINAL_PACKAGE_STATUSES and len(wave_files) > 8:
        decision_log = root / "appendix" / "decision-log.md"
        decision_text = decision_log.read_text(encoding="utf-8").lower() if decision_log.exists() else ""
        if "broader-release-scope-approved" not in decision_text and "non-mvp-release-scope-approved" not in decision_text:
            errors.append("appendix/decision-log.md: 9+ top-level backend waves require broader-release-scope-approved or non-mvp-release-scope-approved decision")

    if package_status == "waves-approved":
        errors.extend(user_approval_errors(root, wave_files))

    rows = collect_question_rows(root, wave_files)
    errors.extend(question_consistency_errors(root, rows))
    if package_status in FINAL_PACKAGE_STATUSES:
        for row in rows:
            if row["status"].lower() == "open" and row["severity"].lower() in BLOCKING_SEVERITIES:
                errors.append(f"final package status blocked by open question {row['id']}")

    return errors


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate shallow PRD wave decomposition docs.")
    parser.add_argument("root", help="PRD waves package root.")
    parser.add_argument("--allow-placeholders", action="store_true", help="Allow PLACEHOLDER and pending-review values for staging validation.")
    args = parser.parse_args()

    errors = validate_root(Path(args.root), args.allow_placeholders)
    if errors:
        for error in errors:
            print(error)
        return 1

    print(f"PRD waves package valid: {args.root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
