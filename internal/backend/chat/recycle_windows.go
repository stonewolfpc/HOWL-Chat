//go:build windows

package chat

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Deletion on Windows goes through the shell so the content lands in the Recycle
// Bin rather than being destroyed. Two rules hold here:
//
//  1. There is no permanent-delete fallback in this code. If the shell refuses,
//     the files stay where they are and the caller reports leftovers.
//  2. Windows may still choose to destroy rather than recycle (a volume without a
//     Recycle Bin, or one whose bin is full), because FOF_ALLOWUNDO preserves undo
//     information "if possible". That is the platform's decision, not a fallback
//     this program performs, and the deletion notice tells the user where the
//     files were so they can verify.

var (
	shell32                = windows.NewLazySystemDLL("shell32.dll")
	procSHFileOperationW   = shell32.NewProc("SHFileOperationW")
	procSHQueryRecycleBinW = shell32.NewProc("SHQueryRecycleBinW")
	ole32                  = windows.NewLazySystemDLL("ole32.dll")
	procCoInitializeEx     = ole32.NewProc("CoInitializeEx")
	procCoUninitialize     = ole32.NewProc("CoUninitialize")
)

const (
	foDelete          = 0x0003
	fofSilent         = 0x0004
	fofNoConfirmation = 0x0010
	fofAllowUndo      = 0x0040
	fofNoErrorUI      = 0x0400

	coinitApartmentThreaded = 0x2
)

// shFileOpStruct mirrors SHFILEOPSTRUCTW. Field order and alignment match the C
// layout on 64-bit Windows; pFrom must point at a double-null-terminated buffer.
type shFileOpStruct struct {
	hwnd                  uintptr
	wFunc                 uint32
	pFrom                 *uint16
	pTo                   *uint16
	fFlags                uint16
	fAnyOperationsAborted int32
	hNameMappings         uintptr
	lpszProgressTitle     *uint16
}

type shQueryRBInfo struct {
	cbSize      uint32
	i64Size     int64
	i64NumItems int64
}

// recycleTree moves a directory tree to the Recycle Bin. It returns an error when
// the shell refuses or aborts the operation; it never deletes permanently itself.
func recycleTree(path string) error {
	// Nothing to recycle is not a failure, and must not become one: the caller
	// treats a missing directory as already handled.
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		return nil
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// The shell APIs require COM on the calling thread. S_FALSE (1) means COM was
	// already initialised here, which is fine and must still be balanced.
	hr, _, _ := procCoInitializeEx.Call(0, coinitApartmentThreaded)
	if hr != 0 && hr != 1 {
		return fmt.Errorf("recycle: COM initialisation failed (0x%x)", hr)
	}
	defer procCoUninitialize.Call()

	from, err := windows.UTF16FromString(path)
	if err != nil {
		return fmt.Errorf("recycle: invalid path: %w", err)
	}
	from = append(from, 0) // the shell requires a second terminating NUL

	op := shFileOpStruct{
		wFunc:  foDelete,
		pFrom:  &from[0],
		fFlags: fofSilent | fofNoConfirmation | fofAllowUndo | fofNoErrorUI,
	}
	ret, _, _ := procSHFileOperationW.Call(uintptr(unsafe.Pointer(&op)))
	runtime.KeepAlive(from)

	if ret != 0 {
		return fmt.Errorf("recycle: the shell refused the delete (code %d)", ret)
	}
	if op.fAnyOperationsAborted != 0 {
		return fmt.Errorf("recycle: the shell aborted the delete")
	}
	return nil
}

// recycleBinAvailable reports whether the volume holding path has a Recycle Bin.
// It is used for diagnostics only: it cannot promise that a later delete will be
// recycled, so it never gates the attempt.
func recycleBinAvailable(path string) bool {
	root := volumeRoot(path)
	if root == "" {
		return false
	}
	ptr, err := windows.UTF16PtrFromString(root)
	if err != nil {
		return false
	}
	var info shQueryRBInfo
	info.cbSize = uint32(unsafe.Sizeof(info))
	hr, _, _ := procSHQueryRecycleBinW.Call(uintptr(unsafe.Pointer(ptr)), uintptr(unsafe.Pointer(&info)))
	return hr == 0
}

// volumeRoot returns the "X:\" root of a path, or "" for a path without one.
func volumeRoot(path string) string {
	volume := filepath.VolumeName(path)
	if len(volume) != 2 || volume[1] != ':' {
		return ""
	}
	return volume + `\`
}
