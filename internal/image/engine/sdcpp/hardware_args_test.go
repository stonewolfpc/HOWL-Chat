package sdcpp

import (
	"strings"
	"testing"

	"howl-chat/internal/image/types"
)

func intPtr(v int) *int { return &v }

func TestFlashAttentionModeIsTriState(t *testing.T) {
	e := &Engine{config: types.EngineConfig{Port: 8081}}
	m := &types.ModelInfo{Family: "sd1.5", Files: map[string]string{"model": `E:\m.safetensors`}}
	for _, mode := range []string{"auto", "on", "off"} {
		args := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{FlashAttentionMode: mode}), " ")
		if !strings.Contains(args, "--flash-attn "+mode) {
			t.Errorf("mode %q did not reach the engine: %s", mode, args)
		}
	}
	legacy := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{FlashAttention: true}), " ")
	if !strings.Contains(legacy, "--flash-attn on") {
		t.Errorf("legacy true must map to on: %s", legacy)
	}
	defaulted := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{}), " ")
	if !strings.Contains(defaulted, "--flash-attn auto") {
		t.Errorf("no selection must explicitly send auto: %s", defaulted)
	}
}

func TestTTLAndDiskBackendReachTheEngine(t *testing.T) {
	e := &Engine{config: types.EngineConfig{Port: 8081}}
	m := &types.ModelInfo{Family: "sd1.5", Files: map[string]string{"model": `E:\m.safetensors`}}
	args := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{
		ParamsBackend:      "disk",
		WeightResidencyTTL: intPtr(0),
		VAETiling:          true,
	}), " ")
	for _, want := range []string{
		"--params-backend disk",
		"--weight-residency-ttl 0",
		"--vae-tiling",
	} {
		if !strings.Contains(args, want) {
			t.Errorf("%q missing: %s", want, args)
		}
	}
}

func TestNegativeAndZeroTTLAreDistinct(t *testing.T) {
	e := &Engine{config: types.EngineConfig{Port: 8081}}
	m := &types.ModelInfo{Family: "sd1.5", Files: map[string]string{"model": `E:\m.safetensors`}}
	args := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{
		WeightResidencyTTL: intPtr(-1),
	}), " ")
	if !strings.Contains(args, "--weight-residency-ttl -1") {
		t.Errorf("a negative TTL must be forwarded unchanged: %s", args)
	}
	if strings.Contains(args, "--weight-staging-ttl") {
		t.Errorf("an unset TTL must remain unset: %s", args)
	}
}
