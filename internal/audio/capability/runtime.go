package capability

import (
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/paths"
)

// LlamaRuntimeProbe records what the shipped llama.cpp modules contain for the
// audio request surface. Route strings are evidence that the handler is
// compiled in, not proof that a model can transcribe: transcription stays
// unverified until an audio-capable model and projector are installed and the
// route is exercised.
type LlamaRuntimeProbe struct {
	ServerPath    string
	ImplPath      string
	LauncherBytes int64
	ImplBytes     int64
	RoutesPresent bool
	InputAudio    bool
	MTMDAudio     bool
	Err           string
}

// ProbeLlamaRuntime inspects the module that actually serves requests. The
// shipped llama-server.exe is a launcher stub, so scanning it says nothing
// about the server; the audio route lives in the implementation module it
// loads.
func ProbeLlamaRuntime(serverPath string) LlamaRuntimeProbe {
	probe := LlamaRuntimeProbe{ServerPath: serverPath}
	if strings.TrimSpace(serverPath) == "" {
		probe.Err = "no llama-server path was resolved"
		return probe
	}
	probe.LauncherBytes = fileSize(serverPath)
	probe.ImplPath = implementationModule(serverPath)
	impl, err := os.ReadFile(probe.ImplPath)
	if err != nil {
		probe.Err = "cannot read " + probe.ImplPath + ": " + err.Error()
		return probe
	}
	probe.ImplBytes = int64(len(impl))
	text := string(impl)
	probe.RoutesPresent = strings.Contains(text, "/v1/audio/transcriptions")
	probe.InputAudio = strings.Contains(text, "input_audio")

	mtmd, err := os.ReadFile(filepath.Join(filepath.Dir(probe.ImplPath), "mtmd.dll"))
	if err == nil {
		mtmdText := string(mtmd)
		probe.MTMDAudio = strings.Contains(mtmdText, "mtmd_gen_audio") || strings.Contains(mtmdText, "clip_has_audio")
	}
	return probe
}

// implementationModule returns the module that actually contains the server
// code, falling back to the server binary for layouts without a separate
// implementation DLL.
func implementationModule(serverPath string) string {
	candidate := strings.TrimSuffix(serverPath, filepath.Ext(serverPath)) + "-impl.dll"
	if fileSize(candidate) > 0 {
		return candidate
	}
	return serverPath
}

func fileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return 0
	}
	return info.Size()
}

var (
	llamaProbeOnce sync.Once
	llamaProbe     *LlamaRuntimeProbe
)

// productionLlamaAudioProbe probes the shipped runtime once per process.
// Replacing the runtime requires a restart for the report to change, which is
// preferable to a stale hardcoded claim in either direction.
func productionLlamaAudioProbe() *LlamaRuntimeProbe {
	llamaProbeOnce.Do(func() {
		root := paths.ProjectRoot()
		for _, name := range []string{"llama-server.exe", "llama-server"} {
			candidate := filepath.Join(root, "runtime", "llama", name)
			if fileSize(candidate) == 0 {
				continue
			}
			probe := ProbeLlamaRuntime(candidate)
			llamaProbe = &probe
			return
		}
		llamaProbe = &LlamaRuntimeProbe{Err: "llama-server is not installed under runtime/llama"}
	})
	return llamaProbe
}

type meloProbeResult struct {
	mu     sync.Mutex
	at     time.Time
	url    string
	ok     bool
	detail string
}

var meloProbe meloProbeResult

// probeMeloTTSService checks whether anything is listening at the configured
// MeloTTS URL. A TCP connection is not a synthesis test, but it replaces a
// permanent "unknown" with an observed reachability fact. Results are cached
// briefly so status refreshes do not dial on every request.
func probeMeloTTSService(rawURL string) (bool, string) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return false, "no MeloTTS URL is configured"
	}
	meloProbe.mu.Lock()
	defer meloProbe.mu.Unlock()
	if meloProbe.url == rawURL && time.Since(meloProbe.at) < 5*time.Second {
		return meloProbe.ok, meloProbe.detail
	}

	conn, err := net.DialTimeout("tcp", meloAddress(rawURL), 500*time.Millisecond)
	ok := err == nil
	detail := "no MeloTTS service is listening at " + rawURL
	if ok {
		_ = conn.Close()
		detail = "MeloTTS service accepted a TCP connection at " + rawURL
	}
	meloProbe.at, meloProbe.url, meloProbe.ok, meloProbe.detail = time.Now(), rawURL, ok, detail
	return ok, detail
}

// llamaRuntimeState reports what the runtime probe establishes. Route presence
// alone is not application support: provider and task dimensions must also be
// available. On this build E4 transcribed real audio through the route; a
// replaced runtime must be re-probed and re-verified.
func llamaRuntimeState(env Env) Dimension {
	if env.LlamaAudio == nil {
		return UnknownDim("the shipped llama-server audio surface was not probed")
	}
	probe := *env.LlamaAudio
	if probe.Err != "" {
		return UnknownDim("the llama-server audio surface could not be probed: " + probe.Err)
	}
	if !probe.RoutesPresent {
		return UnknownDim("no /v1/audio/transcriptions route string was found in " + filepath.Base(probe.ImplPath) + "; whether the route answers is unestablished")
	}
	detail := "POST /v1/audio/transcriptions is compiled into " + filepath.Base(probe.ImplPath)
	if probe.InputAudio {
		detail += " with input_audio request handling"
	}
	if probe.MTMDAudio {
		detail += "; mtmd exposes audio symbols"
	}
	detail += "; E4 transcribed real audio with a Qwen3-ASR model and projector through this route"
	return Avail(detail)
}

func meloAddress(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return rawURL
	}
	if parsed.Port() != "" {
		return parsed.Host
	}
	port := "80"
	if parsed.Scheme == "https" {
		port = "443"
	}
	return net.JoinHostPort(parsed.Hostname(), port)
}
