//go:build native_test_seam

package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"howl-chat/internal/paths"
)

// nativeTestSeamEnv is deliberately build-tagged. Production builds contain
// only the no-op stub in app_native_test_seam_off.go.
const (
	nativeTestSeamEnv         = "HOWL_NATIVE_TEST_SEAM"
	nativeTestSeamTokenHeader = "X-Native-Test-Seam-Token"
	nativeTestSeamOrigin      = "http://wails.localhost"
)

type nativeTestSeam struct {
	app       *App
	server    *http.Server
	listener  net.Listener
	token     string
	mu        sync.Mutex
	pending   map[string]chan nativeTestSeamResult
	closeOnce sync.Once
}

type nativeTestSeamResult struct {
	Value json.RawMessage `json:"value"`
	Error string          `json:"error,omitempty"`
}

var (
	activeNativeTestSeamMu sync.Mutex
	activeNativeTestSeam   *nativeTestSeam
)

// startNativeTestSeam starts a loopback-only JavaScript evaluator in a test
// build. It is a no-op unless the explicit test switch is set and isolation has
// already been validated by paths.ResolveIsolation in main.
func startNativeTestSeam(a *App) {
	if os.Getenv(nativeTestSeamEnv) != "1" {
		return
	}
	if a == nil || a.ctx == nil {
		fmt.Fprintln(os.Stderr, "[native-test-seam] refusing: app context unavailable")
		return
	}
	if !paths.IsolationActive() {
		fmt.Fprintln(os.Stderr, "[native-test-seam] refusing: isolation is not active")
		return
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[native-test-seam] listen failed: %v\n", err)
		return
	}
	token, err := newNativeTestSeamToken()
	if err != nil {
		_ = ln.Close()
		fmt.Fprintf(os.Stderr, "[native-test-seam] token generation failed: %v\n", err)
		return
	}
	tokenPath := paths.StatePath("config", "native-test-seam.token")
	if err := os.MkdirAll(filepath.Dir(tokenPath), 0700); err != nil {
		_ = ln.Close()
		fmt.Fprintf(os.Stderr, "[native-test-seam] token directory failed: %v\n", err)
		return
	}
	if err := os.WriteFile(tokenPath, []byte(token+"\n"), 0600); err != nil {
		_ = ln.Close()
		fmt.Fprintf(os.Stderr, "[native-test-seam] token file failed: %v\n", err)
		return
	}
	if err := restrictNativeTestSeamToken(tokenPath); err != nil {
		_ = ln.Close()
		_ = os.Remove(tokenPath)
		fmt.Fprintf(os.Stderr, "[native-test-seam] token ACL failed: %v\n", err)
		return
	}
	seam := &nativeTestSeam{
		app:      a,
		listener: ln,
		token:    token,
		pending:  make(map[string]chan nativeTestSeamResult),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/__howl_native_test_seam/health", seam.handleHealth)
	mux.HandleFunc("/__howl_native_test_seam/eval", seam.handleEval)
	mux.HandleFunc("/__howl_native_test_seam/result", seam.handleResult)
	mux.HandleFunc("/__howl_native_test_seam/snapshot", seam.handleSnapshot)
	mux.HandleFunc("/__howl_native_test_seam/save-chat", seam.handleSaveChat)
	mux.HandleFunc("/__howl_native_test_seam/prompt-capture", seam.handlePromptCapture)
	seam.server = &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}

	activeNativeTestSeamMu.Lock()
	activeNativeTestSeam = seam
	activeNativeTestSeamMu.Unlock()

	go func() {
		if err := seam.server.Serve(ln); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "[native-test-seam] server stopped: %v\n", err)
		}
	}()
	fmt.Printf("[native-test-seam] listening on http://%s (isolated test build)\n", ln.Addr())
}

func stopNativeTestSeam() {
	activeNativeTestSeamMu.Lock()
	seam := activeNativeTestSeam
	activeNativeTestSeam = nil
	activeNativeTestSeamMu.Unlock()
	if seam == nil {
		return
	}
	seam.closeOnce.Do(func() {
		_ = seam.server.Close()
	})
}

func restrictNativeTestSeamToken(path string) error {
	if goruntime.GOOS != "windows" {
		return os.Chmod(path, 0600)
	}
	current, err := user.Current()
	if err != nil {
		return err
	}
	cmd := exec.Command("icacls", path, "/inheritance:r", "/grant:r", current.Username+":F")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("icacls: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func newNativeTestSeamToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// allowOrigin enforces the browser origin before any token is considered. A
// missing Origin is allowed for the local test driver; a browser origin must be
// the actual Wails page origin. CORS headers are only a response aid, never the
// authorization decision.
func (s *nativeTestSeam) allowOrigin(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	if origin != nativeTestSeamOrigin && origin != "https://wails.localhost" {
		http.Error(w, "origin not allowed", http.StatusForbidden)
		return false
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, "+nativeTestSeamTokenHeader)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	return true
}

func (s *nativeTestSeam) authorize(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get(nativeTestSeamTokenHeader) != s.token {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

func (s *nativeTestSeam) begin(w http.ResponseWriter, r *http.Request, method string) bool {
	if !s.allowOrigin(w, r) {
		return false
	}
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return false
	}
	if r.Method != method {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	if !s.authorize(w, r) {
		return false
	}
	return true
}

func (s *nativeTestSeam) handleHealth(w http.ResponseWriter, r *http.Request) {
	if !s.begin(w, r, http.MethodGet) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"ok":true,"isolation":true}`))
}

type nativeTestSeamEvalRequest struct {
	Expression string `json:"expression"`
}

type nativeTestSeamResultEnvelope struct {
	ID    string          `json:"id"`
	Value json.RawMessage `json:"value"`
	Error string          `json:"error"`
}

func (s *nativeTestSeam) handleEval(w http.ResponseWriter, r *http.Request) {
	if !s.begin(w, r, http.MethodPost) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req nativeTestSeamEvalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, "invalid expression", http.StatusBadRequest)
		return
	}

	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		http.Error(w, "could not allocate request id", http.StatusInternalServerError)
		return
	}
	id := hex.EncodeToString(idBytes)
	ch := make(chan nativeTestSeamResult, 1)
	s.mu.Lock()
	s.pending[id] = ch
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.pending, id)
		s.mu.Unlock()
	}()

	exprJSON, _ := json.Marshal(req.Expression)
	idJSON, _ := json.Marshal(id)
	tokenJSON, _ := json.Marshal(s.token)
	endpointJSON, _ := json.Marshal("http://" + s.listener.Addr().String() + "/__howl_native_test_seam/result")
	script := fmt.Sprintf(`(() => {
  const id = %s;
  const endpoint = %s;
  const token = %s;
  const expression = %s;
  Promise.resolve().then(async () => {
    let value = null;
    let error = '';
    try {
      value = await (0, eval)(expression);
      if (value === undefined) value = null;
    } catch (e) {
      error = String((e && e.stack) || e);
    }
    let body;
    try {
      body = JSON.stringify({id, value, error});
    } catch (e) {
      body = JSON.stringify({id, value: null, error: 'result is not JSON-serializable'});
    }	    try {
	      await fetch(endpoint, {
	        method: 'POST',
	        headers: {'Content-Type': 'application/json', 'X-Native-Test-Seam-Token': token},
	        body
	      });
	    } catch (e) {}
	  })();
	})();`, idJSON, endpointJSON, tokenJSON, exprJSON)

	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				ch <- nativeTestSeamResult{Error: fmt.Sprintf("exec panic: %v", recovered)}
			}
		}()
		runtime.WindowExecJS(s.app.ctx, script)
	}()

	select {
	case result := <-ch:
		w.Header().Set("Content-Type", "application/json")
		if result.Error != "" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": result.Error})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "value": result.Value})
	case <-time.After(10 * time.Second):
		http.Error(w, "native evaluation timed out", http.StatusGatewayTimeout)
	}
}

type nativeTestSeamSnapshotRequest struct {
	WorldID        string   `json:"worldId"`
	ScenarioID     string   `json:"scenarioId"`
	WorldName      string   `json:"worldName"`
	ScenarioName   string   `json:"scenarioName"`
	CharacterNames []string `json:"characterNames"`
	ChatName       string   `json:"chatName"`
}

func (s *nativeTestSeam) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	if !s.begin(w, r, http.MethodPost) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req nativeTestSeamSnapshotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ChatName == "" {
		http.Error(w, "invalid snapshot request", http.StatusBadRequest)
		return
	}
	if s.app.chatService == nil {
		http.Error(w, "chat service unavailable", http.StatusServiceUnavailable)
		return
	}
	if err := s.app.chatService.LoadChatData(req.WorldID, req.ScenarioID, req.WorldName, req.ScenarioName, req.CharacterNames, req.ChatName); err != nil {
		http.Error(w, "snapshot load failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"ids":     s.app.chatService.NativeTestSnapshotIDs(),
		"context": s.app.chatService.GetContextStatus(),
	})
}

type nativeTestSeamSaveChatRequest struct {
	ChatName string `json:"chatName"`
}

func (s *nativeTestSeam) handleSaveChat(w http.ResponseWriter, r *http.Request) {
	if !s.begin(w, r, http.MethodPost) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req nativeTestSeamSaveChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ChatName) == "" {
		http.Error(w, "invalid save-chat request", http.StatusBadRequest)
		return
	}
	if s.app.chatService == nil {
		http.Error(w, "chat service unavailable", http.StatusServiceUnavailable)
		return
	}
	if err := s.app.chatService.NativeTestSaveChatData(strings.TrimSpace(req.ChatName)); err != nil {
		http.Error(w, "save-chat failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":  true,
		"ids": s.app.chatService.NativeTestSnapshotIDs(),
	})
}

func (s *nativeTestSeam) handleResult(w http.ResponseWriter, r *http.Request) {
	if !s.begin(w, r, http.MethodPost) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4<<20)
	var envelope nativeTestSeamResultEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		http.Error(w, "invalid result", http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	ch := s.pending[envelope.ID]
	s.mu.Unlock()
	if ch == nil {
		http.Error(w, "unknown request", http.StatusNotFound)
		return
	}
	select {
	case ch <- nativeTestSeamResult{Value: envelope.Value, Error: envelope.Error}:
	default:
	}
	w.WriteHeader(http.StatusNoContent)
}
