package sdcpp

// Capability answers are made in three layers, and each layer has exactly one
// authority:
//
//	CLI   — parsed from the engine binary's own --help output. Answers "does
//	        this binary accept the option?".
//	API   — fetched from the running server's /capabilities endpoint. Answers
//	        "does the server implement the matching request field?". A flag in
//	        --help does not prove this, so it is probed separately.
//	Model — familyCapabilities plus companion resolution in this package.
//	        Answers "does the selected family support this at all?".
//
// None of the three is a hand-written table of the current engine's features.
// That is deliberate: the two defects this replaced were both table drift — a
// chroma entry emitting --chroma-disable-dit-mask after the binary had dropped
// it, and a clip_g requirement with no way to resolve.
//
// A capability is reported in one of five states, so "the engine took the flag
// away" reads differently from "we never knew about it":

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// CapState is how a capability stands against the engine in use.
type CapState string

const (
	// CapSupported: the binary lists it and does not mark it deprecated.
	CapSupported CapState = "supported"
	// CapDeprecated: the binary still accepts it and says it is deprecated.
	CapDeprecated CapState = "deprecated"
	// CapRemoved: the binary does not have it and the engine is known to have
	// dropped it. Emitting this would be a dead launch, not a slow one.
	CapRemoved CapState = "removed"
	// CapUnknown: the binary does not have it and this package has no record of
	// it either. Treated exactly like Removed when starting a process.
	CapUnknown CapState = "unknown"
	// CapAdvertised: the API says the feature exists. It has not been exercised
	// by HOWL, so it is not called supported.
	CapAdvertised CapState = "advertised"
)

// retiredFlags records flags the engine line removed, so an absent flag can be
// reported as "removed by the engine" instead of "unknown". It is history, not
// capability: nothing here is ever used to decide whether a flag works.
var retiredFlags = map[string]string{
	"--motion-module":               "motion modules were dropped from the engine",
	"--slg-scale":                   "skip-layer guidance was dropped from the engine",
	"--skip-layers":                 "skip-layer guidance was dropped from the engine",
	"--skip-layer-start":            "skip-layer guidance was dropped from the engine",
	"--skip-layer-end":              "skip-layer guidance was dropped from the engine",
	"--high-noise-cfg-scale":        "high-noise sampling was dropped from the engine",
	"--high-noise-guidance":         "high-noise sampling was dropped from the engine",
	"--high-noise-img-cfg-scale":    "high-noise sampling was dropped from the engine",
	"--high-noise-eta":              "high-noise sampling was dropped from the engine",
	"--high-noise-sampling-method":  "high-noise sampling was dropped from the engine",
	"--high-noise-skip-layers":      "high-noise sampling was dropped from the engine",
	"--high-noise-skip-layer-start": "high-noise sampling was dropped from the engine",
	"--high-noise-skip-layer-end":   "high-noise sampling was dropped from the engine",
	"--timestep-shift":              "timestep shifting was replaced by --flow-shift",
	"--ip-adapter":                  "IP-Adapter support was dropped from the engine",
	"--ip-adapter-image":            "IP-Adapter support was dropped from the engine",
	"--ip-adapter-strength":         "IP-Adapter support was dropped from the engine",
	"--pulid-weights":               "PuLID support was dropped from the engine",
	"--pulid-id-embedding":          "PuLID support was dropped from the engine",
	"--pulid-id-weight":             "PuLID support was dropped from the engine",
	"--photo-maker":                 "PhotoMaker support was dropped from the engine",
	"--chroma-disable-dit-mask":     "Chroma's DIT-mask switch became a --model-args key",
}

// CapabilitySet is the CLI layer: what the binary on disk accepts.
type CapabilitySet struct {
	Commit string
	// Flags maps every flag the binary prints in --help to its state.
	Flags map[string]CapState
	// Probed is false when --help could not be run or produced no flags; callers
	// must not read an empty set as "the engine supports nothing".
	Probed bool
}

// State reports how the binary stands on one flag.
func (c CapabilitySet) State(flag string) CapState {
	if st, ok := c.Flags[flag]; ok {
		return st
	}
	if _, ok := retiredFlags[flag]; ok {
		return CapRemoved
	}
	return CapUnknown
}

// Accepts reports whether the binary will take the flag without exiting.
func (c CapabilitySet) Accepts(flag string) bool {
	switch c.State(flag) {
	case CapSupported, CapDeprecated:
		return true
	}
	return false
}

// Explain returns a human sentence for a flag this binary cannot take, or ""
// when it can.
func (c CapabilitySet) Explain(flag string) string {
	switch c.State(flag) {
	case CapSupported, CapDeprecated:
		return ""
	case CapRemoved:
		return fmt.Sprintf("%s is not accepted by sd-server %s: %s", flag, c.Commit, retiredFlags[flag])
	default:
		return fmt.Sprintf("%s is not accepted by sd-server %s and is not a flag this build knows", flag, c.Commit)
	}
}

var (
	capMu     sync.Mutex
	capCache  = map[string]CapabilitySet{}
	helpCache = map[string]string{}
)

// helpOutput runs `<binary> --help` once per binary path and keeps the text.
// Version detection and flag discovery read the same capture: two probes of one
// binary were two process launches.
func helpOutput(binary string) string {
	if binary == "" {
		return ""
	}
	capMu.Lock()
	out, ok := helpCache[binary]
	capMu.Unlock()
	if ok {
		return out
	}
	out = runHelp(binary)
	capMu.Lock()
	helpCache[binary] = out
	capMu.Unlock()
	return out
}

// ProbeCLI parses the binary's own help into a capability set. The result is
// cached per binary path.
func ProbeCLI(binary string) CapabilitySet {
	if binary == "" {
		return CapabilitySet{}
	}
	capMu.Lock()
	if c, ok := capCache[binary]; ok {
		capMu.Unlock()
		return c
	}
	capMu.Unlock()

	out := helpOutput(binary)
	set := CapabilitySet{Commit: BinaryVersion(binary), Flags: parseHelpFlags(out), Probed: len(parseHelpFlags(out)) > 0}

	capMu.Lock()
	capCache[binary] = set
	capMu.Unlock()
	return set
}

// parseHelpFlags reads flag names and their deprecated marker out of a help
// screen. Continuation lines (descriptions that wrapped) are ignored: they do
// not start with a dash.
func parseHelpFlags(help string) map[string]CapState {
	flags := map[string]CapState{}
	for _, line := range strings.Split(help, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "-") {
			continue
		}
		name, desc := splitFlagLine(trimmed)
		if name == "" {
			continue
		}
		if strings.Contains(strings.ToLower(desc), "deprecated") {
			flags[name] = CapDeprecated
			continue
		}
		flags[name] = CapSupported
	}
	return flags
}

// splitFlagLine separates "--flag <arg>    description" from its description.
// A leading short alias ("-p, --prompt") is skipped.
func splitFlagLine(line string) (name, desc string) {
	rest := line
	if strings.HasPrefix(rest, "-") && !strings.HasPrefix(rest, "--") {
		if idx := strings.Index(rest, "--"); idx > 0 {
			rest = rest[idx:]
		}
	}
	if !strings.HasPrefix(rest, "--") {
		return "", ""
	}
	end := strings.IndexAny(rest, " \t")
	if end < 0 {
		return rest, ""
	}
	name = rest[:end]
	name = strings.TrimSuffix(name, ",")
	if i := strings.Index(name, "="); i > 0 {
		name = name[:i]
	}
	desc = strings.TrimSpace(rest[end:])
	// The argument placeholder is part of the flag column, not the description.
	if strings.HasPrefix(desc, "<") {
		if i := strings.Index(desc, ">"); i >= 0 {
			desc = strings.TrimSpace(desc[i+1:])
		}
	}
	return name, desc
}

// probeBinary is a seam for tests: production callers always ask the binary on
// disk. A test can install a capability set without shipping a fake executable.
var probeBinary = ProbeCLI

// UnsupportedFlags returns the flags in an argv this binary would refuse,
// with the sentence to show the user. An unprobed binary (no help text) yields
// nothing: refusing to start because probing failed would replace a real error
// with a worse one.
func UnsupportedFlags(binary string, args []string) []string {
	set := probeBinary(binary)
	if !set.Probed {
		return nil
	}
	seen := map[string]bool{}
	var out []string
	for _, a := range args {
		if !strings.HasPrefix(a, "--") || seen[a] {
			continue
		}
		seen[a] = true
		if why := set.Explain(a); why != "" {
			out = append(out, why)
		}
	}
	sort.Strings(out)
	return out
}

// ---- API layer -------------------------------------------------------------

// capEndpoint is where this engine line advertises what it implements. It was
// named that in the release that introduced it; /capabilities answers 404.
const capEndpoint = "/sdcpp/v1/capabilities"

// APICapabilities is the API layer: what the running server says it implements.
// It is read from /capabilities, which reports a features map (ip_adapter_image,
// ref_images, high_noise_sample_params, …) plus the samplers, schedulers and
// generation modes the build was compiled with.
type APICapabilities struct {
	// Reached is false for an engine build that has no /capabilities endpoint (or
	// a server that is not up). Nothing is claimed in that case.
	Reached    bool
	Features   map[string]bool
	Samplers   []string
	Schedulers []string
	Modes      []string
}

// FeatureState reports a request-level feature. Advertised means the server says
// it implements the field; it does not mean HOWL has verified it end to end.
func (a APICapabilities) FeatureState(name string) CapState {
	if !a.Reached {
		return CapUnknown
	}
	if v, ok := a.Features[name]; ok {
		if v {
			return CapAdvertised
		}
		return CapRemoved
	}
	return CapUnknown
}

var (
	apiMu    sync.Mutex
	apiCache = map[string]APICapabilities{}
)

// ProbeAPI asks the running server what it implements. Failures are reported as
// "not reached", never as "no features": an old build without the endpoint must
// not read as a build without capabilities.
func ProbeAPI(baseURL string) APICapabilities {
	if baseURL == "" {
		return APICapabilities{}
	}
	apiMu.Lock()
	if c, ok := apiCache[baseURL]; ok {
		apiMu.Unlock()
		return c
	}
	apiMu.Unlock()

	caps := fetchCapabilities(baseURL)
	apiMu.Lock()
	apiCache[baseURL] = caps
	apiMu.Unlock()
	return caps
}

// RefreshAPI drops the cached answer for a base URL, for callers that want the
// server's current view (after a model change, for example).
func RefreshAPI(baseURL string) APICapabilities {
	apiMu.Lock()
	delete(apiCache, baseURL)
	apiMu.Unlock()
	return ProbeAPI(baseURL)
}

func fetchCapabilities(baseURL string) APICapabilities {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(strings.TrimSuffix(baseURL, "/") + capEndpoint)
	if err != nil {
		return APICapabilities{}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return APICapabilities{}
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return APICapabilities{}
	}
	var doc struct {
		Features       map[string]bool `json:"features"`
		Samplers       []string        `json:"samplers"`
		Schedulers     []string        `json:"schedulers"`
		Modes          []string        `json:"modes"`
		SupportedModes []string        `json:"supported_modes"`
	}
	if err := json.Unmarshal(body, &doc); err != nil {
		return APICapabilities{}
	}
	if doc.Features == nil {
		doc.Features = map[string]bool{}
	}
	modes := doc.SupportedModes
	if len(modes) == 0 {
		modes = doc.Modes
	}
	return APICapabilities{
		Reached:    true,
		Features:   doc.Features,
		Samplers:   doc.Samplers,
		Schedulers: doc.Schedulers,
		Modes:      modes,
	}
}
