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
	ttsLoadMu     sync.Mutex
	mu            sync.RWMutex
	config        types.PipelineConfig
	status        types.PipelineStatus
	loader        *runtime.ModelLoader
	planner       runtime.Constraints

	asrRecognizer   types.ASRRecognizer
	ttsSynth        types.TTSSynthesizer
	ttsEngines      map[string]types.TTSSynthesizer
	activeTTSEngine string
	multimodal      types.MultimodalBridge
	ttsCache        types.AudioCache
	ttsCacheMu      sync.Mutex
	ttsCacheCalls   map[string]*ttsCacheCall
	loadedVoice     string // Configured voice model for the active synthesizer
	initialization  initOutcome
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

// ttsCacheCall coalesces concurrent requests with the same complete synthesis
// identity. Followers receive the leader's result instead of invoking the
// engine again.
type ttsCacheCall struct {
	done chan struct{}
	path string
	err  error
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
		loader:        loader,
		ttsEngines:    make(map[string]types.TTSSynthesizer),
		ttsCacheCalls: make(map[string]*ttsCacheCall),
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
	p.ttsCacheMu.Lock()
	p.ttsCacheCalls = make(map[string]*ttsCacheCall)
	p.ttsCacheMu.Unlock()
	p.ttsEngines = make(map[string]types.TTSSynthesizer)
	p.activeTTSEngine = ""
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
		var ttsEngine string
		outcome := capabilityInit{requested: true}
		candidates, selectionErr := ttsCandidatesFor(plan.TTS, config.TTSEngine)
		if selectionErr != nil {
			outcome.err = selectionErr
		} else {
			synth, loadedVoice, ttsEngine, outcome = p.initTTS(ctx, candidates)
		}
		p.mu.Lock()
		p.initialization.tts = outcome
		p.ttsSynth = nil
		p.ttsEngines = make(map[string]types.TTSSynthesizer)
		p.activeTTSEngine = ""
		p.loadedVoice = ""
		p.status.TTSReady = false
		if outcome.ready {
			p.ttsSynth = synth
			p.loadedVoice = loadedVoice
			p.activeTTSEngine = ttsEngine
			p.ttsEngines[ttsEngineKey(ttsEngine, "")] = synth
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
		resolved, err := p.resolveTTSRequest(ctx, msg.VoiceOverride)
		if err != nil {
			p.recordError()
			return nil, err
		}
		audioPath, err := p.synthesizeWithCache(ctx, cacheStore, resolved, msg.Text)
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

// resolvedTTS is the request-local result of voice/engine resolution. It is
// never installed into Pipeline.config or Pipeline.ttsSynth.
type resolvedTTS struct {
	synth         types.TTSSynthesizer
	engine        string
	voice         string
	voicePath     string
	modelIdentity string
	language      string
	speed         float64
	characterID   string
}

// resolveTTSRequest resolves one request against the active global binding or
// an explicitly selected engine. Explicit engines are loaded into a separate
// engine map; they never replace the active global synthesizer. Automatic or
// blank engine/voice values inherit the active global configuration.
func (p *Pipeline) resolveTTSRequest(ctx context.Context, override *types.VoiceOverride) (resolvedTTS, error) {
	p.mu.RLock()
	activeSynth := p.ttsSynth
	activeEngine := p.activeTTSEngine
	configuredEngine := p.config.TTSEngine
	configuredVoice := p.config.TTSVoice
	configuredSpeed := p.config.TTSSpeed
	p.mu.RUnlock()

	if activeSynth == nil {
		return resolvedTTS{}, types.NewAudioError(types.ErrCodeModelNotFound, "tts is not enabled", nil)
	}
	if activeEngine == "" {
		activeEngine = configuredEngine
	}

	engine := ""
	speedEngine := ""
	inheritEngine := true
	characterID := ""
	voice := ""
	if override != nil {
		characterID = strings.TrimSpace(override.CharacterID)
		voice = strings.TrimSpace(override.Voice)
		requestedEngine := strings.ToLower(strings.TrimSpace(override.Engine))
		if requestedEngine != "" && requestedEngine != "auto" {
			inheritEngine = false
			engine = requestedEngine
			speedEngine = engine
			if err := validateTTSEngine(engine); err != nil {
				return resolvedTTS{}, wrapCharacterVoiceError(characterID, err)
			}
		}
	}
	if inheritEngine {
		engine = activeEngine
		if voice == "" {
			voice = strings.TrimSpace(configuredVoice)
		}
		if engine == "auto" {
			// The provider selected at initialization is active, but no
			// engine-specific speed contract is available for "auto".
			engine = ""
		}
		// Preserve Group 2A's auto-select contract: the configured engine,
		// not the provider chosen at initialization, decides whether the
		// global speed is applied.
		if configuredEngine != "" && configuredEngine != "auto" {
			speedEngine = configuredEngine
		}
	}

	synth := activeSynth
	if !inheritEngine && engine != activeEngine {
		loaded, err := p.loadTTSForRequest(ctx, engine, voice)
		if err != nil {
			return resolvedTTS{}, wrapCharacterVoiceError(characterID, err)
		}
		synth = loaded
	}

	voicePath, err := ttsVoicePath(engine, voice)
	if err != nil {
		return resolvedTTS{}, wrapCharacterVoiceError(characterID, err)
	}
	modelIdentity := resolvedTTSModelIdentity(synth, engine, voicePath)
	if voice == "" {
		voice = firstVoiceOrDefault(synth, "default")
	}
	if engine != "piper" {
		if err := validateTTSVoice(synth, engine, voice); err != nil {
			return resolvedTTS{}, wrapCharacterVoiceError(characterID, err)
		}
	}
	language := ""
	if engine == string(types.TTSMeloTTS) {
		// MeloTTS receives the speaker and language together. Its voice IDs
		// are language names (EN/ZH/JP/...), so pass the selected voice as
		// the language hint instead of allowing the adapter to default to EN.
		language = voice
	}

	speed := ttsSpeedOrDefault(configuredSpeed)
	if override != nil && override.Speed > 0 {
		// Group 2A precedence: a request speed wins over the global speed.
		// It must still pass the same engine-specific guard as a configured
		// speed. An automatic selection has no engine-specific range, so its
		// established pass-through behavior is retained.
		speed = override.Speed
		if speed < 0.5 || speed > 2.0 {
			return resolvedTTS{}, wrapCharacterVoiceError(characterID,
				types.NewAudioError(types.ErrCodeSynthesisFailed,
					fmt.Sprintf("speech speed must be between 0.5 and 2.0, got %g", speed), nil))
		}
		if engine != "" && engine != "auto" {
			speed, err = ttsSpeedForConfig(engine, speed)
			if err != nil {
				return resolvedTTS{}, wrapCharacterVoiceError(characterID, err)
			}
		}
	} else {
		speed, err = ttsSpeedForConfig(speedEngine, speed)
		if err != nil {
			return resolvedTTS{}, wrapCharacterVoiceError(characterID, err)
		}
	}

	return resolvedTTS{
		synth:         synth,
		engine:        engine,
		voice:         voice,
		voicePath:     voicePath,
		modelIdentity: modelIdentity,
		language:      language,
		speed:         speed,
		characterID:   characterID,
	}, nil
}

func wrapCharacterVoiceError(characterID string, err error) error {
	if characterID == "" {
		return err
	}
	return fmt.Errorf("character %q voice: %w", characterID, err)
}

func validateTTSEngine(engine string) error {
	switch engine {
	case "", "auto":
		return nil
	case "browser":
		return types.NewAudioError(types.ErrCodeFormatNotSupported,
			"browser TTS is a renderer-only engine and cannot be synthesized by the backend", nil)
	}
	if _, ok := types.GetTTSProfile(types.TTSModelType(engine)); !ok {
		return types.NewAudioError(types.ErrCodeModelNotFound,
			fmt.Sprintf("unsupported TTS engine %q", engine), nil)
	}
	return nil
}

func resolvedTTSModelIdentity(synth types.TTSSynthesizer, engine, voicePath string) string {
	cfg := types.TTSConfig{}
	if provider, ok := synth.(interface{ GetConfig() types.TTSConfig }); ok {
		cfg = provider.GetConfig()
	}
	return types.TTSModelIdentity(engine, cfg, voicePath)
}

func validateTTSVoice(synth types.TTSSynthesizer, engine, voice string) error {
	if voice == "" || synth == nil {
		return nil
	}
	voices := synth.GetVoices()
	if len(voices) == 0 {
		// Some adapters enumerate voices dynamically. Let the adapter's own
		// synthesis contract report an invalid voice rather than guessing.
		return nil
	}
	for _, candidate := range voices {
		if candidate.ID == voice || candidate.Name == voice {
			return nil
		}
	}
	return types.NewAudioError(types.ErrCodeSynthesisFailed,
		fmt.Sprintf("voice %q is not supported by TTS engine %q", voice, engine), nil)
}

// loadTTSForRequest loads an explicit engine without changing the active global
// binding. Piper is keyed by its request model path so a later character cannot
// inherit another Piper voice model by accident.
func (p *Pipeline) loadTTSForRequest(ctx context.Context, engine, voice string) (types.TTSSynthesizer, error) {
	p.lifecycleMu.RLock()
	defer p.lifecycleMu.RUnlock()
	p.ttsLoadMu.Lock()
	defer p.ttsLoadMu.Unlock()

	if p.loader == nil {
		return nil, types.NewAudioError(types.ErrCodeModelNotFound,
			fmt.Sprintf("cannot load TTS engine %q without a model loader", engine), nil)
	}
	voicePath, err := ttsVoicePath(engine, voice)
	if err != nil {
		return nil, err
	}
	key := ttsEngineKey(engine, voicePath)
	p.mu.RLock()
	if synth, ok := p.ttsEngines[key]; ok {
		p.mu.RUnlock()
		return synth, nil
	}
	p.mu.RUnlock()

	cfg := ttsConfigForModel(types.TTSModelType(engine), defaultTTSCacheDir())
	if voicePath != "" {
		cfg.ModelPath = voicePath
	}
	synth, err := p.loader.LoadTTS(ctx, engine, cfg)
	if err != nil {
		return nil, types.NewAudioError(types.ErrCodeModelLoadFailed,
			fmt.Sprintf("load TTS engine %q: %v", engine, err), err)
	}
	p.mu.Lock()
	if p.ttsEngines == nil {
		p.ttsEngines = make(map[string]types.TTSSynthesizer)
	}
	p.ttsEngines[key] = synth
	p.mu.Unlock()
	return synth, nil
}

func ttsEngineKey(engine, voicePath string) string {
	if engine == "piper" && voicePath != "" {
		return engine + "|" + voicePath
	}
	return engine
}

func ttsVoicePath(engine, voice string) (string, error) {
	if engine != "piper" || voice == "" {
		return "", nil
	}
	path := voiceModelPath(voice)
	if _, err := os.Stat(path); err != nil {
		return "", types.NewAudioError(types.ErrCodeModelNotFound,
			fmt.Sprintf("piper voice model not found: %s", path), err)
	}
	return path, nil
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
	p.ttsEngines = make(map[string]types.TTSSynthesizer)
	p.activeTTSEngine = ""
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

// Non-Whisper ASR keeps the values this pipeline used before decode settings
// became operational. In particular, llama-server still receives its legacy
// thread count and language value; the four desired decode values must not leak
// into an engine whose adapter does not implement them.
const (
	legacyASRLanguage = "auto"
	legacyASRThreads  = 6
)

// asrDecodeConfigFor forwards Whisper decode options only for an explicit
// Whisper selection. Automatic and non-Whisper engines keep their existing
// adapter values, and the status report marks the four values unsupported or
// unavailable rather than pretending they were consumed.
func asrDecodeConfigFor(config types.PipelineConfig) types.ASRConfig {
	if !isWhisperASREngine(config.ASREngine) {
		return types.ASRConfig{Language: legacyASRLanguage, NumThreads: legacyASRThreads}
	}
	return types.ASRConfig{
		Language:   config.ASRLanguage,
		NumThreads: config.ASRThreads,
		BeamSize:   config.ASRBeamSize,
		BestOf:     config.ASRBestOf,
	}
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
	decode := asrDecodeConfigFor(config)
	for _, modelType := range selected {
		rec, err := p.loader.LoadASR(ctx, string(modelType), types.ASRConfig{
			ModelPath:     ASRModelPathFor(modelType),
			ProjectorPath: ASRProjectorPathFor(modelType),
			ServerPath:    ASRServerPathFor(),
			Language:      decode.Language,
			NumThreads:    decode.NumThreads,
			BeamSize:      decode.BeamSize,
			BestOf:        decode.BestOf,
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
		Language:             decode.Language,
	})
	if err := chain.Initialize(ctx, types.ASRConfig{Language: decode.Language}); err != nil {
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

func (p *Pipeline) initTTS(ctx context.Context, candidates []types.TTSModelType) (types.TTSSynthesizer, string, string, capabilityInit) {
	out := capabilityInit{requested: true}
	if len(candidates) == 0 {
		out.err = types.NewAudioError(types.ErrCodeModelNotFound, "no TTS candidates for current constraints", nil)
		return nil, "", "", out
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
		return synth, cfg.ModelPath, string(candidate), out
	}
	out.err = errors.Join(errs...)
	if out.err == nil {
		out.err = types.NewAudioError(types.ErrCodeModelLoadFailed, "failed to initialize any configured TTS engine", nil)
	}
	return nil, "", "", out
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

// ttsSpeedForConfig applies a configured speed only to an explicit engine whose
// profile declares both speed support and a finite range. Automatic selection
// keeps the established 1.0 default because the effective provider is resolved
// later by the loader, after this configuration is recorded.
func ttsSpeedForConfig(engine string, speed float64) (float64, error) {
	speed = ttsSpeedOrDefault(speed)
	engine = strings.TrimSpace(engine)
	if engine == "" || engine == "auto" {
		return 1.0, nil
	}
	min, max, supported, explicit := ttsSpeedRangeFor(engine)
	if !explicit {
		return 1.0, nil
	}
	if !supported {
		return 0, types.NewAudioError(types.ErrCodeSynthesisFailed,
			fmt.Sprintf("TTS engine %q does not support speech speed", engine), nil)
	}
	if speed < min || speed > max {
		return 0, types.NewAudioError(types.ErrCodeSynthesisFailed,
			fmt.Sprintf("speech speed for %s must be between %g and %g, got %g", engine, min, max, speed), nil)
	}
	return speed, nil
}

func (p *Pipeline) synthesizeWithCache(ctx context.Context, cacheStore types.AudioCache, resolved resolvedTTS, text string) (string, error) {
	req := types.SynthesisRequest{
		Text:          text,
		Voice:         resolved.voice,
		Language:      resolved.language,
		Speed:         ttsSpeedOrDefault(resolved.speed),
		Engine:        resolved.engine,
		CharacterID:   resolved.characterID,
		VoicePath:     resolved.voicePath,
		ModelIdentity: resolved.modelIdentity,
	}
	key := types.GenerateCacheKeyForRequest(req)
	call, leader := p.beginTTSCacheCall(key)
	if !leader {
		select {
		case <-call.done:
			return call.path, call.err
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	// The leader performs the single cache lookup for this identity. Followers
	// wait for that result instead of racing another lookup or synthesis.
	if cached, ok := cacheStore.Get(key); ok {
		p.finishTTSCacheCall(key, call, cached.AudioPath, nil)
		return cached.AudioPath, nil
	}

	out, err := resolved.synth.Synthesize(ctx, req, nil)
	if err != nil {
		p.finishTTSCacheCall(key, call, "", err)
		return "", err
	}
	if err := cacheStore.Store(key, out); err != nil {
		p.finishTTSCacheCall(key, call, "", err)
		return "", err
	}
	p.finishTTSCacheCall(key, call, out.AudioPath, nil)
	return out.AudioPath, nil
}

func (p *Pipeline) beginTTSCacheCall(key string) (*ttsCacheCall, bool) {
	p.ttsCacheMu.Lock()
	defer p.ttsCacheMu.Unlock()
	if p.ttsCacheCalls == nil {
		p.ttsCacheCalls = make(map[string]*ttsCacheCall)
	}
	if call, ok := p.ttsCacheCalls[key]; ok {
		return call, false
	}
	call := &ttsCacheCall{done: make(chan struct{})}
	p.ttsCacheCalls[key] = call
	return call, true
}

func (p *Pipeline) finishTTSCacheCall(key string, call *ttsCacheCall, path string, err error) {
	p.ttsCacheMu.Lock()
	call.path = path
	call.err = err
	delete(p.ttsCacheCalls, key)
	p.ttsCacheMu.Unlock()
	close(call.done)
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
	if paths.IsolationActive() {
		return paths.StatePath("cache", "tts")
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
			ModelName:    meloOpenAIModelName(),
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

func meloOpenAIModelName() string {
	if v := strings.TrimSpace(os.Getenv("HOWL_MELOTTS_OPENAI_MODEL")); v != "" {
		return v
	}
	return "tts-1"
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
