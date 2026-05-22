package errors

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Logger - Silent structured logger. Writes JSONL to file only.
// No stdout/stderr output. Zero console noise during normal operation.
// ---------------------------------------------------------------------------

const (
	maxLogFileSize = 50 << 20 // 50 MB
	maxLogFiles    = 5
)

type Logger struct {
	mu      sync.RWMutex
	file    *os.File
	dir     string
	enabled bool
	path    string
}

// NewLogger creates a silent file logger. Writes to <logDir>/howl-errors-YYYYMMDD.jsonl
func NewLogger(logDir string) (*Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}

	path := filepath.Join(logDir, fmt.Sprintf("howl-errors-%s.jsonl", time.Now().UTC().Format("20060102")))
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	return &Logger{
		file:    f,
		dir:     logDir,
		enabled: true,
		path:    path,
	}, nil
}

// Log writes a TrackedError to the log file as a single JSON line.
// Thread-safe. Silent - no console output.
func (l *Logger) Log(err *TrackedError) {
	if l == nil || err == nil {
		return
	}

	l.mu.RLock()
	if !l.enabled || l.file == nil {
		l.mu.RUnlock()
		return
	}
	l.mu.RUnlock()

	entry := logEntry{
		Time:      err.Time.Format(time.RFC3339Nano),
		Code:      string(err.Code),
		Severity:  err.Severity.String(),
		Operation: err.Operation,
		Location:  err.Location,
		Msg:       err.Err.Error(),
		Context:   sanitizeContext(err.Context),
	}

	data, _ := json.Marshal(entry)

	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file == nil {
		return
	}

	_, _ = l.file.Write(data)
	_, _ = l.file.WriteString("\n")
	_ = l.file.Sync()

	l.rotateIfNeeded()
}

// LogRaw logs a pre-built entry directly. Use for batch operations.
func (l *Logger) LogRaw(entry logEntry) {
	if l == nil {
		return
	}

	l.mu.RLock()
	if !l.enabled || l.file == nil {
		l.mu.RUnlock()
		return
	}
	l.mu.RUnlock()

	data, _ := json.Marshal(entry)

	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file == nil {
		return
	}
	_, _ = l.file.Write(data)
	_, _ = l.file.WriteString("\n")
}

// SetEnabled toggles logging without closing the file.
func (l *Logger) SetEnabled(v bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = v
}

// IsEnabled returns the current enabled state.
func (l *Logger) IsEnabled() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.enabled
}

// Close flushes and closes the log file.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		_ = l.file.Sync()
		err := l.file.Close()
		l.file = nil
		return err
	}
	return nil
}

// Path returns the current log file path.
func (l *Logger) Path() string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.path
}

// ReadRecent returns the last N log entries from the current file.
func (l *Logger) ReadRecent(n int) []logEntry {
	l.mu.RLock()
	path := l.path
	l.mu.RUnlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var entries []logEntry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e logEntry
		if json.Unmarshal(line, &e) == nil {
			entries = append(entries, e)
		}
	}

	if len(entries) <= n {
		return entries
	}
	return entries[len(entries)-n:]
}

func (l *Logger) rotateIfNeeded() {
	info, err := l.file.Stat()
	if err != nil {
		return
	}
	if info.Size() < maxLogFileSize {
		return
	}

	_ = l.file.Close()
	timestamp := time.Now().UTC().Format("20060102_150405")
	newPath := filepath.Join(l.dir, fmt.Sprintf("howl-errors-%s.jsonl", timestamp))
	_ = os.Rename(l.path, newPath)

	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l.file = nil
		return
	}
	l.file = f
	l.pruneOldFiles()
}

func (l *Logger) pruneOldFiles() {
	entries, err := os.ReadDir(l.dir)
	if err != nil {
		return
	}

	var files []os.DirEntry
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".jsonl" {
			files = append(files, e)
		}
	}

	if len(files) <= maxLogFiles {
		return
	}

	// Simple sort by name (timestamp prefix) and remove oldest
	for i := 0; i < len(files)-maxLogFiles; i++ {
		_ = os.Remove(filepath.Join(l.dir, files[i].Name()))
	}
}

// ---------------------------------------------------------------------------
// Internal types
// ---------------------------------------------------------------------------

type logEntry struct {
	Time      string            `json:"t"`
	Code      string            `json:"c"`
	Severity  string            `json:"s"`
	Operation string            `json:"op"`
	Location  string            `json:"loc"`
	Msg       string            `json:"msg"`
	Context   map[string]string `json:"ctx,omitempty"`
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
