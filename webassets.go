package main

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
)

// Wails only injects its runtime into the "main page": the AssetServer treats a
// request as index-like when its path ends in "/" or "/index.html". Every other
// page that is navigated to directly (lorebooks.html, characters.html, ...) is
// served as a plain file, so it runs without window.go and without window.runtime.
// Pages then silently fall back to localStorage instead of the Go backend.
//
// This middleware applies the same two tags Wails uses to every HTML response,
// so one seam covers all current and future pages. It is not a second bridge:
// the runtime and IPC scripts are exactly the ones the AssetServer serves at
// /wails/runtime.js and /wails/ipc.js.
const wailsRuntimeTags = `<script src="/wails/runtime.js"></script><script src="/wails/ipc.js"></script>`

// bridgeMiddleware injects the Wails runtime into HTML responses that the
// AssetServer did not already process.
func bridgeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		recorder := &bridgeResponse{ResponseWriter: w}
		next.ServeHTTP(recorder, r)
		recorder.finish(true)
	})
}

// bridgeResponse buffers successful HTML responses so the runtime tags can be
// inserted, and streams everything else straight through.
type bridgeResponse struct {
	http.ResponseWriter
	status      int
	buffer      bytes.Buffer
	buffered    bool
	wroteHeader bool
}

func (b *bridgeResponse) WriteHeader(code int) {
	if b.wroteHeader {
		return
	}
	b.wroteHeader = true
	b.status = code

	isHTML := strings.Contains(strings.ToLower(b.Header().Get("Content-Type")), "text/html")
	if code == http.StatusOK && isHTML {
		b.buffered = true
		return
	}
	b.ResponseWriter.WriteHeader(code)
}

func (b *bridgeResponse) Write(p []byte) (int, error) {
	if !b.wroteHeader {
		b.WriteHeader(http.StatusOK)
	}
	if !b.buffered {
		return b.ResponseWriter.Write(p)
	}
	return b.buffer.Write(p)
}

// Unwrap lets http.ResponseController reach the underlying writer.
func (b *bridgeResponse) Unwrap() http.ResponseWriter { return b.ResponseWriter }

// Flush commits whatever has been buffered through the injection first, then
// flushes. Flushing the raw writer first would send the unmodified prefix and
// make the later injection impossible.
//
// Two cases defer instead of committing: nothing has been written yet, and the
// buffer holds no body (a handler flushed before the head arrived). Committing
// either would send headers before the document shape is known, and headers
// cannot be revised afterwards, so injection would be impossible.
func (b *bridgeResponse) Flush() {
	if !b.wroteHeader {
		return
	}
	if b.buffered {
		if b.buffer.Len() == 0 {
			return
		}
		b.finish(false)
	}
	if flusher, ok := b.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// finish writes the buffered body through the injection. final is false when a
// Flush forces an early commit: the document is still incomplete, so the fixed
// Content-Length is removed *before* WriteHeader and is never set again, since
// headers already sent cannot be revised. Only a final finish, with the whole
// body in hand, sets an exact length.
func (b *bridgeResponse) finish(final bool) {
	if !b.buffered {
		return
	}
	b.buffered = false
	body := append([]byte(nil), b.buffer.Bytes()...)
	b.buffer.Reset()

	// A compressed body cannot be rewritten by string insertion. Wails does not
	// compress today; this keeps a future encoder from silently corrupting pages.
	encoded := strings.TrimSpace(b.Header().Get("Content-Encoding"))
	if encoded == "" || strings.EqualFold(encoded, "identity") {
		if !bytes.Contains(body, []byte("/wails/runtime.js")) {
			body = injectRuntimeTags(body)
		}
	}

	header := b.Header()
	if final {
		header.Set("Content-Length", strconv.Itoa(len(body)))
	} else {
		header.Del("Content-Length")
	}
	b.ResponseWriter.WriteHeader(b.status)
	_, _ = b.ResponseWriter.Write(body)
}

// injectRuntimeTags puts the tags at the start of <head>, matching how Wails
// inserts them for the main page (as the first child), so both page kinds run
// the runtime before any other head script. Documents without a head are
// returned unchanged.
func injectRuntimeTags(body []byte) []byte {
	lower := bytes.ToLower(body)
	open := bytes.Index(lower, []byte("<head"))
	if open < 0 {
		return body
	}
	close := bytes.IndexByte(lower[open:], '>')
	if close < 0 {
		return body
	}
	at := open + close + 1

	out := make([]byte, 0, len(body)+len(wailsRuntimeTags))
	out = append(out, body[:at]...)
	out = append(out, wailsRuntimeTags...)
	out = append(out, body[at:]...)
	return out
}
