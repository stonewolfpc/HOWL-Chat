package resource

// GPUInfo holds real-time GPU metrics for a single adapter.
type GPUInfo struct {
	Index          int    `json:"index"`
	Name           string `json:"name"`
	Vendor         string `json:"vendor"`
	TotalVRAMBytes int64  `json:"total_vram_bytes"`
	UsedVRAMBytes  int64  `json:"used_vram_bytes"`
	FreeVRAMBytes  int64  `json:"free_vram_bytes"`
	// Source names the detector that produced these values. Live is true only
	// when UsedVRAMBytes came from a live GPU counter, not from Howl's own model
	// registry. Admission refuses GPU starts when Live is false.
	Source string `json:"source,omitempty"`
	Live   bool   `json:"live"`
}

// gpuDetector is the platform-specific implementation.
var gpuDetector func() ([]GPUInfo, error)

// queryGPUs returns live GPU metrics via the platform detector.
func queryGPUs() ([]GPUInfo, error) {
	if gpuDetector == nil {
		return nil, nil
	}
	return gpuDetector()
}

// sumGPUVRAM totals all detected GPU memory.
func sumGPUVRAM(gpus []GPUInfo) (total, used int64) {
	for _, g := range gpus {
		total += g.TotalVRAMBytes
		used += g.UsedVRAMBytes
	}
	return
}
