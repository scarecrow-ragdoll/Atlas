# Main Orchestration

## Run ID

20260618T222231Z

## Selected Wave

WAVE-05: Nutrition

## Phase

preparing — source wave gate, context inventory, scaffold

## Source Wave Gate

See source-wave-gate.md

## Context Inventory

See context-inventory.md

## Next Action

Source wave gate passed. Inventoried context. Run scaffold script, then dispatch wave-orchestrator for WAVE-05.

## Interruption Recovery

If interrupted before wave-orchestrator completion:
1. Resume with RUN_ID=20260618T222231Z, WAVE_ID=WAVE-05
2. Read main-orchestration.md, source-wave-gate.md, context-inventory.md
3. Re-run the scaffold if missing
4. Relaunch wave-orchestrator for remaining planner/reviewer work