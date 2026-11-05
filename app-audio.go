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
	dir := filepath.Join(os.TempDir(), audioTempDirName)
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
	var override audiotypes.VoiceOverride
	if voiceJSON != "" {
		if err := json.Unmarshal([]byte(voiceJSON), &override); err != nil {
			fmt.Printf("WARN: failed to parse voice override JSON (%v), falling back to default voice\n", err)
		}
	}
	msg := audiotypes.AudioMessage{
		Type: "audio_output",
		Text: text,
	}
	if override.Voice != "" || override.Engine != "" {
		msg.VoiceOverride = &override
	}
	result, err := a.processManagedAudio(context.Background(), msg)
	if err != nil {
		a.logError(errors.TTSSynthesizeFail, errors.High, "tts.synthesize.voice", "app-audio.go:SynthesizeSpeechWithVoice", err)
		return "", err
	}
	return result.AudioOutputPath, nil
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
	readiness, readinessDetail := a.audioReadiness(status)
	active := status.Phase != "closing" && status.Phase != "closed" && status.Phase != "closed_with_errors"

	// apply_in_progress is the manager's live build gauge (Status.Builds). It is
	// an active-operation count, not a cumulative counter.
	applyInProgress := status.Builds > 0

	// applied publishes the ACTIVE generation's effective configuration. It is
	// never the desired/newest request, so the panel can compare what it asked
	// for with what is actually installed.
	var applied interface{}
	if a.audioManager != nil {
		if installed, ok := a.audioManager.AppliedSettings(); ok {
			applied = map[string]interface{}{
				"asr_enabled": installed.ASREnabled,
				"tts_enabled": installed.TTSEnabled,
				"asr_engine":  installed.ASREngine,
				"tts_engine":  installed.TTSEngine,
			}
		}
	}
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
		"apply_in_progress": applyInProgress,
		"applied":           applied,
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
func (a *App) audioReadiness(status audio.Status) (string, string) {
	if !status.ASR.Requested {
		return "disabled", "ASR is disabled"
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

// ApplyAudioSettings applies the desired capabilities and explicit engine
// selection as one generation. It is the bridge the settings UI calls when the
// user saves and when persisted settings load at startup; an explicit engine
// that cannot run is refused with its reason and never replaced by another
// provider. Storage remains T2.9.6's job; this only applies what was stored.
func (a *App) ApplyAudioSettings(asrEnabled bool, ttsEnabled bool, asrEngine string, ttsEngine string) error {
	return a.applyAudioSettings(audio.Settings{
		ASREnabled: asrEnabled,
		TTSEnabled: ttsEnabled,
		ASREngine:  asrEngine,
		TTSEngine:  ttsEngine,
	})
}

// applyAudioSettings is the minimal temporary 5A Apply bridge. It deliberately
// does not persist settings or expose D5 status fields.
func (a *App) applyAudioSettings(settings audio.Settings) error {
	if a.audioManager == nil {
		return fmt.Errorf("audio manager not initialized")
	}
	_, _, err := a.audioManager.Apply(context.Background(), settings)
	if err != nil {
		a.logError(errors.ASRModelLoadFail, errors.High, "audio.apply", "app-audio.go:applyAudioSettings",
			fmt.Errorf("apply ASR engine %q: %w", settings.ASREngine, err))
	}
	return err
}

// EnableAudioASR toggles the ASR recognizer without discarding the desired TTS
// engine. The previous implementation rebuilt the settings from one boolean,
// which silently reset an explicit engine chosen in the settings panel.
func (a *App) EnableAudioASR(enabled bool) error {
	if a.audioManager == nil {
		return fmt.Errorf("audio manager not initialized")
	}
	desired := a.audioManager.Desired()
	desired.ASREnabled = enabled
	return a.applyAudioSettings(desired)
}

// EnableAudioTTS toggles the TTS synthesizer without discarding the desired
// ASR engine.
func (a *App) EnableAudioTTS(enabled bool) error {
	if a.audioManager == nil {
		return fmt.Errorf("audio manager not initialized")
	}
	desired := a.audioManager.Desired()
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
