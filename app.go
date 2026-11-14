package main

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"sync/atomic"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"howl-chat/internal/audio"
	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/backend/memory"
	"howl-chat/internal/backend/model"
	"howl-chat/internal/backend/rag"
	"howl-chat/internal/backend/types"
	"howl-chat/internal/contentimport"
	"howl-chat/internal/errors"
	imageengine "howl-chat/internal/image/engine"
	"howl-chat/internal/image/gallery"
	"howl-chat/internal/image/job"
	"howl-chat/internal/image/registry"
	imgsettings "howl-chat/internal/image/settings"
	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/library"
	"howl-chat/internal/lorebookstore"
	"howl-chat/internal/modelintel"
	"howl-chat/internal/paths"
	appsettings "howl-chat/internal/settings"
)

// inferenceClient is the exact inference surface the app uses. *llama.HTTPClient
// implements it in production. It exists so a recorder can sit at the same call
// point in tests, without a second prompt-assembly route or a live model.
type inferenceClient interface {
	GenerateChat(messages []llama.ChatMessage, options *llama.InferenceOptions) (string, error)
	GenerateChatStream(messages []llama.ChatMessage, options *llama.InferenceOptions, callback llama.TokenCallback) error
	GenerateChatStreamWithImage(messages []llama.ChatMessage, imageData []byte, options *llama.InferenceOptions, callback llama.TokenCallback) error
	CancelRequest()
	GetContextSize() int
	GetLoadedModel() (map[string]interface{}, error)
	Health() error
	IsLoaded() bool
}

// App represents Wails application
// Contains application context and backend services
type App struct {
	ctx             context.Context
	httpClient      inferenceClient
	serverProc      *exec.Cmd
	samplerSettings map[string]interface{}
	modelSettings   map[string]interface{}
	modelSettingsMu sync.RWMutex // Protects modelSettings concurrent access
	appSettings     map[string]interface{}
	settingsStore   *appsettings.Store
	chatMessages    []map[string]interface{}
	world           *chat.WorldContext
	scenario        *chat.ScenarioContext
	characters      []*chat.CharacterContext
	lorebooks       []lorebook.Entry
	triggeredLore   map[string]bool
	streamMu        sync.Mutex
	streamCancel    context.CancelFunc
	streamGen       uint64
	chatPersistence *chat.PersistenceManager

	// Lifecycle ownership. shuttingDown is read by all new-work entry points;
	// shutdownCoord makes cleanup single-entry; serverMu protects the llama
	// process and its one Wait owner.
	shuttingDown  atomic.Bool
	shutdownMu    sync.Mutex
	shutdownCoord *shutdownCoordinator
	serverMu      sync.Mutex
	serverMonitor *processMonitor

	// Model loading progress (updated by LoadModel goroutine, read by GetLoadingProgress/Stage)
	loadingProgress atomic.Int64 // stores progress*1000 as integer for lock-free reads
	loadingStage    atomic.Value // stores string

	// Service dependencies
	modelManager *model.Manager
	chatService  *chat.Service

	// Memory Service Integration
	memoryService *memory.Service

	// RAG Service Integration
	ragService *rag.Service

	// Audio Manager Integration
	audioManager *audio.Manager

	// Audio settings contract. The renderer can persist all 34 fields, while
	// audioManager applies only the four operational fields. audioDesired is the
	// authoritative request; audioApplyErr is the last explicit application
	// failure and is never mixed into the manager's effective configuration.
	audioDesiredMu  sync.RWMutex
	audioDesired    audio.DesiredSettings
	audioDesiredSet bool
	audioApplyErr   string

	// asrActivationPending records that an explicit ASR activation was requested
	// but could not complete. The microphone gate refuses before the manager
	// adopts the request, so status uses this flag to keep reporting the actual
	// microphone condition instead of falling back to "disabled".
	asrActivationPending atomic.Bool

	// Image Generation Engine
	imageEngine    imageengine.ImageEngine
	gallerySvc     *gallery.Service
	jobManager     *job.Manager
	lastImageModel string
	// autoTemplate is the chat template this app injected into the last
	// model load. It is tracked so a stale auto template can be replaced when
	// a different family loads, while a user override is left alone.
	autoTemplate       string
	lastImageModelOpts imagetypes.ModelLoadOptions
	imgSettingsStore   *imgsettings.Store
	modelRegistry      *registry.Registry
	contentLibraries   *library.Registry
	contentInventory   *library.InventoryStore
	worldImport        *contentimport.Service
	scenarioImport     *contentimport.Service
	lorebookStore      *lorebookstore.Store
	modelIntel         *modelintel.Cache

	// loaderDiagnosticsDir overrides where model-load failure records are
	// written. Empty means <project root>/logs.
	loaderDiagnosticsDir string

	// discoveryDirOverride redirects directory-discovery records for tests.
	// Empty means the normal settings directory.
	discoveryDirOverride string

	// Error Tracking System
	errorLogger   *errors.Logger
	errorReporter *errors.Reporter
	errorEmitter  *errors.CombinedSink

	// Currently loaded model profile (for capability checks like vision)
	currentModelProfile *types.ModelProfile
}

// NewApp creates a new instance of the App struct
// Initializes backend services with llama-server HTTP client
func NewApp() *App {
	// Children are assigned to the application job at start. Crash cleanup is
	// handled by KILL_ON_JOB_CLOSE, not by matching executable paths: a second
	// Howl instance must never be able to reap this instance's model servers.
	client := llama.NewHTTPClient("localhost", 8080, true)

	// Initialize model manager
	modelManager := model.NewManager(client)

	// Initialize chat persistence manager
	chatPersistence, err := chat.NewPersistenceManager()
	if err != nil {
		// Log error but don't fail app startup
		fmt.Printf("Warning: Failed to initialize chat persistence: %v\n", err)
		chatPersistence = nil
	}

	// Initialize memory service
	memoryService := memory.NewService("default") // Default chat ID for initial setup

	// Initialize chat service
	chatService, err := chat.NewService(modelManager, client, "", "", []string{}, "")
	if err != nil {
		fmt.Printf("Warning: Failed to initialize chat service: %v\n", err)
		chatService = nil
	}

	// Initialize RAG service (uses llama-server embedding endpoint if available)
	var ragService *rag.Service
	var ragEmbedder rag.Embedder
	if client != nil {
		ragEmbedder = rag.NewLlamaEmbedder(client)
	} else {
		ragEmbedder = rag.NewKeywordEmbedder()
	}
	ragService = rag.NewService(paths.StateRoot(), ragEmbedder)
	if chatService != nil {
		chatService.SetRAGService(ragService)
		chatService.SetRAGEnabled(true)
	}

	// Default sampler settings
	defaultSettings := map[string]interface{}{
		"temperature":                0.7,
		"top_p_enabled":              true,
		"top_p":                      0.9,
		"top_k":                      40,
		"min_p":                      0.05,
		"repeat_penalty":             1.1,
		"repeat_last_n":              64,
		"frequency_penalty_enabled":  true,
		"frequency_penalty":          0.0,
		"presence_penalty_enabled":   true,
		"presence_penalty":           0.0,
		"typical_p_enabled":          true,
		"typical_p":                  1.0,
		"mirostat_enabled":           true,
		"mirostat":                   0,
		"mirostat_tau":               5.0,
		"mirostat_eta":               0.1,
		"dynamic_temp_range_enabled": true,
		"dynamic_temp_range":         0.0,
		"dynamic_temp_exponent":      1.0,
		"dry_multiplier":             0.0,
		"dry_allowed_length":         2,
		"dry_base":                   1.0,
		"smoothing_factor":           0.0,
		"smoothing_curve":            1.0,
		"top_a_enabled":              true,
		"top_a":                      0.0,
		"epsilon_cutoff":             0.0,
		"eta_cutoff":                 0.0,
		"encoder_repeat_penalty":     1.0,
		"no_repeat_ngram_size":       0,
		"seed":                       -1,
	}

	app := &App{
		httpClient:      client,
		serverProc:      nil,
		samplerSettings: defaultSettings,
		modelSettings:   make(map[string]interface{}),
		chatMessages:    []map[string]interface{}{},
		chatPersistence: chatPersistence,
		modelManager:    modelManager,
		chatService:     chatService,
		memoryService:   memoryService,
		ragService:      ragService,
		audioManager:    audio.NewProductionManager(),
	}

	// Production is a no-op; a promptcapture build swaps in the test-only recorder.
	app.httpClient = wrapInferenceCapture(app.httpClient)

	// Initialize unified settings store
	app.settingsStore = appsettings.NewStore(app.getSettingsDir())
	if libs, err := library.New(app.getSettingsDir()); err != nil {
		fmt.Printf("Warning: failed to load content libraries: %v\n", err)
	} else {
		app.contentLibraries = libs
		if chatService != nil {
			if err := chatService.SetContentLibraries(libs); err != nil {
				fmt.Printf("Warning: failed to attach world/scenario store: %v\n", err)
			}
		}
		if service, err := contentimport.New(libs, library.Worlds, app.getSettingsDir()); err != nil {
			fmt.Printf("Warning: failed to initialise world import: %v\n", err)
		} else {
			app.worldImport = service
		}
		if service, err := contentimport.New(libs, library.Scenarios, app.getSettingsDir()); err != nil {
			fmt.Printf("Warning: failed to initialise scenario import: %v\n", err)
		} else {
			app.scenarioImport = service
		}
	}
	app.contentInventory = library.NewInventoryStore(app.getSettingsDir())
	app.lorebookStore = lorebookstore.New(app.contentLibraries, app.getSettingsDir())
	if intel, err := modelintel.NewCache(app.getSettingsDir()); err != nil {
		fmt.Printf("Warning: failed to load model intelligence cache: %v\n", err)
	} else {
		app.modelIntel = intel
	}
	loadedSettings, err := app.settingsStore.Load()
	if err != nil {
		fmt.Printf("Warning: failed to load unified settings: %v\n", err)
		loadedSettings = make(map[string]interface{})
	}
	app.appSettings = loadedSettings

	// Restore the full desired audio contract without activating anything. The
	// manager remains untouched until the user explicitly saves or retries.
	app.setAudioDesired(audio.DefaultDesiredSettings())
	if rawAudio, ok := loadedSettings["audio"].(map[string]interface{}); ok {
		app.setAudioDesired(audio.DesiredSettingsFromMap(rawAudio))
	}

	// Migrate legacy sampler/model settings into unified store if unified is empty
	if samplerSection, ok := loadedSettings["sampler"].(map[string]interface{}); ok && len(samplerSection) > 0 {
		for k, v := range samplerSection {
			app.samplerSettings[k] = v
		}
	} else {
		// No unified sampler section yet — load from legacy file
		if err := app.loadSamplerSettingsFromFile(); err != nil {
			fmt.Printf("Warning: failed to load sampler settings from disk: %v\n", err)
		}
	}
	if modelSection, ok := loadedSettings["model"].(map[string]interface{}); ok && len(modelSection) > 0 {
		for k, v := range modelSection {
			app.modelSettings[k] = v
		}
	} else {
		// No unified model section yet — load from legacy file
		if err := app.loadModelSettingsFromFile(); err != nil {
			fmt.Printf("Warning: failed to load model settings from disk: %v\n", err)
		}
	}

	// Initialize image engine from persisted settings (non-fatal)
	app.initImageEngine()

	app.jobManager.SetOnComplete(func(id string, result *imagetypes.ImageResult) {
		if app.gallerySvc != nil && result != nil {
			if err := app.gallerySvc.Save(result); err != nil {
				fmt.Printf("Warning: failed to save gallery from job %s: %v\n", id, err)
			} else {
				fmt.Printf("[App] Job %s result saved to gallery\n", id)
			}
		}
	})

	return app
}

// startup is called when the application is starting up
// Initializes the application context and optional audio pipeline
func (a *App) startup(ctx context.Context) {
	startupMark("OnStartup: webview reached backend")
	a.ctx = ctx
	a.initErrorSystem()
	a.bindErrorEmitter()
	startupMark("OnStartup: init complete")

	// Model scanning is background work: the window is already visible and the
	// user can interact, so a slow or unplugged model directory must never
	// delay the first paint.
	go a.RescanAllDirs()
	// "Never show again" means silent discovery, not no discovery: keep the
	// pickers accurate in the background without interrupting the user.
	go a.RunSilentModelDiscovery()
	startNativeTestSeam(a)
}

// OnShutdown is the Wails shutdown callback. It delegates to the one
// idempotent lifecycle sequence shared with explicit Exit.
func (a *App) OnShutdown(ctx context.Context) {
	if err := a.shutdown(ctx); err != nil {
		fmt.Printf("Application shutdown completed with errors: %v\n", err)
	}
}

// Chat Management Wails Bindings

// GetChatList returns a list of all chat summaries

// emitEvent is runtime.EventsEmit with a validity guard.
//
// Wails' runtime.EventsEmit calls log.Fatalf (os.Exit) when the context
// carries no event dispatcher — "An invalid context was passed". Any
// goroutine that emits outside a live webview (the UI harness, headless
// tests, or a background task that outlives window teardown) would take
// the whole process down. When there is no dispatcher the event is
// dropped instead, which is the only sane behaviour: an event nobody
// can receive is not worth dying for.
// devEventSink, when set, receives every event the app emits. Only the UI
// harness assigns it (harness_main.go), so the browser can see events that
// normally travel over the Wails bridge. It is nil in every release build.
var devEventSink func(name string, args ...interface{})

func (a *App) emitEvent(name string, data ...interface{}) {
	if devEventSink != nil {
		devEventSink(name, data...)
	}
	if a.ctx == nil || a.ctx.Value("events") == nil {
		return
	}
	runtime.EventsEmit(a.ctx, name, data...)
}
