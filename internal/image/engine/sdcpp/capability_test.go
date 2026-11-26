package sdcpp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"howl-chat/internal/image/types"
)

func TestParseHelpFlagsReadsStates(t *testing.T) {
	help := `stable-diffusion.cpp version unknown, commit 3b12781
Svr Options:
  --max-pending-jobs <int>         maximum number of pending async jobs (default: 64)
  --clip-on-cpu                    deprecated; use --backend te=cpu
  -p, --prompt <string>            the prompt to render
  --weight-residency-ttl <int>     idle seconds before unloading model components
                                        from RAM and VRAM (continued line)
`
	set := parseHelpFlags(help)
	cases := map[string]CapState{
		"--max-pending-jobs":     CapSupported,
		"--clip-on-cpu":          CapDeprecated,
		"--prompt":               CapSupported,
		"--weight-residency-ttl": CapSupported,
	}
	for flag, want := range cases {
		if got := set[flag]; got != want {
			t.Errorf("parseHelpFlags[%q] = %q, want %q", flag, got, want)
		}
	}
}

// A flag the engine dropped must read as removed, not as an unexplained
// unknown, because that is the sentence the refusal shows the user.
func TestCapabilityStateDistinguishesRemovedFromUnknown(t *testing.T) {
	set := CapabilitySet{Commit: "3b12781", Flags: map[string]CapState{"--vae": CapSupported}, Probed: true}
	if st := set.State("--motion-module"); st != CapRemoved {
		t.Fatalf("motion module state = %q, want removed", st)
	}
	if st := set.State("--invented-flag"); st != CapUnknown {
		t.Fatalf("unknown flag state = %q, want unknown", st)
	}
	if !set.Accepts("--vae") || set.Accepts("--motion-module") {
		t.Fatal("Accepts must follow the probed set")
	}
	if !strings.Contains(set.Explain("--motion-module"), "motion modules") {
		t.Fatalf("Explain should carry the reason: %q", set.Explain("--motion-module"))
	}
}

// The guard is what turns "sd-server exits on unknown argument" into a named
// refusal before a process is started.
func TestUnsupportedFlagsNamesTheFlag(t *testing.T) {
	restore := probeBinary
	probeBinary = func(string) CapabilitySet {
		return CapabilitySet{Commit: "3b12781", Probed: true, Flags: map[string]CapState{"--vae": CapSupported}}
	}
	defer func() { probeBinary = restore }()

	issues := UnsupportedFlags("sd-server.exe", []string{"--vae", "v", "--motion-module", "mm"})
	if len(issues) != 1 {
		t.Fatalf("want exactly one issue, got %v", issues)
	}
	if !strings.Contains(issues[0], "--motion-module") {
		t.Fatalf("issue must name the flag: %q", issues[0])
	}
}

// An unprobed binary (no help text) must not block a launch: refusing because
// probing failed would replace the engine's own error with a worse one.
func TestUnsupportedFlagsIgnoresUnprobedBinary(t *testing.T) {
	restore := probeBinary
	probeBinary = func(string) CapabilitySet { return CapabilitySet{} }
	defer func() { probeBinary = restore }()

	if issues := UnsupportedFlags("missing.exe", []string{"--motion-module", "mm"}); issues != nil {
		t.Fatalf("unprobed binary must not produce issues, got %v", issues)
	}
}

// The idle-release timers are opt-in: absent options must change nothing.
func TestWeightTTLFlagsAreOptIn(t *testing.T) {
	e := &Engine{config: types.EngineConfig{Port: 8081}}
	m := &types.ModelInfo{Name: "sd15", Family: "sd1.5", Files: map[string]string{"model": `C:\m\sd15.safetensors`}}

	plain := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{}), " ")
	if strings.Contains(plain, "--weight-residency-ttl") || strings.Contains(plain, "--weight-staging-ttl") {
		t.Fatalf("TTL flags must not be emitted by default: %s", plain)
	}

	residency := 30
	withTTL := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{
		WeightResidencyTTL: &residency,
	}), " ")
	if !strings.Contains(withTTL, "--weight-residency-ttl 30") {
		t.Errorf("missing residency TTL in %s", withTTL)
	}
	if strings.Contains(withTTL, "--weight-staging-ttl") {
		t.Errorf("staging TTL must not be forwarded: %s", withTTL)
	}
}

// Motion modules are gone from the engine line; the refusal must happen before
// sd-cli is started, and must name the model and the build.
func TestMotionModuleRefusedOnEngineWithoutIt(t *testing.T) {
	restore := probeBinary
	probeBinary = func(string) CapabilitySet {
		return CapabilitySet{Commit: "3b12781", Probed: true, Flags: map[string]CapState{"-m": CapSupported}}
	}
	defer func() { probeBinary = restore }()

	e := &Engine{config: types.EngineConfig{Port: 8081}, currentModel: &types.ModelInfo{
		Name:   "animatediff-sd15",
		Family: "sd1.5",
		Files:  map[string]string{"model": "sd15.safetensors", "motion_module": "mm.ckpt"},
	}}
	_, err := e.buildVideoArgs(types.ImageOptions{Prompt: "x"}, 16, 8, 512, 512, "o.webm")
	if err == nil {
		t.Fatal("expected a refusal for a motion module on a build without the flag")
	}
	if !strings.Contains(err.Error(), "motion module") || !strings.Contains(err.Error(), "3b12781") {
		t.Fatalf("refusal must name the capability and the build: %v", err)
	}
}

// A build that still takes the flag keeps working, so an older binary is not
// broken by the removal.
func TestMotionModuleStillPassedToEngineThatAcceptsIt(t *testing.T) {
	restore := probeBinary
	probeBinary = func(string) CapabilitySet {
		return CapabilitySet{Commit: "c6dfa9d", Probed: true, Flags: map[string]CapState{"--motion-module": CapSupported}}
	}
	defer func() { probeBinary = restore }()

	e := &Engine{config: types.EngineConfig{Port: 8081}, currentModel: &types.ModelInfo{
		Name:   "animatediff-sd15",
		Family: "sd1.5",
		Files:  map[string]string{"model": "sd15.safetensors", "motion_module": "mm.ckpt"},
	}}
	args, err := e.buildVideoArgs(types.ImageOptions{Prompt: "x"}, 16, 8, 512, 512, "o.webm")
	if err != nil {
		t.Fatalf("buildVideoArgs: %v", err)
	}
	if !strings.Contains(strings.Join(args, " "), "--motion-module mm.ckpt") {
		t.Fatalf("flag must be passed to a build that accepts it: %v", args)
	}
}

// The API layer must read the endpoint the engine actually serves, and a
// feature the server reports false must read as removed rather than absent.
func TestProbeAPIReadsCapabilities(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != capEndpoint {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"features":{"ip_adapter_image":false,"ref_images":true},"samplers":["euler"],"schedulers":["karras"],"supported_modes":["img_gen","vid_gen"]}`))
	}))
	defer srv.Close()

	caps := RefreshAPI(srv.URL)
	if !caps.Reached {
		t.Fatal("capabilities endpoint was reached but not reported")
	}
	if st := caps.FeatureState("ref_images"); st != CapAdvertised {
		t.Errorf("ref_images = %q, want advertised", st)
	}
	if st := caps.FeatureState("ip_adapter_image"); st != CapRemoved {
		t.Errorf("ip_adapter_image = %q, want removed", st)
	}
	if st := caps.FeatureState("never_heard_of_it"); st != CapUnknown {
		t.Errorf("unknown feature = %q, want unknown", st)
	}
	if len(caps.Modes) != 2 || caps.Modes[0] != "img_gen" {
		t.Errorf("modes not decoded: %v", caps.Modes)
	}
}

// An old build without the endpoint must not read as a build without features.
func TestProbeAPIMissingEndpointIsUnknown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	caps := RefreshAPI(srv.URL)
	if caps.Reached {
		t.Fatal("404 endpoint must not count as reached")
	}
	if st := caps.FeatureState("ref_images"); st != CapUnknown {
		t.Fatalf("feature on unreached server = %q, want unknown", st)
	}
}
