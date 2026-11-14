//go:build windows

package paths

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func knownFolder(name string) (string, error) {
	var id *windows.KNOWNFOLDERID
	switch name {
	case "documents":
		id = windows.FOLDERID_Documents
	case "pictures":
		id = windows.FOLDERID_Pictures
	case "videos":
		id = windows.FOLDERID_Videos
	default:
		return "", fmt.Errorf("unknown Windows known folder %q", name)
	}
	path, err := windows.KnownFolderPath(id, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return "", fmt.Errorf("resolve Windows known folder %q: %w", name, err)
	}
	return path, nil
}
