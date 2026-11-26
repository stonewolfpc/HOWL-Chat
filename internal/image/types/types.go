package types

// LoRAStrength represents a LoRA with its strength value.
// Strength range is typically -2.0 to 2.0, where:
//
//	0.0 = disabled (but kept in list for UI)
//	0.5 = reduced effect
//	1.0 = normal effect (default)
//	1.5 = amplified effect
//	2.0 = maximum effect
//	negative = inverts the effect
//
// Enabled allows toggling without removing from the list.
type LoRAStrength struct {
	Name     string  `json:"name"`
	Strength float64 `json:"strength"`
	Enabled  bool    `json:"enabled"`
}

// Image engine profiles. The default is the pinned engine; the IP profile is a
// compatibility bridge for builds that expose native IP-Adapter but not DGL.
const (
	EngineProfileDefault   = "sdcpp-3b12781"
	EngineProfileIPAdapter = "sdcpp-c6dfa9d-ip"
)

// ImageOptions describes a single generation request.
type ImageOptions struct {
	Prompt           string  `json:"prompt"`
	NegativePrompt   string  `json:"negative_prompt"`
	Width            int     `json:"width"`
	Height           int     `json:"height"`
	Steps            int     `json:"steps"`
	Seed             int64   `json:"seed"`
	Sampler          string  `json:"sampler"`
	Scheduler        string  `json:"scheduler"`
	CfgScale         float64 `json:"cfg_scale"`
	Guidance         float64 `json:"guidance"`        // distilled guidance (SD3/Flux)
	Eta              float64 `json:"eta"`             // noise multiplier
	Strength         float64 `json:"strength"`        // img2img / inpaint
	InitImagePath    string  `json:"init_image_path"` // img2img / inpaint
	MaskPath         string  `json:"mask_path"`       // inpaint only
	ModelName        string  `json:"model_name"`
	VAEPath          string  `json:"vae_path"`           // standalone VAE
	ClipLPath        string  `json:"clip_l_path"`        // CLIP-L
	ClipGPath        string  `json:"clip_g_path"`        // CLIP-G
	T5XXLPath        string  `json:"t5xxl_path"`         // T5XXL
	ControlNetPath   string  `json:"control_net_path"`   // ControlNet model
	ControlImagePath string  `json:"control_image_path"` // ControlNet image
	ControlStrength  float64 `json:"control_strength"`   // ControlNet strength
	// IP-Adapter is image-prompt conditioning, not ControlNet. The adapter and
	// CLIP-Vision encoder are startup slots; the reference image and strength
	// are per-request. EngineProfile selects the binary that implements them.
	IPAdapterPath      string  `json:"ip_adapter_path,omitempty"`
	ClipVisionPath     string  `json:"clip_vision_path,omitempty"`
	IPAdapterImagePath string  `json:"ip_adapter_image_path,omitempty"`
	IPAdapterStrength  float64 `json:"ip_adapter_strength,omitempty"`
	ClipVisionCPU      bool    `json:"clip_vision_cpu,omitempty"`
	EngineProfile      string  `json:"engine_profile,omitempty"`
	// UpscalerPath is a DGL/ESRGAN upscaler for a workflow stage. It is not part
	// of a native generation request; the workflow uses it to submit its own DGL
	// job, and the path's directory is added to the engine allowlist at startup.
	UpscalerPath string `json:"upscaler_path,omitempty"`
	UpscalerDir  string `json:"upscaler_dir,omitempty"`
	UpscalerTile int    `json:"upscaler_tile,omitempty"`
	BatchCount   int    `json:"batch_count"`
	Mode         string `json:"mode"` // "txt2img", "img2img", "inpaint", "outpaint", "upscale", "video"
	// JobID identifies the job that owns this request, so a cancellation can be
	// scoped to it. Empty for callers that are not running a job.
	JobID             string         `json:"job_id,omitempty"`
	OutpaintMargin    int            `json:"outpaint_margin,omitempty"`
	OutpaintDirection string         `json:"outpaint_direction,omitempty"` // "all", "left", "right", "up", "down"
	UpscaleFactor     int            `json:"upscale_factor,omitempty"`     // 2 or 4
	LoRAs             []LoRAStrength `json:"loras,omitempty"`              // LoRA with individual strength values
	// LegacyLoRAs is for backward compatibility during migration
	LegacyLoRAs         []string `json:"-"`
	EnableHR            bool     `json:"enable_hr"`
	HRScale             float64  `json:"hr_scale"`
	HRUpscaler          string   `json:"hr_upscaler"`
	DenoisingStrength   float64  `json:"denoising_strength"`
	ResizeMode          string   `json:"resize_mode,omitempty"`    // just_resize, crop, fill, stretch
	StyleStrength       float64  `json:"style_strength,omitempty"` // 0-1, how much original style influences output
	PreserveColors      bool     `json:"preserve_colors,omitempty"`
	MatchHistogram      bool     `json:"match_histogram,omitempty"`
	LatentNoise         float64  `json:"latent_noise,omitempty"` // 0-1, controlled chaos injection
	MaskBlur            int      `json:"mask_blur,omitempty"`
	MaskInvert          bool     `json:"mask_invert,omitempty"`
	MaskFeather         int      `json:"mask_feather,omitempty"`          // inpaint mask edge feathering
	MaskedContent       string   `json:"masked_content,omitempty"`        // fill, preserve, replace, only_modify
	EdgeBehavior        string   `json:"edge_behavior,omitempty"`         // mirror, blur, noise, solid
	UpscaleMethod       string   `json:"upscale_method,omitempty"`        // latent_sd, latent_sdxl, esrgan, realesrgan, swinir, etc.
	UpscaleSharpen      string   `json:"upscale_sharpen,omitempty"`       // none, mild, medium, strong
	UpscaleAntiArtifact string   `json:"upscale_anti_artifact,omitempty"` // off, mild, strong
	Tiling              bool     `json:"tiling"`
	ClipSkip            int      `json:"clip_skip"`
	VideoFrames         int      `json:"video_frames"`                    // SVD/WAN
	FPS                 int      `json:"fps"`                             // SVD/WAN
	VideoMotionStrength float64  `json:"video_motion_strength,omitempty"` // 0-2, how much changes between frames
	VideoSeedMode       string   `json:"video_seed_mode,omitempty"`       // fixed, per_frame, interpolated
	VideoDenoise        float64  `json:"video_denoise,omitempty"`         // per-frame denoise strength
	VideoLoop           bool     `json:"video_loop,omitempty"`            // seamless loop toggle
	OutputFormat        string   `json:"output_format,omitempty"`         // mp4, gif, png_seq, webp
	EndImagePath        string   `json:"end_image_path,omitempty"`        // video end frame
	CacheMode           string   `json:"cache_mode"`                      // easycache, ucache, etc.
	CacheOption         string   `json:"cache_option"`                    // threshold=0.25, etc.
	Prediction          string   `json:"prediction,omitempty"`            // eps, v, edm_v, sd3_flow, flux_flow, flux2_flow
	RNG                 string   `json:"rng,omitempty"`                   // std_default, cuda, cpu
	SamplerRNG          string   `json:"sampler_rng,omitempty"`           // std_default, cuda, cpu
	FlowShift           float64  `json:"flow_shift,omitempty"`            // shift for Flow models
	VAETileOverlap      float64  `json:"vae_tile_overlap,omitempty"`      // 0.0-1.0
	DisableMetadata     bool     `json:"disable_metadata,omitempty"`      // don't embed metadata
	CfgRescale          float64  `json:"cfg_rescale,omitempty"`           // prevents overbaked contrast
	FaceRestore         float64  `json:"face_restore,omitempty"`          // CodeFormer/GFPGAN strength
	SChurn              float64  `json:"s_churn,omitempty"`               // DPM++ sampler shaping
	STmin               float64  `json:"s_tmin,omitempty"`
	STmax               float64  `json:"s_tmax,omitempty"`
	SNoise              float64  `json:"s_noise,omitempty"`
	// Model startup flags — forwarded to lazy-load if engine is not yet running
	VAETiling          bool   `json:"vae_tiling,omitempty"`
	TemporalTiling     bool   `json:"temporal_tiling,omitempty"`
	FlashAttention     bool   `json:"flash_attention,omitempty"`      // legacy boolean; use FlashAttentionMode for new callers
	FlashAttentionMode string `json:"flash_attention_mode,omitempty"` // auto, on, off
	OffloadToCPU       bool   `json:"offload_to_cpu,omitempty"`
	Backend            string `json:"backend,omitempty"`        // compute backend (cpu, cuda, ...)
	ParamsBackend      string `json:"params_backend,omitempty"` // weight placement (cpu, cuda, disk)
	Circular           bool   `json:"circular,omitempty"`       // video temporal circular padding
	// WeightResidencyTTL is sd-server's idle-release timer. Nil keeps the engine
	// default (disabled). A request forwards it so a lazily started engine uses
	// the same value as an explicit load. Staging TTL is intentionally not
	// modelled: the engine releases staged copies when a component finishes, so
	// it has no generation-boundary effect to expose.
	WeightResidencyTTL *int `json:"weight_residency_ttl,omitempty"`
	// Hardware/engine flags. Nil and 0 both mean Auto to the image start path;
	// Howl applies a positive post-reserve budget so the 10% reserve is kept even
	// though sd-server's own -1 default reserves only 1 GiB. A negative value is
	// passed through as the engine's own reserve.
	MaxVRAM               *int   `json:"max_vram,omitempty"`
	WeightType            string `json:"weight_type,omitempty"`             // --type (f16, q8_0, q4_0, ...)
	ConditioningCacheSize int    `json:"conditioning_cache_size,omitempty"` // --conditioning-cache-size
	MMap                  bool   `json:"mmap,omitempty"`                    // --mmap
}

// ModelLoadOptions describes how to load a model, including modular sub-models.
type ModelLoadOptions struct {
	ModelName string `json:"model"`
	ModelPath string `json:"model_path,omitempty"` // absolute path; when set, skips directory scan
	// JobID identifies the job that needs this load, so one job's cancellation
	// cannot abort another job's load. Empty for a load started by the UI.
	JobID string `json:"job_id,omitempty"`
	// Family is the architecture the image registry resolved for this model.
	// The engine can only see its own ModelsDir, so without it a checkpoint that
	// lives in another scanned folder loses its family and is guessed from the
	// filename — which is how a Z-Image GGUF was passed as a full checkpoint.
	Family         string `json:"family,omitempty"`
	VAEPath        string `json:"vae"`
	ClipLPath      string `json:"clip_l"`
	ClipGPath      string `json:"clip_g"`
	T5XXLPath      string `json:"t5xxl"`
	LLMPath        string `json:"llm"` // For Anima's Qwen3 text encoder
	DiffusionModel string `json:"diffusion_model"`
	// ControlNetPath is an explicit startup ControlNet (--control-net). It is a
	// process-wide slot, unlike a per-request control image, and a changed path
	// requires a reload because the engine only reads it at startup.
	ControlNetPath string `json:"control_net_path,omitempty"`
	// IP-Adapter and CLIP-Vision are startup slots on the native-IP engine.
	IPAdapterPath  string `json:"ip_adapter_path,omitempty"`
	ClipVisionPath string `json:"clip_vision_path,omitempty"`
	// ClipVisionCPU keeps the CLIP-Vision encoder's params and compute on CPU
	// (--backend clip_vision=cpu). It is selected from the measured/projected
	// budget, not as a blanket default.
	ClipVisionCPU bool `json:"clip_vision_cpu,omitempty"`
	// EngineProfile selects the engine capability profile for this load. Empty
	// means the default pinned profile; it is never inherited from a previous
	// IP request.
	EngineProfile string `json:"engine_profile,omitempty"`
	// EngineBinaryPath is the resolved runtime path for EngineProfile. It is
	// never persisted as user intent.
	EngineBinaryPath string `json:"-"`
	// UpscalerDir is the directory containing DGL upscaler weights. This
	// sd-server build authorizes DGL file access through model roots, not through
	// --hires-upscalers-dir, so the directory is also passed as --model-dir when
	// it differs from ModelsDir. It is a startup option because the allowlist is
	// fixed when the server starts.
	UpscalerDir   string `json:"upscaler_dir,omitempty"`
	Backend       string `json:"backend,omitempty"`        // cpu, cuda for this build
	ParamsBackend string `json:"params_backend,omitempty"` // cpu, cuda, disk
	CacheMode     string `json:"cache_mode,omitempty"`     // easycache, ucache, etc.
	CacheOption   string `json:"cache_option,omitempty"`   // threshold=0.25, etc.
	VAETiling     bool   `json:"vae_tiling,omitempty"`
	// ControlNetCPU keeps the startup ControlNet's params and compute on CPU
	// (--backend controlnet=cpu). The pinned engine otherwise stages the full
	// 1415.82 MiB weight set to CUDA during sampling, which cannot coexist with
	// the 10% reserve on the 8 GiB card. Set only when admission projects that
	// staging would cross the post-reserve budget.
	ControlNetCPU      bool   `json:"control_net_cpu,omitempty"`
	TemporalTiling     bool   `json:"temporal_tiling,omitempty"`
	Circular           bool   `json:"circular,omitempty"`
	FlashAttention     bool   `json:"flash_attention,omitempty"`      // legacy boolean; empty mode maps true→on, false→auto
	FlashAttentionMode string `json:"flash_attention_mode,omitempty"` // auto, on, off
	OffloadToCPU       bool   `json:"offload_to_cpu,omitempty"`
	DisableMetadata    bool   `json:"disable_metadata,omitempty"`
	// MaxVRAM is sd-server's --max-vram budget in GiB. Nil and 0 mean Auto for a
	// generation; Howl derives a positive post-reserve cap for both. Negative
	// values pass through as the engine's own reserve.
	MaxVRAM *int `json:"max_vram,omitempty"`
	// AutoMaxVRAM marks a runtime budget Howl derived from free VRAM for an
	// offloaded start. It is never saved as user intent and does not replace the
	// 10% default reserve the way an explicit --max-vram does.
	AutoMaxVRAM bool `json:"-"`
	// AutoPlacement marks a placement Howl selected from the requested budget
	// rather than the user's Backend/Params Backend. It is runtime state, not a
	// saved preference: the UI still shows Auto, and a later generation carries
	// the engine's actual placement forward instead of restarting every time.
	AutoPlacement bool `json:"-"`
	// AutoPlacementMode is the selected bounded mode: "cpu-diffusion" keeps CLIP
	// and the VAE on the GPU while diffusion runs on CPU, and "cpu-only" uses
	// no device compute. Empty means no automatic placement was applied.
	AutoPlacementMode string `json:"-"`
	// WeightType is sd-server's --type: weights are quantized to this type while
	// loading (f16, q8_0, q4_0, q5_0, q5_1). Empty leaves the file's own type.
	WeightType string `json:"weight_type,omitempty"`
	// ConditioningCacheSize is --conditioning-cache-size: text-encoder results
	// kept per model context. 0 keeps sd.cpp's own default.
	ConditioningCacheSize int `json:"conditioning_cache_size,omitempty"`
	// MMap enables memory-mapped model loading (faster load, may impact gen speed if paging occurs).
	MMap bool `json:"mmap,omitempty"`
	// Prediction overrides the model's auto-detected prediction type at startup.
	// Required for Flux GGUF single-file checkpoints where sd.cpp cannot detect
	// the architecture from the file. Values: eps, v, edm_v, sd3_flow, flux_flow, flux2_flow.
	// Leave empty for auto-detection (safetensors/split-component bundles).
	Prediction string `json:"prediction,omitempty"`

	// WeightResidencyTTL is sd-server's --weight-residency-ttl: idle seconds
	// before model components are unloaded from RAM and VRAM (-1 disabled, 0
	// immediate). Nil keeps the engine's own default (-1), so this is opt-in and
	// only ever set once a measurement justifies it.
	WeightResidencyTTL *int `json:"weight_residency_ttl,omitempty"`
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
	Path           string          `json:"path"`
	ThumbPath      string          `json:"thumb_path"`
	Metadata       ImageMetadata   `json:"metadata"`
	ImageData      string          `json:"image_data,omitempty"`      // base64 PNG data from sd-server (first image for compat)
	Images         []string        `json:"images,omitempty"`          // all images when batch_count > 1
	VideoData      string          `json:"video_data,omitempty"`      // base64 clip for video mode
	VideoMIME      string          `json:"video_mime,omitempty"`      // container type for the data URL
	ImageMetadatas []ImageMetadata `json:"image_metadatas,omitempty"` // per-image metadata for batches
	// GalleryError reports that generation succeeded but the gallery copy did
	// not. The image bytes above remain usable; this is a partial result.
	GalleryError string `json:"gallery_error,omitempty"`
}

// ImageSource links a generated image to an entity.
type ImageSource struct {
	WorldID     string `json:"world_id,omitempty"`
	ScenarioID  string `json:"scenario_id,omitempty"`
	CharacterID string `json:"character_id,omitempty"`
}

// ImageMetadata is persisted alongside every generated image.
type ImageMetadata struct {
	File           string       `json:"file"`
	Prompt         string       `json:"prompt"`
	NegativePrompt string       `json:"negative_prompt"`
	Model          string       `json:"model"`
	Family         string       `json:"family"`
	Width          int          `json:"width"`
	Height         int          `json:"height"`
	Steps          int          `json:"steps"`
	Seed           int64        `json:"seed"`
	Sampler        string       `json:"sampler"`
	CfgScale       float64      `json:"cfg_scale"`
	Mode           string       `json:"mode"`
	CreatedAt      string       `json:"created_at"` // RFC3339 format, Wails-safe
	Source         ImageSource  `json:"source"`
	Options        ImageOptions `json:"options,omitempty"` // full generation parameters for replay
}

// EngineConfig is passed to ImageEngine.Initialize.
type EngineConfig struct {
	BinaryPath string `json:"binary_path"` // path to sd-server or sd.exe
	// IPAdapterBinaryPath is the optional native-IP engine profile. It is a
	// compatibility bridge, not a second resident engine; the default engine is
	// still BinaryPath.
	IPAdapterBinaryPath string `json:"ip_adapter_binary_path,omitempty"`
	ModelsDir           string `json:"models_dir"` // root models directory
	OutputDir           string `json:"output_dir"` // gallery output directory
	TempDir             string `json:"temp_dir"`   // temp files for img2img/inpaint
	Port                int    `json:"port"`       // sd-server HTTP port
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
	Status       string       `json:"status"`                  // "queued", "loading", "generating", "completed", "failed", "cancelled"
	Progress     float64      `json:"progress"`                // 0.0 .. 1.0
	Step         int          `json:"step"`                    // current sampling step
	TotalSteps   int          `json:"total_steps"`             // total sampling steps
	PreviewImage string       `json:"preview_image,omitempty"` // base64 PNG preview
	Options      ImageOptions `json:"options,omitempty"`
	Result       *ImageResult `json:"result,omitempty"`
	Error        string       `json:"error,omitempty"`
	CreatedAt    string       `json:"created_at"` // RFC3339 format, Wails-safe
}
