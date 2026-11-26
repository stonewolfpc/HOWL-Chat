package scanner

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// writeSafetensors writes a syntactically valid safetensors header. The probe
// only ever reads the header, so the tensor data itself is not written.
func writeSafetensors(t *testing.T, path string, meta map[string]string, tensors map[string][]int) {
	t.Helper()
	header := map[string]interface{}{"__metadata__": meta}
	for name, shape := range tensors {
		header[name] = map[string]interface{}{
			"dtype":        "F16",
			"shape":        shape,
			"data_offsets": []int{0, 0},
		}
	}
	blob, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("marshal header: %v", err)
	}
	buf := make([]byte, 8, 8+len(blob))
	binary.LittleEndian.PutUint64(buf, uint64(len(blob)))
	buf = append(buf, blob...)
	if err := os.WriteFile(path, buf, 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// An Invoke AI main model: the file is called "model.safetensors" inside a UUID
// folder, and the only family evidence is the tensor layout.
func TestProbeIdentifiesFluxInsideUUIDFolder(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "6dfe9131-c2e7-4e6d-8883-1f378e45a103")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(dir, "model.safetensors")
	writeSafetensors(t, path, map[string]string{"format": "pt"}, map[string][]int{
		"double_blocks.0.img_attn.proj.weight": {3072, 3072},
		"single_blocks.0.linear1.weight":       {12288, 3072},
	})

	probe := ProbeWeight(path)
	if got, ok := probe.Family(); !ok || got != FamilyFlux {
		t.Fatalf("family = %q (ok=%v), want %q", got, ok, FamilyFlux)
	}
	if got, ok := probe.AssetType(); !ok || got != AssetDiffusion {
		t.Fatalf("type = %q (ok=%v), want %q", got, ok, AssetDiffusion)
	}
	if name := DisplayName("model", dir, probe, FamilyFlux); name != "Flux 6dfe9131" {
		t.Fatalf("display name = %q, want %q", name, "Flux 6dfe9131")
	}
}

// SDXL is recognised by label_emb and a 2048-wide cross-attention even when the
// filename says nothing at all. The old code called this SD1.5.
func TestProbeIdentifiesSDXLFromStructure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "unholyDesireMixSinister_v90.safetensors")
	writeSafetensors(t, path, nil, map[string][]int{
		"model.diffusion_model.input_blocks.4.1.transformer_blocks.0.attn2.to_k.weight": {320, 2048},
		"model.diffusion_model.label_emb.0.0.weight":                                    {2816, 1280},
	})

	probe := ProbeWeight(path)
	if got, ok := probe.Family(); !ok || got != FamilySDXL {
		t.Fatalf("family = %q (ok=%v), want %q", got, ok, FamilySDXL)
	}
	if probe.CrossAttnDim != 2048 {
		t.Fatalf("cross attention dim = %d, want 2048", probe.CrossAttnDim)
	}
}

// modelspec metadata is the trainer's own declaration and must win.
func TestProbePrefersDeclaredArchitecture(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mystery.safetensors")
	writeSafetensors(t, path, map[string]string{
		"modelspec.architecture": "stable-diffusion-xl-v1-base/lora",
		"modelspec.title":        "Perfect Eyes XL",
		"ss_network_module":      "networks.lora",
	}, map[string][]int{
		"lora_unet_input_blocks_4_1_proj_in.lora_down.weight": {320, 4},
		"lora_unet_input_blocks_4_1_proj_in.lora_up.weight":   {4, 320},
	})

	probe := ProbeWeight(path)
	if got, ok := probe.Family(); !ok || got != FamilySDXL {
		t.Fatalf("family = %q, want %q", got, FamilySDXL)
	}
	if got, ok := probe.AssetType(); !ok || got != AssetLoRA {
		t.Fatalf("type = %q, want %q", got, AssetLoRA)
	}
	if name := DisplayName("mystery", "/tmp", probe, FamilySDXL); name != "Perfect Eyes XL" {
		t.Fatalf("display name = %q, want the modelspec title", name)
	}
}

// A CLIP-Vision image encoder is a vision tower with no text tower. It used to
// be reported as an LLM and then hidden from the image picker, which made the
// downloaded ViT-H encoder impossible to select.
func TestProbeIdentifiesClipVisionFromStructure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "model.safetensors")
	writeSafetensors(t, path, nil, map[string][]int{
		"vision_model.embeddings.class_embedding":               {1280},
		"vision_model.embeddings.patch_embedding.weight":        {1280, 3, 14, 14},
		"vision_model.encoder.layers.0.self_attn.q_proj.weight": {1280, 1280},
		"visual_projection.weight":                              {1024, 1280},
	})

	if got, ok := ProbeWeight(path).AssetType(); !ok || got != AssetClipVision {
		t.Fatalf("type = %q (ok=%v), want %q", got, ok, AssetClipVision)
	}

	scanner := New(dir)
	idx, err := scanner.Scan()
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(idx.ClipVisions) != 1 {
		t.Fatalf("clip visions = %d, want 1", len(idx.ClipVisions))
	}
}

// A standalone VAE is a VAE; a checkpoint that embeds first_stage_model is not.
func TestProbeDistinguishesVAEFromCheckpoint(t *testing.T) {
	dir := t.TempDir()
	vaePath := filepath.Join(dir, "ae.safetensors")
	writeSafetensors(t, vaePath, nil, map[string][]int{
		"encoder.down.0.block.0.conv1.weight": {128, 128, 3, 3},
		"decoder.up.3.block.0.conv1.weight":   {128, 128, 3, 3},
		"quant_conv.weight":                   {8, 8},
	})
	if got, ok := ProbeWeight(vaePath).AssetType(); !ok || got != AssetVAE {
		t.Fatalf("standalone VAE type = %q (ok=%v), want %q", got, ok, AssetVAE)
	}

	checkpoint := filepath.Join(dir, "merge.safetensors")
	writeSafetensors(t, checkpoint, nil, map[string][]int{
		"model.diffusion_model.first_stage_model.encoder.down.0.block.0.conv1.weight": {128, 128, 3, 3},
		"model.diffusion_model.input_blocks.1.0.weight":                               {320, 320},
	})
	if got, ok := ProbeWeight(checkpoint).AssetType(); !ok || got != AssetDiffusion {
		t.Fatalf("checkpoint with embedded VAE type = %q (ok=%v), want %q", got, ok, AssetDiffusion)
	}
}

// A 3.3 GB SigLIP text encoder used to be listed as a diffusion model purely
// because of its size.
func TestProbeClassifiesTextEncodersBeforeSize(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.json"),
		[]byte(`{"architectures":["SiglipModel"],"_class_name":"SiglipModel"}`), 0644); err != nil {
		t.Fatalf("config: %v", err)
	}
	path := filepath.Join(dir, "model.safetensors")
	writeSafetensors(t, path, nil, map[string][]int{
		"text_model.encoder.layers.0.mlp.fc1.weight":     {3072, 768},
		"vision_model.embeddings.patch_embedding.weight": {1152, 3, 14, 14},
	})
	probe := ProbeWeight(path)
	if got, ok := probe.AssetType(); !ok || got != AssetLLM {
		t.Fatalf("siglip type = %q (ok=%v), want %q", got, ok, AssetLLM)
	}

	clip := filepath.Join(t.TempDir(), "clip_g.safetensors")
	writeSafetensors(t, clip, nil, map[string][]int{
		"text_model.encoder.layers.0.mlp.fc1.weight": {5120, 1280},
		"logit_scale": {1},
	})
	if got, ok := ProbeWeight(clip).AssetType(); !ok || got != AssetClipG {
		t.Fatalf("clip type = %q (ok=%v), want %q", got, ok, AssetClipG)
	}
}

// An SD3 transformer is told apart from Flux by joint blocks, and SD3.5 by the
// wider cross-attention.
func TestProbeIdentifiesSD3AndSD35(t *testing.T) {
	dir := t.TempDir()
	medium := filepath.Join(dir, "sd3.safetensors")
	writeSafetensors(t, medium, nil, map[string][]int{
		"joint_blocks.0.attn2.to_k.weight": {1536, 1536},
	})
	if got, _ := ProbeWeight(medium).Family(); got != FamilySD3 {
		t.Fatalf("sd3 medium family = %q, want %q", got, FamilySD3)
	}

	large := filepath.Join(dir, "sd35.safetensors")
	writeSafetensors(t, large, nil, map[string][]int{
		"joint_blocks.0.attn2.to_k.weight": {2432, 2432},
	})
	if got, _ := ProbeWeight(large).Family(); got != FamilySD35 {
		t.Fatalf("sd3.5 large family = %q, want %q", got, FamilySD35)
	}
}

// A Diffusers folder (the shape Invoke AI uses for pipelines) is one model, not
// a pile of unet/vae/text_encoder fragments.
func TestScanTreatsDiffusersFolderAsOneModel(t *testing.T) {
	root := t.TempDir()
	bundle := filepath.Join(root, "4539fc37-98c8-4a9e-8bd6-4a3299bf1c85")
	for _, sub := range []string{"unet", "vae", "text_encoder", "text_encoder_2", "scheduler"} {
		if err := os.MkdirAll(filepath.Join(bundle, sub), 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
	}
	if err := os.WriteFile(filepath.Join(bundle, "model_index.json"),
		[]byte(`{"_class_name":"StableDiffusionXLPipeline"}`), 0644); err != nil {
		t.Fatalf("model_index: %v", err)
	}
	writeSafetensors(t, filepath.Join(bundle, "unet", "diffusion_pytorch_model.fp16.safetensors"),
		nil, map[string][]int{"input_blocks.1.0.weight": {320, 320}, "label_emb.0.0.weight": {2816, 1280}})
	writeSafetensors(t, filepath.Join(bundle, "vae", "diffusion_pytorch_model.fp16.safetensors"),
		nil, map[string][]int{"encoder.down.0.block.0.conv1.weight": {128, 128, 3, 3}})
	writeSafetensors(t, filepath.Join(bundle, "text_encoder", "model.fp16.safetensors"),
		nil, map[string][]int{"text_model.encoder.layers.0.mlp.fc1.weight": {3072, 768}})
	writeSafetensors(t, filepath.Join(bundle, "text_encoder_2", "model.fp16.safetensors"),
		nil, map[string][]int{"text_model.encoder.layers.0.mlp.fc1.weight": {5120, 1280}, "logit_scale": {1}})

	idx, err := New(root).Scan()
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(idx.Models) != 1 {
		t.Fatalf("diffusers bundle should be one model, got %d: %+v", len(idx.Models), idx.Models)
	}
	m := idx.Models[0]
	if m.Format != "diffusers" {
		t.Fatalf("format = %q, want diffusers", m.Format)
	}
	if m.Family != FamilySDXL {
		t.Fatalf("family = %q, want %q", m.Family, FamilySDXL)
	}
	if m.Files["model"] == "" || m.Files["vae"] == "" || m.Files["clip_l"] == "" || m.Files["clip_g"] == "" {
		t.Fatalf("components not mapped: %+v", m.Files)
	}
	if len(idx.VAEs) != 0 || len(idx.ClipLs) != 0 {
		t.Fatalf("bundle components leaked into their own buckets: vaes=%d clips=%d", len(idx.VAEs), len(idx.ClipLs))
	}
}

// A single file inside a plain folder must stay a file entry, not become a
// one-file "bundle" — that is what hid its real format behind "bundle".
func TestScanKeepsSingleFileOutOfBundlePath(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "7bf31271-ea16-4d04-9fff-dd9e748b64a2")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	writeSafetensors(t, filepath.Join(dir, "diffusion_pytorch_model.fp16.safetensors"),
		nil, map[string][]int{"input_blocks.1.0.weight": {320, 320}, "label_emb.0.0.weight": {2816, 1280}})

	idx, err := New(root).Scan()
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(idx.Models) != 1 {
		t.Fatalf("got %d models, want 1", len(idx.Models))
	}
	if got := idx.Models[0].Format; got != "safetensors" {
		t.Fatalf("format = %q, want safetensors", got)
	}
	if got := idx.Models[0].Family; got != FamilySDXL {
		t.Fatalf("family = %q, want %q", got, FamilySDXL)
	}
}

// A real library contained a .gguf whose header declares an impossible string
// length. gguf-parser-go allocated from it and panicked, the Add Folder binding
// never returned, and the Image Generator sat on "Scanning…" for the whole
// session. One unreadable file must cost that file's probe, not the scan.
func TestProbeWeightSurvivesMalformedGGUFHeader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "corrupt.gguf")
	var buf bytes.Buffer
	buf.WriteString("GGUF")
	binary.Write(&buf, binary.LittleEndian, uint32(3))     // version
	binary.Write(&buf, binary.LittleEndian, uint64(0))     // tensor count
	binary.Write(&buf, binary.LittleEndian, uint64(1))     // metadata count
	binary.Write(&buf, binary.LittleEndian, uint64(1<<62)) // key length: impossible
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("ProbeWeight panicked on a malformed header: %v", r)
		}
	}()
	p := ProbeWeight(path)
	if p.TensorCount != 0 || len(p.Metadata) != 0 {
		t.Fatalf("malformed header produced probe data: %+v", p)
	}
}
