package capability

import (
	"strings"
	"testing"

	"howl-chat/internal/audio/types"
)

func presentFiles(paths ...string) func(string) bool {
	set := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		set[path] = struct{}{}
	}
	return func(path string) bool {
		_, ok := set[path]
		return ok
	}
}

// emptyEnv is the observable starting state of this machine: no whisper-cli, no
// ASR model files, no API key, and a runtime with no audio routes.
func emptyEnv() Env {
	return Env{
		Root:         `C:\project`,
		LookPath:     func(string) (string, bool) { return "", false },
		Getenv:       func(string) string { return "" },
		FileExists:   func(string) bool { return false },
		GPUAvailable: true,
		ASRModelPath: func(model types.ASRModelType) string {
			return `C:\project\models\whisper\` + string(model) + ".bin"
		},
	}
}

func find(t *testing.T, records []Record, task Task, id string) Record {
	t.Helper()
	for _, record := range records {
		if record.Task == task && record.ID == id {
			return record
		}
	}
	t.Fatalf("no record for task %q id %q", task, id)
	return Record{}
}

func TestCatalogCoversEveryTask(t *testing.T) {
	records := Catalog(emptyEnv())
	seen := make(map[Task]bool)
	for _, record := range records {
		seen[record.Task] = true
	}
	for _, task := range Tasks() {
		if !seen[task] {
			t.Fatalf("no capability record reports task %q", task)
		}
	}
}

func TestUnavailableProviderIsNotSelectableAndSaysWhy(t *testing.T) {
	records := Catalog(emptyEnv())

	record := find(t, records, TaskFileTranscription, string(types.ASRWhisperTiny))
	if record.Selectable() {
		t.Fatal("whisper is selectable without whisper-cli or a model file")
	}
	if record.ProviderState.State != AvailabilityUnavailable {
		t.Fatalf("provider state = %q, want unavailable", record.ProviderState.State)
	}
	if !strings.Contains(record.Reason, "whisper-cli") {
		t.Fatalf("reason = %q, want the missing binary named", record.Reason)
	}
	if record.Verified() {
		t.Fatal("an unavailable record reported itself verified")
	}
}

func TestConfiguredMachineSelectsWithoutClaimingVerification(t *testing.T) {
	env := emptyEnv()
	env.LookPath = func(name string) (string, bool) {
		if name == "whisper-cli" {
			return `C:\bin\whisper-cli.exe`, true
		}
		return "", false
	}
	env.FileExists = presentFiles(`C:\project\models\whisper\whisper-tiny.bin`)

	records := Catalog(env)
	chosen, err := Select(records, TaskFileTranscription, "whisper-tiny")
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	if !chosen.Selectable() {
		t.Fatalf("record is not selectable: %s", chosen.Reason)
	}
	if chosen.Verified() {
		t.Fatal("an implementation-only capability claimed end-to-end verification")
	}
	if chosen.Verification != types.VerificationImplementation {
		t.Fatalf("verification = %q, want implementation", chosen.Verification)
	}
}

func TestExplicitSelectionIsNeverSilentlyReplaced(t *testing.T) {
	records := Catalog(emptyEnv())

	_, err := Select(records, TaskFileTranscription, "whisper-tiny")
	if err == nil {
		t.Fatal("Select returned a usable capability for an unavailable provider")
	}
	if !strings.Contains(err.Error(), "whisper-tiny") {
		t.Fatalf("error = %v, want the requested id", err)
	}

	_, err = Select(records, TaskFileTranscription, "not-a-model")
	if err == nil {
		t.Fatal("Select accepted an unknown capability id")
	}
}

func TestAutoSelectUsesOnlySelectableRecords(t *testing.T) {
	env := emptyEnv()
	env.LookPath = func(name string) (string, bool) {
		if name == "whisper-cli" {
			return "whisper-cli", true
		}
		return "", false
	}
	env.FileExists = presentFiles(`C:\project\models\whisper\whisper-base.bin`)

	chosen, err := AutoSelect(Catalog(env), TaskFileTranscription)
	if err != nil {
		t.Fatalf("AutoSelect: %v", err)
	}
	if chosen.ID != string(types.ASRWhisperBase) {
		t.Fatalf("auto chose %q, want the only selectable model", chosen.ID)
	}
}

func TestFallbacksStayInsideTaskAndLocality(t *testing.T) {
	env := emptyEnv()
	env.LookPath = func(name string) (string, bool) {
		if name == "whisper-cli" {
			return "whisper-cli", true
		}
		return "", false
	}
	env.FileExists = presentFiles(
		`C:\project\models\whisper\whisper-tiny.bin`,
		`C:\project\models\whisper\whisper-base.bin`,
	)
	env.Getenv = func(key string) string {
		if key == "HOWL_QWEN_API_KEY" {
			return "configured"
		}
		return ""
	}

	records := Catalog(env)
	chosen, err := Select(records, TaskFileTranscription, string(types.ASRWhisperTiny))
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	fallbacks := Fallbacks(records, chosen)
	if len(fallbacks) != 1 || fallbacks[0].ID != string(types.ASRWhisperBase) {
		t.Fatalf("fallbacks = %+v, want only whisper-base", fallbacks)
	}
	for _, fallback := range fallbacks {
		if fallback.Local != chosen.Local {
			t.Fatalf("fallback %q crossed the local/remote boundary", fallback.ID)
		}
		if fallback.Task != chosen.Task {
			t.Fatalf("fallback %q crossed tasks", fallback.ID)
		}
	}
}

// DEBT-T2: a task may not be advertised when its adapter returns unsupported.
func TestStreamingIsReportedUnavailableForEveryAdapter(t *testing.T) {
	records := Catalog(emptyEnv())
	for _, record := range records {
		if record.Task != TaskStreamingTranscription {
			continue
		}
		if record.TaskSupport.OK() {
			t.Fatalf("%q advertises streaming while its adapter returns unsupported", record.ID)
		}
		if record.Selectable() {
			t.Fatalf("%q is selectable for streaming", record.ID)
		}
	}
}

// The audio-understanding and realtime records must not borrow vision support
// from the presence of a projector file.
func TestAudioUnderstandingAndRealtimeAreBlockedByTheRuntime(t *testing.T) {
	records := Catalog(emptyEnv())
	for _, task := range []Task{TaskAudioUnderstanding, TaskRealtimeConversation} {
		record := find(t, records, task, recordIDForTask(records, task))
		if record.Selectable() {
			t.Fatalf("%s is selectable without an adapter or runtime route", task)
		}
		if record.Projector.State == AvailabilityAvailable {
			t.Fatalf("%s claims projector capability without negotiation", task)
		}
	}
}

func recordIDForTask(records []Record, task Task) string {
	for _, record := range records {
		if record.Task == task {
			return record.ID
		}
	}
	return ""
}

func TestProductionEnvProbesTheShippedRuntime(t *testing.T) {
	env := ProductionEnv(nil, nil, nil)
	if env.LlamaAudio == nil {
		t.Fatal("ProductionEnv did not probe the shipped runtime")
	}
	if env.LlamaAudio.ImplPath == "" && env.LlamaAudio.Err == "" {
		t.Fatal("probe reported neither an implementation module nor an error")
	}
}

// A compiled-in route must not be rounded up to application support: the
// runtime dimension is available after E4, but the provider/task dimensions stay
// blocked without an adapter and the record is not selectable.
func TestRuntimeRouteEvidenceDoesNotImplyIntegration(t *testing.T) {
	env := emptyEnv()
	env.LlamaAudio = &LlamaRuntimeProbe{
		ImplPath:      `C:\project\runtime\llama\llama-server-impl.dll`,
		ImplBytes:     10088960,
		RoutesPresent: true,
		InputAudio:    true,
		MTMDAudio:     true,
	}

	record := find(t, Catalog(env), TaskAudioUnderstanding, "llama.cpp/audio-understanding")
	if record.Selectable() {
		t.Fatal("route strings made audio understanding selectable")
	}
	if record.Runtime.State != AvailabilityAvailable {
		t.Fatalf("runtime state = %q, want available after the E4 transcription", record.Runtime.State)
	}
	if record.ProviderState.State != AvailabilityUnavailable {
		t.Fatalf("provider state = %q, want unavailable without an application adapter", record.ProviderState.State)
	}
}

// A missing probe is unknown, never absence.
func TestUnprobedRuntimeIsUnknownNotAbsent(t *testing.T) {
	env := emptyEnv()
	record := find(t, Catalog(env), TaskAudioUnderstanding, "llama.cpp/audio-understanding")
	if record.Runtime.State != AvailabilityUnknown {
		t.Fatalf("runtime state = %q, want unknown", record.Runtime.State)
	}
}
