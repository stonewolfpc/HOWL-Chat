package modelintel

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClassifySeparatesCapabilityAndRuntime(t *testing.T) {
	chat := Classify(Input{Path: "D:/models/qwen.gguf", Type: "llm", Family: "qwen2", Arch: "qwen2", Format: "gguf"})
	if !chat.Supports(CapabilityChat) || chat.Supports(CapabilityImage) {
		t.Fatalf("chat capabilities = %+v", chat)
	}
	if len(chat.Runtimes) != 1 || chat.Runtimes[0] != RuntimeLlamaCPP {
		t.Fatalf("chat runtimes = %+v", chat.Runtimes)
	}
	if chat.Evidence != EvidenceVerified {
		t.Fatalf("chat evidence = %q", chat.Evidence)
	}

	image := Classify(Input{Path: "D:/models/flux.gguf", Type: "diffusion", Family: "flux", Format: "gguf"})
	if !image.Supports(CapabilityImage) || image.Supports(CapabilityChat) {
		t.Fatalf("image capabilities = %+v", image)
	}
	if image.Evidence != EvidenceSourceDerived {
		t.Fatalf("image evidence = %q", image.Evidence)
	}
}

func TestCacheInvalidatesOnFileFingerprintChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "model.gguf")
	if err := os.WriteFile(path, []byte("one"), 0644); err != nil {
		t.Fatal(err)
	}
	cache, err := NewCache(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	profile := Classify(Input{Path: path, Type: "llm", Arch: "llama", SizeBytes: 3, ModTimeNS: 10})
	if err := cache.Put(profile); err != nil {
		t.Fatal(err)
	}
	if _, ok := cache.Get(path, 3, 10, ""); !ok {
		t.Fatal("cached profile was not reused")
	}
	if _, ok := cache.Get(path, 4, 10, ""); ok {
		t.Fatal("cache accepted a changed file size")
	}
}

func TestInspectArtifactsReadsSidecarsAndInvalidatesFingerprint(t *testing.T) {
	dir := t.TempDir()
	model := filepath.Join(dir, "model.gguf")
	if err := os.WriteFile(model, []byte("weights"), 0644); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(dir, "config.json")
	if err := os.WriteFile(config, []byte(`{"architectures":["LlamaForCausalLM"],"model_type":"llama"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "tokenizer_config.json"), []byte(`{"tokenizer_class":"LlamaTokenizer","chat_template":"{{ messages }}"}`), 0644); err != nil {
		t.Fatal(err)
	}
	first := InspectArtifacts(model)
	if first.Architecture != "LlamaForCausalLM" || first.Family != "llama" || first.Tokenizer != "LlamaTokenizer" || first.Fingerprint == "" {
		t.Fatalf("artifacts = %+v", first)
	}
	profile := Classify(Input{Path: model, Type: "llm"})
	if profile.Architecture != "LlamaForCausalLM" || profile.Evidence != EvidenceSourceDerived {
		t.Fatalf("profile = %+v", profile)
	}
	if err := os.WriteFile(config, []byte(`{"architectures":["MistralForCausalLM"],"model_type":"mistral"}`), 0644); err != nil {
		t.Fatal(err)
	}
	second := InspectArtifacts(model)
	if second.Fingerprint == first.Fingerprint {
		t.Fatal("sidecar change did not invalidate fingerprint")
	}
	cache, err := NewCache(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := cache.Put(profile); err != nil {
		t.Fatal(err)
	}
	if _, ok := cache.Get(model, profile.SizeBytes, profile.ModTimeNS, first.Fingerprint); !ok {
		t.Fatal("matching fingerprint was not reused")
	}
	if _, ok := cache.Get(model, profile.SizeBytes, profile.ModTimeNS, second.Fingerprint); ok {
		t.Fatal("changed metadata fingerprint reused stale profile")
	}
}

func TestClassifyConflictingFilenameAndType(t *testing.T) {
	profile := Classify(Input{Path: "D:/models/flux.gguf", Type: "llm", Family: "llama", Arch: "llama", Format: "gguf"})
	if profile.Evidence != EvidenceConflicting {
		t.Fatalf("evidence = %q, want conflicting", profile.Evidence)
	}
}
