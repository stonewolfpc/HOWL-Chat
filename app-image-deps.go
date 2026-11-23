package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/errors"
	"howl-chat/internal/image/deps"
	"howl-chat/internal/image/registry"
	imagetypes "howl-chat/internal/image/types"
)

// -----------------------------------------------------------------------
// Model dependencies
//
// A diffusion file on its own is often not a loadable model: flux2 needs a
// Mistral text encoder, z_image a Qwen3 one, every split-component family a
// VAE. Until now the app said "unsupported: missing llm" and stopped, leaving
// the user to work out which file, how big, and from where.
//
// These bindings answer that question from the same capability table the loader
// uses, and can fetch the missing files — but only after the UI has shown the
// consent text, and never as part of a load attempt on its own.
// -----------------------------------------------------------------------

// DependencyItem describes one companion requirement for the UI.
type DependencyItem struct {
	Role     string `json:"role"`
	Label    string `json:"label"`
	State    string `json:"state"` // present | missing | optional
	Path     string `json:"path,omitempty"`
	FileName string `json:"fileName,omitempty"`
	URL      string `json:"url,omitempty"`
	Source   string `json:"source,omitempty"`
	SizeMB   int    `json:"sizeMB,omitempty"`
	Note     string `json:"note,omitempty"`
	// Mismatch explains why a file that was found cannot be used for this role.
	Mismatch string `json:"mismatch,omitempty"`
	// Error is set when a download was attempted for this item and failed. The
	// message is written for a person: what failed and what to do instead.
	Error string `json:"error,omitempty"`
}

// DependencyReport is everything the UI needs to explain a model's needs.
type DependencyReport struct {
	Model       string           `json:"model"`
	Family      string           `json:"family"`
	Known       bool             `json:"known"`
	Items       []DependencyItem `json:"items"`
	Missing     []string         `json:"missing"`
	DownloadDir string           `json:"downloadDir"`
	TotalSizeMB int              `json:"totalSizeMB"`
	CanDownload bool             `json:"canDownload"`
	// Offline is true when the last download attempt failed because the host
	// was unreachable, so the UI can say so instead of blaming the file.
	Offline bool `json:"offline"`
	// Message carries a one-line explanation when nothing could be fetched.
	Message string `json:"message,omitempty"`
	// ConsentText is the sentence the UI must show before going online.
	ConsentText string `json:"consentText"`
}

// companionDir is where fetched companions are stored: inside the configured
// models directory so the engine and the scanner see them, never in a temp dir
// that a restart would lose.
func (a *App) companionDir() string {
	cfg := a.GetImageEngineConfig()
	if dir := strings.TrimSpace(cfg.ModelsDir); dir != "" {
		return filepath.Join(dir, "companions")
	}
	for _, dir := range a.imageDirectories() {
		if dir == "" {
			continue
		}
		return filepath.Join(dir, "companions")
	}
	return ""
}

// presentCompanions returns the role → path map the loader would use right now.
// It runs the same autofill the load path runs, so the report can never claim a
// file is missing while a load would find it (or the reverse).
func (a *App) presentCompanions(modelName string) map[string]string {
	opts := imagetypes.ModelLoadOptions{ModelName: modelName}
	if e := a.modelEntryByName(modelName); e != nil {
		opts.ModelPath = e.Path
	}
	a.autofillCompanionModels(&opts)
	out := map[string]string{}
	for role, p := range map[string]string{
		"vae":        opts.VAEPath,
		"clip_l":     opts.ClipLPath,
		"clip_g":     opts.ClipGPath,
		"t5xxl":      opts.T5XXLPath,
		"llm":        opts.LLMPath,
		"diffusion_": opts.DiffusionModel,
	} {
		if strings.TrimSpace(p) != "" {
			out[role] = p
		}
	}
	return out
}

// findCompatibleCompanion returns a registered image asset of this role that the
// catalogue accepts for the family. A file of the right kind but the wrong model
// (a Llama GGUF playing Qwen3, a LoRA playing a text encoder, a Qwen-Image VAE
// under Flux) is skipped rather than attached: attaching it does not produce a
// missing-companion error, it produces a load that fails deep inside sd-server
// with a message that names nothing useful.
func (a *App) findCompatibleCompanion(family, role string) string {
	if a.registries.Models() == nil {
		return ""
	}
	bestPath := ""
	bestRank := 3
	for _, e := range a.registries.Models().All() {
		if e.Path == "" || e.Type != role {
			continue
		}
		if a.registries.Models().KindForPath(e.Path) != registry.KindImage {
			continue
		}
		if ok, _ := deps.Compatible(family, role, filepath.Base(e.Path)); !ok {
			continue
		}
		if !companionFamilyUsable(family, e.Family) {
			continue
		}
		// Same architecture first, then a file the scanner could not classify,
		// then a family that legitimately shares the autoencoder. Anything owned
		// by a different architecture is skipped: an sd1.5 VAE under Anima is not
		// a missing-companion error, it is a load that fails inside sd-server.
		f := strings.ToLower(strings.TrimSpace(e.Family))
		rank := 2
		switch {
		case f == strings.ToLower(strings.TrimSpace(family)):
			rank = 0
		case f == "" || f == "unknown":
			rank = 1
		}
		// Ties break on path: Registry.All() walks a map, and that iteration
		// order is how the same model picked ae.safetensors on one run and an
		// sd1.5 model.safetensors on the next.
		if rank < bestRank || (rank == bestRank && (bestPath == "" || e.Path < bestPath)) {
			bestPath, bestRank = e.Path, rank
		}
	}
	return bestPath
}

// dependencyReport builds the report for one model.
func (a *App) dependencyReport(modelName string) (deps.Report, error) {
	e := a.modelEntryByName(modelName)
	if e == nil {
		return deps.Report{}, fmt.Errorf("'%s' is not a model in the image registry. Rescan the models folder first.", modelName)
	}
	return deps.Resolve(modelName, e.Family, a.presentCompanions(modelName), a.companionDir()), nil
}

func toDependencyReport(r deps.Report, items []DependencyItem) DependencyReport {
	if items == nil {
		items = make([]DependencyItem, 0, len(r.Items))
		for _, it := range r.Items {
			items = append(items, DependencyItem{
				Role: it.Role, Label: it.Label, State: it.State, Path: it.Path,
				FileName: it.FileName, URL: it.URL, Source: it.Source,
				SizeMB: it.SizeMB, Note: it.Note, Mismatch: it.Mismatch,
			})
		}
	}
	missing := r.Missing
	if missing == nil {
		missing = []string{}
	}
	return DependencyReport{
		Model:       r.Model,
		Family:      r.Family,
		Known:       r.Known,
		Items:       items,
		Missing:     missing,
		DownloadDir: r.DownloadDir,
		TotalSizeMB: r.TotalSizeMB(),
		CanDownload: len(r.Downloadable()) > 0,
		ConsentText: deps.ConsentText,
	}
}

// GetModelDependencies reports what a model needs and what is already on disk.
// It performs no network access; it only reads the registry.
func (a *App) GetModelDependencies(modelName string) (DependencyReport, error) {
	r, err := a.dependencyReport(modelName)
	if err != nil {
		return DependencyReport{}, err
	}
	return toDependencyReport(r, nil), nil
}

// DownloadModelDependencies fetches every missing companion that has a known
// download location. The caller must have shown ConsentText and received an
// explicit yes first: this method never runs on its own.
//
// It returns a report rather than an error for individual failures so the UI
// can list what worked, what did not, and the manual link for each — a failed
// download must leave the rest of the app usable.
func (a *App) DownloadModelDependencies(modelName string) (DependencyReport, error) {
	before, err := a.dependencyReport(modelName)
	if err != nil {
		return DependencyReport{}, err
	}
	if before.DownloadDir == "" {
		return DependencyReport{}, fmt.Errorf("no models folder is configured, so there is nowhere to save the files")
	}
	if err := os.MkdirAll(before.DownloadDir, 0o755); err != nil {
		return DependencyReport{}, fmt.Errorf("cannot create %s: %w", before.DownloadDir, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Hour)
	defer cancel()

	itemErrors := map[string]string{}
	offline := false
	saved := 0

	for _, it := range before.Downloadable() {
		a.emitEvent("image-deps:progress", map[string]interface{}{
			"role": it.Role, "file": it.FileName, "state": "start", "received": int64(0), "total": int64(0),
		})
		lastEmit := time.Time{}
		path, derr := deps.Download(ctx, nil, it.Artifact, before.DownloadDir, func(name string, got, total int64) {
			// One event per 4 MB keeps the UI smooth without flooding the bridge.
			if time.Since(lastEmit) < 250*time.Millisecond {
				return
			}
			lastEmit = time.Now()
			a.emitEvent("image-deps:progress", map[string]interface{}{
				"role": it.Role, "file": name, "state": "progress", "received": got, "total": total,
			})
		})
		if derr != nil {
			msg := derr.Error()
			if deps.IsOffline(derr) {
				offline = true
				msg = "No internet connection, so this file could not be downloaded. Use the link to get it manually, then rescan."
			}
			itemErrors[it.Role] = msg
			a.logErrorSilent(errors.REGModelAddFail, errors.Medium, "image.deps.download", "app-image-deps.go", derr)
			a.emitEvent("image-deps:progress", map[string]interface{}{
				"role": it.Role, "file": it.FileName, "state": "error", "message": msg,
			})
			continue
		}
		saved++
		a.registerCompanion(path, it.Role, before.Family)
		a.emitEvent("image-deps:progress", map[string]interface{}{
			"role": it.Role, "file": it.FileName, "state": "done", "path": path,
		})
	}

	if saved > 0 {
		if _, err := a.scanDirectoryIntoRegistry(before.DownloadDir, registry.KindImage); err != nil {
			a.logErrorSilent(errors.IMGRegistryFail, errors.Medium, "image.deps.scan", "app-image-deps.go", err)
		}
	}

	after, err := a.dependencyReport(modelName)
	if err != nil {
		return DependencyReport{}, err
	}
	items := make([]DependencyItem, 0, len(after.Items))
	for _, it := range after.Items {
		item := DependencyItem{
			Role: it.Role, Label: it.Label, State: it.State, Path: it.Path,
			FileName: it.FileName, URL: it.URL, Source: it.Source,
			SizeMB: it.SizeMB, Note: it.Note, Mismatch: it.Mismatch,
		}
		if msg, ok := itemErrors[it.Role]; ok && it.State != deps.StatePresent {
			item.Error = msg
		}
		items = append(items, item)
	}

	out := toDependencyReport(after, items)
	out.Offline = offline
	switch {
	case offline:
		out.Message = "No internet connection. The files that could not be fetched are listed with manual links."
	case len(after.Missing) > 0 && saved > 0:
		out.Message = "Some files downloaded, but the model still needs: " + strings.Join(after.Missing, ", ") + "."
	case len(after.Missing) > 0:
		out.Message = "Nothing could be downloaded automatically. Use the links, then rescan."
	case saved > 0:
		out.Message = "All required files are in place. Loading the model now."
	}
	return out, nil
}

// registerCompanion records a downloaded file under the role it was fetched
// for. The scanner classifies by content, and a bare text-encoder safetensors
// has no architecture marker to key on, so the role is taken from the catalogue
// that requested it rather than guessed a second time.
func (a *App) registerCompanion(path, role, family string) {
	if a.registries.Models() == nil || path == "" {
		return
	}
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	entry := registry.Entry{
		Name:     strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		Path:     path,
		Type:     role,
		Family:   family,
		Format:   strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), "."),
		SizeMB:   float64(info.Size()) / (1024 * 1024),
		FileSize: info.Size(),
	}
	if err := a.registries.Models().MergeEntries([]registry.Entry{entry}); err != nil {
		a.logErrorSilent(errors.REGModelAddFail, errors.Medium, "registry.companion", "app-image-deps.go", err)
	}
}
