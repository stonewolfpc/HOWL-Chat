//go:build uiharness

package main

import (
	"context"
	stderrors "errors"
	"os"

	"howl-chat/internal/host"
)

// shutdownTestSeamEnv is the browser host's explicit failure switch. It exists
// so the harness can drive a genuine stage failure through Level 0 — the stage
// runs in the same ordered sequence as every real stage, its error reaches
// host.Stages and the application-supplied sink, and the diagnostics owner
// records and emits it exactly as it would a real failure. The environment
// variable is checked here, in a uiharness-only file; the packaged binary
// compiles the no-op stub instead, so there is no runtime switch in production.
const shutdownTestSeamEnv = "HOWL_TEST_FAIL_SHUTDOWN_STAGE"

func testShutdownStages() []host.Stage {
	if os.Getenv(shutdownTestSeamEnv) != "1" {
		return nil
	}
	return []host.Stage{{
		Name: "shutdown failure seam",
		Run: func(context.Context) error {
			return stderrors.New("injected shutdown failure (" + shutdownTestSeamEnv + "=1)")
		},
	}}
}
