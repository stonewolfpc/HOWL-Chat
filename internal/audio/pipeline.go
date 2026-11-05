package audio

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/audio/asr/fallback"
	"howl-chat/internal/audio/asr/llama"
	"howl-chat/internal/audio/asr/qwen"
	"howl-chat/internal/audio/asr/whisper"
	"howl-chat/internal/audio/capability"
	"howl-chat/internal/audio/runtime"
	"howl-chat/internal/audio/tts/cache"
	"howl-chat/internal/audio/tts/kokoro"
	"howl-chat/internal/audio/tts/melotts"
	"howl-chat/internal/audio/tts/piper"
	"howl-chat/internal/audio/types"
	"howl-chat/internal/paths"
)

// Pipeline composes ASR, TTS, and optional multimodal processing.
type Pipeline struct {
	lifecycleMu   sync.RWMutex
	replacementMu sync.Mutex
	mu            sync.RWMutex
	config        types.PipelineConfig
	status        types.PipelineStatus
	loader        *runtime.ModelLoader
	planner       runtime.Constraints

	asrRecognizer  types.ASRRecognizer
	ttsSynth       types.TTSSynthesizer
	multimodal     types.MultimodalBridge
	ttsCache       types.AudioCache
	loadedVoice    string // Currently loaded TTS voice model name
	initialization initOutcome
}

type capabilityInit struct {
	requested bool
	ready     bool
	err       error
}

type initOutcome struct {
	asr capabilityInit
	tts capabilityInit
}

func (o initOutcome) aggregateError() error {
	var errs []error
	if o.asr.requested && !o.asr.ready {
		err := o.asr.err
		if err == nil {
			err = ErrCapabilityNotReady
		}
		errs = append(errs, fmt.Errorf("ASR: %w", err))
	}
	if o.tts.requested && !o.tts.ready {
		err := o.tts.err
		if err == nil {
			err = ErrCapabilityNotReady
		}
		errs = append(errs, fmt.Errorf("TTS: %w", err))
	}
	return errors.Join(errs...)
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
	if ctx == nil {
		ctx = context.Background()
	}
	p.lifecycleMu.Lock()
	defer p.lifecycleMu.Unlock()

	constraints := runtime.Constraints{
		Tier:             parseTier(config.HardwareTier),
		RequireStreaming: config.RealtimeEnabled,
		RequireGPU:       false,
		PreferredBackend: "local",
	}

	p.mu.Lock()
	p.config = config
	p.status = types.PipelineStatus{}
	p.ttsCache = cache.NewManager(100)
	p.planner = constraints
	p.initialization = initOutcome{
		asr: capabilityInit{requested: config.ASREnabled},
		tts: capabilityInit{requested: config.TTSEnabled},
	}
	p.mu.Unlock()

	plan := runtime.BuildSelectionPlan(constraints)
	if config.ASREnabled {
		var recognizer types.ASRRecognizer
		outcome := capabilityInit{requested: true}
		candidates, selectionErr := asrCandidatesFor(plan.ASR, config.ASREngine)
		if selectionErr != nil {
			outcome.err = selectionErr
		} else {
			recognizer, outcome = p.initASR(ctx, candidates, config)
		}
		p.mu.Lock()
		p.initialization.asr = outcome
		p.asrRecognizer = nil
		p.status.ASRReady = false
		if outcome.ready {
			p.asrRecognizer = recognizer
			p.status.ASRReady = true
		}
		p.mu.Unlock()
	}
	if config.TTSEnabled {
		var synth types.TTSSynthesizer
		var loadedVoice string
		outcome := capabilityInit{requested: true}
		candidates, selectionErr := ttsCandidatesFor(plan.TTS, config.TTSEngine)
		if selectionErr != nil {
			outcome.err = selectionErr
		} else {
			synth, loadedVoice, outcome = p.initTTS(ctx, candidates)
		}
		p.mu.Lock()
		p.initialization.tts = outcome
		p.ttsSynth = nil
		p.loadedVoice = ""
		p.status.TTSReady = false
		if outcome.ready {
			p.ttsSynth = synth
			p.loadedVoice = loadedVoice
			p.status.TTSReady = true
		}
		p.mu.Unlock()
	}
	if config.MultimodalEnabled {
		if err := p.initMultimodal(ctx, plan.Multimodal); err != nil {
			return errors.Join(p.initializationResult().aggregateError(), err)
		}
	}

	return p.initializationResult().aggregateError()
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

	case "audio_output":
		if ttsSynth == nil {
			return nil, types.NewAudioError(types.ErrCodeModelNotFound, "tts is not enabled", nil)
		}
		// If voice override requests a different model, lazily switch
		if msg.VoiceOverride != nil && msg.VoiceOverride.Voice != "" {
			if err := p.ensureTTSVoice(ctx, msg.VoiceOverride.Voice); err != nil {
				fmt.Printf("WARN: voice switch failed (%v), using current voice\n", err)
			}
		}
		audioPath, err := p.synthesizeWithCache(ctx, cacheStore, ttsSynth, msg.Text, msg.VoiceOverride)
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

// GetTTSVoices returns available voices from the current TTS synthesizer.
func (p *Pipeline) GetTTSVoices() []types.Voice {
	p.mu.RLock()
	synth := p.ttsSynth
	p.mu.RUnlock()
	if synth == nil {
		return nil
	}
	return synth.GetVoices()
}

// Shutdown releases loaded resources through loader lifecycle.
func (p *Pipeline) Shutdown(ctx context.Context) error {
	_ = ctx
	p.lifecycleMu.Lock()
	defer p.lifecycleMu.Unlock()

	ids := p.loader.GetLoadedModels()
	var errs []error
	for _, id := range ids {
		if err := p.loader.Unload(id); err != nil {
			errs = append(errs, fmt.Errorf("unload %s: %w", id, err))
		}
	}

	p.mu.Lock()
	p.asrRecognizer = nil
	p.ttsSynth = nil
	p.multimodal = nil
	p.loadedVoice = ""
	p.status.ASRReady = false
	p.status.TTSReady = false
	p.status.MultimodalReady = false
	p.mu.Unlock()

	return errors.Join(errs...)
}

// pipelineCapabilityCatalog is the contract the pipeline enforces for both
// automatic and explicit selection. Package tests replace it so lifecycle and
// provider outcomes can be exercised without installed audio assets;
// production always reads the probed environment.
var pipelineCapabilityCatalog = func() []capability.Record {
	return capability.Catalog(capability.ProductionEnv(ASRModelPathFor, ASRProjectorPathFor, TTSModelPathFor))
}

// asrCandidatesFor resolves the engines one Apply may load.
//
// Automatic selection stays inside local, admitted models: a remote API must
// never be reached because a local model was missing. An explicit engine is
// honoured exactly, or refused with the capability contract's blockers; it is
// never replaced by a different provider.
func asrCandidatesFor(planned []types.ASRModelType, engine string) ([]types.ASRModelType, error) {
	return asrCandidatesForRecords(planned, engine, pipelineCapabilityCatalog())
}

func asrCandidatesForRecords(planned []types.ASRModelType, engine string, records []capability.Record) ([]types.ASRModelType, error) {
	engine = strings.TrimSpace(engine)
	if engine == "" || engine == "auto" {
		filtered := make([]capability.Record, 0, len(planned))
		reasons := make([]string, 0, len(planned))
		for _, candidate := range planned {
			profile, ok := types.GetASRProfile(candidate)
			if !ok || !profile.Local || !profile.ApprovedLocal {
				continue
			}
			record, found := findCapabilityRecord(records, capability.TaskFileTranscription, string(candidate))
			if !found {
				reasons = append(reasons, string(candidate)+": no capability record")
				continue
			}
			filtered = append(filtered, record)
			if !record.Selectable() {
				reasons = append(reasons, string(candidate)+": "+strings.Join(record.Blockers(), "; "))
			}
		}
		primary, err := capability.AutoSelect(filtered, capability.TaskFileTranscription)
		if err != nil {
			if len(reasons) == 0 {
				return nil, fmt.Errorf("no local admitted ASR engine matches the current constraints")
			}
			return nil, fmt.Errorf("no selectable local ASR engine: %s", strings.Join(reasons, " | "))
		}
		out := []types.ASRModelType{types.ASRModelType(primary.ID)}
		for _, fallback := range capability.Fallbacks(filtered, primary) {
			out = append(out, types.ASRModelType(fallback.ID))
		}
		return out, nil
	}

	record, err := capability.Select(records, capability.TaskFileTranscription, engine)
	if err != nil {
		return nil, fmt.Errorf("explicit ASR engine %q cannot be used: %w", engine, err)
	}
	return []types.ASRModelType{types.ASRModelType(record.ID)}, nil
}

func findCapabilityRecord(records []capability.Record, task capability.Task, id string) (capability.Record, bool) {
	for _, record := range records {
		if record.Task == task && record.ID == id {
			return record, true
		}
	}
	return capability.Record{}, false
}

// ttsCandidatesFor resolves the synthesis engines one Apply may load. An
// explicit engine goes through the same capability contract as ASR: it is
// honoured with its evidence or refused with its blockers, never swapped for
// the planner's first choice.
func ttsCandidatesFor(planned []types.TTSModelType, engine string) ([]types.TTSModelType, error) {
	return ttsCandidatesForRecords(planned, engine, pipelineCapabilityCatalog())
}

func ttsCandidatesForRecords(planned []types.TTSModelType, engine string, records []capability.Record) ([]types.TTSModelType, error) {
	engine = strings.TrimSpace(engine)
	if engine == "" || engine == "auto" {
		return planned, nil
	}
	record, err := capability.Select(records, capability.TaskSpeechSynthesis, engine)
	if err != nil {
		return nil, fmt.Errorf("explicit TTS engine %q cannot be used: %w", engine, err)
	}
	return []types.TTSModelType{types.TTSModelType(record.ID)}, nil
}

func (p *Pipeline) initASR(ctx context.Context, candidates []types.ASRModelType, config types.PipelineConfig) (types.ASRRecognizer, capabilityInit) {
	out := capabilityInit{requested: true}
	if len(candidates) == 0 {
		out.err = types.NewAudioError(types.ErrCodeModelNotFound, "no ASR candidates for current constraints", nil)
		return nil, out
	}
	selected := chooseASR(candidates)
	recognizers := make([]types.ASRRecognizer, 0, len(selected))
	recMap := make(map[types.ASRModelType]types.ASRRecognizer)
	ownedIDs := make([]string, 0, len(selected))
	var errs []error
	for _, modelType := range selected {
		rec, err := p.loader.LoadASR(ctx, string(modelType), types.ASRConfig{
			ModelPath:     ASRModelPathFor(modelType),
			ProjectorPath: ASRProjectorPathFor(modelType),
			ServerPath:    ASRServerPathFor(),
			Language:      "auto",
			NumThreads:    6,
		})
		if err != nil {
			errs = append(errs, fmt.Errorf("load ASR %s: %w", modelType, err))
			continue
		}
		ownedIDs = append(ownedIDs, asrLoaderID(string(modelType)))
		recognizers = append(recognizers, rec)
		recMap[modelType] = rec
	}
	if len(recognizers) == 0 {
		out.err = errors.Join(errs...)
		if out.err == nil {
			out.err = types.NewAudioError(types.ErrCodeModelLoadFailed, "failed to initialize any ASR recognizer", nil)
		}
		return nil, out
	}
	chain := fallback.NewChain(fallback.ChainConfig{
		Recognizers:          recognizers,
		RecognizerMap:        recMap,
		MaxAttempts:          1,
		AutoSelectByHardware: true,
		MinMemoryTier:        parseTier(config.HardwareTier),
		Language:             "auto",
	})
	if err := chain.Initialize(ctx, types.ASRConfig{Language: "auto"}); err != nil {
		cleanupErrs := []error{fmt.Errorf("initialize ASR chain: %w", err)}
		for _, id := range ownedIDs {
			if unloadErr := p.loader.Unload(id); unloadErr != nil {
				cleanupErrs = append(cleanupErrs, fmt.Errorf("unload %s: %w", id, unloadErr))
			}
		}
		out.err = errors.Join(cleanupErrs...)
		return nil, out
	}
	out.ready = true
	return chain, out
}

func (p *Pipeline) initTTS(ctx context.Context, candidates []types.TTSModelType) (types.TTSSynthesizer, string, capabilityInit) {
	out := capabilityInit{requested: true}
	if len(candidates) == 0 {
		out.err = types.NewAudioError(types.ErrCodeModelNotFound, "no TTS candidates for current constraints", nil)
		return nil, "", out
	}
	cacheDir := defaultTTSCacheDir()
	var errs []error
	for _, candidate := range candidates {
		cfg := ttsConfigForModel(candidate, cacheDir)
		synth, err := p.loader.LoadTTS(ctx, string(candidate), cfg)
		if err != nil {
			errs = append(errs, fmt.Errorf("load TTS %s: %w", candidate, err))
			continue
		}
		out.ready = true
		return synth, cfg.ModelPath, out
	}
	out.err = errors.Join(errs...)
	if out.err == nil {
		out.err = types.NewAudioError(types.ErrCodeModelLoadFailed, "failed to initialize any configured TTS engine", nil)
	}
	return nil, "", out
}

// ensureTTSVoice lazily loads a different Piper voice model if the requested
// voice differs from the currently loaded one. Returns nil if already loaded.
func (p *Pipeline) ensureTTSVoice(ctx context.Context, voiceName string) error {
	p.lifecycleMu.RLock()
	defer p.lifecycleMu.RUnlock()
	p.replacementMu.Lock()
	defer p.replacementMu.Unlock()

	p.mu.RLock()
	synth := p.ttsSynth
	loadedVoice := p.loadedVoice
	p.mu.RUnlock()
	if synth == nil || loadedVoice == voiceName {
		return nil
	}

	modelPath := voiceModelPath(voiceName)
	if _, err := os.Stat(modelPath); err != nil {
		return types.NewAudioError(types.ErrCodeModelNotFound,
			fmt.Sprintf("voice model not found: %s", modelPath), err)
	}

	cacheDir := defaultTTSCacheDir()
	cfg := types.TTSConfig{
		ModelPath:    modelPath,
		CacheDir:     cacheDir,
		CacheEnabled: true,
	}
	newSynth, err := p.loader.LoadTTS(ctx, string(types.TTSPiper), cfg)
	if err != nil {
		return types.NewAudioError(types.ErrCodeModelLoadFailed,
			fmt.Sprintf("failed to load voice %s", voiceName), err)
	}
	p.mu.Lock()
	p.ttsSynth = newSynth
	p.loadedVoice = voiceName
	p.mu.Unlock()
	fmt.Printf("INFO: loaded TTS voice: %s\n", voiceName)
	return nil
}

func (p *Pipeline) initializationResult() initOutcome {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.initialization
}

func asrLoaderID(modelType string) string {
	return "asr:" + modelType
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

func (p *Pipeline) synthesizeWithCache(ctx context.Context, cacheStore types.AudioCache, synth types.TTSSynthesizer, text string, override *types.VoiceOverride) (string, error) {
	voice := firstVoiceOrDefault(synth, "default")
	req := types.SynthesisRequest{
		Text:  text,
		Voice: voice,
		Speed: 1.0,
	}
	if override != nil {
		if override.Voice != "" {
			req.Voice = override.Voice
		}
		if override.Speed > 0 {
			req.Speed = override.Speed
		}
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
	loader.RegisterASR(types.ASRQwen3ASR17BLocal, func() types.ASRRecognizer { return llama.New() })
	loader.RegisterTTS(types.TTSPiper, func() types.TTSSynthesizer { return piper.New() })
	loader.RegisterTTS(types.TTSMeloTTS, func() types.TTSSynthesizer { return melotts.New() })
	loader.RegisterTTS(types.TTSKokoro, func() types.TTSSynthesizer { return kokoro.New() })
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
	return "models/en_US-lessac-medium.onnx"
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
	case types.TTSKokoro:
		modelPath := defaultKokoroModelPath()
		return types.TTSConfig{
			ModelPath:    modelPath,
			CacheDir:     cacheDir,
			CacheEnabled: true,
			MaxCacheSize: 100,
			ScriptPath:   defaultKokoroScriptPath(),
			VoicesPath:   defaultKokoroVoicesPath(),
			ConfigPath:   defaultKokoroConfigPath(),
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

func defaultKokoroScriptPath() string {
	candidates := []string{
		"scripts/kokoro-tts.py",
		"../scripts/kokoro-tts.py",
		"../../scripts/kokoro-tts.py",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	return "scripts/kokoro-tts.py"
}

// kokoroSearchDirs lists where a Kokoro model set is looked for, in priority
// order. It is rooted at the project rather than at a fixed drive so the paths
// survive the project moving, matching how the rest of the app resolves files.
func kokoroSearchDirs() []string {
	root := paths.ProjectRoot()
	return []string{
		filepath.Join(root, "models", "kokoro"),
		filepath.Join(root, "models"),
		filepath.Join(root, "runtime", "tts", "kokoro"),
		filepath.Join(root, "runtime", "tts"),
	}
}

// findKokoroFile returns the first existing match for name across the search
// directories, or "" when the file is not installed. Callers surface an empty
// path as a clear "model not found" error instead of failing on a stale path.
func findKokoroFile(name string) string {
	for _, dir := range kokoroSearchDirs() {
		candidate := filepath.Join(dir, name)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func defaultKokoroModelPath() string {
	if v := os.Getenv("HOWL_TTS_KOKORO_MODEL"); v != "" {
		return v
	}
	return findKokoroFile("kokoro-v1.0.onnx")
}

func defaultKokoroVoicesPath() string {
	if v := os.Getenv("HOWL_TTS_KOKORO_VOICES"); v != "" {
		return v
	}
	return findKokoroFile("voices-v1.0.bin")
}

func defaultKokoroConfigPath() string {
	if v := os.Getenv("HOWL_TTS_KOKORO_CONFIG"); v != "" {
		return v
	}
	return findKokoroFile(filepath.Join("hf_cache", "config.json"))
}

func firstVoiceOrDefault(synth types.TTSSynthesizer, fallback string) string {
	voices := synth.GetVoices()
	if len(voices) > 0 && voices[0].ID != "" {
		return voices[0].ID
	}
	return fallback
}

// voiceModelPath resolves a voice name (e.g. "en_US-lessac-medium") to its
// Piper .onnx model file path. Looks in models/ first, then checks env.
func voiceModelPath(voiceName string) string {
	if v := os.Getenv("HOWL_TTS_PIPER_MODEL"); v != "" && strings.Contains(v, voiceName) {
		return v
	}
	return filepath.Join("models", voiceName+".onnx")
}

func defaultASRModelPath(model types.ASRModelType) string {
	switch model {
	case types.ASRWhisperTiny:
		return "models/whisper/ggml-tiny.bin"
	case types.ASRWhisperSmall:
		return "models/whisper/ggml-small.bin"
	case types.ASRWhisperBase:
		return "models/whisper/ggml-base.bin"
	case types.ASRQwen3ASR06B, types.ASRQwen3ASR17B:
		// Remote DashScope adapters keep their own model name; there is no
		// local file to point at.
		return ""
	default:
		return string(model)
	}
}
