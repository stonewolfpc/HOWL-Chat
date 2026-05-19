package types

// ImageOptions describes a single generation request.
type ImageOptions struct {
	Prompt            string   `json:"prompt"`
	NegativePrompt    string   `json:"negative_prompt"`
	Width             int      `json:"width"`
	Height            int      `json:"height"`
	Steps             int      `json:"steps"`
	Seed              int64    `json:"seed"`
	Sampler           string   `json:"sampler"`
	Scheduler         string   `json:"scheduler"`
	CfgScale          float64  `json:"cfg_scale"`
	Guidance          float64  `json:"guidance"`  // distilled guidance (SD3/Flux)
	SlgScale          float64  `json:"slg_scale"` // skip layer guidance
	SkipLayerStart    float64  `json:"skip_layer_start"`
	SkipLayerEnd      float64  `json:"skip_layer_end"`
	Eta               float64  `json:"eta"`             // noise multiplier
	Strength          float64  `json:"strength"`        // img2img / inpaint
	InitImagePath     string   `json:"init_image_path"` // img2img / inpaint
	MaskPath          string   `json:"mask_path"`       // inpaint only
	ModelName         string   `json:"model_name"`
	VAEPath           string   `json:"vae_path"`           // standalone VAE
	ClipLPath         string   `json:"clip_l_path"`        // CLIP-L
	ClipGPath         string   `json:"clip_g_path"`        // CLIP-G
	T5XXLPath         string   `json:"t5xxl_path"`         // T5XXL
	ControlNetPath    string   `json:"control_net_path"`   // ControlNet model
	ControlImagePath  string   `json:"control_image_path"` // ControlNet image
	ControlStrength   float64  `json:"control_strength"`   // ControlNet strength
	BatchCount        int      `json:"batch_count"`
	Mode              string   `json:"mode"` // "txt2img", "img2img", "inpaint", "video"
	LoRAs             []string `json:"loras,omitempty"`
	EnableHR          bool     `json:"enable_hr"`
	HRScale           float64  `json:"hr_scale"`
	HRUpscaler        string   `json:"hr_upscaler"`
	DenoisingStrength float64  `json:"denoising_strength"`
	Tiling            bool     `json:"tiling"`
	ClipSkip          int      `json:"clip_skip"`
	VideoFrames       int      `json:"video_frames"` // SVD/WAN
	FPS               int      `json:"fps"`          // SVD/WAN
	CacheMode         string   `json:"cache_mode"`   // easycache, ucache, etc.
	CacheOption       string   `json:"cache_option"` // threshold=0.25, etc.
}

// ModelLoadOptions describes how to load a model, including modular sub-models.
type ModelLoadOptions struct {
	ModelName      string `json:"model"`
	VAEPath        string `json:"vae"`
	ClipLPath      string `json:"clip_l"`
	ClipGPath      string `json:"clip_g"`
	T5XXLPath      string `json:"t5xxl"`
	DiffusionModel string `json:"diffusion_model"`
}

// ModelInfo describes a discovered model bundle.
type ModelInfo struct {
	Name           string            `json:"name"`
	Path           string            `json:"path"`
	Family         string            `json:"family"`       // "flux", "sdxl", "sd3", "sd15"
	Capabilities   []string          `json:"capabilities"` // ["txt2img", "img2img", "inpaint"]
	ResolutionHint string            `json:"resolution_hint"`
	Files          map[string]string `json:"files"`         // family-specific file map
	Format         string            `json:"format"`        // "safetensors", "gguf", "ckpt", etc.
	Compatibility  string            `json:"compatibility"` // "verified", "unknown", "incompatible"
}

// ImageResult is returned after a successful generation.
type ImageResult struct {
	Path      string        `json:"path"`
	ThumbPath string        `json:"thumb_path"`
	Metadata  ImageMetadata `json:"metadata"`
	ImageData string        `json:"image_data,omitempty"` // base64 PNG data from sd-server
}

// ImageSource links a generated image to an entity.
type ImageSource struct {
	WorldID     string `json:"world_id,omitempty"`
	ScenarioID  string `json:"scenario_id,omitempty"`
	CharacterID string `json:"character_id,omitempty"`
}

// ImageMetadata is persisted alongside every generated image.
type ImageMetadata struct {
	File           string      `json:"file"`
	Prompt         string      `json:"prompt"`
	NegativePrompt string      `json:"negative_prompt"`
	Model          string      `json:"model"`
	Family         string      `json:"family"`
	Width          int         `json:"width"`
	Height         int         `json:"height"`
	Steps          int         `json:"steps"`
	Seed           int64       `json:"seed"`
	Sampler        string      `json:"sampler"`
	CfgScale       float64     `json:"cfg_scale"`
	Mode           string      `json:"mode"`
	CreatedAt      string      `json:"created_at"` // RFC3339 format, Wails-safe
	Source         ImageSource `json:"source"`
}

// EngineConfig is passed to ImageEngine.Initialize.
type EngineConfig struct {
	BinaryPath string `json:"binary_path"` // path to sd-server or sd.exe
	ModelsDir  string `json:"models_dir"`  // root models directory
	OutputDir  string `json:"output_dir"`  // gallery output directory
	TempDir    string `json:"temp_dir"`    // temp files for img2img/inpaint
	Port       int    `json:"port"`        // sd-server HTTP port
}

// ProgressInfo is returned by GetProgress.
type ProgressInfo struct {
	Progress     float64 `json:"progress"`                // 0.0 .. 1.0
	Step         int     `json:"step"`                    // current sampling step
	TotalSteps   int     `json:"total_steps"`             // total sampling steps
	PreviewImage string  `json:"preview_image,omitempty"` // base64 PNG preview
}

// ImageJob tracks an async generation job.
type ImageJob struct {
	ID           string       `json:"id"`
	Status       string       `json:"status"`                  // "queued", "generating", "completed", "failed", "cancelled"
	Progress     float64      `json:"progress"`                // 0.0 .. 1.0
	Step         int          `json:"step"`                    // current sampling step
	TotalSteps   int          `json:"total_steps"`             // total sampling steps
	PreviewImage string       `json:"preview_image,omitempty"` // base64 PNG preview
	Options      ImageOptions `json:"options,omitempty"`
	Result       *ImageResult `json:"result,omitempty"`
	Error        string       `json:"error,omitempty"`
	CreatedAt    string       `json:"created_at"` // RFC3339 format, Wails-safe
}
