package resource

import (
	"fmt"
	"sync"
	"time"

	"howl-chat/internal/audio/hw"
)

// ---------------------------------------------------------------------------
// Singleton
// ---------------------------------------------------------------------------

var (
	instance     *ResourceManager
	instanceOnce sync.Once
)

// GetManager returns the singleton ResourceManager, creating it on first call.
func GetManager() *ResourceManager {
	fmt.Println("[ResourceManager] GetManager called")
	instanceOnce.Do(func() {
		fmt.Println("[ResourceManager] initializing singleton")
		instance = &ResourceManager{
			models: make(map[string]*ModelState),
		}
		_ = instance.refreshHardware()
		fmt.Println("[ResourceManager] initialization complete")
	})
	return instance
}

// ---------------------------------------------------------------------------
// ResourceManager
// ---------------------------------------------------------------------------

// ResourceManager tracks loaded models and hardware state.
// It is safe for concurrent use.
type ResourceManager struct {
	mu       sync.RWMutex
	models   map[string]*ModelState // key = "type:name"
	snapshot HardwareSnapshot
	lastScan time.Time
}

// ---------------------------------------------------------------------------
// Model registration
// ---------------------------------------------------------------------------

// RegisterModel records a newly loaded model.
func (rm *ResourceManager) RegisterModel(modelType ResourceType, name, path string, gpuMB, ramMB int64, priority int, isGPUOffload bool) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.models[key(modelType, name)] = &ModelState{
		Type:         modelType,
		Name:         name,
		Path:         path,
		Loaded:       true,
		GPUMemoryMB:  gpuMB,
		SystemRAMMB:  ramMB,
		Priority:     priority,
		LastUsed:     time.Now(),
		IsGPUOffload: isGPUOffload,
	}
}

// UnregisterModel removes a model from tracking (e.g. after unload).
func (rm *ResourceManager) UnregisterModel(modelType ResourceType, name string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	delete(rm.models, key(modelType, name))
}

// UpdateModelUsage bumps the LastUsed timestamp.
func (rm *ResourceManager) UpdateModelUsage(modelType ResourceType, name string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if m, ok := rm.models[key(modelType, name)]; ok {
		m.LastUsed = time.Now()
	}
}

// LoadedModels returns a snapshot of all loaded models.
func (rm *ResourceManager) LoadedModels() []ModelState {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	out := make([]ModelState, 0, len(rm.models))
	for _, m := range rm.models {
		if m.Loaded {
			out = append(out, *m)
		}
	}
	return out
}

// ModelsByType returns loaded models of a single resource type.
func (rm *ResourceManager) ModelsByType(rt ResourceType) []ModelState {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	var out []ModelState
	for _, m := range rm.models {
		if m.Type == rt && m.Loaded {
			out = append(out, *m)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Hardware
// ---------------------------------------------------------------------------

// HardwareSnapshot returns a point-in-time copy of the current hardware state.
func (rm *ResourceManager) HardwareSnapshot() HardwareSnapshot {
	fmt.Println("[ResourceManager] HardwareSnapshot: acquiring RLock")
	rm.mu.RLock()
	fmt.Println("[ResourceManager] HardwareSnapshot: RLock acquired")
	defer func() {
		rm.mu.RUnlock()
		fmt.Println("[ResourceManager] HardwareSnapshot: RLock released")
	}()
	return rm.snapshot
}

// RefreshHardware updates the stored hardware snapshot with live data.
func (rm *ResourceManager) RefreshHardware() {
	_, _ = rm.RefreshSnapshot()
}

// RefreshSnapshot refreshes and returns the snapshot together with a detection
// error. Admission uses this form so a failed read can fail closed instead of
// reusing a stale value.
func (rm *ResourceManager) RefreshSnapshot() (HardwareSnapshot, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	err := rm.refreshHardware()
	return rm.snapshot, err
}

// refreshHardware must be called with rm.mu held.
func (rm *ResourceManager) refreshHardware() error {
	fmt.Println("[ResourceManager] refreshHardware: calling hw.DetectHardware()")
	info := hw.DetectHardware()
	fmt.Println("[ResourceManager] refreshHardware: hw.DetectHardware() returned")

	// Preserve last known good GPU values so a single failed query
	// doesn't wipe out a correct detection.
	lastTotal := rm.snapshot.TotalVRAMBytes
	lastName := rm.snapshot.GPUModel
	lastVendor := rm.snapshot.GPUVendor

	gpuTotal, gpuUsed := int64(0), int64(0)
	gpuName, gpuVendor := "", ""
	gpus := []GPUInfo(nil)
	vramSource := "none"
	vramLive := false
	fmt.Println("[ResourceManager] refreshHardware: calling queryGPUs() with 3s timeout")
	gpuCh := make(chan []GPUInfo, 1)
	go func() {
		if gpus, err := queryGPUs(); err == nil && len(gpus) > 0 {
			gpuCh <- gpus
		} else {
			gpuCh <- nil
		}
	}()
	select {
	case detected := <-gpuCh:
		if detected != nil {
			gpus = detected
			fmt.Printf("[ResourceManager] refreshHardware: queryGPUs() returned %d GPUs\n", len(gpus))
			gpuTotal, gpuUsed = sumGPUVRAM(gpus)
			gpuName = gpus[0].Name
			gpuVendor = gpus[0].Vendor
			vramSource = gpus[0].Source
			if vramSource == "" {
				vramSource = "detector"
			}
			vramLive = gpus[0].Live
		} else {
			fmt.Println("[ResourceManager] refreshHardware: queryGPUs() returned empty/failed, using fallback")
			gpuTotal, gpuUsed = rm.fallbackGPUValues(info, lastTotal)
			gpuName = lastName
			gpuVendor = lastVendor
			vramSource = "fallback"
		}
	case <-time.After(3 * time.Second):
		fmt.Println("[ResourceManager] refreshHardware: queryGPUs() TIMED OUT after 3s, using fallback")
		gpuTotal, gpuUsed = rm.fallbackGPUValues(info, lastTotal)
		gpuName = lastName
		gpuVendor = lastVendor
		vramSource = "timeout"
	}

	// Sanity-check: if detection claims < 1 GB total on a system that
	// previously had a real GPU, keep the last known good total.
	if gpuTotal > 0 && gpuTotal < 1024*1024*1024 && lastTotal >= 1024*1024*1024 {
		fmt.Printf("[ResourceManager] refreshHardware: suspicious GPU total %d MB, restoring last known %d MB\n",
			gpuTotal/1024/1024, lastTotal/1024/1024)
		gpuTotal = lastTotal
		gpuName = lastName
		gpuVendor = lastVendor
		vramSource = "last-known"
		vramLive = false
	}
	fmt.Println("[ResourceManager] refreshHardware: storing snapshot")

	detectionErr := ""
	if info.Error != nil {
		detectionErr = info.Error.Error()
	} else if info.HasGPU && gpuTotal <= 0 {
		detectionErr = "a GPU was detected but its VRAM could not be read"
	}

	rm.snapshot = HardwareSnapshot{
		Timestamp:      time.Now(),
		GPUModel:       gpuName,
		GPUVendor:      gpuVendor,
		TotalVRAMBytes: gpuTotal,
		UsedVRAMBytes:  gpuUsed,
		FreeVRAMBytes:  maxInt64(0, gpuTotal-gpuUsed),
		TotalRAMBytes:  info.TotalRAMBytes,
		UsedRAMBytes:   info.TotalRAMBytes - info.AvailableRAMBytes,
		FreeRAMBytes:   info.AvailableRAMBytes,
		AvailableRAM:   info.AvailableRAMBytes,
		MemoryTier:     string(info.MemoryTier),
		OS:             info.OS,
		Arch:           info.Arch,
		GPUs:           gpus,
		VRAMLive:       vramLive,
		VRAMSource:     vramSource,
		GPUAvailable:   info.HasGPU || gpuTotal > 0,
		DetectionError: detectionErr,
	}
	rm.lastScan = time.Now()
	if detectionErr != "" {
		return fmt.Errorf("%s", detectionErr)
	}
	return nil
}

// fallbackGPUValues returns total and used VRAM when live queries fail.
// It prefers the last known good total; otherwise falls back to hw detection.
func (rm *ResourceManager) fallbackGPUValues(info *hw.HardwareInfo, lastTotal int64) (total, used int64) {
	if lastTotal > 0 {
		return lastTotal, rm.usedVRAMLocked()
	}
	if info.GPUMemoryBytes > 0 {
		return info.GPUMemoryBytes, rm.usedVRAMLocked()
	}
	return 0, 0
}

// usedVRAMLocked sums GPU memory from tracked models. Must hold rm.mu.
func (rm *ResourceManager) usedVRAMLocked() int64 {
	var total int64
	for _, m := range rm.models {
		if m.Loaded && m.IsGPUOffload {
			total += m.GPUMemoryMB * 1024 * 1024
		}
	}
	return total
}

// ---------------------------------------------------------------------------
// Resource summary
// ---------------------------------------------------------------------------

// Summary returns a concise, human-readable resource overview.
func (rm *ResourceManager) Summary() string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	var llm, tts, img int
	for _, m := range rm.models {
		if !m.Loaded {
			continue
		}
		switch m.Type {
		case ResourceLLM:
			llm++
		case ResourceTTS:
			tts++
		case ResourceImage:
			img++
		}
	}

	return fmt.Sprintf(
		"GPU: %s (%s) | VRAM %d/%d MB | RAM %d/%d MB | Models: %d LLM, %d TTS, %d Image",
		rm.snapshot.GPUModel,
		rm.snapshot.MemoryTier,
		rm.snapshot.UsedVRAMBytes/1024/1024,
		rm.snapshot.TotalVRAMBytes/1024/1024,
		rm.snapshot.UsedRAMBytes/1024/1024,
		rm.snapshot.TotalRAMBytes/1024/1024,
		llm, tts, img,
	)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func key(rt ResourceType, name string) string { return string(rt) + ":" + name }

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
