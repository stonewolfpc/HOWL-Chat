package sdcpp

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeSafetensorsHeader(t *testing.T, dir string, tensors map[string][2]int64) string {
	t.Helper()
	header := make(map[string]interface{}, len(tensors))
	for name, offsets := range tensors {
		header[name] = map[string]interface{}{
			"dtype":        "F16",
			"shape":        []int{1},
			"data_offsets": []int64{offsets[0], offsets[1]},
		}
	}
	raw, err := json.Marshal(header)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "model.safetensors")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := binary.Write(f, binary.LittleEndian, uint64(len(raw))); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(raw); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReadComponentBytesSplitsFullCheckpoint(t *testing.T) {
	path := writeSafetensorsHeader(t, t.TempDir(), map[string][2]int64{
		"model.diffusion_model.input_blocks.0.weight": {0, 400},
		"conditioner.text_model.proj.weight":          {400, 550},
		"first_stage_model.decoder.conv.weight":       {550, 600},
		"some.metadata.tensor":                        {600, 610},
	})
	parts, ok := ReadComponentBytes(path)
	if !ok {
		t.Fatal("a valid safetensors header with known components must split")
	}
	if parts.Diffusion != 400 || parts.Conditioner != 150 || parts.VAE != 50 || parts.Other != 10 {
		t.Fatalf("unexpected split: %+v", parts)
	}
	if parts.Total() != 610 {
		t.Fatalf("total = %d, want 610", parts.Total())
	}
}

func TestReadComponentBytesFallsBackWhenNothingIsKnown(t *testing.T) {
	path := writeSafetensorsHeader(t, t.TempDir(), map[string][2]int64{
		"unknown.tensor": {0, 100},
	})
	if _, ok := ReadComponentBytes(path); ok {
		t.Fatal("an unclassified header must fall back to the whole-file estimate")
	}
}

func TestReadComponentBytesRejectsOtherFormatsAndMalformedHeaders(t *testing.T) {
	dir := t.TempDir()
	gguf := filepath.Join(dir, "model.gguf")
	if err := os.WriteFile(gguf, []byte("GGUF"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, ok := ReadComponentBytes(gguf); ok {
		t.Fatal("GGUF must keep the whole-file estimate")
	}

	bad := filepath.Join(dir, "bad.safetensors")
	if err := os.WriteFile(bad, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 0644); err != nil {
		t.Fatal(err)
	}
	if _, ok := ReadComponentBytes(bad); ok {
		t.Fatal("an impossible header length must fall back")
	}
}

func TestComponentRolePrefixes(t *testing.T) {
	cases := map[string]string{
		"model.diffusion_model.x": "diffusion",
		"conditioner.x":           "conditioner",
		"text_encoders.clip.x":    "conditioner",
		"first_stage_model.x":     "vae",
		"vae.x":                   "vae",
		"unet.x":                  "",
	}
	for name, want := range cases {
		if got := componentRole(name); got != want {
			t.Errorf("componentRole(%q) = %q, want %q", name, got, want)
		}
	}
}
