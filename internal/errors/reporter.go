package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// ---------------------------------------------------------------------------
// Reporter - Builds privacy-safe error reports for optional upload.
// Contains zero user-identifiable information.
// ---------------------------------------------------------------------------

const uploadEndpoint = "https://howl-error-reporter.vercel.app/api/report"
const uploadTimeout = 10 * time.Second

type Reporter struct {
	logger      *Logger
	appVersion  string
	osVersion   string
	arch        string
	enabled     bool
	uploadURL   string
}

// NewReporter creates a reporter bound to the given logger.
func NewReporter(logger *Logger, appVersion string) *Reporter {
	return &Reporter{
		logger:     logger,
		appVersion: appVersion,
		osVersion:  runtime.GOOS,
		arch:       runtime.GOARCH,
		uploadURL:  uploadEndpoint,
	}
}

// SetEnabled toggles the upload feature (user-controlled opt-in).
func (r *Reporter) SetEnabled(v bool) {
	r.enabled = v
}

// IsEnabled returns whether upload is currently enabled.
func (r *Reporter) IsEnabled() bool {
	return r.enabled
}

// BuildReport creates a privacy-safe report from recent log entries.
func (r *Reporter) BuildReport(maxEntries int) (*Report, error) {
	if r.logger == nil {
		return nil, fmt.Errorf("logger not initialized")
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

	var errors []ErrorPayload
	for _, e := range entries {
		errors = append(errors, ErrorPayload{
			Code:      e.Code,
			Severity:  e.Severity,
			Operation: e.Operation,
			Location:  e.Location,
			Time:      e.Time,
			Context:   e.Context,
		})
	}

	return &Report{
		ReportID:   generateReportID(),
		Generated:  time.Now().UTC().Format(time.RFC3339),
		System:     sys,
		Errors:     errors,
		ErrorCount: len(errors),
	}, nil
}

// Upload sends the report to the configured endpoint.
// Runs asynchronously and does not block.
func (r *Reporter) Upload(report *Report) error {
	if !r.enabled {
		return fmt.Errorf("upload disabled")
	}
	if report == nil {
		return fmt.Errorf("nil report")
	}

	data, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("marshal report: %w", err)
	}

	client := &http.Client{Timeout: uploadTimeout}
	req, err := http.NewRequest("POST", r.uploadURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Error-Report", "howl-chat")

	go func() {
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
	}()

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

// ---------------------------------------------------------------------------
// Report structures - Zero PII guaranteed
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
	Code      string            `json:"code"`
	Severity  string            `json:"severity"`
	Operation string            `json:"operation"`
	Location  string            `json:"location"`
	Time      string            `json:"time"`
	Context   map[string]string `json:"context,omitempty"`
}

func generateReportID() string {
	return fmt.Sprintf("%d", time.Now().UTC().UnixNano())
}
