package audio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"howl-chat/internal/audio/types"
)

func newCachePipeline(t *testing.T) (*Pipeline, *providerSet) {
	t.Helper()
	p, set := newCharacterVoicePipeline(t)
	set.tts[types.TTSPiper].voices = []types.Voice{{ID: "shared"}}
	set.tts[types.TTSKokoro].voices = []types.Voice{{ID: "shared"}}
	set.tts[types.TTSMeloTTS].voices = []types.Voice{{ID: "shared"}}

	modelPath := filepath.Join(t.TempDir(), "shared.onnx")
	if err := os.WriteFile(modelPath, []byte("model"), 0o600); err != nil {
		t.Fatalf("write piper model: %v", err)
	}
	t.Setenv("HOWL_TTS_PIPER_MODEL", modelPath)

	p.mu.Lock()
	p.config.TTSVoice = ""
	p.mu.Unlock()
	return p, set
}

func cacheStatsFor(p *Pipeline) types.CacheStats {
	p.mu.RLock()
	store := p.ttsCache
	p.mu.RUnlock()
	return store.GetStats()
}

func TestTTSCacheSeparatesEnginesWithSameTextAndVoice(t *testing.T) {
	p, set := newCachePipeline(t)
	piper := set.tts[types.TTSPiper]
	kokoro := set.tts[types.TTSKokoro]

	piperMsg := voiceRequest("piper", "shared", "character-a", "same text")
	kokoroMsg := voiceRequest("kokoro", "shared", "character-b", "same text")

	if _, err := p.ProcessAudioMessage(context.Background(), piperMsg); err != nil {
		t.Fatalf("piper request: %v", err)
	}
	if _, err := p.ProcessAudioMessage(context.Background(), kokoroMsg); err != nil {
		t.Fatalf("kokoro request: %v", err)
	}
	if got := piper.synthCalls.Load(); got != 1 {
		t.Fatalf("piper synthesis calls = %d, want 1", got)
	}
	if got := kokoro.synthCalls.Load(); got != 1 {
		t.Fatalf("kokoro synthesis calls = %d, want 1", got)
	}
	if got := cacheStatsFor(p).HitCount; got != 0 {
		t.Fatalf("cache hits after first engine requests = %d, want 0", got)
	}

	if _, err := p.ProcessAudioMessage(context.Background(), piperMsg); err != nil {
		t.Fatalf("repeat piper request: %v", err)
	}
	if _, err := p.ProcessAudioMessage(context.Background(), kokoroMsg); err != nil {
		t.Fatalf("repeat kokoro request: %v", err)
	}
	if got := piper.synthCalls.Load(); got != 1 {
		t.Fatalf("piper synthesis calls after repeat = %d, want 1", got)
	}
	if got := kokoro.synthCalls.Load(); got != 1 {
		t.Fatalf("kokoro synthesis calls after repeat = %d, want 1", got)
	}
	stats := cacheStatsFor(p)
	if stats.HitCount != 2 {
		t.Fatalf("cache hits = %d, want 2", stats.HitCount)
	}
	if stats.MissCount != 2 {
		t.Fatalf("cache misses = %d, want 2", stats.MissCount)
	}
}

func TestTTSCachePreservesCharacterABASequence(t *testing.T) {
	p, set := newCachePipeline(t)
	piper := set.tts[types.TTSPiper]
	kokoro := set.tts[types.TTSKokoro]

	a := voiceRequest("piper", "shared", "character-a", "A")
	b := voiceRequest("kokoro", "shared", "character-b", "B")
	for _, msg := range []types.AudioMessage{a, b, a} {
		if _, err := p.ProcessAudioMessage(context.Background(), msg); err != nil {
			t.Fatalf("request %q: %v", msg.Text, err)
		}
	}
	if got := piper.synthCalls.Load(); got != 1 {
		t.Fatalf("piper synthesis calls = %d, want 1 for A-B-A", got)
	}
	if got := kokoro.synthCalls.Load(); got != 1 {
		t.Fatalf("kokoro synthesis calls = %d, want 1 for A-B-A", got)
	}
	stats := cacheStatsFor(p)
	if stats.HitCount != 1 {
		t.Fatalf("A-B-A cache hits = %d, want 1", stats.HitCount)
	}
	if stats.MissCount != 2 {
		t.Fatalf("A-B-A cache misses = %d, want 2", stats.MissCount)
	}
}

func TestTTSCacheVersionRejectsLegacyEntry(t *testing.T) {
	p, set := newCachePipeline(t)
	piper := set.tts[types.TTSPiper]

	p.mu.RLock()
	store := p.ttsCache
	p.mu.RUnlock()
	legacyKey := types.GenerateCacheKey("legacy", "shared", 1)
	if err := store.Store(legacyKey, &types.SynthesisResult{AudioPath: "legacy.wav"}); err != nil {
		t.Fatalf("store legacy entry: %v", err)
	}

	resp, err := p.ProcessAudioMessage(context.Background(), voiceRequest("", "", "character", "legacy"))
	if err != nil {
		t.Fatalf("ProcessAudioMessage: %v", err)
	}
	if resp.AudioOutputPath != "scripted.wav" {
		t.Fatalf("audio path = %q, want fresh synthesis result", resp.AudioOutputPath)
	}
	if got := piper.synthCalls.Load(); got != 1 {
		t.Fatalf("piper synthesis calls = %d, want 1 after legacy key miss", got)
	}
}

func TestInheritedEngineSpeedUsesResolvedEngineValidation(t *testing.T) {
	p, set := newCachePipeline(t)
	melo := set.tts[types.TTSMeloTTS]

	p.mu.Lock()
	p.config.TTSEngine = string(types.TTSMeloTTS)
	p.ttsSynth = melo
	p.activeTTSEngine = string(types.TTSMeloTTS)
	p.mu.Unlock()

	invalid := voiceRequest("", "", "character-inherit", "inherited melo invalid")
	invalid.VoiceOverride.Speed = 1.5
	if _, err := p.ProcessAudioMessage(context.Background(), invalid); err == nil {
		t.Fatal("inherited MeloTTS speed above its product guard was accepted")
	} else if !strings.Contains(err.Error(), "between 0.8 and 1.2") {
		t.Fatalf("inherited MeloTTS override error = %v, want its product guard", err)
	}
	if got := melo.synthCalls.Load(); got != 0 {
		t.Fatalf("inherited MeloTTS synthesis calls = %d, want 0 for invalid speed", got)
	}

	valid := voiceRequest("", "", "character-inherit", "inherited melo valid")
	valid.VoiceOverride.Speed = 1.1
	if _, err := p.ProcessAudioMessage(context.Background(), valid); err != nil {
		t.Fatalf("valid inherited MeloTTS override: %v", err)
	}
	req, ok := melo.lastRequest()
	if !ok || req.Speed != 1.1 {
		t.Fatalf("inherited MeloTTS request = %+v, want speed 1.1", req)
	}

	kokoro := set.tts[types.TTSKokoro]
	p.mu.Lock()
	p.config.TTSEngine = string(types.TTSKokoro)
	p.ttsSynth = kokoro
	p.activeTTSEngine = string(types.TTSKokoro)
	p.mu.Unlock()

	invalidKokoro := voiceRequest("", "", "character-inherit", "inherited kokoro invalid")
	invalidKokoro.VoiceOverride.Speed = 0.4
	if _, err := p.ProcessAudioMessage(context.Background(), invalidKokoro); err == nil {
		t.Fatal("inherited Kokoro speed below its range was accepted")
	} else if !strings.Contains(err.Error(), "between 0.5 and 2") {
		t.Fatalf("inherited Kokoro override error = %v, want its engine range", err)
	}
	if got := kokoro.synthCalls.Load(); got != 0 {
		t.Fatalf("inherited Kokoro synthesis calls = %d, want 0 for invalid speed", got)
	}
}

func TestMeloTTSModelIdentityIncludesConfiguredOpenAIModel(t *testing.T) {
	t.Setenv("HOWL_MELOTTS_OPENAI_MODEL", "tts-a")
	cfgA := ttsConfigForModel(types.TTSMeloTTS, t.TempDir())
	t.Setenv("HOWL_MELOTTS_OPENAI_MODEL", "tts-b")
	cfgB := ttsConfigForModel(types.TTSMeloTTS, t.TempDir())

	if cfgA.ModelName != "tts-a" || cfgB.ModelName != "tts-b" {
		t.Fatalf("MeloTTS configured models = %q/%q, want tts-a/tts-b", cfgA.ModelName, cfgB.ModelName)
	}
	identityA := types.TTSModelIdentity(string(types.TTSMeloTTS), cfgA, "")
	identityB := types.TTSModelIdentity(string(types.TTSMeloTTS), cfgB, "")
	if identityA == identityB {
		t.Fatal("MeloTTS OpenAI model setting did not participate in the cache identity")
	}
}

func TestVoiceOverrideSpeedUsesEngineSpecificValidation(t *testing.T) {
	p, set := newCachePipeline(t)
	melo := set.tts[types.TTSMeloTTS]
	kokoro := set.tts[types.TTSKokoro]

	invalidMelo := voiceRequest("melotts", "shared", "character-melo", "melo invalid")
	invalidMelo.VoiceOverride.Speed = 1.5
	if _, err := p.ProcessAudioMessage(context.Background(), invalidMelo); err == nil {
		t.Fatal("MeloTTS override speed above its product guard was accepted")
	} else if !strings.Contains(err.Error(), "between 0.8 and 1.2") {
		t.Fatalf("MeloTTS override error = %v, want its product guard", err)
	}
	if got := melo.synthCalls.Load(); got != 0 {
		t.Fatalf("MeloTTS synthesis calls = %d, want 0 for invalid override", got)
	}

	validMelo := voiceRequest("melotts", "shared", "character-melo", "melo valid")
	validMelo.VoiceOverride.Speed = 1.1
	if _, err := p.ProcessAudioMessage(context.Background(), validMelo); err != nil {
		t.Fatalf("valid MeloTTTS override: %v", err)
	}
	req, ok := melo.lastRequest()
	if !ok || req.Speed != 1.1 {
		t.Fatalf("valid MeloTTTS request = %+v, want speed 1.1", req)
	}

	invalidKokoro := voiceRequest("kokoro", "shared", "character-kokoro", "kokoro invalid")
	invalidKokoro.VoiceOverride.Speed = 0.4
	if _, err := p.ProcessAudioMessage(context.Background(), invalidKokoro); err == nil {
		t.Fatal("Kokoro override speed below its range was accepted")
	} else if !strings.Contains(err.Error(), "between 0.5 and 2") {
		t.Fatalf("Kokoro override error = %v, want its engine range", err)
	}
	if got := kokoro.synthCalls.Load(); got != 0 {
		t.Fatalf("Kokoro synthesis calls = %d, want 0 for invalid override", got)
	}

	validKokoro := voiceRequest("kokoro", "shared", "character-kokoro", "kokoro valid")
	validKokoro.VoiceOverride.Speed = 1.25
	if _, err := p.ProcessAudioMessage(context.Background(), validKokoro); err != nil {
		t.Fatalf("valid Kokoro override: %v", err)
	}
	req, ok = kokoro.lastRequest()
	if !ok || req.Speed != 1.25 {
		t.Fatalf("valid Kokoro request = %+v, want speed 1.25", req)
	}

	invalidPiper := voiceRequest("piper", "shared", "character-piper", "piper invalid")
	invalidPiper.VoiceOverride.Speed = 2.5
	if _, err := p.ProcessAudioMessage(context.Background(), invalidPiper); err == nil {
		t.Fatal("Piper override speed above its range was accepted")
	} else if !strings.Contains(err.Error(), "between 0.5 and 2") {
		t.Fatalf("Piper override error = %v, want its engine range", err)
	}
}

func TestTTSCacheColdWarmAndRepeatedReuse(t *testing.T) {
	p, set := newCachePipeline(t)
	piper := set.tts[types.TTSPiper]
	piper.delay = 20 * time.Millisecond

	msg := voiceRequest("", "", "character", "performance reuse")
	start := time.Now()
	if _, err := p.ProcessAudioMessage(context.Background(), msg); err != nil {
		t.Fatalf("cold request: %v", err)
	}
	cold := time.Since(start)

	const warmBatchRequests = 10_000
	start = time.Now()
	for i := 0; i < warmBatchRequests; i++ {
		if _, err := p.ProcessAudioMessage(context.Background(), msg); err != nil {
			t.Fatalf("warm batch request %d: %v", i, err)
		}
	}
	warmTotal := time.Since(start)
	warmAverage := warmTotal / time.Duration(warmBatchRequests)

	if got := piper.synthCalls.Load(); got != 1 {
		t.Fatalf("piper synthesis calls = %d, want 1 after cold and warm batch", got)
	}
	if warmAverage >= cold {
		t.Fatalf("warm average latency %s did not beat cold synthesis latency %s", warmAverage, cold)
	}
	stats := cacheStatsFor(p)
	t.Logf("latency cold=%s warm_batch_requests=%d warm_batch_total=%s warm_batch_average=%s",
		cold, warmBatchRequests, warmTotal, warmAverage)
	t.Logf("cache reuse synth_calls=%d hits=%d misses=%d entries=%d size_bytes=%d",
		piper.synthCalls.Load(), stats.HitCount, stats.MissCount, stats.EntryCount, stats.SizeBytes)
}

func TestTTSCacheAlternatingCharactersReuse(t *testing.T) {
	p, set := newCachePipeline(t)
	piper := set.tts[types.TTSPiper]
	kokoro := set.tts[types.TTSKokoro]
	piper.delay = 5 * time.Millisecond
	kokoro.delay = 5 * time.Millisecond

	a := voiceRequest("piper", "shared", "character-a", "alternating")
	b := voiceRequest("kokoro", "shared", "character-b", "alternating")
	for i := 0; i < 10; i++ {
		if _, err := p.ProcessAudioMessage(context.Background(), a); err != nil {
			t.Fatalf("A request %d: %v", i, err)
		}
		if _, err := p.ProcessAudioMessage(context.Background(), b); err != nil {
			t.Fatalf("B request %d: %v", i, err)
		}
	}
	if got := piper.synthCalls.Load(); got != 1 {
		t.Fatalf("piper synthesis calls = %d, want 1", got)
	}
	if got := kokoro.synthCalls.Load(); got != 1 {
		t.Fatalf("kokoro synthesis calls = %d, want 1", got)
	}
	stats := cacheStatsFor(p)
	if stats.HitCount != 18 || stats.MissCount != 2 {
		t.Fatalf("alternating stats = hits %d misses %d, want 18/2", stats.HitCount, stats.MissCount)
	}
	t.Logf("alternating requests=20 synth_calls=%d hits=%d misses=%d",
		piper.synthCalls.Load()+kokoro.synthCalls.Load(), stats.HitCount, stats.MissCount)
}

func TestConcurrentIdenticalRequestsSynthesizeOnce(t *testing.T) {
	p, set := newCachePipeline(t)
	piper := set.tts[types.TTSPiper]
	piper.delay = 30 * time.Millisecond

	const requests = 8
	start := make(chan struct{})
	errs := make(chan error, requests)
	var wg sync.WaitGroup
	for i := 0; i < requests; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := p.ProcessAudioMessage(context.Background(), voiceRequest("", "", fmt.Sprintf("character-%d", i), "concurrent identical"))
			errs <- err
		}()
	}
	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent identical request: %v", err)
		}
	}
	if got := piper.synthCalls.Load(); got != 1 {
		t.Fatalf("piper synthesis calls = %d, want 1 for concurrent identical requests", got)
	}
	stats := cacheStatsFor(p)
	t.Logf("concurrent identical requests=%d synth_calls=%d hits=%d misses=%d",
		requests, piper.synthCalls.Load(), stats.HitCount, stats.MissCount)
}

func TestConcurrentDifferentConfigurationsRemainIsolated(t *testing.T) {
	p, set := newCachePipeline(t)
	kokoro := set.tts[types.TTSKokoro]
	melo := set.tts[types.TTSMeloTTS]
	kokoro.delay = 10 * time.Millisecond
	melo.delay = 10 * time.Millisecond

	const rounds = 6
	start := make(chan struct{})
	errs := make(chan error, rounds*2)
	var wg sync.WaitGroup
	launch := func(msg types.AudioMessage) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := p.ProcessAudioMessage(context.Background(), msg)
			errs <- err
		}()
	}
	for i := 0; i < rounds; i++ {
		launch(voiceRequest("kokoro", "shared", fmt.Sprintf("kokoro-%d", i), "concurrent different"))
		launch(voiceRequest("melotts", "shared", fmt.Sprintf("melo-%d", i), "concurrent different"))
	}
	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent different request: %v", err)
		}
	}
	if got := kokoro.synthCalls.Load(); got != 1 {
		t.Fatalf("kokoro synthesis calls = %d, want 1", got)
	}
	if got := melo.synthCalls.Load(); got != 1 {
		t.Fatalf("melo synthesis calls = %d, want 1", got)
	}
	stats := cacheStatsFor(p)
	if stats.HitCount != 0 || stats.MissCount != 2 {
		t.Fatalf("concurrent different stats = hits %d misses %d, want 0/2", stats.HitCount, stats.MissCount)
	}
}
