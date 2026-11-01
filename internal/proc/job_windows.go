//go:build windows

// Package proc manages child processes started by the app so they never
// outlive it.
package proc

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
// Call it right after cmd.Start().
func Attach(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	initJob()
	if job == 0 {
		return
	}
	_ = windows.AssignProcessToJobObject(job, windows.Handle(cmd.Process.Pid))
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

// KillTreeByName terminates every process with the given image name.
func KillTreeByName(name string) {
	if name == "" {
		return
	}
	c := exec.Command("taskkill", "/T", "/F", "/IM", name)
	HideWindow(c)
	_ = c.Run()
}

// CloseJob terminates every process still attached to the job.
func CloseJob() {
	initJob()
	if job == 0 {
		return
	}
	// Closing the handle triggers KILL_ON_JOB_CLOSE for the whole tree.
	_ = windows.CloseHandle(job)
	job = 0
}

// ReapOrphans kills llama-server/sd-server processes that belong to this
// project but survived a previous instance (crash or force-kill). Matching is
// by executable path, so copies living elsewhere are left alone.
func ReapOrphans(projectRoot string) int {
	if projectRoot == "" {
		return 0
	}
	root := strings.TrimRight(projectRoot, `\/`)
	script := fmt.Sprintf(
		"$root='%s'; "+
			"Get-CimInstance Win32_Process -Filter \"Name='llama-server.exe' or Name='sd-server.exe'\" | "+
			"Where-Object { $_.ExecutablePath -and $_.ExecutablePath -like ($root + '*') } | "+
			"ForEach-Object { Stop-Process -Id $_.ProcessId -Force -ErrorAction SilentlyContinue; $_.ProcessId }",
		strings.ReplaceAll(root, "'", "''"),
	)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	HideWindow(cmd)
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	killed := 0
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.TrimSpace(line) != "" {
			killed++
		}
	}
	if killed > 0 {
		fmt.Fprintf(os.Stderr, "[proc] reaped %d orphaned model process(es)\n", killed)
	}
	return killed
}
