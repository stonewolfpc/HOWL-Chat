package sdcpp

import (
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/image/scanner"
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
