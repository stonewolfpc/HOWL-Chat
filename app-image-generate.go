package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/image/engine/sdcpp"
	"howl-chat/internal/image/metadata"
	imagetypes "howl-chat/internal/image/types"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// -----------------------------------------------------------------------
// Generation
// -----------------------------------------------------------------------

// StartImageJob queues an async generation job and returns the job ID.
//
// Engine initialization and model loading run inside the job, not here. When
// they ran in this call, a slow load blocked the RPC for as long as the model
// took to load, the UI could not tell loading from a hang, and the user could
// not cancel it.
func (a *App) StartImageJob(opts imagetypes.ImageOptions) (string, error) {
	if a.jobManager == nil {
		return "", fmt.Errorf("job manager not initialized")
	}
	id := a.jobManager.Start(opts, func(jobID string) error {
		return a.ensureEngineReady(opts.ModelName, jobID, opts)
	})
	return id, nil
}

// GetImageJobStatus returns the current status of a job by ID.
func (a *App) GetImageJobStatus(id string) (*imagetypes.ImageJob, error) {
	if a.jobManager == nil {
		return nil, fmt.Errorf("job manager not initialized")
	}
	j := a.jobManager.Status(id)
	if j == nil {
		return nil, fmt.Errorf("job not found: %s", id)
	}
	return j, nil
}

// CancelImageJob marks a job as cancelled.
//
// The id is required. Generation started directly through GenerateImage has no
// job identity, so there is nothing for a cancel to address; an empty id must
// never be read as "whatever is running".
func (a *App) CancelImageJob(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("CancelImageJob needs a job id: generation started directly with GenerateImage cannot be cancelled")
	}
	if a.jobManager == nil {
		return fmt.Errorf("job manager not initialized")
	}
	if !a.jobManager.Cancel(id) {
		return fmt.Errorf("job not found or already finished: %s", id)
	}
	return nil
}

// GenerateImage creates an image using the active model (lazy-loads if needed).
//
// This binding is synchronous and deliberately not cancellable: it carries no job
// identity, so CancelImageJob cannot address it. The cancellable path is
// StartImageJob; callers that need cancellation must use it instead of this one.
func (a *App) GenerateImage(opts imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	if err := a.ensureEngineReady(opts.ModelName, "", opts); err != nil {
		return nil, fmt.Errorf("engine ready failed: %w", err)
	}

	var result *imagetypes.ImageResult
	var genErr error

	switch opts.Mode {
	case "img2img", "upscale", "outpaint":
		result, genErr = a.imageEngine.GenerateImg2Img(opts)
	case "inpaint":
		result, genErr = a.imageEngine.GenerateInpaint(opts)
	case "video":
		result, genErr = a.imageEngine.GenerateVideo(opts)
	default:
		result, genErr = a.imageEngine.GenerateTxt2Img(opts)
	}
	if genErr != nil {
		a.logImageGenError(genErr, "app-image-generate.go")
		return nil, genErr
	}
	if a.gallerySvc != nil && result != nil && (result.ImageData != "" || len(result.Images) > 0) {
		if err := a.gallerySvc.Save(result); err != nil {
			fmt.Printf("Warning: failed to save to gallery: %v\n", err)
		}
	}
	return result, nil
}

// -----------------------------------------------------------------------
// Gallery
// -----------------------------------------------------------------------

// ListGalleryImages returns all persisted gallery entries.
func (a *App) ListGalleryImages() ([]imagetypes.ImageMetadata, error) {
	if a.gallerySvc == nil {
		return nil, fmt.Errorf("gallery service not initialized")
	}
	return a.gallerySvc.List()
}

// DeleteGalleryImage removes an image from the gallery.
func (a *App) DeleteGalleryImage(baseName string) error {
	if a.gallerySvc == nil {
		return fmt.Errorf("gallery service not initialized")
	}
	return a.gallerySvc.Delete(baseName)
}

// GetGalleryImage reads a gallery image file and returns base64 PNG data.
func (a *App) GetGalleryImage(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read image: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// getGalleryDir returns the gallery output directory.
func (a *App) getGalleryDir() string {
	if a.gallerySvc != nil {
		return a.gallerySvc.BaseDir()
	}
	return filepath.Join(filepath.Dir(os.Args[0]), "images", "gallery")
}

// -----------------------------------------------------------------------
// Image file helpers
// -----------------------------------------------------------------------

// SaveImageAs shows a save dialog and writes the base64 PNG image to the chosen path.
func (a *App) SaveImageAs(imageDataB64 string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(imageDataB64)
	if err != nil {
		return "", fmt.Errorf("decode image data: %w", err)
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            "Save Image As",
		DefaultDirectory: a.getGalleryDir(),
		DefaultFilename:  "generated-image.png",
		Filters: []runtime.FileFilter{
			{DisplayName: "PNG Images (*.png)", Pattern: "*.png"},
			{DisplayName: "JPEG Images (*.jpg)", Pattern: "*.jpg;*.jpeg"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("save dialog: %w", err)
	}
	if path == "" {
		return "", nil
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("write image: %w", err)
	}
	return path, nil
}

// decodeBase64Image accepts either bare base64 or a full data: URL. Canvas
// editors (pose, masks) produce "data:image/png;base64,..." while file pickers
// produce bare base64, and both must reach the engine as bytes.
func decodeBase64Image(encoded string) ([]byte, error) {
	trimmed := strings.TrimSpace(encoded)
	if idx := strings.Index(trimmed, ";base64,"); idx >= 0 {
		trimmed = trimmed[idx+len(";base64,"):]
	}
	return base64.StdEncoding.DecodeString(trimmed)
}

// SaveTempImage decodes a base64 image and writes it to the engine temp directory.
func (a *App) SaveTempImage(imageDataB64 string) (string, error) {
	data, err := decodeBase64Image(imageDataB64)
	if err != nil {
		return "", fmt.Errorf("decode image data: %w", err)
	}
	var tempDir string
	if e, ok := a.imageEngine.(*sdcpp.Engine); ok {
		tempDir = e.GetConfig().TempDir
	}
	if tempDir == "" {
		tempDir = filepath.Join(os.TempDir(), "howl-image-temp")
	}
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}
	timestamp := time.Now().UTC().Format("20060102_150405")
	path := filepath.Join(tempDir, fmt.Sprintf("init_%s.png", timestamp))
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("write temp image: %w", err)
	}
	return path, nil
}

// SaveEntityImage saves a base64 image to the entity-images directory.
func (a *App) SaveEntityImage(imageDataB64, entityType, entityID string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(imageDataB64)
	if err != nil {
		return "", fmt.Errorf("decode image data: %w", err)
	}
	safeType := filepath.Clean(entityType)
	if safeType == "." || safeType == "" {
		safeType = "unknown"
	}
	entityDir := filepath.Join(filepath.Dir(os.Args[0]), "images", "entities", safeType)
	if err := os.MkdirAll(entityDir, 0755); err != nil {
		return "", fmt.Errorf("create entity dir: %w", err)
	}
	safeID := filepath.Clean(entityID)
	if safeID == "." || safeID == "" {
		safeID = "new"
	}
	timestamp := time.Now().UTC().Format("20060102_150405")
	path := filepath.Join(entityDir, fmt.Sprintf("%s_%s.png", safeID, timestamp))
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("write entity image: %w", err)
	}
	return path, nil
}

// -----------------------------------------------------------------------
// Metadata Extraction
// -----------------------------------------------------------------------

// ExtractImageMetadata extracts generation parameters from an image file.
// Supports PNG (tEXt chunks) and JPEG (EXIF UserComment).
// Returns metadata with prompt, negative prompt, seed, CFG, sampler, etc.
func (a *App) ExtractImageMetadata(imagePath string) (*metadata.ExtractedMetadata, error) {
	return metadata.ExtractFromFile(imagePath)
}

// ExtractImageMetadataFromBase64 extracts metadata from a base64-encoded image.
func (a *App) ExtractImageMetadataFromBase64(imageDataB64 string) (*metadata.ExtractedMetadata, error) {
	return metadata.ExtractFromBase64(imageDataB64)
}
