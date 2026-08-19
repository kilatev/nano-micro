local buffer = import("micro/buffer")
local config = import("micro/config")

local tool

function init()
    config.MakeCommand("workbench", open, config.NoComplete)
end

function open(bp, args)
    -- ponytail: each invocation opens a split; M1.2 owns sidebar toggling.
    tool = buffer.NewBuffer("Explorer\nSearch\nGit\n", "")
    tool:SetName("Workbench")
    tool.Type.Readonly = true
    tool.Type.Scratch = true
    bp:VSplitBuf(tool)
end

function onMousePress(bp, event)
    if bp.Buf == tool then
        tool:SetName("Workbench (selected)")
    end
end
