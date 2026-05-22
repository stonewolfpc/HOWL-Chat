package errors

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ---------------------------------------------------------------------------
// WailsEmitter - Bridges TrackedErrors to the frontend via Wails events.
// Emits minimal data (code + severity + operation). No stack traces.
// ---------------------------------------------------------------------------

type WailsEmitter struct {
	ctx context.Context
}

// NewWailsEmitter creates an emitter bound to the Wails context.
func NewWailsEmitter(ctx context.Context) *WailsEmitter {
	return &WailsEmitter{ctx: ctx}
}

// Emit sends an error to the frontend. Only critical/high severity surface.
func (e *WailsEmitter) Emit(err *TrackedError) {
	if e.ctx == nil || err == nil {
		return
	}
	if !err.ShouldNotify() {
		return
	}

	payload := map[string]interface{}{
		"code":      string(err.Code),
		"severity":  err.Severity.String(),
		"operation": err.Operation,
		"location":  err.Location,
		"time":      err.Time.Format(time.RFC3339),
	}

	// Include sanitized context keys only (no values)
	var ctxKeys []string
	for k := range err.Context {
		ctxKeys = append(ctxKeys, k)
	}
	if len(ctxKeys) > 0 {
		payload["context_keys"] = ctxKeys
	}

	data, _ := json.Marshal(payload)
	runtime.EventsEmit(e.ctx, "app:error", string(data))
}

// EmitRecovery notifies the frontend that a panic was recovered.
func (e *WailsEmitter) EmitRecovery(operation string) {
	if e.ctx == nil {
		return
	}
	runtime.EventsEmit(e.ctx, "app:recovered", map[string]string{
		"operation": operation,
		"message":   "A system error was recovered. Functionality may be degraded.",
	})
}

// EmitSilent sends an event that does not show UI but is logged.
func (e *WailsEmitter) EmitSilent(err *TrackedError) {
	if e.ctx == nil || err == nil {
		return
	}
	runtime.EventsEmit(e.ctx, "app:error:silent", map[string]string{
		"code":     string(err.Code),
		"severity": err.Severity.String(),
	})
}

// CombinedSink is a convenience type that logs AND emits in one call.
type CombinedSink struct {
	logger  *Logger
	emitter *WailsEmitter
}

// NewCombinedSink creates a sink that writes to file and emits to frontend.
func NewCombinedSink(logger *Logger, emitter *WailsEmitter) *CombinedSink {
	return &CombinedSink{logger: logger, emitter: emitter}
}

// Handle processes a TrackedError: logs to file, emits to frontend if severe.
func (s *CombinedSink) Handle(err *TrackedError) {
	if err == nil {
		return
	}
	if s.logger != nil {
		s.logger.Log(err)
	}
	if s.emitter != nil {
		s.emitter.Emit(err)
	}
}

// HandleSilent logs to file only, never emits to frontend.
func (s *CombinedSink) HandleSilent(err *TrackedError) {
	if err == nil {
		return
	}
	if s.logger != nil {
		s.logger.Log(err)
	}
}

// ---------------------------------------------------------------------------
// Convenience constructors for common error patterns
// ---------------------------------------------------------------------------

func (s *CombinedSink) LogLLMLoad(err error, location string) {
	s.Handle(New(LLMServerSpawnFail, Critical, "llm.load", location, err))
}

func (s *CombinedSink) LogLLMCrash(err error, location string) {
	s.Handle(New(LLMServerCrash, Critical, "llm.load", location, err))
}

func (s *CombinedSink) LogImageLoad(err error, location string) {
	s.Handle(New(IMGModelLoadFail, High, "image.load", location, err))
}

func (s *CombinedSink) LogImageGen(err error, location string) {
	s.Handle(New(IMGGenerateFail, High, "image.generate", location, err))
}

func (s *CombinedSink) LogTTSSynth(err error, location string) {
	s.Handle(New(TTSSynthesizeFail, Medium, "tts.synthesize", location, err))
}

func (s *CombinedSink) LogChatSave(err error, location string) {
	s.Handle(New(CHATSaveFail, High, "chat.save", location, err))
}

func (s *CombinedSink) LogResourceAlloc(err error, location string) {
	s.Handle(New(RESAllocFail, High, "resource.allocate", location, err))
}

func (s *CombinedSink) LogSystemPanic(operation, location string, r interface{}) {
	s.Handle(New(SYSPanicRecovered, Critical, operation, location, fmt.Errorf("panic: %v", r)))
}
