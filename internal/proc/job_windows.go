//go:build windows

// Package proc manages child processes started by the app so they never
// outlive it.
package proc

import (
	"fmt"
	"os/exec"
	"strconv"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	jobOnce sync.Once
	job     windows.Handle
)

// initJob creates a job object with KILL_ON_JOB_CLOSE. Every attached child
// dies when the last handle closes, which happens when this process exits —
// even if the app is force-killed or crashes.
func initJob() {
	jobOnce.Do(func() {
		h, err := windows.CreateJobObject(nil, nil)
		if err != nil {
			return
		}
		info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
			BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
				LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
			},
		}
		if _, err := windows.SetInformationJobObject(
			h,
			windows.JobObjectExtendedLimitInformation,
			uintptr(unsafe.Pointer(&info)),
			uint32(unsafe.Sizeof(info)),
		); err != nil {
			windows.CloseHandle(h)
			return
		}
		job = h
	})
}

// Attach binds an already-started child process to the job.
// Call it right after cmd.Start() — or use Start, which closes the race between
// creation and assignment.
//
// AssignProcessToJobObject needs a process HANDLE, not a PID. Passing the PID
// failed with "The handle is invalid" and was silently ignored, so no child was
// ever actually attached. The handle must be closed here; the job keeps its own
// reference for as long as the process lives.
func Attach(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return fmt.Errorf("proc: cannot attach a command that has not started")
	}
	initJob()
	if job == 0 {
		return fmt.Errorf("proc: no job object available")
	}
	handle, err := windows.OpenProcess(windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE, false, uint32(cmd.Process.Pid))
	if err != nil {
		return fmt.Errorf("proc: open pid %d for job assignment: %w", cmd.Process.Pid, err)
	}
	defer windows.CloseHandle(handle)
	if err := windows.AssignProcessToJobObject(job, handle); err != nil {
		return fmt.Errorf("proc: assign pid %d to job: %w", cmd.Process.Pid, err)
	}
	return nil
}

// Start launches a child suspended, assigns it to the job, and only then
// resumes it. No child code runs before ownership is established, so a crash at
// any point cannot leave an unowned process behind. If ownership cannot be
// established the child is killed and Start returns the error; it never falls
// back to an unowned process.
func Start(cmd *exec.Cmd) error {
	return StartWithOwnership(cmd, nil)
}

// StartWithOwnership is Start with one extra step between job assignment and
// the child's first instruction. owned runs while the child is still suspended,
// so a caller that must record ownership (the process registry) can fail the
// start instead of racing the first instruction.
func StartWithOwnership(cmd *exec.Cmd, owned func() error) error {
	if cmd == nil {
		return fmt.Errorf("proc: nil command")
	}
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.CreationFlags |= windows.CREATE_SUSPENDED
	if err := cmd.Start(); err != nil {
		return err
	}
	pid := cmd.Process.Pid
	failed := func(err error) error {
		KillTree(pid)
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
		return fmt.Errorf("proc: pid %d not started under ownership: %w", pid, err)
	}
	if err := Attach(cmd); err != nil {
		return failed(err)
	}
	if owned != nil {
		if err := owned(); err != nil {
			return failed(err)
		}
	}
	if err := resumeProcess(pid); err != nil {
		return failed(err)
	}
	return nil
}

// resumeProcess resumes every thread of a process created with CREATE_SUSPENDED.
// os/exec closes the main thread handle after CreateProcess, so the threads are
// found again through a Toolhelp snapshot.
func resumeProcess(pid int) error {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ThreadEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	resumed := 0
	var firstErr error
	for err := windows.Thread32First(snapshot, &entry); err == nil; err = windows.Thread32Next(snapshot, &entry) {
		if entry.OwnerProcessID != uint32(pid) {
			continue
		}
		thread, openErr := windows.OpenThread(windows.THREAD_SUSPEND_RESUME, false, entry.ThreadID)
		if openErr != nil {
			if firstErr == nil {
				firstErr = openErr
			}
			continue
		}
		if _, resumeErr := windows.ResumeThread(thread); resumeErr != nil {
			if firstErr == nil {
				firstErr = resumeErr
			}
		} else {
			resumed++
		}
		windows.CloseHandle(thread)
	}
	if resumed == 0 {
		if firstErr == nil {
			firstErr = fmt.Errorf("no thread found")
		}
		return fmt.Errorf("resume pid %d: %w", pid, firstErr)
	}
	return nil
}

// HideWindow stops a helper process from flashing a console window.
func HideWindow(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.HideWindow = true
}

// KillTree force-terminates a child and everything it spawned.
func KillTree(pid int) {
	if pid <= 0 {
		return
	}
	c := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid))
	HideWindow(c)
	_ = c.Run()
}

// CloseJob terminates every process still attached to the job.
// It is unchanged by the suspended-start path; shutdown still closes the job.
func CloseJob() {
	initJob()
	if job == 0 {
		return
	}
	// Closing the handle triggers KILL_ON_JOB_CLOSE for the whole tree.
	_ = windows.CloseHandle(job)
	job = 0
}
