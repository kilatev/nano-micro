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
- M1.2 sidebar prototype passed: `workbench` toggles a read-only scratch
  sidebar on the left of its source pane, restores file-pane focus after open
  and mouse selection, and renders static Explorer, Search, and Git modes.
- The sidebar remains a normal Micro split: switching tabs retains it only in
  the tab where it was opened, and each new source pane can create another
  sidebar. Split width also follows Micro's normal split layout.
- M1.3 reproduced those limitations and recorded the dock decision in
  [`M1.3-dock-blocker.md`](decisions/M1.3-dock-blocker.md): tab-local and
  per-source-pane sidebars are architectural blockers; split width is polish.

## Current objective

M1.3 is complete. Plan the smallest generic application dock hook described in
[`M1.4-generic-dock-hook.md`](tasks/M1.4-generic-dock-hook.md).

## Blocker

None for M1.4 planning. Do not implement a core dock until its plan is
reviewed.

## Do not start yet

- Core dock implementation.
- Explorer.
- LSP.
- Git UI.
- Sessions and recovery.
- Workspace search.

## Recommended Codex CLI flow

1. Enter Plan mode with `/plan`.
2. Ask Codex to read `AGENTS.md` and the referenced documents.
3. Work only on `docs/tasks/M1.4-generic-dock-hook.md`.
4. Review the plan before allowing edits.
5. Run `/review` before accepting the completed milestone.
