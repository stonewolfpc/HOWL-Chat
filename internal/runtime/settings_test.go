package runtime

import (
	"os"
	"path/filepath"
	"testing"

	appsettings "howl-chat/internal/settings"
)

func newTestSettings(t *testing.T) (*Settings, string) {
	t.Helper()
	dir := t.TempDir()
	return NewSettings(appsettings.NewStore(dir), dir, map[string]interface{}{"temperature": 0.7, "top_k": 40}), dir
}

func TestNewSettingsCopiesDefaultsAndStartsEmpty(t *testing.T) {
	defaults := map[string]interface{}{"temperature": 0.7}
	s := NewSettings(appsettings.NewStore(t.TempDir()), t.TempDir(), defaults)

	defaults["temperature"] = 99.0
	if got := s.Sampler()["temperature"]; got != 0.7 {
		t.Fatalf("sampler temperature = %v, want the copied default 0.7", got)
	}
	if got := s.Model(); len(got) != 0 {
		t.Fatalf("model section = %v, want empty", got)
	}
	if got := s.All(); len(got) != 0 {
		t.Fatalf("document = %v, want empty", got)
	}
}

func TestSettingsLoadRestoresDocumentAndSections(t *testing.T) {
	s, dir := newTestSettings(t)
	doc := map[string]interface{}{
		"general": map[string]interface{}{"theme": "dark"},
		"audio":   map[string]interface{}{"tts_enabled": true},
		"sampler": map[string]interface{}{"temperature": 1.2, "top_p": 0.9},
		"model":   map[string]interface{}{"context_size": float64(8192)},
	}
	if err := appsettings.NewStore(dir).Save(doc); err != nil {
		t.Fatalf("seed store: %v", err)
	}

	if err := s.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}

	sampler := s.Sampler()
	if sampler["temperature"] != 1.2 || sampler["top_p"] != 0.9 {
		t.Fatalf("sampler = %v, want the document's section", sampler)
	}
	if sampler["top_k"] != 40 {
		t.Fatalf("sampler = %v, want the default top_k preserved", sampler)
	}
	if got := s.Model()["context_size"]; got != float64(8192) {
		t.Fatalf("model context_size = %v, want 8192", got)
	}
	if got := s.Section("audio")["tts_enabled"]; got != true {
		t.Fatalf("audio section = %v, want the document's section", s.Section("audio"))
	}
	if got := s.All()["general"].(map[string]interface{})["theme"]; got != "dark" {
		t.Fatalf("document general = %v", s.All()["general"])
	}
}

func TestSettingsSnapshotsAreIsolatedFromTheAuthority(t *testing.T) {
	s, _ := newTestSettings(t)
	s.SetSampler(map[string]interface{}{"temperature": 0.5})
	s.SetModel(map[string]interface{}{"context_size": 4096})
	s.Adopt(map[string]interface{}{"sampler": map[string]interface{}{"temperature": 0.5}})

	sampler := s.Sampler()
	sampler["temperature"] = 2.0
	if got := s.Sampler()["temperature"]; got != 0.5 {
		t.Fatalf("mutating a sampler snapshot changed the authority: %v", got)
	}

	model := s.Model()
	model["context_size"] = 1
	if got := s.Model()["context_size"]; got != 4096 {
		t.Fatalf("mutating a model snapshot changed the authority: %v", got)
	}

	doc := s.All()
	doc["general"] = "changed"
	if _, ok := s.All()["general"]; ok {
		t.Fatalf("mutating a document snapshot changed the authority: %v", s.All())
	}

	values := map[string]interface{}{"top_k": 1}
	s.SetSampler(values)
	values["top_k"] = 2
	if got := s.Sampler()["top_k"]; got != 1 {
		t.Fatalf("SetSampler kept the caller's map: %v", got)
	}
}

func TestSettingsLoadFallsBackToLegacySectionFiles(t *testing.T) {
	s, dir := newTestSettings(t)
	if err := os.WriteFile(filepath.Join(dir, samplerSettingsFile), []byte(`{"temperature":0.3}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, modelSettingsFile), []byte(`{"gpu_layers":12}`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := s.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := s.Sampler()["temperature"]; got != 0.3 {
		t.Fatalf("sampler temperature = %v, want 0.3 from the legacy file", got)
	}
	// The legacy file is a complete section: it replaces the defaults, exactly
	// as the pre-unified boot did, so a key it omits is not invented.
	if _, exists := s.Sampler()["top_k"]; exists {
		t.Fatalf("sampler top_k was re-added by defaults: %v", s.Sampler())
	}
	if got := s.Model()["gpu_layers"]; got != float64(12) {
		t.Fatalf("model gpu_layers = %v, want 12 from the legacy file", got)
	}
}

func TestSettingsLoadKeepsGoingAfterOneBadLegacyFile(t *testing.T) {
	s, dir := newTestSettings(t)
	if err := os.WriteFile(filepath.Join(dir, samplerSettingsFile), []byte(`{not json`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, modelSettingsFile), []byte(`{"batch_size":256}`), 0644); err != nil {
		t.Fatal(err)
	}

	err := s.Load()
	if err == nil {
		t.Fatal("Load reported success for a corrupt sampler file")
	}
	if got := s.Sampler()["top_k"]; got != 40 {
		t.Fatalf("sampler top_k = %v, want the defaults kept after a corrupt file", got)
	}
	if got := s.Model()["batch_size"]; got != float64(256) {
		t.Fatalf("model batch_size = %v, want the healthy file still restored", got)
	}
}

func TestSettingsReplacePersistsAndReplacesSections(t *testing.T) {
	s, dir := newTestSettings(t)
	s.SetSampler(map[string]interface{}{"temperature": 0.7, "stale_key": true})
	s.SetModel(map[string]interface{}{"gpu_layers": 8})

	doc := map[string]interface{}{
		"general": map[string]interface{}{"theme": "light"},
		"sampler": map[string]interface{}{"temperature": 1.5},
	}
	if err := s.Replace(doc); err != nil {
		t.Fatalf("Replace: %v", err)
	}

	sampler := s.Sampler()
	if sampler["temperature"] != 1.5 {
		t.Fatalf("sampler temperature = %v, want 1.5", sampler["temperature"])
	}
	if _, exists := sampler["stale_key"]; exists {
		t.Fatalf("sampler kept a key the caller removed: %v", sampler)
	}
	// A section the caller did not send is left alone, as before.
	if got := s.Model()["gpu_layers"]; got != 8 {
		t.Fatalf("model gpu_layers = %v, want the previous section preserved", got)
	}

	persisted, err := appsettings.NewStore(dir).Load()
	if err != nil {
		t.Fatalf("reload store: %v", err)
	}
	if got := persisted["general"].(map[string]interface{})["theme"]; got != "light" {
		t.Fatalf("persisted document = %v", persisted)
	}

	if err := s.Replace(nil); err == nil {
		t.Fatal("Replace accepted a nil document")
	}
}

func TestSettingsReplaceReportsAnUnwritableStore(t *testing.T) {
	root := t.TempDir()
	blocker := filepath.Join(root, "blocked")
	if err := os.WriteFile(blocker, []byte("not a directory"), 0644); err != nil {
		t.Fatal(err)
	}
	s := NewSettings(appsettings.NewStore(filepath.Join(blocker, "config")), filepath.Join(blocker, "config"), nil)
	if err := s.Replace(map[string]interface{}{"general": map[string]interface{}{}}); err == nil {
		t.Fatal("Replace reported success although the store cannot be written")
	}
}

func TestSettingsMergesCreateAndKeepSections(t *testing.T) {
	var uninitialized *Settings
	uninitialized.MergeModel(map[string]interface{}{"a": 1})
	uninitialized.SetSampler(map[string]interface{}{"a": 1})
	if len(uninitialized.Model()) != 0 || len(uninitialized.Sampler()) != 0 {
		t.Fatal("a nil authority grew a section")
	}
	if err := uninitialized.Replace(map[string]interface{}{}); err == nil {
		t.Fatal("a nil authority accepted a document")
	}
	if err := uninitialized.Load(); err == nil {
		t.Fatal("a nil authority reported a successful load")
	}
	if err := uninitialized.SaveModelFile(); err != nil {
		t.Fatalf("nil SaveModelFile = %v, want nil", err)
	}
	if err := uninitialized.LoadSamplerFile(); err != nil {
		t.Fatalf("nil LoadSamplerFile = %v, want nil", err)
	}

	s, _ := newTestSettings(t)
	s.MergeModel(map[string]interface{}{"gpu_layers": 4})
	s.MergeModel(map[string]interface{}{"context_size": 2048})
	model := s.Model()
	if model["gpu_layers"] != 4 || model["context_size"] != 2048 {
		t.Fatalf("merged model = %v, want both keys", model)
	}

	s.SetModel(map[string]interface{}{"only": true})
	if _, exists := s.Model()["gpu_layers"]; exists {
		t.Fatalf("SetModel merged instead of replacing: %v", s.Model())
	}
}

func TestSettingsLegacyFilesRoundTripAndReportCorruption(t *testing.T) {
	s, dir := newTestSettings(t)
	s.SetSampler(map[string]interface{}{"temperature": 0.25})
	s.SetModel(map[string]interface{}{"threads": 6})

	if err := s.SaveSamplerFile(); err != nil {
		t.Fatalf("SaveSamplerFile: %v", err)
	}
	if err := s.SaveModelFile(); err != nil {
		t.Fatalf("SaveModelFile: %v", err)
	}

	reloaded := NewSettings(appsettings.NewStore(dir), dir, nil)
	if err := reloaded.LoadSamplerFile(); err != nil {
		t.Fatalf("LoadSamplerFile: %v", err)
	}
	if err := reloaded.LoadModelFile(); err != nil {
		t.Fatalf("LoadModelFile: %v", err)
	}
	if got := reloaded.Sampler()["temperature"]; got != 0.25 {
		t.Fatalf("sampler temperature = %v", got)
	}
	if got := reloaded.Model()["threads"]; got != float64(6) {
		t.Fatalf("model threads = %v", got)
	}

	fresh := NewSettings(appsettings.NewStore(t.TempDir()), t.TempDir(), nil)
	if err := fresh.LoadSamplerFile(); err != nil {
		t.Fatalf("missing sampler file = %v, want nil", err)
	}
	if err := fresh.LoadModelFile(); err != nil {
		t.Fatalf("missing model file = %v, want nil", err)
	}

	if err := os.WriteFile(filepath.Join(dir, modelSettingsFile), []byte(`[]`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := reloaded.LoadModelFile(); err == nil {
		t.Fatal("LoadModelFile accepted a non-object file")
	}
}
