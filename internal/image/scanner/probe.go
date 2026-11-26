package scanner

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Image weights used to be identified by filename and size alone. That works
// for the names people type by hand and fails exactly where it matters:
//
//   - Invoke AI stores every model as <uuid>/model.safetensors, so the file is
//     called "model" and the folder is a UUID. Filename heuristics fall through
//     to the conservative SD1.5 default, and a 6.6 GB SDXL merge or a Flux
//     transformer is then loaded as SD1.5 and rendered as noise.
//   - Renamed exports ("diffusion_pytorch_model.fp16") carry no family either.
//   - Text encoders and VAEs are larger than the size thresholds expect, so
//     they were listed as diffusion models.
//
// GGUF chat models never had this problem because the file is read: the header
// names the architecture, the quantisation and the tensor layout. WeightProbe
// gives image weights the same treatment. It reads the safetensors header once,
// inspects __metadata__ and the tensor names/shapes, and consults the Diffusers
// config files that sit beside it. The family then comes from the weights, not
// from a guess.
type WeightProbe struct {
	Metadata      map[string]string // __metadata__ as written by the trainer
	TensorCount   int
	ClassName     string   // sibling config.json: _class_name
	Architectures []string // sibling config.json: architectures
	Pipeline      string   // enclosing model_index.json: _class_name

	// Structural markers, set while scanning tensor names.
	DoubleBlocks  bool // Flux / Chroma mmdit blocks
	SingleBlocks  bool
	JointBlocks   bool // SD3 mmdit
	DistilledGuid bool // Chroma guidance layer
	NoiseRefiner  bool // Z-Image refiner blocks
	Adaln         bool // Anima net.blocks adaLN modulation
	LabelEmb      bool // SDXL pooled embedding projection
	CondEmbedders bool // SDXL dual text encoder stack
	ClipL         bool // CLIP-L text encoder inside the checkpoint
	InputBlocks   bool // classic UNet
	FirstStage    bool // first_stage_model inside a checkpoint
	EncoderDown   bool
	DecoderUp     bool
	QuantConv     bool
	Lora          bool
	ControlNet    bool
	TextModel     bool
	LogitScale    bool
	QwenLayers    bool
	T5Encoder     bool
	VisionModel   bool

	// CrossAttnDim is the projection width of attn2.to_k, the single most
	// reliable way to tell SD1.5 (768) from SD2.1 (1024) from SDXL (2048).
	CrossAttnDim int
}

const (
	// maxSafetensorsHeader bounds the header read. Real headers are a few MB;
	// anything larger is treated as unreadable rather than buffered.
	maxSafetensorsHeader = 64 << 20
	// maxProbeKeys bounds how many tensor names are inspected.
	maxProbeKeys = 200000
)

var uuidDirRe = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// ProbeWeight reads whatever the file itself says about the model. It never
// fails: an unreadable header simply yields an empty probe, and callers fall
// back to the filename logic they used before.
func ProbeWeight(path string) WeightProbe {
	p := WeightProbe{Metadata: map[string]string{}}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".safetensors":
		p.inspectSafetensors(path)
	case ".gguf":
		p.inspectGGUF(path)
	}
	p.inspectConfigs(path)
	return p
}

func (p *WeightProbe) inspectSafetensors(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	var size uint64
	if err := binary.Read(f, binary.LittleEndian, &size); err != nil {
		return
	}
	if size < 2 || size > maxSafetensorsHeader {
		return
	}
	buf := make([]byte, size)
	if _, err := io.ReadFull(f, buf); err != nil {
		return
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(buf, &raw); err != nil {
		return
	}

	if m, ok := raw["__metadata__"]; ok {
		var meta map[string]interface{}
		if json.Unmarshal(m, &meta) == nil {
			for k, v := range meta {
				p.Metadata[k] = fmt.Sprint(v)
			}
		}
	}

	for key, val := range raw {
		if key == "__metadata__" {
			continue
		}
		p.TensorCount++
		if p.TensorCount > maxProbeKeys {
			return
		}
		p.noteKey(key)
		// The first attn2 projection is representative: every block in a family
		// shares the same cross-attention width.
		if p.CrossAttnDim == 0 && strings.Contains(key, "attn2.to_k") {
			var t struct {
				Shape []int `json:"shape"`
			}
			if json.Unmarshal(val, &t) == nil && len(t.Shape) >= 2 && t.Shape[1] > 0 {
				p.CrossAttnDim = t.Shape[1]
			}
		}
	}
}

func (p *WeightProbe) noteKey(key string) {
	k := strings.ToLower(key)
	switch {
	case strings.Contains(k, "double_blocks"):
		p.DoubleBlocks = true
	case strings.Contains(k, "single_blocks"):
		p.SingleBlocks = true
	case strings.Contains(k, "joint_blocks"):
		p.JointBlocks = true
	case strings.Contains(k, "distilled_guidance"):
		p.DistilledGuid = true
	case strings.Contains(k, "noise_refiner"):
		p.NoiseRefiner = true
	// Anima names its blocks net.blocks.N and modulates them with adaLN;
	// nothing else in the image zoo uses both, so the pair is the signature.
	case strings.Contains(k, "net.blocks.") && strings.Contains(k, "adaln_modulation"):
		p.Adaln = true
	case strings.Contains(k, "label_emb"):
		p.LabelEmb = true
	case strings.Contains(k, "conditioner.embedders.1"):
		p.CondEmbedders = true
	case strings.Contains(k, "conditioner.embedders.0"):
		p.ClipL = true
	case strings.Contains(k, "input_blocks"):
		p.InputBlocks = true
	case strings.Contains(k, "first_stage_model"):
		p.FirstStage = true
	case strings.Contains(k, "encoder.down."), strings.Contains(k, "encoder.downsamples"):
		p.EncoderDown = true
	case strings.Contains(k, "decoder.up."), strings.Contains(k, "decoder.upsamples"):
		p.DecoderUp = true
	case strings.Contains(k, "quant_conv"):
		p.QuantConv = true
	case strings.Contains(k, "controlnet"), strings.Contains(k, "control_blocks"),
		strings.Contains(k, "control_all_x_embedder"), strings.Contains(k, "control_layers"),
		strings.Contains(k, "zero_convs"):
		p.ControlNet = true
	case strings.Contains(k, "lora_unet"), strings.Contains(k, "lora_te"),
		strings.Contains(k, "lora_down"), strings.Contains(k, "lora_up"),
		strings.Contains(k, ".lora_a."), strings.Contains(k, ".lora_b."),
		strings.Contains(k, "hada_w1"), strings.Contains(k, "lycoris"):
		p.Lora = true
	case strings.Contains(k, "text_model.encoder"), strings.Contains(k, "text_model.embeddings"):
		p.TextModel = true
	case strings.Contains(k, "logit_scale"):
		p.LogitScale = true
	case strings.Contains(k, "model.embed_tokens"), strings.Contains(k, "model.layers."):
		p.QwenLayers = true
	case strings.Contains(k, "encoder.block.") && strings.Contains(k, "selfattention"):
		p.T5Encoder = true
	case strings.Contains(k, "vision_model"), strings.Contains(k, "vision_tower"):
		p.VisionModel = true
	}
}

// inspectConfigs reads the Diffusers sidecars that Invoke AI and save_pretrained
// write next to the weights.
func (p *WeightProbe) inspectConfigs(path string) {
	dir := filepath.Dir(path)
	if data, err := os.ReadFile(filepath.Join(dir, "config.json")); err == nil {
		var cfg struct {
			ClassName string   `json:"_class_name"`
			Archs     []string `json:"architectures"`
			CrossDim  int      `json:"cross_attention_dim"`
		}
		if json.Unmarshal(data, &cfg) == nil {
			p.ClassName = cfg.ClassName
			p.Architectures = cfg.Archs
			if p.CrossAttnDim == 0 {
				p.CrossAttnDim = cfg.CrossDim
			}
		}
	}
	// A component inside a Diffusers bundle sees the pipeline descriptor one
	// level up, which is the authority on the family.
	if data, err := os.ReadFile(filepath.Join(filepath.Dir(dir), "model_index.json")); err == nil {
		var mi struct {
			ClassName string `json:"_class_name"`
		}
		if json.Unmarshal(data, &mi) == nil {
			p.Pipeline = mi.ClassName
		}
	}
}

// metaString is the haystack of human- and machine-written family hints.
func (p WeightProbe) metaString() string {
	parts := []string{
		p.Metadata["modelspec.architecture"],
		p.Metadata["ss_base_model_version"],
		p.Metadata["modelspec.title"],
		p.Metadata["ss_sd_model_name"],
		p.Metadata["modelspec.implementation"],
		p.Metadata["general.architecture"],
		p.Pipeline,
		p.ClassName,
		strings.Join(p.Architectures, " "),
	}
	return strings.ToLower(strings.Join(parts, " "))
}

// Title returns the name the model was published under, if the weights record
// one. This is what turns <uuid>/model.safetensors back into a readable name.
func (p WeightProbe) Title() string {
	for _, k := range []string{"modelspec.title", "ss_sd_model_name", "general.name"} {
		if v := strings.TrimSpace(p.Metadata[k]); v != "" {
			return v
		}
	}
	return ""
}

// AssetType classifies the file from its contents.
func (p WeightProbe) AssetType() (AssetType, bool) {
	meta := p.metaString()
	hasMMDiT := p.DoubleBlocks || p.SingleBlocks || p.JointBlocks

	switch {
	case p.Lora || contains(meta, "lora"):
		return AssetLoRA, true
	case p.ControlNet:
		return AssetControlNet, true
	// A standalone VAE is nothing but encoder/decoder/quant_conv. A checkpoint
	// carries first_stage_model *inside* the diffusion weights, so it must not
	// be reported as a VAE.
	// Flux/Qwen VAEs name their halves downsamples/upsamples and carry no
	// quant_conv, so requiring both old markers missed them entirely and the
	// encoder was offered as a diffusion model instead of a VAE.
	case (p.FirstStage || (p.EncoderDown && p.DecoderUp) || (p.EncoderDown && p.QuantConv)) && !p.InputBlocks && !hasMMDiT:
		return AssetVAE, true
	case contains(meta, "t5encoder", "t5xxl", "t5-xxl"), p.T5Encoder:
		return AssetT5XXL, true
	case contains(meta, "cliptextmodelwithprojection"):
		return AssetClipG, true
	case contains(meta, "cliptextmodel"):
		return AssetClipL, true
	// Qwen-Image, Z-Image and their derivatives are diffusion models, not
	// language models: their GGUF architecture string just happens to say qwen.
	case contains(meta, "qwen_image", "qwen-image", "qwenimage", "z_image", "z-image", "zimage"):
		return AssetDiffusion, true
	// A GGUF header that declares a language architecture is a language
	// model, whatever folder it sits in: llama-q2_k.gguf was being listed
	// as an SD1.5 checkpoint because its filename said nothing useful.
	case contains(meta, "siglip", "qwen", "llm", "llama", "gemma", "mistral", "mixtral", "deepseek", "phi"), p.QwenLayers:
		return AssetLLM, true
	case p.InputBlocks || hasMMDiT || p.LabelEmb || p.CondEmbedders || p.NoiseRefiner || p.Adaln:
		return AssetDiffusion, true
	case p.TextModel && p.LogitScale:
		return AssetClipG, true
	case p.TextModel:
		return AssetClipL, true
	// CLIP-Vision is an image encoder, not a language model. A vision-only
	// tower is the IP-Adapter image encoder; a multimodal LLM also carries a
	// text model and stays on the LLM path below.
	case p.VisionModel && !p.TextModel:
		return AssetClipVision, true
	case p.VisionModel:
		return AssetLLM, true
	}
	// A declared architecture is enough to call it a diffusion model: Flux.2's
	// blocks are named differently from Flux.1's, but the header is explicit.
	if _, ok := p.Family(); ok {
		return AssetDiffusion, true
	}
	return AssetUnknown, false
}

// Family reports the architecture the weights are built for.
func (p WeightProbe) Family() (ModelFamily, bool) {
	meta := p.metaString()

	// Explicit declarations win: they are written by the trainer and survive
	// renames, which is precisely the case filenames fail at.
	switch {
	// More specific names first: flux2 must not be read as flux.
	case contains(meta, "flux2", "flux.2", "flux-2"):
		return FamilyFlux2, true
	case contains(meta, "z_image", "z-image", "zimage"):
		return FamilyZImage, true
	case contains(meta, "qwen_image", "qwen-image", "qwenimage"):
		return FamilyQwenImage, true
	case p.Adaln:
		return FamilyAnima, true
	case contains(meta, "flux"):
		return FamilyFlux, true
	case contains(meta, "chroma"):
		return FamilyChroma, true
	case contains(meta, "stable-diffusion-3.5", "sd3.5", "sd_3.5", "sd35", "sd3-5"):
		return FamilySD35, true
	case contains(meta, "stable-diffusion-3", "sd_3", "sd3"):
		return FamilySD3, true
	case contains(meta, "sdxl", "sd_xl", "stable-diffusion-xl"):
		return FamilySDXL, true
	case contains(meta, "stable-diffusion-2", "sd2.1", "sd_2.1"):
		return FamilySD21, true
	case contains(meta, "stable-diffusion-v1", "sd1.5", "sd_1.5", "sd-v1"):
		return FamilySD15, true
	case contains(meta, "wan"):
		return FamilyWan, true
	}

	// No declaration: read the layout.
	switch {
	case p.NoiseRefiner:
		return FamilyZImage, true
	case p.JointBlocks:
		if p.CrossAttnDim >= 2432 {
			return FamilySD35, true
		}
		return FamilySD3, true
	case p.DoubleBlocks || p.SingleBlocks:
		if p.DistilledGuid {
			return FamilyChroma, true
		}
		return FamilyFlux, true
	case p.LabelEmb || p.CondEmbedders || p.CrossAttnDim == 2048:
		return FamilySDXL, true
	case p.InputBlocks && p.CrossAttnDim == 1024:
		return FamilySD21, true
	case p.InputBlocks && p.CrossAttnDim == 768:
		return FamilySD15, true
	case p.InputBlocks:
		// Classic UNet with an unreadable projection width: SD1.5 is the only
		// safe reading, and callers may still upgrade it by size.
		return FamilySD15, true
	}
	return FamilyUnknown, false
}

// familyLabel is the short human name used when a model has no usable title.
func familyLabel(f ModelFamily) string {
	switch f {
	case FamilySD15:
		return "SD1.5"
	case FamilySD21:
		return "SD2.1"
	case FamilySDXL:
		return "SDXL"
	case FamilyTurbo:
		return "SDXL Turbo"
	case FamilySD3:
		return "SD3"
	case FamilySD35:
		return "SD3.5"
	case FamilyFlux:
		return "Flux"
	case FamilyFlux2:
		return "Flux.2"
	case FamilyZImage:
		return "Z-Image"
	case FamilyQwenImage:
		return "Qwen-Image"
	case FamilyChroma:
		return "Chroma"
	case FamilyKolors:
		return "Kolors"
	case FamilyCascade:
		return "Cascade"
	case FamilyAuraFlow:
		return "AuraFlow"
	case FamilyHiDream:
		return "HiDream"
	case FamilyWan:
		return "Wan"
	case FamilyLTX:
		return "LTX"
	case FamilyHunyuan:
		return "Hunyuan"
	case FamilyAnimate:
		return "AnimateDiff"
	case FamilyAnima:
		return "Anima"
	}
	return ""
}

// genericStems are the names a model ends up with when an exporter, not a
// person, chose it. They say nothing, so they are never shown as a name.
var genericStems = map[string]bool{
	"model":                        true,
	"model.fp16":                   true,
	"model.bf16":                   true,
	"diffusion_pytorch_model":      true,
	"diffusion_pytorch_model.fp16": true,
	"diffusion_pytorch_model.bf16": true,
	"pytorch_model":                true,
	"pytorch_model.fp16":           true,
	"diffusion_model":              true,
	"unet":                         true,
	"vae":                          true,
	"text_encoder":                 true,
	"text_encoder_2":               true,
}

// DisplayName turns a file into something a person can pick from a menu.
// Order: the published title, then the file stem, then — when the exporter
// named the file "model" inside a UUID folder — family plus a short folder id,
// which is at least stable and unique.
func DisplayName(stem, dir string, probe WeightProbe, family ModelFamily) string {
	if title := probe.Title(); title != "" {
		return title
	}
	key := strings.ToLower(stem)
	// A UUID is not a name either: Invoke AI names its folders that way.
	if !genericStems[key] && !uuidDirRe.MatchString(stem) {
		return stem
	}
	base := filepath.Base(filepath.Clean(dir))
	label := familyLabel(family)
	if uuidDirRe.MatchString(base) {
		short := base[:8]
		if label != "" {
			return label + " " + short
		}
		return short
	}
	if label != "" {
		return base + " " + label
	}
	return base
}
