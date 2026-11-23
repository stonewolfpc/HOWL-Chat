//go:build !windows

package host

import "fmt"

// ChooseFolder exists so callers compile everywhere. Only the packaged Windows
// host has a native folder chooser; the browser harness reports the difference
// instead of pretending the dialog exists.
func ChooseFolder(title, initialDir string) (string, error) {
	return "", fmt.Errorf("folder dialog: not supported on this platform")
}
