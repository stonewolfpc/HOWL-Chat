package settings

import (
	"encoding/json"
	"os"
	"path/filepath"

	imagetypes "howl-chat/internal/image/types"
)

// ImageSettings is the persisted state for the image engine.
type ImageSettings struct {
	EngineConfig   imagetypes.EngineConfig      `json:"engine_config"`
	LastModel      string                       `json:"last_model"`
	LastModelOpts  imagetypes.ModelLoadOptions  `json:"last_model_opts"`
}

// Store handles reading and writing image settings to a JSON file.
type Store struct {
	path string
}

// NewStore creates a Store rooted at configDir/image-settings.json.
func NewStore(configDir string) *Store {
	return &Store{path: filepath.Join(configDir, "image-settings.json")}
}

// Load reads persisted settings from disk. Returns zero-value if the file does not exist yet.
func (s *Store) Load() (ImageSettings, error) {
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

// Save persists settings to disk atomically via a temp file.
func (s *Store) Save(cfg ImageSettings) error {
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
