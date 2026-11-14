//go:build native_test_seam && promptcapture

package main

import (
	"encoding/json"
	"net/http"
)

// handlePromptCapture returns the messages seen at the production inference
// boundary by the test-only promptcapture recorder. Both this route and the
// recorder are compile-time absent from production builds.
func (s *nativeTestSeam) handlePromptCapture(w http.ResponseWriter, r *http.Request) {
	if !s.begin(w, r, http.MethodGet) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.app.GetPromptCapture())
}
