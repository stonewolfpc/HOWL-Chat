package llama

// capability.go — what the shipped llama-server binary actually advertises.
//
// Howl shells out to llama-server, so a flag is only available when the binary
// on this machine accepts it. Probing the binary's own --help is the same rule
// the image engine already applies to sd-server: no hardcoded feature list, no
// optimistic flag, and a stated reason when an option cannot be offered.

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

// CapabilitySet is the parsed answer from one llama-server binary.
type CapabilitySet struct {
	Binary  string          `json:"binary"`
	Version string          `json:"version"`
	Build   string          `json:"build"`
	Commit  string          `json:"commit"`
	Probed  bool            `json:"probed"`
	Error   string          `json:"error,omitempty"`
	Flags   map[string]bool `json:"flags"`
}

// Supports reports whether the probed help advertised flag. An unprobed binary
// supports nothing optional: sending a flag it may not understand is worse than
// omitting an optimisation.
func (c CapabilitySet) Supports(flag string) bool {
	if !c.Probed || c.Flags == nil {
		return false
	}
	return c.Flags[flag]
}

type cachedCapability struct {
	size  int64
	mtime int64
	caps  CapabilitySet
}

var (
	capabilityMu    sync.Mutex
	capabilityCache = map[string]cachedCapability{}
)

var flagPattern = regexp.MustCompile(`--[A-Za-z0-9][A-Za-z0-9-]*`)

// ProbeCLI runs the binary's own help and version commands and returns the
// advertised flag set. Results are cached by path, size and modification time,
// so a replaced binary is re-probed and an unchanged one costs nothing.
func ProbeCLI(binary string) CapabilitySet {
	if strings.TrimSpace(binary) == "" {
		return CapabilitySet{Error: "llama-server path was not resolved", Flags: map[string]bool{}}
	}
	info, statErr := os.Stat(binary)
	if statErr != nil {
		return CapabilitySet{Binary: binary, Error: "llama-server is not readable: " + statErr.Error(), Flags: map[string]bool{}}
	}
	var size, mtime int64
	size, mtime = info.Size(), info.ModTime().UnixNano()

	capabilityMu.Lock()
	if cached, ok := capabilityCache[binary]; ok && cached.size == size && cached.mtime == mtime {
		defer capabilityMu.Unlock()
		return cached.caps
	}
	capabilityMu.Unlock()

	caps := probeBinary(binary)

	capabilityMu.Lock()
	capabilityCache[binary] = cachedCapability{size: size, mtime: mtime, caps: caps}
	capabilityMu.Unlock()
	return caps
}

func probeBinary(binary string) CapabilitySet {
	caps := CapabilitySet{Binary: binary, Flags: map[string]bool{}}

	help, helpErr := runBinary(binary, "--help")
	caps.Flags = parseHelpFlags(help)
	if len(caps.Flags) == 0 {
		if helpErr != nil {
			caps.Error = "llama-server --help failed: " + helpErr.Error()
		} else {
			caps.Error = "llama-server --help advertised no recognizable flags"
		}
		return caps
	}
	caps.Probed = true

	version, _ := runBinary(binary, "--version")
	caps.Version, caps.Build, caps.Commit = parseVersion(version)
	return caps
}

func runBinary(binary, arg string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, binary, arg).CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return string(out), fmt.Errorf("%s timed out", arg)
	}
	// Some builds exit non-zero for --help while still printing the full help.
	// The text, not the exit status, is what this probe reads.
	return string(out), err
}

func parseHelpFlags(help string) map[string]bool {
	flags := map[string]bool{}
	for _, match := range flagPattern.FindAllString(help, -1) {
		flags[match] = true
	}
	return flags
}

var (
	versionPattern = regexp.MustCompile(`version:\s*(\S+)`)
	buildPattern   = regexp.MustCompile(`build\s+([0-9]+)`)
	commitPattern  = regexp.MustCompile(`commit\s+([0-9a-fA-F]+)`)
)

func parseVersion(text string) (version, build, commit string) {
	if m := versionPattern.FindStringSubmatch(text); m != nil {
		version = m[1]
	}
	if m := buildPattern.FindStringSubmatch(text); m != nil {
		build = m[1]
	}
	if m := commitPattern.FindStringSubmatch(text); m != nil {
		commit = m[1]
	}
	return version, build, commit
}
