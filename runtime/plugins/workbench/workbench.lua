local buffer = import("micro/buffer")
local config = import("micro/config")
local micro = import("micro")
local os = import("os")
local filepath = import("path/filepath")

local tool = nil
local projectRoot = nil
local expanded = {}
local rows = {}

local modes = { "Explorer", "Search", "Git" }

function preinit()
    config.RegisterGlobalOption("workbench", "showhidden", false)
end

local function add(path, name, depth, isDir, lines)
    rows[#rows + 1] = { path = path, isDir = isDir }
    lines[#lines + 1] = string.rep("  ", depth) .. (isDir and (expanded[path] and "- " or "+ ") or "  ") .. name
    if not isDir or not expanded[path] then
        return
    end

    local entries, err = os.ReadDir(path)
    if err ~= nil then
        micro.InfoBar():Error("Could not read " .. path .. ": " .. tostring(err))
        return
    end
    local sorted = {}
    for i = 1, #entries do
        sorted[i] = entries[i]
    end
    table.sort(sorted, function(a, b)
        if a:IsDir() ~= b:IsDir() then
            return a:IsDir()
        end
        return a:Name() < b:Name()
    end)
    for _, entry in ipairs(sorted) do
        if config.GetGlobalOption("workbench.showhidden") or string.sub(entry:Name(), 1, 1) ~= "." then
            add(filepath.Join(path, entry:Name()), entry:Name(), depth + 1, entry:IsDir(), lines)
        end
    end
end

local function render(sidebar, selected)
    local tool = sidebar.Buf
    local lines = {}
    for i, mode in ipairs(modes) do
        lines[i] = (i == selected and "> " or "  ") .. mode
    end
    rows = {}
    if selected == 1 then
        add(projectRoot, projectRoot, 0, true, lines)
    end

    tool.Type.Readonly = false
    tool:Remove(buffer.Loc(0, 0), tool:End())
    tool:Insert(buffer.Loc(0, 0), table.concat(lines, "\n") .. "\n")
    tool.Type.Readonly = true
end

function init()
    projectRoot = os.Getwd()
    expanded[projectRoot] = true
    config.MakeCommand("workbench", open, config.NoComplete)
    config.TryBindKey("Ctrl-e", "command:workbench", false)
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
            return
        end
        local row = rows[selected - #modes]
        if row == nil then
            return
        end
        if row.isDir then
            expanded[row.path] = not expanded[row.path]
            render(bp, 1)
        else
            micro.CurPane():HandleCommand("open '" .. string.gsub(row.path, "'", "'\\''") .. "'")
        end
    end
end
