package runtime

import (
	"sync"
	"testing"

	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/types"
)

type stubClient struct{ loaded bool }

func (s *stubClient) GenerateChat([]llama.ChatMessage, *llama.InferenceOptions) (string, error) {
	return "", nil
}
func (s *stubClient) GenerateChatStream([]llama.ChatMessage, *llama.InferenceOptions, llama.TokenCallback) error {
	return nil
}
func (s *stubClient) GenerateChatStreamWithImage([]llama.ChatMessage, []byte, *llama.InferenceOptions, llama.TokenCallback) error {
	return nil
}
func (s *stubClient) CancelRequest()                                  {}
func (s *stubClient) GetContextSize() int                             { return 4096 }
func (s *stubClient) GetLoadedModel() (map[string]interface{}, error) { return nil, nil }
func (s *stubClient) Health() error                                   { return nil }
func (s *stubClient) IsLoaded() bool                                  { return s.loaded }

func TestEngineZeroValueAndNilAreUsable(t *testing.T) {
	var engine *Engine
	if engine.Client() != nil || engine.Models() != nil || engine.Profile() != nil {
		t.Fatal("a nil engine reported state")
	}
	if engine.Progress() != 0 || engine.ProgressMilli() != 0 || engine.Stage() != "" || engine.AutoTemplate() != "" {
		t.Fatal("a nil engine reported load state")
	}
	// Writes are dropped rather than panicking: an App built without boot has
	// no runtime to record into, and must still answer its bindings.
	engine.SetProgressMilli(500)
	engine.SetStage("Loading")
	engine.SetProfile(&types.ModelProfile{Family: "qwen"})
	engine.SetAutoTemplate("tpl")
	engine.SetClient(&stubClient{})
	engine.SetModels(nil)
	if engine.ProgressMilli() != 0 {
		t.Fatal("a nil engine recorded progress")
	}

	zero := &Engine{}
	if zero.Progress() != 0 || zero.Stage() != "" {
		t.Fatal("the zero engine is not an empty session")
	}
	zero.SetStage("Ready")
	if zero.Stage() != "Ready" {
		t.Fatalf("stage = %q, want Ready", zero.Stage())
	}
}

func TestEngineHoldsSessionAndLoadBookkeeping(t *testing.T) {
	client := &stubClient{loaded: true}
	engine := NewEngine(client, nil)

	if engine.Client() != Client(client) {
		t.Fatal("Client did not return the installed client")
	}
	if !engine.Client().IsLoaded() {
		t.Fatal("client reports unloaded")
	}

	engine.SetProgressMilli(475)
	if got := engine.Progress(); got != 0.475 {
		t.Fatalf("Progress = %v, want 0.475", got)
	}
	engine.SetProgressMilli(1000)
	if got := engine.Progress(); got != 1 {
		t.Fatalf("Progress = %v, want 1", got)
	}

	engine.SetStage("Loading weights")
	if got := engine.Stage(); got != "Loading weights" {
		t.Fatalf("Stage = %q", got)
	}

	profile := &types.ModelProfile{Family: "llama", MMPProjPath: "mmproj.gguf"}
	engine.SetProfile(profile)
	if engine.Profile() != profile {
		t.Fatal("Profile did not return the recorded profile")
	}

	if engine.AutoTemplate() != "" {
		t.Fatal("AutoTemplate started non-empty")
	}
	engine.SetAutoTemplate("model-native")
	if got := engine.AutoTemplate(); got != "model-native" {
		t.Fatalf("AutoTemplate = %q", got)
	}
	engine.SetAutoTemplate("replacement")
	if got := engine.AutoTemplate(); got != "replacement" {
		t.Fatalf("AutoTemplate replacement = %q", got)
	}

	replacement := &stubClient{}
	engine.SetClient(replacement)
	if engine.Client() != Client(replacement) {
		t.Fatal("SetClient did not replace the client")
	}
}

func TestEngineProgressAndStageAreSafeUnderConcurrentAccess(t *testing.T) {
	engine := NewEngine(&stubClient{}, nil)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for tick := 0; tick < 200; tick++ {
				engine.SetProgressMilli(int64(tick) * 5)
				engine.SetStage("Loading weights")
				if got := engine.ProgressMilli(); got < 0 || got > 1000 {
					t.Errorf("progress out of range: %d", got)
				}
				if got := engine.Stage(); got != "" && got != "Loading weights" {
					t.Errorf("unexpected stage: %q", got)
				}
				_ = engine.Profile()
				_ = engine.AutoTemplate()
			}
		}()
	}
	wg.Wait()
}
