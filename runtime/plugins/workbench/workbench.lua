local buffer = import("micro/buffer")
local config = import("micro/config")
local micro = import("micro")
local os = import("os")
local filepath = import("path/filepath")
local shell = import("micro/shell")

local tool = nil
local sidebar = nil
local projectRoot = nil
local expanded = {}
local rows = {}
local statuses = {}
local changedDirs = {}
local virtualChildren = {}
local selectedMode = 1

local modes = { "Explorer", "Search", "Git" }

function preinit()
    config.RegisterGlobalOption("workbench", "showhidden", false)
end

local function add(path, name, depth, isDir, isVirtual, lines)
    rows[#rows + 1] = { path = path, isDir = isDir, isVirtual = isVirtual }
    local badge = isDir and changedDirs[path] and " [*]" or (not isDir and statuses[path] and " [" .. statuses[path] .. "]" or "")
    lines[#lines + 1] = string.rep("  ", depth) .. (isDir and (expanded[path] and "- " or "+ ") or "  ") .. name .. badge
    if not isDir or not expanded[path] then
        return
    end

    local sorted = {}
    local seen = {}
    if not isVirtual then
        local entries, err = os.ReadDir(path)
        if err ~= nil then
            micro.InfoBar():Error("Could not read " .. path .. ": " .. tostring(err))
            return
        end
        for i = 1, #entries do
            local entry = entries[i]
            local child = filepath.Join(path, entry:Name())
            sorted[#sorted + 1] = { path = child, name = entry:Name(), isDir = entry:IsDir(), isVirtual = false }
            seen[child] = true
        end
    end
    for child, entry in pairs(virtualChildren[path] or {}) do
        if not seen[child] then
            sorted[#sorted + 1] = entry
        end
    end
    table.sort(sorted, function(a, b)
        if a.isDir ~= b.isDir then
            return a.isDir
        end
        return a.name < b.name
    end)
    for _, entry in ipairs(sorted) do
        if config.GetGlobalOption("workbench.showhidden") or string.sub(entry.name, 1, 1) ~= "." then
            add(entry.path, entry.name, depth + 1, entry.isDir, entry.isVirtual, lines)
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
        add(projectRoot, projectRoot, 0, true, false, lines)
    end

    tool.Type.Readonly = false
    tool:Remove(buffer.Loc(0, 0), tool:End())
    tool:Insert(buffer.Loc(0, 0), table.concat(lines, "\n") .. "\n")
    tool.Type.Readonly = true
end

local function parseStatus(output, repoRoot)
    local parsed = {}
    local start = 1
    while start <= #output do
        local finish = string.find(output, "\0", start, true)
        if finish == nil then
            return nil, "unterminated Git status record"
        end
        local record = string.sub(output, start, finish - 1)
        start = finish + 1
        if #record < 4 or string.sub(record, 3, 3) ~= " " then
            return nil, "malformed Git status record"
        end
        local status = string.sub(record, 1, 2)
        parsed[filepath.Clean(filepath.Join(repoRoot, string.sub(record, 4)))] = status
        local kind = string.sub(status, 1, 1)
        if kind == "R" or kind == "C" then
            finish = string.find(output, "\0", start, true)
            if finish == nil then
                return nil, "unterminated Git rename record"
            end
            start = finish + 1
        end
    end
    return parsed, nil
end

local function decorations(parsed)
    local dirs = {}
    local children = {}
    for path, status in pairs(parsed) do
        local parent = filepath.Dir(path)
        while parent ~= filepath.Dir(parent) do
            dirs[parent] = true
            if parent == projectRoot then
                break
            end
            parent = filepath.Dir(parent)
        end

        local _, err = os.Stat(path)
        if err ~= nil and string.find(status, "D", 1, true) ~= nil then
            local child = path
            parent = filepath.Dir(child)
            children[parent] = children[parent] or {}
            children[parent][child] = { path = child, name = filepath.Base(child), isDir = false, isVirtual = true }
            while parent ~= projectRoot do
                local _, parentErr = os.Stat(parent)
                if parentErr == nil then
                    break
                end
                child = parent
                parent = filepath.Dir(child)
                children[parent] = children[parent] or {}
                children[parent][child] = { path = child, name = filepath.Base(child), isDir = true, isVirtual = true }
            end
        end
    end
    return dirs, children
end

local function refreshGit()
    local rootOutput, err = shell.ExecCommand("git", "-C", projectRoot, "rev-parse", "--show-toplevel")
    if err ~= nil then
        if string.find(rootOutput, "not a git repository", 1, true) ~= nil then
            statuses = {}
            changedDirs = {}
            virtualChildren = {}
            if sidebar ~= nil then
                render(sidebar, selectedMode)
            end
            return
        end
        micro.InfoBar():Error("Could not refresh Git status: " .. (rootOutput ~= "" and rootOutput or tostring(err)))
        return
    end

    local repoRoot = string.gsub(rootOutput, "[\r\n]+$", "")
    local output, statusErr = shell.ExecCommand("git", "-C", projectRoot, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--", ".")
    if statusErr ~= nil then
        micro.InfoBar():Error("Could not refresh Git status: " .. (output ~= "" and output or tostring(statusErr)))
        return
    end
    local parsed, parseErr = parseStatus(output, repoRoot)
    if parseErr ~= nil then
        micro.InfoBar():Error("Could not refresh Git status: " .. parseErr)
        return
    end

    local dirs, children = decorations(parsed)
    statuses = parsed
    changedDirs = dirs
    virtualChildren = children
    if sidebar ~= nil then
        render(sidebar, selectedMode)
    end
end

function init()
    projectRoot = os.Getwd()
    expanded[projectRoot] = true
    config.MakeCommand("workbench", open, config.NoComplete)
    config.MakeCommand("workbench-refresh", refreshGit, config.NoComplete)
    config.TryBindKey("Ctrl-e", "command:workbench", false)
end

function open(bp, args)
    if tool ~= nil then
        micro.ClearDock()
        tool = nil
        sidebar = nil
        return
    end

    selectedMode = 1
    tool = buffer.NewBuffer("", "")
    tool:SetName("Workbench")
    tool.Type.Scratch = true
    sidebar = micro.SetDockBuffer(tool, 24)
    render(sidebar, selectedMode)
    refreshGit()
end

function onMousePress(bp, event)
    if bp.Buf == tool then
        local selected = bp.Cursor.Loc.Y + 1
        if selected >= 1 and selected <= #modes then
            selectedMode = selected
            render(bp, selectedMode)
            return
        end
        local row = rows[selected - #modes]
        if row == nil then
            return
        end
        if row.isDir then
            expanded[row.path] = not expanded[row.path]
            render(bp, 1)
        elseif not row.isVirtual then
            micro.CurPane():HandleCommand("open '" .. string.gsub(row.path, "'", "'\\''") .. "'")
        end
    end
end
