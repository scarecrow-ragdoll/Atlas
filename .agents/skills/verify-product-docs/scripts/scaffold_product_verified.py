#!/usr/bin/env python3
"""Create the required docs/product-verified structure."""

from __future__ import annotations

import argparse
from pathlib import Path


DOC_EXTENSIONS = {
    ".md",
    ".mdx",
    ".txt",
    ".xml",
    ".json",
    ".yaml",
    ".yml",
    ".csv",
}


TEMPLATES: dict[str, str] = {
    "index.md": """# Product Verified

## Status

TBD

## Source Set

TBD

## Document Map

TBD

## Handoff Readiness

TBD
""",
    "product-brief.md": """# Product Brief

## Product Intent

TBD

## Target Users

TBD

## Jobs To Be Done

TBD

## Value Proposition

TBD

## Success Metrics

TBD
""",
    "scope.md": """# Scope

## In Scope

TBD

## Out Of Scope

TBD

## Non-Goals

TBD

## Dependencies

TBD

## Assumptions

TBD
""",
    "actors-and-permissions.md": """# Actors And Permissions

## Actors

TBD

## Roles

TBD

## Permissions Matrix

TBD

## Ownership Rules

TBD

## Privacy And Security Expectations

TBD
""",
    "domain-model.md": """# Domain Model

## Entities

TBD

## Attributes

TBD

## Relationships

TBD

## Lifecycle States

TBD

## Invariants

TBD
""",
    "functional-spec.md": """# Functional Specification

## Capability Map

TBD

## Feature Behavior

TBD

## Validations

TBD

## Notifications

TBD

## Integrations

TBD
""",
    "user-flows.md": """# User Flows

## Primary Flows

TBD

## Alternative Flows

TBD

## Failure And Recovery Flows

TBD

## Empty States

TBD
""",
    "business-rules.md": """# Business Rules

## Validation Rules

TBD

## Calculation Rules

TBD

## State Transition Rules

TBD

## Authorization Rules

TBD

## Integration Rules

TBD
""",
    "edge-cases.md": """# Edge Cases

## Input And Validation

TBD

## Permissions And Ownership

TBD

## State And Concurrency

TBD

## External Dependencies

TBD

## Data Lifecycle

TBD
""",
    "acceptance-criteria.md": """# Acceptance Criteria

## Product-Level Criteria

TBD

## Feature-Level Criteria

TBD

## Negative Criteria

TBD

## Handoff Criteria

TBD
""",
    "open-questions.md": """# Open Questions

## Missing Source Artifacts

TBD

## Blocking

TBD

## Non-Blocking

TBD

## Deferred

TBD
""",
    "features/index.md": """# Features

## Feature Inventory

TBD

## Feature File Map

TBD
""",
    "appendix/subagent-findings.md": """# Subagent Findings

## Reports

TBD

## Cross-Reviewer Conflicts

TBD

## Synthesis Notes

TBD
""",
    "appendix/traceability.md": """# Traceability

## Requirement Map

TBD

## Source Map

TBD

## Assumption Map

TBD

## Open Question Map

TBD
""",
    "appendix/derivation-log.md": """# Derivation Log

## Derived Roles And Permissions

TBD

## Derived Data Fields

TBD

## Derived States

TBD

## Derived Acceptance Criteria

TBD

## Derived Edge Cases

TBD

## Low-Confidence Derivations

TBD
""",
    "appendix/question-ledger.md": """# Question Ledger

## Missing Source Artifacts

TBD

## Blocking Questions

TBD

## Non-Blocking Questions

TBD

## Resolved Questions

TBD

## Deferred Questions

TBD
""",
    "appendix/decision-log.md": """# Decision Log

## Resolved Contradictions

TBD

## Assumptions Adopted

TBD

## Rejected Or Outdated Inputs

TBD
""",
}


def discover_sources(source: Path) -> tuple[list[Path], list[tuple[Path, str]]]:
    if not source.exists():
        raise SystemExit(f"Source folder does not exist: {source}")
    if not source.is_dir():
        raise SystemExit(f"Source path is not a directory: {source}")

    included: list[Path] = []
    excluded: list[tuple[Path, str]] = []

    for path in source.rglob("*"):
        if not path.is_file():
            continue
        rel_parts = path.relative_to(source).parts
        if any(part.startswith(".") for part in rel_parts):
            excluded.append((path, "hidden path"))
            continue
        if path.suffix.lower() not in DOC_EXTENSIONS:
            excluded.append((path, "unsupported extension"))
            continue
        included.append(path)

    sort_key = lambda item: item[0].relative_to(source).as_posix()
    return (
        sorted(included, key=lambda path: path.relative_to(source).as_posix()),
        sorted(excluded, key=sort_key),
    )


def line_count(path: Path) -> int:
    try:
        text = path.read_text(encoding="utf-8")
    except UnicodeDecodeError:
        text = path.read_text(encoding="utf-8", errors="replace")
    return text.count("\n") + (0 if text.endswith("\n") or text == "" else 1)


def source_inventory(source: Path, files: list[Path], excluded: list[tuple[Path, str]]) -> str:
    rows = []
    for path in files:
        rel = path.relative_to(source).as_posix()
        rows.append(f"| `{rel}` | {path.stat().st_size} | {line_count(path)} |")

    table = "\n".join(rows) if rows else "| _No supported source files found_ | 0 | 0 |"
    excluded_rows = []
    for path, reason in excluded:
        rel = path.relative_to(source).as_posix()
        excluded_rows.append(f"| `{rel}` | {reason} |")
    excluded_table = "\n".join(excluded_rows) if excluded_rows else "| _None detected_ | - |"

    return f"""# Source Inventory

## Included Sources

| File | Bytes | Lines |
| --- | ---: | ---: |
{table}

## Excluded Or Noisy Sources

| File | Reason |
| --- | --- |
{excluded_table}

## Source Delta

TBD

## Coverage Gaps

TBD
"""


def write_file(path: Path, content: str, force: bool) -> bool:
    if path.exists() and not force:
        return False
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content.rstrip() + "\n", encoding="utf-8")
    return True


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--source", default="docs/product", help="Raw product docs folder")
    parser.add_argument("--output", default="docs/product-verified", help="Verified docs folder")
    parser.add_argument("--force", action="store_true", help="Overwrite managed files")
    args = parser.parse_args()

    source = Path(args.source)
    output = Path(args.output)
    files, excluded = discover_sources(source)

    written: list[str] = []
    skipped: list[str] = []

    inventory_path = output / "source-inventory.md"
    if write_file(inventory_path, source_inventory(source, files, excluded), args.force):
        written.append(inventory_path.as_posix())
    else:
        skipped.append(inventory_path.as_posix())

    for rel_path, content in TEMPLATES.items():
        target = output / rel_path
        if write_file(target, content, args.force):
            written.append(target.as_posix())
        else:
            skipped.append(target.as_posix())

    print(f"source_files={len(files) + len(excluded)}")
    print(f"included_files={len(files)}")
    print(f"excluded_files={len(excluded)}")
    print(f"written={len(written)}")
    print(f"skipped={len(skipped)}")
    for path in written:
        print(f"+ {path}")
    for path in skipped:
        print(f"= {path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
