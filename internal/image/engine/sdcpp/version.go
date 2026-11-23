package sdcpp

import (
	"os/exec"
	"regexp"
	"sync"

	"howl-chat/internal/proc"
)

// Version probing belongs to the engine package: it is the package that knows
// sd-server's binary contract, and both the image domain's boot and the
// diagnostics binding must ask the same authority instead of running the binary
// themselves.
var (
	versionMu    sync.Mutex
	versionCache = map[string]string{}
	versionRe    = regexp.MustCompile(`commit\s+([0-9a-f]+)`)
)

// BinaryVersion runs sd-server once and reports the commit hash it prints. The
// result is cached per binary path: the diagnostics page polls for it, and each
// call used to spawn a process.
func BinaryVersion(binary string) string {
	if binary == "" {
		return "not found"
	}

	versionMu.Lock()
	cached, ok := versionCache[binary]
	versionMu.Unlock()
	if ok {
		return cached
	}

	version := "unknown"
	cmd := exec.Command(binary, "--help")
	proc.HideWindow(cmd)
	if out, err := cmd.CombinedOutput(); err == nil {
		if m := versionRe.FindStringSubmatch(string(out)); len(m) > 1 {
			version = m[1]
		}
	}

	if version != "unknown" {
		versionMu.Lock()
		versionCache[binary] = version
		versionMu.Unlock()
	}
	return version
}

// IsKnownOldVersion reports commits that cannot load modern GGUF checkpoints.
// SD3.5 support landed in stable-diffusion.cpp around May 2026.
func IsKnownOldVersion(commit string) bool {
	switch commit {
	case "0665a7f", "3c2d18b", "5f3e8a1":
		return true
	}
	return false
}
