package gallery

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/image/metadata"
	"howl-chat/internal/image/types"
)

// ErrOutsideGallery is returned when a requested path is not contained by the
// gallery root after lexical and symlink/reparse-point resolution.
var ErrOutsideGallery = errors.New("path is outside the gallery")

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

	baseName := fmt.Sprintf("img_%s_%d", time.Now().UTC().Format("20060102_150405"), time.Now().UnixNano()%1000000)
	imgPath := filepath.Join(s.baseDir, baseName+".png")
	meta.File = imgPath
	metaBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("encode metadata: %w", err)
	}
	if err := writeFileAtomic(imgPath, imgDataWithMeta, 0644); err != nil {
		return fmt.Errorf("write image: %w", err)
	}

	metaPath := filepath.Join(s.baseDir, baseName+".json")
	if err := writeFileAtomic(metaPath, metaBytes, 0644); err != nil {
		// Do not leave an image without its sidecar when the sidecar could not
		// be committed.
		_ = os.Remove(imgPath)
		return fmt.Errorf("write metadata: %w", err)
	}
	return nil
}

func writeFileAtomic(path string, data []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	cleanup := func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
	}
	if _, err := tmp.Write(data); err != nil {
		cleanup()
		return err
	}
	if err := tmp.Chmod(mode); err != nil {
		cleanup()
		return err
	}
	if err := tmp.Sync(); err != nil {
		cleanup()
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
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
		resolved, err := s.ResolveFile(path)
		if err != nil {
			continue
		}
		data, err := os.ReadFile(resolved)
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
		resolved, err := s.ResolveFile(c)
		if err != nil {
			continue
		}
		return resolved, true
	}
	return "", false
}

// ResolveFile validates a requested gallery path and returns the resolved path.
// Both lexical containment and symlink/reparse-point containment are checked.
func (s *Service) ResolveFile(filePath string) (string, error) {
	if strings.TrimSpace(filePath) == "" {
		return "", fmt.Errorf("gallery path is empty")
	}
	base, err := filepath.Abs(s.baseDir)
	if err != nil {
		return "", fmt.Errorf("resolve gallery root: %w", err)
	}
	target, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("resolve gallery path: %w", err)
	}
	baseResolved := base
	if resolved, err := filepath.EvalSymlinks(base); err == nil {
		baseResolved = resolved
	}
	targetResolved := target
	if resolved, err := filepath.EvalSymlinks(target); err == nil {
		targetResolved = resolved
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("resolve gallery path: %w", err)
	}
	if !isUnderDir(targetResolved, baseResolved) {
		return "", fmt.Errorf("%w: %s", ErrOutsideGallery, filePath)
	}
	info, err := os.Stat(targetResolved)
	if err != nil {
		return "", fmt.Errorf("read gallery image: %w", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("gallery path is a directory")
	}
	return targetResolved, nil
}

// Read returns the bytes of a file inside the gallery root.
func (s *Service) Read(filePath string) ([]byte, error) {
	resolved, err := s.ResolveFile(filePath)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(resolved)
}

func isUnderDir(path, dir string) bool {
	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return false
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return false
	}
	return !filepath.IsAbs(rel)
}

// Delete removes an image, its metadata, and its thumbnail. Missing sidecars
// are harmless, but a missing item or a real filesystem failure is reported.
func (s *Service) Delete(baseName string) error {
	baseName = strings.TrimSpace(baseName)
	for _, ext := range []string{".png", ".json", ".webp"} {
		if strings.EqualFold(filepath.Ext(baseName), ext) {
			baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
			break
		}
	}
	if baseName == "" || baseName == "." || baseName == ".." ||
		filepath.IsAbs(baseName) || filepath.Base(baseName) != baseName ||
		strings.ContainsAny(baseName, `/\`) {
		return fmt.Errorf("invalid gallery item name %q", baseName)
	}

	var errs []error
	removed := false
	for _, ext := range []string{".png", ".json"} {
		err := os.Remove(filepath.Join(s.baseDir, baseName+ext))
		if err == nil {
			removed = true
		} else if !os.IsNotExist(err) {
			errs = append(errs, err)
		}
	}
	err := os.Remove(filepath.Join(s.baseDir, "thumbs", baseName+".webp"))
	if err == nil {
		removed = true
	} else if !os.IsNotExist(err) {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("delete gallery item %q: %w", baseName, errors.Join(errs...))
	}
	if !removed {
		return fmt.Errorf("gallery item %q not found", baseName)
	}
	return nil
}
