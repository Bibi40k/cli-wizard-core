# 2026-03-06 - Wizard Interrupt Contract

## Status
Accepted

## Context
Multiple repositories reuse wizard-style CLI flows. Regressions appeared where `Ctrl+C` behaved inconsistently:
- sometimes local cancel;
- sometimes back to previous menu;
- sometimes full exit.

This creates operator confusion and unpredictable automation behavior.

## Decision
`Ctrl+C` is standardized as a strict interrupt contract:
1. In any wizard step prompt, first `Ctrl+C` must abort the current flow immediately.
2. Interrupt must propagate up to the top-level manager unless the flow explicitly documents local-cancel semantics.
3. No silent conversion of interrupt into default value, `Back`, or `continue`.
4. `Back`/`Exit` remain explicit menu choices only (keyboard navigation + Enter).
5. Draft save behavior is explicit per flow:
   - create flow: may save draft;
   - edit flow: no draft side-effects unless explicitly designed and documented.

## Required Acceptance Checks
For each consumer repo integrating this library:
1. `Ctrl+C` from root menu exits immediately.
2. `Ctrl+C` from nested menu exits according to contract (no accidental back-loop).
3. `Ctrl+C` from selector prompts (`survey`/custom raw mode) exits on first keypress.
4. `Ctrl+C` from line prompts (`readline`) exits on first keypress.
5. No save/rename side-effects after interrupt.

## Consequences
- Predictable UX across repositories.
- Fewer regressions when adding new flows.
- Interrupt behavior becomes testable and reviewable as a contract.
