package sdcpp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"howl-chat/internal/image/types"
)

// client talks to sd-server over HTTP.
type client struct {
	baseURL string
	http    *http.Client
}

func newClient(baseURL string) *client {
	return &client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 600 * time.Second},
	}
}

func (c *client) waitForReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := c.http.Get(c.baseURL + "/sdapi/v1/samplers")
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("sd-server not ready after %v", timeout)
}

// txt2imgRequest mirrors the WebUI /sdapi/v1/txt2img payload.
type txt2imgRequest struct {
	Prompt            string   `json:"prompt"`
	NegativePrompt    string   `json:"negative_prompt"`
	Width             int      `json:"width"`
	Height            int      `json:"height"`
	Steps             int      `json:"steps"`
	Seed              int64    `json:"seed"`
	SamplerName       string   `json:"sampler_name"`
	Scheduler         string   `json:"scheduler"`
	CfgScale          float64  `json:"cfg_scale"`
	Guidance          float64  `json:"guidance"`
	SlgScale          float64  `json:"slg_scale"`
	SkipLayerStart    float64  `json:"skip_layer_start"`
	SkipLayerEnd      float64  `json:"skip_layer_end"`
	Eta               float64  `json:"eta"`
	BatchSize         int      `json:"batch_size"`
	EnableHR          bool     `json:"enable_hr"`
	HRScale           float64  `json:"hr_scale"`
	HRUpscaler        string   `json:"hr_upscaler"`
	DenoisingStrength float64  `json:"denoising_strength"`
	Tiling            bool     `json:"tiling"`
	ClipSkip          int      `json:"clip_skip"`
	LoRAs             []string `json:"loras,omitempty"`
	ControlImage      string   `json:"control_image,omitempty"`
	ControlStrength   float64  `json:"control_strength,omitempty"`
	CacheMode         string   `json:"cache_mode,omitempty"`
	CacheOption       string   `json:"cache_option,omitempty"`
	VideoFrames       int      `json:"video_frames,omitempty"`
	FPS               int      `json:"fps,omitempty"`
}

type txt2imgResponse struct {
	Images     []string `json:"images"`
	Parameters struct{} `json:"parameters"`
	Info       string   `json:"info"`
}

func (c *client) txt2img(opts types.ImageOptions) (*types.ImageResult, error) {
	var ctrlImgB64 string
	if opts.ControlImagePath != "" {
		if data, err := os.ReadFile(opts.ControlImagePath); err == nil {
			ctrlImgB64 = base64.StdEncoding.EncodeToString(data)
		}
	}

	reqBody := txt2imgRequest{
		Prompt:            opts.Prompt,
		NegativePrompt:    opts.NegativePrompt,
		Width:             opts.Width,
		Height:            opts.Height,
		Steps:             opts.Steps,
		Seed:              opts.Seed,
		SamplerName:       opts.Sampler,
		Scheduler:         opts.Scheduler,
		CfgScale:          opts.CfgScale,
		Guidance:          opts.Guidance,
		SlgScale:          opts.SlgScale,
		SkipLayerStart:    opts.SkipLayerStart,
		SkipLayerEnd:      opts.SkipLayerEnd,
		Eta:               opts.Eta,
		BatchSize:         opts.BatchCount,
		EnableHR:          opts.EnableHR,
		HRScale:           opts.HRScale,
		HRUpscaler:        opts.HRUpscaler,
		DenoisingStrength: opts.DenoisingStrength,
		Tiling:            opts.Tiling,
		ClipSkip:          opts.ClipSkip,
		LoRAs:             opts.LoRAs,
		ControlImage:      ctrlImgB64,
		ControlStrength:   opts.ControlStrength,
		CacheMode:         opts.CacheMode,
		CacheOption:       opts.CacheOption,
		VideoFrames:       opts.VideoFrames,
		FPS:               opts.FPS,
	}
	if reqBody.BatchSize <= 0 {
		reqBody.BatchSize = 1
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.http.Post(c.baseURL+"/sdapi/v1/txt2img", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("sd-server txt2img request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sd-server txt2img error %d: %s", resp.StatusCode, string(body))
	}

	var result txt2imgResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode txt2img response: %w", err)
	}
	if len(result.Images) == 0 {
		return nil, fmt.Errorf("no images in response")
	}

	return &types.ImageResult{
		ImageData: result.Images[0],
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
		},
	}, nil
}

type img2imgRequest struct {
	Prompt            string   `json:"prompt"`
	NegativePrompt    string   `json:"negative_prompt"`
	InitImages        []string `json:"init_images"`
	Mask              string   `json:"mask,omitempty"`
	Width             int      `json:"width"`
	Height            int      `json:"height"`
	Steps             int      `json:"steps"`
	Seed              int64    `json:"seed"`
	SamplerName       string   `json:"sampler_name"`
	Scheduler         string   `json:"scheduler"`
	CfgScale          float64  `json:"cfg_scale"`
	Guidance          float64  `json:"guidance"`
	SlgScale          float64  `json:"slg_scale"`
	SkipLayerStart    float64  `json:"skip_layer_start"`
	SkipLayerEnd      float64  `json:"skip_layer_end"`
	Eta               float64  `json:"eta"`
	BatchSize         int      `json:"batch_size"`
	DenoisingStrength float64  `json:"denoising_strength"`
	Tiling            bool     `json:"tiling"`
	ClipSkip          int      `json:"clip_skip"`
	LoRAs             []string `json:"loras,omitempty"`
	ControlImage      string   `json:"control_image,omitempty"`
	ControlStrength   float64  `json:"control_strength,omitempty"`
	CacheMode         string   `json:"cache_mode,omitempty"`
	CacheOption       string   `json:"cache_option,omitempty"`
	VideoFrames       int      `json:"video_frames,omitempty"`
	FPS               int      `json:"fps,omitempty"`
}

func (c *client) img2img(opts types.ImageOptions) (*types.ImageResult, error) {
	if opts.InitImagePath == "" {
		return nil, fmt.Errorf("img2img requires an init image path")
	}

	imgData, err := os.ReadFile(opts.InitImagePath)
	if err != nil {
		return nil, fmt.Errorf("read init image: %w", err)
	}
	b64Img := base64.StdEncoding.EncodeToString(imgData)

	var ctrlImgB64 string
	if opts.ControlImagePath != "" {
		if data, err := os.ReadFile(opts.ControlImagePath); err == nil {
			ctrlImgB64 = base64.StdEncoding.EncodeToString(data)
		}
	}

	reqBody := img2imgRequest{
		Prompt:            opts.Prompt,
		NegativePrompt:    opts.NegativePrompt,
		InitImages:        []string{b64Img},
		Width:             opts.Width,
		Height:            opts.Height,
		Steps:             opts.Steps,
		Seed:              opts.Seed,
		SamplerName:       opts.Sampler,
		Scheduler:         opts.Scheduler,
		CfgScale:          opts.CfgScale,
		Guidance:          opts.Guidance,
		SlgScale:          opts.SlgScale,
		SkipLayerStart:    opts.SkipLayerStart,
		SkipLayerEnd:      opts.SkipLayerEnd,
		Eta:               opts.Eta,
		BatchSize:         opts.BatchCount,
		DenoisingStrength: opts.Strength,
		Tiling:            opts.Tiling,
		ClipSkip:          opts.ClipSkip,
		LoRAs:             opts.LoRAs,
		ControlImage:      ctrlImgB64,
		ControlStrength:   opts.ControlStrength,
		CacheMode:         opts.CacheMode,
		CacheOption:       opts.CacheOption,
	}
	if reqBody.BatchSize <= 0 {
		reqBody.BatchSize = 1
	}
	if reqBody.DenoisingStrength <= 0 {
		reqBody.DenoisingStrength = 0.75
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal img2img request: %w", err)
	}

	resp, err := c.http.Post(c.baseURL+"/sdapi/v1/img2img", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("sd-server img2img request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sd-server img2img error %d: %s", resp.StatusCode, string(body))
	}

	var result txt2imgResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode img2img response: %w", err)
	}
	if len(result.Images) == 0 {
		return nil, fmt.Errorf("no images in img2img response")
	}

	return &types.ImageResult{
		ImageData: result.Images[0],
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
		},
	}, nil
}

func (c *client) inpaint(opts types.ImageOptions) (*types.ImageResult, error) {
	if opts.InitImagePath == "" {
		return nil, fmt.Errorf("inpaint requires an init image path")
	}
	if opts.MaskPath == "" {
		return nil, fmt.Errorf("inpaint requires a mask path")
	}

	imgData, err := os.ReadFile(opts.InitImagePath)
	if err != nil {
		return nil, fmt.Errorf("read init image: %w", err)
	}
	maskData, err := os.ReadFile(opts.MaskPath)
	if err != nil {
		return nil, fmt.Errorf("read mask: %w", err)
	}
	b64Img := base64.StdEncoding.EncodeToString(imgData)
	b64Mask := base64.StdEncoding.EncodeToString(maskData)

	var ctrlImgB64 string
	if opts.ControlImagePath != "" {
		if data, err := os.ReadFile(opts.ControlImagePath); err == nil {
			ctrlImgB64 = base64.StdEncoding.EncodeToString(data)
		}
	}

	reqBody := img2imgRequest{
		Prompt:            opts.Prompt,
		NegativePrompt:    opts.NegativePrompt,
		InitImages:        []string{b64Img},
		Mask:              b64Mask,
		Width:             opts.Width,
		Height:            opts.Height,
		Steps:             opts.Steps,
		Seed:              opts.Seed,
		SamplerName:       opts.Sampler,
		Scheduler:         opts.Scheduler,
		CfgScale:          opts.CfgScale,
		Guidance:          opts.Guidance,
		SlgScale:          opts.SlgScale,
		SkipLayerStart:    opts.SkipLayerStart,
		SkipLayerEnd:      opts.SkipLayerEnd,
		Eta:               opts.Eta,
		BatchSize:         opts.BatchCount,
		DenoisingStrength: opts.Strength,
		Tiling:            opts.Tiling,
		ClipSkip:          opts.ClipSkip,
		LoRAs:             opts.LoRAs,
		ControlImage:      ctrlImgB64,
		ControlStrength:   opts.ControlStrength,
		CacheMode:         opts.CacheMode,
		CacheOption:       opts.CacheOption,
		VideoFrames:       opts.VideoFrames,
		FPS:               opts.FPS,
	}
	if reqBody.BatchSize <= 0 {
		reqBody.BatchSize = 1
	}
	if reqBody.DenoisingStrength <= 0 {
		reqBody.DenoisingStrength = 0.75
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal inpaint request: %w", err)
	}

	resp, err := c.http.Post(c.baseURL+"/sdapi/v1/img2img", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("sd-server inpaint request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sd-server inpaint error %d: %s", resp.StatusCode, string(body))
	}

	var result txt2imgResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode inpaint response: %w", err)
	}
	if len(result.Images) == 0 {
		return nil, fmt.Errorf("no images in inpaint response")
	}

	return &types.ImageResult{
		ImageData: result.Images[0],
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
		},
	}, nil
}

func (c *client) video(opts types.ImageOptions) (*types.ImageResult, error) {
	// Video generation often requires an init image (SVD) or text (WAN).
	// We'll use the img2img endpoint if an init image is present, otherwise txt2img.
	if opts.InitImagePath != "" {
		imgData, err := os.ReadFile(opts.InitImagePath)
		if err != nil {
			return nil, fmt.Errorf("read init image: %w", err)
		}
		b64Img := base64.StdEncoding.EncodeToString(imgData)

		var ctrlImgB64 string
		if opts.ControlImagePath != "" {
			if data, err := os.ReadFile(opts.ControlImagePath); err == nil {
				ctrlImgB64 = base64.StdEncoding.EncodeToString(data)
			}
		}

		reqBody := img2imgRequest{
			Prompt:            opts.Prompt,
			NegativePrompt:    opts.NegativePrompt,
			InitImages:        []string{b64Img},
			Width:             opts.Width,
			Height:            opts.Height,
			Steps:             opts.Steps,
			Seed:              opts.Seed,
			SamplerName:       opts.Sampler,
			Scheduler:         opts.Scheduler,
			CfgScale:          opts.CfgScale,
			Guidance:          opts.Guidance,
			SlgScale:          opts.SlgScale,
			SkipLayerStart:    opts.SkipLayerStart,
			SkipLayerEnd:      opts.SkipLayerEnd,
			Eta:               opts.Eta,
			BatchSize:         opts.BatchCount,
			DenoisingStrength: opts.Strength,
			Tiling:            opts.Tiling,
			ClipSkip:          opts.ClipSkip,
			LoRAs:             opts.LoRAs,
			ControlImage:      ctrlImgB64,
			ControlStrength:   opts.ControlStrength,
			CacheMode:         opts.CacheMode,
			CacheOption:       opts.CacheOption,
			VideoFrames:       opts.VideoFrames,
			FPS:               opts.FPS,
		}

		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("marshal video request: %w", err)
		}

		resp, err := c.http.Post(c.baseURL+"/sdapi/v1/img2img", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("sd-server video request failed: %w", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("sd-server video error %d: %s", resp.StatusCode, string(body))
		}

		var result txt2imgResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("decode video response: %w", err)
		}
		if len(result.Images) == 0 {
			return nil, fmt.Errorf("no images in video response")
		}

		return &types.ImageResult{
			ImageData: result.Images[0],
			Metadata: types.ImageMetadata{
				Prompt: opts.Prompt,
				Width:  opts.Width,
				Height: opts.Height,
				Steps:  opts.Steps,
				Mode:   "video",
			},
		}, nil
	}

	// Txt2Video flow
	return c.txt2img(opts)
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

// progressResponse mirrors AUTOMATIC1111 /sdapi/v1/progress format.
type progressResponse struct {
	Progress     float64 `json:"progress"`
	EtaRelative  float64 `json:"eta_relative"`
	CurrentImage string  `json:"current_image"` // base64 PNG preview
	State        struct {
		Skipped       bool `json:"skipped"`
		Interrupted   bool `json:"interrupted"`
		SamplingStep  int  `json:"sampling_step"`
		SamplingSteps int  `json:"sampling_steps"`
	} `json:"state"`
}

func (c *client) progress() (types.ProgressInfo, error) {
	resp, err := c.http.Get(c.baseURL + "/sdapi/v1/progress")
	if err != nil {
		return types.ProgressInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.ProgressInfo{}, fmt.Errorf("progress endpoint returned %d", resp.StatusCode)
	}

	var pr progressResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return types.ProgressInfo{}, err
	}
	return types.ProgressInfo{
		Progress:     pr.Progress,
		Step:         pr.State.SamplingStep,
		TotalSteps:   pr.State.SamplingSteps,
		PreviewImage: pr.CurrentImage,
	}, nil
}
