package gallery

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/image/types"
)

// Service manages generated images on disk.
type Service struct {
	baseDir string
}

// NewService creates a gallery service.
func NewService(baseDir string) *Service {
	return &Service{baseDir: baseDir}
}

// BaseDir returns the gallery's base directory.
func (s *Service) BaseDir() string {
	return s.baseDir
}

// Save stores a generated image and its metadata.
// When batch_count > 1, all images in result.Images are saved.
func (s *Service) Save(result *types.ImageResult) error {
	if result == nil {
		return fmt.Errorf("nil result")
	}

	// Determine which images to save
	var imagesToSave []string
	if len(result.Images) > 0 {
		imagesToSave = result.Images
	} else if result.ImageData != "" {
		imagesToSave = []string{result.ImageData}
	}
	if len(imagesToSave) == 0 {
		return fmt.Errorf("no image data")
	}

	timestamp := time.Now().UTC().Format("20060102_150405")
	baseName := fmt.Sprintf("img_%s", timestamp)

	// Save all images (first one becomes primary)
	var primaryPath string
	for i, b64 := range imagesToSave {
		imgData, err := decodeBase64Image(b64)
		if err != nil {
			return fmt.Errorf("decode image %d: %w", i, err)
		}
		var imgPath string
		if i == 0 {
			imgPath = filepath.Join(s.baseDir, baseName+".png")
			primaryPath = imgPath
		} else {
			imgPath = filepath.Join(s.baseDir, fmt.Sprintf("%s_%d.png", baseName, i))
		}
		if err := os.WriteFile(imgPath, imgData, 0644); err != nil {
			return fmt.Errorf("write image %d: %w", i, err)
		}
	}
	result.Path = primaryPath
	result.Metadata.File = primaryPath

	// Write metadata JSON
	metaPath := filepath.Join(s.baseDir, baseName+".json")
	metaBytes, _ := json.MarshalIndent(result.Metadata, "", "  ")
	if err := os.WriteFile(metaPath, metaBytes, 0644); err != nil {
		return fmt.Errorf("write metadata: %w", err)
	}

	// Create thumbnail directory
	thumbPath := filepath.Join(s.baseDir, "thumbs", baseName+".webp")
	_ = os.MkdirAll(filepath.Dir(thumbPath), 0755)
	// TODO: generate actual thumbnail
	result.ThumbPath = thumbPath

	return nil
}

// decodeBase64Image extracts bytes from a base64 data URL or raw base64 string.
func decodeBase64Image(data string) ([]byte, error) {
	if idx := strings.Index(data, ","); idx != -1 {
		data = data[idx+1:]
	}
	return base64.StdEncoding.DecodeString(data)
}

// List returns all gallery entries sorted by newest first.
func (s *Service) List() ([]types.ImageMetadata, error) {
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var results []types.ImageMetadata
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		path := filepath.Join(s.baseDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var meta types.ImageMetadata
		if err := json.Unmarshal(data, &meta); err != nil {
			continue
		}
		results = append(results, meta)
	}
	return results, nil
}

// Delete removes an image and its metadata.
func (s *Service) Delete(baseName string) error {
	for _, ext := range []string{".png", ".json"} {
		p := filepath.Join(s.baseDir, baseName+ext)
		_ = os.Remove(p)
	}
	_ = os.Remove(filepath.Join(s.baseDir, "thumbs", baseName+".webp"))
	return nil
}
