package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"howl-chat/internal/audio"
	"howl-chat/internal/audio/capability"
	audiotypes "howl-chat/internal/audio/types"
	"howl-chat/internal/backend/rag"
	"howl-chat/internal/errors"
)

// ---------------------------------------------------------------------------
// Audio Wails Bindings
// ---------------------------------------------------------------------------

// beginManagedAudio registers one operation and acquires one lease. The
// operation covers the request lifetime; the lease pins the generation even if
// a replacement is applied while the request is running.
func (a *App) beginManagedAudio(ctx context.Context) (*audio.Operation, context.Context, *audio.Lease, audio.ManagedPipeline, error) {
	if err := a.ensureAcceptingWork("audio operation"); err != nil {
		return nil, nil, nil, nil, err
	}
	if a.audioManager == nil {
		return nil, nil, nil, nil, fmt.Errorf("audio manager not initialized")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	op, opCtx := a.audioManager.BeginOperation(ctx, 0)
	lease, err := a.audioManager.Acquire(opCtx)
	if err != nil {
		op.End()
		return nil, nil, nil, nil, err
	}
	return op, opCtx, lease, lease.Pipeline(), nil
}

// processManagedAudio runs one audio request under a manager operation and
// lease. T2.9.5A deliberately adds no request deadline here.
func (a *App) processManagedAudio(ctx context.Context, msg audiotypes.AudioMessage) (*audiotypes.AudioResponse, error) {
	op, opCtx, lease, pipeline, err := a.beginManagedAudio(ctx)
	if err != nil {
		return nil, err
	}
	defer op.End()
	defer lease.Release()

	if pipeline == nil {
		return nil, fmt.Errorf("audio pipeline unavailable")
	}
	return pipeline.ProcessAudioMessage(opCtx, msg)
}

const (
	// maxTranscribedAudioBytes bounds one renderer payload before it is written
	// to disk. Two minutes of mono 16 kHz PCM16 is about 3.8 MiB; the larger
	// limit leaves room for other valid rates without allowing unbounded input.
	maxTranscribedAudioBytes = 64 << 20
	audioTempDirName         = "howl-asr"
)

// transcribeAudioPath runs one path-based request through the managed pipeline.
func (a *App) transcribeAudioPath(ctx context.Context, audioPath string) (string, error) {
	result, err := a.processManagedAudio(ctx, audiotypes.AudioMessage{
		Type:      "audio_input",
		AudioPath: audioPath,
	})
	if err != nil {
		a.logError(errors.ASRTranscribeFail, errors.High, "asr.transcribe", "app-audio.go:TranscribeAudio", err)
		return "", err
	}
	return result.TranscribedText, nil
}

// TranscribeAudio converts an audio file to text via the ASR pipeline. The
// microphone path uses TranscribeAudioData so the renderer never needs to
// create a filesystem path of its own.
func (a *App) TranscribeAudio(audioPath string) (string, error) {
	return a.transcribeAudioPath(context.Background(), audioPath)
}

// TranscribeAudioData accepts a base64 WAV payload from the renderer, writes
// one uniquely named temporary file, transcribes it, and removes the file on
// success, error, cancellation, or panic. The renderer cannot supply a path.
func (a *App) TranscribeAudioData(audioDataB64 string) (string, error) {
	data, err := decodeAudioBase64(audioDataB64)
	if err != nil {
		return "", err
	}
	if !isWAVData(data) {
		return "", fmt.Errorf("audio data must be a WAV file")
	}
	return transcribeTempAudio(context.Background(), data, a.transcribeAudioPath)
}

func decodeAudioBase64(encoded string) ([]byte, error) {
	trimmed := strings.TrimSpace(encoded)
	if idx := strings.Index(trimmed, ";base64,"); idx >= 0 {
		trimmed = trimmed[idx+len(";base64,"):]
	}
	if trimmed == "" {
		return nil, fmt.Errorf("audio data is empty")
	}
	if len(trimmed) > base64.StdEncoding.EncodedLen(maxTranscribedAudioBytes) {
		return nil, fmt.Errorf("audio data exceeds %d bytes", maxTranscribedAudioBytes)
	}
	data, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return nil, fmt.Errorf("decode audio data: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("audio data is empty")
	}
	if len(data) > maxTranscribedAudioBytes {
		return nil, fmt.Errorf("audio data exceeds %d bytes", maxTranscribedAudioBytes)
	}
	return data, nil
}

func isWAVData(data []byte) bool {
	return len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WAVE"
}

// writeTempAudio creates a unique file in the app's ASR temp directory. The
// returned cleanup function is safe to call more than once.
func writeTempAudio(data []byte) (string, func(), error) {
	dir := filepath.Join(isolatedTempDir(), audioTempDirName)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", nil, fmt.Errorf("create ASR temp dir: %w", err)
	}
	file, err := os.CreateTemp(dir, "recording-*.wav")
	if err != nil {
		return "", nil, fmt.Errorf("create ASR temp file: %w", err)
	}
	path := file.Name()
	cleanup := func() { _ = os.Remove(path) }
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		cleanup()
		return "", nil, fmt.Errorf("write ASR temp file: %w", err)
	}
	if err := file.Close(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("close ASR temp file: %w", err)
	}
	return path, cleanup, nil
}

// transcribeTempAudio is separated from the Wails binding so the lifecycle can
// be tested with success, error, and cancellation outcomes.
func transcribeTempAudio(ctx context.Context, data []byte, transcribe func(context.Context, string) (string, error)) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if transcribe == nil {
		return "", fmt.Errorf("transcription callback is nil")
	}
	path, cleanup, err := writeTempAudio(data)
	if err != nil {
		return "", err
	}
	defer cleanup()
	return transcribe(ctx, path)
}

// SynthesizeSpeech converts text to speech and returns the path to the generated audio file.
func (a *App) SynthesizeSpeech(text string) (string, error) {
	result, err := a.processManagedAudio(context.Background(), audiotypes.AudioMessage{
		Type: "audio_output",
		Text: text,
	})
	if err != nil {
		a.logError(errors.TTSSynthesizeFail, errors.High, "tts.synthesize", "app-audio.go:SynthesizeSpeech", err)
		return "", err
	}
	return result.AudioOutputPath, nil
}

// SynthesizeSpeechWithVoice converts text to speech using a per-request voice override.
func (a *App) SynthesizeSpeechWithVoice(text string, voiceJSON string) (string, error) {
	override, err := parseVoiceOverrideJSON(voiceJSON)
	if err != nil {
		a.logError(errors.TTSSynthesizeFail, errors.High, "tts.synthesize.voice", "app-audio.go:SynthesizeSpeechWithVoice", err)
		return "", err
	}
	msg := audiotypes.AudioMessage{
		Type:          "audio_output",
		Text:          text,
		VoiceOverride: override,
	}
	result, err := a.processManagedAudio(context.Background(), msg)
	if err != nil {
		a.logError(errors.TTSSynthesizeFail, errors.High, "tts.synthesize.voice", "app-audio.go:SynthesizeSpeechWithVoice", err)
		return "", err
	}
	return result.AudioOutputPath, nil
}

// parseVoiceOverrideJSON preserves malformed requests as errors rather than
// silently asking the active voice to speak. An empty payload means no
// override; a non-empty payload is always an explicit request, including an
// inheritance selection with blank engine/voice fields.
func parseVoiceOverrideJSON(voiceJSON string) (*audiotypes.VoiceOverride, error) {
	if strings.TrimSpace(voiceJSON) == "" {
		return nil, nil
	}
	var override audiotypes.VoiceOverride
	if err := json.Unmarshal([]byte(voiceJSON), &override); err != nil {
		return nil, fmt.Errorf("parse voice override JSON: %w", err)
	}
	return &override, nil
}

// GetAudioStatus returns the current manager snapshot in the temporary 5A
// shape expected by the existing frontend.
func (a *App) GetAudioStatus() map[string]interface{} {
	if a.audioManager == nil {
		return map[string]interface{}{
			"asr_ready":         false,
			"asr_readiness":     "not-ready",
			"asr_error":         "audio manager not initialized",
			"tts_ready":         false,
			"tts_error":         "audio manager not initialized",
			"multimodal_ready":  false,
			"realtime_ready":    false,
			"active_sessions":   int64(0),
			"error_count":       int64(0),
			"active":            false,
			"apply_in_progress": false,
			"applied":           nil,
			"desired_gen":       uint64(0),
			"applied_gen":       uint64(0),
			// The contract is available before the manager: it reports what
			// the build and the installed assets can do, not what a running
			// generation happens to have loaded.
			"capabilities": audioCapabilityRecords(),
			"mic":          microphoneStatusRecord(microphoneStatusNow()),
			"settings":     audio.BuildSettingsReport(a.currentAudioDesired(), audio.Settings{}, false, audio.Status{}, "audio manager not initialized"),
		}
	}
	status := a.audioManager.Status()
	errorCount := 0
	if status.ASR.Error != "" {
		errorCount++
	}
	if status.TTS.Error != "" {
		errorCount++
	}
	mic := microphoneStatusNow()
	readiness, readinessDetail := a.audioReadiness(status, mic)
	active := status.Phase != "closing" && status.Phase != "closed" && status.Phase != "closed_with_errors"

	// apply_in_progress is the manager's live build gauge (Status.Builds). It is
	// an active-operation count, not a cumulative counter.
	applyInProgress := status.Builds > 0

	// applied is a compatibility projection of the ACTIVE generation's four
	// capability fields. It is not the settings contract and must never be used
	// as proof that the full desired configuration was applied; settings is the
	// authoritative per-setting report.
	installed, installedOK := a.audioManager.AppliedSettings()
	var applied interface{}
	if installedOK {
		applied = map[string]interface{}{
			"asr_enabled": installed.ASREnabled,
			"tts_enabled": installed.TTSEnabled,
			"asr_engine":  installed.ASREngine,
			"tts_engine":  installed.TTSEngine,
		}
	}
	settingsReport := audio.BuildSettingsReport(a.currentAudioDesired(), installed, installedOK, status, a.audioLastApplyError())
	return map[string]interface{}{
		// asr_ready is limited to locally initialized engines. A remote engine
		// can be applied successfully and still have no locally proven readiness;
		// asr_readiness carries that distinction instead of turning Apply into
		// an availability claim.
		"asr_ready":        status.ASR.Ready && readiness == "local-initialized",
		"asr_readiness":    readiness,
		"asr_error":        readinessDetail,
		"tts_ready":        status.TTS.Ready,
		"tts_error":        status.TTS.Error,
		"multimodal_ready": false,
		"realtime_ready":   false,
		"active_sessions":  status.InFlight,
		"error_count":      int64(errorCount),
		"active":           active,
		// capabilities is the provider-neutral contract: one record per task,
		// with every dimension and the reason each blocked option is blocked.
		// Status and the settings UI read the same list, so they cannot
		// disagree about what is offered.
		"capabilities":      audioCapabilityRecords(),
		"mic":               microphoneStatusRecord(mic),
		"apply_in_progress": applyInProgress,
		"applied":           applied,
		"settings":          settingsReport,
		// Observability only: a queued caller advances DesiredGen before it
		// waits on applyMu, so equality with AppliedGen shows no outstanding
		// generation. These fields never authorize releasing a lock.
		"desired_gen": status.DesiredGen,
		"applied_gen": status.AppliedGen,
	}
}

// audioReadiness reports the strongest ASR state this build has established.
// Local adapters must be initialized; remote adapters are never presented as
// locally ready merely because their client accepted configuration.
func (a *App) audioReadiness(status audio.Status, mic MicrophoneStatus) (string, string) {
	if !status.ASR.Requested && !a.asrActivationPending.Load() {
		return "disabled", "ASR is disabled"
	}
	// Capture availability is checked before engine readiness: a loaded model
	// with no usable microphone is not ready for speech input.
	if !mic.Usable() {
		return "mic-unavailable", mic.Message
	}
	if a.audioManager == nil {
		return "not-ready", "audio manager not initialized"
	}
	engine := strings.TrimSpace(a.audioManager.Desired().ASREngine)
	if engine == "" || engine == "auto" {
		return "unconfigured", "select an explicit ASR engine"
	}
	profile, ok := audiotypes.GetASRProfile(audiotypes.ASRModelType(engine))
	if !ok {
		return "not-ready", fmt.Sprintf("unknown ASR engine %q", engine)
	}
	if !profile.Local || profile.Backend == "http" {
		detail := "remote engine availability is not locally verified"
		if status.ASR.Error != "" {
			detail = status.ASR.Error
		}
		return "remote-unverified", detail
	}
	if status.ASR.Ready {
		return "local-initialized", status.ASR.Error
	}
	if status.ASR.Error != "" {
		return "not-ready", status.ASR.Error
	}
	return "not-ready", "the selected local ASR engine is not initialized"
}

// audioCapabilityRecords builds the capability contract from the same path
// policies the pipeline loads with.
func audioCapabilityRecords() []capability.Record {
	return capability.Catalog(capability.ProductionEnv(audio.ASRModelPathFor, audio.ASRProjectorPathFor, audio.TTSModelPathFor))
}

// GetTTSAvailableVoices returns the list of available voice names from the
// current TTS engine. Voice enumeration only needs a generation lease.
func (a *App) GetTTSAvailableVoices() []string {
	if a.audioManager == nil {
		return nil
	}
	lease, err := a.audioManager.Acquire(context.Background())
	if err != nil {
		return nil
	}
	defer lease.Release()

	pipeline := lease.Pipeline()
	if pipeline == nil {
		return nil
	}
	voices := pipeline.GetTTSVoices()
	result := make([]string, 0, len(voices))
	for _, v := range voices {
		result = append(result, v.ID)
	}
	return result
}

// ApplyAudioSettings applies the complete desired audio contract. The
// operational fields reach audioManager; every other field is retained and
// reported through GetAudioStatus with an explicit owner/state. It is the
// bridge the settings UI calls when the user explicitly saves or retries;
// persisted selections are restored into the panel but never applied at
// startup. An explicit engine that cannot run is refused with its reason and
// never replaced by another provider. Storage remains the settings store's job.
func (a *App) ApplyAudioSettings(desired audio.DesiredSettings) error {
	return a.applyAudioSettings(desired)
}

func (a *App) setAudioDesired(desired audio.DesiredSettings) {
	a.audioDesiredMu.Lock()
	a.audioDesired = desired
	a.audioDesiredSet = true
	a.audioDesiredMu.Unlock()
}

func (a *App) currentAudioDesired() audio.DesiredSettings {
	a.audioDesiredMu.RLock()
	defer a.audioDesiredMu.RUnlock()
	if !a.audioDesiredSet {
		return audio.DefaultDesiredSettings()
	}
	return a.audioDesired
}

func (a *App) setAudioApplyError(message string) {
	a.audioDesiredMu.Lock()
	a.audioApplyErr = message
	a.audioDesiredMu.Unlock()
}

func (a *App) audioLastApplyError() string {
	a.audioDesiredMu.RLock()
	defer a.audioDesiredMu.RUnlock()
	return a.audioApplyErr
}

// applyAudioSettings is the explicit activation path. The requested full
// contract is stored before validation so status can explain exactly which
// values were rejected. Only settings.Backend() is eligible for adoption by the
// manager; unsupported fields never become Applied merely because they were
// saved.
//
// ASR activation is gated on Windows capture availability. The probe runs before
// the manager builds a generation, so a missing or unusable microphone cannot
// allocate a model or reserve VRAM, and it runs again after the build so a device
// that disappears mid-initialization is torn down instead of left loaded.
func (a *App) applyAudioSettings(desired audio.DesiredSettings) error {
	if err := a.ensureAcceptingWork("audio settings"); err != nil {
		return err
	}
	a.setAudioDesired(desired)
	issues := desired.Validate()
	if audio.HasBackendValidationIssue(desired, issues) {
		err := &audio.SettingsValidationError{Issues: issues}
		a.setAudioApplyError(err.Error())
		return err
	}
	if a.audioManager == nil {
		err := fmt.Errorf("audio manager not initialized")
		a.setAudioApplyError(err.Error())
		return err
	}

	// A previous failure belongs to the previous attempt. Clear it only after
	// the new request has been recorded and the manager is available.
	a.setAudioApplyError("")
	settings := desired.Backend()

	if !settings.ASREnabled {
		a.asrActivationPending.Store(false)
		_, _, err := a.audioManager.Apply(context.Background(), settings)
		if err != nil {
			a.setAudioApplyError(err.Error())
			a.logError(errors.ASRModelLoadFail, errors.High, "audio.apply", "app-audio.go:applyAudioSettings",
				fmt.Errorf("apply ASR engine %q: %w", settings.ASREngine, err))
		}
		return err
	}

	// From here on the user has explicitly asked for ASR. Keep the request
	// visible to status until a generation is both applied and still backed by
	// a usable microphone.
	a.asrActivationPending.Store(true)

	mic := microphoneStatusNow()
	if !mic.Usable() {
		err := &MicrophoneUnavailableError{Status: mic}
		a.setAudioApplyError(err.Error())
		a.logError(errors.ASRMicUnavailable, errors.High, "audio.apply.mic", "app-audio.go:applyAudioSettings", err)
		return err
	}

	if _, _, err := a.audioManager.Apply(context.Background(), settings); err != nil {
		a.setAudioApplyError(err.Error())
		a.logError(errors.ASRModelLoadFail, errors.High, "audio.apply", "app-audio.go:applyAudioSettings",
			fmt.Errorf("apply ASR engine %q: %w", settings.ASREngine, err))
		return err
	}

	after := microphoneStatusNow()
	if !after.Usable() {
		// The microphone disappeared while the model was loading. Retire the ASR
		// capability in the same manager so no generation keeps a loaded model
		// with no capture path. TTS, if requested, is preserved.
		ttsOnly := settings
		ttsOnly.ASREnabled = false
		if _, _, abortErr := a.audioManager.Apply(context.Background(), ttsOnly); abortErr != nil {
			a.logError(errors.ASRMicUnavailable, errors.High, "audio.apply.mic.abort",
				"app-audio.go:applyAudioSettings", abortErr)
		}
		a.asrActivationPending.Store(true)
		err := &MicrophoneUnavailableError{Status: after}
		a.setAudioApplyError(err.Error())
		a.logError(errors.ASRMicUnavailable, errors.High, "audio.apply.mic.lost", "app-audio.go:applyAudioSettings", err)
		return err
	}
	a.asrActivationPending.Store(false)
	a.setAudioApplyError("")
	return nil
}

// EnableAudioASR toggles the ASR recognizer without discarding the desired TTS
// engine or the persisted non-operational settings. The previous implementation
// rebuilt the settings from one boolean, which silently reset an explicit engine
// chosen in the settings panel.
func (a *App) EnableAudioASR(enabled bool) error {
	desired := a.currentAudioDesired()
	desired.ASREnabled = enabled
	return a.applyAudioSettings(desired)
}

// EnableAudioTTS toggles the TTS synthesizer without discarding the desired
// ASR engine or the persisted non-operational settings.
func (a *App) EnableAudioTTS(enabled bool) error {
	desired := a.currentAudioDesired()
	desired.TTSEnabled = enabled
	return a.applyAudioSettings(desired)
}

// ---------------------------------------------------------------------------
// RAG Wails Bindings
// ---------------------------------------------------------------------------

// UploadRAGDocument uploads a document into the RAG knowledge base.
func (a *App) UploadRAGDocument(name string, scope string, scopeID string, sourceType string, filePath string) (string, error) {
	if a.ragService == nil {
		return "", fmt.Errorf("RAG service not initialized")
	}
	rawText, err := rag.ExtractTextFromFile(filePath)
	if err != nil {
		return "", err
	}
	doc := &rag.Document{
		Name:       name,
		Scope:      rag.DocumentScope(scope),
		ScopeID:    scopeID,
		SourceType: rag.SourceType(sourceType),
		SourcePath: filePath,
		SizeBytes:  int64(len(rawText)),
		Enabled:    true,
	}
	if err := a.ragService.IngestDocument(doc, rawText); err != nil {
		return "", fmt.Errorf("failed to ingest document: %w", err)
	}
	return doc.ID, nil
}

// ListRAGDocuments returns documents for a given scope.
func (a *App) ListRAGDocuments(scope string, scopeID string) ([]rag.Document, error) {
	if a.ragService == nil {
		return nil, fmt.Errorf("RAG service not initialized")
	}
	return a.ragService.ListDocuments(rag.DocumentScope(scope), scopeID)
}

// DeleteRAGDocument removes a document from the RAG knowledge base.
func (a *App) DeleteRAGDocument(docID string) error {
	if a.ragService == nil {
		return fmt.Errorf("RAG service not initialized")
	}
	return a.ragService.DeleteDocument(docID)
}

// SetRAGEnabled toggles RAG retrieval for the current chat.
func (a *App) SetRAGEnabled(enabled bool) {
	if a.chatService != nil {
		a.chatService.SetRAGEnabled(enabled)
	}
}

// IsRAGEnabled returns whether RAG is currently enabled.
func (a *App) IsRAGEnabled() bool {
	if a.chatService == nil {
		return false
	}
	return a.chatService.IsRAGEnabled()
}
