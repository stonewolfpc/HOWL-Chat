// Package loader holds process diagnostics for model loading.
//
// A packaged GUI application is started without a usable console. On Windows
// its stdout/stderr handles can exist but be invalid, so writing to os.Stderr
// fails. The old loader used io.MultiWriter(os.Stderr, &stderrBuf): when the
// console write failed, MultiWriter stopped and the byte buffer never received
// the child's output. The failure was then invisible, and the only message the
// user saw was "llama-server exited unexpectedly".
//
// This package makes the bounded capture the primary destination and treats
// the console as an optional mirror whose write errors are ignored. Capture
// therefore keeps working even when os.Stderr is nil, closed or invalid. The
// original MultiWriter hypothesis is not proven; this removes the dependency
// on it either way.
package loader

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultCaptureBytes bounds retained child output. A runaway process can
	// produce megabytes of log lines; memory stays bounded at this size.
	DefaultCaptureBytes = 256 * 1024
	// FailureTailBytes is how much of the capture is written into a failure
	// record and shown to the user.
	FailureTailBytes = 8 * 1024
)

// Capture is an io.Writer that retains only the last Limit bytes written to
// it. It is safe for concurrent writes (stdout and stderr both point at it)
// and never returns an error, so a failing console can never stop the child's
// output from being captured.
type Capture struct {
	mu        sync.Mutex
	limit     int
	buf       []byte
	total     int64
	truncated bool
}

// NewCapture creates a bounded capture. A non-positive limit uses the default.
func NewCapture(limit int) *Capture {
	if limit <= 0 {
		limit = DefaultCaptureBytes
	}
	return &Capture{limit: limit}
}

// Write implements io.Writer.
func (c *Capture) Write(p []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.total += int64(len(p))
	c.buf = append(c.buf, p...)
	if len(c.buf) > c.limit {
		drop := len(c.buf) - c.limit
		c.buf = append(c.buf[:0], c.buf[drop:]...)
		c.truncated = true
	}
	return len(p), nil
}

// Tail returns the retained output, clipped to the last max bytes when max is
// positive.
func (c *Capture) Tail(max int) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	data := c.buf
	if max > 0 && len(data) > max {
		data = data[len(data)-max:]
	}
	return string(data)
}

// Stats reports the total bytes seen, the retained bytes and whether any
// output was dropped.
func (c *Capture) Stats() (total int64, retained int, truncated bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.total, len(c.buf), c.truncated
}

// Reset clears the capture for a retry.
func (c *Capture) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.buf = nil
	c.total = 0
	c.truncated = false
}

// SafeConsole mirrors child output to w while ignoring write failures. A nil
// writer becomes io.Discard, which is what a packaged GUI without a console
// effectively has. Returning len(p), nil is deliberate: io.MultiWriter stops
// on the first error, and console availability must never stop capture.
func SafeConsole(w io.Writer) io.Writer {
	if w == nil {
		return io.Discard
	}
	return consoleWriter{w: w}
}

type consoleWriter struct{ w io.Writer }

func (c consoleWriter) Write(p []byte) (int, error) {
	_, _ = c.w.Write(p)
	return len(p), nil
}

// FailureStage names where a model load stopped.
type FailureStage string

const (
	StageAdmission FailureStage = "admission"
	StageProjector FailureStage = "projector"
	StageSpawn     FailureStage = "spawn"
	StageStartup   FailureStage = "startup"
	StageTimeout   FailureStage = "timeout"
)

func (s FailureStage) describe() string {
	switch s {
	case StageAdmission:
		return "classifying the model"
	case StageProjector:
		return "pairing the companion projector"
	case StageSpawn, StageStartup:
		return "starting llama-server"
	case StageTimeout:
		return "waiting for llama-server to become ready"
	default:
		return string(s)
	}
}

// FailureRecord is one diagnosable model-load failure.
type FailureRecord struct {
	Stage      FailureStage
	ModelPath  string
	ModelName  string
	Arch       string
	Projector  string
	Executable string
	Args       []string
	ExitStatus string
	// Output is an already-bounded tail. WriteFailureRecord does not expand it.
	Output          string
	OutputBytes     int64
	OutputTruncated bool
	Detail          string
	// RecordPath is filled in after the record is persisted.
	RecordPath string
}

// Format renders the full record for the log file.
func (r FailureRecord) Format() string {
	var b strings.Builder
	fmt.Fprintf(&b, "time: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(&b, "stage: %s\n", r.Stage)
	if r.ModelPath != "" {
		fmt.Fprintf(&b, "model: %s\n", r.ModelPath)
	}
	if r.ModelName != "" {
		fmt.Fprintf(&b, "model_name: %s\n", r.ModelName)
	}
	if r.Arch != "" {
		fmt.Fprintf(&b, "architecture: %s\n", r.Arch)
	}
	if r.Projector != "" {
		fmt.Fprintf(&b, "projector: %s\n", r.Projector)
	}
	if r.Executable != "" {
		fmt.Fprintf(&b, "executable: %s\n", r.Executable)
	}
	if len(r.Args) > 0 {
		fmt.Fprintf(&b, "args: %s\n", strings.Join(r.Args, " "))
	}
	if r.ExitStatus != "" {
		fmt.Fprintf(&b, "exit: %s\n", r.ExitStatus)
	}
	if r.Detail != "" {
		fmt.Fprintf(&b, "detail: %s\n", r.Detail)
	}
	if r.OutputBytes > 0 {
		fmt.Fprintf(&b, "output_bytes: %d (truncated=%t)\n", r.OutputBytes, r.OutputTruncated)
	}
	if strings.TrimSpace(r.Output) != "" {
		b.WriteString("output_tail:\n")
		for _, line := range strings.Split(strings.TrimRight(r.Output, "\r\n"), "\n") {
			fmt.Fprintf(&b, "  %s\n", strings.TrimRight(line, "\r"))
		}
	}
	return b.String()
}

// UserMessage is the concise, actionable text shown in the UI and error log.
func (r FailureRecord) UserMessage() string {
	subject := r.ModelName
	if subject == "" {
		subject = filepath.Base(r.ModelPath)
	}
	if subject == "" {
		subject = "model"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Model load failed while %s: %s", r.Stage.describe(), subject)
	if r.Arch != "" {
		fmt.Fprintf(&b, " (architecture %s)", r.Arch)
	}
	if r.ExitStatus != "" {
		fmt.Fprintf(&b, ", llama-server %s", r.ExitStatus)
	}
	b.WriteString(".")
	if r.Detail != "" {
		b.WriteString(" ")
		b.WriteString(r.Detail)
	}
	if tail := compactTail(r.Output); tail != "" {
		b.WriteString(" Last output: ")
		b.WriteString(tail)
	}
	if r.RecordPath != "" {
		b.WriteString(" Diagnostics: ")
		b.WriteString(r.RecordPath)
	}
	return b.String()
}

// compactTail keeps the last few non-empty lines within max characters so the
// user message stays readable even when the capture is full of progress bars.
func compactTail(output string) string {
	const max = 480
	output = strings.TrimSpace(output)
	if output == "" {
		return ""
	}
	lines := strings.Split(output, "\n")
	var kept []string
	for i := len(lines) - 1; i >= 0 && len(kept) < 3; i-- {
		line := strings.TrimSpace(strings.TrimRight(lines[i], "\r"))
		if line == "" {
			continue
		}
		kept = append([]string{line}, kept...)
	}
	joined := strings.Join(kept, " | ")
	if len(joined) > max {
		joined = joined[len(joined)-max:]
	}
	return joined
}

// WriteFailureRecord appends a record to <dir>/model-load-failures.log and
// returns the log path. Failure to persist is returned but must not mask the
// original load error; callers keep the in-memory message either way.
func WriteFailureRecord(dir string, rec FailureRecord) (string, error) {
	if strings.TrimSpace(dir) == "" {
		return "", fmt.Errorf("diagnostics directory is empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create diagnostics directory: %w", err)
	}
	path := filepath.Join(dir, "model-load-failures.log")

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return "", fmt.Errorf("open failure log: %w", err)
	}
	defer file.Close()

	entry := fmt.Sprintf("===== %s =====\n%s\n", time.Now().Format(time.RFC3339), rec.Format())
	if _, err := file.WriteString(entry); err != nil {
		return "", fmt.Errorf("write failure log: %w", err)
	}
	return path, nil
}

// LogDir is the diagnostics directory next to the app's data folders.
func LogDir(projectRoot string) string {
	return filepath.Join(projectRoot, "logs")
}

// ExitStatus renders a process result without losing the exit code.
func ExitStatus(err error) string {
	if err == nil {
		return ""
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if exitErr.ProcessState != nil {
			status := exitErr.ProcessState.String()
			if code := exitErr.ExitCode(); code >= 0 && !strings.Contains(status, strconv.Itoa(code)) {
				return fmt.Sprintf("%s (code %d)", status, code)
			}
			return status
		}
		if code := exitErr.ExitCode(); code >= 0 {
			return "exit status " + strconv.Itoa(code)
		}
	}
	return err.Error()
}
