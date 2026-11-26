package main

import "testing"

func TestImageWorkflowPreservesExplicitVRAMBudget(t *testing.T) {
	five := 5
	opts := normalizeImageWorkflowOptions(ImageWorkflowOptions{MaxVRAM: &five})
	if opts.MaxVRAM == nil || *opts.MaxVRAM != 5 {
		t.Fatalf("workflow must forward explicit max VRAM 5, got %v", opts.MaxVRAM)
	}
}

func TestImageWorkflowOmitsIPStrengthWithoutAdapter(t *testing.T) {
	opts := normalizeImageWorkflowOptions(ImageWorkflowOptions{})
	if opts.IPAdapterStrength != 0 {
		t.Fatalf("plain workflow must not send IP-Adapter strength, got %v", opts.IPAdapterStrength)
	}
	withIP := normalizeImageWorkflowOptions(ImageWorkflowOptions{IPAdapterPath: "adapter.safetensors"})
	if withIP.IPAdapterStrength != 0.7 {
		t.Fatalf("IP workflow should default strength 0.7, got %v", withIP.IPAdapterStrength)
	}
}

func TestImageWorkflowTreatsZeroVRAMBudgetAsAuto(t *testing.T) {
	zero := 0
	opts := normalizeImageWorkflowOptions(ImageWorkflowOptions{MaxVRAM: &zero})
	if opts.MaxVRAM != nil {
		t.Fatalf("zero max VRAM means Auto to the panel, got %v", *opts.MaxVRAM)
	}
}
