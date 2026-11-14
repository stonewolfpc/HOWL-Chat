package main

import "time"

// MicrophoneState is one capture state this build can prove. A state is only
// reported after a probe; enumeration alone never produces MicStateReady.
type MicrophoneState string

const (
	// MicStateReady means the default capture endpoint was initialized in WASAPI
	// shared mode with its mix format. It does not mean recording started, audio
	// samples were received, or transcription succeeded.
	MicStateReady MicrophoneState = "ready"
	// MicStateNoDevice means no default capture endpoint was found.
	MicStateNoDevice MicrophoneState = "no-device"
	// MicStateUnavailable means the default endpoint exists but is inactive,
	// disconnected, in use, or otherwise unavailable.
	MicStateUnavailable MicrophoneState = "unavailable"
	// MicStateAccessDenied is reported only when the API returns access denied.
	MicStateAccessDenied MicrophoneState = "access-denied"
	// MicStateCaptureInitFailed means the endpoint exists but WASAPI
	// initialization failed for a reason other than the states above.
	MicStateCaptureInitFailed MicrophoneState = "capture-init-failed"
	// MicStateUnsupported means this platform has no implemented probe. It fails
	// closed and never claims capture readiness.
	MicStateUnsupported MicrophoneState = "unsupported"
)

// MicrophoneNoDeviceMessage is the exact user-facing copy for an install with no
// selected microphone. It is the contract the audit requires; do not reword it.
const MicrophoneNoDeviceMessage = "Windows shows a mic isn't selected. Please connect a microphone to use this option"

// MicrophoneStatus is the truthful capture state plus the copy the UI shows.
type MicrophoneStatus struct {
	State   MicrophoneState `json:"state"`
	Message string          `json:"message"`
	Detail  string          `json:"detail,omitempty"`
	// DeviceCount is 0 or 1 for this probe: whether the role's default endpoint
	// was found, not a count of every capture device Windows enumerates.
	DeviceCount int    `json:"device_count"`
	CheckedAt   string `json:"checked_at"`
	// Role is the Windows endpoint role that was probed (console for the
	// implemented probe). It documents which Windows selection the status
	// describes; console and communications defaults can differ.
	Role string `json:"role,omitempty"`
}

// Usable reports whether ASR may allocate a model. Ready here means the default
// endpoint initialized; browser capture and transcription remain unverified.
func (s MicrophoneStatus) Usable() bool { return s.State == MicStateReady }

// microphoneProbe is the platform seam. Tests replace it to exercise every
// state without physical hardware.
var microphoneProbe = probeMicrophone

// microphoneStatusNow probes capture availability now. Probes are cheap and
// side-effect free (initialize, then release), so every apply and status read
// can use a fresh result instead of a startup cache that may be stale.
func microphoneStatusNow() MicrophoneStatus {
	status := microphoneProbe()
	if status.State == "" {
		status.State = MicStateCaptureInitFailed
	}
	if status.State == MicStateNoDevice && status.Message == "" {
		status.Message = MicrophoneNoDeviceMessage
	}
	if status.Message == "" {
		status.Message = microphoneMessageFor(status.State)
	}
	if status.CheckedAt == "" {
		status.CheckedAt = time.Now().UTC().Format(time.RFC3339)
	}
	return status
}

func microphoneMessageFor(state MicrophoneState) string {
	switch state {
	case MicStateReady:
		return "The Windows default microphone endpoint initialized for shared-mode capture. Recording has not started."
	case MicStateNoDevice:
		return MicrophoneNoDeviceMessage
	case MicStateUnavailable:
		return "The Windows default microphone is unavailable (disabled, disconnected, or in use)."
	case MicStateAccessDenied:
		return "Windows is blocking microphone access. Allow desktop apps to use the microphone in Windows Settings."
	case MicStateCaptureInitFailed:
		return "The Windows default microphone is present but could not be initialized for recording."
	case MicStateUnsupported:
		return "Microphone availability cannot be verified on this platform, so speech recognition is unavailable."
	default:
		return "Microphone availability could not be established."
	}
}

func microphoneStatusRecord(status MicrophoneStatus) map[string]interface{} {
	return map[string]interface{}{
		"state":        string(status.State),
		"message":      status.Message,
		"detail":       status.Detail,
		"device_count": status.DeviceCount,
		"checked_at":   status.CheckedAt,
		"role":         status.Role,
	}
}

// unsupportedMicrophoneStatus is shared by the non-Windows stub so its
// fail-closed behavior is testable on every platform.
func unsupportedMicrophoneStatus() MicrophoneStatus {
	return newMicrophoneStatus(
		MicStateUnsupported,
		"WASAPI microphone probing is only implemented on Windows; capture availability is unverified",
		0,
	)
}

func newMicrophoneStatus(state MicrophoneState, detail string, deviceCount int) MicrophoneStatus {
	return MicrophoneStatus{
		State:       state,
		Message:     microphoneMessageFor(state),
		Detail:      detail,
		DeviceCount: deviceCount,
		CheckedAt:   time.Now().UTC().Format(time.RFC3339),
	}
}

// MicrophoneUnavailableError carries the state that refused activation. Its
// message is the exact copy the UI must show, so a failed Apply never degrades
// into a generic backend error.
type MicrophoneUnavailableError struct {
	Status MicrophoneStatus
}

func (e *MicrophoneUnavailableError) Error() string {
	if e == nil {
		return "microphone unavailable"
	}
	if e.Status.Message != "" {
		return e.Status.Message
	}
	return "microphone unavailable"
}

func (e *MicrophoneUnavailableError) State() MicrophoneState {
	if e == nil {
		return MicStateCaptureInitFailed
	}
	return e.Status.State
}
