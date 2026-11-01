// Package paths resolves on-disk locations the app needs at runtime.
//
// The project previously hardcoded "d:\HOWL_Chat" throughout the backend,
// which breaks as soon as the drive letter changes (for example after a
// Windows reinstall). Everything that needs the data folders (Chats,
// Worlds, Scenarios, Characters, models, images) should ask this package.
package paths

import (
	"os"
	"path/filepath"
)

// markers are files that only exist in the project root.
var markers = []string{"wails.json", "go.mod"}

// maxDepth bounds the upward search so a stray go.mod high in the tree
// can never send us somewhere unexpected.
const maxDepth = 6

var cachedRoot string

// ProjectRoot returns the directory holding the app's data folders.
//
// It searches upward from the running executable first, which covers both
// `wails dev` (temporary binary inside the project) and `wails build`
// (build\bin\howl-chat.exe). If no marker is found it falls back to the
// current working directory, then to the executable's own directory.
func ProjectRoot() string {
	if cachedRoot != "" {
		return cachedRoot
	}
	cachedRoot = resolve()
	return cachedRoot
}

func resolve() string {
	exeDir := ""
	if exe, err := os.Executable(); err == nil {
		exeDir = filepath.Dir(exe)
		if root, ok := searchUp(exeDir); ok {
			return root
		}
	}
	if cwd, err := os.Getwd(); err == nil {
		if root, ok := searchUp(cwd); ok {
			return root
		}
		if exeDir == "" {
			return cwd
		}
		// Nothing looked like a project root; prefer where we were launched
		// from, since that is the historical behaviour.
		return cwd
	}
	if exeDir != "" {
		return exeDir
	}
	return "."
}

// searchUp walks from dir towards the filesystem root looking for a marker.
func searchUp(dir string) (string, bool) {
	for depth := 0; depth < maxDepth; depth++ {
		for _, marker := range markers {
			if _, err := os.Stat(filepath.Join(dir, marker)); err == nil {
				return dir, true
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}
