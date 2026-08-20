local buffer = import("micro/buffer")
local config = import("micro/config")

-- ponytail: one sidebar per source pane; use an application dock if tabs must share one sidebar.
local sidebars = {}
local editors = {}

local modes = { "Explorer", "Search", "Git" }

local function render(sidebar, selected)
    local tool = sidebar.Buf
    local lines = {}
    for i, mode in ipairs(modes) do
        lines[i] = (i == selected and "> " or "  ") .. mode
    end

    tool.Type.Readonly = false
    tool:Remove(buffer.Loc(0, 0), tool:End())
    tool:Insert(buffer.Loc(0, 0), table.concat(lines, "\n") .. "\n")
    tool.Type.Readonly = true
end

function init()
    config.MakeCommand("workbench", open, config.NoComplete)
end

function open(bp, args)
    local tab = bp:Tab()
    local source = bp:ID()
    local sidebar = sidebars[source]
    if sidebar ~= nil then
        sidebar:ForceQuit()
        editors[sidebar:ID()] = nil
        sidebars[source] = nil
        return
    end

    local tool = buffer.NewBuffer("", "")
    tool:SetName("Workbench")
    tool.Type.Scratch = true
    sidebar = bp:VSplitIndex(tool, false)
    sidebars[source] = sidebar
    editors[sidebar:ID()] = bp
    render(sidebar, 1)
    tab:SetActive(tab:GetPane(bp:ID()))
end

function onMousePress(bp, event)
    local tab = bp:Tab()
    if bp.Buf.Type.Scratch and bp.Buf:GetName() == "Workbench" then
        local selected = bp.Cursor.Loc.Y + 1
        if selected >= 1 and selected <= #modes then
            render(bp, selected)
            local editor = editors[bp:ID()]
            tab:SetActive(tab:GetPane(editor:ID()))
        end
    end
end
