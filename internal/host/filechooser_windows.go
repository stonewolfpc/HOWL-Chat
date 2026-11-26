//go:build windows

package host

import (
	"fmt"
	"runtime"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

// File picking uses the common file dialog with no owner window for the same
// reason as ChooseFolder: an owned dialog disables the Wails window while it is
// open, and the focus event on the way back can hit go-webview2's fatal
// Chromium.Focus error path. With no owner, Windows never disables the app.
//
// The dialog is intentionally the classic GetOpenFileNameW picker; it is the
// same shell dialog Explorer uses and does not require a Wails window handle.

var (
	comdlg32                 = windows.NewLazySystemDLL("comdlg32.dll")
	procGetOpenFileNameW     = comdlg32.NewProc("GetOpenFileNameW")
	procCommDlgExtendedError = comdlg32.NewProc("CommDlgExtendedError")
)

const (
	ofnReadOnly        = 0x00000001
	ofnOverwritePrompt = 0x00000002
	ofnNoChangeDir     = 0x00000008
	ofnPathMustExist   = 0x00000800
	ofnFileMustExist   = 0x00001000
	ofnExplorer        = 0x00080000
)

type openFileNameW struct {
	StructSize      uint32
	Owner           windows.Handle
	Instance        windows.Handle
	Filter          *uint16
	CustomFilter    *uint16
	MaxCustomFilter uint32
	FilterIndex     uint32
	File            *uint16
	MaxFile         uint32
	FileTitle       *uint16
	MaxFileTitle    uint32
	InitialDir      *uint16
	Title           *uint16
	Flags           uint32
	FileOffset      uint16
	FileExtension   uint16
	DefExt          *uint16
	CustData        uintptr
	Hook            uintptr
	TemplateName    *uint16
	Reserved        unsafe.Pointer
	ReservedSize    uint32
	FlagsEx         uint32
}

// ChooseFile shows the common open dialog with no owner window and returns the
// selected path. An empty result with a nil error means the user cancelled.
func ChooseFile(title, initialDir, filter, defaultExt string) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// The common dialog and shell namespace are COM components. Initialise the
	// current thread as STA, matching ChooseFolder; RPC_E_CHANGED_MODE means the
	// thread already has a usable apartment.
	if hr, _, _ := procCoInitializeEx.Call(0, coinitApartmentThreaded); hr >= 0x80000000 && hr != 0x80010106 {
		return "", fmt.Errorf("file dialog: CoInitializeEx failed (0x%08x)", hr)
	}
	defer procCoUninitialize.Call()

	titlePtr, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return "", err
	}
	var initialPtr *uint16
	if initialDir != "" {
		if p, err := windows.UTF16PtrFromString(initialDir); err == nil {
			initialPtr = p
		}
	}
	var defExtPtr *uint16
	if defaultExt != "" {
		if p, err := windows.UTF16PtrFromString(defaultExt); err == nil {
			defExtPtr = p
		}
	}

	// The filter is a double-NUL-terminated sequence of label/pattern pairs.
	// Build the UTF-16 buffer directly because UTF16PtrFromString rejects the
	// embedded separators.
	filter = filter + "\x00\x00"
	filterUTF16 := utf16.Encode([]rune(filter))
	if len(filterUTF16) == 0 {
		filterUTF16 = []uint16{0, 0}
	}

	path := make([]uint16, 32768)
	ofn := openFileNameW{
		StructSize:  uint32(unsafe.Sizeof(openFileNameW{})),
		Owner:       0,
		Filter:      &filterUTF16[0],
		FilterIndex: 1,
		File:        &path[0],
		MaxFile:     uint32(len(path)),
		InitialDir:  initialPtr,
		Title:       titlePtr,
		Flags:       ofnNoChangeDir | ofnPathMustExist | ofnFileMustExist | ofnExplorer | ofnReadOnly,
		DefExt:      defExtPtr,
	}

	ret, _, _ := procGetOpenFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
	runtime.KeepAlive(filterUTF16)
	runtime.KeepAlive(initialPtr)
	runtime.KeepAlive(titlePtr)
	runtime.KeepAlive(defExtPtr)
	if ret == 0 {
		code, _, _ := procCommDlgExtendedError.Call()
		if code == 0 {
			return "", nil // cancelled
		}
		return "", fmt.Errorf("file dialog failed (0x%x)", code)
	}
	return windows.UTF16ToString(path), nil
}
