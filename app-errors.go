package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"howl-chat/internal/diagnostics"
	apperrors "howl-chat/internal/errors"
)

// ---------------------------------------------------------------------------
// Error System Integration
//
// The diagnostics owner (internal/diagnostics) owns logging, reporting and the
// UI event payloads. App supplies what is genuinely application policy: where
// the log lives, the persisted upload preference, the Wails event transport,
// and the decision to escalate a failure (including Level 0's shutdown sink).
// ---------------------------------------------------------------------------

// appDiagnosticPort is the application-supplied transport for diagnostics
// events: the same guarded Wails emission (and harness sink) every other event
// uses, so an error emitted without a live webview is dropped instead of
// terminating the process.
type appDiagnosticPort struct{ app *App }

func (p appDiagnosticPort) Emit(event string, data ...interface{}) {
	p.app.emitEvent(event, data...)
}

func (a *App) initErrorSystem() {
	owner, err := diagnostics.Open(diagnostics.Config{
		LogDir:        filepath.Join(a.getSettingsDir(), "logs"),
		AppVersion:    "1.0.0",
		UploadEnabled: a.loadErrorReportingSetting(),
	})
	if err != nil {
		// Fatal to error system but not to app - operate silently without logging
		return
	}
	a.diagnostics = owner
}

// bindErrorEmitter attaches the application's transport. It needs no webview:
// emitEvent delivers to a host's sink when one exists and drops the Wails call
// until a context carries an events channel, so binding at composition time
// lets every host receive diagnostics events as soon as it can consume them.
func (a *App) bindErrorEmitter() {
	if a.diagnostics == nil {
		return
	}
	a.diagnostics.Attach(appDiagnosticPort{app: a})
}

// ---------------------------------------------------------------------------
// Wails Bindings - Exposed to frontend
// ---------------------------------------------------------------------------

// GetErrorLog returns the last N error entries for display in settings.
func (a *App) GetErrorLog(count int) []map[string]interface{} {
	if a.diagnostics == nil {
		return nil
	}
	entries := a.diagnostics.Entries(count)
	var result []map[string]interface{}
	for _, e := range entries {
		result = append(result, map[string]interface{}{
			"time":      e.Time,
			"code":      e.Code,
			"severity":  e.Severity,
			"operation": e.Operation,
			"location":  e.Location,
		})
	}
	return result
}

// SetErrorReportingEnabled toggles the Vercel upload feature.
func (a *App) SetErrorReportingEnabled(enabled bool) {
	if a.diagnostics == nil {
		return
	}
	a.diagnostics.SetEnabled(enabled)
	_ = a.saveErrorReportingSetting(enabled)
}

// IsErrorReportingEnabled returns the current upload toggle state.
func (a *App) IsErrorReportingEnabled() bool {
	if a.diagnostics == nil {
		return false
	}
	return a.diagnostics.Enabled()
}

// SendErrorReport manually triggers an upload of recent errors.
func (a *App) SendErrorReport() (string, error) {
	if a.diagnostics == nil {
		return "", fmt.Errorf("error reporter not initialized")
	}
	if !a.diagnostics.Enabled() {
		return "", fmt.Errorf("error reporting is disabled")
	}
	report, err := a.diagnostics.BuildReport(50)
	if err != nil {
		return "", err
	}
	if err := a.diagnostics.Upload(report); err != nil {
		return "", err
	}
	return report.ReportID, nil
}

// GetErrorReportPreview returns a preview of what would be uploaded.
func (a *App) GetErrorReportPreview() (string, error) {
	if a.diagnostics == nil {
		return "", fmt.Errorf("error reporter not initialized")
	}
	report, err := a.diagnostics.BuildReport(20)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ---------------------------------------------------------------------------
// Settings persistence for error reporting toggle
// ---------------------------------------------------------------------------

func (a *App) loadErrorReportingSetting() bool {
	path := filepath.Join(a.getSettingsDir(), "error-reporting.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return false // Default: disabled
	}
	var cfg struct {
		Enabled bool `json:"enabled"`
	}
	_ = json.Unmarshal(data, &cfg)
	return cfg.Enabled
}

func (a *App) saveErrorReportingSetting(enabled bool) error {
	path := filepath.Join(a.getSettingsDir(), "error-reporting.json")
	data, err := json.MarshalIndent(map[string]bool{"enabled": enabled}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ---------------------------------------------------------------------------
// Convenience methods for subsystems
// ---------------------------------------------------------------------------

func (a *App) logError(code apperrors.Code, severity apperrors.Severity, operation, location string, err error) {
	a.diagnostics.Report(apperrors.New(code, severity, operation, location, err))
}

func (a *App) logErrorSilent(code apperrors.Code, severity apperrors.Severity, operation, location string, err error) {
	a.diagnostics.ReportSilent(apperrors.New(code, severity, operation, location, err))
}

func (a *App) logLLMLoadError(err error, location string) {
	a.diagnostics.ReportLLMLoad(err, location)
}

func (a *App) logLLMCrash(err error, location string) {
	a.diagnostics.ReportLLMCrash(err, location)
}

func (a *App) logImageLoadError(err error, location string) {
	a.diagnostics.ReportImageLoad(err, location)
}

func (a *App) logImageGenError(err error, location string) {
	a.diagnostics.ReportImageGen(err, location)
}

func (a *App) logTTSError(err error, location string) {
	a.diagnostics.ReportTTSSynth(err, location)
}

func (a *App) logChatSaveError(err error, location string) {
	a.diagnostics.ReportChatSave(err, location)
}

func (a *App) logResourceError(err error, location string) {
	a.diagnostics.ReportResourceAlloc(err, location)
}

// ---------------------------------------------------------------------------
// Panic recovery wrappers for Wails bindings
// ---------------------------------------------------------------------------

func (a *App) withRecovery(operation string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			a.diagnostics.ReportPanic(operation, "app-errors.go:recovery", r)
		}
	}()
	fn()
}

func (a *App) withRecoveryReturn(operation string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic in %s: %v", operation, r)
			a.diagnostics.ReportPanic(operation, "app-errors.go:recovery-return", r)
		}
	}()
	return fn()
}
