# Current state

Active milestone: M1

## Completed

- Product grilling and design tree.
- Choice of Micro as the editor core.
- Agreement on plugin-first architecture and safety boundaries.
- Outcome-oriented task inventory for M0-M7.
- `make build` and `make test` pass with Go 1.26.5 on Linux.
- The built binary reports `0.0.0-unknown` and commit `e3facebf`; no reachable
  release tag is present in this checkout.
- Manual edit/save/reopen verification passed in a real terminal using a new
  temporary file containing `M0 reopen check ✓`.
- M1.1 workbench plugin spike passed: `workbench` opens a read-only scratch
  buffer in a vertical split, preserving the file buffer; its mouse handler
  was verified in both the simulated screen and a real terminal.

## Current objective

M1.1 is complete. Do not begin M1.2 until its plan is reviewed.

## Blocker

None for M1.1.

## Do not start yet

- Sidebar implementation.
- Explorer.
- LSP.
- Git UI.
- Sessions and recovery.
- Workspace search.

## Recommended Codex CLI flow

1. Enter Plan mode with `/plan`.
2. Ask Codex to read `AGENTS.md` and the referenced documents.
3. Work only on `docs/tasks/M0-bootstrap.md`.
4. Review the plan before allowing edits.
5. Run `/review` before accepting the completed milestone.
