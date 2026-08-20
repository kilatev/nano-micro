local buffer = import("micro/buffer")
local config = import("micro/config")
local micro = import("micro")

local tool = nil

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
    if tool ~= nil then
        micro.ClearDock()
        tool = nil
        return
    end

    tool = buffer.NewBuffer("", "")
    tool:SetName("Workbench")
    tool.Type.Scratch = true
    local dock = micro.SetDockBuffer(tool, 24)
    render(dock, 1)
end

function onMousePress(bp, event)
    if bp.Buf == tool then
        local selected = bp.Cursor.Loc.Y + 1
        if selected >= 1 and selected <= #modes then
            render(bp, selected)
        end
    end
end
