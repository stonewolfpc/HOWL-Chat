package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const indexFile = "model-index.json"

// AssetEntry is one record in the model index — covers all asset types.
type AssetEntry struct {
	Name    string            `json:"name"`
	Path    string            `json:"path"`
	Type    AssetType         `json:"type"`
	Family  ModelFamily       `json:"family"`
	Format  string            `json:"format"`
	SizeMB  float64           `json:"size_mb"`
	Files   map[string]string `json:"files,omitempty"`
	Recommended Recommended   `json:"recommended"`
}

// ModelIndex is the full catalog produced by a scan.
// Each slice contains only entries of the matching type.
type ModelIndex struct {
	Models      []AssetEntry `json:"models"`       // diffusion models only
	VAEs        []AssetEntry `json:"vaes"`
	LoRAs       []AssetEntry `json:"loras"`
	ControlNets []AssetEntry `json:"controlnets"`
	Embeddings  []AssetEntry `json:"embeddings"`
	Upscalers   []AssetEntry `json:"upscalers"`
}

// IndexStore handles reading and writing the cached model index.
type IndexStore struct {
	dir string // the models directory; index lives alongside it
}

// NewIndexStore creates a store rooted at modelsDir.
func NewIndexStore(modelsDir string) *IndexStore {
	return &IndexStore{dir: modelsDir}
}

// Load reads the cached index. Returns an empty index if the file is absent.
func (s *IndexStore) Load() (ModelIndex, error) {
	var idx ModelIndex
	data, err := os.ReadFile(s.indexPath())
	if err != nil {
		if os.IsNotExist(err) {
			return idx, nil
		}
		return idx, err
	}
	err = json.Unmarshal(data, &idx)
	return idx, err
}

// Save persists the index atomically.
func (s *IndexStore) Save(idx ModelIndex) error {
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.indexPath() + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.indexPath())
}

func (s *IndexStore) indexPath() string {
	return filepath.Join(s.dir, indexFile)
}

// FilterVAEsForFamily returns VAE entries compatible with the given model family.
func FilterVAEsForFamily(vaes []AssetEntry, family ModelFamily) []AssetEntry {
	var out []AssetEntry
	for _, v := range vaes {
		if CompatibleVAE(family, v.Family) {
			out = append(out, v)
		}
	}
	return out
}

// FilterLoRAsForFamily returns LoRA entries compatible with the given model family.
func FilterLoRAsForFamily(loras []AssetEntry, family ModelFamily) []AssetEntry {
	var out []AssetEntry
	for _, l := range loras {
		if CompatibleLoRA(family, l.Family) {
			out = append(out, l)
		}
	}
	return out
}

// FilterControlNetsForFamily returns ControlNet entries compatible with the model family.
func FilterControlNetsForFamily(cns []AssetEntry, family ModelFamily) []AssetEntry {
	var out []AssetEntry
	for _, c := range cns {
		if CompatibleControlNet(family, c.Family) {
			out = append(out, c)
		}
	}
	return out
}
