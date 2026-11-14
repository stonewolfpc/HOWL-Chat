//go:build !windows

package main

// probeMicrophone cannot verify capture availability on this platform. Returning
// unsupported fails closed: the ASR gate refuses activation rather than silently
// loading a model without a verified microphone.
func probeMicrophone() MicrophoneStatus {
	return unsupportedMicrophoneStatus()
}
