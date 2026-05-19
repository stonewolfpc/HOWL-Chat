package job

import (
	"fmt"
	"sync"
	"time"

	"howl-chat/internal/image/engine"
	"howl-chat/internal/image/types"
)

// Manager tracks async image generation jobs.
type Manager struct {
	mu         sync.RWMutex
	jobs       map[string]*types.ImageJob
	engine     engine.ImageEngine
	onComplete func(id string, result *types.ImageResult)
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

// Start queues a generation job and returns its ID.
func (m *Manager) Start(opts types.ImageOptions) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	job := &types.ImageJob{
		ID:        fmt.Sprintf("job_%d", time.Now().UnixNano()),
		Status:    "queued",
		Progress:  0,
		Options:   opts,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	m.jobs[job.ID] = job

	go m.run(job)
	return job.ID
}

// Status returns the current status of a job.
func (m *Manager) Status(id string) *types.ImageJob {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.jobs[id]
}

// Cancel marks a job as cancelled. It cannot stop an in-flight HTTP request
// without killing the sd-server process, so cancellation only prevents the
// result from being committed and updates the status flag.
func (m *Manager) Cancel(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, ok := m.jobs[id]
	if !ok {
		return false
	}
	if job.Status == "completed" || job.Status == "failed" {
		return false
	}
	job.Status = "cancelled"
	return true
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

// run executes the generation in a goroutine.
func (m *Manager) run(job *types.ImageJob) {
	if m.engine == nil {
		m.finish(job, nil, fmt.Errorf("image engine not initialized"))
		return
	}

	m.mu.Lock()
	job.Status = "generating"
	m.mu.Unlock()

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
				if failed {
					continue
				}
				p, err := m.engine.GetProgress()
				if err != nil {
					failed = true
					continue
				}
				if p.Progress > 0 || p.Step > 0 {
					m.mu.Lock()
					if job.Status == "generating" {
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

	var result *types.ImageResult
	var err error

	switch job.Options.Mode {
	case "img2img":
		result, err = m.engine.GenerateImg2Img(job.Options)
	case "inpaint":
		result, err = m.engine.GenerateInpaint(job.Options)
	case "video":
		result, err = m.engine.GenerateVideo(job.Options)
	default:
		result, err = m.engine.GenerateTxt2Img(job.Options)
	}

	close(progressDone)
	m.finish(job, result, err)
}

// finish updates job state and cleans up if cancelled.
func (m *Manager) finish(job *types.ImageJob, result *types.ImageResult, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if job.Status == "cancelled" {
		return
	}

	if err != nil {
		job.Status = "failed"
		job.Error = err.Error()
		return
	}

	job.Status = "completed"
	job.Progress = 1.0
	job.Result = result

	if m.onComplete != nil && result != nil {
		go m.onComplete(job.ID, result)
	}
}
