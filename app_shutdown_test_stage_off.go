//go:build !uiharness

package main

import "howl-chat/internal/host"

// testShutdownStages is empty outside the browser test host: the packaged app
// has no environment switch that can inject a failing shutdown stage.
func testShutdownStages() []host.Stage {
	return nil
}
