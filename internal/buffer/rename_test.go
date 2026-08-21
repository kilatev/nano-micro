package buffer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenameRebindsSharedBuffer(t *testing.T) {
	dir := t.TempDir()
	oldPath, newPath := filepath.Join(dir, "old.txt"), filepath.Join(dir, "new.txt")
	require.NoError(t, os.WriteFile(oldPath, []byte("before"), 0644))

	first, err := NewBufferFromFile(oldPath, BTDefault)
	require.NoError(t, err)
	second, err := NewBufferFromFile(oldPath, BTDefault)
	require.NoError(t, err)
	t.Cleanup(func() { first.Close(); second.Close() })
	require.Same(t, first.SharedBuffer, second.SharedBuffer)

	require.NoError(t, Rename(oldPath, newPath))
	assert.Equal(t, newPath, first.Path)
	assert.Equal(t, newPath, second.Path)
	assert.Equal(t, newPath, first.AbsPath)
	assert.Equal(t, newPath, second.AbsPath)

	first.Insert(first.End(), " after")
	assert.Equal(t, "before after", string(second.Bytes()))
	require.NoError(t, first.Save())
	_, err = os.Stat(oldPath)
	assert.True(t, os.IsNotExist(err))
	contents, err := os.ReadFile(newPath)
	require.NoError(t, err)
	assert.Equal(t, "before after\n", string(contents))
}

func TestRenameRejectsOpenDestination(t *testing.T) {
	dir := t.TempDir()
	oldPath, newPath := filepath.Join(dir, "old.txt"), filepath.Join(dir, "new.txt")
	require.NoError(t, os.WriteFile(oldPath, []byte("old"), 0644))
	require.NoError(t, os.WriteFile(newPath, []byte("new"), 0644))

	oldBuffer, err := NewBufferFromFile(oldPath, BTDefault)
	require.NoError(t, err)
	newBuffer, err := NewBufferFromFile(newPath, BTDefault)
	require.NoError(t, err)
	t.Cleanup(func() { oldBuffer.Close(); newBuffer.Close() })

	assert.Error(t, Rename(oldPath, newPath))
	assert.Equal(t, oldPath, oldBuffer.Path)
	assert.Equal(t, newPath, newBuffer.Path)
	assert.FileExists(t, oldPath)
	assert.FileExists(t, newPath)
}

func TestRenameFailureLeavesBufferPathUntouched(t *testing.T) {
	oldPath := filepath.Join(t.TempDir(), "old.txt")
	require.NoError(t, os.WriteFile(oldPath, []byte("old"), 0644))

	b, err := NewBufferFromFile(oldPath, BTDefault)
	require.NoError(t, err)
	t.Cleanup(b.Close)

	newPath := filepath.Join(filepath.Dir(oldPath), "missing", "new.txt")
	assert.Error(t, Rename(oldPath, newPath))
	assert.Equal(t, oldPath, b.Path)
	assert.Equal(t, oldPath, b.AbsPath)
	assert.FileExists(t, oldPath)
}
