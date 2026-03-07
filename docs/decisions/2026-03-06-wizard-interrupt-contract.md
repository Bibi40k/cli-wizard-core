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
`Ctrl+C` is standardized with explicit semantics by flow type:
1. **Create config flow**:
   - first `Ctrl+C` aborts immediately;
   - save plaintext draft;
   - exit application immediately.
2. **Edit config flow**:
   - first `Ctrl+C` aborts immediately;
   - persist current in-memory changes to the edited config;
   - exit application immediately.
3. **Menu/submenu navigation flow**:
   - first `Ctrl+C` exits the application immediately.
4. `Esc` acts as explicit `Back` in selector-style menus.
5. No silent conversion of interrupt into default value.
6. `Back`/`Exit` labels remain plain by default; highlight color is applied by the active selector only when selected.
7. Selector hint is standardized as: `[Use arrows to move, type to filter, Esc=Back, Ctrl+C=Exit]`.
8. Shared wizard behavior changes are implemented in `cli-wizard-core` first; consumer repos must not keep divergent local copies of shared selector/interrupt logic.
9. Action-only menus use aligned verb columns via shared formatter (e.g. `Create    ...`, `Edit      ...`, `Delete    ...`) for consistent scanning UX.

## Required Acceptance Checks
For each consumer repo integrating this library:
1. `Ctrl+C` from root menu exits immediately.
2. `Ctrl+C` from nested submenu exits the app (not implicit back loop).
3. `Ctrl+C` during create saves draft and exits app.
4. `Ctrl+C` during edit saves changes and exits app.
5. `Ctrl+C` from selector prompts (`survey`/custom raw mode) is handled on first keypress.
6. `Ctrl+C` from line prompts (`readline`) is handled on first keypress.
7. `Esc` in selector prompts returns Back (or Exit if Back missing).
8. Consumer repos integrate behavior via core dependency updates, not punctual local rewrites.

## Consequences
- Predictable UX across repositories.
- Fewer regressions when adding new flows.
- Interrupt behavior becomes testable and reviewable as a contract.
