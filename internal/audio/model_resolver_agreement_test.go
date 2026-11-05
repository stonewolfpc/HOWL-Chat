package audio

import (
	"errors"
	"testing"

	"howl-chat/internal/audio/types"
)

// DEBT-T3: one declaration decides what may be resolved. If the registry and
// the resolver disagree, the UI will offer something the pipeline rejects.
func TestResolverApprovedSetDerivesFromRegistry(t *testing.T) {
	want := make(map[types.ASRModelType]bool)
	for modelType, profile := range types.ASRModelRegistry {
		want[modelType] = profile.Local && profile.ApprovedLocal && profile.LocalFile != ""
	}

	got := make(map[types.ASRModelType]bool)
	for _, model := range approvedASRModels {
		got[model.Type] = true
		if file, ok := approvedASRModelFile(model.Type); !ok || file != model.File {
			t.Fatalf("%q resolves to %q, %t", model.Type, file, ok)
		}
	}

	for modelType, approved := range want {
		if got[modelType] != approved {
			t.Fatalf("%q approved=%t in resolver, want %t", modelType, got[modelType], approved)
		}
	}
}

func TestResolverRejectsUnadmittedAndRemoteEngines(t *testing.T) {
	resolver := mustResolver(t, ResolverConfig{Root: t.TempDir(), Getenv: envMap(nil)})

	for _, engine := range []string{
		string(types.ASRQwen3ASR06B), // remote API
		string(types.ASRQwen3ASR17B), // remote API
		"whisper-medium",             // declared identity, no registry entry
	} {
		if _, err := resolver.Resolve(engine, nil); !errors.Is(err, ErrEngineUnsupported) {
			t.Fatalf("Resolve(%q) error = %v, want ErrEngineUnsupported", engine, err)
		}
	}

	// The admitted local Qwen is supported, but nothing is installed in this
	// empty root.
	if _, err := resolver.Resolve(string(types.ASRQwen3ASR17BLocal), nil); !errors.Is(err, ErrModelNotFound) {
		t.Fatalf("admitted local Qwen error = %v, want ErrModelNotFound", err)
	}
}
