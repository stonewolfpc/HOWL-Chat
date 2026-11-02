package sdcpp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/image/types"
)

// sdcppImageList accepts both legacy []string and newer []object image payloads.
type sdcppImageList []string

func (l *sdcppImageList) UnmarshalJSON(data []byte) error {
	// Legacy shape: ["<base64>", ...]
	var asStrings []string
	if err := json.Unmarshal(data, &asStrings); err == nil {
		*l = asStrings
		return nil
	}

	// Newer shape: [{"data":"<base64>"}, {"image":"<base64>"}, ...]
	var asObjects []map[string]any
	if err := json.Unmarshal(data, &asObjects); err == nil {
		out := make([]string, 0, len(asObjects))
		for _, obj := range asObjects {
			if obj == nil {
				continue
			}
			for _, key := range []string{"data", "image", "base64", "b64"} {
				if v, ok := obj[key].(string); ok && v != "" {
					out = append(out, v)
					break
				}
			}
		}
		*l = out
		return nil
	}

	// Single-string fallback: "<base64>"
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		if single == "" {
			*l = nil
		} else {
			*l = []string{single}
		}
		return nil
	}

	return fmt.Errorf("unsupported images payload shape")
}

func extractImageStrings(v any) []string {
	out := []string{}
	var walk func(any)
	walk = func(x any) {
		switch t := x.(type) {
		case string:
			if t != "" {
				out = append(out, t)
			}
		case []any:
			for _, it := range t {
				walk(it)
			}
		case map[string]any:
			// Prefer explicit image-like keys first.
			for _, k := range []string{"data", "image", "base64", "b64", "images", "output", "outputs", "result"} {
				if v2, ok := t[k]; ok {
					walk(v2)
				}
			}
			// Then scan remaining values (for future schema changes).
			for _, v2 := range t {
				walk(v2)
			}
		}
	}
	walk(v)
	return out
}

func normalizeImagePayloads(values []string) []string {
	out := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))

	// Known non-image metadata keywords to filter out
	metadataKeywords := map[string]bool{
		"png": true, "jpg": true, "jpeg": true, "webp": true,
		"completed": true, "failed": true, "cancelled": true,
		"generating": true, "queued": true,
		"img_gen": true, "txt2img": true, "img2img": true,
	}

	for _, s := range values {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		// Filter out known metadata keywords
		if metadataKeywords[s] {
			continue
		}

		// Filter out job IDs (start with "job_" and are short)
		if strings.HasPrefix(s, "job_") && len(s) < 50 {
			continue
		}

		// Filter out very short strings (not real images)
		if len(s) < 1000 {
			continue
		}

		// Convert file path payloads to base64 so UI can render consistently.
		if strings.HasPrefix(s, "/") || strings.Contains(s, `:\`) {
			if data, err := os.ReadFile(s); err == nil && len(data) > 0 {
				s = base64.StdEncoding.EncodeToString(data)
			}
		}

		// Accept raw base64 or data URLs.
		if strings.HasPrefix(s, "data:image/") {
			if idx := strings.Index(s, ","); idx >= 0 && idx+1 < len(s) {
				s = s[idx+1:]
			}
		}

		if _, dup := seen[s]; dup {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func recoverImagesFromJobBody(body []byte, preview string) []string {
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}

	fmt.Printf("[DEBUG] recoverImagesFromJobBody() raw keys: %v\n", getKeys(raw))

	values := extractImageStrings(raw)
	fmt.Printf("[DEBUG] recoverImagesFromJobBody() extracted %d strings\n", len(values))
	for i, v := range values {
		fmt.Printf("[DEBUG]   [%d] len=%d, prefix=%s\n", i, len(v), safePrefix(v))
	}

	values = normalizeImagePayloads(values)
	if len(values) > 0 {
		return values
	}
	if strings.TrimSpace(preview) != "" {
		return []string{preview}
	}
	return nil
}

func getKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func safePrefix(s string) string {
	if len(s) > 20 {
		return s[:20]
	}
	return s
}

// controlTimeout bounds the requests that ask the engine what it is doing.
// Generation itself stays unbounded, because a 1024x1024 run legitimately takes
// minutes — but a status probe that never returns is how the job poller lost
// contact with a dead engine and left the UI spinning on "Generating...".
const controlTimeout = 3 * time.Second

// client talks to sd-server over HTTP.
type client struct {
	baseURL   string
	http      *http.Client
	httpShort *http.Client
	mu        sync.Mutex
	// activeJobID is the async job the progress poller should query, and
	// activeJobToken is the ownership token of whoever set it. A request that
	// finishes late can only clear the id it actually owns.
	activeJobID    string
	activeJobToken uint64
	// pollTimeout bounds how long an async job may stay unfinished before the
	// request gives up. It is a field, not a constant, so tests can exercise the
	// timeout path without waiting ten minutes.
	pollTimeout time.Duration
}

func newClient(baseURL string) *client {
	return &client{
		baseURL:     baseURL,
		http:        &http.Client{Timeout: 0}, // No timeout for long generations
		httpShort:   &http.Client{Timeout: controlTimeout},
		pollTimeout: 600 * time.Second,
	}
}

// resetState clears any stale job IDs and cancellation flags.
func (c *client) resetState() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.activeJobID = ""
	c.activeJobToken++
}

// setActiveJob records the async job this request owns for progress purposes and
// returns its ownership token.
func (c *client) setActiveJob(id string) uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.activeJobID = id
	c.activeJobToken++
	return c.activeJobToken
}

// resetActiveJob clears the id the progress poller queries, but only when the
// token still belongs to the request doing the clearing. Without the token an
// abandoned request unwinding after a newer one had claimed progress would wipe
// the newer job's id.
func (c *client) resetActiveJob(token uint64) {
	c.mu.Lock()
	if c.activeJobToken == token {
		c.activeJobID = ""
	}
	c.mu.Unlock()
}

// activeJob reports the async job progress should be read from, if any.
func (c *client) activeJob() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.activeJobID
}

// generateSequentialBatch splits a multi-image batch into sequential
// single-image requests to avoid CUDA OOM. sd-server allocates all
// batch latents simultaneously in a true batch, which exceeds VRAM.
// Each image is saved immediately upon completion.
func (c *client) generateSequentialBatch(ctx context.Context, singleGen func(context.Context, types.ImageOptions) (*types.ImageResult, error), opts types.ImageOptions) (*types.ImageResult, error) {
	var firstResult *types.ImageResult
	var firstSeed = opts.Seed
	var allImages []string

	fmt.Printf("[ImageEngine] Batch: generating %d images sequentially (initial seed=%d)\n", opts.BatchCount, firstSeed)

	for i := 0; i < opts.BatchCount; i++ {
		// The batch stops between items once its request is cancelled; items already
		// generated are kept.
		if ctx.Err() != nil {
			if firstResult == nil {
				return nil, ctx.Err()
			}
			fmt.Printf("[ImageEngine] Batch interrupted by cancellation at item %d/%d\n", i+1, opts.BatchCount)
			break
		}

		singleOpts := opts
		singleOpts.BatchCount = 1

		// Handle seed logic
		if firstSeed == -1 || firstSeed == 0 {
			// Random seed mode - generate new random seed for each image
			singleOpts.Seed = -1
		} else {
			// Fixed seed mode - increment seed per item for variety
			singleOpts.Seed = firstSeed + int64(i)
		}

		result, err := singleGen(ctx, singleOpts)
		if err != nil {
			if firstResult == nil {
				return nil, fmt.Errorf("batch item %d/%d failed: %w", i+1, opts.BatchCount, err)
			}
			fmt.Printf("[ImageEngine] Batch interrupted at item %d/%d: %v\n", i+1, opts.BatchCount, err)
			break
		}

		// Save the first result to return
		if firstResult == nil {
			firstResult = result
		}

		// Accumulate all images for batch display
		if result.ImageData != "" {
			allImages = append(allImages, result.ImageData)
			fmt.Printf("[ImageEngine] Batch: accumulated image %d/%d\n", i+1, opts.BatchCount)
		}
	}

	if firstResult == nil {
		return nil, fmt.Errorf("no images generated in batch of %d", opts.BatchCount)
	}

	// Set all images in the result for batch display
	firstResult.Images = allImages

	fmt.Printf("[ImageEngine] Batch complete: %d images generated\n", len(allImages))

	return firstResult, nil
}

// probe asks the engine once whether it answers, under a bounded timeout. A
// single request the caller can interleave with a process-liveness check, which
// a blocking readiness loop cannot do: it cannot notice the child has died.
func (c *client) probe() error {
	resp, err := c.httpShort.Get(c.baseURL + "/sdapi/v1/samplers")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *client) waitForReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if err := c.probe(); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("sd-server not ready after %v", timeout)
}

// LoRARequest represents a LoRA with strength for API requests.
// Matches the format expected by sd-server / A1111 compatible APIs.
type LoRARequest struct {
	Name     string  `json:"name"`
	Strength float64 `json:"strength"`
}

// convertLoRAs converts []types.LoRAStrength to []LoRARequest for API calls.
// Filters out disabled LoRAs and maps the strength values.
func convertLoRAs(loras []types.LoRAStrength) []LoRARequest {
	if len(loras) == 0 {
		return nil
	}
	result := make([]LoRARequest, 0, len(loras))
	for _, l := range loras {
		if l.Enabled && l.Name != "" {
			result = append(result, LoRARequest{
				Name:     l.Name,
				Strength: l.Strength,
			})
		}
	}
	return result
}

// txt2imgRequest mirrors the WebUI /sdapi/v1/txt2img payload.
// NOTE: Using synchronous API (/sdapi/v1/txt2img) because async job API (/sdcpp/v1/img_gen) has critical limitations:
// 1. Does not respect steps parameter (always uses default 20)
// 2. Does not return progress/preview data for live preview
// See memory: live-preview-limitation-may2026
type txt2imgRequest struct {
	Prompt                  string        `json:"prompt"`
	NegativePrompt          string        `json:"negative_prompt"`
	Width                   int           `json:"width"`
	Height                  int           `json:"height"`
	Steps                   int           `json:"steps"`
	Seed                    int64         `json:"seed"`
	SamplerName             string        `json:"sampler_name"`
	Preview                 string        `json:"preview,omitempty"`
	Scheduler               string        `json:"scheduler"`
	CfgScale                float64       `json:"cfg_scale"`
	Guidance                float64       `json:"guidance,omitempty"`
	SlgScale                float64       `json:"slg_scale,omitempty"`
	SkipLayerStart          float64       `json:"skip_layer_start,omitempty"`
	SkipLayerEnd            float64       `json:"skip_layer_end,omitempty"`
	Eta                     float64       `json:"eta,omitempty"`
	BatchSize               int           `json:"batch_size"`
	EnableHR                bool          `json:"enable_hr"`
	HRScale                 float64       `json:"hr_scale"`
	HRUpscaler              string        `json:"hr_upscaler"`
	DenoisingStrength       float64       `json:"denoising_strength"`
	ResizeMode              string        `json:"resize_mode,omitempty"`
	StyleStrength           float64       `json:"style_strength,omitempty"`
	PreserveColors          bool          `json:"preserve_colors,omitempty"`
	MatchHistogram          bool          `json:"match_histogram,omitempty"`
	LatentNoise             float64       `json:"latent_noise,omitempty"`
	MaskBlur                int           `json:"mask_blur,omitempty"`
	MaskInvert              bool          `json:"mask_invert,omitempty"`
	MaskedContent           string        `json:"masked_content,omitempty"`
	EdgeBehavior            string        `json:"edge_behavior,omitempty"`
	UpscaleMethod           string        `json:"upscale_method,omitempty"`
	UpscaleSharpen          string        `json:"upscale_sharpen,omitempty"`
	UpscaleAntiArtifact     string        `json:"upscale_anti_artifact,omitempty"`
	Tiling                  bool          `json:"tiling"`
	ClipSkip                int           `json:"clip_skip"`
	LoRAs                   []LoRARequest `json:"loras,omitempty"`
	ControlNetPath          string        `json:"control_net_path,omitempty"`
	ControlImage            string        `json:"control_image,omitempty"`
	ControlStrength         float64       `json:"control_strength,omitempty"`
	CacheMode               string        `json:"cache_mode,omitempty"`
	CacheOption             string        `json:"cache_option,omitempty"`
	VideoFrames             int           `json:"video_frames,omitempty"`
	FPS                     int           `json:"fps,omitempty"`
	EndImage                string        `json:"end_image,omitempty"`
	VideoMotionStrength     float64       `json:"video_motion_strength,omitempty"`
	VideoSeedMode           string        `json:"video_seed_mode,omitempty"`
	VideoDenoise            float64       `json:"video_denoise,omitempty"`
	VideoLoop               bool          `json:"video_loop,omitempty"`
	Prediction              string        `json:"prediction,omitempty"`
	RNG                     string        `json:"rng,omitempty"`
	SamplerRNG              string        `json:"sampler_rng,omitempty"`
	FlowShift               float64       `json:"flow_shift,omitempty"`
	VAETileOverlap          float64       `json:"vae_tile_overlap,omitempty"`
	DisableMetadata         bool          `json:"disable_metadata,omitempty"`
	CfgRescale              float64       `json:"cfg_rescale,omitempty"`
	FaceRestore             float64       `json:"face_restore,omitempty"`
	SChurn                  float64       `json:"s_churn,omitempty"`
	STmin                   float64       `json:"s_tmin,omitempty"`
	STmax                   float64       `json:"s_tmax,omitempty"`
	SNoise                  float64       `json:"s_noise,omitempty"`
	OutputFormat            string        `json:"output_format,omitempty"`
	HighNoiseSteps          int           `json:"high_noise_steps,omitempty"`
	HighNoiseCfgScale       float64       `json:"high_noise_cfg_scale,omitempty"`
	HighNoiseGuidance       float64       `json:"high_noise_guidance,omitempty"`
	HighNoiseSlgScale       float64       `json:"high_noise_slg_scale,omitempty"`
	HighNoiseEta            float64       `json:"high_noise_eta,omitempty"`
	HighNoiseSkipLayerStart float64       `json:"high_noise_skip_layer_start,omitempty"`
	HighNoiseSkipLayerEnd   float64       `json:"high_noise_skip_layer_end,omitempty"`
}

type txt2imgResponse struct {
	Images     []string `json:"images"`
	Parameters struct{} `json:"parameters"`
	Info       string   `json:"info"`
}

func (c *client) txt2img(ctx context.Context, opts types.ImageOptions) (*types.ImageResult, error) {
	if opts.BatchCount > 1 {
		return c.generateSequentialBatch(ctx, c.txt2img, opts)
	}
	var ctrlImgB64 string
	if opts.ControlImagePath != "" {
		if data, err := os.ReadFile(opts.ControlImagePath); err == nil {
			ctrlImgB64 = base64.StdEncoding.EncodeToString(data)
		}
	}
	var endImgB64 string
	if opts.EndImagePath != "" {
		if data, err := os.ReadFile(opts.EndImagePath); err == nil {
			endImgB64 = base64.StdEncoding.EncodeToString(data)
		}
	}

	// Flux/Chroma use 'guidance' instead of 'cfg_scale'
	// If guidance is 0, set default (3.5 for Chroma per docs)
	guidance := opts.Guidance
	prediction := opts.Prediction
	isFluxFamily := strings.Contains(strings.ToLower(opts.ModelName), "flux") ||
		strings.Contains(strings.ToLower(opts.ModelName), "chroma")

	// DiT-based models: Flux, Chroma, SD3, SD3.5, SDXL (some), Anima, Kolors, AuraFlow, HiDream
	// These support DiT-specific parameters (flow_shift, slg_scale, skip_layer, eta, vae_tile_overlap)
	// SD1.5, SD2.1, and most SDXL models are NOT DiT-based and should not receive these parameters
	isDiTModel := isFluxFamily ||
		strings.Contains(strings.ToLower(opts.ModelName), "sd3") ||
		strings.Contains(strings.ToLower(opts.ModelName), "anima") ||
		strings.Contains(strings.ToLower(opts.ModelName), "kolors") ||
		strings.Contains(strings.ToLower(opts.ModelName), "auraflow") ||
		strings.Contains(strings.ToLower(opts.ModelName), "hidream")

	// Debug logging to verify model detection
	fmt.Printf("[DEBUG] txt2img: ModelName=%q, isFluxFamily=%v, isDiTModel=%v\n", opts.ModelName, isFluxFamily, isDiTModel)

	if isFluxFamily {
		if guidance == 0 {
			guidance = 3.5 // Default guidance for Flux/Chroma
		}
		if prediction == "" {
			prediction = "flux_flow" // Flux/Chroma require flow prediction
		}
	}

	reqBody := txt2imgRequest{
		Prompt:              opts.Prompt,
		NegativePrompt:      opts.NegativePrompt,
		Width:               opts.Width,
		Height:              opts.Height,
		Steps:               opts.Steps,
		Seed:                opts.Seed,
		SamplerName:         opts.Sampler,
		Preview:             "tae",
		Scheduler:           opts.Scheduler,
		CfgScale:            opts.CfgScale,
		BatchSize:           opts.BatchCount,
		EnableHR:            opts.EnableHR,
		HRScale:             opts.HRScale,
		HRUpscaler:          opts.HRUpscaler,
		DenoisingStrength:   opts.DenoisingStrength,
		Tiling:              opts.Tiling,
		ClipSkip:            opts.ClipSkip,
		LoRAs:               convertLoRAs(opts.LoRAs),
		ControlNetPath:      opts.ControlNetPath,
		ControlImage:        ctrlImgB64,
		ControlStrength:     opts.ControlStrength,
		VideoFrames:         opts.VideoFrames,
		FPS:                 opts.FPS,
		EndImage:            endImgB64,
		RNG:                 opts.RNG,
		SamplerRNG:          opts.SamplerRNG,
		DisableMetadata:     opts.DisableMetadata,
		CfgRescale:          opts.CfgRescale,
		FaceRestore:         opts.FaceRestore,
		VideoMotionStrength: opts.VideoMotionStrength,
		VideoSeedMode:       opts.VideoSeedMode,
		VideoDenoise:        opts.VideoDenoise,
		VideoLoop:           opts.VideoLoop,
		OutputFormat:        opts.OutputFormat,
	}

	// Only set Flux-specific parameters for Flux/Chroma models
	if isFluxFamily {
		reqBody.Guidance = guidance
		reqBody.Prediction = prediction
		reqBody.FlowShift = opts.FlowShift
	}

	// Only set DiT-specific parameters for DiT-based models
	if isDiTModel {
		reqBody.SlgScale = opts.SlgScale
		reqBody.SkipLayerStart = opts.SkipLayerStart
		reqBody.SkipLayerEnd = opts.SkipLayerEnd
		reqBody.Eta = opts.Eta
		reqBody.VAETileOverlap = opts.VAETileOverlap
		reqBody.SChurn = opts.SChurn
		reqBody.STmin = opts.STmin
		reqBody.STmax = opts.STmax
		reqBody.SNoise = opts.SNoise
		reqBody.HighNoiseSteps = opts.HighNoiseSteps
		reqBody.HighNoiseCfgScale = opts.HighNoiseCfgScale
		reqBody.HighNoiseGuidance = opts.HighNoiseGuidance
		reqBody.HighNoiseSlgScale = opts.HighNoiseSlgScale
		reqBody.HighNoiseEta = opts.HighNoiseEta
		reqBody.HighNoiseSkipLayerStart = opts.HighNoiseSkipLayerStart
		reqBody.HighNoiseSkipLayerEnd = opts.HighNoiseSkipLayerEnd
		fmt.Printf("[DEBUG] txt2img: Setting DiT parameters - Eta=%v, VAETileOverlap=%v, SlgScale=%v\n", reqBody.Eta, reqBody.VAETileOverlap, reqBody.SlgScale)
	} else {
		fmt.Printf("[DEBUG] txt2img: Skipping DiT parameters for non-DiT model\n")
	}
	if reqBody.BatchSize <= 0 {
		reqBody.BatchSize = 1
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, notDispatched(fmt.Errorf("marshal request: %w", err))
	}

	fmt.Printf("[DEBUG] txt2img request JSON: %s\n", string(jsonData))

	// Use synchronous API (/sdapi/v1/txt2img) because async job API has critical limitations:
	// 1. Does not respect steps parameter (always uses default 20)
	// 2. Does not return progress/preview data for live preview
	// The context is what makes cancellation real on this path: http.Post builds a
	// request nothing can abort, so a cancelled job used to leave its goroutine
	// blocked in the call until the server answered.
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/sdapi/v1/txt2img", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, notDispatched(fmt.Errorf("build txt2img request: %w", err))
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sd-server txt2img request failed: %w", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sd-server txt2img error %d: %s", resp.StatusCode, string(body))
	}

	var result txt2imgResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode txt2img response: %w", err)
	}
	if len(result.Images) == 0 {
		return nil, fmt.Errorf("no images in txt2img response")
	}

	return &types.ImageResult{
		ImageData: result.Images[0],
		Images:    result.Images,
		Metadata: types.ImageMetadata{
			Prompt:         opts.Prompt,
			NegativePrompt: opts.NegativePrompt,
			Width:          opts.Width,
			Height:         opts.Height,
			Steps:          opts.Steps,
			Seed:           opts.Seed,
			Sampler:        opts.Sampler,
			CfgScale:       opts.CfgScale,
			Mode:           opts.Mode,
			CreatedAt:      time.Now().Format(time.RFC3339),
			Options:        opts,
		},
	}, nil
}

// This client deliberately has no cancel-endpoint call.
//
// A DELETE /sdcpp/v1/jobs/{id} used to be sent for async jobs, but whether this
// build honours it is unestablished
// (build/logs/cancelled-generation-contract-audit.md). Calling a guessed endpoint
// is how an app comes to tell a user that work stopped when it did not.
// Cancellation closes the client request; ownership of what the server may still
// be doing lives in generation.go.

func (c *client) pollJob(ctx context.Context, jobID string) (*sdcppJob, error) {
	deadline := time.NewTimer(c.pollTimeout)
	defer deadline.Stop()
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/sdcpp/v1/jobs/"+jobID, nil)
		if err != nil {
			return nil, err
		}
		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("job status returned %d: %s", resp.StatusCode, string(body))
		}

		var job sdcppJob
		if err := json.Unmarshal(body, &job); err != nil {
			return nil, err
		}

		fmt.Printf("[DEBUG] pollJob() status=%s, images=%d, preview_len=%d\n",
			job.Status, len(job.Result.Images), len(job.PreviewImage))

		switch job.Status {
		case "completed":
			if len(job.Result.Images) == 0 {
				fmt.Printf("[DEBUG] pollJob() no images, attempting recovery\n")
				recovered := recoverImagesFromJobBody(body, job.PreviewImage)
				fmt.Printf("[DEBUG] pollJob() recovered %d images\n", len(recovered))
				if len(recovered) == 0 {
					return nil, fmt.Errorf("no images in job result")
				}
				job.Result.Images = recovered
			}
			return &job, nil
		case "failed":
			return nil, fmt.Errorf("job failed: %s", job.Error)
		case "cancelled":
			// The server reported it stopped this one, which is a backend statement
			// this client can pass on.
			return nil, fmt.Errorf("job cancelled by the engine")
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-deadline.C:
			return nil, fmt.Errorf("job poll timeout")
		case <-tick.C:
		}
	}
}

type img2imgRequest struct {
	Prompt                  string        `json:"prompt"`
	NegativePrompt          string        `json:"negative_prompt"`
	InitImages              []string      `json:"init_images"`
	Mask                    string        `json:"mask,omitempty"`
	Width                   int           `json:"width"`
	Height                  int           `json:"height"`
	Steps                   int           `json:"steps"`
	Seed                    int64         `json:"seed"`
	SamplerName             string        `json:"sampler_name"`
	Preview                 string        `json:"preview,omitempty"`
	Scheduler               string        `json:"scheduler"`
	CfgScale                float64       `json:"cfg_scale"`
	Guidance                float64       `json:"guidance,omitempty"`
	SlgScale                float64       `json:"slg_scale,omitempty"`
	SkipLayerStart          float64       `json:"skip_layer_start,omitempty"`
	SkipLayerEnd            float64       `json:"skip_layer_end,omitempty"`
	Eta                     float64       `json:"eta,omitempty"`
	BatchSize               int           `json:"batch_size"`
	DenoisingStrength       float64       `json:"denoising_strength"`
	ResizeMode              string        `json:"resize_mode,omitempty"`
	StyleStrength           float64       `json:"style_strength,omitempty"`
	PreserveColors          bool          `json:"preserve_colors,omitempty"`
	MatchHistogram          bool          `json:"match_histogram,omitempty"`
	LatentNoise             float64       `json:"latent_noise,omitempty"`
	MaskBlur                int           `json:"mask_blur,omitempty"`
	MaskInvert              bool          `json:"mask_invert,omitempty"`
	MaskFeather             int           `json:"mask_feather,omitempty"`
	MaskedContent           string        `json:"masked_content,omitempty"`
	EdgeBehavior            string        `json:"edge_behavior,omitempty"`
	UpscaleMethod           string        `json:"upscale_method,omitempty"`
	UpscaleSharpen          string        `json:"upscale_sharpen,omitempty"`
	UpscaleAntiArtifact     string        `json:"upscale_anti_artifact,omitempty"`
	Tiling                  bool          `json:"tiling"`
	ClipSkip                int           `json:"clip_skip"`
	LoRAs                   []LoRARequest `json:"loras,omitempty"`
	ControlNetPath          string        `json:"control_net_path,omitempty"`
	ControlImage            string        `json:"control_image,omitempty"`
	ControlStrength         float64       `json:"control_strength,omitempty"`
	CacheMode               string        `json:"cache_mode,omitempty"`
	CacheOption             string        `json:"cache_option,omitempty"`
	VideoFrames             int           `json:"video_frames,omitempty"`
	FPS                     int           `json:"fps,omitempty"`
	EndImage                string        `json:"end_image,omitempty"`
	VideoMotionStrength     float64       `json:"video_motion_strength,omitempty"`
	VideoSeedMode           string        `json:"video_seed_mode,omitempty"`
	VideoDenoise            float64       `json:"video_denoise,omitempty"`
	VideoLoop               bool          `json:"video_loop,omitempty"`
	Prediction              string        `json:"prediction,omitempty"`
	RNG                     string        `json:"rng,omitempty"`
	SamplerRNG              string        `json:"sampler_rng,omitempty"`
	FlowShift               float64       `json:"flow_shift,omitempty"`
	VAETileOverlap          float64       `json:"vae_tile_overlap,omitempty"`
	DisableMetadata         bool          `json:"disable_metadata,omitempty"`
	CfgRescale              float64       `json:"cfg_rescale,omitempty"`
	FaceRestore             float64       `json:"face_restore,omitempty"`
	SChurn                  float64       `json:"s_churn,omitempty"`
	STmin                   float64       `json:"s_tmin,omitempty"`
	STmax                   float64       `json:"s_tmax,omitempty"`
	SNoise                  float64       `json:"s_noise,omitempty"`
	OutputFormat            string        `json:"output_format,omitempty"`
	HighNoiseSteps          int           `json:"high_noise_steps,omitempty"`
	HighNoiseCfgScale       float64       `json:"high_noise_cfg_scale,omitempty"`
	HighNoiseGuidance       float64       `json:"high_noise_guidance,omitempty"`
	HighNoiseSlgScale       float64       `json:"high_noise_slg_scale,omitempty"`
	HighNoiseEta            float64       `json:"high_noise_eta,omitempty"`
	HighNoiseSkipLayerStart float64       `json:"high_noise_skip_layer_start,omitempty"`
	HighNoiseSkipLayerEnd   float64       `json:"high_noise_skip_layer_end,omitempty"`
}

func (c *client) img2img(ctx context.Context, opts types.ImageOptions) (*types.ImageResult, error) {
	if opts.BatchCount > 1 {
		return c.generateSequentialBatch(ctx, c.img2img, opts)
	}
	if opts.InitImagePath == "" {
		return nil, notDispatched(fmt.Errorf("img2img requires an init image path"))
	}

	imgData, err := os.ReadFile(opts.InitImagePath)
	if err != nil {
		return nil, notDispatched(fmt.Errorf("read init image: %w", err))
	}
	b64Img := base64.StdEncoding.EncodeToString(imgData)

	var ctrlImgB64 string
	if opts.ControlImagePath != "" {
		if data, err := os.ReadFile(opts.ControlImagePath); err == nil {
			ctrlImgB64 = base64.StdEncoding.EncodeToString(data)
		}
	}

	var maskB64 string
	if opts.MaskPath != "" {
		if data, err := os.ReadFile(opts.MaskPath); err == nil {
			maskB64 = base64.StdEncoding.EncodeToString(data)
		}
	}

	// Flux/Chroma use 'guidance' instead of 'cfg_scale'
	guidance := opts.Guidance
	prediction := opts.Prediction
	isFluxFamily := strings.Contains(strings.ToLower(opts.ModelName), "flux") ||
		strings.Contains(strings.ToLower(opts.ModelName), "chroma")

	// DiT-based models: Flux, Chroma, SD3, SD3.5, SDXL (some), Anima, Kolors, AuraFlow, HiDream
	isDiTModel := isFluxFamily ||
		strings.Contains(strings.ToLower(opts.ModelName), "sd3") ||
		strings.Contains(strings.ToLower(opts.ModelName), "anima") ||
		strings.Contains(strings.ToLower(opts.ModelName), "kolors") ||
		strings.Contains(strings.ToLower(opts.ModelName), "auraflow") ||
		strings.Contains(strings.ToLower(opts.ModelName), "hidream")

	if isFluxFamily {
		if guidance == 0 {
			guidance = 3.5
		}
		if prediction == "" {
			prediction = "flux_flow"
		}
	}

	reqBody := img2imgRequest{
		Prompt:              opts.Prompt,
		NegativePrompt:      opts.NegativePrompt,
		InitImages:          []string{b64Img},
		Mask:                maskB64,
		Width:               opts.Width,
		Height:              opts.Height,
		Steps:               opts.Steps,
		Seed:                opts.Seed,
		SamplerName:         opts.Sampler,
		Preview:             "tae",
		Scheduler:           opts.Scheduler,
		CfgScale:            opts.CfgScale,
		BatchSize:           opts.BatchCount,
		DenoisingStrength:   opts.Strength,
		ResizeMode:          opts.ResizeMode,
		StyleStrength:       opts.StyleStrength,
		PreserveColors:      opts.PreserveColors,
		MatchHistogram:      opts.MatchHistogram,
		LatentNoise:         opts.LatentNoise,
		MaskBlur:            opts.MaskBlur,
		MaskInvert:          opts.MaskInvert,
		UpscaleMethod:       opts.UpscaleMethod,
		UpscaleSharpen:      opts.UpscaleSharpen,
		UpscaleAntiArtifact: opts.UpscaleAntiArtifact,
		Tiling:              opts.Tiling,
		ClipSkip:            opts.ClipSkip,
		LoRAs:               convertLoRAs(opts.LoRAs),
		ControlNetPath:      opts.ControlNetPath,
		ControlImage:        ctrlImgB64,
		ControlStrength:     opts.ControlStrength,
		RNG:                 opts.RNG,
		SamplerRNG:          opts.SamplerRNG,
		DisableMetadata:     opts.DisableMetadata,
		CfgRescale:          opts.CfgRescale,
		FaceRestore:         opts.FaceRestore,
		VideoMotionStrength: opts.VideoMotionStrength,
		VideoSeedMode:       opts.VideoSeedMode,
		VideoDenoise:        opts.VideoDenoise,
		VideoLoop:           opts.VideoLoop,
		OutputFormat:        opts.OutputFormat,
	}

	// Only set Flux-specific parameters for Flux/Chroma models
	if isFluxFamily {
		reqBody.Guidance = guidance
		reqBody.Prediction = prediction
		reqBody.FlowShift = opts.FlowShift
	}

	// Only set DiT-specific parameters for DiT-based models
	if isDiTModel {
		reqBody.SlgScale = opts.SlgScale
		reqBody.SkipLayerStart = opts.SkipLayerStart
		reqBody.SkipLayerEnd = opts.SkipLayerEnd
		reqBody.Eta = opts.Eta
		reqBody.VAETileOverlap = opts.VAETileOverlap
		reqBody.SChurn = opts.SChurn
		reqBody.STmin = opts.STmin
		reqBody.STmax = opts.STmax
		reqBody.SNoise = opts.SNoise
		reqBody.HighNoiseSteps = opts.HighNoiseSteps
		reqBody.HighNoiseCfgScale = opts.HighNoiseCfgScale
		reqBody.HighNoiseGuidance = opts.HighNoiseGuidance
		reqBody.HighNoiseSlgScale = opts.HighNoiseSlgScale
		reqBody.HighNoiseEta = opts.HighNoiseEta
		reqBody.HighNoiseSkipLayerStart = opts.HighNoiseSkipLayerStart
		reqBody.HighNoiseSkipLayerEnd = opts.HighNoiseSkipLayerEnd
	}
	if reqBody.BatchSize <= 0 {
		reqBody.BatchSize = 1
	}
	if reqBody.DenoisingStrength <= 0 {
		reqBody.DenoisingStrength = 0.75
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, notDispatched(fmt.Errorf("marshal img2img request: %w", err))
	}

	// Use async job API for live preview support
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/sdcpp/v1/img_gen", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, notDispatched(fmt.Errorf("build img_gen request: %w", err))
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sd-server img_gen request failed: %w", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("sd-server img_gen error %d: %s", resp.StatusCode, string(body))
	}

	var jobResp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &jobResp); err != nil {
		return nil, fmt.Errorf("decode job response: %w", err)
	}
	if jobResp.ID == "" {
		return nil, fmt.Errorf("no job ID in response")
	}

	// Own the progress id with a token: a late cleanup may only clear its own.
	token := c.setActiveJob(jobResp.ID)

	// Poll job until completion.
	job, err := c.pollJob(ctx, jobResp.ID)
	c.resetActiveJob(token)
	if err != nil {
		return nil, err
	}

	if len(job.Result.Images) == 0 {
		return nil, fmt.Errorf("no images in job result")
	}

	return &types.ImageResult{
		ImageData: job.Result.Images[0],
		Images:    job.Result.Images,
		Metadata: types.ImageMetadata{
			Prompt:         opts.Prompt,
			NegativePrompt: opts.NegativePrompt,
			Width:          opts.Width,
			Height:         opts.Height,
			Steps:          opts.Steps,
			Seed:           opts.Seed,
			Sampler:        opts.Sampler,
			CfgScale:       opts.CfgScale,
			Mode:           opts.Mode,
			CreatedAt:      time.Now().Format(time.RFC3339),
			Options:        opts,
		},
	}, nil
}

func (c *client) inpaint(ctx context.Context, opts types.ImageOptions) (*types.ImageResult, error) {
	if opts.BatchCount > 1 {
		return c.generateSequentialBatch(ctx, c.inpaint, opts)
	}
	if opts.InitImagePath == "" {
		return nil, notDispatched(fmt.Errorf("inpaint requires an init image path"))
	}
	if opts.MaskPath == "" {
		return nil, notDispatched(fmt.Errorf("inpaint requires a mask path"))
	}

	imgData, err := os.ReadFile(opts.InitImagePath)
	if err != nil {
		return nil, notDispatched(fmt.Errorf("read init image: %w", err))
	}
	maskData, err := os.ReadFile(opts.MaskPath)
	if err != nil {
		return nil, notDispatched(fmt.Errorf("read mask: %w", err))
	}
	b64Img := base64.StdEncoding.EncodeToString(imgData)
	b64Mask := base64.StdEncoding.EncodeToString(maskData)

	var ctrlImgB64 string
	if opts.ControlImagePath != "" {
		if data, err := os.ReadFile(opts.ControlImagePath); err == nil {
			ctrlImgB64 = base64.StdEncoding.EncodeToString(data)
		}
	}

	// Flux/Chroma use 'guidance' instead of 'cfg_scale'
	guidance := opts.Guidance
	prediction := opts.Prediction
	isFluxFamily := strings.Contains(strings.ToLower(opts.ModelName), "flux") ||
		strings.Contains(strings.ToLower(opts.ModelName), "chroma")

	// DiT-based models: Flux, Chroma, SD3, SD3.5, SDXL (some), Anima, Kolors, AuraFlow, HiDream
	isDiTModel := isFluxFamily ||
		strings.Contains(strings.ToLower(opts.ModelName), "sd3") ||
		strings.Contains(strings.ToLower(opts.ModelName), "anima") ||
		strings.Contains(strings.ToLower(opts.ModelName), "kolors") ||
		strings.Contains(strings.ToLower(opts.ModelName), "auraflow") ||
		strings.Contains(strings.ToLower(opts.ModelName), "hidream")

	if isFluxFamily {
		if guidance == 0 {
			guidance = 3.5
		}
		if prediction == "" {
			prediction = "flux_flow"
		}
	}

	reqBody := img2imgRequest{
		Prompt:              opts.Prompt,
		NegativePrompt:      opts.NegativePrompt,
		InitImages:          []string{b64Img},
		Mask:                b64Mask,
		Width:               opts.Width,
		Height:              opts.Height,
		Steps:               opts.Steps,
		Seed:                opts.Seed,
		SamplerName:         opts.Sampler,
		Preview:             "tae",
		Scheduler:           opts.Scheduler,
		CfgScale:            opts.CfgScale,
		BatchSize:           opts.BatchCount,
		DenoisingStrength:   opts.Strength,
		ResizeMode:          opts.ResizeMode,
		StyleStrength:       opts.StyleStrength,
		PreserveColors:      opts.PreserveColors,
		MatchHistogram:      opts.MatchHistogram,
		LatentNoise:         opts.LatentNoise,
		MaskBlur:            opts.MaskBlur,
		MaskInvert:          opts.MaskInvert,
		MaskFeather:         opts.MaskFeather,
		MaskedContent:       opts.MaskedContent,
		EdgeBehavior:        opts.EdgeBehavior,
		Tiling:              opts.Tiling,
		ClipSkip:            opts.ClipSkip,
		LoRAs:               convertLoRAs(opts.LoRAs),
		ControlNetPath:      opts.ControlNetPath,
		ControlImage:        ctrlImgB64,
		ControlStrength:     opts.ControlStrength,
		VideoFrames:         opts.VideoFrames,
		FPS:                 opts.FPS,
		RNG:                 opts.RNG,
		SamplerRNG:          opts.SamplerRNG,
		DisableMetadata:     opts.DisableMetadata,
		CfgRescale:          opts.CfgRescale,
		FaceRestore:         opts.FaceRestore,
		VideoMotionStrength: opts.VideoMotionStrength,
		VideoSeedMode:       opts.VideoSeedMode,
		VideoDenoise:        opts.VideoDenoise,
		VideoLoop:           opts.VideoLoop,
		OutputFormat:        opts.OutputFormat,
	}
	if reqBody.BatchSize <= 0 {
		reqBody.BatchSize = 1
	}
	if reqBody.DenoisingStrength <= 0 {
		reqBody.DenoisingStrength = 0.75
	}

	// Only set Flux-specific parameters for Flux/Chroma models
	if isFluxFamily {
		reqBody.Guidance = guidance
		reqBody.Prediction = prediction
		reqBody.FlowShift = opts.FlowShift
	}

	// Only set DiT-specific parameters for DiT-based models
	if isDiTModel {
		reqBody.SlgScale = opts.SlgScale
		reqBody.SkipLayerStart = opts.SkipLayerStart
		reqBody.SkipLayerEnd = opts.SkipLayerEnd
		reqBody.Eta = opts.Eta
		reqBody.VAETileOverlap = opts.VAETileOverlap
		reqBody.SChurn = opts.SChurn
		reqBody.STmin = opts.STmin
		reqBody.STmax = opts.STmax
		reqBody.SNoise = opts.SNoise
		reqBody.HighNoiseSteps = opts.HighNoiseSteps
		reqBody.HighNoiseCfgScale = opts.HighNoiseCfgScale
		reqBody.HighNoiseGuidance = opts.HighNoiseGuidance
		reqBody.HighNoiseSlgScale = opts.HighNoiseSlgScale
		reqBody.HighNoiseEta = opts.HighNoiseEta
		reqBody.HighNoiseSkipLayerStart = opts.HighNoiseSkipLayerStart
		reqBody.HighNoiseSkipLayerEnd = opts.HighNoiseSkipLayerEnd
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, notDispatched(fmt.Errorf("marshal inpaint request: %w", err))
	}

	// Use async job API for live preview support
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/sdcpp/v1/img_gen", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, notDispatched(fmt.Errorf("build img_gen request: %w", err))
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sd-server img_gen request failed: %w", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("sd-server img_gen error %d: %s", resp.StatusCode, string(body))
	}

	var jobResp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &jobResp); err != nil {
		return nil, fmt.Errorf("decode job response: %w", err)
	}
	if jobResp.ID == "" {
		return nil, fmt.Errorf("no job ID in response")
	}

	// Own the progress id with a token: a late cleanup may only clear its own.
	token := c.setActiveJob(jobResp.ID)

	// Poll job until completion.
	job, err := c.pollJob(ctx, jobResp.ID)
	c.resetActiveJob(token)
	if err != nil {
		return nil, err
	}

	if len(job.Result.Images) == 0 {
		return nil, fmt.Errorf("no images in job result")
	}

	return &types.ImageResult{
		ImageData: job.Result.Images[0],
		Images:    job.Result.Images,
		Metadata: types.ImageMetadata{
			Prompt:         opts.Prompt,
			NegativePrompt: opts.NegativePrompt,
			Width:          opts.Width,
			Height:         opts.Height,
			Steps:          opts.Steps,
			Seed:           opts.Seed,
			Sampler:        opts.Sampler,
			CfgScale:       opts.CfgScale,
			Mode:           opts.Mode,
			CreatedAt:      time.Now().Format(time.RFC3339),
			Options:        opts,
		},
	}, nil
}

func (c *client) getSamplers() ([]string, error) {
	resp, err := c.http.Get(c.baseURL + "/sdapi/v1/samplers")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var list []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}
	names := make([]string, len(list))
	for i, s := range list {
		names[i] = s.Name
	}
	return names, nil
}

// sdcppJob represents the native async job response from /sdcpp/v1/jobs/{id}.
type sdcppJob struct {
	ID           string  `json:"id"`
	Status       string  `json:"status"`
	Progress     float64 `json:"progress"`
	Step         int     `json:"step"`
	TotalSteps   int     `json:"total_steps"`
	PreviewImage string  `json:"preview_image"`
	Error        string  `json:"error"`
	Result       struct {
		Images sdcppImageList `json:"images"`
	} `json:"result"`
}

func (c *client) progress() (types.ProgressInfo, error) {
	jobID := c.activeJob()

	// /sdapi/v1/progress is job-agnostic and is the only place a decoded
	// preview frame shows up (current_image), so it has to be queried even
	// when no async job id exists. The id is only needed by the fallback.
	fmt.Printf("[DEBUG] progress() called (jobID=%q)", jobID)

	// Try legacy /sdapi/v1/progress endpoint first (may have preview data).
	// Bounded, so a stalled engine hands control back to the job poller.
	resp, err := c.httpShort.Get(c.baseURL + "/sdapi/v1/progress")
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			var legacyProgress struct {
				Progress     float64 `json:"progress"`
				EtaRelative  float64 `json:"eta_relative"`
				EtaAbsolute  float64 `json:"eta_absolute"`
				CurrentImage string  `json:"current_image"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&legacyProgress); err == nil {
				if legacyProgress.Progress > 0 || legacyProgress.CurrentImage != "" {
					fmt.Printf("[DEBUG] progress() legacy endpoint: progress=%.2f, preview_len=%d\n",
						legacyProgress.Progress, len(legacyProgress.CurrentImage))
					return types.ProgressInfo{
						Progress:     legacyProgress.Progress,
						PreviewImage: legacyProgress.CurrentImage,
					}, nil
				}
			}
		}
	}

	if jobID == "" {
		return types.ProgressInfo{}, nil
	}

	// Fallback to async job endpoint
	resp, err = c.httpShort.Get(c.baseURL + "/sdcpp/v1/jobs/" + jobID)
	if err != nil {
		fmt.Printf("[DEBUG] progress() HTTP error: %v\n", err)
		return types.ProgressInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("[DEBUG] progress() job not found (404)\n")
		return types.ProgressInfo{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[DEBUG] progress() unexpected status: %d\n", resp.StatusCode)
		return types.ProgressInfo{}, fmt.Errorf("job status returned %d", resp.StatusCode)
	}

	var job sdcppJob
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		fmt.Printf("[DEBUG] progress() decode error: %v\n", err)
		return types.ProgressInfo{}, err
	}

	fmt.Printf("[DEBUG] progress() job status: %s, step: %d/%d, progress: %.2f, preview_len: %d\n",
		job.Status, job.Step, job.TotalSteps, job.Progress, len(job.PreviewImage))

	return types.ProgressInfo{
		Progress:     job.Progress,
		Step:         job.Step,
		TotalSteps:   job.TotalSteps,
		PreviewImage: job.PreviewImage,
	}, nil
}
