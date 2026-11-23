package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	imagetypes "howl-chat/internal/image/types"
)

// The picker and the validator must agree, or a model the UI reports as fully
// equipped is refused at load time with a message about the file it just chose.
func TestCompanionFamilyUsable(t *testing.T) {
	cases := []struct {
		target, file string
		want         bool
		why          string
	}{
		{"anima", "anima", true, "same architecture"},
		{"anima", "", true, "unclassified file cannot be ruled out"},
		{"anima", "unknown", true, "scanner gave up, not a mismatch"},
		{"anima", "sd1.5", false, "an sd1.5 VAE decodes noise under Anima"},
		{"anima", "sdxl", false, "same reason"},
		{"z_image", "flux", true, "both decode with ae.safetensors"},
		{"flux", "z_image", true, "sharing is symmetric"},
		{"qwen_image", "flux", false, "qwen_image has its own VAE"},
		{"z_image", "qwen_image", false, "and does not borrow it"},
		{"Flux", "FLUX", true, "family names are case-insensitive"},
	}
	for _, c := range cases {
		if got := companionFamilyUsable(c.target, c.file); got != c.want {
			t.Errorf("companionFamilyUsable(%q, %q) = %v, want %v (%s)", c.target, c.file, got, c.want, c.why)
		}
	}
}

// A full SDXL checkpoint carries its VAE and both text encoders. The companion
// autofill used to attach a separate CLIP-L and VAE anyway, so sd-server loaded
// the file's UNet without its CLIP-G and rendered noise.
func TestClearEmbeddedCompanionsForFullCheckpoint(t *testing.T) {
	dir := t.TempDir()
	full := filepath.Join(dir, "full.safetensors")
	header := `{"first_stage_model.decoder.conv_in.weight":{"dtype":"F16","shape":[1],"data_offsets":[0,1]},` +
		`"conditioner.embedders.0.model.ln_final.weight":{"dtype":"F16","shape":[1],"data_offsets":[0,1]},` +
		`"conditioner.embedders.1.model.ln_final.weight":{"dtype":"F16","shape":[1],"data_offsets":[0,1]},` +
		`"model.diffusion_model.input_blocks.0.0.weight":{"dtype":"F16","shape":[1],"data_offsets":[0,1]}}`
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint64(len(header)))
	buf.WriteString(header)
	if err := os.WriteFile(full, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	opts := imagetypes.ModelLoadOptions{
		VAEPath:   filepath.Join(dir, "borrowed.vae.safetensors"),
		ClipLPath: filepath.Join(dir, "borrowed.clip_l.safetensors"),
		ClipGPath: filepath.Join(dir, "borrowed.clip_g.safetensors"),
	}
	clearEmbeddedCompanions(&opts, full)
	if opts.VAEPath != "" || opts.ClipLPath != "" || opts.ClipGPath != "" {
		t.Fatalf("borrowed companions survived: vae=%q clip_l=%q clip_g=%q", opts.VAEPath, opts.ClipLPath, opts.ClipGPath)
	}

	// A GGUF diffusion model carries none of those roles, so its companions stay.
	gguf := imagetypes.ModelLoadOptions{VAEPath: "keep.vae", ClipLPath: "keep.clip_l"}
	clearEmbeddedCompanions(&gguf, filepath.Join(dir, "model.gguf"))
	if gguf.VAEPath != "keep.vae" || gguf.ClipLPath != "keep.clip_l" {
		t.Fatalf("GGUF companions were dropped: vae=%q clip_l=%q", gguf.VAEPath, gguf.ClipLPath)
	}
}
