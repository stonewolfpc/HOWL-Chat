package main

import (
	"os"
	"path/filepath"
	"testing"

	imagetypes "howl-chat/internal/image/types"
)

// Switching from a Flux model to an SDXL checkpoint used to carry ae.safetensors
// into the new load, which then failed validation instead of using the built-in
// autoencoder.
func TestPruneStaleCompanionsDropsFluxVAEForSDXL(t *testing.T) {
	opts := imagetypes.ModelLoadOptions{
		VAEPath: filepath.Join("G:"+string(filepath.Separator), "models", "ae.safetensors"),
	}
	famFor := func(string) string { return "flux" }
	pruneStaleCompanions(&opts, "sdxl", famFor)
	if opts.VAEPath != "" {
		t.Errorf("flux VAE survived the switch to sdxl: %q", opts.VAEPath)
	}
}

func TestPruneStaleCompanionsKeepsTheRightVAE(t *testing.T) {
	p := filepath.Join("G:"+string(filepath.Separator), "models", "qwen_image_vae.safetensors")
	opts := imagetypes.ModelLoadOptions{VAEPath: p}
	pruneStaleCompanions(&opts, "qwen_image", func(string) string { return "qwen_image" })
	if opts.VAEPath != p {
		t.Errorf("matching VAE was dropped: %q", opts.VAEPath)
	}
}

// A file the registry never classified must not be discarded on a guess.
func TestPruneStaleCompanionsKeepsUnknownFiles(t *testing.T) {
	p := filepath.Join("G:"+string(filepath.Separator), "models", "custom_vae.safetensors")
	opts := imagetypes.ModelLoadOptions{VAEPath: p}
	pruneStaleCompanions(&opts, "sdxl", func(string) string { return "" })
	if opts.VAEPath != p {
		t.Errorf("unclassified file was dropped: %q", opts.VAEPath)
	}
}

// A registry path is only a candidate. This is the exact file that killed the
// flux run: sd-server aborted on GGML_ASSERT(weight_scale->type ==
// GGML_TYPE_F32) when handed a bitsandbytes int8 T5.
func TestUsableCompanionRejectsBitsAndBytesT5(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bnb_llm_int8_model.safetensors")
	if err := os.WriteFile(bad, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := usableCompanion("flux", "t5xxl", bad); got != "" {
		t.Errorf("usableCompanion accepted a bitsandbytes int8 encoder: %q", got)
	}
}

func TestUsableCompanionAcceptsTheRightFile(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "t5xxl_fp8_e4m3fn.safetensors")
	if err := os.WriteFile(good, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := usableCompanion("flux", "t5xxl", good); got != good {
		t.Errorf("usableCompanion(%q) = %q, want it accepted", good, got)
	}
}

func TestUsableCompanionRejectsMissingAndEmpty(t *testing.T) {
	if got := usableCompanion("flux", "t5xxl", ""); got != "" {
		t.Errorf("empty path accepted as %q", got)
	}
	missing := filepath.Join(t.TempDir(), "t5xxl_gone.safetensors")
	if got := usableCompanion("flux", "t5xxl", missing); got != "" {
		t.Errorf("deleted file accepted as %q", got)
	}
}
