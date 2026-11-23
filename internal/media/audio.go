// Package media owns the application's media state and the operations that
// mutate it: audio configuration and activation, and the image engine session.
//
// Level 0 (internal/host) decides when work may start and stop; this package
// owns what the media subsystems know. App keeps the Wails bindings and the
// presentation projections, and holds no shadow copy of this state.
package media

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"howl-chat/internal/audio"
)

// CaptureUnavailableError reports that an ASR activation was refused because no
// usable capture device existed. Status carries the probe result as the caller
// found it, so the caller can return its own typed error without re-probing a
// device that may have changed since.
type CaptureUnavailableError struct {
	Detail string
	Status interface{}
}

func (e *CaptureUnavailableError) Error() string { return e.Detail }

// Capture is the probe port. It reports the current capture state; it must not
// open a stream. Returning the caller's own status type keeps the microphone
// presentation owned by the caller.
type Capture func() (status interface{}, usable bool, detail string)

// Audio owns the managed audio manager, the desired settings contract, the last
// explicit apply failure, and the ASR activation intent.
//
// It is the only writer of those four things. Status projection is derived by
// the caller from these readings.
type Audio struct {
	mu         sync.RWMutex
	manager    *audio.Manager
	desired    audio.DesiredSettings
	desiredSet bool
	applyErr   string

	asrPending atomic.Bool

	admit   func(operation string) error
	capture Capture
	report  func(kind string, err error)
}

// NewAudio builds the audio owner. admit refuses new work during shutdown;
// capture probes a device without opening it; report escalates a failure.
func NewAudio(manager *audio.Manager, admit func(string) error, capture Capture, report func(kind string, err error)) *Audio {
	return &Audio{manager: manager, admit: admit, capture: capture, report: report}
}

// Manager returns the managed pipeline surface. It may be nil when the manager
// could not be constructed.
func (a *Audio) Manager() *audio.Manager {
	if a == nil {
		return nil
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.manager
}

// SetManager replaces the manager. Only construction and tests use it.
func (a *Audio) SetManager(manager *audio.Manager) {
	if a == nil {
		return
	}
	a.mu.Lock()
	a.manager = manager
	a.mu.Unlock()
}

// SetDesired records the full requested contract.
func (a *Audio) SetDesired(desired audio.DesiredSettings) {
	if a == nil {
		return
	}
	a.mu.Lock()
	a.desired = desired
	a.desiredSet = true
	a.mu.Unlock()
}

// Desired returns the requested contract, defaulting when nothing was recorded.
func (a *Audio) Desired() audio.DesiredSettings {
	if a == nil {
		return audio.DefaultDesiredSettings()
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	if !a.desiredSet {
		return audio.DefaultDesiredSettings()
	}
	return a.desired
}

// SetApplyError records the last explicit apply failure. An empty string clears
// it; a failure belongs to the attempt that produced it.
func (a *Audio) SetApplyError(message string) {
	if a == nil {
		return
	}
	a.mu.Lock()
	a.applyErr = message
	a.mu.Unlock()
}

// LastApplyError returns the last explicit apply failure, if any.
func (a *Audio) LastApplyError() string {
	if a == nil {
		return ""
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.applyErr
}

// ASRPending reports that ASR was explicitly requested and not yet confirmed.
func (a *Audio) ASRPending() bool { return a != nil && a.asrPending.Load() }

// Apply is the explicit activation path. The requested full contract is stored
// before validation so status can explain exactly which values were rejected.
// Only settings.Backend() is eligible for adoption by the manager; unsupported
// fields never become applied merely because they were saved.
//
// ASR activation is gated on capture availability. The probe runs before the
// manager builds a generation, so a missing or unusable device cannot allocate a
// model or reserve memory, and it runs again after the build so a device that
// disappears mid-initialization is torn down instead of left loaded.
func (a *Audio) Apply(ctx context.Context, desired audio.DesiredSettings) error {
	if a.admit != nil {
		if err := a.admit("audio settings"); err != nil {
			return err
		}
	}
	a.SetDesired(desired)

	issues := desired.Validate()
	if audio.HasBackendValidationIssue(desired, issues) {
		err := &audio.SettingsValidationError{Issues: issues}
		a.SetApplyError(err.Error())
		return err
	}
	manager := a.Manager()
	if manager == nil {
		err := fmt.Errorf("audio manager not initialized")
		a.SetApplyError(err.Error())
		return err
	}

	// A previous failure belongs to the previous attempt. Clear it only after
	// the new request has been recorded and the manager is available.
	a.SetApplyError("")
	settings := desired.Backend()

	if !settings.ASREnabled {
		a.asrPending.Store(false)
		_, _, err := manager.Apply(context.Background(), settings)
		if err != nil {
			a.SetApplyError(err.Error())
			a.escalate("apply", fmt.Errorf("apply ASR engine %q: %w", settings.ASREngine, err))
		}
		return err
	}

	// From here on the user has explicitly asked for ASR. Keep the request
	// visible to status until a generation is applied and still backed by a
	// usable capture device.
	a.asrPending.Store(true)

	status, usable, detail := a.probe()
	if !usable {
		err := &CaptureUnavailableError{Detail: detail, Status: status}
		a.SetApplyError(err.Error())
		a.escalate("mic", err)
		return err
	}

	if _, _, err := manager.Apply(context.Background(), settings); err != nil {
		a.SetApplyError(err.Error())
		a.escalate("apply", fmt.Errorf("apply ASR engine %q: %w", settings.ASREngine, err))
		return err
	}

	afterStatus, afterUsable, afterDetail := a.probe()
	if !afterUsable {
		// The device disappeared while the model was loading. Retire the ASR
		// capability in the same manager so no generation keeps a loaded model
		// with no capture path. TTS, if requested, is preserved.
		ttsOnly := settings
		ttsOnly.ASREnabled = false
		if _, _, abortErr := manager.Apply(context.Background(), ttsOnly); abortErr != nil {
			a.escalate("mic.abort", abortErr)
		}
		a.asrPending.Store(true)
		err := &CaptureUnavailableError{Detail: afterDetail, Status: afterStatus}
		a.SetApplyError(err.Error())
		a.escalate("mic.lost", err)
		return err
	}
	a.asrPending.Store(false)
	a.SetApplyError("")
	return nil
}

// Close releases the manager and everything it owns.
func (a *Audio) Close(ctx context.Context) error {
	manager := a.Manager()
	if manager == nil {
		return nil
	}
	return manager.Close(ctx)
}

func (a *Audio) probe() (interface{}, bool, string) {
	if a.capture == nil {
		return nil, true, ""
	}
	return a.capture()
}

func (a *Audio) escalate(kind string, err error) {
	if a.report != nil {
		a.report(kind, err)
	}
}
