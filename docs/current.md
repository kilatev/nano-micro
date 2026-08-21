# Current state

Active milestone: none

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
- M2.4 adds plugin-only Explorer create, rename, and move actions. A shared
  helper cleans each destination and rejects empty, project-root, and
  out-of-project targets; exclusive file creation and preflight destination
  checks prevent silent overwrite. Files and folders use standard `0644` and
  `0755` permissions. Move accepts either a final path or an existing target
  folder, retaining the selected basename in the latter case. Failed
  operations preserve the existing tree, while successful operations re-render
  it. The simulated workbench covers traversal, collisions, permissions,
  successful rename, and cross-directory move.
- M2.5 adds plugin-only trash deletion. Explorer displays the exact target in
  an Enter-to-confirm, Escape-to-cancel prompt, then calls `gio trash --` with
  argument-separated paths. Successful calls re-render the tree; a simulated
  `gio` handoff verifies the exact invocation without modifying the real trash.
- M2.6 adds `buffer.Rename(old, new)`, the documented core blocker for Lua
  Explorer coordination. It rejects a destination owned by a different open
  shared buffer, renames on disk, then updates the shared path, backup
  ownership, modification time, and path settings. Explorer uses it for both
  rename and move and refreshes tab labels. State tests cover Explorer
  rename/save, shared open views, conflicting destinations, and filesystem
  failure; focused tests, `make test`, `make build`, and `./micro -version`
  pass.
- M2.7 replaces the Move path prompt with a project-folder menu. It lists
  project-relative folders, excludes a moved folder and its descendants, and
  rechecks stale selections, destination existence, and overwrite protection
  before using `buffer.Rename`. Simulated tests cover nested selection, stale
  callbacks, and collisions; manual Move verification passed.
- M2.8 adds a case-insensitive fuzzy subsequence filter before the Move
  destination menu. It retains M2.7's destination checks and defers menu
  display through the plugin event hook to avoid re-entrant Lua callbacks.
  Simulated fuzzy-match and no-match tests pass; manual verification passed.

## Current objective

M2 Explorer work is complete. Refine and activate M3.1 before beginning
sessions or recovery work.

## Blocker

M3.1 is still a skeleton task and must be refined before activation.

## Do not start yet

- M3 work until M3.1 is refined and activated.
- LSP.
- Git UI.
- Sessions and recovery.
- Workspace search.

## Recommended Codex CLI flow

1. Enter Plan mode with `/plan`.
2. Ask Codex to read `AGENTS.md` and the referenced documents.
3. Work only on `docs/tasks/M2.6-open-buffer-coordination.md`.
4. Review the plan before allowing edits.
5. Run `/review` before accepting the completed milestone.
