//go:build windows

package hw

import (
	"os/exec"
	"strconv"
	"strings"

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
	out, err := exec.Command("wmic", "path", "Win32_VideoController", "get", "AdapterRAM").CombinedOutput()
	if err != nil {
		return 0
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
	return sum
}
