package audio

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/audio/asr/fallback"
	"howl-chat/internal/audio/asr/qwen"
	"howl-chat/internal/audio/asr/whisper"
	"howl-chat/internal/audio/runtime"
	"howl-chat/internal/audio/tts/cache"
	"howl-chat/internal/audio/tts/melotts"
	"howl-chat/internal/audio/tts/piper"
	"howl-chat/internal/audio/types"
)

// Pipeline composes ASR, TTS, and optional multimodal processing.
type Pipeline struct {
	mu      sync.RWMutex
	config  types.PipelineConfig
	status  types.PipelineStatus
	loader  *runtime.ModelLoader
	planner runtime.Constraints

	asrRecognizer types.ASRRecognizer
	ttsSynth      types.TTSSynthesizer
	multimodal    types.MultimodalBridge
	ttsCache      types.AudioCache
}

// NewPipeline creates an empty pipeline with modular runtime loader.
func NewPipeline() *Pipeline {
	loader := runtime.NewModelLoader()
	registerBuiltins(loader)
	return &Pipeline{
		loader: loader,
	}
}

// Initialize resolves and loads configured capability providers.
func (p *Pipeline) Initialize(ctx context.Context, config types.PipelineConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.config = config
	p.status = types.PipelineStatus{}
	p.ttsCache = cache.NewManager(100)

	tier := parseTier(config.HardwareTier)
	p.planner = runtime.Constraints{
		Tier:             tier,
		RequireStreaming: config.RealtimeEnabled,
		RequireGPU:       false,
		PreferredBackend: "local",
	}

	plan := runtime.BuildSelectionPlan(p.planner)
	if config.ASREnabled {
		if err := p.initASR(ctx, plan.ASR); err != nil {
			return err
		}
	}
	if config.TTSEnabled {
		if err := p.initTTS(ctx, plan.TTS); err != nil {
			return err
		}
	}
	if config.MultimodalEnabled {
		if err := p.initMultimodal(ctx, plan.Multimodal); err != nil {
			return err
		}
	}

	return nil
}

// ProcessAudioMessage runs the configured flow for incoming messages.
func (p *Pipeline) ProcessAudioMessage(ctx context.Context, msg types.AudioMessage) (*types.AudioResponse, error) {
	start := time.Now()
	p.mu.RLock()
	asrRec := p.asrRecognizer
	ttsSynth := p.ttsSynth
	cacheStore := p.ttsCache
	p.mu.RUnlock()

	response := &types.AudioResponse{}
	switch msg.Type {
	case "audio_input":
		if asrRec == nil {
			return nil, types.NewAudioError(types.ErrCodeModelNotFound, "asr is not enabled", nil)
		}
		result, err := asrRec.Transcribe(ctx, msg.AudioPath, nil)
		if err != nil {
			p.recordError()
			return nil, err
		}
		response.TranscribedText = strings.TrimSpace(result.Text)
		response.ResponseText = response.TranscribedText

		if ttsSynth != nil && response.ResponseText != "" {
			audioPath, err := p.synthesizeWithCache(ctx, cacheStore, ttsSynth, response.ResponseText)
			if err != nil {
				p.recordError()
				return nil, err
			}
			response.AudioOutputPath = audioPath
		}
	case "audio_output":
		if ttsSynth == nil {
			return nil, types.NewAudioError(types.ErrCodeModelNotFound, "tts is not enabled", nil)
		}
		audioPath, err := p.synthesizeWithCache(ctx, cacheStore, ttsSynth, msg.Text)
		if err != nil {
			p.recordError()
			return nil, err
		}
		response.ResponseText = msg.Text
		response.AudioOutputPath = audioPath
	default:
		return nil, types.NewAudioError(types.ErrCodeFormatNotSupported, "unknown audio message type", nil)
	}

	response.ProcessingTime = time.Since(start).Seconds()
	return response, nil
}

// EnableASR overrides the ASR recognizer.
func (p *Pipeline) EnableASR(recognizer types.ASRRecognizer) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.asrRecognizer = recognizer
	p.status.ASRReady = recognizer != nil
	return nil
}

// EnableTTS overrides the TTS synthesizer.
func (p *Pipeline) EnableTTS(synthesizer types.TTSSynthesizer) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ttsSynth = synthesizer
	p.status.TTSReady = synthesizer != nil
	return nil
}

// EnableMultimodal overrides the multimodal bridge.
func (p *Pipeline) EnableMultimodal(bridge types.MultimodalBridge) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.multimodal = bridge
	p.status.MultimodalReady = bridge != nil
	return nil
}

// GetStatus returns current readiness and error metrics.
func (p *Pipeline) GetStatus() types.PipelineStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.status
}

// Shutdown releases loaded resources through loader lifecycle.
func (p *Pipeline) Shutdown(ctx context.Context) error {
	_ = ctx
	ids := p.loader.GetLoadedModels()
	var errs []string
	for _, id := range ids {
		if err := p.loader.Unload(id); err != nil {
			errs = append(errs, err.Error())
		}
	}

	p.mu.Lock()
	p.asrRecognizer = nil
	p.ttsSynth = nil
	p.multimodal = nil
	p.status.ASRReady = false
	p.status.TTSReady = false
	p.status.MultimodalReady = false
	p.mu.Unlock()

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

func (p *Pipeline) initASR(ctx context.Context, candidates []types.ASRModelType) error {
	if len(candidates) == 0 {
		return types.NewAudioError(types.ErrCodeModelNotFound, "no ASR candidates for current constraints", nil)
	}
	selected := chooseASR(candidates)
	recognizers := make([]types.ASRRecognizer, 0, len(selected))
	recMap := make(map[types.ASRModelType]types.ASRRecognizer)
	for _, modelType := range selected {
		rec, err := p.loader.LoadASR(string(modelType), types.ASRConfig{
			ModelPath: defaultASRModelPath(modelType),
			Language:  "auto",
		})
		if err != nil {
			continue
		}
		recognizers = append(recognizers, rec)
		recMap[modelType] = rec
	}
	if len(recognizers) == 0 {
		return types.NewAudioError(types.ErrCodeModelLoadFailed, "failed to initialize any ASR recognizer", nil)
	}
	chain := fallback.NewChain(fallback.ChainConfig{
		Recognizers:           recognizers,
		RecognizerMap:         recMap,
		MaxAttempts:           1,
		AutoSelectByHardware:  true,
		MinMemoryTier:         parseTier(p.config.HardwareTier),
		Language:              "auto",
	})
	if err := chain.Initialize(ctx, types.ASRConfig{Language: "auto"}); err != nil {
		return err
	}
	p.asrRecognizer = chain
	p.status.ASRReady = true
	return nil
}

func (p *Pipeline) initTTS(ctx context.Context, candidates []types.TTSModelType) error {
	_ = ctx
	if len(candidates) == 0 {
		return types.NewAudioError(types.ErrCodeModelNotFound, "no TTS candidates for current constraints", nil)
	}
	cacheDir := defaultTTSCacheDir()
	var errs []string
	for _, candidate := range candidates {
		cfg := ttsConfigForModel(candidate, cacheDir)
		synth, err := p.loader.LoadTTS(string(candidate), cfg)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", candidate, err))
			continue
		}
		p.ttsSynth = synth
		p.status.TTSReady = true
		return nil
	}
	return types.NewAudioError(types.ErrCodeModelLoadFailed,
		fmt.Sprintf("failed to initialise any configured TTS engine; attempts: %s", strings.Join(errs, "; ")), nil)
}

func (p *Pipeline) initMultimodal(ctx context.Context, candidates []types.MultimodalModelType) error {
	_ = ctx
	if len(candidates) == 0 {
		return nil
	}
	return types.NewAudioError(types.ErrCodeModelNotFound,
		"multimodal models are referenced by the planner but no MultimodalBridge adapters are registered; register one with runtime.ModelLoader before enabling MultimodalEnabled",
		nil)
}

func (p *Pipeline) synthesizeWithCache(ctx context.Context, cacheStore types.AudioCache, synth types.TTSSynthesizer, text string) (string, error) {
	voice := firstVoiceOrDefault(synth, "default")
	req := types.SynthesisRequest{
		Text:  text,
		Voice: voice,
		Speed: 1.0,
	}
	key := types.GenerateCacheKey(req.Text, req.Voice, req.Speed)
	if cached, ok := cacheStore.Get(key); ok {
		return cached.AudioPath, nil
	}
	out, err := synth.Synthesize(ctx, req, nil)
	if err != nil {
		return "", err
	}
	if err := cacheStore.Store(key, out); err != nil {
		return "", err
	}
	return out.AudioPath, nil
}

func (p *Pipeline) recordError() {
	p.mu.Lock()
	p.status.ErrorCount++
	p.mu.Unlock()
}

func registerBuiltins(loader *runtime.ModelLoader) {
	loader.RegisterASR(types.ASRWhisperTiny, func() types.ASRRecognizer { return whisper.New(types.ASRWhisperTiny) })
	loader.RegisterASR(types.ASRWhisperBase, func() types.ASRRecognizer { return whisper.New(types.ASRWhisperBase) })
	loader.RegisterASR(types.ASRWhisperSmall, func() types.ASRRecognizer { return whisper.New(types.ASRWhisperSmall) })
	loader.RegisterASR(types.ASRQwen3ASR06B, func() types.ASRRecognizer { return qwen.NewClient(os.Getenv("HOWL_QWEN_API_KEY")) })
	loader.RegisterASR(types.ASRQwen3ASR17B, func() types.ASRRecognizer { return qwen.NewClient(os.Getenv("HOWL_QWEN_API_KEY")) })
	loader.RegisterTTS(types.TTSPiper, func() types.TTSSynthesizer { return piper.New() })
	loader.RegisterTTS(types.TTSMeloTTS, func() types.TTSSynthesizer { return melotts.New() })
}

func parseTier(raw string) types.MemoryTier {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "tiny":
		return types.TierTiny
	case "small":
		return types.TierSmall
	case "medium":
		return types.TierMedium
	case "large":
		return types.TierLarge
	default:
		return types.TierAuto
	}
}

func chooseASR(candidates []types.ASRModelType) []types.ASRModelType {
	if len(candidates) > 3 {
		return candidates[:3]
	}
	return candidates
}

func defaultPiperModelPath() string {
	if v := os.Getenv("HOWL_TTS_PIPER_MODEL"); v != "" {
		return v
	}
	return "models/piper/en_US-lessac-medium.onnx"
}

func defaultTTSCacheDir() string {
	if v := os.Getenv("HOWL_AUDIO_CACHE_DIR"); v != "" {
		return v
	}
	return "cache/tts"
}

func ttsConfigForModel(model types.TTSModelType, cacheDir string) types.TTSConfig {
	switch model {
	case types.TTSMeloTTS:
		baseURL := strings.TrimSpace(os.Getenv("HOWL_MELOTTS_URL"))
		if baseURL == "" {
			baseURL = "http://127.0.0.1:8888"
		}
		return types.TTSConfig{
			ModelPath:    baseURL,
			CacheDir:     cacheDir,
			CacheEnabled: true,
			MaxCacheSize: 100,
		}
	default:
		return types.TTSConfig{
			ModelPath:    defaultPiperModelPath(),
			CacheDir:     cacheDir,
			CacheEnabled: true,
			MaxCacheSize: 100,
		}
	}
}

func firstVoiceOrDefault(synth types.TTSSynthesizer, fallback string) string {
	voices := synth.GetVoices()
	if len(voices) > 0 && voices[0].ID != "" {
		return voices[0].ID
	}
	return fallback
}

func defaultASRModelPath(model types.ASRModelType) string {
	switch model {
	case types.ASRWhisperTiny:
		return "models/whisper/ggml-tiny.bin"
	case types.ASRWhisperSmall:
		return "models/whisper/ggml-small.bin"
	case types.ASRWhisperBase:
		return "models/whisper/ggml-base.bin"
	default:
		return string(model)
	}
}
