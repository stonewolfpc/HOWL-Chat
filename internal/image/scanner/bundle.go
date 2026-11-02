package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Diffusers bundles are how Invoke AI stores a full pipeline: a
// model_index.json naming the pipeline class, plus component folders (unet or
// transformer, vae, text_encoder, text_encoder_2). Before this file the walker
// treated each of those folders as its own model and skipped the descriptor
// entirely, so an SDXL pipeline appeared as several nameless fragments and the
// real family was lost.

// loadModelIndexClass returns the Diffusers pipeline class declared by a folder,
// or "" when the folder is not a Diffusers bundle.
func loadModelIndexClass(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "model_index.json"))
	if err != nil {
		return ""
	}
	var mi struct {
		ClassName string `json:"_class_name"`
	}
	if json.Unmarshal(data, &mi) != nil {
		return ""
	}
	return mi.ClassName
}

// hasBundleShape reports whether a folder holds more than one weight file, the
// minimum shape of a real multi-file bundle. One weight and no companions is a
// single model that happens to live in a folder.
func hasBundleShape(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	weights := 0
	for _, e := range entries {
		if !e.IsDir() && IsModelFile(e.Name()) {
			weights++
		}
	}
	return weights >= 2
}

// componentWeight picks the weight file inside a Diffusers component folder.
func componentWeight(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	fallback := ""
	for _, e := range entries {
		if e.IsDir() || !IsModelFile(e.Name()) {
			continue
		}
		name := strings.ToLower(e.Name())
		if strings.Contains(name, "diffusion_pytorch_model") || name == "model.safetensors" {
			return filepath.Join(dir, e.Name())
		}
		if fallback == "" {
			fallback = filepath.Join(dir, e.Name())
		}
	}
	return fallback
}

// familyFromPipeline maps a Diffusers pipeline class to a model family.
func familyFromPipeline(class string) (ModelFamily, bool) {
	c := strings.ToLower(class)
	switch {
	case c == "":
		return FamilyUnknown, false
	case strings.Contains(c, "flux"):
		return FamilyFlux, true
	case strings.Contains(c, "chroma"):
		return FamilyChroma, true
	case strings.Contains(c, "cascade"):
		return FamilyCascade, true
	case strings.Contains(c, "stablediffusion3"), strings.Contains(c, "sd3"):
		if strings.Contains(c, "3.5") || strings.Contains(c, "3_5") {
			return FamilySD35, true
		}
		return FamilySD3, true
	case strings.Contains(c, "xl"):
		return FamilySDXL, true
	case strings.Contains(c, "stablediffusion2"):
		return FamilySD21, true
	case strings.Contains(c, "stablediffusion"):
		return FamilySD15, true
	case strings.Contains(c, "wan"):
		return FamilyWan, true
	}
	return FamilyUnknown, false
}

// scanDiffusersBundle catalogues a Diffusers pipeline as a single model whose
// Files map points at the components the engine loads by path.
func (s *Scanner) scanDiffusersBundle(dir, pipeline string) (AssetEntry, bool) {
	modelPath := componentWeight(filepath.Join(dir, "unet"))
	if modelPath == "" {
		modelPath = componentWeight(filepath.Join(dir, "transformer"))
	}
	if modelPath == "" {
		return AssetEntry{}, false
	}

	fileMap := map[string]string{"model": modelPath}
	if p := componentWeight(filepath.Join(dir, "vae")); p != "" {
		fileMap["vae"] = p
	}
	if p := componentWeight(filepath.Join(dir, "text_encoder")); p != "" {
		fileMap["clip_l"] = p
	}

	family, _ := familyFromPipeline(pipeline)
	probe := ProbeWeight(modelPath)

	// SD3/Flux pipelines ship T5 as the second encoder; SDXL ships CLIP-G.
	if second := componentWeight(filepath.Join(dir, "text_encoder_2")); second != "" {
		if family == FamilyFlux || family == FamilySD3 || family == FamilySD35 {
			fileMap["t5xxl"] = second
		} else {
			fileMap["clip_g"] = second
		}
	}

	// The pipeline class names the architecture; the weights are the tie-break.
	if f, ok := probe.Family(); ok {
		family = f
	}

	var totalMB float64
	for _, p := range fileMap {
		if fi, err := os.Stat(p); err == nil {
			totalMB += float64(fi.Size()) / (1024 * 1024)
		}
	}

	rec := buildRecommended(AssetDiffusion, family)
	m := &Manifest{
		Name:        DisplayName(filepath.Base(dir), dir, probe, family),
		Type:        AssetDiffusion,
		Family:      family,
		Format:      "diffusers",
		SizeMB:      totalMB,
		Files:       fileMap,
		Logic:       manifestLogicVersion,
		Recommended: rec,
	}
	_ = SaveDirManifest(dir, m)
	entry := manifestToEntry(m, dir)
	if entry == nil {
		return AssetEntry{}, false
	}
	return *entry, true
}
