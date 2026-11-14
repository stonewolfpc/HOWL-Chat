//go:build windows

package fakehelper

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func hasReparseComponent(path string) (bool, error) {
	if path == "" {
		return false, fmt.Errorf("path is empty")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}
	for current := filepath.Clean(abs); ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err != nil {
			return false, err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return true, nil
		}
		ptr, err := windows.UTF16PtrFromString(current)
		if err != nil {
			return false, err
		}
		attrs, err := windows.GetFileAttributes(ptr)
		if err != nil {
			return false, err
		}
		if attrs&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
			return true, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}
	return false, nil
}
