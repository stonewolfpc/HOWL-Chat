//go:build windows

package paths

import (
	"fmt"

	"golang.org/x/sys/windows"
)

// FileIdentity is the filesystem identity of a path: volume serial number plus
// file index. It is compared as a value, never as a string.
type FileIdentity struct {
	Volume uint64
	Index  uint64
}

// resolvePhysical returns the canonical physical path for an existing path,
// following junctions and symlinks.
//
// filepath.EvalSymlinks does not resolve junctions on this toolchain, so the
// handle-based final path is used instead. That is the only reliable way to learn
// where a reparse point actually lands.
func resolvePhysical(path string) (string, error) {
	ptr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return "", fmt.Errorf("encode %q: %w", path, err)
	}
	handle, err := windows.CreateFile(
		ptr,
		0,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		return "", fmt.Errorf("open %q: %w", path, err)
	}
	defer windows.CloseHandle(handle)

	buf := make([]uint16, 1024)
	n, err := windows.GetFinalPathNameByHandle(handle, &buf[0], uint32(len(buf)), 0)
	if err != nil {
		return "", fmt.Errorf("final path of %q: %w", path, err)
	}
	if int(n) > len(buf) {
		return "", fmt.Errorf("final path of %q exceeds buffer", path)
	}
	final := windows.UTF16ToString(buf[:n])
	return stripExtendedPrefix(final), nil
}

// stripExtendedPrefix removes the \\?\ and \\?\UNC\ length-extending prefixes
// that GetFinalPathNameByHandle returns.
func stripExtendedPrefix(p string) string {
	for _, prefix := range []string{`\\?\UNC\`, `\\?\`} {
		if len(p) > len(prefix) && p[:len(prefix)] == prefix {
			if prefix == `\\?\UNC\` {
				return `\\` + p[len(prefix):]
			}
			return p[len(prefix):]
		}
	}
	return p
}

// identityOf opens the path and reads its identity. FILE_FLAG_BACKUP_SEMANTICS
// is required to open a directory without read access to its contents.
func identityOf(path string) (FileIdentity, error) {
	ptr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return FileIdentity{}, fmt.Errorf("encode %q: %w", path, err)
	}
	handle, err := windows.CreateFile(
		ptr,
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		return FileIdentity{}, fmt.Errorf("open %q: %w", path, err)
	}
	defer windows.CloseHandle(handle)

	var info windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(handle, &info); err != nil {
		return FileIdentity{}, fmt.Errorf("identity of %q: %w", path, err)
	}
	return FileIdentity{
		Volume: uint64(info.VolumeSerialNumber),
		Index:  uint64(info.FileIndexHigh)<<32 | uint64(info.FileIndexLow),
	}, nil
}
