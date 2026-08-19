# Project instructions

Before planning or editing, read:

- `docs/product.md`
- `docs/architecture.md`
- `docs/roadmap.md`
- `docs/backlog.md`
- `docs/current.md`
- the active task under `docs/tasks/`

## Product

This is a public fork of Micro for a personal, minimal, modeless terminal code
editor with a VS Code-like workbench.

## Architecture rules

- Preserve compatibility with upstream Micro.
- Keep the Micro editor core stable.
- Implement workbench features as bundled Lua plugins first.
- Modify Go core only after documenting a concrete plugin API blocker.
- Core additions must be generic UI or lifecycle primitives, not
  feature-specific Git, LSP, Explorer, or Search implementations.
- Preserve existing Micro settings, themes, keybindings, and plugins.

## Agent workflow

- Work on one small, observable outcome at a time.
- Implement only the active task named in `docs/current.md`.
- Inspect relevant code before proposing changes.
- Present a plan before implementing non-trivial changes.
- Do not begin the next milestone automatically.
- Avoid unrelated refactoring.
- Do not add dependencies without explaining why.
- Add or update tests for changed behavior.
- Run formatting, focused tests, and the relevant full test suite.
- Finish with:
  - files changed;
  - tests run;
  - manual verification;
  - remaining risks.
- Do not commit unless explicitly requested.

## Safety invariants

- Never lose saved or recovery text.
- Never run a Git mutation against an ambiguous or unintended path.
- Text operations must handle Unicode correctly.
- Recovery data must never silently overwrite disk content.
