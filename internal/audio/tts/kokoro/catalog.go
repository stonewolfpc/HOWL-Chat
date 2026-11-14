package kokoro

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// ListVoicesFromArchive reads voice names directly from the Kokoro voices ZIP.
// It does not start Python or load the ONNX model, so the Character tab can
// show the real catalog before synthesis is attempted.
func ListVoicesFromArchive(path string) ([]string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	seen := make(map[string]bool)
	voices := make([]string, 0)
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		name := filepath.Base(f.Name)
		if !strings.HasSuffix(strings.ToLower(name), ".npy") {
			continue
		}
		voice := strings.TrimSuffix(name, filepath.Ext(name))
		if voice == "" || seen[voice] {
			continue
		}
		seen[voice] = true
		voices = append(voices, voice)
	}
	if len(voices) == 0 {
		return nil, fmt.Errorf("voices archive contains no .npy entries")
	}
	sort.Strings(voices)
	return voices, nil
}
