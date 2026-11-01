package main

import (
	"context"
	"fmt"
	"howl-chat/internal/proc"
	"os/exec"
	"sync"
	"sync/atomic"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"howl-chat/internal/audio"
	audiotypes "howl-chat/internal/audio/types"
	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/backend/memory"
	"howl-chat/internal/backend/model"
	"howl-chat/internal/backend/rag"
	"howl-chat/internal/backend/types"
	"howl-chat/internal/errors"
	imageengine "howl-chat/internal/image/engine"
	"howl-chat/internal/image/gallery"
	"howl-chat/internal/image/job"
	"howl-chat/internal/image/registry"
	imgsettings "howl-chat/internal/image/settings"
	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/paths"
	appsettings "howl-chat/internal/settings"
)

// App represents Wails application
// Contains application context and backend services
type App struct {
	ctx             context.Context
	httpClient      *llama.HTTPClient
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
	streamCancel    context.CancelFunc
	chatPersistence *chat.PersistenceManager

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

	// Audio Pipeline Integration
	audioPipeline *audio.Pipeline

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
	// Kill model servers left behind by a crashed or force-killed previous instance.
	// Deferred to a goroutine: this shells out to PowerShell (~0.3-0.5s measured)
	// and nothing needs it before the window exists, so it must not sit on the
	// pre-window startup path where every millisecond delays the first paint.
	go proc.ReapOrphans(paths.ProjectRoot())

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
	ragService = rag.NewService(paths.ProjectRoot(), ragEmbedder)
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
	}

	// Initialize unified settings store
	app.settingsStore = appsettings.NewStore(app.getSettingsDir())
	loadedSettings, err := app.settingsStore.Load()
	if err != nil {
		fmt.Printf("Warning: failed to load unified settings: %v\n", err)
		loadedSettings = make(map[string]interface{})
	}
	app.appSettings = loadedSettings

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
	a.initAudioPipeline()
	startupMark("OnStartup: init complete")

	// Model scanning is background work: the window is already visible and the
	// user can interact, so a slow or unplugged model directory must never
	// delay the first paint.
	go a.RescanAllDirs()
}

// initAudioPipeline creates and initializes the audio pipeline if binaries are available.
// It is best-effort: if whisper-cli or piper are not installed, audio stays disabled.
func (a *App) initAudioPipeline() {
	pipe := audio.NewPipeline()
	cfg := audiotypes.PipelineConfig{
		ASREnabled:        false, // User must enable via settings
		TTSEnabled:        false, // User must enable via settings
		MultimodalEnabled: false,
		RealtimeEnabled:   false,
		HardwareTier:      "auto",
	}
	if err := pipe.Initialize(context.Background(), cfg); err != nil {
		fmt.Printf("Audio pipeline init skipped: %v\n", err)
		return
	}
	a.audioPipeline = pipe
	fmt.Println("Audio pipeline initialized successfully")
}

// OnShutdown is called when the application is shutting down
// Performs cleanup operations before application termination
func (a *App) OnShutdown(ctx context.Context) {
	// Save current chat before shutting down
	if a.chatPersistence != nil {
		a.performFinalSave()
	}

	// Stop the model server (and anything it spawned) before anything else,
	// so port 8080 and VRAM are released immediately.
	if a.serverProc != nil && a.serverProc.Process != nil {
		proc.KillTree(a.serverProc.Process.Pid)
		_ = a.serverProc.Process.Kill()
		_, _ = a.serverProc.Process.Wait()
		a.serverProc = nil
	}

	// Shut down the image engine, which owns its own sd-server subprocess.
	if a.imageEngine != nil {
		_ = a.imageEngine.Shutdown()
	}

	// Shutdown audio pipeline if initialized
	if a.audioPipeline != nil {
		_ = a.audioPipeline.Shutdown(ctx)
	}

	// Reap anything still attached to the job object, so llama-server and
	// sd-server can never outlive the app window.
	proc.CloseJob()
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
