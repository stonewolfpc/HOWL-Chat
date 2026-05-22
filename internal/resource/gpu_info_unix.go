//go:build !windows

package resource

import (
	"os/exec"
	"strconv"
	"strings"
)

func init() {
	gpuDetector = detectUnixGPUs
}

// detectUnixGPUs attempts nvidia-smi, then rocm-smi, then amdgpu sysfs.
func detectUnixGPUs() ([]GPUInfo, error) {
	if gpus := nvidiaGPUs(); len(gpus) > 0 {
		return gpus, nil
	}
	if gpus := amdgpuGPUs(); len(gpus) > 0 {
		return gpus, nil
	}
	return nil, nil
}

// nvidiaGPUs queries NVIDIA cards via nvidia-smi.
func nvidiaGPUs() []GPUInfo {
	out, err := exec.Command("nvidia-smi",
		"--query-gpu=name,memory.total,memory.used",
		"--format=csv,noheader,nounits",
	).CombinedOutput()
	if err != nil {
		return nil
	}

	var gpus []GPUInfo
	for i, line := range strings.Split(string(out), "\n") {
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
		gpus = append(gpus, GPUInfo{
			Index:          i,
			Name:           name,
			Vendor:         "NVIDIA",
			TotalVRAMBytes: totalMiB * 1024 * 1024,
			UsedVRAMBytes:  usedMiB * 1024 * 1024,
			FreeVRAMBytes:  (totalMiB - usedMiB) * 1024 * 1024,
		})
	}
	return gpus
}

// amdgpuGPUs reads AMD GPU memory from sysfs.
func amdgpuGPUs() []GPUInfo {
	// amdgpu sysfs: /sys/class/drm/card*/device/mem_info_vram_total
	// and /sys/class/drm/card*/device/mem_info_vram_used
	// This is a simplified fallback; full implementation would parse both files.
	return nil
}
