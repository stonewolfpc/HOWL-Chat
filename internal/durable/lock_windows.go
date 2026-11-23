//go:build windows

package durable

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows"
)

// fileLock is an exclusive whole-byte-range lock on the lock file. Windows
// releases it when the owning process exits, so a crashed writer cannot wedge the
// application; the lock *file* may remain and is harmless.
type fileLock struct {
	f  *os.File
	ov windows.Overlapped
}

// acquireLock takes the exclusive lock, waiting up to timeout. Returning ErrBusy
// is the only correct outcome when the wait expires: there is no bypass.
func acquireLock(path string, timeout time.Duration) (*fileLock, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open lock: %w", err)
	}
	deadline := time.Now().Add(timeout)
	for {
		ov := new(windows.Overlapped)
		err := windows.LockFileEx(
			windows.Handle(f.Fd()),
			windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY,
			0, 1, 0, ov,
		)
		if err == nil {
			return &fileLock{f: f, ov: *ov}, nil
		}
		if !isLockBusy(err) {
			_ = f.Close()
			return nil, fmt.Errorf("lock: %w", err)
		}
		if time.Now().After(deadline) {
			_ = f.Close()
			return nil, ErrBusy
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func isLockBusy(err error) bool {
	return err == windows.ERROR_LOCK_VIOLATION
}

func (l *fileLock) release() {
	if l == nil || l.f == nil {
		return
	}
	_ = windows.UnlockFileEx(windows.Handle(l.f.Fd()), 0, 1, 0, &l.ov)
	_ = l.f.Close()
	l.f = nil
}
