package buffer

import (
	"fmt"
	"os"

	"github.com/micro-editor/micro/v2/internal/util"
)

// Rename moves a path and rebinds its open shared buffer, if any.
func Rename(oldPath, newPath string) error {
	oldAbs, newAbs := util.ResolvePath(oldPath), util.ResolvePath(newPath)
	var renamed *Buffer
	for _, b := range OpenBuffers {
		if b.AbsPath == oldAbs && b.Type != BTInfo {
			renamed = b
			break
		}
	}
	for _, b := range OpenBuffers {
		if b.AbsPath == newAbs && b.Type != BTInfo && (renamed == nil || b.SharedBuffer != renamed.SharedBuffer) {
			return fmt.Errorf("destination is already open: %s", newPath)
		}
	}
	if err := os.Rename(oldAbs, newAbs); err != nil {
		return err
	}
	if renamed == nil {
		return nil
	}

	renamed.RemoveBackup()
	renamed.Path = newPath
	renamed.AbsPath = newAbs
	renamed.UpdateModTime()
	renamed.ReloadSettings(true)
	if renamed.Modified() {
		renamed.RequestBackup()
	}
	return nil
}
