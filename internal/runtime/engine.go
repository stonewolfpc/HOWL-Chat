// Package runtime owns Howl's runtime and resources: the inference session and
// its load state, and the settings authority every host reads through.
//
// It deliberately sits beside internal/resource (hardware arbitration) rather
// than inside it: resource decides what may run, runtime holds what is running
// and how its configuration is stored.
package runtime

import (
	"sync"
	"sync/atomic"

	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/model"
	"howl-chat/internal/backend/types"
)

// Client is the exact inference surface the application uses. *llama.HTTPClient
// implements it in production. It exists so a recorder can sit at the same call
// point in tests, without a second prompt-assembly route or a live model.
type Client interface {
	GenerateChat(messages []llama.ChatMessage, options *llama.InferenceOptions) (string, error)
	GenerateChatStream(messages []llama.ChatMessage, options *llama.InferenceOptions, callback llama.TokenCallback) error
	GenerateChatStreamWithImage(messages []llama.ChatMessage, imageData []byte, options *llama.InferenceOptions, callback llama.TokenCallback) error
	CancelRequest()
	GetContextSize() int
	GetLoadedModel() (map[string]interface{}, error)
	Health() error
	IsLoaded() bool
}

// ChatPlacement is the last chat load's plan and the placement the engine
// reported. It is a record for the user, not an admission input: the next load
// reads the resource ledger, never this report.
type ChatPlacement struct {
	Model            string          `json:"model"`
	PlanSource       string          `json:"plan_source"`
	PlannedGPULayers int             `json:"planned_gpu_layers"`
	PlannedContext   int             `json:"planned_context"`
	PlannedBatch     int             `json:"planned_batch"`
	PlannedVRAMMB    int64           `json:"planned_vram_mb"`
	PlannedRAMMB     int64           `json:"planned_ram_mb"`
	Actual           llama.Placement `json:"actual"`
	MatchesPlan      bool            `json:"matches_plan"`
	Message          string          `json:"message"`
}

// Engine is the runtime's engine owner: the inference client, the model
// manager built on it, load progress and stage, and the bookkeeping for the
// last load (the profile, the chat template this app injected, and the
// placement the engine reported).
//
// A nil *Engine is an absent runtime, not a second state store: reads return
// zero values and writes are dropped. That keeps a partially constructed App
// (tests, probes) working without duplicating progress fields on it.
type Engine struct {
	mu            sync.RWMutex
	client        Client
	models        *model.Manager
	profile       *types.ModelProfile
	autoTemplate  string
	chatPlacement *ChatPlacement

	// Load progress is lock-free so the progress poller never waits on a load.
	// progress stores progress*1000 as an integer.
	progress atomic.Int64
	stage    atomic.Value // stores string
}

// NewEngine builds the engine owner. Either argument may be nil; boot replaces
// the client when it has a real server to talk to.
func NewEngine(client Client, models *model.Manager) *Engine {
	return &Engine{client: client, models: models}
}

// Client returns the inference client, or nil when no session exists.
func (e *Engine) Client() Client {
	if e == nil {
		return nil
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.client
}

// SetClient installs the inference client for this engine.
func (e *Engine) SetClient(client Client) {
	if e == nil {
		return
	}
	e.mu.Lock()
	e.client = client
	e.mu.Unlock()
}

// Models returns the model manager bound to the client, or nil.
func (e *Engine) Models() *model.Manager {
	if e == nil {
		return nil
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.models
}

// SetModels installs the model manager for this engine.
func (e *Engine) SetModels(models *model.Manager) {
	if e == nil {
		return
	}
	e.mu.Lock()
	e.models = models
	e.mu.Unlock()
}

// ProgressMilli returns load progress in thousandths (0–1000).
func (e *Engine) ProgressMilli() int64 {
	if e == nil {
		return 0
	}
	return e.progress.Load()
}

// SetProgressMilli records load progress in thousandths (0–1000).
func (e *Engine) SetProgressMilli(value int64) {
	if e == nil {
		return
	}
	e.progress.Store(value)
}

// Progress returns load progress as a fraction (0.0–1.0).
func (e *Engine) Progress() float64 {
	return float64(e.ProgressMilli()) / 1000.0
}

// Stage returns the current load stage description, or "" before any load.
func (e *Engine) Stage() string {
	if e == nil {
		return ""
	}
	if v := e.stage.Load(); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// SetStage records the current load stage description.
func (e *Engine) SetStage(stage string) {
	if e == nil {
		return
	}
	e.stage.Store(stage)
}

// Profile returns the profile of the last model this app read, or nil.
func (e *Engine) Profile() *types.ModelProfile {
	if e == nil {
		return nil
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.profile
}

// SetProfile records the profile of the last model this app read.
func (e *Engine) SetProfile(profile *types.ModelProfile) {
	if e == nil {
		return
	}
	e.mu.Lock()
	e.profile = profile
	e.mu.Unlock()
}

// AutoTemplate returns the chat template this app injected into the last load.
func (e *Engine) AutoTemplate() string {
	if e == nil {
		return ""
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.autoTemplate
}

// SetAutoTemplate records the chat template this app injected into a load.
func (e *Engine) SetAutoTemplate(template string) {
	if e == nil {
		return
	}
	e.mu.Lock()
	e.autoTemplate = template
	e.mu.Unlock()
}

// ChatPlacement returns the last load's plan/actual report, or nil before a
// load has reached a healthy server.
func (e *Engine) ChatPlacement() *ChatPlacement {
	if e == nil {
		return nil
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.chatPlacement
}

// SetChatPlacement records the plan and observed placement for the last load.
// A nil report clears it, which is what an unload or a failed load does.
func (e *Engine) SetChatPlacement(placement *ChatPlacement) {
	if e == nil {
		return
	}
	e.mu.Lock()
	e.chatPlacement = placement
	e.mu.Unlock()
}
