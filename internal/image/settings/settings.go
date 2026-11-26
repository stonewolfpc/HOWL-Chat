package settings

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	imagetypes "howl-chat/internal/image/types"
)

// ImageSettings is the persisted state for the image engine.
type ImageSettings struct {
	EngineConfig       imagetypes.EngineConfig     `json:"engine_config"`
	LastModel          string                      `json:"last_model"`
	LastModelOpts      imagetypes.ModelLoadOptions `json:"last_model_opts"`
	Generation         map[string]interface{}      `json:"generation,omitempty"`
	GenerationRevision string                      `json:"generation_revision,omitempty"`
}

// GenerationSaveResult is the structured outcome of a generation-settings
// write. Callers never infer success from an error and never write a cache
// before Status is "saved" or "updated".
type GenerationSaveResult struct {
	Status          string                 `json:"status"`
	Revision        string                 `json:"revision,omitempty"`
	CurrentRevision string                 `json:"current_revision,omitempty"`
	Message         string                 `json:"message,omitempty"`
	Document        map[string]interface{} `json:"document,omitempty"`
}

// Store handles reading and writing image settings to a JSON file.
type Store struct {
	mu   sync.Mutex
	path string
}

// NewStore creates a Store rooted at configDir/image-settings.json.
func NewStore(configDir string) *Store {
	return &Store{path: filepath.Join(configDir, "image-settings.json")}
}

// Load reads persisted settings from disk. Returns zero-value if the file does not exist yet.
func (s *Store) Load() (ImageSettings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadLocked()
}

func (s *Store) loadLocked() (ImageSettings, error) {
	var out ImageSettings
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return out, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return out, err
	}
	return out, nil
}

// Save persists a complete settings record atomically.
func (s *Store) Save(cfg ImageSettings) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveLocked(cfg)
}

func (s *Store) saveLocked(cfg ImageSettings) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// SaveEngineSettings updates the engine/model half of the same record while
// preserving the generation half. It exists so engine path changes cannot drop
// a generation settings save.
func (s *Store) SaveEngineSettings(cfg imagetypes.EngineConfig, last string, opts imagetypes.ModelLoadOptions) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, err := s.loadLocked()
	if err != nil {
		return err
	}
	current.EngineConfig = cfg
	current.LastModel = last
	current.LastModelOpts = opts
	return s.saveLocked(current)
}

// LoadGeneration returns the generation document and its revision. An absent
// document is reported as an empty map with an empty revision, never invented
// defaults: the renderer applies its own sane defaults.
func (s *Store) LoadGeneration() (map[string]interface{}, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, err := s.loadLocked()
	if err != nil {
		return nil, "", err
	}
	document := cloneDocument(current.Generation)
	return document, effectiveRevision(current.Generation, current.GenerationRevision), nil
}

// SaveGeneration replaces the generation document after checking the revision
// the caller loaded. It preserves fields it does not understand.
func (s *Store) SaveGeneration(document map[string]interface{}, baseRevision string) (GenerationSaveResult, error) {
	result := GenerationSaveResult{Status: "write_failed"}
	if document == nil {
		result.Message = "generation settings are required"
		return result, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	current, err := s.loadLocked()
	if err != nil {
		return result, err
	}
	currentRevision := effectiveRevision(current.Generation, current.GenerationRevision)
	result.CurrentRevision = currentRevision
	exists := currentRevision != "" || len(current.Generation) > 0

	if exists && baseRevision == "" {
		result.Status = "revision_conflict"
		result.Message = "A base revision is required to update image generation settings."
		return result, nil
	}
	if !exists && baseRevision != "" {
		result.Status = "revision_conflict"
		result.Message = "The image generation settings have no saved revision."
		return result, nil
	}
	if exists && baseRevision != currentRevision {
		result.Status = "revision_conflict"
		result.Message = "Image generation settings changed since they were loaded. Reload before saving."
		return result, nil
	}

	merged := cloneDocument(current.Generation)
	for key, value := range document {
		if key == "revision" {
			continue
		}
		merged[key] = value
	}
	revision := generationRevision(merged)
	current.Generation = merged
	current.GenerationRevision = revision
	if err := s.saveLocked(current); err != nil {
		result.Message = err.Error()
		return result, err
	}
	result.Status = "updated"
	if !exists {
		result.Status = "saved"
	}
	result.Revision = revision
	result.Document = cloneDocument(merged)
	return result, nil
}

// ClearGeneration removes the generation document. It is used only by the
// explicit rollback of the image generation cutover.
func (s *Store) ClearGeneration() (GenerationSaveResult, error) {
	result := GenerationSaveResult{Status: "write_failed"}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, err := s.loadLocked()
	if err != nil {
		return result, err
	}
	result.CurrentRevision = effectiveRevision(current.Generation, current.GenerationRevision)
	current.Generation = nil
	current.GenerationRevision = ""
	if err := s.saveLocked(current); err != nil {
		result.Message = err.Error()
		return result, err
	}
	result.Status = "cleared"
	return result, nil
}

func effectiveRevision(document map[string]interface{}, stored string) string {
	if stored != "" {
		return stored
	}
	if len(document) > 0 {
		return generationRevision(document)
	}
	return ""
}

func generationRevision(document map[string]interface{}) string {
	data, _ := json.Marshal(document)
	sum := sha256.Sum256(data)
	return "rv1-" + hex.EncodeToString(sum[:8])
}

func cloneDocument(document map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for key, value := range document {
		out[key] = value
	}
	return out
}

// String is used only in diagnostics, never as a persistence format.
func (c GenerationSaveResult) String() string {
	if c.Status == "" {
		return fmt.Sprintf("generation save %s", c.Revision)
	}
	return fmt.Sprintf("generation save %s: %s", c.Status, c.Message)
}
