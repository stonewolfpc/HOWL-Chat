package main

import (
	"context"
	stderrors "errors"

	"howl-chat/internal/audio"
	apperrors "howl-chat/internal/errors"
	"howl-chat/internal/media"
)

// captureProbe is the Media domain's capture port. It is a precondition check,
// not a capture capability: it reports whether Windows selects a default console
// capture endpoint that is present and enabled (DEVICE_STATE_ACTIVE). It does not
// prove the endpoint can be opened, that its mix format is accepted, or that
// recording and transcription will succeed. The renderer owns actual capture and
// reports those failures itself.
//
// It deliberately uses the presence-level probe. Activating a shared-mode
// capture client is a system-wide audio event: on a Bluetooth headset it makes
// Windows select HFP over A2DP, resample every other application's playback on
// that device to 8/16 kHz mono, and change the device's playback level until the
// client is released. A status read or a precondition check must never cause
// that; only the user's own recording may.
func (a *App) captureProbe() (interface{}, bool, string) {
	status := microphonePresenceNow()
	return status, status.Usable(), status.Message
}

// reportAudioFailure escalates a Media audio failure through the application's
// error system. The domain names the kind; the application owns the codes.
func (a *App) reportAudioFailure(kind string, err error) {
	switch kind {
	case "mic":
		a.logError(apperrors.ASRMicUnavailable, apperrors.High, "audio.apply.mic", "internal/media/audio.go:Apply", err)
	case "mic.abort":
		a.logError(apperrors.ASRMicUnavailable, apperrors.High, "audio.apply.mic.abort", "internal/media/audio.go:Apply", err)
	case "mic.lost":
		a.logError(apperrors.ASRMicUnavailable, apperrors.High, "audio.apply.mic.lost", "internal/media/audio.go:Apply", err)
	default:
		a.logError(apperrors.ASRModelLoadFail, apperrors.High, "audio.apply", "internal/media/audio.go:Apply", err)
	}
}

// applyAudioSettings is the explicit activation path. The Media domain owns the
// desired contract, the apply route and the activation intent; App only maps the
// domain's capture refusal onto the typed status the UI already understands.
func (a *App) applyAudioSettings(desired audio.DesiredSettings) error {
	if a.audio == nil {
		return stderrors.New("audio manager not initialized")
	}
	err := a.audio.Apply(context.Background(), desired)
	var captureErr *media.CaptureUnavailableError
	if stderrors.As(err, &captureErr) {
		if status, ok := captureErr.Status.(MicrophoneStatus); ok {
			return &MicrophoneUnavailableError{Status: status}
		}
		return &MicrophoneUnavailableError{Status: MicrophoneStatus{
			State:   MicStateUnavailable,
			Message: captureErr.Detail,
		}}
	}
	return err
}
