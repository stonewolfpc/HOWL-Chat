package main

import (
	"os"

	"howl-chat/internal/audio/tts/fakehelper"
)

func main() {
	os.Exit(fakehelper.Run(
		os.Stdin,
		os.Stdout,
		os.Stderr,
		os.Getenv("HOWL_FAKE_HELPER_MODE"),
		os.Getenv("HOWL_FAKE_HELPER_MARKER_DIR"),
	))
}
