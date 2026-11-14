package main

import (
	"howl-chat/internal/paths"
)

// Isolated native test support.
//
// HOWL_WEBVIEW_DATA_DIR is the single opt-in switch (see internal/paths). These
// helpers keep the condition in one place: when isolation is off every function
// returns "" so the historical cwd/exe-relative resolution stays in charge and
// production behaviour is unchanged.

// usingIsolation reports whether this process is an isolated test run.
func usingIsolation() bool { return paths.IsolationActive() }

// isolatedStateDir returns <state root>/<name> when isolated, "" otherwise.
func isolatedStateDir(name string) string {
	if !paths.IsolationActive() {
		return ""
	}
	return paths.StatePath(name)
}

// isolatedStatePath joins elements onto the state root when isolated, "" otherwise.
func isolatedStatePath(elem ...string) string {
	if !paths.IsolationActive() {
		return ""
	}
	return paths.StatePath(elem...)
}

// isolatedTempDir returns the isolated scratch directory when isolated, and the
// normal OS temp directory otherwise.
func isolatedTempDir() string { return paths.TempDir() }
