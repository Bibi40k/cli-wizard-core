Release Notes

Unreleased
Highlights:

TBD

Notes:

TBD

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
