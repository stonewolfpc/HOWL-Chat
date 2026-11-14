package main

import (
	"context"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"howl-chat/internal/audio"
	audiotypes "howl-chat/internal/audio/types"
)

type appFakePipeline struct {
	mu sync.Mutex

	status    audiotypes.PipelineStatus
	voices    []audiotypes.Voice
	getVoices func() []audiotypes.Voice
	process   func(context.Context, audiotypes.AudioMessage) (*audiotypes.AudioResponse, error)
	shutdown  error
	shutCalls int
}

func (f *appFakePipeline) ProcessAudioMessage(ctx context.Context, msg audiotypes.AudioMessage) (*audiotypes.AudioResponse, error) {
	if f.process != nil {
		return f.process(ctx, msg)
	}
	return &audiotypes.AudioResponse{}, nil
}

func (f *appFakePipeline) GetStatus() audiotypes.PipelineStatus {
	return f.status
}

func (f *appFakePipeline) GetTTSVoices() []audiotypes.Voice {
	f.mu.Lock()
	getVoices := f.getVoices
	voices := append([]audiotypes.Voice(nil), f.voices...)
	f.mu.Unlock()
	if getVoices != nil {
		return getVoices()
	}
	return voices
}

func (f *appFakePipeline) Shutdown(context.Context) error {
	f.mu.Lock()
	f.shutCalls++
	f.mu.Unlock()
	return f.shutdown
}

func (f *appFakePipeline) shutdownCalls() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.shutCalls
}

func appCandidate(pipeline *appFakePipeline, asrReady, ttsReady bool, asrErr, ttsErr error) *audio.Candidate {
	return &audio.Candidate{
		Pipeline: pipeline,
		ASR:      audio.CapabilityResult{Ready: asrReady, Err: asrErr},
		TTS:      audio.CapabilityResult{Ready: ttsReady, Err: ttsErr},
	}
}

func waitAppCondition(t *testing.T, what string, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(time.Millisecond)
	}
}

func desiredWithBackend(settings audio.Settings) audio.DesiredSettings {
	desired := audio.DefaultDesiredSettings()
	desired.ASREnabled = settings.ASREnabled
	desired.TTSEnabled = settings.TTSEnabled
	desired.ASREngine = settings.ASREngine
	desired.TTSEngine = settings.TTSEngine
	if settings.TTSSpeed != 0 {
		desired.TTSSpeed = settings.TTSSpeed
	}
	if settings.ASRLanguage != "" {
		desired.ASRLanguage = settings.ASRLanguage
	}
	if settings.ASRThreads != 0 {
		desired.ASRThreads = settings.ASRThreads
	}
	if settings.ASRBeamSize != 0 {
		desired.ASRBeamSize = settings.ASRBeamSize
	}
	if settings.ASRBestOf != 0 {
		desired.ASRBestOf = settings.ASRBestOf
	}
	return desired
}

func newAudioTestApp(t *testing.T, factory audio.PipelineFactory) (*App, *audio.Manager) {
	t.Helper()
	manager := audio.NewManager(factory)
	app := &App{audioManager: manager}
	t.Cleanup(func() {
		_ = manager.Close(context.Background())
	})
	return app, manager
}

// A1: App owns exactly one manager and no direct pipeline field.
func TestA1SingleManagerOwnership(t *testing.T) {
	appType := reflect.TypeOf(App{})
	if _, ok := appType.FieldByName("audioPipeline"); ok {
		t.Fatal("App still has an audioPipeline field")
	}
	field, ok := appType.FieldByName("audioManager")
	if !ok {
		t.Fatal("App has no audioManager field")
	}
	managerType := reflect.TypeOf((*audio.Manager)(nil))
	if field.Type != managerType {
		t.Fatalf("audioManager type = %v, want %v", field.Type, managerType)
	}

	managerFields := 0
	pipelineFields := 0
	pipelineType := reflect.TypeOf((*audio.Pipeline)(nil))
	for i := 0; i < appType.NumField(); i++ {
		fieldType := appType.Field(i).Type
		if fieldType == managerType {
			managerFields++
		}
		if fieldType == pipelineType {
			pipelineFields++
		}
	}
	if managerFields != 1 {
		t.Fatalf("manager fields = %d, want exactly 1", managerFields)
	}
	if pipelineFields != 0 {
		t.Fatalf("direct pipeline fields = %d, want 0", pipelineFields)
	}
}

// A3: enable uses Apply and a failed replacement preserves the old generation.
func TestA3TemporaryEnableBridgePreservesPreviousReadiness(t *testing.T) {
	useMicProbe(t, constantMicProbe(newMicrophoneStatus(MicStateReady, "", 1)))
	p1 := &appFakePipeline{status: audiotypes.PipelineStatus{ASRReady: true}}
	p2 := &appFakePipeline{}
	replacementErr := errors.New("replacement ASR failed")
	call := 0
	app, manager := newAudioTestApp(t, func(_ context.Context, desired audio.Settings, candidate *audio.Candidate) error {
		call++
		switch call {
		case 1:
			*candidate = *appCandidate(p1, true, false, nil, nil)
		case 2:
			*candidate = *appCandidate(p2, false, false, replacementErr, nil)
		default:
			t.Fatalf("unexpected factory call %d", call)
		}
		return nil
	})

	if err := app.EnableAudioASR(true); err != nil {
		t.Fatalf("initial EnableAudioASR: %v", err)
	}
	if err := app.EnableAudioASR(true); err == nil {
		t.Fatal("failed replacement returned nil")
	} else if !strings.Contains(err.Error(), replacementErr.Error()) {
		t.Fatalf("replacement error = %v, want %v", err, replacementErr)
	}

	status := manager.Status()
	if !status.ASR.Requested || !status.ASR.Ready || status.ASR.Error == "" {
		t.Fatalf("status after failed replacement = %+v, want requested old-ready with error", status.ASR)
	}
	if status.Phase != "failed" {
		t.Fatalf("phase after failed replacement = %q, want failed", status.Phase)
	}
	lease, err := manager.Acquire(context.Background())
	if err != nil {
		t.Fatalf("Acquire after failed replacement: %v", err)
	}
	if lease.Pipeline() != p1 {
		t.Fatal("failed replacement did not preserve the old generation")
	}
	lease.Release()
	if p2.shutdownCalls() != 1 {
		t.Fatalf("failed candidate shutdown calls = %d, want 1", p2.shutdownCalls())
	}
	if p1.shutdownCalls() != 0 {
		t.Fatalf("old generation shutdown calls = %d, want 0", p1.shutdownCalls())
	}

	if err := app.EnableAudioASR(false); err != nil {
		t.Fatalf("disable ASR: %v", err)
	}
	if status := manager.Status(); status.ASR.Requested || status.ASR.Ready {
		t.Fatalf("status after disable = %+v, want ASR unrequested and not ready", status.ASR)
	}
	waitAppCondition(t, "old generation shutdown after disable", func() bool {
		return p1.shutdownCalls() == 1
	})
}

// A4: status comes from the manager and voice enumeration holds one lease.
func TestA4StatusAndVoiceLease(t *testing.T) {
	voiceStarted := make(chan struct{})
	voiceRelease := make(chan struct{})
	pipeline := &appFakePipeline{
		status: audiotypes.PipelineStatus{TTSReady: true},
		voices: []audiotypes.Voice{{ID: "voice-a"}, {ID: "voice-b"}},
	}
	pipeline.getVoices = func() []audiotypes.Voice {
		close(voiceStarted)
		<-voiceRelease
		return []audiotypes.Voice{{ID: "voice-a"}, {ID: "voice-b"}}
	}

	app, manager := newAudioTestApp(t, func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
		*candidate = *appCandidate(pipeline, false, true, nil, nil)
		return nil
	})
	if _, _, err := manager.Apply(context.Background(), audio.Settings{TTSEnabled: true}); err != nil {
		t.Fatalf("Apply TTS: %v", err)
	}

	status := app.GetAudioStatus()
	if status["tts_ready"] != true || status["asr_ready"] != false {
		t.Fatalf("status = %+v, want TTS ready and ASR not ready", status)
	}
	if status["active"] != true {
		t.Fatalf("active = %v, want true", status["active"])
	}
	if status["active_sessions"] != 0 {
		t.Fatalf("active_sessions = %v, want 0", status["active_sessions"])
	}

	done := make(chan []string, 1)
	go func() { done <- app.GetTTSAvailableVoices() }()
	select {
	case <-voiceStarted:
	case <-time.After(3 * time.Second):
		t.Fatal("voice enumeration did not start")
	}
	if got := manager.Status().InFlight; got != 1 {
		t.Fatalf("in-flight during voice enumeration = %d, want one lease", got)
	}
	close(voiceRelease)
	var voices []string
	select {
	case voices = <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("voice enumeration did not finish")
	}
	if len(voices) != 2 || voices[0] != "voice-a" || voices[1] != "voice-b" {
		t.Fatalf("voices = %v, want [voice-a voice-b]", voices)
	}
	if got := manager.Status().InFlight; got != 0 {
		t.Fatalf("in-flight after voice enumeration = %d, want 0", got)
	}
}

// A5: TTS and legacy ASR each hold an operation and lease with no deadline.
func TestA5TTSAndLegacyASRUseOperationAndLeaseWithoutDeadline(t *testing.T) {
	tests := []struct {
		name string
		call func(*App) error
	}{
		{
			name: "TTS",
			call: func(app *App) error {
				_, err := app.SynthesizeSpeech("hello")
				return err
			},
		},
		{
			name: "legacy ASR",
			call: func(app *App) error {
				_, err := app.TranscribeAudio("clip.wav")
				return err
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			started := make(chan struct{})
			release := make(chan struct{})
			var deadline bool
			pipeline := &appFakePipeline{}
			pipeline.process = func(ctx context.Context, msg audiotypes.AudioMessage) (*audiotypes.AudioResponse, error) {
				_, deadline = ctx.Deadline()
				close(started)
				<-release
				if msg.Type == "audio_input" {
					return &audiotypes.AudioResponse{TranscribedText: "text"}, nil
				}
				return &audiotypes.AudioResponse{AudioOutputPath: "speech.wav"}, nil
			}

			app, manager := newAudioTestApp(t, func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
				*candidate = *appCandidate(pipeline, true, true, nil, nil)
				return nil
			})
			if _, _, err := manager.Apply(context.Background(), audio.Settings{ASREnabled: true, TTSEnabled: true}); err != nil {
				t.Fatalf("Apply: %v", err)
			}

			done := make(chan error, 1)
			go func() { done <- tc.call(app) }()
			select {
			case <-started:
			case <-time.After(3 * time.Second):
				t.Fatal("request did not start")
			}
			if got := manager.Status().InFlight; got != 2 {
				t.Fatalf("in-flight during request = %d, want operation + lease = 2", got)
			}
			if deadline {
				t.Fatal("request context has an artificial deadline")
			}
			close(release)
			select {
			case err := <-done:
				if err != nil {
					t.Fatalf("request: %v", err)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("request did not finish")
			}
			if got := manager.Status().InFlight; got != 0 {
				t.Fatalf("in-flight after request = %d, want 0", got)
			}
		})
	}
}

// A6: OnShutdown closes the manager and reports a shutdown error.
func TestA6OnShutdownReportsManagerCloseError(t *testing.T) {
	shutdownErr := errors.New("shutdown failed")
	pipeline := &appFakePipeline{status: audiotypes.PipelineStatus{ASRReady: true}, shutdown: shutdownErr}
	app, manager := newAudioTestApp(t, func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
		*candidate = *appCandidate(pipeline, true, false, nil, nil)
		return nil
	})
	if _, _, err := manager.Apply(context.Background(), audio.Settings{ASREnabled: true}); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	oldStdout := os.Stdout
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = write
	app.OnShutdown(context.Background())
	_ = write.Close()
	os.Stdout = oldStdout
	defer func() { os.Stdout = oldStdout }()
	output, err := io.ReadAll(read)
	if err != nil {
		t.Fatalf("read shutdown output: %v", err)
	}
	_ = read.Close()

	if !strings.Contains(string(output), "Audio manager shutdown error") || !strings.Contains(string(output), shutdownErr.Error()) {
		t.Fatalf("shutdown output = %q, want reported manager error", output)
	}
	if status := manager.Status(); status.Phase != "closed_with_errors" {
		t.Fatalf("phase = %q, want closed_with_errors", status.Phase)
	}
	if pipeline.shutdownCalls() != 1 {
		t.Fatalf("pipeline shutdown calls = %d, want 1", pipeline.shutdownCalls())
	}
}

// A7: Close's default five-second timeout preserves the lease and leaves the
// manager closing; a later Close finalizes after the lease is released.
func TestA7FiveSecondCloseTimeout(t *testing.T) {
	pipeline := &appFakePipeline{status: audiotypes.PipelineStatus{ASRReady: true}}
	manager := audio.NewManager(func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
		*candidate = *appCandidate(pipeline, true, false, nil, nil)
		return nil
	})
	if _, _, err := manager.Apply(context.Background(), audio.Settings{ASREnabled: true}); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	lease, err := manager.Acquire(context.Background())
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}

	start := time.Now()
	result := make(chan error, 1)
	go func() {
		result <- manager.Close(context.Background())
	}()

	var closeErr error
	select {
	case closeErr = <-result:
	case <-time.After(7 * time.Second):
		t.Fatal("Close did not return after the default timeout")
	}
	elapsed := time.Since(start)
	if elapsed < 4*time.Second {
		t.Fatalf("Close returned after %v, want approximately five seconds", elapsed)
	}
	if !errors.Is(closeErr, audio.ErrDrainTimeout) {
		t.Fatalf("Close error = %v, want ErrDrainTimeout", closeErr)
	}
	if status := manager.Status(); status.Phase != "closing" {
		t.Fatalf("phase after timeout = %q, want closing", status.Phase)
	}
	if pipeline.shutdownCalls() != 0 {
		t.Fatalf("pipeline shutdown calls while lease held = %d, want 0", pipeline.shutdownCalls())
	}

	lease.Release()
	waitAppCondition(t, "pipeline shutdown after lease release", func() bool {
		return pipeline.shutdownCalls() == 1
	})
	closeCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := manager.Close(closeCtx); err != nil {
		t.Fatalf("later Close: %v", err)
	}
	if status := manager.Status(); status.Phase != "closed" {
		t.Fatalf("phase after later Close = %q, want closed", status.Phase)
	}
	if pipeline.shutdownCalls() != 1 {
		t.Fatalf("pipeline shutdown calls = %d, want exactly 1", pipeline.shutdownCalls())
	}
}
