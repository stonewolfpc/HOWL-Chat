//go:build !windows

package host

import "fmt"

// ChooseFile exists so callers compile everywhere. Only the packaged Windows
// host has a native file chooser; the browser harness uses its HTML input path.
func ChooseFile(title, initialDir, filter, defaultExt string) (string, error) {
	return "", fmt.Errorf("file dialog: not supported on this platform")
}
