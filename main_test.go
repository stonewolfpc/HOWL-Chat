package main

import (
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/paths"
)

// The app resolves its isolation root once. Give the whole package an isolated
// state root before any test can touch a real config directory.
func TestMain(m *testing.M) {
	if ok, code := runFakeLlamaServerIfRequested(os.Args); ok {
		os.Exit(code)
	}
	root, err := os.MkdirTemp("", "howl-chat-tests-*")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(root, paths.IsolationMarkerName), []byte(paths.IsolationMarkerContent+"\n"), 0644); err != nil {
		panic(err)
	}
	if err := os.Setenv(paths.IsolationEnvVar, root); err != nil {
		panic(err)
	}
	code := m.Run()
	_ = os.RemoveAll(root)
	os.Exit(code)
}
