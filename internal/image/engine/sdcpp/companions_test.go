package sdcpp

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/image/scanner"
	"howl-chat/internal/image/types"
)

// sep builds a native path so the volume logic is exercised on Windows and
// accepted everywhere else.
func sep() string { return string(filepath.Separator) }

func TestPickCompanionPrefersTheModelDirectory(t *testing.T) {
	model := "G:" + sep() + "AI GGUF" + sep() + "Image models" + sep() + "flux1-schnell-Q4_K_M.gguf"
	candidates := []scanner.AssetEntry{
		{Path: "H:" + sep() + "InvokeAI" + sep() + "text_encoder_2" + sep() + "model.safetensors"},
		{Path: "G:" + sep() + "AI GGUF" + sep() + "Image models" + sep() + "t5xxl_fp8_e4m3fn.safetensors"},
	}
	got := pickCompanion("t5xxl", candidates, model)
	if !strings.Contains(got, "t5xxl_fp8_e4m3fn.safetensors") {
		t.Errorf("pickCompanion = %q, want the encoder that sits beside the model", got)
	}
}

func TestPickCompanionRejectsBitsAndBytesEncoder(t *testing.T) {
	model := "G:" + sep() + "models" + sep() + "flux1-schnell-Q4_K_M.gguf"
	candidates := []scanner.AssetEntry{
		{Path: "H:" + sep() + "InvokeAI" + sep() + "text_encoder_2" + sep() + "bnb_llm_int8_model.safetensors"},
	}
	if got := pickCompanion("t5xxl", candidates, model); got != "" {
		t.Errorf("pickCompanion = %q; a bitsandbytes int8 encoder must never be offered", got)
	}
}

func TestPickCompanionIsStableWithoutALocalFile(t *testing.T) {
	model := "G:" + sep() + "models" + sep() + "flux.gguf"
	candidates := []scanner.AssetEntry{
		{Path: "H:" + sep() + "b" + sep() + "t5xxl_a.safetensors"},
		{Path: "H:" + sep() + "a" + sep() + "t5xxl_b.safetensors"},
	}
	first := pickCompanion("t5xxl", candidates, model)
	reverse := []scanner.AssetEntry{candidates[1], candidates[0]}
	if second := pickCompanion("t5xxl", reverse, model); first != second {
		t.Errorf("order changed the pick: %q vs %q", first, second)
	}
	// The tie breaks on the whole path, so the directory decides: "H:\\a\\..."
	// sorts before "H:\\b\\..." no matter what the file is called.
	if !strings.Contains(first, "t5xxl_b.safetensors") {
		t.Errorf("pickCompanion = %q, want the file under the lowest directory", first)
	}
}

// A full checkpoint carries its own VAE and both text encoders. Attaching a
// borrowed CLIP-L from the index left SDXL without its CLIP-G, and the render
// was noise. Only roles the file does not carry may be borrowed.
func TestDropEmbeddedCompanionsKeepsFullCheckpointSelfContained(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "full.safetensors")
	header := `{"first_stage_model.decoder.conv_in.weight":{"dtype":"F16","shape":[1],"data_offsets":[0,1]},` +
		`"conditioner.embedders.1.model.ln_final.weight":{"dtype":"F16","shape":[1],"data_offsets":[0,1]},` +
		`"model.diffusion_model.input_blocks.0.0.weight":{"dtype":"F16","shape":[1],"data_offsets":[0,1]}}`
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint64(len(header)))
	buf.WriteString(header)
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &types.ModelInfo{
		Family: "sdxl",
		Format: "safetensors",
		Files: map[string]string{
			"model":  path,
			"vae":    filepath.Join(dir, "borrowed.vae.safetensors"),
			"clip_l": filepath.Join(dir, "borrowed.clip_l.safetensors"),
			"clip_g": filepath.Join(dir, "borrowed.clip_g.safetensors"),
		},
	}
	dropEmbeddedCompanions(m, types.ModelLoadOptions{})

	for _, role := range []string{"vae", "clip_l", "clip_g"} {
		if m.Files[role] != "" {
			t.Fatalf("%s still borrowed from the index: %q", role, m.Files[role])
		}
	}
}

func TestDropEmbeddedCompanionsKeepsExplicitChoices(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "full.safetensors")
	header := `{"first_stage_model.decoder.conv_in.weight":{"dtype":"F16","shape":[1],"data_offsets":[0,1]},` +
		`"conditioner.embedders.1.model.ln_final.weight":{"dtype":"F16","shape":[1],"data_offsets":[0,1]}}`
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint64(len(header)))
	buf.WriteString(header)
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	chosenVAE := filepath.Join(dir, "chosen.vae.safetensors")
	m := &types.ModelInfo{
		Family: "sdxl",
		Format: "safetensors",
		Files:  map[string]string{"model": path, "vae": chosenVAE, "clip_l": "borrowed.clip_l.safetensors"},
	}
	dropEmbeddedCompanions(m, types.ModelLoadOptions{VAEPath: chosenVAE})

	if m.Files["vae"] != chosenVAE {
		t.Fatalf("explicit VAE was dropped: %q", m.Files["vae"])
	}
	if m.Files["clip_l"] != "" {
		t.Fatalf("borrowed CLIP-L survived: %q", m.Files["clip_l"])
	}
}

// The app layer asks the same question before it borrows a file: a full SDXL
// checkpoint reports the roles it carries so nothing overwrites them.
func TestEmbeddedRolesReportsFullCheckpointRoles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "full.safetensors")
	header := `{"first_stage_model.decoder.conv_in.weight":{"dtype":"F16","shape":[1],"data_offsets":[0,1]},` +
		`"conditioner.embedders.0.model.ln_final.weight":{"dtype":"F16","shape":[1],"data_offsets":[0,1]},` +
		`"conditioner.embedders.1.model.ln_final.weight":{"dtype":"F16","shape":[1],"data_offsets":[0,1]},` +
		`"model.diffusion_model.input_blocks.0.0.weight":{"dtype":"F16","shape":[1],"data_offsets":[0,1]}}`
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint64(len(header)))
	buf.WriteString(header)
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	roles := EmbeddedRoles(path)
	for _, want := range []string{"vae", "clip_l", "clip_g"} {
		if !roles[want] {
			t.Fatalf("EmbeddedRoles missing %q: %v", want, roles)
		}
	}
	if roles := EmbeddedRoles(filepath.Join(dir, "unet.gguf")); len(roles) != 0 {
		t.Fatalf("a GGUF diffusion model must not report embedded roles: %v", roles)
	}
}
