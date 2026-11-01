package llama

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestClient builds an HTTPClient pointed at a fake llama-server.
func newTestClient(t *testing.T, handler http.HandlerFunc) *HTTPClient {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return &HTTPClient{
		baseURL:    srv.URL,
		httpClient: srv.Client(),
		enabled:    true,
	}
}

func TestGetModelInfoReadsProps(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/props" {
			t.Errorf("unexpected path %q", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"model_path": `E:\AI GGUF\GGUF\mira-12b.Q4_K_M.gguf`,
			"default_generation_settings": map[string]interface{}{
				"n_ctx": float64(8192),
			},
		})
	})

	info, err := client.GetModelInfo()
	if err != nil {
		t.Fatalf("GetModelInfo returned error: %v", err)
	}
	if info.Name != "mira-12b.Q4_K_M.gguf" {
		t.Errorf("Name = %q, want mira-12b.Q4_K_M.gguf", info.Name)
	}
	if info.Parameters.ContextSize != 8192 {
		t.Errorf("ContextSize = %d, want 8192", info.Parameters.ContextSize)
	}
	if info.FilePath != `E:\AI GGUF\GGUF\mira-12b.Q4_K_M.gguf` {
		t.Errorf("FilePath = %q", info.FilePath)
	}
}

func TestGetContextSizeFallsBackToTopLevel(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"n_ctx": float64(4096)})
	})

	if got := client.GetContextSize(); got != 4096 {
		t.Errorf("GetContextSize = %d, want 4096", got)
	}
}

func TestTokenizeAndDetokenize(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/tokenize":
			var req map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("bad tokenize body: %v", err)
			}
			if req["content"] != "hello world" {
				t.Errorf("tokenize content = %v", req["content"])
			}
			json.NewEncoder(w).Encode(map[string]interface{}{
				"tokens": []interface{}{float64(15043), float64(3186)},
			})
		case "/detokenize":
			var req struct {
				Tokens []int `json:"tokens"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("bad detokenize body: %v", err)
			}
			if len(req.Tokens) != 2 {
				t.Errorf("detokenize tokens = %v", req.Tokens)
			}
			json.NewEncoder(w).Encode(map[string]interface{}{"content": "hello world"})
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
			http.NotFound(w, r)
		}
	})

	tokens, err := client.Tokenize("hello world", nil)
	if err != nil {
		t.Fatalf("Tokenize error: %v", err)
	}
	if len(tokens) != 2 || tokens[0] != 15043 || tokens[1] != 3186 {
		t.Fatalf("Tokenize = %v, want [15043 3186]", tokens)
	}

	text, err := client.Detokenize(tokens)
	if err != nil {
		t.Fatalf("Detokenize error: %v", err)
	}
	if text != "hello world" {
		t.Errorf("Detokenize = %q", text)
	}
}

func TestTokenizeSurfacesServerErrors(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "no vocab", http.StatusInternalServerError)
	})

	if _, err := client.Tokenize("x", nil); err == nil {
		t.Fatal("expected an error from a 500 response")
	}
}

func TestPropsRequestFailsWhenDisabled(t *testing.T) {
	client := &HTTPClient{enabled: false}
	if _, err := client.GetModelInfo(); err == nil {
		t.Fatal("expected disabled client to error")
	}
}
