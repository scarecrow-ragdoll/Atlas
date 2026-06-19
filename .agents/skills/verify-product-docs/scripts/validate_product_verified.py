#!/usr/bin/env python3
"""Validate the docs/product-verified output contract."""

from __future__ import annotations

import argparse
import re
from pathlib import Path


REQUIRED_HEADINGS: dict[str, list[str]] = {
    "index.md": [
        "# Product Verified",
        "## Status",
        "## Source Set",
        "## Document Map",
        "## Handoff Readiness",
    ],
    "source-inventory.md": [
        "# Source Inventory",
        "## Included Sources",
        "## Excluded Or Noisy Sources",
        "## Source Delta",
        "## Coverage Gaps",
    ],
    "product-brief.md": [
        "# Product Brief",
        "## Product Intent",
        "## Target Users",
        "## Jobs To Be Done",
        "## Value Proposition",
        "## Success Metrics",
    ],
    "scope.md": [
        "# Scope",
        "## In Scope",
        "## Out Of Scope",
        "## Non-Goals",
        "## Dependencies",
        "## Assumptions",
    ],
    "actors-and-permissions.md": [
        "# Actors And Permissions",
        "## Actors",
        "## Roles",
        "## Permissions Matrix",
        "## Ownership Rules",
        "## Privacy And Security Expectations",
    ],
    "domain-model.md": [
        "# Domain Model",
        "## Entities",
        "## Attributes",
        "## Relationships",
        "## Lifecycle States",
        "## Invariants",
    ],
    "functional-spec.md": [
        "# Functional Specification",
        "## Capability Map",
        "## Feature Behavior",
        "## Validations",
        "## Notifications",
        "## Integrations",
    ],
    "user-flows.md": [
        "# User Flows",
        "## Primary Flows",
        "## Alternative Flows",
        "## Failure And Recovery Flows",
        "## Empty States",
    ],
    "business-rules.md": [
        "# Business Rules",
        "## Validation Rules",
        "## Calculation Rules",
        "## State Transition Rules",
        "## Authorization Rules",
        "## Integration Rules",
    ],
    "edge-cases.md": [
        "# Edge Cases",
        "## Input And Validation",
        "## Permissions And Ownership",
        "## State And Concurrency",
        "## External Dependencies",
        "## Data Lifecycle",
    ],
    "acceptance-criteria.md": [
        "# Acceptance Criteria",
        "## Product-Level Criteria",
        "## Feature-Level Criteria",
        "## Negative Criteria",
        "## Handoff Criteria",
    ],
    "open-questions.md": [
        "# Open Questions",
        "## Missing Source Artifacts",
        "## Blocking",
        "## Non-Blocking",
        "## Deferred",
    ],
    "features/index.md": [
        "# Features",
        "## Feature Inventory",
        "## Feature File Map",
    ],
    "appendix/subagent-findings.md": [
        "# Subagent Findings",
        "## Reports",
        "## Cross-Reviewer Conflicts",
        "## Synthesis Notes",
    ],
    "appendix/traceability.md": [
        "# Traceability",
        "## Requirement Map",
        "## Source Map",
        "## Assumption Map",
        "## Open Question Map",
    ],
    "appendix/derivation-log.md": [
        "# Derivation Log",
        "## Derived Roles And Permissions",
        "## Derived Data Fields",
        "## Derived States",
        "## Derived Acceptance Criteria",
        "## Derived Edge Cases",
        "## Low-Confidence Derivations",
    ],
    "appendix/question-ledger.md": [
        "# Question Ledger",
        "## Missing Source Artifacts",
        "## Blocking Questions",
        "## Non-Blocking Questions",
        "## Resolved Questions",
        "## Deferred Questions",
    ],
    "appendix/decision-log.md": [
        "# Decision Log",
        "## Resolved Contradictions",
        "## Assumptions Adopted",
        "## Rejected Or Outdated Inputs",
    ],
}

FEATURE_HEADINGS = [
    "## Source Evidence",
    "## User Problem",
    "## Scope",
    "## Behavior",
    "## Derived Requirements",
    "## Edge Cases",
    "## Acceptance Criteria",
    "## Dependencies",
    "## Open Questions",
]

PLACEHOLDER_PATTERN = re.compile(r"\b(TBD|TODO|FIXME)\b|\[\s*\]|\{[^}\n]+\}")
STABLE_ID_PATTERN = re.compile(r"\b(?:(?:REQ|RULE|EDGE|AC|DEC)-\d{3}|Q(?:-[A-Z]+)?-\d{3})\b")
TRACEABLE_ID_PATTERN = re.compile(r"\b(?:REQ|RULE|EDGE|AC)-\d{3}\b")
EVIDENCE_TOKEN_PATTERN = re.compile(r"\b(?:Source|Subagent|Derivation|Assumption|Open question):")
CONFIDENCE_PATTERN = re.compile(r"\bConfidence:\s*(?:high|medium|low)\b", re.IGNORECASE)
DERIVATION_TOKEN_PATTERN = re.compile(r"\bDerivation:")

KEY_ID_REQUIREMENTS: dict[str, re.Pattern[str]] = {
    "functional-spec.md": re.compile(r"\bREQ-\d{3}\b"),
    "business-rules.md": re.compile(r"\bRULE-\d{3}\b"),
    "edge-cases.md": re.compile(r"\bEDGE-\d{3}\b"),
    "acceptance-criteria.md": re.compile(r"\bAC-\d{3}\b"),
}


def missing_headings(text: str, headings: list[str]) -> list[str]:
    lines = {line.strip() for line in text.splitlines()}
    return [heading for heading in headings if heading not in lines]


def validate_required_file(root: Path, rel_path: str, allow_placeholders: bool) -> list[str]:
    errors: list[str] = []
    path = root / rel_path
    if not path.exists():
        return [f"missing required file: {rel_path}"]
    text = path.read_text(encoding="utf-8")
    for heading in missing_headings(text, REQUIRED_HEADINGS[rel_path]):
        errors.append(f"{rel_path}: missing heading `{heading}`")
    if not allow_placeholders and PLACEHOLDER_PATTERN.search(text):
        errors.append(f"{rel_path}: contains placeholder text")
    return errors


def validate_feature_file(root: Path, path: Path, allow_placeholders: bool) -> list[str]:
    rel_path = path.relative_to(root).as_posix()
    text = path.read_text(encoding="utf-8")
    errors = []
    if not text.startswith("# "):
        errors.append(f"{rel_path}: must start with a level-1 heading")
    for heading in missing_headings(text, FEATURE_HEADINGS):
        errors.append(f"{rel_path}: missing heading `{heading}`")
    if not allow_placeholders and PLACEHOLDER_PATTERN.search(text):
        errors.append(f"{rel_path}: contains placeholder text")
    return errors


def collect_traceable_ids(root: Path) -> set[str]:
    ids: set[str] = set()
    ignored = {
        "appendix/traceability.md",
        "appendix/question-ledger.md",
        "appendix/subagent-findings.md",
    }
    for path in root.rglob("*.md"):
        rel_path = path.relative_to(root).as_posix()
        if rel_path in ignored:
            continue
        ids.update(TRACEABLE_ID_PATTERN.findall(path.read_text(encoding="utf-8")))
    return ids


def validate_key_ids(root: Path) -> list[str]:
    errors: list[str] = []
    for rel_path, pattern in KEY_ID_REQUIREMENTS.items():
        path = root / rel_path
        if path.exists() and not pattern.search(path.read_text(encoding="utf-8")):
            errors.append(f"{rel_path}: missing expected stable ids matching `{pattern.pattern}`")
    return errors


def validate_traceability(root: Path, allow_placeholders: bool) -> list[str]:
    if allow_placeholders:
        return []

    errors: list[str] = []
    path = root / "appendix" / "traceability.md"
    if not path.exists():
        return errors

    text = path.read_text(encoding="utf-8")
    ids = set(STABLE_ID_PATTERN.findall(text))
    if not ids:
        errors.append("appendix/traceability.md: missing stable requirement/rule/edge/acceptance ids")
    if not EVIDENCE_TOKEN_PATTERN.search(text):
        errors.append("appendix/traceability.md: missing evidence links (`Source:`, `Subagent:`, `Derivation:`, `Assumption:`, or `Open question:`)")

    for line_number, line in enumerate(text.splitlines(), start=1):
        if TRACEABLE_ID_PATTERN.search(line) and not EVIDENCE_TOKEN_PATTERN.search(line):
            errors.append(f"appendix/traceability.md:{line_number}: traced id row is missing an evidence token")

    missing_ids = sorted(collect_traceable_ids(root) - set(TRACEABLE_ID_PATTERN.findall(text)))
    if missing_ids:
        joined = ", ".join(missing_ids[:20])
        suffix = " ..." if len(missing_ids) > 20 else ""
        errors.append(f"appendix/traceability.md: missing ids used elsewhere: {joined}{suffix}")

    if DERIVATION_TOKEN_PATTERN.search(text):
        derivation_log = root / "appendix" / "derivation-log.md"
        derivation_text = derivation_log.read_text(encoding="utf-8") if derivation_log.exists() else ""
        if not CONFIDENCE_PATTERN.search(derivation_text):
            errors.append("appendix/derivation-log.md: derivations referenced from traceability require `Confidence: high|medium|low`")

    return errors


def validate_question_sync(root: Path, allow_placeholders: bool) -> list[str]:
    if allow_placeholders:
        return []

    open_questions = root / "open-questions.md"
    question_ledger = root / "appendix" / "question-ledger.md"
    if not open_questions.exists() or not question_ledger.exists():
        return []

    open_ids = set(re.findall(r"\bQ(?:-[A-Z]+)?-\d{3}\b", open_questions.read_text(encoding="utf-8")))
    ledger_text = question_ledger.read_text(encoding="utf-8")
    missing = sorted(question_id for question_id in open_ids if question_id not in ledger_text)
    if not missing:
        return []
    return [f"appendix/question-ledger.md: missing open-question ids: {', '.join(missing)}"]


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("output", nargs="?", default="docs/product-verified")
    parser.add_argument("--allow-placeholders", action="store_true")
    args = parser.parse_args()

    root = Path(args.output)
    if not root.exists():
        print(f"missing output folder: {root}")
        return 1
    if not root.is_dir():
        print(f"output path is not a directory: {root}")
        return 1

    errors: list[str] = []
    for rel_path in REQUIRED_HEADINGS:
        errors.extend(validate_required_file(root, rel_path, args.allow_placeholders))

    features_dir = root / "features"
    feature_files = [
        path
        for path in sorted(features_dir.glob("*.md"))
        if path.name != "index.md" and not path.name.startswith("_")
    ]
    for path in feature_files:
        errors.extend(validate_feature_file(root, path, args.allow_placeholders))

    if not args.allow_placeholders:
        errors.extend(validate_key_ids(root))
    errors.extend(validate_traceability(root, args.allow_placeholders))
    errors.extend(validate_question_sync(root, args.allow_placeholders))

    if errors:
        for error in errors:
            print(f"ERROR: {error}")
        return 1

    print(f"OK: {root} matches product-verified structure")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
