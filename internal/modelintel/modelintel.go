// Package modelintel provides shared, cached model identification metadata.
// It does not load models or replace the existing runtime registries.
package modelintel

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type Evidence string

const (
	EvidenceVerified      Evidence = "verified"
	EvidenceSourceDerived Evidence = "source-derived"
	EvidenceInferred      Evidence = "inferred"
	EvidenceUnknown       Evidence = "unknown"
	EvidenceConflicting   Evidence = "conflicting"
)

type Capability string

const (
	CapabilityChat      Capability = "chat"
	CapabilityImage     Capability = "image"
	CapabilityVideo     Capability = "video"
	CapabilityTTS       Capability = "tts"
	CapabilityASR       Capability = "asr"
	CapabilityEmbedding Capability = "embedding"
	CapabilityUnknown   Capability = "unknown"
)

type Runtime string

const (
	RuntimeLlamaCPP Runtime = "llama.cpp"
	RuntimeSD       Runtime = "sd-server"
	RuntimeSDCLI    Runtime = "sd-cli"
	RuntimePiper    Runtime = "piper"
	RuntimeKokoro   Runtime = "kokoro"
	RuntimeUnknown  Runtime = "unknown"
)

type Input struct {
	Path      string
	Name      string
	Type      string
	Family    string
	Format    string
	Arch      string
	SizeBytes int64
	ModTimeNS int64
}

type Profile struct {
	Path         string       `json:"path"`
	Name         string       `json:"name"`
	Format       string       `json:"format"`
	Architecture string       `json:"architecture,omitempty"`
	Family       string       `json:"family,omitempty"`
	Tokenizer    string       `json:"tokenizer,omitempty"`
	ChatTemplate string       `json:"chat_template,omitempty"`
	Fingerprint  string       `json:"fingerprint,omitempty"`
	Capabilities []Capability `json:"capabilities"`
	Runtimes     []Runtime    `json:"runtimes"`
	Evidence     Evidence     `json:"evidence"`
	Source       string       `json:"source"`
	Notes        []string     `json:"notes,omitempty"`
	Available    bool         `json:"available"`
	SizeBytes    int64        `json:"size_bytes"`
	ModTimeNS    int64        `json:"mod_time_ns"`
	UpdatedAt    string       `json:"updated_at"`
}

// Artifacts is bounded metadata found beside a model. It never reads model
// weights; fingerprint changes only when a small sidecar file changes.
type Artifacts struct {
	Architecture string
	Family       string
	Tokenizer    string
	ChatTemplate string
	Fingerprint  string
	Source       string
}

const maxMetadataBytes = 1 << 20

// InspectArtifacts reads config/tokenizer metadata beside a model. Ordinary
// discovery never calls this on multi-gigabyte weights; it only reads small
// JSON files and caches a content fingerprint.
func InspectArtifacts(path string) Artifacts {
	a := Artifacts{Source: "sidecar metadata"}
	dir := filepath.Dir(path)
	h := sha256.New()
	hasSidecar := false
	for _, name := range []string{"config.json", "tokenizer_config.json", "generation_config.json", "model_index.json"} {
		p := filepath.Join(dir, name)
		info, err := os.Stat(p)
		if err != nil || info.IsDir() || info.Size() > maxMetadataBytes {
			continue
		}
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		h.Write([]byte(name))
		h.Write([]byte{0})
		h.Write(data)
		hasSidecar = true
		a.inspectFile(name, data)
	}
	if hasSidecar {
		// A hash of only the files that actually exist is stable. Cache callers
		// can tell "no sidecars" from "sidecars present".
		a.Fingerprint = hex.EncodeToString(h.Sum(nil))
	}
	return a
}

func (a *Artifacts) inspectFile(name string, data []byte) {
	var raw map[string]interface{}
	if json.Unmarshal(data, &raw) != nil {
		return
	}
	switch name {
	case "config.json":
		if a.Architecture == "" {
			a.Architecture = firstMetadataString(raw, "architectures")
		}
		if a.Family == "" {
			a.Family = firstMetadataString(raw, "model_type", "_class_name")
		}
	case "model_index.json":
		if a.Family == "" {
			a.Family = firstMetadataString(raw, "_class_name")
		}
	case "tokenizer_config.json":
		if a.Tokenizer == "" {
			a.Tokenizer = firstMetadataString(raw, "tokenizer_class", "name_or_path")
		}
		if a.ChatTemplate == "" {
			a.ChatTemplate = firstMetadataString(raw, "chat_template")
		}
	case "generation_config.json":
		if a.Tokenizer == "" {
			a.Tokenizer = firstMetadataString(raw, "tokenizer_class")
		}
	}
}

func firstMetadataString(raw map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		value := raw[key]
		switch v := value.(type) {
		case string:
			if v != "" {
				return v
			}
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok && s != "" {
					return s
				}
			}
		}
	}
	return ""
}

func Classify(in Input) Profile {
	return ClassifyWithArtifacts(in, InspectArtifacts(in.Path))
}

func ClassifyWithArtifacts(in Input, art Artifacts) Profile {
	profile := Profile{
		Path:         filepath.Clean(in.Path),
		Name:         in.Name,
		Format:       in.Format,
		Architecture: in.Arch,
		Family:       in.Family,
		Tokenizer:    art.Tokenizer,
		ChatTemplate: art.ChatTemplate,
		Fingerprint:  art.Fingerprint,
		Evidence:     EvidenceInferred,
		Source:       "local filesystem",
		SizeBytes:    in.SizeBytes,
		ModTimeNS:    in.ModTimeNS,
		UpdatedAt:    time.Now().UTC().Format(time.RFC3339),
	}
	if profile.Name == "" {
		profile.Name = filepath.Base(in.Path)
	}
	if in.Arch != "" {
		profile.Evidence = EvidenceVerified
		profile.Source = "model metadata"
	}
	if in.Family != "" && in.Arch == "" {
		profile.Evidence = EvidenceSourceDerived
	}
	if in.Arch == "" && art.Architecture != "" {
		profile.Architecture = art.Architecture
		profile.Evidence = EvidenceSourceDerived
		profile.Source = art.Source
	}
	if in.Family == "" && art.Family != "" {
		profile.Family = art.Family
		if profile.Evidence != EvidenceVerified {
			profile.Evidence = EvidenceSourceDerived
		}
		profile.Source = art.Source
	}
	if art.Tokenizer != "" {
		profile.Notes = append(profile.Notes, "Tokenizer metadata: "+art.Tokenizer)
	}
	if art.ChatTemplate != "" {
		profile.Notes = append(profile.Notes, "Chat template metadata present")
	}

	switch in.Type {
	case "llm":
		profile.Capabilities = append(profile.Capabilities, CapabilityChat)
		profile.Runtimes = append(profile.Runtimes, RuntimeLlamaCPP)
	case "tts":
		profile.Capabilities = append(profile.Capabilities, CapabilityTTS)
		profile.Runtimes = append(profile.Runtimes, RuntimePiper)
	case "diffusion", "vae", "lora", "controlnet", "embedding", "upscaler",
		"t5xxl", "clip_l", "clip_g":
		profile.Capabilities = append(profile.Capabilities, CapabilityImage)
		profile.Runtimes = append(profile.Runtimes, RuntimeSD, RuntimeSDCLI)
	case "asr":
		profile.Capabilities = append(profile.Capabilities, CapabilityASR)
		profile.Runtimes = append(profile.Runtimes, RuntimeUnknown)
	default:
		profile.Capabilities = append(profile.Capabilities, CapabilityUnknown)
		profile.Runtimes = append(profile.Runtimes, RuntimeUnknown)
		profile.Evidence = EvidenceUnknown
	}

	lowerName := filepath.Base(in.Path)
	lowerName = stringLower(lowerName)
	if (contains(lowerName, "flux") || contains(lowerName, "sdxl") || contains(lowerName, "sd3")) &&
		!hasCapability(profile.Capabilities, CapabilityImage) {
		profile.Evidence = EvidenceConflicting
		profile.Notes = append(profile.Notes, "filename suggests an image model, but the catalog type does not")
	}
	if in.Arch != "" && in.Family != "" && !stringsEqualFold(in.Arch, in.Family) {
		profile.Evidence = EvidenceConflicting
		profile.Notes = append(profile.Notes, "architecture and declared family disagree")
	}
	if in.Arch != "" && art.Architecture != "" && !stringsEqualFold(in.Arch, art.Architecture) {
		profile.Evidence = EvidenceConflicting
		profile.Notes = append(profile.Notes, "sidecar architecture disagrees with model metadata")
	}
	return profile
}

func (p Profile) Supports(capability Capability) bool {
	return hasCapability(p.Capabilities, capability)
}

type Cache struct {
	mu      sync.RWMutex
	path    string
	entries map[string]Profile
}

type persistedCache struct {
	Version int       `json:"version"`
	Models  []Profile `json:"models"`
}

func NewCache(configDir string) (*Cache, error) {
	c := &Cache{
		path:    filepath.Join(configDir, "model-intelligence.json"),
		entries: map[string]Profile{},
	}
	data, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return c, nil
		}
		return nil, fmt.Errorf("read model intelligence cache: %w", err)
	}
	var stored persistedCache
	if err := json.Unmarshal(data, &stored); err != nil {
		return nil, fmt.Errorf("parse model intelligence cache: %w", err)
	}
	for _, profile := range stored.Models {
		if profile.Path == "" {
			continue
		}
		c.entries[cacheKey(profile.Path)] = profile
	}
	return c, nil
}

func (c *Cache) Get(path string, size, modNS int64, fingerprint string) (Profile, bool) {
	c.mu.RLock()
	profile, ok := c.entries[cacheKey(path)]
	c.mu.RUnlock()
	if !ok || profile.SizeBytes != size || profile.ModTimeNS != modNS || profile.Fingerprint != fingerprint {
		return Profile{}, false
	}
	return profile, true
}

func (c *Cache) Put(profile Profile) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if profile.Path == "" {
		return fmt.Errorf("model intelligence profile has no path")
	}
	if profile.UpdatedAt == "" {
		profile.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	c.entries[cacheKey(profile.Path)] = profile
	return c.saveLocked()
}

func (c *Cache) All() []Profile {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Profile, 0, len(c.entries))
	for _, profile := range c.entries {
		out = append(out, profile)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

func (c *Cache) saveLocked() error {
	list := make([]Profile, 0, len(c.entries))
	for _, profile := range c.entries {
		list = append(list, profile)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Path < list[j].Path })
	data, err := json.MarshalIndent(persistedCache{Version: 1, Models: list}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode model intelligence cache: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(c.path), 0755); err != nil {
		return fmt.Errorf("create model intelligence config directory: %w", err)
	}
	tmp := c.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("write model intelligence cache: %w", err)
	}
	if err := os.Rename(tmp, c.path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("commit model intelligence cache: %w", err)
	}
	return nil
}

func cacheKey(path string) string {
	return stringLower(filepath.Clean(path))
}

func hasCapability(values []Capability, want Capability) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func contains(value, part string) bool {
	return len(part) > 0 && len(value) >= len(part) && indexOf(value, part) >= 0
}

func indexOf(value, part string) int {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return i
		}
	}
	return -1
}

func stringLower(value string) string {
	out := []byte(value)
	for i, b := range out {
		if b >= 'A' && b <= 'Z' {
			out[i] = b + ('a' - 'A')
		}
	}
	return string(out)
}

func stringsEqualFold(a, b string) bool {
	return stringLower(a) == stringLower(b)
}
