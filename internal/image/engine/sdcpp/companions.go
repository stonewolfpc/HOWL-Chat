package sdcpp

import (
	"path/filepath"
	"strings"

	"howl-chat/internal/image/scanner"
)

// unusableEncoderQuant reports whether a text-encoder file uses a quantisation
// scheme stable-diffusion.cpp cannot read.
//
// bitsandbytes exports (bnb / int8 / nf4 / int4) come from Diffusers checkpoints,
// not from sd.cpp: they carry their scale factors in a separate module and sd.cpp
// aborts the moment it looks for them —
//
//	GGML_ASSERT(weight_scale != NULL && weight_scale->type == GGML_TYPE_F32) failed
//
// which is exactly how flux1-schnell died after being handed
// H:\InvokeAI\...\text_encoder_2\bnb_llm_int8_model.safetensors. The scheme is
// always in the filename on these exports, and nothing else is available to
// inspect before a load, so the name is the check.
func unusableEncoderQuant(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	for _, marker := range []string{"bnb", "bitsandbytes", "int8", "int4", "nf4", "fp4", "awq", "gptq"} {
		if strings.Contains(base, marker) {
			return true
		}
	}
	return false
}

// pickCompanion chooses one companion file for a role out of the scanned index.
//
// Two rules, both learned from loads that failed inside sd-server:
//
//   - A file sitting beside the model wins over one from another tree. Flux and
//     Z-Image ship their encoders next to the checkpoint; reaching into an
//     unrelated directory instead is how flux1-schnell picked up an InvokeAI
//     bitsandbytes encoder while the correct fp8 T5 sat in its own folder.
//   - A bitsandbytes-quantised encoder is never eligible, wherever it lives.
//
// Within the same tier the lowest path wins, so repeated scans pick the same
// file. Before this, the choice was simply index[0] — whatever the scan happened
// to walk first.
func pickCompanion(role string, candidates []scanner.AssetEntry, modelPath string) string {
	modelDir := ""
	if modelPath != "" {
		modelDir = strings.ToLower(filepath.Dir(modelPath))
	}
	modelVolume := ""
	if modelPath != "" {
		modelVolume = strings.ToLower(filepath.VolumeName(modelPath))
	}

	best := ""
	bestRank := 3
	for _, c := range candidates {
		if c.Path == "" || unusableEncoderQuant(c.Path) {
			continue
		}
		rank := 2
		switch {
		case modelDir != "" && strings.ToLower(filepath.Dir(c.Path)) == modelDir:
			rank = 0
		case modelVolume != "" && strings.ToLower(filepath.VolumeName(c.Path)) == modelVolume:
			rank = 1
		}
		if rank < bestRank || (rank == bestRank && (best == "" || c.Path < best)) {
			best, bestRank = c.Path, rank
		}
	}
	return best
}
