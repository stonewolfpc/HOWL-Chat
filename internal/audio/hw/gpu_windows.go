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

func detectVRAMViaRegistry() int64 {
	// Try NVIDIA hardware info registry
	classes := []string{
		`{4d36e968-e325-11ce-bfc1-08002be10318}`,
		`{4d36e967-e325-11ce-bfc1-08002be10318}`,
	}
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
			// Try HardwareInformation.qwMemorySize (hex string, in bytes)
			memStr, _, err := sub.GetStringValue("HardwareInformation.qwMemorySize")
			if err == nil && memStr != "" {
				memStr = strings.TrimSpace(memStr)
				if strings.HasPrefix(memStr, "0x") {
					memStr = memStr[2:]
				}
				if n, err := strconv.ParseInt(memStr, 16, 64); err == nil && n > 0 {
					sub.Close()
					fmt.Printf("[hw] detected VRAM via registry: %d MB\n", n/1024/1024)
					return n
				}
			}
			// Try AdapterString or DriverDesc to identify NVIDIA/AMD
			driverDesc, _, _ := sub.GetStringValue("DriverDesc")
			sub.Close()
			if driverDesc != "" {
				// If we found a real GPU but no memory value, keep looking
				continue
			}
		}
		k.Close()
	}
	return 0
}

func detectVRAMViaWMI() int64 {
	ch := make(chan int64, 1)
	go func() {
		out, err := exec.Command("wmic", "path", "Win32_VideoController", "get", "AdapterRAM").CombinedOutput()
		if err != nil {
			ch <- 0
			return
		}
		var sum int64
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.EqualFold(line, "AdapterRAM") {
				continue
			}
			n, err := strconv.ParseInt(line, 10, 64)
			if err != nil || n <= 0 {
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
