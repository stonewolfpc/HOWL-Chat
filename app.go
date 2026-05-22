package main

import (
	"context"
	"fmt"
	"os/exec"
	"sync/atomic"

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
)

// App represents Wails application
// Contains application context and backend services
type App struct {
	ctx             context.Context
	httpClient      *llama.HTTPClient
	serverProc      *exec.Cmd
	samplerSettings map[string]interface{}
	modelSettings   map[string]interface{}
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
	imageEngine        imageengine.ImageEngine
	gallerySvc         *gallery.Service
	jobManager         *job.Manager
	lastImageModel     string
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
	ragService = rag.NewService("d:\\HOWL_Chat", ragEmbedder)
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

	// Initialize image engine from persisted settings (non-fatal)
	app.initImageEngine()

	// Load persisted settings from disk (overwrites defaults if files exist)
	if err := app.loadSamplerSettingsFromFile(); err != nil {
		fmt.Printf("Warning: failed to load sampler settings from disk: %v\n", err)
	}
	if err := app.loadModelSettingsFromFile(); err != nil {
		fmt.Printf("Warning: failed to load model settings from disk: %v\n", err)
	}

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
	a.ctx = ctx
	a.initErrorSystem()
	a.bindErrorEmitter()
	a.initAudioPipeline()
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

	if a.serverProc != nil && a.serverProc.Process != nil {
		_ = a.serverProc.Process.Kill()
		_, _ = a.serverProc.Process.Wait()
	}

	// Shutdown audio pipeline if initialized
	if a.audioPipeline != nil {
		_ = a.audioPipeline.Shutdown(ctx)
	}
}

// Chat Management Wails Bindings

// GetChatList returns a list of all chat summaries
