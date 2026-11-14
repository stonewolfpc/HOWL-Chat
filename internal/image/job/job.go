package job

import (
	"fmt"
	"sync"
	"time"

	"howl-chat/internal/image/engine"
	"howl-chat/internal/image/types"
)

// Job status values.
//
// Loading is a phase of the job, not a separate step in front of it: the engine
// can spend minutes reading weights off a slow disk, and the caller has to be
// able to see that happening and to cancel it.
const (
	statusQueued     = "queued"
	statusLoading    = "loading"
	statusGenerating = "generating"
	statusCompleted  = "completed"
	statusFailed     = "failed"
	statusCancelled  = "cancelled"
)

// Manager tracks async image generation jobs.
type Manager struct {
	mu         sync.RWMutex
	jobs       map[string]*types.ImageJob
	engine     engine.ImageEngine
	onComplete func(id string, result *types.ImageResult)
	closed     bool
}

// NewManager creates a job manager bound to an image engine.
func NewManager(eng engine.ImageEngine) *Manager {
	return &Manager{
		jobs:   make(map[string]*types.ImageJob),
		engine: eng,
	}
}

// SetEngine updates the engine reference (e.g. after re-initialization).
func (m *Manager) SetEngine(eng engine.ImageEngine) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.engine = eng
}

// SetOnComplete registers a callback fired when any job finishes successfully.
func (m *Manager) SetOnComplete(fn func(id string, result *types.ImageResult)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onComplete = fn
}

// Start registers a job, launches it, and returns its id.
//
// prepare, when non-nil, runs inside the job before generation. That is where
// engine initialization and model loading now live: run from the Wails RPC
// (where this used to happen) a slow load blocked the caller, could not be
// observed, and could not be cancelled.
func (m *Manager) Start(opts types.ImageOptions, prepare func(jobID string) error) string {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return ""
	}
	m.mu.Unlock()

	job := &types.ImageJob{
		ID:        fmt.Sprintf("job_%d", time.Now().UnixNano()),
		Status:    statusQueued,
		Progress:  0,
		Options:   opts,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return ""
	}
	m.jobs[job.ID] = job
	m.mu.Unlock()

	// Registered before the goroutine starts, so a caller that polls the returned
	// id immediately always finds the job.
	go m.run(job, prepare)
	return job.ID
}

// Status returns a snapshot of a job, or nil when the id is unknown.
//
// A copy, not the live pointer: the job is mutated by its own goroutine while
// the caller serializes it to JSON for the UI, which the race detector reports
// the moment a job is polled during generation.
func (m *Manager) Status(id string) *types.ImageJob {
	m.mu.RLock()
	defer m.mu.RUnlock()
	j, ok := m.jobs[id]
	if !ok {
		return nil
	}
	snapshot := *j
	return &snapshot
}

// Cancel marks a job as cancelled and signals the backend to abort. It works in
// every non-terminal state, including while the model is still loading.
func (m *Manager) Cancel(id string) bool {
	m.mu.Lock()
	job, ok := m.jobs[id]
	if !ok {
		m.mu.Unlock()
		return false
	}
	if job.Status == statusCompleted || job.Status == statusFailed {
		m.mu.Unlock()
		return false
	}
	job.Status = statusCancelled
	eng := m.engine
	jobID := job.ID
	m.mu.Unlock()

	// The engine call happens outside the manager lock: stopping a load can block
	// on a process, and holding the manager mutex there would stall every status
	// poll the UI makes.
	if eng != nil {
		_ = eng.CancelJob(jobID)
	}
	return true
}

// CancelAll marks every non-terminal job cancelled and asks the engine to
// abort each one. It is used during application shutdown so an in-flight load
// or generation cannot outlive the window.
func (m *Manager) CancelAll() {
	m.mu.Lock()
	m.closed = true
	ids := make([]string, 0, len(m.jobs))
	for id, job := range m.jobs {
		if job.Status == statusCompleted || job.Status == statusFailed || job.Status == statusCancelled {
			continue
		}
		job.Status = statusCancelled
		ids = append(ids, id)
	}
	eng := m.engine
	m.mu.Unlock()

	if eng != nil {
		for _, id := range ids {
			_ = eng.CancelJob(id)
		}
	}
}

// List returns all jobs newest-first.
func (m *Manager) List() []*types.ImageJob {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]*types.ImageJob, 0, len(m.jobs))
	for _, j := range m.jobs {
		list = append(list, j)
	}
	// Sort newest first
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}

// transition moves a job from one status to another, and reports whether it
// applied. Every status change that a cancellation can race with goes through
// here: a cancel that lands first changes the status, this then fails, and the
// job is never dragged out of its terminal state.
func (m *Manager) transition(job *types.ImageJob, from, to string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if job.Status != from {
		return false
	}
	job.Status = to
	return true
}

// run executes one job: load the model if needed, then generate.
func (m *Manager) run(job *types.ImageJob, prepare func(jobID string) error) {
	if m.engine == nil {
		m.finish(job, nil, fmt.Errorf("image engine not initialized"))
		return
	}

	// Loading phase. The progress poller deliberately does not run yet: it treats
	// a backend that is not answering as a dead one, and before the child exists
	// every probe fails for a reason that has nothing to do with failure.
	if prepare != nil {
		if !m.transition(job, statusQueued, statusLoading) {
			return // cancelled before it started
		}
		if err := prepare(job.ID); err != nil {
			// finish() leaves a cancellation that arrived during loading alone.
			m.finish(job, nil, err)
			return
		}
	}

	from := statusLoading
	if prepare == nil {
		from = statusQueued
	}
	if !m.transition(job, from, statusGenerating) {
		// Cancelled while loading. The load may well have succeeded; that must not
		// start generating for a job the user already abandoned.
		return
	}

	// Progress poller: updates job.Progress while generation runs.
	// Stops on first endpoint error so a missing /sdapi/v1/progress doesn't
	// waste cycles or compete with the generation request.
	progressDone := make(chan struct{})
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		failed := false
		for {
			select {
			case <-ticker.C:
				// Liveness is checked on every tick, unconditionally. It used to
				// come after the progress call and be skipped once that call had
				// errored, so the first progress failure permanently disabled the
				// one check that could end the job: sd-server could die and the job
				// would sit on "Step 8 of 8" forever with no failure recorded.
				if !m.engine.IsRunning() {
					m.mu.RLock()
					stillGenerating := job.Status == statusGenerating
					m.mu.RUnlock()
					if stillGenerating {
						m.finish(job, nil, fmt.Errorf("image engine stopped before finishing: sd-server exited. Check the model file and the engine path."))
						return
					}
				}
				if failed {
					continue
				}
				p, err := m.engine.GetProgress()
				if err != nil {
					failed = true
					continue
				}
				if p.Progress > 0 || p.Step > 0 || p.PreviewImage != "" {
					m.mu.Lock()
					if job.Status == statusGenerating {
						job.Progress = p.Progress
						job.Step = p.Step
						job.TotalSteps = p.TotalSteps
						if p.PreviewImage != "" {
							job.PreviewImage = p.PreviewImage
						}
					}
					m.mu.Unlock()
				}
			case <-progressDone:
				return
			}
		}
	}()

	// The job id travels with the request so the engine can scope a cancellation
	// to this job and refuse a second request from overlapping it.
	opts := job.Options
	opts.JobID = job.ID

	var result *types.ImageResult
	var err error

	switch opts.Mode {
	case "img2img":
		result, err = m.engine.GenerateImg2Img(opts)
	case "inpaint":
		result, err = m.engine.GenerateInpaint(opts)
	case "video":
		result, err = m.engine.GenerateVideo(opts)
	default:
		result, err = m.engine.GenerateTxt2Img(opts)
	}

	close(progressDone)
	m.finish(job, result, err)
}

// finish updates job state and cleans up if cancelled.
func (m *Manager) finish(job *types.ImageJob, result *types.ImageResult, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Terminal states are final. A cancellation or a completed result must not be
	// overwritten by whatever a late response from the engine has to say, in
	// either direction.
	if job.Status == statusCancelled || job.Status == statusCompleted {
		return
	}

	// A failure already reported by the health check must not be overwritten by
	// a late success from the blocked generation call.
	if job.Status == statusFailed && err == nil {
		return
	}

	if err != nil {
		job.Status = statusFailed
		job.Error = err.Error()
		return
	}

	job.Status = statusCompleted
	job.Progress = 1.0
	job.Result = result

	if m.onComplete != nil && result != nil {
		go m.onComplete(job.ID, result)
	}
}
