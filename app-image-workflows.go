package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Image workflows are saved generation recipes: every control on the image page,
// the mode they belong to, the LoRA stack and an optional pose skeleton. They are
// plain JSON files in the user's config directory so they survive reinstalls and
// can be shared by copying one file.
//
// stable-diffusion.cpp has no workflow format of its own (it is not a graph
// runner), so a workflow here is a recipe the UI can restore exactly, not a
// node graph the engine executes.

const (
	workflowFormat      = "howl-workflow"
	workflowVersion     = 1
	workflowMaxDocBytes = 2 << 20 // 2 MiB is far more than any recipe needs
)

// ImageWorkflowInfo is the summary the picker shows; the full document stays on disk.
type ImageWorkflowInfo struct {
	Name     string `json:"name"`
	FileName string `json:"file_name"`
	Mode     string `json:"mode"`
	Model    string `json:"model"`
	SavedAt  string `json:"saved_at"`
	HasPose  bool   `json:"has_pose"`
}

type imageWorkflowDoc struct {
	Format   string          `json:"format"`
	Version  int             `json:"version"`
	Name     string          `json:"name"`
	SavedAt  string          `json:"saved_at"`
	Mode     string          `json:"mode"`
	Model    string          `json:"model"`
	Settings json.RawMessage `json:"settings"`
	LoRAs    json.RawMessage `json:"loras"`
	Pose     json.RawMessage `json:"pose"`
}

// imageWorkflowDir returns (and creates) the folder that holds workflow files.
func imageWorkflowDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot locate the user config directory: %w", err)
	}
	dir := filepath.Join(base, "HOWL_Chat", "workflows")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("cannot create the workflows folder: %w", err)
	}
	return dir, nil
}

// safeWorkflowName turns a display name into a file name. Anything that is not a
// letter, digit, space, dash, underscore or dot becomes a dash, so a name can
// never escape the workflows folder.
var unsafeWorkflowChars = regexp.MustCompile(`[^A-Za-z0-9 _.\-]+`)

func safeWorkflowName(name string) (string, error) {
	clean := strings.TrimSpace(unsafeWorkflowChars.ReplaceAllString(name, "-"))
	clean = strings.Trim(clean, " .-")
	if clean == "" {
		return "", errors.New("workflow name is empty after removing invalid characters")
	}
	if len(clean) > 80 {
		clean = clean[:80]
	}
	return clean + ".json", nil
}

// ListImageWorkflows returns the saved recipes, newest first.
func (a *App) ListImageWorkflows() ([]ImageWorkflowInfo, error) {
	dir, err := imageWorkflowDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := make([]ImageWorkflowInfo, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		var doc imageWorkflowDoc
		if err := json.Unmarshal(raw, &doc); err != nil {
			// A hand-edited or half-written file should not break the picker.
			out = append(out, ImageWorkflowInfo{
				Name:     strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())),
				FileName: entry.Name(),
			})
			continue
		}
		name := doc.Name
		if name == "" {
			name = strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		}
		out = append(out, ImageWorkflowInfo{
			Name:     name,
			FileName: entry.Name(),
			Mode:     doc.Mode,
			Model:    doc.Model,
			SavedAt:  doc.SavedAt,
			HasPose:  len(doc.Pose) > 0 && string(doc.Pose) != "null",
		})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].SavedAt > out[j].SavedAt })
	return out, nil
}

// SaveImageWorkflow writes a recipe document. Saving under an existing name
// replaces that recipe, which is what "Save" means in the UI.
func (a *App) SaveImageWorkflow(docJSON string) (ImageWorkflowInfo, error) {
	if len(docJSON) > workflowMaxDocBytes {
		return ImageWorkflowInfo{}, errors.New("workflow is too large to save")
	}
	var doc imageWorkflowDoc
	if err := json.Unmarshal([]byte(docJSON), &doc); err != nil {
		return ImageWorkflowInfo{}, fmt.Errorf("workflow is not valid JSON: %w", err)
	}
	if doc.Format != "" && doc.Format != workflowFormat {
		return ImageWorkflowInfo{}, fmt.Errorf("unsupported workflow format %q", doc.Format)
	}
	if strings.TrimSpace(doc.Name) == "" {
		return ImageWorkflowInfo{}, errors.New("workflow needs a name")
	}
	if doc.Version > workflowVersion {
		return ImageWorkflowInfo{}, fmt.Errorf("workflow was made by a newer HOWL Chat (version %d)", doc.Version)
	}

	fileName, err := safeWorkflowName(doc.Name)
	if err != nil {
		return ImageWorkflowInfo{}, err
	}
	dir, err := imageWorkflowDir()
	if err != nil {
		return ImageWorkflowInfo{}, err
	}

	// Always rewrite the tags we own so the summary stays trustworthy.
	doc.Format = workflowFormat
	doc.Version = workflowVersion
	if strings.TrimSpace(doc.SavedAt) == "" {
		doc.SavedAt = time.Now().UTC().Format(time.RFC3339)
	}
	encoded, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return ImageWorkflowInfo{}, err
	}

	path := filepath.Join(dir, fileName)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, encoded, 0o644); err != nil {
		return ImageWorkflowInfo{}, fmt.Errorf("cannot write the workflow: %w", err)
	}
	// Write-then-rename so an interrupted save cannot truncate an existing recipe.
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return ImageWorkflowInfo{}, fmt.Errorf("cannot finalise the workflow: %w", err)
	}
	return ImageWorkflowInfo{
		Name:     doc.Name,
		FileName: fileName,
		Mode:     doc.Mode,
		Model:    doc.Model,
		SavedAt:  doc.SavedAt,
		HasPose:  len(doc.Pose) > 0 && string(doc.Pose) != "null",
	}, nil
}

// GetImageWorkflow returns the raw document for a saved recipe.
func (a *App) GetImageWorkflow(fileName string) (string, error) {
	dir, err := imageWorkflowDir()
	if err != nil {
		return "", err
	}
	base := filepath.Base(fileName)
	if !strings.HasSuffix(strings.ToLower(base), ".json") {
		return "", errors.New("not a workflow file")
	}
	raw, err := os.ReadFile(filepath.Join(dir, base))
	if err != nil {
		return "", fmt.Errorf("cannot read the workflow: %w", err)
	}
	return string(raw), nil
}

// DeleteImageWorkflow removes a recipe. Only .json files directly inside the
// workflows folder can be deleted.
func (a *App) DeleteImageWorkflow(fileName string) error {
	dir, err := imageWorkflowDir()
	if err != nil {
		return err
	}
	base := filepath.Base(fileName)
	if !strings.HasSuffix(strings.ToLower(base), ".json") {
		return errors.New("not a workflow file")
	}
	path := filepath.Join(dir, base)
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("workflow not found: %w", err)
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("cannot delete the workflow: %w", err)
	}
	return nil
}

// OpenImageWorkflowsFolder shows the recipes in the file manager.
func (a *App) OpenImageWorkflowsFolder() error {
	dir, err := imageWorkflowDir()
	if err != nil {
		return err
	}
	cmd := exec.Command("explorer", dir)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cannot open the workflows folder: %w", err)
	}
	go cmd.Wait()
	return nil
}
