package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"howl-chat/internal/audio/asr/whisper"
	"howl-chat/internal/audio/types"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: test-asr <audio-file>")
		fmt.Println("Example: test-asr test.wav")
		os.Exit(1)
	}

	audioPath := os.Args[1]

	// Check file exists
	if _, err := os.Stat(audioPath); err != nil {
		fmt.Printf("Error: Audio file not found: %s\n", audioPath)
		os.Exit(1)
	}

	fmt.Println("=== HOWL Chat ASR Test ===")
	fmt.Printf("Audio file: %s\n\n", audioPath)

	// Create recognizer
	recognizer := whisper.New(types.ASRWhisperBase)

	// Initialize with model path
	ctx := context.Background()
	config := types.ASRConfig{
		ModelPath: "models\\ggml-base.bin",
		Language:  "auto",
		BeamSize:  5,
		BestOf:    5,
	}

	fmt.Println("Initializing recognizer...")
	start := time.Now()
	if err := recognizer.Initialize(ctx, config); err != nil {
		fmt.Printf("Error initializing: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Initialized in %.1fs\n\n", time.Since(start).Seconds())

	// Transcribe with progress callback
	fmt.Println("Transcribing...")
	start = time.Now()
	result, err := recognizer.Transcribe(ctx, audioPath, func(progress int, status string) {
		fmt.Printf("  [%3d%%] %s\n", progress, status)
	})
	if err != nil {
		fmt.Printf("Error transcribing: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n=== Results ===\n")
	fmt.Printf("Duration:  %.1fs\n", result.Duration)
	fmt.Printf("Language:  %s\n", result.Language)
	fmt.Printf("Confidence: %.2f\n", result.Confidence)
	fmt.Printf("Text:      %s\n", result.Text)

	if len(result.Segments) > 0 {
		fmt.Printf("\nSegments:\n")
		for _, seg := range result.Segments {
			fmt.Printf("  [%5.1f - %5.1f] %.2f | %s\n", seg.StartTime, seg.EndTime, seg.Confidence, seg.Text)
		}
	}

	// Cleanup
	recognizer.Release()
	fmt.Println("\nDone!")
}
