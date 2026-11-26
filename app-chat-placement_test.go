package main

import (
	"strings"
	"testing"

	llamabackend "howl-chat/internal/backend/llama"
	"howl-chat/internal/resource"
)

func placementDecision(vram, ram int64) resource.AdmissionDecision {
	return resource.AdmissionDecision{Estimate: resource.AdmissionEstimate{
		RequiredVRAMMB: vram,
		RequiredRAMMB:  ram,
	}}
}

func TestChatPlacementReportsPlannedAndActualSeparately(t *testing.T) {
	settings := map[string]interface{}{
		"gpu_layers":   0,
		"context_size": 3096,
		"batch_size":   512,
	}
	actual := llamabackend.Placement{
		Known: true, GPULayers: 0, TotalLayers: 15, Context: 3328, Batch: 512, RAMMB: 140,
		KVRAMMB: 40, ComputeRAMMB: 10,
		Buffers: map[string]float64{"CPU_Mapped": 140},
	}
	report := chatPlacementReport("m.gguf", settings, placementDecision(0, 865), actual)
	if !report.MatchesPlan {
		t.Fatalf("llama's 256-token context rounding must not be reported as a different plan: %s", report.Message)
	}
	for _, want := range []string{"as planned", "0/15", "context 3328, batch 512", "rounded", "requested context 3096", "0 MiB VRAM / 140 MiB RAM", "engine total 0 MiB VRAM / 190 MiB RAM"} {
		if !strings.Contains(report.Message, want) {
			t.Fatalf("message %q does not contain %q", report.Message, want)
		}
	}
	if strings.Contains(report.Message, "different placement") {
		t.Fatalf("a rounded context must not be called a different placement: %s", report.Message)
	}
}

func TestChatPlacementReportsLowerContextAsMismatch(t *testing.T) {
	settings := map[string]interface{}{"gpu_layers": 10, "context_size": 8192, "batch_size": 512}
	actual := llamabackend.Placement{
		Known: true, GPULayers: 10, TotalLayers: 33, Context: 4096, GPUMB: 1200,
		Buffers: map[string]float64{"CUDA0": 1200},
	}
	report := chatPlacementReport("m.gguf", settings, placementDecision(1200, 200), actual)
	if report.MatchesPlan {
		t.Fatal("a lower actual context must not be reported as the plan")
	}
	if !strings.Contains(report.Message, "different placement") || !strings.Contains(report.Message, "context 4096") {
		t.Fatalf("mismatch message lost the observed context: %s", report.Message)
	}
}

func TestChatPlacementDoesNotInventVerification(t *testing.T) {
	settings := map[string]interface{}{"gpu_layers": 4, "context_size": 4096, "batch_size": 512}
	report := chatPlacementReport("m.gguf", settings, placementDecision(100, 200), llamabackend.Placement{})
	if report.MatchesPlan {
		t.Fatal("unknown placement must not match the plan")
	}
	if !strings.Contains(report.Message, "unverified") {
		t.Fatalf("unknown placement must be called unverified: %s", report.Message)
	}
}
