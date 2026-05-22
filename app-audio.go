package main

import (
	"context"
	"encoding/json"
	"fmt"

	"howl-chat/internal/audio"
	audiotypes "howl-chat/internal/audio/types"
	"howl-chat/internal/backend/rag"
	"howl-chat/internal/errors"
)

// ---------------------------------------------------------------------------
// Audio Wails Bindings
// ---------------------------------------------------------------------------

// TranscribeAudio converts an audio file to text via the ASR pipeline.
func (a *App) TranscribeAudio(audioPath string) (string, error) {
	if a.audioPipeline == nil {
		return "", fmt.Errorf("audio pipeline not initialized")
	}
	result, err := a.audioPipeline.ProcessAudioMessage(context.Background(), audiotypes.AudioMessage{
		Type:      "audio_input",
		AudioPath: audioPath,
	})
	if err != nil {
		a.logError(errors.ASRTranscribeFail, errors.High, "asr.transcribe", "app-audio.go:18", err)
		return "", err
	}
	return result.TranscribedText, nil
}

// SynthesizeSpeech converts text to speech and returns the path to the generated audio file.
func (a *App) SynthesizeSpeech(text string) (string, error) {
	if a.audioPipeline == nil {
		return "", fmt.Errorf("audio pipeline not initialized")
	}
	result, err := a.audioPipeline.ProcessAudioMessage(context.Background(), audiotypes.AudioMessage{
		Type: "audio_output",
		Text: text,
	})
	if err != nil {
		a.logError(errors.TTSSynthesizeFail, errors.High, "tts.synthesize", "app-audio.go:33", err)
		return "", err
	}
	return result.AudioOutputPath, nil
}

// SynthesizeSpeechWithVoice converts text to speech using a per-request voice override.
func (a *App) SynthesizeSpeechWithVoice(text string, voiceJSON string) (string, error) {
	if a.audioPipeline == nil {
		return "", fmt.Errorf("audio pipeline not initialized")
	}
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
	result, err := a.audioPipeline.ProcessAudioMessage(context.Background(), msg)
	if err != nil {
		a.logError(errors.TTSSynthesizeFail, errors.High, "tts.synthesize.voice", "app-audio.go:65", err)
		return "", err
	}
	return result.AudioOutputPath, nil
}

// GetAudioStatus returns the current audio pipeline readiness.
func (a *App) GetAudioStatus() map[string]interface{} {
	if a.audioPipeline == nil {
		return map[string]interface{}{
			"asr_ready":        false,
			"tts_ready":        false,
			"multimodal_ready": false,
			"realtime_ready":   false,
			"active":           false,
		}
	}
	status := a.audioPipeline.GetStatus()
	return map[string]interface{}{
		"asr_ready":        status.ASRReady,
		"tts_ready":        status.TTSReady,
		"multimodal_ready": status.MultimodalReady,
		"realtime_ready":   status.RealtimeReady,
		"active_sessions":  status.ActiveSessions,
		"error_count":      status.ErrorCount,
		"active":           true,
	}
}

// GetTTSAvailableVoices returns the list of available voice names from the current TTS engine.
func (a *App) GetTTSAvailableVoices() []string {
	if a.audioPipeline == nil {
		return nil
	}
	voices := a.audioPipeline.GetTTSVoices()
	result := make([]string, 0, len(voices))
	for _, v := range voices {
		result = append(result, v.ID)
	}
	return result
}

// EnableAudioASR toggles the ASR recognizer on/off.
func (a *App) EnableAudioASR(enabled bool) error {
	if a.audioPipeline == nil {
		return fmt.Errorf("audio pipeline not initialized")
	}
	if enabled {
		_ = a.audioPipeline.Shutdown(context.Background())
		pipe := audio.NewPipeline()
		cfg := audiotypes.PipelineConfig{
			ASREnabled:        true,
			TTSEnabled:        false,
			MultimodalEnabled: false,
			RealtimeEnabled:   false,
			HardwareTier:      "auto",
		}
		if err := pipe.Initialize(context.Background(), cfg); err != nil {
			return err
		}
		a.audioPipeline = pipe
	} else {
		_ = a.audioPipeline.EnableASR(nil)
	}
	return nil
}

// EnableAudioTTS toggles the TTS synthesizer on/off.
func (a *App) EnableAudioTTS(enabled bool) error {
	if a.audioPipeline == nil {
		return fmt.Errorf("audio pipeline not initialized")
	}
	if enabled {
		_ = a.audioPipeline.Shutdown(context.Background())
		pipe := audio.NewPipeline()
		cfg := audiotypes.PipelineConfig{
			ASREnabled:        false,
			TTSEnabled:        true,
			MultimodalEnabled: false,
			RealtimeEnabled:   false,
			HardwareTier:      "auto",
		}
		if err := pipe.Initialize(context.Background(), cfg); err != nil {
			return err
		}
		a.audioPipeline = pipe
	} else {
		_ = a.audioPipeline.EnableTTS(nil)
	}
	return nil
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
