package diagnostics

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"sync"
	"time"

	"howl-chat/internal/errors"
)

// errNotOpen is what report operations return when the owner was never opened.
var errNotOpen = stderrors.New("diagnostics not initialized")

// Emitter is the application-supplied transport for user-facing diagnostics.
// The owner decides which event and payload to send; the application supplies
// the connection (Wails events in the desktop app, the harness sink in the
// browser host, a recorder in tests) and any environment guard.
type Emitter interface {
	Emit(event string, data ...interface{})
}

// Config is the application's explicit wiring of the diagnostics owner: where
// the log lives, which version reports identify, and the persisted upload
// preference. The owner never reads settings itself.
type Config struct {
	LogDir        string
	AppVersion    string
	UploadEnabled bool
}

// Owner is the one holder of the diagnostic responsibilities: the structured
// error log and its rotation, privacy-safe report building/upload/saving, and
// the UI event payloads derived from tracked errors. It is not an escalation
// policy: the application decides what to report and what to do about a
// failure, and Level 0 keeps its own ordering/accumulation.
//
// A nil *Owner is an uninitialised diagnostic system: every method is safe and
// does nothing (or returns an empty value), so a partially built App behaves
// exactly as the old nil-field checks did.
type Owner struct {
	logger   *Logger
	reporter *Reporter

	mu   sync.RWMutex
	port Emitter
}

// Open builds the owner's logging and reporting parts. The event transport is
// attached later, when the application has a live UI; logging works without it.
func Open(cfg Config) (*Owner, error) {
	logger, err := NewLogger(cfg.LogDir)
	if err != nil {
		return nil, err
	}
	reporter := NewReporter(logger, cfg.AppVersion)
	reporter.SetEnabled(cfg.UploadEnabled)
	return &Owner{logger: logger, reporter: reporter}, nil
}

// Attach installs the application's event transport. Attaching a new port
// replaces the old one; nil detaches emission (logging continues).
func (o *Owner) Attach(port Emitter) {
	if o == nil {
		return
	}
	o.mu.Lock()
	o.port = port
	o.mu.Unlock()
}

// Close flushes and closes the log file. Reporting after Close still reads the
// log but cannot append; a second Close is a no-op.
func (o *Owner) Close() error {
	if o == nil || o.logger == nil {
		return nil
	}
	return o.logger.Close()
}

// SetEnabled records the user's upload preference.
func (o *Owner) SetEnabled(enabled bool) {
	if o == nil || o.reporter == nil {
		return
	}
	o.reporter.SetEnabled(enabled)
}

// Enabled reports the current upload preference.
func (o *Owner) Enabled() bool {
	if o == nil || o.reporter == nil {
		return false
	}
	return o.reporter.IsEnabled()
}

// Entry is a flat view of one logged error, for the settings UI.
type Entry struct {
	Time      string
	Code      string
	Severity  string
	Operation string
	Location  string
	Message   string
	Context   map[string]string
}

// Entries returns the most recent n logged errors.
func (o *Owner) Entries(n int) []Entry {
	if o == nil || o.logger == nil {
		return nil
	}
	raw := o.logger.ReadRecent(n)
	out := make([]Entry, 0, len(raw))
	for _, e := range raw {
		out = append(out, Entry{
			Time:      e.Time,
			Code:      e.Code,
			Severity:  e.Severity,
			Operation: e.Operation,
			Location:  e.Location,
			Message:   e.Msg,
			Context:   e.Context,
		})
	}
	return out
}

// BuildReport builds a privacy-safe report from the most recent n errors.
func (o *Owner) BuildReport(n int) (*Report, error) {
	if o == nil || o.reporter == nil {
		return nil, errNotOpen
	}
	return o.reporter.BuildReport(n)
}

// Upload sends a report. The caller decides whether uploads are enabled.
func (o *Owner) Upload(report *Report) error {
	if o == nil || o.reporter == nil {
		return errNotOpen
	}
	return o.reporter.Upload(report)
}

// SaveReport writes a report beside the log for manual inspection.
func (o *Owner) SaveReport(report *Report) (string, error) {
	if o == nil || o.reporter == nil {
		return "", errNotOpen
	}
	return o.reporter.SaveReport(report)
}

// Report records a tracked error and, when the error is severe enough to
// surface, emits the user-facing event. Reporting a nil error does nothing.
func (o *Owner) Report(err *errors.TrackedError) {
	if o == nil || err == nil {
		return
	}
	if o.logger != nil {
		o.logger.Log(err)
	}
	if err.ShouldNotify() {
		o.emit("app:error", errorEventPayload(err))
	}
}

// ReportSilent records a tracked error without emitting a UI event.
func (o *Owner) ReportSilent(err *errors.TrackedError) {
	if o == nil || err == nil {
		return
	}
	if o.logger != nil {
		o.logger.Log(err)
	}
}

// ReportLLMLoad records a model-load failure (critical, surfaced).
func (o *Owner) ReportLLMLoad(err error, location string) {
	o.Report(errors.New(errors.LLMServerSpawnFail, errors.Critical, "llm.load", location, err))
}

// ReportLLMCrash records a llama-server crash (critical, surfaced).
func (o *Owner) ReportLLMCrash(err error, location string) {
	o.Report(errors.New(errors.LLMServerCrash, errors.Critical, "llm.load", location, err))
}

// ReportImageLoad records an image-model load failure (high, surfaced).
func (o *Owner) ReportImageLoad(err error, location string) {
	o.Report(errors.New(errors.IMGModelLoadFail, errors.High, "image.load", location, err))
}

// ReportImageGen records an image-generation failure (high, surfaced).
func (o *Owner) ReportImageGen(err error, location string) {
	o.Report(errors.New(errors.IMGGenerateFail, errors.High, "image.generate", location, err))
}

// ReportTTSSynth records a synthesis failure (medium, logged only).
func (o *Owner) ReportTTSSynth(err error, location string) {
	o.Report(errors.New(errors.TTSSynthesizeFail, errors.Medium, "tts.synthesize", location, err))
}

// ReportChatSave records a chat persistence failure (high, surfaced).
func (o *Owner) ReportChatSave(err error, location string) {
	o.Report(errors.New(errors.CHATSaveFail, errors.High, "chat.save", location, err))
}

// ReportResourceAlloc records a resource allocation failure (high, surfaced).
func (o *Owner) ReportResourceAlloc(err error, location string) {
	o.Report(errors.New(errors.RESAllocFail, errors.High, "resource.allocate", location, err))
}

// ReportPanic records a recovered panic (critical, surfaced).
func (o *Owner) ReportPanic(operation, location string, recovered interface{}) {
	o.Report(errors.New(errors.SYSPanicRecovered, errors.Critical, operation, location, fmt.Errorf("panic: %v", recovered)))
}

// errorEventPayload is the exact payload the desktop event has always carried:
// code, severity, operation, location, time, and sanitized context keys (never
// values), marshalled to a JSON string.
func errorEventPayload(err *errors.TrackedError) string {
	payload := map[string]interface{}{
		"code":      string(err.Code),
		"severity":  err.Severity.String(),
		"operation": err.Operation,
		"location":  err.Location,
		"time":      err.Time.Format(time.RFC3339),
	}

	ctxKeys := make([]string, 0, len(err.Context))
	for k := range err.Context {
		ctxKeys = append(ctxKeys, k)
	}
	if len(ctxKeys) > 0 {
		payload["context_keys"] = ctxKeys
	}

	data, _ := json.Marshal(payload)
	return string(data)
}

// emit sends through the attached port. A nil port drops the event; a port that
// fails mid-delivery must not take the reporting path down with it, so the
// failure is contained here.
func (o *Owner) emit(event string, data ...interface{}) {
	o.mu.RLock()
	port := o.port
	o.mu.RUnlock()
	if port == nil {
		return
	}
	defer func() { _ = recover() }()
	port.Emit(event, data...)
}
