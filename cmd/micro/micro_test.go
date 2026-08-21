package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-errors/errors"
	"github.com/micro-editor/micro/v2/internal/action"
	"github.com/micro-editor/micro/v2/internal/buffer"
	"github.com/micro-editor/micro/v2/internal/config"
	"github.com/micro-editor/micro/v2/internal/screen"
	"github.com/micro-editor/micro/v2/internal/util"
	"github.com/micro-editor/tcell/v2"
	"github.com/stretchr/testify/assert"
)

var tempDir string
var projectDir string
var workingDir string
var sim tcell.SimulationScreen

func init() {
	screen.Events = make(chan tcell.Event, 8)
}

func startup(args []string) (tcell.SimulationScreen, error) {
	var err error

	tempDir, err = os.MkdirTemp("", "micro_test")
	if err != nil {
		return nil, err
	}
	err = config.InitConfigDir(tempDir)
	if err != nil {
		return nil, err
	}

	config.InitRuntimeFiles(true)
	config.InitPlugins()

	err = config.ReadSettings()
	if err != nil {
		return nil, err
	}
	err = config.InitGlobalSettings()
	if err != nil {
		return nil, err
	}

	s, err := screen.InitSimScreen()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := recover(); err != nil {
			screen.Screen.Fini()
			fmt.Println("Micro encountered an error:", err)
			// immediately backup all buffers with unsaved changes
			for _, b := range buffer.OpenBuffers {
				if b.Modified() {
					b.Backup()
				}
			}
			// Print the stack trace too
			log.Fatalf(errors.Wrap(err, 2).ErrorStack())
		}
	}()

	err = config.LoadAllPlugins()
	if err != nil {
		screen.TermMessage(err)
	}

	action.InitBindings()
	action.InitCommands()
	err = config.RunPluginFn("preinit")
	if err != nil {
		return nil, err
	}

	err = config.InitColorscheme()
	if err != nil {
		return nil, err
	}

	b := LoadInput(args)

	if len(b) == 0 {
		return nil, errors.New("No buffers opened")
	}

	action.InitTabs(b)
	action.InitGlobals()

	err = config.RunPluginFn("init")
	if err != nil {
		return nil, err
	}

	s.InjectResize()
	handleEvent()

	return s, nil
}

func cleanup() {
	if err := os.Chdir(workingDir); err != nil {
		log.Println(err)
	}
	os.RemoveAll(tempDir)
	os.RemoveAll(projectDir)
}

func handleEvent() {
	screen.Lock()
	e := screen.Screen.PollEvent()
	screen.Unlock()
	if e != nil {
		screen.Events <- e
	}

	for len(screen.DrawChan()) > 0 || len(screen.Events) > 0 {
		DoEvent()
	}
}

func injectKey(key tcell.Key, r rune, mod tcell.ModMask) {
	sim.InjectKey(key, r, mod)
	handleEvent()
}

func injectMouse(x, y int, buttons tcell.ButtonMask, mod tcell.ModMask) {
	sim.InjectMouse(x, y, buttons, mod)
	handleEvent()
}

func injectString(str string) {
	// the tcell simulation screen event channel can only handle
	// 10 events at once, so we need to divide up the key events
	// into chunks of 10 and handle the 10 events before sending
	// another chunk of events
	iters := len(str) / 10
	extra := len(str) % 10

	for i := 0; i < iters; i++ {
		s := i * 10
		e := i*10 + 10
		sim.InjectKeyBytes([]byte(str[s:e]))
		for i := 0; i < 10; i++ {
			handleEvent()
		}
	}

	sim.InjectKeyBytes([]byte(str[len(str)-extra:]))
	for i := 0; i < extra; i++ {
		handleEvent()
	}
}

func openFile(file string) {
	action.MainTab().CurPane().HandleCommand(fmt.Sprintf("open %s", file))
}

func runCommand(command string) {
	action.MainTab().CurPane().HandleCommand(command)
}

func findBuffer(file string) *buffer.Buffer {
	var buf *buffer.Buffer
	file = util.ResolvePath(file)
	for _, b := range buffer.OpenBuffers {
		if b.AbsPath == file {
			buf = b
		}
	}
	return buf
}

func createTestFile(t *testing.T, content string) string {
	f, err := os.CreateTemp(t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}

	return f.Name()
}

func runGit(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", projectDir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func TestMain(m *testing.M) {
	var err error
	workingDir, err = os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	projectDir, err = os.MkdirTemp("", "micro_workbench")
	if err != nil {
		log.Fatalln(err)
	}
	if err := os.Mkdir(filepath.Join(projectDir, "alpha"), 0755); err != nil {
		log.Fatalln(err)
	}
	if err := os.Mkdir(filepath.Join(projectDir, "zeta"), 0755); err != nil {
		log.Fatalln(err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "alpha", "nested.txt"), []byte("nested content"), 0644); err != nil {
		log.Fatalln(err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "a.txt"), []byte("root content"), 0644); err != nil {
		log.Fatalln(err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, ".hidden"), []byte("hidden content"), 0644); err != nil {
		log.Fatalln(err)
	}
	if err := os.Chdir(projectDir); err != nil {
		log.Fatalln(err)
	}
	sim, err = startup([]string{})
	if err != nil {
		log.Fatalln(err)
	}

	retval := m.Run()
	cleanup()

	os.Exit(retval)
}

func TestSimpleEdit(t *testing.T) {
	file := createTestFile(t, "base content")

	openFile(file)

	if findBuffer(file) == nil {
		t.Fatalf("Could not find buffer %s", file)
	}

	injectKey(tcell.KeyEnter, rune(tcell.KeyEnter), tcell.ModNone)
	injectKey(tcell.KeyUp, 0, tcell.ModNone)
	injectString("first line")

	// test both kinds of backspace
	for i := 0; i < len("ne"); i++ {
		injectKey(tcell.KeyBackspace, rune(tcell.KeyBackspace), tcell.ModNone)
	}
	for i := 0; i < len(" li"); i++ {
		injectKey(tcell.KeyBackspace2, rune(tcell.KeyBackspace2), tcell.ModNone)
	}
	injectString("foobar")

	injectKey(tcell.KeyCtrlS, rune(tcell.KeyCtrlS), tcell.ModCtrl)

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "firstfoobar\nbase content\n", string(data))
}

func TestWorkbench(t *testing.T) {
	file := createTestFile(t, "base content")
	openFile(file)
	original := findBuffer(file)
	tab := action.MainTab()
	source := tab.CurPane()
	panes := len(tab.Panes)
	root := tab.Node
	children := len(root.Children())
	originalView := source.GetView()
	width, _ := screen.Screen.Size()

	injectKey(tcell.KeyCtrlE, rune(tcell.KeyCtrlE), tcell.ModCtrl)

	assert.Same(t, source, tab.CurPane())
	assert.Len(t, tab.Panes, panes)
	assert.Same(t, root, tab.Node)
	assert.Len(t, root.Children(), children)
	toolPane := action.Tabs.Dock
	if toolPane == nil {
		t.Fatal("workbench dock was not opened")
	}
	tool := toolPane.Buf
	assert.Equal(t, "> Explorer\n  Search\n  Git\n- "+projectDir+"\n  + alpha\n  + zeta\n    a.txt\n", string(tool.Bytes()))
	assert.True(t, tool.Type.Readonly)
	assert.True(t, tool.Type.Scratch)
	assert.Same(t, original, findBuffer(file))
	assert.Equal(t, 0, toolPane.GetView().X)
	assert.Equal(t, 24, toolPane.GetView().Width)
	assert.Equal(t, 25, source.GetView().X)
	assert.Equal(t, width-25, source.GetView().Width)
	assert.Nil(t, action.SetDockBuffer(original, 24))
	assert.Same(t, toolPane, action.Tabs.Dock)

	view := toolPane.GetView()
	injectMouse(view.X, view.Y+1, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y+1, tcell.ButtonNone, tcell.ModNone)
	assert.Equal(t, "  Explorer\n> Search\n  Git\n", string(tool.Bytes()))
	assert.Same(t, source, tab.CurPane())

	injectMouse(view.X, view.Y+2, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y+2, tcell.ButtonNone, tcell.ModNone)
	assert.Equal(t, "  Explorer\n  Search\n> Git\n", string(tool.Bytes()))
	assert.Same(t, source, tab.CurPane())
	assert.Equal(t, "base content", string(original.Bytes()))
	dockContent := string(tool.Bytes())

	source.AddTab()
	assert.Len(t, tab.Panes, panes)
	secondTab := action.MainTab()
	assert.Len(t, secondTab.Panes, 1)
	assert.Same(t, toolPane, action.Tabs.Dock)
	assert.Same(t, secondTab, toolPane.Tab())
	assert.Equal(t, 0, toolPane.GetView().X)
	assert.Equal(t, 24, toolPane.GetView().Width)
	assert.Equal(t, 25, secondTab.CurPane().GetView().X)
	assert.Equal(t, dockContent, string(tool.Bytes()))
	secondTab.CurPane().ForceQuit()
	assert.Same(t, tab, action.MainTab())
	assert.Same(t, tab, toolPane.Tab())

	view = source.GetView()
	injectMouse(view.X, view.Y, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y, tcell.ButtonNone, tcell.ModNone)
	injectString("!")
	assert.Equal(t, "!base content", string(original.Bytes()))
	injectKey(tcell.KeyCtrlS, rune(tcell.KeyCtrlS), tcell.ModCtrl)

	injectKey(tcell.KeyCtrlE, rune(tcell.KeyCtrlE), tcell.ModCtrl)
	assert.Nil(t, action.Tabs.Dock)
	assert.Len(t, tab.Panes, panes)
	assert.Same(t, source, tab.CurPane())
	assert.Equal(t, originalView, source.GetView())
}

func TestWorkbenchExplorer(t *testing.T) {
	source := action.MainTab().CurPane()
	previous := source.Buf
	injectKey(tcell.KeyCtrlE, rune(tcell.KeyCtrlE), tcell.ModCtrl)
	dock := action.Tabs.Dock
	if dock == nil {
		t.Fatal("workbench dock was not opened")
	}

	assert.Equal(t, "> Explorer\n  Search\n  Git\n- "+projectDir+"\n  + alpha\n  + zeta\n    a.txt\n", string(dock.Buf.Bytes()))
	assert.NotContains(t, string(dock.Buf.Bytes()), ".hidden")
	runCommand("workbench-refresh")
	assert.Equal(t, "> Explorer\n  Search\n  Git\n- "+projectDir+"\n  + alpha\n  + zeta\n    a.txt\n", string(dock.Buf.Bytes()))

	view := dock.GetView()
	injectMouse(view.X, view.Y+4, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y+4, tcell.ButtonNone, tcell.ModNone)
	assert.Equal(t, "> Explorer\n  Search\n  Git\n- "+projectDir+"\n  - alpha\n      nested.txt\n  + zeta\n    a.txt\n", string(dock.Buf.Bytes()))
	assert.Same(t, source, action.MainTab().CurPane())

	injectMouse(view.X, view.Y+4, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y+4, tcell.ButtonNone, tcell.ModNone)
	assert.Equal(t, "> Explorer\n  Search\n  Git\n- "+projectDir+"\n  + alpha\n  + zeta\n    a.txt\n", string(dock.Buf.Bytes()))

	runCommand("set workbench.showhidden true")
	assert.True(t, config.GetGlobalOption("workbench.showhidden").(bool))
	injectMouse(view.X, view.Y+1, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y+1, tcell.ButtonNone, tcell.ModNone)
	injectMouse(view.X, view.Y, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y, tcell.ButtonNone, tcell.ModNone)
	assert.Contains(t, string(dock.Buf.Bytes()), "    .hidden\n")

	injectMouse(view.X, view.Y+4, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y+4, tcell.ButtonNone, tcell.ModNone)
	injectMouse(view.X, view.Y+5, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y+5, tcell.ButtonNone, tcell.ModNone)
	assert.NotNil(t, findBuffer(filepath.Join(projectDir, "alpha", "nested.txt")))
	assert.Same(t, dock, action.Tabs.Dock)
	assert.Same(t, source, action.MainTab().CurPane())
	source.OpenBuffer(previous)
	injectKey(tcell.KeyCtrlE, rune(tcell.KeyCtrlE), tcell.ModCtrl)
}

func TestWorkbenchGitStatusBadges(t *testing.T) {
	runCommand("set workbench.showhidden false")
	oldName := filepath.Join(projectDir, "rename old.txt")
	deleted := filepath.Join(projectDir, "alpha", "delete me.txt")
	if err := os.WriteFile(oldName, []byte("rename me"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(deleted, []byte("delete me"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, "init")
	runGit(t, "add", "--", ".")
	runGit(t, "-c", "user.name=Workbench Test", "-c", "user.email=workbench@example.invalid", "commit", "-m", "base")

	if err := os.WriteFile(filepath.Join(projectDir, "a.txt"), []byte("staged"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, "add", "--", "a.txt")
	runGit(t, "mv", "--", "rename old.txt", "zeta/renamed ü.txt")
	if err := os.Remove(deleted); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "alpha", "nested.txt"), []byte("unstaged"), 0644); err != nil {
		t.Fatal(err)
	}
	untracked := filepath.Join(projectDir, "alpha", "new ü file.txt")
	if err := os.WriteFile(untracked, []byte("untracked"), 0644); err != nil {
		t.Fatal(err)
	}

	source := action.MainTab().CurPane()
	previous := source.Buf
	injectKey(tcell.KeyCtrlE, rune(tcell.KeyCtrlE), tcell.ModCtrl)
	dock := action.Tabs.Dock
	if dock == nil {
		t.Fatal("workbench dock was not opened")
	}
	view := dock.GetView()
	if strings.Contains(string(dock.Buf.Bytes()), "  - alpha [*]") {
		injectMouse(view.X, view.Y+4, tcell.Button1, tcell.ModNone)
		injectMouse(view.X, view.Y+4, tcell.ButtonNone, tcell.ModNone)
	}
	assert.Equal(t, "> Explorer\n  Search\n  Git\n- "+projectDir+" [*]\n  + alpha [*]\n  + zeta [*]\n    a.txt [M ]\n", string(dock.Buf.Bytes()))

	injectMouse(view.X, view.Y+4, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y+4, tcell.ButtonNone, tcell.ModNone)
	assert.Equal(t, "> Explorer\n  Search\n  Git\n- "+projectDir+" [*]\n  - alpha [*]\n      delete me.txt [ D]\n      nested.txt [ M]\n      new ü file.txt [??]\n  + zeta [*]\n    a.txt [M ]\n", string(dock.Buf.Bytes()))
	injectMouse(view.X, view.Y+5, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y+5, tcell.ButtonNone, tcell.ModNone)
	assert.Nil(t, findBuffer(deleted))
	injectMouse(view.X, view.Y+6, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y+6, tcell.ButtonNone, tcell.ModNone)
	assert.NotNil(t, findBuffer(filepath.Join(projectDir, "alpha", "nested.txt")))
	assert.Same(t, source, action.MainTab().CurPane())

	injectMouse(view.X, view.Y+4, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y+4, tcell.ButtonNone, tcell.ModNone)
	injectMouse(view.X, view.Y+5, tcell.Button1, tcell.ModNone)
	injectMouse(view.X, view.Y+5, tcell.ButtonNone, tcell.ModNone)
	assert.Contains(t, string(dock.Buf.Bytes()), "      renamed ü.txt [R ]\n")
	assert.NotContains(t, string(dock.Buf.Bytes()), "rename old.txt")

	afterRefresh := filepath.Join(projectDir, "zeta", "after refresh.txt")
	if err := os.WriteFile(afterRefresh, []byte("refresh"), 0644); err != nil {
		t.Fatal(err)
	}
	assert.NotContains(t, string(dock.Buf.Bytes()), "after refresh.txt")
	runCommand("workbench-refresh")
	assert.Contains(t, string(dock.Buf.Bytes()), "      after refresh.txt [??]\n")
	lastTree := string(dock.Buf.Bytes())
	t.Setenv("PATH", t.TempDir())
	runCommand("workbench-refresh")
	assert.Equal(t, lastTree, string(dock.Buf.Bytes()))

	source.OpenBuffer(previous)
	injectKey(tcell.KeyCtrlE, rune(tcell.KeyCtrlE), tcell.ModCtrl)
}

func TestMouse(t *testing.T) {
	file := createTestFile(t, "base content")

	openFile(file)

	if findBuffer(file) == nil {
		t.Fatalf("Could not find buffer %s", file)
	}

	// buffer:
	// base content
	// the selections need to happen at different locations to avoid a double click
	injectMouse(3, 0, tcell.Button1, tcell.ModNone)
	injectKey(tcell.KeyLeft, 0, tcell.ModNone)
	injectMouse(0, 0, tcell.ButtonNone, tcell.ModNone)
	injectString("secondline")
	// buffer:
	// secondlinebase content
	injectKey(tcell.KeyEnter, rune(tcell.KeyEnter), tcell.ModNone)
	// buffer:
	// secondline
	// base content
	injectMouse(2, 0, tcell.Button1, tcell.ModNone)
	injectMouse(0, 0, tcell.ButtonNone, tcell.ModNone)
	injectKey(tcell.KeyEnter, rune(tcell.KeyEnter), tcell.ModNone)
	// buffer:
	//
	// secondline
	// base content
	injectKey(tcell.KeyUp, 0, tcell.ModNone)
	injectString("firstline")
	// buffer:
	// firstline
	// secondline
	// base content
	injectKey(tcell.KeyCtrlS, rune(tcell.KeyCtrlS), tcell.ModCtrl)

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "firstline\nsecondline\nbase content\n", string(data))
}

var srTestStart = `foo
foo
foofoofoo
Ernleȝe foo æðelen
`
var srTest2 = `test_string
test_string
test_stringtest_stringtest_string
Ernleȝe test_string æðelen
`
var srTest3 = `test_foo
test_string
test_footest_stringtest_foo
Ernleȝe test_string æðelen
`

func TestSearchAndReplace(t *testing.T) {
	file := createTestFile(t, srTestStart)

	openFile(file)

	if findBuffer(file) == nil {
		t.Fatalf("Could not find buffer %s", file)
	}

	action.MainTab().CurPane().HandleCommand(fmt.Sprintf("replaceall %s %s", "foo", "test_string"))

	injectKey(tcell.KeyCtrlS, rune(tcell.KeyCtrlS), tcell.ModCtrl)

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, srTest2, string(data))

	action.MainTab().CurPane().HandleCommand(fmt.Sprintf("replace %s %s", "string", "foo"))
	injectString("ynyny")
	injectKey(tcell.KeyEscape, 0, tcell.ModNone)

	injectKey(tcell.KeyCtrlS, rune(tcell.KeyCtrlS), tcell.ModCtrl)

	data, err = os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, srTest3, string(data))
}

func TestMultiCursor(t *testing.T) {
	// TODO
}

func TestSettingsPersistence(t *testing.T) {
	// TODO
}

// more tests (rendering, tabs, plugins)?
