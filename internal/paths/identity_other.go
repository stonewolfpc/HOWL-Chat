//go:build !windows

package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// FileIdentity is the filesystem identity of a path. Howl ships on Windows; this
// keeps development builds and tests portable without pretending the Windows
// volume-serial API exists here.
type FileIdentity struct {
	Volume uint64
	Index  uint64
}

// resolvePhysical returns the canonical physical path for an existing path.
func resolvePhysical(path string) (string, error) {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("resolve %q: %w", path, err)
	}
	return resolved, nil
}

func identityOf(path string) (FileIdentity, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileIdentity{}, fmt.Errorf("stat %q: %w", path, err)
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return FileIdentity{}, fmt.Errorf("identity of %q: unsupported platform", path)
	}
	return FileIdentity{
		Volume: uint64(stat.Dev),
		Index:  uint64(stat.Ino),
	}, nil
}
