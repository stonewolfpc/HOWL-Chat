package main

import (
	"strings"
	"testing"
)

// The confirm dialog must name the resident that actually holds the GPU. A CPU
// speech model is not the reason a 7 GB VRAM request cannot start.
func TestFitMessageBlamesTheGPUHolderNotACPUResident(t *testing.T) {
	message := fitMessage(ModelFitResult{
		CanFit:      false,
		RequiredMB:  6945,
		FreeVRAMMB:  7090,
		TotalVRAMMB: 8192,
		Blocking: []LoadedModelInfo{
			{Type: "tts", Name: "piper", GPUMemoryMB: 0, SystemRAMMB: 360},
		},
	})
	if strings.Contains(message, "piper") {
		t.Fatalf("a RAM-only resident must not be named as the VRAM blocker: %s", message)
	}
	if !strings.Contains(message, "cannot start") {
		t.Fatalf("without a GPU blocker the message must state the shortfall: %s", message)
	}
}

func TestFitMessageNamesEachWorkloadKind(t *testing.T) {
	cases := []struct {
		kind string
		name string
		want string
	}{
		{"llm", "qwen", "chat model"},
		{"image", "illustrious", "image model"},
		{"asr", "qwen3-asr", "speech recognition model"},
		{"tts", "piper", "speech model"},
	}
	for _, tc := range cases {
		message := fitMessage(ModelFitResult{
			RequiredMB:  6945,
			FreeVRAMMB:  1860,
			TotalVRAMMB: 8192,
			Blocking: []LoadedModelInfo{
				{Type: tc.kind, Name: tc.name, GPUMemoryMB: 4096},
			},
		})
		if !strings.Contains(message, tc.want) {
			t.Fatalf("%s resident labelled wrong: %s", tc.kind, message)
		}
		if !strings.Contains(message, tc.name) {
			t.Fatalf("message must name the model: %s", message)
		}
		if !strings.Contains(message, "Unload it and continue?") {
			t.Fatalf("message must offer the handoff: %s", message)
		}
	}
}
