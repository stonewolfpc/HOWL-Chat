//go:build windows

package host

import (
	"fmt"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Folder picking goes through the shell's own chooser with no owner window.
//
// Wails' runtime.OpenDirectoryDialog runs the common item dialog with the main
// window as its parent, so Windows disables that window for the duration. When
// the dialog closes, the main window receives focus while it is still disabled,
// Wails forwards that to go-webview2's Chromium.Focus, MoveFocus returns
// E_INVALIDARG, and go-webview2's error callback answers with os.Exit(1): the
// whole app disappears the moment Add Folder is used (observed twice on
// 2026-11-22, once from the packaged app on the real profile).
//
// A chooser with no owner does not disable the app window, so that path is never
// taken. The trade-off is the shell's classic folder tree instead of the modern
// picker; it is the same dialog Explorer shows for "Select Folder" on older
// builds, and it still honours the starting folder.

var (
	shell32                  = windows.NewLazySystemDLL("shell32.dll")
	ole32                    = windows.NewLazySystemDLL("ole32.dll")
	user32                   = windows.NewLazySystemDLL("user32.dll")
	procSHBrowseForFolderW   = shell32.NewProc("SHBrowseForFolderW")
	procSHGetPathFromIDListW = shell32.NewProc("SHGetPathFromIDListW")
	procCoInitializeEx       = ole32.NewProc("CoInitializeEx")
	procCoUninitialize       = ole32.NewProc("CoUninitialize")
	procSendMessageW         = user32.NewProc("SendMessageW")
	procSetForegroundWindow  = user32.NewProc("SetForegroundWindow")
	procCoTaskMemFree        = ole32.NewProc("CoTaskMemFree")
)

const (
	bifReturnOnlyFSDirs     = 0x0001
	bifEditBox              = 0x0010
	bifNewDialogStyle       = 0x0040
	coinitApartmentThreaded = 0x2

	bffmInitialized       = 0x0001
	bffmSetSelectionValue = 0x0467 // BFFM_SETSELECTIONW
)

type browseInfoW struct {
	hwndOwner      windows.Handle
	pidlRoot       uintptr
	pszDisplayName *uint16
	lpszTitle      *uint16
	ulFlags        uint32
	lpfn           uintptr
	lParam         unsafe.Pointer // LPARAM: BFFM_SETSELECTIONW's path, kept alive by KeepAlive
	iImage         int32
}

// ChooseFolder shows the shell folder chooser and returns the selected path.
// An empty result with a nil error means the user cancelled.
func ChooseFolder(title, initialDir string) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// The shell dialog is an STA component; a fresh OS thread has no apartment.
	if hr, _, _ := procCoInitializeEx.Call(0, coinitApartmentThreaded); hr >= 0x80000000 && hr != 0x80010106 { // RPC_E_CHANGED_MODE
		return "", fmt.Errorf("folder dialog: CoInitializeEx failed (0x%08x)", hr)
	}
	defer procCoUninitialize.Call()

	titlePtr, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return "", err
	}
	display := make([]uint16, windows.MAX_PATH)

	var initialPtr *uint16
	if initialDir != "" {
		if p, err := windows.UTF16PtrFromString(initialDir); err == nil {
			initialPtr = p
		}
	}

	// BFFM_INITIALIZED is the only message answered: it sets the starting folder
	// through the pointer handed in lParam, which stays alive for the whole call.
	callback := windows.NewCallback(func(hwnd uintptr, msg uint32, lParam, lpData uintptr) uintptr {
		if msg == bffmInitialized {
			if lpData != 0 {
				procSendMessageW.Call(hwnd, bffmSetSelectionValue, 1, lpData)
			}
			// The chooser has no owner window, so Windows does not raise it over
			// the app on its own; without this it opens behind the WebView and
			// Add Folder looks like it did nothing.
			procSetForegroundWindow.Call(hwnd)
		}
		return 0
	})

	bi := browseInfoW{
		pszDisplayName: &display[0],
		lpszTitle:      titlePtr,
		ulFlags:        bifReturnOnlyFSDirs | bifEditBox | bifNewDialogStyle,
		lpfn:           callback,
	}
	if initialPtr != nil {
		bi.lParam = unsafe.Pointer(initialPtr)
	}

	pidl, _, _ := procSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(&bi)))
	runtime.KeepAlive(initialPtr)
	if pidl == 0 {
		return "", nil // cancelled
	}
	defer procCoTaskMemFree.Call(pidl) // ITEMIDLIST is freed through the proc itself: a uintptr -> unsafe.Pointer conversion here is what vet rejects, and this needs none.

	path := make([]uint16, windows.MAX_PATH)
	ok, _, err := procSHGetPathFromIDListW.Call(pidl, uintptr(unsafe.Pointer(&path[0])))
	if ok == 0 {
		if err != nil && err != windows.ERROR_SUCCESS {
			return "", fmt.Errorf("folder dialog: cannot read the chosen path: %w", err)
		}
		return "", fmt.Errorf("folder dialog: cannot read the chosen path")
	}
	return windows.UTF16ToString(path), nil
}
