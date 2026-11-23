package diagnostics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	apperrors "howl-chat/internal/errors"
)

// ---------------------------------------------------------------------------
// Reporter - Builds privacy-safe error reports for optional upload.
//
// The wire payload carries no free-form application data: error entries carry
// code, severity, operation, location, time and the names of any context keys,
// never their values, and never the log message. System fields are fixed
// runtime identifiers.
// ---------------------------------------------------------------------------

const (
	uploadEndpoint    = "https://howl-error-reporter.vercel.app/api/report"
	uploadTimeout     = 10 * time.Second
	maxReportEntries  = 50
	maxUploadReplyLen = 512
)

// Reporter turns the log's recent entries into a privacy-safe Report and can
// upload it when the user has opted in. It holds the toggle and the endpoint;
// the application decides persistence of the preference.
type Reporter struct {
	mu         sync.RWMutex
	logger     *Logger
	appVersion string
	osVersion  string
	arch       string
	enabled    bool
	uploadURL  string
	client     *http.Client
	// onResult observes the asynchronous outcome of an upload. Production
	// leaves it nil; tests use it to await a specific outcome. Failures are
	// surfaced through the log below either way.
	onResult func(error)
}

// NewReporter creates a reporter bound to the given logger.
func NewReporter(logger *Logger, appVersion string) *Reporter {
	return &Reporter{
		logger:     logger,
		appVersion: appVersion,
		osVersion:  runtime.GOOS,
		arch:       runtime.GOARCH,
		uploadURL:  uploadEndpoint,
		client:     &http.Client{Timeout: uploadTimeout},
	}
}

// SetEnabled toggles the upload feature (user-controlled opt-in).
func (r *Reporter) SetEnabled(v bool) {
	r.mu.Lock()
	r.enabled = v
	r.mu.Unlock()
}

// IsEnabled returns whether upload is currently enabled.
func (r *Reporter) IsEnabled() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.enabled
}

// BuildReport creates a privacy-safe report from recent log entries. The entry
// count is capped at maxReportEntries whatever the caller asks for.
func (r *Reporter) BuildReport(maxEntries int) (*Report, error) {
	if r.logger == nil {
		return nil, fmt.Errorf("logger not initialized")
	}
	if maxEntries > maxReportEntries {
		maxEntries = maxReportEntries
	}

	entries := r.logger.ReadRecent(maxEntries)
	if len(entries) == 0 {
		return nil, fmt.Errorf("no recent errors to report")
	}

	// Build minimal system fingerprint (no PII)
	sys := SystemInfo{
		OS:         r.osVersion,
		Arch:       r.arch,
		GoVersion:  runtime.Version(),
		AppVersion: r.appVersion,
		CPUCount:   runtime.NumCPU(),
	}

	payload := make([]ErrorPayload, 0, len(entries))
	for _, e := range entries {
		payload = append(payload, ErrorPayload{
			Code:        e.Code,
			Severity:    e.Severity,
			Operation:   e.Operation,
			Location:    e.Location,
			Time:        e.Time,
			ContextKeys: contextKeys(e.Context),
		})
	}

	return &Report{
		ReportID:   generateReportID(),
		Generated:  time.Now().UTC().Format(time.RFC3339),
		System:     sys,
		Errors:     payload,
		ErrorCount: len(payload),
	}, nil
}

// Upload sends the report to the configured endpoint.
//
// It returns only the synchronous refusals (disabled, nil report, unbuildable
// request). The network exchange runs in its own goroutine so the binding
// answers without waiting; its outcome is never silent: a failure is written to
// the diagnostics log as NET-REQ, and an observer (tests) may watch it.
func (r *Reporter) Upload(report *Report) error {
	if !r.IsEnabled() {
		return fmt.Errorf("upload disabled")
	}
	if report == nil {
		return fmt.Errorf("nil report")
	}

	data, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("marshal report: %w", err)
	}

	r.mu.RLock()
	url := r.uploadURL
	client := r.client
	logger := r.logger
	observe := r.onResult
	r.mu.RUnlock()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Error-Report", "howl-chat")

	if client == nil {
		client = &http.Client{Timeout: uploadTimeout}
	}

	go func() {
		err := sendReport(client, req)
		if err != nil && logger != nil {
			logger.Log(apperrors.New(apperrors.NETRequestFail, apperrors.Medium,
				"error_report.upload", "internal/diagnostics/reporter.go:Upload", err))
		}
		if observe != nil {
			observe(err)
		}
	}()

	return nil
}

// sendReport performs one bounded POST and classifies its outcome. A non-2xx
// reply is an error, not a success; a hung endpoint ends at the client timeout.
func sendReport(client *http.Client, req *http.Request) error {
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("post report: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxUploadReplyLen))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			return fmt.Errorf("report endpoint returned %s", resp.Status)
		}
		return fmt.Errorf("report endpoint returned %s: %s", resp.Status, msg)
	}
	return nil
}

// SaveReport writes a report to the local logs directory for manual inspection.
func (r *Reporter) SaveReport(report *Report) (string, error) {
	if r.logger == nil {
		return "", fmt.Errorf("logger not initialized")
	}

	dir := filepath.Dir(r.logger.Path())
	path := filepath.Join(dir, fmt.Sprintf("report-%s.json", report.ReportID))

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal report: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("write report: %w", err)
	}
	return path, nil
}

// contextKeys returns the sorted names of a logged context map. Report payloads
// carry the names only, so an identifying value under an unremarkable key
// cannot leave the machine.
func contextKeys(ctx map[string]string) []string {
	if len(ctx) == 0 {
		return nil
	}
	keys := make([]string, 0, len(ctx))
	for k := range ctx {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ---------------------------------------------------------------------------
// Report structures - no free-form application data
// ---------------------------------------------------------------------------

type Report struct {
	ReportID   string         `json:"report_id"`
	Generated  string         `json:"generated"`
	System     SystemInfo     `json:"system"`
	Errors     []ErrorPayload `json:"errors"`
	ErrorCount int            `json:"error_count"`
}

type SystemInfo struct {
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	GoVersion  string `json:"go_version"`
	AppVersion string `json:"app_version"`
	CPUCount   int    `json:"cpu_count"`
}

type ErrorPayload struct {
	Code        string   `json:"code"`
	Severity    string   `json:"severity"`
	Operation   string   `json:"operation"`
	Location    string   `json:"location"`
	Time        string   `json:"time"`
	ContextKeys []string `json:"context_keys,omitempty"`
}

func generateReportID() string {
	return fmt.Sprintf("%d", time.Now().UTC().UnixNano())
}
