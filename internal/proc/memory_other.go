//go:build !windows

package proc

// WorkingSetMB is unavailable on hosts without the Windows process APIs. The
// report says "not measured" rather than inventing a number.
func WorkingSetMB(pid int) (float64, bool) {
	return 0, false
}
