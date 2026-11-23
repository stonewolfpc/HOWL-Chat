package main

import (
	"context"
	"fmt"
	"time"

	apperrors "howl-chat/internal/errors"
	"howl-chat/internal/host"
	"howl-chat/internal/proc"
)

// llamaShutdownTimeout bounds the wait for a killed llama-server. The process
// is killed first; this only prevents shutdown from blocking forever on a
// driver or child-process failure.
const llamaShutdownTimeout = 10 * time.Second

// level0 returns the application's lifecycle and coordination authority.
//
// The authority owns admission, the ordered boot and shutdown sequences,
// background joining, the active stream token and child-process supervision. It
// deliberately stores no conversation, model, audio, image or configuration
// value: those stay with their own owners. See
// docs/architecture/HOC-H0-LEVEL0-2026-11-22.md for the state that still lives
// in App and where it is meant to go.
func (a *App) level0() *host.Authority {
	a.authorityOnce.Do(func() { a.authority = host.NewAuthority() })
	return a.authority
}

// shutdown runs the one authoritative cleanup sequence through Level 0. It is
// safe to call from both the Wails shutdown callback and an explicit Exit
// action: the lifecycle's gate runs the sequence once and shares its result.
func (a *App) shutdown(ctx context.Context) error {
	err := a.level0().Lifecycle.Shutdown(ctx, a.shutdownStages(), host.ShutdownOptions{
		// Error reporting belongs to the application; the lifecycle owns only
		// ordering and accumulation, so every refusal names the failing stage.
		OnFailure: func(stage string, err error, panicked bool) {
			severity := apperrors.High
			if panicked {
				severity = apperrors.Critical
			}
			a.logError(apperrors.SYSShutdownFail, severity, "app.shutdown", "app-lifecycle.go:shutdownStages", err)
		},
	})
	// Diagnostics owns its log file's lifecycle; the sequence has finished, so
	// every shutdown failure above has already been recorded.
	if a.diagnostics != nil {
		_ = a.diagnostics.Close()
	}
	return err
}

// shutdownStages is the application's ordered teardown: stop admitting and
// issuing work, join what is still running, persist durable state, then release
// owned resources in dependency order. Level 0 runs it so one broken component
// cannot prevent the rest of the cleanup from running.
func (a *App) shutdownStages() []host.Stage {
	stages := []host.Stage{
		{Name: "stop in-flight work", Run: a.stopInFlightWork},
		{Name: "background work", Run: func(context.Context) error {
			return a.level0().Background.Join(10 * time.Second)
		}},
		{Name: "final chat save", Run: a.finalChatSave},
		{Name: "llama-server stop", Run: func(context.Context) error {
			return a.level0().Processes.Stop(llamaShutdownTimeout)
		}},
		{Name: "image engine shutdown", Run: func(context.Context) error { return a.images.Close() }},
		{Name: "audio manager shutdown", Run: a.shutdownAudio},
		{Name: "process job close", Run: func(context.Context) error {
			// All persistent children are assigned at start, so
			// KILL_ON_JOB_CLOSE also covers a crash before this point; this
			// closes the handle deliberately on normal shutdown.
			proc.CloseJob()
			return nil
		}},
	}
	// The browser test host can append one genuinely failing stage to prove
	// shutdown escalation end to end; every other build compiles a no-op stub.
	return append(stages, testShutdownStages()...)
}

// stopInFlightWork cancels whatever this process is currently doing. The HTTP
// client's request cancellation is separate from the stream-owned context: both
// are needed for a blocking request and for a stream that has not reached its
// callback.
func (a *App) stopInFlightWork(context.Context) error {
	level0 := a.level0()
	level0.Streams.Cancel()
	if a.engine.Client() != nil {
		a.engine.Client().CancelRequest()
	}
	if jobs := a.images.Jobs(); jobs != nil {
		jobs.CancelAll()
	}
	stopNativeTestSeam()
	return nil
}

// finalChatSave persists the conversation before durable stores are finalised.
func (a *App) finalChatSave(context.Context) error {
	if !a.conversation.Bound() {
		return nil
	}
	return a.performFinalSave()
}

// shutdownAudio closes the managed audio pipeline, keeping the existing
// user-visible diagnostic expected by the audio ownership test and the desktop
// console.
func (a *App) shutdownAudio(ctx context.Context) error {
	if a.audio.Manager() == nil {
		return nil
	}
	err := a.audio.Close(ctx)
	if err != nil {
		fmt.Printf("Audio manager shutdown error: %v\n", err)
	}
	return err
}

// ensureAcceptingWork asks the gate whether this operation may start. It is the
// single refusal route every domain uses, so a shutdown refusal is always
// attributable to the operation that asked.
func (a *App) ensureAcceptingWork(operation string) error {
	return a.level0().Gate.Ensure(operation)
}
