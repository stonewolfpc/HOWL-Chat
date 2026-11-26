package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/memory"
	"howl-chat/internal/backend/rag"
	"howl-chat/internal/characterstore"
	"howl-chat/internal/diagnostics"
	"howl-chat/internal/host"
	"howl-chat/internal/media"
	"howl-chat/internal/registries"
	"howl-chat/internal/resource"
	appruntime "howl-chat/internal/runtime"
)

// App represents Wails application
// Contains application context and backend services
type App struct {
	ctx context.Context

	// Runtime & Resources (internal/runtime). The engine owns the inference
	// client, the model manager and load progress/stage plus the last load's
	// profile and injected chat template; the settings authority owns the
	// unified document, its sampler and model sections and their persistence.
	// App holds the owners and the Wails bindings, not the state.
	engine   *appruntime.Engine
	settings *appruntime.Settings
	// conversation owns the active conversation: its session (turns, world,
	// scenario, characters, lorebooks, trigger bookkeeping) and its collection
	// (membership, append log, snapshots). App holds one handle to the owner and
	// the Wails bindings; it keeps none of that state in its own fields.
	conversation *chat.ActiveConversation

	// characters is the canonical character store. The renderer reads and
	// writes characters only through the bindings in app-characters.go; App
	// composes the owner and keeps no character maps of its own.
	characters  *characterstore.Store
	characterMu sync.Mutex

	// admission is the single reserve-aware gate every model start consults.
	// The pointer is the process-wide authority in production; tests replace it
	// with an isolated instance so they never consult live hardware.
	admission *resource.Authority
	// imageStartMu does the same for image loads, whose job manager can run
	// several queued starts concurrently.
	imageStartMu sync.Mutex
	// lifecycleGate serializes conflicting *transitions* across workloads — a
	// start, stop, unload, recovery or reload — so one workload cannot stop
	// another while it is starting. It is never held during execution; the
	// workload leases above own that span.
	lifecycleGate *resource.LifecycleGate
	lifecycleOnce sync.Once

	// Level 0 — application lifecycle and coordination (internal/host).
	// The authority owns admission, shutdown sequencing, background joining,
	// the active stream token and child-process supervision. It owns no domain
	// state: conversation, model, audio, image and configuration values stay
	// with their own owners. Transitional ownership in this struct is recorded
	// in docs/architecture/HOC-H0-LEVEL0-2026-11-22.md.
	authorityOnce sync.Once
	authority     *host.Authority

	// Service dependencies
	chatService *chat.Service

	// Memory Service Integration
	memoryService *memory.Service

	// RAG Service Integration
	ragService *rag.Service

	// Media domain (internal/media). Audio owns the manager, the desired
	// settings contract, the apply route and the ASR activation intent; Image
	// owns the image engine session, gallery, jobs, settings store and last
	// model. App holds the owners and the Wails bindings, not the state.
	audio  *media.Audio
	images *media.Image
	// imageReclaim is the resident image engine's post-generation policy. It
	// only decides whether to recycle the already-admitted process; admission
	// and placement remain with resource.Manager and sd.cpp respectively.
	imageReclaim imageReclaimPolicy
	// imageWorkflow is the active workflow's parent lease. It holds engine
	// identity and the component claims retained across stages. Reclamation is
	// refused while it is active; it is cleared before the terminal decision.
	imageWorkflow *resource.WorkflowLease
	// imageWorkflowRun is the cancellable execution scope for the active
	// workflow. It is separate from imageStartMu so CancelImageWorkflow can
	// request cancellation without waiting for the workflow to finish.
	imageWorkflowRunMu sync.Mutex
	imageWorkflowRun   *imageWorkflowExecution

	// registries is the one handle to the registry authorities: registered
	// model directories/types, content roots, the inventory snapshot, the
	// world/scenario import services, canonical lorebooks and the model-profile
	// cache. The state lives in those domain packages; App holds the composition
	// owner and the Wails bindings, not seven separate fields.
	registries *registries.Set

	// Content-library readiness (app-content-roots.go). These fields record what
	// the filesystem actually allows, as measured by the provisioning probe; the
	// content-library registry remains the authority for which roots exist.
	contentMu     sync.Mutex
	contentProbes map[string]ProvisionReport
	contentInit   map[string]ProvisionReport
	contentOps    *contentFSOps // tests inject failing filesystem operations

	// loaderDiagnosticsDir overrides where model-load failure records are
	// written. Empty means <project root>/logs.
	loaderDiagnosticsDir string

	// discoveryDirOverride redirects directory-discovery records for tests.
	// Empty means the normal settings directory.
	discoveryDirOverride string

	// Diagnostics (internal/diagnostics). The owner holds the structured error
	// log and its rotation, privacy-safe report building/upload, and the UI
	// event payloads; App supplies the log location, the upload preference and
	// the Wails event transport (app-errors.go).
	diagnostics *diagnostics.Owner
	// diagnosticsFallbackErr records why the configured diagnostics log could not
	// be opened, when startup fell back to a temporary location. It is reported
	// once so an unwritable log is a visible limitation, not a silent one.
	diagnosticsFallbackErr error
}

// NewApp builds the application and runs its ordered composition.
//
// Level 0 (internal/host) owns that sequence: it runs the stages from
// app-boot.go in dependency order, marks readiness, and reports the stage that
// failed. App supplies the bodies and holds the owner references; it no longer
// performs startup ordering itself, and the image session is constructed by the
// image domain rather than here.
func NewApp() *App {
	app := &App{
		engine:       appruntime.NewEngine(nil, nil),
		conversation: chat.NewActiveConversation(nil),
		admission:    resource.DefaultAuthority(),
	}
	// Children of this process are assigned to the application job at start.
	// Crash cleanup relies on KILL_ON_JOB_CLOSE, not on matching executable
	// paths: a second Howl instance must never reap this instance's children.
	if err := app.level0().Lifecycle.Boot(context.Background(), app.bootStages()); err != nil {
		fmt.Println("[Lifecycle] startup incomplete:", err)
	}
	// Diagnostics is host-independent: its log and reporter need no webview.
	// Only the event transport needs a UI, and appDiagnosticPort already drops
	// events until one is attached, so every host that composes an App gets the
	// same logging and reporting instead of only the Wails startup path.
	app.initErrorSystem()
	app.bindErrorEmitter()
	return app
}

// startup is called when the application is starting up
// Initializes the application context and optional audio pipeline
func (a *App) startup(ctx context.Context) {
	startupMark("OnStartup: webview reached backend")
	a.ctx = ctx
	lifecycle := a.level0().Lifecycle
	fmt.Println("[lifecycle] webview attached: boot", lifecycle.Epoch(), "ready", lifecycle.Ready())
	startupMark("OnStartup: init complete")

	// The startup log is optional: if it cannot be written, say so and keep
	// going. Content storage is not optional, so its state is provisioned next
	// and reported per library (app-content-roots.go); neither step fails a boot.
	a.verifyOptionalStartupLog()
	a.initializeContentRoots()

	// Model scanning is background work: the window is already visible and the
	// user can interact, so a slow or unplugged model directory must never
	// delay the first paint.
	a.level0().Background.Start("rescan", a.RescanAllDirs)
	// "Never show again" means silent discovery, not no discovery: keep the
	// pickers accurate in the background without interrupting the user.
	a.level0().Background.Start("model discovery", a.RunSilentModelDiscovery)
	startNativeTestSeam(a)
}

// OnShutdown is the Wails shutdown callback. It delegates to the one
// idempotent lifecycle sequence shared with explicit Exit.
func (a *App) OnShutdown(ctx context.Context) {
	if err := a.shutdown(ctx); err != nil {
		fmt.Printf("Application shutdown completed with errors: %v\n", err)
	}
}

// OnShutdownError runs the same sequence and returns its accumulated failure.
// Wails' shutdown callback cannot return a value, so a host that wants to report
// a failed stage (the browser harness) uses this form instead of printing and
// forgetting it.
func (a *App) OnShutdownError(ctx context.Context) error {
	return a.shutdown(ctx)
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
