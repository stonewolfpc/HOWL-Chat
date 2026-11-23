package sdcpp

import (
	"path/filepath"
	"strings"

	"howl-chat/internal/image/scanner"
	"howl-chat/internal/image/types"
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

// dropEmbeddedCompanions stops a self-contained checkpoint from being overridden
// by files borrowed from the scanned index.
//
// A full SDXL .safetensors already carries its UNet, both text encoders and the
// VAE. autoFillCompanions attached an external CLIP-L and VAE anyway, so
// sd-server loaded the file's UNet but not its CLIP-G: Illustrious-XL-v2.0 at
// 768x768 / 20 steps / euler_a produced coloured noise instead of an image.
// A role the file itself carries stays the file's. An explicit choice from the
// UI still wins, because the user asked for that file by name.
func dropEmbeddedCompanions(m *types.ModelInfo, opts types.ModelLoadOptions) {
	if m == nil || m.Files["model"] == "" || !strings.EqualFold(filepath.Ext(m.Files["model"]), ".safetensors") {
		return
	}
	probe := scanner.ProbeWeight(m.Files["model"])
	if probe.FirstStage && opts.VAEPath == "" {
		delete(m.Files, "vae")
	}
	if probe.CondEmbedders {
		if opts.ClipLPath == "" {
			delete(m.Files, "clip_l")
		}
		if opts.ClipGPath == "" {
			delete(m.Files, "clip_g")
		}
	}
}

// EmbeddedRoles reports which companion roles a checkpoint file already carries
// inside itself.
//
// Only safetensors is inspected: it stores a full SDXL checkpoint's VAE
// (first_stage_model.*) and both text encoders (conditioner.embedders.0/1.*)
// in one file. Borrowing files for those roles replaces working weights with
// unrelated ones — an external CLIP-L with no CLIP-G left SDXL unable to form
// its pooled embedding, and every render came back as coloured noise.
func EmbeddedRoles(path string) map[string]bool {
	roles := map[string]bool{}
	if path == "" || !strings.EqualFold(filepath.Ext(path), ".safetensors") {
		return roles
	}
	probe := scanner.ProbeWeight(path)
	if probe.FirstStage {
		roles["vae"] = true
	}
	if probe.ClipL {
		roles["clip_l"] = true
	}
	if probe.CondEmbedders {
		roles["clip_g"] = true
	}
	return roles
}
