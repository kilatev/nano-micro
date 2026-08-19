# Roadmap

## M0 — Bootstrap

Create a reproducible upstream baseline: clean fork, build, tests, project
instructions, and manual edit/save/reopen verification.

## M1 — Workbench shell

Add a bundled workbench plugin with persistent sidebar scaffolding and compact
Explorer, Search, and Git tabs. Begin plugin-only. Add a generic dock core hook
only if the prototype proves it necessary.

## M2 — Explorer

Add the file tree, filtering, Git badges, context menu, safe create/rename/move,
and trash-based deletion. Coordinate filesystem operations with open buffers,
tabs, sessions, and LSP URIs.

## M3 — Sessions and recovery

Add XDG project sessions, restored tabs, unsaved-buffer snapshots, atomic
`Ctrl+Q` recovery, and external-file conflict handling.

## M4 — Minimal LSP

Add the server registry, PATH discovery, lifecycle management, completion,
diagnostics, definition, references, and tested UTF-8/UTF-16 conversions.

## M5 — Git review

Add Changes and Staged groups, side-by-side diff, narrow-pane behavior, and
whole-file stage/unstage with exact-path safeguards.

## M6 — Workspace search

Add asynchronous `rg --json` search, file-grouped results, mouse navigation,
and literal/word/case/regex controls.

## M7 — Previewed project replace

Add stale-result detection, selectable replacement previews, transactional
backups, and safe application across open and unopened files.

