package audio

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"howl-chat/internal/audio/tts/cache"
	"howl-chat/internal/audio/types"
)

func newCharacterVoicePipeline(t *testing.T) (*Pipeline, *providerSet) {
	t.Helper()
	p, set := newScriptedPipeline(providerBehavior{})

	piper := set.tts[types.TTSPiper]
	kokoro := set.tts[types.TTSKokoro]
	melo := set.tts[types.TTSMeloTTS]
	piper.voices = []types.Voice{{ID: "piper-default"}}
	kokoro.voices = []types.Voice{{ID: "af_heart"}, {ID: "am_echo"}}
	melo.voices = []types.Voice{{ID: "EN"}, {ID: "ZH"}}

	p.mu.Lock()
	p.config = types.PipelineConfig{
		TTSEnabled: true,
		TTSEngine:  string(types.TTSPiper),
		TTSSpeed:   1.0,
	}
	p.ttsSynth = piper
	p.activeTTSEngine = string(types.TTSPiper)
	p.ttsProviders = map[string]*ttsProvider{
		string(types.TTSPiper): {key: string(types.TTSPiper), synth: piper, engine: string(types.TTSPiper)},
	}
	p.ttsCache = cache.NewManager(10)
	p.mu.Unlock()
	return p, set
}

func voiceRequest(engine, voice, characterID, text string) types.AudioMessage {
	return types.AudioMessage{
		Type: "audio_output",
		Text: text,
		VoiceOverride: &types.VoiceOverride{
			CharacterID: characterID,
			Engine:      engine,
			Voice:       voice,
		},
	}
}

func TestCharacterVoiceRequestsResolveToTheRequestedEngineAndVoice(t *testing.T) {
	p, set := newCharacterVoicePipeline(t)
	kokoro := set.tts[types.TTSKokoro]
	melo := set.tts[types.TTSMeloTTS]
	piper := set.tts[types.TTSPiper]

	sequence := []struct {
		name  string
		msg   types.AudioMessage
		synth *scriptedTTS
		want  string
	}{
		{
			name:  "A kokoro",
			msg:   voiceRequest("kokoro", "af_heart", "char_a", "A kokoro"),
			synth: kokoro,
			want:  "af_heart",
		},
		{
			name:  "B melotts",
			msg:   voiceRequest("melotts", "ZH", "char_b", "B melotts"),
			synth: melo,
			want:  "ZH",
		},
		{
			name:  "A kokoro again",
			msg:   voiceRequest("kokoro", "am_echo", "char_a", "A kokoro again"),
			synth: kokoro,
			want:  "am_echo",
		},
	}
	for _, step := range sequence {
		if _, err := p.ProcessAudioMessage(context.Background(), step.msg); err != nil {
			t.Fatalf("%s: ProcessAudioMessage: %v", step.name, err)
		}
		req, ok := step.synth.lastRequest()
		if !ok {
			t.Fatalf("%s: synthesizer received no request", step.name)
		}
		if req.Engine != step.msg.VoiceOverride.Engine || req.Voice != step.want {
			t.Fatalf("%s: request engine/voice = %q/%q, want %q/%q",
				step.name, req.Engine, req.Voice, step.msg.VoiceOverride.Engine, step.want)
		}
		if step.msg.VoiceOverride.Engine == string(types.TTSMeloTTS) && req.Language != step.want {
			t.Fatalf("%s: melo language = %q, want %q", step.name, req.Language, step.want)
		}
		if req.CharacterID != step.msg.VoiceOverride.CharacterID {
			t.Fatalf("%s: character id = %q, want %q", step.name, req.CharacterID, step.msg.VoiceOverride.CharacterID)
		}
	}

	if piper.synthCalls.Load() != 0 {
		t.Fatalf("active piper synth calls = %d, want 0", piper.synthCalls.Load())
	}
	p.mu.RLock()
	activeSynth := p.ttsSynth
	activeEngine := p.activeTTSEngine
	p.mu.RUnlock()
	if activeSynth != piper || activeEngine != string(types.TTSPiper) {
		t.Fatal("explicit character engines mutated the active global TTS binding")
	}
}

func TestCharacterVoiceInheritanceUsesTheGlobalDefault(t *testing.T) {
	p, set := newCharacterVoicePipeline(t)
	piper := set.tts[types.TTSPiper]

	if _, err := p.ProcessAudioMessage(context.Background(), types.AudioMessage{
		Type: "audio_output",
		Text: "no override",
	}); err != nil {
		t.Fatalf("ProcessAudioMessage without override: %v", err)
	}
	req, ok := piper.lastRequest()
	if !ok {
		t.Fatal("active piper received no request")
	}
	if req.Engine != string(types.TTSPiper) || req.Voice != "piper-default" || req.CharacterID != "" {
		t.Fatalf("inherited request = engine %q voice %q character %q, want piper/piper-default/empty",
			req.Engine, req.Voice, req.CharacterID)
	}

	globalVoicePath := filepath.Join(t.TempDir(), "piper-default.onnx")
	if err := os.WriteFile(globalVoicePath, []byte("model"), 0o600); err != nil {
		t.Fatalf("write global piper model: %v", err)
	}
	t.Setenv("HOWL_TTS_PIPER_MODEL", globalVoicePath)
	p.mu.Lock()
	p.config.TTSVoice = "piper-default"
	p.mu.Unlock()
	if _, err := p.ProcessAudioMessage(context.Background(), types.AudioMessage{
		Type: "audio_output",
		Text: "configured global voice",
	}); err != nil {
		t.Fatalf("ProcessAudioMessage with configured global voice: %v", err)
	}
	req, ok = piper.lastRequest()
	if !ok || req.Voice != "piper-default" {
		t.Fatalf("configured global voice request = %+v, want piper-default", req)
	}

	if _, err := p.ProcessAudioMessage(context.Background(), voiceRequest("", "", "char_inherit", "explicit inheritance")); err != nil {
		t.Fatalf("ProcessAudioMessage with explicit inheritance: %v", err)
	}
	req, ok = piper.lastRequest()
	if !ok {
		t.Fatal("active piper received no explicit-inheritance request")
	}
	if req.Engine != string(types.TTSPiper) || req.Voice != "piper-default" || req.CharacterID != "char_inherit" {
		t.Fatalf("explicit inheritance = engine %q voice %q character %q, want piper/piper-default/char_inherit",
			req.Engine, req.Voice, req.CharacterID)
	}
}

func TestCharacterVoiceUnsupportedConfigurationFailsWithoutFallback(t *testing.T) {
	p, set := newCharacterVoicePipeline(t)
	piper := set.tts[types.TTSPiper]
	kokoro := set.tts[types.TTSKokoro]

	before := piper.synthCalls.Load()
	_, err := p.ProcessAudioMessage(context.Background(), voiceRequest("browser", "", "char_bad", "browser"))
	if err == nil || !strings.Contains(err.Error(), "renderer-only") {
		t.Fatalf("browser engine error = %v, want renderer-only failure", err)
	}
	_, err = p.ProcessAudioMessage(context.Background(), voiceRequest("not-an-engine", "", "char_bad", "unknown"))
	if err == nil || !strings.Contains(err.Error(), "unsupported TTS engine") {
		t.Fatalf("unknown engine error = %v, want unsupported engine failure", err)
	}
	_, err = p.ProcessAudioMessage(context.Background(), voiceRequest("kokoro", "not-a-voice", "char_bad", "unknown voice"))
	if err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("unknown voice error = %v, want unsupported voice failure", err)
	}
	if got := piper.synthCalls.Load(); got != before {
		t.Fatalf("active piper synth calls = %d, want %d after failures", got, before)
	}
	if got := kokoro.synthCalls.Load(); got != 0 {
		t.Fatalf("kokoro synth calls = %d, want 0 for an invalid voice", got)
	}

	if _, err := p.ProcessAudioMessage(context.Background(), voiceRequest("", "", "char_next", "next")); err != nil {
		t.Fatalf("request after failures: %v", err)
	}
	req, ok := piper.lastRequest()
	if !ok || req.Voice != "piper-default" || req.CharacterID != "char_next" {
		t.Fatalf("request after failures = %+v, want active piper-default/char_next", req)
	}
}

func TestPiperCharacterVoiceUsesRequestScopedModelPath(t *testing.T) {
	voice := "character-test-voice"
	path := filepath.Join(t.TempDir(), voice+".onnx")
	if err := os.WriteFile(path, []byte("model"), 0o600); err != nil {
		t.Fatalf("write piper model: %v", err)
	}
	t.Setenv("HOWL_TTS_PIPER_MODEL", path)

	p, set := newCharacterVoicePipeline(t)
	if _, err := p.ProcessAudioMessage(context.Background(), voiceRequest("piper", voice, "char_piper", "piper voice")); err != nil {
		t.Fatalf("ProcessAudioMessage: %v", err)
	}
	req, ok := set.tts[types.TTSPiper].lastRequest()
	if !ok {
		t.Fatal("piper received no request")
	}
	if req.Voice != voice || req.VoicePath != path || req.Engine != string(types.TTSPiper) {
		t.Fatalf("piper request = voice %q path %q engine %q, want %q/%q/piper",
			req.Voice, req.VoicePath, req.Engine, voice, path)
	}
}

func TestCharacterVoiceInheritedAutoEngineKeepsEstablishedSpeed(t *testing.T) {
	p, set := newCharacterVoicePipeline(t)
	p.mu.Lock()
	p.config.TTSEngine = "auto"
	p.config.TTSSpeed = 1.25
	p.mu.Unlock()

	if _, err := p.ProcessAudioMessage(context.Background(), types.AudioMessage{
		Type: "audio_output",
		Text: "auto inherited speed",
	}); err != nil {
		t.Fatalf("ProcessAudioMessage: %v", err)
	}
	req, ok := set.tts[types.TTSPiper].lastRequest()
	if !ok {
		t.Fatal("active piper received no request")
	}
	if req.Speed != 1.0 {
		t.Fatalf("auto inherited speed = %g, want established 1.0", req.Speed)
	}
}

func TestCharacterVoiceSpeedPrecedenceIsPreserved(t *testing.T) {
	p, set := newCharacterVoicePipeline(t)
	p.mu.Lock()
	p.config.TTSSpeed = 1.25
	p.mu.Unlock()
	piper := set.tts[types.TTSPiper]

	msg := voiceRequest("piper", "", "char_speed", "speed default")
	if _, err := p.ProcessAudioMessage(context.Background(), msg); err != nil {
		t.Fatalf("ProcessAudioMessage default speed: %v", err)
	}
	req, _ := piper.lastRequest()
	if req.Speed != 1.25 {
		t.Fatalf("inherited speed = %g, want 1.25", req.Speed)
	}

	msg.VoiceOverride.Speed = 1.5
	msg.Text = "speed override"
	if _, err := p.ProcessAudioMessage(context.Background(), msg); err != nil {
		t.Fatalf("ProcessAudioMessage override speed: %v", err)
	}
	req, _ = piper.lastRequest()
	if req.Speed != 1.5 {
		t.Fatalf("override speed = %g, want 1.5", req.Speed)
	}
}

func TestConcurrentCharacterVoiceRequestsRemainIsolated(t *testing.T) {
	p, set := newCharacterVoicePipeline(t)
	kokoro := set.tts[types.TTSKokoro]
	melo := set.tts[types.TTSMeloTTS]

	type request struct {
		msg   types.AudioMessage
		synth *scriptedTTS
		want  string
	}
	requests := []request{
		{msg: voiceRequest("kokoro", "af_heart", "char_a", "concurrent A"), synth: kokoro, want: "af_heart"},
		{msg: voiceRequest("melotts", "ZH", "char_b", "concurrent B"), synth: melo, want: "ZH"},
		{msg: voiceRequest("kokoro", "am_echo", "char_c", "concurrent C"), synth: kokoro, want: "am_echo"},
	}

	var wg sync.WaitGroup
	errs := make(chan error, len(requests))
	for _, step := range requests {
		step := step
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.ProcessAudioMessage(context.Background(), step.msg)
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent ProcessAudioMessage: %v", err)
		}
	}

	for _, step := range requests {
		step.synth.mu.Lock()
		found := false
		for _, req := range step.synth.requests {
			if req.CharacterID == step.msg.VoiceOverride.CharacterID && req.Voice == step.want && req.Engine == step.msg.VoiceOverride.Engine {
				found = true
				break
			}
		}
		step.synth.mu.Unlock()
		if !found {
			t.Fatalf("request %q did not reach its engine with voice %q", step.msg.Text, step.want)
		}
	}

	p.mu.RLock()
	active := p.ttsSynth
	engine := p.activeTTSEngine
	p.mu.RUnlock()
	if active != set.tts[types.TTSPiper] || engine != string(types.TTSPiper) {
		t.Fatal("concurrent explicit requests changed the active global TTS binding")
	}
}
