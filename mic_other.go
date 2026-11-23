//go:build !windows

package main

// probeMicrophone cannot verify capture availability on this platform. Returning
// unsupported fails closed: the ASR gate refuses activation rather than silently
// loading a model without a verified microphone.
func probeMicrophone() MicrophoneStatus {
	return unsupportedMicrophoneStatus()
}

// probeCapturePresence is the presence-level counterpart and fails closed for the
// same reason: this platform has no implemented way to prove the capture
// endpoint exists, so it never reports ready.
func probeCapturePresence() MicrophoneStatus {
	return unsupportedMicrophoneStatus()
}
