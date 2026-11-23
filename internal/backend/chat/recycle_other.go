//go:build !windows

package chat

import "os"

// recycleTree removes a directory tree on platforms without a Recycle Bin. Howl
// ships on Windows, where recycle_windows.go is the real path; this exists so the
// package builds and its tests run elsewhere, and it is not a fallback chosen at
// runtime on Windows.
func recycleTree(path string) error {
	return os.RemoveAll(path)
}

// recycleBinAvailable has no meaning off Windows.
func recycleBinAvailable(string) bool { return false }
