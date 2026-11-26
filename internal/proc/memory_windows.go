//go:build windows

package proc

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// The pinned sd-server does not print a per-component residency summary on this
// build. The process working set is the allocation evidence it does give us:
// it is the operating system's measurement of what the engine actually holds,
// not Howl's estimate, and it can be captured after load and after generation.
type processMemoryCounters struct {
	cb                         uint32
	pageFaultCount             uint32
	peakWorkingSetSize         uintptr
	workingSetSize             uintptr
	quotaPeakPagedPoolUsage    uintptr
	quotaPagedPoolUsage        uintptr
	quotaPeakNonPagedPoolUsage uintptr
	quotaNonPagedPoolUsage     uintptr
	pagefileUsage              uintptr
	peakPagefileUsage          uintptr
}

var (
	psapi                    = syscall.NewLazyDLL("psapi.dll")
	procGetProcessMemoryInfo = psapi.NewProc("GetProcessMemoryInfo")
)

// WorkingSetMB returns the current working set of pid in MiB. False means the
// measurement could not be taken; callers must not turn that into zero.
func WorkingSetMB(pid int) (float64, bool) {
	if pid <= 0 {
		return 0, false
	}
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return 0, false
	}
	defer windows.CloseHandle(h)

	var counters processMemoryCounters
	counters.cb = uint32(unsafe.Sizeof(counters))
	r, _, _ := procGetProcessMemoryInfo.Call(
		uintptr(h),
		uintptr(unsafe.Pointer(&counters)),
		uintptr(counters.cb),
	)
	if r == 0 {
		return 0, false
	}
	return float64(counters.workingSetSize) / (1024 * 1024), true
}
