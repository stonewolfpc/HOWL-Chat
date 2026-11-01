package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Store handles reading and writing a unified settings JSON file.
type Store struct {
	path string
}

// NewStore creates a Store rooted at configDir/settings.json.
func NewStore(configDir string) *Store {
	return &Store{path: filepath.Join(configDir, "settings.json")}
}

// Load reads persisted settings from disk.
// Returns an empty map if the file does not exist or is unreadable.
func (s *Store) Load() (map[string]interface{}, error) {
	out := make(map[string]interface{})
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

// Save persists settings to disk with a direct write.
//
// The settings JSON is tiny (a few KB), so a direct os.WriteFile is
// instantaneous and eliminates the complexity of temp-file rename
// dances that cause Windows file-locking errors and leave stale
// .tmp / .bak files behind. A 150ms frontend debounce batches
// rapid changes so this only fires when the user pauses.
func (s *Store) Save(settings map[string]interface{}) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}
