# Product design

## Goal

- Build a personal, minimal daily editor containing only the required features.
- Use the project to practice agentic development.
- Let an LLM implement small tasks under human review.
- Maintain the project as a public fork of Micro.
- Produce a single locally built binary.
- There is no fixed deadline; the first subjective readiness criterion is that
  the editor builds and feels pleasant to use.

## Base editing experience

- Modeless editing with familiar Ctrl shortcuts.
- Full mouse interaction.
- Multiple cursors.
- Tabs that can be opened and closed, switched by mouse and keyboard, and
  restored after restart.
- Manual saving with `Ctrl+S`; no automatic saving in the first version.
- Command palette on `Ctrl+Shift+P`.
- Flat, minimal visual style without unnecessary borders.
- Preserve compatibility with Micro themes, settings, keybindings, and plugins.

## Project model

- The project root is the directory passed when launching the editor.
- One persistent left sidebar is shared across all tabs.
- Compact tabs at the top of the sidebar switch between Explorer, Search, and
  Git.
- The editor does not include an integrated terminal; shell processes remain in
  Zellij.

## Explorer

- Mouse-driven collapsible and filterable file tree.
- Display Git status for files.
- Create, rename, move, and delete through a mouse context menu.
- Delete moves items to the system trash after confirmation.
- Operations on open files must atomically update the buffer path, tab, session
  state, and LSP URI.

## LSP

- Generic Language Server Protocol client.
- Users install language servers themselves.
- The editor starts supported servers found in `PATH`.
- A built-in registry maps file types to candidate executables, arguments,
  root markers, and language IDs.
- User configuration can override the registry.
- First-version features:
  - automatic completion plus manual `Ctrl+Space` completion;
  - diagnostics;
  - go to definition;
  - find references.
- Diagnostics appear as gutter markers with a message on hover.
- Rename, code actions, formatting, and signature help are deferred.

## Git

- File list on the left and side-by-side diff on the right.
- Side-by-side diff is implemented before unified diff.
- Two groups:
  - Changes compares index with working tree;
  - Staged compares HEAD with index.
- Stage and unstage whole files only.
- Hunk and line staging are deferred.
- In a narrow Zellij pane, hide the sidebar and preserve the two diff columns.

## Project search

- Open with `Ctrl+Shift+F`.
- Run asynchronous searches through `ripgrep`.
- Support literal text, whole-word, case-sensitive, and regular-expression
  modes.
- Group results by file and open matches with the mouse.
- Implement search first.
- Add previewed project replace later, allowing individual matches to be
  selected before replacement.

## Sessions and recovery

- Store project state under XDG state using an identifier derived from the
  project root.
- Restore open tabs, active file, and unsaved changes.
- `Ctrl+Q` atomically saves a recovery snapshot and exits without save/discard
  prompts.
- Recovery must never silently overwrite disk content.
- If another process changes a file:
  - reload a clean buffer automatically;
  - show a version conflict for a modified buffer.

## Engineering guarantees

- Never lose saved or recovery text.
- Never apply a Git operation to the wrong path.
- Do not crash on Unicode input or large pastes.
- Use unit tests for editor and workbench state.
- Use property-based tests for text, Unicode, and coordinate conversions.
- Use exact paths and `--` for Git mutations.
- Keep LLM-generated changes small and reviewable.

## Explicit first-version non-goals

- Integrated terminal.
- Runtime LLM chat or agents inside the editor.
- Automatic installation of language servers.
- Hunk or line staging.
- Unified diff.
- Project-wide replace.
- Debugger or test runner.
- Full VS Code feature parity.
- Windows or macOS packaging.

