package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"howl-chat/internal/errors"
)

// ---------------------------------------------------------------------------
// Error System Integration
// ---------------------------------------------------------------------------

func (a *App) initErrorSystem() {
	logDir := filepath.Join(a.getSettingsDir(), "logs")
	logger, err := errors.NewLogger(logDir)
	if err != nil {
		// Fatal to error system but not to app - operate silently without logging
		return
	}
	a.errorLogger = logger

	reporter := errors.NewReporter(logger, "1.0.0")
	a.errorReporter = reporter

	// Load user preference for upload
	if enabled := a.loadErrorReportingSetting(); enabled {
		reporter.SetEnabled(true)
	}

	// Only create emitter if context is available (after startup)
	if a.ctx != nil {
		a.errorEmitter = errors.NewCombinedSink(logger, errors.NewWailsEmitter(a.ctx))
	}
}

func (a *App) bindErrorEmitter() {
	if a.errorLogger == nil {
		return
	}
	if a.ctx != nil && a.errorEmitter == nil {
		a.errorEmitter = errors.NewCombinedSink(
			a.errorLogger,
			errors.NewWailsEmitter(a.ctx),
		)
	}
}

// ---------------------------------------------------------------------------
// Wails Bindings - Exposed to frontend
// ---------------------------------------------------------------------------

// GetErrorLog returns the last N error entries for display in settings.
func (a *App) GetErrorLog(count int) []map[string]interface{} {
	if a.errorLogger == nil {
		return nil
	}
	entries := a.errorLogger.ReadRecent(count)
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
	if a.errorReporter == nil {
		return
	}
	a.errorReporter.SetEnabled(enabled)
	_ = a.saveErrorReportingSetting(enabled)
}

// IsErrorReportingEnabled returns the current upload toggle state.
func (a *App) IsErrorReportingEnabled() bool {
	if a.errorReporter == nil {
		return false
	}
	return a.errorReporter.IsEnabled()
}

// SendErrorReport manually triggers an upload of recent errors.
func (a *App) SendErrorReport() (string, error) {
	if a.errorReporter == nil {
		return "", fmt.Errorf("error reporter not initialized")
	}
	if !a.errorReporter.IsEnabled() {
		return "", fmt.Errorf("error reporting is disabled")
	}
	report, err := a.errorReporter.BuildReport(50)
	if err != nil {
		return "", err
	}
	if err := a.errorReporter.Upload(report); err != nil {
		return "", err
	}
	return report.ReportID, nil
}

// GetErrorReportPreview returns a preview of what would be uploaded.
func (a *App) GetErrorReportPreview() (string, error) {
	if a.errorReporter == nil {
		return "", fmt.Errorf("error reporter not initialized")
	}
	report, err := a.errorReporter.BuildReport(20)
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

func (a *App) logError(code errors.Code, severity errors.Severity, operation, location string, err error) {
	if a.errorEmitter == nil {
		return
	}
	a.errorEmitter.Handle(errors.New(code, severity, operation, location, err))
}

func (a *App) logErrorSilent(code errors.Code, severity errors.Severity, operation, location string, err error) {
	if a.errorEmitter == nil {
		return
	}
	a.errorEmitter.HandleSilent(errors.New(code, severity, operation, location, err))
}

func (a *App) logLLMLoadError(err error, location string) {
	if a.errorEmitter == nil {
		return
	}
	a.errorEmitter.LogLLMLoad(err, location)
}

func (a *App) logLLMCrash(err error, location string) {
	if a.errorEmitter == nil {
		return
	}
	a.errorEmitter.LogLLMCrash(err, location)
}

func (a *App) logImageLoadError(err error, location string) {
	if a.errorEmitter == nil {
		return
	}
	a.errorEmitter.LogImageLoad(err, location)
}

func (a *App) logImageGenError(err error, location string) {
	if a.errorEmitter == nil {
		return
	}
	a.errorEmitter.LogImageGen(err, location)
}

func (a *App) logTTSError(err error, location string) {
	if a.errorEmitter == nil {
		return
	}
	a.errorEmitter.LogTTSSynth(err, location)
}

func (a *App) logChatSaveError(err error, location string) {
	if a.errorEmitter == nil {
		return
	}
	a.errorEmitter.LogChatSave(err, location)
}

func (a *App) logResourceError(err error, location string) {
	if a.errorEmitter == nil {
		return
	}
	a.errorEmitter.LogResourceAlloc(err, location)
}

// ---------------------------------------------------------------------------
// Panic recovery wrappers for Wails bindings
// ---------------------------------------------------------------------------

func (a *App) withRecovery(operation string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			if a.errorEmitter != nil {
				a.errorEmitter.LogSystemPanic(operation, "app-errors.go:recovery", r)
			}
		}
	}()
	fn()
}

func (a *App) withRecoveryReturn(operation string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic in %s: %v", operation, r)
			if a.errorEmitter != nil {
				a.errorEmitter.LogSystemPanic(operation, "app-errors.go:recovery-return", r)
			}
		}
	}()
	return fn()
}
