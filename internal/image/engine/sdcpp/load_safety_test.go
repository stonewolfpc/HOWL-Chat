package sdcpp

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"howl-chat/internal/image/types"
)

// writeGGUF writes a file with a real GGUF prefix and the given metadata count.
// kv == 0 is the shape of a stripped checkpoint: 453 tensors, no
// general.architecture, and sd.cpp can only answer "get sd version failed".
func writeGGUF(t *testing.T, dir string, kv uint64) string {
	t.Helper()
	path := filepath.Join(dir, "checkpoint.gguf")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString("GGUF"); err != nil {
		t.Fatal(err)
	}
	for _, v := range []interface{}{uint32(3), uint64(453), kv} {
		if err := binary.Write(f, binary.LittleEndian, v); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := f.Write(make([]byte, 64)); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestValidateCheckpointRejectsGGUFWithoutMetadata(t *testing.T) {
	path := writeGGUF(t, t.TempDir(), 0)
	err := validateCheckpoint(&types.ModelInfo{Files: map[string]string{"model": path}})
	if err == nil {
		t.Fatal("a GGUF with no metadata must be rejected before launch")
	}
	if !strings.Contains(err.Error(), "no metadata") {
		t.Errorf("error does not name the real cause: %v", err)
	}
	if strings.Contains(err.Error(), "update your sd-server") {
		t.Error("error blames the binary version for a broken file")
	}
}

func TestValidateCheckpointAcceptsUsableGGUF(t *testing.T) {
	path := writeGGUF(t, t.TempDir(), 28)
	if err := validateCheckpoint(&types.ModelInfo{Files: map[string]string{"model": path}}); err != nil {
		t.Fatalf("usable GGUF rejected: %v", err)
	}
}

func TestValidateCheckpointIgnoresNonGGUF(t *testing.T) {
	path := filepath.Join(t.TempDir(), "model.safetensors")
	if err := os.WriteFile(path, []byte("not a gguf at all"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := validateCheckpoint(&types.ModelInfo{Files: map[string]string{"model": path}}); err != nil {
		t.Fatalf("safetensors must not be probe-rejected: %v", err)
	}
}

// The wait has to grow with the bytes on disk, or a 17 GB encoder on a slow
// drive is declared a failure while the machine still has room for it.
func TestLoadTimeoutGrowsWithWeightsAndStaysBounded(t *testing.T) {
	dir := t.TempDir()
	small := filepath.Join(dir, "small.bin")
	if err := os.WriteFile(small, make([]byte, 1024), 0644); err != nil {
		t.Fatal(err)
	}

	base := loadTimeoutFor(map[string]string{"model": small})
	if base < loadBaseTimeout {
		t.Fatalf("timeout %s is below the base %s", base, loadBaseTimeout)
	}
	if got := loadTimeoutFor(nil); got != loadBaseTimeout {
		t.Errorf("no files: %s, want %s", got, loadBaseTimeout)
	}

	// 2 GiB at 40 MB/s adds ~51s, which must show up in the budget.
	big := filepath.Join(dir, "big.bin")
	f, err := os.Create(big)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(2 << 30); err != nil {
		f.Close()
		t.Skipf("cannot create sparse file: %v", err)
	}
	f.Close()
	got := loadTimeoutFor(map[string]string{"model": big})
	if got <= base {
		t.Errorf("timeout did not grow with weight size: %s vs %s", got, base)
	}
	if got > loadMaxTimeout {
		t.Errorf("timeout %s exceeds the cap %s", got, loadMaxTimeout)
	}
}

// A companion the caller supplied must reach sd-server even when the family has
// no entry in the capability table.
func TestUnknownFamilyStillReceivesSuppliedCompanions(t *testing.T) {
	e := &Engine{config: types.EngineConfig{Port: 8081}}
	m := &types.ModelInfo{
		Family: "mystery-arch",
		Files: map[string]string{
			"model": `E:\models\m.gguf`,
			"llm":   `E:\models\enc.safetensors`,
			"vae":   `E:\models\ae.safetensors`,
		},
	}
	args := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{}), " ")
	for _, want := range []string{"--llm", "--vae"} {
		if !strings.Contains(args, want) {
			t.Errorf("%s missing for an unknown family: %s", want, args)
		}
	}
}

func TestFlashAttentionUsesCurrentFlagName(t *testing.T) {
	e := &Engine{config: types.EngineConfig{Port: 8081}}
	m := &types.ModelInfo{Family: "sd1.5", Files: map[string]string{"model": `E:\m.safetensors`}}
	args := e.buildServerArgs(m, types.ModelLoadOptions{FlashAttention: true})
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--flash-attn on") {
		t.Errorf("flash attention did not use --flash-attn on: %s", joined)
	}
	for i, a := range args {
		if a == "--fa" {
			t.Errorf("deprecated --fa still emitted at %d: %s", i, joined)
		}
	}
}

// Auto-fit is what makes an 8 GB card run a 6 GB transformer: any explicit
// params backend disables it, so it must not be sent unless the user asked.
func TestParamsBackendOnlySentWhenRequested(t *testing.T) {
	e := &Engine{config: types.EngineConfig{Port: 8081}}
	m := &types.ModelInfo{Family: "z_image", Files: map[string]string{"model": `E:\m.gguf`}}

	def := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{}), " ")
	if strings.Contains(def, "--params-backend") {
		t.Errorf("params backend sent by default, disabling auto-fit: %s", def)
	}
	if !strings.Contains(def, "--max-vram -1") {
		t.Errorf("VRAM reserve missing: %s", def)
	}

	asked := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{ParamsBackend: "cpu"}), " ")
	if !strings.Contains(asked, "--params-backend cpu") {
		t.Errorf("explicit params backend dropped: %s", asked)
	}
}

// Z-Image must arrive with the split-component flag and both companions; that
// combination is what actually loads on this machine.
func TestZImageArgsCarryDiffusionModelAndCompanions(t *testing.T) {
	e := &Engine{config: types.EngineConfig{Port: 8081}}
	m := &types.ModelInfo{
		Family: "z_image",
		Format: "gguf",
		Files: map[string]string{
			"model": `E:\image\z.gguf`,
			"llm":   `H:\models\qwen.safetensors`,
			"vae":   `E:\image\ae.safetensors`,
		},
	}
	args := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{}), " ")
	for _, want := range []string{"--diffusion-model", "--llm", "--vae"} {
		if !strings.Contains(args, want) {
			t.Errorf("%s missing from z_image launch: %s", want, args)
		}
	}
	if strings.Contains(args, " -m ") {
		t.Errorf("split-component model passed as -m: %s", args)
	}
}

var _ = time.Second
