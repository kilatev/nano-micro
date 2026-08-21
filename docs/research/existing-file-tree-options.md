# Existing file-tree reuse options

Research date: 2026-08-21

## Recommendation

Keep the small Explorer owned by the bundled `workbench` plugin. No existing
project is a drop-in component for Micro's application-owned dock.

The closest reusable implementation is
[filemanager2](https://github.com/Neko-Box-Coder/filemanager2), an MIT-licensed,
currently maintained Micro Lua plugin. It is useful as a UX and implementation
reference, but adopting it as our base would increase rather than reduce the
code we must own: its main Lua file is about 1,500 source lines, owns a normal
Micro split/tab, and its file mutation and Git behavior do not satisfy this
project's requirements. NERDTree and the Neovim trees are coupled to their host
editors and would require ports, not integration.

The pre-M2.2 workbench implementation already crossed the project-specific
seam: it rendered a recursive tree into the dock, handled mouse selection, and
opened files in the active editor pane in roughly 100 lines
([source](../../runtime/plugins/workbench/workbench.lua)). The core dock API
accepts a Micro buffer and creates a `BufPane`
([source](../../internal/action/tab.go#L128)); that ownership model is the main
reason existing split-based plugins and terminal applications cannot simply be
placed in the left pane.

## Candidate comparison

| Candidate | What it provides | Runtime and UI fit | License and maintenance | Practical reuse |
|---|---|---|---|---|
| [Micro filemanager2](https://github.com/Neko-Box-Coder/filemanager2) | Tree navigation, mouse opening, hidden and Git-ignored filtering, persistent tree across tabs, create/rename/delete/copy ([README](https://github.com/Neko-Box-Coder/filemanager2#readme)) | Correct Micro Lua runtime, but it creates/manages ordinary splits and tabs; it does not use `SetDockBuffer`. Git support only asks Git for ignored paths. Delete ultimately calls `os.RemoveAll`; rename is a direct `os.Rename`; copy shells out to `cp`/`xcopy` ([source](https://github.com/Neko-Box-Coder/filemanager2/blob/main/filemanager2.lua)). It has no modified/staged/untracked/renamed/deleted badges, system trash, or atomic coordination with open buffers. | MIT ([license](https://github.com/Neko-Box-Coder/filemanager2/blob/main/LICENSE)); active releases include [v1.5.2](https://github.com/Neko-Box-Coder/filemanager2/releases/tag/v1.5.2). | Not drop-in. A fork is technically possible, but removing its split lifecycle and replacing its mutation/Git layers would leave us maintaining a large fork. Borrow interaction ideas or small reviewed MIT-licensed routines only. |
| [Official-channel filemanager](https://github.com/micro-editor/updated-plugins/tree/master/filemanager-plugin) | Earlier Micro tree with mouse navigation, hidden/Git-ignored filtering, and basic create/rename/delete operations; it remains listed in Micro's [official plugin channel](https://github.com/micro-editor/plugin-channel). | Micro Lua, but also owns a vertical split and opens files through split-oriented behavior. Its source has the same direct `RemoveAll`/`Rename` model and ignored-only Git query ([source](https://github.com/micro-editor/updated-plugins/blob/master/filemanager-plugin/filemanager.lua)). | MIT in the [original repository](https://github.com/NicolaiSoeborg/filemanager-plugin); the original changelog's latest release is [3.4.0 from 2018](https://github.com/NicolaiSoeborg/filemanager-plugin/blob/master/CHANGELOG.md), while Micro maintains an updated-plugin fork. | Not drop-in. filemanager2 is the better reference because it contains the later fixes and features. The even older official [filetree prototype](https://github.com/micro-editor/filetree-plugin) is only a small WIP historical implementation. |
| [NERDTree](https://github.com/preservim/nerdtree) | Mature hierarchical tree, mouse support, filtering, bookmarks, and a programmable file-system menu ([manual](https://github.com/preservim/nerdtree/blob/master/doc/NERDTree.txt)). Git decorations are supplied by a separate [NERDTree Git plugin](https://github.com/Xuyuanp/nerdtree-git-plugin). | Vimscript tightly coupled to Vim buffers, window variables, mappings, commands, and `wincmd` ([source](https://github.com/preservim/nerdtree/blob/master/autoload/nerdtree/ui_glue.vim)). Micro cannot load it. | WTFPL ([license](https://github.com/preservim/nerdtree/blob/master/LICENCE)). | Ideas only. A Micro version would be a rewrite, so calling it reuse would obscure the real maintenance cost. |
| [nvim-tree.lua](https://github.com/nvim-tree/nvim-tree.lua) (and similar Neo-tree) | Rich tree, live filtering, file operations, automatic refresh, Git integration, and distinct unstaged/staged/unmerged/renamed/untracked/deleted/ignored states ([documentation](https://github.com/nvim-tree/nvim-tree.lua/blob/master/doc/nvim-tree-lua.txt)). | Despite being Lua, it requires Neovim and its `vim.api`, buffers/windows/tabpages, scheduler, and libuv APIs ([view source](https://github.com/nvim-tree/nvim-tree.lua/blob/master/lua/nvim-tree/view.lua), [core source](https://github.com/nvim-tree/nvim-tree.lua/blob/master/lua/nvim-tree/core.lua)). Micro's embedded Lua API is unrelated. | GPL-3.0 ([license](https://github.com/nvim-tree/nvim-tree.lua/blob/master/LICENSE)); the project describes itself as stable and focused on maintenance ([README](https://github.com/nvim-tree/nvim-tree.lua#roadmap)). | Not embeddable, and copying substantial code would also introduce GPL licensing constraints into this MIT project. Use its behavior as a checklist, not its implementation. |
| External TUI or picker: [Yazi](https://github.com/sxyazi/yazi), [fzf](https://github.com/junegunn/fzf), or [micro-fzfinder](https://github.com/MuratovAS/micro-fzfinder) | Fast browsing or file selection with mature standalone tooling. | They run as external processes. Yazi has its own plugin runtime and terminal UI ([plugin docs](https://yazi-rs.github.io/docs/plugins/overview/)); fzf integrations return a selected path rather than maintain a dock tree. The current dock cannot host Micro's separate terminal-pane type. | Yazi and fzf are MIT licensed and actively released in their official repositories. | Good scope-reduction alternative if a persistent Explorer is dropped: launch a picker and open the selected file. Not a way to implement the specified left dock without first adding a generic terminal-dock core primitive and bundling/requiring another executable. |

## Consequences for the plan

The existing M2.1 tree was not unnecessary reinvention: it is the adapter
between filesystem data, Micro buffers, active-pane opening, and the new
application dock. Replacing it with filemanager2 would trade a small tailored
implementation for a much larger split manager whose risky portions we would
need to replace.

Proceed with M2.2's planned single `git status --porcelain=v1 -z` refresh and
path-to-status map
([active task](../tasks/M2.2-git-status-badges.md)). This is smaller and better
matched to the required badge distinctions than any Micro plugin found. Before
M2.3/M2.4, use filemanager2 and NERDTree as interaction references for prompts,
keyboard behavior, and menu vocabulary, but retain the planned safety rules:
system trash, explicit confirmation, and coordinated open-buffer/session/LSP
updates. Do not reuse their direct permanent-delete or rename implementations.

If product scope changes and a persistent left tree is no longer important, an
fzf-style picker is the only option here that materially removes code rather
than moving complexity into a port or long-lived fork.
