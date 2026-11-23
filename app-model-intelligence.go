package main

import (
	"fmt"
	"os"

	"howl-chat/internal/image/registry"
	"howl-chat/internal/modelintel"
)

// GetModelCatalog returns one shared profile per registered model. The profile
// is cached by file size and modification time; it is not a claim that the
// runtime can load the model.
func (a *App) GetModelCatalog() ([]modelintel.Profile, error) {
	if a.registries.Models() == nil {
		return []modelintel.Profile{}, fmt.Errorf("model registry not initialised")
	}
	if a.registries.Intel() == nil {
		return []modelintel.Profile{}, fmt.Errorf("model intelligence cache not initialised")
	}
	entries := a.registries.Models().All()
	out := make([]modelintel.Profile, 0, len(entries))
	for _, entry := range entries {
		profile, err := a.modelProfile(entry)
		if err != nil {
			continue
		}
		out = append(out, profile)
	}
	return out, nil
}

// GetModelIntelligence returns the shared profile for one registered path.
func (a *App) GetModelIntelligence(path string) (modelintel.Profile, error) {
	if a.registries.Intel() == nil {
		return modelintel.Profile{}, fmt.Errorf("model intelligence cache not initialised")
	}
	if a.registries.Models() == nil {
		return modelintel.Profile{}, fmt.Errorf("model registry not initialised")
	}
	entry, ok := a.registries.Models().Lookup(path)
	if !ok {
		return modelintel.Profile{}, fmt.Errorf("model is not registered: %s", path)
	}
	return a.modelProfile(entry)
}

func (a *App) modelProfile(entry registry.Entry) (modelintel.Profile, error) {
	size, modNS := entry.FileSize, entry.ModTimeNS
	available := true
	if info, err := os.Stat(entry.Path); err != nil {
		available = false
		size, modNS = 0, 0
	} else {
		size, modNS = info.Size(), info.ModTime().UnixNano()
	}
	artifacts := modelintel.InspectArtifacts(entry.Path)
	if cached, ok := a.registries.Intel().Get(entry.Path, size, modNS, artifacts.Fingerprint); ok {
		cached.Available = available
		return cached, nil
	}
	profile := modelintel.ClassifyWithArtifacts(modelintel.Input{
		Path:      entry.Path,
		Name:      entry.Name,
		Type:      entry.Type,
		Family:    entry.Family,
		Format:    entry.Format,
		Arch:      entry.Arch,
		SizeBytes: size,
		ModTimeNS: modNS,
	}, artifacts)
	profile.Available = available
	if err := a.registries.Intel().Put(profile); err != nil {
		return profile, err
	}
	return profile, nil
}
