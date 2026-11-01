package gallery

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/image/metadata"
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
// When batch_count > 1, each image is saved as a separate gallery entry.
func (s *Service) Save(result *types.ImageResult) error {
	if result == nil {
		return fmt.Errorf("nil result")
	}

	var imagesToSave []string
	if len(result.Images) > 0 {
		imagesToSave = result.Images
	} else if result.ImageData != "" {
		imagesToSave = []string{result.ImageData}
	}
	if len(imagesToSave) == 0 {
		return fmt.Errorf("no image data")
	}

	// If batch with per-image metadata, save each as a separate entry
	if len(imagesToSave) > 1 && len(result.ImageMetadatas) == len(imagesToSave) {
		for i, b64 := range imagesToSave {
			meta := result.ImageMetadatas[i]
			if err := s.saveSingle(b64, meta); err != nil {
				return fmt.Errorf("save batch image %d: %w", i, err)
			}
		}
		// Set primary path from first image for the caller
		firstMeta := result.ImageMetadatas[0]
		result.Path = firstMeta.File
		result.Metadata = firstMeta
		return nil
	}

	// Single image or batch without per-image metadata: legacy behavior
	return s.saveSingle(imagesToSave[0], result.Metadata)
}

func (s *Service) saveSingle(b64 string, meta types.ImageMetadata) error {
	imgData, err := decodeBase64Image(b64)
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	// Build generation parameters for metadata embedding
	params := metadata.MetadataParams{
		Prompt:         meta.Prompt,
		NegativePrompt: meta.NegativePrompt,
		Steps:          meta.Steps,
		Sampler:        meta.Sampler,
		CFGScale:       meta.CfgScale,
		Seed:           meta.Seed,
		Width:          meta.Width,
		Height:         meta.Height,
	}

	// Embed metadata into the PNG image
	imgDataWithMeta, err := metadata.EmbedMetadata(imgData, params)
	if err != nil {
		// If embedding fails, fall back to original image (don't fail the save)
		imgDataWithMeta = imgData
	}

	timestamp := time.Now().UTC().Format("20060102_150405_") + fmt.Sprintf("%d", time.Now().UnixNano()%1000)
	baseName := fmt.Sprintf("img_%s", timestamp)
	imgPath := filepath.Join(s.baseDir, baseName+".png")
	if err := os.WriteFile(imgPath, imgDataWithMeta, 0644); err != nil {
		return fmt.Errorf("write image: %w", err)
	}

	meta.File = imgPath
	metaPath := filepath.Join(s.baseDir, baseName+".json")
	metaBytes, _ := json.MarshalIndent(meta, "", "  ")
	if err := os.WriteFile(metaPath, metaBytes, 0644); err != nil {
		return fmt.Errorf("write metadata: %w", err)
	}

	thumbDir := filepath.Join(s.baseDir, "thumbs")
	_ = os.MkdirAll(thumbDir, 0755)
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
		imgPath, ok := s.resolveImagePath(entry.Name(), meta.File)
		if !ok {
			continue
		}
		meta.File = imgPath
		results = append(results, meta)
	}
	return results, nil
}

// resolveImagePath finds the image that belongs to a metadata sidecar.
//
// Sidecars record the absolute path captured when the image was generated, so
// that value goes stale as soon as the project folder moves (drive letter
// change or a new install location). The gallery then hands the frontend paths
// it can no longer read, which renders as blank cards. The sidecar's own
// filename is authoritative, so prefer the image sitting next to it and only
// fall back to the base name of the recorded path.
func (s *Service) resolveImagePath(metaName, recorded string) (string, bool) {
	base := strings.TrimSuffix(filepath.Base(metaName), ".json")
	candidates := []string{filepath.Join(s.baseDir, base+".png")}
	if recorded != "" {
		candidates = append(candidates, filepath.Join(s.baseDir, filepath.Base(recorded)))
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && !info.IsDir() {
			return c, true
		}
	}
	return "", false
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
