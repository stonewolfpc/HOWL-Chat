package main

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func serveThroughBridge(t *testing.T, payload string, headers map[string]string, status int) *httptest.ResponseRecorder {
	t.Helper()

	handler := bridgeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, value := range headers {
			w.Header().Set(key, value)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(payload))
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/lorebooks.html", nil))
	return recorder
}

func TestStandalonePageGetsWailsRuntime(t *testing.T) {
	page := "<html><head><title>Lorebooks</title></head><body>page</body></html>"
	recorder := serveThroughBridge(t, page, map[string]string{"Content-Type": "text/html; charset=utf-8"}, http.StatusOK)

	body := recorder.Body.String()
	if count := strings.Count(body, "/wails/runtime.js"); count != 1 {
		t.Fatalf("runtime.js injected %d times, want 1: %s", count, body)
	}
	if count := strings.Count(body, "/wails/ipc.js"); count != 1 {
		t.Fatalf("ipc.js injected %d times, want 1: %s", count, body)
	}
	if !strings.Contains(body, "<title>Lorebooks</title>") {
		t.Fatalf("page content lost: %s", body)
	}
	if got := recorder.Header().Get("Content-Length"); got != itoa(len(body)) {
		t.Fatalf("Content-Length = %q, want %d", got, len(body))
	}
}

// Wails inserts the runtime as the first child of <head>; a page must not be
// able to run its own head script before the bridge exists.
func TestInjectionLandsAtStartOfHead(t *testing.T) {
	page := `<html><head><script src="js/early.js"></script></head><body></body></html>`
	body := serveThroughBridge(t, page, map[string]string{"Content-Type": "text/html"}, http.StatusOK).Body.String()

	runtimeAt := strings.Index(body, "/wails/runtime.js")
	earlyAt := strings.Index(body, "js/early.js")
	if runtimeAt < 0 || earlyAt < 0 {
		t.Fatalf("missing markers: %s", body)
	}
	if runtimeAt > earlyAt {
		t.Fatalf("runtime injected after a head script: %s", body)
	}
	if strings.Count(body, "/wails/ipc.js") != 1 {
		t.Fatalf("ipc.js must be injected exactly once: %s", body)
	}
}

func TestMainPageIsNotInjectedTwice(t *testing.T) {
	page := `<html><head><script src="/wails/runtime.js"></script><script src="/wails/ipc.js"></script></head><body></body></html>`
	recorder := serveThroughBridge(t, page, map[string]string{"Content-Type": "text/html"}, http.StatusOK)

	if count := strings.Count(recorder.Body.String(), "/wails/runtime.js"); count != 1 {
		t.Fatalf("runtime.js present %d times, want 1", count)
	}
}

// A compressed body cannot be rewritten by insertion; it must pass through.
func TestCompressedHTMLIsNotRewritten(t *testing.T) {
	page := "<html><head></head><body>compressed</body></html>"
	recorder := serveThroughBridge(t, page, map[string]string{
		"Content-Type":     "text/html",
		"Content-Encoding": "gzip",
	}, http.StatusOK)

	if strings.Contains(recorder.Body.String(), "/wails/") {
		t.Fatalf("compressed body was rewritten: %s", recorder.Body.String())
	}
	if recorder.Header().Get("Content-Encoding") != "gzip" {
		t.Fatal("Content-Encoding was not preserved")
	}
}

func TestResponseHeadersSurviveInjection(t *testing.T) {
	page := "<html><head></head><body>x</body></html>"
	recorder := serveThroughBridge(t, page, map[string]string{
		"Content-Type":  "text/html; charset=utf-8",
		"Cache-Control": "no-cache",
		"ETag":          `"abc123"`,
	}, http.StatusOK)

	if got := recorder.Header().Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("Cache-Control = %q", got)
	}
	if got := recorder.Header().Get("ETag"); got != `"abc123"` {
		t.Fatalf("ETag = %q", got)
	}
	if got := recorder.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
		t.Fatalf("Content-Type = %q", got)
	}
}

// Wails sets Content-Length for the file it serves. Injection changes the body,
// so the header must describe the injected length, not the original file.
func TestContentLengthIsRecalculatedAfterInjection(t *testing.T) {
	page := "<html><head></head><body>abc</body></html>"
	recorder := serveThroughBridge(t, page, map[string]string{
		"Content-Type":   "text/html",
		"Content-Length": itoa(len(page)),
	}, http.StatusOK)

	body := recorder.Body.String()
	if got, want := recorder.Header().Get("Content-Length"), itoa(len(body)); got != want {
		t.Fatalf("Content-Length = %q, want %q (injected body length)", got, want)
	}
	if recorder.Header().Get("Content-Length") == itoa(len(page)) {
		t.Fatal("the original Content-Length was preserved across injection")
	}
}

// A flush must commit the injection rather than the raw prefix, must not lose
// content written afterwards, and cannot advertise a fixed length while the
// document is still arriving.
func TestFlushCommitsInjectedBody(t *testing.T) {
	handler := bridgeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", "999")
		_, _ = io.WriteString(w, "<html><head>")
		w.(http.Flusher).Flush()
		_, _ = io.WriteString(w, "</head><body>tail</body></html>")
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/lorebooks.html", nil))

	body := recorder.Body.String()
	if count := strings.Count(body, "/wails/runtime.js"); count != 1 {
		t.Fatalf("runtime.js injected %d times after flush: %s", count, body)
	}
	if !strings.Contains(body, "tail") {
		t.Fatalf("content written after the flush was lost: %s", body)
	}
	if got := recorder.Header().Get("Content-Length"); got == "999" {
		t.Fatal("flush kept a stale Content-Length")
	}
}

// Headers are committed by the first write, so an early flush must not leave a
// Content-Length behind and the later final write must not try to set one.
func TestEarlyFlushRemovesContentLength(t *testing.T) {
	handler := bridgeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", "999")
		_, _ = io.WriteString(w, "<html><head>")
		w.(http.Flusher).Flush()
		_, _ = io.WriteString(w, "</head><body>tail</body></html>")
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/lorebooks.html", nil))

	committed := recorder.Result()
	defer committed.Body.Close()
	if got := committed.Header.Get("Content-Length"); got != "" {
		t.Fatalf("committed Content-Length after an early flush = %q, want none", got)
	}
	if committed.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", committed.StatusCode)
	}
	body, err := io.ReadAll(committed.Body)
	if err != nil {
		t.Fatal(err)
	}
	if count := strings.Count(string(body), "/wails/runtime.js"); count != 1 {
		t.Fatalf("runtime.js injected %d times: %s", count, body)
	}
	if !strings.Contains(string(body), "tail") {
		t.Fatalf("content after the flush was lost: %s", body)
	}
}

// A flush before any body has arrived has nothing to send; committing headers
// there would make the eventual injection impossible.
func TestFlushBeforeBodyKeepsInjection(t *testing.T) {
	page := "<html><head></head><body>page</body></html>"
	handler := bridgeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.(http.Flusher).Flush()
		_, _ = io.WriteString(w, page)
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/lorebooks.html", nil))

	body := recorder.Body.String()
	if count := strings.Count(body, "/wails/runtime.js"); count != 1 {
		t.Fatalf("runtime.js injected %d times: %s", count, body)
	}
	if got, want := recorder.Header().Get("Content-Length"), itoa(len(body)); got != want {
		t.Fatalf("Content-Length = %q, want %q", got, want)
	}
}

func TestNonHTMLAndErrorsPassThrough(t *testing.T) {
	script := serveThroughBridge(t, "console.log(1)", map[string]string{"Content-Type": "application/javascript"}, http.StatusOK)
	if strings.Contains(script.Body.String(), "/wails/") {
		t.Fatalf("script response was modified: %s", script.Body.String())
	}

	missing := serveThroughBridge(t, "not found", map[string]string{"Content-Type": "text/plain"}, http.StatusNotFound)
	if missing.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", missing.Code)
	}
	if strings.Contains(missing.Body.String(), "/wails/") {
		t.Fatalf("error response was modified: %s", missing.Body.String())
	}
}

func TestUnwrapExposesUnderlyingWriter(t *testing.T) {
	recorder := httptest.NewRecorder()
	wrapper := &bridgeResponse{ResponseWriter: recorder}
	if wrapper.Unwrap() != recorder {
		t.Fatal("Unwrap must return the underlying writer for http.ResponseController")
	}
}

func TestEveryEmbeddedPageReceivesWailsRuntime(t *testing.T) {
	pages := 0
	err := fs.WalkDir(assets, "frontend/dist", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".html") {
			return nil
		}

		content, readErr := fs.ReadFile(assets, path)
		if readErr != nil {
			return readErr
		}
		pages++

		handler := bridgeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(content)
		}))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/"+path, nil))

		body := recorder.Body.String()
		if strings.Count(body, "/wails/runtime.js") != 1 || strings.Count(body, "/wails/ipc.js") != 1 {
			t.Fatalf("%s does not receive exactly one Wails runtime inject", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk pages: %v", err)
	}
	if pages == 0 {
		t.Fatal("no HTML pages found in embedded assets")
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
