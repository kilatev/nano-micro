# Backlog

The backlog uses rolling-wave planning. Only the active task is implementation
ready. Near-term tasks are more detailed; later tasks specify outcomes and
acceptance criteria without guessing at Micro internals that have not yet been
inspected.

## Status legend

- `ACTIVE`: the only task that may be implemented now.
- `COMPLETE`: verified and closed.
- `READY`: sufficiently specified, but blocked by an earlier task.
- `SKELETON`: outcome agreed; detailed plan must be written just in time.

## M0 — Bootstrap

- `COMPLETE` [`M0-bootstrap.md`](tasks/M0-bootstrap.md)

## M1 — Workbench shell

- `COMPLETE` [`M1.1-plugin-spike.md`](tasks/M1.1-plugin-spike.md)
- `COMPLETE` [`M1.2-sidebar-prototype.md`](tasks/M1.2-sidebar-prototype.md)
- `COMPLETE` [`M1.3-dock-blocker-decision.md`](tasks/M1.3-dock-blocker-decision.md)
- `COMPLETE` [`M1.4-generic-dock-hook.md`](tasks/M1.4-generic-dock-hook.md)
- `COMPLETE` [`M1.5-workbench-shortcut.md`](tasks/M1.5-workbench-shortcut.md)

## M2 — Explorer

- `COMPLETE` [`M2.1-read-only-tree.md`](tasks/M2.1-read-only-tree.md)
- `COMPLETE` [`M2.2-git-status-badges.md`](tasks/M2.2-git-status-badges.md)
- `COMPLETE` [`M2.3-context-menu.md`](tasks/M2.3-context-menu.md)
- `COMPLETE` [`M2.4-create-rename-move.md`](tasks/M2.4-create-rename-move.md)
- `COMPLETE` [`M2.5-trash-delete.md`](tasks/M2.5-trash-delete.md)
- `COMPLETE` [`M2.6-open-buffer-coordination.md`](tasks/M2.6-open-buffer-coordination.md)
- `COMPLETE` [`M2.7-move-destination-picker.md`](tasks/M2.7-move-destination-picker.md)
- `COMPLETE` [`M2.8-move-destination-fuzzy-search.md`](tasks/M2.8-move-destination-fuzzy-search.md)

## M3 — Sessions and recovery

- `SKELETON` [`M3.1-session-format.md`](tasks/M3.1-session-format.md)
- `SKELETON` [`M3.2-tab-restore.md`](tasks/M3.2-tab-restore.md)
- `SKELETON` [`M3.3-recovery-snapshots.md`](tasks/M3.3-recovery-snapshots.md)
- `SKELETON` [`M3.4-external-file-conflicts.md`](tasks/M3.4-external-file-conflicts.md)

## M4 — Minimal LSP

- `SKELETON` [`M4.1-lsp-transport.md`](tasks/M4.1-lsp-transport.md)
- `SKELETON` [`M4.2-server-registry.md`](tasks/M4.2-server-registry.md)
- `SKELETON` [`M4.3-document-sync.md`](tasks/M4.3-document-sync.md)
- `SKELETON` [`M4.4-completion.md`](tasks/M4.4-completion.md)
- `SKELETON` [`M4.5-diagnostics.md`](tasks/M4.5-diagnostics.md)
- `SKELETON` [`M4.6-definition-references.md`](tasks/M4.6-definition-references.md)

## M5 — Git review

- `SKELETON` [`M5.1-git-status-model.md`](tasks/M5.1-git-status-model.md)
- `SKELETON` [`M5.2-side-by-side-rendering.md`](tasks/M5.2-side-by-side-rendering.md)
- `SKELETON` [`M5.3-stage-unstage-file.md`](tasks/M5.3-stage-unstage-file.md)

## M6 — Workspace search

- `SKELETON` [`M6.1-ripgrep-runner.md`](tasks/M6.1-ripgrep-runner.md)
- `SKELETON` [`M6.2-search-sidebar.md`](tasks/M6.2-search-sidebar.md)
- `SKELETON` [`M6.3-search-controls.md`](tasks/M6.3-search-controls.md)

## M7 — Previewed project replace

- `SKELETON` [`M7-previewed-replace.md`](tasks/M7-previewed-replace.md)

## Promotion rule

At the end of each task:

1. Update `docs/current.md` with verified facts and blockers.
2. Refine only the next blocked task using those facts.
3. Change exactly one task to `ACTIVE`.
4. Do not expand later skeletons unless a new dependency is discovered.
