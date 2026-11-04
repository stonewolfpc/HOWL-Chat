package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const manifestSuffix = ".howl.json"

// manifestLogicVersion is bumped whenever classification improves. A sidecar
// written by an older classifier is re-read rather than trusted, so a fixed
// rule reaches files that were already scanned.
const manifestLogicVersion = 5

// Recommended holds generation defaults derived from the model family.
type Recommended struct {
	Sampler   string  `json:"sampler"`
	Scheduler string  `json:"scheduler"`
	Steps     int     `json:"steps"`
	CFG       float64 `json:"cfg"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
}

// Manifest is the .howl.json sidecar written next to each discovered asset.
// It is generated once and read on every subsequent scan — making rescans instant.
type Manifest struct {
	Name         string            `json:"name"`
	Type         AssetType         `json:"type"`
	Family       ModelFamily       `json:"family"`
	Format       string            `json:"format"` // gguf, safetensors, ckpt, …
	SizeMB       float64           `json:"size_mb"`
	FileSize     int64             `json:"file_size,omitempty"`
	ModTimeNS    int64             `json:"mod_time_ns,omitempty"`
	Hash         string            `json:"hash,omitempty"`              // SHA256 of the primary model file
	CivitaiName  string            `json:"civitai_name,omitempty"`     // name from Civitai hash lookup, when metadata title is absent
	Logic        int               `json:"logic,omitempty"`
	Files        map[string]string `json:"files,omitempty"`             // sub-file role → abs path
	VAEPath      string            `json:"vae_path,omitempty"`          // suggested companion VAE
	Recommended  Recommended       `json:"recommended"`
}

// manifestPath returns the sidecar path for a given model file or directory.
func manifestPath(modelPath string) string {
	base := filepath.Base(modelPath)
	ext := filepath.Ext(base)
	stem := base[:len(base)-len(ext)]
	return filepath.Join(filepath.Dir(modelPath), stem+manifestSuffix)
}

// manifestPathForDir returns the sidecar path for a model directory bundle.
func manifestPathForDir(dir string) string {
	return filepath.Join(dir, "model"+manifestSuffix)
}

// LoadManifest reads the sidecar for a model file. Returns nil if absent.
func LoadManifest(modelPath string) *Manifest {
	data, err := os.ReadFile(manifestPath(modelPath))
	if err != nil {
		return nil
	}
	var m Manifest
	if json.Unmarshal(data, &m) != nil {
		return nil
	}
	return &m
}

// LoadDirManifest reads the sidecar for a model directory bundle. Returns nil if absent.
func LoadDirManifest(dir string) *Manifest {
	data, err := os.ReadFile(manifestPathForDir(dir))
	if err != nil {
		return nil
	}
	var m Manifest
	if json.Unmarshal(data, &m) != nil {
		return nil
	}
	return &m
}

// SaveManifest writes a manifest sidecar atomically via a temp file.
func SaveManifest(modelPath string, m *Manifest) error {
	return writeJSON(manifestPath(modelPath), m)
}

// SaveDirManifest writes a manifest sidecar for a directory bundle.
func SaveDirManifest(dir string, m *Manifest) error {
	return writeJSON(manifestPathForDir(dir), m)
}

// ClearManifests removes all .howl.json sidecar files under dir (recursively).
// Call this before a forced rescan so the classifier re-evaluates every file
// with the current logic instead of reusing stale cached results.
func ClearManifests(dir string) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), manifestSuffix) {
			_ = os.Remove(path)
		}
		return nil
	})
}

func writeJSON(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
