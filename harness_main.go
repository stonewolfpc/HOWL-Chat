//go:build uiharness

package main

// harness_main.go — development-only UI harness.
//
// Runs the real HOWL Chat backend (same NewApp, same bindings) and serves
// frontend/dist to an ordinary browser over HTTP, so the interface can be
// clicked and screenshotted by an agent instead of a human at the window.
//
// Built with:  go build -tags uiharness -o build/bin/howl-ui-harness .
// Run with:    build/bin/howl-ui-harness.exe      (listens on 127.0.0.1:8799)
//
// Every window.go.main.App.* call the frontend makes is proxied by reflection
// to the same method the desktop build would call, so this exercises real
// backend code rather than mocks. Calls that need a native dialog (file
// pickers, message boxes) have no webview to attach to and return an error;
// everything else behaves like the shipped app.
//
// This file never participates in the release build: main.go carries the
// inverse !uiharness constraint.

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

//go:embed harness/shim.js
var harnessShim embed.FS

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// eventHub fans app-emitted events out to browser pages over Server-Sent
// Events. Without it the harness can proxy binding calls but never deliver
// backend-originated events, so streamed replies and reasoning never render.
type eventHub struct {
	mu   sync.Mutex
	subs map[chan []byte]struct{}
}

func newEventHub() *eventHub {
	return &eventHub{subs: make(map[chan []byte]struct{})}
}

func (h *eventHub) add() chan []byte {
	ch := make(chan []byte, 256)
	h.mu.Lock()
	h.subs[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *eventHub) remove(ch chan []byte) {
	h.mu.Lock()
	if _, ok := h.subs[ch]; ok {
		delete(h.subs, ch)
		close(ch)
	}
	h.mu.Unlock()
}

func (h *eventHub) publish(name string, args []interface{}) {
	if args == nil {
		args = []interface{}{}
	}
	payload, err := json.Marshal(map[string]interface{}{"name": name, "args": args})
	if err != nil {
		log.Printf("[events] marshal %s: %v", name, err)
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.subs {
		select {
		case ch <- payload:
		default: // slow page: drop rather than stall generation
		}
	}
}

func main() {
	addr := "127.0.0.1:8799"
	if v := os.Getenv("HOWL_HARNESS_ADDR"); v != "" {
		addr = v
	}

	hub := newEventHub()
	devEventSink = func(name string, args ...interface{}) { hub.publish(name, args) }

	app := NewApp()
	// NewApp composes the diagnostics owner and attaches its emitter, so errors
	// reach devEventSink even though this host never runs Wails' startup. Set
	// the context the bindings expect and run the model scan that startup would
	// normally kick off in the background.
	app.ctx = context.Background()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[harness] background rescan panic: %v", r)
			}
		}()
		app.RescanAllDirs()
	}()

	dist := "frontend/dist"
	if _, err := os.Stat(dist); err != nil {
		log.Fatalf("[harness] frontend/dist not found (run from the project root): %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/__rpc", func(w http.ResponseWriter, r *http.Request) { handleRPC(w, r, app) })
	mux.HandleFunc("/__events", func(w http.ResponseWriter, r *http.Request) { handleEvents(w, r, hub) })
	// The harness has no window to close, so this is its shutdown control: it
	// runs the same Level 0 sequence Wails' OnShutdown runs in the packaged app,
	// then exits. Without it the harness could only be killed, and an orderly
	// shutdown could never be exercised in this host.
	mux.HandleFunc("/__shutdown", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[harness] shutdown requested")
		err := app.OnShutdownError(context.Background())
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			_, _ = fmt.Fprintf(w, `{"ok":false,"error":%q}`, err.Error())
		} else {
			_, _ = w.Write([]byte(`{"ok":true}`))
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		go func() {
			time.Sleep(250 * time.Millisecond)
			os.Exit(0)
		}()
	})
	mux.HandleFunc("/__harness.js", func(w http.ResponseWriter, r *http.Request) {
		b, err := harnessShim.ReadFile("harness/shim.js")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		_, _ = w.Write(b)
	})
	mux.Handle("/", http.FileServer(http.Dir(dist)))

	// HTML pages are served through a rewriter so the shim always loads before
	// any page script asks for window.go.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Never let a browser cache assets between edits: a stale stylesheet
		// makes a CSS fix look broken and wastes a whole debugging pass.
		if r.URL.Path != "/__events" {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}
		if r.URL.Path == "/" || strings.HasSuffix(r.URL.Path, ".html") {
			serveHTMLWithShim(w, r, dist)
			return
		}
		mux.ServeHTTP(w, r)
	})

	log.Printf("[harness] serving %s on http://%s", dist, addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

// handleEvents streams app events to the page as Server-Sent Events.
func handleEvents(w http.ResponseWriter, r *http.Request, hub *eventHub) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	ch := hub.add()
	defer hub.remove(ch)

	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-heartbeat.C:
			if _, err := fmt.Fprint(w, ": keepalive\n\n"); err != nil {
				return
			}
			flusher.Flush()
		case payload, ok := <-ch:
			if !ok {
				return
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", payload); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func serveHTMLWithShim(w http.ResponseWriter, r *http.Request, dist string) {
	rel := strings.TrimPrefix(r.URL.Path, "/")
	if rel == "" {
		rel = "index.html"
	}
	path := dist + "/" + rel
	body, err := os.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	html := string(body)
	tag := `<script src="/__harness.js"></script>`
	inserted := false
	for _, marker := range []string{"<head>", "<head ", "<HEAD>"} {
		if idx := strings.Index(html, marker); idx >= 0 {
			end := strings.Index(html[idx:], ">")
			if end >= 0 {
				cut := idx + end + 1
				html = html[:cut] + "\n    " + tag + html[cut:]
				inserted = true
				break
			}
		}
	}
	if !inserted {
		html = tag + "\n" + html
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write([]byte(html))
}

type rpcRequest struct {
	Method string            `json:"method"`
	Args   []json.RawMessage `json:"args"`
}

type rpcResponse struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
}

func handleRPC(w http.ResponseWriter, r *http.Request, app *App) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req rpcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeRPCError(w, fmt.Errorf("bad request: %w", err))
		return
	}

	start := time.Now()
	result, err := invoke(app, req.Method, req.Args)
	elapsed := time.Since(start)

	if err != nil {
		log.Printf("[rpc] %-45s FAILED  %8s  %v", req.Method, elapsed.Round(time.Millisecond), err)
		writeRPCError(w, err)
		return
	}
	log.Printf("[rpc] %-45s ok      %8s", req.Method, elapsed.Round(time.Millisecond))
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rpcResponse{Result: result})
}

func writeRPCError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(rpcResponse{Error: err.Error()})
}

// invoke calls the named method on the App exactly as Wails' reflection
// binding would, converting JSON arguments to the declared parameter types.
func invoke(app *App, name string, args []json.RawMessage) (out interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in %s: %v", name, r)
		}
	}()

	method := reflect.ValueOf(app).MethodByName(name)
	if !method.IsValid() {
		return nil, fmt.Errorf("no such binding: %s", name)
	}
	mt := method.Type()
	if mt.NumIn() != len(args) {
		return nil, fmt.Errorf("%s expects %d arguments, got %d", name, mt.NumIn(), len(args))
	}

	in := make([]reflect.Value, mt.NumIn())
	for i := 0; i < mt.NumIn(); i++ {
		ptr := reflect.New(mt.In(i))
		if len(args[i]) > 0 && string(args[i]) != "null" {
			if uerr := json.Unmarshal(args[i], ptr.Interface()); uerr != nil {
				return nil, fmt.Errorf("%s argument %d: %w", name, i, uerr)
			}
		}
		in[i] = ptr.Elem()
	}

	results := method.Call(in)
	if n := len(results); n > 0 {
		last := results[n-1]
		if last.Type().Implements(errorType) {
			if !last.IsNil() {
				return nil, last.Interface().(error)
			}
			results = results[:n-1]
		}
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0].Interface(), nil
}
