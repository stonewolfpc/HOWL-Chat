//go:build !windows

package hw

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func detectGPUPresent() bool {
	if _, err := os.Stat("/dev/dri"); err == nil {
		return true
	}
	if path, err := exec.LookPath("nvidia-smi"); err == nil && path != "" {
		return true
	}
	if path, err := exec.LookPath("rocm-smi"); err == nil && path != "" {
		return true
	}
	if runtime.GOOS == "darwin" {
		return detectGPUmacOS()
	}
	return detectGPULinuxPCI()
}

func detectGPULinuxPCI() bool {
	matches, _ := filepath.Glob("/sys/bus/pci/devices/*/class")
	for _, classFile := range matches {
		data, err := os.ReadFile(classFile)
		if err != nil {
			continue
		}
		s := strings.TrimSpace(string(data))
		if strings.HasPrefix(s, "0x03") {
			return true
		}
	}
	return false
}

func detectGPUmacOS() bool {
	if runtime.GOOS != "darwin" {
		return false
	}
	out, err := exec.Command("system_profiler", "SPDisplaysDataType").CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "Chipset Model:")
}

func detectVRAMBytes() int64 {
	if v := nvidiaVRAM(); v > 0 {
		return v
	}
	if v := amdgpuSysfsVRAM(); v > 0 {
		return v
	}
	if v := darwinVRAMFromProfiler(); v > 0 {
		return v
	}
	return 0
}

func nvidiaVRAM() int64 {
	out, err := exec.Command("nvidia-smi", "--query-gpu=memory.total", "--format=csv,noheader,nounits").CombinedOutput()
	if err != nil {
		return 0
	}
	var sum int64
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			continue
		}
		sum += n * 1024 * 1024
	}
	return sum
}

func amdgpuSysfsVRAM() int64 {
	matches, _ := filepath.Glob("/sys/class/drm/card*/device/mem_info_vram_total")
	var max int64
	for _, p := range matches {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		n, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
		if err != nil {
			continue
		}
		if n > max {
			max = n
		}
	}
	return max
}

var vramMiBRe = regexp.MustCompile(`(\d+)\s*(MB|GB)\s*\(`)

func darwinVRAMFromProfiler() int64 {
	if runtime.GOOS != "darwin" {
		return 0
	}
	out, err := exec.Command("system_profiler", "SPDisplaysDataType").CombinedOutput()
	if err != nil {
		return 0
	}
	text := string(out)
	if i := strings.Index(text, "VRAM"); i >= 0 {
		rest := text[i:]
		sub := rest[:minInt(120, len(rest))]
		ms := vramMiBRe.FindStringSubmatch(sub)
		if len(ms) == 3 {
			n, err := strconv.ParseInt(ms[1], 10, 64)
			if err != nil {
				return 0
			}
			if strings.EqualFold(ms[2], "GB") {
				return n * 1024 * 1024 * 1024
			}
			return n * 1024 * 1024
		}
	}
	return unifiedMemoryDarwin(text)
}

func unifiedMemoryDarwin(text string) int64 {
	idx := strings.Index(text, "Memory:")
	if idx < 0 {
		return 0
	}
	rest := strings.TrimSpace(text[idx+len("Memory:") :])
	end := strings.IndexAny(rest, "\n")
	if end > 0 {
		rest = rest[:end]
	}
	fields := strings.Fields(rest)
	if len(fields) == 0 {
		return 0
	}
	n, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0
	}
	switch {
	case strings.Contains(rest, "GB"):
		return n * 1024 * 1024 * 1024
	case strings.Contains(rest, "MB"):
		return n * 1024 * 1024
	default:
		return 0
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
