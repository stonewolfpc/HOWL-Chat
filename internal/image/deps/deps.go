// Package deps resolves what an image model needs from companion files, tells
// the user what is missing, and knows where each missing file can be obtained.
//
// Why this exists: a family like flux2 cannot load without a Mistral text
// encoder and a VAE, and a user who downloaded only the diffusion weights gets a
// load failure that names a role ("missing llm") and nothing else. This package
// turns that into a concrete list: which role, which file, how big, where to get
// it — and, when no direct link is known, the page to open by hand.
//
// It also answers the harder question: is the file that is already on disk the
// *right* file for this family? Attaching whatever LLM happens to be in the
// registry made z_image load with a Llama GGUF and qwen_image with a LoRA, then
// fail deep inside sd-server. The catalogue records what each family's companion
// must look like, so a wrong file is reported as wrong instead of "present".
//
// The catalogue is baked into the binary on purpose. HOWL Chat never contacts
// the network by itself: a lookup here is pure data, and every download is a
// separate, explicitly consented action performed by the caller.
package deps

import (
	"path/filepath"
	"sort"
	"strings"

	"howl-chat/internal/image/engine/sdcpp"
)

// Artifact is one concrete file that satisfies a companion role.
type Artifact struct {
	// URL is a direct download. Empty means the file is not fetched
	// automatically — Source is the page the user opens instead.
	URL string
	// FileName is the name the file should have on disk.
	FileName string
	// Source is the human-readable page for the file, shown next to a download
	// and used as the fallback when URL is empty.
	Source string
	// SizeMB is the approximate download size, or 0 when it is not known.
	SizeMB int
	// Note explains anything the user should know before downloading.
	Note string
	// Match lists lowercase name fragments that identify a file as suitable for
	// this role and family. A file whose name contains none of them is a
	// different model of the same kind (a Llama GGUF where a Qwen3 encoder
	// belongs) and is reported as a mismatch rather than silently attached.
	// Empty means this package cannot judge identity for the role.
	Match []string
}

// ConsentText is the sentence the UI must show before any network access. It
// lives here so there is one wording of the promise, not one per screen.
const ConsentText = "HOWL Chat needs to connect to the internet to check the dependencies this model requires and download them. Do you agree? Click Continue to proceed, or Cancel to stop."

// Catalog maps a family and role to the file that satisfies it.
//
// Every URL here was checked against the publisher's repository listing. A role
// with an empty URL is one where no stable direct link is known: the user gets
// the source page instead of a button that 404s.
var Catalog = map[string]map[string]Artifact{
	"flux": {
		"t5xxl": {
			FileName: "t5xxl_fp8_e4m3fn.safetensors",
			URL:      "https://huggingface.co/comfyanonymous/flux_text_encoders/resolve/main/t5xxl_fp8_e4m3fn.safetensors",
			Source:   "https://huggingface.co/comfyanonymous/flux_text_encoders",
			SizeMB:   4890,
			Note:     "fp8 T5-XXL. The fp16 file is ~9.8 GB and needs more RAM.",
			Match:    []string{"t5xxl", "t5-xxl", "t5_xxl"},
		},
		"clip_l": {
			FileName: "clip_l.safetensors",
			URL:      "https://huggingface.co/comfyanonymous/flux_text_encoders/resolve/main/clip_l.safetensors",
			Source:   "https://huggingface.co/comfyanonymous/flux_text_encoders",
			SizeMB:   246,
			Match:    []string{"clip_l", "clip-l", "clip.l"},
		},
		"vae": {
			FileName: "ae.safetensors",
			URL:      "https://huggingface.co/black-forest-labs/FLUX.1-schnell/resolve/main/ae.safetensors",
			Source:   "https://huggingface.co/black-forest-labs/FLUX.1-schnell/tree/main",
			SizeMB:   336,
			Note:     "The Flux autoencoder. Other VAEs decode Flux latents as noise.",
			Match:    []string{"ae"},
		},
	},
	"flux2": {
		"llm": {
			FileName: "mistral_3_small_flux2_fp8.safetensors",
			URL:      "https://huggingface.co/Comfy-Org/Flux-2-Dev/resolve/main/split_files/text_encoders/mistral_3_small_flux2_fp8.safetensors",
			Source:   "https://huggingface.co/Comfy-Org/Flux-2-Dev/tree/main/split_files/text_encoders",
			SizeMB:   17500,
			Note:     "Mistral Small 3 text encoder required by FLUX.2. It is the largest dependency here.",
			Match:    []string{"mistral"},
		},
		"vae": {
			FileName: "flux2-vae.safetensors",
			URL:      "https://huggingface.co/Comfy-Org/Flux-2-Dev/resolve/main/split_files/vae/flux2-vae.safetensors",
			Source:   "https://huggingface.co/Comfy-Org/Flux-2-Dev/tree/main/split_files/vae",
			SizeMB:   336,
			Match:    []string{"flux2-vae", "flux2_vae", "flux.2-vae"},
		},
	},
	"z_image": {
		"llm": {
			FileName: "qwen_3_4b.safetensors",
			URL:      "https://huggingface.co/Comfy-Org/Z-Image-Turbo/resolve/main/split_files/text_encoders/qwen_3_4b.safetensors",
			Source:   "https://huggingface.co/Comfy-Org/Z-Image-Turbo/tree/main/split_files",
			SizeMB:   8040,
			Note:     "Qwen3-4B (2560-dim). A Qwen3-0.6B encoder loads but cannot drive Z-Image: the DiT expects 4B hidden states.",
			Match:    []string{"qwen_3_4b", "qwen3-4b", "qwen3_4b"},
		},
		"vae": {
			FileName: "ae.safetensors",
			URL:      "https://huggingface.co/Comfy-Org/Z-Image-Turbo/resolve/main/split_files/vae/ae.safetensors",
			Source:   "https://huggingface.co/Comfy-Org/Z-Image-Turbo/tree/main/split_files/vae",
			SizeMB:   336,
			Note:     "Z-Image uses the Flux-style autoencoder.",
			Match:    []string{"ae"},
		},
	},
	"qwen_image": {
		"llm": {
			FileName: "qwen_2.5_vl_7b_fp8_scaled.safetensors",
			URL:      "https://huggingface.co/Comfy-Org/Qwen-Image/resolve/main/split_files/text_encoders/qwen_2.5_vl_7b_fp8_scaled.safetensors",
			Source:   "https://huggingface.co/Comfy-Org/Qwen-Image/tree/main/split_files/text_encoders",
			SizeMB:   8600,
			Note:     "Qwen2.5-VL 7B text encoder used by Qwen-Image.",
			Match:    []string{"qwen_2.5_vl", "qwen2.5-vl", "qwen2_5_vl"},
		},
		"vae": {
			FileName: "qwen_image_vae.safetensors",
			URL:      "https://huggingface.co/Comfy-Org/Qwen-Image/resolve/main/split_files/vae/qwen_image_vae.safetensors",
			Source:   "https://huggingface.co/Comfy-Org/Qwen-Image/tree/main/split_files/vae",
			SizeMB:   254,
			Match:    []string{"qwen_image_vae", "qwen-image-vae"},
		},
	},
	"anima": {
		"llm": {
			FileName: "qwen_3_06b_base.safetensors",
			URL:      "https://huggingface.co/circlestone-labs/Anima/resolve/main/split_files/text_encoders/qwen_3_06b_base.safetensors",
			Source:   "https://huggingface.co/circlestone-labs/Anima/tree/main/split_files/text_encoders",
			SizeMB:   1190,
			Note:     "Anima (Cosmos-Predict2-2B) takes the Qwen3-0.6B-Base encoder, not the 4B one Z-Image uses.",
			Match:    []string{"qwen_3_06b", "qwen3-0.6b", "qwen3_0.6b", "qwen3-06b"},
		},
		"vae": {
			FileName: "qwen_image_vae.safetensors",
			URL:      "https://huggingface.co/circlestone-labs/Anima/resolve/main/split_files/vae/qwen_image_vae.safetensors",
			Source:   "https://huggingface.co/circlestone-labs/Anima/tree/main/split_files/vae",
			SizeMB:   254,
			Note:     "Anima ships the Qwen-Image VAE. The Flux ae.safetensors decodes its latents as noise.",
			Match:    []string{"qwen_image_vae", "qwen-image-vae"},
		},
	},
}

// RoleLabel is the user-facing name of a companion role.
func RoleLabel(role string) string {
	switch role {
	case "vae":
		return "VAE"
	case "clip_l":
		return "CLIP-L text encoder"
	case "clip_g":
		return "CLIP-G text encoder"
	case "t5xxl":
		return "T5-XXL text encoder"
	case "llm":
		return "Text encoder (LLM)"
	case "clip_vision":
		return "CLIP vision encoder"
	case "diffusion_model":
		return "Diffusion model"
	}
	return role
}

// Compatible reports whether a file name is the right kind of companion for a
// family and role. The reason is written for the user: it names what the file
// looks like and what the family needs instead.
//
// An empty Match means this package cannot judge identity, so the file is
// accepted rather than rejected on a guess.
func Compatible(family, role, name string) (bool, string) {
	art := artifactFor(strings.ToLower(strings.TrimSpace(family)), role)
	if len(art.Match) == 0 {
		return true, ""
	}
	lower := strings.ToLower(filepath.Base(name))
	for _, m := range art.Match {
		if matchesName(lower, m) {
			return true, ""
		}
	}
	why := filepath.Base(name) + " is not a " + RoleLabel(role) + " for " + family + " models"
	if art.FileName != "" {
		why += " (expected something like " + art.FileName + ")"
	}
	return false, why
}

// matchesName reports whether a file name starts with a catalogue pattern at a
// token boundary. Plain substring matching is wrong here: "ae" inside
// "qwen_image_vae.safetensors" would accept the wrong autoencoder for Flux,
// which is exactly the kind of silent mismatch this check exists to catch.
func matchesName(lower, pattern string) bool {
	if !strings.HasPrefix(lower, pattern) {
		return false
	}
	if len(lower) == len(pattern) {
		return true
	}
	return strings.ContainsRune("._- ", rune(lower[len(pattern)]))
}

// Item is one requirement in a model's report.
type Item struct {
	Role  string
	Label string
	// State is "present", "missing", "mismatch" or "optional". "missing" and
	// "mismatch" both block a load.
	State string
	// Path is the file the loader would use, when one was found.
	Path string
	// Mismatch explains why a file that was found cannot be used.
	Mismatch string
	// Artifact is filled for every state, so a present file and a missing one
	// are described the same way.
	Artifact
}

// Report is the full dependency picture for one model.
type Report struct {
	Model string
	// Family is the architecture the model was classified as.
	Family string
	// Known is false when the family is not implemented by this engine build,
	// which is a different problem from a missing file.
	Known bool
	Items []Item
	// Required lists the roles that must be present before a load can succeed.
	Required []string
	// Missing lists the required roles that are absent or unusable right now.
	Missing []string
	// DownloadDir is where this app puts fetched companions.
	DownloadDir string
}

// BlocksLoad reports whether anything is missing that would fail a load.
func (r Report) BlocksLoad() bool { return len(r.Missing) > 0 }

// StatePresent / StateMissing / StateMismatch / StateOptional are item states.
const (
	StatePresent  = "present"
	StateMissing  = "missing"
	StateMismatch = "mismatch"
	StateOptional = "optional"
)

// Resolve builds the report for a model.
//
// present maps a role to the file the loader would use for it. A file that
// exists but is the wrong kind for this family is reported as a mismatch, not as
// present: that distinction is the difference between "you have what you need"
// and "this will fail on load".
func Resolve(model, family string, present map[string]string, downloadDir string) Report {
	family = strings.ToLower(strings.TrimSpace(family))
	r := Report{Model: model, Family: family, DownloadDir: downloadDir}

	cap, known := sdcpp.Capability(family)
	r.Known = known
	if !known {
		cap = sdcpp.DefaultCapability
	}

	for _, role := range cap.Roles() {
		required := false
		for _, c := range cap.Companions {
			if c == role {
				required = true
				break
			}
		}
		item := Item{Role: role, Label: RoleLabel(role), Artifact: artifactFor(family, role)}
		if path := strings.TrimSpace(present[role]); path != "" {
			item.Path = path
			if ok, why := Compatible(family, role, path); ok {
				item.State = StatePresent
			} else {
				item.State = StateMismatch
				item.Mismatch = why
				if required {
					r.Required = append(r.Required, role)
					r.Missing = append(r.Missing, role)
				}
			}
		} else if required {
			item.State = StateMissing
			r.Required = append(r.Required, role)
			r.Missing = append(r.Missing, role)
		} else {
			item.State = StateOptional
		}
		r.Items = append(r.Items, item)
	}
	return r
}

// artifactFor returns the catalogue entry for a family and role, falling back to
// a related family when one ships the same file, so a Flux VAE can satisfy a
// Flux-derived architecture without a second catalogue entry.
func artifactFor(family, role string) Artifact {
	if fam, ok := Catalog[family]; ok {
		if a, ok := fam[role]; ok {
			return a
		}
	}
	for _, src := range sharedRoles[role] {
		if a, ok := Catalog[src][role]; ok {
			return a
		}
	}
	return Artifact{}
}

// sharedRoles lists families that ship byte-identical files for a role. Only
// genuinely interchangeable files belong here: encoders and autoencoders differ
// per family (Z-Image is 4B/2560-dim, Anima is 0.6B, Qwen-Image has its own VAE)
// and borrowing across them is what produced a load that succeeded while the
// output was noise. Every such file is named explicitly in Catalog instead.
var sharedRoles = map[string][]string{
	"t5xxl":  {"flux"},
	"clip_l": {"flux"},
}

// Downloadable lists the missing or mismatched items that can be fetched.
func (r Report) Downloadable() []Item {
	var out []Item
	for _, it := range r.Items {
		if (it.State == StateMissing || it.State == StateMismatch) && it.URL != "" {
			out = append(out, it)
		}
	}
	return out
}

// TotalSizeMB is the approximate size of everything that still needs fetching.
func (r Report) TotalSizeMB() int {
	total := 0
	for _, it := range r.Downloadable() {
		total += it.SizeMB
	}
	return total
}

// Roles returns the role names in a stable order, for tests and logs.
func (r Report) Roles() []string {
	out := make([]string, 0, len(r.Items))
	for _, it := range r.Items {
		out = append(out, it.Role)
	}
	sort.Strings(out)
	return out
}
