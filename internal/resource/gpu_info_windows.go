//go:build windows

package resource

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/yusufpapurcu/wmi"
	"golang.org/x/sys/windows/registry"
)

// Win32_VideoController provides static GPU information via WMI.
type win32VideoController struct {
	Name           string `json:"Name"`
	AdapterRAM     uint64 `json:"AdapterRAM"`
	VideoProcessor string `json:"VideoProcessor"`
}

// Win32_PerfFormattedData_GPUPerformanceCounters_GPUAdapterMemory provides
// live GPU memory counters via WMI.
type gpuAdapterMemory struct {
	Name           string `json:"Name"`
	DedicatedUsage uint64 `json:"DedicatedUsage"`
	SharedUsage    uint64 `json:"SharedUsage"`
}

func init() {
	gpuDetector = detectWindowsGPUs
}

// detectWindowsGPUs queries GPU info using the fastest available method.
// Priority: nvidia-smi -> registry -> WMI (with timeout).
func detectWindowsGPUs() ([]GPUInfo, error) {
	// 1. Try nvidia-smi first (fast, accurate for NVIDIA)
	if gpus := detectViaNvidiaSMI(); len(gpus) > 0 {
		return gpus, nil
	}
	// 2. Try registry (no external process)
	if gpus := detectViaRegistry(); len(gpus) > 0 {
		return gpus, nil
	}
	// 3. WMI with timeout as last resort
	return detectViaWMI()
}

func detectViaNvidiaSMI() []GPUInfo {
	ch := make(chan []GPUInfo, 1)
	go func() {
		// Query name, total memory, used memory
		out, err := exec.Command("nvidia-smi",
			"--query-gpu=name,memory.total,memory.used",
			"--format=csv,noheader,nounits").CombinedOutput()
		if err != nil {
			ch <- nil
			return
		}
		var gpus []GPUInfo
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.Split(line, ", ")
			if len(parts) < 3 {
				continue
			}
			name := strings.TrimSpace(parts[0])
			totalMiB, _ := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
			usedMiB, _ := strconv.ParseInt(strings.TrimSpace(parts[2]), 10, 64)
			if totalMiB <= 0 {
				continue
			}
			totalBytes := totalMiB * 1024 * 1024
			usedBytes := usedMiB * 1024 * 1024
			gpus = append(gpus, GPUInfo{
				Index:          len(gpus),
				Name:           name,
				Vendor:         inferVendor(name, ""),
				TotalVRAMBytes: totalBytes,
				UsedVRAMBytes:  usedBytes,
				FreeVRAMBytes:  maxInt64(0, totalBytes-usedBytes),
				Source:         "nvidia-smi",
				Live:           true,
			})
		}
		ch <- gpus
	}()
	select {
	case gpus := <-ch:
		if len(gpus) > 0 {
			fmt.Printf("[resource] GPU detected via nvidia-smi: %s (%d MB)\n", gpus[0].Name, gpus[0].TotalVRAMBytes/1024/1024)
		}
		return gpus
	case <-time.After(3 * time.Second):
		fmt.Println("[resource] nvidia-smi timeout")
		return nil
	}
}

// minUsableVRAMBytes filters out integrated GPUs and virtual adapters.
const minUsableVRAMBytes = 512 * 1024 * 1024

func detectViaRegistry() []GPUInfo {
	classes := []string{
		`{4d36e968-e325-11ce-bfc1-08002be10318}`,
		`{4d36e967-e325-11ce-bfc1-08002be10318}`,
	}
	var gpus []GPUInfo
	var foundDiscrete bool
	for _, classGUID := range classes {
		keyPath := `SYSTEM\CurrentControlSet\Control\Class\` + classGUID
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.READ)
		if err != nil {
			continue
		}
		names, _ := k.ReadSubKeyNames(-1)
		for _, name := range names {
			sub, err := registry.OpenKey(k, name, registry.READ)
			if err != nil {
				continue
			}
			driverDesc, _, _ := sub.GetStringValue("DriverDesc")
			memStr, _, _ := sub.GetStringValue("HardwareInformation.qwMemorySize")
			sub.Close()
			if driverDesc == "" {
				continue
			}
			lowerDesc := strings.ToLower(driverDesc)
			isIntelIGP := strings.Contains(lowerDesc, "intel") &&
				(strings.Contains(lowerDesc, "uhd") ||
					strings.Contains(lowerDesc, "iris") ||
					strings.Contains(lowerDesc, "hd graphics"))

			var totalBytes int64
			if memStr != "" {
				memStr = strings.TrimPrefix(strings.TrimSpace(memStr), "0x")
				n, err := strconv.ParseInt(memStr, 16, 64)
				if err == nil && n > 0 {
					totalBytes = n
				}
			}
			// If registry lacks memory size but we know the GPU series, estimate it
			if totalBytes == 0 && !isIntelIGP {
				totalBytes = estimateVRAMFromGPUName(driverDesc)
			}
			if totalBytes < minUsableVRAMBytes {
				continue
			}
			// If we already found a discrete GPU, skip Intel integrated
			if foundDiscrete && isIntelIGP {
				continue
			}
			used := estimateUsedVRAM()
			gpus = append(gpus, GPUInfo{
				Index:          len(gpus),
				Name:           driverDesc,
				Vendor:         inferVendor(driverDesc, ""),
				TotalVRAMBytes: totalBytes,
				UsedVRAMBytes:  used,
				FreeVRAMBytes:  maxInt64(0, totalBytes-used),
				Source:         "registry",
				Live:           false,
			})
			if !isIntelIGP {
				foundDiscrete = true
			}
		}
		k.Close()
	}
	if len(gpus) > 0 {
		fmt.Printf("[resource] GPU detected via registry: %s (%d MB)\n", gpus[0].Name, gpus[0].TotalVRAMBytes/1024/1024)
	}
	return gpus
}

// estimateVRAMFromGPUName guesses VRAM from GPU model name patterns.
func estimateVRAMFromGPUName(name string) int64 {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "rtx 4090"):
		return 24 * 1024 * 1024 * 1024
	case strings.Contains(lower, "rtx 4080"), strings.Contains(lower, "rtx 3090"):
		return 16 * 1024 * 1024 * 1024
	case strings.Contains(lower, "rtx 4070"), strings.Contains(lower, "rtx 3080"):
		return 12 * 1024 * 1024 * 1024
	case strings.Contains(lower, "rtx 4060"), strings.Contains(lower, "rtx 3070"):
		return 8 * 1024 * 1024 * 1024
	case strings.Contains(lower, "rtx 3060"), strings.Contains(lower, "gtx 1080"):
		return 8 * 1024 * 1024 * 1024
	case strings.Contains(lower, "gtx 1070"), strings.Contains(lower, "gtx 1660"):
		return 8 * 1024 * 1024 * 1024
	case strings.Contains(lower, "gtx 1060"):
		return 6 * 1024 * 1024 * 1024
	case strings.Contains(lower, "rtx 50"), strings.Contains(lower, "rtx 40"),
		strings.Contains(lower, "rtx 30"):
		return 8 * 1024 * 1024 * 1024
	case strings.Contains(lower, "radeon rx 79"), strings.Contains(lower, "radeon rx 78"):
		return 24 * 1024 * 1024 * 1024
	case strings.Contains(lower, "radeon rx 69"), strings.Contains(lower, "radeon rx 68"):
		return 16 * 1024 * 1024 * 1024
	case strings.Contains(lower, "radeon rx 67"), strings.Contains(lower, "radeon rx 66"):
		return 12 * 1024 * 1024 * 1024
	case strings.Contains(lower, "radeon rx 6650"), strings.Contains(lower, "radeon rx 6600"):
		return 8 * 1024 * 1024 * 1024
	case strings.Contains(lower, "radeon rx 5700"), strings.Contains(lower, "radeon rx 5600"):
		return 8 * 1024 * 1024 * 1024
	default:
		return 0
	}
}

func detectViaWMI() ([]GPUInfo, error) {
	ch := make(chan []GPUInfo, 1)
	go func() {
		var controllers []win32VideoController
		q := wmi.Query("SELECT Name, AdapterRAM, VideoProcessor FROM Win32_VideoController", &controllers, nil, nil)
		if q != nil {
			ch <- nil
			return
		}
		var gpus []GPUInfo
		for i, c := range controllers {
			name := strings.TrimSpace(c.Name)
			if isDisplayAdapter(name) {
				continue
			}
			used := estimateUsedVRAM()
			gpus = append(gpus, GPUInfo{
				Index:          i,
				Name:           name,
				Vendor:         inferVendor(name, c.VideoProcessor),
				TotalVRAMBytes: int64(c.AdapterRAM),
				UsedVRAMBytes:  used,
				FreeVRAMBytes:  maxInt64(0, int64(c.AdapterRAM)-used),
				Source:         "wmi",
				Live:           false,
			})
		}
		ch <- gpus
	}()
	select {
	case gpus := <-ch:
		if len(gpus) > 0 {
			fmt.Printf("[resource] GPU detected via WMI: %s (%d MB)\n", gpus[0].Name, gpus[0].TotalVRAMBytes/1024/1024)
		}
		return gpus, nil
	case <-time.After(3 * time.Second):
		fmt.Println("[resource] WMI timeout")
		return nil, fmt.Errorf("WMI GPU detection timed out after 3s")
	}
}

// isDisplayAdapter filters out non-physical GPU entries from WMI.
func isDisplayAdapter(name string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, "rdp") ||
		strings.Contains(lower, "remote") ||
		strings.Contains(lower, "basic display") ||
		strings.Contains(lower, "microsoft basic")
}

// inferVendor guesses GPU vendor from name or video processor string.
func inferVendor(name, videoProcessor string) string {
	lower := strings.ToLower(name + " " + videoProcessor)
	switch {
	case strings.Contains(lower, "nvidia") || strings.Contains(lower, "geforce") || strings.Contains(lower, "rtx"):
		return "NVIDIA"
	case strings.Contains(lower, "amd") || strings.Contains(lower, "radeon"):
		return "AMD"
	case strings.Contains(lower, "intel") || strings.Contains(lower, "arc") || strings.Contains(lower, "iris"):
		return "Intel"
	case strings.Contains(lower, "apple"):
		return "Apple"
	default:
		return "Unknown"
	}
}

// estimateUsedVRAM sums GPU memory from models registered with the manager.
// Called as a fallback when live counters are unavailable.
func estimateUsedVRAM() int64 {
	m := GetManager()
	m.mu.RLock()
	defer m.mu.RUnlock()
	var total int64
	for _, model := range m.models {
		if model.Loaded && model.IsGPUOffload {
			total += model.GPUMemoryMB * 1024 * 1024
		}
	}
	return total
}
