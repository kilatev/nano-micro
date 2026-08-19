# Architecture

## Base

- Public fork of Micro.
- Minimal upstream-compatible patch.
- Single locally built binary.
- Ubuntu and Zellij are the initial environment.

Micro remains responsible for the mature editor core: text buffers, undo,
mouse input, multiple cursors, tabs, splits, clipboard, Unicode rendering,
themes, configuration, and plugins.

## Workbench layer

Start with bundled Lua plugins:

- `workbench-shell`
- `explorer`
- `lsp`
- `git-review`
- `workspace-search`
- `session`

Prototype each capability through existing buffers, splits, mouse callbacks,
commands, and background processes before changing the Go core.

## Permitted core additions

Add a Go core primitive only after a plugin prototype documents a concrete
blocker. Prefer generic, upstreamable primitives:

1. Application-level dock/sidebar shared by tabs.
2. Shutdown/session lifecycle hook.
3. Generic popup, list, and context-menu primitives.
4. Optional filesystem-watching hook.

Feature-specific Git, LSP, Explorer, Search, or Session policy remains outside
the editor core.

## Known high-risk boundaries

- UTF-8 editor positions versus UTF-16 LSP positions.
- Late LSP responses for buffers that have changed.
- File rename or move while the file is open.
- Atomic recovery snapshots and disk conflicts.
- Side-by-side diff alignment, renamed files, binaries, untracked files, and
  repositories without HEAD.
- Stale project-search matches and multi-file replacement.

## Development strategy

- Build one thin vertical slice at a time.
- Preserve existing Micro behavior by default.
- Add new actions and options instead of changing existing semantics.
- Keep feature logic independently testable from terminal rendering.
- Require a documented blocker before expanding the Go core surface.

