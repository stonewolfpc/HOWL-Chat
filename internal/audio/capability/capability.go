// Package capability is the provider-neutral audio capability contract.
//
// A capability is reported per task, not per provider name, and every fact that
// decides whether a task can run is a separate dimension: provider presence,
// model presence, projector compatibility, runtime/build support, task support,
// hardware eligibility, and how far the capability has actually been verified.
//
// Nothing here is inferred from a file name. A record is selectable only when
// every dimension says so; an unknown dimension is reported as unknown and
// blocks selection rather than being rounded up to support.
package capability

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"howl-chat/internal/audio/hw"
	"howl-chat/internal/audio/types"
	"howl-chat/internal/paths"
)

// Task identifies what the user is asking the audio system to do.
type Task string

const (
	TaskFileTranscription      Task = "file_transcription"
	TaskStreamingTranscription Task = "streaming_transcription"
	TaskAudioUnderstanding     Task = "audio_understanding"
	TaskRealtimeConversation   Task = "realtime_conversation"
	TaskSpeechSynthesis        Task = "speech_synthesis"
)

// Tasks returns every task the contract represents, in reporting order.
func Tasks() []Task {
	return []Task{
		TaskFileTranscription,
		TaskStreamingTranscription,
		TaskAudioUnderstanding,
		TaskRealtimeConversation,
		TaskSpeechSynthesis,
	}
}

// Availability is deliberately tri-state: "unknown" is not "unavailable", and
// neither is "available".
type Availability string

const (
	AvailabilityUnknown       Availability = "unknown"
	AvailabilityAvailable     Availability = "available"
	AvailabilityUnavailable   Availability = "unavailable"
	AvailabilityNotApplicable Availability = "not_applicable"
)

// Dimension is one fact about a capability, carrying its evidence or blocker.
type Dimension struct {
	State  Availability `json:"state"`
	Detail string       `json:"detail,omitempty"`
}

// Avail reports a satisfied dimension with the evidence that satisfied it.
func Avail(detail string) Dimension {
	return Dimension{State: AvailabilityAvailable, Detail: detail}
}

// Blocked reports a dimension that prevents use, with the reason.
func Blocked(detail string) Dimension {
	return Dimension{State: AvailabilityUnavailable, Detail: detail}
}

// UnknownDim reports a dimension that could not be established.
func UnknownDim(detail string) Dimension {
	return Dimension{State: AvailabilityUnknown, Detail: detail}
}

// NoRequirement reports a dimension that does not apply to this capability.
func NoRequirement(detail string) Dimension {
	return Dimension{State: AvailabilityNotApplicable, Detail: detail}
}

// OK reports whether the dimension permits use. Only an established
// availability or a genuine non-requirement permits it.
func (d Dimension) OK() bool {
	return d.State == AvailabilityAvailable || d.State == AvailabilityNotApplicable
}

// Record is one task/provider/model capability.
type Record struct {
	ID       string `json:"id"`
	Task     Task   `json:"task"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Local    bool   `json:"local"`

	ProviderState Dimension `json:"provider_state"`
	ModelState    Dimension `json:"model_state"`
	Projector     Dimension `json:"projector"`
	Runtime       Dimension `json:"runtime"`
	TaskSupport   Dimension `json:"task_support"`
	Hardware      Dimension `json:"hardware"`

	Verification types.Verification `json:"verification"`
	Reason       string             `json:"reason,omitempty"`
}

// Selectable reports whether every dimension permits use. Verification is not
// part of selectability: a provider whose adapter exists can run before it has
// been graded end to end, but it may not be reported as verified.
func (r Record) Selectable() bool {
	return r.ProviderState.OK() &&
		r.ModelState.OK() &&
		r.Projector.OK() &&
		r.Runtime.OK() &&
		r.TaskSupport.OK() &&
		r.Hardware.OK()
}

// Verified reports whether the capability has been proven end to end on real
// assets. Only this permits claiming the user-facing feature works.
func (r Record) Verified() bool {
	return r.Selectable() && r.Verification == types.VerificationEndToEnd
}

// Blockers lists every reason the capability cannot be selected, in a stable
// order, so UI and status can explain rather than just disable.
func (r Record) Blockers() []string {
	type named struct {
		label string
		dim   Dimension
	}
	dims := []named{
		{"provider", r.ProviderState},
		{"model", r.ModelState},
		{"projector", r.Projector},
		{"runtime", r.Runtime},
		{"task", r.TaskSupport},
		{"hardware", r.Hardware},
	}
	var out []string
	for _, entry := range dims {
		if entry.dim.OK() {
			continue
		}
		reason := string(entry.dim.State)
		if entry.dim.Detail != "" {
			reason += ": " + entry.dim.Detail
		}
		out = append(out, entry.label+" "+reason)
	}
	return out
}

func (r Record) withReason() Record {
	if blockers := r.Blockers(); len(blockers) > 0 {
		r.Reason = strings.Join(blockers, "; ")
	}
	return r
}

// Env carries every external fact the contract needs. It is injectable so the
// catalog can be tested without a machine, and so the contract never shells out
// to discover capability while reporting status.
type Env struct {
	Root         string
	LookPath     func(string) (string, bool)
	Getenv       func(string) string
	FileExists   func(string) bool
	GPUAvailable bool

	// ASRModelPath, ASRProjectorPath and TTSModelPath are the same path
	// policies the pipeline loads from. Supplying them here keeps one source of
	// truth.
	ASRModelPath     func(types.ASRModelType) string
	ASRProjectorPath func(types.ASRModelType) string
	TTSModelPath     func(types.TTSModelType) string

	// LlamaAudio is the shipped runtime's audio-surface probe. A nil probe
	// means the surface was never inspected and is reported as unknown, never
	// as absent. Route strings prove the handler is compiled in; they do not
	// prove a model can transcribe.
	LlamaAudio *LlamaRuntimeProbe

	// MeloTTSReachable probes local MeloTTS readiness. Nil leaves it unknown
	// (tests). The dimension is "is this a ready Howl provider", not "does a
	// port answer": the application supplies the managed service's ownership
	// and liveness probe. The capability package's own fallback is a bounded
	// reachability check and must not be used to advertise an owned provider.
	MeloTTSReachable func(url string) (bool, string)
}

// ProductionEnv builds the environment the application reports.
//
// The llama.cpp audio surface is probed from the shipped implementation module,
// never asserted from a launcher stub or an upstream release note. The probe is
// cached for the process, so replacing the runtime requires a restart before
// the report changes.
func ProductionEnv(asrModelPath, asrProjectorPath func(types.ASRModelType) string, ttsModelPath func(types.TTSModelType) string) Env {
	return Env{
		Root:             paths.ProjectRoot(),
		ASRModelPath:     asrModelPath,
		ASRProjectorPath: asrProjectorPath,
		TTSModelPath:     ttsModelPath,
		GPUAvailable:     detectGPU(),
		LlamaAudio:       productionLlamaAudioProbe(),
		MeloTTSReachable: probeMeloTTSService,
	}
}

var (
	gpuOnce    sync.Once
	gpuPresent bool
)

func detectGPU() bool {
	gpuOnce.Do(func() {
		gpuPresent = hw.DetectHardware().HasGPU
	})
	return gpuPresent
}

// Catalog returns every capability record for the current environment.
func Catalog(env Env) []Record {
	env = env.withDefaults()
	records := make([]Record, 0, 20)
	records = append(records, asrRecords(env)...)
	records = append(records, runtimeAudioRecords(env)...)
	records = append(records, ttsRecords(env)...)
	for i := range records {
		records[i] = records[i].withReason()
	}
	return records
}

func (e Env) withDefaults() Env {
	if e.LookPath == nil {
		e.LookPath = func(name string) (string, bool) {
			path, err := exec.LookPath(name)
			return path, err == nil
		}
	}
	if e.Getenv == nil {
		e.Getenv = os.Getenv
	}
	if e.FileExists == nil {
		e.FileExists = func(path string) bool {
			_, err := os.Stat(path)
			return err == nil
		}
	}
	if e.Root == "" {
		e.Root = paths.ProjectRoot()
	}
	return e
}

// fileExists resolves relative paths against the project root, matching how the
// pipeline's own path policies behave.
func (e Env) fileExists(path string) bool {
	if path == "" {
		return false
	}
	resolved := path
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(e.Root, resolved)
	}
	if e.FileExists != nil && e.FileExists(path) {
		return true
	}
	if resolved != path && e.FileExists != nil {
		return e.FileExists(resolved)
	}
	return false
}

func (e Env) lookWhisperCLI() (string, bool) {
	for _, name := range []string{"whisper-cli", "whisper-cli.exe"} {
		if path, ok := e.LookPath(name); ok {
			return path, true
		}
	}
	for _, candidate := range []string{
		filepath.Join(e.Root, "bin", "whisper", "whisper-cli.exe"),
		filepath.Join(e.Root, "bin", "whisper-cli.exe"),
	} {
		if e.fileExists(candidate) {
			return candidate, true
		}
	}
	return "", false
}

func asrRecords(env Env) []Record {
	modelTypes := make([]types.ASRModelType, 0, len(types.ASRModelRegistry))
	for modelType := range types.ASRModelRegistry {
		modelTypes = append(modelTypes, modelType)
	}
	sort.Slice(modelTypes, func(i, j int) bool { return modelTypes[i] < modelTypes[j] })

	out := make([]Record, 0, len(modelTypes)*2)
	for _, modelType := range modelTypes {
		profile := types.ASRModelRegistry[modelType]
		base := asrBase(env, modelType, profile)

		file := base
		file.Task = TaskFileTranscription
		file.TaskSupport = asrTaskSupport(profile, TaskFileTranscription)
		out = append(out, file.withReason())

		stream := base
		stream.Task = TaskStreamingTranscription
		stream.TaskSupport = asrTaskSupport(profile, TaskStreamingTranscription)
		out = append(out, stream.withReason())
	}
	return out
}

func asrBase(env Env, modelType types.ASRModelType, profile types.ASRModelProfile) Record {
	record := Record{
		ID:           string(modelType),
		Provider:     profile.Runtime,
		Model:        profile.Name,
		Local:        profile.Local,
		Verification: profile.Verification,
		Projector:    NoRequirement("this adapter needs no projector"),
		Hardware:     hardwareState(env, profile.ModelProfile),
	}
	if record.Provider == "" {
		record.Provider = string(modelType)
	}
	if record.Verification == "" {
		record.Verification = types.VerificationUnverified
	}
	record.ProviderState, record.Runtime = providerState(env, profile.ModelProfile)
	record.ModelState = modelState(env, modelType, profile.ModelProfile, record.ProviderState)
	if profile.ProjectorKind == types.ProjectorKindAudio {
		path := ""
		if env.ASRProjectorPath != nil {
			path = env.ASRProjectorPath(modelType)
		}
		switch {
		case strings.TrimSpace(path) == "":
			record.Projector = Blocked("an audio projector is required but no matching projector was found")
		case !env.fileExists(path):
			record.Projector = Blocked("audio projector is not installed: " + path)
		default:
			record.Projector = Avail("matching audio projector found: " + path)
		}
	}
	return record
}

func providerState(env Env, profile types.ModelProfile) (provider Dimension, runtime Dimension) {
	switch profile.Runtime {
	case types.RuntimeWhisperCLI:
		path, ok := env.lookWhisperCLI()
		if !ok {
			dim := Blocked("whisper-cli is not on PATH and not in bin/whisper")
			return dim, dim
		}
		dim := Avail("whisper-cli found at " + path)
		return dim, dim

	case types.RuntimeDashScopeASR:
		if strings.TrimSpace(env.Getenv("HOWL_QWEN_API_KEY")) == "" {
			dim := Blocked("HOWL_QWEN_API_KEY is not set")
			return dim, dim
		}
		return Avail("DashScope API credentials are configured"),
			Avail("the application HTTP client implements the DashScope ASR task")

	case types.RuntimeLlamaServerASR:
		return Avail("the llama-server ASR adapter is registered in this build"), llamaRuntimeState(env)
	}

	return UnknownDim("no provider is declared for this capability"),
		UnknownDim("no runtime support is declared for this capability")
}

func modelState(env Env, modelType types.ASRModelType, profile types.ModelProfile, provider Dimension) Dimension {
	if !profile.Local {
		if provider.OK() {
			return Avail("the model is served by the configured remote provider")
		}
		return Blocked("a remote model cannot be used while its provider is unavailable")
	}
	if !profile.ApprovedLocal {
		return Blocked("this model is not admitted for local use; the admission gate has not passed")
	}

	path := ""
	if env.ASRModelPath != nil {
		path = env.ASRModelPath(modelType)
	}
	if strings.TrimSpace(path) == "" {
		return Blocked("no local model path was resolved")
	}
	if !env.fileExists(path) {
		return Blocked("model file is not installed: " + path)
	}
	return Avail("model file is present: " + path)
}

func hardwareState(env Env, profile types.ModelProfile) Dimension {
	if profile.SupportsCPU {
		return Avail("the adapter supports CPU execution")
	}
	if profile.SupportsGPU {
		if env.GPUAvailable {
			return Avail("the adapter requires a GPU and one is present")
		}
		return Blocked("the adapter requires a GPU and none was detected")
	}
	return Blocked("the adapter declares no supported compute device")
}

func asrTaskSupport(profile types.ASRModelProfile, task Task) Dimension {
	switch task {
	case TaskFileTranscription:
		if profile.Local && !profile.ApprovedLocal {
			return Blocked("no adapter is admitted for this model; the admission gate has not passed")
		}
		return Avail("the adapter implements file transcription")
	case TaskStreamingTranscription:
		if !profile.SupportsStreaming {
			return Blocked("the adapter cannot stream; TranscribeStream returns unsupported")
		}
		return Avail("the adapter declares streaming transcription")
	}
	return UnknownDim("the adapter does not represent this task")
}

func runtimeAudioRecords(env Env) []Record {
	runtimeDim := llamaRuntimeState(env)

	return []Record{
		{
			ID:            "llama.cpp/audio-understanding",
			Task:          TaskAudioUnderstanding,
			Provider:      "llama.cpp chat completions (input_audio)",
			Model:         "multimodal GGUF with an audio projector",
			Local:         true,
			ProviderState: Blocked("no application code constructs a multimodal audio request"),
			ModelState:    UnknownDim("audio model compatibility is not negotiated"),
			Projector:     UnknownDim("projector kind is read from GGUF metadata; no audio projector is enumerated"),
			Runtime:       runtimeDim,
			TaskSupport:   Blocked("no adapter implements direct audio understanding"),
			Hardware:      UnknownDim("no audio model is loaded, so eligibility is unmeasured"),
			Verification:  types.VerificationUnverified,
		},
		{
			ID:            "llama.cpp/realtime-conversation",
			Task:          TaskRealtimeConversation,
			Provider:      "llama.cpp realtime audio",
			Model:         "realtime-capable audio model",
			Local:         true,
			ProviderState: Blocked("no realtime audio adapter exists in the application"),
			ModelState:    UnknownDim("no realtime-capable model has been identified"),
			Projector:     UnknownDim("no projector has been negotiated for realtime audio"),
			Runtime:       runtimeDim,
			TaskSupport:   Blocked("no adapter implements realtime conversation"),
			Hardware:      UnknownDim("no realtime model is loaded, so eligibility is unmeasured"),
			Verification:  types.VerificationUnverified,
		},
	}
}

func ttsRecords(env Env) []Record {
	order := []types.TTSModelType{types.TTSPiper, types.TTSMeloTTS, types.TTSKokoro}
	out := make([]Record, 0, len(order))
	for _, modelType := range order {
		profile, ok := types.ModelRegistry.TTSModels[modelType]
		if !ok {
			continue
		}
		record := Record{
			ID:            string(modelType),
			Task:          TaskSpeechSynthesis,
			Provider:      "speech synthesis adapter",
			Model:         profile.Name,
			Local:         profile.Backend == "local",
			ProviderState: Avail("the synthesis adapter is registered in this build"),
			Runtime:       Avail("the adapter is part of the application build"),
			Projector:     NoRequirement("speech synthesis needs no projector"),
			TaskSupport:   Avail("the adapter implements speech synthesis"),
			Hardware:      hardwareState(env, profile.ModelProfile),
			Verification:  profile.Verification,
		}
		if record.Verification == "" {
			record.Verification = types.VerificationUnverified
		}

		if profile.Backend == "http" {
			url := strings.TrimSpace(env.Getenv("HOWL_MELOTTS_URL"))
			if url == "" {
				url = "http://127.0.0.1:8888"
			}
			switch {
			case env.MeloTTSReachable == nil:
				record.ModelState = UnknownDim("service reachability at " + url + " is not verified")
			default:
				if ok, detail := env.MeloTTSReachable(url); ok {
					record.ModelState = Avail(detail)
				} else {
					record.ModelState = Blocked(detail)
				}
			}
		} else {
			path := ""
			if env.TTSModelPath != nil {
				path = env.TTSModelPath(modelType)
			}
			switch {
			case strings.TrimSpace(path) == "":
				record.ModelState = Blocked("no voice model path was resolved")
			case !env.fileExists(path):
				record.ModelState = Blocked("voice model is not installed: " + path)
			default:
				record.ModelState = Avail("voice model is present: " + path)
			}
		}
		out = append(out, record.withReason())
	}
	return out
}

// Select returns the capability to use for a task. An explicit id is honoured
// or refused with its blockers; it is never quietly replaced by another
// provider. "auto" (or blank) chooses among selectable records.
func Select(records []Record, task Task, id string) (Record, error) {
	id = strings.TrimSpace(id)
	if id == "" || id == "auto" {
		return AutoSelect(records, task)
	}
	for _, record := range records {
		if record.Task != task || record.ID != id {
			continue
		}
		if !record.Selectable() {
			return Record{}, fmt.Errorf("capability %q is not available for %s: %s", id, task, strings.Join(record.Blockers(), "; "))
		}
		return record, nil
	}
	return Record{}, fmt.Errorf("no capability %q for task %q", id, task)
}

// AutoSelect chooses deterministically: verified first, then implemented, then
// local, then by id. It never crosses tasks.
func AutoSelect(records []Record, task Task) (Record, error) {
	candidates := make([]Record, 0, len(records))
	for _, record := range records {
		if record.Task == task && record.Selectable() {
			candidates = append(candidates, record)
		}
	}
	if len(candidates) == 0 {
		return Record{}, fmt.Errorf("no selectable capability for task %q", task)
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		left, right := verificationRank(candidates[i].Verification), verificationRank(candidates[j].Verification)
		if left != right {
			return left > right
		}
		if candidates[i].Local != candidates[j].Local {
			return candidates[i].Local
		}
		return candidates[i].ID < candidates[j].ID
	})
	return candidates[0], nil
}

// Fallbacks returns providers that may replace chosen after a runtime failure.
// They must serve the same task and the same locality: a local transcription
// failure may not silently become a remote API call.
func Fallbacks(records []Record, chosen Record) []Record {
	var out []Record
	for _, record := range records {
		if record.Task != chosen.Task || record.ID == chosen.ID || record.Local != chosen.Local {
			continue
		}
		if !record.Selectable() {
			continue
		}
		out = append(out, record)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func verificationRank(v types.Verification) int {
	switch v {
	case types.VerificationEndToEnd:
		return 3
	case types.VerificationImplementation:
		return 2
	case types.VerificationUnverified:
		return 1
	default:
		return 0
	}
}
