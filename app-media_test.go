package main

import (
	"howl-chat/internal/audio"
	"howl-chat/internal/image/gallery"
	"howl-chat/internal/media"
)

// newMediaTestApp builds an App whose Media owners exist but hold nothing.
// Tests install the manager or gallery they intend to exercise.
func newMediaTestApp() *App {
	app := &App{}
	app.audio = media.NewAudio(nil,
		func(operation string) error { return app.level0().Gate.Ensure(operation) },
		app.captureProbe,
		app.reportAudioFailure)
	app.images = media.NewImage("")
	return app
}

// newAudioTestAppWithManager builds an App whose audio owner holds manager and
// whose ports are the application's real ones, so the microphone seams the
// tests replace are the ones the domain actually calls.
func newAudioTestAppWithManager(manager *audio.Manager) *App {
	app := newMediaTestApp()
	app.audio.SetManager(manager)
	return app
}

// newGalleryTestApp builds an App whose image owner holds one gallery.
func newGalleryTestApp(service *gallery.Service) *App {
	app := newMediaTestApp()
	app.images.SetGallery(service)
	return app
}
