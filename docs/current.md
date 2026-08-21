# Current state

Active milestone: M2.4 — Explorer create, rename, and move

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
- M1.4 added one optional application-owned left dock. It holds a scratch
  buffer outside tab split trees, persists across editor tabs, reserves its
  caller-supplied width, receives mouse input without taking editor focus, and
  restores the normal layout when detached. `workbench` now uses a 24-column
  dock; its simulation and pseudo-terminal checks passed.
- M1.5 binds `Ctrl+E` to toggle the bundled workbench. An explicit user
  `Ctrl-e` binding remains untouched; users can keep command mode or choose
  another shortcut in `bindings.json`.
- M2.1 is complete as a plugin-only Explorer: it captures the launch working
  directory, starts with that root expanded, sorts directories before files,
  hides dot entries by default, and opens clicked files through Micro's normal
  `open` command. The 24-column dock stays attached and does not take focus.
- M2.2 adds plugin-only Git badges on dock open and `workbench-refresh`. It
  parses NUL-delimited porcelain-v1 `XY` records, marks changed directories,
  and renders missing deleted paths as non-openable rows. Non-repositories keep
  the ordinary Explorer, and failed refreshes preserve the last valid tree.
- The workbench simulation covers staged, unstaged, untracked, renamed, and
  deleted paths, including spaces and non-ASCII names. Focused tests,
  `make test`, `make build`, and `./micro -version` pass.
- M2.3 adds an application-owned, Lua-exposed menu overlay. It clamps to the
  terminal, supports mouse plus Up/Down/Enter, closes on Escape, and resets
  captured mouse-release state. Explorer right-clicks use it without replacing
  explicit user bindings, provide entry-specific deferred actions, and reject
  a callback after the tree rerenders. The concrete plugin blocker is recorded
  in [`M2.3-menu-blocker.md`](decisions/M2.3-menu-blocker.md).

## Current objective

Implement M2.4: Explorer create, rename, and move. Keep deletion deferred to
M2.5 and open-buffer coordination deferred to M2.6.

## Blocker

None.

## Do not start yet

- Explorer work beyond M2.4.
- LSP.
- Git UI.
- Sessions and recovery.
- Workspace search.

## Recommended Codex CLI flow

1. Enter Plan mode with `/plan`.
2. Ask Codex to read `AGENTS.md` and the referenced documents.
3. Work only on `docs/tasks/M2.4-create-rename-move.md`.
4. Review the plan before allowing edits.
5. Run `/review` before accepting the completed milestone.
