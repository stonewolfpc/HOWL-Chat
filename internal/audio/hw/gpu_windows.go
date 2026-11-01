//go:build windows

package hw

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

func detectGPUPresent() bool {
	if detectGPURegistryClass(`{4d36e968-e325-11ce-bfc1-08002be10318}`) {
		return true
	}
	return detectGPURegistryClass(`{4d36e967-e325-11ce-bfc1-08002be10318}`)
}

func detectGPURegistryClass(classGUID string) bool {
	keyPath := `SYSTEM\CurrentControlSet\Control\Class\` + classGUID
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.READ)
	if err != nil {
		return false
	}
	defer k.Close()
	names, err := k.ReadSubKeyNames(-1)
	if err != nil || len(names) == 0 {
		return false
	}
	for _, name := range names {
		sub, err := registry.OpenKey(k, name, registry.READ)
		if err != nil {
			continue
		}
		driverDesc, _, _ := sub.GetStringValue("DriverDesc")
		sub.Close()
		if driverDesc != "" {
			return true
		}
	}
	return false
}

func detectVRAMBytes() int64 {
	// 1. Try nvidia-smi first (fast, reliable for NVIDIA)
	if v := detectVRAMViaNvidiaSMI(); v > 0 {
		return v
	}
	// 2. Try registry read (no external process)
	if v := detectVRAMViaRegistry(); v > 0 {
		return v
	}
	// 3. Last resort: wmic with a 3s timeout
	if v := detectVRAMViaWMI(); v > 0 {
		return v
	}
	return 0
}

func detectVRAMViaNvidiaSMI() int64 {
	ch := make(chan int64, 1)
	go func() {
		out, err := exec.Command("nvidia-smi", "--query-gpu=memory.total", "--format=csv,noheader,nounits").CombinedOutput()
		if err != nil {
			ch <- 0
			return
		}
		var sum int64
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// nvidia-smi returns MiB
			n, err := strconv.ParseInt(line, 10, 64)
			if err != nil || n <= 0 {
				continue
			}
			sum += n * 1024 * 1024 // convert MiB to bytes
		}
		ch <- sum
	}()
	select {
	case v := <-ch:
		if v > 0 {
			fmt.Printf("[hw] detected VRAM via nvidia-smi: %d MB\n", v/1024/1024)
		}
		return v
	case <-time.After(3 * time.Second):
		fmt.Println("[hw] nvidia-smi timeout")
		return 0
	}
}

// minDiscreteVRAMBytes filters out integrated GPUs and virtual adapters.
// 512 MB is a reasonable floor for any usable discrete GPU in 2025.
const minDiscreteVRAMBytes = 512 * 1024 * 1024

func detectVRAMViaRegistry() int64 {
	classes := []string{
		`{4d36e968-e325-11ce-bfc1-08002be10318}`,
		`{4d36e967-e325-11ce-bfc1-08002be10318}`,
	}
	var totalBytes int64
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
			// Skip entries without a driver description (not real GPUs)
			if driverDesc == "" {
				sub.Close()
				continue
			}
			lowerDesc := strings.ToLower(driverDesc)
			isIntelIGP := strings.Contains(lowerDesc, "intel") &&
				(strings.Contains(lowerDesc, "uhd") ||
					strings.Contains(lowerDesc, "iris") ||
					strings.Contains(lowerDesc, "hd graphics"))
			// Try HardwareInformation.qwMemorySize (hex string, in bytes)
			memStr, _, err := sub.GetStringValue("HardwareInformation.qwMemorySize")
			sub.Close()
			var gpuBytes int64
			if err == nil && memStr != "" {
				memStr = strings.TrimSpace(memStr)
				if strings.HasPrefix(memStr, "0x") {
					memStr = memStr[2:]
				}
				if n, err := strconv.ParseInt(memStr, 16, 64); err == nil && n > 0 {
					gpuBytes = n
				}
			}
			// If no memory value but we have a real NVIDIA/AMD GPU, estimate
			// from known series patterns so we don't skip it entirely.
			if gpuBytes == 0 && !isIntelIGP {
				gpuBytes = estimateVRAMFromName(driverDesc)
			}
			if gpuBytes < minDiscreteVRAMBytes {
				continue
			}
			// If we already found a discrete GPU, skip Intel integrated graphics
			if foundDiscrete && isIntelIGP {
				continue
			}
			if gpuBytes > 0 {
				if !isIntelIGP {
					foundDiscrete = true
				}
				totalBytes += gpuBytes
			}
		}
		k.Close()
	}
	if totalBytes > 0 {
		fmt.Printf("[hw] detected VRAM via registry: %d MB (discrete=%v)\n",
			totalBytes/1024/1024, foundDiscrete)
	}
	return totalBytes
}

// estimateVRAMFromName guesses VRAM from GPU model name patterns.
func estimateVRAMFromName(name string) int64 {
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

func detectVRAMViaWMI() int64 {
	ch := make(chan int64, 1)
	go func() {
		out, err := exec.Command("wmic", "path", "Win32_VideoController", "get", "Name,AdapterRAM").CombinedOutput()
		if err != nil {
			ch <- 0
			return
		}
		var sum int64
		lines := strings.Split(string(out), "\n")
		for i := 0; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line == "" || strings.Contains(line, "AdapterRAM") {
				continue
			}
			// WMIC output can be tricky; look for numeric values
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			name := strings.ToLower(strings.Join(fields[:len(fields)-1], " "))
			lastField := fields[len(fields)-1]
			n, err := strconv.ParseInt(lastField, 10, 64)
			if err != nil || n <= 0 {
				continue
			}
			// Skip display adapters, RDP, and tiny GPUs
			if strings.Contains(name, "rdp") || strings.Contains(name, "remote") ||
				strings.Contains(name, "basic display") || strings.Contains(name, "microsoft basic") {
				continue
			}
			if n < minDiscreteVRAMBytes {
				continue
			}
			sum += n
		}
		ch <- sum
	}()
	select {
	case v := <-ch:
		if v > 0 {
			fmt.Printf("[hw] detected VRAM via WMI: %d MB\n", v/1024/1024)
		}
		return v
	case <-time.After(3 * time.Second):
		fmt.Println("[hw] WMI timeout")
		return 0
	}
}
