//go:build !windows

package durable

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// fileLock mirrors the Windows lock with an advisory exclusive lock, so the
// package builds and its tests run off Windows. The application ships on Windows;
// this path exists for portability of the primitive, not as a supported target.
type fileLock struct {
	f *os.File
}

func acquireLock(path string, timeout time.Duration) (*fileLock, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open lock: %w", err)
	}
	deadline := time.Now().Add(timeout)
	for {
		err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return &fileLock{f: f}, nil
		}
		if time.Now().After(deadline) {
			_ = f.Close()
			return nil, ErrBusy
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func (l *fileLock) release() {
	if l == nil || l.f == nil {
		return
	}
	_ = syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
	_ = l.f.Close()
	l.f = nil
}
