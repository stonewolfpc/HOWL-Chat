package sdcpp

import (
	"sort"
	"strings"
)

// FamilyCapability is a family's loading contract with sd.cpp: which flag
// carries the diffusion weights, which companion roles the architecture needs,
// any flags it always requires, and the defaults the UI needs to know.
//
// Before this table that knowledge was spread across five places —
// buildServerArgs, isModelLoadable, autoFillCompanions, the validator and the
// resolution switches — which had already drifted apart (HiDream was
// "self-contained" in one and "needs T5XXL + VAE" in another, and the UI's
// family lists disagreed with both). Adding an architecture meant editing all of
// them and hoping. Now it is one entry.
type FamilyCapability struct {
	// SplitComponents passes the diffusion model with --diffusion-model instead
	// of -m. Architectures whose encoders are not inside the checkpoint need it.
	SplitComponents bool
	// Companions are roles the architecture cannot run without.
	Companions []string
	// Optional roles are attached when the user has them, never required.
	Optional []string
	// ExtraFlags are architecture flags sd.cpp needs unconditionally.
	ExtraFlags []string
	// Resolution is the base square resolution the architecture is trained at.
	Resolution int
	// FluxAE marks families that decode with the Flux-style ae.safetensors, which
	// must never be attached to an SD checkpoint.
	FluxAE bool
}

// companionFlag maps a file role to the sd-server flag that accepts it.
var companionFlag = map[string]string{
	"vae":         "--vae",
	"clip_l":      "--clip_l",
	"clip_g":      "--clip_g",
	"t5xxl":       "--t5xxl",
	"llm":         "--llm",
	"clip_vision": "--clip_vision",
}

// familyCapabilities is the single source of truth for how each architecture is
// handed to sd.cpp. Only families this engine build implements belong here; the
// values come from sd-server --help for the shipped commit and from the
// companion requirements the app already enforced.
var familyCapabilities = map[string]FamilyCapability{
	// Checkpoints with encoders and VAE baked in.
	"sd1.5":      {Optional: []string{"clip_l"}, Resolution: 512},
	"sd2.1":      {Optional: []string{"clip_l"}, Resolution: 768},
	"sdxl":       {Optional: []string{"clip_l", "clip_g", "vae"}, Resolution: 1024},
	"sdxl-turbo": {Optional: []string{"clip_l", "clip_g", "vae"}, Resolution: 1024},
	"kolors":     {Optional: []string{"clip_l", "clip_g", "vae"}, Resolution: 1024, FluxAE: true},
	"cascade":    {Optional: []string{"clip_l", "clip_g", "vae"}, Resolution: 768},

	// Split-component families.
	"sd3":        {Companions: []string{"clip_l", "clip_g", "t5xxl", "vae"}, Resolution: 1024},
	"sd3.5":      {Companions: []string{"clip_l", "clip_g", "t5xxl", "vae"}, Resolution: 1024},
	"flux":       {SplitComponents: true, Companions: []string{"t5xxl", "clip_l", "vae"}, Resolution: 1024, FluxAE: true},
	"flux2":      {SplitComponents: true, Companions: []string{"llm", "vae"}, Resolution: 1024, FluxAE: true},
	"chroma":     {SplitComponents: true, Companions: []string{"t5xxl", "vae"}, ExtraFlags: []string{"--chroma-disable-dit-mask", "--clip-on-cpu"}, Resolution: 1024, FluxAE: true},
	"auraflow":   {SplitComponents: true, Companions: []string{"t5xxl", "vae"}, Resolution: 1024, FluxAE: true},
	"anima":      {SplitComponents: true, Companions: []string{"llm", "vae"}, Resolution: 1024},
	"z_image":    {SplitComponents: true, Companions: []string{"llm", "vae"}, Resolution: 1024, FluxAE: true},
	"qwen_image": {SplitComponents: true, Companions: []string{"llm", "vae"}, Resolution: 1024},

	// HiDream ships as a full model with external encoders. The engine used to
	// skip them while the discovery pass required them; the requirement wins.
	"hidream": {Companions: []string{"t5xxl", "vae"}, Resolution: 1024, FluxAE: true},
}

// DefaultCapability describes an architecture we have no entry for. It matches
// the historical behaviour exactly — single-file loading with whatever
// companions happen to be present — so an unlisted family cannot regress.
var DefaultCapability = FamilyCapability{
	Optional: []string{"vae", "clip_l", "clip_g"},
}

// Capability returns the loading contract for a family.
func Capability(family string) (FamilyCapability, bool) {
	c, ok := familyCapabilities[strings.ToLower(strings.TrimSpace(family))]
	return c, ok
}

// SupportsFamily reports whether this engine build implements the family.
func SupportsFamily(family string) bool {
	_, ok := Capability(family)
	return ok
}

// FamilyResolution returns the base square resolution for a family, or 0 when
// the family is unknown (the caller decides the fallback).
func FamilyResolution(family string) int {
	if c, ok := Capability(family); ok {
		return c.Resolution
	}
	return 0
}

// UsesFluxAE reports whether the family decodes with the Flux ae.safetensors.
func UsesFluxAE(family string) bool {
	if c, ok := Capability(family); ok {
		return c.FluxAE
	}
	return false
}

// Roles returns every role the family may use, required or optional.
func (c FamilyCapability) Roles() []string {
	out := make([]string, 0, len(c.Companions)+len(c.Optional))
	out = append(out, c.Companions...)
	out = append(out, c.Optional...)
	return out
}

// Uses reports whether the family ever attaches the given role.
func (c FamilyCapability) Uses(role string) bool {
	for _, r := range c.Companions {
		if r == role {
			return true
		}
	}
	for _, r := range c.Optional {
		if r == role {
			return true
		}
	}
	return false
}

// SupportedFamilies lists the implemented families for user-facing errors.
func SupportedFamilies() []string {
	out := make([]string, 0, len(familyCapabilities))
	for name := range familyCapabilities {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
