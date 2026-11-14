//go:build native_test_seam && !promptcapture

package main

import "net/http"

func (s *nativeTestSeam) handlePromptCapture(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "prompt capture is not compiled", http.StatusNotImplemented)
}
