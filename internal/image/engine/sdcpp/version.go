package sdcpp

import (
	"os/exec"
	"regexp"
	"sync"

	"howl-chat/internal/proc"
)

// EngineTag and EngineCommit pin the stable-diffusion.cpp release HOWL ships.
//
// The engine cuts several builds a day, so the binary in the tree is matched to
// an exact tag rather than to "latest": a swap is a deliberate act with a
// recorded target, not whatever a download happened to return. When these are
// updated, capture the new help screen as a baseline before comparing (see
// build/engine-baseline/) and re-run the family matrix.
const (
	EngineTag    = "master-1122-3b12781"
	EngineCommit = "3b12781"
)

// Version probing belongs to the engine package: it is the package that knows
// sd-server's binary contract, and both the image domain's boot and the
// diagnostics binding must ask the same authority instead of running the binary
// themselves.
var (
	versionMu sync.Mutex
	versionRe = regexp.MustCompile(`commit\s+([0-9a-f]+)`)
)

// runHelp executes the binary's help screen. Callers reach it through
// helpOutput so every probe of one path costs one launch.
func runHelp(binary string) string {
	cmd := exec.Command(binary, "--help")
	proc.HideWindow(cmd)
	out, _ := cmd.CombinedOutput()
	return string(out)
}

// BinaryVersion reports the commit hash sd-server prints, cached per binary
// path. The diagnostics page polls for it, and each call used to spawn a
// process.
func BinaryVersion(binary string) string {
	if binary == "" {
		return "not found"
	}
	out := helpOutput(binary)
	if m := versionRe.FindStringSubmatch(out); len(m) > 1 {
		return m[1]
	}
	return "unknown"
}

// PinnedVersionMatches reports whether the binary answers with the commit HOWL
// pinned. A mismatch is not fatal — it is the signal that the capability tables
// and the binary may describe different engines.
func PinnedVersionMatches(binary string) bool {
	v := BinaryVersion(binary)
	return v == EngineCommit || v == "unknown"
}
