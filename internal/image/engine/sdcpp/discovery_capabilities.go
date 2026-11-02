package sdcpp

import (
	"strings"

	"howl-chat/internal/image/scanner"
	"howl-chat/internal/image/types"
)

// Companion resolution now reads the shared capability table instead of a
// per-family switch. The old arrangement had the engine and the discovery pass
// disagreeing about which roles a family needs (HiDream was "self-contained" in
// one and "needs T5XXL + VAE" in the other), which is how a model can appear
// loadable and then fail on load.

// indexCount returns how many scanned assets of a role exist, so a model can
// borrow a companion from any registered directory.
func indexCount(idx scanner.ModelIndex, role string) int {
	switch role {
	case "vae":
		return len(idx.VAEs)
	case "clip_l":
		return len(idx.ClipLs)
	case "clip_g":
		return len(idx.ClipGs)
	case "t5xxl":
		return len(idx.T5XXLs)
	case "llm":
		return len(idx.LLMs)
	}
	return 0
}

// isModelLoadable checks that every required companion exists either on the
// model itself or somewhere in the scanned index.
func isModelLoadable(info *types.ModelInfo, idx scanner.ModelIndex) bool {
	family := strings.ToLower(strings.TrimSpace(info.Family))

	// SD3.5 GGUF is broken in stable-diffusion.cpp (upstream issue #1105): the
	// safetensors version loads, the quantised one crashes during load.
	if family == "sd3.5" && info.Format == "gguf" {
		info.Compatibility = "unsupported: SD3.5 GGUF is currently broken in stable-diffusion.cpp. Use the .safetensors version instead."
		return false
	}

	cap, known := Capability(family)
	if !known {
		// Unknown family: list it with a warning rather than hiding it. The user
		// can see what was found and why it will not load.
		info.Compatibility = "unsupported: unknown model family"
		return true
	}

	var missing []string
	for _, role := range cap.Companions {
		if info.Files[role] != "" || indexCount(idx, role) > 0 {
			continue
		}
		missing = append(missing, role)
	}
	if len(missing) > 0 {
		info.Compatibility = "unsupported: missing " + strings.Join(missing, ", ")
		return false
	}
	return true
}

// autoFillCompanions attaches companions this family can use, and removes roles
// it cannot: a leftover T5 or CLIP-L from an earlier family guess must not
// follow a model into a load.
func autoFillCompanions(info *types.ModelInfo, idx scanner.ModelIndex) {
	family := strings.ToLower(strings.TrimSpace(info.Family))
	cap, known := Capability(family)
	if !known {
		cap = DefaultCapability
	}
	if info.Files == nil {
		info.Files = map[string]string{}
	}

	for _, role := range []string{"vae", "clip_l", "clip_g", "t5xxl", "llm", "clip_vision"} {
		if !cap.Uses(role) {
			delete(info.Files, role)
		}
	}

	for _, role := range cap.Roles() {
		if info.Files[role] != "" {
			continue
		}
		switch role {
		case "vae":
			if p := pickCompatibleVAE(family, idx); p != "" {
				info.Files["vae"] = p
			}
		// Encoders go through pickCompanion, not index[0]: taking whatever the
		// scan walked first recorded a bitsandbytes int8 T5 from an InvokeAI
		// tree onto flux1-schnell, and sd-server aborted on the F32 weight-scale
		// assert the moment it tried to load it.
		case "clip_l":
			if p := pickCompanion(role, idx.ClipLs, info.Path); p != "" {
				info.Files["clip_l"] = p
			}
		case "clip_g":
			if p := pickCompanion(role, idx.ClipGs, info.Path); p != "" {
				info.Files["clip_g"] = p
			}
		case "t5xxl":
			if p := pickCompanion(role, idx.T5XXLs, info.Path); p != "" {
				info.Files["t5xxl"] = p
			}
		case "llm":
			if p := pickCompanion(role, idx.LLMs, info.Path); p != "" {
				info.Files["llm"] = p
			}
		}
	}
}

// pickCompatibleVAE returns a VAE this family can decode with. The Flux-style
// ae.safetensors is never attached to an SD checkpoint: it decodes noise.
func pickCompatibleVAE(family string, idx scanner.ModelIndex) string {
	for _, vae := range idx.VAEs {
		lower := strings.ToLower(vae.Name)
		pathLower := strings.ToLower(vae.Path)
		isFluxAE := strings.Contains(lower, "ae.safetensors") || strings.Contains(pathLower, "ae.safetensors")
		if isFluxAE && !UsesFluxAE(family) {
			continue
		}
		if scanner.CompatibleVAE(scanner.ModelFamily(family), vae.Family) {
			return vae.Path
		}
	}
	return ""
}
