package main

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"howl-chat/internal/backend/chat"
	"howl-chat/internal/diagnostics"
	apperrors "howl-chat/internal/errors"
	"howl-chat/internal/image/engine"
	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/media"
	appruntime "howl-chat/internal/runtime"
)

type recorderPort struct{ events []string }

func (p *recorderPort) Emit(event string, _ ...interface{}) { p.events = append(p.events, event) }

// A build without a diagnostics owner keeps the old refusals: the bindings must
// not panic and must answer exactly as the nil-field checks did.
func TestDiagnosticsBindingsRefuseWhenTheOwnerIsAbsent(t *testing.T) {
	app := &App{}
	if entries := app.GetErrorLog(10); entries != nil {
		t.Fatalf("GetErrorLog = %v, want nil", entries)
	}
	app.SetErrorReportingEnabled(true) // must not panic or write
	if app.IsErrorReportingEnabled() {
		t.Fatal("IsErrorReportingEnabled = true without an owner")
	}
	if _, err := app.SendErrorReport(); err == nil || err.Error() != "error reporter not initialized" {
		t.Fatalf("SendErrorReport error = %v", err)
	}
	if _, err := app.GetErrorReportPreview(); err == nil || err.Error() != "error reporter not initialized" {
		t.Fatalf("GetErrorReportPreview error = %v", err)
	}
}

func TestDiagnosticsOwnerBacksTheBindings(t *testing.T) {
	dir := t.TempDir()
	owner, err := diagnostics.Open(diagnostics.Config{LogDir: dir, AppVersion: "test"})
	if err != nil {
		t.Fatal(err)
	}
	defer owner.Close()
	port := &recorderPort{}
	owner.Attach(port)

	app := &App{diagnostics: owner}
	app.logLLMLoadError(stderrors.New("cannot spawn"), "test")

	entries := app.GetErrorLog(10)
	if len(entries) != 1 || entries[0]["code"] != string(apperrors.LLMServerSpawnFail) {
		t.Fatalf("GetErrorLog = %v", entries)
	}
	if len(port.events) != 1 || port.events[0] != "app:error" {
		t.Fatalf("events = %v, want one app:error", port.events)
	}

	preview, err := app.GetErrorReportPreview()
	if err != nil {
		t.Fatalf("GetErrorReportPreview: %v", err)
	}
	if !strings.Contains(preview, string(apperrors.LLMServerSpawnFail)) {
		t.Fatalf("preview does not contain the logged code: %s", preview)
	}

	app.SetErrorReportingEnabled(true)
	if !app.IsErrorReportingEnabled() {
		t.Fatal("toggle did not reach the owner")
	}
}

// The emitter is transport, not UI: a host with a sink but no Wails context
// (the browser harness) must still receive diagnostics events.
func TestBindingTheEmitterNeedsNoWebviewContext(t *testing.T) {
	owner, err := diagnostics.Open(diagnostics.Config{LogDir: t.TempDir(), AppVersion: "test"})
	if err != nil {
		t.Fatal(err)
	}
	defer owner.Close()

	var got []string
	previous := devEventSink
	devEventSink = func(name string, _ ...interface{}) { got = append(got, name) }
	defer func() { devEventSink = previous }()

	app := &App{diagnostics: owner}
	app.bindErrorEmitter()
	app.logLLMLoadError(stderrors.New("boom"), "test")

	if len(got) != 1 || got[0] != "app:error" {
		t.Fatalf("sink events = %v, want one app:error without a ctx", got)
	}
}

// failingImageEngine makes the "image engine shutdown" stage fail so the
// application-supplied sink has a real shutdown failure to escalate.
type failingImageEngine struct{}

func (failingImageEngine) Initialize(imagetypes.EngineConfig) error { return nil }
func (failingImageEngine) ListModels() ([]imagetypes.ModelInfo, error) {
	return nil, nil
}
func (failingImageEngine) ListLoRAs() ([]string, error)                { return nil, nil }
func (failingImageEngine) LoadModel(imagetypes.ModelLoadOptions) error { return nil }
func (failingImageEngine) IsRunning() bool                             { return false }
func (failingImageEngine) EnsureReady(string, imagetypes.ModelLoadOptions) error {
	return nil
}
func (failingImageEngine) GenerateTxt2Img(imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	return nil, nil
}
func (failingImageEngine) GenerateImg2Img(imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	return nil, nil
}
func (failingImageEngine) GenerateInpaint(imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	return nil, nil
}
func (failingImageEngine) GenerateVideo(imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	return nil, nil
}
func (failingImageEngine) GetProgress() (imagetypes.ProgressInfo, error) {
	return imagetypes.ProgressInfo{}, nil
}
func (failingImageEngine) GetLoadedModel() *imagetypes.ModelInfo { return nil }
func (failingImageEngine) CancelJob(string) error                { return nil }
func (failingImageEngine) Shutdown() error                       { return stderrors.New("shutdown boom") }

var _ engine.ImageEngine = failingImageEngine{}

func TestShutdownFailureReachesTheApplicationSuppliedSink(t *testing.T) {
	owner, err := diagnostics.Open(diagnostics.Config{LogDir: t.TempDir(), AppVersion: "test"})
	if err != nil {
		t.Fatal(err)
	}
	images := media.NewImage("")
	images.SetEngine(failingImageEngine{})

	app := &App{
		diagnostics:  owner,
		images:       images,
		engine:       appruntime.NewEngine(nil, nil),
		conversation: chat.NewActiveConversation(nil),
	}

	if err := app.shutdown(context.Background()); err == nil {
		t.Fatal("shutdown reported success despite a failing stage")
	}

	entries := app.GetErrorLog(20)
	found := false
	for _, entry := range entries {
		if entry["code"] != string(apperrors.SYSShutdownFail) {
			continue
		}
		found = true
		if entry["operation"] != "app.shutdown" {
			t.Fatalf("shutdown entry operation = %v", entry["operation"])
		}
		if entry["location"] != "app-lifecycle.go:shutdownStages" {
			t.Fatalf("shutdown entry location = %v", entry["location"])
		}
	}
	if !found {
		t.Fatalf("no SYS-SHUTDOWN entry in the diagnostics log: %v", entries)
	}
}
