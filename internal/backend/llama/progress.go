/**
 * Llama Progress
 *
 * This package defines progress tracking for model loading and inference.
 * Provides structured progress reporting with stages and percentages.
 *
 * @package llama
 */

package llama

import (
	"sync"
	"time"
)

// ProgressReporter handles progress reporting for model operations
type ProgressReporter struct {
	mu                sync.RWMutex
	stage             LoadingStage
	progress          float64
	startTime         time.Time
	currentStageStart time.Time
	listeners         []ProgressListener
}

// ProgressListener is a callback that receives progress updates
type ProgressListener func(progress ProgressUpdate)

// ProgressUpdate represents a single progress update
type ProgressUpdate struct {
	Stage         LoadingStage  `json:"stage"`          // Current loading stage
	Progress      float64       `json:"progress"`       // Progress percentage (0.0-1.0)
	StageProgress float64       `json:"stage_progress"` // Progress within current stage
	ElapsedTime   time.Duration `json:"elapsed_time"`   // Total elapsed time
	StageTime     time.Duration `json:"stage_time"`     // Time spent in current stage
	Message       string        `json:"message"`        // Human-readable message
}

// NewProgressReporter creates a new progress reporter
func NewProgressReporter() *ProgressReporter {
	now := time.Now()
	return &ProgressReporter{
		stage:             StageModelLoadStart,
		progress:          0.0,
		startTime:         now,
		currentStageStart: now,
		listeners:         []ProgressListener{},
	}
}

// UpdateProgress updates the overall progress
func (p *ProgressReporter) UpdateProgress(progress float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.progress = progress
	p.notifyListeners()
}

// UpdateStage updates the current loading stage
func (p *ProgressReporter) UpdateStage(stage LoadingStage) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.stage = stage
	p.currentStageStart = time.Now()
	p.notifyListeners()
}

// UpdateStageProgress updates progress within the current stage
func (p *ProgressReporter) UpdateStageProgress(stageProgress float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.notifyListeners()
}

// AddListener adds a progress listener
func (p *ProgressReporter) AddListener(listener ProgressListener) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.listeners = append(p.listeners, listener)
}

// RemoveListener removes a progress listener
func (p *ProgressReporter) RemoveListener(listener ProgressListener) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, l := range p.listeners {
		if &l == &listener {
			p.listeners = append(p.listeners[:i], p.listeners[i+1:]...)
			break
		}
	}
}

// GetCurrentProgress returns the current progress state
func (p *ProgressReporter) GetCurrentProgress() ProgressUpdate {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return ProgressUpdate{
		Stage:       p.stage,
		Progress:    p.progress,
		ElapsedTime: time.Since(p.startTime),
		StageTime:   time.Since(p.currentStageStart),
		Message:     p.getStageMessage(),
	}
}

// Reset resets the progress reporter to initial state
func (p *ProgressReporter) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	p.stage = StageModelLoadStart
	p.progress = 0.0
	p.startTime = now
	p.currentStageStart = now
}

// notifyListeners notifies all registered listeners of progress updates
func (p *ProgressReporter) notifyListeners() {
	update := ProgressUpdate{
		Stage:       p.stage,
		Progress:    p.progress,
		ElapsedTime: time.Since(p.startTime),
		StageTime:   time.Since(p.currentStageStart),
		Message:     p.getStageMessage(),
	}

	for _, listener := range p.listeners {
		listener(update)
	}
}

// getStageMessage returns a human-readable message for the current stage
func (p *ProgressReporter) getStageMessage() string {
	switch p.stage {
	case StageModelLoadStart:
		return "Starting model load"
	case StageTokenizerLoad:
		return "Loading tokenizer"
	case StageTensorAllocation:
		return "Allocating tensors"
	case StageKVCacheInit:
		return "Initializing KV cache"
	case StageModelReady:
		return "Model ready"
	default:
		return "Unknown stage"
	}
}

// ProgressTracker tracks progress across multiple stages
type ProgressTracker struct {
	mu            sync.RWMutex
	stageWeights  map[LoadingStage]float64
	currentStage  LoadingStage
	stageProgress float64
	totalProgress float64
}

// NewProgressTracker creates a new progress tracker with default stage weights
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		stageWeights: map[LoadingStage]float64{
			StageModelLoadStart:   0.05,
			StageTokenizerLoad:    0.10,
			StageTensorAllocation: 0.60,
			StageKVCacheInit:      0.20,
			StageModelReady:       0.05,
		},
		currentStage:  StageModelLoadStart,
		stageProgress: 0.0,
		totalProgress: 0.0,
	}
}

// SetStageWeights sets custom weights for each stage
func (t *ProgressTracker) SetStageWeights(weights map[LoadingStage]float64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.stageWeights = weights
}

// UpdateStageProgress updates progress within the current stage
func (t *ProgressTracker) UpdateStageProgress(progress float64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.stageProgress = progress
	t.calculateTotalProgress()
}

// SetCurrentStage sets the current stage
func (t *ProgressTracker) SetCurrentStage(stage LoadingStage) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.currentStage = stage
	t.stageProgress = 0.0
	t.calculateTotalProgress()
}

// GetTotalProgress returns the overall progress (0.0-1.0)
func (t *ProgressTracker) GetTotalProgress() float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.totalProgress
}

// GetCurrentStage returns the current loading stage
func (t *ProgressTracker) GetCurrentStage() LoadingStage {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.currentStage
}

// calculateTotalProgress calculates total progress based on stage weights
func (t *ProgressTracker) calculateTotalProgress() {
	var completedWeight float64
	var currentStageWeight float64

	for stage, weight := range t.stageWeights {
		if t.isStageComplete(stage) {
			completedWeight += weight
		} else if stage == t.currentStage {
			currentStageWeight = weight
			break
		}
	}

	stageContribution := currentStageWeight * t.stageProgress
	t.totalProgress = completedWeight + stageContribution
}

// isStageComplete checks if a stage is before the current stage
func (t *ProgressTracker) isStageComplete(stage LoadingStage) bool {
	stageOrder := []LoadingStage{
		StageModelLoadStart,
		StageTokenizerLoad,
		StageTensorAllocation,
		StageKVCacheInit,
		StageModelReady,
	}

	for _, s := range stageOrder {
		if s == stage {
			return true
		}
		if s == t.currentStage {
			return false
		}
	}

	return false
}

// Reset resets the progress tracker
func (t *ProgressTracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.currentStage = StageModelLoadStart
	t.stageProgress = 0.0
	t.totalProgress = 0.0
}
