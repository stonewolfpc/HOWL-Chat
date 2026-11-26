package sdcpp

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"howl-chat/internal/image/types"
	"howl-chat/internal/resource"
)

// blockingTransport holds a generation request open after Shutdown cancels the
// request context. It deliberately ignores cancellation so the test can tell
// "shutdown waited for the execution scope" from "the cancelled HTTP request
// happened to return first".
type blockingTransport struct {
	entered chan struct{}
	release chan struct{}
	once    sync.Once
}

func (b *blockingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	b.once.Do(func() { close(b.entered) })
	<-b.release
	return nil, fmt.Errorf("transport released")
}

// Shutdown cancels and joins the active generation scope before it returns.
func TestShutdownJoinsGeneration(t *testing.T) {
	e := NewEngine()
	e.httpClient = newClient("http://unused.invalid")
	e.currentModel = &types.ModelInfo{Name: "test-model", Family: "z-image"}
	transport := &blockingTransport{
		entered: make(chan struct{}),
		release: make(chan struct{}),
	}
	e.httpClient.http.Transport = transport

	generationDone := make(chan struct{})
	go func() {
		defer close(generationDone)
		_, _ = e.GenerateTxt2Img(types.ImageOptions{JobID: "job-shutdown", Prompt: "x"})
	}()

	select {
	case <-transport.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("generation never reached the transport")
	}

	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)
		_ = e.Shutdown()
	}()

	// Let the test end cleanly even when the assertion fails.
	defer func() {
		close(transport.release)
		select {
		case <-generationDone:
		case <-time.After(2 * time.Second):
		}
	}()

	select {
	case <-shutdownDone:
		select {
		case <-generationDone:
			// Correct: shutdown waited for the execution scope.
		default:
			t.Fatal("Shutdown returned while the generation was still blocked")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Shutdown did not return")
	}
}

// The engine reports its exit through the hook but leaves resource
// registrations to the app/domain authority.
func TestEngineShutdownReportsExitWithoutUnregistering(t *testing.T) {
	const name = "shutdown-behavioral-image"
	resource.GetManager().UnregisterModel(resource.ResourceImage, name)
	t.Cleanup(func() {
		resource.GetManager().UnregisterModel(resource.ResourceImage, name)
	})
	resource.GetManager().RegisterModel(
		resource.ResourceImage, name, "/image/model.safetensors", 100, 0, resource.PriorityImage, true,
	)

	e := NewEngine()
	hookCalls := 0
	e.SetShutdownHook(func() { hookCalls++ })

	if err := e.Shutdown(); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
	if hookCalls != 1 {
		t.Fatalf("shutdown hook calls = %d, want 1", hookCalls)
	}
	if got := resource.GetManager().ModelsByType(resource.ResourceImage); len(got) == 0 {
		t.Fatal("engine Shutdown removed the image registration; authority owns that transition")
	}
}
