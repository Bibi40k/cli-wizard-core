Release Notes

Unreleased
Highlights:

- Added generic `Manage` flow for list-based resources:
  create, edit existing, delete, set default, cancel.
- Added `FormatMenuLabel` helper for consistent two-column CLI menu alignment.
- `FormatMenuLabel` now supports empty tag values for aligned untagged rows (e.g. `Exit`).
- Added `Colorize` helper for consistent ANSI color wrapping without duplicated reset logic.
- Added reusable user-facing error helpers:
  `NewUserError`, `WithHint`, `ErrorHint`, `FormatCLIError`.

Notes:

- Core remains provider-agnostic (no Talos/VMware domain logic).

v0.1.2 (2026-03-01)
Highlights:

- Fixed Shields badges owner casing in README to avoid intermittent `repo not found`.

Notes:

- Documentation-only patch release; no API/runtime changes.

v0.1.1 (2026-03-01)
Highlights:

- Added repository hardening docs: `SECURITY.md` and `CONTRIBUTING.md`.
- Standardized release notes format with explicit release date entries.

Notes:

- No API breaking changes.
- No functional runtime changes in wizard session/step behavior.

v0.1.0 (2026-03-01)
Highlights:

- Initial `Session` core for wizard draft lifecycle orchestration.
- Initial `Step` runner (`RunSteps`) with start/done callbacks.
- CI workflow (tests, vet, lint, vulncheck).
- Release workflow + repository badges.

Notes:

- CI lint configuration schema compatibility (`.golangci.yml` version set to `2`).
